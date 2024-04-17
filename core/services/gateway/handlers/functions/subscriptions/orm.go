package subscriptions

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

//go:generate mockery --quiet --name ORM --output ./mocks/ --case=underscore
type ORM interface {
	GetSubscriptions(ctx context.Context, offset, limit uint) ([]StoredSubscription, error)
	UpsertSubscription(ctx context.Context, subscription StoredSubscription) error
}

type orm struct {
	ds                    sqlutil.DataSource
	lggr                  logger.Logger
	routerContractAddress common.Address
}

var _ ORM = (*orm)(nil)
var (
	ErrInvalidParameters = errors.New("invalid parameters provided to create a subscription contract ORM")
)

const (
	tableName = "functions_subscriptions"
)

type storedSubscriptionRow struct {
	SubscriptionID        uint64
	Owner                 common.Address
	Balance               int64
	BlockedBalance        int64
	ProposedOwner         common.Address
	Consumers             pq.ByteaArray
	Flags                 []uint8
	RouterContractAddress common.Address
}

func NewORM(ds sqlutil.DataSource, lggr logger.Logger, routerContractAddress common.Address) (ORM, error) {
	if ds == nil || lggr == nil || routerContractAddress == (common.Address{}) {
		return nil, ErrInvalidParameters
	}

	return &orm{
		ds:                    ds,
		lggr:                  lggr,
		routerContractAddress: routerContractAddress,
	}, nil
}

func (o *orm) GetSubscriptions(ctx context.Context, offset, limit uint) ([]StoredSubscription, error) {
	var storedSubscriptions []StoredSubscription
	var storedSubscriptionRows []storedSubscriptionRow
	stmt := fmt.Sprintf(`
		SELECT subscription_id, owner, balance, blocked_balance, proposed_owner, consumers, flags, router_contract_address
		FROM %s
		WHERE router_contract_address = $1
		ORDER BY subscription_id ASC
		OFFSET $2
		LIMIT $3;
	`, tableName)
	err := o.ds.SelectContext(ctx, &storedSubscriptionRows, stmt, o.routerContractAddress, offset, limit)
	if err != nil {
		return storedSubscriptions, err
	}

	for _, cs := range storedSubscriptionRows {
		storedSubscriptions = append(storedSubscriptions, cs.encode())
	}

	return storedSubscriptions, nil
}

// UpsertSubscription will update if a subscription exists or create if it does not.
// In case a subscription gets deleted we will update it with an owner address equal to 0x0.
func (o *orm) UpsertSubscription(ctx context.Context, subscription StoredSubscription) error {
	stmt := fmt.Sprintf(`
		INSERT INTO %s (subscription_id, owner, balance, blocked_balance, proposed_owner, consumers, flags, router_contract_address)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT (subscription_id, router_contract_address) DO UPDATE
		SET owner=$2, balance=$3, blocked_balance=$4, proposed_owner=$5, consumers=$6, flags=$7, router_contract_address=$8;`, tableName)

	if subscription.Balance == nil {
		subscription.Balance = big.NewInt(0)
	}

	if subscription.BlockedBalance == nil {
		subscription.BlockedBalance = big.NewInt(0)
	}

	var consumers [][]byte
	for _, c := range subscription.Consumers {
		consumers = append(consumers, c.Bytes())
	}

	_, err := o.ds.ExecContext(
		ctx,
		stmt,
		subscription.SubscriptionID,
		subscription.Owner,
		subscription.Balance.Int64(),
		subscription.BlockedBalance.Int64(),
		subscription.ProposedOwner,
		consumers,
		subscription.Flags[:],
		o.routerContractAddress,
	)
	if err != nil {
		return err
	}

	o.lggr.Debugf("Successfully updated subscription: %d for routerContractAddress: %s", subscription.SubscriptionID, o.routerContractAddress)

	return nil
}

func (cs *storedSubscriptionRow) encode() StoredSubscription {
	consumers := make([]common.Address, 0)
	for _, csc := range cs.Consumers {
		consumers = append(consumers, common.BytesToAddress(csc))
	}

	return StoredSubscription{
		SubscriptionID: cs.SubscriptionID,
		IFunctionsSubscriptionsSubscription: functions_router.IFunctionsSubscriptionsSubscription{
			Balance:        big.NewInt(cs.Balance),
			Owner:          cs.Owner,
			BlockedBalance: big.NewInt(cs.BlockedBalance),
			ProposedOwner:  cs.ProposedOwner,
			Consumers:      consumers,
			Flags:          [32]byte(cs.Flags),
		},
	}
}
