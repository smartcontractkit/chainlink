package adapters

import (
	"context"
	crand "crypto/rand"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
	"github.com/smartcontractkit/chainlink/integration-tests/smoke/ccip/canonical/types"
)

var _ types.Adapter = &evmAdapter{}

type evmAdapter struct {
	chain      cldf.Chain
	chainState evm.CCIPChainState
}

// GetExtraArgs implements types.Adapter.
// TODO: apply opts.
func (a *evmAdapter) GetExtraArgs(_ []byte, sourceFamily string, opts ...types.ExtraArgOpt) ([]byte, error) {
	switch sourceFamily {
	case chain_selectors.FamilyEVM:
		return []byte{}, nil // default extra args are empty for EVM
	case chain_selectors.FamilySolana:
		return []byte{}, nil // default extra args are empty for Solana
	default:
		return nil, fmt.Errorf("unsupported source family: %s", sourceFamily)
	}
}

// CCIPReceiver implements types.Adapter.
func (a *evmAdapter) CCIPReceiver() []byte {
	return common.LeftPadBytes(a.chainState.Receiver.Address().Bytes(), 32)
}

// ValidateCommit implements types.Adapter.
func (a *evmAdapter) ValidateCommit(
	t *testing.T,
	sourceSelector uint64,
	seqNumRange ccipocr3.SeqNumRange) {
	testhelpers.ConfirmCommitWithExpectedSeqNumRange(
		t,
		sourceSelector,
		a.chain,
		a.chainState.OffRamp,
		nil,
		seqNumRange,
		false,
	)
}

// ValidateExec implements types.Adapter.
func (a *evmAdapter) ValidateExec(
	t *testing.T,
	sourceSelector uint64,
	seqNrs []uint64) {
	testhelpers.ConfirmExecWithSeqNrs(
		t,
		sourceSelector,
		a.chain,
		a.chainState.OffRamp,
		nil,
		seqNrs,
	)
}

// BuildMessage implements types.Adapter.
func (a *evmAdapter) BuildMessage(components types.MessageComponents) (any, error) {
	var tokenAmounts []router.ClientEVMTokenAmount
	for _, tokenAmount := range components.TokenAmounts {
		tokenAmounts = append(tokenAmounts, router.ClientEVMTokenAmount{
			Token:  common.HexToAddress(tokenAmount.Token),
			Amount: tokenAmount.Amount,
		})
	}

	msg := router.ClientEVM2AnyMessage{
		Receiver:     components.Receiver,
		Data:         components.Data,
		TokenAmounts: tokenAmounts,
		FeeToken:     common.HexToAddress(components.FeeToken),
		ExtraArgs:    components.ExtraArgs,
	}

	return msg, nil
}

// NativeFeeToken implements types.Adapter.
func (a *evmAdapter) NativeFeeToken() string {
	return common.HexToAddress("0x0").Hex()
}

// RandomReceiver implements types.Adapter.
func (a *evmAdapter) RandomReceiver() []byte {
	b := make([]byte, 20)
	_, _ = crand.Read(b) // Assignment for errcheck. Only used in tests so we can ignore.
	// return a random address as a left-padded 32 byte array
	addr := common.LeftPadBytes(b, 32)

	return addr
}

func NewEVMAdapter(chain cldf.Chain, chainState evm.CCIPChainState) types.Adapter {
	return &evmAdapter{chain: chain, chainState: chainState}
}

func (a *evmAdapter) ChainFamily() string {
	return chain_selectors.FamilyEVM
}

func (a *evmAdapter) GetInboundNonce(ctx context.Context, sender []byte, srcSel uint64) (uint64, error) {
	return a.chainState.NonceManager.GetInboundNonce(&bind.CallOpts{
		Context: ctx,
	}, srcSel, sender)
}
