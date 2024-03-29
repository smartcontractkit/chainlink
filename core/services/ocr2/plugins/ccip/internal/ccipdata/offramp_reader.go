package ccipdata

import (
	"context"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

const (
	ManuallyExecute = "manuallyExecute"
)

//go:generate mockery --quiet --name OffRampReader --filename offramp_reader_mock.go --case=underscore
type OffRampReader interface {
	cciptypes.OffRampReader
	//TODO Move to chainlink-common
	GetSendersNonce(ctx context.Context, senders []cciptypes.Address) (map[cciptypes.Address]uint64, error)
}
