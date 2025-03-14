// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package keeper_registry_logic_a_wrapper_2_1

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

var KeeperRegistryLogicA21MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"logicB\",\"type\":\"address\",\"internalType\":\"contractKeeperRegistryLogicB2_1\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addFunds\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"checkCallback\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"values\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"upkeepNeeded\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"performData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"upkeepFailureReason\",\"type\":\"uint8\",\"internalType\":\"enumKeeperRegistryBase2_1.UpkeepFailureReason\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"checkUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"triggerData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"upkeepNeeded\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"performData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"upkeepFailureReason\",\"type\":\"uint8\",\"internalType\":\"enumKeeperRegistryBase2_1.UpkeepFailureReason\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fastGasWei\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"linkNative\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"checkUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"upkeepNeeded\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"performData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"upkeepFailureReason\",\"type\":\"uint8\",\"internalType\":\"enumKeeperRegistryBase2_1.UpkeepFailureReason\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fastGasWei\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"linkNative\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"executeCallback\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"payload\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"upkeepNeeded\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"performData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"upkeepFailureReason\",\"type\":\"uint8\",\"internalType\":\"enumKeeperRegistryBase2_1.UpkeepFailureReason\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"fallbackTo\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"migrateUpkeeps\",\"inputs\":[{\"name\":\"ids\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"destination\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"receiveUpkeeps\",\"inputs\":[{\"name\":\"encodedUpkeeps\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerUpkeep\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"triggerType\",\"type\":\"uint8\",\"internalType\":\"enumKeeperRegistryBase2_1.Trigger\"},{\"name\":\"checkData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"triggerConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerUpkeep\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"checkData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUpkeepTriggerConfig\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"triggerConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AdminPrivilegeConfigSet\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"privilegeConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CancelledUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DedupKeyAdded\",\"inputs\":[{\"name\":\"dedupKey\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsAdded\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsWithdrawn\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InsufficientFundsUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnerFundsWithdrawn\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeesUpdated\",\"inputs\":[{\"name\":\"transmitters\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"payees\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeeshipTransferRequested\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeeshipTransferred\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PaymentWithdrawn\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"payee\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReorgedUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StaleUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepAdminTransferRequested\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepAdminTransferred\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepCanceled\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"atBlockHeight\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepCheckDataSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"newCheckData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepGasLimitSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"gasLimit\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepMigrated\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"remainingBalance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"destination\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepOffchainConfigSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPaused\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPerformed\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"success\",\"type\":\"bool\",\"indexed\":true,\"internalType\":\"bool\"},{\"name\":\"totalPayment\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"gasOverhead\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPrivilegeConfigSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"privilegeConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepReceived\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"startingBalance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"importedFrom\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepRegistered\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"performGas\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"admin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepTriggerConfigSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"triggerConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepUnpaused\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ArrayHasNoEntries\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotCancel\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CheckDataExceedsLimit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConfigDigestMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateEntry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GasLimitCanOnlyIncrease\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GasLimitOutsideRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfFaultyOracles\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfSignatures\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IndexOutOfRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientFunds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidDataLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRecipient\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReport\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSigner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTransmitter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTrigger\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTriggerType\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MigrationNotPermitted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAContract\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyActiveSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyActiveTransmitters\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByLINKToken\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwnerOrAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByProposedAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByProposedPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyPausedUpkeep\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlySimulatedBackend\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyUnpausedUpkeep\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterLengthError\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PaymentGreaterThanAllLINK\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RegistryPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RepeatedSigner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RepeatedTransmitter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TargetCheckReverted\",\"inputs\":[{\"name\":\"reason\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"TooManyOracles\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TranscoderNotSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepAlreadyExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepCancelled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepNotCanceled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepNotNeeded\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ValueNotChanged\",\"inputs\":[]}]",
	Bin: "0x6101406040523480156200001257600080fd5b5060405162006295380380620062958339810160408190526200003591620003df565b80816001600160a01b0316634b4fd03b6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200009b919062000406565b826001600160a01b031663ca30e6036040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000da573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001009190620003df565b836001600160a01b031663b10b673c6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200013f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001659190620003df565b846001600160a01b0316636709d0e56040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca9190620003df565b856001600160a01b0316635425d8ac6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000209573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200022f9190620003df565b3380600081620002865760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620002b957620002b9816200031b565b505050846002811115620002d157620002d162000429565b60e0816002811115620002e857620002e862000429565b9052506001600160a01b0393841660805291831660a052821660c0528116610100529190911661012052506200043f9050565b336001600160a01b03821603620003755760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200027d565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b0381168114620003dc57600080fd5b50565b600060208284031215620003f257600080fd5b8151620003ff81620003c6565b9392505050565b6000602082840312156200041957600080fd5b815160038110620003ff57600080fd5b634e487b7160e01b600052602160045260246000fd5b60805160a05160c05160e0516101005161012051615ddc620004b96000396000818161010001526101c50152600081816104a501526120650152600081816135fb0152818161383101528181613a790152613c21015260006131a501526000613289015260008181611ea701526124730152615ddc6000f3fe608060405260043610620000fe5760003560e01c806385c1b0ba1162000097578063c80480221162000061578063c80480221462000343578063ce7dc5b41462000368578063f2fde38b146200038d578063f7d334ba14620003b257620000fe565b806385c1b0ba14620002a75780638da5cb5b14620002cc5780638e86139b14620002f9578063948108f7146200031e57620000fe565b80634ee88d3511620000d95780634ee88d35146200020b5780636ded9eae146200023057806371791aa0146200025557806379ba5097146200028f57620000fe565b806328f32f38146200014657806329c5efad146200017e578063349e8cca14620001b5575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e8080156200013f573d6000f35b3d6000fd5b005b3480156200015357600080fd5b506200016b62000165366004620042ae565b620003d7565b6040519081526020015b60405180910390f35b3480156200018b57600080fd5b50620001a36200019d36600462004394565b62000750565b604051620001759493929190620044bc565b348015620001c257600080fd5b507f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200162000175565b3480156200021857600080fd5b50620001446200022a366004620044f9565b620009f4565b3480156200023d57600080fd5b506200016b6200024f36600462004549565b62000a5c565b3480156200026257600080fd5b506200027a6200027436600462004394565b62000ac2565b604051620001759796959493929190620045fc565b3480156200029c57600080fd5b5062000144620011b4565b348015620002b457600080fd5b5062000144620002c63660046200464e565b620012b7565b348015620002d957600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff16620001e5565b3480156200030657600080fd5b506200014462000318366004620046db565b62001f28565b3480156200032b57600080fd5b50620001446200033d3660046200473e565b620022b0565b3480156200035057600080fd5b5062000144620003623660046200476d565b62002543565b3480156200037557600080fd5b50620001a36200038736600462004843565b6200290a565b3480156200039a57600080fd5b5062000144620003ac366004620048ba565b620029da565b348015620003bf57600080fd5b506200027a620003d13660046200476d565b620029f2565b6000805473ffffffffffffffffffffffffffffffffffffffff1633148015906200040b57506200040960093362002a30565b155b1562000443576040517fd48b678b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff89163b62000492576040517f09ee12d500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6200049d8662002a64565b9050600089307f0000000000000000000000000000000000000000000000000000000000000000604051620004d2906200403f565b73ffffffffffffffffffffffffffffffffffffffff938416815291831660208301529091166040820152606001604051809103906000f0801580156200051c573d6000803e3d6000fd5b509050620005e3826040518060e001604052806000151581526020018c63ffffffff16815260200163ffffffff801681526020018473ffffffffffffffffffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff168152602001600063ffffffff168152508a89898080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508b92508a915062002c089050565b6014805474010000000000000000000000000000000000000000900463ffffffff169080620006128362004909565b91906101000a81548163ffffffff021916908363ffffffff16021790555050817fbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d0128a8a6040516200068b92919063ffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a2817fcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d8787604051620006c792919062004978565b60405180910390a2817f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664856040516200070191906200498e565b60405180910390a2817f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850846040516200073b91906200498e565b60405180910390a25098975050505050505050565b600060606000806200076162002fe3565b600086815260046020908152604091829020825160e081018452815460ff811615158252610100810463ffffffff90811694830194909452650100000000008104841694820194909452690100000000000000000090930473ffffffffffffffffffffffffffffffffffffffff166060840152600101546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a0840152780100000000000000000000000000000000000000000000000090041660c08201525a9150600080826060015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200087b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620008a19190620049b0565b73ffffffffffffffffffffffffffffffffffffffff166013600101600c9054906101000a900463ffffffff1663ffffffff1689604051620008e39190620049d0565b60006040518083038160008787f1925050503d806000811462000923576040519150601f19603f3d011682016040523d82523d6000602084013e62000928565b606091505b50915091505a6200093a9085620049ee565b93508162000965576000604051806020016040528060008152506007965096509650505050620009eb565b808060200190518101906200097b919062004a5f565b909750955086620009a9576000604051806020016040528060008152506004965096509650505050620009eb565b601554865164010000000090910463ffffffff161015620009e7576000604051806020016040528060008152506005965096509650505050620009eb565b5050505b92959194509250565b620009ff836200301e565b6000838152601a6020526040902062000a1a82848362004b54565b50827f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664838360405162000a4f92919062004978565b60405180910390a2505050565b600062000ab688888860008989604051806020016040528060008152508a8a8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250620003d792505050565b98975050505050505050565b60006060600080600080600062000ad862002fe3565b600062000ae58a620030d4565b905060006012604051806101200160405290816000820160009054906101000a900460ff1660ff1660ff1681526020016000820160019054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160059054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160099054906101000a900462ffffff1662ffffff1662ffffff16815260200160008201600c9054906101000a900461ffff1661ffff1661ffff16815260200160008201600e9054906101000a900460ff1615151515815260200160008201600f9054906101000a900460ff161515151581526020016000820160109054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160008201601c9054906101000a900463ffffffff1663ffffffff1663ffffffff168152505090506000600460008d81526020019081526020016000206040518060e00160405290816000820160009054906101000a900460ff161515151581526020016000820160019054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160059054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160099054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160018201600c9054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff1681526020016001820160189054906101000a900463ffffffff1663ffffffff1663ffffffff168152505090508160a001511562000e09576000604051806020016040528060008152506009600084602001516000808263ffffffff1692509950995099509950995099509950505050620011a8565b604081015163ffffffff9081161462000e5a576000604051806020016040528060008152506001600084602001516000808263ffffffff1692509950995099509950995099509950505050620011a8565b80511562000ea0576000604051806020016040528060008152506002600084602001516000808263ffffffff1692509950995099509950995099509950505050620011a8565b62000eab8262003182565b602083015160155492975090955060009162000edd918591879190640100000000900463ffffffff168a8a8762003374565b9050806bffffffffffffffffffffffff168260a001516bffffffffffffffffffffffff16101562000f47576000604051806020016040528060008152506006600085602001516000808263ffffffff1692509a509a509a509a509a509a509a5050505050620011a8565b600062000f568e868f620033c5565b90505a9850600080846060015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000fae573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000fd49190620049b0565b73ffffffffffffffffffffffffffffffffffffffff166013600101600c9054906101000a900463ffffffff1663ffffffff1684604051620010169190620049d0565b60006040518083038160008787f1925050503d806000811462001056576040519150601f19603f3d011682016040523d82523d6000602084013e6200105b565b606091505b50915091505a6200106d908c620049ee565b9a5081620010ed5760155481516801000000000000000090910463ffffffff161015620010ca57505060408051602080820190925260008082529490910151939c509a50600899505063ffffffff9091169550620011a892505050565b602090940151939b5060039a505063ffffffff9092169650620011a89350505050565b8080602001905181019062001103919062004a5f565b909e509c508d6200114457505060408051602080820190925260008082529490910151939c509a50600499505063ffffffff9091169550620011a892505050565b6015548d5164010000000090910463ffffffff1610156200119557505060408051602080820190925260008082529490910151939c509a50600599505063ffffffff9091169550620011a892505050565b505050506020015163ffffffff16945050505b92959891949750929550565b60015473ffffffffffffffffffffffffffffffffffffffff1633146200123b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600173ffffffffffffffffffffffffffffffffffffffff821660009081526019602052604090205460ff166003811115620012f657620012f662004451565b14158015620013425750600373ffffffffffffffffffffffffffffffffffffffff821660009081526019602052604090205460ff1660038111156200133f576200133f62004451565b14155b156200137a576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6013546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff16620013da576040517fd12d7d8d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082900362001416576040517f2c2fc94100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160e081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c081018290526000808567ffffffffffffffff8111156200146d576200146d62004135565b60405190808252806020026020018201604052801562001497578160200160208202803683370190505b50905060008667ffffffffffffffff811115620014b857620014b862004135565b6040519080825280602002602001820160405280156200153f57816020015b6040805160e08101825260008082526020808301829052928201819052606082018190526080820181905260a0820181905260c082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181620014d75790505b50905060008767ffffffffffffffff81111562001560576200156062004135565b6040519080825280602002602001820160405280156200159557816020015b60608152602001906001900390816200157f5790505b50905060008867ffffffffffffffff811115620015b657620015b662004135565b604051908082528060200260200182016040528015620015eb57816020015b6060815260200190600190039081620015d55790505b50905060008967ffffffffffffffff8111156200160c576200160c62004135565b6040519080825280602002602001820160405280156200164157816020015b60608152602001906001900390816200162b5790505b50905060005b8a81101562001c25578b8b8281811062001665576200166562004c7c565b60209081029290920135600081815260048452604090819020815160e081018352815460ff811615158252610100810463ffffffff90811697830197909752650100000000008104871693820193909352690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff166060830152600101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490931660c08401529a50909850620017449050896200301e565b60608801516040517f1a5da6c800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8c8116600483015290911690631a5da6c890602401600060405180830381600087803b158015620017b457600080fd5b505af1158015620017c9573d6000803e3d6000fd5b5050505087858281518110620017e357620017e362004c7c565b6020026020010181905250600560008a815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1686828151811062001837576200183762004c7c565b73ffffffffffffffffffffffffffffffffffffffff90921660209283029190910182015260008a81526007909152604090208054620018769062004aac565b80601f0160208091040260200160405190810160405280929190818152602001828054620018a49062004aac565b8015620018f55780601f10620018c957610100808354040283529160200191620018f5565b820191906000526020600020905b815481529060010190602001808311620018d757829003601f168201915b50505050508482815181106200190f576200190f62004c7c565b6020026020010181905250601a60008a815260200190815260200160002080546200193a9062004aac565b80601f0160208091040260200160405190810160405280929190818152602001828054620019689062004aac565b8015620019b95780601f106200198d57610100808354040283529160200191620019b9565b820191906000526020600020905b8154815290600101906020018083116200199b57829003601f168201915b5050505050838281518110620019d357620019d362004c7c565b6020026020010181905250601b60008a81526020019081526020016000208054620019fe9062004aac565b80601f016020809104026020016040519081016040528092919081815260200182805462001a2c9062004aac565b801562001a7d5780601f1062001a515761010080835404028352916020019162001a7d565b820191906000526020600020905b81548152906001019060200180831162001a5f57829003601f168201915b505050505082828151811062001a975762001a9762004c7c565b60200260200101819052508760a001516bffffffffffffffffffffffff168762001ac2919062004cab565b60008a815260046020908152604080832080547fffffff000000000000000000000000000000000000000000000000000000000016815560010180547fffffffff000000000000000000000000000000000000000000000000000000001690556007909152812091985062001b3891906200404d565b6000898152601a6020526040812062001b51916200404d565b6000898152601b6020526040812062001b6a916200404d565b600089815260066020526040902080547fffffffffffffffffffffffff000000000000000000000000000000000000000016905562001bab60028a620035e7565b5060a0880151604080516bffffffffffffffffffffffff909216825273ffffffffffffffffffffffffffffffffffffffff8c1660208301528a917fb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff910160405180910390a28062001c1c8162004cc1565b91505062001647565b508560185462001c369190620049ee565b60185560008b8b868167ffffffffffffffff81111562001c5a5762001c5a62004135565b60405190808252806020026020018201604052801562001c84578160200160208202803683370190505b508988888860405160200162001ca298979695949392919062004e67565b60405160208183030381529060405290508973ffffffffffffffffffffffffffffffffffffffff16638e86139b6013600001600c9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663c71249ab60038e73ffffffffffffffffffffffffffffffffffffffff1663aab9edd66040518163ffffffff1660e01b8152600401602060405180830381865afa15801562001d5e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001d84919062004f37565b866040518463ffffffff1660e01b815260040162001da59392919062004f5c565b600060405180830381865afa15801562001dc3573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405262001e0b919081019062004f83565b6040518263ffffffff1660e01b815260040162001e2991906200498e565b600060405180830381600087803b15801562001e4457600080fd5b505af115801562001e59573d6000803e3d6000fd5b50506040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8d81166004830152602482018b90527f000000000000000000000000000000000000000000000000000000000000000016925063a9059cbb91506044016020604051808303816000875af115801562001ef3573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001f19919062004fbc565b50505050505050505050505050565b60023360009081526019602052604090205460ff16600381111562001f515762001f5162004451565b1415801562001f87575060033360009081526019602052604090205460ff16600381111562001f845762001f8462004451565b14155b1562001fbf576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080808080808062001fd5888a018a620051b4565b965096509650965096509650965060005b8751811015620022a457600073ffffffffffffffffffffffffffffffffffffffff168782815181106200201d576200201d62004c7c565b60200260200101516060015173ffffffffffffffffffffffffffffffffffffffff160362002131578581815181106200205a576200205a62004c7c565b6020026020010151307f000000000000000000000000000000000000000000000000000000000000000060405162002092906200403f565b73ffffffffffffffffffffffffffffffffffffffff938416815291831660208301529091166040820152606001604051809103906000f080158015620020dc573d6000803e3d6000fd5b50878281518110620020f257620020f262004c7c565b60200260200101516060019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250505b620021e98882815181106200214a576200214a62004c7c565b602002602001015188838151811062002167576200216762004c7c565b602002602001015187848151811062002184576200218462004c7c565b6020026020010151878581518110620021a157620021a162004c7c565b6020026020010151878681518110620021be57620021be62004c7c565b6020026020010151878781518110620021db57620021db62004c7c565b602002602001015162002c08565b878181518110620021fe57620021fe62004c7c565b60200260200101517f74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a718883815181106200223c576200223c62004c7c565b602002602001015160a0015133604051620022879291906bffffffffffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a2806200229b8162004cc1565b91505062001fe6565b50505050505050505050565b600082815260046020908152604091829020825160e081018452815460ff81161515825263ffffffff6101008204811694830194909452650100000000008104841694820185905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004821660c08201529114620023ae576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818160a00151620023c09190620052e5565b600084815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff93841602179055601854620024289184169062004cab565b6018556040517f23b872dd0000000000000000000000000000000000000000000000000000000081523360048201523060248201526bffffffffffffffffffffffff831660448201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906323b872dd906064016020604051808303816000875af1158015620024d2573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620024f8919062004fbc565b506040516bffffffffffffffffffffffff83168152339084907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a3505050565b6000818152600460209081526040808320815160e081018352815460ff81161515825263ffffffff6101008204811695830195909552650100000000008104851693820184905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004831660c082015292911415906200262c60005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614905081801562002687575080801562002685575062002678620035f5565b836040015163ffffffff16115b155b15620026bf576040517ffbc0357800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80158015620026f2575060008481526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314155b156200272a576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600062002736620035f5565b9050816200274e576200274b60328262004cab565b90505b6000858152600460205260409020805463ffffffff80841665010000000000027fffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffffff90921691909117909155620027aa906002908790620035e716565b5060135460808501516bffffffffffffffffffffffff918216916000911682111562002813576080860151620027e190836200530d565b90508560a001516bffffffffffffffffffffffff16816bffffffffffffffffffffffff16111562002813575060a08501515b808660a001516200282591906200530d565b600088815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff938416021790556014546200288d91839116620052e5565b601480547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff9290921691909117905560405167ffffffffffffffff84169088907f91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f79118190600090a350505050505050565b600060606000806200291b62002fe3565b6000634b56a42e60e01b8888886040516024016200293c9392919062005335565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091529050620029c7898262000750565b929c919b50995090975095505050505050565b620029e4620036b1565b620029ef8162003734565b50565b60006060600080600080600062002a19886040518060200160405280600081525062000ac2565b959e949d50929b5090995097509550909350915050565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415155b90505b92915050565b600080600062002a8b600162002a79620035f5565b62002a859190620049ee565b6200382b565b601454604080516020810193909352309083015274010000000000000000000000000000000000000000900463ffffffff166060820152608001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201209083015201604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052905060045b600f81101562002b97578282828151811062002b535762002b5362004c7c565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508062002b8e8162004cc1565b91505062002b33565b5083600181111562002bad5762002bad62004451565b60f81b81600f8151811062002bc65762002bc662004c7c565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535062002c008162005369565b949350505050565b6012546e010000000000000000000000000000900460ff161562002c58576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601554835163ffffffff909116101562002c9e576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108fc856020015163ffffffff16108062002cdc5750601454602086015163ffffffff70010000000000000000000000000000000090920482169116115b1562002d14576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000868152600460205260409020546901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff161562002d7e576040517f6e3b930b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000868152600460209081526040808320885181548a8501518b85015160608d01517fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000009093169315157fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000ff169390931761010063ffffffff92831602177fffffff000000000000000000000000000000000000000000000000ffffffffff1665010000000000938216939093027fffffff0000000000000000000000000000000000000000ffffffffffffffffff1692909217690100000000000000000073ffffffffffffffffffffffffffffffffffffffff9283160217835560808b01516001909301805460a08d015160c08e01516bffffffffffffffffffffffff9687167fffffffffffffffff000000000000000000000000000000000000000000000000909316929092176c010000000000000000000000009690911695909502949094177fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff1678010000000000000000000000000000000000000000000000009490931693909302919091179091556005835281842080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169189169190911790556007909152902062002f718482620053ac565b508460a001516bffffffffffffffffffffffff1660185462002f94919062004cab565b6018556000868152601a6020526040902062002fb18382620053ac565b506000868152601b6020526040902062002fcc8282620053ac565b5062002fda60028762003993565b50505050505050565b32156200301c576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60008181526005602052604090205473ffffffffffffffffffffffffffffffffffffffff1633146200307c576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526004602052604090205465010000000000900463ffffffff90811614620029ef576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000818160045b600f81101562003169577fff0000000000000000000000000000000000000000000000000000000000000082168382602081106200311d576200311d62004c7c565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916146200315457506000949350505050565b80620031608162004cc1565b915050620030db565b5081600f1a600181111562002c005762002c0062004451565b6000806000836060015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa1580156200320f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620032359190620054ee565b50945090925050506000811315806200324d57508142105b80620032725750828015620032725750620032698242620049ee565b8463ffffffff16105b156200328357601654955062003287565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa158015620032f3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620033199190620054ee565b50945090925050506000811315806200333157508142105b806200335657508280156200335657506200334d8242620049ee565b8463ffffffff16105b15620033675760175494506200336b565b8094505b50505050915091565b6000806200338888878b60000151620039a1565b9050600080620033a58b8a63ffffffff16858a8a60018b62003a40565b9092509050620033b68183620052e5565b9b9a5050505050505050505050565b60606000836001811115620033de57620033de62004451565b03620034ab576000848152600760205260409081902090517f6e04ff0d00000000000000000000000000000000000000000000000000000000916200342691602401620055e6565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091529050620035e0565b6001836001811115620034c257620034c262004451565b03620035ae57600082806020019051810190620034e091906200565d565b6000868152600760205260409081902090519192507f40691db400000000000000000000000000000000000000000000000000000000916200352791849160240162005771565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091529150620035e09050565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b9392505050565b600062002a5b838362003ee2565b600060017f000000000000000000000000000000000000000000000000000000000000000060028111156200362e576200362e62004451565b03620036ac57606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562003681573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620036a7919062005839565b905090565b504390565b60005473ffffffffffffffffffffffffffffffffffffffff1633146200301c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162001232565b3373ffffffffffffffffffffffffffffffffffffffff821603620037b5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162001232565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060017f0000000000000000000000000000000000000000000000000000000000000000600281111562003864576200386462004451565b0362003989576000606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620038b9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620038df919062005839565b90508083101580620038fd5750610100620038fb8483620049ee565b115b156200390c5750600092915050565b6040517f2b407a8200000000000000000000000000000000000000000000000000000000815260048101849052606490632b407a8290602401602060405180830381865afa15801562003963573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620035e0919062005839565b504090565b919050565b600062002a5b838362003fed565b60008080856001811115620039ba57620039ba62004451565b03620039cb575062015f90620039ee565b6001856001811115620039e257620039e262004451565b03620035ae57506201adb05b62003a0163ffffffff8516601462005853565b62003a0e84600162005893565b62003a1f9060ff16611d4c62005853565b62003a2b908362004cab565b62003a37919062004cab565b95945050505050565b6000806000896080015161ffff168762003a5b919062005853565b905083801562003a6a5750803a105b1562003a7357503a5b600060027f0000000000000000000000000000000000000000000000000000000000000000600281111562003aac5762003aac62004451565b0362003c1d57604080516000815260208101909152851562003b105760003660405180608001604052806048815260200162005d886048913960405160200162003af993929190620058af565b604051602081830303815290604052905062003b7e565b60155462003b2e90640100000000900463ffffffff166004620058d8565b63ffffffff1667ffffffffffffffff81111562003b4f5762003b4f62004135565b6040519080825280601f01601f19166020018201604052801562003b7a576020820181803683370190505b5090505b6040517f49948e0e00000000000000000000000000000000000000000000000000000000815273420000000000000000000000000000000000000f906349948e0e9062003bd09084906004016200498e565b602060405180830381865afa15801562003bee573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003c14919062005839565b91505062003d87565b60017f0000000000000000000000000000000000000000000000000000000000000000600281111562003c545762003c5462004451565b0362003d8757841562003cdc57606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562003cae573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003cd4919062005839565b905062003d87565b6000606c73ffffffffffffffffffffffffffffffffffffffff166341b247a86040518163ffffffff1660e01b815260040160c060405180830381865afa15801562003d2b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003d51919062005907565b505060155492945062003d7693505050640100000000900463ffffffff168262005853565b62003d8390601062005853565b9150505b8462003da657808b6080015161ffff1662003da3919062005853565b90505b62003db661ffff87168262005952565b90506000878262003dc88c8e62004cab565b62003dd4908662005853565b62003de0919062004cab565b62003df490670de0b6b3a764000062005853565b62003e00919062005952565b905060008c6040015163ffffffff1664e8d4a5100062003e21919062005853565b898e6020015163ffffffff16858f8862003e3c919062005853565b62003e48919062004cab565b62003e5890633b9aca0062005853565b62003e64919062005853565b62003e70919062005952565b62003e7c919062004cab565b90506b033b2e3c9fd0803ce800000062003e97828462004cab565b111562003ed0576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b6000818152600183016020526040812054801562003fdb57600062003f09600183620049ee565b855490915060009062003f1f90600190620049ee565b905081811462003f8b57600086600001828154811062003f435762003f4362004c7c565b906000526020600020015490508087600001848154811062003f695762003f6962004c7c565b6000918252602080832090910192909255918252600188019052604090208390555b855486908062003f9f5762003f9f6200598e565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505062002a5e565b600091505062002a5e565b5092915050565b6000818152600183016020526040812054620040365750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915562002a5e565b50600062002a5e565b6103ca80620059be83390190565b5080546200405b9062004aac565b6000825580601f106200406c575050565b601f016020900490600052602060002090810190620029ef91905b808211156200409d576000815560010162004087565b5090565b73ffffffffffffffffffffffffffffffffffffffff81168114620029ef57600080fd5b803563ffffffff811681146200398e57600080fd5b8035600281106200398e57600080fd5b60008083601f840112620040fc57600080fd5b50813567ffffffffffffffff8111156200411557600080fd5b6020830191508360208285010111156200412e57600080fd5b9250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160e0810167ffffffffffffffff811182821017156200418a576200418a62004135565b60405290565b604051610100810167ffffffffffffffff811182821017156200418a576200418a62004135565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171562004201576200420162004135565b604052919050565b600067ffffffffffffffff82111562004226576200422662004135565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f8301126200426457600080fd5b81356200427b620042758262004209565b620041b7565b8181528460208386010111156200429157600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060008060e0898b031215620042cb57600080fd5b8835620042d881620040a1565b9750620042e860208a01620040c4565b96506040890135620042fa81620040a1565b95506200430a60608a01620040d9565b9450608089013567ffffffffffffffff808211156200432857600080fd5b620043368c838d01620040e9565b909650945060a08b01359150808211156200435057600080fd5b6200435e8c838d0162004252565b935060c08b01359150808211156200437557600080fd5b50620043848b828c0162004252565b9150509295985092959890939650565b60008060408385031215620043a857600080fd5b82359150602083013567ffffffffffffffff811115620043c757600080fd5b620043d58582860162004252565b9150509250929050565b60005b83811015620043fc578181015183820152602001620043e2565b50506000910152565b600081518084526200441f816020860160208601620043df565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600a8110620044b8577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b9052565b8415158152608060208201526000620044d9608083018662004405565b9050620044ea604083018562004480565b82606083015295945050505050565b6000806000604084860312156200450f57600080fd5b83359250602084013567ffffffffffffffff8111156200452e57600080fd5b6200453c86828701620040e9565b9497909650939450505050565b600080600080600080600060a0888a0312156200456557600080fd5b87356200457281620040a1565b96506200458260208901620040c4565b955060408801356200459481620040a1565b9450606088013567ffffffffffffffff80821115620045b257600080fd5b620045c08b838c01620040e9565b909650945060808a0135915080821115620045da57600080fd5b50620045e98a828b01620040e9565b989b979a50959850939692959293505050565b871515815260e0602082015260006200461960e083018962004405565b90506200462a604083018862004480565b8560608301528460808301528360a08301528260c083015298975050505050505050565b6000806000604084860312156200466457600080fd5b833567ffffffffffffffff808211156200467d57600080fd5b818601915086601f8301126200469257600080fd5b813581811115620046a257600080fd5b8760208260051b8501011115620046b857600080fd5b60209283019550935050840135620046d081620040a1565b809150509250925092565b60008060208385031215620046ef57600080fd5b823567ffffffffffffffff8111156200470757600080fd5b6200471585828601620040e9565b90969095509350505050565b80356bffffffffffffffffffffffff811681146200398e57600080fd5b600080604083850312156200475257600080fd5b82359150620047646020840162004721565b90509250929050565b6000602082840312156200478057600080fd5b5035919050565b600067ffffffffffffffff821115620047a457620047a462004135565b5060051b60200190565b600082601f830112620047c057600080fd5b81356020620047d3620042758362004787565b82815260059290921b84018101918181019086841115620047f357600080fd5b8286015b848110156200483857803567ffffffffffffffff811115620048195760008081fd5b620048298986838b010162004252565b845250918301918301620047f7565b509695505050505050565b600080600080606085870312156200485a57600080fd5b84359350602085013567ffffffffffffffff808211156200487a57600080fd5b6200488888838901620047ae565b945060408701359150808211156200489f57600080fd5b50620048ae87828801620040e9565b95989497509550505050565b600060208284031215620048cd57600080fd5b8135620035e081620040a1565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600063ffffffff808316818103620049255762004925620048da565b6001019392505050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b60208152600062002c006020830184866200492f565b60208152600062002a5b602083018462004405565b80516200398e81620040a1565b600060208284031215620049c357600080fd5b8151620035e081620040a1565b60008251620049e4818460208701620043df565b9190910192915050565b8181038181111562002a5e5762002a5e620048da565b8015158114620029ef57600080fd5b600082601f83011262004a2557600080fd5b815162004a36620042758262004209565b81815284602083860101111562004a4c57600080fd5b62002c00826020830160208701620043df565b6000806040838503121562004a7357600080fd5b825162004a808162004a04565b602084015190925067ffffffffffffffff81111562004a9e57600080fd5b620043d58582860162004a13565b600181811c9082168062004ac157607f821691505b60208210810362004afb577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f82111562004b4f57600081815260208120601f850160051c8101602086101562004b2a5750805b601f850160051c820191505b8181101562004b4b5782815560010162004b36565b5050505b505050565b67ffffffffffffffff83111562004b6f5762004b6f62004135565b62004b878362004b80835462004aac565b8362004b01565b6000601f84116001811462004bdc576000851562004ba55750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b17835562004c75565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b8281101562004c2d578685013582556020948501946001909201910162004c0b565b508682101562004c69577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b8082018082111562002a5e5762002a5e620048da565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820362004cf55762004cf5620048da565b5060010190565b600081518084526020808501945080840160005b8381101562004dbb5781518051151588528381015163ffffffff908116858a01526040808301519091169089015260608082015173ffffffffffffffffffffffffffffffffffffffff16908901526080808201516bffffffffffffffffffffffff169089015260a08082015162004d96828b01826bffffffffffffffffffffffff169052565b505060c09081015163ffffffff169088015260e0909601959082019060010162004d10565b509495945050505050565b600081518084526020808501945080840160005b8381101562004dbb57815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010162004dda565b600081518084526020808501808196508360051b8101915082860160005b8581101562004e5a57828403895262004e4784835162004405565b9885019893509084019060010162004e2c565b5091979650505050505050565b60e081528760e082015260006101007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8a111562004ea457600080fd5b8960051b808c8386013783018381038201602085015262004ec88282018b62004cfc565b915050828103604084015262004edf818962004dc6565b9050828103606084015262004ef5818862004dc6565b9050828103608084015262004f0b818762004e0e565b905082810360a084015262004f21818662004e0e565b905082810360c0840152620033b6818562004e0e565b60006020828403121562004f4a57600080fd5b815160ff81168114620035e057600080fd5b60ff8416815260ff8316602082015260606040820152600062003a37606083018462004405565b60006020828403121562004f9657600080fd5b815167ffffffffffffffff81111562004fae57600080fd5b62002c008482850162004a13565b60006020828403121562004fcf57600080fd5b8151620035e08162004a04565b600082601f83011262004fee57600080fd5b8135602062005001620042758362004787565b82815260059290921b840181019181810190868411156200502157600080fd5b8286015b8481101562004838578035835291830191830162005025565b600082601f8301126200505057600080fd5b8135602062005063620042758362004787565b82815260e092830285018201928282019190878511156200508357600080fd5b8387015b858110156200513a5781818a031215620050a15760008081fd5b620050ab62004164565b8135620050b88162004a04565b8152620050c7828701620040c4565b868201526040620050da818401620040c4565b90820152606082810135620050ef81620040a1565b9082015260806200510283820162004721565b9082015260a06200511583820162004721565b9082015260c062005128838201620040c4565b90820152845292840192810162005087565b5090979650505050505050565b600082601f8301126200515957600080fd5b813560206200516c620042758362004787565b82815260059290921b840181019181810190868411156200518c57600080fd5b8286015b8481101562004838578035620051a681620040a1565b835291830191830162005190565b600080600080600080600060e0888a031215620051d057600080fd5b873567ffffffffffffffff80821115620051e957600080fd5b620051f78b838c0162004fdc565b985060208a01359150808211156200520e57600080fd5b6200521c8b838c016200503e565b975060408a01359150808211156200523357600080fd5b620052418b838c0162005147565b965060608a01359150808211156200525857600080fd5b620052668b838c0162005147565b955060808a01359150808211156200527d57600080fd5b6200528b8b838c01620047ae565b945060a08a0135915080821115620052a257600080fd5b620052b08b838c01620047ae565b935060c08a0135915080821115620052c757600080fd5b50620052d68a828b01620047ae565b91505092959891949750929550565b6bffffffffffffffffffffffff81811683821601908082111562003fe65762003fe6620048da565b6bffffffffffffffffffffffff82811682821603908082111562003fe65762003fe6620048da565b6040815260006200534a604083018662004e0e565b82810360208401526200535f8185876200492f565b9695505050505050565b8051602080830151919081101562004afb577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60209190910360031b1b16919050565b815167ffffffffffffffff811115620053c957620053c962004135565b620053e181620053da845462004aac565b8462004b01565b602080601f831160018114620054375760008415620054005750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b17855562004b4b565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015620054865788860151825594840194600190910190840162005465565b5085821015620054c357878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b805169ffffffffffffffffffff811681146200398e57600080fd5b600080600080600060a086880312156200550757600080fd5b6200551286620054d3565b94506020860151935060408601519250606086015191506200553760808701620054d3565b90509295509295909350565b60008154620055528162004aac565b808552602060018381168015620055725760018114620055ab57620055db565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b8901019550620055db565b866000528260002060005b85811015620055d35781548a8201860152908301908401620055b6565b890184019650505b505050505092915050565b60208152600062002a5b602083018462005543565b600082601f8301126200560d57600080fd5b8151602062005620620042758362004787565b82815260059290921b840181019181810190868411156200564057600080fd5b8286015b8481101562004838578051835291830191830162005644565b6000602082840312156200567057600080fd5b815167ffffffffffffffff808211156200568957600080fd5b9083019061010082860312156200569f57600080fd5b620056a962004190565b8251815260208301516020820152604083015160408201526060830151606082015260808301516080820152620056e360a08401620049a3565b60a082015260c083015182811115620056fb57600080fd5b6200570987828601620055fb565b60c08301525060e0830151828111156200572257600080fd5b620057308782860162004a13565b60e08301525095945050505050565b600081518084526020808501945080840160005b8381101562004dbb5781518752958201959082019060010162005753565b60408152825160408201526020830151606082015260408301516080820152606083015160a0820152608083015160c082015273ffffffffffffffffffffffffffffffffffffffff60a08401511660e0820152600060c0840151610100808185015250620057e46101408401826200573f565b905060e08501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08483030161012085015262005822828262004405565b915050828103602084015262003a37818562005543565b6000602082840312156200584c57600080fd5b5051919050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156200588e576200588e620048da565b500290565b60ff818116838216019081111562002a5e5762002a5e620048da565b828482376000838201600081528351620058ce818360208801620043df565b0195945050505050565b600063ffffffff80831681851681830481118215151615620058fe57620058fe620048da565b02949350505050565b60008060008060008060c087890312156200592157600080fd5b865195506020870151945060408701519350606087015192506080870151915060a087015190509295509295509295565b60008262005989577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfe60c060405234801561001057600080fd5b506040516103ca3803806103ca83398101604081905261002f91610076565b600080546001600160a01b0319166001600160a01b039384161790559181166080521660a0526100b9565b80516001600160a01b038116811461007157600080fd5b919050565b60008060006060848603121561008b57600080fd5b6100948461005a565b92506100a26020850161005a565b91506100b06040850161005a565b90509250925092565b60805160a0516102e76100e36000396000603801526000818160c4015261011701526102e76000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806379188d161461007b578063f00e6a2a146100aa575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e808015610076573d6000f35b3d6000fd5b61008e6100893660046101c1565b6100ee565b6040805192151583526020830191909152015b60405180910390f35b60405173ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001681526020016100a1565b60008054819073ffffffffffffffffffffffffffffffffffffffff16331461011557600080fd5b7f00000000000000000000000000000000000000000000000000000000000000005a91505a61138881101561014957600080fd5b61138881039050856040820482031161016157600080fd5b50803b61016d57600080fd5b6000808551602087016000858af192505a610188908361029a565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080604083850312156101d457600080fd5b82359150602083013567ffffffffffffffff808211156101f357600080fd5b818501915085601f83011261020757600080fd5b81358181111561021957610219610192565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561025f5761025f610192565b8160405282815288602084870101111561027857600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b818103818111156102d4577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b9291505056fea164736f6c6343000810000a307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000810000a",
}

