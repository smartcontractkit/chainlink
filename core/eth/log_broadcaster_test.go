package eth_test

import (
	"math/big"
	"testing"
	"time"

	"chainlink/core/eth"
	"chainlink/core/eth/contracts"
	"chainlink/core/internal/cltest"
	"chainlink/core/internal/mocks"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLogBroadcaster_ResubscribesOnAddOrRemoveContract(t *testing.T) {
	const numContracts = 3

	ethClient := new(mocks.Client)
	sub := new(mocks.Subscription)

	var subscribeCalls int
	var unsubscribeCalls int
	ethClient.On("SubscribeToLogs", mock.Anything, mock.Anything).
		Return(sub, nil).
		Run(func(mock.Arguments) { subscribeCalls++ })
	sub.On("Unsubscribe").
		Return().
		Run(func(mock.Arguments) { unsubscribeCalls++ })

	lb := eth.NewLogBroadcaster(ethClient)
	lb.Start()

	type registration struct {
		common.Address
		eth.LogListener
	}
	registrations := make([]registration, numContracts)
	for i := 0; i < numContracts; i++ {
		registrations[i] = registration{cltest.NewAddress(), new(mocks.LogListener)}
		lb.Register(registrations[i].Address, registrations[i].LogListener)
	}
	require.Eventually(t, func() bool { return subscribeCalls == numContracts }, time.Second, 10*time.Millisecond)
	require.Eventually(t, func() bool { return unsubscribeCalls == numContracts-1 }, time.Second, 10*time.Millisecond)

	for _, r := range registrations {
		lb.Unregister(r.Address, r.LogListener)
	}
	require.Eventually(t, func() bool { return subscribeCalls == (2*numContracts)-1 }, time.Second, 10*time.Millisecond)
	require.Eventually(t, func() bool { return unsubscribeCalls == (2*numContracts)-1 }, time.Second, 10*time.Millisecond)

	lb.Stop()
	require.Eventually(t, func() bool { return unsubscribeCalls == (2*numContracts)-1 }, time.Second, 10*time.Millisecond)

	ethClient.AssertExpectations(t)
	sub.AssertExpectations(t)
}

func TestLogBroadcaster_BroadcastsToCorrectRecipients(t *testing.T) {
	ethClient := new(mocks.Client)
	sub := new(mocks.Subscription)

	chchRawLogs := make(chan chan<- eth.Log, 1)
	ethClient.On("SubscribeToLogs", mock.Anything, mock.Anything).
		Return(sub, nil).
		Once()
	ethClient.On("SubscribeToLogs", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) { chchRawLogs <- args.Get(0).(chan<- eth.Log) }).
		Return(sub, nil).
		Once()

	sub.On("Unsubscribe").Return()

	lb := eth.NewLogBroadcaster(ethClient)
	lb.Start()
	defer lb.Stop()

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
	lb.Register(addr1, eth.NewFuncLogListener(func(log interface{}, err error) {
		require.NoError(t, err)
		addr1Logs1 = append(addr1Logs1, log)
	}))
	lb.Register(addr1, eth.NewFuncLogListener(func(log interface{}, err error) {
		require.NoError(t, err)
		addr1Logs2 = append(addr1Logs2, log)
	}))
	lb.Register(addr2, eth.NewFuncLogListener(func(log interface{}, err error) {
		require.NoError(t, err)
		addr2Logs1 = append(addr2Logs1, log)
	}))
	lb.Register(addr2, eth.NewFuncLogListener(func(log interface{}, err error) {
		require.NoError(t, err)
		addr2Logs2 = append(addr2Logs2, log)
	}))
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

	for i := range addr1SentLogs {
		require.Equal(t, addr1SentLogs[i], addr1Logs1[i])
		require.Equal(t, addr1SentLogs[i], addr1Logs2[i])
	}
	for i := range addr2SentLogs {
		require.Equal(t, addr2SentLogs[i], addr2Logs1[i])
		require.Equal(t, addr2SentLogs[i], addr2Logs2[i])
	}

	ethClient.AssertExpectations(t)
}

