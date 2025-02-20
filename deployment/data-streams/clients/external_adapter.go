package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ExternalAdapterRequest struct {
	AdapterURL string `json:"adapterUrl"`
	Asset      string `json:"asset"`
	Quote      string `json:"quote"`
	Endpoint   string `json:"endpoint"`
}

type ExternalAdapterResponse struct {
	Pass         bool                   `json:"pass"`
	StatusCode   int                    `json:"statusCode"`
	Data         map[string]interface{} `json:"data,omitempty"`
	ErrorMessage string                 `json:"errorMessage"`
}

type ExternalAdapterClient interface {
	Query(ctx context.Context, req ExternalAdapterRequest) (ExternalAdapterResponse, error)
}

type externalAdapterClientImpl struct {
	client *http.Client
}

// If client is nil, http.DefaultClient is used.
func NewExternalAdapterClient(client *http.Client) ExternalAdapterClient {
	if client == nil {
		client = http.DefaultClient
	}
	return &externalAdapterClientImpl{
		client: client,
	}
}

func (c *externalAdapterClientImpl) Query(ctx context.Context, reqData ExternalAdapterRequest) (ExternalAdapterResponse, error) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"from":     reqData.Asset,
			"to":       reqData.Quote,
			"endpoint": reqData.Endpoint,
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return ExternalAdapterResponse{}, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, reqData.AdapterURL, bytes.NewBuffer(body))
	if err != nil {
		return ExternalAdapterResponse{}, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return ExternalAdapterResponse{
			Pass:         false,
			StatusCode:   http.StatusInternalServerError,
			ErrorMessage: err.Error(),
		}, nil
	}
	defer resp.Body.Close()

	var result ExternalAdapterResponse
	if resp.StatusCode == http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(&result.Data); err != nil {
			return ExternalAdapterResponse{
				Pass:         false,
				StatusCode:   resp.StatusCode,
				ErrorMessage: fmt.Sprintf("failed to decode JSON: %v", err),
			}, nil
		}
		result.Pass = true
		result.StatusCode = resp.StatusCode
		result.ErrorMessage = ""
	} else {
		var errData map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&errData)
		result = ExternalAdapterResponse{
			Pass:         false,
			StatusCode:   resp.StatusCode,
			ErrorMessage: resp.Status,
			Data:         errData,
		}
	}

	return result, nil
}
