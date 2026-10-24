package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Ticket .1. The ONLY test seam: the real router over HTTPS (Secure cookies need TLS) with a
// temp SQLite file, real migrations and a temp DATA_DIR. Fakes: Verifier and clock.
// Feature tickets add fixtures as methods in their own fixtures_<feature>_test.go file and never
// edit this file.

const (
	basePath   = "/mmtc34" // tests always run under a prefix, so prefix bugs fail tests
	adminLine  = "U-admin" // listed in ADMIN_LINE_USER_IDS
	start      = "2026-11-01T10:00:00+07:00"
)

type H struct {
	t   *testing.T
	srv *httptest.Server // httptest.NewTLSServer(server.New(...))
	now time.Time
}

// New: t.TempDir() as DATA_DIR, db.Open (migrations), Config{BasePath, AdminLineUserIDs:{adminLine},
// SessionSecret: fixed 32 bytes}, Verify: fakeVerifier, Now: func() time.Time { return h.now }.
func New(t *testing.T) *H { panic("not implemented") }

// Advance moves the injected clock, e.g. across a Campaign's close time.
func (h *H) Advance(d time.Duration) { h.now = h.now.Add(d) }

// fakeVerifier: token "ok:<userId>" -> userId; anything else -> error.
type fakeVerifier struct{}

// Anonymous returns a Client with its own cookie jar and no session.
func (h *H) Anonymous() *Client { panic("not implemented") }

// Client is one browser. Paths are written as "/api/..." and prefixed with basePath.
type Client struct {
	h    *H
	http *http.Client
}

// Do sends body as JSON (struct), multipart (Upload) or nothing (nil).
func (c *Client) Do(method, path string, body any) *Res { panic("not implemented") }
func (c *Client) GET(path string) *Res                  { return c.Do(http.MethodGet, path, nil) }
func (c *Client) POST(path string, body any) *Res       { return c.Do(http.MethodPost, path, body) }
func (c *Client) PUT(path string, body any) *Res        { return c.Do(http.MethodPut, path, body) }
func (c *Client) DELETE(path string) *Res               { return c.Do(http.MethodDelete, path, nil) }

// Upload is a multipart body with field "file".
type Upload struct {
	Filename string
	Bytes    []byte
}

// JPEG is a tiny valid JPEG for Slip/image uploads.
func JPEG() Upload { panic("not implemented") }

type Res struct {
	t      *testing.T
	Code   int
	Body   []byte
	Header http.Header
}

// Status fails the test (printing method, path and body) unless Code == want.
func (r *Res) Status(want int) *Res { panic("not implemented") }

// Error returns the Thai {error} message.
func (r *Res) Error() string { panic("not implemented") }

// Text returns the body as a string (CSV assertions).
func (r *Res) Text() string { return string(r.Body) }

// As decodes the JSON body into T or fails the test.
func As[T any](r *Res) T { panic("not implemented") }

// Equal fails with a diff when got != want.
func Equal[T comparable](t *testing.T, got, want T) { panic("not implemented") }
