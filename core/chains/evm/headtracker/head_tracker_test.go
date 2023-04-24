package headtracker_test

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"testing"
	"time"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"

	"github.com/ethereum/go-ethereum"
	gethCommon "github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	txmmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/types/mocks"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func firstHead(t *testing.T, db *sqlx.DB) (h evmtypes.Head) {
	if err := db.Get(&h, `SELECT * FROM evm_heads ORDER BY number ASC LIMIT 1`); err != nil {
		t.Fatal(err)
	}
	return h
}

func TestHeadTracker_New(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	logger := logger.TestLogger(t)
	config := configtest.NewGeneralConfig(t, nil)
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(cltest.Head(0), nil)

	orm := headtracker.NewORM(db, logger, config, cltest.FixtureChainID)
	assert.Nil(t, orm.IdempotentInsertHead(testutils.Context(t), cltest.Head(1)))
	last := cltest.Head(16)
	assert.Nil(t, orm.IdempotentInsertHead(testutils.Context(t), last))
	assert.Nil(t, orm.IdempotentInsertHead(testutils.Context(t), cltest.Head(10)))

	evmcfg := cltest.NewTestChainScopedConfig(t)
	ht := createHeadTracker(t, ethClient, evmcfg, orm)
	ht.Start(t)

	latest := ht.headSaver.LatestChain()
	require.NotNil(t, latest)
	assert.Equal(t, last.Number, latest.Number)
}

func TestHeadTracker_Save_InsertsAndTrimsTable(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	logger := logger.TestLogger(t)
	config := cltest.NewTestChainScopedConfig(t)

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	orm := headtracker.NewORM(db, logger, config, cltest.FixtureChainID)

	for idx := 0; idx < 200; idx++ {
		assert.Nil(t, orm.IdempotentInsertHead(testutils.Context(t), cltest.Head(idx)))
	}

	ht := createHeadTracker(t, ethClient, config, orm)

	h := cltest.Head(200)
	require.NoError(t, ht.headSaver.Save(testutils.Context(t), h))
	assert.Equal(t, big.NewInt(200), ht.headSaver.LatestChain().ToInt())

	firstHead := firstHead(t, db)
	assert.Equal(t, big.NewInt(101), firstHead.ToInt())

	lastHead, err := orm.LatestHead(testutils.Context(t))
	require.NoError(t, err)
	assert.Equal(t, int64(200), lastHead.Number)
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
			logger := logger.TestLogger(t)
			config := cltest.NewTestChainScopedConfig(t)
			orm := headtracker.NewORM(db, logger, config, cltest.FixtureChainID)

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

			ht := createHeadTracker(t, ethClient, config, orm)
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
	logger := logger.TestLogger(t)
	config := cltest.NewTestChainScopedConfig(t)
	orm := headtracker.NewORM(db, logger, config, cltest.FixtureChainID)

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	chStarted := make(chan struct{})
	mockEth := &evmtest.MockEth{EthClient: ethClient}
	sub := mockEth.NewSub(t)
	ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(cltest.Head(0), nil)
	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).
		Run(func(mock.Arguments) {
			close(chStarted)
		}).
		Return(sub, nil)

	ht := createHeadTracker(t, ethClient, config, orm)
	ht.Start(t)

	<-chStarted
}

func TestHeadTracker_Start_CancelContext(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	logger := logger.TestLogger(t)
	config := cltest.NewTestChainScopedConfig(t)
	orm := headtracker.NewORM(db, logger, config, cltest.FixtureChainID)
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	chStarted := make(chan struct{})
	ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Run(func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)
		select {
		case <-ctx.Done():
			return
		case <-time.After(10 * time.Second):
			assert.FailNow(t, "context was not cancelled within 10s")
		}
	}).Return(cltest.Head(0), nil)
	mockEth := &evmtest.MockEth{EthClient: ethClient}
	sub := mockEth.NewSub(t)
	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).
		Run(func(mock.Arguments) {
			close(chStarted)
		}).
		Return(sub, nil).
		Maybe()

	ht := createHeadTracker(t, ethClient, config, orm)

	ctx, cancel := context.WithCancel(testutils.Context(t))
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()
	err := ht.headTracker.Start(ctx)
	require.NoError(t, err)
	require.NoError(t, ht.headTracker.Close())
}

