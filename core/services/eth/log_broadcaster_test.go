package eth_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/eth"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	ethsvc "github.com/smartcontractkit/chainlink/core/services/eth"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLogBroadcaster_AwaitsInitialSubscribersOnStartup(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	const (
		blockHeight uint64 = 123
	)

	ethClient := new(mocks.Client)
	sub := new(mocks.Subscription)
	listener := new(mocks.LogListener)

	chOkayToAssert := make(chan struct{}) // avoid flaky tests

	listener.On("OnConnect").Return()
	listener.On("OnDisconnect").Return().Run(func(mock.Arguments) { close(chOkayToAssert) })

	sub.On("Unsubscribe").Return()
	sub.On("Err").Return(nil)

	chSubscribe := make(chan struct{}, 10)
	ethClient.On("SubscribeToLogs", mock.Anything, mock.Anything, mock.Anything).
		Return(sub, nil).
		Run(func(mock.Arguments) { chSubscribe <- struct{}{} })
	ethClient.On("GetBlockHeight").Return(blockHeight, nil)

	lb := ethsvc.NewLogBroadcaster(ethClient, store.ORM)
	lb.AddDependents(2)
	lb.Start()

	lb.Register(common.Address{}, listener)

	g.Consistently(func() int { return len(chSubscribe) }).Should(gomega.Equal(0))
	lb.DependentReady()
	g.Consistently(func() int { return len(chSubscribe) }).Should(gomega.Equal(0))
	lb.DependentReady()
	g.Eventually(func() int { return len(chSubscribe) }).Should(gomega.Equal(1))
	g.Consistently(func() int { return len(chSubscribe) }).Should(gomega.Equal(1))

	lb.Stop()

	<-chOkayToAssert

	ethClient.AssertExpectations(t)
	sub.AssertExpectations(t)
}

func TestLogBroadcaster_ResubscribesOnAddOrRemoveContract(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	const (
		numContracts        = 3
		blockHeight  uint64 = 123
	)

	ethClient := new(mocks.Client)
	sub := new(mocks.Subscription)

	var subscribeCalls int
	var unsubscribeCalls int
	ethClient.On("SubscribeToLogs", mock.Anything, mock.Anything, mock.Anything).
		Return(sub, nil).
		Run(func(args mock.Arguments) {
			subscribeCalls++
			q := args.Get(2).(ethereum.FilterQuery)
			require.Equal(t, int64(blockHeight), q.FromBlock.Int64())
		})
	ethClient.On("GetLatestBlock").
		Return(eth.Block{Number: hexutil.Uint64(blockHeight)}, nil)
	ethClient.On("GetLogs", mock.Anything).
		Return(nil, nil)
	sub.On("Unsubscribe").
		Return().
		Run(func(mock.Arguments) { unsubscribeCalls++ })
	sub.On("Err").Return(nil)

	lb := ethsvc.NewLogBroadcaster(ethClient, store.ORM, 10)
	lb.Start()

	type registration struct {
		common.Address
		ethsvc.LogListener
	}
	registrations := make([]registration, numContracts)
	for i := 0; i < numContracts; i++ {
		listener := new(mocks.LogListener)
		listener.On("OnConnect").Return()
		listener.On("OnDisconnect").Return()
		registrations[i] = registration{cltest.NewAddress(), listener}
		lb.Register(registrations[i].Address, registrations[i].LogListener)
	}

	require.Eventually(t, func() bool { return subscribeCalls == 1 }, 5*time.Second, 10*time.Millisecond)
	gomega.NewGomegaWithT(t).Consistently(subscribeCalls).Should(gomega.Equal(1))
	gomega.NewGomegaWithT(t).Consistently(unsubscribeCalls).Should(gomega.Equal(0))

	for _, r := range registrations {
		lb.Unregister(r.Address, r.LogListener)
	}
	require.Eventually(t, func() bool { return unsubscribeCalls == 1 }, 5*time.Second, 10*time.Millisecond)
	gomega.NewGomegaWithT(t).Consistently(subscribeCalls).Should(gomega.Equal(1))

	lb.Stop()
	gomega.NewGomegaWithT(t).Consistently(unsubscribeCalls).Should(gomega.Equal(1))

	ethClient.AssertExpectations(t)
	sub.AssertExpectations(t)
}

type funcLogListener struct {
	fn func(log interface{}, err error)
}

