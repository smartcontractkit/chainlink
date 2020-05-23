package services_test

import (
	"errors"
	"math/big"
	"sync/atomic"
	"testing"

	"github.com/smartcontractkit/chainlink/core/eth"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/ethereum/go-ethereum/common/hexutil"
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

	headers := make(chan eth.BlockHeader)
	mocketh.RegisterSubscription("newHeads", headers)
	mocketh.Register("eth_chainId", store.Config.ChainID())

	assert.Nil(t, ht.Start())
	g.Eventually(func() int32 { return checker.ConnectedCount() }).Should(gomega.Equal(int32(1)))
	assert.Equal(t, int32(0), checker.DisconnectedCount())
	assert.Equal(t, int32(0), checker.OnNewLongestChainCount())

	headers <- eth.BlockHeader{Number: cltest.BigHexInt(1)}
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
	headers := make(chan eth.BlockHeader)
	mocketh.RegisterSubscription("newHeads", headers)

	lastSavedBN := big.NewInt(1)
	currentBN := big.NewInt(2)
	var connectedValue atomic.Value

	checker := &cltest.MockHeadTrackable{ConnectedCallback: func(bn *models.Head) {
		connectedValue.Store(bn.ToInt())
	}}
	ht := services.NewHeadTracker(store, []strpkg.HeadTrackable{checker}, cltest.NeverSleeper{})

	require.NoError(t, ht.Save(models.NewHead(lastSavedBN, cltest.NewHash(), cltest.NewHash(), big.NewInt(0))))

	assert.Nil(t, ht.Start())
	headers <- eth.BlockHeader{Number: hexutil.Big(*currentBN)}
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

	headers := make(chan eth.BlockHeader)
	mocketh.RegisterSubscription("newHeads", headers)
	mocketh.Register("eth_chainId", store.Config.ChainID())

	checker.On("Connect", mock.MatchedBy(func(h *models.Head) bool {
		return h == nil
	})).Return(nil).Once()
	checker.On("Disconnect").Return(nil).Once()

	assert.Nil(t, ht.Start())

	lastHead := make(chan struct{})
	blockHeaders := []eth.BlockHeader{}

	// First block comes in
	blockHeaders = append(blockHeaders, eth.BlockHeader{Number: cltest.BigHexInt(1), ParentHash: cltest.NewHash(), GethHash: cltest.NewHash(), Time: cltest.BigHexInt(1)})
	// Blocks 2 and 3 are out of order
	block2Hash := cltest.NewHash()
	blockHeaders = append(blockHeaders, eth.BlockHeader{Number: cltest.BigHexInt(3), ParentHash: block2Hash, GethHash: cltest.NewHash(), Time: cltest.BigHexInt(3)})
	blockHeaders = append(blockHeaders, eth.BlockHeader{Number: cltest.BigHexInt(2), ParentHash: blockHeaders[0].Hash(), GethHash: block2Hash, Time: cltest.BigHexInt(2)})
	// Block 4 comes in
	blockHeaders = append(blockHeaders, eth.BlockHeader{Number: cltest.BigHexInt(4), ParentHash: blockHeaders[1].Hash(), GethHash: cltest.NewHash(), Time: cltest.BigHexInt(4)})
	// Another block at level 4 comes in, that will be uncled
	blockHeaders = append(blockHeaders, eth.BlockHeader{Number: cltest.BigHexInt(4), ParentHash: blockHeaders[1].Hash(), GethHash: cltest.NewHash(), Time: cltest.BigHexInt(5)})
	// Reorg happened forking from block 2
	blockHeaders = append(blockHeaders, eth.BlockHeader{Number: cltest.BigHexInt(2), ParentHash: blockHeaders[0].Hash(), GethHash: cltest.NewHash(), Time: cltest.BigHexInt(6)})
	blockHeaders = append(blockHeaders, eth.BlockHeader{Number: cltest.BigHexInt(3), ParentHash: blockHeaders[5].Hash(), GethHash: cltest.NewHash(), Time: cltest.BigHexInt(7)})
	blockHeaders = append(blockHeaders, eth.BlockHeader{Number: cltest.BigHexInt(4), ParentHash: blockHeaders[6].Hash(), GethHash: cltest.NewHash(), Time: cltest.BigHexInt(8)})
	// Now the new chain is longer
	blockHeaders = append(blockHeaders, eth.BlockHeader{Number: cltest.BigHexInt(5), ParentHash: blockHeaders[7].Hash(), GethHash: cltest.NewHash(), Time: cltest.BigHexInt(9)})

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
		assert.Equal(t, c.Timestamp.Unix(), h.Time.ToInt().Int64())
		assert.Equal(t, c.Number, h.Number.ToInt().Int64())
	}

	checker.AssertExpectations(t)
}
