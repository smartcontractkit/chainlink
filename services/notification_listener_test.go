package services_test

import (
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
		{"ethlog all address", "ethlog", noAddr, newAddr(), 1, hexutil.Bytes{}},
		{"runlog w/o address", "runlog", noAddr, newAddr(), 1, cltest.StringToRunLogData(`{"value":"100"}`)},
		{"runlog matching address", "runlog", sharedAddr, sharedAddr, 1, cltest.StringToRunLogData(`{"value":"100"}`)},
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

func TestNotificationListener_newHeadsNotification(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	store := app.Store

	ethMock := app.MockEthClient()
	nhChan := ethMock.RegisterNewHeads()
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{})
	sentAt := uint64(23456)
	confirmationAt := sentAt + 1
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(confirmationAt))

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
	result := cltest.RunResultWithValue(txa.Hash.String())
	tr.Result = result.MarkPending()
	tr.Status = models.StatusPending
	jr.TaskRuns[0] = tr
	jr.Status = models.StatusPending
	assert.Nil(t, store.Save(&jr))

	blockNumber := cltest.BigHexInt(1)
	nhChan <- models.BlockHeader{blockNumber}

	ethMock.EnsureAllCalled(t)
	assert.Equal(t, blockNumber, app.Store.HeadTracker.Get().Number)
}
