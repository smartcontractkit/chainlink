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
	ABI: "[{\"inputs\":[{\"internalType\":\"contractKeeperRegistryLogicB2_1\",\"name\":\"logicB\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTrigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTriggerType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByUpkeepManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PipelineDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"adminOffchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepAdminOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newPipelineData\",\"type\":\"bytes\"}],\"name\":\"UpkeepPipelineDataSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepTriggerConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumKeeperRegistryBase2_1.UpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumKeeperRegistryBase2_1.UpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumKeeperRegistryBase2_1.UpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"executeCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumKeeperRegistryBase2_1.UpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fallbackTo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFastGasFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkNativeFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMode\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase2_1.Mode\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTriggerType\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase2_1.Trigger\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"migrateUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedUpkeeps\",\"type\":\"bytes\"}],\"name\":\"receiveUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"enumKeeperRegistryBase2_1.Trigger\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"setUpkeepTriggerConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTranscoderVersion\",\"outputs\":[{\"internalType\":\"enumUpkeepFormat\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepVersion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
<<<<<<< HEAD
	Bin: "0x6101206040523480156200001257600080fd5b5060405162005e1938038062005e19833981016040819052620000359162000386565b80816001600160a01b0316634b4fd03b6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200009b9190620003ad565b826001600160a01b031663ca30e6036040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000da573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000100919062000386565b836001600160a01b031663b10b673c6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200013f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000165919062000386565b846001600160a01b0316636709d0e56040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca919062000386565b3380600081620002215760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b038481169190911790915581161562000254576200025481620002c2565b5050508360028111156200026c576200026c620003d0565b60e0816002811115620002835762000283620003d0565b9052506001600160a01b0392831660805290821660a052811660c052601980546001600160a01b03191633179055919091166101005250620003e69050565b336001600160a01b038216036200031c5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000218565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200038357600080fd5b50565b6000602082840312156200039957600080fd5b8151620003a6816200036d565b9392505050565b600060208284031215620003c057600080fd5b815160038110620003a657600080fd5b634e487b7160e01b600052602160045260246000fd5b60805160a05160c05160e051610100516159af6200046a600039600081816101aa01526102450152600081816102a50152818161343201528181613668015281816138ee0152613a9601526000818161030901526131fe0152600081816103fd01526132e201526000818161043b01528181611feb015261248101526159af6000f3fe60806040523480156200001157600080fd5b5060043610620001a85760003560e01c806385c1b0ba11620000ed578063b10b673c1162000099578063ce7dc5b4116200006f578063ce7dc5b41462000460578063f2fde38b1462000477578063f7d334ba146200048e57620001a8565b8063b10b673c14620003fb578063c80480221462000422578063ca30e603146200043957620001a8565b80638e86139b11620000cf5780638e86139b14620003b1578063948108f714620003c8578063aab9edd614620003df57620001a8565b806385c1b0ba146200037b5780638da5cb5b146200039257620001a8565b80634ee88d3511620001595780636ded9eae116200012f5780636ded9eae146200032e57806371791aa0146200034557806379ba5097146200037157620001a8565b80634ee88d3514620002ca5780635147cd5914620002e15780636709d0e5146200030757620001a8565b8063349e8cca116200018f578063349e8cca146200024357806348013d7b146200028b5780634b4fd03b14620002a357620001a8565b806328f32f3814620001f057806329c5efad146200021a575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e808015620001e9573d6000f35b3d6000fd5b005b620002076200020136600462004053565b620004a5565b6040519081526020015b60405180910390f35b620002316200022b36600462004139565b620007c8565b6040516200021194939291906200423c565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200162000211565b62000294600081565b6040516200021191906200428c565b7f000000000000000000000000000000000000000000000000000000000000000062000294565b620001ee620002db366004620042a1565b62000995565b620002f8620002f2366004620042f1565b620009fd565b6040516200021191906200430b565b7f000000000000000000000000000000000000000000000000000000000000000062000265565b620002076200033f36600462004322565b62000ab3565b6200035c6200035636600462004139565b62000b5f565b604051620002119796959493929190620043d5565b620001ee620012fa565b620001ee6200038c36600462004427565b620013fd565b60005473ffffffffffffffffffffffffffffffffffffffff1662000265565b620001ee620003c2366004620044b4565b6200206c565b620001ee620003d936600462004517565b620022af565b620003e8600281565b60405160ff909116815260200162000211565b7f000000000000000000000000000000000000000000000000000000000000000062000265565b620001ee62000433366004620042f1565b62002551565b7f000000000000000000000000000000000000000000000000000000000000000062000265565b620002316200047136600462004602565b62002927565b620001ee6200048836600462004679565b620029f7565b6200035c6200049f366004620042f1565b62002a0f565b6000805473ffffffffffffffffffffffffffffffffffffffff163314801590620004d95750620004d760093362002ae1565b155b1562000511576040517fd48b678b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6200051c8662002b15565b90506000818a30604051620005319062003dfd565b92835273ffffffffffffffffffffffffffffffffffffffff9182166020840152166040820152606001604051809103906000f08015801562000577573d6000803e3d6000fd5b5090506200065b826040518061010001604052806000151581526020018c63ffffffff16815260200163ffffffff801681526020018473ffffffffffffffffffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff168152602001600063ffffffff1681526020018d73ffffffffffffffffffffffffffffffffffffffff168152508a89898080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508b92508a915062002cb19050565b6014805474010000000000000000000000000000000000000000900463ffffffff1690806200068a83620046c8565b91906101000a81548163ffffffff021916908363ffffffff16021790555050817fbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d0128a8a6040516200070392919063ffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a2817f787b2e556c364148324d7dc0cee4322a048aaf798bda184136c8636c72ed305687876040516200073f92919062004737565b60405180910390a2817f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664856040516200077991906200474d565b60405180910390a2817f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf485084604051620007b391906200474d565b60405180910390a25098975050505050505050565b60006060600080620007d9620030ea565b60008681526004602090815260409182902082516101008082018552825460ff81161515835263ffffffff91810482169483019490945265010000000000840481169482019490945273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009093048316606082015260018201546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490931660c0840152600201541660e08201525a60e0820151601454604051929450600092839273ffffffffffffffffffffffffffffffffffffffff16916c01000000000000000000000000900463ffffffff1690620008f9908b9062004762565b60006040518083038160008787f1925050503d806000811462000939576040519150601f19603f3d011682016040523d82523d6000602084013e6200093e565b606091505b50915091505a62000950908562004780565b9350816200096257600794506200097e565b80806020019051810190620009789190620047f1565b90975095505b866200098957600494505b50505092959194509250565b620009a08362003125565b6000838152601b60205260409020620009bb828483620048e6565b50827f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d56648383604051620009f092919062004737565b60405180910390a2505050565b6000818160045b600f81101562000a92577fff00000000000000000000000000000000000000000000000000000000000000821683826020811062000a465762000a4662004a0e565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161462000a7d57506000949350505050565b8062000a898162004a3d565b91505062000a04565b5081600f1a600181111562000aab5762000aab620041f6565b949350505050565b600062000b53888888600089896040518060200160405280600163ffffffff1681525060405160200162000af0915163ffffffff16815260200190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181526020601f8d018190048102840181019092528b835291908c908c9081908401838280828437600092019190915250620004a592505050565b98975050505050505050565b60006060600080600080600062000b75620030ea565b600062000b828a620009fd565b905060006012604051806101200160405290816000820160009054906101000a900460ff1660ff1660ff1681526020016000820160019054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160059054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160099054906101000a900462ffffff1662ffffff1662ffffff16815260200160008201600c9054906101000a900461ffff1661ffff1661ffff16815260200160008201600e9054906101000a900460ff1615151515815260200160008201600f9054906101000a900460ff161515151581526020016000820160109054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160008201601c9054906101000a900463ffffffff1663ffffffff1663ffffffff168152505090506000600460008d8152602001908152602001600020604051806101000160405290816000820160009054906101000a900460ff161515151581526020016000820160019054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160059054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160099054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160018201600c9054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff1681526020016001820160189054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016002820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681525050905063ffffffff8016816040015163ffffffff161462000f0a576000604051806020016040528060008152506001600084602001516000808263ffffffff1692509950995099509950995099509950505050620012ee565b80511562000f50576000604051806020016040528060008152506002600084602001516000808263ffffffff1692509950995099509950995099509950505050620012ee565b62000f5b82620031db565b602083015160155492975090955060009162000f8d918591879190640100000000900463ffffffff168a8a87620033cd565b9050806bffffffffffffffffffffffff168260a001516bffffffffffffffffffffffff16101562000ff7576000604051806020016040528060008152506006600085602001516000808263ffffffff1692509a509a509a509a509a509a509a5050505050620012ee565b60019a5060606000856001811115620010145762001014620041f6565b03620010d4576040517f6e04ff0d00000000000000000000000000000000000000000000000000000000906200104f908f906024016200474d565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152905062001189565b6040517fbe61b775000000000000000000000000000000000000000000000000000000009062001109908f906024016200474d565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915290505b5a98508260e0015173ffffffffffffffffffffffffffffffffffffffff166013600101600c9054906101000a900463ffffffff1663ffffffff1682604051620011d3919062004762565b60006040518083038160008787f1925050503d806000811462001213576040519150601f19603f3d011682016040523d82523d6000602084013e62001218565b606091505b50909c509a505a6200122b908a62004780565b98508b6200123d576003995062001290565b8a806020019051810190620012539190620047f1565b909c509a508b6200129057505060408051602080820190925260008082529290910151919a5098506004975063ffffffff169450620012ee915050565b6015548b5164010000000090910463ffffffff161015620012dd57505060408051602080820190925260008082529290910151919a5098506005975063ffffffff169450620012ee915050565b50506020015163ffffffff16945050505b92959891949750929550565b60015473ffffffffffffffffffffffffffffffffffffffff16331462001381576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600173ffffffffffffffffffffffffffffffffffffffff82166000908152601a602052604090205460ff1660038111156200143c576200143c620041f6565b14158015620014885750600373ffffffffffffffffffffffffffffffffffffffff82166000908152601a602052604090205460ff166003811115620014855762001485620041f6565b14155b15620014c0576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6013546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1662001520576040517fd12d7d8d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008290036200155c576040517f2c2fc94100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805161010081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526000808567ffffffffffffffff811115620015bb57620015bb62003f00565b604051908082528060200260200182016040528015620015e5578160200160208202803683370190505b50905060008667ffffffffffffffff81111562001606576200160662003f00565b6040519080825280602002602001820160405280156200169557816020015b604080516101008101825260008082526020808301829052928201819052606082018190526080820181905260a0820181905260c0820181905260e082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181620016255790505b50905060008767ffffffffffffffff811115620016b657620016b662003f00565b604051908082528060200260200182016040528015620016eb57816020015b6060815260200190600190039081620016d55790505b50905060008867ffffffffffffffff8111156200170c576200170c62003f00565b6040519080825280602002602001820160405280156200174157816020015b60608152602001906001900390816200172b5790505b50905060008967ffffffffffffffff81111562001762576200176262003f00565b6040519080825280602002602001820160405280156200179757816020015b6060815260200190600190039081620017815790505b50905060005b8a81101562001db2578b8b82818110620017bb57620017bb62004a0e565b6020908102929092013560008181526004845260409081902081516101008082018452825460ff81161515835263ffffffff91810482169783019790975265010000000000870481169382019390935273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009096048616606082015260018201546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490921660c08301526002015490931660e08401529a50909850620018a790508962003125565b60608801516040517f1a5da6c800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8c8116600483015290911690631a5da6c890602401600060405180830381600087803b1580156200191757600080fd5b505af11580156200192c573d6000803e3d6000fd5b505050508785828151811062001946576200194662004a0e565b6020026020010181905250600560008a815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff168682815181106200199a576200199a62004a0e565b73ffffffffffffffffffffffffffffffffffffffff90921660209283029190910182015260008a81526007909152604090208054620019d9906200483e565b80601f016020809104026020016040519081016040528092919081815260200182805462001a07906200483e565b801562001a585780601f1062001a2c5761010080835404028352916020019162001a58565b820191906000526020600020905b81548152906001019060200180831162001a3a57829003601f168201915b505050505084828151811062001a725762001a7262004a0e565b6020026020010181905250601b60008a8152602001908152602001600020805462001a9d906200483e565b80601f016020809104026020016040519081016040528092919081815260200182805462001acb906200483e565b801562001b1c5780601f1062001af05761010080835404028352916020019162001b1c565b820191906000526020600020905b81548152906001019060200180831162001afe57829003601f168201915b505050505083828151811062001b365762001b3662004a0e565b6020026020010181905250601c60008a8152602001908152602001600020805462001b61906200483e565b80601f016020809104026020016040519081016040528092919081815260200182805462001b8f906200483e565b801562001be05780601f1062001bb45761010080835404028352916020019162001be0565b820191906000526020600020905b81548152906001019060200180831162001bc257829003601f168201915b505050505082828151811062001bfa5762001bfa62004a0e565b60200260200101819052508760a001516bffffffffffffffffffffffff168762001c25919062004a78565b60008a815260046020908152604080832080547fffffff00000000000000000000000000000000000000000000000000000000001681556001810180547fffffffff0000000000000000000000000000000000000000000000000000000016905560020180547fffffffffffffffffffffffff00000000000000000000000000000000000000001690556007909152812091985062001cc5919062003e0b565b6000898152601b6020526040812062001cde9162003e0b565b6000898152601c6020526040812062001cf79162003e0b565b600089815260066020526040902080547fffffffffffffffffffffffff000000000000000000000000000000000000000016905562001d3860028a6200341e565b5060a0880151604080516bffffffffffffffffffffffff909216825273ffffffffffffffffffffffffffffffffffffffff8c1660208301528a917fb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff910160405180910390a28062001da98162004a3d565b9150506200179d565b508560185462001dc3919062004780565b60185560405160009062001de8908d908d9088908a9089908990899060200162004b3a565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905260135490915073ffffffffffffffffffffffffffffffffffffffff808c1691638e86139b916c010000000000000000000000009091041663c71249ab60028e73ffffffffffffffffffffffffffffffffffffffff1663aab9edd66040518163ffffffff1660e01b81526004016020604051808303816000875af115801562001ea2573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001ec8919062004cb4565b866040518463ffffffff1660e01b815260040162001ee99392919062004cd9565b600060405180830381865afa15801562001f07573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405262001f4f919081019062004d00565b6040518263ffffffff1660e01b815260040162001f6d91906200474d565b600060405180830381600087803b15801562001f8857600080fd5b505af115801562001f9d573d6000803e3d6000fd5b50506040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8d81166004830152602482018b90527f000000000000000000000000000000000000000000000000000000000000000016925063a9059cbb91506044016020604051808303816000875af115801562002037573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200205d919062004d39565b50505050505050505050505050565b6002336000908152601a602052604090205460ff166003811115620020955762002095620041f6565b14158015620020cb57506003336000908152601a602052604090205460ff166003811115620020c857620020c8620041f6565b14155b1562002103576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008080808080620021188789018962004f39565b95509550955095509550955060005b8651811015620022a457620021e98782815181106200214a576200214a62004a0e565b602002602001015187838151811062002167576200216762004a0e565b602002602001015187848151811062002184576200218462004a0e565b6020026020010151878581518110620021a157620021a162004a0e565b6020026020010151878681518110620021be57620021be62004a0e565b6020026020010151878781518110620021db57620021db62004a0e565b602002602001015162002cb1565b868181518110620021fe57620021fe62004a0e565b60200260200101517f74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a718783815181106200223c576200223c62004a0e565b602002602001015160a0015133604051620022879291906bffffffffffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a2806200229b8162004a3d565b91505062002127565b505050505050505050565b60008281526004602090815260409182902082516101008082018552825460ff81161515835263ffffffff918104821694830194909452650100000000008404811694820185905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009094048416606083015260018301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a084015278010000000000000000000000000000000000000000000000009004811660c083015260029092015490921660e0830152909114620023bc576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818160a00151620023ce919062005041565b600084815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff93841602179055601854620024369184169062004a78565b6018556040517f23b872dd0000000000000000000000000000000000000000000000000000000081523360048201523060248201526bffffffffffffffffffffffff831660448201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906323b872dd906064016020604051808303816000875af1158015620024e0573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062002506919062004d39565b506040516bffffffffffffffffffffffff83168152339084907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a3505050565b600081815260046020908152604080832081516101008082018452825460ff81161515835263ffffffff918104821695830195909552650100000000008504811693820184905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009095048516606083015260018301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a084015278010000000000000000000000000000000000000000000000009004811660c083015260029092015490931660e08401529192911415906200264960005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16149050818015620026a45750808015620026a25750620026956200342c565b836040015163ffffffff16115b155b15620026dc576040517ffbc0357800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b801580156200270f575060008481526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314155b1562002747576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000620027536200342c565b9050816200276b576200276860328262004a78565b90505b6000858152600460205260409020805463ffffffff80841665010000000000027fffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffffff90921691909117909155620027c79060029087906200341e16565b5060135460808501516bffffffffffffffffffffffff918216916000911682111562002830576080860151620027fe908362005069565b90508560a001516bffffffffffffffffffffffff16816bffffffffffffffffffffffff16111562002830575060a08501515b808660a0015162002842919062005069565b600088815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff93841602179055601454620028aa9183911662005041565b601480547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff9290921691909117905560405167ffffffffffffffff84169088907f91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f79118190600090a350505050505050565b6000606060008062002938620030ea565b6000634b56a42e60e01b888888604051602401620029599392919062005091565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091529050620029e48982620007c8565b929c919b50995090975095505050505050565b62002a01620034e8565b62002a0c816200356b565b50565b60006060600080600080600062002aca88600760008b8152602001908152602001600020805462002a40906200483e565b80601f016020809104026020016040519081016040528092919081815260200182805462002a6e906200483e565b801562002abf5780601f1062002a935761010080835404028352916020019162002abf565b820191906000526020600020905b81548152906001019060200180831162002aa157829003601f168201915b505050505062000b5f565b959e949d50929b5090995097509550909350915050565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415155b90505b92915050565b600080600062002b3c600162002b2a6200342c565b62002b36919062004780565b62003662565b601454604080516020810193909352309083015274010000000000000000000000000000000000000000900463ffffffff166060820152608001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201209083015201604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052905060045b600f81101562002c48578282828151811062002c045762002c0462004a0e565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508062002c3f8162004a3d565b91505062002be4565b5083600181111562002c5e5762002c5e620041f6565b60f81b81600f8151811062002c775762002c7762004a0e565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535062000aab81620050c5565b6012546e010000000000000000000000000000900460ff161562002d01576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60e085015173ffffffffffffffffffffffffffffffffffffffff163b62002d54576040517f09ee12d500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601554835163ffffffff909116101562002d9a576040517f2b49197900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108fc856020015163ffffffff16108062002dd85750601454602086015163ffffffff70010000000000000000000000000000000090920482169116115b1562002e10576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008681526004602052604090206002015473ffffffffffffffffffffffffffffffffffffffff161562002e70576040517f6e3b930b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600086815260046020908152604080832088518154848b0151848c015160608d015173ffffffffffffffffffffffffffffffffffffffff9081166901000000000000000000027fffffff0000000000000000000000000000000000000000ffffffffffffffffff63ffffffff9384166501000000000002167fffffff000000000000000000000000000000000000000000000000ffffffffff948416610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000ff971515979097167fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000009096169590951795909517929092169290921792909217835560808b015160018401805460a08e015160c08f01519094167801000000000000000000000000000000000000000000000000027fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff6bffffffffffffffffffffffff9586166c01000000000000000000000000027fffffffffffffffff0000000000000000000000000000000000000000000000009093169590941694909417179190911691909117905560e08a0151600290920180549282167fffffffffffffffffffffffff0000000000000000000000000000000000000000938416179055600584528285208054918a16919092161790556007909152902062003078848262005108565b508460a001516bffffffffffffffffffffffff166018546200309b919062004a78565b6018556000868152601b60205260409020620030b8838262005108565b506000868152601c60205260409020620030d3828262005108565b50620030e1600287620037d1565b50505050505050565b321562003123576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60008181526005602052604090205473ffffffffffffffffffffffffffffffffffffffff16331462003183576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526004602052604090205465010000000000900463ffffffff9081161462002a0c576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000806000836060015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa15801562003268573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200328e91906200524a565b5094509092505050600081131580620032a657508142105b80620032cb5750828015620032cb5750620032c2824262004780565b8463ffffffff16105b15620032dc576016549550620032e0565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa1580156200334c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200337291906200524a565b50945090925050506000811315806200338a57508142105b80620033af5750828015620033af5750620033a6824262004780565b8463ffffffff16105b15620033c0576017549450620033c4565b8094505b50505050915091565b600080620033e188878b60000151620037df565b9050600080620033fe8b8a63ffffffff16858a8a60018b620038b5565b90925090506200340f818362005041565b9b9a5050505050505050505050565b600062002b0c838362003ca0565b600060017f00000000000000000000000000000000000000000000000000000000000000006002811115620034655762003465620041f6565b03620034e357606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620034b8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620034de91906200529f565b905090565b504390565b60005473ffffffffffffffffffffffffffffffffffffffff16331462003123576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162001378565b3373ffffffffffffffffffffffffffffffffffffffff821603620035ec576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162001378565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060017f000000000000000000000000000000000000000000000000000000000000000060028111156200369b576200369b620041f6565b03620037c7576000606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620036f0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200371691906200529f565b9050808310158062003734575061010062003732848362004780565b115b15620037435750600092915050565b6040517f2b407a8200000000000000000000000000000000000000000000000000000000815260048101849052606490632b407a8290602401602060405180830381865afa1580156200379a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620037c091906200529f565b9392505050565b504090565b919050565b600062002b0c838362003dab565b60008080856001811115620037f857620037f8620041f6565b036200380957506201388062003863565b6001856001811115620038205762003820620041f6565b03620038315750620186a062003863565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6200387663ffffffff85166014620052b9565b62003883846001620052f9565b620038949060ff16611d4c620052b9565b620038a0908362004a78565b620038ac919062004a78565b95945050505050565b6000806000896080015161ffff1687620038d09190620052b9565b9050838015620038df5750803a105b15620038e857503a5b600060027f00000000000000000000000000000000000000000000000000000000000000006002811115620039215762003921620041f6565b0362003a9257604080516000815260208101909152851562003985576000366040518060800160405280604881526020016200595b604891396040516020016200396e9392919062005315565b6040516020818303038152906040529050620039f3565b601554620039a390640100000000900463ffffffff1660046200533e565b63ffffffff1667ffffffffffffffff811115620039c457620039c462003f00565b6040519080825280601f01601f191660200182016040528015620039ef576020820181803683370190505b5090505b6040517f49948e0e00000000000000000000000000000000000000000000000000000000815273420000000000000000000000000000000000000f906349948e0e9062003a459084906004016200474d565b602060405180830381865afa15801562003a63573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003a8991906200529f565b91505062003b45565b60017f0000000000000000000000000000000000000000000000000000000000000000600281111562003ac95762003ac9620041f6565b0362003b4557606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562003b1c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003b4291906200529f565b90505b8462003b6457808b6080015161ffff1662003b619190620052b9565b90505b62003b7461ffff8716826200536d565b90506000878262003b868c8e62004a78565b62003b929086620052b9565b62003b9e919062004a78565b62003bb290670de0b6b3a7640000620052b9565b62003bbe91906200536d565b905060008c6040015163ffffffff1664e8d4a5100062003bdf9190620052b9565b898e6020015163ffffffff16858f8862003bfa9190620052b9565b62003c06919062004a78565b62003c1690633b9aca00620052b9565b62003c229190620052b9565b62003c2e91906200536d565b62003c3a919062004a78565b90506b033b2e3c9fd0803ce800000062003c55828462004a78565b111562003c8e576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b6000818152600183016020526040812054801562003d9957600062003cc760018362004780565b855490915060009062003cdd9060019062004780565b905081811462003d4957600086600001828154811062003d015762003d0162004a0e565b906000526020600020015490508087600001848154811062003d275762003d2762004a0e565b6000918252602080832090910192909255918252600188019052604090208390555b855486908062003d5d5762003d5d620053a9565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505062002b0f565b600091505062002b0f565b5092915050565b600081815260018301602052604081205462003df45750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915562002b0f565b50600062002b0f565b61058280620053d983390190565b50805462003e19906200483e565b6000825580601f1062003e2a575050565b601f01602090049060005260206000209081019062002a0c91905b8082111562003e5b576000815560010162003e45565b5090565b73ffffffffffffffffffffffffffffffffffffffff8116811462002a0c57600080fd5b8035620037cc8162003e5f565b803563ffffffff81168114620037cc57600080fd5b803560028110620037cc57600080fd5b60008083601f84011262003ec757600080fd5b50813567ffffffffffffffff81111562003ee057600080fd5b60208301915083602082850101111562003ef957600080fd5b9250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610100810167ffffffffffffffff8111828210171562003f565762003f5662003f00565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171562003fa65762003fa662003f00565b604052919050565b600067ffffffffffffffff82111562003fcb5762003fcb62003f00565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f8301126200400957600080fd5b8135620040206200401a8262003fae565b62003f5c565b8181528460208386010111156200403657600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060008060e0898b0312156200407057600080fd5b88356200407d8162003e5f565b97506200408d60208a0162003e8f565b965060408901356200409f8162003e5f565b9550620040af60608a0162003ea4565b9450608089013567ffffffffffffffff80821115620040cd57600080fd5b620040db8c838d0162003eb4565b909650945060a08b0135915080821115620040f557600080fd5b620041038c838d0162003ff7565b935060c08b01359150808211156200411a57600080fd5b50620041298b828c0162003ff7565b9150509295985092959890939650565b600080604083850312156200414d57600080fd5b82359150602083013567ffffffffffffffff8111156200416c57600080fd5b6200417a8582860162003ff7565b9150509250929050565b60005b83811015620041a157818101518382015260200162004187565b50506000910152565b60008151808452620041c481602086016020860162004184565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b60088110620042385762004238620041f6565b9052565b8415158152608060208201526000620042596080830186620041aa565b90506200426a604083018562004225565b82606083015295945050505050565b6003811062002a0c5762002a0c620041f6565b602081016200429b8362004279565b91905290565b600080600060408486031215620042b757600080fd5b83359250602084013567ffffffffffffffff811115620042d657600080fd5b620042e48682870162003eb4565b9497909650939450505050565b6000602082840312156200430457600080fd5b5035919050565b60208101600283106200429b576200429b620041f6565b600080600080600080600060a0888a0312156200433e57600080fd5b87356200434b8162003e5f565b96506200435b6020890162003e8f565b955060408801356200436d8162003e5f565b9450606088013567ffffffffffffffff808211156200438b57600080fd5b620043998b838c0162003eb4565b909650945060808a0135915080821115620043b357600080fd5b50620043c28a828b0162003eb4565b989b979a50959850939692959293505050565b871515815260e060208201526000620043f260e0830189620041aa565b905062004403604083018862004225565b8560608301528460808301528360a08301528260c083015298975050505050505050565b6000806000604084860312156200443d57600080fd5b833567ffffffffffffffff808211156200445657600080fd5b818601915086601f8301126200446b57600080fd5b8135818111156200447b57600080fd5b8760208260051b85010111156200449157600080fd5b60209283019550935050840135620044a98162003e5f565b809150509250925092565b60008060208385031215620044c857600080fd5b823567ffffffffffffffff811115620044e057600080fd5b620044ee8582860162003eb4565b90969095509350505050565b80356bffffffffffffffffffffffff81168114620037cc57600080fd5b600080604083850312156200452b57600080fd5b823591506200453d60208401620044fa565b90509250929050565b600067ffffffffffffffff82111562004563576200456362003f00565b5060051b60200190565b600082601f8301126200457f57600080fd5b81356020620045926200401a8362004546565b82815260059290921b84018101918181019086841115620045b257600080fd5b8286015b84811015620045f757803567ffffffffffffffff811115620045d85760008081fd5b620045e88986838b010162003ff7565b845250918301918301620045b6565b509695505050505050565b600080600080606085870312156200461957600080fd5b84359350602085013567ffffffffffffffff808211156200463957600080fd5b62004647888389016200456d565b945060408701359150808211156200465e57600080fd5b506200466d8782880162003eb4565b95989497509550505050565b6000602082840312156200468c57600080fd5b8135620037c08162003e5f565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600063ffffffff808316818103620046e457620046e462004699565b6001019392505050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b60208152600062000aab602083018486620046ee565b60208152600062002b0c6020830184620041aa565b600082516200477681846020870162004184565b9190910192915050565b8181038181111562002b0f5762002b0f62004699565b801515811462002a0c57600080fd5b600082601f830112620047b757600080fd5b8151620047c86200401a8262003fae565b818152846020838601011115620047de57600080fd5b62000aab82602083016020870162004184565b600080604083850312156200480557600080fd5b8251620048128162004796565b602084015190925067ffffffffffffffff8111156200483057600080fd5b6200417a85828601620047a5565b600181811c908216806200485357607f821691505b6020821081036200488d577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f821115620048e157600081815260208120601f850160051c81016020861015620048bc5750805b601f850160051c820191505b81811015620048dd57828155600101620048c8565b5050505b505050565b67ffffffffffffffff83111562004901576200490162003f00565b62004919836200491283546200483e565b8362004893565b6000601f8411600181146200496e5760008515620049375750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b17835562004a07565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b82811015620049bf57868501358255602094850194600190920191016200499d565b5086821015620049fb577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820362004a715762004a7162004699565b5060010190565b8082018082111562002b0f5762002b0f62004699565b600081518084526020808501945080840160005b8381101562004ad657815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010162004aa2565b509495945050505050565b600081518084526020808501808196508360051b8101915082860160005b8581101562004b2d57828403895262004b1a848351620041aa565b9885019893509084019060010162004aff565b5091979650505050505050565b600060c0808352888184015260e07f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8a111562004b7657600080fd5b8960051b808c83870137840184810382016020808701919091528a518383018190528b82019261010092919083019060005b8181101562004c495785518051151584528481015163ffffffff9081168686015260408083015182169086015260608083015173ffffffffffffffffffffffffffffffffffffffff908116918701919091526080808401516bffffffffffffffffffffffff9081169188019190915260a080850151909116908701528a8301519091168a860152908801511687840152948301949184019160010162004ba8565b5050878103604089015262004c5f818d62004a8e565b95505050505050828103606084015262004c7a818762004ae1565b9050828103608084015262004c90818662004ae1565b905082810360a084015262004ca6818562004ae1565b9a9950505050505050505050565b60006020828403121562004cc757600080fd5b815160ff81168114620037c057600080fd5b60ff8416815260ff83166020820152606060408201526000620038ac6060830184620041aa565b60006020828403121562004d1357600080fd5b815167ffffffffffffffff81111562004d2b57600080fd5b62000aab84828501620047a5565b60006020828403121562004d4c57600080fd5b8151620037c08162004796565b600082601f83011262004d6b57600080fd5b8135602062004d7e6200401a8362004546565b82815260059290921b8401810191818101908684111562004d9e57600080fd5b8286015b84811015620045f7578035835291830191830162004da2565b600082601f83011262004dcd57600080fd5b8135602062004de06200401a8362004546565b82815260089290921b8401810191818101908684111562004e0057600080fd5b8286015b84811015620045f757610100818903121562004e205760008081fd5b62004e2a62003f2f565b813562004e378162004796565b815262004e4682860162003e8f565b85820152604062004e5981840162003e8f565b90820152606062004e6c83820162003e82565b90820152608062004e7f838201620044fa565b9082015260a062004e92838201620044fa565b9082015260c062004ea583820162003e8f565b9082015260e062004eb883820162003e82565b908201528352918301916101000162004e04565b600082601f83011262004ede57600080fd5b8135602062004ef16200401a8362004546565b82815260059290921b8401810191818101908684111562004f1157600080fd5b8286015b84811015620045f757803562004f2b8162003e5f565b835291830191830162004f15565b60008060008060008060c0878903121562004f5357600080fd5b863567ffffffffffffffff8082111562004f6c57600080fd5b62004f7a8a838b0162004d59565b9750602089013591508082111562004f9157600080fd5b62004f9f8a838b0162004dbb565b9650604089013591508082111562004fb657600080fd5b62004fc48a838b0162004ecc565b9550606089013591508082111562004fdb57600080fd5b62004fe98a838b016200456d565b945060808901359150808211156200500057600080fd5b6200500e8a838b016200456d565b935060a08901359150808211156200502557600080fd5b506200503489828a016200456d565b9150509295509295509295565b6bffffffffffffffffffffffff81811683821601908082111562003da45762003da462004699565b6bffffffffffffffffffffffff82811682821603908082111562003da45762003da462004699565b604081526000620050a6604083018662004ae1565b8281036020840152620050bb818587620046ee565b9695505050505050565b805160208083015191908110156200488d577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60209190910360031b1b16919050565b815167ffffffffffffffff81111562005125576200512562003f00565b6200513d816200513684546200483e565b8462004893565b602080601f8311600181146200519357600084156200515c5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555620048dd565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015620051e257888601518255948401946001909101908401620051c1565b50858210156200521f57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b805169ffffffffffffffffffff81168114620037cc57600080fd5b600080600080600060a086880312156200526357600080fd5b6200526e866200522f565b945060208601519350604086015192506060860151915062005293608087016200522f565b90509295509295909350565b600060208284031215620052b257600080fd5b5051919050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615620052f457620052f462004699565b500290565b60ff818116838216019081111562002b0f5762002b0f62004699565b8284823760008382016000815283516200533481836020880162004184565b0195945050505050565b600063ffffffff8083168185168183048111821515161562005364576200536462004699565b02949350505050565b600082620053a4577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfe60c060405234801561001057600080fd5b5060405161058238038061058283398101604081905261002f91610073565b600080546001600160a01b0319166001600160a01b039283161790551660805260a0526100af565b80516001600160a01b038116811461006e57600080fd5b919050565b60008060006060848603121561008857600080fd5b8351925061009860208501610057565b91506100a660408501610057565b90509250925092565b60805160a0516104a76100db6000396000610145015260008181610170015261028001526104a76000f3fe608060405234801561001057600080fd5b50600436106100725760003560e01c806379188d161161005057806379188d161461011d5780638ee489b214610140578063f00e6a2a1461016e57600080fd5b8063181f5a77146100775780631a5da6c8146100c95780635ab1bd53146100de575b600080fd5b6100b36040518060400160405280601981526020017f4175746f6d6174696f6e466f7277617264657220312e302e300000000000000081525081565b6040516100c091906102e9565b60405180910390f35b6100dc6100d7366004610355565b610194565b005b60005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100c0565b61013061012b3660046103c1565b61022c565b60405190151581526020016100c0565b6040517f000000000000000000000000000000000000000000000000000000000000000081526020016100c0565b7f00000000000000000000000000000000000000000000000000000000000000006100f8565b60005473ffffffffffffffffffffffffffffffffffffffff1633146101e5576040517fea8e4eb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b6000805473ffffffffffffffffffffffffffffffffffffffff16331461027e576040517fea8e4eb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f00000000000000000000000000000000000000000000000000000000000000005a6113888110156102af57600080fd5b6113888103905084604082048203116102c757600080fd5b50803b6102d357600080fd5b60008084516020860160008589f1949350505050565b600060208083528351808285015260005b81811015610316578581018301518582016040015282016102fa565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b60006020828403121561036757600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461038b57600080fd5b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080604083850312156103d457600080fd5b82359150602083013567ffffffffffffffff808211156103f357600080fd5b818501915085601f83011261040757600080fd5b81358181111561041957610419610392565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561045f5761045f610392565b8160405282815288602084870101111561047857600080fd5b826020860160208301376000602084830101528095505050505050925092905056fea164736f6c6343000810000a307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000810000a",
