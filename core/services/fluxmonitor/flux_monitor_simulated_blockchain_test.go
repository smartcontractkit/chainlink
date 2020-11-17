package fluxmonitor_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flags_wrapper"
	faw "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	goEthereumEth "github.com/ethereum/go-ethereum/eth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const description = "exactly thirty-three characters!!"
const decimals = 8
const fee = int64(100) // Amount paid by FA contract, in LINK-wei
const faTimeout = uint32(1)

var pollTimerPeriod = 200 * time.Millisecond // if failing due to timeouts, increase this
var oneEth = big.NewInt(1000000000000000000)
var emptyList = []common.Address{}

// fluxAggregatorUniverse represents the universe with which the aggregator
// contract interacts
type fluxAggregatorUniverse struct {
	aggregatorContract        *faw.FluxAggregator
	aggregatorContractAddress common.Address
	linkContract              *link_token_interface.LinkToken
	flagsContract             *flags_wrapper.Flags
	flagsContractAddress      common.Address
	// Abstraction representation of the ethereum blockchain
	backend       *backends.SimulatedBackend
	aggregatorABI abi.ABI
	// Cast of participants
	sergey  *bind.TransactOpts // Owns all the LINK initially
	neil    *bind.TransactOpts // Node operator Flux Monitor Oracle
	ned     *bind.TransactOpts // Node operator Flux Monitor Oracle
	nallory *bind.TransactOpts // Node operator Flux Monitor Oracle (Baddie.)
}

// newIdentity returns a go-ethereum abstraction of an ethereum account for
// interacting with contract golang wrappers
func newIdentity(t *testing.T) *bind.TransactOpts {
	key, err := crypto.GenerateKey()
	require.NoError(t, err, "failed to generate ethereum identity")
	return bind.NewKeyedTransactor(key)
}

// setupFluxAggregatorUniverse returns a fully initialized fluxAggregator universe. The
// arguments match the arguments of the same name in the FluxAggregator
// constructor.
func setupFluxAggregatorUniverse(t *testing.T) fluxAggregatorUniverse {
	var f fluxAggregatorUniverse
	f.sergey = newIdentity(t)
	f.neil = newIdentity(t)
	f.ned = newIdentity(t)
	f.nallory = cltest.OracleTransactor
	genesisData := core.GenesisAlloc{
		f.sergey.From:  {Balance: oneEth},
		f.neil.From:    {Balance: oneEth},
		f.ned.From:     {Balance: oneEth},
		f.nallory.From: {Balance: oneEth},
	}
	gasLimit := goEthereumEth.DefaultConfig.Miner.GasCeil * 2
	f.backend = backends.NewSimulatedBackend(genesisData, gasLimit)
	var err error
	f.aggregatorABI, err = abi.JSON(strings.NewReader(faw.FluxAggregatorABI))
	require.NoError(t, err, "could not parse FluxAggregator ABI")

	var linkAddress common.Address
	linkAddress, _, f.linkContract, err = link_token_interface.DeployLinkToken(f.sergey, f.backend)
	require.NoError(t, err, "failed to deploy link contract to simulated ethereum blockchain")

	f.flagsContractAddress, _, f.flagsContract, err = flags_wrapper.DeployFlags(f.sergey, f.backend, f.sergey.From)
	require.NoError(t, err, "failed to deploy flags contract to simulated ethereum blockchain")

	f.backend.Commit()

	// FluxAggregator contract subtracts timeout from block timestamp, which will
	// be less than the timeout, leading to a SafeMath error. Wait for longer than
	// the timeout... Golang is unpleasant about mixing int64 and time.Duration in
	// arithmetic operations, so do everything as int64 and then convert.
	waitTimeMs := int64(faTimeout * 5000)
	time.Sleep(time.Duration((waitTimeMs + waitTimeMs/20) * int64(time.Millisecond)))
	oldGasLimit := f.sergey.GasLimit
	f.sergey.GasLimit = gasLimit
	minSubmissionValue := big.NewInt(0)
	maxSubmissionValue := big.NewInt(100000000000)
	f.aggregatorContractAddress, _, f.aggregatorContract, err = faw.DeployFluxAggregator(
		f.sergey,
		f.backend,
		linkAddress,
		big.NewInt(fee),
		faTimeout,
		common.Address{},
		minSubmissionValue,
		maxSubmissionValue,
		decimals,
		description,
	)
	f.backend.Commit() // Must commit contract to chain before we can fund with LINK
	require.NoError(t, err, "failed to deploy FluxAggregator contract to simulated ethereum blockchain")

	f.sergey.GasLimit = oldGasLimit

	_, err = f.linkContract.Transfer(f.sergey, f.aggregatorContractAddress, oneEth) // Actually, LINK
	require.NoError(t, err, "failed to fund FluxAggregator contract with LINK")

	_, err = f.aggregatorContract.UpdateAvailableFunds(f.sergey)
	require.NoError(t, err, "failed to update aggregator's availableFunds field")

	f.backend.Commit()
	availableFunds, err := f.aggregatorContract.AvailableFunds(nil)
	require.NoError(t, err, "failed to retrieve AvailableFunds")
	require.Equal(t, availableFunds, oneEth)

	ilogs, err := f.aggregatorContract.FilterAvailableFundsUpdated(nil, []*big.Int{oneEth})
	require.NoError(t, err, "failed to gather AvailableFundsUpdated logs")

	logs := cltest.GetLogs(t, nil, ilogs)
	require.Len(t, logs, 1, "a single AvailableFundsUpdated log should be emitted")

	return f
}

