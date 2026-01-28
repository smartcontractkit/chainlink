package testhelpers

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"

	solconfig "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/latest/ccip_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_0/ccip_router"
	solccip "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/ccip"
	solcommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip/consts"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	svm_stateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
)

type SVMAdapter struct {
	state svm_stateview.CCIPChainState
	cldf_solana.Chain
}

func NewSVMAdapter(chain cldf.BlockChain, env deployment.Environment) Adapter {
	c, ok := chain.(cldf_solana.Chain)
	if !ok {
		panic(fmt.Sprintf("invalid chain type: %T", chain))
	}
	state, err := stateview.LoadOnchainStateSolana(env)
	if err != nil {
		panic(fmt.Sprintf("failed to load onchain state: %T", err))
	}
	// NOTE: since this returns a copy, adapters shouldn't be constructed until everything is deployed
	s := state.SolChains[c.ChainSelector()]
	return &SVMAdapter{
		state: s,
		Chain: c,
	}
}

func (a *SVMAdapter) BuildMessage(components MessageComponents) (any, error) {
	feeToken := solana.PublicKey{}
	if len(components.FeeToken) > 0 {
		var err error
		feeToken, err = solana.PublicKeyFromBase58(components.FeeToken)
		if err != nil {
			return nil, err
		}
	}

	return ccip_router.SVM2AnyMessage{
		Receiver:     common.LeftPadBytes(components.Receiver, 32),
		TokenAmounts: nil,
		Data:         components.Data,
		FeeToken:     feeToken,
		ExtraArgs:    components.ExtraArgs,
	}, nil
}

func (a *SVMAdapter) NativeFeeToken() string {
	return solana.SolMint.String()
}

func (a *SVMAdapter) GetExtraArgs(receiver []byte, sourceFamily string, opts ...ExtraArgOpt) ([]byte, error) {
	return nil, nil
}

func (a *SVMAdapter) GetInboundNonce(ctx context.Context, sender []byte, srcSel uint64) (uint64, error) {
	// TODO: solcommon.FindNoncePDA expected the sender to be a solana pubkey
	chainSelectorLE := solcommon.Uint64ToLE(a.Selector)
	noncePDA, _, err := solana.FindProgramAddress([][]byte{[]byte("nonce"), chainSelectorLE, sender}, a.state.Router)
	if err != nil {
		return 0, err
	}
	var nonceCounterAccount ccip_router.Nonce
	// we ignore the error because the account might not exist yet
	_ = solcommon.GetAccountDataBorshInto(ctx, a.Client, noncePDA, solconfig.DefaultCommitment, &nonceCounterAccount)
	return nonceCounterAccount.Counter, nil
}

func (a *SVMAdapter) ValidateCommit(t *testing.T, sourceSelector uint64, startBlock *uint64, seqNumRange ccipocr3.SeqNumRange) {
	var startSlot uint64
	if startBlock != nil {
		startSlot = *startBlock
	}
	_, err := confirmCommitWithExpectedSeqNumRangeSol(
		t,
		sourceSelector,
		a.Chain,
		a.state.OffRamp,
		startSlot,
		seqNumRange,
		true,
	)
	require.NoError(t, err)
}

func (a *SVMAdapter) ValidateExec(t *testing.T, sourceSelector uint64, startBlock *uint64, seqNrs []uint64) (executionStates map[uint64]int) {
	var startSlot uint64
	if startBlock != nil {
		startSlot = *startBlock
	}
	executionStates, err := confirmExecWithSeqNrsSol(
		t,
		sourceSelector,
		a.Chain,
		a.state.OffRamp,
		startSlot,
		seqNrs,
	)
	require.NoError(t, err)
	return executionStates
}

type EventWithTxn[T any] struct {
	Event T
	Txn   *solrpc.GetTransactionResult
}

