package internal_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/auth"
	ethpkg "github.com/smartcontractkit/chainlink/core/eth"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/models/vrfkey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestIntegration_Scheduler(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t, cltest.EthMockRegisterChainID)
	defer cleanup()
	require.NoError(t, app.Start())

	j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/scheduler_job.json")

	cltest.WaitForRunsAtLeast(t, j, app.Store, 1)

	initr := j.Initiators[0]
	assert.Equal(t, models.InitiatorCron, initr.Type)
	assert.Equal(t, "CRON_TZ=UTC * * * * * *", string(initr.Schedule), "Wrong cron schedule saved")
}

func TestIntegration_HttpRequestWithHeaders(t *testing.T) {
	tickerHeaders := http.Header{
		"Key1": []string{"value"},
		"Key2": []string{"value", "value"},
	}
	tickerResponse := `{"high": "10744.00", "last": "10583.75", "timestamp": "1512156162", "bid": "10555.13", "vwap": "10097.98", "volume": "17861.33960013", "low": "9370.11", "ask": "10583.00", "open": "9927.29"}`
	mockServer, assertCalled := cltest.NewHTTPMockServer(t, http.StatusOK, "GET", tickerResponse,
		func(header http.Header, _ string) {
			for key, values := range tickerHeaders {
				assert.Equal(t, values, header[key])
			}
		})
	defer assertCalled()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	config := app.Config
	eth := app.EthMock

	newHeads := make(chan ethpkg.BlockHeader)
	attempt1Hash := common.HexToHash("0xb7862c896a6ba2711bccc0410184e46d793ea83b3e05470f1d359ea276d16bb5")
	sentAt := uint64(23456)
	confirmed := sentAt + config.EthGasBumpThreshold() + 1
	safe := confirmed + config.MinOutgoingConfirmations() - 1
	unconfirmedReceipt := ethpkg.TxReceipt{}
	confirmedReceipt := ethpkg.TxReceipt{
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
		eth.Register("eth_sendRawTransaction", attempt1Hash) // Initial tx attempt sent
		eth.Register("eth_getTransactionReceipt", unconfirmedReceipt)
	})
	j := cltest.CreateHelloWorldJobViaWeb(t, app, mockServer.URL)
	jr := cltest.WaitForJobRunToPendConfirmations(t, app.Store, cltest.CreateJobRunViaWeb(t, app, j))
	eth.EventuallyAllCalled(t)
	cltest.WaitForTxAttemptCount(t, app.Store, 1)

	jr.ObservedHeight = confirmedReceipt.BlockNumber
	require.NoError(t, app.Store.SaveJobRun(&jr))

	eth.Context("ethTx.Perform()#4 at block 23465", func(eth *cltest.EthMock) {
		eth.Register("eth_getTransactionReceipt", confirmedReceipt) // confirmed for gas bumped txat
		eth.Register("eth_getBalance", "0x0100")
		eth.Register("eth_call", "0x0100")
	})
	newHeads <- ethpkg.BlockHeader{Number: cltest.BigHexInt(safe)} // 23465
	eth.EventuallyAllCalled(t)
	cltest.WaitForTxAttemptCount(t, app.Store, 1)

	cltest.WaitForJobRunToComplete(t, app.Store, jr)

	eth.EventuallyAllCalled(t)
}

