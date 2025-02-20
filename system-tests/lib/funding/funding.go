package funding

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/conversions"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

func PrivateKeyToAddress(privateKey *ecdsa.PrivateKey) (common.Address, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, errors.New("error casting public key to ECDSA")
	}
	return crypto.PubkeyToAddress(*publicKeyECDSA), nil
}

func SendFunds(logger zerolog.Logger, client *seth.Client, payload libtypes.FundsToSend) (*types.Receipt, error) {
	fromAddress, err := PrivateKeyToAddress(payload.PrivateKey)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), client.Cfg.Network.TxnTimeout.Duration())
	nonce, err := client.Client.PendingNonceAt(ctx, fromAddress)
	defer cancel()
	if err != nil {
		return nil, err
	}

	gasLimit, err := client.EstimateGasLimitForFundTransfer(fromAddress, payload.ToAddress, payload.Amount)
	if err != nil {
		transferGasFee := client.Cfg.Network.TransferGasFee
		if transferGasFee < 0 {
			return nil, fmt.Errorf("negative transfer gas fee: %d", transferGasFee)
		}
		gasLimit = uint64(transferGasFee)
	}

	gasPrice := big.NewInt(0)
	gasFeeCap := big.NewInt(0)
	gasTipCap := big.NewInt(0)

	if payload.GasLimit != nil {
		if *payload.GasLimit < 0 {
			return nil, fmt.Errorf("negative gas limit: %d", *payload.GasLimit)
		}
		gasLimit = uint64(*payload.GasLimit)
	}

	if client.Cfg.Network.EIP1559DynamicFees {
		// if any of the dynamic fees are not set, we need to either estimate them or read them from config
		if payload.GasFeeCap == nil || payload.GasTipCap == nil {
			// estimation or config reading happens here
			txOptions := client.NewTXOpts(seth.WithGasLimit(gasLimit))
			gasFeeCap = txOptions.GasFeeCap
			gasTipCap = txOptions.GasTipCap
		}

		// override with payload values if they are set
		if payload.GasFeeCap != nil {
			gasFeeCap = payload.GasFeeCap
		}

		if payload.GasTipCap != nil {
			gasTipCap = payload.GasTipCap
		}
	} else {
		if payload.GasPrice == nil {
			txOptions := client.NewTXOpts(seth.WithGasLimit(gasLimit))
			gasPrice = txOptions.GasPrice
		} else {
			gasPrice = payload.GasPrice
		}
	}

	var rawTx types.TxData

	if client.Cfg.Network.EIP1559DynamicFees {
		rawTx = &types.DynamicFeeTx{
			Nonce:     nonce,
			To:        &payload.ToAddress,
			Value:     payload.Amount,
			Gas:       gasLimit,
			GasFeeCap: gasFeeCap,
			GasTipCap: gasTipCap,
		}
	} else {
		rawTx = &types.LegacyTx{
			Nonce:    nonce,
			To:       &payload.ToAddress,
			Value:    payload.Amount,
			Gas:      gasLimit,
			GasPrice: gasPrice,
		}
	}

	signedTx, err := types.SignNewTx(payload.PrivateKey, types.LatestSignerForChainID(big.NewInt(client.ChainID)), rawTx)

	if err != nil {
		return nil, errors.Wrap(err, "failed to sign tx")
	}

	txTimeout := client.Cfg.Network.TxnTimeout.Duration()
	if payload.TxTimeout != nil {
		txTimeout = *payload.TxTimeout
	}

	logger.Debug().
		Str("From", fromAddress.Hex()).
		Str("To", payload.ToAddress.Hex()).
		Str("Amount (wei/ether)", fmt.Sprintf("%s/%s", payload.Amount, conversions.WeiToEther(payload.Amount).Text('f', -1))).
		Uint64("Nonce", nonce).
		Uint64("Gas Limit", gasLimit).
		Str("Gas Price", gasPrice.String()).
		Str("Gas Fee Cap", gasFeeCap.String()).
		Str("Gas Tip Cap", gasTipCap.String()).
		Bool("Dynamic fees", client.Cfg.Network.EIP1559DynamicFees).
		Msg("About to send funds")

	ctx, cancel = context.WithTimeout(ctx, txTimeout)
	defer cancel()
	err = client.Client.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send transaction")
	}

	logger.Debug().
		Str("From", fromAddress.Hex()).
		Str("To", payload.ToAddress.Hex()).
		Str("TxHash", signedTx.Hash().String()).
		Str("Amount (wei/ether)", fmt.Sprintf("%s/%s", payload.Amount, conversions.WeiToEther(payload.Amount).Text('f', -1))).
		Uint64("Nonce", nonce).
		Uint64("Gas Limit", gasLimit).
		Str("Gas Price", gasPrice.String()).
		Str("Gas Fee Cap", gasFeeCap.String()).
		Str("Gas Tip Cap", gasTipCap.String()).
		Bool("Dynamic fees", client.Cfg.Network.EIP1559DynamicFees).
		Msg("Sent funds")

	receipt, receiptErr := client.WaitMined(ctx, logger, client.Client, signedTx)
	if receiptErr != nil {
		return nil, errors.Wrap(receiptErr, "failed to wait for transaction to be mined")
	}

	if receipt.Status == 1 {
		return receipt, nil
	}

	tx, _, err := client.Client.TransactionByHash(ctx, signedTx.Hash())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get transaction by hash ")
	}

	_, err = client.Decode(tx, receiptErr)
	if err != nil {
		return nil, err
	}

	return receipt, nil
}
