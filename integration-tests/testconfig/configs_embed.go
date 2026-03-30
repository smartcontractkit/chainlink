//go:build embed
// +build embed

package testconfig

import "embed"

//go:embed default.toml
//go:embed ccip/ccip.toml

var embeddedConfigsFs embed.FS

func init() {
	areConfigsEmbedded = true
}
