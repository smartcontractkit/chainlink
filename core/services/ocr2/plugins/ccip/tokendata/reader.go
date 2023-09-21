package tokendata

import (
	"context"
	"errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal"
)

var (
	ErrNotReady = errors.New("token data not ready")
)

// Reader is an interface for fetching offchain token data
//
//go:generate mockery --quiet --name Reader --output . --filename reader_mock.go --inpackage --case=underscore
type Reader interface {
	// ReadTokenData returns the attestation bytes if ready, and throws an error if not ready.
	ReadTokenData(ctx context.Context, msg internal.EVM2EVMOnRampCCIPSendRequestedWithMeta) (tokenData []byte, err error)

	// GetSourceLogPollerFilters returns the filters that should be used for the source chain log poller
	GetSourceLogPollerFilters() []logpoller.Filter
}
