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
	storeInterval   int
	fileStoragePath string
	restore         bool
	databaseDsn     *DatabaseDsn
}

func NewConfig() *Config {
	address := &NetAddress{}
	address.Set(":8080")
	databaseDsn := &DatabaseDsn{value: "host=localhost user=app password=pass dbname=app sslmode=disable"}
	_ = flag.Value(address)
	_ = flag.Value(databaseDsn)
	flag.Var(address, "a", "Net serverAddress host:port")
	flag.Var(databaseDsn, "d", "Database DSN")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		address.Set(envRunAddr)
	}
	if envDatabaseDsn := os.Getenv("DATABASE_DSN"); envDatabaseDsn != "" {
		databaseDsn.Set(envDatabaseDsn)
	}

	return &Config{
		serverAddress:   address,
		storeInterval:   300,
		fileStoragePath: "/tmp/metrics-db.json",
		restore:         true,
		databaseDsn:     databaseDsn,
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
