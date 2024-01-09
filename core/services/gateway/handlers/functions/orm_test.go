package functions_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_router"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions"
)

var (
	defaultFlags = [32]byte{0x1, 0x2, 0x3}
)

func setupORM(t *testing.T) (functions.ORM, error) {
	t.Helper()

	var (
		db   = pgtest.NewSqlxDB(t)
		lggr = logger.TestLogger(t)
	)

	return functions.NewORM(db, lggr, pgtest.NewQConfig(true), testutils.NewAddress())
}

func createSubscriptions(t *testing.T, orm functions.ORM, amount int) []functions.CachedSubscription {
	cachedSubscriptions := make([]functions.CachedSubscription, 0)
	for i := amount; i > 0; i-- {
		cs := functions.CachedSubscription{
			SubscriptionID: uint64(i),
			IFunctionsSubscriptionsSubscription: functions_router.IFunctionsSubscriptionsSubscription{
				Balance:        big.NewInt(10),
				Owner:          testutils.NewAddress(),
				BlockedBalance: big.NewInt(20),
				ProposedOwner:  common.Address{},
				Consumers:      []common.Address{},
				Flags:          defaultFlags,
			},
		}
		cachedSubscriptions = append(cachedSubscriptions, cs)
		err := orm.UpsertSubscription(cs)
		require.NoError(t, err)
	}
	return cachedSubscriptions
}

func TestORM_GetSubscriptions(t *testing.T) {
	t.Parallel()
	t.Run("fetch first page", func(t *testing.T) {
		orm, err := setupORM(t)
		require.NoError(t, err)
		cachedSubscriptions := createSubscriptions(t, orm, 2)
		results, err := orm.GetSubscriptions(0, 1)
		require.NoError(t, err)
		require.Equal(t, 1, len(results), "incorrect results length")
		require.Equal(t, cachedSubscriptions[1], results[0])
	})

	t.Run("fetch second page", func(t *testing.T) {
		orm, err := setupORM(t)
		require.NoError(t, err)
		cachedSubscriptions := createSubscriptions(t, orm, 2)
		results, err := orm.GetSubscriptions(1, 5)
		require.NoError(t, err)
		require.Equal(t, 1, len(results), "incorrect results length")
		require.Equal(t, cachedSubscriptions[0], results[0])
	})
}

