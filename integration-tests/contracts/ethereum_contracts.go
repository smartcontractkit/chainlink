package contracts

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	ocrConfigHelper "github.com/smartcontractkit/libocr/offchainreporting/confighelper"
	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting/types"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	contractsethereum "github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_coordinator"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_load_test_client"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_v1_events_mock"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/authorized_forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/flags_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/functions_billing_registry_events_mock"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/functions_oracle_events_mock"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/gas_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/gas_wrapper_mock"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_automation_registry_master_wrapper_2_2"
	iregistry22 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_automation_registry_master_wrapper_2_2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registrar_wrapper1_2_mock"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_1_mock"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_3"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_aggregator_proxy"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_ethlink_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_gas_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/operator_factory"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/operator_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/oracle_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/test_api_consumer_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/fee_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/reward_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_proxy"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/werc20_mock"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	eth_contracts "github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
)

// LegacyEthereumOracle oracle for "directrequest" job tests
type LegacyEthereumOracle struct {
	address *common.Address
	client  blockchain.EVMClient
	oracle  *oracle_wrapper.Oracle
}

func (e *LegacyEthereumOracle) Address() string {
	return e.address.Hex()
}

func (e *LegacyEthereumOracle) Fund(ethAmount *big.Float) error {
	gasEstimates, err := e.client.EstimateGas(ethereum.CallMsg{
		To: e.address,
	})
	if err != nil {
		return err
	}
	return e.client.Fund(e.address.Hex(), ethAmount, gasEstimates)
}

// SetFulfillmentPermission sets fulfillment permission for particular address
func (e *LegacyEthereumOracle) SetFulfillmentPermission(address string, allowed bool) error {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := e.oracle.SetFulfillmentPermission(opts, common.HexToAddress(address), allowed)
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx)
}

// LegacyEthereumAPIConsumer API consumer for job type "directrequest" tests
type LegacyEthereumAPIConsumer struct {
	address  *common.Address
	client   blockchain.EVMClient
	consumer *test_api_consumer_wrapper.TestAPIConsumer
}

func (e *LegacyEthereumAPIConsumer) Address() string {
	return e.address.Hex()
}

func (e *LegacyEthereumAPIConsumer) RoundID(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	return e.consumer.CurrentRoundID(opts)
}

func (e *LegacyEthereumAPIConsumer) Fund(ethAmount *big.Float) error {
	gasEstimates, err := e.client.EstimateGas(ethereum.CallMsg{
		To: e.address,
	})
	if err != nil {
		return err
	}
	return e.client.Fund(e.address.Hex(), ethAmount, gasEstimates)
}

func (e *LegacyEthereumAPIConsumer) Data(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	data, err := e.consumer.Data(opts)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// CreateRequestTo creates request to an oracle for particular jobID with params
func (e *LegacyEthereumAPIConsumer) CreateRequestTo(
	oracleAddr string,
	jobID [32]byte,
	payment *big.Int,
	url string,
	path string,
	times *big.Int,
) error {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := e.consumer.CreateRequestTo(opts, common.HexToAddress(oracleAddr), jobID, payment, url, path, times)
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx)
}

// EthereumStaking
type EthereumStaking struct {
	client  blockchain.EVMClient
	staking *eth_contracts.Staking
	address *common.Address
}

func (f *EthereumStaking) Address() string {
	return f.address.Hex()
}

// Fund sends specified currencies to the contract
func (f *EthereumStaking) Fund(ethAmount *big.Float) error {
	gasEstimates, err := f.client.EstimateGas(ethereum.CallMsg{
		To: f.address,
	})
	if err != nil {
		return err
	}
	return f.client.Fund(f.address.Hex(), ethAmount, gasEstimates)
}

func (f *EthereumStaking) AddOperators(operators []common.Address) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.staking.AddOperators(opts, operators)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumStaking) RemoveOperators(operators []common.Address) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.staking.RemoveOperators(opts, operators)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumStaking) SetFeedOperators(operators []common.Address) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.staking.SetFeedOperators(opts, operators)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumStaking) RaiseAlert() error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.staking.RaiseAlert(opts)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumStaking) Start(amount *big.Int, initialRewardRate *big.Int) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.staking.Start(opts, amount, initialRewardRate)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumStaking) SetMerkleRoot(newMerkleRoot [32]byte) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.staking.SetMerkleRoot(opts, newMerkleRoot)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

// EthereumFunctionsOracleEventsMock represents the basic events mock contract
type EthereumFunctionsOracleEventsMock struct {
	client     blockchain.EVMClient
	eventsMock *functions_oracle_events_mock.FunctionsOracleEventsMock
	address    *common.Address
}

func (f *EthereumFunctionsOracleEventsMock) Address() string {
	return f.address.Hex()
}

func (f *EthereumFunctionsOracleEventsMock) OracleResponse(requestId [32]byte) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitOracleResponse(opts, requestId)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumFunctionsOracleEventsMock) OracleRequest(requestId [32]byte, requestingContract common.Address, requestInitiator common.Address, subscriptionId uint64, subscriptionOwner common.Address, data []byte) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitOracleRequest(opts, requestId, requestingContract, requestInitiator, subscriptionId, subscriptionOwner, data)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumFunctionsOracleEventsMock) UserCallbackError(requestId [32]byte, reason string) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitUserCallbackError(opts, requestId, reason)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumFunctionsOracleEventsMock) UserCallbackRawError(requestId [32]byte, lowLevelData []byte) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitUserCallbackRawError(opts, requestId, lowLevelData)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

// EthereumFunctionsBillingRegistryEventsMock represents the basic events mock contract
type EthereumFunctionsBillingRegistryEventsMock struct {
	client     blockchain.EVMClient
	eventsMock *functions_billing_registry_events_mock.FunctionsBillingRegistryEventsMock
	address    *common.Address
}

func (f *EthereumFunctionsBillingRegistryEventsMock) Address() string {
	return f.address.Hex()
}

func (f *EthereumFunctionsBillingRegistryEventsMock) SubscriptionFunded(subscriptionId uint64, oldBalance *big.Int, newBalance *big.Int) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitSubscriptionFunded(opts, subscriptionId, oldBalance, newBalance)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumFunctionsBillingRegistryEventsMock) BillingStart(requestId [32]byte, commitment functions_billing_registry_events_mock.FunctionsBillingRegistryEventsMockCommitment) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitBillingStart(opts, requestId, commitment)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumFunctionsBillingRegistryEventsMock) BillingEnd(requestId [32]byte, subscriptionId uint64, signerPayment *big.Int, transmitterPayment *big.Int, totalCost *big.Int, success bool) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitBillingEnd(opts, requestId, subscriptionId, signerPayment, transmitterPayment, totalCost, success)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

// LegacyEthereumStakingEventsMock represents the basic events mock contract
type LegacyEthereumStakingEventsMock struct {
	client     blockchain.EVMClient
	eventsMock *eth_contracts.StakingEventsMock
	address    *common.Address
}

func (f *LegacyEthereumStakingEventsMock) Address() string {
	return f.address.Hex()
}

func (f *LegacyEthereumStakingEventsMock) MaxCommunityStakeAmountIncreased(maxStakeAmount *big.Int) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitMaxCommunityStakeAmountIncreased(opts, maxStakeAmount)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumStakingEventsMock) PoolSizeIncreased(maxPoolSize *big.Int) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitPoolSizeIncreased(opts, maxPoolSize)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumStakingEventsMock) MaxOperatorStakeAmountIncreased(maxStakeAmount *big.Int) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitMaxOperatorStakeAmountIncreased(opts, maxStakeAmount)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumStakingEventsMock) RewardInitialized(rate *big.Int, available *big.Int, startTimestamp *big.Int, endTimestamp *big.Int) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitRewardInitialized(opts, rate, available, startTimestamp, endTimestamp)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumStakingEventsMock) AlertRaised(alerter common.Address, roundId *big.Int, rewardAmount *big.Int) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitAlertRaised(opts, alerter, roundId, rewardAmount)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumStakingEventsMock) Staked(staker common.Address, newStake *big.Int, totalStake *big.Int) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitStaked(opts, staker, newStake, totalStake)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumStakingEventsMock) OperatorAdded(operator common.Address) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitOperatorAdded(opts, operator)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumStakingEventsMock) OperatorRemoved(operator common.Address, amount *big.Int) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitOperatorRemoved(opts, operator, amount)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumStakingEventsMock) FeedOperatorsSet(feedOperators []common.Address) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitFeedOperatorsSet(opts, feedOperators)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

