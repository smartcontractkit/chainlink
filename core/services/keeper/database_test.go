package keeper_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var registryAddress = cltest.NewEIP55Address()
var fromAddress = cltest.NewEIP55Address()
var checkDataStr = "ABC123"
var checkData = common.Hex2Bytes(checkDataStr)
var executeGas = int32(10_000)
var checkGas = int32(2_000_000)
var blockCountPerTurn = int32(20)

func setupKeeperDB(t *testing.T) (*store.Store, keeper.DB, func()) {
	store, cleanup := cltest.NewStore(t)
	keeperDB := keeper.NewDBInterface(store.ORM)
	return store, keeperDB, cleanup
}

func newUpkeep(reg keeper.Registry, upkeepID int64) keeper.UpkeepRegistration {
	return keeper.UpkeepRegistration{
		UpkeepID:   upkeepID,
		ExecuteGas: executeGas,
		Registry:   reg,
		CheckData:  checkData,
	}
}

func TestKeeperDB_Registries(t *testing.T) {
	store, keeperDB, cleanup := setupKeeperDB(t)
	defer cleanup()

	_ = cltest.MustInsertKeeperRegistry(t, store)
	_ = cltest.MustInsertKeeperRegistry(t, store)

	existingRegistries, err := keeperDB.Registries()
	require.NoError(t, err)
	require.Equal(t, 2, len(existingRegistries))
}

func TestKeeperDB_UpsertUpkeep(t *testing.T) {
	store, keeperDB, cleanup := setupKeeperDB(t)
	defer cleanup()

	reg := cltest.MustInsertKeeperRegistry(t, store)

	// create upkeep
	upkeep := newUpkeep(reg, 0)
	err := keeperDB.UpsertUpkeep(upkeep)
	require.NoError(t, err)

	cltest.AssertCount(t, store, &keeper.UpkeepRegistration{}, 1)
	var upkeepFromDB keeper.UpkeepRegistration
	err = store.DB.First(&upkeepFromDB).Error
	require.NoError(t, err)
	require.Equal(t, executeGas, upkeepFromDB.ExecuteGas)
	require.Equal(t, checkData, upkeepFromDB.CheckData)

	// update upkeep
	upkeep.ExecuteGas = 20_000
	upkeep.CheckData = common.Hex2Bytes("8888")

	err = keeperDB.UpsertUpkeep(upkeep)
	require.NoError(t, err)
	cltest.AssertCount(t, store, &keeper.UpkeepRegistration{}, 1)
	err = store.DB.First(&upkeepFromDB).Error
	require.NoError(t, err)
	require.Equal(t, int32(20_000), upkeepFromDB.ExecuteGas)
	require.Equal(t, "8888", common.Bytes2Hex(upkeepFromDB.CheckData))
}

func TestKeeperDB_BatchDelete(t *testing.T) {
	store, keeperDB, cleanup := setupKeeperDB(t)
	defer cleanup()

	reg := cltest.MustInsertKeeperRegistry(t, store)

	registrations := [3]keeper.UpkeepRegistration{
		newUpkeep(reg, 0),
		newUpkeep(reg, 1),
		newUpkeep(reg, 2),
	}

	for _, reg := range registrations {
		err := store.DB.Create(&reg).Error
		require.NoError(t, err)
	}

	cltest.AssertCount(t, store, &keeper.UpkeepRegistration{}, 3)
	err := keeperDB.BatchDeleteUpkeeps(reg.ID, []int64{0, 2})
	require.NoError(t, err)
	cltest.AssertCount(t, store, &keeper.UpkeepRegistration{}, 1)
}

func TestKeeperDB_DeleteRegistryByJobID(t *testing.T) {
	store, keeperDB, cleanup := setupKeeperDB(t)
	defer cleanup()

	reg := cltest.MustInsertKeeperRegistry(t, store)

	registrations := [3]keeper.UpkeepRegistration{
		newUpkeep(reg, 0),
		newUpkeep(reg, 1),
		newUpkeep(reg, 2),
	}

	for _, reg := range registrations {
		err := store.DB.Create(&reg).Error
		require.NoError(t, err)
	}

	cltest.AssertCount(t, store, &keeper.UpkeepRegistration{}, 3)

	err := keeperDB.DeleteRegistryByJobID(reg.JobID)
	require.NoError(t, err)

	cltest.AssertCount(t, store, keeper.Registry{}, 0)
	cltest.AssertCount(t, store, &keeper.UpkeepRegistration{}, 0)
}

