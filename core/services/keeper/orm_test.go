package keeper_test

import (
	"fmt"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jmoiron/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	evmutils "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keeper"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	bigmath "github.com/smartcontractkit/chainlink/v2/core/utils/big_math"
)

var (
	checkData  = common.Hex2Bytes("ABC123")
	executeGas = uint32(10_000)
)

func setupKeeperDB(t *testing.T) (
	*sqlx.DB,
	chainlink.GeneralConfig,
	*keeper.ORM,
) {
	cfg := configtest.NewGeneralConfig(t, nil)
	db := pgtest.NewSqlxDB(t)
	orm := keeper.NewORM(db, logger.TestLogger(t))
	return db, cfg, orm
}

func newUpkeep(registry keeper.Registry, upkeepID int64) keeper.UpkeepRegistration {
	return keeper.UpkeepRegistration{
		UpkeepID:   ubig.NewI(upkeepID),
		ExecuteGas: executeGas,
		Registry:   registry,
		RegistryID: registry.ID,
		CheckData:  checkData,
	}
}

func waitLastRunHeight(t *testing.T, db *sqlx.DB, upkeep keeper.UpkeepRegistration, height int64) {
	t.Helper()

	gomega.NewWithT(t).Eventually(func() int64 {
		err := db.Get(&upkeep, `SELECT * FROM upkeep_registrations WHERE id = $1`, upkeep.ID)
		require.NoError(t, err)
		return upkeep.LastRunBlockHeight
	}, time.Second*2, time.Millisecond*100).Should(gomega.Equal(height))
}

func assertLastRunHeight(t *testing.T, db *sqlx.DB, upkeep keeper.UpkeepRegistration, lastRunBlockHeight int64, lastKeeperIndex int64) {
	err := db.Get(&upkeep, `SELECT * FROM upkeep_registrations WHERE id = $1`, upkeep.ID)
	require.NoError(t, err)
	require.Equal(t, lastRunBlockHeight, upkeep.LastRunBlockHeight)
	require.Equal(t, lastKeeperIndex, upkeep.LastKeeperIndex.Int64)
}

func TestKeeperDB_Registries(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	db, _, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 1, 20)
	cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 1, 20)

	existingRegistries, err := orm.Registries(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(existingRegistries))
}

func TestKeeperDB_RegistryByContractAddress(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	db, _, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	registry, _ := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 1, 20)
	cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 1, 20)

	registryByContractAddress, err := orm.RegistryByContractAddress(ctx, registry.ContractAddress)
	require.NoError(t, err)
	require.Equal(t, registry, registryByContractAddress)
}

func TestKeeperDB_UpsertUpkeep(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	db, _, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	registry, _ := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 1, 20)
	upkeep := keeper.UpkeepRegistration{
		UpkeepID:            ubig.NewI(0),
		ExecuteGas:          executeGas,
		Registry:            registry,
		RegistryID:          registry.ID,
		CheckData:           checkData,
		LastRunBlockHeight:  1,
		PositioningConstant: 1,
	}
	require.NoError(t, orm.UpsertUpkeep(ctx, &upkeep))
	cltest.AssertCount(t, db, "upkeep_registrations", 1)

	// update upkeep
	upkeep.ExecuteGas = 20_000
	upkeep.CheckData = common.Hex2Bytes("8888")
	upkeep.LastRunBlockHeight = 2

	err := orm.UpsertUpkeep(ctx, &upkeep)
	require.NoError(t, err)
	cltest.AssertCount(t, db, "upkeep_registrations", 1)

	var upkeepFromDB keeper.UpkeepRegistration
	err = db.Get(&upkeepFromDB, `SELECT * FROM upkeep_registrations ORDER BY id LIMIT 1`)
	require.NoError(t, err)
	require.Equal(t, uint32(20_000), upkeepFromDB.ExecuteGas)
	require.Equal(t, "8888", common.Bytes2Hex(upkeepFromDB.CheckData))
	require.Equal(t, int64(1), upkeepFromDB.LastRunBlockHeight) // shouldn't change on upsert
}