func TestIntegration_FeeBump(t *testing.T) {
	tickerResponse := `{"high": "10744.00", "last": "10583.75", "timestamp": "1512156162", "bid": "10555.13", "vwap": "10097.98", "volume": "17861.33960013", "low": "9370.11", "ask": "10583.00", "open": "9927.29"}`
	mockServer, assertCalled := cltest.NewHTTPMockServer(t, http.StatusOK, "GET", tickerResponse)
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

	unconfirmedReceipt := ethpkg.TxReceipt{}

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
	thirdTxConfirmedReceipt := ethpkg.TxReceipt{
		Hash:        attempt1Hash,
		BlockNumber: cltest.Int(thirdTxConfirmedAt),
	}
	thirdTxSafeAt := thirdTxSentAt + config.MinOutgoingConfirmations()

	newHeads := make(chan ethpkg.BlockHeader)
	eth := app.EthMock
	eth.Context("app.Start()", func(eth *cltest.EthMock) {
		eth.RegisterSubscription("newHeads", newHeads)
		eth.Register("eth_chainId", config.ChainID())
		eth.Register("eth_getTransactionCount", `0x0100`) // TxManager.ActivateAccount()
	})
	require.NoError(t, app.Store.ORM.CreateHead(cltest.Head(firstTxSentAt)))
	assert.NoError(t, app.Start())
	eth.EventuallyAllCalled(t)

	// This first run of the EthTx adapter creates an initial transaction which
	// starts unconfirmed
	eth.Context("ethTx.Perform()#1", func(eth *cltest.EthMock) {
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
		eth.Register("eth_getTransactionReceipt", unconfirmedReceipt)
	})
	newHeads <- ethpkg.BlockHeader{Number: cltest.BigHexInt(firstTxRemainsUnconfirmedAt)}
	eth.EventuallyAllCalled(t)
	jr = cltest.WaitForJobRunToPendConfirmations(t, app.Store, jr)

	// At the next head, the transaction remains unconfirmed but the gas bump
	// threshold has been met, so a new transaction is made with a higher amount
	// of gas ("bumped gas")
	eth.Context("ethTx.Perform()#3", func(eth *cltest.EthMock) {
		eth.Register("eth_getTransactionReceipt", unconfirmedReceipt)
		eth.Register("eth_sendRawTransaction", attempt2Hash)
	})
	newHeads <- ethpkg.BlockHeader{Number: cltest.BigHexInt(firstTxGasBumpAt)}
	eth.EventuallyAllCalled(t)
	jr = cltest.WaitForJobRunToPendConfirmations(t, app.Store, jr)
	cltest.WaitForTxAttemptCount(t, app.Store, 2)

	// Another head comes in and both transactions are still unconfirmed, more
	// waiting...
	eth.Context("ethTx.Perform()#4", func(eth *cltest.EthMock) {
		eth.Register("eth_getTransactionReceipt", unconfirmedReceipt)
		eth.Register("eth_getTransactionReceipt", unconfirmedReceipt)
	})
	newHeads <- ethpkg.BlockHeader{Number: cltest.BigHexInt(secondTxRemainsUnconfirmedAt)}
	eth.EventuallyAllCalled(t)
	jr = cltest.WaitForJobRunToPendConfirmations(t, app.Store, jr)

	// Now the second transaction attempt meets the gas bump threshold, so a
	// final transaction attempt shoud be made
	eth.Context("ethTx.Perform()#5", func(eth *cltest.EthMock) {
		eth.Register("eth_getTransactionReceipt", unconfirmedReceipt)
		eth.Register("eth_getTransactionReceipt", unconfirmedReceipt)
		eth.Register("eth_sendRawTransaction", attempt3Hash)
	})
	newHeads <- ethpkg.BlockHeader{Number: cltest.BigHexInt(secondTxGasBumpAt)}
	eth.EventuallyAllCalled(t)
	jr = cltest.WaitForJobRunToPendConfirmations(t, app.Store, jr)
	cltest.WaitForTxAttemptCount(t, app.Store, 3)

	// This third attempt has enough gas and gets confirmed, but has not yet
	// received sufficient confirmations, so we wait again...
	eth.Context("ethTx.Perform()#6", func(eth *cltest.EthMock) {
		eth.Register("eth_getTransactionReceipt", thirdTxConfirmedReceipt)
		eth.Register("eth_getBalance", "0x0000")
	})
	newHeads <- ethpkg.BlockHeader{Number: cltest.BigHexInt(thirdTxConfirmedAt)}
	eth.EventuallyAllCalled(t)
	jr = cltest.WaitForJobRunToPendConfirmations(t, app.Store, jr)

	// Finally the third attempt gets to a minimum number of safe confirmations,
	// the amount remaining in the account is printed (eth_getBalance, eth_call)
	eth.Context("ethTx.Perform()#7", func(eth *cltest.EthMock) {
		eth.Register("eth_getTransactionReceipt", thirdTxConfirmedReceipt)
		eth.Register("eth_getBalance", "0x100")
		eth.Register("eth_call", "0x100")
	})
	newHeads <- ethpkg.BlockHeader{Number: cltest.BigHexInt(thirdTxSafeAt)}
	eth.EventuallyAllCalled(t)
	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)

	require.Len(t, jr.TaskRuns, 4)
	value := cltest.MustResultString(t, jr.TaskRuns[0].Result)
	assert.Equal(t, tickerResponse, value)
	value = cltest.MustResultString(t, jr.TaskRuns[1].Result)
	assert.Equal(t, "10583.75", value)
	value = cltest.MustResultString(t, jr.TaskRuns[3].Result)
	assert.Equal(t, attempt1Hash.String(), value)
	value = cltest.MustResultString(t, jr.Result)
	assert.Equal(t, attempt1Hash.String(), value)
}

func TestIntegration_RunAt(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t, cltest.EthMockRegisterChainID)
	defer cleanup()
	app.InstantClock()

	require.NoError(t, app.Start())
	j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/run_at_job.json")

	initr := j.Initiators[0]
	assert.Equal(t, models.InitiatorRunAt, initr.Type)
	assert.Equal(t, "2018-01-08T18:12:01Z", utils.ISO8601UTC(initr.Time.Time))

	jrs := cltest.WaitForRuns(t, j, app.Store, 1)
	cltest.WaitForJobRunToComplete(t, app.Store, jrs[0])
}

func TestIntegration_EthLog(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	defer cleanup()

	eth := app.EthMock
	logs := make(chan ethpkg.Log, 1)
	eth.Context("app.Start()", func(eth *cltest.EthMock) {
		eth.Register("eth_chainId", app.Store.Config.ChainID())
		eth.RegisterSubscription("logs", logs)
	})
	require.NoError(t, app.Start())

	j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/eth_log_job.json")
	address := common.HexToAddress("0x3cCad4715152693fE3BC4460591e3D3Fbd071b42")

	initr := j.Initiators[0]
	assert.Equal(t, models.InitiatorEthLog, initr.Type)
	assert.Equal(t, address, initr.Address)

	logs <- cltest.LogFromFixture(t, "testdata/requestLog0original.json")
	jrs := cltest.WaitForRuns(t, j, app.Store, 1)
	cltest.WaitForJobRunToComplete(t, app.Store, jrs[0])
}

