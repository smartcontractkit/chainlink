package testhelpers

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	evm_stateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
	"github.com/stretchr/testify/require"
)

type EVMAdapter struct {
	state evm_stateview.CCIPChainState
	*evm.Chain
}

func NewEVMAdapter(chain cldf.BlockChain, state stateview.CCIPOnChainState) Adapter {
	c, ok := chain.(*evm.Chain)
	if !ok {
		panic("invalid chain type")
	}
	// NOTE: since this returns a copy, adapters shouldn't be constructed until everything is deployed
	s := state.Chains[c.ChainSelector()]
	return &EVMAdapter{
		state: s,
		Chain: c,
	}
}

func (a *EVMAdapter) BuildMessage(components MessageComponents) (any, error) {
	receiver := common.LeftPadBytes(components.Receiver, 32)
	feeToken := common.HexToAddress(a.NativeFeeToken())
	if len(components.FeeToken) > 0 {
		feeToken = common.HexToAddress(components.FeeToken)
	}

	return router.ClientEVM2AnyMessage{
		Receiver:     receiver,
		Data:         components.Data,
		TokenAmounts: nil, // TODO:
		FeeToken:     feeToken,
		ExtraArgs:    components.ExtraArgs,
	}, nil
}

func (a *EVMAdapter) NativeFeeToken() string {
	return "0x0"
}

func (a *EVMAdapter) GetExtraArgs(receiver []byte, sourceFamily string, opts ...ExtraArgOpt) ([]byte, error) {
	return nil, nil
}

func (a *EVMAdapter) GetInboundNonce(ctx context.Context, sender []byte, srcSel uint64) (uint64, error) {
	return a.state.NonceManager.GetInboundNonce(&bind.CallOpts{Context: ctx}, srcSel, sender)
}

func (a *EVMAdapter) ValidateCommit(t *testing.T, sourceSelector uint64, startBlock *uint64, seqNumRange ccipocr3.SeqNumRange) {
	_, err := ConfirmCommitWithExpectedSeqNumRange(
		t,
		sourceSelector,
		*a.Chain,
		a.state.OffRamp,
		startBlock,
		seqNumRange,
		true,
	)
	require.NoError(t, err)
}

func (a *EVMAdapter) ValidateExec(t *testing.T, sourceSelector uint64, startBlock *uint64, seqNrs []uint64) (executionStates map[uint64]int) {
	executionStates, err := ConfirmExecWithSeqNrs(
		t,
		sourceSelector,
		*a.Chain,
		a.state.OffRamp,
		startBlock,
		seqNrs,
	)
	require.NoError(t, err)
	return executionStates
}

// ConfirmCommitWithExpectedSeqNumRange waits for a commit report on the destination chain with the expected sequence number range.
// startBlock is the block number to start watching from.
// If startBlock is nil, it will start watching from the latest block.
func ConfirmCommitWithExpectedSeqNumRange(
	t *testing.T,
	srcSelector uint64,
	dest cldf_evm.Chain,
	offRamp offramp.OffRampInterface,
	startBlock *uint64,
	expectedSeqNumRange ccipocr3.SeqNumRange,
	enforceSingleCommit bool,
) (*offramp.OffRampCommitReportAccepted, error) {
	sink := make(chan *offramp.OffRampCommitReportAccepted)
	subscription, err := offRamp.WatchCommitReportAccepted(&bind.WatchOpts{
		Context: context.Background(),
		Start:   startBlock,
	}, sink)
	if err != nil {
		return nil, fmt.Errorf("error to subscribe CommitReportAccepted : %w", err)
	}

	seenMessages := NewCommitReportTracker(srcSelector, expectedSeqNumRange)

	verifyCommitReport := func(report *offramp.OffRampCommitReportAccepted) bool {
		processRoots := func(roots []offramp.InternalMerkleRoot) bool {
			for _, mr := range roots {
				t.Logf(
					"Received commit report for [%d, %d] on selector %d from source selector %d expected seq nr range %s, token prices: %v",
					mr.MinSeqNr, mr.MaxSeqNr, dest.Selector, srcSelector, expectedSeqNumRange.String(), report.PriceUpdates.TokenPriceUpdates,
				)
				seenMessages.visitCommitReport(srcSelector, mr.MinSeqNr, mr.MaxSeqNr)

				if mr.SourceChainSelector == srcSelector &&
					uint64(expectedSeqNumRange.Start()) >= mr.MinSeqNr &&
					uint64(expectedSeqNumRange.End()) <= mr.MaxSeqNr {
					t.Logf(
						"All sequence numbers committed in a single report [%d, %d]",
						expectedSeqNumRange.Start(), expectedSeqNumRange.End(),
					)
					return true
				}

				if !enforceSingleCommit && seenMessages.allCommited(srcSelector) {
					t.Logf(
						"All sequence numbers already committed from range [%d, %d]",
						expectedSeqNumRange.Start(), expectedSeqNumRange.End(),
					)
					return true
				}
			}
			return false
		}

		return processRoots(report.BlessedMerkleRoots) || processRoots(report.UnblessedMerkleRoots)
	}

	defer subscription.Unsubscribe()
	timeout := time.NewTimer(tests.WaitTimeout(t))
	defer timeout.Stop()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-t.Context().Done():
			return nil, nil
		case <-ticker.C:
			t.Logf("Waiting for commit report on chain selector %d from source selector %d expected seq nr range %s",
				dest.Selector, srcSelector, expectedSeqNumRange.String())

			// Need to do this because the subscription sometimes fails to get the event.
			iter, err := offRamp.FilterCommitReportAccepted(&bind.FilterOpts{
				Context: t.Context(),
			})

			// In some test case the test ends while the filter is still running resulting in a context.Canceled error.
			if err != nil && !errors.Is(err, context.Canceled) {
				require.NoError(t, err)
			}
			for iter.Next() {
				event := iter.Event
				verified := verifyCommitReport(event)
				if verified {
					return event, nil
				}
			}
		case subErr := <-subscription.Err():
			return nil, fmt.Errorf("subscription error: %w", subErr)
		case <-timeout.C:
			return nil, fmt.Errorf("timed out after waiting for commit report on chain selector %d from source selector %d expected seq nr range %s",
				dest.Selector, srcSelector, expectedSeqNumRange.String())
		case report := <-sink:
			verified := verifyCommitReport(report)
			if verified {
				return report, nil
			}
		}
	}
}

