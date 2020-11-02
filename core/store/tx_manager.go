package store

import (
	"math/big"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store/orm"

	"github.com/ethereum/go-ethereum/common"
)

// GetLINKBalance returns the balance of LINK at the given address
func GetLINKBalance(config orm.ConfigReader, client eth.Client, address common.Address) (*assets.Link, error) {
	contractAddress := common.HexToAddress(config.LinkContractAddress())
	balance, err := client.GetERC20Balance(address, contractAddress)
	if err != nil {
		return assets.NewLink(0), err
	}
	return (*assets.Link)(balance), nil
}

// BumpGas computes the next gas price to attempt as the largest of:
// - A configured percentage bump (ETH_GAS_BUMP_PERCENT) on top of the baseline price.
// - A configured fixed amount of Wei (ETH_GAS_PRICE_WEI) on top of the baseline price.
// The baseline price is the maximum of the previous gas price attempt and the node's current gas price.
func BumpGas(config orm.ConfigReader, originalGasPrice *big.Int) (*big.Int, error) {
	baselinePrice := max(originalGasPrice, config.EthGasPriceDefault())

	var priceByPercentage = new(big.Int)
	priceByPercentage.Mul(baselinePrice, big.NewInt(int64(100+config.EthGasBumpPercent())))
	priceByPercentage.Div(priceByPercentage, big.NewInt(100))

	var priceByIncrement = new(big.Int)
	priceByIncrement.Add(baselinePrice, config.EthGasBumpWei())

	bumpedGasPrice := max(priceByPercentage, priceByIncrement)
	if bumpedGasPrice.Cmp(config.EthMaxGasPriceWei()) > 0 {
		return config.EthMaxGasPriceWei(), errors.Errorf("bumped gas price of %s would exceed configured max gas price of %s (original price was %s)",
			bumpedGasPrice.String(), config.EthMaxGasPriceWei(), originalGasPrice.String())
	} else if bumpedGasPrice.Cmp(originalGasPrice) == 0 {
		// NOTE: This really shouldn't happen since we enforce minimums for
		// ETH_GAS_BUMP_PERCENT and ETH_GAS_BUMP_WEI in the config validation,
		// but it's here anyway for a "belts and braces" approach
		return bumpedGasPrice, errors.Errorf("bumped gas price of %s is equal to original gas price of %s."+
			" ACTION REQUIRED: This is a configuration error, you must increase either "+
			"ETH_GAS_BUMP_PERCENT or ETH_GAS_BUMP_WEI", bumpedGasPrice.String(), originalGasPrice.String())
	}
	return bumpedGasPrice, nil
}

func max(a, b *big.Int) *big.Int {
	if a.Cmp(b) >= 0 {
		return a
	}
	return b
}
