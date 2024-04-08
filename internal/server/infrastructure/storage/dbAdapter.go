package storage

import (
	"context"
	"database/sql"
	"log"

	"github.com/gofrs/uuid"

	"github.com/psfpro/metrics/internal/server/domain"
)

type DBAdapter struct {
	db                      *sql.DB
	counterMetricRepository *CounterMetricRepository
	gaugeMetricRepository   *GaugeMetricRepository
}

func NewDBAdapter(db *sql.DB, counterMetricRepository *CounterMetricRepository, gaugeMetricRepository *GaugeMetricRepository) *DBAdapter {
	return &DBAdapter{db: db, counterMetricRepository: counterMetricRepository, gaugeMetricRepository: gaugeMetricRepository}
}

func (a *DBAdapter) Flush(ctx context.Context) error {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	saveGaugeMetricQuery := `
INSERT INTO gauge_metric (id, metric_name, metric_value)
VALUES ($1, $2, $3)
ON CONFLICT (id)
DO UPDATE SET
    metric_name = excluded.metric_name,
    metric_value = excluded.metric_value,
    recorded_at = CURRENT_TIMESTAMP
`
	saveGaugeMetricStmt, _ := tx.PrepareContext(ctx, saveGaugeMetricQuery)
	defer saveGaugeMetricStmt.Close()
	saveCounterMetricQuery := `
INSERT INTO counter_metric (id, metric_name, metric_value)
VALUES ($1, $2, $3)
ON CONFLICT (id)
DO UPDATE SET
    metric_name = excluded.metric_name,
    metric_value = excluded.metric_value,
    recorded_at = CURRENT_TIMESTAMP
`
	saveCounterMetricStmt, _ := tx.PrepareContext(ctx, saveCounterMetricQuery)
	defer saveCounterMetricStmt.Close()

	for _, v := range a.gaugeMetricRepository.data {
		_, err := saveGaugeMetricStmt.ExecContext(ctx, a.uuidByName(v.Name()).String(), v.Name(), v.Value())
		if err != nil {
			log.Printf("Save gauge metric error: %v", err)
		}
	}
	for _, v := range a.counterMetricRepository.data {
		_, err2 := saveCounterMetricStmt.ExecContext(ctx, a.uuidByName(v.Name()).String(), v.Name(), v.Value())
		if err2 != nil {
			log.Printf("Save counter metric error: %v", err2)
		}
	}
	return tx.Commit()
}

func (a *DBAdapter) Restore(ctx context.Context) error {
	createGaugeMetricTableQuery := `
CREATE TABLE IF NOT EXISTS gauge_metric (
    id UUID PRIMARY KEY,
    metric_name VARCHAR(255),
    metric_value DOUBLE PRECISION,
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`
	createCounterMetricTableQuery := `
CREATE TABLE IF NOT EXISTS counter_metric (
    id UUID PRIMARY KEY,
    metric_name VARCHAR(255),
    metric_value INT,
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`
	a.db.ExecContext(context.TODO(), createGaugeMetricTableQuery)
	a.db.ExecContext(context.TODO(), createCounterMetricTableQuery)
	a.gaugeMetricRepository.data = make(map[string]*domain.GaugeMetric)
	a.counterMetricRepository.data = make(map[string]*domain.CounterMetric)

	// Gauge data restore
	gaugeRows, err := a.db.QueryContext(ctx, "SELECT metric_name, metric_value from gauge_metric")
	if err != nil {
		return err
	}
	defer gaugeRows.Close()

	for gaugeRows.Next() {
		var name string
		var value float64
		err = gaugeRows.Scan(&name, &value)
		if err != nil {
			return err
		}
		a.gaugeMetricRepository.data[name] = domain.NewGaugeMetric(name, value)
	}

	err = gaugeRows.Err()
	if err != nil {
		return err
	}

	// Counter data restore
	counterRows, err2 := a.db.QueryContext(ctx, "SELECT metric_name, metric_value from counter_metric")
	if err2 != nil {
		return err2
	}
	defer counterRows.Close()

	for counterRows.Next() {
		var name string
		var value int64
		err = counterRows.Scan(&name, &value)
		if err != nil {
			return err
		}
		a.counterMetricRepository.data[name] = domain.NewCounterMetric(name, value)
	}

	err = counterRows.Err()
	if err != nil {
		return err
	}

	return nil
}

func (a *DBAdapter) uuidByName(name string) uuid.UUID {
	ns, _ := uuid.FromString("6ebd7718-6855-499c-886e-54a4aea9c5c5")
	return uuid.NewV5(ns, name)
}
