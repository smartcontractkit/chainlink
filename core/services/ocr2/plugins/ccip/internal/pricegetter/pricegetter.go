package pricegetter

import cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

type PriceGetter interface {
	cciptypes.PriceGetter
}
