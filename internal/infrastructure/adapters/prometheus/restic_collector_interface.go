package prometheus

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"restic-exporter/internal/application/cqrs"
	"restic-exporter/internal/domain/restic"
	"sync"
)

type ResticCollectorInterface interface {
	prometheus.Collector

	AddRepo(repoPath string) error
	RemoveRepo(repoPath string)
	InitRepos(reposMap restic.ReposMap) *ResticCollector

	CollectMetrics(wg *sync.WaitGroup, ctx context.Context, rootDir string, dispatcher cqrs.DispatcherInterface)

	CollectRepoSnapshotsInfo(ctx context.Context, dispatcher cqrs.DispatcherInterface, repo restic.Repo) error
}
