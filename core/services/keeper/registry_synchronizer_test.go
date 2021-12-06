package keeper_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	ethmocks "github.com/smartcontractkit/chainlink/core/services/eth/mocks"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/services/log"
	logmocks "github.com/smartcontractkit/chainlink/core/services/log/mocks"
)

const syncInterval = 1000 * time.Hour // prevents sync timer from triggering during test
const syncUpkeepQueueSize = 10

var registryConfig = keeper_registry_wrapper.GetConfig{
	PaymentPremiumPPB: 100,
	BlockCountPerTurn: big.NewInt(20),
	CheckGasLimit:     2_000_000,
	StalenessSeconds:  big.NewInt(3600),
	FallbackGasPrice:  big.NewInt(1000000),
	FallbackLinkPrice: big.NewInt(1000000),
}

var upkeepConfig = keeper_registry_wrapper.GetUpkeep{
	Target:              cltest.NewAddress(),
	ExecuteGas:          2_000_000,
	CheckData:           common.Hex2Bytes("1234"),
	Balance:             big.NewInt(1000000000000000000),
	LastKeeper:          cltest.NewAddress(),
	Admin:               cltest.NewAddress(),
	MaxValidBlocknumber: 1_000_000_000,
}

func setupRegistrySync(t *testing.T) (
	*sqlx.DB,
	*keeper.RegistrySynchronizer,
	*ethmocks.Client,
	*logmocks.Broadcaster,
	job.Job,
) {
	db := pgtest.NewSqlxDB(t)
	korm := keeper.NewORM(db, logger.TestLogger(t), nil, nil, nil)
	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	lbMock := new(logmocks.Broadcaster)
	lbMock.Test(t)
	lbMock.On("AddDependents", 1).Maybe()
	j := cltest.MustInsertKeeperJob(t, db, korm, cltest.NewEIP55Address(), cltest.NewEIP55Address())
	cfg := cltest.NewTestGeneralConfig(t)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, Client: ethClient, LogBroadcaster: lbMock, GeneralConfig: cfg})
	ch := evmtest.MustGetDefaultChain(t, cc)
	keyStore := cltest.NewKeyStore(t, db, cfg)
	jpv2 := cltest.NewJobPipelineV2(t, cfg, cc, db, keyStore)
	contractAddress := j.KeeperSpec.ContractAddress.Address()
	contract, err := keeper_registry_wrapper.NewKeeperRegistry(
		contractAddress,
		ethClient,
	)
	require.NoError(t, err)

	lbMock.On("Register", mock.Anything, mock.MatchedBy(func(opts log.ListenerOpts) bool {
		return opts.Contract == contractAddress
	})).Return(func() {})
	lbMock.On("IsConnected").Return(true).Maybe()

	orm := keeper.NewORM(db, logger.TestLogger(t), nil, ch.Config(), bulletprooftxmanager.SendEveryStrategy{})
	synchronizer := keeper.NewRegistrySynchronizer(keeper.RegistrySynchronizerOptions{
		Job:                      j,
		Contract:                 contract,
		ORM:                      orm,
		JRM:                      jpv2.Jrm,
		LogBroadcaster:           lbMock,
		SyncInterval:             syncInterval,
		MinIncomingConfirmations: 1,
		Logger:                   logger.TestLogger(t),
		SyncUpkeepQueueSize:      syncUpkeepQueueSize,
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

func Test_RegistrySynchronizer_Start(t *testing.T) {
	db, synchronizer, ethMock, _, job := setupRegistrySync(t)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.RegistryABI, contractAddress)
	canceledUpkeeps := []*big.Int{big.NewInt(1)}
	registryMock.MockResponse("getConfig", registryConfig).Once()
	registryMock.MockResponse("getKeeperList", []common.Address{fromAddress}).Once()
	registryMock.MockResponse("getCanceledUpkeepList", canceledUpkeeps).Once()
	registryMock.MockResponse("getUpkeepCount", big.NewInt(0)).Once()

	err := synchronizer.Start()
	require.NoError(t, err)
	defer synchronizer.Close()

	cltest.WaitForCount(t, db, "keeper_registries", 1)

	err = synchronizer.Start()
	require.Error(t, err)
}

func Test_RegistrySynchronizer_CalcPositioningConstant(t *testing.T) {
	t.Parallel()
	for _, upkeepID := range []int64{0, 1, 100, 10_000} {
		_, err := keeper.CalcPositioningConstant(upkeepID, cltest.NewEIP55Address())
		require.NoError(t, err)
	}
}

func Test_RegistrySynchronizer_FullSync(t *testing.T) {
	db, synchronizer, ethMock, _, job := setupRegistrySync(t)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.RegistryABI, contractAddress)
	canceledUpkeeps := []*big.Int{big.NewInt(1)}
	registryMock.MockResponse("getConfig", registryConfig).Once()
	registryMock.MockResponse("getKeeperList", []common.Address{fromAddress}).Once()
	registryMock.MockResponse("getCanceledUpkeepList", canceledUpkeeps).Once()
	registryMock.MockResponse("getUpkeepCount", big.NewInt(3)).Once()
	registryMock.MockResponse("getUpkeep", upkeepConfig).Times(3) // sync all 3, then delete

	synchronizer.ExportedFullSync()

	cltest.AssertCount(t, db, "keeper_registries", 1)
	cltest.AssertCount(t, db, "upkeep_registrations", 2)

	var registry keeper.Registry
	var upkeepRegistration keeper.UpkeepRegistration
	require.NoError(t, db.Get(&registry, `SELECT * FROM keeper_registries`))
	require.NoError(t, db.Get(&upkeepRegistration, `SELECT * FROM upkeep_registrations`))
	require.Equal(t, job.KeeperSpec.ContractAddress, registry.ContractAddress)
	require.Equal(t, job.KeeperSpec.FromAddress, registry.FromAddress)
	require.Equal(t, int32(20), registry.BlockCountPerTurn)
	require.Equal(t, int32(0), registry.KeeperIndex)
	require.Equal(t, int32(1), registry.NumKeepers)
	require.Equal(t, upkeepConfig.CheckData, upkeepRegistration.CheckData)
	require.Equal(t, uint64(upkeepConfig.ExecuteGas), upkeepRegistration.ExecuteGas)

	assertUpkeepIDs(t, db, []int64{0, 2})
	ethMock.AssertExpectations(t)

	// 2nd sync
	canceledUpkeeps = []*big.Int{big.NewInt(0), big.NewInt(1), big.NewInt(3)}
	registryMock.MockResponse("getConfig", registryConfig).Once()
	registryMock.MockResponse("getKeeperList", []common.Address{fromAddress}).Once()
	registryMock.MockResponse("getCanceledUpkeepList", canceledUpkeeps).Once()
	registryMock.MockResponse("getUpkeepCount", big.NewInt(5)).Once()
	registryMock.MockResponse("getUpkeep", upkeepConfig).Times(2) // two new upkeeps to sync

	synchronizer.ExportedFullSync()

	cltest.AssertCount(t, db, "keeper_registries", 1)
	cltest.AssertCount(t, db, "upkeep_registrations", 2)
	assertUpkeepIDs(t, db, []int64{2, 4})
	ethMock.AssertExpectations(t)
}

