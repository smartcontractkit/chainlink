package validate_test

import (
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	medianconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/median/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
)

func TestValidateOracleSpec(t *testing.T) {
	var tt = []struct {
		name      string
		toml      string
		overrides func(c *chainlink.Config, s *chainlink.Secrets)
		assertion func(t *testing.T, os job.Job, err error)
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
[onchainSigningStrategy]
strategyName = "single-chain"
[onchainSigningStrategy.config]
evm = ""
publicKey = "0x1234567890123456789012345678901234567890"
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
				require.NoError(t, pc.ValidatePluginConfig())
				var oss validate.OCR2OnchainSigningStrategy
				require.NoError(t, json.Unmarshal(r.OnchainSigningStrategy.Bytes(), &oss))
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
p2pv2Bootstrappers = [
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
[onchainSigningStrategy]
strategyName = "single-chain"
[onchainSigningStrategy.config]
evm = ""
publicKey = "0x1234567890123456789012345678901234567890"
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
p2pv2Bootstrappers = [
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
[onchainSigningStrategy]
strategyName = "single-chain"
[onchainSigningStrategy.config]
evm = ""
publicKey = "0x1234567890123456789012345678901234567890"
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
p2pv2Bootstrappers = []
[relayConfig]
chainID = 1337
[onchainSigningStrategy]
strategyName = "single-chain"
[onchainSigningStrategy.config]
evm = ""
publicKey = "0x1234567890123456789012345678901234567890"
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
p2pv2Bootstrappers = []
observationSource = """
->
"""
[relayConfig]
chainID = 1337
[onchainSigningStrategy]
strategyName = "single-chain"
[onchainSigningStrategy.config]
evm = ""
publicKey = "0x1234567890123456789012345678901234567890"
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
p2pv2Bootstrappers = ["/invalid/peer/address"]
observationSource = """
blah
"""
[relayConfig]
chainID = 1337
[onchainSigningStrategy]
strategyName = "single-chain"
[onchainSigningStrategy.config]
evm = ""
publicKey = "0x1234567890123456789012345678901234567890"
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
p2pv2Bootstrappers = ["12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001"]
blockchainTimeout  = "0s"
observationSource = """
blah
"""
[relayConfig]
chainID = 1337
[onchainSigningStrategy]
strategyName = "single-chain"
[onchainSigningStrategy.config]
evm = ""
publicKey = "0x1234567890123456789012345678901234567890"
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
p2pv2Bootstrappers = ["12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001"]
observationSource = """
blah
"""
[relayConfig]
chainID = 1337
[onchainSigningStrategy]
strategyName = "single-chain"
[onchainSigningStrategy.config]
evm = ""
publicKey = "0x1234567890123456789012345678901234567890"
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
p2pv2Bootstrappers = []
monitoringEndpoint = "\t/fd\2ff )(*&^%$#@"
[relayConfig]
chainID = 1337
[onchainSigningStrategy]
strategyName = "single-chain"
[onchainSigningStrategy.config]
evm = ""
publicKey = "0x1234567890123456789012345678901234567890"
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
p2pv2Bootstrappers = [
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
[onchainSigningStrategy]
strategyName = "single-chain"
[onchainSigningStrategy.config]
evm = ""
publicKey = "0x1234567890123456789012345678901234567890"
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
			overrides: func(c *chainlink.Config, s *chainlink.Secrets) {
				c.OCR2.DatabaseTimeout = commonconfig.MustNewDuration(20 * time.Minute)
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
[onchainSigningStrategy]
strategyName = "single-chain"
[onchainSigningStrategy.config]
evm = ""
publicKey = "0x1234567890123456789012345678901234567890"
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
[onchainSigningStrategy]
strategyName = "single-chain"
[onchainSigningStrategy.config]
evm = ""
publicKey = "0x1234567890123456789012345678901234567890"
[pluginConfig]
juelsPerFeeCoinSource = """
ds1          [type=bridge name=voter_turnout];
"""
[relayConfig]
chainID = 1337
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				t.Log("relay", os.OCR2OracleSpec.Relay)
				require.Error(t, err)
				require.Contains(t, err.Error(), "no such relay blerg supported")
			},
		},
		{
			name: "Generic public onchain signing strategy with no public key",
			toml: `
type               = "offchainreporting2"
pluginType         = "plugin"
schemaVersion      = 1
relay              = "evm"
contractID         = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pv2Bootstrappers = [
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
[onchainSigningStrategy]
strategyName = "single-chain"
[onchainSigningStrategy.config]
evm = ""
publicKey = ""
[pluginConfig]
pluginName = "median"
telemetryType = "median"
OCRVersion=2
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "must provide public key for the onchain signing strategy")
			},
		},
		{
			name: "Valid ccip-execute",
			toml: `
type = "offchainreporting2"
schemaVersion = 1
relay = "evm"
contractID = "0x1234567"
pluginType = "ccip-execution"

[relayConfig]
chainID = 1337

[pluginConfig]
SourceStartBlock = 1
DestStartBlock = 2
USDCConfig.SourceTokenAddress = "0x1234567890123456789012345678901234567890"
USDCConfig.SourceMessageTransmitterAddress = "0x0987654321098765432109876543210987654321"
USDCConfig.AttestationAPI = "some api"
USDCConfig.AttestationAPITimeoutSeconds = 12
USDCConfig.AttestationAPIIntervalMilliseconds = 100
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.NoError(t, err)
				expected := config.ExecPluginJobSpecConfig{
					SourceStartBlock: 1,
					DestStartBlock:   2,
					USDCConfig: config.USDCConfig{
						SourceTokenAddress:                 common.HexToAddress("0x1234567890123456789012345678901234567890"),
						SourceMessageTransmitterAddress:    common.HexToAddress("0x0987654321098765432109876543210987654321"),
						AttestationAPI:                     "some api",
						AttestationAPITimeoutSeconds:       12,
						AttestationAPIIntervalMilliseconds: 100,
					},
				}
				var cfg config.ExecPluginJobSpecConfig
				err = json.Unmarshal(os.OCR2OracleSpec.PluginConfig.Bytes(), &cfg)
				require.NoError(t, err)
				require.Equal(t, expected, cfg)
			},
		},
		{
			name: "ccip-execute non hex address unmarshalling",
			toml: `
type = "offchainreporting2"
schemaVersion = 1
relay = "evm"
contractID = "0x1234567"
pluginType = "ccip-execution"

[relayConfig]
chainID = 1337

[pluginConfig]
SourceStartBlock = 1
DestStartBlock = 2
USDCConfig.SourceTokenAddress = "non-hex"
USDCConfig.SourceMessageTransmitterAddress = "0x0987654321098765432109876543210987654321"
USDCConfig.AttestationAPI = "some api"
USDCConfig.AttestationAPITimeoutSeconds = 12
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "cannot unmarshal hex string without 0x prefix into Go struct field USDCConfig.USDCConfig.SourceTokenAddress of type common.Address")
			},
		},
		{
			name: "ccip-execute usdcconfig validation failure",
			toml: `
type = "offchainreporting2"
schemaVersion = 1
relay = "evm"
contractID = "0x1234567"
pluginType = "ccip-execution"

[relayConfig]
chainID = 1337

[pluginConfig]
SourceStartBlock = 1
DestStartBlock = 2
USDCConfig.SourceTokenAddress = "0x1234567890123456789012345678901234567890"
USDCConfig.SourceMessageTransmitterAddress = "0x0987654321098765432109876543210987654321"
USDCConfig.AttestationAPI = "some api"
USDCConfig.AttestationAPIIntervalMilliseconds = 100
USDCConfig.AttestationAPITimeoutSeconds = -12
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "error while unmarshalling plugin config: json: cannot unmarshal number -12 into Go struct field USDCConfig.USDCConfig.AttestationAPITimeoutSeconds of type uint")
			},
		},
		{
			name: "Valid ccip-commit pipeline",
			toml: `
type = "offchainreporting2"
schemaVersion = 1
relay = "evm"
contractID = "0x1234567"
pluginType = "ccip-commit"

[relayConfig]
chainID = 1337

[pluginConfig]
SourceStartBlock = 1
DestStartBlock = 2
offRamp = "0x1234567890123456789012345678901234567890"
tokenPricesUSDPipeline = "merge [type=merge left=\"{}\" right=\"{\\\"0xC79b96044906550A5652BCf20a6EA02f139B9Ae5\\\":\\\"1000000000000000000\\\"}\"];"
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.NoError(t, err)
				expected := config.CommitPluginJobSpecConfig{
					SourceStartBlock:       1,
					DestStartBlock:         2,
					OffRamp:                cciptypes.Address(common.HexToAddress("0x1234567890123456789012345678901234567890").String()),
					TokenPricesUSDPipeline: `merge [type=merge left="{}" right="{\"0xC79b96044906550A5652BCf20a6EA02f139B9Ae5\":\"1000000000000000000\"}"];`,
					PriceGetterConfig:      nil,
				}
				var cfg config.CommitPluginJobSpecConfig
				err = json.Unmarshal(os.OCR2OracleSpec.PluginConfig.Bytes(), &cfg)
				require.NoError(t, err)
				require.Equal(t, expected, cfg)
			},
		},
		{
			name: "Valid ccip-commit dynamic price getter",
			toml: `
type = "offchainreporting2"
schemaVersion = 1
relay = "evm"
contractID = "0x1234567"
pluginType = "ccip-commit"

[relayConfig]
chainID = 1337

[pluginConfig]
SourceStartBlock = 1
DestStartBlock = 2
offRamp = "0x1234567890123456789012345678901234567890"
priceGetterConfig = """
{
	"aggregatorPrices": {
		"0x0820c05e1fba1244763a494a52272170c321cad3": {
			"chainID": "1000",
			"contractAddress": "0xb8dabd288955d302d05ca6b011bb46dfa3ea7acf"
		},
		"0x4a98bb4d65347016a7ab6f85bea24b129c9a1272": {
			"chainID": "1337",
			"contractAddress": "0xb80244cc8b0bb18db071c150b36e9bcb8310b236"
		}
	},
	"staticPrices": {
		"0xec8c353470ccaa4f43067fcde40558e084a12927": {
			"chainID": "1057",
			"price": 1000000000000000000
		}
	}
}
"""
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.NoError(t, err)
				expected := config.CommitPluginJobSpecConfig{
					SourceStartBlock:       1,
					DestStartBlock:         2,
					OffRamp:                cciptypes.Address(common.HexToAddress("0x1234567890123456789012345678901234567890").String()),
					TokenPricesUSDPipeline: "",
					PriceGetterConfig: &config.DynamicPriceGetterConfig{
						AggregatorPrices: map[common.Address]config.AggregatorPriceConfig{
							common.HexToAddress("0x0820c05e1fba1244763a494a52272170c321cad3"): {
								ChainID:                   1000,
								AggregatorContractAddress: common.HexToAddress("0xb8dabd288955d302d05ca6b011bb46dfa3ea7acf"),
							},
							common.HexToAddress("0x4a98bb4d65347016a7ab6f85bea24b129c9a1272"): {
								ChainID:                   1337,
								AggregatorContractAddress: common.HexToAddress("0xb80244cc8b0bb18db071c150b36e9bcb8310b236"),
							},
						},
						StaticPrices: map[common.Address]config.StaticPriceConfig{
							common.HexToAddress("0xec8c353470ccaa4f43067fcde40558e084a12927"): {
								ChainID: 1057,
								Price:   big.NewInt(1000000000000000000),
							},
						},
					},
				}
				var cfg config.CommitPluginJobSpecConfig
				err = json.Unmarshal(os.OCR2OracleSpec.PluginConfig.Bytes(), &cfg)
				require.NoError(t, err)
				require.Equal(t, expected, cfg)
			},
		},
		{
			name: "ccip-commit dual price getter configuration",
			toml: `
type = "offchainreporting2"
schemaVersion = 1
relay = "evm"
contractID = "0x1234567"
pluginType = "ccip-commit"

[relayConfig]
chainID = 1337

[pluginConfig]
SourceStartBlock = 1
DestStartBlock = 2
offRamp = "0x1234567890123456789012345678901234567890"
tokenPricesUSDPipeline = "merge [type=merge left=\"{}\" right=\"{\\\"0xC79b96044906550A5652BCf20a6EA02f139B9Ae5\\\":\\\"1000000000000000000\\\"}\"];"
priceGetterConfig = """
{
	"aggregatorPrices": {
		"0x0820c05e1fba1244763a494a52272170c321cad3": {
			"chainID": "1000",
			"contractAddress": "0xb8dabd288955d302d05ca6b011bb46dfa3ea7acf"
		},
		"0x4a98bb4d65347016a7ab6f85bea24b129c9a1272": {
			"chainID": "1337",
			"contractAddress": "0xb80244cc8b0bb18db071c150b36e9bcb8310b236"
		}
	},
	"staticPrices": {
		"0xec8c353470ccaa4f43067fcde40558e084a12927": {
			"chainID": "1057",
			"price": 1000000000000000000
		}
	}
}
"""
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "only one of tokenPricesUSDPipeline or priceGetterConfig must be set")
			},
		},
		{
			name: "ccip-commit invalid pipeline",
			toml: `
type = "offchainreporting2"
schemaVersion = 1
relay = "evm"
contractID = "0x1234567"
pluginType = "ccip-commit"

[relayConfig]
chainID = 1337

[pluginConfig]
SourceStartBlock = 1
DestStartBlock = 2
offRamp = "0x1234567890123456789012345678901234567890"
tokenPricesUSDPipeline = "this is not a pipeline"
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "invalid token prices pipeline")
			},
		},
		{
			name: "ccip-commit invalid dynamic token prices config",
			toml: `
type = "offchainreporting2"
schemaVersion = 1
relay = "evm"
contractID = "0x1234567"
pluginType = "ccip-commit"

[relayConfig]
chainID = 1337

[pluginConfig]
SourceStartBlock = 1
DestStartBlock = 2
offRamp = "0x1234567890123456789012345678901234567890"
priceGetterConfig = "this is not a proper dynamic price config"
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "error while unmarshalling plugin config")
			},
		},
		{
			name: "Generic plugin config validation - nothing provided",
			toml: `
type = "offchainreporting2"
schemaVersion = 1
name = "dkg"
externalJobID = "6d46d85f-d38c-4f4a-9f00-ac29a25b6330"
maxTaskDuration = "1s"
contractID = "0x3e54dCc49F16411A3aaa4cDbC41A25bCa9763Cee"
ocrKeyBundleID = "08d14c6eed757414d72055d28de6caf06535806c6a14e450f3a2f1c854420e17"
p2pv2Bootstrappers = [
	"12D3KooWSbPRwXY4gxFRJT7LWCnjgGbR4S839nfCRCDgQUiNenxa@127.0.0.1:8000"
]
relay = "evm"
pluginType = "plugin"
transmitterID = "0x74103Cf8b436465870b26aa9Fa2F62AD62b22E35"

[relayConfig]
chainID = 4

[onchainSigningStrategy]
strategyName = "single-chain"
[onchainSigningStrategy.config]
evm = ""
publicKey = "0x1234567890123456789012345678901234567890"

[pluginConfig]
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "must provide plugin name")
			},
		}, {
			name: "Generic plugin config validation - ocr version",
			toml: `
type = "offchainreporting2"
schemaVersion = 1
name = "dkg"
externalJobID = "6d46d85f-d38c-4f4a-9f00-ac29a25b6330"
maxTaskDuration = "1s"
contractID = "0x3e54dCc49F16411A3aaa4cDbC41A25bCa9763Cee"
ocrKeyBundleID = "08d14c6eed757414d72055d28de6caf06535806c6a14e450f3a2f1c854420e17"
p2pv2Bootstrappers = [
	"12D3KooWSbPRwXY4gxFRJT7LWCnjgGbR4S839nfCRCDgQUiNenxa@127.0.0.1:8000"
]
relay = "evm"
pluginType = "plugin"
transmitterID = "0x74103Cf8b436465870b26aa9Fa2F62AD62b22E35"

[relayConfig]
chainID = 4

[onchainSigningStrategy]
strategyName = "single-chain"
[onchainSigningStrategy.config]
evm = ""
publicKey = "0x1234567890123456789012345678901234567890"

[pluginConfig]
PluginName="some random name"
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "only OCR version 2 and 3 are supported")
			},
		},
		{
			name: "Generic plugin config validation - no command",
			toml: `
type = "offchainreporting2"
schemaVersion = 1
name = "dkg"
externalJobID = "6d46d85f-d38c-4f4a-9f00-ac29a25b6330"
maxTaskDuration = "1s"
contractID = "0x3e54dCc49F16411A3aaa4cDbC41A25bCa9763Cee"
ocrKeyBundleID = "08d14c6eed757414d72055d28de6caf06535806c6a14e450f3a2f1c854420e17"
p2pv2Bootstrappers = [
	"12D3KooWSbPRwXY4gxFRJT7LWCnjgGbR4S839nfCRCDgQUiNenxa@127.0.0.1:8000"
]
relay = "evm"
pluginType = "plugin"
transmitterID = "0x74103Cf8b436465870b26aa9Fa2F62AD62b22E35"

[relayConfig]
chainID = 4

[onchainSigningStrategy]
strategyName = "single-chain"
[onchainSigningStrategy.config]
evm = ""
publicKey = "0x1234567890123456789012345678901234567890"

[pluginConfig]
PluginName="some random name"
OCRVersion=2
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "no command found")
			},
		},
		{
			name: "Generic plugin config validation - no binary",
			toml: `
type = "offchainreporting2"
schemaVersion = 1
name = "dkg"
externalJobID = "6d46d85f-d38c-4f4a-9f00-ac29a25b6330"
maxTaskDuration = "1s"
contractID = "0x3e54dCc49F16411A3aaa4cDbC41A25bCa9763Cee"
ocrKeyBundleID = "08d14c6eed757414d72055d28de6caf06535806c6a14e450f3a2f1c854420e17"
p2pv2Bootstrappers = [
	"12D3KooWSbPRwXY4gxFRJT7LWCnjgGbR4S839nfCRCDgQUiNenxa@127.0.0.1:8000"
]
relay = "evm"
pluginType = "plugin"
transmitterID = "0x74103Cf8b436465870b26aa9Fa2F62AD62b22E35"

[relayConfig]
chainID = 4

[onchainSigningStrategy]
strategyName = "single-chain"
[onchainSigningStrategy.config]
evm = ""
publicKey = "0x1234567890123456789012345678901234567890"

[pluginConfig]
PluginName="some random name"
OCRVersion=2
Command="some random command"
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "failed to find binary")
			},
		}, {
			name: "minimal OCR2 oracle spec with JuelsPerFeeCoinCache",
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
[pluginConfig.JuelsPerFeeCoinCache]
Disable=false
UpdateInterval="1m"
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.NoError(t, err)
				b, err := jsonapi.Marshal(os.OCR2OracleSpec)
				require.NoError(t, err)
				var r job.OCR2OracleSpec
				err = jsonapi.Unmarshal(b, &r)
				require.NoError(t, err)
				assert.Equal(t, "median", string(r.PluginType))
				var pc medianconfig.PluginConfig
				require.NoError(t, json.Unmarshal(r.PluginConfig.Bytes(), &pc))
				require.NoError(t, pc.ValidatePluginConfig())
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.Insecure.OCRDevelopmentMode = testutils.Ptr(false) // tests run with OCRDevelopmentMode by default.
				if tc.overrides != nil {
					tc.overrides(c, s)
				}
			})
			s, err := validate.ValidatedOracleSpecToml(testutils.Context(t), c.OCR2(), c.Insecure(), tc.toml, nil)
			tc.assertion(t, s, err)
		})
	}
}

type envelope struct {
	PluginConfig *validate.OCR2GenericPluginConfig
}

func TestOCR2GenericPluginConfig_Unmarshal(t *testing.T) {
	payload := `
[pluginConfig]
pluginName = "median"
telemetryType = "median"
foo = "bar"

[[pluginConfig.pipelines]]
name = "default"
spec = "a spec"
`
	tree, err := toml.Load(payload)
	require.NoError(t, err)

	// Load the toml how we load it in the plugin, i.e. convert to
	// map[string]any first, then treat as JSON
	o := map[string]any{}
	err = tree.Unmarshal(&o)
	require.NoError(t, err)

	b, err := json.Marshal(o)
	require.NoError(t, err)

	e := &envelope{}
	err = json.Unmarshal(b, e)
	require.NoError(t, err)

	pc := e.PluginConfig
	assert.Equal(t, "bar", pc.PluginConfig["foo"])
	assert.Len(t, pc.Pipelines, 1)
	assert.Equal(t, validate.PipelineSpec{Name: "default", Spec: "a spec"}, pc.Pipelines[0])
	assert.Equal(t, "median", pc.PluginName)
	assert.Equal(t, "median", pc.TelemetryType)
}

type envelope2 struct {
	OnchainSigningStrategy *validate.OCR2OnchainSigningStrategy
}

func TestOCR2OnchainSigningStrategy_Unmarshal(t *testing.T) {
	payload := `
[onchainSigningStrategy]
strategyName = "single-chain"
[onchainSigningStrategy.config]
evm = "08d14c6eed757414d72055d28de6caf06535806c6a14e450f3a2f1c854420e17"
publicKey = "0x1234567890123456789012345678901234567890"
`
	oss := &envelope2{}
	tree, err := toml.Load(payload)
	require.NoError(t, err)
	o := map[string]any{}
	err = tree.Unmarshal(&o)
	require.NoError(t, err)
	b, err := json.Marshal(o)
	require.NoError(t, err)
	err = json.Unmarshal(b, oss)
	require.NoError(t, err)

	pk, err := oss.OnchainSigningStrategy.PublicKey()
	require.NoError(t, err)
	kbID, err := oss.OnchainSigningStrategy.KeyBundleID("evm")
	require.NoError(t, err)

	assert.False(t, oss.OnchainSigningStrategy.IsMultiChain())
	assert.Equal(t, "0x1234567890123456789012345678901234567890", pk)
	assert.Equal(t, "08d14c6eed757414d72055d28de6caf06535806c6a14e450f3a2f1c854420e17", kbID)
}
