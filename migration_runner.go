package pgdbtemplategoose

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/andrei-polukhin/pgdbtemplate"
	pgdbtemplatepgx "github.com/andrei-polukhin/pgdbtemplate-pgx"
	pgdbtemplatepq "github.com/andrei-polukhin/pgdbtemplate-pq"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// MigrationRunner implements pgdbtemplate.MigrationRunner using goose.
type MigrationRunner struct {
	migrationsDir string
	dialect       goose.Dialect
	opts          []goose.ProviderOption
}

// NewMigrationRunner creates a new goose-based migration runner.
//
// The migrationsDir should point to a directory containing goose migration files.
// By default, uses "postgres" dialect.
//
// Example:
//
//	runner := pgdbtemplategoose.NewMigrationRunner(
//	    "./migrations",
//	    pgdbtemplategoose.WithDialect("postgres"),
//	)
func NewMigrationRunner(migrationsDir string, options ...Option) *MigrationRunner {
	runner := &MigrationRunner{
		migrationsDir: migrationsDir,
		dialect:       goose.DialectPostgres,
		opts:          []goose.ProviderOption{},
	}

	for _, opt := range options {
		opt(runner)
	}

	return runner
}

// RunMigrations implements pgdbtemplate.MigrationRunner.RunMigrations.
//
// It runs all pending goose migrations on the provided database connection.
// Supports both pgdbtemplate-pq (database/sql) and pgdbtemplate-pgx (pgx/v5).
func (r *MigrationRunner) RunMigrations(ctx context.Context, conn pgdbtemplate.DatabaseConnection) error {
	// Extract *sql.DB from connection.
	// This assumes the connection is from pgdbtemplate-pq which embeds *sql.DB.
	db, err := r.extractSQLDB(conn)
	if err != nil {
		return fmt.Errorf("goose adapter requires database/sql connection: %w", err)
	}

	// Create goose provider with dialect.
	// Using NewProvider directly is thread-safe and doesn't require global SetDialect().
	provider, err := goose.NewProvider(r.dialect, db, os.DirFS(r.migrationsDir), r.opts...)
	if err != nil {
		return fmt.Errorf("failed to create goose provider: %w", err)
	}

	// Run migrations up to the latest version.
	_, err = provider.Up(ctx)
	if err != nil {
		return fmt.Errorf("failed to run goose migrations: %w", err)
	}

	return nil
}

// extractSQLDB attempts to extract *sql.DB from the connection.
// Supports both pgdbtemplate-pq and pgdbtemplate-pgx.
func (r *MigrationRunner) extractSQLDB(conn pgdbtemplate.DatabaseConnection) (*sql.DB, error) {
	// Try pgdbtemplate-pq first (embeds *sql.DB).
	if pqConn, ok := conn.(*pgdbtemplatepq.DatabaseConnection); ok {
		return pqConn.DB, nil
	}

	// Try pgdbtemplate-pgx (has Pool field).
	if pgxConn, ok := conn.(*pgdbtemplatepgx.DatabaseConnection); ok {
		// Wrap the pool with stdlib to get *sql.DB.
		db := stdlib.OpenDBFromPool(pgxConn.Pool)
		return db, nil
	}

	return nil, fmt.Errorf("goose adapter requires pgdbtemplate-pq or pgdbtemplate-pgx connection, got %T", conn)
}
