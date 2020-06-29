package services_test

import (
	"context"
	"errors"
	"math/big"
	"sync/atomic"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"

	gethCommon "github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func firstHead(t *testing.T, store *strpkg.Store) models.Head {
	h := models.Head{}
	if err := store.GetRawDB().Order("number asc").First(&h).Error; err != nil {
		t.Fatal(err)
	}
	return h
}

func TestHeadTracker_New(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	cltest.MockEthOnStore(t, store, cltest.EthMockRegisterChainID)

	assert.Nil(t, store.IdempotentInsertHead(*cltest.Head(1)))
	last := cltest.Head(16)
	assert.Nil(t, store.IdempotentInsertHead(*last))
	assert.Nil(t, store.IdempotentInsertHead(*cltest.Head(10)))

	ht := services.NewHeadTracker(store, []strpkg.HeadTrackable{})
	assert.Nil(t, ht.Start())
	assert.Equal(t, last.Number, ht.HighestSeenHead().Number)
}

func TestHeadTracker_Save_InsertsAndTrimsTable(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	store.Config.Set("ETH_HEAD_TRACKER_HISTORY_DEPTH", 100)
	defer cleanup()

	cltest.MockEthOnStore(t, store, cltest.EthMockRegisterChainID)

	for idx := 0; idx < 200; idx++ {
		assert.Nil(t, store.IdempotentInsertHead(*cltest.Head(idx)))
	}

	ht := services.NewHeadTracker(store, []strpkg.HeadTrackable{})

	h := cltest.Head(200)
	require.NoError(t, ht.Save(*h))
	assert.Equal(t, big.NewInt(200), ht.HighestSeenHead().ToInt())

	firstHead := firstHead(t, store)
	assert.Equal(t, big.NewInt(101), firstHead.ToInt())

	lastHead, err := store.LastHead()
	require.NoError(t, err)
	assert.Equal(t, int64(200), lastHead.Number)
}

func TestHeadTracker_Get(t *testing.T) {
	t.Parallel()

	start := cltest.Head(5)

	tests := []struct {
		name    string
		initial *models.Head
		toSave  *models.Head
		want    *big.Int
	}{
		{"greater", start, cltest.Head(6), big.NewInt(6)},
		{"less than", start, cltest.Head(1), big.NewInt(5)},
		{"zero", start, cltest.Head(0), big.NewInt(5)},
		{"nil", start, nil, big.NewInt(5)},
		{"nil no initial", nil, nil, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			cltest.MockEthOnStore(t, store, cltest.EthMockRegisterChainID)
			if test.initial != nil {
				assert.Nil(t, store.IdempotentInsertHead(*test.initial))
			}

			ht := services.NewHeadTracker(store, []strpkg.HeadTrackable{})
			ht.Start()
			defer ht.Stop()

			if test.toSave != nil {
				err := ht.Save(*test.toSave)
				assert.NoError(t, err)
			}

			assert.Equal(t, test.want, ht.HighestSeenHead().ToInt())
		})
	}
}

func TestHeadTracker_Start_NewHeads(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	eth := cltest.MockEthOnStore(t, store, cltest.EthMockRegisterChainID)
	ht := services.NewHeadTracker(store, []strpkg.HeadTrackable{})
	defer ht.Stop()

	eth.RegisterSubscription("newHeads")

	assert.Nil(t, ht.Start())
	eth.EventuallyAllCalled(t)
}

