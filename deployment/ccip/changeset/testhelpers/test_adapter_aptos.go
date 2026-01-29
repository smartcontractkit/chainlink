package testhelpers

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	aptos_ccip_offramp "github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp"
	module_offramp "github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp/offramp"
	"github.com/smartcontractkit/chainlink-aptos/relayer/codec"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	aptosstate "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/aptos"

	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
)

type AptosAdapter struct {
	state aptosstate.CCIPChainState
	cldf_aptos.Chain
}

func NewAptosAdapter(chain cldf.BlockChain, env deployment.Environment) Adapter {
	c, ok := chain.(cldf_aptos.Chain)
	if !ok {
		panic(fmt.Sprintf("invalid chain type: %T", chain))
	}
	state, err := aptosstate.LoadOnchainStateAptos(env)
	if err != nil {
		panic(fmt.Sprintf("failed to load onchain state: %T", err))
	}
	// NOTE: since this returns a copy, adapters shouldn't be constructed until everything is deployed
	s := state[c.ChainSelector()]
	return &AptosAdapter{
		state: s,
		Chain: c,
	}
}

func (a *AptosAdapter) BuildMessage(components MessageComponents) (any, error) {
	feeToken := aptos.AccountAddress{}
	if len(components.FeeToken) > 0 {
		err := feeToken.ParseStringRelaxed(components.FeeToken)
		if err != nil {
			return nil, err
		}
	}
	return AptosSendRequest{
		Data:         components.Data,
		Receiver:     common.LeftPadBytes(components.Receiver, 32),
		ExtraArgs:    components.ExtraArgs,
		FeeToken:     feeToken,
		TokenAmounts: nil,
	}, nil
}

func (a *AptosAdapter) NativeFeeToken() string {
	return "0xa"
}

func (a *AptosAdapter) GetExtraArgs(receiver []byte, sourceFamily string, opts ...ExtraArgOpt) ([]byte, error) {
	return nil, nil
}

func (a *AptosAdapter) GetInboundNonce(ctx context.Context, sender []byte, srcSel uint64) (uint64, error) {
	return 0, errors.ErrUnsupported
}

func (a *AptosAdapter) ValidateCommit(t *testing.T, sourceSelector uint64, startBlock *uint64, seqNumRange ccipocr3.SeqNumRange) {
	_, err := confirmCommitWithExpectedSeqNumRangeAptos(
		t,
		sourceSelector,
		a.Chain,
		a.state.CCIPAddress,
		startBlock,
		seqNumRange,
		true,
	)
	require.NoError(t, err)
}

func (a *AptosAdapter) ValidateExec(t *testing.T, sourceSelector uint64, startBlock *uint64, seqNrs []uint64) (executionStates map[uint64]int) {
	executionStates, err := confirmExecWithExpectedSeqNrsAptos(
		t,
		sourceSelector,
		a.Chain,
		a.state.CCIPAddress,
		startBlock,
		seqNrs,
	)
	require.NoError(t, err)
	return executionStates
}

func AptosEventEmitter[T any](
	t *testing.T,
	client aptos.AptosRpcClient,
	address aptos.AccountAddress,
	eventHandle, fieldname string,
	startVersion *uint64,
	done chan any,
) (<-chan struct {
	Event   T
	Version uint64
}, <-chan error) {
	ch := make(chan struct {
		Event   T
		Version uint64
	}, 200)
	errChan := make(chan error)
	limit := uint64(100)
	seqNum := uint64(0)
	go func() {
		ticker := time.NewTicker(time.Second * 2)
		defer ticker.Stop()

		for {
			for {
				// As this can take a few iterations if there are many events, check for done before each request
				select {
				case <-done:
					return
				default:
				}
				events, err := client.EventsByHandle(address, eventHandle, fieldname, &seqNum, &limit)
				if err != nil {
					errChan <- err
					return
				}
				if len(events) == 0 {
					// No new events found
					break
				}
				for _, event := range events {
					seqNum = event.SequenceNumber + 1
					if startVersion != nil && event.Version < *startVersion {
						continue
					}
					var out T
					if err := codec.DecodeAptosJsonValue(event.Data, &out); err != nil {
						errChan <- err
						continue
					}
					ch <- struct {
						Event   T
						Version uint64
					}{
						Event:   out,
						Version: event.Version,
					}
				}
			}
			select {
			case <-done:
				return
			case <-ticker.C:
				continue
			}
		}
	}()
	return ch, errChan
}

