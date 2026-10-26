package db

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

//go:embed seed/*.sql seed/images/*.jpg
var seeds embed.FS

// Seed copies the embedded images into dataDir/files/images, then runs every embedded seed/*.sql
// in name order inside one transaction. Each SQL file must be idempotent (INSERT ... WHERE NOT
// EXISTS) and each image is copied only when absent, so running the command twice changes nothing.
func Seed(ctx context.Context, sqlDB *sql.DB, dataDir string) error {
	if err := seedImages(filepath.Join(dataDir, "files", "images")); err != nil {
		return err
	}
	files, err := fs.Glob(seeds, "seed/*.sql")
	if err != nil {
		return err
	}
	tx, err := sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, f := range files {
		body, err := seeds.ReadFile(f)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, string(body)); err != nil {
			return fmt.Errorf("seed %s: %w", f, err)
		}
	}
	return tx.Commit()
}

func seedImages(dir string) error {
	entries, err := fs.ReadDir(seeds, "seed/images")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	for _, e := range entries {
		dst := filepath.Join(dir, e.Name())
		if _, err := os.Stat(dst); err == nil {
			continue
		} else if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		body, err := seeds.ReadFile(path.Join("seed/images", e.Name()))
		if err != nil {
			return err
		}
		if err := writeFileAtomic(dst, body); err != nil {
			return fmt.Errorf("seed image %s: %w", e.Name(), err)
		}
	}
	return nil
}

// writeFileAtomic keeps a crash from leaving a truncated image that later runs would skip.
func writeFileAtomic(dst string, body []byte) error {
	tmp, err := os.CreateTemp(filepath.Dir(dst), ".seed-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(body); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmp.Name(), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), dst)
}
