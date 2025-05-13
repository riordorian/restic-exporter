package prometheus

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"os"
	"restic-exporter/internal/application/cqrs"
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
	mu sync.Mutex
}

func NewResticCollector() *ResticCollector {
	once.Do(func() {
		collectorInstance = &ResticCollector{
			metrics:                 make(map[string]*RepoMetrics, 100), //TODO: need dynamic size
			vecContainerFactory:     &MetricsContainerFactory[*prometheus.GaugeVec]{},
			counterContainerFactory: &MetricsContainerFactory[*prometheus.CounterVec]{},
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

func (c *ResticCollector) CollectMetrics(
	wg *sync.WaitGroup,
	ctx context.Context,
	rootDir string,
	dispatcher cqrs.DispatcherInterface,
) {
	defer wg.Done()
	ticker := time.NewTicker(5 * time.Second) //TODO: need dynamic interval
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.SetOutput(os.Stdout)
			log.Println("Metrics collection stopped")
			return
		case <-ticker.C:
			repos, err := dispatcher.DispatchQuery(ctx, queries.CollectReposQuery{
				RootDir: rootDir,
			})

			if err != nil {
				log.SetOutput(os.Stderr)
				log.Println(err.Error())
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
							log.SetOutput(os.Stderr)
							log.Println(err.Error())
						}
					}()
					fmt.Println(path, repo)
				}
			} else {
				log.SetOutput(os.Stderr)
				log.Fatal("Invalid type of repos struct. Expected restic.ReposMap")
				return
			}

			// TODO: 1. Collect repos snapshots information in another goroutines
			// TODO: 2. Initialize ResticCollector repos by InitRepos() method
			// TODO: 3. Each goroutine should run CollectSnapshots query handler
			// TODO: 4. After that, must update ResticCollector metrics data by  gauge.Set(float64(snapshotCount))...
		}
	}
}

func (c *ResticCollector) CollectRepoSnapshotsInfo(
	ctx context.Context,
	dispatcher cqrs.DispatcherInterface,
	repo restic.Repo,
) error {
	snapshot, err := dispatcher.DispatchQuery(ctx, queries.GetSnapshotQuery{Repo: repo})
	if err != nil {
		return err
	}

	if _, ok := snapshot.(restic.Snapshot); !ok {
		return fmt.Errorf("invalid type of snapshot struct. Expected restic.Snapshot struct")
	}

	snapshotInfo := snapshot.(restic.Snapshot)

	c.metrics[repo.Path].Vector.metrics["restic_repo_snapshot_avg_size_bytes"] = baseMetrics["restic_repo_snapshot_avg_size_bytes"]
	c.metrics[repo.Path].Vector.metrics["restic_repo_snapshot_avg_size_bytes"].
		WithLabelValues(repo.Path, "latest").
		Set(float64(snapshotInfo.Size))

	c.metrics[repo.Path].Vector.metrics["restic_latest_snapshot_files_count"] = baseMetrics["restic_latest_snapshot_files_count"]
	c.metrics[repo.Path].Vector.metrics["restic_latest_snapshot_files_count"].
		WithLabelValues(repo.Path, "latest").
		Set(float64(snapshotInfo.FilesCount))

	repoStatistic, err := dispatcher.DispatchQuery(ctx, queries.GetRepoStatisticQuery{Repo: repo})

	if err != nil {
		return err
	}

	if _, ok := repoStatistic.(restic.Repo); !ok {
		return fmt.Errorf("invalid type of repo statistic struct. Expected restic.Repo struct")
	}
	repoStat := repoStatistic.(restic.Repo)

	c.metrics[repo.Path].Vector.metrics["restic_snapshots_count"] = baseMetrics["restic_snapshots_count"]
	c.metrics[repo.Path].Vector.metrics["restic_snapshots_count"].
		WithLabelValues(repo.Path).
		Set(float64(repoStat.SnapshotsCount))

	c.metrics[repo.Path].Vector.metrics["restic_repo_total_size"] = baseMetrics["restic_repo_total_size"]
	c.metrics[repo.Path].Vector.metrics["restic_repo_total_size"].
		WithLabelValues(repo.Path).
		Set(float64(repoStat.TotalSize))

	return nil
}

func (c *ResticCollector) GetGauges(repoPath string) *MetricsContainer[*prometheus.GaugeVec] {
	return c.metrics[repoPath].Vector
}

func (c *ResticCollector) GetCounters(repoPath string) *MetricsContainer[*prometheus.CounterVec] {
	return c.metrics[repoPath].CounterVec
}
