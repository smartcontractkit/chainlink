package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"

	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	logmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/log/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	registry1_3 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_3"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keeper"
)

var registryConfig1_3 = registry1_3.Config{
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

var registryState1_3 = registry1_3.State{
	Nonce:               uint32(0),
	OwnerLinkBalance:    big.NewInt(1000000000000000000),
	ExpectedLinkBalance: big.NewInt(1000000000000000000),
	NumUpkeeps:          big.NewInt(0),
}

var upkeepConfig1_3 = registry1_3.GetUpkeep{
	Target:              testutils.NewAddress(),
	ExecuteGas:          2_000_000,
	CheckData:           common.Hex2Bytes("1234"),
	Balance:             big.NewInt(1000000000000000000),
	LastKeeper:          testutils.NewAddress(),
	Admin:               testutils.NewAddress(),
	MaxValidBlocknumber: 1_000_000_000,
	AmountSpent:         big.NewInt(0),
}

func mockRegistry1_3(
	t *testing.T,
	ethMock *evmclimocks.Client,
	contractAddress common.Address,
	config registry1_3.Config,
	activeUpkeepIDs []*big.Int,
	keeperList []common.Address,
	upkeepConfig registry1_3.GetUpkeep,
	timesGetUpkeepMock int,
	getStateTime int,
	getActiveUpkeepIDsTime int,
) {
	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_3ABI, contractAddress)

	state := registryState1_3
	state.NumUpkeeps = big.NewInt(int64(len(activeUpkeepIDs)))
	var getState = registry1_3.GetState{
		State:   state,
		Config:  config,
		Keepers: keeperList,
	}
	ethMock.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).
		Return(&evmtypes.Head{Number: 10}, nil)
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

func Test_LogListenerOpts1_3(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	korm := keeper.NewORM(db, logger.TestLogger(t))
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	j := cltest.MustInsertKeeperJob(t, db, korm, cltest.NewEIP55Address(), cltest.NewEIP55Address())

	contractAddress := j.KeeperSpec.ContractAddress.Address()
	registryMock := cltest.NewContractMockReceiver(t, ethClient, keeper.Registry1_1ABI, contractAddress)
	registryMock.MockResponse("typeAndVersion", "KeeperRegistry 1.3.0").Once()

	registryWrapper, err := keeper.NewRegistryWrapper(j.KeeperSpec.ContractAddress, ethClient)
	require.NoError(t, err)

	logListenerOpts, err := registryWrapper.GetLogListenerOpts(1, nil)
	require.NoError(t, err)

	require.Contains(t, logListenerOpts.LogsWithTopics, registry1_3.KeeperRegistryKeepersUpdated{}.Topic(), "Registry should listen to KeeperRegistryKeepersUpdated log")
	require.Contains(t, logListenerOpts.LogsWithTopics, registry1_3.KeeperRegistryConfigSet{}.Topic(), "Registry should listen to KeeperRegistryConfigSet log")
	require.Contains(t, logListenerOpts.LogsWithTopics, registry1_3.KeeperRegistryUpkeepCanceled{}.Topic(), "Registry should listen to KeeperRegistryUpkeepCanceled log")
	require.Contains(t, logListenerOpts.LogsWithTopics, registry1_3.KeeperRegistryUpkeepRegistered{}.Topic(), "Registry should listen to KeeperRegistryUpkeepRegistered log")
	require.Contains(t, logListenerOpts.LogsWithTopics, registry1_3.KeeperRegistryUpkeepPerformed{}.Topic(), "Registry should listen to KeeperRegistryUpkeepPerformed log")
	require.Contains(t, logListenerOpts.LogsWithTopics, registry1_3.KeeperRegistryUpkeepGasLimitSet{}.Topic(), "Registry should listen to KeeperRegistryUpkeepGasLimitSet log")
	require.Contains(t, logListenerOpts.LogsWithTopics, registry1_3.KeeperRegistryUpkeepMigrated{}.Topic(), "Registry should listen to KeeperRegistryUpkeepMigrated log")
	require.Contains(t, logListenerOpts.LogsWithTopics, registry1_3.KeeperRegistryUpkeepReceived{}.Topic(), "Registry should listen to KeeperRegistryUpkeepReceived log")
	require.Contains(t, logListenerOpts.LogsWithTopics, registry1_3.KeeperRegistryUpkeepPaused{}.Topic(), "Registry should listen to KeeperRegistryUpkeepPaused log")
	require.Contains(t, logListenerOpts.LogsWithTopics, registry1_3.KeeperRegistryUpkeepUnpaused{}.Topic(), "Registry should listen to KeeperRegistryUpkeepUnpaused log")
	require.Contains(t, logListenerOpts.LogsWithTopics, registry1_3.KeeperRegistryUpkeepCheckDataUpdated{}.Topic(), "Registry should listen to KeeperRegistryUpkeepCheckDataUpdated log")
}

