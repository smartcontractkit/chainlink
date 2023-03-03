package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	ethcontracts "github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

const optimismGasOracleAddress string = "0x420000000000000000000000000000000000000F"

// OptimismMultinodeClient represents a multi-node, EVM compatible client for the Optimism network
type OptimismMultinodeClient struct {
	*EthereumMultinodeClient
}

// OptimismClient represents a single node, EVM compatible client for the Optimism network
type OptimismClient struct {
	*EthereumClient
}

// Fund sends some ETH to an address using the default wallet
func (o *OptimismClient) Fund(
	toAddress string,
	amount *big.Float,
) error {
	privateKey, err := crypto.HexToECDSA(o.DefaultWallet.PrivateKey())
	if err != nil {
		return fmt.Errorf("invalid private key: %v", err)
	}
	to := common.HexToAddress(toAddress)

	// Optimism is unique in its usage of an L1 data fee on top of regular gas costs. Need to call their oracle
	// https://community.optimism.io/docs/developers/build/transaction-fees/#the-l1-data-fee
	gasOracle, err := ethcontracts.NewOptimismGas(common.HexToAddress(optimismGasOracleAddress), o.Client)
	if err != nil {
		return err
	}
	opts := &bind.CallOpts{
		From:    common.HexToAddress(o.GetDefaultWallet().Address()),
		Context: context.Background(),
	}
	l1Fee, err := gasOracle.GetL1Fee(opts, types.DynamicFeeTx{To: &to}.Data)
	if err != nil {
		return err
	}

	suggestedGasTipCap, err := o.Client.SuggestGasTipCap(context.Background())
	if err != nil {
		return err
	}

	// Bump Tip Cap
	gasPriceBuffer := big.NewInt(0).SetUint64(o.NetworkConfig.GasEstimationBuffer)
	suggestedGasTipCap.Add(suggestedGasTipCap, gasPriceBuffer)

	nonce, err := o.GetNonce(context.Background(), common.HexToAddress(o.DefaultWallet.Address()))
	if err != nil {
		return err
	}
	latestHeader, err := o.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}
	baseFeeMult := new(big.Int).Mul(latestHeader.BaseFee, big.NewInt(2))
	gasFeeCap := new(big.Int).Add(baseFeeMult, suggestedGasTipCap)

	estimatedGas, err := o.Client.EstimateGas(context.Background(), ethereum.CallMsg{})
	if err != nil {
		return err
	}
	totalEstimatedGasCost := new(big.Int).Mul(gasFeeCap, new(big.Int).SetUint64(estimatedGas))
	totalEstimatedGasCost.Add(totalEstimatedGasCost, l1Fee)

	unsignedTx := &types.DynamicFeeTx{
		ChainID:   o.GetChainID(),
		Nonce:     nonce,
		To:        &to,
		Value:     utils.EtherToWei(amount),
		GasTipCap: suggestedGasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       estimatedGas,
	}

	tx, err := types.SignNewTx(privateKey, types.LatestSignerForChainID(o.GetChainID()), unsignedTx)
	if err != nil {
		return err
	}

	log.Info().
		Str("Token", "OP").
		Str("From", o.DefaultWallet.Address()).
		Str("To", toAddress).
		Str("Amount", amount.String()).
		Uint64("Estimated Gas Cost", totalEstimatedGasCost.Uint64()).
		Msg("Funding Address")
	if err := o.SendTransaction(context.Background(), tx); err != nil {
		if strings.Contains(err.Error(), "nonce") {
			err = errors.Wrap(err, fmt.Sprintf("using nonce %d", nonce))
		}
		return err
	}

	return o.ProcessTransaction(tx)
}

func (o *OptimismClient) ReturnFunds(fromKey *ecdsa.PrivateKey) error {
	var tx *types.Transaction
	var err error
	for attempt := 0; attempt < 10; attempt++ {
		tx, err = o.attemptReturn(fromKey, attempt)
		if err == nil {
			return o.ProcessTransaction(tx)
		}
		log.Debug().Err(err).Int("Attempt", attempt+1).Msg("Error returning funds from Chainlink node, trying again")
	}
	return err
}

