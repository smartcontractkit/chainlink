package internal_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestIntegration_Scheduler(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	defer cleanup()
	eth := app.MockEthClient(cltest.Strict)
	eth.Register("eth_chainId", *cltest.Int(app.Store.Config.ChainID()))
	app.Start()

	j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/scheduler_job.json")

	cltest.WaitForRunsAtLeast(t, j, app.Store, 1)

	initr := j.Initiators[0]
	assert.Equal(t, models.InitiatorCron, initr.Type)
	assert.Equal(t, "* * * * *", string(initr.Schedule), "Wrong cron schedule saved")
}

func TestIntegration_HttpRequestWithHeaders(t *testing.T) {
	tickerHeaders := http.Header{
		"Key1": []string{"value"},
		"Key2": []string{"value", "value"},
	}
	tickerResponse := `{"high": "10744.00", "last": "10583.75", "timestamp": "1512156162", "bid": "10555.13", "vwap": "10097.98", "volume": "17861.33960013", "low": "9370.11", "ask": "10583.00", "open": "9927.29"}`
	mockServer, assertCalled := cltest.NewHTTPMockServer(t, 200, "GET", tickerResponse,
		func(header http.Header, _ string) {
			for key, values := range tickerHeaders {
				assert.Equal(t, values, header[key])
			}
		})
	defer assertCalled()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	config := app.Config
	eth := app.MockEthClient(cltest.Strict)

	newHeads := make(chan models.BlockHeader)
	attempt1Hash := common.HexToHash("0xb7862c896a6ba2711bccc0410184e46d793ea83b3e05470f1d359ea276d16bb5")
	sentAt := uint64(23456)
	confirmed := sentAt + config.EthGasBumpThreshold() + 1
	safe := confirmed + config.MinOutgoingConfirmations() - 1
	unconfirmedReceipt := models.TxReceipt{}
	confirmedReceipt := models.TxReceipt{
		Hash:        attempt1Hash,
		BlockNumber: cltest.Int(confirmed),
	}

	eth.Context("app.Start()", func(eth *cltest.EthMock) {
		eth.RegisterSubscription("newHeads", newHeads)
		eth.Register("eth_getTransactionCount", `0x0100`) // TxManager.ActivateAccount()
		eth.Register("eth_chainId", config.ChainID())
	})
	assert.NoError(t, app.Start())
	eth.EventuallyAllCalled(t)

	eth.Context("ethTx.Perform()#1 at block 23456", func(eth *cltest.EthMock) {
		eth.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))
		eth.Register("eth_sendRawTransaction", attempt1Hash) // Initial tx attempt sent
		eth.Register("eth_getTransactionReceipt", unconfirmedReceipt)
	})
	j := cltest.CreateHelloWorldJobViaWeb(t, app, mockServer.URL)
	jr := cltest.WaitForJobRunToPendConfirmations(t, app.Store, cltest.CreateJobRunViaWeb(t, app, j))
	eth.EventuallyAllCalled(t)
	cltest.WaitForTxAttemptCount(t, app.Store, 1)

	eth.Context("ethTx.Perform()#4 at block 23465", func(eth *cltest.EthMock) {
		eth.Register("eth_blockNumber", utils.Uint64ToHex(safe))
		eth.Register("eth_getTransactionReceipt", confirmedReceipt) // confirmed for gas bumped txat
		eth.Register("eth_getBalance", "0x0100")
		eth.Register("eth_call", "0x0100")
	})
	newHeads <- models.BlockHeader{Number: cltest.BigHexInt(safe)} // 23465
	eth.EventuallyAllCalled(t)
	cltest.WaitForTxAttemptCount(t, app.Store, 1)

	cltest.WaitForJobRunToComplete(t, app.Store, jr)

	eth.EventuallyAllCalled(t)
}

