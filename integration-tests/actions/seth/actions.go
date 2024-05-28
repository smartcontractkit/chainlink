package actions_seth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math"
	"math/big"
	"strings"
	"testing"
	"time"

	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/seth"
	"github.com/test-go/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/operator_factory"

	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/conversions"
	seth_utils "github.com/smartcontractkit/chainlink-testing-framework/utils/seth"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
)

var ContractDeploymentInterval = 200

// FundChainlinkNodesFromRootAddress sends native token amount (expressed in human-scale) to each Chainlink Node
// from root private key. It returns an error if any of the transactions failed.
func FundChainlinkNodesFromRootAddress(
	logger zerolog.Logger,
	client *seth.Client,
	nodes []contracts.ChainlinkNodeWithKeysAndAddress,
	amount *big.Float,
) error {
	if len(client.PrivateKeys) == 0 {
		return errors.Wrap(errors.New(seth.ErrNoKeyLoaded), fmt.Sprintf("requested key: %d", 0))
	}

	return FundChainlinkNodes(logger, client, nodes, client.PrivateKeys[0], amount)
}

// FundChainlinkNodes sends native token amount (expressed in human-scale) to each Chainlink Node
// from private key's address. It returns an error if any of the transactions failed.
func FundChainlinkNodes(
	logger zerolog.Logger,
	client *seth.Client,
	nodes []contracts.ChainlinkNodeWithKeysAndAddress,
	privateKey *ecdsa.PrivateKey,
	amount *big.Float,
) error {
	keyAddressFn := func(cl contracts.ChainlinkNodeWithKeysAndAddress) (string, error) {
		return cl.PrimaryEthAddress()
	}
	return fundChainlinkNodesAtAnyKey(logger, client, nodes, privateKey, amount, keyAddressFn)
}

// FundChainlinkNodesAtKeyIndexFromRootAddress sends native token amount (expressed in human-scale) to each Chainlink Node
// from root private key.It returns an error if any of the transactions failed. It sends the funds to
// node address at keyIndex (as each node can have multiple addresses).
func FundChainlinkNodesAtKeyIndexFromRootAddress(
	logger zerolog.Logger,
	client *seth.Client,
	nodes []contracts.ChainlinkNodeWithKeysAndAddress,
	amount *big.Float,
	keyIndex int,
) error {
	if len(client.PrivateKeys) == 0 {
		return errors.Wrap(errors.New(seth.ErrNoKeyLoaded), fmt.Sprintf("requested key: %d", 0))
	}

	return FundChainlinkNodesAtKeyIndex(logger, client, nodes, client.PrivateKeys[0], amount, keyIndex)
}

// FundChainlinkNodesAtKeyIndex sends native token amount (expressed in human-scale) to each Chainlink Node
// from private key's address. It returns an error if any of the transactions failed. It sends the funds to
// node address at keyIndex (as each node can have multiple addresses).
func FundChainlinkNodesAtKeyIndex(
	logger zerolog.Logger,
	client *seth.Client,
	nodes []contracts.ChainlinkNodeWithKeysAndAddress,
	privateKey *ecdsa.PrivateKey,
	amount *big.Float,
	keyIndex int,
) error {
	keyAddressFn := func(cl contracts.ChainlinkNodeWithKeysAndAddress) (string, error) {
		toAddress, err := cl.EthAddresses()
		if err != nil {
			return "", err
		}
		return toAddress[keyIndex], nil
	}
	return fundChainlinkNodesAtAnyKey(logger, client, nodes, privateKey, amount, keyAddressFn)
}

func fundChainlinkNodesAtAnyKey(
	logger zerolog.Logger,
	client *seth.Client,
	nodes []contracts.ChainlinkNodeWithKeysAndAddress,
	privateKey *ecdsa.PrivateKey,
	amount *big.Float,
	keyAddressFn func(contracts.ChainlinkNodeWithKeysAndAddress) (string, error),
) error {
	for _, cl := range nodes {
		toAddress, err := keyAddressFn(cl)
		if err != nil {
			return err
		}

		fromAddress, err := privateKeyToAddress(privateKey)
		if err != nil {
			return err
		}

		receipt, err := SendFunds(logger, client, FundsToSendPayload{
			ToAddress:  common.HexToAddress(toAddress),
			Amount:     conversions.EtherToWei(amount),
			PrivateKey: privateKey,
		})
		if err != nil {
			logger.Err(err).
				Str("From", fromAddress.Hex()).
				Str("To", toAddress).
				Msg("Failed to fund Chainlink node")

			return err
		}

		txHash := "(none)"
		if receipt != nil {
			txHash = receipt.TxHash.String()
		}

		logger.Info().
			Str("From", fromAddress.Hex()).
			Str("To", toAddress).
			Str("TxHash", txHash).
			Str("Amount", amount.String()).
			Msg("Funded Chainlink node")
	}

	return nil
}

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

