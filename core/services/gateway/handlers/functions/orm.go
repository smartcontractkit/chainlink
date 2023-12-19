package functions

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

//go:generate mockery --quiet --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	FetchSubscriptions(offset, limit uint, qopts ...pg.QOpt) ([]CachedSubscription, error)
	CreateSubscription(subscription CachedSubscription, qopts ...pg.QOpt) error
}

type orm struct {
	q pg.Q
}

var _ ORM = (*orm)(nil)
var (
	ErrDuplicateSubscriptionID = errors.New("Functions ORM: duplicate subscription ID")
	ErrInvalidParameters       = errors.New("invalid parameters provided to create a subscription cache ORM")
)

const (
	tableName = "functions_subscriptions"
)

type cachedSubscriptionRow struct {
	SubscriptionID uint64
	Balance        int64
	Owner          common.Address
	BlockedBalance int64
	ProposedOwner  common.Address
	Consumers      pq.ByteaArray
	Flags          []uint8
}

func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig) (ORM, error) {
	if db == nil || cfg == nil || lggr == nil {
		return nil, ErrInvalidParameters
	}

	return &orm{
		q: pg.NewQ(db, lggr, cfg),
	}, nil
}

func (o *orm) FetchSubscriptions(offset, limit uint, qopts ...pg.QOpt) ([]CachedSubscription, error) {
	var cacheSubscriptions []CachedSubscription
	var cacheSubscriptionRows []cachedSubscriptionRow
	stmt := fmt.Sprintf(`
		SELECT subscription_id, owner, balance, blocked_balance, proposed_owner, consumers, flags
		FROM %s
		ORDER BY subscription_id DESC
		OFFSET $1
		LIMIT $2;
	`, tableName)
	err := o.q.WithOpts(qopts...).Select(&cacheSubscriptionRows, stmt, offset, limit)
	if err != nil {
		return cacheSubscriptions, err
	}

	for _, cs := range cacheSubscriptionRows {
		consumers := make([]common.Address, 0)
		for _, csc := range cs.Consumers {
			consumers = append(consumers, common.BytesToAddress(csc))
		}
		cacheSubscriptions = append(cacheSubscriptions, CachedSubscription{
			SubscriptionID: cs.SubscriptionID,
			IFunctionsSubscriptionsSubscription: functions_router.IFunctionsSubscriptionsSubscription{
				Balance:        big.NewInt(cs.Balance),
				Owner:          cs.Owner,
				BlockedBalance: big.NewInt(cs.BlockedBalance),
				ProposedOwner:  cs.ProposedOwner,
				Consumers:      consumers,
				Flags:          [32]byte(cs.Flags),
			},
		})
	}

	return cacheSubscriptions, nil
}

func (o *orm) CreateSubscription(subscription CachedSubscription, qopts ...pg.QOpt) error {
	stmt := fmt.Sprintf(`
		INSERT INTO %s (subscription_id, owner, balance, blocked_balance, proposed_owner, consumers, flags)
		VALUES ($1,$2,$3,$4,$5,$6,$7) ON CONFLICT (subscription_id) DO NOTHING;
	`, tableName)

	if subscription.Balance == nil {
		subscription.Balance = big.NewInt(0)
	}

	if subscription.BlockedBalance == nil {
		subscription.BlockedBalance = big.NewInt(0)
	}

	result, err := o.q.WithOpts(qopts...).Exec(
		stmt,
		subscription.SubscriptionID,
		subscription.Owner,
		subscription.Balance.Int64(),
		subscription.BlockedBalance.Int64(),
		subscription.ProposedOwner,
		subscription.Consumers,
		subscription.Flags[:],
	)
	if err != nil {
		return err
	}
	nrows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if nrows == 0 {
		return ErrDuplicateSubscriptionID
	}
	return nil
}

type noopORM struct{}

func NewNoopORM() ORM {
	return &noopORM{}
}
func (o *noopORM) FetchSubscriptions(offset, limit uint, qopts ...pg.QOpt) ([]CachedSubscription, error) {
	return nil, nil
}
func (o *noopORM) CreateSubscription(subscription CachedSubscription, qopts ...pg.QOpt) error {
	return nil
}
