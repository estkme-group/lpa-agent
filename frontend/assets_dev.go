//go:build dev

package frontend

import "os"

func DistFS() fs.FS {
	return os.DirFS("./frontend/dist")
}