// EthereumOffchainAggregatorEventsMock represents the basic events mock contract
type EthereumOffchainAggregatorEventsMock struct {
	client     blockchain.EVMClient
	eventsMock *eth_contracts.OffchainAggregatorEventsMock
	address    *common.Address
}

func (f *EthereumOffchainAggregatorEventsMock) Address() string {
	return f.address.Hex()
}

func (f *EthereumOffchainAggregatorEventsMock) ConfigSet(previousConfigBlockNumber uint32, configCount uint64, signers []common.Address, transmitters []common.Address, threshold uint8, encodedConfigVersion uint64, encoded []byte) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitConfigSet(opts, previousConfigBlockNumber, configCount, signers, transmitters, threshold, encodedConfigVersion, encoded)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumOffchainAggregatorEventsMock) NewTransmission(aggregatorRoundId uint32, answer *big.Int, transmitter common.Address, observations []*big.Int, observers []byte, rawReportContext [32]byte) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitNewTransmission(opts, aggregatorRoundId, answer, transmitter, observations, observers, rawReportContext)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

// LegacyEthereumKeeperRegistry11Mock represents the basic keeper registry 1.1 mock contract
type LegacyEthereumKeeperRegistry11Mock struct {
	client       blockchain.EVMClient
	registryMock *keeper_registry_wrapper1_1_mock.KeeperRegistryMock
	address      *common.Address
}

func (f *LegacyEthereumKeeperRegistry11Mock) Address() string {
	return f.address.Hex()
}

func (f *LegacyEthereumKeeperRegistry11Mock) EmitUpkeepPerformed(id *big.Int, success bool, from common.Address, payment *big.Int, performData []byte) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.registryMock.EmitUpkeepPerformed(opts, id, success, from, payment, performData)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumKeeperRegistry11Mock) EmitUpkeepCanceled(id *big.Int, atBlockHeight uint64) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.registryMock.EmitUpkeepCanceled(opts, id, atBlockHeight)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumKeeperRegistry11Mock) EmitFundsWithdrawn(id *big.Int, amount *big.Int, to common.Address) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.registryMock.EmitFundsWithdrawn(opts, id, amount, to)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumKeeperRegistry11Mock) EmitKeepersUpdated(keepers []common.Address, payees []common.Address) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.registryMock.EmitKeepersUpdated(opts, keepers, payees)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumKeeperRegistry11Mock) EmitUpkeepRegistered(id *big.Int, executeGas uint32, admin common.Address) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.registryMock.EmitUpkeepRegistered(opts, id, executeGas, admin)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumKeeperRegistry11Mock) EmitFundsAdded(id *big.Int, from common.Address, amount *big.Int) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.registryMock.EmitFundsAdded(opts, id, from, amount)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumKeeperRegistry11Mock) SetUpkeepCount(_upkeepCount *big.Int) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.registryMock.SetUpkeepCount(opts, _upkeepCount)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumKeeperRegistry11Mock) SetCanceledUpkeepList(_canceledUpkeepList []*big.Int) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.registryMock.SetCanceledUpkeepList(opts, _canceledUpkeepList)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumKeeperRegistry11Mock) SetKeeperList(_keepers []common.Address) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.registryMock.SetKeeperList(opts, _keepers)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumKeeperRegistry11Mock) SetConfig(_paymentPremiumPPB uint32, _flatFeeMicroLink uint32, _blockCountPerTurn *big.Int, _checkGasLimit uint32, _stalenessSeconds *big.Int, _gasCeilingMultiplier uint16, _fallbackGasPrice *big.Int, _fallbackLinkPrice *big.Int) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.registryMock.SetConfig(opts, _paymentPremiumPPB, _flatFeeMicroLink, _blockCountPerTurn, _checkGasLimit, _stalenessSeconds, _gasCeilingMultiplier, _fallbackGasPrice, _fallbackLinkPrice)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumKeeperRegistry11Mock) SetUpkeep(id *big.Int, _target common.Address, _executeGas uint32, _balance *big.Int, _admin common.Address, _maxValidBlocknumber uint64, _lastKeeper common.Address, _checkData []byte) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.registryMock.SetUpkeep(opts, id, _target, _executeGas, _balance, _admin, _maxValidBlocknumber, _lastKeeper, _checkData)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumKeeperRegistry11Mock) SetMinBalance(id *big.Int, minBalance *big.Int) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.registryMock.SetMinBalance(opts, id, minBalance)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumKeeperRegistry11Mock) SetCheckUpkeepData(id *big.Int, performData []byte, maxLinkPayment *big.Int, gasLimit *big.Int, adjustedGasWei *big.Int, linkEth *big.Int) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.registryMock.SetCheckUpkeepData(opts, id, performData, maxLinkPayment, gasLimit, adjustedGasWei, linkEth)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumKeeperRegistry11Mock) SetPerformUpkeepSuccess(id *big.Int, success bool) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.registryMock.SetPerformUpkeepSuccess(opts, id, success)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

// LegacyEthereumKeeperRegistrar12Mock represents the basic keeper registrar 1.2 mock contract
type LegacyEthereumKeeperRegistrar12Mock struct {
	client        blockchain.EVMClient
	registrarMock *keeper_registrar_wrapper1_2_mock.KeeperRegistrarMock
	address       *common.Address
}

func (f *LegacyEthereumKeeperRegistrar12Mock) Address() string {
	return f.address.Hex()
}

func (f *LegacyEthereumKeeperRegistrar12Mock) EmitRegistrationRequested(hash [32]byte, name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, amount *big.Int, source uint8) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.registrarMock.EmitRegistrationRequested(opts, hash, name, encryptedEmail, upkeepContract, gasLimit, adminAddress, checkData, amount, source)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumKeeperRegistrar12Mock) EmitRegistrationApproved(hash [32]byte, displayName string, upkeepId *big.Int) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.registrarMock.EmitRegistrationApproved(opts, hash, displayName, upkeepId)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumKeeperRegistrar12Mock) SetRegistrationConfig(_autoApproveConfigType uint8, _autoApproveMaxAllowed uint32, _approvedCount uint32, _keeperRegistry common.Address, _minLINKJuels *big.Int) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.registrarMock.SetRegistrationConfig(opts, _autoApproveConfigType, _autoApproveMaxAllowed, _approvedCount, _keeperRegistry, _minLINKJuels)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

// LegacyEthereumKeeperGasWrapperMock represents the basic keeper gas wrapper mock contract
type LegacyEthereumKeeperGasWrapperMock struct {
	client         blockchain.EVMClient
	gasWrapperMock *gas_wrapper_mock.KeeperRegistryCheckUpkeepGasUsageWrapperMock
	address        *common.Address
}

func (f *LegacyEthereumKeeperGasWrapperMock) Address() string {
	return f.address.Hex()
}

func (f *LegacyEthereumKeeperGasWrapperMock) SetMeasureCheckGasResult(result bool, payload []byte, gas *big.Int) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.gasWrapperMock.SetMeasureCheckGasResult(opts, result, payload, gas)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

// LegacyEthereumFunctionsV1EventsMock represents the basic functions v1 events mock contract
type LegacyEthereumFunctionsV1EventsMock struct {
	client     blockchain.EVMClient
	eventsMock *functions_v1_events_mock.FunctionsV1EventsMock
	address    *common.Address
}

func (f *LegacyEthereumFunctionsV1EventsMock) Address() string {
	return f.address.Hex()
}

