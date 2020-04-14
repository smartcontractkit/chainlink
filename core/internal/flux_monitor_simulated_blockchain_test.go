package internal_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	faw "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitor"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	goEthereumEth "github.com/ethereum/go-ethereum/eth"
	"github.com/shopspring/decimal"
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

	// Add the participating oracles. Ends up with minAnswers=restartDelay=2,
	// maxAnswers=3
	oracleList := []common.Address{f.neil.From, f.ned.From, f.nallory.From}
	for numOracles, o := range oracleList {
		n := uint32(numOracles)
		_, err = f.aggregatorContract.AddOracle(f.sergey, o, o, n, n+1, n)
		require.NoError(t, err, "failed to update oracles list")
	}
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
	// Could filter for the known values here, but while that would be more
	// succinct it leads to less informative error messages... Did the log not
	// appear at all, or did it just have a wrong value?
	ilogs, err := fa.aggregatorContract.FilterSubmissionReceived(fromBlock,
		[]*big.Int{}, []uint32{}, []common.Address{})
	require.NoError(t, err, "failed to get SubmissionReceived logs")
	srlogs := cltest.GetLogs(ilogs)
	assert.Len(t, srlogs, 1,
		"FluxAggregator did not correct SubmissionReceived log")
	srlog := srlogs[0].(*faw.FluxAggregatorSubmissionReceived)
	assert.True(t, srlog.Answer.Cmp(answer) == 0,
		"SubmissionReceived log has wrong answer")
	assert.Equal(t, uint32(roundId.Int64()), srlog.Round,
		"SubmissionReceived log has wrong round")
	assert.Equal(t, from.From, srlog.Oracle,
		"SubmissionReceived log has wrong oracle")
	inrlogs, err := fa.aggregatorContract.FilterNewRound(fromBlock, []*big.Int{},
		[]common.Address{})
	require.NoError(t, err, "failed to get NewRound logs")
	if isNewRound {
		nrlogs := cltest.GetLogs(inrlogs)
		require.Len(t, nrlogs, 1,
			"FluxAggregator did not emit correct NewRound log")
		nrlog := nrlogs[0].(*faw.FluxAggregatorNewRound)
		assert.Equal(t, roundId, nrlog.RoundId, "NewRound log has wrong roundId")
		assert.Equal(t, from.From, nrlog.StartedBy,
			"NewRound log started by wrong oracle")
	} else {
		assert.Len(t, cltest.GetLogs(inrlogs), 0,
			"FluxAggregator emitted unexpected NewRound log")
	}
	iaflogs, err := fa.aggregatorContract.FilterAvailableFundsUpdated(fromBlock,
		[]*big.Int{})
	require.NoError(t, err, "failed to get AvailableFundsUpdated logs")
	aflogs := cltest.GetLogs(iaflogs)
	assert.Len(t, aflogs, 1,
		"FluxAggregator did not emit correct AvailableFundsUpdated log")
	aflog := aflogs[0].(*faw.FluxAggregatorAvailableFundsUpdated)
	assert.True(t, big.NewInt(0).Sub(currentBalance, fee).Cmp(aflog.Amount) == 0,
		"AvailableFundsUpdated log has wrong amount")
	iaulogs, err := fa.aggregatorContract.FilterAnswerUpdated(fromBlock,
		[]*big.Int{answer}, []*big.Int{roundId})
	require.NoError(t, err, "failed to get AnswerUpdated logs")
	if completesAnswer {
		aulogs := cltest.GetLogs(iaulogs)
		assert.Len(t, aulogs, 1,
			"FluxAggregator did not emit correct AnswerUpdated log")
		aulog := aulogs[0].(*faw.FluxAggregatorAnswerUpdated)
		assert.Equal(t, roundId, aulog.RoundId,
			"AnswerUpdated log has wrong roundId")
		assert.True(t, answer.Cmp(aulog.Current) == 0,
			"AnswerUpdated log has wrong current value")
	}
}

// currentbalance returns the current balance of fa's FluxAggregator
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
	tx, err := fa.aggregatorContract.UpdateAnswer(from, roundId, answer)
	require.NoError(t, err, "failed to initialize first flux aggregation round:")
	fa.backend.Commit()
	receipt, err := fa.backend.TransactionReceipt(context.TODO(), tx.Hash())
	spew.Dump(receipt)
	checkUpdateAnswer(t, fa, roundId, answer, cb, from, isNewRound,
		completesAnswer, 0)
}

