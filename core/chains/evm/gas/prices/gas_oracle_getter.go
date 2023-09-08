package prices

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type OracleType string

const (
	ARBITRUM OracleType = "ARBITRUM"
	OP_STACK OracleType = "OP_STACK"
)

type gasOracleGetter struct {
	client             evmclient.Client
	logger             logger.Logger
	address            string
	selector           string
	componentPriceType gas.PriceType

	componentPriceMu sync.RWMutex
	componentPrice   *assets.Wei
}

const (
	// ArbGasInfoAddress is the address of the "Precompiled contract that exists in every Arbitrum chain."
	// https://github.com/OffchainLabs/nitro/blob/f7645453cfc77bf3e3644ea1ac031eff629df325/contracts/src/precompiles/ArbGasInfo.sol
	ArbGasInfoAddress = "0x000000000000000000000000000000000000006C"
	// ArbGasInfo_getL1BaseFeeEstimate is the a hex encoded call to:
	// `function getL1BaseFeeEstimate() external view returns (uint256);`
	ArbGasInfo_getL1BaseFeeEstimate = "f5d6ded7"

	// GasOracleAddress is the address of the "Precompiled contract that exists in every OP stack chain."
	OPGasOracleAddress = "0x420000000000000000000000000000000000000F"
	// GasOracle_l1BaseFee is the a hex encoded call to:
	// `function l1BaseFee() external view returns (uint256);`
	OPGasOracle_getL1BaseFeeEstimate = "519b4bd3"
)

func NewGasOracleGetter(lggr logger.Logger, ethClient evmclient.Client, oracleType OracleType, componentPriceType gas.PriceType) gas.PriceComponentGetter {
	var address, selector string
	switch oracleType {
	case ARBITRUM:
		address = ArbGasInfoAddress
		selector = ArbGasInfo_getL1BaseFeeEstimate
	case OP_STACK:
		address = OPGasOracleAddress
		selector = OPGasOracle_getL1BaseFeeEstimate
	default:
		panic(fmt.Errorf("unsupportd oracle type: %s", oracleType))
	}

	return &gasOracleGetter{
		client:             ethClient,
		logger:             lggr.Named("NewGasOracleGetter"),
		address:            address,
		selector:           selector,
		componentPriceType: componentPriceType,
	}
}

func (o *gasOracleGetter) RefreshComponents(ctx context.Context) (err error) {
	precompile := common.HexToAddress(o.address)
	b, err := o.client.CallContract(ctx, ethereum.CallMsg{
		To:   &precompile,
		Data: common.Hex2Bytes(o.selector),
	}, big.NewInt(-1))
	if err != nil {
		return err
	}

	if len(b) != 32 { // returns uint256;
		return fmt.Errorf("return data length (%d) different than expected (%d)", len(b), 32)
	}
	price := new(big.Int).SetBytes(b)

	o.componentPriceMu.Lock()
	defer o.componentPriceMu.Unlock()
	o.componentPrice = assets.NewWei(price)

	return
}

func (o *gasOracleGetter) GetPriceComponents(_ context.Context, gasPrice *assets.Wei) (prices []gas.PriceComponent, err error) {
	o.componentPriceMu.RLock()
	defer o.componentPriceMu.RUnlock()
	return []gas.PriceComponent{
		{
			Price:     gasPrice,
			PriceType: gas.GAS_PRICE,
		},
		{
			Price:     o.componentPrice,
			PriceType: o.componentPriceType,
		},
	}, nil
}
