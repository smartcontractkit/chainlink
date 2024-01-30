// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package keeper_registry_logic_a_wrapper_2_2

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
<<<<<<< HEAD
	ABI: "[{\"inputs\":[{\"internalType\":\"contractAutomationRegistryLogicB2_2\",\"name\":\"logicB\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTrigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTriggerType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"AdminPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"DedupKeyAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepTriggerConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumAutomationRegistryBase2_2.UpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"triggerData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumAutomationRegistryBase2_2.UpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumAutomationRegistryBase2_2.UpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"executeCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumAutomationRegistryBase2_2.UpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fallbackTo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"migrateUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedUpkeeps\",\"type\":\"bytes\"}],\"name\":\"receiveUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"enumAutomationRegistryBase2_2.Trigger\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"setUpkeepTriggerConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101606040523480156200001257600080fd5b50604051620062773803806200627783398101604081905262000035916200044b565b80816001600160a01b0316634b4fd03b6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200009b919062000472565b826001600160a01b031663ca30e6036040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000da573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200010091906200044b565b836001600160a01b031663b10b673c6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200013f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200016591906200044b565b846001600160a01b0316636709d0e56040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca91906200044b565b856001600160a01b0316635425d8ac6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000209573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200022f91906200044b565b866001600160a01b031663a08714c06040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200026e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200029491906200044b565b3380600081620002eb5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200031e576200031e8162000387565b50505085600281111562000336576200033662000495565b60e08160028111156200034d576200034d62000495565b9052506001600160a01b0394851660805292841660a05290831660c052821661010052811661012052919091166101405250620004ab9050565b336001600160a01b03821603620003e15760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620002e2565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200044857600080fd5b50565b6000602082840312156200045e57600080fd5b81516200046b8162000432565b9392505050565b6000602082840312156200048557600080fd5b8151600381106200046b57600080fd5b634e487b7160e01b600052602160045260246000fd5b60805160a05160c05160e051610100516101205161014051615d47620005306000396000818161010e01526101a901526000612f590152600081816103e10152611fbd015260008181613590015281816137c601528181613a0e0152613bb60152600061313a0152600061321e015260008181611dff01526123cb0152615d476000f3fe60806040523480156200001157600080fd5b50600436106200010c5760003560e01c806385c1b0ba11620000a5578063c8048022116200006f578063c804802214620002b7578063ce7dc5b414620002ce578063f2fde38b14620002e5578063f7d334ba14620002fc576200010c565b806385c1b0ba14620002535780638da5cb5b146200026a5780638e86139b1462000289578063948108f714620002a0576200010c565b80634ee88d3511620000e75780634ee88d3514620001ef5780636ded9eae146200020657806371791aa0146200021d57806379ba50971462000249576200010c565b806328f32f38146200015457806329c5efad146200017e578063349e8cca14620001a7575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e8080156200014d573d6000f35b3d6000fd5b005b6200016b6200016536600462004243565b62000313565b6040519081526020015b60405180910390f35b620001956200018f36600462004329565b6200068d565b60405162000175949392919062004451565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200162000175565b62000152620002003660046200448e565b62000931565b6200016b62000217366004620044de565b62000999565b620002346200022e36600462004329565b620009ff565b60405162000175979695949392919062004591565b620001526200110c565b6200015262000264366004620045e3565b6200120f565b60005473ffffffffffffffffffffffffffffffffffffffff16620001c9565b620001526200029a36600462004670565b62001e80565b62000152620002b1366004620046d3565b62002208565b62000152620002c836600462004702565b6200249b565b62000195620002df366004620047d8565b62002862565b62000152620002f63660046200484f565b62002928565b620002346200030d36600462004702565b62002940565b6000805473ffffffffffffffffffffffffffffffffffffffff163314801590620003475750620003456009336200297e565b155b156200037f576040517fd48b678b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff89163b620003ce576040517f09ee12d500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b620003d986620029b2565b9050600089307f00000000000000000000000000000000000000000000000000000000000000006040516200040e9062003fd4565b73ffffffffffffffffffffffffffffffffffffffff938416815291831660208301529091166040820152606001604051809103906000f08015801562000458573d6000803e3d6000fd5b5090506200051f826040518060e001604052806000151581526020018c63ffffffff16815260200163ffffffff801681526020018473ffffffffffffffffffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff168152602001600063ffffffff168152508a89898080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508b92508a915062002b569050565b6015805474010000000000000000000000000000000000000000900463ffffffff169060146200054f836200489e565b91906101000a81548163ffffffff021916908363ffffffff16021790555050817fbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d0128a8a604051620005c892919063ffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a2817fcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d8787604051620006049291906200490d565b60405180910390a2817f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664856040516200063e919062004923565b60405180910390a2817f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf48508460405162000678919062004923565b60405180910390a25098975050505050505050565b600060606000806200069e62002f41565b600086815260046020908152604091829020825160e081018452815460ff811615158252610100810463ffffffff90811694830194909452650100000000008104841694820194909452690100000000000000000090930473ffffffffffffffffffffffffffffffffffffffff166060840152600101546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a0840152780100000000000000000000000000000000000000000000000090041660c08201525a9150600080826060015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620007b8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620007de919062004945565b73ffffffffffffffffffffffffffffffffffffffff166014600101600c9054906101000a900463ffffffff1663ffffffff168960405162000820919062004965565b60006040518083038160008787f1925050503d806000811462000860576040519150601f19603f3d011682016040523d82523d6000602084013e62000865565b606091505b50915091505a62000877908562004983565b935081620008a257600060405180602001604052806000815250600796509650965050505062000928565b80806020019051810190620008b89190620049f4565b909750955086620008e657600060405180602001604052806000815250600496509650965050505062000928565b601654865164010000000090910463ffffffff1610156200092457600060405180602001604052806000815250600596509650965050505062000928565b5050505b92959194509250565b6200093c8362002fb3565b6000838152601b602052604090206200095782848362004ae9565b50827f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d566483836040516200098c9291906200490d565b60405180910390a2505050565b6000620009f388888860008989604051806020016040528060008152508a8a8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506200031392505050565b98975050505050505050565b60006060600080600080600062000a1562002f41565b600062000a228a62003069565b905060006012604051806101400160405290816000820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160008201600c9054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160109054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160149054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160189054906101000a900462ffffff1662ffffff1662ffffff16815260200160008201601b9054906101000a900461ffff1661ffff1661ffff16815260200160008201601d9054906101000a900460ff1660ff1660ff16815260200160008201601e9054906101000a900460ff1615151515815260200160008201601f9054906101000a900460ff161515151581526020016001820160009054906101000a900460ff16151515158152505090506000600460008d81526020019081526020016000206040518060e00160405290816000820160009054906101000a900460ff161515151581526020016000820160019054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160059054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160099054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160018201600c9054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff1681526020016001820160189054906101000a900463ffffffff1663ffffffff1663ffffffff168152505090508160e001511562000d61576000604051806020016040528060008152506009600084602001516000808263ffffffff169250995099509950995099509950995050505062001100565b604081015163ffffffff9081161462000db2576000604051806020016040528060008152506001600084602001516000808263ffffffff169250995099509950995099509950995050505062001100565b80511562000df8576000604051806020016040528060008152506002600084602001516000808263ffffffff169250995099509950995099509950995050505062001100565b62000e038262003117565b602083015160165492975090955060009162000e35918591879190640100000000900463ffffffff168a8a8762003309565b9050806bffffffffffffffffffffffff168260a001516bffffffffffffffffffffffff16101562000e9f576000604051806020016040528060008152506006600085602001516000808263ffffffff1692509a509a509a509a509a509a509a505050505062001100565b600062000eae8e868f6200335a565b90505a9850600080846060015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000f06573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000f2c919062004945565b73ffffffffffffffffffffffffffffffffffffffff166014600101600c9054906101000a900463ffffffff1663ffffffff168460405162000f6e919062004965565b60006040518083038160008787f1925050503d806000811462000fae576040519150601f19603f3d011682016040523d82523d6000602084013e62000fb3565b606091505b50915091505a62000fc5908c62004983565b9a5081620010455760165481516801000000000000000090910463ffffffff1610156200102257505060408051602080820190925260008082529490910151939c509a50600899505063ffffffff90911695506200110092505050565b602090940151939b5060039a505063ffffffff9092169650620011009350505050565b808060200190518101906200105b9190620049f4565b909e509c508d6200109c57505060408051602080820190925260008082529490910151939c509a50600499505063ffffffff90911695506200110092505050565b6016548d5164010000000090910463ffffffff161015620010ed57505060408051602080820190925260008082529490910151939c509a50600599505063ffffffff90911695506200110092505050565b505050506020015163ffffffff16945050505b92959891949750929550565b60015473ffffffffffffffffffffffffffffffffffffffff16331462001193576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600173ffffffffffffffffffffffffffffffffffffffff82166000908152601a602052604090205460ff1660038111156200124e576200124e620043e6565b141580156200129a5750600373ffffffffffffffffffffffffffffffffffffffff82166000908152601a602052604090205460ff166003811115620012975762001297620043e6565b14155b15620012d2576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6014546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1662001332576040517fd12d7d8d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008290036200136e576040517f2c2fc94100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160e081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c081018290526000808567ffffffffffffffff811115620013c557620013c5620040ca565b604051908082528060200260200182016040528015620013ef578160200160208202803683370190505b50905060008667ffffffffffffffff811115620014105762001410620040ca565b6040519080825280602002602001820160405280156200149757816020015b6040805160e08101825260008082526020808301829052928201819052606082018190526080820181905260a0820181905260c082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816200142f5790505b50905060008767ffffffffffffffff811115620014b857620014b8620040ca565b604051908082528060200260200182016040528015620014ed57816020015b6060815260200190600190039081620014d75790505b50905060008867ffffffffffffffff8111156200150e576200150e620040ca565b6040519080825280602002602001820160405280156200154357816020015b60608152602001906001900390816200152d5790505b50905060008967ffffffffffffffff811115620015645762001564620040ca565b6040519080825280602002602001820160405280156200159957816020015b6060815260200190600190039081620015835790505b50905060005b8a81101562001b7d578b8b82818110620015bd57620015bd62004c11565b60209081029290920135600081815260048452604090819020815160e081018352815460ff811615158252610100810463ffffffff90811697830197909752650100000000008104871693820193909352690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff166060830152600101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490931660c08401529a509098506200169c90508962002fb3565b60608801516040517f1a5da6c800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8c8116600483015290911690631a5da6c890602401600060405180830381600087803b1580156200170c57600080fd5b505af115801562001721573d6000803e3d6000fd5b50505050878582815181106200173b576200173b62004c11565b6020026020010181905250600560008a815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff168682815181106200178f576200178f62004c11565b73ffffffffffffffffffffffffffffffffffffffff90921660209283029190910182015260008a81526007909152604090208054620017ce9062004a41565b80601f0160208091040260200160405190810160405280929190818152602001828054620017fc9062004a41565b80156200184d5780601f1062001821576101008083540402835291602001916200184d565b820191906000526020600020905b8154815290600101906020018083116200182f57829003601f168201915b505050505084828151811062001867576200186762004c11565b6020026020010181905250601b60008a81526020019081526020016000208054620018929062004a41565b80601f0160208091040260200160405190810160405280929190818152602001828054620018c09062004a41565b8015620019115780601f10620018e55761010080835404028352916020019162001911565b820191906000526020600020905b815481529060010190602001808311620018f357829003601f168201915b50505050508382815181106200192b576200192b62004c11565b6020026020010181905250601c60008a81526020019081526020016000208054620019569062004a41565b80601f0160208091040260200160405190810160405280929190818152602001828054620019849062004a41565b8015620019d55780601f10620019a957610100808354040283529160200191620019d5565b820191906000526020600020905b815481529060010190602001808311620019b757829003601f168201915b5050505050828281518110620019ef57620019ef62004c11565b60200260200101819052508760a001516bffffffffffffffffffffffff168762001a1a919062004c40565b60008a815260046020908152604080832080547fffffff000000000000000000000000000000000000000000000000000000000016815560010180547fffffffff000000000000000000000000000000000000000000000000000000001690556007909152812091985062001a90919062003fe2565b6000898152601b6020526040812062001aa99162003fe2565b6000898152601c6020526040812062001ac29162003fe2565b600089815260066020526040902080547fffffffffffffffffffffffff000000000000000000000000000000000000000016905562001b0360028a6200357c565b5060a0880151604080516bffffffffffffffffffffffff909216825273ffffffffffffffffffffffffffffffffffffffff8c1660208301528a917fb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff910160405180910390a28062001b748162004c56565b9150506200159f565b508560195462001b8e919062004983565b60195560008b8b868167ffffffffffffffff81111562001bb25762001bb2620040ca565b60405190808252806020026020018201604052801562001bdc578160200160208202803683370190505b508988888860405160200162001bfa98979695949392919062004dfc565b60405160208183030381529060405290508973ffffffffffffffffffffffffffffffffffffffff16638e86139b6014600001600c9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663c71249ab60038e73ffffffffffffffffffffffffffffffffffffffff1663aab9edd66040518163ffffffff1660e01b8152600401602060405180830381865afa15801562001cb6573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001cdc919062004ecc565b866040518463ffffffff1660e01b815260040162001cfd9392919062004ef1565b600060405180830381865afa15801562001d1b573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405262001d63919081019062004f18565b6040518263ffffffff1660e01b815260040162001d81919062004923565b600060405180830381600087803b15801562001d9c57600080fd5b505af115801562001db1573d6000803e3d6000fd5b50506040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8d81166004830152602482018b90527f000000000000000000000000000000000000000000000000000000000000000016925063a9059cbb91506044016020604051808303816000875af115801562001e4b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001e71919062004f51565b50505050505050505050505050565b6002336000908152601a602052604090205460ff16600381111562001ea95762001ea9620043e6565b1415801562001edf57506003336000908152601a602052604090205460ff16600381111562001edc5762001edc620043e6565b14155b1562001f17576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080808080808062001f2d888a018a62005149565b965096509650965096509650965060005b8751811015620021fc57600073ffffffffffffffffffffffffffffffffffffffff1687828151811062001f755762001f7562004c11565b60200260200101516060015173ffffffffffffffffffffffffffffffffffffffff1603620020895785818151811062001fb25762001fb262004c11565b6020026020010151307f000000000000000000000000000000000000000000000000000000000000000060405162001fea9062003fd4565b73ffffffffffffffffffffffffffffffffffffffff938416815291831660208301529091166040820152606001604051809103906000f08015801562002034573d6000803e3d6000fd5b508782815181106200204a576200204a62004c11565b60200260200101516060019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250505b62002141888281518110620020a257620020a262004c11565b6020026020010151888381518110620020bf57620020bf62004c11565b6020026020010151878481518110620020dc57620020dc62004c11565b6020026020010151878581518110620020f957620020f962004c11565b602002602001015187868151811062002116576200211662004c11565b602002602001015187878151811062002133576200213362004c11565b602002602001015162002b56565b87818151811062002156576200215662004c11565b60200260200101517f74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a7188838151811062002194576200219462004c11565b602002602001015160a0015133604051620021df9291906bffffffffffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a280620021f38162004c56565b91505062001f3e565b50505050505050505050565b600082815260046020908152604091829020825160e081018452815460ff81161515825263ffffffff6101008204811694830194909452650100000000008104841694820185905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004821660c0820152911462002306576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818160a001516200231891906200527a565b600084815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff93841602179055601954620023809184169062004c40565b6019556040517f23b872dd0000000000000000000000000000000000000000000000000000000081523360048201523060248201526bffffffffffffffffffffffff831660448201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906323b872dd906064016020604051808303816000875af11580156200242a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062002450919062004f51565b506040516bffffffffffffffffffffffff83168152339084907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a3505050565b6000818152600460209081526040808320815160e081018352815460ff81161515825263ffffffff6101008204811695830195909552650100000000008104851693820184905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004831660c082015292911415906200258460005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16149050818015620025df5750808015620025dd5750620025d06200358a565b836040015163ffffffff16115b155b1562002617576040517ffbc0357800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b801580156200264a575060008481526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314155b1562002682576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006200268e6200358a565b905081620026a657620026a360328262004c40565b90505b6000858152600460205260409020805463ffffffff80841665010000000000027fffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffffff90921691909117909155620027029060029087906200357c16565b5060145460808501516bffffffffffffffffffffffff91821691600091168211156200276b576080860151620027399083620052a2565b90508560a001516bffffffffffffffffffffffff16816bffffffffffffffffffffffff1611156200276b575060a08501515b808660a001516200277d9190620052a2565b600088815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff93841602179055601554620027e5918391166200527a565b601580547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff9290921691909117905560405167ffffffffffffffff84169088907f91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f79118190600090a350505050505050565b600060606000806000634b56a42e60e01b8888886040516024016200288a93929190620052ca565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915290506200291589826200068d565b929c919b50995090975095505050505050565b6200293262003646565b6200293d81620036c9565b50565b600060606000806000806000620029678860405180602001604052806000815250620009ff565b959e949d50929b5090995097509550909350915050565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415155b90505b92915050565b6000806000620029d96001620029c76200358a565b620029d3919062004983565b620037c0565b601554604080516020810193909352309083015274010000000000000000000000000000000000000000900463ffffffff166060820152608001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201209083015201604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052905060045b600f81101562002ae5578282828151811062002aa15762002aa162004c11565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508062002adc8162004c56565b91505062002a81565b5083600181111562002afb5762002afb620043e6565b60f81b81600f8151811062002b145762002b1462004c11565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535062002b4e81620052fe565b949350505050565b6012547e01000000000000000000000000000000000000000000000000000000000000900460ff161562002bb6576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601654835163ffffffff909116101562002bfc576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108fc856020015163ffffffff16108062002c3a5750601554602086015163ffffffff70010000000000000000000000000000000090920482169116115b1562002c72576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000868152600460205260409020546901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff161562002cdc576040517f6e3b930b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000868152600460209081526040808320885181548a8501518b85015160608d01517fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000009093169315157fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000ff169390931761010063ffffffff92831602177fffffff000000000000000000000000000000000000000000000000ffffffffff1665010000000000938216939093027fffffff0000000000000000000000000000000000000000ffffffffffffffffff1692909217690100000000000000000073ffffffffffffffffffffffffffffffffffffffff9283160217835560808b01516001909301805460a08d015160c08e01516bffffffffffffffffffffffff9687167fffffffffffffffff000000000000000000000000000000000000000000000000909316929092176c010000000000000000000000009690911695909502949094177fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff1678010000000000000000000000000000000000000000000000009490931693909302919091179091556005835281842080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169189169190911790556007909152902062002ecf848262005341565b508460a001516bffffffffffffffffffffffff1660195462002ef2919062004c40565b6019556000868152601b6020526040902062002f0f838262005341565b506000868152601c6020526040902062002f2a828262005341565b5062002f3860028762003928565b50505050505050565b3273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161462002fb1576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60008181526005602052604090205473ffffffffffffffffffffffffffffffffffffffff16331462003011576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526004602052604090205465010000000000900463ffffffff908116146200293d576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000818160045b600f811015620030fe577fff000000000000000000000000000000000000000000000000000000000000008216838260208110620030b257620030b262004c11565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614620030e957506000949350505050565b80620030f58162004c56565b91505062003070565b5081600f1a600181111562002b4e5762002b4e620043e6565b6000806000836080015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa158015620031a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620031ca919062005483565b5094509092505050600081131580620031e257508142105b80620032075750828015620032075750620031fe824262004983565b8463ffffffff16105b15620032185760175495506200321c565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa15801562003288573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620032ae919062005483565b5094509092505050600081131580620032c657508142105b80620032eb5750828015620032eb5750620032e2824262004983565b8463ffffffff16105b15620032fc57601854945062003300565b8094505b50505050915091565b6000806200331d88878b60c0015162003936565b90506000806200333a8b8a63ffffffff16858a8a60018b620039d5565b90925090506200334b81836200527a565b9b9a5050505050505050505050565b60606000836001811115620033735762003373620043e6565b0362003440576000848152600760205260409081902090517f6e04ff0d0000000000000000000000000000000000000000000000000000000091620033bb916024016200557b565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152905062003575565b6001836001811115620034575762003457620043e6565b036200354357600082806020019051810190620034759190620055f2565b6000868152600760205260409081902090519192507f40691db40000000000000000000000000000000000000000000000000000000091620034bc91849160240162005706565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091529150620035759050565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b9392505050565b6000620029a9838362003e77565b600060017f00000000000000000000000000000000000000000000000000000000000000006002811115620035c357620035c3620043e6565b036200364157606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562003616573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200363c9190620057ce565b905090565b504390565b60005473ffffffffffffffffffffffffffffffffffffffff16331462002fb1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016200118a565b3373ffffffffffffffffffffffffffffffffffffffff8216036200374a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200118a565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060017f00000000000000000000000000000000000000000000000000000000000000006002811115620037f957620037f9620043e6565b036200391e576000606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200384e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620038749190620057ce565b9050808310158062003892575061010062003890848362004983565b115b15620038a15750600092915050565b6040517f2b407a8200000000000000000000000000000000000000000000000000000000815260048101849052606490632b407a8290602401602060405180830381865afa158015620038f8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620035759190620057ce565b504090565b919050565b6000620029a9838362003f82565b600080808560018111156200394f576200394f620043e6565b0362003960575062015f9062003983565b6001856001811115620039775762003977620043e6565b036200354357506201adb05b6200399663ffffffff85166014620057e8565b620039a384600162005802565b620039b49060ff16611d4c620057e8565b620039c0908362004c40565b620039cc919062004c40565b95945050505050565b60008060008960a0015161ffff1687620039f09190620057e8565b9050838015620039ff5750803a105b1562003a0857503a5b600060027f0000000000000000000000000000000000000000000000000000000000000000600281111562003a415762003a41620043e6565b0362003bb257604080516000815260208101909152851562003aa55760003660405180608001604052806048815260200162005cf36048913960405160200162003a8e939291906200581e565b604051602081830303815290604052905062003b13565b60165462003ac390640100000000900463ffffffff16600462005847565b63ffffffff1667ffffffffffffffff81111562003ae45762003ae4620040ca565b6040519080825280601f01601f19166020018201604052801562003b0f576020820181803683370190505b5090505b6040517f49948e0e00000000000000000000000000000000000000000000000000000000815273420000000000000000000000000000000000000f906349948e0e9062003b6590849060040162004923565b602060405180830381865afa15801562003b83573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003ba99190620057ce565b91505062003d1c565b60017f0000000000000000000000000000000000000000000000000000000000000000600281111562003be95762003be9620043e6565b0362003d1c57841562003c7157606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562003c43573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003c699190620057ce565b905062003d1c565b6000606c73ffffffffffffffffffffffffffffffffffffffff166341b247a86040518163ffffffff1660e01b815260040160c060405180830381865afa15801562003cc0573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003ce6919062005872565b505060165492945062003d0b93505050640100000000900463ffffffff1682620057e8565b62003d18906010620057e8565b9150505b8462003d3b57808b60a0015161ffff1662003d389190620057e8565b90505b62003d4b61ffff871682620058bd565b90506000878262003d5d8c8e62004c40565b62003d699086620057e8565b62003d75919062004c40565b62003d8990670de0b6b3a7640000620057e8565b62003d959190620058bd565b905060008c6040015163ffffffff1664e8d4a5100062003db69190620057e8565b898e6020015163ffffffff16858f8862003dd19190620057e8565b62003ddd919062004c40565b62003ded90633b9aca00620057e8565b62003df99190620057e8565b62003e059190620058bd565b62003e11919062004c40565b90506b033b2e3c9fd0803ce800000062003e2c828462004c40565b111562003e65576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b6000818152600183016020526040812054801562003f7057600062003e9e60018362004983565b855490915060009062003eb49060019062004983565b905081811462003f2057600086600001828154811062003ed85762003ed862004c11565b906000526020600020015490508087600001848154811062003efe5762003efe62004c11565b6000918252602080832090910192909255918252600188019052604090208390555b855486908062003f345762003f34620058f9565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050620029ac565b6000915050620029ac565b5092915050565b600081815260018301602052604081205462003fcb57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155620029ac565b506000620029ac565b6103ca806200592983390190565b50805462003ff09062004a41565b6000825580601f1062004001575050565b601f0160209004906000526020600020908101906200293d91905b808211156200403257600081556001016200401c565b5090565b73ffffffffffffffffffffffffffffffffffffffff811681146200293d57600080fd5b803563ffffffff811681146200392357600080fd5b8035600281106200392357600080fd5b60008083601f8401126200409157600080fd5b50813567ffffffffffffffff811115620040aa57600080fd5b602083019150836020828501011115620040c357600080fd5b9250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160e0810167ffffffffffffffff811182821017156200411f576200411f620040ca565b60405290565b604051610100810167ffffffffffffffff811182821017156200411f576200411f620040ca565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715620041965762004196620040ca565b604052919050565b600067ffffffffffffffff821115620041bb57620041bb620040ca565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f830112620041f957600080fd5b8135620042106200420a826200419e565b6200414c565b8181528460208386010111156200422657600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060008060e0898b0312156200426057600080fd5b88356200426d8162004036565b97506200427d60208a0162004059565b965060408901356200428f8162004036565b95506200429f60608a016200406e565b9450608089013567ffffffffffffffff80821115620042bd57600080fd5b620042cb8c838d016200407e565b909650945060a08b0135915080821115620042e557600080fd5b620042f38c838d01620041e7565b935060c08b01359150808211156200430a57600080fd5b50620043198b828c01620041e7565b9150509295985092959890939650565b600080604083850312156200433d57600080fd5b82359150602083013567ffffffffffffffff8111156200435c57600080fd5b6200436a85828601620041e7565b9150509250929050565b60005b838110156200439157818101518382015260200162004377565b50506000910152565b60008151808452620043b481602086016020860162004374565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600a81106200444d577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b9052565b84151581526080602082015260006200446e60808301866200439a565b90506200447f604083018562004415565b82606083015295945050505050565b600080600060408486031215620044a457600080fd5b83359250602084013567ffffffffffffffff811115620044c357600080fd5b620044d1868287016200407e565b9497909650939450505050565b600080600080600080600060a0888a031215620044fa57600080fd5b8735620045078162004036565b9650620045176020890162004059565b95506040880135620045298162004036565b9450606088013567ffffffffffffffff808211156200454757600080fd5b620045558b838c016200407e565b909650945060808a01359150808211156200456f57600080fd5b506200457e8a828b016200407e565b989b979a50959850939692959293505050565b871515815260e060208201526000620045ae60e08301896200439a565b9050620045bf604083018862004415565b8560608301528460808301528360a08301528260c083015298975050505050505050565b600080600060408486031215620045f957600080fd5b833567ffffffffffffffff808211156200461257600080fd5b818601915086601f8301126200462757600080fd5b8135818111156200463757600080fd5b8760208260051b85010111156200464d57600080fd5b60209283019550935050840135620046658162004036565b809150509250925092565b600080602083850312156200468457600080fd5b823567ffffffffffffffff8111156200469c57600080fd5b620046aa858286016200407e565b90969095509350505050565b80356bffffffffffffffffffffffff811681146200392357600080fd5b60008060408385031215620046e757600080fd5b82359150620046f960208401620046b6565b90509250929050565b6000602082840312156200471557600080fd5b5035919050565b600067ffffffffffffffff821115620047395762004739620040ca565b5060051b60200190565b600082601f8301126200475557600080fd5b81356020620047686200420a836200471c565b82815260059290921b840181019181810190868411156200478857600080fd5b8286015b84811015620047cd57803567ffffffffffffffff811115620047ae5760008081fd5b620047be8986838b0101620041e7565b8452509183019183016200478c565b509695505050505050565b60008060008060608587031215620047ef57600080fd5b84359350602085013567ffffffffffffffff808211156200480f57600080fd5b6200481d8883890162004743565b945060408701359150808211156200483457600080fd5b5062004843878288016200407e565b95989497509550505050565b6000602082840312156200486257600080fd5b8135620035758162004036565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600063ffffffff808316818103620048ba57620048ba6200486f565b6001019392505050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b60208152600062002b4e602083018486620048c4565b602081526000620029a960208301846200439a565b8051620039238162004036565b6000602082840312156200495857600080fd5b8151620035758162004036565b600082516200497981846020870162004374565b9190910192915050565b81810381811115620029ac57620029ac6200486f565b80151581146200293d57600080fd5b600082601f830112620049ba57600080fd5b8151620049cb6200420a826200419e565b818152846020838601011115620049e157600080fd5b62002b4e82602083016020870162004374565b6000806040838503121562004a0857600080fd5b825162004a158162004999565b602084015190925067ffffffffffffffff81111562004a3357600080fd5b6200436a85828601620049a8565b600181811c9082168062004a5657607f821691505b60208210810362004a90577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f82111562004ae457600081815260208120601f850160051c8101602086101562004abf5750805b601f850160051c820191505b8181101562004ae05782815560010162004acb565b5050505b505050565b67ffffffffffffffff83111562004b045762004b04620040ca565b62004b1c8362004b15835462004a41565b8362004a96565b6000601f84116001811462004b71576000851562004b3a5750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b17835562004c0a565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b8281101562004bc2578685013582556020948501946001909201910162004ba0565b508682101562004bfe577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b80820180821115620029ac57620029ac6200486f565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820362004c8a5762004c8a6200486f565b5060010190565b600081518084526020808501945080840160005b8381101562004d505781518051151588528381015163ffffffff908116858a01526040808301519091169089015260608082015173ffffffffffffffffffffffffffffffffffffffff16908901526080808201516bffffffffffffffffffffffff169089015260a08082015162004d2b828b01826bffffffffffffffffffffffff169052565b505060c09081015163ffffffff169088015260e0909601959082019060010162004ca5565b509495945050505050565b600081518084526020808501945080840160005b8381101562004d5057815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010162004d6f565b600081518084526020808501808196508360051b8101915082860160005b8581101562004def57828403895262004ddc8483516200439a565b9885019893509084019060010162004dc1565b5091979650505050505050565b60e081528760e082015260006101007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8a111562004e3957600080fd5b8960051b808c8386013783018381038201602085015262004e5d8282018b62004c91565b915050828103604084015262004e74818962004d5b565b9050828103606084015262004e8a818862004d5b565b9050828103608084015262004ea0818762004da3565b905082810360a084015262004eb6818662004da3565b905082810360c08401526200334b818562004da3565b60006020828403121562004edf57600080fd5b815160ff811681146200357557600080fd5b60ff8416815260ff83166020820152606060408201526000620039cc60608301846200439a565b60006020828403121562004f2b57600080fd5b815167ffffffffffffffff81111562004f4357600080fd5b62002b4e84828501620049a8565b60006020828403121562004f6457600080fd5b8151620035758162004999565b600082601f83011262004f8357600080fd5b8135602062004f966200420a836200471c565b82815260059290921b8401810191818101908684111562004fb657600080fd5b8286015b84811015620047cd578035835291830191830162004fba565b600082601f83011262004fe557600080fd5b8135602062004ff86200420a836200471c565b82815260e092830285018201928282019190878511156200501857600080fd5b8387015b85811015620050cf5781818a031215620050365760008081fd5b62005040620040f9565b81356200504d8162004999565b81526200505c82870162004059565b8682015260406200506f81840162004059565b90820152606082810135620050848162004036565b90820152608062005097838201620046b6565b9082015260a0620050aa838201620046b6565b9082015260c0620050bd83820162004059565b9082015284529284019281016200501c565b5090979650505050505050565b600082601f830112620050ee57600080fd5b81356020620051016200420a836200471c565b82815260059290921b840181019181810190868411156200512157600080fd5b8286015b84811015620047cd5780356200513b8162004036565b835291830191830162005125565b600080600080600080600060e0888a0312156200516557600080fd5b873567ffffffffffffffff808211156200517e57600080fd5b6200518c8b838c0162004f71565b985060208a0135915080821115620051a357600080fd5b620051b18b838c0162004fd3565b975060408a0135915080821115620051c857600080fd5b620051d68b838c01620050dc565b965060608a0135915080821115620051ed57600080fd5b620051fb8b838c01620050dc565b955060808a01359150808211156200521257600080fd5b620052208b838c0162004743565b945060a08a01359150808211156200523757600080fd5b620052458b838c0162004743565b935060c08a01359150808211156200525c57600080fd5b506200526b8a828b0162004743565b91505092959891949750929550565b6bffffffffffffffffffffffff81811683821601908082111562003f7b5762003f7b6200486f565b6bffffffffffffffffffffffff82811682821603908082111562003f7b5762003f7b6200486f565b604081526000620052df604083018662004da3565b8281036020840152620052f4818587620048c4565b9695505050505050565b8051602080830151919081101562004a90577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60209190910360031b1b16919050565b815167ffffffffffffffff8111156200535e576200535e620040ca565b62005376816200536f845462004a41565b8462004a96565b602080601f831160018114620053cc5760008415620053955750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b17855562004ae0565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b828110156200541b57888601518255948401946001909101908401620053fa565b50858210156200545857878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b805169ffffffffffffffffffff811681146200392357600080fd5b600080600080600060a086880312156200549c57600080fd5b620054a78662005468565b9450602086015193506040860151925060608601519150620054cc6080870162005468565b90509295509295909350565b60008154620054e78162004a41565b808552602060018381168015620055075760018114620055405762005570565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b890101955062005570565b866000528260002060005b85811015620055685781548a82018601529083019084016200554b565b890184019650505b505050505092915050565b602081526000620029a96020830184620054d8565b600082601f830112620055a257600080fd5b81516020620055b56200420a836200471c565b82815260059290921b84018101918181019086841115620055d557600080fd5b8286015b84811015620047cd5780518352918301918301620055d9565b6000602082840312156200560557600080fd5b815167ffffffffffffffff808211156200561e57600080fd5b9083019061010082860312156200563457600080fd5b6200563e62004125565b82518152602083015160208201526040830151604082015260608301516060820152608083015160808201526200567860a0840162004938565b60a082015260c0830151828111156200569057600080fd5b6200569e8782860162005590565b60c08301525060e083015182811115620056b757600080fd5b620056c587828601620049a8565b60e08301525095945050505050565b600081518084526020808501945080840160005b8381101562004d5057815187529582019590820190600101620056e8565b60408152825160408201526020830151606082015260408301516080820152606083015160a0820152608083015160c082015273ffffffffffffffffffffffffffffffffffffffff60a08401511660e0820152600060c084015161010080818501525062005779610140840182620056d4565b905060e08501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc084830301610120850152620057b782826200439a565b9150508281036020840152620039cc8185620054d8565b600060208284031215620057e157600080fd5b5051919050565b8082028115828204841417620029ac57620029ac6200486f565b60ff8181168382160190811115620029ac57620029ac6200486f565b8284823760008382016000815283516200583d81836020880162004374565b0195945050505050565b63ffffffff8181168382160280821691908281146200586a576200586a6200486f565b505092915050565b60008060008060008060c087890312156200588c57600080fd5b865195506020870151945060408701519350606087015192506080870151915060a087015190509295509295509295565b600082620058f4577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfe60c060405234801561001057600080fd5b506040516103ca3803806103ca83398101604081905261002f91610076565b600080546001600160a01b0319166001600160a01b039384161790559181166080521660a0526100b9565b80516001600160a01b038116811461007157600080fd5b919050565b60008060006060848603121561008b57600080fd5b6100948461005a565b92506100a26020850161005a565b91506100b06040850161005a565b90509250925092565b60805160a0516102e76100e36000396000603801526000818160c4015261011701526102e76000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806379188d161461007b578063f00e6a2a146100aa575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e808015610076573d6000f35b3d6000fd5b61008e6100893660046101c1565b6100ee565b6040805192151583526020830191909152015b60405180910390f35b60405173ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001681526020016100a1565b60008054819073ffffffffffffffffffffffffffffffffffffffff16331461011557600080fd5b7f00000000000000000000000000000000000000000000000000000000000000005a91505a61138881101561014957600080fd5b61138881039050856040820482031161016157600080fd5b50803b61016d57600080fd5b6000808551602087016000858af192505a610188908361029a565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080604083850312156101d457600080fd5b82359150602083013567ffffffffffffffff808211156101f357600080fd5b818501915085601f83011261020757600080fd5b81358181111561021957610219610192565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561025f5761025f610192565b8160405282815288602084870101111561027857600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b818103818111156102d4577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b9291505056fea164736f6c6343000813000a307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000813000a",
