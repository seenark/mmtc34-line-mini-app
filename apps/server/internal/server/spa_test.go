package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/gin-gonic/gin"
)

func TestCheckSPABase(t *testing.T) {
	built := func(base string) fstest.MapFS {
		return fstest.MapFS{"index.html": {Data: []byte(`<script src="` + base + `/assets/a.js"></script>`)}}
	}
	cases := []struct {
		name    string
		dist    fstest.MapFS
		wantErr bool
	}{
		{"matching base", built("/x"), false},
		{"other base", built("/mmtc34"), true},
		{"web not built", fstest.MapFS{"index.html": {Data: []byte("<html></html>")}}, false},
		{"no index", fstest.MapFS{".keep": {}}, false},
	}
	for _, tc := range cases {
		if err := CheckSPABase(tc.dist, "/x"); (err != nil) != tc.wantErr {
			t.Errorf("%s: err = %v, wantErr %v", tc.name, err, tc.wantErr)
		}
	}
}

func TestSPAHandler(t *testing.T) {
	dist := fstest.MapFS{
		"index.html":  {Data: []byte("<html>app</html>")},
		"assets/a.js": {Data: []byte("js")},
	}
	engine := gin.New()
	engine.NoRoute(spaHandler(dist, "/x"))

	cases := []struct {
		method, path string
		code         int
		cache        string
	}{
		{"GET", "/x", 301, ""},
		{"GET", "/y/campaigns", 404, ""},
		{"GET", "/x/api/anything", 404, ""},
		{"POST", "/x/campaigns", 404, ""},
		{"GET", "/x/assets/a.js", 200, "public, max-age=31536000, immutable"},
		{"GET", "/x/campaigns/1", 200, "no-cache"},
		{"GET", "/x/", 200, "no-cache"},
	}
	for _, tc := range cases {
		rec := httptest.NewRecorder()
		engine.ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, nil))
		if rec.Code != tc.code || rec.Header().Get("Cache-Control") != tc.cache {
			t.Errorf("%s %s: got %d %q, want %d %q", tc.method, tc.path, rec.Code, rec.Header().Get("Cache-Control"), tc.code, tc.cache)
		}
	}

	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x/campaigns/1", nil))
	if rec.Body.String() != "<html>app</html>" {
		t.Errorf("client route body = %q, want index.html", rec.Body.String())
	}
}

func TestSPAResponsesCarryHardeningHeaders(t *testing.T) {
	dist := fstest.MapFS{
		"index.html":  {Data: []byte("<html>app</html>")},
		"assets/a.js": {Data: []byte("js")},
	}
	engine := gin.New()
	engine.NoRoute(spaHandler(dist, "/x"))

	for _, p := range []string{"/x/", "/x/campaigns/1", "/x/assets/a.js"} {
		rec := httptest.NewRecorder()
		engine.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, p, nil))
		if got := rec.Header().Get("Referrer-Policy"); got != "same-origin" {
			t.Errorf("%s: Referrer-Policy = %q, want same-origin", p, got)
		}
		if got := rec.Header().Get("X-Content-Type-Options"); got != "nosniff" {
			t.Errorf("%s: X-Content-Type-Options = %q, want nosniff", p, got)
		}
	}
}
