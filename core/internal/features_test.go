package internal_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"

	"github.com/pborman/uuid"

	"github.com/onsi/gomega"

	"github.com/smartcontractkit/chainlink/core/store/models/ocrkey"
	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	"github.com/smartcontractkit/libocr/gethwrappers/testoffchainaggregator"
	"github.com/smartcontractkit/libocr/gethwrappers/testvalidator"
	"github.com/smartcontractkit/libocr/offchainreporting/confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
	"gopkg.in/guregu/null.v4"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	goEthereumEth "github.com/ethereum/go-ethereum/eth"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/multiwordconsumer_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/operator_wrapper"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/static"

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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

var oneETH = assets.Eth(*big.NewInt(1000000000000000000))

func TestIntegration_Scheduler(t *testing.T) {
	t.Parallel()

	rpcClient, gethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		eth.NewClientWith(rpcClient, gethClient),
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

	rpcClient, gethClient, sub, assertMocksCalled := cltest.NewEthMocks(t)
	defer assertMocksCalled()
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

	rpcClient, gethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		eth.NewClientWith(rpcClient, gethClient),
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
	rpcClient, gethClient, sub, assertMockCalls := cltest.NewEthMocks(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		eth.NewClientWith(rpcClient, gethClient),
	)
	defer cleanup()

	sub.On("Err").Return(nil).Maybe()
	sub.On("Unsubscribe").Return(nil).Maybe()
	gethClient.On("ChainID", mock.Anything).Return(app.Store.Config.ChainID(), nil)
	gethClient.On("FilterLogs", mock.Anything, mock.Anything).Maybe().Return([]models.Log{}, nil)
	rpcClient.On("EthSubscribe", mock.Anything, mock.Anything, "newHeads").Return(sub, nil)
	logsCh := cltest.MockSubscribeToLogsCh(gethClient, sub)
	gethClient.On("TransactionReceipt", mock.Anything, mock.Anything).
		Return(&types.Receipt{}, nil)
	require.NoError(t, app.StartAndConnect())

	j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/eth_log_job.json")
	address := common.HexToAddress("0x3cCad4715152693fE3BC4460591e3D3Fbd071b42")
	initr := j.Initiators[0]
	assert.Equal(t, models.InitiatorEthLog, initr.Type)
	assert.Equal(t, address, initr.Address)

	logs := <-logsCh
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

			rpcClient, gethClient, sub, assertMockCalls := cltest.NewEthMocks(t)
			defer assertMockCalls()
			app, cleanup := cltest.NewApplication(t,
				eth.NewClientWith(rpcClient, gethClient),
			)
			defer cleanup()
			sub.On("Err").Return(nil).Maybe()
			sub.On("Unsubscribe").Return(nil).Maybe()
			rpcClient.On("CallContext", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			gethClient.On("ChainID", mock.Anything).Return(app.Store.Config.ChainID(), nil)
			gethClient.On("FilterLogs", mock.Anything, mock.Anything).Maybe().Return([]models.Log{}, nil)
			logsCh := cltest.MockSubscribeToLogsCh(gethClient, sub)
			newHeads := make(chan<- *models.Head, 10)
			rpcClient.On("EthSubscribe", mock.Anything, mock.Anything, "newHeads").
				Run(func(args mock.Arguments) {
					newHeads = args.Get(1).(chan<- *models.Head)
				}).
				Return(sub, nil)
			require.NoError(t, app.StartAndConnect())
			j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/runlog_noop_job.json")
			requiredConfs := int64(100)
			initr := j.Initiators[0]
			assert.Equal(t, models.InitiatorRunLog, initr.Type)

			creationHeight := int64(1)
			runlog := cltest.NewRunLog(t, j.ID, cltest.NewAddress(), cltest.NewAddress(), int(creationHeight), `{}`)
			runlog.BlockHash = test.logBlockHash
			logs := <-logsCh
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
			gethClient.On("BlockByNumber", mock.Anything, mock.Anything).Return(&types.Block{}, nil)
			gethClient.On("TransactionReceipt", mock.Anything, mock.Anything).
				Return(confirmedReceipt, nil)

			app.EthBroadcaster.Trigger()
			jr = cltest.WaitForJobRunStatus(t, app.Store, jr, test.wantStatus)
			assert.True(t, jr.FinishedAt.Valid)
			assert.Equal(t, int64(requiredConfs), int64(jr.TaskRuns[0].ObservedIncomingConfirmations.Uint32))
		})
	}
}

