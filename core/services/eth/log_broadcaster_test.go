package eth_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/eth"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	ethsvc "github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/eth/contracts"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLogBroadcaster_ResubscribesOnAddOrRemoveContract(t *testing.T) {
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
	ethClient.On("GetBlockHeight").
		Return(blockHeight, nil)
	sub.On("Unsubscribe").
		Return().
		Run(func(mock.Arguments) { unsubscribeCalls++ })
	sub.On("Err").Return(nil)

	lb := ethsvc.NewLogBroadcaster(ethClient, store.ORM)
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
	ethClient.On("GetBlockHeight").Return(blockHeight, nil)
	sub.On("Err").Return(nil)
	sub.On("Unsubscribe").Return()

	lb := ethsvc.NewLogBroadcaster(ethClient, store.ORM)
	lb.Start()

	addr1 := cltest.NewAddress()
	addr2 := cltest.NewAddress()
	addr1SentLogs := []eth.Log{
		{Address: addr1, BlockNumber: 0},
		{Address: addr1, BlockNumber: 1},
		{Address: addr1, BlockNumber: 2},
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
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	ethClient := new(mocks.Client)
	sub := new(mocks.Subscription)

	ethClient.On("GetBlockHeight").
		Return(uint64(0), nil)
	chchRawLogs := make(chan chan<- eth.Log, 1)
	ethClient.On("SubscribeToLogs", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) { chchRawLogs <- args.Get(1).(chan<- eth.Log) }).
		Return(sub, nil).
		Once()

	sub.On("Unsubscribe").Return()
	sub.On("Err").Return(nil)

	lb := ethsvc.NewLogBroadcaster(ethClient, store.ORM)
	lb.Start()

	addr := cltest.NewAddress()
	logs := []eth.Log{
		{Address: addr, BlockNumber: 0, Index: 0},
		{Address: addr, BlockNumber: 0, Index: 1},
		{Address: addr, BlockNumber: 0, Index: 2},
		{Address: addr, BlockNumber: 1, Index: 0},
		{Address: addr, BlockNumber: 1, Index: 1},
		{Address: addr, BlockNumber: 1, Index: 2},
		{Address: addr, BlockNumber: 2, Index: 0},
		{Address: addr, BlockNumber: 2, Index: 1},
		{Address: addr, BlockNumber: 2, Index: 2},
	}

	var recvd []interface{}
	lb.Register(addr, &funcLogListener{func(log interface{}, err error) {
		require.NoError(t, err)
		recvd = append(recvd, log)
	}})

	chRawLogs := <-chchRawLogs

	// Simulates resuming the subscription repeatedly as new blocks are coming in
	for i := 0; i < len(logs); i++ {
		for _, log := range logs[0 : i+1] {
			chRawLogs <- log
		}
	}

	lb.Stop() // This should ensure that all sending is complete

	require.Len(t, recvd, len(logs))
	for i := range recvd {
		require.Equal(t, recvd[i], logs[i])
	}

	ethClient.AssertExpectations(t)
}

func TestLogBroadcaster_Register_ResubscribesToMostRecentlySeenBlock(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	const expectedBlock = 3

	ethClient := new(mocks.Client)
	sub := new(mocks.Subscription)

	addr1 := cltest.NewAddress()
	addr2 := cltest.NewAddress()

	ethClient.On("GetBlockHeight").Return(uint64(0), nil)
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

	sub.On("Unsubscribe").Return()
	sub.On("Err").Return(nil)

	listener1 := new(mocks.LogListener)
	listener2 := new(mocks.LogListener)
	listener1.On("OnConnect").Return()
	listener2.On("OnConnect").Return()
	listener1.On("OnDisconnect").Return()
	listener2.On("OnDisconnect").Return()

	lb := ethsvc.NewLogBroadcaster(ethClient, store.ORM)
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
	contract, err := eth.GetV6ContractCodec("FluxAggregator")
	require.NoError(t, err)

	logTypes := map[common.Hash]interface{}{
		eth.MustGetV6ContractEventID("FluxAggregator", "NewRound"): contracts.LogNewRound{},
	}

	var decodedLog interface{}
	listener := ethsvc.NewDecodingLogListener(contract, logTypes, &funcLogListener{func(decoded interface{}, innerErr error) {
		err = innerErr
		decodedLog = decoded
	}})
	rawLog := cltest.LogFromFixture(t, "../testdata/new_round_log.json")
	listener.HandleLog(rawLog, nil)
	require.NoError(t, err)
	newRoundLog := decodedLog.(*contracts.LogNewRound)
	require.Equal(t, newRoundLog.Log, rawLog)
	require.True(t, newRoundLog.RoundId.Cmp(big.NewInt(1)) == 0)
	require.Equal(t, newRoundLog.StartedBy, common.HexToAddress("f17f52151ebef6c7334fad080c5704d77216b732"))
	require.True(t, newRoundLog.StartedAt.Cmp(big.NewInt(15)) == 0)

	expectedErr := errors.New("oh no!")
	listener.HandleLog(nil, expectedErr)
	require.Equal(t, err, expectedErr)
}

func TestLogBroadcaster_ReceivesAllLogsWhenResubscribing(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	const blockHeight uint64 = 0

	ethClient := new(mocks.Client)
	sub := new(mocks.Subscription)

	chchRawLogs := make(chan chan<- eth.Log, 1)

	ethClient.On("SubscribeToLogs", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			chRawLogs := args.Get(1).(chan<- eth.Log)
			chchRawLogs <- chRawLogs
		}).
		Return(sub, nil).
		Twice()

	ethClient.On("GetBlockHeight").Return(blockHeight, nil)
	sub.On("Err").Return(nil)
	sub.On("Unsubscribe").Return()

	lb := ethsvc.NewLogBroadcaster(ethClient, store.ORM)
	lb.Start()

	logCount := 0
	logListener := funcLogListener{
		fn: func(log interface{}, err error) { logCount++ },
	}
	logListener2 := funcLogListener{
		fn: func(log interface{}, err error) {},
	}

	lb.Register(common.Address{}, &logListener)
	chRawLogs1 := <-chchRawLogs
	chRawLogs1 <- eth.Log{BlockNumber: 0, Index: 0}
	chRawLogs1 <- eth.Log{BlockNumber: 1, Index: 0}

	lb.Register(common.Address{1}, &logListener2) // trigger resubscription
	chRawLogs2 := <-chchRawLogs
	chRawLogs2 <- eth.Log{BlockNumber: 1, Index: 0} // send overlapping logs
	chRawLogs2 <- eth.Log{BlockNumber: 2, Index: 0}

	require.Eventually(t, func() bool { return logCount == 3 }, 5*time.Second, 10*time.Millisecond)
}