func TestHeadTracker_CallsHeadTrackableCallbacks(t *testing.T) {
	t.Parallel()
	g := gomega.NewGomegaWithT(t)

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	mocketh := cltest.MockEthOnStore(t, store, cltest.EthMockRegisterChainID)

	checker := &cltest.MockHeadTrackable{}
	ht := services.NewHeadTracker(store, []strpkg.HeadTrackable{checker}, cltest.NeverSleeper{})

	headers := make(chan gethTypes.Header)
	mocketh.RegisterSubscription("newHeads", headers)
	mocketh.Register("eth_chainId", store.Config.ChainID())

	assert.Nil(t, ht.Start())
	g.Eventually(func() int32 { return checker.ConnectedCount() }).Should(gomega.Equal(int32(1)))
	assert.Equal(t, int32(0), checker.DisconnectedCount())
	assert.Equal(t, int32(0), checker.OnNewLongestChainCount())

	headers <- gethTypes.Header{Number: big.NewInt(1)}
	g.Eventually(func() int32 { return checker.OnNewLongestChainCount() }).Should(gomega.Equal(int32(1)))
	assert.Equal(t, int32(1), checker.ConnectedCount())
	assert.Equal(t, int32(0), checker.DisconnectedCount())

	require.NoError(t, ht.Stop())
	assert.Equal(t, int32(1), checker.DisconnectedCount())
	assert.Equal(t, int32(1), checker.ConnectedCount())
	assert.Equal(t, int32(1), checker.OnNewLongestChainCount())
}

func TestHeadTracker_ReconnectOnError(t *testing.T) {
	t.Parallel()
	g := gomega.NewGomegaWithT(t)

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	txManager := new(mocks.TxManager)
	subscription := cltest.EmptyMockSubscription()
	txManager.On("GetChainID").Maybe().Return(store.Config.ChainID(), nil)
	txManager.On("SubscribeToNewHeads", mock.Anything, mock.Anything, mock.Anything).Return(subscription, nil)
	txManager.On("SubscribeToNewHeads", mock.Anything, mock.Anything).Return(nil, errors.New("cannot reconnect"))
	txManager.On("SubscribeToNewHeads", mock.Anything, mock.Anything).Return(subscription, nil)
	store.TxManager = txManager

	checker := &cltest.MockHeadTrackable{}
	ht := services.NewHeadTracker(store, []strpkg.HeadTrackable{checker}, cltest.NeverSleeper{})

	// connect
	assert.Nil(t, ht.Start())
	g.Eventually(func() int32 { return checker.ConnectedCount() }).Should(gomega.Equal(int32(1)))
	assert.Equal(t, int32(0), checker.DisconnectedCount())
	assert.Equal(t, int32(0), checker.OnNewLongestChainCount())

	// trigger reconnect loop
	subscription.Errors <- errors.New("Test error to force reconnect")
	g.Eventually(func() int32 { return checker.ConnectedCount() }).Should(gomega.Equal(int32(2)))
	g.Consistently(func() int32 { return checker.ConnectedCount() }).Should(gomega.Equal(int32(2)))
	assert.Equal(t, int32(1), checker.DisconnectedCount())
	assert.Equal(t, int32(0), checker.OnNewLongestChainCount())

	// stop
	assert.NoError(t, ht.Stop())
}

func TestHeadTracker_StartConnectsFromLastSavedHeader(t *testing.T) {
	t.Parallel()
	g := gomega.NewGomegaWithT(t)

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	mocketh := cltest.MockEthOnStore(t, store, cltest.EthMockRegisterChainID)
	headers := make(chan gethTypes.Header)
	mocketh.RegisterSubscription("newHeads", headers)

	lastSavedBN := big.NewInt(1)
	currentBN := big.NewInt(2)
	var connectedValue atomic.Value

	checker := &cltest.MockHeadTrackable{ConnectedCallback: func(bn *models.Head) {
		connectedValue.Store(bn.ToInt())
	}}
	ht := services.NewHeadTracker(store, []strpkg.HeadTrackable{checker}, cltest.NeverSleeper{})

	require.NoError(t, ht.Save(models.NewHead(lastSavedBN, cltest.NewHash(), cltest.NewHash(), 0)))

	assert.Nil(t, ht.Start())
	headers <- gethTypes.Header{Number: currentBN}
	g.Eventually(func() int32 { return checker.ConnectedCount() }).Should(gomega.Equal(int32(1)))

	connectedBN := connectedValue.Load().(*big.Int)
	assert.Equal(t, lastSavedBN, connectedBN)

	assert.NoError(t, ht.Stop())

	// Check that it saved the head
	h, err := store.LastHead()
	require.NoError(t, err)
	assert.Equal(t, h.Number, currentBN.Int64())
}

