package ccip

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	ctf_client "github.com/smartcontractkit/chainlink-testing-framework/lib/client"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/osutil"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	ccipseq "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_home"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_remote"

	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

func TestRMN_IncorrectSig(t *testing.T) {
	runRmnTestCase(t, rmnTestCase{
		nodesWithIncorrectSigner: []int{0, 1},
		name:                     "messages with incorrect RMN signature",
		waitForExec:              true,
		passIfNoCommitAfter:      15 * time.Second,
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
		},
	})
}

func TestRMN_TwoMessagesOnTwoLanesIncludingBatching(t *testing.T) {
	runRmnTestCase(t, rmnTestCase{
		name:        "messages on two lanes including batching one lane RMN-enabled the other RMN-disabled",
		waitForExec: true,
		homeChainConfig: homeChainConfig{
			f: map[int]int{
				chain0: 1,
				//chain1: RMN-Disabled if no f defined
			},
		},
		remoteChainsConfig: []remoteChainConfig{
			{chainIdx: chain0, f: 1},
			{chainIdx: chain1, f: 1},
		},
		rmnNodes: []rmnNode{
			{id: 0, isSigner: true, observedChainIdxs: []int{chain0}},
			{id: 1, isSigner: true, observedChainIdxs: []int{chain0}},
			{id: 2, isSigner: true, observedChainIdxs: []int{chain0}},
		},
		messagesToSend: []messageToSend{
			{fromChainIdx: chain0, toChainIdx: chain1, count: 1},
			{fromChainIdx: chain1, toChainIdx: chain0, count: 5},
		},
	})
}

func TestRMN_SimpleVerificationDisabledOnDestination(t *testing.T) {
	runRmnTestCase(t, rmnTestCase{
		name:        "messages on two lanes one lane RMN-enabled the other RMN-disabled",
		waitForExec: true,
		homeChainConfig: homeChainConfig{
			f: map[int]int{
				chain1: 1,
			},
		},
		remoteChainsConfig: []remoteChainConfig{
			{chainIdx: chain2, f: 1},
		},
		rmnNodes: []rmnNode{
			{id: 0, isSigner: true, observedChainIdxs: []int{chain1, chain2}},
			{id: 1, isSigner: true, observedChainIdxs: []int{chain1, chain2}},
			{id: 2, isSigner: true, observedChainIdxs: []int{chain1, chain2}},
		},
		messagesToSend: []messageToSend{
			{fromChainIdx: chain0, toChainIdx: chain2, count: 1},
			{fromChainIdx: chain1, toChainIdx: chain2, count: 1},
		},
	})
}

func TestRMN_TwoMessagesOnTwoLanesIncludingBatchingWithTemporaryPause(t *testing.T) {
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
			{id: 0, isSigner: true, observedChainIdxs: []int{chain0, chain1}, forceExit: true, restart: true},
			{id: 1, isSigner: true, observedChainIdxs: []int{chain0, chain1}, forceExit: true, restart: true},
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
		passIfNoCommitAfter: 15 * time.Second,
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
		passIfNoCommitAfter: 15 * time.Second,
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

func TestRMN_TwoMessagesOneSourceChainCursed(t *testing.T) {
	runRmnTestCase(t, rmnTestCase{
		name:                "two messages, one source chain is cursed the other chain was cursed but curse is revoked",
		passIfNoCommitAfter: 15 * time.Second,
		cursedSubjectsPerChain: map[int][]int{
			chain1: {chain0},
		},
		revokedCursedSubjectsPerChain: map[int]map[int]time.Duration{
			chain0: {globalCurse: 5 * time.Second}, // chain0 will be globally cursed and curse will be revoked later
		},
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
			{fromChainIdx: chain0, toChainIdx: chain1, count: 1}, // <----- this message should not be committed
			{fromChainIdx: chain1, toChainIdx: chain0, count: 1},
		},
	})
}

func TestRMN_GlobalCurseTwoMessagesOnTwoLanes(t *testing.T) {
	runRmnTestCase(t, rmnTestCase{
		name:        "global curse messages on two lanes",
		waitForExec: false,
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
		cursedSubjectsPerChain: map[int][]int{
			chain1: {globalCurse},
			chain0: {globalCurse},
		},
		passIfNoCommitAfter: 15 * time.Second,
	})
}

const (
	chain0      = 0
	chain1      = 1
	chain2      = 2
	globalCurse = 1000
)

