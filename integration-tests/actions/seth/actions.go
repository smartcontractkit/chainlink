package actions_seth

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"

	"github.com/smartcontractkit/chainlink-testing-framework/utils/conversions"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/operator_factory"
	"github.com/smartcontractkit/seth"
)

var ContractDeploymentInterval = 200

func FundChainlinkNodes(
	logger zerolog.Logger,
	client *seth.Client,
	nodes []contracts.ChainlinkNodeWithAddress,
	fromKeyNum int,
	amount *big.Float,
) error {
	for _, cl := range nodes {
		toAddress, err := cl.PrimaryEthAddress()
		if err != nil {
			return err
		}

		if fromKeyNum > len(client.PrivateKeys) || fromKeyNum > len(client.Addresses) {
			return errors.Wrap(errors.New(seth.ErrNoKeyLoaded), fmt.Sprintf("requested key: %d", fromKeyNum))
		}
		toAddr := common.HexToAddress(toAddress)
		chainID, err := client.Client.NetworkID(context.Background())
		if err != nil {
			return errors.Wrap(err, "failed to get network ID")
		}

		ctx, cancel := context.WithTimeout(context.Background(), client.Cfg.Network.TxnTimeout.Duration())
		defer cancel()
		nonce, err := client.Client.PendingNonceAt(ctx, client.Addresses[fromKeyNum])
		if err != nil {
			return err
		}

		rawTx := &types.LegacyTx{
			Nonce:    nonce,
			To:       &toAddr,
			Value:    conversions.EtherToWei(amount),
			Gas:      uint64(client.Cfg.Network.TransferGasFee),
			GasPrice: big.NewInt(client.Cfg.Network.GasPrice),
		}
		signedTx, err := types.SignNewTx(client.PrivateKeys[fromKeyNum], types.NewEIP155Signer(chainID), rawTx)
		if err != nil {
			return errors.Wrap(err, "failed to sign tx")
		}

		ctx, cancel = context.WithTimeout(ctx, client.Cfg.Network.TxnTimeout.Duration())
		defer cancel()
		err = client.Client.SendTransaction(ctx, signedTx)
		if err != nil {
			return errors.Wrap(err, "failed to send transaction")
		}
		_, err = client.WaitMined(ctx, logger, client.Client, signedTx)
		if err != nil {
			return err
		}
		return err

	}

	return nil
}

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

