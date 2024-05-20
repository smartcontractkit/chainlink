package headtracker_test

import (
	"context"
	"errors"
	"math/big"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	gethCommon "github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/jmoiron/sqlx"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox/mailboxtest"

	htmocks "github.com/smartcontractkit/chainlink/v2/common/headtracker/mocks"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func firstHead(t *testing.T, db *sqlx.DB) (h evmtypes.Head) {
	if err := db.Get(&h, `SELECT * FROM evm.heads ORDER BY number ASC LIMIT 1`); err != nil {
		t.Fatal(err)
	}
	return h
}

func TestHeadTracker_New(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(cltest.Head(0), nil)
	// finalized
	ethClient.On("HeadByNumber", mock.Anything, big.NewInt(0)).Return(cltest.Head(0), nil)

	orm := headtracker.NewORM(cltest.FixtureChainID, db)
	assert.Nil(t, orm.IdempotentInsertHead(testutils.Context(t), cltest.Head(1)))
	last := cltest.Head(16)
	assert.Nil(t, orm.IdempotentInsertHead(testutils.Context(t), last))
	assert.Nil(t, orm.IdempotentInsertHead(testutils.Context(t), cltest.Head(10)))

	evmcfg := cltest.NewTestChainScopedConfig(t)
	ht := createHeadTracker(t, ethClient, evmcfg.EVM(), evmcfg.EVM().HeadTracker(), orm)
	ht.Start(t)

	latest := ht.headSaver.LatestChain()
	require.NotNil(t, latest)
	assert.Equal(t, last.Number, latest.Number)
}

func TestHeadTracker_MarkFinalized_MarksAndTrimsTable(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	gCfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, _ *chainlink.Secrets) {
		c.EVM[0].HeadTracker.HistoryDepth = ptr[uint32](100)
	})
	config := evmtest.NewChainScopedConfig(t, gCfg)

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	orm := headtracker.NewORM(cltest.FixtureChainID, db)

	for idx := 0; idx < 200; idx++ {
		assert.Nil(t, orm.IdempotentInsertHead(testutils.Context(t), cltest.Head(idx)))
	}

	latest := cltest.Head(201)
	assert.Nil(t, orm.IdempotentInsertHead(testutils.Context(t), latest))

	ht := createHeadTracker(t, ethClient, config.EVM(), config.EVM().HeadTracker(), orm)
	_, err := ht.headSaver.Load(testutils.Context(t), latest.Number)
	require.NoError(t, err)
	require.NoError(t, ht.headSaver.MarkFinalized(testutils.Context(t), latest))
	assert.Equal(t, big.NewInt(201), ht.headSaver.LatestChain().ToInt())

	firstHead := firstHead(t, db)
	assert.Equal(t, big.NewInt(101), firstHead.ToInt())

	lastHead, err := orm.LatestHead(testutils.Context(t))
	require.NoError(t, err)
	assert.Equal(t, int64(201), lastHead.Number)
}

func TestHeadTracker_Get(t *testing.T) {
	t.Parallel()

	start := cltest.Head(5)

	tests := []struct {
		name    string
		initial *evmtypes.Head
		toSave  *evmtypes.Head
		want    *big.Int
	}{
		{"greater", start, cltest.Head(6), big.NewInt(6)},
		{"less than", start, cltest.Head(1), big.NewInt(5)},
		{"zero", start, cltest.Head(0), big.NewInt(5)},
		{"nil", start, nil, big.NewInt(5)},
		{"nil no initial", nil, nil, big.NewInt(0)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := pgtest.NewSqlxDB(t)
			config := cltest.NewTestChainScopedConfig(t)
			orm := headtracker.NewORM(cltest.FixtureChainID, db)

			ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
			chStarted := make(chan struct{})
			mockEth := &evmtest.MockEth{
				EthClient: ethClient,
			}
			ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).
				Maybe().
				Return(
					func(ctx context.Context, ch chan<- *evmtypes.Head) ethereum.Subscription {
						defer close(chStarted)
						return mockEth.NewSub(t)
					},
					func(ctx context.Context, ch chan<- *evmtypes.Head) error { return nil },
				)
			ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(cltest.Head(0), nil)

			fnCall := ethClient.On("HeadByNumber", mock.Anything, mock.Anything)
			fnCall.RunFn = func(args mock.Arguments) {
				num := args.Get(1).(*big.Int)
				fnCall.ReturnArguments = mock.Arguments{cltest.Head(num.Int64()), nil}
			}

			if test.initial != nil {
				assert.Nil(t, orm.IdempotentInsertHead(testutils.Context(t), test.initial))
			}

			ht := createHeadTracker(t, ethClient, config.EVM(), config.EVM().HeadTracker(), orm)
			ht.Start(t)

			if test.toSave != nil {
				err := ht.headSaver.Save(testutils.Context(t), test.toSave)
				assert.NoError(t, err)
			}

			assert.Equal(t, test.want, ht.headSaver.LatestChain().ToInt())
		})
	}
}

func TestHeadTracker_Start_NewHeads(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	config := cltest.NewTestChainScopedConfig(t)
	orm := headtracker.NewORM(cltest.FixtureChainID, db)

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	chStarted := make(chan struct{})
	mockEth := &evmtest.MockEth{EthClient: ethClient}
	sub := mockEth.NewSub(t)
	// for initial load
	ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(cltest.Head(0), nil).Once()
	ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(cltest.Head(0), nil).Once()
	// for backfill
	ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(cltest.Head(0), nil).Maybe()
	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).
		Run(func(mock.Arguments) {
			close(chStarted)
		}).
		Return(sub, nil)

	ht := createHeadTracker(t, ethClient, config.EVM(), config.EVM().HeadTracker(), orm)
	ht.Start(t)

	<-chStarted
}

