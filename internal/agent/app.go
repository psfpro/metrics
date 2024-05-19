package agent

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/psfpro/metrics/internal/agent/infrastructure/collector"
	"github.com/psfpro/metrics/internal/agent/infrastructure/metrics/grpc"
	"github.com/psfpro/metrics/internal/agent/infrastructure/metrics/http"
	"github.com/psfpro/metrics/internal/agent/model"
)

type App struct {
	config            *Config
	closer            context.CancelFunc
	closed            chan struct{}
	collectorWorker   *collector.Worker
	metricsHTTPClient *http.Client
	metricsGrpcClient *grpc.Client
}

func NewApp(config *Config) *App {
	collectorWorker := collector.NewWorker()
	metricsHTTPClient := http.NewClient(config.ServerAddress, config.ReportInterval, config.HashKey, config.CryptoKey)
	metricsGrpcClient := grpc.NewClient(config.ReportInterval)

	return &App{
		config:            config,
		collectorWorker:   collectorWorker,
		metricsHTTPClient: metricsHTTPClient,
		metricsGrpcClient: metricsGrpcClient,
	}
}

func (obj *App) Run() {
	obj.runWorkers()
	obj.waitSignal()
}

func (obj *App) runWorkers() {
	ctx, closer := context.WithCancel(context.Background())
	obj.closer = closer
	obj.closed = make(chan struct{})

	collectJobs := make(chan int, obj.config.RateLimit)
	collectResults := make(chan []model.Metrics, obj.config.RateLimit)
	sendResults := make(chan error, obj.config.RateLimit)
	// запускаем сбор метрик, ждем завершения воркеров и закрываем канал
	go obj.collectorWorker.Run(collectJobs, collectResults)
	// запускаем отправку метрик, ждем завершения воркеров и закрываем канал
	//go obj.metricsHTTPClient.Run(collectResults, sendResults, obj.closed)
	go obj.metricsGrpcClient.Run(collectResults, sendResults, obj.closed)

	go func() {
		defer close(collectJobs)
		ticker := time.NewTicker(time.Duration(obj.config.PollInterval) * time.Second)
		for {
			select {
			case <-ctx.Done(): // Проверка на сигнал об отмене
				return
			case <-ticker.C:
				collectJobs <- 1
			}
		}
	}()

	go func() {
		for {
			err := <-sendResults
			if err != nil {
				log.Printf("Ошибка отправки метрик: %v", err)
			}
		}
	}()
}

func (obj *App) waitSignal() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	sig := <-signalChan
	signal.Stop(signalChan)
	log.Printf("received signal %s, shutting down", sig.String())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	obj.shutdown(ctx)
}

func (obj *App) shutdown(ctx context.Context) {
	obj.closer()
	for {
		select {
		case <-obj.closed: // ждём завершения процедуры graceful shutdown
			log.Println("shutdown gracefully")
			return
		case <-ctx.Done():
			log.Println(ctx.Err())
			return
		}
	}
}
