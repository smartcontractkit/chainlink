package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
)

const (
	testEVMCap    = "evm:ChainSelector:3379446385462418246@1.0.0"
	testAptosCap  = "aptos:ChainSelector:4457093679053095497@1.0.0"
	testStellar   = "stellar:ChainSelector:7777@1.0.0"
	testSolanaCap = "solana:ChainSelector:124615329519749607@1.0.0"
)

func chainSpec(t *testing.T, kv map[string]any) *capabilitiespb.CapabilityConfig {
	t.Helper()
	m, err := values.WrapMap(kv)
	require.NoError(t, err)
	return &capabilitiespb.CapabilityConfig{SpecConfig: values.ProtoMap(m)}
}

func fullSpec() map[string]any {
	return map[string]any{"chainId": "1337", "network": "evm", "creForwarderAddress": "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0"}
}

// payload marshals with EmitDefaultValues, so localOnly: false is stated.
func payload(t *testing.T, dons map[string]map[string]*capabilitiespb.CapabilityConfig) string {
	t.Helper()
	reg := &capabilitiespb.OffchainCapabilitiesRegistry{Domain: "cre", Env: "test", Version: 1, Dons: map[string]*capabilitiespb.DONConfig{}}
	for donName, caps := range dons {
		reg.Dons[donName] = &capabilitiespb.DONConfig{Capabilities: caps}
	}
	out, err := protojson.MarshalOptions{EmitDefaultValues: true}.Marshal(reg)
	require.NoError(t, err)
	return string(out)
}

func TestValidateChainCapabilityEntries_AcceptsChainEntries(t *testing.T) {
	t.Parallel()
	raw := payload(t, map[string]map[string]*capabilitiespb.CapabilityConfig{
		"workflow": {
			testEVMCap:   chainSpec(t, fullSpec()),
			testAptosCap: chainSpec(t, map[string]any{"chainId": "4", "network": "aptos", "creForwarderAddress": "0xa", "isLocal": false}),
			testStellar:  chainSpec(t, map[string]any{"chainId": "net", "network": "stellar", "creForwarderAddress": "C1"}),
			"cron@1.0.0": {}, // other capabilities are not checked
		},
	})
	require.NoError(t, ValidateChainCapabilityEntries(raw))
}

func TestValidateChainCapabilityEntries_LocalOnlyMustNotBeTrue(t *testing.T) {
	t.Parallel()

	// Omitted (what protojson.Marshal produces for false) and explicit false both match the onchain value.
	reg := &capabilitiespb.OffchainCapabilitiesRegistry{Version: 1, Dons: map[string]*capabilitiespb.DONConfig{
		"workflow": {Capabilities: map[string]*capabilitiespb.CapabilityConfig{testEVMCap: chainSpec(t, fullSpec())}},
	}}
	omitted, err := protojson.Marshal(reg)
	require.NoError(t, err)
	require.NoError(t, ValidateChainCapabilityEntries(string(omitted)))
	require.NoError(t, ValidateChainCapabilityEntries(payload(t, map[string]map[string]*capabilitiespb.CapabilityConfig{
		"workflow": {testEVMCap: chainSpec(t, fullSpec())},
	})))

	// true contradicts the onchain value.
	cfg := chainSpec(t, fullSpec())
	cfg.LocalOnly = true
	err = ValidateChainCapabilityEntries(payload(t, map[string]map[string]*capabilitiespb.CapabilityConfig{"workflow": {testEVMCap: cfg}}))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "localOnly must be false")
}

func TestValidateChainCapabilityEntries_RequiresSpecConfigKeys(t *testing.T) {
	t.Parallel()
	raw := payload(t, map[string]map[string]*capabilitiespb.CapabilityConfig{
		"workflow": {testEVMCap: chainSpec(t, map[string]any{"chainId": "1337"})},
	})
	err := ValidateChainCapabilityEntries(raw)
	require.Error(t, err)
	assert.Contains(t, err.Error(), `missing "network"`)
	assert.Contains(t, err.Error(), `missing "creForwarderAddress"`)
	assert.NotContains(t, err.Error(), `missing "chainId"`)
}

func TestValidateChainCapabilityEntries_RefusesSolanaAndMalformedIDs(t *testing.T) {
	t.Parallel()
	raw := payload(t, map[string]map[string]*capabilitiespb.CapabilityConfig{
		"workflow": {
			testSolanaCap:                          chainSpec(t, fullSpec()),
			"evm:ChainSelector:notanumber@1.0.0":   chainSpec(t, fullSpec()),
			"aptos:ChainSelector:4457093679053095": chainSpec(t, fullSpec()), // no version
		},
	})
	err := ValidateChainCapabilityEntries(raw)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "solana stays on job specs")
	assert.Contains(t, err.Error(), "malformed chain capability ID")
}

func TestValidateChainCapabilityEntries_RejectsInvalidPayload(t *testing.T) {
	t.Parallel()
	require.Error(t, ValidateChainCapabilityEntries(``))
	require.Error(t, ValidateChainCapabilityEntries(`{"version":1,"capabilityConfigs":{}}`), "unknown fields are rejected")
}
