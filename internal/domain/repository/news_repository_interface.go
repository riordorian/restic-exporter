package repository

import (
	"context"
	"github.com/google/uuid"
	"grpc/internal/domain/news"
	"grpc/internal/shared/dto"
)

type NewsRepositoryInterface interface {
	BaseRepositoryInterface
	GetList(ctx context.Context, request dto.ListRequest) ([]news.News, error)
	GetById(id uuid.UUID) (news.News, error)
	Insert(ctx context.Context, fields news.News) (id uuid.UUID, err error)
	Update(id uuid.UUID, fields news.News) (bool, error)
}
