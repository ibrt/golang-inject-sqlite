package testsqlz_test

import (
	"context"
	"testing"

	"github.com/ibrt/golang-fixtures/fixturez"

	"github.com/ibrt/golang-inject-sqlite/sqlz"
	"github.com/ibrt/golang-inject-sqlite/sqlz/internal"
	"github.com/ibrt/golang-inject-sqlite/sqlz/testsqlz"
)

func TestHelpers(t *testing.T) {
	fixturez.RunSuite(t, &Suite{})
}

type Suite struct {
	*fixturez.DefaultConfigMixin
	SQLiteConfig *internal.ConfigHelper
	SQLite       *testsqlz.Helper
}

func (s *Suite) TestHelper(ctx context.Context, t *testing.T) {
	_, err := sqlz.Get(ctx).ExecContext(ctx, `INSERT INTO first_table (id, value) VALUES ("id", "value")`)
	fixturez.RequireNoError(t, err)

	_, err = sqlz.Get(ctx).ExecContext(ctx, `INSERT INTO second_table (id, value) VALUES ("id", "value")`)
	fixturez.RequireNoError(t, err)
}
