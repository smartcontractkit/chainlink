package services_test

import (
	"math/big"
	"sync/atomic"
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/eth"

	"github.com/ethereum/go-ethereum/core/types"

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
	ethClient := new(mocks.Client)
	defer ethClient.AssertExpectations(t)
	store.EthClient = ethClient

	job := cltest.NewJobWithLogInitiator()
	initr := job.Initiators[0]
	log := cltest.LogFromFixture(t, "testdata/subscription_logs.json")
	ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(cltest.EmptyMockSubscription(), nil)
	b := types.NewBlockWithHeader(&types.Header{
		Number: big.NewInt(2),
	})
	ethClient.On("BlockByNumber", mock.Anything, mock.Anything).Maybe().Return(b, nil)
	ethClient.On("FilterLogs", mock.Anything, mock.Anything).Maybe().Return([]types.Log{log}, nil)

	var count int32
	callback := func(services.RunManager, models.LogRequest) { atomic.AddInt32(&count, 1) }
	fromBlock := cltest.Head(0)
	jm := new(mocks.RunManager)
	sub, err := services.NewInitiatorSubscription(initr, store.EthClient, jm, fromBlock.NextInt(), store.Config, callback)
	assert.NoError(t, err)
	defer sub.Unsubscribe()
	gomega.NewGomegaWithT(t).Eventually(func() int32 {
		return atomic.LoadInt32(&count)
	}).Should(gomega.Equal(int32(1)))
}

func TestServices_NewInitiatorSubscription_BackfillLogs_BatchWindows(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ethClient := new(mocks.Client)
	defer ethClient.AssertExpectations(t)
	store.EthClient = ethClient

	job := cltest.NewJobWithLogInitiator()
	initr := job.Initiators[0]
	log := cltest.LogFromFixture(t, "testdata/subscription_logs.json")
	ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(cltest.EmptyMockSubscription(), nil)
	b := types.NewBlockWithHeader(&types.Header{
		Number: big.NewInt(213),
	})
	ethClient.On("BlockByNumber", mock.Anything, mock.Anything).Maybe().Return(b, nil)
	ethClient.On("FilterLogs", mock.Anything, mock.Anything).Once().Return([]types.Log{log}, nil).Run(func(args mock.Arguments) {
		query := args.Get(1).(ethereum.FilterQuery)
		assert.Equal(t, big.NewInt(1), query.FromBlock)
		assert.Equal(t, big.NewInt(100), query.ToBlock)
	})
	ethClient.On("FilterLogs", mock.Anything, mock.Anything).Once().Return([]types.Log{log}, nil).Run(func(args mock.Arguments) {
		query := args.Get(1).(ethereum.FilterQuery)
		assert.Equal(t, big.NewInt(101), query.FromBlock)
		assert.Equal(t, big.NewInt(200), query.ToBlock)
	})
	ethClient.On("FilterLogs", mock.Anything, mock.Anything).Once().Return([]types.Log{log}, nil).Run(func(args mock.Arguments) {
		query := args.Get(1).(ethereum.FilterQuery)
		assert.Equal(t, big.NewInt(201), query.FromBlock)
		assert.Equal(t, big.NewInt(213), query.ToBlock)
	})

	var count int32
	callback := func(services.RunManager, models.LogRequest) { atomic.AddInt32(&count, 1) }
	fromBlock := cltest.Head(0)
	jm := new(mocks.RunManager)
	sub, err := services.NewInitiatorSubscription(initr, store.EthClient, jm, fromBlock.NextInt(), store.Config, callback)
	assert.NoError(t, err)
	defer sub.Unsubscribe()
	gomega.NewGomegaWithT(t).Eventually(func() int32 {
		return atomic.LoadInt32(&count)
	}).Should(gomega.Equal(int32(3)))
}

func TestServices_NewInitiatorSubscription_BackfillLogs_WithNoHead(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ethClient := new(mocks.Client)
	defer ethClient.AssertExpectations(t)
	store.EthClient = ethClient

	job := cltest.NewJobWithLogInitiator()
	initr := job.Initiators[0]
	b := types.NewBlockWithHeader(&types.Header{
		Number: big.NewInt(2),
	})
	ethClient.On("BlockByNumber", mock.Anything, mock.Anything).Maybe().Return(b, nil)
	ethClient.On("FilterLogs", mock.Anything, mock.Anything).Maybe().Return([]models.Log{}, nil)
	ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).Return(cltest.EmptyMockSubscription(), nil)

	var count int32
	callback := func(services.RunManager, models.LogRequest) { atomic.AddInt32(&count, 1) }
	jm := new(mocks.RunManager)
	sub, err := services.NewInitiatorSubscription(initr, store.EthClient, jm, nil, store.Config, callback)
	assert.NoError(t, err)
	defer sub.Unsubscribe()
	assert.Equal(t, int32(0), atomic.LoadInt32(&count))
}