func TestKeeperDB_BatchDeleteUpkeepsForJob(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	db, _, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	registry, job := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 1, 20)

	expectedUpkeepID := cltest.MustInsertUpkeepForRegistry(t, db, registry).UpkeepID
	var upkeepIDs []ubig.Big
	for i := 0; i < 2; i++ {
		upkeep := cltest.MustInsertUpkeepForRegistry(t, db, registry)
		upkeepIDs = append(upkeepIDs, *upkeep.UpkeepID)
	}

	cltest.AssertCount(t, db, "upkeep_registrations", 3)

	_, err := orm.BatchDeleteUpkeepsForJob(ctx, job.ID, upkeepIDs)
	require.NoError(t, err)
	cltest.AssertCount(t, db, "upkeep_registrations", 1)

	var remainingUpkeep keeper.UpkeepRegistration
	err = db.Get(&remainingUpkeep, `SELECT * FROM upkeep_registrations ORDER BY id LIMIT 1`)
	require.NoError(t, err)
	require.Equal(t, expectedUpkeepID, remainingUpkeep.UpkeepID)
}

func TestKeeperDB_EligibleUpkeeps_Shuffle(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	db, _, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	blockheight := int64(63)
	gracePeriod := int64(10)

	registry, _ := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 1, 20)

	ordered := [100]int64{}
	for i := 0; i < 100; i++ {
		k := newUpkeep(registry, int64(i))
		ordered[i] = int64(i)
		err := orm.UpsertUpkeep(ctx, &k)
		require.NoError(t, err)
	}
	cltest.AssertCount(t, db, "upkeep_registrations", 100)

	eligibleUpkeeps, err := orm.EligibleUpkeepsForRegistry(ctx, registry.ContractAddress, blockheight, gracePeriod, fmt.Sprintf("%b", evmutils.NewHash().Big()))
	assert.NoError(t, err)

	require.Len(t, eligibleUpkeeps, 100)
	shuffled := [100]*ubig.Big{}
	for i := 0; i < 100; i++ {
		shuffled[i] = eligibleUpkeeps[i].UpkeepID
	}
	assert.NotEqualValues(t, ordered, shuffled)
}

func TestKeeperDB_NewEligibleUpkeeps_GracePeriod(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	db, _, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	registry, _ := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 2, 20)

	for i := 0; i < 100; i++ {
		cltest.MustInsertUpkeepForRegistry(t, db, registry)
	}

	cltest.AssertCount(t, db, "keeper_registries", 1)
	cltest.AssertCount(t, db, "upkeep_registrations", 100)

	// if current keeper index = 0 and all upkeeps last perform was done by index = 0 and still within grace period
	upkeep := keeper.UpkeepRegistration{}
	require.NoError(t, db.Get(&upkeep, `UPDATE upkeep_registrations SET last_keeper_index = 0, last_run_block_height = 10 RETURNING *`))
	list0, err := orm.EligibleUpkeepsForRegistry(ctx, registry.ContractAddress, 21, 100, fmt.Sprintf("%b", evmutils.NewHash().Big())) // none eligible
	require.NoError(t, err)
	require.Equal(t, 0, len(list0), "should be 0 as all last perform was done by current node")

	// once passed grace period
	list1, err := orm.EligibleUpkeepsForRegistry(ctx, registry.ContractAddress, 121, 100, fmt.Sprintf("%b", evmutils.NewHash().Big())) // none eligible
	require.NoError(t, err)
	require.NotEqual(t, 0, len(list1), "should get some eligible upkeeps now that they are outside grace period")
}

