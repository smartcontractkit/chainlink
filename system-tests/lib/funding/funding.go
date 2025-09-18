package funding

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/conversions"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	crecrypto "github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
)

type FundsToSend struct {
	ToAddress  common.Address
	Amount     *big.Int
	PrivateKey *ecdsa.PrivateKey
	GasLimit   *int64
	GasPrice   *big.Int
	GasFeeCap  *big.Int
	GasTipCap  *big.Int
	TxTimeout  *time.Duration
	Nonce      *uint64
}

type FundsToSendSol struct {
	Recipent   solana.PublicKey
	PrivateKey solana.PrivateKey
	Amount     uint64
}

func SendFundsSol(ctx context.Context, logger zerolog.Logger, client *rpc.Client, payload FundsToSendSol) error {
	funder := payload.PrivateKey
	recipient := payload.Recipent
	if recipient.IsZero() {
		return errors.New("recipient is zero")
	}
	bal, err := client.GetBalance(ctx, funder.PublicKey(), rpc.CommitmentConfirmed)
	if err != nil {
		return fmt.Errorf("failed to get funder balance: %w", err)
	}
	logger.Debug().
		Uint64("Sender balance:", bal.Value)

	recent, err := client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return fmt.Errorf("failed to get recent block hash: %w", err)
	}

	tx, err := solana.NewTransaction([]solana.Instruction{
		system.NewTransferInstruction(
			payload.Amount,
			funder.PublicKey(),
			recipient,
		).Build(),
	},
		recent.Value.Blockhash,
		solana.TransactionPayer(funder.PublicKey()))
	if err != nil {
		return fmt.Errorf("failed to build fund transaction: %w", err)
	}

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if funder.PublicKey().Equals(key) {
				return &funder
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("failed to sign fund transaction: %w", err)
	}

	_, err = client.SendTransaction(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to send fund transaction: %w", err)
	}

	bal2, err := client.GetBalance(ctx, funder.PublicKey(), rpc.CommitmentConfirmed)
	if err != nil {
		return fmt.Errorf("failed to get recipient balance: %w", err)
	}

	logger.Debug().
		Uint64("Recipient balance: ", bal2.Value)
	return nil
}

func SendFunds(ctx context.Context, logger zerolog.Logger, client *seth.Client, payload FundsToSend) (*types.Receipt, error) {
	fromAddress, err := crecrypto.PrivateKeyToAddress(payload.PrivateKey)
	if err != nil {
		return nil, err
	}

	var nonce uint64
	if payload.Nonce == nil {
		nonceCtx, cancel := context.WithTimeout(ctx, client.Cfg.Network.TxnTimeout.Duration())
		nonce, err = client.Client.PendingNonceAt(nonceCtx, fromAddress)
		defer cancel()
		if err != nil {
			return nil, err
		}
	} else {
		nonce = *payload.Nonce
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

	sendCtx, cancel := context.WithTimeout(ctx, txTimeout)
	defer cancel()
	err = client.Client.SendTransaction(sendCtx, signedTx)
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

	minedCtx, mineCancel := context.WithTimeout(ctx, txTimeout)
	defer mineCancel()
	receipt, receiptErr := client.WaitMined(minedCtx, logger, client.Client, signedTx)
	if receiptErr != nil {
		return nil, errors.Wrap(receiptErr, "failed to wait for transaction to be mined")
	}

	if receipt.Status == 1 {
		return receipt, nil
	}

	txCtx, txCancel := context.WithTimeout(ctx, txTimeout)
	defer txCancel()
	tx, _, err := client.Client.TransactionByHash(txCtx, signedTx.Hash())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get transaction by hash ")
	}

	_, err = client.Decode(tx, receiptErr)
	if err != nil {
		return nil, err
	}

	return receipt, nil
}