// TODO: move to CTF?
// SendFunds sends native token amount (expressed in human-scale) from address controlled by private key
// to given address. You can override any or none of the following: gas limit, gas price, gas fee cap, gas tip cap.
// Values that are not set will be estimated or taken from config.
func SendFunds(logger zerolog.Logger, client *seth.Client, payload FundsToSendPayload) (*types.Receipt, error) {
	fromAddress, err := privateKeyToAddress(payload.PrivateKey)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), client.Cfg.Network.TxnTimeout.Duration())
	nonce, err := client.Client.PendingNonceAt(ctx, fromAddress)
	defer cancel()
	if err != nil {
		return nil, err
	}

	var gasLimit int64
	gasLimitRaw, err := client.EstimateGasLimitForFundTransfer(fromAddress, payload.ToAddress, payload.Amount)
	if err != nil {
		gasLimit = client.Cfg.Network.TransferGasFee
	} else {
		gasLimit = int64(gasLimitRaw)
	}

	gasPrice := big.NewInt(0)
	gasFeeCap := big.NewInt(0)
	gasTipCap := big.NewInt(0)

	if payload.GasLimit != nil {
		gasLimit = *payload.GasLimit
	}

	if client.Cfg.Network.EIP1559DynamicFees {
		// if any of the dynamic fees are not set, we need to either estimate them or read them from config
		if payload.GasFeeCap == nil || payload.GasTipCap == nil {
			// estimation or config reading happens here
			txOptions := client.NewTXOpts(seth.WithGasLimit(uint64(gasLimit)))
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
			txOptions := client.NewTXOpts(seth.WithGasLimit(uint64(gasLimit)))
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
			Gas:       uint64(gasLimit),
			GasFeeCap: gasFeeCap,
			GasTipCap: gasTipCap,
		}
	} else {
		rawTx = &types.LegacyTx{
			Nonce:    nonce,
			To:       &payload.ToAddress,
			Value:    payload.Amount,
			Gas:      uint64(gasLimit),
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
		Int64("Gas Limit", gasLimit).
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
		Int64("Gas Limit", gasLimit).
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

// DeployForwarderContracts first deploys Operator Factory and then uses it to deploy given number of
// operator and forwarder pairs. It waits for each transaction to be mined and then extracts operator and
// forwarder addresses from emitted events.
func DeployForwarderContracts(
	t *testing.T,
	seth *seth.Client,
	linkTokenAddress common.Address,
	numberOfOperatorForwarderPairs int,
) (operators []common.Address, authorizedForwarders []common.Address, operatorFactoryInstance contracts.OperatorFactory) {
	instance, err := contracts.DeployEthereumOperatorFactory(seth, linkTokenAddress)
	require.NoError(t, err, "failed to create new instance of operator factory")
	operatorFactoryInstance = &instance

	for i := 0; i < numberOfOperatorForwarderPairs; i++ {
		decodedTx, err := seth.Decode(operatorFactoryInstance.DeployNewOperatorAndForwarder())
		require.NoError(t, err, "Deploying new operator with proposed ownership with forwarder shouldn't fail")

		for i, event := range decodedTx.Events {
			require.True(t, len(event.Topics) > 0, fmt.Sprintf("Event %d should have topics", i))
			switch event.Topics[0] {
			case operator_factory.OperatorFactoryOperatorCreated{}.Topic().String():
				if address, ok := event.EventData["operator"]; ok {
					operators = append(operators, address.(common.Address))
				} else {
					require.Fail(t, "Operator address not found in event", event)
				}
			case operator_factory.OperatorFactoryAuthorizedForwarderCreated{}.Topic().String():
				if address, ok := event.EventData["forwarder"]; ok {
					authorizedForwarders = append(authorizedForwarders, address.(common.Address))
				} else {
					require.Fail(t, "Forwarder address not found in event", event)
				}
			}
		}
	}
	return operators, authorizedForwarders, operatorFactoryInstance
}

// WatchNewOCRRound watches for a new OCR round, similarly to StartNewRound, but it does not explicitly request a new
// round from the contract, as this can cause some odd behavior in some cases. It announces success if latest round
// is >= roundNumber.
func WatchNewOCRRound(
	l zerolog.Logger,
	seth *seth.Client,
	roundNumber int64,
	ocrInstances []contracts.OffChainAggregatorWithRounds,
	timeout time.Duration,
) error {
	confirmed := make(map[string]bool)
	timeoutC := time.After(timeout)
	ticker := time.NewTicker(time.Millisecond * 200)
	defer ticker.Stop()

	l.Info().Msgf("Waiting for round %d to be confirmed by all nodes", roundNumber)

	for {
		select {
		case <-timeoutC:
			return fmt.Errorf("timeout waiting for round %d to be confirmed. %d/%d nodes confirmed it", roundNumber, len(confirmed), len(ocrInstances))
		case <-ticker.C:
			for i := 0; i < len(ocrInstances); i++ {
				if confirmed[ocrInstances[i].Address()] {
					continue
				}
				ctx, cancel := context.WithTimeout(context.Background(), seth.Cfg.Network.TxnTimeout.Duration())
				roundData, err := ocrInstances[i].GetLatestRound(ctx)
				if err != nil {
					cancel()
					return fmt.Errorf("getting latest round from OCR instance %d have failed: %w", i+1, err)
				}
				cancel()
				if roundData.RoundId.Cmp(big.NewInt(roundNumber)) >= 0 {
					l.Debug().Msgf("OCR instance %d/%d confirmed round %d", i+1, len(ocrInstances), roundNumber)
					confirmed[ocrInstances[i].Address()] = true
				}
			}
			if len(confirmed) == len(ocrInstances) {
				return nil
			}
		}
	}
}

// AcceptAuthorizedReceiversOperator sets authorized receivers for each operator contract to
// authorizedForwarder and authorized EA to nodeAddresses. Once done, it confirms that authorizations
// were set correctly.
func AcceptAuthorizedReceiversOperator(
	t *testing.T,
	logger zerolog.Logger,
	seth *seth.Client,
	operator common.Address,
	authorizedForwarder common.Address,
	nodeAddresses []common.Address,
) {
	operatorInstance, err := contracts.LoadEthereumOperator(logger, seth, operator)
	require.NoError(t, err, "Loading operator contract shouldn't fail")
	forwarderInstance, err := contracts.LoadEthereumAuthorizedForwarder(seth, authorizedForwarder)
	require.NoError(t, err, "Loading authorized forwarder contract shouldn't fail")

	err = operatorInstance.AcceptAuthorizedReceivers([]common.Address{authorizedForwarder}, nodeAddresses)
	require.NoError(t, err, "Accepting authorized forwarder shouldn't fail")

	senders, err := forwarderInstance.GetAuthorizedSenders(testcontext.Get(t))
	require.NoError(t, err, "Getting authorized senders shouldn't fail")
	var nodesAddrs []string
	for _, o := range nodeAddresses {
		nodesAddrs = append(nodesAddrs, o.Hex())
	}
	require.Equal(t, nodesAddrs, senders, "Senders addresses should match node addresses")

	owner, err := forwarderInstance.Owner(testcontext.Get(t))
	require.NoError(t, err, "Getting authorized forwarder owner shouldn't fail")
	require.Equal(t, operator.Hex(), owner, "Forwarder owner should match operator")
}

// TrackForwarder creates forwarder track for a given Chainlink node
func TrackForwarder(
	t *testing.T,
	seth *seth.Client,
	authorizedForwarder common.Address,
	node contracts.ChainlinkNodeWithForwarder,
) {
	l := logging.GetTestLogger(t)
	chainID := big.NewInt(seth.ChainID)
	_, _, err := node.TrackForwarder(chainID, authorizedForwarder)
	require.NoError(t, err, "Forwarder track should be created")
	l.Info().Str("NodeURL", node.GetConfig().URL).
		Str("ForwarderAddress", authorizedForwarder.Hex()).
		Str("ChaindID", chainID.String()).
		Msg("Forwarder tracked")
}

// DeployOCRv2Contracts deploys a number of OCRv2 contracts and configures them with defaults
func DeployOCRv2Contracts(
	l zerolog.Logger,
	seth *seth.Client,
	numberOfContracts int,
	linkTokenAddress common.Address,
	transmitters []string,
	ocrOptions contracts.OffchainOptions,
) ([]contracts.OffchainAggregatorV2, error) {
	var ocrInstances []contracts.OffchainAggregatorV2
	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		ocrInstance, err := contracts.DeployOffchainAggregatorV2(
			l,
			seth,
			linkTokenAddress,
			ocrOptions,
		)
		if err != nil {
			return nil, fmt.Errorf("OCRv2 instance deployment have failed: %w", err)
		}
		ocrInstances = append(ocrInstances, &ocrInstance)
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			time.Sleep(2 * time.Second)
		}
	}

	// Gather address payees
	var payees []string
	for range transmitters {
		payees = append(payees, seth.Addresses[0].Hex())
	}

	// Set Payees
	for contractCount, ocrInstance := range ocrInstances {
		err := ocrInstance.SetPayees(transmitters, payees)
		if err != nil {
			return nil, fmt.Errorf("error settings OCR payees: %w", err)
		}
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			time.Sleep(2 * time.Second)
		}
	}
	return ocrInstances, nil
}

