package services_test

import (
	"errors"
	"math/big"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadTracker_New(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	cltest.MockEthOnStore(store)
	assert.Nil(t, store.SaveHead(cltest.Head(1)))
	last := cltest.Head(16)
	assert.Nil(t, store.SaveHead(last))
	assert.Nil(t, store.SaveHead(cltest.Head(10)))

	ht := services.NewHeadTracker(store)
	assert.Nil(t, ht.Start())
	assert.Equal(t, last.Number, ht.Head().Number)
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
			store, cleanup := cltest.NewStore()
			defer cleanup()

			cltest.MockEthOnStore(store)
			if test.initial != nil {
				assert.Nil(t, store.SaveHead(test.initial))
			}

			ht := services.NewHeadTracker(store)
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

	store, cleanup := cltest.NewStore()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)
	ht := services.NewHeadTracker(store)
	defer ht.Stop()

	eth.RegisterSubscription("newHeads")

	assert.Nil(t, ht.Start())
	eth.EventuallyAllCalled(t)
}

func TestHeadTracker_HeadTrackableCallbacks(t *testing.T) {
	t.Parallel()
	g := gomega.NewGomegaWithT(t)

	store, cleanup := cltest.NewStore()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)
	ht := services.NewHeadTracker(store, cltest.NeverSleeper{})

	checker := &cltest.MockHeadTrackable{}
	ht.Attach(checker)

	headers := make(chan models.BlockHeader)
	eth.RegisterSubscription("newHeads", headers)

	assert.Nil(t, ht.Start())
	g.Eventually(func() int32 { return checker.ConnectedCount() }).Should(gomega.Equal(int32(1)))
	assert.Equal(t, int32(0), checker.DisconnectedCount())
	assert.Equal(t, int32(0), checker.OnNewHeadCount())

	headers <- models.BlockHeader{Number: cltest.BigHexInt(1)}
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

	store, cleanup := cltest.NewStore()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)
	ht := services.NewHeadTracker(store, cltest.NeverSleeper{})

	firstSub := eth.RegisterSubscription("newHeads")
	headers := make(chan models.BlockHeader)
	eth.RegisterSubscription("newHeads", headers)

	checker := &cltest.MockHeadTrackable{}
	ht.Attach(checker)

	// connect
	assert.Nil(t, ht.Start())
	g.Eventually(func() int32 { return checker.ConnectedCount() }).Should(gomega.Equal(int32(1)))
	assert.Equal(t, int32(0), checker.DisconnectedCount())
	assert.Equal(t, int32(0), checker.OnNewHeadCount())

	// disconnect
	firstSub.Errors <- errors.New("Test error to force reconnect")
	g.Eventually(func() int32 { return checker.ConnectedCount() }).Should(gomega.Equal(int32(2)))
	assert.Equal(t, int32(1), checker.DisconnectedCount())
	assert.Equal(t, int32(0), checker.OnNewHeadCount())

	// new head
	headers <- models.BlockHeader{Number: cltest.BigHexInt(1)}
	g.Eventually(func() int32 { return checker.OnNewHeadCount() }).Should(gomega.Equal(int32(1)))
	assert.Equal(t, int32(2), checker.ConnectedCount())
	assert.Equal(t, int32(1), checker.DisconnectedCount())
}

func TestHeadTracker_ReconnectAndStopDoesntDeadlock(t *testing.T) {
	t.Parallel()
	g := gomega.NewGomegaWithT(t)

	store, cleanup := cltest.NewStore()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)
	eth.NoMagic()
	ht := services.NewHeadTracker(store, cltest.NeverSleeper{})

	firstConnection := eth.RegisterSubscription("newHeads")
	checker := &cltest.MockHeadTrackable{}
	ht.Attach(checker)

	// connect
	assert.Nil(t, ht.Start())
	g.Eventually(func() int32 { return checker.ConnectedCount() }).Should(gomega.Equal(int32(1)))
	assert.Equal(t, int32(0), checker.DisconnectedCount())
	assert.Equal(t, int32(0), checker.OnNewHeadCount())

	// trigger reconnect loop
	firstConnection.Errors <- errors.New("Test error to force reconnect")
	g.Consistently(func() int32 { return checker.ConnectedCount() }).Should(gomega.Equal(int32(1)))
	assert.Equal(t, int32(1), checker.DisconnectedCount())
	assert.Equal(t, int32(0), checker.OnNewHeadCount())

	// stop
	assert.NoError(t, ht.Stop())
}

func TestHeadTracker_StartConnectsFromLastSavedHeader(t *testing.T) {
	t.Parallel()
	g := gomega.NewGomegaWithT(t)

	store, cleanup := cltest.NewStore()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)
	headers := make(chan models.BlockHeader)
	eth.RegisterSubscription("newHeads", headers)

	ht := services.NewHeadTracker(store, cltest.NeverSleeper{})

	lastSavedBN := big.NewInt(1)
	currentBN := big.NewInt(2)
	var connectedValue atomic.Value

	require.NoError(t, ht.Save(models.NewHead(lastSavedBN, cltest.NewHash())))
	checker := &cltest.MockHeadTrackable{ConnectedCallback: func(bn *models.Head) {
		connectedValue.Store(bn.ToInt())
	}}

	ht.Attach(checker)

	assert.Nil(t, ht.Start())
	headers <- models.BlockHeader{Number: hexutil.Big(*currentBN)}
	g.Eventually(func() int32 { return checker.ConnectedCount() }).Should(gomega.Equal(int32(1)))

	connectedBN := connectedValue.Load().(*big.Int)
	assert.Equal(t, lastSavedBN, connectedBN)
	g.Eventually(func() *big.Int { return ht.Head().ToInt() }).Should(gomega.Equal(currentBN))
	assert.NoError(t, ht.Stop())
}

func TestHeadTracker_OnlyCallsConnectCallbackOnce(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	cltest.MockEthOnStore(store)

	heads := make(chan int)
	headTrackable := &lockableHeadTrackable{heads: heads, id: 1}

	ht := services.NewHeadTracker(store)
	assert.Nil(t, ht.Start())
}

type lockableHeadTrackable struct {
	heads chan int
	id    int
}

func (t *lockableHeadTrackable) Connect(*models.Head) error {
	t.heads <- t.id
	return nil
}

func (t *lockableHeadTrackable) Disconnect()            {}
func (t *lockableHeadTrackable) OnNewHead(*models.Head) {}

func TestHeadTracker_AttachAlwaysRespectsConnectOrder(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	cltest.MockEthOnStore(store)

	heads := make(chan int)
	firstHeadTrackable := &lockableHeadTrackable{heads: heads, id: 1}
	secondHeadTrackable := &lockableHeadTrackable{heads: heads, id: 2}
	thirdHeadTrackable := &lockableHeadTrackable{heads: heads, id: 3}

	ht := services.NewHeadTracker(store)
	assert.Nil(t, ht.Start())

	go func() {
		ht.Attach(firstHeadTrackable)
		ht.Attach(secondHeadTrackable)
		ht.Attach(thirdHeadTrackable)
	}()

	logger.Debug("waiting for firstHead")
	firstHead := <-heads
	require.Equal(t, 1, firstHead)

	logger.Debug("waiting for secondHead")
	secondHead := <-heads
	require.Equal(t, 2, secondHead)

	logger.Debug("waiting for thirdHead")
	thirdHead := <-heads
	require.Equal(t, 3, thirdHead)
}
