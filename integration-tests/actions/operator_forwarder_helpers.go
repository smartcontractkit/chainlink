package actions

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

	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

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
			}
		}
	}()
}
