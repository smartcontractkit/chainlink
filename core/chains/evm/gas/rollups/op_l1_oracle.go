package rollups

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	gethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/v2/common/client"
	"github.com/smartcontractkit/chainlink/v2/common/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
)

// Reads L2-specific precompiles and caches the l1GasPrice set by the L2.
type optimismL1Oracle struct {
	services.StateMachine
	client     ethClient
	pollPeriod time.Duration
	logger     logger.SugaredLogger
	chainType  config.ChainType

	l1GasPriceAddress   string
	gasPriceMethod      string
	l1GasPriceMethodAbi abi.ABI
	l1GasPriceMu        sync.RWMutex
	l1GasPrice          priceEntry

	l1GasCostAddress   string
	gasCostMethod      string
	l1GasCostMethodAbi abi.ABI

	priceReader daPriceReader

	chInitialised chan struct{}
	chStop        services.StopChan
	chDone        chan struct{}
}

const (
	// OPGasOracleAddress is the address of the precompiled contract that exists on OP stack chain.
	// This is the case for Optimism and Base.
	OPGasOracleAddress = "0x420000000000000000000000000000000000000F"
	// OPGasOracle_l1BaseFee is a hex encoded call to:
	// `function l1BaseFee() external view returns (uint256);`
	OPGasOracle_l1BaseFee = "l1BaseFee"
	// OPGasOracle_getL1Fee is a hex encoded call to:
	// `function getL1Fee(bytes) external view returns (uint256);`
	OPGasOracle_getL1Fee = "getL1Fee"

	// GasOracleAddress is the address of the precompiled contract that exists on Kroma chain.
	// This is the case for Kroma.
	KromaGasOracleAddress = "0x4200000000000000000000000000000000000005"
	// GasOracle_l1BaseFee is the a hex encoded call to:
	// `function l1BaseFee() external view returns (uint256);`
	KromaGasOracle_l1BaseFee = "l1BaseFee"
)

func NewOpStackL1GasOracle(lggr logger.Logger, ethClient ethClient, chainType config.ChainType) L1Oracle {
	var precompileAddress string
	switch chainType {
	case config.ChainOptimismBedrock:
		precompileAddress = OPGasOracleAddress
	case config.ChainKroma:
		precompileAddress = KromaGasOracleAddress
	default:
		panic(fmt.Sprintf("Received unspported chaintype %s", chainType))
	}
	priceReader := newOPPriceReader(lggr, ethClient, chainType, precompileAddress)
	return newOpStackL1GasOracle(lggr, ethClient, priceReader, chainType)
}

func newOpStackL1GasOracle(lggr logger.Logger, ethClient ethClient, priceReader daPriceReader, chainType config.ChainType) L1Oracle {
	var l1GasPriceAddress, gasPriceMethod, l1GasCostAddress, gasCostMethod string
	var l1GasPriceMethodAbi, l1GasCostMethodAbi abi.ABI
	var gasPriceErr, gasCostErr error

	l1GasPriceAddress = OPGasOracleAddress
	gasPriceMethod = OPGasOracle_l1BaseFee
	l1GasPriceMethodAbi, gasPriceErr = abi.JSON(strings.NewReader(L1BaseFeeAbiString))
	l1GasCostAddress = OPGasOracleAddress
	gasCostMethod = OPGasOracle_getL1Fee
	l1GasCostMethodAbi, gasCostErr = abi.JSON(strings.NewReader(GetL1FeeAbiString))

	if gasPriceErr != nil {
		panic(fmt.Sprintf("Failed to parse L1 gas price method ABI for chain: optimismBedrock"))
	}
	if gasCostErr != nil {
		panic(fmt.Sprintf("Failed to parse L1 gas cost method ABI for chain: optimismBedrock"))
	}

	return &optimismL1Oracle{
		client:     ethClient,
		pollPeriod: PollPeriod,
		logger:     logger.Sugared(logger.Named(lggr, fmt.Sprintf("L1GasOracle(optimismBedrock)"))),
		chainType:  "optimismBedrock",

		l1GasPriceAddress:   l1GasPriceAddress,
		gasPriceMethod:      gasPriceMethod,
		l1GasPriceMethodAbi: l1GasPriceMethodAbi,
		l1GasCostAddress:    l1GasCostAddress,
		gasCostMethod:       gasCostMethod,
		l1GasCostMethodAbi:  l1GasCostMethodAbi,

		priceReader: priceReader,

		chInitialised: make(chan struct{}),
		chStop:        make(chan struct{}),
		chDone:        make(chan struct{}),
	}
}

