package testhelpers

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	sui_module_offramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_offramp/offramp"
	sui_ccip_offramp "github.com/smartcontractkit/chainlink-sui/bindings/packages/offramp"
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

func SuiEventEmitter[T any](
	t *testing.T,
	client sui.ISuiAPI,
	packageID, moduleName, event string,
	done chan any,
) (<-chan struct {
	Event   T
	Version string
}, <-chan error) {
	startTime := time.Now()
	t.Logf("[DEBUG] SuiEventEmitter: Starting at %s - will capture ALL historical events plus new ones", startTime.Format(time.RFC3339))
	ch := make(chan struct {
		Event   T
		Version string
	}, 200)
	errChan := make(chan error)
	limit := uint64(50)

	go func() {
		defer close(ch)
		defer close(errChan)

		ticker := time.NewTicker(time.Second * 2)
		defer ticker.Stop()

		var cursor interface{}

		for {
			// Drain all available pages from the current cursor position before waiting.
			for {
				select {
				case <-done:
					t.Logf("[DEBUG] SuiEventEmitter: Stopping due to done signal")
					return
				default:
				}

				eventFilter := models.EventFilterByMoveEventType{
					MoveEventType: fmt.Sprintf("%s::%s::%s", packageID, moduleName, event),
				}

				events, err := client.SuiXQueryEvents(t.Context(), models.SuiXQueryEventsRequest{
					SuiEventFilter:  eventFilter,
					Cursor:          cursor,
					Limit:           limit,
					DescendingOrder: false,
				})
				if err != nil {
					t.Logf("[DEBUG] SuiEventEmitter: Query error: %v", err)
					select {
					case errChan <- err:
					case <-done:
						return
					}
					return
				}

				if len(events.Data) == 0 {
					t.Logf("[DEBUG] SuiEventEmitter: No new events found")
					break
				}

				t.Logf("[DEBUG] SuiEventEmitter: Processing %d events", len(events.Data))

				for _, ev := range events.Data {
					eventID := fmt.Sprintf("%s:%s", ev.Id.TxDigest, ev.Id.EventSeq)

					var out T
					if err := codec.DecodeSuiJsonValue(ev.ParsedJson, &out); err != nil {
						t.Logf("[DEBUG] SuiEventEmitter: Decode error for event %s: %v (skipping)", eventID, err)
						continue
					}

					eventData := struct {
						Event   T
						Version string
					}{
						Event:   out,
						Version: ev.Id.EventSeq,
					}

					select {
					case ch <- eventData:
						t.Logf("[DEBUG] SuiEventEmitter: Sent event %s with type %s at timestamp %s", eventID, ev.Type, ev.TimestampMs)
					case <-done:
						t.Logf("[DEBUG] SuiEventEmitter: Stopping due to done signal during send")
						return
					default:
						t.Logf("[WARNING] SuiEventEmitter: Channel full, dropping event %s", eventID)
					}
				}

				// Advance the cursor so the next query picks up where we left off.
				cursor = events.NextCursor

				if !events.HasNextPage {
					break
				}
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