=======
	ABI: "[{\"inputs\":[{\"internalType\":\"contractAutomationRegistryLogicB2_2\",\"name\":\"logicB\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTrigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTriggerType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"AdminPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newModule\",\"type\":\"address\"}],\"name\":\"ChainSpecificModuleUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"DedupKeyAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepTriggerConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumAutomationRegistryBase2_2.UpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"triggerData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumAutomationRegistryBase2_2.UpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumAutomationRegistryBase2_2.UpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"executeCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumAutomationRegistryBase2_2.UpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fallbackTo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"migrateUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedUpkeeps\",\"type\":\"bytes\"}],\"name\":\"receiveUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"enumAutomationRegistryBase2_2.Trigger\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"setUpkeepTriggerConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	Bin: "0x6101206040523480156200001257600080fd5b5060405162005e0738038062005e07833981016040819052620000359162000345565b80816001600160a01b031663ca30e6036040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200009b919062000345565b826001600160a01b031663b10b673c6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000da573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000100919062000345565b836001600160a01b0316636709d0e56040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200013f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000165919062000345565b846001600160a01b0316635425d8ac6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca919062000345565b3380600081620002215760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200025457620002548162000281565b5050506001600160a01b0393841660805291831660a052821660c052811660e0521661010052506200036c565b336001600160a01b03821603620002db5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000218565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200034257600080fd5b50565b6000602082840312156200035857600080fd5b815162000365816200032c565b9392505050565b60805160a05160c05160e05161010051615a41620003c66000396000818161010e01526101a90152600081816103e10152612013015260006132f7015260006133db015260008181611e5501526124210152615a416000f3fe60806040523480156200001157600080fd5b50600436106200010c5760003560e01c806385c1b0ba11620000a5578063c8048022116200006f578063c804802214620002b7578063ce7dc5b414620002ce578063f2fde38b14620002e5578063f7d334ba14620002fc576200010c565b806385c1b0ba14620002535780638da5cb5b146200026a5780638e86139b1462000289578063948108f714620002a0576200010c565b80634ee88d3511620000e75780634ee88d3514620001ef5780636ded9eae146200020657806371791aa0146200021d57806379ba50971462000249576200010c565b806328f32f38146200015457806329c5efad146200017e578063349e8cca14620001a7575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e8080156200014d573d6000f35b3d6000fd5b005b6200016b6200016536600462003ffe565b62000313565b6040519081526020015b60405180910390f35b620001956200018f366004620040e4565b6200068d565b6040516200017594939291906200420c565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200162000175565b620001526200020036600462004249565b62000931565b6200016b6200021736600462004299565b62000999565b620002346200022e366004620040e4565b620009ff565b6040516200017597969594939291906200434c565b6200015262001162565b62000152620002643660046200439e565b62001265565b60005473ffffffffffffffffffffffffffffffffffffffff16620001c9565b620001526200029a3660046200442b565b62001ed6565b62000152620002b13660046200448e565b6200225e565b62000152620002c8366004620044bd565b620024f1565b62000195620002df36600462004593565b6200293c565b62000152620002f63660046200460a565b62002a0c565b620002346200030d366004620044bd565b62002a24565b6000805473ffffffffffffffffffffffffffffffffffffffff1633148015906200034757506200034560093362002a62565b155b156200037f576040517fd48b678b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff89163b620003ce576040517f09ee12d500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b620003d98662002a96565b9050600089307f00000000000000000000000000000000000000000000000000000000000000006040516200040e9062003d8a565b73ffffffffffffffffffffffffffffffffffffffff938416815291831660208301529091166040820152606001604051809103906000f08015801562000458573d6000803e3d6000fd5b5090506200051f826040518060e001604052806000151581526020018c63ffffffff16815260200163ffffffff801681526020018473ffffffffffffffffffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff168152602001600063ffffffff168152508a89898080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508b92508a915062002d429050565b6015805474010000000000000000000000000000000000000000900463ffffffff169060146200054f8362004659565b91906101000a81548163ffffffff021916908363ffffffff16021790555050817fbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d0128a8a604051620005c892919063ffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a2817fcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d878760405162000604929190620046c8565b60405180910390a2817f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664856040516200063e9190620046de565b60405180910390a2817f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf485084604051620006789190620046de565b60405180910390a25098975050505050505050565b600060606000806200069e6200312d565b600086815260046020908152604091829020825160e081018452815460ff811615158252610100810463ffffffff90811694830194909452650100000000008104841694820194909452690100000000000000000090930473ffffffffffffffffffffffffffffffffffffffff166060840152600101546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a0840152780100000000000000000000000000000000000000000000000090041660c08201525a9150600080826060015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620007b8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620007de919062004700565b73ffffffffffffffffffffffffffffffffffffffff166014600101600c9054906101000a900463ffffffff1663ffffffff168960405162000820919062004720565b60006040518083038160008787f1925050503d806000811462000860576040519150601f19603f3d011682016040523d82523d6000602084013e62000865565b606091505b50915091505a6200087790856200473e565b935081620008a257600060405180602001604052806000815250600796509650965050505062000928565b80806020019051810190620008b89190620047af565b909750955086620008e657600060405180602001604052806000815250600496509650965050505062000928565b601654865164010000000090910463ffffffff1610156200092457600060405180602001604052806000815250600596509650965050505062000928565b5050505b92959194509250565b6200093c8362003168565b6000838152601b6020526040902062000957828483620048a4565b50827f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d566483836040516200098c929190620046c8565b60405180910390a2505050565b6000620009f388888860008989604051806020016040528060008152508a8a8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506200031392505050565b98975050505050505050565b60006060600080600080600062000a156200312d565b600062000a228a6200321e565b905060006012604051806101600160405290816000820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160008201600c9054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160109054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160149054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160189054906101000a900462ffffff1662ffffff1662ffffff16815260200160008201601b9054906101000a900461ffff1661ffff1661ffff16815260200160008201601d9054906101000a900460ff1660ff1660ff16815260200160008201601e9054906101000a900460ff1615151515815260200160008201601f9054906101000a900460ff161515151581526020016001820160009054906101000a900460ff161515151581526020016001820160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152505090506000600460008d81526020019081526020016000206040518060e00160405290816000820160009054906101000a900460ff161515151581526020016000820160019054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160059054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160099054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160018201600c9054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff1681526020016001820160189054906101000a900463ffffffff1663ffffffff1663ffffffff168152505090508160e001511562000db7576000604051806020016040528060008152506009600084602001516000808263ffffffff169250995099509950995099509950995050505062001156565b604081015163ffffffff9081161462000e08576000604051806020016040528060008152506001600084602001516000808263ffffffff169250995099509950995099509950995050505062001156565b80511562000e4e576000604051806020016040528060008152506002600084602001516000808263ffffffff169250995099509950995099509950995050505062001156565b62000e5982620032d4565b602083015160165492975090955060009162000e8b918591879190640100000000900463ffffffff168a8a87620034c6565b9050806bffffffffffffffffffffffff168260a001516bffffffffffffffffffffffff16101562000ef5576000604051806020016040528060008152506006600085602001516000808263ffffffff1692509a509a509a509a509a509a509a505050505062001156565b600062000f048e868f62003517565b90505a9850600080846060015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000f5c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000f82919062004700565b73ffffffffffffffffffffffffffffffffffffffff166014600101600c9054906101000a900463ffffffff1663ffffffff168460405162000fc4919062004720565b60006040518083038160008787f1925050503d806000811462001004576040519150601f19603f3d011682016040523d82523d6000602084013e62001009565b606091505b50915091505a6200101b908c6200473e565b9a50816200109b5760165481516801000000000000000090910463ffffffff1610156200107857505060408051602080820190925260008082529490910151939c509a50600899505063ffffffff90911695506200115692505050565b602090940151939b5060039a505063ffffffff9092169650620011569350505050565b80806020019051810190620010b19190620047af565b909e509c508d620010f257505060408051602080820190925260008082529490910151939c509a50600499505063ffffffff90911695506200115692505050565b6016548d5164010000000090910463ffffffff1610156200114357505060408051602080820190925260008082529490910151939c509a50600599505063ffffffff90911695506200115692505050565b505050506020015163ffffffff16945050505b92959891949750929550565b60015473ffffffffffffffffffffffffffffffffffffffff163314620011e9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600173ffffffffffffffffffffffffffffffffffffffff82166000908152601a602052604090205460ff166003811115620012a457620012a4620041a1565b14158015620012f05750600373ffffffffffffffffffffffffffffffffffffffff82166000908152601a602052604090205460ff166003811115620012ed57620012ed620041a1565b14155b1562001328576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6014546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1662001388576040517fd12d7d8d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000829003620013c4576040517f2c2fc94100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160e081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c081018290526000808567ffffffffffffffff8111156200141b576200141b62003e85565b60405190808252806020026020018201604052801562001445578160200160208202803683370190505b50905060008667ffffffffffffffff81111562001466576200146662003e85565b604051908082528060200260200182016040528015620014ed57816020015b6040805160e08101825260008082526020808301829052928201819052606082018190526080820181905260a0820181905260c082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181620014855790505b50905060008767ffffffffffffffff8111156200150e576200150e62003e85565b6040519080825280602002602001820160405280156200154357816020015b60608152602001906001900390816200152d5790505b50905060008867ffffffffffffffff81111562001564576200156462003e85565b6040519080825280602002602001820160405280156200159957816020015b6060815260200190600190039081620015835790505b50905060008967ffffffffffffffff811115620015ba57620015ba62003e85565b604051908082528060200260200182016040528015620015ef57816020015b6060815260200190600190039081620015d95790505b50905060005b8a81101562001bd3578b8b82818110620016135762001613620049cc565b60209081029290920135600081815260048452604090819020815160e081018352815460ff811615158252610100810463ffffffff90811697830197909752650100000000008104871693820193909352690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff166060830152600101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490931660c08401529a50909850620016f290508962003168565b60608801516040517f1a5da6c800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8c8116600483015290911690631a5da6c890602401600060405180830381600087803b1580156200176257600080fd5b505af115801562001777573d6000803e3d6000fd5b5050505087858281518110620017915762001791620049cc565b6020026020010181905250600560008a815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16868281518110620017e557620017e5620049cc565b73ffffffffffffffffffffffffffffffffffffffff90921660209283029190910182015260008a815260079091526040902080546200182490620047fc565b80601f01602080910402602001604051908101604052809291908181526020018280546200185290620047fc565b8015620018a35780601f106200187757610100808354040283529160200191620018a3565b820191906000526020600020905b8154815290600101906020018083116200188557829003601f168201915b5050505050848281518110620018bd57620018bd620049cc565b6020026020010181905250601b60008a81526020019081526020016000208054620018e890620047fc565b80601f01602080910402602001604051908101604052809291908181526020018280546200191690620047fc565b8015620019675780601f106200193b5761010080835404028352916020019162001967565b820191906000526020600020905b8154815290600101906020018083116200194957829003601f168201915b5050505050838281518110620019815762001981620049cc565b6020026020010181905250601c60008a81526020019081526020016000208054620019ac90620047fc565b80601f0160208091040260200160405190810160405280929190818152602001828054620019da90620047fc565b801562001a2b5780601f10620019ff5761010080835404028352916020019162001a2b565b820191906000526020600020905b81548152906001019060200180831162001a0d57829003601f168201915b505050505082828151811062001a455762001a45620049cc565b60200260200101819052508760a001516bffffffffffffffffffffffff168762001a709190620049fb565b60008a815260046020908152604080832080547fffffff000000000000000000000000000000000000000000000000000000000016815560010180547fffffffff000000000000000000000000000000000000000000000000000000001690556007909152812091985062001ae6919062003d98565b6000898152601b6020526040812062001aff9162003d98565b6000898152601c6020526040812062001b189162003d98565b600089815260066020526040902080547fffffffffffffffffffffffff000000000000000000000000000000000000000016905562001b5960028a62003739565b5060a0880151604080516bffffffffffffffffffffffff909216825273ffffffffffffffffffffffffffffffffffffffff8c1660208301528a917fb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff910160405180910390a28062001bca8162004a11565b915050620015f5565b508560195462001be491906200473e565b60195560008b8b868167ffffffffffffffff81111562001c085762001c0862003e85565b60405190808252806020026020018201604052801562001c32578160200160208202803683370190505b508988888860405160200162001c5098979695949392919062004bb7565b60405160208183030381529060405290508973ffffffffffffffffffffffffffffffffffffffff16638e86139b6014600001600c9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663c71249ab60038e73ffffffffffffffffffffffffffffffffffffffff1663aab9edd66040518163ffffffff1660e01b8152600401602060405180830381865afa15801562001d0c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001d32919062004c87565b866040518463ffffffff1660e01b815260040162001d539392919062004cac565b600060405180830381865afa15801562001d71573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405262001db9919081019062004cd3565b6040518263ffffffff1660e01b815260040162001dd79190620046de565b600060405180830381600087803b15801562001df257600080fd5b505af115801562001e07573d6000803e3d6000fd5b50506040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8d81166004830152602482018b90527f000000000000000000000000000000000000000000000000000000000000000016925063a9059cbb91506044016020604051808303816000875af115801562001ea1573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001ec7919062004d0c565b50505050505050505050505050565b6002336000908152601a602052604090205460ff16600381111562001eff5762001eff620041a1565b1415801562001f3557506003336000908152601a602052604090205460ff16600381111562001f325762001f32620041a1565b14155b1562001f6d576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080808080808062001f83888a018a62004f04565b965096509650965096509650965060005b87518110156200225257600073ffffffffffffffffffffffffffffffffffffffff1687828151811062001fcb5762001fcb620049cc565b60200260200101516060015173ffffffffffffffffffffffffffffffffffffffff1603620020df57858181518110620020085762002008620049cc565b6020026020010151307f0000000000000000000000000000000000000000000000000000000000000000604051620020409062003d8a565b73ffffffffffffffffffffffffffffffffffffffff938416815291831660208301529091166040820152606001604051809103906000f0801580156200208a573d6000803e3d6000fd5b50878281518110620020a057620020a0620049cc565b60200260200101516060019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250505b62002197888281518110620020f857620020f8620049cc565b6020026020010151888381518110620021155762002115620049cc565b6020026020010151878481518110620021325762002132620049cc565b60200260200101518785815181106200214f576200214f620049cc565b60200260200101518786815181106200216c576200216c620049cc565b6020026020010151878781518110620021895762002189620049cc565b602002602001015162002d42565b878181518110620021ac57620021ac620049cc565b60200260200101517f74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71888381518110620021ea57620021ea620049cc565b602002602001015160a0015133604051620022359291906bffffffffffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a280620022498162004a11565b91505062001f94565b50505050505050505050565b600082815260046020908152604091829020825160e081018452815460ff81161515825263ffffffff6101008204811694830194909452650100000000008104841694820185905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004821660c082015291146200235c576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818160a001516200236e919062005035565b600084815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff93841602179055601954620023d691841690620049fb565b6019556040517f23b872dd0000000000000000000000000000000000000000000000000000000081523360048201523060248201526bffffffffffffffffffffffff831660448201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906323b872dd906064016020604051808303816000875af115801562002480573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620024a6919062004d0c565b506040516bffffffffffffffffffffffff83168152339084907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a3505050565b6000818152600460209081526040808320815160e081018352815460ff81161515825263ffffffff6101008204811695830195909552650100000000008104851693820184905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004831660c08201529291141590620025da60005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161490506000601260010160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200267d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620026a391906200505d565b9050828015620026c75750818015620026c5575080846040015163ffffffff16115b155b15620026ff576040517ffbc0357800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8115801562002732575060008581526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314155b156200276a576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8162002780576200277d603282620049fb565b90505b6000858152600460205260409020805463ffffffff80841665010000000000027fffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffffff90921691909117909155620027dc9060029087906200373916565b5060145460808501516bffffffffffffffffffffffff91821691600091168211156200284557608086015162002813908362005077565b90508560a001516bffffffffffffffffffffffff16816bffffffffffffffffffffffff16111562002845575060a08501515b808660a0015162002857919062005077565b600088815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff93841602179055601554620028bf9183911662005035565b601580547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff9290921691909117905560405167ffffffffffffffff84169088907f91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f79118190600090a350505050505050565b600060606000806200294d6200312d565b6000634b56a42e60e01b8888886040516024016200296e939291906200509f565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091529050620029f989826200068d565b929c919b50995090975095505050505050565b62002a1662003747565b62002a2181620037ca565b50565b60006060600080600080600062002a4b8860405180602001604052806000815250620009ff565b959e949d50929b5090995097509550909350915050565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415155b90505b92915050565b6000806000601260010160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905060008173ffffffffffffffffffffffffffffffffffffffff166385df51fd60018473ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa15801562002b2f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062002b5591906200505d565b62002b6191906200473e565b6040518263ffffffff1660e01b815260040162002b8091815260200190565b602060405180830381865afa15801562002b9e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062002bc491906200505d565b601554604080516020810193909352309083015274010000000000000000000000000000000000000000900463ffffffff166060820152608001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201209083015201604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052905060045b600f81101562002cd0578382828151811062002c8c5762002c8c620049cc565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508062002cc78162004a11565b91505062002c6c565b5084600181111562002ce65762002ce6620041a1565b60f81b81600f8151811062002cff5762002cff620049cc565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535062002d3981620050d3565b95945050505050565b6012547e01000000000000000000000000000000000000000000000000000000000000900460ff161562002da2576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601654835163ffffffff909116101562002de8576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108fc856020015163ffffffff16108062002e265750601554602086015163ffffffff70010000000000000000000000000000000090920482169116115b1562002e5e576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000868152600460205260409020546901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff161562002ec8576040517f6e3b930b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000868152600460209081526040808320885181548a8501518b85015160608d01517fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000009093169315157fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000ff169390931761010063ffffffff92831602177fffffff000000000000000000000000000000000000000000000000ffffffffff1665010000000000938216939093027fffffff0000000000000000000000000000000000000000ffffffffffffffffff1692909217690100000000000000000073ffffffffffffffffffffffffffffffffffffffff9283160217835560808b01516001909301805460a08d015160c08e01516bffffffffffffffffffffffff9687167fffffffffffffffff000000000000000000000000000000000000000000000000909316929092176c010000000000000000000000009690911695909502949094177fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff1678010000000000000000000000000000000000000000000000009490931693909302919091179091556005835281842080547fffffffffffffffffffffffff00000000000000000000000000000000000000001691891691909117905560079091529020620030bb848262005116565b508460a001516bffffffffffffffffffffffff16601954620030de9190620049fb565b6019556000868152601b60205260409020620030fb838262005116565b506000868152601c6020526040902062003116828262005116565b5062003124600287620038c1565b50505050505050565b321562003166576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60008181526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314620031c6576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526004602052604090205465010000000000900463ffffffff9081161462002a21576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000818160045b600f811015620032b3577fff000000000000000000000000000000000000000000000000000000000000008216838260208110620032675762003267620049cc565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916146200329e57506000949350505050565b80620032aa8162004a11565b91505062003225565b5081600f1a6001811115620032cc57620032cc620041a1565b949350505050565b6000806000836080015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa15801562003361573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003387919062005258565b50945090925050506000811315806200339f57508142105b80620033c45750828015620033c45750620033bb82426200473e565b8463ffffffff16105b15620033d5576017549550620033d9565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa15801562003445573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200346b919062005258565b50945090925050506000811315806200348357508142105b80620034a85750828015620034a857506200349f82426200473e565b8463ffffffff16105b15620034b9576018549450620034bd565b8094505b50505050915091565b600080620034da88878b60c00151620038cf565b9050600080620034f78b8a63ffffffff16858a8a60018b62003965565b909250905062003508818362005035565b9b9a5050505050505050505050565b60606000836001811115620035305762003530620041a1565b03620035fd576000848152600760205260409081902090517f6e04ff0d0000000000000000000000000000000000000000000000000000000091620035789160240162005350565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152905062003732565b6001836001811115620036145762003614620041a1565b036200370057600082806020019051810190620036329190620053c7565b6000868152600760205260409081902090519192507f40691db4000000000000000000000000000000000000000000000000000000009162003679918491602401620054db565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091529150620037329050565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b9392505050565b600062002a8d838362003c2d565b60005473ffffffffffffffffffffffffffffffffffffffff16331462003166576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401620011e0565b3373ffffffffffffffffffffffffffffffffffffffff8216036200384b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620011e0565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600062002a8d838362003d38565b60008080856001811115620038e857620038e8620041a1565b03620038f9575062015f906200391c565b6001856001811115620039105762003910620041a1565b036200370057506201adb05b6200392f63ffffffff85166014620055a3565b6200393c846001620055e3565b6200394d9060ff16611d4c620055a3565b620039599083620049fb565b62002d399190620049fb565b60008060008960a0015161ffff1687620039809190620055a3565b90508380156200398f5750803a105b156200399857503a5b6000841562003a30578a610140015173ffffffffffffffffffffffffffffffffffffffff166349948e0e6000366040518363ffffffff1660e01b8152600401620039e4929190620046c8565b602060405180830381865afa15801562003a02573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003a2891906200505d565b905062003af1565b6101408b01516016546040517f1254414000000000000000000000000000000000000000000000000000000000815264010000000090910463ffffffff16600482015273ffffffffffffffffffffffffffffffffffffffff90911690631254414090602401602060405180830381865afa15801562003ab3573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003ad991906200505d565b8b60a0015161ffff1662003aee9190620055a3565b90505b62003b0161ffff871682620055ff565b90506000878262003b138c8e620049fb565b62003b1f9086620055a3565b62003b2b9190620049fb565b62003b3f90670de0b6b3a7640000620055a3565b62003b4b9190620055ff565b905060008c6040015163ffffffff1664e8d4a5100062003b6c9190620055a3565b898e6020015163ffffffff16858f8862003b879190620055a3565b62003b939190620049fb565b62003ba390633b9aca00620055a3565b62003baf9190620055a3565b62003bbb9190620055ff565b62003bc79190620049fb565b90506b033b2e3c9fd0803ce800000062003be28284620049fb565b111562003c1b576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b6000818152600183016020526040812054801562003d2657600062003c546001836200473e565b855490915060009062003c6a906001906200473e565b905081811462003cd657600086600001828154811062003c8e5762003c8e620049cc565b906000526020600020015490508087600001848154811062003cb45762003cb4620049cc565b6000918252602080832090910192909255918252600188019052604090208390555b855486908062003cea5762003cea6200563b565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505062002a90565b600091505062002a90565b5092915050565b600081815260018301602052604081205462003d815750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915562002a90565b50600062002a90565b6103ca806200566b83390190565b50805462003da690620047fc565b6000825580601f1062003db7575050565b601f01602090049060005260206000209081019062002a2191905b8082111562003de8576000815560010162003dd2565b5090565b73ffffffffffffffffffffffffffffffffffffffff8116811462002a2157600080fd5b803563ffffffff8116811462003e2457600080fd5b919050565b80356002811062003e2457600080fd5b60008083601f84011262003e4c57600080fd5b50813567ffffffffffffffff81111562003e6557600080fd5b60208301915083602082850101111562003e7e57600080fd5b9250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160e0810167ffffffffffffffff8111828210171562003eda5762003eda62003e85565b60405290565b604051610100810167ffffffffffffffff8111828210171562003eda5762003eda62003e85565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171562003f515762003f5162003e85565b604052919050565b600067ffffffffffffffff82111562003f765762003f7662003e85565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f83011262003fb457600080fd5b813562003fcb62003fc58262003f59565b62003f07565b81815284602083860101111562003fe157600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060008060e0898b0312156200401b57600080fd5b8835620040288162003dec565b97506200403860208a0162003e0f565b965060408901356200404a8162003dec565b95506200405a60608a0162003e29565b9450608089013567ffffffffffffffff808211156200407857600080fd5b620040868c838d0162003e39565b909650945060a08b0135915080821115620040a057600080fd5b620040ae8c838d0162003fa2565b935060c08b0135915080821115620040c557600080fd5b50620040d48b828c0162003fa2565b9150509295985092959890939650565b60008060408385031215620040f857600080fd5b82359150602083013567ffffffffffffffff8111156200411757600080fd5b620041258582860162003fa2565b9150509250929050565b60005b838110156200414c57818101518382015260200162004132565b50506000910152565b600081518084526200416f8160208601602086016200412f565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600a811062004208577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b9052565b841515815260806020820152600062004229608083018662004155565b90506200423a6040830185620041d0565b82606083015295945050505050565b6000806000604084860312156200425f57600080fd5b83359250602084013567ffffffffffffffff8111156200427e57600080fd5b6200428c8682870162003e39565b9497909650939450505050565b600080600080600080600060a0888a031215620042b557600080fd5b8735620042c28162003dec565b9650620042d26020890162003e0f565b95506040880135620042e48162003dec565b9450606088013567ffffffffffffffff808211156200430257600080fd5b620043108b838c0162003e39565b909650945060808a01359150808211156200432a57600080fd5b50620043398a828b0162003e39565b989b979a50959850939692959293505050565b871515815260e0602082015260006200436960e083018962004155565b90506200437a6040830188620041d0565b8560608301528460808301528360a08301528260c083015298975050505050505050565b600080600060408486031215620043b457600080fd5b833567ffffffffffffffff80821115620043cd57600080fd5b818601915086601f830112620043e257600080fd5b813581811115620043f257600080fd5b8760208260051b85010111156200440857600080fd5b60209283019550935050840135620044208162003dec565b809150509250925092565b600080602083850312156200443f57600080fd5b823567ffffffffffffffff8111156200445757600080fd5b620044658582860162003e39565b90969095509350505050565b80356bffffffffffffffffffffffff8116811462003e2457600080fd5b60008060408385031215620044a257600080fd5b82359150620044b46020840162004471565b90509250929050565b600060208284031215620044d057600080fd5b5035919050565b600067ffffffffffffffff821115620044f457620044f462003e85565b5060051b60200190565b600082601f8301126200451057600080fd5b813560206200452362003fc583620044d7565b82815260059290921b840181019181810190868411156200454357600080fd5b8286015b848110156200458857803567ffffffffffffffff811115620045695760008081fd5b620045798986838b010162003fa2565b84525091830191830162004547565b509695505050505050565b60008060008060608587031215620045aa57600080fd5b84359350602085013567ffffffffffffffff80821115620045ca57600080fd5b620045d888838901620044fe565b94506040870135915080821115620045ef57600080fd5b50620045fe8782880162003e39565b95989497509550505050565b6000602082840312156200461d57600080fd5b8135620037328162003dec565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600063ffffffff8083168181036200467557620046756200462a565b6001019392505050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b602081526000620032cc6020830184866200467f565b60208152600062002a8d602083018462004155565b805162003e248162003dec565b6000602082840312156200471357600080fd5b8151620037328162003dec565b60008251620047348184602087016200412f565b9190910192915050565b8181038181111562002a905762002a906200462a565b801515811462002a2157600080fd5b600082601f8301126200477557600080fd5b81516200478662003fc58262003f59565b8181528460208386010111156200479c57600080fd5b620032cc8260208301602087016200412f565b60008060408385031215620047c357600080fd5b8251620047d08162004754565b602084015190925067ffffffffffffffff811115620047ee57600080fd5b620041258582860162004763565b600181811c908216806200481157607f821691505b6020821081036200484b577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f8211156200489f57600081815260208120601f850160051c810160208610156200487a5750805b601f850160051c820191505b818110156200489b5782815560010162004886565b5050505b505050565b67ffffffffffffffff831115620048bf57620048bf62003e85565b620048d783620048d08354620047fc565b8362004851565b6000601f8411600181146200492c5760008515620048f55750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355620049c5565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b828110156200497d57868501358255602094850194600190920191016200495b565b5086821015620049b9577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b8082018082111562002a905762002a906200462a565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820362004a455762004a456200462a565b5060010190565b600081518084526020808501945080840160005b8381101562004b0b5781518051151588528381015163ffffffff908116858a01526040808301519091169089015260608082015173ffffffffffffffffffffffffffffffffffffffff16908901526080808201516bffffffffffffffffffffffff169089015260a08082015162004ae6828b01826bffffffffffffffffffffffff169052565b505060c09081015163ffffffff169088015260e0909601959082019060010162004a60565b509495945050505050565b600081518084526020808501945080840160005b8381101562004b0b57815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010162004b2a565b600081518084526020808501808196508360051b8101915082860160005b8581101562004baa57828403895262004b9784835162004155565b9885019893509084019060010162004b7c565b5091979650505050505050565b60e081528760e082015260006101007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8a111562004bf457600080fd5b8960051b808c8386013783018381038201602085015262004c188282018b62004a4c565b915050828103604084015262004c2f818962004b16565b9050828103606084015262004c45818862004b16565b9050828103608084015262004c5b818762004b5e565b905082810360a084015262004c71818662004b5e565b905082810360c084015262003508818562004b5e565b60006020828403121562004c9a57600080fd5b815160ff811681146200373257600080fd5b60ff8416815260ff8316602082015260606040820152600062002d39606083018462004155565b60006020828403121562004ce657600080fd5b815167ffffffffffffffff81111562004cfe57600080fd5b620032cc8482850162004763565b60006020828403121562004d1f57600080fd5b8151620037328162004754565b600082601f83011262004d3e57600080fd5b8135602062004d5162003fc583620044d7565b82815260059290921b8401810191818101908684111562004d7157600080fd5b8286015b8481101562004588578035835291830191830162004d75565b600082601f83011262004da057600080fd5b8135602062004db362003fc583620044d7565b82815260e0928302850182019282820191908785111562004dd357600080fd5b8387015b8581101562004e8a5781818a03121562004df15760008081fd5b62004dfb62003eb4565b813562004e088162004754565b815262004e1782870162003e0f565b86820152604062004e2a81840162003e0f565b9082015260608281013562004e3f8162003dec565b90820152608062004e5283820162004471565b9082015260a062004e6583820162004471565b9082015260c062004e7883820162003e0f565b90820152845292840192810162004dd7565b5090979650505050505050565b600082601f83011262004ea957600080fd5b8135602062004ebc62003fc583620044d7565b82815260059290921b8401810191818101908684111562004edc57600080fd5b8286015b848110156200458857803562004ef68162003dec565b835291830191830162004ee0565b600080600080600080600060e0888a03121562004f2057600080fd5b873567ffffffffffffffff8082111562004f3957600080fd5b62004f478b838c0162004d2c565b985060208a013591508082111562004f5e57600080fd5b62004f6c8b838c0162004d8e565b975060408a013591508082111562004f8357600080fd5b62004f918b838c0162004e97565b965060608a013591508082111562004fa857600080fd5b62004fb68b838c0162004e97565b955060808a013591508082111562004fcd57600080fd5b62004fdb8b838c01620044fe565b945060a08a013591508082111562004ff257600080fd5b620050008b838c01620044fe565b935060c08a01359150808211156200501757600080fd5b50620050268a828b01620044fe565b91505092959891949750929550565b6bffffffffffffffffffffffff81811683821601908082111562003d315762003d316200462a565b6000602082840312156200507057600080fd5b5051919050565b6bffffffffffffffffffffffff82811682821603908082111562003d315762003d316200462a565b604081526000620050b4604083018662004b5e565b8281036020840152620050c98185876200467f565b9695505050505050565b805160208083015191908110156200484b577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60209190910360031b1b16919050565b815167ffffffffffffffff81111562005133576200513362003e85565b6200514b81620051448454620047fc565b8462004851565b602080601f831160018114620051a157600084156200516a5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556200489b565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015620051f057888601518255948401946001909101908401620051cf565b50858210156200522d57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b805169ffffffffffffffffffff8116811462003e2457600080fd5b600080600080600060a086880312156200527157600080fd5b6200527c866200523d565b9450602086015193506040860151925060608601519150620052a1608087016200523d565b90509295509295909350565b60008154620052bc81620047fc565b808552602060018381168015620052dc5760018114620053155762005345565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b890101955062005345565b866000528260002060005b858110156200533d5781548a820186015290830190840162005320565b890184019650505b505050505092915050565b60208152600062002a8d6020830184620052ad565b600082601f8301126200537757600080fd5b815160206200538a62003fc583620044d7565b82815260059290921b84018101918181019086841115620053aa57600080fd5b8286015b84811015620045885780518352918301918301620053ae565b600060208284031215620053da57600080fd5b815167ffffffffffffffff80821115620053f357600080fd5b9083019061010082860312156200540957600080fd5b6200541362003ee0565b82518152602083015160208201526040830151604082015260608301516060820152608083015160808201526200544d60a08401620046f3565b60a082015260c0830151828111156200546557600080fd5b620054738782860162005365565b60c08301525060e0830151828111156200548c57600080fd5b6200549a8782860162004763565b60e08301525095945050505050565b600081518084526020808501945080840160005b8381101562004b0b57815187529582019590820190600101620054bd565b60408152825160408201526020830151606082015260408301516080820152606083015160a0820152608083015160c082015273ffffffffffffffffffffffffffffffffffffffff60a08401511660e0820152600060c08401516101008081850152506200554e610140840182620054a9565b905060e08501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0848303016101208501526200558c828262004155565b915050828103602084015262002d398185620052ad565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615620055de57620055de6200462a565b500290565b60ff818116838216019081111562002a905762002a906200462a565b60008262005636577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfe60c060405234801561001057600080fd5b506040516103ca3803806103ca83398101604081905261002f91610076565b600080546001600160a01b0319166001600160a01b039384161790559181166080521660a0526100b9565b80516001600160a01b038116811461007157600080fd5b919050565b60008060006060848603121561008b57600080fd5b6100948461005a565b92506100a26020850161005a565b91506100b06040850161005a565b90509250925092565b60805160a0516102e76100e36000396000603801526000818160c4015261011701526102e76000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806379188d161461007b578063f00e6a2a146100aa575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e808015610076573d6000f35b3d6000fd5b61008e6100893660046101c1565b6100ee565b6040805192151583526020830191909152015b60405180910390f35b60405173ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001681526020016100a1565b60008054819073ffffffffffffffffffffffffffffffffffffffff16331461011557600080fd5b7f00000000000000000000000000000000000000000000000000000000000000005a91505a61138881101561014957600080fd5b61138881039050856040820482031161016157600080fd5b50803b61016d57600080fd5b6000808551602087016000858af192505a610188908361029a565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080604083850312156101d457600080fd5b82359150602083013567ffffffffffffffff808211156101f357600080fd5b818501915085601f83011261020757600080fd5b81358181111561021957610219610192565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561025f5761025f610192565b8160405282815288602084870101111561027857600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b818103818111156102d4577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b9291505056fea164736f6c6343000810000aa164736f6c6343000810000a",