func runRmnTestCase(t *testing.T, tc rmnTestCase) {
	require.NoError(t, tc.validate())

	ctx := testcontext.Get(t)
	t.Logf("Running RMN test case: %s", tc.name)

	envWithRMN, rmnCluster, _ := testsetups.NewIntegrationEnvironment(t,
		testhelpers.WithRMNEnabled(len(tc.rmnNodes)),
	)
	tc.populateFields(t, envWithRMN, rmnCluster)

	onChainState, err := stateview.LoadOnchainState(envWithRMN.Env)
	require.NoError(t, err)

	homeChainState, ok := onChainState.EVMChainState(envWithRMN.HomeChainSel)
	require.True(t, ok)

	allDigests, err := homeChainState.RMNHome.GetConfigDigests(&bind.CallOpts{Context: ctx})
	require.NoError(t, err)

	t.Logf("RMNHome candidateDigest before setting new candidate: %x, activeDigest: %x",
		allDigests.CandidateConfigDigest[:], allDigests.ActiveConfigDigest[:])

	staticConfig := rmn_home.RMNHomeStaticConfig{Nodes: tc.pf.rmnHomeNodes, OffchainConfig: []byte{}}
	dynamicConfig := rmn_home.RMNHomeDynamicConfig{SourceChains: tc.pf.rmnHomeSourceChains, OffchainConfig: []byte{}}
	t.Logf("Setting RMNHome candidate with staticConfig: %+v, dynamicConfig: %+v, current candidateDigest: %x",
		staticConfig, dynamicConfig, allDigests.CandidateConfigDigest[:])

	candidateDigest, err := homeChainState.RMNHome.GetCandidateDigest(&bind.CallOpts{Context: ctx})
	require.NoError(t, err)

	_, err = v1_6.SetRMNHomeCandidateConfigChangeset(envWithRMN.Env, v1_6.SetRMNHomeCandidateConfig{
		HomeChainSelector: envWithRMN.HomeChainSel,
		RMNStaticConfig:   staticConfig,
		RMNDynamicConfig:  dynamicConfig,
		DigestToOverride:  candidateDigest,
	})
	require.NoError(t, err)

	candidateDigest, err = homeChainState.RMNHome.GetCandidateDigest(&bind.CallOpts{Context: ctx})
	require.NoError(t, err)

	t.Logf("RMNHome candidateDigest after setting new candidate: %x", candidateDigest[:])
	t.Logf("Promoting RMNHome candidate with candidateDigest: %x", candidateDigest[:])

	_, err = v1_6.PromoteRMNHomeCandidateConfigChangeset(envWithRMN.Env, v1_6.PromoteRMNHomeCandidateConfig{
		HomeChainSelector: envWithRMN.HomeChainSel,
		DigestToPromote:   candidateDigest,
	})
	require.NoError(t, err)

	// check the active digest is the same as the candidate digest
	activeDigest, err := homeChainState.RMNHome.GetActiveDigest(&bind.CallOpts{Context: ctx})
	require.NoError(t, err)
	require.Equalf(t, candidateDigest, activeDigest,
		"active digest should be the same as the previously candidate digest after promotion, previous candidate: %x, active: %x",
		candidateDigest[:], activeDigest[:])

	rmnRemoteConfig := make(map[uint64]ccipops.RMNRemoteConfig)
	for _, remoteCfg := range tc.remoteChainsConfig {
		selector := tc.pf.chainSelectors[remoteCfg.chainIdx]
		if remoteCfg.f < 0 {
			t.Fatalf("remoteCfg.f is negative: %d", remoteCfg.f)
		}
		rmnRemoteConfig[selector] = ccipops.RMNRemoteConfig{
			F:       uint64(remoteCfg.f),
			Signers: tc.alterSigners(t, tc.pf.rmnRemoteSigners),
		}
	}

	_, err = v1_6.SetRMNRemoteConfigChangeset(envWithRMN.Env, ccipseq.SetRMNRemoteConfig{
		RMNRemoteConfigs: rmnRemoteConfig,
	})
	require.NoError(t, err)

	tc.killMarkedRmnNodes(t, rmnCluster)

	envWithRMN.RmnEnabledSourceChains = make(map[uint64]bool)
	for chainIdx := range tc.homeChainConfig.f {
		chainSel := tc.pf.chainSelectors[chainIdx]
		envWithRMN.RmnEnabledSourceChains[chainSel] = true
	}

	testhelpers.ReplayLogs(t, envWithRMN.Env.Offchain, envWithRMN.ReplayBlocks)
	testhelpers.AddLanesForAll(t, &envWithRMN, onChainState)
	disabledNodes := tc.disableOraclesIfThisIsACursingTestCase(ctx, t, envWithRMN)

	startBlocks, seqNumCommit, seqNumExec := tc.sendMessages(t, onChainState, envWithRMN)
	t.Logf("Sent all messages, seqNumCommit: %v seqNumExec: %v", seqNumCommit, seqNumExec)

	cleanup := tc.restartNode(t, rmnCluster)
	defer cleanup()

	eg := errgroup.Group{}
	tc.callContractsToCurseChains(ctx, t, onChainState, envWithRMN)
	tc.callContractsToCurseAndRevokeCurse(ctx, &eg, t, onChainState, envWithRMN)

	tc.enableOracles(ctx, t, envWithRMN, disabledNodes)

	expectedSeqNum := make(map[testhelpers.SourceDestPair]uint64)
	for k, v := range seqNumCommit {
		cursedSubjectsOfDest, exists := tc.pf.cursedSubjectsPerChainSel[k.DestChainSelector]
		shouldSkip := exists && (slices.Contains(cursedSubjectsOfDest, globalCurse) ||
			slices.Contains(cursedSubjectsOfDest, k.SourceChainSelector))

		if !shouldSkip {
			expectedSeqNum[k] = v
		}
	}

	t.Logf("expectedSeqNums: %v", expectedSeqNum)
	t.Logf("expectedSeqNums including cursed chains: %v", seqNumCommit)

	if len(tc.cursedSubjectsPerChain) > 0 && len(seqNumCommit) == len(expectedSeqNum) {
		t.Fatalf("test case is wrong: no message was sent to non-cursed chains when you " +
			"define curse subjects, your test case should have at least one message not expected to be delivered")
	}

	hasCommitReportBeenReceived := false
	// Trying to replay logs at intervals to avoid test flakiness
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if (hasCommitReportBeenReceived) || tc.passIfNoCommitAfter > 0 {
					return
				}
				// Do not assert on error as we replay logs to avoid race condition where nodes are being shut down and we call replay
				replayBlocks := make(map[uint64]uint64)
				for srcDestPair := range seqNumCommit {
					replayBlocks[srcDestPair.SourceChainSelector] = 1
					replayBlocks[srcDestPair.DestChainSelector] = 1
				}
				t.Logf("replaying logs after waiting for more than 1 minute (%v)", replayBlocks)
				testhelpers.ReplayLogs(t, envWithRMN.Env.Offchain, replayBlocks, testhelpers.WithAssertOnError(false))
			case <-t.Context().Done():
				return
			}
		}
	}()

	commitReportReceived := make(chan struct{})
	go func() {
		if len(expectedSeqNum) > 0 {
			testhelpers.ConfirmCommitForAllWithExpectedSeqNums(t, envWithRMN.Env, onChainState, expectedSeqNum, startBlocks)
			commitReportReceived <- struct{}{}
		}

		if len(seqNumCommit) > 0 && len(seqNumCommit) > len(expectedSeqNum) {
			// wait for a duration and assert that commit reports were not delivered for cursed source chains
			testhelpers.ConfirmCommitForAllWithExpectedSeqNums(t, envWithRMN.Env, onChainState, seqNumCommit, startBlocks)
			commitReportReceived <- struct{}{}
		}
	}()

	if tc.passIfNoCommitAfter > 0 { // wait for a duration and assert that commit reports were not delivered
		if len(expectedSeqNum) > 0 && len(seqNumCommit) > len(expectedSeqNum) {
			t.Logf("⌛ Waiting for commit reports of non-cursed chains...")
			<-commitReportReceived
			t.Logf("✅ Commit reports of non-cursed chains received")
		}

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
	hasCommitReportBeenReceived = true
	t.Logf("✅ Commit report")

	require.NoError(t, eg.Wait())

	if tc.waitForExec {
		t.Logf("⌛ Waiting for exec reports...")
		testhelpers.ConfirmExecWithSeqNrsForAll(t, envWithRMN.Env, onChainState, seqNumExec, startBlocks)
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
	restart           bool // restart will restart the rmn node to simulate failure scenarios
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
	passIfNoCommitAfter    time.Duration
	cursedSubjectsPerChain map[int][]int
	// revokedCursedSubjectsPerChain is used to revoke this specific curses after a timer expires
	revokedCursedSubjectsPerChain map[int]map[int]time.Duration // chainIdx -> subjectIdx -> timer to revoke
	waitForExec                   bool
	homeChainConfig               homeChainConfig
	remoteChainsConfig            []remoteChainConfig
	rmnNodes                      []rmnNode
	messagesToSend                []messageToSend
	nodesWithIncorrectSigner      []int

	// populated fields after environment setup
	pf testCasePopulatedFields
}

type testCasePopulatedFields struct {
	chainSelectors                   []uint64
	rmnHomeNodes                     []rmn_home.RMNHomeNode
	rmnRemoteSigners                 []rmn_remote.RMNRemoteSigner
	rmnHomeSourceChains              []rmn_home.RMNHomeSourceChain
	cursedSubjectsPerChainSel        map[uint64][]uint64
	revokedCursedSubjectsPerChainSel map[uint64]map[uint64]time.Duration
}

func (tc *rmnTestCase) alterSigners(t *testing.T, signers []rmn_remote.RMNRemoteSigner) []rmn_remote.RMNRemoteSigner {
	for _, n := range tc.nodesWithIncorrectSigner {
		for i, s := range signers {
			if n >= 0 && s.NodeIndex == uint64(n) {
				// Random address ethereum private key
				privateKey, err := crypto.GenerateKey()
				if err != nil {
					t.Fatalf("failed to generate private key: %v", err)
				}
				publicKey := privateKey.Public()
				publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
				if !ok {
					t.Fatalf("failed to cast public key to ECDSA")
				}
				address := crypto.PubkeyToAddress(*publicKeyECDSA)
				signers[i].OnchainPublicKey = address
			}
		}
	}

	return signers
}

func (tc *rmnTestCase) populateFields(t *testing.T, envWithRMN testhelpers.DeployedEnv, rmnCluster devenv.RMNCluster) {
	require.GreaterOrEqual(t, len(envWithRMN.Env.BlockChains.EVMChains()), 2, "test assumes at least two chains")
	for _, chain := range envWithRMN.Env.BlockChains.EVMChains() {
		tc.pf.chainSelectors = append(tc.pf.chainSelectors, chain.Selector)
	}

	for _, rmnNodeInfo := range tc.rmnNodes {
		rmn := rmnCluster.Nodes["rmn_"+strconv.Itoa(rmnNodeInfo.id)]

		var offchainPublicKey [32]byte
		copy(offchainPublicKey[:], rmn.RMN.OffchainPublicKey)

		tc.pf.rmnHomeNodes = append(tc.pf.rmnHomeNodes, rmn_home.RMNHomeNode{
			PeerId:            rmn.Proxy.PeerID,
			OffchainPublicKey: offchainPublicKey,
		})

		if rmnNodeInfo.isSigner {
			if rmnNodeInfo.id < 0 {
				t.Fatalf("node id is negative: %d", rmnNodeInfo.id)
			}
			tc.pf.rmnRemoteSigners = append(tc.pf.rmnRemoteSigners, rmn_remote.RMNRemoteSigner{
				OnchainPublicKey: rmn.RMN.EVMOnchainPublicKey,
				NodeIndex:        uint64(rmnNodeInfo.id),
			})
		}
	}

	for remoteChainIdx, remoteF := range tc.homeChainConfig.f {
		if remoteF < 0 {
			t.Fatalf("negative remote F: %d", remoteF)
		}
		// configure remote chain details on the home contract
		tc.pf.rmnHomeSourceChains = append(tc.pf.rmnHomeSourceChains, rmn_home.RMNHomeSourceChain{
			ChainSelector:       tc.pf.chainSelectors[remoteChainIdx],
			FObserve:            uint64(remoteF),
			ObserverNodesBitmap: createObserverNodesBitmap(tc.pf.chainSelectors[remoteChainIdx], tc.rmnNodes, tc.pf.chainSelectors),
		})
	}

	// populate cursed subjects with actual chain selectors
	tc.pf.cursedSubjectsPerChainSel = make(map[uint64][]uint64)
	for chainIdx, subjects := range tc.cursedSubjectsPerChain {
		chainSel := tc.pf.chainSelectors[chainIdx]
		for _, subject := range subjects {
			subjSel := uint64(globalCurse)
			if subject != globalCurse {
				subjSel = tc.pf.chainSelectors[subject]
			}
			tc.pf.cursedSubjectsPerChainSel[chainSel] = append(tc.pf.cursedSubjectsPerChainSel[chainSel], subjSel)
		}
	}

	// populate revoked cursed subjects with actual chain selectors
	tc.pf.revokedCursedSubjectsPerChainSel = make(map[uint64]map[uint64]time.Duration)
	for chainIdx, subjects := range tc.revokedCursedSubjectsPerChain {
		chainSel := tc.pf.chainSelectors[chainIdx]
		for subject, revokeAfter := range subjects {
			subjSel := uint64(globalCurse)
			if subject != globalCurse {
				subjSel = tc.pf.chainSelectors[subject]
			}
			if _, ok := tc.pf.revokedCursedSubjectsPerChainSel[chainSel]; !ok {
				tc.pf.revokedCursedSubjectsPerChainSel[chainSel] = make(map[uint64]time.Duration)
			}
			tc.pf.revokedCursedSubjectsPerChainSel[chainSel][subjSel] = revokeAfter
		}
	}
}

func (tc rmnTestCase) validate() error {
	if len(tc.cursedSubjectsPerChain) > 0 && tc.passIfNoCommitAfter == 0 {
		return errors.New("when you define cursed subjects you also need to define the duration that the " +
			"test will wait for non-transmitted roots")
	}
	return nil
}

func (tc rmnTestCase) killMarkedRmnNodes(t *testing.T, rmnCluster devenv.RMNCluster) {
	for _, n := range tc.rmnNodes {
		if n.forceExit {
			t.Logf("Pausing RMN node %d", n.id)
			rmnN := rmnCluster.Nodes["rmn_"+strconv.Itoa(n.id)]
			require.NoError(t, osutil.ExecCmd(zerolog.Nop(), "docker kill "+rmnN.Proxy.ContainerName))
			t.Logf("Paused RMN node %d", n.id)
		}
	}
}

func (tc rmnTestCase) restartNode(t *testing.T, rmnCluster devenv.RMNCluster) func() {
	errCh := make(chan error, 1)
	go func() {
		time.Sleep(10 * time.Second)
		for _, n := range tc.rmnNodes {
			if n.restart {
				t.Logf("Restarting RMN node %d", n.id)
				rmnN := rmnCluster.Nodes["rmn_"+strconv.Itoa(n.id)]
				if err := osutil.ExecCmd(zerolog.Nop(), "docker start "+rmnN.Proxy.ContainerName); err != nil {
					errCh <- err
					return
				}
				t.Logf("Restarted RMN node %d", n.id)
			}
		}
		errCh <- nil
	}()
	require.NoError(t, <-errCh)
	return func() {
		for _, n := range tc.rmnNodes {
			if n.restart {
				t.Logf("Stopping RMN node %d", n.id)
				rmnN := rmnCluster.Nodes["rmn_"+strconv.Itoa(n.id)]
				require.NoError(t, osutil.ExecCmd(zerolog.Nop(), "docker stop "+rmnN.Proxy.ContainerName))
				t.Logf("Stopped RMN node %d", n.id)
			}
		}
	}
}

func (tc rmnTestCase) disableOraclesIfThisIsACursingTestCase(ctx context.Context, t *testing.T, envWithRMN testhelpers.DeployedEnv) []string {
	disabledNodes := make([]string, 0)

	if len(tc.cursedSubjectsPerChain) > 0 {
		listNodesResp, err := envWithRMN.Env.Offchain.ListNodes(ctx, &node.ListNodesRequest{})
		require.NoError(t, err)

		for _, n := range listNodesResp.Nodes {
			if strings.HasPrefix(n.Name, "bootstrap") {
				continue
			}
			_, err := envWithRMN.Env.Offchain.DisableNode(ctx, &node.DisableNodeRequest{Id: n.Id})
			require.NoError(t, err)
			disabledNodes = append(disabledNodes, n.Id)
			t.Logf("node %s disabled", n.Id)
		}
	}

	return disabledNodes
}

func (tc rmnTestCase) sendMessages(t *testing.T, onChainState stateview.CCIPOnChainState, envWithRMN testhelpers.DeployedEnv) (map[uint64]*uint64, map[testhelpers.SourceDestPair]uint64, map[testhelpers.SourceDestPair][]uint64) {
	startBlocks := make(map[uint64]*uint64)
	seqNumCommit := make(map[testhelpers.SourceDestPair]uint64)
	seqNumExec := make(map[testhelpers.SourceDestPair][]uint64)

	for _, msg := range tc.messagesToSend {
		fromChain := tc.pf.chainSelectors[msg.fromChainIdx]
		toChain := tc.pf.chainSelectors[msg.toChainIdx]

		for i := 0; i < msg.count; i++ {
			msgSentEvent := testhelpers.TestSendRequest(t, envWithRMN.Env, onChainState, fromChain, toChain, false, router.ClientEVM2AnyMessage{
				Receiver:     common.LeftPadBytes(onChainState.MustGetEVMChainState(toChain).Receiver.Address().Bytes(), 32),
				Data:         []byte("hello world"),
				TokenAmounts: nil,
				FeeToken:     common.HexToAddress("0x0"),
				ExtraArgs:    nil,
			}, testhelpers.WithMaxRetries(5))
			seqNumCommit[testhelpers.SourceDestPair{
				SourceChainSelector: fromChain,
				DestChainSelector:   toChain,
			}] = msgSentEvent.SequenceNumber
			seqNumExec[testhelpers.SourceDestPair{
				SourceChainSelector: fromChain,
				DestChainSelector:   toChain,
			}] = []uint64{msgSentEvent.SequenceNumber}
			t.Logf("Sent message from chain %d to chain %d with seqNum %d", fromChain, toChain, msgSentEvent.SequenceNumber)
		}

		zero := uint64(0)
		startBlocks[toChain] = &zero
	}

	return startBlocks, seqNumCommit, seqNumExec
}

func (tc rmnTestCase) callContractsToCurseChains(ctx context.Context, t *testing.T, onChainState stateview.CCIPOnChainState, envWithRMN testhelpers.DeployedEnv) {
	for _, remoteCfg := range tc.remoteChainsConfig {
		remoteSel := tc.pf.chainSelectors[remoteCfg.chainIdx]
		chState, ok := onChainState.EVMChainState(remoteSel)
		require.True(t, ok)
		_, ok = envWithRMN.Env.BlockChains.EVMChains()[remoteSel]
		require.True(t, ok)

		cursedSubjects, ok := tc.cursedSubjectsPerChain[remoteCfg.chainIdx]
		if !ok {
			continue // nothing to curse on this chain
		}

		for _, subjectDescription := range cursedSubjects {
			curseActions := make([]v1_6.CurseAction, 0)

			if subjectDescription == globalCurse {
				curseActions = append(curseActions, v1_6.CurseGloballyOnlyOnChain(remoteSel))
			} else {
				curseActions = append(curseActions, v1_6.CurseLaneOnlyOnSource(remoteSel, tc.pf.chainSelectors[subjectDescription]))
			}

			_, err := v1_6.RMNCurseChangeset(envWithRMN.Env, v1_6.RMNCurseConfig{
				CurseActions: curseActions,
				Reason:       "test curse",
			})
			require.NoError(t, err)
		}

		cs, err := chState.RMNRemote.GetCursedSubjects(&bind.CallOpts{Context: ctx})
		require.NoError(t, err)
		t.Logf("Cursed subjects: %v", cs)
	}
}

func (tc rmnTestCase) callContractsToCurseAndRevokeCurse(ctx context.Context, eg *errgroup.Group, t *testing.T, onChainState stateview.CCIPOnChainState, envWithRMN testhelpers.DeployedEnv) {
	for _, remoteCfg := range tc.remoteChainsConfig {
		remoteSel := tc.pf.chainSelectors[remoteCfg.chainIdx]
		chState, ok := onChainState.EVMChainState(remoteSel)
		require.True(t, ok)
		_, ok = envWithRMN.Env.BlockChains.EVMChains()[remoteSel]
		require.True(t, ok)

		cursedSubjects := tc.revokedCursedSubjectsPerChain[remoteCfg.chainIdx]

		for subjectDescription, revokeAfter := range cursedSubjects {
			curseActions := make([]v1_6.CurseAction, 0)

			if subjectDescription == globalCurse {
				curseActions = append(curseActions, v1_6.CurseGloballyOnlyOnChain(remoteSel))
			} else {
				curseActions = append(curseActions, v1_6.CurseLaneOnlyOnSource(remoteSel, tc.pf.chainSelectors[subjectDescription]))
			}

			_, err := v1_6.RMNCurseChangeset(envWithRMN.Env, v1_6.RMNCurseConfig{
				CurseActions: curseActions,
				Reason:       "test curse",
			})
			require.NoError(t, err)

			eg.Go(func() error {
				<-time.NewTimer(revokeAfter).C
				t.Logf("revoking curse on subject %d (%d)", subjectDescription, subjectDescription)

				_, err := v1_6.RMNUncurseChangeset(envWithRMN.Env, v1_6.RMNCurseConfig{
					CurseActions: curseActions,
					Reason:       "test uncurse",
				})
				if err != nil {
					return err
				}
				return nil
			})
		}
		cs, err := chState.RMNRemote.GetCursedSubjects(&bind.CallOpts{Context: ctx})
		require.NoError(t, err)
		t.Logf("Cursed subjects: %v, %v", cs, remoteSel)
		eg.Go(func() error {
			<-time.NewTimer(time.Second * 10).C
			cs, err := chState.RMNRemote.GetCursedSubjects(&bind.CallOpts{Context: ctx})

			if err != nil {
				return err
			}

			t.Logf("Cursed subjects after revoking: %v, %v", cs, remoteSel)
			return nil
		})
	}
}

func (tc rmnTestCase) enableOracles(ctx context.Context, t *testing.T, envWithRMN testhelpers.DeployedEnv, nodeIDs []string) {
	for _, n := range nodeIDs {
		_, err := envWithRMN.Env.Offchain.EnableNode(ctx, &node.EnableNodeRequest{Id: n})
		require.NoError(t, err)
		t.Logf("node %s enabled", n)
	}
}

func configureAndPromoteRMNHome(
	t *testing.T,
	tc *rmnTestCase,
	envWithRMN testhelpers.DeployedEnv,
	rmnCluster devenv.RMNCluster,
) stateview.CCIPOnChainState {
	ctx := testcontext.Get(t)
	tc.populateFields(t, envWithRMN, rmnCluster)

	// Load on-chain state
	onChainState, err := stateview.LoadOnchainState(envWithRMN.Env)
	require.NoError(t, err)

	// Get the home chain state and the candidate/active digests
	homeChainState, ok := onChainState.EVMChainState(envWithRMN.HomeChainSel)
	require.True(t, ok)

	allDigests, err := homeChainState.RMNHome.GetConfigDigests(&bind.CallOpts{Context: ctx})
	require.NoError(t, err)
	t.Logf("RMNHome candidateDigest before setting new candidate: %x, activeDigest: %x",
		allDigests.CandidateConfigDigest[:], allDigests.ActiveConfigDigest[:])

	// Configure candidate using the populated test-case fields
	staticConfig := rmn_home.RMNHomeStaticConfig{Nodes: tc.pf.rmnHomeNodes, OffchainConfig: []byte{}}
	dynamicConfig := rmn_home.RMNHomeDynamicConfig{SourceChains: tc.pf.rmnHomeSourceChains, OffchainConfig: []byte{}}
	t.Logf("Setting RMNHome candidate with staticConfig: %+v, dynamicConfig: %+v, current candidateDigest: %x",
		staticConfig, dynamicConfig, allDigests.CandidateConfigDigest[:])

	candidateDigest, err := homeChainState.RMNHome.GetCandidateDigest(&bind.CallOpts{Context: ctx})
	require.NoError(t, err)

	_, err = v1_6.SetRMNHomeCandidateConfigChangeset(envWithRMN.Env, v1_6.SetRMNHomeCandidateConfig{
		HomeChainSelector: envWithRMN.HomeChainSel,
		RMNStaticConfig:   staticConfig,
		RMNDynamicConfig:  dynamicConfig,
		DigestToOverride:  candidateDigest,
	})
	require.NoError(t, err)

	candidateDigest, err = homeChainState.RMNHome.GetCandidateDigest(&bind.CallOpts{Context: ctx})
	require.NoError(t, err)
	t.Logf("RMNHome candidateDigest after setting new candidate: %x", candidateDigest[:])
	t.Logf("Promoting RMNHome candidate with candidateDigest: %x", candidateDigest[:])

	// Promote candidate
	_, err = v1_6.PromoteRMNHomeCandidateConfigChangeset(envWithRMN.Env, v1_6.PromoteRMNHomeCandidateConfig{
		HomeChainSelector: envWithRMN.HomeChainSel,
		DigestToPromote:   candidateDigest,
	})
	require.NoError(t, err)

	// Validate that candidate promotion is successful
	activeDigest, err := homeChainState.RMNHome.GetActiveDigest(&bind.CallOpts{Context: ctx})
	require.NoError(t, err)
	require.Equalf(t, candidateDigest, activeDigest,
		"active digest should be the same as the previously candidate digest after promotion, previous candidate: %x, active: %x",
		candidateDigest[:], activeDigest[:])

	// Configure remote chain settings
	rmnRemoteConfig := make(map[uint64]ccipops.RMNRemoteConfig)
	for _, remoteCfg := range tc.remoteChainsConfig {
		selector := tc.pf.chainSelectors[remoteCfg.chainIdx]
		if remoteCfg.f < 0 {
			t.Fatalf("remoteCfg.f is negative: %d", remoteCfg.f)
		}
		rmnRemoteConfig[selector] = ccipops.RMNRemoteConfig{
			F:       uint64(remoteCfg.f),
			Signers: tc.pf.rmnRemoteSigners,
		}
	}
	_, err = v1_6.SetRMNRemoteConfigChangeset(envWithRMN.Env, ccipseq.SetRMNRemoteConfig{
		RMNRemoteConfigs: rmnRemoteConfig,
	})
	require.NoError(t, err)

	return onChainState
}

func performReorgTest(t *testing.T, e testhelpers.DeployedEnv, l logging.Logger, dockerEnv *testsetups.DeployedLocalDevEnvironment, state stateview.CCIPOnChainState, nonBootstrapP2PIDs []string) (sourceSelector uint64, destSelector uint64) {
	// Chain setup
	allChains := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyEVM))
	require.GreaterOrEqual(t, len(allChains), 2)
	sourceSelector = allChains[0]
	destSelector = allChains[1]

	// Build RPC map and get clients
	chainSelToRPCURL := buildChainSelectorToRPCURLMap(t, dockerEnv)
	sourceClient := ctf_client.NewRPCClient(chainSelToRPCURL[sourceSelector], nil)

	// Setup CCIP lane
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceSelector, destSelector, false)
	waitForLogPollerFilters(l)

	// Send initial message
	msgBeforeReorg := sendCCIPMessage(t, e.Env, state, sourceSelector, destSelector, l)

	// Wait and perform reorg
	minBlock := msgBeforeReorg.Raw.BlockNumber + lessThanFinalityReorgDepth - 1
	waitForBlockNumber(t, sourceClient, minBlock, 1*time.Minute, 500*time.Millisecond, l)
	performReorg(t, sourceClient, lessThanFinalityReorgDepth, l)

	// Verify message consistency
	msgAfterReorg := sendCCIPMessage(t, e.Env, state, sourceSelector, destSelector, l)
	require.Equal(t, msgBeforeReorg.Message.Header.SequenceNumber, msgAfterReorg.Message.Header.SequenceNumber)
	require.Equal(t, msgBeforeReorg.Message.Header.MessageId, msgAfterReorg.Message.Header.MessageId)

	// Check node health
	nodeAPIs := dockerEnv.GetCLClusterTestEnv().ClCluster.NodeAPIs()
	checkFinalityViolations(
		t,
		nodeAPIs,
		nonBootstrapP2PIDs,
		getHeadTrackerService(t, sourceSelector),
		getLogPollerService(t, sourceSelector),
		l,
		0,              // no nodes reporting finality violation
		1*time.Minute,  // timeout
		10*time.Second, // interval
	)

	return sourceSelector, destSelector
}

