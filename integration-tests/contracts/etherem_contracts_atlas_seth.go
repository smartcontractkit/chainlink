package contracts

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/seth"

	eth_contracts "github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/integration-tests/wrappers"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_v1_events_mock"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/gas_wrapper_mock"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registrar_wrapper1_2_mock"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_1_mock"
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

// EthereumKeeperRegistry11Mock represents the basic keeper registry 1.1 mock contract
type EthereumKeeperRegistry11Mock struct {
	client       *seth.Client
	registryMock *keeper_registry_wrapper1_1_mock.KeeperRegistryMock
	address      *common.Address
}

func (f *EthereumKeeperRegistry11Mock) Address() string {
	return f.address.Hex()
}

func (f *EthereumKeeperRegistry11Mock) EmitUpkeepPerformed(id *big.Int, success bool, from common.Address, payment *big.Int, performData []byte) error {
	_, err := f.client.Decode(f.registryMock.EmitUpkeepPerformed(f.client.NewTXOpts(), id, success, from, payment, performData))
	return err
}

func (f *EthereumKeeperRegistry11Mock) EmitUpkeepCanceled(id *big.Int, atBlockHeight uint64) error {
	_, err := f.client.Decode(f.registryMock.EmitUpkeepCanceled(f.client.NewTXOpts(), id, atBlockHeight))
	return err
}

func (f *EthereumKeeperRegistry11Mock) EmitFundsWithdrawn(id *big.Int, amount *big.Int, to common.Address) error {
	_, err := f.client.Decode(f.registryMock.EmitFundsWithdrawn(f.client.NewTXOpts(), id, amount, to))
	return err
}

func (f *EthereumKeeperRegistry11Mock) EmitKeepersUpdated(keepers []common.Address, payees []common.Address) error {
	_, err := f.client.Decode(f.registryMock.EmitKeepersUpdated(f.client.NewTXOpts(), keepers, payees))
	return err
}

func (f *EthereumKeeperRegistry11Mock) EmitUpkeepRegistered(id *big.Int, executeGas uint32, admin common.Address) error {
	_, err := f.client.Decode(f.registryMock.EmitUpkeepRegistered(f.client.NewTXOpts(), id, executeGas, admin))
	return err
}

func (f *EthereumKeeperRegistry11Mock) EmitFundsAdded(id *big.Int, from common.Address, amount *big.Int) error {
	_, err := f.client.Decode(f.registryMock.EmitFundsAdded(f.client.NewTXOpts(), id, from, amount))
	return err
}

func (f *EthereumKeeperRegistry11Mock) SetUpkeepCount(upkeepCount *big.Int) error {
	_, err := f.client.Decode(f.registryMock.SetUpkeepCount(f.client.NewTXOpts(), upkeepCount))
	return err
}

func (f *EthereumKeeperRegistry11Mock) SetCanceledUpkeepList(canceledUpkeepList []*big.Int) error {
	_, err := f.client.Decode(f.registryMock.SetCanceledUpkeepList(f.client.NewTXOpts(), canceledUpkeepList))
	return err
}

func (f *EthereumKeeperRegistry11Mock) SetKeeperList(keepers []common.Address) error {
	_, err := f.client.Decode(f.registryMock.SetKeeperList(f.client.NewTXOpts(), keepers))
	return err
}

func (f *EthereumKeeperRegistry11Mock) SetConfig(paymentPremiumPPB uint32, flatFeeMicroLink uint32, blockCountPerTurn *big.Int, checkGasLimit uint32, stalenessSeconds *big.Int, gasCeilingMultiplier uint16, fallbackGasPrice *big.Int, fallbackLinkPrice *big.Int) error {
	_, err := f.client.Decode(f.registryMock.SetConfig(f.client.NewTXOpts(), paymentPremiumPPB, flatFeeMicroLink, blockCountPerTurn, checkGasLimit, stalenessSeconds, gasCeilingMultiplier, fallbackGasPrice, fallbackLinkPrice))
	return err
}

