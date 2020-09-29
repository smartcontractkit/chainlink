package pipeline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type HTTPTask struct {
	BaseTask    `mapstructure:",squash"`
	Method      string
	URL         models.WebURL
	RequestData HttpRequestData `json:"requestData"`

	config Config
}

type PossibleErrorResponses struct {
	Error        string `json:"error"`
	ErrorMessage string `json:"errorMessage"`
}

var _ Task = (*HTTPTask)(nil)

func (t *HTTPTask) Type() TaskType {
	return TaskTypeHTTP
}

func (t *HTTPTask) Run(taskRun TaskRun, inputs []Result) Result {
	if len(inputs) > 0 {
		return Result{Error: errors.Wrapf(ErrWrongInputCardinality, "HTTPTask requires 0 inputs")}
	}

	buf := &bytes.Buffer{}
	var bodyReader io.Reader
	if t.RequestData != nil {
		bodyBytes, err := json.Marshal(t.RequestData)
		if err != nil {
			return Result{Error: errors.Wrap(err, "failed to encode request body as JSON")}
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	request, err := http.NewRequest(t.Method, t.URL.String(), bodyReader)
	if err != nil {
		return Result{Error: errors.Wrap(err, "failed to create http.Request")}
	}
	request.Header.Set("Content-Type", "application/json")

	var client http.Client
	response, err := client.Do(request)
	if err != nil {
		return Result{Error: errors.Wrapf(err, "could not fetch answer from %s with payload '%s'", t.URL.String(), t.RequestData)}
	}
	defer logger.ErrorIfCalling(response.Body.Close)

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return Result{Error: errors.Wrapf(err, "could not read response body")}
	}

	if response.StatusCode >= 400 {
		maybeErr := bestEffortExtractError(responseBytes)
		return Result{Error: errors.Errorf("got error from %s: (status %s) %s", t.URL.String(), response.Status, maybeErr)}
	}

	fmt.Println("ASDF ~>", string(responseBytes))

	logger.Debugw("HTTP task got response",
		"response", string(responseBytes),
		"url", t.URL.String(),
	)
	return Result{Value: responseBytes}
}

func bestEffortExtractError(responseBytes []byte) string {
	var resp PossibleErrorResponses
	err := json.Unmarshal(responseBytes, &resp)
	if err != nil {
		return ""
	}
	if resp.Error != "" {
		return resp.Error
	} else if resp.ErrorMessage != "" {
		return resp.ErrorMessage
	}
	return string(responseBytes)
}
