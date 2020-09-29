package pipeline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/guregu/null.v4"
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

var _ Task = (*HTTPTask)(nil)

func (t *HTTPTask) Type() TaskType {
	return TaskTypeHTTP
}

func (t *HTTPTask) Run(inputs []Result) (result Result) {
	if len(inputs) > 0 {
		return Result{Error: errors.Wrapf(ErrWrongInputCardinality, "HTTPTask requires 0 inputs")}
	}

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

	r, err := client.Do(request)
	if err != nil {
		return Result{Error: errors.Wrapf(err, "could not fetch answer from %s with payload '%s'", t.URL.String(), t.RequestData)}
	}
	defer logger.ErrorIfCalling(r.Body.Close)

	if r.StatusCode >= 400 {
		return Result{Error: errors.Errorf("got error status code %d; unable to retrieve answer from %s", r.StatusCode, t.URL.String())}
	}

	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return Result{Error: errors.Wrapf(err, "could not read response body")}
	}

	logger.Debugw("HTTP task got response",
		"response", string(bs),
		"url", t.URL.String(),
	)
	return Result{Value: bs}
}
