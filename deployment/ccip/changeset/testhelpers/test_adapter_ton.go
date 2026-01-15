package testhelpers

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tvm/cell"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_ton "github.com/smartcontractkit/chainlink-deployments-framework/chain/ton"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	ops "github.com/smartcontractkit/chainlink-ton/deployment/ccip"
	tonstate "github.com/smartcontractkit/chainlink-ton/deployment/state"
	tonrouter "github.com/smartcontractkit/chainlink-ton/pkg/ccip/bindings/router"
)

type TonAdapter struct {
	state tonstate.CCIPChainState
	cldf_ton.Chain
}

func NewTonAdapter(chain cldf.BlockChain, env deployment.Environment) Adapter {
	c, ok := chain.(cldf_ton.Chain)
	if !ok {
		panic(fmt.Sprintf("invalid chain type: %T", chain))
	}

	state, err := tonstate.LoadOnchainState(env)
	if err != nil {
		panic(fmt.Sprintf("failed to load onchain state: %T", err))
	}

	// NOTE: since this returns a copy, adapters shouldn't be constructed until everything is deployed
	s := state[c.ChainSelector()]
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

	// TODO: add TokenAmounts support for TON token transfers
	return tonrouter.CCIPSend{
		QueryID:           rand.Uint64(),
		DestChainSelector: components.DestChainSelector,
		Data:              components.Data,
		Receiver:          components.Receiver,
		ExtraArgs:         c, // TODO handle ExtraArgs properly
		FeeToken:          feeToken,
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
		a.Chain,
		a.state.OffRamp,
		seqNumRange,
	)
	require.NoError(t, err)
}

func (a *TonAdapter) ValidateExec(t *testing.T, sourceSelector uint64, startBlock *uint64, seqNrs []uint64) (executionStates map[uint64]int) {
	executionStates, err := confirmExecWithExpectedSeqNrsTON(
		t,
		sourceSelector,
		a.Chain,
		a.state.OffRamp,
		startBlock,
		seqNrs,
	)
	require.NoError(t, err)
	return executionStates
}
