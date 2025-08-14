package contracts_test

import (
	"fmt"
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/wsrpc/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations/optest"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
)

func Test_DeployForwardersSeq(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Chains: 2,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	registrySel := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]
	otherChainSel := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[1]
	b := optest.NewBundle(t)
	deps := contracts.DeployKeystoneForwardersSequenceDeps{
		Env: &env,
	}
	input := contracts.DeployKeystoneForwardersInput{
		Targets: []uint64{registrySel, otherChainSel},
	}

	got, err := operations.ExecuteSequence(b, contracts.DeployKeystoneForwardersSequence, deps, input)
	require.NoError(t, err)
	// Check that the output has the address
	addrRefs, err := got.Output.Addresses.Fetch()
	require.NoError(t, err)
	require.Len(t, addrRefs, len(input.Targets))
}

func Test_DeployRegistryOp(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Chains: 2,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	registrySel := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]
	b := optest.NewBundle(t)
	deps := contracts.DeployCapabilityRegistryOpDeps{
		Env: &env,
	}
	input := contracts.DeployCapabilityRegistryInput{
		ChainSelector: registrySel,
	}

	got, err := operations.ExecuteOperation(b, contracts.DeployCapabilityRegistryOp, deps, input)
	require.NoError(t, err)
	addrRefs, err := got.Output.Addresses.Fetch()
	require.NoError(t, err)
	require.Len(t, addrRefs, 1)

	fmt.Println(env.DataStore.Addresses())
}

func Test_ConfigureForwardersSeq(t *testing.T) {
	te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
		WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
		AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
		WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
		NumChains:       2,
	})

	var wfNodes []string
	for _, id := range te.GetP2PIDs("wfDon") {
		wfNodes = append(wfNodes, id.String())
	}
	deps2 := contracts.ConfigureForwardersSeqDeps{
		Env:      &te.Env,
		Registry: te.CapabilitiesRegistry(),
	}
	input2 := contracts.ConfigureForwardersSeqInput{
		RegistryChainSel: te.RegistrySelector,
		DONs: []contracts.ConfigureKeystoneDON{
			{
				Name:    "wfDon",
				NodeIDs: wfNodes,
			},
		},
	}
	b := optest.NewBundle(t)
	got2, err := operations.ExecuteSequence(b, contracts.ConfigureForwardersSeq, deps2, input2)
	require.NoError(t, err)
	fmt.Println(got2)
}

func Test_DeployConfigureForwardersSeq(t *testing.T) {
	te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
		WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
		AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
		WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
		NumChains:       2,
	})
	deps := contracts.DeployConfigureForwardersSeqDeps{
		Env:      &te.Env,
		Registry: te.CapabilitiesRegistry(),
	}
	var wfNodes []string
	for _, id := range te.GetP2PIDs("wfDon") {
		wfNodes = append(wfNodes, id.String())
	}
	registrySel := te.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]
	require.Equal(t, registrySel, te.RegistrySelector)
	input := contracts.DeployConfigureForwardersSeqInput{
		// ForwaderDeploymentChains: []uint64{registrySel},
		RegistryChainSel: te.RegistrySelector,
		Chain2WfDonMap: map[uint64][]contracts.ConfigureKeystoneDON{
			registrySel: {
				{
					Name:    "wfDon",
					NodeIDs: wfNodes,
				},
				{
					Name:    "wfDon2",
					NodeIDs: wfNodes,
				},
			},
		},
	}
	b := optest.NewBundle(t)
	_, err := operations.ExecuteSequence(b, contracts.DeployConfigureForwardersSeq, deps, input)
	require.NoError(t, err)
}
