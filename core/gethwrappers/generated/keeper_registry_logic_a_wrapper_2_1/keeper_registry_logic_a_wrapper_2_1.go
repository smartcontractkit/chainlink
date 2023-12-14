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

var KeeperRegistryLogicAMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractKeeperRegistryLogicB2_1\",\"name\":\"logicB\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTrigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTriggerType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"AdminPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"DedupKeyAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepTriggerConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumKeeperRegistryBase2_1.UpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"triggerData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumKeeperRegistryBase2_1.UpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumKeeperRegistryBase2_1.UpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"executeCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumKeeperRegistryBase2_1.UpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fallbackTo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"migrateUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedUpkeeps\",\"type\":\"bytes\"}],\"name\":\"receiveUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"enumKeeperRegistryBase2_1.Trigger\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"setUpkeepTriggerConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101406040523480156200001257600080fd5b50604051620061d1380380620061d18339810160408190526200003591620003df565b80816001600160a01b0316634b4fd03b6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200009b919062000406565b826001600160a01b031663ca30e6036040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000da573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001009190620003df565b836001600160a01b031663b10b673c6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200013f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001659190620003df565b846001600160a01b0316636709d0e56040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca9190620003df565b856001600160a01b0316635425d8ac6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000209573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200022f9190620003df565b3380600081620002865760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620002b957620002b9816200031b565b505050846002811115620002d157620002d162000429565b60e0816002811115620002e857620002e862000429565b9052506001600160a01b0393841660805291831660a052821660c0528116610100529190911661012052506200043f9050565b336001600160a01b03821603620003755760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200027d565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b0381168114620003dc57600080fd5b50565b600060208284031215620003f257600080fd5b8151620003ff81620003c6565b9392505050565b6000602082840312156200041957600080fd5b815160038110620003ff57600080fd5b634e487b7160e01b600052602160045260246000fd5b60805160a05160c05160e0516101005161012051615d18620004b96000396000818161010e01526101a90152600081816103e10152611fa10152600081816135370152818161376d015281816139b50152613b5d015260006130e1015260006131c5015260008181611de301526123af0152615d186000f3fe60806040523480156200001157600080fd5b50600436106200010c5760003560e01c806385c1b0ba11620000a5578063c8048022116200006f578063c804802214620002b7578063ce7dc5b414620002ce578063f2fde38b14620002e5578063f7d334ba14620002fc576200010c565b806385c1b0ba14620002535780638da5cb5b146200026a5780638e86139b1462000289578063948108f714620002a0576200010c565b80634ee88d3511620000e75780634ee88d3514620001ef5780636ded9eae146200020657806371791aa0146200021d57806379ba50971462000249576200010c565b806328f32f38146200015457806329c5efad146200017e578063349e8cca14620001a7575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e8080156200014d573d6000f35b3d6000fd5b005b6200016b62000165366004620041ea565b62000313565b6040519081526020015b60405180910390f35b620001956200018f366004620042d0565b6200068c565b604051620001759493929190620043f8565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200162000175565b620001526200020036600462004435565b62000930565b6200016b6200021736600462004485565b62000998565b620002346200022e366004620042d0565b620009fe565b60405162000175979695949392919062004538565b62000152620010f0565b62000152620002643660046200458a565b620011f3565b60005473ffffffffffffffffffffffffffffffffffffffff16620001c9565b620001526200029a36600462004617565b62001e64565b62000152620002b13660046200467a565b620021ec565b62000152620002c8366004620046a9565b6200247f565b62000195620002df3660046200477f565b62002846565b62000152620002f6366004620047f6565b62002916565b620002346200030d366004620046a9565b6200292e565b6000805473ffffffffffffffffffffffffffffffffffffffff163314801590620003475750620003456009336200296c565b155b156200037f576040517fd48b678b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff89163b620003ce576040517f09ee12d500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b620003d986620029a0565b9050600089307f00000000000000000000000000000000000000000000000000000000000000006040516200040e9062003f7b565b73ffffffffffffffffffffffffffffffffffffffff938416815291831660208301529091166040820152606001604051809103906000f08015801562000458573d6000803e3d6000fd5b5090506200051f826040518060e001604052806000151581526020018c63ffffffff16815260200163ffffffff801681526020018473ffffffffffffffffffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff168152602001600063ffffffff168152508a89898080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508b92508a915062002b449050565b6014805474010000000000000000000000000000000000000000900463ffffffff1690806200054e8362004845565b91906101000a81548163ffffffff021916908363ffffffff16021790555050817fbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d0128a8a604051620005c792919063ffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a2817fcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d878760405162000603929190620048b4565b60405180910390a2817f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664856040516200063d9190620048ca565b60405180910390a2817f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf485084604051620006779190620048ca565b60405180910390a25098975050505050505050565b600060606000806200069d62002f1f565b600086815260046020908152604091829020825160e081018452815460ff811615158252610100810463ffffffff90811694830194909452650100000000008104841694820194909452690100000000000000000090930473ffffffffffffffffffffffffffffffffffffffff166060840152600101546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a0840152780100000000000000000000000000000000000000000000000090041660c08201525a9150600080826060015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620007b7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620007dd9190620048ec565b73ffffffffffffffffffffffffffffffffffffffff166013600101600c9054906101000a900463ffffffff1663ffffffff16896040516200081f91906200490c565b60006040518083038160008787f1925050503d80600081146200085f576040519150601f19603f3d011682016040523d82523d6000602084013e62000864565b606091505b50915091505a6200087690856200492a565b935081620008a157600060405180602001604052806000815250600796509650965050505062000927565b80806020019051810190620008b791906200499b565b909750955086620008e557600060405180602001604052806000815250600496509650965050505062000927565b601554865164010000000090910463ffffffff1610156200092357600060405180602001604052806000815250600596509650965050505062000927565b5050505b92959194509250565b6200093b8362002f5a565b6000838152601a602052604090206200095682848362004a90565b50827f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d566483836040516200098b929190620048b4565b60405180910390a2505050565b6000620009f288888860008989604051806020016040528060008152508a8a8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506200031392505050565b98975050505050505050565b60006060600080600080600062000a1462002f1f565b600062000a218a62003010565b905060006012604051806101200160405290816000820160009054906101000a900460ff1660ff1660ff1681526020016000820160019054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160059054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160099054906101000a900462ffffff1662ffffff1662ffffff16815260200160008201600c9054906101000a900461ffff1661ffff1661ffff16815260200160008201600e9054906101000a900460ff1615151515815260200160008201600f9054906101000a900460ff161515151581526020016000820160109054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160008201601c9054906101000a900463ffffffff1663ffffffff1663ffffffff168152505090506000600460008d81526020019081526020016000206040518060e00160405290816000820160009054906101000a900460ff161515151581526020016000820160019054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160059054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160099054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160018201600c9054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff1681526020016001820160189054906101000a900463ffffffff1663ffffffff1663ffffffff168152505090508160a001511562000d45576000604051806020016040528060008152506009600084602001516000808263ffffffff1692509950995099509950995099509950505050620010e4565b604081015163ffffffff9081161462000d96576000604051806020016040528060008152506001600084602001516000808263ffffffff1692509950995099509950995099509950505050620010e4565b80511562000ddc576000604051806020016040528060008152506002600084602001516000808263ffffffff1692509950995099509950995099509950505050620010e4565b62000de782620030be565b602083015160155492975090955060009162000e19918591879190640100000000900463ffffffff168a8a87620032b0565b9050806bffffffffffffffffffffffff168260a001516bffffffffffffffffffffffff16101562000e83576000604051806020016040528060008152506006600085602001516000808263ffffffff1692509a509a509a509a509a509a509a5050505050620010e4565b600062000e928e868f62003301565b90505a9850600080846060015173ffffffffffffffffffffffffffffffffffffffff1663f00e6a2a6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000eea573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000f109190620048ec565b73ffffffffffffffffffffffffffffffffffffffff166013600101600c9054906101000a900463ffffffff1663ffffffff168460405162000f5291906200490c565b60006040518083038160008787f1925050503d806000811462000f92576040519150601f19603f3d011682016040523d82523d6000602084013e62000f97565b606091505b50915091505a62000fa9908c6200492a565b9a5081620010295760155481516801000000000000000090910463ffffffff1610156200100657505060408051602080820190925260008082529490910151939c509a50600899505063ffffffff9091169550620010e492505050565b602090940151939b5060039a505063ffffffff9092169650620010e49350505050565b808060200190518101906200103f91906200499b565b909e509c508d6200108057505060408051602080820190925260008082529490910151939c509a50600499505063ffffffff9091169550620010e492505050565b6015548d5164010000000090910463ffffffff161015620010d157505060408051602080820190925260008082529490910151939c509a50600599505063ffffffff9091169550620010e492505050565b505050506020015163ffffffff16945050505b92959891949750929550565b60015473ffffffffffffffffffffffffffffffffffffffff16331462001177576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600173ffffffffffffffffffffffffffffffffffffffff821660009081526019602052604090205460ff1660038111156200123257620012326200438d565b141580156200127e5750600373ffffffffffffffffffffffffffffffffffffffff821660009081526019602052604090205460ff1660038111156200127b576200127b6200438d565b14155b15620012b6576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6013546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1662001316576040517fd12d7d8d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082900362001352576040517f2c2fc94100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160e081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c081018290526000808567ffffffffffffffff811115620013a957620013a962004071565b604051908082528060200260200182016040528015620013d3578160200160208202803683370190505b50905060008667ffffffffffffffff811115620013f457620013f462004071565b6040519080825280602002602001820160405280156200147b57816020015b6040805160e08101825260008082526020808301829052928201819052606082018190526080820181905260a0820181905260c082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181620014135790505b50905060008767ffffffffffffffff8111156200149c576200149c62004071565b604051908082528060200260200182016040528015620014d157816020015b6060815260200190600190039081620014bb5790505b50905060008867ffffffffffffffff811115620014f257620014f262004071565b6040519080825280602002602001820160405280156200152757816020015b6060815260200190600190039081620015115790505b50905060008967ffffffffffffffff81111562001548576200154862004071565b6040519080825280602002602001820160405280156200157d57816020015b6060815260200190600190039081620015675790505b50905060005b8a81101562001b61578b8b82818110620015a157620015a162004bb8565b60209081029290920135600081815260048452604090819020815160e081018352815460ff811615158252610100810463ffffffff90811697830197909752650100000000008104871693820193909352690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff166060830152600101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490931660c08401529a509098506200168090508962002f5a565b60608801516040517f1a5da6c800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8c8116600483015290911690631a5da6c890602401600060405180830381600087803b158015620016f057600080fd5b505af115801562001705573d6000803e3d6000fd5b50505050878582815181106200171f576200171f62004bb8565b6020026020010181905250600560008a815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1686828151811062001773576200177362004bb8565b73ffffffffffffffffffffffffffffffffffffffff90921660209283029190910182015260008a81526007909152604090208054620017b290620049e8565b80601f0160208091040260200160405190810160405280929190818152602001828054620017e090620049e8565b8015620018315780601f10620018055761010080835404028352916020019162001831565b820191906000526020600020905b8154815290600101906020018083116200181357829003601f168201915b50505050508482815181106200184b576200184b62004bb8565b6020026020010181905250601a60008a815260200190815260200160002080546200187690620049e8565b80601f0160208091040260200160405190810160405280929190818152602001828054620018a490620049e8565b8015620018f55780601f10620018c957610100808354040283529160200191620018f5565b820191906000526020600020905b815481529060010190602001808311620018d757829003601f168201915b50505050508382815181106200190f576200190f62004bb8565b6020026020010181905250601b60008a815260200190815260200160002080546200193a90620049e8565b80601f01602080910402602001604051908101604052809291908181526020018280546200196890620049e8565b8015620019b95780601f106200198d57610100808354040283529160200191620019b9565b820191906000526020600020905b8154815290600101906020018083116200199b57829003601f168201915b5050505050828281518110620019d357620019d362004bb8565b60200260200101819052508760a001516bffffffffffffffffffffffff1687620019fe919062004be7565b60008a815260046020908152604080832080547fffffff000000000000000000000000000000000000000000000000000000000016815560010180547fffffffff000000000000000000000000000000000000000000000000000000001690556007909152812091985062001a74919062003f89565b6000898152601a6020526040812062001a8d9162003f89565b6000898152601b6020526040812062001aa69162003f89565b600089815260066020526040902080547fffffffffffffffffffffffff000000000000000000000000000000000000000016905562001ae760028a62003523565b5060a0880151604080516bffffffffffffffffffffffff909216825273ffffffffffffffffffffffffffffffffffffffff8c1660208301528a917fb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff910160405180910390a28062001b588162004bfd565b91505062001583565b508560185462001b7291906200492a565b60185560008b8b868167ffffffffffffffff81111562001b965762001b9662004071565b60405190808252806020026020018201604052801562001bc0578160200160208202803683370190505b508988888860405160200162001bde98979695949392919062004da3565b60405160208183030381529060405290508973ffffffffffffffffffffffffffffffffffffffff16638e86139b6013600001600c9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663c71249ab60038e73ffffffffffffffffffffffffffffffffffffffff1663aab9edd66040518163ffffffff1660e01b8152600401602060405180830381865afa15801562001c9a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001cc0919062004e73565b866040518463ffffffff1660e01b815260040162001ce19392919062004e98565b600060405180830381865afa15801562001cff573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405262001d47919081019062004ebf565b6040518263ffffffff1660e01b815260040162001d659190620048ca565b600060405180830381600087803b15801562001d8057600080fd5b505af115801562001d95573d6000803e3d6000fd5b50506040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8d81166004830152602482018b90527f000000000000000000000000000000000000000000000000000000000000000016925063a9059cbb91506044016020604051808303816000875af115801562001e2f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001e55919062004ef8565b50505050505050505050505050565b60023360009081526019602052604090205460ff16600381111562001e8d5762001e8d6200438d565b1415801562001ec3575060033360009081526019602052604090205460ff16600381111562001ec05762001ec06200438d565b14155b1562001efb576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080808080808062001f11888a018a620050f0565b965096509650965096509650965060005b8751811015620021e057600073ffffffffffffffffffffffffffffffffffffffff1687828151811062001f595762001f5962004bb8565b60200260200101516060015173ffffffffffffffffffffffffffffffffffffffff16036200206d5785818151811062001f965762001f9662004bb8565b6020026020010151307f000000000000000000000000000000000000000000000000000000000000000060405162001fce9062003f7b565b73ffffffffffffffffffffffffffffffffffffffff938416815291831660208301529091166040820152606001604051809103906000f08015801562002018573d6000803e3d6000fd5b508782815181106200202e576200202e62004bb8565b60200260200101516060019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250505b6200212588828151811062002086576200208662004bb8565b6020026020010151888381518110620020a357620020a362004bb8565b6020026020010151878481518110620020c057620020c062004bb8565b6020026020010151878581518110620020dd57620020dd62004bb8565b6020026020010151878681518110620020fa57620020fa62004bb8565b602002602001015187878151811062002117576200211762004bb8565b602002602001015162002b44565b8781815181106200213a576200213a62004bb8565b60200260200101517f74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a7188838151811062002178576200217862004bb8565b602002602001015160a0015133604051620021c39291906bffffffffffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a280620021d78162004bfd565b91505062001f22565b50505050505050505050565b600082815260046020908152604091829020825160e081018452815460ff81161515825263ffffffff6101008204811694830194909452650100000000008104841694820185905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004821660c08201529114620022ea576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818160a00151620022fc919062005221565b600084815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff93841602179055601854620023649184169062004be7565b6018556040517f23b872dd0000000000000000000000000000000000000000000000000000000081523360048201523060248201526bffffffffffffffffffffffff831660448201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906323b872dd906064016020604051808303816000875af11580156200240e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062002434919062004ef8565b506040516bffffffffffffffffffffffff83168152339084907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a3505050565b6000818152600460209081526040808320815160e081018352815460ff81161515825263ffffffff6101008204811695830195909552650100000000008104851693820184905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004831660c082015292911415906200256860005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16149050818015620025c35750808015620025c15750620025b462003531565b836040015163ffffffff16115b155b15620025fb576040517ffbc0357800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b801580156200262e575060008481526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314155b1562002666576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006200267262003531565b9050816200268a576200268760328262004be7565b90505b6000858152600460205260409020805463ffffffff80841665010000000000027fffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffffff90921691909117909155620026e69060029087906200352316565b5060135460808501516bffffffffffffffffffffffff91821691600091168211156200274f5760808601516200271d908362005249565b90508560a001516bffffffffffffffffffffffff16816bffffffffffffffffffffffff1611156200274f575060a08501515b808660a0015162002761919062005249565b600088815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff93841602179055601454620027c99183911662005221565b601480547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff9290921691909117905560405167ffffffffffffffff84169088907f91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f79118190600090a350505050505050565b600060606000806200285762002f1f565b6000634b56a42e60e01b888888604051602401620028789392919062005271565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915290506200290389826200068c565b929c919b50995090975095505050505050565b62002920620035ed565b6200292b8162003670565b50565b600060606000806000806000620029558860405180602001604052806000815250620009fe565b959e949d50929b5090995097509550909350915050565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415155b90505b92915050565b6000806000620029c76001620029b562003531565b620029c191906200492a565b62003767565b601454604080516020810193909352309083015274010000000000000000000000000000000000000000900463ffffffff166060820152608001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201209083015201604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052905060045b600f81101562002ad3578282828151811062002a8f5762002a8f62004bb8565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508062002aca8162004bfd565b91505062002a6f565b5083600181111562002ae95762002ae96200438d565b60f81b81600f8151811062002b025762002b0262004bb8565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535062002b3c81620052a5565b949350505050565b6012546e010000000000000000000000000000900460ff161562002b94576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601554835163ffffffff909116101562002bda576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108fc856020015163ffffffff16108062002c185750601454602086015163ffffffff70010000000000000000000000000000000090920482169116115b1562002c50576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000868152600460205260409020546901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff161562002cba576040517f6e3b930b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000868152600460209081526040808320885181548a8501518b85015160608d01517fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000009093169315157fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000ff169390931761010063ffffffff92831602177fffffff000000000000000000000000000000000000000000000000ffffffffff1665010000000000938216939093027fffffff0000000000000000000000000000000000000000ffffffffffffffffff1692909217690100000000000000000073ffffffffffffffffffffffffffffffffffffffff9283160217835560808b01516001909301805460a08d015160c08e01516bffffffffffffffffffffffff9687167fffffffffffffffff000000000000000000000000000000000000000000000000909316929092176c010000000000000000000000009690911695909502949094177fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff1678010000000000000000000000000000000000000000000000009490931693909302919091179091556005835281842080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169189169190911790556007909152902062002ead8482620052e8565b508460a001516bffffffffffffffffffffffff1660185462002ed0919062004be7565b6018556000868152601a6020526040902062002eed8382620052e8565b506000868152601b6020526040902062002f088282620052e8565b5062002f16600287620038cf565b50505050505050565b321562002f58576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60008181526005602052604090205473ffffffffffffffffffffffffffffffffffffffff16331462002fb8576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526004602052604090205465010000000000900463ffffffff908116146200292b576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000818160045b600f811015620030a5577fff00000000000000000000000000000000000000000000000000000000000000821683826020811062003059576200305962004bb8565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916146200309057506000949350505050565b806200309c8162004bfd565b91505062003017565b5081600f1a600181111562002b3c5762002b3c6200438d565b6000806000836060015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa1580156200314b573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200317191906200542a565b50945090925050506000811315806200318957508142105b80620031ae5750828015620031ae5750620031a582426200492a565b8463ffffffff16105b15620031bf576016549550620031c3565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa1580156200322f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200325591906200542a565b50945090925050506000811315806200326d57508142105b806200329257508280156200329257506200328982426200492a565b8463ffffffff16105b15620032a3576017549450620032a7565b8094505b50505050915091565b600080620032c488878b60000151620038dd565b9050600080620032e18b8a63ffffffff16858a8a60018b6200397c565b9092509050620032f2818362005221565b9b9a5050505050505050505050565b606060008360018111156200331a576200331a6200438d565b03620033e7576000848152600760205260409081902090517f6e04ff0d0000000000000000000000000000000000000000000000000000000091620033629160240162005522565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915290506200351c565b6001836001811115620033fe57620033fe6200438d565b03620034ea576000828060200190518101906200341c919062005599565b6000868152600760205260409081902090519192507f40691db4000000000000000000000000000000000000000000000000000000009162003463918491602401620056ad565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915291506200351c9050565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b9392505050565b600062002997838362003e1e565b600060017f000000000000000000000000000000000000000000000000000000000000000060028111156200356a576200356a6200438d565b03620035e857606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620035bd573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620035e3919062005775565b905090565b504390565b60005473ffffffffffffffffffffffffffffffffffffffff16331462002f58576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016200116e565b3373ffffffffffffffffffffffffffffffffffffffff821603620036f1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200116e565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060017f00000000000000000000000000000000000000000000000000000000000000006002811115620037a057620037a06200438d565b03620038c5576000606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620037f5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200381b919062005775565b905080831015806200383957506101006200383784836200492a565b115b15620038485750600092915050565b6040517f2b407a8200000000000000000000000000000000000000000000000000000000815260048101849052606490632b407a8290602401602060405180830381865afa1580156200389f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200351c919062005775565b504090565b919050565b600062002997838362003f29565b60008080856001811115620038f657620038f66200438d565b0362003907575062015f906200392a565b60018560018111156200391e576200391e6200438d565b03620034ea57506201adb05b6200393d63ffffffff851660146200578f565b6200394a846001620057cf565b6200395b9060ff16611d4c6200578f565b62003967908362004be7565b62003973919062004be7565b95945050505050565b6000806000896080015161ffff16876200399791906200578f565b9050838015620039a65750803a105b15620039af57503a5b600060027f00000000000000000000000000000000000000000000000000000000000000006002811115620039e857620039e86200438d565b0362003b5957604080516000815260208101909152851562003a4c5760003660405180608001604052806048815260200162005cc46048913960405160200162003a3593929190620057eb565b604051602081830303815290604052905062003aba565b60155462003a6a90640100000000900463ffffffff16600462005814565b63ffffffff1667ffffffffffffffff81111562003a8b5762003a8b62004071565b6040519080825280601f01601f19166020018201604052801562003ab6576020820181803683370190505b5090505b6040517f49948e0e00000000000000000000000000000000000000000000000000000000815273420000000000000000000000000000000000000f906349948e0e9062003b0c908490600401620048ca565b602060405180830381865afa15801562003b2a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003b50919062005775565b91505062003cc3565b60017f0000000000000000000000000000000000000000000000000000000000000000600281111562003b905762003b906200438d565b0362003cc357841562003c1857606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562003bea573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003c10919062005775565b905062003cc3565b6000606c73ffffffffffffffffffffffffffffffffffffffff166341b247a86040518163ffffffff1660e01b815260040160c060405180830381865afa15801562003c67573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003c8d919062005843565b505060155492945062003cb293505050640100000000900463ffffffff16826200578f565b62003cbf9060106200578f565b9150505b8462003ce257808b6080015161ffff1662003cdf91906200578f565b90505b62003cf261ffff8716826200588e565b90506000878262003d048c8e62004be7565b62003d1090866200578f565b62003d1c919062004be7565b62003d3090670de0b6b3a76400006200578f565b62003d3c91906200588e565b905060008c6040015163ffffffff1664e8d4a5100062003d5d91906200578f565b898e6020015163ffffffff16858f8862003d7891906200578f565b62003d84919062004be7565b62003d9490633b9aca006200578f565b62003da091906200578f565b62003dac91906200588e565b62003db8919062004be7565b90506b033b2e3c9fd0803ce800000062003dd3828462004be7565b111562003e0c576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b6000818152600183016020526040812054801562003f1757600062003e456001836200492a565b855490915060009062003e5b906001906200492a565b905081811462003ec757600086600001828154811062003e7f5762003e7f62004bb8565b906000526020600020015490508087600001848154811062003ea55762003ea562004bb8565b6000918252602080832090910192909255918252600188019052604090208390555b855486908062003edb5762003edb620058ca565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506200299a565b60009150506200299a565b5092915050565b600081815260018301602052604081205462003f72575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556200299a565b5060006200299a565b6103ca80620058fa83390190565b50805462003f9790620049e8565b6000825580601f1062003fa8575050565b601f0160209004906000526020600020908101906200292b91905b8082111562003fd9576000815560010162003fc3565b5090565b73ffffffffffffffffffffffffffffffffffffffff811681146200292b57600080fd5b803563ffffffff81168114620038ca57600080fd5b803560028110620038ca57600080fd5b60008083601f8401126200403857600080fd5b50813567ffffffffffffffff8111156200405157600080fd5b6020830191508360208285010111156200406a57600080fd5b9250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160e0810167ffffffffffffffff81118282101715620040c657620040c662004071565b60405290565b604051610100810167ffffffffffffffff81118282101715620040c657620040c662004071565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156200413d576200413d62004071565b604052919050565b600067ffffffffffffffff82111562004162576200416262004071565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f830112620041a057600080fd5b8135620041b7620041b18262004145565b620040f3565b818152846020838601011115620041cd57600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060008060e0898b0312156200420757600080fd5b8835620042148162003fdd565b97506200422460208a0162004000565b96506040890135620042368162003fdd565b95506200424660608a0162004015565b9450608089013567ffffffffffffffff808211156200426457600080fd5b620042728c838d0162004025565b909650945060a08b01359150808211156200428c57600080fd5b6200429a8c838d016200418e565b935060c08b0135915080821115620042b157600080fd5b50620042c08b828c016200418e565b9150509295985092959890939650565b60008060408385031215620042e457600080fd5b82359150602083013567ffffffffffffffff8111156200430357600080fd5b62004311858286016200418e565b9150509250929050565b60005b83811015620043385781810151838201526020016200431e565b50506000910152565b600081518084526200435b8160208601602086016200431b565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600a8110620043f4577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b9052565b841515815260806020820152600062004415608083018662004341565b9050620044266040830185620043bc565b82606083015295945050505050565b6000806000604084860312156200444b57600080fd5b83359250602084013567ffffffffffffffff8111156200446a57600080fd5b620044788682870162004025565b9497909650939450505050565b600080600080600080600060a0888a031215620044a157600080fd5b8735620044ae8162003fdd565b9650620044be6020890162004000565b95506040880135620044d08162003fdd565b9450606088013567ffffffffffffffff80821115620044ee57600080fd5b620044fc8b838c0162004025565b909650945060808a01359150808211156200451657600080fd5b50620045258a828b0162004025565b989b979a50959850939692959293505050565b871515815260e0602082015260006200455560e083018962004341565b9050620045666040830188620043bc565b8560608301528460808301528360a08301528260c083015298975050505050505050565b600080600060408486031215620045a057600080fd5b833567ffffffffffffffff80821115620045b957600080fd5b818601915086601f830112620045ce57600080fd5b813581811115620045de57600080fd5b8760208260051b8501011115620045f457600080fd5b602092830195509350508401356200460c8162003fdd565b809150509250925092565b600080602083850312156200462b57600080fd5b823567ffffffffffffffff8111156200464357600080fd5b620046518582860162004025565b90969095509350505050565b80356bffffffffffffffffffffffff81168114620038ca57600080fd5b600080604083850312156200468e57600080fd5b82359150620046a0602084016200465d565b90509250929050565b600060208284031215620046bc57600080fd5b5035919050565b600067ffffffffffffffff821115620046e057620046e062004071565b5060051b60200190565b600082601f830112620046fc57600080fd5b813560206200470f620041b183620046c3565b82815260059290921b840181019181810190868411156200472f57600080fd5b8286015b848110156200477457803567ffffffffffffffff811115620047555760008081fd5b620047658986838b01016200418e565b84525091830191830162004733565b509695505050505050565b600080600080606085870312156200479657600080fd5b84359350602085013567ffffffffffffffff80821115620047b657600080fd5b620047c488838901620046ea565b94506040870135915080821115620047db57600080fd5b50620047ea8782880162004025565b95989497509550505050565b6000602082840312156200480957600080fd5b81356200351c8162003fdd565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600063ffffffff80831681810362004861576200486162004816565b6001019392505050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b60208152600062002b3c6020830184866200486b565b60208152600062002997602083018462004341565b8051620038ca8162003fdd565b600060208284031215620048ff57600080fd5b81516200351c8162003fdd565b60008251620049208184602087016200431b565b9190910192915050565b818103818111156200299a576200299a62004816565b80151581146200292b57600080fd5b600082601f8301126200496157600080fd5b815162004972620041b18262004145565b8181528460208386010111156200498857600080fd5b62002b3c8260208301602087016200431b565b60008060408385031215620049af57600080fd5b8251620049bc8162004940565b602084015190925067ffffffffffffffff811115620049da57600080fd5b62004311858286016200494f565b600181811c90821680620049fd57607f821691505b60208210810362004a37577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f82111562004a8b57600081815260208120601f850160051c8101602086101562004a665750805b601f850160051c820191505b8181101562004a875782815560010162004a72565b5050505b505050565b67ffffffffffffffff83111562004aab5762004aab62004071565b62004ac38362004abc8354620049e8565b8362004a3d565b6000601f84116001811462004b18576000851562004ae15750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b17835562004bb1565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b8281101562004b69578685013582556020948501946001909201910162004b47565b508682101562004ba5577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b808201808211156200299a576200299a62004816565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820362004c315762004c3162004816565b5060010190565b600081518084526020808501945080840160005b8381101562004cf75781518051151588528381015163ffffffff908116858a01526040808301519091169089015260608082015173ffffffffffffffffffffffffffffffffffffffff16908901526080808201516bffffffffffffffffffffffff169089015260a08082015162004cd2828b01826bffffffffffffffffffffffff169052565b505060c09081015163ffffffff169088015260e0909601959082019060010162004c4c565b509495945050505050565b600081518084526020808501945080840160005b8381101562004cf757815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010162004d16565b600081518084526020808501808196508360051b8101915082860160005b8581101562004d9657828403895262004d8384835162004341565b9885019893509084019060010162004d68565b5091979650505050505050565b60e081528760e082015260006101007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8a111562004de057600080fd5b8960051b808c8386013783018381038201602085015262004e048282018b62004c38565b915050828103604084015262004e1b818962004d02565b9050828103606084015262004e31818862004d02565b9050828103608084015262004e47818762004d4a565b905082810360a084015262004e5d818662004d4a565b905082810360c0840152620032f2818562004d4a565b60006020828403121562004e8657600080fd5b815160ff811681146200351c57600080fd5b60ff8416815260ff8316602082015260606040820152600062003973606083018462004341565b60006020828403121562004ed257600080fd5b815167ffffffffffffffff81111562004eea57600080fd5b62002b3c848285016200494f565b60006020828403121562004f0b57600080fd5b81516200351c8162004940565b600082601f83011262004f2a57600080fd5b8135602062004f3d620041b183620046c3565b82815260059290921b8401810191818101908684111562004f5d57600080fd5b8286015b8481101562004774578035835291830191830162004f61565b600082601f83011262004f8c57600080fd5b8135602062004f9f620041b183620046c3565b82815260e0928302850182019282820191908785111562004fbf57600080fd5b8387015b85811015620050765781818a03121562004fdd5760008081fd5b62004fe7620040a0565b813562004ff48162004940565b81526200500382870162004000565b8682015260406200501681840162004000565b908201526060828101356200502b8162003fdd565b9082015260806200503e8382016200465d565b9082015260a0620050518382016200465d565b9082015260c06200506483820162004000565b90820152845292840192810162004fc3565b5090979650505050505050565b600082601f8301126200509557600080fd5b81356020620050a8620041b183620046c3565b82815260059290921b84018101918181019086841115620050c857600080fd5b8286015b8481101562004774578035620050e28162003fdd565b8352918301918301620050cc565b600080600080600080600060e0888a0312156200510c57600080fd5b873567ffffffffffffffff808211156200512557600080fd5b620051338b838c0162004f18565b985060208a01359150808211156200514a57600080fd5b620051588b838c0162004f7a565b975060408a01359150808211156200516f57600080fd5b6200517d8b838c0162005083565b965060608a01359150808211156200519457600080fd5b620051a28b838c0162005083565b955060808a0135915080821115620051b957600080fd5b620051c78b838c01620046ea565b945060a08a0135915080821115620051de57600080fd5b620051ec8b838c01620046ea565b935060c08a01359150808211156200520357600080fd5b50620052128a828b01620046ea565b91505092959891949750929550565b6bffffffffffffffffffffffff81811683821601908082111562003f225762003f2262004816565b6bffffffffffffffffffffffff82811682821603908082111562003f225762003f2262004816565b60408152600062005286604083018662004d4a565b82810360208401526200529b8185876200486b565b9695505050505050565b8051602080830151919081101562004a37577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60209190910360031b1b16919050565b815167ffffffffffffffff81111562005305576200530562004071565b6200531d81620053168454620049e8565b8462004a3d565b602080601f8311600181146200537357600084156200533c5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b17855562004a87565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015620053c257888601518255948401946001909101908401620053a1565b5085821015620053ff57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b805169ffffffffffffffffffff81168114620038ca57600080fd5b600080600080600060a086880312156200544357600080fd5b6200544e866200540f565b945060208601519350604086015192506060860151915062005473608087016200540f565b90509295509295909350565b600081546200548e81620049e8565b808552602060018381168015620054ae5760018114620054e75762005517565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b890101955062005517565b866000528260002060005b858110156200550f5781548a8201860152908301908401620054f2565b890184019650505b505050505092915050565b6020815260006200299760208301846200547f565b600082601f8301126200554957600080fd5b815160206200555c620041b183620046c3565b82815260059290921b840181019181810190868411156200557c57600080fd5b8286015b8481101562004774578051835291830191830162005580565b600060208284031215620055ac57600080fd5b815167ffffffffffffffff80821115620055c557600080fd5b908301906101008286031215620055db57600080fd5b620055e5620040cc565b82518152602083015160208201526040830151604082015260608301516060820152608083015160808201526200561f60a08401620048df565b60a082015260c0830151828111156200563757600080fd5b620056458782860162005537565b60c08301525060e0830151828111156200565e57600080fd5b6200566c878286016200494f565b60e08301525095945050505050565b600081518084526020808501945080840160005b8381101562004cf7578151875295820195908201906001016200568f565b60408152825160408201526020830151606082015260408301516080820152606083015160a0820152608083015160c082015273ffffffffffffffffffffffffffffffffffffffff60a08401511660e0820152600060c0840151610100808185015250620057206101408401826200567b565b905060e08501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0848303016101208501526200575e828262004341565b91505082810360208401526200397381856200547f565b6000602082840312156200578857600080fd5b5051919050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615620057ca57620057ca62004816565b500290565b60ff81811683821601908111156200299a576200299a62004816565b8284823760008382016000815283516200580a8183602088016200431b565b0195945050505050565b600063ffffffff808316818516818304811182151516156200583a576200583a62004816565b02949350505050565b60008060008060008060c087890312156200585d57600080fd5b865195506020870151945060408701519350606087015192506080870151915060a087015190509295509295509295565b600082620058c5577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfe60c060405234801561001057600080fd5b506040516103ca3803806103ca83398101604081905261002f91610076565b600080546001600160a01b0319166001600160a01b039384161790559181166080521660a0526100b9565b80516001600160a01b038116811461007157600080fd5b919050565b60008060006060848603121561008b57600080fd5b6100948461005a565b92506100a26020850161005a565b91506100b06040850161005a565b90509250925092565b60805160a0516102e76100e36000396000603801526000818160c4015261011701526102e76000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806379188d161461007b578063f00e6a2a146100aa575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e808015610076573d6000f35b3d6000fd5b61008e6100893660046101c1565b6100ee565b6040805192151583526020830191909152015b60405180910390f35b60405173ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001681526020016100a1565b60008054819073ffffffffffffffffffffffffffffffffffffffff16331461011557600080fd5b7f00000000000000000000000000000000000000000000000000000000000000005a91505a61138881101561014957600080fd5b61138881039050856040820482031161016157600080fd5b50803b61016d57600080fd5b6000808551602087016000858af192505a610188908361029a565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080604083850312156101d457600080fd5b82359150602083013567ffffffffffffffff808211156101f357600080fd5b818501915085601f83011261020757600080fd5b81358181111561021957610219610192565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561025f5761025f610192565b8160405282815288602084870101111561027857600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b818103818111156102d4577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b9291505056fea164736f6c6343000810000a307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000810000a",
}