// a single fund return attempt, further attempts exponentially raise the error margin for fund returns
func (o *OptimismClient) attemptReturn(fromKey *ecdsa.PrivateKey, attemptCount int) (*types.Transaction, error) {
	to := common.HexToAddress(o.DefaultWallet.Address())

	// Optimism is unique in its usage of an L1 data fee on top of regular gas costs. Need to call their oracle
	// https://community.optimism.io/docs/developers/build/transaction-fees/#the-l1-data-fee
	gasOracle, err := ethcontracts.NewOptimismGas(common.HexToAddress(optimismGasOracleAddress), o.Client)
	if err != nil {
		return nil, err
	}
	opts := &bind.CallOpts{
		From:    common.HexToAddress(o.GetDefaultWallet().Address()),
		Context: context.Background(),
	}
	l1Fee, err := gasOracle.GetL1Fee(opts, types.DynamicFeeTx{To: &to}.Data)
	if err != nil {
		return nil, err
	}

	suggestedGasTipCap, err := o.Client.SuggestGasTipCap(context.Background())
	if err != nil {
		return nil, err
	}
	// Bump Tip Cap
	gasPriceBuffer := big.NewInt(0).SetUint64(o.NetworkConfig.GasEstimationBuffer)
	suggestedGasTipCap.Add(suggestedGasTipCap, gasPriceBuffer)
	latestHeader, err := o.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	baseFeeMult := new(big.Int).Mul(latestHeader.BaseFee, big.NewInt(2))
	gasFeeCap := new(big.Int).Add(baseFeeMult, suggestedGasTipCap)

	fromAddress, err := utils.PrivateKeyToAddress(fromKey)
	if err != nil {
		return nil, err
	}
	originalBalance, err := o.Client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		return nil, err
	}

	nonce, err := o.GetNonce(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}
	estimatedGas, err := o.Client.EstimateGas(context.Background(), ethereum.CallMsg{})
	if err != nil {
		return nil, err
	}
	totalEstimatedGasCost := new(big.Int).Mul(gasFeeCap, new(big.Int).SetUint64(estimatedGas))
	buffer := big.NewInt(int64(math.Pow(float64(attemptCount), 2) * 1000000000)) // exponentially increase error margin)
	totalEstimatedGasCost.Add(totalEstimatedGasCost, l1Fee)
	totalEstimatedGasCost.Add(totalEstimatedGasCost, buffer)

	sendBalance := new(big.Int).Sub(originalBalance, totalEstimatedGasCost)
	sendBalance.Sub(sendBalance, totalEstimatedGasCost)

	unsignedTx := &types.DynamicFeeTx{
		ChainID:   o.GetChainID(),
		Nonce:     nonce,
		To:        &to,
		Value:     sendBalance,
		GasTipCap: suggestedGasTipCap, // eth_maxPriorityFeePerGas
		GasFeeCap: gasFeeCap,          // maxFeePerGas = eth_maxPriorityFeePerGas + (BASEFEE * 2) + gasOracle.GetL1Fee(txData) + buffer
		Gas:       estimatedGas,       // eth_estimateGas (51,000 for a normal Tx)
	}

	signedTx, err := types.SignNewTx(fromKey, types.LatestSignerForChainID(o.GetChainID()), unsignedTx)
	if err != nil {
		return nil, err
	}
	log.Info().
		Str("Token", "OP").
		Uint64("Original Balance", originalBalance.Uint64()).
		Uint64("Send Amount", signedTx.Value().Uint64()).
		Str("From", fromAddress.Hex()).
		Uint64("L1 Fee", l1Fee.Uint64()).
		Uint64("Buffer", buffer.Uint64()).
		Uint64("Estimated Gas Cost", totalEstimatedGasCost.Uint64()).
		Msg("Returning Funds to Default Wallet")
	return signedTx, o.SendTransaction(context.Background(), signedTx)
}
