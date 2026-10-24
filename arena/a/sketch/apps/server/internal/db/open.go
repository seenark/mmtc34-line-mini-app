// Package db is sqlc output (db.go, models.go, <queryfile>.sql.go) plus Open and Seed.
// Hand-written files here: open.go, seed.go. Everything else is generated; never edit it.
package db

import (
	"context"
	"database/sql"
	"embed"
)

//go:embed migrations/*.sql
var migrations embed.FS

//go:embed seed/*.sql
var seeds embed.FS

// Open opens <dir>/app.db with modernc.org/sqlite and runs goose migrations from the embedded FS.
// DSN pragmas: foreign_keys(1), journal_mode(WAL), busy_timeout(5000), synchronous(NORMAL);
// _txlock=immediate so every transaction takes the write lock up front (no SQLITE_BUSY on upgrade).
func Open(ctx context.Context, dir string) (*sql.DB, error) {
	panic("not implemented")
}

// Seed runs every embedded seed/*.sql file in name order in one transaction. Each file is
// idempotent (INSERT ... WHERE NOT EXISTS), so `mise run seed` can run repeatedly.
// Files: seed/10_jersey_2569.sql (ticket .3), seed/20_jersey_2569_sponsors.sql (ticket .10).
func Seed(ctx context.Context, sqlDB *sql.DB) error {
	panic("not implemented")
}
