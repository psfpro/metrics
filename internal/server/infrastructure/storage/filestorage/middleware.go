package filestorage

import (
	"net/http"
)

type Middleware struct {
	entityManager *EntityManager
}

func NewMiddleware(entityManager *EntityManager) *Middleware {
	return &Middleware{entityManager: entityManager}
}

func (obj *Middleware) Handle(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		obj.entityManager.Flush()
	}

	return http.HandlerFunc(fn)
}