func Test_RegistrySynchronizer_ConfigSetLog(t *testing.T) {
	db, synchronizer, ethMock, lb, job := setupRegistrySync(t)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.RegistryABI, contractAddress)
	registryMock.MockResponse("getKeeperList", []common.Address{fromAddress}).Once()
	registryMock.MockResponse("getConfig", registryConfig).Once()
	registryMock.MockResponse("getCanceledUpkeepList", []*big.Int{}).Once()
	registryMock.MockResponse("getUpkeepCount", big.NewInt(0)).Once()

	require.NoError(t, synchronizer.Start())
	defer synchronizer.Close()
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	var registry keeper.Registry
	require.NoError(t, db.Get(&registry, `SELECT * FROM keeper_registries`))

	registryConfig.BlockCountPerTurn = big.NewInt(40) // change from default
	registryMock.MockResponse("getKeeperList", []common.Address{fromAddress}).Once()
	registryMock.MockResponse("getConfig", registryConfig).Once()

	cfg := cltest.NewTestGeneralConfig(t)
	head := cltest.MustInsertHead(t, db, cfg, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := keeper_registry_wrapper.KeeperRegistryConfigSet{}
	logBroadcast := new(logmocks.Broadcast)
	logBroadcast.On("DecodedLog").Return(&log)
	logBroadcast.On("RawLog").Return(rawLog)
	logBroadcast.On("String").Maybe().Return("")
	lb.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

	// Do the thing
	synchronizer.HandleLog(logBroadcast)

	cltest.AssertRecordEventually(t, db, &registry, fmt.Sprintf(`SELECT * FROM keeper_registries WHERE id = %d`, registry.ID), func() bool {
		return registry.BlockCountPerTurn == 40
	})
	cltest.AssertCount(t, db, "keeper_registries", 1)
	ethMock.AssertExpectations(t)
	logBroadcast.AssertExpectations(t)
}

