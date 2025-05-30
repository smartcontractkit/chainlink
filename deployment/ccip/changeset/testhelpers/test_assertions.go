package testhelpers

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

	solconfig "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solccip "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/ccip"
	solcommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"

	chainsel "github.com/smartcontractkit/chain-selectors"

	commonutils "github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
)

func ConfirmGasPriceUpdatedForAll(
	t *testing.T,
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	startBlocks map[uint64]*uint64,
	gasPrice *big.Int,
) {
	var wg errgroup.Group
	evmChains := e.BlockChains.EVMChains()
	for src, srcChain := range evmChains {
		for dest, dstChain := range evmChains {
			if src == dest {
				continue
			}
			srcChain := srcChain
			dstChain := dstChain
			wg.Go(func() error {
				var startBlock *uint64
				if startBlocks != nil {
					startBlock = startBlocks[srcChain.Selector]
				}
				return ConfirmGasPriceUpdated(
					t,
					dstChain,
					state.MustGetEVMChainState(srcChain.Selector).FeeQuoter,
					*startBlock,
					gasPrice,
				)
			})
		}
	}
	require.NoError(t, wg.Wait())
}

func ConfirmGasPriceUpdated(
	t *testing.T,
	dest cldf_evm.Chain,
	srcFeeQuoter *fee_quoter.FeeQuoter,
	startBlock uint64,
	gasPrice *big.Int,
) error {
	it, err := srcFeeQuoter.FilterUsdPerUnitGasUpdated(&bind.FilterOpts{
		Context: context.Background(),
		Start:   startBlock,
	}, []uint64{dest.Selector})

	require.NoError(t, err)
	require.Truef(t, it.Next(), "No gas price update event found on chain %d, fee quoter %s",
		dest.Selector, srcFeeQuoter.Address().String())
	require.NotEqualf(t, gasPrice, it.Event.Value, "Gas price not updated on chain %d, fee quoter %s",
		dest.Selector, srcFeeQuoter.Address().String())
	return nil
}

func ConfirmTokenPriceUpdatedForAll(
	t *testing.T,
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	startBlocks map[uint64]*uint64,
	linkPrice *big.Int,
	wethPrice *big.Int,
) {
	var wg errgroup.Group
	for _, chain := range e.BlockChains.EVMChains() {
		chain := chain
		wg.Go(func() error {
			var startBlock *uint64
			if startBlocks != nil {
				startBlock = startBlocks[chain.Selector]
			}
			linkAddress := state.MustGetEVMChainState(chain.Selector).LinkToken.Address()
			wethAddress := state.MustGetEVMChainState(chain.Selector).Weth9.Address()
			tokenToPrice := make(map[common.Address]*big.Int)
			tokenToPrice[linkAddress] = linkPrice
			tokenToPrice[wethAddress] = wethPrice
			return ConfirmTokenPriceUpdated(
				t,
				chain,
				state.MustGetEVMChainState(chain.Selector).FeeQuoter,
				*startBlock,
				tokenToPrice,
			)
		})
	}
	require.NoError(t, wg.Wait())
}

func ConfirmTokenPriceUpdated(
	t *testing.T,
	chain cldf_evm.Chain,
	feeQuoter *fee_quoter.FeeQuoter,
	startBlock uint64,
	tokenToInitialPrice map[common.Address]*big.Int,
) error {
	tokens := make([]common.Address, 0, len(tokenToInitialPrice))
	for token := range tokenToInitialPrice {
		tokens = append(tokens, token)
	}
	it, err := feeQuoter.FilterUsdPerTokenUpdated(&bind.FilterOpts{
		Context: context.Background(),
		Start:   startBlock,
	}, tokens)
	require.NoError(t, err)
	for it.Next() {
		token := it.Event.Token
		initialValue, ok := tokenToInitialPrice[token]
		if ok {
			require.Contains(t, tokens, token)
			// Initial Value should be changed
			require.NotEqual(t, initialValue, it.Event.Value)
		}

		// Remove the token from the map until we assert all tokens are updated
		delete(tokenToInitialPrice, token)
		if len(tokenToInitialPrice) == 0 {
			return nil
		}
	}

	if len(tokenToInitialPrice) > 0 {
		return fmt.Errorf("not all tokens updated on chain  %d", chain.Selector)
	}

	return nil
}

