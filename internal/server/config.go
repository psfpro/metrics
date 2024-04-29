package server

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	serverAddress   *NetAddress
	storeInterval   int64
	fileStoragePath string
	restore         bool
	databaseDsn     *DatabaseDsn
	hashKey         string
	cryptoKey       string
}

func NewConfig() *Config {
	address := &NetAddress{}
	err := address.Set(":8080")
	if err != nil {
		panic(err)
	}
	databaseDsn := &DatabaseDsn{}
	_ = flag.Value(address)
	_ = flag.Value(databaseDsn)
	storeInterval := flag.Int64("i", 300, "Store interval")
	fileStoragePath := flag.String("f", "/tmp/metrics-db.json", "File storage path")
	restore := flag.Bool("r", true, "Restore")
	hashKey := flag.String("k", "", "Hash key")
	cryptoKey := flag.String("crypto-key", "", "crypto key")
	flag.Var(address, "a", "Net serverAddress host:port")
	flag.Var(databaseDsn, "d", "Database DSN")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		err = address.Set(envRunAddr)
		if err != nil {
			panic(err)
		}
	}
	if envDatabaseDsn := os.Getenv("DATABASE_DSN"); envDatabaseDsn != "" {
		err = databaseDsn.Set(envDatabaseDsn)
		if err != nil {
			panic(err)
		}
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		i, _ := strconv.ParseInt(envStoreInterval, 10, 64)
		storeInterval = &i
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		fileStoragePath = &envFileStoragePath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		r, _ := strconv.ParseBool(envRestore)
		restore = &r
	}
	if envHashKey := os.Getenv("KEY"); envHashKey != "" {
		hashKey = &envHashKey
	}
	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		cryptoKey = &envCryptoKey
	}

	return &Config{
		serverAddress:   address,
		storeInterval:   *storeInterval,
		fileStoragePath: *fileStoragePath,
		restore:         *restore,
		databaseDsn:     databaseDsn,
		hashKey:         *hashKey,
		cryptoKey:       *cryptoKey,
	}
}

type NetAddress struct {
	Host string
	Port int
}

func (a NetAddress) String() string {
	return a.Host + ":" + strconv.Itoa(a.Port)
}

func (a *NetAddress) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("need serverAddress in a form host:port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	a.Host = hp[0]
	a.Port = port
	return nil
}

type DatabaseDsn struct {
	value string
}

func (d DatabaseDsn) String() string {
	return d.value
}

func (d *DatabaseDsn) Set(s string) error {
	d.value = s

	return nil
}
