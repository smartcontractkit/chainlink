package internal_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/testdata/testspecs"

	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/smartcontractkit/chainlink/core/store/dialects"
	"github.com/smartcontractkit/chainlink/core/web/presenters"

	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/gasupdater"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"

	"github.com/onsi/gomega"

	"github.com/smartcontractkit/chainlink/core/store/models/ocrkey"
	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	"github.com/smartcontractkit/libocr/gethwrappers/testoffchainaggregator"
	"github.com/smartcontractkit/libocr/offchainreporting/confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
	"gopkg.in/guregu/null.v4"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/multiwordconsumer_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/operator_wrapper"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/static"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web"

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

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()
	app.Start()

	j := cltest.FixtureCreateJobViaWeb(t, app, "../testdata/jsonspecs/scheduler_job.json")

	cltest.WaitForRunsAtLeast(t, j, app.Store, 1)

	initr := j.Initiators[0]
	assert.Equal(t, models.InitiatorCron, initr.Type)
	assert.Equal(t, "CRON_TZ=UTC * * * * * *", string(initr.Schedule), "Wrong cron schedule saved")
}

func TestIntegration_HttpRequestWithHeaders(t *testing.T) {
	config, cfgCleanup := cltest.NewConfig(t)
	defer cfgCleanup()
	config.Set("ADMIN_CREDENTIALS_FILE", "")
	config.Set("ETH_HEAD_TRACKER_MAX_BUFFER_SIZE", 99)
	config.Set("ETH_FINALITY_DEPTH", 3)

	ethClient, sub, assertMocksCalled := cltest.NewEthMocks(t)
	t.Cleanup(assertMocksCalled)
	chchNewHeads := make(chan chan<- *models.Head, 1)

	app, appCleanup := cltest.NewApplicationWithConfigAndKey(t, config,
		ethClient,
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

	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) { chchNewHeads <- args.Get(1).(chan<- *models.Head) }).
		Return(sub, nil)

	ethClient.On("HeaderByNumber", mock.Anything, mock.AnythingOfType("*big.Int")).Return(cltest.Head(inLongestChain), nil)
	ethClient.On("Dial", mock.Anything).Return(nil)
	ethClient.On("ChainID", mock.Anything).Return(config.ChainID(), nil)
	ethClient.On("PendingNonceAt", mock.Anything, mock.Anything).Maybe().Return(uint64(0), nil)
	ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(uint64(100), nil)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(oneETH.ToInt(), nil)

	ethClient.On("SendTransaction", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			tx, ok := args.Get(1).(*types.Transaction)
			require.True(t, ok)
			ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
				return len(b) == 1 && cltest.BatchElemMatchesHash(b[0], tx.Hash())
			})).Return(nil).Run(func(args mock.Arguments) {
				elems := args.Get(1).([]rpc.BatchElem)
				elems[0].Result = &bulletprooftxmanager.Receipt{TxHash: tx.Hash(), BlockNumber: big.NewInt(confirmed), BlockHash: cltest.NewHash()}
			})
		}).
		Return(nil).Once()

	sub.On("Err").Return(nil)
	sub.On("Unsubscribe").Return(nil).Maybe()

	assert.NoError(t, app.StartAndConnect())

	newHeads := <-chchNewHeads

	j := cltest.CreateHelloWorldJobViaWeb(t, app, mockServer.URL)
	jr := cltest.WaitForJobRunToPendOutgoingConfirmations(t, app.Store, cltest.CreateJobRunViaWeb(t, app, j))

	triggerAllKeys(t, app)
	cltest.WaitForEthTxAttemptCount(t, app.Store, 1)

	// Do the thing
	newHeads <- cltest.Head(safe)

	// sending another head to make sure EthTx executes after EthConfirmer is done
	newHeads <- cltest.Head(safe + 1)

	cltest.WaitForJobRunToComplete(t, app.Store, jr)
}

func TestIntegration_RunAt(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()
	app.InstantClock()

	require.NoError(t, app.Start())
	j := cltest.FixtureCreateJobViaWeb(t, app, "../testdata/jsonspecs/run_at_job.json")

	initr := j.Initiators[0]
	assert.Equal(t, models.InitiatorRunAt, initr.Type)
	assert.Equal(t, "2018-01-08T18:12:01Z", utils.ISO8601UTC(initr.Time.Time))

	jrs := cltest.WaitForRuns(t, j, app.Store, 1)
	cltest.WaitForJobRunToComplete(t, app.Store, jrs[0])
}

func TestIntegration_EthLog(t *testing.T) {
	t.Parallel()

	ethClient, sub, assertMockCalls := cltest.NewEthMocks(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()

	sub.On("Err").Return(nil).Maybe()
	sub.On("Unsubscribe").Return(nil).Maybe()
	ethClient.On("Dial", mock.Anything).Return(nil)
	ethClient.On("ChainID", mock.Anything).Return(app.Store.Config.ChainID(), nil)
	ethClient.On("FilterLogs", mock.Anything, mock.Anything).Maybe().Return([]types.Log{}, nil)
	b := types.NewBlockWithHeader(&types.Header{
		Number: big.NewInt(100),
	})
	ethClient.On("BlockByNumber", mock.Anything, mock.Anything).Maybe().Return(b, nil)
	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).Return(sub, nil)
	logsCh := cltest.MockSubscribeToLogsCh(ethClient, sub)
	ethClient.On("TransactionReceipt", mock.Anything, mock.Anything).
		Return(&types.Receipt{}, nil)
	require.NoError(t, app.StartAndConnect())

	j := cltest.FixtureCreateJobViaWeb(t, app, "../testdata/jsonspecs/eth_log_job.json")
	address := common.HexToAddress("0x3cCad4715152693fE3BC4460591e3D3Fbd071b42")
	initr := j.Initiators[0]
	assert.Equal(t, models.InitiatorEthLog, initr.Type)
	assert.Equal(t, address, initr.Address)

	logs := <-logsCh
	logs <- cltest.LogFromFixture(t, "../testdata/jsonrpc/requestLog0original.json")
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
			t.Cleanup(cfgCleanup)
			config.Set("MIN_INCOMING_CONFIRMATIONS", 6)

			ethClient, sub, assertMockCalls := cltest.NewEthMocks(t)
			defer assertMockCalls()
			app, cleanup := cltest.NewApplication(t,
				ethClient,
			)
			t.Cleanup(cleanup)
			sub.On("Err").Return(nil).Maybe()
			sub.On("Unsubscribe").Return(nil).Maybe()
			ethClient.On("Dial", mock.Anything).Return(nil)
			ethClient.On("ChainID", mock.Anything).Return(app.Store.Config.ChainID(), nil)
			ethClient.On("FilterLogs", mock.Anything, mock.Anything).Maybe().Return([]types.Log{}, nil)
			logsCh := cltest.MockSubscribeToLogsCh(ethClient, sub)
			newHeads := make(chan<- *models.Head, 10)
			ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).
				Run(func(args mock.Arguments) {
					newHeads = args.Get(1).(chan<- *models.Head)
				}).
				Return(sub, nil)

			b := types.NewBlockWithHeader(&types.Header{
				Number: big.NewInt(100),
			})
			ethClient.On("BlockByNumber", mock.Anything, mock.Anything).Maybe().Return(b, nil)

			require.NoError(t, app.StartAndConnect())
			j := cltest.FixtureCreateJobViaWeb(t, app, "../testdata/jsonspecs/runlog_noop_job.json")
			requiredConfs := int64(100)
			initr := j.Initiators[0]
			assert.Equal(t, models.InitiatorRunLog, initr.Type)

			creationHeight := int64(1)
			runlog := cltest.NewRunLog(t, models.IDToTopic(j.ID), cltest.NewAddress(), cltest.NewAddress(), int(creationHeight), `{}`)
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

			ethClient.On("TransactionReceipt", mock.Anything, mock.Anything).
				Return(confirmedReceipt, nil)

			triggerAllKeys(t, app)
			jr = cltest.WaitForJobRunStatus(t, app.Store, jr, test.wantStatus)
			assert.True(t, jr.FinishedAt.Valid)
			assert.Equal(t, int64(requiredConfs), int64(jr.TaskRuns[0].ObservedIncomingConfirmations.Uint32))
		})
	}
}

