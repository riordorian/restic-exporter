package prometheus

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"restic-exporter/internal/application/cqrs"
	logger "restic-exporter/internal/application/log"
	"restic-exporter/internal/application/prometheus/queries"
	"restic-exporter/internal/domain/restic"
	"sync"
	"time"
)

var (
	collectorInstance *ResticCollector
	once              sync.Once
)

type RepoMetrics struct {
	Vector     *MetricsContainer[*prometheus.GaugeVec]
	CounterVec *MetricsContainer[*prometheus.CounterVec]
}

type ResticCollector struct {
	repos                   restic.ReposMap
	metrics                 map[string]*RepoMetrics
	vecContainerFactory     *MetricsContainerFactory[*prometheus.GaugeVec]
	counterContainerFactory *MetricsContainerFactory[*prometheus.CounterVec]
	//metrics map[string]RepoMetrics
	mu                 sync.Mutex
	log                logger.LoggerInterface
	collectingInterval time.Duration
}

func NewResticCollector(log logger.LoggerInterface, collectingInterval time.Duration) *ResticCollector {
	once.Do(func() {
		collectorInstance = &ResticCollector{
			metrics:                 make(map[string]*RepoMetrics, 100), //TODO: need dynamic size
			vecContainerFactory:     &MetricsContainerFactory[*prometheus.GaugeVec]{},
			counterContainerFactory: &MetricsContainerFactory[*prometheus.CounterVec]{},
			log:                     log,
			collectingInterval:      collectingInterval,
		}
	})
	return collectorInstance
}

func (c *ResticCollector) AddRepo(repoPath string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.metrics[repoPath]; !exists {
		c.metrics[repoPath] = &RepoMetrics{
			Vector:     c.vecContainerFactory.NewMetricsContainer(),
			CounterVec: c.counterContainerFactory.NewMetricsContainer(),
		}

		return nil

	} else {
		return fmt.Errorf("repo %s already exists", repoPath)
	}
}

// RemoveRepo удаляет репозиторий из коллектора
func (c *ResticCollector) RemoveRepo(repoPath string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.metrics, repoPath)
}

func (c *ResticCollector) InitRepos(
	reposMap restic.ReposMap,
) *ResticCollector {
	c.mu.Lock()
	defer c.mu.Unlock()

	for repoPath := range reposMap {

		if _, exists := c.metrics[repoPath]; !exists {
			c.metrics[repoPath] = &RepoMetrics{
				Vector:     c.vecContainerFactory.NewMetricsContainer(),
				CounterVec: c.counterContainerFactory.NewMetricsContainer(),
			}

			for name, metricVec := range baseMetrics {
				c.metrics[repoPath].Vector.Set(name, metricVec)
			}
		}
	}

	return c
}

func (c *ResticCollector) Describe(ch chan<- *prometheus.Desc) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, metric := range baseMetrics {
		metric.Describe(ch)
	}

	/*for _, repoMetrics := range c.metrics {
		for _, metric := range repoMetrics.Vector.metrics {
			metric.Describe(ch)
		}
		for _, metric := range repoMetrics.CounterVec.metrics {
			metric.Describe(ch)
		}
	}*/
}

func (c *ResticCollector) Collect(ch chan<- prometheus.Metric) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, repoMetrics := range c.metrics {
		for _, metric := range repoMetrics.Vector.metrics {
			metric.Collect(ch)
		}
		for _, metric := range repoMetrics.CounterVec.metrics {
			metric.Collect(ch)
		}
	}
}

func (c *ResticCollector) CollectMetricsFlow(
	wg *sync.WaitGroup,
	ctx context.Context,
	rootDir string,
	dispatcher cqrs.DispatcherInterface,
) {
	defer wg.Done()
	ticker := time.NewTicker(c.collectingInterval)
	defer ticker.Stop()
	startTimer := time.After(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			c.log.Info("Restic collection stopped")
			return
		case <-startTimer:
			startTimer = nil
			c.collectMetrics(wg, ctx, rootDir, dispatcher)
		case <-ticker.C:
			c.collectMetrics(wg, ctx, rootDir, dispatcher)
		}
	}
}

