package docs

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pelletier/go-toml/v2"

	evmcfg "github.com/smartcontractkit/chainlink-integrations/evm/config/toml"
)

// evmChainDefaults returns generated Markdown for the EVM per-chain defaults. See v2.Defaults.
func evmChainDefaults() (string, error) {
	var sb strings.Builder
	for _, id := range evmcfg.DefaultIDs {
		config, name := evmcfg.DefaultsNamed(id)
		fmt.Fprintf(&sb, "\n<details><summary>%s (%s)</summary><p>\n\n", name, id)
		sb.WriteString("```toml\n")
		b, err := toml.Marshal(config)
		if err != nil {
			return "", err
		}
		sb.Write(bytes.TrimSpace(b))
		sb.WriteString("\n```\n\n")
		sb.WriteString("</p></details>\n")
	}

	return sb.String(), nil
}