func TestIntegration_RandomnessReorgProtection(t *testing.T) {
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("MIN_INCOMING_CONFIRMATIONS", 6)

	ethClient, sub, assertMockCalls := cltest.NewEthMocks(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	t.Cleanup(cleanup)
	sub.On("Err").Return(nil).Maybe()
	sub.On("Unsubscribe").Return(nil).Maybe()
	ethClient.On("Dial", mock.Anything).Return(nil)
	ethClient.On("ChainID", mock.Anything).Return(app.Store.Config.ChainID(), nil)
	ethClient.On("FilterLogs", mock.Anything, mock.Anything).Maybe().Return([]types.Log{}, nil)
	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).Return(cltest.EmptyMockSubscription(), nil)

	logsCh := cltest.MockSubscribeToLogsCh(ethClient, sub)
	ethClient.On("BlockByNumber", mock.Anything, mock.Anything).Maybe().Return(types.NewBlockWithHeader(&types.Header{
		Number: big.NewInt(100),
	}), nil)

	require.NoError(t, app.StartAndConnect())
	// Fixture values
	sender := common.HexToAddress("0xABA5eDc1a551E55b1A570c0e1f1055e5BE11eca7")
	keyHash := common.HexToHash("0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800")
	jb := cltest.CreateSpecViaWeb(t, app, testspecs.RandomnessJob)
	logs := <-logsCh
	fee := assets.Link(*big.NewInt(100)) // Default min link is 100
	randLog := models.RandomnessRequestLog{
		KeyHash:   keyHash,
		Seed:      big.NewInt(1),
		JobID:     cltest.NewHash(),
		Sender:    sender,
		Fee:       &fee,
		RequestID: cltest.NewHash(),
		Raw:       models.RawRandomnessRequestLog{},
	}
	log := cltest.NewRandomnessRequestLog(t, randLog, sender, 101)
	logs <- log
	runs := cltest.WaitForRuns(t, jb, app.Store, 1)
	cltest.WaitForJobRunToPendIncomingConfirmations(t, app.Store, runs[0])
	assert.Equal(t, uint32(6), runs[0].TaskRuns[0].MinRequiredIncomingConfirmations.Uint32)

	// Same requestID log again should result in a doubling of incoming confs
	log.TxHash = cltest.NewHash()
	log.BlockHash = cltest.NewHash()
	log.BlockNumber = 102
	logs <- log
	runs = cltest.WaitForRuns(t, jb, app.Store, 2)
	cltest.WaitForJobRunToPendIncomingConfirmations(t, app.Store, runs[0])
	assert.Equal(t, uint32(6)*2, runs[0].TaskRuns[0].MinRequiredIncomingConfirmations.Uint32)

	// Same requestID log again should result in a doubling of incoming confs
	log.TxHash = cltest.NewHash()
	log.BlockHash = cltest.NewHash()
	log.BlockNumber = 103
	logs <- log
	runs = cltest.WaitForRuns(t, jb, app.Store, 3)
	cltest.WaitForJobRunToPendIncomingConfirmations(t, app.Store, runs[0])
	assert.Equal(t, uint32(6)*2*2, runs[0].TaskRuns[0].MinRequiredIncomingConfirmations.Uint32)

	// New requestID should be back to original
	randLog.RequestID = cltest.NewHash()
	newReqLog := cltest.NewRandomnessRequestLog(t, randLog, sender, 104)
	logs <- newReqLog
	runs = cltest.WaitForRuns(t, jb, app.Store, 4)
	cltest.WaitForJobRunToPendIncomingConfirmations(t, app.Store, runs[0])
	assert.Equal(t, uint32(6), runs[0].TaskRuns[0].MinRequiredIncomingConfirmations.Uint32)
}

func TestIntegration_StartAt(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	j := cltest.FixtureCreateJobViaWeb(t, app, "../testdata/jsonspecs/start_at_job.json")
	startAt := cltest.ParseISO8601(t, "1970-01-01T00:00:00.000Z")
	assert.Equal(t, startAt, j.StartAt.Time)

	jr := cltest.CreateJobRunViaWeb(t, app, j)
	cltest.WaitForJobRunToComplete(t, app.Store, jr)
}

