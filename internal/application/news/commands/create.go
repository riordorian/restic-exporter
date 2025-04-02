package commands

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"restic-exporter/internal/application/storage"
	"restic-exporter/internal/domain/news"
	"restic-exporter/internal/domain/repository"
	"restic-exporter/internal/shared"
	"restic-exporter/internal/shared/dto"
)

type CreateHandler struct {
	NewsRepo   repository.NewsRepositoryInterface
	Transactor storage.TransactorInterface
}

type CreateHandlerInterface interface {
	Handle(ctx context.Context, req dto.CreateRequest) error
}

func NewCreateHandler(newsRepo repository.NewsRepositoryInterface,
	t storage.TransactorInterface) CreateHandlerInterface {
	return CreateHandler{
		NewsRepo:   newsRepo,
		Transactor: t,
	}
}

func (c CreateHandler) Handle(ctx context.Context, req dto.CreateRequest) error {
	var media []shared.Media
	userId, _ := ctx.Value("userId").(uuid.UUID)

	err := c.Transactor.MakeTransaction(ctx, func(ctx context.Context) error {
		newItem := news.News{
			Title:     req.Title,
			Text:      req.Text,
			Status:    1,
			CreatedBy: userId,
			Media:     media,
		}
		createdNewsId, err := c.NewsRepo.Insert(ctx, newItem)
		fmt.Println(createdNewsId)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
