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
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	registry1_2 "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper1_2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
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
) {
	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_2ABI, contractAddress)

	state := registryState1_2
	state.NumUpkeeps = big.NewInt(int64(len(activeUpkeepIDs)))
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
		3) // sync all 3
	synchronizer.ExportedFullSync()

	cltest.AssertCount(t, db, "keeper_registries", 1)
	cltest.AssertCount(t, db, "upkeep_registrations", 3)

	// Last keeper index should be set correctly on upkeep
	g.Eventually(func() bool {
		var upkeep keeper.UpkeepRegistration
		err := db.Get(&upkeep, `SELECT * FROM upkeep_registrations`)
		require.NoError(t, err)
		return upkeep.LastKeeperIndex.Valid
	}, cltest.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal(true))
	g.Eventually(func() int64 {
		var upkeep keeper.UpkeepRegistration
		err := db.Get(&upkeep, `SELECT * FROM upkeep_registrations`)
		require.NoError(t, err)
		return upkeep.LastKeeperIndex.Int64
	}, cltest.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal(int64(0)))

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

	cfg := cltest.NewTestGeneralConfig(t)
	head := cltest.MustInsertHead(t, db, cfg, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_2.KeeperRegistryConfigSet{}
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

	cfg := cltest.NewTestGeneralConfig(t)
	head := cltest.MustInsertHead(t, db, cfg, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_2.KeeperRegistryKeepersUpdated{}
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
		3)

	require.NoError(t, synchronizer.Start(testutils.Context(t)))
	defer func() { require.NoError(t, synchronizer.Close()) }()
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	cltest.WaitForCount(t, db, "upkeep_registrations", 3)

	cfg := cltest.NewTestGeneralConfig(t)
	head := cltest.MustInsertHead(t, db, cfg, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_2.KeeperRegistryUpkeepCanceled{Id: big.NewInt(3)}
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
		1)

	require.NoError(t, synchronizer.Start(testutils.Context(t)))
	defer synchronizer.Close()
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	cltest.WaitForCount(t, db, "upkeep_registrations", 1)

	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_2ABI, contractAddress)
	registryMock.MockResponse("getUpkeep", upkeepConfig1_2).Once()

	cfg := cltest.NewTestGeneralConfig(t)
	head := cltest.MustInsertHead(t, db, cfg, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_2.KeeperRegistryUpkeepRegistered{Id: big.NewInt(420)}
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
		1)

	require.NoError(t, synchronizer.Start(testutils.Context(t)))
	defer synchronizer.Close()
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	cltest.WaitForCount(t, db, "upkeep_registrations", 1)

	pgtest.MustExec(t, db, `UPDATE upkeep_registrations SET last_run_block_height = 100`)

	cfg := cltest.NewTestGeneralConfig(t)
	head := cltest.MustInsertHead(t, db, cfg, 1)
	rawLog := types.Log{BlockHash: head.Hash, BlockNumber: 200}
	log := registry1_2.KeeperRegistryUpkeepPerformed{Id: big.NewInt(3), From: fromAddress}
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
		1)

	require.NoError(t, synchronizer.Start(testutils.Context(t)))
	defer synchronizer.Close()
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	getExecuteGas := func() uint64 {
		var upkeep keeper.UpkeepRegistration
		err := db.Get(&upkeep, `SELECT * FROM upkeep_registrations`)
		require.NoError(t, err)
		return upkeep.ExecuteGas
	}
	g.Eventually(getExecuteGas, cltest.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal(uint64(2_000_000)))

	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_2ABI, contractAddress)
	newConfig := upkeepConfig1_2
	newConfig.ExecuteGas = 4_000_000 // change from default
	registryMock.MockResponse("getUpkeep", newConfig).Once()

	cfg := cltest.NewTestGeneralConfig(t)
	head := cltest.MustInsertHead(t, db, cfg, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_2.KeeperRegistryUpkeepGasLimitSet{Id: big.NewInt(3), GasLimit: big.NewInt(4_000_000)}
	logBroadcast := new(logmocks.Broadcast)
	logBroadcast.On("DecodedLog").Return(&log)
	logBroadcast.On("RawLog").Return(rawLog)
	logBroadcast.On("String").Maybe().Return("")
	lb.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)
	lb.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

	// Do the thing
	synchronizer.HandleLog(logBroadcast)

	g.Eventually(getExecuteGas, cltest.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal(uint64(4_000_000)))
	ethMock.AssertExpectations(t)
	logBroadcast.AssertExpectations(t)
}
