package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/guregu/null"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ethUSDPairing has the ETH/USD parameters needed when POSTing to the price
// external adapters.
// https://github.com/smartcontractkit/price-adapters
const ethUSDPairing = `{"data":{"coin":"ETH","market":"USD"}}`

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
				assert.NoError(t, err)
				urls = append(urls, newURL)
			}

			medianFetcher, err := newMedianFetcherFromURLs(defaultHTTPTimeout, ethUSDPairing, urls)
			require.NoError(t, err)

			medianPrice, err := medianFetcher.Fetch()
			assert.NoError(t, err)
			assert.Equal(t, test.expect, medianPrice.String())
		})
	}
}

func TestNewMedianFetcherFromURLs_Error(t *testing.T) {
	s1 := httptest.NewServer(fakePriceResponder(t, ethUSDPairing, decimal.NewFromInt(101)))
	defer s1.Close()
	var urls []*url.URL

	_, err := newMedianFetcherFromURLs(defaultHTTPTimeout, ethUSDPairing, urls)
	require.Error(t, err)
}

func TestHTTPFetcher_Happy(t *testing.T) {
	btcUSDPairing := `{"data":{"coin":"BTC","market":"USD"}}`
	s1 := httptest.NewServer(fakePriceResponder(t, btcUSDPairing, decimal.NewFromInt(9700)))
	defer s1.Close()
	feedURL, err := url.ParseRequestURI(s1.URL)
	assert.NoError(t, err)

	fetcher := newHTTPFetcher(defaultHTTPTimeout, btcUSDPairing, feedURL)
	price, err := fetcher.Fetch()
	assert.NoError(t, err)
	assert.Equal(t, decimal.NewFromInt(9700), price)
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
	assert.NoError(t, err)

	fetcher := newHTTPFetcher(defaultHTTPTimeout, ethUSDPairing, feedURL)
	price, err := fetcher.Fetch()
	assert.Error(t, err)
	assert.Equal(t, decimal.NewFromInt(0).String(), price.String())
	assert.Contains(t, err.Error(), "could not hit data fetcher")
}

func TestHTTPFetcher_OnlyErrorMessage(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_, err := w.Write([]byte(mustReadFile(t, "testdata/coinmarketcap.error.json")))
		require.NoError(t, err)
	})

	server := httptest.NewServer(handler)
	defer server.Close()
	feedURL, err := url.ParseRequestURI(server.URL)
	assert.NoError(t, err)

	fetcher := newHTTPFetcher(defaultHTTPTimeout, ethUSDPairing, feedURL)
	price, err := fetcher.Fetch()
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
	assert.NoError(t, err)

	fetcher := newHTTPFetcher(defaultHTTPTimeout, ethUSDPairing, feedURL)
	price, err := fetcher.Fetch()
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
		{"bravenewcoin", mustReadFile(t, "testdata/bravenewcoin.json"), decimal.NewFromFloat(306.52036004)},
		{"coinmarketcap", mustReadFile(t, "testdata/coinmarketcap.json"), decimal.NewFromFloat(305.5574615)},
		{"cryptocompare", mustReadFile(t, "testdata/cryptocompare.json"), decimal.NewFromFloat(305.76)},
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
	price, err := medianFetcher.Fetch()
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

			medianPrice, err := medianFetcher.Fetch()
			assert.NoError(t, err)
			assert.True(t, decimal.NewFromInt(100).Equal(medianPrice))
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

			medianPrice, err := medianFetcher.Fetch()
			assert.Error(t, err)
			assert.True(t, decimal.NewFromInt(0).Equal(medianPrice))
		})
	}
}
