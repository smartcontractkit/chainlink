package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	logmocks "github.com/smartcontractkit/chainlink/core/chains/evm/log/mocks"
	evmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	registry1_2 "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper1_2"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
)

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

var registryState1_2 = registry1_2.State{
	Nonce:               uint32(0),
	OwnerLinkBalance:    big.NewInt(1000000000000000000),
	ExpectedLinkBalance: big.NewInt(1000000000000000000),
	NumUpkeeps:          big.NewInt(0),
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

func mockRegistry1_2(
	t *testing.T,
	ethMock *evmmocks.Client,
	contractAddress common.Address,
	config registry1_2.Config,
	activeUpkeepIDs []*big.Int,
	keeperList []common.Address,
	upkeepConfig registry1_2.GetUpkeep,
	timesGetUpkeepMock int,
	getStateTime int,
	getActiveUpkeepIDsTime int,
) {
	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_2ABI, contractAddress)

	state := registryState1_2
	state.NumUpkeeps = big.NewInt(int64(len(activeUpkeepIDs)))
	var getState = registry1_2.GetState{
		State:   state,
		Config:  config,
		Keepers: keeperList,
	}
	ethMock.On("HeaderByNumber", mock.Anything, mock.Anything).
		Return(&types.Header{Number: big.NewInt(10)}, nil)
	if getStateTime > 0 {
		registryMock.MockResponse("getState", getState).Times(getStateTime)
	}
	if getActiveUpkeepIDsTime > 0 {
		registryMock.MockResponse("getActiveUpkeepIDs", activeUpkeepIDs).Times(getActiveUpkeepIDsTime)
	}
	if timesGetUpkeepMock > 0 {
		registryMock.MockResponse("getUpkeep", upkeepConfig).Times(timesGetUpkeepMock)
	}
}

func Test_LogListenerOpts1_2(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	scopedConfig := evmtest.NewChainScopedConfig(t, configtest.NewGeneralConfig(t, nil))
	korm := keeper.NewORM(db, logger.TestLogger(t), scopedConfig, nil)
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	j := cltest.MustInsertKeeperJob(t, db, korm, cltest.NewEIP55Address(), cltest.NewEIP55Address())

	contractAddress := j.KeeperSpec.ContractAddress.Address()
	registryMock := cltest.NewContractMockReceiver(t, ethClient, keeper.Registry1_1ABI, contractAddress)
	registryMock.MockResponse("typeAndVersion", "KeeperRegistry 1.2.0").Once()

	registryWrapper, err := keeper.NewRegistryWrapper(j.KeeperSpec.ContractAddress, ethClient)
	require.NoError(t, err)

	logListenerOpts, err := registryWrapper.GetLogListenerOpts(1, nil)
	require.NoError(t, err)

	require.Contains(t, logListenerOpts.LogsWithTopics, registry1_2.KeeperRegistryKeepersUpdated{}.Topic(), "Registry should listen to KeeperRegistryKeepersUpdated log")
	require.Contains(t, logListenerOpts.LogsWithTopics, registry1_2.KeeperRegistryConfigSet{}.Topic(), "Registry should listen to KeeperRegistryConfigSet log")
	require.Contains(t, logListenerOpts.LogsWithTopics, registry1_2.KeeperRegistryUpkeepCanceled{}.Topic(), "Registry should listen to KeeperRegistryUpkeepCanceled log")
	require.Contains(t, logListenerOpts.LogsWithTopics, registry1_2.KeeperRegistryUpkeepRegistered{}.Topic(), "Registry should listen to KeeperRegistryUpkeepRegistered log")
	require.Contains(t, logListenerOpts.LogsWithTopics, registry1_2.KeeperRegistryUpkeepPerformed{}.Topic(), "Registry should listen to KeeperRegistryUpkeepPerformed log")
	require.Contains(t, logListenerOpts.LogsWithTopics, registry1_2.KeeperRegistryUpkeepGasLimitSet{}.Topic(), "Registry should listen to KeeperRegistryUpkeepGasLimitSet log")
	require.Contains(t, logListenerOpts.LogsWithTopics, registry1_2.KeeperRegistryUpkeepMigrated{}.Topic(), "Registry should listen to KeeperRegistryUpkeepMigrated log")
	require.Contains(t, logListenerOpts.LogsWithTopics, registry1_2.KeeperRegistryUpkeepReceived{}.Topic(), "Registry should listen to KeeperRegistryUpkeepReceived log")
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
		0,
		2,
		0)

	err := synchronizer.Start(testutils.Context(t))
	require.NoError(t, err)
	defer synchronizer.Close()

	cltest.WaitForCount(t, db, "keeper_registries", 1)

	err = synchronizer.Start(testutils.Context(t))
	require.Error(t, err)
}