func TestKeeperDB_EligibleUpkeeps_BlockCountPerTurn(t *testing.T) {
	store, keeperDB, cleanup := setupKeeperDB(t)
	defer cleanup()

	blockheight := int64(40)

	reg1 := cltest.MustInsertKeeperRegistry(t, store)
	reg1.BlockCountPerTurn = 20
	require.NoError(t, store.DB.Save(&reg1).Error)
	reg2 := cltest.MustInsertKeeperRegistry(t, store)
	reg2.BlockCountPerTurn = 30
	require.NoError(t, store.DB.Save(&reg2).Error)

	upkeeps := [3]keeper.UpkeepRegistration{
		newUpkeep(reg1, 0), // our turn
		newUpkeep(reg1, 1), // our turn
		newUpkeep(reg2, 0), // not our turn
	}

	for _, upkeep := range upkeeps {
		err := keeperDB.UpsertUpkeep(upkeep)
		require.NoError(t, err)
	}

	cltest.AssertCount(t, store, &keeper.UpkeepRegistration{}, 3)

	elligibleUpkeeps, err := keeperDB.EligibleUpkeeps(blockheight)
	assert.NoError(t, err)
	assert.Len(t, elligibleUpkeeps, 2)
	assert.Equal(t, int64(0), elligibleUpkeeps[0].UpkeepID)
	assert.Equal(t, int64(1), elligibleUpkeeps[1].UpkeepID)

	// preloads registry data
	assert.Equal(t, reg1.ID, elligibleUpkeeps[0].RegistryID)
	assert.Equal(t, reg1.ID, elligibleUpkeeps[1].RegistryID)
	assert.Equal(t, reg1.CheckGas, elligibleUpkeeps[0].Registry.CheckGas)
	assert.Equal(t, reg1.CheckGas, elligibleUpkeeps[1].Registry.CheckGas)
	assert.Equal(t, reg1.ContractAddress, elligibleUpkeeps[0].Registry.ContractAddress)
	assert.Equal(t, reg1.ContractAddress, elligibleUpkeeps[1].Registry.ContractAddress)
}

func TestKeeperDB_EligibleUpkeeps_KeepersRotate(t *testing.T) {
	store, keeperDB, cleanup := setupKeeperDB(t)
	defer cleanup()

	reg := cltest.MustInsertKeeperRegistry(t, store)
	reg.NumKeepers = 5
	require.NoError(t, store.DB.Save(&reg).Error)
	upkeep := newUpkeep(reg, 0)
	err := keeperDB.UpsertUpkeep(upkeep)
	require.NoError(t, err)

	cltest.AssertCount(t, store, keeper.Registry{}, 1)
	cltest.AssertCount(t, store, &keeper.UpkeepRegistration{}, 1)

	// out of 5 valid block heights, with 5 keepers, we are eligible
	// to submit on exactly 1 of them
	list1, err := keeperDB.EligibleUpkeeps(20) // someone eligible
	require.NoError(t, err)
	list2, err := keeperDB.EligibleUpkeeps(30) // noone eligible
	require.NoError(t, err)
	list3, err := keeperDB.EligibleUpkeeps(40) // someone eligible
	require.NoError(t, err)
	list4, err := keeperDB.EligibleUpkeeps(41) // noone eligible
	require.NoError(t, err)
	list5, err := keeperDB.EligibleUpkeeps(60) // someone eligible
	require.NoError(t, err)
	list6, err := keeperDB.EligibleUpkeeps(80) // someone eligible
	require.NoError(t, err)
	list7, err := keeperDB.EligibleUpkeeps(99) // noone eligible
	require.NoError(t, err)
	list8, err := keeperDB.EligibleUpkeeps(100) // someone eligible
	require.NoError(t, err)

	totalEligible := len(list1) + len(list2) + len(list3) + len(list4) + len(list5) + len(list6) + len(list7) + len(list8)
	require.Equal(t, 1, totalEligible)
}

func TestKeeperDB_NextUpkeepID(t *testing.T) {
	store, keeperDB, cleanup := setupKeeperDB(t)
	defer cleanup()

	reg := cltest.MustInsertKeeperRegistry(t, store)

	nextID, err := keeperDB.NextUpkeepIDForRegistry(reg)
	require.NoError(t, err)
	require.Equal(t, int64(0), nextID)

	upkeep := newUpkeep(reg, 0)
	err = keeperDB.UpsertUpkeep(upkeep)
	require.NoError(t, err)

	nextID, err = keeperDB.NextUpkeepIDForRegistry(reg)
	require.NoError(t, err)
	require.Equal(t, int64(1), nextID)

	upkeep = newUpkeep(reg, 3)
	err = keeperDB.UpsertUpkeep(upkeep)
	require.NoError(t, err)

	nextID, err = keeperDB.NextUpkeepIDForRegistry(reg)
	require.NoError(t, err)
	require.Equal(t, int64(4), nextID)
}
