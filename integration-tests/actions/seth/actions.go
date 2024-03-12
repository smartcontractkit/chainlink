package actions_seth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/seth"
	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/operator_factory"

	"github.com/smartcontractkit/chainlink-testing-framework/utils/conversions"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
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
	for _, cl := range nodes {
		toAddress, err := cl.PrimaryEthAddress()
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
	GasLimit   *uint64
}

// TODO: move to CTF?
// SendFunds sends native token amount (expressed in human-scale) from address controlled by private key
// to given address. If no gas limit is set, then network's default will be used.
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

	gasLimit := uint64(client.Cfg.Network.GasLimit)
	if payload.GasLimit != nil {
		gasLimit = *payload.GasLimit
	}

	var signedTx *types.Transaction

	if client.Cfg.Network.EIP1559DynamicFees {
		rawTx := &types.DynamicFeeTx{
			Nonce:     nonce,
			To:        &payload.ToAddress,
			Value:     payload.Amount,
			Gas:       gasLimit,
			GasFeeCap: big.NewInt(client.Cfg.Network.GasFeeCap),
			GasTipCap: big.NewInt(client.Cfg.Network.GasTipCap),
		}
		signedTx, err = types.SignNewTx(payload.PrivateKey, types.NewLondonSigner(big.NewInt(client.ChainID)), rawTx)
	} else {
		rawTx := &types.LegacyTx{
			Nonce:    nonce,
			To:       &payload.ToAddress,
			Value:    payload.Amount,
			Gas:      gasLimit,
			GasPrice: big.NewInt(client.Cfg.Network.GasPrice),
		}
		signedTx, err = types.SignNewTx(payload.PrivateKey, types.NewEIP155Signer(big.NewInt(client.ChainID)), rawTx)
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to sign tx")
	}

	ctx, cancel = context.WithTimeout(ctx, client.Cfg.Network.TxnTimeout.Duration())
	defer cancel()
	err = client.Client.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send transaction")
	}

	logger.Debug().
		Str("From", fromAddress.Hex()).
		Str("To", payload.ToAddress.Hex()).
		Str("TxHash", signedTx.Hash().String()).
		Str("Amount", conversions.WeiToEther(payload.Amount).String()).
		Uint64("Nonce", nonce).
		Uint64("Gas Limit", gasLimit).
		Int64("Gas Price", client.Cfg.Network.GasPrice).
		Int64("Gas Fee Cap", client.Cfg.Network.GasFeeCap).
		Int64("Gas Tip Cap", client.Cfg.Network.GasTipCap).
		Bool("Dynamic fees", client.Cfg.Network.EIP1559DynamicFees).
		Msg("Sent funds")

	return client.WaitMined(ctx, logger, client.Client, signedTx)
}

// DeployForwarderContracts first deploys Operator Factory and then uses it to deploy given number of
// operator and forwarder pairs. It waits for each transaction to be mined and then extracts operator and
// forwarder addresses from emitted events.
func DeployForwarderContracts(
	t *testing.T,
	seth *seth.Client,
	linkTokenData seth.DeploymentData,
	numberOfOperatorForwarderPairs int,
) (operators []common.Address, authorizedForwarders []common.Address, operatorFactoryInstance contracts.OperatorFactory) {
	instance, err := contracts.DeployEthereumOperatorFactory(seth, linkTokenData.Address)
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

// WatchNewRound watches for a new OCR round, similarly to StartNewRound, but it does not explicitly request a new
// round from the contract, as this can cause some odd behavior in some cases. It announces success if latest round
// is >= roundNumber.
func WatchNewRound(
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
	node *client.ChainlinkK8sClient,
) {
	l := logging.GetTestLogger(t)
	chainID := big.NewInt(seth.ChainID)
	_, _, err := node.TrackForwarder(chainID, authorizedForwarder)
	require.NoError(t, err, "Forwarder track should be created")
	l.Info().Str("NodeURL", node.Config.URL).
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

	if err = ReturnFunds(l, client, contracts.ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(chainlinkNodes)); err != nil {
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