func Test_RegistrySynchronizer1_3_Start(t *testing.T) {
	db, synchronizer, ethMock, _, job := setupRegistrySync(t, keeper.RegistryVersion_1_3)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()
	mockRegistry1_3(
		t,
		ethMock,
		contractAddress,
		registryConfig1_3,
		[]*big.Int{},
		[]common.Address{fromAddress},
		upkeepConfig1_3,
		0,
		2,
		0)

	err := synchronizer.Start(testutils.Context(t))
	require.NoError(t, err)
	defer func() { assert.NoError(t, synchronizer.Close()) }()

	cltest.WaitForCount(t, db, "keeper_registries", 1)

	err = synchronizer.Start(testutils.Context(t))
	require.Error(t, err)
}

func Test_RegistrySynchronizer1_3_FullSync(t *testing.T) {
	ctx := testutils.Context(t)
	g := gomega.NewWithT(t)
	db, synchronizer, ethMock, _, job := setupRegistrySync(t, keeper.RegistryVersion_1_3)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	upkeepConfig := upkeepConfig1_3
	upkeepConfig.LastKeeper = fromAddress
	mockRegistry1_3(
		t,
		ethMock,
		contractAddress,
		registryConfig1_3,
		[]*big.Int{big.NewInt(3), big.NewInt(69), big.NewInt(420)}, // Upkeep IDs
		[]common.Address{fromAddress},
		upkeepConfig,
		3, // sync all 3
		2,
		1)
	synchronizer.ExportedFullSync(ctx)

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
	require.Equal(t, upkeepConfig1_3.CheckData, upkeepRegistration.CheckData)
	require.Equal(t, upkeepConfig1_3.ExecuteGas, upkeepRegistration.ExecuteGas)

	assertUpkeepIDs(t, db, []int64{3, 69, 420})

	// 2nd sync. Cancel upkeep (id 3) and add a new upkeep (id 2022)
	mockRegistry1_3(
		t,
		ethMock,
		contractAddress,
		registryConfig1_3,
		[]*big.Int{big.NewInt(69), big.NewInt(420), big.NewInt(2022)}, // Upkeep IDs
		[]common.Address{fromAddress},
		upkeepConfig1_3,
		3, // sync all 3 upkeeps
		2,
		1)
	synchronizer.ExportedFullSync(ctx)

	cltest.AssertCount(t, db, "keeper_registries", 1)
	cltest.AssertCount(t, db, "upkeep_registrations", 3)
	assertUpkeepIDs(t, db, []int64{69, 420, 2022})
}

