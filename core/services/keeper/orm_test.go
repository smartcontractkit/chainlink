package keeper_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var checkData = common.Hex2Bytes("ABC123")
var executeGas = int32(10_000)

func setupKeeperDB(t *testing.T) (*store.Store, keeper.ORM, func()) {
	store, cleanup := cltest.NewStore(t)
	orm := keeper.NewORM(store.DB)
	return store, orm, cleanup
}

func newUpkeep(registry keeper.Registry, upkeepID int64) keeper.UpkeepRegistration {
	return keeper.UpkeepRegistration{
		UpkeepID:   upkeepID,
		ExecuteGas: executeGas,
		Registry:   registry,
		RegistryID: registry.ID,
		CheckData:  checkData,
	}
}

func assertLastRunHeight(t *testing.T, store *store.Store, upkeep keeper.UpkeepRegistration, height int64) {
	err := store.DB.Find(&upkeep).Error
	require.NoError(t, err)
	require.Equal(t, height, upkeep.LastRunBlockHeight)
}

func TestKeeperDB_Registries(t *testing.T) {
	t.Parallel()
	store, orm, cleanup := setupKeeperDB(t)
	defer cleanup()

	cltest.MustInsertKeeperRegistry(t, store)
	cltest.MustInsertKeeperRegistry(t, store)

	existingRegistries, err := orm.Registries(context.Background())
	require.NoError(t, err)
	require.Equal(t, 2, len(existingRegistries))
}

func TestKeeperDB_UpsertUpkeep(t *testing.T) {
	t.Parallel()
	store, orm, cleanup := setupKeeperDB(t)
	defer cleanup()

	registry, _ := cltest.MustInsertKeeperRegistry(t, store)
	upkeep := keeper.UpkeepRegistration{
		UpkeepID:            0,
		ExecuteGas:          executeGas,
		Registry:            registry,
		RegistryID:          registry.ID,
		CheckData:           checkData,
		LastRunBlockHeight:  1,
		PositioningConstant: 1,
	}
	require.NoError(t, store.DB.Create(&upkeep).Error)
	cltest.AssertCount(t, store, &keeper.UpkeepRegistration{}, 1)

	// update upkeep
	upkeep.ExecuteGas = 20_000
	upkeep.CheckData = common.Hex2Bytes("8888")
	upkeep.PositioningConstant = 2
	upkeep.LastRunBlockHeight = 2

	err := orm.UpsertUpkeep(context.Background(), &upkeep)
	require.NoError(t, err)
	cltest.AssertCount(t, store, &keeper.UpkeepRegistration{}, 1)

	var upkeepFromDB keeper.UpkeepRegistration
	err = store.DB.First(&upkeepFromDB).Error
	require.NoError(t, err)
	require.Equal(t, int32(20_000), upkeepFromDB.ExecuteGas)
	require.Equal(t, "8888", common.Bytes2Hex(upkeepFromDB.CheckData))
	require.Equal(t, int32(2), upkeepFromDB.PositioningConstant)
	require.Equal(t, int64(1), upkeepFromDB.LastRunBlockHeight) // shouldn't change on upsert
}

func TestKeeperDB_BatchDeleteUpkeepsForJob(t *testing.T) {
	t.Parallel()
	store, orm, cleanup := setupKeeperDB(t)
	defer cleanup()

	registry, job := cltest.MustInsertKeeperRegistry(t, store)

	for i := int64(0); i < 3; i++ {
		cltest.MustInsertUpkeepForRegistry(t, store, registry)
	}

	cltest.AssertCount(t, store, &keeper.UpkeepRegistration{}, 3)

	_, err := orm.BatchDeleteUpkeepsForJob(context.Background(), job.ID, []int64{0, 2})
	require.NoError(t, err)
	cltest.AssertCount(t, store, &keeper.UpkeepRegistration{}, 1)

	var remainingUpkeep keeper.UpkeepRegistration
	err = store.DB.First(&remainingUpkeep).Error
	require.NoError(t, err)
	require.Equal(t, int64(1), remainingUpkeep.UpkeepID)
}