func Test_RegistrySynchronizer_KeepersUpdatedLog(t *testing.T) {
	db, synchronizer, ethMock, lb, job := setupRegistrySync(t)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.RegistryABI, contractAddress)
	registryMock.MockResponse("getKeeperList", []common.Address{fromAddress}).Once()
	registryMock.MockResponse("getConfig", registryConfig).Once()
	registryMock.MockResponse("getCanceledUpkeepList", []*big.Int{}).Once()
	registryMock.MockResponse("getUpkeepCount", big.NewInt(0)).Once()

	require.NoError(t, synchronizer.Start())
	defer synchronizer.Close()
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	var registry keeper.Registry
	require.NoError(t, db.Get(&registry, `SELECT * FROM keeper_registries`))

	addresses := []common.Address{fromAddress, cltest.NewAddress()} // change from default
	registryMock.MockResponse("getConfig", registryConfig).Once()
	registryMock.MockResponse("getKeeperList", addresses).Once()

	cfg := cltest.NewTestGeneralConfig(t)
	head := cltest.MustInsertHead(t, db, cfg, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := keeper_registry_wrapper.KeeperRegistryKeepersUpdated{}
	logBroadcast := new(logmocks.Broadcast)
	logBroadcast.On("DecodedLog").Return(&log)
	logBroadcast.On("RawLog").Return(rawLog)
	logBroadcast.On("String").Maybe().Return("")
	lb.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

	// Do the thing
	synchronizer.HandleLog(logBroadcast)

	cltest.AssertRecordEventually(t, db, &registry, fmt.Sprintf(`SELECT * FROM keeper_registries WHERE id = %d`, registry.ID), func() bool {
		return registry.NumKeepers == 2
	})
	cltest.AssertCount(t, db, "keeper_registries", 1)
	ethMock.AssertExpectations(t)
	logBroadcast.AssertExpectations(t)
}

func Test_RegistrySynchronizer_UpkeepCanceledLog(t *testing.T) {
	db, synchronizer, ethMock, lb, job := setupRegistrySync(t)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.RegistryABI, contractAddress)
	registryMock.MockResponse("getConfig", registryConfig).Once()
	registryMock.MockResponse("getKeeperList", []common.Address{fromAddress}).Once()
	registryMock.MockResponse("getCanceledUpkeepList", []*big.Int{}).Once()
	registryMock.MockResponse("getUpkeepCount", big.NewInt(3)).Once()
	registryMock.MockResponse("getUpkeep", upkeepConfig).Times(3)

	require.NoError(t, synchronizer.Start())
	defer func() { require.NoError(t, synchronizer.Close()) }()
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	cltest.WaitForCount(t, db, "upkeep_registrations", 3)

	cfg := cltest.NewTestGeneralConfig(t)
	head := cltest.MustInsertHead(t, db, cfg, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := keeper_registry_wrapper.KeeperRegistryUpkeepCanceled{Id: big.NewInt(1)}
	logBroadcast := new(logmocks.Broadcast)
	logBroadcast.On("DecodedLog").Return(&log)
	logBroadcast.On("RawLog").Return(rawLog)
	logBroadcast.On("String").Maybe().Return("")
	lb.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

	// Do the thing
	synchronizer.HandleLog(logBroadcast)

	cltest.WaitForCount(t, db, "upkeep_registrations", 2)
	ethMock.AssertExpectations(t)
	logBroadcast.AssertExpectations(t)
}