func confirmCommitWithExpectedSeqNumRangeAptos(
	t *testing.T,
	srcSelector uint64,
	dest cldf_aptos.Chain,
	offRampAddress aptos.AccountAddress,
	startVersion *uint64,
	expectedSeqNumRange ccipocr3.SeqNumRange,
	enforceSingleCommit bool,
) (*module_offramp.CommitReportAccepted, error) {
	boundOffRamp := aptos_ccip_offramp.Bind(offRampAddress, dest.Client)
	offRampStateAddress, err := boundOffRamp.Offramp().GetStateAddress(nil)
	require.NoError(t, err)

	done := make(chan any)
	defer close(done)
	sink, errChan := AptosEventEmitter[module_offramp.CommitReportAccepted](t, dest.Client, offRampStateAddress, offRampAddress.StringLong()+"::offramp::OffRampState", "commit_report_accepted_events", startVersion, done)

	timeout := time.NewTimer(tests.WaitTimeout(t))
	defer timeout.Stop()

	seenMessages := NewCommitReportTracker(srcSelector, expectedSeqNumRange)

	verifyCommitReport := func(report module_offramp.CommitReportAccepted) bool {
		processRoots := func(roots []module_offramp.MerkleRoot) bool {
			for _, mr := range roots {
				t.Logf("(Aptos) Received commit report for [%d, %d] on selector %d from source selector %d expected seq nr range %s, token prices: %v",
					mr.MinSeqNr, mr.MaxSeqNr, dest.Selector, srcSelector, expectedSeqNumRange.String(), report.PriceUpdates.TokenPriceUpdates,
				)
				seenMessages.visitCommitReport(srcSelector, mr.MinSeqNr, mr.MaxSeqNr)

				if mr.SourceChainSelector == srcSelector && uint64(expectedSeqNumRange.Start()) >= mr.MinSeqNr && uint64(expectedSeqNumRange.End()) <= mr.MaxSeqNr {
					t.Logf("(Aptos) All sequence numbers committed in a single report [%d, %d]",
						expectedSeqNumRange.Start(), expectedSeqNumRange.End(),
					)
					return true
				}

				if !enforceSingleCommit && seenMessages.allCommitted(srcSelector) {
					t.Logf(
						"(Aptos) All sequence numbers already committed from range [%d, %d]",
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
			return nil, fmt.Errorf("(aptos) timed out after waiting for commit report on chain selector %d from source selector %d expected seq nr range %s",
				dest.Selector, srcSelector, expectedSeqNumRange.String())
		}
	}
}

func confirmExecWithExpectedSeqNrsAptos(
	t *testing.T,
	srcSelector uint64,
	dest cldf_aptos.Chain,
	offRampAddress aptos.AccountAddress,
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

	t.Logf("[DEBUG] Binding OffRamp at address %s", offRampAddress.String())
	boundOffRamp := aptos_ccip_offramp.Bind(offRampAddress, dest.Client)
	t.Log("[DEBUG] Fetching OffRamp state address...")
	offRampStateAddress, err := boundOffRamp.Offramp().GetStateAddress(nil)
	require.NoError(t, err)
	t.Logf("[DEBUG] Got OffRamp state address: %s", offRampStateAddress.String())

	done := make(chan any)
	defer close(done)

	t.Log("[DEBUG] Subscribing to Aptos events...")
	sink, errChan := AptosEventEmitter[module_offramp.ExecutionStateChanged](t, dest.Client, offRampStateAddress, offRampAddress.StringLong()+"::offramp::OffRampState", "execution_state_changed_events", startVersion, done)

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
				t.Logf("(Aptos) received ExecutionStateChanged (state %s) on chain %d (offramp %s) with expected sequence number %d (tx %d)",
					executionStateToString(event.Event.State), dest.Selector, offRampAddress.String(), event.Event.SequenceNumber, event.Version,
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
			return nil, fmt.Errorf("(Aptos) timed out waiting for ExecutionStateChanged on chain %d (offramp %s) from chain %d with expected sequence numbers %+v",
				dest.Selector, offRampAddress.String(), srcSelector, expectedSeqNrs)
		}
	}
}
