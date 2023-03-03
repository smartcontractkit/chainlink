package gateway

import (
	"context"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

type StorageAtOptions struct {
	BlockNumber int    `url:"blockNumber,omitempty"`
	BlockHash   string `url:"blockHash,omitempty"`
}

// Get a storage slots value.
//
// [Reference](https://github.com/starkware-libs/cairo-lang/blob/fc97bdd8322a7df043c87c371634b26c15ed6cee/src/starkware/starknet/services/api/feeder_gateway/feeder_gateway_client.py#L70)
func (sg *Gateway) StorageAt(ctx context.Context, address string, key string, opts *StorageAtOptions) (string, error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_storage_at", nil)
	if err != nil {
		return "", err
	}

	appendQueryValues(req, url.Values{
		"contractAddress": []string{address},
		"key":             []string{key},
	})

	if opts != nil {
		vs, err := query.Values(opts)
		if err != nil {
			return "", err
		}
		appendQueryValues(req, vs)
	}

	var resp string
	return resp, sg.do(req, &resp)
}
