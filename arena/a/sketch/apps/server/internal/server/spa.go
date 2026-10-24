package server

import (
	"io/fs"

	"github.com/gin-gonic/gin"
)

// Ticket .1.

// checkSPABase fails startup when index.html was built for another BASE_PATH
// (it must reference "<BasePath>/assets/"). An index without assets (web not built,
// e.g. in `go test`) passes and spaHandler serves 503 "web app not built".
func checkSPABase(dist fs.FS, basePath string) error { panic("not implemented") }

// spaHandler is the router's NoRoute:
//   - outside <BasePath>                  -> 404
//   - exactly <BasePath>                  -> 301 to <BasePath>/
//   - under <BasePath>/api/               -> 404 JSON {error}
//   - existing file in dist               -> the file; /assets/* immutable cache, others no-cache
//   - anything else (client route)        -> index.html, Cache-Control: no-cache
func spaHandler(dist fs.FS, basePath string) gin.HandlerFunc { panic("not implemented") }
