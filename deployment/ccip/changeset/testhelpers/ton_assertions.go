package testhelpers

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ccipocr3common "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	cldf_ton "github.com/smartcontractkit/chainlink-deployments-framework/chain/ton"

	"github.com/smartcontractkit/chainlink-ton/pkg/ccip/bindings/offramp"
	tonlploader "github.com/smartcontractkit/chainlink-ton/pkg/logpoller/backend/loader/account"
	tonlptypes "github.com/smartcontractkit/chainlink-ton/pkg/logpoller/types"
	"github.com/smartcontractkit/chainlink-ton/pkg/ton/event"
	"github.com/smartcontractkit/chainlink-ton/pkg/ton/hash"
)

var (
	// ErrTimeout is returned when event subscription times out
	ErrTimeout = errors.New("timed out waiting for events")
)

// TON blockchain polling configuration
const (
	safeLookbackBlocks  = uint32(50)              // Number of blocks to look back when starting event scan
	pollInterval        = 2500 * time.Millisecond // How often to check for new blocks
	clientRetries       = 3                       // Number of retries for TON client operations
	txBatchSize         = 100                     // Number of transactions to fetch per batch
	progressLogInterval = 15 * time.Second        // How often to log "still waiting" progress updates
)

// eventSubscriber encapsulates subscribing to and waiting for TON events with timeout handling.
type eventSubscriber[T any] struct {
	lggr       logger.Logger
	tonChain   cldf_ton.Chain
	contract   *address.Address
	startBlock uint32
	eventName  string
}

// subscribeToEvents creates a new event subscriber for a specific event type.
// If startBlock is 0, it will calculate an appropriate start block with lookback.
func subscribeToEvents[T any](
	ctx context.Context,
	lggr logger.Logger,
	tonChain cldf_ton.Chain,
	contract *address.Address,
	startBlock uint32,
	eventName string,
) *eventSubscriber[T] {
	sub := &eventSubscriber[T]{
		lggr:       lggr,
		tonChain:   tonChain,
		contract:   contract,
		startBlock: startBlock,
		eventName:  eventName,
	}

	// Calculate start block with lookback if not provided
	if startBlock == 0 {
		sub.startBlock = sub.calculateStartBlock(ctx)
	}

	return sub
}

// calculateStartBlock determines the starting block for event scanning with lookback.
func (s *eventSubscriber[T]) calculateStartBlock(ctx context.Context) uint32 {
	currentBlock, err := s.tonChain.Client.CurrentMasterchainInfo(ctx)
	if err != nil || currentBlock.SeqNo <= safeLookbackBlocks {
		return 0
	}
	return currentBlock.SeqNo - safeLookbackBlocks
}

// waitUntil waits for events and processes them with the provided handler.
// The handler returns (done bool, error) where done=true stops waiting.
func (s *eventSubscriber[T]) waitUntil(ctx context.Context, until time.Duration, processEvent func(T) (bool, error)) error {
	eventCh, errCh := streamEvents[T](ctx, s.lggr, s.tonChain, s.contract, s.startBlock, s.eventName)

	timeout := time.NewTimer(until)
	defer timeout.Stop()

	progressTicker := time.NewTicker(progressLogInterval)
	defer progressTicker.Stop()

	eventCount := 0
	startTime := time.Now()

	for {
		select {
		case event := <-eventCh:
			eventCount++
			done, err := processEvent(event)
			if err != nil {
				return err
			}
			if done {
				return nil
			}

		case err := <-errCh:
			return err

		case <-progressTicker.C:
			s.lggr.Infow("Still waiting",
				"eventName", s.eventName,
				"elapsed", time.Since(startTime).Round(time.Second).String())

		case <-timeout.C:
			return fmt.Errorf("%w: after %s", ErrTimeout, time.Since(startTime).String())
		}
	}
}

