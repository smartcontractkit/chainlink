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
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/seth"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"

	clClient "github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

const (
	InsufficientFundsErr = "insufficient funds"
	GasTooLowErr         = "gas too low"
	OvershotErr          = "overshot"
)

var (
	RetrySuccessfulMsg = "Retry successful"
	NotSupportedMsg    = "Error not supported. Passing to next retrier"
)

// TransactionRetrier is an interface that every retrier of failed funds transfer transaction needs to implement
type TransactionRetrier interface {
	Retry(ctx context.Context, logger zerolog.Logger, client *seth.Client, txErr error, payload FundsToSendPayload, currentAttempt int) error
}

// InsufficientFundTransferRetrier will retry a failed funds transfer transaction if the error is due to insufficient funds
// by subtracting 1 Gwei from amount to send and retrying it up to maxRetries times
type InsufficientFundTransferRetrier struct {
	nextRetrier TransactionRetrier
	maxRetries  int
}

func (r *InsufficientFundTransferRetrier) Retry(ctx context.Context, logger zerolog.Logger, client *seth.Client, txErr error, payload FundsToSendPayload, currentAttempt int) error {
	if currentAttempt >= r.maxRetries {
		if r.nextRetrier != nil {
			logger.Debug().
				Str("retier", "InsufficientFundTransferRetrier").
				Msg("Max retries reached. Passing to next retrier")
			return r.nextRetrier.Retry(ctx, logger, client, txErr, payload, 0)
		}
		return txErr
	}

	for txErr != nil && (strings.Contains(txErr.Error(), InsufficientFundsErr)) {
		logger.Info().
			Msg("Insufficient funds error detected, retrying with less funds")

		newAmount := big.NewInt(0).Sub(payload.Amount, big.NewInt(blockchain.GWei))

		logger.Debug().
			Str("retier", "InsufficientFundTransferRetrier").
			Str("old amount", payload.Amount.String()).
			Str("new amount", newAmount.String()).
			Str("diff", big.NewInt(0).Sub(payload.Amount, newAmount).String()).
			Msg("New amount to send")

		payload.Amount = newAmount

		_, retryErr := SendFunds(logger, client, payload)
		if retryErr == nil {
			logger.Info().
				Str("retier", "InsufficientFundTransferRetrier").
				Msg(RetrySuccessfulMsg)
			return nil
		}

		if strings.Contains(retryErr.Error(), InsufficientFundsErr) {
			return r.Retry(ctx, logger, client, retryErr, payload, currentAttempt+1)
		}
	}

	if r.nextRetrier != nil {
		logger.Debug().
			Str("retier", "InsufficientFundTransferRetrier").
			Msg(NotSupportedMsg)
		return r.nextRetrier.Retry(ctx, logger, client, txErr, payload, 0)
	}

	logger.Warn().
		Str("retier", "InsufficientFundTransferRetrier").
		Msg("No more retriers available. Unable to retry transaction. Returning error.")

	return txErr
}

// GasTooLowTransferRetrier will retry a failed funds transfer transaction if the error is due to gas too low
// by doubling the gas limit and retrying until reaching maxGasLimit
type GasTooLowTransferRetrier struct {
	nextRetrier TransactionRetrier
	maxGasLimit int64
}

func (r *GasTooLowTransferRetrier) Retry(ctx context.Context, logger zerolog.Logger, client *seth.Client, txErr error, payload FundsToSendPayload, currentAttempt int) error {
	if payload.GasLimit != nil && *payload.GasLimit >= r.maxGasLimit {
		if r.nextRetrier != nil {
			logger.Debug().
				Str("retier", "GasTooLowTransferRetrier").
				Msg("Max gas limit reached. Passing to next retrier")
			return r.nextRetrier.Retry(ctx, logger, client, txErr, payload, 0)
		}
		return txErr
	}

	for txErr != nil && strings.Contains(txErr.Error(), GasTooLowErr) {
		logger.Info().
			Msg("Too low gas error detected, retrying with more gas")
		var newGasLimit int64
		if payload.GasLimit != nil {
			newGasLimit = *payload.GasLimit * 2
		} else {
			newGasLimit = client.Cfg.Network.TransferGasFee * 2
		}

		logger.Debug().
			Str("retier", "GasTooLowTransferRetrier").
			Int64("old gas limit", newGasLimit/2).
			Int64("new gas limit", newGasLimit).
			Int64("diff", newGasLimit).
			Msg("New gas limit to use")

		payload.GasLimit = &newGasLimit

		_, retryErr := SendFunds(logger, client, payload)
		if retryErr == nil {
			logger.Info().
				Str("retier", "GasTooLowTransferRetrier").
				Msg(RetrySuccessfulMsg)
			return nil
		}

		if strings.Contains(retryErr.Error(), GasTooLowErr) {
			return r.Retry(ctx, logger, client, retryErr, payload, currentAttempt+1)
		}
	}

	if r.nextRetrier != nil {
		logger.Debug().
			Str("retier", "OvershotTransferRetrier").
			Msg(NotSupportedMsg)
		return r.nextRetrier.Retry(ctx, logger, client, txErr, payload, 0)
	}

	logger.Warn().
		Str("retier", "OvershotTransferRetrier").
		Msg("No more retriers available. Unable to retry transaction. Returning error.")

	return txErr
}

