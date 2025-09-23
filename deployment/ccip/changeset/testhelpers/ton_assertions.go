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

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ccipocr3common "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	cldf_ton "github.com/smartcontractkit/chainlink-deployments-framework/chain/ton"

	"github.com/smartcontractkit/chainlink-ton/pkg/ccip/bindings/offramp"
	tonlploader "github.com/smartcontractkit/chainlink-ton/pkg/logpoller/backend/loader/account"
	tonlptypes "github.com/smartcontractkit/chainlink-ton/pkg/logpoller/types"
)

func ConfirmCommitWithExpectedSeqNumRangeTON(t *testing.T, srcSelector uint64, tonChain cldf_ton.Chain, offRampContract address.Address, expectedSeqNumRange ccipocr3common.SeqNumRange) (bool, error) {
	seenMessages := NewCommitReportTracker(srcSelector, expectedSeqNumRange)

	// Create prefixed logger for TON event assertion
	lggr := logger.Named(logger.Test(t), "TON_EVENT_ASSERTION")
	lggr.Infof("waiting for commit report from srcSelector=%d, expectedSeqNumRange=[%d, %d], timeout=%v",
		srcSelector, expectedSeqNumRange.Start(), expectedSeqNumRange.End(), tests.WaitTimeout(t))

	tonBlockTicker := time.NewTicker(2500 * time.Millisecond)
	// Start from a bit earlier to catch any events we might have missed
	// Use a lookback of about 50 blocks to ensure we don't miss commit events
	var startBlock uint32 = 0
	if currentBlock, err := tonChain.Client.CurrentMasterchainInfo(t.Context()); err == nil && currentBlock.SeqNo > 50 {
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

type TxLoader interface {
	FetchTxsForAddress(ctx context.Context, blockRange *tonlptypes.BlockRange, addr *address.Address) ([]tonlptypes.TxWithBlock, error)
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
				txs, err := loader.(TxLoader).FetchTxsForAddress(ctx, blockRange, contractAddress)
				if err != nil {
					errorCh <- fmt.Errorf("failed to load transactions: %w", err)
					return
				}

				// 3. Get messages in transactions to see if there is event we're looking for
				events, err := extractEventMessage[T](txs)
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

// extractEventMessage processes transactions to extract events of type T from external messages
func extractEventMessage[T any](txs []tonlptypes.TxWithBlock) ([]T, error) {
	var events []T

	for _, tx := range txs {
		if tx.Tx == nil {
			continue
		}
		if tx.Tx.IO.Out != nil {
			msgs, err := tx.Tx.IO.Out.ToSlice()
			if err != nil {
				// skip this tx
				continue
			}

			for _, msg := range msgs {
				if msg.MsgType == tlb.MsgTypeExternalOut {
					if extOut := msg.AsExternalOut(); extOut != nil {

						var event T
						err := tlb.LoadFromCell(&event, extOut.Body.BeginParse())
						if err == nil {
							events = append(events, event)
						}
					}
				}
			}
		}
	}

	return events, nil
}
