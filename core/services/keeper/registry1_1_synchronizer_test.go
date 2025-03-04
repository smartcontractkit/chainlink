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

	"github.com/smartcontractkit/chainlink-integrations/evm/client/clienttest"
	evmtypes "github.com/smartcontractkit/chainlink-integrations/evm/types"
	ubig "github.com/smartcontractkit/chainlink-integrations/evm/utils/big"
	logmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/log/mocks"
	registry1_1 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_1"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keeper"
)

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

func mockRegistry1_1(
	t *testing.T,
	ethMock *clienttest.Client,
	contractAddress common.Address,
	config registry1_1.GetConfig,
	keeperList []common.Address,
	cancelledUpkeeps []*big.Int,
	upkeepCount *big.Int,
	upkeepConfig registry1_1.GetUpkeep,
	timesGetUpkeepMock int,
) {
	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_1ABI, contractAddress)

	ethMock.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).
		Return(&evmtypes.Head{Number: 10}, nil)
	registryMock.MockResponse("getConfig", config).Once()
	registryMock.MockResponse("getKeeperList", keeperList).Once()
	registryMock.MockResponse("getCanceledUpkeepList", cancelledUpkeeps).Once()
	registryMock.MockResponse("getUpkeepCount", upkeepCount).Once()
	if timesGetUpkeepMock > 0 {
		registryMock.MockResponse("getUpkeep", upkeepConfig).Times(timesGetUpkeepMock)
	}
}

func Test_LogListenerOpts1_1(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	korm := keeper.NewORM(db, logger.TestLogger(t))
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	j := cltest.MustInsertKeeperJob(t, db, korm, cltest.NewEIP55Address(), cltest.NewEIP55Address())

	contractAddress := j.KeeperSpec.ContractAddress.Address()
	registryMock := cltest.NewContractMockReceiver(t, ethClient, keeper.Registry1_1ABI, contractAddress)
	registryMock.MockResponse("typeAndVersion", "KeeperRegistry 1.1.0").Once()

	registryWrapper, err := keeper.NewRegistryWrapper(j.KeeperSpec.ContractAddress, ethClient)
	require.NoError(t, err)

	logListenerOpts, err := registryWrapper.GetLogListenerOpts(1, nil)
	require.NoError(t, err)

	require.Contains(t, logListenerOpts.LogsWithTopics, registry1_1.KeeperRegistryKeepersUpdated{}.Topic(), "Registry should listen to KeeperRegistryKeepersUpdated log")
	require.Contains(t, logListenerOpts.LogsWithTopics, registry1_1.KeeperRegistryConfigSet{}.Topic(), "Registry should listen to KeeperRegistryConfigSet log")
	require.Contains(t, logListenerOpts.LogsWithTopics, registry1_1.KeeperRegistryUpkeepCanceled{}.Topic(), "Registry should listen to KeeperRegistryUpkeepCanceled log")
	require.Contains(t, logListenerOpts.LogsWithTopics, registry1_1.KeeperRegistryUpkeepRegistered{}.Topic(), "Registry should listen to KeeperRegistryUpkeepRegistered log")
	require.Contains(t, logListenerOpts.LogsWithTopics, registry1_1.KeeperRegistryUpkeepPerformed{}.Topic(), "Registry should listen to KeeperRegistryUpkeepPerformed log")
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
	defer func() { assert.NoError(t, synchronizer.Close()) }()

	cltest.WaitForCount(t, db, "keeper_registries", 1)

	err = synchronizer.Start(testutils.Context(t))
	require.Error(t, err)
}

func Test_RegistrySynchronizer_CalcPositioningConstant(t *testing.T) {
	t.Parallel()
	for _, upkeepID := range []int64{0, 1, 100, 10_000} {
		_, err := keeper.CalcPositioningConstant(ubig.NewI(upkeepID), cltest.NewEIP55Address())
		require.NoError(t, err)
	}
}

