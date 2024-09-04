package keystone_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/test-go/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/memory"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var (
	adminAddr = "0x6CdfBF967A8ec4C29Fe26aF2a33Eb485d02f22D6"
	testNops  = []kcr.CapabilitiesRegistryNodeOperator{
		{
			Admin: common.HexToAddress("0x6CdfBF967A8ec4C29Fe26aF2a33Eb485d02f22D6"),
			Name:  "NOP_00",
		},
		{
			Admin: common.HexToAddress("0x6CdfBF967A8ec4C29Fe26aF2a33Eb485d02f2200"),
			Name:  "NOP_01",
		},
		{
			Admin: common.HexToAddress("0x11dfBF967A8ec4C29Fe26aF2a33Eb485d02f22D6"),
			Name:  "NOP_02",
		},
		{
			Admin: common.HexToAddress("0x6CdfBF967A8ec4C29Fe26aF2a33Eb485d02f2222"),
			Name:  "NOP_03",
		},
	}
)

func TestDeploy(t *testing.T) {
	lggr := logger.TestLogger(t)
	multDonCfg := memory.MemoryEnvironmentMultiDonConfig{
		Configs: make(map[string]memory.MemoryEnvironmentConfig),
	}
	wfEnvCfg := memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     1,
		Nodes:      4,
	}
	multDonCfg.Configs[keystone.WFDonName] = wfEnvCfg

	targetEnvCfg := memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     4,
		Nodes:      4,
	}
	multDonCfg.Configs[keystone.TargetDonName] = targetEnvCfg

	e := memory.NewMultiDonMemoryEnvironment(t, lggr, zapcore.InfoLevel, multDonCfg)

	var nodeToNop = make(map[string]kcr.CapabilitiesRegistryNodeOperator) //node -> nop
	// assign nops to nodes
	for _, env := range e.DonToEnv {
		for i, nodeID := range env.NodeIDs {
			idx := i % len(testNops)
			nop := testNops[idx]
			nodeToNop[nodeID] = nop
		}
	}

	var donsToDeploy = map[string][]kcr.CapabilitiesRegistryCapability{
		keystone.WFDonName:     []kcr.CapabilitiesRegistryCapability{keystone.OCR3Cap},
		keystone.TargetDonName: []kcr.CapabilitiesRegistryCapability{keystone.WriteChainCap},
	}

	ctx := context.Background()
	// Deploy all the Keystone contracts.
	homeChain := e.Get(keystone.WFDonName).AllChainSelectors()[0]
	deployReq := keystone.DeployRequest{
		RegistryChain:     homeChain,
		Menv:              e,
		DonToCapabilities: donsToDeploy,
		NodeIDToNop:       nodeToNop,
	}

	r, err := keystone.Deploy(ctx, lggr, deployReq)
	require.NoError(t, err)
	ad := r.Changeset.AddressBook
	addrs, err := ad.Addresses()
	require.NoError(t, err)
	lggr.Infow("Deployed Keystone contracts", "address book", addrs)

	// all contracts on home chain
	homeChainAddrs, err := ad.AddressesForChain(homeChain)
	require.NoError(t, err)
	require.Len(t, homeChainAddrs, 3)
	// only forwarder on non-home chain
	for _, chain := range e.Get(keystone.TargetDonName).AllChainSelectors() {
		chainAddrs, err := ad.AddressesForChain(chain)
		require.NoError(t, err)
		if chain != homeChain {
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
		require.True(t, containsForwarder, "no forwarder found in %v on chain %d for target don", chainAddrs, chain)
	}
	require.FailNow(t, "print logs")
}