func TestHeadTracker_CallsHeadTrackableCallbacks(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	db := pgtest.NewSqlxDB(t)
	logger := logger.TestLogger(t)
	config := cltest.NewTestChainScopedConfig(t)
	orm := headtracker.NewORM(db, logger, config, cltest.FixtureChainID)

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

	checker := &cltest.MockHeadTrackable{}
	ht := createHeadTrackerWithChecker(t, ethClient, config, orm, checker)

	ht.Start(t)
	assert.Equal(t, int32(0), checker.OnNewLongestChainCount())

	headers := <-chchHeaders
	headers.TrySend(&evmtypes.Head{Number: 1, Hash: utils.NewHash(), EVMChainID: utils.NewBig(&cltest.FixtureChainID)})
	g.Eventually(checker.OnNewLongestChainCount).Should(gomega.Equal(int32(1)))

	ht.Stop(t)
	assert.Equal(t, int32(1), checker.OnNewLongestChainCount())
}

func TestHeadTracker_ReconnectOnError(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	db := pgtest.NewSqlxDB(t)
	logger := logger.TestLogger(t)
	config := cltest.NewTestChainScopedConfig(t)
	orm := headtracker.NewORM(db, logger, config, cltest.FixtureChainID)

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
	ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(cltest.Head(0), nil)

	checker := &cltest.MockHeadTrackable{}
	ht := createHeadTrackerWithChecker(t, ethClient, config, orm, checker)

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
	logger := logger.TestLogger(t)
	config := cltest.NewTestChainScopedConfig(t)
	orm := headtracker.NewORM(db, logger, config, cltest.FixtureChainID)

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

	checker := &cltest.MockHeadTrackable{}
	ht := createHeadTrackerWithChecker(t, ethClient, config, orm, checker)

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
	logger := logger.TestLogger(t)
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
	ethClient.On("HeadByNumber", mock.Anything, big.NewInt(2)).Return(heads[2], nil).Maybe()
	ethClient.On("HeadByNumber", mock.Anything, big.NewInt(1)).Return(heads[1], nil).Maybe()
	ethClient.On("HeadByNumber", mock.Anything, big.NewInt(0)).Return(heads[0], nil).Maybe()

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

	orm := headtracker.NewORM(db, logger, config, cltest.FixtureChainID)
	trackable := &cltest.MockHeadTrackable{}
	ht := createHeadTrackerWithChecker(t, ethClient, config, orm, trackable)

	require.NoError(t, orm.IdempotentInsertHead(testutils.Context(t), heads[2]))

	ht.Start(t)

	assert.Equal(t, int32(0), trackable.OnNewLongestChainCount())

	headers := <-chchHeaders
	go func() {
		headers.TrySend(cltest.Head(1))
	}()

	gomega.NewWithT(t).Eventually(func() bool {
		report := ht.headTracker.HealthReport()
		maps.Copy(report, ht.headBroadcaster.HealthReport())
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
	logger := logger.TestLogger(t)

	config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].FinalityDepth = ptr[uint32](50)
		// Need to set the buffer to something large since we inject a lot of heads at once and otherwise they will be dropped
		c.EVM[0].HeadTracker.MaxBufferSize = ptr[uint32](100)
		c.EVM[0].HeadTracker.SamplingInterval = models.MustNewDuration(2500 * time.Millisecond)
	})

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

	checker := txmmocks.NewHeadTrackable[*evmtypes.Head](t)
	orm := headtracker.NewORM(db, logger, config, *config.DefaultChainID())
	ht := createHeadTrackerWithChecker(t, ethClient, evmtest.NewChainScopedConfig(t, config), orm, checker)

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
	latestHeadByNumber := make(map[int64]*evmtypes.Head)
	latestHeadByNumberMu := new(sync.Mutex)

	fnCall := ethClient.On("HeadByNumber", mock.Anything, mock.Anything)
	fnCall.RunFn = func(args mock.Arguments) {
		latestHeadByNumberMu.Lock()
		defer latestHeadByNumberMu.Unlock()
		num := args.Get(1).(*big.Int)
		head, exists := latestHeadByNumber[num.Int64()]
		if !exists {
			head = cltest.Head(num.Int64())
			latestHeadByNumber[num.Int64()] = head
		}
		fnCall.ReturnArguments = mock.Arguments{head, nil}
	}

	for _, h := range headSeq.Heads {
		latestHeadByNumberMu.Lock()
		latestHeadByNumber[h.Number] = h
		latestHeadByNumberMu.Unlock()
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
	logger := logger.TestLogger(t)

	config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].FinalityDepth = ptr[uint32](50)
		// Need to set the buffer to something large since we inject a lot of heads at once and otherwise they will be dropped
		c.EVM[0].HeadTracker.MaxBufferSize = ptr[uint32](100)
		c.EVM[0].HeadTracker.SamplingInterval = models.MustNewDuration(0)
	})

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

	checker := txmmocks.NewHeadTrackable[*evmtypes.Head](t)
	orm := headtracker.NewORM(db, logger, config, cltest.FixtureChainID)
	evmcfg := evmtest.NewChainScopedConfig(t, config)
	ht := createHeadTrackerWithChecker(t, ethClient, evmcfg, orm, checker)

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
	latestHeadByNumber := make(map[int64]*evmtypes.Head)
	latestHeadByNumberMu := new(sync.Mutex)

	fnCall := ethClient.On("HeadByNumber", mock.Anything, mock.Anything)
	fnCall.RunFn = func(args mock.Arguments) {
		latestHeadByNumberMu.Lock()
		defer latestHeadByNumberMu.Unlock()
		num := args.Get(1).(*big.Int)
		head, exists := latestHeadByNumber[num.Int64()]
		if !exists {
			head = cltest.Head(num.Int64())
			latestHeadByNumber[num.Int64()] = head
		}
		fnCall.ReturnArguments = mock.Arguments{head, nil}
	}

	for _, h := range headSeq.Heads {
		latestHeadByNumberMu.Lock()
		latestHeadByNumber[h.Number] = h
		latestHeadByNumberMu.Unlock()
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
	head0 := evmtypes.NewHead(gethHead0.Number, utils.NewHash(), gethHead0.ParentHash, gethHead0.Time, utils.NewBig(&cltest.FixtureChainID))

	h1 := *cltest.Head(1)
	h1.ParentHash = head0.Hash

	gethHead8 := &gethTypes.Header{
		Number:     big.NewInt(8),
		ParentHash: utils.NewHash(),
		Time:       now,
	}
	head8 := evmtypes.NewHead(gethHead8.Number, utils.NewHash(), gethHead8.ParentHash, gethHead8.Time, utils.NewBig(&cltest.FixtureChainID))

	h9 := *cltest.Head(9)
	h9.ParentHash = head8.Hash

	gethHead10 := &gethTypes.Header{
		Number:     big.NewInt(10),
		ParentHash: h9.Hash,
		Time:       now,
	}
	head10 := evmtypes.NewHead(gethHead10.Number, utils.NewHash(), gethHead10.ParentHash, gethHead10.Time, utils.NewBig(&cltest.FixtureChainID))

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

	t.Run("does nothing if all the heads are in database", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		cfg := configtest.NewGeneralConfig(t, nil)
		logger := logger.TestLogger(t)
		orm := headtracker.NewORM(db, logger, cfg, cltest.FixtureChainID)
		for i := range heads {
			require.NoError(t, orm.IdempotentInsertHead(testutils.Context(t), &heads[i]))
		}

		ethClient := evmtest.NewEthClientMock(t)
		ethClient.On("ConfiguredChainID", mock.Anything).Return(cfg.DefaultChainID(), nil)
		ht := createHeadTrackerWithNeverSleeper(t, ethClient, cfg, orm)

		err := ht.Backfill(ctx, &h12, 2)
		require.NoError(t, err)
	})

	t.Run("fetches a missing head", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		cfg := configtest.NewGeneralConfig(t, nil)
		logger := logger.TestLogger(t)
		orm := headtracker.NewORM(db, logger, cfg, cltest.FixtureChainID)
		for i := range heads {
			require.NoError(t, orm.IdempotentInsertHead(testutils.Context(t), &heads[i]))
		}

		ethClient := evmtest.NewEthClientMock(t)
		ethClient.On("ConfiguredChainID", mock.Anything).Return(cfg.DefaultChainID(), nil)
		ethClient.On("HeadByNumber", mock.Anything, big.NewInt(10)).
			Return(&head10, nil)

		ht := createHeadTrackerWithNeverSleeper(t, ethClient, cfg, orm)

		var depth uint = 3

		err := ht.Backfill(ctx, &h12, depth)
		require.NoError(t, err)

		h := ht.headSaver.Chain(h12.Hash)

		assert.Equal(t, int64(12), h.Number)
		require.NotNil(t, h.Parent)
		assert.Equal(t, int64(11), h.Parent.Number)
		require.NotNil(t, h.Parent)
		assert.Equal(t, int64(10), h.Parent.Parent.Number)
		require.NotNil(t, h.Parent.Parent.Parent)
		assert.Equal(t, int64(9), h.Parent.Parent.Parent.Number)

		writtenHead, err := orm.HeadByHash(testutils.Context(t), head10.Hash)
		require.NoError(t, err)
		assert.Equal(t, int64(10), writtenHead.Number)
	})

	t.Run("fetches only heads that are missing", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		cfg := configtest.NewGeneralConfig(t, nil)
		logger := logger.TestLogger(t)
		orm := headtracker.NewORM(db, logger, cfg, cltest.FixtureChainID)
		for i := range heads {
			require.NoError(t, orm.IdempotentInsertHead(testutils.Context(t), &heads[i]))
		}

		ethClient := evmtest.NewEthClientMock(t)
		ethClient.On("ConfiguredChainID", mock.Anything).Return(cfg.DefaultChainID(), nil)

		ht := createHeadTrackerWithNeverSleeper(t, ethClient, cfg, orm)

		ethClient.On("HeadByNumber", mock.Anything, big.NewInt(10)).
			Return(&head10, nil)
		ethClient.On("HeadByNumber", mock.Anything, big.NewInt(8)).
			Return(&head8, nil)

		// Needs to be 8 because there are 8 heads in chain (15,14,13,12,11,10,9,8)
		var depth uint = 8

		err := ht.Backfill(ctx, &h15, depth)
		require.NoError(t, err)

		h := ht.headSaver.Chain(h15.Hash)

		require.Equal(t, uint32(8), h.ChainLength())
		earliestInChain := h.EarliestInChain()
		assert.Equal(t, head8.Number, earliestInChain.BlockNumber())
		assert.Equal(t, head8.Hash, earliestInChain.BlockHash())
	})

	t.Run("does not backfill if chain length is already greater than or equal to depth", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		cfg := configtest.NewGeneralConfig(t, nil)
		logger := logger.TestLogger(t)
		orm := headtracker.NewORM(db, logger, cfg, cltest.FixtureChainID)
		for i := range heads {
			require.NoError(t, orm.IdempotentInsertHead(testutils.Context(t), &heads[i]))
		}

		ethClient := evmtest.NewEthClientMock(t)
		ethClient.On("ConfiguredChainID", mock.Anything).Return(cfg.DefaultChainID(), nil)

		ht := createHeadTrackerWithNeverSleeper(t, ethClient, cfg, orm)

		err := ht.Backfill(ctx, &h15, 3)
		require.NoError(t, err)

		err = ht.Backfill(ctx, &h15, 5)
		require.NoError(t, err)
	})

	t.Run("only backfills to height 0 if chain length would otherwise cause it to try and fetch a negative head", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		cfg := configtest.NewGeneralConfig(t, nil)
		logger := logger.TestLogger(t)
		orm := headtracker.NewORM(db, logger, cfg, cltest.FixtureChainID)

		ethClient := evmtest.NewEthClientMock(t)
		ethClient.On("ConfiguredChainID", mock.Anything).Return(cfg.DefaultChainID(), nil)
		ethClient.On("HeadByNumber", mock.Anything, big.NewInt(0)).
			Return(&head0, nil)

		require.NoError(t, orm.IdempotentInsertHead(testutils.Context(t), &h1))

		ht := createHeadTrackerWithNeverSleeper(t, ethClient, cfg, orm)

		err := ht.Backfill(ctx, &h1, 400)
		require.NoError(t, err)

		h := ht.headSaver.Chain(h1.Hash)
		require.NotNil(t, h)

		require.Equal(t, uint32(2), h.ChainLength())
		require.Equal(t, int64(0), h.EarliestInChain().BlockNumber())
	})

	t.Run("abandons backfill and returns error if the eth node returns not found", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		cfg := configtest.NewGeneralConfig(t, nil)
		logger := logger.TestLogger(t)
		orm := headtracker.NewORM(db, logger, cfg, cltest.FixtureChainID)
		for i := range heads {
			require.NoError(t, orm.IdempotentInsertHead(testutils.Context(t), &heads[i]))
		}

		ethClient := evmtest.NewEthClientMock(t)
		ethClient.On("ConfiguredChainID", mock.Anything).Return(cfg.DefaultChainID(), nil)
		ethClient.On("HeadByNumber", mock.Anything, big.NewInt(10)).
			Return(&head10, nil).
			Once()
		ethClient.On("HeadByNumber", mock.Anything, big.NewInt(8)).
			Return(nil, ethereum.NotFound).
			Once()

		ht := createHeadTrackerWithNeverSleeper(t, ethClient, cfg, orm)

		err := ht.Backfill(ctx, &h12, 400)
		require.Error(t, err)
		require.EqualError(t, err, "fetchAndSaveHead failed: not found")

		h := ht.headSaver.Chain(h12.Hash)

		// Should contain 12, 11, 10, 9
		assert.Equal(t, 4, int(h.ChainLength()))
		assert.Equal(t, int64(9), h.EarliestInChain().BlockNumber())
	})

	t.Run("abandons backfill and returns error if the context time budget is exceeded", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		cfg := configtest.NewGeneralConfig(t, nil)
		logger := logger.TestLogger(t)
		orm := headtracker.NewORM(db, logger, cfg, cltest.FixtureChainID)
		for i := range heads {
			require.NoError(t, orm.IdempotentInsertHead(testutils.Context(t), &heads[i]))
		}

		ethClient := evmtest.NewEthClientMock(t)
		ethClient.On("ConfiguredChainID", mock.Anything).Return(cfg.DefaultChainID(), nil)
		ethClient.On("HeadByNumber", mock.Anything, big.NewInt(10)).
			Return(&head10, nil)
		ethClient.On("HeadByNumber", mock.Anything, big.NewInt(8)).
			Return(nil, context.DeadlineExceeded)

		ht := createHeadTrackerWithNeverSleeper(t, ethClient, cfg, orm)

		err := ht.Backfill(ctx, &h12, 400)
		require.Error(t, err)
		require.EqualError(t, err, "fetchAndSaveHead failed: context deadline exceeded")

		h := ht.headSaver.Chain(h12.Hash)

		// Should contain 12, 11, 10, 9
		assert.Equal(t, 4, int(h.ChainLength()))
		assert.Equal(t, int64(9), h.EarliestInChain().BlockNumber())
	})
}

