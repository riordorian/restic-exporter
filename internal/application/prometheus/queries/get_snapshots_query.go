package queries

import "restic-exporter/internal/domain/restic"

type GetSnapshotsQuery struct {
	Repo restic.Repo
}

func (c GetSnapshotsQuery) QueryName() string {
	return "GetSnapshots"
}
