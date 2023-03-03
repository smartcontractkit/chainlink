package gateway

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

type Block struct {
	BlockHash           string               `json:"block_hash"`
	ParentBlockHash     string               `json:"parent_block_hash"`
	BlockNumber         int                  `json:"block_number"`
	StateRoot           string               `json:"state_root"`
	Status              string               `json:"status"`
	Transactions        []Transaction        `json:"transactions"`
	Timestamp           int                  `json:"timestamp"`
	TransactionReceipts []TransactionReceipt `json:"transaction_receipts"`
}

type BlockOptions struct {
	BlockNumber uint64 `url:"blockNumber,omitempty"`
	BlockHash   string `url:"blockHash,omitempty"`
}

// Gets the block information from a block ID.
//
// [Reference](https://github.com/starkware-libs/cairo-lang/blob/f464ec4797361b6be8989e36e02ec690e74ef285/src/starkware/starknet/services/api/feeder_gateway/feeder_gateway_client.py#L27-L31)
func (sg *Gateway) Block(ctx context.Context, opts *BlockOptions) (*Block, error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_block", nil)
	if err != nil {
		return nil, err
	}
	if opts != nil {
		vs, err := query.Values(opts)
		if err != nil {
			return nil, err
		}
		appendQueryValues(req, vs)
	}

	var resp Block
	return &resp, sg.do(req, &resp)
}

func (sg *Gateway) BlockHashByID(ctx context.Context, id uint64) (block string, err error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_block_hash_by_id", nil)
	if err != nil {
		return "", err
	}

	appendQueryValues(req, url.Values{
		"blockId": []string{fmt.Sprint(id)},
	})

	var resp string
	return resp, sg.do(req, &resp)
}

func (sg *Gateway) BlockIDByHash(ctx context.Context, hash string) (block uint64, err error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_block_id_by_hash", nil)
	if err != nil {
		return 0, err
	}

	appendQueryValues(req, url.Values{
		"blockHash": []string{hash},
	})

	var resp uint64
	return resp, sg.do(req, &resp)
}

func (sg *Gateway) BlockByHash(context.Context, string, string) (*Block, error) {
	panic("not implemented")
}

func (sg *Gateway) BlockByNumber(context.Context, *big.Int, string) (*Block, error) {
	panic("not implemented")
}
