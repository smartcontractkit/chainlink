package contracts_test

import (
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"gotest.tools/v3/assert"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations/optest"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func doDeployConfigureForwardersSeq(t *testing.T, useMcms bool) {
	te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
		WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
		AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
		WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
		NumChains:       2,
		UseMCMS:         useMcms,
	})
	testChain := te.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[1]
	require.NotEqual(t, testChain, te.RegistrySelector)
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
	caps := []kcr.CapabilitiesRegistryCapability{capA, capB}
	newCapabilities := make(map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability)
	for _, id := range te.GetP2PIDs("writerDon") {
		newCapabilities[id] = caps
	}
	deps := contracts.DeployConfigureForwardersSeqDeps{
		Env:         &te.Env,
		Registry:    te.CapabilitiesRegistry(),
		RegistryRef: te.CapabilityRegistryAddressRef(),
		WriteCapabilityConfigs: []internal.CapabilityConfig{
			{
				Capability: capA, Config: capACfg,
			},
			{
				Capability: capB, Config: capBCfg,
			},
		},
		P2pToWriteCapabilities: newCapabilities,
	}
	var wfNodes []string
	for _, id := range te.GetP2PIDs("wfDon") {
		wfNodes = append(wfNodes, id.String())
	}
	registrySel := te.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]
	require.Equal(t, registrySel, te.RegistrySelector)
	input := contracts.DeployConfigureForwardersSeqInput{
		ForwarderDeploymentChains: map[uint64]contracts.ForwarderDeploymentOps{
			registrySel: {
				WfDons: []contracts.ConfigureKeystoneDON{
					{
						Name:    "wfDon",
						NodeIDs: wfNodes,
					},
				},
			},
			testChain: {
				WfDons: []contracts.ConfigureKeystoneDON{
					{
						Name:    "wfDon",
						NodeIDs: wfNodes,
					},
				},
			},
		},
		RegistryChainSel: te.RegistrySelector,
	}
	if useMcms {
		input.MCMSConfig = &changeset.MCMSConfig{MinDuration: 0}
	}
	b := optest.NewBundle(t)
	seqOutput, err := operations.ExecuteSequence(b, contracts.DeployConfigureForwardersSeq, deps, input)
	require.NoError(t, err)
	if useMcms {
		// configure forwarders, add capabilities, update don
		assert.Equal(t, len(seqOutput.Output.MCMSTimelockProposals), 3)
	}
}

func Test_DeployConfigureForwardersSeqWithoutMCMSSetup(t *testing.T) {
	t.Parallel()
	doDeployConfigureForwardersSeq(t, false)
}

func Test_DeployConfigureForwardersSeqWithMCMSSetup(t *testing.T) {
	t.Parallel()
	doDeployConfigureForwardersSeq(t, true)
}
