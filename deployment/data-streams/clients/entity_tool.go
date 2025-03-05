package clients

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type GetOverridesRequest struct {
	asset   string
	quote   string
	product *string
}

type GetAssetEAsRequest struct {
	asset   string
	quote   string
	product *string
}

type (
	GetOverridesResponse map[string]string
	GetAssetEAsResponse  []string
)

type EntityToolClient interface {
	GetOverrides(ctx context.Context, in *GetOverridesRequest) (*GetOverridesResponse, error)
	GetAssetEAs(ctx context.Context, in *GetAssetEAsRequest) (*GetAssetEAsResponse, error)
}

type EntityToolClientImpl struct {
	baseURL string
	client  *http.Client
}

func NewEntityToolClient(baseURL string, client *http.Client) EntityToolClient {
	if client == nil {
		client = http.DefaultClient
	}
	return &EntityToolClientImpl{
		baseURL: baseURL,
		client:  client,
	}
}

type feedBuildObjectResponse struct {
	ExternalAdapterRequestParams struct {
		Overrides map[string]string `json:"overrides"`
	} `json:"externalAdapterRequestParams"`
	APIs []string `json:"apis"`
}

func (c *EntityToolClientImpl) getFeedBuildObjectResponse(ctx context.Context, asset, quote, product string) (*feedBuildObjectResponse, error) {
	endpoint := c.baseURL + "/api/feed-build-object"
	values := url.Values{}
	values.Add("base", asset)
	values.Add("quote", quote)
	values.Add("product", product)
	values.Add("configType", "DATA_STREAMS")
	reqURL := endpoint + "?" + values.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, errors.New("failed to create request: " + err.Error())
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.New("request failed: " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code: " + http.StatusText(resp.StatusCode))
	}
	var result feedBuildObjectResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errors.New("failed to decode response: " + err.Error())
	}
	return &result, nil
}

func (c *EntityToolClientImpl) GetOverrides(ctx context.Context, in *GetOverridesRequest) (*GetOverridesResponse, error) {
	product := "crypto"
	if in.product != nil {
		product = *in.product
	}
	res, err := c.getFeedBuildObjectResponse(ctx, in.asset, in.quote, product)
	if err != nil {
		return nil, err
	}
	if res.ExternalAdapterRequestParams.Overrides == nil {
		return nil, errors.New("no APIs found for the given asset and quote")
	}
	overrides := GetOverridesResponse(res.ExternalAdapterRequestParams.Overrides)
	return &overrides, nil
}

func (c *EntityToolClientImpl) GetAssetEAs(ctx context.Context, in *GetAssetEAsRequest) (*GetAssetEAsResponse, error) {
	product := "crypto"
	if in.product != nil {
		product = *in.product
	}
	res, err := c.getFeedBuildObjectResponse(ctx, in.asset, in.quote, product)
	if err != nil {
		return nil, err
	}
	if len(res.APIs) == 0 {
		return nil, errors.New("no APIs found for the given asset and quote")
	}
	eas := GetAssetEAsResponse(res.APIs)
	return &eas, nil
}

func NewGetOverridesRequest(asset, quote string, product ...string) *GetOverridesRequest {
	var prod *string
	if len(product) > 0 {
		prod = &product[0]
	}
	return &GetOverridesRequest{
		asset:   asset,
		quote:   quote,
		product: prod,
	}
}

func NewGetAssetEAsRequest(asset, quote string, product ...string) *GetAssetEAsRequest {
	var prod *string
	if len(product) > 0 {
		prod = &product[0]
	}
	return &GetAssetEAsRequest{
		asset:   asset,
		quote:   quote,
		product: prod,
	}
}
