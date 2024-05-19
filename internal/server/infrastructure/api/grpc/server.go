package grpc

import (
	"context"
	"github.com/psfpro/metrics/internal/server/application"
	"github.com/psfpro/metrics/proto"
	"log"
)

type MetricsServer struct {
	proto.UnimplementedMetricsServer
	updateGaugeMetricHandler     *application.UpdateGaugeMetricHandler
	updateCounterMetricHandler   *application.UpdateCounterMetricHandler
	increaseCounterMetricHandler *application.IncreaseCounterMetricHandler
}

func NewMetricsServer(
	updateGaugeMetricHandler *application.UpdateGaugeMetricHandler,
	updateCounterMetricHandler *application.UpdateCounterMetricHandler,
	increaseCounterMetricHandler *application.IncreaseCounterMetricHandler,
) *MetricsServer {
	return &MetricsServer{
		updateGaugeMetricHandler:     updateGaugeMetricHandler,
		updateCounterMetricHandler:   updateCounterMetricHandler,
		increaseCounterMetricHandler: increaseCounterMetricHandler,
	}
}

func (s *MetricsServer) Update(_ context.Context, in *proto.UpdateRequest) (*proto.UpdateResponse, error) {
	for _, metrics := range in.Metrics {
		if metrics.Type == proto.MetricType_GAUGE {
			if metrics.Id == "" || metrics.Value == 0 {
				continue
			}
			log.Printf("Update gauge %v: %v\n", metrics.Id, metrics.Value)
			s.updateGaugeMetricHandler.Handle(metrics.Id, float64(metrics.Value))
		} else if metrics.Type == proto.MetricType_COUNTER {
			if metrics.Id == "" {
				continue
			}
			if metrics.Delta == 0 {
				log.Printf("Increase counter %v\n", metrics.Id)
				s.increaseCounterMetricHandler.Handle(metrics.Id)
			} else {
				log.Printf("Update counter %v: %v\n", metrics.Id, metrics.Delta)
				s.updateCounterMetricHandler.Handle(metrics.Id, metrics.Delta)
			}
		}
	}

	return &proto.UpdateResponse{}, nil
}
