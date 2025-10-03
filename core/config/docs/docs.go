package docs

import (
	_ "embed"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/config/configdoc"
)

var (
	//go:embed secrets.toml
	secretsTOML string
	//go:embed core.toml
	coreTOML string
	//go:embed chains-evm.toml
	chainsEVMTOML string
	//go:embed chains-cosmos.toml
	chainsCosmosTOML string
	//go:embed chains-solana.toml
	chainsSolanaTOML string
	//go:embed chains-starknet.toml
	chainsStarknetTOML string

	//go:embed example-config.toml
	exampleConfig string
	//go:embed example-secrets.toml
	exampleSecrets string

	docsTOML = coreTOML + chainsEVMTOML + chainsCosmosTOML + chainsSolanaTOML + chainsStarknetTOML
)

// GenerateConfig returns MarkDown documentation generated from core.toml & chains-*.toml.
func GenerateConfig() (string, error) {
	evmDefaults, err := evmChainDefaults()
	if err != nil {
		return "", fmt.Errorf("failed to generate evm chain defaults: %w", err)
	}
	return configdoc.Generate(docsTOML, `[//]: # (Documentation generated from docs/*.toml - DO NOT EDIT.)

This document describes the TOML format for configuration.

See also [SECRETS.md](SECRETS.md)
`, exampleConfig, map[string]string{"EVM": evmDefaults})
}

// GenerateSecrets returns MarkDown documentation generated from secrets.toml.
func GenerateSecrets() (string, error) {
	return configdoc.Generate(secretsTOML, `[//]: # (Documentation generated from docs/secrets.toml - DO NOT EDIT.)

This document describes the TOML format for secrets.

Each secret has an alternative corresponding environment variable.

See also [CONFIG.md](CONFIG.md)
`, exampleSecrets, nil)
}
