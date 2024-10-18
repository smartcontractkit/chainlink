package keystone_test

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/require"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/maps"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/memory"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeploy(t *testing.T) {
	lggr := logger.TestLogger(t)

	// sepolia; all nodes are on the this chain
	sepoliaChainId := uint64(11155111)
	sepoliaArbitrumChainId := uint64(421614)

	sepoliaChainSel, err := chainsel.SelectorFromChainId(sepoliaChainId)
	require.NoError(t, err)
	// sepoliaArbitrumChainSel, err := chainsel.SelectorFromChainId(sepoliaArbitrumChainId)
	// require.NoError(t, err)
	// aptosChainSel := uint64(999) // TODO:

	crConfig := deployment.CapabilityRegistryConfig{
		EVMChainID: sepoliaChainId,
		Contract:   [20]byte{},
	}

	evmChains := memory.NewMemoryChainsWithChainIDs(t, []uint64{sepoliaChainId, sepoliaArbitrumChainId})
	// aptosChain := memory.NewMemoryChain(t, aptosChainSel)

	// TODO: also need to tag these nodes
	wfChains := map[uint64]deployment.Chain{}
	wfChains[sepoliaChainSel] = evmChains[sepoliaChainSel]
	// wfChains[aptosChainSel] = aptosChain
	wfNodes := memory.NewNodes(t, zapcore.InfoLevel, wfChains, 4, 0, crConfig)
	require.Len(t, wfNodes, 4)

	cwNodes := memory.NewNodes(t, zapcore.InfoLevel, evmChains, 4, 0, crConfig)

	assetChains := map[uint64]deployment.Chain{}
	assetChains[sepoliaChainSel] = evmChains[sepoliaChainSel]
	assetNodes := memory.NewNodes(t, zapcore.InfoLevel, assetChains, 4, 0, crConfig)
	require.Len(t, assetNodes, 4)

	wfDon := keystone.DonCapabilities{
		Name:         keystone.WFDonName,
		Nodes:        maps.Keys(wfNodes),
		Capabilities: []kcr.CapabilitiesRegistryCapability{keystone.OCR3Cap},
	}
	cwDon := keystone.DonCapabilities{
		Name:         keystone.TargetDonName,
		Nodes:        maps.Keys(cwNodes),
		Capabilities: []kcr.CapabilitiesRegistryCapability{keystone.WriteChainCap},
	}
	assetDon := keystone.DonCapabilities{
		Name:         keystone.StreamDonName,
		Nodes:        maps.Keys(assetNodes),
		Capabilities: []kcr.CapabilitiesRegistryCapability{keystone.StreamTriggerCap},
	}

	allChains := make(map[uint64]deployment.Chain)
	maps.Copy(allChains, evmChains)
	// allChains[aptosChainSel] = aptosChain

	allNodes := make(map[string]memory.Node)
	maps.Copy(allNodes, wfNodes)
	maps.Copy(allNodes, cwNodes)
	maps.Copy(allNodes, assetNodes)
	env := memory.NewMemoryEnvironmentFromChainsNodes(t, lggr, allChains, allNodes)

	var ocr3Config = keystone.OracleConfigSource{
		MaxFaultyOracles: len(wfNodes) / 3,
	}

	ctx := tests.Context(t)
	// explicitly deploy the contracts
	cs, err := keystone.DeployContracts(lggr, &env, sepoliaChainSel)
	require.NoError(t, err)

	deployReq := keystone.ConfigureContractsRequest{
		RegistryChainSel: sepoliaChainSel,
		Env:              &env,
		OCR3Config:       &ocr3Config,
		Dons:             []keystone.DonCapabilities{wfDon, cwDon, assetDon},
		AddressBook:      cs.AddressBook,
		DoContractDeploy: false,
	}
	deployResp, err := keystone.ConfigureContracts(ctx, lggr, deployReq)
	require.NoError(t, err)
	ad := deployResp.Changeset.AddressBook
	addrs, err := ad.Addresses()
	require.NoError(t, err)
	lggr.Infow("Deployed Keystone contracts", "address book", addrs)

	// all contracts on home chain
	homeChainAddrs, err := ad.AddressesForChain(sepoliaChainSel)
	require.NoError(t, err)
	require.Len(t, homeChainAddrs, 3)
	// only forwarder on non-home chain
	for sel := range env.Chains {
		chainAddrs, err := ad.AddressesForChain(sel)
		require.NoError(t, err)
		if sel != sepoliaChainSel {
			require.Len(t, chainAddrs, 1)
		} else {
			require.Len(t, chainAddrs, 3)
		}
		containsForwarder := false
		for _, tv := range chainAddrs {
			if tv.Type == keystone.KeystoneForwarder {
				containsForwarder = true
				break
			}
		}
		require.True(t, containsForwarder, "no forwarder found in %v on chain %d for target don", chainAddrs, sel)
	}
	req := &keystone.GetContractSetsRequest{
		Chains:      env.Chains,
		AddressBook: ad,
	}

	contractSetsResp, err := keystone.GetContractSets(req)
	require.NoError(t, err)
	require.Len(t, contractSetsResp.ContractSets, len(env.Chains))
	// check the registry
	regChainContracts, ok := contractSetsResp.ContractSets[sepoliaChainSel]
	require.True(t, ok)
	gotRegistry := regChainContracts.CapabilitiesRegistry
	require.NotNil(t, gotRegistry)
	// contract reads
	gotDons, err := gotRegistry.GetDONs(&bind.CallOpts{})
	if err != nil {
		err = keystone.DecodeErr(kcr.CapabilitiesRegistryABI, err)
		require.Fail(t, fmt.Sprintf("failed to get Dons from registry at %s: %s", gotRegistry.Address().String(), err))
	}
	require.NoError(t, err)
	assert.Len(t, gotDons, len(deployReq.Dons))

	for n, info := range deployResp.DonInfos {
		found := false
		for _, gdon := range gotDons {
			if gdon.Id == info.Id {
				found = true
				assert.EqualValues(t, info, gdon)
				break
			}
		}
		require.True(t, found, "don %s not found in registry", n)
	}
	// check the forwarder
	for _, cs := range contractSetsResp.ContractSets {
		forwarder := cs.Forwarder
		require.NotNil(t, forwarder)
		// any read to ensure that the contract is deployed correctly
		_, err := forwarder.Owner(&bind.CallOpts{})
		require.NoError(t, err)
		// TODO expand this test; there is no get method on the forwarder so unclear how to test it
	}
	// check the ocr3 contract
	for chainSel, cs := range contractSetsResp.ContractSets {
		if chainSel != sepoliaChainSel {
			require.Nil(t, cs.OCR3)
			continue
		}
		require.NotNil(t, cs.OCR3)
		// any read to ensure that the contract is deployed correctly
		_, err := cs.OCR3.LatestConfigDetails(&bind.CallOpts{})
		require.NoError(t, err)
	}
}
