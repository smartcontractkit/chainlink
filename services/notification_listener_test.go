package services_test

import (
	"encoding/json"
	"strings"
	"testing"

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

	eth.RegisterSubscription("newHeads", make(chan models.BlockHeader))

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

	j1 := cltest.NewJobWithLogInitiator()
	j2 := cltest.NewJobWithLogInitiator()
	assert.Nil(t, store.SaveJob(&j1))
	assert.Nil(t, store.SaveJob(&j2))
	eth.RegisterSubscription("logs", make(chan types.Log))
	eth.RegisterSubscription("logs", make(chan types.Log))

	err := nl.Start()
	assert.Nil(t, err)

	eth.EnsureAllCalled(t)
}

func newAddr() common.Address {
	return cltest.NewAddress()
}

func TestNotificationListener_AddJob_Listening(t *testing.T) {
	t.Parallel()

	sharedAddr := newAddr()
	noAddr := common.Address{}

	tests := []struct {
		name      string
		initType  string
		initrAddr common.Address
		logAddr   common.Address
		wantCount int
		data      hexutil.Bytes
	}{
		{"ethlog matching address", "ethlog", sharedAddr, sharedAddr, 1, hexutil.Bytes{}},
		{"ethlog non-matching address", "ethlog", newAddr(), newAddr(), 0, hexutil.Bytes{}},
		{"runlog w/o address", "runlog", noAddr, newAddr(), 1, cltest.StringToRunLogPayload(`{"value":"100"}`)},
		{"runlog matching address", "runlog", sharedAddr, sharedAddr, 1, cltest.StringToRunLogPayload(`{"value":"100"}`)},
		{"runlog non-matching", "runlog", newAddr(), newAddr(), 0, hexutil.Bytes{}},
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
			logChan := make(chan types.Log, 1)
			eth.RegisterSubscription("logs", logChan)

			j := cltest.NewJob()
			initr := models.Initiator{Type: test.initType}
			if !utils.IsEmptyAddress(test.initrAddr) {
				initr.Address = test.initrAddr
			}
			j.Initiators = []models.Initiator{initr}
			assert.Nil(t, store.SaveJob(&j))

			nl.AddJob(j)

			logChan <- types.Log{
				Address: test.logAddr,
				Data:    test.data,
				Topics: []common.Hash{
					services.RunLogTopic,
					common.StringToHash("requestID"),
					common.StringToHash(j.ID),
				},
			}

			cltest.WaitForRuns(t, j, store, test.wantCount)

			eth.EnsureAllCalled(t)
		})
	}
}

func jsonFromFixture(path string) models.JSON {
	res := gjson.Get(string(cltest.LoadJSON(path)), "params.result")
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
	nhChan := make(chan models.BlockHeader)
	ethMock.RegisterSubscription("newHeads", nhChan)
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{})
	sentAt := uint64(23456)
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+1))

	app.Start()

	j := models.NewJob()
	j.Tasks = []models.Task{cltest.NewTask("ethtx", "{}")}
	assert.Nil(t, store.SaveJob(&j))

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
	assert.Nil(t, store.Save(&jr))

	nhChan <- models.BlockHeader{}

	ethMock.EnsureAllCalled(t)
}

func TestServices_InitiatorsForLog(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	runLogSig := services.RunLogTopic
	requestID := common.StringToHash("42")

	elj := cltest.NewJob()
	el := models.Initiator{Type: "ethlog", Address: cltest.NewAddress()}
	elj.Initiators = []models.Initiator{el}
	assert.Nil(t, store.SaveJob(&elj))
	assert.Nil(t, store.One("JobID", elj.ID, &el))

	rlj := cltest.NewJob()
	rl := models.Initiator{Type: "runlog"}
	rlj.Initiators = []models.Initiator{rl}
	assert.Nil(t, store.SaveJob(&rlj))
	assert.Nil(t, store.One("JobID", rlj.ID, &rl))
	rljIDHash := common.StringToHash(rlj.ID)

	rlaj := cltest.NewJob()
	rla := models.Initiator{Type: "runlog", Address: cltest.NewAddress()}
	rlaj.Initiators = []models.Initiator{rla}
	assert.Nil(t, store.SaveJob(&rlaj))
	assert.Nil(t, store.One("JobID", rlaj.ID, &rla))
	rlajIDHash := common.StringToHash(rlaj.ID)

	tests := []struct {
		name string
		log  types.Log
		want []models.Initiator
	}{
		{"ethlog matching address", types.Log{Address: el.Address},
			[]models.Initiator{el}},
		{"ethlog non-matching address", types.Log{Address: cltest.NewAddress()},
			[]models.Initiator{}},
		{"runlog w/ required topic", types.Log{
			Topics: []common.Hash{runLogSig, requestID, rljIDHash},
		}, []models.Initiator{rl}},
		{"runlog w/o required topic", types.Log{
			Topics: []common.Hash{runLogSig, requestID, cltest.NewHash()},
		}, []models.Initiator{}},
		{"runlog w/ matching address", types.Log{
			Address: rla.Address,
			Topics:  []common.Hash{runLogSig, requestID, rlajIDHash},
		}, []models.Initiator{rla}},
		{"runlog w/o matching address", types.Log{
			Address: cltest.NewAddress(),
			Topics:  []common.Hash{runLogSig, requestID, rlajIDHash},
		}, []models.Initiator{}},
		{"runlog w/ matching address but not job ID", types.Log{
			Address: cltest.NewAddress(),
			Topics:  []common.Hash{runLogSig, requestID, cltest.NewHash()},
		}, []models.Initiator{}},
		{"runlog matching ethlog address", types.Log{
			Address: el.Address,
			Topics:  []common.Hash{runLogSig, requestID, rljIDHash},
		}, []models.Initiator{el, rl}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := services.InitiatorsForLog(store, test.log)
			assert.Nil(t, err)
			assert.Equal(t, test.want, actual)
		})
	}
}

// If updating this test, be sure to update the truffle suite's "expected event signature" test.
func TestServices_RunLogTopic_ExpectedEventSignature(t *testing.T) {
	t.Parallel()

	expected := "0x06f4bf36b4e011a5c499cef1113c2d166800ce4013f6c2509cab1a0e92b83fb2"
	assert.Equal(t, expected, services.RunLogTopic.Hex())
}
