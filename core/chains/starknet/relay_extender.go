package starknet

import (
	"context"
	"fmt"
	"math/big"

	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"
	starkchain "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/chain"

	"github.com/smartcontractkit/chainlink/v2/core/chains"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

// TODO remove these wrappers after BCF-2441
type RelayExtender struct {
	starkchain.Chain
	chainImpl *chain
}

var _ relay.RelayerExt = &RelayExtender{}

func NewRelayExtender(cfg *StarknetConfig, opts ChainOpts) (*RelayExtender, error) {
	c, err := NewChain(cfg, opts)
	if err != nil {
		return nil, err
	}
	chainImpl, ok := (c).(*chain)
	if !ok {
		return nil, fmt.Errorf("internal error: starkent relay extender got wrong type %t", c)
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
	return chains.ErrLOOPPUnsupported
}
