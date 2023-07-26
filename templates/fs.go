package templates

import "embed"

// embed is used to provide a pattern which is going to be all the files that we want it to include

// FS is a file system here provided by the embed and this is going to be embedded into our binary

//go:embed *.gohtml
var FS embed.FS
