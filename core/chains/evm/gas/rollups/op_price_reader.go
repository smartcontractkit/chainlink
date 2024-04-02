package rollups

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/common/config"
)

const (
	// OPStackGasOracle_l1BaseFee fetches the l1 base fee set in the OP Stack GasPriceOracle contract
	OPStackGasOracle_l1BaseFee = "l1BaseFee"

	// OPStackGasOracle_isEcotone fetches if the OP Stack GasPriceOracle contract has upgraded to Ecotone
	OPStackGasOracle_isEcotone = "isEcotone"

	// OPStackGasOracle_getL1GasUsed fetches the l1 gas used for given tx bytes
	OPStackGasOracle_getL1GasUsed = "getL1GasUsed"

	// OPStackGasOracle_getL1Fee fetches the l1 fee for given tx bytes
	OPStackGasOracle_getL1Fee = "getL1Fee"

	// OPStackGasOracle_isEcotonePollingPeriod is the interval to poll if chain has upgraded to Ecotone
	// Set to poll every 4 hours
	OPStackGasOracle_isEcotonePollingPeriod = 14400
)

type opStackGasPriceReader struct {
	client ethClient
	logger logger.SugaredLogger

	oracleAddress      common.Address
	isEcotoneMethodAbi abi.ABI

	l1BaseFeeCalldata    []byte
	isEcotoneCalldata    []byte
	getL1GasUsedCalldata []byte
	getL1FeeCalldata     []byte

	isEcotone        bool
	isEcotoneCheckTs int64
}

func newOPPriceReader(lggr logger.Logger, ethClient ethClient, chainType config.ChainType, oracleAddress string) daPriceReader {
	// encode calldata for each method; these calldata will remain the same for each call, we can encode them just once
	l1BaseFeeMethodAbi, err := abi.JSON(strings.NewReader(L1BaseFeeAbiString))
	if err != nil {
		panic(fmt.Errorf("failed to parse GasPriceOracle %s() method ABI for chain: %s; %w", OPStackGasOracle_l1BaseFee, chainType, err))
	}
	l1BaseFeeCalldata, err := l1BaseFeeMethodAbi.Pack(OPStackGasOracle_l1BaseFee)
	if err != nil {
		panic(fmt.Errorf("failed to parse GasPriceOracle %s() calldata for chain: %s; %w", OPStackGasOracle_l1BaseFee, chainType, err))
	}

	isEcotoneMethodAbi, err := abi.JSON(strings.NewReader(OPIsEcotoneAbiString))
	if err != nil {
		panic(fmt.Errorf("failed to parse GasPriceOracle %s() method ABI for chain: %s; %w", OPStackGasOracle_isEcotone, chainType, err))
	}
	isEcotoneCalldata, err := isEcotoneMethodAbi.Pack(OPStackGasOracle_isEcotone)
	if err != nil {
		panic(fmt.Errorf("failed to parse GasPriceOracle %s() calldata for chain: %s; %w", OPStackGasOracle_isEcotone, chainType, err))
	}

	getL1GasUsedMethodAbi, err := abi.JSON(strings.NewReader(OPGetL1GasUsedAbiString))
	if err != nil {
		panic(fmt.Errorf("failed to parse GasPriceOracle %s() method ABI for chain: %s; %w", OPStackGasOracle_getL1GasUsed, chainType, err))
	}
	getL1GasUsedCalldata, err := getL1GasUsedMethodAbi.Pack(OPStackGasOracle_getL1GasUsed, []byte{0x1})
	if err != nil {
		panic(fmt.Errorf("failed to parse GasPriceOracle %s() calldata for chain: %s; %w", OPStackGasOracle_getL1GasUsed, chainType, err))
	}

	getL1FeeMethodAbi, err := abi.JSON(strings.NewReader(GetL1FeeAbiString))
	if err != nil {
		panic(fmt.Errorf("failed to parse GasPriceOracle %s() method ABI for chain: %s; %w", OPStackGasOracle_getL1Fee, chainType, err))
	}
	getL1FeeCalldata, err := getL1FeeMethodAbi.Pack(OPStackGasOracle_getL1Fee, []byte{0x1})
	if err != nil {
		panic(fmt.Errorf("failed to parse GasPriceOracle %s() calldata for chain: %s; %w", OPStackGasOracle_getL1Fee, chainType, err))
	}

	return &opStackGasPriceReader{
		client: ethClient,
		logger: logger.Sugared(logger.Named(lggr, fmt.Sprintf("OPStackGasOracle(%s)", chainType))),

		oracleAddress:      common.HexToAddress(oracleAddress),
		isEcotoneMethodAbi: isEcotoneMethodAbi,

		l1BaseFeeCalldata:    l1BaseFeeCalldata,
		isEcotoneCalldata:    isEcotoneCalldata,
		getL1GasUsedCalldata: getL1GasUsedCalldata,
		getL1FeeCalldata:     getL1FeeCalldata,

		isEcotone:        false,
		isEcotoneCheckTs: 0,
	}
}

