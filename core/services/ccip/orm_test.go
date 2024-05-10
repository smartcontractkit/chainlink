package ccip

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

func setupORM(t *testing.T) (ORM, sqlutil.DataSource) {
	t.Helper()

	db := pgtest.NewSqlxDB(t)
	orm, err := NewORM(db)

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

func generateGasPriceUpdates(chainSelector uint64, n int) []GasPriceUpdate {
	updates := make([]GasPriceUpdate, n)
	for i := 0; i < n; i++ {
		// gas prices can take up whole range of uint256
		uint256Max := new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), big.NewInt(1))
		row := GasPriceUpdate{
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

func generateTokenPriceUpdates(tokenAddr string, n int) []TokenPriceUpdate {
	updates := make([]TokenPriceUpdate, n)
	for i := 0; i < n; i++ {
		row := TokenPriceUpdate{
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
	assert.Equal(t, 0, len(prices))
	assert.NoError(t, err)
}

func TestORM_EmptyTokenPrices(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	orm, _ := setupORM(t)

	prices, err := orm.GetTokenPricesByDestChain(ctx, 1)
	assert.Equal(t, 0, len(prices))
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

	updates := make(map[uint64][]GasPriceUpdate)
	for _, selector := range sourceSelectors {
		updates[selector] = generateGasPriceUpdates(selector, numUpdatesPerSourceSelector)
	}

	// 5 jobs, each inserting prices for 10 chains, with 20 updates per chain.
	expectedPrices := make(map[uint64]GasPriceUpdate)
	for i := 0; i < numJobs; i++ {
		for selector, updatesPerSelector := range updates {
			lastIndex := len(updatesPerSelector) - 1

			err := orm.InsertGasPricesForDestChain(ctx, destSelector, int32(i), updatesPerSelector[:lastIndex])
			assert.NoError(t, err)
			err = orm.InsertGasPricesForDestChain(ctx, destSelector, int32(i), updatesPerSelector[lastIndex:])
			assert.NoError(t, err)

			expectedPrices[selector] = updatesPerSelector[lastIndex]
		}
	}

	// verify number of rows inserted
	numRows := getGasTableRowCount(t, db)
	assert.Equal(t, numJobs*numSourceChainSelectors*numUpdatesPerSourceSelector, numRows)

	prices, err := orm.GetGasPricesByDestChain(ctx, destSelector)
	assert.NoError(t, err)
	// should return 1 price per source chain selector
	assert.Equal(t, numSourceChainSelectors, len(prices))

	// verify getGasPrices returns prices of latest timestamp
	for _, price := range prices {
		selector := price.SourceChainSelector
		assert.Equal(t, expectedPrices[selector].GasPrice, price.GasPrice)
	}

	// after the initial inserts, insert new round of prices, 1 price per selector this time
	var combinedUpdates []GasPriceUpdate
	for selector, updatesPerSelector := range updates {
		combinedUpdates = append(combinedUpdates, updatesPerSelector[0])
		expectedPrices[selector] = updatesPerSelector[0]
	}

	err = orm.InsertGasPricesForDestChain(ctx, destSelector, 1, combinedUpdates)
	assert.NoError(t, err)
	assert.Equal(t, numJobs*numSourceChainSelectors*numUpdatesPerSourceSelector+numSourceChainSelectors, getGasTableRowCount(t, db))

	prices, err = orm.GetGasPricesByDestChain(ctx, destSelector)
	assert.NoError(t, err)
	assert.Equal(t, numSourceChainSelectors, len(prices))

	for _, price := range prices {
		selector := price.SourceChainSelector
		assert.Equal(t, expectedPrices[selector].GasPrice, price.GasPrice)
	}
}

func TestORM_InsertAndDeleteGasPrices(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	orm, db := setupORM(t)

	numSourceChainSelectors := 10
	numUpdatesPerSourceSelector := 20
	destSelector := uint64(1)

	sourceSelectors := generateChainSelectors(numSourceChainSelectors)

	updates := make(map[uint64][]GasPriceUpdate)
	for _, selector := range sourceSelectors {
		updates[selector] = generateGasPriceUpdates(selector, numUpdatesPerSourceSelector)
	}

	for _, updatesPerSelector := range updates {
		err := orm.InsertGasPricesForDestChain(ctx, destSelector, 1, updatesPerSelector)
		assert.NoError(t, err)
	}

	sleepSec := 2
	time.Sleep(time.Duration(sleepSec) * time.Second)

	// insert for the 2nd time after interimTimeStamp
	for _, updatesPerSelector := range updates {
		err := orm.InsertGasPricesForDestChain(ctx, destSelector, 1, updatesPerSelector)
		assert.NoError(t, err)
	}

	assert.Equal(t, 2*numSourceChainSelectors*numUpdatesPerSourceSelector, getGasTableRowCount(t, db))

	// clear by sleepSec should delete rows inserted before it
	err := orm.ClearGasPricesByDestChain(ctx, destSelector, sleepSec)
	assert.NoError(t, err)
	assert.Equal(t, numSourceChainSelectors*numUpdatesPerSourceSelector, getGasTableRowCount(t, db))

	// clear by 0 expiration seconds should delete all rows
	err = orm.ClearGasPricesByDestChain(ctx, destSelector, 0)
	assert.NoError(t, err)
	assert.Equal(t, 0, getGasTableRowCount(t, db))
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

	updates := make(map[string][]TokenPriceUpdate)
	for _, addr := range addrs {
		updates[addr] = generateTokenPriceUpdates(addr, numUpdatesPerAddress)
	}

	// 5 jobs, each inserting prices for 10 chains, with 20 updates per chain.
	expectedPrices := make(map[string]TokenPriceUpdate)
	for i := 0; i < numJobs; i++ {
		for addr, updatesPerAddr := range updates {
			lastIndex := len(updatesPerAddr) - 1

			err := orm.InsertTokenPricesForDestChain(ctx, destSelector, int32(i), updatesPerAddr[:lastIndex])
			assert.NoError(t, err)
			err = orm.InsertTokenPricesForDestChain(ctx, destSelector, int32(i), updatesPerAddr[lastIndex:])
			assert.NoError(t, err)

			expectedPrices[addr] = updatesPerAddr[lastIndex]
		}
	}

	// verify number of rows inserted
	numRows := getTokenTableRowCount(t, db)
	assert.Equal(t, numJobs*numAddresses*numUpdatesPerAddress, numRows)

	prices, err := orm.GetTokenPricesByDestChain(ctx, destSelector)
	assert.NoError(t, err)
	// should return 1 price per source chain selector
	assert.Equal(t, numAddresses, len(prices))

	// verify getTokenPrices returns prices of latest timestamp
	for _, price := range prices {
		addr := price.TokenAddr
		assert.Equal(t, expectedPrices[addr].TokenPrice, price.TokenPrice)
	}

	// after the initial inserts, insert new round of prices, 1 price per selector this time
	var combinedUpdates []TokenPriceUpdate
	for addr, updatesPerAddr := range updates {
		combinedUpdates = append(combinedUpdates, updatesPerAddr[0])
		expectedPrices[addr] = updatesPerAddr[0]
	}

	err = orm.InsertTokenPricesForDestChain(ctx, destSelector, 1, combinedUpdates)
	assert.NoError(t, err)
	assert.Equal(t, numJobs*numAddresses*numUpdatesPerAddress+numAddresses, getTokenTableRowCount(t, db))

	prices, err = orm.GetTokenPricesByDestChain(ctx, destSelector)
	assert.NoError(t, err)
	assert.Equal(t, numAddresses, len(prices))

	for _, price := range prices {
		addr := price.TokenAddr
		assert.Equal(t, expectedPrices[addr].TokenPrice, price.TokenPrice)
	}
}

func TestORM_InsertAndDeleteTokenPrices(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	orm, db := setupORM(t)

	numAddresses := 10
	numUpdatesPerAddress := 20
	destSelector := uint64(1)

	addrs := generateTokenAddresses(numAddresses)

	updates := make(map[string][]TokenPriceUpdate)
	for _, addr := range addrs {
		updates[addr] = generateTokenPriceUpdates(addr, numUpdatesPerAddress)
	}

	for _, updatesPerAddr := range updates {
		err := orm.InsertTokenPricesForDestChain(ctx, destSelector, 1, updatesPerAddr)
		assert.NoError(t, err)
	}

	sleepSec := 2
	time.Sleep(time.Duration(sleepSec) * time.Second)

	// insert for the 2nd time after interimTimeStamp
	for _, updatesPerAddr := range updates {
		err := orm.InsertTokenPricesForDestChain(ctx, destSelector, 1, updatesPerAddr)
		assert.NoError(t, err)
	}

	assert.Equal(t, 2*numAddresses*numUpdatesPerAddress, getTokenTableRowCount(t, db))

	// clear by sleepSec should delete rows inserted before it
	err := orm.ClearTokenPricesByDestChain(ctx, destSelector, sleepSec)
	assert.NoError(t, err)
	assert.Equal(t, numAddresses*numUpdatesPerAddress, getTokenTableRowCount(t, db))

	// clear by 0 expiration seconds should delete all rows
	err = orm.ClearTokenPricesByDestChain(ctx, destSelector, 0)
	assert.NoError(t, err)
	assert.Equal(t, 0, getTokenTableRowCount(t, db))
}
