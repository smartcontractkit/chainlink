package ccip

import (
	"fmt"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/jmoiron/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
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
	GetGasPricesByDestChain(destChainSelector uint64, qopts ...pg.QOpt) ([]GasPrice, error)
	GetTokenPricesByDestChain(destChainSelector uint64, qopts ...pg.QOpt) ([]TokenPrice, error)

	InsertGasPricesForDestChain(destChainSelector uint64, jobId int32, gasPrices []GasPrice, qopts ...pg.QOpt) error
	InsertTokenPricesForDestChain(destChainSelector uint64, jobId int32, tokenPrices []TokenPrice, qopts ...pg.QOpt) error

	ClearGasPricesByDestChain(destChainSelector uint64, to time.Time, qopts ...pg.QOpt) error
	ClearTokenPricesByDestChain(destChainSelector uint64, to time.Time, qopts ...pg.QOpt) error
}

type orm struct {
	q    pg.Q
	lggr logger.Logger
}

var _ ORM = (*orm)(nil)

const (
	gasTableName   = "ccip.observed_gas_prices"
	tokenTableName = "ccip.observed_token_prices"
)

func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig) (ORM, error) {
	if db == nil || lggr == nil || cfg == nil {
		return nil, fmt.Errorf("params to CCIP NewORM cannot be nil")
	}

	namedLogger := lggr.Named("CCIP_ORM")

	return &orm{
		q:    pg.NewQ(db, namedLogger, cfg),
		lggr: namedLogger,
	}, nil
}

func (o *orm) GetGasPricesByDestChain(destChainSelector uint64, qopts ...pg.QOpt) ([]GasPrice, error) {
	var gasPrices []GasPrice
	stmt := fmt.Sprintf(`
		SELECT DISTINCT ON (source_chain_selector)
		source_chain_selector, gas_price, created_at
		FROM %s
		WHERE chain_selector = $1
		ORDER BY source_chain_selector, created_at DESC;
	`, gasTableName)
	err := o.q.WithOpts(qopts...).Select(&gasPrices, stmt, destChainSelector)
	if err != nil {
		return nil, err
	}

	return gasPrices, nil
}

func (o *orm) GetTokenPricesByDestChain(destChainSelector uint64, qopts ...pg.QOpt) ([]TokenPrice, error) {
	var tokenPrices []TokenPrice
	stmt := fmt.Sprintf(`
		SELECT DISTINCT ON (token_addr)
		token_addr, token_price, created_at
		FROM %s
		WHERE chain_selector = $1
		ORDER BY token_addr, created_at DESC;

	`, tokenTableName)
	err := o.q.WithOpts(qopts...).Select(&tokenPrices, stmt, destChainSelector)
	if err != nil {
		return nil, err
	}

	return tokenPrices, nil
}

func (o *orm) InsertGasPricesForDestChain(destChainSelector uint64, jobId int32, gasPrices []GasPrice, qopts ...pg.QOpt) error {
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

	_, err := o.q.WithOpts(qopts...).Exec(
		stmt,
		values...,
	)
	if err != nil {
		o.lggr.Errorf("Error inserting gas prices for job %d: %v", jobId, err)
	}
	return err
}

func (o *orm) InsertTokenPricesForDestChain(destChainSelector uint64, jobId int32, tokenPrices []TokenPrice, qopts ...pg.QOpt) error {
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

	_, err := o.q.WithOpts(qopts...).Exec(
		stmt,
		values...,
	)
	if err != nil {
		o.lggr.Errorf("Error inserting token prices for job %d: %v", jobId, err)
	}
	return err
}

func (o *orm) ClearGasPricesByDestChain(destChainSelector uint64, to time.Time, qopts ...pg.QOpt) error {
	stmt := fmt.Sprintf(`DELETE FROM %s WHERE chain_selector = $1 AND created_at < $2`, gasTableName)

	_, err := o.q.WithOpts(qopts...).Exec(
		stmt,
		destChainSelector,
		to,
	)
	return err
}

func (o *orm) ClearTokenPricesByDestChain(destChainSelector uint64, to time.Time, qopts ...pg.QOpt) error {
	stmt := fmt.Sprintf(`DELETE FROM %s WHERE chain_selector = $1 AND created_at < $2`, tokenTableName)

	_, err := o.q.WithOpts(qopts...).Exec(
		stmt,
		destChainSelector,
		to,
	)
	return err
}
