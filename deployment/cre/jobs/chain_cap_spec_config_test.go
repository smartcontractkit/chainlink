package jobs

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
)

const testForwarder = "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0"

func TestChainCapabilityID(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "evm:ChainSelector:3379446385462418246@1.0.0", ChainCapabilityID("evm", 3379446385462418246))
	assert.Equal(t, "stellar:ChainSelector:7777@1.0.0", ChainCapabilityID("stellar", 7777))
}

func TestEVMCapabilityEntry(t *testing.T) {
	t.Parallel()
	sel := chainselectors.ETHEREUM_TESTNET_SEPOLIA.Selector

	entry, err := EVMCapabilityEntry(EVMCapabilityEntryInput{
		ChainSelector:           sel,
		ForwarderAddress:        testForwarder,
		DeltaStage:              15 * time.Second,
		ForwarderLookbackBlocks: 200,
	})
	require.NoError(t, err)
	assert.Equal(t, ChainCapabilityID("evm", sel), entry.CapabilityID)
	assert.False(t, entry.Config.GetLocalOnly())

	fields := entry.Config.GetSpecConfig().GetFields()
	assert.Equal(t, int64(11155111), fields["chainId"].GetInt64Value(), "chainId from the selector, kept integral")
	assert.Equal(t, "evm", fields["network"].GetStringValue())
	assert.Equal(t, testForwarder, fields["creForwarderAddress"].GetStringValue())
	assert.Equal(t, int64(15*time.Second), fields["deltaStage"].GetInt64Value())
	assert.Equal(t, int64(200), fields["forwarderLookbackBlocks"].GetInt64Value())
	assert.Equal(t, int64(logTriggerPollInterval), fields["logTriggerPollInterval"].GetInt64Value(), "default applied")
	assert.Equal(t, int64(receiverGasMinimum), fields["receiverGasMinimum"].GetInt64Value(), "default applied")
	assert.Equal(t, int64(logTriggerSendChanBuf), fields["logTriggerSendChannelBufferSize"].GetInt64Value(), "default applied")
	assert.NotContains(t, fields, "nodeAddress", "per node: the node fills it at launch")
}

func TestEVMCapabilityEntry_OverridesAndNodeAddress(t *testing.T) {
	t.Parallel()
	entry, err := EVMCapabilityEntry(EVMCapabilityEntryInput{
		ChainSelector:           chainselectors.ETHEREUM_TESTNET_SEPOLIA.Selector,
		ForwarderAddress:        testForwarder,
		DeltaStage:              time.Second,
		ForwarderLookbackBlocks: 200,
		Override: OverrideDefaultCfg{
			LogTriggerPollInterval:  logTriggerPollInterval * 2,
			ForwarderLookbackBlocks: 50,
			NodeAddress:             "0xshould-not-be-shared",
		},
	})
	require.NoError(t, err)
	fields := entry.Config.GetSpecConfig().GetFields()
	assert.Equal(t, int64(logTriggerPollInterval*2), fields["logTriggerPollInterval"].GetInt64Value())
	assert.Equal(t, int64(50), fields["forwarderLookbackBlocks"].GetInt64Value(), "a non-zero override wins over the DON-wide value")
	assert.NotContains(t, fields, "nodeAddress")
}

// The payload entry and the job-spec changeset build their config from the same function, so the shared
// spec_config is the job-spec config minus the per-node fields.
func TestEVMCapabilityEntry_MatchesJobSpecConfig(t *testing.T) {
	t.Parallel()
	sel := chainselectors.ETHEREUM_TESTNET_SEPOLIA.Selector
	override := OverrideDefaultCfg{ReceiverGasMinimum: 900}

	entry, err := EVMCapabilityEntry(EVMCapabilityEntryInput{
		ChainSelector: sel, ForwarderAddress: testForwarder, DeltaStage: 5 * time.Second, ForwarderLookbackBlocks: 100, Override: override,
	})
	require.NoError(t, err)

	jobCfg := evmCapConfig(override, 11155111, testForwarder, 5*time.Second, 100)
	raw, err := json.Marshal(jobCfg)
	require.NoError(t, err)
	var want map[string]any
	require.NoError(t, json.Unmarshal(raw, &want))

	got := entry.Config.GetSpecConfig().GetFields()
	require.Len(t, got, len(want))
	for key := range want {
		assert.Contains(t, got, key)
	}
}

func TestEVMCapabilityEntry_RejectsBadInput(t *testing.T) {
	t.Parallel()
	_, err := EVMCapabilityEntry(EVMCapabilityEntryInput{ChainSelector: chainselectors.ETHEREUM_TESTNET_SEPOLIA.Selector})
	require.ErrorContains(t, err, "forwarder address is required")
	_, err = EVMCapabilityEntry(EVMCapabilityEntryInput{ChainSelector: 1, ForwarderAddress: testForwarder})
	require.ErrorContains(t, err, "chain ID from selector")
}

func TestEVMCapabilityEntry_BuildsAValidPayload(t *testing.T) {
	t.Parallel()
	evm, err := EVMCapabilityEntry(EVMCapabilityEntryInput{ChainSelector: chainselectors.ETHEREUM_TESTNET_SEPOLIA.Selector, ForwarderAddress: testForwarder, DeltaStage: time.Second})
	require.NoError(t, err)

	payload, err := pkg.BuildCapabilitiesRegistryPayload("cre", "staging", 3, map[string]map[string]*capabilitiespb.CapabilityConfig{
		"workflow_1": {
			evm.CapabilityID:     evm.Config,
			"cron-trigger@1.0.0": {},
		},
	})
	require.NoError(t, err)
	require.NoError(t, pkg.ValidateChainCapabilityEntries(payload))
}
