package adapters

import (
	"context"
	crand "crypto/rand"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/message_hasher"
	solconfig "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solccip "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/ccip"
	solcommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	solanastate "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
	"github.com/smartcontractkit/chainlink/integration-tests/smoke/ccip/canonical/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
)

var _ types.Adapter = &svmAdapter{}

type svmAdapter struct {
	chain      cldf.SolChain
	chainState solanastate.CCIPChainState
}

// ChainSelector implements types.Adapter.
func (s *svmAdapter) ChainSelector() uint64 {
	return s.chain.Selector
}

// GetExtraArgs implements types.Adapter.
// TODO: apply opts.
func (s *svmAdapter) GetExtraArgs(receiver []byte, sourceFamily string, opts ...types.ExtraArgOpt) ([]byte, error) {
	receiverProgram := solana.PublicKeyFromBytes(receiver)
	receiverTargetAccountPDA, _, _ := solana.FindProgramAddress([][]byte{[]byte("counter")}, receiverProgram)
	receiverExternalExecutionConfigPDA, _, _ := solana.FindProgramAddress([][]byte{[]byte("external_execution_config")}, receiverProgram)

	accounts := [][32]byte{
		receiverExternalExecutionConfigPDA,
		receiverTargetAccountPDA,
		solana.SystemProgramID,
	}

	switch sourceFamily {
	case chain_selectors.FamilyEVM:
		extraArgs, err := ccipevm.SerializeClientSVMExtraArgsV1(message_hasher.ClientSVMExtraArgsV1{
			AccountIsWritableBitmap: solccip.GenerateBitMapForIndexes([]int{0, 1}),
			Accounts:                accounts,
			ComputeUnits:            80_000,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to serialize extra args: %w", err)
		}
		return extraArgs, nil
	default:
		// TODO: add support for other families
		return nil, fmt.Errorf("unsupported source family: %s", sourceFamily)
	}
}

// CCIPReceiver implements types.Adapter.
func (s *svmAdapter) CCIPReceiver() []byte {
	return s.chainState.Receiver.Bytes()
}

func NewSVMAdapter(chain cldf.SolChain, chainState solanastate.CCIPChainState) types.Adapter {
	return &svmAdapter{
		chain:      chain,
		chainState: chainState,
	}
}

// BuildMessage implements types.Adapter.
func (s *svmAdapter) BuildMessage(components types.MessageComponents) (any, error) {
	feeToken := solana.PublicKey{}
	if len(components.FeeToken) > 0 {
		var err error
		feeToken, err = solana.PublicKeyFromBase58(components.FeeToken)
		if err != nil {
			return nil, fmt.Errorf("invalid format for fee token: %w", err)
		}
	}

	var tokenAmounts []ccip_router.SVMTokenAmount
	if len(components.TokenAmounts) > 0 {
		tokenAmounts = make([]ccip_router.SVMTokenAmount, len(components.TokenAmounts))
		for i, amount := range components.TokenAmounts {
			token, err := solana.PublicKeyFromBase58(amount.Token)
			if err != nil {
				return nil, fmt.Errorf("invalid format for token: %w", err)
			}
			tokenAmounts[i] = ccip_router.SVMTokenAmount{
				Token:  token,
				Amount: amount.Amount.Uint64(),
			}
		}
	}

	msg := ccip_router.SVM2AnyMessage{
		Receiver:     components.Receiver,
		TokenAmounts: tokenAmounts,
		Data:         components.Data,
		FeeToken:     feeToken,
		ExtraArgs:    components.ExtraArgs,
	}

	return msg, nil
}

// ChainFamily implements types.Adapter.
func (s *svmAdapter) ChainFamily() string {
	return chain_selectors.FamilySolana
}

// GetInboundNonce implements types.Adapter.
func (s *svmAdapter) GetInboundNonce(ctx context.Context, sender []byte, srcSel uint64) (uint64, error) {
	client := s.chain.Client
	// TODO: solcommon.FindNoncePDA expected the sender to be a solana pubkey
	chainSelectorLE := solcommon.Uint64ToLE(s.chain.Selector)
	noncePDA, _, err := solana.FindProgramAddress([][]byte{[]byte("nonce"), chainSelectorLE, sender}, s.chainState.Router)
	if err != nil {
		return 0, fmt.Errorf("failed to find nonce PDA: %w", err)
	}
	var nonceCounterAccount ccip_router.Nonce
	// we ignore the error because the account might not exist yet
	_ = solcommon.GetAccountDataBorshInto(ctx, client, noncePDA, solconfig.DefaultCommitment, &nonceCounterAccount)
	latestNonce := nonceCounterAccount.Counter
	return latestNonce, nil
}

// NativeFeeToken implements types.Adapter.
func (s *svmAdapter) NativeFeeToken() string {
	return solana.PublicKey{}.String()
}

// RandomReceiver implements types.Adapter.
func (s *svmAdapter) RandomReceiver() []byte {
	b := make([]byte, 20)
	_, _ = crand.Read(b) // Assignment for errcheck. Only used in tests so we can ignore.
	// return a random address as a left-padded 32 byte array
	addr := common.LeftPadBytes(b, 32)

	return addr
}

// ValidateCommit implements types.Adapter.
func (s *svmAdapter) ValidateCommit(t *testing.T, sourceSelector uint64, seqNumRange ccipocr3.SeqNumRange) {
	testhelpers.ConfirmCommitWithExpectedSeqNumRangeSol(
		t,
		sourceSelector,
		s.chain,
		s.chainState.OffRamp,
		0, // startSlot
		seqNumRange,
		true,
	)
}

// ValidateExec implements types.Adapter.
func (s *svmAdapter) ValidateExec(t *testing.T, sourceSelector uint64, seqNrs []uint64) {
	testhelpers.ConfirmExecWithSeqNrsSol(
		t,
		sourceSelector,
		s.chain,
		s.chainState.OffRamp,
		0, // startSlot
		seqNrs,
	)
}
