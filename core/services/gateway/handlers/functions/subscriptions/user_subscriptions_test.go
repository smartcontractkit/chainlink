package subscriptions_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink-integrations/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_router"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions/subscriptions"

	"github.com/stretchr/testify/assert"
)

func TestUserSubscriptions(t *testing.T) {
	t.Parallel()

	us := subscriptions.NewUserSubscriptions()

	t.Run("GetMaxUserBalance for unknown user", func(t *testing.T) {
		_, err := us.GetMaxUserBalance(utils.RandomAddress())
		assert.Error(t, err)
	})

	t.Run("UpdateSubscription then GetMaxUserBalance", func(t *testing.T) {
		user1 := utils.RandomAddress()
		user1Balance := big.NewInt(10)
		user2 := utils.RandomAddress()
		user2Balance1 := big.NewInt(50)
		user2Balance2 := big.NewInt(70)

		updated := us.UpdateSubscription(5, &functions_router.IFunctionsSubscriptionsSubscription{
			Owner:   user1,
			Balance: user1Balance,
		})
		assert.True(t, updated)

		updated = us.UpdateSubscription(3, &functions_router.IFunctionsSubscriptionsSubscription{
			Owner:   user2,
			Balance: user2Balance1,
		})
		assert.True(t, updated)

		updated = us.UpdateSubscription(10, &functions_router.IFunctionsSubscriptionsSubscription{
			Owner:   user2,
			Balance: user2Balance2,
		})
		assert.True(t, updated)

		balance, err := us.GetMaxUserBalance(user1)
		assert.NoError(t, err)
		assert.Zero(t, balance.Cmp(user1Balance))

		balance, err = us.GetMaxUserBalance(user2)
		assert.NoError(t, err)
		assert.Zero(t, balance.Cmp(user2Balance2))
	})
}

func TestUserSubscriptions_UpdateSubscription(t *testing.T) {
	t.Parallel()

	t.Run("update balance", func(t *testing.T) {
		us := subscriptions.NewUserSubscriptions()
		owner := utils.RandomAddress()

		updated := us.UpdateSubscription(1, &functions_router.IFunctionsSubscriptionsSubscription{
			Owner:   owner,
			Balance: big.NewInt(10),
		})
		assert.True(t, updated)

		updated = us.UpdateSubscription(1, &functions_router.IFunctionsSubscriptionsSubscription{
			Owner:   owner,
			Balance: big.NewInt(100),
		})
		assert.True(t, updated)
	})

	t.Run("updated proposed owner", func(t *testing.T) {
		us := subscriptions.NewUserSubscriptions()
		owner := utils.RandomAddress()

		updated := us.UpdateSubscription(1, &functions_router.IFunctionsSubscriptionsSubscription{
			Owner:   owner,
			Balance: big.NewInt(10),
		})
		assert.True(t, updated)

		updated = us.UpdateSubscription(1, &functions_router.IFunctionsSubscriptionsSubscription{
			Owner:         owner,
			Balance:       big.NewInt(10),
			ProposedOwner: utils.RandomAddress(),
		})
		assert.True(t, updated)
	})
	t.Run("remove subscriptions", func(t *testing.T) {
		us := subscriptions.NewUserSubscriptions()
		user2 := utils.RandomAddress()
		user2Balance1 := big.NewInt(50)
		user2Balance2 := big.NewInt(70)

		updated := us.UpdateSubscription(3, &functions_router.IFunctionsSubscriptionsSubscription{
			Owner:   user2,
			Balance: user2Balance1,
		})
		assert.True(t, updated)

		updated = us.UpdateSubscription(10, &functions_router.IFunctionsSubscriptionsSubscription{
			Owner:   user2,
			Balance: user2Balance2,
		})
		assert.True(t, updated)

		updated = us.UpdateSubscription(3, &functions_router.IFunctionsSubscriptionsSubscription{
			Owner: utils.ZeroAddress,
		})
		assert.True(t, updated)

		updated = us.UpdateSubscription(10, &functions_router.IFunctionsSubscriptionsSubscription{
			Owner: utils.ZeroAddress,
		})
		assert.True(t, updated)

		_, err := us.GetMaxUserBalance(user2)
		assert.Error(t, err)
	})

	t.Run("remove a non existing subscription", func(t *testing.T) {
		us := subscriptions.NewUserSubscriptions()
		updated := us.UpdateSubscription(3, &functions_router.IFunctionsSubscriptionsSubscription{
			Owner: utils.ZeroAddress,
		})
		assert.False(t, updated)
	})

	t.Run("no actual changes", func(t *testing.T) {
		us := subscriptions.NewUserSubscriptions()
		subscription := functions_router.IFunctionsSubscriptionsSubscription{
			Owner:          utils.RandomAddress(),
			Balance:        big.NewInt(25),
			BlockedBalance: big.NewInt(25),
		}
		identicalSubscription := subscription
		updated := us.UpdateSubscription(5, &subscription)
		assert.True(t, updated)

		updated = us.UpdateSubscription(5, &identicalSubscription)
		assert.False(t, updated)
	})
}
