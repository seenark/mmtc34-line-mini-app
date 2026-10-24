package server

import (
	"bytes"
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

// CheckSPABase fails startup when index.html was built for another BASE_PATH
// (it must reference "<basePath>/assets/"). An absent index, or one without assets
// (web not built, e.g. in `go test`), passes and spaHandler serves 503.
func CheckSPABase(dist fs.FS, basePath string) error {
	index, err := fs.ReadFile(dist, "index.html")
	if err != nil {
		return nil
	}
	if bytes.Contains(index, []byte("/assets/")) && !bytes.Contains(index, []byte(basePath+"/assets/")) {
		return fmt.Errorf("embedded web app was not built for BASE_PATH %q; rebuild it with the same BASE_PATH", basePath)
	}
	return nil
}

// spaHandler is the router's NoRoute:
//   - outside <basePath>                -> 404
//   - exactly <basePath>                -> 301 to <basePath>/
//   - under <basePath>/api/             -> 404 JSON {error}
//   - existing file in dist             -> the file; /assets/* immutable cache, others no-cache
//   - anything else (client route)      -> index.html, Cache-Control: no-cache
func spaHandler(dist fs.FS, basePath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := c.Request.URL.Path
		if p == basePath {
			c.Redirect(http.StatusMovedPermanently, basePath+"/")
			return
		}
		rel, found := strings.CutPrefix(p, basePath+"/")
		if !found || rel == "api" || strings.HasPrefix(rel, "api/") ||
			(c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead) {
			abort(c, http.StatusNotFound, "ไม่พบข้อมูล")
			return
		}

		name := path.Clean("/" + rel)[1:]
		if info, err := fs.Stat(dist, name); err == nil && !info.IsDir() && name != "index.html" {
			if strings.HasPrefix(name, "assets/") {
				c.Header("Cache-Control", "public, max-age=31536000, immutable")
			} else {
				c.Header("Cache-Control", "no-cache")
			}
			http.ServeFileFS(c.Writer, c.Request, dist, name)
			return
		}

		index, err := fs.ReadFile(dist, "index.html")
		if err != nil {
			c.String(http.StatusServiceUnavailable, "web app not built")
			return
		}
		c.Header("Cache-Control", "no-cache")
		c.Data(http.StatusOK, "text/html; charset=utf-8", index)
	}
}