// checkOraclesAdded asserts that the correct logs were emitted for each oracle added
func checkOraclesAdded(t *testing.T, f fluxAggregatorUniverse, oracleList []common.Address) {
	iaddedLogs, err := f.aggregatorContract.FilterOraclePermissionsUpdated(nil, oracleList, []bool{true})
	require.NoError(t, err, "failed to gather OraclePermissionsUpdated logs")

	addedLogs := cltest.GetLogs(t, nil, iaddedLogs)
	require.Len(t, addedLogs, len(oracleList), "should have log for each oracle")

	iadminLogs, err := f.aggregatorContract.FilterOracleAdminUpdated(nil, oracleList, oracleList)
	require.NoError(t, err, "failed to gather OracleAdminUpdated logs")

	adminLogs := cltest.GetLogs(t, nil, iadminLogs)
	require.Len(t, adminLogs, len(oracleList), "should have log for each oracle")

	for oracleIdx, oracle := range oracleList {
		require.Equal(t, oracle, addedLogs[oracleIdx].(*faw.FluxAggregatorOraclePermissionsUpdated).Oracle, "log for wrong oracle emitted")
		require.Equal(t, oracle, adminLogs[oracleIdx].(*faw.FluxAggregatorOracleAdminUpdated).Oracle, "log for wrong oracle emitted")
	}
}

// type answerParams struct {
//     fa                          *fluxAggregatorUniverse
//     roundId, answer             int64
//     from                        *bind.TransactOpts
//     isNewRound, completesAnswer bool
// }

// // checkSubmission verifies all the logs emitted by fa's FluxAggregator
// // contract after an updateAnswer with the given values.
// func checkSubmission(t *testing.T, p answerParams,
//     currentBalance int64, receiptBlock uint64) {
//     t.Helper()
//     if receiptBlock == 0 {
//         receiptBlock = p.fa.backend.Blockchain().CurrentBlock().Number().Uint64()
//     }
//     fromBlock := &bind.FilterOpts{Start: receiptBlock, End: &receiptBlock}

//     // Could filter for the known values here, but while that would be more
//     // succinct it leads to less informative error messages... Did the log not
//     // appear at all, or did it just have a wrong value?
//     ilogs, err := p.fa.aggregatorContract.FilterSubmissionReceived(fromBlock, []*big.Int{}, []uint32{}, []common.Address{})
//     require.NoError(t, err, "failed to get SubmissionReceived logs")

//     var srlogs []*faw.FluxAggregatorSubmissionReceived
//     _ = cltest.GetLogs(t, &srlogs, ilogs)
//     require.Len(t, srlogs, 1, "FluxAggregator did not emit correct SubmissionReceived log")
//     assert.Equal(t, uint32(p.roundId), srlogs[0].Round, "SubmissionReceived log has wrong round")
//     assert.Equal(t, p.from.From, srlogs[0].Oracle, "SubmissionReceived log has wrong oracle")