// I think this is the same as the original WatchNewRound
func WatchNewRound(
	l zerolog.Logger,
	seth *seth.Client,
	roundNumber int64,
	ocrInstances []contracts.OffChainAggregatorWithRounds,
	timeout time.Duration,
) error {
	endTime := time.Now().Add(timeout)
	confirmed := make(map[string]bool)

	l.Info().Msgf("Waiting for round %d to be confirmed by all nodes", roundNumber)

	for {
		if time.Now().After(endTime) {
			return fmt.Errorf("timeout waiting for round %d to be confirmed. %d/%d nodes confirmed it", roundNumber, len(confirmed), len(ocrInstances))
		}
		for i := 0; i < len(ocrInstances); i++ {
			if confirmed[ocrInstances[i].Address()] {
				continue
			}
			ctx, cancel := context.WithTimeout(context.Background(), seth.Cfg.Network.TxnTimeout.Duration())
			roundData, err := ocrInstances[i].GetLatestRound(ctx)
			if err != nil {
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

func TeardownRemoteSuite(
	t *testing.T,
	client *seth.Client,
	namespace string,
	chainlinkNodes []*client.ChainlinkK8sClient,
	optionalTestReporter testreporters.TestReporter, // Optionally pass in a test reporter to log further metrics
	grafnaUrlProvider testreporters.GrafanaURLProvider,
) error {
	l := logging.GetTestLogger(t)
	var err error
	if err = testreporters.SendReport(t, namespace, "./", optionalTestReporter, grafnaUrlProvider); err != nil {
		l.Warn().Err(err).Msg("Error writing test report")
	}
	// Delete all jobs to stop depleting the funds
	// err = actions.DeleteAllJobs(chainlinkNodes)
	if err != nil {
		l.Warn().Msgf("Error deleting jobs %+v", err)
	}

	if err = ReturnFunds(l, client, contracts.ChainlinkK8sClientToChainlinkNodeWithKeys(chainlinkNodes)); err != nil {
		l.Error().Err(err).Str("Namespace", namespace).
			Msg("Error attempting to return funds from chainlink nodes to network's default wallet. " +
				"Environment is left running so you can try manually!")
	}
	return err
}

// StartNewRound requests a new round from the ocr contracts and waits for confirmation
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
	workerNodes []*client.ChainlinkK8sClient,
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

	return deployOCRContracts(logger, seth, numberOfContracts, linkTokenContractAddress, workerNodes, transmitterPayeesFn)
}

// DeployOCRContracts deploys and funds a certain number of offchain aggregator contracts
func DeployOCRContracts(
	logger zerolog.Logger,
	seth *seth.Client,
	numberOfContracts int,
	linkTokenContractAddress common.Address,
	workerNodes []*client.ChainlinkK8sClient,
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

	return deployOCRContracts(logger, seth, numberOfContracts, linkTokenContractAddress, workerNodes, transmitterPayeesFn)
}

func deployOCRContracts(
	logger zerolog.Logger,
	seth *seth.Client,
	numberOfContracts int,
	linkTokenContractAddress common.Address,
	workerNodes []*client.ChainlinkK8sClient,
	getTransmitterAndPayeesFn func() ([]string, []string, error),
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
	transmitters, payees, err := getTransmitterAndPayeesFn()

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
	transmitterAddresses := make([]common.Address, 0)
	for _, node := range workerNodes {
		primaryAddress, err := node.PrimaryEthAddress()
		if err != nil {
			return nil, err
		}
		transmitterAddresses = append(transmitterAddresses, common.HexToAddress(primaryAddress))
	}

	if err != nil {
		return nil, fmt.Errorf("getting node common addresses should not fail: %w", err)
	}

	var nodesAsInterface []contracts.ChainlinkNodeWithKeys = make([]contracts.ChainlinkNodeWithKeys, len(workerNodes))
	for i, node := range workerNodes {
		nodesAsInterface[i] = node // Assigning each *ChainlinkK8sClient to the interface type
	}

	for contractCount, ocrInstance := range ocrInstances {
		// Exclude the first node, which will be used as a bootstrapper
		err = ocrInstance.SetConfig(
			nodesAsInterface,
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

func ReturnFunds(log zerolog.Logger, seth *seth.Client, chainlinkNodes []contracts.ChainlinkNodeWithKeys) error {
	if seth == nil {
		return fmt.Errorf("blockchain client is nil, unable to return funds from chainlink nodes")
	}
	log.Info().Msg("Attempting to return Chainlink node funds to default network wallets")
	if seth.Cfg.IsSimulatedNetwork() {
		log.Info().Str("Network Name", seth.Cfg.Network.Name).
			Msg("Network is a simulated network. Skipping fund return.")
		return nil
	}

	panic("implement me")

	//TODO maybe implement fund return in a similar way, but with chain or responsibility
	//for handling different errors? Will be easy to add new errors to handle
	// for _, chainlinkNode := range chainlinkNodes {
	// 	fundedKeys, err := chainlinkNode.ExportEVMKeysForChain(fmt.Sprint(seth.ChainID))
	// 	if err != nil {
	// 		return err
	// 	}
	// 	for _, key := range fundedKeys {
	// 		keyToDecrypt, err := json.Marshal(key)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		// This can take up a good bit of RAM and time. When running on the remote-test-runner, this can lead to OOM
	// 		// issues. So we avoid running in parallel; slower, but safer.
	// 		decryptedKey, err := keystore.DecryptKey(keyToDecrypt, client.ChainlinkKeyPassword)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		err = blockchainClient.ReturnFunds(decryptedKey.PrivateKey)
	// 		if err != nil {
	// 			log.Error().Err(err).Str("Address", fundedKeys[0].Address).Msg("Error returning funds from Chainlink node")
	// 		}
	// 	}
	// }
}
