package ccipdata

import (
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

const (
	ManuallyExecute = "manuallyExecute"
)

type OffRampReader interface {
	cciptypes.OffRampReader
}
