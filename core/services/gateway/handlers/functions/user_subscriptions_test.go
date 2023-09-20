package functions_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_router"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/stretchr/testify/assert"
)

func TestUserSubscriptions(t *testing.T) {
	t.Parallel()

	us := functions.NewUserSubscriptions()

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

		us.UpdateSubscription(5, &functions_router.IFunctionsSubscriptionsSubscription{
			Owner:   user1,
			Balance: user1Balance,
		})
		us.UpdateSubscription(3, &functions_router.IFunctionsSubscriptionsSubscription{
			Owner:   user2,
			Balance: user2Balance1,
		})
		us.UpdateSubscription(10, &functions_router.IFunctionsSubscriptionsSubscription{
			Owner:   user2,
			Balance: user2Balance2,
		})

		balance, err := us.GetMaxUserBalance(user1)
		assert.NoError(t, err)
		assert.Zero(t, balance.Cmp(user1Balance))

		balance, err = us.GetMaxUserBalance(user2)
		assert.NoError(t, err)
		assert.Zero(t, balance.Cmp(user2Balance2))
	})

	t.Run("UpdateSubscription to remove subscriptions", func(t *testing.T) {
		user2 := utils.RandomAddress()
		user2Balance1 := big.NewInt(50)
		user2Balance2 := big.NewInt(70)

		us.UpdateSubscription(3, &functions_router.IFunctionsSubscriptionsSubscription{
			Owner:   user2,
			Balance: user2Balance1,
		})
		us.UpdateSubscription(10, &functions_router.IFunctionsSubscriptionsSubscription{
			Owner:   user2,
			Balance: user2Balance2,
		})

		us.UpdateSubscription(3, &functions_router.IFunctionsSubscriptionsSubscription{
			Owner: utils.ZeroAddress,
		})
		us.UpdateSubscription(10, &functions_router.IFunctionsSubscriptionsSubscription{
			Owner: utils.ZeroAddress,
		})

		_, err := us.GetMaxUserBalance(user2)
		assert.Error(t, err)
	})
}
