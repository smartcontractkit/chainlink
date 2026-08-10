package ui

import (
	_ "embed"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
)

//go:embed web/index.html
var webIndexHTML string

//go:embed web/styles.css
var webStylesCSS string

//go:embed web/app.js
var webAppJS string

//go:embed web/capability_defaults.toml
var capabilityDefaultsTOML string

var webFiles = map[string]string{
	"index.html": webIndexHTML,
	"styles.css": webStylesCSS,
	"app.js":     webAppJS,
}

func LoadCapabilityDefaults() map[string]domain.CapabilityConfig {
	return domain.LoadCapabilityDefaults(capabilityDefaultsTOML)
}
