package testhelpers

import (
	"context"
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

const tonTickerInterval = 2500 * time.Millisecond
const safeLookbackBlocks = 50

// CCIPEventTopics maps TON event topics (CRC32 hashes) to their event names.
// OffRamp events used in these test assertions (EVM→TON flow).
var CCIPEventTopics = map[uint32]string{
	hash.CRC32(consts.EventNameCommitReportAccepted):  consts.EventNameCommitReportAccepted,
	hash.CRC32(consts.EventNameExecutionStateChanged): consts.EventNameExecutionStateChanged,
}

// getEventName returns the event name for a given topic, or a formatted unknown topic string.
func getEventName(topic uint32) string {
	if name, ok := CCIPEventTopics[topic]; ok {
		return name
	}
	return fmt.Sprintf("Unknown_0x%08x", topic)
}

func ConfirmCommitWithExpectedSeqNumRangeTON(t *testing.T, srcSelector uint64, tonChain cldf_ton.Chain, offRampContract address.Address, expectedSeqNumRange ccipocr3common.SeqNumRange) (bool, error) {
	seenMessages := NewCommitReportTracker(srcSelector, expectedSeqNumRange)

	// Create prefixed logger for TON event assertion
	lggr := logger.Named(logger.Test(t), "TON_EVENT_ASSERTION:COMMIT")
	lggr.Infof("waiting for commit report from srcSelector=%d, expectedSeqNumRange=[%d, %d], timeout=%v",
		srcSelector, expectedSeqNumRange.Start(), expectedSeqNumRange.End(), tests.WaitTimeout(t))

	tonBlockTicker := time.NewTicker(tonTickerInterval)
	// Start from a bit earlier to catch any events we might have missed
	// Use a lookback of about 50 blocks to ensure we don't miss commit events
	var startBlock uint32 = 0
	if currentBlock, err := tonChain.Client.CurrentMasterchainInfo(t.Context()); err == nil && currentBlock.SeqNo > safeLookbackBlocks {
		startBlock = currentBlock.SeqNo - safeLookbackBlocks
		lggr.Infof("scan from block %d (50 blocks back from current %d)", startBlock, currentBlock.SeqNo)
	}

	eventCh, errCh := GetEvents[offramp.CommitReportAccepted](t, lggr, t.Context(), tonChain, &offRampContract, startBlock, tonBlockTicker)

	timeout := time.NewTimer(tests.WaitTimeout(t))
	defer timeout.Stop()

	for {
		select {
		case commitEvent := <-eventCh:
			// Log event details with proper dereferencing
			if commitEvent.MerkleRoot != nil {
				lggr.Infof("received CommitReportAccepted event: MerkleRoot={SourceChainSelector:%d, MinSeqNr:%d, MaxSeqNr:%d, MerkleRoot:%x}, PriceUpdates=%+v",
					commitEvent.MerkleRoot.SourceChainSelector, commitEvent.MerkleRoot.MinSeqNr, commitEvent.MerkleRoot.MaxSeqNr,
					commitEvent.MerkleRoot.MerkleRoot, commitEvent.PriceUpdates)
			} else {
				lggr.Infof("received CommitReportAccepted event: MerkleRoot=<nil>, PriceUpdates=%+v", commitEvent.PriceUpdates)
			}

			// if merkle root is zero, it only contains price updates
			if commitEvent.MerkleRoot == nil {
				lggr.Infof("Skipping CommitReportAccepted with only price updates")
				continue
			}

			mr := commitEvent.MerkleRoot
			require.Equal(t, srcSelector, mr.SourceChainSelector)

			// TODO: this logic is duplicated with verifyCommitReport, share
			seenMessages.visitCommitReport(commitEvent.MerkleRoot.SourceChainSelector, mr.MinSeqNr, mr.MaxSeqNr)
			if mr.SourceChainSelector == srcSelector &&
				uint64(expectedSeqNumRange.Start()) >= mr.MinSeqNr &&
				uint64(expectedSeqNumRange.End()) <= mr.MaxSeqNr {
				t.Logf("All sequence numbers committed in a single report [%d, %d]", expectedSeqNumRange.Start(), expectedSeqNumRange.End())
				return true, nil
			}

			if seenMessages.allCommited(srcSelector) {
				t.Logf("All sequence numbers already committed from range [%d, %d]", expectedSeqNumRange.Start(), expectedSeqNumRange.End())
				return true, nil
			}
		case err := <-errCh:
			require.NoError(t, err)
		case <-timeout.C:
			return false, fmt.Errorf("timed out after waiting for commit report on chain selector %d from source selector %d expected seq nr range %s",
				tonChain.Selector, srcSelector, expectedSeqNumRange.String())
		}
	}
}

// ConfirmExecWithSeqNrsTON waits for execution state changes on TON for the given sequence numbers
// Returns a map of sequence number to execution state
func ConfirmExecWithSeqNrsTON(
	t *testing.T,
	srcSelector uint64,
	tonChain cldf_ton.Chain,
	offRampContract address.Address,
	startBlock *uint64,
	expectedSeqNrs []uint64,
) (map[uint64]int, error) {
	if len(expectedSeqNrs) == 0 {
		return nil, fmt.Errorf("no expected sequence numbers provided")
	}

	// Create prefixed logger for TON event assertion
	lggr := logger.Named(logger.Test(t), "TON_EVENT_ASSERTION:EXEC")
	lggr.Infof("waiting for execution state changes from srcSelector=%d, expectedSeqNrs=%v, timeout=%v",
		srcSelector, expectedSeqNrs, tests.WaitTimeout(t))

	executionStates := make(map[uint64]int)

	tonBlockTicker := time.NewTicker(tonTickerInterval)

	// Determine start block
	var scanStartBlock uint32 = 0
	if startBlock != nil {
		scanStartBlock = uint32(*startBlock)
	} else if currentBlock, err := tonChain.Client.CurrentMasterchainInfo(t.Context()); err == nil && currentBlock.SeqNo > safeLookbackBlocks {
		scanStartBlock = currentBlock.SeqNo - safeLookbackBlocks
		lggr.Infof("scan from block %d (50 blocks back from current %d)", scanStartBlock, currentBlock.SeqNo)
	}

	eventCh, errCh := GetEvents[offramp.ExecutionStateChanged](t, lggr, t.Context(), tonChain, &offRampContract, scanStartBlock, tonBlockTicker)

	timeout := time.NewTimer(tests.WaitTimeout(t))
	defer timeout.Stop()

	// Track which sequence numbers we're waiting for
	expectedSeqNums := make(map[uint64]bool)
	for _, seqNum := range expectedSeqNrs {
		expectedSeqNums[seqNum] = true
	}

	for {
		select {
		case execEvent := <-eventCh:
			lggr.Infof("received ExecutionStateChanged event: SourceChainSelector=%d, SequenceNumber=%d, MessageID=%x, State=%d",
				execEvent.SourceChainSelector, execEvent.SequenceNumber, execEvent.MessageID, execEvent.State)

			// Check if this is for our source chain
			if execEvent.SourceChainSelector != srcSelector {
				lggr.Debugf("skipping event from different source chain: %d (expected %d)", execEvent.SourceChainSelector, srcSelector)
				continue
			}

			// Check if this is a sequence number we're waiting for
			if _, expected := expectedSeqNums[execEvent.SequenceNumber]; expected {
				// State 1 = IN_PROGRESS, State 2 = SUCCESS, State 3 = FAILURE
				// Only record and remove from expected if we got SUCCESS state
				if execEvent.State == 1 {
					lggr.Infof("received IN_PROGRESS state for seq num %d, waiting for SUCCESS state", execEvent.SequenceNumber)
					continue
				}

				if execEvent.State == 3 {
					lggr.Errorf("Execution failure detected for seq num %d, source chain %d, message id %x, state %d", execEvent.SequenceNumber, execEvent.SourceChainSelector, execEvent.MessageID, execEvent.State)
					// Don't delete from expected, don't record, just continue
					continue
				}

				executionStates[execEvent.SequenceNumber] = int(execEvent.State)
				delete(expectedSeqNums, execEvent.SequenceNumber)

				lggr.Infof("found execution state for seq num %d: state=%d, remaining=%d",
					execEvent.SequenceNumber, execEvent.State, len(expectedSeqNums))
			}

			// if we've found all expected sequence numbers, return
			if len(expectedSeqNums) == 0 {
				t.Logf("All sequence numbers executed: %v", expectedSeqNrs)
				return executionStates, nil
			}

		case err := <-errCh:
			return nil, fmt.Errorf("error while waiting for execution events: %w", err)

		case <-timeout.C:
			missing := make([]uint64, 0, len(expectedSeqNums))
			for seqNum := range expectedSeqNums {
				missing = append(missing, seqNum)
			}
			return executionStates, fmt.Errorf("timed out after waiting for execution on chain selector %d from source selector %d, missing seq nums: %v",
				tonChain.Selector, srcSelector, missing)
		}
	}
}

func GetEvents[T any](t *testing.T, lggr logger.Logger, ctx context.Context, tonChain cldf_ton.Chain, contractAddress *address.Address, startBlock uint32, ticker *time.Ticker) (<-chan T, <-chan error) {
	ch := make(chan T)
	errorCh := make(chan error)

	go func() {
		defer ticker.Stop()
		defer close(ch)
		defer close(errorCh)

		// lazy client provider
		clientProvider := func(ctx context.Context) (ton.APIClientWrapped, error) {
			return tonChain.Client.WithRetry(3), nil
		}
		// init loader
		loader := tonlploader.NewTxLoader(lggr, clientProvider, 100)

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
				txs, err := loader.FetchTxsForAddress(ctx, blockRange, contractAddress)
				if err != nil {
					errorCh <- fmt.Errorf("failed to load transactions: %w", err)
					return
				}

				// 3. Get messages in transactions to see if there is event we're looking for
				events, err := extractEventMessage[T](txs, lggr)
				if err != nil {
					errorCh <- err
					return
				}

				// Send events to channel
				for _, event := range events {
					lggr.Infof("TON:FOUND EVENT: %+v", event)
					select {
					case ch <- event:
					case <-ctx.Done():
						return
					}
				}

				// Update last processed block
				lastProcessedBlock = newSeqNo
			}
		}
	}()

	return ch, errorCh
}

