package headtracker_test

import (
	"testing"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/headtracker"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHeadBroadcaster_Subscribe(t *testing.T) {
	t.Parallel()
	g := gomega.NewGomegaWithT(t)

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	logger := store.Config.CreateProductionLogger()

	sub := new(mocks.Subscription)
	ethClient := new(mocks.Client)
	store.EthClient = ethClient

	chchHeaders := make(chan chan<- *models.Head, 1)
	ethClient.On("ChainID", mock.Anything).Return(store.Config.ChainID(), nil)
	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			chchHeaders <- args.Get(1).(chan<- *models.Head)
		}).
		Return(sub, nil)
	ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(cltest.Head(1), nil)

	sub.On("Unsubscribe").Return()
	sub.On("Err").Return(nil)

	checker1 := &cltest.MockHeadTrackable{}
	checker2 := &cltest.MockHeadTrackable{}

	hr := headtracker.NewHeadBroadcaster()
	ht := headtracker.NewHeadTracker(logger, store, hr, cltest.NeverSleeper{})
	require.NoError(t, hr.Start())
	defer hr.Close()
	require.NoError(t, ht.Start())
	defer ht.Stop()

	hr.Subscribe(checker1)
	unsubscribe2 := hr.Subscribe(checker2)

	headers := <-chchHeaders
	headers <- &models.Head{Number: 1}
	g.Eventually(func() int32 { return checker1.OnNewLongestChainCount() }).Should(gomega.Equal(int32(1)))
	g.Eventually(func() int32 { return checker2.OnNewLongestChainCount() }).Should(gomega.Equal(int32(1)))

	unsubscribe2()

	headers <- &models.Head{Number: 2}
	g.Eventually(func() int32 { return checker1.OnNewLongestChainCount() }).Should(gomega.Equal(int32(2)))
	g.Eventually(func() int32 { return checker2.OnNewLongestChainCount() }).Should(gomega.Equal(int32(1)))

	require.NoError(t, ht.Stop())
}
