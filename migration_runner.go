package pgdbtemplategoose

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/andrei-polukhin/pgdbtemplate"
	pgdbtemplatepq "github.com/andrei-polukhin/pgdbtemplate-pq"
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
// The connection must be compatible with database/sql (e.g., using pgdbtemplate-pq).
func (r *MigrationRunner) RunMigrations(ctx context.Context, conn pgdbtemplate.DatabaseConnection) error {
	// Extract *sql.DB from connection.
	// This assumes the connection is from pgdbtemplate-pq which embeds *sql.DB.
	db, err := r.extractSQLDB(conn)
	if err != nil {
		return fmt.Errorf("goose adapter requires database/sql connection: %w", err)
	}

	// Set goose configuration.
	if err := goose.SetDialect(string(r.dialect)); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	// Apply options if any.
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
func (r *MigrationRunner) extractSQLDB(conn pgdbtemplate.DatabaseConnection) (*sql.DB, error) {
	// goose requires database/sql, so we expect pgdbtemplate-pq connection.
	pqConn, ok := conn.(*pgdbtemplatepq.DatabaseConnection)
	if !ok {
		return nil, fmt.Errorf("goose adapter requires pgdbtemplate-pq connection, got %T", conn)
	}
	return pqConn.DB, nil
}