func createHeadTracker(t *testing.T, ethClient evmclient.Client, config headtracker.Config, orm headtracker.ORM) *headTrackerUniverse {
	lggr := logger.TestLogger(t)
	hb := headtracker.NewHeadBroadcaster(lggr)
	hs := headtracker.NewHeadSaver(lggr, orm, config)
	mailMon := utils.NewMailboxMonitor(t.Name())
	return &headTrackerUniverse{
		mu:              new(sync.Mutex),
		headTracker:     headtracker.NewHeadTracker(lggr, ethClient, config, hb, hs, mailMon),
		headBroadcaster: hb,
		headSaver:       hs,
		mailMon:         mailMon,
	}
}

func createHeadTrackerWithNeverSleeper(t *testing.T, ethClient evmclient.Client, cfg chainlink.GeneralConfig, orm headtracker.ORM) *headTrackerUniverse {
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	lggr := logger.TestLogger(t)
	hb := headtracker.NewHeadBroadcaster(lggr)
	hs := headtracker.NewHeadSaver(lggr, orm, evmcfg)
	mailMon := utils.NewMailboxMonitor(t.Name())
	ht := headtracker.NewHeadTracker(lggr, ethClient, evmcfg, hb, hs, mailMon)
	_, err := hs.LoadFromDB(testutils.Context(t))
	require.NoError(t, err)
	return &headTrackerUniverse{
		mu:              new(sync.Mutex),
		headTracker:     ht,
		headBroadcaster: hb,
		headSaver:       hs,
		mailMon:         mailMon,
	}
}

func createHeadTrackerWithChecker(t *testing.T, ethClient evmclient.Client, config headtracker.Config, orm headtracker.ORM, checker httypes.HeadTrackable) *headTrackerUniverse {
	lggr := logger.TestLogger(t)
	hb := headtracker.NewHeadBroadcaster(lggr)
	hs := headtracker.NewHeadSaver(lggr, orm, config)
	hb.Subscribe(checker)
	mailMon := utils.NewMailboxMonitor(t.Name())
	ht := headtracker.NewHeadTracker(lggr, ethClient, config, hb, hs, mailMon)
	return &headTrackerUniverse{
		mu:              new(sync.Mutex),
		headTracker:     ht,
		headBroadcaster: hb,
		headSaver:       hs,
		mailMon:         mailMon,
	}
}

type headTrackerUniverse struct {
	mu              *sync.Mutex
	stopped         bool
	headTracker     httypes.HeadTracker
	headBroadcaster httypes.HeadBroadcaster
	headSaver       httypes.HeadSaver
	mailMon         *utils.MailboxMonitor
}

func (u *headTrackerUniverse) Backfill(ctx context.Context, head *evmtypes.Head, depth uint) error {
	return u.headTracker.Backfill(ctx, head, depth)
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