func TestIntegration_StartAt(t *testing.T) {
	t.Parallel()

	rpcClient, gethClient, _, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		eth.NewClientWith(rpcClient, gethClient),
	)
	defer cleanup()
	require.NoError(t, app.Start())

	j := cltest.FixtureCreateJobViaWeb(t, app, "fixtures/web/start_at_job.json")
	startAt := cltest.ParseISO8601(t, "1970-01-01T00:00:00.000Z")
	assert.Equal(t, startAt, j.StartAt.Time)

	jr := cltest.CreateJobRunViaWeb(t, app, j)
	cltest.WaitForJobRunToComplete(t, app.Store, jr)
}

func TestIntegration_ExternalAdapter_RunLogInitiated(t *testing.T) {
	t.Parallel()
	rpcClient, gethClient, sub, assertMockCalls := cltest.NewEthMocks(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		eth.NewClientWith(rpcClient, gethClient),
	)
	defer cleanup()

	gethClient.On("ChainID", mock.Anything).Return(app.Config.ChainID(), nil)
	sub.On("Err").Return(nil)
	sub.On("Unsubscribe").Return(nil)
	newHeadsCh := make(chan chan<- *models.Head, 1)
	logsCh := cltest.MockSubscribeToLogsCh(gethClient, sub)
	rpcClient.On("CallContext", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	rpcClient.On("EthSubscribe", mock.Anything, mock.Anything, "newHeads").
		Run(func(args mock.Arguments) {
			newHeadsCh <- args.Get(1).(chan<- *models.Head)
		}).
		Return(sub, nil)
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
	logs := <-logsCh
	logs <- runlog
	jr := cltest.WaitForRuns(t, j, app.Store, 1)[0]
	cltest.WaitForJobRunToPendIncomingConfirmations(t, app.Store, jr)

	gethClient.On("BlockByNumber", mock.Anything, mock.Anything).Return(types.NewBlockWithHeader(&types.Header{
		Number: big.NewInt(int64(logBlockNumber + 8)),
	}), nil) // Gas updater checks the block by number.
	newHeads := <-newHeadsCh
	newHeads <- cltest.Head(logBlockNumber + 8)
	cltest.WaitForJobRunToPendIncomingConfirmations(t, app.Store, jr)

	confirmedReceipt := &types.Receipt{
		TxHash:      runlog.TxHash,
		BlockHash:   runlog.BlockHash,
		BlockNumber: big.NewInt(int64(logBlockNumber)),
	}

	gethClient.On("BlockByNumber", mock.Anything, mock.Anything).Return(types.NewBlockWithHeader(&types.Header{
		Number: big.NewInt(int64(logBlockNumber + 9)),
	}), nil)
	gethClient.On("TransactionReceipt", mock.Anything, mock.Anything).
		Return(confirmedReceipt, nil)

	newHeads <- cltest.Head(logBlockNumber + 9)
	jr = cltest.SendBlocksUntilComplete(t, app.Store, jr, newHeads, int64(logBlockNumber+9))

	tr := jr.TaskRuns[0]
	assert.Equal(t, "randomnumber", tr.TaskSpec.Type.String())
	value := cltest.MustResultString(t, tr.Result)
	assert.Equal(t, eaValue, value)
	res := tr.Result.Data.Get("extra")
	assert.Equal(t, eaExtra, res.String())
}

// This test ensures that the response body of an external adapter are supplied
// as params to the successive task
func TestIntegration_ExternalAdapter_Copy(t *testing.T) {
	t.Parallel()

	rpcClient, gethClient, _, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		eth.NewClientWith(rpcClient, gethClient),
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

	rpcClient, gethClient, _, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		eth.NewClientWith(rpcClient, gethClient),
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

	rpcClient, gethClient, sub, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		eth.NewClientWith(rpcClient, gethClient),
	)
	defer cleanup()

	logsCh := cltest.MockSubscribeToLogsCh(gethClient, sub)
	gethClient.On("TransactionReceipt", mock.Anything, mock.Anything).
		Return(&types.Receipt{}, nil)

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

	logs := <-logsCh
	logs <- log

	jobRuns := cltest.WaitForRuns(t, j, app.Store, 1)
	cltest.WaitForJobRunToComplete(t, app.Store, jobRuns[0])
}

