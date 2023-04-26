package keeper_test

import (
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	logmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/log/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keeper"
	"github.com/smartcontractkit/chainlink/v2/core/services/srvctest"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const syncInterval = 1000 * time.Hour // prevents sync timer from triggering during test
const syncUpkeepQueueSize = 10

func setupRegistrySync(t *testing.T, version keeper.RegistryVersion) (
	*sqlx.DB,
	*keeper.RegistrySynchronizer,
	*evmclimocks.Client,
	*logmocks.Broadcaster,
	job.Job,
) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	scopedConfig := evmtest.NewChainScopedConfig(t, cfg)
	korm := keeper.NewORM(db, logger.TestLogger(t), scopedConfig)
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	keyStore := cltest.NewKeyStore(t, db, cfg)
	lbMock := logmocks.NewBroadcaster(t)
	lbMock.On("AddDependents", 1).Maybe()
	j := cltest.MustInsertKeeperJob(t, db, korm, cltest.NewEIP55Address(), cltest.NewEIP55Address())
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, Client: ethClient, LogBroadcaster: lbMock, GeneralConfig: cfg, KeyStore: keyStore.Eth()})
	ch := evmtest.MustGetDefaultChain(t, cc)
	jpv2 := cltest.NewJobPipelineV2(t, cfg, cc, db, keyStore, nil, nil)
	contractAddress := j.KeeperSpec.ContractAddress.Address()

	switch version {
	case keeper.RegistryVersion_1_0, keeper.RegistryVersion_1_1:
		registryMock := cltest.NewContractMockReceiver(t, ethClient, keeper.Registry1_1ABI, contractAddress)
		registryMock.MockResponse("typeAndVersion", "KeeperRegistry 1.1.1").Once()
	case keeper.RegistryVersion_1_2:
		registryMock := cltest.NewContractMockReceiver(t, ethClient, keeper.Registry1_2ABI, contractAddress)
		registryMock.MockResponse("typeAndVersion", "KeeperRegistry 1.2.0").Once()
	case keeper.RegistryVersion_1_3:
		registryMock := cltest.NewContractMockReceiver(t, ethClient, keeper.Registry1_3ABI, contractAddress)
		registryMock.MockResponse("typeAndVersion", "KeeperRegistry 1.3.0").Once()
	}

	registryWrapper, err := keeper.NewRegistryWrapper(j.KeeperSpec.ContractAddress, ethClient)
	require.NoError(t, err)

	lbMock.On("Register", mock.Anything, mock.MatchedBy(func(opts log.ListenerOpts) bool {
		return opts.Contract == contractAddress
	})).Maybe().Return(func() {})
	lbMock.On("IsConnected").Return(true).Maybe()

	mailMon := srvctest.Start(t, utils.NewMailboxMonitor(t.Name()))

	orm := keeper.NewORM(db, logger.TestLogger(t), ch.Config())
	synchronizer := keeper.NewRegistrySynchronizer(keeper.RegistrySynchronizerOptions{
		Job:                      j,
		RegistryWrapper:          *registryWrapper,
		ORM:                      orm,
		JRM:                      jpv2.Jrm,
		LogBroadcaster:           lbMock,
		MailMon:                  mailMon,
		SyncInterval:             syncInterval,
		MinIncomingConfirmations: 1,
		Logger:                   logger.TestLogger(t),
		SyncUpkeepQueueSize:      syncUpkeepQueueSize,
		EffectiveKeeperAddress:   j.KeeperSpec.FromAddress.Address(),
	})
	return db, synchronizer, ethClient, lbMock, j
}

func assertUpkeepIDs(t *testing.T, db *sqlx.DB, expected []int64) {
	g := gomega.NewWithT(t)
	var upkeepIDs []int64
	err := db.Select(&upkeepIDs, `SELECT upkeep_id FROM upkeep_registrations`)
	require.NoError(t, err)
	require.Equal(t, len(expected), len(upkeepIDs))
	g.Expect(upkeepIDs).To(gomega.ContainElements(expected))
}