func (f *LegacyEthereumFunctionsV1EventsMock) EmitRequestProcessed(requestId [32]byte, subscriptionId uint64, totalCostJuels *big.Int, transmitter common.Address, resultCode uint8, response []byte, errByte []byte, callbackReturnData []byte) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitRequestProcessed(opts, requestId, subscriptionId, totalCostJuels, transmitter, resultCode, response, errByte, callbackReturnData)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumFunctionsV1EventsMock) EmitRequestStart(requestId [32]byte, donId [32]byte, subscriptionId uint64, subscriptionOwner common.Address, requestingContract common.Address, requestInitiator common.Address, data []byte, dataVersion uint16, callbackGasLimit uint32, estimatedTotalCostJuels *big.Int) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitRequestStart(opts, requestId, donId, subscriptionId, subscriptionOwner, requestingContract, requestInitiator, data, dataVersion, callbackGasLimit, estimatedTotalCostJuels)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumFunctionsV1EventsMock) EmitSubscriptionCanceled(subscriptionId uint64, fundsRecipient common.Address, fundsAmount *big.Int) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitSubscriptionCanceled(opts, subscriptionId, fundsRecipient, fundsAmount)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumFunctionsV1EventsMock) EmitSubscriptionConsumerAdded(subscriptionId uint64, consumer common.Address) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitSubscriptionConsumerAdded(opts, subscriptionId, consumer)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumFunctionsV1EventsMock) EmitSubscriptionConsumerRemoved(subscriptionId uint64, consumer common.Address) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitSubscriptionConsumerRemoved(opts, subscriptionId, consumer)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumFunctionsV1EventsMock) EmitSubscriptionCreated(subscriptionId uint64, owner common.Address) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitSubscriptionCreated(opts, subscriptionId, owner)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumFunctionsV1EventsMock) EmitSubscriptionFunded(subscriptionId uint64, oldBalance *big.Int, newBalance *big.Int) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitSubscriptionFunded(opts, subscriptionId, oldBalance, newBalance)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumFunctionsV1EventsMock) EmitSubscriptionOwnerTransferred(subscriptionId uint64, from common.Address, to common.Address) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitSubscriptionOwnerTransferred(opts, subscriptionId, from, to)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumFunctionsV1EventsMock) EmitSubscriptionOwnerTransferRequested(subscriptionId uint64, from common.Address, to common.Address) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitSubscriptionOwnerTransferRequested(opts, subscriptionId, from, to)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumFunctionsV1EventsMock) EmitRequestNotProcessed(requestId [32]byte, coordinator common.Address, transmitter common.Address, resultCode uint8) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitRequestNotProcessed(opts, requestId, coordinator, transmitter, resultCode)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumFunctionsV1EventsMock) EmitContractUpdated(id [32]byte, from common.Address, to common.Address) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.eventsMock.EmitContractUpdated(opts, id, from, to)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

// LegacyEthereumFluxAggregator represents the basic flux aggregation contract
type LegacyEthereumFluxAggregator struct {
	client         blockchain.EVMClient
	fluxAggregator *flux_aggregator_wrapper.FluxAggregator
	address        *common.Address
}

func (f *LegacyEthereumFluxAggregator) Address() string {
	return f.address.Hex()
}

// Fund sends specified currencies to the contract
func (f *LegacyEthereumFluxAggregator) Fund(ethAmount *big.Float) error {
	gasEstimates, err := f.client.EstimateGas(ethereum.CallMsg{
		To: f.address,
	})
	if err != nil {
		return err
	}
	return f.client.Fund(f.address.Hex(), ethAmount, gasEstimates)
}

func (f *LegacyEthereumFluxAggregator) UpdateAvailableFunds() error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.fluxAggregator.UpdateAvailableFunds(opts)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumFluxAggregator) PaymentAmount(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	payment, err := f.fluxAggregator.PaymentAmount(opts)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (f *LegacyEthereumFluxAggregator) RequestNewRound(_ context.Context) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.fluxAggregator.RequestNewRound(opts)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

// WatchSubmissionReceived subscribes to any submissions on a flux feed
func (f *LegacyEthereumFluxAggregator) WatchSubmissionReceived(ctx context.Context, eventChan chan<- *SubmissionEvent) error {
	ethEventChan := make(chan *flux_aggregator_wrapper.FluxAggregatorSubmissionReceived)
	sub, err := f.fluxAggregator.WatchSubmissionReceived(&bind.WatchOpts{}, ethEventChan, nil, nil, nil)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case event := <-ethEventChan:
			eventChan <- &SubmissionEvent{
				Contract:    event.Raw.Address,
				Submission:  event.Submission,
				Round:       event.Round,
				BlockNumber: event.Raw.BlockNumber,
				Oracle:      event.Oracle,
			}
		case err := <-sub.Err():
			return err
		case <-ctx.Done():
			return nil
		}
	}
}

func (f *LegacyEthereumFluxAggregator) SetRequesterPermissions(_ context.Context, addr common.Address, authorized bool, roundsDelay uint32) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.fluxAggregator.SetRequesterPermissions(opts, addr, authorized, roundsDelay)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumFluxAggregator) GetOracles(ctx context.Context) ([]string, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	addresses, err := f.fluxAggregator.GetOracles(opts)
	if err != nil {
		return nil, err
	}
	var oracleAddrs []string
	for _, o := range addresses {
		oracleAddrs = append(oracleAddrs, o.Hex())
	}
	return oracleAddrs, nil
}

func (f *LegacyEthereumFluxAggregator) LatestRoundID(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	rID, err := f.fluxAggregator.LatestRound(opts)
	if err != nil {
		return nil, err
	}
	return rID, nil
}

func (f *LegacyEthereumFluxAggregator) WithdrawPayment(
	_ context.Context,
	from common.Address,
	to common.Address,
	amount *big.Int) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := f.fluxAggregator.WithdrawPayment(opts, from, to, amount)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *LegacyEthereumFluxAggregator) WithdrawablePayment(ctx context.Context, addr common.Address) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	balance, err := f.fluxAggregator.WithdrawablePayment(opts, addr)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (f *LegacyEthereumFluxAggregator) LatestRoundData(ctx context.Context) (flux_aggregator_wrapper.LatestRoundData, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	lr, err := f.fluxAggregator.LatestRoundData(opts)
	if err != nil {
		return flux_aggregator_wrapper.LatestRoundData{}, err
	}
	return lr, nil
}

// GetContractData retrieves basic data for the flux aggregator contract
func (f *LegacyEthereumFluxAggregator) GetContractData(ctx context.Context) (*FluxAggregatorData, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.client.GetDefaultWallet().Address()),
		Context: ctx,
	}

	allocated, err := f.fluxAggregator.AllocatedFunds(opts)
	if err != nil {
		return &FluxAggregatorData{}, err
	}

	available, err := f.fluxAggregator.AvailableFunds(opts)
	if err != nil {
		return &FluxAggregatorData{}, err
	}

	lr, err := f.fluxAggregator.LatestRoundData(opts)
	if err != nil {
		return &FluxAggregatorData{}, err
	}
	latestRound := RoundData(lr)

	oracles, err := f.fluxAggregator.GetOracles(opts)
	if err != nil {
		return &FluxAggregatorData{}, err
	}

	return &FluxAggregatorData{
		AllocatedFunds:  allocated,
		AvailableFunds:  available,
		LatestRoundData: latestRound,
		Oracles:         oracles,
	}, nil
}

// SetOracles allows the ability to add and/or remove oracles from the contract, and to set admins
func (f *LegacyEthereumFluxAggregator) SetOracles(o FluxAggregatorSetOraclesOptions) error {
	opts, err := f.client.TransactionOpts(f.client.GetDefaultWallet())
	if err != nil {
		return err
	}

	tx, err := f.fluxAggregator.ChangeOracles(opts, o.RemoveList, o.AddList, o.AdminList, o.MinSubmissions, o.MaxSubmissions, o.RestartDelayRounds)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

// Description returns the description of the flux aggregator contract
func (f *LegacyEthereumFluxAggregator) Description(ctxt context.Context) (string, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.client.GetDefaultWallet().Address()),
		Context: ctxt,
	}
	return f.fluxAggregator.Description(opts)
}

// FluxAggregatorRoundConfirmer is a header subscription that awaits for a certain flux round to be completed
type FluxAggregatorRoundConfirmer struct {
	fluxInstance FluxAggregator
	roundID      *big.Int
	doneChan     chan struct{}
	context      context.Context
	cancel       context.CancelFunc
	complete     bool
	l            zerolog.Logger
}

// NewFluxAggregatorRoundConfirmer provides a new instance of a FluxAggregatorRoundConfirmer
func NewFluxAggregatorRoundConfirmer(
	contract FluxAggregator,
	roundID *big.Int,
	timeout time.Duration,
	logger zerolog.Logger,
) *FluxAggregatorRoundConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
	return &FluxAggregatorRoundConfirmer{
		fluxInstance: contract,
		roundID:      roundID,
		doneChan:     make(chan struct{}),
		context:      ctx,
		cancel:       ctxCancel,
		l:            logger,
	}
}

// ReceiveHeader will query the latest FluxAggregator round and check to see whether the round has confirmed
func (f *FluxAggregatorRoundConfirmer) ReceiveHeader(header blockchain.NodeHeader) error {
	if f.complete {
		return nil
	}
	lr, err := f.fluxInstance.LatestRoundID(context.Background())
	if err != nil {
		return err
	}
	logFields := map[string]any{
		"Contract Address":  f.fluxInstance.Address(),
		"Current Round":     lr.Int64(),
		"Waiting for Round": f.roundID.Int64(),
		"Header Number":     header.Number.Uint64(),
	}
	if lr.Cmp(f.roundID) >= 0 {
		f.l.Info().Fields(logFields).Msg("FluxAggregator round completed")
		f.complete = true
		f.doneChan <- struct{}{}
	} else {
		f.l.Debug().Fields(logFields).Msg("Waiting for FluxAggregator round")
	}
	return nil
}

