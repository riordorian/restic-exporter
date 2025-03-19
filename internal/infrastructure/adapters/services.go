package adapters

import (
	"grpc/internal/domain/repository"
	"grpc/internal/domain/search"
	"grpc/internal/infrastructure/db"
)

type Services struct {
	Database            *db.Db
	NewsRepository      repository.NewsRepositoryInterface
}
