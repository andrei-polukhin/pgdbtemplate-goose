package pgdbtemplategoose_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/andrei-polukhin/pgdbtemplate"
	pgdbtemplategoose "github.com/andrei-polukhin/pgdbtemplate-goose"
	pgdbtemplatepq "github.com/andrei-polukhin/pgdbtemplate-pq"
	qt "github.com/frankban/quicktest"
	"github.com/pressly/goose/v3"
)

// testConnectionStringFunc creates a connection string for tests.
func testConnectionStringFunc(dbName string) string {
	return pgdbtemplate.ReplaceDatabaseInConnectionString(testConnectionString, dbName)
}

func TestGooseMigrationRunner(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ctx := context.Background()

	c.Run("Basic migration execution", func(c *qt.C) {
		c.Parallel()

		// Create temporary migration directory.
		tempDir := c.TempDir()
		migrationsDir := filepath.Join(tempDir, "migrations")
		err := os.MkdirAll(migrationsDir, 0755)
		c.Assert(err, qt.IsNil)

		// Create a simple migration file.
		migrationSQL := `-- +goose Up
CREATE TABLE goose_test_users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

-- +goose Down
DROP TABLE goose_test_users;
`
		migrationFile := filepath.Join(migrationsDir, "00001_create_users.sql")
		err = os.WriteFile(migrationFile, []byte(migrationSQL), 0644)
		c.Assert(err, qt.IsNil)

		// Create connection provider.
		provider := pgdbtemplatepq.NewConnectionProvider(testConnectionStringFunc)

		// Create goose migration runner.
		runner := pgdbtemplategoose.NewMigrationRunner(migrationsDir)

		// Create template manager.
		config := pgdbtemplate.Config{
			ConnectionProvider: provider,
			MigrationRunner:    runner,
		}

		tm, err := pgdbtemplate.NewTemplateManager(config)
		c.Assert(err, qt.IsNil)

		// Initialize template database.
		err = tm.Initialize(ctx)
		c.Assert(err, qt.IsNil)
		defer tm.Cleanup(ctx)

		// Create test database.
		testDB, dbName, err := tm.CreateTestDatabase(ctx)
		c.Assert(err, qt.IsNil)
		defer testDB.Close()
		defer tm.DropTestDatabase(ctx, dbName)

		// Verify the migration ran by checking if table exists.
		pqConn, ok := testDB.(*pgdbtemplatepq.DatabaseConnection)
		c.Assert(ok, qt.IsTrue, qt.Commentf("expected *pgdbtemplatepq.DatabaseConnection"))

		var tableName string
		err = pqConn.DB.QueryRowContext(ctx, `
			SELECT table_name 
			FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'goose_test_users'
		`).Scan(&tableName)
		c.Assert(err, qt.IsNil)
		c.Assert(tableName, qt.Equals, "goose_test_users")
	})

	c.Run("Multiple migrations", func(c *qt.C) {
		c.Parallel()

		// Create temporary migration directory.
		tempDir := c.TempDir()
		migrationsDir := filepath.Join(tempDir, "migrations")
		err := os.MkdirAll(migrationsDir, 0755)
		c.Assert(err, qt.IsNil)

		// Create first migration.
		migration1 := `-- +goose Up
CREATE TABLE goose_test_posts (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL
);

-- +goose Down
DROP TABLE goose_test_posts;
`
		err = os.WriteFile(filepath.Join(migrationsDir, "00001_create_posts.sql"), []byte(migration1), 0644)
		c.Assert(err, qt.IsNil)

		// Create second migration.
		migration2 := `-- +goose Up
CREATE TABLE goose_test_comments (
    id SERIAL PRIMARY KEY,
    post_id INTEGER REFERENCES goose_test_posts(id),
    content TEXT NOT NULL
);

-- +goose Down
DROP TABLE goose_test_comments;
`
		err = os.WriteFile(filepath.Join(migrationsDir, "00002_create_comments.sql"), []byte(migration2), 0644)
		c.Assert(err, qt.IsNil)

		// Create connection provider.
		provider := pgdbtemplatepq.NewConnectionProvider(testConnectionStringFunc)

		// Create goose migration runner.
		runner := pgdbtemplategoose.NewMigrationRunner(migrationsDir)

		// Create template manager.
		config := pgdbtemplate.Config{
			ConnectionProvider: provider,
			MigrationRunner:    runner,
		}

		tm, err := pgdbtemplate.NewTemplateManager(config)
		c.Assert(err, qt.IsNil)

		// Initialize template database.
		err = tm.Initialize(ctx)
		c.Assert(err, qt.IsNil)
		defer tm.Cleanup(ctx)

		// Create test database.
		testDB, dbName, err := tm.CreateTestDatabase(ctx)
		c.Assert(err, qt.IsNil)
		defer testDB.Close()
		defer tm.DropTestDatabase(ctx, dbName)

		// Verify both migrations ran.
		pqConn := testDB.(*pgdbtemplatepq.DatabaseConnection)

		var count int
		err = pqConn.DB.QueryRowContext(ctx, `
			SELECT COUNT(*) 
			FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name IN ('goose_test_posts', 'goose_test_comments')
		`).Scan(&count)
		c.Assert(err, qt.IsNil)
		c.Assert(count, qt.Equals, 2)
	})

	c.Run("Custom dialect", func(c *qt.C) {
		c.Parallel()

		// Create temporary migration directory.
		tempDir := c.TempDir()
		migrationsDir := filepath.Join(tempDir, "migrations")
		err := os.MkdirAll(migrationsDir, 0755)
		c.Assert(err, qt.IsNil)

		// Create migration.
		migrationSQL := `-- +goose Up
CREATE TABLE goose_dialect_test (
    id SERIAL PRIMARY KEY
);

-- +goose Down
DROP TABLE goose_dialect_test;
`
		err = os.WriteFile(filepath.Join(migrationsDir, "00001_test.sql"), []byte(migrationSQL), 0644)
		c.Assert(err, qt.IsNil)

		// Create connection provider.
		provider := pgdbtemplatepq.NewConnectionProvider(testConnectionStringFunc)

		// Create goose migration runner with explicit dialect.
		runner := pgdbtemplategoose.NewMigrationRunner(
			migrationsDir,
			pgdbtemplategoose.WithDialect(goose.DialectPostgres),
		)

		// Create template manager.
		config := pgdbtemplate.Config{
			ConnectionProvider: provider,
			MigrationRunner:    runner,
		}

		tm, err := pgdbtemplate.NewTemplateManager(config)
		c.Assert(err, qt.IsNil)

		// Initialize template database.
		err = tm.Initialize(ctx)
		c.Assert(err, qt.IsNil)
		defer tm.Cleanup(ctx)

		// Create test database.
		testDB, dbName, err := tm.CreateTestDatabase(ctx)
		c.Assert(err, qt.IsNil)
		defer testDB.Close()
		defer tm.DropTestDatabase(ctx, dbName)

		// Verify migration ran.
		pqConn := testDB.(*pgdbtemplatepq.DatabaseConnection)

		var tableName string
		err = pqConn.DB.QueryRowContext(ctx, `
			SELECT table_name 
			FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'goose_dialect_test'
		`).Scan(&tableName)
		c.Assert(err, qt.IsNil)
		c.Assert(tableName, qt.Equals, "goose_dialect_test")
	})

	c.Run("Wrong connection type error", func(c *qt.C) {
		c.Parallel()

		// Create temporary migration directory.
		tempDir := c.TempDir()
		migrationsDir := filepath.Join(tempDir, "migrations")
		err := os.MkdirAll(migrationsDir, 0755)
		c.Assert(err, qt.IsNil)

		// Create migration.
		migrationSQL := `-- +goose Up
CREATE TABLE test (id SERIAL);

-- +goose Down
DROP TABLE test;
`
		err = os.WriteFile(filepath.Join(migrationsDir, "00001_test.sql"), []byte(migrationSQL), 0644)
		c.Assert(err, qt.IsNil)

		// Create goose migration runner.
		runner := pgdbtemplategoose.NewMigrationRunner(migrationsDir)

		// Create a mock connection that is not pgdbtemplate-pq.
		type mockConnection struct {
			pgdbtemplate.DatabaseConnection
		}

		ctx := context.Background()
		err = runner.RunMigrations(ctx, &mockConnection{})
		c.Assert(err, qt.ErrorMatches, "goose adapter requires database/sql connection.*")
	})
}
