package gas

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type ArbConfig interface {
	EvmGasLimitMax() uint32
}

//go:generate mockery --quiet --name ethClient --output ./mocks/ --case=underscore --structname ETHClient
type ethClient interface {
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
}

// arbitrumEstimator is an Estimator which extends l2SuggestedPriceEstimator to use getPricesInArbGas() for gas limit estimation.
type arbitrumEstimator struct {
	cfg ArbConfig

	EvmEstimator // *l2SuggestedPriceEstimator

	client     ethClient
	pollPeriod time.Duration
	logger     logger.Logger

	getPricesInArbGasMu sync.RWMutex
	perL2Tx             uint32
	perL1CalldataUnit   uint32

	chForceRefetch chan (chan struct{})
	chInitialised  chan struct{}
	chStop         utils.StopChan
	chDone         chan struct{}

	utils.StartStopOnce
}

func NewArbitrumEstimator(lggr logger.Logger, cfg ArbConfig, rpcClient rpcClient, ethClient ethClient) EvmEstimator {
	lggr = lggr.Named("ArbitrumEstimator")
	return &arbitrumEstimator{
		cfg:            cfg,
		EvmEstimator:   NewL2SuggestedPriceEstimator(lggr, rpcClient),
		client:         ethClient,
		pollPeriod:     10 * time.Second,
		logger:         lggr,
		chForceRefetch: make(chan (chan struct{})),
		chInitialised:  make(chan struct{}),
		chStop:         make(chan struct{}),
		chDone:         make(chan struct{}),
	}
}

func (a *arbitrumEstimator) Name() string {
	return a.logger.Name()
}

func (a *arbitrumEstimator) Start(ctx context.Context) error {
	return a.StartOnce("ArbitrumEstimator", func() error {
		if err := a.EvmEstimator.Start(ctx); err != nil {
			return errors.Wrap(err, "failed to start gas price estimator")
		}
		go a.run()
		<-a.chInitialised
		return nil
	})
}
func (a *arbitrumEstimator) Close() error {
	return a.StopOnce("ArbitrumEstimator", func() (err error) {
		close(a.chStop)
		err = errors.Wrap(a.EvmEstimator.Close(), "failed to stop gas price estimator")
		<-a.chDone
		return
	})
}

func (a *arbitrumEstimator) Ready() error { return a.StartStopOnce.Ready() }

func (a *arbitrumEstimator) HealthReport() map[string]error {
	return map[string]error{a.Name(): a.StartStopOnce.Healthy()}
}

// GetLegacyGas estimates both the gas price and the gas limit.
//   - Price is delegated to the embedded l2SuggestedPriceEstimator.
//   - Limit is computed from the dynamic values perL2Tx and perL1CalldataUnit, provided by the getPricesInArbGas() method
//     of the precompilie contract at ArbGasInfoAddress. perL2Tx is a constant amount of gas, and perL1CalldataUnit is
//     multiplied by the length of the tx calldata. The sum of these two values plus the original l2GasLimit is returned.
func (a *arbitrumEstimator) GetLegacyGas(ctx context.Context, calldata []byte, l2GasLimit uint32, maxGasPriceWei *assets.Wei, opts ...txmgrtypes.Opt) (gasPrice *assets.Wei, chainSpecificGasLimit uint32, err error) {
	gasPrice, _, err = a.EvmEstimator.GetLegacyGas(ctx, calldata, l2GasLimit, maxGasPriceWei, opts...)
	if err != nil {
		return
	}
	gasPrice = a.gasPriceWithBuffer(gasPrice, maxGasPriceWei)
	ok := a.IfStarted(func() {
		if slices.Contains(opts, txmgrtypes.OptForceRefetch) {
			ch := make(chan struct{})
			select {
			case a.chForceRefetch <- ch:
			case <-a.chStop:
				err = errors.New("estimator stopped")
				return
			case <-ctx.Done():
				err = ctx.Err()
				return
			}
			select {
			case <-ch:
			case <-a.chStop:
				err = errors.New("estimator stopped")
				return
			case <-ctx.Done():
				err = ctx.Err()
				return
			}
		}
		perL2Tx, perL1CalldataUnit := a.getPricesInArbGas()
		chainSpecificGasLimit = l2GasLimit + perL2Tx + uint32(len(calldata))*perL1CalldataUnit
		a.logger.Debugw("GetLegacyGas", "l2GasLimit", l2GasLimit, "calldataLen", len(calldata), "perL2Tx", perL2Tx,
			"perL1CalldataUnit", perL1CalldataUnit, "chainSpecificGasLimit", chainSpecificGasLimit)
	})
	if !ok {
		return nil, 0, errors.New("estimator is not started")
	} else if err != nil {
		return
	}
	if max := a.cfg.EvmGasLimitMax(); chainSpecificGasLimit > max {
		err = fmt.Errorf("estimated gas limit: %d is greater than the maximum gas limit configured: %d", chainSpecificGasLimit, max)
		return
	}
	return
}

