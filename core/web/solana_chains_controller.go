package web

import (
	"context"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func NewSolanaChainsController(app chainlink.Application) ChainsController {
	chainSet := &relayerChainSet{app.GetChains().SolanaRelayer}
	return newChainsController("solana", chainSet, ErrSolanaNotEnabled,
		presenters.NewSolanaChainResource, app.GetLogger(), app.GetAuditLogger())
}

type relayerChainSet struct {
	relayerFn func() (loop.Relayer, error)
}

func (r *relayerChainSet) ChainStatus(ctx context.Context, id string) (types.ChainStatus, error) {
	relayer, err := r.relayerFn()
	if err != nil {
		return types.ChainStatus{}, err
	}
	return relayer.ChainStatus(ctx, id)
}

func (r *relayerChainSet) ChainStatuses(ctx context.Context, offset, limit int) ([]types.ChainStatus, int, error) {
	relayer, err := r.relayerFn()
	if err != nil {
		return []types.ChainStatus{}, -1, err
	}
	return relayer.ChainStatuses(ctx, offset, limit)
}