func (fn funcLogListener) HandleLog(log interface{}, err error) {
	fn.fn(log, err)
}
func (fn funcLogListener) OnConnect()    {}
func (fn funcLogListener) OnDisconnect() {}

func TestLogBroadcaster_BroadcastsToCorrectRecipients(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	const blockHeight uint64 = 0

	ethClient := new(mocks.Client)
	sub := new(mocks.Subscription)

	chchRawLogs := make(chan chan<- eth.Log, 1)
	ethClient.On("SubscribeToLogs", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			q := args.Get(2).(ethereum.FilterQuery)
			require.Equal(t, int64(blockHeight), q.FromBlock.Int64())

			chchRawLogs <- args.Get(1).(chan<- eth.Log)
		}).
		Return(sub, nil).
		Once()
	ethClient.On("GetLatestBlock").
		Return(eth.Block{Number: hexutil.Uint64(blockHeight)}, nil)
	ethClient.On("GetLogs", mock.Anything).
		Return(nil, nil)
	sub.On("Err").Return(nil)
	sub.On("Unsubscribe").Return()

	lb := ethsvc.NewLogBroadcaster(ethClient, store.ORM, 10)
	lb.Start()

	addr1 := cltest.NewAddress()
	addr2 := cltest.NewAddress()
	addr1SentLogs := []eth.Log{
		{Address: addr1, BlockNumber: 1},
		{Address: addr1, BlockNumber: 2},
		{Address: addr1, BlockNumber: 3},
	}
	addr2SentLogs := []eth.Log{
		{Address: addr2, BlockNumber: 4},
		{Address: addr2, BlockNumber: 5},
		{Address: addr2, BlockNumber: 6},
	}

	var addr1Logs1, addr1Logs2, addr2Logs1, addr2Logs2 []interface{}
	lb.Register(addr1, &funcLogListener{func(log interface{}, err error) {
		require.NoError(t, err)
		addr1Logs1 = append(addr1Logs1, log)
	}})
	lb.Register(addr1, &funcLogListener{func(log interface{}, err error) {
		require.NoError(t, err)
		addr1Logs2 = append(addr1Logs2, log)
	}})
	lb.Register(addr2, &funcLogListener{func(log interface{}, err error) {
		require.NoError(t, err)
		addr2Logs1 = append(addr2Logs1, log)
	}})
	lb.Register(addr2, &funcLogListener{func(log interface{}, err error) {
		require.NoError(t, err)
		addr2Logs2 = append(addr2Logs2, log)
	}})
	chRawLogs := <-chchRawLogs

	for _, log := range addr1SentLogs {
		chRawLogs <- log
	}
	for _, log := range addr2SentLogs {
		chRawLogs <- log
	}

	require.Eventually(t, func() bool { return len(addr1Logs1) == len(addr1SentLogs) }, time.Second, 10*time.Millisecond)
	require.Eventually(t, func() bool { return len(addr1Logs2) == len(addr1SentLogs) }, time.Second, 10*time.Millisecond)
	require.Eventually(t, func() bool { return len(addr2Logs1) == len(addr2SentLogs) }, time.Second, 10*time.Millisecond)
	require.Eventually(t, func() bool { return len(addr2Logs2) == len(addr2SentLogs) }, time.Second, 10*time.Millisecond)

	lb.Stop()

	for i := range addr1SentLogs {
		require.Equal(t, addr1SentLogs[i], addr1Logs1[i])
		require.Equal(t, addr1SentLogs[i], addr1Logs2[i])
	}
	for i := range addr2SentLogs {
		require.Equal(t, addr2SentLogs[i], addr2Logs1[i])
		require.Equal(t, addr2SentLogs[i], addr2Logs2[i])
	}

	ethClient.AssertExpectations(t)
	sub.AssertExpectations(t)
}

