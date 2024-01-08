package handler

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

type encoderResponseWriter struct {
	http.ResponseWriter
	EncoderWriter io.Writer
}

func (obj *encoderResponseWriter) Write(b []byte) (int, error) {
	return obj.EncoderWriter.Write(b)
}

func Compressor(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		isContentEncoding := strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
		log.Printf("Is content encoding: %v", isContentEncoding)
		if isContentEncoding {
			body := r.Body
			bodyBytes, err := gzip.NewReader(body)
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}
			defer bodyBytes.Close()
			r.Body = io.NopCloser(bodyBytes)
		}

		isAcceptEncoding := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
		log.Printf("Is accept encoding: %v", isAcceptEncoding)

		if isAcceptEncoding {
			w.Header().Set("Content-Encoding", "gzip")
			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}
			defer gz.Close()

			encoderWriter := encoderResponseWriter{
				ResponseWriter: w,
				EncoderWriter:  gz,
			}
			next.ServeHTTP(&encoderWriter, r)
		} else {
			next.ServeHTTP(w, r)
		}
	}

	return http.HandlerFunc(fn)
}