//     inrlogs, err := p.fa.aggregatorContract.FilterNewRound(fromBlock, []*big.Int{}, []common.Address{})
//     require.NoError(t, err, "failed to get NewRound logs")

//     if p.isNewRound {
//         var nrlogs []*faw.FluxAggregatorNewRound
//         cltest.GetLogs(t, &nrlogs, inrlogs)
//         require.Len(t, nrlogs, 1, "FluxAggregator did not emit correct NewRound log")
//         assert.Equal(t, p.roundId, nrlogs[0].RoundId.Int64(), "NewRound log has wrong roundId")
//         assert.Equal(t, p.from.From, nrlogs[0].StartedBy, "NewRound log started by wrong oracle")
//     } else {
//         assert.Len(t, cltest.GetLogs(t, nil, inrlogs), 0, "FluxAggregator emitted unexpected NewRound log")
//     }

//     iaflogs, err := p.fa.aggregatorContract.FilterAvailableFundsUpdated(fromBlock, []*big.Int{})
//     require.NoError(t, err, "failed to get AvailableFundsUpdated logs")

//     var aflogs []*faw.FluxAggregatorAvailableFundsUpdated
//     _ = cltest.GetLogs(t, &aflogs, iaflogs)
//     assert.Len(t, aflogs, 1, "FluxAggregator did not emit correct AvailableFundsUpdated log")
//     assert.Equal(t, currentBalance-fee, aflogs[0].Amount.Int64(), "AvailableFundsUpdated log has wrong amount")

//     iaulogs, err := p.fa.aggregatorContract.FilterAnswerUpdated(fromBlock,
//         []*big.Int{big.NewInt(p.answer)}, []*big.Int{big.NewInt(p.roundId)})
//     require.NoError(t, err, "failed to get AnswerUpdated logs")
//     if p.completesAnswer {
//         var aulogs []*faw.FluxAggregatorAnswerUpdated
//         _ = cltest.GetLogs(t, &aulogs, iaulogs)
//         assert.Len(t, aulogs, 1, "FluxAggregator did not emit correct AnswerUpdated log")
//         assert.Equal(t, p.roundId, aulogs[0].RoundId.Int64(), "AnswerUpdated log has wrong roundId")
//         assert.Equal(t, p.answer, aulogs[0].Current.Int64(), "AnswerUpdated log has wrong current value")
//     }
// }

// // currentbalance returns the current balance of fa's FluxAggregator
// func currentBalance(t *testing.T, fa *fluxAggregatorUniverse) *big.Int {
//     currentBalance, err := fa.aggregatorContract.AvailableFunds(nil)
//     require.NoError(t, err, "failed to get current FA balance")
//     return currentBalance
// }

// // submitAnswer simulates a call to fa's FluxAggregator contract from from, with
// // the given roundId and answer, and checks that all the logs emitted by the
// // contract are correct
// func submitAnswer(t *testing.T, p answerParams) {
//     cb := currentBalance(t, p.fa)
//     _, err := p.fa.aggregatorContract.Submit(p.from, big.NewInt(p.roundId), big.NewInt(p.answer))
//     require.NoError(t, err, "failed to initialize first flux aggregation round:")

//     p.fa.backend.Commit()
//     checkSubmission(t, p, cb.Int64(), 0)
// }

// type maliciousFluxMonitor interface {
//     CreateJob(t *testing.T, jobSpecId *models.ID, polledAnswer decimal.Decimal, nextRound *big.Int) error
// }

func waitForRunsAndEthTxCount(
	t *testing.T,
	job models.JobSpec,
	runCount int,
	app *cltest.TestApplication,
	backend *backends.SimulatedBackend,
) []models.JobRun {
	t.Helper()
	store := app.Store
	jrs := cltest.WaitForRuns(t, job, store, runCount) // Submit answer from
	app.EthBroadcaster.Trigger()
	txes := cltest.WaitForEthTxCount(t, store, runCount)
	txas := cltest.WaitForEthTxAttemptsForEthTx(t, store, txes[0])
	cltest.WaitForTxInMempool(t, backend, txas[0].Hash)

	backend.Commit()
	return jrs
}

