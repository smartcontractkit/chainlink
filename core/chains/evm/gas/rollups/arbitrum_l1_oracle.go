package rollups

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	gethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/common/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
)

type ArbL1GasOracle interface {
	L1Oracle
	GetPricesInArbGas() (perL2Tx uint32, perL1CalldataUnit uint32, err error)
}

// Reads L2-specific precompiles and caches the l1GasPrice set by the L2.
type arbitrumL1Oracle struct {
	services.StateMachine
	client     l1OracleClient
	pollPeriod time.Duration
	logger     logger.SugaredLogger
	chainType  chaintype.ChainType

	l1GasPriceAddress   string
	gasPriceMethod      string
	l1GasPriceMethodAbi abi.ABI
	l1GasPriceMu        sync.RWMutex
	l1GasPrice          priceEntry

	l1GasCostAddress   string
	gasCostMethod      string
	l1GasCostMethodAbi abi.ABI

	chInitialised chan struct{}
	chStop        services.StopChan
	chDone        chan struct{}
}

const (
	// ArbGasInfoAddress is the address of the "Precompiled contract that exists in every Arbitrum chain."
	// https://github.com/OffchainLabs/nitro/blob/f7645453cfc77bf3e3644ea1ac031eff629df325/contracts/src/precompiles/ArbGasInfo.sol
	ArbGasInfoAddress = "0x000000000000000000000000000000000000006C"
	// ArbGasInfo_getL1BaseFeeEstimate is the a hex encoded call to:
	// `function getL1BaseFeeEstimate() external view returns (uint256);`
	ArbGasInfo_getL1BaseFeeEstimate = "getL1BaseFeeEstimate"
	// NodeInterfaceAddress is the address of the precompiled contract that is only available through RPC
	// https://github.com/OffchainLabs/nitro/blob/e815395d2e91fb17f4634cad72198f6de79c6e61/nodeInterface/NodeInterface.go#L37
	ArbNodeInterfaceAddress = "0x00000000000000000000000000000000000000C8"
	// ArbGasInfo_getPricesInArbGas is the a hex encoded call to:
	// `function gasEstimateL1Component(address to, bool contractCreation, bytes calldata data) external payable returns (uint64 gasEstimateForL1, uint256 baseFee, uint256 l1BaseFeeEstimate);`
	ArbNodeInterface_gasEstimateL1Component = "gasEstimateL1Component"
	// ArbGasInfo_getPricesInArbGas is the a hex encoded call to:
	// `function getPricesInArbGas() external view returns (uint256, uint256, uint256);`
	ArbGasInfo_getPricesInArbGas = "02199f34"
)

func NewArbitrumL1GasOracle(lggr logger.Logger, ethClient l1OracleClient) (*arbitrumL1Oracle, error) {
	var l1GasPriceAddress, gasPriceMethod, l1GasCostAddress, gasCostMethod string
	var l1GasPriceMethodAbi, l1GasCostMethodAbi abi.ABI
	var gasPriceErr, gasCostErr error

	l1GasPriceAddress = ArbGasInfoAddress
	gasPriceMethod = ArbGasInfo_getL1BaseFeeEstimate
	l1GasPriceMethodAbi, gasPriceErr = abi.JSON(strings.NewReader(GetL1BaseFeeEstimateAbiString))
	l1GasCostAddress = ArbNodeInterfaceAddress
	gasCostMethod = ArbNodeInterface_gasEstimateL1Component
	l1GasCostMethodAbi, gasCostErr = abi.JSON(strings.NewReader(GasEstimateL1ComponentAbiString))

	if gasPriceErr != nil {
		return nil, fmt.Errorf("failed to parse L1 gas price method ABI for chain: arbitrum: %w", gasPriceErr)
	}
	if gasCostErr != nil {
		return nil, fmt.Errorf("failed to parse L1 gas cost method ABI for chain: arbitrum: %w", gasCostErr)
	}

	return &arbitrumL1Oracle{
		client:     ethClient,
		pollPeriod: PollPeriod,
		logger:     logger.Sugared(logger.Named(lggr, "L1GasOracle(arbitrum)")),
		chainType:  chaintype.ChainArbitrum,

		l1GasPriceAddress:   l1GasPriceAddress,
		gasPriceMethod:      gasPriceMethod,
		l1GasPriceMethodAbi: l1GasPriceMethodAbi,
		l1GasCostAddress:    l1GasCostAddress,
		gasCostMethod:       gasCostMethod,
		l1GasCostMethodAbi:  l1GasCostMethodAbi,

		chInitialised: make(chan struct{}),
		chStop:        make(chan struct{}),
		chDone:        make(chan struct{}),
	}, nil
}

