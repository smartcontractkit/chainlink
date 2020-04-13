package internal_test

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"chainlink/core/internal/cltest"
	faw "chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"chainlink/core/internal/gethwrappers/generated/link_token_interface"
	"chainlink/core/store/models"

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

// fluxAggregator represents the universe with which the aggregator contract
// interacts
type fluxAggregator struct {
	aggregatorContract        *faw.FluxAggregator
	aggregatorContractAddress common.Address
	linkContract              *link_token_interface.LinkToken
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

var oneEth = big.NewInt(1000000000000000000)
var fee = big.NewInt(100) // Amount paid by FA contract, in LINK-wei

// deployFluxAggregator returns a fully initialized fluxAggregator universe. The
// arguments match the arguments of the same name in the FluxAggregator
// constructor.
func deployFluxAggregator(t *testing.T, paymentAmount *big.Int, timeout uint32,
	decimals uint8, description [32]byte) fluxAggregator {
	var f fluxAggregator
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
	gasLimit := goEthereumEth.DefaultConfig.Miner.GasCeil
	f.backend = backends.NewSimulatedBackend(genesisData, gasLimit)
	var err error
	f.aggregatorABI, err = abi.JSON(strings.NewReader(
		faw.FluxAggregatorABI))
	require.NoError(t, err, "could not parse FluxAggregator ABI")
	var linkAddress common.Address
	linkAddress, _, f.linkContract, err = link_token_interface.DeployLinkToken(
		f.sergey, f.backend)
	require.NoError(t, err,
		"failed to deploy link contract to simulated ethereum blockchain")
	f.backend.Commit()
	// FluxAggregator contract subtracts timeout from block timestamp, which will
	// be less than the timeout, leading to a SafeMath error. Wait for longer than
	// the timeout... Golang is unpleasant about mixing int64 and time.Duration in
	// arithmetic operations, so do everything as int64 and then convert.
	waitTimeMs := int64(timeout * 1000)
	time.Sleep(time.Duration((waitTimeMs + waitTimeMs/20) * int64(time.Millisecond)))
	f.aggregatorContractAddress, _, f.aggregatorContract, err =
		faw.DeployFluxAggregator(f.sergey, f.backend,
			linkAddress, paymentAmount, timeout, decimals, description)
	f.backend.Commit() // Must commit contract to chain before we can fund with LINK
	require.NoError(t, err,
		"failed to deploy FluxAggregator contract to simulated ethereum blockchain")
	_, err = f.linkContract.Transfer(f.sergey, f.aggregatorContractAddress,
		oneEth) // Actually, LINK
	require.NoError(t, err, "failed to fund FluxAggregator contract with LINK")
	_, err = f.aggregatorContract.UpdateAvailableFunds(f.sergey)
	require.NoError(t, err, "failed to update aggregator's vailableFunds field")
	f.backend.Commit()
	availableFunds, err := f.aggregatorContract.AvailableFunds(nil)
	require.NoError(t, err, "failed to retrieve AvailableFunds")
	require.Equal(t, availableFunds, oneEth)
	ilogs, err := f.aggregatorContract.FilterAvailableFundsUpdated(nil,
		[]*big.Int{oneEth})
	require.NoError(t, err, "failed to gather AvailableFundsUpdated logs")
	logs := cltest.GetLogs(ilogs)
	require.Len(t, logs, 1, "a single AvailableFundsUpdated log should be emitted")

	oracleList := []common.Address{f.neil.From, f.ned.From, f.nallory.From}
	_, err = f.aggregatorContract.AddOracles(
		f.sergey, oracleList, oracleList, 2, 3, 2)
	require.NoError(t, err, "failed to update oracles list")
	f.backend.Commit()
	iaddedLogs, err := f.aggregatorContract.FilterOracleAdded(nil, oracleList)
	require.NoError(t, err, "failed to gather OracleAdded logs")
	addedLogs := cltest.GetLogs(iaddedLogs)
	require.Len(t, addedLogs, len(oracleList), "should have log for each oracle")
	iadminLogs, err := f.aggregatorContract.FilterOracleAdminUpdated(nil,
		oracleList, oracleList)
	require.NoError(t, err, "failed to gather OracleAdminUpdated logs")
	adminLogs := cltest.GetLogs(iadminLogs)
	require.Len(t, adminLogs, len(oracleList), "should have log for each oracle")
	for oracleIdx, oracle := range oracleList {
		require.Equal(t, oracle,
			addedLogs[oracleIdx].(*faw.FluxAggregatorOracleAdded).Oracle,
			"log for wrong oracle emitted")
		require.Equal(t, oracle,
			adminLogs[oracleIdx].(*faw.FluxAggregatorOracleAdminUpdated).Oracle,
			"log for wrong oracle emmitted")
	}
	return f
}

// checkUpdateAnswer verifies all the logs emitted by fa's FluxAggregator
// contract after an updateAnswer with the given values.
func checkUpdateAnswer(t *testing.T, fa *fluxAggregator, roundId,
	answer, currentBalance *big.Int, from *bind.TransactOpts, isNewRound,
	completesAnswer bool, receiptBlock uint64) {
	if receiptBlock == 0 {
		receiptBlock = fa.backend.Blockchain().CurrentBlock().Number().Uint64()
	}
	fromBlock := &bind.FilterOpts{Start: receiptBlock, End: &receiptBlock}
	ilogs, err := fa.aggregatorContract.FilterSubmissionReceived(fromBlock,
		[]*big.Int{answer}, []uint32{uint32(roundId.Uint64())},
		[]common.Address{from.From})
	require.NoError(t, err, "failed to get SubmissionReceived logs")
	assert.Len(t, cltest.GetLogs(ilogs), 1,
		"FluxAggregator did not emit correct SubmissionReceived log")
	inrlogs, err := fa.aggregatorContract.FilterNewRound(fromBlock,
		[]*big.Int{roundId}, []common.Address{from.From})
	require.NoError(t, err, "failed to get NewRound logs")
	if isNewRound {
		assert.Len(t, cltest.GetLogs(inrlogs), 1,
			"FluxAggregator did not emit correct NewRound log")
	} else {
		assert.Len(t, cltest.GetLogs(inrlogs), 0,
			"FluxAggregator emitted unexpected NewRound log")
	}
	iaflogs, err := fa.aggregatorContract.FilterAvailableFundsUpdated(fromBlock,
		[]*big.Int{big.NewInt(0).Sub(currentBalance, fee)})
	require.NoError(t, err, "failed to get AvailableFundsUpdated logs")
	assert.Len(t, cltest.GetLogs(iaflogs), 1,
		"FluxAggregator did not emit correct AvailableFundsUpdated log")
	iaulogs, err := fa.aggregatorContract.FilterAnswerUpdated(fromBlock,
		[]*big.Int{answer}, []*big.Int{roundId})
	require.NoError(t, err, "failed to get AnswerUpdated logs")
	if completesAnswer {
		assert.Len(t, cltest.GetLogs(iaulogs), 1,
			"FluxAggregator did not emit correct AnswerUpdated log")
	}
}

// currentBalance returns the current balance of fa's FluxAggregator
func currentBalance(t *testing.T, fa *fluxAggregator) *big.Int {
	currentBalance, err := fa.aggregatorContract.AvailableFunds(nil)
	require.NoError(t, err, "failed to get current FA balance")
	return currentBalance
}

// updateAnswer simulates a call to fa's FluxAggregator contract from from, with
// the given roundId and answer, and checks that all the logs emitted by the
// contract are correct
func updateAnswer(t *testing.T, fa *fluxAggregator, roundId, answer *big.Int,
	from *bind.TransactOpts, isNewRound, completesAnswer bool) {
	cb := currentBalance(t, fa)
	_, err := fa.aggregatorContract.UpdateAnswer(from, roundId, answer)
	require.NoError(t, err, "failed to initialize first flux aggregation round:")
	fa.backend.Commit()
	checkUpdateAnswer(t, fa, roundId, answer, cb, from, isNewRound,
		completesAnswer, 0)
}

var currentPrice int

var mockServer *httptest.Server

func init() {
	mockServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, fmt.Sprintf("%d", currentPrice))
		}))
}

