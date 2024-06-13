package actions

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/operator_factory"

	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

func ProcessNewEvent(
	t *testing.T,
	operatorCreated chan *operator_factory.OperatorFactoryOperatorCreated,
	authorizedForwarderCreated chan *operator_factory.OperatorFactoryAuthorizedForwarderCreated,
	event *types.Log,
	eventDetails *abi.Event,
	operatorFactoryInstance contracts.OperatorFactory,
	chainClient blockchain.EVMClient,
) {
	l := logging.GetTestLogger(t)
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