func TestIntegration_MultiplierInt256(t *testing.T) {
	rpcClient, gethClient, _, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		eth.NewClientWith(rpcClient, gethClient),
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
	rpcClient, gethClient, _, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		eth.NewClientWith(rpcClient, gethClient),
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
	rpcClient, gethClient, _, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplicationWithConfig(t,
		config,
		eth.NewClientWith(rpcClient, gethClient),
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
		message = <-wsserver.ReceivedText
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
	rpcClient, gethClient, _, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		eth.NewClientWith(rpcClient, gethClient),
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

	rpcClient, gethClient, _, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		eth.NewClientWith(rpcClient, gethClient),
		services.NewExternalInitiatorManager(),
	)
	defer cleanup()
	require.NoError(t, app.Start())

	exInitr := struct {
		Header http.Header
		Body   services.JobSpecNotice
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
		exInitr.Header.Get(static.ExternalInitiatorAccessKeyHeader),
	)
	assert.Equal(t,
		eip.OutgoingSecret,
		exInitr.Header.Get(static.ExternalInitiatorSecretHeader),
	)
	expected := services.JobSpecNotice{
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

	rpcClient, gethClient, _, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		eth.NewClientWith(rpcClient, gethClient),
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
	rpcClient, gethClient, _, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		eth.NewClientWith(rpcClient, gethClient),
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

	_, address := cltest.MustAddRandomKeyToKeystore(t, app.Store, 0)

	kst := new(mocks.KeyStoreInterface)
	kst.On("HasAccountWithAddress", address).Return(true)
	kst.On("GetAccountByAddress", mock.Anything).Maybe().Return(accounts.Account{}, nil)
	kst.On("SignTx", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(&types.Transaction{}, nil)
	kst.On("Accounts").Return([]accounts.Account{})

	app.Store.KeyStore = kst

	// Start, connect, and initialize node
	sub.On("Err").Return(nil).Maybe()
	sub.On("Unsubscribe").Return(nil).Maybe()
	gethClient.On("ChainID", mock.Anything).Return(app.Store.Config.ChainID(), nil)
	gethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(oneETH.ToInt(), nil)
	newHeads := make(chan<- *models.Head, 1)
	rpcClient.On("EthSubscribe", mock.Anything, mock.Anything, "newHeads").
		Run(func(args mock.Arguments) { newHeads = args.Get(1).(chan<- *models.Head) }).
		Return(sub, nil)

	logsSub := new(mocks.Subscription)
	logsSub.On("Err").Return(nil)
	logsSub.On("Unsubscribe").Return(nil).Maybe()

	err := app.StartAndConnect()
	require.NoError(t, err)

	gethClient.AssertExpectations(t)
	sub.AssertExpectations(t)

	// Configure fake Eth Node to return 10,000 cents when FM initiates price.
	minPayment := app.Store.Config.MinimumContractPayment().ToInt().Uint64()
	availableFunds := minPayment * 100

	// getOracles()
	getOraclesResult, err := cltest.GenericEncode([]string{"address[]"}, []common.Address{address})
	require.NoError(t, err)
	cltest.MockFluxAggCall(gethClient, cltest.FluxAggAddress, "getOracles").
		Return(getOraclesResult, nil).Once()

	// latestRoundData()
	lrdTypes := []string{"uint80", "int256", "uint256", "uint256", "uint80"}
	latestRoundDataResult, err := cltest.GenericEncode(
		lrdTypes, big.NewInt(2), big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1),
	)
	require.NoError(t, err)
	cltest.MockFluxAggCall(gethClient, cltest.FluxAggAddress, "latestRoundData").
		Return(latestRoundDataResult, nil).Once()

	// oracleRoundState()
	result := cltest.MakeRoundStateReturnData(2, true, 10000, 7, 0, availableFunds, minPayment, 1)
	cltest.MockFluxAggCall(gethClient, cltest.FluxAggAddress, "oracleRoundState").
		Return(result, nil).Twice()

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
	newHeadsCh := make(chan chan<- *models.Head, 1)
	rpcClient.On("EthSubscribe", mock.Anything, mock.Anything, "newHeads").
		Run(func(args mock.Arguments) { newHeadsCh <- args.Get(1).(chan<- *models.Head) }).
		Return(sub, nil)

	err := app.StartAndConnect()
	require.NoError(t, err)

	gethClient.AssertExpectations(t)
	rpcClient.AssertExpectations(t)
	sub.AssertExpectations(t)

	// Configure fake Eth Node to return 10,000 cents when FM initiates price.
	getOraclesResult, err := cltest.GenericEncode([]string{"address[]"}, []common.Address{})
	require.NoError(t, err)
	cltest.MockFluxAggCall(gethClient, cltest.FluxAggAddress, "getOracles").
		Return(getOraclesResult, nil).Once()

	cltest.MockFluxAggCall(gethClient, cltest.FluxAggAddress, "latestRoundData").
		Return(nil, errors.New("first round")).Once()

	result := cltest.MakeRoundStateReturnData(2, true, 10000, 7, 0, availableFunds, minPayment, 1)
	cltest.MockFluxAggCall(gethClient, cltest.FluxAggAddress, "oracleRoundState").
		Return(result, nil)

	// Have price adapter server respond with 100 for price on initialization,
	// NOT enough for deviation.
	priceResponse := `{"data":{"result": 100}}`
	mockServer, assertCalled := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", priceResponse)
	defer assertCalled()

	confirmed := int64(23456)
	safe := confirmed + int64(config.MinRequiredOutgoingConfirmations())
	inLongestChain := safe - int64(config.GasUpdaterBlockDelay())

	// Prepare new rounds logs subscription to be called by new FM job
	logs := make(chan<- types.Log, 1)
	gethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) { logs = args.Get(2).(chan<- types.Log) }).
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

	logs <- log

	jrs := cltest.WaitForRuns(t, j, app.Store, 1)
	_ = cltest.WaitForJobRunToPendOutgoingConfirmations(t, app.Store, jrs[0])
	app.EthBroadcaster.Trigger()
	cltest.WaitForEthTxAttemptCount(t, app.Store, 1)

	newHeads := <-newHeadsCh
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
	headsCh := make(chan chan<- *models.Head, 1)
	rpcClient.On("EthSubscribe", mock.Anything, mock.Anything, "newHeads").
		Run(func(args mock.Arguments) { headsCh <- args.Get(1).(chan<- *models.Head) }).
		Return(sub, nil)
	gethClient.On("SendTransaction", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			tx, ok := args.Get(1).(*types.Transaction)
			require.True(t, ok)
			assert.Equal(t, cltest.MustHexDecodeString(
				"0000000000000000000000000000000000000000000000000000000000000001"+ // reqID
					"00000000000000000000000000000000000000000000000000000000000000c0"+ // fixed offset
					"0000000000000000000000000000000000000000000000000000000000000060"+ // length 3 * 32
					"0000000000000000000000000000000000000000000000000000000000000001"+ // reqID
					"3130302e31000000000000000000000000000000000000000000000000000000"+ // bid
					"3130302e31350000000000000000000000000000000000000000000000000000"), // ask
				tx.Data()[4:])
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
	spec := string(cltest.MustReadFile(t, "testdata/multiword_v1_web.json"))
	spec = strings.Replace(spec, "https://bitstamp.net/api/ticker/", mockServer.URL, 2)
	j := cltest.CreateSpecViaWeb(t, app, spec)
	jr := cltest.CreateJobRunViaWeb(t, app, j)
	_ = cltest.WaitForJobRunStatus(t, app.Store, jr, models.RunStatusPendingOutgoingConfirmations)
	app.EthBroadcaster.Trigger()
	cltest.WaitForEthTxAttemptCount(t, app.Store, 1)

	// Feed the subscriber a block head so the transaction completes.
	heads := <-headsCh
	heads <- cltest.Head(safe)
	// Job should complete successfully.
	_ = cltest.WaitForJobRunToComplete(t, app.Store, jr)
	jr2, err := app.Store.ORM.FindJobRun(jr.ID)
	require.NoError(t, err)
	assert.Equal(t, 9, len(jr2.TaskRuns))
	// We expect 2 results collected, the bid and ask
	assert.Equal(t, 2, len(jr2.TaskRuns[8].Result.Data.Get(models.ResultCollectionKey).Array()))
}

