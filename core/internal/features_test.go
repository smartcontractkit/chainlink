package internal_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	goEthereumEth "github.com/ethereum/go-ethereum/eth"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/multiwordconsumer"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/operator"
	"github.com/smartcontractkit/libocr/gethwrappers/linktoken"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

var oneETH = assets.Eth(*big.NewInt(1000000000000000000))

func TestIntegration_Scheduler(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	app.Start()

	j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/scheduler_job.json")

	cltest.WaitForRunsAtLeast(t, j, app.Store, 1)

	initr := j.Initiators[0]
	assert.Equal(t, models.InitiatorCron, initr.Type)
	assert.Equal(t, "CRON_TZ=UTC * * * * * *", string(initr.Schedule), "Wrong cron schedule saved")
}

func TestIntegration_HttpRequestWithHeaders(t *testing.T) {
	config, cfgCleanup := cltest.NewConfig(t)
	defer cfgCleanup()

	gethClient := new(mocks.GethClient)
	rpcClient := new(mocks.RPCClient)
	sub := new(mocks.Subscription)
	chchNewHeads := make(chan chan<- *models.Head, 1)

	app, appCleanup := cltest.NewApplicationWithConfigAndKey(t, config,
		eth.NewClientWith(rpcClient, gethClient),
	)
	defer appCleanup()

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

	confirmed := int64(23456)
	safe := confirmed + int64(config.MinRequiredOutgoingConfirmations())
	inLongestChain := safe - int64(config.GasUpdaterBlockDelay())

	rpcClient.On("EthSubscribe", mock.Anything, mock.Anything, "newHeads").
		Run(func(args mock.Arguments) { chchNewHeads <- args.Get(1).(chan<- *models.Head) }).
		Return(sub, nil)
	rpcClient.On("CallContext", mock.Anything, mock.Anything, "eth_getBlockByNumber", mock.Anything, false).
		Run(func(args mock.Arguments) {
			head := args.Get(1).(**models.Head)
			*head = cltest.Head(inLongestChain)
		}).
		Return(nil)

	gethClient.On("ChainID", mock.Anything).Return(config.ChainID(), nil)
	gethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(oneETH.ToInt(), nil)
	gethClient.On("BlockByNumber", mock.Anything, big.NewInt(inLongestChain)).Return(cltest.BlockWithTransactions(), nil)

	gethClient.On("SendTransaction", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			tx, ok := args.Get(1).(*types.Transaction)
			require.True(t, ok)
			gethClient.On("TransactionReceipt", mock.Anything, mock.Anything).
				Return(&types.Receipt{TxHash: tx.Hash(), BlockNumber: big.NewInt(confirmed)}, nil)
		}).
		Return(nil).Once()

	sub.On("Err").Return(nil)
	sub.On("Unsubscribe").Return(nil).Maybe()

	assert.NoError(t, app.StartAndConnect())

	newHeads := <-chchNewHeads

	j := cltest.CreateHelloWorldJobViaWeb(t, app, mockServer.URL)
	jr := cltest.WaitForJobRunToPendOutgoingConfirmations(t, app.Store, cltest.CreateJobRunViaWeb(t, app, j))

	app.EthBroadcaster.Trigger()
	cltest.WaitForEthTxAttemptCount(t, app.Store, 1)

	// Do the thing
	newHeads <- cltest.Head(safe)

	cltest.WaitForJobRunToComplete(t, app.Store, jr)
}

