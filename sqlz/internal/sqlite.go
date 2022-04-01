package internal

import (
	"context"
	"testing"

	"github.com/ibrt/golang-fixtures/fixturez"

	"github.com/ibrt/golang-inject-sqlite/sqlz"
	"github.com/ibrt/golang-inject-sqlite/sqlz/internal/assets"
)

var (
	_ fixturez.BeforeSuite = &ConfigHelper{}
)

// ConfigHelper is a test helper for *Config.
type ConfigHelper struct {
	// intentionally empty
}

// BeforeSuite implements fixturez.BeforeSuite.
func (h *ConfigHelper) BeforeSuite(ctx context.Context, t *testing.T) context.Context {
	t.Helper()

	return sqlz.NewConfigSingletonInjector(&sqlz.Config{
		DBSpec:          sqlz.NewDefaultConfigDBSpec("unused"),
		ApplyMigrations: sqlz.NewConfigMigrations(assets.MigrationsAssetFS, "migrations"),
	})(ctx)
}