>>>>>>> ab4212059a (implement modules)
=======
	Bin: "0x6101206040523480156200001257600080fd5b5060405162005de138038062005de1833981016040819052620000359162000345565b80816001600160a01b031663ca30e6036040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200009b919062000345565b826001600160a01b031663b10b673c6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000da573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000100919062000345565b836001600160a01b0316636709d0e56040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200013f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000165919062000345565b846001600160a01b0316635425d8ac6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca919062000345565b3380600081620002215760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200025457620002548162000281565b5050506001600160a01b0393841660805291831660a052821660c052811660e0521661010052506200036c565b336001600160a01b03821603620002db5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000218565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200034257600080fd5b50565b6000602082840312156200035857600080fd5b815162000365816200032c565b9392505050565b60805160a05160c05160e05161010051615a1b620003c66000396000818161010e01526101a90152600081816103e10152612013015260006132f7015260006133db015260008181611e5501526124210152615a1b6000f3fe60806040523480156200001157600080fd5b50600436106200010c5760003560e01c806385c1b0ba11620000a5578063c8048022116200006f578063c804802214620002b7578063ce7dc5b414620002ce578063f2fde38b14620002e5578063f7d334ba14620002fc576200010c565b806385c1b0ba14620002535780638da5cb5b146200026a5780638e86139b1462000289578063948108f714620002a0576200010c565b80634ee88d3511620000e75780634ee88d3514620001ef5780636ded9eae146200020657806371791aa0146200021d57806379ba50971462000249576200010c565b806328f32f38146200015457806329c5efad146200017e578063349e8cca14620001a7575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e8080156200014d573d6000f35b3d6000fd5b005b6200016b6200016536600462003ffe565b62000313565b6040519081526020015b60405180910390f35b620001956200018f366004620040e4565b6200068d565b6040516200017594939291906200420c565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200162000175565b620001526200020036600462004249565b62000931565b6200016b6200021736600462004299565b62000999565b620002346200022e366004620040e4565b620009ff565b6040516200017597969594939291906200434c565b6200015262001162565b62000152620002643660046200439e565b62001265565b60005473ffffffffffffffffffffffffffffffffffffffff16620001c9565b620001526200029a3660046200442b565b62001ed6565b62000152620002b13660046200448e565b6200225e565b62000152620002c8366004620044bd565b620024f1565b62000195620002df36600462004593565b6200293c565b62000152620002f63660046200460a565b62002a0c565b620002346200030d366004620044bd565b62002a24565b6000805473ffffffffffffffffffffffffffffffffffffffff1633148015906200034757506200034560093362002a62565b155b156200037f576040517fd48b678b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff89163b620003ce576040517f09ee12d500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b620003d98662002a96565b9050600089307f00000000000000000000000000000000000000000000000000000000000000006040516200040e9062003d8a565b73ffffffffffffffffffffffffffffffffffffffff938416815291831660208301529091166040820152606001604051809103906000f08015801562000458573d6000803e3d6000fd5b5090506200051f826040518060e001604052806000151581526020018c63ffffffff16815260200163ffffffff801681526020018473ffffffffffffffffffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff168152602001600063ffffffff168152508a89898080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508b92508a915062002d429050565b6015805474010000000000000000000000000000000000000000900463ffffffff169060146200054f8362004659565b91906101000a81548163ffffffff021916908363ffffffff16021790555050817fbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d0128a8a604051620005c892919063ffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a2817fcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d878760405162000604929190620046c8565b60405180910390a2817f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664856040516200063e9190620046de565b60405180910390a2817f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf485084604051620006789190620046de565b60405180910390a25098975050505050505050565b600060606000806200069e6200312d565b600086815260046020908152604091829020825160e081018452815460ff811615158252610100810463ffffffff90811694830194909452650100000000008104841694820194909452690100000000000000000090930473ffffffffffffffffffffffffffffffffffffffff166060840152600101546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a0840152780100000000000000000000000000000000000000000000000090041660c08201525a9150600080826060015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620007b8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620007de919062004700565b73ffffffffffffffffffffffffffffffffffffffff166014600101600c9054906101000a900463ffffffff1663ffffffff168960405162000820919062004720565b60006040518083038160008787f1925050503d806000811462000860576040519150601f19603f3d011682016040523d82523d6000602084013e62000865565b606091505b50915091505a6200087790856200473e565b935081620008a257600060405180602001604052806000815250600796509650965050505062000928565b80806020019051810190620008b89190620047af565b909750955086620008e657600060405180602001604052806000815250600496509650965050505062000928565b601654865164010000000090910463ffffffff1610156200092457600060405180602001604052806000815250600596509650965050505062000928565b5050505b92959194509250565b6200093c8362003168565b6000838152601b6020526040902062000957828483620048a4565b50827f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d566483836040516200098c929190620046c8565b60405180910390a2505050565b6000620009f388888860008989604051806020016040528060008152508a8a8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506200031392505050565b98975050505050505050565b60006060600080600080600062000a156200312d565b600062000a228a6200321e565b905060006012604051806101600160405290816000820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160008201600c9054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160109054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160149054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160189054906101000a900462ffffff1662ffffff1662ffffff16815260200160008201601b9054906101000a900461ffff1661ffff1661ffff16815260200160008201601d9054906101000a900460ff1660ff1660ff16815260200160008201601e9054906101000a900460ff1615151515815260200160008201601f9054906101000a900460ff161515151581526020016001820160009054906101000a900460ff161515151581526020016001820160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152505090506000600460008d81526020019081526020016000206040518060e00160405290816000820160009054906101000a900460ff161515151581526020016000820160019054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160059054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160099054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160018201600c9054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff1681526020016001820160189054906101000a900463ffffffff1663ffffffff1663ffffffff168152505090508160e001511562000db7576000604051806020016040528060008152506009600084602001516000808263ffffffff169250995099509950995099509950995050505062001156565b604081015163ffffffff9081161462000e08576000604051806020016040528060008152506001600084602001516000808263ffffffff169250995099509950995099509950995050505062001156565b80511562000e4e576000604051806020016040528060008152506002600084602001516000808263ffffffff169250995099509950995099509950995050505062001156565b62000e5982620032d4565b602083015160165492975090955060009162000e8b918591879190640100000000900463ffffffff168a8a87620034c6565b9050806bffffffffffffffffffffffff168260a001516bffffffffffffffffffffffff16101562000ef5576000604051806020016040528060008152506006600085602001516000808263ffffffff1692509a509a509a509a509a509a509a505050505062001156565b600062000f048e868f62003517565b90505a9850600080846060015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000f5c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000f82919062004700565b73ffffffffffffffffffffffffffffffffffffffff166014600101600c9054906101000a900463ffffffff1663ffffffff168460405162000fc4919062004720565b60006040518083038160008787f1925050503d806000811462001004576040519150601f19603f3d011682016040523d82523d6000602084013e62001009565b606091505b50915091505a6200101b908c6200473e565b9a50816200109b5760165481516801000000000000000090910463ffffffff1610156200107857505060408051602080820190925260008082529490910151939c509a50600899505063ffffffff90911695506200115692505050565b602090940151939b5060039a505063ffffffff9092169650620011569350505050565b80806020019051810190620010b19190620047af565b909e509c508d620010f257505060408051602080820190925260008082529490910151939c509a50600499505063ffffffff90911695506200115692505050565b6016548d5164010000000090910463ffffffff1610156200114357505060408051602080820190925260008082529490910151939c509a50600599505063ffffffff90911695506200115692505050565b505050506020015163ffffffff16945050505b92959891949750929550565b60015473ffffffffffffffffffffffffffffffffffffffff163314620011e9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600173ffffffffffffffffffffffffffffffffffffffff82166000908152601a602052604090205460ff166003811115620012a457620012a4620041a1565b14158015620012f05750600373ffffffffffffffffffffffffffffffffffffffff82166000908152601a602052604090205460ff166003811115620012ed57620012ed620041a1565b14155b1562001328576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6014546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1662001388576040517fd12d7d8d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000829003620013c4576040517f2c2fc94100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160e081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c081018290526000808567ffffffffffffffff8111156200141b576200141b62003e85565b60405190808252806020026020018201604052801562001445578160200160208202803683370190505b50905060008667ffffffffffffffff81111562001466576200146662003e85565b604051908082528060200260200182016040528015620014ed57816020015b6040805160e08101825260008082526020808301829052928201819052606082018190526080820181905260a0820181905260c082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181620014855790505b50905060008767ffffffffffffffff8111156200150e576200150e62003e85565b6040519080825280602002602001820160405280156200154357816020015b60608152602001906001900390816200152d5790505b50905060008867ffffffffffffffff81111562001564576200156462003e85565b6040519080825280602002602001820160405280156200159957816020015b6060815260200190600190039081620015835790505b50905060008967ffffffffffffffff811115620015ba57620015ba62003e85565b604051908082528060200260200182016040528015620015ef57816020015b6060815260200190600190039081620015d95790505b50905060005b8a81101562001bd3578b8b82818110620016135762001613620049cc565b60209081029290920135600081815260048452604090819020815160e081018352815460ff811615158252610100810463ffffffff90811697830197909752650100000000008104871693820193909352690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff166060830152600101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490931660c08401529a50909850620016f290508962003168565b60608801516040517f1a5da6c800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8c8116600483015290911690631a5da6c890602401600060405180830381600087803b1580156200176257600080fd5b505af115801562001777573d6000803e3d6000fd5b5050505087858281518110620017915762001791620049cc565b6020026020010181905250600560008a815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16868281518110620017e557620017e5620049cc565b73ffffffffffffffffffffffffffffffffffffffff90921660209283029190910182015260008a815260079091526040902080546200182490620047fc565b80601f01602080910402602001604051908101604052809291908181526020018280546200185290620047fc565b8015620018a35780601f106200187757610100808354040283529160200191620018a3565b820191906000526020600020905b8154815290600101906020018083116200188557829003601f168201915b5050505050848281518110620018bd57620018bd620049cc565b6020026020010181905250601b60008a81526020019081526020016000208054620018e890620047fc565b80601f01602080910402602001604051908101604052809291908181526020018280546200191690620047fc565b8015620019675780601f106200193b5761010080835404028352916020019162001967565b820191906000526020600020905b8154815290600101906020018083116200194957829003601f168201915b5050505050838281518110620019815762001981620049cc565b6020026020010181905250601c60008a81526020019081526020016000208054620019ac90620047fc565b80601f0160208091040260200160405190810160405280929190818152602001828054620019da90620047fc565b801562001a2b5780601f10620019ff5761010080835404028352916020019162001a2b565b820191906000526020600020905b81548152906001019060200180831162001a0d57829003601f168201915b505050505082828151811062001a455762001a45620049cc565b60200260200101819052508760a001516bffffffffffffffffffffffff168762001a709190620049fb565b60008a815260046020908152604080832080547fffffff000000000000000000000000000000000000000000000000000000000016815560010180547fffffffff000000000000000000000000000000000000000000000000000000001690556007909152812091985062001ae6919062003d98565b6000898152601b6020526040812062001aff9162003d98565b6000898152601c6020526040812062001b189162003d98565b600089815260066020526040902080547fffffffffffffffffffffffff000000000000000000000000000000000000000016905562001b5960028a62003739565b5060a0880151604080516bffffffffffffffffffffffff909216825273ffffffffffffffffffffffffffffffffffffffff8c1660208301528a917fb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff910160405180910390a28062001bca8162004a11565b915050620015f5565b508560195462001be491906200473e565b60195560008b8b868167ffffffffffffffff81111562001c085762001c0862003e85565b60405190808252806020026020018201604052801562001c32578160200160208202803683370190505b508988888860405160200162001c5098979695949392919062004bb7565b60405160208183030381529060405290508973ffffffffffffffffffffffffffffffffffffffff16638e86139b6014600001600c9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663c71249ab60038e73ffffffffffffffffffffffffffffffffffffffff1663aab9edd66040518163ffffffff1660e01b8152600401602060405180830381865afa15801562001d0c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001d32919062004c87565b866040518463ffffffff1660e01b815260040162001d539392919062004cac565b600060405180830381865afa15801562001d71573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405262001db9919081019062004cd3565b6040518263ffffffff1660e01b815260040162001dd79190620046de565b600060405180830381600087803b15801562001df257600080fd5b505af115801562001e07573d6000803e3d6000fd5b50506040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8d81166004830152602482018b90527f000000000000000000000000000000000000000000000000000000000000000016925063a9059cbb91506044016020604051808303816000875af115801562001ea1573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001ec7919062004d0c565b50505050505050505050505050565b6002336000908152601a602052604090205460ff16600381111562001eff5762001eff620041a1565b1415801562001f3557506003336000908152601a602052604090205460ff16600381111562001f325762001f32620041a1565b14155b1562001f6d576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080808080808062001f83888a018a62004f04565b965096509650965096509650965060005b87518110156200225257600073ffffffffffffffffffffffffffffffffffffffff1687828151811062001fcb5762001fcb620049cc565b60200260200101516060015173ffffffffffffffffffffffffffffffffffffffff1603620020df57858181518110620020085762002008620049cc565b6020026020010151307f0000000000000000000000000000000000000000000000000000000000000000604051620020409062003d8a565b73ffffffffffffffffffffffffffffffffffffffff938416815291831660208301529091166040820152606001604051809103906000f0801580156200208a573d6000803e3d6000fd5b50878281518110620020a057620020a0620049cc565b60200260200101516060019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250505b62002197888281518110620020f857620020f8620049cc565b6020026020010151888381518110620021155762002115620049cc565b6020026020010151878481518110620021325762002132620049cc565b60200260200101518785815181106200214f576200214f620049cc565b60200260200101518786815181106200216c576200216c620049cc565b6020026020010151878781518110620021895762002189620049cc565b602002602001015162002d42565b878181518110620021ac57620021ac620049cc565b60200260200101517f74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71888381518110620021ea57620021ea620049cc565b602002602001015160a0015133604051620022359291906bffffffffffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a280620022498162004a11565b91505062001f94565b50505050505050505050565b600082815260046020908152604091829020825160e081018452815460ff81161515825263ffffffff6101008204811694830194909452650100000000008104841694820185905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004821660c082015291146200235c576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818160a001516200236e919062005035565b600084815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff93841602179055601954620023d691841690620049fb565b6019556040517f23b872dd0000000000000000000000000000000000000000000000000000000081523360048201523060248201526bffffffffffffffffffffffff831660448201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906323b872dd906064016020604051808303816000875af115801562002480573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620024a6919062004d0c565b506040516bffffffffffffffffffffffff83168152339084907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a3505050565b6000818152600460209081526040808320815160e081018352815460ff81161515825263ffffffff6101008204811695830195909552650100000000008104851693820184905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004831660c08201529291141590620025da60005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161490506000601260010160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200267d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620026a391906200505d565b9050828015620026c75750818015620026c5575080846040015163ffffffff16115b155b15620026ff576040517ffbc0357800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8115801562002732575060008581526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314155b156200276a576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8162002780576200277d603282620049fb565b90505b6000858152600460205260409020805463ffffffff80841665010000000000027fffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffffff90921691909117909155620027dc9060029087906200373916565b5060145460808501516bffffffffffffffffffffffff91821691600091168211156200284557608086015162002813908362005077565b90508560a001516bffffffffffffffffffffffff16816bffffffffffffffffffffffff16111562002845575060a08501515b808660a0015162002857919062005077565b600088815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff93841602179055601554620028bf9183911662005035565b601580547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff9290921691909117905560405167ffffffffffffffff84169088907f91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f79118190600090a350505050505050565b600060606000806200294d6200312d565b6000634b56a42e60e01b8888886040516024016200296e939291906200509f565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091529050620029f989826200068d565b929c919b50995090975095505050505050565b62002a1662003747565b62002a2181620037ca565b50565b60006060600080600080600062002a4b8860405180602001604052806000815250620009ff565b959e949d50929b5090995097509550909350915050565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415155b90505b92915050565b6000806000601260010160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905060008173ffffffffffffffffffffffffffffffffffffffff166385df51fd60018473ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa15801562002b2f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062002b5591906200505d565b62002b6191906200473e565b6040518263ffffffff1660e01b815260040162002b8091815260200190565b602060405180830381865afa15801562002b9e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062002bc491906200505d565b601554604080516020810193909352309083015274010000000000000000000000000000000000000000900463ffffffff166060820152608001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201209083015201604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052905060045b600f81101562002cd0578382828151811062002c8c5762002c8c620049cc565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508062002cc78162004a11565b91505062002c6c565b5084600181111562002ce65762002ce6620041a1565b60f81b81600f8151811062002cff5762002cff620049cc565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535062002d3981620050d3565b95945050505050565b6012547e01000000000000000000000000000000000000000000000000000000000000900460ff161562002da2576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601654835163ffffffff909116101562002de8576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108fc856020015163ffffffff16108062002e265750601554602086015163ffffffff70010000000000000000000000000000000090920482169116115b1562002e5e576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000868152600460205260409020546901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff161562002ec8576040517f6e3b930b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000868152600460209081526040808320885181548a8501518b85015160608d01517fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000009093169315157fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000ff169390931761010063ffffffff92831602177fffffff000000000000000000000000000000000000000000000000ffffffffff1665010000000000938216939093027fffffff0000000000000000000000000000000000000000ffffffffffffffffff1692909217690100000000000000000073ffffffffffffffffffffffffffffffffffffffff9283160217835560808b01516001909301805460a08d015160c08e01516bffffffffffffffffffffffff9687167fffffffffffffffff000000000000000000000000000000000000000000000000909316929092176c010000000000000000000000009690911695909502949094177fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff1678010000000000000000000000000000000000000000000000009490931693909302919091179091556005835281842080547fffffffffffffffffffffffff00000000000000000000000000000000000000001691891691909117905560079091529020620030bb848262005116565b508460a001516bffffffffffffffffffffffff16601954620030de9190620049fb565b6019556000868152601b60205260409020620030fb838262005116565b506000868152601c6020526040902062003116828262005116565b5062003124600287620038c1565b50505050505050565b321562003166576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60008181526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314620031c6576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526004602052604090205465010000000000900463ffffffff9081161462002a21576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000818160045b600f811015620032b3577fff000000000000000000000000000000000000000000000000000000000000008216838260208110620032675762003267620049cc565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916146200329e57506000949350505050565b80620032aa8162004a11565b91505062003225565b5081600f1a6001811115620032cc57620032cc620041a1565b949350505050565b6000806000836080015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa15801562003361573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003387919062005258565b50945090925050506000811315806200339f57508142105b80620033c45750828015620033c45750620033bb82426200473e565b8463ffffffff16105b15620033d5576017549550620033d9565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa15801562003445573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200346b919062005258565b50945090925050506000811315806200348357508142105b80620034a85750828015620034a857506200349f82426200473e565b8463ffffffff16105b15620034b9576018549450620034bd565b8094505b50505050915091565b600080620034da88878b60c00151620038cf565b9050600080620034f78b8a63ffffffff16858a8a60018b62003965565b909250905062003508818362005035565b9b9a5050505050505050505050565b60606000836001811115620035305762003530620041a1565b03620035fd576000848152600760205260409081902090517f6e04ff0d0000000000000000000000000000000000000000000000000000000091620035789160240162005350565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152905062003732565b6001836001811115620036145762003614620041a1565b036200370057600082806020019051810190620036329190620053c7565b6000868152600760205260409081902090519192507f40691db4000000000000000000000000000000000000000000000000000000009162003679918491602401620054db565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091529150620037329050565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b9392505050565b600062002a8d838362003c2d565b60005473ffffffffffffffffffffffffffffffffffffffff16331462003166576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401620011e0565b3373ffffffffffffffffffffffffffffffffffffffff8216036200384b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620011e0565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600062002a8d838362003d38565b60008080856001811115620038e857620038e8620041a1565b03620038f9575062015f906200391c565b6001856001811115620039105762003910620041a1565b036200370057506201b5e45b6200392f63ffffffff85166014620055a3565b6200393c846001620055bd565b6200394d9060ff16611d4c620055a3565b620039599083620049fb565b62002d399190620049fb565b60008060008960a0015161ffff1687620039809190620055a3565b90508380156200398f5750803a105b156200399857503a5b6000841562003a30578a610140015173ffffffffffffffffffffffffffffffffffffffff166349948e0e6000366040518363ffffffff1660e01b8152600401620039e4929190620046c8565b602060405180830381865afa15801562003a02573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003a2891906200505d565b905062003af1565b6101408b01516016546040517f1254414000000000000000000000000000000000000000000000000000000000815264010000000090910463ffffffff16600482015273ffffffffffffffffffffffffffffffffffffffff90911690631254414090602401602060405180830381865afa15801562003ab3573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003ad991906200505d565b8b60a0015161ffff1662003aee9190620055a3565b90505b62003b0161ffff871682620055d9565b90506000878262003b138c8e620049fb565b62003b1f9086620055a3565b62003b2b9190620049fb565b62003b3f90670de0b6b3a7640000620055a3565b62003b4b9190620055d9565b905060008c6040015163ffffffff1664e8d4a5100062003b6c9190620055a3565b898e6020015163ffffffff16858f8862003b879190620055a3565b62003b939190620049fb565b62003ba390633b9aca00620055a3565b62003baf9190620055a3565b62003bbb9190620055d9565b62003bc79190620049fb565b90506b033b2e3c9fd0803ce800000062003be28284620049fb565b111562003c1b576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b6000818152600183016020526040812054801562003d2657600062003c546001836200473e565b855490915060009062003c6a906001906200473e565b905081811462003cd657600086600001828154811062003c8e5762003c8e620049cc565b906000526020600020015490508087600001848154811062003cb45762003cb4620049cc565b6000918252602080832090910192909255918252600188019052604090208390555b855486908062003cea5762003cea62005615565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505062002a90565b600091505062002a90565b5092915050565b600081815260018301602052604081205462003d815750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915562002a90565b50600062002a90565b6103ca806200564583390190565b50805462003da690620047fc565b6000825580601f1062003db7575050565b601f01602090049060005260206000209081019062002a2191905b8082111562003de8576000815560010162003dd2565b5090565b73ffffffffffffffffffffffffffffffffffffffff8116811462002a2157600080fd5b803563ffffffff8116811462003e2457600080fd5b919050565b80356002811062003e2457600080fd5b60008083601f84011262003e4c57600080fd5b50813567ffffffffffffffff81111562003e6557600080fd5b60208301915083602082850101111562003e7e57600080fd5b9250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160e0810167ffffffffffffffff8111828210171562003eda5762003eda62003e85565b60405290565b604051610100810167ffffffffffffffff8111828210171562003eda5762003eda62003e85565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171562003f515762003f5162003e85565b604052919050565b600067ffffffffffffffff82111562003f765762003f7662003e85565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f83011262003fb457600080fd5b813562003fcb62003fc58262003f59565b62003f07565b81815284602083860101111562003fe157600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060008060e0898b0312156200401b57600080fd5b8835620040288162003dec565b97506200403860208a0162003e0f565b965060408901356200404a8162003dec565b95506200405a60608a0162003e29565b9450608089013567ffffffffffffffff808211156200407857600080fd5b620040868c838d0162003e39565b909650945060a08b0135915080821115620040a057600080fd5b620040ae8c838d0162003fa2565b935060c08b0135915080821115620040c557600080fd5b50620040d48b828c0162003fa2565b9150509295985092959890939650565b60008060408385031215620040f857600080fd5b82359150602083013567ffffffffffffffff8111156200411757600080fd5b620041258582860162003fa2565b9150509250929050565b60005b838110156200414c57818101518382015260200162004132565b50506000910152565b600081518084526200416f8160208601602086016200412f565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600a811062004208577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b9052565b841515815260806020820152600062004229608083018662004155565b90506200423a6040830185620041d0565b82606083015295945050505050565b6000806000604084860312156200425f57600080fd5b83359250602084013567ffffffffffffffff8111156200427e57600080fd5b6200428c8682870162003e39565b9497909650939450505050565b600080600080600080600060a0888a031215620042b557600080fd5b8735620042c28162003dec565b9650620042d26020890162003e0f565b95506040880135620042e48162003dec565b9450606088013567ffffffffffffffff808211156200430257600080fd5b620043108b838c0162003e39565b909650945060808a01359150808211156200432a57600080fd5b50620043398a828b0162003e39565b989b979a50959850939692959293505050565b871515815260e0602082015260006200436960e083018962004155565b90506200437a6040830188620041d0565b8560608301528460808301528360a08301528260c083015298975050505050505050565b600080600060408486031215620043b457600080fd5b833567ffffffffffffffff80821115620043cd57600080fd5b818601915086601f830112620043e257600080fd5b813581811115620043f257600080fd5b8760208260051b85010111156200440857600080fd5b60209283019550935050840135620044208162003dec565b809150509250925092565b600080602083850312156200443f57600080fd5b823567ffffffffffffffff8111156200445757600080fd5b620044658582860162003e39565b90969095509350505050565b80356bffffffffffffffffffffffff8116811462003e2457600080fd5b60008060408385031215620044a257600080fd5b82359150620044b46020840162004471565b90509250929050565b600060208284031215620044d057600080fd5b5035919050565b600067ffffffffffffffff821115620044f457620044f462003e85565b5060051b60200190565b600082601f8301126200451057600080fd5b813560206200452362003fc583620044d7565b82815260059290921b840181019181810190868411156200454357600080fd5b8286015b848110156200458857803567ffffffffffffffff811115620045695760008081fd5b620045798986838b010162003fa2565b84525091830191830162004547565b509695505050505050565b60008060008060608587031215620045aa57600080fd5b84359350602085013567ffffffffffffffff80821115620045ca57600080fd5b620045d888838901620044fe565b94506040870135915080821115620045ef57600080fd5b50620045fe8782880162003e39565b95989497509550505050565b6000602082840312156200461d57600080fd5b8135620037328162003dec565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600063ffffffff8083168181036200467557620046756200462a565b6001019392505050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b602081526000620032cc6020830184866200467f565b60208152600062002a8d602083018462004155565b805162003e248162003dec565b6000602082840312156200471357600080fd5b8151620037328162003dec565b60008251620047348184602087016200412f565b9190910192915050565b8181038181111562002a905762002a906200462a565b801515811462002a2157600080fd5b600082601f8301126200477557600080fd5b81516200478662003fc58262003f59565b8181528460208386010111156200479c57600080fd5b620032cc8260208301602087016200412f565b60008060408385031215620047c357600080fd5b8251620047d08162004754565b602084015190925067ffffffffffffffff811115620047ee57600080fd5b620041258582860162004763565b600181811c908216806200481157607f821691505b6020821081036200484b577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f8211156200489f57600081815260208120601f850160051c810160208610156200487a5750805b601f850160051c820191505b818110156200489b5782815560010162004886565b5050505b505050565b67ffffffffffffffff831115620048bf57620048bf62003e85565b620048d783620048d08354620047fc565b8362004851565b6000601f8411600181146200492c5760008515620048f55750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355620049c5565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b828110156200497d57868501358255602094850194600190920191016200495b565b5086821015620049b9577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b8082018082111562002a905762002a906200462a565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820362004a455762004a456200462a565b5060010190565b600081518084526020808501945080840160005b8381101562004b0b5781518051151588528381015163ffffffff908116858a01526040808301519091169089015260608082015173ffffffffffffffffffffffffffffffffffffffff16908901526080808201516bffffffffffffffffffffffff169089015260a08082015162004ae6828b01826bffffffffffffffffffffffff169052565b505060c09081015163ffffffff169088015260e0909601959082019060010162004a60565b509495945050505050565b600081518084526020808501945080840160005b8381101562004b0b57815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010162004b2a565b600081518084526020808501808196508360051b8101915082860160005b8581101562004baa57828403895262004b9784835162004155565b9885019893509084019060010162004b7c565b5091979650505050505050565b60e081528760e082015260006101007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8a111562004bf457600080fd5b8960051b808c8386013783018381038201602085015262004c188282018b62004a4c565b915050828103604084015262004c2f818962004b16565b9050828103606084015262004c45818862004b16565b9050828103608084015262004c5b818762004b5e565b905082810360a084015262004c71818662004b5e565b905082810360c084015262003508818562004b5e565b60006020828403121562004c9a57600080fd5b815160ff811681146200373257600080fd5b60ff8416815260ff8316602082015260606040820152600062002d39606083018462004155565b60006020828403121562004ce657600080fd5b815167ffffffffffffffff81111562004cfe57600080fd5b620032cc8482850162004763565b60006020828403121562004d1f57600080fd5b8151620037328162004754565b600082601f83011262004d3e57600080fd5b8135602062004d5162003fc583620044d7565b82815260059290921b8401810191818101908684111562004d7157600080fd5b8286015b8481101562004588578035835291830191830162004d75565b600082601f83011262004da057600080fd5b8135602062004db362003fc583620044d7565b82815260e0928302850182019282820191908785111562004dd357600080fd5b8387015b8581101562004e8a5781818a03121562004df15760008081fd5b62004dfb62003eb4565b813562004e088162004754565b815262004e1782870162003e0f565b86820152604062004e2a81840162003e0f565b9082015260608281013562004e3f8162003dec565b90820152608062004e5283820162004471565b9082015260a062004e6583820162004471565b9082015260c062004e7883820162003e0f565b90820152845292840192810162004dd7565b5090979650505050505050565b600082601f83011262004ea957600080fd5b8135602062004ebc62003fc583620044d7565b82815260059290921b8401810191818101908684111562004edc57600080fd5b8286015b848110156200458857803562004ef68162003dec565b835291830191830162004ee0565b600080600080600080600060e0888a03121562004f2057600080fd5b873567ffffffffffffffff8082111562004f3957600080fd5b62004f478b838c0162004d2c565b985060208a013591508082111562004f5e57600080fd5b62004f6c8b838c0162004d8e565b975060408a013591508082111562004f8357600080fd5b62004f918b838c0162004e97565b965060608a013591508082111562004fa857600080fd5b62004fb68b838c0162004e97565b955060808a013591508082111562004fcd57600080fd5b62004fdb8b838c01620044fe565b945060a08a013591508082111562004ff257600080fd5b620050008b838c01620044fe565b935060c08a01359150808211156200501757600080fd5b50620050268a828b01620044fe565b91505092959891949750929550565b6bffffffffffffffffffffffff81811683821601908082111562003d315762003d316200462a565b6000602082840312156200507057600080fd5b5051919050565b6bffffffffffffffffffffffff82811682821603908082111562003d315762003d316200462a565b604081526000620050b4604083018662004b5e565b8281036020840152620050c98185876200467f565b9695505050505050565b805160208083015191908110156200484b577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60209190910360031b1b16919050565b815167ffffffffffffffff81111562005133576200513362003e85565b6200514b81620051448454620047fc565b8462004851565b602080601f831160018114620051a157600084156200516a5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556200489b565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015620051f057888601518255948401946001909101908401620051cf565b50858210156200522d57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b805169ffffffffffffffffffff8116811462003e2457600080fd5b600080600080600060a086880312156200527157600080fd5b6200527c866200523d565b9450602086015193506040860151925060608601519150620052a1608087016200523d565b90509295509295909350565b60008154620052bc81620047fc565b808552602060018381168015620052dc5760018114620053155762005345565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b890101955062005345565b866000528260002060005b858110156200533d5781548a820186015290830190840162005320565b890184019650505b505050505092915050565b60208152600062002a8d6020830184620052ad565b600082601f8301126200537757600080fd5b815160206200538a62003fc583620044d7565b82815260059290921b84018101918181019086841115620053aa57600080fd5b8286015b84811015620045885780518352918301918301620053ae565b600060208284031215620053da57600080fd5b815167ffffffffffffffff80821115620053f357600080fd5b9083019061010082860312156200540957600080fd5b6200541362003ee0565b82518152602083015160208201526040830151604082015260608301516060820152608083015160808201526200544d60a08401620046f3565b60a082015260c0830151828111156200546557600080fd5b620054738782860162005365565b60c08301525060e0830151828111156200548c57600080fd5b6200549a8782860162004763565b60e08301525095945050505050565b600081518084526020808501945080840160005b8381101562004b0b57815187529582019590820190600101620054bd565b60408152825160408201526020830151606082015260408301516080820152606083015160a0820152608083015160c082015273ffffffffffffffffffffffffffffffffffffffff60a08401511660e0820152600060c08401516101008081850152506200554e610140840182620054a9565b905060e08501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0848303016101208501526200558c828262004155565b915050828103602084015262002d398185620052ad565b808202811582820484141762002a905762002a906200462a565b60ff818116838216019081111562002a905762002a906200462a565b60008262005610577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfe60c060405234801561001057600080fd5b506040516103ca3803806103ca83398101604081905261002f91610076565b600080546001600160a01b0319166001600160a01b039384161790559181166080521660a0526100b9565b80516001600160a01b038116811461007157600080fd5b919050565b60008060006060848603121561008b57600080fd5b6100948461005a565b92506100a26020850161005a565b91506100b06040850161005a565b90509250925092565b60805160a0516102e76100e36000396000603801526000818160c4015261011701526102e76000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806379188d161461007b578063f00e6a2a146100aa575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e808015610076573d6000f35b3d6000fd5b61008e6100893660046101c1565b6100ee565b6040805192151583526020830191909152015b60405180910390f35b60405173ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001681526020016100a1565b60008054819073ffffffffffffffffffffffffffffffffffffffff16331461011557600080fd5b7f00000000000000000000000000000000000000000000000000000000000000005a91505a61138881101561014957600080fd5b61138881039050856040820482031161016157600080fd5b50803b61016d57600080fd5b6000808551602087016000858af192505a610188908361029a565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080604083850312156101d457600080fd5b82359150602083013567ffffffffffffffff808211156101f357600080fd5b818501915085601f83011261020757600080fd5b81358181111561021957610219610192565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561025f5761025f610192565b8160405282815288602084870101111561027857600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b818103818111156102d4577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b9291505056fea164736f6c6343000813000aa164736f6c6343000813000a",
