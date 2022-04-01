package testsqlz

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ibrt/golang-errors/errorz"
	"github.com/ibrt/golang-fixtures/fixturez"

	"github.com/ibrt/golang-inject-sqlite/sqlz"
)

var (
	_ fixturez.BeforeSuite = &Helper{}
	_ fixturez.AfterSuite  = &Helper{}
)

// Helper is a test helper for PG.
type Helper struct {
	origDBSpec  *sqlz.ConfigDBSpec
	tmpFilePath string
	releaser    func()
}

// BeforeSuite implements fixturez.BeforeSuite.
func (h *Helper) BeforeSuite(ctx context.Context, t *testing.T) context.Context {
	t.Helper()

	fd, err := ioutil.TempFile("", "golang-inject-sqlite")
	errorz.MaybeMustWrap(err, errorz.SkipPackage())
	h.tmpFilePath = fd.Name()
	errorz.MaybeMustWrap(fd.Close(), errorz.SkipPackage())

	cfg := sqlz.GetConfig(ctx)
	h.origDBSpec = cfg.DBSpec
	cfg.DBSpec.FilePath = h.tmpFilePath

	injector, releaser := sqlz.Initializer(ctx)
	h.releaser = releaser
	return injector(ctx)
}

// AfterSuite implements fixturez.AfterSuite.
func (h *Helper) AfterSuite(ctx context.Context, t *testing.T) {
	t.Helper()

	h.releaser()
	h.releaser = nil

	errorz.MaybeMustWrap(os.RemoveAll(h.tmpFilePath), errorz.SkipPackage())
	h.tmpFilePath = ""

	cfg := sqlz.GetConfig(ctx)
	cfg.DBSpec = h.origDBSpec
	h.origDBSpec = nil
}