func TestHeadTracker_SwitchesToLongestChain(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	mocketh := cltest.MockEthOnStore(t, store, cltest.EthMockRegisterChainID)

	checker := new(mocks.HeadTrackable)
	ht := services.NewHeadTracker(store, []strpkg.HeadTrackable{checker}, cltest.NeverSleeper{})

	headers := make(chan gethTypes.Header)
	mocketh.RegisterSubscription("newHeads", headers)
	mocketh.Register("eth_chainId", store.Config.ChainID())

	checker.On("Connect", mock.MatchedBy(func(h *models.Head) bool {
		return h == nil
	})).Return(nil).Once()
	checker.On("Disconnect").Return(nil).Once()

	assert.Nil(t, ht.Start())

	lastHead := make(chan struct{})
	blockHeaders := []gethTypes.Header{}

	// First block comes in
	blockHeaders = append(blockHeaders, gethTypes.Header{Number: big.NewInt(1), ParentHash: cltest.NewHash(), Time: 1})
	// Blocks 2 and 3 are out of order
	head2 := gethTypes.Header{Number: big.NewInt(2), ParentHash: blockHeaders[0].Hash(), Time: 2}
	head3 := gethTypes.Header{Number: big.NewInt(3), ParentHash: head2.Hash(), Time: 3}
	blockHeaders = append(blockHeaders, head3)
	blockHeaders = append(blockHeaders, head2)
	// Block 4 comes in
	blockHeaders = append(blockHeaders, gethTypes.Header{Number: big.NewInt(4), ParentHash: blockHeaders[1].Hash(), Time: 4})
	// Another block at level 4 comes in, that will be uncled
	blockHeaders = append(blockHeaders, gethTypes.Header{Number: big.NewInt(4), ParentHash: blockHeaders[1].Hash(), Time: 5})
	// Reorg happened forking from block 2
	blockHeaders = append(blockHeaders, gethTypes.Header{Number: big.NewInt(2), ParentHash: blockHeaders[0].Hash(), Time: 6})
	blockHeaders = append(blockHeaders, gethTypes.Header{Number: big.NewInt(3), ParentHash: blockHeaders[5].Hash(), Time: 7})
	blockHeaders = append(blockHeaders, gethTypes.Header{Number: big.NewInt(4), ParentHash: blockHeaders[6].Hash(), Time: 8})
	// Now the new chain is longer
	blockHeaders = append(blockHeaders, gethTypes.Header{Number: big.NewInt(5), ParentHash: blockHeaders[7].Hash(), Time: 9})

	checker.On("OnNewLongestChain", mock.MatchedBy(func(h models.Head) bool {
		return h.Number == 1 && h.Hash == blockHeaders[0].Hash()
	})).Return().Once()
	checker.On("OnNewLongestChain", mock.MatchedBy(func(h models.Head) bool {
		return h.Number == 3 && h.Hash == blockHeaders[1].Hash()
	})).Return().Once()
	checker.On("OnNewLongestChain", mock.MatchedBy(func(h models.Head) bool {
		if h.Number == 4 && h.Hash == blockHeaders[3].Hash() {
			// Check that the block came with its parents
			require.NotNil(t, h.Parent)
			require.Equal(t, h.Parent.Hash, blockHeaders[1].Hash())
			require.NotNil(t, h.Parent.Parent.Hash)
			require.Equal(t, h.Parent.Parent.Hash, blockHeaders[2].Hash())
			require.NotNil(t, h.Parent.Parent.Parent)
			require.NotNil(t, h.Parent.Parent.Parent.Hash, blockHeaders[0].Hash())
			return true
		}
		return false
	})).Return().Once()
	checker.On("OnNewLongestChain", mock.MatchedBy(func(h models.Head) bool {
		if h.Number == 5 && h.Hash == blockHeaders[8].Hash() {
			// This is the new longest chain, check that it came with its parents
			require.NotNil(t, h.Parent)
			require.Equal(t, h.Parent.Hash, blockHeaders[7].Hash())
			require.NotNil(t, h.Parent.Parent.Hash)
			require.Equal(t, h.Parent.Parent.Hash, blockHeaders[6].Hash())
			require.NotNil(t, h.Parent.Parent.Parent)
			require.NotNil(t, h.Parent.Parent.Parent.Hash, blockHeaders[5].Hash())
			require.NotNil(t, h.Parent.Parent.Parent.Parent)
			require.NotNil(t, h.Parent.Parent.Parent.Parent.Hash, blockHeaders[0].Hash())

			return true
		}
		return false
	})).Return().Once().Run(func(_ mock.Arguments) {
		close(lastHead)
	})

	for _, h := range blockHeaders {
		headers <- h
	}

	gomega.NewGomegaWithT(t).Eventually(lastHead).Should(gomega.BeClosed())
	require.NoError(t, ht.Stop())
	assert.Equal(t, int64(5), ht.HighestSeenHead().Number)

	for _, h := range blockHeaders {
		c, err := store.Chain(h.Hash(), 1)
		require.NoError(t, err)
		require.NotNil(t, c)
		assert.Equal(t, c.ParentHash, h.ParentHash)
		assert.Equal(t, c.Timestamp.Unix(), int64(h.Time))
		assert.Equal(t, c.Number, h.Number.Int64())
	}

	checker.AssertExpectations(t)
}