var KeeperRegistryLogicA21ABI = KeeperRegistryLogicA21MetaData.ABI

var KeeperRegistryLogicA21Bin = KeeperRegistryLogicA21MetaData.Bin

func DeployKeeperRegistryLogicA21(auth *bind.TransactOpts, backend bind.ContractBackend, logicB common.Address) (common.Address, *types.Transaction, *KeeperRegistryLogicA21, error) {
	parsed, err := KeeperRegistryLogicA21MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistryLogicA21Bin), backend, logicB)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperRegistryLogicA21{address: address, abi: *parsed, KeeperRegistryLogicA21Caller: KeeperRegistryLogicA21Caller{contract: contract}, KeeperRegistryLogicA21Transactor: KeeperRegistryLogicA21Transactor{contract: contract}, KeeperRegistryLogicA21Filterer: KeeperRegistryLogicA21Filterer{contract: contract}}, nil
}

type KeeperRegistryLogicA21 struct {
	address common.Address
	abi     abi.ABI
	KeeperRegistryLogicA21Caller
	KeeperRegistryLogicA21Transactor
	KeeperRegistryLogicA21Filterer
}

type KeeperRegistryLogicA21Caller struct {
	contract *bind.BoundContract
}

type KeeperRegistryLogicA21Transactor struct {
	contract *bind.BoundContract
}

