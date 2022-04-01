package sqlz_test

import (
	"context"
	"testing"

	"github.com/ibrt/golang-errors/errorz"
	"github.com/ibrt/golang-fixtures/fixturez"
	"github.com/stretchr/testify/require"

	"github.com/ibrt/golang-inject-sqlite/sqlz"
	"github.com/ibrt/golang-inject-sqlite/sqlz/internal"
	"github.com/ibrt/golang-inject-sqlite/sqlz/internal/assets"
	"github.com/ibrt/golang-inject-sqlite/sqlz/testsqlz"
)

func TestModule(t *testing.T) {
	fixturez.RunSuite(t, &Suite{})
}

type Suite struct {
	*fixturez.DefaultConfigMixin
	SQLiteConfig *internal.ConfigHelper
	SQLite       *testsqlz.Helper
}

func (s *Suite) TestSQLite(ctx context.Context, t *testing.T) {
	migrations := sqlz.NewConfigMigrations(assets.MigrationsAssetFS, "migrations")

	_, err := sqlz.Get(ctx).ExecContext(ctx, `INSERT INTO first_table (id, value) VALUES ("id1", "value1")`)
	fixturez.RequireNoError(t, err)

	_, err = sqlz.GetCtx(ctx).Exec(`INSERT INTO second_table (id, value) VALUES ("id1", "value1")`)
	fixturez.RequireNoError(t, err)

	count, err := sqlz.RevertMigrations(sqlz.GetDB(ctx), migrations, 1)
	fixturez.RequireNoError(t, err)
	require.Equal(t, 1, count)

	_, err = sqlz.Get(ctx).ExecContext(ctx, `INSERT INTO first_table (id, value) VALUES ("id2", "value2")`)
	fixturez.RequireNoError(t, err)

	_, err = sqlz.GetCtx(ctx).Exec(`INSERT INTO second_table (id, value) VALUES ("id2", "value2")`)
	require.EqualError(t, err, "no such table: second_table")

	count, err = sqlz.ApplyMigrations(sqlz.GetDB(ctx), migrations, 1)
	fixturez.RequireNoError(t, err)
	require.Equal(t, 1, count)

	_, err = sqlz.Get(ctx).ExecContext(ctx, `INSERT INTO first_table (id, value) VALUES ("id3", "value3")`)
	fixturez.RequireNoError(t, err)

	_, err = sqlz.GetCtx(ctx).Exec(`INSERT INTO second_table (id, value) VALUES ("id2", "value2")`)
	fixturez.RequireNoError(t, err)

	rows, err := sqlz.GetCtx(ctx).Query(`SELECT 1`)
	fixturez.RequireNoError(t, err)
	defer errorz.IgnoreClose(rows)
	var i int64
	require.True(t, rows.Next())
	fixturez.RequireNoError(t, rows.Scan(&i))
	require.False(t, rows.Next())
	require.Equal(t, int64(1), i)

	row := sqlz.GetCtx(ctx).QueryRow(`SELECT 2`)
	fixturez.RequireNoError(t, row.Err())
	fixturez.RequireNoError(t, row.Scan(&i))
	fixturez.RequireNoError(t, row.Err())
	require.Equal(t, int64(2), i)

	dbSpec := sqlz.NewConfigDBSpec("filePath", map[string]interface{}{"k": "v"})
	require.Equal(t, "filePath", dbSpec.FilePath)
	require.Equal(t, map[string]interface{}{"k": "v"}, dbSpec.Pragmas)
}
