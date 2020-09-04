package job

import (
	"fmt"
	"net/url"

	// "github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type BridgeTask struct {
	BaseTask

	BridgeName  string          `json:"name"`
	RequestData HttpRequestData `json:"requestData" gorm:"type:jsonb"`

	ORM                BridgeTaskORM `json:"-" gorm:"-"`
	defaultHTTPTimeout models.Duration
}

type BridgeTaskORM interface {
	FindBridge(name models.TaskType) (models.BridgeType, error)
}

func (f *BridgeTask) Run(inputs []Result) (out interface{}, err error) {
	if len(inputs) > 0 {
		return nil, errors.Wrapf(ErrWrongInputCardinality, "BridgeTask requires 0 inputs")
	}

	url, err := f.getBridgeURLFromName()
	if err != nil {
		return nil, err
	}

	// client := &http.Client{Timeout: f.defaultHTTPTimeout.Duration(), Transport: http.DefaultTransport}
	// client.Transport = promhttp.InstrumentRoundTripperDuration(promFMResponseTime, client.Transport)
	// client.Transport = instrumentRoundTripperReponseSize(promFMResponseSize, client.Transport)

	// add an arbitrary "id" field to the request json
	// this is done in order to keep request payloads consistent in format
	// between flux monitor polling requests and http/bridge adapters
	f.RequestData["id"] = models.NewID()

	result, err := (&HTTPTask{URL: models.WebURL(url), Method: "POST", RequestData: f.RequestData}).Fetch()
	if err != nil {
		return nil, err
	}
	logger.Debugw(
		fmt.Sprintf("Fetched answer", result, url.String()),
		"answer", result,
		"url", url.String(),
	)
	return result, nil
}

func (f BridgeTask) getBridgeURLFromName() (url.URL, error) {
	task := models.TaskType(f.BridgeName)
	bridge, err := f.ORM.FindBridge(task)
	if err != nil {
		return url.URL{}, err
	}
	bridgeURL := url.URL(bridge.URL)
	return bridgeURL, nil
}