=======
	Bin: "0x6101206040523480156200001257600080fd5b5060405162005d8738038062005d87833981016040819052620000359162000386565b80816001600160a01b0316634b4fd03b6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200009b9190620003ad565b826001600160a01b031663ca30e6036040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000da573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000100919062000386565b836001600160a01b031663b10b673c6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200013f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000165919062000386565b846001600160a01b0316636709d0e56040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca919062000386565b3380600081620002215760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b038481169190911790915581161562000254576200025481620002c2565b5050508360028111156200026c576200026c620003d0565b60e0816002811115620002835762000283620003d0565b9052506001600160a01b0392831660805290821660a052811660c052601980546001600160a01b03191633179055919091166101005250620003e69050565b336001600160a01b038216036200031c5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000218565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200038357600080fd5b50565b6000602082840312156200039957600080fd5b8151620003a6816200036d565b9392505050565b600060208284031215620003c057600080fd5b815160038110620003a657600080fd5b634e487b7160e01b600052602160045260246000fd5b60805160a05160c05160e0516101005161591d6200046a600039600081816101aa01526102450152600081816102a50152818161342d01528181613663015281816138610152613a0901526000818161030901526131fb0152600081816103fd01526132df01526000818161043b01528181611fe8015261247e015261591d6000f3fe60806040523480156200001157600080fd5b5060043610620001a85760003560e01c806385c1b0ba11620000ed578063b10b673c1162000099578063ce7dc5b4116200006f578063ce7dc5b41462000460578063f2fde38b1462000477578063f7d334ba146200048e57620001a8565b8063b10b673c14620003fb578063c80480221462000422578063ca30e603146200043957620001a8565b80638e86139b11620000cf5780638e86139b14620003b1578063948108f714620003c8578063aab9edd614620003df57620001a8565b806385c1b0ba146200037b5780638da5cb5b146200039257620001a8565b80634ee88d3511620001595780636ded9eae116200012f5780636ded9eae146200032e57806371791aa0146200034557806379ba5097146200037157620001a8565b80634ee88d3514620002ca5780635147cd5914620002e15780636709d0e5146200030757620001a8565b8063349e8cca116200018f578063349e8cca146200024357806348013d7b146200028b5780634b4fd03b14620002a357620001a8565b806328f32f3814620001f057806329c5efad146200021a575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e808015620001e9573d6000f35b3d6000fd5b005b620002076200020136600462003fc6565b620004a5565b6040519081526020015b60405180910390f35b620002316200022b366004620040ac565b620007c8565b604051620002119493929190620041af565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200162000211565b62000294600081565b604051620002119190620041ff565b7f000000000000000000000000000000000000000000000000000000000000000062000294565b620001ee620002db36600462004214565b62000995565b620002f8620002f236600462004264565b620009fd565b6040516200021191906200427e565b7f000000000000000000000000000000000000000000000000000000000000000062000265565b620002076200033f36600462004295565b62000ab3565b6200035c62000356366004620040ac565b62000b5f565b60405162000211979695949392919062004348565b620001ee620012f7565b620001ee6200038c3660046200439a565b620013fa565b60005473ffffffffffffffffffffffffffffffffffffffff1662000265565b620001ee620003c236600462004427565b62002069565b620001ee620003d93660046200448a565b620022ac565b620003e8600281565b60405160ff909116815260200162000211565b7f000000000000000000000000000000000000000000000000000000000000000062000265565b620001ee6200043336600462004264565b6200254e565b7f000000000000000000000000000000000000000000000000000000000000000062000265565b620002316200047136600462004575565b62002924565b620001ee62000488366004620045ec565b620029f4565b6200035c6200049f36600462004264565b62002a0c565b6000805473ffffffffffffffffffffffffffffffffffffffff163314801590620004d95750620004d760093362002ade565b155b1562000511576040517fd48b678b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6200051c8662002b12565b90506000818a30604051620005319062003d70565b92835273ffffffffffffffffffffffffffffffffffffffff9182166020840152166040820152606001604051809103906000f08015801562000577573d6000803e3d6000fd5b5090506200065b826040518061010001604052806000151581526020018c63ffffffff16815260200163ffffffff801681526020018473ffffffffffffffffffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff168152602001600063ffffffff1681526020018d73ffffffffffffffffffffffffffffffffffffffff168152508a89898080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508b92508a915062002cae9050565b6014805474010000000000000000000000000000000000000000900463ffffffff1690806200068a836200463b565b91906101000a81548163ffffffff021916908363ffffffff16021790555050817fbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d0128a8a6040516200070392919063ffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a2817f787b2e556c364148324d7dc0cee4322a048aaf798bda184136c8636c72ed305687876040516200073f929190620046aa565b60405180910390a2817f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d566485604051620007799190620046c0565b60405180910390a2817f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf485084604051620007b39190620046c0565b60405180910390a25098975050505050505050565b60006060600080620007d9620030e7565b60008681526004602090815260409182902082516101008082018552825460ff81161515835263ffffffff91810482169483019490945265010000000000840481169482019490945273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009093048316606082015260018201546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490931660c0840152600201541660e08201525a60e0820151601454604051929450600092839273ffffffffffffffffffffffffffffffffffffffff16916c01000000000000000000000000900463ffffffff1690620008f9908b90620046d5565b60006040518083038160008787f1925050503d806000811462000939576040519150601f19603f3d011682016040523d82523d6000602084013e6200093e565b606091505b50915091505a620009509085620046f3565b9350816200096257600794506200097e565b8080602001905181019062000978919062004764565b90975095505b866200098957600494505b50505092959194509250565b620009a08362003122565b6000838152601b60205260409020620009bb82848362004859565b50827f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d56648383604051620009f0929190620046aa565b60405180910390a2505050565b6000818160045b600f81101562000a92577fff00000000000000000000000000000000000000000000000000000000000000821683826020811062000a465762000a4662004981565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161462000a7d57506000949350505050565b8062000a8981620049b0565b91505062000a04565b5081600f1a600181111562000aab5762000aab62004169565b949350505050565b600062000b53888888600089896040518060200160405280600163ffffffff1681525060405160200162000af0915163ffffffff16815260200190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181526020601f8d018190048102840181019092528b835291908c908c9081908401838280828437600092019190915250620004a592505050565b98975050505050505050565b60006060600080600080600062000b75620030e7565b600062000b828a620009fd565b905060006012604051806101200160405290816000820160009054906101000a900460ff1660ff1660ff1681526020016000820160019054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160059054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160099054906101000a900462ffffff1662ffffff1662ffffff16815260200160008201600c9054906101000a900461ffff1661ffff1661ffff16815260200160008201600e9054906101000a900460ff1615151515815260200160008201600f9054906101000a900460ff161515151581526020016000820160109054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160008201601c9054906101000a900463ffffffff1663ffffffff1663ffffffff168152505090506000600460008d8152602001908152602001600020604051806101000160405290816000820160009054906101000a900460ff161515151581526020016000820160019054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160059054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160099054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160018201600c9054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff1681526020016001820160189054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016002820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681525050905063ffffffff8016816040015163ffffffff161462000f0a576000604051806020016040528060008152506001600084602001516000808263ffffffff1692509950995099509950995099509950505050620012eb565b80511562000f50576000604051806020016040528060008152506002600084602001516000808263ffffffff1692509950995099509950995099509950505050620012eb565b62000f5b82620031d8565b602083015160155492975090955060009162000f8a918591640100000000900463ffffffff16898986620033ca565b9050806bffffffffffffffffffffffff168260a001516bffffffffffffffffffffffff16101562000ff4576000604051806020016040528060008152506006600085602001516000808263ffffffff1692509a509a509a509a509a509a509a5050505050620012eb565b60019a506060600085600181111562001011576200101162004169565b03620010d1576040517f6e04ff0d00000000000000000000000000000000000000000000000000000000906200104c908f90602401620046c0565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152905062001186565b6040517fbe61b775000000000000000000000000000000000000000000000000000000009062001106908f90602401620046c0565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915290505b5a98508260e0015173ffffffffffffffffffffffffffffffffffffffff166013600101600c9054906101000a900463ffffffff1663ffffffff1682604051620011d09190620046d5565b60006040518083038160008787f1925050503d806000811462001210576040519150601f19603f3d011682016040523d82523d6000602084013e62001215565b606091505b50909c509a505a62001228908a620046f3565b98508b6200123a57600399506200128d565b8a80602001905181019062001250919062004764565b909c509a508b6200128d57505060408051602080820190925260008082529290910151919a5098506004975063ffffffff169450620012eb915050565b6015548b5164010000000090910463ffffffff161015620012da57505060408051602080820190925260008082529290910151919a5098506005975063ffffffff169450620012eb915050565b50506020015163ffffffff16945050505b92959891949750929550565b60015473ffffffffffffffffffffffffffffffffffffffff1633146200137e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600173ffffffffffffffffffffffffffffffffffffffff82166000908152601a602052604090205460ff16600381111562001439576200143962004169565b14158015620014855750600373ffffffffffffffffffffffffffffffffffffffff82166000908152601a602052604090205460ff16600381111562001482576200148262004169565b14155b15620014bd576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6013546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff166200151d576040517fd12d7d8d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082900362001559576040517f2c2fc94100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805161010081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526000808567ffffffffffffffff811115620015b857620015b862003e73565b604051908082528060200260200182016040528015620015e2578160200160208202803683370190505b50905060008667ffffffffffffffff81111562001603576200160362003e73565b6040519080825280602002602001820160405280156200169257816020015b604080516101008101825260008082526020808301829052928201819052606082018190526080820181905260a0820181905260c0820181905260e082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181620016225790505b50905060008767ffffffffffffffff811115620016b357620016b362003e73565b604051908082528060200260200182016040528015620016e857816020015b6060815260200190600190039081620016d25790505b50905060008867ffffffffffffffff81111562001709576200170962003e73565b6040519080825280602002602001820160405280156200173e57816020015b6060815260200190600190039081620017285790505b50905060008967ffffffffffffffff8111156200175f576200175f62003e73565b6040519080825280602002602001820160405280156200179457816020015b60608152602001906001900390816200177e5790505b50905060005b8a81101562001daf578b8b82818110620017b857620017b862004981565b6020908102929092013560008181526004845260409081902081516101008082018452825460ff81161515835263ffffffff91810482169783019790975265010000000000870481169382019390935273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009096048616606082015260018201546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490921660c08301526002015490931660e08401529a50909850620018a490508962003122565b60608801516040517f1a5da6c800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8c8116600483015290911690631a5da6c890602401600060405180830381600087803b1580156200191457600080fd5b505af115801562001929573d6000803e3d6000fd5b505050508785828151811062001943576200194362004981565b6020026020010181905250600560008a815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1686828151811062001997576200199762004981565b73ffffffffffffffffffffffffffffffffffffffff90921660209283029190910182015260008a81526007909152604090208054620019d690620047b1565b80601f016020809104026020016040519081016040528092919081815260200182805462001a0490620047b1565b801562001a555780601f1062001a295761010080835404028352916020019162001a55565b820191906000526020600020905b81548152906001019060200180831162001a3757829003601f168201915b505050505084828151811062001a6f5762001a6f62004981565b6020026020010181905250601b60008a8152602001908152602001600020805462001a9a90620047b1565b80601f016020809104026020016040519081016040528092919081815260200182805462001ac890620047b1565b801562001b195780601f1062001aed5761010080835404028352916020019162001b19565b820191906000526020600020905b81548152906001019060200180831162001afb57829003601f168201915b505050505083828151811062001b335762001b3362004981565b6020026020010181905250601c60008a8152602001908152602001600020805462001b5e90620047b1565b80601f016020809104026020016040519081016040528092919081815260200182805462001b8c90620047b1565b801562001bdd5780601f1062001bb15761010080835404028352916020019162001bdd565b820191906000526020600020905b81548152906001019060200180831162001bbf57829003601f168201915b505050505082828151811062001bf75762001bf762004981565b60200260200101819052508760a001516bffffffffffffffffffffffff168762001c229190620049eb565b60008a815260046020908152604080832080547fffffff00000000000000000000000000000000000000000000000000000000001681556001810180547fffffffff0000000000000000000000000000000000000000000000000000000016905560020180547fffffffffffffffffffffffff00000000000000000000000000000000000000001690556007909152812091985062001cc2919062003d7e565b6000898152601b6020526040812062001cdb9162003d7e565b6000898152601c6020526040812062001cf49162003d7e565b600089815260066020526040902080547fffffffffffffffffffffffff000000000000000000000000000000000000000016905562001d3560028a62003419565b5060a0880151604080516bffffffffffffffffffffffff909216825273ffffffffffffffffffffffffffffffffffffffff8c1660208301528a917fb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff910160405180910390a28062001da681620049b0565b9150506200179a565b508560185462001dc09190620046f3565b60185560405160009062001de5908d908d9088908a9089908990899060200162004aad565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905260135490915073ffffffffffffffffffffffffffffffffffffffff808c1691638e86139b916c010000000000000000000000009091041663c71249ab60028e73ffffffffffffffffffffffffffffffffffffffff1663aab9edd66040518163ffffffff1660e01b81526004016020604051808303816000875af115801562001e9f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001ec5919062004c19565b866040518463ffffffff1660e01b815260040162001ee69392919062004c3e565b600060405180830381865afa15801562001f04573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405262001f4c919081019062004c6e565b6040518263ffffffff1660e01b815260040162001f6a9190620046c0565b600060405180830381600087803b15801562001f8557600080fd5b505af115801562001f9a573d6000803e3d6000fd5b50506040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8d81166004830152602482018b90527f000000000000000000000000000000000000000000000000000000000000000016925063a9059cbb91506044016020604051808303816000875af115801562002034573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200205a919062004ca7565b50505050505050505050505050565b6002336000908152601a602052604090205460ff16600381111562002092576200209262004169565b14158015620020c857506003336000908152601a602052604090205460ff166003811115620020c557620020c562004169565b14155b1562002100576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008080808080620021158789018962004ea7565b95509550955095509550955060005b8651811015620022a157620021e687828151811062002147576200214762004981565b602002602001015187838151811062002164576200216462004981565b602002602001015187848151811062002181576200218162004981565b60200260200101518785815181106200219e576200219e62004981565b6020026020010151878681518110620021bb57620021bb62004981565b6020026020010151878781518110620021d857620021d862004981565b602002602001015162002cae565b868181518110620021fb57620021fb62004981565b60200260200101517f74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a7187838151811062002239576200223962004981565b602002602001015160a0015133604051620022849291906bffffffffffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a2806200229881620049b0565b91505062002124565b505050505050505050565b60008281526004602090815260409182902082516101008082018552825460ff81161515835263ffffffff918104821694830194909452650100000000008404811694820185905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009094048416606083015260018301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a084015278010000000000000000000000000000000000000000000000009004811660c083015260029092015490921660e0830152909114620023b9576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818160a00151620023cb919062004faf565b600084815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff938416021790556018546200243391841690620049eb565b6018556040517f23b872dd0000000000000000000000000000000000000000000000000000000081523360048201523060248201526bffffffffffffffffffffffff831660448201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906323b872dd906064016020604051808303816000875af1158015620024dd573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062002503919062004ca7565b506040516bffffffffffffffffffffffff83168152339084907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a3505050565b600081815260046020908152604080832081516101008082018452825460ff81161515835263ffffffff918104821695830195909552650100000000008504811693820184905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009095048516606083015260018301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a084015278010000000000000000000000000000000000000000000000009004811660c083015260029092015490931660e08401529192911415906200264660005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16149050818015620026a157508080156200269f57506200269262003427565b836040015163ffffffff16115b155b15620026d9576040517ffbc0357800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b801580156200270c575060008481526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314155b1562002744576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006200275062003427565b905081620027685762002765603282620049eb565b90505b6000858152600460205260409020805463ffffffff80841665010000000000027fffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffffff90921691909117909155620027c49060029087906200341916565b5060135460808501516bffffffffffffffffffffffff91821691600091168211156200282d576080860151620027fb908362004fd7565b90508560a001516bffffffffffffffffffffffff16816bffffffffffffffffffffffff1611156200282d575060a08501515b808660a001516200283f919062004fd7565b600088815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff93841602179055601454620028a79183911662004faf565b601480547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff9290921691909117905560405167ffffffffffffffff84169088907f91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f79118190600090a350505050505050565b6000606060008062002935620030e7565b6000634b56a42e60e01b888888604051602401620029569392919062004fff565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091529050620029e18982620007c8565b929c919b50995090975095505050505050565b620029fe620034e3565b62002a098162003566565b50565b60006060600080600080600062002ac788600760008b8152602001908152602001600020805462002a3d90620047b1565b80601f016020809104026020016040519081016040528092919081815260200182805462002a6b90620047b1565b801562002abc5780601f1062002a905761010080835404028352916020019162002abc565b820191906000526020600020905b81548152906001019060200180831162002a9e57829003601f168201915b505050505062000b5f565b959e949d50929b5090995097509550909350915050565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415155b90505b92915050565b600080600062002b39600162002b2762003427565b62002b339190620046f3565b6200365d565b601454604080516020810193909352309083015274010000000000000000000000000000000000000000900463ffffffff166060820152608001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201209083015201604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052905060045b600f81101562002c45578282828151811062002c015762002c0162004981565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508062002c3c81620049b0565b91505062002be1565b5083600181111562002c5b5762002c5b62004169565b60f81b81600f8151811062002c745762002c7462004981565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535062000aab8162005033565b6012546e010000000000000000000000000000900460ff161562002cfe576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60e085015173ffffffffffffffffffffffffffffffffffffffff163b62002d51576040517f09ee12d500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601554835163ffffffff909116101562002d97576040517f2b49197900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108fc856020015163ffffffff16108062002dd55750601454602086015163ffffffff70010000000000000000000000000000000090920482169116115b1562002e0d576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008681526004602052604090206002015473ffffffffffffffffffffffffffffffffffffffff161562002e6d576040517f6e3b930b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600086815260046020908152604080832088518154848b0151848c015160608d015173ffffffffffffffffffffffffffffffffffffffff9081166901000000000000000000027fffffff0000000000000000000000000000000000000000ffffffffffffffffff63ffffffff9384166501000000000002167fffffff000000000000000000000000000000000000000000000000ffffffffff948416610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000ff971515979097167fffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000009096169590951795909517929092169290921792909217835560808b015160018401805460a08e015160c08f01519094167801000000000000000000000000000000000000000000000000027fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff6bffffffffffffffffffffffff9586166c01000000000000000000000000027fffffffffffffffff0000000000000000000000000000000000000000000000009093169590941694909417179190911691909117905560e08a0151600290920180549282167fffffffffffffffffffffffff0000000000000000000000000000000000000000938416179055600584528285208054918a16919092161790556007909152902062003075848262005076565b508460a001516bffffffffffffffffffffffff16601854620030989190620049eb565b6018556000868152601b60205260409020620030b5838262005076565b506000868152601c60205260409020620030d0828262005076565b50620030de600287620037cc565b50505050505050565b321562003120576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60008181526005602052604090205473ffffffffffffffffffffffffffffffffffffffff16331462003180576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526004602052604090205465010000000000900463ffffffff9081161462002a09576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000806000836060015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa15801562003265573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200328b9190620051b8565b5094509092505050600081131580620032a357508142105b80620032c85750828015620032c85750620032bf8242620046f3565b8463ffffffff16105b15620032d9576016549550620032dd565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa15801562003349573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200336f9190620051b8565b50945090925050506000811315806200338757508142105b80620033ac5750828015620033ac5750620033a38242620046f3565b8463ffffffff16105b15620033bd576017549450620033c1565b8094505b50505050915091565b600080620033dd868960000151620037da565b9050600080620033fa8a8a63ffffffff16858a8a60018b62003828565b90925090506200340b818362004faf565b9a9950505050505050505050565b600062002b09838362003c13565b600060017f0000000000000000000000000000000000000000000000000000000000000000600281111562003460576200346062004169565b03620034de57606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620034b3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620034d991906200520d565b905090565b504390565b60005473ffffffffffffffffffffffffffffffffffffffff16331462003120576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162001375565b3373ffffffffffffffffffffffffffffffffffffffff821603620035e7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162001375565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060017f0000000000000000000000000000000000000000000000000000000000000000600281111562003696576200369662004169565b03620037c2576000606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620036eb573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200371191906200520d565b905080831015806200372f57506101006200372d8483620046f3565b115b156200373e5750600092915050565b6040517f2b407a8200000000000000000000000000000000000000000000000000000000815260048101849052606490632b407a8290602401602060405180830381865afa15801562003795573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620037bb91906200520d565b9392505050565b504090565b919050565b600062002b09838362003d1e565b6000620037ef63ffffffff8416601462005227565b620037fc83600162005267565b6200380d9060ff16611d4c62005227565b6200381c9062013880620049eb565b62002b099190620049eb565b6000806000896080015161ffff168762003843919062005227565b9050838015620038525750803a105b156200385b57503a5b600060027f0000000000000000000000000000000000000000000000000000000000000000600281111562003894576200389462004169565b0362003a05576040805160008152602081019091528515620038f857600036604051806080016040528060488152602001620058c960489139604051602001620038e19392919062005283565b604051602081830303815290604052905062003966565b6015546200391690640100000000900463ffffffff166004620052ac565b63ffffffff1667ffffffffffffffff81111562003937576200393762003e73565b6040519080825280601f01601f19166020018201604052801562003962576020820181803683370190505b5090505b6040517f49948e0e00000000000000000000000000000000000000000000000000000000815273420000000000000000000000000000000000000f906349948e0e90620039b8908490600401620046c0565b602060405180830381865afa158015620039d6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620039fc91906200520d565b91505062003ab8565b60017f0000000000000000000000000000000000000000000000000000000000000000600281111562003a3c5762003a3c62004169565b0362003ab857606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562003a8f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003ab591906200520d565b90505b8462003ad757808b6080015161ffff1662003ad4919062005227565b90505b62003ae761ffff871682620052db565b90506000878262003af98c8e620049eb565b62003b05908662005227565b62003b119190620049eb565b62003b2590670de0b6b3a764000062005227565b62003b319190620052db565b905060008c6040015163ffffffff1664e8d4a5100062003b52919062005227565b898e6020015163ffffffff16858f8862003b6d919062005227565b62003b799190620049eb565b62003b8990633b9aca0062005227565b62003b95919062005227565b62003ba19190620052db565b62003bad9190620049eb565b90506b033b2e3c9fd0803ce800000062003bc88284620049eb565b111562003c01576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b6000818152600183016020526040812054801562003d0c57600062003c3a600183620046f3565b855490915060009062003c5090600190620046f3565b905081811462003cbc57600086600001828154811062003c745762003c7462004981565b906000526020600020015490508087600001848154811062003c9a5762003c9a62004981565b6000918252602080832090910192909255918252600188019052604090208390555b855486908062003cd05762003cd062005317565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505062002b0c565b600091505062002b0c565b5092915050565b600081815260018301602052604081205462003d675750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915562002b0c565b50600062002b0c565b610582806200534783390190565b50805462003d8c90620047b1565b6000825580601f1062003d9d575050565b601f01602090049060005260206000209081019062002a0991905b8082111562003dce576000815560010162003db8565b5090565b73ffffffffffffffffffffffffffffffffffffffff8116811462002a0957600080fd5b8035620037c78162003dd2565b803563ffffffff81168114620037c757600080fd5b803560028110620037c757600080fd5b60008083601f84011262003e3a57600080fd5b50813567ffffffffffffffff81111562003e5357600080fd5b60208301915083602082850101111562003e6c57600080fd5b9250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610100810167ffffffffffffffff8111828210171562003ec95762003ec962003e73565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171562003f195762003f1962003e73565b604052919050565b600067ffffffffffffffff82111562003f3e5762003f3e62003e73565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f83011262003f7c57600080fd5b813562003f9362003f8d8262003f21565b62003ecf565b81815284602083860101111562003fa957600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060008060e0898b03121562003fe357600080fd5b883562003ff08162003dd2565b97506200400060208a0162003e02565b96506040890135620040128162003dd2565b95506200402260608a0162003e17565b9450608089013567ffffffffffffffff808211156200404057600080fd5b6200404e8c838d0162003e27565b909650945060a08b01359150808211156200406857600080fd5b620040768c838d0162003f6a565b935060c08b01359150808211156200408d57600080fd5b506200409c8b828c0162003f6a565b9150509295985092959890939650565b60008060408385031215620040c057600080fd5b82359150602083013567ffffffffffffffff811115620040df57600080fd5b620040ed8582860162003f6a565b9150509250929050565b60005b8381101562004114578181015183820152602001620040fa565b50506000910152565b6000815180845262004137816020860160208601620040f7565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b60088110620041ab57620041ab62004169565b9052565b8415158152608060208201526000620041cc60808301866200411d565b9050620041dd604083018562004198565b82606083015295945050505050565b6003811062002a095762002a0962004169565b602081016200420e83620041ec565b91905290565b6000806000604084860312156200422a57600080fd5b83359250602084013567ffffffffffffffff8111156200424957600080fd5b620042578682870162003e27565b9497909650939450505050565b6000602082840312156200427757600080fd5b5035919050565b60208101600283106200420e576200420e62004169565b600080600080600080600060a0888a031215620042b157600080fd5b8735620042be8162003dd2565b9650620042ce6020890162003e02565b95506040880135620042e08162003dd2565b9450606088013567ffffffffffffffff80821115620042fe57600080fd5b6200430c8b838c0162003e27565b909650945060808a01359150808211156200432657600080fd5b50620043358a828b0162003e27565b989b979a50959850939692959293505050565b871515815260e0602082015260006200436560e08301896200411d565b905062004376604083018862004198565b8560608301528460808301528360a08301528260c083015298975050505050505050565b600080600060408486031215620043b057600080fd5b833567ffffffffffffffff80821115620043c957600080fd5b818601915086601f830112620043de57600080fd5b813581811115620043ee57600080fd5b8760208260051b85010111156200440457600080fd5b602092830195509350508401356200441c8162003dd2565b809150509250925092565b600080602083850312156200443b57600080fd5b823567ffffffffffffffff8111156200445357600080fd5b620044618582860162003e27565b90969095509350505050565b80356bffffffffffffffffffffffff81168114620037c757600080fd5b600080604083850312156200449e57600080fd5b82359150620044b0602084016200446d565b90509250929050565b600067ffffffffffffffff821115620044d657620044d662003e73565b5060051b60200190565b600082601f830112620044f257600080fd5b813560206200450562003f8d83620044b9565b82815260059290921b840181019181810190868411156200452557600080fd5b8286015b848110156200456a57803567ffffffffffffffff8111156200454b5760008081fd5b6200455b8986838b010162003f6a565b84525091830191830162004529565b509695505050505050565b600080600080606085870312156200458c57600080fd5b84359350602085013567ffffffffffffffff80821115620045ac57600080fd5b620045ba88838901620044e0565b94506040870135915080821115620045d157600080fd5b50620045e08782880162003e27565b95989497509550505050565b600060208284031215620045ff57600080fd5b8135620037bb8162003dd2565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600063ffffffff8083168181036200465757620046576200460c565b6001019392505050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b60208152600062000aab60208301848662004661565b60208152600062002b0960208301846200411d565b60008251620046e9818460208701620040f7565b9190910192915050565b8181038181111562002b0c5762002b0c6200460c565b801515811462002a0957600080fd5b600082601f8301126200472a57600080fd5b81516200473b62003f8d8262003f21565b8181528460208386010111156200475157600080fd5b62000aab826020830160208701620040f7565b600080604083850312156200477857600080fd5b8251620047858162004709565b602084015190925067ffffffffffffffff811115620047a357600080fd5b620040ed8582860162004718565b600181811c90821680620047c657607f821691505b60208210810362004800577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f8211156200485457600081815260208120601f850160051c810160208610156200482f5750805b601f850160051c820191505b8181101562004850578281556001016200483b565b5050505b505050565b67ffffffffffffffff83111562004874576200487462003e73565b6200488c83620048858354620047b1565b8362004806565b6000601f841160018114620048e15760008515620048aa5750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b1783556200497a565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b8281101562004932578685013582556020948501946001909201910162004910565b50868210156200496e577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203620049e457620049e46200460c565b5060010190565b8082018082111562002b0c5762002b0c6200460c565b600081518084526020808501945080840160005b8381101562004a4957815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010162004a15565b509495945050505050565b600081518084526020808501808196508360051b8101915082860160005b8581101562004aa057828403895262004a8d8483516200411d565b9885019893509084019060010162004a72565b5091979650505050505050565b600060c0808352888184015260e07f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8a111562004ae957600080fd5b8960051b808c83870137840184810382016020808701919091528a518383018190528b82019261010092919083019060005b8181101562004bbc5785518051151584528481015163ffffffff9081168686015260408083015182169086015260608083015173ffffffffffffffffffffffffffffffffffffffff908116918701919091526080808401516bffffffffffffffffffffffff9081169188019190915260a080850151909116908701528a8301519091168a860152908801511687840152948301949184019160010162004b1b565b5050878103604089015262004bd2818d62004a01565b95505050505050828103606084015262004bed818762004a54565b9050828103608084015262004c03818662004a54565b905082810360a08401526200340b818562004a54565b60006020828403121562004c2c57600080fd5b815160ff81168114620037bb57600080fd5b60ff8416815260ff8316602082015260606040820152600062004c6560608301846200411d565b95945050505050565b60006020828403121562004c8157600080fd5b815167ffffffffffffffff81111562004c9957600080fd5b62000aab8482850162004718565b60006020828403121562004cba57600080fd5b8151620037bb8162004709565b600082601f83011262004cd957600080fd5b8135602062004cec62003f8d83620044b9565b82815260059290921b8401810191818101908684111562004d0c57600080fd5b8286015b848110156200456a578035835291830191830162004d10565b600082601f83011262004d3b57600080fd5b8135602062004d4e62003f8d83620044b9565b82815260089290921b8401810191818101908684111562004d6e57600080fd5b8286015b848110156200456a57610100818903121562004d8e5760008081fd5b62004d9862003ea2565b813562004da58162004709565b815262004db482860162003e02565b85820152604062004dc781840162003e02565b90820152606062004dda83820162003df5565b90820152608062004ded8382016200446d565b9082015260a062004e008382016200446d565b9082015260c062004e1383820162003e02565b9082015260e062004e2683820162003df5565b908201528352918301916101000162004d72565b600082601f83011262004e4c57600080fd5b8135602062004e5f62003f8d83620044b9565b82815260059290921b8401810191818101908684111562004e7f57600080fd5b8286015b848110156200456a57803562004e998162003dd2565b835291830191830162004e83565b60008060008060008060c0878903121562004ec157600080fd5b863567ffffffffffffffff8082111562004eda57600080fd5b62004ee88a838b0162004cc7565b9750602089013591508082111562004eff57600080fd5b62004f0d8a838b0162004d29565b9650604089013591508082111562004f2457600080fd5b62004f328a838b0162004e3a565b9550606089013591508082111562004f4957600080fd5b62004f578a838b01620044e0565b9450608089013591508082111562004f6e57600080fd5b62004f7c8a838b01620044e0565b935060a089013591508082111562004f9357600080fd5b5062004fa289828a01620044e0565b9150509295509295509295565b6bffffffffffffffffffffffff81811683821601908082111562003d175762003d176200460c565b6bffffffffffffffffffffffff82811682821603908082111562003d175762003d176200460c565b60408152600062005014604083018662004a54565b82810360208401526200502981858762004661565b9695505050505050565b8051602080830151919081101562004800577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60209190910360031b1b16919050565b815167ffffffffffffffff81111562005093576200509362003e73565b620050ab81620050a48454620047b1565b8462004806565b602080601f831160018114620051015760008415620050ca5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b17855562004850565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101562005150578886015182559484019460019091019084016200512f565b50858210156200518d57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b805169ffffffffffffffffffff81168114620037c757600080fd5b600080600080600060a08688031215620051d157600080fd5b620051dc866200519d565b945060208601519350604086015192506060860151915062005201608087016200519d565b90509295509295909350565b6000602082840312156200522057600080fd5b5051919050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156200526257620052626200460c565b500290565b60ff818116838216019081111562002b0c5762002b0c6200460c565b828482376000838201600081528351620052a2818360208801620040f7565b0195945050505050565b600063ffffffff80831681851681830481118215151615620052d257620052d26200460c565b02949350505050565b60008262005312577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfe60c060405234801561001057600080fd5b5060405161058238038061058283398101604081905261002f91610073565b600080546001600160a01b0319166001600160a01b039283161790551660805260a0526100af565b80516001600160a01b038116811461006e57600080fd5b919050565b60008060006060848603121561008857600080fd5b8351925061009860208501610057565b91506100a660408501610057565b90509250925092565b60805160a0516104a76100db6000396000610145015260008181610170015261028001526104a76000f3fe608060405234801561001057600080fd5b50600436106100725760003560e01c806379188d161161005057806379188d161461011d5780638ee489b214610140578063f00e6a2a1461016e57600080fd5b8063181f5a77146100775780631a5da6c8146100c95780635ab1bd53146100de575b600080fd5b6100b36040518060400160405280601981526020017f4175746f6d6174696f6e466f7277617264657220312e302e300000000000000081525081565b6040516100c091906102e9565b60405180910390f35b6100dc6100d7366004610355565b610194565b005b60005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100c0565b61013061012b3660046103c1565b61022c565b60405190151581526020016100c0565b6040517f000000000000000000000000000000000000000000000000000000000000000081526020016100c0565b7f00000000000000000000000000000000000000000000000000000000000000006100f8565b60005473ffffffffffffffffffffffffffffffffffffffff1633146101e5576040517fea8e4eb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b6000805473ffffffffffffffffffffffffffffffffffffffff16331461027e576040517fea8e4eb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f00000000000000000000000000000000000000000000000000000000000000005a6113888110156102af57600080fd5b6113888103905084604082048203116102c757600080fd5b50803b6102d357600080fd5b60008084516020860160008589f1949350505050565b600060208083528351808285015260005b81811015610316578581018301518582016040015282016102fa565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b60006020828403121561036757600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461038b57600080fd5b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080604083850312156103d457600080fd5b82359150602083013567ffffffffffffffff808211156103f357600080fd5b818501915085601f83011261040757600080fd5b81358181111561041957610419610392565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561045f5761045f610392565b8160405282815288602084870101111561047857600080fd5b826020860160208301376000602084830101528095505050505050925092905056fea164736f6c6343000810000a307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000810000a",