func TestIntegration_FeeBump(t *testing.T) {
	tickerResponse := `{"high": "10744.00", "last": "10583.75", "timestamp": "1512156162", "bid": "10555.13", "vwap": "10097.98", "volume": "17861.33960013", "low": "9370.11", "ask": "10583.00", "open": "9927.29"}`
	mockServer, assertCalled := cltest.NewHTTPMockServer(t, 200, "GET", tickerResponse)
	defer assertCalled()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	config := app.Config

	// Put some distance between these two values so we can explore more of the state space
	config.Set("ETH_GAS_BUMP_THRESHOLD", 10)
	config.Set("MIN_OUTGOING_CONFIRMATIONS", 20)

	attempt1Hash := common.HexToHash("0xb7862c896a6ba2711bccc0410184e46d793ea83b3e05470f1d359ea276d16bb5")
	attempt2Hash := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000002")
	attempt3Hash := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000003")

	unconfirmedReceipt := models.TxReceipt{}

	// Enumerate the different block heights at which various state changes
	// happen for the transaction attempts created during this test
	firstTxSentAt := uint64(23456)
	firstTxGasBumpAt := firstTxSentAt + config.EthGasBumpThreshold()
	firstTxRemainsUnconfirmedAt := firstTxGasBumpAt - 1

	secondTxSentAt := firstTxGasBumpAt
	secondTxGasBumpAt := secondTxSentAt + config.EthGasBumpThreshold()
	secondTxRemainsUnconfirmedAt := secondTxGasBumpAt - 1

	thirdTxSentAt := secondTxGasBumpAt
	thirdTxConfirmedAt := thirdTxSentAt + 1
	thirdTxConfirmedReceipt := models.TxReceipt{
		Hash:        attempt1Hash,
		BlockNumber: cltest.Int(thirdTxConfirmedAt),
	}
	thirdTxSafeAt := thirdTxSentAt + config.MinOutgoingConfirmations()

	newHeads := make(chan models.BlockHeader)
	eth := app.MockEthClient(cltest.Strict)
	eth.Context("app.Start()", func(eth *cltest.EthMock) {
		eth.RegisterSubscription("newHeads", newHeads)
		eth.Register("eth_chainId", config.ChainID())
		eth.Register("eth_getTransactionCount", `0x0100`) // TxManager.ActivateAccount()
	})
	assert.NoError(t, app.Start())
	eth.EventuallyAllCalled(t)

	// This first run of the EthTx adapter creates an initial transaction which
	// starts unconfirmed
	eth.Context("ethTx.Perform()#1", func(eth *cltest.EthMock) {
		eth.Register("eth_blockNumber", utils.Uint64ToHex(firstTxSentAt))
		eth.Register("eth_sendRawTransaction", attempt1Hash)
		eth.Register("eth_getTransactionReceipt", unconfirmedReceipt)
	})
	j := cltest.CreateHelloWorldJobViaWeb(t, app, mockServer.URL)
	jr := cltest.WaitForJobRunToPendConfirmations(t, app.Store, cltest.CreateJobRunViaWeb(t, app, j))
	eth.EventuallyAllCalled(t)
	cltest.WaitForTxAttemptCount(t, app.Store, 1)

	// At the next head, the transaction is still unconfirmed, but no thresholds
	// have been met so we just wait...
	eth.Context("ethTx.Perform()#2", func(eth *cltest.EthMock) {
		eth.Register("eth_blockNumber", utils.Uint64ToHex(firstTxRemainsUnconfirmedAt))
		eth.Register("eth_getTransactionReceipt", unconfirmedReceipt)
	})
	newHeads <- models.BlockHeader{Number: cltest.BigHexInt(firstTxRemainsUnconfirmedAt)}
	eth.EventuallyAllCalled(t)
	jr = cltest.WaitForJobRunToPendConfirmations(t, app.Store, jr)

	// At the next head, the transaction remains unconfirmed but the gas bump
	// threshold has been met, so a new transaction is made with a higher amount
	// of gas ("bumped gas")
	eth.Context("ethTx.Perform()#3", func(eth *cltest.EthMock) {
		eth.Register("eth_blockNumber", utils.Uint64ToHex(firstTxGasBumpAt))
		eth.Register("eth_getTransactionReceipt", unconfirmedReceipt)
		eth.Register("eth_sendRawTransaction", attempt2Hash)
	})
	newHeads <- models.BlockHeader{Number: cltest.BigHexInt(firstTxGasBumpAt)}
	eth.EventuallyAllCalled(t)
	jr = cltest.WaitForJobRunToPendConfirmations(t, app.Store, jr)
	cltest.WaitForTxAttemptCount(t, app.Store, 2)

	// Another head comes in and both transactions are still unconfirmed, more
	// waiting...
	eth.Context("ethTx.Perform()#4", func(eth *cltest.EthMock) {
		eth.Register("eth_blockNumber", utils.Uint64ToHex(secondTxRemainsUnconfirmedAt))
		eth.Register("eth_getTransactionReceipt", unconfirmedReceipt)
		eth.Register("eth_getTransactionReceipt", unconfirmedReceipt)
	})
	newHeads <- models.BlockHeader{Number: cltest.BigHexInt(secondTxRemainsUnconfirmedAt)}
	eth.EventuallyAllCalled(t)
	jr = cltest.WaitForJobRunToPendConfirmations(t, app.Store, jr)

	// Now the second transaction attempt meets the gas bump threshold, so a
	// final transaction attempt shoud be made
	eth.Context("ethTx.Perform()#5", func(eth *cltest.EthMock) {
		eth.Register("eth_blockNumber", utils.Uint64ToHex(secondTxGasBumpAt))
		eth.Register("eth_getTransactionReceipt", unconfirmedReceipt)
		eth.Register("eth_getTransactionReceipt", unconfirmedReceipt)
		eth.Register("eth_sendRawTransaction", attempt3Hash)
	})
	newHeads <- models.BlockHeader{Number: cltest.BigHexInt(secondTxGasBumpAt)}
	eth.EventuallyAllCalled(t)
	jr = cltest.WaitForJobRunToPendConfirmations(t, app.Store, jr)
	cltest.WaitForTxAttemptCount(t, app.Store, 3)

	// This third attempt has enough gas and gets confirmed, but has not yet
	// received sufficient confirmations, so we wait again...
	eth.Context("ethTx.Perform()#6", func(eth *cltest.EthMock) {
		eth.Register("eth_blockNumber", utils.Uint64ToHex(thirdTxConfirmedAt))
		eth.Register("eth_getTransactionReceipt", unconfirmedReceipt)
		eth.Register("eth_getTransactionReceipt", unconfirmedReceipt)
		eth.Register("eth_getTransactionReceipt", thirdTxConfirmedReceipt)
	})
	newHeads <- models.BlockHeader{Number: cltest.BigHexInt(thirdTxConfirmedAt)}
	eth.EventuallyAllCalled(t)
	jr = cltest.WaitForJobRunToPendConfirmations(t, app.Store, jr)

	// Finally the third attempt gets to a minimum number of safe confirmations,
	// the amount remaining in the account is printed (eth_getBalance, eth_call)
	eth.Context("ethTx.Perform()#7", func(eth *cltest.EthMock) {
		eth.Register("eth_blockNumber", utils.Uint64ToHex(thirdTxSafeAt))
		eth.Register("eth_getTransactionReceipt", unconfirmedReceipt)
		eth.Register("eth_getTransactionReceipt", unconfirmedReceipt)
		eth.Register("eth_getTransactionReceipt", thirdTxConfirmedReceipt)
		eth.Register("eth_getBalance", "0x100")
		eth.Register("eth_call", "0x100")
	})
	newHeads <- models.BlockHeader{Number: cltest.BigHexInt(thirdTxSafeAt)}
	eth.EventuallyAllCalled(t)
	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)

	val, err := jr.TaskRuns[0].Result.ResultString()
	assert.NoError(t, err)
	assert.Equal(t, tickerResponse, val)
	val, err = jr.TaskRuns[1].Result.ResultString()
	assert.Equal(t, "10583.75", val)
	assert.NoError(t, err)
	val, err = jr.TaskRuns[3].Result.ResultString()
	assert.Equal(t, attempt1Hash.String(), val)
	assert.NoError(t, err)
	val, err = jr.Result.ResultString()
	assert.Equal(t, attempt1Hash.String(), val)
	assert.NoError(t, err)
}

