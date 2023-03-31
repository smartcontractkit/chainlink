package actions

import (
	"context"
	"math/big"
	"testing"

	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/operator_factory"
)

func DeployForwarderContracts(
	t *testing.T,
	contractDeployer contracts.ContractDeployer,
	linkToken contracts.LinkToken,
	chainClient blockchain.EVMClient,
	numberOfOperatorForwarderPairs int,
) (operators []common.Address, authorizedForwarders []common.Address, operatorFactoryInstance contracts.OperatorFactory) {
	operatorFactoryInstance, err := contractDeployer.DeployOperatorFactory(linkToken.Address())
	require.NoError(t, err, "Deploying OperatorFactory Contract shouldn't fail")
	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Failed waiting for deployment of flux aggregator contract")

	operatorCreated := make(chan *operator_factory.OperatorFactoryOperatorCreated)
	authorizedForwarderCreated := make(chan *operator_factory.OperatorFactoryAuthorizedForwarderCreated)
	for i := 0; i < numberOfOperatorForwarderPairs; i++ {
		SubscribeOperatorFactoryEvents(t, authorizedForwarderCreated, operatorCreated, chainClient, operatorFactoryInstance)
		_, err = operatorFactoryInstance.DeployNewOperatorAndForwarder()
		require.NoError(t, err, "Deploying new operator with proposed ownership with forwarder shouldn't fail")
		err = chainClient.WaitForEvents()
		require.NoError(t, err, "Waiting for events in nodes shouldn't fail")
		eventDataAuthorizedForwarder, eventDataOperatorCreated := <-authorizedForwarderCreated, <-operatorCreated
		operator, authorizedForwarder := eventDataOperatorCreated.Operator, eventDataAuthorizedForwarder.Forwarder
		operators = append(operators, operator)
		authorizedForwarders = append(authorizedForwarders, authorizedForwarder)
	}
	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error waiting for events")
	return operators, authorizedForwarders, operatorFactoryInstance
}

func AcceptAuthorizedReceiversOperator(
	t *testing.T,
	operator common.Address,
	authorizedForwarder common.Address,
	nodeAddresses []common.Address,
	chainClient blockchain.EVMClient,
	contractLoader contracts.ContractLoader,
) {
	operatorInstance, err := contractLoader.LoadOperatorContract(operator)
	require.NoError(t, err, "Loading operator contract shouldn't fail")
	forwarderInstance, err := contractLoader.LoadAuthorizedForwarder(authorizedForwarder)
	require.NoError(t, err, "Loading authorized forwarder contract shouldn't fail")

	err = operatorInstance.AcceptAuthorizedReceivers([]common.Address{authorizedForwarder}, nodeAddresses)
	require.NoError(t, err, "Accepting authorized receivers shouldn't fail")
	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Waiting for events in nodes shouldn't fail")

	senders, err := forwarderInstance.GetAuthorizedSenders(context.Background())
	require.NoError(t, err, "Getting authorized senders shouldn't fail")
	var nodesAddrs []string
	for _, o := range nodeAddresses {
		nodesAddrs = append(nodesAddrs, o.Hex())
	}
	require.Equal(t, nodesAddrs, senders, "Senders addresses should match node addresses")

	owner, err := forwarderInstance.Owner(context.Background())
	require.NoError(t, err, "Getting authorized forwarder owner shouldn't fail")
	require.Equal(t, operator.Hex(), owner, "Forwarder owner should match operator")
}