// SolEventEmitter listens for events of type T emitted by the Solana program at the given address. Failed transactions
// can be included by setting the includeFailed flag to true.
func SolEventEmitter[T any](ctx context.Context, client *solrpc.Client, address solana.PublicKey, eventType string, startSlot uint64, done chan any, ticker *time.Ticker, includeFailed bool) (<-chan EventWithTxn[T], <-chan error) {
	ch := make(chan EventWithTxn[T])
	errorCh := make(chan error)
	go func() {
		defer ticker.Stop()
		var until solana.Signature
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				// Scan for transactions referencing the address
				txSigs, err := client.GetSignaturesForAddressWithOpts(
					ctx,
					address,
					&solrpc.GetSignaturesForAddressOpts{
						Commitment: solrpc.CommitmentConfirmed,
						Until:      until,
					},
				)
				if err != nil {
					errorCh <- err
					return
				}

				if len(txSigs) == 0 {
					continue
				}

				// values are returned ordered newest to oldest, so we replay them backwards
				for _, txSig := range slices.Backward(txSigs) {
					if txSig.Err != nil && !includeFailed {
						// We're not interested in failed transactions.
						continue
					}
					if txSig.Slot < startSlot {
						// Skip any signatures that are before the starting slot
						continue
					}
					v := uint64(0) // v0 = latest, supports address table lookups
					tx, err := client.GetTransaction(
						ctx,
						txSig.Signature,
						&solrpc.GetTransactionOpts{
							Commitment:                     solrpc.CommitmentConfirmed,
							Encoding:                       solana.EncodingBase64,
							MaxSupportedTransactionVersion: &v,
						},
					)
					if err != nil {
						errorCh <- err
						return
					}

					events, err := solcommon.ParseMultipleEvents[T](tx.Meta.LogMessages, eventType, solconfig.PrintEvents)
					if err != nil && strings.Contains(err.Error(), "event not found") {
						continue
					}
					if err != nil {
						errorCh <- err
						return
					}

					for _, event := range events {
						select {
						case ch <- EventWithTxn[T]{
							Event: event,
							Txn:   tx,
						}:
						case <-done:
							return
						}
					}
				}
				// next scan should stop at the newest signature we've received
				until = txSigs[0].Signature
			}
		}
	}()

	return ch, errorCh
}

func confirmCommitWithExpectedSeqNumRangeSol(
	t *testing.T,
	srcSelector uint64,
	dest cldf_solana.Chain,
	offrampAddress solana.PublicKey,
	startSlot uint64,
	expectedSeqNumRange ccipocr3.SeqNumRange,
	enforceSingleCommit bool,
) (bool, error) {
	seenMessages := NewCommitReportTracker(srcSelector, expectedSeqNumRange)

	done := make(chan any)
	defer close(done)
	sink, errCh := SolEventEmitter[solcommon.EventCommitReportAccepted](t.Context(), dest.Client, offrampAddress, consts.EventNameCommitReportAccepted, startSlot, done, time.NewTicker(2*time.Second), false)

	timeout := time.NewTimer(tests.WaitTimeout(t))
	defer timeout.Stop()

	for {
		select {
		case eventWithTxn := <-sink:
			commitEvent := eventWithTxn.Event
			// if merkle root is zero, it only contains price updates
			if commitEvent.Report == nil {
				t.Logf("Skipping CommitReportAccepted with only price updates")
				continue
			}
			require.Equal(t, srcSelector, commitEvent.Report.SourceChainSelector)

			// TODO: this logic is duplicated with verifyCommitReport, share
			mr := commitEvent.Report
			seenMessages.visitCommitReport(mr.SourceChainSelector, mr.MinSeqNr, mr.MaxSeqNr)
			if mr.SourceChainSelector == srcSelector &&
				uint64(expectedSeqNumRange.Start()) >= mr.MinSeqNr &&
				uint64(expectedSeqNumRange.End()) <= mr.MaxSeqNr {
				t.Logf("All sequence numbers committed in a single report [%d, %d]", expectedSeqNumRange.Start(), expectedSeqNumRange.End())
				return true, nil
			}

			if !enforceSingleCommit && seenMessages.allCommitted(srcSelector) {
				t.Logf("All sequence numbers already committed from range [%d, %d]", expectedSeqNumRange.Start(), expectedSeqNumRange.End())
				return true, nil
			}
		case err := <-errCh:
			require.NoError(t, err)
		case <-timeout.C:
			return false, fmt.Errorf("timed out after waiting for commit report on chain selector %d from source selector %d expected seq nr range %s",
				dest.Selector, srcSelector, expectedSeqNumRange.String())
		}
	}
}

type MessageStateEvent struct {
	SequenceNumber uint64
	Block          uint64
	State          ccip_offramp.MessageExecutionState
}

