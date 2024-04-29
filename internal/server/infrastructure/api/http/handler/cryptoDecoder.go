package handler

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type CryptoDecoder struct {
	cryptoKey string
	key       *rsa.PrivateKey
}

func NewCryptoDecoder(cryptoKey string) *CryptoDecoder {
	return &CryptoDecoder{cryptoKey: cryptoKey}
}

func (h *CryptoDecoder) Decode(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if h.cryptoKey != "" {
			bodyBytes, _ := io.ReadAll(r.Body)
			privateKey, err := h.privateKey()
			if err != nil {
				panic(err)
			}
			encryptedData := bodyBytes[:len(bodyBytes)-256-12]
			encryptedKey := bodyBytes[len(bodyBytes)-256-12 : len(bodyBytes)-12]
			nonce := bodyBytes[len(bodyBytes)-12:]
			aesKey, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedKey)

			// Decrypt the data using AES-GCM
			blockCipher, err := aes.NewCipher(aesKey)
			if err != nil {
				fmt.Println("Failed to create cipher:", err)
				return
			}

			gcm, err := cipher.NewGCM(blockCipher)
			if err != nil {
				fmt.Println("Failed to create GCM:", err)
				return
			}

			// Decrypt the data
			bodyBytes, err = gcm.Open(nil, nonce, encryptedData, nil)
			if err != nil {
				fmt.Println("Failed to decrypt data:", err)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (h *CryptoDecoder) privateKey() (*rsa.PrivateKey, error) {
	if h.key == nil {
		privateKey, err := loadRSAPrivateKey(h.cryptoKey)
		if err != nil {
			return nil, err
		}
		h.key = privateKey
	}

	return h.key, nil
}

func loadRSAPrivateKey(filePath string) (*rsa.PrivateKey, error) {
	// Read the file containing the private key
	keyBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Decode the PEM data
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, err
	}
	key := block.Bytes

	// Parse the private key
	privateKey, err := x509.ParsePKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}

	// Assert the type to *rsa.PrivateKey
	rsaPrivate, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not an RSA private key")
	}

	return rsaPrivate, nil
}
