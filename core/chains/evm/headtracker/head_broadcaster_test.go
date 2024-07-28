package headtracker_test

import (
	"context"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox/mailboxtest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	commonhtrk "github.com/smartcontractkit/chainlink/v2/common/headtracker"
	commonmocks "github.com/smartcontractkit/chainlink/v2/common/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/testutils"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

func waitHeadBroadcasterToStart(t *testing.T, hb types.HeadBroadcaster) {
	t.Helper()

	subscriber := &mocks.MockHeadTrackable{}
	_, unsubscribe := hb.Subscribe(subscriber)
	defer unsubscribe()

	hb.BroadcastNewLongestChain(testutils.Head(1))
	g := gomega.NewWithT(t)
	g.Eventually(subscriber.OnNewLongestChainCount).Should(gomega.Equal(int32(1)))
}

func TestHeadBroadcaster_Subscribe(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	evmCfg := testutils.NewTestChainScopedConfig(t, func(c *toml.EVMConfig) {
		c.HeadTracker.SamplingInterval = &commonconfig.Duration{}
	})
	db := pgtest.NewSqlxDB(t)
	logger := logger.Test(t)

	sub := commonmocks.NewSubscription(t)
	ethClient := testutils.NewEthClientMockWithDefaultChain(t)

	chchHeaders := make(chan chan<- *evmtypes.Head, 1)
	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			chchHeaders <- args.Get(1).(chan<- *evmtypes.Head)
		}).
		Return(sub, nil)
	ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(testutils.Head(1), nil)

	sub.On("Unsubscribe").Return()
	sub.On("Err").Return(nil)

	checker1 := &mocks.MockHeadTrackable{}
	checker2 := &mocks.MockHeadTrackable{}

	orm := headtracker.NewORM(*ethClient.ConfiguredChainID(), db)
	hs := headtracker.NewHeadSaver(logger, orm, evmCfg.EVM(), evmCfg.EVM().HeadTracker())
	mailMon := mailboxtest.NewMonitor(t)
	servicetest.Run(t, mailMon)
	hb := headtracker.NewHeadBroadcaster(logger)
	servicetest.Run(t, hb)
	ht := headtracker.NewHeadTracker(logger, ethClient, evmCfg.EVM(), evmCfg.EVM().HeadTracker(), hb, hs, mailMon)
	servicetest.Run(t, ht)

	latest1, unsubscribe1 := hb.Subscribe(checker1)
	// "latest head" is nil here because we didn't receive any yet
	assert.Equal(t, (*evmtypes.Head)(nil), latest1)

	headers := <-chchHeaders
	h := evmtypes.Head{Number: 1, Hash: utils.NewHash(), ParentHash: utils.NewHash(), EVMChainID: big.New(testutils.FixtureChainID)}
	headers <- &h
	g.Eventually(checker1.OnNewLongestChainCount).Should(gomega.Equal(int32(1)))

	latest2, _ := hb.Subscribe(checker2)
	// "latest head" is set here to the most recent head received
	assert.NotNil(t, latest2)
	assert.Equal(t, h.Number, latest2.Number)

	unsubscribe1()

	headers <- &evmtypes.Head{Number: 2, Hash: utils.NewHash(), ParentHash: h.Hash, EVMChainID: big.New(testutils.FixtureChainID)}
	g.Eventually(checker2.OnNewLongestChainCount).Should(gomega.Equal(int32(1)))
}

func TestHeadBroadcaster_BroadcastNewLongestChain(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	lggr := logger.Test(t)
	broadcaster := headtracker.NewHeadBroadcaster(lggr)

	err := broadcaster.Start(tests.Context(t))
	require.NoError(t, err)

	waitHeadBroadcasterToStart(t, broadcaster)

	subscriber1 := &mocks.MockHeadTrackable{}
	subscriber2 := &mocks.MockHeadTrackable{}
	_, unsubscribe1 := broadcaster.Subscribe(subscriber1)
	_, unsubscribe2 := broadcaster.Subscribe(subscriber2)

	broadcaster.BroadcastNewLongestChain(testutils.Head(1))
	g.Eventually(subscriber1.OnNewLongestChainCount).Should(gomega.Equal(int32(1)))

	unsubscribe1()

	broadcaster.BroadcastNewLongestChain(testutils.Head(2))
	g.Eventually(subscriber2.OnNewLongestChainCount).Should(gomega.Equal(int32(2)))

	unsubscribe2()

	subscriber3 := &mocks.MockHeadTrackable{}
	_, unsubscribe3 := broadcaster.Subscribe(subscriber3)
	broadcaster.BroadcastNewLongestChain(testutils.Head(1))
	g.Eventually(subscriber3.OnNewLongestChainCount).Should(gomega.Equal(int32(1)))

	unsubscribe3()

	// no subscribers - shall do nothing
	broadcaster.BroadcastNewLongestChain(testutils.Head(0))

	err = broadcaster.Close()
	require.NoError(t, err)

	require.Equal(t, int32(1), subscriber3.OnNewLongestChainCount())
}

func TestHeadBroadcaster_TrackableCallbackTimeout(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	broadcaster := headtracker.NewHeadBroadcaster(lggr)

	err := broadcaster.Start(tests.Context(t))
	require.NoError(t, err)

	waitHeadBroadcasterToStart(t, broadcaster)

	slowAwaiter := testutils.NewAwaiter()
	fastAwaiter := testutils.NewAwaiter()
	slow := &sleepySubscriber{awaiter: slowAwaiter, delay: commonhtrk.TrackableCallbackTimeout * 2}
	fast := &sleepySubscriber{awaiter: fastAwaiter, delay: commonhtrk.TrackableCallbackTimeout / 2}
	_, unsubscribe1 := broadcaster.Subscribe(slow)
	_, unsubscribe2 := broadcaster.Subscribe(fast)

	broadcaster.BroadcastNewLongestChain(testutils.Head(1))
	slowAwaiter.AwaitOrFail(t, tests.WaitTimeout(t))
	fastAwaiter.AwaitOrFail(t, tests.WaitTimeout(t))

	require.True(t, slow.contextDone)
	require.False(t, fast.contextDone)

	unsubscribe1()
	unsubscribe2()

	err = broadcaster.Close()
	require.NoError(t, err)
}

type sleepySubscriber struct {
	awaiter     testutils.Awaiter
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
