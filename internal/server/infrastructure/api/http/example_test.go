package http

import (
	"fmt"
	"github.com/go-chi/chi/v5"
)

func Example() {
	router := chi.NewRouter()
	app := NewApp(":8080", router)

	fmt.Printf("addr: %+v", app.addr)

	// Output:
	// addr: :8080
}
