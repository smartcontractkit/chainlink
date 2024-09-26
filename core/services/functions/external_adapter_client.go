package functions

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/hex"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
)

// ExternalAdapterClient supports two endpoints:
//  1. Request (aka "lambda") for executing Functions requests via RunComputation()
//  2. Secrets (aka "fetcher") for fetching offchain secrets via FetchEncryptedSecrets()
//
// Both endpoints share the same response format.
// All methods are thread-safe.
type ExternalAdapterClient interface {
	RunComputation(
		ctx context.Context,
		requestId string,
		jobName string,
		subscriptionOwner string,
		subscriptionId uint64,
		flags RequestFlags,
		nodeProvidedSecrets string,
		requestData *RequestData,
	) (userResult, userError []byte, domains []string, err error)

	FetchEncryptedSecrets(ctx context.Context, encryptedSecretsUrls []byte, requestId string, jobName string) (encryptedSecrets, userError []byte, err error)
}

type externalAdapterClient struct {
	adapterURL             url.URL
	maxResponseBytes       int64
	maxRetries             int
	exponentialBackoffBase time.Duration
}

var _ ExternalAdapterClient = (*externalAdapterClient)(nil)

type BridgeAccessor interface {
	NewExternalAdapterClient(context.Context) (ExternalAdapterClient, error)
}

type bridgeAccessor struct {
	bridgeORM              bridges.ORM
	bridgeName             string
	maxResponseBytes       int64
	maxRetries             int
	exponentialBackoffBase time.Duration
}

var _ BridgeAccessor = (*bridgeAccessor)(nil)

type requestPayload struct {
	Endpoint            string       `json:"endpoint"`
	RequestId           string       `json:"requestId"`
	JobName             string       `json:"jobName"`
	SubscriptionOwner   string       `json:"subscriptionOwner"`
	SubscriptionId      uint64       `json:"subscriptionId"`
	Flags               RequestFlags `json:"flags"` // marshalled as an array of numbers
	NodeProvidedSecrets string       `json:"nodeProvidedSecrets"`
	Data                *RequestData `json:"data"`
}

type secretsPayload struct {
	Endpoint  string      `json:"endpoint"`
	RequestId string      `json:"requestId"`
	JobName   string      `json:"jobName"`
	Data      secretsData `json:"data"`
}

type secretsData struct {
	RequestType          string `json:"requestType"`
	EncryptedSecretsUrls []byte `json:"encryptedSecretsUrls"`
}

type response struct {
	Result     string        `json:"result"`
	Data       *responseData `json:"data"`
	StatusCode int           `json:"statusCode"`
}

type responseData struct {
	Result      string   `json:"result"`
	Error       string   `json:"error"`
	ErrorString string   `json:"errorString"`
	Domains     []string `json:"domains"`
}

var (
	promEAClientLatency = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "functions_external_adapter_client_latency",
		Help: "Functions EA client latency in seconds scoped by endpoint",
	},
		[]string{"name"},
	)
	promEAClientErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "functions_external_adapter_client_errors_total",
		Help: "Functions EA client error count scoped by endpoint",
	},
		[]string{"name"},
	)
)

func NewExternalAdapterClient(adapterURL url.URL, maxResponseBytes int64, maxRetries int, exponentialBackoffBase time.Duration) ExternalAdapterClient {
	return &externalAdapterClient{
		adapterURL:             adapterURL,
		maxResponseBytes:       maxResponseBytes,
		maxRetries:             maxRetries,
		exponentialBackoffBase: exponentialBackoffBase,
	}
}

func (ea *externalAdapterClient) RunComputation(
	ctx context.Context,
	requestId string,
	jobName string,
	subscriptionOwner string,
	subscriptionId uint64,
	flags RequestFlags,
	nodeProvidedSecrets string,
	requestData *RequestData,
) (userResult, userError []byte, domains []string, err error) {
	requestData.Secrets = nil // secrets are passed in nodeProvidedSecrets

	payload := requestPayload{
		Endpoint:            "lambda",
		RequestId:           requestId,
		JobName:             jobName,
		SubscriptionOwner:   subscriptionOwner,
		SubscriptionId:      subscriptionId,
		Flags:               flags,
		NodeProvidedSecrets: nodeProvidedSecrets,
		Data:                requestData,
	}

	userResult, userError, domains, err = ea.request(ctx, payload, requestId, jobName, "run_computation")
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "error running computation")
	}

	return userResult, userError, domains, nil
}

