package application

import (
	appnewscommands "grpc/internal/application/news/commands"
	appnews "grpc/internal/application/news/queries"
)

type Services struct {
	Handlers Handlers
}

type Handlers struct {
	Queries  Queries
	Commands Commands
}

type Queries struct {
	GetList appnews.GetListHandlerInterface
}

type Commands struct {
	Create appnewscommands.CreateHandlerInterface
}
