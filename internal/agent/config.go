package agent

import "time"

type Config struct {
	HashKey        string
	ServerAddress  string
	PollInterval   time.Duration
	ReportInterval time.Duration
	RateLimit      int
}
