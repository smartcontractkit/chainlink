package ccip

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

type GasPrice struct {
	SourceChainSelector uint64
	GasPrice            *big.Int
	CreatedAt           time.Time
}

type GasPriceRow struct {
	SourceChainSelector uint64
	GasPrice            string
	CreatedAt           time.Time
}

type GasPriceUpdate struct {
	SourceChainSelector uint64
	GasPrice            *big.Int
}

type TokenPrice struct {
	TokenAddr  ccip.Address
	TokenPrice *big.Int
	CreatedAt  time.Time
}

type TokenPriceRow struct {
	TokenAddr  string
	TokenPrice string
	CreatedAt  time.Time
}

type TokenPriceUpdate struct {
	TokenAddr  ccip.Address
	TokenPrice *big.Int
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

const (
	gasTableName   = "ccip.observed_gas_prices"
	tokenTableName = "ccip.observed_token_prices"
)

func NewORM(ds sqlutil.DataSource) (ORM, error) {
	if ds == nil {
		return nil, fmt.Errorf("datasource to CCIP NewORM cannot be nil")
	}

	return &orm{
		ds: ds,
	}, nil
}

func (o *orm) GetGasPricesByDestChain(ctx context.Context, destChainSelector uint64) ([]GasPrice, error) {
	var gasPriceRows []GasPriceRow
	stmt := fmt.Sprintf(`
		SELECT DISTINCT ON (source_chain_selector)
		source_chain_selector, gas_price, created_at
		FROM %s
		WHERE chain_selector = $1
		ORDER BY source_chain_selector, created_at DESC;
	`, gasTableName)
	err := o.ds.SelectContext(ctx, &gasPriceRows, stmt, destChainSelector)
	if err != nil {
		return nil, err
	}

	gasPrices := make([]GasPrice, len(gasPriceRows))
	for i, row := range gasPriceRows {
		price, ok := new(big.Int).SetString(row.GasPrice, 10)
		if !ok {
			return nil, fmt.Errorf("error parsing gas price fetched from db: %s", row.GasPrice)
		}
		gasPrices[i] = GasPrice{
			SourceChainSelector: row.SourceChainSelector,
			GasPrice:            price,
			CreatedAt:           row.CreatedAt,
		}
	}

	return gasPrices, nil
}

func (o *orm) GetTokenPricesByDestChain(ctx context.Context, destChainSelector uint64) ([]TokenPrice, error) {
	var tokenPriceRows []TokenPriceRow
	stmt := fmt.Sprintf(`
		SELECT DISTINCT ON (token_addr)
		token_addr, token_price, created_at
		FROM %s
		WHERE chain_selector = $1
		ORDER BY token_addr, created_at DESC;

	`, tokenTableName)
	err := o.ds.SelectContext(ctx, &tokenPriceRows, stmt, destChainSelector)
	if err != nil {
		return nil, err
	}

	tokenPrices := make([]TokenPrice, len(tokenPriceRows))
	for i, row := range tokenPriceRows {
		price, ok := new(big.Int).SetString(row.TokenPrice, 10)
		if !ok {
			return nil, fmt.Errorf("error parsing token price fetched from db: %s", row.TokenPrice)
		}
		tokenPrices[i] = TokenPrice{
			TokenAddr:  ccip.Address(row.TokenAddr),
			TokenPrice: price,
			CreatedAt:  row.CreatedAt,
		}
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
		values = append(values, destChainSelector, jobId, price.SourceChainSelector, price.GasPrice.String(), now)
	}
	// Trim the last comma
	sqlStr = sqlStr[0 : len(sqlStr)-1]

	stmt := fmt.Sprintf(`
		INSERT INTO %s (chain_selector, job_id, source_chain_selector, gas_price, created_at)
		VALUES %s;`,
		gasTableName, sqlStr)

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
		values = append(values, destChainSelector, jobId, string(price.TokenAddr), price.TokenPrice.String(), now)
	}
	// Trim the last comma
	sqlStr = sqlStr[0 : len(sqlStr)-1]

	stmt := fmt.Sprintf(`
		INSERT INTO %s (chain_selector, job_id, token_addr, token_price, created_at)
		VALUES %s;`,
		tokenTableName, sqlStr)

	_, err := o.ds.ExecContext(ctx, stmt, values...)
	if err != nil {
		err = fmt.Errorf("error inserting gas prices for job %d: %w", jobId, err)
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
