package server

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	httpServer *http.Server
	grpcServer *grpc.Server
}

func NewApp(httpServer *http.Server, grpcServer *grpc.Server) *App {
	return &App{httpServer: httpServer, grpcServer: grpcServer}
}

func (a *App) Run() {
	a.runHTTPServer()
	a.runGrpcServer()
	a.waitSignal()
}

func (a *App) runHTTPServer() {
	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()
}

func (a *App) runGrpcServer() {
	go func() {
		listen, err := net.Listen("tcp", ":3200")
		if err != nil {
			log.Fatal(err)
		}
		if err := a.grpcServer.Serve(listen); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()
}

func (a *App) waitSignal() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	sig := <-signalChan
	signal.Stop(signalChan)
	log.Printf("received signal %s, shutting down", sig.String())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	a.shutdown(ctx)
}

func (a *App) shutdown(ctx context.Context) {
	if err := a.httpServer.Shutdown(ctx); err != nil {
		log.Printf("shutdown http server error %v", err)
	}
	a.grpcServer.GracefulStop()
}
