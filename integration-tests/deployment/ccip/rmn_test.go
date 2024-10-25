package ccipdeployment

import (
	"math/big"
	"os"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_remote"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestRMN(t *testing.T) {
	t.Skip("Local only")

	require.NoError(t, os.Setenv("ENABLE_RMN", "true"))

	envWithRMN, rmnCluster := NewLocalDevEnvironmentWithRMN(t, logger.TestLogger(t))
	t.Logf("envWithRmn: %#v", envWithRMN)

	var (
		rmnHomeNodes     []rmn_home.RMNHomeNode
		rmnRemoteSigners []rmn_remote.RMNRemoteSigner
		nodeIndex        uint64
	)

	for rmnNode, rmn := range rmnCluster.Nodes {
		t.Log(rmnNode, rmn.Proxy.PeerID, rmn.RMN.OffchainPublicKey, rmn.RMN.EVMOnchainPublicKey)
		var offchainPublicKey [32]byte
		copy(offchainPublicKey[:], rmn.RMN.OffchainPublicKey)
		rmnHomeNodes = append(rmnHomeNodes, rmn_home.RMNHomeNode{
			PeerId:            rmn.Proxy.PeerID,
			OffchainPublicKey: offchainPublicKey,
		})
		rmnRemoteSigners = append(rmnRemoteSigners, rmn_remote.RMNRemoteSigner{
			OnchainPublicKey: rmn.RMN.EVMOnchainPublicKey,
			NodeIndex:        nodeIndex,
		})
		nodeIndex++
	}

	var rmnHomeSourceChains []rmn_home.RMNHomeSourceChain
	for _, chain := range envWithRMN.Env.Chains {
		rmnHomeSourceChains = append(rmnHomeSourceChains, rmn_home.RMNHomeSourceChain{
			ChainSelector:       chain.Selector,
			MinObservers:        1,
			ObserverNodesBitmap: createObserverNodesBitmap(len(rmnHomeNodes)),
		})
	}

	onChainState, err := LoadOnchainState(envWithRMN.Env, envWithRMN.Ab)
	require.NoError(t, err)
	t.Logf("onChainState: %#v", onChainState)

	homeChain, ok := envWithRMN.Env.Chains[envWithRMN.HomeChainSel]
	require.True(t, ok)

	homeChainState, ok := onChainState.Chains[envWithRMN.HomeChainSel]
	require.True(t, ok)

	allDigests, err := homeChainState.RMNHome.GetConfigDigests(&bind.CallOpts{
		Context: testcontext.Get(t),
	})
	require.NoError(t, err)

	t.Logf("RMNHome candidateDigest before setting new candidate: %x, activeDigest: %x",
		allDigests.CandidateConfigDigest[:], allDigests.ActiveConfigDigest[:])

	staticConfig := rmn_home.RMNHomeStaticConfig{
		Nodes:          rmnHomeNodes,
		OffchainConfig: []byte{},
	}
	dynamicConfig := rmn_home.RMNHomeDynamicConfig{
		SourceChains:   rmnHomeSourceChains,
		OffchainConfig: []byte{},
	}
	t.Logf("Setting RMNHome candidate with staticConfig: %+v, dynamicConfig: %+v, current candidateDigest: %x",
		staticConfig, dynamicConfig, allDigests.CandidateConfigDigest[:])
	tx, err := homeChainState.RMNHome.SetCandidate(homeChain.DeployerKey, staticConfig, dynamicConfig, allDigests.CandidateConfigDigest)
	require.NoError(t, err)

	_, err = deployment.ConfirmIfNoError(homeChain, tx, err)
	require.NoError(t, err)

	candidateDigest, err := homeChainState.RMNHome.GetCandidateDigest(&bind.CallOpts{
		Context: testcontext.Get(t),
	})
	require.NoError(t, err)

	t.Logf("RMNHome candidateDigest after setting new candidate: %x", candidateDigest[:])
	t.Logf("Promoting RMNHome candidate with candidateDigest: %x", candidateDigest[:])

	tx, err = homeChainState.RMNHome.PromoteCandidateAndRevokeActive(
		homeChain.DeployerKey, candidateDigest, allDigests.ActiveConfigDigest)
	require.NoError(t, err)

	_, err = deployment.ConfirmIfNoError(homeChain, tx, err)
	require.NoError(t, err)

	// check the active digest is the same as the candidate digest
	activeDigest, err := homeChainState.RMNHome.GetActiveDigest(&bind.CallOpts{
		Context: testcontext.Get(t),
	})
	require.NoError(t, err)
	require.Equalf(t, candidateDigest, activeDigest,
		"active digest should be the same as the previously candidate digest after promotion, previous candidate: %x, active: %x",
		candidateDigest[:], activeDigest[:])

	// Set RMN remote config appropriately
	for _, chain := range envWithRMN.Env.Chains {
		chState, ok := onChainState.Chains[chain.Selector]
		require.True(t, ok)
		rmnRemoteConfig := rmn_remote.RMNRemoteConfig{
			RmnHomeContractConfigDigest: activeDigest,
			Signers:                     rmnRemoteSigners,
			MinSigners:                  1,
		}
		t.Logf("Setting RMNRemote config with RMNHome active digest: %x, cfg: %+v", activeDigest[:], rmnRemoteConfig)
		tx2, err2 := chState.RMNRemote.SetConfig(chain.DeployerKey, rmnRemoteConfig)
		require.NoError(t, err2)
		_, err2 = deployment.ConfirmIfNoError(chain, tx2, err2)
		require.NoError(t, err2)

		// confirm the config is set correctly
		config, err2 := chState.RMNRemote.GetVersionedConfig(&bind.CallOpts{
			Context: testcontext.Get(t),
		})
		require.NoError(t, err2)
		require.Equalf(t,
			activeDigest,
			config.Config.RmnHomeContractConfigDigest,
			"RMNRemote config digest should be the same as the active digest of RMNHome after setting, RMNHome active: %x, RMNRemote config: %x",
			activeDigest[:], config.Config.RmnHomeContractConfigDigest[:])

		t.Logf("RMNRemote config digest after setting: %x", config.Config.RmnHomeContractConfigDigest[:])
	}

	jobSpecs, err := NewCCIPJobSpecs(envWithRMN.Env.NodeIDs, envWithRMN.Env.Offchain)
	require.NoError(t, err)

	ctx := Context(t)

	ReplayLogs(t, envWithRMN.Env.Offchain, envWithRMN.ReplayBlocks)

	for nodeID, jobs := range jobSpecs {
		for _, job := range jobs {
			_, err := envWithRMN.Env.Offchain.ProposeJob(ctx,
				&jobv1.ProposeJobRequest{
					NodeId: nodeID,
					Spec:   job,
				})
			require.NoError(t, err)
		}
	}

	// Add all lanes
	require.NoError(t, AddLanesForAll(envWithRMN.Env, onChainState))

	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)

	// Send one message from one chain to another.
	chains := maps.Values(envWithRMN.Env.Chains)
	t.Logf("all chains: %v", chains)

	sort.Slice(chains, func(i int, j int) bool { return chains[i].Selector < chains[j].Selector })
	srcChain := chains[0]
	dstChain := chains[1]
	require.True(t, srcChain.Selector != dstChain.Selector)
	t.Logf("source chain is %d dest chain is %d", srcChain.Selector, dstChain.Selector)

	latesthdr, err := dstChain.Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	block := latesthdr.Number.Uint64()
	startBlocks[dstChain.Selector] = &block
	seqNum := TestSendRequest(t, envWithRMN.Env, onChainState, srcChain.Selector, dstChain.Selector, false)
	t.Logf("expected seqNum: %d", seqNum)

	expectedSeqNum := make(map[uint64]uint64)
	expectedSeqNum[dstChain.Selector] = seqNum

	t.Logf("⌛ Waiting for commit report...")
	ConfirmCommitForAllWithExpectedSeqNums(t, envWithRMN.Env, onChainState, expectedSeqNum, startBlocks)
	t.Logf("✅ Commit report")

	t.Logf("⌛ Waiting for exec report...")
	ConfirmExecWithSeqNrForAll(t, envWithRMN.Env, onChainState, expectedSeqNum, startBlocks)
	t.Logf("✅ Exec report")
}

func createObserverNodesBitmap(numNodes int) *big.Int {
	// for now, all nodes support all chains, so the bitmap is all 1s.
	// first, initialize a big.Int with all bits set to 0.
	// then, set the first numNodes bits to 1.
	bitmap := new(big.Int)
	for i := 0; i < numNodes; i++ {
		bitmap.SetBit(bitmap, i, 1)
	}
	return bitmap
}