func TestORM_UpsertSubscription(t *testing.T) {
	t.Parallel()

	t.Run("create a subscription", func(t *testing.T) {
		orm, err := setupORM(t)
		require.NoError(t, err)
		expected := functions.CachedSubscription{
			SubscriptionID: uint64(1),
			IFunctionsSubscriptionsSubscription: functions_router.IFunctionsSubscriptionsSubscription{
				Balance:        big.NewInt(10),
				Owner:          testutils.NewAddress(),
				BlockedBalance: big.NewInt(20),
				ProposedOwner:  common.Address{},
				Consumers:      []common.Address{},
				Flags:          defaultFlags,
			},
		}
		err = orm.UpsertSubscription(expected)
		require.NoError(t, err)

		results, err := orm.GetSubscriptions(0, 1)
		require.NoError(t, err)
		require.Equal(t, 1, len(results), "incorrect results length")
		require.Equal(t, expected, results[0])
	})

	t.Run("update a subscription", func(t *testing.T) {
		orm, err := setupORM(t)
		require.NoError(t, err)

		expectedUpdated := functions.CachedSubscription{
			SubscriptionID: uint64(1),
			IFunctionsSubscriptionsSubscription: functions_router.IFunctionsSubscriptionsSubscription{
				Balance:        big.NewInt(10),
				Owner:          testutils.NewAddress(),
				BlockedBalance: big.NewInt(20),
				ProposedOwner:  common.Address{},
				Consumers:      []common.Address{},
				Flags:          defaultFlags,
			},
		}
		err = orm.UpsertSubscription(expectedUpdated)
		require.NoError(t, err)

		expectedNotUpdated := functions.CachedSubscription{
			SubscriptionID: uint64(2),
			IFunctionsSubscriptionsSubscription: functions_router.IFunctionsSubscriptionsSubscription{
				Balance:        big.NewInt(10),
				Owner:          testutils.NewAddress(),
				BlockedBalance: big.NewInt(20),
				ProposedOwner:  common.Address{},
				Consumers:      []common.Address{},
				Flags:          defaultFlags,
			},
		}
		err = orm.UpsertSubscription(expectedNotUpdated)
		require.NoError(t, err)

		// update the balance value
		expectedUpdated.Balance = big.NewInt(20)
		err = orm.UpsertSubscription(expectedUpdated)
		require.NoError(t, err)

		results, err := orm.GetSubscriptions(0, 5)
		require.NoError(t, err)
		require.Equal(t, 2, len(results), "incorrect results length")
		require.Equal(t, expectedNotUpdated, results[1])
		require.Equal(t, expectedUpdated, results[0])
	})

	t.Run("update a deleted subscription", func(t *testing.T) {
		orm, err := setupORM(t)
		require.NoError(t, err)

		subscription := functions.CachedSubscription{
			SubscriptionID: uint64(1),
			IFunctionsSubscriptionsSubscription: functions_router.IFunctionsSubscriptionsSubscription{
				Balance:        big.NewInt(10),
				Owner:          testutils.NewAddress(),
				BlockedBalance: big.NewInt(20),
				ProposedOwner:  common.Address{},
				Consumers:      []common.Address{},
				Flags:          defaultFlags,
			},
		}
		err = orm.UpsertSubscription(subscription)
		require.NoError(t, err)

		// empty subscription
		subscription.IFunctionsSubscriptionsSubscription = functions_router.IFunctionsSubscriptionsSubscription{
			Balance:        big.NewInt(0),
			Owner:          common.Address{},
			BlockedBalance: big.NewInt(0),
			ProposedOwner:  common.Address{},
			Consumers:      []common.Address{},
			Flags:          [32]byte{},
		}

		err = orm.UpsertSubscription(subscription)
		require.NoError(t, err)

		results, err := orm.GetSubscriptions(0, 5)
		require.NoError(t, err)
		require.Equal(t, 1, len(results), "incorrect results length")
		require.Equal(t, subscription, results[0])
	})

	t.Run("create a subscription with same id but different router address", func(t *testing.T) {
		var (
			db   = pgtest.NewSqlxDB(t)
			lggr = logger.TestLogger(t)
		)

		orm1, err := functions.NewORM(db, lggr, pgtest.NewQConfig(true), testutils.NewAddress())
		require.NoError(t, err)
		orm2, err := functions.NewORM(db, lggr, pgtest.NewQConfig(true), testutils.NewAddress())
		require.NoError(t, err)

		subscription := functions.CachedSubscription{
			SubscriptionID: uint64(1),
			IFunctionsSubscriptionsSubscription: functions_router.IFunctionsSubscriptionsSubscription{
				Balance:        assets.Ether(10).ToInt(),
				Owner:          testutils.NewAddress(),
				BlockedBalance: assets.Ether(20).ToInt(),
				ProposedOwner:  common.Address{},
				Consumers:      []common.Address{},
				Flags:          defaultFlags,
			},
		}

		err = orm1.UpsertSubscription(subscription)
		require.NoError(t, err)

		// should update the existing subscription
		subscription.Balance = assets.Ether(12).ToInt()
		err = orm1.UpsertSubscription(subscription)
		require.NoError(t, err)

		results, err := orm1.GetSubscriptions(0, 10)
		require.NoError(t, err)
		require.Equal(t, 1, len(results), "incorrect results length")

		// should create a new subscription because it comes from a different router contract
		err = orm2.UpsertSubscription(subscription)
		require.NoError(t, err)

		results, err = orm1.GetSubscriptions(0, 10)
		require.NoError(t, err)
		require.Equal(t, 1, len(results), "incorrect results length")

		results, err = orm2.GetSubscriptions(0, 10)
		require.NoError(t, err)
		require.Equal(t, 1, len(results), "incorrect results length")
	})
}

func Test_NewORM(t *testing.T) {
	t.Run("OK-create_ORM", func(t *testing.T) {
		_, err := functions.NewORM(pgtest.NewSqlxDB(t), logger.TestLogger(t), pgtest.NewQConfig(true), testutils.NewAddress())
		require.NoError(t, err)
	})
	t.Run("NOK-create_ORM_with_nil_fields", func(t *testing.T) {
		_, err := functions.NewORM(nil, nil, nil, common.Address{})
		require.Error(t, err)
	})
	t.Run("NOK-create_ORM_with_empty_address", func(t *testing.T) {
		_, err := functions.NewORM(pgtest.NewSqlxDB(t), logger.TestLogger(t), pgtest.NewQConfig(true), common.Address{})
		require.Error(t, err)
	})
}