func (o *opStackGasPriceReader) GetDAGasPrice(ctx context.Context) (*big.Int, error) {
	isEcotone, err := o.checkIsEcotone(ctx)
	if err != nil {
		return nil, err
	}

	o.logger.Infof("Chain isEcotone result: %t", isEcotone)

	if isEcotone {
		return o.getEcotoneGasPrice(ctx)
	}

	return o.getV1GasPrice(ctx)
}

func (o *opStackGasPriceReader) checkIsEcotone(ctx context.Context) (bool, error) {
	// if chain is already Ecotone, NOOP
	if o.isEcotone {
		return true, nil
	}
	// if time since last check has not exceeded polling period, NOOP
	if time.Now().Unix()-o.isEcotoneCheckTs < OPStackGasOracle_isEcotonePollingPeriod {
		return false, nil
	}
	o.isEcotoneCheckTs = time.Now().Unix()

	// confirmed with OP team that isEcotone() is the canonical way to check if the chain has upgraded
	b, err := o.client.CallContract(ctx, ethereum.CallMsg{
		To:   &o.oracleAddress,
		Data: o.isEcotoneCalldata,
	}, nil)

	// if the chain has not upgraded to Ecotone, the isEcotone call will revert, this would be expected
	if err != nil {
		o.logger.Infof("isEcotone() call failed, this can happen if chain has not upgraded: %w", err)
		return false, nil
	}

	res, err := o.isEcotoneMethodAbi.Unpack(OPStackGasOracle_isEcotone, b)
	if err != nil {
		return false, fmt.Errorf("failed to unpack isEcotone() return data: %w", err)
	}
	o.isEcotone = res[0].(bool)
	return o.isEcotone, nil
}

func (o *opStackGasPriceReader) getV1GasPrice(ctx context.Context) (*big.Int, error) {
	b, err := o.client.CallContract(ctx, ethereum.CallMsg{
		To:   &o.oracleAddress,
		Data: o.l1BaseFeeCalldata,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("l1BaseFee() call failed: %w", err)
	}

	if len(b) != 32 {
		return nil, fmt.Errorf("l1BaseFee() return data length (%d) different than expected (%d)", len(b), 32)
	}
	return new(big.Int).SetBytes(b), nil
}

func (o *opStackGasPriceReader) getEcotoneGasPrice(ctx context.Context) (*big.Int, error) {
	rpcBatchCalls := []rpc.BatchElem{
		{
			Method: "eth_call",
			Args: []any{
				map[string]interface{}{
					"from": common.Address{},
					"to":   o.oracleAddress,
					"data": hexutil.Bytes(o.getL1GasUsedCalldata),
				},
				"latest",
			},
			Result: new(string),
		},
		{
			Method: "eth_call",
			Args: []any{
				map[string]interface{}{
					"from": common.Address{},
					"to":   o.oracleAddress,
					"data": hexutil.Bytes(o.getL1FeeCalldata),
				},
				"latest",
			},
			Result: new(string),
		},
	}

	err := o.client.BatchCallContext(ctx, rpcBatchCalls)
	if err != nil {
		return nil, fmt.Errorf("getEcotoneGasPrice batch call failed: %w", err)
	}
	if rpcBatchCalls[0].Error != nil {
		return nil, fmt.Errorf("%s call failed in a batch: %w", OPStackGasOracle_getL1GasUsed, err)
	}
	if rpcBatchCalls[1].Error != nil {
		return nil, fmt.Errorf("%s call failed in a batch: %w", OPStackGasOracle_getL1Fee, err)
	}

	l1GasUsedResult := *(rpcBatchCalls[0].Result.(*string))
	l1FeeResult := *(rpcBatchCalls[1].Result.(*string))

	l1GasUsedBytes, err := hexutil.Decode(l1GasUsedResult)
	if err != nil {
		return nil, fmt.Errorf("failed to decode %s rpc result: %w", OPStackGasOracle_getL1GasUsed, err)
	}
	l1FeeBytes, err := hexutil.Decode(l1FeeResult)
	if err != nil {
		return nil, fmt.Errorf("failed to decode %s rpc result: %w", OPStackGasOracle_getL1Fee, err)
	}

	l1GasUsed := new(big.Int).SetBytes(l1GasUsedBytes)
	l1Fee := new(big.Int).SetBytes(l1FeeBytes)

	// for the same tx byte, l1Fee / l1GasUsed will give the l1 gas price
	// note this price is per l1 gas, not l1 data byte
	return new(big.Int).Div(l1Fee, l1GasUsed), nil
}