func assertPrices(t *testing.T, usd, eur, jpy []byte, consumer *multiwordconsumer_wrapper.MultiWordConsumer) {
	var tmp [32]byte
	copy(tmp[:], usd)
	haveUsd, err := consumer.Usd(nil)
	require.NoError(t, err)
	assert.Equal(t, tmp[:], haveUsd[:])
	copy(tmp[:], eur)
	haveEur, err := consumer.Eur(nil)
	require.NoError(t, err)
	assert.Equal(t, tmp[:], haveEur[:])
	copy(tmp[:], jpy)
	haveJpy, err := consumer.Jpy(nil)
	require.NoError(t, err)
	assert.Equal(t, tmp[:], haveJpy[:])
}

func setupMultiWordContracts(t *testing.T) (*bind.TransactOpts, common.Address, *link_token_interface.LinkToken, *multiwordconsumer_wrapper.MultiWordConsumer, *operator_wrapper.Operator, *backends.SimulatedBackend) {
	key, err := crypto.GenerateKey()
	require.NoError(t, err, "failed to generate ethereum identity")
	user := cltest.MustNewSimulatedBackendKeyedTransactor(t, key)
	sb := new(big.Int)
	sb, _ = sb.SetString("100000000000000000000", 10)
	genesisData := core.GenesisAlloc{
		user.From: {Balance: sb}, // 1 eth
	}
	gasLimit := goEthereumEth.DefaultConfig.Miner.GasCeil * 2
	b := backends.NewSimulatedBackend(genesisData, gasLimit)
	linkTokenAddress, _, linkContract, err := link_token_interface.DeployLinkToken(user, b)
	require.NoError(t, err)
	b.Commit()

	operatorAddress, _, operatorContract, err := operator_wrapper.DeployOperator(user, b, linkTokenAddress, user.From)
	require.NoError(t, err)
	b.Commit()

	var empty [32]byte
	consumerAddress, _, consumerContract, err := multiwordconsumer_wrapper.DeployMultiWordConsumer(user, b, linkTokenAddress, operatorAddress, empty)
	require.NoError(t, err)
	b.Commit()

	// The consumer contract needs to have link in it to be able to pay
	// for the data request.
	_, err = linkContract.Transfer(user, consumerAddress, big.NewInt(1000))
	require.NoError(t, err)
	return user, consumerAddress, linkContract, consumerContract, operatorContract, b
}

