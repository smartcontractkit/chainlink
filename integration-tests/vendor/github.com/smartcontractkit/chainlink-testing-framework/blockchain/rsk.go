package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

// Handles specific issues with the RSK EVM chain: https://developers.rsk.co/rsk/node/architecture/json-rpc/

// RSKMultinodeClient represents a multi-node, EVM compatible client for the RSK network
type RSKMultinodeClient struct {
	*EthereumMultinodeClient
}

// RSKClient represents a single node, EVM compatible client for the RSK network
type RSKClient struct {
	*EthereumClient
}

// Fund sends some ETH to an address using the default wallet
func (r *RSKClient) Fund(toAddress string, amount *big.Float) error {
	privateKey, err := crypto.HexToECDSA(r.DefaultWallet.PrivateKey())
	to := common.HexToAddress(toAddress)
	if err != nil {
		return fmt.Errorf("invalid private key: %v", err)
	}
	// RSK uses legacy transactions and gas estimations, is behind London fork as of 04/27/2022
	suggestedGasPrice, err := r.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	// Bump gas price
	gasPriceBuffer := big.NewInt(0).SetUint64(r.NetworkConfig.GasEstimationBuffer)
	suggestedGasPrice.Add(suggestedGasPrice, gasPriceBuffer)

	nonce, err := r.GetNonce(context.Background(), common.HexToAddress(r.DefaultWallet.Address()))
	if err != nil {
		return err
	}
	estimatedGas, err := r.Client.EstimateGas(context.Background(), ethereum.CallMsg{})
	if err != nil {
		return err
	}

	tx, err := types.SignNewTx(privateKey, types.LatestSignerForChainID(r.GetChainID()), &types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    utils.EtherToWei(amount),
		GasPrice: suggestedGasPrice,
		Gas:      estimatedGas,
	})
	if err != nil {
		return err
	}

	log.Info().
		Str("Token", "RSK").
		Str("From", r.DefaultWallet.Address()).
		Str("To", toAddress).
		Str("Amount", amount.String()).
		Uint64("Estimated Gas Cost", new(big.Int).Mul(suggestedGasPrice, new(big.Int).SetUint64(estimatedGas)).Uint64()).
		Msg("Funding Address")
	if err := r.SendTransaction(context.Background(), tx); err != nil {
		return err
	}

	return r.ProcessTransaction(tx)
}

// DeployContract acts as a general contract deployment tool to an EVM chain
func (r *RSKClient) DeployContract(
	contractName string,
	deployer ContractDeployer,
) (*common.Address, *types.Transaction, interface{}, error) {
	opts, err := r.TransactionOpts(r.DefaultWallet)
	if err != nil {
		return nil, nil, nil, err
	}

	// RSK uses legacy transactions and gas estimations, is behind London fork as of 04/27/2022
	suggestedGasPrice, err := r.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, nil, nil, err
	}

	// Bump gas price
	gasPriceBuffer := big.NewInt(0).SetUint64(r.NetworkConfig.GasEstimationBuffer)
	suggestedGasPrice.Add(suggestedGasPrice, gasPriceBuffer)

	opts.GasPrice = suggestedGasPrice

	contractAddress, transaction, contractInstance, err := deployer(opts, r.Client)
	if err != nil {
		return nil, nil, nil, err
	}

	if err = r.ProcessTransaction(transaction); err != nil {
		return nil, nil, nil, err
	}

	log.Info().
		Str("Contract Address", contractAddress.Hex()).
		Str("Contract Name", contractName).
		Str("From", r.DefaultWallet.Address()).
		Str("Total Gas Cost (RSK)", utils.WeiToEther(transaction.Cost()).String()).
		Str("Network Name", r.NetworkConfig.Name).
		Msg("Deployed contract")
	return &contractAddress, transaction, contractInstance, err
}

// Fund sends some ARB to an address using the default wallet
func (r *RSKClient) ReturnFunds(fromPrivateKey *ecdsa.PrivateKey) error {
	to := common.HexToAddress(r.DefaultWallet.Address())

	// Arbitrum uses legacy transactions and gas estimations
	suggestedGasPrice, err := r.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	fromAddress, err := utils.PrivateKeyToAddress(fromPrivateKey)
	if err != nil {
		return err
	}

	balance, err := r.Client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		return err
	}
	balance.Sub(balance, big.NewInt(1).Mul(suggestedGasPrice, big.NewInt(21000)))

	nonce, err := r.GetNonce(context.Background(), fromAddress)
	if err != nil {
		return err
	}
	estimatedGas, err := r.Client.EstimateGas(context.Background(), ethereum.CallMsg{})
	if err != nil {
		return err
	}

	tx, err := types.SignNewTx(fromPrivateKey, types.LatestSignerForChainID(r.GetChainID()), &types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    balance,
		GasPrice: suggestedGasPrice,
		Gas:      estimatedGas,
	})
	if err != nil {
		return err
	}

	log.Info().
		Str("Token", "RSK").
		Str("From", fromAddress.Hex()).
		Str("Amount", balance.String()).
		Msg("Returning Funds to Default Wallet")
	if err := r.SendTransaction(context.Background(), tx); err != nil {
		return err
	}

	return r.ProcessTransaction(tx)
}