// ConfigureOCRv2AggregatorContracts sets configuration for a number of OCRv2 contracts
func ConfigureOCRv2AggregatorContracts(
	contractConfig *contracts.OCRv2Config,
	ocrv2Contracts []contracts.OffchainAggregatorV2,
) error {
	for contractCount, ocrInstance := range ocrv2Contracts {
		// Exclude the first node, which will be used as a bootstrapper
		err := ocrInstance.SetConfig(contractConfig)
		if err != nil {
			return fmt.Errorf("error setting OCR config for contract '%s': %w", ocrInstance.Address(), err)
		}
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			time.Sleep(2 * time.Second)
		}
	}
	return nil
}

// TeardownSuite tears down networks/clients and environment and creates a logs folder for failed tests in the
// specified path. Can also accept a testreporter (if one was used) to log further results
func TeardownSuite(
	t *testing.T,
	chainClient *seth.Client,
	env *environment.Environment,
	chainlinkNodes []*client.ChainlinkK8sClient,
	optionalTestReporter testreporters.TestReporter, // Optionally pass in a test reporter to log further metrics
	failingLogLevel zapcore.Level, // Examines logs after the test, and fails the test if any Chainlink logs are found at or above provided level
	grafnaUrlProvider testreporters.GrafanaURLProvider,
) error {
	l := logging.GetTestLogger(t)
	if err := testreporters.WriteTeardownLogs(t, env, optionalTestReporter, failingLogLevel, grafnaUrlProvider); err != nil {
		return fmt.Errorf("Error dumping environment logs, leaving environment running for manual retrieval, err: %w", err)
	}
	// Delete all jobs to stop depleting the funds
	err := DeleteAllJobs(chainlinkNodes)
	if err != nil {
		l.Warn().Msgf("Error deleting jobs %+v", err)
	}

	if chainlinkNodes != nil {
		if err := ReturnFundsFromNodes(l, chainClient, contracts.ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(chainlinkNodes)); err != nil {
			// This printed line is required for tests that use real funds to propagate the failure
			// out to the system running the test. Do not remove
			fmt.Println(environment.FAILED_FUND_RETURN)
			l.Error().Err(err).Str("Namespace", env.Cfg.Namespace).
				Msg("Error attempting to return funds from chainlink nodes to network's default wallet. " +
					"Environment is left running so you can try manually!")
		}
	} else {
		l.Info().Msg("Successfully returned funds from chainlink nodes to default network wallets")
	}

	return env.Shutdown()
}

