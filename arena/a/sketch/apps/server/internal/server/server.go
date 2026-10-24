// Package server is the HTTP shell: Gin routes, session, uploads, and one handler file per
// feature. Handlers parse input, call domain rules, run sqlc queries, write api DTOs.
// There is no service/repository layer: a request is traceable in
// <feature>.go -> domain/*.go -> db/queries/<feature>.sql.
package server

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/db"
)

// Config is every environment setting. Parsed once in main (ConfigFromEnv); tests build it directly.
type Config struct {
	BasePath         string          // BASE_PATH, normalised "/mmtc34" (leading slash, no trailing)
	DataDir          string          // DATA_DIR: app.db and files/{images,slips}
	LiffID           string          // LIFF_ID, served by GET /api/config
	LineChannelID    string          // LINE_CHANNEL_ID, used by line.Client
	SessionSecret    []byte          // SESSION_SECRET, >= 32 bytes or startup fails
	AdminLineUserIDs map[string]bool // ADMIN_LINE_USER_IDS, comma-separated
	DevLogin         bool            // DEV_LOGIN, default false
	Port             string          // PORT, default 8080
}

// ConfigFromEnv reads and validates the environment; any invalid value is a startup error.
func ConfigFromEnv() (Config, error) {
	panic("not implemented")
}

// Verifier turns a LIFF ID token into a LINE userId. line.Client in production; tests fake it.
type Verifier interface {
	Verify(ctx context.Context, idToken string) (lineUserID string, err error)
}

// Deps are the only seams. Tests substitute Verify and Now; everything else is real.
type Deps struct {
	Config Config
	DB     *sql.DB // opened and migrated by db.Open
	Verify Verifier
	Now    func() time.Time
}

type Server struct {
	cfg    Config
	sqlDB  *sql.DB
	q      *db.Queries
	verify Verifier
	now    func() time.Time
}

// New builds the full router: API under <BasePath>/api, SPA for every other path under
// <BasePath>. It fails if the embedded SPA was built for a different BASE_PATH.
func New(d Deps) (http.Handler, error) {
	panic("not implemented")
}

// tx runs fn in one IMMEDIATE transaction (DSN _txlock=immediate), so check-then-write
// sequences cannot interleave with another writer.
func (s *Server) tx(ctx context.Context, fn func(q *db.Queries) error) error {
	panic("not implemented")
}

func (s *Server) mountConfig(r routes) {
	r.public.GET("/config", func(c *gin.Context) { panic("not implemented") })
}
