package smoke

import (
	"math/big"
	"os"
	"strconv"
	"testing"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/osutil"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipdeployment "github.com/smartcontractkit/chainlink/deployment/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_remote"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"

	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testsetups"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestRMN_TwoMessagesOnTwoLanesIncludingBatching(t *testing.T) {
	runRmnTestCase(t, rmnTestCase{
		name:        "messages on two lanes including batching",
		waitForExec: true,
		homeChainConfig: homeChainConfig{
			f: map[int]int{chain0: 1, chain1: 1},
		},
		remoteChainsConfig: []remoteChainConfig{
			{chainIdx: chain0, f: 1},
			{chainIdx: chain1, f: 1},
		},
		rmnNodes: []rmnNode{
			{id: 0, isSigner: true, observedChainIdxs: []int{chain0, chain1}},
			{id: 1, isSigner: true, observedChainIdxs: []int{chain0, chain1}},
			{id: 2, isSigner: true, observedChainIdxs: []int{chain0, chain1}},
		},
		messagesToSend: []messageToSend{
			{fromChainIdx: chain0, toChainIdx: chain1, count: 1},
			{fromChainIdx: chain1, toChainIdx: chain0, count: 5},
		},
	})
}

func TestRMN_MultipleMessagesOnOneLaneNoWaitForExec(t *testing.T) {
	runRmnTestCase(t, rmnTestCase{
		name:        "multiple messages for rmn batching inspection and one rmn node down",
		waitForExec: false, // do not wait for execution reports
		homeChainConfig: homeChainConfig{
			f: map[int]int{chain0: 1, chain1: 1},
		},
		remoteChainsConfig: []remoteChainConfig{
			{chainIdx: chain0, f: 1},
			{chainIdx: chain1, f: 1},
		},
		rmnNodes: []rmnNode{
			{id: 0, isSigner: true, observedChainIdxs: []int{chain0, chain1}},
			{id: 1, isSigner: true, observedChainIdxs: []int{chain0, chain1}},
			{id: 2, isSigner: true, observedChainIdxs: []int{chain0, chain1}, forceExit: true}, // one rmn node is down
		},
		messagesToSend: []messageToSend{
			{fromChainIdx: chain1, toChainIdx: chain0, count: 10},
		},
	})
}

func TestRMN_NotEnoughObservers(t *testing.T) {
	runRmnTestCase(t, rmnTestCase{
		name:                "one message but not enough observers, should not get a commit report",
		passIfNoCommitAfter: time.Minute, // wait for a minute and assert that commit report was not delivered
		homeChainConfig: homeChainConfig{
			f: map[int]int{chain0: 1, chain1: 1},
		},
		remoteChainsConfig: []remoteChainConfig{
			{chainIdx: chain0, f: 1},
			{chainIdx: chain1, f: 1},
		},
		rmnNodes: []rmnNode{
			{id: 0, isSigner: true, observedChainIdxs: []int{chain0, chain1}},
			{id: 1, isSigner: true, observedChainIdxs: []int{chain0, chain1}, forceExit: true},
			{id: 2, isSigner: true, observedChainIdxs: []int{chain0, chain1}, forceExit: true},
		},
		messagesToSend: []messageToSend{
			{fromChainIdx: chain0, toChainIdx: chain1, count: 1},
		},
	})
}

func TestRMN_DifferentSigners(t *testing.T) {
	runRmnTestCase(t, rmnTestCase{
		name: "different signers and different observers",
		homeChainConfig: homeChainConfig{
			f: map[int]int{chain0: 1, chain1: 1},
		},
		remoteChainsConfig: []remoteChainConfig{
			{chainIdx: chain0, f: 1},
			{chainIdx: chain1, f: 1},
		},
		rmnNodes: []rmnNode{
			{id: 0, isSigner: false, observedChainIdxs: []int{chain0, chain1}},
			{id: 1, isSigner: false, observedChainIdxs: []int{chain0, chain1}},
			{id: 2, isSigner: false, observedChainIdxs: []int{chain0, chain1}},
			{id: 3, isSigner: true, observedChainIdxs: []int{}},
			{id: 4, isSigner: true, observedChainIdxs: []int{}},
			{id: 5, isSigner: true, observedChainIdxs: []int{}},
		},
		messagesToSend: []messageToSend{
			{fromChainIdx: chain0, toChainIdx: chain1, count: 1},
		},
	})
}

