package functions

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// ExternalAdapterClient supports two endpoints:
//  1. Request (aka "lambda") for executing Functions requests via RunComputation()
//  2. Secrets (aka "fetcher") for fetching offchain secrets via FetchEncryptedSecrets()
//
// Both endpoints share the same response format.
// All methods are thread-safe.
//
//go:generate mockery --quiet --name ExternalAdapterClient --output ./mocks/ --case=underscore
type ExternalAdapterClient interface {
	RunComputation(
		ctx context.Context,
		requestId string,
		jobName string,
		subscriptionOwner string,
		subscriptionId uint64,
		nodeProvidedSecrets string,
		jsonData json.RawMessage,
	) (userResult, userError []byte, domains []string, err error)

	FetchEncryptedSecrets(ctx context.Context, encryptedSecretsUrls []byte, requestId string, jobName string) (encryptedSecrets, userError []byte, err error)
}

type externalAdapterClient struct {
	AdapterURL       url.URL
	MaxResponseBytes int64
}

var _ ExternalAdapterClient = (*externalAdapterClient)(nil)

type requestPayload struct {
	Endpoint            string       `json:"endpoint"`
	RequestId           string       `json:"requestId"`
	JobName             string       `json:"jobName"`
	SubscriptionOwner   string       `json:"subscriptionOwner"`
	SubscriptionId      string       `json:"subscriptionId"`
	NodeProvidedSecrets string       `json:"nodeProvidedSecrets"`
	Data                *requestData `json:"data"`
}

type requestData struct {
	Source          string   `json:"source"`
	Language        int      `json:"language"`
	CodeLocation    int      `json:"codeLocation"`
	Secrets         string   `json:"secrets"`
	SecretsLocation int      `json:"secretsLocation"`
	Args            []string `json:"args"`
}

type secretsPayload struct {
	Endpoint  string      `json:"endpoint"`
	RequestId string      `json:"requestId"`
	JobName   string      `json:"jobName"`
	Data      secretsData `json:"data"`
}

type secretsData struct {
	RequestType          string `json:"requestType"`
	EncryptedSecretsUrls string `json:"encryptedSecretsUrls"`
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

func NewExternalAdapterClient(adapterURL url.URL, maxResponseBytes int64) ExternalAdapterClient {
	return &externalAdapterClient{
		AdapterURL:       adapterURL,
		MaxResponseBytes: maxResponseBytes,
	}
}

func (ea *externalAdapterClient) RunComputation(
	ctx context.Context,
	requestId string,
	jobName string,
	subscriptionOwner string,
	subscriptionId uint64,
	nodeProvidedSecrets string,
	jsonData json.RawMessage,
) (userResult, userError []byte, domains []string, err error) {
	var data requestData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed unmarshalling json data")
	}

	payload := requestPayload{
		Endpoint:            "lambda",
		RequestId:           requestId,
		JobName:             jobName,
		SubscriptionOwner:   subscriptionOwner,
		SubscriptionId:      strconv.FormatUint(subscriptionId, 10),
		NodeProvidedSecrets: nodeProvidedSecrets,
		Data:                &data,
	}

	userResult, userError, domains, err = ea.request(ctx, payload, requestId, jobName, "run_computation")
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "error running computation")
	}

	return userResult, userError, domains, nil
}

func (ea *externalAdapterClient) FetchEncryptedSecrets(ctx context.Context, encryptedSecretsUrls []byte, requestId string, jobName string) (encryptedSecrets, userError []byte, err error) {
	encodedSecretsUrls := base64.StdEncoding.EncodeToString(encryptedSecretsUrls)

	data := secretsData{
		RequestType:          "fetchThresholdEncryptedSecrets",
		EncryptedSecretsUrls: encodedSecretsUrls,
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

	req, err := http.NewRequestWithContext(ctx, "POST", ea.AdapterURL.String(), bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "error constructing external adapter request")
	}
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		promEAClientErrors.WithLabelValues(label).Inc()
		return nil, nil, nil, errors.Wrap(err, "error during external adapter request")
	}
	defer resp.Body.Close()

	source := http.MaxBytesReader(nil, resp.Body, ea.MaxResponseBytes)
	body, err := io.ReadAll(source)
	elapsed := time.Since(start)
	promEAClientLatency.WithLabelValues(label).Set(elapsed.Seconds())
	if err != nil {
		promEAClientErrors.WithLabelValues(label).Inc()
		return nil, nil, nil, errors.Wrap(err, "error reading external adapter response")
	}

	var eaResp response
	err = json.Unmarshal(body, &eaResp)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, fmt.Sprintf("error parsing external adapter response %s", body))
	}

	if resp.StatusCode != http.StatusOK || eaResp.StatusCode != http.StatusOK {
		return nil, nil, nil, fmt.Errorf("external adapter responded with error code %d", eaResp.StatusCode)
	}

	if eaResp.Data == nil {
		return nil, nil, nil, errors.New("external adapter response data was empty")
	}

	switch eaResp.Result {
	case "error":
		userError, err = utils.TryParseHex(eaResp.Data.Error)
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "error decoding userError hex string")
		}
		return nil, userError, eaResp.Data.Domains, nil
	case "success":
		userResult, err = utils.TryParseHex(eaResp.Data.Result)
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "error decoding result hex string")
		}
		return userResult, nil, eaResp.Data.Domains, nil
	default:
		return nil, nil, nil, fmt.Errorf("unexpected result in response: '%+v'", eaResp.Result)
	}
}