func TestHeadTracker_Start(t *testing.T) {
	t.Parallel()

	const historyDepth = 100
	newHeadTracker := func(t *testing.T) *headTrackerUniverse {
		db := pgtest.NewSqlxDB(t)
		gCfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, _ *chainlink.Secrets) {
			c.EVM[0].FinalityTagEnabled = ptr[bool](true)
			c.EVM[0].HeadTracker.HistoryDepth = ptr[uint32](historyDepth)
		})
		config := evmtest.NewChainScopedConfig(t, gCfg)
		orm := headtracker.NewORM(cltest.FixtureChainID, db)
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		return createHeadTracker(t, ethClient, config.EVM(), config.EVM().HeadTracker(), orm)
	}

	t.Run("Fail start if context was canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(testutils.Context(t))
		ht := newHeadTracker(t)
		ht.ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Run(func(args mock.Arguments) {
			cancel()
		}).Return(cltest.Head(0), context.Canceled)
		err := ht.headTracker.Start(ctx)
		require.ErrorIs(t, err, context.Canceled)
	})
	t.Run("Starts even if failed to get initialHead", func(t *testing.T) {
		ht := newHeadTracker(t)
		ht.ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(cltest.Head(0), errors.New("failed to get init head"))
		ht.Start(t)
		tests.AssertLogEventually(t, ht.observer, "Error handling initial head")
	})
	t.Run("Starts even if received invalid head", func(t *testing.T) {
		ht := newHeadTracker(t)
		ht.ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(nil, nil)
		ht.Start(t)
		tests.AssertLogEventually(t, ht.observer, "Got nil initial head")
	})
	t.Run("Starts even if fails to get finalizedHead", func(t *testing.T) {
		ht := newHeadTracker(t)
		head := cltest.Head(1000)
		ht.ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(head, nil).Once()
		ht.ethClient.On("LatestFinalizedBlock", mock.Anything).Return(nil, errors.New("failed to load latest finalized")).Once()
		ht.Start(t)
		tests.AssertLogEventually(t, ht.observer, "Error handling initial head")
	})
	t.Run("Starts even if latest finalizedHead is nil", func(t *testing.T) {
		ht := newHeadTracker(t)
		head := cltest.Head(1000)
		ht.ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(head, nil).Once()
		ht.ethClient.On("LatestFinalizedBlock", mock.Anything).Return(nil, nil).Once()
		ht.Start(t)
		tests.AssertLogEventually(t, ht.observer, "Error handling initial head")
	})
	t.Run("Happy path", func(t *testing.T) {
		head := cltest.Head(1000)
		ht := newHeadTracker(t)
		ctx := testutils.Context(t)
		require.NoError(t, ht.orm.IdempotentInsertHead(ctx, cltest.Head(799)))
		ht.ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(head, nil).Once()
		finalizedHead := cltest.Head(800)
		// on start
		ht.ethClient.On("LatestFinalizedBlock", mock.Anything).Return(finalizedHead, nil).Once()
		// on backfill
		ht.ethClient.On("LatestFinalizedBlock", mock.Anything).Return(nil, errors.New("backfill call to finalized failed")).Maybe()
		ht.Start(t)
		tests.AssertLogEventually(t, ht.observer, "Loaded chain from DB")
	})
}

func TestHeadTracker_CallsHeadTrackableCallbacks(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	db := pgtest.NewSqlxDB(t)
	config := cltest.NewTestChainScopedConfig(t)
	orm := headtracker.NewORM(cltest.FixtureChainID, db)

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

	chchHeaders := make(chan evmtest.RawSub[*evmtypes.Head], 1)
	mockEth := &evmtest.MockEth{EthClient: ethClient}
	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).
		Return(
			func(ctx context.Context, ch chan<- *evmtypes.Head) ethereum.Subscription {
				sub := mockEth.NewSub(t)
				chchHeaders <- evmtest.NewRawSub(ch, sub.Err())
				return sub
			},
			func(ctx context.Context, ch chan<- *evmtypes.Head) error { return nil },
		)
	ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(cltest.Head(0), nil)
	ethClient.On("HeadByHash", mock.Anything, mock.Anything).Return(cltest.Head(0), nil).Maybe()

	checker := &cltest.MockHeadTrackable{}
	ht := createHeadTrackerWithChecker(t, ethClient, config.EVM(), config.EVM().HeadTracker(), orm, checker)

	ht.Start(t)
	assert.Equal(t, int32(0), checker.OnNewLongestChainCount())

	headers := <-chchHeaders
	headers.TrySend(&evmtypes.Head{Number: 1, Hash: utils.NewHash(), EVMChainID: ubig.New(&cltest.FixtureChainID)})
	g.Eventually(checker.OnNewLongestChainCount).Should(gomega.Equal(int32(1)))

	ht.Stop(t)
	assert.Equal(t, int32(1), checker.OnNewLongestChainCount())
}

