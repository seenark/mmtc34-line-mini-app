// Package spa embeds the built web app. Vite writes into ./dist (build.outDir in
// apps/web/vite.config.ts). dist/.keep is committed so `go build`/`go test` work before a web build.
package spa

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var files embed.FS

// Dist is the built SPA rooted at dist/.
func Dist() fs.FS {
	d, err := fs.Sub(files, "dist")
	if err != nil {
		panic(err) // only fails on an invalid path literal
	}
	return d
}
