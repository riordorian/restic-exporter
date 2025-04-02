package queries

import (
	"context"
	"errors"
	"restic-exporter/internal/application/cqrs"
	"restic-exporter/internal/application/storage"
)

type CollectReposQueryHandler struct {
	FileStorage storage.FilesystemInterface
}

func (c CollectReposQueryHandler) Handle(ctx context.Context, query cqrs.QueryInterface) (interface{}, error) {
	if q, ok := query.(CollectReposQuery); ok {
		repos, err := c.FileStorage.FindAllRepos(ctx, q.RootDir)
		if err != nil {
			return nil, err
		}

		return repos, nil
	}

	return nil, errors.New("invalid query type. Expected CollectReposQuery")
}