// OvershotTransferRetrier will retry a failed funds transfer transaction if the error is due to overshot
// by subtracting the overshot amount from the amount to send and retrying it up to maxRetries times
type OvershotTransferRetrier struct {
	nextRetrier TransactionRetrier
	maxRetries  int
}

func (r *OvershotTransferRetrier) Retry(ctx context.Context, logger zerolog.Logger, client *seth.Client, txErr error, payload FundsToSendPayload, currentAttempt int) error {
	if currentAttempt >= r.maxRetries {
		logger.Debug().
			Str("retier", "OvershotTransferRetrier").
			Msg("Max retries reached. Passing to next retrier")
		if r.nextRetrier != nil {
			return r.nextRetrier.Retry(ctx, logger, client, txErr, payload, 0)
		}
		return txErr
	}

	overshotRe := regexp.MustCompile(`overshot (\d+)`)
	if txErr != nil && strings.Contains(txErr.Error(), OvershotErr) {
		logger.Info().
			Msg("Overshot error detected, retrying with less funds")
		submatches := overshotRe.FindStringSubmatch(txErr.Error())
		if len(submatches) < 1 {
			return fmt.Errorf("error parsing overshot amount in error: %w", txErr)
		}
		numberString := submatches[1]
		overshotAmount, err := strconv.Atoi(numberString)
		if err != nil {
			return err
		}

		newAmount := big.NewInt(0).Sub(payload.Amount, big.NewInt(int64(overshotAmount)))
		logger.Debug().
			Str("retier", "OvershotTransferRetrier").
			Str("old amount", payload.Amount.String()).
			Str("new amount", newAmount.String()).
			Str("diff", big.NewInt(0).Sub(payload.Amount, newAmount).String()).
			Msg("New amount to send")

		payload.Amount = newAmount

		_, retryErr := SendFunds(logger, client, payload)
		if retryErr == nil {
			logger.Info().
				Str("retier", "OvershotTransferRetrier").
				Msg(RetrySuccessfulMsg)
			return nil
		}

		if strings.Contains(retryErr.Error(), OvershotErr) {
			return r.Retry(ctx, logger, client, retryErr, payload, currentAttempt+1)
		}
	}

	if r.nextRetrier != nil {
		logger.Debug().
			Str("retier", "OvershotTransferRetrier").
			Msg(NotSupportedMsg)
		return r.nextRetrier.Retry(ctx, logger, client, txErr, payload, 0)
	}

	return txErr
}

