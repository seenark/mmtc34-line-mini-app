// Command server runs the app's HTTP server.
package main

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/db"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/line"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/server"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/spa"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

// runSeed needs only DATA_DIR, not the full server environment.
func runSeed() error {
	dir := cmp.Or(os.Getenv("DATA_DIR"), "data")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	sqlDB, err := db.Open(filepath.Join(dir, "app.db"))
	if err != nil {
		return err
	}
	defer sqlDB.Close()
	return db.Seed(context.Background(), sqlDB, dir)
}

func run() error {
	if len(os.Args) > 1 {
		if os.Args[1] != "seed" {
			return fmt.Errorf("unknown command %q (only: seed)", os.Args[1])
		}
		return runSeed()
	}
	cfg, err := server.ConfigFromEnv()
	if err != nil {
		return err
	}
	if err := server.CheckSPABase(spa.Dist(), cfg.BasePath); err != nil {
		return err
	}
	if err := os.MkdirAll(cfg.DataDir, 0o755); err != nil {
		return err
	}
	sqlDB, err := db.Open(filepath.Join(cfg.DataDir, "app.db"))
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	handler, err := server.New(server.Deps{
		Config: cfg,
		DB:     sqlDB,
		Verify: line.Client{ChannelID: cfg.LineChannelID},
		Now:    time.Now,
	})
	if err != nil {
		return err
	}
	srv := &http.Server{
		Addr: ":" + cfg.Port, Handler: handler,
		ReadHeaderTimeout: 10 * time.Second, ReadTimeout: 60 * time.Second,
		WriteTimeout: 60 * time.Second, IdleTimeout: 120 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.ListenAndServe() }()
	slog.Info("listening", "addr", srv.Addr, "basePath", cfg.BasePath)

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
