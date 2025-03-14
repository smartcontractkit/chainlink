// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package automation_registry_logic_a_wrapper_2_2

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

var AutomationRegistryLogicA22MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"logicB\",\"type\":\"address\",\"internalType\":\"contractAutomationRegistryLogicB2_2\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addFunds\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"checkCallback\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"values\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"upkeepNeeded\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"performData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"upkeepFailureReason\",\"type\":\"uint8\",\"internalType\":\"enumAutomationRegistryBase2_2.UpkeepFailureReason\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"checkUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"triggerData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"upkeepNeeded\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"performData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"upkeepFailureReason\",\"type\":\"uint8\",\"internalType\":\"enumAutomationRegistryBase2_2.UpkeepFailureReason\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fastGasWei\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"linkNative\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"checkUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"upkeepNeeded\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"performData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"upkeepFailureReason\",\"type\":\"uint8\",\"internalType\":\"enumAutomationRegistryBase2_2.UpkeepFailureReason\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fastGasWei\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"linkNative\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"executeCallback\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"payload\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"upkeepNeeded\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"performData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"upkeepFailureReason\",\"type\":\"uint8\",\"internalType\":\"enumAutomationRegistryBase2_2.UpkeepFailureReason\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"fallbackTo\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"migrateUpkeeps\",\"inputs\":[{\"name\":\"ids\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"destination\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"receiveUpkeeps\",\"inputs\":[{\"name\":\"encodedUpkeeps\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerUpkeep\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"triggerType\",\"type\":\"uint8\",\"internalType\":\"enumAutomationRegistryBase2_2.Trigger\"},{\"name\":\"checkData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"triggerConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerUpkeep\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"checkData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUpkeepTriggerConfig\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"triggerConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AdminPrivilegeConfigSet\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"privilegeConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CancelledUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainSpecificModuleUpdated\",\"inputs\":[{\"name\":\"newModule\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DedupKeyAdded\",\"inputs\":[{\"name\":\"dedupKey\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsAdded\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsWithdrawn\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InsufficientFundsUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnerFundsWithdrawn\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeesUpdated\",\"inputs\":[{\"name\":\"transmitters\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"payees\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeeshipTransferRequested\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeeshipTransferred\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PaymentWithdrawn\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"payee\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReorgedUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StaleUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepAdminTransferRequested\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepAdminTransferred\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepCanceled\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"atBlockHeight\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepCheckDataSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"newCheckData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepGasLimitSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"gasLimit\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepMigrated\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"remainingBalance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"destination\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepOffchainConfigSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPaused\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPerformed\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"success\",\"type\":\"bool\",\"indexed\":true,\"internalType\":\"bool\"},{\"name\":\"totalPayment\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"gasOverhead\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPrivilegeConfigSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"privilegeConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepReceived\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"startingBalance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"importedFrom\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepRegistered\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"performGas\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"admin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepTriggerConfigSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"triggerConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepUnpaused\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ArrayHasNoEntries\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotCancel\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CheckDataExceedsLimit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConfigDigestMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateEntry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GasLimitCanOnlyIncrease\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GasLimitOutsideRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfFaultyOracles\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfSignatures\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IndexOutOfRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidDataLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRecipient\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReport\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSigner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTransmitter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTrigger\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTriggerType\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MigrationNotPermitted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAContract\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyActiveSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyActiveTransmitters\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByLINKToken\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwnerOrAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByProposedAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByProposedPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyPausedUpkeep\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlySimulatedBackend\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyUnpausedUpkeep\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterLengthError\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PaymentGreaterThanAllLINK\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RegistryPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RepeatedSigner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RepeatedTransmitter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TargetCheckReverted\",\"inputs\":[{\"name\":\"reason\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"TooManyOracles\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TranscoderNotSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepAlreadyExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepCancelled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepNotCanceled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepNotNeeded\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ValueNotChanged\",\"inputs\":[]}]",
	Bin: "0x6101406040523480156200001257600080fd5b5060405162005fd938038062005fd98339810160408190526200003591620003b1565b80816001600160a01b031663ca30e6036040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200009b9190620003b1565b826001600160a01b031663b10b673c6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000da573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001009190620003b1565b836001600160a01b0316636709d0e56040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200013f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001659190620003b1565b846001600160a01b0316635425d8ac6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca9190620003b1565b856001600160a01b031663a08714c06040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000209573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200022f9190620003b1565b3380600081620002865760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620002b957620002b981620002ed565b5050506001600160a01b0394851660805292841660a05290831660c052821660e052811661010052166101205250620003d8565b336001600160a01b03821603620003475760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200027d565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b0381168114620003ae57600080fd5b50565b600060208284031215620003c457600080fd5b8151620003d18162000398565b9392505050565b60805160a05160c05160e0516101005161012051615b9c6200043d6000396000818161010001526101c5015260006132180152600081816104a501526120c201526000613401015260006134e5015260008181611f0401526124d00152615b9c6000f3fe608060405260043610620000fe5760003560e01c806385c1b0ba1162000097578063c80480221162000061578063c80480221462000343578063ce7dc5b41462000368578063f2fde38b146200038d578063f7d334ba14620003b257620000fe565b806385c1b0ba14620002a75780638da5cb5b14620002cc5780638e86139b14620002f9578063948108f7146200031e57620000fe565b80634ee88d3511620000d95780634ee88d35146200020b5780636ded9eae146200023057806371791aa0146200025557806379ba5097146200028f57620000fe565b806328f32f38146200014657806329c5efad146200017e578063349e8cca14620001b5575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e8080156200013f573d6000f35b3d6000fd5b005b3480156200015357600080fd5b506200016b6200016536600462004137565b620003d7565b6040519081526020015b60405180910390f35b3480156200018b57600080fd5b50620001a36200019d3660046200421d565b62000751565b60405162000175949392919062004345565b348015620001c257600080fd5b507f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200162000175565b3480156200021857600080fd5b50620001446200022a36600462004382565b620009f5565b3480156200023d57600080fd5b506200016b6200024f366004620043d2565b62000a5d565b3480156200026257600080fd5b506200027a620002743660046200421d565b62000ac3565b60405162000175979695949392919062004485565b3480156200029c57600080fd5b506200014462001211565b348015620002b457600080fd5b5062000144620002c6366004620044d7565b62001314565b348015620002d957600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff16620001e5565b3480156200030657600080fd5b50620001446200031836600462004564565b62001f85565b3480156200032b57600080fd5b50620001446200033d366004620045c7565b6200230d565b3480156200035057600080fd5b506200014462000362366004620045f6565b620025a0565b3480156200037557600080fd5b50620001a362000387366004620046cc565b62002a19565b3480156200039a57600080fd5b5062000144620003ac36600462004743565b62002adf565b348015620003bf57600080fd5b506200027a620003d1366004620045f6565b62002af7565b6000805473ffffffffffffffffffffffffffffffffffffffff1633148015906200040b57506200040960093362002b35565b155b1562000443576040517fd48b678b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff89163b62000492576040517f09ee12d500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6200049d8662002b69565b9050600089307f0000000000000000000000000000000000000000000000000000000000000000604051620004d29062003ec3565b73ffffffffffffffffffffffffffffffffffffffff938416815291831660208301529091166040820152606001604051809103906000f0801580156200051c573d6000803e3d6000fd5b509050620005e3826040518060e001604052806000151581526020018c63ffffffff16815260200163ffffffff801681526020018473ffffffffffffffffffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff168152602001600063ffffffff168152508a89898080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508b92508a915062002e159050565b6015805474010000000000000000000000000000000000000000900463ffffffff16906014620006138362004792565b91906101000a81548163ffffffff021916908363ffffffff16021790555050817fbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d0128a8a6040516200068c92919063ffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a2817fcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d8787604051620006c892919062004801565b60405180910390a2817f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d56648560405162000702919062004817565b60405180910390a2817f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850846040516200073c919062004817565b60405180910390a25098975050505050505050565b600060606000806200076262003200565b600086815260046020908152604091829020825160e081018452815460ff811615158252610100810463ffffffff90811694830194909452650100000000008104841694820194909452690100000000000000000090930473ffffffffffffffffffffffffffffffffffffffff166060840152600101546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a0840152780100000000000000000000000000000000000000000000000090041660c08201525a9150600080826060015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200087c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620008a2919062004839565b73ffffffffffffffffffffffffffffffffffffffff166014600101600c9054906101000a900463ffffffff1663ffffffff1689604051620008e4919062004859565b60006040518083038160008787f1925050503d806000811462000924576040519150601f19603f3d011682016040523d82523d6000602084013e62000929565b606091505b50915091505a6200093b908562004877565b93508162000966576000604051806020016040528060008152506007965096509650505050620009ec565b808060200190518101906200097c9190620048e8565b909750955086620009aa576000604051806020016040528060008152506004965096509650505050620009ec565b601654865164010000000090910463ffffffff161015620009e8576000604051806020016040528060008152506005965096509650505050620009ec565b5050505b92959194509250565b62000a008362003272565b6000838152601b6020526040902062000a1b828483620049dd565b50827f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664838360405162000a5092919062004801565b60405180910390a2505050565b600062000ab788888860008989604051806020016040528060008152508a8a8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250620003d792505050565b98975050505050505050565b60006060600080600080600062000ad962003200565b600062000ae68a62003328565b905060006012604051806101600160405290816000820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160008201600c9054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160109054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160149054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160189054906101000a900462ffffff1662ffffff1662ffffff16815260200160008201601b9054906101000a900461ffff1661ffff1661ffff16815260200160008201601d9054906101000a900460ff1660ff1660ff16815260200160008201601e9054906101000a900460ff1615151515815260200160008201601f9054906101000a900460ff161515151581526020016001820160009054906101000a900460ff161515151581526020016001820160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152505090506000600460008d81526020019081526020016000206040518060e00160405290816000820160009054906101000a900460ff161515151581526020016000820160019054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160059054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160099054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160018201600c9054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff1681526020016001820160189054906101000a900463ffffffff1663ffffffff1663ffffffff168152505090508160e001511562000e7b576000604051806020016040528060008152506009600084602001516000808263ffffffff169250995099509950995099509950995050505062001205565b604081015163ffffffff9081161462000ecc576000604051806020016040528060008152506001600084602001516000808263ffffffff169250995099509950995099509950995050505062001205565b80511562000f12576000604051806020016040528060008152506002600084602001516000808263ffffffff169250995099509950995099509950995050505062001205565b62000f1d82620033de565b8095508196505050600062000f3a838584602001518989620035d0565b9050806bffffffffffffffffffffffff168260a001516bffffffffffffffffffffffff16101562000fa4576000604051806020016040528060008152506006600085602001516000808263ffffffff1692509a509a509a509a509a509a509a505050505062001205565b600062000fb38e868f62003885565b90505a9850600080846060015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200100b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001031919062004839565b73ffffffffffffffffffffffffffffffffffffffff166014600101600c9054906101000a900463ffffffff1663ffffffff168460405162001073919062004859565b60006040518083038160008787f1925050503d8060008114620010b3576040519150601f19603f3d011682016040523d82523d6000602084013e620010b8565b606091505b50915091505a620010ca908c62004877565b9a50816200114a5760165481516801000000000000000090910463ffffffff1610156200112757505060408051602080820190925260008082529490910151939c509a50600899505063ffffffff90911695506200120592505050565b602090940151939b5060039a505063ffffffff9092169650620012059350505050565b80806020019051810190620011609190620048e8565b909e509c508d620011a157505060408051602080820190925260008082529490910151939c509a50600499505063ffffffff90911695506200120592505050565b6016548d5164010000000090910463ffffffff161015620011f257505060408051602080820190925260008082529490910151939c509a50600599505063ffffffff90911695506200120592505050565b505050506020015163ffffffff16945050505b92959891949750929550565b60015473ffffffffffffffffffffffffffffffffffffffff16331462001298576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600173ffffffffffffffffffffffffffffffffffffffff82166000908152601a602052604090205460ff166003811115620013535762001353620042da565b141580156200139f5750600373ffffffffffffffffffffffffffffffffffffffff82166000908152601a602052604090205460ff1660038111156200139c576200139c620042da565b14155b15620013d7576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6014546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1662001437576040517fd12d7d8d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082900362001473576040517f2c2fc94100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160e081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c081018290526000808567ffffffffffffffff811115620014ca57620014ca62003fbe565b604051908082528060200260200182016040528015620014f4578160200160208202803683370190505b50905060008667ffffffffffffffff81111562001515576200151562003fbe565b6040519080825280602002602001820160405280156200159c57816020015b6040805160e08101825260008082526020808301829052928201819052606082018190526080820181905260a0820181905260c082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181620015345790505b50905060008767ffffffffffffffff811115620015bd57620015bd62003fbe565b604051908082528060200260200182016040528015620015f257816020015b6060815260200190600190039081620015dc5790505b50905060008867ffffffffffffffff81111562001613576200161362003fbe565b6040519080825280602002602001820160405280156200164857816020015b6060815260200190600190039081620016325790505b50905060008967ffffffffffffffff81111562001669576200166962003fbe565b6040519080825280602002602001820160405280156200169e57816020015b6060815260200190600190039081620016885790505b50905060005b8a81101562001c82578b8b82818110620016c257620016c262004b05565b60209081029290920135600081815260048452604090819020815160e081018352815460ff811615158252610100810463ffffffff90811697830197909752650100000000008104871693820193909352690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff166060830152600101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490931660c08401529a50909850620017a190508962003272565b60608801516040517f1a5da6c800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8c8116600483015290911690631a5da6c890602401600060405180830381600087803b1580156200181157600080fd5b505af115801562001826573d6000803e3d6000fd5b505050508785828151811062001840576200184062004b05565b6020026020010181905250600560008a815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1686828151811062001894576200189462004b05565b73ffffffffffffffffffffffffffffffffffffffff90921660209283029190910182015260008a81526007909152604090208054620018d39062004935565b80601f0160208091040260200160405190810160405280929190818152602001828054620019019062004935565b8015620019525780601f10620019265761010080835404028352916020019162001952565b820191906000526020600020905b8154815290600101906020018083116200193457829003601f168201915b50505050508482815181106200196c576200196c62004b05565b6020026020010181905250601b60008a81526020019081526020016000208054620019979062004935565b80601f0160208091040260200160405190810160405280929190818152602001828054620019c59062004935565b801562001a165780601f10620019ea5761010080835404028352916020019162001a16565b820191906000526020600020905b815481529060010190602001808311620019f857829003601f168201915b505050505083828151811062001a305762001a3062004b05565b6020026020010181905250601c60008a8152602001908152602001600020805462001a5b9062004935565b80601f016020809104026020016040519081016040528092919081815260200182805462001a899062004935565b801562001ada5780601f1062001aae5761010080835404028352916020019162001ada565b820191906000526020600020905b81548152906001019060200180831162001abc57829003601f168201915b505050505082828151811062001af45762001af462004b05565b60200260200101819052508760a001516bffffffffffffffffffffffff168762001b1f919062004b34565b60008a815260046020908152604080832080547fffffff000000000000000000000000000000000000000000000000000000000016815560010180547fffffffff000000000000000000000000000000000000000000000000000000001690556007909152812091985062001b95919062003ed1565b6000898152601b6020526040812062001bae9162003ed1565b6000898152601c6020526040812062001bc79162003ed1565b600089815260066020526040902080547fffffffffffffffffffffffff000000000000000000000000000000000000000016905562001c0860028a62003a75565b5060a0880151604080516bffffffffffffffffffffffff909216825273ffffffffffffffffffffffffffffffffffffffff8c1660208301528a917fb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff910160405180910390a28062001c798162004b4a565b915050620016a4565b508560195462001c93919062004877565b60195560008b8b868167ffffffffffffffff81111562001cb75762001cb762003fbe565b60405190808252806020026020018201604052801562001ce1578160200160208202803683370190505b508988888860405160200162001cff98979695949392919062004d11565b60405160208183030381529060405290508973ffffffffffffffffffffffffffffffffffffffff16638e86139b6014600001600c9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663c71249ab60038e73ffffffffffffffffffffffffffffffffffffffff1663aab9edd66040518163ffffffff1660e01b8152600401602060405180830381865afa15801562001dbb573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001de1919062004df0565b866040518463ffffffff1660e01b815260040162001e029392919062004e15565b600060405180830381865afa15801562001e20573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405262001e68919081019062004e3c565b6040518263ffffffff1660e01b815260040162001e86919062004817565b600060405180830381600087803b15801562001ea157600080fd5b505af115801562001eb6573d6000803e3d6000fd5b50506040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8d81166004830152602482018b90527f000000000000000000000000000000000000000000000000000000000000000016925063a9059cbb91506044016020604051808303816000875af115801562001f50573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001f76919062004e75565b50505050505050505050505050565b6002336000908152601a602052604090205460ff16600381111562001fae5762001fae620042da565b1415801562001fe457506003336000908152601a602052604090205460ff16600381111562001fe15762001fe1620042da565b14155b156200201c576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080808080808062002032888a018a62005060565b965096509650965096509650965060005b87518110156200230157600073ffffffffffffffffffffffffffffffffffffffff168782815181106200207a576200207a62004b05565b60200260200101516060015173ffffffffffffffffffffffffffffffffffffffff16036200218e57858181518110620020b757620020b762004b05565b6020026020010151307f0000000000000000000000000000000000000000000000000000000000000000604051620020ef9062003ec3565b73ffffffffffffffffffffffffffffffffffffffff938416815291831660208301529091166040820152606001604051809103906000f08015801562002139573d6000803e3d6000fd5b508782815181106200214f576200214f62004b05565b60200260200101516060019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250505b62002246888281518110620021a757620021a762004b05565b6020026020010151888381518110620021c457620021c462004b05565b6020026020010151878481518110620021e157620021e162004b05565b6020026020010151878581518110620021fe57620021fe62004b05565b60200260200101518786815181106200221b576200221b62004b05565b602002602001015187878151811062002238576200223862004b05565b602002602001015162002e15565b8781815181106200225b576200225b62004b05565b60200260200101517f74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a7188838151811062002299576200229962004b05565b602002602001015160a0015133604051620022e49291906bffffffffffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a280620022f88162004b4a565b91505062002043565b50505050505050505050565b600082815260046020908152604091829020825160e081018452815460ff81161515825263ffffffff6101008204811694830194909452650100000000008104841694820185905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004821660c082015291146200240b576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818160a001516200241d919062005191565b600084815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff93841602179055601954620024859184169062004b34565b6019556040517f23b872dd0000000000000000000000000000000000000000000000000000000081523360048201523060248201526bffffffffffffffffffffffff831660448201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906323b872dd906064016020604051808303816000875af11580156200252f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062002555919062004e75565b506040516bffffffffffffffffffffffff83168152339084907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a3505050565b6000818152600460209081526040808320815160e081018352815460ff81161515825263ffffffff610100820481169583019590955265010000000000810485169382019390935273ffffffffffffffffffffffffffffffffffffffff6901000000000000000000909304929092166060830152600101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490911660c0820152906200268860005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161490506000601260010160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200272b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620027519190620051b9565b9050826040015163ffffffff1660000362002798576040517ffbc0357800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604083015163ffffffff90811614620027dd576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8115801562002810575060008481526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314155b1562002848576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b816200285e576200285b60328262004b34565b90505b6000848152600460205260409020805463ffffffff80841665010000000000027fffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffffff90921691909117909155620028ba90600290869062003a7516565b5060145460808401516bffffffffffffffffffffffff918216916000911682111562002923576080850151620028f19083620051d3565b90508460a001516bffffffffffffffffffffffff16816bffffffffffffffffffffffff16111562002923575060a08401515b808560a00151620029359190620051d3565b600087815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff938416021790556015546200299d9183911662005191565b601580547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff9290921691909117905560405167ffffffffffffffff84169087907f91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f79118190600090a3505050505050565b600060606000806000634b56a42e60e01b88888860405160240162002a4193929190620051fb565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152905062002acc898262000751565b929c919b50995090975095505050505050565b62002ae962003a83565b62002af48162003b06565b50565b60006060600080600080600062002b1e886040518060200160405280600081525062000ac3565b959e949d50929b5090995097509550909350915050565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415155b90505b92915050565b6000806000601260010160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905060008173ffffffffffffffffffffffffffffffffffffffff166385df51fd60018473ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa15801562002c02573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062002c289190620051b9565b62002c34919062004877565b6040518263ffffffff1660e01b815260040162002c5391815260200190565b602060405180830381865afa15801562002c71573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062002c979190620051b9565b601554604080516020810193909352309083015274010000000000000000000000000000000000000000900463ffffffff166060820152608001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201209083015201604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052905060045b600f81101562002da3578382828151811062002d5f5762002d5f62004b05565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508062002d9a8162004b4a565b91505062002d3f565b5084600181111562002db95762002db9620042da565b60f81b81600f8151811062002dd25762002dd262004b05565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535062002e0c816200522f565b95945050505050565b6012547e01000000000000000000000000000000000000000000000000000000000000900460ff161562002e75576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601654835163ffffffff909116101562002ebb576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108fc856020015163ffffffff16108062002ef95750601554602086015163ffffffff70010000000000000000000000000000000090920482169116115b1562002f31576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000868152600460205260409020546901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff161562002f9b576040517f6e3b930b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000868152600460209081526040808320885181548a8501518b85015160608d01517fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000009093169315157fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000ff169390931761010063ffffffff92831602177fffffff000000000000000000000000000000000000000000000000ffffffffff1665010000000000938216939093027fffffff0000000000000000000000000000000000000000ffffffffffffffffff1692909217690100000000000000000073ffffffffffffffffffffffffffffffffffffffff9283160217835560808b01516001909301805460a08d015160c08e01516bffffffffffffffffffffffff9687167fffffffffffffffff000000000000000000000000000000000000000000000000909316929092176c010000000000000000000000009690911695909502949094177fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff1678010000000000000000000000000000000000000000000000009490931693909302919091179091556005835281842080547fffffffffffffffffffffffff000000000000000000000000000000000000000016918916919091179055600790915290206200318e848262005272565b508460a001516bffffffffffffffffffffffff16601954620031b1919062004b34565b6019556000868152601b60205260409020620031ce838262005272565b506000868152601c60205260409020620031e9828262005272565b50620031f760028762003bfd565b50505050505050565b3273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161462003270576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60008181526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314620032d0576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526004602052604090205465010000000000900463ffffffff9081161462002af4576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000818160045b600f811015620033bd577fff00000000000000000000000000000000000000000000000000000000000000821683826020811062003371576200337162004b05565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614620033a857506000949350505050565b80620033b48162004b4a565b9150506200332f565b5081600f1a6001811115620033d657620033d6620042da565b949350505050565b6000806000836080015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa1580156200346b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620034919190620053b4565b5094509092505050600081131580620034a957508142105b80620034ce5750828015620034ce5750620034c5824262004877565b8463ffffffff16105b15620034df576017549550620034e3565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa1580156200354f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620035759190620053b4565b50945090925050506000811315806200358d57508142105b80620035b25750828015620035b25750620035a9824262004877565b8463ffffffff16105b15620035c3576018549450620035c7565b8094505b50505050915091565b60008080866001811115620035e957620035e9620042da565b03620035f9575061ea6062003653565b6001866001811115620036105762003610620042da565b0362003621575062014c0862003653565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008760c00151600162003668919062005409565b620036789060ff16604062005425565b60165462003698906103a490640100000000900463ffffffff1662004b34565b620036a4919062004b34565b601354604080517fde9ee35e00000000000000000000000000000000000000000000000000000000815281519394506000938493610100900473ffffffffffffffffffffffffffffffffffffffff169263de9ee35e92600480820193918290030181865afa1580156200371b573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200374191906200543f565b909250905081836200375583601862004b34565b62003761919062005425565b60c08c01516200377390600162005409565b620037849060ff166115e062005425565b62003790919062004b34565b6200379c919062004b34565b620037a8908562004b34565b935060008a610140015173ffffffffffffffffffffffffffffffffffffffff166312544140856040518263ffffffff1660e01b8152600401620037ed91815260200190565b602060405180830381865afa1580156200380b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620038319190620051b9565b8b60a0015161ffff1662003846919062005425565b9050600080620038638d8c63ffffffff1689868e8e600062003c0b565b909250905062003874818362005191565b9d9c50505050505050505050505050565b606060008360018111156200389e576200389e620042da565b036200396b576000848152600760205260409081902090517f6e04ff0d0000000000000000000000000000000000000000000000000000000091620038e69160240162005507565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152905062003a6e565b6001836001811115620039825762003982620042da565b036200362157600082806020019051810190620039a091906200557e565b6000868152600760205260409081902090519192507f40691db40000000000000000000000000000000000000000000000000000000091620039e791849160240162005692565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152915062003a6e9050565b9392505050565b600062002b60838362003d66565b60005473ffffffffffffffffffffffffffffffffffffffff16331462003270576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016200128f565b3373ffffffffffffffffffffffffffffffffffffffff82160362003b87576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200128f565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600062002b60838362003e71565b60008060008960a0015161ffff168662003c26919062005425565b905083801562003c355750803a105b1562003c3e57503a5b6000858862003c4e8b8d62004b34565b62003c5a908562005425565b62003c66919062004b34565b62003c7a90670de0b6b3a764000062005425565b62003c8691906200575a565b905060008b6040015163ffffffff1664e8d4a5100062003ca7919062005425565b60208d0151889063ffffffff168b62003cc18f8862005425565b62003ccd919062004b34565b62003cdd90633b9aca0062005425565b62003ce9919062005425565b62003cf591906200575a565b62003d01919062004b34565b90506b033b2e3c9fd0803ce800000062003d1c828462004b34565b111562003d55576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909b909a5098505050505050505050565b6000818152600183016020526040812054801562003e5f57600062003d8d60018362004877565b855490915060009062003da39060019062004877565b905081811462003e0f57600086600001828154811062003dc75762003dc762004b05565b906000526020600020015490508087600001848154811062003ded5762003ded62004b05565b6000918252602080832090910192909255918252600188019052604090208390555b855486908062003e235762003e2362005796565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505062002b63565b600091505062002b63565b5092915050565b600081815260018301602052604081205462003eba5750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915562002b63565b50600062002b63565b6103ca80620057c683390190565b50805462003edf9062004935565b6000825580601f1062003ef0575050565b601f01602090049060005260206000209081019062002af491905b8082111562003f21576000815560010162003f0b565b5090565b73ffffffffffffffffffffffffffffffffffffffff8116811462002af457600080fd5b803563ffffffff8116811462003f5d57600080fd5b919050565b80356002811062003f5d57600080fd5b60008083601f84011262003f8557600080fd5b50813567ffffffffffffffff81111562003f9e57600080fd5b60208301915083602082850101111562003fb757600080fd5b9250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160e0810167ffffffffffffffff8111828210171562004013576200401362003fbe565b60405290565b604051610100810167ffffffffffffffff8111828210171562004013576200401362003fbe565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156200408a576200408a62003fbe565b604052919050565b600067ffffffffffffffff821115620040af57620040af62003fbe565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f830112620040ed57600080fd5b813562004104620040fe8262004092565b62004040565b8181528460208386010111156200411a57600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060008060e0898b0312156200415457600080fd5b8835620041618162003f25565b97506200417160208a0162003f48565b96506040890135620041838162003f25565b95506200419360608a0162003f62565b9450608089013567ffffffffffffffff80821115620041b157600080fd5b620041bf8c838d0162003f72565b909650945060a08b0135915080821115620041d957600080fd5b620041e78c838d01620040db565b935060c08b0135915080821115620041fe57600080fd5b506200420d8b828c01620040db565b9150509295985092959890939650565b600080604083850312156200423157600080fd5b82359150602083013567ffffffffffffffff8111156200425057600080fd5b6200425e85828601620040db565b9150509250929050565b60005b83811015620042855781810151838201526020016200426b565b50506000910152565b60008151808452620042a881602086016020860162004268565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600a811062004341577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b9052565b84151581526080602082015260006200436260808301866200428e565b905062004373604083018562004309565b82606083015295945050505050565b6000806000604084860312156200439857600080fd5b83359250602084013567ffffffffffffffff811115620043b757600080fd5b620043c58682870162003f72565b9497909650939450505050565b600080600080600080600060a0888a031215620043ee57600080fd5b8735620043fb8162003f25565b96506200440b6020890162003f48565b955060408801356200441d8162003f25565b9450606088013567ffffffffffffffff808211156200443b57600080fd5b620044498b838c0162003f72565b909650945060808a01359150808211156200446357600080fd5b50620044728a828b0162003f72565b989b979a50959850939692959293505050565b871515815260e060208201526000620044a260e08301896200428e565b9050620044b3604083018862004309565b8560608301528460808301528360a08301528260c083015298975050505050505050565b600080600060408486031215620044ed57600080fd5b833567ffffffffffffffff808211156200450657600080fd5b818601915086601f8301126200451b57600080fd5b8135818111156200452b57600080fd5b8760208260051b85010111156200454157600080fd5b60209283019550935050840135620045598162003f25565b809150509250925092565b600080602083850312156200457857600080fd5b823567ffffffffffffffff8111156200459057600080fd5b6200459e8582860162003f72565b90969095509350505050565b80356bffffffffffffffffffffffff8116811462003f5d57600080fd5b60008060408385031215620045db57600080fd5b82359150620045ed60208401620045aa565b90509250929050565b6000602082840312156200460957600080fd5b5035919050565b600067ffffffffffffffff8211156200462d576200462d62003fbe565b5060051b60200190565b600082601f8301126200464957600080fd5b813560206200465c620040fe8362004610565b82815260059290921b840181019181810190868411156200467c57600080fd5b8286015b84811015620046c157803567ffffffffffffffff811115620046a25760008081fd5b620046b28986838b0101620040db565b84525091830191830162004680565b509695505050505050565b60008060008060608587031215620046e357600080fd5b84359350602085013567ffffffffffffffff808211156200470357600080fd5b620047118883890162004637565b945060408701359150808211156200472857600080fd5b50620047378782880162003f72565b95989497509550505050565b6000602082840312156200475657600080fd5b813562003a6e8162003f25565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600063ffffffff808316818103620047ae57620047ae62004763565b6001019392505050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b602081526000620033d6602083018486620047b8565b60208152600062002b6060208301846200428e565b805162003f5d8162003f25565b6000602082840312156200484c57600080fd5b815162003a6e8162003f25565b600082516200486d81846020870162004268565b9190910192915050565b8181038181111562002b635762002b6362004763565b801515811462002af457600080fd5b600082601f830112620048ae57600080fd5b8151620048bf620040fe8262004092565b818152846020838601011115620048d557600080fd5b620033d682602083016020870162004268565b60008060408385031215620048fc57600080fd5b825162004909816200488d565b602084015190925067ffffffffffffffff8111156200492757600080fd5b6200425e858286016200489c565b600181811c908216806200494a57607f821691505b60208210810362004984577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f821115620049d857600081815260208120601f850160051c81016020861015620049b35750805b601f850160051c820191505b81811015620049d457828155600101620049bf565b5050505b505050565b67ffffffffffffffff831115620049f857620049f862003fbe565b62004a108362004a09835462004935565b836200498a565b6000601f84116001811462004a65576000851562004a2e5750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b17835562004afe565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b8281101562004ab6578685013582556020948501946001909201910162004a94565b508682101562004af2577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b8082018082111562002b635762002b6362004763565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820362004b7e5762004b7e62004763565b5060010190565b600081518084526020808501945080840160005b8381101562004c445781518051151588528381015163ffffffff908116858a01526040808301519091169089015260608082015173ffffffffffffffffffffffffffffffffffffffff16908901526080808201516bffffffffffffffffffffffff169089015260a08082015162004c1f828b01826bffffffffffffffffffffffff169052565b505060c09081015163ffffffff169088015260e0909601959082019060010162004b99565b509495945050505050565b600081518084526020808501945080840160005b8381101562004c4457815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010162004c63565b600082825180855260208086019550808260051b84010181860160005b8481101562004d04577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe086840301895262004cf18383516200428e565b9884019892509083019060010162004cb4565b5090979650505050505050565b60e081528760e082015260006101007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8a111562004d4e57600080fd5b8960051b808c8386013783018381038201602085015262004d728282018b62004b85565b915050828103604084015262004d89818962004c4f565b9050828103606084015262004d9f818862004c4f565b9050828103608084015262004db5818762004c97565b905082810360a084015262004dcb818662004c97565b905082810360c084015262004de1818562004c97565b9b9a5050505050505050505050565b60006020828403121562004e0357600080fd5b815160ff8116811462003a6e57600080fd5b60ff8416815260ff8316602082015260606040820152600062002e0c60608301846200428e565b60006020828403121562004e4f57600080fd5b815167ffffffffffffffff81111562004e6757600080fd5b620033d6848285016200489c565b60006020828403121562004e8857600080fd5b815162003a6e816200488d565b600082601f83011262004ea757600080fd5b8135602062004eba620040fe8362004610565b82815260059290921b8401810191818101908684111562004eda57600080fd5b8286015b84811015620046c1578035835291830191830162004ede565b600082601f83011262004f0957600080fd5b8135602062004f1c620040fe8362004610565b82815260e0928302850182019282820191908785111562004f3c57600080fd5b8387015b8581101562004d045781818a03121562004f5a5760008081fd5b62004f6462003fed565b813562004f71816200488d565b815262004f8082870162003f48565b86820152604062004f9381840162003f48565b9082015260608281013562004fa88162003f25565b90820152608062004fbb838201620045aa565b9082015260a062004fce838201620045aa565b9082015260c062004fe183820162003f48565b90820152845292840192810162004f40565b600082601f8301126200500557600080fd5b8135602062005018620040fe8362004610565b82815260059290921b840181019181810190868411156200503857600080fd5b8286015b84811015620046c1578035620050528162003f25565b83529183019183016200503c565b600080600080600080600060e0888a0312156200507c57600080fd5b873567ffffffffffffffff808211156200509557600080fd5b620050a38b838c0162004e95565b985060208a0135915080821115620050ba57600080fd5b620050c88b838c0162004ef7565b975060408a0135915080821115620050df57600080fd5b620050ed8b838c0162004ff3565b965060608a01359150808211156200510457600080fd5b620051128b838c0162004ff3565b955060808a01359150808211156200512957600080fd5b620051378b838c0162004637565b945060a08a01359150808211156200514e57600080fd5b6200515c8b838c0162004637565b935060c08a01359150808211156200517357600080fd5b50620051828a828b0162004637565b91505092959891949750929550565b6bffffffffffffffffffffffff81811683821601908082111562003e6a5762003e6a62004763565b600060208284031215620051cc57600080fd5b5051919050565b6bffffffffffffffffffffffff82811682821603908082111562003e6a5762003e6a62004763565b60408152600062005210604083018662004c97565b828103602084015262005225818587620047b8565b9695505050505050565b8051602080830151919081101562004984577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60209190910360031b1b16919050565b815167ffffffffffffffff8111156200528f576200528f62003fbe565b620052a781620052a0845462004935565b846200498a565b602080601f831160018114620052fd5760008415620052c65750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555620049d4565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b828110156200534c578886015182559484019460019091019084016200532b565b50858210156200538957878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b805169ffffffffffffffffffff8116811462003f5d57600080fd5b600080600080600060a08688031215620053cd57600080fd5b620053d88662005399565b9450602086015193506040860151925060608601519150620053fd6080870162005399565b90509295509295909350565b60ff818116838216019081111562002b635762002b6362004763565b808202811582820484141762002b635762002b6362004763565b600080604083850312156200545357600080fd5b505080516020909101519092909150565b60008154620054738162004935565b808552602060018381168015620054935760018114620054cc57620054fc565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b8901019550620054fc565b866000528260002060005b85811015620054f45781548a8201860152908301908401620054d7565b890184019650505b505050505092915050565b60208152600062002b60602083018462005464565b600082601f8301126200552e57600080fd5b8151602062005541620040fe8362004610565b82815260059290921b840181019181810190868411156200556157600080fd5b8286015b84811015620046c1578051835291830191830162005565565b6000602082840312156200559157600080fd5b815167ffffffffffffffff80821115620055aa57600080fd5b908301906101008286031215620055c057600080fd5b620055ca62004019565b82518152602083015160208201526040830151604082015260608301516060820152608083015160808201526200560460a084016200482c565b60a082015260c0830151828111156200561c57600080fd5b6200562a878286016200551c565b60c08301525060e0830151828111156200564357600080fd5b62005651878286016200489c565b60e08301525095945050505050565b600081518084526020808501945080840160005b8381101562004c445781518752958201959082019060010162005674565b60408152825160408201526020830151606082015260408301516080820152606083015160a0820152608083015160c082015273ffffffffffffffffffffffffffffffffffffffff60a08401511660e0820152600060c08401516101008081850152506200570561014084018262005660565b905060e08501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0848303016101208501526200574382826200428e565b915050828103602084015262002e0c818562005464565b60008262005791577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfe60c060405234801561001057600080fd5b506040516103ca3803806103ca83398101604081905261002f91610076565b600080546001600160a01b0319166001600160a01b039384161790559181166080521660a0526100b9565b80516001600160a01b038116811461007157600080fd5b919050565b60008060006060848603121561008b57600080fd5b6100948461005a565b92506100a26020850161005a565b91506100b06040850161005a565b90509250925092565b60805160a0516102e76100e36000396000603801526000818160c4015261011701526102e76000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806379188d161461007b578063f00e6a2a146100aa575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e808015610076573d6000f35b3d6000fd5b61008e6100893660046101c1565b6100ee565b6040805192151583526020830191909152015b60405180910390f35b60405173ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001681526020016100a1565b60008054819073ffffffffffffffffffffffffffffffffffffffff16331461011557600080fd5b7f00000000000000000000000000000000000000000000000000000000000000005a91505a61138881101561014957600080fd5b61138881039050856040820482031161016157600080fd5b50803b61016d57600080fd5b6000808551602087016000858af192505a610188908361029a565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080604083850312156101d457600080fd5b82359150602083013567ffffffffffffffff808211156101f357600080fd5b818501915085601f83011261020757600080fd5b81358181111561021957610219610192565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561025f5761025f610192565b8160405282815288602084870101111561027857600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b818103818111156102d4577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b9291505056fea164736f6c6343000813000aa164736f6c6343000813000a",
}

