package web

import (
	"context"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// ErrSolanaNotEnabled is returned when Solana.Enabled is not true.
var ErrSolanaNotEnabled = errChainDisabled{name: "Solana", tomlKey: "Solana.Enabled"}

func NewSolanaNodesController(app chainlink.Application) NodesController {
	nodeSet := &relayerNodeSet{app.GetChains().Solana}
	return newNodesController[presenters.SolanaNodeResource](
		nodeSet, ErrSolanaNotEnabled, presenters.NewSolanaNodeResource, app.GetAuditLogger())
}

var _ chains.Nodes = (*relayerNodeSet)(nil)

type relayerNodeSet struct {
	relay.RelayerService
}

func (r *relayerNodeSet) NodeStatuses(ctx context.Context, offset, limit int, chainIDs ...string) (nodes []types.NodeStatus, count int, err error) {
	relayer, err := r.Relayer()
	if err != nil {
		return nil, -1, err
	}
	return relayer.NodeStatuses(ctx, offset, limit, chainIDs...)
}