// TODO: This test is non-deterministic and needs to be rewritten or rethought
// See: https://www.pivotaltracker.com/story/show/175757546
// func TestFluxMonitorAntiSpamLogic(t *testing.T) {
//     t.Skip()
//     // Comments starting with "-" describe the steps this test executes.

//     // - deploy a brand new FM contract
//     fa := setupFluxAggregatorUniverse(t)

//     // - add oracles
//     oracleList := []common.Address{fa.neil.From, fa.ned.From, fa.nallory.From}
//     _, err := fa.aggregatorContract.ChangeOracles(fa.sergey, emptyList, oracleList, oracleList, 2, 3, 2)
//     assert.NoError(t, err, "failed to add oracles to aggregator")
//     fa.backend.Commit()
//     checkOraclesAdded(t, fa, oracleList)

//     // Set up chainlink app
//     config, cfgCleanup := cltest.NewConfig(t)
//     config.Config.Set("DEFAULT_HTTP_TIMEOUT", "100ms")
//     config.Config.Set("TRIGGER_FALLBACK_DB_POLL_INTERVAL", "1s")
//     defer cfgCleanup()
//     app, cleanup := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, fa.backend)
//     defer cleanup()
//     require.NoError(t, app.StartAndConnect())
//     minFee := app.Store.Config.MinimumContractPayment().ToInt().Int64()
//     require.Equal(t, fee, minFee, "fee paid by FluxAggregator (%d) must at "+
//         "least match MinimumContractPayment (%s). (Which is currently set in "+
//         "cltest.go.)", fee, minFee)

//     answer := int64(1) // Answer the nodes give on the first round

//     //- have one of the fake nodes start a round.
//     roundId := int64(1)
//     processedAnswer := answer * 100 // [> job has multiply times 100 <]
//     submitAnswer(t, answerParams{
//         fa:              &fa,
//         roundId:         roundId,
//         answer:          processedAnswer,
//         from:            fa.neil,
//         isNewRound:      true,
//         completesAnswer: false,
//     })

//     // - successfully close the round through the submissions of the other nodes
//     // Response by malicious chainlink node, nallory
//     initialBalance := currentBalance(t, &fa).Int64()
//     reportPrice := answer
//     priceResponse := func() string {
//         return fmt.Sprintf(`{"data":{"result": %d}}`, atomic.LoadInt64(&reportPrice))
//     }
//     mockServer := cltest.NewHTTPMockServerWithAlterableResponse(t, priceResponse)
//     defer mockServer.Close()

//     // When event appears on submissionReceived, flux monitor job run is complete
//     submissionReceived := make(chan *faw.FluxAggregatorSubmissionReceived)
//     subscription, err := fa.aggregatorContract.WatchSubmissionReceived(
//         nil,
//         submissionReceived,
//         []*big.Int{},
//         []uint32{},
//         []common.Address{fa.nallory.From},
//     )
//     require.NoError(t, err, "failed to subscribe to SubmissionReceived events")
//     defer subscription.Unsubscribe()

//     // Create FM Job, and wait for job run to start (the above UpdateAnswer call
//     // to FluxAggregator contract initiates a run.)
//     buffer := cltest.MustReadFile(t, "../../internal/testdata/flux_monitor_job.json")
//     var job models.JobSpec
//     require.NoError(t, json.Unmarshal(buffer, &job))
//     initr := &job.Initiators[0]
//     initr.InitiatorParams.Feeds = cltest.JSONFromString(t, fmt.Sprintf(`["%s"]`, mockServer.URL))
//     initr.InitiatorParams.PollTimer.Period = models.MustMakeDuration(pollTimerPeriod)
//     initr.InitiatorParams.Address = fa.aggregatorContractAddress

//     j := cltest.CreateJobSpecViaWeb(t, app, job)
//     jrs := waitForRunsAndEthTxCount(t, j, 1, app, fa.backend)

