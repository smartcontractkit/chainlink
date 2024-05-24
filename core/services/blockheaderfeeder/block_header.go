package blockheaderfeeder

import (
	"bytes"
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
)

type Client interface {
	BatchCallContext(ctx context.Context, b []rpc.BatchElem) error
}

type GethBlockHeaderProvider struct {
	client Client
}

func NewGethBlockHeaderProvider(client Client) *GethBlockHeaderProvider {
	return &GethBlockHeaderProvider{
		client: client,
	}
}

// RlpHeadersBatch retrieves RLP-encoded block headers
// this function is not supported for Avax because Avalanche
// block header format is different from go-ethereum types.Header.
// validation for invalid chain ID is done upstream in blockheaderfeeder.validate.go
func (p *GethBlockHeaderProvider) RlpHeadersBatch(ctx context.Context, blockRange []*big.Int) ([][]byte, error) {
	var reqs []rpc.BatchElem
	for _, num := range blockRange {
		parentBlockNum := big.NewInt(num.Int64() + 1)
		req := rpc.BatchElem{
			Method: "eth_getHeaderByNumber",
			// Get child block since it's the one that has the parent hash in its header.
			Args:   []interface{}{hexutil.EncodeBig(parentBlockNum)},
			Result: &types.Header{},
		}
		reqs = append(reqs, req)
	}
	err := p.client.BatchCallContext(ctx, reqs)
	if err != nil {
		return nil, err
	}

	var headers [][]byte
	for _, req := range reqs {
		header, ok := req.Result.(*types.Header)
		if !ok {
			return nil, errors.Errorf("received invalid type: %T", req.Result)
		}
		if header == nil {
			return nil, errors.New("invariant violation: got nil header")
		}
		headerBuffer := new(bytes.Buffer)
		err := header.EncodeRLP(headerBuffer)
		if err != nil {
			return nil, err
		}
		headers = append(headers, headerBuffer.Bytes())
	}

	return headers, nil
}