func TestIntegration_RunAt(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t)
	defer cleanup()
	eth := app.MockEthClient(cltest.Strict)
	eth.Register("eth_chainId", app.Store.Config.ChainID())
	app.InstantClock()

	j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/run_at_job.json")

	initr := j.Initiators[0]
	assert.Equal(t, models.InitiatorRunAt, initr.Type)
	assert.Equal(t, "2018-01-08T18:12:01Z", utils.ISO8601UTC(initr.Time.Time))

	app.Start()

	cltest.WaitForRuns(t, j, app.Store, 1)
}

func TestIntegration_EthLog(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t)
	defer cleanup()

	eth := app.MockEthClient()
	logs := make(chan models.Log, 1)
	eth.Context("app.Start()", func(eth *cltest.EthMock) {
		eth.Register("eth_chainId", app.Store.Config.ChainID())
		eth.RegisterSubscription("logs", logs)
	})
	app.Start()

	j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/eth_log_job.json")
	address := common.HexToAddress("0x3cCad4715152693fE3BC4460591e3D3Fbd071b42")

	initr := j.Initiators[0]
	assert.Equal(t, models.InitiatorEthLog, initr.Type)
	assert.Equal(t, address, initr.Address)

	logs <- cltest.LogFromFixture(t, "testdata/requestLog0original.json")
	cltest.WaitForRuns(t, j, app.Store, 1)
}

