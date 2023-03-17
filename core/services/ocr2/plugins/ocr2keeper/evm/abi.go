package evm

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	types2 "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/smartcontractkit/ocr2keepers/pkg/chain"
	"github.com/smartcontractkit/ocr2keepers/pkg/types"

	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper2_0"
)

type evmRegistryPackerV2_0 struct {
	abi abi.ABI
}

// enum UpkeepFailureReason
// https://github.com/smartcontractkit/chainlink/blob/d9dee8ea6af26bc82463510cb8786b951fa98585/contracts/src/v0.8/interfaces/AutomationRegistryInterface2_0.sol#L94
const (
	UPKEEP_FAILURE_REASON_NONE = iota
	UPKEEP_FAILURE_REASON_UPKEEP_CANCELLED
	UPKEEP_FAILURE_REASON_UPKEEP_PAUSED
	UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED
	UPKEEP_FAILURE_REASON_UPKEEP_NOT_NEEDED
	UPKEEP_FAILURE_REASON_PERFORM_DATA_EXCEEDS_LIMIT
	UPKEEP_FAILURE_REASON_INSUFFICIENT_BALANCE
	// below not from AutomationRegistryInterface2_0
	UPKEEP_FAILURE_REASON_OFFCHAIN_LOOKUP_ERROR
	UPKEEP_FAILURE_REASON_TARGET_PERFORM_REVERTED
)

func NewEvmRegistryPackerV2_0(abi abi.ABI) *evmRegistryPackerV2_0 {
	return &evmRegistryPackerV2_0{abi: abi}
}

func (rp *evmRegistryPackerV2_0) UnpackCheckResult(key types.UpkeepKey, raw string) (types.UpkeepResult, error) {
	b, err := hexutil.Decode(raw)
	if err != nil {
		return types.UpkeepResult{}, err
	}

	out, err := rp.abi.Methods["checkUpkeep"].Outputs.UnpackValues(b)
	if err != nil {
		return types.UpkeepResult{}, fmt.Errorf("%w: unpack checkUpkeep return: %s", err, raw)
	}

	result := types.UpkeepResult{
		Key:   key,
		State: types.Eligible,
	}

	upkeepNeeded := *abi.ConvertType(out[0], new(bool)).(*bool)
	rawPerformData := *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	result.FailureReason = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	result.GasUsed = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	result.FastGasWei = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	result.LinkNative = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	if !upkeepNeeded {
		result.State = types.NotEligible
	}
	// if NONE we expect the perform data. if TARGET_CHECK_REVERTED we will have the error data in the perform data used for off chain lookup
	if result.FailureReason == UPKEEP_FAILURE_REASON_NONE || result.FailureReason == UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED {
		var ret0 = new(performDataWrapper)
		err = pdataABI.UnpackIntoInterface(ret0, "check", rawPerformData)
		if err != nil {
			return types.UpkeepResult{}, err
		}

		result.CheckBlockNumber = ret0.Result.CheckBlockNumber
		result.CheckBlockHash = ret0.Result.CheckBlockhash
		result.PerformData = ret0.Result.PerformData
	}

	// This is a default placeholder which is used since we do not get the execute gas
	// from checkUpkeep result. This field is overwritten later from the execute gas
	// we have for an upkeep in memory. TODO (AUTO-1482): Refactor this
	result.ExecuteGas = 5_000_000

	return result, nil
}

func (rp *evmRegistryPackerV2_0) UnpackPerformResult(raw string) (bool, error) {
	b, err := hexutil.Decode(raw)
	if err != nil {
		return false, err
	}

	out, err := rp.abi.Methods["simulatePerformUpkeep"].
		Outputs.UnpackValues(b)
	if err != nil {
		return false, fmt.Errorf("%w: unpack simulatePerformUpkeep return: %s", err, raw)
	}

	return *abi.ConvertType(out[0], new(bool)).(*bool), nil
}