func TestLogBroadcaster_SkipsOldLogs(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	const blockHeight = 0

	ethClient := new(mocks.Client)
	sub := new(mocks.Subscription)

	chchRawLogs := make(chan chan<- eth.Log, 1)
	ethClient.On("SubscribeToLogs", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) { chchRawLogs <- args.Get(1).(chan<- eth.Log) }).
		Return(sub, nil).
		Once()
	ethClient.On("GetLatestBlock").
		Return(eth.Block{Number: hexutil.Uint64(blockHeight)}, nil)
	ethClient.On("GetLogs", mock.Anything).
		Return(nil, nil)

	sub.On("Unsubscribe").Return()
	sub.On("Err").Return(nil)

	lb := ethsvc.NewLogBroadcaster(ethClient, store.ORM, 10)
	lb.Start()

	addr := cltest.NewAddress()
	logs := []eth.Log{
		{Address: addr, BlockNumber: 1, Index: 0},
		{Address: addr, BlockNumber: 1, Index: 1},
		{Address: addr, BlockNumber: 1, Index: 2},
		{Address: addr, BlockNumber: 1, Index: 0}, // old log
		{Address: addr, BlockNumber: 2, Index: 0},
		{Address: addr, BlockNumber: 2, Index: 1},
		{Address: addr, BlockNumber: 2, Index: 1}, // old log
		{Address: addr, BlockNumber: 2, Index: 2},
		{Address: addr, BlockNumber: 3, Index: 0},
		{Address: addr, BlockNumber: 2, Index: 2}, // old log
		{Address: addr, BlockNumber: 3, Index: 1},
		{Address: addr, BlockNumber: 3, Index: 2},
	}

	var recvd []eth.Log

	lb.Register(addr, &funcLogListener{func(log interface{}, err error) {
		require.NoError(t, err)
		ethLog := log.(eth.Log)
		recvd = append(recvd, ethLog)
	}})

	chRawLogs := <-chchRawLogs

	for i := 0; i < len(logs); i++ {
		chRawLogs <- logs[i]
	}

	require.Eventually(t, func() bool { return len(recvd) == 9 }, 5*time.Second, 10*time.Millisecond)

	// check that all 9 received logs are unique
	recvdIdx := 0
	for blockNum := 1; blockNum <= 3; blockNum++ {
		for index := 0; index < 3; index++ {
			require.Equal(t, recvd[recvdIdx].BlockNumber, uint64(blockNum))
			require.Equal(t, recvd[recvdIdx].Index, uint(index))
			recvdIdx++
		}
	}

	ethClient.AssertExpectations(t)
}

func TestLogBroadcaster_Register_ResubscribesToMostRecentlySeenBlock(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	const (
		blockHeight   = 0
		expectedBlock = 3
	)

	ethClient := new(mocks.Client)
	sub := new(mocks.Subscription)

	addr1 := cltest.NewAddress()
	addr2 := cltest.NewAddress()

	chchRawLogs := make(chan chan<- eth.Log, 1)
	ethClient.On("SubscribeToLogs", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			chchRawLogs <- args.Get(1).(chan<- eth.Log)
		}).
		Return(sub, nil).
		Once()
	ethClient.On("SubscribeToLogs", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			query := args.Get(2).(ethereum.FilterQuery)
			require.Equal(t, big.NewInt(expectedBlock), query.FromBlock)
			require.Contains(t, query.Addresses, addr1)
			require.Contains(t, query.Addresses, addr2)
			require.Len(t, query.Addresses, 2)
			chchRawLogs <- args.Get(1).(chan<- eth.Log)
		}).
		Return(sub, nil).
		Once()

	ethClient.On("GetLatestBlock").
		Return(eth.Block{Number: hexutil.Uint64(blockHeight)}, nil)
	ethClient.On("GetLogs", mock.Anything).
		Return(nil, nil)

	sub.On("Unsubscribe").Return()
	sub.On("Err").Return(nil)

	listener1 := new(mocks.LogListener)
	listener2 := new(mocks.LogListener)
	listener1.On("OnConnect").Return()
	listener2.On("OnConnect").Return()
	listener1.On("OnDisconnect").Return()
	listener2.On("OnDisconnect").Return()

	lb := ethsvc.NewLogBroadcaster(ethClient, store.ORM, 10)
	lb.Start()                    // Subscribe #1
	lb.Register(addr1, listener1) // Subscribe #2
	chRawLogs := <-chchRawLogs
	chRawLogs <- eth.Log{BlockNumber: expectedBlock}
	lb.Register(addr2, listener2) // Subscribe #3
	<-chchRawLogs

	lb.Stop()

	ethClient.AssertExpectations(t)
	listener1.AssertExpectations(t)
	listener2.AssertExpectations(t)
	sub.AssertExpectations(t)
}

