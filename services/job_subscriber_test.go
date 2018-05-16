package services_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
)

func TestJobSubscriber_Connect_WithJobs(t *testing.T) {
	t.Parallel()

	el, cleanup := cltest.NewJobSubscriber()
	defer cleanup()
	eth := cltest.MockEthOnStore(el.Store)

	j1, _ := cltest.NewJobWithLogInitiator()
	j2, _ := cltest.NewJobWithLogInitiator()
	assert.Nil(t, el.Store.SaveJob(&j1))
	assert.Nil(t, el.Store.SaveJob(&j2))
	eth.RegisterSubscription("logs")
	eth.RegisterSubscription("logs")

	assert.Nil(t, el.Connect(cltest.IndexableBlockNumber(1)))
	eth.EventuallyAllCalled(t)
}

func newAddr() common.Address {
	return cltest.NewAddress()
}

func TestJobSubscriber_reconnectLoop_Resubscribing(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)
	j1, _ := cltest.NewJobWithLogInitiator()
	j2, _ := cltest.NewJobWithLogInitiator()
	assert.Nil(t, store.SaveJob(&j1))
	assert.Nil(t, store.SaveJob(&j2))

	eth.RegisterSubscription("logs")
	eth.RegisterSubscription("logs")

	el := services.JobSubscriber{Store: store}
	assert.Nil(t, el.Connect(cltest.IndexableBlockNumber(1)))
	assert.Equal(t, 2, len(el.Jobs()))
	el.Disconnect()
	assert.Equal(t, 0, len(el.Jobs()))

	eth.RegisterSubscription("logs")
	eth.RegisterSubscription("logs")
	assert.Nil(t, el.Connect(cltest.IndexableBlockNumber(2)))
	assert.Equal(t, 2, len(el.Jobs()))
	el.Disconnect()
	assert.Equal(t, 0, len(el.Jobs()))
	eth.EventuallyAllCalled(t)
}

func TestJobSubscriber_AttachedToHeadTracker(t *testing.T) {
	t.Parallel()

	el, cleanup := cltest.NewJobSubscriber()
	store := el.Store
	defer cleanup()
	eth := cltest.MockEthOnStore(store)
	j1, _ := cltest.NewJobWithLogInitiator()
	j2, _ := cltest.NewJobWithLogInitiator()
	assert.Nil(t, store.SaveJob(&j1))
	assert.Nil(t, store.SaveJob(&j2))

	eth.RegisterSubscription("logs")
	eth.RegisterSubscription("logs")

	ht := services.NewHeadTracker(store)
	assert.Nil(t, ht.Start())
	id := ht.Attach(el)
	assert.Equal(t, 2, len(el.Jobs()))
	eth.EventuallyAllCalled(t)

	ht.Detach(id)
	assert.Equal(t, 0, len(el.Jobs()))
}

func TestJobSubscriber_AddJob_Listening(t *testing.T) {
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
			el, cleanup := cltest.NewJobSubscriber()
			defer cleanup()
			store := el.Store

			eth := cltest.MockEthOnStore(store)
			logChan := make(chan types.Log, 1)
			eth.RegisterSubscription("logs", logChan)

			j := cltest.NewJob()
			initr := models.Initiator{Type: test.initType}
			if !utils.IsEmptyAddress(test.initrAddr) {
				initr.Address = test.initrAddr
			}
			j.Initiators = []models.Initiator{initr}
			el.AddJob(j, cltest.IndexableBlockNumber(1))

			ht := services.NewHeadTracker(store)
			ht.Attach(el)
			assert.Nil(t, ht.Start())

			logChan <- types.Log{
				Address: test.logAddr,
				Data:    test.data,
				Topics: []common.Hash{
					services.RunLogTopic,
					cltest.StringToHash("internalID"),
					cltest.StringToHash(j.ID),
					common.BigToHash(big.NewInt(0)),
				},
			}

			cltest.WaitForRuns(t, j, store, test.wantCount)

			eth.EventuallyAllCalled(t)
		})
	}
}

func TestJobSubscriber_OnNewHead_OnlyRunPendingConfirmations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status     models.RunStatus
		wantStatus models.RunStatus
	}{
		{models.RunStatusPendingBridge, models.RunStatusPendingBridge},
		{models.RunStatusPendingConfirmations, models.RunStatusCompleted},
	}

	for _, test := range tests {
		t.Run(string(test.status), func(t *testing.T) {
			el, cleanup := cltest.NewJobSubscriber()
			defer cleanup()
			store := el.Store

			job, initr := cltest.NewJobWithWebInitiator()
			run := job.NewRun(initr)
			run.Status = test.status

			assert.Nil(t, store.SaveJob(&job))
			assert.Nil(t, store.Save(&run))
			el.OnNewHead(cltest.NewBlockHeader(10))

			refreshed, err := store.FindJobRun(run.ID)
			assert.Nil(t, err)
			assert.Equal(t, test.wantStatus, refreshed.Status)
		})
	}
}
