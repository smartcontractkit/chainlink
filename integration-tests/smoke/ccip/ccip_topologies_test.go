package ccip

import (
	"context"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_home"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/cciptesthelpertypes"
	mt "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

func Test_CCIPTopologies_EVM2EVM_RoleDON_AllSupportSource_SomeSupportDest(t *testing.T) {
	// Setup 3 chains and a single lane between chains that doesn't
	// include the home chain.
	// We need at least 7 nodes to set up a suitable role DON.
	// In this scenario, 7 nodes will support the source chain, and
	// only 4 nodes will support the destination chain. At least one
	// node will end up supporting both source and dest in this particular
	// scenario.
	// This will help us test the scenario where not all nodes support the destination chain.
	// There are two parts to the setup:
	// 1. the chainlink nodes themselves must be configured to support only a subset of the chains
	// that are spun up in the environment.
	// 2. the fChain values must be set correctly on CCIPHome according to the determined topology.
	// When spinning up an environment we don't know the chain IDs / selectors of the chains being
	// spun up, so we can't have an explicit mapping of chain IDs to fChain values / node indices.
	ctx := t.Context()

	const (
		// fRoleDON needs to be at least 2, since for a role don of size 4, we can't have
		// different node/chain committees.
		fRoleDON = 2
		nRoleDON = 3*fRoleDON + 1

		// we need at least 3 chains to test the role DON topology.
		// this is because one chain will be the home chain, which all nodes
		// must support, and the other two chains will be the source and dest
		// chains, of which only some nodes will support the dest chain.
		numChains = 3

		fChainSource = 2
		fChainDest   = 1
	)

	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(numChains),
		testhelpers.WithNumOfNodes(nRoleDON),
		testhelpers.WithRoleDONTopology(cciptesthelpertypes.NewRandomTopology(
			cciptesthelpertypes.RandomTopologyArgs{
				FChainToNumChains: map[int]int{
					fChainSource: 1, // 1 chain with fChain fChainSource
					fChainDest:   1, // 1 chain with fChain fChainDest
				},
				Seed: 42, // for reproducible setups.
			},
		)),
	)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	allChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	require.Len(t, allChainSelectors, 3)
	// filter out the home chain
	var nonHomeChains []uint64
	for _, chainSel := range allChainSelectors {
		if chainSel != e.HomeChainSel {
			nonHomeChains = append(nonHomeChains, chainSel)
		}
	}
	require.Len(t, nonHomeChains, 2)

	homeChainConfig, err := state.Chains[e.HomeChainSel].CCIPHome.GetChainConfig(&bind.CallOpts{
		Context: ctx,
	}, e.HomeChainSel)
	require.NoError(t, err)
	require.Equalf(t, homeChainConfig.FChain, uint8(fRoleDON), "home chain fChain must be %d, got %d", fRoleDON, homeChainConfig.FChain)

	// get the fChain values of the chains to determine which will be the source
	// and which will be the dest.
	chainConfig0, err := state.Chains[e.HomeChainSel].CCIPHome.GetChainConfig(&bind.CallOpts{
		Context: ctx,
	}, nonHomeChains[0])
	require.NoError(t, err)
	chainConfig1, err := state.Chains[e.HomeChainSel].CCIPHome.GetChainConfig(&bind.CallOpts{
		Context: ctx,
	}, nonHomeChains[1])
	require.NoError(t, err)
	// the fChain values must be different, otherwise setup is incorrect.
	require.NotEqual(t, chainConfig0.FChain, chainConfig1.FChain)

	var sourceChain, destChain uint64
	var sourceChainConfig, destChainConfig ccip_home.CCIPHomeChainConfig
	if chainConfig0.FChain == fChainSource {
		sourceChain = nonHomeChains[0]
		destChain = nonHomeChains[1]
		sourceChainConfig = chainConfig0
		destChainConfig = chainConfig1
	} else {
		sourceChain = nonHomeChains[1]
		destChain = nonHomeChains[0]
		sourceChainConfig = chainConfig1
		destChainConfig = chainConfig0
	}

	t.Logf("home chain: %d, source chain: %d, dest chain: %d", e.HomeChainSel, sourceChain, destChain)
	t.Logf("source chain is supported by %d readers", len(sourceChainConfig.Readers))
	t.Logf("dest chain is supported by %d readers", len(destChainConfig.Readers))

	// now we can set up the lane.
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)

	var (
		replayed bool
		nonce    uint64
		sender   = common.LeftPadBytes(e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey.From.Bytes(), 32)
		out      mt.TestCaseOutput
		setup    = mt.NewTestSetupWithDeployedEnv(
			t,
			e,
			state,
			sourceChain,
			destChain,
			sender,
			false, // testRouter
		)
	)

	monitorCtx, monitorCancel := context.WithCancel(ctx)
	ms := &monitorState{}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		monitorReExecutions(monitorCtx, t, state, destChain, ms)
	}()

	t.Run("data message to eoa", func(t *testing.T) {
		out = mt.Run(
			t,
			mt.TestCase{
				ValidationType:         mt.ValidationTypeExec,
				TestSetup:              setup,
				Replayed:               replayed,
				Nonce:                  &nonce,
				Receiver:               common.HexToAddress("0xdead").Bytes(),
				MsgData:                []byte("hello eoa"),
				ExtraArgs:              nil,                                 // default extraArgs
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS, // success because offRamp won't call an EOA
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) {
					},
				},
			},
		)
	})

	_ = out

	monitorCancel()
	wg.Wait()
	// there should be no re-executions.
	require.Equal(t, int32(0), ms.reExecutionsObserved.Load())
}
