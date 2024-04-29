package agent

import (
	"flag"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HashKey        string
	CryptoKey      string
	ServerAddress  string
	PollInterval   time.Duration
	ReportInterval time.Duration
	RateLimit      int
}

func NewConfig() *Config {
	hashKey := flag.String("k", "", "hash key")
	cryptoKey := flag.String("crypto-key", "", "crypto key")
	flagRunAddr := flag.String("a", ":8080", "address and port to run server")
	flagReportInterval := flag.Int("r", 10, "frequency of sending metrics to the server")
	flagPollInterval := flag.Int("p", 2, "metrics polling rate")
	rateLimit := flag.Int("l", 2, "rete limit")
	flag.Parse()
	if envHashKey := os.Getenv("KEY"); envHashKey != "" {
		hashKey = &envHashKey
	}
	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		hashKey = &envCryptoKey
	}
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = &envRunAddr
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		r, _ := strconv.Atoi(envReportInterval)
		flagReportInterval = &r
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		r, _ := strconv.Atoi(envPollInterval)
		flagPollInterval = &r
	}
	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		r, _ := strconv.Atoi(envRateLimit)
		rateLimit = &r
	}

	return &Config{
		HashKey:        *hashKey,
		CryptoKey:      *cryptoKey,
		ServerAddress:  "http://" + *flagRunAddr,
		PollInterval:   time.Duration(*flagPollInterval) * time.Second,
		ReportInterval: time.Duration(*flagReportInterval) * time.Second,
		RateLimit:      *rateLimit,
	}
}
