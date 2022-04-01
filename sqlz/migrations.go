package sqlz

import (
	"database/sql"

	"github.com/ibrt/golang-errors/errorz"
	migrate "github.com/rubenv/sql-migrate"
)

const (
	// AllMigrations means apply/revert all migrations.
	AllMigrations = -1
)

// ApplyMigrations applies the given migrations.
func ApplyMigrations(db *sql.DB, migrations *ConfigMigrations, max int) (int, error) {
	return execMigrations(db, migrations, migrate.Up, max)
}

// RevertMigrations reverts the given migrations.
func RevertMigrations(db *sql.DB, migrations *ConfigMigrations, max int) (int, error) {
	return execMigrations(db, migrations, migrate.Down, max)
}

func execMigrations(db *sql.DB, migrations *ConfigMigrations, dir migrate.MigrationDirection, max int) (int, error) {
	migrationSource := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: migrations.EmbedFS,
		Root:       migrations.EmbedMigrationsDirPath,
	}

	count, err := migrate.ExecMax(db, "sqlite3", migrationSource, dir, max)
	return count, errorz.MaybeWrap(err, errorz.SkipPackage())
}
