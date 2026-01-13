package testhelpers

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	chainsel "github.com/smartcontractkit/chain-selectors"

	commonutils "github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

// SourceDestPair is represents a pair of source and destination chain selectors.
// Use this as a key in maps that need to identify sequence numbers, nonces, or
// other things that require identification.
type SourceDestPair struct {
	SourceChainSelector uint64
	DestChainSelector   uint64
}

func ToSeqRangeMap(seqNrs map[SourceDestPair]uint64) map[SourceDestPair]ccipocr3.SeqNumRange {
	seqRangeMap := make(map[SourceDestPair]ccipocr3.SeqNumRange)
	for sourceDest, seqNr := range seqNrs {
		seqRangeMap[sourceDest] = ccipocr3.SeqNumRange{
			ccipocr3.SeqNum(seqNr), ccipocr3.SeqNum(seqNr),
		}
	}
	return seqRangeMap
}

// ConfirmCommitForAllWithExpectedSeqNums waits for all chains in the environment to commit the given expectedSeqNums.
// expectedSeqNums is a map that maps a (source, dest) selector pair to the expected sequence number
// to confirm the commit for.
// startBlocks is a map of destination chain selector to start block number to start watching from.
// If startBlocks is nil, it will start watching from the latest block.
func ConfirmCommitForAllWithExpectedSeqNums(
	t *testing.T,
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	expectedSeqNums map[SourceDestPair]ccipocr3.SeqNumRange,
	startBlocks map[uint64]*uint64,
) {
	var wg errgroup.Group
	for sourceDest, expectedSeqNum := range expectedSeqNums {
		srcChain := sourceDest.SourceChainSelector
		dstChain := sourceDest.DestChainSelector
		if expectedSeqNum.Start() == 0 {
			continue
		}
		wg.Go(func() error {
			var startBlock *uint64
			if startBlocks != nil {
				startBlock = startBlocks[dstChain]
			}

			family, err := chainsel.GetSelectorFamily(dstChain)
			if err != nil {
				return err
			}
			switch family {
			case chainsel.FamilyEVM:
				return commonutils.JustError(ConfirmCommitWithExpectedSeqNumRange(
					t,
					srcChain,
					e.BlockChains.EVMChains()[dstChain],
					state.MustGetEVMChainState(dstChain).OffRamp,
					startBlock,
					expectedSeqNum,
					true,
				))
			case chainsel.FamilySolana:
				var startSlot uint64
				if startBlock != nil {
					startSlot = *startBlock
				}
				return commonutils.JustError(confirmCommitWithExpectedSeqNumRangeSol(
					t,
					srcChain,
					e.BlockChains.SolanaChains()[dstChain],
					state.SolChains[dstChain].OffRamp,
					startSlot,
					expectedSeqNum,
					true,
				))
			case chainsel.FamilySui:
				return commonutils.JustError(confirmCommitWithExpectedSeqNumRangeSui(
					t,
					srcChain,
					e.BlockChains.SuiChains()[dstChain],
					state.SuiChains[dstChain].OffRampAddress,
					startBlock,
					expectedSeqNum,
					true,
				))
			case chainsel.FamilyAptos:
				return commonutils.JustError(confirmCommitWithExpectedSeqNumRangeAptos(
					t,
					srcChain,
					e.BlockChains.AptosChains()[dstChain],
					state.AptosChains[dstChain].CCIPAddress,
					startBlock,
					expectedSeqNum,
					true,
				))
			case chainsel.FamilyTon:
				return commonutils.JustError(confirmCommitWithExpectedSeqNumRangeTON(
					t,
					srcChain,
					e.BlockChains.TonChains()[dstChain],
					state.TonChains[dstChain].OffRamp,
					expectedSeqNum,
				))
			default:
				return fmt.Errorf("unsupported chain family; %v", family)
			}
		})
	}

	done := make(chan struct{})
	go func() {
		assert.NoError(t, wg.Wait())
		close(done)
	}()

	require.Eventually(t, func() bool {
		select {
		case <-done:
			return true
		default:
			return false
		}
	},
		tests.WaitTimeout(t),
		2*time.Second,
		"all commitments did not confirm",
	)
}

type CommitReportTracker struct {
	seenMessages map[uint64]map[uint64]bool
}

func NewCommitReportTracker(sourceChainSelector uint64, seqNrs ccipocr3.SeqNumRange) CommitReportTracker {
	seenMessages := make(map[uint64]map[uint64]bool)
	seenMessages[sourceChainSelector] = make(map[uint64]bool)

	for i := seqNrs.Start(); i <= seqNrs.End(); i++ {
		seenMessages[sourceChainSelector][uint64(i)] = false
	}
	return CommitReportTracker{seenMessages: seenMessages}
}