func TestServices_NewInitiatorSubscription_PreventsDoubleDispatch(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	rpcClient, gethClient, subMock, assertMocksCalled := cltest.NewEthMocks(t)
	defer assertMocksCalled()
	store.EthClient = eth.NewClientWith(rpcClient, gethClient)
	subMock.On("Unsubscribe").Return(nil)
	subMock.On("Err").Return(nil)

	job := cltest.NewJobWithLogInitiator()
	initr := job.Initiators[0]

	log := cltest.LogFromFixture(t, "testdata/subscription_logs.json")
	b := types.NewBlockWithHeader(&types.Header{
		Number: big.NewInt(2),
	})
	gethClient.On("BlockByNumber", mock.Anything, mock.Anything).Maybe().Return(b, nil)
	gethClient.On("FilterLogs", mock.Anything, mock.Anything).Maybe().Return([]models.Log{log}, nil)
	logsCh := cltest.MockSubscribeToLogsCh(gethClient, subMock)
	var count int32
	callback := func(services.RunManager, models.LogRequest) { atomic.AddInt32(&count, 1) }
	head := cltest.Head(0)
	jm := new(mocks.RunManager)
	sub, err := services.NewInitiatorSubscription(initr, store.EthClient, jm, head.NextInt(), store.Config, callback)
	assert.NoError(t, err)
	defer sub.Unsubscribe()
	logs := <-logsCh
	logs <- log
	// Add the same original log
	logs <- log
	// Add a log after the repeated log to make sure it gets processed
	log2 := cltest.LogFromFixture(t, "testdata/requestLog0original.json")
	logs <- log2

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

			rpcClient, gethClient, subMock, assertMocksCalled := cltest.NewEthMocks(t)
			defer assertMocksCalled()
			store.EthClient = eth.NewClientWith(rpcClient, gethClient)
			subMock.On("Err").Return(nil)
			b := types.NewBlockWithHeader(&types.Header{
				Number: big.NewInt(100),
			})
			gethClient.On("BlockByNumber", mock.Anything, mock.Anything).Maybe().Return(b, nil)
			gethClient.On("FilterLogs", mock.Anything, mock.Anything).Maybe().Return([]models.Log{}, nil)
			logsCh := cltest.MockSubscribeToLogsCh(gethClient, subMock)
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
			logs := <-logsCh
			logs <- models.Log{
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

			rpcClient, gethClient, subMock, assertMocksCalled := cltest.NewEthMocks(t)
			defer assertMocksCalled()
			store.EthClient = eth.NewClientWith(rpcClient, gethClient)
			subMock.On("Err").Maybe().Return(nil)

			logsCh := cltest.MockSubscribeToLogsCh(gethClient, subMock)
			b := types.NewBlockWithHeader(&types.Header{
				Number: big.NewInt(100),
			})
			gethClient.On("BlockByNumber", mock.Anything, mock.Anything).Maybe().Return(b, nil)
			gethClient.On("FilterLogs", mock.Anything, mock.Anything).Maybe().Return([]models.Log{}, nil)
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
			logs := <-logsCh
			logs <- models.Log{
				Address: sharedAddr,
				Data:    models.UntrustedBytes(test.data),
				Topics: []common.Hash{
					common.Hash{},
					models.IDToTopic(job.ID),
					cltest.NewAddress().Hash(),
					common.BigToHash(big.NewInt(0)),
				},
			}

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
		{"head > ReplayFromBlock, no initr fromBlock", 14, nil, big.NewInt(15)},
		{"head < ReplayFromBlock, initr fromBlock > ReplayFromBlock", 5, utils.NewBig(big.NewInt(12)), big.NewInt(12)},
		{"head < ReplayFromBlock, initr fromBlock < ReplayFromBlock", 5, utils.NewBig(big.NewInt(8)), big.NewInt(10)},
		{"head > ReplayFromBlock, initr fromBlock > ReplayFromBlock", 14, utils.NewBig(big.NewInt(12)), big.NewInt(15)},
		{"head > ReplayFromBlock, initr fromBlock < ReplayFromBlock", 14, utils.NewBig(big.NewInt(8)), big.NewInt(15)},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			ethClient := new(mocks.Client)
			defer ethClient.AssertExpectations(t)
			store.EthClient = ethClient

			currentHead := cltest.Head(test.currentHead)

			store.Config.Set(orm.EnvVarName("ReplayFromBlock"), 10)

			job := cltest.NewJobWithLogInitiator()
			job.Initiators[0].InitiatorParams.FromBlock = test.initrParamFromBlock

			b := types.NewBlockWithHeader(&types.Header{
				Number: big.NewInt(100),
			})
			expectedQuery := ethereum.FilterQuery{
				FromBlock: test.wantFromBlock,
				Addresses: []common.Address{job.Initiators[0].InitiatorParams.Address},
				Topics:    [][]common.Hash{},
			}

			log := cltest.LogFromFixture(t, "testdata/subscription_logs.json")

			ethClient.On("BlockByNumber", mock.Anything, mock.Anything).Maybe().Return(b, nil)
			ethClient.On("SubscribeFilterLogs", mock.Anything, expectedQuery, mock.Anything).Return(cltest.EmptyMockSubscription(), nil)
			expectedQuery.ToBlock = b.Number()
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
		{"head > ReplayFromBlock", 14, big.NewInt(15)},
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

			b := types.NewBlockWithHeader(&types.Header{
				Number: big.NewInt(100),
			})
			ethClient.On("BlockByNumber", mock.Anything, mock.Anything).Maybe().Return(b, nil)
			ethClient.On("SubscribeFilterLogs", mock.Anything, expectedQuery, mock.Anything).Return(cltest.EmptyMockSubscription(), nil)
			expectedQuery.ToBlock = b.Number()
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
