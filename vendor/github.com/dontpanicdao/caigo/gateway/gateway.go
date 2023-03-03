package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	INVOKE         string = "INVOKE_FUNCTION"
	DEPLOY         string = "DEPLOY"
	DEPLOY_ACCOUNT string = "DEPLOY_ACCOUNT"
	DECLARE        string = "DECLARE"
	GOERLI_ID      string = "SN_GOERLI"
	MAINNET_ID     string = "SN_MAIN"
	LOCAL_BASE     string = "http://localhost:5050"
	GOERLI_BASE    string = "https://alpha4.starknet.io"
	MAINNET_BASE   string = "https://alpha-mainnet.starknet.io"
)

type doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Gateway struct {
	Base         string `json:"base"`
	Feeder       string `json:"feeder"`
	Gateway      string `json:"gateway"`
	ChainId      string `json:"chainId"`
	client       doer
	errorHandler func(e error) error
}

/*
Instantiate a new StarkNet Gateway client
- defaults to the GOERLI endpoints
*/
func NewClient(opts ...Option) *Gateway {
	gopts := options{
		chainID: GOERLI_ID,
		client:  http.DefaultClient,
	}

	for _, opt := range opts {
		opt.apply(&gopts)
	}

	switch id := strings.ToLower(gopts.chainID); {
	case strings.Contains(id, "main"):
		gopts.chainID = MAINNET_ID
		if gopts.baseUrl == "" {
			gopts.baseUrl = MAINNET_BASE
		}
	case strings.Contains(id, "local"):
		fallthrough
	case strings.Contains(id, "dev"):
		gopts.chainID = GOERLI_ID
		if gopts.baseUrl == "" {
			gopts.baseUrl = LOCAL_BASE
		}
	default:
		gopts.chainID = GOERLI_ID
		if gopts.baseUrl == "" {
			gopts.baseUrl = GOERLI_BASE
		}
	}

	return &Gateway{
		Base:         gopts.baseUrl,
		Feeder:       gopts.baseUrl + "/feeder_gateway",
		Gateway:      gopts.baseUrl + "/gateway",
		ChainId:      gopts.chainID,
		client:       gopts.client,
		errorHandler: gopts.errorHandler,
	}
}

func (sg *Gateway) newRequest(
	ctx context.Context, method, endpoint string, body interface{},
) (*http.Request, error) {
	url := sg.Feeder + endpoint
	if strings.HasSuffix(endpoint, "add_transaction") {
		url = sg.Gateway + endpoint
	}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
		req.Body = io.NopCloser(bytes.NewBuffer(data))
		req.Header.Add("Content-Type", "application/json; charset=utf")
	}
	return req, nil
}

type Error struct {
	StatusCode int    `json:"-"`
	Body       []byte `json:"-"`

	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface.
func (e Error) Error() string {
	return fmt.Sprintf("%d: %s %s", e.StatusCode, e.Code, e.Message)
}

// NewError creates a new Error from an API response.
func NewError(resp *http.Response) error {
	apiErr := Error{StatusCode: resp.StatusCode}
	data, err := io.ReadAll(resp.Body)
	if err == nil && data != nil {
		apiErr.Body = data
		if err := json.Unmarshal(data, &apiErr); err != nil {
			apiErr.Code = "unknown_error_format"
			apiErr.Message = string(data)
		}
	}
	return &apiErr
}

func (sg *Gateway) do(req *http.Request, v interface{}) error {
	resp, err := sg.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // nolint: errcheck

	if resp.StatusCode >= 299 {
		e := NewError(resp)
		if sg.errorHandler != nil {
			return sg.errorHandler(e)
		}
		return e
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

func appendQueryValues(req *http.Request, values url.Values) {
	q := req.URL.Query()
	for k, vs := range values {
		for _, v := range vs {
			q.Add(k, v)
		}
	}
	req.URL.RawQuery = q.Encode()
}
