package ccip

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type GasPrice struct {
	SourceChainSelector uint64
	GasPrice            *assets.Wei
}

type TokenPrice struct {
	TokenAddr  string
	TokenPrice *assets.Wei
}

type ORM interface {
	GetGasPricesByDestChain(ctx context.Context, destChainSelector uint64) ([]GasPrice, error)
	GetTokenPricesByDestChain(ctx context.Context, destChainSelector uint64) ([]TokenPrice, error)

	UpsertGasPricesForDestChain(ctx context.Context, destChainSelector uint64, gasPrices []GasPrice) (int64, error)
	UpsertTokenPricesForDestChain(ctx context.Context, destChainSelector uint64, tokenPrices []TokenPrice, interval time.Duration) (int64, error)
}

type orm struct {
	ds   sqlutil.DataSource
	lggr logger.Logger
}

var _ ORM = (*orm)(nil)

func NewORM(ds sqlutil.DataSource, lggr logger.Logger) (ORM, error) {
	if ds == nil {
		return nil, fmt.Errorf("datasource to CCIP NewORM cannot be nil")
	}

	return &orm{
		ds:   ds,
		lggr: lggr,
	}, nil
}

func (o *orm) GetGasPricesByDestChain(ctx context.Context, destChainSelector uint64) ([]GasPrice, error) {
	var gasPrices []GasPrice
	stmt := `
		SELECT source_chain_selector, gas_price
		FROM ccip.observed_gas_prices
		WHERE chain_selector = $1;
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
		SELECT token_addr, token_price
		FROM ccip.observed_token_prices
		WHERE chain_selector = $1;
	`
	err := o.ds.SelectContext(ctx, &tokenPrices, stmt, destChainSelector)
	if err != nil {
		return nil, err
	}
	return tokenPrices, nil
}

func (o *orm) UpsertGasPricesForDestChain(ctx context.Context, destChainSelector uint64, gasPrices []GasPrice) (int64, error) {
	if len(gasPrices) == 0 {
		return 0, nil
	}

	uniqueGasUpdates := make(map[string]GasPrice)
	for _, gasPrice := range gasPrices {
		key := fmt.Sprintf("%d-%d", gasPrice.SourceChainSelector, destChainSelector)
		uniqueGasUpdates[key] = gasPrice
	}

	insertData := make([]map[string]interface{}, 0, len(uniqueGasUpdates))
	for _, price := range uniqueGasUpdates {
		insertData = append(insertData, map[string]interface{}{
			"chain_selector":        destChainSelector,
			"source_chain_selector": price.SourceChainSelector,
			"gas_price":             price.GasPrice,
		})
	}

	stmt := `INSERT INTO ccip.observed_gas_prices (chain_selector, source_chain_selector, gas_price, updated_at)
		VALUES (:chain_selector, :source_chain_selector, :gas_price, statement_timestamp())
		ON CONFLICT (source_chain_selector, chain_selector)
		DO UPDATE SET gas_price = EXCLUDED.gas_price, updated_at = EXCLUDED.updated_at;`

	result, err := o.ds.NamedExecContext(ctx, stmt, insertData)
	if err != nil {
		return 0, fmt.Errorf("error inserting gas prices %w", err)
	}
	return result.RowsAffected()
}