// TeardownRemoteSuite sends a report and returns funds from chainlink nodes to network's default wallet
func TeardownRemoteSuite(
	t *testing.T,
	client *seth.Client,
	namespace string,
	chainlinkNodes []*client.ChainlinkK8sClient,
	optionalTestReporter testreporters.TestReporter, // Optionally pass in a test reporter to log further metrics
	grafnaUrlProvider testreporters.GrafanaURLProvider,
) error {
	l := logging.GetTestLogger(t)
	if err := testreporters.SendReport(t, namespace, "./", optionalTestReporter, grafnaUrlProvider); err != nil {
		l.Warn().Err(err).Msg("Error writing test report")
	}
	// Delete all jobs to stop depleting the funds
	err := DeleteAllJobs(chainlinkNodes)
	if err != nil {
		l.Warn().Msgf("Error deleting jobs %+v", err)
	}

	if err = ReturnFundsFromNodes(l, client, contracts.ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(chainlinkNodes)); err != nil {
		l.Error().Err(err).Str("Namespace", namespace).
			Msg("Error attempting to return funds from chainlink nodes to network's default wallet. " +
				"Environment is left running so you can try manually!")
	}

	return err
}

// DeleteAllJobs deletes all jobs from all chainlink nodes
// added here temporarily to avoid circular import
func DeleteAllJobs(chainlinkNodes []*client.ChainlinkK8sClient) error {
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
			_, err := node.DeleteJob(id)
			if err != nil {
				return fmt.Errorf("error deleting job from chainlink node, err: %w", err)
			}
		}
	}
	return nil
}

// StartNewRound requests a new round from the ocr contracts and returns once transaction was mined
func StartNewRound(
	ocrInstances []contracts.OffChainAggregatorWithRounds,
) error {
	for i := 0; i < len(ocrInstances); i++ {
		err := ocrInstances[i].RequestNewRound()
		if err != nil {
			return fmt.Errorf("requesting new OCR round %d have failed: %w", i+1, err)
		}
	}
	return nil
}

// DeployOCRContractsForwarderFlow deploys and funds a certain number of offchain
// aggregator contracts with forwarders as effectiveTransmitters
func DeployOCRContractsForwarderFlow(
	logger zerolog.Logger,
	seth *seth.Client,
	numberOfContracts int,
	linkTokenContractAddress common.Address,
	workerNodes []contracts.ChainlinkNodeWithKeysAndAddress,
	forwarderAddresses []common.Address,
) ([]contracts.OffchainAggregator, error) {
	transmitterPayeesFn := func() (transmitters []string, payees []string, err error) {
		transmitters = make([]string, 0)
		payees = make([]string, 0)
		for _, forwarderCommonAddress := range forwarderAddresses {
			forwarderAddress := forwarderCommonAddress.Hex()
			transmitters = append(transmitters, forwarderAddress)
			payees = append(payees, seth.Addresses[0].Hex())
		}

		return
	}

	transmitterAddressesFn := func() ([]common.Address, error) {
		return forwarderAddresses, nil
	}

	return deployAnyOCRv1Contracts(logger, seth, numberOfContracts, linkTokenContractAddress, workerNodes, transmitterPayeesFn, transmitterAddressesFn)
}

