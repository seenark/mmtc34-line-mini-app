// Package server is the HTTP shell: Gin routes, session, uploads, and one handler file per
// feature. Handlers parse input, call domain rules, run sqlc queries, write api DTOs.
// There is no service/repository layer: a request is traceable in
// <feature>.go -> domain/*.go -> db/queries/<feature>.sql.
package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/db"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/spa"
)

// Config is every environment setting. Parsed once in main (ConfigFromEnv); tests build it directly.
type Config struct {
	BasePath         string          // BASE_PATH, normalised "/mmtc34" (leading slash, no trailing)
	DataDir          string          // DATA_DIR: app.db and files/{images,slips}
	LiffID           string          // LIFF_ID, served by GET /api/config
	LineChannelID    string          // LINE_CHANNEL_ID, used to verify ID tokens
	SessionSecret    []byte          // SESSION_SECRET, >= 32 bytes
	AdminLineUserIDs map[string]bool // ADMIN_LINE_USER_IDS, comma-separated
	DevLogin         bool            // DEV_LOGIN, default false
	Port             string          // PORT, default 8080
}

// ConfigFromEnv reads and validates the environment; any invalid value is a startup error.
func ConfigFromEnv() (Config, error) {
	cfg := Config{
		BasePath:         "/" + strings.Trim(envOr("BASE_PATH", "/mmtc34"), "/"),
		DataDir:          envOr("DATA_DIR", "data"),
		LiffID:           os.Getenv("LIFF_ID"),
		LineChannelID:    os.Getenv("LINE_CHANNEL_ID"),
		SessionSecret:    []byte(os.Getenv("SESSION_SECRET")),
		AdminLineUserIDs: map[string]bool{},
		Port:             envOr("PORT", "8080"),
	}
	if cfg.BasePath == "/" {
		return Config{}, errors.New("BASE_PATH must not be empty")
	}
	if len(cfg.SessionSecret) < 32 {
		return Config{}, errors.New("SESSION_SECRET must be at least 32 bytes")
	}
	for _, id := range strings.Split(os.Getenv("ADMIN_LINE_USER_IDS"), ",") {
		if id = strings.TrimSpace(id); id != "" {
			cfg.AdminLineUserIDs[id] = true
		}
	}
	if v := os.Getenv("DEV_LOGIN"); v != "" {
		dev, err := strconv.ParseBool(v)
		if err != nil {
			return Config{}, fmt.Errorf("DEV_LOGIN: %w", err)
		}
		cfg.DevLogin = dev
	}
	if !cfg.DevLogin && (cfg.LiffID == "" || cfg.LineChannelID == "") {
		return Config{}, errors.New("LIFF_ID and LINE_CHANNEL_ID are required unless DEV_LOGIN is on")
	}
	return cfg, nil
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
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
// <BasePath>.
func New(d Deps) (http.Handler, error) {
	dist := spa.Dist()
	s := &Server{cfg: d.Config, sqlDB: d.DB, q: db.New(d.DB), verify: d.Verify, now: d.Now}

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.RedirectTrailingSlash = false
	engine.RedirectFixedPath = false
	apiGroup := engine.Group(d.Config.BasePath).Group("/api")
	apiGroup.Use(rejectCrossSiteWrites)
	s.mount(apiGroup)
	engine.NoRoute(spaHandler(dist, d.Config.BasePath))
	return engine, nil
}

// rejectCrossSiteWrites 403s a state-changing request that a browser marks as not
// same-origin. A missing Sec-Fetch-Site (non-browser client, old webview) passes.
func rejectCrossSiteWrites(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return
	}
	if site := c.GetHeader("Sec-Fetch-Site"); site != "" && site != "same-origin" {
		abort(c, http.StatusForbidden, "ไม่อนุญาตคำขอข้ามเว็บไซต์")
	}
}

// tx runs fn in one IMMEDIATE transaction (DSN _txlock=immediate), so check-then-write
// sequences cannot interleave with another writer.
func (s *Server) tx(ctx context.Context, fn func(q *db.Queries) error) error {
	t, err := s.sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer t.Rollback()
	if err := fn(s.q.WithTx(t)); err != nil {
		return err
	}
	return t.Commit()
}

func (s *Server) mountConfig(r routes) {
	r.public.GET("/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, api.Config{LiffID: s.cfg.LiffID, DevLogin: s.cfg.DevLogin})
	})
}
