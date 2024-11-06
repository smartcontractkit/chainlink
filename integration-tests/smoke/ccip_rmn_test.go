package smoke

import (
	"math/big"
	"os"
	"strconv"
	"testing"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink/deployment"
	ccipdeployment "github.com/smartcontractkit/chainlink/deployment/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_remote"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testsetups"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// Set false to run the RMN tests
const skipRmnTest = false

const (
	// TODO: get the selectors dynamically from the config
	chain0 = 12922642891491394802
	chain1 = 3379446385462418246
)

func TestRMN_TwoMessagesOnTwoLanes(t *testing.T) {
	runRmnTestCase(t, rmnTestCase{
		name:        "two messages on two lanes",
		waitTimeout: 10 * time.Minute,
		waitForExec: true,
		homeChainConfig: homeChainConfig{
			f: map[uint64]int{chain0: 1, chain1: 1},
		},
		remoteChainsConfig: []remoteChainConfig{
			{selector: chain0, f: 1},
			{selector: chain1, f: 1},
		},
		rmnNodes: []rmnNode{
			{id: 0, isSigner: true, observedChains: []uint64{chain0, chain1}},
			{id: 1, isSigner: true, observedChains: []uint64{chain0, chain1}},
			{id: 2, isSigner: true, observedChains: []uint64{chain0, chain1}},
		},
		messagesToSend: []messageToSend{
			{fromChain: chain0, toChain: chain1, count: 1, expectedDelivered: true},
			{fromChain: chain1, toChain: chain0, count: 1, expectedDelivered: true},
		},
	})
}

func TestRMN_MultipleMessagesOnOneLaneNoWaitForExec(t *testing.T) {
	runRmnTestCase(t, rmnTestCase{
		name:        "multiple messages on two lanes for rmn batching inspection",
		waitTimeout: 10 * time.Minute,
		waitForExec: false, // do not wait for execution reports
		homeChainConfig: homeChainConfig{
			f: map[uint64]int{chain0: 1, chain1: 1},
		},
		remoteChainsConfig: []remoteChainConfig{
			{selector: chain0, f: 1},
			{selector: chain1, f: 1},
		},
		rmnNodes: []rmnNode{
			{id: 0, isSigner: true, observedChains: []uint64{chain0, chain1}},
			{id: 1, isSigner: true, observedChains: []uint64{chain0, chain1}},
			{id: 2, isSigner: true, observedChains: []uint64{chain0, chain1}},
		},
		messagesToSend: []messageToSend{
			{fromChain: chain1, toChain: chain0, count: 10, expectedDelivered: true},
		},
	})
}

func runRmnTestCase(t *testing.T, tc rmnTestCase) {
	if skipRmnTest {
		t.Skip("Local only")
	}
	require.NoError(t, os.Setenv("ENABLE_RMN", "true"))

	envWithRMN, rmnCluster := testsetups.NewLocalDevEnvironmentWithRMN(t, logger.TestLogger(t), len(tc.rmnNodes))
	t.Logf("envWithRmn: %#v", envWithRMN)

	remoteChainSelectors := make([]uint64, 0, len(envWithRMN.Env.Chains)-1)
	for _, chain := range envWithRMN.Env.Chains {
		remoteChainSelectors = append(remoteChainSelectors, chain.Selector)
	}
	require.Greater(t, len(remoteChainSelectors), 0, "There should be at least one remote chain")

	var (
		rmnHomeNodes     []rmn_home.RMNHomeNode
		rmnRemoteSigners []rmn_remote.RMNRemoteSigner
	)

	for _, rmnNodeInfo := range tc.rmnNodes {
		rmn := rmnCluster.Nodes["rmn_"+strconv.Itoa(rmnNodeInfo.id)]

		t.Log(rmnNodeInfo.id, rmn.Proxy.PeerID, rmn.RMN.OffchainPublicKey, rmn.RMN.EVMOnchainPublicKey)

		var offchainPublicKey [32]byte
		copy(offchainPublicKey[:], rmn.RMN.OffchainPublicKey)

		rmnHomeNodes = append(rmnHomeNodes, rmn_home.RMNHomeNode{
			PeerId:            rmn.Proxy.PeerID,
			OffchainPublicKey: offchainPublicKey,
		})

		rmnRemoteSigners = append(rmnRemoteSigners, rmn_remote.RMNRemoteSigner{
			OnchainPublicKey: rmn.RMN.EVMOnchainPublicKey,
			NodeIndex:        uint64(rmnNodeInfo.id),
		})
	}

	var rmnHomeSourceChains []rmn_home.RMNHomeSourceChain
	for remoteChainSel, remoteF := range tc.homeChainConfig.f {
		// configure remote chain details on the home contract
		rmnHomeSourceChains = append(rmnHomeSourceChains, rmn_home.RMNHomeSourceChain{
			ChainSelector:       remoteChainSel,
			F:                   uint64(remoteF),
			ObserverNodesBitmap: createObserverNodesBitmap(remoteChainSel, tc.rmnNodes),
		})
	}

	onChainState, err := ccipdeployment.LoadOnchainState(envWithRMN.Env, envWithRMN.Ab)
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
	for _, remoteCfg := range tc.remoteChainsConfig {
		chState, ok := onChainState.Chains[remoteCfg.selector]
		require.True(t, ok)
		rmnRemoteConfig := rmn_remote.RMNRemoteConfig{
			RmnHomeContractConfigDigest: activeDigest,
			Signers:                     rmnRemoteSigners,
			F:                           uint64(remoteCfg.f),
		}

		chain := envWithRMN.Env.Chains[remoteCfg.selector]

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

	jobSpecs, err := ccipdeployment.NewCCIPJobSpecs(envWithRMN.Env.NodeIDs, envWithRMN.Env.Offchain)
	require.NoError(t, err)

	ctx := ccipdeployment.Context(t)

	ccipdeployment.ReplayLogs(t, envWithRMN.Env.Offchain, envWithRMN.ReplayBlocks)

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
	require.NoError(t, ccipdeployment.AddLanesForAll(envWithRMN.Env, onChainState))

	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	expectedSeqNum := make(map[uint64]uint64)
	for _, msg := range tc.messagesToSend {
		for i := 0; i < msg.count; i++ {
			seqNum := ccipdeployment.TestSendRequest(t, envWithRMN.Env, onChainState, msg.fromChain, msg.toChain, false)
			expectedSeqNum[msg.toChain] = seqNum
			t.Logf("Sent message from chain %d to chain %d with seqNum %d", msg.fromChain, msg.toChain, seqNum)
		}
		zero := uint64(0)
		startBlocks[msg.toChain] = &zero
	}
	t.Logf("Sent all messages, expectedSeqNum: %v", expectedSeqNum)

	t.Logf("⌛ Waiting for commit reports...")
	ccipdeployment.ConfirmCommitForAllWithExpectedSeqNums(t, envWithRMN.Env, onChainState, expectedSeqNum, startBlocks)
	t.Logf("✅ Commit report")

	if tc.waitForExec {
		t.Logf("⌛ Waiting for exec reports...")
		ccipdeployment.ConfirmExecWithSeqNrForAll(t, envWithRMN.Env, onChainState, expectedSeqNum, startBlocks)
		t.Logf("✅ Exec report")
	}
}

func createObserverNodesBitmap(chainSel uint64, rmnNodes []rmnNode) *big.Int {
	bitmap := new(big.Int)
	for _, n := range rmnNodes {
		observedChains := mapset.NewSet(n.observedChains...)
		if !observedChains.Contains(chainSel) {
			continue
		}
		bitmap.SetBit(bitmap, n.id, 1)
	}
	return bitmap
}

type homeChainConfig struct {
	f map[uint64]int
}

type remoteChainConfig struct {
	selector uint64
	f        int
}

type rmnNode struct {
	id             int
	isSigner       bool
	observedChains []uint64
}

type messageToSend struct {
	fromChain         uint64
	toChain           uint64
	count             int
	expectedDelivered bool
}

type rmnTestCase struct {
	name               string
	waitTimeout        time.Duration
	waitForExec        bool
	homeChainConfig    homeChainConfig
	remoteChainsConfig []remoteChainConfig
	rmnNodes           []rmnNode
	messagesToSend     []messageToSend
}
