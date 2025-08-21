package ccipton

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// TONChainAccessorFactory implements cciptypes.ChainAccessorFactory for TON chains.
type TONChainAccessorFactory struct{}

// NewChainAccessor creates a new chain accessor to be used for TON chains.
func (f TONChainAccessorFactory) NewChainAccessor(params common.ChainAccessorFactoryParams) (ccipocr3.ChainAccessor, error) {
	p, e := params.Relayer.NewCCIPProvider(context.Background(), types.RelayArgs{})
	if e != nil {
		return nil, e
	}

	return p.ChainAccessor(), nil
}