func TestIntegration_RunLog(t *testing.T) {
	triggeringBlockHash := cltest.NewHash()
	otherBlockHash := cltest.NewHash()

	tests := []struct {
		name             string
		logBlockHash     common.Hash
		receiptBlockHash common.Hash
		wantStatus       models.RunStatus
		wantFinishedAt   bool
	}{
		{
			name:             "completed",
			logBlockHash:     triggeringBlockHash,
			receiptBlockHash: triggeringBlockHash,
			wantStatus:       models.RunStatusCompleted,
			wantFinishedAt:   true,
		},
		{
			name:             "ommered request",
			logBlockHash:     triggeringBlockHash,
			receiptBlockHash: otherBlockHash,
			wantStatus:       models.RunStatusErrored,
			wantFinishedAt:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, cfgCleanup := cltest.NewConfig(t)
			defer cfgCleanup()
			config.Set("MIN_INCOMING_CONFIRMATIONS", 6)
			app, cleanup := cltest.NewApplicationWithConfig(t, config)
			defer cleanup()

			eth := app.MockEthClient()
			logs := make(chan models.Log, 1)
			newHeads := eth.RegisterNewHeads()
			eth.Context("app.Start()", func(eth *cltest.EthMock) {
				eth.RegisterSubscription("logs", logs)
			})
			eth.Register("eth_chainId", config.ChainID())
			app.Start()

			j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/runlog_noop_job.json")
			requiredConfs := uint32(100)

			initr := j.Initiators[0]
			assert.Equal(t, models.InitiatorRunLog, initr.Type)

			creationHeight := uint32(1)
			runlog := cltest.NewRunLog(t, j.ID, cltest.NewAddress(), cltest.NewAddress(), int(creationHeight), `{}`)
			runlog.BlockHash = test.logBlockHash
			logs <- runlog
			cltest.WaitForRuns(t, j, app.Store, 1)

			runs, err := app.Store.JobRunsFor(j.ID)
			assert.NoError(t, err)
			jr := runs[0]
			cltest.WaitForJobRunToPendConfirmations(t, app.Store, jr)
			assert.False(t, jr.TaskRuns[0].Confirmations.Valid)

			blockIncrease := app.Store.Config.MinIncomingConfirmations()
			minGlobalHeight := creationHeight + blockIncrease
			newHeads <- models.BlockHeader{Number: cltest.BigHexInt(minGlobalHeight)}
			<-time.After(time.Second)
			jr = cltest.JobRunStaysPendingConfirmations(t, app.Store, jr)
			assert.Equal(t, uint32(creationHeight+blockIncrease), jr.TaskRuns[0].Confirmations.Uint32)

			safeNumber := creationHeight + requiredConfs
			newHeads <- models.BlockHeader{Number: cltest.BigHexInt(safeNumber)}
			confirmedReceipt := models.TxReceipt{
				Hash:        runlog.TxHash,
				BlockHash:   test.receiptBlockHash,
				BlockNumber: cltest.Int(creationHeight),
			}
			eth.Context("validateOnMainChain", func(ethMock *cltest.EthMock) {
				eth.Register("eth_getTransactionReceipt", confirmedReceipt)
			})

			jr = cltest.WaitForJobRunStatus(t, app.Store, jr, test.wantStatus)
			assert.Equal(t, test.wantFinishedAt, jr.FinishedAt.Valid)
			assert.Equal(t, requiredConfs, jr.TaskRuns[0].Confirmations.Uint32)
			assert.True(t, eth.AllCalled(), eth.Remaining())
		})
	}
}