func (f *EthereumKeeperRegistry11Mock) SetUpkeep(id *big.Int, target common.Address, executeGas uint32, balance *big.Int, admin common.Address, maxValidBlocknumber uint64, lastKeeper common.Address, checkData []byte) error {
	_, err := f.client.Decode(f.registryMock.SetUpkeep(f.client.NewTXOpts(), id, target, executeGas, balance, admin, maxValidBlocknumber, lastKeeper, checkData))
	return err
}

func (f *EthereumKeeperRegistry11Mock) SetMinBalance(id *big.Int, minBalance *big.Int) error {
	_, err := f.client.Decode(f.registryMock.SetMinBalance(f.client.NewTXOpts(), id, minBalance))
	return err
}

func (f *EthereumKeeperRegistry11Mock) SetCheckUpkeepData(id *big.Int, performData []byte, maxLinkPayment *big.Int, gasLimit *big.Int, adjustedGasWei *big.Int, linkEth *big.Int) error {
	_, err := f.client.Decode(f.registryMock.SetCheckUpkeepData(f.client.NewTXOpts(), id, performData, maxLinkPayment, gasLimit, adjustedGasWei, linkEth))
	return err
}

func (f *EthereumKeeperRegistry11Mock) SetPerformUpkeepSuccess(id *big.Int, success bool) error {
	_, err := f.client.Decode(f.registryMock.SetPerformUpkeepSuccess(f.client.NewTXOpts(), id, success))
	return err
}

func DeployKeeperRegistry11Mock(client *seth.Client) (KeeperRegistry11Mock, error) {
	abi, err := keeper_registry_wrapper1_1_mock.KeeperRegistryMockMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperRegistry11Mock{}, fmt.Errorf("failed to get KeeperRegistry11Mock ABI: %w", err)
	}
	client.ContractStore.AddABI("KeeperRegistry11Mock", *abi)
	client.ContractStore.AddBIN("KeeperRegistry11Mock", common.FromHex(keeper_registry_wrapper1_1_mock.KeeperRegistryMockMetaData.Bin))

	data, err := client.DeployContract(client.NewTXOpts(), "KeeperRegistry11Mock", *abi, common.FromHex(keeper_registry_wrapper1_1_mock.KeeperRegistryMockMetaData.Bin))

	if err != nil {
		return &EthereumKeeperRegistry11Mock{}, fmt.Errorf("KeeperRegistry11Mock instance deployment have failed: %w", err)
	}

	instance, err := keeper_registry_wrapper1_1_mock.NewKeeperRegistryMock(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumKeeperRegistry11Mock{}, fmt.Errorf("failed to instantiate KeeperRegistry11Mock instance: %w", err)
	}

	return &EthereumKeeperRegistry11Mock{
		client:       client,
		registryMock: instance,
		address:      &data.Address,
	}, nil
}

// EthereumKeeperRegistrar12Mock represents the basic keeper registrar 1.2 mock contract
type EthereumKeeperRegistrar12Mock struct {
	client        *seth.Client
	registrarMock *keeper_registrar_wrapper1_2_mock.KeeperRegistrarMock
	address       *common.Address
}

func (f *EthereumKeeperRegistrar12Mock) Address() string {
	return f.address.Hex()
}

func (f *EthereumKeeperRegistrar12Mock) EmitRegistrationRequested(hash [32]byte, name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, amount *big.Int, source uint8) error {
	_, err := f.client.Decode(f.registrarMock.EmitRegistrationRequested(f.client.NewTXOpts(), hash, name, encryptedEmail, upkeepContract, gasLimit, adminAddress, checkData, amount, source))
	return err
}

func (f *EthereumKeeperRegistrar12Mock) EmitRegistrationApproved(hash [32]byte, displayName string, upkeepId *big.Int) error {
	_, err := f.client.Decode(f.registrarMock.EmitRegistrationApproved(f.client.NewTXOpts(), hash, displayName, upkeepId))
	return err
}