func TestIntegration_MultiwordV1_Sim(t *testing.T) {
	// Simulate a consumer contract calling to obtain ETH quotes in 3 different currencies
	// in a single callback.
	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	user, _, _, consumerContract, operatorContract, b := setupMultiWordContracts(t)
	app, cleanup := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, b)
	defer cleanup()
	app.Config.Set("ETH_HEAD_TRACKER_MAX_BUFFER_SIZE", 100)
	app.Config.Set("MIN_OUTGOING_CONFIRMATIONS", 1)

	_, err := operatorContract.SetAuthorizedSender(user, app.Store.KeyStore.Accounts()[0].Address, true)
	require.NoError(t, err)
	b.Commit()

	// Fund node account with ETH.
	n, err := b.NonceAt(context.Background(), user.From, nil)
	require.NoError(t, err)
	tx := types.NewTransaction(n, app.Store.KeyStore.Accounts()[0].Address, big.NewInt(1000000000000000000), 21000, big.NewInt(1), nil)
	signedTx, err := user.Signer(user.From, tx)
	require.NoError(t, err)
	err = b.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)
	b.Commit()

	err = app.StartAndConnect()
	require.NoError(t, err)

	var call int64
	response := func() string {
		defer func() { atomic.AddInt64(&call, 1) }()
		switch call {
		case 0:
			return `{"USD": 614.64}`
		case 1:
			return `{"EUR": 507.07}`
		case 2:
			return `{"JPY":63818.86}`
		}
		require.Fail(t, "only 3 calls expected")
		return ""
	}
	mockServer := cltest.NewHTTPMockServerWithAlterableResponse(t, response)
	spec := string(cltest.MustReadFile(t, "testdata/multiword_v1_runlog.json"))
	spec = strings.Replace(spec, "{url}", mockServer.URL, 1)
	spec = strings.Replace(spec, "{url}", mockServer.URL, 1)
	spec = strings.Replace(spec, "{url}", mockServer.URL, 1)
	j := cltest.CreateSpecViaWeb(t, app, spec)

	var specID [32]byte
	by, err := hex.DecodeString(j.ID.String())
	require.NoError(t, err)
	copy(specID[:], by[:])
	_, err = consumerContract.SetSpecID(user, specID)
	require.NoError(t, err)

	user.GasPrice = big.NewInt(1)
	user.GasLimit = 1000000
	_, err = consumerContract.RequestMultipleParameters(user, "", big.NewInt(1000))
	require.NoError(t, err)
	b.Commit()

	var empty [32]byte
	assertPrices(t, empty[:], empty[:], empty[:], consumerContract)

	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()
	go func() {
		for range tick.C {
			app.EthBroadcaster.Trigger()
			b.Commit()
		}
	}()
	cltest.WaitForRuns(t, j, app.Store, 1)
	jr, err := app.Store.JobRunsFor(j.ID)
	require.NoError(t, err)
	cltest.WaitForEthTxAttemptCount(t, app.Store, 1)

	// Job should complete successfully.
	_ = cltest.WaitForJobRunStatus(t, app.Store, jr[0], models.RunStatusCompleted)
	assertPrices(t, []byte("614.64"), []byte("507.07"), []byte("63818.86"), consumerContract)
}

