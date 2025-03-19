package di

import (
	"github.com/sarulabs/di"
	"grpc/internal/application"
	appnewscommands "grpc/internal/application/news/commands"
	appnews "grpc/internal/application/news/queries"
	"grpc/internal/domain/repository"
	"grpc/internal/infrastructure/db"
)

var ApplicationServices = []di.Def{
	{
		Name:  "ApplicationHandlers",
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			transactor := ctn.Get("Database").(*db.Db)
			newsRepository := ctn.Get("NewsRepository").(repository.NewsRepositoryInterface)

			return application.Handlers{
				Queries: application.Queries{
					GetList: appnews.NewGetListHandler(newsRepository, transactor),
				},
				Commands: application.Commands{
					Create: appnewscommands.NewCreateHandler(newsRepository, transactor),
				},
			}, nil
		},
	},
}
