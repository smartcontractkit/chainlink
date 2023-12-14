package functions

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_router"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// Methods are NOT thread-safe.

type UserSubscriptions interface {
	UpdateSubscription(subscriptionId uint64, subscription *functions_router.IFunctionsSubscriptionsSubscription)
	GetMaxUserBalance(user common.Address) (*big.Int, error)
}

type userSubscriptions struct {
	userSubscriptionsMap map[common.Address]map[uint64]*functions_router.IFunctionsSubscriptionsSubscription
	subscriptionIdsMap   map[uint64]common.Address
}

func NewUserSubscriptions() UserSubscriptions {
	return &userSubscriptions{
		userSubscriptionsMap: make(map[common.Address]map[uint64]*functions_router.IFunctionsSubscriptionsSubscription),
		subscriptionIdsMap:   make(map[uint64]common.Address),
	}
}

func (us *userSubscriptions) UpdateSubscription(subscriptionId uint64, subscription *functions_router.IFunctionsSubscriptionsSubscription) {
	if subscription == nil || subscription.Owner == utils.ZeroAddress {
		user, ok := us.subscriptionIdsMap[subscriptionId]
		if ok {
			delete(us.userSubscriptionsMap[user], subscriptionId)
			if len(us.userSubscriptionsMap[user]) == 0 {
				delete(us.userSubscriptionsMap, user)
			}
		}
		delete(us.subscriptionIdsMap, subscriptionId)
	} else {
		us.subscriptionIdsMap[subscriptionId] = subscription.Owner
		if _, ok := us.userSubscriptionsMap[subscription.Owner]; !ok {
			us.userSubscriptionsMap[subscription.Owner] = make(map[uint64]*functions_router.IFunctionsSubscriptionsSubscription)
		}
		us.userSubscriptionsMap[subscription.Owner][subscriptionId] = subscription
	}
}

func (us *userSubscriptions) GetMaxUserBalance(user common.Address) (*big.Int, error) {
	subs, exists := us.userSubscriptionsMap[user]
	if !exists {
		return nil, errors.New("user has no subscriptions")
	}

	maxBalance := big.NewInt(0)
	for _, sub := range subs {
		if sub.Balance.Cmp(maxBalance) > 0 {
			maxBalance = sub.Balance
		}
	}
	return maxBalance, nil
}