func TestHeadTracker_ReconnectOnError(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	db := pgtest.NewSqlxDB(t)
	config := cltest.NewTestChainScopedConfig(t)
	orm := headtracker.NewORM(cltest.FixtureChainID, db)

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	mockEth := &evmtest.MockEth{EthClient: ethClient}
	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).
		Return(
			func(ctx context.Context, ch chan<- *evmtypes.Head) ethereum.Subscription { return mockEth.NewSub(t) },
			func(ctx context.Context, ch chan<- *evmtypes.Head) error { return nil },
		)
	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).Return(nil, errors.New("cannot reconnect"))
	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).
		Return(
			func(ctx context.Context, ch chan<- *evmtypes.Head) ethereum.Subscription { return mockEth.NewSub(t) },
			func(ctx context.Context, ch chan<- *evmtypes.Head) error { return nil },
		)
	ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(cltest.Head(0), nil)

	checker := &cltest.MockHeadTrackable{}
	ht := createHeadTrackerWithChecker(t, ethClient, config.EVM(), config.EVM().HeadTracker(), orm, checker)

	// connect
	ht.Start(t)
	assert.Equal(t, int32(0), checker.OnNewLongestChainCount())

	// trigger reconnect loop
	mockEth.SubsErr(errors.New("test error to force reconnect"))
	g.Eventually(checker.OnNewLongestChainCount).Should(gomega.Equal(int32(1)))
}

func TestHeadTracker_ResubscribeOnSubscriptionError(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	db := pgtest.NewSqlxDB(t)
	config := cltest.NewTestChainScopedConfig(t)
	orm := headtracker.NewORM(cltest.FixtureChainID, db)

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

	chchHeaders := make(chan evmtest.RawSub[*evmtypes.Head], 1)
	mockEth := &evmtest.MockEth{EthClient: ethClient}
	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).
		Return(
			func(ctx context.Context, ch chan<- *evmtypes.Head) ethereum.Subscription {
				sub := mockEth.NewSub(t)
				chchHeaders <- evmtest.NewRawSub(ch, sub.Err())
				return sub
			},
			func(ctx context.Context, ch chan<- *evmtypes.Head) error { return nil },
		)
	ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(cltest.Head(0), nil)
	ethClient.On("HeadByHash", mock.Anything, mock.Anything).Return(cltest.Head(0), nil).Maybe()

	checker := &cltest.MockHeadTrackable{}
	ht := createHeadTrackerWithChecker(t, ethClient, config.EVM(), config.EVM().HeadTracker(), orm, checker)

	ht.Start(t)
	assert.Equal(t, int32(0), checker.OnNewLongestChainCount())

	headers := <-chchHeaders
	go func() {
		headers.TrySend(cltest.Head(1))
	}()

	g.Eventually(func() bool {
		report := ht.headTracker.HealthReport()
		return !slices.ContainsFunc(maps.Values(report), func(e error) bool { return e != nil })
	}, 5*time.Second, testutils.TestInterval).Should(gomega.Equal(true))

	// trigger reconnect loop
	headers.CloseCh()

	// wait for full disconnect and a new subscription
	g.Eventually(checker.OnNewLongestChainCount, 5*time.Second, testutils.TestInterval).Should(gomega.Equal(int32(1)))
}

func TestHeadTracker_Start_LoadsLatestChain(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	config := cltest.NewTestChainScopedConfig(t)
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

	heads := []*evmtypes.Head{
		cltest.Head(0),
		cltest.Head(1),
		cltest.Head(2),
		cltest.Head(3),
	}
	var parentHash gethCommon.Hash
	for i := 0; i < len(heads); i++ {
		if parentHash != (gethCommon.Hash{}) {
			heads[i].ParentHash = parentHash
		}
		parentHash = heads[i].Hash
	}
	ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(heads[3], nil).Maybe()
	ethClient.On("HeadByNumber", mock.Anything, big.NewInt(0)).Return(heads[0], nil).Maybe()
	ethClient.On("HeadByHash", mock.Anything, heads[2].Hash).Return(heads[2], nil).Maybe()
	ethClient.On("HeadByHash", mock.Anything, heads[1].Hash).Return(heads[1], nil).Maybe()
	ethClient.On("HeadByHash", mock.Anything, heads[0].Hash).Return(heads[0], nil).Maybe()

	chchHeaders := make(chan evmtest.RawSub[*evmtypes.Head], 1)
	mockEth := &evmtest.MockEth{EthClient: ethClient}
	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).
		Return(
			func(ctx context.Context, ch chan<- *evmtypes.Head) ethereum.Subscription {
				sub := mockEth.NewSub(t)
				chchHeaders <- evmtest.NewRawSub(ch, sub.Err())
				return sub
			},
			func(ctx context.Context, ch chan<- *evmtypes.Head) error { return nil },
		)

	orm := headtracker.NewORM(cltest.FixtureChainID, db)
	trackable := &cltest.MockHeadTrackable{}
	ht := createHeadTrackerWithChecker(t, ethClient, config.EVM(), config.EVM().HeadTracker(), orm, trackable)

	require.NoError(t, orm.IdempotentInsertHead(testutils.Context(t), heads[2]))

	ht.Start(t)

	assert.Equal(t, int32(0), trackable.OnNewLongestChainCount())

	headers := <-chchHeaders
	go func() {
		headers.TrySend(cltest.Head(1))
	}()

	gomega.NewWithT(t).Eventually(func() bool {
		report := ht.headTracker.HealthReport()
		services.CopyHealth(report, ht.headBroadcaster.HealthReport())
		return !slices.ContainsFunc(maps.Values(report), func(e error) bool { return e != nil })
	}, 5*time.Second, testutils.TestInterval).Should(gomega.Equal(true))

	h, err := orm.LatestHead(testutils.Context(t))
	require.NoError(t, err)
	require.NotNil(t, h)
	assert.Equal(t, h.Number, int64(3))
}

