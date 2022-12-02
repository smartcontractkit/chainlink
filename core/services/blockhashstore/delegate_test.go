package blockhashstore_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/blockhashstore"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestDelegate_JobType(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	delegate := blockhashstore.NewDelegate(lggr, nil, nil)

	assert.Equal(t, job.BlockhashStore, delegate.JobType())
}

type testData struct {
	ethClient   *mocks.Client
	ethKeyStore keystore.Eth
	chainSet    evm.ChainSet
	sendingKey  ethkey.KeyV2
	logs        *observer.ObservedLogs
}

func createTestDelegate(t *testing.T) (*blockhashstore.Delegate, *testData) {
	t.Helper()

	lggr, logs := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	db := pgtest.NewSqlxDB(t)
	kst := cltest.NewKeyStore(t, db, cfg).Eth()
	sendingKey, _ := cltest.MustAddRandomKeyToKeystore(t, kst)
	chainSet := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, KeyStore: kst, GeneralConfig: cfg, Client: ethClient})

	return blockhashstore.NewDelegate(lggr, chainSet, kst), &testData{
		ethClient:   ethClient,
		ethKeyStore: kst,
		chainSet:    chainSet,
		sendingKey:  sendingKey,
		logs:        logs,
	}
}

func TestDelegate_ServicesForSpec(t *testing.T) {
	t.Parallel()

	delegate, testData := createTestDelegate(t)

	require.NotEmpty(t, testData.chainSet.Chains())
	defaultWaitBlocks := (int32)(testData.chainSet.Chains()[0].Config().EvmFinalityDepth())

	t.Run("happy", func(t *testing.T) {
		spec := job.Job{BlockhashStoreSpec: &job.BlockhashStoreSpec{WaitBlocks: defaultWaitBlocks}}
		services, err := delegate.ServicesForSpec(spec)

		require.NoError(t, err)
		require.Len(t, services, 1)
	})

	t.Run("happy with coordinators", func(t *testing.T) {
		coordinatorV1 := cltest.NewEIP55Address()
		coordinatorV2 := cltest.NewEIP55Address()

		spec := job.Job{BlockhashStoreSpec: &job.BlockhashStoreSpec{
			WaitBlocks:           defaultWaitBlocks,
			CoordinatorV1Address: &coordinatorV1,
			CoordinatorV2Address: &coordinatorV2,
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

	t.Run("WaitBlocks less than EvmFinalityDepth", func(t *testing.T) {
		spec := job.Job{BlockhashStoreSpec: &job.BlockhashStoreSpec{
			WaitBlocks: defaultWaitBlocks - 1,
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

	require.NotEmpty(t, testData.chainSet.Chains())
	defaultWaitBlocks := (int32)(testData.chainSet.Chains()[0].Config().EvmFinalityDepth())
	spec := job.Job{BlockhashStoreSpec: &job.BlockhashStoreSpec{
		WaitBlocks: defaultWaitBlocks,
		PollPeriod: time.Second,
		RunTimeout: testutils.WaitTimeout(t),
	}}
	services, err := delegate.ServicesForSpec(spec)

	require.NoError(t, err)
	require.Len(t, services, 1)

	blocks := cltest.NewBlocks(t, 1)
	testData.ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(blocks.Head(0), nil)
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
