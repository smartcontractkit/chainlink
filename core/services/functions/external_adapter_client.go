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

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type ExternalAdapterClient struct {
	AdapterURL       url.URL
	MaxResponseBytes int64
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
	Result      string `json:"result"`
	Error       string `json:"error"`
	ErrorString string `json:"errorString"`
}

func (ea ExternalAdapterClient) FetchEncryptedSecrets(ctx context.Context, encryptedSecretsUrls []byte, requestId string, jobName string) (encryptedSecrets, userError []byte, err error) {
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

	encryptedSecrets, userError, err = ea.request(ctx, payload, requestId, jobName)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error fetching encrypted secrets")
	}

	return encryptedSecrets, userError, nil
}

func (ea ExternalAdapterClient) request(
	ctx context.Context,
	payload interface{},
	requestId string,
	jobName string,
) (userResult, userError []byte, err error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error constructing external adapter request payload")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", ea.AdapterURL.String(), bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, nil, errors.Wrap(err, "error constructing external adapter request")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error during external adapter request")
	}
	defer resp.Body.Close()

	source := http.MaxBytesReader(nil, resp.Body, ea.MaxResponseBytes)
	body, err := io.ReadAll(source)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error reading external adapter response")
	}

	var eaResp response
	err = json.Unmarshal(body, &eaResp)
	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("error parsing external adapter response %s", body))
	}

	if eaResp.StatusCode != 200 {
		return nil, nil, fmt.Errorf("external adapter responded with error code %d", eaResp.StatusCode)
	}

	if eaResp.Data == nil {
		return nil, nil, errors.New("external adapter response data was empty")
	}

	switch eaResp.Result {
	case "error":
		userError, err = utils.TryParseHex(eaResp.Data.Error)
		if err != nil {
			return nil, nil, errors.Wrap(err, "error decoding userError hex string")
		}
		return nil, userError, nil
	case "success":
		userResult, err = utils.TryParseHex(eaResp.Data.Result)
		if err != nil {
			return nil, nil, errors.Wrap(err, "error decoding result hex string")
		}
		return userResult, nil, nil
	default:
		return nil, nil, fmt.Errorf("unexpected result in response: '%+v'", eaResp.Result)
	}
}