// Wait is a blocking function that will wait until the round has confirmed, and timeout if the deadline has passed
func (f *FluxAggregatorRoundConfirmer) Wait() error {
	defer func() { f.complete = true }()
	for {
		select {
		case <-f.doneChan:
			f.cancel()
			return nil
		case <-f.context.Done():
			return fmt.Errorf("timeout waiting for flux round to confirm: %d", f.roundID)
		}
	}
}

func (f *FluxAggregatorRoundConfirmer) Complete() bool {
	return f.complete
}

// LegacyEthereumLinkToken represents a LinkToken address
type LegacyEthereumLinkToken struct {
	client   blockchain.EVMClient
	instance *link_token_interface.LinkToken
	address  common.Address
	l        zerolog.Logger
}

// Fund the LINK Token contract with ETH to distribute the token
func (l *LegacyEthereumLinkToken) Fund(ethAmount *big.Float) error {
	gasEstimates, err := l.client.EstimateGas(ethereum.CallMsg{
		To: &l.address,
	})
	if err != nil {
		return err
	}
	return l.client.Fund(l.address.Hex(), ethAmount, gasEstimates)
}

func (l *LegacyEthereumLinkToken) BalanceOf(ctx context.Context, addr string) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(l.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	balance, err := l.instance.BalanceOf(opts, common.HexToAddress(addr))
	if err != nil {
		return nil, err
	}
	return balance, nil
}

// Name returns the name of the link token
func (l *LegacyEthereumLinkToken) Name(ctxt context.Context) (string, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(l.client.GetDefaultWallet().Address()),
		Context: ctxt,
	}
	return l.instance.Name(opts)
}

func (l *LegacyEthereumLinkToken) Address() string {
	return l.address.Hex()
}