// GetMessageStatesWithSeqNrsSol waits for execution state changes on the destination Solana chain with the expected
// sequence numbers. The timeout is configurable.
// If "inProgress" is true, and there are unexecutable messages, it will continue to watch for additional "InProgress"
// states for the entire timeoutDuration.
func GetMessageStatesWithSeqNrsSol(
	t *testing.T,
	timeoutDuration time.Duration,
	srcSelector uint64,
	dest cldf_solana.Chain,
	offrampAddress solana.PublicKey,
	startSlot uint64,
	expectedSeqNrs []uint64,
	inProgress bool,
) (executionStates map[uint64][]MessageStateEvent, err error) {
	// TODO: share with EVM
	// some state to efficiently track the execution states
	// of all the expected sequence numbers.
	executionStates = make(map[uint64][]MessageStateEvent)
	seqNrsInProgress := make(map[uint64]struct{})
	seqNrsToWatch := make(map[uint64]struct{})
	for _, seqNr := range expectedSeqNrs {
		seqNrsToWatch[seqNr] = struct{}{}
		seqNrsInProgress[seqNr] = struct{}{}
	}

	done := make(chan any)
	defer close(done)
	sink, errCh := SolEventEmitter[solccip.EventExecutionStateChanged](t.Context(), dest.Client, offrampAddress, consts.EventNameExecutionStateChanged, startSlot, done, time.NewTicker(2*time.Second), inProgress)

	timeout := time.NewTimer(timeoutDuration)
	defer timeout.Stop()

	for {
		select {
		case eventWithTxn := <-sink:
			execEvent := eventWithTxn.Event
			// TODO: share with EVM
			_, found := seqNrsToWatch[execEvent.SequenceNumber]
			if found && execEvent.SourceChainSelector == srcSelector {
				t.Logf("Received ExecutionStateChanged (state %s) on chain %d (offramp %s) from chain %d with expected sequence number %d",
					execEvent.State.String(), dest.Selector, offrampAddress.String(), srcSelector, execEvent.SequenceNumber)
				executionStates[execEvent.SequenceNumber] = append(executionStates[execEvent.SequenceNumber],
					MessageStateEvent{
						SequenceNumber: execEvent.SequenceNumber,
						Block:          eventWithTxn.Txn.Slot,
						State:          execEvent.State,
					})
				if execEvent.State == ccip_offramp.InProgress_MessageExecutionState {
					delete(seqNrsInProgress, execEvent.SequenceNumber)
					// continue watching for final state or timeout
					continue
				}
				delete(seqNrsToWatch, execEvent.SequenceNumber)
				delete(seqNrsInProgress, execEvent.SequenceNumber)
				if len(seqNrsToWatch) == 0 {
					return executionStates, nil
				}
			}
		case err := <-errCh:
			require.NoError(t, err)
		case <-timeout.C:
			// If we have some status for every seqNr, return what we have instead of an error.
			if len(seqNrsInProgress) == 0 {
				return executionStates, nil
			}
			// Otherwise, return a timeout error.
			return nil, fmt.Errorf("timed out waiting for ExecutionStateChanged on chain %d (offramp %s) from chain %d with expected sequence numbers %+v",
				dest.Selector, offrampAddress.String(), srcSelector, expectedSeqNrs)
		}
	}
}

// ConfirmExecWithSeqNrsSol waits for an execution state change on the destination Solana chain with the expected
// sequence numbers. The timeout is automatically set to 30 seconds or 90% of the test timeout (if there is one).
// This is a wrapper around the more general GetMessageStatesWithSeqNrsSol.
func confirmExecWithSeqNrsSol(
	t *testing.T,
	srcSelector uint64,
	dest cldf_solana.Chain,
	offrampAddress solana.PublicKey,
	startSlot uint64,
	expectedSeqNrs []uint64,
) (executionStates map[uint64]int, err error) {
	timeout := tests.WaitTimeout(t)
	states, err := GetMessageStatesWithSeqNrsSol(t, timeout, srcSelector, dest, offrampAddress, startSlot, expectedSeqNrs, false)
	if err != nil {
		return nil, err
	}

	executionStates = make(map[uint64]int)

	for seqNr, stateList := range states {
		if len(stateList) == 0 {
			return nil, fmt.Errorf("no execution states found for seqNr %d", seqNr)
		}
		// check that the final state is either success or failure
		state := stateList[len(stateList)-1].State
		if state != ccip_offramp.Success_MessageExecutionState && state != ccip_offramp.Failure_MessageExecutionState {
			return nil, fmt.Errorf("expected final execution state for seqNr %d, got %s", seqNr, state.String())
		}

		executionStates[seqNr] = int(state)
	}

	return executionStates, nil
}