func setupOCRContracts(t *testing.T) (*bind.TransactOpts, *backends.SimulatedBackend, common.Address, *offchainaggregator.OffchainAggregator) {
	key, err := crypto.GenerateKey()
	require.NoError(t, err, "failed to generate ethereum identity")
	owner := cltest.MustNewSimulatedBackendKeyedTransactor(t, key)
	sb := new(big.Int)
	sb, _ = sb.SetString("100000000000000000000", 10) // 1 eth
	genesisData := core.GenesisAlloc{
		owner.From: {Balance: sb},
	}
	gasLimit := goEthereumEth.DefaultConfig.Miner.GasCeil * 2
	b := backends.NewSimulatedBackend(genesisData, gasLimit)
	linkTokenAddress, _, linkContract, err := link_token_interface.DeployLinkToken(owner, b)
	require.NoError(t, err)
	testValidatorAddress, _, _, err := testvalidator.DeployTestValidator(owner, b)
	require.NoError(t, err)
	accessAddress, _, _, err :=
		testoffchainaggregator.DeploySimpleWriteAccessController(owner, b)
	require.NoError(t, err, "failed to deploy test access controller contract")
	b.Commit()

	min, max := new(big.Int), new(big.Int)
	min.Exp(big.NewInt(-2), big.NewInt(191), nil)
	max.Exp(big.NewInt(2), big.NewInt(191), nil)
	max.Sub(max, big.NewInt(1))
	ocrContractAddress, _, ocrContract, err := offchainaggregator.DeployOffchainAggregator(owner, b,
		1000,             // _maximumGasPrice uint32,
		200,              //_reasonableGasPrice uint32,
		3.6e7,            // 3.6e7 microLINK, or 36 LINK
		1e8,              // _linkGweiPerObservation uint32,
		4e8,              // _linkGweiPerTransmission uint32,
		linkTokenAddress, //_link common.Address,
		testValidatorAddress,
		min, // -2**191
		max, // 2**191 - 1
		accessAddress,
		0,
		"TEST")
	require.NoError(t, err)
	_, err = linkContract.Transfer(owner, ocrContractAddress, big.NewInt(1000))
	require.NoError(t, err)
	b.Commit()
	return owner, b, ocrContractAddress, ocrContract
}