// ConfirmCommitWithExpectedSeqNumRangeTON waits for a commit report that covers the expected sequence number range.
func ConfirmCommitWithExpectedSeqNumRangeTON(
	t *testing.T,
	srcChainSelector uint64,
	tonChain cldf_ton.Chain,
	offRamp address.Address,
	expectedSeqNums ccipocr3common.SeqNumRange,
) (bool, error) {
	seqNumTracker := NewCommitReportTracker(srcChainSelector, expectedSeqNums)

	lggr := logger.Named(logger.Test(t), "TON_EVENT_ASSERTION:COMMIT")
	lggr.Infow("Waiting for commit report",
		"sourceChain", srcChainSelector,
		"destChain", tonChain.Selector,
		"expectedSeqNums", fmt.Sprintf("[%d, %d]", expectedSeqNums.Start(), expectedSeqNums.End()))

	ctx := t.Context()
	subscriber := subscribeToEvents[offramp.CommitReportAccepted](ctx, lggr, tonChain, &offRamp, 0, consts.EventNameCommitReportAccepted)

	reportsProcessed := 0
	err := subscriber.waitUntil(ctx, tests.WaitTimeout(t), func(commitEvent offramp.CommitReportAccepted) (bool, error) {
		reportsProcessed++

		// skip price-only OR empty updates
		if commitEvent.MerkleRoot == nil {
			return false, nil // continue waiting
		}

		mr := commitEvent.MerkleRoot
		require.Equal(t, srcChainSelector, mr.SourceChainSelector,
			"Commit report source chain mismatch")

		lggr.Infow("Received commit",
			"seqNums", fmt.Sprintf("[%d, %d]", mr.MinSeqNr, mr.MaxSeqNr))

		// Track this commit report
		seqNumTracker.visitCommitReport(srcChainSelector, mr.MinSeqNr, mr.MaxSeqNr)

		// Check if single report covers entire range
		if uint64(expectedSeqNums.Start()) >= mr.MinSeqNr &&
			uint64(expectedSeqNums.End()) <= mr.MaxSeqNr {
			t.Logf("All sequence numbers committed in a single report [%d, %d]",
				expectedSeqNums.Start(), expectedSeqNums.End())
			return true, nil
		}

		// Check if all messages committed across multiple reports
		if seqNumTracker.allCommited(srcChainSelector) {
			t.Logf("All sequence numbers committed across multiple reports [%d, %d]",
				expectedSeqNums.Start(), expectedSeqNums.End())
			return true, nil
		}

		return false, nil // continue waiting
	})

	if err != nil {
		// Add detailed context to timeout error
		if errors.Is(err, ErrTimeout) {
			lggr.Errorw("Commit confirmation timed out",
				"destChain", tonChain.Selector,
				"sourceChain", srcChainSelector,
				"expectedSeqNums", expectedSeqNums.String(),
				"reportsProcessed", reportsProcessed)
			return false, fmt.Errorf("timed out after waiting for commit report on chain selector %d from source chain %d expected seq nums %s (processed %d reports): %w",
				tonChain.Selector, srcChainSelector, expectedSeqNums.String(), reportsProcessed, err)
		}
		require.NoError(t, err)
		return false, err
	}

	return true, nil
}

