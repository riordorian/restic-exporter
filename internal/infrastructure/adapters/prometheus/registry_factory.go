package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	promapp "restic-exporter/internal/application/prometheus"
)

type RegistryFactory struct{}

func (r *RegistryFactory) NewRegistry() prometheus.Registerer {
	return prometheus.NewRegistry()
}

func NewRegistryFactory() promapp.RegistryFactoryInterface {
	return &RegistryFactory{}
}