func TestKeeperDB_EligibleUpkeeps_BlockCountPerTurn(t *testing.T) {
	t.Parallel()
	store, orm, cleanup := setupKeeperDB(t)
	defer cleanup()

	blockheight := int64(61)
	gracePeriod := int64(5)

	reg1, _ := cltest.MustInsertKeeperRegistry(t, store)
	reg1.BlockCountPerTurn = 20
	require.NoError(t, store.DB.Save(&reg1).Error)
	reg2, _ := cltest.MustInsertKeeperRegistry(t, store)
	reg2.BlockCountPerTurn = 25
	require.NoError(t, store.DB.Save(&reg2).Error)

	upkeeps := [5]keeper.UpkeepRegistration{
		newUpkeep(reg1, 0),
		newUpkeep(reg1, 1),
		newUpkeep(reg1, 2),
		newUpkeep(reg2, 0),
		newUpkeep(reg2, 1),
	}

	upkeeps[0].LastRunBlockHeight = 0  // Never run
	upkeeps[1].LastRunBlockHeight = 41 // Run last turn
	upkeeps[2].LastRunBlockHeight = 59 // Run last turn but still inside grace period (EXCLUDE)
	upkeeps[3].LastRunBlockHeight = 46 // Run last turn
	upkeeps[4].LastRunBlockHeight = 51 // Run this turn, outside grace period (EXCLUDE)

	for _, upkeep := range upkeeps {
		err := orm.UpsertUpkeep(context.Background(), &upkeep)
		require.NoError(t, err)
	}

	cltest.AssertCount(t, store, &keeper.UpkeepRegistration{}, 5)

	eligibleUpkeeps, err := orm.EligibleUpkeeps(context.Background(), blockheight, gracePeriod)
	assert.NoError(t, err)

	require.Len(t, eligibleUpkeeps, 3)
	assert.Equal(t, int64(0), eligibleUpkeeps[0].UpkeepID)
	assert.Equal(t, int64(1), eligibleUpkeeps[1].UpkeepID)
	assert.Equal(t, int64(0), eligibleUpkeeps[2].UpkeepID)

	// preloads registry data
	assert.Equal(t, reg1.ID, eligibleUpkeeps[0].RegistryID)
	assert.Equal(t, reg1.ID, eligibleUpkeeps[1].RegistryID)
	assert.Equal(t, reg2.ID, eligibleUpkeeps[2].RegistryID)
	assert.Equal(t, reg1.CheckGas, eligibleUpkeeps[0].Registry.CheckGas)
	assert.Equal(t, reg1.CheckGas, eligibleUpkeeps[1].Registry.CheckGas)
	assert.Equal(t, reg2.CheckGas, eligibleUpkeeps[2].Registry.CheckGas)
	assert.Equal(t, reg1.ContractAddress, eligibleUpkeeps[0].Registry.ContractAddress)
	assert.Equal(t, reg1.ContractAddress, eligibleUpkeeps[1].Registry.ContractAddress)
	assert.Equal(t, reg2.ContractAddress, eligibleUpkeeps[2].Registry.ContractAddress)
}

func TestKeeperDB_EligibleUpkeeps_GracePeriod(t *testing.T) {
	t.Parallel()
	store, orm, cleanup := setupKeeperDB(t)
	defer cleanup()

	blockheight := int64(120)
	gracePeriod := int64(100)

	registry, _ := cltest.MustInsertKeeperRegistry(t, store)
	upkeep1 := newUpkeep(registry, 0)
	upkeep1.LastRunBlockHeight = 0
	upkeep2 := newUpkeep(registry, 1)
	upkeep2.LastRunBlockHeight = 19
	upkeep3 := newUpkeep(registry, 2)
	upkeep3.LastRunBlockHeight = 20

	for _, upkeep := range [3]keeper.UpkeepRegistration{upkeep1, upkeep2, upkeep3} {
		err := orm.UpsertUpkeep(context.Background(), &upkeep)
		require.NoError(t, err)
	}

	cltest.AssertCount(t, store, &keeper.UpkeepRegistration{}, 3)

	eligibleUpkeeps, err := orm.EligibleUpkeeps(context.Background(), blockheight, gracePeriod)
	assert.NoError(t, err)
	assert.Len(t, eligibleUpkeeps, 2)
	assert.Equal(t, int64(0), eligibleUpkeeps[0].UpkeepID)
	assert.Equal(t, int64(1), eligibleUpkeeps[1].UpkeepID)
}

func TestKeeperDB_EligibleUpkeeps_KeepersRotate(t *testing.T) {
	t.Parallel()
	store, orm, cleanup := setupKeeperDB(t)
	defer cleanup()

	registry, _ := cltest.MustInsertKeeperRegistry(t, store)
	registry.NumKeepers = 5
	require.NoError(t, store.DB.Save(&registry).Error)
	cltest.MustInsertUpkeepForRegistry(t, store, registry)

	cltest.AssertCount(t, store, keeper.Registry{}, 1)
	cltest.AssertCount(t, store, &keeper.UpkeepRegistration{}, 1)

	// out of 5 valid block ranges, with 5 keepers, we are eligible
	// to submit on exactly 1 of them
	list1, err := orm.EligibleUpkeeps(context.Background(), 20, 0)
	require.NoError(t, err)
	list2, err := orm.EligibleUpkeeps(context.Background(), 41, 0)
	require.NoError(t, err)
	list3, err := orm.EligibleUpkeeps(context.Background(), 62, 0)
	require.NoError(t, err)
	list4, err := orm.EligibleUpkeeps(context.Background(), 83, 0)
	require.NoError(t, err)
	list5, err := orm.EligibleUpkeeps(context.Background(), 104, 0)
	require.NoError(t, err)

	totalEligible := len(list1) + len(list2) + len(list3) + len(list4) + len(list5)
	require.Equal(t, 1, totalEligible)
}