func Test_CCIPReorg_BelowFinality_OnSource_WithRMN(t *testing.T) {
	tc := rmnTestCase{
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
			{id: 2, isSigner: true, observedChainIdxs: []int{chain0, chain1}}, // one rmn node is down
		},
	}
	e, l, dockerEnv, nonBootstrapP2PIDs, state, rmnCluster := setupReorgTest(t,
		testhelpers.WithExtraConfigTomls([]string{"Test_CCIPReorg_BelowFinality_OnSource_WithRMN.toml"}),
		testhelpers.WithRMNEnabled(len(tc.rmnNodes)),
		testhelpers.WithRMNConfDepth(20),
	)

	configureAndPromoteRMNHome(t, &tc, e, rmnCluster)

	e.RmnEnabledSourceChains = make(map[uint64]bool)
	for chainIdx := range tc.homeChainConfig.f {
		chainSel := tc.pf.chainSelectors[chainIdx]
		e.RmnEnabledSourceChains[chainSel] = true
	}

	sourceSelector, destSelector := performReorgTest(t, e, l, dockerEnv, state, nonBootstrapP2PIDs)

	_, err := testhelpers.ConfirmCommitWithExpectedSeqNumRange(
		t,
		sourceSelector,
		e.Env.BlockChains.EVMChains()[destSelector],
		state.MustGetEVMChainState(destSelector).OffRamp,
		nil, // startBlock
		ccipocr3.NewSeqNumRange(1, 1),
		false, // enforceSingleCommit
	)
	require.NoError(t, err)
}

