package ccip

import (
	"context"
	"errors"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

const (
	transmitted = iota
	committed
	executed
	LokiLoadLabel  = "ccip_load_test"
	ErrLokiClient  = "failed to create Loki client for monitoring"
	ErrLokiPush    = "failed to push metrics to Loki"
	tickerDuration = 30 * time.Second
)

// todo: Have a different struct for commit/exec?
type LokiMetric struct {
	EventType      int       `json:"event_type"`
	Timestamp      time.Time `json:"timestamp"`
	GasUsed        uint64    `json:"gas_used"`
	SequenceNumber uint64    `json:"sequence_number"`
}

func SendMetricsToLoki(l logger.Logger, lc *wasp.LokiClient, updatedLabels map[string]string, metrics *LokiMetric) {
	if err := lc.HandleStruct(wasp.LabelsMapToModel(updatedLabels), time.Now(), metrics); err != nil {
		l.Error(ErrLokiPush)
	}
}

func setLokiLabels(src, dst uint64) (map[string]string, error) {
	srcChainID, err := chainselectors.GetChainIDFromSelector(src)
	if err != nil {
		return nil, err
	}
	dstChainID, err := chainselectors.GetChainIDFromSelector(dst)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"sourceEvmChainId":    srcChainID,
		"sourceSelector":      strconv.FormatUint(src, 10),
		"destEvmChainId":      dstChainID,
		"destinationSelector": strconv.FormatUint(dst, 10),
		"testType":            LokiLoadLabel,
	}, nil
}

type finalSeqNrReport struct {
	sourceChainSelector uint64
	expectedSeqNrRange  ccipocr3.SeqNumRange
}

