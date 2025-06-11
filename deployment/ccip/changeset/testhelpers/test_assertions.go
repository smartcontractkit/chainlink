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
	"github.com/xssnick/tonutils-go/address"
	"golang.org/x/sync/errgroup"

	solconfig "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solccip "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/ccip"
	solcommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf_ton "github.com/smartcontractkit/chainlink-deployments-framework/chain/ton"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
)

type EVMDestConfirmer struct {
	DestChain cldf_evm.Chain
	OffRamp   offramp.OffRampInterface
}

func NewEVMDestConfirmer(dest cldf_evm.Chain, offramp offramp.OffRampInterface) *EVMDestConfirmer {
	return &EVMDestConfirmer{DestChain: dest, OffRamp: offramp}
}

type SolanaDestConfirmer struct {
	DestChain      cldf_solana.Chain
	OffRampAddress solana.PublicKey
}

func NewSolanaDestConfirmer(dest cldf_solana.Chain, offramp solana.PublicKey) *SolanaDestConfirmer {
	return &SolanaDestConfirmer{DestChain: dest, OffRampAddress: offramp}
}

type TonDestConfirmer struct {
	DestChain   cldf_ton.Chain
	OffRampAddr address.Address
}

func NewTonDestConfirmer(dest cldf_ton.Chain, offRampAddr address.Address) *TonDestConfirmer {
	return &TonDestConfirmer{DestChain: dest, OffRampAddr: offRampAddr}
}

type ConfirmCommitArgs struct {
	T                   *testing.T
	SrcSelector         uint64
	StartBlock          *uint64
	ExpectedSeqNumRange ccipocr3.SeqNumRange
	EnforceSingleCommit bool
}

