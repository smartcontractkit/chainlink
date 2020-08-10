package fluxmonitor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/eth/contracts"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func ExportedSetCheckerFactory(fm Service, fac DeviationCheckerFactory) {
	impl := fm.(*concreteFluxMonitor)
	impl.checkerFactory = fac
}

func (p *PollingDeviationChecker) ExportedPollIfEligible(threshold, absoluteThreshold float64) {
	p.pollIfEligible(DeviationThresholds{Rel: threshold, Abs: absoluteThreshold})
}

func (p *PollingDeviationChecker) ExportedRespondToNewRoundLog(log *contracts.LogNewRound) {
	p.respondToNewRoundLog(*log)
}

func (p *PollingDeviationChecker) ExportedSufficientFunds(state contracts.FluxAggregatorRoundState) bool {
	return p.sufficientFunds(state)
}

func (p *PollingDeviationChecker) ExportedSufficientPayment(payment *big.Int) bool {
	return p.sufficientPayment(payment)
}

func (p *PollingDeviationChecker) ExportedProcessLogs() {
	p.processLogs()
}

func (p *PollingDeviationChecker) ExportedBacklog() *utils.BoundedPriorityQueue {
	return p.backlog
}

func (p *PollingDeviationChecker) ExportedFluxAggregator() contracts.FluxAggregator {
	return p.fluxAggregator
}

func (p *PollingDeviationChecker) ExportedRoundState() {
	p.roundState(0)
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

type fetcherRequest struct {
	Data interface{} `json:"data"`
	ID   string      `json:"id"`
}

func fakePriceResponder(t *testing.T, requestData string, result decimal.Decimal) http.Handler {
	t.Helper()

	var expectedRequest fetcherRequest
	err := json.Unmarshal([]byte(requestData), &expectedRequest)
	require.NoError(t, err)
	response := adapterResponse{Data: dataWithResult(t, result)}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody fetcherRequest
		payload, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		defer r.Body.Close()
		err = json.Unmarshal(payload, &reqBody)
		require.NoError(t, err)
		assert.Equal(t, expectedRequest.Data, reqBody.Data)
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

// CreateJob is used in TestFluxMonitorAntiSpamLogic to create a
// job with a specific answer and round, for testing nodes with malicious
// behavior
func (fm *concreteFluxMonitor) CreateJob(t *testing.T,
	jobSpecId *models.ID, polledAnswer decimal.Decimal,
	nextRound *big.Int) error {
	jobSpec, err := fm.store.ORM.FindJob(jobSpecId)
	require.NoError(t, err, "could not find job spec with that ID")
	checker, err := fm.checkerFactory.New(jobSpec.Initiators[0], nil, fm.runManager,
		fm.store.ORM, models.MustMakeDuration(100*time.Second))
	require.NoError(t, err, "could not create deviation checker")
	payment := fm.store.Config.MinimumContractPayment()
	return checker.(*PollingDeviationChecker).createJobRun(polledAnswer, uint32(nextRound.Uint64()), payment)
}
