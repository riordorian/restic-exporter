package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

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