func (f *EthereumKeeperRegistrar12Mock) SetRegistrationConfig(autoApproveConfigType uint8, autoApproveMaxAllowed uint32, approvedCount uint32, keeperRegistry common.Address, minLINKJuels *big.Int) error {
	_, err := f.client.Decode(f.registrarMock.SetRegistrationConfig(f.client.NewTXOpts(), autoApproveConfigType, autoApproveMaxAllowed, approvedCount, keeperRegistry, minLINKJuels))
	return err
}

func DeployKeeperRegistrar12Mock(client *seth.Client) (KeeperRegistrar12Mock, error) {
	abi, err := keeper_registrar_wrapper1_2_mock.KeeperRegistrarMockMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperRegistrar12Mock{}, fmt.Errorf("failed to get KeeperRegistrar12Mock ABI: %w", err)
	}
	client.ContractStore.AddABI("KeeperRegistrar12Mock", *abi)
	client.ContractStore.AddBIN("KeeperRegistrar12Mock", common.FromHex(keeper_registrar_wrapper1_2_mock.KeeperRegistrarMockMetaData.Bin))

	data, err := client.DeployContract(client.NewTXOpts(), "KeeperRegistrar12Mock", *abi, common.FromHex(keeper_registrar_wrapper1_2_mock.KeeperRegistrarMockMetaData.Bin))

	if err != nil {
		return &EthereumKeeperRegistrar12Mock{}, fmt.Errorf("KeeperRegistrar12Mock instance deployment have failed: %w", err)
	}

	instance, err := keeper_registrar_wrapper1_2_mock.NewKeeperRegistrarMock(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumKeeperRegistrar12Mock{}, fmt.Errorf("failed to instantiate KeeperRegistrar12Mock instance: %w", err)
	}

	return &EthereumKeeperRegistrar12Mock{
		client:        client,
		registrarMock: instance,
		address:       &data.Address,
	}, nil
}

// EthereumKeeperGasWrapperMock represents the basic keeper gas wrapper mock contract
type EthereumKeeperGasWrapperMock struct {
	client         *seth.Client
	gasWrapperMock *gas_wrapper_mock.KeeperRegistryCheckUpkeepGasUsageWrapperMock
	address        *common.Address
}

func (f *EthereumKeeperGasWrapperMock) Address() string {
	return f.address.Hex()
}

func (f *EthereumKeeperGasWrapperMock) SetMeasureCheckGasResult(result bool, payload []byte, gas *big.Int) error {
	_, err := f.client.Decode(f.gasWrapperMock.SetMeasureCheckGasResult(f.client.NewTXOpts(), result, payload, gas))
	return err
}

func DeployKeeperGasWrapperMock(client *seth.Client) (KeeperGasWrapperMock, error) {
	abi, err := gas_wrapper_mock.KeeperRegistryCheckUpkeepGasUsageWrapperMockMetaData.GetAbi()
	if err != nil {
		return &EthereumKeeperGasWrapperMock{}, fmt.Errorf("failed to get KeeperGasWrapperMock ABI: %w", err)
	}
	client.ContractStore.AddABI("KeeperGasWrapperMock", *abi)
	client.ContractStore.AddBIN("KeeperGasWrapperMock", common.FromHex(gas_wrapper_mock.KeeperRegistryCheckUpkeepGasUsageWrapperMockMetaData.Bin))

	data, err := client.DeployContract(client.NewTXOpts(), "KeeperGasWrapperMock", *abi, common.FromHex(gas_wrapper_mock.KeeperRegistryCheckUpkeepGasUsageWrapperMockMetaData.Bin))

	if err != nil {
		return &EthereumKeeperGasWrapperMock{}, fmt.Errorf("KeeperGasWrapperMock instance deployment have failed: %w", err)
	}

	instance, err := gas_wrapper_mock.NewKeeperRegistryCheckUpkeepGasUsageWrapperMock(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumKeeperGasWrapperMock{}, fmt.Errorf("failed to instantiate KeeperGasWrapperMock instance: %w", err)
	}

	return &EthereumKeeperGasWrapperMock{
		client:         client,
		gasWrapperMock: instance,
		address:        &data.Address,
	}, nil
}

