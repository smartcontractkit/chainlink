package pricegetter

import "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"

//go:generate mockery --quiet --name PriceGetter --output . --filename mock.go --inpackage --case=underscore
type PriceGetter interface {
	cciptypes.PriceGetter
}
