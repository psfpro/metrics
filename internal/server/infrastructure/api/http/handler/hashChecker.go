package handler

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
)

type HashChecker struct {
	hashKey string
}

func NewHashChecker(hashKey string) *HashChecker {
	return &HashChecker{hashKey: hashKey}
}

func (h *HashChecker) Check(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if h.hashKey != "" {
			hmac := hmac.New(sha256.New, []byte(h.hashKey))
			hmac.Write(bodyBytes)
			signature := hex.EncodeToString(hmac.Sum(nil))
			header := r.Header.Get("HashSHA256")
			isHashHeaderNotEmpty := header != ""
			isInvalidHash := signature != header
			log.Printf("Hash header not empty: %v", isHashHeaderNotEmpty)
			log.Printf("Is invalid hash: %v", isInvalidHash)
			if isHashHeaderNotEmpty && isInvalidHash {
				log.Printf("Signature: %v", signature)
				w.Header().Set("HashSHA256", signature)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
