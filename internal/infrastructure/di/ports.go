package di

import (
	"fmt"
	"github.com/sarulabs/di"
	"github.com/spf13/viper"
	"grpc/internal/infrastructure/ports"
)

var PortsServices = []di.Def{
	{
		Name:  "GrpcServer",
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			config := ctn.Get("ConfigProvider").(*viper.Viper)
			address := fmt.Sprintf("%s:%s", config.GetString("GRPC_SERVER_HOST"), config.GetString("GRPC_SERVER_PORT"))
			fmt.Println(address)

			return nil, nil
		},
	},
	{
		Name:  "PortsServices",
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return ports.Services{
				//GrpcServer: ctn.Get("GrpcServer").(*grpc.NewsServer),
				//HttpServer: http.GetServer(appServices.Handler),
			}, nil
		},
	},
}