type KeeperRegistryLogicA21Filterer struct {
	contract *bind.BoundContract
}

type KeeperRegistryLogicA21Session struct {
	Contract     *KeeperRegistryLogicA21
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeeperRegistryLogicA21CallerSession struct {
	Contract *KeeperRegistryLogicA21Caller
	CallOpts bind.CallOpts
}

type KeeperRegistryLogicA21TransactorSession struct {
	Contract     *KeeperRegistryLogicA21Transactor
	TransactOpts bind.TransactOpts
}

type KeeperRegistryLogicA21Raw struct {
	Contract *KeeperRegistryLogicA21
}

type KeeperRegistryLogicA21CallerRaw struct {
	Contract *KeeperRegistryLogicA21Caller
}

type KeeperRegistryLogicA21TransactorRaw struct {
	Contract *KeeperRegistryLogicA21Transactor
}

func NewKeeperRegistryLogicA21(address common.Address, backend bind.ContractBackend) (*KeeperRegistryLogicA21, error) {
	abi, err := abi.JSON(strings.NewReader(KeeperRegistryLogicA21ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeeperRegistryLogicA21(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21{address: address, abi: abi, KeeperRegistryLogicA21Caller: KeeperRegistryLogicA21Caller{contract: contract}, KeeperRegistryLogicA21Transactor: KeeperRegistryLogicA21Transactor{contract: contract}, KeeperRegistryLogicA21Filterer: KeeperRegistryLogicA21Filterer{contract: contract}}, nil
}

func NewKeeperRegistryLogicA21Caller(address common.Address, caller bind.ContractCaller) (*KeeperRegistryLogicA21Caller, error) {
	contract, err := bindKeeperRegistryLogicA21(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21Caller{contract: contract}, nil
}

func NewKeeperRegistryLogicA21Transactor(address common.Address, transactor bind.ContractTransactor) (*KeeperRegistryLogicA21Transactor, error) {
	contract, err := bindKeeperRegistryLogicA21(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21Transactor{contract: contract}, nil
}

func NewKeeperRegistryLogicA21Filterer(address common.Address, filterer bind.ContractFilterer) (*KeeperRegistryLogicA21Filterer, error) {
	contract, err := bindKeeperRegistryLogicA21(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21Filterer{contract: contract}, nil
}

func bindKeeperRegistryLogicA21(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeeperRegistryLogicA21MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistryLogicA21.Contract.KeeperRegistryLogicA21Caller.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.KeeperRegistryLogicA21Transactor.contract.Transfer(opts)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.KeeperRegistryLogicA21Transactor.contract.Transact(opts, method, params...)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistryLogicA21.Contract.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.contract.Transfer(opts)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.contract.Transact(opts, method, params...)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Caller) FallbackTo(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogicA21.contract.Call(opts, &out, "fallbackTo")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Session) FallbackTo() (common.Address, error) {
	return _KeeperRegistryLogicA21.Contract.FallbackTo(&_KeeperRegistryLogicA21.CallOpts)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21CallerSession) FallbackTo() (common.Address, error) {
	return _KeeperRegistryLogicA21.Contract.FallbackTo(&_KeeperRegistryLogicA21.CallOpts)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogicA21.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Session) Owner() (common.Address, error) {
	return _KeeperRegistryLogicA21.Contract.Owner(&_KeeperRegistryLogicA21.CallOpts)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21CallerSession) Owner() (common.Address, error) {
	return _KeeperRegistryLogicA21.Contract.Owner(&_KeeperRegistryLogicA21.CallOpts)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.contract.Transact(opts, "acceptOwnership")
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Session) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.AcceptOwnership(&_KeeperRegistryLogicA21.TransactOpts)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.AcceptOwnership(&_KeeperRegistryLogicA21.TransactOpts)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Transactor) AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.contract.Transact(opts, "addFunds", id, amount)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Session) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.AddFunds(&_KeeperRegistryLogicA21.TransactOpts, id, amount)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21TransactorSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.AddFunds(&_KeeperRegistryLogicA21.TransactOpts, id, amount)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Transactor) CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.contract.Transact(opts, "cancelUpkeep", id)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Session) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.CancelUpkeep(&_KeeperRegistryLogicA21.TransactOpts, id)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21TransactorSession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.CancelUpkeep(&_KeeperRegistryLogicA21.TransactOpts, id)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Transactor) CheckCallback(opts *bind.TransactOpts, id *big.Int, values [][]byte, extraData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.contract.Transact(opts, "checkCallback", id, values, extraData)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Session) CheckCallback(id *big.Int, values [][]byte, extraData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.CheckCallback(&_KeeperRegistryLogicA21.TransactOpts, id, values, extraData)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21TransactorSession) CheckCallback(id *big.Int, values [][]byte, extraData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.CheckCallback(&_KeeperRegistryLogicA21.TransactOpts, id, values, extraData)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Transactor) CheckUpkeep(opts *bind.TransactOpts, id *big.Int, triggerData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.contract.Transact(opts, "checkUpkeep", id, triggerData)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Session) CheckUpkeep(id *big.Int, triggerData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.CheckUpkeep(&_KeeperRegistryLogicA21.TransactOpts, id, triggerData)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21TransactorSession) CheckUpkeep(id *big.Int, triggerData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.CheckUpkeep(&_KeeperRegistryLogicA21.TransactOpts, id, triggerData)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Transactor) CheckUpkeep0(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.contract.Transact(opts, "checkUpkeep0", id)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Session) CheckUpkeep0(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.CheckUpkeep0(&_KeeperRegistryLogicA21.TransactOpts, id)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21TransactorSession) CheckUpkeep0(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.CheckUpkeep0(&_KeeperRegistryLogicA21.TransactOpts, id)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Transactor) ExecuteCallback(opts *bind.TransactOpts, id *big.Int, payload []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.contract.Transact(opts, "executeCallback", id, payload)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Session) ExecuteCallback(id *big.Int, payload []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.ExecuteCallback(&_KeeperRegistryLogicA21.TransactOpts, id, payload)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21TransactorSession) ExecuteCallback(id *big.Int, payload []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.ExecuteCallback(&_KeeperRegistryLogicA21.TransactOpts, id, payload)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Transactor) MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.contract.Transact(opts, "migrateUpkeeps", ids, destination)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Session) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.MigrateUpkeeps(&_KeeperRegistryLogicA21.TransactOpts, ids, destination)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21TransactorSession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.MigrateUpkeeps(&_KeeperRegistryLogicA21.TransactOpts, ids, destination)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Transactor) ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.contract.Transact(opts, "receiveUpkeeps", encodedUpkeeps)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Session) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.ReceiveUpkeeps(&_KeeperRegistryLogicA21.TransactOpts, encodedUpkeeps)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21TransactorSession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.ReceiveUpkeeps(&_KeeperRegistryLogicA21.TransactOpts, encodedUpkeeps)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Transactor) RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.contract.Transact(opts, "registerUpkeep", target, gasLimit, admin, triggerType, checkData, triggerConfig, offchainConfig)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Session) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.RegisterUpkeep(&_KeeperRegistryLogicA21.TransactOpts, target, gasLimit, admin, triggerType, checkData, triggerConfig, offchainConfig)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21TransactorSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.RegisterUpkeep(&_KeeperRegistryLogicA21.TransactOpts, target, gasLimit, admin, triggerType, checkData, triggerConfig, offchainConfig)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Transactor) RegisterUpkeep0(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.contract.Transact(opts, "registerUpkeep0", target, gasLimit, admin, checkData, offchainConfig)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Session) RegisterUpkeep0(target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.RegisterUpkeep0(&_KeeperRegistryLogicA21.TransactOpts, target, gasLimit, admin, checkData, offchainConfig)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21TransactorSession) RegisterUpkeep0(target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.RegisterUpkeep0(&_KeeperRegistryLogicA21.TransactOpts, target, gasLimit, admin, checkData, offchainConfig)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Transactor) SetUpkeepTriggerConfig(opts *bind.TransactOpts, id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.contract.Transact(opts, "setUpkeepTriggerConfig", id, triggerConfig)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Session) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.SetUpkeepTriggerConfig(&_KeeperRegistryLogicA21.TransactOpts, id, triggerConfig)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21TransactorSession) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.SetUpkeepTriggerConfig(&_KeeperRegistryLogicA21.TransactOpts, id, triggerConfig)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.contract.Transact(opts, "transferOwnership", to)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.TransferOwnership(&_KeeperRegistryLogicA21.TransactOpts, to)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.TransferOwnership(&_KeeperRegistryLogicA21.TransactOpts, to)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Transactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.contract.RawTransact(opts, calldata)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Session) Fallback(calldata []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.Fallback(&_KeeperRegistryLogicA21.TransactOpts, calldata)
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21TransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA21.Contract.Fallback(&_KeeperRegistryLogicA21.TransactOpts, calldata)
}

