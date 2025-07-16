package ccipton

import (
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// TONChainAccessorFactory implements cciptypes.ChainAccessorFactory for TON chains.
type TONChainAccessorFactory struct{}

// NewChainAccessor creates a new chain accessor to be used for TON chains.
func (f TONChainAccessorFactory) NewChainAccessor(_ common.ChainAccessorFactoryParams) (ccipocr3.ChainAccessor, error) {
	// TODO(NONEVM-1460): Return TONAccessor from the chainlink-ton repo. This should not be called yet since TON is
	// not yet supported.
	return nil, nil
}
