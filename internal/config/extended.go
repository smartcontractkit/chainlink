package config

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pelletier/go-toml/v2"

	evmcfg "github.com/smartcontractkit/chainlink/core/chains/evm/config/v2"
)

func extended(key string) (string, error) {
	if key != "EVM" {
		return "", fmt.Errorf("%s: no extended description available", key)
	}

	// EVM Per-Chain Defaults
	var sb strings.Builder
	for _, id := range evmcfg.DefaultIDs {
		fmt.Fprintf(&sb, "<details><summary>%s) %s</summary><p>\n\n", id, evmcfg.DefaultName(id))
		sb.WriteString("```toml\n")
		var config evmcfg.Chain
		config.SetDefaults(id)
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
