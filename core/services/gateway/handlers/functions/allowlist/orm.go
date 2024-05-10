package allowlist

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

//go:generate mockery --quiet --name ORM --output ./mocks/ --case=underscore
type ORM interface {
	GetAllowedSenders(ctx context.Context, offset, limit uint) ([]common.Address, error)
	CreateAllowedSenders(ctx context.Context, allowedSenders []common.Address) error
	DeleteAllowedSenders(ctx context.Context, blockedSenders []common.Address) error
	PurgeAllowedSenders(ctx context.Context) error
}

type orm struct {
	ds                    sqlutil.DataSource
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

func (o *orm) GetAllowedSenders(ctx context.Context, offset, limit uint) ([]common.Address, error) {
	var addresses []common.Address
	stmt := fmt.Sprintf(`
		SELECT allowed_address
		FROM %s
		WHERE router_contract_address = $1
		ORDER BY id ASC
		OFFSET $2
		LIMIT $3;
	`, tableName)
	err := o.ds.SelectContext(ctx, &addresses, stmt, o.routerContractAddress, offset, limit)
	if err != nil {
		return addresses, err
	}
	o.lggr.Debugf("Successfully fetched allowed sender list from DB. offset: %d, limit: %d, length: %d", offset, limit, len(addresses))

	return addresses, nil
}

func (o *orm) CreateAllowedSenders(ctx context.Context, allowedSenders []common.Address) error {
	if len(allowedSenders) == 0 {
		o.lggr.Debugf("empty allowed senders list: %v for routerContractAddress: %s. skipping...", allowedSenders, o.routerContractAddress)
		return nil
	}

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

	_, err := o.ds.ExecContext(ctx, stmt, args...)
	if err != nil {
		return err
	}

	o.lggr.Debugf("Successfully stored allowed senders: %v for routerContractAddress: %s", allowedSenders, o.routerContractAddress)

	return nil
}

// DeleteAllowedSenders is used to remove blocked senders from the functions_allowlist table.
// This is achieved by specifying a list of blockedSenders to remove.
func (o *orm) DeleteAllowedSenders(ctx context.Context, blockedSenders []common.Address) error {
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

	res, err := o.ds.ExecContext(ctx, stmt, args...)
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

// PurgeAllowedSenders will remove all the allowed senders for the configured orm routerContractAddress
func (o *orm) PurgeAllowedSenders(ctx context.Context) error {
	stmt := fmt.Sprintf(`
		DELETE FROM %s
		WHERE router_contract_address = $1;`, tableName)

	res, err := o.ds.ExecContext(ctx, stmt, o.routerContractAddress)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	o.lggr.Debugf("Successfully purged allowed senders for routerContractAddress: %s. rowsAffected: %d", o.routerContractAddress, rowsAffected)

	return nil
}
