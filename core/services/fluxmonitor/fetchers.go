package fluxmonitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/guregu/null"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shopspring/decimal"
	"go.uber.org/multierr"
)

//go:generate mockery --name Fetcher --output ../../internal/mocks/ --case=underscore

// Fetcher is the interface encapsulating all functionality needed to retrieve
// a price.
type Fetcher interface {
	Fetch() (decimal.Decimal, error)
}

// httpFetcher retrieves data via HTTP from an external price adapter source.
type httpFetcher struct {
	client      *http.Client
	url         *url.URL
	requestData string
}

func newHTTPFetcher(
	timeout models.Duration,
	requestData string,
	url *url.URL,
) Fetcher {
	client := &http.Client{Timeout: timeout.Duration(), Transport: http.DefaultTransport}
	client.Transport = promhttp.InstrumentRoundTripperDuration(promFMResponseTime, client.Transport)
	client.Transport = instrumentRoundTripperReponseSize(promFMResponseSize, client.Transport)

	return &httpFetcher{
		client:      client,
		url:         url,
		requestData: requestData,
	}
}

func (p *httpFetcher) Fetch() (decimal.Decimal, error) {
	request, err := withRandomID(p.requestData)
	if err != nil {
		return decimal.Decimal{}, errors.Wrap(err, fmt.Sprintf("unable to fetch price from %s, cannot add request ID", p.url.String()))
	}
	r, err := p.client.Post(p.url.String(), "application/json", strings.NewReader(request))
	if err != nil {
		return decimal.Decimal{}, errors.Wrap(err, fmt.Sprintf("unable to fetch price from %s with payload '%s'", p.url.String(), p.requestData))
	}

	defer logger.ErrorIfCalling(r.Body.Close)
	target := adapterResponse{}
	if err = json.NewDecoder(r.Body).Decode(&target); err != nil {
		return decimal.Decimal{}, errors.Wrap(err, fmt.Sprintf("unable to decode price from %s", p.url.String()))
	}
	if target.ErrorMessage.Valid {
		return decimal.Decimal{}, errors.Wrap(errors.New(target.ErrorMessage.String), fmt.Sprintf("price fetcher %s returned error", p.url.String()))
	}
	if r.StatusCode >= 400 {
		return decimal.Decimal{}, fmt.Errorf("status code: %d, no error message; unable to retrieve price from %s", r.StatusCode, p.url.String())
	}

	result := target.Result()
	if result == nil {
		return decimal.Decimal{}, errors.Wrap(errors.New("no result returned"), fmt.Sprintf("unable to fetch price from %s", p.url.String()))
	}

	resultFloat, _ := result.Float64()
	promFMIndividualReportedValue.WithLabelValues(p.url.String()).Set(resultFloat)
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

// withRandomID add an arbitrary "id" field to the request json
// this is done in order to keep request payloads consistent in format
// between flux monitor polling requests and http/bridge adapters
func withRandomID(rawReqData string) (string, error) {
	rawReqData = strings.TrimSpace(rawReqData)
	valid := json.Valid([]byte(rawReqData))
	if !valid {
		return "", errors.New(fmt.Sprintf("invalid raw request json: %s", rawReqData))
	}
	return fmt.Sprintf(`{"id":"%s",%s`, models.NewID(), rawReqData[1:]), nil
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
	requestData string,
	priceURLs []*url.URL,
) (Fetcher, error) {
	fetchers := []Fetcher{}
	for _, url := range priceURLs {
		ps := newHTTPFetcher(timeout, requestData, url)
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

func (m *medianFetcher) Fetch() (decimal.Decimal, error) {
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
			price, err := fetcher.Fetch()
			if err != nil {
				logger.Error(err)
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

	errorRate := float64(len(fetchErrors)) / float64(len(m.fetchers))
	if errorRate >= 0.5 {
		return decimal.Decimal{}, errors.Wrap(multierr.Combine(fetchErrors...), "majority of fetchers in median failed")
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
