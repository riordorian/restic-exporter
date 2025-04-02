package di

import (
	"github.com/sarulabs/di"
	promapp "restic-exporter/internal/application/prometheus"
	storageapp "restic-exporter/internal/application/storage"
	"restic-exporter/internal/infrastructure/adapters"
	"restic-exporter/internal/infrastructure/adapters/prometheus"
	"restic-exporter/internal/infrastructure/adapters/storage"
)

var AdaptersServices = []di.Def{
	{
		Name:  "PrometheusRegistryFactory",
		Scope: di.App,
		Build: func(c di.Container) (interface{}, error) {
			return prometheus.NewRegistryFactory(), nil
		},
	},
	{
		Name:  "FilesystemStorage",
		Scope: di.App,
		Build: func(c di.Container) (interface{}, error) {
			fsStorage := new(storage.Filesystem)
			return fsStorage, nil
		},
	},
	{
		Name:  "AdaptersServices",
		Scope: di.App,
		Build: func(c di.Container) (interface{}, error) {
			prometheusRegistryFactory := c.Get("PrometheusRegistryFactory").(promapp.RegistryFactoryInterface)
			filesystemStorage := c.Get("FilesystemStorage").(storageapp.FilesystemInterface)
			return adapters.Services{
				PrometheusRegistryFactory: prometheusRegistryFactory,
				FilesystemStorage:         filesystemStorage,
				ResticCollector:           prometheus.NewResticCollector(),
			}, nil
		},
	},
}
