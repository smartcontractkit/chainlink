package blockhashstore_test

import (
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ethkey"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	lpmocks "github.com/smartcontractkit/chainlink/v2/common/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockhashstore"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
)

func TestDelegate_JobType(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	delegate := blockhashstore.NewDelegate(nil, lggr, nil, nil)

	assert.Equal(t, job.BlockhashStore, delegate.JobType())
}

type testData struct {
	ethKeyStore  keystore.Eth
	legacyChains legacyevm.LegacyChainContainer
	sendingKey   ethkey.KeyV2
	logs         *observer.ObservedLogs
}

func createTestDelegate(t *testing.T) (*blockhashstore.Delegate, *testData) {
	t.Helper()

	lggr, logs := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	ethClient := client.NewNullClient(big.NewInt(evmtest.NullClientChainID), logger.TestLogger(t))
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Feature.LogPoller = func(b bool) *bool { return &b }(true)
	})
	db := pgtest.NewSqlxDB(t)
	kst := cltest.NewKeyStore(t, db).Eth()
	sendingKey, _ := cltest.MustInsertRandomKey(t, kst)
	lp := lpmocks.NewLogPoller(t)
	servicetest.SetupNoOpMock(lp)
	lp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil).Maybe()
	lp.On("LatestBlock", mock.Anything).Return(logpoller.Block{}, nil).Maybe()

	legacyChains := evmtest.NewLegacyChains(
		t,
		evmtest.TestChainOpts{
			ChainConfigs:   cfg.EVMConfigs(),
			DatabaseConfig: cfg.Database(),
			FeatureConfig:  cfg.Feature(),
			ListenerConfig: cfg.Database().Listener(),
			DB:             db,
			KeyStore:       kst,
			Client:         ethClient,
			LogPoller:      lp,
		},
	)
	return blockhashstore.NewDelegate(cfg, lggr, legacyChains, kst), &testData{
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
	chain, ok := testData.legacyChains.Slice()[0].(legacyevm.Chain)
	require.True(t, ok)
	finalityDepth := chain.Config().EVM().FinalityDepth()
	if finalityDepth > math.MaxInt32 {
		t.Fatalf("finality depth overflows int32: %d", finalityDepth)
	}
	defaultWaitBlocks := (int32)(finalityDepth)

	t.Run("happy", func(t *testing.T) {
		spec := job.Job{BlockhashStoreSpec: &job.BlockhashStoreSpec{WaitBlocks: defaultWaitBlocks, EVMChainID: (*sqlutil.Big)(testutils.FixtureChainID)}}
		services, err := delegate.ServicesForSpec(testutils.Context(t), spec)

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
			EVMChainID:               (*sqlutil.Big)(testutils.FixtureChainID),
		}}
		services, err := delegate.ServicesForSpec(testutils.Context(t), spec)

		require.NoError(t, err)
		require.Len(t, services, 1)
	})

	t.Run("missing BlockhashStoreSpec", func(t *testing.T) {
		spec := job.Job{BlockhashStoreSpec: nil}
		_, err := delegate.ServicesForSpec(testutils.Context(t), spec)
		assert.Error(t, err)
	})

	t.Run("wrong EVMChainID", func(t *testing.T) {
		spec := job.Job{BlockhashStoreSpec: &job.BlockhashStoreSpec{
			EVMChainID: sqlutil.NewI(123),
		}}
		_, err := delegate.ServicesForSpec(testutils.Context(t), spec)
		assert.Error(t, err)
	})

	t.Run("missing EnabledKeysForChain", func(t *testing.T) {
		ctx := testutils.Context(t)
		_, err := testData.ethKeyStore.Delete(ctx, testData.sendingKey.ID())
		require.NoError(t, err)

		spec := job.Job{BlockhashStoreSpec: &job.BlockhashStoreSpec{
			WaitBlocks: defaultWaitBlocks,
		}}
		_, err = delegate.ServicesForSpec(testutils.Context(t), spec)
		assert.Error(t, err)
	})
}

func TestDelegate_StartStop(t *testing.T) {
	t.Parallel()

	delegate, testData := createTestDelegate(t)

	require.NotEmpty(t, testData.legacyChains.Slice())
	chain, ok := testData.legacyChains.Slice()[0].(legacyevm.Chain)
	require.True(t, ok)

	finalityDepth := chain.Config().EVM().FinalityDepth()
	if finalityDepth > math.MaxInt32 {
		t.Fatalf("finality depth overflows int32: %d", finalityDepth)
	}
	defaultWaitBlocks := (int32)(finalityDepth)
	spec := job.Job{BlockhashStoreSpec: &job.BlockhashStoreSpec{
		WaitBlocks: defaultWaitBlocks,
		PollPeriod: time.Second,
		RunTimeout: testutils.WaitTimeout(t),
		EVMChainID: (*sqlutil.Big)(testutils.FixtureChainID),
	}}
	services, err := delegate.ServicesForSpec(testutils.Context(t), spec)

	require.NoError(t, err)
	require.Len(t, services, 1)

	err = services[0].Start(testutils.Context(t))
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		return testData.logs.FilterMessage("Starting BHS feeder").Len() > 0 &&
			testData.logs.FilterMessage("Running BHS feeder").Len() > 0 &&
			testData.logs.FilterMessage("BHS feeder run completed successfully").Len() > 0
	}, testutils.WaitTimeout(t), testutils.TestInterval)

	err = services[0].Close()
	require.NoError(t, err)

	assert.NotZero(t, testData.logs.FilterMessage("Stopping BHS feeder").Len())
}