// SourceDestPair is represents a pair of source and destination chain selectors.
// Use this as a key in maps that need to identify sequence numbers, nonces, or
// other things that require identification.
type SourceDestPair struct {
	SourceChainSelector uint64
	DestChainSelector   uint64
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
	expectedSeqNums map[SourceDestPair]uint64,
	startBlocks map[uint64]*uint64,
) {
	var wg errgroup.Group
	for sourceDest, expectedSeqNum := range expectedSeqNums {
		srcChain := sourceDest.SourceChainSelector
		dstChain := sourceDest.DestChainSelector
		if expectedSeqNum == 0 {
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
					ccipocr3.SeqNumRange{
						ccipocr3.SeqNum(expectedSeqNum),
						ccipocr3.SeqNum(expectedSeqNum),
					},
					true,
				))
			case chainsel.FamilySolana:
				var startSlot uint64
				if startBlock != nil {
					startSlot = *startBlock
				}
				return commonutils.JustError(ConfirmCommitWithExpectedSeqNumRangeSol(
					t,
					srcChain,
					e.BlockChains.SolanaChains()[dstChain],
					state.SolChains[dstChain].OffRamp,
					startSlot,
					ccipocr3.SeqNumRange{
						ccipocr3.SeqNum(expectedSeqNum),
						ccipocr3.SeqNum(expectedSeqNum),
					},
					true,
				))
			default:
				return fmt.Errorf("unsupported chain family; %v", family)
			}
		})
	}

	done := make(chan struct{})
	go func() {
		require.NoError(t, wg.Wait())
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

func (c *CommitReportTracker) allCommited(sourceChainSelector uint64) bool {
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
		seqRange := seqRange
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
				_, err := ConfirmCommitWithExpectedSeqNumRangeSol(
					t,
					srcChain,
					env.BlockChains.SolanaChains()[destChain],
					state.SolChains[destChain].OffRamp,
					startSlot,
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

// ConfirmCommitWithExpectedSeqNumRange waits for a commit report on the destination chain with the expected sequence number range.
// startBlock is the block number to start watching from.
// If startBlock is nil, it will start watching from the latest block.
func ConfirmCommitWithExpectedSeqNumRange(
	t *testing.T,
	srcSelector uint64,
	dest cldf_evm.Chain,
	offRamp offramp.OffRampInterface,
	startBlock *uint64,
	expectedSeqNumRange ccipocr3.SeqNumRange,
	enforceSingleCommit bool,
) (*offramp.OffRampCommitReportAccepted, error) {
	sink := make(chan *offramp.OffRampCommitReportAccepted)
	subscription, err := offRamp.WatchCommitReportAccepted(&bind.WatchOpts{
		Context: context.Background(),
		Start:   startBlock,
	}, sink)
	if err != nil {
		return nil, fmt.Errorf("error to subscribe CommitReportAccepted : %w", err)
	}

	seenMessages := NewCommitReportTracker(srcSelector, expectedSeqNumRange)

	verifyCommitReport := func(report *offramp.OffRampCommitReportAccepted) bool {
		processRoots := func(roots []offramp.InternalMerkleRoot) bool {
			for _, mr := range roots {
				t.Logf(
					"Received commit report for [%d, %d] on selector %d from source selector %d expected seq nr range %s, token prices: %v",
					mr.MinSeqNr, mr.MaxSeqNr, dest.Selector, srcSelector, expectedSeqNumRange.String(), report.PriceUpdates.TokenPriceUpdates,
				)
				seenMessages.visitCommitReport(srcSelector, mr.MinSeqNr, mr.MaxSeqNr)

				if mr.SourceChainSelector == srcSelector &&
					uint64(expectedSeqNumRange.Start()) >= mr.MinSeqNr &&
					uint64(expectedSeqNumRange.End()) <= mr.MaxSeqNr {
					t.Logf(
						"All sequence numbers committed in a single report [%d, %d]",
						expectedSeqNumRange.Start(), expectedSeqNumRange.End(),
					)
					return true
				}

				if !enforceSingleCommit && seenMessages.allCommited(srcSelector) {
					t.Logf(
						"All sequence numbers already committed from range [%d, %d]",
						expectedSeqNumRange.Start(), expectedSeqNumRange.End(),
					)
					return true
				}
			}
			return false
		}

		return processRoots(report.BlessedMerkleRoots) || processRoots(report.UnblessedMerkleRoots)
	}

	defer subscription.Unsubscribe()
	timeout := time.NewTimer(tests.WaitTimeout(t))
	defer timeout.Stop()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-t.Context().Done():
			return nil, nil
		case <-ticker.C:
			t.Logf("Waiting for commit report on chain selector %d from source selector %d expected seq nr range %s",
				dest.Selector, srcSelector, expectedSeqNumRange.String())

			// Need to do this because the subscription sometimes fails to get the event.
			iter, err := offRamp.FilterCommitReportAccepted(&bind.FilterOpts{
				Context: t.Context(),
			})
			// In some test case the test ends while the filter is still running resulting in a context.Canceled error.
			if err != nil && !errors.Is(err, context.Canceled) {
				require.NoError(t, err)
			}
			for iter.Next() {
				event := iter.Event
				verified := verifyCommitReport(event)
				if verified {
					return event, nil
				}
			}
		case subErr := <-subscription.Err():
			return nil, fmt.Errorf("subscription error: %w", subErr)
		case <-timeout.C:
			return nil, fmt.Errorf("timed out after waiting for commit report on chain selector %d from source selector %d expected seq nr range %s",
				dest.Selector, srcSelector, expectedSeqNumRange.String())
		case report := <-sink:
			verified := verifyCommitReport(report)
			if verified {
				return report, nil
			}
		}
	}
}

// Scan for events referencing address
func SolEventEmitter[T any](
	t *testing.T,
	client *solrpc.Client,
	address solana.PublicKey,
	eventType string,
	startSlot uint64,
	done chan any,
) (<-chan T, <-chan error) {
	ch := make(chan T)
	errorCh := make(chan error)
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		var until solana.Signature
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				// Scan for transactions referencing the address
				ctx := context.Background()
				txSigs, err := client.GetSignaturesForAddressWithOpts(
					ctx,
					address,
					&solrpc.GetSignaturesForAddressOpts{
						Commitment: solrpc.CommitmentConfirmed,
						Until:      until,
					},
				)
				if err != nil {
					errorCh <- err
					return
				}

				if len(txSigs) == 0 {
					continue
				}

				// values are returned ordered newest to oldest, so we replay them backwards
				for _, txSig := range slices.Backward(txSigs) {
					if txSig.Err != nil {
						// We're not interested in failed transactions.
						continue
					}
					if txSig.Slot < startSlot {
						// Skip any signatures that are before the starting slot
						continue
					}
					v := uint64(0) // v0 = latest, supports address table lookups
					tx, err := client.GetTransaction(
						ctx,
						txSig.Signature,
						&solrpc.GetTransactionOpts{
							Commitment:                     solrpc.CommitmentConfirmed,
							Encoding:                       solana.EncodingBase64,
							MaxSupportedTransactionVersion: &v,
						},
					)
					if err != nil {
						errorCh <- err
						return
					}

					events, err := solcommon.ParseMultipleEvents[T](tx.Meta.LogMessages, eventType, solconfig.PrintEvents)
					if err != nil && strings.Contains(err.Error(), "event not found") {
						continue
					}
					if err != nil {
						errorCh <- err
						return
					}

					for _, event := range events {
						select {
						case ch <- event:
						case <-done:
							return
						}
					}
				}
				// next scan should stop at the newest signature we've received
				until = txSigs[0].Signature
			}
		}
	}()

	return ch, errorCh
}

