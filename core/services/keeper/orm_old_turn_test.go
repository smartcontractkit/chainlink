package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestKeeperDB_EligibleUpkeeps_BlockCountPerTurn(t *testing.T) {
	t.Parallel()
	db, config, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db, config).Eth()

	blockheight := int64(63)
	gracePeriod := int64(10)

	registry, _ := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 1, 20)

	upkeeps := [5]keeper.UpkeepRegistration{
		newUpkeep(registry, 0),
		newUpkeep(registry, 1),
		newUpkeep(registry, 2),
		newUpkeep(registry, 3),
		newUpkeep(registry, 4),
	}

	upkeeps[0].LastRunBlockHeight = 0  // Never run
	upkeeps[1].LastRunBlockHeight = 41 // Run last turn, outside grade period
	upkeeps[2].LastRunBlockHeight = 46 // Run last turn, outside grade period
	upkeeps[3].LastRunBlockHeight = 59 // Run last turn, inside grace period (EXCLUDE)
	upkeeps[4].LastRunBlockHeight = 61 // Run this turn, inside grace period (EXCLUDE)

	for i := range upkeeps {
		upkeeps[i].PositioningConstant = int32(i)
		err := orm.UpsertUpkeep(&upkeeps[i])
		require.NoError(t, err)
	}

	cltest.AssertCount(t, db, "upkeep_registrations", 5)

	eligibleUpkeeps, err := orm.EligibleUpkeepsForRegistry(registry.ContractAddress, blockheight, gracePeriod)
	assert.NoError(t, err)

	// 3 out of 5 are eligible, check that ids are 0,1 or 2 but order is shuffled so can not use equals
	require.Len(t, eligibleUpkeeps, 3)
	assert.Equal(t, eligibleUpkeeps[0].UpkeepID.Cmp(utils.NewBigI(3)), -1)
	assert.Equal(t, eligibleUpkeeps[1].UpkeepID.Cmp(utils.NewBigI(3)), -1)
	assert.Equal(t, eligibleUpkeeps[2].UpkeepID.Cmp(utils.NewBigI(3)), -1)

	// preloads registry data
	assert.Equal(t, registry.ID, eligibleUpkeeps[0].RegistryID)
	assert.Equal(t, registry.ID, eligibleUpkeeps[1].RegistryID)
	assert.Equal(t, registry.ID, eligibleUpkeeps[2].RegistryID)
	assert.Equal(t, registry.CheckGas, eligibleUpkeeps[0].Registry.CheckGas)
	assert.Equal(t, registry.CheckGas, eligibleUpkeeps[1].Registry.CheckGas)
	assert.Equal(t, registry.CheckGas, eligibleUpkeeps[2].Registry.CheckGas)
	assert.Equal(t, registry.ContractAddress, eligibleUpkeeps[0].Registry.ContractAddress)
	assert.Equal(t, registry.ContractAddress, eligibleUpkeeps[1].Registry.ContractAddress)
	assert.Equal(t, registry.ContractAddress, eligibleUpkeeps[2].Registry.ContractAddress)
}

func TestKeeperDB_EligibleUpkeeps_GracePeriod(t *testing.T) {
	t.Parallel()
	db, config, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db, config).Eth()

	blockheight := int64(120)
	gracePeriod := int64(100)

	registry, _ := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 2, 20)
	upkeep1 := newUpkeep(registry, 0)
	upkeep1.LastRunBlockHeight = 0
	upkeep2 := newUpkeep(registry, 1)
	upkeep2.LastRunBlockHeight = 19
	upkeep3 := newUpkeep(registry, 2)
	upkeep3.LastRunBlockHeight = 20

	upkeeps := [3]keeper.UpkeepRegistration{upkeep1, upkeep2, upkeep3}
	for i := range upkeeps {
		err := orm.UpsertUpkeep(&upkeeps[i])
		require.NoError(t, err)
	}

	cltest.AssertCount(t, db, "upkeep_registrations", 3)

	eligibleUpkeeps, err := orm.EligibleUpkeepsForRegistry(registry.ContractAddress, blockheight, gracePeriod)
	assert.NoError(t, err)
	// 2 out of 3 are eligible, check that ids are 0 or 1 but order is shuffled so can not use equals
	assert.Len(t, eligibleUpkeeps, 2)
	assert.Equal(t, eligibleUpkeeps[0].UpkeepID.Cmp(utils.NewBigI(2)), -1)
	assert.Equal(t, eligibleUpkeeps[1].UpkeepID.Cmp(utils.NewBigI(2)), -1)
}

