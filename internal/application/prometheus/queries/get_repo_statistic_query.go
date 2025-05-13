package queries

import "restic-exporter/internal/domain/restic"

type GetRepoStatisticQuery struct {
	Repo restic.Repo
}

func (c GetRepoStatisticQuery) QueryName() string {
	return "GetRepoStatistic"
}
