package main

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/config/docs"
)

var (
	//go:embed docs/CONFIG.md
	configMD string
	//go:embed docs/SECRETS.md
	secretsMD string
)

func TestConfigDocs(t *testing.T) {
	config, err := docs.GenerateConfig()
	assert.NoError(t, err, "invalid config docs")
	assert.Equal(t, configMD, config, "docs/CONFIG.md is out of date. Run 'make config-docs' to regenerate.")

	secrets, err := docs.GenerateSecrets()
	assert.NoError(t, err, "invalid secrets docs")
	assert.Equal(t, secretsMD, secrets, "docs/SECRETS.md is out of date. Run 'make config-docs' to regenerate.")
}