func (c *ResticCollector) collectMetrics(
	wg *sync.WaitGroup,
	ctx context.Context,
	rootDir string,
	dispatcher cqrs.DispatcherInterface,
) {
	repos, err := dispatcher.DispatchQuery(ctx, queries.CollectReposQuery{
		RootDir: rootDir,
	})

	if err != nil {
		c.log.Error(err.Error())
		return
	}

	if repos, ok := repos.(restic.ReposMap); ok {
		c.InitRepos(repos)

		for path, repo := range repos {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := c.CollectRepoSnapshotsInfo(ctx, dispatcher, repo)
				if err != nil {
					c.log.Error(err.Error())
				}
			}()
			c.log.Info("Restic repo collect info: ", "path", path)
		}
	} else {
		c.log.Error("Invalid type of repos struct. Expected restic.ReposMap")
		return
	}
}

func (c *ResticCollector) CollectRepoSnapshotsInfo(
	ctx context.Context,
	dispatcher cqrs.DispatcherInterface,
	repo restic.Repo,
) error {
	c.log.Info("Collect repo info: ", "path", repo.Path)
	snapshot, err := dispatcher.DispatchQuery(ctx, queries.GetSnapshotQuery{Repo: repo})
	if err != nil {
		c.log.Info("Failed to collect repo info: ", "path", repo.Path, "snapshot", snapshot)
		c.log.Error(err.Error())
		return err

	}

	if _, ok := snapshot.(restic.Snapshot); !ok {
		return fmt.Errorf("invalid type of snapshot struct. Expected restic.Snapshot struct")
	}

	snapshotInfo := snapshot.(restic.Snapshot)

	c.metrics[repo.Path].Vector.metrics["restic_repo_snapshot_avg_size_bytes"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "restic_repo_snapshot_avg_size_bytes",
			Help: "Size of the restic repository latest snapshot in bytes",
		},
		[]string{"repo_path", "repo_name", "snapshot_id"},
	)
	c.metrics[repo.Path].Vector.metrics["restic_repo_snapshot_avg_size_bytes"].
		WithLabelValues(repo.Path, repo.Name, "latest").
		Set(float64(snapshotInfo.Size))

	c.metrics[repo.Path].Vector.metrics["restic_last_backup_timestamp"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "restic_last_backup_timestamp",
			Help: "Timestamp of the last backup in the repository",
		},
		[]string{"repo_path", "repo_name", "snapshot_id"},
	)
	c.metrics[repo.Path].Vector.metrics["restic_last_backup_timestamp"].
		WithLabelValues(repo.Path, repo.Name, "latest").
		Set(float64(snapshotInfo.Timestamp.Unix()))

	c.metrics[repo.Path].Vector.metrics["restic_latest_snapshot_files_count"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "restic_latest_snapshot_files_count",
			Help: "Number of files in the restic repository latest snapshot",
		},
		[]string{"repo_path", "repo_name", "snapshot_id"},
	)
	c.metrics[repo.Path].Vector.metrics["restic_latest_snapshot_files_count"].
		WithLabelValues(repo.Path, repo.Name, "latest").
		Set(float64(snapshotInfo.FilesCount))

	repoStatistic, err := dispatcher.DispatchQuery(ctx, queries.GetRepoStatisticQuery{Repo: repo})

	if err != nil {
		return err
	}

	if _, ok := repoStatistic.(restic.Repo); !ok {
		return fmt.Errorf("invalid type of repo statistic struct. Expected restic.Repo struct")
	}
	repoStat := repoStatistic.(restic.Repo)

	c.metrics[repo.Path].Vector.metrics["restic_snapshots_count"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "restic_snapshots_count",
			Help: "Number of snapshots in the repository",
		},
		[]string{"repo_path", "repo_name"},
	)
	c.metrics[repo.Path].Vector.metrics["restic_snapshots_count"].
		WithLabelValues(repo.Path, repo.Name).
		Set(float64(repoStat.SnapshotsCount))

	c.metrics[repo.Path].Vector.metrics["restic_repo_total_size"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "restic_repo_total_size",
			Help: "Total repository size in bytes",
		},
		[]string{"repo_path", "repo_name"},
	)
	c.metrics[repo.Path].Vector.metrics["restic_repo_total_size"].
		WithLabelValues(repo.Path, repo.Name).
		Set(float64(repoStat.TotalSize))

	return nil
}

func (c *ResticCollector) GetGauges(repoPath string) *MetricsContainer[*prometheus.GaugeVec] {
	return c.metrics[repoPath].Vector
}

func (c *ResticCollector) GetCounters(repoPath string) *MetricsContainer[*prometheus.CounterVec] {
	return c.metrics[repoPath].CounterVec
}
