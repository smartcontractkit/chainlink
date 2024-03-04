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

var AutomationRegistryLogicAMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractAutomationRegistryLogicB2_2\",\"name\":\"logicB\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTrigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTriggerType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"AdminPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newModule\",\"type\":\"address\"}],\"name\":\"ChainSpecificModuleUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"DedupKeyAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepTriggerConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumAutomationRegistryBase2_2.UpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"triggerData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumAutomationRegistryBase2_2.UpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumAutomationRegistryBase2_2.UpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"executeCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumAutomationRegistryBase2_2.UpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fallbackTo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"migrateUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedUpkeeps\",\"type\":\"bytes\"}],\"name\":\"receiveUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"enumAutomationRegistryBase2_2.Trigger\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"setUpkeepTriggerConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101406040523480156200001257600080fd5b5060405162005ec638038062005ec68339810160408190526200003591620003b1565b80816001600160a01b031663ca30e6036040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200009b9190620003b1565b826001600160a01b031663b10b673c6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000da573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001009190620003b1565b836001600160a01b0316636709d0e56040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200013f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001659190620003b1565b846001600160a01b0316635425d8ac6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca9190620003b1565b856001600160a01b031663a08714c06040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000209573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200022f9190620003b1565b3380600081620002865760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620002b957620002b981620002ed565b5050506001600160a01b0394851660805292841660a05290831660c052821660e052811661010052166101205250620003d8565b336001600160a01b03821603620003475760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200027d565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b0381168114620003ae57600080fd5b50565b600060208284031215620003c457600080fd5b8151620003d18162000398565b9392505050565b60805160a05160c05160e0516101005161012051615a896200043d6000396000818161010e01526101a9015260006131050152600081816103e10152611ffe015260006132ee015260006133d2015260008181611e40015261240c0152615a896000f3fe60806040523480156200001157600080fd5b50600436106200010c5760003560e01c806385c1b0ba11620000a5578063c8048022116200006f578063c804802214620002b7578063ce7dc5b414620002ce578063f2fde38b14620002e5578063f7d334ba14620002fc576200010c565b806385c1b0ba14620002535780638da5cb5b146200026a5780638e86139b1462000289578063948108f714620002a0576200010c565b80634ee88d3511620000e75780634ee88d3514620001ef5780636ded9eae146200020657806371791aa0146200021d57806379ba50971462000249576200010c565b806328f32f38146200015457806329c5efad146200017e578063349e8cca14620001a7575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e8080156200014d573d6000f35b3d6000fd5b005b6200016b6200016536600462004024565b62000313565b6040519081526020015b60405180910390f35b620001956200018f3660046200410a565b6200068d565b60405162000175949392919062004232565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200162000175565b62000152620002003660046200426f565b62000931565b6200016b62000217366004620042bf565b62000999565b620002346200022e3660046200410a565b620009ff565b60405162000175979695949392919062004372565b620001526200114d565b6200015262000264366004620043c4565b62001250565b60005473ffffffffffffffffffffffffffffffffffffffff16620001c9565b620001526200029a36600462004451565b62001ec1565b62000152620002b1366004620044b4565b62002249565b62000152620002c8366004620044e3565b620024dc565b62000195620002df366004620045b9565b62002906565b62000152620002f636600462004630565b620029cc565b620002346200030d366004620044e3565b620029e4565b6000805473ffffffffffffffffffffffffffffffffffffffff1633148015906200034757506200034560093362002a22565b155b156200037f576040517fd48b678b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff89163b620003ce576040517f09ee12d500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b620003d98662002a56565b9050600089307f00000000000000000000000000000000000000000000000000000000000000006040516200040e9062003db0565b73ffffffffffffffffffffffffffffffffffffffff938416815291831660208301529091166040820152606001604051809103906000f08015801562000458573d6000803e3d6000fd5b5090506200051f826040518060e001604052806000151581526020018c63ffffffff16815260200163ffffffff801681526020018473ffffffffffffffffffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff168152602001600063ffffffff168152508a89898080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508b92508a915062002d029050565b6015805474010000000000000000000000000000000000000000900463ffffffff169060146200054f836200467f565b91906101000a81548163ffffffff021916908363ffffffff16021790555050817fbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d0128a8a604051620005c892919063ffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a2817fcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d878760405162000604929190620046ee565b60405180910390a2817f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664856040516200063e919062004704565b60405180910390a2817f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf48508460405162000678919062004704565b60405180910390a25098975050505050505050565b600060606000806200069e620030ed565b600086815260046020908152604091829020825160e081018452815460ff811615158252610100810463ffffffff90811694830194909452650100000000008104841694820194909452690100000000000000000090930473ffffffffffffffffffffffffffffffffffffffff166060840152600101546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a0840152780100000000000000000000000000000000000000000000000090041660c08201525a9150600080826060015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620007b8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620007de919062004726565b73ffffffffffffffffffffffffffffffffffffffff166014600101600c9054906101000a900463ffffffff1663ffffffff168960405162000820919062004746565b60006040518083038160008787f1925050503d806000811462000860576040519150601f19603f3d011682016040523d82523d6000602084013e62000865565b606091505b50915091505a62000877908562004764565b935081620008a257600060405180602001604052806000815250600796509650965050505062000928565b80806020019051810190620008b89190620047d5565b909750955086620008e657600060405180602001604052806000815250600496509650965050505062000928565b601654865164010000000090910463ffffffff1610156200092457600060405180602001604052806000815250600596509650965050505062000928565b5050505b92959194509250565b6200093c836200315f565b6000838152601b6020526040902062000957828483620048ca565b50827f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d566483836040516200098c929190620046ee565b60405180910390a2505050565b6000620009f388888860008989604051806020016040528060008152508a8a8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506200031392505050565b98975050505050505050565b60006060600080600080600062000a15620030ed565b600062000a228a62003215565b905060006012604051806101600160405290816000820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160008201600c9054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160109054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160149054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160189054906101000a900462ffffff1662ffffff1662ffffff16815260200160008201601b9054906101000a900461ffff1661ffff1661ffff16815260200160008201601d9054906101000a900460ff1660ff1660ff16815260200160008201601e9054906101000a900460ff1615151515815260200160008201601f9054906101000a900460ff161515151581526020016001820160009054906101000a900460ff161515151581526020016001820160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152505090506000600460008d81526020019081526020016000206040518060e00160405290816000820160009054906101000a900460ff161515151581526020016000820160019054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160059054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160099054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160018201600c9054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff1681526020016001820160189054906101000a900463ffffffff1663ffffffff1663ffffffff168152505090508160e001511562000db7576000604051806020016040528060008152506009600084602001516000808263ffffffff169250995099509950995099509950995050505062001141565b604081015163ffffffff9081161462000e08576000604051806020016040528060008152506001600084602001516000808263ffffffff169250995099509950995099509950995050505062001141565b80511562000e4e576000604051806020016040528060008152506002600084602001516000808263ffffffff169250995099509950995099509950995050505062001141565b62000e5982620032cb565b8095508196505050600062000e76838584602001518989620034bd565b9050806bffffffffffffffffffffffff168260a001516bffffffffffffffffffffffff16101562000ee0576000604051806020016040528060008152506006600085602001516000808263ffffffff1692509a509a509a509a509a509a509a505050505062001141565b600062000eef8e868f62003772565b90505a9850600080846060015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000f47573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000f6d919062004726565b73ffffffffffffffffffffffffffffffffffffffff166014600101600c9054906101000a900463ffffffff1663ffffffff168460405162000faf919062004746565b60006040518083038160008787f1925050503d806000811462000fef576040519150601f19603f3d011682016040523d82523d6000602084013e62000ff4565b606091505b50915091505a62001006908c62004764565b9a5081620010865760165481516801000000000000000090910463ffffffff1610156200106357505060408051602080820190925260008082529490910151939c509a50600899505063ffffffff90911695506200114192505050565b602090940151939b5060039a505063ffffffff9092169650620011419350505050565b808060200190518101906200109c9190620047d5565b909e509c508d620010dd57505060408051602080820190925260008082529490910151939c509a50600499505063ffffffff90911695506200114192505050565b6016548d5164010000000090910463ffffffff1610156200112e57505060408051602080820190925260008082529490910151939c509a50600599505063ffffffff90911695506200114192505050565b505050506020015163ffffffff16945050505b92959891949750929550565b60015473ffffffffffffffffffffffffffffffffffffffff163314620011d4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600173ffffffffffffffffffffffffffffffffffffffff82166000908152601a602052604090205460ff1660038111156200128f576200128f620041c7565b14158015620012db5750600373ffffffffffffffffffffffffffffffffffffffff82166000908152601a602052604090205460ff166003811115620012d857620012d8620041c7565b14155b1562001313576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6014546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1662001373576040517fd12d7d8d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000829003620013af576040517f2c2fc94100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160e081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c081018290526000808567ffffffffffffffff81111562001406576200140662003eab565b60405190808252806020026020018201604052801562001430578160200160208202803683370190505b50905060008667ffffffffffffffff81111562001451576200145162003eab565b604051908082528060200260200182016040528015620014d857816020015b6040805160e08101825260008082526020808301829052928201819052606082018190526080820181905260a0820181905260c082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181620014705790505b50905060008767ffffffffffffffff811115620014f957620014f962003eab565b6040519080825280602002602001820160405280156200152e57816020015b6060815260200190600190039081620015185790505b50905060008867ffffffffffffffff8111156200154f576200154f62003eab565b6040519080825280602002602001820160405280156200158457816020015b60608152602001906001900390816200156e5790505b50905060008967ffffffffffffffff811115620015a557620015a562003eab565b604051908082528060200260200182016040528015620015da57816020015b6060815260200190600190039081620015c45790505b50905060005b8a81101562001bbe578b8b82818110620015fe57620015fe620049f2565b60209081029290920135600081815260048452604090819020815160e081018352815460ff811615158252610100810463ffffffff90811697830197909752650100000000008104871693820193909352690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff166060830152600101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490931660c08401529a50909850620016dd9050896200315f565b60608801516040517f1a5da6c800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8c8116600483015290911690631a5da6c890602401600060405180830381600087803b1580156200174d57600080fd5b505af115801562001762573d6000803e3d6000fd5b50505050878582815181106200177c576200177c620049f2565b6020026020010181905250600560008a815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16868281518110620017d057620017d0620049f2565b73ffffffffffffffffffffffffffffffffffffffff90921660209283029190910182015260008a815260079091526040902080546200180f9062004822565b80601f01602080910402602001604051908101604052809291908181526020018280546200183d9062004822565b80156200188e5780601f1062001862576101008083540402835291602001916200188e565b820191906000526020600020905b8154815290600101906020018083116200187057829003601f168201915b5050505050848281518110620018a857620018a8620049f2565b6020026020010181905250601b60008a81526020019081526020016000208054620018d39062004822565b80601f0160208091040260200160405190810160405280929190818152602001828054620019019062004822565b8015620019525780601f10620019265761010080835404028352916020019162001952565b820191906000526020600020905b8154815290600101906020018083116200193457829003601f168201915b50505050508382815181106200196c576200196c620049f2565b6020026020010181905250601c60008a81526020019081526020016000208054620019979062004822565b80601f0160208091040260200160405190810160405280929190818152602001828054620019c59062004822565b801562001a165780601f10620019ea5761010080835404028352916020019162001a16565b820191906000526020600020905b815481529060010190602001808311620019f857829003601f168201915b505050505082828151811062001a305762001a30620049f2565b60200260200101819052508760a001516bffffffffffffffffffffffff168762001a5b919062004a21565b60008a815260046020908152604080832080547fffffff000000000000000000000000000000000000000000000000000000000016815560010180547fffffffff000000000000000000000000000000000000000000000000000000001690556007909152812091985062001ad1919062003dbe565b6000898152601b6020526040812062001aea9162003dbe565b6000898152601c6020526040812062001b039162003dbe565b600089815260066020526040902080547fffffffffffffffffffffffff000000000000000000000000000000000000000016905562001b4460028a62003962565b5060a0880151604080516bffffffffffffffffffffffff909216825273ffffffffffffffffffffffffffffffffffffffff8c1660208301528a917fb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff910160405180910390a28062001bb58162004a37565b915050620015e0565b508560195462001bcf919062004764565b60195560008b8b868167ffffffffffffffff81111562001bf35762001bf362003eab565b60405190808252806020026020018201604052801562001c1d578160200160208202803683370190505b508988888860405160200162001c3b98979695949392919062004bfe565b60405160208183030381529060405290508973ffffffffffffffffffffffffffffffffffffffff16638e86139b6014600001600c9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663c71249ab60038e73ffffffffffffffffffffffffffffffffffffffff1663aab9edd66040518163ffffffff1660e01b8152600401602060405180830381865afa15801562001cf7573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001d1d919062004cdd565b866040518463ffffffff1660e01b815260040162001d3e9392919062004d02565b600060405180830381865afa15801562001d5c573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405262001da4919081019062004d29565b6040518263ffffffff1660e01b815260040162001dc2919062004704565b600060405180830381600087803b15801562001ddd57600080fd5b505af115801562001df2573d6000803e3d6000fd5b50506040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8d81166004830152602482018b90527f000000000000000000000000000000000000000000000000000000000000000016925063a9059cbb91506044016020604051808303816000875af115801562001e8c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001eb2919062004d62565b50505050505050505050505050565b6002336000908152601a602052604090205460ff16600381111562001eea5762001eea620041c7565b1415801562001f2057506003336000908152601a602052604090205460ff16600381111562001f1d5762001f1d620041c7565b14155b1562001f58576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080808080808062001f6e888a018a62004f4d565b965096509650965096509650965060005b87518110156200223d57600073ffffffffffffffffffffffffffffffffffffffff1687828151811062001fb65762001fb6620049f2565b60200260200101516060015173ffffffffffffffffffffffffffffffffffffffff1603620020ca5785818151811062001ff35762001ff3620049f2565b6020026020010151307f00000000000000000000000000000000000000000000000000000000000000006040516200202b9062003db0565b73ffffffffffffffffffffffffffffffffffffffff938416815291831660208301529091166040820152606001604051809103906000f08015801562002075573d6000803e3d6000fd5b508782815181106200208b576200208b620049f2565b60200260200101516060019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250505b62002182888281518110620020e357620020e3620049f2565b6020026020010151888381518110620021005762002100620049f2565b60200260200101518784815181106200211d576200211d620049f2565b60200260200101518785815181106200213a576200213a620049f2565b6020026020010151878681518110620021575762002157620049f2565b6020026020010151878781518110620021745762002174620049f2565b602002602001015162002d02565b878181518110620021975762002197620049f2565b60200260200101517f74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71888381518110620021d557620021d5620049f2565b602002602001015160a0015133604051620022209291906bffffffffffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a280620022348162004a37565b91505062001f7f565b50505050505050505050565b600082815260046020908152604091829020825160e081018452815460ff81161515825263ffffffff6101008204811694830194909452650100000000008104841694820185905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004821660c0820152911462002347576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818160a001516200235991906200507e565b600084815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff93841602179055601954620023c19184169062004a21565b6019556040517f23b872dd0000000000000000000000000000000000000000000000000000000081523360048201523060248201526bffffffffffffffffffffffff831660448201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906323b872dd906064016020604051808303816000875af11580156200246b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062002491919062004d62565b506040516bffffffffffffffffffffffff83168152339084907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a3505050565b6000818152600460209081526040808320815160e081018352815460ff81161515825263ffffffff6101008204811695830195909552650100000000008104851693820184905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004831660c08201529291141590620025c560005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161490506000601260010160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa15801562002668573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200268e9190620050a6565b90508215620026c9576040517ffbc0357800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81158015620026fc575060008581526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314155b1562002734576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b816200274a576200274760328262004a21565b90505b6000858152600460205260409020805463ffffffff80841665010000000000027fffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffffff90921691909117909155620027a69060029087906200396216565b5060145460808501516bffffffffffffffffffffffff91821691600091168211156200280f576080860151620027dd9083620050c0565b90508560a001516bffffffffffffffffffffffff16816bffffffffffffffffffffffff1611156200280f575060a08501515b808660a00151620028219190620050c0565b600088815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff9384160217905560155462002889918391166200507e565b601580547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff9290921691909117905560405167ffffffffffffffff84169088907f91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f79118190600090a350505050505050565b600060606000806000634b56a42e60e01b8888886040516024016200292e93929190620050e8565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091529050620029b989826200068d565b929c919b50995090975095505050505050565b620029d662003970565b620029e181620039f3565b50565b60006060600080600080600062002a0b8860405180602001604052806000815250620009ff565b959e949d50929b5090995097509550909350915050565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415155b90505b92915050565b6000806000601260010160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905060008173ffffffffffffffffffffffffffffffffffffffff166385df51fd60018473ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa15801562002aef573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062002b159190620050a6565b62002b21919062004764565b6040518263ffffffff1660e01b815260040162002b4091815260200190565b602060405180830381865afa15801562002b5e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062002b849190620050a6565b601554604080516020810193909352309083015274010000000000000000000000000000000000000000900463ffffffff166060820152608001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201209083015201604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052905060045b600f81101562002c90578382828151811062002c4c5762002c4c620049f2565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508062002c878162004a37565b91505062002c2c565b5084600181111562002ca65762002ca6620041c7565b60f81b81600f8151811062002cbf5762002cbf620049f2565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535062002cf9816200511c565b95945050505050565b6012547e01000000000000000000000000000000000000000000000000000000000000900460ff161562002d62576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601654835163ffffffff909116101562002da8576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108fc856020015163ffffffff16108062002de65750601554602086015163ffffffff70010000000000000000000000000000000090920482169116115b1562002e1e576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000868152600460205260409020546901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff161562002e88576040517f6e3b930b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000868152600460209081526040808320885181548a8501518b85015160608d01517fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000009093169315157fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000ff169390931761010063ffffffff92831602177fffffff000000000000000000000000000000000000000000000000ffffffffff1665010000000000938216939093027fffffff0000000000000000000000000000000000000000ffffffffffffffffff1692909217690100000000000000000073ffffffffffffffffffffffffffffffffffffffff9283160217835560808b01516001909301805460a08d015160c08e01516bffffffffffffffffffffffff9687167fffffffffffffffff000000000000000000000000000000000000000000000000909316929092176c010000000000000000000000009690911695909502949094177fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff1678010000000000000000000000000000000000000000000000009490931693909302919091179091556005835281842080547fffffffffffffffffffffffff000000000000000000000000000000000000000016918916919091179055600790915290206200307b84826200515f565b508460a001516bffffffffffffffffffffffff166019546200309e919062004a21565b6019556000868152601b60205260409020620030bb83826200515f565b506000868152601c60205260409020620030d682826200515f565b50620030e460028762003aea565b50505050505050565b3273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146200315d576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60008181526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314620031bd576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526004602052604090205465010000000000900463ffffffff90811614620029e1576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000818160045b600f811015620032aa577fff0000000000000000000000000000000000000000000000000000000000000082168382602081106200325e576200325e620049f2565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916146200329557506000949350505050565b80620032a18162004a37565b9150506200321c565b5081600f1a6001811115620032c357620032c3620041c7565b949350505050565b6000806000836080015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa15801562003358573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200337e9190620052a1565b50945090925050506000811315806200339657508142105b80620033bb5750828015620033bb5750620033b2824262004764565b8463ffffffff16105b15620033cc576017549550620033d0565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa1580156200343c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620034629190620052a1565b50945090925050506000811315806200347a57508142105b806200349f57508280156200349f575062003496824262004764565b8463ffffffff16105b15620034b0576018549450620034b4565b8094505b50505050915091565b60008080866001811115620034d657620034d6620041c7565b03620034e6575061ea6062003540565b6001866001811115620034fd57620034fd620041c7565b036200350e575062014c0862003540565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008760c001516001620035559190620052f6565b620035659060ff16604062005312565b60165462003585906103a490640100000000900463ffffffff1662004a21565b62003591919062004a21565b601354604080517fde9ee35e00000000000000000000000000000000000000000000000000000000815281519394506000938493610100900473ffffffffffffffffffffffffffffffffffffffff169263de9ee35e92600480820193918290030181865afa15801562003608573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200362e91906200532c565b909250905081836200364283601862004a21565b6200364e919062005312565b60c08c015162003660906001620052f6565b620036719060ff166115e062005312565b6200367d919062004a21565b62003689919062004a21565b62003695908562004a21565b935060008a610140015173ffffffffffffffffffffffffffffffffffffffff166312544140856040518263ffffffff1660e01b8152600401620036da91815260200190565b602060405180830381865afa158015620036f8573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200371e9190620050a6565b8b60a0015161ffff1662003733919062005312565b9050600080620037508d8c63ffffffff1689868e8e600062003af8565b90925090506200376181836200507e565b9d9c50505050505050505050505050565b606060008360018111156200378b576200378b620041c7565b0362003858576000848152600760205260409081902090517f6e04ff0d0000000000000000000000000000000000000000000000000000000091620037d391602401620053f4565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915290506200395b565b60018360018111156200386f576200386f620041c7565b036200350e576000828060200190518101906200388d91906200546b565b6000868152600760205260409081902090519192507f40691db40000000000000000000000000000000000000000000000000000000091620038d49184916024016200557f565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915291506200395b9050565b9392505050565b600062002a4d838362003c53565b60005473ffffffffffffffffffffffffffffffffffffffff1633146200315d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401620011cb565b3373ffffffffffffffffffffffffffffffffffffffff82160362003a74576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620011cb565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600062002a4d838362003d5e565b60008060008960a0015161ffff168662003b13919062005312565b905083801562003b225750803a105b1562003b2b57503a5b6000858862003b3b8b8d62004a21565b62003b47908562005312565b62003b53919062004a21565b62003b6790670de0b6b3a764000062005312565b62003b73919062005647565b905060008b6040015163ffffffff1664e8d4a5100062003b94919062005312565b60208d0151889063ffffffff168b62003bae8f8862005312565b62003bba919062004a21565b62003bca90633b9aca0062005312565b62003bd6919062005312565b62003be2919062005647565b62003bee919062004a21565b90506b033b2e3c9fd0803ce800000062003c09828462004a21565b111562003c42576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909b909a5098505050505050505050565b6000818152600183016020526040812054801562003d4c57600062003c7a60018362004764565b855490915060009062003c909060019062004764565b905081811462003cfc57600086600001828154811062003cb45762003cb4620049f2565b906000526020600020015490508087600001848154811062003cda5762003cda620049f2565b6000918252602080832090910192909255918252600188019052604090208390555b855486908062003d105762003d1062005683565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505062002a50565b600091505062002a50565b5092915050565b600081815260018301602052604081205462003da75750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915562002a50565b50600062002a50565b6103ca80620056b383390190565b50805462003dcc9062004822565b6000825580601f1062003ddd575050565b601f016020900490600052602060002090810190620029e191905b8082111562003e0e576000815560010162003df8565b5090565b73ffffffffffffffffffffffffffffffffffffffff81168114620029e157600080fd5b803563ffffffff8116811462003e4a57600080fd5b919050565b80356002811062003e4a57600080fd5b60008083601f84011262003e7257600080fd5b50813567ffffffffffffffff81111562003e8b57600080fd5b60208301915083602082850101111562003ea457600080fd5b9250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160e0810167ffffffffffffffff8111828210171562003f005762003f0062003eab565b60405290565b604051610100810167ffffffffffffffff8111828210171562003f005762003f0062003eab565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171562003f775762003f7762003eab565b604052919050565b600067ffffffffffffffff82111562003f9c5762003f9c62003eab565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f83011262003fda57600080fd5b813562003ff162003feb8262003f7f565b62003f2d565b8181528460208386010111156200400757600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060008060e0898b0312156200404157600080fd5b88356200404e8162003e12565b97506200405e60208a0162003e35565b96506040890135620040708162003e12565b95506200408060608a0162003e4f565b9450608089013567ffffffffffffffff808211156200409e57600080fd5b620040ac8c838d0162003e5f565b909650945060a08b0135915080821115620040c657600080fd5b620040d48c838d0162003fc8565b935060c08b0135915080821115620040eb57600080fd5b50620040fa8b828c0162003fc8565b9150509295985092959890939650565b600080604083850312156200411e57600080fd5b82359150602083013567ffffffffffffffff8111156200413d57600080fd5b6200414b8582860162003fc8565b9150509250929050565b60005b838110156200417257818101518382015260200162004158565b50506000910152565b600081518084526200419581602086016020860162004155565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600a81106200422e577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b9052565b84151581526080602082015260006200424f60808301866200417b565b9050620042606040830185620041f6565b82606083015295945050505050565b6000806000604084860312156200428557600080fd5b83359250602084013567ffffffffffffffff811115620042a457600080fd5b620042b28682870162003e5f565b9497909650939450505050565b600080600080600080600060a0888a031215620042db57600080fd5b8735620042e88162003e12565b9650620042f86020890162003e35565b955060408801356200430a8162003e12565b9450606088013567ffffffffffffffff808211156200432857600080fd5b620043368b838c0162003e5f565b909650945060808a01359150808211156200435057600080fd5b506200435f8a828b0162003e5f565b989b979a50959850939692959293505050565b871515815260e0602082015260006200438f60e08301896200417b565b9050620043a06040830188620041f6565b8560608301528460808301528360a08301528260c083015298975050505050505050565b600080600060408486031215620043da57600080fd5b833567ffffffffffffffff80821115620043f357600080fd5b818601915086601f8301126200440857600080fd5b8135818111156200441857600080fd5b8760208260051b85010111156200442e57600080fd5b60209283019550935050840135620044468162003e12565b809150509250925092565b600080602083850312156200446557600080fd5b823567ffffffffffffffff8111156200447d57600080fd5b6200448b8582860162003e5f565b90969095509350505050565b80356bffffffffffffffffffffffff8116811462003e4a57600080fd5b60008060408385031215620044c857600080fd5b82359150620044da6020840162004497565b90509250929050565b600060208284031215620044f657600080fd5b5035919050565b600067ffffffffffffffff8211156200451a576200451a62003eab565b5060051b60200190565b600082601f8301126200453657600080fd5b813560206200454962003feb83620044fd565b82815260059290921b840181019181810190868411156200456957600080fd5b8286015b84811015620045ae57803567ffffffffffffffff8111156200458f5760008081fd5b6200459f8986838b010162003fc8565b8452509183019183016200456d565b509695505050505050565b60008060008060608587031215620045d057600080fd5b84359350602085013567ffffffffffffffff80821115620045f057600080fd5b620045fe8883890162004524565b945060408701359150808211156200461557600080fd5b50620046248782880162003e5f565b95989497509550505050565b6000602082840312156200464357600080fd5b81356200395b8162003e12565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600063ffffffff8083168181036200469b576200469b62004650565b6001019392505050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b602081526000620032c3602083018486620046a5565b60208152600062002a4d60208301846200417b565b805162003e4a8162003e12565b6000602082840312156200473957600080fd5b81516200395b8162003e12565b600082516200475a81846020870162004155565b9190910192915050565b8181038181111562002a505762002a5062004650565b8015158114620029e157600080fd5b600082601f8301126200479b57600080fd5b8151620047ac62003feb8262003f7f565b818152846020838601011115620047c257600080fd5b620032c382602083016020870162004155565b60008060408385031215620047e957600080fd5b8251620047f6816200477a565b602084015190925067ffffffffffffffff8111156200481457600080fd5b6200414b8582860162004789565b600181811c908216806200483757607f821691505b60208210810362004871577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f821115620048c557600081815260208120601f850160051c81016020861015620048a05750805b601f850160051c820191505b81811015620048c157828155600101620048ac565b5050505b505050565b67ffffffffffffffff831115620048e557620048e562003eab565b620048fd83620048f6835462004822565b8362004877565b6000601f8411600181146200495257600085156200491b5750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355620049eb565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b82811015620049a3578685013582556020948501946001909201910162004981565b5086821015620049df577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b8082018082111562002a505762002a5062004650565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820362004a6b5762004a6b62004650565b5060010190565b600081518084526020808501945080840160005b8381101562004b315781518051151588528381015163ffffffff908116858a01526040808301519091169089015260608082015173ffffffffffffffffffffffffffffffffffffffff16908901526080808201516bffffffffffffffffffffffff169089015260a08082015162004b0c828b01826bffffffffffffffffffffffff169052565b505060c09081015163ffffffff169088015260e0909601959082019060010162004a86565b509495945050505050565b600081518084526020808501945080840160005b8381101562004b3157815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010162004b50565b600082825180855260208086019550808260051b84010181860160005b8481101562004bf1577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe086840301895262004bde8383516200417b565b9884019892509083019060010162004ba1565b5090979650505050505050565b60e081528760e082015260006101007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8a111562004c3b57600080fd5b8960051b808c8386013783018381038201602085015262004c5f8282018b62004a72565b915050828103604084015262004c76818962004b3c565b9050828103606084015262004c8c818862004b3c565b9050828103608084015262004ca2818762004b84565b905082810360a084015262004cb8818662004b84565b905082810360c084015262004cce818562004b84565b9b9a5050505050505050505050565b60006020828403121562004cf057600080fd5b815160ff811681146200395b57600080fd5b60ff8416815260ff8316602082015260606040820152600062002cf960608301846200417b565b60006020828403121562004d3c57600080fd5b815167ffffffffffffffff81111562004d5457600080fd5b620032c38482850162004789565b60006020828403121562004d7557600080fd5b81516200395b816200477a565b600082601f83011262004d9457600080fd5b8135602062004da762003feb83620044fd565b82815260059290921b8401810191818101908684111562004dc757600080fd5b8286015b84811015620045ae578035835291830191830162004dcb565b600082601f83011262004df657600080fd5b8135602062004e0962003feb83620044fd565b82815260e0928302850182019282820191908785111562004e2957600080fd5b8387015b8581101562004bf15781818a03121562004e475760008081fd5b62004e5162003eda565b813562004e5e816200477a565b815262004e6d82870162003e35565b86820152604062004e8081840162003e35565b9082015260608281013562004e958162003e12565b90820152608062004ea883820162004497565b9082015260a062004ebb83820162004497565b9082015260c062004ece83820162003e35565b90820152845292840192810162004e2d565b600082601f83011262004ef257600080fd5b8135602062004f0562003feb83620044fd565b82815260059290921b8401810191818101908684111562004f2557600080fd5b8286015b84811015620045ae57803562004f3f8162003e12565b835291830191830162004f29565b600080600080600080600060e0888a03121562004f6957600080fd5b873567ffffffffffffffff8082111562004f8257600080fd5b62004f908b838c0162004d82565b985060208a013591508082111562004fa757600080fd5b62004fb58b838c0162004de4565b975060408a013591508082111562004fcc57600080fd5b62004fda8b838c0162004ee0565b965060608a013591508082111562004ff157600080fd5b62004fff8b838c0162004ee0565b955060808a01359150808211156200501657600080fd5b620050248b838c0162004524565b945060a08a01359150808211156200503b57600080fd5b620050498b838c0162004524565b935060c08a01359150808211156200506057600080fd5b506200506f8a828b0162004524565b91505092959891949750929550565b6bffffffffffffffffffffffff81811683821601908082111562003d575762003d5762004650565b600060208284031215620050b957600080fd5b5051919050565b6bffffffffffffffffffffffff82811682821603908082111562003d575762003d5762004650565b604081526000620050fd604083018662004b84565b828103602084015262005112818587620046a5565b9695505050505050565b8051602080830151919081101562004871577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60209190910360031b1b16919050565b815167ffffffffffffffff8111156200517c576200517c62003eab565b62005194816200518d845462004822565b8462004877565b602080601f831160018114620051ea5760008415620051b35750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555620048c1565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015620052395788860151825594840194600190910190840162005218565b50858210156200527657878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b805169ffffffffffffffffffff8116811462003e4a57600080fd5b600080600080600060a08688031215620052ba57600080fd5b620052c58662005286565b9450602086015193506040860151925060608601519150620052ea6080870162005286565b90509295509295909350565b60ff818116838216019081111562002a505762002a5062004650565b808202811582820484141762002a505762002a5062004650565b600080604083850312156200534057600080fd5b505080516020909101519092909150565b60008154620053608162004822565b808552602060018381168015620053805760018114620053b957620053e9565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b8901019550620053e9565b866000528260002060005b85811015620053e15781548a8201860152908301908401620053c4565b890184019650505b505050505092915050565b60208152600062002a4d602083018462005351565b600082601f8301126200541b57600080fd5b815160206200542e62003feb83620044fd565b82815260059290921b840181019181810190868411156200544e57600080fd5b8286015b84811015620045ae578051835291830191830162005452565b6000602082840312156200547e57600080fd5b815167ffffffffffffffff808211156200549757600080fd5b908301906101008286031215620054ad57600080fd5b620054b762003f06565b8251815260208301516020820152604083015160408201526060830151606082015260808301516080820152620054f160a0840162004719565b60a082015260c0830151828111156200550957600080fd5b620055178782860162005409565b60c08301525060e0830151828111156200553057600080fd5b6200553e8782860162004789565b60e08301525095945050505050565b600081518084526020808501945080840160005b8381101562004b315781518752958201959082019060010162005561565b60408152825160408201526020830151606082015260408301516080820152606083015160a0820152608083015160c082015273ffffffffffffffffffffffffffffffffffffffff60a08401511660e0820152600060c0840151610100808185015250620055f26101408401826200554d565b905060e08501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0848303016101208501526200563082826200417b565b915050828103602084015262002cf9818562005351565b6000826200567e577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfe60c060405234801561001057600080fd5b506040516103ca3803806103ca83398101604081905261002f91610076565b600080546001600160a01b0319166001600160a01b039384161790559181166080521660a0526100b9565b80516001600160a01b038116811461007157600080fd5b919050565b60008060006060848603121561008b57600080fd5b6100948461005a565b92506100a26020850161005a565b91506100b06040850161005a565b90509250925092565b60805160a0516102e76100e36000396000603801526000818160c4015261011701526102e76000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806379188d161461007b578063f00e6a2a146100aa575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e808015610076573d6000f35b3d6000fd5b61008e6100893660046101c1565b6100ee565b6040805192151583526020830191909152015b60405180910390f35b60405173ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001681526020016100a1565b60008054819073ffffffffffffffffffffffffffffffffffffffff16331461011557600080fd5b7f00000000000000000000000000000000000000000000000000000000000000005a91505a61138881101561014957600080fd5b61138881039050856040820482031161016157600080fd5b50803b61016d57600080fd5b6000808551602087016000858af192505a610188908361029a565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080604083850312156101d457600080fd5b82359150602083013567ffffffffffffffff808211156101f357600080fd5b818501915085601f83011261020757600080fd5b81358181111561021957610219610192565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561025f5761025f610192565b8160405282815288602084870101111561027857600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b818103818111156102d4577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b9291505056fea164736f6c6343000813000aa164736f6c6343000813000a",
}