func (c *CommitReportTracker) visitCommitReport(sourceChainSelector uint64, minSeqNr uint64, maxSeqNr uint64) {
	if _, ok := c.seenMessages[sourceChainSelector]; !ok {
		return
	}

	for i := minSeqNr; i <= maxSeqNr; i++ {
		c.seenMessages[sourceChainSelector][i] = true
	}
}

func (c *CommitReportTracker) allCommitted(sourceChainSelector uint64) bool {
	for _, v := range c.seenMessages[sourceChainSelector] {
		if !v {
			return false
		}
	}
	return true
}

// ConfirmMultipleCommits waits for multiple ccipocr3.SeqNumRange to be committed by the Offramp.
// Waiting is done in parallel per every sourceChain/destChain (lane) passed as argument.
func ConfirmMultipleCommits(
	t *testing.T,
	env cldf.Environment,
	state stateview.CCIPOnChainState,
	startBlocks map[uint64]*uint64,
	enforceSingleCommit bool,
	expectedSeqNums map[SourceDestPair]ccipocr3.SeqNumRange,
) error {
	errGrp := &errgroup.Group{}

	for sourceDest, seqRange := range expectedSeqNums {
		srcChain := sourceDest.SourceChainSelector
		destChain := sourceDest.DestChainSelector

		errGrp.Go(func() error {
			family, err := chainsel.GetSelectorFamily(destChain)
			if err != nil {
				return err
			}
			switch family {
			case chainsel.FamilyEVM:
				_, err := ConfirmCommitWithExpectedSeqNumRange(
					t,
					srcChain,
					env.BlockChains.EVMChains()[destChain],
					state.MustGetEVMChainState(destChain).OffRamp,
					startBlocks[destChain],
					seqRange,
					enforceSingleCommit,
				)
				return err
			case chainsel.FamilySolana:
				var startSlot uint64
				if startBlocks[destChain] != nil {
					startSlot = *startBlocks[destChain]
				}
				_, err := confirmCommitWithExpectedSeqNumRangeSol(
					t,
					srcChain,
					env.BlockChains.SolanaChains()[destChain],
					state.SolChains[destChain].OffRamp,
					startSlot,
					seqRange,
					enforceSingleCommit,
				)
				return err
			case chainsel.FamilySui:
				_, err := confirmCommitWithExpectedSeqNumRangeSui(
					t,
					srcChain,
					env.BlockChains.SuiChains()[destChain],
					state.SuiChains[destChain].OffRampAddress,
					startBlocks[destChain],
					seqRange,
					enforceSingleCommit,
				)
				return err
			case chainsel.FamilyAptos:
				_, err := confirmCommitWithExpectedSeqNumRangeAptos(
					t,
					srcChain,
					env.BlockChains.AptosChains()[destChain],
					state.AptosChains[destChain].CCIPAddress,
					startBlocks[destChain],
					seqRange,
					enforceSingleCommit,
				)
				return err
			default:
				return fmt.Errorf("unsupported chain family; %v", family)
			}
		})
	}

	return errGrp.Wait()
}