func TestFluxMonitorAntiSpamLogic(t *testing.T) {
	// Comments starting with "-" describe the steps this test executes.

	// - deploy a brand new FM contract
	var description [32]byte
	copy(description[:], "exactly thirty-two characters!!!")
	fa := deployFluxAggregator(t, fee, 1, 8, description)

	// Set up chainlink app
	config, cfgCleanup := cltest.NewConfig(t)
	config.Config.Set("DEFAULT_HTTP_TIMEOUT", "100ms")
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
	priceResponse := func() string {
		return fmt.Sprintf(`{"data":{"result": %d}}`, reportPrice)
	}
	mockServer := cltest.NewHTTPMockServerWithAlterableResponse(t, priceResponse)
	defer mockServer.Close()

	// When event appears on submissionReceived, flux monitor job run is complete
	submissionReceived := make(chan *faw.FluxAggregatorSubmissionReceived)
	_, err := fa.aggregatorContract.WatchSubmissionReceived(nil,
		submissionReceived, []*big.Int{}, []uint32{},
		[]common.Address{fa.nallory.From})
	require.NoError(t, err, "failed to subscribe to SubmissionReceived events")

	// Create FM Job, and wait for job run to start (the above UpdateAnswer calls
	// to FluxAggregator contract initiate a run.)
	//
	// Emits SubmissionReceived, AnswerUpdated and AvailableFundsUpdated
	buffer := cltest.MustReadFile(t, "testdata/flux_monitor_job.json")
	var job models.JobSpec
	require.NoError(t, json.Unmarshal(buffer, &job))
	initr := &job.Initiators[0]
	initr.InitiatorParams.Feeds = cltest.JSONFromString(t, fmt.Sprintf(`["%s"]`,
		mockServer.URL))
	initr.InitiatorParams.PollingInterval =
		models.Duration(100 * time.Millisecond)
	initr.InitiatorParams.Address = fa.aggregatorContractAddress
	j := cltest.CreateJobSpecViaWeb(t, app, job)
	jrs := cltest.WaitForRuns(t, j, app.Store, 1) // Submit answer from
	reportedPrice := jrs[0].RunRequest.RequestParams.Get("result").String()
	assert.Equal(t, reportedPrice, fmt.Sprintf("%d", reportPrice),
		"failed to report correct price to contract")
	var receiptBlock uint64
	select { // block until FluxAggregator contract acknowledges chainlink message
	case log := <-submissionReceived:
		receiptBlock = log.Raw.BlockNumber
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("chainlink failed to submit answer to FluxAggregator contract")
	}
	checkUpdateAnswer(t, &fa, roundId, processedAnswer, initialBalance,
		fa.nallory, false, true, receiptBlock)

	//- have the malicious node start the next round.
	initialBalance = initialBalance.Sub(initialBalance, fee)
	// Triggers a new round, since price deviation exceeds threshold
	reportPrice = answer + 1
	select {
	case log := <-submissionReceived:
		receiptBlock = log.Raw.BlockNumber
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("chainlink failed to submit answer to FluxAggregator contract")
	}
	newRound := big.NewInt(0).Add(roundId, big.NewInt(1))
	processedAnswer = big.NewInt(int64(100 * reportPrice))
	checkUpdateAnswer(t, &fa, newRound, processedAnswer, initialBalance,
		fa.nallory, true, false,
		receiptBlock)
	//- successfully close the round through the submissions of the other nodes
	updateAnswer(t, &fa, newRound, processedAnswer, fa.neil, false, true)

	//- have the malicious node try to start another round repeatedly until the
	//roundDelay is reached, making sure that it isn't successful
	// Triggers a new round, since price deviation exceeds threshold
	reportPrice = answer + 1
	select {
	case <-submissionReceived:
		t.Fatalf("chainlink node updated FA, even though it's not allowed to")
	case <-time.After(500 * time.Millisecond):
	}
	// Could add a check for "not eligible to submit here", using the memory log
	newRound = big.NewInt(0).Add(newRound, big.NewInt(1))
	processedAnswer = big.NewInt(int64(100 * reportPrice))
	precision := job.Initiators[0].InitiatorParams.Precision
	// FORCE node to try to start a new round
	err = app.FluxMonitor.(*fluxmonitor.ConcreteFluxMonitor).
		XXXTestingOnlyCreateJob(t, j.ID,
			decimal.New(processedAnswer.Int64(), precision), newRound)
	require.NoError(t, err)
	select {
	case <-submissionReceived:
		t.Fatalf("FA allowed chainlink node to start a new round early")
	case <-time.After(500 * time.Millisecond):
	}
	// Try to start a new round directly, should fail
	_, err = fa.aggregatorContract.StartNewRound(fa.nallory)
	assert.Error(t, err, "FA allowed chainlink node to start a new round early")

	//- finally, ensure it can start a legitimate round after roundDelay is reached
	updateAnswer(t, &fa, newRound, processedAnswer, fa.ned, true, false)
	// Triggers a new round, since price deviation exceeds threshold
	reportPrice = answer + 1
	select {
	case <-submissionReceived:
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("could not start a new round, even though delay has passed")
	}
}

// XAU/XAG happened partly because you can update the entire state all at once.
// Having to add oracles one-by-one slows you down, so you can avoid some
// mistakes.
