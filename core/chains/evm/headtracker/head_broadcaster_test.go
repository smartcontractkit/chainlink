package headtracker_test

import (
	"context"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	evmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func waitHeadBroadcasterToStart(t *testing.T, hb types.HeadBroadcaster) {
	t.Helper()

	subscriber := &cltest.MockHeadTrackable{}
	_, unsubscribe := hb.Subscribe(subscriber)
	defer unsubscribe()

	hb.BroadcastNewLongestChain(cltest.Head(1))
	g := gomega.NewWithT(t)
	g.Eventually(subscriber.OnNewLongestChainCount).Should(gomega.Equal(int32(1)))
}

func TestHeadBroadcaster_Subscribe(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].HeadTracker.SamplingInterval = &models.Duration{}
	})
	evmCfg := evmtest.NewChainScopedConfig(t, cfg)
	db := pgtest.NewSqlxDB(t)
	logger := logger.TestLogger(t)

	sub := new(evmmocks.Subscription)
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

	chchHeaders := make(chan chan<- *evmtypes.Head, 1)
	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			chchHeaders <- args.Get(1).(chan<- *evmtypes.Head)
		}).
		Return(sub, nil)
	ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(cltest.Head(1), nil)

	sub.On("Unsubscribe").Return()
	sub.On("Err").Return(nil)

	checker1 := &cltest.MockHeadTrackable{}
	checker2 := &cltest.MockHeadTrackable{}

	hb := headtracker.NewHeadBroadcaster(logger)
	orm := headtracker.NewORM(db, logger, cfg, *ethClient.ChainID())
	hs := headtracker.NewHeadSaver(logger, orm, evmCfg)
	mailMon := utils.NewMailboxMonitor(t.Name())
	ht := headtracker.NewHeadTracker(logger, ethClient, evmCfg, hb, hs, mailMon)
	var ms services.MultiStart
	require.NoError(t, ms.Start(testutils.Context(t), mailMon, hb, ht))
	t.Cleanup(func() { require.NoError(t, services.MultiClose{mailMon, hb, ht}.Close()) })

	latest1, unsubscribe1 := hb.Subscribe(checker1)
	// "latest head" is nil here because we didn't receive any yet
	assert.Equal(t, (*evmtypes.Head)(nil), latest1)

	headers := <-chchHeaders
	h := evmtypes.Head{Number: 1, Hash: utils.NewHash(), ParentHash: utils.NewHash(), EVMChainID: utils.NewBig(&cltest.FixtureChainID)}
	headers <- &h
	g.Eventually(checker1.OnNewLongestChainCount).Should(gomega.Equal(int32(1)))

	latest2, _ := hb.Subscribe(checker2)
	// "latest head" is set here to the most recent head received
	assert.NotNil(t, latest2)
	assert.Equal(t, h.Number, latest2.Number)

	unsubscribe1()

	headers <- &evmtypes.Head{Number: 2, Hash: utils.NewHash(), ParentHash: h.Hash, EVMChainID: utils.NewBig(&cltest.FixtureChainID)}
	g.Eventually(checker2.OnNewLongestChainCount).Should(gomega.Equal(int32(1)))
}

func TestHeadBroadcaster_BroadcastNewLongestChain(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	lggr := logger.TestLogger(t)
	broadcaster := headtracker.NewHeadBroadcaster(lggr)

	err := broadcaster.Start(testutils.Context(t))
	require.NoError(t, err)

	waitHeadBroadcasterToStart(t, broadcaster)

	subscriber1 := &cltest.MockHeadTrackable{}
	subscriber2 := &cltest.MockHeadTrackable{}
	_, unsubscribe1 := broadcaster.Subscribe(subscriber1)
	_, unsubscribe2 := broadcaster.Subscribe(subscriber2)

	broadcaster.BroadcastNewLongestChain(cltest.Head(1))
	g.Eventually(subscriber1.OnNewLongestChainCount).Should(gomega.Equal(int32(1)))

	unsubscribe1()

	broadcaster.BroadcastNewLongestChain(cltest.Head(2))
	g.Eventually(subscriber2.OnNewLongestChainCount).Should(gomega.Equal(int32(2)))

	unsubscribe2()

	subscriber3 := &cltest.MockHeadTrackable{}
	_, unsubscribe3 := broadcaster.Subscribe(subscriber3)
	broadcaster.BroadcastNewLongestChain(cltest.Head(1))
	g.Eventually(subscriber3.OnNewLongestChainCount).Should(gomega.Equal(int32(1)))

	unsubscribe3()

	// no subscribers - shall do nothing
	broadcaster.BroadcastNewLongestChain(cltest.Head(0))

	err = broadcaster.Close()
	require.NoError(t, err)

	require.Equal(t, int32(1), subscriber3.OnNewLongestChainCount())
}

func TestHeadBroadcaster_TrackableCallbackTimeout(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	broadcaster := headtracker.NewHeadBroadcaster(lggr)

	err := broadcaster.Start(testutils.Context(t))
	require.NoError(t, err)

	waitHeadBroadcasterToStart(t, broadcaster)

	slowAwaiter := cltest.NewAwaiter()
	fastAwaiter := cltest.NewAwaiter()
	slow := &sleepySubscriber{awaiter: slowAwaiter, delay: headtracker.TrackableCallbackTimeout * 2}
	fast := &sleepySubscriber{awaiter: fastAwaiter, delay: headtracker.TrackableCallbackTimeout / 2}
	_, unsubscribe1 := broadcaster.Subscribe(slow)
	_, unsubscribe2 := broadcaster.Subscribe(fast)

	broadcaster.BroadcastNewLongestChain(cltest.Head(1))
	slowAwaiter.AwaitOrFail(t, testutils.WaitTimeout(t))
	fastAwaiter.AwaitOrFail(t, testutils.WaitTimeout(t))

	require.True(t, slow.contextDone)
	require.False(t, fast.contextDone)

	unsubscribe1()
	unsubscribe2()

	err = broadcaster.Close()
	require.NoError(t, err)
}

type sleepySubscriber struct {
	awaiter     cltest.Awaiter
	delay       time.Duration
	contextDone bool
}

func (ss *sleepySubscriber) OnNewLongestChain(ctx context.Context, head *evmtypes.Head) {
	time.Sleep(ss.delay)
	select {
	case <-ctx.Done():
		ss.contextDone = true
	default:
	}
	ss.awaiter.ItHappened()
}
