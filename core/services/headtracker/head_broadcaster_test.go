package headtracker_test

import (
	"testing"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/services/headtracker"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHeadBroadcaster_Subscribe(t *testing.T) {
	t.Parallel()
	g := gomega.NewGomegaWithT(t)

	cfg := cltest.NewTestGeneralConfig(t)
	evmCfg := evmtest.NewChainScopedConfig(t, cfg)
	store, cleanup := cltest.NewStoreWithConfig(t, cfg)
	defer cleanup()
	logger := store.Config.CreateProductionLogger()

	sub := new(mocks.Subscription)
	ethClient := cltest.NewEthClientMockWithDefaultChain(t)

	chchHeaders := make(chan chan<- *models.Head, 1)
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

	hr := headtracker.NewHeadBroadcaster(logger)
	orm := headtracker.NewORM(store.DB, *ethClient.ChainID())
	ht := headtracker.NewHeadTracker(logger, ethClient, evmCfg, orm, hr, cltest.NeverSleeper{})
	require.NoError(t, hr.Start())
	defer hr.Close()
	require.NoError(t, ht.Start())
	defer ht.Stop()

	latest1, unsubscribe1 := hr.Subscribe(checker1)
	// "latest head" is nil here because we didn't receive any yet
	assert.Equal(t, (*models.Head)(nil), latest1)

	headers := <-chchHeaders
	h := models.Head{Number: 1}
	headers <- &h
	g.Eventually(func() int32 { return checker1.OnNewLongestChainCount() }).Should(gomega.Equal(int32(1)))

	latest2, _ := hr.Subscribe(checker2)
	// "latest head" is set here to the most recent head received
	assert.NotNil(t, latest2)
	assert.Equal(t, h.Number, latest2.Number)

	unsubscribe1()

	headers <- &models.Head{Number: 2}
	g.Eventually(func() int32 { return checker2.OnNewLongestChainCount() }).Should(gomega.Equal(int32(1)))

	require.NoError(t, ht.Stop())
}
