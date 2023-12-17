package main

import (
	"flag"
	"github.com/psfpro/metrics/internal/server/infrastructure/api/http"
	"log"
)

func main() {
	log.Println("Start server")
	addr := &http.NetAddress{}
	addr.Set(":8080")
	_ = flag.Value(addr)
	flag.Var(addr, "a", "Net address host:port")
	flag.Parse()
	log.Println(addr.Host)
	log.Println(addr.Port)
	app := http.NewApp(&http.Config{
		Address: addr,
	})
	app.Run()
}
