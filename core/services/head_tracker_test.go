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

func TestHeadTracker_New(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	cltest.MockEthOnStore(t, store, cltest.EthMockRegisterChainID)

	assert.Nil(t, store.CreateHead(cltest.Head(1)))
	last := cltest.Head(16)
	assert.Nil(t, store.CreateHead(last))
	assert.Nil(t, store.CreateHead(cltest.Head(10)))

	ht := services.NewHeadTracker(store, []strpkg.HeadTrackable{})
	assert.Nil(t, ht.Start())
	assert.Equal(t, last.Number, ht.Head().Number)
}

func TestHeadTracker_New_Limit_At_100(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	cltest.MockEthOnStore(t, store, cltest.EthMockRegisterChainID)

	for idx := 0; idx <= 200; idx++ {
		assert.Nil(t, store.CreateHead(cltest.Head(idx)))
	}
	firstHead, err := store.FirstHead()
	assert.Nil(t, err)
	assert.Equal(t, big.NewInt(101), firstHead.ToInt())
	ht := services.NewHeadTracker(store, []strpkg.HeadTrackable{})
	assert.Nil(t, ht.Start())
	assert.Equal(t, big.NewInt(200), ht.Head().ToInt())
}

func TestHeadTracker_Get(t *testing.T) {
	t.Parallel()

	start := cltest.Head(5)

	tests := []struct {
		name      string
		initial   *models.Head
		toSave    *models.Head
		want      *big.Int
		wantError bool
	}{
		{"greater", start, cltest.Head(6), big.NewInt(6), false},
		{"less than", start, cltest.Head(1), big.NewInt(5), true},
		{"zero", start, cltest.Head(0), big.NewInt(5), true},
		{"nil", start, nil, big.NewInt(5), true},
		{"nil no initial", nil, nil, nil, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			cltest.MockEthOnStore(t, store, cltest.EthMockRegisterChainID)
			if test.initial != nil {
				assert.Nil(t, store.CreateHead(test.initial))
			}

			ht := services.NewHeadTracker(store, []strpkg.HeadTrackable{})
			ht.Start()
			defer ht.Stop()

			err := ht.Save(test.toSave)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.want, ht.Head().ToInt())
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

func TestHeadTracker_HeadTrackableCallbacks(t *testing.T) {
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
	assert.Equal(t, int32(0), checker.OnNewHeadCount())

	headers <- eth.BlockHeader{Number: cltest.BigHexInt(1)}
	g.Eventually(func() int32 { return checker.OnNewHeadCount() }).Should(gomega.Equal(int32(1)))
	assert.Equal(t, int32(1), checker.ConnectedCount())
	assert.Equal(t, int32(0), checker.DisconnectedCount())

	ht.Stop()
	assert.Equal(t, int32(1), checker.DisconnectedCount())
	assert.Equal(t, int32(1), checker.ConnectedCount())
	assert.Equal(t, int32(1), checker.OnNewHeadCount())
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
	assert.Equal(t, int32(0), checker.OnNewHeadCount())

	// trigger reconnect loop
	subscription.Errors <- errors.New("Test error to force reconnect")
	g.Eventually(func() int32 { return checker.ConnectedCount() }).Should(gomega.Equal(int32(2)))
	g.Consistently(func() int32 { return checker.ConnectedCount() }).Should(gomega.Equal(int32(2)))
	assert.Equal(t, int32(1), checker.DisconnectedCount())
	assert.Equal(t, int32(0), checker.OnNewHeadCount())

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

	require.NoError(t, ht.Save(models.NewHead(lastSavedBN, cltest.NewHash())))

	assert.Nil(t, ht.Start())
	headers <- eth.BlockHeader{Number: hexutil.Big(*currentBN)}
	g.Eventually(func() int32 { return checker.ConnectedCount() }).Should(gomega.Equal(int32(1)))

	connectedBN := connectedValue.Load().(*big.Int)
	assert.Equal(t, lastSavedBN, connectedBN)
	g.Eventually(func() *big.Int { return ht.Head().ToInt() }).Should(gomega.Equal(currentBN))
	assert.NoError(t, ht.Stop())
}