// ConfirmExecWithSeqNrsForAll waits for all chains in the environment to execute the given expectedSeqNums.
// If successful, it returns a map that maps the SourceDestPair to the expected sequence number
// to its execution state.
// expectedSeqNums is a map of SourceDestPair to a slice of expected sequence numbers to be executed.
// startBlocks is a map of destination chain selector to start block number to start watching from.
// If startBlocks is nil, it will start watching from the latest block.
func ConfirmExecWithSeqNrsForAll(
	t *testing.T,
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	expectedSeqNums map[SourceDestPair][]uint64,
	startBlocks map[uint64]*uint64,
) (executionStates map[SourceDestPair]map[uint64]int) {
	var (
		wg errgroup.Group
		mx sync.Mutex
	)
	executionStates = make(map[SourceDestPair]map[uint64]int)
	for sourceDest, seqRange := range expectedSeqNums {
		srcChain := sourceDest.SourceChainSelector
		dstChain := sourceDest.DestChainSelector

		var startBlock *uint64
		if startBlocks != nil {
			startBlock = startBlocks[dstChain]
		}

		wg.Go(func() error {
			family, err := chainsel.GetSelectorFamily(dstChain)
			if err != nil {
				return err
			}

			var innerExecutionStates map[uint64]int
			switch family {
			case chainsel.FamilyEVM:
				innerExecutionStates, err = ConfirmExecWithSeqNrs(
					t,
					srcChain,
					e.BlockChains.EVMChains()[dstChain],
					state.MustGetEVMChainState(dstChain).OffRamp,
					startBlock,
					seqRange,
				)
				if err != nil {
					return err
				}
			case chainsel.FamilySolana:
				var startSlot uint64
				if startBlock != nil {
					startSlot = *startBlock
				}
				innerExecutionStates, err = confirmExecWithSeqNrsSol(
					t,
					srcChain,
					e.BlockChains.SolanaChains()[dstChain],
					state.SolChains[dstChain].OffRamp,
					startSlot,
					seqRange,
				)
				if err != nil {
					return err
				}
			case chainsel.FamilyAptos:
				innerExecutionStates, err = confirmExecWithExpectedSeqNrsAptos(
					t,
					srcChain,
					e.BlockChains.AptosChains()[dstChain],
					state.AptosChains[dstChain].CCIPAddress,
					startBlock,
					seqRange,
				)
				if err != nil {
					return err
				}
			case chainsel.FamilySui:
				innerExecutionStates, err = confirmExecWithExpectedSeqNrsSui(
					t,
					srcChain,
					e.BlockChains.SuiChains()[dstChain],
					state.SuiChains[dstChain].OffRampAddress,
					startBlock,
					seqRange,
				)
				if err != nil {
					return err
				}
			case chainsel.FamilyTon:
				innerExecutionStates, err = confirmExecWithExpectedSeqNrsTON(
					t,
					srcChain,
					e.BlockChains.TonChains()[dstChain],
					state.TonChains[dstChain].OffRamp,
					startBlock,
					seqRange,
				)
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("unsupported chain family; %v", family)
			}

			mx.Lock()
			executionStates[sourceDest] = innerExecutionStates
			mx.Unlock()

			return nil
		})
	}

	require.NoError(t, wg.Wait())
	return executionStates
}

func ConfirmNoExecSuccessConsistentlyWithSeqNr(
	t *testing.T,
	sourceSelector uint64,
	dest cldf_evm.Chain,
	offRamp offramp.OffRampInterface,
	expectedSeqNr uint64,
	timeout time.Duration,
) {
	RequireConsistently(t, func() bool {
		scc, executionState := getExecutionState(t, sourceSelector, offRamp, expectedSeqNr)
		t.Logf("Waiting for ExecutionStateChanged on chain %d (offramp %s) from chain %d with expected sequence number %d, current onchain minSeqNr: %d, execution state: %s",
			dest.Selector, offRamp.Address().String(), sourceSelector, expectedSeqNr, scc.MinSeqNr, executionStateToString(executionState))
		if executionState == EXECUTION_STATE_SUCCESS {
			t.Logf("Observed %s execution state on chain %d (offramp %s) from chain %d with expected sequence number %d",
				executionStateToString(executionState), dest.Selector, offRamp.Address().String(), sourceSelector, expectedSeqNr)
			return false
		}
		return true
	}, timeout, 3*time.Second, "Expected no execution success on chain %d (offramp %s) from chain %d with expected sequence number %d", dest.Selector, offRamp.Address().String(), sourceSelector, expectedSeqNr)
}

func getExecutionState(t *testing.T, sourceSelector uint64, offRamp offramp.OffRampInterface, expectedSeqNr uint64) (offramp.OffRampSourceChainConfig, uint8) {
	scc, err := offRamp.GetSourceChainConfig(nil, sourceSelector)
	require.NoError(t, err)
	executionState, err := offRamp.GetExecutionState(nil, sourceSelector, expectedSeqNr)
	require.NoError(t, err)
	return scc, executionState
}

func RequireConsistently(t *testing.T, condition func() bool, duration time.Duration, tick time.Duration, msgAndArgs ...any) {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	tickTimer := time.NewTicker(tick)
	defer tickTimer.Stop()
	for {
		select {
		case <-tickTimer.C:
			if !condition() {
				require.FailNow(t, "Condition failed", msgAndArgs...)
			}
		case <-timer.C:
			return
		}
	}
}

func SeqNumberRangeToSlice(seqRanges map[SourceDestPair]ccipocr3.SeqNumRange) map[SourceDestPair][]uint64 {
	flatten := make(map[SourceDestPair][]uint64)

	for srcDst, seqRange := range seqRanges {
		if _, ok := flatten[srcDst]; !ok {
			flatten[srcDst] = make([]uint64, 0, seqRange.End()-seqRange.Start()+1)
		}

		for i := seqRange.Start(); i <= seqRange.End(); i++ {
			flatten[srcDst] = append(flatten[srcDst], uint64(i))
		}
	}

	return flatten
}

