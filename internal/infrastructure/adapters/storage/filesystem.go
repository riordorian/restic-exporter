package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	logger "restic-exporter/internal/application/log"
	"restic-exporter/internal/domain/restic"
	"strings"
	"time"
)

type Filesystem struct {
	Repos restic.ReposMap
	log   logger.LoggerInterface
}

func (f *Filesystem) SetLogger(log logger.LoggerInterface) {
	f.log = log
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
		if info.IsDir() && f.isResticRepo(path) {
			repoPath := strings.Split(filepath.Clean(path), string(filepath.Separator))
			n := len(repoPath)
			repoName := path
			if 2 == n {
				repoName = filepath.Join(repoPath[0], repoPath[1])
			}
			if n > 2 {
				repoName = filepath.Join(repoPath[n-1], repoPath[n-2])
			}

			repo := restic.Repo{Path: path, Name: repoName}
			repos[path] = repo
			return filepath.SkipDir
		}

		return nil
	})

	return repos, err
}

func (f *Filesystem) FindAccessFiles(ctx context.Context, rootDir string) ([]restic.RepoAccess, error) {

	accessFiles := make([]restic.RepoAccess, 0)
	root := "."

	if rootDir != "" {
		root = rootDir
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		repoPath, repoPassword := "", ""
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fileName := filepath.Base(path)
		if !strings.Contains(fileName, "access") {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			if strings.HasPrefix(line, "export RESTIC_REPOSITORY=") {
				repoPath = strings.TrimPrefix(line, "export RESTIC_REPOSITORY=")
			}

			if strings.HasPrefix(line, "export RESTIC_PASSWORD=") {
				repoPassword = strings.TrimPrefix(line, "export RESTIC_PASSWORD=")
				repoPassword = strings.Trim(repoPassword, `"'`)
			}
		}

		if repoPath != "" && repoPassword != "" {
			accessFile := restic.RepoAccess{Path: repoPath, Password: repoPassword}
			accessFiles = append(accessFiles, accessFile)
		}

		return nil
	})

	return accessFiles, err
}

func (f *Filesystem) GetSnapshots(repo restic.Repo) ([]restic.Snapshot, error) {
	cmd := exec.Command("restic", "-r", repo.Path, "snapshots", "--json", "--no-lock")
	cmd.Env = append(os.Environ(), "RESTIC_PASSWORD=1")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("restic failed: %v\nOutput: %s", err, string(output))
	}

	count := strings.Count(string(output), `"time":`)

	fmt.Println(string(output), count)

	// Todo: parse output
	return nil, nil
}

func NewFilesystem() *Filesystem {
	return &Filesystem{}
}

func (f *Filesystem) isResticRepo(path string) bool {
	f.log.Info(path)
	required := []string{"config", "data", "index", "keys", "snapshots"}
	for _, file := range required {
		if _, err := os.Stat(filepath.Join(path, file)); os.IsNotExist(err) {
			f.log.Info(err.Error())

			return false
		}
	}

	return true

	/*cmd := exec.Command("restic", "-r", path, "stats", "--password-file", "./.restic-password", "--json", "--no-lock")
	output, err := cmd.CombinedOutput()

	if err != nil {
		f.log.Info(cmd.String())
		f.log.Error(err.Error())
	}

	var snapshot map[string]interface{}

	err = json.Unmarshal([]byte(output), &snapshot)
	if err != nil {
		f.log.Info(cmd.String())
		f.log.Error(string(output))
		f.log.Error("JSON decode error:", err)
		return false
	}

	// Проверяем наличие ключа "snapshots_count"
	if _, exists := snapshot["snapshots_count"]; exists {
		return true
	} else {
		return false
	}*/
}

// TODO: Need Refactor. Get snapshot info by snapshot id

func (f *Filesystem) GetLatestSnapshotInfo(repo restic.Repo) (restic.Snapshot, error) {
	snapshot := restic.Snapshot{}
	cmd := exec.Command("restic", "-r", repo.Path, "--password-file", "./.restic-password", "snapshots", "latest", "--json", "--no-lock")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return restic.Snapshot{}, fmt.Errorf("restic failed: %v\nOutput: %s", err, string(output))
	}

	var snapshots []map[string]interface{}
	err = json.Unmarshal(output, &snapshots)
	if err != nil {
		return snapshot, err
	}

	if len(snapshots) == 0 {
		return snapshot, fmt.Errorf("no snapshots found")
	}

	snapshotInfo := snapshots[0]
	if _, exists := snapshotInfo["time"]; exists {
		snapshotTimestamp, err := time.Parse(time.RFC3339Nano, snapshotInfo["time"].(string))
		if err != nil {
			return snapshot, err
		}

		snapshot.Timestamp = snapshotTimestamp
		snapshot.Id = snapshotInfo["short_id"].(string)
	} else {
		return snapshot, err
	}

	cmd = exec.Command("restic", "-r", repo.Path, "--password-file", "./.restic-password", "stats", "--json", "--no-lock", "latest")

	output, err = cmd.CombinedOutput()
	if err != nil {
		return restic.Snapshot{}, fmt.Errorf("restic failed: %v\nOutput: %s", err, string(output))
	}
	latestSnapshot, err := f.jsonUnmarshal(string(output))
	if err != nil {
		return snapshot, err
	}

	keys := []string{"total_file_count", "total_size"}
	exists := true
	for _, key := range keys {
		if _, found := latestSnapshot[key]; !found {
			exists = false
			break
		}
	}

	if exists {
		snapshot.FilesCount = int(latestSnapshot["total_file_count"].(float64))
		snapshot.Size = int(latestSnapshot["total_size"].(float64))
	}

	return snapshot, nil
}

func (f *Filesystem) GetRepoStatistic(repo restic.Repo) (restic.Repo, error) {
	cmd := exec.Command("restic", "-r", repo.Path, "--password-file", "./.restic-password", "stats", "--json", "--no-lock")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return restic.Repo{}, fmt.Errorf("restic failed: %v\nOutput: %s", err, string(output))
	}

	stat, err := f.jsonUnmarshal(string(output))
	if err != nil {
		return repo, err
	}

	keys := []string{"total_size", "snapshots_count"}
	exists := true
	for _, key := range keys {
		if _, found := stat[key]; !found {
			exists = false
			break
		}
	}

	if exists {
		repo.SnapshotsCount = int(stat["snapshots_count"].(float64))
		repo.TotalSize = int(stat["total_size"].(float64))
	}

	return repo, nil
}

func (f *Filesystem) jsonUnmarshal(data string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, err
	}
	return result, nil

}
