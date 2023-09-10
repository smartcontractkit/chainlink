package chainoracles

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type OracleType string

const (
	Arbitrum OracleType = "ARBITRUM"
	OPStack  OracleType = "OP_STACK"
)

// Reads L2-specific precompiles and caches the l1BaseFee set by the L2.
type l1BaseFeeOracle struct {
	client     evmclient.Client
	pollPeriod time.Duration
	logger     logger.Logger
	address    string
	selector   string

	l1BaseFeeMu sync.RWMutex
	l1BaseFee   *assets.Wei

	chInitialised chan struct{}
	chStop        utils.StopChan
	chDone        chan struct{}
	utils.StartStopOnce
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
	OPGasOracle_l1BaseFee = "519b4bd3"

	// Interval at which to poll for L1BaseFee. A good starting point is the L1 block time.
	PollPeriod = 12 * time.Second
)

func NewL1BaeFeeOracle(lggr logger.Logger, ethClient evmclient.Client, oracleType OracleType) L1Oracle {
	var address, selector string
	switch oracleType {
	case Arbitrum:
		address = ArbGasInfoAddress
		selector = ArbGasInfo_getL1BaseFeeEstimate
	case OPStack:
		address = OPGasOracleAddress
		selector = OPGasOracle_l1BaseFee
	default:
		panic(fmt.Errorf("unsupportd oracle type: %s", oracleType))
	}

	return &l1BaseFeeOracle{
		client:     ethClient,
		pollPeriod: PollPeriod,
		logger:     lggr.Named(fmt.Sprintf("%s L1BaseFeeOracle", oracleType)),
		address:    address,
		selector:   selector,
	}
}

func (o *l1BaseFeeOracle) Name() string {
	return o.logger.Name()
}

func (o *l1BaseFeeOracle) Start(ctx context.Context) error {
	return o.StartOnce(o.Name(), func() error {
		go o.run()
		<-o.chInitialised
		return nil
	})
}
func (o *l1BaseFeeOracle) Close() error {
	return o.StopOnce(o.Name(), func() (err error) {
		close(o.chStop)
		<-o.chDone
		return
	})
}

func (o *l1BaseFeeOracle) Ready() error { return o.StartStopOnce.Ready() }

func (o *l1BaseFeeOracle) HealthReport() map[string]error {
	return map[string]error{o.Name(): o.StartStopOnce.Healthy()}
}

func (o *l1BaseFeeOracle) run() {
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

func (o *l1BaseFeeOracle) refresh() (t *time.Timer) {
	t = time.NewTimer(utils.WithJitter(o.pollPeriod))

	ctx, cancel := o.chStop.CtxCancel(evmclient.ContextWithDefaultTimeout())
	defer cancel()

	precompile := common.HexToAddress(o.address)
	b, err := o.client.CallContract(ctx, ethereum.CallMsg{
		To:   &precompile,
		Data: common.Hex2Bytes(o.selector),
	}, big.NewInt(-1))
	if err != nil {
		return
	}

	if len(b) != 32 { // returns uint256;
		o.logger.Warnf("return data length (%d) different than expected (%d)", len(b), 32)
		return
	}
	price := new(big.Int).SetBytes(b)

	o.l1BaseFeeMu.Lock()
	defer o.l1BaseFeeMu.Unlock()
	o.l1BaseFee = assets.NewWei(price)
	return
}

func (o *l1BaseFeeOracle) L1GasPrice(_ context.Context) (*assets.Wei, error) {
	o.l1BaseFeeMu.RLock()
	defer o.l1BaseFeeMu.RUnlock()
	return o.l1BaseFee, nil
}