var AutomationRegistryLogicA22ABI = AutomationRegistryLogicA22MetaData.ABI

var AutomationRegistryLogicA22Bin = AutomationRegistryLogicA22MetaData.Bin

func DeployAutomationRegistryLogicA22(auth *bind.TransactOpts, backend bind.ContractBackend, logicB common.Address) (common.Address, *types.Transaction, *AutomationRegistryLogicA22, error) {
	parsed, err := AutomationRegistryLogicA22MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AutomationRegistryLogicA22Bin), backend, logicB)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AutomationRegistryLogicA22{address: address, abi: *parsed, AutomationRegistryLogicA22Caller: AutomationRegistryLogicA22Caller{contract: contract}, AutomationRegistryLogicA22Transactor: AutomationRegistryLogicA22Transactor{contract: contract}, AutomationRegistryLogicA22Filterer: AutomationRegistryLogicA22Filterer{contract: contract}}, nil
}

type AutomationRegistryLogicA22 struct {
	address common.Address
	abi     abi.ABI
	AutomationRegistryLogicA22Caller
	AutomationRegistryLogicA22Transactor
	AutomationRegistryLogicA22Filterer
}

type AutomationRegistryLogicA22Caller struct {
	contract *bind.BoundContract
}

type AutomationRegistryLogicA22Transactor struct {
	contract *bind.BoundContract
}