>>>>>>> e15f3cf45d (regen wrappers)
=======
	Bin: "0x6101206040523480156200001257600080fd5b5060405162005dd138038062005dd1833981016040819052620000359162000345565b80816001600160a01b031663ca30e6036040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200009b919062000345565b826001600160a01b031663b10b673c6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000da573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000100919062000345565b836001600160a01b0316636709d0e56040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200013f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000165919062000345565b846001600160a01b0316635425d8ac6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca919062000345565b3380600081620002215760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200025457620002548162000281565b5050506001600160a01b0393841660805291831660a052821660c052811660e0521661010052506200036c565b336001600160a01b03821603620002db5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000218565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200034257600080fd5b50565b6000602082840312156200035857600080fd5b815162000365816200032c565b9392505050565b60805160a05160c05160e05161010051615a0b620003c66000396000818161010e01526101a90152600081816103e10152612013015260006132f7015260006133db015260008181611e5501526124210152615a0b6000f3fe60806040523480156200001157600080fd5b50600436106200010c5760003560e01c806385c1b0ba11620000a5578063c8048022116200006f578063c804802214620002b7578063ce7dc5b414620002ce578063f2fde38b14620002e5578063f7d334ba14620002fc576200010c565b806385c1b0ba14620002535780638da5cb5b146200026a5780638e86139b1462000289578063948108f714620002a0576200010c565b80634ee88d3511620000e75780634ee88d3514620001ef5780636ded9eae146200020657806371791aa0146200021d57806379ba50971462000249576200010c565b806328f32f38146200015457806329c5efad146200017e578063349e8cca14620001a7575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e8080156200014d573d6000f35b3d6000fd5b005b6200016b6200016536600462003fee565b62000313565b6040519081526020015b60405180910390f35b620001956200018f366004620040d4565b6200068d565b604051620001759493929190620041fc565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200162000175565b620001526200020036600462004239565b62000931565b6200016b6200021736600462004289565b62000999565b620002346200022e366004620040d4565b620009ff565b6040516200017597969594939291906200433c565b6200015262001162565b62000152620002643660046200438e565b62001265565b60005473ffffffffffffffffffffffffffffffffffffffff16620001c9565b620001526200029a3660046200441b565b62001ed6565b62000152620002b13660046200447e565b6200225e565b62000152620002c8366004620044ad565b620024f1565b62000195620002df36600462004583565b6200293c565b62000152620002f6366004620045fa565b62002a0c565b620002346200030d366004620044ad565b62002a24565b6000805473ffffffffffffffffffffffffffffffffffffffff1633148015906200034757506200034560093362002a62565b155b156200037f576040517fd48b678b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff89163b620003ce576040517f09ee12d500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b620003d98662002a96565b9050600089307f00000000000000000000000000000000000000000000000000000000000000006040516200040e9062003d7a565b73ffffffffffffffffffffffffffffffffffffffff938416815291831660208301529091166040820152606001604051809103906000f08015801562000458573d6000803e3d6000fd5b5090506200051f826040518060e001604052806000151581526020018c63ffffffff16815260200163ffffffff801681526020018473ffffffffffffffffffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff168152602001600063ffffffff168152508a89898080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508b92508a915062002d429050565b6015805474010000000000000000000000000000000000000000900463ffffffff169060146200054f8362004649565b91906101000a81548163ffffffff021916908363ffffffff16021790555050817fbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d0128a8a604051620005c892919063ffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a2817fcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d878760405162000604929190620046b8565b60405180910390a2817f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664856040516200063e9190620046ce565b60405180910390a2817f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf485084604051620006789190620046ce565b60405180910390a25098975050505050505050565b600060606000806200069e6200312d565b600086815260046020908152604091829020825160e081018452815460ff811615158252610100810463ffffffff90811694830194909452650100000000008104841694820194909452690100000000000000000090930473ffffffffffffffffffffffffffffffffffffffff166060840152600101546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a0840152780100000000000000000000000000000000000000000000000090041660c08201525a9150600080826060015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620007b8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620007de9190620046f0565b73ffffffffffffffffffffffffffffffffffffffff166014600101600c9054906101000a900463ffffffff1663ffffffff168960405162000820919062004710565b60006040518083038160008787f1925050503d806000811462000860576040519150601f19603f3d011682016040523d82523d6000602084013e62000865565b606091505b50915091505a6200087790856200472e565b935081620008a257600060405180602001604052806000815250600796509650965050505062000928565b80806020019051810190620008b891906200479f565b909750955086620008e657600060405180602001604052806000815250600496509650965050505062000928565b601654865164010000000090910463ffffffff1610156200092457600060405180602001604052806000815250600596509650965050505062000928565b5050505b92959194509250565b6200093c8362003168565b6000838152601b602052604090206200095782848362004894565b50827f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d566483836040516200098c929190620046b8565b60405180910390a2505050565b6000620009f388888860008989604051806020016040528060008152508a8a8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506200031392505050565b98975050505050505050565b60006060600080600080600062000a156200312d565b600062000a228a6200321e565b905060006012604051806101600160405290816000820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160008201600c9054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160109054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160149054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160189054906101000a900462ffffff1662ffffff1662ffffff16815260200160008201601b9054906101000a900461ffff1661ffff1661ffff16815260200160008201601d9054906101000a900460ff1660ff1660ff16815260200160008201601e9054906101000a900460ff1615151515815260200160008201601f9054906101000a900460ff161515151581526020016001820160009054906101000a900460ff161515151581526020016001820160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152505090506000600460008d81526020019081526020016000206040518060e00160405290816000820160009054906101000a900460ff161515151581526020016000820160019054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160059054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160099054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160018201600c9054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff1681526020016001820160189054906101000a900463ffffffff1663ffffffff1663ffffffff168152505090508160e001511562000db7576000604051806020016040528060008152506009600084602001516000808263ffffffff169250995099509950995099509950995050505062001156565b604081015163ffffffff9081161462000e08576000604051806020016040528060008152506001600084602001516000808263ffffffff169250995099509950995099509950995050505062001156565b80511562000e4e576000604051806020016040528060008152506002600084602001516000808263ffffffff169250995099509950995099509950995050505062001156565b62000e5982620032d4565b602083015160165492975090955060009162000e8b918591879190640100000000900463ffffffff168a8a87620034c6565b9050806bffffffffffffffffffffffff168260a001516bffffffffffffffffffffffff16101562000ef5576000604051806020016040528060008152506006600085602001516000808263ffffffff1692509a509a509a509a509a509a509a505050505062001156565b600062000f048e868f62003517565b90505a9850600080846060015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000f5c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000f829190620046f0565b73ffffffffffffffffffffffffffffffffffffffff166014600101600c9054906101000a900463ffffffff1663ffffffff168460405162000fc4919062004710565b60006040518083038160008787f1925050503d806000811462001004576040519150601f19603f3d011682016040523d82523d6000602084013e62001009565b606091505b50915091505a6200101b908c6200472e565b9a50816200109b5760165481516801000000000000000090910463ffffffff1610156200107857505060408051602080820190925260008082529490910151939c509a50600899505063ffffffff90911695506200115692505050565b602090940151939b5060039a505063ffffffff9092169650620011569350505050565b80806020019051810190620010b191906200479f565b909e509c508d620010f257505060408051602080820190925260008082529490910151939c509a50600499505063ffffffff90911695506200115692505050565b6016548d5164010000000090910463ffffffff1610156200114357505060408051602080820190925260008082529490910151939c509a50600599505063ffffffff90911695506200115692505050565b505050506020015163ffffffff16945050505b92959891949750929550565b60015473ffffffffffffffffffffffffffffffffffffffff163314620011e9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600173ffffffffffffffffffffffffffffffffffffffff82166000908152601a602052604090205460ff166003811115620012a457620012a462004191565b14158015620012f05750600373ffffffffffffffffffffffffffffffffffffffff82166000908152601a602052604090205460ff166003811115620012ed57620012ed62004191565b14155b1562001328576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6014546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1662001388576040517fd12d7d8d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000829003620013c4576040517f2c2fc94100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160e081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c081018290526000808567ffffffffffffffff8111156200141b576200141b62003e75565b60405190808252806020026020018201604052801562001445578160200160208202803683370190505b50905060008667ffffffffffffffff81111562001466576200146662003e75565b604051908082528060200260200182016040528015620014ed57816020015b6040805160e08101825260008082526020808301829052928201819052606082018190526080820181905260a0820181905260c082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181620014855790505b50905060008767ffffffffffffffff8111156200150e576200150e62003e75565b6040519080825280602002602001820160405280156200154357816020015b60608152602001906001900390816200152d5790505b50905060008867ffffffffffffffff81111562001564576200156462003e75565b6040519080825280602002602001820160405280156200159957816020015b6060815260200190600190039081620015835790505b50905060008967ffffffffffffffff811115620015ba57620015ba62003e75565b604051908082528060200260200182016040528015620015ef57816020015b6060815260200190600190039081620015d95790505b50905060005b8a81101562001bd3578b8b82818110620016135762001613620049bc565b60209081029290920135600081815260048452604090819020815160e081018352815460ff811615158252610100810463ffffffff90811697830197909752650100000000008104871693820193909352690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff166060830152600101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490931660c08401529a50909850620016f290508962003168565b60608801516040517f1a5da6c800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8c8116600483015290911690631a5da6c890602401600060405180830381600087803b1580156200176257600080fd5b505af115801562001777573d6000803e3d6000fd5b5050505087858281518110620017915762001791620049bc565b6020026020010181905250600560008a815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16868281518110620017e557620017e5620049bc565b73ffffffffffffffffffffffffffffffffffffffff90921660209283029190910182015260008a815260079091526040902080546200182490620047ec565b80601f01602080910402602001604051908101604052809291908181526020018280546200185290620047ec565b8015620018a35780601f106200187757610100808354040283529160200191620018a3565b820191906000526020600020905b8154815290600101906020018083116200188557829003601f168201915b5050505050848281518110620018bd57620018bd620049bc565b6020026020010181905250601b60008a81526020019081526020016000208054620018e890620047ec565b80601f01602080910402602001604051908101604052809291908181526020018280546200191690620047ec565b8015620019675780601f106200193b5761010080835404028352916020019162001967565b820191906000526020600020905b8154815290600101906020018083116200194957829003601f168201915b5050505050838281518110620019815762001981620049bc565b6020026020010181905250601c60008a81526020019081526020016000208054620019ac90620047ec565b80601f0160208091040260200160405190810160405280929190818152602001828054620019da90620047ec565b801562001a2b5780601f10620019ff5761010080835404028352916020019162001a2b565b820191906000526020600020905b81548152906001019060200180831162001a0d57829003601f168201915b505050505082828151811062001a455762001a45620049bc565b60200260200101819052508760a001516bffffffffffffffffffffffff168762001a709190620049eb565b60008a815260046020908152604080832080547fffffff000000000000000000000000000000000000000000000000000000000016815560010180547fffffffff000000000000000000000000000000000000000000000000000000001690556007909152812091985062001ae6919062003d88565b6000898152601b6020526040812062001aff9162003d88565b6000898152601c6020526040812062001b189162003d88565b600089815260066020526040902080547fffffffffffffffffffffffff000000000000000000000000000000000000000016905562001b5960028a62003739565b5060a0880151604080516bffffffffffffffffffffffff909216825273ffffffffffffffffffffffffffffffffffffffff8c1660208301528a917fb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff910160405180910390a28062001bca8162004a01565b915050620015f5565b508560195462001be491906200472e565b60195560008b8b868167ffffffffffffffff81111562001c085762001c0862003e75565b60405190808252806020026020018201604052801562001c32578160200160208202803683370190505b508988888860405160200162001c5098979695949392919062004ba7565b60405160208183030381529060405290508973ffffffffffffffffffffffffffffffffffffffff16638e86139b6014600001600c9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663c71249ab60038e73ffffffffffffffffffffffffffffffffffffffff1663aab9edd66040518163ffffffff1660e01b8152600401602060405180830381865afa15801562001d0c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001d32919062004c77565b866040518463ffffffff1660e01b815260040162001d539392919062004c9c565b600060405180830381865afa15801562001d71573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405262001db9919081019062004cc3565b6040518263ffffffff1660e01b815260040162001dd79190620046ce565b600060405180830381600087803b15801562001df257600080fd5b505af115801562001e07573d6000803e3d6000fd5b50506040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8d81166004830152602482018b90527f000000000000000000000000000000000000000000000000000000000000000016925063a9059cbb91506044016020604051808303816000875af115801562001ea1573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001ec7919062004cfc565b50505050505050505050505050565b6002336000908152601a602052604090205460ff16600381111562001eff5762001eff62004191565b1415801562001f3557506003336000908152601a602052604090205460ff16600381111562001f325762001f3262004191565b14155b1562001f6d576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080808080808062001f83888a018a62004ef4565b965096509650965096509650965060005b87518110156200225257600073ffffffffffffffffffffffffffffffffffffffff1687828151811062001fcb5762001fcb620049bc565b60200260200101516060015173ffffffffffffffffffffffffffffffffffffffff1603620020df57858181518110620020085762002008620049bc565b6020026020010151307f0000000000000000000000000000000000000000000000000000000000000000604051620020409062003d7a565b73ffffffffffffffffffffffffffffffffffffffff938416815291831660208301529091166040820152606001604051809103906000f0801580156200208a573d6000803e3d6000fd5b50878281518110620020a057620020a0620049bc565b60200260200101516060019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250505b62002197888281518110620020f857620020f8620049bc565b6020026020010151888381518110620021155762002115620049bc565b6020026020010151878481518110620021325762002132620049bc565b60200260200101518785815181106200214f576200214f620049bc565b60200260200101518786815181106200216c576200216c620049bc565b6020026020010151878781518110620021895762002189620049bc565b602002602001015162002d42565b878181518110620021ac57620021ac620049bc565b60200260200101517f74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71888381518110620021ea57620021ea620049bc565b602002602001015160a0015133604051620022359291906bffffffffffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a280620022498162004a01565b91505062001f94565b50505050505050505050565b600082815260046020908152604091829020825160e081018452815460ff81161515825263ffffffff6101008204811694830194909452650100000000008104841694820185905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004821660c082015291146200235c576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818160a001516200236e919062005025565b600084815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff93841602179055601954620023d691841690620049eb565b6019556040517f23b872dd0000000000000000000000000000000000000000000000000000000081523360048201523060248201526bffffffffffffffffffffffff831660448201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906323b872dd906064016020604051808303816000875af115801562002480573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620024a6919062004cfc565b506040516bffffffffffffffffffffffff83168152339084907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a3505050565b6000818152600460209081526040808320815160e081018352815460ff81161515825263ffffffff6101008204811695830195909552650100000000008104851693820184905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004831660c08201529291141590620025da60005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161490506000601260010160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200267d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620026a391906200504d565b9050828015620026c75750818015620026c5575080846040015163ffffffff16115b155b15620026ff576040517ffbc0357800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8115801562002732575060008581526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314155b156200276a576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8162002780576200277d603282620049eb565b90505b6000858152600460205260409020805463ffffffff80841665010000000000027fffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffffff90921691909117909155620027dc9060029087906200373916565b5060145460808501516bffffffffffffffffffffffff91821691600091168211156200284557608086015162002813908362005067565b90508560a001516bffffffffffffffffffffffff16816bffffffffffffffffffffffff16111562002845575060a08501515b808660a0015162002857919062005067565b600088815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff93841602179055601554620028bf9183911662005025565b601580547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff9290921691909117905560405167ffffffffffffffff84169088907f91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f79118190600090a350505050505050565b600060606000806200294d6200312d565b6000634b56a42e60e01b8888886040516024016200296e939291906200508f565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091529050620029f989826200068d565b929c919b50995090975095505050505050565b62002a1662003747565b62002a2181620037ca565b50565b60006060600080600080600062002a4b8860405180602001604052806000815250620009ff565b959e949d50929b5090995097509550909350915050565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415155b90505b92915050565b6000806000601260010160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905060008173ffffffffffffffffffffffffffffffffffffffff166385df51fd60018473ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa15801562002b2f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062002b5591906200504d565b62002b6191906200472e565b6040518263ffffffff1660e01b815260040162002b8091815260200190565b602060405180830381865afa15801562002b9e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062002bc491906200504d565b601554604080516020810193909352309083015274010000000000000000000000000000000000000000900463ffffffff166060820152608001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201209083015201604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052905060045b600f81101562002cd0578382828151811062002c8c5762002c8c620049bc565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508062002cc78162004a01565b91505062002c6c565b5084600181111562002ce65762002ce662004191565b60f81b81600f8151811062002cff5762002cff620049bc565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535062002d3981620050c3565b95945050505050565b6012547e01000000000000000000000000000000000000000000000000000000000000900460ff161562002da2576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601654835163ffffffff909116101562002de8576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108fc856020015163ffffffff16108062002e265750601554602086015163ffffffff70010000000000000000000000000000000090920482169116115b1562002e5e576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000868152600460205260409020546901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff161562002ec8576040517f6e3b930b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000868152600460209081526040808320885181548a8501518b85015160608d01517fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000009093169315157fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000ff169390931761010063ffffffff92831602177fffffff000000000000000000000000000000000000000000000000ffffffffff1665010000000000938216939093027fffffff0000000000000000000000000000000000000000ffffffffffffffffff1692909217690100000000000000000073ffffffffffffffffffffffffffffffffffffffff9283160217835560808b01516001909301805460a08d015160c08e01516bffffffffffffffffffffffff9687167fffffffffffffffff000000000000000000000000000000000000000000000000909316929092176c010000000000000000000000009690911695909502949094177fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff1678010000000000000000000000000000000000000000000000009490931693909302919091179091556005835281842080547fffffffffffffffffffffffff00000000000000000000000000000000000000001691891691909117905560079091529020620030bb848262005106565b508460a001516bffffffffffffffffffffffff16601954620030de9190620049eb565b6019556000868152601b60205260409020620030fb838262005106565b506000868152601c6020526040902062003116828262005106565b5062003124600287620038c1565b50505050505050565b321562003166576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60008181526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314620031c6576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526004602052604090205465010000000000900463ffffffff9081161462002a21576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000818160045b600f811015620032b3577fff000000000000000000000000000000000000000000000000000000000000008216838260208110620032675762003267620049bc565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916146200329e57506000949350505050565b80620032aa8162004a01565b91505062003225565b5081600f1a6001811115620032cc57620032cc62004191565b949350505050565b6000806000836080015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa15801562003361573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003387919062005248565b50945090925050506000811315806200339f57508142105b80620033c45750828015620033c45750620033bb82426200472e565b8463ffffffff16105b15620033d5576017549550620033d9565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa15801562003445573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200346b919062005248565b50945090925050506000811315806200348357508142105b80620034a85750828015620034a857506200349f82426200472e565b8463ffffffff16105b15620034b9576018549450620034bd565b8094505b50505050915091565b600080620034da88878b60c00151620038cf565b9050600080620034f78b8a63ffffffff16858a8a60018b62003965565b909250905062003508818362005025565b9b9a5050505050505050505050565b6060600083600181111562003530576200353062004191565b03620035fd576000848152600760205260409081902090517f6e04ff0d0000000000000000000000000000000000000000000000000000000091620035789160240162005340565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152905062003732565b600183600181111562003614576200361462004191565b036200370057600082806020019051810190620036329190620053b7565b6000868152600760205260409081902090519192507f40691db4000000000000000000000000000000000000000000000000000000009162003679918491602401620054cb565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091529150620037329050565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b9392505050565b600062002a8d838362003c1d565b60005473ffffffffffffffffffffffffffffffffffffffff16331462003166576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401620011e0565b3373ffffffffffffffffffffffffffffffffffffffff8216036200384b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620011e0565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600062002a8d838362003d28565b60008080856001811115620038e857620038e862004191565b03620038f9575062015f906200391c565b600185600181111562003910576200391062004191565b036200370057506201b38c5b6200392f63ffffffff8516601462005593565b6200393c846001620055ad565b6200394d9060ff16611d4c62005593565b620039599083620049eb565b62002d399190620049eb565b60008060008960a0015161ffff168762003980919062005593565b90508380156200398f5750803a105b156200399857503a5b6000841562003a20578a610140015173ffffffffffffffffffffffffffffffffffffffff166318b8f6136040518163ffffffff1660e01b8152600401602060405180830381865afa158015620039f2573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003a1891906200504d565b905062003ae1565b6101408b01516016546040517f1254414000000000000000000000000000000000000000000000000000000000815264010000000090910463ffffffff16600482015273ffffffffffffffffffffffffffffffffffffffff90911690631254414090602401602060405180830381865afa15801562003aa3573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003ac991906200504d565b8b60a0015161ffff1662003ade919062005593565b90505b62003af161ffff871682620055c9565b90506000878262003b038c8e620049eb565b62003b0f908662005593565b62003b1b9190620049eb565b62003b2f90670de0b6b3a764000062005593565b62003b3b9190620055c9565b905060008c6040015163ffffffff1664e8d4a5100062003b5c919062005593565b898e6020015163ffffffff16858f8862003b77919062005593565b62003b839190620049eb565b62003b9390633b9aca0062005593565b62003b9f919062005593565b62003bab9190620055c9565b62003bb79190620049eb565b90506b033b2e3c9fd0803ce800000062003bd28284620049eb565b111562003c0b576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b6000818152600183016020526040812054801562003d1657600062003c446001836200472e565b855490915060009062003c5a906001906200472e565b905081811462003cc657600086600001828154811062003c7e5762003c7e620049bc565b906000526020600020015490508087600001848154811062003ca45762003ca4620049bc565b6000918252602080832090910192909255918252600188019052604090208390555b855486908062003cda5762003cda62005605565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505062002a90565b600091505062002a90565b5092915050565b600081815260018301602052604081205462003d715750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915562002a90565b50600062002a90565b6103ca806200563583390190565b50805462003d9690620047ec565b6000825580601f1062003da7575050565b601f01602090049060005260206000209081019062002a2191905b8082111562003dd8576000815560010162003dc2565b5090565b73ffffffffffffffffffffffffffffffffffffffff8116811462002a2157600080fd5b803563ffffffff8116811462003e1457600080fd5b919050565b80356002811062003e1457600080fd5b60008083601f84011262003e3c57600080fd5b50813567ffffffffffffffff81111562003e5557600080fd5b60208301915083602082850101111562003e6e57600080fd5b9250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160e0810167ffffffffffffffff8111828210171562003eca5762003eca62003e75565b60405290565b604051610100810167ffffffffffffffff8111828210171562003eca5762003eca62003e75565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171562003f415762003f4162003e75565b604052919050565b600067ffffffffffffffff82111562003f665762003f6662003e75565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f83011262003fa457600080fd5b813562003fbb62003fb58262003f49565b62003ef7565b81815284602083860101111562003fd157600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060008060e0898b0312156200400b57600080fd5b8835620040188162003ddc565b97506200402860208a0162003dff565b965060408901356200403a8162003ddc565b95506200404a60608a0162003e19565b9450608089013567ffffffffffffffff808211156200406857600080fd5b620040768c838d0162003e29565b909650945060a08b01359150808211156200409057600080fd5b6200409e8c838d0162003f92565b935060c08b0135915080821115620040b557600080fd5b50620040c48b828c0162003f92565b9150509295985092959890939650565b60008060408385031215620040e857600080fd5b82359150602083013567ffffffffffffffff8111156200410757600080fd5b620041158582860162003f92565b9150509250929050565b60005b838110156200413c57818101518382015260200162004122565b50506000910152565b600081518084526200415f8160208601602086016200411f565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600a8110620041f8577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b9052565b841515815260806020820152600062004219608083018662004145565b90506200422a6040830185620041c0565b82606083015295945050505050565b6000806000604084860312156200424f57600080fd5b83359250602084013567ffffffffffffffff8111156200426e57600080fd5b6200427c8682870162003e29565b9497909650939450505050565b600080600080600080600060a0888a031215620042a557600080fd5b8735620042b28162003ddc565b9650620042c26020890162003dff565b95506040880135620042d48162003ddc565b9450606088013567ffffffffffffffff80821115620042f257600080fd5b620043008b838c0162003e29565b909650945060808a01359150808211156200431a57600080fd5b50620043298a828b0162003e29565b989b979a50959850939692959293505050565b871515815260e0602082015260006200435960e083018962004145565b90506200436a6040830188620041c0565b8560608301528460808301528360a08301528260c083015298975050505050505050565b600080600060408486031215620043a457600080fd5b833567ffffffffffffffff80821115620043bd57600080fd5b818601915086601f830112620043d257600080fd5b813581811115620043e257600080fd5b8760208260051b8501011115620043f857600080fd5b60209283019550935050840135620044108162003ddc565b809150509250925092565b600080602083850312156200442f57600080fd5b823567ffffffffffffffff8111156200444757600080fd5b620044558582860162003e29565b90969095509350505050565b80356bffffffffffffffffffffffff8116811462003e1457600080fd5b600080604083850312156200449257600080fd5b82359150620044a46020840162004461565b90509250929050565b600060208284031215620044c057600080fd5b5035919050565b600067ffffffffffffffff821115620044e457620044e462003e75565b5060051b60200190565b600082601f8301126200450057600080fd5b813560206200451362003fb583620044c7565b82815260059290921b840181019181810190868411156200453357600080fd5b8286015b848110156200457857803567ffffffffffffffff811115620045595760008081fd5b620045698986838b010162003f92565b84525091830191830162004537565b509695505050505050565b600080600080606085870312156200459a57600080fd5b84359350602085013567ffffffffffffffff80821115620045ba57600080fd5b620045c888838901620044ee565b94506040870135915080821115620045df57600080fd5b50620045ee8782880162003e29565b95989497509550505050565b6000602082840312156200460d57600080fd5b8135620037328162003ddc565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600063ffffffff8083168181036200466557620046656200461a565b6001019392505050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b602081526000620032cc6020830184866200466f565b60208152600062002a8d602083018462004145565b805162003e148162003ddc565b6000602082840312156200470357600080fd5b8151620037328162003ddc565b60008251620047248184602087016200411f565b9190910192915050565b8181038181111562002a905762002a906200461a565b801515811462002a2157600080fd5b600082601f8301126200476557600080fd5b81516200477662003fb58262003f49565b8181528460208386010111156200478c57600080fd5b620032cc8260208301602087016200411f565b60008060408385031215620047b357600080fd5b8251620047c08162004744565b602084015190925067ffffffffffffffff811115620047de57600080fd5b620041158582860162004753565b600181811c908216806200480157607f821691505b6020821081036200483b577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f8211156200488f57600081815260208120601f850160051c810160208610156200486a5750805b601f850160051c820191505b818110156200488b5782815560010162004876565b5050505b505050565b67ffffffffffffffff831115620048af57620048af62003e75565b620048c783620048c08354620047ec565b8362004841565b6000601f8411600181146200491c5760008515620048e55750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355620049b5565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b828110156200496d57868501358255602094850194600190920191016200494b565b5086821015620049a9577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b8082018082111562002a905762002a906200461a565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820362004a355762004a356200461a565b5060010190565b600081518084526020808501945080840160005b8381101562004afb5781518051151588528381015163ffffffff908116858a01526040808301519091169089015260608082015173ffffffffffffffffffffffffffffffffffffffff16908901526080808201516bffffffffffffffffffffffff169089015260a08082015162004ad6828b01826bffffffffffffffffffffffff169052565b505060c09081015163ffffffff169088015260e0909601959082019060010162004a50565b509495945050505050565b600081518084526020808501945080840160005b8381101562004afb57815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010162004b1a565b600081518084526020808501808196508360051b8101915082860160005b8581101562004b9a57828403895262004b8784835162004145565b9885019893509084019060010162004b6c565b5091979650505050505050565b60e081528760e082015260006101007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8a111562004be457600080fd5b8960051b808c8386013783018381038201602085015262004c088282018b62004a3c565b915050828103604084015262004c1f818962004b06565b9050828103606084015262004c35818862004b06565b9050828103608084015262004c4b818762004b4e565b905082810360a084015262004c61818662004b4e565b905082810360c084015262003508818562004b4e565b60006020828403121562004c8a57600080fd5b815160ff811681146200373257600080fd5b60ff8416815260ff8316602082015260606040820152600062002d39606083018462004145565b60006020828403121562004cd657600080fd5b815167ffffffffffffffff81111562004cee57600080fd5b620032cc8482850162004753565b60006020828403121562004d0f57600080fd5b8151620037328162004744565b600082601f83011262004d2e57600080fd5b8135602062004d4162003fb583620044c7565b82815260059290921b8401810191818101908684111562004d6157600080fd5b8286015b8481101562004578578035835291830191830162004d65565b600082601f83011262004d9057600080fd5b8135602062004da362003fb583620044c7565b82815260e0928302850182019282820191908785111562004dc357600080fd5b8387015b8581101562004e7a5781818a03121562004de15760008081fd5b62004deb62003ea4565b813562004df88162004744565b815262004e0782870162003dff565b86820152604062004e1a81840162003dff565b9082015260608281013562004e2f8162003ddc565b90820152608062004e4283820162004461565b9082015260a062004e5583820162004461565b9082015260c062004e6883820162003dff565b90820152845292840192810162004dc7565b5090979650505050505050565b600082601f83011262004e9957600080fd5b8135602062004eac62003fb583620044c7565b82815260059290921b8401810191818101908684111562004ecc57600080fd5b8286015b848110156200457857803562004ee68162003ddc565b835291830191830162004ed0565b600080600080600080600060e0888a03121562004f1057600080fd5b873567ffffffffffffffff8082111562004f2957600080fd5b62004f378b838c0162004d1c565b985060208a013591508082111562004f4e57600080fd5b62004f5c8b838c0162004d7e565b975060408a013591508082111562004f7357600080fd5b62004f818b838c0162004e87565b965060608a013591508082111562004f9857600080fd5b62004fa68b838c0162004e87565b955060808a013591508082111562004fbd57600080fd5b62004fcb8b838c01620044ee565b945060a08a013591508082111562004fe257600080fd5b62004ff08b838c01620044ee565b935060c08a01359150808211156200500757600080fd5b50620050168a828b01620044ee565b91505092959891949750929550565b6bffffffffffffffffffffffff81811683821601908082111562003d215762003d216200461a565b6000602082840312156200506057600080fd5b5051919050565b6bffffffffffffffffffffffff82811682821603908082111562003d215762003d216200461a565b604081526000620050a4604083018662004b4e565b8281036020840152620050b98185876200466f565b9695505050505050565b805160208083015191908110156200483b577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60209190910360031b1b16919050565b815167ffffffffffffffff81111562005123576200512362003e75565b6200513b81620051348454620047ec565b8462004841565b602080601f8311600181146200519157600084156200515a5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556200488b565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015620051e057888601518255948401946001909101908401620051bf565b50858210156200521d57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b805169ffffffffffffffffffff8116811462003e1457600080fd5b600080600080600060a086880312156200526157600080fd5b6200526c866200522d565b945060208601519350604086015192506060860151915062005291608087016200522d565b90509295509295909350565b60008154620052ac81620047ec565b808552602060018381168015620052cc5760018114620053055762005335565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b890101955062005335565b866000528260002060005b858110156200532d5781548a820186015290830190840162005310565b890184019650505b505050505092915050565b60208152600062002a8d60208301846200529d565b600082601f8301126200536757600080fd5b815160206200537a62003fb583620044c7565b82815260059290921b840181019181810190868411156200539a57600080fd5b8286015b848110156200457857805183529183019183016200539e565b600060208284031215620053ca57600080fd5b815167ffffffffffffffff80821115620053e357600080fd5b908301906101008286031215620053f957600080fd5b6200540362003ed0565b82518152602083015160208201526040830151604082015260608301516060820152608083015160808201526200543d60a08401620046e3565b60a082015260c0830151828111156200545557600080fd5b620054638782860162005355565b60c08301525060e0830151828111156200547c57600080fd5b6200548a8782860162004753565b60e08301525095945050505050565b600081518084526020808501945080840160005b8381101562004afb57815187529582019590820190600101620054ad565b60408152825160408201526020830151606082015260408301516080820152606083015160a0820152608083015160c082015273ffffffffffffffffffffffffffffffffffffffff60a08401511660e0820152600060c08401516101008081850152506200553e61014084018262005499565b905060e08501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0848303016101208501526200557c828262004145565b915050828103602084015262002d3981856200529d565b808202811582820484141762002a905762002a906200461a565b60ff818116838216019081111562002a905762002a906200461a565b60008262005600577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfe60c060405234801561001057600080fd5b506040516103ca3803806103ca83398101604081905261002f91610076565b600080546001600160a01b0319166001600160a01b039384161790559181166080521660a0526100b9565b80516001600160a01b038116811461007157600080fd5b919050565b60008060006060848603121561008b57600080fd5b6100948461005a565b92506100a26020850161005a565b91506100b06040850161005a565b90509250925092565b60805160a0516102e76100e36000396000603801526000818160c4015261011701526102e76000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806379188d161461007b578063f00e6a2a146100aa575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e808015610076573d6000f35b3d6000fd5b61008e6100893660046101c1565b6100ee565b6040805192151583526020830191909152015b60405180910390f35b60405173ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001681526020016100a1565b60008054819073ffffffffffffffffffffffffffffffffffffffff16331461011557600080fd5b7f00000000000000000000000000000000000000000000000000000000000000005a91505a61138881101561014957600080fd5b61138881039050856040820482031161016157600080fd5b50803b61016d57600080fd5b6000808551602087016000858af192505a610188908361029a565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080604083850312156101d457600080fd5b82359150602083013567ffffffffffffffff808211156101f357600080fd5b818501915085601f83011261020757600080fd5b81358181111561021957610219610192565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561025f5761025f610192565b8160405282815288602084870101111561027857600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b818103818111156102d4577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b9291505056fea164736f6c6343000813000aa164736f6c6343000813000a",
