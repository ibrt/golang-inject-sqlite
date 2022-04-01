package testsqlz

import (
	"database/sql"

	"github.com/ibrt/golang-errors/errorz"

	"github.com/ibrt/golang-inject-sqlite/sqlz"
)

// MustOpen opens the SQLite database, panics on error.
func MustOpen(dbSpec *sqlz.ConfigDBSpec) *sql.DB {
	db, err := sql.Open("sqlite3", dbSpec.ToOpen())
	errorz.MaybeMustWrap(err, errorz.SkipPackage())
	errorz.MaybeMustWrap(db.Ping(), errorz.SkipPackage())
	return db
}

// MustApplyMigrations applies the given migrations, panics on error.
func MustApplyMigrations(dbSpec *sqlz.ConfigDBSpec, migrations *sqlz.ConfigMigrations, max int) int {
	db := MustOpen(dbSpec)
	defer errorz.IgnoreClose(db)
	count, err := sqlz.ApplyMigrations(db, migrations, max)
	errorz.MaybeMustWrap(err, errorz.SkipPackage())
	return count
}

// MustRevertMigrations reverts the given migrations, panics on error.
func MustRevertMigrations(dbSpec *sqlz.ConfigDBSpec, migrations *sqlz.ConfigMigrations, max int) int {
	db := MustOpen(dbSpec)
	defer errorz.IgnoreClose(db)
	count, err := sqlz.RevertMigrations(db, migrations, max)
	errorz.MaybeMustWrap(err, errorz.SkipPackage())
	return count
}
