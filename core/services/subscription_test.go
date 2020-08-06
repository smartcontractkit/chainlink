package services_test

import (
	"math/big"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	callback := func(services.RunManager, models.LogRequest) { atomic.AddInt32(&count, 1) }
	fromBlock := cltest.Head(0)
	jm := new(mocks.RunManager)
	sub, err := services.NewInitiatorSubscription(initr, store.TxManager, jm, fromBlock.NextInt(), callback)
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
	callback := func(services.RunManager, models.LogRequest) { atomic.AddInt32(&count, 1) }
	jm := new(mocks.RunManager)
	sub, err := services.NewInitiatorSubscription(initr, store.TxManager, jm, nil, callback)
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
	callback := func(services.RunManager, models.LogRequest) { atomic.AddInt32(&count, 1) }
	head := cltest.Head(0)
	jm := new(mocks.RunManager)
	sub, err := services.NewInitiatorSubscription(initr, store.TxManager, jm, head.NextInt(), callback)
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
		Initiator: jobSpec.Initiators[0],
		Log: models.Log{
			Removed: true,
		},
	}

	_, err := store.ORM.CountOf(&models.JobRun{})
	require.NoError(t, err)

	jm := new(mocks.RunManager)
	services.ReceiveLogRequest(jm, log)
	jm.AssertExpectations(t)
}

func TestServices_StartJobSubscription(t *testing.T) {
	t.Parallel()

	sharedAddr := cltest.NewAddress()
	noAddr := common.Address{}

	tests := []struct {
		name      string
		initType  string
		initrAddr common.Address
		logAddr   common.Address
		topic0    common.Hash
		data      hexutil.Bytes
	}{
		{
			"ethlog matching address",
			"ethlog",
			sharedAddr,
			sharedAddr,
			common.Hash{},
			hexutil.Bytes{},
		},
		{
			"ethlog all address",
			"ethlog",
			noAddr,
			cltest.NewAddress(),
			common.Hash{},
			hexutil.Bytes{},
		},
		{
			"runlog v20190207 w/o address",
			"runlog",
			noAddr,
			cltest.NewAddress(),
			models.RunLogTopic20190207withoutIndexes,
			cltest.StringToVersionedLogData20190207withoutIndexes(t, "id", cltest.NewAddress(), `{"value":"100"}`),
		},
		{
			"runlog v20190207 matching address",
			"runlog",
			sharedAddr,
			sharedAddr,
			models.RunLogTopic20190207withoutIndexes,
			cltest.StringToVersionedLogData20190207withoutIndexes(t, "id", cltest.NewAddress(), `{"value":"100"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			eth := cltest.MockEthOnStore(t, store)
			eth.Register("eth_getLogs", []models.Log{})
			logChan := make(chan models.Log, 1)
			eth.RegisterSubscription("logs", logChan)

			job := cltest.NewJob()
			initr := models.Initiator{Type: test.initType}
			initr.Address = test.initrAddr
			job.Initiators = []models.Initiator{initr}
			require.NoError(t, store.CreateJob(&job))

			executeJobChannel := make(chan struct{})

			runManager := new(mocks.RunManager)
			runManager.On("Create", job.ID, mock.Anything, big.NewInt(0), mock.Anything).
				Return(nil, nil).
				Run(func(mock.Arguments) {
					executeJobChannel <- struct{}{}
				})

			subscription, err := services.StartJobSubscription(job, cltest.Head(91), store, runManager)
			require.NoError(t, err)
			assert.NotNil(t, subscription)

			logChan <- models.Log{
				Address: test.logAddr,
				Data:    models.UntrustedBytes(test.data),
				Topics: []common.Hash{
					test.topic0,
					models.IDToTopic(job.ID),
					cltest.NewAddress().Hash(),
					common.BigToHash(big.NewInt(0)),
				},
			}

			cltest.CallbackOrTimeout(t, "Create", func() {
				<-executeJobChannel
			})

			runManager.AssertExpectations(t)
			eth.EventuallyAllCalled(t)
		})
	}
}

func TestServices_StartJobSubscription_RunlogNoTopicMatch(t *testing.T) {
	t.Parallel()

	sharedAddr := cltest.NewAddress()

	tests := []struct {
		name string
		data hexutil.Bytes
	}{
		{
			"runlog w non-matching topic",
			cltest.StringToVersionedLogData20190207withoutIndexes(t, "id", cltest.NewAddress(), `{"value":"100"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			eth := cltest.MockEthOnStore(t, store)
			eth.Register("eth_getLogs", []models.Log{})
			logChan := make(chan models.Log, 1)
			eth.RegisterSubscription("logs", logChan)

			job := cltest.NewJob()
			initr := models.Initiator{Type: "runlog"}
			initr.Address = sharedAddr
			job.Initiators = []models.Initiator{initr}
			require.NoError(t, store.CreateJob(&job))

			runManager := new(mocks.RunManager)
			runManager.On("CreateErrored", mock.Anything, mock.Anything, mock.Anything).
				Return(nil, nil)

			subscription, err := services.StartJobSubscription(job, cltest.Head(91), store, runManager)
			require.NoError(t, err)
			assert.NotNil(t, subscription)

			logChan <- models.Log{
				Address: sharedAddr,
				Data:    models.UntrustedBytes(test.data),
				Topics: []common.Hash{
					common.Hash{},
					models.IDToTopic(job.ID),
					cltest.NewAddress().Hash(),
					common.BigToHash(big.NewInt(0)),
				},
			}

			eth.EventuallyAllCalled(t)
		})
	}
}

func TestServices_NewInitiatorSubscription_EthLog_ReplayFromBlock(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name                string
		currentHead         int
		initrParamFromBlock *utils.Big
		wantFromBlock       *big.Int
	}{
		{"head < ReplayFromBlock, no initr fromBlock", 5, nil, big.NewInt(10)},
		{"head > ReplayFromBlock, no initr fromBlock", 14, nil, big.NewInt(10)},
		{"head < ReplayFromBlock, initr fromBlock > ReplayFromBlock", 5, utils.NewBig(big.NewInt(12)), big.NewInt(12)},
		{"head < ReplayFromBlock, initr fromBlock < ReplayFromBlock", 5, utils.NewBig(big.NewInt(8)), big.NewInt(10)},
		{"head > ReplayFromBlock, initr fromBlock > ReplayFromBlock", 14, utils.NewBig(big.NewInt(12)), big.NewInt(12)},
		{"head > ReplayFromBlock, initr fromBlock < ReplayFromBlock", 14, utils.NewBig(big.NewInt(8)), big.NewInt(10)},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			ethClient := new(mocks.Client)
			store.EthClient = ethClient

			currentHead := cltest.Head(test.currentHead)

			store.Config.Set(orm.EnvVarName("ReplayFromBlock"), 10)

			job := cltest.NewJobWithLogInitiator()
			job.Initiators[0].InitiatorParams.FromBlock = test.initrParamFromBlock

			expectedQuery := ethereum.FilterQuery{
				FromBlock: test.wantFromBlock,
				Addresses: []common.Address{job.Initiators[0].InitiatorParams.Address},
				Topics:    [][]common.Hash{},
			}

			log := cltest.LogFromFixture(t, "testdata/subscription_logs.json")

			ethClient.On("SubscribeFilterLogs", mock.Anything, expectedQuery, mock.Anything).Return(cltest.EmptyMockSubscription(), nil)
			ethClient.On("FilterLogs", mock.Anything, expectedQuery).Return([]models.Log{log}, nil)

			executeJobChannel := make(chan struct{})

			runManager := new(mocks.RunManager)
			runManager.On("Create", job.ID, mock.Anything, big.NewInt(int64(log.BlockNumber)), mock.Anything).
				Return(nil, nil).
				Run(func(mock.Arguments) {
					executeJobChannel <- struct{}{}
				})

			_, err := services.StartJobSubscription(job, currentHead, store, runManager)
			require.NoError(t, err)

			<-executeJobChannel

			ethClient.AssertExpectations(t)
			runManager.AssertExpectations(t)
		})
	}
}

