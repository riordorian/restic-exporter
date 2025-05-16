package di

import (
	"github.com/gorilla/mux"
	"github.com/sarulabs/di"
	"github.com/spf13/viper"
	logger "restic-exporter/internal/application/log"
	"restic-exporter/internal/infrastructure/adapters/prometheus"
	"restic-exporter/internal/infrastructure/ports"
	"restic-exporter/internal/infrastructure/ports/http"
)

var PortsServices = []di.Def{
	{
		Name:  "Router",
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return mux.NewRouter(), nil
		},
	},
	{
		Name:  "HttpServer",
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			config := ctn.Get("ConfigProvider").(*viper.Viper)
			router := ctn.Get("Router").(*mux.Router)
			log := ctn.Get("LoggerService").(logger.LoggerInterface)
			collector := ctn.Get("ResticCollector").(prometheus.ResticCollectorInterface)

			return http.GetServer(config.GetInt("EXPOSE_PORT"), router, log, collector), nil
		},
	},
	{
		Name:  "PortsServices",
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return ports.Services{
				HttpServer: ctn.Get("HttpServer").(*http.Server),
			}, nil
		},
	},
}
