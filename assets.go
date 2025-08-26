package main

import "embed"

// embeddedStatic holds all files under the ./static directory.
// The directory structure is preserved for serving via http.FS.
//
//go:embed static
var embeddedStatic embed.FS
