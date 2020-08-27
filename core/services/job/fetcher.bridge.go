package job

import (
	"encoding/json"
	"fmt"
	"net/url"

	// "github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type BridgeFetcher struct {
	BaseFetcher

	BridgeName   string          `json:"name"`
	RequestData  HttpRequestData `json:"requestData" gorm:"type:jsonb"`
	Transformers Transformers    `json:"transformPipeline,omitempty" gorm:"-"`

	ORM                BridgeFetcherORM `json:"-" gorm:"-"`
	defaultHTTPTimeout models.Duration
}

type BridgeFetcherORM interface {
	FindBridge(name models.TaskType) (models.BridgeType, error)
}

func (f *BridgeFetcher) Fetch() (out interface{}, err error) {
	defer func() { f.notifiee.OnEndStage(f, out, err) }()
	f.notifiee.OnBeginStage(f, nil)

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

	result, err := (&HttpFetcher{URL: models.WebURL(url), Method: "POST", RequestData: f.RequestData}).Fetch()
	if err != nil {
		return nil, err
	}

	logger.Debugw(
		fmt.Sprintf("Fetched answer", result, url.String()),
		"answer", result,
		"url", url.String(),
	)
	return f.Transformers.Transform(result)
}

func (f *BridgeFetcher) SetNotifiee(n Notifiee) {
	f.notifiee = n
	f.Transformers.SetNotifiee(n)
}

func (f BridgeFetcher) MarshalJSON() ([]byte, error) {
	type preventInfiniteRecursion BridgeFetcher
	type fetcherWithType struct {
		Type FetcherType `json:"type"`
		preventInfiniteRecursion
	}
	f2 := fetcherWithType{FetcherTypeBridge, preventInfiniteRecursion(f)}
	return json.Marshal(f2)
}

func (f BridgeFetcher) getBridgeURLFromName() (url.URL, error) {
	task := models.TaskType(f.BridgeName)
	bridge, err := f.ORM.FindBridge(task)
	if err != nil {
		return url.URL{}, err
	}
	bridgeURL := url.URL(bridge.URL)
	return bridgeURL, nil
}
