package queries

import (
	"context"
	"errors"
	"restic-exporter/internal/application/cqrs"
	"restic-exporter/internal/application/storage"
)

type GetRepoStatisticQueryHandler struct {
	FileStorage storage.FilesystemInterface
}

func (c GetRepoStatisticQueryHandler) Handle(ctx context.Context, query cqrs.QueryInterface) (interface{}, error) {
	if q, ok := query.(GetRepoStatisticQuery); ok {
		//RESTIC_PASSWORD=1 restic -r . stats --json --no-lock - размер последнего, кол-во файлов
		//RESTIC_PASSWORD=1 restic -r . snapshots latest --json --no-lock // дата последнего

		repoStat, err := c.FileStorage.GetRepoStatistic(q.Repo)
		if err != nil {
			return nil, err
		}

		return repoStat, nil
	}

	return nil, errors.New("invalid query type. Expected GetRepoStatisticQuery")
}
