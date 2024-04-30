package agent

import (
	"flag"
	"github.com/mailru/easyjson"
	"github.com/psfpro/metrics/pkg/config"
	"log"
	"os"
)

//easyjson:json
type Config struct {
	ServerAddress  string `json:"address"`
	ReportInterval int64  `json:"report_interval"`
	PollInterval   int64  `json:"poll_interval"`
	RateLimit      int64  `json:"rate_limit"`
	HashKey        string `json:"hash_key"`
	CryptoKey      string `json:"crypto_key"`
}

func NewConfig() *Config {
	var cfg Config
	file := flag.String("cfg", "config-agent.json", "Configuration file")
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
	config.StringVar(&cfg.ServerAddress, "ADDRESS", "a", "Net serverAddress host:port")
	config.Int64Var(&cfg.ReportInterval, "REPORT_INTERVAL", "r", "frequency of sending metrics to the server")
	config.Int64Var(&cfg.PollInterval, "POLL_INTERVAL", "p", "metrics polling rate")
	config.Int64Var(&cfg.RateLimit, "RATE_LIMIT", "l", "rete limit")
	config.StringVar(&cfg.HashKey, "KEY", "k", "Hash key")
	config.StringVar(&cfg.CryptoKey, "CRYPTO_KEY", "crypto-key", "crypto key")

	return &cfg
}
