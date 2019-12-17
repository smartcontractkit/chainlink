package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/guregu/null"
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
		prices []float64
		expect float64
	}{
		{"single", []float64{101}, 101},
		{"odd", []float64{101, 102, 103}, 102},
		{"even", []float64{101, 102, 103, 104}, 102.5},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			urls := []string{}
			for _, price := range test.prices {
				s := httptest.NewServer(fakePriceResponder(t, ethUSDPairing, price))
				defer s.Close()
				urls = append(urls, s.URL)
			}

			medianFetcher, err := newMedianFetcherFromURLs(defaultHTTPTimeout, ethUSDPairing, urls...)
			require.NoError(t, err)

			medianPrice, err := medianFetcher.Fetch()
			assert.NoError(t, err)
			assert.Equal(t, test.expect, medianPrice)
		})
	}
}

func TestNewMedianFetcherFromURLs_Error(t *testing.T) {
	s1 := httptest.NewServer(fakePriceResponder(t, ethUSDPairing, 101))
	defer s1.Close()

	_, err := newMedianFetcherFromURLs(defaultHTTPTimeout, ethUSDPairing, s1.URL, "garbage")
	require.Error(t, err)
}

func TestHTTPFetcher_Happy(t *testing.T) {
	btcUSDPairing := `{"data":{"coin":"BTC","market":"USD"}}`
	s1 := httptest.NewServer(fakePriceResponder(t, btcUSDPairing, 9700))
	defer s1.Close()

	fetcher, err := newHTTPFetcher(defaultHTTPTimeout, btcUSDPairing, s1.URL)
	require.NoError(t, err)
	price, err := fetcher.Fetch()
	assert.NoError(t, err)
	assert.Equal(t, float64(9700), price)
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

	fetcher, err := newHTTPFetcher(defaultHTTPTimeout, ethUSDPairing, server.URL)
	require.NoError(t, err)
	price, err := fetcher.Fetch()
	assert.Error(t, err)
	assert.Equal(t, float64(0), price)
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

	fetcher, err := newHTTPFetcher(defaultHTTPTimeout, ethUSDPairing, server.URL)
	require.NoError(t, err)
	price, err := fetcher.Fetch()
	assert.Error(t, err)
	assert.Equal(t, float64(0), price)
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

	fetcher, err := newHTTPFetcher(defaultHTTPTimeout, ethUSDPairing, server.URL)
	require.NoError(t, err)
	price, err := fetcher.Fetch()
	assert.Error(t, err)
	assert.Equal(t, float64(0), price)
}

// Sample input taken from
// https://github.com/smartcontractkit/price-adapters#chainlink-price-request-adapters
func TestAdapterResponse_UnmarshalJSON_Happy(t *testing.T) {
	tests := []struct {
		name, content string
		expect        float64
	}{
		{"basic", `{"data":{"result":123.4567890},"jobRunID":"1","statusCode":200}`, 123.456789},
		{"bravenewcoin", mustReadFile(t, "testdata/bravenewcoin.json"), 306.52036004},
		{"coinmarketcap", mustReadFile(t, "testdata/coinmarketcap.json"), 305.5574615},
		{"cryptocompare", mustReadFile(t, "testdata/cryptocompare.json"), 305.76},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var response adapterResponse
			err := json.Unmarshal([]byte(test.content), &response)
			assert.NoError(t, err)
			result := response.Result()
			assert.Equal(t, test.expect, *result)
		})
	}
}

func TestAdapterResponse_Result_float(t *testing.T) {
	var pr adapterResponse
	input := `{"data":{"result":100.1}}`
	assert.NoError(t, json.Unmarshal([]byte(input), &pr))

	result := pr.Result()
	assert.Equal(t, 100.1, *result)
}

func TestAdapterResponse_Result_empty(t *testing.T) {
	var pr adapterResponse
	input := `{"data":{"other":"100.1"}}`
	assert.NoError(t, json.Unmarshal([]byte(input), &pr))

	assert.Nil(t, pr.Result())
}

func TestNewMedianFetcher_EmptyFetchersError(t *testing.T) {
	_, err := newMedianFetcher()
	require.Error(t, err)
}

func TestMedianFetcher_FetchError(t *testing.T) {
	s1 := newFixedPricedFetcher(102)
	s2 := newErroringPricedFetcher()
	medianFetcher, err := newMedianFetcher(s1, s2)
	require.NoError(t, err)
	price, err := medianFetcher.Fetch()
	assert.Error(t, err)
	assert.Equal(t, float64(0), price)
}

func TestMedianFetcher_MajorityFetches(t *testing.T) {
	hf := newFixedPricedFetcher(100) // healthy fetcher
	ef := newErroringPricedFetcher() // erroring fetcher

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
			assert.Equal(t, float64(100), medianPrice)
		})
	}
}

func TestMedianFetcher_MinorityErrors(t *testing.T) {
	hf := newFixedPricedFetcher(100) // healthy fetcher
	ef := newErroringPricedFetcher() // erroring fetcher

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
			assert.Equal(t, float64(0), medianPrice)
		})
	}
}
