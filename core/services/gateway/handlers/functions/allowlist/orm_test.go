package allowlist_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions/allowlist"
)

func setupORM(t *testing.T) (allowlist.ORM, error) {
	t.Helper()

	var (
		db   = pgtest.NewSqlxDB(t)
		lggr = logger.TestLogger(t)
	)

	return allowlist.NewORM(db, lggr, testutils.NewAddress())
}

func seedAllowedSenders(t *testing.T, orm allowlist.ORM, amount int) []common.Address {
	ctx := testutils.Context(t)
	storedAllowedSenders := make([]common.Address, amount)
	for i := 0; i < amount; i++ {
		address := testutils.NewAddress()
		storedAllowedSenders[i] = address
	}

	err := orm.CreateAllowedSenders(ctx, storedAllowedSenders)
	require.NoError(t, err)

	return storedAllowedSenders
}
func TestORM_GetAllowedSenders(t *testing.T) {
	t.Parallel()
	t.Run("fetch first page", func(t *testing.T) {
		ctx := testutils.Context(t)
		orm, err := setupORM(t)
		require.NoError(t, err)
		storedAllowedSenders := seedAllowedSenders(t, orm, 2)
		results, err := orm.GetAllowedSenders(ctx, 0, 1)
		require.NoError(t, err)
		require.Equal(t, 1, len(results), "incorrect results length")
		require.Equal(t, storedAllowedSenders[0], results[0])
	})

	t.Run("fetch second page", func(t *testing.T) {
		ctx := testutils.Context(t)
		orm, err := setupORM(t)
		require.NoError(t, err)
		storedAllowedSenders := seedAllowedSenders(t, orm, 2)
		results, err := orm.GetAllowedSenders(ctx, 1, 5)
		require.NoError(t, err)
		require.Equal(t, 1, len(results), "incorrect results length")
		require.Equal(t, storedAllowedSenders[1], results[0])
	})
}

func TestORM_CreateAllowedSenders(t *testing.T) {
	t.Parallel()

	t.Run("OK-create_an_allowed_sender", func(t *testing.T) {
		ctx := testutils.Context(t)
		orm, err := setupORM(t)
		require.NoError(t, err)
		expected := testutils.NewAddress()
		err = orm.CreateAllowedSenders(ctx, []common.Address{expected})
		require.NoError(t, err)

		results, err := orm.GetAllowedSenders(ctx, 0, 1)
		require.NoError(t, err)
		require.Equal(t, 1, len(results), "incorrect results length")
		require.Equal(t, expected, results[0])
	})

	t.Run("OK-create_an_existing_allowed_sender", func(t *testing.T) {
		ctx := testutils.Context(t)
		orm, err := setupORM(t)
		require.NoError(t, err)
		expected := testutils.NewAddress()
		err = orm.CreateAllowedSenders(ctx, []common.Address{expected})
		require.NoError(t, err)

		err = orm.CreateAllowedSenders(ctx, []common.Address{expected})
		require.NoError(t, err)

		results, err := orm.GetAllowedSenders(ctx, 0, 5)
		require.NoError(t, err)
		require.Equal(t, 1, len(results), "incorrect results length")
		require.Equal(t, expected, results[0])
	})

	t.Run("OK-create_multiple_allowed_senders_in_one_query", func(t *testing.T) {
		ctx := testutils.Context(t)
		orm, err := setupORM(t)
		require.NoError(t, err)
		expected := []common.Address{testutils.NewAddress(), testutils.NewAddress()}
		err = orm.CreateAllowedSenders(ctx, expected)
		require.NoError(t, err)

		results, err := orm.GetAllowedSenders(ctx, 0, 2)
		require.NoError(t, err)
		require.Equal(t, 2, len(results), "incorrect results length")
		require.Equal(t, expected[0], results[0])
		require.Equal(t, expected[1], results[1])
	})

	t.Run("OK-create_multiple_allowed_senders_with_duplicates", func(t *testing.T) {
		ctx := testutils.Context(t)
		orm, err := setupORM(t)
		require.NoError(t, err)
		addr1 := testutils.NewAddress()
		addr2 := testutils.NewAddress()
		expected := []common.Address{addr1, addr2}

		duplicatedAddressInput := []common.Address{addr1, addr1, addr1, addr2}
		err = orm.CreateAllowedSenders(ctx, duplicatedAddressInput)
		require.NoError(t, err)

		results, err := orm.GetAllowedSenders(ctx, 0, 10)
		require.NoError(t, err)
		require.Equal(t, 2, len(results), "incorrect results length")
		require.Equal(t, expected[0], results[0])
		require.Equal(t, expected[1], results[1])
	})

	// this scenario can happen if the allowlist is empty but we call CreateAllowedSenders
	t.Run("OK-empty_list", func(t *testing.T) {
		ctx := testutils.Context(t)
		orm, err := setupORM(t)
		require.NoError(t, err)
		err = orm.CreateAllowedSenders(ctx, []common.Address{})
		require.NoError(t, err)
	})
}

