package pgdbtemplategoose

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/andrei-polukhin/pgdbtemplate"
	pgdbtemplatepgx "github.com/andrei-polukhin/pgdbtemplate-pgx"
	pgdbtemplatepq "github.com/andrei-polukhin/pgdbtemplate-pq"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// MigrationRunner implements pgdbtemplate.MigrationRunner using goose.
type MigrationRunner struct {
	migrationsFs fs.FS
	dialect      goose.Dialect
	opts         []goose.ProviderOption
}

// NewMigrationRunner creates a new goose-based migration runner.
//
// The migrationsFs parameter accepts any fs.FS implementation containing goose migration files.
// This provides flexibility to use os.DirFS(), embed.FS, or any custom fs.FS implementation.
// By default, uses goose.DialectPostgres dialect.
//
// Example with filesystem directory:
//
//	migrationsFs := os.DirFS("./migrations")
//	runner := pgdbtemplategoose.NewMigrationRunner(
//	    migrationsFs,
//	    pgdbtemplategoose.WithDialect(goose.DialectPostgres),
//	)
//
// Example with embedded filesystem:
//
//	//go:embed migrations/*.sql
//	var migrationsFS embed.FS
//	runner := pgdbtemplategoose.NewMigrationRunner(migrationsFS)
func NewMigrationRunner(migrationsFs fs.FS, options ...Option) *MigrationRunner {
	runner := &MigrationRunner{
		migrationsFs: migrationsFs,
		dialect:      goose.DialectPostgres,
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
	provider, err := goose.NewProvider(r.dialect, db, r.migrationsFs, r.opts...)
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