func TestIntegration_EndAt(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	defer cleanup()
	eth := app.MockEthClient(cltest.Strict)
<<<<<<< HEAD
	eth.Register("eth_chainId", app.Store.Config.ChainID())
=======
	eth.Register("eth_chainId", *cltest.Int(app.Store.Config.ChainID()))
>>>>>>> Mock eth_chainId in integration tests
	clock := cltest.UseSettableClock(app.Store)
	app.Start()
	client := app.NewHTTPClient()

	j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/end_at_job.json")
	endAt := cltest.ParseISO8601(t, "3000-01-01T00:00:00.000Z")
	assert.Equal(t, endAt, j.EndAt.Time)

	cltest.CreateJobRunViaWeb(t, app, j)

	clock.SetTime(endAt.Add(time.Nanosecond))

	resp, cleanup := client.Post("/v2/specs/"+j.ID.String()+"/runs", &bytes.Buffer{})
	defer cleanup()
	assert.Equal(t, 500, resp.StatusCode)
	gomega.NewGomegaWithT(t).Consistently(func() []models.JobRun {
		jobRuns, err := app.Store.JobRunsFor(j.ID)
		assert.NoError(t, err)
		return jobRuns
	}).Should(gomega.HaveLen(1))
}

func TestIntegration_StartAt(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	defer cleanup()
	eth := app.MockEthClient(cltest.Strict)
<<<<<<< HEAD
	eth.Register("eth_chainId", app.Store.Config.ChainID())
=======
	eth.Register("eth_chainId", *cltest.Int(app.Store.Config.ChainID()))
>>>>>>> Mock eth_chainId in integration tests
	clock := cltest.UseSettableClock(app.Store)
	app.Start()
	client := app.NewHTTPClient()

	j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/start_at_job.json")
	startAt := cltest.ParseISO8601(t, "3000-01-01T00:00:00.000Z")
	assert.Equal(t, startAt, j.StartAt.Time)

	resp, cleanup := client.Post("/v2/specs/"+j.ID.String()+"/runs", &bytes.Buffer{})
	defer cleanup()
	assert.Equal(t, 500, resp.StatusCode)
	cltest.WaitForRuns(t, j, app.Store, 0)

	clock.SetTime(startAt)

	cltest.CreateJobRunViaWeb(t, app, j)
}

func TestIntegration_ExternalAdapter_RunLogInitiated(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	defer cleanup()

	eth := app.MockEthClient()
	eth.Register("eth_chainId", app.Store.Config.ChainID())
	logs := make(chan models.Log, 1)
	newHeads := make(chan models.BlockHeader, 10)
	eth.Context("app.Start()", func(eth *cltest.EthMock) {
		eth.RegisterSubscription("logs", logs)
		eth.RegisterSubscription("newHeads", newHeads)
	})
	app.Start()

	eaValue := "87698118359"
	eaExtra := "other values to be used by external adapters"
	eaResponse := fmt.Sprintf(`{"data":{"result": "%v", "extra": "%v"}}`, eaValue, eaExtra)
	mockServer, ensureRequest := cltest.NewHTTPMockServer(t, 200, "POST", eaResponse)
	defer ensureRequest()

	bridgeJSON := fmt.Sprintf(`{"name":"randomNumber","url":"%v","confirmations":10}`, mockServer.URL)
	cltest.CreateBridgeTypeViaWeb(t, app, bridgeJSON)
	j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/log_initiated_bridge_type_job.json")

	logBlockNumber := 1
	runlog := cltest.NewRunLog(t, j.ID, cltest.NewAddress(), cltest.NewAddress(), logBlockNumber, `{}`)
	logs <- runlog
	jr := cltest.WaitForRuns(t, j, app.Store, 1)[0]
	cltest.WaitForJobRunToPendConfirmations(t, app.Store, jr)

	newHeads <- models.BlockHeader{Number: cltest.BigHexInt(logBlockNumber + 8)}
	cltest.WaitForJobRunToPendConfirmations(t, app.Store, jr)

	confirmedReceipt := models.TxReceipt{
		Hash:        runlog.TxHash,
		BlockHash:   runlog.BlockHash,
		BlockNumber: cltest.Int(logBlockNumber),
	}
	eth.Context("validateOnMainChain", func(ethMock *cltest.EthMock) {
		eth.Register("eth_getTransactionReceipt", confirmedReceipt)
	})

	newHeads <- models.BlockHeader{Number: cltest.BigHexInt(logBlockNumber + 9)}
	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)

	tr := jr.TaskRuns[0]
	assert.Equal(t, "randomnumber", tr.TaskSpec.Type.String())
	val, err := tr.Result.ResultString()
	assert.NoError(t, err)
	assert.Equal(t, eaValue, val)
	res := tr.Result.Get("extra")
	assert.Equal(t, eaExtra, res.String())

	assert.True(t, eth.AllCalled(), eth.Remaining())
}