func TestIntegration_RunLog(t *testing.T) {
	triggeringBlockHash := cltest.NewHash()
	otherBlockHash := cltest.NewHash()

	tests := []struct {
		name             string
		logBlockHash     common.Hash
		receiptBlockHash common.Hash
		wantStatus       models.RunStatus
	}{
		{
			name:             "completed",
			logBlockHash:     triggeringBlockHash,
			receiptBlockHash: triggeringBlockHash,
			wantStatus:       models.RunStatusCompleted,
		},
		{
			name:             "ommered request",
			logBlockHash:     triggeringBlockHash,
			receiptBlockHash: otherBlockHash,
			wantStatus:       models.RunStatusErrored,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, cfgCleanup := cltest.NewConfig(t)
			defer cfgCleanup()
			config.Set("MIN_INCOMING_CONFIRMATIONS", 6)
			app, cleanup := cltest.NewApplicationWithConfig(t, config, cltest.LenientEthMock)
			defer cleanup()

			eth := app.EthMock
			logs := make(chan ethpkg.Log, 1)
			newHeads := eth.RegisterNewHeads()
			eth.Context("app.Start()", func(eth *cltest.EthMock) {
				eth.RegisterSubscription("logs", logs)
			})
			eth.Register("eth_chainId", config.ChainID())
			require.NoError(t, app.Start())

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
			require.Len(t, jr.TaskRuns, 1)
			assert.False(t, jr.TaskRuns[0].Confirmations.Valid)

			blockIncrease := app.Store.Config.MinIncomingConfirmations()
			minGlobalHeight := creationHeight + blockIncrease
			newHeads <- ethpkg.BlockHeader{Number: cltest.BigHexInt(minGlobalHeight)}
			<-time.After(time.Second)
			jr = cltest.JobRunStaysPendingConfirmations(t, app.Store, jr)
			assert.Equal(t, uint32(creationHeight+blockIncrease), jr.TaskRuns[0].Confirmations.Uint32)

			safeNumber := creationHeight + requiredConfs
			newHeads <- ethpkg.BlockHeader{Number: cltest.BigHexInt(safeNumber)}
			confirmedReceipt := ethpkg.TxReceipt{
				Hash:        runlog.TxHash,
				BlockHash:   &test.receiptBlockHash,
				BlockNumber: cltest.Int(creationHeight),
			}
			eth.Context("validateOnMainChain", func(ethMock *cltest.EthMock) {
				eth.Register("eth_getTransactionReceipt", confirmedReceipt)
			})

			jr = cltest.WaitForJobRunStatus(t, app.Store, jr, test.wantStatus)
			assert.True(t, jr.FinishedAt.Valid)
			assert.Equal(t, requiredConfs, jr.TaskRuns[0].Confirmations.Uint32)
			assert.True(t, eth.AllCalled(), eth.Remaining())
		})
	}
}

func TestIntegration_StartAt(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	defer cleanup()
	eth := app.EthMock
	eth.Register("eth_chainId", app.Store.Config.ChainID())
	require.NoError(t, app.Start())

	j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/start_at_job.json")
	startAt := cltest.ParseISO8601(t, "1970-01-01T00:00:00.000Z")
	assert.Equal(t, startAt, j.StartAt.Time)

	jr := cltest.CreateJobRunViaWeb(t, app, j)
	cltest.WaitForJobRunToComplete(t, app.Store, jr)
}

func TestIntegration_ExternalAdapter_RunLogInitiated(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	defer cleanup()

	eth := app.EthMock
	eth.Register("eth_chainId", app.Store.Config.ChainID())
	logs := make(chan ethpkg.Log, 1)
	newHeads := make(chan ethpkg.BlockHeader, 10)
	eth.Context("app.Start()", func(eth *cltest.EthMock) {
		eth.RegisterSubscription("logs", logs)
		eth.RegisterSubscription("newHeads", newHeads)
	})
	require.NoError(t, app.Start())

	eaValue := "87698118359"
	eaExtra := "other values to be used by external adapters"
	eaResponse := fmt.Sprintf(`{"data":{"result": "%v", "extra": "%v"}}`, eaValue, eaExtra)
	mockServer, ensureRequest := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", eaResponse)
	defer ensureRequest()

	bridgeJSON := fmt.Sprintf(`{"name":"randomNumber","url":"%v","confirmations":10}`, mockServer.URL)
	cltest.CreateBridgeTypeViaWeb(t, app, bridgeJSON)
	j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/log_initiated_bridge_type_job.json")

	logBlockNumber := 1
	runlog := cltest.NewRunLog(t, j.ID, cltest.NewAddress(), cltest.NewAddress(), logBlockNumber, `{}`)
	logs <- runlog
	jr := cltest.WaitForRuns(t, j, app.Store, 1)[0]
	cltest.WaitForJobRunToPendConfirmations(t, app.Store, jr)

	newHeads <- ethpkg.BlockHeader{Number: cltest.BigHexInt(logBlockNumber + 8)}
	cltest.WaitForJobRunToPendConfirmations(t, app.Store, jr)

	confirmedReceipt := ethpkg.TxReceipt{
		Hash:        runlog.TxHash,
		BlockHash:   &runlog.BlockHash,
		BlockNumber: cltest.Int(logBlockNumber),
	}
	eth.Context("validateOnMainChain", func(ethMock *cltest.EthMock) {
		eth.Register("eth_getTransactionReceipt", confirmedReceipt)
	})

	newHeads <- ethpkg.BlockHeader{Number: cltest.BigHexInt(logBlockNumber + 9)}
	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)

	tr := jr.TaskRuns[0]
	assert.Equal(t, "randomnumber", tr.TaskSpec.Type.String())
	value := cltest.MustResultString(t, tr.Result)
	assert.Equal(t, eaValue, value)
	res := tr.Result.Data.Get("extra")
	assert.Equal(t, eaExtra, res.String())

	assert.True(t, eth.AllCalled(), eth.Remaining())
}