func TestIntegration_ExternalAdapter_RunLogInitiated(t *testing.T) {
	t.Parallel()

	ethClient, sub, assertMockCalls := cltest.NewEthMocks(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()

	ethClient.On("Dial", mock.Anything).Return(nil)
	ethClient.On("ChainID", mock.Anything).Return(app.Config.ChainID(), nil)
	b := types.NewBlockWithHeader(&types.Header{
		Number: big.NewInt(100),
	})
	ethClient.On("BlockByNumber", mock.Anything, mock.Anything).Maybe().Return(b, nil)
	ethClient.On("FilterLogs", mock.Anything, mock.Anything).Maybe().Return([]types.Log{}, nil)
	sub.On("Err").Return(nil)
	sub.On("Unsubscribe").Return(nil)
	newHeadsCh := make(chan chan<- *models.Head, 1)
	logsCh := cltest.MockSubscribeToLogsCh(ethClient, sub)
	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).
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
	j := cltest.FixtureCreateJobViaWeb(t, app, "../testdata/jsonspecs/log_initiated_bridge_type_job.json")

	logBlockNumber := 1
	runlog := cltest.NewRunLog(t, models.IDToTopic(j.ID), cltest.NewAddress(), cltest.NewAddress(), logBlockNumber, `{}`)
	logs := <-logsCh
	logs <- runlog
	jr := cltest.WaitForRuns(t, j, app.Store, 1)[0]
	cltest.WaitForJobRunToPendIncomingConfirmations(t, app.Store, jr)

	newHeads := <-newHeadsCh
	newHeads <- cltest.Head(logBlockNumber + 8)
	cltest.WaitForJobRunToPendIncomingConfirmations(t, app.Store, jr)

	confirmedReceipt := &types.Receipt{
		TxHash:      runlog.TxHash,
		BlockHash:   runlog.BlockHash,
		BlockNumber: big.NewInt(int64(logBlockNumber)),
	}

	ethClient.On("TransactionReceipt", mock.Anything, mock.Anything).
		Return(confirmedReceipt, nil)

	newHeads <- cltest.Head(logBlockNumber + 9)
	jr = cltest.SendBlocksUntilComplete(t, app.Store, jr, newHeads, int64(logBlockNumber+9), ethClient)

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

	ethClient, _, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
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
	j := cltest.FixtureCreateJobViaWeb(t, app, "../testdata/jsonspecs/bridge_type_copy_job.json")
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

	ethClient, _, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	resource := &presenters.BridgeResource{}
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
			assert.Equal(t, resource.OutgoingToken, token)
		})
	defer cleanup()

	bridgeJSON := fmt.Sprintf(`{"name":"randomNumber","url":"%v"}`, mockServer.URL)
	resource = cltest.CreateBridgeTypeViaWeb(t, app, bridgeJSON)
	j = cltest.FixtureCreateJobViaWeb(t, app, "../testdata/jsonspecs/random_number_bridge_type_job.json")
	jr := cltest.CreateJobRunViaWeb(t, app, j)
	jr = cltest.WaitForJobRunToPendBridge(t, app.Store, jr)

	tr := jr.TaskRuns[0]
	assert.Equal(t, models.RunStatusPendingBridge, tr.Status)
	assert.Equal(t, gjson.Null, tr.Result.Data.Get("result").Type)

	jr = cltest.UpdateJobRunViaWeb(t, app, jr, resource, `{"data":{"result":"100"}}`)
	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)
	tr = jr.TaskRuns[0]
	assert.Equal(t, models.RunStatusCompleted, tr.Status)

	value := cltest.MustResultString(t, tr.Result)
	assert.Equal(t, "100", value)
}

func TestIntegration_WeiWatchers(t *testing.T) {
	t.Parallel()

	ethClient, sub, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()

	logsCh := cltest.MockSubscribeToLogsCh(ethClient, sub)
	ethClient.On("Dial", mock.Anything).Return(nil)
	ethClient.On("TransactionReceipt", mock.Anything, mock.Anything).
		Return(&types.Receipt{}, nil)

	b := types.NewBlockWithHeader(&types.Header{
		Number: big.NewInt(100),
	})
	ethClient.On("BlockByNumber", mock.Anything, mock.Anything).Maybe().Return(b, nil)
	ethClient.On("FilterLogs", mock.Anything, mock.Anything).Maybe().Return([]types.Log{}, nil)

	log := cltest.LogFromFixture(t, "../testdata/jsonrpc/requestLog0original.json")
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
	t.Parallel()

	ethClient, _, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	j := cltest.FixtureCreateJobViaWeb(t, app, "../testdata/jsonspecs/int256_job.json")
	jr := cltest.CreateJobRunViaWeb(t, app, j, `{"result":"-10221.30"}`)
	jr = cltest.WaitForJobRunToComplete(t, app.Store, jr)

	value := cltest.MustResultString(t, jr.Result)
	assert.Equal(t, "0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0674e", value)
}

func TestIntegration_MultiplierUint256(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	j := cltest.FixtureCreateJobViaWeb(t, app, "../testdata/jsonspecs/uint256_job.json")
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
	ethClient, _, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplicationWithConfig(t,
		config,
		ethClient,
	)
	defer cleanup()
	cltest.MustAddRandomKeyToKeystore(t, app.Store)

	app.InstantClock()
	require.NoError(t, app.Start())

	j := cltest.FixtureCreateJobViaWeb(t, app, "../testdata/jsonspecs/run_at_job.json")

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
}

func TestIntegration_SleepAdapter(t *testing.T) {
	t.Parallel()

	sleepSeconds := 4
	ethClient, _, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	app.Config.Set("ENABLE_EXPERIMENTAL_ADAPTERS", "true")
	defer cleanup()
	require.NoError(t, app.Start())

	j := cltest.FixtureCreateJobViaWeb(t, app, "../testdata/jsonspecs/sleep_job.json")

	runInput := fmt.Sprintf("{\"until\": \"%s\"}", time.Now().Local().Add(time.Second*time.Duration(sleepSeconds)))
	jr := cltest.CreateJobRunViaWeb(t, app, j, runInput)

	cltest.WaitForJobRunStatus(t, app.Store, jr, models.RunStatusInProgress)
	cltest.JobRunStays(t, app.Store, jr, models.RunStatusInProgress, 3*time.Second)
	cltest.WaitForJobRunToComplete(t, app.Store, jr)
}

func TestIntegration_ExternalInitiator(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
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

	jobSpec := cltest.FixtureCreateJobViaWeb(t, app, "../testdata/jsonspecs/external_initiator_job.json")
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
	_, err = app.Store.JobRunsFor(jobRun.JobSpecID)
	assert.NoError(t, err)
	cltest.WaitForJobRunToComplete(t, app.Store, jobRun)
}

func TestIntegration_ExternalInitiator_WithoutURL(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
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

	jobSpec := cltest.FixtureCreateJobViaWeb(t, app, "../testdata/jsonspecs/external_initiator_job.json")

	jobRun := cltest.CreateJobRunViaExternalInitiator(t, app, jobSpec, *eia, "")
	_, err = app.Store.JobRunsFor(jobRun.JobSpecID)
	assert.NoError(t, err)
	cltest.WaitForJobRunToComplete(t, app.Store, jobRun)
}

