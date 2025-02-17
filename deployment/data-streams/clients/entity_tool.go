package clients

import (
	"context"
	"encoding/json"
	"fmt"
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
	baseUrl string
	client  *http.Client
}

// If client is nil, http.DefaultClient is used.
func NewEntityToolClient(baseUrl string, client *http.Client) EntityToolClient {
	if client == nil {
		client = http.DefaultClient
	}
	return &EntityToolClientImpl{
		baseUrl: baseUrl,
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
	endpoint := fmt.Sprintf("%s/api/feed-build-object", c.baseUrl)

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

	reqURL := fmt.Sprintf("%s?%s", endpoint, values.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetOverrides request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GetOverrides request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetOverrides unexpected status code: %d", resp.StatusCode)
	}

	var result feedBuildObjectResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode GetOverrides response: %w", err)
	}

	if result.ExternalAdapterRequestParams.Overrides == nil {
		return nil, fmt.Errorf("no APIs found for the given asset and quote")
	}

	overrides := GetOverridesResponse(result.ExternalAdapterRequestParams.Overrides)
	return &overrides, nil
}

func (c *EntityToolClientImpl) GetAssetEAs(ctx context.Context, in *GetAssetEAsRequest) (*GetAssetEAsResponse, error) {
	endpoint := fmt.Sprintf("%s/api/feed-build-object", c.baseUrl)

	product := "crypto"
	if in.product != nil {
		product = *in.product
	}

	values := url.Values{}
	values.Add("base", in.asset)
	values.Add("quote", in.quote)
	values.Add("product", product)
	values.Add("configType", "DATA_STREAMS")

	reqURL := fmt.Sprintf("%s?%s", endpoint, values.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetAssetEAs request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GetAssetEAs request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetAssetEAs unexpected status code: %d", resp.StatusCode)
	}

	var result feedBuildObjectResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode GetAssetEAs response: %w", err)
	}

	if len(result.APIs) == 0 {
		return nil, fmt.Errorf("no APIs found for the given asset and quote")
	}

	eas := GetAssetEAsResponse(result.APIs)
	return &eas, nil
}