func (rp *evmRegistryPackerV2_0) UnpackUpkeepResult(id *big.Int, raw string) (activeUpkeep, error) {
	b, err := hexutil.Decode(raw)
	if err != nil {
		return activeUpkeep{}, err
	}

	out, err := rp.abi.Methods["getUpkeep"].Outputs.UnpackValues(b)
	if err != nil {
		return activeUpkeep{}, fmt.Errorf("%w: unpack getUpkeep return: %s", err, raw)
	}

	type upkeepInfo struct {
		Target                 common.Address
		ExecuteGas             uint32
		CheckData              []byte
		Balance                *big.Int
		Admin                  common.Address
		MaxValidBlocknumber    uint64
		LastPerformBlockNumber uint32
		AmountSpent            *big.Int
		Paused                 bool
		OffchainConfig         []byte
	}
	temp := *abi.ConvertType(out[0], new(upkeepInfo)).(*upkeepInfo)

	au := activeUpkeep{
		ID:              id,
		PerformGasLimit: temp.ExecuteGas,
		CheckData:       temp.CheckData,
	}

	return au, nil
}

func (rp *evmRegistryPackerV2_0) UnpackTransmitTxInput(raw []byte) ([]types.UpkeepResult, error) {
	out, err := rp.abi.Methods["transmit"].Inputs.UnpackValues(raw)
	if err != nil {
		return nil, fmt.Errorf("%w: unpack TransmitTxInput return: %s", err, raw)
	}

	if len(out) < 2 {
		return nil, fmt.Errorf("invalid unpacking of TransmitTxInput in %s", raw)
	}
	decodedReport, err := chain.NewEVMReportEncoder().DecodeReport(out[1].([]byte))
	if err != nil {
		return nil, fmt.Errorf("error during decoding report while unpacking TransmitTxInput: %w", err)
	}
	return decodedReport, nil
}

var (
	// rawPerformData is abi encoded tuple(uint32, bytes32, bytes). We create an ABI with dummy
	// function which returns this tuple in order to decode the bytes
	pdataABI, _ = abi.JSON(strings.NewReader(`[{
		"name":"check",
		"type":"function",
		"outputs":[{
			"name":"ret",
			"type":"tuple",
			"components":[
				{"type":"uint32","name":"checkBlockNumber"},
				{"type":"bytes32","name":"checkBlockhash"},
				{"type":"bytes","name":"performData"}
				]
			}]
		}]`,
	))
)

type performDataWrapper struct {
	Result performDataStruct
}
type performDataStruct struct {
	CheckBlockNumber uint32   `abi:"checkBlockNumber"`
	CheckBlockhash   [32]byte `abi:"checkBlockhash"`
	PerformData      []byte   `abi:"performData"`
}