func TestHeadTracker_SwitchesToLongestChainWithHeadSamplingEnabled(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)

	config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].FinalityDepth = ptr[uint32](50)
		// Need to set the buffer to something large since we inject a lot of heads at once and otherwise they will be dropped
		c.EVM[0].HeadTracker.MaxBufferSize = ptr[uint32](100)
		c.EVM[0].HeadTracker.SamplingInterval = commonconfig.MustNewDuration(2500 * time.Millisecond)
	})

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

	checker := htmocks.NewHeadTrackable[*evmtypes.Head, gethCommon.Hash](t)
	orm := headtracker.NewORM(*evmtest.MustGetDefaultChainID(t, config.EVMConfigs()), db)
	csCfg := evmtest.NewChainScopedConfig(t, config)
	ht := createHeadTrackerWithChecker(t, ethClient, csCfg.EVM(), csCfg.EVM().HeadTracker(), orm, checker)

	chchHeaders := make(chan evmtest.RawSub[*evmtypes.Head], 1)
	mockEth := &evmtest.MockEth{EthClient: ethClient}
	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).
		Return(
			func(ctx context.Context, ch chan<- *evmtypes.Head) ethereum.Subscription {
				sub := mockEth.NewSub(t)
				chchHeaders <- evmtest.NewRawSub(ch, sub.Err())
				return sub
			},
			func(ctx context.Context, ch chan<- *evmtypes.Head) error { return nil },
		)

	// ---------------------
	blocks := cltest.NewBlocks(t, 10)

	head0 := blocks.Head(0)
	// Initial query
	ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(head0, nil)
	// backfill query
	ethClient.On("HeadByNumber", mock.Anything, big.NewInt(0)).Return(head0, nil)
	ht.Start(t)

	headSeq := cltest.NewHeadBuffer(t)
	headSeq.Append(blocks.Head(0))
	headSeq.Append(blocks.Head(1))

	// Blocks 2 and 3 are out of order
	headSeq.Append(blocks.Head(3))
	headSeq.Append(blocks.Head(2))

	// Block 4 comes in
	headSeq.Append(blocks.Head(4))

	// Another block at level 4 comes in, that will be uncled
	headSeq.Append(blocks.NewHead(4))

	// Reorg happened forking from block 2
	blocksForked := blocks.ForkAt(t, 2, 5)
	headSeq.Append(blocksForked.Head(2))
	headSeq.Append(blocksForked.Head(3))
	headSeq.Append(blocksForked.Head(4))
	headSeq.Append(blocksForked.Head(5)) // Now the new chain is longer

	lastLongestChainAwaiter := cltest.NewAwaiter()

	// the callback is only called for head number 5 because of head sampling
	checker.On("OnNewLongestChain", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			h := args.Get(1).(*evmtypes.Head)

			assert.Equal(t, int64(5), h.Number)
			assert.Equal(t, blocksForked.Head(5).Hash, h.Hash)

			// This is the new longest chain, check that it came with its parents
			if !assert.NotNil(t, h.Parent) {
				return
			}
			assert.Equal(t, h.Parent.Hash, blocksForked.Head(4).Hash)
			if !assert.NotNil(t, h.Parent.Parent) {
				return
			}
			assert.Equal(t, h.Parent.Parent.Hash, blocksForked.Head(3).Hash)
			if !assert.NotNil(t, h.Parent.Parent.Parent) {
				return
			}
			assert.Equal(t, h.Parent.Parent.Parent.Hash, blocksForked.Head(2).Hash)
			if !assert.NotNil(t, h.Parent.Parent.Parent.Parent) {
				return
			}
			assert.Equal(t, h.Parent.Parent.Parent.Parent.Hash, blocksForked.Head(1).Hash)
			lastLongestChainAwaiter.ItHappened()
		}).Return().Once()

	headers := <-chchHeaders

	// This grotesque construction is the only way to do dynamic return values using
	// the mock package.  We need dynamic returns because we're simulating reorgs.
	latestHeadByHash := make(map[gethCommon.Hash]*evmtypes.Head)
	latestHeadByHashMu := new(sync.Mutex)

	fnCall := ethClient.On("HeadByHash", mock.Anything, mock.Anything).Maybe()
	fnCall.RunFn = func(args mock.Arguments) {
		latestHeadByHashMu.Lock()
		defer latestHeadByHashMu.Unlock()
		hash := args.Get(1).(gethCommon.Hash)
		head := latestHeadByHash[hash]
		fnCall.ReturnArguments = mock.Arguments{head, nil}
	}

	for _, h := range headSeq.Heads {
		latestHeadByHashMu.Lock()
		latestHeadByHash[h.Hash] = h
		latestHeadByHashMu.Unlock()
		headers.TrySend(h)
	}

	// default 10s may not be sufficient, so using testutils.WaitTimeout(t)
	lastLongestChainAwaiter.AwaitOrFail(t, testutils.WaitTimeout(t))
	ht.Stop(t)
	assert.Equal(t, int64(5), ht.headSaver.LatestChain().Number)

	for _, h := range headSeq.Heads {
		c := ht.headSaver.Chain(h.Hash)
		require.NotNil(t, c)
		assert.Equal(t, c.ParentHash, h.ParentHash)
		assert.Equal(t, c.Timestamp.Unix(), h.Timestamp.UTC().Unix())
		assert.Equal(t, c.Number, h.Number)
	}
}

