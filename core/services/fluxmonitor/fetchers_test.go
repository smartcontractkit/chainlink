package fluxmonitor

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// ethUSDPairing has the ETH/USD parameters needed when POSTing to the price
// external adapters.
// https://github.com/smartcontractkit/price-adapters
var (
	ethUSDPairing      = utils.MustUnmarshalToMap(`{"data":{"coin":"ETH","market":"USD"}}`)
	defaultHTTPTimeout = models.MustMakeDuration(15 * time.Second)
	emptyMeta          = utils.MustUnmarshalToMap("{}")
)

func TestNewMedianFetcherFromURLs_Happy(t *testing.T) {
	tests := []struct {
		name   string
		prices []decimal.Decimal
		expect string
	}{
		{
			"single",
			[]decimal.Decimal{
				decimal.NewFromInt(101),
			},
			"101",
		},
		{
			"odd",
			[]decimal.Decimal{
				decimal.NewFromInt(101),
				decimal.NewFromInt(102),
				decimal.NewFromInt(103),
			},
			"102",
		},
		{
			"even",
			[]decimal.Decimal{
				decimal.NewFromInt(101),
				decimal.NewFromInt(102),
				decimal.NewFromInt(103),
				decimal.NewFromInt(104),
			},
			"102.5",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var urls []*url.URL
			for _, price := range test.prices {
				s := httptest.NewServer(fakePriceResponder(t, ethUSDPairing, price))
				defer s.Close()
				newURL, err := url.ParseRequestURI(s.URL)
				require.NoError(t, err)
				urls = append(urls, newURL)
			}

			medianFetcher, err := newMedianFetcherFromURLs(defaultHTTPTimeout, ethUSDPairing, urls, 32768)
			require.NoError(t, err)

			medianPrice, err := medianFetcher.Fetch(context.Background(), emptyMeta, *logger.Default)
			require.NoError(t, err)
			assert.Equal(t, test.expect, medianPrice.String())
		})
	}
}

func TestNewMedianFetcherFromURLs_EmptyError(t *testing.T) {
	s1 := httptest.NewServer(fakePriceResponder(t, ethUSDPairing, decimal.NewFromInt(101)))
	defer s1.Close()
	var urls []*url.URL

	_, err := newMedianFetcherFromURLs(defaultHTTPTimeout, ethUSDPairing, urls, 32768)
	require.Error(t, err)
}

func TestHTTPFetcher_Happy(t *testing.T) {
	btcUSDPairing := utils.MustUnmarshalToMap(`{"data":{"coin":"BTC","market":"USD"}}`)
	s1 := httptest.NewServer(fakePriceResponder(t, btcUSDPairing, decimal.NewFromInt(9700)))
	defer s1.Close()
	feedURL, err := url.ParseRequestURI(s1.URL)
	require.NoError(t, err)

	fetcher := newHTTPFetcher(defaultHTTPTimeout, btcUSDPairing, feedURL, 32768)
	price, err := fetcher.Fetch(context.Background(), emptyMeta, *logger.Default)
	require.NoError(t, err)
	assert.Equal(t, decimal.NewFromInt(9700), price)
}

func TestHTTPFetcher_Meta(t *testing.T) {
	empty := adapterResponse{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req fetcherRequest
		body, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(body, &req)
		require.NoError(t, err)
		require.Equal(t, float64(10), req.Meta["latestAnswer"])
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(empty))
	})

	md, err := models.MarshalBridgeMetaData(big.NewInt(10), big.NewInt(1616447984))
	require.NoError(t, err)

	s1 := httptest.NewServer(handler)

	defer s1.Close()
	feedURL, err := url.ParseRequestURI(s1.URL)
	require.NoError(t, err)

	fetcher := newHTTPFetcher(defaultHTTPTimeout, ethUSDPairing, feedURL, 32768)
	fetcher.Fetch(context.Background(), md, *logger.Default)
}

func TestHTTPFetcher_ErrorMessage(t *testing.T) {
	data := adapterResponse{
		ErrorMessage: null.StringFrom("could not hit data fetcher"),
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		require.NoError(t, json.NewEncoder(w).Encode(data))
	})

	server := httptest.NewServer(handler)
	defer server.Close()
	feedURL, err := url.ParseRequestURI(server.URL)
	require.NoError(t, err)

	fetcher := newHTTPFetcher(defaultHTTPTimeout, ethUSDPairing, feedURL, 32768)
	price, err := fetcher.Fetch(context.Background(), emptyMeta, *logger.Default)
	assert.Error(t, err)
	assert.Equal(t, decimal.NewFromInt(0).String(), price.String())
	assert.Contains(t, err.Error(), "could not hit data fetcher")
}

func TestHTTPFetcher_OnlyErrorMessage(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_, err := w.Write([]byte(mustReadFile(t, "../../testdata/apiresponses/coinmarketcap.error.json")))
		require.NoError(t, err)
	})

	server := httptest.NewServer(handler)
	defer server.Close()
	feedURL, err := url.ParseRequestURI(server.URL)
	require.NoError(t, err)

	fetcher := newHTTPFetcher(defaultHTTPTimeout, ethUSDPairing, feedURL, 32768)
	price, err := fetcher.Fetch(context.Background(), emptyMeta, *logger.Default)
	assert.Error(t, err)
	assert.Equal(t, decimal.NewFromInt(0).String(), price.String())
	assert.Contains(t, err.Error(), "RequestId")
}

