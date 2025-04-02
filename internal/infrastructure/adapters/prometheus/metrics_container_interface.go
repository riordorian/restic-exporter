package prometheus

type MetricsContainerInterface[T MetricType | MetricVecType] interface {
	Set(name string, metric T)
	Get(name string) (T, bool)
}

type MetricContainerFactoryInterface[T MetricType | MetricVecType] interface {
	NewMetricContainer() *MetricsContainer[T]
}
