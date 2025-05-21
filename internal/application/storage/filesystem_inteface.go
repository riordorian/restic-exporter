package storage

import (
	"context"
	"restic-exporter/internal/domain/restic"
)

type FilesystemInterface interface {
	FindAllRepos(ctx context.Context, rootDir string) (restic.ReposMap, error)
	GetSnapshots(repo restic.Repo) ([]restic.Snapshot, error)
	GetLatestSnapshotInfo(repo restic.Repo) (restic.Snapshot, error)
	GetRepoStatistic(repo restic.Repo) (restic.Repo, error)
	FindAccessFiles(ctx context.Context, rootDir string) ([]restic.RepoAccess, error)
}
