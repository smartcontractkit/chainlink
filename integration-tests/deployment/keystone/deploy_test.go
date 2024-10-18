package keystone_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/clo"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/clo/models"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestCLOdata(t *testing.T) {
	// hack to test cli
	var wantHash = "70900d840775d0b53b5bcd825ed52267619b71c3a31d7c877062059d7caed3ae"
	b, err := os.ReadFile("../clo/testdata/workflow_nodes.json")
	require.NoError(t, err)
	b1 := sha256.Sum256(b)
	got := hex.EncodeToString(b1[:])
	t.Log("sha256 of workflow_nodes.json", got)
	require.Equal(t, wantHash, got)
	wfNops := loadTestNops(t, "../clo/testdata/workflow_nodes.json")
	for _, nop := range wfNops {
		for _, node := range nop.Nodes {
			debug, err := json.MarshalIndent(node, "", "  ")
			require.NoError(t, err)
			var hasSepolia bool
			var hasAptos bool
			for _, c := range node.ChainConfigs {
				if c.Network.ChainID == "11155111" {
					hasSepolia = true
				}
				if c.Network.ChainID == "2" {
					hasAptos = true
				}
			}
			if !hasSepolia || !hasAptos {
				t.Logf("Node %s\n%s", node.Name, debug)
				t.Logf("expected sepolia: %t, aptos: %t", hasSepolia, hasAptos)
				t.Fail()
			}
		}
	}
}

func TestDeploy(t *testing.T) {
	t.Skip("TODO: KS-478 fix this test")
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

	ctx := tests.Context(t)
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
	// chain selector lib doesn't support chain id 2 and we don't use it in tests
	// because it's not an evm chain
	ignoreAptos := func(c *models.NodeChainConfig) bool {
		return c.Network.ChainID == "2" // aptos chain
	}
	for _, don := range dons {
		env := clo.NewDonEnvWithMemoryChains(t, clo.DonEnvConfig{
			DonName: don.Name,
			Nops:    don.Nops,
			Logger:  lggr,
		}, ignoreAptos)
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
