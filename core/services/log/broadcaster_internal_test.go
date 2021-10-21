package log

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/logger"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	ethmocks "github.com/smartcontractkit/chainlink/core/services/eth/mocks"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type tc struct {
}

func (tc) BlockBackfillDepth() uint64 {
	return 0
}
func (tc) BlockBackfillSkip() bool {
	return true
}
func (tc) EvmFinalityDepth() uint {
	return 1
}
func (tc) EvmLogBackfillBatchSize() uint32 {
	return 1
}

type listener struct {
	logs chan Broadcast
}

func (l listener) HandleLog(b Broadcast) {
	l.logs <- b
}

func (l listener) JobID() int32 {
	return 1
}

type sub struct {
}

func (s sub) Unsubscribe() {
}

func (s sub) Err() <-chan error {
	return nil
}

func TestBroadcaster_BroadcastsWithZeroConfirmations(t *testing.T) {
	gm := gomega.NewGomegaWithT(t)
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
	lb := NewBroadcaster(dborm, ec, tc{}, logger.Default, nil)
	lb.Start()
	defer lb.Close()

	addr := common.HexToAddress("0xf0d54349aDdcf704F77AE15b96510dEA15cb7952")
	contract1, err := flux_aggregator_wrapper.NewFluxAggregator(addr, nil)
	require.NoError(t, err)

	// 3 logs all in the same block
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

	// Give these listeners a big buffer of logs
	// in case the log broadcaster erroneously sends logs more than once.
	broadcastsToListener1 := make(chan Broadcast, 100)
	broadcastsToListener2 := make(chan Broadcast, 100)
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
	// Wait until the logpool has the 3 logs
	gm.Eventually(func() bool {
		return len(lb.logPool.logsByBlockHash[bh]) == len(addr1SentLogs)
	}, 2*time.Second, 100*time.Millisecond).Should(gomega.BeTrue())

	// Send a block to trigger sending the logs from the pool
	// to the subscribers
	lb.OnNewLongestChain(context.Background(), models.Head{
		Number: 2,
	})

	// The subs should each get exactly 3 broadcasts each
	// If we do not receive a broadcast for 1 second
	// we assume the log broadcaster is done sending.
	gm.Eventually(func() bool {
		return len(broadcastsToListener1) == len(addr1SentLogs) && len(broadcastsToListener2) == len(addr1SentLogs)
	}, 2*time.Second, 100*time.Millisecond).Should(gomega.BeTrue())
	gm.Consistently(func() bool {
		return len(broadcastsToListener1) == len(addr1SentLogs) && len(broadcastsToListener2) == len(addr1SentLogs)
	}, 1*time.Second).Should(gomega.BeTrue())
}