func Test_RegistrySynchronizer_UpkeepRegisteredLog(t *testing.T) {
	db, synchronizer, ethMock, lb, job := setupRegistrySync(t)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.RegistryABI, contractAddress)
	registryMock.MockResponse("getConfig", registryConfig).Once()
	registryMock.MockResponse("getKeeperList", []common.Address{fromAddress}).Once()
	registryMock.MockResponse("getCanceledUpkeepList", []*big.Int{}).Once()
	registryMock.MockResponse("getUpkeepCount", big.NewInt(0)).Once()

	require.NoError(t, synchronizer.Start())
	defer synchronizer.Close()
	cltest.WaitForCount(t, db, "keeper_registries", 1)

	registryMock.MockResponse("getUpkeep", upkeepConfig).Once()

	cfg := cltest.NewTestGeneralConfig(t)
	head := cltest.MustInsertHead(t, db, cfg, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := keeper_registry_wrapper.KeeperRegistryUpkeepRegistered{Id: big.NewInt(3)}
	logBroadcast := new(logmocks.Broadcast)
	logBroadcast.On("DecodedLog").Return(&log)
	logBroadcast.On("RawLog").Return(rawLog)
	logBroadcast.On("String").Maybe().Return("")
	lb.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

	// Do the thing
	synchronizer.HandleLog(logBroadcast)

	cltest.WaitForCount(t, db, "upkeep_registrations", 1)
	ethMock.AssertExpectations(t)
	logBroadcast.AssertExpectations(t)
}

func Test_RegistrySynchronizer_UpkeepPerformedLog(t *testing.T) {
	g := gomega.NewWithT(t)

	db, synchronizer, ethMock, lb, job := setupRegistrySync(t)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.RegistryABI, contractAddress)
	registryMock.MockResponse("getConfig", registryConfig).Once()
	registryMock.MockResponse("getKeeperList", []common.Address{fromAddress}).Once()
	registryMock.MockResponse("getCanceledUpkeepList", []*big.Int{}).Once()
	registryMock.MockResponse("getUpkeepCount", big.NewInt(1)).Once()
	registryMock.MockResponse("getUpkeep", upkeepConfig).Once()

	require.NoError(t, synchronizer.Start())
	defer synchronizer.Close()
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	cltest.WaitForCount(t, db, "upkeep_registrations", 1)

	pgtest.MustExec(t, db, `UPDATE upkeep_registrations SET last_run_block_height = 100`)

	cfg := cltest.NewTestGeneralConfig(t)
	head := cltest.MustInsertHead(t, db, cfg, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := keeper_registry_wrapper.KeeperRegistryUpkeepPerformed{Id: big.NewInt(0)}
	logBroadcast := new(logmocks.Broadcast)
	logBroadcast.On("DecodedLog").Return(&log)
	logBroadcast.On("RawLog").Return(rawLog)
	logBroadcast.On("String").Maybe().Return("")
	lb.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

	// Do the thing
	synchronizer.HandleLog(logBroadcast)

	g.Eventually(func() int64 {
		var upkeep keeper.UpkeepRegistration
		err := db.Get(&upkeep, `SELECT * FROM upkeep_registrations`)
		require.NoError(t, err)
		return upkeep.LastRunBlockHeight
	}, cltest.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal(int64(0)))

	ethMock.AssertExpectations(t)
	logBroadcast.AssertExpectations(t)
}