type AutomationRegistryLogicA22Filterer struct {
	contract *bind.BoundContract
}

type AutomationRegistryLogicA22Session struct {
	Contract     *AutomationRegistryLogicA22
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type AutomationRegistryLogicA22CallerSession struct {
	Contract *AutomationRegistryLogicA22Caller
	CallOpts bind.CallOpts
}

type AutomationRegistryLogicA22TransactorSession struct {
	Contract     *AutomationRegistryLogicA22Transactor
	TransactOpts bind.TransactOpts
}

type AutomationRegistryLogicA22Raw struct {
	Contract *AutomationRegistryLogicA22
}

type AutomationRegistryLogicA22CallerRaw struct {
	Contract *AutomationRegistryLogicA22Caller
}

type AutomationRegistryLogicA22TransactorRaw struct {
	Contract *AutomationRegistryLogicA22Transactor
}

func NewAutomationRegistryLogicA22(address common.Address, backend bind.ContractBackend) (*AutomationRegistryLogicA22, error) {
	abi, err := abi.JSON(strings.NewReader(AutomationRegistryLogicA22ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindAutomationRegistryLogicA22(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22{address: address, abi: abi, AutomationRegistryLogicA22Caller: AutomationRegistryLogicA22Caller{contract: contract}, AutomationRegistryLogicA22Transactor: AutomationRegistryLogicA22Transactor{contract: contract}, AutomationRegistryLogicA22Filterer: AutomationRegistryLogicA22Filterer{contract: contract}}, nil
}

func NewAutomationRegistryLogicA22Caller(address common.Address, caller bind.ContractCaller) (*AutomationRegistryLogicA22Caller, error) {
	contract, err := bindAutomationRegistryLogicA22(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22Caller{contract: contract}, nil
}

func NewAutomationRegistryLogicA22Transactor(address common.Address, transactor bind.ContractTransactor) (*AutomationRegistryLogicA22Transactor, error) {
	contract, err := bindAutomationRegistryLogicA22(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22Transactor{contract: contract}, nil
}

func NewAutomationRegistryLogicA22Filterer(address common.Address, filterer bind.ContractFilterer) (*AutomationRegistryLogicA22Filterer, error) {
	contract, err := bindAutomationRegistryLogicA22(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22Filterer{contract: contract}, nil
}

func bindAutomationRegistryLogicA22(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AutomationRegistryLogicA22MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationRegistryLogicA22.Contract.AutomationRegistryLogicA22Caller.contract.Call(opts, result, method, params...)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.AutomationRegistryLogicA22Transactor.contract.Transfer(opts)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.AutomationRegistryLogicA22Transactor.contract.Transact(opts, method, params...)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationRegistryLogicA22.Contract.contract.Call(opts, result, method, params...)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.contract.Transfer(opts)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.contract.Transact(opts, method, params...)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Caller) FallbackTo(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicA22.contract.Call(opts, &out, "fallbackTo")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Session) FallbackTo() (common.Address, error) {
	return _AutomationRegistryLogicA22.Contract.FallbackTo(&_AutomationRegistryLogicA22.CallOpts)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22CallerSession) FallbackTo() (common.Address, error) {
	return _AutomationRegistryLogicA22.Contract.FallbackTo(&_AutomationRegistryLogicA22.CallOpts)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicA22.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Session) Owner() (common.Address, error) {
	return _AutomationRegistryLogicA22.Contract.Owner(&_AutomationRegistryLogicA22.CallOpts)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22CallerSession) Owner() (common.Address, error) {
	return _AutomationRegistryLogicA22.Contract.Owner(&_AutomationRegistryLogicA22.CallOpts)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.contract.Transact(opts, "acceptOwnership")
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Session) AcceptOwnership() (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.AcceptOwnership(&_AutomationRegistryLogicA22.TransactOpts)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.AcceptOwnership(&_AutomationRegistryLogicA22.TransactOpts)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Transactor) AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.contract.Transact(opts, "addFunds", id, amount)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Session) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.AddFunds(&_AutomationRegistryLogicA22.TransactOpts, id, amount)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22TransactorSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.AddFunds(&_AutomationRegistryLogicA22.TransactOpts, id, amount)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Transactor) CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.contract.Transact(opts, "cancelUpkeep", id)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Session) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.CancelUpkeep(&_AutomationRegistryLogicA22.TransactOpts, id)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22TransactorSession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.CancelUpkeep(&_AutomationRegistryLogicA22.TransactOpts, id)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Transactor) CheckCallback(opts *bind.TransactOpts, id *big.Int, values [][]byte, extraData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.contract.Transact(opts, "checkCallback", id, values, extraData)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Session) CheckCallback(id *big.Int, values [][]byte, extraData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.CheckCallback(&_AutomationRegistryLogicA22.TransactOpts, id, values, extraData)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22TransactorSession) CheckCallback(id *big.Int, values [][]byte, extraData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.CheckCallback(&_AutomationRegistryLogicA22.TransactOpts, id, values, extraData)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Transactor) CheckUpkeep(opts *bind.TransactOpts, id *big.Int, triggerData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.contract.Transact(opts, "checkUpkeep", id, triggerData)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Session) CheckUpkeep(id *big.Int, triggerData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.CheckUpkeep(&_AutomationRegistryLogicA22.TransactOpts, id, triggerData)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22TransactorSession) CheckUpkeep(id *big.Int, triggerData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.CheckUpkeep(&_AutomationRegistryLogicA22.TransactOpts, id, triggerData)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Transactor) CheckUpkeep0(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.contract.Transact(opts, "checkUpkeep0", id)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Session) CheckUpkeep0(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.CheckUpkeep0(&_AutomationRegistryLogicA22.TransactOpts, id)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22TransactorSession) CheckUpkeep0(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.CheckUpkeep0(&_AutomationRegistryLogicA22.TransactOpts, id)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Transactor) ExecuteCallback(opts *bind.TransactOpts, id *big.Int, payload []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.contract.Transact(opts, "executeCallback", id, payload)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Session) ExecuteCallback(id *big.Int, payload []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.ExecuteCallback(&_AutomationRegistryLogicA22.TransactOpts, id, payload)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22TransactorSession) ExecuteCallback(id *big.Int, payload []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.ExecuteCallback(&_AutomationRegistryLogicA22.TransactOpts, id, payload)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Transactor) MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.contract.Transact(opts, "migrateUpkeeps", ids, destination)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Session) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.MigrateUpkeeps(&_AutomationRegistryLogicA22.TransactOpts, ids, destination)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22TransactorSession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.MigrateUpkeeps(&_AutomationRegistryLogicA22.TransactOpts, ids, destination)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Transactor) ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.contract.Transact(opts, "receiveUpkeeps", encodedUpkeeps)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Session) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.ReceiveUpkeeps(&_AutomationRegistryLogicA22.TransactOpts, encodedUpkeeps)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22TransactorSession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.ReceiveUpkeeps(&_AutomationRegistryLogicA22.TransactOpts, encodedUpkeeps)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Transactor) RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.contract.Transact(opts, "registerUpkeep", target, gasLimit, admin, triggerType, checkData, triggerConfig, offchainConfig)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Session) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.RegisterUpkeep(&_AutomationRegistryLogicA22.TransactOpts, target, gasLimit, admin, triggerType, checkData, triggerConfig, offchainConfig)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22TransactorSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.RegisterUpkeep(&_AutomationRegistryLogicA22.TransactOpts, target, gasLimit, admin, triggerType, checkData, triggerConfig, offchainConfig)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Transactor) RegisterUpkeep0(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.contract.Transact(opts, "registerUpkeep0", target, gasLimit, admin, checkData, offchainConfig)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Session) RegisterUpkeep0(target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.RegisterUpkeep0(&_AutomationRegistryLogicA22.TransactOpts, target, gasLimit, admin, checkData, offchainConfig)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22TransactorSession) RegisterUpkeep0(target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.RegisterUpkeep0(&_AutomationRegistryLogicA22.TransactOpts, target, gasLimit, admin, checkData, offchainConfig)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Transactor) SetUpkeepTriggerConfig(opts *bind.TransactOpts, id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.contract.Transact(opts, "setUpkeepTriggerConfig", id, triggerConfig)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Session) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.SetUpkeepTriggerConfig(&_AutomationRegistryLogicA22.TransactOpts, id, triggerConfig)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22TransactorSession) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.SetUpkeepTriggerConfig(&_AutomationRegistryLogicA22.TransactOpts, id, triggerConfig)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.contract.Transact(opts, "transferOwnership", to)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.TransferOwnership(&_AutomationRegistryLogicA22.TransactOpts, to)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.TransferOwnership(&_AutomationRegistryLogicA22.TransactOpts, to)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Transactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.contract.RawTransact(opts, calldata)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Session) Fallback(calldata []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.Fallback(&_AutomationRegistryLogicA22.TransactOpts, calldata)
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22TransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA22.Contract.Fallback(&_AutomationRegistryLogicA22.TransactOpts, calldata)
}