func Test_RegistrySynchronizer1_3_ConfigSetLog(t *testing.T) {
	ctx := testutils.Context(t)
	db, synchronizer, ethMock, lb, job := setupRegistrySync(t, keeper.RegistryVersion_1_3)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	mockRegistry1_3(
		t,
		ethMock,
		contractAddress,
		registryConfig1_3,
		[]*big.Int{}, // Upkeep IDs
		[]common.Address{fromAddress},
		upkeepConfig1_3,
		0,
		2,
		0)

	servicetest.Run(t, synchronizer)
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	var registry keeper.Registry
	require.NoError(t, db.Get(&registry, `SELECT * FROM keeper_registries`))

	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_3ABI, contractAddress)
	newConfig := registryConfig1_3
	newConfig.BlockCountPerTurn = big.NewInt(40) // change from default
	registryMock.MockResponse("getState", registry1_3.GetState{
		State:   registryState1_3,
		Config:  newConfig,
		Keepers: []common.Address{fromAddress},
	}).Once()

	head := cltest.MustInsertHead(t, db, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_3.KeeperRegistryConfigSet{}
	logBroadcast := logmocks.NewBroadcast(t)
	logBroadcast.On("DecodedLog").Return(&log)
	logBroadcast.On("RawLog").Return(rawLog)
	logBroadcast.On("String").Maybe().Return("")
	lb.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

	// Do the thing
	synchronizer.HandleLog(ctx, logBroadcast)

	cltest.AssertRecordEventually(t, db, &registry, fmt.Sprintf(`SELECT * FROM keeper_registries WHERE id = %d`, registry.ID), func() bool {
		return registry.BlockCountPerTurn == 40
	})
	cltest.AssertCount(t, db, "keeper_registries", 1)
}

func Test_RegistrySynchronizer1_3_KeepersUpdatedLog(t *testing.T) {
	ctx := testutils.Context(t)
	db, synchronizer, ethMock, lb, job := setupRegistrySync(t, keeper.RegistryVersion_1_3)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	mockRegistry1_3(
		t,
		ethMock,
		contractAddress,
		registryConfig1_3,
		[]*big.Int{}, // Upkeep IDs
		[]common.Address{fromAddress},
		upkeepConfig1_3,
		0,
		2,
		0)

	servicetest.Run(t, synchronizer)
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	var registry keeper.Registry
	require.NoError(t, db.Get(&registry, `SELECT * FROM keeper_registries`))

	addresses := []common.Address{fromAddress, testutils.NewAddress()} // change from default
	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_3ABI, contractAddress)
	registryMock.MockResponse("getState", registry1_3.GetState{
		State:   registryState1_3,
		Config:  registryConfig1_3,
		Keepers: addresses,
	}).Once()

	head := cltest.MustInsertHead(t, db, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_3.KeeperRegistryKeepersUpdated{}
	logBroadcast := logmocks.NewBroadcast(t)
	logBroadcast.On("DecodedLog").Return(&log)
	logBroadcast.On("RawLog").Return(rawLog)
	logBroadcast.On("String").Maybe().Return("")
	lb.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

	// Do the thing
	synchronizer.HandleLog(ctx, logBroadcast)

	cltest.AssertRecordEventually(t, db, &registry, fmt.Sprintf(`SELECT * FROM keeper_registries WHERE id = %d`, registry.ID), func() bool {
		return registry.NumKeepers == 2
	})
	cltest.AssertCount(t, db, "keeper_registries", 1)
}

