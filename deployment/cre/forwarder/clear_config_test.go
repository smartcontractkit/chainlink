package forwarder_test

import (
	"crypto/ecdsa"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldftesthelpers "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils/testhelpers"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations/optest"
	forwarderwrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"

	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/forwarder"
)

// TestClearConfigSeq exercises the sequence directly: configure a DON, then clear it and assert
// the forwarder emitted a ConfigSet with f == 0.
func TestClearConfigSeq(t *testing.T) {
	h, donConfig := setupForwarderTest(t, false)
	env := h.Runtime.Environment()
	chains := map[uint64]struct{}{h.RegistrySelector: {}}

	b := optest.NewBundle(t)

	// clearConfig only zeroes an existing entry, so configure first to have something to clear.
	_, err := operations.ExecuteSequence(b, forwarder.ConfigureSeq,
		forwarder.ConfigureSeqDeps{Env: &env},
		forwarder.ConfigureSeqInput{
			DON:       donConfig,
			Qualifier: setupForwarderQualifier,
			Chains:    chains,
		})
	require.NoError(t, err, "ConfigureSeq should execute successfully")
	requireLatestConfig(t, env, h.RegistrySelector, donConfig.ID, donConfig.Version, donConfig.F)

	out, err := operations.ExecuteSequence(b, forwarder.ClearConfigSeq,
		forwarder.ClearConfigSeqDeps{Env: &env},
		forwarder.ClearConfigSeqInput{
			DonID:         donConfig.ID,
			ConfigVersion: donConfig.Version,
			Qualifier:     setupForwarderQualifier,
			Chains:        chains,
		})
	require.NoError(t, err, "ClearConfigSeq should execute successfully")
	require.Empty(t, out.Output.MCMSTimelockProposals, "should not have MCMS proposals when not using MCMS")

	requireLatestConfig(t, env, h.RegistrySelector, donConfig.ID, donConfig.Version, 0)
}

// TestClearForwardersConfig covers the changeset wrapper and, by passing an empty Chains map,
// pins the documented "empty means all available chains" behaviour.
func TestClearForwardersConfig_AllChains(t *testing.T) {
	h, donConfig := setupForwarderTest(t, false)

	err := h.Runtime.Exec(runtime.ChangesetTask(forwarder.ConfigureForwarders{}, forwarder.ConfigureSeqInput{
		DON:       donConfig,
		Qualifier: setupForwarderQualifier,
		Chains:    map[uint64]struct{}{}, // Empty means all chains
	}))
	require.NoError(t, err, "configure changeset apply failed")

	task := runtime.ChangesetTask(forwarder.ClearForwardersConfig{}, forwarder.ClearConfigSeqInput{
		DonID:         donConfig.ID,
		ConfigVersion: donConfig.Version,
		Qualifier:     setupForwarderQualifier,
		Chains:        map[uint64]struct{}{}, // Empty means all chains
	})
	require.NoError(t, h.Runtime.Exec(task), "clear changeset apply failed")

	out := h.Runtime.State().Outputs[task.ID()]
	require.NotNil(t, out, "changeset output should not be nil")
	require.NotEmpty(t, out.Reports, "should have at least one report for the cleared chain")
	require.Empty(t, out.MCMSTimelockProposals, "should not have MCMS proposals when not using MCMS")

	env := h.Runtime.Environment()
	for chainSel := range env.BlockChains.EVMChains() {
		requireLatestConfig(t, env, chainSel, donConfig.ID, donConfig.Version, 0)
	}
}

// TestClearForwardersConfig_WithMCMS drives the full MCMS round trip: the config is set and then
// cleared through timelock proposals that the test signs and executes.
func TestClearForwardersConfig_WithMCMS(t *testing.T) {
	h, donConfig := setupForwarderTest(t, true)
	registryChainSel := h.RegistrySelector

	// MinDelay 0 so the scheduled operation is immediately executable without waiting.
	mcmsConfig := &contracts.MCMSConfig{
		MinDelay: 0,
		TimelockQualifierPerChain: map[uint64]string{
			registryChainSel: "",
		},
	}
	chains := map[uint64]struct{}{registryChainSel: {}}
	signers := []*ecdsa.PrivateKey{cldftesthelpers.TestXXXMCMSSigner}

	// The forwarder is MCMS-owned by now, so setting the config also has to go through a proposal.
	err := h.Runtime.Exec(
		runtime.ChangesetTask(forwarder.ConfigureForwarders{}, forwarder.ConfigureSeqInput{
			DON:        donConfig,
			Qualifier:  setupForwarderQualifier,
			MCMSConfig: mcmsConfig,
			Chains:     chains,
		}),
		runtime.SignAndExecuteProposalsTask(signers),
	)
	require.NoError(t, err, "configure changeset apply failed")
	requireLatestConfig(t, h.Runtime.Environment(), registryChainSel, donConfig.ID, donConfig.Version, donConfig.F)

	task := runtime.ChangesetTask(forwarder.ClearForwardersConfig{}, forwarder.ClearConfigSeqInput{
		DonID:         donConfig.ID,
		ConfigVersion: donConfig.Version,
		Qualifier:     setupForwarderQualifier,
		MCMSConfig:    mcmsConfig,
		Chains:        chains,
	})
	require.NoError(t, h.Runtime.Exec(task, runtime.SignAndExecuteProposalsTask(signers)), "clear changeset apply failed")

	out := h.Runtime.State().Outputs[task.ID()]
	require.NotNil(t, out, "changeset output should not be nil")
	require.NotEmpty(t, out.MCMSTimelockProposals, "should have MCMS proposals when using MCMS")

	requireLatestConfig(t, h.Runtime.Environment(), registryChainSel, donConfig.ID, donConfig.Version, 0)
}

