package shared

import "path/filepath"

type MediaType int

const (
	IMAGE = iota
	VIDEO
	DOCUMENT
)

type Media struct {
	Name string
	Path string
	Type MediaType
}

func GetMediaInstanceByPath(path string) Media {
	filename := path
	base := filepath.Base(filename)
	extension := filepath.Ext(base)
	filenameWithoutExt := base[0 : len(base)-len(extension)]

	return Media{
		Name: filenameWithoutExt,
		Path: path,
		Type: 0,
	}
}
