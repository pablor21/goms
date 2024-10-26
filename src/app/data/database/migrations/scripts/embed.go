package scripts

import (
	"embed"
)

//go:embed *.sql
var Scripts embed.FS
