package ocrbootstrap

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateBootstrapSpec(t *testing.T) {
	var tt = []struct {
		name         string
		toml         string
		setGlobalCfg func(t *testing.T, c *configtest.TestGeneralConfig)
		assertion    func(t *testing.T, os job.Job, err error)
	}{
		{
			name: "decodes valid bootstrap spec toml",
			toml: `
type				= "bootstrap"
name				= "bootstrap"
schemaVersion		= 1
contractID			= "0x613a38AC1659769640aaE063C651F48E0250454C"
monitoringEndpoint	= "chain.link:4321"
relay				= "evm"
[relayConfig]
chainID 			= 1337
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.NoError(t, err)
				assert.Equal(t, 1, int(os.SchemaVersion))
			},
		},
		{
			name: "raises error on missing key",
			toml: `
type				= "bootstrap"
schemaVersion		= 1
monitoringEndpoint	= "chain.link:4321"
relay				= "evm"
[relayConfig]
chainID 			= 1337
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "missing required key contractID")
			},
		},
		{
			name: "raises error on unexpected key",
			toml: `
type				= "bootstrap"
schemaVersion		= 1
contractID			= "0x613a38AC1659769640aaE063C651F48E0250454C"
monitoringEndpoint	= "chain.link:4321"
isBootstrapPeer		= true
relay				= "evm"
[relayConfig]
chainID			= 1337
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "unrecognised key for bootstrap peer: isBootstrapPeer")
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := ValidatedBootstrapSpecToml(tc.toml)
			tc.assertion(t, s, err)
		})
	}
}
