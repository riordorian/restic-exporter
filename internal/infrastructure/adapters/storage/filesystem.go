package storage

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"restic-exporter/internal/domain/restic"
	"strings"
)

type Filesystem struct {
	Repos restic.ReposMap
}

func (f *Filesystem) FindAllRepos(ctx context.Context, rootDir string) (restic.ReposMap, error) {

	repos := make(restic.ReposMap, 100) //TODO: need dynamic size
	root := "."

	if rootDir != "" {
		root = rootDir
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && isResticRepo(path) {
			repo := restic.Repo{Path: path}

			cmd := exec.Command("restic", "-r", repo.Path, "snapshots", "--json", "--no-lock")
			cmd.Env = append(os.Environ(), "RESTIC_PASSWORD=1")
			output, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("restic failed: %v\nOutput: %s", err, string(output))
			}

			count := strings.Count(string(output), `"time":`)

			fmt.Println(string(output), count)
			repos[path] = repo
			return filepath.SkipDir
		}

		return nil
	})

	return repos, err
}

func (f *Filesystem) GetSnapshots(repo restic.Repo) ([]restic.Snapshot, error) {
	return nil, nil
}

func NewFilesystem() *Filesystem {
	return &Filesystem{}
}

func isResticRepo(path string) bool {
	required := []string{"config", "data", "index", "keys", "snapshots"}
	for _, file := range required {
		if _, err := os.Stat(filepath.Join(path, file)); os.IsNotExist(err) {
			return false
		}
	}

	return true
}