func (l *LegacyEthereumLinkToken) Approve(to string, amount *big.Int) error {
	opts, err := l.client.TransactionOpts(l.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	l.l.Info().
		Str("From", l.client.GetDefaultWallet().Address()).
		Str("To", to).
		Str("Amount", amount.String()).
		Uint64("Nonce", opts.Nonce.Uint64()).
		Msg("Approving LINK Transfer")
	tx, err := l.instance.Approve(opts, common.HexToAddress(to), amount)
	if err != nil {
		return err
	}
	return l.client.ProcessTransaction(tx)
}

func (l *LegacyEthereumLinkToken) Transfer(to string, amount *big.Int) error {
	opts, err := l.client.TransactionOpts(l.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	l.l.Info().
		Str("From", l.client.GetDefaultWallet().Address()).
		Str("To", to).
		Str("Amount", amount.String()).
		Uint64("Nonce", opts.Nonce.Uint64()).
		Msg("Transferring LINK")
	tx, err := l.instance.Transfer(opts, common.HexToAddress(to), amount)
	if err != nil {
		return err
	}
	return l.client.ProcessTransaction(tx)
}

func (l *LegacyEthereumLinkToken) TransferAndCall(to string, amount *big.Int, data []byte) (*types.Transaction, error) {
	opts, err := l.client.TransactionOpts(l.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	tx, err := l.instance.TransferAndCall(opts, common.HexToAddress(to), amount, data)
	if err != nil {
		return nil, err
	}
	l.l.Info().
		Str("From", l.client.GetDefaultWallet().Address()).
		Str("To", to).
		Str("Amount", amount.String()).
		Uint64("Nonce", opts.Nonce.Uint64()).
		Str("TxHash", tx.Hash().String()).
		Msg("Transferring and Calling LINK")
	return tx, l.client.ProcessTransaction(tx)
}

func (l *LegacyEthereumLinkToken) TransferAndCallFromKey(_ string, _ *big.Int, _ []byte, _ int) (*types.Transaction, error) {
	panic("supported only with Seth")
}

// LegacyEthereumOffchainAggregator represents the offchain aggregation contract
// Deprecated: we are moving away from blockchain.EVMClient, use EthereumOffchainAggregator instead
type LegacyEthereumOffchainAggregator struct {
	client  blockchain.EVMClient
	ocr     *offchainaggregator.OffchainAggregator
	address *common.Address
	l       zerolog.Logger
}

// SetPayees sets wallets for the contract to pay out to?
func (o *LegacyEthereumOffchainAggregator) SetPayees(
	transmitters, payees []string,
) error {
	opts, err := o.client.TransactionOpts(o.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	var transmittersAddr, payeesAddr []common.Address
	for _, tr := range transmitters {
		transmittersAddr = append(transmittersAddr, common.HexToAddress(tr))
	}
	for _, p := range payees {
		payeesAddr = append(payeesAddr, common.HexToAddress(p))
	}

	o.l.Info().
		Str("Transmitters", fmt.Sprintf("%v", transmitters)).
		Str("Payees", fmt.Sprintf("%v", payees)).
		Str("OCR Address", o.Address()).
		Msg("Setting OCR Payees")

	tx, err := o.ocr.SetPayees(opts, transmittersAddr, payeesAddr)
	if err != nil {
		return err
	}
	return o.client.ProcessTransaction(tx)
}

// SetConfig sets the payees and the offchain reporting protocol configuration
func (o *LegacyEthereumOffchainAggregator) SetConfig(
	chainlinkNodes []ChainlinkNodeWithKeysAndAddress,
	ocrConfig OffChainAggregatorConfig,
	transmitters []common.Address,
) error {
	// Gather necessary addresses and keys from our chainlink nodes to properly configure the OCR contract
	log.Info().Str("Contract Address", o.address.Hex()).Msg("Configuring OCR Contract")
	for i, node := range chainlinkNodes {
		ocrKeys, err := node.MustReadOCRKeys()
		if err != nil {
			return err
		}
		if len(ocrKeys.Data) == 0 {
			return fmt.Errorf("no OCR keys found for node %v", node)
		}
		primaryOCRKey := ocrKeys.Data[0]
		p2pKeys, err := node.MustReadP2PKeys()
		if err != nil {
			return err
		}
		primaryP2PKey := p2pKeys.Data[0]

		// Need to convert the key representations
		var onChainSigningAddress [20]byte
		var configPublicKey [32]byte
		offchainSigningAddress, err := hex.DecodeString(primaryOCRKey.Attributes.OffChainPublicKey)
		if err != nil {
			return err
		}
		decodeConfigKey, err := hex.DecodeString(primaryOCRKey.Attributes.ConfigPublicKey)
		if err != nil {
			return err
		}

		// https://stackoverflow.com/questions/8032170/how-to-assign-string-to-bytes-array
		copy(onChainSigningAddress[:], common.HexToAddress(primaryOCRKey.Attributes.OnChainSigningAddress).Bytes())
		copy(configPublicKey[:], decodeConfigKey)

		oracleIdentity := ocrConfigHelper.OracleIdentity{
			TransmitAddress:       transmitters[i],
			OnChainSigningAddress: onChainSigningAddress,
			PeerID:                primaryP2PKey.Attributes.PeerID,
			OffchainPublicKey:     offchainSigningAddress,
		}
		oracleIdentityExtra := ocrConfigHelper.OracleIdentityExtra{
			OracleIdentity:                  oracleIdentity,
			SharedSecretEncryptionPublicKey: ocrTypes.SharedSecretEncryptionPublicKey(configPublicKey),
		}

		ocrConfig.OracleIdentities = append(ocrConfig.OracleIdentities, oracleIdentityExtra)
	}

	signers, transmitters, threshold, encodedConfigVersion, encodedConfig, err := ocrConfigHelper.ContractSetConfigArgs(
		ocrConfig.DeltaProgress,
		ocrConfig.DeltaResend,
		ocrConfig.DeltaRound,
		ocrConfig.DeltaGrace,
		ocrConfig.DeltaC,
		ocrConfig.AlphaPPB,
		ocrConfig.DeltaStage,
		ocrConfig.RMax,
		ocrConfig.S,
		ocrConfig.OracleIdentities,
		ocrConfig.F,
	)
	if err != nil {
		return err
	}

	// Set Config
	opts, err := o.client.TransactionOpts(o.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := o.ocr.SetConfig(opts, signers, transmitters, threshold, encodedConfigVersion, encodedConfig)
	if err != nil {
		return err
	}
	return o.client.ProcessTransaction(tx)
}

// RequestNewRound requests the OCR contract to create a new round
func (o *LegacyEthereumOffchainAggregator) RequestNewRound() error {
	opts, err := o.client.TransactionOpts(o.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := o.ocr.RequestNewRound(opts)
	if err != nil {
		return err
	}
	o.l.Info().Str("Contract Address", o.address.Hex()).Msg("New OCR round requested")

	return o.client.ProcessTransaction(tx)
}

// GetLatestAnswer returns the latest answer from the OCR contract
func (o *LegacyEthereumOffchainAggregator) GetLatestAnswer(ctxt context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(o.client.GetDefaultWallet().Address()),
		Context: ctxt,
	}
	return o.ocr.LatestAnswer(opts)
}

func (o *LegacyEthereumOffchainAggregator) Address() string {
	return o.address.Hex()
}

// GetLatestRound returns data from the latest round
func (o *LegacyEthereumOffchainAggregator) GetLatestRound(ctx context.Context) (*RoundData, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(o.client.GetDefaultWallet().Address()),
		Context: ctx,
	}

	roundData, err := o.ocr.LatestRoundData(opts)
	if err != nil {
		return nil, err
	}

	return &RoundData{
		RoundId:         roundData.RoundId,
		Answer:          roundData.Answer,
		AnsweredInRound: roundData.AnsweredInRound,
		StartedAt:       roundData.StartedAt,
		UpdatedAt:       roundData.UpdatedAt,
	}, err
}

func (o *LegacyEthereumOffchainAggregator) LatestRoundDataUpdatedAt() (*big.Int, error) {
	data, err := o.ocr.LatestRoundData(&bind.CallOpts{
		From:    common.HexToAddress(o.client.GetDefaultWallet().Address()),
		Context: context.Background(),
	})
	if err != nil {
		return nil, err
	}
	return data.UpdatedAt, nil
}

// GetRound retrieves an OCR round by the round ID
func (o *LegacyEthereumOffchainAggregator) GetRound(ctx context.Context, roundID *big.Int) (*RoundData, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(o.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	roundData, err := o.ocr.GetRoundData(opts, roundID)
	if err != nil {
		return nil, err
	}

	return &RoundData{
		RoundId:         roundData.RoundId,
		Answer:          roundData.Answer,
		AnsweredInRound: roundData.AnsweredInRound,
		StartedAt:       roundData.StartedAt,
		UpdatedAt:       roundData.UpdatedAt,
	}, nil
}

// ParseEventAnswerUpdated parses the log for event AnswerUpdated
func (o *LegacyEthereumOffchainAggregator) ParseEventAnswerUpdated(eventLog types.Log) (*offchainaggregator.OffchainAggregatorAnswerUpdated, error) {
	return o.ocr.ParseAnswerUpdated(eventLog)
}

// RunlogRoundConfirmer is a header subscription that awaits for a certain Runlog round to be completed
type RunlogRoundConfirmer struct {
	consumer APIConsumer
	roundID  *big.Int
	doneChan chan struct{}
	context  context.Context
	cancel   context.CancelFunc
	l        zerolog.Logger
}

// NewRunlogRoundConfirmer provides a new instance of a RunlogRoundConfirmer
func NewRunlogRoundConfirmer(
	contract APIConsumer,
	roundID *big.Int,
	timeout time.Duration,
	logger zerolog.Logger,
) *RunlogRoundConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
	return &RunlogRoundConfirmer{
		consumer: contract,
		roundID:  roundID,
		doneChan: make(chan struct{}),
		context:  ctx,
		cancel:   ctxCancel,
		l:        logger,
	}
}

// ReceiveHeader will query the latest Runlog round and check to see whether the round has confirmed
func (o *RunlogRoundConfirmer) ReceiveHeader(_ blockchain.NodeHeader) error {
	currentRoundID, err := o.consumer.RoundID(context.Background())
	if err != nil {
		return err
	}
	logFields := map[string]any{
		"Contract Address":  o.consumer.Address(),
		"Current Round":     currentRoundID.Int64(),
		"Waiting for Round": o.roundID.Int64(),
	}
	if currentRoundID.Cmp(o.roundID) >= 0 {
		o.l.Info().Fields(logFields).Msg("Runlog round completed")
		o.doneChan <- struct{}{}
	} else {
		o.l.Debug().Fields(logFields).Msg("Waiting for Runlog round")
	}
	return nil
}

// Wait is a blocking function that will wait until the round has confirmed, and timeout if the deadline has passed
func (o *RunlogRoundConfirmer) Wait() error {
	for {
		select {
		case <-o.doneChan:
			o.cancel()
			return nil
		case <-o.context.Done():
			return fmt.Errorf("timeout waiting for OCR round to confirm: %d", o.roundID)
		}
	}
}

// OffchainAggregatorRoundConfirmer is a header subscription that awaits for a certain OCR round to be completed
type OffchainAggregatorRoundConfirmer struct {
	ocrInstance       OffchainAggregator
	roundID           *big.Int
	doneChan          chan struct{}
	context           context.Context
	cancel            context.CancelFunc
	blocksSinceAnswer uint
	complete          bool
	l                 zerolog.Logger
}

// NewOffchainAggregatorRoundConfirmer provides a new instance of a OffchainAggregatorRoundConfirmer
func NewOffchainAggregatorRoundConfirmer(
	contract OffchainAggregator,
	roundID *big.Int,
	timeout time.Duration,
	logger zerolog.Logger,
) *OffchainAggregatorRoundConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
	return &OffchainAggregatorRoundConfirmer{
		ocrInstance: contract,
		roundID:     roundID,
		doneChan:    make(chan struct{}),
		context:     ctx,
		cancel:      ctxCancel,
		complete:    false,
		l:           logger,
	}
}

// ReceiveHeader will query the latest OffchainAggregator round and check to see whether the round has confirmed
func (o *OffchainAggregatorRoundConfirmer) ReceiveHeader(_ blockchain.NodeHeader) error {
	if channelClosed(o.doneChan) {
		return nil
	}

	lr, err := o.ocrInstance.GetLatestRound(context.Background())
	if err != nil {
		return err
	}
	o.blocksSinceAnswer++
	currRound := lr.RoundId
	logFields := map[string]any{
		"Contract Address":  o.ocrInstance.Address(),
		"Current Round":     currRound.Int64(),
		"Waiting for Round": o.roundID.Int64(),
	}
	if currRound.Cmp(o.roundID) >= 0 {
		o.l.Info().Fields(logFields).Msg("OCR round completed")
		o.doneChan <- struct{}{}
		o.complete = true
	} else {
		o.l.Debug().Fields(logFields).Msg("Waiting on OCR Round")
	}
	return nil
}

// Wait is a blocking function that will wait until the round has confirmed, and timeout if the deadline has passed
func (o *OffchainAggregatorRoundConfirmer) Wait() error {
	defer func() { o.complete = true }()
	for {
		select {
		case <-o.doneChan:
			o.cancel()
			close(o.doneChan)
			return nil
		case <-o.context.Done():
			return fmt.Errorf("timeout waiting for OCR round to confirm: %d", o.roundID)
		}
	}
}

func (o *OffchainAggregatorRoundConfirmer) Complete() bool {
	return o.complete
}

// OffchainAggregatorRoundConfirmer is a header subscription that awaits for a certain OCR round to be completed
type OffchainAggregatorV2RoundConfirmer struct {
	ocrInstance       OffchainAggregatorV2
	roundID           *big.Int
	doneChan          chan struct{}
	context           context.Context
	cancel            context.CancelFunc
	blocksSinceAnswer uint
	complete          bool
	l                 zerolog.Logger
}

// NewOffchainAggregatorRoundConfirmer provides a new instance of a OffchainAggregatorRoundConfirmer
func NewOffchainAggregatorV2RoundConfirmer(
	contract OffchainAggregatorV2,
	roundID *big.Int,
	timeout time.Duration,
	logger zerolog.Logger,
) *OffchainAggregatorV2RoundConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
	return &OffchainAggregatorV2RoundConfirmer{
		ocrInstance: contract,
		roundID:     roundID,
		doneChan:    make(chan struct{}),
		context:     ctx,
		cancel:      ctxCancel,
		complete:    false,
		l:           logger,
	}
}

// ReceiveHeader will query the latest OffchainAggregator round and check to see whether the round has confirmed
func (o *OffchainAggregatorV2RoundConfirmer) ReceiveHeader(_ blockchain.NodeHeader) error {
	if channelClosed(o.doneChan) {
		return nil
	}

	lr, err := o.ocrInstance.GetLatestRound(context.Background())
	if err != nil {
		return err
	}
	o.blocksSinceAnswer++
	currRound := lr.RoundId
	logFields := map[string]any{
		"Contract Address":  o.ocrInstance.Address(),
		"Current Round":     currRound.Int64(),
		"Waiting for Round": o.roundID.Int64(),
	}
	if currRound.Cmp(o.roundID) >= 0 {
		o.l.Info().Fields(logFields).Msg("OCR round completed")
		o.doneChan <- struct{}{}
		o.complete = true
	} else {
		o.l.Debug().Fields(logFields).Msg("Waiting on OCR Round")
	}
	return nil
}

// Wait is a blocking function that will wait until the round has confirmed, and timeout if the deadline has passed
func (o *OffchainAggregatorV2RoundConfirmer) Wait() error {
	defer func() { o.complete = true }()
	for {
		select {
		case <-o.doneChan:
			o.cancel()
			close(o.doneChan)
			return nil
		case <-o.context.Done():
			return fmt.Errorf("timeout waiting for OCR round to confirm: %d", o.roundID)
		}
	}
}

func (o *OffchainAggregatorV2RoundConfirmer) Complete() bool {
	return o.complete
}

// LegacyEthereumMockETHLINKFeed represents mocked ETH/LINK feed contract
type LegacyEthereumMockETHLINKFeed struct {
	client  blockchain.EVMClient
	feed    *mock_ethlink_aggregator_wrapper.MockETHLINKAggregator
	address *common.Address
}

func (v *LegacyEthereumMockETHLINKFeed) Address() string {
	return v.address.Hex()
}

func (v *LegacyEthereumMockETHLINKFeed) LatestRoundData() (*big.Int, error) {
	data, err := v.feed.LatestRoundData(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: context.Background(),
	})
	if err != nil {
		return nil, err
	}
	return data.Ans, nil
}

func (v *LegacyEthereumMockETHLINKFeed) LatestRoundDataUpdatedAt() (*big.Int, error) {
	data, err := v.feed.LatestRoundData(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: context.Background(),
	})
	if err != nil {
		return nil, err
	}
	return data.UpdatedAt, nil
}

