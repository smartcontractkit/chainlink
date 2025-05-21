package ccip

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var (
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func setupORM(t *testing.T) (ORM, sqlutil.DataSource) {
	t.Helper()

	db := pgtest.NewSqlxDB(t)
	orm, err := NewORM(db, logger.TestLogger(t))

	require.NoError(t, err)

	return orm, db
}

func generateChainSelectors(n int) []uint64 {
	selectors := make([]uint64, n)
	for i := 0; i < n; i++ {
		selectors[i] = rand.Uint64()
	}

	return selectors
}

func generateGasPrices(chainSelector uint64, n int) []GasPrice {
	updates := make([]GasPrice, n)
	for i := 0; i < n; i++ {
		// gas prices can take up whole range of uint256
		uint256Max := new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), big.NewInt(1))
		row := GasPrice{
			SourceChainSelector: chainSelector,
			GasPrice:            assets.NewWei(new(big.Int).Sub(uint256Max, big.NewInt(int64(i)))),
		}
		updates[i] = row
	}

	return updates
}

func generateTokenAddresses(n int) []string {
	addrs := make([]string, n)
	for i := 0; i < n; i++ {
		addrs[i] = utils.RandomAddress().Hex()
	}

	return addrs
}

func generateRandomTokenPrices(tokenAddrs []string) []TokenPrice {
	updates := make([]TokenPrice, 0, len(tokenAddrs))
	for _, addr := range tokenAddrs {
		updates = append(updates, TokenPrice{
			TokenAddr:  addr,
			TokenPrice: assets.NewWei(new(big.Int).Rand(r, big.NewInt(1e18))),
		})
	}
	return updates
}

func generateTokenPrices(tokenAddr string, n int) []TokenPrice {
	updates := make([]TokenPrice, n)
	for i := 0; i < n; i++ {
		row := TokenPrice{
			TokenAddr:  tokenAddr,
			TokenPrice: assets.NewWei(new(big.Int).Mul(big.NewInt(1e18), big.NewInt(int64(i)))),
		}
		updates[i] = row
	}

	return updates
}

func getGasTableRowCount(t *testing.T, ds sqlutil.DataSource) int {
	var count int
	stmt := `SELECT COUNT(*) FROM ccip.observed_gas_prices;`
	err := ds.QueryRowxContext(testutils.Context(t), stmt).Scan(&count)
	require.NoError(t, err)

	return count
}

func getTokenTableRowCount(t *testing.T, ds sqlutil.DataSource) int {
	var count int
	stmt := `SELECT COUNT(*) FROM ccip.observed_token_prices;`
	err := ds.QueryRowxContext(testutils.Context(t), stmt).Scan(&count)
	require.NoError(t, err)

	return count
}

func TestInitORM(t *testing.T) {
	t.Parallel()

	orm, _ := setupORM(t)
	assert.NotNil(t, orm)
}

func TestORM_EmptyGasPrices(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	orm, _ := setupORM(t)

	prices, err := orm.GetGasPricesByDestChain(ctx, 1)
	assert.Empty(t, prices)
	assert.NoError(t, err)
}

func TestORM_EmptyTokenPrices(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	orm, _ := setupORM(t)

	prices, err := orm.GetTokenPricesByDestChain(ctx, 1)
	assert.Empty(t, prices)
	assert.NoError(t, err)
}

