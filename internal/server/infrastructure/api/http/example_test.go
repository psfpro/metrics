package http

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Example() {
	router := chi.NewRouter()
	srv := &http.Server{Addr: ":8080", Handler: router}
	app := NewApp(srv)

	fmt.Printf("addr: %+v", app.httpServer.Addr)

	// Output:
	// addr: :8080
}