type KeeperRegistryLogicA21AdminPrivilegeConfigSetIterator struct {
	Event *KeeperRegistryLogicA21AdminPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21AdminPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21AdminPrivilegeConfigSet)
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
		it.Event = new(KeeperRegistryLogicA21AdminPrivilegeConfigSet)
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

func (it *KeeperRegistryLogicA21AdminPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21AdminPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21AdminPrivilegeConfigSet struct {
	Admin           common.Address
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*KeeperRegistryLogicA21AdminPrivilegeConfigSetIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21AdminPrivilegeConfigSetIterator{contract: _KeeperRegistryLogicA21.contract, event: "AdminPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21AdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21AdminPrivilegeConfigSet)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParseAdminPrivilegeConfigSet(log types.Log) (*KeeperRegistryLogicA21AdminPrivilegeConfigSet, error) {
	event := new(KeeperRegistryLogicA21AdminPrivilegeConfigSet)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21CancelledUpkeepReportIterator struct {
	Event *KeeperRegistryLogicA21CancelledUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21CancelledUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21CancelledUpkeepReport)
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
		it.Event = new(KeeperRegistryLogicA21CancelledUpkeepReport)
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

func (it *KeeperRegistryLogicA21CancelledUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21CancelledUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21CancelledUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21CancelledUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21CancelledUpkeepReportIterator{contract: _KeeperRegistryLogicA21.contract, event: "CancelledUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21CancelledUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21CancelledUpkeepReport)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParseCancelledUpkeepReport(log types.Log) (*KeeperRegistryLogicA21CancelledUpkeepReport, error) {
	event := new(KeeperRegistryLogicA21CancelledUpkeepReport)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21DedupKeyAddedIterator struct {
	Event *KeeperRegistryLogicA21DedupKeyAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21DedupKeyAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21DedupKeyAdded)
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
		it.Event = new(KeeperRegistryLogicA21DedupKeyAdded)
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

func (it *KeeperRegistryLogicA21DedupKeyAddedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21DedupKeyAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21DedupKeyAdded struct {
	DedupKey [32]byte
	Raw      types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*KeeperRegistryLogicA21DedupKeyAddedIterator, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21DedupKeyAddedIterator{contract: _KeeperRegistryLogicA21.contract, event: "DedupKeyAdded", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21DedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21DedupKeyAdded)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParseDedupKeyAdded(log types.Log) (*KeeperRegistryLogicA21DedupKeyAdded, error) {
	event := new(KeeperRegistryLogicA21DedupKeyAdded)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21FundsAddedIterator struct {
	Event *KeeperRegistryLogicA21FundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21FundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21FundsAdded)
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
		it.Event = new(KeeperRegistryLogicA21FundsAdded)
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

func (it *KeeperRegistryLogicA21FundsAddedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21FundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21FundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistryLogicA21FundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21FundsAddedIterator{contract: _KeeperRegistryLogicA21.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21FundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21FundsAdded)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "FundsAdded", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParseFundsAdded(log types.Log) (*KeeperRegistryLogicA21FundsAdded, error) {
	event := new(KeeperRegistryLogicA21FundsAdded)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21FundsWithdrawnIterator struct {
	Event *KeeperRegistryLogicA21FundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21FundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21FundsWithdrawn)
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
		it.Event = new(KeeperRegistryLogicA21FundsWithdrawn)
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

func (it *KeeperRegistryLogicA21FundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21FundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21FundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21FundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21FundsWithdrawnIterator{contract: _KeeperRegistryLogicA21.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21FundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21FundsWithdrawn)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParseFundsWithdrawn(log types.Log) (*KeeperRegistryLogicA21FundsWithdrawn, error) {
	event := new(KeeperRegistryLogicA21FundsWithdrawn)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21InsufficientFundsUpkeepReportIterator struct {
	Event *KeeperRegistryLogicA21InsufficientFundsUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21InsufficientFundsUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21InsufficientFundsUpkeepReport)
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
		it.Event = new(KeeperRegistryLogicA21InsufficientFundsUpkeepReport)
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

func (it *KeeperRegistryLogicA21InsufficientFundsUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21InsufficientFundsUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21InsufficientFundsUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21InsufficientFundsUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21InsufficientFundsUpkeepReportIterator{contract: _KeeperRegistryLogicA21.contract, event: "InsufficientFundsUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21InsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21InsufficientFundsUpkeepReport)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParseInsufficientFundsUpkeepReport(log types.Log) (*KeeperRegistryLogicA21InsufficientFundsUpkeepReport, error) {
	event := new(KeeperRegistryLogicA21InsufficientFundsUpkeepReport)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21OwnerFundsWithdrawnIterator struct {
	Event *KeeperRegistryLogicA21OwnerFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21OwnerFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21OwnerFundsWithdrawn)
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
		it.Event = new(KeeperRegistryLogicA21OwnerFundsWithdrawn)
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

func (it *KeeperRegistryLogicA21OwnerFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21OwnerFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21OwnerFundsWithdrawn struct {
	Amount *big.Int
	Raw    types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*KeeperRegistryLogicA21OwnerFundsWithdrawnIterator, error) {

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21OwnerFundsWithdrawnIterator{contract: _KeeperRegistryLogicA21.contract, event: "OwnerFundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21OwnerFundsWithdrawn) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21OwnerFundsWithdrawn)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParseOwnerFundsWithdrawn(log types.Log) (*KeeperRegistryLogicA21OwnerFundsWithdrawn, error) {
	event := new(KeeperRegistryLogicA21OwnerFundsWithdrawn)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21OwnershipTransferRequestedIterator struct {
	Event *KeeperRegistryLogicA21OwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21OwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21OwnershipTransferRequested)
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
		it.Event = new(KeeperRegistryLogicA21OwnershipTransferRequested)
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

func (it *KeeperRegistryLogicA21OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryLogicA21OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21OwnershipTransferRequestedIterator{contract: _KeeperRegistryLogicA21.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21OwnershipTransferRequested)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryLogicA21OwnershipTransferRequested, error) {
	event := new(KeeperRegistryLogicA21OwnershipTransferRequested)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21OwnershipTransferredIterator struct {
	Event *KeeperRegistryLogicA21OwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21OwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21OwnershipTransferred)
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
		it.Event = new(KeeperRegistryLogicA21OwnershipTransferred)
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

func (it *KeeperRegistryLogicA21OwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryLogicA21OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21OwnershipTransferredIterator{contract: _KeeperRegistryLogicA21.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21OwnershipTransferred)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParseOwnershipTransferred(log types.Log) (*KeeperRegistryLogicA21OwnershipTransferred, error) {
	event := new(KeeperRegistryLogicA21OwnershipTransferred)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21PausedIterator struct {
	Event *KeeperRegistryLogicA21Paused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21PausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21Paused)
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
		it.Event = new(KeeperRegistryLogicA21Paused)
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

func (it *KeeperRegistryLogicA21PausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21PausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21Paused struct {
	Account common.Address
	Raw     types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterPaused(opts *bind.FilterOpts) (*KeeperRegistryLogicA21PausedIterator, error) {

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21PausedIterator{contract: _KeeperRegistryLogicA21.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21Paused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21Paused)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "Paused", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParsePaused(log types.Log) (*KeeperRegistryLogicA21Paused, error) {
	event := new(KeeperRegistryLogicA21Paused)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21PayeesUpdatedIterator struct {
	Event *KeeperRegistryLogicA21PayeesUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21PayeesUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21PayeesUpdated)
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
		it.Event = new(KeeperRegistryLogicA21PayeesUpdated)
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

func (it *KeeperRegistryLogicA21PayeesUpdatedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21PayeesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21PayeesUpdated struct {
	Transmitters []common.Address
	Payees       []common.Address
	Raw          types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterPayeesUpdated(opts *bind.FilterOpts) (*KeeperRegistryLogicA21PayeesUpdatedIterator, error) {

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21PayeesUpdatedIterator{contract: _KeeperRegistryLogicA21.contract, event: "PayeesUpdated", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21PayeesUpdated) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21PayeesUpdated)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParsePayeesUpdated(log types.Log) (*KeeperRegistryLogicA21PayeesUpdated, error) {
	event := new(KeeperRegistryLogicA21PayeesUpdated)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21PayeeshipTransferRequestedIterator struct {
	Event *KeeperRegistryLogicA21PayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21PayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21PayeeshipTransferRequested)
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
		it.Event = new(KeeperRegistryLogicA21PayeeshipTransferRequested)
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

func (it *KeeperRegistryLogicA21PayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21PayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21PayeeshipTransferRequested struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogicA21PayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21PayeeshipTransferRequestedIterator{contract: _KeeperRegistryLogicA21.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21PayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21PayeeshipTransferRequested)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistryLogicA21PayeeshipTransferRequested, error) {
	event := new(KeeperRegistryLogicA21PayeeshipTransferRequested)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21PayeeshipTransferredIterator struct {
	Event *KeeperRegistryLogicA21PayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21PayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21PayeeshipTransferred)
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
		it.Event = new(KeeperRegistryLogicA21PayeeshipTransferred)
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

func (it *KeeperRegistryLogicA21PayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21PayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21PayeeshipTransferred struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogicA21PayeeshipTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21PayeeshipTransferredIterator{contract: _KeeperRegistryLogicA21.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21PayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21PayeeshipTransferred)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParsePayeeshipTransferred(log types.Log) (*KeeperRegistryLogicA21PayeeshipTransferred, error) {
	event := new(KeeperRegistryLogicA21PayeeshipTransferred)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21PaymentWithdrawnIterator struct {
	Event *KeeperRegistryLogicA21PaymentWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21PaymentWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21PaymentWithdrawn)
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
		it.Event = new(KeeperRegistryLogicA21PaymentWithdrawn)
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

func (it *KeeperRegistryLogicA21PaymentWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21PaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21PaymentWithdrawn struct {
	Transmitter common.Address
	Amount      *big.Int
	To          common.Address
	Payee       common.Address
	Raw         types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistryLogicA21PaymentWithdrawnIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21PaymentWithdrawnIterator{contract: _KeeperRegistryLogicA21.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21PaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21PaymentWithdrawn)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParsePaymentWithdrawn(log types.Log) (*KeeperRegistryLogicA21PaymentWithdrawn, error) {
	event := new(KeeperRegistryLogicA21PaymentWithdrawn)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21ReorgedUpkeepReportIterator struct {
	Event *KeeperRegistryLogicA21ReorgedUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21ReorgedUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21ReorgedUpkeepReport)
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
		it.Event = new(KeeperRegistryLogicA21ReorgedUpkeepReport)
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

func (it *KeeperRegistryLogicA21ReorgedUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21ReorgedUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21ReorgedUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21ReorgedUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21ReorgedUpkeepReportIterator{contract: _KeeperRegistryLogicA21.contract, event: "ReorgedUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21ReorgedUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21ReorgedUpkeepReport)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParseReorgedUpkeepReport(log types.Log) (*KeeperRegistryLogicA21ReorgedUpkeepReport, error) {
	event := new(KeeperRegistryLogicA21ReorgedUpkeepReport)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21StaleUpkeepReportIterator struct {
	Event *KeeperRegistryLogicA21StaleUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21StaleUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21StaleUpkeepReport)
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
		it.Event = new(KeeperRegistryLogicA21StaleUpkeepReport)
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

func (it *KeeperRegistryLogicA21StaleUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21StaleUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21StaleUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21StaleUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21StaleUpkeepReportIterator{contract: _KeeperRegistryLogicA21.contract, event: "StaleUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21StaleUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21StaleUpkeepReport)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParseStaleUpkeepReport(log types.Log) (*KeeperRegistryLogicA21StaleUpkeepReport, error) {
	event := new(KeeperRegistryLogicA21StaleUpkeepReport)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21UnpausedIterator struct {
	Event *KeeperRegistryLogicA21Unpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21UnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21Unpaused)
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
		it.Event = new(KeeperRegistryLogicA21Unpaused)
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

func (it *KeeperRegistryLogicA21UnpausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21UnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21Unpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistryLogicA21UnpausedIterator, error) {

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21UnpausedIterator{contract: _KeeperRegistryLogicA21.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21Unpaused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21Unpaused)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParseUnpaused(log types.Log) (*KeeperRegistryLogicA21Unpaused, error) {
	event := new(KeeperRegistryLogicA21Unpaused)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21UpkeepAdminTransferRequestedIterator struct {
	Event *KeeperRegistryLogicA21UpkeepAdminTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21UpkeepAdminTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21UpkeepAdminTransferRequested)
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
		it.Event = new(KeeperRegistryLogicA21UpkeepAdminTransferRequested)
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

func (it *KeeperRegistryLogicA21UpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21UpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21UpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogicA21UpkeepAdminTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21UpkeepAdminTransferRequestedIterator{contract: _KeeperRegistryLogicA21.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21UpkeepAdminTransferRequested)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParseUpkeepAdminTransferRequested(log types.Log) (*KeeperRegistryLogicA21UpkeepAdminTransferRequested, error) {
	event := new(KeeperRegistryLogicA21UpkeepAdminTransferRequested)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21UpkeepAdminTransferredIterator struct {
	Event *KeeperRegistryLogicA21UpkeepAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21UpkeepAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21UpkeepAdminTransferred)
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
		it.Event = new(KeeperRegistryLogicA21UpkeepAdminTransferred)
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

func (it *KeeperRegistryLogicA21UpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21UpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21UpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogicA21UpkeepAdminTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21UpkeepAdminTransferredIterator{contract: _KeeperRegistryLogicA21.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21UpkeepAdminTransferred)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParseUpkeepAdminTransferred(log types.Log) (*KeeperRegistryLogicA21UpkeepAdminTransferred, error) {
	event := new(KeeperRegistryLogicA21UpkeepAdminTransferred)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21UpkeepCanceledIterator struct {
	Event *KeeperRegistryLogicA21UpkeepCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21UpkeepCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21UpkeepCanceled)
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
		it.Event = new(KeeperRegistryLogicA21UpkeepCanceled)
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

func (it *KeeperRegistryLogicA21UpkeepCanceledIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21UpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21UpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistryLogicA21UpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21UpkeepCanceledIterator{contract: _KeeperRegistryLogicA21.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21UpkeepCanceled)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParseUpkeepCanceled(log types.Log) (*KeeperRegistryLogicA21UpkeepCanceled, error) {
	event := new(KeeperRegistryLogicA21UpkeepCanceled)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21UpkeepCheckDataSetIterator struct {
	Event *KeeperRegistryLogicA21UpkeepCheckDataSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21UpkeepCheckDataSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21UpkeepCheckDataSet)
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
		it.Event = new(KeeperRegistryLogicA21UpkeepCheckDataSet)
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

func (it *KeeperRegistryLogicA21UpkeepCheckDataSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21UpkeepCheckDataSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21UpkeepCheckDataSet struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21UpkeepCheckDataSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21UpkeepCheckDataSetIterator{contract: _KeeperRegistryLogicA21.contract, event: "UpkeepCheckDataSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepCheckDataSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21UpkeepCheckDataSet)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParseUpkeepCheckDataSet(log types.Log) (*KeeperRegistryLogicA21UpkeepCheckDataSet, error) {
	event := new(KeeperRegistryLogicA21UpkeepCheckDataSet)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21UpkeepGasLimitSetIterator struct {
	Event *KeeperRegistryLogicA21UpkeepGasLimitSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21UpkeepGasLimitSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21UpkeepGasLimitSet)
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
		it.Event = new(KeeperRegistryLogicA21UpkeepGasLimitSet)
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

func (it *KeeperRegistryLogicA21UpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21UpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21UpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21UpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21UpkeepGasLimitSetIterator{contract: _KeeperRegistryLogicA21.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21UpkeepGasLimitSet)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParseUpkeepGasLimitSet(log types.Log) (*KeeperRegistryLogicA21UpkeepGasLimitSet, error) {
	event := new(KeeperRegistryLogicA21UpkeepGasLimitSet)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21UpkeepMigratedIterator struct {
	Event *KeeperRegistryLogicA21UpkeepMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21UpkeepMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21UpkeepMigrated)
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
		it.Event = new(KeeperRegistryLogicA21UpkeepMigrated)
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

func (it *KeeperRegistryLogicA21UpkeepMigratedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21UpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21UpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21UpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21UpkeepMigratedIterator{contract: _KeeperRegistryLogicA21.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21UpkeepMigrated)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParseUpkeepMigrated(log types.Log) (*KeeperRegistryLogicA21UpkeepMigrated, error) {
	event := new(KeeperRegistryLogicA21UpkeepMigrated)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21UpkeepOffchainConfigSetIterator struct {
	Event *KeeperRegistryLogicA21UpkeepOffchainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21UpkeepOffchainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21UpkeepOffchainConfigSet)
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
		it.Event = new(KeeperRegistryLogicA21UpkeepOffchainConfigSet)
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

func (it *KeeperRegistryLogicA21UpkeepOffchainConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21UpkeepOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21UpkeepOffchainConfigSet struct {
	Id             *big.Int
	OffchainConfig []byte
	Raw            types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21UpkeepOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21UpkeepOffchainConfigSetIterator{contract: _KeeperRegistryLogicA21.contract, event: "UpkeepOffchainConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21UpkeepOffchainConfigSet)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParseUpkeepOffchainConfigSet(log types.Log) (*KeeperRegistryLogicA21UpkeepOffchainConfigSet, error) {
	event := new(KeeperRegistryLogicA21UpkeepOffchainConfigSet)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21UpkeepPausedIterator struct {
	Event *KeeperRegistryLogicA21UpkeepPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21UpkeepPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21UpkeepPaused)
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
		it.Event = new(KeeperRegistryLogicA21UpkeepPaused)
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

func (it *KeeperRegistryLogicA21UpkeepPausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21UpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21UpkeepPaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21UpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21UpkeepPausedIterator{contract: _KeeperRegistryLogicA21.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21UpkeepPaused)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParseUpkeepPaused(log types.Log) (*KeeperRegistryLogicA21UpkeepPaused, error) {
	event := new(KeeperRegistryLogicA21UpkeepPaused)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21UpkeepPerformedIterator struct {
	Event *KeeperRegistryLogicA21UpkeepPerformed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21UpkeepPerformedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21UpkeepPerformed)
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
		it.Event = new(KeeperRegistryLogicA21UpkeepPerformed)
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

func (it *KeeperRegistryLogicA21UpkeepPerformedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21UpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21UpkeepPerformed struct {
	Id           *big.Int
	Success      bool
	TotalPayment *big.Int
	GasUsed      *big.Int
	GasOverhead  *big.Int
	Trigger      []byte
	Raw          types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*KeeperRegistryLogicA21UpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21UpkeepPerformedIterator{contract: _KeeperRegistryLogicA21.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21UpkeepPerformed)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParseUpkeepPerformed(log types.Log) (*KeeperRegistryLogicA21UpkeepPerformed, error) {
	event := new(KeeperRegistryLogicA21UpkeepPerformed)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21UpkeepPrivilegeConfigSetIterator struct {
	Event *KeeperRegistryLogicA21UpkeepPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21UpkeepPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21UpkeepPrivilegeConfigSet)
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
		it.Event = new(KeeperRegistryLogicA21UpkeepPrivilegeConfigSet)
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

func (it *KeeperRegistryLogicA21UpkeepPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21UpkeepPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21UpkeepPrivilegeConfigSet struct {
	Id              *big.Int
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21UpkeepPrivilegeConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21UpkeepPrivilegeConfigSetIterator{contract: _KeeperRegistryLogicA21.contract, event: "UpkeepPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21UpkeepPrivilegeConfigSet)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParseUpkeepPrivilegeConfigSet(log types.Log) (*KeeperRegistryLogicA21UpkeepPrivilegeConfigSet, error) {
	event := new(KeeperRegistryLogicA21UpkeepPrivilegeConfigSet)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21UpkeepReceivedIterator struct {
	Event *KeeperRegistryLogicA21UpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21UpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21UpkeepReceived)
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
		it.Event = new(KeeperRegistryLogicA21UpkeepReceived)
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

func (it *KeeperRegistryLogicA21UpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21UpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21UpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21UpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21UpkeepReceivedIterator{contract: _KeeperRegistryLogicA21.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21UpkeepReceived)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParseUpkeepReceived(log types.Log) (*KeeperRegistryLogicA21UpkeepReceived, error) {
	event := new(KeeperRegistryLogicA21UpkeepReceived)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21UpkeepRegisteredIterator struct {
	Event *KeeperRegistryLogicA21UpkeepRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21UpkeepRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21UpkeepRegistered)
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
		it.Event = new(KeeperRegistryLogicA21UpkeepRegistered)
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

func (it *KeeperRegistryLogicA21UpkeepRegisteredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21UpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21UpkeepRegistered struct {
	Id         *big.Int
	PerformGas uint32
	Admin      common.Address
	Raw        types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21UpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21UpkeepRegisteredIterator{contract: _KeeperRegistryLogicA21.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21UpkeepRegistered)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParseUpkeepRegistered(log types.Log) (*KeeperRegistryLogicA21UpkeepRegistered, error) {
	event := new(KeeperRegistryLogicA21UpkeepRegistered)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21UpkeepTriggerConfigSetIterator struct {
	Event *KeeperRegistryLogicA21UpkeepTriggerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21UpkeepTriggerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21UpkeepTriggerConfigSet)
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
		it.Event = new(KeeperRegistryLogicA21UpkeepTriggerConfigSet)
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

func (it *KeeperRegistryLogicA21UpkeepTriggerConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21UpkeepTriggerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21UpkeepTriggerConfigSet struct {
	Id            *big.Int
	TriggerConfig []byte
	Raw           types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21UpkeepTriggerConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21UpkeepTriggerConfigSetIterator{contract: _KeeperRegistryLogicA21.contract, event: "UpkeepTriggerConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21UpkeepTriggerConfigSet)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParseUpkeepTriggerConfigSet(log types.Log) (*KeeperRegistryLogicA21UpkeepTriggerConfigSet, error) {
	event := new(KeeperRegistryLogicA21UpkeepTriggerConfigSet)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicA21UpkeepUnpausedIterator struct {
	Event *KeeperRegistryLogicA21UpkeepUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicA21UpkeepUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicA21UpkeepUnpaused)
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
		it.Event = new(KeeperRegistryLogicA21UpkeepUnpaused)
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

func (it *KeeperRegistryLogicA21UpkeepUnpausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicA21UpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicA21UpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21UpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA21UpkeepUnpausedIterator{contract: _KeeperRegistryLogicA21.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA21.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicA21UpkeepUnpaused)
				if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21Filterer) ParseUpkeepUnpaused(log types.Log) (*KeeperRegistryLogicA21UpkeepUnpaused, error) {
	event := new(KeeperRegistryLogicA21UpkeepUnpaused)
	if err := _KeeperRegistryLogicA21.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _KeeperRegistryLogicA21.abi.Events["AdminPrivilegeConfigSet"].ID:
		return _KeeperRegistryLogicA21.ParseAdminPrivilegeConfigSet(log)
	case _KeeperRegistryLogicA21.abi.Events["CancelledUpkeepReport"].ID:
		return _KeeperRegistryLogicA21.ParseCancelledUpkeepReport(log)
	case _KeeperRegistryLogicA21.abi.Events["DedupKeyAdded"].ID:
		return _KeeperRegistryLogicA21.ParseDedupKeyAdded(log)
	case _KeeperRegistryLogicA21.abi.Events["FundsAdded"].ID:
		return _KeeperRegistryLogicA21.ParseFundsAdded(log)
	case _KeeperRegistryLogicA21.abi.Events["FundsWithdrawn"].ID:
		return _KeeperRegistryLogicA21.ParseFundsWithdrawn(log)
	case _KeeperRegistryLogicA21.abi.Events["InsufficientFundsUpkeepReport"].ID:
		return _KeeperRegistryLogicA21.ParseInsufficientFundsUpkeepReport(log)
	case _KeeperRegistryLogicA21.abi.Events["OwnerFundsWithdrawn"].ID:
		return _KeeperRegistryLogicA21.ParseOwnerFundsWithdrawn(log)
	case _KeeperRegistryLogicA21.abi.Events["OwnershipTransferRequested"].ID:
		return _KeeperRegistryLogicA21.ParseOwnershipTransferRequested(log)
	case _KeeperRegistryLogicA21.abi.Events["OwnershipTransferred"].ID:
		return _KeeperRegistryLogicA21.ParseOwnershipTransferred(log)
	case _KeeperRegistryLogicA21.abi.Events["Paused"].ID:
		return _KeeperRegistryLogicA21.ParsePaused(log)
	case _KeeperRegistryLogicA21.abi.Events["PayeesUpdated"].ID:
		return _KeeperRegistryLogicA21.ParsePayeesUpdated(log)
	case _KeeperRegistryLogicA21.abi.Events["PayeeshipTransferRequested"].ID:
		return _KeeperRegistryLogicA21.ParsePayeeshipTransferRequested(log)
	case _KeeperRegistryLogicA21.abi.Events["PayeeshipTransferred"].ID:
		return _KeeperRegistryLogicA21.ParsePayeeshipTransferred(log)
	case _KeeperRegistryLogicA21.abi.Events["PaymentWithdrawn"].ID:
		return _KeeperRegistryLogicA21.ParsePaymentWithdrawn(log)
	case _KeeperRegistryLogicA21.abi.Events["ReorgedUpkeepReport"].ID:
		return _KeeperRegistryLogicA21.ParseReorgedUpkeepReport(log)
	case _KeeperRegistryLogicA21.abi.Events["StaleUpkeepReport"].ID:
		return _KeeperRegistryLogicA21.ParseStaleUpkeepReport(log)
	case _KeeperRegistryLogicA21.abi.Events["Unpaused"].ID:
		return _KeeperRegistryLogicA21.ParseUnpaused(log)
	case _KeeperRegistryLogicA21.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _KeeperRegistryLogicA21.ParseUpkeepAdminTransferRequested(log)
	case _KeeperRegistryLogicA21.abi.Events["UpkeepAdminTransferred"].ID:
		return _KeeperRegistryLogicA21.ParseUpkeepAdminTransferred(log)
	case _KeeperRegistryLogicA21.abi.Events["UpkeepCanceled"].ID:
		return _KeeperRegistryLogicA21.ParseUpkeepCanceled(log)
	case _KeeperRegistryLogicA21.abi.Events["UpkeepCheckDataSet"].ID:
		return _KeeperRegistryLogicA21.ParseUpkeepCheckDataSet(log)
	case _KeeperRegistryLogicA21.abi.Events["UpkeepGasLimitSet"].ID:
		return _KeeperRegistryLogicA21.ParseUpkeepGasLimitSet(log)
	case _KeeperRegistryLogicA21.abi.Events["UpkeepMigrated"].ID:
		return _KeeperRegistryLogicA21.ParseUpkeepMigrated(log)
	case _KeeperRegistryLogicA21.abi.Events["UpkeepOffchainConfigSet"].ID:
		return _KeeperRegistryLogicA21.ParseUpkeepOffchainConfigSet(log)
	case _KeeperRegistryLogicA21.abi.Events["UpkeepPaused"].ID:
		return _KeeperRegistryLogicA21.ParseUpkeepPaused(log)
	case _KeeperRegistryLogicA21.abi.Events["UpkeepPerformed"].ID:
		return _KeeperRegistryLogicA21.ParseUpkeepPerformed(log)
	case _KeeperRegistryLogicA21.abi.Events["UpkeepPrivilegeConfigSet"].ID:
		return _KeeperRegistryLogicA21.ParseUpkeepPrivilegeConfigSet(log)
	case _KeeperRegistryLogicA21.abi.Events["UpkeepReceived"].ID:
		return _KeeperRegistryLogicA21.ParseUpkeepReceived(log)
	case _KeeperRegistryLogicA21.abi.Events["UpkeepRegistered"].ID:
		return _KeeperRegistryLogicA21.ParseUpkeepRegistered(log)
	case _KeeperRegistryLogicA21.abi.Events["UpkeepTriggerConfigSet"].ID:
		return _KeeperRegistryLogicA21.ParseUpkeepTriggerConfigSet(log)
	case _KeeperRegistryLogicA21.abi.Events["UpkeepUnpaused"].ID:
		return _KeeperRegistryLogicA21.ParseUpkeepUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeeperRegistryLogicA21AdminPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x7c44b4eb59ee7873514e7e43e7718c269d872965938b288aa143befca62f99d2")
}

func (KeeperRegistryLogicA21CancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd5636")
}

func (KeeperRegistryLogicA21DedupKeyAdded) Topic() common.Hash {
	return common.HexToHash("0xa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f2")
}

func (KeeperRegistryLogicA21FundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (KeeperRegistryLogicA21FundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (KeeperRegistryLogicA21InsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e02")
}

func (KeeperRegistryLogicA21OwnerFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1")
}

func (KeeperRegistryLogicA21OwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (KeeperRegistryLogicA21OwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (KeeperRegistryLogicA21Paused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (KeeperRegistryLogicA21PayeesUpdated) Topic() common.Hash {
	return common.HexToHash("0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725")
}

func (KeeperRegistryLogicA21PayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (KeeperRegistryLogicA21PayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (KeeperRegistryLogicA21PaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (KeeperRegistryLogicA21ReorgedUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301")
}

func (KeeperRegistryLogicA21StaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8")
}

func (KeeperRegistryLogicA21Unpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (KeeperRegistryLogicA21UpkeepAdminTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35")
}

func (KeeperRegistryLogicA21UpkeepAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c")
}

func (KeeperRegistryLogicA21UpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (KeeperRegistryLogicA21UpkeepCheckDataSet) Topic() common.Hash {
	return common.HexToHash("0xcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d")
}

func (KeeperRegistryLogicA21UpkeepGasLimitSet) Topic() common.Hash {
	return common.HexToHash("0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c")
}

func (KeeperRegistryLogicA21UpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (KeeperRegistryLogicA21UpkeepOffchainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850")
}

func (KeeperRegistryLogicA21UpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (KeeperRegistryLogicA21UpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b")
}

func (KeeperRegistryLogicA21UpkeepPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae7769")
}

func (KeeperRegistryLogicA21UpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (KeeperRegistryLogicA21UpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (KeeperRegistryLogicA21UpkeepTriggerConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664")
}

func (KeeperRegistryLogicA21UpkeepUnpaused) Topic() common.Hash {
	return common.HexToHash("0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456")
}

func (_KeeperRegistryLogicA21 *KeeperRegistryLogicA21) Address() common.Address {
	return _KeeperRegistryLogicA21.address
}

type KeeperRegistryLogicA21Interface interface {
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

	FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*KeeperRegistryLogicA21AdminPrivilegeConfigSetIterator, error)

	WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21AdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error)

	ParseAdminPrivilegeConfigSet(log types.Log) (*KeeperRegistryLogicA21AdminPrivilegeConfigSet, error)

	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21CancelledUpkeepReportIterator, error)

	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21CancelledUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseCancelledUpkeepReport(log types.Log) (*KeeperRegistryLogicA21CancelledUpkeepReport, error)

	FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*KeeperRegistryLogicA21DedupKeyAddedIterator, error)

	WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21DedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error)

	ParseDedupKeyAdded(log types.Log) (*KeeperRegistryLogicA21DedupKeyAdded, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistryLogicA21FundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21FundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*KeeperRegistryLogicA21FundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21FundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21FundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*KeeperRegistryLogicA21FundsWithdrawn, error)

	FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21InsufficientFundsUpkeepReportIterator, error)

	WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21InsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseInsufficientFundsUpkeepReport(log types.Log) (*KeeperRegistryLogicA21InsufficientFundsUpkeepReport, error)

	FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*KeeperRegistryLogicA21OwnerFundsWithdrawnIterator, error)

	WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21OwnerFundsWithdrawn) (event.Subscription, error)

	ParseOwnerFundsWithdrawn(log types.Log) (*KeeperRegistryLogicA21OwnerFundsWithdrawn, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryLogicA21OwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryLogicA21OwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryLogicA21OwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*KeeperRegistryLogicA21OwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*KeeperRegistryLogicA21PausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21Paused) (event.Subscription, error)

	ParsePaused(log types.Log) (*KeeperRegistryLogicA21Paused, error)

	FilterPayeesUpdated(opts *bind.FilterOpts) (*KeeperRegistryLogicA21PayeesUpdatedIterator, error)

	WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21PayeesUpdated) (event.Subscription, error)

	ParsePayeesUpdated(log types.Log) (*KeeperRegistryLogicA21PayeesUpdated, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogicA21PayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21PayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistryLogicA21PayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogicA21PayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21PayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*KeeperRegistryLogicA21PayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistryLogicA21PaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21PaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*KeeperRegistryLogicA21PaymentWithdrawn, error)

	FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21ReorgedUpkeepReportIterator, error)

	WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21ReorgedUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseReorgedUpkeepReport(log types.Log) (*KeeperRegistryLogicA21ReorgedUpkeepReport, error)

	FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21StaleUpkeepReportIterator, error)

	WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21StaleUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseStaleUpkeepReport(log types.Log) (*KeeperRegistryLogicA21StaleUpkeepReport, error)

	FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistryLogicA21UnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21Unpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*KeeperRegistryLogicA21Unpaused, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogicA21UpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*KeeperRegistryLogicA21UpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogicA21UpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*KeeperRegistryLogicA21UpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistryLogicA21UpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*KeeperRegistryLogicA21UpkeepCanceled, error)

	FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21UpkeepCheckDataSetIterator, error)

	WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepCheckDataSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCheckDataSet(log types.Log) (*KeeperRegistryLogicA21UpkeepCheckDataSet, error)

	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21UpkeepGasLimitSetIterator, error)

	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepGasLimitSet(log types.Log) (*KeeperRegistryLogicA21UpkeepGasLimitSet, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21UpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*KeeperRegistryLogicA21UpkeepMigrated, error)

	FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21UpkeepOffchainConfigSetIterator, error)

	WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepOffchainConfigSet(log types.Log) (*KeeperRegistryLogicA21UpkeepOffchainConfigSet, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21UpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*KeeperRegistryLogicA21UpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*KeeperRegistryLogicA21UpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*KeeperRegistryLogicA21UpkeepPerformed, error)

	FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21UpkeepPrivilegeConfigSetIterator, error)

	WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPrivilegeConfigSet(log types.Log) (*KeeperRegistryLogicA21UpkeepPrivilegeConfigSet, error)

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21UpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*KeeperRegistryLogicA21UpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21UpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*KeeperRegistryLogicA21UpkeepRegistered, error)

	FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21UpkeepTriggerConfigSetIterator, error)

	WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepTriggerConfigSet(log types.Log) (*KeeperRegistryLogicA21UpkeepTriggerConfigSet, error)

	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicA21UpkeepUnpausedIterator, error)

	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicA21UpkeepUnpaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepUnpaused(log types.Log) (*KeeperRegistryLogicA21UpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