const (
	EXECUTION_STATE_UNTOUCHED  = 0
	EXECUTION_STATE_INPROGRESS = 1
	EXECUTION_STATE_SUCCESS    = 2
	EXECUTION_STATE_FAILURE    = 3
)

func executionStateToString(state uint8) string {
	switch state {
	case EXECUTION_STATE_UNTOUCHED:
		return "UNTOUCHED"
	case EXECUTION_STATE_INPROGRESS:
		return "IN_PROGRESS"
	case EXECUTION_STATE_SUCCESS:
		return "SUCCESS"
	case EXECUTION_STATE_FAILURE:
		return "FAILURE"
	default:
		return "UNKNOWN"
	}
}

func AssertEqualFeeConfig(t *testing.T, want, have fee_quoter.FeeQuoterDestChainConfig) {
	assert.Equal(t, want.DestGasOverhead, have.DestGasOverhead)
	assert.Equal(t, want.IsEnabled, have.IsEnabled)
	assert.Equal(t, want.ChainFamilySelector, have.ChainFamilySelector)
	assert.Equal(t, want.DefaultTokenDestGasOverhead, have.DefaultTokenDestGasOverhead)
	assert.Equal(t, want.DefaultTokenFeeUSDCents, have.DefaultTokenFeeUSDCents)
	assert.Equal(t, want.DefaultTxGasLimit, have.DefaultTxGasLimit)
	assert.Equal(t, want.DestGasPerPayloadByteBase, have.DestGasPerPayloadByteBase)
	assert.Equal(t, want.DestGasPerPayloadByteHigh, have.DestGasPerPayloadByteHigh)
	assert.Equal(t, want.DestGasPerPayloadByteThreshold, have.DestGasPerPayloadByteThreshold)
	assert.Equal(t, want.DestGasPerDataAvailabilityByte, have.DestGasPerDataAvailabilityByte)
	assert.Equal(t, want.DestDataAvailabilityMultiplierBps, have.DestDataAvailabilityMultiplierBps)
	assert.Equal(t, want.DestDataAvailabilityOverheadGas, have.DestDataAvailabilityOverheadGas)
	assert.Equal(t, want.MaxDataBytes, have.MaxDataBytes)
	assert.Equal(t, want.MaxNumberOfTokensPerMsg, have.MaxNumberOfTokensPerMsg)
	assert.Equal(t, want.MaxPerMsgGasLimit, have.MaxPerMsgGasLimit)
}

// AssertTimelockOwnership asserts that the ownership of the contracts has been transferred
// to the appropriate timelock contract on each chain.
func AssertTimelockOwnership(
	t *testing.T,
	e DeployedEnv,
	chains []uint64,
	state stateview.CCIPOnChainState,
	withTestRouterTransfer bool,
) {
	// check that the ownership has been transferred correctly
	for _, chain := range chains {
		allContracts := []common.Address{
			state.MustGetEVMChainState(chain).OnRamp.Address(),
			state.MustGetEVMChainState(chain).OffRamp.Address(),
			state.MustGetEVMChainState(chain).FeeQuoter.Address(),
			state.MustGetEVMChainState(chain).NonceManager.Address(),
			state.MustGetEVMChainState(chain).RMNRemote.Address(),
			state.MustGetEVMChainState(chain).Router.Address(),
			state.MustGetEVMChainState(chain).TokenAdminRegistry.Address(),
			state.MustGetEVMChainState(chain).RMNProxy.Address(),
		}
		if withTestRouterTransfer {
			allContracts = append(allContracts, state.MustGetEVMChainState(chain).TestRouter.Address())
		}
		for _, contract := range allContracts {
			owner, _, err := commonchangeset.LoadOwnableContract(contract, e.Env.BlockChains.EVMChains()[chain].Client)
			require.NoError(t, err)
			require.Equal(t, state.MustGetEVMChainState(chain).Timelock.Address(), owner)
		}
	}

	// check home chain contracts ownership
	homeChainTimelockAddress := state.MustGetEVMChainState(e.HomeChainSel).Timelock.Address()
	for _, contract := range []common.Address{
		state.MustGetEVMChainState(e.HomeChainSel).CapabilityRegistry.Address(),
		state.MustGetEVMChainState(e.HomeChainSel).CCIPHome.Address(),
		state.MustGetEVMChainState(e.HomeChainSel).RMNHome.Address(),
	} {
		owner, _, err := commonchangeset.LoadOwnableContract(contract, e.Env.BlockChains.EVMChains()[e.HomeChainSel].Client)
		require.NoError(t, err)
		require.Equal(t, homeChainTimelockAddress, owner)
	}
}
