package rollups

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	gethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
)

// Reads L2-specific precompiles and caches the l1GasPrice set by the L2.
type zkSyncL1Oracle struct {
	services.StateMachine
	client     l1OracleClient
	pollPeriod time.Duration
	logger     logger.SugaredLogger
	chainType  chaintype.ChainType

	systemContextAddress  string
	gasPerPubdataMethod   string
	gasPerPubdataSelector string
	l2GasPriceMethod      string
	l2GasPriceSelector    string

	l1GasPriceMu sync.RWMutex
	l1GasPrice   priceEntry

	chInitialised chan struct{}
	chStop        services.StopChan
	chDone        chan struct{}
}

const (
	// SystemContextAddress is the address of the "Precompiled contract that calls that holds the current gas per pubdata byte"
	// https://sepolia.explorer.zksync.io/address/0x000000000000000000000000000000000000800b#contract
	SystemContextAddress = "0x000000000000000000000000000000000000800B"

	// ZksyncGasInfo_GetL2GasPerPubDataBytes is the a hex encoded call to:
	// function gasPerPubdataByte() external view returns (uint256 gasPerPubdataByte);
	SystemContext_gasPerPubdataByteMethod = "gasPerPubdataByte"
	ZksyncGasInfo_getGasPerPubdataByteL2  = "0x7cb9357e"

	// ZksyncGasInfo_GetL2GasPrice is the a hex encoded call to:
	// `function gasPrice() external view returns (uint256);`
	SystemContext_gasPriceMethod = "gasPrice"
	ZksyncGasInfo_getGasPriceL2  = "0xfe173b97"
)

func NewZkSyncL1GasOracle(lggr logger.Logger, ethClient l1OracleClient) *zkSyncL1Oracle {
	return &zkSyncL1Oracle{
		client:     ethClient,
		pollPeriod: PollPeriod,
		logger:     logger.Sugared(logger.Named(lggr, "L1GasOracle(zkSync)")),
		chainType:  chaintype.ChainZkSync,

		systemContextAddress:  SystemContextAddress,
		gasPerPubdataMethod:   SystemContext_gasPerPubdataByteMethod,
		gasPerPubdataSelector: ZksyncGasInfo_getGasPerPubdataByteL2,
		l2GasPriceMethod:      SystemContext_gasPriceMethod,
		l2GasPriceSelector:    ZksyncGasInfo_getGasPriceL2,

		chInitialised: make(chan struct{}),
		chStop:        make(chan struct{}),
		chDone:        make(chan struct{}),
	}
}

func (o *zkSyncL1Oracle) Name() string {
	return o.logger.Name()
}

func (o *zkSyncL1Oracle) ChainType(_ context.Context) chaintype.ChainType {
	return o.chainType
}

func (o *zkSyncL1Oracle) Start(ctx context.Context) error {
	return o.StartOnce(o.Name(), func() error {
		go o.run()
		<-o.chInitialised
		return nil
	})
}
func (o *zkSyncL1Oracle) Close() error {
	return o.StopOnce(o.Name(), func() error {
		close(o.chStop)
		<-o.chDone
		return nil
	})
}

func (o *zkSyncL1Oracle) HealthReport() map[string]error {
	return map[string]error{o.Name(): o.Healthy()}
}