// This test ensures that the response body of an external adapter are supplied
// as params to the successive task
func TestIntegration_ExternalAdapter_Copy(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t, cltest.EthMockRegisterChainID)
	defer cleanup()
	bridgeURL := cltest.WebURL(t, "https://test.chain.link/always")
	app.Store.Config.Set("BRIDGE_RESPONSE_URL", bridgeURL)
	require.NoError(t, app.Start())

	eaPrice := "1234"
	eaQuote := "USD"
	eaResponse := fmt.Sprintf(`{"data":{"price": "%v", "quote": "%v"}}`, eaPrice, eaQuote)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "POST", r.Method)
		require.Equal(t, "/", r.URL.Path)

		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		body := cltest.JSONFromBytes(t, b)
		data := body.Get("data")
		require.True(t, data.Exists())
		bodyParam := data.Get("bodyParam")
		require.True(t, bodyParam.Exists())
		require.Equal(t, true, bodyParam.Bool())

		url := body.Get("responseURL")
		require.Contains(t, url.String(), "https://test.chain.link/always/v2/runs")

		w.WriteHeader(http.StatusOK)
		_, err = io.WriteString(w, eaResponse)
		require.NoError(t, err)
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
	value := cltest.MustResultString(t, tr.Result)
	assert.Equal(t, eaPrice, value)
}

// This test ensures that an bridge adapter task is resumed from pending after
// sending out a request to an external adapter and waiting to receive a
// request back
func TestIntegration_ExternalAdapter_Pending(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t, cltest.EthMockRegisterChainID)
	defer cleanup()
	require.NoError(t, app.Start())

	bta := &models.BridgeTypeAuthentication{}
	var j models.JobSpec
	mockServer, cleanup := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", `{"pending":true}`,
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
	assert.Equal(t, gjson.Null, tr.Result.Data.Get("result").Type)

	jr = cltest.UpdateJobRunViaWeb(t, app, jr, bta, `{"data":{"result":"100"}}`)
	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)
	tr = jr.TaskRuns[0]
	assert.Equal(t, models.RunStatusCompleted, tr.Status)

	value := cltest.MustResultString(t, tr.Result)
	assert.Equal(t, "100", value)
}

func TestIntegration_WeiWatchers(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	defer cleanup()

	eth := app.EthMock
	eth.RegisterNewHead(1)
	logs := make(chan ethpkg.Log, 1)
	eth.Context("app.Start()", func(eth *cltest.EthMock) {
		eth.Register("eth_chainId", app.Config.ChainID())
		eth.RegisterSubscription("logs", logs)
	})

	log := cltest.LogFromFixture(t, "testdata/requestLog0original.json")
	mockServer, cleanup := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", `{"pending":true}`,
		func(_ http.Header, body string) {
			marshaledLog, err := json.Marshal(&log)
			assert.NoError(t, err)
			assert.JSONEq(t, string(marshaledLog), body)
		})
	defer cleanup()

	require.NoError(t, app.Start())

	j := cltest.NewJobWithLogInitiator()
	post := cltest.NewTask(t, "httppostwithunrestrictednetworkaccess", fmt.Sprintf(`{"url":"%v"}`, mockServer.URL))
	tasks := []models.TaskSpec{post}
	j.Tasks = tasks
	j = cltest.CreateJobSpecViaWeb(t, app, j)

	logs <- log

	jobRuns := cltest.WaitForRuns(t, j, app.Store, 1)
	cltest.WaitForJobRunToComplete(t, app.Store, jobRuns[0])
}

func TestIntegration_MultiplierInt256(t *testing.T) {
	app, cleanup := cltest.NewApplication(t, cltest.EthMockRegisterChainID)
	defer cleanup()
	require.NoError(t, app.Start())

	j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/int256_job.json")
	jr := cltest.CreateJobRunViaWeb(t, app, j, `{"result":"-10221.30"}`)
	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)

	value := cltest.MustResultString(t, jr.Result)
	assert.Equal(t, "0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0674e", value)
}

func TestIntegration_MultiplierUint256(t *testing.T) {
	app, cleanup := cltest.NewApplication(t, cltest.EthMockRegisterChainID)
	defer cleanup()
	require.NoError(t, app.Start())

	j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/uint256_job.json")
	jr := cltest.CreateJobRunViaWeb(t, app, j, `{"result":"10221.30"}`)
	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)

	value := cltest.MustResultString(t, jr.Result)
	assert.Equal(t, "0x00000000000000000000000000000000000000000000000000000000000f98b2", value)
}