// This test ensures that the response body of an external adapter are supplied
// as params to the successive task
func TestIntegration_ExternalAdapter_Copy(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	defer cleanup()
	eth := app.MockEthClient(cltest.Strict)
<<<<<<< HEAD
	eth.Register("eth_chainId", app.Store.Config.ChainID())
=======
	eth.Register("eth_chainId", *cltest.Int(app.Store.Config.ChainID()))
>>>>>>> Mock eth_chainId in integration tests
	bridgeURL := cltest.WebURL(t, "https://test.chain.link/always")
	app.Store.Config.Set("BRIDGE_RESPONSE_URL", bridgeURL)
	app.Start()

	eaPrice := "1234"
	eaQuote := "USD"
	eaResponse := fmt.Sprintf(`{"data":{"price": "%v", "quote": "%v"}}`, eaPrice, eaQuote)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/", r.URL.Path)

		b, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)
		body := cltest.JSONFromBytes(t, b)
		data := body.Get("data")
		assert.True(t, data.Exists())
		bodyParam := data.Get("bodyParam")
		assert.True(t, bodyParam.Exists())
		assert.Equal(t, true, bodyParam.Bool())

		url := body.Get("responseURL")
		assert.Contains(t, url.String(), "https://test.chain.link/always/v2/runs")

		w.WriteHeader(200)
		io.WriteString(w, eaResponse)
	}))
	defer ts.Close()

	bridgeJSON := fmt.Sprintf(`{"name":"assetPrice","url":"%v"}`, ts.URL)
	cltest.CreateBridgeTypeViaWeb(t, app, bridgeJSON)
	j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/bridge_type_copy_job.json")
	jr := cltest.WaitForJobRunToComplete(t, app.Store, cltest.CreateJobRunViaWeb(t, app, j, `{"copyPath": ["price"]}`))

	tr := jr.TaskRuns[0]
	assert.Equal(t, "assetprice", tr.TaskSpec.Type.String())
	tr = jr.TaskRuns[1]
	assert.Equal(t, "copy", tr.TaskSpec.Type.String())
	val, err := tr.Result.ResultString()
	assert.NoError(t, err)
	assert.Equal(t, eaPrice, val)
	res := tr.Result.Get("quote")
	assert.Equal(t, eaQuote, res.String())
}

// This test ensures that an bridge adapter task is resumed from pending after
// sending out a request to an external adapter and waiting to receive a
// request back
func TestIntegration_ExternalAdapter_Pending(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	defer cleanup()
	eth := app.MockEthClient(cltest.Strict)
<<<<<<< HEAD
	eth.Register("eth_chainId", app.Store.Config.ChainID())
=======
	eth.Register("eth_chainId", *cltest.Int(app.Store.Config.ChainID()))
>>>>>>> Mock eth_chainId in integration tests
	app.Start()

	bta := &models.BridgeTypeAuthentication{}
	var j models.JobSpec
	mockServer, cleanup := cltest.NewHTTPMockServer(t, 200, "POST", `{"pending":true}`,
		func(h http.Header, b string) {
			body := cltest.JSONFromString(t, b)

			jrs := cltest.WaitForRuns(t, j, app.Store, 1)
			jr := jrs[0]
			id := body.Get("id")
			assert.True(t, id.Exists())
			assert.Equal(t, jr.ID.String(), id.String())

			data := body.Get("data")
			assert.True(t, data.Exists())
			assert.Equal(t, data.Type, gjson.JSON)

			token := utils.StripBearer(h.Get("Authorization"))
			assert.Equal(t, bta.OutgoingToken, token)
		})
	defer cleanup()

	bridgeJSON := fmt.Sprintf(`{"name":"randomNumber","url":"%v"}`, mockServer.URL)
	bta = cltest.CreateBridgeTypeViaWeb(t, app, bridgeJSON)
	j = cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/random_number_bridge_type_job.json")
	jr := cltest.CreateJobRunViaWeb(t, app, j)
	jr = cltest.WaitForJobRunToPendBridge(t, app.Store, jr)

	tr := jr.TaskRuns[0]
	assert.Equal(t, models.RunStatusPendingBridge, tr.Status)
	val, err := tr.Result.ResultString()
	assert.Error(t, err)
	assert.Equal(t, "", val)

	jr = cltest.UpdateJobRunViaWeb(t, app, jr, bta, `{"data":{"result":"100"}}`)
	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)
	tr = jr.TaskRuns[0]
	assert.Equal(t, models.RunStatusCompleted, tr.Status)
	val, err = tr.Result.ResultString()
	assert.NoError(t, err)
	assert.Equal(t, "100", val)
}

