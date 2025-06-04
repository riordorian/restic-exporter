package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

var baseMetrics = map[string]*prometheus.GaugeVec{
	"restic_repo_snapshot_avg_size_bytes": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "restic_repo_snapshot_avg_size_bytes",
			Help: "Size of the restic repository latest snapshot in bytes",
		},
		[]string{"repo_path", "snapshot_id"},
	),
	"restic_latest_snapshot_files_count": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "restic_latest_snapshot_files_count",
			Help: "Number of files in the restic repository latest snapshot",
		},
		[]string{"repo_path", "snapshot_id"},
	),
	"restic_snapshots_count": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "restic_snapshots_count",
			Help: "Number of snapshots in the repository",
		},
		[]string{"repo_path"},
	),
	"restic_repo_total_size": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "restic_repo_total_size",
			Help: "Total repository size in bytes",
		},
		[]string{"repo_path"},
	),
	"restic_last_backup_timestamp": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "restic_last_backup_timestamp",
			Help: "Timestamp of the last backup in the repository",
		},
		[]string{"repo_path", "latest"},
	),
}

// MetricType объединяет поддерживаемые типы метрик Prometheus
type MetricType interface {
	*prometheus.Gauge | *prometheus.Counter | *prometheus.Histogram | *prometheus.Summary
}

// MetricVecType объединяет векторные типы метрик
type MetricVecType interface {
	*prometheus.GaugeVec | *prometheus.CounterVec | *prometheus.HistogramVec | *prometheus.SummaryVec
}

// MetricsContainer универсальный контейнер для метрик
type MetricsContainer[T MetricType | MetricVecType] struct {
	metrics map[string]T
	mu      sync.Mutex
}

type MetricsContainerFactory[T MetricType | MetricVecType] struct{}

/*
	func (mcf *MetricsContainerFactory) NewMetricsContainer() *MetricsContainer[T] {
		return &MetricsContainer[T]{
			metrics: make(map[string]T, len(BaseMetrics)),
		}
	}
*/
func (mcf *MetricsContainerFactory[T]) NewMetricsContainer() *MetricsContainer[T] {
	return &MetricsContainer[T]{
		metrics: make(map[string]T),
	}
}

// Set добавляет или обновляет метрику
func (mc *MetricsContainer[T]) Set(name string, metric T) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.metrics[name] = metric
}

// Get возвращает метрику по имени
func (mc *MetricsContainer[T]) Get(name string) (T, bool) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	metric, exists := mc.metrics[name]
	return metric, exists
}

func (mcf *MetricsContainerFactory[T]) NewMetric(name, help string, labels []string) *prometheus.GaugeVec {
	metric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: help,
		},
		labels,
	)
	prometheus.MustRegister(metric)
	return metric
}