func TestIntegration_NonceManagement_firstRunWithExistingTxs(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t, cltest.EthMockRegisterChainID)
	defer cleanup()

	eth := app.EthMock
	newHeads := make(chan ethpkg.BlockHeader)
	eth.Context("app.Start()", func(eth *cltest.EthMock) {
		eth.RegisterSubscription("newHeads", newHeads)
		eth.Register("eth_getTransactionCount", `0x100`) // activate account nonce
	})
	require.NoError(t, app.Store.ORM.CreateHead(cltest.Head(100)))
	require.NoError(t, app.StartAndConnect())

	j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/web_initiated_eth_tx_job.json")
	hash := common.HexToHash("0xb7862c896a6ba2711bccc0410184e46d793ea83b3e05470f1d359ea276d16bb5")

	createCompletedJobRun := func(blockNumber uint64, expectedNonce uint64) {
		confirmedBlockNumber := blockNumber - app.Store.Config.MinOutgoingConfirmations()

		eth.Context("ethTx.Perform()", func(eth *cltest.EthMock) {
			eth.Register("eth_sendRawTransaction", hash)
			eth.Register("eth_getTransactionReceipt", ethpkg.TxReceipt{
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

	newHeads <- ethpkg.BlockHeader{Number: cltest.BigHexInt(200)}
	createCompletedJobRun(200, uint64(0x101))
}

func TestIntegration_CreateServiceAgreement(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()

	eth := app.EthMock
	logs := make(chan ethpkg.Log, 1)
	eth.Context("app.Start()", func(eth *cltest.EthMock) {
		eth.RegisterSubscription("logs", logs)
		eth.Register("eth_getTransactionCount", `0x100`)
		eth.Register("eth_chainId", app.Store.Config.ChainID())
	})
	assert.NoError(t, app.StartAndConnect())
	endAt := time.Now().AddDate(0, 10, 0).Round(time.Second).UTC()
	sa := cltest.CreateServiceAgreementViaWeb(t, app, "fixtures/web/noop_agreement.json", endAt)

	assert.NotEqual(t, "", sa.ID)
	j := cltest.FindJob(t, app.Store, sa.JobSpecID)

	assert.Equal(t, cltest.NewLink(t, "1000000000000000000"), sa.Encumbrance.Payment)
	assert.Equal(t, uint64(300), sa.Encumbrance.Expiration)

	assert.Equal(t, endAt, sa.Encumbrance.EndAt.Time)
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
	app, cleanup := cltest.NewApplicationWithConfig(t, config, cltest.EthMockRegisterChainID)
	defer cleanup()

	app.InstantClock()
	require.NoError(t, app.Start())

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
	cltest.WaitForJobRunToComplete(t, app.Store, run)
}

func TestIntegration_SleepAdapter(t *testing.T) {
	t.Parallel()

	sleepSeconds := 4
	app, cleanup := cltest.NewApplication(t, cltest.EthMockRegisterChainID)
	app.Config.Set("ENABLE_EXPERIMENTAL_ADAPTERS", "true")
	defer cleanup()
	require.NoError(t, app.Start())

	j := cltest.FixtureCreateJobViaWeb(t, app, "./testdata/sleep_job.json")

	runInput := fmt.Sprintf("{\"until\": \"%s\"}", time.Now().Local().Add(time.Second*time.Duration(sleepSeconds)))
	jr := cltest.CreateJobRunViaWeb(t, app, j, runInput)

	cltest.WaitForJobRunStatus(t, app.Store, jr, models.RunStatusInProgress)
	cltest.JobRunStays(t, app.Store, jr, models.RunStatusInProgress, time.Second)
	cltest.WaitForJobRunToComplete(t, app.Store, jr)
}

func TestIntegration_ExternalInitiator(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t, cltest.EthMockRegisterChainID)
	defer cleanup()
	require.NoError(t, app.Start())

	exInitr := struct {
		Header http.Header
		Body   web.JobSpecNotice
	}{}
	eiMockServer, assertCalled := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", "",
		func(header http.Header, body string) {
			exInitr.Header = header
			err := json.Unmarshal([]byte(body), &exInitr.Body)
			require.NoError(t, err)
		},
	)
	defer assertCalled()

	eiCreate := map[string]string{
		"name": "someCoin",
		"url":  eiMockServer.URL,
	}
	eiCreateJSON, err := json.Marshal(eiCreate)
	require.NoError(t, err)
	eip := cltest.CreateExternalInitiatorViaWeb(t, app, string(eiCreateJSON))

	eia := &auth.Token{
		AccessKey: eip.AccessKey,
		Secret:    eip.Secret,
	}
	ei, err := app.Store.FindExternalInitiator(eia)
	require.NoError(t, err)

	require.Equal(t, eiCreate["url"], ei.URL.String())
	require.Equal(t, strings.ToLower(eiCreate["name"]), ei.Name)
	require.Equal(t, eip.AccessKey, ei.AccessKey)
	require.Equal(t, eip.OutgoingSecret, ei.OutgoingSecret)

	jobSpec := cltest.FixtureCreateJobViaWeb(t, app, "./testdata/external_initiator_job.json")
	assert.Equal(t,
		eip.OutgoingToken,
		exInitr.Header.Get(web.ExternalInitiatorAccessKeyHeader),
	)
	assert.Equal(t,
		eip.OutgoingSecret,
		exInitr.Header.Get(web.ExternalInitiatorSecretHeader),
	)
	expected := web.JobSpecNotice{
		JobID:  jobSpec.ID,
		Type:   models.InitiatorExternal,
		Params: cltest.JSONFromString(t, `{"foo":"bar"}`),
	}
	assert.Equal(t, expected, exInitr.Body)

	jobRun := cltest.CreateJobRunViaExternalInitiator(t, app, jobSpec, *eia, "")
	_, err = app.Store.JobRunsFor(jobRun.ID)
	assert.NoError(t, err)
	cltest.WaitForJobRunToComplete(t, app.Store, jobRun)
}

func TestIntegration_ExternalInitiator_WithoutURL(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t, cltest.EthMockRegisterChainID)
	defer cleanup()
	require.NoError(t, app.Start())

	eiCreate := map[string]string{
		"name": "someCoin",
	}
	eiCreateJSON, err := json.Marshal(eiCreate)
	require.NoError(t, err)
	eip := cltest.CreateExternalInitiatorViaWeb(t, app, string(eiCreateJSON))

	eia := &auth.Token{
		AccessKey: eip.AccessKey,
		Secret:    eip.Secret,
	}
	ei, err := app.Store.FindExternalInitiator(eia)
	require.NoError(t, err)

	require.Equal(t, strings.ToLower(eiCreate["name"]), ei.Name)
	require.Equal(t, eip.AccessKey, ei.AccessKey)
	require.Equal(t, eip.OutgoingSecret, ei.OutgoingSecret)

	jobSpec := cltest.FixtureCreateJobViaWeb(t, app, "./testdata/external_initiator_job.json")

	jobRun := cltest.CreateJobRunViaExternalInitiator(t, app, jobSpec, *eia, "")
	_, err = app.Store.JobRunsFor(jobRun.ID)
	assert.NoError(t, err)
	cltest.WaitForJobRunToComplete(t, app.Store, jobRun)
}

func TestIntegration_AuthToken(t *testing.T) {
	app, cleanup := cltest.NewApplication(t, cltest.EthMockRegisterChainID)
	defer cleanup()

	require.NoError(t, app.Start())

	// set up user
	app.MustSeedUserAPIKey()

	url := app.Config.ClientNodeURL() + "/v2/config"
	headers := make(map[string]string)
	headers[web.APIKey] = cltest.APIKey
	headers[web.APISecret] = cltest.APISecret
	buf := bytes.NewBufferString(`{"ethGasPriceDefault":15000000}`)

	resp, cleanup := cltest.UnauthenticatedPatch(t, url, buf, headers)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)
}

