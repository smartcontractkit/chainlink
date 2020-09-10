package pipeline

import (
	"net/url"

	// "github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type BridgeTask struct {
	BaseTask

	Name        string          `json:"name"`
	RequestData HttpRequestData `json:"requestData"`

	ORM                BridgeTaskORM   `json:"-"`
	defaultHTTPTimeout models.Duration `json:"-"`
}

type BridgeTaskORM interface {
	FindBridge(name models.TaskType) (models.BridgeType, error)
}

func (f *BridgeTask) Run(inputs []Result) Result {
	if len(inputs) > 0 {
		return Result{Error: errors.Wrapf(ErrWrongInputCardinality, "BridgeTask requires 0 inputs")}
	}

	url, err := f.getBridgeURLFromName()
	if err != nil {
		return Result{Error: err}
	}

	// client := &http.Client{Timeout: f.defaultHTTPTimeout.Duration(), Transport: http.DefaultTransport}
	// client.Transport = promhttp.InstrumentRoundTripperDuration(promFMResponseTime, client.Transport)
	// client.Transport = instrumentRoundTripperReponseSize(promFMResponseSize, client.Transport)

	// add an arbitrary "id" field to the request json
	// this is done in order to keep request payloads consistent in format
	// between flux monitor polling requests and http/bridge adapters
	f.RequestData["id"] = models.NewID()

	result := (&HTTPTask{URL: models.WebURL(url), Method: "POST", RequestData: f.RequestData}).Run(inputs)
	if result.Error != nil {
		return result
	}
	logger.Debugw("Bridge: fetched answer",
		"answer", result,
		"url", url.String(),
	)
	return result
}

func (f BridgeTask) getBridgeURLFromName() (url.URL, error) {
	task := models.TaskType(f.Name)
	bridge, err := f.ORM.FindBridge(task)
	if err != nil {
		return url.URL{}, err
	}
	bridgeURL := url.URL(bridge.URL)
	return bridgeURL, nil
}
