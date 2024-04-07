package storage

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
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

		var err error
		retryDelays := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

		for _, delay := range retryDelays {
			err = obj.adapter.Flush(r.Context())
			if err == nil {
				return
			}
			var pgErr *pgconn.ConnectError
			if !errors.As(err, &pgErr) {
				log.Printf("Неизвестная ошибка %v: %v", reflect.TypeOf(err), err)
				return
			}

			time.Sleep(delay)
		}

		log.Printf("Storage flush error: %v", fmt.Errorf("после нескольких попыток: %w", err))
	}

	return http.HandlerFunc(fn)
}
