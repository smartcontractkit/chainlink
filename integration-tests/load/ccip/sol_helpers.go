package ccip

import (
	"context"
	"math"
	"slices"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"go.uber.org/atomic"

	solccip "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/ccip"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
)

func subscribeSolTransmitEvents(
	ctx context.Context,
	lggr logger.Logger,
	onrampAddress solana.PublicKey,
	otherChains []uint64,
	startSlot uint64,
	srcChainSel uint64,
	loadFinished chan struct{},
	client *solrpc.Client,
	wg *sync.WaitGroup,
	metricPipe chan messageData,
	finalSeqNrCommitChannels map[uint64]chan finalSeqNrReport,
	finalSeqNrExecChannels map[uint64]chan finalSeqNrReport,
) {
	defer wg.Done()
	lggr.Infow("starting solana chain transmit event subscriber for ",
		"srcChain", srcChainSel,
		"startSlot", startSlot,
	)

	seqNums := make(map[testhelpers.SourceDestPair]SeqNumRange)
	for _, cs := range otherChains {
		seqNums[testhelpers.SourceDestPair{
			SourceChainSelector: srcChainSel,
			DestChainSelector:   cs,
		}] = SeqNumRange{
			// we use the maxuint as a sentinel value here to ensure we always get the lowest possible seqnum
			Start: atomic.NewUint64(math.MaxUint64),
			End:   atomic.NewUint64(0),
		}
	}

	done := make(chan any)
	sink, errCh := testhelpers.SolEventEmitter[solccip.EventCCIPMessageSent](client, onrampAddress, "CCIPMessageSent", startSlot, done, time.NewTicker(15*time.Second))
	defer close(done)
	for {
		select {
		case err := <-errCh:
			lggr.Errorw("error in solana event emitter for subscribing transmit events",
				"srcChain", srcChainSel,
				"err", err)
			return

		case eventWithTxn := <-sink:
			event := eventWithTxn.Event
			lggr.Debugw("received solana transmit event for",
				"srcChain", srcChainSel,
				"destChain", event.DestinationChainSelector,
				"sequenceNumber", event.SequenceNumber,
				"timestamp", int64(*eventWithTxn.Txn.BlockTime))

			data := messageData{
				eventType: transmitted,
				srcDstSeqNum: srcDstSeqNum{
					src:    srcChainSel,
					dst:    event.DestinationChainSelector,
					seqNum: event.SequenceNumber,
				},
				timestamp: uint64(*eventWithTxn.Txn.BlockTime), //nolint:gosec // G115
			}

			metricPipe <- data
			csPair := testhelpers.SourceDestPair{
				SourceChainSelector: srcChainSel,
				DestChainSelector:   event.DestinationChainSelector,
			}
			// always store the lowest seen number as the start seq num
			if event.SequenceNumber < seqNums[csPair].Start.Load() {
				seqNums[csPair].Start.Store(event.SequenceNumber)
			}

			// always store the greatest sequence number we have seen as the maximum
			if event.SequenceNumber > seqNums[csPair].End.Load() {
				seqNums[csPair].End.Store(event.SequenceNumber)
			}
		case <-ctx.Done():
			lggr.Errorw("received context cancel signal for transmit watcher",
				"srcChain", srcChainSel)
			done <- struct{}{}
			return
		case <-loadFinished:
			lggr.Debugw("load finished, closing transmit watchers", "srcChainSel", srcChainSel)
			for csPair, seqNums := range seqNums {
				lggr.Infow("pushing finalized sequence numbers for ",
					"srcChainSelector", srcChainSel,
					"destChainSelector", csPair.DestChainSelector,
					"seqNums", seqNums)
				finalSeqNrCommitChannels[csPair.DestChainSelector] <- finalSeqNrReport{
					sourceChainSelector: csPair.SourceChainSelector,
					expectedSeqNrRange: ccipocr3.SeqNumRange{
						ccipocr3.SeqNum(seqNums.Start.Load()), ccipocr3.SeqNum(seqNums.End.Load()),
					},
				}

				finalSeqNrExecChannels[csPair.DestChainSelector] <- finalSeqNrReport{
					sourceChainSelector: csPair.SourceChainSelector,
					expectedSeqNrRange: ccipocr3.SeqNumRange{
						ccipocr3.SeqNum(seqNums.Start.Load()), ccipocr3.SeqNum(seqNums.End.Load()),
					},
				}
			}
			return
		}
	}
}