func TestIntegration_FluxMonitor_Deviation(t *testing.T) {
	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	minPayment := app.Store.Config.MinimumContractPayment().ToInt().Uint64()
	availableFunds := minPayment * 100

	// Start, connect, and initialize node
	newHeads := make(chan ethpkg.BlockHeader)
	eth := app.EthMock
	eth.Context("app.Start()", func(eth *cltest.EthMock) {
		eth.RegisterSubscription("newHeads", newHeads)
		eth.Register("eth_getTransactionCount", `0x0100`) // TxManager.ActivateAccount()
		eth.Register("eth_chainId", app.Store.Config.ChainID())
	})
	require.NoError(t, app.StartAndConnect())
	eth.EventuallyAllCalled(t)

	// Configure fake Eth Node to return 10,000 cents when FM initiates price.
	eth.Context("Flux Monitor initializes price", func(mock *cltest.EthMock) {
		hex := cltest.MakeRoundStateReturnData(2, true, 10000, 7, 0, availableFunds, minPayment, 1)
		mock.Register("eth_call", hex)
	})

	// Have server respond with 102 for price when FM checks external price
	// adapter for deviation. 102 is enough deviation to trigger a job run.
	priceResponse := `{"data":{"result": 102}}`
	mockServer, assertCalled := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", priceResponse)
	defer assertCalled()

	// Single task ethTx receives configuration from FM init and writes to chain.
	attemptHash := cltest.NewHash()
	confirmedReceipt := ethpkg.TxReceipt{
		Hash:        attemptHash,
		BlockNumber: cltest.Int(1),
	}

	eth.Context("ethTx.Perform() for initial send", func(eth *cltest.EthMock) {
		eth.Register("eth_sendRawTransaction", attemptHash)         // Initial tx attempt sent
		eth.Register("eth_getTransactionReceipt", confirmedReceipt) // confirmed for gas bumped txat
		eth.Register("eth_getBlockByNumber", ethpkg.Block{Number: hexutil.Uint64(1)})
		eth.Register("eth_getLogs", []ethpkg.Log{})
	})

	// Create FM Job, and wait for job run to start because the above criteria initiates a run.
	buffer := cltest.MustReadFile(t, "testdata/flux_monitor_job.json")
	var job models.JobSpec
	err := json.Unmarshal(buffer, &job)
	require.NoError(t, err)
	job.Initiators[0].InitiatorParams.Feeds = cltest.JSONFromString(t, fmt.Sprintf(`["%s"]`, mockServer.URL))
	job.Initiators[0].InitiatorParams.PollingInterval = models.MustMakeDuration(15 * time.Second)

	j := cltest.CreateJobSpecViaWeb(t, app, job)
	jrs := cltest.WaitForRuns(t, j, app.Store, 1)

	eth.EventuallyAllCalled(t)

	// Send a head w block number 10, high enough to mark ethtx as safe.
	eth.Context("ethTx.Perform() for safe", func(eth *cltest.EthMock) {
		eth.Register("eth_getTransactionReceipt", confirmedReceipt)
		eth.Register("eth_getBalance", "0x100")
		eth.Register("eth_call", cltest.MustEVMUintHexFromBase10String(t, "256"))
	})
	newHeads <- ethpkg.BlockHeader{Number: cltest.BigHexInt(10)}

	// Check the FM price on completed run output
	jr := cltest.WaitForJobRunToComplete(t, app.GetStore(), jrs[0])

	requestParams := jr.RunRequest.RequestParams
	assert.Equal(t, "102", requestParams.Get("result").String())
	assert.Equal(
		t,
		"0x3cCad4715152693fE3BC4460591e3D3Fbd071b42", // from testdata/flux_monitor_job.json
		requestParams.Get("address").String())
	assert.Equal(t, "0x202ee0ed", requestParams.Get("functionSelector").String())
	assert.Equal(
		t,
		"0x0000000000000000000000000000000000000000000000000000000000000002",
		requestParams.Get("dataPrefix").String())
}

