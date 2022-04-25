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

	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	logmocks "github.com/smartcontractkit/chainlink/core/chains/evm/log/mocks"
	evmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	registry1_1 "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper1_1"
	registry1_2 "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper1_2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
)

const syncInterval = 1000 * time.Hour // prevents sync timer from triggering during test
const syncUpkeepQueueSize = 10

var registryConfig1_1 = registry1_1.GetConfig{
	PaymentPremiumPPB: 100,
	BlockCountPerTurn: big.NewInt(20),
	CheckGasLimit:     2_000_000,
	StalenessSeconds:  big.NewInt(3600),
	FallbackGasPrice:  big.NewInt(1000000),
	FallbackLinkPrice: big.NewInt(1000000),
}

var upkeepConfig1_1 = registry1_1.GetUpkeep{
	Target:              testutils.NewAddress(),
	ExecuteGas:          2_000_000,
	CheckData:           common.Hex2Bytes("1234"),
	Balance:             big.NewInt(1000000000000000000),
	LastKeeper:          testutils.NewAddress(),
	Admin:               testutils.NewAddress(),
	MaxValidBlocknumber: 1_000_000_000,
}

var registryConfig1_2 = registry1_2.Config{
	PaymentPremiumPPB:    100,
	FlatFeeMicroLink:     uint32(0),
	BlockCountPerTurn:    big.NewInt(20),
	CheckGasLimit:        2_000_000,
	StalenessSeconds:     big.NewInt(3600),
	GasCeilingMultiplier: uint16(2),
	MinUpkeepSpend:       big.NewInt(0),
	MaxPerformGas:        uint32(5000000),
	FallbackGasPrice:     big.NewInt(1000000),
	FallbackLinkPrice:    big.NewInt(1000000),
	Transcoder:           cltest.NewEIP55Address().Address(),
	Registrar:            cltest.NewEIP55Address().Address(),
}

var upkeepConfig1_2 = registry1_2.GetUpkeep{
	Target:              testutils.NewAddress(),
	ExecuteGas:          2_000_000,
	CheckData:           common.Hex2Bytes("1234"),
	Balance:             big.NewInt(1000000000000000000),
	LastKeeper:          testutils.NewAddress(),
	Admin:               testutils.NewAddress(),
	MaxValidBlocknumber: 1_000_000_000,
	AmountSpent:         big.NewInt(0),
}

