package agentcommerce

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type X402Submission struct {
	RequestHash string
	PaymentID   string
	Phase       SettlementPhase
}

type X402Client interface {
	SubmitX402(context.Context, RailSettlementRequest) (X402Submission, error)
	ObserveX402(context.Context, RailSettlementRequest, string) (RailSettlementReceipt, error)
}

type X402Authorizer interface {
	AuthorizeX402(context.Context, X402Challenge, string) (string, error)
}

type X402Challenge struct {
	X402Version int `json:"x402Version"`
	Accepted    struct {
		Scheme  string `json:"scheme"`
		Network string `json:"network"`
		Asset   string `json:"asset"`
		Amount  string `json:"amount"`
		PayTo   string `json:"payTo"`
	} `json:"accepted"`
}

type HTTPX402Client struct {
	client     *http.Client
	endpoint   *url.URL
	authorizer X402Authorizer
}

func NewHTTPX402Client(client *http.Client, endpoint string, authorizer X402Authorizer) (*HTTPX402Client, error) {
	if client == nil || authorizer == nil {
		return nil, errors.New("x402 HTTP client and authorizer are required")
	}
	parsed, err := url.Parse(endpoint)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return nil, errors.New("valid x402 endpoint is required")
	}
	return &HTTPX402Client{client: client, endpoint: parsed, authorizer: authorizer}, nil
}

func (c *HTTPX402Client) SubmitX402(ctx context.Context, request RailSettlementRequest) (X402Submission, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return X402Submission{}, err
	}
	initial, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return X402Submission{}, err
	}
	initial.Header.Set("Content-Type", "application/json")
	response, err := c.client.Do(initial)
	if err != nil {
		return X402Submission{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusPaymentRequired {
		return X402Submission{}, fmt.Errorf("x402 service returned %d before payment", response.StatusCode)
	}
	challenge, err := parseX402Challenge(response.Header.Get("PAYMENT-REQUIRED"), request)
	if err != nil {
		return X402Submission{}, err
	}
	requestHash, _ := SettlementRequestHash(request)
	paymentHeader, err := c.authorizer.AuthorizeX402(ctx, challenge, requestHash)
	if err != nil {
		return X402Submission{}, err
	}
	if paymentHeader == "" {
		return X402Submission{}, errors.New("x402 authorizer returned empty payment")
	}
	retry, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return X402Submission{}, err
	}
	retry.Header.Set("Content-Type", "application/json")
	retry.Header.Set("PAYMENT-SIGNATURE", paymentHeader)
	paid, err := c.client.Do(retry)
	if err != nil {
		return X402Submission{}, err
	}
	defer paid.Body.Close()
	if paid.StatusCode != http.StatusAccepted {
		return X402Submission{}, fmt.Errorf("x402 payment rejected with status %d", paid.StatusCode)
	}
	var result X402Submission
	if err := decodeLimitedJSON(paid.Body, &result); err != nil {
		return X402Submission{}, err
	}
	return result, nil
}

func (c *HTTPX402Client) ObserveX402(ctx context.Context, request RailSettlementRequest, paymentID string) (RailSettlementReceipt, error) {
	if paymentID == "" {
		return RailSettlementReceipt{}, errors.New("x402 payment id is required")
	}
	endpoint := *c.endpoint
	query := endpoint.Query()
	query.Set("payment_id", paymentID)
	endpoint.RawQuery = query.Encode()
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return RailSettlementReceipt{}, err
	}
	response, err := c.client.Do(httpRequest)
	if err != nil {
		return RailSettlementReceipt{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return RailSettlementReceipt{}, fmt.Errorf("x402 observation returned %d", response.StatusCode)
	}
	var receipt RailSettlementReceipt
	if err := decodeLimitedJSON(response.Body, &receipt); err != nil {
		return RailSettlementReceipt{}, err
	}
	return receipt, nil
}

func parseX402Challenge(encoded string, request RailSettlementRequest) (X402Challenge, error) {
	if encoded == "" {
		return X402Challenge{}, errors.New("PAYMENT-REQUIRED header is missing")
	}
	payload, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return X402Challenge{}, errors.New("malformed x402 payment challenge")
	}
	var challenge X402Challenge
	if err := json.Unmarshal(payload, &challenge); err != nil {
		return X402Challenge{}, errors.New("malformed x402 payment challenge")
	}
	if challenge.X402Version != 2 || challenge.Accepted.Scheme != "exact" || challenge.Accepted.Network != request.Network || challenge.Accepted.Asset != request.Asset || challenge.Accepted.Amount != request.Amount || challenge.Accepted.PayTo != request.Destination {
		return X402Challenge{}, errors.New("x402 challenge does not match settlement request")
	}
	return challenge, nil
}

func decodeLimitedJSON(reader io.Reader, target any) error {
	decoder := json.NewDecoder(io.LimitReader(reader, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return errors.New("malformed x402 response")
	}
	return nil
}

type X402SettlementAdapter struct{ client X402Client }

func NewX402SettlementAdapter(client X402Client) (*X402SettlementAdapter, error) {
	if client == nil {
		return nil, errors.New("x402 client is required")
	}
	return &X402SettlementAdapter{client: client}, nil
}

func (a *X402SettlementAdapter) Submit(ctx context.Context, req RailSettlementRequest) (X402Submission, error) {
	if req.Rail != "x402" {
		return X402Submission{}, errors.New("x402 adapter requires x402 rail")
	}
	wantHash, err := SettlementRequestHash(req)
	if err != nil {
		return X402Submission{}, err
	}
	result, err := a.client.SubmitX402(ctx, req)
	if err != nil {
		return X402Submission{}, fmt.Errorf("submit x402 payment: %w", err)
	}
	if result.RequestHash != wantHash || result.PaymentID == "" {
		return X402Submission{}, errors.New("x402 submission identity mismatch")
	}
	if result.Phase != SettlementSubmitted && result.Phase != SettlementBroadcast {
		return X402Submission{}, errors.New("x402 submission returned an unsupported phase")
	}
	return result, nil
}

func (a *X402SettlementAdapter) Observe(ctx context.Context, req RailSettlementRequest, paymentID string) (RailSettlementReceipt, error) {
	wantHash, err := SettlementRequestHash(req)
	if err != nil {
		return RailSettlementReceipt{}, err
	}
	receipt, err := a.client.ObserveX402(ctx, req, paymentID)
	if err != nil {
		return RailSettlementReceipt{}, err
	}
	if receipt.RequestHash != wantHash || receipt.Rail != "x402" || receipt.ExternalID != paymentID {
		return RailSettlementReceipt{}, errors.New("x402 observation identity mismatch")
	}
	if receipt.ObservedAt.IsZero() {
		return RailSettlementReceipt{}, errors.New("invalid x402 observation time")
	}
	return SealSettlementReceipt(receipt)
}
