package keeper_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestKeeperDB_RegistryByContractAddress(t *testing.T) {
	t.Parallel()
	db, config, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db, config).Eth()

	registry, _ := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 1, 20)
	cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 1, 20)

	registryByContractAddress, err := orm.RegistryByContractAddress(registry.ContractAddress)
	require.NoError(t, err)
	require.Equal(t, registry, registryByContractAddress)
}

func TestKeeperDB_NewEligibleUpkeeps_Shuffle(t *testing.T) {
	t.Parallel()
	db, config, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db, config).Eth()

	blockheight := int64(63)
	gracePeriod := int64(10)

	registry, _ := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 1, 20)

	ordered := [100]int64{}
	for i := 0; i < 100; i++ {
		k := newUpkeep(registry, int64(i))
		ordered[i] = int64(i)
		err := orm.UpsertUpkeep(&k)
		require.NoError(t, err)
	}
	cltest.AssertCount(t, db, "upkeep_registrations", 100)

	eligibleUpkeeps, err := orm.NewEligibleUpkeepsForRegistry(registry.ContractAddress, blockheight, gracePeriod, fmt.Sprintf("%b", utils.NewHash().Big()))
	assert.NoError(t, err)

	require.Len(t, eligibleUpkeeps, 100)
	shuffled := [100]int64{}
	for i := 0; i < 100; i++ {
		shuffled[i] = eligibleUpkeeps[i].UpkeepID
	}
	assert.NotEqualValues(t, ordered, shuffled)
}

func TestKeeperDB_NewEligibleUpkeeps_GracePeriod(t *testing.T) {
	t.Parallel()
	db, config, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db, config).Eth()

	registry, _ := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 2, 20)

	for i := 0; i < 100; i++ {
		cltest.MustInsertUpkeepForRegistry(t, db, config, registry)
	}

	cltest.AssertCount(t, db, "keeper_registries", 1)
	cltest.AssertCount(t, db, "upkeep_registrations", 100)

	// if current keeper index = 0 and all upkeeps last perform was done by index = 0 and still within grace period
	upkeep := keeper.UpkeepRegistration{}
	require.NoError(t, db.Get(&upkeep, `UPDATE upkeep_registrations SET last_keeper_index = 0, last_run_block_height = 10 RETURNING *`))
	list0, err := orm.NewEligibleUpkeepsForRegistry(registry.ContractAddress, 21, 100, fmt.Sprintf("%b", utils.NewHash().Big())) // none eligible
	require.NoError(t, err)
	require.Equal(t, 0, len(list0), "should be 0 as all last perform was done by current node")

	// once passed grace period
	list1, err := orm.NewEligibleUpkeepsForRegistry(registry.ContractAddress, 121, 100, fmt.Sprintf("%b", utils.NewHash().Big())) // none eligible
	require.NoError(t, err)
	require.NotEqual(t, 0, len(list1), "should get some eligible upkeeps now that they are outside grace period")
}

func TestKeeperDB_NewEligibleUpkeeps_TurnsRandom(t *testing.T) {
	t.Parallel()
	db, config, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db, config).Eth()

	registry, _ := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 3, 10)

	for i := 0; i < 1000; i++ {
		cltest.MustInsertUpkeepForRegistry(t, db, config, registry)
	}

	cltest.AssertCount(t, db, "keeper_registries", 1)
	cltest.AssertCount(t, db, "upkeep_registrations", 1000)

	// 3 keepers 10 block turns should be different every turn
	list1, err := orm.NewEligibleUpkeepsForRegistry(registry.ContractAddress, 20, 100, fmt.Sprintf("%b", utils.NewHash().Big()))
	require.NoError(t, err)
	list2, err := orm.NewEligibleUpkeepsForRegistry(registry.ContractAddress, 31, 100, fmt.Sprintf("%b", utils.NewHash().Big()))
	require.NoError(t, err)
	list3, err := orm.NewEligibleUpkeepsForRegistry(registry.ContractAddress, 42, 100, fmt.Sprintf("%b", utils.NewHash().Big()))
	require.NoError(t, err)
	list4, err := orm.NewEligibleUpkeepsForRegistry(registry.ContractAddress, 53, 100, fmt.Sprintf("%b", utils.NewHash().Big()))
	require.NoError(t, err)

	// sort before compare
	sort.Slice(list1, func(i, j int) bool {
		return list1[i].UpkeepID < list1[j].UpkeepID
	})
	sort.Slice(list2, func(i, j int) bool {
		return list2[i].UpkeepID < list2[j].UpkeepID
	})
	sort.Slice(list3, func(i, j int) bool {
		return list3[i].UpkeepID < list3[j].UpkeepID
	})
	sort.Slice(list4, func(i, j int) bool {
		return list4[i].UpkeepID < list4[j].UpkeepID
	})

	assert.NotEqual(t, list1, list2, "list1 vs list2")
	assert.NotEqual(t, list1, list3, "list1 vs list3")
	assert.NotEqual(t, list1, list4, "list1 vs list4")
}