func (o *optimismL1Oracle) Name() string {
	return o.logger.Name()
}

func (o *optimismL1Oracle) Start(ctx context.Context) error {
	return o.StartOnce(o.Name(), func() error {
		go o.run()
		<-o.chInitialised
		return nil
	})
}
func (o *optimismL1Oracle) Close() error {
	return o.StopOnce(o.Name(), func() error {
		close(o.chStop)
		<-o.chDone
		return nil
	})
}

func (o *optimismL1Oracle) HealthReport() map[string]error {
	return map[string]error{o.Name(): o.Healthy()}
}

func (o *optimismL1Oracle) run() {
	defer close(o.chDone)

	t := o.refresh()
	close(o.chInitialised)

	for {
		select {
		case <-o.chStop:
			return
		case <-t.C:
			t = o.refresh()
		}
	}
}
func (o *optimismL1Oracle) refresh() (t *time.Timer) {
	t, err := o.refreshWithError()
	if err != nil {
		o.SvcErrBuffer.Append(err)
	}
	return
}

func (o *optimismL1Oracle) refreshWithError() (t *time.Timer, err error) {
	t = time.NewTimer(utils.WithJitter(o.pollPeriod))

	ctx, cancel := o.chStop.CtxCancel(evmclient.ContextWithDefaultTimeout())
	defer cancel()

	price, err := o.fetchL1GasPrice(ctx)
	if err != nil {
		return t, err
	}

	o.l1GasPriceMu.Lock()
	defer o.l1GasPriceMu.Unlock()
	o.l1GasPrice = priceEntry{price: assets.NewWei(price), timestamp: time.Now()}
	return
}

func (o *optimismL1Oracle) fetchL1GasPrice(ctx context.Context) (price *big.Int, err error) {
	// if dedicated priceReader exists, use the reader
	if o.priceReader != nil {
		return o.priceReader.GetDAGasPrice(ctx)
	}

	var callData, b []byte
	precompile := common.HexToAddress(o.l1GasPriceAddress)
	callData, err = o.l1GasPriceMethodAbi.Pack(o.gasPriceMethod)
	if err != nil {
		errMsg := fmt.Sprintf("failed to pack calldata for %s L1 gas price method", o.chainType)
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

func (o *optimismL1Oracle) GasPrice(_ context.Context) (l1GasPrice *assets.Wei, err error) {
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
func (o *optimismL1Oracle) GetGasCost(ctx context.Context, tx *gethtypes.Transaction, blockNum *big.Int) (*assets.Wei, error) {
	ctx, cancel := context.WithTimeout(ctx, client.QueryTimeout)
	defer cancel()
	var callData, b []byte
	var err error
	// Append rlp-encoded tx
	var encodedtx []byte
	if encodedtx, err = tx.MarshalBinary(); err != nil {
		return nil, fmt.Errorf("failed to marshal tx for gas cost estimation: %w", err)
	}
	if callData, err = o.l1GasCostMethodAbi.Pack(o.gasCostMethod, encodedtx); err != nil {
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
	if len(b) != 32 { // returns uint256;
		errorMsg := fmt.Sprintf("return data length (%d) different than expected (%d)", len(b), 32)
		o.logger.Critical(errorMsg)
		return nil, fmt.Errorf(errorMsg)
	}
	l1GasCost = new(big.Int).SetBytes(b)

	return assets.NewWei(l1GasCost), nil
}