func setupNode(t *testing.T, owner *bind.TransactOpts, port int, dbName string, b *backends.SimulatedBackend) (*cltest.TestApplication, string, common.Address, ocrkey.EncryptedKeyBundle, func()) {
	config, _, ormCleanup := cltest.BootstrapThrowawayORM(t, fmt.Sprintf("%s%s", dbName, strings.Replace(uuid.New(), "-", "", -1)), true)
	config.Dialect = orm.DialectPostgresWithoutLock
	app, appCleanup := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, b)
	_, _, err := app.Store.OCRKeyStore.GenerateEncryptedP2PKey()
	require.NoError(t, err)
	p2pIDs := app.Store.OCRKeyStore.DecryptedP2PKeys()
	require.NoError(t, err)
	require.Len(t, p2pIDs, 1)
	peerID := p2pIDs[0].MustGetPeerID().String()

	app.Config.Set("P2P_PEER_ID", peerID)
	app.Config.Set("P2P_LISTEN_PORT", port)
	app.Config.Set("ETH_HEAD_TRACKER_MAX_BUFFER_SIZE", 100)
	app.Config.Set("MIN_OUTGOING_CONFIRMATIONS", 1)
	app.Config.Set("CHAINLINK_DEV", true) // Disables ocr spec validation so we can have fast polling for the test.

	transmitter := app.Store.KeyStore.Accounts()[0].Address

	// Fund the transmitter address with some ETH
	n, err := b.NonceAt(context.Background(), owner.From, nil)
	require.NoError(t, err)

	tx := types.NewTransaction(n, transmitter, big.NewInt(1000000000000000000), 21000, big.NewInt(1), nil)
	signedTx, err := owner.Signer(owner.From, tx)
	require.NoError(t, err)
	err = b.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)
	b.Commit()

	_, kb, err := app.Store.OCRKeyStore.GenerateEncryptedOCRKeyBundle()
	require.NoError(t, err)
	return app, peerID, transmitter, kb, func() {
		ormCleanup()
		appCleanup()
	}
}