// KeeperRegistryInterface copied from wrapper
//
//go:generate mockery --quiet --name KeeperRegistryInterface --output ./mocks/ --case=underscore
type KeeperRegistryInterface interface {
	GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)
	GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error)
	GetKeeperRegistryLogicAddress(opts *bind.CallOpts) (common.Address, error)
	GetLinkAddress(opts *bind.CallOpts) (common.Address, error)
	GetLinkNativeFeedAddress(opts *bind.CallOpts) (common.Address, error)
	GetMaxPaymentForGas(opts *bind.CallOpts, gasLimit uint32) (*big.Int, error)
	GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error)
	// GetMode(opts *bind.CallOpts) (uint8, error)
	GetPeerRegistryMigrationPermission(opts *bind.CallOpts, peer common.Address) (uint8, error)
	GetSignerInfo(opts *bind.CallOpts, query common.Address) (keeper_registry_wrapper2_0.GetSignerInfo, error)
	GetState(opts *bind.CallOpts) (keeper_registry_wrapper2_0.GetState, error)
	GetTransmitterInfo(opts *bind.CallOpts, query common.Address) (keeper_registry_wrapper2_0.GetTransmitterInfo, error)
	GetUpkeep(opts *bind.CallOpts, id *big.Int) (keeper_registry_wrapper2_0.UpkeepInfo, error)
	LatestConfigDetails(opts *bind.CallOpts) (keeper_registry_wrapper2_0.LatestConfigDetails, error)
	LatestConfigDigestAndEpoch(opts *bind.CallOpts) (keeper_registry_wrapper2_0.LatestConfigDigestAndEpoch, error)
	Owner(opts *bind.CallOpts) (common.Address, error)
	TypeAndVersion(opts *bind.CallOpts) (string, error)
	UpkeepTranscoderVersion(opts *bind.CallOpts) (uint8, error)
	UpkeepVersion(opts *bind.CallOpts) (uint8, error)
	AcceptOwnership(opts *bind.TransactOpts) (*types2.Transaction, error)
	AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types2.Transaction, error)
	AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types2.Transaction, error)
	AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types2.Transaction, error)
	CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types2.Transaction, error)
	CheckUpkeep(opts *bind.TransactOpts, id *big.Int) (*types2.Transaction, error)
	MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types2.Transaction, error)
	OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types2.Transaction, error)
	Pause(opts *bind.TransactOpts) (*types2.Transaction, error)
	PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types2.Transaction, error)
	ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types2.Transaction, error)
	RecoverFunds(opts *bind.TransactOpts) (*types2.Transaction, error)
	RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types2.Transaction, error)
	SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types2.Transaction, error)
	SetPayees(opts *bind.TransactOpts, payees []common.Address) (*types2.Transaction, error)
	SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types2.Transaction, error)
	SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types2.Transaction, error)
	SetUpkeepOffchainConfig(opts *bind.TransactOpts, id *big.Int, config []byte) (*types2.Transaction, error)
	SimulatePerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types2.Transaction, error)
	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types2.Transaction, error)
	TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types2.Transaction, error)
	TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types2.Transaction, error)
	Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types2.Transaction, error)
	Unpause(opts *bind.TransactOpts) (*types2.Transaction, error)
	UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types2.Transaction, error)
	UpdateCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types2.Transaction, error)
	WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types2.Transaction, error)
	WithdrawOwnerFunds(opts *bind.TransactOpts) (*types2.Transaction, error)
	WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types2.Transaction, error)
	Fallback(opts *bind.TransactOpts, calldata []byte) (*types2.Transaction, error)
	Receive(opts *bind.TransactOpts) (*types2.Transaction, error)
	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*keeper_registry_wrapper2_0.KeeperRegistryCancelledUpkeepReportIterator, error)
	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryCancelledUpkeepReport, id []*big.Int) (event.Subscription, error)
	ParseCancelledUpkeepReport(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryCancelledUpkeepReport, error)
	FilterConfigSet(opts *bind.FilterOpts) (*keeper_registry_wrapper2_0.KeeperRegistryConfigSetIterator, error)
	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryConfigSet) (event.Subscription, error)
	ParseConfigSet(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryConfigSet, error)
	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*keeper_registry_wrapper2_0.KeeperRegistryFundsAddedIterator, error)
	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)
	ParseFundsAdded(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryFundsAdded, error)
	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*keeper_registry_wrapper2_0.KeeperRegistryFundsWithdrawnIterator, error)
	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryFundsWithdrawn, id []*big.Int) (event.Subscription, error)
	ParseFundsWithdrawn(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryFundsWithdrawn, error)
	FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*keeper_registry_wrapper2_0.KeeperRegistryInsufficientFundsUpkeepReportIterator, error)
	WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error)
	ParseInsufficientFundsUpkeepReport(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryInsufficientFundsUpkeepReport, error)
	FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*keeper_registry_wrapper2_0.KeeperRegistryOwnerFundsWithdrawnIterator, error)
	WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryOwnerFundsWithdrawn) (event.Subscription, error)
	ParseOwnerFundsWithdrawn(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryOwnerFundsWithdrawn, error)
	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*keeper_registry_wrapper2_0.KeeperRegistryOwnershipTransferRequestedIterator, error)
	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)
	ParseOwnershipTransferRequested(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryOwnershipTransferRequested, error)
	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*keeper_registry_wrapper2_0.KeeperRegistryOwnershipTransferredIterator, error)
	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)
	ParseOwnershipTransferred(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryOwnershipTransferred, error)
	FilterPaused(opts *bind.FilterOpts) (*keeper_registry_wrapper2_0.KeeperRegistryPausedIterator, error)
	WatchPaused(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryPaused) (event.Subscription, error)
	ParsePaused(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryPaused, error)
	FilterPayeesUpdated(opts *bind.FilterOpts) (*keeper_registry_wrapper2_0.KeeperRegistryPayeesUpdatedIterator, error)
	WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryPayeesUpdated) (event.Subscription, error)
	ParsePayeesUpdated(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryPayeesUpdated, error)
	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*keeper_registry_wrapper2_0.KeeperRegistryPayeeshipTransferRequestedIterator, error)
	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)
	ParsePayeeshipTransferRequested(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryPayeeshipTransferRequested, error)
	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*keeper_registry_wrapper2_0.KeeperRegistryPayeeshipTransferredIterator, error)
	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)
	ParsePayeeshipTransferred(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryPayeeshipTransferred, error)
	FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*keeper_registry_wrapper2_0.KeeperRegistryPaymentWithdrawnIterator, error)
	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)
	ParsePaymentWithdrawn(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryPaymentWithdrawn, error)
	FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*keeper_registry_wrapper2_0.KeeperRegistryReorgedUpkeepReportIterator, error)
	WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryReorgedUpkeepReport, id []*big.Int) (event.Subscription, error)
	ParseReorgedUpkeepReport(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryReorgedUpkeepReport, error)
	FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*keeper_registry_wrapper2_0.KeeperRegistryStaleUpkeepReportIterator, error)
	WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryStaleUpkeepReport, id []*big.Int) (event.Subscription, error)
	ParseStaleUpkeepReport(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryStaleUpkeepReport, error)
	FilterTransmitted(opts *bind.FilterOpts) (*keeper_registry_wrapper2_0.KeeperRegistryTransmittedIterator, error)
	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryTransmitted) (event.Subscription, error)
	ParseTransmitted(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryTransmitted, error)
	FilterUnpaused(opts *bind.FilterOpts) (*keeper_registry_wrapper2_0.KeeperRegistryUnpausedIterator, error)
	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryUnpaused) (event.Subscription, error)
	ParseUnpaused(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryUnpaused, error)
	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*keeper_registry_wrapper2_0.KeeperRegistryUpkeepAdminTransferRequestedIterator, error)
	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)
	ParseUpkeepAdminTransferRequested(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryUpkeepAdminTransferRequested, error)
	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*keeper_registry_wrapper2_0.KeeperRegistryUpkeepAdminTransferredIterator, error)
	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)
	ParseUpkeepAdminTransferred(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryUpkeepAdminTransferred, error)
	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*keeper_registry_wrapper2_0.KeeperRegistryUpkeepCanceledIterator, error)
	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)
	ParseUpkeepCanceled(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryUpkeepCanceled, error)
	FilterUpkeepCheckDataUpdated(opts *bind.FilterOpts, id []*big.Int) (*keeper_registry_wrapper2_0.KeeperRegistryUpkeepCheckDataUpdatedIterator, error)
	WatchUpkeepCheckDataUpdated(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryUpkeepCheckDataUpdated, id []*big.Int) (event.Subscription, error)
	ParseUpkeepCheckDataUpdated(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryUpkeepCheckDataUpdated, error)
	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*keeper_registry_wrapper2_0.KeeperRegistryUpkeepGasLimitSetIterator, error)
	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)
	ParseUpkeepGasLimitSet(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryUpkeepGasLimitSet, error)
	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*keeper_registry_wrapper2_0.KeeperRegistryUpkeepMigratedIterator, error)
	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryUpkeepMigrated, id []*big.Int) (event.Subscription, error)
	ParseUpkeepMigrated(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryUpkeepMigrated, error)
	FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*keeper_registry_wrapper2_0.KeeperRegistryUpkeepOffchainConfigSetIterator, error)
	WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error)
	ParseUpkeepOffchainConfigSet(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryUpkeepOffchainConfigSet, error)
	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*keeper_registry_wrapper2_0.KeeperRegistryUpkeepPausedIterator, error)
	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryUpkeepPaused, id []*big.Int) (event.Subscription, error)
	ParseUpkeepPaused(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryUpkeepPaused, error)
	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*keeper_registry_wrapper2_0.KeeperRegistryUpkeepPerformedIterator, error)
	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error)
	ParseUpkeepPerformed(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryUpkeepPerformed, error)
	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*keeper_registry_wrapper2_0.KeeperRegistryUpkeepReceivedIterator, error)
	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryUpkeepReceived, id []*big.Int) (event.Subscription, error)
	ParseUpkeepReceived(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryUpkeepReceived, error)
	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*keeper_registry_wrapper2_0.KeeperRegistryUpkeepRegisteredIterator, error)
	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryUpkeepRegistered, id []*big.Int) (event.Subscription, error)
	ParseUpkeepRegistered(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryUpkeepRegistered, error)
	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*keeper_registry_wrapper2_0.KeeperRegistryUpkeepUnpausedIterator, error)
	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *keeper_registry_wrapper2_0.KeeperRegistryUpkeepUnpaused, id []*big.Int) (event.Subscription, error)
	ParseUpkeepUnpaused(log types2.Log) (*keeper_registry_wrapper2_0.KeeperRegistryUpkeepUnpaused, error)
	ParseLog(log types2.Log) (generated.AbigenLog, error)
	Address() common.Address
}
