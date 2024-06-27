package blockhashstore_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	mocklp "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockhashstore"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	evmrelay "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestDelegate_JobType(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	delegate := blockhashstore.NewDelegate(lggr, nil, nil)

	assert.Equal(t, job.BlockhashStore, delegate.JobType())
}

type testData struct {
	ethClient    *mocks.Client
	ethKeyStore  keystore.Eth
	legacyChains evm.LegacyChainContainer
	sendingKey   ethkey.KeyV2
	logs         *observer.ObservedLogs
}

func createTestDelegate(t *testing.T) (*blockhashstore.Delegate, *testData) {
	t.Helper()

	lggr, logs := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Feature.LogPoller = func(b bool) *bool { return &b }(true)
	})
	db := pgtest.NewSqlxDB(t)
	kst := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	sendingKey, _ := cltest.MustInsertRandomKey(t, kst)
	lp := &mocklp.LogPoller{}
	lp.On("RegisterFilter", mock.Anything).Return(nil)
	lp.On("LatestBlock", mock.Anything, mock.Anything).Return(int64(0), nil)

	relayExtenders := evmtest.NewChainRelayExtenders(
		t,
		evmtest.TestChainOpts{
			DB:            db,
			KeyStore:      kst,
			GeneralConfig: cfg,
			Client:        ethClient,
			LogPoller:     lp,
		},
	)
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)
	return blockhashstore.NewDelegate(lggr, legacyChains, kst), &testData{
		ethClient:    ethClient,
		ethKeyStore:  kst,
		legacyChains: legacyChains,
		sendingKey:   sendingKey,
		logs:         logs,
	}
}

func TestDelegate_ServicesForSpec(t *testing.T) {
	t.Parallel()

	delegate, testData := createTestDelegate(t)

	require.NotEmpty(t, testData.legacyChains.Slice())
	defaultWaitBlocks := (int32)(testData.legacyChains.Slice()[0].Config().EVM().FinalityDepth())

	t.Run("happy", func(t *testing.T) {
		spec := job.Job{BlockhashStoreSpec: &job.BlockhashStoreSpec{WaitBlocks: defaultWaitBlocks, EVMChainID: (*utils.Big)(testutils.FixtureChainID)}}
		services, err := delegate.ServicesForSpec(spec)

		require.NoError(t, err)
		require.Len(t, services, 1)
	})

	t.Run("happy with coordinators", func(t *testing.T) {
		coordinatorV1 := cltest.NewEIP55Address()
		coordinatorV2 := cltest.NewEIP55Address()
		coordinatorV2Plus := cltest.NewEIP55Address()

		spec := job.Job{BlockhashStoreSpec: &job.BlockhashStoreSpec{
			WaitBlocks:               defaultWaitBlocks,
			CoordinatorV1Address:     &coordinatorV1,
			CoordinatorV2Address:     &coordinatorV2,
			CoordinatorV2PlusAddress: &coordinatorV2Plus,
			EVMChainID:               (*utils.Big)(testutils.FixtureChainID),
		}}
		services, err := delegate.ServicesForSpec(spec)

		require.NoError(t, err)
		require.Len(t, services, 1)
	})

	t.Run("missing BlockhashStoreSpec", func(t *testing.T) {
		spec := job.Job{BlockhashStoreSpec: nil}
		_, err := delegate.ServicesForSpec(spec)
		assert.Error(t, err)
	})

	t.Run("wrong EVMChainID", func(t *testing.T) {
		spec := job.Job{BlockhashStoreSpec: &job.BlockhashStoreSpec{
			EVMChainID: utils.NewBigI(123),
		}}
		_, err := delegate.ServicesForSpec(spec)
		assert.Error(t, err)
	})

	t.Run("missing EnabledKeysForChain", func(t *testing.T) {
		_, err := testData.ethKeyStore.Delete(testData.sendingKey.ID())
		require.NoError(t, err)

		spec := job.Job{BlockhashStoreSpec: &job.BlockhashStoreSpec{
			WaitBlocks: defaultWaitBlocks,
		}}
		_, err = delegate.ServicesForSpec(spec)
		assert.Error(t, err)
	})
}

func TestDelegate_StartStop(t *testing.T) {
	t.Parallel()

	delegate, testData := createTestDelegate(t)

	require.NotEmpty(t, testData.legacyChains.Slice())
	defaultWaitBlocks := (int32)(testData.legacyChains.Slice()[0].Config().EVM().FinalityDepth())
	spec := job.Job{BlockhashStoreSpec: &job.BlockhashStoreSpec{
		WaitBlocks: defaultWaitBlocks,
		PollPeriod: time.Second,
		RunTimeout: testutils.WaitTimeout(t),
		EVMChainID: (*utils.Big)(testutils.FixtureChainID),
	}}
	services, err := delegate.ServicesForSpec(spec)

	require.NoError(t, err)
	require.Len(t, services, 1)

	err = services[0].Start(testutils.Context(t))
	require.NoError(t, err)

	assert.Eventually(t, func() bool {
		return testData.logs.FilterMessage("Starting BHS feeder").Len() > 0 &&
			testData.logs.FilterMessage("Running BHS feeder").Len() > 0 &&
			testData.logs.FilterMessage("BHS feeder run completed successfully").Len() > 0
	}, testutils.WaitTimeout(t), testutils.TestInterval)

	err = services[0].Close()
	require.NoError(t, err)

	assert.NotZero(t, testData.logs.FilterMessage("Stopping BHS feeder").Len())
}