// ConfirmExecWithSeqNrs waits for an execution state change on the destination chain with the expected sequence number.
// startBlock is the block number to start watching from.
// If startBlock is nil, it will start watching from the latest block.
// Returns a map that maps the expected sequence number to its execution state.
func ConfirmExecWithSeqNrs(
	t *testing.T,
	sourceSelector uint64,
	dest cldf_evm.Chain,
	offRamp offramp.OffRampInterface,
	startBlock *uint64,
	expectedSeqNrs []uint64,
) (executionStates map[uint64]int, err error) {
	if len(expectedSeqNrs) == 0 {
		return nil, errors.New("no expected sequence numbers provided")
	}

	timeout := time.NewTimer(tests.WaitTimeout(t))
	defer timeout.Stop()
	tick := time.NewTicker(3 * time.Second)
	defer tick.Stop()
	sink := make(chan *offramp.OffRampExecutionStateChanged)
	subscription, err := offRamp.WatchExecutionStateChanged(&bind.WatchOpts{
		Context: context.Background(),
		Start:   startBlock,
	}, sink, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error to subscribe ExecutionStateChanged : %w", err)
	}
	defer subscription.Unsubscribe()

	// some state to efficiently track the execution states
	// of all the expected sequence numbers.
	executionStates = make(map[uint64]int)
	seqNrsToWatch := make(map[uint64]struct{})
	for _, seqNr := range expectedSeqNrs {
		seqNrsToWatch[seqNr] = struct{}{}
	}
	for {
		select {
		case <-tick.C:
			for expectedSeqNr := range seqNrsToWatch {
				scc, executionState := getExecutionState(t, sourceSelector, offRamp, expectedSeqNr)
				t.Logf("Waiting for ExecutionStateChanged on chain %d (offramp %s) from chain %d with expected sequence number %d, current onchain minSeqNr: %d, execution state: %s",
					dest.Selector, offRamp.Address().String(), sourceSelector, expectedSeqNr, scc.MinSeqNr, executionStateToString(executionState))
				if executionState == EXECUTION_STATE_SUCCESS || executionState == EXECUTION_STATE_FAILURE {
					t.Logf("Observed %s execution state on chain %d (offramp %s) from chain %d with expected sequence number %d",
						executionStateToString(executionState), dest.Selector, offRamp.Address().String(), sourceSelector, expectedSeqNr)
					executionStates[expectedSeqNr] = int(executionState)
					delete(seqNrsToWatch, expectedSeqNr)
					if len(seqNrsToWatch) == 0 {
						return executionStates, nil
					}
				}
			}
		case execEvent := <-sink:
			t.Logf("Received ExecutionStateChanged (state %s) for seqNum %d on chain %d (offramp %s) from chain %d",
				executionStateToString(execEvent.State), execEvent.SequenceNumber, dest.Selector, offRamp.Address().String(),
				sourceSelector,
			)

			_, found := seqNrsToWatch[execEvent.SequenceNumber]
			if found && execEvent.SourceChainSelector == sourceSelector {
				t.Logf("Received ExecutionStateChanged (state %s) on chain %d (offramp %s) from chain %d with expected sequence number %d",
					executionStateToString(execEvent.State), dest.Selector, offRamp.Address().String(), sourceSelector, execEvent.SequenceNumber)
				executionStates[execEvent.SequenceNumber] = int(execEvent.State)
				delete(seqNrsToWatch, execEvent.SequenceNumber)
				if len(seqNrsToWatch) == 0 {
					return executionStates, nil
				}
			}
		case <-timeout.C:
			return nil, fmt.Errorf("timed out waiting for ExecutionStateChanged on chain %d (offramp %s) from chain %d with expected sequence numbers %+v",
				dest.Selector, offRamp.Address().String(), sourceSelector, expectedSeqNrs)
		case subErr := <-subscription.Err():
			return nil, fmt.Errorf("subscription error: %w", subErr)
		}
	}
}