func TestIntegration_FluxMonitor_NewRound(t *testing.T) {
	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	minPayment := app.Store.Config.MinimumContractPayment().ToInt().Uint64()
	availableFunds := minPayment * 100

	// Start, connect, and initialize node
	eth := app.EthMock
	eth.Context("app.StartAndConnect()", func(eth *cltest.EthMock) {
		eth.Register("eth_getTransactionCount", `0x0100`) // TxManager.ActivateAccount()
		eth.Register("eth_chainId", app.Store.Config.ChainID())
	})
	require.NoError(t, app.StartAndConnect())
	eth.EventuallyAllCalled(t)

	// Configure fake Eth Node to return 10,000 cents when FM initiates price.
	eth.Context("Flux Monitor queries FluxAggregator.RoundState()", func(mock *cltest.EthMock) {
		hex := cltest.MakeRoundStateReturnData(2, true, 10000, 7, 0, availableFunds, minPayment, 1)
		mock.Register("eth_call", hex)
	})

	// Have price adapter server respond with 100 for price on initialization,
	// NOT enough for deviation.
	priceResponse := `{"data":{"result": 100}}`
	mockServer, assertCalled := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", priceResponse)
	defer assertCalled()

	// Prepare new rounds logs subscription to be called by new FM job
	newRounds := make(chan ethpkg.Log)
	eth.RegisterSubscription("logs", newRounds)

	eth.Context("Log Broadcaster backfills logs", func(mock *cltest.EthMock) {
		eth.Register("eth_getBlockByNumber", ethpkg.Block{Number: hexutil.Uint64(1)})
		eth.Register("eth_getLogs", []ethpkg.Log{})
	})

	// Create FM Job, and ensure no runs because above criteria has no deviation.
	buffer := cltest.MustReadFile(t, "testdata/flux_monitor_job.json")
	var job models.JobSpec
	err := json.Unmarshal(buffer, &job)
	require.NoError(t, err)
	job.Initiators[0].InitiatorParams.Feeds = cltest.JSONFromString(t, fmt.Sprintf(`["%s"]`, mockServer.URL))
	job.Initiators[0].InitiatorParams.PollingInterval = models.MustMakeDuration(15 * time.Second)
	// NOTE: Must set a very large idle timeout here because log.StartedAt is
	// hardcoded to 15 i.e. 1970-01-01T00:00:15Z.
	// Otherwise this makes the idle timer immediately fire and the subsequent
	// eth_call is not handled.
	job.Initiators[0].InitiatorParams.IdleThreshold = models.MustMakeDuration(math.MaxInt64)

	j := cltest.CreateJobSpecViaWeb(t, app, job)
	_ = cltest.WaitForRuns(t, j, app.Store, 0)

	eth.EventuallyAllCalled(t)

	// Send a NewRound log event to trigger a run.
	log := cltest.LogFromFixture(t, "testdata/new_round_log.json")
	log.Address = job.Initiators[0].InitiatorParams.Address

	attemptHash := cltest.NewHash()
	confirmedReceipt := ethpkg.TxReceipt{
		Hash:        attemptHash,
		BlockNumber: cltest.Int(1),
	}
	eth.Context("ethTx.Perform() for new round send", func(eth *cltest.EthMock) {
		eth.Register("eth_sendRawTransaction", attemptHash)         // Initial tx attempt sent
		eth.Register("eth_getTransactionReceipt", confirmedReceipt) // confirmed for gas bumped txat
	})
	eth.Context("Flux Monitor queries FluxAggregator.RoundState()", func(mock *cltest.EthMock) {
		hex := cltest.MakeRoundStateReturnData(2, true, 10000, 7, 0, availableFunds, minPayment, 1)
		mock.Register("eth_call", hex)
	})
	newRounds <- log
	jrs := cltest.WaitForRuns(t, j, app.Store, 1)
	_ = cltest.WaitForJobRunToPendConfirmations(t, app.Store, jrs[0])
	eth.EventuallyAllCalled(t)
}

