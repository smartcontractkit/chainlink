package pkg

import "embed"

//go:embed *tmpl
var tmplFS embed.FS //nolint:unused // reserved for future template loading; kept alongside the go:embed directive