func TestHeadTracker_SwitchesToLongestChainWithHeadSamplingDisabled(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)

	config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].FinalityDepth = ptr[uint32](50)
		// Need to set the buffer to something large since we inject a lot of heads at once and otherwise they will be dropped
		c.EVM[0].HeadTracker.MaxBufferSize = ptr[uint32](100)
		c.EVM[0].HeadTracker.SamplingInterval = commonconfig.MustNewDuration(0)
	})

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

	checker := htmocks.NewHeadTrackable[*evmtypes.Head, gethCommon.Hash](t)
	orm := headtracker.NewORM(cltest.FixtureChainID, db)
	evmcfg := evmtest.NewChainScopedConfig(t, config)
	ht := createHeadTrackerWithChecker(t, ethClient, evmcfg.EVM(), evmcfg.EVM().HeadTracker(), orm, checker)

	chchHeaders := make(chan evmtest.RawSub[*evmtypes.Head], 1)
	mockEth := &evmtest.MockEth{EthClient: ethClient}
	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).
		Return(
			func(ctx context.Context, ch chan<- *evmtypes.Head) ethereum.Subscription {
				sub := mockEth.NewSub(t)
				chchHeaders <- evmtest.NewRawSub(ch, sub.Err())
				return sub
			},
			func(ctx context.Context, ch chan<- *evmtypes.Head) error { return nil },
		)

	// ---------------------
	blocks := cltest.NewBlocks(t, 10)

	head0 := blocks.Head(0) // evmtypes.Head{Number: 0, Hash: utils.NewHash(), ParentHash: utils.NewHash(), Timestamp: time.Unix(0, 0)}
	// Initial query
	ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(head0, nil)
	// backfill
	ethClient.On("HeadByNumber", mock.Anything, big.NewInt(0)).Return(head0, nil)

	headSeq := cltest.NewHeadBuffer(t)
	headSeq.Append(blocks.Head(0))
	headSeq.Append(blocks.Head(1))

	// Blocks 2 and 3 are out of order
	headSeq.Append(blocks.Head(3))
	headSeq.Append(blocks.Head(2))

	// Block 4 comes in
	headSeq.Append(blocks.Head(4))

	// Another block at level 4 comes in, that will be uncled
	headSeq.Append(blocks.NewHead(4))

	// Reorg happened forking from block 2
	blocksForked := blocks.ForkAt(t, 2, 5)
	headSeq.Append(blocksForked.Head(2))
	headSeq.Append(blocksForked.Head(3))
	headSeq.Append(blocksForked.Head(4))
	headSeq.Append(blocksForked.Head(5)) // Now the new chain is longer

	lastLongestChainAwaiter := cltest.NewAwaiter()

	checker.On("OnNewLongestChain", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			h := args.Get(1).(*evmtypes.Head)
			require.Equal(t, int64(0), h.Number)
			require.Equal(t, blocks.Head(0).Hash, h.Hash)
		}).Return().Once()

	checker.On("OnNewLongestChain", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			h := args.Get(1).(*evmtypes.Head)
			require.Equal(t, int64(1), h.Number)
			require.Equal(t, blocks.Head(1).Hash, h.Hash)
		}).Return().Once()

	checker.On("OnNewLongestChain", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			h := args.Get(1).(*evmtypes.Head)
			require.Equal(t, int64(3), h.Number)
			require.Equal(t, blocks.Head(3).Hash, h.Hash)
		}).Return().Once()

	checker.On("OnNewLongestChain", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			h := args.Get(1).(*evmtypes.Head)
			require.Equal(t, int64(4), h.Number)
			require.Equal(t, blocks.Head(4).Hash, h.Hash)

			// Check that the block came with its parents
			require.NotNil(t, h.Parent)
			require.Equal(t, h.Parent.Hash, blocks.Head(3).Hash)
			require.NotNil(t, h.Parent.Parent.Hash)
			require.Equal(t, h.Parent.Parent.Hash, blocks.Head(2).Hash)
			require.NotNil(t, h.Parent.Parent.Parent)
			require.Equal(t, h.Parent.Parent.Parent.Hash, blocks.Head(1).Hash)
		}).Return().Once()

	checker.On("OnNewLongestChain", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			h := args.Get(1).(*evmtypes.Head)

			require.Equal(t, int64(5), h.Number)
			require.Equal(t, blocksForked.Head(5).Hash, h.Hash)

			// This is the new longest chain, check that it came with its parents
			require.NotNil(t, h.Parent)
			require.Equal(t, h.Parent.Hash, blocksForked.Head(4).Hash)
			require.NotNil(t, h.Parent.Parent)
			require.Equal(t, h.Parent.Parent.Hash, blocksForked.Head(3).Hash)
			require.NotNil(t, h.Parent.Parent.Parent)
			require.Equal(t, h.Parent.Parent.Parent.Hash, blocksForked.Head(2).Hash)
			require.NotNil(t, h.Parent.Parent.Parent.Parent)
			require.Equal(t, h.Parent.Parent.Parent.Parent.Hash, blocksForked.Head(1).Hash)
			lastLongestChainAwaiter.ItHappened()
		}).Return().Once()

	ht.Start(t)

	headers := <-chchHeaders

	// This grotesque construction is the only way to do dynamic return values using
	// the mock package.  We need dynamic returns because we're simulating reorgs.
	latestHeadByHash := make(map[gethCommon.Hash]*evmtypes.Head)
	latestHeadByHashMu := new(sync.Mutex)

	fnCall := ethClient.On("HeadByHash", mock.Anything, mock.Anything).Maybe()
	fnCall.RunFn = func(args mock.Arguments) {
		latestHeadByHashMu.Lock()
		defer latestHeadByHashMu.Unlock()
		hash := args.Get(1).(gethCommon.Hash)
		head := latestHeadByHash[hash]
		fnCall.ReturnArguments = mock.Arguments{head, nil}
	}

	for _, h := range headSeq.Heads {
		latestHeadByHashMu.Lock()
		latestHeadByHash[h.Hash] = h
		latestHeadByHashMu.Unlock()
		headers.TrySend(h)
		time.Sleep(testutils.TestInterval)
	}

	// default 10s may not be sufficient, so using testutils.WaitTimeout(t)
	lastLongestChainAwaiter.AwaitOrFail(t, testutils.WaitTimeout(t))
	ht.Stop(t)
	assert.Equal(t, int64(5), ht.headSaver.LatestChain().Number)

	for _, h := range headSeq.Heads {
		c := ht.headSaver.Chain(h.Hash)
		require.NotNil(t, c)
		assert.Equal(t, c.ParentHash, h.ParentHash)
		assert.Equal(t, c.Timestamp.Unix(), h.Timestamp.UTC().Unix())
		assert.Equal(t, c.Number, h.Number)
	}
}

