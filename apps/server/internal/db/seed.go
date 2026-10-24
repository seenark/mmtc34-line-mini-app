package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
)

//go:embed seed/*.sql
var seeds embed.FS

// Seed runs every embedded seed/*.sql in name order inside one transaction. Each file must be
// idempotent (INSERT ... WHERE NOT EXISTS), so running the command twice changes nothing.
func Seed(ctx context.Context, sqlDB *sql.DB) error {
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