func Test_RegistrySynchronizer1_2_FullSync(t *testing.T) {
	g := gomega.NewWithT(t)
	db, synchronizer, ethMock, _, job := setupRegistrySync(t, keeper.RegistryVersion_1_2)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	upkeepConfig := upkeepConfig1_2
	upkeepConfig.LastKeeper = fromAddress
	mockRegistry1_2(
		t,
		ethMock,
		contractAddress,
		registryConfig1_2,
		[]*big.Int{big.NewInt(3), big.NewInt(69), big.NewInt(420)}, // Upkeep IDs
		[]common.Address{fromAddress},
		upkeepConfig,
		3, // sync all 3
		2,
		1)
	synchronizer.ExportedFullSync()

	cltest.AssertCount(t, db, "keeper_registries", 1)
	cltest.AssertCount(t, db, "upkeep_registrations", 3)

	// Last keeper index should be set correctly on upkeep
	g.Eventually(func() bool {
		var upkeep keeper.UpkeepRegistration
		err := db.Get(&upkeep, `SELECT * FROM upkeep_registrations`)
		require.NoError(t, err)
		return upkeep.LastKeeperIndex.Valid
	}, testutils.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal(true))
	g.Eventually(func() int64 {
		var upkeep keeper.UpkeepRegistration
		err := db.Get(&upkeep, `SELECT * FROM upkeep_registrations`)
		require.NoError(t, err)
		return upkeep.LastKeeperIndex.Int64
	}, testutils.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal(int64(0)))

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
	require.Equal(t, upkeepConfig1_2.ExecuteGas, upkeepRegistration.ExecuteGas)

	assertUpkeepIDs(t, db, []int64{3, 69, 420})

	// 2nd sync. Cancel upkeep (id 3) and add a new upkeep (id 2022)
	mockRegistry1_2(
		t,
		ethMock,
		contractAddress,
		registryConfig1_2,
		[]*big.Int{big.NewInt(69), big.NewInt(420), big.NewInt(2022)}, // Upkeep IDs
		[]common.Address{fromAddress},
		upkeepConfig1_2,
		1, // 1 new upkeep to sync
		2,
		1)
	synchronizer.ExportedFullSync()

	cltest.AssertCount(t, db, "keeper_registries", 1)
	cltest.AssertCount(t, db, "upkeep_registrations", 3)
	assertUpkeepIDs(t, db, []int64{69, 420, 2022})
}

func Test_RegistrySynchronizer1_2_ConfigSetLog(t *testing.T) {
	db, synchronizer, ethMock, lb, job := setupRegistrySync(t, keeper.RegistryVersion_1_2)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	mockRegistry1_2(
		t,
		ethMock,
		contractAddress,
		registryConfig1_2,
		[]*big.Int{}, // Upkeep IDs
		[]common.Address{fromAddress},
		upkeepConfig1_2,
		0,
		2,
		0)

	require.NoError(t, synchronizer.Start(testutils.Context(t)))
	defer synchronizer.Close()
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	var registry keeper.Registry
	require.NoError(t, db.Get(&registry, `SELECT * FROM keeper_registries`))

	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_2ABI, contractAddress)
	newConfig := registryConfig1_2
	newConfig.BlockCountPerTurn = big.NewInt(40) // change from default
	registryMock.MockResponse("getState", registry1_2.GetState{
		State:   registryState1_2,
		Config:  newConfig,
		Keepers: []common.Address{fromAddress},
	}).Once()

	cfg := configtest.NewGeneralConfig(t, nil)
	head := cltest.MustInsertHead(t, db, cfg, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_2.KeeperRegistryConfigSet{}
	logBroadcast := logmocks.NewBroadcast(t)
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
}

func Test_RegistrySynchronizer1_2_KeepersUpdatedLog(t *testing.T) {
	db, synchronizer, ethMock, lb, job := setupRegistrySync(t, keeper.RegistryVersion_1_2)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	mockRegistry1_2(
		t,
		ethMock,
		contractAddress,
		registryConfig1_2,
		[]*big.Int{}, // Upkeep IDs
		[]common.Address{fromAddress},
		upkeepConfig1_2,
		0,
		2,
		0)

	require.NoError(t, synchronizer.Start(testutils.Context(t)))
	defer synchronizer.Close()
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	var registry keeper.Registry
	require.NoError(t, db.Get(&registry, `SELECT * FROM keeper_registries`))

	addresses := []common.Address{fromAddress, testutils.NewAddress()} // change from default
	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_2ABI, contractAddress)
	registryMock.MockResponse("getState", registry1_2.GetState{
		State:   registryState1_2,
		Config:  registryConfig1_2,
		Keepers: addresses,
	}).Once()

	cfg := configtest.NewGeneralConfig(t, nil)
	head := cltest.MustInsertHead(t, db, cfg, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_2.KeeperRegistryKeepersUpdated{}
	logBroadcast := logmocks.NewBroadcast(t)
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
}