func Test_RegistrySynchronizer1_3_UpkeepCanceledLog(t *testing.T) {
	ctx := testutils.Context(t)
	db, synchronizer, ethMock, lb, job := setupRegistrySync(t, keeper.RegistryVersion_1_3)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	mockRegistry1_3(
		t,
		ethMock,
		contractAddress,
		registryConfig1_3,
		[]*big.Int{big.NewInt(3), big.NewInt(69), big.NewInt(420)}, // Upkeep IDs
		[]common.Address{fromAddress},
		upkeepConfig1_3,
		3,
		2,
		1)

	servicetest.Run(t, synchronizer)
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	cltest.WaitForCount(t, db, "upkeep_registrations", 3)

	head := cltest.MustInsertHead(t, db, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_3.KeeperRegistryUpkeepCanceled{Id: big.NewInt(3)}
	logBroadcast := logmocks.NewBroadcast(t)
	logBroadcast.On("DecodedLog").Return(&log)
	logBroadcast.On("RawLog").Return(rawLog)
	logBroadcast.On("String").Maybe().Return("")
	lb.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

	// Do the thing
	synchronizer.HandleLog(ctx, logBroadcast)

	cltest.WaitForCount(t, db, "upkeep_registrations", 2)
}

func Test_RegistrySynchronizer1_3_UpkeepRegisteredLog(t *testing.T) {
	ctx := testutils.Context(t)
	db, synchronizer, ethMock, lb, job := setupRegistrySync(t, keeper.RegistryVersion_1_3)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	mockRegistry1_3(
		t,
		ethMock,
		contractAddress,
		registryConfig1_3,
		[]*big.Int{big.NewInt(3)}, // Upkeep IDs
		[]common.Address{fromAddress},
		upkeepConfig1_3,
		1,
		2,
		1)

	servicetest.Run(t, synchronizer)
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	cltest.WaitForCount(t, db, "upkeep_registrations", 1)

	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_3ABI, contractAddress)
	registryMock.MockResponse("getUpkeep", upkeepConfig1_3).Once()

	head := cltest.MustInsertHead(t, db, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_3.KeeperRegistryUpkeepRegistered{Id: big.NewInt(420)}
	logBroadcast := logmocks.NewBroadcast(t)
	logBroadcast.On("DecodedLog").Return(&log)
	logBroadcast.On("RawLog").Return(rawLog)
	logBroadcast.On("String").Maybe().Return("")
	lb.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

	// Do the thing
	synchronizer.HandleLog(ctx, logBroadcast)

	cltest.WaitForCount(t, db, "upkeep_registrations", 2)
}

func Test_RegistrySynchronizer1_3_UpkeepPerformedLog(t *testing.T) {
	ctx := testutils.Context(t)
	g := gomega.NewWithT(t)

	db, synchronizer, ethMock, lb, job := setupRegistrySync(t, keeper.RegistryVersion_1_3)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	mockRegistry1_3(
		t,
		ethMock,
		contractAddress,
		registryConfig1_3,
		[]*big.Int{big.NewInt(3)}, // Upkeep IDs
		[]common.Address{fromAddress},
		upkeepConfig1_3,
		1,
		2,
		1)

	servicetest.Run(t, synchronizer)
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	cltest.WaitForCount(t, db, "upkeep_registrations", 1)

	pgtest.MustExec(t, db, `UPDATE upkeep_registrations SET last_run_block_height = 100`)

	head := cltest.MustInsertHead(t, db, 1)
	rawLog := types.Log{BlockHash: head.Hash, BlockNumber: 200}
	log := registry1_3.KeeperRegistryUpkeepPerformed{Id: big.NewInt(3), From: fromAddress}
	logBroadcast := logmocks.NewBroadcast(t)
	logBroadcast.On("DecodedLog").Return(&log)
	logBroadcast.On("RawLog").Return(rawLog)
	logBroadcast.On("String").Maybe().Return("")
	lb.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

	// Do the thing
	synchronizer.HandleLog(ctx, logBroadcast)

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

func Test_RegistrySynchronizer1_3_UpkeepGasLimitSetLog(t *testing.T) {
	ctx := testutils.Context(t)
	g := gomega.NewWithT(t)
	db, synchronizer, ethMock, lb, job := setupRegistrySync(t, keeper.RegistryVersion_1_3)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	mockRegistry1_3(
		t,
		ethMock,
		contractAddress,
		registryConfig1_3,
		[]*big.Int{big.NewInt(3)}, // Upkeep IDs
		[]common.Address{fromAddress},
		upkeepConfig1_3,
		1,
		2,
		1)

	servicetest.Run(t, synchronizer)
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	cltest.WaitForCount(t, db, "upkeep_registrations", 1)

	getExecuteGas := func() uint32 {
		var upkeep keeper.UpkeepRegistration
		err := db.Get(&upkeep, `SELECT * FROM upkeep_registrations`)
		require.NoError(t, err)
		return upkeep.ExecuteGas
	}
	g.Eventually(getExecuteGas, testutils.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal(uint32(2_000_000)))

	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_3ABI, contractAddress)
	newConfig := upkeepConfig1_3
	newConfig.ExecuteGas = 4_000_000 // change from default
	registryMock.MockResponse("getUpkeep", newConfig).Once()

	head := cltest.MustInsertHead(t, db, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_3.KeeperRegistryUpkeepGasLimitSet{Id: big.NewInt(3), GasLimit: big.NewInt(4_000_000)}
	logBroadcast := logmocks.NewBroadcast(t)
	logBroadcast.On("DecodedLog").Return(&log)
	logBroadcast.On("RawLog").Return(rawLog)
	logBroadcast.On("String").Maybe().Return("")
	lb.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

	// Do the thing
	synchronizer.HandleLog(ctx, logBroadcast)

	g.Eventually(getExecuteGas, testutils.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal(uint32(4_000_000)))
}

func Test_RegistrySynchronizer1_3_UpkeepReceivedLog(t *testing.T) {
	ctx := testutils.Context(t)
	db, synchronizer, ethMock, lb, job := setupRegistrySync(t, keeper.RegistryVersion_1_3)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	mockRegistry1_3(
		t,
		ethMock,
		contractAddress,
		registryConfig1_3,
		[]*big.Int{big.NewInt(3)}, // Upkeep IDs
		[]common.Address{fromAddress},
		upkeepConfig1_3,
		1,
		2,
		1)

	servicetest.Run(t, synchronizer)
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	cltest.WaitForCount(t, db, "upkeep_registrations", 1)

	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_3ABI, contractAddress)
	registryMock.MockResponse("getUpkeep", upkeepConfig1_3).Once()

	head := cltest.MustInsertHead(t, db, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_3.KeeperRegistryUpkeepReceived{Id: big.NewInt(420)}
	logBroadcast := logmocks.NewBroadcast(t)
	logBroadcast.On("DecodedLog").Return(&log)
	logBroadcast.On("RawLog").Return(rawLog)
	logBroadcast.On("String").Maybe().Return("")
	lb.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

	// Do the thing
	synchronizer.HandleLog(ctx, logBroadcast)

	cltest.WaitForCount(t, db, "upkeep_registrations", 2)
}

func Test_RegistrySynchronizer1_3_UpkeepMigratedLog(t *testing.T) {
	ctx := testutils.Context(t)
	db, synchronizer, ethMock, lb, job := setupRegistrySync(t, keeper.RegistryVersion_1_3)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()

	mockRegistry1_3(
		t,
		ethMock,
		contractAddress,
		registryConfig1_3,
		[]*big.Int{big.NewInt(3), big.NewInt(69), big.NewInt(420)}, // Upkeep IDs
		[]common.Address{fromAddress},
		upkeepConfig1_3,
		3,
		2,
		1)

	servicetest.Run(t, synchronizer)
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	cltest.WaitForCount(t, db, "upkeep_registrations", 3)

	head := cltest.MustInsertHead(t, db, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_3.KeeperRegistryUpkeepMigrated{Id: big.NewInt(3)}
	logBroadcast := logmocks.NewBroadcast(t)
	logBroadcast.On("DecodedLog").Return(&log)
	logBroadcast.On("RawLog").Return(rawLog)
	logBroadcast.On("String").Maybe().Return("")
	lb.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

	// Do the thing
	synchronizer.HandleLog(ctx, logBroadcast)

	// race condition: "wait for count"
	cltest.WaitForCount(t, db, "upkeep_registrations", 2)
}

func Test_RegistrySynchronizer1_3_UpkeepPausedLog_UpkeepUnpausedLog(t *testing.T) {
	ctx := testutils.Context(t)
	db, synchronizer, ethMock, lb, job := setupRegistrySync(t, keeper.RegistryVersion_1_3)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()
	upkeepId := big.NewInt(3)

	mockRegistry1_3(
		t,
		ethMock,
		contractAddress,
		registryConfig1_3,
		[]*big.Int{big.NewInt(3), big.NewInt(69), big.NewInt(420)}, // Upkeep IDs
		[]common.Address{fromAddress},
		upkeepConfig1_3,
		4,
		2,
		1)

	servicetest.Run(t, synchronizer)
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	cltest.WaitForCount(t, db, "upkeep_registrations", 3)

	head := cltest.MustInsertHead(t, db, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_3.KeeperRegistryUpkeepPaused{Id: upkeepId}
	logBroadcast := logmocks.NewBroadcast(t)
	logBroadcast.On("DecodedLog").Return(&log)
	logBroadcast.On("RawLog").Return(rawLog)
	logBroadcast.On("String").Maybe().Return("")
	lb.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

	// Do the thing
	synchronizer.HandleLog(ctx, logBroadcast)

	cltest.WaitForCount(t, db, "upkeep_registrations", 2)

	head = cltest.MustInsertHead(t, db, 2)
	rawLog = types.Log{BlockHash: head.Hash}
	unpausedlog := registry1_3.KeeperRegistryUpkeepUnpaused{Id: upkeepId}
	logBroadcast = logmocks.NewBroadcast(t)
	logBroadcast.On("DecodedLog").Return(&unpausedlog)
	logBroadcast.On("RawLog").Return(rawLog)
	logBroadcast.On("String").Maybe().Return("")
	lb.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

	// Do the thing
	synchronizer.HandleLog(ctx, logBroadcast)

	cltest.WaitForCount(t, db, "upkeep_registrations", 3)
	var upkeep keeper.UpkeepRegistration
	err := db.Get(&upkeep, `SELECT * FROM upkeep_registrations WHERE upkeep_id = $1`, ubig.New(upkeepId))
	require.NoError(t, err)

	require.Equal(t, upkeepId.String(), upkeep.UpkeepID.String())
	require.Equal(t, upkeepConfig1_3.CheckData, upkeep.CheckData)
	require.Equal(t, upkeepConfig1_3.ExecuteGas, upkeep.ExecuteGas)

	var registryId int64
	err = db.Get(&registryId, `SELECT id from keeper_registries WHERE job_id = $1`, job.ID)
	require.NoError(t, err)
	require.Equal(t, registryId, upkeep.RegistryID)
}

func Test_RegistrySynchronizer1_3_UpkeepCheckDataUpdatedLog(t *testing.T) {
	ctx := testutils.Context(t)
	g := gomega.NewWithT(t)
	db, synchronizer, ethMock, lb, job := setupRegistrySync(t, keeper.RegistryVersion_1_3)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()
	upkeepId := big.NewInt(3)

	mockRegistry1_3(
		t,
		ethMock,
		contractAddress,
		registryConfig1_3,
		[]*big.Int{upkeepId}, // Upkeep IDs
		[]common.Address{fromAddress},
		upkeepConfig1_3,
		1,
		2,
		1)

	servicetest.Run(t, synchronizer)
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	cltest.WaitForCount(t, db, "upkeep_registrations", 1)

	head := cltest.MustInsertHead(t, db, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	_ = logmocks.NewBroadcast(t)
	newCheckData := []byte("Chainlink")
	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_3ABI, contractAddress)
	newConfig := upkeepConfig1_3
	newConfig.CheckData = newCheckData // changed from default
	registryMock.MockResponse("getUpkeep", newConfig).Once()

	updatedLog := registry1_3.KeeperRegistryUpkeepCheckDataUpdated{Id: upkeepId, NewCheckData: newCheckData}
	logBroadcast := logmocks.NewBroadcast(t)
	logBroadcast.On("DecodedLog").Return(&updatedLog)
	logBroadcast.On("RawLog").Return(rawLog)
	logBroadcast.On("String").Maybe().Return("")
	lb.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

	// Do the thing
	synchronizer.HandleLog(ctx, logBroadcast)

	g.Eventually(func() []byte {
		var upkeep keeper.UpkeepRegistration
		err := db.Get(&upkeep, `SELECT * FROM upkeep_registrations WHERE upkeep_id = $1`, ubig.New(upkeepId))
		require.NoError(t, err)
		return upkeep.CheckData
	}, testutils.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal(newCheckData))
}