func TestIntegration_OCR(t *testing.T) {
	owner, b, ocrContractAddress, ocrContract := setupOCRContracts(t)

	// Note it's plausible these ports could be occupied on a CI machine.
	// May need a port randomize + retry approach if we observe collisions.
	appBootstrap, bootstrapPeerID, _, _, cleanup := setupNode(t, owner, 19999, "bootstrap", b)
	defer cleanup()

	var (
		oracles      []confighelper.OracleIdentity
		transmitters []common.Address
		kbs          []ocrkey.EncryptedKeyBundle
		apps         []*cltest.TestApplication
	)
	for i := 0; i < 4; i++ {
		app, peerID, transmitter, kb, cleanup := setupNode(t, owner, 20000+i, fmt.Sprintf("oracle%d", i), b)
		defer cleanup()
		// We want to quickly poll for the bootstrap node to come up, but if we poll too quickly
		// we'll flood it with messages and slow things down. 5s is about how long it takes the
		// bootstrap node to come up.
		app.Config.Set("OCR_BOOTSTRAP_CHECK_INTERVAL", "5s")

		kbs = append(kbs, kb)
		apps = append(apps, app)
		transmitters = append(transmitters, transmitter)

		oracles = append(oracles, confighelper.OracleIdentity{
			OnChainSigningAddress:           ocrtypes.OnChainSigningAddress(kb.OnChainSigningAddress),
			TransmitAddress:                 transmitter,
			OffchainPublicKey:               ocrtypes.OffchainPublicKey(kb.OffChainPublicKey),
			PeerID:                          peerID,
			SharedSecretEncryptionPublicKey: ocrtypes.SharedSecretEncryptionPublicKey(kb.ConfigPublicKey),
		})
	}

	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	go func() {
		for range tick.C {
			b.Commit()
		}
	}()

	_, err := ocrContract.SetPayees(owner,
		transmitters,
		transmitters,
	)
	require.NoError(t, err)
	signers, transmitters, threshold, encodedConfigVersion, encodedConfig, err := confighelper.ContractSetConfigArgsForIntegrationTest(
		oracles,
		1,
		1000000000/100, // threshold PPB
	)
	require.NoError(t, err)
	_, err = ocrContract.SetConfig(owner,
		signers,
		transmitters,
		threshold,
		encodedConfigVersion,
		encodedConfig,
	)
	require.NoError(t, err)
	b.Commit()

	err = appBootstrap.StartAndConnect()
	require.NoError(t, err)
	defer appBootstrap.Stop()

	ocrJob, err := offchainreporting.ValidatedOracleSpecToml(appBootstrap.Config.Config, fmt.Sprintf(`
type               = "offchainreporting"
schemaVersion      = 1
name               = "boot"
contractAddress    = "%s"
isBootstrapPeer    = true 
`, ocrContractAddress))
	require.NoError(t, err)
	_, err = appBootstrap.AddJobV2(context.Background(), ocrJob, null.NewString("boot", true))
	require.NoError(t, err)

	var jids []int32
	for i := 0; i < 4; i++ {
		err = apps[i].StartAndConnect()
		require.NoError(t, err)
		defer apps[i].Stop()

		mockHTTP, cleanupHTTP := cltest.NewHTTPMockServer(t, http.StatusOK, "GET", `{"data": 10}`)
		defer cleanupHTTP()
		ocrJob, err := offchainreporting.ValidatedOracleSpecToml(apps[i].Config.Config, fmt.Sprintf(`
type               = "offchainreporting"
schemaVersion      = 1
name               = "web oracle spec"
contractAddress    = "%s"
isBootstrapPeer    = false
p2pBootstrapPeers  = [
    "/ip4/127.0.0.1/tcp/19999/p2p/%s"
]
keyBundleID        = "%s"
transmitterAddress = "%s"
observationTimeout = "20s"
contractConfigConfirmations = 1 
contractConfigTrackerPollInterval = "1s"
observationSource = """
    // data source 1
    ds1          [type=http method=GET url="%s"];
    ds1_parse    [type=jsonparse path="data"];
    ds1_multiply [type=multiply times=%d];
	ds1->ds1_parse->ds1_multiply;
"""
`, ocrContractAddress, bootstrapPeerID, kbs[i].ID, transmitters[i], mockHTTP.URL, i))
		require.NoError(t, err)
		jid, err := apps[i].AddJobV2(context.Background(), ocrJob, null.NewString("testocr", true))
		require.NoError(t, err)
		jids = append(jids, jid)
	}

	// Assert that all the OCR jobs get a run with valid values eventually.
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		ic := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			pr := cltest.WaitForPipelineComplete(t, ic, jids[ic], apps[ic].GetJobORM(), 1*time.Minute, 1*time.Second)
			jb, err := pr.Outputs.MarshalJSON()
			require.NoError(t, err)
			assert.Equal(t, []byte(fmt.Sprintf("[\"%d\"]", 10*ic)), jb)
			require.NoError(t, err)
		}()
	}
	wg.Wait()

	// 4 oracles reporting 0, 10, 20, 30. Answer should be 20 (results[4/2]).
	gomega.NewGomegaWithT(t).Eventually(func() string {
		answer, err := ocrContract.LatestAnswer(nil)
		require.NoError(t, err)
		return answer.String()
	}, 5*time.Second, 200*time.Millisecond).Should(gomega.Equal("20"))
}