func TestIntegration_RunAt(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
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
	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()

	eth := app.EthMock
	logs := make(chan models.Log, 1)
	eth.Context("app.Start()", func(eth *cltest.EthMock) {
		eth.RegisterSubscription("logs", logs)
		eth.Register("eth_getTransactionReceipt", &types.Receipt{})
	})
	require.NoError(t, app.StartAndConnect())

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
			app, cleanup := cltest.NewApplicationWithConfig(t, config,
				cltest.LenientEthMock,
				cltest.EthMockRegisterGetBlockByNumber,
				cltest.EthMockRegisterGetBalance,
			)
			defer cleanup()

			eth := app.EthMock
			logs := make(chan types.Log, 1)
			newHeads := eth.RegisterNewHeads()
			eth.Context("app.Start()", func(eth *cltest.EthMock) {
				eth.RegisterSubscription("logs", logs)
			})
			eth.Register("eth_chainId", config.ChainID())
			require.NoError(t, app.Start())

			j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/runlog_noop_job.json")
			requiredConfs := int64(100)

			initr := j.Initiators[0]
			assert.Equal(t, models.InitiatorRunLog, initr.Type)

			creationHeight := int64(1)
			runlog := cltest.NewRunLog(t, j.ID, cltest.NewAddress(), cltest.NewAddress(), int(creationHeight), `{}`)
			runlog.BlockHash = test.logBlockHash
			logs <- runlog
			cltest.WaitForRuns(t, j, app.Store, 1)

			runs, err := app.Store.JobRunsFor(j.ID)
			assert.NoError(t, err)
			jr := runs[0]
			cltest.WaitForJobRunToPendIncomingConfirmations(t, app.Store, jr)
			require.Len(t, jr.TaskRuns, 1)
			assert.False(t, jr.TaskRuns[0].ObservedIncomingConfirmations.Valid)

			blockIncrease := int64(app.Store.Config.MinIncomingConfirmations())
			minGlobalHeight := creationHeight + blockIncrease
			newHeads <- cltest.Head(minGlobalHeight)
			<-time.After(time.Second)
			jr = cltest.JobRunStaysPendingIncomingConfirmations(t, app.Store, jr)
			assert.Equal(t, int64(creationHeight+blockIncrease), int64(jr.TaskRuns[0].ObservedIncomingConfirmations.Uint32))

			safeNumber := creationHeight + requiredConfs
			newHeads <- cltest.Head(safeNumber)
			confirmedReceipt := &types.Receipt{
				TxHash:      runlog.TxHash,
				BlockHash:   test.receiptBlockHash,
				BlockNumber: big.NewInt(creationHeight),
			}
			eth.Context("validateOnMainChain", func(ethMock *cltest.EthMock) {
				eth.Register("eth_getTransactionReceipt", confirmedReceipt)
			})

			app.EthBroadcaster.Trigger()
			jr = cltest.WaitForJobRunStatus(t, app.Store, jr, test.wantStatus)
			assert.True(t, jr.FinishedAt.Valid)
			assert.Equal(t, int64(requiredConfs), int64(jr.TaskRuns[0].ObservedIncomingConfirmations.Uint32))
			assert.True(t, eth.AllCalled(), eth.Remaining())
		})
	}
}

func TestIntegration_StartAt(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterGetBalance,
	)
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

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterGetBlockByNumber,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()

	eth := app.EthMock
	eth.Register("eth_chainId", app.Store.Config.ChainID())
	logs := make(chan models.Log, 1)
	newHeads := make(chan *models.Head, 10)
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
	cltest.WaitForJobRunToPendIncomingConfirmations(t, app.Store, jr)

	newHeads <- cltest.Head(logBlockNumber + 8)
	cltest.WaitForJobRunToPendIncomingConfirmations(t, app.Store, jr)

	confirmedReceipt := &types.Receipt{
		TxHash:      runlog.TxHash,
		BlockHash:   runlog.BlockHash,
		BlockNumber: big.NewInt(int64(logBlockNumber)),
	}
	eth.Context("validateOnMainChain", func(ethMock *cltest.EthMock) {
		eth.Register("eth_getTransactionReceipt", confirmedReceipt)
	})

	newHeads <- cltest.Head(logBlockNumber + 9)
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

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
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
	value := cltest.MustResultString(t, tr.Result)
	assert.Equal(t, eaPrice, value)
}

// This test ensures that an bridge adapter task is resumed from pending after
// sending out a request to an external adapter and waiting to receive a
// request back
func TestIntegration_ExternalAdapter_Pending(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
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

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterGetBlockByNumber,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()

	eth := app.EthMock
	eth.RegisterNewHead(1)
	logs := make(chan models.Log, 1)
	eth.Context("app.Start()", func(eth *cltest.EthMock) {
		eth.Register("eth_chainId", app.Config.ChainID())
		eth.RegisterSubscription("logs", logs)
		eth.Register("eth_getTransactionReceipt", &types.Receipt{})
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
	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/int256_job.json")
	jr := cltest.CreateJobRunViaWeb(t, app, j, `{"result":"-10221.30"}`)
	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)

	value := cltest.MustResultString(t, jr.Result)
	assert.Equal(t, "0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0674e", value)
}

func TestIntegration_MultiplierUint256(t *testing.T) {
	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/uint256_job.json")
	jr := cltest.CreateJobRunViaWeb(t, app, j, `{"result":"10221.30"}`)
	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)

	value := cltest.MustResultString(t, jr.Result)
	assert.Equal(t, "0x00000000000000000000000000000000000000000000000000000000000f98b2", value)
}

