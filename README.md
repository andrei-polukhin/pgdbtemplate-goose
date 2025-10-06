# pgdbtemplate-goose

[![Go Reference](https://pkg.go.dev/badge/github.com/andrei-polukhin/pgdbtemplate-goose.svg)](https://pkg.go.dev/github.com/andrei-polukhin/pgdbtemplate-goose)
[![CI](https://github.com/andrei-polukhin/pgdbtemplate-goose/actions/workflows/test.yml/badge.svg)](https://github.com/andrei-polukhin/pgdbtemplate-goose/actions/workflows/test.yml)
[![Coverage](https://codecov.io/gh/andrei-polukhin/pgdbtemplate-goose/branch/main/graph/badge.svg)](https://codecov.io/gh/andrei-polukhin/pgdbtemplate-goose)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/andrei-polukhin/pgdbtemplate-goose/blob/main/LICENSE)

A [goose](https://github.com/pressly/goose) migration adapter for
[pgdbtemplate](https://github.com/andrei-polukhin/pgdbtemplate).

## Features

- 🦆 **Goose Integration** - Seamless integration with `pressly/goose`
- 🚀 **Fast Test Databases** - Leverage pgdbtemplate's template database speed
- 🔒 **Thread-safe** - Safe for concurrent test execution
- 📁 **Flexible Migration Sources** - Supports `fs.FS` interface (filesystem, embed.FS, etc.)
- ⚙️ **Configurable** - Support for custom goose options and dialects
- 🎯 **Flexible Driver Support** - Works with both `pq` and `pgx/v5`
- 🧪 **Well-tested** - 95%+ test coverage with real PostgreSQL integration

## Installation

```bash
go get github.com/andrei-polukhin/pgdbtemplate-goose
```

## Usage

### Using filesystem directory

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/andrei-polukhin/pgdbtemplate"
	pgdbtemplatepq "github.com/andrei-polukhin/pgdbtemplate-pq"
	pgdbtemplategoose "github.com/andrei-polukhin/pgdbtemplate-goose"
)

func main() {
	ctx := context.Background()

	// Create connection provider (works with both pq and pgx).
	connStringFunc := func(dbName string) string {
		return fmt.Sprintf("postgres://user:pass@localhost/%s?sslmode=disable", dbName)
	}
	provider := pgdbtemplatepq.NewConnectionProvider(connStringFunc)
	// Or use: provider := pgdbtemplatepgx.NewConnectionProvider(connStringFunc)

	// Create goose migration runner from filesystem directory.
	migrationsFs := os.DirFS("./migrations")
	migrationRunner := pgdbtemplategoose.NewMigrationRunner(migrationsFs)

	// Create template manager.
	config := pgdbtemplate.Config{
		ConnectionProvider: provider,
		MigrationRunner:    migrationRunner,
	}

	tm, err := pgdbtemplate.NewTemplateManager(config)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize template with goose migrations.
	if err := tm.Initialize(ctx); err != nil {
		log.Fatal(err)
	}
	defer tm.Cleanup(ctx)

	// Create test databases (fast!).
	testDB, dbName, err := tm.CreateTestDatabase(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer testDB.Close()
	defer tm.DropTestDatabase(ctx, dbName)

	log.Printf("Test database %s ready with goose migrations!", dbName)
}
```

### Using embed.FS (recommended for production)

```go
package main

import (
	"context"
	"embed"
	"fmt"
	"log"

	"github.com/andrei-polukhin/pgdbtemplate"
	pgdbtemplatepq "github.com/andrei-polukhin/pgdbtemplate-pq"
	pgdbtemplategoose "github.com/andrei-polukhin/pgdbtemplate-goose"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func main() {
	ctx := context.Background()

	// Create connection provider.
	connStringFunc := func(dbName string) string {
		return fmt.Sprintf("postgres://user:pass@localhost/%s?sslmode=disable", dbName)
	}
	provider := pgdbtemplatepq.NewConnectionProvider(connStringFunc)

	// Create goose migration runner from embedded filesystem.
	migrationRunner := pgdbtemplategoose.NewMigrationRunner(migrationsFS)

	// Create template manager.
	config := pgdbtemplate.Config{
		ConnectionProvider: provider,
		MigrationRunner:    migrationRunner,
	}

	tm, err := pgdbtemplate.NewTemplateManager(config)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize and use...
	if err := tm.Initialize(ctx); err != nil {
		log.Fatal(err)
	}
	defer tm.Cleanup(ctx)

	testDB, dbName, err := tm.CreateTestDatabase(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer testDB.Close()
	defer tm.DropTestDatabase(ctx, dbName)

	log.Printf("Test database %s ready!", dbName)
}
```

## Advanced Configuration

```go
// With custom goose options.
migrationsFs := os.DirFS("./migrations")
runner := pgdbtemplategoose.NewMigrationRunner(
	migrationsFs,
	pgdbtemplategoose.WithDialect(goose.DialectPostgres),
	pgdbtemplategoose.WithGooseOptions(
		goose.WithVerbose(true),
		goose.WithAllowOutofOrder(true),
	),
)
```

### Custom fs.FS Implementation

You can provide any `fs.FS` implementation:

```go
import (
	"io/fs"
	"os"
)

// From filesystem
migrationsFs := os.DirFS("./db/migrations")

// From embedded filesystem
//go:embed sql/migrations/*.sql
var embeddedMigrations embed.FS

// From sub-directory of embedded filesystem
subFs, _ := fs.Sub(embeddedMigrations, "sql/migrations")

// Use any of them
runner := pgdbtemplategoose.NewMigrationRunner(migrationsFs)
```

## Examples

See the [`examples/`](examples/) directory for complete working examples:

- [`examples/filesystem/`](examples/filesystem/) - Using `os.DirFS()` to load migrations from disk
- [`examples/embed/`](examples/embed/) - Using `embed.FS` to bundle migrations into binary

## Requirements

- Go 1.21+
- PostgreSQL 9.5+
- **Driver**: Works with both `pgdbtemplate-pq` (database/sql) and
  `pgdbtemplate-pgx` (pgx/v5)

## License

MIT license.
