package queries

type CollectReposQuery struct {
	RootDir string
}

func (c CollectReposQuery) QueryName() string {
	return "CollectRepos"
}