func TestIntegration_SyncJobRuns(t *testing.T) {
	t.Parallel()
	wsserver, wsserverCleanup := cltest.NewEventWebSocketServer(t)
	defer wsserverCleanup()

	config, _ := cltest.NewConfig(t)
	config.Set("EXPLORER_URL", wsserver.URL.String())
	app, cleanup := cltest.NewApplicationWithConfig(t, config,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	kst := new(mocks.KeyStoreInterface)
	app.Store.KeyStore = kst
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
	kst.AssertExpectations(t)
}

func TestIntegration_SleepAdapter(t *testing.T) {
	t.Parallel()

	sleepSeconds := 4
	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	app.Config.Set("ENABLE_EXPERIMENTAL_ADAPTERS", "true")
	defer cleanup()
	require.NoError(t, app.Start())

	j := cltest.FixtureCreateJobViaWeb(t, app, "./testdata/sleep_job.json")

	runInput := fmt.Sprintf("{\"until\": \"%s\"}", time.Now().Local().Add(time.Second*time.Duration(sleepSeconds)))
	jr := cltest.CreateJobRunViaWeb(t, app, j, runInput)

	cltest.WaitForJobRunStatus(t, app.Store, jr, models.RunStatusInProgress)
	cltest.JobRunStays(t, app.Store, jr, models.RunStatusInProgress, 3*time.Second)
	cltest.WaitForJobRunToComplete(t, app.Store, jr)
}

func TestIntegration_ExternalInitiator(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
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

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
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
	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()

	require.NoError(t, app.Start())

	// set up user
	mockUser := cltest.MustRandomUser()
	apiToken := auth.Token{AccessKey: cltest.APIKey, Secret: cltest.APISecret}
	require.NoError(t, mockUser.SetAuthToken(&apiToken))
	require.NoError(t, app.Store.SaveUser(&mockUser))

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
	gethClient := new(mocks.GethClient)
	rpcClient := new(mocks.RPCClient)
	sub := new(mocks.Subscription)

	config, cfgCleanup := cltest.NewConfig(t)
	defer cfgCleanup()
	app, appCleanup := cltest.NewApplicationWithConfig(t, config,
		eth.NewClientWith(rpcClient, gethClient),
	)
	defer appCleanup()

	// Helps to avoid deadlocks caused by multiple tests simultaneously writing
	// to eth_txes in different order and locking on nonce
	//
	// We have two tests (transactions because we're using txdb) which, for the same DefaultKeyAddress, try to:
	// Tx1:
	//     Read the nonce from keys (locking on address, nonce)
	//     Write to eth_txes
	// Tx 2:
	//     Write to eth_txes (I guess updating the state of the tx as completed?)
	//     Update nonce in keys (blocked on lock from tx 1)
	//
	// If every test/transaction has a different random nonce, then tx 2 doesn't have to wait for tx 1 to release the lock on that (address, nonce)

	cltest.RandomizeNonce(t, app.Store)

	kst := new(mocks.KeyStoreInterface)
	kst.On("HasAccountWithAddress", cltest.DefaultKeyAddress).Return(true)
	kst.On("GetAccountByAddress", mock.Anything).Maybe().Return(accounts.Account{}, nil)
	kst.On("SignTx", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(&types.Transaction{}, nil)

	app.Store.KeyStore = kst

	// Start, connect, and initialize node
	sub.On("Err").Return(nil).Maybe()
	sub.On("Unsubscribe").Return(nil).Maybe()
	gethClient.On("ChainID", mock.Anything).Return(app.Store.Config.ChainID(), nil)
	gethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(oneETH.ToInt(), nil)
	chchNewHeads := make(chan chan<- *models.Head, 1)
	rpcClient.On("EthSubscribe", mock.Anything, mock.Anything, "newHeads").
		Run(func(args mock.Arguments) { chchNewHeads <- args.Get(1).(chan<- *models.Head) }).
		Return(sub, nil)

	logsSub := new(mocks.Subscription)
	logsSub.On("Err").Return(nil)
	logsSub.On("Unsubscribe").Return(nil).Maybe()

	// GetOracles()
	rpcClient.On(
		"Call",
		mock.Anything,
		"eth_call",
		mock.MatchedBy(func(callArgs eth.CallArgs) bool {
			if (callArgs.To.Hex() == "0x3cCad4715152693fE3BC4460591e3D3Fbd071b42") && (hexutil.Encode(callArgs.Data) == "0x40884c52") {
				return true
			}
			return false
		}),
		mock.Anything,
		mock.Anything,
	).Return(nil).Once()

	err := app.StartAndConnect()
	require.NoError(t, err)

	gethClient.AssertExpectations(t)
	sub.AssertExpectations(t)

	// Configure fake Eth Node to return 10,000 cents when FM initiates price.
	minPayment := app.Store.Config.MinimumContractPayment().ToInt().Uint64()
	availableFunds := minPayment * 100
	rpcClient.On("Call", mock.Anything, "eth_call", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			*args.Get(0).(*hexutil.Bytes) = cltest.MakeRoundStateReturnData(2, true, 10000, 7, 0, availableFunds, minPayment, 1)
		}).
		Return(nil).
		Twice()

	// Have server respond with 102 for price when FM checks external price
	// adapter for deviation. 102 is enough deviation to trigger a job run.
	priceResponse := `{"data":{"result": 102}}`
	mockServer, assertCalled := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", priceResponse)
	defer assertCalled()

	confirmed := int64(23456)
	safe := confirmed + int64(config.MinRequiredOutgoingConfirmations())
	inLongestChain := safe - int64(config.GasUpdaterBlockDelay())

	gethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(logsSub, nil)
	gethClient.On("FilterLogs", mock.Anything, mock.Anything).Maybe().Return([]models.Log{}, nil)

	// Initial tx attempt sent
	gethClient.On("SendTransaction", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			tx, ok := args.Get(1).(*types.Transaction)
			require.True(t, ok)
			gethClient.On("TransactionReceipt", mock.Anything, mock.Anything).
				Return(&types.Receipt{TxHash: tx.Hash(), BlockNumber: big.NewInt(confirmed)}, nil)
		}).
		Return(nil).Once()

	rpcClient.On("CallContext", mock.Anything, mock.Anything, "eth_getBlockByNumber", mock.Anything, false).
		Run(func(args mock.Arguments) {
			head := args.Get(1).(**models.Head)
			*head = cltest.Head(inLongestChain)
		}).
		Return(nil)

	gethClient.On("BlockByNumber", mock.Anything, big.NewInt(inLongestChain)).Return(cltest.BlockWithTransactions(), nil)

	// Create FM Job, and wait for job run to start because the above criteria initiates a run.
	buffer := cltest.MustReadFile(t, "testdata/flux_monitor_job.json")
	var job models.JobSpec
	err = json.Unmarshal(buffer, &job)
	require.NoError(t, err)
	job.Initiators[0].InitiatorParams.Feeds = cltest.JSONFromString(t, fmt.Sprintf(`["%s"]`, mockServer.URL))
	job.Initiators[0].InitiatorParams.PollTimer.Period = models.MustMakeDuration(15 * time.Second)

	j := cltest.CreateJobSpecViaWeb(t, app, job)
	jrs := cltest.WaitForRuns(t, j, app.Store, 1)
	jr := cltest.WaitForJobRunToPendOutgoingConfirmations(t, app.Store, jrs[0])

	app.EthBroadcaster.Trigger()
	cltest.WaitForEthTxAttemptCount(t, app.Store, 1)

	newHeads := <-chchNewHeads
	newHeads <- cltest.Head(safe)

	// Check the FM price on completed run output
	jr = cltest.WaitForJobRunToComplete(t, app.GetStore(), jr)

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

	linkEarned, err := app.GetStore().LinkEarnedFor(&j)
	require.NoError(t, err)
	assert.Equal(t, app.Store.Config.MinimumContractPayment(), linkEarned)

	gethClient.AssertExpectations(t)
	rpcClient.AssertExpectations(t)
	sub.AssertExpectations(t)
}