func TestIntegration_ExternalInitiator_WithMultiplyAndBridge(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()

	app, cleanup := cltest.NewApplication(t,
		ethClient,
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

	// HTTP passthrough
	var httpCalled bool
	mockServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			httpCalled = true
			body, err2 := ioutil.ReadAll(r.Body)
			require.NoError(t, err2)
			requestStr := string(body)
			require.Equal(t, `{"result":"4200000000"}`, requestStr)
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, `4200000000`)
		}))
	defer func() {
		mockServer.Close()
		assert.True(t, httpCalled, "expected http server to be called")
	}()

	// Bridge
	var bridgeCalled bool
	bridgeServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bridgeCalled = true
			body, err2 := ioutil.ReadAll(r.Body)
			require.NoError(t, err2)
			json := models.MustParseJSON(body)
			require.Equal(t, `{"result":"4200000000"}`, json.Get("data").Raw)
			w.WriteHeader(http.StatusOK)
			require.NoError(t, err)
			io.WriteString(w, json.String())
		}))
	u, _ := url.Parse(bridgeServer.URL)
	app.Store.CreateBridgeType(&models.BridgeType{
		Name: models.TaskType("custombridge"),
		URL:  models.WebURL(*u),
	})
	defer func() {
		bridgeServer.Close()
		assert.True(t, bridgeCalled, "expected bridge server to be called")
	}()

	buffer := cltest.MustReadFile(t, "../testdata/jsonspecs/external_initiator_job_multiply.json")
	var jobSpec models.JobSpec
	err = json.Unmarshal(buffer, &jobSpec)
	require.NoError(t, err)
	httpParams, err := jobSpec.Tasks[1].Params.Add("post", mockServer.URL)
	require.NoError(t, err)
	jobSpec.Tasks[1].Params = httpParams

	jobSpec = cltest.CreateJobSpecViaWeb(t, app, jobSpec)

	jobRun := cltest.CreateJobRunViaExternalInitiator(t, app, jobSpec, *eia, `{"result": 42}`)
	_, err = app.Store.JobRunsFor(jobRun.JobSpecID)
	assert.NoError(t, err)
	cltest.WaitForJobRunToComplete(t, app.Store, jobRun)

	jobRun, err = app.Store.FindJobRun(jobRun.ID)
	require.NoError(t, err)
	finalResult := jobRun.Result

	assert.Equal(t, `{"result": "4200000000"}`, finalResult.Data.Raw)
}

func TestIntegration_AuthToken(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
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
	ethClient := new(mocks.Client)

	sub := new(mocks.Subscription)

	config, cfgCleanup := cltest.NewConfig(t)
	defer cfgCleanup()
	config.Set("ETH_FINALITY_DEPTH", 3)
	app, appCleanup := cltest.NewApplicationWithConfig(t, config,
		ethClient,
	)
	defer appCleanup()

	_, address := cltest.MustAddRandomKeyToKeystore(t, app.Store)

	// Start, connect, and initialize node
	sub.On("Err").Return(nil).Maybe()
	sub.On("Unsubscribe").Return(nil).Maybe()
	ethClient.On("Dial", mock.Anything).Return(nil)
	ethClient.On("ChainID", mock.Anything).Return(app.Store.Config.ChainID(), nil)
	ethClient.On("PendingNonceAt", mock.Anything, mock.Anything).Maybe().Return(uint64(0), nil)
	ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(uint64(0), nil)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(oneETH.ToInt(), nil)

	newHeads := make(chan<- *models.Head, 1)
	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) { newHeads = args.Get(1).(chan<- *models.Head) }).
		Return(sub, nil)

	logsSub := new(mocks.Subscription)
	logsSub.On("Err").Return(nil)
	logsSub.On("Unsubscribe").Return(nil).Maybe()

	err := app.StartAndConnect()
	require.NoError(t, err)

	ethClient.AssertExpectations(t)
	sub.AssertExpectations(t)

	cltest.MockFluxAggCall(ethClient, cltest.FluxAggAddress, "minSubmissionValue").
		Return(cltest.MustGenericEncode([]string{"uint256"}, big.NewInt(0)), nil).Once()
	cltest.MockFluxAggCall(ethClient, cltest.FluxAggAddress, "maxSubmissionValue").
		Return(cltest.MustGenericEncode([]string{"uint256"}, big.NewInt(10000000)), nil).Once()
	cltest.MockFluxAggCall(ethClient, cltest.FluxAggAddress, "latestRoundData").
		Return(cltest.MustGenericEncode(
			[]string{"uint80", "int256", "uint256", "uint256", "uint80"},
			big.NewInt(2), big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1),
		), nil).Maybe() // Called 3-4 times.

	// Configure fake Eth Node to return 10,000 cents when FM initiates price.
	minPayment := app.Store.Config.MinimumContractPayment().ToInt().Uint64()
	availableFunds := minPayment * 100

	// getOracles()
	getOraclesResult, err := cltest.GenericEncode([]string{"address[]"}, []common.Address{address})
	require.NoError(t, err)
	cltest.MockFluxAggCall(ethClient, cltest.FluxAggAddress, "getOracles").
		Return(getOraclesResult, nil).Once()

	// oracleRoundState()
	result := cltest.MakeRoundStateReturnData(2, true, 10000, 7, 0, availableFunds, minPayment, 1)
	cltest.MockFluxAggCall(ethClient, cltest.FluxAggAddress, "oracleRoundState").
		Return(result, nil).Twice()

	// Have server respond with 102 for price when FM checks external price
	// adapter for deviation. 102 is enough deviation to trigger a job run.
	priceResponse := `{"data":{"result": 102}}`
	mockServer, assertCalled := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", priceResponse)
	defer assertCalled()

	confirmed := int64(23456)
	safe := confirmed + int64(config.MinRequiredOutgoingConfirmations())
	inLongestChain := safe - int64(config.GasUpdaterBlockDelay())

	ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(logsSub, nil)
	ethClient.On("FilterLogs", mock.Anything, mock.Anything).Maybe().Return([]types.Log{}, nil)

	// Initial tx attempt sent
	ethClient.On("SendTransaction", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			tx, ok := args.Get(1).(*types.Transaction)
			require.True(t, ok)
			ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
				return len(b) == 1 && cltest.BatchElemMatchesHash(b[0], tx.Hash())
			})).Return(nil).Run(func(args mock.Arguments) {
				elems := args.Get(1).([]rpc.BatchElem)
				elems[0].Result = &bulletprooftxmanager.Receipt{TxHash: tx.Hash(), BlockNumber: big.NewInt(confirmed), BlockHash: cltest.NewHash()}
			})
		}).
		Return(nil).Once()
	ethClient.On("HeaderByNumber", mock.Anything, mock.AnythingOfType("*big.Int")).Return(cltest.Head(inLongestChain), nil)

	// Create FM Job, and wait for job run to start because the above criteria initiates a run.
	buffer := cltest.MustReadFile(t, "../testdata/jsonspecs/flux_monitor_job.json")
	var job models.JobSpec
	err = json.Unmarshal(buffer, &job)
	require.NoError(t, err)
	job.Initiators[0].InitiatorParams.Feeds = cltest.JSONFromString(t, fmt.Sprintf(`["%s"]`, mockServer.URL))
	job.Initiators[0].InitiatorParams.PollTimer.Period = models.MustMakeDuration(15 * time.Second)

	j := cltest.CreateJobSpecViaWeb(t, app, job)
	jrs := cltest.WaitForRuns(t, j, app.Store, 1)
	jr := cltest.WaitForJobRunToPendOutgoingConfirmations(t, app.Store, jrs[0])

	triggerAllKeys(t, app)
	cltest.WaitForEthTxAttemptCount(t, app.Store, 1)

	// Check the FM price on completed run output
	jr = cltest.SendBlocksUntilComplete(t, app.GetStore(), jr, newHeads, safe, ethClient)

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

	ethClient.AssertExpectations(t)
	sub.AssertExpectations(t)
}

