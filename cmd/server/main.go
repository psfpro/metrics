package main

import (
	"flag"
	"github.com/psfpro/metrics/internal/server/infrastructure/api/http"
	"log"
	"os"
)

func main() {
	log.Println("Start server")
	addr := &http.NetAddress{}
	addr.Set(":8080")
	_ = flag.Value(addr)
	flag.Var(addr, "a", "Net address host:port")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		addr.Set(envRunAddr)
	}
	log.Println(addr.Host)
	log.Println(addr.Port)
	app := http.NewApp(&http.Config{
		Address: addr,
	})
	app.Run()
}
