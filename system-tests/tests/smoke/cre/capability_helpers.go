package cre

import (
	"bufio"
	"context"
	"errors"
	"io"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	evmreadcontracts "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evm/evmread/contracts"
)

// chainlink-common/pkg/capabilities/base_trigger.go logs this when AckEvent is called.
var triggerEventACKLogPattern = regexp.MustCompile(`Event ACK`)

// verifyTriggerEventACKLogs starts parallel readers (one per container stream) that scan for BaseTrigger
// Event ACK lines. Call the returned cleanup after the test step finishes to cancel scanners and close
// readers. The returned atomic is set to true when any container reports a matching line.
func verifyTriggerEventACKLogs(ctx context.Context, lggr zerolog.Logger, logStreams map[string]io.ReadCloser) (*atomic.Bool, func()) {
	found := &atomic.Bool{}
	if len(logStreams) == 0 {
		lggr.Error().Msg("verifyTriggerEventACKLogs: no container log streams to scan")
		return found, func() {}
	}

	scanCtx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	// Snapshot readers so cleanup can close them — closing the reader is what unblocks a goroutine
	// stuck in scanner.Scan() on a Follow=true stream; context cancel alone is not enough.
	readers := make([]io.ReadCloser, 0, len(logStreams))
	for containerName, r := range logStreams {
		wg.Add(1)
		readers = append(readers, r)
		go func(name string, reader io.ReadCloser) {
			defer wg.Done()
			scanOneContainerForTriggerEventACK(scanCtx, cancel, lggr, name, reader, found)
		}(containerName, r)
	}
	return found, func() {
		cancel()
		for _, r := range readers {
			_ = r.Close()
		}
		wg.Wait()
	}
}

func scanOneContainerForTriggerEventACK(ctx context.Context, cancel context.CancelFunc, lggr zerolog.Logger, containerName string, reader io.Reader, found *atomic.Bool) {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
		}
		line := scanner.Text()
		if triggerEventACKLogPattern.MatchString(line) {
			lggr.Info().Str("container", containerName).Str("line", line).Msg("detected BaseTrigger Event ACK in container logs")
			found.Store(true)
			cancel()
			return
		}
	}
	if err := scanner.Err(); err != nil && !isExpectedLogStreamCloseErr(ctx, err) {
		lggr.Error().Err(err).Str("container", containerName).Msg("error reading container logs while scanning for Event ACK")
	}
}

// isExpectedLogStreamCloseErr returns true when a follower goroutine exits because cleanup
// closed the Docker log stream or cancelled the scan context after another container matched.
func isExpectedLogStreamCloseErr(ctx context.Context, err error) bool {
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrClosedPipe) || ctx.Err() != nil {
		return true
	}
	msg := err.Error()
	return strings.Contains(msg, "use of closed network connection") ||
		strings.Contains(msg, "file already closed")
}

// startTriggerEventACKLogWatch streams live container logs and scans for BaseTrigger Event ACK lines.
// Call the returned cleanup after the test step finishes.
func startTriggerEventACKLogWatch(t *testing.T, lggr zerolog.Logger) (*atomic.Bool, func()) {
	t.Helper()
	// Follow=true is required so we see "Event ACK" lines that are logged AFTER we start streaming
	// (they only happen after events are emitted further down). Tail="0" skips historical lines so
	// we don't match ACKs from previous test iterations or earlier workflow activity.
	logsOpts := framework.CTFContainersLogsOpts()
	logsOpts.Follow = true
	logsOpts.Tail = "0"
	logstream, err := framework.StreamContainerLogs(framework.CTFContainersListOpts(), logsOpts)
	require.NoError(t, err, "failed to stream container logs for Event ACK check")
	return verifyTriggerEventACKLogs(t.Context(), lggr, logstream)
}

// requireTriggerEventACKLog asserts that a BaseTrigger Event ACK line appeared in container logs.
func requireTriggerEventACKLog(t *testing.T, lggr zerolog.Logger, singleAckFound *atomic.Bool) {
	t.Helper()
	require.Eventually(t, func() bool {
		return singleAckFound.Load()
	}, 4*time.Minute, 500*time.Millisecond,
		"expected BaseTrigger Event ACK log line in container logs (BaseTriggerCapability.AckEvent logs msg Event ACK)")
	lggr.Info().Msg("found BaseTrigger Event ACK log")
}

func startEVMLogTriggerEventEmitter(
	ctx context.Context,
	t *testing.T,
	logger zerolog.Logger,
	chainID string,
	bcOutput blockchains.Blockchain,
	msgEmitter *evmreadcontracts.MessageEmitter,
	message string,
) {
	t.Helper()

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		defer ticker.Stop()

		var emittedEventCount int64
		for {
			logger.Info().Int64("eventCount", emittedEventCount).Str("chainID", chainID).Msg("Emitting EVM log trigger event")
			blockNumber := emitEvent(t, logger, chainID, bcOutput, msgEmitter, message)
			logger.Info().Str("chainID", chainID).Uint64("blockNumber", blockNumber).Msg("EVM log trigger event emitted")
			emittedEventCount++

			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}
	}()
}

func emitEvent(t *testing.T, lggr zerolog.Logger, chainID string, bcOutput blockchains.Blockchain, msgEmitter *evmreadcontracts.MessageEmitter, message string) uint64 {
	lggr.Info().Msgf("Emitting event to be picked up by workflow for chain '%s'", chainID)
	sethClient := bcOutput.(*evm.Blockchain).SethClient
	emittingTx, err := msgEmitter.EmitMessage(sethClient.NewTXOpts(), message)
	if err != nil {
		lggr.Info().Msgf("Failed to emit transaction for chain '%s': %v", chainID, err)
		return 0
	}

	emittingReceipt, err := sethClient.WaitMined(t.Context(), lggr, sethClient.Client, emittingTx)
	if err != nil {
		lggr.Info().Msgf("Failed to emit receipt for chain '%s': %v", chainID, err)
		return 0
	}
	lggr.Info().Msgf("Transaction for chain '%s' mined at '%d' with emitted message %q", chainID, emittingReceipt.BlockNumber.Uint64(), message)
	return emittingReceipt.BlockNumber.Uint64()
}