func TestORM_InsertAndGetGasPrices(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	orm, db := setupORM(t)

	numJobs := 5
	numSourceChainSelectors := 10
	numUpdatesPerSourceSelector := 20
	destSelector := uint64(1)

	sourceSelectors := generateChainSelectors(numSourceChainSelectors)

	updates := make(map[uint64][]GasPrice)
	for _, selector := range sourceSelectors {
		updates[selector] = generateGasPrices(selector, numUpdatesPerSourceSelector)
	}

	// 5 jobs, each inserting prices for 10 chains, with 20 updates per chain.
	expectedPrices := make(map[uint64]GasPrice)
	for i := 0; i < numJobs; i++ {
		for selector, updatesPerSelector := range updates {
			lastIndex := len(updatesPerSelector) - 1

			_, err := orm.UpsertGasPricesForDestChain(ctx, destSelector, updatesPerSelector[:lastIndex])
			assert.NoError(t, err)
			_, err = orm.UpsertGasPricesForDestChain(ctx, destSelector, updatesPerSelector[lastIndex:])
			assert.NoError(t, err)

			expectedPrices[selector] = updatesPerSelector[lastIndex]
		}
	}

	// verify number of rows inserted
	numRows := getGasTableRowCount(t, db)
	assert.Equal(t, numSourceChainSelectors, numRows)

	prices, err := orm.GetGasPricesByDestChain(ctx, destSelector)
	assert.NoError(t, err)
	// should return 1 price per source chain selector
	assert.Len(t, prices, numSourceChainSelectors)

	// verify getGasPrices returns prices of latest timestamp
	for _, price := range prices {
		selector := price.SourceChainSelector
		assert.Equal(t, expectedPrices[selector].GasPrice, price.GasPrice)
	}

	// after the initial inserts, insert new round of prices, 1 price per selector this time
	var combinedUpdates []GasPrice
	for selector, updatesPerSelector := range updates {
		combinedUpdates = append(combinedUpdates, updatesPerSelector[0])
		expectedPrices[selector] = updatesPerSelector[0]
	}

	_, err = orm.UpsertGasPricesForDestChain(ctx, destSelector, combinedUpdates)
	assert.NoError(t, err)
	assert.Equal(t, numSourceChainSelectors, getGasTableRowCount(t, db))

	prices, err = orm.GetGasPricesByDestChain(ctx, destSelector)
	assert.NoError(t, err)
	assert.Len(t, prices, numSourceChainSelectors)

	for _, price := range prices {
		selector := price.SourceChainSelector
		assert.Equal(t, expectedPrices[selector].GasPrice, price.GasPrice)
	}
}

func TestORM_UpsertGasPrices(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	orm, db := setupORM(t)

	numSourceChainSelectors := 10
	numUpdatesPerSourceSelector := 20
	destSelector := uint64(1)

	sourceSelectors := generateChainSelectors(numSourceChainSelectors)

	updates := make(map[uint64][]GasPrice)
	for _, selector := range sourceSelectors {
		updates[selector] = generateGasPrices(selector, numUpdatesPerSourceSelector)
	}

	for _, updatesPerSelector := range updates {
		_, err := orm.UpsertGasPricesForDestChain(ctx, destSelector, updatesPerSelector)
		assert.NoError(t, err)
	}

	sleepSec := 2
	time.Sleep(time.Duration(sleepSec) * time.Second)

	// insert for the 2nd time after interimTimeStamp
	for _, updatesPerSelector := range updates {
		_, err := orm.UpsertGasPricesForDestChain(ctx, destSelector, updatesPerSelector)
		assert.NoError(t, err)
	}

	assert.Equal(t, numSourceChainSelectors, getGasTableRowCount(t, db))
}

