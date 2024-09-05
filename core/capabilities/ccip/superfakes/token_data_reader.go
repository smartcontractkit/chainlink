package superfakes

import (
	"context"

	"github.com/smartcontractkit/chainlink-ccip/execute/exectypes"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

// NewNilTokenDataReader returns a new nilTokenDataReader.
// This token data reader always returns nil for the token data.
func NewNilTokenDataReader() exectypes.TokenDataReader {
	return &nilTokenDataReader{}
}

type nilTokenDataReader struct{}

// ReadTokenData implements exectypes.TokenDataReader.
func (t *nilTokenDataReader) ReadTokenData(ctx context.Context, srcChain ccipocr3.ChainSelector, num ccipocr3.SeqNum) (r [][]byte, err error) {
	return nil, nil
}

var _ exectypes.TokenDataReader = (*nilTokenDataReader)(nil)
