package fluxmonitorv2_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/web"

	"github.com/ethereum/go-ethereum/eth/ethconfig"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flags_wrapper"
	faw "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitorv2"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
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
	key                       ethkey.Key
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
	nallory *bind.TransactOpts // Node operator Flux Monitor Oracle running this node
}

// newIdentity returns a go-ethereum abstraction of an ethereum account for
// interacting with contract golang wrappers
func newIdentity(t *testing.T) *bind.TransactOpts {
	key, err := crypto.GenerateKey()
	require.NoError(t, err, "failed to generate ethereum identity")
	return cltest.MustNewSimulatedBackendKeyedTransactor(t, key)
}

type fluxAggregatorUniverseConfig struct {
	MinSubmission *big.Int
	MaxSubmission *big.Int
}

func WithMinMaxSubmission(min, max *big.Int) func(cfg *fluxAggregatorUniverseConfig) {
	return func(cfg *fluxAggregatorUniverseConfig) {
		cfg.MinSubmission = min
		cfg.MaxSubmission = max
	}
}

// setupFluxAggregatorUniverse returns a fully initialized fluxAggregator universe. The
// arguments match the arguments of the same name in the FluxAggregator
// constructor.
func setupFluxAggregatorUniverse(t *testing.T, configOptions ...func(cfg *fluxAggregatorUniverseConfig)) fluxAggregatorUniverse {
	cfg := &fluxAggregatorUniverseConfig{
		MinSubmission: big.NewInt(0),
		MaxSubmission: big.NewInt(100000000000),
	}

	for _, optFn := range configOptions {
		optFn(cfg)
	}

	key := cltest.MustGenerateRandomKey(t)
	k, err := keystore.DecryptKey(key.JSON.RawMessage[:], cltest.Password)
	require.NoError(t, err)
	oracleTransactor := cltest.MustNewSimulatedBackendKeyedTransactor(t, k.PrivateKey)

	var f fluxAggregatorUniverse
	f.key = key
	f.sergey = newIdentity(t)
	f.neil = newIdentity(t)
	f.ned = newIdentity(t)
	f.nallory = oracleTransactor
	genesisData := core.GenesisAlloc{
		f.sergey.From:  {Balance: oneEth},
		f.neil.From:    {Balance: oneEth},
		f.ned.From:     {Balance: oneEth},
		f.nallory.From: {Balance: oneEth},
	}
	gasLimit := ethconfig.Defaults.Miner.GasCeil * 2
	f.backend = backends.NewSimulatedBackend(genesisData, gasLimit)

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
	f.aggregatorContractAddress, _, f.aggregatorContract, err = faw.DeployFluxAggregator(
		f.sergey,
		f.backend,
		linkAddress,
		big.NewInt(fee),
		faTimeout,
		common.Address{},
		cfg.MinSubmission,
		cfg.MaxSubmission,
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

// watchSubmissionReceived creates a channel which sends the log when a
// submission is received. When event appears on submissionReceived,
// it indicates that flux monitor job run is complete.
//
// It will only watch for logs from addresses that are provided
func (fau fluxAggregatorUniverse) WatchSubmissionReceived(t *testing.T, addresses []common.Address) chan *faw.FluxAggregatorSubmissionReceived {
	submissionReceived := make(chan *faw.FluxAggregatorSubmissionReceived)
	subscription, err := fau.aggregatorContract.WatchSubmissionReceived(
		nil,
		submissionReceived,
		[]*big.Int{},
		[]uint32{},
		addresses,
	)
	require.NoError(t, err, "failed to subscribe to SubmissionReceived events")
	t.Cleanup(subscription.Unsubscribe)

	return submissionReceived
}

func setupApplication(
	t *testing.T,
	fa fluxAggregatorUniverse,
	setConfig func(cfg *orm.Config),
) *cltest.TestApplication {
	config, cfgCleanup := cltest.NewConfig(t)
	setConfig(config.Config)

	t.Cleanup(cfgCleanup)

	app, cleanup := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, fa.backend, fa.key)
	t.Cleanup(cleanup)

	return app
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

func generatePriceResponseFn(price *int64) func() string {
	return func() string {
		return fmt.Sprintf(`{"data":{"result": %d}}`, atomic.LoadInt64(price))
	}
}

type answerParams struct {
	fa                          *fluxAggregatorUniverse
	roundId, answer             int64
	from                        *bind.TransactOpts
	isNewRound, completesAnswer bool
}

// checkSubmission verifies all the logs emitted by fa's FluxAggregator
// contract after an updateAnswer with the given values.
func checkSubmission(t *testing.T, p answerParams, currentBalance int64, receiptBlock uint64) {
	t.Helper()
	if receiptBlock == 0 {
		receiptBlock = p.fa.backend.Blockchain().CurrentBlock().Number().Uint64()
	}
	blockRange := &bind.FilterOpts{Start: 0, End: &receiptBlock}

	// Could filter for the known values here, but while that would be more
	// succinct it leads to less informative error messages... Did the log not
	// appear at all, or did it just have a wrong value?
	ilogs, err := p.fa.aggregatorContract.FilterSubmissionReceived(
		blockRange,
		[]*big.Int{big.NewInt(p.answer)},
		[]uint32{uint32(p.roundId)},
		[]common.Address{p.from.From},
	)
	require.NoError(t, err, "failed to get SubmissionReceived logs")

	var srlogs []*faw.FluxAggregatorSubmissionReceived
	_ = cltest.GetLogs(t, &srlogs, ilogs)
	require.Len(t, srlogs, 1, "FluxAggregator did not emit correct "+
		"SubmissionReceived log")

	inrlogs, err := p.fa.aggregatorContract.FilterNewRound(
		blockRange, []*big.Int{big.NewInt(p.roundId)}, []common.Address{p.from.From},
	)
	require.NoError(t, err, "failed to get NewRound logs")

	if p.isNewRound {
		var nrlogs []*faw.FluxAggregatorNewRound
		cltest.GetLogs(t, &nrlogs, inrlogs)
		require.Len(t, nrlogs, 1, "FluxAggregator did not emit correct NewRound "+
			"log")
	} else {
		assert.Len(t, cltest.GetLogs(t, nil, inrlogs), 0, "FluxAggregator emitted "+
			"unexpected NewRound log")
	}

	iaflogs, err := p.fa.aggregatorContract.FilterAvailableFundsUpdated(
		blockRange, []*big.Int{big.NewInt(currentBalance - fee)},
	)
	require.NoError(t, err, "failed to get AvailableFundsUpdated logs")
	var aflogs []*faw.FluxAggregatorAvailableFundsUpdated
	_ = cltest.GetLogs(t, &aflogs, iaflogs)
	assert.Len(t, aflogs, 1, "FluxAggregator did not emit correct "+
		"AvailableFundsUpdated log")

	iaulogs, err := p.fa.aggregatorContract.FilterAnswerUpdated(blockRange,
		[]*big.Int{big.NewInt(p.answer)}, []*big.Int{big.NewInt(p.roundId)},
	)
	require.NoError(t, err, "failed to get AnswerUpdated logs")
	if p.completesAnswer {
		var aulogs []*faw.FluxAggregatorAnswerUpdated
		_ = cltest.GetLogs(t, &aulogs, iaulogs)
		// XXX: sometimes this log is repeated; don't know why...
		assert.NotEmpty(t, aulogs, "FluxAggregator did not emit correct "+
			"AnswerUpdated log")
	}
}

// currentbalance returns the current balance of fa's FluxAggregator
func currentBalance(t *testing.T, fa *fluxAggregatorUniverse) *big.Int {
	currentBalance, err := fa.aggregatorContract.AvailableFunds(nil)
	require.NoError(t, err, "failed to get current FA balance")
	return currentBalance
}

// submitAnswer simulates a call to fa's FluxAggregator contract from a fake
// node (neil or ned), with the given roundId and answer, and checks that all
// the logs emitted by the contract are correct
func submitAnswer(t *testing.T, p answerParams) {
	cb := currentBalance(t, p.fa)

	// used to ensure that the simulated backend has processed the submission,
	// before we search for the log and checek it.
	srCh := make(chan *faw.FluxAggregatorSubmissionReceived)
	fromBlock := uint64(0)
	srSubscription, err := p.fa.aggregatorContract.WatchSubmissionReceived(
		&bind.WatchOpts{Start: &fromBlock},
		srCh,
		[]*big.Int{big.NewInt(p.answer)},
		[]uint32{uint32(p.roundId)},
		[]common.Address{p.from.From},
	)
	defer func() {
		srSubscription.Unsubscribe()
		err = <-srSubscription.Err()
		require.NoError(t, err, "failed to unsubscribe from AvailableFundsUpdated logs")
	}()

	_, err = p.fa.aggregatorContract.Submit(
		p.from, big.NewInt(p.roundId), big.NewInt(p.answer),
	)
	require.NoError(t, err, "failed to submit answer to flux aggregator")

	p.fa.backend.Commit()

	select {
	case <-srCh:
	case <-time.After(5 * time.Second):
		t.Fatalf("failed to complete submission to flux aggregator")
	}
	checkSubmission(t, p, cb.Int64(), 0)
}

func awaitSubmission(t *testing.T, submissionReceived chan *faw.FluxAggregatorSubmissionReceived) (
	receiptBlock uint64, answer int64,
) {
	select { // block until FluxAggregator contract acknowledges chainlink message
	case log := <-submissionReceived:
		return log.Raw.BlockNumber, log.Submission.Int64()
	case <-time.After(20 * pollTimerPeriod):
		t.Fatalf("chainlink failed to submit answer to FluxAggregator contract")
		return 0, 0 // unreachable
	}
}

// assertNoSubmission asserts that no submission was sent for a given duration
func assertNoSubmission(t *testing.T,
	submissionReceived chan *faw.FluxAggregatorSubmissionReceived,
	duration time.Duration,
	msg string,
) {
	select {
	case <-submissionReceived:
		assert.Fail(t, "flags are up, but submission was sent")
	case <-time.After(2 * time.Second):
	}
}

// assertPipelineRunCreated checks that a pipeline exists for a given round and
// verifies the answer
func assertPipelineRunCreated(t *testing.T, db *gorm.DB, roundID int64, result float64) {
	// Fetch the stats to extract the run id
	stats := fluxmonitorv2.FluxMonitorRoundStatsV2{}
	require.NoError(t, db.Where("round_id = ?", roundID).Find(&stats).Error)

	// Verify the pipeline run data
	run := pipeline.Run{}
	require.NoError(t, db.Find(&run, stats.PipelineRunID).Error, "runID %v", stats.PipelineRunID)
	assert.Equal(t, []interface{}{result}, run.Outputs.Val)
}

func TestFluxMonitor_Deviation(t *testing.T) {
	fa := setupFluxAggregatorUniverse(t)

	// - add oracles
	oracleList := []common.Address{fa.nallory.From}
	_, err := fa.aggregatorContract.ChangeOracles(fa.sergey, emptyList, oracleList, oracleList, 1, 1, 0)
	assert.NoError(t, err, "failed to add oracles to aggregator")
	fa.backend.Commit()
	checkOraclesAdded(t, fa, oracleList)

	// Set up chainlink app
	app := setupApplication(t, fa, func(cfg *orm.Config) {
		cfg.Set("DEFAULT_HTTP_TIMEOUT", "100ms")
		cfg.Set("TRIGGER_FALLBACK_DB_POLL_INTERVAL", "1s")
	})
	require.NoError(t, app.StartAndConnect())

	// Create mock server
	// We expect metadata of:
	//  latestAnswer:nil updatedAt:nil // First call
	//  latestAnswer:100 updatedAt:50
	//  latestAnswer:103 updatedAt:60
	type k struct{ latestAnswer, updatedAt string }
	expectedMeta := map[k]struct{}{
		{"100", "50"}: {},
		{"103", "60"}: {},
	}
	reportPrice := int64(100)
	mockServer := cltest.NewHTTPMockServerWithAlterableResponseAndRequest(t,
		generatePriceResponseFn(&reportPrice),
		func(r *http.Request) {
			b, err1 := ioutil.ReadAll(r.Body)
			require.NoError(t, err1)
			var m models.BridgeMetaDataJSON
			require.NoError(t, json.Unmarshal(b, &m))
			if m.Meta.LatestAnswer != nil && m.Meta.UpdatedAt != nil {
				delete(expectedMeta, k{m.Meta.LatestAnswer.String(), m.Meta.UpdatedAt.String()})
			}
		},
	)
	t.Cleanup(mockServer.Close)
	u, _ := url.Parse(mockServer.URL)
	app.Store.CreateBridgeType(&models.BridgeType{
		Name: "bridge",
		URL:  models.WebURL(*u),
	})

	// When event appears on submissionReceived, flux monitor job run is complete
	submissionReceived := fa.WatchSubmissionReceived(t,
		[]common.Address{fa.nallory.From},
	)

	// Create the job
	s := `
	type              = "fluxmonitor"
	schemaVersion     = 1
	name              = "integration test"
	contractAddress   = "%s"
	precision = 0
	threshold = 2.0
	absoluteThreshold = 0.0

	idleTimerPeriod = "1s"
	idleTimerDisabled = false

	pollTimerPeriod = "%s"
	pollTimerDisabled = false

	observationSource = """
	ds1 [type=bridge name=bridge];
	ds1_parse [type=jsonparse path="data,result"];

	ds1 -> ds1_parse
	"""
		`

	s = fmt.Sprintf(s, fa.aggregatorContractAddress, pollTimerPeriod)

	requestBody, err := json.Marshal(web.CreateJobRequest{
		TOML: string(s),
	})
	assert.NoError(t, err)

	initialBalance := currentBalance(t, &fa).Int64()

	cltest.CreateJobViaWeb2(t, app, string(requestBody))

	// Initial Poll
	receiptBlock, answer := awaitSubmission(t, submissionReceived)

	assert.Equal(t, atomic.LoadInt64(&reportPrice), answer,
		"failed to report correct price to contract")
	checkSubmission(t,
		answerParams{
			fa:              &fa,
			roundId:         1,
			answer:          int64(100),
			from:            fa.nallory,
			isNewRound:      true,
			completesAnswer: true,
		},
		initialBalance,
		receiptBlock,
	)
	assertPipelineRunCreated(t, app.Store.DB, 1, float64(100))

	// Change reported price to a value outside the deviation
	reportPrice = int64(103)
	receiptBlock, answer = awaitSubmission(t, submissionReceived)

	assert.Equal(t, atomic.LoadInt64(&reportPrice), answer,
		"failed to report correct price to contract")
	checkSubmission(t,
		answerParams{
			fa:              &fa,
			roundId:         2,
			answer:          int64(103),
			from:            fa.nallory,
			isNewRound:      true,
			completesAnswer: true,
		},
		initialBalance-fee,
		receiptBlock,
	)
	assertPipelineRunCreated(t, app.Store.DB, 2, float64(103))

	// Should not received a submission as it is inside the deviation
	reportPrice = int64(104)
	assertNoSubmission(t, submissionReceived, 2*time.Second, "Should not receive a submission")
	assert.Len(t, expectedMeta, 0, "expected metadata %v", expectedMeta)
}

func TestFluxMonitor_NewRound(t *testing.T) {
	fa := setupFluxAggregatorUniverse(t)

	// - add oracles
	oracleList := []common.Address{fa.neil.From, fa.nallory.From}
	_, err := fa.aggregatorContract.ChangeOracles(fa.sergey, emptyList, oracleList, oracleList, 1, 2, 1)
	assert.NoError(t, err, "failed to add oracles to aggregator")
	fa.backend.Commit()
	checkOraclesAdded(t, fa, oracleList)

	// Set up chainlink app
	app := setupApplication(t, fa, func(cfg *orm.Config) {
		cfg.Set("DEFAULT_HTTP_TIMEOUT", "100ms")
		cfg.Set("FLAGS_CONTRACT_ADDRESS", fa.flagsContractAddress.Hex())
		cfg.Set("TRIGGER_FALLBACK_DB_POLL_INTERVAL", "1s")
	})
	require.NoError(t, app.StartAndConnect())

	initialBalance := currentBalance(t, &fa).Int64()

	// Create mock server
	reportPrice := int64(1)
	mockServer := cltest.NewHTTPMockServerWithAlterableResponse(t,
		generatePriceResponseFn(&reportPrice),
	)
	t.Cleanup(mockServer.Close)

	// When event appears on submissionReceived, flux monitor job run is complete
	submissionReceived := fa.WatchSubmissionReceived(t,
		[]common.Address{fa.nallory.From},
	)

	// Create the job
	s := `
type              = "fluxmonitor"
schemaVersion     = 1
name              = "example flux monitor spec"
contractAddress   = "%s"
precision = 2
threshold = 0.5
absoluteThreshold = 0.0

idleTimerPeriod = "1s"
idleTimerDisabled = false

pollTimerPeriod = "%s"
pollTimerDisabled = false

observationSource = """
ds1 [type=http method=GET url="%s"];
ds1_parse [type=jsonparse path="data,result"];

ds1 -> ds1_parse
"""
	`

	s = fmt.Sprintf(s, fa.aggregatorContractAddress, pollTimerPeriod, mockServer.URL)

	// raise flags to disable polling
	fa.flagsContract.RaiseFlag(fa.sergey, utils.ZeroAddress) // global kill switch
	fa.flagsContract.RaiseFlag(fa.sergey, fa.aggregatorContractAddress)
	fa.backend.Commit()

	requestBody, err := json.Marshal(web.CreateJobRequest{
		TOML: string(s),
	})
	assert.NoError(t, err)

	cltest.CreateJobViaWeb2(t, app, string(requestBody))

	// Have the the fake node start a new round
	submitAnswer(t, answerParams{
		fa:              &fa,
		roundId:         1,
		answer:          2,
		from:            fa.neil,
		isNewRound:      true,
		completesAnswer: false,
	})

	// Waiting for flux monitor to finish Register process in log broadcaster
	// and then to have log broadcaster backfill logs after the debounceResubscribe period of ~ 1 sec
	assert.Eventually(t, func() bool {
		return app.LogBroadcaster.TrackedAddressesCount() >= 2
	}, 3*time.Second, 200*time.Millisecond)

	// Finally, the logs from log broadcaster are sent only after a next block is received.
	fa.backend.Commit()

	// Wait for the node's submission, and ensure it submits to the round
	// started by the fake node
	receiptBlock, _ := awaitSubmission(t, submissionReceived)
	checkSubmission(t,
		answerParams{
			fa:              &fa,
			roundId:         1,
			answer:          int64(1),
			from:            fa.nallory,
			isNewRound:      false,
			completesAnswer: true,
		},
		initialBalance-fee,
		receiptBlock,
	)
}

func TestFluxMonitor_HibernationMode(t *testing.T) {
	fa := setupFluxAggregatorUniverse(t)

	// - add oracles
	oracleList := []common.Address{fa.nallory.From}
	_, err := fa.aggregatorContract.ChangeOracles(fa.sergey, emptyList, oracleList, oracleList, 1, 1, 0)
	assert.NoError(t, err, "failed to add oracles to aggregator")
	fa.backend.Commit()
	checkOraclesAdded(t, fa, oracleList)

	// Set up chainlink app
	app := setupApplication(t, fa, func(cfg *orm.Config) {
		cfg.Set("DEFAULT_HTTP_TIMEOUT", "100ms")
		cfg.Set("FLAGS_CONTRACT_ADDRESS", fa.flagsContractAddress.Hex())
		cfg.Set("TRIGGER_FALLBACK_DB_POLL_INTERVAL", "1s")
	})
	require.NoError(t, app.StartAndConnect())

	// Create mock server
	reportPrice := int64(1)
	mockServer := cltest.NewHTTPMockServerWithAlterableResponse(t,
		generatePriceResponseFn(&reportPrice),
	)
	t.Cleanup(mockServer.Close)

	// When event appears on submissionReceived, flux monitor job run is complete
	submissionReceived := fa.WatchSubmissionReceived(t,
		[]common.Address{fa.nallory.From},
	)

	// Create the job
	s := `
type              = "fluxmonitor"
schemaVersion     = 1
name              = "example flux monitor spec"
contractAddress   = "%s"
precision = 0
threshold = 0.5
absoluteThreshold = 0.0

idleTimerPeriod = "1s"
idleTimerDisabled = false

pollTimerPeriod = "%s"
pollTimerDisabled = false

observationSource = """
ds1 [type=http method=GET url="%s"];
ds1_parse [type=jsonparse path="data,result"];

ds1 -> ds1_parse
"""
	`

	s = fmt.Sprintf(s, fa.aggregatorContractAddress, "1000ms", mockServer.URL)

	// raise flags
	fa.flagsContract.RaiseFlag(fa.sergey, utils.ZeroAddress) // global kill switch
	fa.flagsContract.RaiseFlag(fa.sergey, fa.aggregatorContractAddress)
	fa.backend.Commit()

	requestBody, err := json.Marshal(web.CreateJobRequest{
		TOML: string(s),
	})
	assert.NoError(t, err)

	j := cltest.CreateJobViaWeb2(t, app, string(requestBody))

	// node doesn't submit initial response, because flag is up
	// Wait here so the next lower flags doesn't trigger immediately
	cltest.AssertPipelineRunsStays(t, j.PipelineSpec.ID, app.Store, 0)

	// lower global kill switch flag - should trigger job run
	fa.flagsContract.LowerFlags(fa.sergey, []common.Address{utils.ZeroAddress})
	fa.backend.Commit()
	awaitSubmission(t, submissionReceived)

	reportPrice = int64(2) // change in price should trigger run
	awaitSubmission(t, submissionReceived)

	// lower contract's flag - should have no effect
	fa.flagsContract.LowerFlags(fa.sergey, []common.Address{fa.aggregatorContractAddress})
	fa.backend.Commit()
	assertNoSubmission(t, submissionReceived, 5*pollTimerPeriod, "should not trigger a new run because FM is already hibernating")

	// change in price should trigger run
	reportPrice = int64(4)
	awaitSubmission(t, submissionReceived)

	// raise both flags
	fa.flagsContract.RaiseFlag(fa.sergey, fa.aggregatorContractAddress)
	fa.flagsContract.RaiseFlag(fa.sergey, utils.ZeroAddress)
	fa.backend.Commit()

	// wait for FM to receive flags raised logs
	assert.Eventually(t, func() bool {
		ilogs, err := fa.flagsContract.FilterFlagRaised(nil, []common.Address{})
		require.NoError(t, err)
		logs := cltest.GetLogs(t, nil, ilogs)
		return len(logs) == 4
	}, 7*time.Second, 100*time.Millisecond)

	// change in price should not trigger run
	reportPrice = int64(8)
	assertNoSubmission(t, submissionReceived, 5*pollTimerPeriod, "should not trigger a new run, while flag is raised")
}

func TestFluxMonitor_InvalidSubmission(t *testing.T) {
	// 8 decimals places used for prices.
	fa := setupFluxAggregatorUniverse(t, WithMinMaxSubmission(
		big.NewInt(100000000),     // 1 * 10^8
		big.NewInt(1000000000000), // 10000 * 10^8
	))

	oracleList := []common.Address{fa.neil.From, fa.ned.From, fa.nallory.From}
	_, err := fa.aggregatorContract.ChangeOracles(fa.sergey, emptyList, oracleList, oracleList, 1, 3, 2)
	assert.NoError(t, err, "failed to add oracles to aggregator")
	fa.backend.Commit()

	// Set up chainlink app
	app := setupApplication(t, fa, func(cfg *orm.Config) {
		cfg.Set("DEFAULT_HTTP_TIMEOUT", "100ms")
		cfg.Set("TRIGGER_FALLBACK_DB_POLL_INTERVAL", "1s")
		cfg.Set("MIN_OUTGOING_CONFIRMATIONS", "2")
		cfg.Set("MIN_OUTGOING_CONFIRMATIONS", "2")
		cfg.Set("ETH_HEAD_TRACKER_MAX_BUFFER_SIZE", "100")
	})
	require.NoError(t, app.StartAndConnect())

	// Report a price that is above the maximum allowed value,
	// causing it to revert.
	reportPrice := int64(10001) // 10001 ETH/USD price is outside the range.
	mockServer := cltest.NewHTTPMockServerWithAlterableResponse(t,
		generatePriceResponseFn(&reportPrice),
	)
	t.Cleanup(mockServer.Close)

	// Generate custom TOML for this test due to precision change
	toml := `
type              = "fluxmonitor"
schemaVersion     = 1
name              = "example flux monitor spec"
contractAddress   = "%s"
precision = %d
threshold = 0.5
absoluteThreshold = 0.01

idleTimerPeriod = "1h"
idleTimerDisabled = false

pollTimerPeriod = "%s"
pollTimerDisabled = false

observationSource = """
ds1 [type=http method=GET url="%s"];
ds1_parse [type=jsonparse path="data,result"];

ds1 -> ds1_parse
"""
`

	s := fmt.Sprintf(toml, fa.aggregatorContractAddress, 8, "100ms", mockServer.URL)

	// raise flags
	fa.flagsContract.RaiseFlag(fa.sergey, utils.ZeroAddress) // global kill switch
	fa.flagsContract.RaiseFlag(fa.sergey, fa.aggregatorContractAddress)
	fa.backend.Commit()

	requestBody, err := json.Marshal(web.CreateJobRequest{
		TOML: string(s),
	})
	assert.NoError(t, err)

	j := cltest.CreateJobViaWeb2(t, app, string(requestBody))

	closer := cltest.Mine(fa.backend, 500*time.Millisecond)
	defer closer()

	// We should see a spec error because the value is too large to submit on-chain.
	jobID, err := strconv.ParseInt(j.ID, 10, 32)
	require.NoError(t, err)

	jse := cltest.WaitForSpecErrorV2(t, app.Store, int32(jobID), 1)
	assert.Contains(t, jse[0].Description, "Answer is outside acceptable range")
}

func TestFluxMonitorAntiSpamLogic(t *testing.T) {
	// - deploy a brand new FM contract
	fa := setupFluxAggregatorUniverse(t)

	// - add oracles
	oracleList := []common.Address{fa.neil.From, fa.ned.From, fa.nallory.From}
	_, err := fa.aggregatorContract.ChangeOracles(fa.sergey, emptyList, oracleList, oracleList, 2, 3, 2)
	assert.NoError(t, err, "failed to add oracles to aggregator")
	fa.backend.Commit()
	checkOraclesAdded(t, fa, oracleList)

	// Set up chainlink app
	app := setupApplication(t, fa, func(cfg *orm.Config) {
		cfg.Set("DEFAULT_HTTP_TIMEOUT", "100ms")
		cfg.Set("TRIGGER_FALLBACK_DB_POLL_INTERVAL", "1s")
	})
	require.NoError(t, app.StartAndConnect())

	minFee := app.Store.Config.MinimumContractPayment().ToInt().Int64()
	require.Equal(t, fee, minFee, "fee paid by FluxAggregator (%d) must at "+
		"least match MinimumContractPayment (%s). (Which is currently set in "+
		"cltest.go.)", fee, minFee)

	answer := int64(1) // Answer the nodes give on the first round

	//- have one of the fake nodes start a round.
	roundId := int64(1)
	processedAnswer := answer * 100 /* job has multiply times 100 */
	submitAnswer(t, answerParams{
		fa:              &fa,
		roundId:         roundId,
		answer:          processedAnswer,
		from:            fa.neil,
		isNewRound:      true,
		completesAnswer: false,
	})

	// - successfully close the round through the submissions of the other nodes
	// Response by spammy chainlink node, nallory
	//
	// The initial balance is the LINK balance of flux aggregator contract. We
	// use it to check that the fee for submitting an answer has been paid out.
	initialBalance := currentBalance(t, &fa).Int64()
	reportPrice := answer
	priceResponse := func() string {
		return fmt.Sprintf(`{"data":{"result": %d}}`, atomic.LoadInt64(&reportPrice))
	}
	mockServer := cltest.NewHTTPMockServerWithAlterableResponse(t, priceResponse)
	t.Cleanup(mockServer.Close)

	// When event appears on submissionReceived, flux monitor job run is complete
	submissionReceived := fa.WatchSubmissionReceived(t,
		[]common.Address{fa.nallory.From},
	)

	// Create FM Job, and wait for job run to start (the above submitAnswr call
	// to FluxAggregator contract initiates a run.)
	s := `
type              = "fluxmonitor"
schemaVersion     = 1
name              = "example flux monitor spec"
contractAddress   = "%s"
precision = 2
threshold = 0.5
absoluteThreshold = 0.0

idleTimerPeriod = "1s"
idleTimerDisabled = false

pollTimerPeriod = "%s"
pollTimerDisabled = false

observationSource = """
ds1 [type=http method=GET url="%s"];
ds1_parse [type=jsonparse path="data,result"];
ds1_multiply [type=multiply times=100]

ds1 -> ds1_parse -> ds1_multiply
"""
	`

	s = fmt.Sprintf(s, fa.aggregatorContractAddress, "200ms", mockServer.URL)
	requestBody, err := json.Marshal(web.CreateJobRequest{
		TOML: string(s),
	})
	assert.NoError(t, err)

	cltest.CreateJobViaWeb2(t, app, string(requestBody))

	receiptBlock, answer := awaitSubmission(t, submissionReceived)

	assert.Equal(t, 100*atomic.LoadInt64(&reportPrice), answer,
		"failed to report correct price to contract")
	checkSubmission(t,
		answerParams{
			fa:              &fa,
			roundId:         roundId,
			answer:          processedAnswer,
			from:            fa.nallory,
			isNewRound:      false,
			completesAnswer: true},
		initialBalance,
		receiptBlock,
	)

	//- have the malicious node start the next round.
	nextRoundBalance := initialBalance - fee
	// Triggers a new round, since price deviation exceeds threshold
	atomic.StoreInt64(&reportPrice, answer+1)

	receiptBlock, _ = awaitSubmission(t, submissionReceived)
	newRound := roundId + 1
	processedAnswer = 100 * atomic.LoadInt64(&reportPrice)
	checkSubmission(t,
		answerParams{
			fa:              &fa,
			roundId:         newRound,
			answer:          processedAnswer,
			from:            fa.nallory,
			isNewRound:      true,
			completesAnswer: false},
		nextRoundBalance,
		receiptBlock,
	)

	// Successfully close the round through the submissions of the other nodes
	submitAnswer(t,
		answerParams{
			fa:              &fa,
			roundId:         newRound,
			answer:          processedAnswer,
			from:            fa.neil,
			isNewRound:      false,
			completesAnswer: true},
	)

	// Have the malicious node try to start another round. It should not pass as
	// restartDelay has not been reached.
	newRound = newRound + 1
	processedAnswer = 100 * atomic.LoadInt64(&reportPrice)

	submitMaliciousAnswer(t,
		answerParams{
			fa:              &fa,
			roundId:         newRound,
			answer:          processedAnswer,
			from:            fa.nallory,
			isNewRound:      true,
			completesAnswer: false},
	)

	assertNoSubmission(t, submissionReceived, 5*pollTimerPeriod, "FA allowed chainlink node to start a new round early")

	// Try to start a new round directly, should fail because of delay
	_, err = fa.aggregatorContract.RequestNewRound(fa.nallory)
	assert.Error(t, err, "FA allowed chainlink node to start a new round early")

	//- finally, ensure it can start a legitimate round after restartDelay is
	//reached start an intervening round
	submitAnswer(t, answerParams{fa: &fa, roundId: newRound,
		answer: processedAnswer, from: fa.ned, isNewRound: true,
		completesAnswer: false})
	submitAnswer(t, answerParams{fa: &fa, roundId: newRound,
		answer: processedAnswer, from: fa.neil, isNewRound: false,
		completesAnswer: true})

	// start a legitimate new round
	atomic.StoreInt64(&reportPrice, reportPrice+3)

	// Wait for the node's submission, and ensure it submits to the round
	// started by the fake node
	awaitSubmission(t, submissionReceived)
}

// submitMaliciousAnswer simulates a call to fa's FluxAggregator contract from
// nallory, with the given roundId and answer and errors
func submitMaliciousAnswer(t *testing.T, p answerParams) {
	_, err := p.fa.aggregatorContract.Submit(
		p.from, big.NewInt(p.roundId), big.NewInt(p.answer),
	)
	require.Error(t, err)

	p.fa.backend.Commit()
}