func TestKeeperDB_EligibleUpkeeps_KeepersRotate(t *testing.T) {
	t.Parallel()
	db, config, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db, config).Eth()

	registry, _ := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 2, 20)
	registry.NumKeepers = 5
	require.NoError(t, db.Get(&registry, `UPDATE keeper_registries SET num_keepers = 5 WHERE id = $1 RETURNING *`, registry.ID))
	cltest.MustInsertUpkeepForRegistry(t, db, config, registry)

	cltest.AssertCount(t, db, "keeper_registries", 1)
	cltest.AssertCount(t, db, "upkeep_registrations", 1)

	// out of 5 valid block ranges, with 5 keepers, we are eligible
	// to submit on exactly 1 of them
	list1, err := orm.EligibleUpkeepsForRegistry(registry.ContractAddress, 20, 0)
	require.NoError(t, err)
	list2, err := orm.EligibleUpkeepsForRegistry(registry.ContractAddress, 41, 0)
	require.NoError(t, err)
	list3, err := orm.EligibleUpkeepsForRegistry(registry.ContractAddress, 62, 0)
	require.NoError(t, err)
	list4, err := orm.EligibleUpkeepsForRegistry(registry.ContractAddress, 83, 0)
	require.NoError(t, err)
	list5, err := orm.EligibleUpkeepsForRegistry(registry.ContractAddress, 104, 0)
	require.NoError(t, err)

	totalEligible := len(list1) + len(list2) + len(list3) + len(list4) + len(list5)
	require.Equal(t, 1, totalEligible)
}

func TestKeeperDB_EligibleUpkeeps_KeepersCycleAllUpkeeps(t *testing.T) {
	t.Parallel()
	db, config, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db, config).Eth()

	registry, _ := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 2, 20)
	require.NoError(t, db.Get(&registry, `UPDATE keeper_registries SET num_keepers = 5, keeper_index = 3 WHERE id = $1 RETURNING *`, registry.ID))

	for i := 0; i < 1000; i++ {
		cltest.MustInsertUpkeepForRegistry(t, db, config, registry)
	}

	cltest.AssertCount(t, db, "keeper_registries", 1)
	cltest.AssertCount(t, db, "upkeep_registrations", 1000)

	// in a full cycle, each node should be responsible for each upkeep exactly once
	list1, err := orm.EligibleUpkeepsForRegistry(registry.ContractAddress, 20, 0) // someone eligible
	require.NoError(t, err)
	list2, err := orm.EligibleUpkeepsForRegistry(registry.ContractAddress, 40, 0) // someone eligible
	require.NoError(t, err)
	list3, err := orm.EligibleUpkeepsForRegistry(registry.ContractAddress, 60, 0) // someone eligible
	require.NoError(t, err)
	list4, err := orm.EligibleUpkeepsForRegistry(registry.ContractAddress, 80, 0) // someone eligible
	require.NoError(t, err)
	list5, err := orm.EligibleUpkeepsForRegistry(registry.ContractAddress, 100, 0) // someone eligible
	require.NoError(t, err)

	totalEligible := len(list1) + len(list2) + len(list3) + len(list4) + len(list5)
	require.Equal(t, 1000, totalEligible)
}

func TestKeeperDB_EligibleUpkeeps_FiltersByRegistry(t *testing.T) {
	t.Parallel()
	db, config, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db, config).Eth()

	registry1, _ := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 1, 20)
	registry2, _ := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 1, 20)

	cltest.MustInsertUpkeepForRegistry(t, db, config, registry1)
	cltest.MustInsertUpkeepForRegistry(t, db, config, registry2)

	cltest.AssertCount(t, db, "keeper_registries", 2)
	cltest.AssertCount(t, db, "upkeep_registrations", 2)

	list1, err := orm.EligibleUpkeepsForRegistry(registry1.ContractAddress, 20, 0)
	require.NoError(t, err)
	list2, err := orm.EligibleUpkeepsForRegistry(registry2.ContractAddress, 20, 0)
	require.NoError(t, err)

	assert.Equal(t, 1, len(list1))
	assert.Equal(t, 1, len(list2))
}

func TestKeeperDB_SetLastRunInfoForUpkeepOnJob(t *testing.T) {
	t.Parallel()
	db, config, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db, config).Eth()

	registry, j := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 2, 20)
	upkeep := cltest.MustInsertUpkeepForRegistry(t, db, config, registry)

	// check normal behavior
	err := orm.SetLastRunInfoForUpkeepOnJob(j.ID, upkeep.UpkeepID, 100, registry.FromAddress)
	require.NoError(t, err)
	assertLastRunHeight(t, db, upkeep, 100, 0)
	// check that if we put in an unknown from address nothing breaks
	err = orm.SetLastRunInfoForUpkeepOnJob(j.ID, upkeep.UpkeepID, 0, cltest.NewEIP55Address())
	require.NoError(t, err)
	assertLastRunHeight(t, db, upkeep, 100, 0)
}