func TestIntegration_FluxMonitor_NewRound(t *testing.T) {
	ethClient := new(mocks.Client)

	sub := new(mocks.Subscription)

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, config,
		ethClient,
	)
	defer cleanup()

	app.GetStore().Config.Set(orm.EnvVarName("MinRequiredOutgoingConfirmations"), 1)
	minPayment := app.Store.Config.MinimumContractPayment().ToInt().Uint64()
	availableFunds := minPayment * 100

	// Start, connect, and initialize node
	ethClient.On("Dial", mock.Anything).Return(nil)
	ethClient.On("ChainID", mock.Anything).Maybe().Return(app.Store.Config.ChainID(), nil)
	ethClient.On("PendingNonceAt", mock.Anything, mock.Anything).Maybe().Return(uint64(0), nil)
	ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(uint64(0), nil)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(oneETH.ToInt(), nil)
	// Log backfill
	ethClient.On("HeaderByNumber", mock.Anything, mock.AnythingOfType("*big.Int")).Return(cltest.Head(0), nil).Maybe()

	newHeadsCh := make(chan chan<- *models.Head, 1)
	ethClientDone := cltest.NewAwaiter()
	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			ethClientDone.ItHappened()
			newHeadsCh <- args.Get(1).(chan<- *models.Head)
		}).
		Return(sub, nil)

	sub.On("Err").Maybe().Return(nil)
	sub.On("Unsubscribe").Maybe().Return(nil)

	err := app.StartAndConnect()
	require.NoError(t, err)

	ethClientDone.AwaitOrFail(t)
	ethClient.AssertExpectations(t)
	ethClient.AssertExpectations(t)
	sub.AssertExpectations(t)

	cltest.MockFluxAggCall(ethClient, cltest.FluxAggAddress, "minSubmissionValue").
		Return(cltest.MustGenericEncode([]string{"uint256"}, big.NewInt(0)), nil).Once()
	cltest.MockFluxAggCall(ethClient, cltest.FluxAggAddress, "maxSubmissionValue").
		Return(cltest.MustGenericEncode([]string{"uint256"}, big.NewInt(10000000)), nil).Once()
	require.NoError(t, err)
	cltest.MockFluxAggCall(ethClient, cltest.FluxAggAddress, "latestRoundData").
		Return(cltest.MustGenericEncode(
			[]string{"uint80", "int256", "uint256", "uint256", "uint80"},
			big.NewInt(2), big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1),
		), nil).Maybe() // Called 3-4 times.

	// Configure fake Eth Node to return 10,000 cents when FM initiates price.
	getOraclesResult, err := cltest.GenericEncode([]string{"address[]"}, []common.Address{})
	require.NoError(t, err)
	cltest.MockFluxAggCall(ethClient, cltest.FluxAggAddress, "getOracles").
		Return(getOraclesResult, nil).Once()

	result := cltest.MakeRoundStateReturnData(2, true, 10000, 7, 0, availableFunds, minPayment, 1)
	cltest.MockFluxAggCall(ethClient, cltest.FluxAggAddress, "oracleRoundState").
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
	ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) { logs = args.Get(2).(chan<- types.Log) }).
		Return(sub, nil)

	// Log Broadcaster backfills logs
	ethClient.On("HeaderByNumber", mock.Anything, mock.AnythingOfType("*big.Int")).Return(cltest.Head(1), nil)
	ethClient.On("FilterLogs", mock.Anything, mock.Anything).Return([]types.Log{}, nil)

	// Create FM Job, and ensure no runs because above criteria has no deviation.
	buffer := cltest.MustReadFile(t, "../testdata/jsonspecs/flux_monitor_job.json")

	var job models.JobSpec
	err = json.Unmarshal(buffer, &job)
	require.NoError(t, err)
	job.Initiators[0].InitiatorParams.Feeds = cltest.JSONFromString(t, fmt.Sprintf(`["%s"]`, mockServer.URL))
	job.Initiators[0].InitiatorParams.PollTimer.Period = models.MustMakeDuration(15 * time.Second)
	job.Initiators[0].InitiatorParams.IdleTimer.Disabled = true
	job.Initiators[0].InitiatorParams.IdleTimer.Duration = models.MustMakeDuration(0)

	j := cltest.CreateJobSpecViaWeb(t, app, job)
	_ = cltest.AssertRunsStays(t, j, app.Store, 0)

	ethClient.AssertExpectations(t)
	ethClient.AssertExpectations(t)
	sub.AssertExpectations(t)

	// Send a NewRound log event to trigger a run.
	log := cltest.LogFromFixture(t, "../testdata/jsonrpc/new_round_log.json")
	log.Address = job.Initiators[0].InitiatorParams.Address

	ethClient.On("SendTransaction", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			tx, ok := args.Get(1).(*types.Transaction)
			require.True(t, ok)
			ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
				return len(b) == 1 && cltest.BatchElemMatchesHash(b[0], tx.Hash())
			})).Return(nil).Run(func(args mock.Arguments) {
				elems := args.Get(1).([]rpc.BatchElem)
				elems[0].Result = &bulletprooftxmanager.Receipt{TxHash: tx.Hash(), BlockNumber: big.NewInt(confirmed), BlockHash: cltest.NewHash()}
			})
		}).
		Return(nil).Once()

	ethClient.On("HeaderByNumber", mock.Anything, mock.AnythingOfType("*big.Int")).Return(cltest.Head(inLongestChain), nil)

	logs <- log

	newHeads := <-newHeadsCh
	newHeads <- cltest.Head(log.BlockNumber)

	jrs := cltest.WaitForRuns(t, j, app.Store, 1)
	_ = cltest.WaitForJobRunToPendOutgoingConfirmations(t, app.Store, jrs[0])
	triggerAllKeys(t, app)
	cltest.WaitForEthTxAttemptCount(t, app.Store, 1)

	_ = cltest.SendBlocksUntilComplete(t, app.Store, jrs[0], newHeads, safe, ethClient)
	linkEarned, err := app.GetStore().LinkEarnedFor(&j)
	require.NoError(t, err)
	assert.Equal(t, app.Store.Config.MinimumContractPayment(), linkEarned)

	ethClient.AssertExpectations(t)
	ethClient.AssertExpectations(t)
	sub.AssertExpectations(t)
}