func TestHTTPFetcher_NoResultNorErrorMessage(t *testing.T) {
	empty := adapterResponse{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(empty))
	})

	server := httptest.NewServer(handler)
	defer server.Close()
	feedURL, err := url.ParseRequestURI(server.URL)
	require.NoError(t, err)

	fetcher := newHTTPFetcher(defaultHTTPTimeout, ethUSDPairing, feedURL, 32768)
	price, err := fetcher.Fetch(context.Background(), emptyMeta, *logger.Default)
	assert.Error(t, err)
	assert.True(t, decimal.NewFromInt(0).Equal(price))
}

// Sample input taken from
// https://github.com/smartcontractkit/price-adapters#chainlink-price-request-adapters
func TestAdapterResponse_UnmarshalJSON_Happy(t *testing.T) {
	tests := []struct {
		name, content string
		expect        decimal.Decimal
	}{
		{"basic", `{"data":{"result":123.4567890},"jobRunID":"1","statusCode":200}`, decimal.NewFromFloat(123.456789)},
		{"bravenewcoin", mustReadFile(t, "../../testdata/apiresponses/bravenewcoin.json"), decimal.NewFromFloat(306.52036004)},
		{"coinmarketcap", mustReadFile(t, "../../testdata/apiresponses/coinmarketcap.json"), decimal.NewFromFloat(305.5574615)},
		{"cryptocompare", mustReadFile(t, "../../testdata/apiresponses/cryptocompare.json"), decimal.NewFromFloat(305.76)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var response adapterResponse
			err := json.Unmarshal([]byte(test.content), &response)
			assert.NoError(t, err)
			result := response.Result()
			assert.Equal(t, test.expect.String(), result.String())
		})
	}
}

func TestNewMedianFetcher_EmptyFetchersError(t *testing.T) {
	_, err := newMedianFetcher()
	require.Error(t, err)
}

func TestMedianFetcher_FetchError(t *testing.T) {
	s1 := newFixedPricedFetcher(decimal.NewFromInt(102))
	s2 := newErroringPricedFetcher()
	medianFetcher, err := newMedianFetcher(s1, s2)
	require.NoError(t, err)
	price, err := medianFetcher.Fetch(context.Background(), emptyMeta, *logger.Default)
	assert.Error(t, err)
	assert.Equal(t, decimal.NewFromInt(0).String(), price.String())
}

func TestMedianFetcher_MajorityFetches(t *testing.T) {
	hf := newFixedPricedFetcher(decimal.NewFromInt(100)) // healthy fetcher)
	ef := newErroringPricedFetcher()                     // erroring fetcher

	tests := []struct {
		name     string
		fetchers []Fetcher
	}{
		{"2/3", []Fetcher{hf, hf, ef}},
		{"3/3", []Fetcher{hf, hf, hf}},
		{"3/4", []Fetcher{hf, hf, hf, ef}},
		{"3/5", []Fetcher{hf, hf, hf, ef, ef}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			medianFetcher, err := newMedianFetcher(test.fetchers...)
			require.NoError(t, err)

			medianPrice, err := medianFetcher.Fetch(context.Background(), emptyMeta, *logger.Default)
			assert.NoError(t, err)
			assert.True(t, decimal.NewFromInt(100).Equal(medianPrice))
		})
	}
}

func TestMedianFetcher_MajorityFetchesCalculatesCorrectMedian(t *testing.T) {
	hf50 := newFixedPricedFetcher(decimal.NewFromInt(50))
	hf75 := newFixedPricedFetcher(decimal.NewFromInt(75))
	hf100 := newFixedPricedFetcher(decimal.NewFromInt(100))
	hf999 := newFixedPricedFetcher(decimal.NewFromInt(999))
	ef := newErroringPricedFetcher()

	tests := []struct {
		name           string
		fetchers       []Fetcher
		expectedMedian string
	}{
		{"3/5", []Fetcher{hf50, hf75, hf100, ef, ef}, "75"},
		{"4/7", []Fetcher{hf50, hf75, hf100, hf999, ef, ef, ef}, "87.5"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			medianFetcher, err := newMedianFetcher(test.fetchers...)
			require.NoError(t, err)

			medianPrice, err := medianFetcher.Fetch(context.Background(), emptyMeta, *logger.Default)
			assert.NoError(t, err)
			assert.Equal(t, medianPrice.String(), test.expectedMedian)
		})
	}
}

func TestMedianFetcher_MinorityErrors(t *testing.T) {
	hf := newFixedPricedFetcher(decimal.NewFromInt(100)) // healthy fetcher
	ef := newErroringPricedFetcher()                     // erroring fetcher

	tests := []struct {
		name     string
		fetchers []Fetcher
	}{
		{"1/2", []Fetcher{hf, ef}},
		{"1/3", []Fetcher{hf, ef, ef}},
		{"2/4", []Fetcher{hf, hf, ef, ef}},
		{"2/5", []Fetcher{hf, hf, ef, ef, ef}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			medianFetcher, err := newMedianFetcher(test.fetchers...)
			require.NoError(t, err)

			medianPrice, err := medianFetcher.Fetch(context.Background(), emptyMeta, *logger.Default)
			assert.Error(t, err)
			assert.True(t, decimal.NewFromInt(0).Equal(medianPrice))
		})
	}
}

func TestHTTPFetcher_AddsArbitraryRequestID(t *testing.T) {
	empty := adapterResponse{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req fetcherRequest
		body, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(body, &req)
		require.NoError(t, err)
		require.NotEmpty(t, req.ID)
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(empty))
	})

	s1 := httptest.NewServer(handler)

	defer s1.Close()
	feedURL, err := url.ParseRequestURI(s1.URL)
	require.NoError(t, err)

	fetcher := newHTTPFetcher(defaultHTTPTimeout, ethUSDPairing, feedURL, 32768)
	fetcher.Fetch(context.Background(), emptyMeta, *logger.Default)
}
