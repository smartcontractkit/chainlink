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

// If client is nil, http.DefaultClient is used.
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

func (c *EntityToolClientImpl) GetOverrides(ctx context.Context, in *GetOverridesRequest) (*GetOverridesResponse, error) {
	endpoint := c.baseURL + "/api/feed-build-object"

	// Default product to "crypto" if not provided.
	product := "crypto"
	if in.product != nil {
		product = *in.product
	}

	values := url.Values{}
	values.Add("base", in.asset)
	values.Add("quote", in.quote)
	values.Add("product", product)
	values.Add("configType", "DATA_STREAMS")

	reqURL := endpoint + "?" + values.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, errors.New("failed to create GetOverrides request: " + err.Error())
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.New("GetOverrides request failed: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("GetOverrides unexpected status code: " + http.StatusText(resp.StatusCode))
	}

	var result feedBuildObjectResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errors.New("failed to decode GetOverrides response: " + err.Error())
	}

	if result.ExternalAdapterRequestParams.Overrides == nil {
		return nil, errors.New("no APIs found for the given asset and quote")
	}

	overrides := GetOverridesResponse(result.ExternalAdapterRequestParams.Overrides)
	return &overrides, nil
}

func (c *EntityToolClientImpl) GetAssetEAs(ctx context.Context, in *GetAssetEAsRequest) (*GetAssetEAsResponse, error) {
	endpoint := c.baseURL + "/api/feed-build-object"

	product := "crypto"
	if in.product != nil {
		product = *in.product
	}

	values := url.Values{}
	values.Add("base", in.asset)
	values.Add("quote", in.quote)
	values.Add("product", product)
	values.Add("configType", "DATA_STREAMS")

	reqURL := endpoint + "?" + values.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, errors.New("failed to create GetAssetEAs request: " + err.Error())
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.New("GetAssetEAs request failed: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("GetAssetEAs unexpected status code: " + http.StatusText(resp.StatusCode))
	}

	var result feedBuildObjectResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errors.New("failed to decode GetAssetEAs response: " + err.Error())
	}

	if len(result.APIs) == 0 {
		return nil, errors.New("no APIs found for the given asset and quote")
	}

	eas := GetAssetEAsResponse(result.APIs)
	return &eas, nil
}