>>>>>>> develop
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
	return address, tx, &KeeperRegistryLogicA{KeeperRegistryLogicACaller: KeeperRegistryLogicACaller{contract: contract}, KeeperRegistryLogicATransactor: KeeperRegistryLogicATransactor{contract: contract}, KeeperRegistryLogicAFilterer: KeeperRegistryLogicAFilterer{contract: contract}}, nil
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicACaller) GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogicA.contract.Call(opts, &out, "getFastGasFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) GetFastGasFeedAddress() (common.Address, error) {
	return _KeeperRegistryLogicA.Contract.GetFastGasFeedAddress(&_KeeperRegistryLogicA.CallOpts)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicACallerSession) GetFastGasFeedAddress() (common.Address, error) {
	return _KeeperRegistryLogicA.Contract.GetFastGasFeedAddress(&_KeeperRegistryLogicA.CallOpts)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicACaller) GetLinkAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogicA.contract.Call(opts, &out, "getLinkAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) GetLinkAddress() (common.Address, error) {
	return _KeeperRegistryLogicA.Contract.GetLinkAddress(&_KeeperRegistryLogicA.CallOpts)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicACallerSession) GetLinkAddress() (common.Address, error) {
	return _KeeperRegistryLogicA.Contract.GetLinkAddress(&_KeeperRegistryLogicA.CallOpts)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicACaller) GetLinkNativeFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogicA.contract.Call(opts, &out, "getLinkNativeFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) GetLinkNativeFeedAddress() (common.Address, error) {
	return _KeeperRegistryLogicA.Contract.GetLinkNativeFeedAddress(&_KeeperRegistryLogicA.CallOpts)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicACallerSession) GetLinkNativeFeedAddress() (common.Address, error) {
	return _KeeperRegistryLogicA.Contract.GetLinkNativeFeedAddress(&_KeeperRegistryLogicA.CallOpts)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicACaller) GetMode(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistryLogicA.contract.Call(opts, &out, "getMode")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) GetMode() (uint8, error) {
	return _KeeperRegistryLogicA.Contract.GetMode(&_KeeperRegistryLogicA.CallOpts)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicACallerSession) GetMode() (uint8, error) {
	return _KeeperRegistryLogicA.Contract.GetMode(&_KeeperRegistryLogicA.CallOpts)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicACaller) GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistryLogicA.contract.Call(opts, &out, "getTriggerType", upkeepId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _KeeperRegistryLogicA.Contract.GetTriggerType(&_KeeperRegistryLogicA.CallOpts, upkeepId)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicACallerSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _KeeperRegistryLogicA.Contract.GetTriggerType(&_KeeperRegistryLogicA.CallOpts, upkeepId)
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicACaller) UpkeepTranscoderVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistryLogicA.contract.Call(opts, &out, "upkeepTranscoderVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) UpkeepTranscoderVersion() (uint8, error) {
	return _KeeperRegistryLogicA.Contract.UpkeepTranscoderVersion(&_KeeperRegistryLogicA.CallOpts)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicACallerSession) UpkeepTranscoderVersion() (uint8, error) {
	return _KeeperRegistryLogicA.Contract.UpkeepTranscoderVersion(&_KeeperRegistryLogicA.CallOpts)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicACaller) UpkeepVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistryLogicA.contract.Call(opts, &out, "upkeepVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) UpkeepVersion() (uint8, error) {
	return _KeeperRegistryLogicA.Contract.UpkeepVersion(&_KeeperRegistryLogicA.CallOpts)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicACallerSession) UpkeepVersion() (uint8, error) {
	return _KeeperRegistryLogicA.Contract.UpkeepVersion(&_KeeperRegistryLogicA.CallOpts)
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) CheckUpkeep(opts *bind.TransactOpts, id *big.Int, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "checkUpkeep", id, checkData)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) CheckUpkeep(id *big.Int, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.CheckUpkeep(&_KeeperRegistryLogicA.TransactOpts, id, checkData)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) CheckUpkeep(id *big.Int, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.CheckUpkeep(&_KeeperRegistryLogicA.TransactOpts, id, checkData)
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

type KeeperRegistryLogicAUpkeepAdminOffchainConfigSetIterator struct {
	Event *KeeperRegistryLogicAUpkeepAdminOffchainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAUpkeepAdminOffchainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAUpkeepAdminOffchainConfigSet)
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
		it.Event = new(KeeperRegistryLogicAUpkeepAdminOffchainConfigSet)
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

func (it *KeeperRegistryLogicAUpkeepAdminOffchainConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAUpkeepAdminOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAUpkeepAdminOffchainConfigSet struct {
	Id                  *big.Int
	AdminOffchainConfig []byte
	Raw                 types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterUpkeepAdminOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepAdminOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "UpkeepAdminOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAUpkeepAdminOffchainConfigSetIterator{contract: _KeeperRegistryLogicA.contract, event: "UpkeepAdminOffchainConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchUpkeepAdminOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepAdminOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "UpkeepAdminOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAUpkeepAdminOffchainConfigSet)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepAdminOffchainConfigSet", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseUpkeepAdminOffchainConfigSet(log types.Log) (*KeeperRegistryLogicAUpkeepAdminOffchainConfigSet, error) {
	event := new(KeeperRegistryLogicAUpkeepAdminOffchainConfigSet)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepAdminOffchainConfigSet", log); err != nil {
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

type KeeperRegistryLogicAUpkeepPipelineDataSetIterator struct {
	Event *KeeperRegistryLogicAUpkeepPipelineDataSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAUpkeepPipelineDataSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAUpkeepPipelineDataSet)
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
		it.Event = new(KeeperRegistryLogicAUpkeepPipelineDataSet)
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

func (it *KeeperRegistryLogicAUpkeepPipelineDataSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAUpkeepPipelineDataSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAUpkeepPipelineDataSet struct {
	Id              *big.Int
	NewPipelineData []byte
	Raw             types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterUpkeepPipelineDataSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepPipelineDataSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "UpkeepPipelineDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAUpkeepPipelineDataSetIterator{contract: _KeeperRegistryLogicA.contract, event: "UpkeepPipelineDataSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchUpkeepPipelineDataSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepPipelineDataSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "UpkeepPipelineDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAUpkeepPipelineDataSet)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepPipelineDataSet", log); err != nil {
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseUpkeepPipelineDataSet(log types.Log) (*KeeperRegistryLogicAUpkeepPipelineDataSet, error) {
	event := new(KeeperRegistryLogicAUpkeepPipelineDataSet)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepPipelineDataSet", log); err != nil {
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
	ExecuteGas uint32
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
	case _KeeperRegistryLogicA.abi.Events["CancelledUpkeepReport"].ID:
		return _KeeperRegistryLogicA.ParseCancelledUpkeepReport(log)
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
	case _KeeperRegistryLogicA.abi.Events["UpkeepAdminOffchainConfigSet"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepAdminOffchainConfigSet(log)
	case _KeeperRegistryLogicA.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepAdminTransferRequested(log)
	case _KeeperRegistryLogicA.abi.Events["UpkeepAdminTransferred"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepAdminTransferred(log)
	case _KeeperRegistryLogicA.abi.Events["UpkeepCanceled"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepCanceled(log)
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
	case _KeeperRegistryLogicA.abi.Events["UpkeepPipelineDataSet"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepPipelineDataSet(log)
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

func (KeeperRegistryLogicACancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd5636")
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

func (KeeperRegistryLogicAUpkeepAdminOffchainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x09a658476c5597979b9948f488ec2958cfead97bc8f46b19ca0b21cdab93cdee")
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

func (KeeperRegistryLogicAUpkeepPipelineDataSet) Topic() common.Hash {
	return common.HexToHash("0x787b2e556c364148324d7dc0cee4322a048aaf798bda184136c8636c72ed3056")
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

	GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error)

	GetLinkAddress(opts *bind.CallOpts) (common.Address, error)

	GetLinkNativeFeedAddress(opts *bind.CallOpts) (common.Address, error)

	GetMode(opts *bind.CallOpts) (uint8, error)

	GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	UpkeepTranscoderVersion(opts *bind.CallOpts) (uint8, error)

	UpkeepVersion(opts *bind.CallOpts) (uint8, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error)

	CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	CheckCallback(opts *bind.TransactOpts, id *big.Int, values [][]byte, extraData []byte) (*types.Transaction, error)

	CheckUpkeep(opts *bind.TransactOpts, id *big.Int, checkData []byte) (*types.Transaction, error)

	CheckUpkeep0(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	ExecuteCallback(opts *bind.TransactOpts, id *big.Int, payload []byte) (*types.Transaction, error)

	MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error)

	ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error)

	RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, triggerType uint8, checkData []byte, triggerConfig []byte, offchainConfig []byte) (*types.Transaction, error)

	RegisterUpkeep0(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error)

	SetUpkeepTriggerConfig(opts *bind.TransactOpts, id *big.Int, triggerConfig []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error)

	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicACancelledUpkeepReportIterator, error)

	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicACancelledUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseCancelledUpkeepReport(log types.Log) (*KeeperRegistryLogicACancelledUpkeepReport, error)

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

	FilterUpkeepAdminOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepAdminOffchainConfigSetIterator, error)

	WatchUpkeepAdminOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepAdminOffchainConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepAdminOffchainConfigSet(log types.Log) (*KeeperRegistryLogicAUpkeepAdminOffchainConfigSet, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogicAUpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*KeeperRegistryLogicAUpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogicAUpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*KeeperRegistryLogicAUpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistryLogicAUpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*KeeperRegistryLogicAUpkeepCanceled, error)

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

	FilterUpkeepPipelineDataSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepPipelineDataSetIterator, error)

	WatchUpkeepPipelineDataSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepPipelineDataSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPipelineDataSet(log types.Log) (*KeeperRegistryLogicAUpkeepPipelineDataSet, error)

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
