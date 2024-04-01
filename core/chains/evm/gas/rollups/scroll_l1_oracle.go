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
type scrollL1Oracle struct {
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
	// ScrollGasOracleAddress is the address of the precompiled contract that exists on Scroll chain.
	ScrollGasOracleAddress = "0x5300000000000000000000000000000000000002"
	// ScrollGasOracle_l1BaseFee is a hex encoded call to:
	// `function l1BaseFee() external view returns (uint256);`
	ScrollGasOracle_l1BaseFee = "l1BaseFee"
	// ScrollGasOracle_getL1Fee is a hex encoded call to:
	// `function getL1Fee(bytes) external view returns (uint256);`
	ScrollGasOracle_getL1Fee = "getL1Fee"
)

func NewScrollL1GasOracle(lggr logger.Logger, ethClient ethClient) L1Oracle {
	l1Oracle := newScrollL1GasOracle(lggr, ethClient)
	return l1Oracle
}

func newScrollL1GasOracle(lggr logger.Logger, ethClient ethClient) L1Oracle {
	var l1GasPriceAddress, gasPriceMethod, l1GasCostAddress, gasCostMethod string
	var l1GasPriceMethodAbi, l1GasCostMethodAbi abi.ABI
	var gasPriceErr, gasCostErr error

	l1GasPriceAddress = ScrollGasOracleAddress
	gasPriceMethod = ScrollGasOracle_l1BaseFee
	l1GasPriceMethodAbi, gasPriceErr = abi.JSON(strings.NewReader(L1BaseFeeAbiString))
	l1GasCostAddress = ScrollGasOracleAddress
	gasCostMethod = ScrollGasOracle_getL1Fee
	l1GasCostMethodAbi, gasCostErr = abi.JSON(strings.NewReader(GetL1FeeAbiString))

	if gasPriceErr != nil {
		panic(fmt.Sprintf("Failed to parse L1 gas price method ABI for chain: scroll"))
	}
	if gasCostErr != nil {
		panic(fmt.Sprintf("Failed to parse L1 gas cost method ABI for chain: scroll"))
	}

	return &scrollL1Oracle{
		client:     ethClient,
		pollPeriod: PollPeriod,
		logger:     logger.Sugared(logger.Named(lggr, fmt.Sprintf("L1GasOracle(scroll)"))),
		chainType:  "scroll",

		l1GasPriceAddress:   l1GasPriceAddress,
		gasPriceMethod:      gasPriceMethod,
		l1GasPriceMethodAbi: l1GasPriceMethodAbi,
		l1GasCostAddress:    l1GasCostAddress,
		gasCostMethod:       gasCostMethod,
		l1GasCostMethodAbi:  l1GasCostMethodAbi,

		priceReader: nil,

		chInitialised: make(chan struct{}),
		chStop:        make(chan struct{}),
		chDone:        make(chan struct{}),
	}
}

func (o *scrollL1Oracle) Name() string {
	return o.logger.Name()
}

func (o *scrollL1Oracle) Start(ctx context.Context) error {
	return o.StartOnce(o.Name(), func() error {
		go o.run()
		<-o.chInitialised
		return nil
	})
}
func (o *scrollL1Oracle) Close() error {
	return o.StopOnce(o.Name(), func() error {
		close(o.chStop)
		<-o.chDone
		return nil
	})
}

func (o *scrollL1Oracle) HealthReport() map[string]error {
	return map[string]error{o.Name(): o.Healthy()}
}

func (o *scrollL1Oracle) run() {
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
func (o *scrollL1Oracle) refresh() (t *time.Timer) {
	t, err := o.refreshWithError()
	if err != nil {
		o.SvcErrBuffer.Append(err)
	}
	return
}

func (o *scrollL1Oracle) refreshWithError() (t *time.Timer, err error) {
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

func (o *scrollL1Oracle) fetchL1GasPrice(ctx context.Context) (price *big.Int, err error) {
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

func (o *scrollL1Oracle) GasPrice(_ context.Context) (l1GasPrice *assets.Wei, err error) {
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
func (o *scrollL1Oracle) GetGasCost(ctx context.Context, tx *gethtypes.Transaction, blockNum *big.Int) (*assets.Wei, error) {
	ctx, cancel := context.WithTimeout(ctx, client.QueryTimeout)
	defer cancel()
	var callData, b []byte
	var err error
	if o.chainType == config.ChainOptimismBedrock || o.chainType == config.ChainScroll {
		// Append rlp-encoded tx
		var encodedtx []byte
		if encodedtx, err = tx.MarshalBinary(); err != nil {
			return nil, fmt.Errorf("failed to marshal tx for gas cost estimation: %w", err)
		}
		if callData, err = o.l1GasCostMethodAbi.Pack(o.gasCostMethod, encodedtx); err != nil {
			return nil, fmt.Errorf("failed to pack calldata for %s L1 gas cost estimation method: %w", o.chainType, err)
		}
	} else if o.chainType == config.ChainArbitrum {
		if callData, err = o.l1GasCostMethodAbi.Pack(o.gasCostMethod, tx.To(), false, tx.Data()); err != nil {
			return nil, fmt.Errorf("failed to pack calldata for %s L1 gas cost estimation method: %w", o.chainType, err)
		}
	} else {
		return nil, fmt.Errorf("L1 gas cost not supported for this chain: %s", o.chainType)
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
	if o.chainType == config.ChainOptimismBedrock || o.chainType == config.ChainScroll {
		if len(b) != 32 { // returns uint256;
			errorMsg := fmt.Sprintf("return data length (%d) different than expected (%d)", len(b), 32)
			o.logger.Critical(errorMsg)
			return nil, fmt.Errorf(errorMsg)
		}
		l1GasCost = new(big.Int).SetBytes(b)
	} else if o.chainType == config.ChainArbitrum {
		if len(b) != 8+2*32 { // returns (uint64 gasEstimateForL1, uint256 baseFee, uint256 l1BaseFeeEstimate);
			errorMsg := fmt.Sprintf("return data length (%d) different than expected (%d)", len(b), 8+2*32)
			o.logger.Critical(errorMsg)
			return nil, fmt.Errorf(errorMsg)
		}
		l1GasCost = new(big.Int).SetBytes(b[:8])
	}

	return assets.NewWei(l1GasCost), nil
}