func TestIntegration_MultiwordV1(t *testing.T) {
	t.Parallel()

	ethClient := new(mocks.Client)

	sub := new(mocks.Subscription)

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, config,
		ethClient,
	)
	defer cleanup()
	app.Config.Set(orm.EnvVarName("DefaultHTTPAllowUnrestrictedNetworkAccess"), true)
	confirmed := int64(23456)
	safe := confirmed + int64(config.MinRequiredOutgoingConfirmations())
	inLongestChain := safe - int64(config.GasUpdaterBlockDelay())

	sub.On("Err").Return(nil)
	sub.On("Unsubscribe").Return(nil).Maybe()

	ethClient.On("Dial", mock.Anything).Return(nil)
	ethClient.On("ChainID", mock.Anything).Return(app.Store.Config.ChainID(), nil)
	ethClient.On("PendingNonceAt", mock.Anything, mock.Anything).Maybe().Return(uint64(0), nil)
	ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(uint64(0), nil)

	headsCh := make(chan chan<- *models.Head, 1)
	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) { headsCh <- args.Get(1).(chan<- *models.Head) }).
		Return(sub, nil)
	ethClient.On("SendTransaction", mock.Anything, mock.Anything).
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
			ethClient.On("TransactionReceipt", mock.Anything, mock.Anything).
				Return(&types.Receipt{TxHash: tx.Hash(), BlockNumber: big.NewInt(confirmed)}, nil)
			ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
				return len(b) == 1 && cltest.BatchElemMatchesHash(b[0], tx.Hash())
			})).Return(nil).Run(func(args mock.Arguments) {
				elems := args.Get(1).([]rpc.BatchElem)
				elems[0].Result = &bulletprooftxmanager.Receipt{TxHash: tx.Hash(), BlockNumber: big.NewInt(confirmed), BlockHash: cltest.NewHash()}
			}).Maybe()
		}).
		Return(nil).Once()
	ethClient.On("HeaderByNumber", mock.Anything, mock.AnythingOfType("*big.Int")).Return(cltest.Head(inLongestChain), nil)

	err := app.StartAndConnect()
	require.NoError(t, err)
	priceResponse := `{"bid": 100.10, "ask": 100.15}`
	mockServer, assertCalled := cltest.NewHTTPMockServer(t, http.StatusOK, "GET", priceResponse)
	defer assertCalled()
	spec := string(cltest.MustReadFile(t, "../testdata/jsonspecs/multiword_v1_web.json"))
	spec = strings.Replace(spec, "https://bitstamp.net/api/ticker/", mockServer.URL, 2)
	j := cltest.CreateSpecViaWeb(t, app, spec)
	jr := cltest.CreateJobRunViaWeb(t, app, j)
	_ = cltest.WaitForJobRunStatus(t, app.Store, jr, models.RunStatusPendingOutgoingConfirmations)
	triggerAllKeys(t, app)
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
	gasLimit := ethconfig.Defaults.Miner.GasCeil * 2
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
	t.Parallel()

	// Simulate a consumer contract calling to obtain ETH quotes in 3 different currencies
	// in a single callback.
	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	user, _, _, consumerContract, operatorContract, b := setupMultiWordContracts(t)
	app, cleanup := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, b)
	defer cleanup()
	app.Config.Set("ETH_HEAD_TRACKER_MAX_BUFFER_SIZE", 100)
	app.Config.Set("MIN_OUTGOING_CONFIRMATIONS", 1)

	sendingKeys, err := app.Store.KeyStore.SendingKeys()
	require.NoError(t, err)
	authorizedSenders := []common.Address{sendingKeys[0].Address.Address()}
	_, err = operatorContract.SetAuthorizedSenders(user, authorizedSenders)
	require.NoError(t, err)
	b.Commit()

	// Fund node account with ETH.
	n, err := b.NonceAt(context.Background(), user.From, nil)
	require.NoError(t, err)
	tx := types.NewTransaction(n, sendingKeys[0].Address.Address(), big.NewInt(1000000000000000000), 21000, big.NewInt(1), nil)
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
	spec := string(cltest.MustReadFile(t, "../testdata/jsonspecs/multiword_v1_runlog.json"))
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
			triggerAllKeys(t, app)
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
	gasLimit := ethconfig.Defaults.Miner.GasCeil * 2
	b := backends.NewSimulatedBackend(genesisData, gasLimit)
	linkTokenAddress, _, linkContract, err := link_token_interface.DeployLinkToken(owner, b)
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
		min,              // -2**191
		max,              // 2**191 - 1
		accessAddress,
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
	config, _, ormCleanup := cltest.BootstrapThrowawayORM(t, fmt.Sprintf("%s%d", dbName, port), true)
	config.Dialect = dialects.PostgresWithoutLock
	app, appCleanup := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, b)
	_, _, err := app.OCRKeyStore.GenerateEncryptedP2PKey()
	require.NoError(t, err)
	p2pIDs := app.OCRKeyStore.DecryptedP2PKeys()
	require.NoError(t, err)
	require.Len(t, p2pIDs, 1)
	peerID := p2pIDs[0].MustGetPeerID().Raw()

	app.Config.Set("P2P_PEER_ID", peerID)
	app.Config.Set("P2P_LISTEN_PORT", port)
	app.Config.Set("ETH_HEAD_TRACKER_MAX_BUFFER_SIZE", 100)
	app.Config.Set("MIN_OUTGOING_CONFIRMATIONS", 1)
	app.Config.Set("CHAINLINK_DEV", true) // Disables ocr spec validation so we can have fast polling for the test.

	sendingKeys, err := app.Store.KeyStore.SendingKeys()
	require.NoError(t, err)
	transmitter := sendingKeys[0].Address.Address()

	// Fund the transmitter address with some ETH
	n, err := b.NonceAt(context.Background(), owner.From, nil)
	require.NoError(t, err)

	tx := types.NewTransaction(n, transmitter, big.NewInt(1000000000000000000), 21000, big.NewInt(1), nil)
	signedTx, err := owner.Signer(owner.From, tx)
	require.NoError(t, err)
	err = b.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)
	b.Commit()

	_, kb, err := app.OCRKeyStore.GenerateEncryptedOCRKeyBundle()
	require.NoError(t, err)
	return app, peerID, transmitter, kb, func() {
		ormCleanup()
		appCleanup()
	}
}

