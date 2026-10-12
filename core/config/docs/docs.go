package docs

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/config/configdoc"
)

var (
	//go:embed secrets.toml
	secretsTOML string
	//go:embed core.toml
	coreTOML string
	//go:embed chains-evm.toml
	chainsEVMTOML string
	//go:embed chains-solana.toml
	chainsSolanaTOML string
	//go:embed chains-starknet.toml
	chainsStarknetTOML string

	//go:embed example-config.toml
	exampleConfig string
	//go:embed example-secrets.toml
	exampleSecrets string

	docsTOML = coreTOML + chainsEVMTOML + chainsSolanaTOML + chainsStarknetTOML
)

// GenerateConfig returns MarkDown documentation generated from core.toml & chains-*.toml.
func GenerateConfig() (string, error) {
	evmDefaults, err := evmChainDefaults()
	if err != nil {
		return "", fmt.Errorf("failed to generate evm chain defaults: %w", err)
	}
	content, err := configdoc.Generate(docsTOML, `[//]: # (Documentation generated from docs/*.toml - DO NOT EDIT.)

This document describes the TOML format for configuration.

See also [SECRETS.md](SECRETS.md)
`, exampleConfig, map[string]string{"EVM": evmDefaults})
	if err != nil {
		return "", err
	}
	return strings.TrimRight(content, "\r\n") + "\n", nil
}

// GenerateSecrets returns MarkDown documentation generated from secrets.toml.
func GenerateSecrets() (string, error) {
	content, err := configdoc.Generate(secretsTOML, `[//]: # (Documentation generated from docs/secrets.toml - DO NOT EDIT.)

This document describes the TOML format for secrets.

Each secret has an alternative corresponding environment variable.

See also [CONFIG.md](CONFIG.md)
`, exampleSecrets, nil)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(content, "\r\n") + "\n", nil
}