func TestORM_DeleteAllowedSenders(t *testing.T) {
	t.Parallel()

	t.Run("OK-delete_blocked_sender_from_allowed_list", func(t *testing.T) {
		ctx := testutils.Context(t)
		orm, err := setupORM(t)
		require.NoError(t, err)
		add1 := testutils.NewAddress()
		add2 := testutils.NewAddress()
		add3 := testutils.NewAddress()
		err = orm.CreateAllowedSenders(ctx, []common.Address{add1, add2, add3})
		require.NoError(t, err)

		results, err := orm.GetAllowedSenders(ctx, 0, 10)
		require.NoError(t, err)
		require.Equal(t, 3, len(results), "incorrect results length")
		require.Equal(t, add1, results[0])

		err = orm.DeleteAllowedSenders(ctx, []common.Address{add1, add3})
		require.NoError(t, err)

		results, err = orm.GetAllowedSenders(ctx, 0, 10)
		require.NoError(t, err)
		require.Equal(t, 1, len(results), "incorrect results length")
		require.Equal(t, add2, results[0])
	})

	t.Run("OK-delete_non_existing_blocked_sender_from_allowed_list", func(t *testing.T) {
		ctx := testutils.Context(t)
		orm, err := setupORM(t)
		require.NoError(t, err)
		add1 := testutils.NewAddress()
		add2 := testutils.NewAddress()
		err = orm.CreateAllowedSenders(ctx, []common.Address{add1, add2})
		require.NoError(t, err)

		results, err := orm.GetAllowedSenders(ctx, 0, 10)
		require.NoError(t, err)
		require.Equal(t, 2, len(results), "incorrect results length")
		require.Equal(t, add1, results[0])

		add3 := testutils.NewAddress()
		err = orm.DeleteAllowedSenders(ctx, []common.Address{add3})
		require.NoError(t, err)

		results, err = orm.GetAllowedSenders(ctx, 0, 10)
		require.NoError(t, err)
		require.Equal(t, 2, len(results), "incorrect results length")
		require.Equal(t, add1, results[0])
		require.Equal(t, add2, results[1])
	})
}

func TestORM_PurgeAllowedSenders(t *testing.T) {
	t.Parallel()

	t.Run("OK-purge_allowed_list", func(t *testing.T) {
		ctx := testutils.Context(t)
		orm, err := setupORM(t)
		require.NoError(t, err)
		add1 := testutils.NewAddress()
		add2 := testutils.NewAddress()
		add3 := testutils.NewAddress()
		err = orm.CreateAllowedSenders(ctx, []common.Address{add1, add2, add3})
		require.NoError(t, err)

		results, err := orm.GetAllowedSenders(ctx, 0, 10)
		require.NoError(t, err)
		require.Equal(t, 3, len(results), "incorrect results length")
		require.Equal(t, add1, results[0])

		err = orm.PurgeAllowedSenders(ctx)
		require.NoError(t, err)

		results, err = orm.GetAllowedSenders(ctx, 0, 10)
		require.NoError(t, err)
		require.Equal(t, 0, len(results), "incorrect results length")
	})

	t.Run("OK-purge_allowed_list_for_contract_address", func(t *testing.T) {
		ctx := testutils.Context(t)
		orm1, err := setupORM(t)
		require.NoError(t, err)
		add1 := testutils.NewAddress()
		add2 := testutils.NewAddress()
		err = orm1.CreateAllowedSenders(ctx, []common.Address{add1, add2})
		require.NoError(t, err)

		results, err := orm1.GetAllowedSenders(ctx, 0, 10)
		require.NoError(t, err)
		require.Equal(t, 2, len(results), "incorrect results length")
		require.Equal(t, add1, results[0])

		orm2, err := setupORM(t)
		require.NoError(t, err)
		add3 := testutils.NewAddress()
		add4 := testutils.NewAddress()
		err = orm2.CreateAllowedSenders(ctx, []common.Address{add3, add4})
		require.NoError(t, err)

		results, err = orm2.GetAllowedSenders(ctx, 0, 10)
		require.NoError(t, err)
		require.Equal(t, 2, len(results), "incorrect results length")
		require.Equal(t, add3, results[0])

		err = orm2.PurgeAllowedSenders(ctx)
		require.NoError(t, err)

		results, err = orm2.GetAllowedSenders(ctx, 0, 10)
		require.NoError(t, err)
		require.Equal(t, 0, len(results), "incorrect results length")

		results, err = orm1.GetAllowedSenders(ctx, 0, 10)
		require.NoError(t, err)
		require.Equal(t, 2, len(results), "incorrect results length")
		require.Equal(t, add1, results[0])
		require.Equal(t, add2, results[1])
	})
}

func Test_NewORM(t *testing.T) {
	t.Run("OK-create_ORM", func(t *testing.T) {
		_, err := allowlist.NewORM(pgtest.NewSqlxDB(t), logger.TestLogger(t), testutils.NewAddress())
		require.NoError(t, err)
	})
	t.Run("NOK-create_ORM_with_nil_fields", func(t *testing.T) {
		_, err := allowlist.NewORM(nil, nil, common.Address{})
		require.Error(t, err)
	})
	t.Run("NOK-create_ORM_with_empty_address", func(t *testing.T) {
		_, err := allowlist.NewORM(pgtest.NewSqlxDB(t), logger.TestLogger(t), common.Address{})
		require.Error(t, err)
	})
}
