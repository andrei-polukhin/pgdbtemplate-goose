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
- 📁 **Standard Migrations** - Use your existing goose migration files
- ⚙️ **Configurable** - Support for custom goose options and dialects
- 🧪 **Well-tested** - 75%+ test coverage with real PostgreSQL integration

## Installation

```bash
go get github.com/andrei-polukhin/pgdbtemplate-goose
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/andrei-polukhin/pgdbtemplate"
	pgdbtemplatepq "github.com/andrei-polukhin/pgdbtemplate-pq"
	pgdbtemplategoose "github.com/andrei-polukhin/pgdbtemplate-goose"
)

func main() {
	ctx := context.Background()

	// Create connection provider (must use pq for goose).
	connStringFunc := func(dbName string) string {
		return fmt.Sprintf("postgres://user:pass@localhost/%s?sslmode=disable", dbName)
	}
	provider := pgdbtemplatepq.NewConnectionProvider(connStringFunc)

	// Create goose migration runner.
	migrationRunner := pgdbtemplategoose.NewMigrationRunner("./migrations")

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

## Advanced Configuration

```go
// With custom goose options.
runner := pgdbtemplategoose.NewMigrationRunner(
	"./migrations",
	pgdbtemplategoose.WithDialect(goose.DialectPostgres),
	pgdbtemplategoose.WithGooseOptions(
		goose.WithAllowMissing(),
	),
)
```

## Requirements

- Go 1.21+
- PostgreSQL 9.5+
- **Note**: This adapter requires `pgdbtemplate-pq` (database/sql) driver.
  It is not compatible with `pgdbtemplate-pgx` due to goose's database/sql dependency.

## License

MIT license.
