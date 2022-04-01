package sqlz

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"net/url"

	"github.com/ibrt/golang-errors/errorz"
	"github.com/ibrt/golang-inject/injectz"
	"github.com/ibrt/golang-validation/vz"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

type contextKey int

const (
	dbContextKey contextKey = iota
	configContextKey
)

var (
	// DefaultPragmas provides a default, sane set of PRAGMAs for opening SQLite databases.
	DefaultPragmas = map[string]interface{}{
		"auto_vacuum":              2, // INCREMENTAL
		"cache_size":               -2000,
		"foreign_keys":             1,
		"ignore_check_constraints": 0,
		"locking_mode":             "normal",
		"synchronous":              2, // FULL
		"writable_schema":          0,
	}
)

// Config describes the configuration for SQLite.
type Config struct {
	DBSpec          *ConfigDBSpec     `json:"dbSpec" validate:"required"`
	ApplyMigrations *ConfigMigrations `json:"-"`
}

// ConfigDBSpec describes part of the configuration for SQLite.
type ConfigDBSpec struct {
	FilePath string                 `json:"filePath" validate:"required"`
	Pragmas  map[string]interface{} `json:"pragmas" validate:"required"`
}

// ToOpen converts the *ConfigDBSpec to a string that can be passed to sql.Open().
func (c *ConfigDBSpec) ToOpen() string {
	query := url.Values{}
	for k, v := range c.Pragmas {
		query.Set("_"+k, fmt.Sprintf("%v", v))
	}
	return fmt.Sprintf("%v?%v", c.FilePath, query.Encode())
}

// NewConfigDBSpec initializes a new *ConfigDBSpec.
func NewConfigDBSpec(filePath string, pragmas map[string]interface{}) *ConfigDBSpec {
	return &ConfigDBSpec{
		FilePath: filePath,
		Pragmas:  pragmas,
	}
}

// NewDefaultConfigDBSpec initializes a new *ConfigDBSpec with the default PRAGMAs.
func NewDefaultConfigDBSpec(filePath string) *ConfigDBSpec {
	return &ConfigDBSpec{
		FilePath: filePath,
		Pragmas:  DefaultPragmas,
	}
}

// ConfigMigrations describes part of the configuration for SQLite.
type ConfigMigrations struct {
	EmbedFS                embed.FS `json:"-" validate:"required"`
	EmbedMigrationsDirPath string   `json:"embedMigrationsDirPath" validate:"required"`
}

// NewConfigMigrations initializes a new *ConfigMigrations.
func NewConfigMigrations(embedFS embed.FS, embedMigrationsDirPath string) *ConfigMigrations {
	return &ConfigMigrations{
		EmbedFS:                embedFS,
		EmbedMigrationsDirPath: embedMigrationsDirPath,
	}
}

// Validate implements the vz.Validator interface.
func (c *Config) Validate() error {
	return errorz.MaybeWrap(vz.ValidateStruct(c), errorz.SkipPackage())
}

// NewConfigSingletonInjector always inject the given *Config.
func NewConfigSingletonInjector(cfg *Config) injectz.Injector {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, configContextKey, cfg)
	}
}

// GetConfig extracts the *Config from context, panics if not found.
func GetConfig(ctx context.Context) *Config {
	return ctx.Value(configContextKey).(*Config)
}

// SQLite describes the sqlite module (a subset of *sql.DB).
type SQLite interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// ContextSQLite describes a SQLite with a cached context.
type ContextSQLite interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type contextSQLiteImpl struct {
	ctx context.Context
	pg  SQLite
}

// Exec executes a query.
func (p *contextSQLiteImpl) Exec(query string, args ...interface{}) (sql.Result, error) {
	return p.pg.ExecContext(p.ctx, query, args...)
}

// Query executes a query.
func (p *contextSQLiteImpl) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return p.pg.QueryContext(p.ctx, query, args...)
}

// QueryRow executes a query.
func (p *contextSQLiteImpl) QueryRow(query string, args ...interface{}) *sql.Row {
	return p.pg.QueryRowContext(p.ctx, query, args...)
}

// Initializer is a SQLite initializer.
func Initializer(ctx context.Context) (injectz.Injector, injectz.Releaser) {
	cfg := ctx.Value(configContextKey).(*Config)
	errorz.MaybeMustWrap(cfg.Validate(), errorz.SkipPackage())

	db, err := sql.Open("sqlite3", cfg.DBSpec.ToOpen())
	errorz.MaybeMustWrap(err, errorz.SkipPackage())
	errorz.MaybeMustWrap(db.Ping(), errorz.SkipPackage())

	if cfg.ApplyMigrations != nil {
		_, err = ApplyMigrations(db, cfg.ApplyMigrations, AllMigrations)
		errorz.MaybeMustWrap(err, errorz.SkipPackage())
	}

	injector := func(ctx context.Context) context.Context {
		return context.WithValue(ctx, dbContextKey, db)
	}

	releaser := func() {
		errorz.IgnoreClose(db)
	}

	return injector, releaser
}

// Get extracts the SQLite from context, panics if not found.
func Get(ctx context.Context) SQLite {
	return ctx.Value(dbContextKey).(SQLite)
}

// GetCtx extracts the SQLite from context and wraps it as ContextSQLite, panics if not found.
func GetCtx(ctx context.Context) ContextSQLite {
	return &contextSQLiteImpl{
		ctx: ctx,
		pg:  Get(ctx),
	}
}

// GetDB extracts the SQLite from context as *sql.DB, panics if not found.
func GetDB(ctx context.Context) *sql.DB {
	return ctx.Value(dbContextKey).(*sql.DB)
}
