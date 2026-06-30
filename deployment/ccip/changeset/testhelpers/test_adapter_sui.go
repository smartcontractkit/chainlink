package testhelpers

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	sui_module_offramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_offramp/offramp"
	sui_ccip_offramp "github.com/smartcontractkit/chainlink-sui/bindings/packages/offramp"
	cslclient "github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_sui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	suistate "github.com/smartcontractkit/chainlink-sui/deployment"
)

type SuiAdapter struct {
	state suistate.CCIPChainState
	cldf_sui.Chain
}

// offRampOriginalPkgID returns the original (V1) package ID, which must be used
// for event queries. In Sui, struct types (including events) always carry the
// original defining package's ID regardless of which upgraded version emitted them.
func (a *SuiAdapter) offRampOriginalPkgID() string {
	return a.state.OffRampAddress
}

func NewSuiAdapter(chain cldf.BlockChain, env deployment.Environment) Adapter {
	c, ok := chain.(cldf_sui.Chain)
	if !ok {
		panic(fmt.Sprintf("invalid chain type: %T", chain))
	}
	state, err := suistate.LoadOnchainStatesui(env)
	if err != nil {
		panic(fmt.Sprintf("failed to load onchain state: %T", err))
	}
	// NOTE: since this returns a copy, adapters shouldn't be constructed until everything is deployed
	s := state[c.ChainSelector()]
	return &SuiAdapter{
		state: s,
		Chain: c,
	}
}

func (a *SuiAdapter) BuildMessage(components MessageComponents) (any, error) {
	return SuiSendRequest{
		Data:      components.Data,
		Receiver:  common.LeftPadBytes(components.Receiver, 32),
		ExtraArgs: components.ExtraArgs,
		FeeToken:  components.FeeToken,
	}, nil
}

func (a *SuiAdapter) NativeFeeToken() string {
	// TODO:
	return ""
}

func (a *SuiAdapter) GetExtraArgs(receiver []byte, sourceFamily string, opts ...ExtraArgOpt) ([]byte, error) {
	return nil, nil
}

func (a *SuiAdapter) GetInboundNonce(ctx context.Context, sender []byte, srcSel uint64) (uint64, error) {
	return 0, errors.ErrUnsupported
}

func (a *SuiAdapter) ValidateCommit(t *testing.T, sourceSelector uint64, startBlock *uint64, seqNumRange ccipocr3.SeqNumRange) {
	_, err := confirmCommitWithExpectedSeqNumRangeSui(
		t,
		sourceSelector,
		a.Chain,
		a.offRampOriginalPkgID(),
		startBlock,
		seqNumRange,
		true,
	)
	require.NoError(t, err)
}

func (a *SuiAdapter) ValidateExec(t *testing.T, sourceSelector uint64, startBlock *uint64, seqNrs []uint64) (executionStates map[uint64]int) {
	executionStates, err := confirmExecWithExpectedSeqNrsSui(
		t,
		sourceSelector,
		a.Chain,
		a.offRampOriginalPkgID(),
		startBlock,
		seqNrs,
	)
	require.NoError(t, err)
	return executionStates
}

// suiEventCheckpointBackfill is how many checkpoints behind the current tip the
// emitter starts scanning when it subscribes, so events emitted just before the
// subscription (e.g. a commit that landed between send and validation) are not missed.
const suiEventCheckpointBackfill = uint64(50)

