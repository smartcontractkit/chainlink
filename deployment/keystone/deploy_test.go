package keystone_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/require"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/maps"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/clo"
	"github.com/smartcontractkit/chainlink/deployment/environment/clo/models"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/deployment/keystone"
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

	// TODO: partition nodes into multiple nops

	wfDon := keystone.DonCapabilities{
		Name: keystone.WFDonName,
		Nops: []keystone.NOP{
			{
				Name:  "nop 1",
				Nodes: maps.Keys(wfNodes),
			},
		},
		Capabilities: []kcr.CapabilitiesRegistryCapability{keystone.OCR3Cap},
	}
	cwDon := keystone.DonCapabilities{
		Name: keystone.TargetDonName,
		Nops: []keystone.NOP{
			{
				Name:  "nop 2",
				Nodes: maps.Keys(cwNodes),
			},
		},
		Capabilities: []kcr.CapabilitiesRegistryCapability{keystone.WriteChainCap},
	}
	assetDon := keystone.DonCapabilities{
		Name: keystone.StreamDonName,
		Nops: []keystone.NOP{
			{
				Name:  "nop 3",
				Nodes: maps.Keys(assetNodes),
			},
		},
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

	var ocr3Config = keystone.OracleConfigWithSecrets{
		OracleConfig: keystone.OracleConfig{
			MaxFaultyOracles: len(wfNodes) / 3,
		},
		OCRSecrets: deployment.XXXGenerateTestOCRSecrets(),
	}

	ctx := tests.Context(t)
	// explicitly deploy the contracts
	cs, err := keystone.DeployContracts(lggr, &env, sepoliaChainSel)
	require.NoError(t, err)
	env.ExistingAddresses = cs.AddressBook
	deployReq := keystone.ConfigureContractsRequest{
		RegistryChainSel: sepoliaChainSel,
		Env:              &env,
		OCR3Config:       &ocr3Config,
		Dons:             []keystone.DonCapabilities{wfDon, cwDon, assetDon},
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

	contractSetsResp, err := keystone.GetContractSets(lggr, req)
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

// TODO: Deprecated, remove everything below that leverages CLO

func nodeOperatorsToIDs(t *testing.T, nops []*models.NodeOperator) (nodeIDs []keystone.NOP) {
	for _, nop := range nops {
		nodeOperator := keystone.NOP{
			Name: nop.Name,
		}
		for _, node := range nop.Nodes {
			p2pID, err := clo.NodeP2PId(node)
			require.NoError(t, err)

			nodeOperator.Nodes = append(nodeOperator.Nodes, p2pID)
		}
		nodeIDs = append(nodeIDs, nodeOperator)
	}
	return nodeIDs
}

func TestDeployCLO(t *testing.T) {
	lggr := logger.TestLogger(t)

	wfNops := loadTestNops(t, "testdata/workflow_nodes.json")
	cwNops := loadTestNops(t, "testdata/chain_writer_nodes.json")
	assetNops := loadTestNops(t, "testdata/asset_nodes.json")
	require.Len(t, wfNops, 10)
	requireChains(t, wfNops, []models.ChainType{models.ChainTypeEvm, models.ChainTypeAptos})
	require.Len(t, cwNops, 10)
	requireChains(t, cwNops, []models.ChainType{models.ChainTypeEvm, models.ChainTypeEvm})
	require.Len(t, assetNops, 16)
	requireChains(t, assetNops, []models.ChainType{models.ChainTypeEvm})

	wfNodes := nodeOperatorsToIDs(t, wfNops)
	cwNodes := nodeOperatorsToIDs(t, cwNops)
	assetNodes := nodeOperatorsToIDs(t, assetNops)

	wfDon := keystone.DonCapabilities{
		Name:         keystone.WFDonName,
		Nops:         wfNodes,
		Capabilities: []kcr.CapabilitiesRegistryCapability{keystone.OCR3Cap},
	}
	cwDon := keystone.DonCapabilities{
		Name:         keystone.TargetDonName,
		Nops:         cwNodes,
		Capabilities: []kcr.CapabilitiesRegistryCapability{keystone.WriteChainCap},
	}
	assetDon := keystone.DonCapabilities{
		Name:         keystone.StreamDonName,
		Nops:         assetNodes,
		Capabilities: []kcr.CapabilitiesRegistryCapability{keystone.StreamTriggerCap},
	}

	var allNops []*models.NodeOperator
	allNops = append(allNops, wfNops...)
	allNops = append(allNops, cwNops...)
	allNops = append(allNops, assetNops...)

	chains := make(map[uint64]struct{})
	for _, nop := range allNops {
		for _, node := range nop.Nodes {
			for _, chain := range node.ChainConfigs {
				// chain selector lib doesn't support chain id 2 and we don't use it in tests
				// because it's not an evm chain
				if chain.Network.ChainID == "2" { // aptos chain
					continue
				}
				id, err := strconv.ParseUint(chain.Network.ChainID, 10, 64)
				require.NoError(t, err, "failed to parse chain id to uint64")
				chains[id] = struct{}{}
			}
		}
	}
	var chainIDs []uint64
	for c := range chains {
		chainIDs = append(chainIDs, c)
	}
	allChains := memory.NewMemoryChainsWithChainIDs(t, chainIDs)

	env := &deployment.Environment{
		Name:              "CLO",
		ExistingAddresses: deployment.NewMemoryAddressBook(),
		Offchain:          clo.NewJobClient(lggr, clo.JobClientConfig{Nops: allNops}),
		Chains:            allChains,
		Logger:            lggr,
	}
	// assume that all the nodes in the provided input nops are part of the don
	for _, nop := range allNops {
		for _, node := range nop.Nodes {
			env.NodeIDs = append(env.NodeIDs, node.ID)
		}
	}

	// sepolia; all nodes are on the this chain
	registryChainSel, err := chainsel.SelectorFromChainId(11155111)
	require.NoError(t, err)

	var ocr3Config = keystone.OracleConfigWithSecrets{
		OracleConfig: keystone.OracleConfig{
			MaxFaultyOracles: len(wfNops) / 3,
		},
		OCRSecrets: deployment.XXXGenerateTestOCRSecrets(),
	}

	ctx := tests.Context(t)
	// explicitly deploy the contracts
	cs, err := keystone.DeployContracts(lggr, env, registryChainSel)
	require.NoError(t, err)
	// Deploy successful these are now part of our env.
	require.NoError(t, env.ExistingAddresses.Merge(cs.AddressBook))

	deployReq := keystone.ConfigureContractsRequest{
		RegistryChainSel: registryChainSel,
		Env:              env,
		OCR3Config:       &ocr3Config,
		Dons:             []keystone.DonCapabilities{wfDon, cwDon, assetDon},
		DoContractDeploy: false,
	}
	deployResp, err := keystone.ConfigureContracts(ctx, lggr, deployReq)
	require.NoError(t, err)
	ad := deployResp.Changeset.AddressBook
	addrs, err := ad.Addresses()
	require.NoError(t, err)
	lggr.Infow("Deployed Keystone contracts", "address book", addrs)

	// all contracts on home chain
	homeChainAddrs, err := ad.AddressesForChain(registryChainSel)
	require.NoError(t, err)
	require.Len(t, homeChainAddrs, 3)
	// only forwarder on non-home chain
	for sel := range env.Chains {
		chainAddrs, err := ad.AddressesForChain(sel)
		require.NoError(t, err)
		if sel != registryChainSel {
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

	contractSetsResp, err := keystone.GetContractSets(lggr, req)
	require.NoError(t, err)
	require.Len(t, contractSetsResp.ContractSets, len(env.Chains))
	// check the registry
	regChainContracts, ok := contractSetsResp.ContractSets[registryChainSel]
	require.True(t, ok)
	gotRegistry := regChainContracts.CapabilitiesRegistry
	require.NotNil(t, gotRegistry)
	// check DONs
	gotDons, err := gotRegistry.GetDONs(&bind.CallOpts{})
	if err != nil {
		err = keystone.DecodeErr(kcr.CapabilitiesRegistryABI, err)
		require.Fail(t, fmt.Sprintf("failed to get DONs from registry at %s: %s", gotRegistry.Address().String(), err))
	}
	require.NoError(t, err)
	assert.Len(t, gotDons, len(deployReq.Dons))
	// check NOPs
	nops, err := gotRegistry.GetNodeOperators(&bind.CallOpts{})
	if err != nil {
		err = keystone.DecodeErr(kcr.CapabilitiesRegistryABI, err)
		require.Fail(t, fmt.Sprintf("failed to get NOPs from registry at %s: %s", gotRegistry.Address().String(), err))
	}
	require.NoError(t, err)
	assert.Len(t, nops, 26) // 10 NOPs owning workflow & writer DONs + 16 NOPs owning Asset DON

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
		if chainSel != registryChainSel {
			require.Nil(t, cs.OCR3)
			continue
		}
		require.NotNil(t, cs.OCR3)
		// any read to ensure that the contract is deployed correctly
		_, err := cs.OCR3.LatestConfigDetails(&bind.CallOpts{})
		require.NoError(t, err)
	}
}

func requireChains(t *testing.T, donNops []*models.NodeOperator, cs []models.ChainType) {
	got := make(map[models.ChainType]struct{})
	want := make(map[models.ChainType]struct{})
	for _, c := range cs {
		want[c] = struct{}{}
	}
	for _, nop := range donNops {
		for _, node := range nop.Nodes {
			for _, cc := range node.ChainConfigs {
				got[cc.Network.ChainType] = struct{}{}
			}
		}
		require.EqualValues(t, want, got, "did not find all chains in node %s", nop.Name)
	}
}

func loadTestNops(t *testing.T, pth string) []*models.NodeOperator {
	f, err := os.ReadFile(pth)
	require.NoError(t, err)
	var nops []*models.NodeOperator
	require.NoError(t, json.Unmarshal(f, &nops))
	return nops
}