func ProcessNewEvent(
	t *testing.T,
	eventSub geth.Subscription,
	operatorCreated chan *operator_factory.OperatorFactoryOperatorCreated,
	authorizedForwarderCreated chan *operator_factory.OperatorFactoryAuthorizedForwarderCreated,
	event *types.Log,
	eventDetails *abi.Event,
	operatorFactoryInstance contracts.OperatorFactory,
	contractABI *abi.ABI,
	chainClient blockchain.EVMClient,
) {
	l := utils.GetTestLogger(t)
	errorChan := make(chan error)
	eventConfirmed := make(chan bool)
	err := chainClient.ProcessEvent(eventDetails.Name, event, eventConfirmed, errorChan)
	if err != nil {
		l.Error().Err(err).Str("Hash", event.TxHash.Hex()).Str("Event", eventDetails.Name).Msg("Error trying to process event")
		return
	}
	l.Debug().
		Str("Event", eventDetails.Name).
		Str("Address", event.Address.Hex()).
		Str("Hash", event.TxHash.Hex()).
		Msg("Attempting to Confirm Event")
	for {
		select {
		case err := <-errorChan:
			l.Error().Err(err).Msg("Error while confirming event")
			return
		case confirmed := <-eventConfirmed:
			if confirmed {
				if eventDetails.Name == "AuthorizedForwarderCreated" { // AuthorizedForwarderCreated event to authorizedForwarderCreated channel to handle in main loop
					eventData, err := operatorFactoryInstance.ParseAuthorizedForwarderCreated(*event)
					require.NoError(t, err, "Parsing OperatorFactoryAuthorizedForwarderCreated event log in "+
						"OperatorFactory instance shouldn't fail")
					authorizedForwarderCreated <- eventData
				}
				if eventDetails.Name == "OperatorCreated" { // OperatorCreated event to operatorCreated channel to handle in main loop
					eventData, err := operatorFactoryInstance.ParseOperatorCreated(*event)
					require.NoError(t, err, "Parsing OperatorFactoryAuthorizedForwarderCreated event log in "+
						"OperatorFactory instance shouldn't fail")
					operatorCreated <- eventData
				}
			}
			return
		}
	}
}

// SubscribeOperatorFactoryEvents subscribes to the event log for authorizedForwarderCreated and operatorCreated events
// from OperatorFactory contract
func SubscribeOperatorFactoryEvents(
	t *testing.T,
	authorizedForwarderCreated chan *operator_factory.OperatorFactoryAuthorizedForwarderCreated,
	operatorCreated chan *operator_factory.OperatorFactoryOperatorCreated,
	chainClient blockchain.EVMClient,
	operatorFactoryInstance contracts.OperatorFactory,
) {
	l := utils.GetTestLogger(t)
	contractABI, err := operator_factory.OperatorFactoryMetaData.GetAbi()
	require.NoError(t, err, "Getting contract abi for OperatorFactory shouldn't fail")
	latestBlockNum, err := chainClient.LatestBlockNumber(context.Background())
	require.NoError(t, err, "Subscribing to contract event log for OperatorFactory instance shouldn't fail")
	query := geth.FilterQuery{
		FromBlock: big.NewInt(0).SetUint64(latestBlockNum),
		Addresses: []common.Address{common.HexToAddress(operatorFactoryInstance.Address())},
	}

	eventLogs := make(chan types.Log)
	sub, err := chainClient.SubscribeFilterLogs(context.Background(), query, eventLogs)
	require.NoError(t, err, "Subscribing to contract event log for OperatorFactory instance shouldn't fail")
	go func() {
		defer sub.Unsubscribe()
		remainingExpectedEvents := 2
		for {
			select {
			case err := <-sub.Err():
				l.Error().Err(err).Msg("Error while watching for new contract events. Retrying Subscription")
				sub.Unsubscribe()

				sub, err = chainClient.SubscribeFilterLogs(context.Background(), query, eventLogs)
				require.NoError(t, err, "Subscribing to contract event log for OperatorFactory instance shouldn't fail")
			case vLog := <-eventLogs:
				eventDetails, err := contractABI.EventByID(vLog.Topics[0])
				require.NoError(t, err, "Getting event details for OperatorFactory instance shouldn't fail")
				go ProcessNewEvent(
					t, sub, operatorCreated, authorizedForwarderCreated, &vLog,
					eventDetails, operatorFactoryInstance, contractABI, chainClient,
				)
				if eventDetails.Name == "AuthorizedForwarderCreated" || eventDetails.Name == "OperatorCreated" {
					remainingExpectedEvents--
					if remainingExpectedEvents <= 0 {
						return
					}
				}
			}
		}
	}()
}

func TrackForwarder(
	t *testing.T,
	chainClient blockchain.EVMClient,
	authorizedForwarder common.Address,
	node *client.Chainlink,
) {
	l := utils.GetTestLogger(t)
	chainID := chainClient.GetChainID()
	_, _, err := node.TrackForwarder(chainID, authorizedForwarder)
	require.NoError(t, err, "Forwarder track should be created")
	l.Info().Str("NodeURL", node.Config.URL).
		Str("ForwarderAddress", authorizedForwarder.Hex()).
		Str("ChaindID", chainID.String()).
		Msg("Forwarder tracked")
}