func TestORM_InsertAndGetTokenPrices(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	orm, db := setupORM(t)

	numJobs := 5
	numAddresses := 10
	numUpdatesPerAddress := 20
	destSelector := uint64(1)

	addrs := generateTokenAddresses(numAddresses)

	updates := make(map[string][]TokenPrice)
	for _, addr := range addrs {
		updates[addr] = generateTokenPrices(addr, numUpdatesPerAddress)
	}

	// 5 jobs, each inserting prices for 10 chains, with 20 updates per chain.
	expectedPrices := make(map[string]TokenPrice)
	for i := 0; i < numJobs; i++ {
		for addr, updatesPerAddr := range updates {
			lastIndex := len(updatesPerAddr) - 1

			_, err := orm.UpsertTokenPricesForDestChain(ctx, destSelector, updatesPerAddr[:lastIndex], 0)
			assert.NoError(t, err)
			_, err = orm.UpsertTokenPricesForDestChain(ctx, destSelector, updatesPerAddr[lastIndex:], 0)
			assert.NoError(t, err)

			expectedPrices[addr] = updatesPerAddr[lastIndex]
		}
	}

	// verify number of rows inserted
	numRows := getTokenTableRowCount(t, db)
	assert.Equal(t, numAddresses, numRows)

	prices, err := orm.GetTokenPricesByDestChain(ctx, destSelector)
	assert.NoError(t, err)
	// should return 1 price per source chain selector
	assert.Len(t, prices, numAddresses)

	// verify getTokenPrices returns prices of latest timestamp
	for _, price := range prices {
		addr := price.TokenAddr
		assert.Equal(t, expectedPrices[addr].TokenPrice, price.TokenPrice)
	}

	// after the initial inserts, insert new round of prices, 1 price per selector this time
	var combinedUpdates []TokenPrice
	for addr, updatesPerAddr := range updates {
		combinedUpdates = append(combinedUpdates, updatesPerAddr[0])
		expectedPrices[addr] = updatesPerAddr[0]
	}

	_, err = orm.UpsertTokenPricesForDestChain(ctx, destSelector, combinedUpdates, 0)
	assert.NoError(t, err)
	assert.Equal(t, numAddresses, getTokenTableRowCount(t, db))

	prices, err = orm.GetTokenPricesByDestChain(ctx, destSelector)
	assert.NoError(t, err)
	assert.Len(t, prices, numAddresses)

	for _, price := range prices {
		addr := price.TokenAddr
		assert.Equal(t, expectedPrices[addr].TokenPrice, price.TokenPrice)
	}
}

func TestORM_InsertTokenPricesWhenExpired(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	orm, _ := setupORM(t)

	numAddresses := 10
	destSelector := rand.Uint64()
	addrs := generateTokenAddresses(numAddresses)
	initTokenUpdates := generateRandomTokenPrices(addrs)

	// Insert the first time, table is initialized
	rowsUpdated, err := orm.UpsertTokenPricesForDestChain(ctx, destSelector, initTokenUpdates, time.Minute)
	require.NoError(t, err)
	assert.Equal(t, int64(numAddresses), rowsUpdated)

	// time.Sleep(100 * time.Millisecond)

	// Insert the second time, no updates, because prices haven't changed
	rowsUpdated, err = orm.UpsertTokenPricesForDestChain(ctx, destSelector, initTokenUpdates, time.Minute)
	require.NoError(t, err)
	assert.Equal(t, int64(0), rowsUpdated)

	// There are new prices, but we still haven't reached interval
	newPrices := generateRandomTokenPrices(addrs)
	rowsUpdated, err = orm.UpsertTokenPricesForDestChain(ctx, destSelector, newPrices, time.Minute)
	require.NoError(t, err)
	assert.Equal(t, int64(0), rowsUpdated)

	time.Sleep(100 * time.Millisecond)

	// Again with the same new prices, but this time interval is reached
	rowsUpdated, err = orm.UpsertTokenPricesForDestChain(ctx, destSelector, newPrices, time.Millisecond)
	require.NoError(t, err)
	assert.Equal(t, int64(numAddresses), rowsUpdated)

	dbTokenPrices, err := orm.GetTokenPricesByDestChain(ctx, destSelector)
	require.NoError(t, err)
	assert.Len(t, dbTokenPrices, numAddresses)

	dbTokenPricesByAddr := toTokensByAddress(dbTokenPrices)
	for _, tkPrice := range newPrices {
		dbToken, ok := dbTokenPricesByAddr[tkPrice.TokenAddr]
		assert.True(t, ok)
		assert.Equal(t, dbToken, tkPrice.TokenPrice)
	}
}

func Benchmark_UpsertsTheSameTokenPrices(b *testing.B) {
	db := pgtest.NewSqlxDB(b)
	orm, err := NewORM(db, logger.NullLogger)
	require.NoError(b, err)

	ctx := testutils.Context(b)
	numAddresses := 50
	destSelector := rand.Uint64()
	addrs := generateTokenAddresses(numAddresses)
	tokenUpdates := generateRandomTokenPrices(addrs)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err1 := orm.UpsertTokenPricesForDestChain(ctx, destSelector, tokenUpdates, time.Second)
		require.NoError(b, err1)
	}
}
