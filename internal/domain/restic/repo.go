package restic

const (
	INACTIVE = iota
	ACTIVE
)

type Repo struct {
	Path      string
	Name      string
	Source    string
	Snapshots []Snapshot
}

type ReposMap map[string]Repo
