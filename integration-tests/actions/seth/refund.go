package actions_seth

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/seth"
)

// TransactionRetrier is an interface that every retrier of failed funds transfer transaction needs to implement
type TransactionRetrier interface {
	Retry(ctx context.Context, logger zerolog.Logger, client *seth.Client, txErr error, payload FundsToSendPayload, currentAttempt int) error
}

// InsufficientFundTransferRetrier will retry a failed funds transfer transaction if the error is due to insufficient funds
// by substracting 1 Gwei from amount to send and retrying it up to maxRetries times
type InsufficientFundTransferRetrier struct {
	nextRetrier TransactionRetrier
	maxRetries  int
}

func (r *InsufficientFundTransferRetrier) Retry(ctx context.Context, logger zerolog.Logger, client *seth.Client, txErr error, payload FundsToSendPayload, currentAttempt int) error {
	if currentAttempt >= r.maxRetries {
		if r.nextRetrier != nil {
			return r.nextRetrier.Retry(ctx, logger, client, txErr, payload, 0)
		}
		return txErr
	}

	for txErr != nil && (strings.Contains(txErr.Error(), "insufficient funds")) {
		payload.Amount = payload.Amount.Sub(payload.Amount, big.NewInt(blockchain.GWei))

		retryErr := SendFunds(logger, client, payload)
		if retryErr == nil {
			return nil
		}

		if strings.Contains(retryErr.Error(), "insufficient funds") {
			r.Retry(ctx, logger, client, retryErr, payload, currentAttempt+1)
		}
	}

	if r.nextRetrier != nil {
		return r.nextRetrier.Retry(ctx, logger, client, txErr, payload, 0)
	}

	return txErr
}

// GasTooLowTransferRetrier will retry a failed funds transfer transaction if the error is due to gas too low
// by doubling the gas limit and retrying until reaching maxGasLimit
type GasTooLowTransferRetrier struct {
	nextRetrier TransactionRetrier
	maxGasLimit uint64
}

func (r *GasTooLowTransferRetrier) Retry(ctx context.Context, logger zerolog.Logger, client *seth.Client, txErr error, payload FundsToSendPayload, currentAttempt int) error {
	if payload.GasLimit != nil && *payload.GasLimit >= r.maxGasLimit {
		if r.nextRetrier != nil {
			return r.nextRetrier.Retry(ctx, logger, client, txErr, payload, 0)
		}
		return txErr
	}

	for txErr != nil && strings.Contains(txErr.Error(), "gas too low") {
		var newGasLimit uint64
		if payload.GasLimit != nil {
			newGasLimit = *payload.GasLimit * 2
		} else {
			newGasLimit = uint64(client.Cfg.Network.TransferGasFee) * 2
		}

		payload.GasLimit = &newGasLimit

		retryErr := SendFunds(logger, client, payload)
		if retryErr == nil {
			return nil
		}

		if strings.Contains(retryErr.Error(), "insufficient funds") {
			r.Retry(ctx, logger, client, retryErr, payload, currentAttempt+1)
		}
	}

	if r.nextRetrier != nil {
		return r.nextRetrier.Retry(ctx, logger, client, txErr, payload, 0)
	}

	return txErr
}

// OvershotTransferRetrier will retry a failed funds transfer transaction if the error is due to overshot
// by substracting the overshot amount from the amount to send and retrying it up to maxRetries times
type OvershotTransferRetrier struct {
	nextRetrier TransactionRetrier
	maxRetries  int
}

func (r *OvershotTransferRetrier) Retry(ctx context.Context, logger zerolog.Logger, client *seth.Client, txErr error, payload FundsToSendPayload, currentAttempt int) error {
	if currentAttempt >= r.maxRetries {
		if r.nextRetrier != nil {
			return r.nextRetrier.Retry(ctx, logger, client, txErr, payload, 0)
		}
		return txErr
	}

	overshotRe := regexp.MustCompile(`overshot (\d+)`)
	if txErr != nil && strings.Contains(txErr.Error(), "overshot") {
		logger.Info().Msg("Overshot error detected, retrying with less funds")
		submatches := overshotRe.FindStringSubmatch(txErr.Error())
		if len(submatches) < 1 {
			return fmt.Errorf("error parsing overshot amount in error: %w", txErr)
		}
		numberString := submatches[1]
		overshotAmount, err := strconv.Atoi(numberString)
		if err != nil {
			return err
		}

		payload.Amount = payload.Amount.Sub(payload.Amount, big.NewInt(int64(overshotAmount)))

		retryErr := SendFunds(logger, client, payload)
		if retryErr == nil {
			return nil
		}

		if strings.Contains(retryErr.Error(), "overshot") {
			r.Retry(ctx, logger, client, retryErr, payload, currentAttempt+1)
		}
	}

	if r.nextRetrier != nil {
		return r.nextRetrier.Retry(ctx, logger, client, txErr, payload, 0)
	}

	return txErr
}

// ReturnFunds returns funds from the given chainlink nodes to the default network wallet. It will use a variety
// of strategies to attempt to return funds, including retrying with less funds if the transaction fails due to
// insufficient funds, and retrying with a higher gas limit if the transaction fails due to gas too low.
func ReturnFunds(log zerolog.Logger, seth *seth.Client, chainlinkNodes []contracts.ChainlinkNodeWithKeys) error {
	if seth == nil {
		return fmt.Errorf("Seth client is nil, unable to return funds from chainlink nodes")
	}
	log.Info().Msg("Attempting to return Chainlink node funds to default network wallets")
	if seth.Cfg.IsSimulatedNetwork() {
		log.Info().Str("Network Name", seth.Cfg.Network.Name).
			Msg("Network is a simulated network. Skipping fund return.")
		return nil
	}

	for _, chainlinkNode := range chainlinkNodes {
		fundedKeys, err := chainlinkNode.ExportEVMKeysForChain(fmt.Sprint(seth.ChainID))
		if err != nil {
			return err
		}
		for _, key := range fundedKeys {
			keyToDecrypt, err := json.Marshal(key)
			if err != nil {
				return err
			}
			// This can take up a good bit of RAM and time. When running on the remote-test-runner, this can lead to OOM
			// issues. So we avoid running in parallel; slower, but safer.
			decryptedKey, err := keystore.DecryptKey(keyToDecrypt, client.ChainlinkKeyPassword)
			if err != nil {
				return err
			}

			publicKey := decryptedKey.PrivateKey.Public()
			publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
			if !ok {
				return errors.New("error casting public key to ECDSA")
			}
			fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

			balance, err := seth.Client.BalanceAt(context.Background(), fromAddress, nil)
			if err != nil {
				return err
			}

			totalGasCost := new(big.Int).Mul(big.NewInt(0).SetUint64(seth.Cfg.Network.GasLimit), big.NewInt(0).SetInt64(seth.Cfg.Network.GasPrice))
			toSend := new(big.Int).Sub(balance, totalGasCost)

			payload := FundsToSendPayload{ToAddress: seth.Addresses[0], Amount: toSend, PrivateKey: decryptedKey.PrivateKey}

			err = SendFunds(log, seth, payload)
			if err != nil {
				handler := OvershotTransferRetrier{maxRetries: 3, nextRetrier: &InsufficientFundTransferRetrier{maxRetries: 3, nextRetrier: &GasTooLowTransferRetrier{maxGasLimit: seth.Cfg.Network.GasLimit * 3}}}
				return handler.Retry(context.Background(), log, seth, err, payload, 0)
			}
		}
	}

	return nil
}
