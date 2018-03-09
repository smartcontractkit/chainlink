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
	strpkg "github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
)

func TestEthereumListener_Start_WithJobs(t *testing.T) {
	t.Parallel()

	el, cleanup := cltest.NewEthereumListener()
	defer cleanup()
	eth := cltest.MockEthOnStore(el.Store)
	assert.Nil(t, el.HeadTracker.Start())

	j1 := cltest.NewJobWithLogInitiator()
	j2 := cltest.NewJobWithLogInitiator()
	assert.Nil(t, el.Store.SaveJob(&j1))
	assert.Nil(t, el.Store.SaveJob(&j2))
	eth.RegisterSubscription("logs")
	eth.RegisterSubscription("logs")

	assert.Nil(t, el.Start())
	eth.EnsureAllCalled(t)
}

func newAddr() common.Address {
	return cltest.NewAddress()
}

func TestEthereumListener_Restart(t *testing.T) {
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

	ht := services.NewHeadTracker(store)
	ht.Start()

	el := services.EthereumListener{Store: store, HeadTracker: ht}
	assert.Nil(t, el.Start())
	assert.Equal(t, 2, len(el.Jobs()))
	assert.Nil(t, el.Stop())
	assert.Equal(t, 0, len(el.Jobs()))

	eth.RegisterSubscription("logs")
	eth.RegisterSubscription("logs")
	assert.Nil(t, el.Start())
	assert.Equal(t, 2, len(el.Jobs()))
	assert.Nil(t, el.Stop())
	assert.Equal(t, 0, len(el.Jobs()))
	eth.EnsureAllCalled(t)
}

func TestEthereumListener_Reconnected(t *testing.T) {
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

	ht := services.NewHeadTracker(store)
	el := services.EthereumListener{Store: store, HeadTracker: ht}
	el.Start()
	assert.Nil(t, ht.Start())
	assert.Equal(t, 2, len(el.Jobs()))
	assert.Nil(t, ht.Stop())
	assert.Equal(t, 0, len(el.Jobs()))

	eth.RegisterNewHeads()
	eth.RegisterSubscription("logs")
	eth.RegisterSubscription("logs")
	assert.Nil(t, ht.Start())
	assert.Equal(t, 2, len(el.Jobs()))
	assert.Nil(t, ht.Stop())
	assert.Equal(t, 0, len(el.Jobs()))
	eth.EnsureAllCalled(t)
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
			cltest.MockEthOnStore(store)
			assert.Nil(t, el.HeadTracker.Start())
			assert.Nil(t, el.Start())

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

			el.AddJob(j)

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

func TestEthereumListener_newHeadsNotification(t *testing.T) {
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
	j.Tasks = []models.Task{cltest.NewTask("ethtx")}
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
	nhChan <- models.BlockHeader{Number: blockNumber}

	ethMock.EnsureAllCalled(t)
	assert.Equal(t, blockNumber, app.EthereumListener.HeadTracker.Get().Number)
}

func TestHeadTracker_New(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	cltest.MockEthOnStore(store)
	assert.Nil(t, store.Save(models.NewIndexableBlockNumber(big.NewInt(1))))
	last := models.NewIndexableBlockNumber(big.NewInt(0x10))
	assert.Nil(t, store.Save(last))
	assert.Nil(t, store.Save(models.NewIndexableBlockNumber(big.NewInt(0xf))))

	ht := services.NewHeadTracker(store)
	assert.Nil(t, ht.Start())
	assert.Equal(t, last.Number, ht.Get().Number)
}

func TestHeadTracker_Get(t *testing.T) {
	t.Parallel()

	start := models.NewIndexableBlockNumber(big.NewInt(5))

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

			assert.Equal(t, test.want, ht.Get().ToInt())
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

	eth.RegisterSubscription("newHeads", make(chan models.BlockHeader))

	assert.Nil(t, ht.Start())
	eth.EnsureAllCalled(t)
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

	firstSub := eth.RegisterSubscription("newHeads", make(chan models.BlockHeader))
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
