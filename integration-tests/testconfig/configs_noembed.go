//go:build !embed
// +build !embed

package testconfig

import "embed"

var embeddedConfigsFs embed.FS

func init() {
	configsEmbedded = false
}