//     reportedPrice := jrs[0].RunRequest.RequestParams.Get("result").String()
//     assert.Equal(t, reportedPrice, fmt.Sprintf("%d", atomic.LoadInt64(&reportPrice)), "failed to report correct price to contract")
//     var receiptBlock uint64
//     select { // block until FluxAggregator contract acknowledges chainlink message
//     case log := <-submissionReceived:
//         receiptBlock = log.Raw.BlockNumber
//     case <-time.After(pollTimerPeriod):
//         t.Fatalf("chainlink failed to submit answer to FluxAggregator contract")
//     }
//     checkSubmission(t,
//         answerParams{
//             fa:              &fa,
//             roundId:         roundId,
//             answer:          processedAnswer,
//             from:            fa.nallory,
//             isNewRound:      false,
//             completesAnswer: true},
//         initialBalance,
//         receiptBlock,
//     )

//     //- have the malicious node start the next round.
//     nextRoundBalance := initialBalance - fee
//     // Triggers a new round, since price deviation exceeds threshold
//     atomic.StoreInt64(&reportPrice, answer+1)

//     waitForRunsAndEthTxCount(t, j, 2, app, fa.backend)

//     select {
//     case log := <-submissionReceived:
//         receiptBlock = log.Raw.BlockNumber
//     case <-time.After(2 * pollTimerPeriod):
//         t.Fatalf("chainlink failed to submit answer to FluxAggregator contract")
//     }
//     newRound := roundId + 1
//     processedAnswer = 100 * atomic.LoadInt64(&reportPrice)
//     checkSubmission(t,
//         answerParams{
//             fa:              &fa,
//             roundId:         newRound,
//             answer:          processedAnswer,
//             from:            fa.nallory,
//             isNewRound:      true,
//             completesAnswer: false},
//         nextRoundBalance,
//         receiptBlock,
//     )

//     // Successfully close the round through the submissions of the other nodes
//     submitAnswer(t,
//         answerParams{
//             fa:              &fa,
//             roundId:         newRound,
//             answer:          processedAnswer,
//             from:            fa.neil,
//             isNewRound:      false,
//             completesAnswer: true},
//     )

//     // Have the malicious node try to start another round repeatedly until the
//     // roundDelay is reached, making sure that it isn't successful
//     newRound = newRound + 1
//     processedAnswer = 100 * atomic.LoadInt64(&reportPrice)
//     precision := job.Initiators[0].InitiatorParams.Precision
//     // FORCE node to try to start a new round
//     err = app.FluxMonitor.(maliciousFluxMonitor).CreateJob(t, j.ID, decimal.New(processedAnswer, precision), big.NewInt(newRound))
//     require.NoError(t, err)

//     waitForRunsAndEthTxCount(t, j, 3, app, fa.backend)

//     select {
//     case <-submissionReceived:
//         t.Fatalf("FA allowed chainlink node to start a new round early")
//     case <-time.After(5 * pollTimerPeriod):
//     }
//     // Remove the record of the submitted round, or else FM's reorg protection will cause the test to fail
//     err = app.Store.DeleteFluxMonitorRoundsBackThrough(fa.aggregatorContractAddress, uint32(newRound))
//     require.NoError(t, err)

//     // Try to start a new round directly, should fail
//     _, err = fa.aggregatorContract.RequestNewRound(fa.nallory)
//     assert.Error(t, err, "FA allowed chainlink node to start a new round early")

//     //- finally, ensure it can start a legitimate round after roundDelay is reached
//     // start an intervening round
//     submitAnswer(t, answerParams{fa: &fa, roundId: newRound,
//         answer: processedAnswer, from: fa.ned, isNewRound: true,
//         completesAnswer: false})
//     submitAnswer(t, answerParams{fa: &fa, roundId: newRound,
//         answer: processedAnswer, from: fa.neil, isNewRound: false,
//         completesAnswer: true})
//     // start a legitimate new round
//     atomic.StoreInt64(&reportPrice, reportPrice+3)
//     waitForRunsAndEthTxCount(t, j, 4, app, fa.backend)

//     select {
//     case <-submissionReceived:
//     case <-time.After(5 * pollTimerPeriod):
//         t.Fatalf("could not start a new round, even though delay has passed")
//     }
// }

