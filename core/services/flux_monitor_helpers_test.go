package services

import (
	"chainlink/core/eth"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExportedSetCheckerFactory(fm FluxMonitor, fac DeviationCheckerFactory) {
	impl := fm.(*concreteFluxMonitor)
	impl.checkerFactory = fac
}

func (p *PollingDeviationChecker) ExportedFetchAggregatorData(client eth.Client) error {
	return p.fetchAggregatorData(client)
}

func (p *PollingDeviationChecker) ExportedRespondToNewRound(log eth.Log) error {
	return p.respondToNewRound(log)
}

func (p *PollingDeviationChecker) ExportedPoll() error {
	return p.poll(p.threshold)
}

// ExportedCurrentPrice returns the private current price for assertions;
// technically thread unsafe because it can be set in parallel from
// the CSP consumer.
func (p *PollingDeviationChecker) ExportedCurrentPrice() decimal.Decimal {
	return p.currentPrice
}

// ExportedCurrentRound returns the private current round for assertions;
// technically thread unsafe because it can be set in parallel from
// the CSP consumer.
func (p *PollingDeviationChecker) ExportedCurrentRound() *big.Int {
	return new(big.Int).Set(p.currentRound)
}

func mustReadFile(t testing.TB, file string) string {
	t.Helper()

	content, err := ioutil.ReadFile(file)
	require.NoError(t, err)
	return string(content)
}

type fixedFetcher struct {
	price decimal.Decimal
}

func newFixedPricedFetcher(price decimal.Decimal) *fixedFetcher {
	return &fixedFetcher{price: price}
}

func (ps *fixedFetcher) Fetch() (decimal.Decimal, error) {
	return ps.price, nil
}

type erroringFetcher struct{}

func newErroringPricedFetcher() *erroringFetcher {
	return &erroringFetcher{}
}

func (*erroringFetcher) Fetch() (decimal.Decimal, error) {
	return decimal.NewFromInt(0), errors.New("failed to fetch; I always error")
}

func fakePriceResponder(t *testing.T, requestData string, result decimal.Decimal) http.Handler {
	t.Helper()

	response := adapterResponse{Data: dataWithResult(t, result)}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		defer r.Body.Close()
		assert.Equal(t, requestData, string(payload))
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(response))
	})
}

func dataWithResult(t *testing.T, result decimal.Decimal) adapterResponseData {
	t.Helper()
	var data adapterResponseData
	body := []byte(fmt.Sprintf(`{"result":%v}`, result))
	require.NoError(t, json.Unmarshal(body, &data))
	return data
}
