package testhelpers

import (
	"context"
	"errors"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tvm/cell"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_ton "github.com/smartcontractkit/chainlink-deployments-framework/chain/ton"
	ops "github.com/smartcontractkit/chainlink-ton/deployment/ccip"
	tonstate "github.com/smartcontractkit/chainlink-ton/deployment/state"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

type TonAdapter struct {
	state tonstate.CCIPChainState
	*cldf_ton.Chain
}

func NewTonAdapter(chain cldf.BlockChain, state stateview.CCIPOnChainState) Adapter {
	c, ok := chain.(*cldf_ton.Chain)
	if !ok {
		panic("invalid chain type")
	}
	// NOTE: since this returns a copy, adapters shouldn't be constructed until everything is deployed
	s := state.TonChains[c.ChainSelector()]
	return &TonAdapter{
		state: s,
		Chain: c,
	}
}

func (a *TonAdapter) BuildMessage(components MessageComponents) (any, error) {
	feeToken := ops.TonTokenAddr
	if len(components.FeeToken) > 0 {
		var err error
		feeToken, err = address.ParseAddr(components.FeeToken)
		if err != nil {
			return nil, err
		}
	}

	c, err := cell.FromBOC(components.ExtraArgs)
	if err != nil {
		return nil, err
	}

	return ops.TonSendRequest{
		QueryID:   rand.Uint64(),
		Data:      components.Data,
		Receiver:  components.Receiver,
		ExtraArgs: c, // TODO handle ExtraArgs properly
		FeeToken:  feeToken,
	}, nil

}

func (a *TonAdapter) NativeFeeToken() string {
	// TODO:
	return ""
}

func (a *TonAdapter) GetExtraArgs(receiver []byte, sourceFamily string, opts ...ExtraArgOpt) ([]byte, error) {
	return nil, nil
}

func (a *TonAdapter) GetInboundNonce(ctx context.Context, sender []byte, srcSel uint64) (uint64, error) {
	return 0, errors.ErrUnsupported
}

func (a *TonAdapter) ValidateCommit(t *testing.T, sourceSelector uint64, startBlock *uint64, seqNumRange ccipocr3.SeqNumRange) {
	_, err := confirmCommitWithExpectedSeqNumRangeTON(
		t,
		sourceSelector,
		*a.Chain,
		a.state.OffRamp,
		seqNumRange,
	)
	require.NoError(t, err)
}

func (a *TonAdapter) ValidateExec(t *testing.T, sourceSelector uint64, startBlock *uint64, seqNrs []uint64) (executionStates map[uint64]int) {
	executionStates, err := confirmExecWithExpectedSeqNrsTON(
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