// UpsertTokenPricesForDestChain inserts or updates only relevant token prices.
// In order to reduce locking an unnecessary writes to the table, we start with fetching current prices.
// If price for a token doesn't change or was updated recently we don't include that token to the upsert query.
// We don't run in TX intentionally, because we don't want to lock the table and conflicts are resolved on the insert level
func (o *orm) UpsertTokenPricesForDestChain(ctx context.Context, destChainSelector uint64, tokenPrices []TokenPrice, interval time.Duration) (int64, error) {
	if len(tokenPrices) == 0 {
		return 0, nil
	}

	tokensToUpdate, err := o.pickOnlyRelevantTokensForUpdate(ctx, destChainSelector, tokenPrices, interval)
	if err != nil || len(tokensToUpdate) == 0 {
		return 0, err
	}

	insertData := make([]map[string]interface{}, 0, len(tokensToUpdate))
	for _, price := range tokensToUpdate {
		insertData = append(insertData, map[string]interface{}{
			"chain_selector": destChainSelector,
			"token_addr":     price.TokenAddr,
			"token_price":    price.TokenPrice,
		})
	}

	stmt := `INSERT INTO ccip.observed_token_prices (chain_selector, token_addr, token_price, updated_at)
		VALUES (:chain_selector, :token_addr, :token_price, statement_timestamp())
		ON CONFLICT (token_addr, chain_selector) 
		DO UPDATE SET token_price = EXCLUDED.token_price, updated_at = EXCLUDED.updated_at;`
	result, err := o.ds.NamedExecContext(ctx, stmt, insertData)
	if err != nil {
		return 0, fmt.Errorf("error inserting token prices %w", err)
	}
	return result.RowsAffected()
}

// pickOnlyRelevantTokensForUpdate returns only tokens that need to be updated. Multiple jobs can be updating the same tokens,
// in order to reduce table locking and redundant upserts we start with reading the table and checking which tokens are eligible for update.
// A token is eligible for update when time since last update is greater than the interval.
func (o *orm) pickOnlyRelevantTokensForUpdate(
	ctx context.Context,
	destChainSelector uint64,
	tokenPrices []TokenPrice,
	interval time.Duration,
) ([]TokenPrice, error) {
	tokenPricesByAddress := toTokensByAddress(tokenPrices)

	// Picks only tokens which were recently updated and can be ignored,
	// we will filter out these tokens from the upsert query.
	stmt := `
		SELECT 
		    token_addr
		FROM ccip.observed_token_prices
		WHERE 
		    chain_selector = $1
			and token_addr = any($2)
			and updated_at >= statement_timestamp() - $3::interval
	`

	pgInterval := fmt.Sprintf("%d milliseconds", interval.Milliseconds())
	args := []interface{}{destChainSelector, tokenAddrsToBytes(tokenPricesByAddress), pgInterval}
	var dbTokensToIgnore []string
	if err := o.ds.SelectContext(ctx, &dbTokensToIgnore, stmt, args...); err != nil {
		return nil, err
	}

	tokensToIgnore := make(map[string]struct{}, len(dbTokensToIgnore))
	for _, tk := range dbTokensToIgnore {
		tokensToIgnore[tk] = struct{}{}
	}

	tokenPricesToUpdate := make([]TokenPrice, 0, len(tokenPrices))
	for tokenAddr, tokenPrice := range tokenPricesByAddress {
		eligibleForUpdate := false
		if _, ok := tokensToIgnore[tokenAddr]; !ok {
			eligibleForUpdate = true
			tokenPricesToUpdate = append(tokenPricesToUpdate, TokenPrice{TokenAddr: tokenAddr, TokenPrice: tokenPrice})
		}
		o.lggr.Debugw(
			"Token price eligibility for database update",
			"eligibleForUpdate", eligibleForUpdate,
			"token", tokenAddr,
			"price", tokenPrice,
		)
	}
	return tokenPricesToUpdate, nil
}

func toTokensByAddress(tokens []TokenPrice) map[string]*assets.Wei {
	tokensByAddr := make(map[string]*assets.Wei, len(tokens))
	for _, tk := range tokens {
		tokensByAddr[tk.TokenAddr] = tk.TokenPrice
	}
	return tokensByAddr
}

func tokenAddrsToBytes(tokens map[string]*assets.Wei) [][]byte {
	addrs := make([][]byte, 0, len(tokens))
	for tkAddr := range tokens {
		addrs = append(addrs, []byte(tkAddr))
	}
	return addrs
}