func Test_RegistrySynchronizer1_2_UpkeepCanceledLog(t *testing.T) {
	db, synchronizer, ethMock, lb, job := setupRegistrySync(t, keeper.RegistryVersion_1_2)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	mockRegistry1_2(
		t,
		ethMock,
		contractAddress,
		registryConfig1_2,
		[]*big.Int{big.NewInt(3), big.NewInt(69), big.NewInt(420)}, // Upkeep IDs
		[]common.Address{fromAddress},
		upkeepConfig1_2,
		3,
		2,
		1)

	require.NoError(t, synchronizer.Start(testutils.Context(t)))
	defer func() { require.NoError(t, synchronizer.Close()) }()
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	cltest.WaitForCount(t, db, "upkeep_registrations", 3)

	cfg := configtest.NewGeneralConfig(t, nil)
	head := cltest.MustInsertHead(t, db, cfg, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_2.KeeperRegistryUpkeepCanceled{Id: big.NewInt(3)}
	logBroadcast := logmocks.NewBroadcast(t)
	logBroadcast.On("DecodedLog").Return(&log)
	logBroadcast.On("RawLog").Return(rawLog)
	logBroadcast.On("String").Maybe().Return("")
	lb.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

	// Do the thing
	synchronizer.HandleLog(logBroadcast)

	cltest.WaitForCount(t, db, "upkeep_registrations", 2)
}

func Test_RegistrySynchronizer1_2_UpkeepRegisteredLog(t *testing.T) {
	db, synchronizer, ethMock, lb, job := setupRegistrySync(t, keeper.RegistryVersion_1_2)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	mockRegistry1_2(
		t,
		ethMock,
		contractAddress,
		registryConfig1_2,
		[]*big.Int{big.NewInt(3)}, // Upkeep IDs
		[]common.Address{fromAddress},
		upkeepConfig1_2,
		1,
		2,
		1)

	require.NoError(t, synchronizer.Start(testutils.Context(t)))
	defer synchronizer.Close()
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	cltest.WaitForCount(t, db, "upkeep_registrations", 1)

	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_2ABI, contractAddress)
	registryMock.MockResponse("getUpkeep", upkeepConfig1_2).Once()

	cfg := configtest.NewGeneralConfig(t, nil)
	head := cltest.MustInsertHead(t, db, cfg, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_2.KeeperRegistryUpkeepRegistered{Id: big.NewInt(420)}
	logBroadcast := logmocks.NewBroadcast(t)
	logBroadcast.On("DecodedLog").Return(&log)
	logBroadcast.On("RawLog").Return(rawLog)
	logBroadcast.On("String").Maybe().Return("")
	lb.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

	// Do the thing
	synchronizer.HandleLog(logBroadcast)

	cltest.WaitForCount(t, db, "upkeep_registrations", 2)
}

func Test_RegistrySynchronizer1_2_UpkeepPerformedLog(t *testing.T) {
	g := gomega.NewWithT(t)

	db, synchronizer, ethMock, lb, job := setupRegistrySync(t, keeper.RegistryVersion_1_2)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	mockRegistry1_2(
		t,
		ethMock,
		contractAddress,
		registryConfig1_2,
		[]*big.Int{big.NewInt(3)}, // Upkeep IDs
		[]common.Address{fromAddress},
		upkeepConfig1_2,
		1,
		2,
		1)

	require.NoError(t, synchronizer.Start(testutils.Context(t)))
	defer synchronizer.Close()
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	cltest.WaitForCount(t, db, "upkeep_registrations", 1)

	pgtest.MustExec(t, db, `UPDATE upkeep_registrations SET last_run_block_height = 100`)

	cfg := configtest.NewGeneralConfig(t, nil)
	head := cltest.MustInsertHead(t, db, cfg, 1)
	rawLog := types.Log{BlockHash: head.Hash, BlockNumber: 200}
	log := registry1_2.KeeperRegistryUpkeepPerformed{Id: big.NewInt(3), From: fromAddress}
	logBroadcast := logmocks.NewBroadcast(t)
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
	}, testutils.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal(int64(200)))

	g.Eventually(func() int64 {
		var upkeep keeper.UpkeepRegistration
		err := db.Get(&upkeep, `SELECT * FROM upkeep_registrations`)
		require.NoError(t, err)
		return upkeep.LastKeeperIndex.Int64
	}, testutils.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal(int64(0)))
}