type AutomationRegistryLogicA22AdminPrivilegeConfigSetIterator struct {
	Event *AutomationRegistryLogicA22AdminPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22AdminPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22AdminPrivilegeConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22AdminPrivilegeConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22AdminPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22AdminPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22AdminPrivilegeConfigSet struct {
	Admin           common.Address
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*AutomationRegistryLogicA22AdminPrivilegeConfigSetIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22AdminPrivilegeConfigSetIterator{contract: _AutomationRegistryLogicA22.contract, event: "AdminPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22AdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22AdminPrivilegeConfigSet)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseAdminPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicA22AdminPrivilegeConfigSet, error) {
	event := new(AutomationRegistryLogicA22AdminPrivilegeConfigSet)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22CancelledUpkeepReportIterator struct {
	Event *AutomationRegistryLogicA22CancelledUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22CancelledUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22CancelledUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22CancelledUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22CancelledUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22CancelledUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22CancelledUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22CancelledUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22CancelledUpkeepReportIterator{contract: _AutomationRegistryLogicA22.contract, event: "CancelledUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22CancelledUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22CancelledUpkeepReport)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseCancelledUpkeepReport(log types.Log) (*AutomationRegistryLogicA22CancelledUpkeepReport, error) {
	event := new(AutomationRegistryLogicA22CancelledUpkeepReport)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22ChainSpecificModuleUpdatedIterator struct {
	Event *AutomationRegistryLogicA22ChainSpecificModuleUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22ChainSpecificModuleUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22ChainSpecificModuleUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22ChainSpecificModuleUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22ChainSpecificModuleUpdatedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22ChainSpecificModuleUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22ChainSpecificModuleUpdated struct {
	NewModule common.Address
	Raw       types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterChainSpecificModuleUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicA22ChainSpecificModuleUpdatedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "ChainSpecificModuleUpdated")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22ChainSpecificModuleUpdatedIterator{contract: _AutomationRegistryLogicA22.contract, event: "ChainSpecificModuleUpdated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchChainSpecificModuleUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22ChainSpecificModuleUpdated) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "ChainSpecificModuleUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22ChainSpecificModuleUpdated)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "ChainSpecificModuleUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseChainSpecificModuleUpdated(log types.Log) (*AutomationRegistryLogicA22ChainSpecificModuleUpdated, error) {
	event := new(AutomationRegistryLogicA22ChainSpecificModuleUpdated)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "ChainSpecificModuleUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22DedupKeyAddedIterator struct {
	Event *AutomationRegistryLogicA22DedupKeyAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22DedupKeyAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22DedupKeyAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22DedupKeyAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22DedupKeyAddedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22DedupKeyAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22DedupKeyAdded struct {
	DedupKey [32]byte
	Raw      types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*AutomationRegistryLogicA22DedupKeyAddedIterator, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22DedupKeyAddedIterator{contract: _AutomationRegistryLogicA22.contract, event: "DedupKeyAdded", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22DedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22DedupKeyAdded)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseDedupKeyAdded(log types.Log) (*AutomationRegistryLogicA22DedupKeyAdded, error) {
	event := new(AutomationRegistryLogicA22DedupKeyAdded)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22FundsAddedIterator struct {
	Event *AutomationRegistryLogicA22FundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22FundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22FundsAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22FundsAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22FundsAddedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22FundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22FundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*AutomationRegistryLogicA22FundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22FundsAddedIterator{contract: _AutomationRegistryLogicA22.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22FundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22FundsAdded)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "FundsAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseFundsAdded(log types.Log) (*AutomationRegistryLogicA22FundsAdded, error) {
	event := new(AutomationRegistryLogicA22FundsAdded)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22FundsWithdrawnIterator struct {
	Event *AutomationRegistryLogicA22FundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22FundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22FundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22FundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22FundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22FundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22FundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22FundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22FundsWithdrawnIterator{contract: _AutomationRegistryLogicA22.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22FundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22FundsWithdrawn)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseFundsWithdrawn(log types.Log) (*AutomationRegistryLogicA22FundsWithdrawn, error) {
	event := new(AutomationRegistryLogicA22FundsWithdrawn)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22InsufficientFundsUpkeepReportIterator struct {
	Event *AutomationRegistryLogicA22InsufficientFundsUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22InsufficientFundsUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22InsufficientFundsUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22InsufficientFundsUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22InsufficientFundsUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22InsufficientFundsUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22InsufficientFundsUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22InsufficientFundsUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22InsufficientFundsUpkeepReportIterator{contract: _AutomationRegistryLogicA22.contract, event: "InsufficientFundsUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22InsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22InsufficientFundsUpkeepReport)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseInsufficientFundsUpkeepReport(log types.Log) (*AutomationRegistryLogicA22InsufficientFundsUpkeepReport, error) {
	event := new(AutomationRegistryLogicA22InsufficientFundsUpkeepReport)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22OwnerFundsWithdrawnIterator struct {
	Event *AutomationRegistryLogicA22OwnerFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22OwnerFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22OwnerFundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22OwnerFundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22OwnerFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22OwnerFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22OwnerFundsWithdrawn struct {
	Amount *big.Int
	Raw    types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*AutomationRegistryLogicA22OwnerFundsWithdrawnIterator, error) {

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22OwnerFundsWithdrawnIterator{contract: _AutomationRegistryLogicA22.contract, event: "OwnerFundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22OwnerFundsWithdrawn) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22OwnerFundsWithdrawn)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseOwnerFundsWithdrawn(log types.Log) (*AutomationRegistryLogicA22OwnerFundsWithdrawn, error) {
	event := new(AutomationRegistryLogicA22OwnerFundsWithdrawn)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22OwnershipTransferRequestedIterator struct {
	Event *AutomationRegistryLogicA22OwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22OwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22OwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22OwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicA22OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22OwnershipTransferRequestedIterator{contract: _AutomationRegistryLogicA22.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22OwnershipTransferRequested)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseOwnershipTransferRequested(log types.Log) (*AutomationRegistryLogicA22OwnershipTransferRequested, error) {
	event := new(AutomationRegistryLogicA22OwnershipTransferRequested)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22OwnershipTransferredIterator struct {
	Event *AutomationRegistryLogicA22OwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22OwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22OwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22OwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22OwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicA22OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22OwnershipTransferredIterator{contract: _AutomationRegistryLogicA22.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22OwnershipTransferred)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseOwnershipTransferred(log types.Log) (*AutomationRegistryLogicA22OwnershipTransferred, error) {
	event := new(AutomationRegistryLogicA22OwnershipTransferred)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22PausedIterator struct {
	Event *AutomationRegistryLogicA22Paused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22PausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22Paused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22Paused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22PausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22PausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22Paused struct {
	Account common.Address
	Raw     types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterPaused(opts *bind.FilterOpts) (*AutomationRegistryLogicA22PausedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22PausedIterator{contract: _AutomationRegistryLogicA22.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22Paused) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22Paused)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParsePaused(log types.Log) (*AutomationRegistryLogicA22Paused, error) {
	event := new(AutomationRegistryLogicA22Paused)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22PayeesUpdatedIterator struct {
	Event *AutomationRegistryLogicA22PayeesUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22PayeesUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22PayeesUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22PayeesUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22PayeesUpdatedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22PayeesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22PayeesUpdated struct {
	Transmitters []common.Address
	Payees       []common.Address
	Raw          types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterPayeesUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicA22PayeesUpdatedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22PayeesUpdatedIterator{contract: _AutomationRegistryLogicA22.contract, event: "PayeesUpdated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22PayeesUpdated) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22PayeesUpdated)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParsePayeesUpdated(log types.Log) (*AutomationRegistryLogicA22PayeesUpdated, error) {
	event := new(AutomationRegistryLogicA22PayeesUpdated)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22PayeeshipTransferRequestedIterator struct {
	Event *AutomationRegistryLogicA22PayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22PayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22PayeeshipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22PayeeshipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22PayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22PayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22PayeeshipTransferRequested struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicA22PayeeshipTransferRequestedIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22PayeeshipTransferRequestedIterator{contract: _AutomationRegistryLogicA22.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22PayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22PayeeshipTransferRequested)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParsePayeeshipTransferRequested(log types.Log) (*AutomationRegistryLogicA22PayeeshipTransferRequested, error) {
	event := new(AutomationRegistryLogicA22PayeeshipTransferRequested)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22PayeeshipTransferredIterator struct {
	Event *AutomationRegistryLogicA22PayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22PayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22PayeeshipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22PayeeshipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22PayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22PayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22PayeeshipTransferred struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicA22PayeeshipTransferredIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22PayeeshipTransferredIterator{contract: _AutomationRegistryLogicA22.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22PayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22PayeeshipTransferred)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParsePayeeshipTransferred(log types.Log) (*AutomationRegistryLogicA22PayeeshipTransferred, error) {
	event := new(AutomationRegistryLogicA22PayeeshipTransferred)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22PaymentWithdrawnIterator struct {
	Event *AutomationRegistryLogicA22PaymentWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22PaymentWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22PaymentWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22PaymentWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22PaymentWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22PaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22PaymentWithdrawn struct {
	Transmitter common.Address
	Amount      *big.Int
	To          common.Address
	Payee       common.Address
	Raw         types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*AutomationRegistryLogicA22PaymentWithdrawnIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22PaymentWithdrawnIterator{contract: _AutomationRegistryLogicA22.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22PaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22PaymentWithdrawn)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParsePaymentWithdrawn(log types.Log) (*AutomationRegistryLogicA22PaymentWithdrawn, error) {
	event := new(AutomationRegistryLogicA22PaymentWithdrawn)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22ReorgedUpkeepReportIterator struct {
	Event *AutomationRegistryLogicA22ReorgedUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22ReorgedUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22ReorgedUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22ReorgedUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22ReorgedUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22ReorgedUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22ReorgedUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22ReorgedUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22ReorgedUpkeepReportIterator{contract: _AutomationRegistryLogicA22.contract, event: "ReorgedUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22ReorgedUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22ReorgedUpkeepReport)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseReorgedUpkeepReport(log types.Log) (*AutomationRegistryLogicA22ReorgedUpkeepReport, error) {
	event := new(AutomationRegistryLogicA22ReorgedUpkeepReport)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22StaleUpkeepReportIterator struct {
	Event *AutomationRegistryLogicA22StaleUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22StaleUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22StaleUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22StaleUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22StaleUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22StaleUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22StaleUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22StaleUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22StaleUpkeepReportIterator{contract: _AutomationRegistryLogicA22.contract, event: "StaleUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22StaleUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22StaleUpkeepReport)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseStaleUpkeepReport(log types.Log) (*AutomationRegistryLogicA22StaleUpkeepReport, error) {
	event := new(AutomationRegistryLogicA22StaleUpkeepReport)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22UnpausedIterator struct {
	Event *AutomationRegistryLogicA22Unpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22UnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22Unpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22Unpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22UnpausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22UnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22Unpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterUnpaused(opts *bind.FilterOpts) (*AutomationRegistryLogicA22UnpausedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22UnpausedIterator{contract: _AutomationRegistryLogicA22.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22Unpaused) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22Unpaused)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseUnpaused(log types.Log) (*AutomationRegistryLogicA22Unpaused, error) {
	event := new(AutomationRegistryLogicA22Unpaused)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22UpkeepAdminTransferRequestedIterator struct {
	Event *AutomationRegistryLogicA22UpkeepAdminTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22UpkeepAdminTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22UpkeepAdminTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22UpkeepAdminTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22UpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22UpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22UpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicA22UpkeepAdminTransferRequestedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22UpkeepAdminTransferRequestedIterator{contract: _AutomationRegistryLogicA22.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22UpkeepAdminTransferRequested)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseUpkeepAdminTransferRequested(log types.Log) (*AutomationRegistryLogicA22UpkeepAdminTransferRequested, error) {
	event := new(AutomationRegistryLogicA22UpkeepAdminTransferRequested)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22UpkeepAdminTransferredIterator struct {
	Event *AutomationRegistryLogicA22UpkeepAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22UpkeepAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22UpkeepAdminTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22UpkeepAdminTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22UpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22UpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22UpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicA22UpkeepAdminTransferredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22UpkeepAdminTransferredIterator{contract: _AutomationRegistryLogicA22.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22UpkeepAdminTransferred)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseUpkeepAdminTransferred(log types.Log) (*AutomationRegistryLogicA22UpkeepAdminTransferred, error) {
	event := new(AutomationRegistryLogicA22UpkeepAdminTransferred)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22UpkeepCanceledIterator struct {
	Event *AutomationRegistryLogicA22UpkeepCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22UpkeepCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22UpkeepCanceled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22UpkeepCanceled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22UpkeepCanceledIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22UpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22UpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*AutomationRegistryLogicA22UpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22UpkeepCanceledIterator{contract: _AutomationRegistryLogicA22.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22UpkeepCanceled)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseUpkeepCanceled(log types.Log) (*AutomationRegistryLogicA22UpkeepCanceled, error) {
	event := new(AutomationRegistryLogicA22UpkeepCanceled)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22UpkeepCheckDataSetIterator struct {
	Event *AutomationRegistryLogicA22UpkeepCheckDataSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22UpkeepCheckDataSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22UpkeepCheckDataSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22UpkeepCheckDataSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22UpkeepCheckDataSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22UpkeepCheckDataSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22UpkeepCheckDataSet struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22UpkeepCheckDataSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22UpkeepCheckDataSetIterator{contract: _AutomationRegistryLogicA22.contract, event: "UpkeepCheckDataSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepCheckDataSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22UpkeepCheckDataSet)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseUpkeepCheckDataSet(log types.Log) (*AutomationRegistryLogicA22UpkeepCheckDataSet, error) {
	event := new(AutomationRegistryLogicA22UpkeepCheckDataSet)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22UpkeepGasLimitSetIterator struct {
	Event *AutomationRegistryLogicA22UpkeepGasLimitSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22UpkeepGasLimitSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22UpkeepGasLimitSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22UpkeepGasLimitSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22UpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22UpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22UpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22UpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22UpkeepGasLimitSetIterator{contract: _AutomationRegistryLogicA22.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22UpkeepGasLimitSet)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseUpkeepGasLimitSet(log types.Log) (*AutomationRegistryLogicA22UpkeepGasLimitSet, error) {
	event := new(AutomationRegistryLogicA22UpkeepGasLimitSet)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22UpkeepMigratedIterator struct {
	Event *AutomationRegistryLogicA22UpkeepMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22UpkeepMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22UpkeepMigrated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22UpkeepMigrated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22UpkeepMigratedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22UpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22UpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22UpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22UpkeepMigratedIterator{contract: _AutomationRegistryLogicA22.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22UpkeepMigrated)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseUpkeepMigrated(log types.Log) (*AutomationRegistryLogicA22UpkeepMigrated, error) {
	event := new(AutomationRegistryLogicA22UpkeepMigrated)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22UpkeepOffchainConfigSetIterator struct {
	Event *AutomationRegistryLogicA22UpkeepOffchainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22UpkeepOffchainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22UpkeepOffchainConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22UpkeepOffchainConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22UpkeepOffchainConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22UpkeepOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22UpkeepOffchainConfigSet struct {
	Id             *big.Int
	OffchainConfig []byte
	Raw            types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22UpkeepOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22UpkeepOffchainConfigSetIterator{contract: _AutomationRegistryLogicA22.contract, event: "UpkeepOffchainConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22UpkeepOffchainConfigSet)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseUpkeepOffchainConfigSet(log types.Log) (*AutomationRegistryLogicA22UpkeepOffchainConfigSet, error) {
	event := new(AutomationRegistryLogicA22UpkeepOffchainConfigSet)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22UpkeepPausedIterator struct {
	Event *AutomationRegistryLogicA22UpkeepPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22UpkeepPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22UpkeepPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22UpkeepPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22UpkeepPausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22UpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22UpkeepPaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22UpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22UpkeepPausedIterator{contract: _AutomationRegistryLogicA22.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22UpkeepPaused)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseUpkeepPaused(log types.Log) (*AutomationRegistryLogicA22UpkeepPaused, error) {
	event := new(AutomationRegistryLogicA22UpkeepPaused)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22UpkeepPerformedIterator struct {
	Event *AutomationRegistryLogicA22UpkeepPerformed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22UpkeepPerformedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22UpkeepPerformed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22UpkeepPerformed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22UpkeepPerformedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22UpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22UpkeepPerformed struct {
	Id           *big.Int
	Success      bool
	TotalPayment *big.Int
	GasUsed      *big.Int
	GasOverhead  *big.Int
	Trigger      []byte
	Raw          types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*AutomationRegistryLogicA22UpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22UpkeepPerformedIterator{contract: _AutomationRegistryLogicA22.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22UpkeepPerformed)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseUpkeepPerformed(log types.Log) (*AutomationRegistryLogicA22UpkeepPerformed, error) {
	event := new(AutomationRegistryLogicA22UpkeepPerformed)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22UpkeepPrivilegeConfigSetIterator struct {
	Event *AutomationRegistryLogicA22UpkeepPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22UpkeepPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22UpkeepPrivilegeConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22UpkeepPrivilegeConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22UpkeepPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22UpkeepPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22UpkeepPrivilegeConfigSet struct {
	Id              *big.Int
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22UpkeepPrivilegeConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22UpkeepPrivilegeConfigSetIterator{contract: _AutomationRegistryLogicA22.contract, event: "UpkeepPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22UpkeepPrivilegeConfigSet)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseUpkeepPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicA22UpkeepPrivilegeConfigSet, error) {
	event := new(AutomationRegistryLogicA22UpkeepPrivilegeConfigSet)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22UpkeepReceivedIterator struct {
	Event *AutomationRegistryLogicA22UpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22UpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22UpkeepReceived)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22UpkeepReceived)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22UpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22UpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22UpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22UpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22UpkeepReceivedIterator{contract: _AutomationRegistryLogicA22.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22UpkeepReceived)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseUpkeepReceived(log types.Log) (*AutomationRegistryLogicA22UpkeepReceived, error) {
	event := new(AutomationRegistryLogicA22UpkeepReceived)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22UpkeepRegisteredIterator struct {
	Event *AutomationRegistryLogicA22UpkeepRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22UpkeepRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22UpkeepRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22UpkeepRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22UpkeepRegisteredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22UpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22UpkeepRegistered struct {
	Id         *big.Int
	PerformGas uint32
	Admin      common.Address
	Raw        types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22UpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22UpkeepRegisteredIterator{contract: _AutomationRegistryLogicA22.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22UpkeepRegistered)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseUpkeepRegistered(log types.Log) (*AutomationRegistryLogicA22UpkeepRegistered, error) {
	event := new(AutomationRegistryLogicA22UpkeepRegistered)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22UpkeepTriggerConfigSetIterator struct {
	Event *AutomationRegistryLogicA22UpkeepTriggerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22UpkeepTriggerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22UpkeepTriggerConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22UpkeepTriggerConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22UpkeepTriggerConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22UpkeepTriggerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22UpkeepTriggerConfigSet struct {
	Id            *big.Int
	TriggerConfig []byte
	Raw           types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22UpkeepTriggerConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22UpkeepTriggerConfigSetIterator{contract: _AutomationRegistryLogicA22.contract, event: "UpkeepTriggerConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22UpkeepTriggerConfigSet)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseUpkeepTriggerConfigSet(log types.Log) (*AutomationRegistryLogicA22UpkeepTriggerConfigSet, error) {
	event := new(AutomationRegistryLogicA22UpkeepTriggerConfigSet)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicA22UpkeepUnpausedIterator struct {
	Event *AutomationRegistryLogicA22UpkeepUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicA22UpkeepUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicA22UpkeepUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(AutomationRegistryLogicA22UpkeepUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *AutomationRegistryLogicA22UpkeepUnpausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicA22UpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicA22UpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22UpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA22UpkeepUnpausedIterator{contract: _AutomationRegistryLogicA22.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA22.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicA22UpkeepUnpaused)
				if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22Filterer) ParseUpkeepUnpaused(log types.Log) (*AutomationRegistryLogicA22UpkeepUnpaused, error) {
	event := new(AutomationRegistryLogicA22UpkeepUnpaused)
	if err := _AutomationRegistryLogicA22.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _AutomationRegistryLogicA22.abi.Events["AdminPrivilegeConfigSet"].ID:
		return _AutomationRegistryLogicA22.ParseAdminPrivilegeConfigSet(log)
	case _AutomationRegistryLogicA22.abi.Events["CancelledUpkeepReport"].ID:
		return _AutomationRegistryLogicA22.ParseCancelledUpkeepReport(log)
	case _AutomationRegistryLogicA22.abi.Events["ChainSpecificModuleUpdated"].ID:
		return _AutomationRegistryLogicA22.ParseChainSpecificModuleUpdated(log)
	case _AutomationRegistryLogicA22.abi.Events["DedupKeyAdded"].ID:
		return _AutomationRegistryLogicA22.ParseDedupKeyAdded(log)
	case _AutomationRegistryLogicA22.abi.Events["FundsAdded"].ID:
		return _AutomationRegistryLogicA22.ParseFundsAdded(log)
	case _AutomationRegistryLogicA22.abi.Events["FundsWithdrawn"].ID:
		return _AutomationRegistryLogicA22.ParseFundsWithdrawn(log)
	case _AutomationRegistryLogicA22.abi.Events["InsufficientFundsUpkeepReport"].ID:
		return _AutomationRegistryLogicA22.ParseInsufficientFundsUpkeepReport(log)
	case _AutomationRegistryLogicA22.abi.Events["OwnerFundsWithdrawn"].ID:
		return _AutomationRegistryLogicA22.ParseOwnerFundsWithdrawn(log)
	case _AutomationRegistryLogicA22.abi.Events["OwnershipTransferRequested"].ID:
		return _AutomationRegistryLogicA22.ParseOwnershipTransferRequested(log)
	case _AutomationRegistryLogicA22.abi.Events["OwnershipTransferred"].ID:
		return _AutomationRegistryLogicA22.ParseOwnershipTransferred(log)
	case _AutomationRegistryLogicA22.abi.Events["Paused"].ID:
		return _AutomationRegistryLogicA22.ParsePaused(log)
	case _AutomationRegistryLogicA22.abi.Events["PayeesUpdated"].ID:
		return _AutomationRegistryLogicA22.ParsePayeesUpdated(log)
	case _AutomationRegistryLogicA22.abi.Events["PayeeshipTransferRequested"].ID:
		return _AutomationRegistryLogicA22.ParsePayeeshipTransferRequested(log)
	case _AutomationRegistryLogicA22.abi.Events["PayeeshipTransferred"].ID:
		return _AutomationRegistryLogicA22.ParsePayeeshipTransferred(log)
	case _AutomationRegistryLogicA22.abi.Events["PaymentWithdrawn"].ID:
		return _AutomationRegistryLogicA22.ParsePaymentWithdrawn(log)
	case _AutomationRegistryLogicA22.abi.Events["ReorgedUpkeepReport"].ID:
		return _AutomationRegistryLogicA22.ParseReorgedUpkeepReport(log)
	case _AutomationRegistryLogicA22.abi.Events["StaleUpkeepReport"].ID:
		return _AutomationRegistryLogicA22.ParseStaleUpkeepReport(log)
	case _AutomationRegistryLogicA22.abi.Events["Unpaused"].ID:
		return _AutomationRegistryLogicA22.ParseUnpaused(log)
	case _AutomationRegistryLogicA22.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _AutomationRegistryLogicA22.ParseUpkeepAdminTransferRequested(log)
	case _AutomationRegistryLogicA22.abi.Events["UpkeepAdminTransferred"].ID:
		return _AutomationRegistryLogicA22.ParseUpkeepAdminTransferred(log)
	case _AutomationRegistryLogicA22.abi.Events["UpkeepCanceled"].ID:
		return _AutomationRegistryLogicA22.ParseUpkeepCanceled(log)
	case _AutomationRegistryLogicA22.abi.Events["UpkeepCheckDataSet"].ID:
		return _AutomationRegistryLogicA22.ParseUpkeepCheckDataSet(log)
	case _AutomationRegistryLogicA22.abi.Events["UpkeepGasLimitSet"].ID:
		return _AutomationRegistryLogicA22.ParseUpkeepGasLimitSet(log)
	case _AutomationRegistryLogicA22.abi.Events["UpkeepMigrated"].ID:
		return _AutomationRegistryLogicA22.ParseUpkeepMigrated(log)
	case _AutomationRegistryLogicA22.abi.Events["UpkeepOffchainConfigSet"].ID:
		return _AutomationRegistryLogicA22.ParseUpkeepOffchainConfigSet(log)
	case _AutomationRegistryLogicA22.abi.Events["UpkeepPaused"].ID:
		return _AutomationRegistryLogicA22.ParseUpkeepPaused(log)
	case _AutomationRegistryLogicA22.abi.Events["UpkeepPerformed"].ID:
		return _AutomationRegistryLogicA22.ParseUpkeepPerformed(log)
	case _AutomationRegistryLogicA22.abi.Events["UpkeepPrivilegeConfigSet"].ID:
		return _AutomationRegistryLogicA22.ParseUpkeepPrivilegeConfigSet(log)
	case _AutomationRegistryLogicA22.abi.Events["UpkeepReceived"].ID:
		return _AutomationRegistryLogicA22.ParseUpkeepReceived(log)
	case _AutomationRegistryLogicA22.abi.Events["UpkeepRegistered"].ID:
		return _AutomationRegistryLogicA22.ParseUpkeepRegistered(log)
	case _AutomationRegistryLogicA22.abi.Events["UpkeepTriggerConfigSet"].ID:
		return _AutomationRegistryLogicA22.ParseUpkeepTriggerConfigSet(log)
	case _AutomationRegistryLogicA22.abi.Events["UpkeepUnpaused"].ID:
		return _AutomationRegistryLogicA22.ParseUpkeepUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (AutomationRegistryLogicA22AdminPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x7c44b4eb59ee7873514e7e43e7718c269d872965938b288aa143befca62f99d2")
}

func (AutomationRegistryLogicA22CancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd5636")
}

func (AutomationRegistryLogicA22ChainSpecificModuleUpdated) Topic() common.Hash {
	return common.HexToHash("0xdefc28b11a7980dbe0c49dbbd7055a1584bc8075097d1e8b3b57fb7283df2ad7")
}

func (AutomationRegistryLogicA22DedupKeyAdded) Topic() common.Hash {
	return common.HexToHash("0xa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f2")
}

func (AutomationRegistryLogicA22FundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (AutomationRegistryLogicA22FundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (AutomationRegistryLogicA22InsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e02")
}

func (AutomationRegistryLogicA22OwnerFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1")
}

func (AutomationRegistryLogicA22OwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (AutomationRegistryLogicA22OwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (AutomationRegistryLogicA22Paused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (AutomationRegistryLogicA22PayeesUpdated) Topic() common.Hash {
	return common.HexToHash("0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725")
}

func (AutomationRegistryLogicA22PayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (AutomationRegistryLogicA22PayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (AutomationRegistryLogicA22PaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (AutomationRegistryLogicA22ReorgedUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301")
}

func (AutomationRegistryLogicA22StaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8")
}

func (AutomationRegistryLogicA22Unpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (AutomationRegistryLogicA22UpkeepAdminTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35")
}

func (AutomationRegistryLogicA22UpkeepAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c")
}

func (AutomationRegistryLogicA22UpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (AutomationRegistryLogicA22UpkeepCheckDataSet) Topic() common.Hash {
	return common.HexToHash("0xcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d")
}

func (AutomationRegistryLogicA22UpkeepGasLimitSet) Topic() common.Hash {
	return common.HexToHash("0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c")
}

func (AutomationRegistryLogicA22UpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (AutomationRegistryLogicA22UpkeepOffchainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850")
}

func (AutomationRegistryLogicA22UpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (AutomationRegistryLogicA22UpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b")
}

func (AutomationRegistryLogicA22UpkeepPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae7769")
}

func (AutomationRegistryLogicA22UpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (AutomationRegistryLogicA22UpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (AutomationRegistryLogicA22UpkeepTriggerConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664")
}

func (AutomationRegistryLogicA22UpkeepUnpaused) Topic() common.Hash {
	return common.HexToHash("0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456")
}

func (_AutomationRegistryLogicA22 *AutomationRegistryLogicA22) Address() common.Address {
	return _AutomationRegistryLogicA22.address
}

type AutomationRegistryLogicA22Interface interface {
	FallbackTo(opts *bind.CallOpts) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error)

	CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	CheckCallback(opts *bind.TransactOpts, id *big.Int, values [][]byte, extraData []byte) (*types.Transaction, error)

	CheckUpkeep(opts *bind.TransactOpts, id *big.Int, triggerData []byte) (*types.Transaction, error)

	CheckUpkeep0(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	ExecuteCallback(opts *bind.TransactOpts, id *big.Int, payload []byte) (*types.Transaction, error)

	MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error)

	ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error)

	RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error)

	RegisterUpkeep0(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error)

	SetUpkeepTriggerConfig(opts *bind.TransactOpts, id *big.Int, triggerConfig []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error)

	FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*AutomationRegistryLogicA22AdminPrivilegeConfigSetIterator, error)

	WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22AdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error)

	ParseAdminPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicA22AdminPrivilegeConfigSet, error)

	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22CancelledUpkeepReportIterator, error)

	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22CancelledUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseCancelledUpkeepReport(log types.Log) (*AutomationRegistryLogicA22CancelledUpkeepReport, error)

	FilterChainSpecificModuleUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicA22ChainSpecificModuleUpdatedIterator, error)

	WatchChainSpecificModuleUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22ChainSpecificModuleUpdated) (event.Subscription, error)

	ParseChainSpecificModuleUpdated(log types.Log) (*AutomationRegistryLogicA22ChainSpecificModuleUpdated, error)

	FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*AutomationRegistryLogicA22DedupKeyAddedIterator, error)

	WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22DedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error)

	ParseDedupKeyAdded(log types.Log) (*AutomationRegistryLogicA22DedupKeyAdded, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*AutomationRegistryLogicA22FundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22FundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*AutomationRegistryLogicA22FundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22FundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22FundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*AutomationRegistryLogicA22FundsWithdrawn, error)

	FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22InsufficientFundsUpkeepReportIterator, error)

	WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22InsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseInsufficientFundsUpkeepReport(log types.Log) (*AutomationRegistryLogicA22InsufficientFundsUpkeepReport, error)

	FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*AutomationRegistryLogicA22OwnerFundsWithdrawnIterator, error)

	WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22OwnerFundsWithdrawn) (event.Subscription, error)

	ParseOwnerFundsWithdrawn(log types.Log) (*AutomationRegistryLogicA22OwnerFundsWithdrawn, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicA22OwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*AutomationRegistryLogicA22OwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicA22OwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*AutomationRegistryLogicA22OwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*AutomationRegistryLogicA22PausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22Paused) (event.Subscription, error)

	ParsePaused(log types.Log) (*AutomationRegistryLogicA22Paused, error)

	FilterPayeesUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicA22PayeesUpdatedIterator, error)

	WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22PayeesUpdated) (event.Subscription, error)

	ParsePayeesUpdated(log types.Log) (*AutomationRegistryLogicA22PayeesUpdated, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicA22PayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22PayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*AutomationRegistryLogicA22PayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicA22PayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22PayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*AutomationRegistryLogicA22PayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*AutomationRegistryLogicA22PaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22PaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*AutomationRegistryLogicA22PaymentWithdrawn, error)

	FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22ReorgedUpkeepReportIterator, error)

	WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22ReorgedUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseReorgedUpkeepReport(log types.Log) (*AutomationRegistryLogicA22ReorgedUpkeepReport, error)

	FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22StaleUpkeepReportIterator, error)

	WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22StaleUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseStaleUpkeepReport(log types.Log) (*AutomationRegistryLogicA22StaleUpkeepReport, error)

	FilterUnpaused(opts *bind.FilterOpts) (*AutomationRegistryLogicA22UnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22Unpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*AutomationRegistryLogicA22Unpaused, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicA22UpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*AutomationRegistryLogicA22UpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicA22UpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*AutomationRegistryLogicA22UpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*AutomationRegistryLogicA22UpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*AutomationRegistryLogicA22UpkeepCanceled, error)

	FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22UpkeepCheckDataSetIterator, error)

	WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepCheckDataSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCheckDataSet(log types.Log) (*AutomationRegistryLogicA22UpkeepCheckDataSet, error)

	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22UpkeepGasLimitSetIterator, error)

	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepGasLimitSet(log types.Log) (*AutomationRegistryLogicA22UpkeepGasLimitSet, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22UpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*AutomationRegistryLogicA22UpkeepMigrated, error)

	FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22UpkeepOffchainConfigSetIterator, error)

	WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepOffchainConfigSet(log types.Log) (*AutomationRegistryLogicA22UpkeepOffchainConfigSet, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22UpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*AutomationRegistryLogicA22UpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*AutomationRegistryLogicA22UpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*AutomationRegistryLogicA22UpkeepPerformed, error)

	FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22UpkeepPrivilegeConfigSetIterator, error)

	WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicA22UpkeepPrivilegeConfigSet, error)

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22UpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*AutomationRegistryLogicA22UpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22UpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*AutomationRegistryLogicA22UpkeepRegistered, error)

	FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22UpkeepTriggerConfigSetIterator, error)

	WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepTriggerConfigSet(log types.Log) (*AutomationRegistryLogicA22UpkeepTriggerConfigSet, error)

	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicA22UpkeepUnpausedIterator, error)

	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicA22UpkeepUnpaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepUnpaused(log types.Log) (*AutomationRegistryLogicA22UpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
