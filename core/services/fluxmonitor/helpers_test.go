package fluxmonitor

import (
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

	"github.com/smartcontractkit/chainlink/core/services/eth/contracts"
)

func ExportedSetCheckerFactory(fm Service, fac DeviationCheckerFactory) {
	impl := fm.(*concreteFluxMonitor)
	impl.checkerFactory = fac
}

func (p *PollingDeviationChecker) ExportedPollIfEligible(threshold float64) bool {
	return p.pollIfEligible(threshold)
}

func (p *PollingDeviationChecker) ExportedSetStoredReportableRoundID(roundID *big.Int) {
	p.reportableRoundID = roundID
}

func (p *PollingDeviationChecker) ExportedRespondToNewRoundLog(log *contracts.LogNewRound) {
	p.respondToNewRoundLog(*log)
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
