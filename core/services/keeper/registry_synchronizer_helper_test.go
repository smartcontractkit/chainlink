package keeper_test

import (
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/jmoiron/sqlx"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox/mailboxtest"

	"github.com/smartcontractkit/chainlink-integrations/evm/client/clienttest"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	logmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/log/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keeper"
)

const syncInterval = 1000 * time.Hour // prevents sync timer from triggering during test
const syncUpkeepQueueSize = 10

func setupRegistrySync(t *testing.T, version keeper.RegistryVersion) (
	*sqlx.DB,
	*keeper.RegistrySynchronizer,
	*clienttest.Client,
	*logmocks.Broadcaster,
	job.Job,
) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	korm := keeper.NewORM(db, logger.TestLogger(t))
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	keyStore := cltest.NewKeyStore(t, db)
	lbMock := logmocks.NewBroadcaster(t)
	lbMock.On("AddDependents", 1).Maybe()
	j := cltest.MustInsertKeeperJob(t, db, korm, cltest.NewEIP55Address(), cltest.NewEIP55Address())
	legacyChains := evmtest.NewLegacyChains(t, evmtest.TestChainOpts{
		DB:             db,
		Client:         ethClient,
		LogBroadcaster: lbMock,
		ChainConfigs:   cfg.EVMConfigs(),
		DatabaseConfig: cfg.Database(),
		FeatureConfig:  cfg.Feature(),
		ListenerConfig: cfg.Database().Listener(),
		KeyStore:       keyStore.Eth(),
	})
	jpv2 := cltest.NewJobPipelineV2(t, cfg.WebServer(), cfg.JobPipeline(), legacyChains, db, keyStore, nil, nil)
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
	case keeper.RegistryVersion_2_0, keeper.RegistryVersion_2_1:
		t.Fatalf("Unsupported version: %s", version)
	}

	registryWrapper, err := keeper.NewRegistryWrapper(j.KeeperSpec.ContractAddress, ethClient)
	require.NoError(t, err)

	lbMock.On("Register", mock.Anything, mock.MatchedBy(func(opts log.ListenerOpts) bool {
		return opts.Contract == contractAddress
	})).Maybe().Return(func() {})
	lbMock.On("IsConnected").Return(true).Maybe()

	mailMon := servicetest.Run(t, mailboxtest.NewMonitor(t))

	orm := keeper.NewORM(db, logger.TestLogger(t))
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
