package validate_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
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
				require.NoError(t, medianconfig.ValidatePluginConfig(pc))
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
			name: "valid DKG pluginConfig",
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
pluginType = "dkg"
transmitterID = "0x74103Cf8b436465870b26aa9Fa2F62AD62b22E35"

[relayConfig]
chainID = 4

[onchainSigningStrategy]
strategyName = "single-chain"
[onchainSigningStrategy.config]
evm = ""
publicKey = "0x1234567890123456789012345678901234567890"

[pluginConfig]
EncryptionPublicKey = "0e86e8cf899ae9a1b43e023bbe8825b103659bb8d6d4e54f6a3cfae7b106069c"
SigningPublicKey    = "eb62dbd2beb7c1524275a8019022f6ce6a7e86c9e65e3099452a2b96fc2432b1"
KeyID               = "6f3b82406688b8ddb944c6f2e6d808f014c8fa8d568d639c25019568c715fbf0"
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "DKG encryption key is not hex",
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
pluginType = "dkg"
transmitterID = "0x74103Cf8b436465870b26aa9Fa2F62AD62b22E35"

[relayConfig]
chainID = 4

[onchainSigningStrategy]
strategyName = "single-chain"
[onchainSigningStrategy.config]
evm = ""
publicKey = "0x1234567890123456789012345678901234567890"

[pluginConfig]
EncryptionPublicKey = "frog"
SigningPublicKey    = "eb62dbd2beb7c1524275a8019022f6ce6a7e86c9e65e3099452a2b96fc2432b1"
KeyID               = "6f3b82406688b8ddb944c6f2e6d808f014c8fa8d568d639c25019568c715fbf0"
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "expected hex string but received frog")
				require.Contains(t, err.Error(), "validation error for encryptedPublicKey")
			},
		},
		{
			name: "DKG encryption key is too short",
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
pluginType = "dkg"
transmitterID = "0x74103Cf8b436465870b26aa9Fa2F62AD62b22E35"

[relayConfig]
chainID = 4

[onchainSigningStrategy]
strategyName = "single-chain"
[onchainSigningStrategy.config]
evm = ""
publicKey = "0x1234567890123456789012345678901234567890"

[pluginConfig]
EncryptionPublicKey = "0e86e8cf899ae9a1b43e023bbe8825b103659bb8d6d4e54f6a3cfae7b10606"
SigningPublicKey    = "eb62dbd2beb7c1524275a8019022f6ce6a7e86c9e65e3099452a2b96fc2432b1"
KeyID               = "6f3b82406688b8ddb944c6f2e6d808f014c8fa8d568d639c25019568c715fbf0"
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "value: 0e86e8cf899ae9a1b43e023bbe8825b103659bb8d6d4e54f6a3cfae7b10606 has unexpected length. Expected 32 bytes")
				require.Contains(t, err.Error(), "validation error for encryptedPublicKey")
			},
		},
		{
			name: "DKG signing key is not hex",
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
pluginType = "dkg"
transmitterID = "0x74103Cf8b436465870b26aa9Fa2F62AD62b22E35"

[relayConfig]
chainID = 4

[onchainSigningStrategy]
strategyName = "single-chain"
[onchainSigningStrategy.config]
evm = ""
publicKey = "0x1234567890123456789012345678901234567890"

[pluginConfig]
EncryptionPublicKey = "0e86e8cf899ae9a1b43e023bbe8825b103659bb8d6d4e54f6a3cfae7b106069c"
SigningPublicKey    = "frog"
KeyID               = "6f3b82406688b8ddb944c6f2e6d808f014c8fa8d568d639c25019568c715fbf0"
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "expected hex string but received frog")
				require.Contains(t, err.Error(), "validation error for signingPublicKey")
			},
		},
		{
			name: "DKG signing key is too short",
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
pluginType = "dkg"
transmitterID = "0x74103Cf8b436465870b26aa9Fa2F62AD62b22E35"

[relayConfig]
chainID = 4

[onchainSigningStrategy]
strategyName = "single-chain"
[onchainSigningStrategy.config]
evm = ""
publicKey = "0x1234567890123456789012345678901234567890"

[pluginConfig]
EncryptionPublicKey = "0e86e8cf899ae9a1b43e023bbe8825b103659bb8d6d4e54f6a3cfae7b106069c"
SigningPublicKey    = "eb62dbd2beb7c1524275a8019022f6ce6a7e86c9e65e3099452a2b96fc24"
KeyID               = "6f3b82406688b8ddb944c6f2e6d808f014c8fa8d568d639c25019568c715fbf0"
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "value: eb62dbd2beb7c1524275a8019022f6ce6a7e86c9e65e3099452a2b96fc24 has unexpected length. Expected 32 bytes")
				require.Contains(t, err.Error(), "validation error for signingPublicKey")
			},
		},
		{
			name: "DKG keyID is not hex",
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
pluginType = "dkg"
transmitterID = "0x74103Cf8b436465870b26aa9Fa2F62AD62b22E35"

[relayConfig]
chainID = 4

[onchainSigningStrategy]
strategyName = "single-chain"
[onchainSigningStrategy.config]
evm = ""
publicKey = "0x1234567890123456789012345678901234567890"

[pluginConfig]
EncryptionPublicKey = "0e86e8cf899ae9a1b43e023bbe8825b103659bb8d6d4e54f6a3cfae7b106069c"
SigningPublicKey    = "eb62dbd2beb7c1524275a8019022f6ce6a7e86c9e65e3099452a2b96fc2432b1"
KeyID               = "frog"
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "expected hex string but received frog")
				require.Contains(t, err.Error(), "validation error for keyID")
			},
		},
		{
			name: "DKG keyID is too long",
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
pluginType = "dkg"
transmitterID = "0x74103Cf8b436465870b26aa9Fa2F62AD62b22E35"

[relayConfig]
chainID = 4

[onchainSigningStrategy]
strategyName = "single-chain"
[onchainSigningStrategy.config]
evm = ""
publicKey = "0x1234567890123456789012345678901234567890"

[pluginConfig]
EncryptionPublicKey = "0e86e8cf899ae9a1b43e023bbe8825b103659bb8d6d4e54f6a3cfae7b106069c"
SigningPublicKey    = "eb62dbd2beb7c1524275a8019022f6ce6a7e86c9e65e3099452a2b96fc2432b1"
KeyID               = "6f3b82406688b8ddb944c6f2e6d808f014c8fa8d568d639c25019568c715fbaaaabc"
`,
			assertion: func(t *testing.T, os job.Job, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "value: 6f3b82406688b8ddb944c6f2e6d808f014c8fa8d568d639c25019568c715fbaaaabc has unexpected length. Expected 32 bytes")
				require.Contains(t, err.Error(), "validation error for keyID")
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
				require.NoError(t, medianconfig.ValidatePluginConfig(pc))
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