func TestKeeperDB_EligibleUpkeeps_TurnsRandom(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	db, _, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	registry, _ := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 3, 10)

	for i := 0; i < 1000; i++ {
		cltest.MustInsertUpkeepForRegistry(t, db, registry)
	}

	cltest.AssertCount(t, db, "keeper_registries", 1)
	cltest.AssertCount(t, db, "upkeep_registrations", 1000)

	// 3 keepers 10 block turns should be different every turn
	list1, err := orm.EligibleUpkeepsForRegistry(ctx, registry.ContractAddress, 20, 100, fmt.Sprintf("%b", evmutils.NewHash().Big()))
	require.NoError(t, err)
	list2, err := orm.EligibleUpkeepsForRegistry(ctx, registry.ContractAddress, 31, 100, fmt.Sprintf("%b", evmutils.NewHash().Big()))
	require.NoError(t, err)
	list3, err := orm.EligibleUpkeepsForRegistry(ctx, registry.ContractAddress, 42, 100, fmt.Sprintf("%b", evmutils.NewHash().Big()))
	require.NoError(t, err)
	list4, err := orm.EligibleUpkeepsForRegistry(ctx, registry.ContractAddress, 53, 100, fmt.Sprintf("%b", evmutils.NewHash().Big()))
	require.NoError(t, err)

	// sort before compare
	sort.Slice(list1, func(i, j int) bool {
		return list1[i].UpkeepID.Cmp(list1[j].UpkeepID) == -1
	})
	sort.Slice(list2, func(i, j int) bool {
		return list2[i].UpkeepID.Cmp(list2[j].UpkeepID) == -1
	})
	sort.Slice(list3, func(i, j int) bool {
		return list3[i].UpkeepID.Cmp(list3[j].UpkeepID) == -1
	})
	sort.Slice(list4, func(i, j int) bool {
		return list4[i].UpkeepID.Cmp(list4[j].UpkeepID) == -1
	})

	assert.NotEqual(t, list1, list2, "list1 vs list2")
	assert.NotEqual(t, list1, list3, "list1 vs list3")
	assert.NotEqual(t, list1, list4, "list1 vs list4")
}

func TestKeeperDB_NewEligibleUpkeeps_SkipIfLastPerformedByCurrentKeeper(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	db, _, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	registry, _ := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 2, 20)

	for i := 0; i < 100; i++ {
		cltest.MustInsertUpkeepForRegistry(t, db, registry)
	}

	cltest.AssertCount(t, db, "keeper_registries", 1)
	cltest.AssertCount(t, db, "upkeep_registrations", 100)

	// if current keeper index = 0 and all upkeeps last perform was done by index = 0 then skip as it would not pass required turn taking
	upkeep := keeper.UpkeepRegistration{}
	require.NoError(t, db.Get(&upkeep, `UPDATE upkeep_registrations SET last_keeper_index = 0 RETURNING *`))
	list0, err := orm.EligibleUpkeepsForRegistry(ctx, registry.ContractAddress, 21, 100, fmt.Sprintf("%b", evmutils.NewHash().Big())) // none eligible
	require.NoError(t, err)
	require.Equal(t, 0, len(list0), "should be 0 as all last perform was done by current node")
}

func TestKeeperDB_NewEligibleUpkeeps_CoverBuddy(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	db, _, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	registry, _ := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 1, 2, 20)

	for i := 0; i < 100; i++ {
		cltest.MustInsertUpkeepForRegistry(t, db, registry)
	}

	cltest.AssertCount(t, db, "keeper_registries", 1)
	cltest.AssertCount(t, db, "upkeep_registrations", 100)

	upkeep := keeper.UpkeepRegistration{}
	binaryHash := fmt.Sprintf("%b", evmutils.NewHash().Big())
	listBefore, err := orm.EligibleUpkeepsForRegistry(ctx, registry.ContractAddress, 21, 100, binaryHash) // normal
	require.NoError(t, err)
	require.NoError(t, db.Get(&upkeep, `UPDATE upkeep_registrations SET last_keeper_index = 0 RETURNING *`))
	listAfter, err := orm.EligibleUpkeepsForRegistry(ctx, registry.ContractAddress, 21, 100, binaryHash) // covering buddy
	require.NoError(t, err)
	require.Greater(t, len(listAfter), len(listBefore), "after our buddy runs all the performs we should have more eligible then a normal turn")
}