func setupRegistrySync(t *testing.T, version keeper.RegistryVersion) (
	*sqlx.DB,
	*keeper.RegistrySynchronizer,
	*evmmocks.Client,
	*logmocks.Broadcaster,
	job.Job,
) {
	db := pgtest.NewSqlxDB(t)
	korm := keeper.NewORM(db, logger.TestLogger(t), nil, nil)
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

	registryMock := cltest.NewContractMockReceiver(t, ethClient, keeper.Registry1_1ABI, contractAddress)
	switch version {
	case keeper.RegistryVersion_1_0, keeper.RegistryVersion_1_1:
		registryMock.MockResponse("typeAndVersion", "KeeperRegistry 1.1.0").Once()
	case keeper.RegistryVersion_1_2:
		registryMock.MockResponse("typeAndVersion", "KeeperRegistry 1.2.0").Once()
	}

	registryWrapper, err := keeper.NewRegistryWrapper(j.KeeperSpec.ContractAddress, ethClient)
	require.NoError(t, err)

	lbMock.On("Register", mock.Anything, mock.MatchedBy(func(opts log.ListenerOpts) bool {
		return opts.Contract == contractAddress
	})).Return(func() {})
	lbMock.On("IsConnected").Return(true).Maybe()

	orm := keeper.NewORM(db, logger.TestLogger(t), ch.Config(), txmgr.SendEveryStrategy{})
	synchronizer := keeper.NewRegistrySynchronizer(keeper.RegistrySynchronizerOptions{
		Job:                      j,
		RegistryWrapper:          registryWrapper,
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

func mockRegistry1_1(
	t *testing.T,
	ethMock *evmmocks.Client,
	contractAddress common.Address,
	config registry1_1.GetConfig,
	keeperList []common.Address,
	cancelledUpkeeps []*big.Int,
	upkeepCount *big.Int,
	upkeepConfig registry1_1.GetUpkeep,
	timesGetUpkeepMock int,
) {
	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_1ABI, contractAddress)

	registryMock.MockResponse("getConfig", config).Once()
	registryMock.MockResponse("getKeeperList", keeperList).Once()
	registryMock.MockResponse("getCanceledUpkeepList", cancelledUpkeeps).Once()
	registryMock.MockResponse("getUpkeepCount", upkeepCount).Once()
	if timesGetUpkeepMock > 0 {
		registryMock.MockResponse("getUpkeep", upkeepConfig).Times(timesGetUpkeepMock)
	}
}

func mockRegistry1_2(
	t *testing.T,
	ethMock *evmmocks.Client,
	contractAddress common.Address,
	config registry1_2.Config,
	activeUpkeepIDs []*big.Int,
	keeperList []common.Address,
	upkeepConfig registry1_2.GetUpkeep,
	timesGetUpkeepMock int,
) {
	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_2ABI, contractAddress)

	var state = registry1_2.State{
		Nonce:            uint32(0),
		OwnerLinkBalance: big.NewInt(1000000000000000000),
		NumUpkeeps:       big.NewInt(int64(len(activeUpkeepIDs))),
	}
	var getState = registry1_2.GetState{
		State:   state,
		Config:  config,
		Keepers: keeperList,
	}
	registryMock.MockResponse("getState", getState).Once()
	registryMock.MockResponse("getActiveUpkeepIDs", activeUpkeepIDs).Once()
	if timesGetUpkeepMock > 0 {
		registryMock.MockResponse("getUpkeep", upkeepConfig).Times(timesGetUpkeepMock)
	}
}

func Test_RegistrySynchronizer1_1_Start(t *testing.T) {
	db, synchronizer, ethMock, _, job := setupRegistrySync(t, keeper.RegistryVersion_1_1)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()
	mockRegistry1_1(
		t,
		ethMock,
		contractAddress,
		registryConfig1_1,
		[]common.Address{fromAddress},
		[]*big.Int{},
		big.NewInt(0),
		upkeepConfig1_1,
		0)

	err := synchronizer.Start(testutils.Context(t))
	require.NoError(t, err)
	defer synchronizer.Close()

	cltest.WaitForCount(t, db, "keeper_registries", 1)

	err = synchronizer.Start(testutils.Context(t))
	require.Error(t, err)
	ethMock.AssertExpectations(t)
}

func Test_RegistrySynchronizer1_2_Start(t *testing.T) {
	db, synchronizer, ethMock, _, job := setupRegistrySync(t, keeper.RegistryVersion_1_2)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()
	mockRegistry1_2(
		t,
		ethMock,
		contractAddress,
		registryConfig1_2,
		[]*big.Int{},
		[]common.Address{fromAddress},
		upkeepConfig1_2,
		0)

	err := synchronizer.Start(testutils.Context(t))
	require.NoError(t, err)
	defer synchronizer.Close()

	cltest.WaitForCount(t, db, "keeper_registries", 1)

	err = synchronizer.Start(testutils.Context(t))
	require.Error(t, err)
	ethMock.AssertExpectations(t)
}

func Test_RegistrySynchronizer_CalcPositioningConstant(t *testing.T) {
	t.Parallel()
	for _, upkeepID := range []int64{0, 1, 100, 10_000} {
		_, err := keeper.CalcPositioningConstant(upkeepID, cltest.NewEIP55Address())
		require.NoError(t, err)
	}
}

func Test_RegistrySynchronizer1_1_FullSync(t *testing.T) {
	db, synchronizer, ethMock, _, job := setupRegistrySync(t, keeper.RegistryVersion_1_1)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()
	canceledUpkeeps := []*big.Int{big.NewInt(1)}

	mockRegistry1_1(
		t,
		ethMock,
		contractAddress,
		registryConfig1_1,
		[]common.Address{fromAddress},
		canceledUpkeeps,
		big.NewInt(3),
		upkeepConfig1_1,
		3) // sync all 3, then delete

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
	require.Equal(t, upkeepConfig1_1.CheckData, upkeepRegistration.CheckData)
	require.Equal(t, uint64(upkeepConfig1_1.ExecuteGas), upkeepRegistration.ExecuteGas)

	assertUpkeepIDs(t, db, []int64{0, 2})
	ethMock.AssertExpectations(t)

	// 2nd sync
	canceledUpkeeps = []*big.Int{big.NewInt(0), big.NewInt(1), big.NewInt(3)}
	mockRegistry1_1(
		t,
		ethMock,
		contractAddress,
		registryConfig1_1,
		[]common.Address{fromAddress},
		canceledUpkeeps,
		big.NewInt(5),
		upkeepConfig1_1,
		2) // two new upkeeps to sync
	synchronizer.ExportedFullSync()

	cltest.AssertCount(t, db, "keeper_registries", 1)
	cltest.AssertCount(t, db, "upkeep_registrations", 2)
	assertUpkeepIDs(t, db, []int64{2, 4})
	ethMock.AssertExpectations(t)
}

func Test_RegistrySynchronizer1_2_FullSync(t *testing.T) {
	db, synchronizer, ethMock, _, job := setupRegistrySync(t, keeper.RegistryVersion_1_2)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()
	fmt.Println(contractAddress)
	mockRegistry1_2(
		t,
		ethMock,
		contractAddress,
		registryConfig1_2,
		[]*big.Int{big.NewInt(3), big.NewInt(69), big.NewInt(420)}, // Upkeep IDs
		[]common.Address{fromAddress},
		upkeepConfig1_2,
		3) // sync all 3
	synchronizer.ExportedFullSync()

	cltest.AssertCount(t, db, "keeper_registries", 1)
	cltest.AssertCount(t, db, "upkeep_registrations", 3)

	var registry keeper.Registry
	var upkeepRegistration keeper.UpkeepRegistration
	require.NoError(t, db.Get(&registry, `SELECT * FROM keeper_registries`))
	require.Equal(t, job.KeeperSpec.ContractAddress, registry.ContractAddress)
	require.Equal(t, job.KeeperSpec.FromAddress, registry.FromAddress)
	require.Equal(t, int32(20), registry.BlockCountPerTurn)
	require.Equal(t, int32(0), registry.KeeperIndex)
	require.Equal(t, int32(1), registry.NumKeepers)

	require.NoError(t, db.Get(&upkeepRegistration, `SELECT * FROM upkeep_registrations`))
	require.Equal(t, upkeepConfig1_2.CheckData, upkeepRegistration.CheckData)
	require.Equal(t, uint64(upkeepConfig1_2.ExecuteGas), upkeepRegistration.ExecuteGas)

	assertUpkeepIDs(t, db, []int64{3, 69, 420})
	ethMock.AssertExpectations(t)

	// 2nd sync. Cancel upkeep (id 3) and add a new upkeep (id 2022)
	mockRegistry1_2(
		t,
		ethMock,
		contractAddress,
		registryConfig1_2,
		[]*big.Int{big.NewInt(69), big.NewInt(420), big.NewInt(2022)}, // Upkeep IDs
		[]common.Address{fromAddress},
		upkeepConfig1_2,
		1) // 1 new upkeep to sync
	synchronizer.ExportedFullSync()

	cltest.AssertCount(t, db, "keeper_registries", 1)
	cltest.AssertCount(t, db, "upkeep_registrations", 3)
	assertUpkeepIDs(t, db, []int64{69, 420, 2022})

	ethMock.AssertExpectations(t)
}

// TODO (sc-36399) support all v1.2 logs
func Test_RegistrySynchronizer_ConfigSetLog(t *testing.T) {
	db, synchronizer, ethMock, lb, job := setupRegistrySync(t, keeper.RegistryVersion_1_1)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	mockRegistry1_1(
		t,
		ethMock,
		contractAddress,
		registryConfig1_1,
		[]common.Address{fromAddress},
		[]*big.Int{},
		big.NewInt(0),
		upkeepConfig1_1,
		0)

	require.NoError(t, synchronizer.Start(testutils.Context(t)))
	defer synchronizer.Close()
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	var registry keeper.Registry
	require.NoError(t, db.Get(&registry, `SELECT * FROM keeper_registries`))

	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_1ABI, contractAddress)
	newConfig := registryConfig1_1
	newConfig.BlockCountPerTurn = big.NewInt(40) // change from default
	registryMock.MockResponse("getKeeperList", []common.Address{fromAddress}).Once()
	registryMock.MockResponse("getConfig", newConfig).Once()

	cfg := cltest.NewTestGeneralConfig(t)
	head := cltest.MustInsertHead(t, db, cfg, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_1.KeeperRegistryConfigSet{}
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

// TODO (sc-36399) support all v1.2 logs
func Test_RegistrySynchronizer_KeepersUpdatedLog(t *testing.T) {
	db, synchronizer, ethMock, lb, job := setupRegistrySync(t, keeper.RegistryVersion_1_1)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	mockRegistry1_1(
		t,
		ethMock,
		contractAddress,
		registryConfig1_1,
		[]common.Address{fromAddress},
		[]*big.Int{},
		big.NewInt(0),
		upkeepConfig1_1,
		0)

	require.NoError(t, synchronizer.Start(testutils.Context(t)))
	defer synchronizer.Close()
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	var registry keeper.Registry
	require.NoError(t, db.Get(&registry, `SELECT * FROM keeper_registries`))

	addresses := []common.Address{fromAddress, testutils.NewAddress()} // change from default
	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_1ABI, contractAddress)
	registryMock.MockResponse("getConfig", registryConfig1_1).Once()
	registryMock.MockResponse("getKeeperList", addresses).Once()

	cfg := cltest.NewTestGeneralConfig(t)
	head := cltest.MustInsertHead(t, db, cfg, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_1.KeeperRegistryKeepersUpdated{}
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

// TODO (sc-36399) support all v1.2 logs
func Test_RegistrySynchronizer_UpkeepCanceledLog(t *testing.T) {
	db, synchronizer, ethMock, lb, job := setupRegistrySync(t, keeper.RegistryVersion_1_1)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	mockRegistry1_1(
		t,
		ethMock,
		contractAddress,
		registryConfig1_1,
		[]common.Address{fromAddress},
		[]*big.Int{},
		big.NewInt(3),
		upkeepConfig1_1,
		3)

	require.NoError(t, synchronizer.Start(testutils.Context(t)))
	defer func() { require.NoError(t, synchronizer.Close()) }()
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	cltest.WaitForCount(t, db, "upkeep_registrations", 3)

	cfg := cltest.NewTestGeneralConfig(t)
	head := cltest.MustInsertHead(t, db, cfg, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_1.KeeperRegistryUpkeepCanceled{Id: big.NewInt(1)}
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

// TODO (sc-36399) support all v1.2 logs
func Test_RegistrySynchronizer_UpkeepRegisteredLog(t *testing.T) {
	db, synchronizer, ethMock, lb, job := setupRegistrySync(t, keeper.RegistryVersion_1_1)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	mockRegistry1_1(
		t,
		ethMock,
		contractAddress,
		registryConfig1_1,
		[]common.Address{fromAddress},
		[]*big.Int{},
		big.NewInt(1),
		upkeepConfig1_1,
		1)

	require.NoError(t, synchronizer.Start(testutils.Context(t)))
	defer synchronizer.Close()
	cltest.WaitForCount(t, db, "keeper_registries", 1)

	cfg := cltest.NewTestGeneralConfig(t)
	head := cltest.MustInsertHead(t, db, cfg, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_1.KeeperRegistryUpkeepRegistered{Id: big.NewInt(3)}
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

// TODO (sc-36399) support all v1.2 logs
func Test_RegistrySynchronizer_UpkeepPerformedLog(t *testing.T) {
	g := gomega.NewWithT(t)

	db, synchronizer, ethMock, lb, job := setupRegistrySync(t, keeper.RegistryVersion_1_1)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	mockRegistry1_1(
		t,
		ethMock,
		contractAddress,
		registryConfig1_1,
		[]common.Address{fromAddress},
		[]*big.Int{},
		big.NewInt(1),
		upkeepConfig1_1,
		1)

	require.NoError(t, synchronizer.Start(testutils.Context(t)))
	defer synchronizer.Close()
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	cltest.WaitForCount(t, db, "upkeep_registrations", 1)

	pgtest.MustExec(t, db, `UPDATE upkeep_registrations SET last_run_block_height = 100`)

	cfg := cltest.NewTestGeneralConfig(t)
	head := cltest.MustInsertHead(t, db, cfg, 1)
	rawLog := types.Log{BlockHash: head.Hash, BlockNumber: 200}
	log := registry1_1.KeeperRegistryUpkeepPerformed{Id: big.NewInt(0), From: fromAddress}
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
	}, cltest.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal(int64(200)))

	g.Eventually(func() int64 {
		var upkeep keeper.UpkeepRegistration
		err := db.Get(&upkeep, `SELECT * FROM upkeep_registrations`)
		require.NoError(t, err)
		return upkeep.LastKeeperIndex.Int64
	}, cltest.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal(int64(0)))

	ethMock.AssertExpectations(t)
	logBroadcast.AssertExpectations(t)
}