//- successfully close the round through the submissions of the other nodes
//- have the malicious node try to start another round repeatedly until the roundDelay is reached, making sure that it isn't successful
//- finally, ensure it can start a legitimate round after roundDelay is reached
func TestFluxMonitorAntiSpamLogic(t *testing.T) {
	// Comments starting with "-" describe the steps this test executes.

	// - deploy a brand new FM contract
	var description [32]byte
	copy(description[:], "exactly thirty-two characters!!!")
	fa := deployFluxAggregator(t, fee, 1, 8, description)

	// Set up chainlink app
	config, cfgCleanup := cltest.NewConfig(t)
	defer cfgCleanup()
	app, cleanup := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t,
		config, fa.backend)
	defer cleanup()
	require.NoError(t, app.StartAndConnect(), "failed to start chainlink")
	minFee := app.Store.Config.MinimumContractPayment().ToInt()
	require.True(t, fee.Cmp(minFee) >= 0, "fee paid by FluxAggregator (%d) must "+
		"at least match MinimumContractPayment (%s). (Which is currently set in "+
		"cltest.go.)", fee, minFee)

	answer := int64(1) // Answer the nodes give on the first round

	//- have one of the fake nodes start a round.
	roundId := big.NewInt(1)
	processedAnswer := big.NewInt(answer * 100 /* job has multiply times 100 */)
	updateAnswer(t, &fa, roundId, processedAnswer, fa.neil, true, false)

	// - successfully close the round through the submissions of the other nodes
	// Response by malicious chainlink node, nallory
	initialBalance := currentBalance(t, &fa)
	reportPrice := answer
	priceResponse := fmt.Sprintf(`{"data":{"result": %d}}`, reportPrice)
	mockServer, assertCalled := cltest.NewHTTPMockServer(t, http.StatusOK, "POST",
		priceResponse)
	defer assertCalled()

	// When event appears on submissionReceived, flux monitor job run is complete
	submissionReceived := make(chan *faw.FluxAggregatorSubmissionReceived)
	subscription, err := fa.aggregatorContract.WatchSubmissionReceived(nil,
		submissionReceived, []*big.Int{ /* big.NewInt(int64(reportPrice)) */ }, // XXX: I would expect the reportPrice, here
		[]uint32{uint32(roundId.Uint64())}, []common.Address{fa.nallory.From})
	require.NoError(t, err, "failed to subscribe to SubmissionReceived events")

	// Create FM Job, and wait for job run to start (the above UpdateAnswer calls
	// to FluxAggregator contract initiate a run.)
	//
	// Emits SubmissionReceived, AnswerUpdated and AvailableFundsUpdated
	buffer := cltest.MustReadFile(t, "testdata/flux_monitor_job.json")
	var job models.JobSpec
	require.NoError(t, json.Unmarshal(buffer, &job))
	job.Initiators[0].InitiatorParams.Feeds = cltest.JSONFromString(t, fmt.Sprintf(`["%s"]`, mockServer.URL))
	job.Initiators[0].InitiatorParams.PollingInterval = models.Duration(15 * time.Second)
	job.Initiators[0].InitiatorParams.Address = fa.aggregatorContractAddress
	j := cltest.CreateJobSpecViaWeb(t, app, job)
	jrs := cltest.WaitForRuns(t, j, app.Store, 1) // Submit answer from
	reportedPrice := jrs[0].RunRequest.RequestParams.Get("result").String()
	assert.Equal(t, reportedPrice, fmt.Sprintf("%d", reportPrice),
		"failed to report correct price to contract")
	var receiptBlock uint64
	select { // block until FluxAggregator contract acknowledges chainlink message
	case log := <-submissionReceived:
		receiptBlock = log.Raw.BlockNumber
	case <-time.After(2 * time.Second):
		t.Fatalf("chainlink failed to submit answer to FluxAggregator contract")
	}
	subscription.Unsubscribe()
	checkUpdateAnswer(t, &fa, roundId, processedAnswer, initialBalance,
		fa.nallory, false, true, receiptBlock)

	//- have the malicious node start the next round.
	initialBalance = initialBalance.Sub(initialBalance, fee)
	reportPrice = answer + 1
	priceResponse = fmt.Sprintf(`{"data":{"result": %d}}`, reportPrice)
	mockServer, assertCalled2 := cltest.NewHTTPMockServer(t, http.StatusOK, "POST",
		priceResponse)
	defer assertCalled2()

}

// XAU/XAG happened partly because you can update the entire state all at once.
// Having to add oracles one-by-one slows you down, so you can avoid some
// mistakes.
