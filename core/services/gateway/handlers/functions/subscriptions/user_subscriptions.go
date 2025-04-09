package subscriptions

import (
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/functions/generated/functions_router"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
)

// Methods are NOT thread-safe.

var ErrUserHasNoSubscription = errors.New("user has no subscriptions")

type UserSubscriptions interface {
	UpdateSubscription(subscriptionId uint64, subscription *functions_router.IFunctionsSubscriptionsSubscription) bool
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

// StoredSubscription is used to populate the user subscription maps from a persistent layer like postgres.
type StoredSubscription struct {
	SubscriptionID uint64
	functions_router.IFunctionsSubscriptionsSubscription
}

// UpdateSubscription updates a subscription returning false in case there was no variation to the current state.
func (us *userSubscriptions) UpdateSubscription(subscriptionId uint64, subscription *functions_router.IFunctionsSubscriptionsSubscription) bool {
	if subscription == nil || subscription.Owner == utils.ZeroAddress {
		user, ok := us.subscriptionIdsMap[subscriptionId]
		if !ok {
			return false
		}

		delete(us.userSubscriptionsMap[user], subscriptionId)
		delete(us.subscriptionIdsMap, subscriptionId)
		if len(us.userSubscriptionsMap[user]) == 0 {
			delete(us.userSubscriptionsMap, user)
		}
		return true
	}

	// there is no change to the subscription
	if reflect.DeepEqual(us.userSubscriptionsMap[subscription.Owner][subscriptionId], subscription) {
		return false
	}

	us.subscriptionIdsMap[subscriptionId] = subscription.Owner
	if _, ok := us.userSubscriptionsMap[subscription.Owner]; !ok {
		us.userSubscriptionsMap[subscription.Owner] = make(map[uint64]*functions_router.IFunctionsSubscriptionsSubscription)
	}
	us.userSubscriptionsMap[subscription.Owner][subscriptionId] = subscription
	return true
}

func (us *userSubscriptions) GetMaxUserBalance(user common.Address) (*big.Int, error) {
	subs, exists := us.userSubscriptionsMap[user]
	if !exists {
		return nil, ErrUserHasNoSubscription
	}

	maxBalance := big.NewInt(0)
	for _, sub := range subs {
		if sub.Balance.Cmp(maxBalance) > 0 {
			maxBalance = sub.Balance
		}
	}
	return maxBalance, nil
}
