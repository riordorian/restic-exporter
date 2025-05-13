package queries

import (
	"context"
	"errors"
	"restic-exporter/internal/application/cqrs"
	"restic-exporter/internal/application/storage"
)

type GetSnapshotsQueryHandler struct {
	FileStorage storage.FilesystemInterface
}

func (c GetSnapshotsQueryHandler) Handle(ctx context.Context, query cqrs.QueryInterface) (interface{}, error) {
	if q, ok := query.(GetSnapshotQuery); ok {
		//RESTIC_PASSWORD=1 restic -r . stats --json latest --no-lock - размер последнего, кол-во файлов
		//RESTIC_PASSWORD=1 restic -r . snapshots latest --json --no-lock // дата последнего

		snapshot, err := c.FileStorage.GetLatestSnapshotInfo(q.Repo)
		if err != nil {
			return nil, err
		}

		return snapshot, nil
	}

	return nil, errors.New("invalid query type. Expected GetSnapshotsQuery")
}