var AutomationRegistryLogicAABI = AutomationRegistryLogicAMetaData.ABI

var AutomationRegistryLogicABin = AutomationRegistryLogicAMetaData.Bin

func DeployAutomationRegistryLogicA(auth *bind.TransactOpts, backend bind.ContractBackend, logicB common.Address) (common.Address, *types.Transaction, *AutomationRegistryLogicA, error) {
	parsed, err := AutomationRegistryLogicAMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AutomationRegistryLogicABin), backend, logicB)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AutomationRegistryLogicA{address: address, abi: *parsed, AutomationRegistryLogicACaller: AutomationRegistryLogicACaller{contract: contract}, AutomationRegistryLogicATransactor: AutomationRegistryLogicATransactor{contract: contract}, AutomationRegistryLogicAFilterer: AutomationRegistryLogicAFilterer{contract: contract}}, nil
}

type AutomationRegistryLogicA struct {
	address common.Address
	abi     abi.ABI
	AutomationRegistryLogicACaller
	AutomationRegistryLogicATransactor
	AutomationRegistryLogicAFilterer
}

type AutomationRegistryLogicACaller struct {
	contract *bind.BoundContract
}

type AutomationRegistryLogicATransactor struct {
	contract *bind.BoundContract
}

type AutomationRegistryLogicAFilterer struct {
	contract *bind.BoundContract
}

