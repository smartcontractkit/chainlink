package contracts

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/integration-tests/wrappers"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_v1_events_mock"
	"github.com/smartcontractkit/seth"
)

// EthereumFunctionsV1EventsMock represents the basic functions v1 events mock contract
type EthereumFunctionsV1EventsMock struct {
	client     *seth.Client
	eventsMock *functions_v1_events_mock.FunctionsV1EventsMock
	address    *common.Address
}

func (f *EthereumFunctionsV1EventsMock) Address() string {
	return f.address.Hex()
}

func (f *EthereumFunctionsV1EventsMock) EmitRequestProcessed(requestId [32]byte, subscriptionId uint64, totalCostJuels *big.Int, transmitter common.Address, resultCode uint8, response []byte, errByte []byte, callbackReturnData []byte) error {
	_, err := f.client.Decode(f.eventsMock.EmitRequestProcessed(f.client.NewTXOpts(), requestId, subscriptionId, totalCostJuels, transmitter, resultCode, response, errByte, callbackReturnData))
	return err
}

func (f *EthereumFunctionsV1EventsMock) EmitRequestStart(requestId [32]byte, donId [32]byte, subscriptionId uint64, subscriptionOwner common.Address, requestingContract common.Address, requestInitiator common.Address, data []byte, dataVersion uint16, callbackGasLimit uint32, estimatedTotalCostJuels *big.Int) error {
	_, err := f.client.Decode(f.eventsMock.EmitRequestStart(f.client.NewTXOpts(), requestId, donId, subscriptionId, subscriptionOwner, requestingContract, requestInitiator, data, dataVersion, callbackGasLimit, estimatedTotalCostJuels))
	return err
}

func (f *EthereumFunctionsV1EventsMock) EmitSubscriptionCanceled(subscriptionId uint64, fundsRecipient common.Address, fundsAmount *big.Int) error {
	_, err := f.client.Decode(f.eventsMock.EmitSubscriptionCanceled(f.client.NewTXOpts(), subscriptionId, fundsRecipient, fundsAmount))
	return err
}

func (f *EthereumFunctionsV1EventsMock) EmitSubscriptionConsumerAdded(subscriptionId uint64, consumer common.Address) error {
	_, err := f.client.Decode(f.eventsMock.EmitSubscriptionConsumerAdded(f.client.NewTXOpts(), subscriptionId, consumer))
	return err
}

func (f *EthereumFunctionsV1EventsMock) EmitSubscriptionConsumerRemoved(subscriptionId uint64, consumer common.Address) error {
	_, err := f.client.Decode(f.eventsMock.EmitSubscriptionConsumerRemoved(f.client.NewTXOpts(), subscriptionId, consumer))
	return err
}

func (f *EthereumFunctionsV1EventsMock) EmitSubscriptionCreated(subscriptionId uint64, owner common.Address) error {
	_, err := f.client.Decode(f.eventsMock.EmitSubscriptionCreated(f.client.NewTXOpts(), subscriptionId, owner))
	return err
}

func (f *EthereumFunctionsV1EventsMock) EmitSubscriptionFunded(subscriptionId uint64, oldBalance *big.Int, newBalance *big.Int) error {
	_, err := f.client.Decode(f.eventsMock.EmitSubscriptionFunded(f.client.NewTXOpts(), subscriptionId, oldBalance, newBalance))
	return err
}

func (f *EthereumFunctionsV1EventsMock) EmitSubscriptionOwnerTransferred(subscriptionId uint64, from common.Address, to common.Address) error {
	_, err := f.client.Decode(f.eventsMock.EmitSubscriptionOwnerTransferred(f.client.NewTXOpts(), subscriptionId, from, to))
	return err
}

func (f *EthereumFunctionsV1EventsMock) EmitSubscriptionOwnerTransferRequested(subscriptionId uint64, from common.Address, to common.Address) error {
	_, err := f.client.Decode(f.eventsMock.EmitSubscriptionOwnerTransferRequested(f.client.NewTXOpts(), subscriptionId, from, to))
	return err
}

func (f *EthereumFunctionsV1EventsMock) EmitRequestNotProcessed(requestId [32]byte, coordinator common.Address, transmitter common.Address, resultCode uint8) error {
	_, err := f.client.Decode(f.eventsMock.EmitRequestNotProcessed(f.client.NewTXOpts(), requestId, coordinator, transmitter, resultCode))
	return err
}

func (f *EthereumFunctionsV1EventsMock) EmitContractUpdated(id [32]byte, from common.Address, to common.Address) error {
	_, err := f.client.Decode(f.eventsMock.EmitContractUpdated(f.client.NewTXOpts(), id, from, to))
	return err
}

// DeployFunctionsV1EventsMock deploys a new instance of the FunctionsV1EventsMock contract
func DeployFunctionsV1EventsMock(client *seth.Client) (FunctionsV1EventsMock, error) {
	abi, err := functions_v1_events_mock.FunctionsV1EventsMockMetaData.GetAbi()
	if err != nil {
		return &EthereumFunctionsV1EventsMock{}, fmt.Errorf("failed to get FunctionsV1EventsMock ABI: %w", err)
	}
	client.ContractStore.AddABI("FunctionsV1EventsMock", *abi)
	client.ContractStore.AddBIN("FunctionsV1EventsMock", common.FromHex(functions_v1_events_mock.FunctionsV1EventsMockMetaData.Bin))

	data, err := client.DeployContract(client.NewTXOpts(), "FunctionsV1EventsMock", *abi, common.FromHex(functions_v1_events_mock.FunctionsV1EventsMockMetaData.Bin))

	if err != nil {
		return &EthereumFunctionsV1EventsMock{}, fmt.Errorf("FunctionsV1EventsMock instance deployment have failed: %w", err)
	}

	instance, err := functions_v1_events_mock.NewFunctionsV1EventsMock(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumFunctionsV1EventsMock{}, fmt.Errorf("failed to instantiate FunctionsV1EventsMock instance: %w", err)
	}

	return &EthereumFunctionsV1EventsMock{
		client:     client,
		eventsMock: instance,
		address:    &data.Address,
	}, nil
}