func (o *arbitrumL1Oracle) Name() string {
	return o.logger.Name()
}

func (o *arbitrumL1Oracle) ChainType(_ context.Context) chaintype.ChainType {
	return o.chainType
}

func (o *arbitrumL1Oracle) Start(ctx context.Context) error {
	return o.StartOnce(o.Name(), func() error {
		go o.run()
		<-o.chInitialised
		return nil
	})
}
func (o *arbitrumL1Oracle) Close() error {
	return o.StopOnce(o.Name(), func() error {
		close(o.chStop)
		<-o.chDone
		return nil
	})
}

func (o *arbitrumL1Oracle) HealthReport() map[string]error {
	return map[string]error{o.Name(): o.Healthy()}
}

func (o *arbitrumL1Oracle) run() {
	defer close(o.chDone)

	o.refresh()
	close(o.chInitialised)

	t := services.TickerConfig{
		Initial:   o.pollPeriod,
		JitterPct: services.DefaultJitter,
	}.NewTicker(o.pollPeriod)
	defer t.Stop()

	for {
		select {
		case <-o.chStop:
			return
		case <-t.C:
			o.refresh()
		}
	}
}
func (o *arbitrumL1Oracle) refresh() {
	err := o.refreshWithError()
	if err != nil {
		o.logger.Criticalw("Failed to refresh gas price", "err", err)
		o.SvcErrBuffer.Append(err)
	}
}

func (o *arbitrumL1Oracle) refreshWithError() error {
	ctx, cancel := o.chStop.CtxCancel(evmclient.ContextWithDefaultTimeout())
	defer cancel()

	price, err := o.fetchL1GasPrice(ctx)
	if err != nil {
		return err
	}

	o.l1GasPriceMu.Lock()
	defer o.l1GasPriceMu.Unlock()
	o.l1GasPrice = priceEntry{price: assets.NewWei(price), timestamp: time.Now()}
	return nil
}

