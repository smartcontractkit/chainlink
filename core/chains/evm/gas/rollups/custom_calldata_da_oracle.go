package rollups

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
)

type customCalldataDAOracle struct {
	services.StateMachine
	client     l1OracleClient
	pollPeriod time.Duration
	logger     logger.SugaredLogger

	daOracleAddress  common.Address
	daCustomCalldata string
	daGasPriceMu     sync.RWMutex
	daGasPrice       priceEntry

	chInitialized chan struct{}
	chStop        services.StopChan
	chDone        chan struct{}
}

// NewCustomCalldataDAOracle creates a new custom calldata DA oracle. The CustomCalldataDAOracle fetches gas price from
// whatever function is specified in the DAOracle's CustomGasPriceCalldata field. This allows for more flexibility when
// chains have custom DA gas calculation methods.
func NewCustomCalldataDAOracle(lggr logger.Logger, ethClient l1OracleClient, chainType chaintype.ChainType, daOracleConfig evmconfig.DAOracle) (*customCalldataDAOracle, error) {
	if daOracleConfig.OracleType() == nil {
		return nil, errors.New("OracleType is required for CustomCalldataDAOracle but was nil")
	}
	if *daOracleConfig.OracleType() != toml.DAOracleCustomCalldata {
		return nil, fmt.Errorf("expected %s oracle type, got %s", toml.DAOracleCustomCalldata, *daOracleConfig.OracleType())
	}
	if daOracleConfig.OracleAddress() == nil || *daOracleConfig.OracleAddress() == "" {
		return nil, errors.New("OracleAddress is required for CustomCalldataDAOracle but was nil or empty")
	}
	if daOracleConfig.CustomGasPriceCalldata() == nil || *daOracleConfig.CustomGasPriceCalldata() == "" {
		return nil, errors.New("CustomGasPriceCalldata is required")
	}
	oracleAddress := *daOracleConfig.OracleAddress()
	return &customCalldataDAOracle{
		client:     ethClient,
		pollPeriod: PollPeriod,
		logger:     logger.Sugared(logger.Named(lggr, fmt.Sprintf("CustomCalldataDAOracle(%s)", chainType))),

		daOracleAddress:  oracleAddress.Address(),
		daCustomCalldata: *daOracleConfig.CustomGasPriceCalldata(),

		chInitialized: make(chan struct{}),
		chStop:        make(chan struct{}),
		chDone:        make(chan struct{}),
	}, nil
}

func (o *customCalldataDAOracle) Name() string {
	return o.logger.Name()
}

func (o *customCalldataDAOracle) Start(_ context.Context) error {
	return o.StartOnce(o.Name(), func() error {
		go o.run()
		<-o.chInitialized
		return nil
	})
}

func (o *customCalldataDAOracle) Close() error {
	return o.StopOnce(o.Name(), func() error {
		close(o.chStop)
		<-o.chDone
		return nil
	})
}

func (o *customCalldataDAOracle) HealthReport() map[string]error {
	return map[string]error{o.Name(): o.Healthy()}
}

func (o *customCalldataDAOracle) run() {
	defer close(o.chDone)

	o.refresh()
	close(o.chInitialized)

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

func (o *customCalldataDAOracle) refresh() {
	err := o.refreshWithError()
	if err != nil {
		o.logger.Criticalw("Failed to refresh gas price", "err", err)
		o.SvcErrBuffer.Append(err)
	}
}

func (o *customCalldataDAOracle) refreshWithError() error {
	ctx, cancel := o.chStop.CtxWithTimeout(client.QueryTimeout)
	defer cancel()

	price, err := o.getCustomCalldataGasPrice(ctx)
	if err != nil {
		return err
	}

	o.daGasPriceMu.Lock()
	defer o.daGasPriceMu.Unlock()
	o.daGasPrice = priceEntry{price: assets.NewWei(price), timestamp: time.Now()}
	return nil
}

func (o *customCalldataDAOracle) GasPrice(_ context.Context) (daGasPrice *assets.Wei, err error) {
	var timestamp time.Time
	ok := o.IfStarted(func() {
		o.daGasPriceMu.RLock()
		daGasPrice = o.daGasPrice.price
		timestamp = o.daGasPrice.timestamp
		o.daGasPriceMu.RUnlock()
	})
	if !ok {
		return daGasPrice, errors.New("DAGasOracle is not started; cannot estimate gas")
	}
	if daGasPrice == nil {
		return daGasPrice, errors.New("failed to get DA gas price; gas price not set")
	}
	// Validate the price has been updated within the pollPeriod * 2
	// Allowing double the poll period before declaring the price stale to give ample time for the refresh to process
	if time.Since(timestamp) > o.pollPeriod*2 {
		return daGasPrice, errors.New("gas price is stale")
	}
	return
}

func (o *customCalldataDAOracle) getCustomCalldataGasPrice(ctx context.Context) (*big.Int, error) {
	calldata := strings.TrimPrefix(o.daCustomCalldata, "0x")
	calldataBytes, err := hex.DecodeString(calldata)
	if err != nil {
		return nil, fmt.Errorf("failed to decode custom fee method calldata: %w", err)
	}

	b, err := o.client.CallContract(ctx, ethereum.CallMsg{
		To:   &o.daOracleAddress,
		Data: calldataBytes,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("custom fee method call failed: %w", err)
	}
	return new(big.Int).SetBytes(b), nil
}