func TestHeadTracker_GetChainWithBackfill(t *testing.T) {
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

	head0 := gethTypes.Header{
		Number:     big.NewInt(0),
		ParentHash: gethCommon.BigToHash(big.NewInt(0)),
	}
	h1 := *cltest.Head(1)
	h1.ParentHash = head0.Hash()

	head8 := gethTypes.Header{
		Number:     big.NewInt(8),
		ParentHash: cltest.NewHash(),
	}

	h9 := *cltest.Head(9)
	h9.ParentHash = head8.Hash()

	head10 := gethTypes.Header{
		Number:     big.NewInt(10),
		ParentHash: h9.Hash,
	}

	h11 := *cltest.Head(11)
	h11.ParentHash = head10.Hash()

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

	heads := []models.Head{
		h9,
		h11,
		h12,
		h13,
		h14Orphaned,
		h14,
		h15,
	}

	ctx := context.Background()

	t.Run("returns chain if all the heads are in database", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		for _, h := range heads {
			require.NoError(t, store.IdempotentInsertHead(h))
		}

		mocketh := cltest.MockEthOnStore(t, store, cltest.EthMockRegisterChainID)
		ht := services.NewHeadTracker(store, []strpkg.HeadTrackable{}, cltest.NeverSleeper{})

		mocketh.RegisterSubscription("newHeads")
		mocketh.Register("eth_chainId", store.Config.ChainID())

		h, err := ht.GetChainWithBackfill(ctx, h12, 2)
		require.NoError(t, err)

		assert.Equal(t, int64(12), h.Number)
		require.NotNil(t, h.Parent)
		assert.Equal(t, int64(11), h.Parent.Number)
		require.Nil(t, h.Parent.Parent)
	})

	t.Run("fetches a missing head", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		for _, h := range heads {
			require.NoError(t, store.IdempotentInsertHead(h))
		}

		gethClient := new(mocks.GethClient)
		store.GethClientWrapper = cltest.NewSimpleGethWrapper(gethClient)

		mocketh := cltest.MockEthOnStore(t, store, cltest.EthMockRegisterChainID)
		ht := services.NewHeadTracker(store, []strpkg.HeadTrackable{}, cltest.NeverSleeper{})

		mocketh.RegisterSubscription("newHeads")
		mocketh.Register("eth_chainId", store.Config.ChainID())

		gethClient.On("HeaderByNumber", mock.Anything, mock.MatchedBy(func(b *big.Int) bool {
			return b.Cmp(big.NewInt(10)) == 0
		})).Return(&head10, nil).Once()

		h, err := ht.GetChainWithBackfill(ctx, h12, 3)
		require.NoError(t, err)

		assert.Equal(t, int64(12), h.Number)
		require.NotNil(t, h.Parent)
		assert.Equal(t, int64(11), h.Parent.Number)
		require.NotNil(t, h.Parent)
		assert.Equal(t, int64(10), h.Parent.Parent.Number)
		require.Nil(t, h.Parent.Parent.Parent)

		writtenHead, err := store.HeadByHash(head10.Hash())
		require.NoError(t, err)
		assert.Equal(t, int64(10), writtenHead.Number)

		gethClient.AssertExpectations(t)
	})

	t.Run("fetches only heads that are missing", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		for _, h := range heads {
			require.NoError(t, store.IdempotentInsertHead(h))
		}

		gethClient := new(mocks.GethClient)
		store.GethClientWrapper = cltest.NewSimpleGethWrapper(gethClient)

		mocketh := cltest.MockEthOnStore(t, store, cltest.EthMockRegisterChainID)
		ht := services.NewHeadTracker(store, []strpkg.HeadTrackable{}, cltest.NeverSleeper{})

		mocketh.RegisterSubscription("newHeads")
		mocketh.Register("eth_chainId", store.Config.ChainID())

		gethClient.On("HeaderByNumber", mock.Anything, mock.MatchedBy(func(b *big.Int) bool {
			return b.Cmp(big.NewInt(10)) == 0
		})).Return(&head10, nil).Once()
		gethClient.On("HeaderByNumber", mock.Anything, mock.MatchedBy(func(b *big.Int) bool {
			return b.Cmp(big.NewInt(8)) == 0
		})).Return(&head8, nil).Once()

		// Needs to be 8 because there are 8 heads in chain (15,14,13,12,11,10,9,8)
		h, err := ht.GetChainWithBackfill(ctx, h15, 8)
		require.NoError(t, err)

		require.Equal(t, uint32(8), h.ChainLength())
		earliestInChain := h.EarliestInChain()
		assert.Equal(t, head8.Number.Int64(), earliestInChain.Number)
		assert.Equal(t, head8.Hash(), earliestInChain.Hash)

		gethClient.AssertExpectations(t)
	})

	t.Run("returns error if first head is not in database", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()

		gethClient := new(mocks.GethClient)
		store.GethClientWrapper = cltest.NewSimpleGethWrapper(gethClient)

		mocketh := cltest.MockEthOnStore(t, store, cltest.EthMockRegisterChainID)
		ht := services.NewHeadTracker(store, []strpkg.HeadTrackable{}, cltest.NeverSleeper{})

		mocketh.RegisterSubscription("newHeads")
		mocketh.Register("eth_chainId", store.Config.ChainID())

		h16 := *cltest.Head(16)
		h16.ParentHash = h15.Hash

		_, err := ht.GetChainWithBackfill(ctx, h16, 3)
		require.Contains(t, err.Error(), "record not found")

		gethClient.AssertExpectations(t)
	})

	t.Run("does not backfill if chain length is already greater than or equal to depth", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		for _, h := range heads {
			require.NoError(t, store.IdempotentInsertHead(h))
		}

		gethClient := new(mocks.GethClient)
		store.GethClientWrapper = cltest.NewSimpleGethWrapper(gethClient)

		mocketh := cltest.MockEthOnStore(t, store, cltest.EthMockRegisterChainID)
		ht := services.NewHeadTracker(store, []strpkg.HeadTrackable{}, cltest.NeverSleeper{})

		mocketh.RegisterSubscription("newHeads")
		mocketh.Register("eth_chainId", store.Config.ChainID())

		h, err := ht.GetChainWithBackfill(ctx, h15, 3)
		require.NoError(t, err)
		require.Equal(t, uint32(3), h.ChainLength())

		h, err = ht.GetChainWithBackfill(ctx, h15, 5)
		require.NoError(t, err)
		require.Equal(t, uint32(5), h.ChainLength())

		gethClient.AssertExpectations(t)
	})

	t.Run("only backfills to height 0 if chain length would otherwise cause it to try and fetch a negative head", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()

		gethClient := new(mocks.GethClient)
		store.GethClientWrapper = cltest.NewSimpleGethWrapper(gethClient)

		mocketh := cltest.MockEthOnStore(t, store, cltest.EthMockRegisterChainID)
		ht := services.NewHeadTracker(store, []strpkg.HeadTrackable{}, cltest.NeverSleeper{})

		mocketh.RegisterSubscription("newHeads")
		mocketh.Register("eth_chainId", store.Config.ChainID())

		require.NoError(t, store.IdempotentInsertHead(h1))

		gethClient.On("HeaderByNumber", mock.Anything, mock.MatchedBy(func(b *big.Int) bool {
			return b.Cmp(big.NewInt(0)) == 0
		})).Return(&head0, nil).Once()

		h, err := ht.GetChainWithBackfill(ctx, h1, 400)
		require.NoError(t, err)
		require.Equal(t, uint32(2), h.ChainLength())
		require.Equal(t, int64(0), h.EarliestInChain().Number)

		gethClient.AssertExpectations(t)
	})

	t.Run("abandons backfill and returns whatever we have if the eth node returns not found", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		for _, h := range heads {
			require.NoError(t, store.IdempotentInsertHead(h))
		}

		gethClient := new(mocks.GethClient)
		store.GethClientWrapper = cltest.NewSimpleGethWrapper(gethClient)

		mocketh := cltest.MockEthOnStore(t, store, cltest.EthMockRegisterChainID)
		ht := services.NewHeadTracker(store, []strpkg.HeadTrackable{}, cltest.NeverSleeper{})

		mocketh.RegisterSubscription("newHeads")
		mocketh.Register("eth_chainId", store.Config.ChainID())

		gethClient.On("HeaderByNumber", mock.Anything, mock.MatchedBy(func(b *big.Int) bool {
			return b.Cmp(big.NewInt(10)) == 0
		})).Return(&head10, nil).Once()
		gethClient.On("HeaderByNumber", mock.Anything, mock.MatchedBy(func(b *big.Int) bool {
			return b.Cmp(big.NewInt(8)) == 0
		})).Return(nil, errors.New("not found")).Once()

		h, err := ht.GetChainWithBackfill(ctx, h12, 400)
		require.NoError(t, err)

		// Should contain 12, 11, 10, 9
		assert.Equal(t, 4, int(h.ChainLength()))
		assert.Equal(t, int64(9), h.EarliestInChain().Number)

		gethClient.AssertExpectations(t)
	})

	t.Run("abandons backfill and returns whatever we have if the context time budget is exceeded", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		for _, h := range heads {
			require.NoError(t, store.IdempotentInsertHead(h))
		}

		gethClient := new(mocks.GethClient)
		store.GethClientWrapper = cltest.NewSimpleGethWrapper(gethClient)

		mocketh := cltest.MockEthOnStore(t, store, cltest.EthMockRegisterChainID)
		ht := services.NewHeadTracker(store, []strpkg.HeadTrackable{}, cltest.NeverSleeper{})

		mocketh.RegisterSubscription("newHeads")
		mocketh.Register("eth_chainId", store.Config.ChainID())

		gethClient.On("HeaderByNumber", mock.Anything, mock.MatchedBy(func(b *big.Int) bool {
			return b.Cmp(big.NewInt(10)) == 0
		})).Return(&head10, nil).Once()
		gethClient.On("HeaderByNumber", mock.Anything, mock.MatchedBy(func(b *big.Int) bool {
			return b.Cmp(big.NewInt(8)) == 0
		})).Return(nil, context.DeadlineExceeded).Once()

		h, err := ht.GetChainWithBackfill(ctx, h12, 400)
		require.NoError(t, err)

		// Should contain 12, 11, 10, 9
		assert.Equal(t, 4, int(h.ChainLength()))
		assert.Equal(t, int64(9), h.EarliestInChain().Number)

		gethClient.AssertExpectations(t)
	})
}
