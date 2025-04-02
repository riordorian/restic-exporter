package adapters

import (
	"restic-exporter/internal/application/prometheus"
	"restic-exporter/internal/application/storage"
	promadapt "restic-exporter/internal/infrastructure/adapters/prometheus"
)

type Services struct {
	PrometheusRegistryFactory prometheus.RegistryFactoryInterface
	FilesystemStorage         storage.FilesystemInterface
	ResticCollector           promadapt.ResticCollectorInterface
}