func TestFluxMonitor_HibernationMode(t *testing.T) {
	fa := setupFluxAggregatorUniverse(t)

	// - add oracles
	oracleList := []common.Address{fa.nallory.From}
	_, err := fa.aggregatorContract.ChangeOracles(fa.sergey, emptyList, oracleList, oracleList, 1, 1, 0)
	assert.NoError(t, err, "failed to add oracles to aggregator")
	fa.backend.Commit()
	checkOraclesAdded(t, fa, oracleList)

	// Set up chainlink app
	config, cfgCleanup := cltest.NewConfig(t)
	config.Config.Set("DEFAULT_HTTP_TIMEOUT", "100ms")
	config.Config.Set("FLAGS_CONTRACT_ADDRESS", fa.flagsContractAddress.Hex())
	config.Config.Set("TRIGGER_FALLBACK_DB_POLL_INTERVAL", "1s")
	defer cfgCleanup()
	app, cleanup := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, fa.backend)
	defer cleanup()
	require.NoError(t, app.StartAndConnect(), "failed to start chainlink")

	// // create mock server
	reportPrice := int64(1)
	priceResponse := func() string {
		return fmt.Sprintf(`{"data":{"result": %d}}`, atomic.LoadInt64(&reportPrice))
	}
	mockServer := cltest.NewHTTPMockServerWithAlterableResponse(t, priceResponse)
	defer mockServer.Close()

	// // When event appears on submissionReceived, flux monitor job run is complete
	submissionReceived := make(chan *faw.FluxAggregatorSubmissionReceived)
	subscription, err := fa.aggregatorContract.WatchSubmissionReceived(
		nil,
		submissionReceived,
		[]*big.Int{},
		[]uint32{},
		[]common.Address{fa.nallory.From},
	)
	require.NoError(t, err, "failed to subscribe to SubmissionReceived events")
	defer subscription.Unsubscribe()

	// // Create FM Job, and wait for job run to start (the above UpdateAnswer call
	// // to FluxAggregator contract initiates a run.)
	buffer := cltest.MustReadFile(t, "../../internal/testdata/flux_monitor_job.json")
	var job models.JobSpec
	require.NoError(t, json.Unmarshal(buffer, &job))
	initr := &job.Initiators[0]
	initr.InitiatorParams.Feeds = cltest.JSONFromString(t, fmt.Sprintf(`["%s"]`, mockServer.URL))
	initr.InitiatorParams.PollTimer.Period = models.MustMakeDuration(pollTimerPeriod)
	initr.InitiatorParams.Address = fa.aggregatorContractAddress

	// raise flags
	fa.flagsContract.RaiseFlag(fa.sergey, utils.ZeroAddress) // global kill switch
	fa.flagsContract.RaiseFlag(fa.sergey, initr.Address)
	fa.backend.Commit()

	job = cltest.CreateJobSpecViaWeb(t, app, job)

	// node doesn't submit initial response, because flag is up
	cltest.AssertRunsStays(t, job, app.Store, 0)

	// lower global kill switch flag - should trigger job run
	fa.flagsContract.LowerFlags(fa.sergey, []common.Address{utils.ZeroAddress})
	fa.backend.Commit()
	waitForRunsAndEthTxCount(t, job, 1, app, fa.backend)

	// change in price should trigger run
	reportPrice = int64(2)
	waitForRunsAndEthTxCount(t, job, 2, app, fa.backend)

	// lower contract's flag - should have no effect (but currently does)
	// TODO - https://www.pivotaltracker.com/story/show/175419789
	fa.flagsContract.LowerFlags(fa.sergey, []common.Address{initr.Address})
	fa.backend.Commit()
	waitForRunsAndEthTxCount(t, job, 3, app, fa.backend)

	// change in price should trigger run
	reportPrice = int64(4)
	waitForRunsAndEthTxCount(t, job, 4, app, fa.backend)

	// raise both flags
	fa.flagsContract.RaiseFlag(fa.sergey, initr.Address)
	fa.flagsContract.RaiseFlag(fa.sergey, utils.ZeroAddress)
	fa.backend.Commit()

	// wait for FM to receive flags raised logs
	assert.Eventually(t, func() bool {
		ilogs, err := fa.flagsContract.FilterFlagRaised(nil, []common.Address{})
		require.NoError(t, err)
		logs := cltest.GetLogs(t, nil, ilogs)
		return len(logs) == 4
	}, 5*time.Second, 100*time.Millisecond)

	// change in price should not trigger run
	reportPrice = int64(8)
	cltest.AssertRunsStays(t, job, app.Store, 4)
}