// LegacyEthereumMockGASFeed represents mocked Gas feed contract
type LegacyEthereumMockGASFeed struct {
	client  blockchain.EVMClient
	feed    *mock_gas_aggregator_wrapper.MockGASAggregator
	address *common.Address
}

func (v *LegacyEthereumMockGASFeed) Address() string {
	return v.address.Hex()
}

// EthereumFlags represents flags contract
type EthereumFlags struct {
	client  blockchain.EVMClient
	flags   *flags_wrapper.Flags
	address *common.Address
}

func (e *EthereumFlags) Address() string {
	return e.address.Hex()
}

// GetFlag returns boolean if a flag was set for particular address
func (e *EthereumFlags) GetFlag(ctx context.Context, addr string) (bool, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	flag, err := e.flags.GetFlag(opts, common.HexToAddress(addr))
	if err != nil {
		return false, err
	}
	return flag, nil
}

// LegacyEthereumOperatorFactory represents operator factory contract
// Deprecated: we are moving away from blockchain.EVMClient, use EthereumOperatorFactory instead
type LegacyEthereumOperatorFactory struct {
	address         *common.Address
	client          blockchain.EVMClient
	operatorFactory *operator_factory.OperatorFactory
}

func (e *LegacyEthereumOperatorFactory) ParseAuthorizedForwarderCreated(eventLog types.Log) (*operator_factory.OperatorFactoryAuthorizedForwarderCreated, error) {
	return e.operatorFactory.ParseAuthorizedForwarderCreated(eventLog)
}

func (e *LegacyEthereumOperatorFactory) ParseOperatorCreated(eventLog types.Log) (*operator_factory.OperatorFactoryOperatorCreated, error) {
	return e.operatorFactory.ParseOperatorCreated(eventLog)
}

func (e *LegacyEthereumOperatorFactory) Address() string {
	return e.address.Hex()
}

func (e *LegacyEthereumOperatorFactory) DeployNewOperatorAndForwarder() (*types.Transaction, error) {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	tx, err := e.operatorFactory.DeployNewOperatorAndForwarder(opts)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// LegacyEthereumOperator represents operator contract
// Deprecated: we are moving away from blockchain.EVMClient, use EthereumOperator instead
type LegacyEthereumOperator struct {
	address  common.Address
	client   blockchain.EVMClient
	operator *operator_wrapper.Operator
	l        zerolog.Logger
}

func (e *LegacyEthereumOperator) Address() string {
	return e.address.Hex()
}

func (e *LegacyEthereumOperator) AcceptAuthorizedReceivers(forwarders []common.Address, eoa []common.Address) error {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	e.l.Info().
		Str("ForwardersAddresses", fmt.Sprint(forwarders)).
		Str("EoaAddresses", fmt.Sprint(eoa)).
		Msg("Accepting Authorized Receivers")
	tx, err := e.operator.AcceptAuthorizedReceivers(opts, forwarders, eoa)
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx)
}

// LegacyEthereumAuthorizedForwarder represents authorized forwarder contract
// Deprecated: we are moving away from blockchain.EVMClient, use EthereumAuthorizedForwarder instead
type LegacyEthereumAuthorizedForwarder struct {
	address             common.Address
	client              blockchain.EVMClient
	authorizedForwarder *authorized_forwarder.AuthorizedForwarder
}

// Owner return authorized forwarder owner address
func (e *LegacyEthereumAuthorizedForwarder) Owner(ctx context.Context) (string, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	owner, err := e.authorizedForwarder.Owner(opts)

	return owner.Hex(), err
}

func (e *LegacyEthereumAuthorizedForwarder) GetAuthorizedSenders(ctx context.Context) ([]string, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	authorizedSenders, err := e.authorizedForwarder.GetAuthorizedSenders(opts)
	if err != nil {
		return nil, err
	}
	var sendersAddrs []string
	for _, o := range authorizedSenders {
		sendersAddrs = append(sendersAddrs, o.Hex())
	}
	return sendersAddrs, nil
}

func (e *LegacyEthereumAuthorizedForwarder) Address() string {
	return e.address.Hex()
}

// EthereumMockAggregatorProxy represents mock aggregator proxy contract
type EthereumMockAggregatorProxy struct {
	address             *common.Address
	client              blockchain.EVMClient
	mockAggregatorProxy *mock_aggregator_proxy.MockAggregatorProxy
}

func (e *EthereumMockAggregatorProxy) Address() string {
	return e.address.Hex()
}

func (e *EthereumMockAggregatorProxy) UpdateAggregator(aggregator common.Address) error {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := e.mockAggregatorProxy.UpdateAggregator(opts, aggregator)
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx)
}

func (e *EthereumMockAggregatorProxy) Aggregator() (common.Address, error) {
	addr, err := e.mockAggregatorProxy.Aggregator(&bind.CallOpts{
		From:    common.HexToAddress(e.client.GetDefaultWallet().Address()),
		Context: context.Background(),
	})
	if err != nil {
		return common.Address{}, err
	}
	return addr, nil
}

func channelClosed(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

// Deprecated: we are moving away from blockchain.EVMClient, use EthereumOffchainAggregatorV2 instead
type LegacyEthereumOffchainAggregatorV2 struct {
	address  *common.Address
	client   blockchain.EVMClient
	contract *ocr2aggregator.OCR2Aggregator
	l        zerolog.Logger
}

// OCRv2Config represents the config for the OCRv2 contract
type OCRv2Config struct {
	Signers               []common.Address
	Transmitters          []common.Address
	F                     uint8
	OnchainConfig         []byte
	TypedOnchainConfig21  i_keeper_registry_master_wrapper_2_1.IAutomationV21PlusCommonOnchainConfigLegacy
	TypedOnchainConfig22  i_automation_registry_master_wrapper_2_2.AutomationRegistryBase22OnchainConfig
	OffchainConfigVersion uint64
	OffchainConfig        []byte
}

func (e *LegacyEthereumOffchainAggregatorV2) Address() string {
	return e.address.Hex()
}

func (e *LegacyEthereumOffchainAggregatorV2) RequestNewRound() error {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := e.contract.RequestNewRound(opts)
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx)
}

func (e *LegacyEthereumOffchainAggregatorV2) GetLatestAnswer(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	return e.contract.LatestAnswer(opts)
}

func (e *LegacyEthereumOffchainAggregatorV2) GetLatestRound(ctx context.Context) (*RoundData, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	data, err := e.contract.LatestRoundData(opts)
	if err != nil {
		return nil, err
	}
	return &RoundData{
		RoundId:         data.RoundId,
		StartedAt:       data.StartedAt,
		UpdatedAt:       data.UpdatedAt,
		AnsweredInRound: data.AnsweredInRound,
		Answer:          data.Answer,
	}, nil
}

