package ccipdata

import (
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

const (
	ManuallyExecute = "manuallyExecute"
)

//go:generate mockery --quiet --name OffRampReader --filename offramp_reader_mock.go --case=underscore
type OffRampReader interface {
	cciptypes.OffRampReader
}
