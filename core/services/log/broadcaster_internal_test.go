package log

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	ethmocks "github.com/smartcontractkit/chainlink/core/services/eth/mocks"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

type tc struct {
}

func (tc) BlockBackfillDepth() uint64 {
	return 0
}
func (tc) BlockBackfillSkip() bool {
	return true
}
func (tc) EthFinalityDepth() uint {
	return 1
}
func (tc) EthLogBackfillBatchSize() uint32 {
	return 1
}

type listener struct {
	logs chan Broadcast
}

func (l listener) HandleLog(b Broadcast) {
	l.logs <- b
}

func (l listener) JobID() models.JobID {
	return models.NewJobID()
}

func (l listener) JobIDV2() int32 {
	return 1
}

func (l listener) IsV2Job() bool {
	return true
}

type sub struct {
}

func (s sub) Unsubscribe() {
}

func (s sub) Err() <-chan error {
	return nil
}

func TestBroadcaster_BroadcastsWithZeroConfirmations(t *testing.T) {
	logsChCh := make(chan chan<- types.Log)
	ec := new(ethmocks.Client)
	ec.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			logsChCh <- args.Get(2).(chan<- types.Log)
		}).
		Return(sub{}, nil)
	ec.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).
		Return(&models.Head{Number: 1}, nil)
	ec.On("FilterLogs", mock.Anything, mock.Anything).
		Return(nil, nil)
	db := pgtest.NewGormDB(t)
	dborm := NewORM(db)
	lb := NewBroadcaster(dborm, ec, tc{}, nil)
	lb.Start()
	// TODO: make this not hang
	//defer lb.Close()

	addr := common.HexToAddress("0xf0d54349aDdcf704F77AE15b96510dEA15cb7952")
	contract1, err := flux_aggregator_wrapper.NewFluxAggregator(addr, nil)
	require.NoError(t, err)

	bh := utils.NewHash()
	addr1SentLogs := []types.Log{
		{
			Address:     addr,
			BlockHash:   bh,
			BlockNumber: 2,
			Index:       0,
			Topics: []common.Hash{
				(flux_aggregator_wrapper.FluxAggregatorNewRound{}).Topic(),
				utils.NewHash(),
				utils.NewHash(),
			},
			Data: []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		},
		{
			Address:     addr,
			BlockHash:   bh,
			BlockNumber: 2,
			Index:       1,
			Topics: []common.Hash{
				(flux_aggregator_wrapper.FluxAggregatorNewRound{}).Topic(),
				utils.NewHash(),
				utils.NewHash(),
			},
			Data: []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		},
		{
			Address:     addr,
			BlockHash:   bh,
			BlockNumber: 2,
			Index:       2,
			Topics: []common.Hash{
				(flux_aggregator_wrapper.FluxAggregatorNewRound{}).Topic(),
				utils.NewHash(),
				utils.NewHash(),
			},
			Data: []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		},
	}

	broadcastsToListener1 := make(chan Broadcast, 3)
	broadcastsToListener2 := make(chan Broadcast, 3)
	lt := make(map[common.Hash][][]Topic)
	lt[flux_aggregator_wrapper.FluxAggregatorNewRound{}.Topic()] = nil
	lt[flux_aggregator_wrapper.FluxAggregatorAnswerUpdated{}.Topic()] = nil
	lb.Register(listener{broadcastsToListener1}, ListenerOpts{
		Contract:         addr,
		LogsWithTopics:   lt,
		ParseLog:         contract1.ParseLog,
		NumConfirmations: 0,
	})
	lb.Register(listener{broadcastsToListener2}, ListenerOpts{
		Contract:         addr,
		LogsWithTopics:   lt,
		ParseLog:         contract1.ParseLog,
		NumConfirmations: 0,
	})
	logs := <-logsChCh

	for _, log := range addr1SentLogs {
		select {
		case logs <- log:
		case <-time.After(time.Second):
			t.Error("failed to send log to log broadcaster")
		}
	}
	for i := 0; i < 10; i++ {
		if len(lb.logPool.logsByBlockHash[bh]) == len(addr1SentLogs) {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	// All logs should be sitting in the the lb
	// Send a head to fire them off
	lb.OnNewLongestChain(context.Background(), models.Head{
		Number: 2,
	})
	for i := 0; i < 2*len(addr1SentLogs); i++ {
		select {
		case <-broadcastsToListener1:
		case <-broadcastsToListener2:
		case <-time.After(5 * time.Second):
			t.Error("failed to get broadcasts")
		}
	}
	// TODO: I think there may be a bug in getLogsToSend
	// I'm seeing 14 "sending logs" debugs when there should only be 6?
}