// DeployOCRv1Contracts deploys and funds a certain number of offchain aggregator contracts
func DeployOCRv1Contracts(
	logger zerolog.Logger,
	seth *seth.Client,
	numberOfContracts int,
	linkTokenContractAddress common.Address,
	workerNodes []contracts.ChainlinkNodeWithKeysAndAddress,
) ([]contracts.OffchainAggregator, error) {
	transmitterPayeesFn := func() (transmitters []string, payees []string, err error) {
		transmitters = make([]string, 0)
		payees = make([]string, 0)
		for _, node := range workerNodes {
			var addr string
			addr, err = node.PrimaryEthAddress()
			if err != nil {
				err = fmt.Errorf("error getting node's primary ETH address: %w", err)
				return
			}
			transmitters = append(transmitters, addr)
			payees = append(payees, seth.Addresses[0].Hex())
		}

		return
	}

	transmitterAddressesFn := func() ([]common.Address, error) {
		transmitterAddresses := make([]common.Address, 0)
		for _, node := range workerNodes {
			primaryAddress, err := node.PrimaryEthAddress()
			if err != nil {
				return nil, err
			}
			transmitterAddresses = append(transmitterAddresses, common.HexToAddress(primaryAddress))
		}

		return transmitterAddresses, nil
	}

	return deployAnyOCRv1Contracts(logger, seth, numberOfContracts, linkTokenContractAddress, workerNodes, transmitterPayeesFn, transmitterAddressesFn)
}

func deployAnyOCRv1Contracts(
	logger zerolog.Logger,
	seth *seth.Client,
	numberOfContracts int,
	linkTokenContractAddress common.Address,
	workerNodes []contracts.ChainlinkNodeWithKeysAndAddress,
	getTransmitterAndPayeesFn func() ([]string, []string, error),
	getTransmitterAddressesFn func() ([]common.Address, error),
) ([]contracts.OffchainAggregator, error) {
	// Deploy contracts
	var ocrInstances []contracts.OffchainAggregator
	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		ocrInstance, err := contracts.DeployOffchainAggregator(logger, seth, linkTokenContractAddress, contracts.DefaultOffChainAggregatorOptions())
		if err != nil {
			return nil, fmt.Errorf("OCR instance deployment have failed: %w", err)
		}
		ocrInstances = append(ocrInstances, &ocrInstance)
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			time.Sleep(2 * time.Second)
		}
	}

	// Gather transmitter and address payees
	var transmitters, payees []string
	var err error
	transmitters, payees, err = getTransmitterAndPayeesFn()
	if err != nil {
		return nil, fmt.Errorf("error getting transmitter and payees: %w", err)
	}

	// Set Payees
	for contractCount, ocrInstance := range ocrInstances {
		err := ocrInstance.SetPayees(transmitters, payees)
		if err != nil {
			return nil, fmt.Errorf("error settings OCR payees: %w", err)
		}
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			time.Sleep(2 * time.Second)
		}
	}

	// Set Config
	transmitterAddresses, err := getTransmitterAddressesFn()
	if err != nil {
		return nil, fmt.Errorf("getting transmitter addresses should not fail: %w", err)
	}

	for contractCount, ocrInstance := range ocrInstances {
		// Exclude the first node, which will be used as a bootstrapper
		err = ocrInstance.SetConfig(
			workerNodes,
			contracts.DefaultOffChainAggregatorConfig(len(workerNodes)),
			transmitterAddresses,
		)
		if err != nil {
			return nil, fmt.Errorf("error setting OCR config for contract '%s': %w", ocrInstance.Address(), err)
		}
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			time.Sleep(2 * time.Second)
		}
	}

	return ocrInstances, nil
}

func privateKeyToAddress(privateKey *ecdsa.PrivateKey) (common.Address, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, errors.New("error casting public key to ECDSA")
	}
	return crypto.PubkeyToAddress(*publicKeyECDSA), nil
}

func WatchNewFluxRound(
	l zerolog.Logger,
	seth *seth.Client,
	roundNumber int64,
	fluxInstance contracts.FluxAggregator,
	timeout time.Duration,
) error {
	timeoutC := time.After(timeout)
	ticker := time.NewTicker(time.Millisecond * 200)
	defer ticker.Stop()

	l.Info().Msgf("Waiting for flux round %d to be confirmed by flux aggregator", roundNumber)

	for {
		select {
		case <-timeoutC:
			return fmt.Errorf("timeout waiting for round %d to be confirmed", roundNumber)
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), seth.Cfg.Network.TxnTimeout.Duration())
			roundId, err := fluxInstance.LatestRoundID(ctx)
			if err != nil {
				cancel()
				return fmt.Errorf("getting latest round from flux instance has failed: %w", err)
			}
			cancel()
			if roundId.Cmp(big.NewInt(roundNumber)) >= 0 {
				l.Debug().Msgf("Flux instance confirmed round %d", roundNumber)
				return nil
			}
		}
	}
}

