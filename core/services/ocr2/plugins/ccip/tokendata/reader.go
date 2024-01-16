package tokendata

import (
	"context"
	"errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal"
)

var (
	ErrNotReady        = errors.New("token data not ready")
	ErrRateLimit       = errors.New("token data API is being rate limited")
	ErrTimeout         = errors.New("token data API timed out")
	ErrRequestsBlocked = errors.New("requests are currently blocked")
)

// Reader is an interface for fetching offchain token data
//
//go:generate mockery --quiet --name Reader --output . --filename reader_mock.go --inpackage --case=underscore
type Reader interface {
	// ReadTokenData returns the attestation bytes if ready, and throws an error if not ready.
	// It supports messages with a single token transfer, the returned []byte has the token data for the first token of the msg.
	ReadTokenData(ctx context.Context, msg internal.EVM2EVMOnRampCCIPSendRequestedWithMeta, tokenIndex int) (tokenData []byte, err error)
}
