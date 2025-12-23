package database

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
)

// migrateSetup helper initializes the migration setup by setting the base file system for migrations.
func migrateSetup(conn *Conn, fs fs.FS) (*sql.DB, error) {
	handle := conn.Handle()
	if handle == nil {
		return nil, errors.New("database handle was nil")
	}

	goose.SetBaseFS(fs)
	if err := goose.SetDialect("postgres"); err != nil {
		return nil, fmt.Errorf("failed to set goose dialect: %w", err)
	}

	return handle, nil
}

// MigrateUpAll applies all pending migrations to the database using the embedded migration files.
func MigrateUpAll(conn *Conn, fs fs.FS) error {
	handle, err := migrateSetup(conn, fs)
	if err != nil {
		return err
	}

	if err := goose.Up(handle, "migration"); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

// MigrateDownAll rolls back all applied migrations in the database using the embedded migration files.
func MigrateDownAll(conn *Conn, fs fs.FS) error {
	handle, err := migrateSetup(conn, fs)
	if err != nil {
		return err
	}

	if err := goose.DownTo(handle, "migration", 0); err != nil {
		return fmt.Errorf("failed to down migrations: %w", err)
	}

	return nil
}