func TestHeadTracker_Backfill(t *testing.T) {
	t.Parallel()

	// Heads are arranged as follows:
	// headN indicates an unpersisted ethereum header
	// hN indicates a persisted head record
	//
	// (1)->(H0)
	//
	//       (14Orphaned)-+
	//                    +->(13)->(12)->(11)->(H10)->(9)->(H8)
	// (15)->(14)---------+

	now := uint64(time.Now().UTC().Unix())

	gethHead0 := &gethTypes.Header{
		Number:     big.NewInt(0),
		ParentHash: gethCommon.BigToHash(big.NewInt(0)),
		Time:       now,
	}
	head0 := evmtypes.NewHead(gethHead0.Number, utils.NewHash(), gethHead0.ParentHash, gethHead0.Time, ubig.New(&cltest.FixtureChainID))

	h1 := *cltest.Head(1)
	h1.ParentHash = head0.Hash

	gethHead8 := &gethTypes.Header{
		Number:     big.NewInt(8),
		ParentHash: utils.NewHash(),
		Time:       now,
	}
	head8 := evmtypes.NewHead(gethHead8.Number, utils.NewHash(), gethHead8.ParentHash, gethHead8.Time, ubig.New(&cltest.FixtureChainID))

	h9 := *cltest.Head(9)
	h9.ParentHash = head8.Hash

	gethHead10 := &gethTypes.Header{
		Number:     big.NewInt(10),
		ParentHash: h9.Hash,
		Time:       now,
	}
	head10 := evmtypes.NewHead(gethHead10.Number, utils.NewHash(), gethHead10.ParentHash, gethHead10.Time, ubig.New(&cltest.FixtureChainID))

	h11 := *cltest.Head(11)
	h11.ParentHash = head10.Hash

	h12 := *cltest.Head(12)
	h12.ParentHash = h11.Hash

	h13 := *cltest.Head(13)
	h13.ParentHash = h12.Hash

	h14Orphaned := *cltest.Head(14)
	h14Orphaned.ParentHash = h13.Hash

	h14 := *cltest.Head(14)
	h14.ParentHash = h13.Hash

	h15 := *cltest.Head(15)
	h15.ParentHash = h14.Hash

	heads := []evmtypes.Head{
		h9,
		h11,
		h12,
		h13,
		h14Orphaned,
		h14,
		h15,
	}

	ctx := testutils.Context(t)

	type opts struct {
		Heads []evmtypes.Head
	}
	newHeadTrackerUniverse := func(t *testing.T, opts opts) *headTrackerUniverse {
		cfg := configtest.NewGeneralConfig(t, nil)

		evmcfg := evmtest.NewChainScopedConfig(t, cfg)
		db := pgtest.NewSqlxDB(t)
		orm := headtracker.NewORM(cltest.FixtureChainID, db)
		for i := range opts.Heads {
			require.NoError(t, orm.IdempotentInsertHead(testutils.Context(t), &opts.Heads[i]))
		}
		ethClient := evmtest.NewEthClientMock(t)
		ethClient.On("ConfiguredChainID", mock.Anything).Return(evmtest.MustGetDefaultChainID(t, cfg.EVMConfigs()), nil)
		ht := createHeadTracker(t, ethClient, evmcfg.EVM(), evmcfg.EVM().HeadTracker(), orm)
		_, err := ht.headSaver.Load(testutils.Context(t), 0)
		require.NoError(t, err)
		return ht
	}

	t.Run("returns error if latestFinalized is not valid", func(t *testing.T) {
		htu := newHeadTrackerUniverse(t, opts{})

		err := htu.headTracker.Backfill(ctx, &h12, nil)
		require.EqualError(t, err, "can not perform backfill without a valid latestFinalized head")
	})
	t.Run("Returns error if finalized head is ahead of canonical", func(t *testing.T) {
		htu := newHeadTrackerUniverse(t, opts{})

		err := htu.headTracker.Backfill(ctx, &h12, &h14Orphaned)
		require.EqualError(t, err, "invariant violation: expected head of canonical chain to be ahead of the latestFinalized")
	})
	t.Run("Returns error if finalizedHead is not present in the canonical chain", func(t *testing.T) {
		htu := newHeadTrackerUniverse(t, opts{Heads: heads})

		err := htu.headTracker.Backfill(ctx, &h15, &h14Orphaned)
		require.EqualError(t, err, "expected finalized block to be present in canonical chain")
	})
	t.Run("Marks all blocks in chain that are older than finalized", func(t *testing.T) {
		htu := newHeadTrackerUniverse(t, opts{Heads: heads})

		assertFinalized := func(expectedFinalized bool, msg string, heads ...evmtypes.Head) {
			for _, h := range heads {
				storedHead := htu.headSaver.Chain(h.Hash)
				assert.Equal(t, expectedFinalized, storedHead != nil && storedHead.IsFinalized, msg, "block_number", h.Number)
			}
		}

		err := htu.headTracker.Backfill(ctx, &h15, &h14)
		require.NoError(t, err)
		assertFinalized(true, "expected heads to be marked as finalized after backfill", h14, h13, h12, h11)
		assertFinalized(false, "expected heads to remain unfinalized", h15, head10)
	})

	t.Run("fetches a missing head", func(t *testing.T) {
		htu := newHeadTrackerUniverse(t, opts{Heads: heads})
		htu.ethClient.On("HeadByHash", mock.Anything, head10.Hash).
			Return(&head10, nil)

		err := htu.headTracker.Backfill(ctx, &h12, &h9)
		require.NoError(t, err)

		h := htu.headSaver.Chain(h12.Hash)

		assert.Equal(t, int64(12), h.Number)
		require.NotNil(t, h.Parent)
		assert.Equal(t, int64(11), h.Parent.Number)
		require.NotNil(t, h.Parent.Parent)
		assert.Equal(t, int64(10), h.Parent.Parent.Number)
		require.NotNil(t, h.Parent.Parent.Parent)
		assert.Equal(t, int64(9), h.Parent.Parent.Parent.Number)

		writtenHead, err := htu.orm.HeadByHash(testutils.Context(t), head10.Hash)
		require.NoError(t, err)
		assert.Equal(t, int64(10), writtenHead.Number)
	})

	t.Run("fetches only heads that are missing", func(t *testing.T) {
		htu := newHeadTrackerUniverse(t, opts{Heads: heads})

		htu.ethClient.On("HeadByHash", mock.Anything, head10.Hash).
			Return(&head10, nil)
		htu.ethClient.On("HeadByHash", mock.Anything, head8.Hash).
			Return(&head8, nil)

		err := htu.headTracker.Backfill(ctx, &h15, &head8)
		require.NoError(t, err)

		h := htu.headSaver.Chain(h15.Hash)

		require.Equal(t, uint32(8), h.ChainLength())
		earliestInChain := h.EarliestInChain()
		assert.Equal(t, head8.Number, earliestInChain.BlockNumber())
		assert.Equal(t, head8.Hash, earliestInChain.BlockHash())
	})

	t.Run("abandons backfill and returns error if the eth node returns not found", func(t *testing.T) {
		htu := newHeadTrackerUniverse(t, opts{Heads: heads})
		htu.ethClient.On("HeadByHash", mock.Anything, head10.Hash).
			Return(&head10, nil).
			Once()
		htu.ethClient.On("HeadByHash", mock.Anything, head8.Hash).
			Return(nil, ethereum.NotFound).
			Once()

		err := htu.headTracker.Backfill(ctx, &h12, &head8)
		require.Error(t, err)
		require.EqualError(t, err, "fetchAndSaveHead failed: not found")

		h := htu.headSaver.Chain(h12.Hash)

		// Should contain 12, 11, 10, 9
		assert.Equal(t, 4, int(h.ChainLength()))
		assert.Equal(t, int64(9), h.EarliestInChain().BlockNumber())
	})

	t.Run("abandons backfill and returns error if the context time budget is exceeded", func(t *testing.T) {
		htu := newHeadTrackerUniverse(t, opts{Heads: heads})
		htu.ethClient.On("HeadByHash", mock.Anything, head10.Hash).
			Return(&head10, nil)
		lctx, cancel := context.WithCancel(ctx)
		htu.ethClient.On("HeadByHash", mock.Anything, head8.Hash).
			Return(nil, context.DeadlineExceeded).Run(func(args mock.Arguments) {
			cancel()
		})

		err := htu.headTracker.Backfill(lctx, &h12, &head8)
		require.Error(t, err)
		require.EqualError(t, err, "fetchAndSaveHead failed: context canceled")

		h := htu.headSaver.Chain(h12.Hash)

		// Should contain 12, 11, 10, 9
		assert.Equal(t, 4, int(h.ChainLength()))
		assert.Equal(t, int64(9), h.EarliestInChain().BlockNumber())
	})

	t.Run("abandons backfill and returns error when fetching a block by hash fails, indicating a reorg", func(t *testing.T) {
		htu := newHeadTrackerUniverse(t, opts{})
		htu.ethClient.On("HeadByHash", mock.Anything, h14.Hash).Return(&h14, nil).Once()
		htu.ethClient.On("HeadByHash", mock.Anything, h13.Hash).Return(&h13, nil).Once()
		htu.ethClient.On("HeadByHash", mock.Anything, h12.Hash).Return(nil, errors.New("not found")).Once()

		err := htu.headTracker.Backfill(ctx, &h15, &h11)

		require.Error(t, err)
		require.EqualError(t, err, "fetchAndSaveHead failed: not found")

		h := htu.headSaver.Chain(h14.Hash)

		// Should contain 14, 13 (15 was never added). When trying to get the parent of h13 by hash, a reorg happened and backfill exited.
		assert.Equal(t, 2, int(h.ChainLength()))
		assert.Equal(t, int64(13), h.EarliestInChain().BlockNumber())
	})
	t.Run("marks head as finalized, if latestHead = finalizedHead (0 finality depth)", func(t *testing.T) {
		htu := newHeadTrackerUniverse(t, opts{Heads: []evmtypes.Head{h15}})
		finalizedH15 := h15 // copy h15 to have different addresses
		err := htu.headTracker.Backfill(ctx, &h15, &finalizedH15)
		require.NoError(t, err)

		h := htu.headSaver.LatestChain()

		// Should contain 14, 13 (15 was never added). When trying to get the parent of h13 by hash, a reorg happened and backfill exited.
		assert.Equal(t, 1, int(h.ChainLength()))
		assert.True(t, h.IsFinalized)
		assert.Equal(t, h15.BlockNumber(), h.BlockNumber())
		assert.Equal(t, h15.Hash, h.Hash)
	})
}

