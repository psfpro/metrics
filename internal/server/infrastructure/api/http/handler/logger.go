package handler

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"go.uber.org/zap"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func Logger(next http.Handler) http.Handler {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer func() {
		err = logger.Sync()
		if err != nil {
			fmt.Println(err)
		}
	}()

	sugar := *logger.Sugar()

	name := runtime.FuncForPC(reflect.ValueOf(next).Pointer()).Name()
	log.Printf("Handler function called: %s", name)

	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		body := r.Body
		bodyBytes, _ := io.ReadAll(body)
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
			"body", string(bodyBytes),
		)
	}

	return http.HandlerFunc(fn)
}