func TestLogBroadcaster_SkipsOldLogs(t *testing.T) {
	ethClient := new(mocks.Client)
	sub := new(mocks.Subscription)

	chchRawLogs := make(chan chan<- eth.Log, 1)
	ethClient.On("SubscribeToLogs", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) { chchRawLogs <- args.Get(0).(chan<- eth.Log) }).
		Return(sub, nil).
		Once()

	sub.On("Unsubscribe").Return()

	lb := eth.NewLogBroadcaster(ethClient)
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
	lb.Register(addr, eth.NewFuncLogListener(func(log interface{}, err error) {
		require.NoError(t, err)
		recvd = append(recvd, log)
	}))

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

func TestLogBroadcaster_ResubscribesToMostRecentlySeenBlock(t *testing.T) {
	const expectedBlock = 3

	ethClient := new(mocks.Client)
	sub := new(mocks.Subscription)

	addr1 := cltest.NewAddress()
	addr2 := cltest.NewAddress()

	chchRawLogs := make(chan chan<- eth.Log, 1)
	ethClient.On("SubscribeToLogs", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			chchRawLogs <- args.Get(0).(chan<- eth.Log)
		}).
		Return(sub, nil).
		Once()
	ethClient.On("SubscribeToLogs", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			query := args.Get(1).(ethereum.FilterQuery)
			require.Equal(t, big.NewInt(expectedBlock), query.FromBlock)
			require.Contains(t, query.Addresses, addr1)
			require.Contains(t, query.Addresses, addr2)
			require.Len(t, query.Addresses, 2)
		}).
		Return(sub, nil).
		Once()

	sub.On("Unsubscribe").Return()

	lb := eth.NewLogBroadcaster(ethClient)
	lb.Start()                                 // Subscribe #1
	lb.Register(addr1, new(mocks.LogListener)) // Subscribe #2
	chRawLogs := <-chchRawLogs
	chRawLogs <- eth.Log{BlockNumber: expectedBlock}
	lb.Register(addr2, new(mocks.LogListener)) // Subscribe #3

	lb.Stop()

	ethClient.AssertExpectations(t)
}

func TestDecodingLogListener(t *testing.T) {
	contract, err := eth.GetV6Contract("FluxAggregator")
	require.NoError(t, err)

	logTypes := map[common.Hash]interface{}{
		eth.MustGetV6ContractEventID("FluxAggregator", "NewRound"):      contracts.LogNewRound{},
		eth.MustGetV6ContractEventID("FluxAggregator", "AnswerUpdated"): contracts.LogAnswerUpdated{},
	}

	var decodedLog interface{}
	listener := eth.NewDecodingLogListener(contract, logTypes, func(decoded interface{}, innerErr error) {
		err = innerErr
		decodedLog = decoded
	})
	rawLog := cltest.LogFromFixture(t, "../services/testdata/new_round_log.json")
	listener.HandleLog(rawLog, nil)
	require.NoError(t, err)
	newRoundLog := decodedLog.(*contracts.LogNewRound)
	require.Equal(t, newRoundLog.Log, rawLog)
	require.True(t, newRoundLog.RoundId.Cmp(big.NewInt(1)) == 0)
	require.Equal(t, newRoundLog.StartedBy, common.HexToAddress("f17f52151ebef6c7334fad080c5704d77216b732"))
	require.True(t, newRoundLog.StartedAt.Cmp(big.NewInt(9)) == 0)

	rawLog = cltest.LogFromFixture(t, "../services/testdata/answer_updated_log.json")
	listener.HandleLog(rawLog, nil)
	require.NoError(t, err)
	answerUpdatedLog := decodedLog.(*contracts.LogAnswerUpdated)
	require.Equal(t, answerUpdatedLog.Log, rawLog)
	require.True(t, answerUpdatedLog.Current.Cmp(big.NewInt(1)) == 0)
	require.True(t, answerUpdatedLog.RoundId.Cmp(big.NewInt(2)) == 0)
	require.True(t, answerUpdatedLog.Timestamp.Cmp(big.NewInt(3)) == 0)

	expectedErr := errors.New("oh no!")
	listener.HandleLog(nil, expectedErr)
	require.Equal(t, err, expectedErr)
}