func (e *LegacyEthereumOffchainAggregatorV2) GetRound(ctx context.Context, roundID *big.Int) (*RoundData, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	data, err := e.contract.GetRoundData(opts, roundID)
	if err != nil {
		return nil, err
	}
	return &RoundData{
		RoundId:         data.RoundId,
		StartedAt:       data.StartedAt,
		UpdatedAt:       data.UpdatedAt,
		AnsweredInRound: data.AnsweredInRound,
		Answer:          data.Answer,
	}, nil
}

func (e *LegacyEthereumOffchainAggregatorV2) SetPayees(transmitters, payees []string) error {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	e.l.Info().
		Str("Transmitters", fmt.Sprintf("%v", transmitters)).
		Str("Payees", fmt.Sprintf("%v", payees)).
		Str("OCRv2 Address", e.Address()).
		Msg("Setting OCRv2 Payees")

	var addTransmitters, addrPayees []common.Address
	for _, t := range transmitters {
		addTransmitters = append(addTransmitters, common.HexToAddress(t))
	}
	for _, p := range payees {
		addrPayees = append(addrPayees, common.HexToAddress(p))
	}

	tx, err := e.contract.SetPayees(opts, addTransmitters, addrPayees)
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx)
}

func (e *LegacyEthereumOffchainAggregatorV2) SetConfig(ocrConfig *OCRv2Config) error {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	e.l.Info().
		Str("Address", e.Address()).
		Interface("Signers", ocrConfig.Signers).
		Interface("Transmitters", ocrConfig.Transmitters).
		Uint8("F", ocrConfig.F).
		Bytes("OnchainConfig", ocrConfig.OnchainConfig).
		Uint64("OffchainConfigVersion", ocrConfig.OffchainConfigVersion).
		Bytes("OffchainConfig", ocrConfig.OffchainConfig).
		Msg("Setting OCRv2 Config")
	tx, err := e.contract.SetConfig(
		opts,
		ocrConfig.Signers,
		ocrConfig.Transmitters,
		ocrConfig.F,
		ocrConfig.OnchainConfig,
		ocrConfig.OffchainConfigVersion,
		ocrConfig.OffchainConfig,
	)
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx)
}

func (e *LegacyEthereumOffchainAggregatorV2) ParseEventAnswerUpdated(log types.Log) (*ocr2aggregator.OCR2AggregatorAnswerUpdated, error) {
	return e.contract.ParseAnswerUpdated(log)
}

// EthereumKeeperRegistryCheckUpkeepGasUsageWrapper represents a gas wrapper for keeper registry
type EthereumKeeperRegistryCheckUpkeepGasUsageWrapper struct {
	address         *common.Address
	client          blockchain.EVMClient
	gasUsageWrapper *gas_wrapper.KeeperRegistryCheckUpkeepGasUsageWrapper
}

func (e *EthereumKeeperRegistryCheckUpkeepGasUsageWrapper) Address() string {
	return e.address.Hex()
}

/* Functions 1_0_0 */

type LegacyEthereumFunctionsRouter struct {
	address  common.Address
	client   blockchain.EVMClient
	instance *functions_router.FunctionsRouter
	l        zerolog.Logger
}

func (e *LegacyEthereumFunctionsRouter) Address() string {
	return e.address.Hex()
}

func (e *LegacyEthereumFunctionsRouter) CreateSubscriptionWithConsumer(consumer string) (uint64, error) {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return 0, err
	}
	tx, err := e.instance.CreateSubscriptionWithConsumer(opts, common.HexToAddress(consumer))
	if err != nil {
		return 0, err
	}
	if err := e.client.ProcessTransaction(tx); err != nil {
		return 0, err
	}
	r, err := e.client.GetTxReceipt(tx.Hash())
	if err != nil {
		return 0, err
	}
	for _, l := range r.Logs {
		e.l.Info().Interface("Log", common.Bytes2Hex(l.Data)).Send()
	}
	topicsMap := map[string]interface{}{}

	fabi, err := abi.JSON(strings.NewReader(functions_router.FunctionsRouterABI))
	if err != nil {
		return 0, err
	}
	for _, ev := range fabi.Events {
		e.l.Info().Str("EventName", ev.Name).Send()
	}
	topicOneInputs := abi.Arguments{fabi.Events["SubscriptionCreated"].Inputs[0]}
	topicOneHash := []common.Hash{r.Logs[0].Topics[1:][0]}
	if err := abi.ParseTopicsIntoMap(topicsMap, topicOneInputs, topicOneHash); err != nil {
		return 0, fmt.Errorf("failed to decode topic value, err: %w", err)
	}
	e.l.Info().Interface("NewTopicsDecoded", topicsMap).Send()
	if topicsMap["subscriptionId"] == 0 {
		return 0, fmt.Errorf("failed to decode subscription ID after creation")
	}
	return topicsMap["subscriptionId"].(uint64), nil
}

type LegacyEthereumFunctionsCoordinator struct {
	address  common.Address
	client   blockchain.EVMClient
	instance *functions_coordinator.FunctionsCoordinator
}

func (e *LegacyEthereumFunctionsCoordinator) GetThresholdPublicKey() ([]byte, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.client.GetDefaultWallet().Address()),
		Context: context.Background(),
	}
	return e.instance.GetThresholdPublicKey(opts)
}

func (e *LegacyEthereumFunctionsCoordinator) GetDONPublicKey() ([]byte, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.client.GetDefaultWallet().Address()),
		Context: context.Background(),
	}
	return e.instance.GetDONPublicKey(opts)
}

func (e *LegacyEthereumFunctionsCoordinator) Address() string {
	return e.address.Hex()
}

type LegacyEthereumFunctionsLoadTestClient struct {
	address  common.Address
	client   blockchain.EVMClient
	instance *functions_load_test_client.FunctionsLoadTestClient
}

func (e *LegacyEthereumFunctionsLoadTestClient) Address() string {
	return e.address.Hex()
}

type EthereumFunctionsLoadStats struct {
	LastRequestID string
	LastResponse  string
	LastError     string
	Total         uint32
	Succeeded     uint32
	Errored       uint32
	Empty         uint32
}

func Bytes32ToSlice(a [32]byte) (r []byte) {
	r = append(r, a[:]...)
	return
}

func (e *LegacyEthereumFunctionsLoadTestClient) GetStats() (*EthereumFunctionsLoadStats, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.client.GetDefaultWallet().Address()),
		Context: context.Background(),
	}
	lr, lbody, lerr, total, succeeded, errored, empty, err := e.instance.GetStats(opts)
	if err != nil {
		return nil, err
	}
	return &EthereumFunctionsLoadStats{
		LastRequestID: string(Bytes32ToSlice(lr)),
		LastResponse:  string(lbody),
		LastError:     string(lerr),
		Total:         total,
		Succeeded:     succeeded,
		Errored:       errored,
		Empty:         empty,
	}, nil
}

func (e *LegacyEthereumFunctionsLoadTestClient) ResetStats() error {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := e.instance.ResetStats(opts)
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx)
}

func (e *LegacyEthereumFunctionsLoadTestClient) SendRequest(times uint32, source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, jobId [32]byte) error {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := e.instance.SendRequest(opts, times, source, encryptedSecretsReferences, args, subscriptionId, jobId)
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx)
}

func (e *LegacyEthereumFunctionsLoadTestClient) SendRequestWithDONHostedSecrets(times uint32, source string, slotID uint8, slotVersion uint64, args []string, subscriptionId uint64, donID [32]byte) error {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := e.instance.SendRequestWithDONHostedSecrets(opts, times, source, slotID, slotVersion, args, subscriptionId, donID)
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx)
}

type EthereumMercuryVerifier struct {
	address  common.Address
	client   blockchain.EVMClient
	instance *verifier.Verifier
	l        zerolog.Logger
}

func (e *EthereumMercuryVerifier) Address() common.Address {
	return e.address
}

func (e *EthereumMercuryVerifier) Verify(signedReport []byte, sender common.Address) error {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := e.instance.Verify(opts, signedReport, sender)
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx)
}

func (e *EthereumMercuryVerifier) SetConfig(feedId [32]byte, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte, recipientAddressesAndWeights []verifier.CommonAddressAndWeight) (*types.Transaction, error) {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	tx, err := e.instance.SetConfig(opts, feedId, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, recipientAddressesAndWeights)
	e.l.Info().Err(err).Str("contractAddress", e.address.Hex()).Hex("feedId", feedId[:]).Msg("Called EthereumMercuryVerifier.SetConfig()")
	if err != nil {
		return nil, err
	}
	return tx, e.client.ProcessTransaction(tx)
}

