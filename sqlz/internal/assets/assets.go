package assets

import (
	"embed"
)

// Embedded assets.
var (
	//go:embed migrations
	MigrationsAssetFS embed.FS
)
