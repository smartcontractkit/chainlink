package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

// ArbitrumMultinodeClient represents a multi-node, EVM compatible client for the Arbitrum network
type ArbitrumMultinodeClient struct {
	*EthereumMultinodeClient
}

// ArbitrumClient represents a single node, EVM compatible client for the Arbitrum network
type ArbitrumClient struct {
	*EthereumClient
}

// Fund sends some ARB to an address using the default wallet
func (a *ArbitrumClient) Fund(toAddress string, amount *big.Float) error {
	privateKey, err := crypto.HexToECDSA(a.DefaultWallet.PrivateKey())
	if err != nil {
		return fmt.Errorf("invalid private key: %v", err)
	}
	to := common.HexToAddress(toAddress)

	suggestedGasTipCap, err := a.Client.SuggestGasTipCap(context.Background())
	if err != nil {
		return err
	}

	// Bump Tip Cap
	gasPriceBuffer := big.NewInt(0).SetUint64(a.NetworkConfig.GasEstimationBuffer)
	suggestedGasTipCap.Add(suggestedGasTipCap, gasPriceBuffer)

	nonce, err := a.GetNonce(context.Background(), common.HexToAddress(a.DefaultWallet.Address()))
	if err != nil {
		return err
	}
	latestHeader, err := a.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}
	baseFeeMult := big.NewInt(1).Mul(latestHeader.BaseFee, big.NewInt(2))
	gasFeeCap := baseFeeMult.Add(baseFeeMult, suggestedGasTipCap)

	estimatedGas, err := a.Client.EstimateGas(context.Background(), ethereum.CallMsg{})
	if err != nil {
		return err
	}

	tx, err := types.SignNewTx(privateKey, types.LatestSignerForChainID(a.GetChainID()), &types.DynamicFeeTx{
		ChainID:   a.GetChainID(),
		Nonce:     nonce,
		To:        &to,
		Value:     utils.EtherToWei(amount),
		GasTipCap: suggestedGasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       estimatedGas,
	})
	if err != nil {
		return err
	}

	log.Info().
		Str("Token", "ARB").
		Str("From", a.DefaultWallet.Address()).
		Str("To", toAddress).
		Str("Hash", tx.Hash().Hex()).
		Uint64("Nonce", tx.Nonce()).
		Str("Network Name", a.GetNetworkName()).
		Str("Amount", amount.String()).
		Uint64("Estimated Gas Cost", new(big.Int).Mul(gasFeeCap, new(big.Int).SetUint64(estimatedGas)).Uint64()).
		Msg("Funding Address")
	if err := a.SendTransaction(context.Background(), tx); err != nil {
		if strings.Contains(err.Error(), "nonce") {
			err = errors.Wrap(err, fmt.Sprintf("using nonce %d", nonce))
		}
		return err
	}

	return a.ProcessTransaction(tx)
}

func (a *ArbitrumClient) ReturnFunds(fromKey *ecdsa.PrivateKey) error {
	var tx *types.Transaction
	var err error
	for attempt := 1; attempt < 10; attempt++ {
		tx, err = attemptArbReturn(a, fromKey, attempt)
		if err == nil {
			return a.ProcessTransaction(tx)
		}
		log.Debug().Err(err).Int("Attempt", attempt+1).Msg("Error returning funds from Chainlink node, trying again")
	}
	return err
}

// a single fund return attempt, further attempts exponentially raise the error margin for fund returns
func attemptArbReturn(a *ArbitrumClient, fromKey *ecdsa.PrivateKey, attemptCount int) (*types.Transaction, error) {
	to := common.HexToAddress(a.DefaultWallet.Address())

	suggestedGasTipCap, err := a.Client.SuggestGasTipCap(context.Background())
	if err != nil {
		return nil, err
	}
	latestHeader, err := a.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	baseFeeMult := big.NewInt(1).Mul(latestHeader.BaseFee, big.NewInt(2))
	gasFeeCap := baseFeeMult.Add(baseFeeMult, suggestedGasTipCap)
	gasFeeCap.Add(gasFeeCap, big.NewInt(int64(math.Pow(float64(attemptCount), 2)*1000))) // exponentially increase error margin

	fromAddress, err := utils.PrivateKeyToAddress(fromKey)
	if err != nil {
		return nil, err
	}
	balance, err := a.Client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		return nil, err
	}
	estGas, err := a.Client.EstimateGas(context.Background(), ethereum.CallMsg{})
	if err != nil {
		return nil, err
	}
	balance.Sub(balance, big.NewInt(1).Mul(gasFeeCap, big.NewInt(0).SetUint64(estGas)))

	nonce, err := a.GetNonce(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}
	estimatedGas, err := a.Client.EstimateGas(context.Background(), ethereum.CallMsg{})
	if err != nil {
		return nil, err
	}

	tx, err := types.SignNewTx(fromKey, types.LatestSignerForChainID(a.GetChainID()), &types.DynamicFeeTx{
		ChainID:   a.GetChainID(),
		Nonce:     nonce,
		To:        &to,
		Value:     balance,
		GasTipCap: suggestedGasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       estimatedGas,
	})
	if err != nil {
		return nil, err
	}
	log.Info().
		Str("Token", "ARB").
		Str("Amount", balance.String()).
		Str("From", fromAddress.Hex()).
		Msg("Returning Funds to Default Wallet")
	return tx, a.SendTransaction(context.Background(), tx)
}
