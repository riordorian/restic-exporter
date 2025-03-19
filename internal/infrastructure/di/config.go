package di

import (
	"github.com/sarulabs/di"
	"grpc/config"
)

var ConfigServices = []di.Def{
	{
		Name:  "ConfigProvider",
		Scope: di.App,
		Build: func(c di.Container) (interface{}, error) {
			return config.InitConfig(), nil
		},
	},
}
