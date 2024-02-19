package ccipdata

import "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"

const (
	ManuallyExecute = "manuallyExecute"
)

//go:generate mockery --quiet --name OffRampReader --filename offramp_reader_mock.go --case=underscore
type OffRampReader interface {
	cciptypes.OffRampReader
}