// getBlockRange creates a block range from lastProcessedBlock to current masterchain head
func getBlockRange(ctx context.Context, tonChain cldf_ton.Chain, lastProcessedBlock uint32) (*tonlptypes.BlockRange, uint32, error) {
	// lazy client provider
	clientProvider := func(_ context.Context) (ton.APIClientWrapped, error) {
		return tonChain.Client.WithRetry(3), nil
	}

	// 1. Get latest block number
	client, err := clientProvider(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get client: %w", err)
	}

	toBlock, err := client.CurrentMasterchainInfo(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get current masterchain info: %w", err)
	}

	// Skip if no new blocks
	if toBlock.SeqNo <= lastProcessedBlock {
		return nil, toBlock.SeqNo, nil
	}

	// 2. Create a block range
	var prevBlock *ton.BlockIDExt
	if lastProcessedBlock > 0 {
		// Get the previous block reference
		prevBlock, err = client.LookupBlock(ctx, toBlock.Workchain, toBlock.Shard, lastProcessedBlock)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to lookup previous block %d: %w", lastProcessedBlock, err)
		}
	}

	blockRange := &tonlptypes.BlockRange{
		Prev: prevBlock,
		To:   toBlock,
	}

	return blockRange, toBlock.SeqNo, nil
}

// extractEventMessage processes transactions to extract events of type T from external messages.
// Only processes events registered in CCIPEventTopics to filter out noise from OCR and other events.
func extractEventMessage[T any](txs []tonlptypes.TxWithBlock, lggr logger.Logger) ([]T, error) {
	var events []T

	for _, tx := range txs {
		if tx.Tx == nil {
			continue
		}
		if tx.Tx.IO.Out != nil {
			msgs, err := tx.Tx.IO.Out.ToSlice()
			if err != nil {
				continue
			}

			for _, msg := range msgs {
				if msg.MsgType == tlb.MsgTypeExternalOut {
					if extOut := msg.AsExternalOut(); extOut != nil {
						// decode event topic
						bucket := event.NewExtOutLogBucket(extOut.DestAddr())
						topic, err := bucket.DecodeEventTopic()
						if err != nil {
							lggr.Debugw("Failed to decode event topic", "error", err)
							continue
						}

						// skip events not in our assertion list (filters out OCR events, etc.)
						if _, isKnown := CCIPEventTopics[topic]; !isKnown {
							continue
						}

						bodyCell := extOut.Payload()
						if bodyCell == nil {
							lggr.Debugw("No payload in external out message")
							continue
						}

						lggr.Debugw("Processing event",
							"name", getEventName(topic),
							"topic", fmt.Sprintf("0x%08x", topic),
							"bits", bodyCell.BeginParse().BitsLeft(),
							"refs", bodyCell.RefsNum())

						// parse event using TLB - only succeeds if structure matches
						var event T
						err = tlb.LoadFromCell(&event, bodyCell.BeginParse())
						if err != nil {
							lggr.Debugw("Failed to parse event",
								"event", getEventName(topic),
								"error", err)
							continue
						}

						events = append(events, event)
					}
				}
			}
		}
	}

	return events, nil
}
