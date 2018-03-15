package services_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
)

func TestEthereumListener_Connect_WithJobs(t *testing.T) {
	t.Parallel()

	el, cleanup := cltest.NewEthereumListener()
	defer cleanup()
	eth := cltest.MockEthOnStore(el.Store)

	j1 := cltest.NewJobWithLogInitiator()
	j2 := cltest.NewJobWithLogInitiator()
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

func TestEthereumListener_reconnectLoop_Resubscribing(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)
	j1 := cltest.NewJobWithLogInitiator()
	j2 := cltest.NewJobWithLogInitiator()
	assert.Nil(t, store.SaveJob(&j1))
	assert.Nil(t, store.SaveJob(&j2))

	eth.RegisterSubscription("logs")
	eth.RegisterSubscription("logs")

	el := services.EthereumListener{Store: store}
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

func TestEthereumListener_AttachedToHeadTracker(t *testing.T) {
	t.Parallel()

	el, cleanup := cltest.NewEthereumListener()
	store := el.Store
	defer cleanup()
	eth := cltest.MockEthOnStore(store)
	j1 := cltest.NewJobWithLogInitiator()
	j2 := cltest.NewJobWithLogInitiator()
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

func TestEthereumListener_AddJob_Listening(t *testing.T) {
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
			el, cleanup := cltest.NewEthereumListener()
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
					common.StringToHash("requestID"),
					common.StringToHash(j.ID),
				},
			}

			cltest.WaitForRuns(t, j, store, test.wantCount)

			eth.EventuallyAllCalled(t)
		})
	}
}

func TestHeadTracker_New(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	cltest.MockEthOnStore(store)
	assert.Nil(t, store.Save(cltest.IndexableBlockNumber(1)))
	last := cltest.IndexableBlockNumber(16)
	assert.Nil(t, store.Save(last))
	assert.Nil(t, store.Save(cltest.IndexableBlockNumber(10)))

	ht := services.NewHeadTracker(store)
	assert.Nil(t, ht.Start())
	assert.Equal(t, last.Number, ht.LastRecord().Number)
}

func TestHeadTracker_Get(t *testing.T) {
	t.Parallel()

	start := cltest.IndexableBlockNumber(5)

	tests := []struct {
		name      string
		initial   *models.IndexableBlockNumber
		toSave    *models.IndexableBlockNumber
		want      *big.Int
		wantError bool
	}{
		{"greater", start, cltest.IndexableBlockNumber(6), big.NewInt(6), false},
		{"less than", start, cltest.IndexableBlockNumber(1), big.NewInt(5), false},
		{"zero", start, cltest.IndexableBlockNumber(0), big.NewInt(5), true},
		{"nil", start, nil, big.NewInt(5), true},
		{"nil no initial", nil, nil, nil, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore()
			defer cleanup()
			cltest.MockEthOnStore(store)
			if test.initial != nil {
				assert.Nil(t, store.Save(test.initial))
			}

			ht := services.NewHeadTracker(store)
			ht.Start()
			defer ht.Stop()

			err := ht.Save(test.toSave)
			if test.wantError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, test.want, ht.LastRecord().ToInt())
		})
	}
}

func TestHeadTracker_Start_NewHeads(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)
	ht := services.NewHeadTracker(store)
	defer ht.Stop()

	eth.RegisterSubscription("newHeads")

	assert.Nil(t, ht.Start())
	eth.EventuallyAllCalled(t)
}

func TestHeadTracker_HeadTrackableCallbacks(t *testing.T) {
	t.Parallel()
	g := gomega.NewGomegaWithT(t)

	store, cleanup := cltest.NewStore()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)
	ht := services.NewHeadTracker(store, cltest.NeverSleeper{})

	checker := &cltest.MockHeadTrackable{}
	ht.Attach(checker)

	headers := make(chan models.BlockHeader)
	eth.RegisterSubscription("newHeads", headers)

	assert.Nil(t, ht.Start())
	assert.Equal(t, 1, checker.ConnectedCount)
	assert.Equal(t, 0, checker.DisconnectedCount)
	assert.Equal(t, 0, checker.OnNewHeadCount)

	headers <- models.BlockHeader{Number: cltest.BigHexInt(1)}
	g.Eventually(func() int { return checker.OnNewHeadCount }).Should(gomega.Equal(1))
	assert.Equal(t, 1, checker.ConnectedCount)
	assert.Equal(t, 0, checker.DisconnectedCount)

	ht.Stop()
	assert.Equal(t, 1, checker.DisconnectedCount)
	assert.Equal(t, 1, checker.ConnectedCount)
	assert.Equal(t, 1, checker.OnNewHeadCount)
}

func TestHeadTracker_ReconnectOnError(t *testing.T) {
	t.Parallel()
	g := gomega.NewGomegaWithT(t)

	store, cleanup := cltest.NewStore()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)
	ht := services.NewHeadTracker(store, cltest.NeverSleeper{})

	firstSub := eth.RegisterSubscription("newHeads")
	headers := make(chan models.BlockHeader)
	eth.RegisterSubscription("newHeads", headers)

	checker := &cltest.MockHeadTrackable{}
	ht.Attach(checker)

	// connect
	assert.Nil(t, ht.Start())
	assert.Equal(t, 1, checker.ConnectedCount)
	assert.Equal(t, 0, checker.DisconnectedCount)
	assert.Equal(t, 0, checker.OnNewHeadCount)

	// disconnect
	firstSub.Errors <- errors.New("Test error to force reconnect")
	g.Eventually(func() int { return checker.ConnectedCount }).Should(gomega.Equal(2))
	assert.Equal(t, 1, checker.DisconnectedCount)
	assert.Equal(t, 0, checker.OnNewHeadCount)

	// new head
	headers <- models.BlockHeader{Number: cltest.BigHexInt(1)}
	g.Eventually(func() int { return checker.OnNewHeadCount }).Should(gomega.Equal(1))
	assert.Equal(t, 2, checker.ConnectedCount)
	assert.Equal(t, 1, checker.DisconnectedCount)
}
