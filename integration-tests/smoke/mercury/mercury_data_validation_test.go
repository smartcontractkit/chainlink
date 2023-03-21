package smoke

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"sync"
	"testing"
	"text/tabwriter"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	mercuryactions "github.com/smartcontractkit/chainlink/integration-tests/actions/mercury"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups/mercury"
)

type expectedResult struct {
	FeedId [32]byte
	Value  *big.Int
}

type feedIdResult struct {
	Ok       bool
	Expected [32]byte
	Actual   [32]byte
}

type numberResult struct {
	Ok       bool
	Expected *big.Int
	Actual   *big.Int
}

type blockNumberResult struct {
	Ok                      bool
	Expected                uint64
	ActualValidFromBlockNum uint64
	ActualCurrentBlockNum   uint64
}

type observationTimestampResult struct {
	Ok                          bool
	BlockNumber                 uint64
	ActualObservationsTimestamp uint32
}

type RoundResult struct {
	Duration   time.Duration
	Validators map[string]ValidatorResults
}

type ValidatorResults struct {
	AllResults []testResult
	//MissingReports int
	// CalculateResults()
}

func (vr *ValidatorResults) RequestsCount() int {
	return len(vr.AllResults)
}

func (vr *ValidatorResults) MissingReports() []testResult {
	results := []testResult{}
	for _, r := range vr.AllResults {
		if r.MissingReport {
			results = append(results, r)
		}
	}
	return results
}

func (vr *ValidatorResults) WrongFeedId() []testResult {
	results := []testResult{}
	for _, r := range vr.AllResults {
		if !r.MissingReport && !r.FeedId.Ok {
			results = append(results, r)
		}
	}
	return results
}

func (vr *ValidatorResults) WrongBid() []testResult {
	results := []testResult{}
	for _, r := range vr.AllResults {
		if !r.MissingReport && !r.Bid.Ok {
			results = append(results, r)
		}
	}
	return results
}

func (vr *ValidatorResults) WrongAsk() []testResult {
	results := []testResult{}
	for _, r := range vr.AllResults {
		if !r.MissingReport && !r.Ask.Ok {
			results = append(results, r)
		}
	}
	return results
}

func (vr *ValidatorResults) WrongBenchmark() []testResult {
	results := []testResult{}
	for _, r := range vr.AllResults {
		if !r.MissingReport && !r.BenchmarkPrice.Ok {
			results = append(results, r)
		}
	}
	return results
}

func (vr *ValidatorResults) WrongBlockNumbers() []testResult {
	results := []testResult{}
	for _, r := range vr.AllResults {
		if !r.MissingReport && !r.BlockNumber.Ok {
			results = append(results, r)
		}
	}
	return results
}

type testResult struct {
	Id                    string
	MissingReport         bool
	FeedId                feedIdResult
	BenchmarkPrice        numberResult
	Bid                   numberResult
	Ask                   numberResult
	BlockNumber           blockNumberResult
	ObservationsTimestamp observationTimestampResult
	Err                   error
}

func validateNewReportsEveryBlock(
	validatorId string,
	duration time.Duration,
	callTimeout time.Duration,
	feedId string, er expectedResult,
	evmClient blockchain.EVMClient, msClient *client.MercuryServer,
	wg *sync.WaitGroup, resultChan chan testResult) {

	defer wg.Done()

	headerChan := make(chan *blockchain.SafeEVMHeader)
	sub, err := evmClient.SubscribeNewHeaders(context.Background(), headerChan)
	if err != nil {
		resultChan <- testResult{Err: err}
		return
	}
	defer sub.Unsubscribe()

	timeout := time.After(duration)

	for {
		select {
		case err := <-sub.Err():
			sub.Unsubscribe()
			resultChan <- testResult{Id: validatorId, Err: err}
			return

		case <-timeout:
			log.Info().Msgf("Validator %s is done!", validatorId)
			return

		case header := <-headerChan:
			// Ask mercury for report every time new block is generated
			bn := header.Number.Uint64()

			// Wait for some time between new header and first mercury api call
			time.Sleep(callTimeout)

			requestTime := time.Now()
			reportData, _, err := msClient.GetReports(feedId, bn)
			if err != nil {
				resultChan <- testResult{Id: validatorId, Err: err}
				break
			}

			if reportData.ChainlinkBlob == "" {
				resultChan <- testResult{Id: validatorId, MissingReport: true}
				log.Warn().Msgf("No report for block: %d for %s", bn, validatorId)
			} else {
				reportBytes, err := hex.DecodeString(reportData.ChainlinkBlob[2:])
				if err != nil {
					resultChan <- testResult{Err: err}
				}
				reportCtx, err := mercuryactions.DecodeReport(reportBytes)
				if err != nil {
					resultChan <- testResult{Err: err}
				}
				report := reportCtx.Report
				result := testResult{
					Id: validatorId,
				}
				// Compare feed id
				result.FeedId = feedIdResult{
					Ok:       bytes.Equal(er.FeedId[:], report.FeedId[:]),
					Expected: er.FeedId,
					Actual:   report.FeedId,
				}
				// Compare bid
				result.Bid = numberResult{
					Ok:       er.Value.Cmp(report.Bid) == 0,
					Expected: er.Value,
					Actual:   report.Bid,
				}
				// Compare ask
				result.Ask = numberResult{
					Ok:       er.Value.Cmp(report.Ask) == 0,
					Expected: er.Value,
					Actual:   report.Ask,
				}
				// Compare benchmark price
				result.BenchmarkPrice = numberResult{
					Ok:       er.Value.Cmp(report.BenchmarkPrice) == 0,
					Expected: er.Value,
					Actual:   report.BenchmarkPrice,
				}
				// Compare block number
				result.BlockNumber = blockNumberResult{
					Ok:                      bn == report.CurrentBlockNum || bn == report.ValidFromBlockNum,
					Expected:                bn,
					ActualValidFromBlockNum: report.ValidFromBlockNum,
					ActualCurrentBlockNum:   report.CurrentBlockNum,
				}
				// TODO: Compare with observation timestamp for report.ValidFromBlockNum timestamp?
				ot := time.Unix(int64(report.ObservationsTimestamp), 0)
				result.ObservationsTimestamp = observationTimestampResult{
					Ok:                          requestTime.After(ot),
					BlockNumber:                 bn,
					ActualObservationsTimestamp: report.ObservationsTimestamp,
				}

				resultChan <- result
			}
		}
	}
}

