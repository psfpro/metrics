package filestorage

import (
	"github.com/psfpro/metrics/internal/server/domain"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestEntityManager_Flush(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive test",
		},
	}
	for _, tt := range tests {
		tmpFile, _ := os.CreateTemp("/tmp", "testem")
		defer os.Remove(tmpFile.Name())
		t.Run(tt.name, func(t *testing.T) {
			gauge := domain.NewGaugeMetric("gauge", 1.1)
			counter := domain.NewCounterMetric("counter", 1)
			em := NewEntityManager(tmpFile.Name())
			em.persistGaugeMetric(gauge)
			em.persistCounterMetric(counter)
			em.Flush()
			em.Flush()

			em2 := NewEntityManager(tmpFile.Name())
			em2.Restore()

			gaugeResult, _ := em2.findGaugeMetric("gauge")
			assert.Equal(t, gaugeResult.Value(), gauge.Value())
			counterResult, _ := em2.findCounterMetric("counter")
			assert.Equal(t, counterResult.Value(), counter.Value())
		})
	}
}