type AutomationRegistryLogicASession struct {
	Contract     *AutomationRegistryLogicA
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type AutomationRegistryLogicACallerSession struct {
	Contract *AutomationRegistryLogicACaller
	CallOpts bind.CallOpts
}

type AutomationRegistryLogicATransactorSession struct {
	Contract     *AutomationRegistryLogicATransactor
	TransactOpts bind.TransactOpts
}

type AutomationRegistryLogicARaw struct {
	Contract *AutomationRegistryLogicA
}

type AutomationRegistryLogicACallerRaw struct {
	Contract *AutomationRegistryLogicACaller
}

type AutomationRegistryLogicATransactorRaw struct {
	Contract *AutomationRegistryLogicATransactor
}

func NewAutomationRegistryLogicA(address common.Address, backend bind.ContractBackend) (*AutomationRegistryLogicA, error) {
	abi, err := abi.JSON(strings.NewReader(AutomationRegistryLogicAABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindAutomationRegistryLogicA(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicA{address: address, abi: abi, AutomationRegistryLogicACaller: AutomationRegistryLogicACaller{contract: contract}, AutomationRegistryLogicATransactor: AutomationRegistryLogicATransactor{contract: contract}, AutomationRegistryLogicAFilterer: AutomationRegistryLogicAFilterer{contract: contract}}, nil
}

func NewAutomationRegistryLogicACaller(address common.Address, caller bind.ContractCaller) (*AutomationRegistryLogicACaller, error) {
	contract, err := bindAutomationRegistryLogicA(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicACaller{contract: contract}, nil
}

func NewAutomationRegistryLogicATransactor(address common.Address, transactor bind.ContractTransactor) (*AutomationRegistryLogicATransactor, error) {
	contract, err := bindAutomationRegistryLogicA(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicATransactor{contract: contract}, nil
}

func NewAutomationRegistryLogicAFilterer(address common.Address, filterer bind.ContractFilterer) (*AutomationRegistryLogicAFilterer, error) {
	contract, err := bindAutomationRegistryLogicA(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAFilterer{contract: contract}, nil
}

func bindAutomationRegistryLogicA(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AutomationRegistryLogicAMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicARaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationRegistryLogicA.Contract.AutomationRegistryLogicACaller.contract.Call(opts, result, method, params...)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicARaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.AutomationRegistryLogicATransactor.contract.Transfer(opts)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicARaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.AutomationRegistryLogicATransactor.contract.Transact(opts, method, params...)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicACallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationRegistryLogicA.Contract.contract.Call(opts, result, method, params...)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.contract.Transfer(opts)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.contract.Transact(opts, method, params...)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicACaller) FallbackTo(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicA.contract.Call(opts, &out, "fallbackTo")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicA *AutomationRegistryLogicASession) FallbackTo() (common.Address, error) {
	return _AutomationRegistryLogicA.Contract.FallbackTo(&_AutomationRegistryLogicA.CallOpts)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicACallerSession) FallbackTo() (common.Address, error) {
	return _AutomationRegistryLogicA.Contract.FallbackTo(&_AutomationRegistryLogicA.CallOpts)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicACaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistryLogicA.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistryLogicA *AutomationRegistryLogicASession) Owner() (common.Address, error) {
	return _AutomationRegistryLogicA.Contract.Owner(&_AutomationRegistryLogicA.CallOpts)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicACallerSession) Owner() (common.Address, error) {
	return _AutomationRegistryLogicA.Contract.Owner(&_AutomationRegistryLogicA.CallOpts)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.contract.Transact(opts, "acceptOwnership")
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicASession) AcceptOwnership() (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.AcceptOwnership(&_AutomationRegistryLogicA.TransactOpts)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.AcceptOwnership(&_AutomationRegistryLogicA.TransactOpts)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactor) AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.contract.Transact(opts, "addFunds", id, amount)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicASession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.AddFunds(&_AutomationRegistryLogicA.TransactOpts, id, amount)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactorSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.AddFunds(&_AutomationRegistryLogicA.TransactOpts, id, amount)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactor) CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.contract.Transact(opts, "cancelUpkeep", id)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicASession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.CancelUpkeep(&_AutomationRegistryLogicA.TransactOpts, id)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactorSession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.CancelUpkeep(&_AutomationRegistryLogicA.TransactOpts, id)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactor) CheckCallback(opts *bind.TransactOpts, id *big.Int, values [][]byte, extraData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.contract.Transact(opts, "checkCallback", id, values, extraData)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicASession) CheckCallback(id *big.Int, values [][]byte, extraData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.CheckCallback(&_AutomationRegistryLogicA.TransactOpts, id, values, extraData)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactorSession) CheckCallback(id *big.Int, values [][]byte, extraData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.CheckCallback(&_AutomationRegistryLogicA.TransactOpts, id, values, extraData)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactor) CheckUpkeep(opts *bind.TransactOpts, id *big.Int, triggerData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.contract.Transact(opts, "checkUpkeep", id, triggerData)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicASession) CheckUpkeep(id *big.Int, triggerData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.CheckUpkeep(&_AutomationRegistryLogicA.TransactOpts, id, triggerData)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactorSession) CheckUpkeep(id *big.Int, triggerData []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.CheckUpkeep(&_AutomationRegistryLogicA.TransactOpts, id, triggerData)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactor) CheckUpkeep0(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.contract.Transact(opts, "checkUpkeep0", id)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicASession) CheckUpkeep0(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.CheckUpkeep0(&_AutomationRegistryLogicA.TransactOpts, id)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactorSession) CheckUpkeep0(id *big.Int) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.CheckUpkeep0(&_AutomationRegistryLogicA.TransactOpts, id)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactor) ExecuteCallback(opts *bind.TransactOpts, id *big.Int, payload []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.contract.Transact(opts, "executeCallback", id, payload)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicASession) ExecuteCallback(id *big.Int, payload []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.ExecuteCallback(&_AutomationRegistryLogicA.TransactOpts, id, payload)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactorSession) ExecuteCallback(id *big.Int, payload []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.ExecuteCallback(&_AutomationRegistryLogicA.TransactOpts, id, payload)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactor) MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.contract.Transact(opts, "migrateUpkeeps", ids, destination)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicASession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.MigrateUpkeeps(&_AutomationRegistryLogicA.TransactOpts, ids, destination)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactorSession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.MigrateUpkeeps(&_AutomationRegistryLogicA.TransactOpts, ids, destination)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactor) ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.contract.Transact(opts, "receiveUpkeeps", encodedUpkeeps)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicASession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.ReceiveUpkeeps(&_AutomationRegistryLogicA.TransactOpts, encodedUpkeeps)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactorSession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.ReceiveUpkeeps(&_AutomationRegistryLogicA.TransactOpts, encodedUpkeeps)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactor) RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.contract.Transact(opts, "registerUpkeep", target, gasLimit, admin, triggerType, checkData, triggerConfig, offchainConfig)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicASession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.RegisterUpkeep(&_AutomationRegistryLogicA.TransactOpts, target, gasLimit, admin, triggerType, checkData, triggerConfig, offchainConfig)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactorSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.RegisterUpkeep(&_AutomationRegistryLogicA.TransactOpts, target, gasLimit, admin, triggerType, checkData, triggerConfig, offchainConfig)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactor) RegisterUpkeep0(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.contract.Transact(opts, "registerUpkeep0", target, gasLimit, admin, checkData, offchainConfig)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicASession) RegisterUpkeep0(target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.RegisterUpkeep0(&_AutomationRegistryLogicA.TransactOpts, target, gasLimit, admin, checkData, offchainConfig)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactorSession) RegisterUpkeep0(target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.RegisterUpkeep0(&_AutomationRegistryLogicA.TransactOpts, target, gasLimit, admin, checkData, offchainConfig)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactor) SetUpkeepTriggerConfig(opts *bind.TransactOpts, id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.contract.Transact(opts, "setUpkeepTriggerConfig", id, triggerConfig)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicASession) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.SetUpkeepTriggerConfig(&_AutomationRegistryLogicA.TransactOpts, id, triggerConfig)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactorSession) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.SetUpkeepTriggerConfig(&_AutomationRegistryLogicA.TransactOpts, id, triggerConfig)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.contract.Transact(opts, "transferOwnership", to)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicASession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.TransferOwnership(&_AutomationRegistryLogicA.TransactOpts, to)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.TransferOwnership(&_AutomationRegistryLogicA.TransactOpts, to)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.contract.RawTransact(opts, calldata)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicASession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.Fallback(&_AutomationRegistryLogicA.TransactOpts, calldata)
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicATransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _AutomationRegistryLogicA.Contract.Fallback(&_AutomationRegistryLogicA.TransactOpts, calldata)
}