func Test_CCIPReorg_BelowFinality_OnSource_WithRMN_Recover(t *testing.T) {
	tc := rmnTestCase{
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
			{id: 2, isSigner: true, observedChainIdxs: []int{chain0, chain1}}, // one rmn node is down
		},
	}
	e, l, dockerEnv, nonBootstrapP2PIDs, state, rmnCluster := setupReorgTest(t,
		testhelpers.WithExtraConfigTomls([]string{"Test_CCIPReorg_BelowFinality_OnSource_WithRMN.toml"}),
		testhelpers.WithRMNEnabled(len(tc.rmnNodes)),
		testhelpers.WithRMNConfDepth(0),
	)

	configureAndPromoteRMNHome(t, &tc, e, rmnCluster)

	e.RmnEnabledSourceChains = make(map[uint64]bool)
	fmt.Printf("Setup RMN enabled")
	for chainIdx := range tc.homeChainConfig.f {
		chainSel := tc.pf.chainSelectors[chainIdx]
		e.RmnEnabledSourceChains[chainSel] = true
		fmt.Printf("Setup RMN enabled for chain %d", chainSel)
	}

	sourceSelector, destSelector := performReorgTest(t, e, l, dockerEnv, state, nonBootstrapP2PIDs)

	err := rmnCluster.Restart(t.Context())
	require.NoError(t, err)

	_, err = testhelpers.ConfirmCommitWithExpectedSeqNumRange(
		t,
		sourceSelector,
		e.Env.BlockChains.EVMChains()[destSelector],
		state.MustGetEVMChainState(destSelector).OffRamp,
		nil, // startBlock
		ccipocr3.NewSeqNumRange(1, 1),
		false, // enforceSingleCommit
	)
	require.NoError(t, err)
}

