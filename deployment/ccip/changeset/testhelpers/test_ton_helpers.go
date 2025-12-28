package testhelpers

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
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
	"github.com/smartcontractkit/chainlink-ton/pkg/ton/hash"

	tonlogpoller "github.com/smartcontractkit/chainlink-ton/pkg/logpoller"
	tonlploader "github.com/smartcontractkit/chainlink-ton/pkg/logpoller/loader"
	tonlptypes "github.com/smartcontractkit/chainlink-ton/pkg/logpoller/models"
	tonlpquery "github.com/smartcontractkit/chainlink-ton/pkg/logpoller/query"
	tonlpstore "github.com/smartcontractkit/chainlink-ton/pkg/logpoller/store/memory"
)

var (
	// ErrTimeout is returned when event subscription times out
	ErrTimeout = errors.New("timed out waiting for events")
)

// TON blockchain polling configuration
const (
	clientRetries       = 3                      // Number of retries for TON client operations
	queryInterval       = 500 * time.Millisecond // How often to query logpoller for new events
	progressLogInterval = 5 * time.Second        // How often to log "still waiting" progress updates
)

// setupLogPoller creates and starts a logpoller service with in-memory stores for the given contract and event.
func setupLogPoller(
	ctx context.Context,
	t *testing.T,
	lggr logger.Logger,
	tonChain cldf_ton.Chain,
	contract *address.Address,
	eventName string,
) tonlogpoller.Service {
	chainID := strconv.FormatUint(tonChain.Selector, 10)
	clientProvider := func(ctx context.Context) (ton.APIClientWrapped, error) {
		return tonChain.Client.WithRetry(clientRetries), nil
	}

	// Create logpoller with in-memory stores for testing
	service, err := tonlogpoller.NewService(lggr, chainID, clientProvider, &tonlogpoller.ServiceOptions{
		Config:      tonlogpoller.DefaultConfigSet,
		FilterStore: tonlpstore.NewFilterStore(chainID, lggr),
		TxLoader:    tonlploader.New(lggr, clientProvider),
		LogStore:    tonlpstore.NewLogStore(chainID, lggr),
	})
	require.NoError(t, err)

	_, err = service.RegisterFilter(ctx, tonlptypes.Filter{
		Name:     fmt.Sprintf("%s-%s", contract.String(), eventName),
		Address:  contract,
		EventSig: hash.CRC32(eventName),
		MsgType:  tlb.MsgTypeExternalOut,
	})
	require.NoError(t, err)
	require.NoError(t, service.Start(ctx))

	return service
}