func subscribeDeferredCommitEvents(
	ctx context.Context,
	lggr logger.Logger,
	offRamp offramp.OffRampInterface,
	srcChains []uint64,
	startBlock *uint64,
	loki *wasp.LokiClient,
	chainSelector uint64,
	client deployment.OnchainClient,
	finalSeqNrs chan finalSeqNrReport,
	errChan chan error,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	lggr.Infow("starting commit event subscriber for ",
		"destChainSelector", chainSelector,
		"offRamp", offRamp.Address().String(),
		"startblock", startBlock,
	)
	seenMessages := make(map[uint64][]uint64)
	expectedRange := make(map[uint64]ccipocr3.SeqNumRange)
	completedSrcChains := make(map[uint64]bool)
	for _, srcChain := range srcChains {
		// todo: seenMessages should hold a range to avoid hitting memory constraints
		seenMessages[srcChain] = make([]uint64, 0)
		completedSrcChains[srcChain] = false
	}

	sink := make(chan *offramp.OffRampCommitReportAccepted)
	subscription, err := offRamp.WatchCommitReportAccepted(&bind.WatchOpts{
		Context: context.Background(),
		Start:   startBlock,
	}, sink)
	if err != nil {
		errChan <- err
		return
	}
	defer subscription.Unsubscribe()
	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

	for {
		select {
		case subErr := <-subscription.Err():
			errChan <- subErr
			return
		case report := <-sink:
			if len(report.BlessedMerkleRoots)+len(report.UnblessedMerkleRoots) > 0 {
				for _, mr := range append(report.BlessedMerkleRoots, report.UnblessedMerkleRoots...) {
					lggr.Infow("Received commit report ",
						"sourceChain", mr.SourceChainSelector,
						"offRamp", offRamp.Address().String(),
						"minSeqNr", mr.MinSeqNr,
						"maxSeqNr", mr.MaxSeqNr)

					lokiLabels, err := setLokiLabels(mr.SourceChainSelector, chainSelector)
					if err != nil {
						errChan <- err
						// don't return here, we still want to push metrics to loki
					}
					// push metrics to loki here
					for i := mr.MinSeqNr; i <= mr.MaxSeqNr; i++ {
						blockNum := report.Raw.BlockNumber
						header, err := client.HeaderByNumber(ctx, new(big.Int).SetUint64(blockNum))
						if err != nil {
							errChan <- err
						}
						timestamp := time.Unix(int64(header.Time), 0) //nolint:gosec // disable G115
						SendMetricsToLoki(lggr, loki, lokiLabels, &LokiMetric{
							EventType:      committed,
							Timestamp:      timestamp,
							SequenceNumber: i,
						})
						seenMessages[mr.SourceChainSelector] = append(seenMessages[mr.SourceChainSelector], i)
					}
				}
			}
		case <-ctx.Done():
			lggr.Errorw("timed out waiting for commit report",
				"offRamp", offRamp.Address().String(),
				"sourceChains", srcChains,
				"expectedSeqNumbers", expectedRange)
			errChan <- errors.New("timed out waiting for commit report")
			return

		case finalSeqNrUpdate, ok := <-finalSeqNrs:
			if ok { // only add to range if channel is still open
				expectedRange[finalSeqNrUpdate.sourceChainSelector] = finalSeqNrUpdate.expectedSeqNrRange
			}

		case <-ticker.C:
			lggr.Infow("ticking, checking committed events",
				"seenMessages", seenMessages,
				"expectedRange", expectedRange,
				"completedSrcChains", completedSrcChains)
			for srcChain, seqNumRange := range expectedRange {
				// if this chain has already been marked as completed, skip
				if !completedSrcChains[srcChain] {
					// else, check if all expected sequence numbers have been seen
					// todo: We might need to modify if there are other non-load test txns on network
					if len(seenMessages[srcChain]) >= seqNumRange.Length() {
						completedSrcChains[srcChain] = true
						lggr.Infow("committed all sequence numbers for ",
							"sourceChain", srcChain,
							"destChain", chainSelector,
							"seqNumRange", seqNumRange)
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
				lggr.Infow("all chains have committed all expected sequence numbers")
				return
			}
		}
	}
}

func subscribeExecutionEvents(
	ctx context.Context,
	lggr logger.Logger,
	offRamp offramp.OffRampInterface,
	srcChains []uint64,
	startBlock *uint64,
	loki *wasp.LokiClient,
	chainSelector uint64,
	client deployment.OnchainClient,
	finalSeqNrs chan finalSeqNrReport,
	errChan chan error,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	defer close(errChan)

	lggr.Infow("starting execution event subscriber for ",
		"destChainSelector", chainSelector,
		"offRamp", offRamp.Address().String(),
		"startblock", startBlock,
	)
	seenMessages := make(map[uint64][]uint64)
	expectedRange := make(map[uint64]ccipocr3.SeqNumRange)
	completedSrcChains := make(map[uint64]bool)
	for _, srcChain := range srcChains {
		seenMessages[srcChain] = make([]uint64, 0)
		completedSrcChains[srcChain] = false
	}

	sink := make(chan *offramp.OffRampExecutionStateChanged)
	subscription, err := offRamp.WatchExecutionStateChanged(&bind.WatchOpts{
		Context: context.Background(),
		Start:   startBlock,
	}, sink, srcChains, nil, nil)
	if err != nil {
		errChan <- err
		return
	}
	defer subscription.Unsubscribe()
	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

	for {
		select {
		case subErr := <-subscription.Err():
			lggr.Errorw("error in execution subscription",
				"err", subErr)
			errChan <- subErr
			return
		case event := <-sink:
			lggr.Debugw("received execution event",
				"event", event)
			lokiLabels, err := setLokiLabels(event.SourceChainSelector, chainSelector)
			if err != nil {
				errChan <- err
				// don't return here, we still want to push metrics to loki
			}
			// push metrics to loki here
			blockNum := event.Raw.BlockNumber
			header, err := client.HeaderByNumber(ctx, new(big.Int).SetUint64(blockNum))
			if err != nil {
				errChan <- err
			}
			timestamp := time.Unix(int64(header.Time), 0) //nolint:gosec // disable G115
			SendMetricsToLoki(lggr, loki, lokiLabels, &LokiMetric{
				EventType:      executed,
				Timestamp:      timestamp,
				GasUsed:        event.GasUsed.Uint64(),
				SequenceNumber: event.SequenceNumber,
			})
			seenMessages[event.SourceChainSelector] = append(seenMessages[event.SourceChainSelector], event.SequenceNumber)

		case <-ctx.Done():
			lggr.Errorw("timed out waiting for execution event",
				"offRamp", offRamp.Address().String(),
				"sourceChains", srcChains,
				"expectedSeqNumbers", expectedRange,
				"seenMessages", seenMessages,
				"completedSrcChains", completedSrcChains)
			errChan <- errors.New("timed out waiting for execution event")
			return

		case finalSeqNrUpdate := <-finalSeqNrs:
			expectedRange[finalSeqNrUpdate.sourceChainSelector] = finalSeqNrUpdate.expectedSeqNrRange

		case <-ticker.C:
			lggr.Infow("ticking, checking executed events",
				"seenMessages", seenMessages,
				"expectedRange", expectedRange,
				"completedSrcChains", completedSrcChains)

			for srcChain, seqNumRange := range expectedRange {
				// if this chain has already been marked as completed, skip
				if !completedSrcChains[srcChain] {
					// else, check if all expected sequence numbers have been seen
					if len(seenMessages[srcChain]) >= seqNumRange.Length() {
						completedSrcChains[srcChain] = true
						lggr.Infow("executed all sequence numbers for ",
							"sourceChain", srcChain,
							"destChain", chainSelector,
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
				lggr.Infow("all chains have executed all expected sequence numbers",
					"expectedSeqNumbers", expectedRange)
				return
			}
		}
	}
}