type AutomationRegistryLogicAAdminPrivilegeConfigSetIterator struct {
	Event *AutomationRegistryLogicAAdminPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAAdminPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAAdminPrivilegeConfigSet)
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
		it.Event = new(AutomationRegistryLogicAAdminPrivilegeConfigSet)
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

func (it *AutomationRegistryLogicAAdminPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAAdminPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAAdminPrivilegeConfigSet struct {
	Admin           common.Address
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*AutomationRegistryLogicAAdminPrivilegeConfigSetIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAAdminPrivilegeConfigSetIterator{contract: _AutomationRegistryLogicA.contract, event: "AdminPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAAdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAAdminPrivilegeConfigSet)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseAdminPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicAAdminPrivilegeConfigSet, error) {
	event := new(AutomationRegistryLogicAAdminPrivilegeConfigSet)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicACancelledUpkeepReportIterator struct {
	Event *AutomationRegistryLogicACancelledUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicACancelledUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicACancelledUpkeepReport)
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
		it.Event = new(AutomationRegistryLogicACancelledUpkeepReport)
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

func (it *AutomationRegistryLogicACancelledUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicACancelledUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicACancelledUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicACancelledUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicACancelledUpkeepReportIterator{contract: _AutomationRegistryLogicA.contract, event: "CancelledUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicACancelledUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicACancelledUpkeepReport)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseCancelledUpkeepReport(log types.Log) (*AutomationRegistryLogicACancelledUpkeepReport, error) {
	event := new(AutomationRegistryLogicACancelledUpkeepReport)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAChainSpecificModuleUpdatedIterator struct {
	Event *AutomationRegistryLogicAChainSpecificModuleUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAChainSpecificModuleUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAChainSpecificModuleUpdated)
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
		it.Event = new(AutomationRegistryLogicAChainSpecificModuleUpdated)
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

func (it *AutomationRegistryLogicAChainSpecificModuleUpdatedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAChainSpecificModuleUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAChainSpecificModuleUpdated struct {
	NewModule common.Address
	Raw       types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterChainSpecificModuleUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicAChainSpecificModuleUpdatedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "ChainSpecificModuleUpdated")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAChainSpecificModuleUpdatedIterator{contract: _AutomationRegistryLogicA.contract, event: "ChainSpecificModuleUpdated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchChainSpecificModuleUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAChainSpecificModuleUpdated) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "ChainSpecificModuleUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAChainSpecificModuleUpdated)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "ChainSpecificModuleUpdated", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseChainSpecificModuleUpdated(log types.Log) (*AutomationRegistryLogicAChainSpecificModuleUpdated, error) {
	event := new(AutomationRegistryLogicAChainSpecificModuleUpdated)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "ChainSpecificModuleUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicADedupKeyAddedIterator struct {
	Event *AutomationRegistryLogicADedupKeyAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicADedupKeyAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicADedupKeyAdded)
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
		it.Event = new(AutomationRegistryLogicADedupKeyAdded)
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

func (it *AutomationRegistryLogicADedupKeyAddedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicADedupKeyAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicADedupKeyAdded struct {
	DedupKey [32]byte
	Raw      types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*AutomationRegistryLogicADedupKeyAddedIterator, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicADedupKeyAddedIterator{contract: _AutomationRegistryLogicA.contract, event: "DedupKeyAdded", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicADedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicADedupKeyAdded)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseDedupKeyAdded(log types.Log) (*AutomationRegistryLogicADedupKeyAdded, error) {
	event := new(AutomationRegistryLogicADedupKeyAdded)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAFundsAddedIterator struct {
	Event *AutomationRegistryLogicAFundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAFundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAFundsAdded)
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
		it.Event = new(AutomationRegistryLogicAFundsAdded)
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

func (it *AutomationRegistryLogicAFundsAddedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAFundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAFundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*AutomationRegistryLogicAFundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAFundsAddedIterator{contract: _AutomationRegistryLogicA.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAFundsAdded)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "FundsAdded", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseFundsAdded(log types.Log) (*AutomationRegistryLogicAFundsAdded, error) {
	event := new(AutomationRegistryLogicAFundsAdded)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAFundsWithdrawnIterator struct {
	Event *AutomationRegistryLogicAFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAFundsWithdrawn)
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
		it.Event = new(AutomationRegistryLogicAFundsWithdrawn)
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

func (it *AutomationRegistryLogicAFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAFundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAFundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAFundsWithdrawnIterator{contract: _AutomationRegistryLogicA.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAFundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAFundsWithdrawn)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseFundsWithdrawn(log types.Log) (*AutomationRegistryLogicAFundsWithdrawn, error) {
	event := new(AutomationRegistryLogicAFundsWithdrawn)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAInsufficientFundsUpkeepReportIterator struct {
	Event *AutomationRegistryLogicAInsufficientFundsUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAInsufficientFundsUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAInsufficientFundsUpkeepReport)
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
		it.Event = new(AutomationRegistryLogicAInsufficientFundsUpkeepReport)
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

func (it *AutomationRegistryLogicAInsufficientFundsUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAInsufficientFundsUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAInsufficientFundsUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAInsufficientFundsUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAInsufficientFundsUpkeepReportIterator{contract: _AutomationRegistryLogicA.contract, event: "InsufficientFundsUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAInsufficientFundsUpkeepReport)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseInsufficientFundsUpkeepReport(log types.Log) (*AutomationRegistryLogicAInsufficientFundsUpkeepReport, error) {
	event := new(AutomationRegistryLogicAInsufficientFundsUpkeepReport)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAOwnerFundsWithdrawnIterator struct {
	Event *AutomationRegistryLogicAOwnerFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAOwnerFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAOwnerFundsWithdrawn)
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
		it.Event = new(AutomationRegistryLogicAOwnerFundsWithdrawn)
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

func (it *AutomationRegistryLogicAOwnerFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAOwnerFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAOwnerFundsWithdrawn struct {
	Amount *big.Int
	Raw    types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*AutomationRegistryLogicAOwnerFundsWithdrawnIterator, error) {

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAOwnerFundsWithdrawnIterator{contract: _AutomationRegistryLogicA.contract, event: "OwnerFundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAOwnerFundsWithdrawn) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAOwnerFundsWithdrawn)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseOwnerFundsWithdrawn(log types.Log) (*AutomationRegistryLogicAOwnerFundsWithdrawn, error) {
	event := new(AutomationRegistryLogicAOwnerFundsWithdrawn)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAOwnershipTransferRequestedIterator struct {
	Event *AutomationRegistryLogicAOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAOwnershipTransferRequested)
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
		it.Event = new(AutomationRegistryLogicAOwnershipTransferRequested)
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

func (it *AutomationRegistryLogicAOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicAOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAOwnershipTransferRequestedIterator{contract: _AutomationRegistryLogicA.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAOwnershipTransferRequested)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseOwnershipTransferRequested(log types.Log) (*AutomationRegistryLogicAOwnershipTransferRequested, error) {
	event := new(AutomationRegistryLogicAOwnershipTransferRequested)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAOwnershipTransferredIterator struct {
	Event *AutomationRegistryLogicAOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAOwnershipTransferred)
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
		it.Event = new(AutomationRegistryLogicAOwnershipTransferred)
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

func (it *AutomationRegistryLogicAOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicAOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAOwnershipTransferredIterator{contract: _AutomationRegistryLogicA.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAOwnershipTransferred)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseOwnershipTransferred(log types.Log) (*AutomationRegistryLogicAOwnershipTransferred, error) {
	event := new(AutomationRegistryLogicAOwnershipTransferred)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAPausedIterator struct {
	Event *AutomationRegistryLogicAPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAPaused)
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
		it.Event = new(AutomationRegistryLogicAPaused)
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

func (it *AutomationRegistryLogicAPausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterPaused(opts *bind.FilterOpts) (*AutomationRegistryLogicAPausedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAPausedIterator{contract: _AutomationRegistryLogicA.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAPaused) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAPaused)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "Paused", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParsePaused(log types.Log) (*AutomationRegistryLogicAPaused, error) {
	event := new(AutomationRegistryLogicAPaused)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAPayeesUpdatedIterator struct {
	Event *AutomationRegistryLogicAPayeesUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAPayeesUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAPayeesUpdated)
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
		it.Event = new(AutomationRegistryLogicAPayeesUpdated)
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

func (it *AutomationRegistryLogicAPayeesUpdatedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAPayeesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAPayeesUpdated struct {
	Transmitters []common.Address
	Payees       []common.Address
	Raw          types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterPayeesUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicAPayeesUpdatedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAPayeesUpdatedIterator{contract: _AutomationRegistryLogicA.contract, event: "PayeesUpdated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAPayeesUpdated) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAPayeesUpdated)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParsePayeesUpdated(log types.Log) (*AutomationRegistryLogicAPayeesUpdated, error) {
	event := new(AutomationRegistryLogicAPayeesUpdated)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAPayeeshipTransferRequestedIterator struct {
	Event *AutomationRegistryLogicAPayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAPayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAPayeeshipTransferRequested)
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
		it.Event = new(AutomationRegistryLogicAPayeeshipTransferRequested)
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

func (it *AutomationRegistryLogicAPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAPayeeshipTransferRequested struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicAPayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAPayeeshipTransferRequestedIterator{contract: _AutomationRegistryLogicA.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAPayeeshipTransferRequested)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParsePayeeshipTransferRequested(log types.Log) (*AutomationRegistryLogicAPayeeshipTransferRequested, error) {
	event := new(AutomationRegistryLogicAPayeeshipTransferRequested)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAPayeeshipTransferredIterator struct {
	Event *AutomationRegistryLogicAPayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAPayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAPayeeshipTransferred)
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
		it.Event = new(AutomationRegistryLogicAPayeeshipTransferred)
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

func (it *AutomationRegistryLogicAPayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAPayeeshipTransferred struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicAPayeeshipTransferredIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAPayeeshipTransferredIterator{contract: _AutomationRegistryLogicA.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAPayeeshipTransferred)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParsePayeeshipTransferred(log types.Log) (*AutomationRegistryLogicAPayeeshipTransferred, error) {
	event := new(AutomationRegistryLogicAPayeeshipTransferred)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAPaymentWithdrawnIterator struct {
	Event *AutomationRegistryLogicAPaymentWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAPaymentWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAPaymentWithdrawn)
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
		it.Event = new(AutomationRegistryLogicAPaymentWithdrawn)
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

func (it *AutomationRegistryLogicAPaymentWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAPaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAPaymentWithdrawn struct {
	Transmitter common.Address
	Amount      *big.Int
	To          common.Address
	Payee       common.Address
	Raw         types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*AutomationRegistryLogicAPaymentWithdrawnIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAPaymentWithdrawnIterator{contract: _AutomationRegistryLogicA.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAPaymentWithdrawn)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParsePaymentWithdrawn(log types.Log) (*AutomationRegistryLogicAPaymentWithdrawn, error) {
	event := new(AutomationRegistryLogicAPaymentWithdrawn)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAReorgedUpkeepReportIterator struct {
	Event *AutomationRegistryLogicAReorgedUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAReorgedUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAReorgedUpkeepReport)
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
		it.Event = new(AutomationRegistryLogicAReorgedUpkeepReport)
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

func (it *AutomationRegistryLogicAReorgedUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAReorgedUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAReorgedUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAReorgedUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAReorgedUpkeepReportIterator{contract: _AutomationRegistryLogicA.contract, event: "ReorgedUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAReorgedUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAReorgedUpkeepReport)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseReorgedUpkeepReport(log types.Log) (*AutomationRegistryLogicAReorgedUpkeepReport, error) {
	event := new(AutomationRegistryLogicAReorgedUpkeepReport)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAStaleUpkeepReportIterator struct {
	Event *AutomationRegistryLogicAStaleUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAStaleUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAStaleUpkeepReport)
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
		it.Event = new(AutomationRegistryLogicAStaleUpkeepReport)
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

func (it *AutomationRegistryLogicAStaleUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAStaleUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAStaleUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAStaleUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAStaleUpkeepReportIterator{contract: _AutomationRegistryLogicA.contract, event: "StaleUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAStaleUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAStaleUpkeepReport)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseStaleUpkeepReport(log types.Log) (*AutomationRegistryLogicAStaleUpkeepReport, error) {
	event := new(AutomationRegistryLogicAStaleUpkeepReport)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAUnpausedIterator struct {
	Event *AutomationRegistryLogicAUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAUnpaused)
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
		it.Event = new(AutomationRegistryLogicAUnpaused)
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

func (it *AutomationRegistryLogicAUnpausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterUnpaused(opts *bind.FilterOpts) (*AutomationRegistryLogicAUnpausedIterator, error) {

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAUnpausedIterator{contract: _AutomationRegistryLogicA.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUnpaused) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAUnpaused)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseUnpaused(log types.Log) (*AutomationRegistryLogicAUnpaused, error) {
	event := new(AutomationRegistryLogicAUnpaused)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAUpkeepAdminTransferRequestedIterator struct {
	Event *AutomationRegistryLogicAUpkeepAdminTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAUpkeepAdminTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAUpkeepAdminTransferRequested)
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
		it.Event = new(AutomationRegistryLogicAUpkeepAdminTransferRequested)
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

func (it *AutomationRegistryLogicAUpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAUpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAUpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicAUpkeepAdminTransferRequestedIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAUpkeepAdminTransferRequestedIterator{contract: _AutomationRegistryLogicA.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAUpkeepAdminTransferRequested)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseUpkeepAdminTransferRequested(log types.Log) (*AutomationRegistryLogicAUpkeepAdminTransferRequested, error) {
	event := new(AutomationRegistryLogicAUpkeepAdminTransferRequested)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAUpkeepAdminTransferredIterator struct {
	Event *AutomationRegistryLogicAUpkeepAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAUpkeepAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAUpkeepAdminTransferred)
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
		it.Event = new(AutomationRegistryLogicAUpkeepAdminTransferred)
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

func (it *AutomationRegistryLogicAUpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAUpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAUpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicAUpkeepAdminTransferredIterator, error) {

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

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAUpkeepAdminTransferredIterator{contract: _AutomationRegistryLogicA.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAUpkeepAdminTransferred)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseUpkeepAdminTransferred(log types.Log) (*AutomationRegistryLogicAUpkeepAdminTransferred, error) {
	event := new(AutomationRegistryLogicAUpkeepAdminTransferred)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAUpkeepCanceledIterator struct {
	Event *AutomationRegistryLogicAUpkeepCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAUpkeepCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAUpkeepCanceled)
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
		it.Event = new(AutomationRegistryLogicAUpkeepCanceled)
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

func (it *AutomationRegistryLogicAUpkeepCanceledIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAUpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAUpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*AutomationRegistryLogicAUpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAUpkeepCanceledIterator{contract: _AutomationRegistryLogicA.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAUpkeepCanceled)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseUpkeepCanceled(log types.Log) (*AutomationRegistryLogicAUpkeepCanceled, error) {
	event := new(AutomationRegistryLogicAUpkeepCanceled)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAUpkeepCheckDataSetIterator struct {
	Event *AutomationRegistryLogicAUpkeepCheckDataSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAUpkeepCheckDataSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAUpkeepCheckDataSet)
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
		it.Event = new(AutomationRegistryLogicAUpkeepCheckDataSet)
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

func (it *AutomationRegistryLogicAUpkeepCheckDataSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAUpkeepCheckDataSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAUpkeepCheckDataSet struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAUpkeepCheckDataSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAUpkeepCheckDataSetIterator{contract: _AutomationRegistryLogicA.contract, event: "UpkeepCheckDataSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepCheckDataSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAUpkeepCheckDataSet)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseUpkeepCheckDataSet(log types.Log) (*AutomationRegistryLogicAUpkeepCheckDataSet, error) {
	event := new(AutomationRegistryLogicAUpkeepCheckDataSet)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAUpkeepGasLimitSetIterator struct {
	Event *AutomationRegistryLogicAUpkeepGasLimitSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAUpkeepGasLimitSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAUpkeepGasLimitSet)
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
		it.Event = new(AutomationRegistryLogicAUpkeepGasLimitSet)
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

func (it *AutomationRegistryLogicAUpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAUpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAUpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAUpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAUpkeepGasLimitSetIterator{contract: _AutomationRegistryLogicA.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAUpkeepGasLimitSet)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseUpkeepGasLimitSet(log types.Log) (*AutomationRegistryLogicAUpkeepGasLimitSet, error) {
	event := new(AutomationRegistryLogicAUpkeepGasLimitSet)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAUpkeepMigratedIterator struct {
	Event *AutomationRegistryLogicAUpkeepMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAUpkeepMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAUpkeepMigrated)
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
		it.Event = new(AutomationRegistryLogicAUpkeepMigrated)
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

func (it *AutomationRegistryLogicAUpkeepMigratedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAUpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAUpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAUpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAUpkeepMigratedIterator{contract: _AutomationRegistryLogicA.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAUpkeepMigrated)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseUpkeepMigrated(log types.Log) (*AutomationRegistryLogicAUpkeepMigrated, error) {
	event := new(AutomationRegistryLogicAUpkeepMigrated)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAUpkeepOffchainConfigSetIterator struct {
	Event *AutomationRegistryLogicAUpkeepOffchainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAUpkeepOffchainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAUpkeepOffchainConfigSet)
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
		it.Event = new(AutomationRegistryLogicAUpkeepOffchainConfigSet)
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

func (it *AutomationRegistryLogicAUpkeepOffchainConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAUpkeepOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAUpkeepOffchainConfigSet struct {
	Id             *big.Int
	OffchainConfig []byte
	Raw            types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAUpkeepOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAUpkeepOffchainConfigSetIterator{contract: _AutomationRegistryLogicA.contract, event: "UpkeepOffchainConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAUpkeepOffchainConfigSet)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseUpkeepOffchainConfigSet(log types.Log) (*AutomationRegistryLogicAUpkeepOffchainConfigSet, error) {
	event := new(AutomationRegistryLogicAUpkeepOffchainConfigSet)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAUpkeepPausedIterator struct {
	Event *AutomationRegistryLogicAUpkeepPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAUpkeepPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAUpkeepPaused)
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
		it.Event = new(AutomationRegistryLogicAUpkeepPaused)
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

func (it *AutomationRegistryLogicAUpkeepPausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAUpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAUpkeepPaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAUpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAUpkeepPausedIterator{contract: _AutomationRegistryLogicA.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAUpkeepPaused)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseUpkeepPaused(log types.Log) (*AutomationRegistryLogicAUpkeepPaused, error) {
	event := new(AutomationRegistryLogicAUpkeepPaused)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAUpkeepPerformedIterator struct {
	Event *AutomationRegistryLogicAUpkeepPerformed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAUpkeepPerformedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAUpkeepPerformed)
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
		it.Event = new(AutomationRegistryLogicAUpkeepPerformed)
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

func (it *AutomationRegistryLogicAUpkeepPerformedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAUpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAUpkeepPerformed struct {
	Id           *big.Int
	Success      bool
	TotalPayment *big.Int
	GasUsed      *big.Int
	GasOverhead  *big.Int
	Trigger      []byte
	Raw          types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*AutomationRegistryLogicAUpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAUpkeepPerformedIterator{contract: _AutomationRegistryLogicA.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAUpkeepPerformed)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseUpkeepPerformed(log types.Log) (*AutomationRegistryLogicAUpkeepPerformed, error) {
	event := new(AutomationRegistryLogicAUpkeepPerformed)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAUpkeepPrivilegeConfigSetIterator struct {
	Event *AutomationRegistryLogicAUpkeepPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAUpkeepPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAUpkeepPrivilegeConfigSet)
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
		it.Event = new(AutomationRegistryLogicAUpkeepPrivilegeConfigSet)
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

func (it *AutomationRegistryLogicAUpkeepPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAUpkeepPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAUpkeepPrivilegeConfigSet struct {
	Id              *big.Int
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAUpkeepPrivilegeConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAUpkeepPrivilegeConfigSetIterator{contract: _AutomationRegistryLogicA.contract, event: "UpkeepPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAUpkeepPrivilegeConfigSet)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseUpkeepPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicAUpkeepPrivilegeConfigSet, error) {
	event := new(AutomationRegistryLogicAUpkeepPrivilegeConfigSet)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAUpkeepReceivedIterator struct {
	Event *AutomationRegistryLogicAUpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAUpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAUpkeepReceived)
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
		it.Event = new(AutomationRegistryLogicAUpkeepReceived)
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

func (it *AutomationRegistryLogicAUpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAUpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAUpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAUpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAUpkeepReceivedIterator{contract: _AutomationRegistryLogicA.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAUpkeepReceived)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseUpkeepReceived(log types.Log) (*AutomationRegistryLogicAUpkeepReceived, error) {
	event := new(AutomationRegistryLogicAUpkeepReceived)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAUpkeepRegisteredIterator struct {
	Event *AutomationRegistryLogicAUpkeepRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAUpkeepRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAUpkeepRegistered)
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
		it.Event = new(AutomationRegistryLogicAUpkeepRegistered)
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

func (it *AutomationRegistryLogicAUpkeepRegisteredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAUpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAUpkeepRegistered struct {
	Id         *big.Int
	PerformGas uint32
	Admin      common.Address
	Raw        types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAUpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAUpkeepRegisteredIterator{contract: _AutomationRegistryLogicA.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAUpkeepRegistered)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseUpkeepRegistered(log types.Log) (*AutomationRegistryLogicAUpkeepRegistered, error) {
	event := new(AutomationRegistryLogicAUpkeepRegistered)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAUpkeepTriggerConfigSetIterator struct {
	Event *AutomationRegistryLogicAUpkeepTriggerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAUpkeepTriggerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAUpkeepTriggerConfigSet)
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
		it.Event = new(AutomationRegistryLogicAUpkeepTriggerConfigSet)
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

func (it *AutomationRegistryLogicAUpkeepTriggerConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAUpkeepTriggerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAUpkeepTriggerConfigSet struct {
	Id            *big.Int
	TriggerConfig []byte
	Raw           types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAUpkeepTriggerConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAUpkeepTriggerConfigSetIterator{contract: _AutomationRegistryLogicA.contract, event: "UpkeepTriggerConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAUpkeepTriggerConfigSet)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseUpkeepTriggerConfigSet(log types.Log) (*AutomationRegistryLogicAUpkeepTriggerConfigSet, error) {
	event := new(AutomationRegistryLogicAUpkeepTriggerConfigSet)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryLogicAUpkeepUnpausedIterator struct {
	Event *AutomationRegistryLogicAUpkeepUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryLogicAUpkeepUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryLogicAUpkeepUnpaused)
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
		it.Event = new(AutomationRegistryLogicAUpkeepUnpaused)
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

func (it *AutomationRegistryLogicAUpkeepUnpausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryLogicAUpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryLogicAUpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAUpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryLogicAUpkeepUnpausedIterator{contract: _AutomationRegistryLogicA.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistryLogicA.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryLogicAUpkeepUnpaused)
				if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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

func (_AutomationRegistryLogicA *AutomationRegistryLogicAFilterer) ParseUpkeepUnpaused(log types.Log) (*AutomationRegistryLogicAUpkeepUnpaused, error) {
	event := new(AutomationRegistryLogicAUpkeepUnpaused)
	if err := _AutomationRegistryLogicA.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicA) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _AutomationRegistryLogicA.abi.Events["AdminPrivilegeConfigSet"].ID:
		return _AutomationRegistryLogicA.ParseAdminPrivilegeConfigSet(log)
	case _AutomationRegistryLogicA.abi.Events["CancelledUpkeepReport"].ID:
		return _AutomationRegistryLogicA.ParseCancelledUpkeepReport(log)
	case _AutomationRegistryLogicA.abi.Events["ChainSpecificModuleUpdated"].ID:
		return _AutomationRegistryLogicA.ParseChainSpecificModuleUpdated(log)
	case _AutomationRegistryLogicA.abi.Events["DedupKeyAdded"].ID:
		return _AutomationRegistryLogicA.ParseDedupKeyAdded(log)
	case _AutomationRegistryLogicA.abi.Events["FundsAdded"].ID:
		return _AutomationRegistryLogicA.ParseFundsAdded(log)
	case _AutomationRegistryLogicA.abi.Events["FundsWithdrawn"].ID:
		return _AutomationRegistryLogicA.ParseFundsWithdrawn(log)
	case _AutomationRegistryLogicA.abi.Events["InsufficientFundsUpkeepReport"].ID:
		return _AutomationRegistryLogicA.ParseInsufficientFundsUpkeepReport(log)
	case _AutomationRegistryLogicA.abi.Events["OwnerFundsWithdrawn"].ID:
		return _AutomationRegistryLogicA.ParseOwnerFundsWithdrawn(log)
	case _AutomationRegistryLogicA.abi.Events["OwnershipTransferRequested"].ID:
		return _AutomationRegistryLogicA.ParseOwnershipTransferRequested(log)
	case _AutomationRegistryLogicA.abi.Events["OwnershipTransferred"].ID:
		return _AutomationRegistryLogicA.ParseOwnershipTransferred(log)
	case _AutomationRegistryLogicA.abi.Events["Paused"].ID:
		return _AutomationRegistryLogicA.ParsePaused(log)
	case _AutomationRegistryLogicA.abi.Events["PayeesUpdated"].ID:
		return _AutomationRegistryLogicA.ParsePayeesUpdated(log)
	case _AutomationRegistryLogicA.abi.Events["PayeeshipTransferRequested"].ID:
		return _AutomationRegistryLogicA.ParsePayeeshipTransferRequested(log)
	case _AutomationRegistryLogicA.abi.Events["PayeeshipTransferred"].ID:
		return _AutomationRegistryLogicA.ParsePayeeshipTransferred(log)
	case _AutomationRegistryLogicA.abi.Events["PaymentWithdrawn"].ID:
		return _AutomationRegistryLogicA.ParsePaymentWithdrawn(log)
	case _AutomationRegistryLogicA.abi.Events["ReorgedUpkeepReport"].ID:
		return _AutomationRegistryLogicA.ParseReorgedUpkeepReport(log)
	case _AutomationRegistryLogicA.abi.Events["StaleUpkeepReport"].ID:
		return _AutomationRegistryLogicA.ParseStaleUpkeepReport(log)
	case _AutomationRegistryLogicA.abi.Events["Unpaused"].ID:
		return _AutomationRegistryLogicA.ParseUnpaused(log)
	case _AutomationRegistryLogicA.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _AutomationRegistryLogicA.ParseUpkeepAdminTransferRequested(log)
	case _AutomationRegistryLogicA.abi.Events["UpkeepAdminTransferred"].ID:
		return _AutomationRegistryLogicA.ParseUpkeepAdminTransferred(log)
	case _AutomationRegistryLogicA.abi.Events["UpkeepCanceled"].ID:
		return _AutomationRegistryLogicA.ParseUpkeepCanceled(log)
	case _AutomationRegistryLogicA.abi.Events["UpkeepCheckDataSet"].ID:
		return _AutomationRegistryLogicA.ParseUpkeepCheckDataSet(log)
	case _AutomationRegistryLogicA.abi.Events["UpkeepGasLimitSet"].ID:
		return _AutomationRegistryLogicA.ParseUpkeepGasLimitSet(log)
	case _AutomationRegistryLogicA.abi.Events["UpkeepMigrated"].ID:
		return _AutomationRegistryLogicA.ParseUpkeepMigrated(log)
	case _AutomationRegistryLogicA.abi.Events["UpkeepOffchainConfigSet"].ID:
		return _AutomationRegistryLogicA.ParseUpkeepOffchainConfigSet(log)
	case _AutomationRegistryLogicA.abi.Events["UpkeepPaused"].ID:
		return _AutomationRegistryLogicA.ParseUpkeepPaused(log)
	case _AutomationRegistryLogicA.abi.Events["UpkeepPerformed"].ID:
		return _AutomationRegistryLogicA.ParseUpkeepPerformed(log)
	case _AutomationRegistryLogicA.abi.Events["UpkeepPrivilegeConfigSet"].ID:
		return _AutomationRegistryLogicA.ParseUpkeepPrivilegeConfigSet(log)
	case _AutomationRegistryLogicA.abi.Events["UpkeepReceived"].ID:
		return _AutomationRegistryLogicA.ParseUpkeepReceived(log)
	case _AutomationRegistryLogicA.abi.Events["UpkeepRegistered"].ID:
		return _AutomationRegistryLogicA.ParseUpkeepRegistered(log)
	case _AutomationRegistryLogicA.abi.Events["UpkeepTriggerConfigSet"].ID:
		return _AutomationRegistryLogicA.ParseUpkeepTriggerConfigSet(log)
	case _AutomationRegistryLogicA.abi.Events["UpkeepUnpaused"].ID:
		return _AutomationRegistryLogicA.ParseUpkeepUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (AutomationRegistryLogicAAdminPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x7c44b4eb59ee7873514e7e43e7718c269d872965938b288aa143befca62f99d2")
}

func (AutomationRegistryLogicACancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd5636")
}

func (AutomationRegistryLogicAChainSpecificModuleUpdated) Topic() common.Hash {
	return common.HexToHash("0xdefc28b11a7980dbe0c49dbbd7055a1584bc8075097d1e8b3b57fb7283df2ad7")
}

func (AutomationRegistryLogicADedupKeyAdded) Topic() common.Hash {
	return common.HexToHash("0xa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f2")
}

func (AutomationRegistryLogicAFundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (AutomationRegistryLogicAFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (AutomationRegistryLogicAInsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e02")
}

func (AutomationRegistryLogicAOwnerFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1")
}

func (AutomationRegistryLogicAOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (AutomationRegistryLogicAOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (AutomationRegistryLogicAPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (AutomationRegistryLogicAPayeesUpdated) Topic() common.Hash {
	return common.HexToHash("0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725")
}

func (AutomationRegistryLogicAPayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (AutomationRegistryLogicAPayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (AutomationRegistryLogicAPaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (AutomationRegistryLogicAReorgedUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301")
}

func (AutomationRegistryLogicAStaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8")
}

func (AutomationRegistryLogicAUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (AutomationRegistryLogicAUpkeepAdminTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35")
}

func (AutomationRegistryLogicAUpkeepAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c")
}

func (AutomationRegistryLogicAUpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (AutomationRegistryLogicAUpkeepCheckDataSet) Topic() common.Hash {
	return common.HexToHash("0xcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d")
}

func (AutomationRegistryLogicAUpkeepGasLimitSet) Topic() common.Hash {
	return common.HexToHash("0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c")
}

func (AutomationRegistryLogicAUpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (AutomationRegistryLogicAUpkeepOffchainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850")
}

func (AutomationRegistryLogicAUpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (AutomationRegistryLogicAUpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b")
}

func (AutomationRegistryLogicAUpkeepPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae7769")
}

func (AutomationRegistryLogicAUpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (AutomationRegistryLogicAUpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (AutomationRegistryLogicAUpkeepTriggerConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664")
}

func (AutomationRegistryLogicAUpkeepUnpaused) Topic() common.Hash {
	return common.HexToHash("0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456")
}

func (_AutomationRegistryLogicA *AutomationRegistryLogicA) Address() common.Address {
	return _AutomationRegistryLogicA.address
}

type AutomationRegistryLogicAInterface interface {
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

	FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*AutomationRegistryLogicAAdminPrivilegeConfigSetIterator, error)

	WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAAdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error)

	ParseAdminPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicAAdminPrivilegeConfigSet, error)

	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicACancelledUpkeepReportIterator, error)

	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicACancelledUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseCancelledUpkeepReport(log types.Log) (*AutomationRegistryLogicACancelledUpkeepReport, error)

	FilterChainSpecificModuleUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicAChainSpecificModuleUpdatedIterator, error)

	WatchChainSpecificModuleUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAChainSpecificModuleUpdated) (event.Subscription, error)

	ParseChainSpecificModuleUpdated(log types.Log) (*AutomationRegistryLogicAChainSpecificModuleUpdated, error)

	FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*AutomationRegistryLogicADedupKeyAddedIterator, error)

	WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicADedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error)

	ParseDedupKeyAdded(log types.Log) (*AutomationRegistryLogicADedupKeyAdded, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*AutomationRegistryLogicAFundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*AutomationRegistryLogicAFundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAFundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAFundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*AutomationRegistryLogicAFundsWithdrawn, error)

	FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAInsufficientFundsUpkeepReportIterator, error)

	WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseInsufficientFundsUpkeepReport(log types.Log) (*AutomationRegistryLogicAInsufficientFundsUpkeepReport, error)

	FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*AutomationRegistryLogicAOwnerFundsWithdrawnIterator, error)

	WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAOwnerFundsWithdrawn) (event.Subscription, error)

	ParseOwnerFundsWithdrawn(log types.Log) (*AutomationRegistryLogicAOwnerFundsWithdrawn, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicAOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*AutomationRegistryLogicAOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryLogicAOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*AutomationRegistryLogicAOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*AutomationRegistryLogicAPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*AutomationRegistryLogicAPaused, error)

	FilterPayeesUpdated(opts *bind.FilterOpts) (*AutomationRegistryLogicAPayeesUpdatedIterator, error)

	WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAPayeesUpdated) (event.Subscription, error)

	ParsePayeesUpdated(log types.Log) (*AutomationRegistryLogicAPayeesUpdated, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicAPayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*AutomationRegistryLogicAPayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryLogicAPayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*AutomationRegistryLogicAPayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*AutomationRegistryLogicAPaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*AutomationRegistryLogicAPaymentWithdrawn, error)

	FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAReorgedUpkeepReportIterator, error)

	WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAReorgedUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseReorgedUpkeepReport(log types.Log) (*AutomationRegistryLogicAReorgedUpkeepReport, error)

	FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAStaleUpkeepReportIterator, error)

	WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAStaleUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseStaleUpkeepReport(log types.Log) (*AutomationRegistryLogicAStaleUpkeepReport, error)

	FilterUnpaused(opts *bind.FilterOpts) (*AutomationRegistryLogicAUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*AutomationRegistryLogicAUnpaused, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicAUpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*AutomationRegistryLogicAUpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryLogicAUpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*AutomationRegistryLogicAUpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*AutomationRegistryLogicAUpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*AutomationRegistryLogicAUpkeepCanceled, error)

	FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAUpkeepCheckDataSetIterator, error)

	WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepCheckDataSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCheckDataSet(log types.Log) (*AutomationRegistryLogicAUpkeepCheckDataSet, error)

	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAUpkeepGasLimitSetIterator, error)

	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepGasLimitSet(log types.Log) (*AutomationRegistryLogicAUpkeepGasLimitSet, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAUpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*AutomationRegistryLogicAUpkeepMigrated, error)

	FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAUpkeepOffchainConfigSetIterator, error)

	WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepOffchainConfigSet(log types.Log) (*AutomationRegistryLogicAUpkeepOffchainConfigSet, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAUpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*AutomationRegistryLogicAUpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*AutomationRegistryLogicAUpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*AutomationRegistryLogicAUpkeepPerformed, error)

	FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAUpkeepPrivilegeConfigSetIterator, error)

	WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPrivilegeConfigSet(log types.Log) (*AutomationRegistryLogicAUpkeepPrivilegeConfigSet, error)

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAUpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*AutomationRegistryLogicAUpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAUpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*AutomationRegistryLogicAUpkeepRegistered, error)

	FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAUpkeepTriggerConfigSetIterator, error)

	WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepTriggerConfigSet(log types.Log) (*AutomationRegistryLogicAUpkeepTriggerConfigSet, error)

	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryLogicAUpkeepUnpausedIterator, error)

	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryLogicAUpkeepUnpaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepUnpaused(log types.Log) (*AutomationRegistryLogicAUpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
