package mantle

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"

	evmClient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
)

const (
	// tokenRatio is not volatile and can be requested not often.
	tokenRatioUpdateInterval = 60 * time.Minute
	// tokenRatio fetches the tokenRatio used for Mantle's gas price calculation
	tokenRatioMethod          = "tokenRatio"
	mantleTokenRatioAbiString = `[{"inputs":[],"name":"tokenRatio","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`
)

type Interceptor struct {
	client               evmClient.Client
	tokenRatioCallData   []byte
	tokenRatio           *big.Int
	tokenRatioLastUpdate time.Time
}

func NewInterceptor(_ context.Context, client evmClient.Client) (*Interceptor, error) {
	// Encode calldata for tokenRatio method
	tokenRatioMethodAbi, err := abi.JSON(strings.NewReader(mantleTokenRatioAbiString))
	if err != nil {
		return nil, fmt.Errorf("failed to parse GasPriceOracle %s() method ABI for Mantle; %w", tokenRatioMethod, err)
	}
	tokenRatioCallData, err := tokenRatioMethodAbi.Pack(tokenRatioMethod)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GasPriceOracle %s() calldata for Mantle; %w", tokenRatioMethod, err)
	}

	return &Interceptor{
		client:             client,
		tokenRatioCallData: tokenRatioCallData,
	}, nil
}

// ModifyGasPriceComponents returns modified gasPrice.
func (i *Interceptor) ModifyGasPriceComponents(ctx context.Context, execGasPrice, daGasPrice *big.Int) (*big.Int, *big.Int, error) {
	if time.Since(i.tokenRatioLastUpdate) > tokenRatioUpdateInterval {
		mantleTokenRatio, err := i.getMantleTokenRatio(ctx)
		if err != nil {
			return nil, nil, err
		}

		i.tokenRatio, i.tokenRatioLastUpdate = mantleTokenRatio, time.Now()
	}

	// multiply daGasPrice and execGas price by tokenRatio
	newExecGasPrice := new(big.Int).Mul(execGasPrice, i.tokenRatio)
	newDAGasPrice := new(big.Int).Mul(daGasPrice, i.tokenRatio)
	return newExecGasPrice, newDAGasPrice, nil
}

// getMantleTokenRatio Requests and returns a token ratio value for the Mantle chain.
func (i *Interceptor) getMantleTokenRatio(ctx context.Context) (*big.Int, error) {
	// FIXME it's removed from chainlink repo
	// precompile := common.HexToAddress(rollups.OPGasOracleAddress)
	precompile := utils.RandomAddress()
	tokenRatio, err := i.client.CallContract(ctx, ethereum.CallMsg{
		To:   &precompile,
		Data: i.tokenRatioCallData,
	}, nil)

	if err != nil {
		return nil, fmt.Errorf("getMantleTokenRatio call failed: %w", err)
	}

	return new(big.Int).SetBytes(tokenRatio), nil
}