// During network congestion Arbitrum's suggested gas price can be extremely volatile, making gas estimations less accurate. For any transaction, Arbitrum will only charge
// the block's base fee. If the base fee increases rapidly there is a chance the suggested gas price will fall under that value, resulting in a fee too low error.
// We use gasPriceWithBuffer to increase the estimated gas price by some percentage to avoid fee too low errors. Eventually, only the base fee will be paid, regardless of the price.
func (a *arbitrumEstimator) gasPriceWithBuffer(gasPrice *assets.Wei, maxGasPriceWei *assets.Wei) *assets.Wei {
	const gasPriceBufferPercentage = 50

	gasPrice = gasPrice.AddPercentage(gasPriceBufferPercentage)
	if gasPrice.Cmp(maxGasPriceWei) > 0 {
		a.logger.Warnw("Updated gasPrice with buffer is higher than the max gas price limit. Falling back to max gas price", "gasPriceWithBuffer", gasPrice, "maxGasPriceWei", maxGasPriceWei)
		gasPrice = maxGasPriceWei
	}
	a.logger.Debugw("gasPriceWithBuffer", "updatedGasPrice", gasPrice)
	return gasPrice
}

func (a *arbitrumEstimator) getPricesInArbGas() (perL2Tx uint32, perL1CalldataUnit uint32) {
	a.getPricesInArbGasMu.RLock()
	perL2Tx, perL1CalldataUnit = a.perL2Tx, a.perL1CalldataUnit
	a.getPricesInArbGasMu.RUnlock()
	return
}

func (a *arbitrumEstimator) run() {
	defer close(a.chDone)

	t := a.refreshPricesInArbGas()
	close(a.chInitialised)

	for {
		select {
		case <-a.chStop:
			return
		case ch := <-a.chForceRefetch:
			t.Stop()
			t = a.refreshPricesInArbGas()
			close(ch)
		case <-t.C:
			t = a.refreshPricesInArbGas()
		}
	}
}

// refreshPricesInArbGas calls getPricesInArbGas() and caches the refreshed prices.
func (a *arbitrumEstimator) refreshPricesInArbGas() (t *time.Timer) {
	t = time.NewTimer(utils.WithJitter(a.pollPeriod))

	perL2Tx, perL1CalldataUnit, err := a.callGetPricesInArbGas()
	if err != nil {
		a.logger.Warnw("Failed to refresh prices", "err", err)
		return
	}

	a.logger.Debugw("refreshPricesInArbGas", "perL2Tx", perL2Tx, "perL2CalldataUnit", perL1CalldataUnit)

	a.getPricesInArbGasMu.Lock()
	a.perL2Tx = perL2Tx
	a.perL1CalldataUnit = perL1CalldataUnit
	a.getPricesInArbGasMu.Unlock()
	return
}

const (
	// ArbGasInfoAddress is the address of the "Precompiled contract that exists in every Arbitrum chain."
	// https://github.com/OffchainLabs/nitro/blob/f7645453cfc77bf3e3644ea1ac031eff629df325/contracts/src/precompiles/ArbGasInfo.sol
	ArbGasInfoAddress = "0x000000000000000000000000000000000000006C"
	// ArbGasInfo_getPricesInArbGas is the a hex encoded call to:
	// `function getPricesInArbGas() external view returns (uint256, uint256, uint256);`
	ArbGasInfo_getPricesInArbGas = "02199f34"
)

// callGetPricesInArbGas calls ArbGasInfo.getPricesInArbGas() on the precompile contract ArbGasInfoAddress.
//
// @return (per L2 tx, per L1 calldata unit, per storage allocation)
// function getPricesInArbGas() external view returns (uint256, uint256, uint256);
//
// https://github.com/OffchainLabs/nitro/blob/f7645453cfc77bf3e3644ea1ac031eff629df325/contracts/src/precompiles/ArbGasInfo.sol#L69
func (a *arbitrumEstimator) callGetPricesInArbGas() (perL2Tx uint32, perL1CalldataUnit uint32, err error) {
	ctx, cancel := a.chStop.CtxCancel(evmclient.ContextWithDefaultTimeout())
	defer cancel()

	precompile := common.HexToAddress(ArbGasInfoAddress)
	b, err := a.client.CallContract(ctx, ethereum.CallMsg{
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
