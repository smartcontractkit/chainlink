package functions_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_router"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions"
	"github.com/stretchr/testify/require"
)

var (
	defaultFlags = [32]byte{0x1, 0x2, 0x3}
)

func setupORM(t *testing.T) functions.ORM {
	t.Helper()

	var (
		db   = pgtest.NewSqlxDB(t)
		lggr = logger.TestLogger(t)
		orm  = functions.NewORM(db, lggr, pgtest.NewQConfig(true))
	)

	return orm
}

func createSubscription(t *testing.T, orm functions.ORM, amount int) []functions.CachedSubscription {
	cachedSubscriptions := make([]functions.CachedSubscription, 0)
	for i := amount; i > 0; i-- {
		address := testutils.NewAddress()
		cs := functions.CachedSubscription{
			SubscriptionID: uint64(i),
			IFunctionsSubscriptionsSubscription: functions_router.IFunctionsSubscriptionsSubscription{
				Balance:        assets.Ether(10).ToInt(),
				Owner:          address,
				BlockedBalance: assets.Ether(20).ToInt(),
				ProposedOwner:  common.Address{},
				Consumers:      []common.Address{},
				Flags:          defaultFlags,
			},
		}
		cachedSubscriptions = append(cachedSubscriptions, cs)
		err := orm.CreateSubscription(cs)
		require.NoError(t, err)
	}
	return cachedSubscriptions
}

func TestORM_FetchSubscriptions(t *testing.T) {
	t.Parallel()
	t.Run("fetch first page", func(t *testing.T) {
		orm := setupORM(t)
		cachedSubscriptions := createSubscription(t, orm, 2)
		results, err := orm.FetchSubscriptions(0, 1)
		require.NoError(t, err)
		require.Equal(t, 1, len(results), "incorrect results length")
		require.Equal(t, results[0].Owner, cachedSubscriptions[0].Owner)
	})

	t.Run("fetch second page", func(t *testing.T) {
		orm := setupORM(t)
		cachedSubscriptions := createSubscription(t, orm, 2)
		results, err := orm.FetchSubscriptions(1, 5)
		require.NoError(t, err)
		require.Equal(t, 1, len(results), "incorrect results length")
		require.Equal(t, results[0].Owner, cachedSubscriptions[1].Owner)
	})
}