// BenchmarkHeadTracker_Backfill - benchmarks HeadTracker's Backfill with focus on efficiency after initial
// backfill on start up
func BenchmarkHeadTracker_Backfill(b *testing.B) {
	cfg := configtest.NewGeneralConfig(b, nil)

	evmcfg := evmtest.NewChainScopedConfig(b, cfg)
	db := pgtest.NewSqlxDB(b)
	chainID := big.NewInt(evmclient.NullClientChainID)
	orm := headtracker.NewORM(*chainID, db)
	ethClient := evmclimocks.NewClient(b)
	ethClient.On("ConfiguredChainID").Return(chainID)
	ht := createHeadTracker(b, ethClient, evmcfg.EVM(), evmcfg.EVM().HeadTracker(), orm)
	ctx := tests.Context(b)
	makeHash := func(n int64) gethCommon.Hash {
		return gethCommon.BigToHash(big.NewInt(n))
	}
	const finalityDepth = 12000 // observed value on Arbitrum
	makeBlock := func(n int64) *evmtypes.Head {
		return &evmtypes.Head{Number: n, Hash: makeHash(n), ParentHash: makeHash(n - 1)}
	}
	latest := makeBlock(finalityDepth)
	finalized := makeBlock(1)
	ethClient.On("HeadByHash", mock.Anything, mock.Anything).Return(func(_ context.Context, hash gethCommon.Hash) (*evmtypes.Head, error) {
		number := hash.Big().Int64()
		return makeBlock(number), nil
	})
	// run initial backfill to populate the database
	err := ht.headTracker.Backfill(ctx, latest, finalized)
	require.NoError(b, err)
	b.ResetTimer()
	// focus benchmark on processing of a new latest block
	for i := 0; i < b.N; i++ {
		latest = makeBlock(int64(finalityDepth + i))
		finalized = makeBlock(int64(i + 1))
		err := ht.headTracker.Backfill(ctx, latest, finalized)
		require.NoError(b, err)
	}
}