func TestIntegration_FluxMonitor_NewRound(t *testing.T) {
	gethClient := new(mocks.GethClient)
	rpcClient := new(mocks.RPCClient)
	sub := new(mocks.Subscription)

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, config,
		eth.NewClientWith(rpcClient, gethClient),
	)
	defer cleanup()

	app.GetStore().Config.Set(orm.EnvVarName("MinRequiredOutgoingConfirmations"), 1)
	minPayment := app.Store.Config.MinimumContractPayment().ToInt().Uint64()
	availableFunds := minPayment * 100

	// Start, connect, and initialize node
	sub.On("Err").Return(nil)
	sub.On("Unsubscribe").Return(nil).Maybe()
	gethClient.On("ChainID", mock.Anything).Return(app.Store.Config.ChainID(), nil)
	gethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(oneETH.ToInt(), nil)
	chchNewHeads := make(chan chan<- *models.Head, 1)
	rpcClient.On("EthSubscribe", mock.Anything, mock.Anything, "newHeads").
		Run(func(args mock.Arguments) { chchNewHeads <- args.Get(1).(chan<- *models.Head) }).
		Return(sub, nil)

	err := app.StartAndConnect()
	require.NoError(t, err)

	gethClient.AssertExpectations(t)
	rpcClient.AssertExpectations(t)
	sub.AssertExpectations(t)

	// Configure fake Eth Node to return 10,000 cents when FM initiates price.
	rpcClient.On("Call", mock.Anything, "eth_call", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			*args.Get(0).(*hexutil.Bytes) = cltest.MakeRoundStateReturnData(2, true, 10000, 7, 0, availableFunds, minPayment, 1)
		}).
		Return(nil)

	// Have price adapter server respond with 100 for price on initialization,
	// NOT enough for deviation.
	priceResponse := `{"data":{"result": 100}}`
	mockServer, assertCalled := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", priceResponse)
	defer assertCalled()

	confirmed := int64(23456)
	safe := confirmed + int64(config.MinRequiredOutgoingConfirmations())
	inLongestChain := safe - int64(config.GasUpdaterBlockDelay())

	// Prepare new rounds logs subscription to be called by new FM job
	chchLogs := make(chan chan<- types.Log, 1)
	gethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) { chchLogs <- args.Get(2).(chan<- types.Log) }).
		Return(sub, nil)

	// Log Broadcaster backfills logs
	rpcClient.On("CallContext", mock.Anything, mock.Anything, "eth_getBlockByNumber", mock.Anything, false).
		Run(func(args mock.Arguments) {
			head := args.Get(1).(**models.Head)
			*head = cltest.Head(1)
		}).
		Return(nil)
	gethClient.On("FilterLogs", mock.Anything, mock.Anything).Return([]models.Log{}, nil)

	// Create FM Job, and ensure no runs because above criteria has no deviation.
	buffer := cltest.MustReadFile(t, "testdata/flux_monitor_job.json")
	var job models.JobSpec
	err = json.Unmarshal(buffer, &job)
	require.NoError(t, err)
	job.Initiators[0].InitiatorParams.Feeds = cltest.JSONFromString(t, fmt.Sprintf(`["%s"]`, mockServer.URL))
	job.Initiators[0].InitiatorParams.PollTimer.Period = models.MustMakeDuration(15 * time.Second)
	job.Initiators[0].InitiatorParams.IdleTimer.Disabled = true
	job.Initiators[0].InitiatorParams.IdleTimer.Duration = models.MustMakeDuration(0)

	j := cltest.CreateJobSpecViaWeb(t, app, job)
	_ = cltest.AssertRunsStays(t, j, app.Store, 0)

	gethClient.AssertExpectations(t)
	rpcClient.AssertExpectations(t)
	sub.AssertExpectations(t)

	// Send a NewRound log event to trigger a run.
	log := cltest.LogFromFixture(t, "testdata/new_round_log.json")
	log.Address = job.Initiators[0].InitiatorParams.Address

	gethClient.On("SendTransaction", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			tx, ok := args.Get(1).(*types.Transaction)
			require.True(t, ok)
			gethClient.On("TransactionReceipt", mock.Anything, mock.Anything).
				Return(&types.Receipt{TxHash: tx.Hash(), BlockNumber: big.NewInt(confirmed)}, nil)
		}).
		Return(nil).Once()

	rpcClient.On("CallContext", mock.Anything, mock.Anything, "eth_getBlockByNumber", mock.Anything, false).
		Run(func(args mock.Arguments) {
			head := args.Get(1).(**models.Head)
			*head = cltest.Head(inLongestChain)
		}).
		Return(nil)

	gethClient.On("BlockByNumber", mock.Anything, big.NewInt(inLongestChain)).Return(cltest.BlockWithTransactions(), nil)

	// Flux Monitor queries FluxAggregator.RoundState()
	rpcClient.On("Call", mock.Anything, "eth_call", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			*args.Get(1).(*hexutil.Bytes) = cltest.MakeRoundStateReturnData(2, true, 10000, 7, 0, availableFunds, minPayment, 1)
		}).
		Return(nil)

	newRounds := <-chchLogs
	newRounds <- log

	jrs := cltest.WaitForRuns(t, j, app.Store, 1)
	_ = cltest.WaitForJobRunToPendOutgoingConfirmations(t, app.Store, jrs[0])
	app.EthBroadcaster.Trigger()
	cltest.WaitForEthTxAttemptCount(t, app.Store, 1)

	newHeads := <-chchNewHeads
	newHeads <- cltest.Head(safe)
	_ = cltest.WaitForJobRunToComplete(t, app.Store, jrs[0])
	linkEarned, err := app.GetStore().LinkEarnedFor(&j)
	require.NoError(t, err)
	assert.Equal(t, app.Store.Config.MinimumContractPayment(), linkEarned)

	gethClient.AssertExpectations(t)
	rpcClient.AssertExpectations(t)
	sub.AssertExpectations(t)
}

