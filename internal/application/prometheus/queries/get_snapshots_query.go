package queries

import "restic-exporter/internal/domain/restic"

type GetSnapshotQuery struct {
	Repo restic.Repo
}

func (c GetSnapshotQuery) QueryName() string {
	return "GetSnapshots"
}
