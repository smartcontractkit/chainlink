package actions

//revive:disable:dot-imports
import (
	"context"
	"math/big"

	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"

	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/operator_factory"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

func DeployForwarderContracts(
	contractDeployer contracts.ContractDeployer,
	linkToken contracts.LinkToken,
	chainClient blockchain.EVMClient,
	numberOfOperatorForwarderPairs int,
) (operators []common.Address, authorizedForwarders []common.Address, operatorFactoryInstance contracts.OperatorFactory) {
	By("Deploying OperatorFactory contract")
	operatorFactoryInstance, err := contractDeployer.DeployOperatorFactory(linkToken.Address())
	Expect(err).ShouldNot(HaveOccurred(), "Deploying OperatorFactory Contract shouldn't fail")
	err = chainClient.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for deployment of flux aggregator contract")

	operatorCreated := make(chan *operator_factory.OperatorFactoryOperatorCreated)
	authorizedForwarderCreated := make(chan *operator_factory.OperatorFactoryAuthorizedForwarderCreated)
	for i := 0; i < numberOfOperatorForwarderPairs; i++ {
		By("Subscribe to Operator factory Events")
		SubscribeOperatorFactoryEvents(authorizedForwarderCreated, operatorCreated, chainClient, operatorFactoryInstance)
		By("Create new operator and forwarder")
		_, err = operatorFactoryInstance.DeployNewOperatorAndForwarder()
		Expect(err).ShouldNot(HaveOccurred(), "Deploying new operator with proposed ownership with forwarder shouldn't fail")
		err = chainClient.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred(), "Waiting for events in nodes shouldn't fail")
		eventDataAuthorizedForwarder, eventDataOperatorCreated := <-authorizedForwarderCreated, <-operatorCreated
		operator, authorizedForwarder := eventDataOperatorCreated.Operator, eventDataAuthorizedForwarder.Forwarder
		operators = append(operators, operator)
		authorizedForwarders = append(authorizedForwarders, authorizedForwarder)
	}
	err = chainClient.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Error waiting for events")
	return operators, authorizedForwarders, operatorFactoryInstance
}

func AcceptAuthorizedReceiversOperator(
	operator common.Address,
	authorizedForwarder common.Address,
	nodeAddresses []common.Address,
	chainClient blockchain.EVMClient,
	contractLoader contracts.ContractLoader,
) {
	operatorInstance, err := contractLoader.LoadOperatorContract(operator)
	Expect(err).ShouldNot(HaveOccurred(), "Loading operator contract shouldn't fail")
	forwarderInstance, err := contractLoader.LoadAuthorizedForwarder(authorizedForwarder)
	Expect(err).ShouldNot(HaveOccurred(), "Loading authorized forwarder contract shouldn't fail")

	By("Accept authorized receivers")
	err = operatorInstance.AcceptAuthorizedReceivers([]common.Address{authorizedForwarder}, nodeAddresses)
	Expect(err).ShouldNot(HaveOccurred(), "Accepting authorized receivers shouldn't fail")
	err = chainClient.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Waiting for events in nodes shouldn't fail")

	By("Verify authorized senders on forwarder")
	senders, err := forwarderInstance.GetAuthorizedSenders(context.Background())
	Expect(err).ShouldNot(HaveOccurred(), "Getting authorized senders shouldn't fail")
	var nodesAddrs []string
	for _, o := range nodeAddresses {
		nodesAddrs = append(nodesAddrs, o.Hex())
	}
	Expect(senders).Should(Equal(nodesAddrs), "Senders addresses should match node addresses")

	By("Verify forwarder Owner")
	owner, err := forwarderInstance.Owner(context.Background())
	Expect(err).ShouldNot(HaveOccurred(), "Getting authorized forwarder owner shouldn't fail")
	Expect(owner).Should(Equal(operator.Hex()), "Forwarder owner should match operator")
}

