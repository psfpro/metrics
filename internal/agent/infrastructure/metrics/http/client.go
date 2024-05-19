package http

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/psfpro/metrics/internal/agent/model"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

type Client struct {
	serverAddress  string
	reportInterval int64
	hashKey        string
	cryptoKey      string
	wg             sync.WaitGroup
	ip             net.IP
}

func NewClient(serverAddress string, reportInterval int64, hashKey string, cryptoKey string) *Client {
	ip, _ := detectIPAddress()
	return &Client{serverAddress: serverAddress, reportInterval: reportInterval, hashKey: hashKey, cryptoKey: cryptoKey, ip: ip}
}

func (c *Client) Run(collectResults chan []model.Metrics, sendResults chan error, closed chan struct{}) {
	for w := 1; w <= 3; w++ {
		c.wg.Add(1)
		go c.send(w, collectResults, sendResults)
	}
	c.wg.Wait()
	close(sendResults)
	closed <- struct{}{} // можно завершать приложение gracefully
}

func (c *Client) sendBatchMetrics(batch []model.Metrics) error {
	var err error
	retryDelays := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

	for _, delay := range retryDelays {
		err = c.sendBatch(batch)
		if err == nil {
			return nil
		}

		time.Sleep(delay)
	}

	return fmt.Errorf("после нескольких попыток: %w", err)
}

func (c *Client) sendBatch(metric []model.Metrics) (err error) {
	reqBytes, err := json.Marshal(metric)
	if err != nil {
		fmt.Println(err)
		return err
	}
	urlString := fmt.Sprintf("%s/updates", c.serverAddress)
	var encryption string
	if c.cryptoKey != "" {
		rsaPublicKey, err := loadRSAPublicKey(c.cryptoKey)
		if err != nil {
			return err
		}
		// Generate a new AES key
		aesKey := make([]byte, aes.BlockSize)
		if _, err := rand.Read(aesKey); err != nil {
			return err
		}
		// Encrypt the data with AES
		blockCipher, err := aes.NewCipher(aesKey)
		if err != nil {
			return err
		}
		gcm, err := cipher.NewGCM(blockCipher)
		if err != nil {
			return err
		}
		nonce := make([]byte, gcm.NonceSize())
		if _, err := rand.Read(nonce); err != nil {
			return err
		}
		encryptedData := gcm.Seal(nil, nonce, reqBytes, nil)
		// Encrypt the AES key with RSA
		encryptedKey, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, aesKey)
		if err != nil {
			return err
		}
		reqBytes = append(encryptedData, encryptedKey...)
		reqBytes = append(reqBytes, nonce...)
		encryption = "rsa"
	}
	body, err := c.compress(reqBytes)
	if err != nil {
		return err
	}
	request, _ := http.NewRequest("POST", urlString, &body)
	if encryption != "" {
		request.Header.Set("X-Encryption", encryption)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Content-Encoding", "gzip")
	if c.ip != nil {
		request.Header.Set("X-Real-IP", c.ip.String())
	}
	if c.hashKey != "" {
		hash := hmac.New(sha256.New, []byte(c.hashKey))
		hash.Write(reqBytes)
		signature := hex.EncodeToString(hash.Sum(nil))
		request.Header.Set("HashSHA256", signature)
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		err = resp.Body.Close()
	}()
	return nil
}

func (c *Client) compress(data []byte) (bytes.Buffer, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)

	_, err := w.Write(data)
	if err != nil {
		return b, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}
	err = w.Close()
	if err != nil {
		return b, fmt.Errorf("failed compress data: %v", err)
	}

	return b, nil
}

func (c *Client) send(id int, jobs <-chan []model.Metrics, results chan<- error) {
	dataForSend := make(map[string]model.Metrics)
	ticker := time.NewTicker(time.Duration(c.reportInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case j, ok := <-jobs:
			if !ok {
				log.Printf("send %d stopping\n", id)
				c.wg.Done()
				return
			}
			log.Printf("send %d starting task\n", id)
			for _, value := range j {
				dataForSend[value.ID] = value
			}
		case <-ticker.C:
			if len(dataForSend) == 0 {
				continue
			}
			log.Printf("send %d performing action\n", id)
			var values []model.Metrics
			for _, value := range dataForSend {
				values = append(values, value)
			}
			err := c.sendBatchMetrics(values)
			if err == nil {
				dataForSend = make(map[string]model.Metrics)
			}

			results <- err
		}
	}
}

func loadRSAPublicKey(filename string) (*rsa.PublicKey, error) {
	// Read the file containing the public key
	keyData, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read public key file: %w", err)
	}

	// Decode the PEM data
	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	// Parse the public key
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse public key: %w", err)
	}

	// Assert the type to *rsa.PublicKey
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	return rsaPub, nil
}

func detectIPAddress() (net.IP, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println("Error getting network interfaces:", err)
		return nil, err
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Println("Error getting addresses for interface:", i.Name, err)
			continue
		}

		for _, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			// Check if the IP address is IPv4 (you can skip this check if you want both IPv4 and IPv6)
			if ip.To4() != nil {
				log.Printf("Interface: %s, IP: %s\n", i.Name, ip.String())
				return ip, nil
			}
		}
	}
	return nil, errors.New("no IP address found")
}
