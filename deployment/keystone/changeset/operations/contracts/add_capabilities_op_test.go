package contracts_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations/optest"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func doAppendCapabilitiesOp(t *testing.T, useMcms bool) {
	te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
		WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
		AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
		WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
		NumChains:       2,
		UseMCMS:         useMcms,
	})
	b := optest.NewBundle(t)

	// write capabilities
	capA := kcr.CapabilitiesRegistryCapability{
		LabelledName:   "write_test_chain1",
		Version:        "0.4.2",
		CapabilityType: uint8(3),
	}
	capB := kcr.CapabilitiesRegistryCapability{
		LabelledName:   "write_test_chain2",
		Version:        "3.16.0",
		CapabilityType: uint8(3),
	}
	caps := []kcr.CapabilitiesRegistryCapability{capA, capB}
	newCapabilities := make(map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability)
	for _, id := range te.GetP2PIDs("writerDon") {
		newCapabilities[id] = caps
	}
	deps := contracts.AppendCapabilitiesOpDeps{
		Env:               &te.Env,
		RegistryRef:       te.CapabilityRegistryAddressRef(),
		P2pToCapabilities: newCapabilities,
	}
	input := contracts.AppendCapabilitiesOpInput{
		RegistryChainSel: te.RegistrySelector,
	}
	if useMcms {
		input.MCMSConfig = &changeset.MCMSConfig{MinDuration: 0}
	}
	opOutput, err := operations.ExecuteOperation(b, contracts.AppendCapabilitiesOp, deps, input)
	require.NoError(t, err)
	if useMcms {
		require.Len(t, opOutput.Output.MCMSTimelockProposals, 1)
		require.Len(t, opOutput.Output.MCMSTimelockProposals[0].Operations, 1)
		require.Len(t, opOutput.Output.MCMSTimelockProposals[0].Operations[0].Transactions, 2) // add capabilities, update nodes
	}
}
func TestAppendCapabilitiesWithMCMS(t *testing.T) {
	t.Parallel()
	doAppendCapabilitiesOp(t, true)
}
func TestAppendCapabilitiesWithoutMCMS(t *testing.T) {
	t.Parallel()
	doAppendCapabilitiesOp(t, false)
}

func doUpdateDonOp(t *testing.T, useMcms bool) {
	te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
		WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
		AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
		WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
		NumChains:       2,
		UseMCMS:         useMcms,
	})
	b := optest.NewBundle(t)

	// write capabilities
	capA := kcr.CapabilitiesRegistryCapability{
		LabelledName:   "write_test_chain1",
		Version:        "0.4.2",
		CapabilityType: uint8(3),
	}
	capB := kcr.CapabilitiesRegistryCapability{
		LabelledName:   "write_test_chain2",
		Version:        "3.16.0",
		CapabilityType: uint8(3),
	}
	capACfg, err := proto.Marshal(test.GetDefaultCapConfig(t, capA))
	require.NoError(t, err)
	capBCfg, err := proto.Marshal(test.GetDefaultCapConfig(t, capB))
	require.NoError(t, err)

	deps := contracts.UpdateDonOpDeps{
		Env:         &te.Env,
		RegistryRef: te.CapabilityRegistryAddressRef(),
		P2PIDs:      te.GetP2PIDs("writerDon"),
		CapabilityConfigs: []changeset.CapabilityConfig{
			{
				Capability: capA, Config: capACfg,
			},
			{
				Capability: capB, Config: capBCfg,
			},
		},
	}
	input := contracts.UpdateDonOpInput{
		RegistryChainSel: te.RegistrySelector,
	}
	if useMcms {
		input.MCMSConfig = &changeset.MCMSConfig{MinDuration: 0}
	}

	opOutput, err := operations.ExecuteOperation(b, contracts.UpdateDonOp, deps, input)
	require.NoError(t, err)
	if useMcms {
		require.Len(t, opOutput.Output.MCMSTimelockProposals, 1)
		require.Len(t, opOutput.Output.MCMSTimelockProposals[0].Operations, 1)
		require.Len(t, opOutput.Output.MCMSTimelockProposals[0].Operations[0].Transactions, 3)
	}
}
func TestUpdateDonOpWithMCMS(t *testing.T) {
	t.Parallel()
	doUpdateDonOp(t, true)
}
func TestUpdateDonOpWithoutMCMS(t *testing.T) {
	t.Parallel()
	doUpdateDonOp(t, false)
}
