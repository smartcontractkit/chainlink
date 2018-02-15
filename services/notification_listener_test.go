package services_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	strpkg "github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestNotificationListener_Start_NewHeads(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)
	nl := services.NotificationListener{Store: store}
	defer nl.Stop()

	eth.RegisterSubscription("newHeads", make(chan types.Header))

	assert.Nil(t, nl.Start())
	eth.EnsureAllCalled(t)
}

func TestNotificationListener_Start_WithJobs(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)
	nl := services.NotificationListener{Store: store}
	defer nl.Stop()

	assert.Nil(t, store.SaveJob(cltest.NewJobWithLogInitiator()))
	assert.Nil(t, store.SaveJob(cltest.NewJobWithLogInitiator()))
	eth.RegisterSubscription("logs", make(chan []types.Log))
	eth.RegisterSubscription("logs", make(chan []types.Log))

	err := nl.Start()
	assert.Nil(t, err)

	eth.EnsureAllCalled(t)
}

func TestNotificationListener_AddJob(t *testing.T) {
	t.Parallel()

	initrAddress := cltest.NewAddress()

	tests := []struct {
		name       string
		initType   string
		logAddress common.Address
		wantCount  int
		data       hexutil.Bytes
	}{
		{"basic eth log", "ethlog", initrAddress, 1, hexutil.Bytes{}},
		{"non-matching eth log", "ethlog", cltest.NewAddress(), 0, hexutil.Bytes{}},
		{"basic cllog", "runlog", initrAddress, 1, cltest.StringToRunLogPayload(`{"value":"100"}`)},
		{"cllog non-matching", "runlog", cltest.NewAddress(), 0, hexutil.Bytes{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore()
			defer cleanup()
			cltest.MockEthOnStore(store)

			nl := services.NotificationListener{Store: store}
			defer nl.Stop()
			assert.Nil(t, nl.Start())

			eth := cltest.MockEthOnStore(store)
			logChan := make(chan []types.Log, 1)
			eth.RegisterSubscription("logs", logChan)

			j := cltest.NewJob()
			j.Initiators = []models.Initiator{{
				Type:    test.initType,
				Address: initrAddress,
			}}
			assert.Nil(t, store.SaveJob(j))

			nl.AddJob(*j)

			logChan <- []types.Log{{
				Address: test.logAddress,
				Data:    test.data,
				Topics:  []common.Hash{common.HexToHash("0x00"), common.HexToHash("0x01"), common.HexToHash("0x22")},
			}}
			<-time.After(100 * time.Millisecond)

			cltest.WaitForRuns(t, j, store, test.wantCount)

			eth.EnsureAllCalled(t)
		})
	}
}

func jsonFromFixture(path string) models.JSON {
	res := gjson.Get(string(cltest.LoadJSON(path)), "params.result.0")
	out := cltest.JSONFromString(res.String())
	return out
}

func TestStore_FormatLogJSON(t *testing.T) {
	t.Parallel()

	var clData models.JSON
	clDataFixture := `{"url":"https://etherprice.com/api","path":["recent","usd"],"address":"0x3cCad4715152693fE3BC4460591e3D3Fbd071b42","dataPrefix":"0x0000000000000000000000000000000000000000000000000000000000000001","functionSelector":"76005c26"}`
	assert.Nil(t, json.Unmarshal([]byte(clDataFixture), &clData))

	hwLog := cltest.LogFromFixture("../internal/fixtures/eth/subscription_logs_hello_world.json")
	exampleLog := cltest.LogFromFixture("../internal/fixtures/eth/subscription_logs.json")
	tests := []struct {
		name        string
		el          types.Log
		initr       models.Initiator
		wantErrored bool
		wantOutput  models.JSON
	}{
		{"example ethLog", exampleLog, models.Initiator{Type: "ethlog"}, false,
			jsonFromFixture("../internal/fixtures/eth/subscription_logs.json")},
		{"hello world ethLog", hwLog, models.Initiator{Type: "ethlog"}, false,
			jsonFromFixture("../internal/fixtures/eth/subscription_logs_hello_world.json")},
		{"hello world runLog", hwLog, models.Initiator{Type: "runlog"}, false,
			clData},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := services.FormatLogJSON(test.initr, test.el)
			assert.JSONEq(t, strings.ToLower(test.wantOutput.String()), strings.ToLower(output.String()))
			assert.Equal(t, test.wantErrored, (err != nil))
		})
	}
}

func TestNotificationListener_newHeadsNotification(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	store := app.Store

	ethMock := app.MockEthClient()
	nhChan := make(chan types.Header)
	ethMock.RegisterSubscription("newHeads", nhChan)
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{})
	sentAt := uint64(23456)
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+1))

	app.Start()

	j := models.NewJob()
	j.Tasks = []models.Task{cltest.NewTask("ethtx", "{}")}
	assert.Nil(t, store.SaveJob(j))

	tx := cltest.CreateTxAndAttempt(store, cltest.NewAddress(), sentAt)
	txas, err := store.AttemptsFor(tx.ID)
	assert.Nil(t, err)
	txa := txas[0]

	jr := j.NewRun()
	tr := jr.TaskRuns[0]
	result := models.RunResultWithValue(txa.Hash.String())
	tr.Result = models.RunResultPending(result)
	tr.Status = models.StatusPending
	jr.TaskRuns[0] = tr
	jr.Status = models.StatusPending
	assert.Nil(t, store.Save(jr))

	nhChan <- types.Header{}

	ethMock.EnsureAllCalled(t)
}