func SuiEventEmitter[T any](
	t *testing.T,
	client cslclient.SuiPTBClient,
	packageID, moduleName, event string,
	done chan any,
) (<-chan struct {
	Event   T
	Version string
}, <-chan error) {
	startTime := time.Now()
	t.Logf("[DEBUG] SuiEventEmitter: Starting at %s - polling checkpoints for events", startTime.Format(time.RFC3339))
	ch := make(chan struct {
		Event   T
		Version string
	}, 200)
	errChan := make(chan error)

	// The gRPC client does not implement cursor-based event queries (QueryEvents is
	// "pending gRPC migration"), so we poll checkpoints and filter their events by the
	// fully-qualified event type, mirroring the relayer's chain_poller.
	eventType := fmt.Sprintf("%s::%s::%s", packageID, moduleName, event)

	go func() {
		defer close(ch)
		defer close(errChan)

		ctx := t.Context()

		emitErr := func(err error) {
			select {
			case errChan <- err:
			case <-done:
			}
		}

		// Seed the scan position from the current chain tip, backfilling a small window.
		latest, err := client.GetLatestCheckpoint(ctx)
		if err != nil {
			emitErr(fmt.Errorf("failed to get latest checkpoint: %w", err))
			return
		}
		var nextSeq uint64
		if tip := latest.GetSequenceNumber(); tip > suiEventCheckpointBackfill {
			nextSeq = tip - suiEventCheckpointBackfill
		}

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				t.Logf("[DEBUG] SuiEventEmitter: Stopping due to done signal")
				return
			default:
			}

			latest, err := client.GetLatestCheckpoint(ctx)
			if err != nil {
				emitErr(fmt.Errorf("failed to get latest checkpoint: %w", err))
				return
			}
			tip := latest.GetSequenceNumber()

			for seq := nextSeq; seq <= tip; seq++ {
				select {
				case <-done:
					return
				default:
				}

				data, err := client.GetCheckpointData(ctx, seq)
				if err != nil {
					if isSuiCheckpointNotFound(err) {
						// Tip advanced past a checkpoint not yet available; retry next tick.
						break
					}
					emitErr(err)
					return
				}

				for _, tx := range data.Transactions {
					for _, ev := range tx.GetEvents().GetEvents() {
						// Sui event struct types always carry the original defining
						// package's ID, so match on the fully-qualified handle.
						qualified := strings.Join([]string{ev.GetPackageId(), ev.GetModule(), ev.GetEventType()}, "::")
						if qualified != eventType {
							continue
						}
						if ev.GetJson() == nil {
							continue
						}

						var out T
						if err := codec.DecodeSuiJsonValue(ev.GetJson().AsInterface(), &out); err != nil {
							t.Logf("[DEBUG] SuiEventEmitter: Decode error at checkpoint %d: %v (skipping)", seq, err)
							continue
						}

						eventData := struct {
							Event   T
							Version string
						}{
							Event:   out,
							Version: strconv.FormatUint(seq, 10),
						}

						select {
						case ch <- eventData:
							t.Logf("[DEBUG] SuiEventEmitter: Sent %s event from checkpoint %d", eventType, seq)
						case <-done:
							t.Logf("[DEBUG] SuiEventEmitter: Stopping due to done signal during send")
							return
						default:
							t.Logf("[WARNING] SuiEventEmitter: Channel full, dropping event at checkpoint %d", seq)
						}
					}
				}

				nextSeq = seq + 1
			}

			select {
			case <-done:
				t.Logf("[DEBUG] SuiEventEmitter: Stopping due to done signal in ticker loop")
				return
			case <-ticker.C:
				continue
			}
		}
	}()
	return ch, errChan
}

// isSuiCheckpointNotFound reports whether err indicates a checkpoint that is not yet
// available on the fullnode (the tip can advance past the latest indexed checkpoint).
func isSuiCheckpointNotFound(err error) bool {
	for err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return true
		}
		err = errors.Unwrap(err)
	}
	return false
}

func confirmCommitWithExpectedSeqNumRangeSui(
	t *testing.T,
	srcSelector uint64,
	dest cldf_sui.Chain,
	offRampAddress string,
	startVersion *uint64,
	expectedSeqNumRange ccipocr3.SeqNumRange,
	enforceSingleCommit bool,
) (any, error) {
	// Bound the offRamp
	boundOffRamp, err := sui_ccip_offramp.NewOfframp(offRampAddress, dest.Client)
	require.NoError(t, err)

	done := make(chan any)
	defer close(done)
	sink, errChan := SuiEventEmitter[sui_module_offramp.CommitReportAccepted](t, dest.Client, boundOffRamp.Address(), "offramp", "CommitReportAccepted", done)

	timeout := time.NewTimer(tests.WaitTimeout(t))
	defer timeout.Stop()

	seenMessages := NewCommitReportTracker(srcSelector, expectedSeqNumRange)

	verifyCommitReport := func(report sui_module_offramp.CommitReportAccepted) bool {
		processRoots := func(roots []sui_module_offramp.MerkleRoot) bool {
			for _, mr := range roots {
				t.Logf("(Sui) Received commit report for [%d, %d] on selector %d from source selector %d expected seq nr range %s, token prices: %v",
					mr.MinSeqNr, mr.MaxSeqNr, dest.Selector, srcSelector, expectedSeqNumRange.String(), report.PriceUpdates.TokenPriceUpdates,
				)
				seenMessages.visitCommitReport(srcSelector, mr.MinSeqNr, mr.MaxSeqNr)

				if mr.SourceChainSelector == srcSelector && uint64(expectedSeqNumRange.Start()) >= mr.MinSeqNr && uint64(expectedSeqNumRange.End()) <= mr.MaxSeqNr {
					t.Logf("(Sui) All sequence numbers committed in a single report [%d, %d]",
						expectedSeqNumRange.Start(), expectedSeqNumRange.End(),
					)
					return true
				}

				if !enforceSingleCommit && seenMessages.allCommitted(srcSelector) {
					t.Logf(
						"(Sui) All sequence numbers already committed from range [%d, %d]",
						expectedSeqNumRange.Start(), expectedSeqNumRange.End(),
					)
					return true
				}
			}
			return false
		}

		return processRoots(report.BlessedMerkleRoots) || processRoots(report.UnblessedMerkleRoots)
	}

	for {
		select {
		case event := <-sink:
			verified := verifyCommitReport(event.Event)
			if verified {
				return &event.Event, nil
			}
		case err := <-errChan:
			require.NoError(t, err)
		case <-timeout.C:
			return nil, fmt.Errorf("(sui) timed out after waiting for commit report on chain selector %d from source selector %d expected seq nr range %s",
				dest.Selector, srcSelector, expectedSeqNumRange.String())
		}
	}
}

