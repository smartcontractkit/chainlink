package validate

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	medianconfig "github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/median/config"
	"github.com/stretchr/testify/assert"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/services/job"
)

func TestValidateOracleSpec(t *testing.T) {
	var tt = []struct {
		name       string
		toml       string
		setGlobals func(t *testing.T, c *configtest.TestGeneralConfig)
		assertion  func(t *testing.T, os job.Job, err error)
	}{
		{
			name: "minimal OCR2 oracle spec",
			toml: `
type               = "offchainreporting2"
pluginType         = "median"
schemaVersion      = 1
relay              = "evm"
contractID         = "0x613a38AC1659769640aaE063C651F48E0250454C"
observationSource  = """
ds1          [type=bridge name=voter_turnout];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
answer1      [type=median index=0];
"""
[relayConfig]
chainID = 1337
[pluginConfig]
juelsPerFeeCoinSource = """
ds1          [type=bridge name=voter_turnout];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
answer1      [type=median index=0];
"""
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.NoError(t, err)
				// Should be able to jsonapi marshal/unmarshal the minimum spec.
				// This ensures the UnmarshalJSON's defined on the fields handle a min spec correctly.
				b, err := jsonapi.Marshal(os.OCR2OracleSpec)
				require.NoError(t, err)
				var r job.OCR2OracleSpec
				err = jsonapi.Unmarshal(b, &r)
				require.NoError(t, err)
				assert.Equal(t, "median", string(r.PluginType))
				var pc medianconfig.PluginConfig
				require.NoError(t, json.Unmarshal(r.PluginConfig.Bytes(), &pc))
				require.NoError(t, medianconfig.ValidatePluginConfig(pc))
			},
		},
		{
			name: "decodes valid oracle spec toml",
			toml: `
type               = "offchainreporting2"
pluginType         = "median"
schemaVersion      = 1
relay              = "evm"
contractID         = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = [
"12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001",
]
ocrKeyBundleID     = "73e8966a78ca09bb912e9565cfb79fbe8a6048fab1f0cf49b18047c3895e0447"
monitoringEndpoint = "chain.link:4321"
transmitterID = "0xF67D0290337bca0847005C7ffD1BC75BA9AAE6e4"
observationTimeout = "10s"
observationSource  = """
ds1          [type=bridge name=voter_turnout];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
answer1      [type=median index=0];
"""
[relayConfig]
chainID = 1337
[pluginConfig]
juelsPerFeeCoinSource = """
ds1          [type=bridge name=voter_turnout];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
answer1      [type=median index=0];
"""
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.NoError(t, err)
				assert.Equal(t, 1, int(os.SchemaVersion))
			},
		},
		{
			name: "raises error on extra keys",
			toml: `
type               = "offchainreporting2"
pluginType         = "median"
schemaVersion      = 1
relay              = "evm"
contractID         = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = [
"12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001",
]
isBootstrapPeer    = true
ocrKeyBundleID     = "73e8966a78ca09bb912e9565cfb79fbe8a6048fab1f0cf49b18047c3895e0447"
monitoringEndpoint = "chain.link:4321"
transmitterID      = "0xF67D0290337bca0847005C7ffD1BC75BA9AAE6e4"
observationTimeout = "10s"
observationSource  = """
ds1          [type=bridge name=voter_turnout];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
answer1      [type=median index=0];
"""
[relayConfig]
chainID = 1337
[pluginConfig]
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "unrecognised key for ocr2 peer: isBootstrapPeer")
			},
		},
		{
			name: "empty pipeline string",
			toml: `
type               = "offchainreporting2"
pluginType         = "median"
schemaVersion      = 1
relay              = "evm"
contractID         = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = []
[relayConfig]
chainID = 1337
[pluginConfig]
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "invalid dot",
			toml: `
type               = "offchainreporting2"
pluginType         = "median"
schemaVersion      = 1
relay              = "evm"
contractID         = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = []
observationSource = """
->
"""
[relayConfig]
chainID = 1337
[pluginConfig]
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "invalid peer address",
			toml: `
type               = "offchainreporting2"
pluginType         = "median"
schemaVersion      = 1
relay              = "evm"
contractID         = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = ["/invalid/peer/address"]
observationSource = """
blah
"""
[relayConfig]
chainID = 1337
[pluginConfig]
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "non-zero timeouts",
			toml: `
type               = "offchainreporting2"
pluginType         = "median"
schemaVersion      = 1
relay              = "evm"
contractID         = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = ["12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001"]
blockchainTimeout  = "0s"
observationSource = """
blah
"""
[relayConfig]
chainID = 1337
[pluginConfig]
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "non-zero intervals",
			toml: `
type               = "offchainreporting2"
pluginType         = "median"
schemaVersion      = 1
relay              = "evm"
contractID         = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = ["12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001"]
observationSource = """
blah
"""
[relayConfig]
chainID = 1337
[pluginConfig]
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "broken monitoring endpoint",
			toml: `
type               = "offchainreporting2"
pluginType         = "median"
schemaVersion      = 1
relay              = "evm"
contractID         = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = []
monitoringEndpoint = "\t/fd\2ff )(*&^%$#@"
[relayConfig]
chainID = 1337
[pluginConfig]
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "invalid escape sequence")
			},
		},
		{
			name: "toml parse doesn't panic",
			toml: string(hexutil.MustDecode("0x2222220d5c22223b22225c0d21222222")),
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "invalid global default",
			toml: `
type               = "offchainreporting2"
pluginType         = "median"
schemaVersion      = 1
maxTaskDuration    = "30m"
relay              = "evm"
contractID         = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = [
"12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001",
]
monitoringEndpoint = "chain.link:4321"
transmitterID = "0xF67D0290337bca0847005C7ffD1BC75BA9AAE6e4"
observationSource  = """
ds1          [type=bridge name=voter_turnout];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
answer1      [type=median index=0];
"""
[pluginConfig]
juelsPerFeeCoinSource = """
ds1          [type=bridge name=voter_turnout];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
answer1      [type=median index=0];
"""
[relayConfig]
chainID = 1337
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "database timeout must be between 100ms and 10s, but is currently 20m0s")
			},
			setGlobals: func(t *testing.T, c *configtest.TestGeneralConfig) {
				to := 20 * time.Minute
				c.Overrides.OCR2DatabaseTimeout = &to
			},
		},
		{
			name: "invalid pluginType",
			toml: `
type               = "offchainreporting2"
pluginType         = "medion"
schemaVersion      = 1
relay              = "evm"
contractID         = "0x613a38AC1659769640aaE063C651F48E0250454C"
observationSource  = """
ds1          [type=bridge name=voter_turnout];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
answer1      [type=median index=0];
"""
[pluginConfig]
juelsPerFeeCoinSource = """
->
"""
[relayConfig]
chainID = 1337
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "invalid pluginType medion")
			},
		},
		{
			name: "invalid relay",
			toml: `
type               = "offchainreporting2"
pluginType         = "median"
schemaVersion      = 1
relay              = "blerg"
contractID         = "0x613a38AC1659769640aaE063C651F48E0250454C"
observationSource  = """
ds1          [type=bridge name=voter_turnout];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
answer1      [type=median index=0];
"""
[pluginConfig]
juelsPerFeeCoinSource = """
ds1          [type=bridge name=voter_turnout];
"""
[relayConfig]
chainID = 1337
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				fmt.Println("relay", os.OCR2OracleSpec.Relay)
				require.Error(t, err)
				require.Contains(t, err.Error(), "no such relay blerg supported")
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := configtest.NewTestGeneralConfig(t)
			c.Overrides.Dev = null.BoolFrom(false)
			c.Overrides.EVMRPCEnabled = null.BoolFrom(false)
			if tc.setGlobals != nil {
				tc.setGlobals(t, c)
			}
			s, err := ValidatedOracleSpecToml(c, tc.toml)
			tc.assertion(t, s, err)
		})
	}
}