type ConfirmExecArgs struct {
	T              *testing.T
	SourceChain    uint64
	StartBlock     *uint64
	ExpectedSeqNrs []uint64
}

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
		srcChainSel := sourceDest.SourceChainSelector
		dstChainSel := sourceDest.DestChainSelector
		if expectedSeqNum.Start() == 0 {
			continue
		}

		wg.Go(func() error {
			commitArgs := &ConfirmCommitArgs{
				T:                   t,
				SrcSelector:         srcChainSel,
				StartBlock:          startBlocks[dstChainSel],
				ExpectedSeqNumRange: expectedSeqNum,
				EnforceSingleCommit: true,
			}

			family, err := chainsel.GetSelectorFamily(dstChainSel)
			if err != nil {
				return err
			}

			switch family {
			case chainsel.FamilyEVM:
				confirmer := NewEVMDestConfirmer(
					e.BlockChains.EVMChains()[dstChainSel],
					state.MustGetEVMChainState(dstChainSel).OffRamp,
				)
				_, err = confirmer.ConfirmCommitWithExpectedSeqNumRange(commitArgs)
				return err
			case chainsel.FamilySolana:
				confirmer := NewSolanaDestConfirmer(
					e.BlockChains.SolanaChains()[dstChainSel],
					state.SolChains[dstChainSel].OffRamp,
				)
				_, err = confirmer.ConfirmCommitWithExpectedSeqNumRangeSol(commitArgs)
				return err
			case chainsel.FamilyTon:
				confirmer := NewTonDestConfirmer(
					e.BlockChains.TonChains()[dstChainSel],
					state.TonChains[dstChainSel].OffRamp,
				)
				_, err = confirmer.ConfirmCommitWithExpectedSeqNumRangeTon(commitArgs)
				return err
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
		srcChainSel := sourceDest.SourceChainSelector
		destChainSel := sourceDest.DestChainSelector
		seqRange := seqRange

		errGrp.Go(func() error {
			commitArgs := &ConfirmCommitArgs{
				T:                   t,
				SrcSelector:         srcChainSel,
				StartBlock:          startBlocks[destChainSel],
				ExpectedSeqNumRange: seqRange,
				EnforceSingleCommit: enforceSingleCommit,
			}

			family, err := chainsel.GetSelectorFamily(destChainSel)
			if err != nil {
				return err
			}
			switch family {
			case chainsel.FamilyEVM:
				confirmer := NewEVMDestConfirmer(
					env.BlockChains.EVMChains()[destChainSel],
					state.MustGetEVMChainState(destChainSel).OffRamp,
				)
				_, err = confirmer.ConfirmCommitWithExpectedSeqNumRange(commitArgs)
				return err
			case chainsel.FamilySolana:
				confirmer := NewSolanaDestConfirmer(
					env.BlockChains.SolanaChains()[destChainSel],
					state.SolChains[destChainSel].OffRamp,
				)
				_, err = confirmer.ConfirmCommitWithExpectedSeqNumRangeSol(commitArgs)
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
func (c *EVMDestConfirmer) ConfirmCommitWithExpectedSeqNumRange(args *ConfirmCommitArgs) (*offramp.OffRampCommitReportAccepted, error) {
	sink := make(chan *offramp.OffRampCommitReportAccepted)
	subscription, err := c.OffRamp.WatchCommitReportAccepted(&bind.WatchOpts{
		Context: context.Background(),
		Start:   args.StartBlock,
	}, sink)
	if err != nil {
		return nil, fmt.Errorf("error to subscribe CommitReportAccepted : %w", err)
	}

	seenMessages := NewCommitReportTracker(args.SrcSelector, args.ExpectedSeqNumRange)

	verifyCommitReport := func(report *offramp.OffRampCommitReportAccepted) bool {
		processRoots := func(roots []offramp.InternalMerkleRoot) bool {
			for _, mr := range roots {
				args.T.Logf(
					"Received commit report for [%d, %d] on selector %d from source selector %d expected seq nr range %s, token prices: %v",
					mr.MinSeqNr, mr.MaxSeqNr, c.DestChain.Selector, args.SrcSelector, args.ExpectedSeqNumRange.String(), report.PriceUpdates.TokenPriceUpdates,
				)
				seenMessages.visitCommitReport(args.SrcSelector, mr.MinSeqNr, mr.MaxSeqNr)

				if mr.SourceChainSelector == args.SrcSelector &&
					uint64(args.ExpectedSeqNumRange.Start()) >= mr.MinSeqNr &&
					uint64(args.ExpectedSeqNumRange.End()) <= mr.MaxSeqNr {
					args.T.Logf(
						"All sequence numbers committed in a single report [%d, %d]",
						args.ExpectedSeqNumRange.Start(), args.ExpectedSeqNumRange.End(),
					)
					return true
				}

				if !args.EnforceSingleCommit && seenMessages.allCommited(args.SrcSelector) {
					args.T.Logf(
						"All sequence numbers already committed from range [%d, %d]",
						args.ExpectedSeqNumRange.Start(), args.ExpectedSeqNumRange.End(),
					)
					return true
				}
			}
			return false
		}

		return processRoots(report.BlessedMerkleRoots) || processRoots(report.UnblessedMerkleRoots)
	}

	defer subscription.Unsubscribe()
	timeout := time.NewTimer(tests.WaitTimeout(args.T))
	defer timeout.Stop()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-args.T.Context().Done():
			return nil, nil
		case <-ticker.C:
			args.T.Logf("Waiting for commit report on chain selector %d from source selector %d expected seq nr range %s",
				c.DestChain.Selector, args.SrcSelector, args.ExpectedSeqNumRange.String())

			// Need to do this because the subscription sometimes fails to get the event.
			iter, err := c.OffRamp.FilterCommitReportAccepted(&bind.FilterOpts{
				Context: args.T.Context(),
			})
			// In some test case the test ends while the filter is still running resulting in a context.Canceled error.
			if err != nil && !errors.Is(err, context.Canceled) {
				require.NoError(args.T, err)
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
				c.DestChain.Selector, args.SrcSelector, args.ExpectedSeqNumRange.String())
		case report := <-sink:
			verified := verifyCommitReport(report)
			if verified {
				return report, nil
			}
		}
	}
}

type EventWithTxn[T any] struct {
	Event T
	Txn   *solrpc.GetTransactionResult
}

// Scan for events referencing address
func SolEventEmitter[T any](ctx context.Context, client *solrpc.Client, address solana.PublicKey, eventType string, startSlot uint64, done chan any, ticker *time.Ticker) (<-chan EventWithTxn[T], <-chan error) {
	ch := make(chan EventWithTxn[T])
	errorCh := make(chan error)
	go func() {
		defer ticker.Stop()
		var until solana.Signature
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				// Scan for transactions referencing the address
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
						case ch <- EventWithTxn[T]{
							Event: event,
							Txn:   tx,
						}:
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

func (c *SolanaDestConfirmer) ConfirmCommitWithExpectedSeqNumRangeSol(args *ConfirmCommitArgs) (bool, error) {
	seenMessages := NewCommitReportTracker(args.SrcSelector, args.ExpectedSeqNumRange)

	var startSlot uint64
	if args.StartBlock != nil {
		startSlot = *args.StartBlock
	}

	done := make(chan any)
	defer close(done)
	sink, errCh := SolEventEmitter[solccip.EventCommitReportAccepted](t.Context(), dest.Client, offrampAddress, consts.EventNameCommitReportAccepted, startSlot, done, time.NewTicker(2*time.Second))

	timeout := time.NewTimer(tests.WaitTimeout(args.T))
	defer timeout.Stop()

	for {
		select {
		case eventWithTxn := <-sink:
			commitEvent := eventWithTxn.Event
			// if merkle root is zero, it only contains price updates
			if commitEvent.Report == nil {
				args.T.Logf("Skipping CommitReportAccepted with only price updates")
				continue
			}
			require.Equal(args.T, args.SrcSelector, commitEvent.Report.SourceChainSelector)

			// TODO: this logic is duplicated with verifyCommitReport, share
			mr := commitEvent.Report
			seenMessages.visitCommitReport(mr.SourceChainSelector, mr.MinSeqNr, mr.MaxSeqNr)
			if mr.SourceChainSelector == args.SrcSelector &&
				uint64(args.ExpectedSeqNumRange.Start()) >= mr.MinSeqNr &&
				uint64(args.ExpectedSeqNumRange.End()) <= mr.MaxSeqNr {
				args.T.Logf("All sequence numbers committed in a single report [%d, %d]", args.ExpectedSeqNumRange.Start(), args.ExpectedSeqNumRange.End())
				return true, nil
			}

			if !args.EnforceSingleCommit && seenMessages.allCommited(args.SrcSelector) {
				args.T.Logf("All sequence numbers already committed from range [%d, %d]", args.ExpectedSeqNumRange.Start(), args.ExpectedSeqNumRange.End())
				return true, nil
			}
		case err := <-errCh:
			require.NoError(args.T, err)
		case <-timeout.C:
			return false, fmt.Errorf("timed out after waiting for commit report on chain selector %d from source selector %d expected seq nr range %s",
				c.DestChain.Selector, args.SrcSelector, args.ExpectedSeqNumRange.String())
		}
	}
}

func (c *TonDestConfirmer) ConfirmCommitWithExpectedSeqNumRangeTon(args *ConfirmCommitArgs) (any, error) {
	// TODO once offramp contracts are supported, we can add the logic to confirm commit with expected sequence number range
	return true, nil
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
	for sourceDest, seqNrs := range expectedSeqNums {
		sourceDest := sourceDest
		seqNrs := seqNrs

		wg.Go(func() error {
			execArgs := &ConfirmExecArgs{
				T:              t,
				SourceChain:    sourceDest.SourceChainSelector,
				StartBlock:     startBlocks[sourceDest.DestChainSelector],
				ExpectedSeqNrs: seqNrs,
			}

			family, err := chainsel.GetSelectorFamily(sourceDest.DestChainSelector)
			if err != nil {
				return err
			}

			var innerExecutionStates map[uint64]int
			switch family {
			case chainsel.FamilyEVM:
				confirmer := NewEVMDestConfirmer(
					e.BlockChains.EVMChains()[sourceDest.DestChainSelector],
					state.MustGetEVMChainState(sourceDest.DestChainSelector).OffRamp,
				)
				innerExecutionStates, err = confirmer.ConfirmExecWithSeqNrs(execArgs)
			case chainsel.FamilySolana:
				confirmer := NewSolanaDestConfirmer(
					e.BlockChains.SolanaChains()[sourceDest.DestChainSelector],
					state.SolChains[sourceDest.DestChainSelector].OffRamp,
				)
				innerExecutionStates, err = confirmer.ConfirmExecWithSeqNrsSol(execArgs)
			case chainsel.FamilyTon:
				confirmer := NewTonDestConfirmer(
					e.BlockChains.TonChains()[sourceDest.DestChainSelector],
					state.TonChains[sourceDest.DestChainSelector].OffRamp,
				)
				innerExecutionStates, err = confirmer.ConfirmExecWithSeqNrsTon(execArgs)
			}

			if err != nil {
				return err
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
func (c *EVMDestConfirmer) ConfirmExecWithSeqNrs(args *ConfirmExecArgs) (executionStates map[uint64]int, err error) {
	if len(args.ExpectedSeqNrs) == 0 {
		return nil, errors.New("no expected sequence numbers provided")
	}

	timeout := time.NewTimer(tests.WaitTimeout(args.T))
	defer timeout.Stop()
	tick := time.NewTicker(3 * time.Second)
	defer tick.Stop()
	sink := make(chan *offramp.OffRampExecutionStateChanged)
	subscription, err := c.OffRamp.WatchExecutionStateChanged(&bind.WatchOpts{
		Context: context.Background(),
		Start:   args.StartBlock,
	}, sink, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error to subscribe ExecutionStateChanged : %w", err)
	}
	defer subscription.Unsubscribe()

	// some state to efficiently track the execution states
	// of all the expected sequence numbers.
	executionStates = make(map[uint64]int)
	seqNrsToWatch := make(map[uint64]struct{})
	for _, seqNr := range args.ExpectedSeqNrs {
		seqNrsToWatch[seqNr] = struct{}{}
	}
	for {
		select {
		case <-tick.C:
			for expectedSeqNr := range seqNrsToWatch {
				scc, executionState := getExecutionState(args.T, args.SourceChain, c.OffRamp, expectedSeqNr)
				args.T.Logf("Waiting for ExecutionStateChanged on chain %d (offramp %s) from chain %d with expected sequence number %d, current onchain minSeqNr: %d, execution state: %s",
					c.DestChain.Selector, c.OffRamp.Address().String(), args.SourceChain, expectedSeqNr, scc.MinSeqNr, executionStateToString(executionState))
				if executionState == EXECUTION_STATE_SUCCESS || executionState == EXECUTION_STATE_FAILURE {
					args.T.Logf("Observed %s execution state on chain %d (offramp %s) from chain %d with expected sequence number %d",
						executionStateToString(executionState), c.DestChain.Selector, c.OffRamp.Address().String(), args.SourceChain, expectedSeqNr)
					executionStates[expectedSeqNr] = int(executionState)
					delete(seqNrsToWatch, expectedSeqNr)
					if len(seqNrsToWatch) == 0 {
						return executionStates, nil
					}
				}
			}
		case execEvent := <-sink:
			args.T.Logf("Received ExecutionStateChanged (state %s) for seqNum %d on chain %d (offramp %s) from chain %d",
				executionStateToString(execEvent.State), execEvent.SequenceNumber, c.DestChain.Selector, c.OffRamp.Address().String(),
				args.SourceChain,
			)

			_, found := seqNrsToWatch[execEvent.SequenceNumber]
			if found && execEvent.SourceChainSelector == args.SourceChain {
				args.T.Logf("Received ExecutionStateChanged (state %s) on chain %d (offramp %s) from chain %d with expected sequence number %d",
					executionStateToString(execEvent.State), c.DestChain.Selector, c.OffRamp.Address().String(), args.SourceChain, execEvent.SequenceNumber)
				executionStates[execEvent.SequenceNumber] = int(execEvent.State)
				delete(seqNrsToWatch, execEvent.SequenceNumber)
				if len(seqNrsToWatch) == 0 {
					return executionStates, nil
				}
			}
		case <-timeout.C:
			return nil, fmt.Errorf("timed out waiting for ExecutionStateChanged on chain %d (offramp %s) from chain %d with expected sequence numbers %+v",
				c.DestChain.Selector, c.OffRamp.Address().String(), args.SourceChain, args.ExpectedSeqNrs)
		case subErr := <-subscription.Err():
			return nil, fmt.Errorf("subscription error: %w", subErr)
		}
	}
}

func (c *SolanaDestConfirmer) ConfirmExecWithSeqNrsSol(args *ConfirmExecArgs) (executionStates map[uint64]int, err error) {
	// TODO: share with EVM
	// some state to efficiently track the execution states
	// of all the expected sequence numbers.
	executionStates = make(map[uint64]int)
	seqNrsToWatch := make(map[uint64]struct{})
	for _, seqNr := range args.ExpectedSeqNrs {
		seqNrsToWatch[seqNr] = struct{}{}
	}

	var startSlot uint64
	if args.StartBlock != nil {
		startSlot = *args.StartBlock
	}

	done := make(chan any)
	defer close(done)
	sink, errCh := SolEventEmitter[solccip.EventExecutionStateChanged](t.Context(), dest.Client, offrampAddress, consts.EventNameExecutionStateChanged, startSlot, done, time.NewTicker(2*time.Second))

	timeout := time.NewTimer(tests.WaitTimeout(args.T))
	defer timeout.Stop()

	for {
		select {
		case eventWithTxn := <-sink:
			execEvent := eventWithTxn.Event
			// TODO: share with EVM
			_, found := seqNrsToWatch[execEvent.SequenceNumber]
			if found && execEvent.SourceChainSelector == args.SourceChain {
				args.T.Logf("Received ExecutionStateChanged (state %s) on chain %d (offramp %s) from chain %d with expected sequence number %d",
					execEvent.State.String(), c.DestChain.Selector, c.OffRampAddress.String(), args.SourceChain, execEvent.SequenceNumber)
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
			require.NoError(args.T, err)
		case <-timeout.C:
			return nil, fmt.Errorf("timed out waiting for ExecutionStateChanged on chain %d (offramp %s) from chain %d with expected sequence numbers %+v",
				c.DestChain.Selector, c.OffRampAddress.String(), args.SourceChain, args.ExpectedSeqNrs)
		}
	}
}

func (c *TonDestConfirmer) ConfirmExecWithSeqNrsTon(args *ConfirmExecArgs) (executionStates map[uint64]int, err error) {
	args.T.Logf("DEBUG: TODO(ton): ConfirmExecWithSeqNrsTon\n")
	seqNrsToWatch := make(map[uint64]int)
	for _, seqNr := range args.ExpectedSeqNrs {
		seqNrsToWatch[seqNr] = EXECUTION_STATE_SUCCESS
	}
	return seqNrsToWatch, nil
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