func TestKeeperDB_EligibleUpkeeps_KeepersCycleAllUpkeeps(t *testing.T) {
	t.Parallel()
	store, orm, cleanup := setupKeeperDB(t)
	defer cleanup()

	registry, _ := cltest.MustInsertKeeperRegistry(t, store)
	registry.NumKeepers = 5
	registry.KeeperIndex = 3
	require.NoError(t, store.DB.Save(&registry).Error)

	for i := 0; i < 1000; i++ {
		cltest.MustInsertUpkeepForRegistry(t, store, registry)
	}

	cltest.AssertCount(t, store, keeper.Registry{}, 1)
	cltest.AssertCount(t, store, &keeper.UpkeepRegistration{}, 1000)

	// in a full cycle, each node should be responsible for each upkeep exactly once
	list1, err := orm.EligibleUpkeeps(context.Background(), 20, 0) // someone eligible
	require.NoError(t, err)
	list2, err := orm.EligibleUpkeeps(context.Background(), 40, 0) // someone eligible
	require.NoError(t, err)
	list3, err := orm.EligibleUpkeeps(context.Background(), 60, 0) // someone eligible
	require.NoError(t, err)
	list4, err := orm.EligibleUpkeeps(context.Background(), 80, 0) // someone eligible
	require.NoError(t, err)
	list5, err := orm.EligibleUpkeeps(context.Background(), 100, 0) // someone eligible
	require.NoError(t, err)

	totalEligible := len(list1) + len(list2) + len(list3) + len(list4) + len(list5)
	require.Equal(t, 1000, totalEligible)
}

func TestKeeperDB_NextUpkeepID(t *testing.T) {
	t.Parallel()
	store, orm, cleanup := setupKeeperDB(t)
	defer cleanup()

	registry, _ := cltest.MustInsertKeeperRegistry(t, store)

	nextID, err := orm.LowestUnsyncedID(context.Background(), registry)
	require.NoError(t, err)
	require.Equal(t, int64(0), nextID)

	upkeep := newUpkeep(registry, 0)
	err = orm.UpsertUpkeep(context.Background(), &upkeep)
	require.NoError(t, err)

	nextID, err = orm.LowestUnsyncedID(context.Background(), registry)
	require.NoError(t, err)
	require.Equal(t, int64(1), nextID)

	upkeep = newUpkeep(registry, 3)
	err = orm.UpsertUpkeep(context.Background(), &upkeep)
	require.NoError(t, err)

	nextID, err = orm.LowestUnsyncedID(context.Background(), registry)
	require.NoError(t, err)
	require.Equal(t, int64(4), nextID)
}

func TestKeeperDB_SetLastRunHeightForUpkeepOnJob(t *testing.T) {
	t.Parallel()
	store, orm, cleanup := setupKeeperDB(t)
	defer cleanup()

	registry, j := cltest.MustInsertKeeperRegistry(t, store)
	upkeep := cltest.MustInsertUpkeepForRegistry(t, store, registry)

	orm.SetLastRunHeightForUpkeepOnJob(orm.DB, j.ID, upkeep.UpkeepID, 100)
	assertLastRunHeight(t, store, upkeep, 100)
	orm.SetLastRunHeightForUpkeepOnJob(orm.DB, j.ID, upkeep.UpkeepID, 0)
	assertLastRunHeight(t, store, upkeep, 0)
}

func TestKeeperDB_CreateEthTransactionForUpkeep(t *testing.T) {
	t.Parallel()
	store, orm, cleanup := setupKeeperDB(t)
	defer cleanup()

	registry, _ := cltest.MustInsertKeeperRegistry(t, store)
	upkeep := cltest.MustInsertUpkeepForRegistry(t, store, registry)

	payload := common.Hex2Bytes("1234")
	gasBuffer := int32(200_000)

	ethTX, err := orm.CreateEthTransactionForUpkeep(postgres.MustSQLDB(orm.DB), upkeep, payload, 500)
	require.NoError(t, err)

	require.Equal(t, registry.FromAddress.Address(), ethTX.FromAddress)
	require.Equal(t, registry.ContractAddress.Address(), ethTX.ToAddress)
	require.Equal(t, payload, ethTX.EncodedPayload)
	require.Equal(t, upkeep.ExecuteGas+gasBuffer, int32(ethTX.GasLimit))
}
