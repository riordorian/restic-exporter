package prometheus

import "github.com/prometheus/client_golang/prometheus"

type RegistryFactoryInterface interface {
	NewRegistry() prometheus.Registerer
}