var KeeperRegistryLogicAABI = KeeperRegistryLogicAMetaData.ABI

var KeeperRegistryLogicABin = KeeperRegistryLogicAMetaData.Bin

func DeployKeeperRegistryLogicA(auth *bind.TransactOpts, backend bind.ContractBackend, logicB common.Address) (common.Address, *types.Transaction, *KeeperRegistryLogicA, error) {
	parsed, err := KeeperRegistryLogicAMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistryLogicABin), backend, logicB)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperRegistryLogicA{address: address, abi: *parsed, KeeperRegistryLogicACaller: KeeperRegistryLogicACaller{contract: contract}, KeeperRegistryLogicATransactor: KeeperRegistryLogicATransactor{contract: contract}, KeeperRegistryLogicAFilterer: KeeperRegistryLogicAFilterer{contract: contract}}, nil
}

type KeeperRegistryLogicA struct {
	address common.Address
	abi     abi.ABI
	KeeperRegistryLogicACaller
	KeeperRegistryLogicATransactor
	KeeperRegistryLogicAFilterer
}

type KeeperRegistryLogicACaller struct {
	contract *bind.BoundContract
}

type KeeperRegistryLogicATransactor struct {
	contract *bind.BoundContract
}

type KeeperRegistryLogicAFilterer struct {
	contract *bind.BoundContract
}