func subscribeSolCommitEvents(
	ctx context.Context,
	lggr logger.Logger,
	offrampAddress solana.PublicKey,
	srcChains []uint64,
	startSlot uint64,
	chainSelector uint64,
	client *solrpc.Client,
	finalSeqNrs chan finalSeqNrReport,
	wg *sync.WaitGroup,
	metricPipe chan messageData,
) {
	defer wg.Done()
	defer close(finalSeqNrs)

	lggr.Infow("starting solana commit event subscriber for ",
		"destChain", chainSelector,
		"startSlot", startSlot,
	)
	seenMessages := make(map[uint64][]uint64)
	expectedRange := make(map[uint64]ccipocr3.SeqNumRange)
	completedSrcChains := make(map[uint64]bool)
	for _, srcChain := range srcChains {
		// todo: seenMessages should hold a range to avoid hitting memory constraints
		seenMessages[srcChain] = make([]uint64, 0)
		completedSrcChains[srcChain] = false
	}

	done := make(chan any)
	sink, errCh := testhelpers.SolEventEmitter[solccip.EventCommitReportAccepted](client, offrampAddress, "CommitReportAccepted", startSlot, done, time.NewTicker(15*time.Second))
	defer close(done)

	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

	for {
		select {
		case err := <-errCh:
			lggr.Errorw("error in solana event emitter for subscribing commit events",
				"destChain", chainSelector,
				"seenMessages", seenMessages,
				"expectedRange", expectedRange,
				"completedSrcChains", completedSrcChains,
				"err", err)
			return

		case eventWithTx := <-sink:
			mr := eventWithTx.Event.Report
			if mr == nil {
				continue
			}

			lggr.Infow("Received solana commit report ",
				"sourceChain", mr.SourceChainSelector,
				"destChain", chainSelector,
				"minSeqNr", mr.MinSeqNr,
				"maxSeqNr", mr.MaxSeqNr,
				"timestamp", int64(*eventWithTx.Txn.BlockTime))

			// push metrics to state manager for eventual distribution to loki
			for i := mr.MinSeqNr; i <= mr.MaxSeqNr; i++ {
				data := messageData{
					eventType: committed,
					srcDstSeqNum: srcDstSeqNum{
						src:    mr.SourceChainSelector,
						dst:    chainSelector,
						seqNum: i,
					},
					timestamp: uint64(*eventWithTx.Txn.BlockTime), //nolint:gosec // G115
				}
				metricPipe <- data
				seenMessages[mr.SourceChainSelector] = append(seenMessages[mr.SourceChainSelector], i)
			}
		case <-ctx.Done():
			lggr.Errorw("timed out waiting for commit report",
				"destChain", chainSelector,
				"sourceChains", srcChains,
				"expectedSeqNumbers", expectedRange)
			done <- struct{}{}
			return

		case finalSeqNrUpdate, ok := <-finalSeqNrs:
			if finalSeqNrUpdate.expectedSeqNrRange.Start() == math.MaxUint64 || finalSeqNrUpdate.expectedSeqNrRange.End() == 0 {
				delete(completedSrcChains, finalSeqNrUpdate.sourceChainSelector)
				delete(seenMessages, finalSeqNrUpdate.sourceChainSelector)
			} else if ok {
				// only add to range if channel is still open
				expectedRange[finalSeqNrUpdate.sourceChainSelector] = finalSeqNrUpdate.expectedSeqNrRange
			}

		case <-ticker.C:
			lggr.Infow("ticking, checking committed events",
				"destChain", chainSelector,
				"seenMessages", seenMessages,
				"expectedRange", expectedRange,
				"completedSrcChains", completedSrcChains)
			for srcChain, seqNumRange := range expectedRange {
				// if this chain has already been marked as completed, skip
				if !completedSrcChains[srcChain] {
					// else, check if all expected sequence numbers have been seen
					// todo: We might need to modify if there are other non-load test txns on network
					if len(seenMessages[srcChain]) >= seqNumRange.Length() && slices.Contains(seenMessages[srcChain], uint64(seqNumRange.End())) {
						completedSrcChains[srcChain] = true
						delete(expectedRange, srcChain)
						delete(seenMessages, srcChain)
						lggr.Infow("committed all sequence numbers for ",
							"sourceChain", srcChain,
							"destChain", chainSelector)
					}
				}
			}
			// if all chains have hit expected sequence numbers, return
			// we could instead push complete chains to an incrementer and compare size
			allComplete := true
			for c := range completedSrcChains {
				if !completedSrcChains[c] {
					allComplete = false
					break
				}
			}
			if allComplete {
				lggr.Infof("received commits from expected source chains for all expected sequence numbers to chainSelector %d", chainSelector)
				return
			}
		}
	}
}