// EstimateCostForChainlinkOperations estimates the cost of running a number of operations on the Chainlink node based on estimated gas costs. It supports
// both legacy and EIP-1559 transactions.
func EstimateCostForChainlinkOperations(l zerolog.Logger, client *seth.Client, network blockchain.EVMNetwork, amountOfOperations int) (*big.Float, error) {
	bigAmountOfOperations := big.NewInt(int64(amountOfOperations))
	estimations := client.CalculateGasEstimations(client.NewDefaultGasEstimationRequest())

	gasLimit := network.GasEstimationBuffer + network.ChainlinkTransactionLimit

	var gasPriceInWei *big.Int
	if client.Cfg.Network.EIP1559DynamicFees {
		gasPriceInWei = estimations.GasFeeCap
	} else {
		gasPriceInWei = estimations.GasPrice
	}

	gasCostPerOperationWei := big.NewInt(1).Mul(big.NewInt(1).SetUint64(gasLimit), gasPriceInWei)
	gasCostPerOperationETH := conversions.WeiToEther(gasCostPerOperationWei)
	// total Wei needed for all TXs = total value for TX * number of TXs
	totalWeiForAllOperations := big.NewInt(1).Mul(gasCostPerOperationWei, bigAmountOfOperations)
	totalEthForAllOperations := conversions.WeiToEther(totalWeiForAllOperations)

	l.Debug().
		Int("Number of Operations", amountOfOperations).
		Uint64("Gas Limit per Operation", gasLimit).
		Str("Value per Operation (ETH)", gasCostPerOperationETH.String()).
		Str("Total (ETH)", totalEthForAllOperations.String()).
		Msg("Calculated ETH for Chainlink Operations")

	return totalEthForAllOperations, nil
}

// GetLatestFinalizedBlockHeader returns latest finalised block header for given network (taking into account finality tag/depth)
func GetLatestFinalizedBlockHeader(ctx context.Context, client *seth.Client, network blockchain.EVMNetwork) (*types.Header, error) {
	if network.FinalityTag {
		return client.Client.HeaderByNumber(ctx, big.NewInt(rpc.FinalizedBlockNumber.Int64()))
	}
	if network.FinalityDepth == 0 {
		return nil, fmt.Errorf("finality depth is 0 and finality tag is not enabled")
	}
	header, err := client.Client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}
	latestBlockNumber := header.Number.Uint64()
	finalizedBlockNumber := latestBlockNumber - network.FinalityDepth
	return client.Client.HeaderByNumber(ctx, big.NewInt(int64(finalizedBlockNumber)))
}