func TestDecodingLogListener(t *testing.T) {
	t.Parallel()

	contract, err := eth.GetV6ContractCodec("FluxAggregator")
	require.NoError(t, err)

	type LogNewRound struct {
		eth.Log
		RoundId   *big.Int
		StartedBy common.Address
		StartedAt *big.Int
	}

	logTypes := map[common.Hash]interface{}{
		eth.MustGetV6ContractEventID("FluxAggregator", "NewRound"): LogNewRound{},
	}

	var decodedLog interface{}
	listener := ethsvc.NewDecodingLogListener(contract, logTypes, &funcLogListener{func(decoded interface{}, innerErr error) {
		err = innerErr
		decodedLog = decoded
	}})
	rawLog := cltest.LogFromFixture(t, "../testdata/new_round_log.json")
	listener.HandleLog(rawLog, nil)
	require.NoError(t, err)
	newRoundLog := decodedLog.(*LogNewRound)
	require.Equal(t, newRoundLog.Log, rawLog)
	require.True(t, newRoundLog.RoundId.Cmp(big.NewInt(1)) == 0)
	require.Equal(t, newRoundLog.StartedBy, common.HexToAddress("f17f52151ebef6c7334fad080c5704d77216b732"))
	require.True(t, newRoundLog.StartedAt.Cmp(big.NewInt(15)) == 0)

	expectedErr := errors.New("oh no!")
	listener.HandleLog(nil, expectedErr)
	require.Equal(t, err, expectedErr)
}

func TestLogBroadcaster_ReceivesAllLogsWhenResubscribing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		blockHeight1     uint64
		blockHeight2     uint64
		batch1           []eth.Log
		backfillableLogs []eth.Log
		batch2           []eth.Log
		expectedFinal    []eth.Log
	}{
		{
			name:         "no backfilled logs, no overlap",
			blockHeight1: 0,
			blockHeight2: 2,
			batch1: []eth.Log{
				eth.Log{BlockNumber: 1, Index: 0},
				eth.Log{BlockNumber: 2, Index: 0},
			},
			backfillableLogs: nil,
			batch2: []eth.Log{
				eth.Log{BlockNumber: 4, Index: 0},
				eth.Log{BlockNumber: 5, Index: 0},
			},
			expectedFinal: []eth.Log{
				eth.Log{BlockNumber: 1, Index: 0},
				eth.Log{BlockNumber: 2, Index: 0},
				eth.Log{BlockNumber: 4, Index: 0},
				eth.Log{BlockNumber: 5, Index: 0},
			},
		},
		{
			name:         "no backfilled logs, overlap",
			blockHeight1: 0,
			blockHeight2: 2,
			batch1: []eth.Log{
				eth.Log{BlockNumber: 1, Index: 0},
				eth.Log{BlockNumber: 2, Index: 0},
			},
			backfillableLogs: nil,
			batch2: []eth.Log{
				eth.Log{BlockNumber: 2, Index: 0},
				eth.Log{BlockNumber: 3, Index: 0},
			},
			expectedFinal: []eth.Log{
				eth.Log{BlockNumber: 1, Index: 0},
				eth.Log{BlockNumber: 2, Index: 0},
				eth.Log{BlockNumber: 3, Index: 0},
			},
		},
		{
			name:         "backfilled logs, no overlap",
			blockHeight1: 0,
			blockHeight2: 15,
			batch1: []eth.Log{
				eth.Log{BlockNumber: 1, Index: 0},
				eth.Log{BlockNumber: 2, Index: 0},
			},
			backfillableLogs: []eth.Log{
				eth.Log{BlockNumber: 6, Index: 0},
				eth.Log{BlockNumber: 7, Index: 2},
				eth.Log{BlockNumber: 12, Index: 11},
				eth.Log{BlockNumber: 15, Index: 0},
			},
			batch2: []eth.Log{
				eth.Log{BlockNumber: 16, Index: 0},
				eth.Log{BlockNumber: 17, Index: 0},
			},
			expectedFinal: []eth.Log{
				eth.Log{BlockNumber: 1, Index: 0},
				eth.Log{BlockNumber: 2, Index: 0},
				eth.Log{BlockNumber: 6, Index: 0},
				eth.Log{BlockNumber: 7, Index: 2},
				eth.Log{BlockNumber: 12, Index: 11},
				eth.Log{BlockNumber: 15, Index: 0},
				eth.Log{BlockNumber: 16, Index: 0},
				eth.Log{BlockNumber: 17, Index: 0},
			},
		},
		{
			name:         "backfilled logs, overlap",
			blockHeight1: 0,
			blockHeight2: 15,
			batch1: []eth.Log{
				eth.Log{BlockNumber: 1, Index: 0},
				eth.Log{BlockNumber: 9, Index: 0},
			},
			backfillableLogs: []eth.Log{
				eth.Log{BlockNumber: 9, Index: 0},
				eth.Log{BlockNumber: 12, Index: 11},
				eth.Log{BlockNumber: 15, Index: 0},
			},
			batch2: []eth.Log{
				eth.Log{BlockNumber: 16, Index: 0},
				eth.Log{BlockNumber: 17, Index: 0},
			},
			expectedFinal: []eth.Log{
				eth.Log{BlockNumber: 1, Index: 0},
				eth.Log{BlockNumber: 9, Index: 0},
				eth.Log{BlockNumber: 12, Index: 11},
				eth.Log{BlockNumber: 15, Index: 0},
				eth.Log{BlockNumber: 16, Index: 0},
				eth.Log{BlockNumber: 17, Index: 0},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			sub := new(mocks.Subscription)
			ethClient := new(mocks.Client)

			chchRawLogs := make(chan chan<- eth.Log, 1)

			ethClient.On("SubscribeToLogs", mock.Anything, mock.Anything, mock.Anything).
				Run(func(args mock.Arguments) {
					chRawLogs := args.Get(1).(chan<- eth.Log)
					chchRawLogs <- chRawLogs
				}).
				Return(sub, nil).
				Twice()

			ethClient.On("GetLatestBlock").Return(eth.Block{Number: hexutil.Uint64(test.blockHeight1)}, nil).Times(3)
			ethClient.On("GetLatestBlock").Return(eth.Block{Number: hexutil.Uint64(test.blockHeight2)}, nil).Once()
			ethClient.On("GetLogs", mock.Anything).Return(nil, nil).Twice()
			ethClient.On("GetLogs", mock.Anything).Return(test.backfillableLogs, nil).Once()

			sub.On("Err").Return(nil)
			sub.On("Unsubscribe").Return()

			lb := ethsvc.NewLogBroadcaster(ethClient, store.ORM, 10)
			lb.Start()

			var recvd []eth.Log
			logListener := &funcLogListener{
				fn: func(log interface{}, err error) { recvd = append(recvd, log.(eth.Log)) },
			}

			// Send initial logs
			lb.Register(common.Address{0}, logListener)
			chRawLogs1 := <-chchRawLogs
			for _, log := range test.batch1 {
				chRawLogs1 <- log
			}
			require.Eventually(t, func() bool { return len(recvd) == len(test.batch1) }, 5*time.Second, 10*time.Millisecond)
			for i, log := range test.batch1 {
				require.Equal(t, test.batch1[i], log)
			}

			// Trigger resubscription
			lb.Register(common.Address{1}, &funcLogListener{})
			chRawLogs2 := <-chchRawLogs
			for _, log := range test.batch2 {
				chRawLogs2 <- log
			}
			require.Eventually(t, func() bool { return len(recvd) == len(test.expectedFinal) }, 5*time.Second, 10*time.Millisecond)
			for i, log := range test.expectedFinal {
				require.Equal(t, test.expectedFinal[i], log)
			}

			lb.Stop()
		})
	}
}