type SimulationRound struct {
	Duration                      time.Duration
	TimeBetweenNewBlockAndApiCall time.Duration
	DataProviderValue             *big.Int
	ValidatorsCount               int
	Results                       RoundResult
	// Percentage of missing reports allowed in the round
	MissingReportsThreshold float32
}

type Simulation struct {
	Rounds []SimulationRound
}

// Define simulation params for the test
var simulation = Simulation{
	Rounds: []SimulationRound{
		{
			Duration:                      20 * time.Second,
			TimeBetweenNewBlockAndApiCall: 500 * time.Millisecond,
			DataProviderValue:             big.NewInt(800),
			ValidatorsCount:               10,
			MissingReportsThreshold:       0.7,
		},
		{
			Duration:                      25 * time.Second,
			TimeBetweenNewBlockAndApiCall: 700 * time.Millisecond,
			DataProviderValue:             big.NewInt(800),
			ValidatorsCount:               10,
			MissingReportsThreshold:       0.5,
		},
		{
			Duration:                      30 * time.Second,
			TimeBetweenNewBlockAndApiCall: 1000 * time.Millisecond,
			DataProviderValue:             big.NewInt(800),
			ValidatorsCount:               10,
			MissingReportsThreshold:       0.2,
		},
		{
			Duration:                      40 * time.Second,
			TimeBetweenNewBlockAndApiCall: 1500 * time.Millisecond,
			DataProviderValue:             big.NewInt(800),
			ValidatorsCount:               10,
			MissingReportsThreshold:       0.0,
		},
		{
			Duration:                      50 * time.Second,
			TimeBetweenNewBlockAndApiCall: 2000 * time.Millisecond,
			DataProviderValue:             big.NewInt(800),
			ValidatorsCount:               10,
			MissingReportsThreshold:       0.0,
		},
		{
			Duration:                      60 * time.Second,
			TimeBetweenNewBlockAndApiCall: 2500 * time.Millisecond,
			DataProviderValue:             big.NewInt(800),
			ValidatorsCount:               10,
			MissingReportsThreshold:       0.0,
		},
	},
}