func TestServices_NewInitiatorSubscription_RunLog_ReplayFromBlock(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		currentHead   int
		wantFromBlock *big.Int
	}{
		{"head < ReplayFromBlock", 5, big.NewInt(10)},
		{"head > ReplayFromBlock", 14, big.NewInt(10)},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			ethClient := new(mocks.Client)
			store.EthClient = ethClient

			currentHead := cltest.Head(test.currentHead)

			store.Config.Set(orm.EnvVarName("ReplayFromBlock"), 10)

			job := cltest.NewJobWithRunLogInitiator()
			initr := job.Initiators[0]

			expectedQuery := ethereum.FilterQuery{
				FromBlock: test.wantFromBlock,
				Addresses: []common.Address{initr.InitiatorParams.Address},
				Topics: [][]common.Hash{
					models.TopicsForInitiatorsWhichRequireJobSpecIDTopic[models.InitiatorRunLog],
					{models.IDToTopic(initr.JobSpecID), models.IDToHexTopic(initr.JobSpecID)},
				},
			}

			receipt := cltest.TxReceiptFromFixture(t, "./eth/testdata/runlogReceipt.json")
			log := receipt.Logs[3]
			log.Topics[1] = models.IDToTopic(job.ID)

			ethClient.On("SubscribeFilterLogs", mock.Anything, expectedQuery, mock.Anything).Return(cltest.EmptyMockSubscription(), nil)
			ethClient.On("FilterLogs", mock.Anything, expectedQuery).Return([]models.Log{*log}, nil)

			executeJobChannel := make(chan struct{})

			runManager := new(mocks.RunManager)
			runManager.On("Create", job.ID, mock.Anything, big.NewInt(int64(log.BlockNumber)), mock.Anything).
				Return(nil, nil).
				Run(func(mock.Arguments) {
					executeJobChannel <- struct{}{}
				})

			_, err := services.StartJobSubscription(job, currentHead, store, runManager)
			require.NoError(t, err)

			<-executeJobChannel

			runManager.AssertExpectations(t)
			ethClient.AssertExpectations(t)
		})
	}
}
