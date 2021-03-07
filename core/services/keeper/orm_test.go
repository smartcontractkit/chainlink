package keeper_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var checkData = common.Hex2Bytes("ABC123")
var executeGas = int32(10_000)

func setupKeeperDB(t *testing.T) (*store.Store, keeper.KeeperORM, func()) {
	store, cleanup := cltest.NewStore(t)
	keeperORM := keeper.NewORM(store.ORM)
	return store, keeperORM, cleanup
}

func newUpkeep(registry keeper.Registry, upkeepID int64) keeper.UpkeepRegistration {
	return keeper.UpkeepRegistration{
		UpkeepID:   upkeepID,
		ExecuteGas: executeGas,
		Registry:   registry,
		CheckData:  checkData,
	}
}

func TestKeeperDB_Registries(t *testing.T) {
	store, keeperORM, cleanup := setupKeeperDB(t)
	defer cleanup()

	cltest.MustInsertKeeperRegistry(t, store)
	cltest.MustInsertKeeperRegistry(t, store)

	existingRegistries, err := keeperORM.Registries()
	require.NoError(t, err)
	require.Equal(t, 2, len(existingRegistries))
}

func TestKeeperDB_UpsertUpkeep(t *testing.T) {
	store, keeperORM, cleanup := setupKeeperDB(t)
	defer cleanup()

	registry := cltest.MustInsertKeeperRegistry(t, store)
	upkeep := cltest.MustInsertUpkeepForRegistry(t, store, registry)

	cltest.AssertCount(t, store, &keeper.UpkeepRegistration{}, 1)
	var upkeepFromDB keeper.UpkeepRegistration
	err := store.DB.First(&upkeepFromDB).Error
	require.NoError(t, err)
	require.Equal(t, executeGas, upkeepFromDB.ExecuteGas)
	require.Equal(t, checkData, upkeepFromDB.CheckData)

	// update upkeep
	upkeep.ExecuteGas = 20_000
	upkeep.CheckData = common.Hex2Bytes("8888")

	err = keeperORM.UpsertUpkeep(&upkeep)
	require.NoError(t, err)
	cltest.AssertCount(t, store, &keeper.UpkeepRegistration{}, 1)
	err = store.DB.First(&upkeepFromDB).Error
	require.NoError(t, err)
	require.Equal(t, int32(20_000), upkeepFromDB.ExecuteGas)
	require.Equal(t, "8888", common.Bytes2Hex(upkeepFromDB.CheckData))
}

func TestKeeperDB_BatchDelete(t *testing.T) {
	store, keeperORM, cleanup := setupKeeperDB(t)
	defer cleanup()

	registry := cltest.MustInsertKeeperRegistry(t, store)

	for i := int64(0); i < 3; i++ {
		cltest.MustInsertUpkeepForRegistry(t, store, registry)
	}

	cltest.AssertCount(t, store, &keeper.UpkeepRegistration{}, 3)
	err := keeperORM.BatchDeleteUpkeeps(registry.ID, []int64{0, 2})
	require.NoError(t, err)
	cltest.AssertCount(t, store, &keeper.UpkeepRegistration{}, 1)

	var remainingUpkeep keeper.UpkeepRegistration
	err = store.DB.First(&remainingUpkeep).Error
	require.NoError(t, err)
	require.Equal(t, int64(1), remainingUpkeep.UpkeepID)
}

func TestKeeperDB_DeleteRegistryByJobID(t *testing.T) {
	store, keeperORM, cleanup := setupKeeperDB(t)
	defer cleanup()

	registry := cltest.MustInsertKeeperRegistry(t, store)

	for i := int64(0); i < 3; i++ {
		cltest.MustInsertUpkeepForRegistry(t, store, registry)
	}

	cltest.AssertCount(t, store, &keeper.UpkeepRegistration{}, 3)

	err := keeperORM.DeleteRegistryByJobID(registry.JobID)
	require.NoError(t, err)

	cltest.AssertCount(t, store, keeper.Registry{}, 0)
	cltest.AssertCount(t, store, &keeper.UpkeepRegistration{}, 0)
}