func Test_RegistrySynchronizer1_1_FullSync(t *testing.T) {
	ctx := testutils.Context(t)
	g := gomega.NewWithT(t)
	db, synchronizer, ethMock, _, job := setupRegistrySync(t, keeper.RegistryVersion_1_1)

	contractAddress := job.KeeperSpec.ContractAddress.Address()
	fromAddress := job.KeeperSpec.FromAddress.Address()
	canceledUpkeeps := []*big.Int{big.NewInt(1)}

	upkeepConfig := upkeepConfig1_1
	upkeepConfig.LastKeeper = fromAddress
	mockRegistry1_1(
		t,
		ethMock,
		contractAddress,
		registryConfig1_1,
		[]common.Address{fromAddress},
		canceledUpkeeps,
		big.NewInt(3),
		upkeepConfig,
		2) // sync only 2 (#0,#2)

	synchronizer.ExportedFullSync(ctx)

	cltest.AssertCount(t, db, "keeper_registries", 1)
	cltest.AssertCount(t, db, "upkeep_registrations", 2)

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
	require.NoError(t, db.Get(&upkeepRegistration, `SELECT * FROM upkeep_registrations`))
	require.Equal(t, job.KeeperSpec.ContractAddress, registry.ContractAddress)
	require.Equal(t, job.KeeperSpec.FromAddress, registry.FromAddress)
	require.Equal(t, int32(20), registry.BlockCountPerTurn)
	require.Equal(t, int32(0), registry.KeeperIndex)
	require.Equal(t, int32(1), registry.NumKeepers)
	require.Equal(t, upkeepConfig1_1.CheckData, upkeepRegistration.CheckData)
	require.Equal(t, upkeepConfig1_1.ExecuteGas, upkeepRegistration.ExecuteGas)

	assertUpkeepIDs(t, db, []int64{0, 2})

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
		2) // sync all 2 upkeeps (#2, #4)
	synchronizer.ExportedFullSync(ctx)

	cltest.AssertCount(t, db, "keeper_registries", 1)
	cltest.AssertCount(t, db, "upkeep_registrations", 2)
	assertUpkeepIDs(t, db, []int64{2, 4})
}

func Test_RegistrySynchronizer1_1_ConfigSetLog(t *testing.T) {
	ctx := testutils.Context(t)
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

	servicetest.Run(t, synchronizer)
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	var registry keeper.Registry
	require.NoError(t, db.Get(&registry, `SELECT * FROM keeper_registries`))

	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_1ABI, contractAddress)
	newConfig := registryConfig1_1
	newConfig.BlockCountPerTurn = big.NewInt(40) // change from default
	registryMock.MockResponse("getKeeperList", []common.Address{fromAddress}).Once()
	registryMock.MockResponse("getConfig", newConfig).Once()

	head := cltest.MustInsertHead(t, db, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_1.KeeperRegistryConfigSet{}
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

func Test_RegistrySynchronizer1_1_KeepersUpdatedLog(t *testing.T) {
	ctx := testutils.Context(t)
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

	servicetest.Run(t, synchronizer)
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	var registry keeper.Registry
	require.NoError(t, db.Get(&registry, `SELECT * FROM keeper_registries`))

	addresses := []common.Address{fromAddress, testutils.NewAddress()} // change from default
	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_1ABI, contractAddress)
	registryMock.MockResponse("getConfig", registryConfig1_1).Once()
	registryMock.MockResponse("getKeeperList", addresses).Once()

	head := cltest.MustInsertHead(t, db, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_1.KeeperRegistryKeepersUpdated{}
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
func Test_RegistrySynchronizer1_1_UpkeepCanceledLog(t *testing.T) {
	ctx := testutils.Context(t)
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

	servicetest.Run(t, synchronizer)
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	cltest.WaitForCount(t, db, "upkeep_registrations", 3)

	head := cltest.MustInsertHead(t, db, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_1.KeeperRegistryUpkeepCanceled{Id: big.NewInt(1)}
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

func Test_RegistrySynchronizer1_1_UpkeepRegisteredLog(t *testing.T) {
	ctx := testutils.Context(t)
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

	servicetest.Run(t, synchronizer)
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	cltest.WaitForCount(t, db, "upkeep_registrations", 1)

	registryMock := cltest.NewContractMockReceiver(t, ethMock, keeper.Registry1_1ABI, contractAddress)
	registryMock.MockResponse("getUpkeep", upkeepConfig1_1).Once()

	head := cltest.MustInsertHead(t, db, 1)
	rawLog := types.Log{BlockHash: head.Hash}
	log := registry1_1.KeeperRegistryUpkeepRegistered{Id: big.NewInt(1)}
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

func Test_RegistrySynchronizer1_1_UpkeepPerformedLog(t *testing.T) {
	ctx := testutils.Context(t)
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

	servicetest.Run(t, synchronizer)
	cltest.WaitForCount(t, db, "keeper_registries", 1)
	cltest.WaitForCount(t, db, "upkeep_registrations", 1)

	pgtest.MustExec(t, db, `UPDATE upkeep_registrations SET last_run_block_height = 100`)

	head := cltest.MustInsertHead(t, db, 1)
	rawLog := types.Log{BlockHash: head.Hash, BlockNumber: 200}
	log := registry1_1.KeeperRegistryUpkeepPerformed{Id: big.NewInt(0), From: fromAddress}
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
