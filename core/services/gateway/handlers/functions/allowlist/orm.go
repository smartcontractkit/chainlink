package allowlist

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

//go:generate mockery --quiet --name ORM --output ./mocks/ --case=underscore
type ORM interface {
	GetAllowedSenders(offset, limit uint, qopts ...pg.QOpt) ([]common.Address, error)
	CreateAllowedSenders(allowedSenders []common.Address, qopts ...pg.QOpt) error
	DeleteAllowedSenders(blockedSenders []common.Address, qopts ...pg.QOpt) error
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
	tableName = "functions_allowlist"
)

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

func (o *orm) GetAllowedSenders(offset, limit uint, qopts ...pg.QOpt) ([]common.Address, error) {
	var addresses []common.Address
	stmt := fmt.Sprintf(`
		SELECT allowed_address
		FROM %s
		WHERE router_contract_address = $1
		ORDER BY id ASC
		OFFSET $2
		LIMIT $3;
	`, tableName)
	err := o.q.WithOpts(qopts...).Select(&addresses, stmt, o.routerContractAddress, offset, limit)
	if err != nil {
		return addresses, err
	}
	o.lggr.Debugf("Successfully fetched allowed sender list from DB. offset: %d, limit: %d, length: %d", offset, limit, len(addresses))

	return addresses, nil
}

func (o *orm) CreateAllowedSenders(allowedSenders []common.Address, qopts ...pg.QOpt) error {
	var valuesPlaceholder []string
	for i := 1; i <= len(allowedSenders)*2; i += 2 {
		valuesPlaceholder = append(valuesPlaceholder, fmt.Sprintf("($%d, $%d)", i, i+1))
	}

	stmt := fmt.Sprintf(`
		INSERT INTO %s (allowed_address, router_contract_address)
		VALUES %s ON CONFLICT (allowed_address, router_contract_address) DO NOTHING;`, tableName, strings.Join(valuesPlaceholder, ", "))

	var args []interface{}
	for _, as := range allowedSenders {
		args = append(args, as, o.routerContractAddress)
	}

	_, err := o.q.WithOpts(qopts...).Exec(stmt, args...)
	if err != nil {
		return err
	}

	o.lggr.Debugf("Successfully stored allowed senders: %v for routerContractAddress: %s", allowedSenders, o.routerContractAddress)

	return nil
}

func (o *orm) DeleteAllowedSenders(blockedSenders []common.Address, qopts ...pg.QOpt) error {
	var valuesPlaceholder []string
	for i := 1; i <= len(blockedSenders); i++ {
		valuesPlaceholder = append(valuesPlaceholder, fmt.Sprintf("$%d", i+1))
	}

	stmt := fmt.Sprintf(`
		DELETE FROM %s
		WHERE router_contract_address = $1
		AND allowed_address IN (%s);`, tableName, strings.Join(valuesPlaceholder, ", "))

	args := []interface{}{o.routerContractAddress}
	for _, bs := range blockedSenders {
		args = append(args, bs)
	}

	res, err := o.q.WithOpts(qopts...).Exec(stmt, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	o.lggr.Debugf("Successfully removed blocked senders from the allowed list: %v for routerContractAddress: %s. rowsAffected: %d", blockedSenders, o.routerContractAddress, rowsAffected)

	return nil
}
