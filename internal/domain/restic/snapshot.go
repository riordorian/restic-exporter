package restic

import "time"

type Snapshot struct {
	Id         string
	Timestamp  time.Time
	Host       string
	Paths      []string
	Size       int
	FilesCount int
	Duration   time.Duration
}