// waitForTONEvent sets up a logpoller and waits for events matching the given criteria.
// Handles service lifecycle and common error patterns.
func waitForTONEvent[T any](
	t *testing.T,
	tonChain cldf_ton.Chain,
	offRamp *address.Address,
	eventName string,
	loggerName string,
	processEvent func(lggr logger.Logger, event tonlptypes.TypedLog[T]) (done bool, err error),
) error {
	ctx := t.Context()
	lggr := logger.Named(logger.Test(t), loggerName)

	service := setupLogPoller(ctx, t, lggr, tonChain, offRamp, eventName)
	defer service.Close()

	eventSig := hash.CRC32(eventName)
	deadline := time.Now().Add(tests.WaitTimeout(t))
	ticker := time.NewTicker(queryInterval)
	defer ticker.Stop()

	progressTicker := time.NewTicker(progressLogInterval)
	defer progressTicker.Stop()

	startTime := time.Now()
	seenEvents := make(map[string]bool)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-progressTicker.C:
			lggr.Infow("Still waiting",
				"eventName", eventName,
				"elapsed", time.Since(startTime).Round(time.Second).String())

		case <-ticker.C:
			if time.Now().After(deadline) {
				return ErrTimeout
			}

			logs, _, _, err := service.NewQuery().
				WithSource(offRamp).
				WithEventSig(eventSig).
				Execute(ctx)
			if err != nil {
				lggr.Warnw("Failed to query logs", "error", err)
				continue
			}

			events, err := tonlpquery.DecodedLogs[T](logs)
			if err != nil {
				lggr.Warnw("Failed to decode logs", "error", err)
				continue
			}

			for _, event := range events {
				eventKey := fmt.Sprintf("%d-%d", event.TxLT, event.MsgIndex)
				if seenEvents[eventKey] {
					continue
				}
				seenEvents[eventKey] = true

				done, err := processEvent(lggr, event)
				if err != nil {
					return err
				}
				if done {
					return nil
				}
			}
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
	tracker := NewCommitReportTracker(srcChainSelector, expectedSeqNums)
	reportsProcessed := 0

	err := waitForTONEvent(t, tonChain, &offRamp, consts.EventNameCommitReportAccepted, "TON_EVENT_ASSERTION:COMMIT",
		func(lggr logger.Logger, event tonlptypes.TypedLog[offramp.CommitReportAccepted]) (bool, error) {
			mr := event.TypedData.MerkleRoot
			if mr == nil {
				return false, nil // Skip price-only updates
			}

			reportsProcessed++
			require.Equal(t, srcChainSelector, mr.SourceChainSelector, "Commit report source chain mismatch")
			lggr.Infow("Received commit", "seqNums", fmt.Sprintf("[%d, %d]", mr.MinSeqNr, mr.MaxSeqNr))

			tracker.visitCommitReport(srcChainSelector, mr.MinSeqNr, mr.MaxSeqNr)

			// Check if all messages committed (single or multiple reports)
			if (uint64(expectedSeqNums.Start()) >= mr.MinSeqNr && uint64(expectedSeqNums.End()) <= mr.MaxSeqNr) ||
				tracker.allCommited(srcChainSelector) {
				t.Logf("All sequence numbers committed [%d, %d]", expectedSeqNums.Start(), expectedSeqNums.End())
				return true, nil
			}

			return false, nil
		})

	if errors.Is(err, ErrTimeout) {
		return false, fmt.Errorf("timed out waiting for commit on chain %d from source %d, seq nums %s (%d reports processed): %w",
			tonChain.Selector, srcChainSelector, expectedSeqNums.String(), reportsProcessed, err)
	}
	return err == nil, err
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

	executionStates := make(map[uint64]int)
	pending := make(map[uint64]bool)
	for _, seqNum := range expectedSeqNums {
		pending[seqNum] = true
	}
	eventsProcessed := 0

	err := waitForTONEvent(t, tonChain, &offRamp, consts.EventNameExecutionStateChanged, "TON_EVENT_ASSERTION:EXEC",
		func(lggr logger.Logger, event tonlptypes.TypedLog[offramp.ExecutionStateChanged]) (bool, error) {
			exec := event.TypedData

			if exec.SourceChainSelector != srcChainSelector || (!pending[exec.SequenceNumber] && executionStates[exec.SequenceNumber] == 0) {
				return false, nil
			}

			eventsProcessed++

			switch exec.State {
			case EXECUTION_STATE_INPROGRESS:
				return false, nil

			case EXECUTION_STATE_FAILURE:
				lggr.Errorw("Execution failed", "sequenceNumber", exec.SequenceNumber, "messageID", hex.EncodeToString(exec.MessageID))
				return false, fmt.Errorf("execution failed for seq %d on chain %d, message ID: %x",
					exec.SequenceNumber, exec.SourceChainSelector, exec.MessageID)

			case EXECUTION_STATE_SUCCESS:
				executionStates[exec.SequenceNumber] = int(exec.State)
				delete(pending, exec.SequenceNumber)
				lggr.Infow("Execution successful", "sequenceNumber", exec.SequenceNumber, "remaining", len(pending))

				if len(pending) == 0 {
					t.Logf("All sequence numbers executed: %v", expectedSeqNums)
					return true, nil
				}

			default:
				lggr.Warnw("Unknown execution state", "state", exec.State, "sequenceNumber", exec.SequenceNumber)
			}

			return false, nil
		})

	if errors.Is(err, ErrTimeout) {
		missing := make([]uint64, 0, len(pending))
		for seqNum := range pending {
			missing = append(missing, seqNum)
		}
		return executionStates, fmt.Errorf("timed out waiting for execution on chain %d from source %d, missing: %v (%d events, %d successful): %w",
			tonChain.Selector, srcChainSelector, missing, eventsProcessed, len(executionStates), err)
	}
	return executionStates, err
}
