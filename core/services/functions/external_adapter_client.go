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

type ExternalAdapterInterface struct {
	AdapterURL url.URL
}

type Payload struct {
	RequestId            string `json:"requestId"`
	JobName              string `json:"jobName"`
	EncryptedSecretsUrls string `json:"encryptedSecretsUrls"`
}

type Response struct {
	Result     string       `json:"result"`
	Data       ResponseData `json:"data"`
	StatusCode int          `json:"statusCode"`
}

type ResponseData struct {
	Result      string `json:"result"`
	Error       string `json:"error"`
	ErrorString string `json:"errorString"`
}

func (ea ExternalAdapterInterface) FetchEncryptedSecrets(ctx context.Context, encryptedSecretsUrls []byte, requestId string, jobName string) (encryptedSecrets, userError []byte, err error) {
	encodedSecretsUrls := base64.StdEncoding.EncodeToString(encryptedSecretsUrls)

	payload := Payload{
		RequestId:            requestId,
		JobName:              jobName,
		EncryptedSecretsUrls: encodedSecretsUrls,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error constructing adapter encrypted secrets fetch payload")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", ea.AdapterURL.JoinPath("/fetcher").String(), bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, nil, errors.Wrap(err, "error constructing external adapter encrypted secrets fetch request")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error during external adapter encrypted secrets fetch request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error reading external adapter encrypted secrets fetch response")
	}

	var apiResponse Response
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("error parsing external adapter encrypted secrets fetch response %s", body))
	}

	switch apiResponse.Result {
	case "error":
		userError, err = utils.TryParseHex(apiResponse.Data.Error)
		if err != nil {
			return nil, nil, errors.Wrap(err, "error decoding userError hex string")
		}
		return nil, userError, nil
	case "success":
		encryptedSecrets, err = utils.TryParseHex(apiResponse.Data.Result)
		if err != nil {
			return nil, nil, errors.Wrap(err, "error decoding encryptedSecrets hex string")
		}
		return encryptedSecrets, nil, nil
	default:
		return nil, nil, fmt.Errorf("unexpected response %s", apiResponse.Result)
	}
}
