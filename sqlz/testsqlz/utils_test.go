package testsqlz_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/ibrt/golang-fixtures/fixturez"
	"github.com/stretchr/testify/require"

	"github.com/ibrt/golang-inject-sqlite/sqlz"
	"github.com/ibrt/golang-inject-sqlite/sqlz/internal/assets"
	"github.com/ibrt/golang-inject-sqlite/sqlz/testsqlz"
)

func TestUtils(t *testing.T) {
	fd, err := ioutil.TempFile("", "golang-inject-sqlite")
	fixturez.RequireNoError(t, err)
	tmpFilePath := fd.Name()
	defer func() {
		fixturez.RequireNoError(t, os.RemoveAll(tmpFilePath))
	}()
	fixturez.RequireNoError(t, fd.Close())

	dbSpec := sqlz.NewDefaultConfigDBSpec(tmpFilePath)
	migrations := sqlz.NewConfigMigrations(assets.MigrationsAssetFS, "migrations")

	require.NotPanics(t, func() {
		testsqlz.MustApplyMigrations(dbSpec, migrations, -1)
	})

	require.NotPanics(t, func() {
		db := testsqlz.MustOpen(dbSpec)
		defer func() {
			fixturez.RequireNoError(t, db.Close())
		}()

		_, err := db.Exec(`INSERT INTO first_table (id, value) VALUES ("id1", "value1")`)
		fixturez.RequireNoError(t, err)

		_, err = db.Exec(`INSERT INTO second_table (id, value) VALUES ("id1", "value1")`)
		fixturez.RequireNoError(t, err)
	})

	require.NotPanics(t, func() {
		testsqlz.MustRevertMigrations(dbSpec, migrations, 1)
	})

	require.NotPanics(t, func() {
		db := testsqlz.MustOpen(dbSpec)
		defer func() {
			fixturez.RequireNoError(t, db.Close())
		}()

		_, err := db.Exec(`INSERT INTO first_table (id, value) VALUES ("id2", "value2")`)
		fixturez.RequireNoError(t, err)

		_, err = db.Exec(`INSERT INTO second_table (id, value) VALUES ("id2", "value2")`)
		require.EqualError(t, err, "no such table: second_table")
	})
}