func createHeadTracker(t testing.TB, ethClient *evmclimocks.Client, config headtracker.Config, htConfig headtracker.HeadTrackerConfig, orm headtracker.ORM) *headTrackerUniverse {
	lggr, ob := logger.TestObserved(t, zap.DebugLevel)
	hb := headtracker.NewHeadBroadcaster(lggr)
	hs := headtracker.NewHeadSaver(lggr, orm, config, htConfig)
	mailMon := mailboxtest.NewMonitor(t)
	return &headTrackerUniverse{
		mu:              new(sync.Mutex),
		headTracker:     headtracker.NewHeadTracker(lggr, ethClient, config, htConfig, hb, hs, mailMon),
		headBroadcaster: hb,
		headSaver:       hs,
		mailMon:         mailMon,
		observer:        ob,
		orm:             orm,
		ethClient:       ethClient,
	}
}

func createHeadTrackerWithChecker(t *testing.T, ethClient *evmclimocks.Client, config headtracker.Config, htConfig headtracker.HeadTrackerConfig, orm headtracker.ORM, checker httypes.HeadTrackable) *headTrackerUniverse {
	lggr, ob := logger.TestObserved(t, zap.DebugLevel)
	hb := headtracker.NewHeadBroadcaster(lggr)
	hs := headtracker.NewHeadSaver(lggr, orm, config, htConfig)
	hb.Subscribe(checker)
	mailMon := mailboxtest.NewMonitor(t)
	ht := headtracker.NewHeadTracker(lggr, ethClient, config, htConfig, hb, hs, mailMon)
	return &headTrackerUniverse{
		mu:              new(sync.Mutex),
		headTracker:     ht,
		headBroadcaster: hb,
		headSaver:       hs,
		mailMon:         mailMon,
		observer:        ob,
		orm:             orm,
		ethClient:       ethClient,
	}
}

type headTrackerUniverse struct {
	mu              *sync.Mutex
	stopped         bool
	headTracker     httypes.HeadTracker
	headBroadcaster httypes.HeadBroadcaster
	headSaver       httypes.HeadSaver
	mailMon         *mailbox.Monitor
	observer        *observer.ObservedLogs
	orm             headtracker.ORM
	ethClient       *evmclimocks.Client
}

func (u *headTrackerUniverse) Backfill(ctx context.Context, head, finalizedHead *evmtypes.Head) error {
	return u.headTracker.Backfill(ctx, head, finalizedHead)
}

func (u *headTrackerUniverse) Start(t *testing.T) {
	u.mu.Lock()
	defer u.mu.Unlock()
	ctx := testutils.Context(t)
	require.NoError(t, u.headBroadcaster.Start(ctx))
	require.NoError(t, u.headTracker.Start(ctx))
	require.NoError(t, u.mailMon.Start(ctx))

	g := gomega.NewWithT(t)
	g.Eventually(func() bool {
		report := u.headBroadcaster.HealthReport()
		return !slices.ContainsFunc(maps.Values(report), func(e error) bool { return e != nil })
	}, 5*time.Second, testutils.TestInterval).Should(gomega.Equal(true))

	t.Cleanup(func() {
		u.Stop(t)
	})
}

func (u *headTrackerUniverse) Stop(t *testing.T) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.stopped {
		return
	}
	u.stopped = true
	require.NoError(t, u.headBroadcaster.Close())
	require.NoError(t, u.headTracker.Close())
	require.NoError(t, u.mailMon.Close())
}

func ptr[T any](t T) *T { return &t }
