package migrations

import "embed"

//go:embed *.sql
var FS embed.FS

// any sql files in this directory will be embedded into the binary.
// This makes it much easier to run migrations anywhere without also copying the .sql files around.