type KeeperRegistryLogicASession struct {
	Contract     *KeeperRegistryLogicA
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeeperRegistryLogicACallerSession struct {
	Contract *KeeperRegistryLogicACaller
	CallOpts bind.CallOpts
}

type KeeperRegistryLogicATransactorSession struct {
	Contract     *KeeperRegistryLogicATransactor
	TransactOpts bind.TransactOpts
}

type KeeperRegistryLogicARaw struct {
	Contract *KeeperRegistryLogicA
}

type KeeperRegistryLogicACallerRaw struct {
	Contract *KeeperRegistryLogicACaller
}

type KeeperRegistryLogicATransactorRaw struct {
	Contract *KeeperRegistryLogicATransactor
}

func NewKeeperRegistryLogicA(address common.Address, backend bind.ContractBackend) (*KeeperRegistryLogicA, error) {
	abi, err := abi.JSON(strings.NewReader(KeeperRegistryLogicAABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeeperRegistryLogicA(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicA{address: address, abi: abi, KeeperRegistryLogicACaller: KeeperRegistryLogicACaller{contract: contract}, KeeperRegistryLogicATransactor: KeeperRegistryLogicATransactor{contract: contract}, KeeperRegistryLogicAFilterer: KeeperRegistryLogicAFilterer{contract: contract}}, nil
}

func NewKeeperRegistryLogicACaller(address common.Address, caller bind.ContractCaller) (*KeeperRegistryLogicACaller, error) {
	contract, err := bindKeeperRegistryLogicA(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicACaller{contract: contract}, nil
}

func NewKeeperRegistryLogicATransactor(address common.Address, transactor bind.ContractTransactor) (*KeeperRegistryLogicATransactor, error) {
	contract, err := bindKeeperRegistryLogicA(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicATransactor{contract: contract}, nil
}

func NewKeeperRegistryLogicAFilterer(address common.Address, filterer bind.ContractFilterer) (*KeeperRegistryLogicAFilterer, error) {
	contract, err := bindKeeperRegistryLogicA(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAFilterer{contract: contract}, nil
}

func bindKeeperRegistryLogicA(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeeperRegistryLogicAMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicARaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistryLogicA.Contract.KeeperRegistryLogicACaller.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicARaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.KeeperRegistryLogicATransactor.contract.Transfer(opts)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicARaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.KeeperRegistryLogicATransactor.contract.Transact(opts, method, params...)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicACallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistryLogicA.Contract.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.contract.Transfer(opts)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.contract.Transact(opts, method, params...)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicACaller) FallbackTo(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogicA.contract.Call(opts, &out, "fallbackTo")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) FallbackTo() (common.Address, error) {
	return _KeeperRegistryLogicA.Contract.FallbackTo(&_KeeperRegistryLogicA.CallOpts)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicACallerSession) FallbackTo() (common.Address, error) {
	return _KeeperRegistryLogicA.Contract.FallbackTo(&_KeeperRegistryLogicA.CallOpts)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicACaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogicA.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) Owner() (common.Address, error) {
	return _KeeperRegistryLogicA.Contract.Owner(&_KeeperRegistryLogicA.CallOpts)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicACallerSession) Owner() (common.Address, error) {
	return _KeeperRegistryLogicA.Contract.Owner(&_KeeperRegistryLogicA.CallOpts)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "acceptOwnership")
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.AcceptOwnership(&_KeeperRegistryLogicA.TransactOpts)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.AcceptOwnership(&_KeeperRegistryLogicA.TransactOpts)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "addFunds", id, amount)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.AddFunds(&_KeeperRegistryLogicA.TransactOpts, id, amount)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.AddFunds(&_KeeperRegistryLogicA.TransactOpts, id, amount)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "cancelUpkeep", id)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.CancelUpkeep(&_KeeperRegistryLogicA.TransactOpts, id)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.CancelUpkeep(&_KeeperRegistryLogicA.TransactOpts, id)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) CheckCallback(opts *bind.TransactOpts, id *big.Int, values [][]byte, extraData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "checkCallback", id, values, extraData)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) CheckCallback(id *big.Int, values [][]byte, extraData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.CheckCallback(&_KeeperRegistryLogicA.TransactOpts, id, values, extraData)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) CheckCallback(id *big.Int, values [][]byte, extraData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.CheckCallback(&_KeeperRegistryLogicA.TransactOpts, id, values, extraData)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) CheckUpkeep(opts *bind.TransactOpts, id *big.Int, triggerData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "checkUpkeep", id, triggerData)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) CheckUpkeep(id *big.Int, triggerData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.CheckUpkeep(&_KeeperRegistryLogicA.TransactOpts, id, triggerData)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) CheckUpkeep(id *big.Int, triggerData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.CheckUpkeep(&_KeeperRegistryLogicA.TransactOpts, id, triggerData)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) CheckUpkeep0(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "checkUpkeep0", id)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) CheckUpkeep0(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.CheckUpkeep0(&_KeeperRegistryLogicA.TransactOpts, id)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) CheckUpkeep0(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.CheckUpkeep0(&_KeeperRegistryLogicA.TransactOpts, id)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) ExecuteCallback(opts *bind.TransactOpts, id *big.Int, payload []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "executeCallback", id, payload)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) ExecuteCallback(id *big.Int, payload []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.ExecuteCallback(&_KeeperRegistryLogicA.TransactOpts, id, payload)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) ExecuteCallback(id *big.Int, payload []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.ExecuteCallback(&_KeeperRegistryLogicA.TransactOpts, id, payload)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "migrateUpkeeps", ids, destination)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.MigrateUpkeeps(&_KeeperRegistryLogicA.TransactOpts, ids, destination)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.MigrateUpkeeps(&_KeeperRegistryLogicA.TransactOpts, ids, destination)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "receiveUpkeeps", encodedUpkeeps)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.ReceiveUpkeeps(&_KeeperRegistryLogicA.TransactOpts, encodedUpkeeps)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.ReceiveUpkeeps(&_KeeperRegistryLogicA.TransactOpts, encodedUpkeeps)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "registerUpkeep", target, gasLimit, admin, triggerType, checkData, triggerConfig, offchainConfig)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.RegisterUpkeep(&_KeeperRegistryLogicA.TransactOpts, target, gasLimit, admin, triggerType, checkData, triggerConfig, offchainConfig)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.RegisterUpkeep(&_KeeperRegistryLogicA.TransactOpts, target, gasLimit, admin, triggerType, checkData, triggerConfig, offchainConfig)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) RegisterUpkeep0(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "registerUpkeep0", target, gasLimit, admin, checkData, offchainConfig)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) RegisterUpkeep0(target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.RegisterUpkeep0(&_KeeperRegistryLogicA.TransactOpts, target, gasLimit, admin, checkData, offchainConfig)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) RegisterUpkeep0(target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.RegisterUpkeep0(&_KeeperRegistryLogicA.TransactOpts, target, gasLimit, admin, checkData, offchainConfig)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) SetUpkeepTriggerConfig(opts *bind.TransactOpts, id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "setUpkeepTriggerConfig", id, triggerConfig)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.SetUpkeepTriggerConfig(&_KeeperRegistryLogicA.TransactOpts, id, triggerConfig)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.SetUpkeepTriggerConfig(&_KeeperRegistryLogicA.TransactOpts, id, triggerConfig)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "transferOwnership", to)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.TransferOwnership(&_KeeperRegistryLogicA.TransactOpts, to)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.TransferOwnership(&_KeeperRegistryLogicA.TransactOpts, to)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.RawTransact(opts, calldata)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.Fallback(&_KeeperRegistryLogicA.TransactOpts, calldata)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.Fallback(&_KeeperRegistryLogicA.TransactOpts, calldata)
}

