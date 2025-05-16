package di

import (
	"github.com/sarulabs/di"
	"github.com/spf13/viper"
	logger "restic-exporter/internal/application/log"
	promapp "restic-exporter/internal/application/prometheus"
	storageapp "restic-exporter/internal/application/storage"
	"restic-exporter/internal/infrastructure/adapters"
	"restic-exporter/internal/infrastructure/adapters/log"
	"restic-exporter/internal/infrastructure/adapters/prometheus"
	"restic-exporter/internal/infrastructure/adapters/storage"
	"time"
)

var AdaptersServices = []di.Def{
	{
		Name:  "LoggerService",
		Scope: di.App,
		Build: func(c di.Container) (interface{}, error) {
			logAdapter, err := log.NewZapAdapter()
			if err != nil {
				return nil, err
			}

			return logAdapter, nil
		},
	},
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
			logger := c.Get("LoggerService").(logger.LoggerInterface)
			fsStorage.SetLogger(logger)
			return fsStorage, nil
		},
	},
	{
		Name:  "ResticCollector",
		Scope: di.App,
		Build: func(c di.Container) (interface{}, error) {
			logger := c.Get("LoggerService").(logger.LoggerInterface)
			collectingInterval := c.Get("ConfigProvider").(*viper.Viper).GetInt("METRICS_COLLECTING_INTERVAL_SECONDS")
			collector := prometheus.NewResticCollector(logger, time.Duration(collectingInterval)*time.Second)

			return collector, nil
		},
	},
	{
		Name:  "AdaptersServices",
		Scope: di.App,
		Build: func(c di.Container) (interface{}, error) {
			prometheusRegistryFactory := c.Get("PrometheusRegistryFactory").(promapp.RegistryFactoryInterface)
			filesystemStorage := c.Get("FilesystemStorage").(storageapp.FilesystemInterface)
			collector := c.Get("ResticCollector").(prometheus.ResticCollectorInterface)
			return adapters.Services{
				PrometheusRegistryFactory: prometheusRegistryFactory,
				FilesystemStorage:         filesystemStorage,
				ResticCollector:           collector,
			}, nil
		},
	},
}
