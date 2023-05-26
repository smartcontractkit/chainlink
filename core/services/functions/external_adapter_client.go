package functions

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type ExternalAdapterInterface struct {
	AdapterURL string
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

func (ea ExternalAdapterInterface) FetchEncryptedSecrets(encryptedSecretsUrls []byte, requestId string, jobName string) (encryptedSecrets, userError []byte, err error) {
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

	req, err := http.NewRequest("POST", ea.AdapterURL+"/fetcher", bytes.NewBuffer(jsonPayload))
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
		return nil, nil, errors.Wrap(err, "error parsing external adapter encrypted secrets fetch response")
	}

	switch apiResponse.Result {
	case "error":
		userError, err = hex.DecodeString(apiResponse.Data.Error[2:])
		if err != nil {
			return nil, nil, errors.Wrap(err, "error decoding userError hex string")
		}
		return nil, userError, nil
	case "success":
		encryptedSecrets, err = hex.DecodeString(apiResponse.Data.Result[2:])
		if err != nil {
			return nil, nil, errors.Wrap(err, "error decoding encryptedSecrets hex string")
		}
		return encryptedSecrets, nil, nil
	default:
		return nil, nil, fmt.Errorf("unexpected response %s", apiResponse.Result)
	}
}