func TestKeeperDB_NewEligibleUpkeeps_FirstTurn(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	db, _, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	registry, _ := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 2, 20)

	for i := 0; i < 100; i++ {
		cltest.MustInsertUpkeepForRegistry(t, db, registry)
	}

	cltest.AssertCount(t, db, "keeper_registries", 1)
	cltest.AssertCount(t, db, "upkeep_registrations", 100)

	binaryHash := fmt.Sprintf("%b", evmutils.NewHash().Big())
	// last keeper index is null to simulate a normal first run
	listKpr0, err := orm.EligibleUpkeepsForRegistry(ctx, registry.ContractAddress, 21, 100, binaryHash) // someone eligible only kpr0 turn
	require.NoError(t, err)
	require.NotEqual(t, 0, len(listKpr0), "kpr0 should have some eligible as a normal turn")
}

func TestKeeperDB_NewEligibleUpkeeps_FiltersByRegistry(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	db, _, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	registry1, _ := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 1, 20)
	registry2, _ := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 1, 20)

	cltest.MustInsertUpkeepForRegistry(t, db, registry1)
	cltest.MustInsertUpkeepForRegistry(t, db, registry2)

	cltest.AssertCount(t, db, "keeper_registries", 2)
	cltest.AssertCount(t, db, "upkeep_registrations", 2)

	binaryHash := fmt.Sprintf("%b", evmutils.NewHash().Big())
	list1, err := orm.EligibleUpkeepsForRegistry(ctx, registry1.ContractAddress, 20, 100, binaryHash)
	require.NoError(t, err)
	list2, err := orm.EligibleUpkeepsForRegistry(ctx, registry2.ContractAddress, 20, 100, binaryHash)
	require.NoError(t, err)

	assert.Equal(t, 1, len(list1))
	assert.Equal(t, 1, len(list2))
}

func TestKeeperDB_AllUpkeepIDsForRegistry(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	db, _, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	registry, _ := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 1, 20)

	upkeepIDs, err := orm.AllUpkeepIDsForRegistry(ctx, registry.ID)
	require.NoError(t, err)
	// No upkeeps returned
	require.Len(t, upkeepIDs, 0)

	upkeep := newUpkeep(registry, 3)
	err = orm.UpsertUpkeep(ctx, &upkeep)
	require.NoError(t, err)

	upkeep = newUpkeep(registry, 8)
	err = orm.UpsertUpkeep(ctx, &upkeep)
	require.NoError(t, err)

	// We should get two upkeeps IDs, 3 & 8
	upkeepIDs, err = orm.AllUpkeepIDsForRegistry(ctx, registry.ID)
	require.NoError(t, err)
	// No upkeeps returned
	require.Len(t, upkeepIDs, 2)
	require.Contains(t, upkeepIDs, *ubig.New(big.NewInt(3)))
	require.Contains(t, upkeepIDs, *ubig.New(big.NewInt(8)))
}

func TestKeeperDB_UpdateUpkeepLastKeeperIndex(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	db, _, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	registry, j := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 1, 20)
	upkeep := cltest.MustInsertUpkeepForRegistry(t, db, registry)

	require.NoError(t, orm.UpdateUpkeepLastKeeperIndex(ctx, j.ID, upkeep.UpkeepID, registry.FromAddress))

	err := db.Get(&upkeep, `SELECT * FROM upkeep_registrations WHERE upkeep_id = $1`, upkeep.UpkeepID)
	require.NoError(t, err)
	require.Equal(t, int64(0), upkeep.LastKeeperIndex.Int64)
}