func confirmExecWithExpectedSeqNrsSui(
	t *testing.T,
	srcSelector uint64,
	dest cldf_sui.Chain,
	offRampAddress string,
	startVersion *uint64,
	expectedSeqNrs []uint64,
) (executionStates map[uint64]int, err error) {
	if startVersion != nil {
		t.Logf("[DEBUG] startVersion = %d", *startVersion)
	} else {
		t.Log("[DEBUG] startVersion = nil (streaming from latest)")
	}

	if len(expectedSeqNrs) == 0 {
		t.Log("[DEBUG] expectedSeqNrs is empty")
		return nil, errors.New("no expected sequence numbers provided")
	}

	done := make(chan any)
	defer close(done)

	t.Log("[DEBUG] Subscribing to Sui events...", offRampAddress)
	sink, errChan := SuiEventEmitter[sui_module_offramp.ExecutionStateChanged](t, dest.Client, offRampAddress, "offramp", "ExecutionStateChanged", done)

	t.Log("[DEBUG] Event subscription established")

	executionStates = make(map[uint64]int)
	seqNrsToWatch := make(map[uint64]bool)
	for _, seqNr := range expectedSeqNrs {
		seqNrsToWatch[seqNr] = true
	}
	t.Logf("[DEBUG] Watching for sequence numbers: %+v", seqNrsToWatch)

	timeout := time.NewTimer(tests.WaitTimeout(t))
	defer timeout.Stop()

	for {
		select {
		case event := <-sink:
			t.Logf("[DEBUG] Received event: %+v", event)

			if !seqNrsToWatch[event.Event.SequenceNumber] {
				t.Logf("[DEBUG] Ignoring event with unexpected sequence number: %d", event.Event.SequenceNumber)
				continue
			}

			if event.Event.SourceChainSelector != srcSelector {
				t.Logf("[DEBUG] Ignoring event with unexpected source chain selector: got %d, expected %d",
					event.Event.SourceChainSelector, srcSelector)
				continue
			}

			if seqNrsToWatch[event.Event.SequenceNumber] && event.Event.SourceChainSelector == srcSelector {
				t.Logf("(Sui) received ExecutionStateChanged (state %s) on chain %d (offramp %s) with expected sequence number %d (tx %s)",
					executionStateToString(event.Event.State), dest.Selector, offRampAddress, event.Event.SequenceNumber, event.Version,
				)
				if event.Event.State == EXECUTION_STATE_INPROGRESS {
					continue
				}
				executionStates[event.Event.SequenceNumber] = int(event.Event.State)
				delete(seqNrsToWatch, event.Event.SequenceNumber)
				if len(seqNrsToWatch) == 0 {
					return executionStates, nil
				}
			}

		case err := <-errChan:
			require.NoError(t, err)
		case <-timeout.C:
			return nil, fmt.Errorf("(Sui) timed out waiting for ExecutionStateChanged on chain %d (offramp %s) from chain %d with expected sequence numbers %+v",
				dest.Selector, offRampAddress, srcSelector, expectedSeqNrs)
		}
	}
}
