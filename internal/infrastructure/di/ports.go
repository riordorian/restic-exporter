package di

import (
	"github.com/gorilla/mux"
	"github.com/sarulabs/di"
	"github.com/spf13/viper"
	"grpc/internal/infrastructure/ports"
	"grpc/internal/infrastructure/ports/http"
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

			return http.GetServer(config.GetInt("EXPOSE_PORT"), router), nil
		},
	},
	{
		Name:  "PortsServices",
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return ports.Services{
				HttpServer: ctn.Get("httpServer").(*http.Server),
			}, nil
		},
	},
}