func TestAppendLogChannel(t *testing.T) {
	t.Parallel()

	logs1 := []eth.Log{
		{BlockNumber: 1},
		{BlockNumber: 2},
		{BlockNumber: 3},
		{BlockNumber: 4},
		{BlockNumber: 5},
	}

	logs2 := []eth.Log{
		{BlockNumber: 6},
		{BlockNumber: 7},
		{BlockNumber: 8},
		{BlockNumber: 9},
		{BlockNumber: 10},
	}

	logs3 := []eth.Log{
		{BlockNumber: 11},
		{BlockNumber: 12},
		{BlockNumber: 13},
		{BlockNumber: 14},
		{BlockNumber: 15},
	}

	ch1 := make(chan eth.Log)
	ch2 := make(chan eth.Log)
	ch3 := make(chan eth.Log)

	chCombined := ethsvc.ExposedAppendLogChannel(ch1, ch2)
	chCombined = ethsvc.ExposedAppendLogChannel(chCombined, ch3)

	go func() {
		defer close(ch1)
		for _, log := range logs1 {
			ch1 <- log
		}
	}()
	go func() {
		defer close(ch2)
		for _, log := range logs2 {
			ch2 <- log
		}
	}()
	go func() {
		defer close(ch3)
		for _, log := range logs3 {
			ch3 <- log
		}
	}()

	expected := append(logs1, logs2...)
	expected = append(expected, logs3...)

	var i int
	for log := range chCombined {
		require.Equal(t, expected[i], log)
		i++
	}
}