func TestIntegration_WeiWatchers(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	defer cleanup()

	eth := app.MockEthClient()
	eth.RegisterNewHead(1)
	logs := make(chan models.Log, 1)
	eth.Context("app.Start()", func(eth *cltest.EthMock) {
<<<<<<< HEAD
		eth.Register("eth_chainId", app.Config.ChainID())
=======
		eth.Register("eth_chainId", *cltest.Int(app.Config.ChainID()))
>>>>>>> Mock eth_chainId in integration tests
		eth.RegisterSubscription("logs", logs)
	})

	log := cltest.LogFromFixture(t, "testdata/requestLog0original.json")
	mockServer, cleanup := cltest.NewHTTPMockServer(t, 200, "POST", `{"pending":true}`,
		func(_ http.Header, body string) {
			marshaledLog, err := json.Marshal(&log)
			assert.NoError(t, err)
			assert.JSONEq(t, string(marshaledLog), body)
		})
	defer cleanup()

	j := cltest.NewJobWithLogInitiator()
	post := cltest.NewTask(t, "httppost", fmt.Sprintf(`{"url":"%v"}`, mockServer.URL))
	tasks := []models.TaskSpec{post}
	j.Tasks = tasks
	j = cltest.CreateJobSpecViaWeb(t, app, j)

	app.Start()
	logs <- log

	jobRuns := cltest.WaitForRuns(t, j, app.Store, 1)
	cltest.WaitForJobRunToComplete(t, app.Store, jobRuns[0])
}

func TestIntegration_MultiplierInt256(t *testing.T) {
	app, cleanup := cltest.NewApplication(t)
	defer cleanup()
	eth := app.MockEthClient(cltest.Strict)
<<<<<<< HEAD
	eth.Register("eth_chainId", app.Store.Config.ChainID())
=======
	eth.Register("eth_chainId", *cltest.Int(app.Store.Config.ChainID()))
>>>>>>> Mock eth_chainId in integration tests
	app.Start()

	j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/int256_job.json")
	jr := cltest.CreateJobRunViaWeb(t, app, j, `{"result":"-10221.30"}`)
	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)

	val, err := jr.Result.ResultString()
	assert.NoError(t, err)
	assert.Equal(t, "0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0674e", val)
}

func TestIntegration_MultiplierUint256(t *testing.T) {
	app, cleanup := cltest.NewApplication(t)
	defer cleanup()
	eth := app.MockEthClient(cltest.Strict)
<<<<<<< HEAD
	eth.Register("eth_chainId", app.Store.Config.ChainID())
=======
	eth.Register("eth_chainId", *cltest.Int(app.Store.Config.ChainID()))
>>>>>>> Mock eth_chainId in integration tests
	app.Start()

	j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/uint256_job.json")
	jr := cltest.CreateJobRunViaWeb(t, app, j, `{"result":"10221.30"}`)
	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)

	val, err := jr.Result.ResultString()
	assert.NoError(t, err)
	assert.Equal(t, "0x00000000000000000000000000000000000000000000000000000000000f98b2", val)
}

func TestIntegration_NonceManagement_firstRunWithExistingTxs(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()

	config, configCleanup := cltest.NewConfig(t)
	defer configCleanup()

	j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/web_initiated_eth_tx_job.json")

	eth := app.MockEthClient()
	eth.Context("app.Start()", func(eth *cltest.EthMock) {
		eth.Register("eth_getTransactionCount", `0x100`) // activate account nonce
		eth.Register("eth_chainId", config.ChainID())
	})
	require.NoError(t, app.StartAndConnect())

	createCompletedJobRun := func(blockNumber uint64, expectedNonce uint64) {
		hash := common.HexToHash("0xb7862c896a6ba2711bccc0410184e46d793ea83b3e05470f1d359ea276d16bb5")
		confirmedBlockNumber := blockNumber - app.Store.Config.MinOutgoingConfirmations()

		eth.Context("ethTx.Perform()", func(eth *cltest.EthMock) {
			eth.Register("eth_blockNumber", utils.Uint64ToHex(blockNumber))
			eth.Register("eth_sendRawTransaction", hash)
			eth.Register("eth_getTransactionReceipt", models.TxReceipt{
				Hash:        hash,
				BlockNumber: cltest.Int(confirmedBlockNumber),
			})
		})

		jr := cltest.CreateJobRunViaWeb(t, app, j, `{"result":"0x11"}`)
		cltest.WaitForJobRunToComplete(t, app.Store, jr)

		attempt := cltest.GetLastTxAttempt(t, app.Store)
		tx, err := app.Store.FindTx(attempt.TxID)
		assert.NoError(t, err)
		assert.Equal(t, expectedNonce, tx.Nonce)
	}

	createCompletedJobRun(100, uint64(0x100))
	createCompletedJobRun(200, uint64(0x101))
}

