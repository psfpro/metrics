package grpc

import (
	"context"
	"github.com/psfpro/metrics/internal/agent/model"
	"github.com/psfpro/metrics/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"time"
)

type Client struct {
	reportInterval int64
	wg             sync.WaitGroup
	client         proto.MetricsClient
}

func NewClient(reportInterval int64) *Client {
	conn, err := grpc.NewClient(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	c := proto.NewMetricsClient(conn)
	return &Client{reportInterval: reportInterval, client: c}
}

func (c *Client) Run(collectResults chan []model.Metrics, sendResults chan error, closed chan struct{}) {
	for w := 1; w <= 3; w++ {
		c.wg.Add(1)
		go c.send(w, collectResults, sendResults)
	}
	c.wg.Wait()
	close(sendResults)
	closed <- struct{}{} // можно завершать приложение gracefully
}

func (c *Client) send(id int, jobs <-chan []model.Metrics, results chan<- error) {
	dataForSend := make(map[string]model.Metrics)
	ticker := time.NewTicker(time.Duration(c.reportInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case j, ok := <-jobs:
			if !ok {
				log.Printf("send %d stopping\n", id)
				c.wg.Done()
				return
			}
			log.Printf("send %d starting task\n", id)
			for _, value := range j {
				dataForSend[value.ID] = value
			}
		case <-ticker.C:
			if len(dataForSend) == 0 {
				continue
			}
			log.Printf("send %d performing action\n", id)
			var values []model.Metrics
			for _, value := range dataForSend {
				values = append(values, value)
			}
			err := c.sendBatchMetrics(values)
			if err == nil {
				dataForSend = make(map[string]model.Metrics)
			}

			results <- err
		}
	}
}

func (c *Client) sendBatchMetrics(batch []model.Metrics) error {
	metrics := make([]*proto.Metric, len(batch))
	for _, m := range batch {
		if m.MType == "gauge" {
			metrics = append(metrics, &proto.Metric{
				Id:    m.ID,
				Type:  proto.MetricType_GAUGE,
				Delta: 0,
				Value: float32(*m.Value),
			})
		} else {
			metrics = append(metrics, &proto.Metric{
				Id:    m.ID,
				Type:  proto.MetricType_COUNTER,
				Delta: *m.Delta,
				Value: 0,
			})
		}
	}
	ctx := context.Background()
	_, err := c.client.Update(ctx, &proto.UpdateRequest{Metrics: metrics})
	if err != nil {
		return err
	}
	return nil
}
