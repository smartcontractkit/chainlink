package rollups

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/common/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
)

//go:generate mockery --quiet --name ethClient --output ./mocks/ --case=underscore --structname ETHClient
type ethClient interface {
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
}

// Reads L2-specific precompiles and caches the l1GasPrice set by the L2.
type l1GasPriceOracle struct {
	services.StateMachine
	client     ethClient
	pollPeriod time.Duration
	logger     logger.SugaredLogger
	address    string
	callArgs   string

	l1GasPriceMu sync.RWMutex
	l1GasPrice   *assets.Wei

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
	ArbGasInfo_getL1BaseFeeEstimate = "f5d6ded7"

	// GasOracleAddress is the address of the precompiled contract that exists on OP stack chain.
	// This is the case for Optimism and Base.
	OPGasOracleAddress = "0x420000000000000000000000000000000000000F"
	// GasOracle_l1BaseFee is the a hex encoded call to:
	// `function l1BaseFee() external view returns (uint256);`
	OPGasOracle_l1BaseFee = "519b4bd3"

	// GasOracleAddress is the address of the precompiled contract that exists on Kroma chain.
	// This is the case for Kroma.
	KromaGasOracleAddress = "0x4200000000000000000000000000000000000005"
	// GasOracle_l1BaseFee is the a hex encoded call to:
	// `function l1BaseFee() external view returns (uint256);`
	KromaGasOracle_l1BaseFee = "519b4bd3"

	// GasOracleAddress is the address of the precompiled contract that exists on scroll chain.
	// This is the case for Scroll.
	ScrollGasOracleAddress = "0x5300000000000000000000000000000000000002"
	// GasOracle_l1BaseFee is the a hex encoded call to:
	// `function l1BaseFee() external view returns (uint256);`
	ScrollGasOracle_l1BaseFee = "519b4bd3"

	// Interval at which to poll for L1BaseFee. A good starting point is the L1 block time.
	PollPeriod = 12 * time.Second
)

var supportedChainTypes = []config.ChainType{config.ChainArbitrum, config.ChainOptimismBedrock, config.ChainKroma, config.ChainScroll}

func IsRollupWithL1Support(chainType config.ChainType) bool {
	return slices.Contains(supportedChainTypes, chainType)
}

func NewL1GasPriceOracle(lggr logger.Logger, ethClient ethClient, chainType config.ChainType) L1Oracle {
	var address, callArgs string
	switch chainType {
	case config.ChainArbitrum:
		address = ArbGasInfoAddress
		callArgs = ArbGasInfo_getL1BaseFeeEstimate
	case config.ChainOptimismBedrock:
		address = OPGasOracleAddress
		callArgs = OPGasOracle_l1BaseFee
	case config.ChainKroma:
		address = KromaGasOracleAddress
		callArgs = KromaGasOracle_l1BaseFee
	case config.ChainScroll:
		address = ScrollGasOracleAddress
		callArgs = ScrollGasOracle_l1BaseFee
	default:
		panic(fmt.Sprintf("Received unspported chaintype %s", chainType))
	}

	lggr.Infow("Initializing L1GasPriceOracle with address and callargs", "address", address, "callArgs", callArgs)

	return &l1GasPriceOracle{
		client:        ethClient,
		pollPeriod:    PollPeriod,
		logger:        logger.Sugared(logger.Named(lggr, fmt.Sprintf("L1GasPriceOracle(%s)", chainType))),
		address:       address,
		callArgs:      callArgs,
		chInitialised: make(chan struct{}),
		chStop:        make(chan struct{}),
		chDone:        make(chan struct{}),
	}
}

func (o *l1GasPriceOracle) Name() string {
	return o.logger.Name()
}

func (o *l1GasPriceOracle) Start(ctx context.Context) error {
	return o.StartOnce(o.Name(), func() error {
		go o.run()
		<-o.chInitialised
		return nil
	})
}
func (o *l1GasPriceOracle) Close() error {
	return o.StopOnce(o.Name(), func() error {
		close(o.chStop)
		<-o.chDone
		return nil
	})
}

func (o *l1GasPriceOracle) HealthReport() map[string]error {
	return map[string]error{o.Name(): o.Healthy()}
}

func (o *l1GasPriceOracle) run() {
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

func (o *l1GasPriceOracle) refresh() (t *time.Timer) {
	t = time.NewTimer(utils.WithJitter(o.pollPeriod))

	ctx, cancel := o.chStop.CtxCancel(evmclient.ContextWithDefaultTimeout())
	defer cancel()

	precompile := common.HexToAddress(o.address)
	b, err := o.client.CallContract(ctx, ethereum.CallMsg{
		To:   &precompile,
		Data: common.Hex2Bytes(o.callArgs),
	}, nil)
	if err != nil {
		o.logger.Errorf("gas oracle contract call failed: %v", err)
		return
	}

	if len(b) != 32 { // returns uint256;
		o.logger.Criticalf("return data length (%d) different than expected (%d)", len(b), 32)
		return
	}
	price := new(big.Int).SetBytes(b)

	o.logger.Infow("Fetching l1GasPrice", "l1GasPrice", price.String())

	o.l1GasPriceMu.Lock()
	defer o.l1GasPriceMu.Unlock()
	o.l1GasPrice = assets.NewWei(price)
	return
}

func (o *l1GasPriceOracle) GasPrice(_ context.Context) (l1GasPrice *assets.Wei, err error) {
	ok := o.IfStarted(func() {
		o.l1GasPriceMu.RLock()
		l1GasPrice = o.l1GasPrice
		o.l1GasPriceMu.RUnlock()
	})
	if !ok {
		return l1GasPrice, errors.New("L1GasPriceOracle is not started; cannot estimate gas")
	}
	if l1GasPrice == nil {
		return l1GasPrice, errors.New("failed to get l1 gas price; gas price not set")
	}
	return
}