>>>>>>> 6d733192b6 (address some comments)
=======
	Bin: "0x6101206040523480156200001257600080fd5b5060405162005dd138038062005dd1833981016040819052620000359162000345565b80816001600160a01b031663ca30e6036040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200009b919062000345565b826001600160a01b031663b10b673c6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000da573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000100919062000345565b836001600160a01b0316636709d0e56040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200013f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000165919062000345565b846001600160a01b0316635425d8ac6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca919062000345565b3380600081620002215760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200025457620002548162000281565b5050506001600160a01b0393841660805291831660a052821660c052811660e0521661010052506200036c565b336001600160a01b03821603620002db5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000218565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200034257600080fd5b50565b6000602082840312156200035857600080fd5b815162000365816200032c565b9392505050565b60805160a05160c05160e05161010051615a0b620003c66000396000818161010e01526101a90152600081816103e10152612013015260006132f7015260006133db015260008181611e5501526124210152615a0b6000f3fe60806040523480156200001157600080fd5b50600436106200010c5760003560e01c806385c1b0ba11620000a5578063c8048022116200006f578063c804802214620002b7578063ce7dc5b414620002ce578063f2fde38b14620002e5578063f7d334ba14620002fc576200010c565b806385c1b0ba14620002535780638da5cb5b146200026a5780638e86139b1462000289578063948108f714620002a0576200010c565b80634ee88d3511620000e75780634ee88d3514620001ef5780636ded9eae146200020657806371791aa0146200021d57806379ba50971462000249576200010c565b806328f32f38146200015457806329c5efad146200017e578063349e8cca14620001a7575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e8080156200014d573d6000f35b3d6000fd5b005b6200016b6200016536600462003fee565b62000313565b6040519081526020015b60405180910390f35b620001956200018f366004620040d4565b6200068d565b604051620001759493929190620041fc565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200162000175565b620001526200020036600462004239565b62000931565b6200016b6200021736600462004289565b62000999565b620002346200022e366004620040d4565b620009ff565b6040516200017597969594939291906200433c565b6200015262001162565b62000152620002643660046200438e565b62001265565b60005473ffffffffffffffffffffffffffffffffffffffff16620001c9565b620001526200029a3660046200441b565b62001ed6565b62000152620002b13660046200447e565b6200225e565b62000152620002c8366004620044ad565b620024f1565b62000195620002df36600462004583565b6200293c565b62000152620002f6366004620045fa565b62002a0c565b620002346200030d366004620044ad565b62002a24565b6000805473ffffffffffffffffffffffffffffffffffffffff1633148015906200034757506200034560093362002a62565b155b156200037f576040517fd48b678b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff89163b620003ce576040517f09ee12d500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b620003d98662002a96565b9050600089307f00000000000000000000000000000000000000000000000000000000000000006040516200040e9062003d7a565b73ffffffffffffffffffffffffffffffffffffffff938416815291831660208301529091166040820152606001604051809103906000f08015801562000458573d6000803e3d6000fd5b5090506200051f826040518060e001604052806000151581526020018c63ffffffff16815260200163ffffffff801681526020018473ffffffffffffffffffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff168152602001600063ffffffff168152508a89898080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508b92508a915062002d429050565b6015805474010000000000000000000000000000000000000000900463ffffffff169060146200054f8362004649565b91906101000a81548163ffffffff021916908363ffffffff16021790555050817fbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d0128a8a604051620005c892919063ffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a2817fcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d878760405162000604929190620046b8565b60405180910390a2817f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664856040516200063e9190620046ce565b60405180910390a2817f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf485084604051620006789190620046ce565b60405180910390a25098975050505050505050565b600060606000806200069e6200312d565b600086815260046020908152604091829020825160e081018452815460ff811615158252610100810463ffffffff90811694830194909452650100000000008104841694820194909452690100000000000000000090930473ffffffffffffffffffffffffffffffffffffffff166060840152600101546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a0840152780100000000000000000000000000000000000000000000000090041660c08201525a9150600080826060015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620007b8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620007de9190620046f0565b73ffffffffffffffffffffffffffffffffffffffff166014600101600c9054906101000a900463ffffffff1663ffffffff168960405162000820919062004710565b60006040518083038160008787f1925050503d806000811462000860576040519150601f19603f3d011682016040523d82523d6000602084013e62000865565b606091505b50915091505a6200087790856200472e565b935081620008a257600060405180602001604052806000815250600796509650965050505062000928565b80806020019051810190620008b891906200479f565b909750955086620008e657600060405180602001604052806000815250600496509650965050505062000928565b601654865164010000000090910463ffffffff1610156200092457600060405180602001604052806000815250600596509650965050505062000928565b5050505b92959194509250565b6200093c8362003168565b6000838152601b602052604090206200095782848362004894565b50827f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d566483836040516200098c929190620046b8565b60405180910390a2505050565b6000620009f388888860008989604051806020016040528060008152508a8a8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506200031392505050565b98975050505050505050565b60006060600080600080600062000a156200312d565b600062000a228a6200321e565b905060006012604051806101600160405290816000820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160008201600c9054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160109054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160149054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160189054906101000a900462ffffff1662ffffff1662ffffff16815260200160008201601b9054906101000a900461ffff1661ffff1661ffff16815260200160008201601d9054906101000a900460ff1660ff1660ff16815260200160008201601e9054906101000a900460ff1615151515815260200160008201601f9054906101000a900460ff161515151581526020016001820160009054906101000a900460ff161515151581526020016001820160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152505090506000600460008d81526020019081526020016000206040518060e00160405290816000820160009054906101000a900460ff161515151581526020016000820160019054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160059054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160099054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160018201600c9054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff1681526020016001820160189054906101000a900463ffffffff1663ffffffff1663ffffffff168152505090508160e001511562000db7576000604051806020016040528060008152506009600084602001516000808263ffffffff169250995099509950995099509950995050505062001156565b604081015163ffffffff9081161462000e08576000604051806020016040528060008152506001600084602001516000808263ffffffff169250995099509950995099509950995050505062001156565b80511562000e4e576000604051806020016040528060008152506002600084602001516000808263ffffffff169250995099509950995099509950995050505062001156565b62000e5982620032d4565b602083015160165492975090955060009162000e8b918591879190640100000000900463ffffffff168a8a87620034c6565b9050806bffffffffffffffffffffffff168260a001516bffffffffffffffffffffffff16101562000ef5576000604051806020016040528060008152506006600085602001516000808263ffffffff1692509a509a509a509a509a509a509a505050505062001156565b600062000f048e868f62003517565b90505a9850600080846060015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000f5c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000f829190620046f0565b73ffffffffffffffffffffffffffffffffffffffff166014600101600c9054906101000a900463ffffffff1663ffffffff168460405162000fc4919062004710565b60006040518083038160008787f1925050503d806000811462001004576040519150601f19603f3d011682016040523d82523d6000602084013e62001009565b606091505b50915091505a6200101b908c6200472e565b9a50816200109b5760165481516801000000000000000090910463ffffffff1610156200107857505060408051602080820190925260008082529490910151939c509a50600899505063ffffffff90911695506200115692505050565b602090940151939b5060039a505063ffffffff9092169650620011569350505050565b80806020019051810190620010b191906200479f565b909e509c508d620010f257505060408051602080820190925260008082529490910151939c509a50600499505063ffffffff90911695506200115692505050565b6016548d5164010000000090910463ffffffff1610156200114357505060408051602080820190925260008082529490910151939c509a50600599505063ffffffff90911695506200115692505050565b505050506020015163ffffffff16945050505b92959891949750929550565b60015473ffffffffffffffffffffffffffffffffffffffff163314620011e9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600173ffffffffffffffffffffffffffffffffffffffff82166000908152601a602052604090205460ff166003811115620012a457620012a462004191565b14158015620012f05750600373ffffffffffffffffffffffffffffffffffffffff82166000908152601a602052604090205460ff166003811115620012ed57620012ed62004191565b14155b1562001328576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6014546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1662001388576040517fd12d7d8d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000829003620013c4576040517f2c2fc94100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160e081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c081018290526000808567ffffffffffffffff8111156200141b576200141b62003e75565b60405190808252806020026020018201604052801562001445578160200160208202803683370190505b50905060008667ffffffffffffffff81111562001466576200146662003e75565b604051908082528060200260200182016040528015620014ed57816020015b6040805160e08101825260008082526020808301829052928201819052606082018190526080820181905260a0820181905260c082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181620014855790505b50905060008767ffffffffffffffff8111156200150e576200150e62003e75565b6040519080825280602002602001820160405280156200154357816020015b60608152602001906001900390816200152d5790505b50905060008867ffffffffffffffff81111562001564576200156462003e75565b6040519080825280602002602001820160405280156200159957816020015b6060815260200190600190039081620015835790505b50905060008967ffffffffffffffff811115620015ba57620015ba62003e75565b604051908082528060200260200182016040528015620015ef57816020015b6060815260200190600190039081620015d95790505b50905060005b8a81101562001bd3578b8b82818110620016135762001613620049bc565b60209081029290920135600081815260048452604090819020815160e081018352815460ff811615158252610100810463ffffffff90811697830197909752650100000000008104871693820193909352690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff166060830152600101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490931660c08401529a50909850620016f290508962003168565b60608801516040517f1a5da6c800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8c8116600483015290911690631a5da6c890602401600060405180830381600087803b1580156200176257600080fd5b505af115801562001777573d6000803e3d6000fd5b5050505087858281518110620017915762001791620049bc565b6020026020010181905250600560008a815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16868281518110620017e557620017e5620049bc565b73ffffffffffffffffffffffffffffffffffffffff90921660209283029190910182015260008a815260079091526040902080546200182490620047ec565b80601f01602080910402602001604051908101604052809291908181526020018280546200185290620047ec565b8015620018a35780601f106200187757610100808354040283529160200191620018a3565b820191906000526020600020905b8154815290600101906020018083116200188557829003601f168201915b5050505050848281518110620018bd57620018bd620049bc565b6020026020010181905250601b60008a81526020019081526020016000208054620018e890620047ec565b80601f01602080910402602001604051908101604052809291908181526020018280546200191690620047ec565b8015620019675780601f106200193b5761010080835404028352916020019162001967565b820191906000526020600020905b8154815290600101906020018083116200194957829003601f168201915b5050505050838281518110620019815762001981620049bc565b6020026020010181905250601c60008a81526020019081526020016000208054620019ac90620047ec565b80601f0160208091040260200160405190810160405280929190818152602001828054620019da90620047ec565b801562001a2b5780601f10620019ff5761010080835404028352916020019162001a2b565b820191906000526020600020905b81548152906001019060200180831162001a0d57829003601f168201915b505050505082828151811062001a455762001a45620049bc565b60200260200101819052508760a001516bffffffffffffffffffffffff168762001a709190620049eb565b60008a815260046020908152604080832080547fffffff000000000000000000000000000000000000000000000000000000000016815560010180547fffffffff000000000000000000000000000000000000000000000000000000001690556007909152812091985062001ae6919062003d88565b6000898152601b6020526040812062001aff9162003d88565b6000898152601c6020526040812062001b189162003d88565b600089815260066020526040902080547fffffffffffffffffffffffff000000000000000000000000000000000000000016905562001b5960028a62003739565b5060a0880151604080516bffffffffffffffffffffffff909216825273ffffffffffffffffffffffffffffffffffffffff8c1660208301528a917fb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff910160405180910390a28062001bca8162004a01565b915050620015f5565b508560195462001be491906200472e565b60195560008b8b868167ffffffffffffffff81111562001c085762001c0862003e75565b60405190808252806020026020018201604052801562001c32578160200160208202803683370190505b508988888860405160200162001c5098979695949392919062004ba7565b60405160208183030381529060405290508973ffffffffffffffffffffffffffffffffffffffff16638e86139b6014600001600c9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663c71249ab60038e73ffffffffffffffffffffffffffffffffffffffff1663aab9edd66040518163ffffffff1660e01b8152600401602060405180830381865afa15801562001d0c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001d32919062004c77565b866040518463ffffffff1660e01b815260040162001d539392919062004c9c565b600060405180830381865afa15801562001d71573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405262001db9919081019062004cc3565b6040518263ffffffff1660e01b815260040162001dd79190620046ce565b600060405180830381600087803b15801562001df257600080fd5b505af115801562001e07573d6000803e3d6000fd5b50506040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8d81166004830152602482018b90527f000000000000000000000000000000000000000000000000000000000000000016925063a9059cbb91506044016020604051808303816000875af115801562001ea1573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001ec7919062004cfc565b50505050505050505050505050565b6002336000908152601a602052604090205460ff16600381111562001eff5762001eff62004191565b1415801562001f3557506003336000908152601a602052604090205460ff16600381111562001f325762001f3262004191565b14155b1562001f6d576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080808080808062001f83888a018a62004ef4565b965096509650965096509650965060005b87518110156200225257600073ffffffffffffffffffffffffffffffffffffffff1687828151811062001fcb5762001fcb620049bc565b60200260200101516060015173ffffffffffffffffffffffffffffffffffffffff1603620020df57858181518110620020085762002008620049bc565b6020026020010151307f0000000000000000000000000000000000000000000000000000000000000000604051620020409062003d7a565b73ffffffffffffffffffffffffffffffffffffffff938416815291831660208301529091166040820152606001604051809103906000f0801580156200208a573d6000803e3d6000fd5b50878281518110620020a057620020a0620049bc565b60200260200101516060019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250505b62002197888281518110620020f857620020f8620049bc565b6020026020010151888381518110620021155762002115620049bc565b6020026020010151878481518110620021325762002132620049bc565b60200260200101518785815181106200214f576200214f620049bc565b60200260200101518786815181106200216c576200216c620049bc565b6020026020010151878781518110620021895762002189620049bc565b602002602001015162002d42565b878181518110620021ac57620021ac620049bc565b60200260200101517f74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71888381518110620021ea57620021ea620049bc565b602002602001015160a0015133604051620022359291906bffffffffffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a280620022498162004a01565b91505062001f94565b50505050505050505050565b600082815260046020908152604091829020825160e081018452815460ff81161515825263ffffffff6101008204811694830194909452650100000000008104841694820185905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004821660c082015291146200235c576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818160a001516200236e919062005025565b600084815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff93841602179055601954620023d691841690620049eb565b6019556040517f23b872dd0000000000000000000000000000000000000000000000000000000081523360048201523060248201526bffffffffffffffffffffffff831660448201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906323b872dd906064016020604051808303816000875af115801562002480573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620024a6919062004cfc565b506040516bffffffffffffffffffffffff83168152339084907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a3505050565b6000818152600460209081526040808320815160e081018352815460ff81161515825263ffffffff6101008204811695830195909552650100000000008104851693820184905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004831660c08201529291141590620025da60005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161490506000601260010160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200267d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620026a391906200504d565b9050828015620026c75750818015620026c5575080846040015163ffffffff16115b155b15620026ff576040517ffbc0357800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8115801562002732575060008581526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314155b156200276a576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8162002780576200277d603282620049eb565b90505b6000858152600460205260409020805463ffffffff80841665010000000000027fffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffffff90921691909117909155620027dc9060029087906200373916565b5060145460808501516bffffffffffffffffffffffff91821691600091168211156200284557608086015162002813908362005067565b90508560a001516bffffffffffffffffffffffff16816bffffffffffffffffffffffff16111562002845575060a08501515b808660a0015162002857919062005067565b600088815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff93841602179055601554620028bf9183911662005025565b601580547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff9290921691909117905560405167ffffffffffffffff84169088907f91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f79118190600090a350505050505050565b600060606000806200294d6200312d565b6000634b56a42e60e01b8888886040516024016200296e939291906200508f565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091529050620029f989826200068d565b929c919b50995090975095505050505050565b62002a1662003747565b62002a2181620037ca565b50565b60006060600080600080600062002a4b8860405180602001604052806000815250620009ff565b959e949d50929b5090995097509550909350915050565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415155b90505b92915050565b6000806000601260010160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905060008173ffffffffffffffffffffffffffffffffffffffff166385df51fd60018473ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa15801562002b2f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062002b5591906200504d565b62002b6191906200472e565b6040518263ffffffff1660e01b815260040162002b8091815260200190565b602060405180830381865afa15801562002b9e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062002bc491906200504d565b601554604080516020810193909352309083015274010000000000000000000000000000000000000000900463ffffffff166060820152608001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201209083015201604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052905060045b600f81101562002cd0578382828151811062002c8c5762002c8c620049bc565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508062002cc78162004a01565b91505062002c6c565b5084600181111562002ce65762002ce662004191565b60f81b81600f8151811062002cff5762002cff620049bc565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535062002d3981620050c3565b95945050505050565b6012547e01000000000000000000000000000000000000000000000000000000000000900460ff161562002da2576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601654835163ffffffff909116101562002de8576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108fc856020015163ffffffff16108062002e265750601554602086015163ffffffff70010000000000000000000000000000000090920482169116115b1562002e5e576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000868152600460205260409020546901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff161562002ec8576040517f6e3b930b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000868152600460209081526040808320885181548a8501518b85015160608d01517fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000009093169315157fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000ff169390931761010063ffffffff92831602177fffffff000000000000000000000000000000000000000000000000ffffffffff1665010000000000938216939093027fffffff0000000000000000000000000000000000000000ffffffffffffffffff1692909217690100000000000000000073ffffffffffffffffffffffffffffffffffffffff9283160217835560808b01516001909301805460a08d015160c08e01516bffffffffffffffffffffffff9687167fffffffffffffffff000000000000000000000000000000000000000000000000909316929092176c010000000000000000000000009690911695909502949094177fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff1678010000000000000000000000000000000000000000000000009490931693909302919091179091556005835281842080547fffffffffffffffffffffffff00000000000000000000000000000000000000001691891691909117905560079091529020620030bb848262005106565b508460a001516bffffffffffffffffffffffff16601954620030de9190620049eb565b6019556000868152601b60205260409020620030fb838262005106565b506000868152601c6020526040902062003116828262005106565b5062003124600287620038c1565b50505050505050565b321562003166576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60008181526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314620031c6576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526004602052604090205465010000000000900463ffffffff9081161462002a21576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000818160045b600f811015620032b3577fff000000000000000000000000000000000000000000000000000000000000008216838260208110620032675762003267620049bc565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916146200329e57506000949350505050565b80620032aa8162004a01565b91505062003225565b5081600f1a6001811115620032cc57620032cc62004191565b949350505050565b6000806000836080015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa15801562003361573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003387919062005248565b50945090925050506000811315806200339f57508142105b80620033c45750828015620033c45750620033bb82426200472e565b8463ffffffff16105b15620033d5576017549550620033d9565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa15801562003445573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200346b919062005248565b50945090925050506000811315806200348357508142105b80620034a85750828015620034a857506200349f82426200472e565b8463ffffffff16105b15620034b9576018549450620034bd565b8094505b50505050915091565b600080620034da88878b60c00151620038cf565b9050600080620034f78b8a63ffffffff16858a8a60018b62003965565b909250905062003508818362005025565b9b9a5050505050505050505050565b6060600083600181111562003530576200353062004191565b03620035fd576000848152600760205260409081902090517f6e04ff0d0000000000000000000000000000000000000000000000000000000091620035789160240162005340565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152905062003732565b600183600181111562003614576200361462004191565b036200370057600082806020019051810190620036329190620053b7565b6000868152600760205260409081902090519192507f40691db4000000000000000000000000000000000000000000000000000000009162003679918491602401620054cb565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091529150620037329050565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b9392505050565b600062002a8d838362003c1d565b60005473ffffffffffffffffffffffffffffffffffffffff16331462003166576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401620011e0565b3373ffffffffffffffffffffffffffffffffffffffff8216036200384b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620011e0565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600062002a8d838362003d28565b60008080856001811115620038e857620038e862004191565b03620038f9575062015f906200391c565b600185600181111562003910576200391062004191565b036200370057506201ae145b6200392f63ffffffff8516601462005593565b6200393c846001620055ad565b6200394d9060ff16611d4c62005593565b620039599083620049eb565b62002d399190620049eb565b60008060008960a0015161ffff168762003980919062005593565b90508380156200398f5750803a105b156200399857503a5b6000841562003a20578a610140015173ffffffffffffffffffffffffffffffffffffffff166318b8f6136040518163ffffffff1660e01b8152600401602060405180830381865afa158015620039f2573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003a1891906200504d565b905062003ae1565b6101408b01516016546040517f1254414000000000000000000000000000000000000000000000000000000000815264010000000090910463ffffffff16600482015273ffffffffffffffffffffffffffffffffffffffff90911690631254414090602401602060405180830381865afa15801562003aa3573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003ac991906200504d565b8b60a0015161ffff1662003ade919062005593565b90505b62003af161ffff871682620055c9565b90506000878262003b038c8e620049eb565b62003b0f908662005593565b62003b1b9190620049eb565b62003b2f90670de0b6b3a764000062005593565b62003b3b9190620055c9565b905060008c6040015163ffffffff1664e8d4a5100062003b5c919062005593565b898e6020015163ffffffff16858f8862003b77919062005593565b62003b839190620049eb565b62003b9390633b9aca0062005593565b62003b9f919062005593565b62003bab9190620055c9565b62003bb79190620049eb565b90506b033b2e3c9fd0803ce800000062003bd28284620049eb565b111562003c0b576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b6000818152600183016020526040812054801562003d1657600062003c446001836200472e565b855490915060009062003c5a906001906200472e565b905081811462003cc657600086600001828154811062003c7e5762003c7e620049bc565b906000526020600020015490508087600001848154811062003ca45762003ca4620049bc565b6000918252602080832090910192909255918252600188019052604090208390555b855486908062003cda5762003cda62005605565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505062002a90565b600091505062002a90565b5092915050565b600081815260018301602052604081205462003d715750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915562002a90565b50600062002a90565b6103ca806200563583390190565b50805462003d9690620047ec565b6000825580601f1062003da7575050565b601f01602090049060005260206000209081019062002a2191905b8082111562003dd8576000815560010162003dc2565b5090565b73ffffffffffffffffffffffffffffffffffffffff8116811462002a2157600080fd5b803563ffffffff8116811462003e1457600080fd5b919050565b80356002811062003e1457600080fd5b60008083601f84011262003e3c57600080fd5b50813567ffffffffffffffff81111562003e5557600080fd5b60208301915083602082850101111562003e6e57600080fd5b9250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160e0810167ffffffffffffffff8111828210171562003eca5762003eca62003e75565b60405290565b604051610100810167ffffffffffffffff8111828210171562003eca5762003eca62003e75565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171562003f415762003f4162003e75565b604052919050565b600067ffffffffffffffff82111562003f665762003f6662003e75565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f83011262003fa457600080fd5b813562003fbb62003fb58262003f49565b62003ef7565b81815284602083860101111562003fd157600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060008060e0898b0312156200400b57600080fd5b8835620040188162003ddc565b97506200402860208a0162003dff565b965060408901356200403a8162003ddc565b95506200404a60608a0162003e19565b9450608089013567ffffffffffffffff808211156200406857600080fd5b620040768c838d0162003e29565b909650945060a08b01359150808211156200409057600080fd5b6200409e8c838d0162003f92565b935060c08b0135915080821115620040b557600080fd5b50620040c48b828c0162003f92565b9150509295985092959890939650565b60008060408385031215620040e857600080fd5b82359150602083013567ffffffffffffffff8111156200410757600080fd5b620041158582860162003f92565b9150509250929050565b60005b838110156200413c57818101518382015260200162004122565b50506000910152565b600081518084526200415f8160208601602086016200411f565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600a8110620041f8577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b9052565b841515815260806020820152600062004219608083018662004145565b90506200422a6040830185620041c0565b82606083015295945050505050565b6000806000604084860312156200424f57600080fd5b83359250602084013567ffffffffffffffff8111156200426e57600080fd5b6200427c8682870162003e29565b9497909650939450505050565b600080600080600080600060a0888a031215620042a557600080fd5b8735620042b28162003ddc565b9650620042c26020890162003dff565b95506040880135620042d48162003ddc565b9450606088013567ffffffffffffffff80821115620042f257600080fd5b620043008b838c0162003e29565b909650945060808a01359150808211156200431a57600080fd5b50620043298a828b0162003e29565b989b979a50959850939692959293505050565b871515815260e0602082015260006200435960e083018962004145565b90506200436a6040830188620041c0565b8560608301528460808301528360a08301528260c083015298975050505050505050565b600080600060408486031215620043a457600080fd5b833567ffffffffffffffff80821115620043bd57600080fd5b818601915086601f830112620043d257600080fd5b813581811115620043e257600080fd5b8760208260051b8501011115620043f857600080fd5b60209283019550935050840135620044108162003ddc565b809150509250925092565b600080602083850312156200442f57600080fd5b823567ffffffffffffffff8111156200444757600080fd5b620044558582860162003e29565b90969095509350505050565b80356bffffffffffffffffffffffff8116811462003e1457600080fd5b600080604083850312156200449257600080fd5b82359150620044a46020840162004461565b90509250929050565b600060208284031215620044c057600080fd5b5035919050565b600067ffffffffffffffff821115620044e457620044e462003e75565b5060051b60200190565b600082601f8301126200450057600080fd5b813560206200451362003fb583620044c7565b82815260059290921b840181019181810190868411156200453357600080fd5b8286015b848110156200457857803567ffffffffffffffff811115620045595760008081fd5b620045698986838b010162003f92565b84525091830191830162004537565b509695505050505050565b600080600080606085870312156200459a57600080fd5b84359350602085013567ffffffffffffffff80821115620045ba57600080fd5b620045c888838901620044ee565b94506040870135915080821115620045df57600080fd5b50620045ee8782880162003e29565b95989497509550505050565b6000602082840312156200460d57600080fd5b8135620037328162003ddc565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600063ffffffff8083168181036200466557620046656200461a565b6001019392505050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b602081526000620032cc6020830184866200466f565b60208152600062002a8d602083018462004145565b805162003e148162003ddc565b6000602082840312156200470357600080fd5b8151620037328162003ddc565b60008251620047248184602087016200411f565b9190910192915050565b8181038181111562002a905762002a906200461a565b801515811462002a2157600080fd5b600082601f8301126200476557600080fd5b81516200477662003fb58262003f49565b8181528460208386010111156200478c57600080fd5b620032cc8260208301602087016200411f565b60008060408385031215620047b357600080fd5b8251620047c08162004744565b602084015190925067ffffffffffffffff811115620047de57600080fd5b620041158582860162004753565b600181811c908216806200480157607f821691505b6020821081036200483b577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f8211156200488f57600081815260208120601f850160051c810160208610156200486a5750805b601f850160051c820191505b818110156200488b5782815560010162004876565b5050505b505050565b67ffffffffffffffff831115620048af57620048af62003e75565b620048c783620048c08354620047ec565b8362004841565b6000601f8411600181146200491c5760008515620048e55750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355620049b5565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b828110156200496d57868501358255602094850194600190920191016200494b565b5086821015620049a9577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b8082018082111562002a905762002a906200461a565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820362004a355762004a356200461a565b5060010190565b600081518084526020808501945080840160005b8381101562004afb5781518051151588528381015163ffffffff908116858a01526040808301519091169089015260608082015173ffffffffffffffffffffffffffffffffffffffff16908901526080808201516bffffffffffffffffffffffff169089015260a08082015162004ad6828b01826bffffffffffffffffffffffff169052565b505060c09081015163ffffffff169088015260e0909601959082019060010162004a50565b509495945050505050565b600081518084526020808501945080840160005b8381101562004afb57815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010162004b1a565b600081518084526020808501808196508360051b8101915082860160005b8581101562004b9a57828403895262004b8784835162004145565b9885019893509084019060010162004b6c565b5091979650505050505050565b60e081528760e082015260006101007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8a111562004be457600080fd5b8960051b808c8386013783018381038201602085015262004c088282018b62004a3c565b915050828103604084015262004c1f818962004b06565b9050828103606084015262004c35818862004b06565b9050828103608084015262004c4b818762004b4e565b905082810360a084015262004c61818662004b4e565b905082810360c084015262003508818562004b4e565b60006020828403121562004c8a57600080fd5b815160ff811681146200373257600080fd5b60ff8416815260ff8316602082015260606040820152600062002d39606083018462004145565b60006020828403121562004cd657600080fd5b815167ffffffffffffffff81111562004cee57600080fd5b620032cc8482850162004753565b60006020828403121562004d0f57600080fd5b8151620037328162004744565b600082601f83011262004d2e57600080fd5b8135602062004d4162003fb583620044c7565b82815260059290921b8401810191818101908684111562004d6157600080fd5b8286015b8481101562004578578035835291830191830162004d65565b600082601f83011262004d9057600080fd5b8135602062004da362003fb583620044c7565b82815260e0928302850182019282820191908785111562004dc357600080fd5b8387015b8581101562004e7a5781818a03121562004de15760008081fd5b62004deb62003ea4565b813562004df88162004744565b815262004e0782870162003dff565b86820152604062004e1a81840162003dff565b9082015260608281013562004e2f8162003ddc565b90820152608062004e4283820162004461565b9082015260a062004e5583820162004461565b9082015260c062004e6883820162003dff565b90820152845292840192810162004dc7565b5090979650505050505050565b600082601f83011262004e9957600080fd5b8135602062004eac62003fb583620044c7565b82815260059290921b8401810191818101908684111562004ecc57600080fd5b8286015b848110156200457857803562004ee68162003ddc565b835291830191830162004ed0565b600080600080600080600060e0888a03121562004f1057600080fd5b873567ffffffffffffffff8082111562004f2957600080fd5b62004f378b838c0162004d1c565b985060208a013591508082111562004f4e57600080fd5b62004f5c8b838c0162004d7e565b975060408a013591508082111562004f7357600080fd5b62004f818b838c0162004e87565b965060608a013591508082111562004f9857600080fd5b62004fa68b838c0162004e87565b955060808a013591508082111562004fbd57600080fd5b62004fcb8b838c01620044ee565b945060a08a013591508082111562004fe257600080fd5b62004ff08b838c01620044ee565b935060c08a01359150808211156200500757600080fd5b50620050168a828b01620044ee565b91505092959891949750929550565b6bffffffffffffffffffffffff81811683821601908082111562003d215762003d216200461a565b6000602082840312156200506057600080fd5b5051919050565b6bffffffffffffffffffffffff82811682821603908082111562003d215762003d216200461a565b604081526000620050a4604083018662004b4e565b8281036020840152620050b98185876200466f565b9695505050505050565b805160208083015191908110156200483b577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60209190910360031b1b16919050565b815167ffffffffffffffff81111562005123576200512362003e75565b6200513b81620051348454620047ec565b8462004841565b602080601f8311600181146200519157600084156200515a5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556200488b565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015620051e057888601518255948401946001909101908401620051bf565b50858210156200521d57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b805169ffffffffffffffffffff8116811462003e1457600080fd5b600080600080600060a086880312156200526157600080fd5b6200526c866200522d565b945060208601519350604086015192506060860151915062005291608087016200522d565b90509295509295909350565b60008154620052ac81620047ec565b808552602060018381168015620052cc5760018114620053055762005335565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b890101955062005335565b866000528260002060005b858110156200532d5781548a820186015290830190840162005310565b890184019650505b505050505092915050565b60208152600062002a8d60208301846200529d565b600082601f8301126200536757600080fd5b815160206200537a62003fb583620044c7565b82815260059290921b840181019181810190868411156200539a57600080fd5b8286015b848110156200457857805183529183019183016200539e565b600060208284031215620053ca57600080fd5b815167ffffffffffffffff80821115620053e357600080fd5b908301906101008286031215620053f957600080fd5b6200540362003ed0565b82518152602083015160208201526040830151604082015260608301516060820152608083015160808201526200543d60a08401620046e3565b60a082015260c0830151828111156200545557600080fd5b620054638782860162005355565b60c08301525060e0830151828111156200547c57600080fd5b6200548a8782860162004753565b60e08301525095945050505050565b600081518084526020808501945080840160005b8381101562004afb57815187529582019590820190600101620054ad565b60408152825160408201526020830151606082015260408301516080820152606083015160a0820152608083015160c082015273ffffffffffffffffffffffffffffffffffffffff60a08401511660e0820152600060c08401516101008081850152506200553e61014084018262005499565b905060e08501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0848303016101208501526200557c828262004145565b915050828103602084015262002d3981856200529d565b808202811582820484141762002a905762002a906200461a565b60ff818116838216019081111562002a905762002a906200461a565b60008262005600577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfe60c060405234801561001057600080fd5b506040516103ca3803806103ca83398101604081905261002f91610076565b600080546001600160a01b0319166001600160a01b039384161790559181166080521660a0526100b9565b80516001600160a01b038116811461007157600080fd5b919050565b60008060006060848603121561008b57600080fd5b6100948461005a565b92506100a26020850161005a565b91506100b06040850161005a565b90509250925092565b60805160a0516102e76100e36000396000603801526000818160c4015261011701526102e76000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806379188d161461007b578063f00e6a2a146100aa575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e808015610076573d6000f35b3d6000fd5b61008e6100893660046101c1565b6100ee565b6040805192151583526020830191909152015b60405180910390f35b60405173ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001681526020016100a1565b60008054819073ffffffffffffffffffffffffffffffffffffffff16331461011557600080fd5b7f00000000000000000000000000000000000000000000000000000000000000005a91505a61138881101561014957600080fd5b61138881039050856040820482031161016157600080fd5b50803b61016d57600080fd5b6000808551602087016000858af192505a610188908361029a565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080604083850312156101d457600080fd5b82359150602083013567ffffffffffffffff808211156101f357600080fd5b818501915085601f83011261020757600080fd5b81358181111561021957610219610192565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561025f5761025f610192565b8160405282815288602084870101111561027857600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b818103818111156102d4577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b9291505056fea164736f6c6343000813000aa164736f6c6343000813000a",