// ReturnFunds returns funds from the given chainlink nodes to the default network wallet. It will use a variety
// of strategies to attempt to return funds, including retrying with less funds if the transaction fails due to
// insufficient funds, and retrying with a higher gas limit if the transaction fails due to gas too low.
func ReturnFunds(log zerolog.Logger, sethClient *seth.Client, chainlinkNodes []contracts.ChainlinkNodeWithKeysAndAddress) error {
	if sethClient == nil {
		return fmt.Errorf("Seth client is nil, unable to return funds from chainlink nodes")
	}
	log.Info().Msg("Attempting to return Chainlink node funds to default network wallets")
	if sethClient.Cfg.IsSimulatedNetwork() {
		log.Info().Str("Network Name", sethClient.Cfg.Network.Name).
			Msg("Network is a simulated network. Skipping fund return.")
		return nil
	}

	failedReturns := []common.Address{}

	for _, chainlinkNode := range chainlinkNodes {
		fundedKeys, err := chainlinkNode.ExportEVMKeysForChain(fmt.Sprint(sethClient.ChainID))
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
			decryptedKey, err := keystore.DecryptKey(keyToDecrypt, clClient.ChainlinkKeyPassword)
			if err != nil {
				return err
			}

			publicKey := decryptedKey.PrivateKey.Public()
			publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
			if !ok {
				return errors.New("error casting public key to ECDSA")
			}
			fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

			balance, err := sethClient.Client.BalanceAt(context.Background(), fromAddress, nil)
			if err != nil {
				return err
			}

			if balance.Cmp(big.NewInt(0)) == 0 {
				log.Info().
					Str("Address", fromAddress.String()).
					Msg("No balance to return. Skipping return.")
			}

			// if not set, it will be just set to empty string, which is okay as long as gas estimation is disabled
			txPriority := sethClient.Cfg.Network.GasEstimationTxPriority
			txTimeout := sethClient.Cfg.Network.TxnTimeout.Duration()

			if sethClient.Cfg.IsExperimentEnabled(seth.Experiment_SlowFundsReturn) {
				txPriority = "slow"
				thirtyMinutes := time.Duration(30 * time.Minute)
				txTimeout = thirtyMinutes
			}

			estimations := sethClient.CalculateGasEstimations(seth.GasEstimationRequest{
				GasEstimationEnabled: sethClient.Cfg.Network.GasEstimationEnabled,
				FallbackGasPrice:     sethClient.Cfg.Network.GasPrice,
				FallbackGasFeeCap:    sethClient.Cfg.Network.GasFeeCap,
				FallbackGasTipCap:    sethClient.Cfg.Network.GasTipCap,
				Priority:             txPriority,
			})

			var maxTotalGasCost *big.Int
			if sethClient.Cfg.Network.EIP1559DynamicFees {
				maxTotalGasCost = new(big.Int).Mul(big.NewInt(0).SetInt64(sethClient.Cfg.Network.TransferGasFee), estimations.GasFeeCap)
			} else {
				maxTotalGasCost = new(big.Int).Mul(big.NewInt(0).SetInt64(sethClient.Cfg.Network.TransferGasFee), estimations.GasPrice)
			}

			toSend := new(big.Int).Sub(balance, maxTotalGasCost)

			if toSend.Cmp(big.NewInt(0)) <= 0 {
				log.Warn().
					Str("Address", fromAddress.String()).
					Str("Estimated maximum total gas cost", maxTotalGasCost.String()).
					Str("Balance", balance.String()).
					Str("To send", toSend.String()).
					Msg("Not enough balance to cover gas cost. Skipping return.")

				failedReturns = append(failedReturns, fromAddress)
				continue
			}

			payload := FundsToSendPayload{
				ToAddress:  sethClient.Addresses[0],
				Amount:     toSend,
				PrivateKey: decryptedKey.PrivateKey,
				GasLimit:   &sethClient.Cfg.Network.TransferGasFee,
				GasPrice:   estimations.GasPrice,
				GasFeeCap:  estimations.GasFeeCap,
				GasTipCap:  estimations.GasTipCap,
				TxTimeout:  &txTimeout,
			}

			_, err = SendFunds(log, sethClient, payload)
			if err != nil {
				handler := OvershotTransferRetrier{maxRetries: 10, nextRetrier: &InsufficientFundTransferRetrier{maxRetries: 10, nextRetrier: &GasTooLowTransferRetrier{maxGasLimit: sethClient.Cfg.Network.TransferGasFee * 10}}}
				err = handler.Retry(context.Background(), log, sethClient, err, payload, 0)
				if err != nil {
					log.Error().
						Err(err).
						Str("Address", fromAddress.String()).
						Msg("Failed to return funds from Chainlink node to default network wallet")
					failedReturns = append(failedReturns, fromAddress)
				}
			}
		}
	}

	if len(failedReturns) > 0 {
		return fmt.Errorf("failed to return funds from Chainlink nodes to default network wallet for addresses: %v", failedReturns)
	}

	log.Info().Msg("Successfully returned funds from all Chainlink nodes to default network wallets")

	return nil
}
