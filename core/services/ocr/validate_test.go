package ocr_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr"
)

func TestValidateOracleSpec(t *testing.T) {
	var tt = []struct {
		name      string
		toml      string
		overrides func(c *chainlink.Config, s *chainlink.Secrets)
		assertion func(t *testing.T, os job.Job, err error)
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
p2pv2Bootstrappers = ["12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001"]
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
p2pv2Bootstrappers = ["12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001"]
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
p2pv2Bootstrappers = ["12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001"]
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
p2pv2Bootstrappers = ["12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001"]
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
p2pv2Bootstrappers = ["12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001"]
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
			name: "invalid v2 bootstrapper address",
			toml: `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pv2Bootstrappers = ["invalid bootstrapper /#@ address"]
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
p2pv2Bootstrappers = [
"12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001",
]
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
p2pv2Bootstrappers = [
"12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001",
]
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
p2pv2Bootstrappers = ["12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001"]
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
p2pv2Bootstrappers = ["12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001"]
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
p2pv2Bootstrappers = ["12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001"]
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
p2pv2Bootstrappers = ["12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001"]
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
p2pv2Bootstrappers = ["12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001"]
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
p2pv2Bootstrappers = ["12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001"]
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
p2pv2Bootstrappers = ["12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001"]
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
			overrides: func(c *chainlink.Config, s *chainlink.Secrets) {
				c.OCR.ObservationTimeout = commonconfig.MustNewDuration(20 * time.Minute)
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.Insecure.OCRDevelopmentMode = null.BoolFrom(false).Ptr()
				if tc.overrides != nil {
					tc.overrides(c, s)
				}
			})

			s, err := ocr.ValidatedOracleSpecTomlCfg(c, func(id *big.Int) (evmconfig.ChainScopedConfig, error) {
				return evmtest.NewChainScopedConfig(t, c), nil
			}, tc.toml)
			tc.assertion(t, s, err)
		})
	}
}
