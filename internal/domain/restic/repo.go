package restic

const (
	INACTIVE = iota
	ACTIVE
)

type Repo struct {
	Path           string
	Name           string
	Source         string
	Snapshots      []Snapshot
	SnapshotsCount int
	TotalSize      int
	FilesCount     int
}

type ReposMap map[string]Repo
