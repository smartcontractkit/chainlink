package functions

import (
	"fmt"
	"math/big"
	"strings"

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
	GetSubscriptions(offset, limit uint, qopts ...pg.QOpt) ([]CachedSubscription, error)
	UpsertSubscription(subscription CachedSubscription, qopts ...pg.QOpt) error

	GetAllowedSenders(offset, limit uint, qopts ...pg.QOpt) ([]common.Address, error)
	CreateAllowedSenders(allowedSenders []common.Address, qopts ...pg.QOpt) error
}

type orm struct {
	q                     pg.Q
	lggr                  logger.Logger
	routerContractAddress common.Address
}

var _ ORM = (*orm)(nil)
var (
	ErrInvalidParameters = errors.New("invalid parameters provided to create a functions contract cache ORM")
)

const (
	subscriptionsTableName = "functions_subscriptions"
	allowlistTableName     = "functions_allowlist"
)

type cachedSubscriptionRow struct {
	SubscriptionID        uint64
	Owner                 common.Address
	Balance               int64
	BlockedBalance        int64
	ProposedOwner         common.Address
	Consumers             pq.ByteaArray
	Flags                 []uint8
	RouterContractAddress common.Address
}

func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig, routerContractAddress common.Address) (ORM, error) {
	if db == nil || cfg == nil || lggr == nil || routerContractAddress == (common.Address{}) {
		return nil, ErrInvalidParameters
	}

	return &orm{
		q:                     pg.NewQ(db, lggr, cfg),
		lggr:                  lggr,
		routerContractAddress: routerContractAddress,
	}, nil
}

func (o *orm) GetSubscriptions(offset, limit uint, qopts ...pg.QOpt) ([]CachedSubscription, error) {
	var cacheSubscriptions []CachedSubscription
	var cacheSubscriptionRows []cachedSubscriptionRow
	stmt := fmt.Sprintf(`
		SELECT subscription_id, owner, balance, blocked_balance, proposed_owner, consumers, flags, router_contract_address
		FROM %s
		WHERE router_contract_address = $1
		ORDER BY subscription_id ASC
		OFFSET $2
		LIMIT $3;
	`, subscriptionsTableName)
	err := o.q.WithOpts(qopts...).Select(&cacheSubscriptionRows, stmt, o.routerContractAddress, offset, limit)
	if err != nil {
		return cacheSubscriptions, err
	}

	for _, cs := range cacheSubscriptionRows {
		cacheSubscriptions = append(cacheSubscriptions, cs.encode())
	}

	return cacheSubscriptions, nil
}

// UpsertSubscription will update if a subscription exists or create if it does not.
// In case a subscription gets deleted we will update it with an owner address equal to 0x0.
func (o *orm) UpsertSubscription(subscription CachedSubscription, qopts ...pg.QOpt) error {
	stmt := fmt.Sprintf(`
		INSERT INTO %s (subscription_id, owner, balance, blocked_balance, proposed_owner, consumers, flags, router_contract_address)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT (subscription_id, router_contract_address) DO UPDATE
		SET owner=$2, balance=$3, blocked_balance=$4, proposed_owner=$5, consumers=$6, flags=$7, router_contract_address=$8;`, subscriptionsTableName)

	if subscription.Balance == nil {
		subscription.Balance = big.NewInt(0)
	}

	if subscription.BlockedBalance == nil {
		subscription.BlockedBalance = big.NewInt(0)
	}

	_, err := o.q.WithOpts(qopts...).Exec(
		stmt,
		subscription.SubscriptionID,
		subscription.Owner,
		subscription.Balance.Int64(),
		subscription.BlockedBalance.Int64(),
		subscription.ProposedOwner,
		subscription.Consumers,
		subscription.Flags[:],
		o.routerContractAddress,
	)

	return err
}

func (cs *cachedSubscriptionRow) encode() CachedSubscription {
	consumers := make([]common.Address, 0)
	for _, csc := range cs.Consumers {
		consumers = append(consumers, common.BytesToAddress(csc))
	}

	return CachedSubscription{
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

func (o *orm) GetAllowedSenders(offset, limit uint, qopts ...pg.QOpt) ([]common.Address, error) {
	var addresses []common.Address
	stmt := fmt.Sprintf(`
		SELECT allowed_address
		FROM %s
		WHERE router_contract_address = $1
		ORDER BY id ASC
		OFFSET $2
		LIMIT $3;
	`, allowlistTableName)
	err := o.q.WithOpts(qopts...).Select(&addresses, stmt, o.routerContractAddress, offset, limit)
	if err != nil {
		return addresses, err
	}
	o.lggr.Debugf("Successfully fetched allowed sender list from DB. offset: %d, limit: %d, length: %d", offset, limit, len(addresses))

	return addresses, nil
}

func (o *orm) CreateAllowedSenders(allowedSender []common.Address, qopts ...pg.QOpt) error {
	var valuesPlaceholder []string
	for i := 1; i <= len(allowedSender)*2; i += 2 {
		valuesPlaceholder = append(valuesPlaceholder, fmt.Sprintf("($%d, $%d)", i, i+1))
	}

	stmt := fmt.Sprintf(`
		INSERT INTO %s (allowed_address, router_contract_address)
		VALUES %s ON CONFLICT (allowed_address, router_contract_address) DO NOTHING;`, allowlistTableName, strings.Join(valuesPlaceholder, ", "))

	var args []interface{}
	for _, as := range allowedSender {
		args = append(args, as, o.routerContractAddress)
	}

	_, err := o.q.WithOpts(qopts...).Exec(stmt, args...)
	if err != nil {
		return err
	}
	o.lggr.Debugf("Successfully stored allowed sender: %s for routerContractAddress: %s", allowedSender, o.routerContractAddress)

	return nil
}