func (ea *externalAdapterClient) FetchEncryptedSecrets(ctx context.Context, encryptedSecretsUrls []byte, requestId string, jobName string) (encryptedSecrets, userError []byte, err error) {
	data := secretsData{
		RequestType:          "fetchThresholdEncryptedSecrets",
		EncryptedSecretsUrls: encryptedSecretsUrls,
	}

	payload := secretsPayload{
		Endpoint:  "fetcher",
		RequestId: requestId,
		JobName:   jobName,
		Data:      data,
	}

	encryptedSecrets, userError, _, err = ea.request(ctx, payload, requestId, jobName, "fetch_secrets")
	if err != nil {
		return nil, nil, errors.Wrap(err, "error fetching encrypted secrets")
	}

	return encryptedSecrets, userError, nil
}

func (ea *externalAdapterClient) request(
	ctx context.Context,
	payload interface{},
	requestId string,
	jobName string,
	label string,
) (userResult, userError []byte, domains []string, err error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "error constructing external adapter request payload")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", ea.adapterURL.String(), bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "error constructing external adapter request")
	}
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()

	// retry will only happen on a 5XX error response code (except 501)
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = ea.maxRetries
	retryClient.RetryWaitMin = ea.exponentialBackoffBase

	client := retryClient.StandardClient()
	resp, err := client.Do(req)
	if err != nil {
		promEAClientErrors.WithLabelValues(label).Inc()
		return nil, nil, nil, errors.Wrap(err, "error during external adapter request")
	}
	defer resp.Body.Close()

	source := http.MaxBytesReader(nil, resp.Body, ea.maxResponseBytes)
	body, err := io.ReadAll(source)
	elapsed := time.Since(start)
	promEAClientLatency.WithLabelValues(label).Set(elapsed.Seconds())
	if err != nil {
		promEAClientErrors.WithLabelValues(label).Inc()
		return nil, nil, nil, errors.Wrap(err, "error reading external adapter response")
	}

	if resp.StatusCode != http.StatusOK {
		promEAClientErrors.WithLabelValues(label).Inc()
		return nil, nil, nil, fmt.Errorf("external adapter responded with HTTP %d, body: %s", resp.StatusCode, body)
	}

	var eaResp response
	err = json.Unmarshal(body, &eaResp)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, fmt.Sprintf("error parsing external adapter response %s", body))
	}

	if eaResp.StatusCode != http.StatusOK {
		return nil, nil, nil, fmt.Errorf("external adapter invalid StatusCode %d", eaResp.StatusCode)
	}

	if eaResp.Data == nil {
		return nil, nil, nil, errors.New("external adapter response data was empty")
	}

	switch eaResp.Result {
	case "error":
		userError, err = hex.DecodeString(eaResp.Data.Error)
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "error decoding userError hex string")
		}
		return nil, userError, eaResp.Data.Domains, nil
	case "success":
		userResult, err = hex.DecodeString(eaResp.Data.Result)
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "error decoding result hex string")
		}
		return userResult, nil, eaResp.Data.Domains, nil
	default:
		return nil, nil, nil, fmt.Errorf("unexpected result in response: '%+v'", eaResp.Result)
	}
}

func NewBridgeAccessor(bridgeORM bridges.ORM, bridgeName string, maxResponseBytes int64, maxRetries int, exponentialBackoffBase time.Duration) BridgeAccessor {
	return &bridgeAccessor{
		bridgeORM:              bridgeORM,
		bridgeName:             bridgeName,
		maxResponseBytes:       maxResponseBytes,
		maxRetries:             maxRetries,
		exponentialBackoffBase: exponentialBackoffBase,
	}
}

func (b *bridgeAccessor) NewExternalAdapterClient(ctx context.Context) (ExternalAdapterClient, error) {
	bridge, err := b.bridgeORM.FindBridge(ctx, bridges.BridgeName(b.bridgeName))
	if err != nil {
		return nil, err
	}
	return NewExternalAdapterClient(url.URL(bridge.URL), b.maxResponseBytes, b.maxRetries, b.exponentialBackoffBase), nil
}
