package fluxmonitor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitor/promfm"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shopspring/decimal"
	"go.uber.org/multierr"
	"gopkg.in/guregu/null.v4"
)

//go:generate mockery --name Fetcher --output ../../internal/mocks/ --case=underscore

// Fetcher is the interface encapsulating all functionality needed to retrieve
// a price.
type Fetcher interface {
	Fetch(context.Context, map[string]interface{}, logger.Logger) (decimal.Decimal, error)
}

// httpFetcher retrieves data via HTTP from an external price adapter source.
type httpFetcher struct {
	client      *http.Client
	url         *url.URL
	requestData map[string]interface{}
	sizeLimit   int64
}

func newHTTPFetcher(
	timeout models.Duration,
	requestData map[string]interface{},
	url *url.URL,
	sizeLimit int64,
) Fetcher {
	client := &http.Client{Timeout: timeout.Duration(), Transport: http.DefaultTransport}
	client.Transport = promhttp.InstrumentRoundTripperDuration(promfm.ResponseTime, client.Transport)
	client.Transport = promfm.InstrumentRoundTripperReponseSize(promfm.ResponseSize, client.Transport)

	return &httpFetcher{
		client:      client,
		url:         url,
		requestData: requestData,
		sizeLimit:   sizeLimit,
	}
}

func (p *httpFetcher) Fetch(ctx context.Context, meta map[string]interface{}, logger logger.Logger) (decimal.Decimal, error) {
	request := withIDAndMeta(p.requestData, meta)
	body, err := json.Marshal(request)
	if err != nil {
		return decimal.Decimal{}, errors.Wrap(err, "error encoding request body as JSON")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.url.String(), bytes.NewReader(body))
	if err != nil {
		return decimal.Decimal{}, errors.Wrap(err, "unable to create request")
	}

	req.Header.Add("Content-Type", "application/json")
	response, err := p.client.Do(req)
	if err != nil {
		return decimal.Decimal{}, errors.Wrap(err, fmt.Sprintf("unable to fetch price from %s with payload '%s'", p.url.String(), p.requestData))
	}

	defer logger.ErrorIfCalling(response.Body.Close)
	target := adapterResponse{}
	responseReader := utils.NewMaxBytesReader(response.Body, p.sizeLimit)
	if err = json.NewDecoder(responseReader).Decode(&target); err != nil {
		return decimal.Decimal{}, errors.Wrap(err, fmt.Sprintf("unable to decode price from %s", p.url.String()))
	}
	if target.ErrorMessage.Valid {
		return decimal.Decimal{}, errors.Wrap(errors.New(target.ErrorMessage.String), fmt.Sprintf("price fetcher %s returned error", p.url.String()))
	}
	if response.StatusCode >= 400 {
		return decimal.Decimal{}, fmt.Errorf("status code: %d, no error message; unable to retrieve price from %s", response.StatusCode, p.url.String())
	}

	result := target.Result()
	if result == nil {
		return decimal.Decimal{}, errors.Wrap(errors.New("no result returned"), fmt.Sprintf("unable to fetch price from %s", p.url.String()))
	}

	resultFloat, _ := result.Float64()
	promfm.IndividualReportedValue.WithLabelValues(p.url.String()).Set(resultFloat)
	logger.Debugw(
		fmt.Sprintf("fetched price %v from %s", *result, p.url.String()),
		"price", result,
		"url", p.url.String(),
	)
	return *result, nil
}

func (p *httpFetcher) String() string {
	return fmt.Sprintf("http price fetcher: %s", p.url.String())
}

func withIDAndMeta(request, meta map[string]interface{}) map[string]interface{} {
	output := make(map[string]interface{})
	for k, v := range request {
		output[k] = v
	}
	output["id"] = uuid.NewV4()
	output["meta"] = meta
	return output
}

type adapterResponseData struct {
	Result *decimal.Decimal `json:"result"`
}

// adapterResponse is the HTTP response as defined by the external adapter:
// https://github.com/smartcontractkit/bnc-adapter
type adapterResponse struct {
	Data         adapterResponseData `json:"data"`
	ErrorMessage null.String         `json:"errorMessage"`
}

func (pr adapterResponse) Result() *decimal.Decimal {
	return pr.Data.Result
}

// medianFetcher fetches from all fetchers, and returns the median value, or
// average if even number of results.
type medianFetcher struct {
	fetchers []Fetcher
}

// newMedianFetcherFromURLs creates a median fetcher that retrieves a price
// from all passed URLs using httpFetcher, and returns the median.
func newMedianFetcherFromURLs(
	timeout models.Duration,
	requestData map[string]interface{},
	priceURLs []*url.URL,
	sizeLimit int64,
) (Fetcher, error) {
	fetchers := []Fetcher{}
	for _, url := range priceURLs {
		ps := newHTTPFetcher(timeout, requestData, url, sizeLimit)
		fetchers = append(fetchers, ps)
	}

	medianFetcher, err := newMedianFetcher(fetchers...)
	if err != nil {
		return nil, err
	}

	return medianFetcher, nil
}

func newMedianFetcher(fetchers ...Fetcher) (Fetcher, error) {
	if len(fetchers) == 0 {
		return nil, errors.New("must pass in at least one price fetcher to newMedianFetcher")
	}
	return &medianFetcher{
		fetchers: fetchers,
	}, nil
}

func (m *medianFetcher) Fetch(ctx context.Context, meta map[string]interface{}, logger logger.Logger) (decimal.Decimal, error) {
	prices := []decimal.Decimal{}
	fetchErrors := []error{}

	type result struct {
		price decimal.Decimal
		err   error
	}

	chResults := make(chan result)
	for _, fetcher := range m.fetchers {
		fetcher := fetcher
		go func() {
			price, err := fetcher.Fetch(ctx, meta, logger)
			if err != nil {
				logger.Warn(err)
				chResults <- result{err: err}
			} else {
				chResults <- result{price: price}
			}
		}()
	}

	for i := 0; i < len(m.fetchers); i++ {
		r := <-chResults
		if r.err != nil {
			fetchErrors = append(fetchErrors, r.err)
		} else {
			prices = append(prices, r.price)
		}
	}

	fetchersCount := len(m.fetchers)
	fetchErrorsCount := len(fetchErrors)
	errorRate := float64(fetchErrorsCount) / float64(fetchersCount)
	if errorRate >= 0.5 {
		err := errors.Wrap(multierr.Combine(fetchErrors...), fmt.Sprintf("at least 50%% of the fetchers in median failed (%d/%d)", fetchErrorsCount, fetchersCount))
		return decimal.Decimal{}, err
	}

	sort.Slice(prices, func(i, j int) bool {
		return prices[i].LessThan(prices[j])
	})
	k := len(prices) / 2
	if len(prices)%2 == 1 {
		return prices[k], nil
	}
	return prices[k].Add(prices[k-1]).Div(decimal.NewFromInt(2)), nil
}

func (m *medianFetcher) String() string {
	fetcherDescriptions := make([]string, len(m.fetchers))
	for i, fetcher := range m.fetchers {
		fetcherDescriptions[i] = fmt.Sprintf("%s", fetcher)
	}
	return fmt.Sprintf("median fetcher: %s", strings.Join(fetcherDescriptions, ","))
}