func TestIntegration_CreateServiceAgreement(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()

	eth := app.MockEthClient()
	logs := make(chan models.Log, 1)
	eth.Context("app.Start()", func(eth *cltest.EthMock) {
		eth.RegisterSubscription("logs", logs)
		eth.Register("eth_getTransactionCount", `0x100`)
		eth.Register("eth_chainId", app.Store.Config.ChainID())
	})
	assert.NoError(t, app.StartAndConnect())
	sa := cltest.FixtureCreateServiceAgreementViaWeb(t, app, "fixtures/web/noop_agreement.json")

	assert.NotEqual(t, "", sa.ID)
	j := cltest.FindJob(t, app.Store, sa.JobSpecID)

	assert.Equal(t, assets.NewLink(1000000000000000000), sa.Encumbrance.Payment)
	assert.Equal(t, uint64(300), sa.Encumbrance.Expiration)

	assert.Equal(t, time.Unix(1571523439, 0).UTC(), sa.Encumbrance.EndAt.Time)
	assert.NotEqual(t, "", sa.ID)

	// Request execution of the job associated with this ServiceAgreement
	logs <- cltest.NewServiceAgreementExecutionLog(
		t,
		j.ID,
		cltest.NewAddress(),
		cltest.NewAddress(),
		1,
		`{}`)

	runs := cltest.WaitForRuns(t, j, app.Store, 1)
	cltest.WaitForJobRunToComplete(t, app.Store, runs[0])

	eth.EventuallyAllCalled(t)

}

func TestIntegration_SyncJobRuns(t *testing.T) {
	t.Parallel()
	wsserver, wsserverCleanup := cltest.NewEventWebSocketServer(t)
	defer wsserverCleanup()

	config, _ := cltest.NewConfig(t)
	config.Set("EXPLORER_URL", wsserver.URL.String())
	app, cleanup := cltest.NewApplicationWithConfig(t, config)
	defer cleanup()
	eth := app.MockEthClient(cltest.Strict)
<<<<<<< HEAD
	eth.Register("eth_chainId", config.ChainID())
=======
	eth.Register("eth_chainId", *cltest.Int(config.ChainID()))
>>>>>>> Mock eth_chainId in integration tests

	app.InstantClock()

	app.Store.StatsPusher.Period = 300 * time.Millisecond

	app.Start()
	j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/run_at_job.json")

	cltest.CallbackOrTimeout(t, "stats pusher connects", func() {
		<-wsserver.Connected
	}, 5*time.Second)

	var message string
	cltest.CallbackOrTimeout(t, "stats pusher sends", func() {
		message = <-wsserver.Received
	}, 5*time.Second)

	var run models.JobRun
	err := json.Unmarshal([]byte(message), &run)
	require.NoError(t, err)
	assert.Equal(t, j.ID, run.JobSpecID)
}

func TestIntegration_SleepAdapter(t *testing.T) {
	t.Parallel()

	sleepSeconds := 4
	app, cleanup := cltest.NewApplication(t)
	defer cleanup()
	eth := app.MockEthClient(cltest.Strict)
<<<<<<< HEAD
	eth.Register("eth_chainId", app.Store.Config.ChainID())
=======
	eth.Register("eth_chainId", *cltest.Int(app.Store.Config.ChainID()))
>>>>>>> Mock eth_chainId in integration tests
	app.Start()

	j := cltest.FixtureCreateJobViaWeb(t, app, "./testdata/sleep_job.json")

	runInput := fmt.Sprintf("{\"until\": \"%s\"}", time.Now().Local().Add(time.Second*time.Duration(sleepSeconds)))
	jr := cltest.CreateJobRunViaWeb(t, app, j, runInput)

	cltest.WaitForJobRunToPendSleep(t, app.Store, jr)
	cltest.JobRunStays(t, app.Store, jr, models.RunStatusPendingSleep, time.Second)
	cltest.WaitForJobRunToComplete(t, app.Store, jr)
}