>>>>>>> 9507a3be49 (adjust gas overheads)
=======
	Bin: "0x6101206040523480156200001257600080fd5b5060405162005dd138038062005dd1833981016040819052620000359162000345565b80816001600160a01b031663ca30e6036040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200009b919062000345565b826001600160a01b031663b10b673c6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000da573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000100919062000345565b836001600160a01b0316636709d0e56040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200013f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000165919062000345565b846001600160a01b0316635425d8ac6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca919062000345565b3380600081620002215760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200025457620002548162000281565b5050506001600160a01b0393841660805291831660a052821660c052811660e0521661010052506200036c565b336001600160a01b03821603620002db5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000218565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200034257600080fd5b50565b6000602082840312156200035857600080fd5b815162000365816200032c565b9392505050565b60805160a05160c05160e05161010051615a0b620003c66000396000818161010e01526101a90152600081816103e10152612013015260006132f7015260006133db015260008181611e5501526124210152615a0b6000f3fe60806040523480156200001157600080fd5b50600436106200010c5760003560e01c806385c1b0ba11620000a5578063c8048022116200006f578063c804802214620002b7578063ce7dc5b414620002ce578063f2fde38b14620002e5578063f7d334ba14620002fc576200010c565b806385c1b0ba14620002535780638da5cb5b146200026a5780638e86139b1462000289578063948108f714620002a0576200010c565b80634ee88d3511620000e75780634ee88d3514620001ef5780636ded9eae146200020657806371791aa0146200021d57806379ba50971462000249576200010c565b806328f32f38146200015457806329c5efad146200017e578063349e8cca14620001a7575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e8080156200014d573d6000f35b3d6000fd5b005b6200016b6200016536600462003fee565b62000313565b6040519081526020015b60405180910390f35b620001956200018f366004620040d4565b6200068d565b604051620001759493929190620041fc565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200162000175565b620001526200020036600462004239565b62000931565b6200016b6200021736600462004289565b62000999565b620002346200022e366004620040d4565b620009ff565b6040516200017597969594939291906200433c565b6200015262001162565b62000152620002643660046200438e565b62001265565b60005473ffffffffffffffffffffffffffffffffffffffff16620001c9565b620001526200029a3660046200441b565b62001ed6565b62000152620002b13660046200447e565b6200225e565b62000152620002c8366004620044ad565b620024f1565b62000195620002df36600462004583565b6200293c565b62000152620002f6366004620045fa565b62002a0c565b620002346200030d366004620044ad565b62002a24565b6000805473ffffffffffffffffffffffffffffffffffffffff1633148015906200034757506200034560093362002a62565b155b156200037f576040517fd48b678b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff89163b620003ce576040517f09ee12d500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b620003d98662002a96565b9050600089307f00000000000000000000000000000000000000000000000000000000000000006040516200040e9062003d7a565b73ffffffffffffffffffffffffffffffffffffffff938416815291831660208301529091166040820152606001604051809103906000f08015801562000458573d6000803e3d6000fd5b5090506200051f826040518060e001604052806000151581526020018c63ffffffff16815260200163ffffffff801681526020018473ffffffffffffffffffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff168152602001600063ffffffff168152508a89898080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508b92508a915062002d429050565b6015805474010000000000000000000000000000000000000000900463ffffffff169060146200054f8362004649565b91906101000a81548163ffffffff021916908363ffffffff16021790555050817fbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d0128a8a604051620005c892919063ffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a2817fcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d878760405162000604929190620046b8565b60405180910390a2817f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664856040516200063e9190620046ce565b60405180910390a2817f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf485084604051620006789190620046ce565b60405180910390a25098975050505050505050565b600060606000806200069e6200312d565b600086815260046020908152604091829020825160e081018452815460ff811615158252610100810463ffffffff90811694830194909452650100000000008104841694820194909452690100000000000000000090930473ffffffffffffffffffffffffffffffffffffffff166060840152600101546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a0840152780100000000000000000000000000000000000000000000000090041660c08201525a9150600080826060015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620007b8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620007de9190620046f0565b73ffffffffffffffffffffffffffffffffffffffff166014600101600c9054906101000a900463ffffffff1663ffffffff168960405162000820919062004710565b60006040518083038160008787f1925050503d806000811462000860576040519150601f19603f3d011682016040523d82523d6000602084013e62000865565b606091505b50915091505a6200087790856200472e565b935081620008a257600060405180602001604052806000815250600796509650965050505062000928565b80806020019051810190620008b891906200479f565b909750955086620008e657600060405180602001604052806000815250600496509650965050505062000928565b601654865164010000000090910463ffffffff1610156200092457600060405180602001604052806000815250600596509650965050505062000928565b5050505b92959194509250565b6200093c8362003168565b6000838152601b602052604090206200095782848362004894565b50827f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d566483836040516200098c929190620046b8565b60405180910390a2505050565b6000620009f388888860008989604051806020016040528060008152508a8a8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506200031392505050565b98975050505050505050565b60006060600080600080600062000a156200312d565b600062000a228a6200321e565b905060006012604051806101600160405290816000820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160008201600c9054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160109054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160149054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160189054906101000a900462ffffff1662ffffff1662ffffff16815260200160008201601b9054906101000a900461ffff1661ffff1661ffff16815260200160008201601d9054906101000a900460ff1660ff1660ff16815260200160008201601e9054906101000a900460ff1615151515815260200160008201601f9054906101000a900460ff161515151581526020016001820160009054906101000a900460ff161515151581526020016001820160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152505090506000600460008d81526020019081526020016000206040518060e00160405290816000820160009054906101000a900460ff161515151581526020016000820160019054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160059054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160099054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160018201600c9054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff1681526020016001820160189054906101000a900463ffffffff1663ffffffff1663ffffffff168152505090508160e001511562000db7576000604051806020016040528060008152506009600084602001516000808263ffffffff169250995099509950995099509950995050505062001156565b604081015163ffffffff9081161462000e08576000604051806020016040528060008152506001600084602001516000808263ffffffff169250995099509950995099509950995050505062001156565b80511562000e4e576000604051806020016040528060008152506002600084602001516000808263ffffffff169250995099509950995099509950995050505062001156565b62000e5982620032d4565b602083015160165492975090955060009162000e8b918591879190640100000000900463ffffffff168a8a87620034c6565b9050806bffffffffffffffffffffffff168260a001516bffffffffffffffffffffffff16101562000ef5576000604051806020016040528060008152506006600085602001516000808263ffffffff1692509a509a509a509a509a509a509a505050505062001156565b600062000f048e868f62003517565b90505a9850600080846060015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000f5c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000f829190620046f0565b73ffffffffffffffffffffffffffffffffffffffff166014600101600c9054906101000a900463ffffffff1663ffffffff168460405162000fc4919062004710565b60006040518083038160008787f1925050503d806000811462001004576040519150601f19603f3d011682016040523d82523d6000602084013e62001009565b606091505b50915091505a6200101b908c6200472e565b9a50816200109b5760165481516801000000000000000090910463ffffffff1610156200107857505060408051602080820190925260008082529490910151939c509a50600899505063ffffffff90911695506200115692505050565b602090940151939b5060039a505063ffffffff9092169650620011569350505050565b80806020019051810190620010b191906200479f565b909e509c508d620010f257505060408051602080820190925260008082529490910151939c509a50600499505063ffffffff90911695506200115692505050565b6016548d5164010000000090910463ffffffff1610156200114357505060408051602080820190925260008082529490910151939c509a50600599505063ffffffff90911695506200115692505050565b505050506020015163ffffffff16945050505b92959891949750929550565b60015473ffffffffffffffffffffffffffffffffffffffff163314620011e9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600173ffffffffffffffffffffffffffffffffffffffff82166000908152601a602052604090205460ff166003811115620012a457620012a462004191565b14158015620012f05750600373ffffffffffffffffffffffffffffffffffffffff82166000908152601a602052604090205460ff166003811115620012ed57620012ed62004191565b14155b1562001328576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6014546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1662001388576040517fd12d7d8d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000829003620013c4576040517f2c2fc94100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160e081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c081018290526000808567ffffffffffffffff8111156200141b576200141b62003e75565b60405190808252806020026020018201604052801562001445578160200160208202803683370190505b50905060008667ffffffffffffffff81111562001466576200146662003e75565b604051908082528060200260200182016040528015620014ed57816020015b6040805160e08101825260008082526020808301829052928201819052606082018190526080820181905260a0820181905260c082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181620014855790505b50905060008767ffffffffffffffff8111156200150e576200150e62003e75565b6040519080825280602002602001820160405280156200154357816020015b60608152602001906001900390816200152d5790505b50905060008867ffffffffffffffff81111562001564576200156462003e75565b6040519080825280602002602001820160405280156200159957816020015b6060815260200190600190039081620015835790505b50905060008967ffffffffffffffff811115620015ba57620015ba62003e75565b604051908082528060200260200182016040528015620015ef57816020015b6060815260200190600190039081620015d95790505b50905060005b8a81101562001bd3578b8b82818110620016135762001613620049bc565b60209081029290920135600081815260048452604090819020815160e081018352815460ff811615158252610100810463ffffffff90811697830197909752650100000000008104871693820193909352690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff166060830152600101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490931660c08401529a50909850620016f290508962003168565b60608801516040517f1a5da6c800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8c8116600483015290911690631a5da6c890602401600060405180830381600087803b1580156200176257600080fd5b505af115801562001777573d6000803e3d6000fd5b5050505087858281518110620017915762001791620049bc565b6020026020010181905250600560008a815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16868281518110620017e557620017e5620049bc565b73ffffffffffffffffffffffffffffffffffffffff90921660209283029190910182015260008a815260079091526040902080546200182490620047ec565b80601f01602080910402602001604051908101604052809291908181526020018280546200185290620047ec565b8015620018a35780601f106200187757610100808354040283529160200191620018a3565b820191906000526020600020905b8154815290600101906020018083116200188557829003601f168201915b5050505050848281518110620018bd57620018bd620049bc565b6020026020010181905250601b60008a81526020019081526020016000208054620018e890620047ec565b80601f01602080910402602001604051908101604052809291908181526020018280546200191690620047ec565b8015620019675780601f106200193b5761010080835404028352916020019162001967565b820191906000526020600020905b8154815290600101906020018083116200194957829003601f168201915b5050505050838281518110620019815762001981620049bc565b6020026020010181905250601c60008a81526020019081526020016000208054620019ac90620047ec565b80601f0160208091040260200160405190810160405280929190818152602001828054620019da90620047ec565b801562001a2b5780601f10620019ff5761010080835404028352916020019162001a2b565b820191906000526020600020905b81548152906001019060200180831162001a0d57829003601f168201915b505050505082828151811062001a455762001a45620049bc565b60200260200101819052508760a001516bffffffffffffffffffffffff168762001a709190620049eb565b60008a815260046020908152604080832080547fffffff000000000000000000000000000000000000000000000000000000000016815560010180547fffffffff000000000000000000000000000000000000000000000000000000001690556007909152812091985062001ae6919062003d88565b6000898152601b6020526040812062001aff9162003d88565b6000898152601c6020526040812062001b189162003d88565b600089815260066020526040902080547fffffffffffffffffffffffff000000000000000000000000000000000000000016905562001b5960028a62003739565b5060a0880151604080516bffffffffffffffffffffffff909216825273ffffffffffffffffffffffffffffffffffffffff8c1660208301528a917fb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff910160405180910390a28062001bca8162004a01565b915050620015f5565b508560195462001be491906200472e565b60195560008b8b868167ffffffffffffffff81111562001c085762001c0862003e75565b60405190808252806020026020018201604052801562001c32578160200160208202803683370190505b508988888860405160200162001c5098979695949392919062004ba7565b60405160208183030381529060405290508973ffffffffffffffffffffffffffffffffffffffff16638e86139b6014600001600c9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663c71249ab60038e73ffffffffffffffffffffffffffffffffffffffff1663aab9edd66040518163ffffffff1660e01b8152600401602060405180830381865afa15801562001d0c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001d32919062004c77565b866040518463ffffffff1660e01b815260040162001d539392919062004c9c565b600060405180830381865afa15801562001d71573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405262001db9919081019062004cc3565b6040518263ffffffff1660e01b815260040162001dd79190620046ce565b600060405180830381600087803b15801562001df257600080fd5b505af115801562001e07573d6000803e3d6000fd5b50506040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8d81166004830152602482018b90527f000000000000000000000000000000000000000000000000000000000000000016925063a9059cbb91506044016020604051808303816000875af115801562001ea1573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001ec7919062004cfc565b50505050505050505050505050565b6002336000908152601a602052604090205460ff16600381111562001eff5762001eff62004191565b1415801562001f3557506003336000908152601a602052604090205460ff16600381111562001f325762001f3262004191565b14155b1562001f6d576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080808080808062001f83888a018a62004ef4565b965096509650965096509650965060005b87518110156200225257600073ffffffffffffffffffffffffffffffffffffffff1687828151811062001fcb5762001fcb620049bc565b60200260200101516060015173ffffffffffffffffffffffffffffffffffffffff1603620020df57858181518110620020085762002008620049bc565b6020026020010151307f0000000000000000000000000000000000000000000000000000000000000000604051620020409062003d7a565b73ffffffffffffffffffffffffffffffffffffffff938416815291831660208301529091166040820152606001604051809103906000f0801580156200208a573d6000803e3d6000fd5b50878281518110620020a057620020a0620049bc565b60200260200101516060019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250505b62002197888281518110620020f857620020f8620049bc565b6020026020010151888381518110620021155762002115620049bc565b6020026020010151878481518110620021325762002132620049bc565b60200260200101518785815181106200214f576200214f620049bc565b60200260200101518786815181106200216c576200216c620049bc565b6020026020010151878781518110620021895762002189620049bc565b602002602001015162002d42565b878181518110620021ac57620021ac620049bc565b60200260200101517f74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71888381518110620021ea57620021ea620049bc565b602002602001015160a0015133604051620022359291906bffffffffffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a280620022498162004a01565b91505062001f94565b50505050505050505050565b600082815260046020908152604091829020825160e081018452815460ff81161515825263ffffffff6101008204811694830194909452650100000000008104841694820185905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004821660c082015291146200235c576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818160a001516200236e919062005025565b600084815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff93841602179055601954620023d691841690620049eb565b6019556040517f23b872dd0000000000000000000000000000000000000000000000000000000081523360048201523060248201526bffffffffffffffffffffffff831660448201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906323b872dd906064016020604051808303816000875af115801562002480573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620024a6919062004cfc565b506040516bffffffffffffffffffffffff83168152339084907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a3505050565b6000818152600460209081526040808320815160e081018352815460ff81161515825263ffffffff6101008204811695830195909552650100000000008104851693820184905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004831660c08201529291141590620025da60005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161490506000601260010160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200267d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620026a391906200504d565b9050828015620026c75750818015620026c5575080846040015163ffffffff16115b155b15620026ff576040517ffbc0357800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8115801562002732575060008581526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314155b156200276a576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8162002780576200277d603282620049eb565b90505b6000858152600460205260409020805463ffffffff80841665010000000000027fffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffffff90921691909117909155620027dc9060029087906200373916565b5060145460808501516bffffffffffffffffffffffff91821691600091168211156200284557608086015162002813908362005067565b90508560a001516bffffffffffffffffffffffff16816bffffffffffffffffffffffff16111562002845575060a08501515b808660a0015162002857919062005067565b600088815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff93841602179055601554620028bf9183911662005025565b601580547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff9290921691909117905560405167ffffffffffffffff84169088907f91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f79118190600090a350505050505050565b600060606000806200294d6200312d565b6000634b56a42e60e01b8888886040516024016200296e939291906200508f565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091529050620029f989826200068d565b929c919b50995090975095505050505050565b62002a1662003747565b62002a2181620037ca565b50565b60006060600080600080600062002a4b8860405180602001604052806000815250620009ff565b959e949d50929b5090995097509550909350915050565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415155b90505b92915050565b6000806000601260010160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905060008173ffffffffffffffffffffffffffffffffffffffff166385df51fd60018473ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa15801562002b2f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062002b5591906200504d565b62002b6191906200472e565b6040518263ffffffff1660e01b815260040162002b8091815260200190565b602060405180830381865afa15801562002b9e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062002bc491906200504d565b601554604080516020810193909352309083015274010000000000000000000000000000000000000000900463ffffffff166060820152608001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201209083015201604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052905060045b600f81101562002cd0578382828151811062002c8c5762002c8c620049bc565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508062002cc78162004a01565b91505062002c6c565b5084600181111562002ce65762002ce662004191565b60f81b81600f8151811062002cff5762002cff620049bc565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535062002d3981620050c3565b95945050505050565b6012547e01000000000000000000000000000000000000000000000000000000000000900460ff161562002da2576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601654835163ffffffff909116101562002de8576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108fc856020015163ffffffff16108062002e265750601554602086015163ffffffff70010000000000000000000000000000000090920482169116115b1562002e5e576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000868152600460205260409020546901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff161562002ec8576040517f6e3b930b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000868152600460209081526040808320885181548a8501518b85015160608d01517fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000009093169315157fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000ff169390931761010063ffffffff92831602177fffffff000000000000000000000000000000000000000000000000ffffffffff1665010000000000938216939093027fffffff0000000000000000000000000000000000000000ffffffffffffffffff1692909217690100000000000000000073ffffffffffffffffffffffffffffffffffffffff9283160217835560808b01516001909301805460a08d015160c08e01516bffffffffffffffffffffffff9687167fffffffffffffffff000000000000000000000000000000000000000000000000909316929092176c010000000000000000000000009690911695909502949094177fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff1678010000000000000000000000000000000000000000000000009490931693909302919091179091556005835281842080547fffffffffffffffffffffffff00000000000000000000000000000000000000001691891691909117905560079091529020620030bb848262005106565b508460a001516bffffffffffffffffffffffff16601954620030de9190620049eb565b6019556000868152601b60205260409020620030fb838262005106565b506000868152601c6020526040902062003116828262005106565b5062003124600287620038c1565b50505050505050565b321562003166576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60008181526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314620031c6576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526004602052604090205465010000000000900463ffffffff9081161462002a21576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000818160045b600f811015620032b3577fff000000000000000000000000000000000000000000000000000000000000008216838260208110620032675762003267620049bc565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916146200329e57506000949350505050565b80620032aa8162004a01565b91505062003225565b5081600f1a6001811115620032cc57620032cc62004191565b949350505050565b6000806000836080015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa15801562003361573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003387919062005248565b50945090925050506000811315806200339f57508142105b80620033c45750828015620033c45750620033bb82426200472e565b8463ffffffff16105b15620033d5576017549550620033d9565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa15801562003445573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200346b919062005248565b50945090925050506000811315806200348357508142105b80620034a85750828015620034a857506200349f82426200472e565b8463ffffffff16105b15620034b9576018549450620034bd565b8094505b50505050915091565b600080620034da88878b60c00151620038cf565b9050600080620034f78b8a63ffffffff16858a8a60018b62003965565b909250905062003508818362005025565b9b9a5050505050505050505050565b6060600083600181111562003530576200353062004191565b03620035fd576000848152600760205260409081902090517f6e04ff0d0000000000000000000000000000000000000000000000000000000091620035789160240162005340565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152905062003732565b600183600181111562003614576200361462004191565b036200370057600082806020019051810190620036329190620053b7565b6000868152600760205260409081902090519192507f40691db4000000000000000000000000000000000000000000000000000000009162003679918491602401620054cb565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091529150620037329050565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b9392505050565b600062002a8d838362003c1d565b60005473ffffffffffffffffffffffffffffffffffffffff16331462003166576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401620011e0565b3373ffffffffffffffffffffffffffffffffffffffff8216036200384b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620011e0565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600062002a8d838362003d28565b60008080856001811115620038e857620038e862004191565b03620038f9575062015f906200391c565b600185600181111562003910576200391062004191565b036200370057506201af405b6200392f63ffffffff8516601462005593565b6200393c846001620055ad565b6200394d9060ff16611d4c62005593565b620039599083620049eb565b62002d399190620049eb565b60008060008960a0015161ffff168762003980919062005593565b90508380156200398f5750803a105b156200399857503a5b6000841562003a20578a610140015173ffffffffffffffffffffffffffffffffffffffff166318b8f6136040518163ffffffff1660e01b8152600401602060405180830381865afa158015620039f2573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003a1891906200504d565b905062003ae1565b6101408b01516016546040517f1254414000000000000000000000000000000000000000000000000000000000815264010000000090910463ffffffff16600482015273ffffffffffffffffffffffffffffffffffffffff90911690631254414090602401602060405180830381865afa15801562003aa3573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003ac991906200504d565b8b60a0015161ffff1662003ade919062005593565b90505b62003af161ffff871682620055c9565b90506000878262003b038c8e620049eb565b62003b0f908662005593565b62003b1b9190620049eb565b62003b2f90670de0b6b3a764000062005593565b62003b3b9190620055c9565b905060008c6040015163ffffffff1664e8d4a5100062003b5c919062005593565b898e6020015163ffffffff16858f8862003b77919062005593565b62003b839190620049eb565b62003b9390633b9aca0062005593565b62003b9f919062005593565b62003bab9190620055c9565b62003bb79190620049eb565b90506b033b2e3c9fd0803ce800000062003bd28284620049eb565b111562003c0b576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b6000818152600183016020526040812054801562003d1657600062003c446001836200472e565b855490915060009062003c5a906001906200472e565b905081811462003cc657600086600001828154811062003c7e5762003c7e620049bc565b906000526020600020015490508087600001848154811062003ca45762003ca4620049bc565b6000918252602080832090910192909255918252600188019052604090208390555b855486908062003cda5762003cda62005605565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505062002a90565b600091505062002a90565b5092915050565b600081815260018301602052604081205462003d715750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915562002a90565b50600062002a90565b6103ca806200563583390190565b50805462003d9690620047ec565b6000825580601f1062003da7575050565b601f01602090049060005260206000209081019062002a2191905b8082111562003dd8576000815560010162003dc2565b5090565b73ffffffffffffffffffffffffffffffffffffffff8116811462002a2157600080fd5b803563ffffffff8116811462003e1457600080fd5b919050565b80356002811062003e1457600080fd5b60008083601f84011262003e3c57600080fd5b50813567ffffffffffffffff81111562003e5557600080fd5b60208301915083602082850101111562003e6e57600080fd5b9250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160e0810167ffffffffffffffff8111828210171562003eca5762003eca62003e75565b60405290565b604051610100810167ffffffffffffffff8111828210171562003eca5762003eca62003e75565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171562003f415762003f4162003e75565b604052919050565b600067ffffffffffffffff82111562003f665762003f6662003e75565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f83011262003fa457600080fd5b813562003fbb62003fb58262003f49565b62003ef7565b81815284602083860101111562003fd157600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060008060e0898b0312156200400b57600080fd5b8835620040188162003ddc565b97506200402860208a0162003dff565b965060408901356200403a8162003ddc565b95506200404a60608a0162003e19565b9450608089013567ffffffffffffffff808211156200406857600080fd5b620040768c838d0162003e29565b909650945060a08b01359150808211156200409057600080fd5b6200409e8c838d0162003f92565b935060c08b0135915080821115620040b557600080fd5b50620040c48b828c0162003f92565b9150509295985092959890939650565b60008060408385031215620040e857600080fd5b82359150602083013567ffffffffffffffff8111156200410757600080fd5b620041158582860162003f92565b9150509250929050565b60005b838110156200413c57818101518382015260200162004122565b50506000910152565b600081518084526200415f8160208601602086016200411f565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600a8110620041f8577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b9052565b841515815260806020820152600062004219608083018662004145565b90506200422a6040830185620041c0565b82606083015295945050505050565b6000806000604084860312156200424f57600080fd5b83359250602084013567ffffffffffffffff8111156200426e57600080fd5b6200427c8682870162003e29565b9497909650939450505050565b600080600080600080600060a0888a031215620042a557600080fd5b8735620042b28162003ddc565b9650620042c26020890162003dff565b95506040880135620042d48162003ddc565b9450606088013567ffffffffffffffff80821115620042f257600080fd5b620043008b838c0162003e29565b909650945060808a01359150808211156200431a57600080fd5b50620043298a828b0162003e29565b989b979a50959850939692959293505050565b871515815260e0602082015260006200435960e083018962004145565b90506200436a6040830188620041c0565b8560608301528460808301528360a08301528260c083015298975050505050505050565b600080600060408486031215620043a457600080fd5b833567ffffffffffffffff80821115620043bd57600080fd5b818601915086601f830112620043d257600080fd5b813581811115620043e257600080fd5b8760208260051b8501011115620043f857600080fd5b60209283019550935050840135620044108162003ddc565b809150509250925092565b600080602083850312156200442f57600080fd5b823567ffffffffffffffff8111156200444757600080fd5b620044558582860162003e29565b90969095509350505050565b80356bffffffffffffffffffffffff8116811462003e1457600080fd5b600080604083850312156200449257600080fd5b82359150620044a46020840162004461565b90509250929050565b600060208284031215620044c057600080fd5b5035919050565b600067ffffffffffffffff821115620044e457620044e462003e75565b5060051b60200190565b600082601f8301126200450057600080fd5b813560206200451362003fb583620044c7565b82815260059290921b840181019181810190868411156200453357600080fd5b8286015b848110156200457857803567ffffffffffffffff811115620045595760008081fd5b620045698986838b010162003f92565b84525091830191830162004537565b509695505050505050565b600080600080606085870312156200459a57600080fd5b84359350602085013567ffffffffffffffff80821115620045ba57600080fd5b620045c888838901620044ee565b94506040870135915080821115620045df57600080fd5b50620045ee8782880162003e29565b95989497509550505050565b6000602082840312156200460d57600080fd5b8135620037328162003ddc565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600063ffffffff8083168181036200466557620046656200461a565b6001019392505050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b602081526000620032cc6020830184866200466f565b60208152600062002a8d602083018462004145565b805162003e148162003ddc565b6000602082840312156200470357600080fd5b8151620037328162003ddc565b60008251620047248184602087016200411f565b9190910192915050565b8181038181111562002a905762002a906200461a565b801515811462002a2157600080fd5b600082601f8301126200476557600080fd5b81516200477662003fb58262003f49565b8181528460208386010111156200478c57600080fd5b620032cc8260208301602087016200411f565b60008060408385031215620047b357600080fd5b8251620047c08162004744565b602084015190925067ffffffffffffffff811115620047de57600080fd5b620041158582860162004753565b600181811c908216806200480157607f821691505b6020821081036200483b577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f8211156200488f57600081815260208120601f850160051c810160208610156200486a5750805b601f850160051c820191505b818110156200488b5782815560010162004876565b5050505b505050565b67ffffffffffffffff831115620048af57620048af62003e75565b620048c783620048c08354620047ec565b8362004841565b6000601f8411600181146200491c5760008515620048e55750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355620049b5565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b828110156200496d57868501358255602094850194600190920191016200494b565b5086821015620049a9577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b8082018082111562002a905762002a906200461a565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820362004a355762004a356200461a565b5060010190565b600081518084526020808501945080840160005b8381101562004afb5781518051151588528381015163ffffffff908116858a01526040808301519091169089015260608082015173ffffffffffffffffffffffffffffffffffffffff16908901526080808201516bffffffffffffffffffffffff169089015260a08082015162004ad6828b01826bffffffffffffffffffffffff169052565b505060c09081015163ffffffff169088015260e0909601959082019060010162004a50565b509495945050505050565b600081518084526020808501945080840160005b8381101562004afb57815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010162004b1a565b600081518084526020808501808196508360051b8101915082860160005b8581101562004b9a57828403895262004b8784835162004145565b9885019893509084019060010162004b6c565b5091979650505050505050565b60e081528760e082015260006101007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8a111562004be457600080fd5b8960051b808c8386013783018381038201602085015262004c088282018b62004a3c565b915050828103604084015262004c1f818962004b06565b9050828103606084015262004c35818862004b06565b9050828103608084015262004c4b818762004b4e565b905082810360a084015262004c61818662004b4e565b905082810360c084015262003508818562004b4e565b60006020828403121562004c8a57600080fd5b815160ff811681146200373257600080fd5b60ff8416815260ff8316602082015260606040820152600062002d39606083018462004145565b60006020828403121562004cd657600080fd5b815167ffffffffffffffff81111562004cee57600080fd5b620032cc8482850162004753565b60006020828403121562004d0f57600080fd5b8151620037328162004744565b600082601f83011262004d2e57600080fd5b8135602062004d4162003fb583620044c7565b82815260059290921b8401810191818101908684111562004d6157600080fd5b8286015b8481101562004578578035835291830191830162004d65565b600082601f83011262004d9057600080fd5b8135602062004da362003fb583620044c7565b82815260e0928302850182019282820191908785111562004dc357600080fd5b8387015b8581101562004e7a5781818a03121562004de15760008081fd5b62004deb62003ea4565b813562004df88162004744565b815262004e0782870162003dff565b86820152604062004e1a81840162003dff565b9082015260608281013562004e2f8162003ddc565b90820152608062004e4283820162004461565b9082015260a062004e5583820162004461565b9082015260c062004e6883820162003dff565b90820152845292840192810162004dc7565b5090979650505050505050565b600082601f83011262004e9957600080fd5b8135602062004eac62003fb583620044c7565b82815260059290921b8401810191818101908684111562004ecc57600080fd5b8286015b848110156200457857803562004ee68162003ddc565b835291830191830162004ed0565b600080600080600080600060e0888a03121562004f1057600080fd5b873567ffffffffffffffff8082111562004f2957600080fd5b62004f378b838c0162004d1c565b985060208a013591508082111562004f4e57600080fd5b62004f5c8b838c0162004d7e565b975060408a013591508082111562004f7357600080fd5b62004f818b838c0162004e87565b965060608a013591508082111562004f9857600080fd5b62004fa68b838c0162004e87565b955060808a013591508082111562004fbd57600080fd5b62004fcb8b838c01620044ee565b945060a08a013591508082111562004fe257600080fd5b62004ff08b838c01620044ee565b935060c08a01359150808211156200500757600080fd5b50620050168a828b01620044ee565b91505092959891949750929550565b6bffffffffffffffffffffffff81811683821601908082111562003d215762003d216200461a565b6000602082840312156200506057600080fd5b5051919050565b6bffffffffffffffffffffffff82811682821603908082111562003d215762003d216200461a565b604081526000620050a4604083018662004b4e565b8281036020840152620050b98185876200466f565b9695505050505050565b805160208083015191908110156200483b577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60209190910360031b1b16919050565b815167ffffffffffffffff81111562005123576200512362003e75565b6200513b81620051348454620047ec565b8462004841565b602080601f8311600181146200519157600084156200515a5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556200488b565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015620051e057888601518255948401946001909101908401620051bf565b50858210156200521d57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b805169ffffffffffffffffffff8116811462003e1457600080fd5b600080600080600060a086880312156200526157600080fd5b6200526c866200522d565b945060208601519350604086015192506060860151915062005291608087016200522d565b90509295509295909350565b60008154620052ac81620047ec565b808552602060018381168015620052cc5760018114620053055762005335565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b890101955062005335565b866000528260002060005b858110156200532d5781548a820186015290830190840162005310565b890184019650505b505050505092915050565b60208152600062002a8d60208301846200529d565b600082601f8301126200536757600080fd5b815160206200537a62003fb583620044c7565b82815260059290921b840181019181810190868411156200539a57600080fd5b8286015b848110156200457857805183529183019183016200539e565b600060208284031215620053ca57600080fd5b815167ffffffffffffffff80821115620053e357600080fd5b908301906101008286031215620053f957600080fd5b6200540362003ed0565b82518152602083015160208201526040830151604082015260608301516060820152608083015160808201526200543d60a08401620046e3565b60a082015260c0830151828111156200545557600080fd5b620054638782860162005355565b60c08301525060e0830151828111156200547c57600080fd5b6200548a8782860162004753565b60e08301525095945050505050565b600081518084526020808501945080840160005b8381101562004afb57815187529582019590820190600101620054ad565b60408152825160408201526020830151606082015260408301516080820152606083015160a0820152608083015160c082015273ffffffffffffffffffffffffffffffffffffffff60a08401511660e0820152600060c08401516101008081850152506200553e61014084018262005499565b905060e08501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0848303016101208501526200557c828262004145565b915050828103602084015262002d3981856200529d565b808202811582820484141762002a905762002a906200461a565b60ff818116838216019081111562002a905762002a906200461a565b60008262005600577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfe60c060405234801561001057600080fd5b506040516103ca3803806103ca83398101604081905261002f91610076565b600080546001600160a01b0319166001600160a01b039384161790559181166080521660a0526100b9565b80516001600160a01b038116811461007157600080fd5b919050565b60008060006060848603121561008b57600080fd5b6100948461005a565b92506100a26020850161005a565b91506100b06040850161005a565b90509250925092565b60805160a0516102e76100e36000396000603801526000818160c4015261011701526102e76000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806379188d161461007b578063f00e6a2a146100aa575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e808015610076573d6000f35b3d6000fd5b61008e6100893660046101c1565b6100ee565b6040805192151583526020830191909152015b60405180910390f35b60405173ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001681526020016100a1565b60008054819073ffffffffffffffffffffffffffffffffffffffff16331461011557600080fd5b7f00000000000000000000000000000000000000000000000000000000000000005a91505a61138881101561014957600080fd5b61138881039050856040820482031161016157600080fd5b50803b61016d57600080fd5b6000808551602087016000858af192505a610188908361029a565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080604083850312156101d457600080fd5b82359150602083013567ffffffffffffffff808211156101f357600080fd5b818501915085601f83011261020757600080fd5b81358181111561021957610219610192565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561025f5761025f610192565b8160405282815288602084870101111561027857600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b818103818111156102d4577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b9291505056fea164736f6c6343000813000aa164736f6c6343000813000a",
>>>>>>> 911f05e7a3 (adjust gas overhead again)
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
