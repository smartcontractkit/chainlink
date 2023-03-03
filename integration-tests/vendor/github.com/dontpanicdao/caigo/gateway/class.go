package gateway

import (
	"context"
	"net/http"
	"net/url"

	"github.com/dontpanicdao/caigo/types"
)

func (sg *Gateway) ClassByHash(ctx context.Context, hash string) (*types.ContractClass, error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_class_by_hash", nil)
	if err != nil {
		return nil, err
	}

	appendQueryValues(req, url.Values{
		"classHash": []string{hash},
	})

	var resp types.ContractClass
	return &resp, sg.do(req, &resp)
}

func (sg *Gateway) ClassHashAt(ctx context.Context, address string) (*types.Felt, error) {
	req, err := sg.newRequest(ctx, http.MethodGet, "/get_class_hash_at", nil)
	if err != nil {
		return nil, err
	}

	appendQueryValues(req, url.Values{
		"contractAddress": []string{address},
	})

	var resp types.Felt
	return &resp, sg.do(req, &resp)
}

func (sg *Gateway) Class(context.Context, string) (*types.ContractClass, error) {
	panic("not implemented")
}

func (sg *Gateway) ClassAt(context.Context, string) (*types.ContractClass, error) {
	panic("not implemented")
}
