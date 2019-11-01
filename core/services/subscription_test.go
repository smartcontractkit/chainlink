package services_test

import (
	"math/big"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
	"chainlink/core/internal/cltest"
	"chainlink/core/internal/mocks"
	"chainlink/core/services"
	strpkg "chainlink/core/store"
	"chainlink/core/store/models"
	"chainlink/core/store/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServices_NewInitiatorSubscription_BackfillLogs(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	eth := cltest.MockEthOnStore(t, store)

	job := cltest.NewJobWithLogInitiator()
	initr := job.Initiators[0]
	log := cltest.LogFromFixture(t, "testdata/subscription_logs.json")
	eth.Register("eth_getLogs", []models.Log{log})
	eth.RegisterSubscription("logs")

	var count int32
	callback := func(*strpkg.Store, models.LogRequest) { atomic.AddInt32(&count, 1) }
	fromBlock := cltest.Head(0)
	sub, err := services.NewInitiatorSubscription(initr, job, store, fromBlock, callback)
	assert.NoError(t, err)
	defer sub.Unsubscribe()

	eth.EventuallyAllCalled(t)

	gomega.NewGomegaWithT(t).Eventually(func() int32 {
		return atomic.LoadInt32(&count)
	}).Should(gomega.Equal(int32(1)))
}

func TestServices_NewInitiatorSubscription_BackfillLogs_WithNoHead(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	eth := cltest.MockEthOnStore(t, store)

	job := cltest.NewJobWithLogInitiator()
	initr := job.Initiators[0]
	eth.RegisterSubscription("logs")

	var count int32
	callback := func(*strpkg.Store, models.LogRequest) { atomic.AddInt32(&count, 1) }
	sub, err := services.NewInitiatorSubscription(initr, job, store, nil, callback)
	assert.NoError(t, err)
	defer sub.Unsubscribe()

	eth.EventuallyAllCalled(t)
	assert.Equal(t, int32(0), atomic.LoadInt32(&count))
}

func TestServices_NewInitiatorSubscription_PreventsDoubleDispatch(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	eth := cltest.MockEthOnStore(t, store)

	job := cltest.NewJobWithLogInitiator()
	initr := job.Initiators[0]

	log := cltest.LogFromFixture(t, "testdata/subscription_logs.json")
	eth.Register("eth_getLogs", []models.Log{log}) // backfill
	logsChan := make(chan models.Log)
	eth.RegisterSubscription("logs", logsChan)

	var count int32
	callback := func(*strpkg.Store, models.LogRequest) { atomic.AddInt32(&count, 1) }
	head := cltest.Head(0)
	sub, err := services.NewInitiatorSubscription(initr, job, store, head, callback)
	assert.NoError(t, err)
	defer sub.Unsubscribe()

	// Add the same original log
	logsChan <- log
	// Add a log after the repeated log to make sure it gets processed
	log2 := cltest.LogFromFixture(t, "testdata/requestLog0original.json")
	logsChan <- log2

	eth.EventuallyAllCalled(t)
	g := gomega.NewGomegaWithT(t)
	g.Eventually(func() int32 { return atomic.LoadInt32(&count) }).Should(gomega.Equal(int32(2)))
}

func TestServices_ReceiveLogRequest_IgnoredLogWithRemovedFlag(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	jobSpec := cltest.NewJobWithLogInitiator()
	require.NoError(t, store.CreateJob(&jobSpec))

	log := models.InitiatorLogEvent{
		JobSpec: jobSpec,
		Log: models.Log{
			Removed: true,
		},
	}

	originalCount := 0
	err := store.ORM.DB.Model(&models.JobRun{}).Count(&originalCount).Error
	require.NoError(t, err)

	services.ReceiveLogRequest(store, log)

	gomega.NewGomegaWithT(t).Consistently(func() int {
		count := 0
		err := store.ORM.DB.Model(&models.JobRun{}).Count(&count).Error
		require.NoError(t, err)
		return count - originalCount
	}).Should(gomega.Equal(0))
}

func TestServices_NewInitiatorSubscription_ReplayFromBlock(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	txmMock := mocks.NewMockTxManager(ctrl)
	store.TxManager = txmMock

	cases := []struct {
		name                string
		currentHead         int
		initrParamFromBlock *models.Big
		wantFromBlock       *big.Int
	}{
		{"head < ReplayFromBlock, no initr fromBlock", 5, nil, big.NewInt(11)},
		{"head > ReplayFromBlock, no initr fromBlock", 14, nil, big.NewInt(15)},
		{"head < ReplayFromBlock, initr fromBlock > ReplayFromBlock", 5, models.NewBig(big.NewInt(12)), big.NewInt(12)},
		{"head < ReplayFromBlock, initr fromBlock < ReplayFromBlock", 5, models.NewBig(big.NewInt(9)), big.NewInt(11)},
		{"head > ReplayFromBlock, initr fromBlock > ReplayFromBlock", 14, models.NewBig(big.NewInt(12)), big.NewInt(15)},
		{"head > ReplayFromBlock, initr fromBlock < ReplayFromBlock", 14, models.NewBig(big.NewInt(9)), big.NewInt(15)},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			currentHead := cltest.Head(test.currentHead)

			store.Config.Set(orm.EnvVarName("ReplayFromBlock"), 10)

			job := cltest.NewJobWithLogInitiator()
			initr := job.Initiators[0]
			initr.InitiatorParams.FromBlock = test.initrParamFromBlock

			expectedQuery := ethereum.FilterQuery{
				FromBlock: test.wantFromBlock,
				Addresses: []common.Address{initr.InitiatorParams.Address},
				Topics:    [][]common.Hash{},
			}

			log := cltest.LogFromFixture(t, "testdata/subscription_logs.json")

			txmMock.EXPECT().SubscribeToLogs(gomock.Any(), expectedQuery).Return(cltest.EmptyMockSubscription(), nil)
			txmMock.EXPECT().GetLogs(expectedQuery).Return([]models.Log{log}, nil)

			var wg sync.WaitGroup
			wg.Add(1)
			callback := func(*strpkg.Store, models.LogRequest) { wg.Done() }

			_, err := services.NewInitiatorSubscription(initr, job, store, currentHead, callback)
			require.NoError(t, err)

			wg.Wait()
		})
	}
}
