package pricegetter

import cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

//go:generate mockery --quiet --name PriceGetter --output . --filename mock.go --inpackage --case=underscore
type PriceGetter interface {
	cciptypes.PriceGetter
}