func TestRMN_NotEnoughSigners(t *testing.T) {
	runRmnTestCase(t, rmnTestCase{
		name:                "different signers and different observers",
		passIfNoCommitAfter: time.Minute, // wait for a minute and assert that commit report was not delivered
		homeChainConfig: homeChainConfig{
			f: map[int]int{chain0: 1, chain1: 1},
		},
		remoteChainsConfig: []remoteChainConfig{
			{chainIdx: chain0, f: 1},
			{chainIdx: chain1, f: 1},
		},
		rmnNodes: []rmnNode{
			{id: 0, isSigner: false, observedChainIdxs: []int{chain0, chain1}},
			{id: 1, isSigner: false, observedChainIdxs: []int{chain0, chain1}},
			{id: 2, isSigner: false, observedChainIdxs: []int{chain0, chain1}},
			{id: 3, isSigner: true, observedChainIdxs: []int{}},
			{id: 4, isSigner: true, observedChainIdxs: []int{}, forceExit: true}, // signer is down
			{id: 5, isSigner: true, observedChainIdxs: []int{}, forceExit: true}, // signer is down
		},
		messagesToSend: []messageToSend{
			{fromChainIdx: chain0, toChainIdx: chain1, count: 1},
		},
	})
}

func TestRMN_DifferentRmnNodesForDifferentChains(t *testing.T) {
	runRmnTestCase(t, rmnTestCase{
		name:        "different rmn nodes support different chains",
		waitForExec: false,
		homeChainConfig: homeChainConfig{
			f: map[int]int{chain0: 1, chain1: 1},
		},
		remoteChainsConfig: []remoteChainConfig{
			{chainIdx: chain0, f: 1},
			{chainIdx: chain1, f: 1},
		},
		rmnNodes: []rmnNode{
			{id: 0, isSigner: true, observedChainIdxs: []int{chain0}},
			{id: 1, isSigner: true, observedChainIdxs: []int{chain0}},
			{id: 2, isSigner: true, observedChainIdxs: []int{chain0}},
			{id: 3, isSigner: true, observedChainIdxs: []int{chain1}},
			{id: 4, isSigner: true, observedChainIdxs: []int{chain1}},
			{id: 5, isSigner: true, observedChainIdxs: []int{chain1}},
		},
		messagesToSend: []messageToSend{
			{fromChainIdx: chain0, toChainIdx: chain1, count: 1},
			{fromChainIdx: chain1, toChainIdx: chain0, count: 1},
		},
	})
}

const (
	chain0 = 0
	chain1 = 1
)