func TestKeeperDB_NewSetLastRunInfoForUpkeepOnJob(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	db, _, orm := setupKeeperDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	registry, j := cltest.MustInsertKeeperRegistry(t, db, orm, ethKeyStore, 0, 1, 20)
	upkeep := cltest.MustInsertUpkeepForRegistry(t, db, registry)
	registry.NumKeepers = 2
	registry.KeeperIndexMap = map[types.EIP55Address]int32{
		registry.FromAddress: 0,
		types.EIP55AddressFromAddress(evmutils.ZeroAddress): 1,
	}
	err := orm.UpsertRegistry(ctx, &registry)
	require.NoError(t, err, "UPDATE keeper_registries")

	// update
	rowsAffected, err := orm.SetLastRunInfoForUpkeepOnJob(ctx, j.ID, upkeep.UpkeepID, 100, registry.FromAddress)
	require.NoError(t, err)
	require.Equal(t, rowsAffected, int64(1))
	assertLastRunHeight(t, db, upkeep, 100, 0)
	// update to lower block height not allowed
	rowsAffected, err = orm.SetLastRunInfoForUpkeepOnJob(ctx, j.ID, upkeep.UpkeepID, 0, registry.FromAddress)
	require.NoError(t, err)
	require.Equal(t, rowsAffected, int64(0))
	assertLastRunHeight(t, db, upkeep, 100, 0)
	// update to same block height allowed
	rowsAffected, err = orm.SetLastRunInfoForUpkeepOnJob(ctx, j.ID, upkeep.UpkeepID, 100, types.EIP55AddressFromAddress(evmutils.ZeroAddress))
	require.NoError(t, err)
	require.Equal(t, rowsAffected, int64(1))
	assertLastRunHeight(t, db, upkeep, 100, 1)
	// update to higher block height allowed
	rowsAffected, err = orm.SetLastRunInfoForUpkeepOnJob(ctx, j.ID, upkeep.UpkeepID, 101, registry.FromAddress)
	require.NoError(t, err)
	require.Equal(t, rowsAffected, int64(1))
	assertLastRunHeight(t, db, upkeep, 101, 0)
}

func TestKeeperDB_LeastSignificant(t *testing.T) {
	t.Parallel()
	db, _, _ := setupKeeperDB(t)
	sql := `SELECT least_significant($1, $2)`
	inputBytes := "10001000101010101101"
	for _, test := range []struct {
		name           string
		inputLength    int
		expectedOutput string
		expectedError  bool
	}{
		{
			name:           "half slice",
			inputLength:    10,
			expectedOutput: "1010101101",
		},
		{
			name:           "full slice",
			inputLength:    20,
			expectedOutput: "10001000101010101101",
		},
		{
			name:           "no slice",
			inputLength:    0,
			expectedOutput: "",
		},
		{
			name:          "slice too large",
			inputLength:   21,
			expectedError: true,
		},
	} {
		t.Run(test.name, func(tt *testing.T) {
			var test = test
			var result string
			err := db.Get(&result, sql, inputBytes, test.inputLength)
			if test.expectedError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.expectedOutput, result)
		})
	}
}

func TestKeeperDB_Uint256ToBit(t *testing.T) {
	t.Parallel()
	db, _, _ := setupKeeperDB(t)
	sql := `SELECT uint256_to_bit($1)`
	for _, test := range []struct {
		name          string
		input         *big.Int
		errorExpected bool
	}{
		{
			name:  "min",
			input: big.NewInt(0),
		},
		{
			name:  "max",
			input: evmutils.MaxUint256,
		},
		{
			name:  "rand",
			input: evmutils.RandUint256(),
		},
		{
			name:  "needs pading",
			input: big.NewInt(1),
		},
		{
			name:          "overflow",
			input:         bigmath.Add(evmutils.MaxUint256, big.NewInt(1)),
			errorExpected: true,
		},
	} {
		t.Run(test.name, func(tt *testing.T) {
			var test = test
			var result string
			err := db.Get(&result, sql, test.input.String())
			if test.errorExpected {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			expected := utils.LeftPadBitString(fmt.Sprintf("%b", test.input), 256)
			require.Equal(t, expected, result)
		})
	}
}