func TestKeeperDB_EligibleUpkeeps_BlockCountPerTurn(t *testing.T) {
	store, keeperORM, cleanup := setupKeeperDB(t)
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
		err := keeperORM.UpsertUpkeep(&upkeep)
		require.NoError(t, err)
	}

	cltest.AssertCount(t, store, &keeper.UpkeepRegistration{}, 3)

	elligibleUpkeeps, err := keeperORM.EligibleUpkeeps(blockheight)
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
	store, keeperORM, cleanup := setupKeeperDB(t)
	defer cleanup()

	registry := cltest.MustInsertKeeperRegistry(t, store)
	registry.NumKeepers = 5
	require.NoError(t, store.DB.Save(&registry).Error)
	cltest.MustInsertUpkeepForRegistry(t, store, registry)

	cltest.AssertCount(t, store, keeper.Registry{}, 1)
	cltest.AssertCount(t, store, &keeper.UpkeepRegistration{}, 1)

	// out of 5 valid block heights, with 5 keepers, we are eligible
	// to submit on exactly 1 of them
	list1, err := keeperORM.EligibleUpkeeps(20) // someone eligible
	require.NoError(t, err)
	list2, err := keeperORM.EligibleUpkeeps(30) // noone eligible
	require.NoError(t, err)
	list3, err := keeperORM.EligibleUpkeeps(40) // someone eligible
	require.NoError(t, err)
	list4, err := keeperORM.EligibleUpkeeps(41) // noone eligible
	require.NoError(t, err)
	list5, err := keeperORM.EligibleUpkeeps(60) // someone eligible
	require.NoError(t, err)
	list6, err := keeperORM.EligibleUpkeeps(80) // someone eligible
	require.NoError(t, err)
	list7, err := keeperORM.EligibleUpkeeps(99) // noone eligible
	require.NoError(t, err)
	list8, err := keeperORM.EligibleUpkeeps(100) // someone eligible
	require.NoError(t, err)

	totalEligible := len(list1) + len(list2) + len(list3) + len(list4) + len(list5) + len(list6) + len(list7) + len(list8)
	require.Equal(t, 1, totalEligible)
}

func TestKeeperDB_NextUpkeepID(t *testing.T) {
	store, keeperORM, cleanup := setupKeeperDB(t)
	defer cleanup()

	registry := cltest.MustInsertKeeperRegistry(t, store)

	nextID, err := keeperORM.NextUpkeepIDForRegistry(registry)
	require.NoError(t, err)
	require.Equal(t, int64(0), nextID)

	upkeep := newUpkeep(registry, 0)
	err = keeperORM.UpsertUpkeep(&upkeep)
	require.NoError(t, err)

	nextID, err = keeperORM.NextUpkeepIDForRegistry(registry)
	require.NoError(t, err)
	require.Equal(t, int64(1), nextID)

	upkeep = newUpkeep(registry, 3)
	err = keeperORM.UpsertUpkeep(&upkeep)
	require.NoError(t, err)

	nextID, err = keeperORM.NextUpkeepIDForRegistry(registry)
	require.NoError(t, err)
	require.Equal(t, int64(4), nextID)
}

func TestKeeperDB_InsertEthTXForUpkeep(t *testing.T) {
	store, keeperORM, cleanup := setupKeeperDB(t)
	defer cleanup()

	registry := cltest.MustInsertKeeperRegistry(t, store)
	upkeep := cltest.MustInsertUpkeepForRegistry(t, store, registry)

	payload := common.Hex2Bytes("1234")
	gasBuffer := int32(200_000)

	err := keeperORM.InsertEthTXForUpkeep(upkeep, payload)
	require.NoError(t, err)

	var ethTX models.EthTx
	err = store.DB.First(&ethTX).Error
	require.NoError(t, err)
	require.Equal(t, registry.FromAddress.Address(), ethTX.FromAddress)
	require.Equal(t, registry.ContractAddress.Address(), ethTX.ToAddress)
	require.Equal(t, payload, ethTX.EncodedPayload)
	require.Equal(t, upkeep.ExecuteGas+gasBuffer, int32(ethTX.GasLimit))
}