func TestIntegration_OCR(t *testing.T) {
	t.Parallel()

	owner, b, ocrContractAddress, ocrContract := setupOCRContracts(t)

	// Note it's plausible these ports could be occupied on a CI machine.
	// May need a port randomize + retry approach if we observe collisions.
	appBootstrap, bootstrapPeerID, _, _, cleanup := setupNode(t, owner, 19999, "bootstrap", b)
	defer cleanup()

	var (
		oracles      []confighelper.OracleIdentityExtra
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
		// GracePeriod < ObservationTimeout
		app.Config.Set("OCR_OBSERVATION_GRACE_PERIOD", "100ms")

		kbs = append(kbs, kb)
		apps = append(apps, app)
		transmitters = append(transmitters, transmitter)

		oracles = append(oracles, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnChainSigningAddress: ocrtypes.OnChainSigningAddress(kb.OnChainSigningAddress),
				TransmitAddress:       transmitter,
				OffchainPublicKey:     ocrtypes.OffchainPublicKey(kb.OffChainPublicKey),
				PeerID:                peerID,
			},
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
	var servers, slowServers = make([]*httptest.Server, 4), make([]*httptest.Server, 4)
	// We expect metadata of:
	//  latestAnswer:nil // First call
	//  latestAnswer:0
	//  latestAnswer:10
	//  latestAnswer:20
	//  latestAnswer:30
	expectedMeta := map[string]struct{}{
		"0": {}, "10": {}, "20": {}, "30": {},
	}
	for i := 0; i < 4; i++ {
		err = apps[i].StartAndConnect()
		require.NoError(t, err)
		defer apps[i].Stop()

		// Since this API speed is > ObservationTimeout we should ignore it and still produce values.
		slowServers[i] = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			time.Sleep(5 * time.Second)
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(`{"data":10}`))
		}))
		defer slowServers[i].Close()
		servers[i] = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			b, err := ioutil.ReadAll(req.Body)
			require.NoError(t, err)
			var m models.BridgeMetaDataJSON
			require.NoError(t, json.Unmarshal(b, &m))
			if m.Meta.LatestAnswer != nil && m.Meta.UpdatedAt != nil {
				delete(expectedMeta, m.Meta.LatestAnswer.String())
			}
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(`{"data":10}`))
		}))
		defer servers[i].Close()
		u, _ := url.Parse(servers[i].URL)
		apps[i].Store.CreateBridgeType(&models.BridgeType{
			Name: models.TaskType(fmt.Sprintf("bridge%d", i)),
			URL:  models.WebURL(*u),
		})

		// Note we need: observationTimeout + observationGracePeriod + DeltaGrace (500ms) < DeltaRound (1s)
		// So 200ms + 200ms + 500ms < 1s
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
observationTimeout = "100ms"
contractConfigConfirmations = 1
contractConfigTrackerPollInterval = "1s"
observationSource = """
    // data source 1
    ds1          [type=bridge name="%s"];
    ds1_parse    [type=jsonparse path="data"];
    ds1_multiply [type=multiply times=%d];

    // data source 2
    ds2          [type=http method=GET url="%s"];
    ds2_parse    [type=jsonparse path="data"];
    ds2_multiply [type=multiply times=%d];

    ds1 -> ds1_parse -> ds1_multiply -> answer1;
    ds2 -> ds2_parse -> ds2_multiply -> answer1;

	answer1 [type=median index=0];
"""
`, ocrContractAddress, bootstrapPeerID, kbs[i].ID, transmitters[i], fmt.Sprintf("bridge%d", i), i, slowServers[i].URL, i))
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
			// Want at least 2 runs so we see all the metadata.
			pr := cltest.WaitForPipelineComplete(t, ic, jids[ic], 2, 0, apps[ic].GetJobORM(), 1*time.Minute, 1*time.Second)
			jb, err := pr[0].Outputs.MarshalJSON()
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
	}, 10*time.Second, 200*time.Millisecond).Should(gomega.Equal("20"))

	for _, app := range apps {
		jobs, err := app.JobORM.JobsV2()
		require.NoError(t, err)
		// No spec errors
		for _, j := range jobs {
			ignore := 0
			for i := range j.JobSpecErrors {
				// Non-fatal timing related error, ignore for testing.
				if strings.Contains(j.JobSpecErrors[i].Description, "leader's phase conflicts tGrace timeout") {
					ignore++
				}
			}
			require.Len(t, j.JobSpecErrors, ignore)
		}
	}
	assert.Len(t, expectedMeta, 0, "expected metadata %v", expectedMeta)
}

func TestIntegration_DirectRequest(t *testing.T) {
	config, cfgCleanup := cltest.NewConfig(t)
	defer cfgCleanup()

	httpAwaiter := cltest.NewAwaiter()
	httpServer, assertCalled := cltest.NewHTTPMockServer(
		t,
		http.StatusOK,
		"GET",
		`{"USD": "31982"}`,
		func(header http.Header, _ string) {
			httpAwaiter.ItHappened()
		},
	)
	defer assertCalled()

	ethClient, sub, assertMockCalls := cltest.NewEthMocks(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()

	sub.On("Err").Return(nil).Maybe()
	sub.On("Unsubscribe").Return(nil).Maybe()

	ethClient.On("HeaderByNumber", mock.Anything, mock.AnythingOfType("*big.Int")).Return(cltest.Head(10), nil)

	var headCh chan<- *models.Head
	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).Maybe().
		Run(func(args mock.Arguments) {
			headCh = args.Get(1).(chan<- *models.Head)
		}).
		Return(sub, nil)

	ethClient.On("Dial", mock.Anything).Return(nil)
	ethClient.On("ChainID", mock.Anything).Maybe().Return(app.Store.Config.ChainID(), nil)
	ethClient.On("FilterLogs", mock.Anything, mock.Anything).Maybe().Return([]types.Log{}, nil)
	ethClient.On("HeaderByNumber", mock.Anything, mock.AnythingOfType("*big.Int")).Return(cltest.Head(0), nil)
	logsCh := cltest.MockSubscribeToLogsCh(ethClient, sub)

	require.NoError(t, app.StartAndConnect())

	store := app.Store
	eventBroadcaster := postgres.NewEventBroadcaster(config.DatabaseURL(), 0, 0)
	eventBroadcaster.Start()
	defer eventBroadcaster.Close()

	pipelineORM := pipeline.NewORM(store.ORM.DB, config, eventBroadcaster)
	jobORM := job.NewORM(store.ORM.DB, store.Config, pipelineORM, eventBroadcaster, &postgres.NullAdvisoryLocker{})

	directRequestSpec := string(cltest.MustReadFile(t, "../testdata/tomlspecs/direct-request-spec.toml"))
	directRequestSpec = strings.Replace(directRequestSpec, "http://example.com", httpServer.URL, 1)
	request := web.CreateJobRequest{TOML: directRequestSpec}
	output, err := json.Marshal(request)
	require.NoError(t, err)
	job := cltest.CreateJobViaWeb(t, app, output)

	eventBroadcaster.Notify(postgres.ChannelJobCreated, "")

	runLog := cltest.NewRunLog(t, job.DirectRequestSpec.OnChainJobSpecID.Hash(), job.DirectRequestSpec.ContractAddress.Address(), cltest.NewAddress(), 1, `{}`)
	var logs chan<- types.Log
	cltest.CallbackOrTimeout(t, "obtain log channel", func() {
		logs = <-logsCh
	}, 5*time.Second)
	cltest.CallbackOrTimeout(t, "send run log", func() {
		logs <- runLog
	}, 30*time.Second)

	eventBroadcaster.Notify(postgres.ChannelRunStarted, "")
	headCh <- &models.Head{Number: 10}
	headCh <- &models.Head{Number: 11}

	httpAwaiter.AwaitOrFail(t)

	runs := cltest.WaitForPipelineComplete(t, 0, job.ID, 1, 3, jobORM, 5*time.Second, 300*time.Millisecond)
	require.Len(t, runs, 1)
	run := runs[0]
	require.Len(t, run.PipelineTaskRuns, 3)
	require.Empty(t, run.PipelineTaskRuns[0].Error)
	require.Empty(t, run.PipelineTaskRuns[1].Error)
	require.Empty(t, run.PipelineTaskRuns[2].Error)
}

func TestIntegration_GasUpdater(t *testing.T) {
	t.Parallel()

	c, cfgCleanup := cltest.NewConfig(t)
	defer cfgCleanup()
	c.Set("ETH_GAS_PRICE_DEFAULT", 5000000000)
	c.Set("GAS_UPDATER_ENABLED", true)
	c.Set("GAS_UPDATER_BLOCK_DELAY", 0)
	c.Set("GAS_UPDATER_BLOCK_HISTORY_SIZE", 2)
	// Limit the headtracker backfill depth just so we aren't here all week
	c.Set("ETH_FINALITY_DEPTH", 3)

	ethClient, sub, assertMocksCalled := cltest.NewEthMocks(t)
	defer assertMocksCalled()
	chchNewHeads := make(chan chan<- *models.Head, 1)

	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, c,
		ethClient,
	)
	defer cleanup()

	b41 := gasupdater.Block{
		Number:       41,
		Hash:         cltest.NewHash(),
		Transactions: cltest.TransactionsFromGasPrices(41000000000, 41500000000),
	}
	b42 := gasupdater.Block{
		Number:       42,
		Hash:         cltest.NewHash(),
		Transactions: cltest.TransactionsFromGasPrices(44000000000, 45000000000),
	}
	b43 := gasupdater.Block{
		Number:       43,
		Hash:         cltest.NewHash(),
		Transactions: cltest.TransactionsFromGasPrices(48000000000, 49000000000, 31000000000),
	}

	h40 := models.Head{Hash: cltest.NewHash(), Number: 40}
	h41 := models.Head{Hash: b41.Hash, ParentHash: h40.Hash, Number: 41}
	h42 := models.Head{Hash: b42.Hash, ParentHash: h41.Hash, Number: 42}

	sub.On("Err").Return(nil)
	sub.On("Unsubscribe").Return(nil).Maybe()

	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) { chchNewHeads <- args.Get(1).(chan<- *models.Head) }).
		Return(sub, nil)
	// Nonce syncer
	ethClient.On("PendingNonceAt", mock.Anything, mock.Anything).Maybe().Return(uint64(0), nil)

	// GasUpdater boot calls
	ethClient.On("HeaderByNumber", mock.Anything, mock.AnythingOfType("*big.Int")).Return(&h42, nil)
	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 2 &&
			b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == "0x29" &&
			b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == "0x2a"
	})).Return(nil).Run(func(args mock.Arguments) {
		elems := args.Get(1).([]rpc.BatchElem)
		elems[0].Result = &b41
		elems[1].Result = &b42
	})

	ethClient.On("Dial", mock.Anything).Return(nil)
	ethClient.On("ChainID", mock.Anything).Return(c.ChainID(), nil)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(oneETH.ToInt(), nil)

	require.NoError(t, app.Start())
	var newHeads chan<- *models.Head
	select {
	case newHeads = <-chchNewHeads:
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for app to subscribe")
	}

	assert.Equal(t, "41500000000", app.Config.EthGasPriceDefault().String())

	// GasUpdater new blocks
	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 2 &&
			b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == "0x2a" &&
			b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == "0x2b"
	})).Return(nil).Run(func(args mock.Arguments) {
		elems := args.Get(1).([]rpc.BatchElem)
		elems[0].Result = &b43
		elems[1].Result = &b42
	})

	// HeadTracker backfill
	ethClient.On("HeaderByNumber", mock.Anything, big.NewInt(42)).Return(&h42, nil)
	ethClient.On("HeaderByNumber", mock.Anything, big.NewInt(41)).Return(&h41, nil)

	// Simulate one new head and check the gas price got updated
	newHeads <- cltest.Head(43)

	gomega.NewGomegaWithT(t).Eventually(func() string {
		return c.EthGasPriceDefault().String()
	}, cltest.DBWaitTimeout, cltest.DBPollingInterval).Should(gomega.Equal("45000000000"))
}

func triggerAllKeys(t *testing.T, app *cltest.TestApplication) {
	keys, err := app.Store.KeyStore.SendingKeys()
	require.NoError(t, err)
	for _, k := range keys {
		app.BPTXM.Trigger(k.Address.Address())
	}
}