func (o *zkSyncL1Oracle) run() {
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
func (o *zkSyncL1Oracle) refresh() (t *time.Timer) {
	t, err := o.refreshWithError()
	if err != nil {
		o.logger.Criticalw("Failed to refresh gas price", "err", err)
		o.SvcErrBuffer.Append(err)
	}
	return
}

func (o *zkSyncL1Oracle) refreshWithError() (t *time.Timer, err error) {
	t = time.NewTimer(utils.WithJitter(o.pollPeriod))

	ctx, cancel := o.chStop.CtxCancel(evmclient.ContextWithDefaultTimeout())
	defer cancel()

	price, err := o.CalculateL1GasPrice(ctx)
	if err != nil {
		return t, err
	}

	o.l1GasPriceMu.Lock()
	defer o.l1GasPriceMu.Unlock()
	o.l1GasPrice = priceEntry{price: assets.NewWei(price), timestamp: time.Now()}
	return
}

// For zkSync l2_gas_PerPubdataByte = (blob_byte_price_on_l1 + part_of_l1_verification_cost) / (gas_price_on_l2)
// l2_gas_PerPubdataByte = blob_gas_price_on_l1 * gas_per_byte / gas_price_on_l2
// blob_gas_price_on_l1 * gas_per_byte ~= gas_price_on_l2 * l2_gas_PerPubdataByte
func (o *zkSyncL1Oracle) CalculateL1GasPrice(ctx context.Context) (price *big.Int, err error) {
	l2GasPrice, err := o.GetL2GasPrice(ctx)
	if err != nil {
		return nil, err
	}
	l2GasPerPubDataByte, err := o.GetL2GasPerPubDataBytes(ctx)
	if err != nil {
		return nil, err
	}
	price = new(big.Int).Mul(l2GasPrice, l2GasPerPubDataByte)
	return
}

func (o *zkSyncL1Oracle) GasPrice(_ context.Context) (l1GasPrice *assets.Wei, err error) {
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
func (o *zkSyncL1Oracle) GetGasCost(ctx context.Context, tx *gethtypes.Transaction, blockNum *big.Int) (*assets.Wei, error) {
	//Unused method, so not implemented
	// And its not possible to know gas consumption of a transaction before its executed, since zkSync only posts the state difference
	return nil, fmt.Errorf("unimplemented")
}

// GetL2GasPrice calls SystemContract.gasPrice()  on the zksync system precompile contract.
//
// @return (The current gasPrice on L2: same as tx.gasPrice)
//  function gasPrice() external view returns (uint256);
//
// https://github.com/matter-labs/era-contracts/blob/12a7d3bc1777ae5663e7525b2628061502755cbd/system-contracts/contracts/interfaces/ISystemContext.sol#L34C4-L34C57

func (o *zkSyncL1Oracle) GetL2GasPrice(ctx context.Context) (gasPriceL2 *big.Int, err error) {
	precompile := common.HexToAddress(o.systemContextAddress)
	method, err := hex.DecodeString(ZksyncGasInfo_getGasPriceL2[2:])
	if err != nil {
		return common.Big0, fmt.Errorf("cannot decode method: %w", err)
	}
	b, err := o.client.CallContract(ctx, ethereum.CallMsg{
		To:   &precompile,
		Data: method,
	}, nil)
	if err != nil {
		return common.Big0, fmt.Errorf("cannot fetch l2GasPrice from zkSync SystemContract: %w", err)
	}

	if len(b) != 1*32 { // uint256 gasPrice;
		err = fmt.Errorf("return gasPrice (%d) different than expected (%d)", len(b), 3*32)
		return
	}
	gasPriceL2 = new(big.Int).SetBytes(b)
	return
}

// GetL2GasPerPubDataBytes calls SystemContract.gasPerPubdataByte() on the zksync system precompile contract.
//
// @return (The current gas per pubdata byte on L2)
// function gasPerPubdataByte() external view returns (uint256 gasPerPubdataByte);
//
// https://github.com/matter-labs/era-contracts/blob/12a7d3bc1777ae5663e7525b2628061502755cbd/system-contracts/contracts/interfaces/ISystemContext.sol#L58C14-L58C31

func (o *zkSyncL1Oracle) GetL2GasPerPubDataBytes(ctx context.Context) (gasPerPubByteL2 *big.Int, err error) {
	precompile := common.HexToAddress(o.systemContextAddress)
	method, err := hex.DecodeString(ZksyncGasInfo_getGasPerPubdataByteL2[2:])
	if err != nil {
		return common.Big0, fmt.Errorf("cannot decode method: %w", err)
	}
	b, err := o.client.CallContract(ctx, ethereum.CallMsg{
		To:   &precompile,
		Data: method,
	}, nil)
	if err != nil {
		return common.Big0, fmt.Errorf("cannot fetch gasPerPubdataByte from zkSync SystemContract: %w", err)
	}

	if len(b) != 1*32 { // uint256 gasPerPubdataByte;
		err = fmt.Errorf("return data length (%d) different than expected (%d)", len(b), 3*32)
		return
	}
	gasPerPubByteL2 = new(big.Int).SetBytes(b)
	return
}
