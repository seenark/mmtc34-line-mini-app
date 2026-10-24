package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/db"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/server"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

// The only test seam: the real router over HTTPS (Secure cookies need TLS) with a temp SQLite
// file, real migrations and a temp DATA_DIR. Fakes: Verifier and clock.
// Fixtures live in fixtures_<feature>_test.go files.

const (
	basePath = "/mmtc34" // tests run under a prefix by default, so prefix bugs fail tests
	start    = "2026-11-01T10:00:00+07:00"
)

type H struct {
	t    *testing.T
	srv  *httptest.Server
	base string

	mu  sync.Mutex
	now time.Time
}

// Option customises the server under test.
type Option func(*server.Config)

// WithBasePath builds the server under another BASE_PATH, e.g. "/x".
func WithBasePath(p string) Option { return func(c *server.Config) { c.BasePath = p } }

func New(t *testing.T, opts ...Option) *H {
	t.Helper()
	dir := t.TempDir()
	sqlDB, err := db.Open(filepath.Join(dir, "app.db"))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { sqlDB.Close() })

	cfg := server.Config{
		BasePath:           basePath,
		DataDir:            dir,
		LiffID:             "test-liff-id",
		LineChannelID:      "test-channel",
		SessionSecret:      bytes.Repeat([]byte("k"), 32),
		SuperAdminUsername: "admin",
		SuperAdminPassword: superAdminPassword,
		DevLogin:           true,
		Port:               "0",
	}
	for _, o := range opts {
		o(&cfg)
	}

	startAt, err := time.Parse(time.RFC3339, start)
	if err != nil {
		t.Fatal(err)
	}
	h := &H{t: t, base: cfg.BasePath, now: startAt}
	handler, err := server.New(server.Deps{Config: cfg, DB: sqlDB, Verify: fakeVerifier{}, Now: h.clock})
	if err != nil {
		t.Fatalf("server.New: %v", err)
	}
	h.srv = httptest.NewTLSServer(handler)
	t.Cleanup(h.srv.Close)
	return h
}

func (h *H) clock() time.Time {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.now
}

// Advance moves the injected clock, e.g. across a Campaign's close time.
func (h *H) Advance(d time.Duration) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.now = h.now.Add(d)
}

// fakeVerifier: token "ok:<userId>" -> userId; anything else -> error.
type fakeVerifier struct{}

func (fakeVerifier) Verify(_ context.Context, token string) (string, error) {
	if id, ok := strings.CutPrefix(token, "ok:"); ok && id != "" {
		return id, nil
	}
	return "", errors.New("invalid id token")
}

// Anonymous returns a Client with its own cookie jar and no session.
func (h *H) Anonymous() *Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		h.t.Fatal(err)
	}
	return &Client{h: h, http: &http.Client{
		Transport:     h.srv.Client().Transport,
		Jar:           jar,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}}
}

// Client is one browser. Paths are written as "/api/..." and prefixed with the base path.
// Redirects are not followed, so tests can assert them.
type Client struct {
	h    *H
	http *http.Client
}

// Do sends body as JSON (struct), multipart (Upload) or nothing (nil).
func (c *Client) Do(method, path string, body any) *Res { return c.DoAt(method, c.h.base+path, body) }

// DoAt is Do without the base path prefix, for asserting what is served outside or at the prefix itself.
func (c *Client) DoAt(method, fullPath string, body any) *Res {
	t := c.h.t
	t.Helper()

	var payload io.Reader
	contentType := ""
	switch b := body.(type) {
	case nil:
	case Upload:
		var buf bytes.Buffer
		mw := multipart.NewWriter(&buf)
		part, err := mw.CreateFormFile("file", b.Filename)
		if err != nil {
			t.Fatal(err)
		}
		part.Write(b.Bytes)
		mw.Close()
		payload, contentType = &buf, mw.FormDataContentType()
	default:
		raw, err := json.Marshal(b)
		if err != nil {
			t.Fatal(err)
		}
		payload, contentType = bytes.NewReader(raw), "application/json"
	}

	req, err := http.NewRequest(method, c.h.srv.URL+fullPath, payload)
	if err != nil {
		t.Fatal(err)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, fullPath, err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return &Res{t: t, method: method, path: fullPath, Code: resp.StatusCode, Body: raw, Header: resp.Header}
}

func (c *Client) GET(path string) *Res            { return c.Do(http.MethodGet, path, nil) }
func (c *Client) POST(path string, body any) *Res { return c.Do(http.MethodPost, path, body) }
func (c *Client) PUT(path string, body any) *Res  { return c.Do(http.MethodPut, path, body) }
func (c *Client) DELETE(path string) *Res         { return c.Do(http.MethodDelete, path, nil) }
func (c *Client) GETAt(fullPath string) *Res      { return c.DoAt(http.MethodGet, fullPath, nil) }

// Upload is a multipart body with field "file".
type Upload struct {
	Filename string
	Bytes    []byte
}

// JPEG is a tiny valid JPEG for Slip/image uploads.
func JPEG() Upload {
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, image.NewRGBA(image.Rect(0, 0, 1, 1)), nil); err != nil {
		panic(err)
	}
	return Upload{Filename: "slip.jpg", Bytes: buf.Bytes()}
}

type Res struct {
	t      *testing.T
	method string
	path   string
	Code   int
	Body   []byte
	Header http.Header
}

// Status fails the test, printing method, path and body, unless Code == want.
func (r *Res) Status(want int) *Res {
	r.t.Helper()
	if r.Code != want {
		r.t.Fatalf("%s %s: status %d, want %d; body: %s", r.method, r.path, r.Code, want, r.Body)
	}
	return r
}

// Error returns the Thai {error} message.
func (r *Res) Error() string {
	r.t.Helper()
	return As[api.ErrorResponse](r).Error
}

// Text returns the body as a string (CSV assertions).
func (r *Res) Text() string { return string(r.Body) }

// As decodes the JSON body into T or fails the test.
func As[T any](r *Res) T {
	r.t.Helper()
	var v T
	if err := json.Unmarshal(r.Body, &v); err != nil {
		r.t.Fatalf("%s %s: body is not %T JSON (%v): %s", r.method, r.path, v, err, r.Body)
	}
	return v
}

// Equal fails the test when got != want.
func Equal[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Fatalf("got %v, want %v", fmt.Sprint(got), fmt.Sprint(want))
	}
}