func ConfirmCommitWithExpectedSeqNumRangeSol(
	t *testing.T,
	srcSelector uint64,
	dest cldf_solana.Chain,
	offrampAddress solana.PublicKey,
	startSlot uint64,
	expectedSeqNumRange ccipocr3.SeqNumRange,
	enforceSingleCommit bool,
) (bool, error) {
	seenMessages := NewCommitReportTracker(srcSelector, expectedSeqNumRange)

	done := make(chan any)
	defer close(done)
	sink, errCh := SolEventEmitter[solccip.EventCommitReportAccepted](t, dest.Client, offrampAddress, "CommitReportAccepted", startSlot, done)

	timeout := time.NewTimer(tests.WaitTimeout(t))
	defer timeout.Stop()

	for {
		select {
		case commitEvent := <-sink:
			// if merkle root is zero, it only contains price updates
			if commitEvent.Report == nil {
				t.Logf("Skipping CommitReportAccepted with only price updates")
				continue
			}
			require.Equal(t, srcSelector, commitEvent.Report.SourceChainSelector)

			// TODO: this logic is duplicated with verifyCommitReport, share
			mr := commitEvent.Report
			seenMessages.visitCommitReport(mr.SourceChainSelector, mr.MinSeqNr, mr.MaxSeqNr)
			if mr.SourceChainSelector == srcSelector &&
				uint64(expectedSeqNumRange.Start()) >= mr.MinSeqNr &&
				uint64(expectedSeqNumRange.End()) <= mr.MaxSeqNr {
				t.Logf("All sequence numbers committed in a single report [%d, %d]", expectedSeqNumRange.Start(), expectedSeqNumRange.End())
				return true, nil
			}

			if !enforceSingleCommit && seenMessages.allCommited(srcSelector) {
				t.Logf("All sequence numbers already committed from range [%d, %d]", expectedSeqNumRange.Start(), expectedSeqNumRange.End())
				return true, nil
			}
		case err := <-errCh:
			require.NoError(t, err)
		case <-timeout.C:
			return false, fmt.Errorf("timed out after waiting for commit report on chain selector %d from source selector %d expected seq nr range %s",
				dest.Selector, srcSelector, expectedSeqNumRange.String())
		}
	}
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
		seqRange := seqRange
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
				innerExecutionStates, err = ConfirmExecWithSeqNrsSol(
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

// ConfirmExecWithSeqNrs waits for an execution state change on the destination chain with the expected sequence number.
// startBlock is the block number to start watching from.
// If startBlock is nil, it will start watching from the latest block.
// Returns a map that maps the expected sequence number to its execution state.
func ConfirmExecWithSeqNrs(
	t *testing.T,
	sourceSelector uint64,
	dest cldf_evm.Chain,
	offRamp offramp.OffRampInterface,
	startBlock *uint64,
	expectedSeqNrs []uint64,
) (executionStates map[uint64]int, err error) {
	if len(expectedSeqNrs) == 0 {
		return nil, errors.New("no expected sequence numbers provided")
	}

	timeout := time.NewTimer(tests.WaitTimeout(t))
	defer timeout.Stop()
	tick := time.NewTicker(3 * time.Second)
	defer tick.Stop()
	sink := make(chan *offramp.OffRampExecutionStateChanged)
	subscription, err := offRamp.WatchExecutionStateChanged(&bind.WatchOpts{
		Context: context.Background(),
		Start:   startBlock,
	}, sink, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error to subscribe ExecutionStateChanged : %w", err)
	}
	defer subscription.Unsubscribe()

	// some state to efficiently track the execution states
	// of all the expected sequence numbers.
	executionStates = make(map[uint64]int)
	seqNrsToWatch := make(map[uint64]struct{})
	for _, seqNr := range expectedSeqNrs {
		seqNrsToWatch[seqNr] = struct{}{}
	}
	for {
		select {
		case <-tick.C:
			for expectedSeqNr := range seqNrsToWatch {
				scc, executionState := getExecutionState(t, sourceSelector, offRamp, expectedSeqNr)
				t.Logf("Waiting for ExecutionStateChanged on chain %d (offramp %s) from chain %d with expected sequence number %d, current onchain minSeqNr: %d, execution state: %s",
					dest.Selector, offRamp.Address().String(), sourceSelector, expectedSeqNr, scc.MinSeqNr, executionStateToString(executionState))
				if executionState == EXECUTION_STATE_SUCCESS || executionState == EXECUTION_STATE_FAILURE {
					t.Logf("Observed %s execution state on chain %d (offramp %s) from chain %d with expected sequence number %d",
						executionStateToString(executionState), dest.Selector, offRamp.Address().String(), sourceSelector, expectedSeqNr)
					executionStates[expectedSeqNr] = int(executionState)
					delete(seqNrsToWatch, expectedSeqNr)
					if len(seqNrsToWatch) == 0 {
						return executionStates, nil
					}
				}
			}
		case execEvent := <-sink:
			t.Logf("Received ExecutionStateChanged (state %s) for seqNum %d on chain %d (offramp %s) from chain %d",
				executionStateToString(execEvent.State), execEvent.SequenceNumber, dest.Selector, offRamp.Address().String(),
				sourceSelector,
			)

			_, found := seqNrsToWatch[execEvent.SequenceNumber]
			if found && execEvent.SourceChainSelector == sourceSelector {
				t.Logf("Received ExecutionStateChanged (state %s) on chain %d (offramp %s) from chain %d with expected sequence number %d",
					executionStateToString(execEvent.State), dest.Selector, offRamp.Address().String(), sourceSelector, execEvent.SequenceNumber)
				executionStates[execEvent.SequenceNumber] = int(execEvent.State)
				delete(seqNrsToWatch, execEvent.SequenceNumber)
				if len(seqNrsToWatch) == 0 {
					return executionStates, nil
				}
			}
		case <-timeout.C:
			return nil, fmt.Errorf("timed out waiting for ExecutionStateChanged on chain %d (offramp %s) from chain %d with expected sequence numbers %+v",
				dest.Selector, offRamp.Address().String(), sourceSelector, expectedSeqNrs)
		case subErr := <-subscription.Err():
			return nil, fmt.Errorf("subscription error: %w", subErr)
		}
	}
}

func ConfirmExecWithSeqNrsSol(
	t *testing.T,
	srcSelector uint64,
	dest cldf_solana.Chain,
	offrampAddress solana.PublicKey,
	startSlot uint64,
	expectedSeqNrs []uint64,
) (executionStates map[uint64]int, err error) {
	// TODO: share with EVM
	// some state to efficiently track the execution states
	// of all the expected sequence numbers.
	executionStates = make(map[uint64]int)
	seqNrsToWatch := make(map[uint64]struct{})
	for _, seqNr := range expectedSeqNrs {
		seqNrsToWatch[seqNr] = struct{}{}
	}

	done := make(chan any)
	defer close(done)
	sink, errCh := SolEventEmitter[solccip.EventExecutionStateChanged](t, dest.Client, offrampAddress, "ExecutionStateChanged", startSlot, done)

	timeout := time.NewTimer(tests.WaitTimeout(t))
	defer timeout.Stop()

	for {
		select {
		case execEvent := <-sink:
			// TODO: share with EVM
			_, found := seqNrsToWatch[execEvent.SequenceNumber]
			if found && execEvent.SourceChainSelector == srcSelector {
				t.Logf("Received ExecutionStateChanged (state %s) on chain %d (offramp %s) from chain %d with expected sequence number %d",
					execEvent.State.String(), dest.Selector, offrampAddress.String(), srcSelector, execEvent.SequenceNumber)
				if execEvent.State == ccip_offramp.InProgress_MessageExecutionState {
					// skip the in progress state, executed event should follow
					continue
				}
				executionStates[execEvent.SequenceNumber] = int(execEvent.State)
				delete(seqNrsToWatch, execEvent.SequenceNumber)
				if len(seqNrsToWatch) == 0 {
					return executionStates, nil
				}
			}
		case err := <-errCh:
			require.NoError(t, err)
		case <-timeout.C:
			return nil, fmt.Errorf("timed out waiting for ExecutionStateChanged on chain %d (offramp %s) from chain %d with expected sequence numbers %+v",
				dest.Selector, offrampAddress.String(), srcSelector, expectedSeqNrs)
		}
	}
}

func ConfirmNoExecConsistentlyWithSeqNr(
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
		if executionState == EXECUTION_STATE_UNTOUCHED {
			return true
		}
		t.Logf("Observed %s execution state on chain %d (offramp %s) from chain %d with expected sequence number %d",
			executionStateToString(executionState), dest.Selector, offRamp.Address().String(), sourceSelector, expectedSeqNr)
		return false
	}, timeout, 3*time.Second, "Expected no execution state change on chain %d (offramp %s) from chain %d with expected sequence number %d", dest.Selector, offRamp.Address().String(), sourceSelector, expectedSeqNr)
}

func getExecutionState(t *testing.T, sourceSelector uint64, offRamp offramp.OffRampInterface, expectedSeqNr uint64) (offramp.OffRampSourceChainConfig, uint8) {
	scc, err := offRamp.GetSourceChainConfig(nil, sourceSelector)
	require.NoError(t, err)
	executionState, err := offRamp.GetExecutionState(nil, sourceSelector, expectedSeqNr)
	require.NoError(t, err)
	return scc, executionState
}

func RequireConsistently(t *testing.T, condition func() bool, duration time.Duration, tick time.Duration, msgAndArgs ...interface{}) {
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
