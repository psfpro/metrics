package storage

import (
	"log"
	"net/http"
)

type Middleware struct {
	adapter Adapter
}

func NewMiddleware(adapter Adapter) *Middleware {
	return &Middleware{adapter: adapter}
}

func (obj *Middleware) Handle(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		err := obj.adapter.Flush(r.Context())
		if err != nil {
			log.Printf("Storage flush error: %v", err)
		}
	}

	return http.HandlerFunc(fn)
}
