package queries

import (
	"context"
	"grpc/internal/application/storage"
	"grpc/internal/domain/news"
	"grpc/internal/domain/repository"
	"grpc/internal/shared/dto"
)

type ListHandler struct {
	Repo       repository.NewsRepositoryInterface
	Transactor storage.TransactorInterface
}

type GetListHandlerInterface interface {
	Handle(ctx context.Context, req dto.ListRequest) ([]news.News, error)
}

func NewGetListHandler(repo repository.NewsRepositoryInterface, t storage.TransactorInterface) GetListHandlerInterface {
	return ListHandler{
		Repo:       repo,
		Transactor: t,
	}
}

func (l ListHandler) Handle(ctx context.Context, req dto.ListRequest) ([]news.News, error) {
	/*err := l.Transactor.MakeTransaction(ctx, func(ctx context.Context) error {
		id, _ := uuid.Parse("44266dc6-18d0-46bd-a2b5-238de53db2cb")
		new := news.New{
			Title:     "New 5",
			Text:      "New 5 text",
			Status:    1,
			CreatedBy: id,
		}
		_, err := l.Repo.Insert(ctx, new)
		if err != nil {
			return err
		}

		new2 := news.New{
			Title:     "New 6",
			Text:      "New 6 text",
			Status:    1,
			CreatedBy: id,
		}
		_, err2 := l.Repo.Insert(ctx, new2)
		if err2 != nil {
			return err2
		}
		return nil
	})*/
	list, err := l.Repo.GetList(ctx, req)

	if err != nil {
		return nil, err
	}

	return list, nil
}
