//go:build embed
// +build embed

package testconfig

import "embed"

//go:embed default.toml
//go:embed automation/automation.toml
//go:embed functions/functions.toml
//go:embed keeper/keeper.toml
//go:embed log_poller/log_poller.toml
//go:embed node/node.toml
//go:embed ocr/ocr.toml
//go:embed vrfv2/vrfv2.toml
//go:embed vrfv2plus/vrfv2plus.toml
var embeddedConfigsFs embed.FS

func init() {
	areConfigsEmbedded = true
}