func Test_CCIPReorg_BelowFinality_OnSource_WithRMN_Block(t *testing.T) {
	tc := rmnTestCase{
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
			{id: 2, isSigner: true, observedChainIdxs: []int{chain0, chain1}}, // one rmn node is down
		},
	}
	e, l, dockerEnv, nonBootstrapP2PIDs, state, rmnCluster := setupReorgTest(t,
		testhelpers.WithExtraConfigTomls([]string{"Test_CCIPReorg_BelowFinality_OnSource_WithRMN.toml"}),
		testhelpers.WithRMNEnabled(len(tc.rmnNodes)),
		testhelpers.WithRMNConfDepth(0),
	)

	configureAndPromoteRMNHome(t, &tc, e, rmnCluster)

	e.RmnEnabledSourceChains = make(map[uint64]bool)
	fmt.Printf("Setup RMN enabled")
	for chainIdx := range tc.homeChainConfig.f {
		chainSel := tc.pf.chainSelectors[chainIdx]
		e.RmnEnabledSourceChains[chainSel] = true
		fmt.Printf("Setup RMN enabled for chain %d", chainSel)
	}

	sourceSelector, destSelector := performReorgTest(t, e, l, dockerEnv, state, nonBootstrapP2PIDs)

	commitReportReceived := make(chan struct{})
	commitReportError := make(chan error)
	// Verify commit
	go func() {
		_, err := testhelpers.ConfirmCommitWithExpectedSeqNumRange(
			t,
			sourceSelector,
			e.Env.BlockChains.EVMChains()[destSelector],
			state.MustGetEVMChainState(destSelector).OffRamp,
			nil, // startBlock
			ccipocr3.NewSeqNumRange(1, 1),
			false, // enforceSingleCommit
		)

		if err != nil {
			commitReportError <- err
		} else {
			commitReportReceived <- struct{}{}
		}
	}()

	tim := time.NewTimer(15 * time.Second)
	t.Logf("waiting for 15s before asserting that commit report was not received")
	select {
	case err := <-commitReportError:
		t.Errorf("Error while confirming commit: %v", err)
	case <-commitReportReceived:
		t.Errorf("Commit report was received while it was not expected")
		return
	case <-tim.C:
		return
	}
}