func Test_RegistrySynchronizer1_2_UpkeepGasLimitSetLog(t *testing.T) {
	g := gomega.NewWithT(t)
	db, synchronizer, ethMock, lb, job := setupRegistrySync(t, keeper.RegistryVersion_1_2)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	mockRegistry1_2(
		t,
		ethMock,
		contractAddress,
		registryConfig1_2,
		[]*big.Int{big.NewInt(3)}, // Upkeep IDs
		[]common.Address{fromAddress},
		upkeepConfig1_2,
		1,
		2,
		1)

	require.NoError(t, synchronizer.Start(testutils.Context(t)))
	defer synchronizer.Close()
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	cltest.WaitForCount(t, db, "upkeep_registrations", 1)

	getExecuteGas := func() uint32 {
		var upkeep keeper.UpkeepRegistration
		err := db.Get(&upkeep, `SELECT * FROM upkeep_registrations`)
		require.NoError(t, err)
		return upkeep.ExecuteGas
	}
	g.Eventually(getExecuteGas, testutils.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal(uint32(2_000_000)))

	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_2ABI, contractAddress)
	newConfig := upkeepConfig1_2
	newConfig.ExecuteGas = 4_000_000 // change from default
	registryMock.MockResponse("getUpkeep", newConfig).Once()

	cfg := configtest.NewGeneralConfig(t, nil)
	head := cltest.MustInsertHead(t, db, cfg, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_2.KeeperRegistryUpkeepGasLimitSet{Id: big.NewInt(3), GasLimit: big.NewInt(4_000_000)}
	logBroadcast := logmocks.NewBroadcast(t)
	logBroadcast.On("DecodedLog").Return(&log)
	logBroadcast.On("RawLog").Return(rawLog)
	logBroadcast.On("String").Maybe().Return("")
	lb.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

	// Do the thing
	synchronizer.HandleLog(logBroadcast)

	g.Eventually(getExecuteGas, testutils.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal(uint32(4_000_000)))
}

func Test_RegistrySynchronizer1_2_UpkeepReceivedLog(t *testing.T) {
	db, synchronizer, ethMock, lb, job := setupRegistrySync(t, keeper.RegistryVersion_1_2)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	mockRegistry1_2(
		t,
		ethMock,
		contractAddress,
		registryConfig1_2,
		[]*big.Int{big.NewInt(3)}, // Upkeep IDs
		[]common.Address{fromAddress},
		upkeepConfig1_2,
		1,
		2,
		1)

	require.NoError(t, synchronizer.Start(testutils.Context(t)))
	defer synchronizer.Close()
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	cltest.WaitForCount(t, db, "upkeep_registrations", 1)

	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_2ABI, contractAddress)
	registryMock.MockResponse("getUpkeep", upkeepConfig1_2).Once()

	cfg := configtest.NewGeneralConfig(t, nil)
	head := cltest.MustInsertHead(t, db, cfg, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_2.KeeperRegistryUpkeepReceived{Id: big.NewInt(420)}
	logBroadcast := logmocks.NewBroadcast(t)
	logBroadcast.On("DecodedLog").Return(&log)
	logBroadcast.On("RawLog").Return(rawLog)
	logBroadcast.On("String").Maybe().Return("")
	lb.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

	// Do the thing
	synchronizer.HandleLog(logBroadcast)

	cltest.WaitForCount(t, db, "upkeep_registrations", 2)
}

func Test_RegistrySynchronizer1_2_UpkeepMigratedLog(t *testing.T) {
	db, synchronizer, ethMock, lb, job := setupRegistrySync(t, keeper.RegistryVersion_1_2)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	mockRegistry1_2(
		t,
		ethMock,
		contractAddress,
		registryConfig1_2,
		[]*big.Int{big.NewInt(3), big.NewInt(69), big.NewInt(420)}, // Upkeep IDs
		[]common.Address{fromAddress},
		upkeepConfig1_2,
		3,
		2,
		1)

	require.NoError(t, synchronizer.Start(testutils.Context(t)))
	defer func() { require.NoError(t, synchronizer.Close()) }()
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	cltest.WaitForCount(t, db, "upkeep_registrations", 3)

	cfg := configtest.NewGeneralConfig(t, nil)
	head := cltest.MustInsertHead(t, db, cfg, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_2.KeeperRegistryUpkeepMigrated{Id: big.NewInt(3)}
	logBroadcast := logmocks.NewBroadcast(t)
	logBroadcast.On("DecodedLog").Return(&log)
	logBroadcast.On("RawLog").Return(rawLog)
	logBroadcast.On("String").Maybe().Return("")
	lb.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

	// Do the thing
	synchronizer.HandleLog(logBroadcast)

	cltest.WaitForCount(t, db, "upkeep_registrations", 2)
}
