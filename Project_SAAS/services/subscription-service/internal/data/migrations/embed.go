package migrations

import "embed"

// Files exposes the embedded SQL migrations for the subscription service.
//
//go:embed *.sql
var Files embed.FS
