package database

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const createMigrationsTable = `
CREATE TABLE IF NOT EXISTS schema_migrations (
    version    TEXT PRIMARY KEY,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)`

// MigrateUp applies every *.up.sql file that has not been applied yet.
func MigrateUp(ctx context.Context, pool *pgxpool.Pool, files fs.FS) error {
	if _, err := pool.Exec(ctx, createMigrationsTable); err != nil {
		return err
	}

	names, err := fs.Glob(files, "*.up.sql")
	if err != nil {
		return err
	}
	sort.Strings(names)

	applied := 0
	for _, name := range names {
		version := versionOf(name)

		var exists bool
		err := pool.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)`, version,
		).Scan(&exists)
		if err != nil {
			return err
		}
		if exists {
			continue
		}

		sql, err := fs.ReadFile(files, name)
		if err != nil {
			return err
		}
		if err := inTx(ctx, pool, func(tx pgx.Tx) error {
			if _, err := tx.Exec(ctx, string(sql)); err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}
			_, err := tx.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, version)
			return err
		}); err != nil {
			return err
		}
		log.Printf("applied %s", name)
		applied++
	}

	if applied == 0 {
		log.Println("no new migrations")
	}
	return nil
}

// MigrateDown rolls back the most recently applied migration.
func MigrateDown(ctx context.Context, pool *pgxpool.Pool, files fs.FS) error {
	if _, err := pool.Exec(ctx, createMigrationsTable); err != nil {
		return err
	}

	var version string
	err := pool.QueryRow(ctx,
		`SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1`,
	).Scan(&version)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Println("nothing to roll back")
		return nil
	}
	if err != nil {
		return err
	}

	matches, err := fs.Glob(files, version+"_*.down.sql")
	if err != nil {
		return err
	}
	if len(matches) != 1 {
		return fmt.Errorf("no down migration found for version %s", version)
	}

	sql, err := fs.ReadFile(files, matches[0])
	if err != nil {
		return err
	}
	if err := inTx(ctx, pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, string(sql)); err != nil {
			return fmt.Errorf("%s: %w", matches[0], err)
		}
		_, err := tx.Exec(ctx, `DELETE FROM schema_migrations WHERE version = $1`, version)
		return err
	}); err != nil {
		return err
	}
	log.Printf("rolled back %s", matches[0])
	return nil
}

func versionOf(filename string) string {
	v, _, _ := strings.Cut(filename, "_")
	return v
}

func inTx(ctx context.Context, pool *pgxpool.Pool, fn func(pgx.Tx) error) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // no-op after a successful commit

	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
