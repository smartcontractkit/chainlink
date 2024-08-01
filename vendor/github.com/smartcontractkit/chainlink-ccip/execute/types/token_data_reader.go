package types

import (
	"context"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

// TokenDataReader is an interface for reading extra token data from an async process.
// TODO: Build a token data reading process.
type TokenDataReader interface {
	ReadTokenData(ctx context.Context, srcChain cciptypes.ChainSelector, num cciptypes.SeqNum) ([][]byte, error)
}
