package docs

import "embed"

//go:embed all:swagger.json
var ApiDocs embed.FS