// SendLinkFundsToDeploymentAddresses sends LINK token to all addresses, but the root one, from the root address. It uses
// Multicall contract to batch all transfers in a single transaction. It also checks if the funds were transferred correctly.
// It's primary use case is to fund addresses that will be used for Upkeep registration (as that requires LINK balance) during
// Automation/Keeper test setup.
func SendLinkFundsToDeploymentAddresses(
	chainClient *seth.Client,
	concurrency,
	totalUpkeeps,
	operationsPerAddress int,
	multicallAddress common.Address,
	linkAmountPerUpkeep *big.Int,
	linkToken contracts.LinkToken,
) error {
	var generateCallData = func(receiver common.Address, amount *big.Int) ([]byte, error) {
		abi, err := link_token_interface.LinkTokenMetaData.GetAbi()
		if err != nil {
			return nil, err
		}
		data, err := abi.Pack("transfer", receiver, amount)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	toTransferToMultiCallContract := big.NewInt(0).Mul(linkAmountPerUpkeep, big.NewInt(int64(totalUpkeeps+concurrency)))
	toTransferPerClient := big.NewInt(0).Mul(linkAmountPerUpkeep, big.NewInt(int64(operationsPerAddress+1)))
	err := linkToken.Transfer(multicallAddress.Hex(), toTransferToMultiCallContract)
	if err != nil {
		return errors.Wrapf(err, "Error transferring LINK to multicall contract")
	}

	balance, err := linkToken.BalanceOf(context.Background(), multicallAddress.Hex())
	if err != nil {
		return errors.Wrapf(err, "Error getting LINK balance of multicall contract")
	}

	if balance.Cmp(toTransferToMultiCallContract) < 0 {
		return fmt.Errorf("Incorrect LINK balance of multicall contract. Expected at least: %s. Got: %s", toTransferToMultiCallContract.String(), balance.String())
	}

	// Transfer LINK to ephemeral keys
	multiCallData := make([][]byte, 0)
	for i := 1; i <= concurrency; i++ {
		data, err := generateCallData(chainClient.Addresses[i], toTransferPerClient)
		if err != nil {
			return errors.Wrapf(err, "Error generating call data for LINK transfer")
		}
		multiCallData = append(multiCallData, data)
	}

	var call []contracts.Call
	for _, d := range multiCallData {
		data := contracts.Call{Target: common.HexToAddress(linkToken.Address()), AllowFailure: false, CallData: d}
		call = append(call, data)
	}

	multiCallABI, err := abi.JSON(strings.NewReader(contracts.MultiCallABI))
	if err != nil {
		return errors.Wrapf(err, "Error getting Multicall contract ABI")
	}
	boundContract := bind.NewBoundContract(multicallAddress, multiCallABI, chainClient.Client, chainClient.Client, chainClient.Client)
	// call aggregate3 to group all msg call data and send them in a single transaction
	_, err = chainClient.Decode(boundContract.Transact(chainClient.NewTXOpts(), "aggregate3", call))
	if err != nil {
		return errors.Wrapf(err, "Error calling Multicall contract")
	}

	for i := 1; i <= concurrency; i++ {
		balance, err := linkToken.BalanceOf(context.Background(), chainClient.Addresses[i].Hex())
		if err != nil {
			return errors.Wrapf(err, "Error getting LINK balance of ephemeral key %d", i)
		}
		if balance.Cmp(toTransferPerClient) < 0 {
			return fmt.Errorf("Incorrect LINK balance after transfer. Ephemeral key %d. Expected: %s. Got: %s", i, toTransferPerClient.String(), balance.String())
		}
	}

	return nil
}

var noOpSethConfigFn = func(cfg *seth.Config) error { return nil }

type SethConfigFunction = func(*seth.Config) error

// OneEphemeralKeysLiveTestnetCheckFn checks whether there's at least one ephemeral key on a simulated network or at least one static key on a live network,
// and that there are no ephemeral keys on a live network. Root key is excluded from the check.
var OneEphemeralKeysLiveTestnetCheckFn = func(sethCfg *seth.Config) error {
	concurrency := sethCfg.GetMaxConcurrency()

	if sethCfg.IsSimulatedNetwork() {
		if concurrency < 1 {
			return fmt.Errorf(INSUFFICIENT_EPHEMERAL_KEYS, 0)
		}

		return nil
	}

	if sethCfg.EphemeralAddrs != nil && int(*sethCfg.EphemeralAddrs) > 0 {
		ephMsg := `
			Error: Ephemeral Addresses Detected on Live Network

			Ephemeral addresses are currently set for use on a live network, which is not permitted. The number of ephemeral addresses set is %d. Please make the following update to your TOML configuration file to correct this:
			'[Seth] ephemeral_addresses_number = 0'

			Additionally, ensure the following requirements are met to run this test on a live network:
			1. Use more than one private key in your network configuration.
			`

		return errors.New(ephMsg)
	}

	if concurrency < 1 {
		return fmt.Errorf(INSUFFICIENT_STATIC_KEYS, len(sethCfg.Network.PrivateKeys))
	}

	return nil
}

// OneEphemeralKeysLiveTestnetAutoFixFn checks whether there's at least one ephemeral key on a simulated network or at least one static key on a live network,
// and that there are no epehemeral keys on a live network (if ephemeral keys count is different from zero, it will disable them). Root key is excluded from the check.
var OneEphemeralKeysLiveTestnetAutoFixFn = func(sethCfg *seth.Config) error {
	concurrency := sethCfg.GetMaxConcurrency()

	if sethCfg.IsSimulatedNetwork() {
		if concurrency < 1 {
			return fmt.Errorf(INSUFFICIENT_EPHEMERAL_KEYS, 0)
		}

		return nil
	}

	if sethCfg.EphemeralAddrs != nil && int(*sethCfg.EphemeralAddrs) > 0 {
		var zero int64 = 0
		sethCfg.EphemeralAddrs = &zero
	}

	if concurrency < 1 {
		return fmt.Errorf(INSUFFICIENT_STATIC_KEYS, len(sethCfg.Network.PrivateKeys))
	}

	return nil
}

// GetChainClient returns a seth client for the given network after validating the config
func GetChainClient(config ctf_config.SethConfig, network blockchain.EVMNetwork) (*seth.Client, error) {
	return GetChainClientWithConfigFunction(config, network, noOpSethConfigFn)
}

// GetChainClientWithConfigFunction returns a seth client for the given network after validating the config and applying the config function
func GetChainClientWithConfigFunction(config ctf_config.SethConfig, network blockchain.EVMNetwork, configFn SethConfigFunction) (*seth.Client, error) {
	readSethCfg := config.GetSethConfig()
	if readSethCfg == nil {
		return nil, fmt.Errorf("Seth config not found")
	}

	sethCfg, err := seth_utils.MergeSethAndEvmNetworkConfigs(network, *readSethCfg)
	if err != nil {
		return nil, errors.Wrapf(err, "Error merging seth and evm network configs")
	}

	err = configFn(&sethCfg)
	if err != nil {
		return nil, errors.Wrapf(err, "Error applying seth config function")
	}

	err = seth_utils.ValidateSethNetworkConfig(sethCfg.Network)
	if err != nil {
		return nil, errors.Wrapf(err, "Error validating seth network config")
	}

	chainClient, err := seth.NewClientWithConfig(&sethCfg)
	if err != nil {
		return nil, errors.Wrapf(err, "Error creating seth client")
	}

	return chainClient, nil
}

// GenerateUpkeepReport generates a report of performed, successful, reverted and stale upkeeps for a given registry contract based on transaction logs. In case of test failure it can help us
// to triage the issue by providing more context.
func GenerateUpkeepReport(t *testing.T, chainClient *seth.Client, startBlock, endBlock *big.Int, instance contracts.KeeperRegistry, registryVersion ethereum.KeeperRegistryVersion) (performedUpkeeps, successfulUpkeeps, revertedUpkeeps, staleUpkeeps int, err error) {
	registryLogs := []gethtypes.Log{}
	l := logging.GetTestLogger(t)

	var (
		blockBatchSize  int64 = 100
		logs            []gethtypes.Log
		timeout         = 5 * time.Second
		addr            = common.HexToAddress(instance.Address())
		queryStartBlock = startBlock
	)

	// Gather logs from the registry in 100 block chunks to avoid read limits
	for queryStartBlock.Cmp(endBlock) < 0 {
		filterQuery := geth.FilterQuery{
			Addresses: []common.Address{addr},
			FromBlock: queryStartBlock,
			ToBlock:   big.NewInt(0).Add(queryStartBlock, big.NewInt(blockBatchSize)),
		}

		// This RPC call can possibly time out or otherwise die. Failure is not an option, keep retrying to get our stats.
		err = fmt.Errorf("initial error") // to ensure our for loop runs at least once
		for err != nil {
			ctx, cancel := context.WithTimeout(testcontext.Get(t), timeout)
			logs, err = chainClient.Client.FilterLogs(ctx, filterQuery)
			cancel()
			if err != nil {
				l.Error().
					Err(err).
					Interface("Filter Query", filterQuery).
					Str("Timeout", timeout.String()).
					Msg("Error getting logs from chain, trying again")
				timeout = time.Duration(math.Min(float64(timeout)*2, float64(2*time.Minute)))
				continue
			}
			l.Info().
				Uint64("From Block", queryStartBlock.Uint64()).
				Uint64("To Block", filterQuery.ToBlock.Uint64()).
				Int("Log Count", len(logs)).
				Str("Registry Address", addr.Hex()).
				Msg("Collected logs")
			queryStartBlock.Add(queryStartBlock, big.NewInt(blockBatchSize))
			registryLogs = append(registryLogs, logs...)
		}
	}

	var contractABI *abi.ABI
	contractABI, err = contracts.GetRegistryContractABI(registryVersion)
	if err != nil {
		return
	}

	for _, allLogs := range registryLogs {
		log := allLogs
		var eventDetails *abi.Event
		eventDetails, err = contractABI.EventByID(log.Topics[0])
		if err != nil {
			l.Error().Err(err).Str("Log Hash", log.TxHash.Hex()).Msg("Error getting event details for log, report data inaccurate")
			break
		}
		if eventDetails.Name == "UpkeepPerformed" {
			performedUpkeeps++
			var parsedLog *contracts.UpkeepPerformedLog
			parsedLog, err = instance.ParseUpkeepPerformedLog(&log)
			if err != nil {
				l.Error().Err(err).Str("Log Hash", log.TxHash.Hex()).Msg("Error parsing upkeep performed log, report data inaccurate")
				break
			}
			if !parsedLog.Success {
				revertedUpkeeps++
			} else {
				successfulUpkeeps++
			}
		} else if eventDetails.Name == "StaleUpkeepReport" {
			staleUpkeeps++
		}
	}

	return
}

func GetStalenessReportCleanupFn(t *testing.T, logger zerolog.Logger, chainClient *seth.Client, startBlock uint64, registry contracts.KeeperRegistry, registryVersion ethereum.KeeperRegistryVersion) func() {
	return func() {
		if t.Failed() {
			endBlock, err := chainClient.Client.BlockNumber(context.Background())
			require.NoError(t, err, "Failed to get end block")

			total, ok, reverted, stale, err := GenerateUpkeepReport(t, chainClient, big.NewInt(int64(startBlock)), big.NewInt(int64(endBlock)), registry, registryVersion)
			require.NoError(t, err, "Failed to get staleness data")
			if stale > 0 || reverted > 0 {
				logger.Warn().Int("Total upkeeps", total).Int("Successful upkeeps", ok).Int("Reverted Upkeeps", reverted).Int("Stale Upkeeps", stale).Msg("Staleness data")
			} else {
				logger.Info().Int("Total upkeeps", total).Int("Successful upkeeps", ok).Int("Reverted Upkeeps", reverted).Int("Stale Upkeeps", stale).Msg("Staleness data")
			}
		}
	}
}
