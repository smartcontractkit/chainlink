package actions

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/testreporters"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/conversions"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

type FundsToSendPayload struct {
	ToAddress  common.Address
	Amount     *big.Int
	PrivateKey *ecdsa.PrivateKey
	GasLimit   *int64
	GasPrice   *big.Int
	GasFeeCap  *big.Int
	GasTipCap  *big.Int
	TxTimeout  *time.Duration
}

func PrivateKeyToAddress(privateKey *ecdsa.PrivateKey) (common.Address, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, errors.New("error casting public key to ECDSA")
	}
	return crypto.PubkeyToAddress(*publicKeyECDSA), nil
}

// SendFunds sends native funds from a private key to a target address.
func SendFunds(logger zerolog.Logger, client *seth.Client, payload FundsToSendPayload) (*types.Receipt, error) {
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
	if payload.GasLimit != nil {
		if *payload.GasLimit < 0 {
			return nil, fmt.Errorf("negative gas limit: %d", *payload.GasLimit)
		}
		gasLimit = uint64(*payload.GasLimit)
	}

	gasPrice := big.NewInt(0)
	gasFeeCap := big.NewInt(0)
	gasTipCap := big.NewInt(0)
	if client.Cfg.Network.EIP1559DynamicFees {
		if payload.GasFeeCap == nil || payload.GasTipCap == nil {
			txOptions := client.NewTXOpts(seth.WithGasLimit(gasLimit))
			gasFeeCap = txOptions.GasFeeCap
			gasTipCap = txOptions.GasTipCap
		}
		if payload.GasFeeCap != nil {
			gasFeeCap = payload.GasFeeCap
		}
		if payload.GasTipCap != nil {
			gasTipCap = payload.GasTipCap
		}
	} else if payload.GasPrice == nil {
		txOptions := client.NewTXOpts(seth.WithGasLimit(gasLimit))
		gasPrice = txOptions.GasPrice
	} else {
		gasPrice = payload.GasPrice
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
	if err := client.Client.SendTransaction(ctx, signedTx); err != nil {
		return nil, errors.Wrap(err, "failed to send transaction")
	}

	receipt, receiptErr := client.WaitMined(ctx, logger, client.Client, signedTx)
	if receiptErr != nil {
		return nil, errors.Wrap(receiptErr, "failed to wait for transaction to be mined")
	}
	if receipt.Status == 1 {
		return receipt, nil
	}

	tx, _, err := client.Client.TransactionByHash(ctx, signedTx.Hash())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get transaction by hash")
	}
	_, err = client.Decode(tx, receiptErr)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

// ReturnFunds attempts to return all the funds from the chainlink nodes to the network's default address.
func ReturnFunds(lggr zerolog.Logger, chainlinkNodes []*nodeclient.ChainlinkK8sClient, blockchainClient blockchain.EVMClient) error {
	if blockchainClient == nil {
		return errors.New("blockchain client is nil, unable to return funds from chainlink nodes")
	}
	lggr.Info().Msg("Attempting to return Chainlink node funds to default network wallets")
	if blockchainClient.NetworkSimulated() {
		lggr.Info().Str("Network Name", blockchainClient.GetNetworkName()).
			Msg("Network is a simulated network. Skipping fund return.")
		return nil
	}

	for _, chainlinkNode := range chainlinkNodes {
		fundedKeys, err := chainlinkNode.ExportEVMKeysForChain(blockchainClient.GetChainID().String())
		if err != nil {
			return err
		}
		for _, key := range fundedKeys {
			keyToDecrypt, err := json.Marshal(key)
			if err != nil {
				return err
			}
			decryptedKey, err := keystore.DecryptKey(keyToDecrypt, nodeclient.ChainlinkKeyPassword)
			if err != nil {
				return err
			}
			err = blockchainClient.ReturnFunds(decryptedKey.PrivateKey)
			if err != nil {
				lggr.Error().Err(err).Str("Address", fundedKeys[0].Address).Msg("Error returning funds from Chainlink node")
			}
		}
	}
	return blockchainClient.WaitForEvents()
}

// TeardownSuite tears down networks/clients and environment and creates logs for failed tests.
func TeardownSuite(
	t *testing.T,
	chainClient *seth.Client,
	env *environment.Environment,
	chainlinkNodes []*nodeclient.ChainlinkK8sClient,
	optionalTestReporter testreporters.TestReporter,
	failingLogLevel zapcore.Level,
	grafnaUrlProvider testreporters.GrafanaURLProvider,
	evmClients ...blockchain.EVMClient,
) error {
	l := logging.GetTestLogger(t)
	if err := testreporters.WriteTeardownLogs(t, env, optionalTestReporter, failingLogLevel, grafnaUrlProvider); err != nil {
		return fmt.Errorf("error dumping environment logs, leaving environment running for manual retrieval, err: %w", err)
	}
	err := DeleteAllJobs(chainlinkNodes)
	if err != nil {
		l.Warn().Msgf("Error deleting jobs %+v", err)
	}

	if chainlinkNodes != nil && chainClient != nil {
		if err := ReturnFundsFromNodes(l, chainClient, contracts.ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(chainlinkNodes)); err != nil {
			fmt.Println(environment.FAILED_FUND_RETURN)
			l.Error().Err(err).Str("Namespace", env.Cfg.Namespace).
				Msg("Error attempting to return funds from chainlink nodes to network's default wallet. Environment is left running so you can try manually!")
		}
	} else {
		l.Info().Msg("Successfully returned funds from chainlink nodes to default network wallets")
	}

	for _, c := range evmClients {
		if c != nil && chainlinkNodes != nil && len(chainlinkNodes) > 0 {
			if err := ReturnFunds(l, chainlinkNodes, c); err != nil {
				fmt.Println(environment.FAILED_FUND_RETURN)
				l.Error().Err(err).Str("Namespace", env.Cfg.Namespace).
					Msg("Error attempting to return funds from chainlink nodes to network's default wallet. Environment is left running so you can try manually!")
			}
		} else {
			l.Info().Msg("Successfully returned funds from chainlink nodes to default network wallets")
		}
		if c != nil {
			if err := c.Close(); err != nil {
				return err
			}
		}
	}

	return env.Shutdown()
}

// DeleteAllJobs deletes all jobs from all chainlink nodes.
func DeleteAllJobs(chainlinkNodes []*nodeclient.ChainlinkK8sClient) error {
	for _, node := range chainlinkNodes {
		if node == nil {
			return fmt.Errorf("found a nil chainlink node in the list of chainlink nodes while tearing down: %v", chainlinkNodes)
		}
		jobs, _, err := node.ReadJobs()
		if err != nil {
			return fmt.Errorf("error reading jobs from chainlink node, err: %w", err)
		}
		for _, maps := range jobs.Data {
			if _, ok := maps["id"]; !ok {
				return fmt.Errorf("error reading job id from chainlink node's jobs %+v", jobs.Data)
			}
			id := maps["id"].(string)
			if _, err := node.DeleteJob(id); err != nil {
				return fmt.Errorf("error deleting job from chainlink node, err: %w", err)
			}
		}
	}
	return nil
}
