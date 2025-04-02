package di

import (
	"github.com/sarulabs/di"
	"restic-exporter/internal/application"
	"restic-exporter/internal/application/cqrs"
)

var ApplicationServices = []di.Def{
	{
		Name:  "ApplicationServices",
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			//fs := ctn.Get("FilesystemStorage").(storage.FilesystemInterface)
			dispatcher := ctn.Get("Dispatcher").(cqrs.DispatcherInterface)
			return application.Services{
				Dispatcher: dispatcher,
			}, nil
		},
	},
	{
		Name:  "Dispatcher",
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {

			return cqrs.NewDispatcher(), nil
		},
	},
}
