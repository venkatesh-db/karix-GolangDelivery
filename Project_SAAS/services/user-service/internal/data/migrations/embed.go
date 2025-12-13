package migrations

import "embed"

// Files holds the SQL migrations for the user service.
//
//go:embed *.sql
var Files embed.FS
