//go:build embed
// +build embed

package testconfig

import "embed"

//go:embed default.toml
//go:embed automation/automation.toml
//go:embed forwarder_ocr/forwarder_ocr.toml
//go:embed forwarder_ocr2/forwarder_ocr2.toml
//go:embed functions/functions.toml
//go:embed keeper/keeper.toml
//go:embed log_poller/log_poller.toml
//go:embed node/node.toml
//go:embed ocr/ocr.toml
//go:embed ocr2/ocr2.toml
//go:embed vrf/vrf.toml
//go:embed vrfv2/vrfv2.toml
//go:embed vrfv2plus/vrfv2plus.toml
//go:embed ccip/ccip.toml

var embeddedConfigsFs embed.FS

func init() {
	areConfigsEmbedded = true
}
