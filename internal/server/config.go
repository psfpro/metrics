package server

import (
	"flag"
	"github.com/mailru/easyjson"
	"github.com/psfpro/metrics/pkg/config"
	"log"
	"os"
)

//easyjson:json
type Config struct {
	Address       string `json:"address"`
	Restore       bool   `json:"restore"`
	StoreInterval int64  `json:"store_interval"`
	StoragePath   string `json:"store_file"`
	DatabaseDsn   string `json:"database_dsn"`
	HashKey       string `json:"hash_key"`
	CryptoKey     string `json:"crypto_key"`
	TrustedSubnet string `json:"trusted_subnet"`
}

func NewConfig() *Config {
	var cfg Config
	file := flag.String("cfg", "config-server.json", "Configuration file")
	if envFile := os.Getenv("CONFIG"); envFile != "" {
		file = &envFile
	}
	reader, err := os.Open(*file)
	if err != nil {
		log.Println(err)
	} else {
		if err := easyjson.UnmarshalFromReader(reader, &cfg); err != nil {
			log.Println(err)
		}
	}
	config.StringVar(&cfg.Address, "ADDRESS", "a", "Net serverAddress host:port")
	config.Int64Var(&cfg.StoreInterval, "STORE_INTERVAL", "i", "Store interval")
	config.StringVar(&cfg.StoragePath, "FILE_STORAGE_PATH", "f", "File storage path")
	config.BoolVar(&cfg.Restore, "RESTORE", "r", "Restore")
	config.StringVar(&cfg.DatabaseDsn, "DATABASE_DSN", "d", "Database DSN")
	config.StringVar(&cfg.HashKey, "KEY", "k", "Hash key")
	config.StringVar(&cfg.CryptoKey, "CRYPTO_KEY", "crypto-key", "crypto key")
	config.StringVar(&cfg.TrustedSubnet, "TRUSTED_SUBNET", "t", "trusted subnet")

	return &cfg
}
