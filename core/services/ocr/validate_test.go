package ocr

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func TestValidateOracleSpec(t *testing.T) {
	var tt = []struct {
		name         string
		toml         string
		setGlobalCfg func(t *testing.T, c *configtest.TestGeneralConfig)
		assertion    func(t *testing.T, os job.Job, err error)
	}{
		{
			name: "minimal non-bootstrap oracle spec",
			toml: `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
isBootstrapPeer    = false
observationSource = """
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
				b, err := jsonapi.Marshal(os.OCROracleSpec)
				require.NoError(t, err)
				var r job.OCROracleSpec
				err = jsonapi.Unmarshal(b, &r)
				require.NoError(t, err)
			},
		},
		{
			name: "decodes valid oracle spec toml",
			toml: `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = [
"/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju",
]
isBootstrapPeer    = false
keyBundleID        = "73e8966a78ca09bb912e9565cfb79fbe8a6048fab1f0cf49b18047c3895e0447"
monitoringEndpoint = "chain.link:4321"
transmitterAddress = "0xF67D0290337bca0847005C7ffD1BC75BA9AAE6e4"
observationTimeout = "10s"
databaseTimeout = "2s"
observationGracePeriod = "2s"
contractTransmitterTransmitTimeout = "1s"
observationSource = """
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
				assert.False(t, os.OCROracleSpec.IsBootstrapPeer)
			},
		},
		{
			name: "decodes bootstrap toml",
			toml: `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = []
isBootstrapPeer    = true
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.NoError(t, err)
				assert.Equal(t, 1, int(os.SchemaVersion))
				assert.True(t, os.OCROracleSpec.IsBootstrapPeer)
			},
		},
		{
			name: "raises error on extra keys",
			toml: `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = [
"/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju",
]
isBootstrapPeer    = true
keyBundleID        = "73e8966a78ca09bb912e9565cfb79fbe8a6048fab1f0cf49b18047c3895e0447"
monitoringEndpoint = "chain.link:4321"
transmitterAddress = "0xF67D0290337bca0847005C7ffD1BC75BA9AAE6e4"
observationTimeout = "10s"
observationSource = """
ds1          [type=bridge name=voter_turnout];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
answer1      [type=median index=0];
"""
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "unrecognised key for bootstrap peer: observationSource")
			},
		},
		{
			name: "empty pipeline string non-bootstrap node",
			toml: `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = []
isBootstrapPeer    = false
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "invalid dot",
			toml: `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = []
isBootstrapPeer    = false
observationSource = """
->
"""
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "invalid peer address",
			toml: `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = ["/invalid/peer/address"]
isBootstrapPeer    = false
observationSource = """
blah
"""
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "non-zero blockchain timeout",
			toml: `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = ["/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju"]
isBootstrapPeer    = false
blockchainTimeout  = "0s"
observationSource = """
blah
"""
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "non-zero database timeout",
			toml: `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = ["/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju"]
isBootstrapPeer    = false
databaseTimeout  = "0s"
observationSource = """
blah
"""
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "non-zero observation grace period",
			toml: `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = ["/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju"]
isBootstrapPeer    = false
observationGracePeriod = "0s"
observationSource = """
blah
"""
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "non-zero contract transmitter transmit timeout",
			toml: `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = ["/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju"]
isBootstrapPeer    = false
contractTransmitterTransmitTimeout = "0s"
observationSource = """
blah
"""
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "non-zero intervals",
			toml: `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = ["/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju"]
isBootstrapPeer    = false
contractConfigTrackerSubscribeInterval = "0s"
observationSource = """
blah
"""
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "broken monitoring endpoint",
			toml: `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = []
isBootstrapPeer    = true
monitoringEndpoint = "\t/fd\2ff )(*&^%$#@"
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.EqualError(t, err, "toml error on load: (8, 23): invalid escape sequence: \\2")
			},
		},
		{
			name: "max task duration > observation timeout should error",
			toml: `
type               = "offchainreporting"
maxTaskDuration    = "30s"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = [
"/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju",
]
isBootstrapPeer    = false
keyBundleID        = "73e8966a78ca09bb912e9565cfb79fbe8a6048fab1f0cf49b18047c3895e0447"
monitoringEndpoint = "chain.link:4321"
transmitterAddress = "0xF67D0290337bca0847005C7ffD1BC75BA9AAE6e4"
observationTimeout = "10s"
observationSource = """
ds1          [type=bridge name=voter_turnout];
"""
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "max task duration must be < observation timeout")
			},
		},
		{
			name: "individual max task duration > observation timeout should error",
			toml: `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = [
"/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju",
]
isBootstrapPeer    = false
keyBundleID        = "73e8966a78ca09bb912e9565cfb79fbe8a6048fab1f0cf49b18047c3895e0447"
monitoringEndpoint = "chain.link:4321"
transmitterAddress = "0xF67D0290337bca0847005C7ffD1BC75BA9AAE6e4"
observationTimeout = "10s"
observationSource = """
ds1          [type=bridge name=voter_turnout timeout="30s"];
"""
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "individual max task duration must be < observation timeout")
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
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = [
"/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju",
]
isBootstrapPeer    = false
keyBundleID        = "73e8966a78ca09bb912e9565cfb79fbe8a6048fab1f0cf49b18047c3895e0447"
monitoringEndpoint = "chain.link:4321"
transmitterAddress = "0xF67D0290337bca0847005C7ffD1BC75BA9AAE6e4"
observationSource = """
ds1          [type=bridge name=voter_turnout];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
answer1      [type=median index=0];
"""
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "data source timeout must be between 1s and 20s, but is currently 20m0s")
			},
			setGlobalCfg: func(t *testing.T, c *configtest.TestGeneralConfig) {
				d := (20 * time.Minute)
				c.Overrides.OCRObservationTimeout = &d
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := configtest.NewTestGeneralConfig(t)
			c.Overrides.Dev = null.BoolFrom(false)
			c.Overrides.EVMRPCEnabled = null.BoolFrom(false)
			cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{GeneralConfig: c})
			if tc.setGlobalCfg != nil {
				tc.setGlobalCfg(t, c)
			}
			s, err := ValidatedOracleSpecToml(cc, tc.toml)
			tc.assertion(t, s, err)
		})
	}
}