func TestIntegration_RandomnessRequest(t *testing.T) {
	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	eth := app.MockCallerSubscriberClient()
	logs := make(chan ethpkg.Log, 1)
	txHash := cltest.NewHash()
	eth.Context("app.Start()", func(eth *cltest.EthMock) {
		eth.RegisterSubscription("logs", logs)
		eth.Register("eth_getTransactionCount", `0x100`) // activate account nonce
		eth.Register("eth_sendRawTransaction", txHash)
		eth.Register("eth_getTransactionReceipt", ethpkg.TxReceipt{
			Hash:        cltest.NewHash(),
			BlockNumber: cltest.Int(10),
		})
	})
	config, cfgCleanup := cltest.NewConfig(t)
	defer cfgCleanup()
	eth.Register("eth_chainId", config.ChainID())
	require.NoError(t, app.Start())

	j := cltest.FixtureCreateJobViaWeb(t, app, "testdata/randomness_job.json")
	rawKey := j.Tasks[0].Params.Get("publicKey").String()
	pk, err := vrfkey.NewPublicKeyFromHex(rawKey)
	require.NoError(t, err)
	var sk int64 = 1
	coordinatorAddress := j.Initiators[0].Address

	provingKey := vrfkey.NewPrivateKeyXXXTestingOnly(big.NewInt(sk))
	require.Equal(t, &provingKey.PublicKey, pk,
		"public key in fixture %s does not match secret key in test %d (which has public key %s)",
		pk, sk, provingKey.PublicKey.String())
	app.Store.VRFKeyStore.StoreInMemoryXXXTestingOnly(provingKey)
	rawID := []byte(j.ID.String()) // CL requires ASCII hex encoding of jobID
	r := vrf.RandomnessRequestLog{
		KeyHash: provingKey.PublicKey.MustHash(),
		Seed:    big.NewInt(2),
		JobID:   common.BytesToHash(rawID),
		Sender:  cltest.NewAddress(),
		Fee:     assets.NewLink(100),
	}
	requestlog := cltest.NewRandomnessRequestLog(t, r, coordinatorAddress, 1)

	logs <- requestlog
	cltest.WaitForRuns(t, j, app.Store, 1)
	runs, err := app.Store.JobRunsFor(j.ID)
	assert.NoError(t, err)
	require.Len(t, runs, 1)
	jr := runs[0]
	require.Len(t, jr.TaskRuns, 2)
	assert.False(t, jr.TaskRuns[0].Confirmations.Valid)
	attempts := cltest.WaitForTxAttemptCount(t, app.Store, 1)
	require.True(t, eth.AllCalled(), eth.Remaining())
	require.Len(t, attempts, 1)

	rawTx := attempts[0].SignedRawTx
	var tx *types.Transaction
	require.NoError(t, rlp.DecodeBytes(rawTx, &tx))
	fixtureToAddress := j.Tasks[1].Params.Get("address").String()
	require.Equal(t, *tx.To(), common.HexToAddress(fixtureToAddress))
	payload := tx.Data()
	require.Equal(t, hexutil.Encode(payload[:4]), vrf.FulfillSelector())
	proofContainer := make(map[string]interface{})
	err = vrf.FulfillMethod().Inputs.UnpackIntoMap(proofContainer, payload[4:])
	require.NoError(t, err)
	proof, ok := proofContainer["_proof"].([]byte)
	require.True(t, ok)
	require.Len(t, proof, vrf.ProofLength)
	publicPoint, err := provingKey.PublicKey.Point()
	require.NoError(t, err)
	require.Equal(t, proof[:64], secp256k1.LongMarshal(publicPoint))
	goProof, err := vrf.UnmarshalSolidityProof(proof)
	require.NoError(t, err, "problem parsing solidity proof")
	proofValid, err := goProof.VerifyVRFProof()
	require.NoError(t, err, "problem verifying solidity proof")
	require.True(t, proofValid, "vrf proof was invalid: %s", goProof.String())

	// Check that a log from a different address is rejected. (The node will only
	// ever see this situation if the ethereum.FilterQuery for this job breaks,
	// but it's hard to test that without a full integration test.)
	badAddress := common.HexToAddress("0x0000000000000000000000000000000000000001")
	badRequestlog := cltest.NewRandomnessRequestLog(t, r, badAddress, 1)
	logs <- badRequestlog
	expectedLogTemplate := `log received from address %s, but expect logs from %s`
	expectedLog := fmt.Sprintf(expectedLogTemplate, badAddress.String(),
		coordinatorAddress.String())
	millisecondsWaited := 0
	expectedLogDeadline := 200
	for !strings.Contains(cltest.MemoryLogTestingOnly().String(), expectedLog) &&
		millisecondsWaited < expectedLogDeadline {
		time.Sleep(time.Millisecond)
		millisecondsWaited += 1
		if millisecondsWaited >= expectedLogDeadline {
			assert.Fail(t, "message about log with bad source address not found")
		}
	}
}