type KeeperRegistryLogicAAdminPrivilegeConfigSetIterator struct {
	Event *KeeperRegistryLogicAAdminPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAAdminPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAAdminPrivilegeConfigSet)
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
		it.Event = new(KeeperRegistryLogicAAdminPrivilegeConfigSet)
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

func (it *KeeperRegistryLogicAAdminPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAAdminPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAAdminPrivilegeConfigSet struct {
	Admin           common.Address
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*KeeperRegistryLogicAAdminPrivilegeConfigSetIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAAdminPrivilegeConfigSetIterator{contract: _KeeperRegistryLogicA.contract, event: "AdminPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAAdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAAdminPrivilegeConfigSet)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseAdminPrivilegeConfigSet(log types.Log) (*KeeperRegistryLogicAAdminPrivilegeConfigSet, error) {
	event := new(KeeperRegistryLogicAAdminPrivilegeConfigSet)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicACancelledUpkeepReportIterator struct {
	Event *KeeperRegistryLogicACancelledUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicACancelledUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicACancelledUpkeepReport)
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
		it.Event = new(KeeperRegistryLogicACancelledUpkeepReport)
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

func (it *KeeperRegistryLogicACancelledUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicACancelledUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicACancelledUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicACancelledUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicACancelledUpkeepReportIterator{contract: _KeeperRegistryLogicA.contract, event: "CancelledUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicACancelledUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicACancelledUpkeepReport)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseCancelledUpkeepReport(log types.Log) (*KeeperRegistryLogicACancelledUpkeepReport, error) {
	event := new(KeeperRegistryLogicACancelledUpkeepReport)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicADedupKeyAddedIterator struct {
	Event *KeeperRegistryLogicADedupKeyAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicADedupKeyAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicADedupKeyAdded)
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
		it.Event = new(KeeperRegistryLogicADedupKeyAdded)
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

func (it *KeeperRegistryLogicADedupKeyAddedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicADedupKeyAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicADedupKeyAdded struct {
	DedupKey [32]byte
	Raw      types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*KeeperRegistryLogicADedupKeyAddedIterator, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicADedupKeyAddedIterator{contract: _KeeperRegistryLogicA.contract, event: "DedupKeyAdded", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicADedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicADedupKeyAdded)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseDedupKeyAdded(log types.Log) (*KeeperRegistryLogicADedupKeyAdded, error) {
	event := new(KeeperRegistryLogicADedupKeyAdded)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAFundsAddedIterator struct {
	Event *KeeperRegistryLogicAFundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAFundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAFundsAdded)
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
		it.Event = new(KeeperRegistryLogicAFundsAdded)
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

func (it *KeeperRegistryLogicAFundsAddedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAFundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAFundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistryLogicAFundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAFundsAddedIterator{contract: _KeeperRegistryLogicA.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAFundsAdded)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "FundsAdded", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseFundsAdded(log types.Log) (*KeeperRegistryLogicAFundsAdded, error) {
	event := new(KeeperRegistryLogicAFundsAdded)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAFundsWithdrawnIterator struct {
	Event *KeeperRegistryLogicAFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAFundsWithdrawn)
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
		it.Event = new(KeeperRegistryLogicAFundsWithdrawn)
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

func (it *KeeperRegistryLogicAFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAFundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAFundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAFundsWithdrawnIterator{contract: _KeeperRegistryLogicA.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAFundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAFundsWithdrawn)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseFundsWithdrawn(log types.Log) (*KeeperRegistryLogicAFundsWithdrawn, error) {
	event := new(KeeperRegistryLogicAFundsWithdrawn)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAInsufficientFundsUpkeepReportIterator struct {
	Event *KeeperRegistryLogicAInsufficientFundsUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAInsufficientFundsUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAInsufficientFundsUpkeepReport)
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
		it.Event = new(KeeperRegistryLogicAInsufficientFundsUpkeepReport)
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

func (it *KeeperRegistryLogicAInsufficientFundsUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAInsufficientFundsUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAInsufficientFundsUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAInsufficientFundsUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAInsufficientFundsUpkeepReportIterator{contract: _KeeperRegistryLogicA.contract, event: "InsufficientFundsUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAInsufficientFundsUpkeepReport)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseInsufficientFundsUpkeepReport(log types.Log) (*KeeperRegistryLogicAInsufficientFundsUpkeepReport, error) {
	event := new(KeeperRegistryLogicAInsufficientFundsUpkeepReport)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAOwnerFundsWithdrawnIterator struct {
	Event *KeeperRegistryLogicAOwnerFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAOwnerFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAOwnerFundsWithdrawn)
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
		it.Event = new(KeeperRegistryLogicAOwnerFundsWithdrawn)
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

func (it *KeeperRegistryLogicAOwnerFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAOwnerFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAOwnerFundsWithdrawn struct {
	Amount *big.Int
	Raw    types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*KeeperRegistryLogicAOwnerFundsWithdrawnIterator, error) {

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAOwnerFundsWithdrawnIterator{contract: _KeeperRegistryLogicA.contract, event: "OwnerFundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAOwnerFundsWithdrawn) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAOwnerFundsWithdrawn)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseOwnerFundsWithdrawn(log types.Log) (*KeeperRegistryLogicAOwnerFundsWithdrawn, error) {
	event := new(KeeperRegistryLogicAOwnerFundsWithdrawn)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAOwnershipTransferRequestedIterator struct {
	Event *KeeperRegistryLogicAOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAOwnershipTransferRequested)
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
		it.Event = new(KeeperRegistryLogicAOwnershipTransferRequested)
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

func (it *KeeperRegistryLogicAOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryLogicAOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAOwnershipTransferRequestedIterator{contract: _KeeperRegistryLogicA.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAOwnershipTransferRequested)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryLogicAOwnershipTransferRequested, error) {
	event := new(KeeperRegistryLogicAOwnershipTransferRequested)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAOwnershipTransferredIterator struct {
	Event *KeeperRegistryLogicAOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAOwnershipTransferred)
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
		it.Event = new(KeeperRegistryLogicAOwnershipTransferred)
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

func (it *KeeperRegistryLogicAOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryLogicAOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAOwnershipTransferredIterator{contract: _KeeperRegistryLogicA.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAOwnershipTransferred)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseOwnershipTransferred(log types.Log) (*KeeperRegistryLogicAOwnershipTransferred, error) {
	event := new(KeeperRegistryLogicAOwnershipTransferred)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAPausedIterator struct {
	Event *KeeperRegistryLogicAPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAPaused)
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
		it.Event = new(KeeperRegistryLogicAPaused)
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

func (it *KeeperRegistryLogicAPausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterPaused(opts *bind.FilterOpts) (*KeeperRegistryLogicAPausedIterator, error) {

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAPausedIterator{contract: _KeeperRegistryLogicA.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAPaused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAPaused)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "Paused", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParsePaused(log types.Log) (*KeeperRegistryLogicAPaused, error) {
	event := new(KeeperRegistryLogicAPaused)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAPayeesUpdatedIterator struct {
	Event *KeeperRegistryLogicAPayeesUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAPayeesUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAPayeesUpdated)
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
		it.Event = new(KeeperRegistryLogicAPayeesUpdated)
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

func (it *KeeperRegistryLogicAPayeesUpdatedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAPayeesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAPayeesUpdated struct {
	Transmitters []common.Address
	Payees       []common.Address
	Raw          types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterPayeesUpdated(opts *bind.FilterOpts) (*KeeperRegistryLogicAPayeesUpdatedIterator, error) {

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAPayeesUpdatedIterator{contract: _KeeperRegistryLogicA.contract, event: "PayeesUpdated", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAPayeesUpdated) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAPayeesUpdated)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParsePayeesUpdated(log types.Log) (*KeeperRegistryLogicAPayeesUpdated, error) {
	event := new(KeeperRegistryLogicAPayeesUpdated)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAPayeeshipTransferRequestedIterator struct {
	Event *KeeperRegistryLogicAPayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAPayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAPayeeshipTransferRequested)
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
		it.Event = new(KeeperRegistryLogicAPayeeshipTransferRequested)
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

func (it *KeeperRegistryLogicAPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAPayeeshipTransferRequested struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogicAPayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAPayeeshipTransferRequestedIterator{contract: _KeeperRegistryLogicA.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAPayeeshipTransferRequested)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistryLogicAPayeeshipTransferRequested, error) {
	event := new(KeeperRegistryLogicAPayeeshipTransferRequested)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAPayeeshipTransferredIterator struct {
	Event *KeeperRegistryLogicAPayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAPayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAPayeeshipTransferred)
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
		it.Event = new(KeeperRegistryLogicAPayeeshipTransferred)
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

func (it *KeeperRegistryLogicAPayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAPayeeshipTransferred struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogicAPayeeshipTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAPayeeshipTransferredIterator{contract: _KeeperRegistryLogicA.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAPayeeshipTransferred)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParsePayeeshipTransferred(log types.Log) (*KeeperRegistryLogicAPayeeshipTransferred, error) {
	event := new(KeeperRegistryLogicAPayeeshipTransferred)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAPaymentWithdrawnIterator struct {
	Event *KeeperRegistryLogicAPaymentWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAPaymentWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAPaymentWithdrawn)
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
		it.Event = new(KeeperRegistryLogicAPaymentWithdrawn)
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

func (it *KeeperRegistryLogicAPaymentWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAPaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAPaymentWithdrawn struct {
	Transmitter common.Address
	Amount      *big.Int
	To          common.Address
	Payee       common.Address
	Raw         types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistryLogicAPaymentWithdrawnIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAPaymentWithdrawnIterator{contract: _KeeperRegistryLogicA.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAPaymentWithdrawn)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParsePaymentWithdrawn(log types.Log) (*KeeperRegistryLogicAPaymentWithdrawn, error) {
	event := new(KeeperRegistryLogicAPaymentWithdrawn)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAReorgedUpkeepReportIterator struct {
	Event *KeeperRegistryLogicAReorgedUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAReorgedUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAReorgedUpkeepReport)
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
		it.Event = new(KeeperRegistryLogicAReorgedUpkeepReport)
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

func (it *KeeperRegistryLogicAReorgedUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAReorgedUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAReorgedUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAReorgedUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAReorgedUpkeepReportIterator{contract: _KeeperRegistryLogicA.contract, event: "ReorgedUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAReorgedUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAReorgedUpkeepReport)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseReorgedUpkeepReport(log types.Log) (*KeeperRegistryLogicAReorgedUpkeepReport, error) {
	event := new(KeeperRegistryLogicAReorgedUpkeepReport)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAStaleUpkeepReportIterator struct {
	Event *KeeperRegistryLogicAStaleUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAStaleUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAStaleUpkeepReport)
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
		it.Event = new(KeeperRegistryLogicAStaleUpkeepReport)
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

func (it *KeeperRegistryLogicAStaleUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAStaleUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAStaleUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAStaleUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAStaleUpkeepReportIterator{contract: _KeeperRegistryLogicA.contract, event: "StaleUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAStaleUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAStaleUpkeepReport)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseStaleUpkeepReport(log types.Log) (*KeeperRegistryLogicAStaleUpkeepReport, error) {
	event := new(KeeperRegistryLogicAStaleUpkeepReport)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAUnpausedIterator struct {
	Event *KeeperRegistryLogicAUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAUnpaused)
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
		it.Event = new(KeeperRegistryLogicAUnpaused)
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

func (it *KeeperRegistryLogicAUnpausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistryLogicAUnpausedIterator, error) {

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAUnpausedIterator{contract: _KeeperRegistryLogicA.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUnpaused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAUnpaused)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseUnpaused(log types.Log) (*KeeperRegistryLogicAUnpaused, error) {
	event := new(KeeperRegistryLogicAUnpaused)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAUpkeepAdminTransferRequestedIterator struct {
	Event *KeeperRegistryLogicAUpkeepAdminTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAUpkeepAdminTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAUpkeepAdminTransferRequested)
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
		it.Event = new(KeeperRegistryLogicAUpkeepAdminTransferRequested)
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

func (it *KeeperRegistryLogicAUpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAUpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAUpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogicAUpkeepAdminTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAUpkeepAdminTransferRequestedIterator{contract: _KeeperRegistryLogicA.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAUpkeepAdminTransferRequested)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseUpkeepAdminTransferRequested(log types.Log) (*KeeperRegistryLogicAUpkeepAdminTransferRequested, error) {
	event := new(KeeperRegistryLogicAUpkeepAdminTransferRequested)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAUpkeepAdminTransferredIterator struct {
	Event *KeeperRegistryLogicAUpkeepAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAUpkeepAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAUpkeepAdminTransferred)
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
		it.Event = new(KeeperRegistryLogicAUpkeepAdminTransferred)
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

func (it *KeeperRegistryLogicAUpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAUpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAUpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogicAUpkeepAdminTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAUpkeepAdminTransferredIterator{contract: _KeeperRegistryLogicA.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAUpkeepAdminTransferred)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseUpkeepAdminTransferred(log types.Log) (*KeeperRegistryLogicAUpkeepAdminTransferred, error) {
	event := new(KeeperRegistryLogicAUpkeepAdminTransferred)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAUpkeepCanceledIterator struct {
	Event *KeeperRegistryLogicAUpkeepCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAUpkeepCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAUpkeepCanceled)
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
		it.Event = new(KeeperRegistryLogicAUpkeepCanceled)
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

func (it *KeeperRegistryLogicAUpkeepCanceledIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAUpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAUpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistryLogicAUpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAUpkeepCanceledIterator{contract: _KeeperRegistryLogicA.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAUpkeepCanceled)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseUpkeepCanceled(log types.Log) (*KeeperRegistryLogicAUpkeepCanceled, error) {
	event := new(KeeperRegistryLogicAUpkeepCanceled)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAUpkeepCheckDataSetIterator struct {
	Event *KeeperRegistryLogicAUpkeepCheckDataSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAUpkeepCheckDataSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAUpkeepCheckDataSet)
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
		it.Event = new(KeeperRegistryLogicAUpkeepCheckDataSet)
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

func (it *KeeperRegistryLogicAUpkeepCheckDataSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAUpkeepCheckDataSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAUpkeepCheckDataSet struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepCheckDataSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAUpkeepCheckDataSetIterator{contract: _KeeperRegistryLogicA.contract, event: "UpkeepCheckDataSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepCheckDataSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAUpkeepCheckDataSet)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseUpkeepCheckDataSet(log types.Log) (*KeeperRegistryLogicAUpkeepCheckDataSet, error) {
	event := new(KeeperRegistryLogicAUpkeepCheckDataSet)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAUpkeepGasLimitSetIterator struct {
	Event *KeeperRegistryLogicAUpkeepGasLimitSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAUpkeepGasLimitSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAUpkeepGasLimitSet)
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
		it.Event = new(KeeperRegistryLogicAUpkeepGasLimitSet)
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

func (it *KeeperRegistryLogicAUpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAUpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAUpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAUpkeepGasLimitSetIterator{contract: _KeeperRegistryLogicA.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAUpkeepGasLimitSet)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseUpkeepGasLimitSet(log types.Log) (*KeeperRegistryLogicAUpkeepGasLimitSet, error) {
	event := new(KeeperRegistryLogicAUpkeepGasLimitSet)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAUpkeepMigratedIterator struct {
	Event *KeeperRegistryLogicAUpkeepMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAUpkeepMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAUpkeepMigrated)
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
		it.Event = new(KeeperRegistryLogicAUpkeepMigrated)
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

func (it *KeeperRegistryLogicAUpkeepMigratedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAUpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAUpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAUpkeepMigratedIterator{contract: _KeeperRegistryLogicA.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAUpkeepMigrated)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseUpkeepMigrated(log types.Log) (*KeeperRegistryLogicAUpkeepMigrated, error) {
	event := new(KeeperRegistryLogicAUpkeepMigrated)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAUpkeepOffchainConfigSetIterator struct {
	Event *KeeperRegistryLogicAUpkeepOffchainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAUpkeepOffchainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAUpkeepOffchainConfigSet)
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
		it.Event = new(KeeperRegistryLogicAUpkeepOffchainConfigSet)
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

func (it *KeeperRegistryLogicAUpkeepOffchainConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAUpkeepOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAUpkeepOffchainConfigSet struct {
	Id             *big.Int
	OffchainConfig []byte
	Raw            types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAUpkeepOffchainConfigSetIterator{contract: _KeeperRegistryLogicA.contract, event: "UpkeepOffchainConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAUpkeepOffchainConfigSet)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseUpkeepOffchainConfigSet(log types.Log) (*KeeperRegistryLogicAUpkeepOffchainConfigSet, error) {
	event := new(KeeperRegistryLogicAUpkeepOffchainConfigSet)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAUpkeepPausedIterator struct {
	Event *KeeperRegistryLogicAUpkeepPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAUpkeepPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAUpkeepPaused)
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
		it.Event = new(KeeperRegistryLogicAUpkeepPaused)
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

func (it *KeeperRegistryLogicAUpkeepPausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAUpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAUpkeepPaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAUpkeepPausedIterator{contract: _KeeperRegistryLogicA.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAUpkeepPaused)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseUpkeepPaused(log types.Log) (*KeeperRegistryLogicAUpkeepPaused, error) {
	event := new(KeeperRegistryLogicAUpkeepPaused)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAUpkeepPerformedIterator struct {
	Event *KeeperRegistryLogicAUpkeepPerformed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAUpkeepPerformedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAUpkeepPerformed)
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
		it.Event = new(KeeperRegistryLogicAUpkeepPerformed)
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

func (it *KeeperRegistryLogicAUpkeepPerformedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAUpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAUpkeepPerformed struct {
	Id           *big.Int
	Success      bool
	TotalPayment *big.Int
	GasUsed      *big.Int
	GasOverhead  *big.Int
	Trigger      []byte
	Raw          types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*KeeperRegistryLogicAUpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAUpkeepPerformedIterator{contract: _KeeperRegistryLogicA.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAUpkeepPerformed)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseUpkeepPerformed(log types.Log) (*KeeperRegistryLogicAUpkeepPerformed, error) {
	event := new(KeeperRegistryLogicAUpkeepPerformed)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAUpkeepPrivilegeConfigSetIterator struct {
	Event *KeeperRegistryLogicAUpkeepPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAUpkeepPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAUpkeepPrivilegeConfigSet)
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
		it.Event = new(KeeperRegistryLogicAUpkeepPrivilegeConfigSet)
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

func (it *KeeperRegistryLogicAUpkeepPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAUpkeepPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAUpkeepPrivilegeConfigSet struct {
	Id              *big.Int
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepPrivilegeConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAUpkeepPrivilegeConfigSetIterator{contract: _KeeperRegistryLogicA.contract, event: "UpkeepPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAUpkeepPrivilegeConfigSet)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseUpkeepPrivilegeConfigSet(log types.Log) (*KeeperRegistryLogicAUpkeepPrivilegeConfigSet, error) {
	event := new(KeeperRegistryLogicAUpkeepPrivilegeConfigSet)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAUpkeepReceivedIterator struct {
	Event *KeeperRegistryLogicAUpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAUpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAUpkeepReceived)
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
		it.Event = new(KeeperRegistryLogicAUpkeepReceived)
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

func (it *KeeperRegistryLogicAUpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAUpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAUpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAUpkeepReceivedIterator{contract: _KeeperRegistryLogicA.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAUpkeepReceived)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseUpkeepReceived(log types.Log) (*KeeperRegistryLogicAUpkeepReceived, error) {
	event := new(KeeperRegistryLogicAUpkeepReceived)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAUpkeepRegisteredIterator struct {
	Event *KeeperRegistryLogicAUpkeepRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAUpkeepRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAUpkeepRegistered)
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
		it.Event = new(KeeperRegistryLogicAUpkeepRegistered)
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

func (it *KeeperRegistryLogicAUpkeepRegisteredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAUpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAUpkeepRegistered struct {
	Id         *big.Int
	PerformGas uint32
	Admin      common.Address
	Raw        types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAUpkeepRegisteredIterator{contract: _KeeperRegistryLogicA.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAUpkeepRegistered)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseUpkeepRegistered(log types.Log) (*KeeperRegistryLogicAUpkeepRegistered, error) {
	event := new(KeeperRegistryLogicAUpkeepRegistered)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAUpkeepTriggerConfigSetIterator struct {
	Event *KeeperRegistryLogicAUpkeepTriggerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAUpkeepTriggerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAUpkeepTriggerConfigSet)
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
		it.Event = new(KeeperRegistryLogicAUpkeepTriggerConfigSet)
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

func (it *KeeperRegistryLogicAUpkeepTriggerConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAUpkeepTriggerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAUpkeepTriggerConfigSet struct {
	Id            *big.Int
	TriggerConfig []byte
	Raw           types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepTriggerConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAUpkeepTriggerConfigSetIterator{contract: _KeeperRegistryLogicA.contract, event: "UpkeepTriggerConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAUpkeepTriggerConfigSet)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseUpkeepTriggerConfigSet(log types.Log) (*KeeperRegistryLogicAUpkeepTriggerConfigSet, error) {
	event := new(KeeperRegistryLogicAUpkeepTriggerConfigSet)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicAUpkeepUnpausedIterator struct {
	Event *KeeperRegistryLogicAUpkeepUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAUpkeepUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAUpkeepUnpaused)
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
		it.Event = new(KeeperRegistryLogicAUpkeepUnpaused)
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

func (it *KeeperRegistryLogicAUpkeepUnpausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAUpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAUpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAUpkeepUnpausedIterator{contract: _KeeperRegistryLogicA.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAUpkeepUnpaused)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseUpkeepUnpaused(log types.Log) (*KeeperRegistryLogicAUpkeepUnpaused, error) {
	event := new(KeeperRegistryLogicAUpkeepUnpaused)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicA) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _KeeperRegistryLogicA.abi.Events["AdminPrivilegeConfigSet"].ID:
		return _KeeperRegistryLogicA.ParseAdminPrivilegeConfigSet(log)
	case _KeeperRegistryLogicA.abi.Events["CancelledUpkeepReport"].ID:
		return _KeeperRegistryLogicA.ParseCancelledUpkeepReport(log)
	case _KeeperRegistryLogicA.abi.Events["DedupKeyAdded"].ID:
		return _KeeperRegistryLogicA.ParseDedupKeyAdded(log)
	case _KeeperRegistryLogicA.abi.Events["FundsAdded"].ID:
		return _KeeperRegistryLogicA.ParseFundsAdded(log)
	case _KeeperRegistryLogicA.abi.Events["FundsWithdrawn"].ID:
		return _KeeperRegistryLogicA.ParseFundsWithdrawn(log)
	case _KeeperRegistryLogicA.abi.Events["InsufficientFundsUpkeepReport"].ID:
		return _KeeperRegistryLogicA.ParseInsufficientFundsUpkeepReport(log)
	case _KeeperRegistryLogicA.abi.Events["OwnerFundsWithdrawn"].ID:
		return _KeeperRegistryLogicA.ParseOwnerFundsWithdrawn(log)
	case _KeeperRegistryLogicA.abi.Events["OwnershipTransferRequested"].ID:
		return _KeeperRegistryLogicA.ParseOwnershipTransferRequested(log)
	case _KeeperRegistryLogicA.abi.Events["OwnershipTransferred"].ID:
		return _KeeperRegistryLogicA.ParseOwnershipTransferred(log)
	case _KeeperRegistryLogicA.abi.Events["Paused"].ID:
		return _KeeperRegistryLogicA.ParsePaused(log)
	case _KeeperRegistryLogicA.abi.Events["PayeesUpdated"].ID:
		return _KeeperRegistryLogicA.ParsePayeesUpdated(log)
	case _KeeperRegistryLogicA.abi.Events["PayeeshipTransferRequested"].ID:
		return _KeeperRegistryLogicA.ParsePayeeshipTransferRequested(log)
	case _KeeperRegistryLogicA.abi.Events["PayeeshipTransferred"].ID:
		return _KeeperRegistryLogicA.ParsePayeeshipTransferred(log)
	case _KeeperRegistryLogicA.abi.Events["PaymentWithdrawn"].ID:
		return _KeeperRegistryLogicA.ParsePaymentWithdrawn(log)
	case _KeeperRegistryLogicA.abi.Events["ReorgedUpkeepReport"].ID:
		return _KeeperRegistryLogicA.ParseReorgedUpkeepReport(log)
	case _KeeperRegistryLogicA.abi.Events["StaleUpkeepReport"].ID:
		return _KeeperRegistryLogicA.ParseStaleUpkeepReport(log)
	case _KeeperRegistryLogicA.abi.Events["Unpaused"].ID:
		return _KeeperRegistryLogicA.ParseUnpaused(log)
	case _KeeperRegistryLogicA.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepAdminTransferRequested(log)
	case _KeeperRegistryLogicA.abi.Events["UpkeepAdminTransferred"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepAdminTransferred(log)
	case _KeeperRegistryLogicA.abi.Events["UpkeepCanceled"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepCanceled(log)
	case _KeeperRegistryLogicA.abi.Events["UpkeepCheckDataSet"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepCheckDataSet(log)
	case _KeeperRegistryLogicA.abi.Events["UpkeepGasLimitSet"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepGasLimitSet(log)
	case _KeeperRegistryLogicA.abi.Events["UpkeepMigrated"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepMigrated(log)
	case _KeeperRegistryLogicA.abi.Events["UpkeepOffchainConfigSet"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepOffchainConfigSet(log)
	case _KeeperRegistryLogicA.abi.Events["UpkeepPaused"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepPaused(log)
	case _KeeperRegistryLogicA.abi.Events["UpkeepPerformed"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepPerformed(log)
	case _KeeperRegistryLogicA.abi.Events["UpkeepPrivilegeConfigSet"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepPrivilegeConfigSet(log)
	case _KeeperRegistryLogicA.abi.Events["UpkeepReceived"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepReceived(log)
	case _KeeperRegistryLogicA.abi.Events["UpkeepRegistered"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepRegistered(log)
	case _KeeperRegistryLogicA.abi.Events["UpkeepTriggerConfigSet"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepTriggerConfigSet(log)
	case _KeeperRegistryLogicA.abi.Events["UpkeepUnpaused"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeeperRegistryLogicAAdminPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x7c44b4eb59ee7873514e7e43e7718c269d872965938b288aa143befca62f99d2")
}

func (KeeperRegistryLogicACancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd5636")
}

func (KeeperRegistryLogicADedupKeyAdded) Topic() common.Hash {
	return common.HexToHash("0xa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f2")
}

func (KeeperRegistryLogicAFundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (KeeperRegistryLogicAFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (KeeperRegistryLogicAInsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e02")
}

func (KeeperRegistryLogicAOwnerFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1")
}

func (KeeperRegistryLogicAOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (KeeperRegistryLogicAOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (KeeperRegistryLogicAPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (KeeperRegistryLogicAPayeesUpdated) Topic() common.Hash {
	return common.HexToHash("0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725")
}

func (KeeperRegistryLogicAPayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (KeeperRegistryLogicAPayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (KeeperRegistryLogicAPaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (KeeperRegistryLogicAReorgedUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301")
}

func (KeeperRegistryLogicAStaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8")
}

func (KeeperRegistryLogicAUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (KeeperRegistryLogicAUpkeepAdminTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35")
}

func (KeeperRegistryLogicAUpkeepAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c")
}

func (KeeperRegistryLogicAUpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (KeeperRegistryLogicAUpkeepCheckDataSet) Topic() common.Hash {
	return common.HexToHash("0xcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d")
}

func (KeeperRegistryLogicAUpkeepGasLimitSet) Topic() common.Hash {
	return common.HexToHash("0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c")
}

func (KeeperRegistryLogicAUpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (KeeperRegistryLogicAUpkeepOffchainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850")
}

func (KeeperRegistryLogicAUpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (KeeperRegistryLogicAUpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b")
}

func (KeeperRegistryLogicAUpkeepPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae7769")
}

func (KeeperRegistryLogicAUpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (KeeperRegistryLogicAUpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (KeeperRegistryLogicAUpkeepTriggerConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664")
}

func (KeeperRegistryLogicAUpkeepUnpaused) Topic() common.Hash {
	return common.HexToHash("0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456")
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicA) Address() common.Address {
	return _KeeperRegistryLogicA.address
}

type KeeperRegistryLogicAInterface interface {
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

	FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*KeeperRegistryLogicAAdminPrivilegeConfigSetIterator, error)

	WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAAdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error)

	ParseAdminPrivilegeConfigSet(log types.Log) (*KeeperRegistryLogicAAdminPrivilegeConfigSet, error)

	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicACancelledUpkeepReportIterator, error)

	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicACancelledUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseCancelledUpkeepReport(log types.Log) (*KeeperRegistryLogicACancelledUpkeepReport, error)

	FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*KeeperRegistryLogicADedupKeyAddedIterator, error)

	WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicADedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error)

	ParseDedupKeyAdded(log types.Log) (*KeeperRegistryLogicADedupKeyAdded, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistryLogicAFundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*KeeperRegistryLogicAFundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAFundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAFundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*KeeperRegistryLogicAFundsWithdrawn, error)

	FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAInsufficientFundsUpkeepReportIterator, error)

	WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseInsufficientFundsUpkeepReport(log types.Log) (*KeeperRegistryLogicAInsufficientFundsUpkeepReport, error)

	FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*KeeperRegistryLogicAOwnerFundsWithdrawnIterator, error)

	WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAOwnerFundsWithdrawn) (event.Subscription, error)

	ParseOwnerFundsWithdrawn(log types.Log) (*KeeperRegistryLogicAOwnerFundsWithdrawn, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryLogicAOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryLogicAOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryLogicAOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*KeeperRegistryLogicAOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*KeeperRegistryLogicAPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*KeeperRegistryLogicAPaused, error)

	FilterPayeesUpdated(opts *bind.FilterOpts) (*KeeperRegistryLogicAPayeesUpdatedIterator, error)

	WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAPayeesUpdated) (event.Subscription, error)

	ParsePayeesUpdated(log types.Log) (*KeeperRegistryLogicAPayeesUpdated, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogicAPayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistryLogicAPayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogicAPayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*KeeperRegistryLogicAPayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistryLogicAPaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*KeeperRegistryLogicAPaymentWithdrawn, error)

	FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAReorgedUpkeepReportIterator, error)

	WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAReorgedUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseReorgedUpkeepReport(log types.Log) (*KeeperRegistryLogicAReorgedUpkeepReport, error)

	FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAStaleUpkeepReportIterator, error)

	WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAStaleUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseStaleUpkeepReport(log types.Log) (*KeeperRegistryLogicAStaleUpkeepReport, error)

	FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistryLogicAUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*KeeperRegistryLogicAUnpaused, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogicAUpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*KeeperRegistryLogicAUpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogicAUpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*KeeperRegistryLogicAUpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistryLogicAUpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*KeeperRegistryLogicAUpkeepCanceled, error)

	FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepCheckDataSetIterator, error)

	WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepCheckDataSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCheckDataSet(log types.Log) (*KeeperRegistryLogicAUpkeepCheckDataSet, error)

	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepGasLimitSetIterator, error)

	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepGasLimitSet(log types.Log) (*KeeperRegistryLogicAUpkeepGasLimitSet, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*KeeperRegistryLogicAUpkeepMigrated, error)

	FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepOffchainConfigSetIterator, error)

	WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepOffchainConfigSet(log types.Log) (*KeeperRegistryLogicAUpkeepOffchainConfigSet, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*KeeperRegistryLogicAUpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*KeeperRegistryLogicAUpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*KeeperRegistryLogicAUpkeepPerformed, error)

	FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepPrivilegeConfigSetIterator, error)

	WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPrivilegeConfigSet(log types.Log) (*KeeperRegistryLogicAUpkeepPrivilegeConfigSet, error)

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*KeeperRegistryLogicAUpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*KeeperRegistryLogicAUpkeepRegistered, error)

	FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepTriggerConfigSetIterator, error)

	WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepTriggerConfigSet(log types.Log) (*KeeperRegistryLogicAUpkeepTriggerConfigSet, error)

	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepUnpausedIterator, error)

	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepUnpaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepUnpaused(log types.Log) (*KeeperRegistryLogicAUpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
