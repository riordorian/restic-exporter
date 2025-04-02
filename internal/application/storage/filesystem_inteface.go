package storage

import (
	"context"
	"restic-exporter/internal/domain/restic"
)

type FilesystemInterface interface {
	FindAllRepos(ctx context.Context, rootDir string) (restic.ReposMap, error)
	GetSnapshots(repo restic.Repo) ([]restic.Snapshot, error)
}
