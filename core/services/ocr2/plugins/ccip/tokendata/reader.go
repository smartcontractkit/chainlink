package tokendata

import (
	"errors"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

var (
	ErrNotReady        = errors.New("token data not ready")
	ErrRateLimit       = errors.New("token data API is being rate limited")
	ErrTimeout         = errors.New("token data API timed out")
	ErrRequestsBlocked = errors.New("requests are currently blocked")
)

// Reader is an interface for fetching offchain token data
type Reader interface {
	cciptypes.TokenDataReader
}