func runRmnTestCase(t *testing.T, tc rmnTestCase) {
	require.NoError(t, os.Setenv("ENABLE_RMN", "true"))

	envWithRMN, rmnCluster := testsetups.NewLocalDevEnvironmentWithRMN(t, logger.TestLogger(t), len(tc.rmnNodes))
	t.Logf("envWithRmn: %#v", envWithRMN)

	var chainSelectors []uint64
	for _, chain := range envWithRMN.Env.Chains {
		chainSelectors = append(chainSelectors, chain.Selector)
	}
	require.Greater(t, len(chainSelectors), 1, "There should be at least two chains")

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

		var offchainPublicKey [32]byte
		copy(offchainPublicKey[:], rmn.RMN.OffchainPublicKey)

		rmnHomeNodes = append(rmnHomeNodes, rmn_home.RMNHomeNode{
			PeerId:            rmn.Proxy.PeerID,
			OffchainPublicKey: offchainPublicKey,
		})

		if rmnNodeInfo.isSigner {
			rmnRemoteSigners = append(rmnRemoteSigners, rmn_remote.RMNRemoteSigner{
				OnchainPublicKey: rmn.RMN.EVMOnchainPublicKey,
				NodeIndex:        uint64(rmnNodeInfo.id),
			})
		}
	}

	var rmnHomeSourceChains []rmn_home.RMNHomeSourceChain
	for remoteChainIdx, remoteF := range tc.homeChainConfig.f {
		// configure remote chain details on the home contract
		rmnHomeSourceChains = append(rmnHomeSourceChains, rmn_home.RMNHomeSourceChain{
			ChainSelector:       chainSelectors[remoteChainIdx],
			F:                   uint64(remoteF),
			ObserverNodesBitmap: createObserverNodesBitmap(chainSelectors[remoteChainIdx], tc.rmnNodes, chainSelectors),
		})
	}

	onChainState, err := ccipdeployment.LoadOnchainState(envWithRMN.Env)
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
		remoteSel := chainSelectors[remoteCfg.chainIdx]
		chState, ok := onChainState.Chains[remoteSel]
		require.True(t, ok)
		rmnRemoteConfig := rmn_remote.RMNRemoteConfig{
			RmnHomeContractConfigDigest: activeDigest,
			Signers:                     rmnRemoteSigners,
			F:                           uint64(remoteCfg.f),
		}

		chain := envWithRMN.Env.Chains[chainSelectors[remoteCfg.chainIdx]]

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

	// Kill the RMN nodes that are marked for force exit
	for _, n := range tc.rmnNodes {
		if n.forceExit {
			t.Logf("Pausing RMN node %d", n.id)
			rmnN := rmnCluster.Nodes["rmn_"+strconv.Itoa(n.id)]
			require.NoError(t, osutil.ExecCmd(zerolog.Nop(), "docker kill "+rmnN.Proxy.ContainerName))
			t.Logf("Paused RMN node %d", n.id)
		}
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
		fromChain := chainSelectors[msg.fromChainIdx]
		toChain := chainSelectors[msg.toChainIdx]

		for i := 0; i < msg.count; i++ {
			msgSentEvent := ccipdeployment.TestSendRequest(t, envWithRMN.Env, onChainState, fromChain, toChain, false, router.ClientEVM2AnyMessage{
				Receiver:     common.LeftPadBytes(onChainState.Chains[toChain].Receiver.Address().Bytes(), 32),
				Data:         []byte("hello world"),
				TokenAmounts: nil,
				FeeToken:     common.HexToAddress("0x0"),
				ExtraArgs:    nil,
			})
			expectedSeqNum[toChain] = msgSentEvent.SequenceNumber
			t.Logf("Sent message from chain %d to chain %d with seqNum %d", fromChain, toChain, msgSentEvent.SequenceNumber)
		}

		zero := uint64(0)
		startBlocks[toChain] = &zero
	}
	t.Logf("Sent all messages, expectedSeqNum: %v", expectedSeqNum)

	commitReportReceived := make(chan struct{})
	go func() {
		ccipdeployment.ConfirmCommitForAllWithExpectedSeqNums(t, envWithRMN.Env, onChainState, expectedSeqNum, startBlocks)
		commitReportReceived <- struct{}{}
	}()

	if tc.passIfNoCommitAfter > 0 { // wait for a duration and assert that commit reports were not delivered
		tim := time.NewTimer(tc.passIfNoCommitAfter)
		t.Logf("waiting for %s before asserting that commit report was not received", tc.passIfNoCommitAfter)
		select {
		case <-commitReportReceived:
			t.Errorf("Commit report was received while it was not expected")
			return
		case <-tim.C:
			return
		}
	}

	t.Logf("⌛ Waiting for commit reports...")
	<-commitReportReceived // wait for commit reports
	t.Logf("✅ Commit report")

	if tc.waitForExec {
		t.Logf("⌛ Waiting for exec reports...")
		ccipdeployment.ConfirmExecWithSeqNrForAll(t, envWithRMN.Env, onChainState, expectedSeqNum, startBlocks)
		t.Logf("✅ Exec report")
	}
}

func createObserverNodesBitmap(chainSel uint64, rmnNodes []rmnNode, chainSelectors []uint64) *big.Int {
	bitmap := new(big.Int)
	for _, n := range rmnNodes {
		observedChainSelectors := mapset.NewSet[uint64]()
		for _, chainIdx := range n.observedChainIdxs {
			observedChainSelectors.Add(chainSelectors[chainIdx])
		}

		if !observedChainSelectors.Contains(chainSel) {
			continue
		}

		bitmap.SetBit(bitmap, n.id, 1)
	}

	return bitmap
}

type homeChainConfig struct {
	f map[int]int
}

type remoteChainConfig struct {
	chainIdx int
	f        int
}

type rmnNode struct {
	id                int
	isSigner          bool
	observedChainIdxs []int
	forceExit         bool // force exit will simply force exit the rmn node to simulate failure scenarios
}

type messageToSend struct {
	fromChainIdx int
	toChainIdx   int
	count        int
}

type rmnTestCase struct {
	name string
	// If set to 0, the test will wait for commit reports.
	// If set to a positive value, the test will wait for that duration and will assert that commit report was not delivered.
	passIfNoCommitAfter time.Duration
	waitForExec         bool
	homeChainConfig     homeChainConfig
	remoteChainsConfig  []remoteChainConfig
	rmnNodes            []rmnNode
	messagesToSend      []messageToSend
}