func (o *arbitrumL1Oracle) fetchL1GasPrice(ctx context.Context) (price *big.Int, err error) {
	var callData, b []byte
	precompile := common.HexToAddress(o.l1GasPriceAddress)
	callData, err = o.l1GasPriceMethodAbi.Pack(o.gasPriceMethod)
	if err != nil {
		errMsg := "failed to pack calldata for arbitrum L1 gas price method"
		o.logger.Errorf(errMsg)
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}
	b, err = o.client.CallContract(ctx, ethereum.CallMsg{
		To:   &precompile,
		Data: callData,
	}, nil)
	if err != nil {
		errMsg := "gas oracle contract call failed"
		o.logger.Errorf(errMsg)
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	if len(b) != 32 { // returns uint256;
		errMsg := fmt.Sprintf("return data length (%d) different than expected (%d)", len(b), 32)
		o.logger.Criticalf(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	price = new(big.Int).SetBytes(b)
	return price, nil
}

func (o *arbitrumL1Oracle) GasPrice(_ context.Context) (l1GasPrice *assets.Wei, err error) {
	var timestamp time.Time
	ok := o.IfStarted(func() {
		o.l1GasPriceMu.RLock()
		l1GasPrice = o.l1GasPrice.price
		timestamp = o.l1GasPrice.timestamp
		o.l1GasPriceMu.RUnlock()
	})
	if !ok {
		return l1GasPrice, fmt.Errorf("L1GasOracle is not started; cannot estimate gas")
	}
	if l1GasPrice == nil {
		return l1GasPrice, fmt.Errorf("failed to get l1 gas price; gas price not set")
	}
	// Validate the price has been updated within the pollPeriod * 2
	// Allowing double the poll period before declaring the price stale to give ample time for the refresh to process
	if time.Since(timestamp) > o.pollPeriod*2 {
		return l1GasPrice, fmt.Errorf("gas price is stale")
	}
	return
}

// Gets the L1 gas cost for the provided transaction at the specified block num
// If block num is not provided, the value on the latest block num is used
func (o *arbitrumL1Oracle) GetGasCost(ctx context.Context, tx *gethtypes.Transaction, blockNum *big.Int) (*assets.Wei, error) {
	ctx, cancel := context.WithTimeout(ctx, client.QueryTimeout)
	defer cancel()
	var callData, b []byte
	var err error

	if callData, err = o.l1GasCostMethodAbi.Pack(o.gasCostMethod, tx.To(), false, tx.Data()); err != nil {
		return nil, fmt.Errorf("failed to pack calldata for %s L1 gas cost estimation method: %w", o.chainType, err)
	}

	precompile := common.HexToAddress(o.l1GasCostAddress)
	b, err = o.client.CallContract(ctx, ethereum.CallMsg{
		To:   &precompile,
		Data: callData,
	}, blockNum)
	if err != nil {
		errorMsg := fmt.Sprintf("gas oracle contract call failed: %v", err)
		o.logger.Errorf(errorMsg)
		return nil, fmt.Errorf(errorMsg)
	}

	var l1GasCost *big.Int

	if len(b) != 8+2*32 { // returns (uint64 gasEstimateForL1, uint256 baseFee, uint256 l1BaseFeeEstimate);
		errorMsg := fmt.Sprintf("return data length (%d) different than expected (%d)", len(b), 8+2*32)
		o.logger.Critical(errorMsg)
		return nil, fmt.Errorf(errorMsg)
	}
	l1GasCost = new(big.Int).SetBytes(b[:8])

	return assets.NewWei(l1GasCost), nil
}

// callGetPricesInArbGas calls ArbGasInfo.getPricesInArbGas() on the precompile contract ArbGasInfoAddress.
//
// @return (per L2 tx, per L1 calldata unit, per storage allocation)
// function getPricesInArbGas() external view returns (uint256, uint256, uint256);
//
// https://github.com/OffchainLabs/nitro/blob/f7645453cfc77bf3e3644ea1ac031eff629df325/contracts/src/precompiles/ArbGasInfo.sol#L69

func (o *arbitrumL1Oracle) GetPricesInArbGas() (perL2Tx uint32, perL1CalldataUnit uint32, err error) {
	ctx, cancel := o.chStop.CtxCancel(evmclient.ContextWithDefaultTimeout())
	defer cancel()
	precompile := common.HexToAddress(ArbGasInfoAddress)
	b, err := o.client.CallContract(ctx, ethereum.CallMsg{
		To:   &precompile,
		Data: common.Hex2Bytes(ArbGasInfo_getPricesInArbGas),
	}, big.NewInt(-1))
	if err != nil {
		return 0, 0, err
	}

	if len(b) != 3*32 { // returns (uint256, uint256, uint256);
		err = fmt.Errorf("return data length (%d) different than expected (%d)", len(b), 3*32)
		return
	}
	bPerL2Tx := new(big.Int).SetBytes(b[:32])
	bPerL1CalldataUnit := new(big.Int).SetBytes(b[32:64])
	// ignore perStorageAllocation
	if !bPerL2Tx.IsUint64() || !bPerL1CalldataUnit.IsUint64() {
		err = fmt.Errorf("returned integers are not uint64 (%s, %s)", bPerL2Tx.String(), bPerL1CalldataUnit.String())
		return
	}

	perL2TxU64 := bPerL2Tx.Uint64()
	perL1CalldataUnitU64 := bPerL1CalldataUnit.Uint64()
	if perL2TxU64 > math.MaxUint32 || perL1CalldataUnitU64 > math.MaxUint32 {
		err = fmt.Errorf("returned integers are not uint32 (%d, %d)", perL2TxU64, perL1CalldataUnitU64)
		return
	}
	perL2Tx = uint32(perL2TxU64)
	perL1CalldataUnit = uint32(perL1CalldataUnitU64)
	return
}