// ConfirmExecWithExpectedSeqNrsTON waits for execution state changes on TON for the given sequence numbers.
// Returns a map of sequence number to execution state.
func ConfirmExecWithExpectedSeqNrsTON(
	t *testing.T,
	srcChainSelector uint64,
	tonChain cldf_ton.Chain,
	offRamp address.Address,
	startBlock *uint64,
	expectedSeqNums []uint64,
) (map[uint64]int, error) {
	if len(expectedSeqNums) == 0 {
		return nil, errors.New("no expected sequence numbers provided")
	}

	lggr := logger.Named(logger.Test(t), "TON_EVENT_ASSERTION:EXEC")
	lggr.Infow("Waiting for execution",
		"sourceChain", srcChainSelector,
		"destChain", tonChain.Selector,
		"expectedSeqNums", expectedSeqNums)

	// Use provided start block or calculate with lookback
	var scanStartBlock uint32
	if startBlock != nil {
		scanStartBlock = uint32(*startBlock) //nolint:gosec // safe conversion, test node
	}

	ctx := t.Context()
	subscriber := subscribeToEvents[offramp.ExecutionStateChanged](ctx, lggr, tonChain, &offRamp, scanStartBlock, consts.EventNameExecutionStateChanged)

	executionStates := make(map[uint64]int)
	pending := make(map[uint64]bool)
	for _, seqNum := range expectedSeqNums {
		pending[seqNum] = true
	}

	eventsProcessed := 0
	err := subscriber.waitUntil(ctx, tests.WaitTimeout(t), func(execEvent offramp.ExecutionStateChanged) (bool, error) {
		eventsProcessed++

		// Check if this is for our source chain and expected sequence number
		if execEvent.SourceChainSelector != srcChainSelector {
			return false, nil // continue waiting
		}
		if _, expected := pending[execEvent.SequenceNumber]; !expected {
			return false, nil // continue waiting
		}

		// Handle different execution states
		switch execEvent.State {
		case EXECUTION_STATE_INPROGRESS:
			// Don't log IN_PROGRESS, it's just noise
			return false, nil // continue waiting

		case EXECUTION_STATE_FAILURE:
			lggr.Errorw("Execution failed",
				"sequenceNumber", execEvent.SequenceNumber,
				"messageID", hex.EncodeToString(execEvent.MessageID))
			return false, fmt.Errorf("execution failed for sequence number %d on chain %d, message ID: %x",
				execEvent.SequenceNumber, execEvent.SourceChainSelector, execEvent.MessageID)

		case EXECUTION_STATE_SUCCESS:
			executionStates[execEvent.SequenceNumber] = int(execEvent.State)
			delete(pending, execEvent.SequenceNumber)

			lggr.Infow("Execution successful",
				"sequenceNumber", execEvent.SequenceNumber,
				"remaining", len(pending))

			if len(pending) == 0 {
				t.Logf("All sequence numbers executed: %v", expectedSeqNums)
				return true, nil // done
			}

		default:
			lggr.Warnw("Unknown execution state",
				"state", execEvent.State,
				"sequenceNumber", execEvent.SequenceNumber)
		}

		return false, nil // continue waiting
	})

	if err != nil {
		// Add detailed context to timeout error
		if errors.Is(err, ErrTimeout) {
			missing := make([]uint64, 0, len(pending))
			for seqNum := range pending {
				missing = append(missing, seqNum)
			}
			lggr.Errorw("Execution confirmation timed out",
				"destChain", tonChain.Selector,
				"sourceChain", srcChainSelector,
				"expectedSeqNums", expectedSeqNums,
				"missingSeqNums", missing,
				"successfulExecutions", len(executionStates),
				"eventsProcessed", eventsProcessed)
			return executionStates, fmt.Errorf("timed out after waiting for execution on chain selector %d from source chain %d, missing seq nums: %v (processed %d events, %d successful): %w",
				tonChain.Selector, srcChainSelector, missing, eventsProcessed, len(executionStates), err)
		}
		return nil, fmt.Errorf("error while waiting for execution events: %w", err)
	}

	return executionStates, nil
}

// streamEvents continuously polls the TON blockchain for events of type T.
// Returns two channels: one for events and one for errors.
// The polling continues until the context is cancelled.
func streamEvents[T any](
	ctx context.Context,
	lggr logger.Logger,
	tonChain cldf_ton.Chain,
	contract *address.Address,
	startBlock uint32,
	eventName string,
) (<-chan T, <-chan error) {
	eventCh := make(chan T)
	errorCh := make(chan error)

	go func() {
		defer close(eventCh)
		defer close(errorCh)

		ticker := time.NewTicker(pollInterval)
		defer ticker.Stop()

		// Client provider with retry logic
		clientProvider := func(ctx context.Context) (ton.APIClientWrapped, error) {
			return tonChain.Client.WithRetry(clientRetries), nil
		}

		// Initialize transaction loader
		loader := tonlploader.NewTxLoader(lggr, clientProvider, txBatchSize)

		lastProcessedBlock := startBlock

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// 1. Get block range
				blockRange, newSeqNo, err := getBlockRange(ctx, tonChain, lastProcessedBlock)
				if err != nil {
					errorCh <- err
					return
				}

				// Skip if no new blocks
				if blockRange == nil {
					continue
				}

				// 2. Fetch transactions
				txs, err := loader.FetchTxsForAddress(ctx, blockRange, contract)
				if err != nil {
					errorCh <- fmt.Errorf("failed to load transactions: %w", err)
					return
				}

				// 3. Extract and filter events
				events, err := extractEventMessage[T](txs, lggr, eventName)
				if err != nil {
					errorCh <- fmt.Errorf("failed to extract events from block %d: %w", newSeqNo, err)
					return
				}

				// Send events to channel
				for _, event := range events {
					select {
					case eventCh <- event:
					case <-ctx.Done():
						return
					}
				}

				// Update last processed block
				lastProcessedBlock = newSeqNo
			}
		}
	}()

	return eventCh, errorCh
}

