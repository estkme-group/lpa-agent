package frontend

import (
	"embed"
	"io/fs"
)

//go:embed dist/*
var distFS embed.FS

func DistFS() fs.FS {
	dist, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(err)
	}
	return dist
}