func TestKeeperDB_NewEligibleUpkeeps_SkipIfLastPerformedByCurrentKeeper(t *testing.T) {
	t.Parallel()
	db, config, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db, config).Eth()

	registry, _ := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 2, 20)

	for i := 0; i < 100; i++ {
		cltest.MustInsertUpkeepForRegistry(t, db, config, registry)
	}

	cltest.AssertCount(t, db, "keeper_registries", 1)
	cltest.AssertCount(t, db, "upkeep_registrations", 100)

	// if current keeper index = 0 and all upkeeps last perform was done by index = 0 then skip as it would not pass required turn taking
	upkeep := keeper.UpkeepRegistration{}
	require.NoError(t, db.Get(&upkeep, `UPDATE upkeep_registrations SET last_keeper_index = 0 RETURNING *`))
	list0, err := orm.NewEligibleUpkeepsForRegistry(registry.ContractAddress, 21, 100, fmt.Sprintf("%b", utils.NewHash().Big())) // none eligible
	require.NoError(t, err)
	require.Equal(t, 0, len(list0), "should be 0 as all last perform was done by current node")
}

func TestKeeperDB_NewEligibleUpkeeps_CoverBuddy(t *testing.T) {
	t.Parallel()
	db, config, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db, config).Eth()

	registry, _ := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 1, 2, 20)

	for i := 0; i < 100; i++ {
		cltest.MustInsertUpkeepForRegistry(t, db, config, registry)
	}

	cltest.AssertCount(t, db, "keeper_registries", 1)
	cltest.AssertCount(t, db, "upkeep_registrations", 100)

	upkeep := keeper.UpkeepRegistration{}
	binaryHash := fmt.Sprintf("%b", utils.NewHash().Big())
	listBefore, err := orm.NewEligibleUpkeepsForRegistry(registry.ContractAddress, 21, 100, binaryHash) // normal
	require.NoError(t, err)
	require.NoError(t, db.Get(&upkeep, `UPDATE upkeep_registrations SET last_keeper_index = 0 RETURNING *`))
	listAfter, err := orm.NewEligibleUpkeepsForRegistry(registry.ContractAddress, 21, 100, binaryHash) // covering buddy
	require.NoError(t, err)
	require.Greater(t, len(listAfter), len(listBefore), "after our buddy runs all the performs we should have more eligible then a normal turn")
}

func TestKeeperDB_NewEligibleUpkeeps_FirstTurn(t *testing.T) {
	t.Parallel()
	db, config, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db, config).Eth()

	registry, _ := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 2, 20)

	for i := 0; i < 100; i++ {
		cltest.MustInsertUpkeepForRegistry(t, db, config, registry)
	}

	cltest.AssertCount(t, db, "keeper_registries", 1)
	cltest.AssertCount(t, db, "upkeep_registrations", 100)

	binaryHash := fmt.Sprintf("%b", utils.NewHash().Big())
	// last keeper index is null to simulate a normal first run
	listKpr0, err := orm.NewEligibleUpkeepsForRegistry(registry.ContractAddress, 21, 100, binaryHash) // someone eligible only kpr0 turn
	require.NoError(t, err)
	require.NotEqual(t, 0, len(listKpr0), "kpr0 should have some eligible as a normal turn")
}

func TestKeeperDB_NewEligibleUpkeeps_FiltersByRegistry(t *testing.T) {
	t.Parallel()
	db, config, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db, config).Eth()

	registry1, _ := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 1, 20)
	registry2, _ := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 1, 20)

	cltest.MustInsertUpkeepForRegistry(t, db, config, registry1)
	cltest.MustInsertUpkeepForRegistry(t, db, config, registry2)

	cltest.AssertCount(t, db, "keeper_registries", 2)
	cltest.AssertCount(t, db, "upkeep_registrations", 2)

	binaryHash := fmt.Sprintf("%b", utils.NewHash().Big())
	list1, err := orm.NewEligibleUpkeepsForRegistry(registry1.ContractAddress, 20, 100, binaryHash)
	require.NoError(t, err)
	list2, err := orm.NewEligibleUpkeepsForRegistry(registry2.ContractAddress, 20, 100, binaryHash)
	require.NoError(t, err)

	assert.Equal(t, 1, len(list1))
	assert.Equal(t, 1, len(list2))
}

func TestKeeperDB_NewSetLastRunInfoForUpkeepOnJob(t *testing.T) {
	t.Parallel()
	db, config, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db, config).Eth()

	registry, j := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 1, 20)
	upkeep := cltest.MustInsertUpkeepForRegistry(t, db, config, registry)

	// update
	require.NoError(t, orm.SetLastRunInfoForUpkeepOnJob(j.ID, upkeep.UpkeepID, 100, registry.FromAddress))
	assertLastRunHeight(t, db, upkeep, 100, 0)
	// update to lower block not allowed
	require.NoError(t, orm.SetLastRunInfoForUpkeepOnJob(j.ID, upkeep.UpkeepID, 0, registry.FromAddress))
	assertLastRunHeight(t, db, upkeep, 100, 0)
	// update to higher block allowed
	require.NoError(t, orm.SetLastRunInfoForUpkeepOnJob(j.ID, upkeep.UpkeepID, 101, registry.FromAddress))
	assertLastRunHeight(t, db, upkeep, 101, 0)
}