func ProcessNewEvent(
	eventSub geth.Subscription,
	operatorCreated chan *operator_factory.OperatorFactoryOperatorCreated,
	authorizedForwarderCreated chan *operator_factory.OperatorFactoryAuthorizedForwarderCreated,
	event *types.Log,
	eventDetails *abi.Event,
	operatorFactoryInstance contracts.OperatorFactory,
	contractABI *abi.ABI,
	chainClient blockchain.EVMClient,
) {
	defer GinkgoRecover()

	errorChan := make(chan error)
	eventConfirmed := make(chan bool)
	err := chainClient.ProcessEvent(eventDetails.Name, event, eventConfirmed, errorChan)
	if err != nil {
		log.Error().Err(err).Str("Hash", event.TxHash.Hex()).Str("Event", eventDetails.Name).Msg("Error trying to process event")
		return
	}
	log.Debug().
		Str("Event", eventDetails.Name).
		Str("Address", event.Address.Hex()).
		Str("Hash", event.TxHash.Hex()).
		Msg("Attempting to Confirm Event")
	for {
		select {
		case err := <-errorChan:
			log.Error().Err(err).Msg("Error while confirming event")
			return
		case confirmed := <-eventConfirmed:
			if confirmed {
				if eventDetails.Name == "AuthorizedForwarderCreated" { // AuthorizedForwarderCreated event to authorizedForwarderCreated channel to handle in main loop
					eventData, err := operatorFactoryInstance.ParseAuthorizedForwarderCreated(*event)
					Expect(err).ShouldNot(HaveOccurred(), "Parsing OperatorFactoryAuthorizedForwarderCreated event log in "+
						"OperatorFactory instance shouldn't fail")
					authorizedForwarderCreated <- eventData
				}
				if eventDetails.Name == "OperatorCreated" { // OperatorCreated event to operatorCreated channel to handle in main loop
					eventData, err := operatorFactoryInstance.ParseOperatorCreated(*event)
					Expect(err).ShouldNot(HaveOccurred(), "Parsing OperatorFactoryAuthorizedForwarderCreated event log in "+
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
	authorizedForwarderCreated chan *operator_factory.OperatorFactoryAuthorizedForwarderCreated,
	operatorCreated chan *operator_factory.OperatorFactoryOperatorCreated,
	chainClient blockchain.EVMClient,
	operatorFactoryInstance contracts.OperatorFactory,
) {
	contractABI, err := operator_factory.OperatorFactoryMetaData.GetAbi()
	Expect(err).ShouldNot(HaveOccurred(), "Getting contract abi for OperatorFactory shouldn't fail")
	latestBlockNum, err := chainClient.LatestBlockNumber(context.Background())
	Expect(err).ShouldNot(HaveOccurred(), "Subscribing to contract event log for OperatorFactory instance shouldn't fail")
	query := geth.FilterQuery{
		FromBlock: big.NewInt(0).SetUint64(latestBlockNum),
		Addresses: []common.Address{common.HexToAddress(operatorFactoryInstance.Address())},
	}

	eventLogs := make(chan types.Log)
	sub, err := chainClient.SubscribeFilterLogs(context.Background(), query, eventLogs)
	Expect(err).ShouldNot(HaveOccurred(), "Subscribing to contract event log for OperatorFactory instance shouldn't fail")
	go func() {
		defer GinkgoRecover()
		defer sub.Unsubscribe()
		remainingExpectedEvents := 2
		for {
			select {
			case err := <-sub.Err():
				log.Error().Err(err).Msg("Error while watching for new contract events. Retrying Subscription")
				sub.Unsubscribe()

				sub, err = chainClient.SubscribeFilterLogs(context.Background(), query, eventLogs)
				Expect(err).ShouldNot(HaveOccurred(), "Subscribing to contract event log for OperatorFactory instance shouldn't fail")
			case vLog := <-eventLogs:
				eventDetails, err := contractABI.EventByID(vLog.Topics[0])
				Expect(err).ShouldNot(HaveOccurred(), "Getting event details for OperatorFactory instance shouldn't fail")
				go ProcessNewEvent(sub, operatorCreated, authorizedForwarderCreated, &vLog, eventDetails, operatorFactoryInstance, contractABI, chainClient)
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

func TrackForwarder(chainClient blockchain.EVMClient, authorizedForwarder common.Address, node *client.Chainlink) {
	chainID := chainClient.GetChainID()
	_, _, err := node.TrackForwarder(chainID, authorizedForwarder)
	Expect(err).ShouldNot(HaveOccurred(), "Forwarder track should be created")
	log.Info().Str("NodeURL", node.Config.URL).
		Str("ForwarderAddress", authorizedForwarder.Hex()).
		Str("ChaindID", chainID.String()).
		Msg("Forwarder tracked")
}
