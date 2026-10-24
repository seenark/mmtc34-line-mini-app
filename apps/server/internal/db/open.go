// Package db is sqlc output (db.go, models.go, <queryfile>.sql.go) plus Open.
// Hand-written files here: open.go, and seed.go once ticket .3 adds it. Everything else is generated; never edit it.
package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"net/url"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrations embed.FS

// Open opens the SQLite file at path with modernc.org/sqlite and runs the embedded goose migrations.
// _txlock=immediate makes every transaction take the write lock up front, so a read-then-write
// sequence never fails with SQLITE_BUSY on lock upgrade.
func Open(path string) (*sql.DB, error) {
	q := url.Values{}
	for _, p := range []string{"foreign_keys(1)", "journal_mode(WAL)", "busy_timeout(5000)", "synchronous(NORMAL)"} {
		q.Add("_pragma", p)
	}
	q.Set("_txlock", "immediate")

	sqlDB, err := sql.Open("sqlite", path+"?"+q.Encode())
	if err != nil {
		return nil, err
	}
	if err := migrate(sqlDB); err != nil {
		sqlDB.Close()
		return nil, err
	}
	return sqlDB, nil
}

func migrate(sqlDB *sql.DB) error {
	fsys, err := fs.Sub(migrations, "migrations")
	if err != nil {
		return err
	}
	p, err := goose.NewProvider(goose.DialectSQLite3, sqlDB, fsys)
	if err != nil {
		return err
	}
	if _, err := p.Up(context.Background()); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}
