package tokendata

import (
	"context"
	"errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

var (
	ErrNotReady  = errors.New("token data not ready")
	ErrRateLimit = errors.New("token data API is being rate limited")
	ErrTimeout   = errors.New("token data API timed out")
)

// Reader is an interface for fetching offchain token data
//
//go:generate mockery --quiet --name Reader --output . --filename reader_mock.go --inpackage --case=underscore
type Reader interface {
	// ReadTokenData returns the attestation bytes if ready, and throws an error if not ready.
	ReadTokenData(ctx context.Context, msg internal.EVM2EVMOnRampCCIPSendRequestedWithMeta) (tokenData []byte, err error)
	Close(qopts ...pg.QOpt) error
}
