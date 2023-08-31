package solana

import (
	"context"
	"fmt"
	"math/big"

	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

// TODO remove these wrappers after BCF-2441
type RelayExtender struct {
	solana.Chain
	chainImpl *chain
}

var _ relay.RelayerExt = &RelayExtender{}

func NewRelayExtender(cfg *SolanaConfig, opts ChainOpts) (*RelayExtender, error) {
	c, err := NewChain(cfg, opts)
	if err != nil {
		return nil, err
	}
	chainImpl, ok := (c).(*chain)
	if !ok {
		return nil, fmt.Errorf("internal error: cosmos relay extender got wrong type %t", c)
	}
	return &RelayExtender{Chain: chainImpl, chainImpl: chainImpl}, nil
}
func (r *RelayExtender) GetChainStatus(ctx context.Context) (relaytypes.ChainStatus, error) {
	return r.chainImpl.GetChainStatus(ctx)
}
func (r *RelayExtender) ListNodeStatuses(ctx context.Context, pageSize int32, pageToken string) (stats []relaytypes.NodeStatus, nextPageToken string, total int, err error) {
	return r.chainImpl.ListNodeStatuses(ctx, pageSize, pageToken)
}
func (r *RelayExtender) Transact(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error {
	return r.chainImpl.SendTx(ctx, from, to, amount, balanceCheck)
}