// TestClearConfigSeq_Idempotent documents that clearing an unconfigured (donId, configVersion) is
// not an error: clearConfig zeroes the fault tolerance unconditionally.
func TestClearConfigSeq_Idempotent(t *testing.T) {
	h, donConfig := setupForwarderTest(t, false)
	env := h.Runtime.Environment()

	input := forwarder.ClearConfigSeqInput{
		DonID:         donConfig.ID,
		ConfigVersion: donConfig.Version,
		Qualifier:     setupForwarderQualifier,
		Chains:        map[uint64]struct{}{h.RegistrySelector: {}},
	}

	// Never configured, so there is nothing to clear.
	_, err := operations.ExecuteSequence(optest.NewBundle(t), forwarder.ClearConfigSeq,
		forwarder.ClearConfigSeqDeps{Env: &env}, input)
	require.NoError(t, err, "clearing an unconfigured DON should not fail")
	requireLatestConfig(t, env, h.RegistrySelector, donConfig.ID, donConfig.Version, 0)

	// Clearing twice is also fine.
	_, err = operations.ExecuteSequence(optest.NewBundle(t), forwarder.ClearConfigSeq,
		forwarder.ClearConfigSeqDeps{Env: &env}, input)
	require.NoError(t, err, "clearing twice should not fail")
	requireLatestConfig(t, env, h.RegistrySelector, donConfig.ID, donConfig.Version, 0)
}

func TestClearConfigSeq_NoForwarderDeployed(t *testing.T) {
	t.Parallel()

	env := newSimulatedEnv(t)
	selector := chainsel.TEST_90000001.Selector

	_, err := operations.ExecuteSequence(optest.NewBundle(t), forwarder.ClearConfigSeq,
		forwarder.ClearConfigSeqDeps{Env: env},
		forwarder.ClearConfigSeqInput{
			DonID:         1,
			ConfigVersion: 1,
			Qualifier:     "does-not-exist",
			Chains:        map[uint64]struct{}{selector: {}},
		})
	require.Error(t, err)
	require.Contains(t, err.Error(), "no KeystoneForwarder contract found")
}

func TestClearForwardersConfig_VerifyPreconditions(t *testing.T) {
	t.Parallel()

	env := newSimulatedEnv(t)
	selector := chainsel.TEST_90000001.Selector

	valid := forwarder.ClearConfigSeqInput{
		DonID:         1,
		ConfigVersion: 1,
		Qualifier:     setupForwarderQualifier,
		Chains:        map[uint64]struct{}{selector: {}},
	}

	tests := []struct {
		name    string
		mutate  func(*forwarder.ClearConfigSeqInput)
		wantErr string
	}{
		{
			name:   "valid input",
			mutate: func(*forwarder.ClearConfigSeqInput) {},
		},
		{
			name:   "no chains means all chains",
			mutate: func(in *forwarder.ClearConfigSeqInput) { in.Chains = nil },
		},
		{
			name:    "unknown chain selector",
			mutate:  func(in *forwarder.ClearConfigSeqInput) { in.Chains = map[uint64]struct{}{1: {}} },
			wantErr: "not found in environment",
		},
		{
			name:    "zero DON ID",
			mutate:  func(in *forwarder.ClearConfigSeqInput) { in.DonID = 0 },
			wantErr: "DON ID must be non-zero",
		},
		{
			name:    "zero config version",
			mutate:  func(in *forwarder.ClearConfigSeqInput) { in.ConfigVersion = 0 },
			wantErr: "config version must be non-zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			input := valid
			tt.mutate(&input)

			err := forwarder.ClearForwardersConfig{}.VerifyPreconditions(*env, input)
			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func newSimulatedEnv(t *testing.T) *cldf.Environment {
	t.Helper()

	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{chainsel.TEST_90000001.Selector}),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)

	return env
}

// requireLatestConfig asserts the fault tolerance carried by the most recent ConfigSet event for
// the given (donID, configVersion). The forwarder exposes no config getter, but both setConfig and
// clearConfig emit ConfigSet — clearConfig with f == 0 and no signers.
func requireLatestConfig(t *testing.T, env cldf.Environment, chainSel uint64, donID, configVersion uint32, wantF uint8) {
	t.Helper()

	refs := env.DataStore.Addresses().Filter(
		datastore.AddressRefByChainSelector(chainSel),
		datastore.AddressRefByType(datastore.ContractType(contracts.KeystoneForwarder)),
		datastore.AddressRefByQualifier(setupForwarderQualifier),
	)
	require.Len(t, refs, 1, "expected exactly one forwarder for chain %d", chainSel)

	chain := env.BlockChains.EVMChains()[chainSel]
	fwdr, err := forwarderwrapper.NewKeystoneForwarder(common.HexToAddress(refs[0].Address), chain.Client)
	require.NoError(t, err)

	iter, err := fwdr.FilterConfigSet(&bind.FilterOpts{Start: 0}, []uint32{donID}, []uint32{configVersion})
	require.NoError(t, err)
	defer iter.Close()

	var latest *forwarderwrapper.KeystoneForwarderConfigSet
	for iter.Next() {
		event := *iter.Event
		latest = &event
	}
	require.NoError(t, iter.Error())
	require.NotNil(t, latest, "expected a ConfigSet event for don %d version %d on chain %d", donID, configVersion, chainSel)

	require.Equal(t, wantF, latest.F, "unexpected fault tolerance on chain %d", chainSel)
	if wantF == 0 {
		require.Empty(t, latest.Signers, "cleared config should carry no signers")
	}
}
