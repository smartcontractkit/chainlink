package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
)

type BridgeFetcher struct {
	BaseFetcher

	ID           uint64                 `json:"-" gorm:"primary_key;auto_increment"`
	BridgeName   string                 `json:"name"`
	RequestData  map[string]interface{} `json:"requestData"`
	Transformers Transformers           `json:"transformPipeline,omitempty"`

	ORM    BridgeFetcherORM `json:"-" gorm:"-"`
	Config *orm.Config      `json:"-" gorm:"-"`
}

type BridgeFetcherORM interface {
	FindBridge(name models.TaskType) (models.BridgeType, error)
}

func (f BridgeFetcher) Fetch() (out interface{}, err error) {
	defer func() { f.notifiee.OnEndStage(f, out, err) }()
	f.notifiee.OnBeginStage(f, nil)

	url, err := getBridgeURLFromName(f.BridgeName, f.ORM)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: timeout.Duration(), Transport: http.DefaultTransport}
	client.Transport = promhttp.InstrumentRoundTripperDuration(promFMResponseTime, client.Transport)
	client.Transport = instrumentRoundTripperReponseSize(promFMResponseSize, client.Transport)

	requestDataMap, err := withRandomID(p.RequestData)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to fetch price from %s, cannot add request ID", url.String())
	}

	bodyBytes, err := json.Marshal(requestDataMap)
	if err != nil {
		return nil, errors.Wrapf(err, "could not encode request data to json: %v", err)
	}

	result, err := HttpFetcher{URL: url, Method: "POST", Body: requestDataMap, Config: f.Config}.Fetch()
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

// withRandomID add an arbitrary "id" field to the request json
// this is done in order to keep request payloads consistent in format
// between flux monitor polling requests and http/bridge adapters
func withRandomID(rawReqData map[string]interface{}) (map[string]interface{}, error) {
	rawReqData = strings.TrimSpace(rawReqData)
	valid := json.Valid([]byte(rawReqData))
	if !valid {
		return "", errors.New(fmt.Sprintf("invalid raw request json: %s", rawReqData))
	}
	return fmt.Sprintf(`{"id":"%s",%s`, models.NewID(), rawReqData[1:]), nil
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

func getBridgeURLFromName(name string, orm *orm.ORM) (*url.URL, error) {
	task := models.TaskType(name)
	bridge, err := orm.FindBridge(task)
	if err != nil {
		return nil, err
	}
	bridgeURL := url.URL(bridge.URL)
	return &bridgeURL, nil
}