func TestIntegration_MultiwordV1(t *testing.T) {
	gethClient := new(mocks.GethClient)
	rpcClient := new(mocks.RPCClient)
	sub := new(mocks.Subscription)

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, config,
		eth.NewClientWith(rpcClient, gethClient),
	)
	defer cleanup()
	app.Config.Set(orm.EnvVarName("DefaultHTTPAllowUnrestrictedNetworkAccess"), true)
	confirmed := int64(23456)
	safe := confirmed + int64(config.MinRequiredOutgoingConfirmations())
	inLongestChain := safe - int64(config.GasUpdaterBlockDelay())

	sub.On("Err").Return(nil)
	sub.On("Unsubscribe").Return(nil).Maybe()
	gethClient.On("ChainID", mock.Anything).Return(app.Store.Config.ChainID(), nil)
	chchNewHeads := make(chan chan<- *models.Head, 1)
	rpcClient.On("EthSubscribe", mock.Anything, mock.Anything, "newHeads").
		Run(func(args mock.Arguments) { chchNewHeads <- args.Get(1).(chan<- *models.Head) }).
		Return(sub, nil)
	gethClient.On("SendTransaction", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			tx, ok := args.Get(1).(*types.Transaction)
			require.True(t, ok)
			// Should expect the function arguments to be
			// {reqID,bid,ask}={bytes32,bytes32,bytes32}
			bytes32, _ := abi.NewType("bytes32", "", nil)
			var a abi.Arguments = []abi.Argument{{Type: bytes32}, {Type: bytes32}, {Type: bytes32}}
			args, err := a.UnpackValues(tx.Data()[4:])
			require.NoError(t, err)
			assert.Equal(t, cltest.MustHexDecode32ByteString("0000000000000000000000000000000000000000000000000000000000000002"), args[0])
			assert.Equal(t, cltest.MustHexDecode32ByteString("3130302e31000000000000000000000000000000000000000000000000000000"), args[1])
			assert.Equal(t, cltest.MustHexDecode32ByteString("3130302e31350000000000000000000000000000000000000000000000000000"), args[2])
			gethClient.On("TransactionReceipt", mock.Anything, mock.Anything).
				Return(&types.Receipt{TxHash: tx.Hash(), BlockNumber: big.NewInt(confirmed)}, nil)
		}).
		Return(nil).Once()
	rpcClient.On("CallContext", mock.Anything, mock.Anything, "eth_getBlockByNumber", mock.Anything, false).
		Run(func(args mock.Arguments) {
			head := args.Get(1).(**models.Head)
			*head = cltest.Head(inLongestChain)
		}).
		Return(nil)

	gethClient.On("BlockByNumber", mock.Anything, big.NewInt(inLongestChain)).
		Return(cltest.BlockWithTransactions(), nil)

	err := app.StartAndConnect()
	require.NoError(t, err)
	priceResponse := `{"bid": 100.10, "ask": 100.15}`
	mockServer, assertCalled := cltest.NewHTTPMockServer(t, http.StatusOK, "GET", priceResponse)
	defer assertCalled()
	spec := string(cltest.MustReadFile(t, "fixtures/web/multi_word_v1.json"))
	spec = strings.Replace(spec, "https://bitstamp.net/api/ticker/", mockServer.URL, 2)
	j := cltest.CreateSpecViaWeb(t, app, spec)
	jr := cltest.CreateJobRunViaWeb(t, app, j)
	_ = cltest.WaitForJobRunStatus(t, app.Store, jr, models.RunStatusPendingOutgoingConfirmations)
	app.EthBroadcaster.Trigger()
	cltest.WaitForEthTxAttemptCount(t, app.Store, 1)

	// Feed the subscriber a block head so the transaction completes.
	newHeads := <-chchNewHeads
	newHeads <- cltest.Head(safe)
	// Job should complete successfully.
	_ = cltest.WaitForJobRunToComplete(t, app.Store, jr)
	jr2, err := app.Store.ORM.FindJobRun(jr.ID)
	assert.Equal(t, 9, len(jr2.TaskRuns))
	// We expect 2 results collected, the bid and ask
	assert.Equal(t, 2, len(jr2.TaskRuns[8].Result.Data.Get(models.ResultCollectionKey).Array()))
}

