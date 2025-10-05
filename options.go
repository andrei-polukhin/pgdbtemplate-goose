package pgdbtemplategoose

import "github.com/pressly/goose/v3"

// Option configures the goose migration runner.
type Option func(*MigrationRunner)

// WithDialect sets the SQL dialect for goose migrations.
// Default is "postgres".
func WithDialect(dialect goose.Dialect) Option {
	return func(r *MigrationRunner) {
		r.dialect = dialect
	}
}

// WithGooseOptions sets additional goose provider options.
//
// Example:
//
//	runner := NewMigrationRunner(
//	    "./migrations",
//	    WithGooseOptions(
//	        goose.WithAllowMissing(),
//	        goose.WithNoVersioning(),
//	    ),
//	)
func WithGooseOptions(opts ...goose.ProviderOption) Option {
	return func(r *MigrationRunner) {
		r.opts = append(r.opts, opts...)
	}
}