// EthereumStakingEventsMock represents the basic events mock contract
type EthereumStakingEventsMock struct {
	client     *seth.Client
	eventsMock *eth_contracts.StakingEventsMock
	address    *common.Address
}

func (f *EthereumStakingEventsMock) Address() string {
	return f.address.Hex()
}

func (f *EthereumStakingEventsMock) MaxCommunityStakeAmountIncreased(maxStakeAmount *big.Int) error {
	_, err := f.client.Decode(f.eventsMock.EmitMaxCommunityStakeAmountIncreased(f.client.NewTXOpts(), maxStakeAmount))
	return err
}

func (f *EthereumStakingEventsMock) PoolSizeIncreased(maxPoolSize *big.Int) error {
	_, err := f.client.Decode(f.eventsMock.EmitPoolSizeIncreased(f.client.NewTXOpts(), maxPoolSize))
	return err
}

func (f *EthereumStakingEventsMock) MaxOperatorStakeAmountIncreased(maxStakeAmount *big.Int) error {
	_, err := f.client.Decode(f.eventsMock.EmitMaxOperatorStakeAmountIncreased(f.client.NewTXOpts(), maxStakeAmount))
	return err
}

func (f *EthereumStakingEventsMock) RewardInitialized(rate *big.Int, available *big.Int, startTimestamp *big.Int, endTimestamp *big.Int) error {
	_, err := f.client.Decode(f.eventsMock.EmitRewardInitialized(f.client.NewTXOpts(), rate, available, startTimestamp, endTimestamp))
	return err
}

func (f *EthereumStakingEventsMock) AlertRaised(alerter common.Address, roundId *big.Int, rewardAmount *big.Int) error {
	_, err := f.client.Decode(f.eventsMock.EmitAlertRaised(f.client.NewTXOpts(), alerter, roundId, rewardAmount))
	return err
}

func (f *EthereumStakingEventsMock) Staked(staker common.Address, newStake *big.Int, totalStake *big.Int) error {
	_, err := f.client.Decode(f.eventsMock.EmitStaked(f.client.NewTXOpts(), staker, newStake, totalStake))
	return err
}

func (f *EthereumStakingEventsMock) OperatorAdded(operator common.Address) error {
	_, err := f.client.Decode(f.eventsMock.EmitOperatorAdded(f.client.NewTXOpts(), operator))
	return err
}

func (f *EthereumStakingEventsMock) OperatorRemoved(operator common.Address, amount *big.Int) error {
	_, err := f.client.Decode(f.eventsMock.EmitOperatorRemoved(f.client.NewTXOpts(), operator, amount))
	return err
}

func (f *EthereumStakingEventsMock) FeedOperatorsSet(feedOperators []common.Address) error {
	_, err := f.client.Decode(f.eventsMock.EmitFeedOperatorsSet(f.client.NewTXOpts(), feedOperators))
	return err
}

func DeployStakingEventsMock(client *seth.Client) (StakingEventsMock, error) {
	abi, err := eth_contracts.StakingEventsMockMetaData.GetAbi()
	if err != nil {
		return &EthereumStakingEventsMock{}, fmt.Errorf("failed to get StakingEventsMock ABI: %w", err)
	}
	client.ContractStore.AddABI("StakingEventsMock", *abi)
	client.ContractStore.AddBIN("StakingEventsMock", common.FromHex(eth_contracts.StakingEventsMockMetaData.Bin))

	data, err := client.DeployContract(client.NewTXOpts(), "StakingEventsMock", *abi, common.FromHex(eth_contracts.StakingEventsMockMetaData.Bin))

	if err != nil {
		return &EthereumStakingEventsMock{}, fmt.Errorf("StakingEventsMock instance deployment have failed: %w", err)
	}

	instance, err := eth_contracts.NewStakingEventsMock(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumStakingEventsMock{}, fmt.Errorf("failed to instantiate StakingEventsMock instance: %w", err)
	}

	return &EthereumStakingEventsMock{
		client:     client,
		eventsMock: instance,
		address:    &data.Address,
	}, nil
}