func TestIntegration_MultiwordV1_Sim(t *testing.T) {
	//rpcClient := new(mocks.RPCClient)
	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	key, err := crypto.GenerateKey()
	require.NoError(t, err, "failed to generate ethereum identity")
	user := bind.NewKeyedTransactor(key)
	sb := new(big.Int)
	sb, _ = sb.SetString("100000000000000000000", 10)
	genesisData := core.GenesisAlloc{
		user.From:  {Balance: sb}, // 1 eth
	}
	gasLimit := goEthereumEth.DefaultConfig.Miner.GasCeil * 2
	b := backends.NewSimulatedBackend(genesisData, gasLimit)
	linkTokenAddress, _, linkToken, err := linktoken.DeployLinkToken(user, b)
	require.NoError(t, err)
	b.Commit()

	a, err := linkToken.BalanceOf(nil, user.From)
	require.NoError(t, err)
	t.Log(a)
	b.Commit()
	operatorAddress, _, _, err := operator.DeployOperator(user, b, linkTokenAddress, user.From)
	require.NoError(t, err)
	b.Commit()
	config.Set("OPERATOR_CONTRACT_ADDRESS", operatorAddress.String())
	app, cleanup := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, b)
	app.Config.Set("ETH_HEAD_TRACKER_MAX_BUFFER_SIZE", 100)
	app.Config.Set("OPERATOR_CONTRACT_ADDRESS", operatorAddress.String())
	t.Log("min pay", app.Config.MinimumContractPayment().String())

	tx := types.NewTransaction(2, app.Store.KeyStore.Accounts()[0].Address, big.NewInt(1000000000000000000), 21000, big.NewInt(1), nil)
	signedTx, err := user.Signer(types.HomesteadSigner{}, user.From, tx)
	require.NoError(t, err)
	err = b.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)
	b.Commit()

	err = app.StartAndConnect()
	require.NoError(t, err)
	spec := string(cltest.MustReadFile(t, "testdata/multiword_v1.json"))
	//spec = strings.Replace(spec, "0x794240253727e8d030B09223a9DD13D6dA323401", operatorAddress.String(), 1)
	j := cltest.CreateSpecViaWeb(t, app, spec)
	var specID [32]byte
	by, err := hex.DecodeString(j.ID.String())
	require.NoError(t, err)
	copy(specID[:], by[:])
	t.Log("spec ID", specID, "spec", spec)
	consumerAddress, _, consumer, err := multiwordconsumer.DeployMultiwordConsumer(user, b, linkTokenAddress, operatorAddress, specID)
	t.Log(operatorAddress.String(), consumerAddress.String())
	require.NoError(t, err)
	b.Commit()
	// The consumer contract needs to have ETH in it to be able to pay
	// fo
	_, err = linkToken.Transfer(user, consumerAddress, big.NewInt(1000))
	require.NoError(t, err)

	user.GasPrice = big.NewInt(1)
	user.GasLimit = 1000000
	tx, err = consumer.RequestMultipleParameters(user, "", big.NewInt(1000))
	require.NoError(t, err)
	b.Commit()
	t.Log("request hash", tx.Hash().String())
	p, err := consumer.First(nil)
	require.NoError(t, err)
	t.Log(p)
	//block, err := b.BlockByNumber(context.Background(), nil)
	//for _, tr := range block.Transactions() {
	//	t.Log(tr.Hash().String())
	//	re, err := b.TransactionReceipt(context.Background(), tr.Hash())
	//	require.NoError(t, err)
	//	t.Logf("%+v", re)
	//}
	//
	//
	//logs2, err := consumer.FilterTest(nil)
	//require.NoError(t, err)
	//b.Commit()
	//t.Log("logs", logs2.Event, logs2.Next())
	//var reqID [32]byte
	//reqID[0] = 0x01
	//logs1, err := consumer.FilterChainlinkRequested(nil, [][32]byte{reqID})
	//require.NoError(t, err)
	//b.Commit()
	//logs1, err := operatorContract.OperatorFilterer.FilterOracleRequest(nil, [][32]byte{specID})
	//require.NoError(t, err)
	//b.Commit()
	//t.Log("logs", logs1.Event)
	//t.Log(operatorContract)
	//logs, err := b.FilterLogs(context.Background(), ethereum.FilterQuery{
	//	nil, nil, nil, []common.Address{operatorAddress}, nil,
	//})
	//require.NoError(t, err)
	//t.Log("logs", logs)
	//logs, err := b.FilterLogs(context.Background(), ethereum.FilterQuery{
	//	nil, nil, nil, nil, nil,
	//})
	//require.NoError(t, err)
	//t.Log("logs", logs)

	//jr := cltest.CreateJobRunViaWeb(t, app, j)
	//_ = cltest.WaitForJobRunStatus(t, app.Store, jr, models.RunStatusPendingOutgoingConfirmations)

	app.EthBroadcaster.Trigger()
	tick := time.NewTicker(100*time.Millisecond)
	defer tick.Stop()
	go func() {
		for {
			select {
			case <-tick.C:
				app.EthBroadcaster.Trigger()
				b.Commit()
			}
		}
	}()
	t.Log("waiting for a run")
	cltest.WaitForRuns(t, j, app.Store, 1)
	app.EthBroadcaster.Trigger()

	jr, err := app.Store.JobRunsFor(j.ID)
	for _, task := range jr[0].TaskRuns {
		t.Log(task)
	}
	t.Log(jr[0].RunRequest.TxHash.String())
	cltest.WaitForEthTxAttemptCount(t, app.Store, 1)
	attempts, _, _ := app.Store.EthTxAttempts(0, 1)
	b.Commit()
	recp, err := b.TransactionReceipt(context.Background(), attempts[0].Hash)
	if err != nil {
		panic(err)
	}
	fmt.Println("status", recp.Status)

	// Job should complete successfully.
	_ = cltest.WaitForJobRunStatus(t, app.Store, jr[0], models.RunStatusCompleted)
	t.Log(jr[0].TaskRuns[len(jr[0].TaskRuns)-1], err)
	//tx2, _, err := b.TransactionByHash(context.Background(),common.HexToHash(jrr.Result.Data.Get("result").String()))
	//require.NoError(t, err)
	//t.Log(tx2.To())
	p2, err := consumer.First(nil)
	require.NoError(t, err)
	t.Log(p2)
	logs, err := b.FilterLogs(context.Background(), ethereum.FilterQuery{
		nil, nil, nil, nil, nil,
	})
	require.NoError(t, err)
	for _, l := range logs {
		t.Log("logs", string(l.Data))
	}
}