func (e *EthereumMercuryVerifier) LatestConfigDetails(ctx context.Context, feedId [32]byte) (verifier.LatestConfigDetails, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	d, err := e.instance.LatestConfigDetails(opts, feedId)
	e.l.Info().Err(err).Str("contractAddress", e.address.Hex()).Hex("feedId", feedId[:]).
		Interface("details", d).
		Msg("Called EthereumMercuryVerifier.LatestConfigDetails()")
	if err != nil {
		return verifier.LatestConfigDetails{}, err
	}
	return d, nil
}

type EthereumMercuryVerifierProxy struct {
	address  common.Address
	client   blockchain.EVMClient
	instance *verifier_proxy.VerifierProxy
	l        zerolog.Logger
}

func (e *EthereumMercuryVerifierProxy) Address() common.Address {
	return e.address
}

func (e *EthereumMercuryVerifierProxy) InitializeVerifier(verifierAddress common.Address) (*types.Transaction, error) {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	tx, err := e.instance.InitializeVerifier(opts, verifierAddress)
	e.l.Info().Err(err).Str("contractAddress", e.address.Hex()).Str("verifierAddress", verifierAddress.Hex()).
		Msg("Called EthereumMercuryVerifierProxy.InitializeVerifier()")
	if err != nil {
		return nil, err
	}
	return tx, e.client.ProcessTransaction(tx)
}

func (e *EthereumMercuryVerifierProxy) Verify(signedReport []byte, parameterPayload []byte, value *big.Int) (*types.Transaction, error) {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if value != nil {
		opts.Value = value
	}
	if err != nil {
		return nil, err
	}
	tx, err := e.instance.Verify(opts, signedReport, parameterPayload)
	if err != nil {
		return nil, err
	}
	return tx, e.client.ProcessTransaction(tx)
}

func (e *EthereumMercuryVerifierProxy) VerifyBulk(signedReports [][]byte, parameterPayload []byte, value *big.Int) (*types.Transaction, error) {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if value != nil {
		opts.Value = value
	}
	if err != nil {
		return nil, err
	}
	tx, err := e.instance.VerifyBulk(opts, signedReports, parameterPayload)
	if err != nil {
		return nil, err
	}
	return tx, e.client.ProcessTransaction(tx)
}

func (e *EthereumMercuryVerifierProxy) SetFeeManager(feeManager common.Address) (*types.Transaction, error) {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	tx, err := e.instance.SetFeeManager(opts, feeManager)
	e.l.Info().Err(err).Str("feeManager", feeManager.Hex()).Msg("Called MercuryVerifierProxy.SetFeeManager()")
	if err != nil {
		return nil, err
	}
	return tx, e.client.ProcessTransaction(tx)
}

type EthereumMercuryFeeManager struct {
	address  common.Address
	client   blockchain.EVMClient
	instance *fee_manager.FeeManager
	l        zerolog.Logger
}

func (e *EthereumMercuryFeeManager) Address() common.Address {
	return e.address
}

func (e *EthereumMercuryFeeManager) UpdateSubscriberDiscount(subscriber common.Address, feedId [32]byte, token common.Address, discount uint64) (*types.Transaction, error) {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	tx, err := e.instance.UpdateSubscriberDiscount(opts, subscriber, feedId, token, discount)
	e.l.Info().Err(err).Msg("Called EthereumMercuryFeeManager.UpdateSubscriberDiscount()")
	if err != nil {
		return nil, err
	}
	return tx, e.client.ProcessTransaction(tx)
}

type EthereumMercuryRewardManager struct {
	address  common.Address
	client   blockchain.EVMClient
	instance *reward_manager.RewardManager
	l        zerolog.Logger
}

func (e *EthereumMercuryRewardManager) Address() common.Address {
	return e.address
}

func (e *EthereumMercuryRewardManager) SetFeeManager(feeManager common.Address) (*types.Transaction, error) {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	tx, err := e.instance.SetFeeManager(opts, feeManager)
	e.l.Info().Err(err).Str("feeManager", feeManager.Hex()).Msg("Called EthereumMercuryRewardManager.SetFeeManager()")
	if err != nil {
		return nil, err
	}
	return tx, e.client.ProcessTransaction(tx)
}

type EthereumWERC20Mock struct {
	address  common.Address
	client   blockchain.EVMClient
	instance *werc20_mock.WERC20Mock
	l        zerolog.Logger
}

func (e *EthereumWERC20Mock) Address() common.Address {
	return e.address
}

func (e *EthereumWERC20Mock) Approve(to string, amount *big.Int) error {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	e.l.Info().
		Str("From", e.client.GetDefaultWallet().Address()).
		Str("To", to).
		Str("Amount", amount.String()).
		Uint64("Nonce", opts.Nonce.Uint64()).
		Msg("Approving LINK Transfer")
	tx, err := e.instance.Approve(opts, common.HexToAddress(to), amount)
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx)
}

func (e *EthereumWERC20Mock) BalanceOf(ctx context.Context, addr string) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	balance, err := e.instance.BalanceOf(opts, common.HexToAddress(addr))
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (e *EthereumWERC20Mock) Transfer(to string, amount *big.Int) error {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	e.l.Info().
		Str("From", e.client.GetDefaultWallet().Address()).
		Str("To", to).
		Str("Amount", amount.String()).
		Uint64("Nonce", opts.Nonce.Uint64()).
		Msg("EthereumWERC20Mock.Transfer()")
	tx, err := e.instance.Transfer(opts, common.HexToAddress(to), amount)
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx)
}

func (e *EthereumWERC20Mock) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	opts, err := e.client.TransactionOpts(e.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	e.l.Info().
		Str("account", account.Hex()).
		Str("amount", amount.String()).
		Msg("EthereumWERC20Mock.Mint()")
	tx, err := e.instance.Mint(opts, account, amount)
	if err != nil {
		return tx, err
	}
	return tx, e.client.ProcessTransaction(tx)
}

func ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(k8sNodes []*client.ChainlinkK8sClient) []ChainlinkNodeWithKeysAndAddress {
	var nodesAsInterface = make([]ChainlinkNodeWithKeysAndAddress, len(k8sNodes))
	for i, node := range k8sNodes {
		nodesAsInterface[i] = node
	}

	return nodesAsInterface
}

func ChainlinkClientToChainlinkNodeWithKeysAndAddress(k8sNodes []*client.ChainlinkClient) []ChainlinkNodeWithKeysAndAddress {
	var nodesAsInterface = make([]ChainlinkNodeWithKeysAndAddress, len(k8sNodes))
	for i, node := range k8sNodes {
		nodesAsInterface[i] = node
	}

	return nodesAsInterface
}

func V2OffChainAgrregatorToOffChainAggregatorWithRounds(contracts []OffchainAggregatorV2) []OffChainAggregatorWithRounds {
	var contractsAsInterface = make([]OffChainAggregatorWithRounds, len(contracts))
	for i, contract := range contracts {
		contractsAsInterface[i] = contract
	}

	return contractsAsInterface
}

func V1OffChainAgrregatorToOffChainAggregatorWithRounds(contracts []OffchainAggregator) []OffChainAggregatorWithRounds {
	var contractsAsInterface = make([]OffChainAggregatorWithRounds, len(contracts))
	for i, contract := range contracts {
		contractsAsInterface[i] = contract
	}

	return contractsAsInterface
}

func GetRegistryContractABI(version contractsethereum.KeeperRegistryVersion) (*abi.ABI, error) {
	var (
		contractABI *abi.ABI
		err         error
	)
	switch version {
	case contractsethereum.RegistryVersion_1_0, contractsethereum.RegistryVersion_1_1:
		contractABI, err = keeper_registry_wrapper1_1.KeeperRegistryMetaData.GetAbi()
	case contractsethereum.RegistryVersion_1_2:
		contractABI, err = keeper_registry_wrapper1_2.KeeperRegistryMetaData.GetAbi()
	case contractsethereum.RegistryVersion_1_3:
		contractABI, err = keeper_registry_wrapper1_3.KeeperRegistryMetaData.GetAbi()
	case contractsethereum.RegistryVersion_2_0:
		contractABI, err = keeper_registry_wrapper2_0.KeeperRegistryMetaData.GetAbi()
	case contractsethereum.RegistryVersion_2_1:
		contractABI, err = iregistry21.IKeeperRegistryMasterMetaData.GetAbi()
	case contractsethereum.RegistryVersion_2_2:
		contractABI, err = iregistry22.IAutomationRegistryMasterMetaData.GetAbi()
	default:
		contractABI, err = keeper_registry_wrapper2_0.KeeperRegistryMetaData.GetAbi()
	}

	return contractABI, err
}