func subscribeSolExecutionEvents(
	ctx context.Context,
	lggr logger.Logger,
	offrampAddress solana.PublicKey,
	srcChains []uint64,
	startSlot uint64,
	chainSelector uint64,
	client *solrpc.Client,
	finalSeqNrs chan finalSeqNrReport,
	wg *sync.WaitGroup,
	metricPipe chan messageData,
) {
	defer wg.Done()
	defer close(finalSeqNrs)

	lggr.Infow("starting solana chain execution event subscriber for ",
		"destChain", chainSelector,
		"startblock", startSlot,
	)
	seenMessages := make(map[uint64][]uint64)
	expectedRange := make(map[uint64]ccipocr3.SeqNumRange)
	completedSrcChains := make(map[uint64]bool)
	for _, srcChain := range srcChains {
		seenMessages[srcChain] = make([]uint64, 0)
		completedSrcChains[srcChain] = false
	}
	done := make(chan any)
	sink, errCh := testhelpers.SolEventEmitter[solccip.EventExecutionStateChanged](client, offrampAddress, "ExecutionStateChanged", startSlot, done, time.NewTicker(15*time.Second))
	defer close(done)

	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

	for {
		select {
		case err := <-errCh:
			lggr.Errorw("error in solana event emitter for execution events",
				"destChain", chainSelector,
				"seenMessages", seenMessages,
				"expectedRange", expectedRange,
				"completedSrcChains", completedSrcChains,
				"err", err)
			return

		case eventWithTxn := <-sink:
			event := eventWithTxn.Event
			if event.State.String() != "Success" {
				continue
			}
			lggr.Debugw("received solana execution event for",
				"destChain", chainSelector,
				"sourceChain", event.SourceChainSelector,
				"sequenceNumber", event.SequenceNumber,
				"timestamp", uint64(*eventWithTxn.Txn.BlockTime)) //nolint:gosec // G115
			// push metrics to loki here
			data := messageData{
				eventType: executed,
				srcDstSeqNum: srcDstSeqNum{
					src:    event.SourceChainSelector,
					dst:    chainSelector,
					seqNum: event.SequenceNumber,
				},
				timestamp: uint64(*eventWithTxn.Txn.BlockTime), //nolint:gosec // G115
			}
			metricPipe <- data
			seenMessages[event.SourceChainSelector] = append(seenMessages[event.SourceChainSelector], event.SequenceNumber)

		case <-ctx.Done():
			lggr.Errorw("timed out waiting for execution event",
				"destChain", chainSelector,
				"sourceChains", srcChains,
				"expectedSeqNumbers", expectedRange,
				"seenMessages", seenMessages,
				"completedSrcChains", completedSrcChains)
			done <- struct{}{}
			return

		case finalSeqNrUpdate := <-finalSeqNrs:
			if finalSeqNrUpdate.expectedSeqNrRange.Start() == math.MaxUint64 || finalSeqNrUpdate.expectedSeqNrRange.End() == 0 {
				delete(completedSrcChains, finalSeqNrUpdate.sourceChainSelector)
				delete(seenMessages, finalSeqNrUpdate.sourceChainSelector)
			} else {
				expectedRange[finalSeqNrUpdate.sourceChainSelector] = finalSeqNrUpdate.expectedSeqNrRange
			}

		case <-ticker.C:
			lggr.Infow("ticking, checking executed events",
				"destChain", chainSelector,
				"seenMessages", seenMessages,
				"expectedRange", expectedRange,
				"completedSrcChains", completedSrcChains)

			for srcChain, seqNumRange := range expectedRange {
				// if this chain has already been marked as completed, skip
				if !completedSrcChains[srcChain] {
					// else, check if all expected sequence numbers have been seen
					if len(seenMessages[srcChain]) >= seqNumRange.Length() && slices.Contains(seenMessages[srcChain], uint64(seqNumRange.End())) {
						completedSrcChains[srcChain] = true
						lggr.Infow("executed all sequence numbers for ",
							"destChain", chainSelector,
							"sourceChain", srcChain,
							"seqNumRange", seqNumRange)
					}
				}
			}
			// if all chains have hit expected sequence numbers, return
			allComplete := true
			for c := range completedSrcChains {
				if !completedSrcChains[c] {
					allComplete = false
					break
				}
			}
			if allComplete {
				lggr.Infow("all messages have been executed for all expected sequence numbers",
					"destChain", chainSelector)
				return
			}
		}
	}
}