func TestMercuryReportsHaveValidValues(t *testing.T) {
	l := utils.GetTestLogger(t)

	var (
		feedIds = [][32]byte{
			mercury.StringToByte32("feed-1"),
			// mercury.StringToByte32("feed-2"),
		}
	)

	testEnv, err := mercury.NewEnv(t.Name(), "smoke", mercury.DefaultResources)

	t.Cleanup(func() {
		testEnv.Cleanup(t)
	})
	require.NoError(t, err)

	testEnv.AddEvmNetwork()

	err = testEnv.AddDON()
	require.NoError(t, err)

	ocrConfig, err := testEnv.BuildOCRConfig()
	require.NoError(t, err)

	_, _, err = testEnv.AddMercuryServer(nil)
	require.NoError(t, err)

	verifierProxyContract, err := testEnv.AddVerifierProxyContract("verifierProxy1")
	require.NoError(t, err)
	verifierContract, err := testEnv.AddVerifierContract("verifier1", verifierProxyContract.Address())
	require.NoError(t, err)

	for i, feedId := range feedIds {
		blockNumber, err := testEnv.SetConfigAndInitializeVerifierContract(
			fmt.Sprintf("setAndInitializeVerifier%d", i),
			"verifier1",
			"verifierProxy1",
			feedId,
			*ocrConfig,
		)
		require.NoError(t, err)

		err = testEnv.AddBootstrapJob(fmt.Sprintf("createBoostrap%d", i), verifierContract.Address(), uint64(blockNumber), feedId)
		require.NoError(t, err)

		err = testEnv.AddOCRJobs(fmt.Sprintf("createOcrJobs%d", i), verifierContract.Address(), uint64(blockNumber), feedId)
		require.NoError(t, err)
	}

	err = testEnv.WaitForReportsInMercuryDb(feedIds)
	require.NoError(t, err)

	for _, feedIdBytes := range feedIds {
		feedIdStr := mercury.Byte32ToString(feedIdBytes)

		t.Run(fmt.Sprintf("validate report values for latest block number for %s", feedIdStr),
			func(t *testing.T) {
				evmClient := testEnv.EvmClient
				msClient := testEnv.MSClient

				for i, _ := range simulation.Rounds {
					simulation.Rounds[i].Results = RoundResult{
						Validators: map[string]ValidatorResults{},
					}

					er := expectedResult{
						FeedId: feedIdBytes,
						Value:  simulation.Rounds[i].DataProviderValue,
					}

					setMockserver(t, &testEnv, simulation.Rounds[i].DataProviderValue)

					l.Info().Msgf("Validation round %d starts now! Expected result: %+v", i, er)

					var wg sync.WaitGroup
					var roundDoneChan = make(chan bool)
					var resultChan = make(chan testResult)

					// Start multiple validators requesting for the latest block number reports
					for j := 0; j < simulation.Rounds[i].ValidatorsCount; j++ {
						validatorId := fmt.Sprintf("validator_%d", j)
						simulation.Rounds[i].Results.Validators[validatorId] = ValidatorResults{
							AllResults: []testResult{},
						}

						wg.Add(1)
						go validateNewReportsEveryBlock(
							validatorId, simulation.Rounds[i].Duration,
							simulation.Rounds[i].TimeBetweenNewBlockAndApiCall,
							feedIdStr, er, evmClient, msClient, &wg, resultChan,
						)
					}

					// Wait for validators to finish
					go func() {
						wg.Wait()
						roundDoneChan <- true
					}()

				loop:
					for {
						select {
						case <-roundDoneChan:
							// All validators finished
							// Continue to next validation round
							l.Info().Msgf("Round %d is done", i)
							break loop
						case r := <-resultChan:
							// l.Info().Msgf("Received test result from %s: %+v", r.Id, r)
							v := simulation.Rounds[i].Results.Validators[r.Id]
							v.AllResults = append(v.AllResults, r)
							simulation.Rounds[i].Results.Validators[r.Id] = v
						}
					}
				}

				l.Info().Msgf("All simulations done!")

				checkResults(t, simulation)
			})
	}
}

func setMockserver(t *testing.T, testEnv *mercury.TestEnv, value *big.Int) {
	// Update job spec to not multiply the value by 10? Otherwise division is needed
	err := testEnv.MockserverClient.SetValuePath("/variable", int(value.Int64()/10))
	require.NoError(t, err)
}

func checkResults(t *testing.T, simulation Simulation) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "ID\tDuration\tHeader Sleep\tValidator ID\tRequests\tMissing Reports\tMissing Reports %\tBN\tIncorrect Benchmark\tBid\tAsk\tFeed ID")

	for i, round := range simulation.Rounds {
		for validatorId, validator := range round.Results.Validators {
			requestsCount := validator.RequestsCount()
			missingReports := validator.MissingReports()
			missingReportsP := float32(len(missingReports)) / float32(requestsCount) * 100
			wrongBlockNumbers := validator.WrongBlockNumbers()
			wrongBenchmark := validator.WrongBenchmark()
			wrongBid := validator.WrongBid()
			wrongAsk := validator.WrongAsk()
			wrongFeedId := validator.WrongFeedId()

			// Build test results summary table
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%d\t%d\t%.1f%%\t%d\t%d\t%d\t%d\t%d\n",
				i,
				round.Duration,
				round.TimeBetweenNewBlockAndApiCall,
				validatorId,
				requestsCount,
				len(missingReports),
				missingReportsP,
				len(wrongBlockNumbers),
				len(wrongBenchmark),
				len(wrongBid),
				len(wrongAsk),
				len(wrongFeedId),
			)

			// Fail test if thresholds not met
			if missingReportsP > round.MissingReportsThreshold {
				t.Errorf("Too many missing reports for %s", validatorId)
			}
			if len(wrongBenchmark) > 0 {
				t.Errorf("Wrong benchmark price for %+v", wrongBenchmark)
			}
			if len(wrongAsk) > 0 {
				t.Errorf("Wrong ask price for %+v", wrongAsk)
			}
			if len(wrongBid) > 0 {
				t.Errorf("Wrong bid price for %+v", wrongBid)
			}
			if len(wrongBlockNumbers) > 0 {
				t.Errorf("Wrong block numbers for %+v", wrongBlockNumbers)
			}
			if len(wrongFeedId) > 0 {
				t.Errorf("Wrong feed id %+v", wrongFeedId)
			}
		}
	}

	w.Flush()
}
