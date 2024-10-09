package keystone_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/clo"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/clo/models"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeploy(t *testing.T) {
	lggr := logger.TestLogger(t)

	wfNops := loadTestNops(t, "../clo/testdata/workflow_nodes.json")
	cwNops := loadTestNops(t, "../clo/testdata/chain_writer_nodes.json")
	assetNops := loadTestNops(t, "../clo/testdata/asset_nodes.json")
	require.Len(t, wfNops, 10)
	require.Len(t, cwNops, 10)
	require.Len(t, assetNops, 16)

	wfDon := keystone.DonCapabilities{
		Name:         keystone.WFDonName,
		Nops:         wfNops,
		Capabilities: []kcr.CapabilitiesRegistryCapability{keystone.OCR3Cap},
	}
	cwDon := keystone.DonCapabilities{
		Name:         keystone.TargetDonName,
		Nops:         cwNops,
		Capabilities: []kcr.CapabilitiesRegistryCapability{keystone.WriteChainCap},
	}
	assetDon := keystone.DonCapabilities{
		Name:         keystone.StreamDonName,
		Nops:         assetNops,
		Capabilities: []kcr.CapabilitiesRegistryCapability{keystone.StreamTriggerCap},
	}

	env := makeMultiDonTestEnv(t, lggr, []keystone.DonCapabilities{wfDon, cwDon, assetDon})

	// sepolia; all nodes are on the this chain
	registryChainSel, err := chainsel.SelectorFromChainId(11155111)
	require.NoError(t, err)

	var ocr3Config = keystone.OracleConfigSource{
		MaxFaultyOracles: len(wfNops) / 3,
	}

	ctx := context.Background()

	// explicitly deploy the contracts
	cs, err := keystone.DeployContracts(lggr, env, registryChainSel)
	require.NoError(t, err)

	deployReq := keystone.ConfigureContractsRequest{
		RegistryChainSel: registryChainSel,
		Env:              env,
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

	contractSetsResp, err := keystone.GetContractSets(req)
	require.NoError(t, err)
	require.Len(t, contractSetsResp.ContractSets, len(env.Chains))
	// check the registry
	regChainContracts, ok := contractSetsResp.ContractSets[registryChainSel]
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

func makeMultiDonTestEnv(t *testing.T, lggr logger.Logger, dons []keystone.DonCapabilities) *deployment.Environment {
	var donToEnv = make(map[string]*deployment.Environment)
	for _, don := range dons {
		env := clo.NewDonEnvWithMemoryChains(t, clo.DonEnvConfig{
			DonName: don.Name,
			Nops:    don.Nops,
			Logger:  lggr,
		})
		donToEnv[don.Name] = env
	}
	menv := clo.NewTestEnv(t, lggr, donToEnv)
	return menv.Flatten("testing-env")
}

func loadTestNops(t *testing.T, pth string) []*models.NodeOperator {
	f, err := os.ReadFile(pth)
	require.NoError(t, err)
	var nops []*models.NodeOperator
	require.NoError(t, json.Unmarshal(f, &nops))
	return nops
}
