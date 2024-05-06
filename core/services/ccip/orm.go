package ccip

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
)

type GasPrice struct {
	SourceChainSelector uint64
	GasPrice            *assets.Wei
	CreatedAt           time.Time
}

type GasPriceUpdate struct {
	SourceChainSelector uint64
	GasPrice            *assets.Wei
}

type TokenPrice struct {
	TokenAddr  string
	TokenPrice *assets.Wei
	CreatedAt  time.Time
}

type TokenPriceUpdate struct {
	TokenAddr  string
	TokenPrice *assets.Wei
}

//go:generate mockery --quiet --name ORM --output ./mocks/ --case=underscore
type ORM interface {
	GetGasPricesByDestChain(ctx context.Context, destChainSelector uint64) ([]GasPrice, error)
	GetTokenPricesByDestChain(ctx context.Context, destChainSelector uint64) ([]TokenPrice, error)

	InsertGasPricesForDestChain(ctx context.Context, destChainSelector uint64, jobId int32, gasPrices []GasPriceUpdate) error
	InsertTokenPricesForDestChain(ctx context.Context, destChainSelector uint64, jobId int32, tokenPrices []TokenPriceUpdate) error

	ClearGasPricesByDestChain(ctx context.Context, destChainSelector uint64, to time.Time) error
	ClearTokenPricesByDestChain(ctx context.Context, destChainSelector uint64, to time.Time) error
}

type orm struct {
	ds sqlutil.DataSource
}

var _ ORM = (*orm)(nil)

func NewORM(ds sqlutil.DataSource) (ORM, error) {
	if ds == nil {
		return nil, fmt.Errorf("datasource to CCIP NewORM cannot be nil")
	}

	return &orm{
		ds: ds,
	}, nil
}

func (o *orm) GetGasPricesByDestChain(ctx context.Context, destChainSelector uint64) ([]GasPrice, error) {
	var gasPrices []GasPrice
	stmt := `
		SELECT DISTINCT ON (source_chain_selector)
		source_chain_selector, gas_price, created_at
		FROM ccip.observed_gas_prices
		WHERE chain_selector = $1
		ORDER BY source_chain_selector, created_at DESC;
	`
	err := o.ds.SelectContext(ctx, &gasPrices, stmt, destChainSelector)
	if err != nil {
		return nil, err
	}

	return gasPrices, nil
}

func (o *orm) GetTokenPricesByDestChain(ctx context.Context, destChainSelector uint64) ([]TokenPrice, error) {
	var tokenPrices []TokenPrice
	stmt := `
		SELECT DISTINCT ON (token_addr)
		token_addr, token_price, created_at
		FROM ccip.observed_token_prices
		WHERE chain_selector = $1
		ORDER BY token_addr, created_at DESC;
	`
	err := o.ds.SelectContext(ctx, &tokenPrices, stmt, destChainSelector)
	if err != nil {
		return nil, err
	}

	return tokenPrices, nil
}

func (o *orm) InsertGasPricesForDestChain(ctx context.Context, destChainSelector uint64, jobId int32, gasPrices []GasPriceUpdate) error {
	if len(gasPrices) == 0 {
		return nil
	}

	now := time.Now()
	sqlStr := ""
	var values []interface{}
	for i, price := range gasPrices {
		sqlStr += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5)
		values = append(values, destChainSelector, jobId, price.SourceChainSelector, price.GasPrice, now)
	}
	// Trim the last comma
	sqlStr = sqlStr[0 : len(sqlStr)-1]

	stmt := fmt.Sprintf(`
		INSERT INTO ccip.observed_gas_prices (chain_selector, job_id, source_chain_selector, gas_price, created_at)
		VALUES %s;`,
		sqlStr)

	_, err := o.ds.ExecContext(ctx, stmt, values...)
	if err != nil {
		err = fmt.Errorf("error inserting gas prices for job %d: %w", jobId, err)
	}
	return err
}

func (o *orm) InsertTokenPricesForDestChain(ctx context.Context, destChainSelector uint64, jobId int32, tokenPrices []TokenPriceUpdate) error {
	if len(tokenPrices) == 0 {
		return nil
	}

	now := time.Now()
	sqlStr := ""
	var values []interface{}
	for i, price := range tokenPrices {
		sqlStr += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5)
		values = append(values, destChainSelector, jobId, price.TokenAddr, price.TokenPrice, now)
	}
	// Trim the last comma
	sqlStr = sqlStr[0 : len(sqlStr)-1]

	stmt := fmt.Sprintf(`
		INSERT INTO ccip.observed_token_prices (chain_selector, job_id, token_addr, token_price, created_at)
		VALUES %s;`,
		sqlStr)

	_, err := o.ds.ExecContext(ctx, stmt, values...)
	if err != nil {
		err = fmt.Errorf("error inserting gas prices for job %d: %w", jobId, err)
	}
	return err
}

func (o *orm) ClearGasPricesByDestChain(ctx context.Context, destChainSelector uint64, to time.Time) error {
	stmt := `DELETE FROM ccip.observed_gas_prices WHERE chain_selector = $1 AND created_at < $2`

	_, err := o.ds.ExecContext(ctx, stmt, destChainSelector, to)
	return err
}

func (o *orm) ClearTokenPricesByDestChain(ctx context.Context, destChainSelector uint64, to time.Time) error {
	stmt := `DELETE FROM ccip.observed_token_prices WHERE chain_selector = $1 AND created_at < $2`

	_, err := o.ds.ExecContext(ctx, stmt, destChainSelector, to)
	return err
}