// getBlockRange creates a block range from lastProcessedBlock to current masterchain head.
// Returns nil blockRange if there are no new blocks to process.
func getBlockRange(ctx context.Context, tonChain cldf_ton.Chain, lastProcessedBlock uint32) (*tonlptypes.BlockRange, uint32, error) {
	client := tonChain.Client.WithRetry(clientRetries)

	// Get current block
	toBlock, err := client.CurrentMasterchainInfo(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get current masterchain info: %w", err)
	}

	// No new blocks to process
	if toBlock.SeqNo <= lastProcessedBlock {
		return nil, toBlock.SeqNo, nil
	}

	// Lookup previous block if we have a starting point
	var prevBlock *ton.BlockIDExt
	if lastProcessedBlock > 0 {
		prevBlock, err = client.LookupBlock(ctx, toBlock.Workchain, toBlock.Shard, lastProcessedBlock)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to lookup previous block %d: %w", lastProcessedBlock, err)
		}
	}

	return &tonlptypes.BlockRange{Prev: prevBlock, To: toBlock}, toBlock.SeqNo, nil
}

// tryParseEvent attempts to parse an event of type T from a single TON message.
// Only parses events that match the expected eventName topic.
// Returns the event and true if successful, or zero value and false if parsing failed.
func tryParseEvent[T any](msg *tlb.Message, lggr logger.Logger, expectedEventName string) (T, bool) {
	var zero T

	if msg.MsgType != tlb.MsgTypeExternalOut {
		return zero, false
	}

	extOut := msg.AsExternalOut()
	if extOut == nil {
		return zero, false
	}

	// decode event topic
	bucket := event.NewExtOutLogBucket(extOut.DestAddr())
	topic, err := bucket.DecodeEventTopic()
	if err != nil {
		return zero, false
	}

	// Filter by expected event topic first - avoid parsing wrong event types
	expectedTopic := hash.CRC32(expectedEventName)
	if topic != expectedTopic {
		return zero, false
	}

	bodyCell := extOut.Payload()
	if bodyCell == nil {
		return zero, false
	}

	// parse event using TLB
	var parsedEvent T
	if err := tlb.LoadFromCell(&parsedEvent, bodyCell.BeginParse()); err != nil {
		lggr.Warnw("Failed to parse event body",
			"eventName", expectedEventName,
			"topic", fmt.Sprintf("0x%08x", topic),
			"error", err)
		return zero, false
	}

	return parsedEvent, true
}

// extractEventMessage processes transactions to extract events of type T from external messages.
// Only processes events matching the specified eventName topic.
func extractEventMessage[T any](txs []tonlptypes.TxWithBlock, lggr logger.Logger, eventName string) ([]T, error) {
	var events []T

	for _, tx := range txs {
		if tx.Tx == nil || tx.Tx.IO.Out == nil {
			continue
		}

		msgs, err := tx.Tx.IO.Out.ToSlice()
		if err != nil {
			continue
		}

		for _, msg := range msgs {
			if event, ok := tryParseEvent[T](&msg, lggr, eventName); ok {
				events = append(events, event)
			}
		}
	}

	return events, nil
}
