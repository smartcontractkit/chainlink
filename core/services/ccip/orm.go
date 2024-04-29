package ccip

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type GasPrice struct {
	SourceChainSelector uint64
	GasPrice            *big.Int
	CreatedAt           time.Time
}

type TokenPrice struct {
	TokenAddr  ccip.Address
	TokenPrice *big.Int
	CreatedAt  time.Time
}

//go:generate mockery --quiet --name ORM --output ./mocks/ --case=underscore
type ORM interface {
	GetGasPricesByDestChain(ctx context.Context, destChainSelector uint64) ([]GasPrice, error)
	GetTokenPricesByDestChain(ctx context.Context, destChainSelector uint64) ([]TokenPrice, error)

	InsertGasPricesForDestChain(ctx context.Context, destChainSelector uint64, jobId int32, gasPrices []GasPrice) error
	InsertTokenPricesForDestChain(ctx context.Context, destChainSelector uint64, jobId int32, tokenPrices []TokenPrice) error

	ClearGasPricesByDestChain(ctx context.Context, destChainSelector uint64, to time.Time) error
	ClearTokenPricesByDestChain(ctx context.Context, destChainSelector uint64, to time.Time) error
}

type orm struct {
	ds   sqlutil.DataSource
	lggr logger.Logger
}

var _ ORM = (*orm)(nil)

const (
	gasTableName   = "ccip.observed_gas_prices"
	tokenTableName = "ccip.observed_token_prices"
)

func NewORM(ds sqlutil.DataSource, lggr logger.Logger) (ORM, error) {
	if ds == nil || lggr == nil {
		return nil, fmt.Errorf("params to CCIP NewORM cannot be nil")
	}

	namedLogger := lggr.Named("CCIP_ORM")

	return &orm{
		ds:   ds,
		lggr: namedLogger,
	}, nil
}

func (o *orm) GetGasPricesByDestChain(ctx context.Context, destChainSelector uint64) ([]GasPrice, error) {
	var gasPrices []GasPrice
	stmt := fmt.Sprintf(`
		SELECT DISTINCT ON (source_chain_selector)
		source_chain_selector, gas_price, created_at
		FROM %s
		WHERE chain_selector = $1
		ORDER BY source_chain_selector, created_at DESC;
	`, gasTableName)
	err := o.ds.SelectContext(ctx, &gasPrices, stmt, destChainSelector)
	if err != nil {
		return nil, err
	}

	return gasPrices, nil
}

func (o *orm) GetTokenPricesByDestChain(ctx context.Context, destChainSelector uint64) ([]TokenPrice, error) {
	var tokenPrices []TokenPrice
	stmt := fmt.Sprintf(`
		SELECT DISTINCT ON (token_addr)
		token_addr, token_price, created_at
		FROM %s
		WHERE chain_selector = $1
		ORDER BY token_addr, created_at DESC;

	`, tokenTableName)
	err := o.ds.SelectContext(ctx, &tokenPrices, stmt, destChainSelector)
	if err != nil {
		return nil, err
	}

	return tokenPrices, nil
}

func (o *orm) InsertGasPricesForDestChain(ctx context.Context, destChainSelector uint64, jobId int32, gasPrices []GasPrice) error {
	if len(gasPrices) == 0 {
		return nil
	}

	sqlStr := ""
	var values []interface{}
	for _, price := range gasPrices {
		sqlStr += "($1,$2,$3,$4,NOW()),"
		values = append(values, destChainSelector, jobId, price.SourceChainSelector, price.GasPrice)
	}
	// Trim the last comma
	sqlStr = sqlStr[0 : len(sqlStr)-1]

	stmt := fmt.Sprintf(`
		INSERT INTO %s (chain_selector, job_id, source_chain_selector, gas_price, created_at)
		VALUES %s;`,
		gasTableName, sqlStr)

	_, err := o.ds.ExecContext(ctx, stmt, values...)
	if err != nil {
		o.lggr.Errorf("Error inserting gas prices for job %d: %v", jobId, err)
	}
	return err
}

func (o *orm) InsertTokenPricesForDestChain(ctx context.Context, destChainSelector uint64, jobId int32, tokenPrices []TokenPrice) error {
	if len(tokenPrices) == 0 {
		return nil
	}

	sqlStr := ""
	var values []interface{}
	for _, price := range tokenPrices {
		sqlStr += "($1,$2,$3,$4,NOW()),"
		values = append(values, destChainSelector, jobId, price.TokenAddr, price.TokenPrice)
	}
	// Trim the last comma
	sqlStr = sqlStr[0 : len(sqlStr)-1]

	stmt := fmt.Sprintf(`
		INSERT INTO %s (chain_selector, job_id, token_addr, token_price, created_at)
		VALUES %s;`,
		tokenTableName, sqlStr)

	_, err := o.ds.ExecContext(ctx, stmt, values...)
	if err != nil {
		o.lggr.Errorf("Error inserting token prices for job %d: %v", jobId, err)
	}
	return err
}

func (o *orm) ClearGasPricesByDestChain(ctx context.Context, destChainSelector uint64, to time.Time) error {
	stmt := fmt.Sprintf(`DELETE FROM %s WHERE chain_selector = $1 AND created_at < $2`, gasTableName)

	_, err := o.ds.ExecContext(ctx, stmt, destChainSelector, to)
	return err
}

func (o *orm) ClearTokenPricesByDestChain(ctx context.Context, destChainSelector uint64, to time.Time) error {
	stmt := fmt.Sprintf(`DELETE FROM %s WHERE chain_selector = $1 AND created_at < $2`, tokenTableName)

	_, err := o.ds.ExecContext(ctx, stmt, destChainSelector, to)
	return err
}
