package di

import (
	"github.com/sarulabs/di"
	"restic-exporter/internal/application"
	"restic-exporter/internal/application/cqrs"
	"restic-exporter/internal/application/prometheus/queries"
	"restic-exporter/internal/application/storage"
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
			fs := ctn.Get("FilesystemStorage").(storage.FilesystemInterface)

			dispatcher := cqrs.NewDispatcher()

			collectReposHandler := queries.CollectReposQueryHandler{
				FileStorage: fs,
			}
			dispatcher.RegisterQuery("CollectRepos", collectReposHandler)

			getSnapshotsHandler := queries.GetSnapshotsQueryHandler{
				FileStorage: fs,
			}
			dispatcher.RegisterQuery("GetSnapshots", getSnapshotsHandler)

			getRepoStatisticHandler := queries.GetRepoStatisticQueryHandler{
				FileStorage: fs,
			}
			dispatcher.RegisterQuery("GetRepoStatistic", getRepoStatisticHandler)

			return dispatcher, nil
		},
	},
}
