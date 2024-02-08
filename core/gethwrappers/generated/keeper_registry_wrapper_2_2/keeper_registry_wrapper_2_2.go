// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package keeper_registry_wrapper_2_2

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

type AutomationRegistryBase22OnchainConfig struct {
	PaymentPremiumPPB      uint32
	FlatFeeMicroLink       uint32
	CheckGasLimit          uint32
	StalenessSeconds       *big.Int
	GasCeilingMultiplier   uint16
	MinUpkeepSpend         *big.Int
	MaxPerformGas          uint32
	MaxCheckDataSize       uint32
	MaxPerformDataSize     uint32
	MaxRevertDataSize      uint32
	FallbackGasPrice       *big.Int
	FallbackLinkPrice      *big.Int
	Transcoder             common.Address
	Registrars             []common.Address
	UpkeepPrivilegeManager common.Address
	ReorgProtectionEnabled bool
}

var AutomationRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractAutomationRegistryLogicB2_2\",\"name\":\"logicA\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTrigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTriggerType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"AdminPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"DedupKeyAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepTriggerConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fallbackTo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfigBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"reorgProtectionEnabled\",\"type\":\"bool\"}],\"internalType\":\"structAutomationRegistryBase2_2.OnchainConfig\",\"name\":\"onchainConfig\",\"type\":\"tuple\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfigTypeSafe\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"simulatePerformUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"rawReport\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6101606040523480156200001257600080fd5b506040516200564e3803806200564e83398101604081905262000035916200044b565b80816001600160a01b0316634b4fd03b6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200009b919062000472565b826001600160a01b031663ca30e6036040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000da573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200010091906200044b565b836001600160a01b031663b10b673c6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200013f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200016591906200044b565b846001600160a01b0316636709d0e56040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca91906200044b565b856001600160a01b0316635425d8ac6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000209573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200022f91906200044b565b866001600160a01b031663a08714c06040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200026e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200029491906200044b565b3380600081620002eb5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200031e576200031e8162000387565b50505085600281111562000336576200033662000495565b60e08160028111156200034d576200034d62000495565b9052506001600160a01b0394851660805292841660a05290831660c052821661010052811661012052919091166101405250620004ab9050565b336001600160a01b03821603620003e15760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620002e2565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200044857600080fd5b50565b6000602082840312156200045e57600080fd5b81516200046b8162000432565b9392505050565b6000602082840312156200048557600080fd5b8151600381106200046b57600080fd5b634e487b7160e01b600052602160045260246000fd5b60805160a05160c05160e051610100516101205161014051615136620005186000396000818160d6015261016f015260006127150152600050506000818161253d0152818161347e015281816136110152613b2e01526000505060005050600061134a01526151366000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c8063a4c0ed3611610081578063b1dc65a41161005b578063b1dc65a4146102e0578063e3d0e712146102f3578063f2fde38b14610306576100d4565b8063a4c0ed3614610262578063aed2e92914610275578063afcb95d71461029f576100d4565b806379ba5097116100b257806379ba5097146101c757806381ff7048146101cf5780638da5cb5b14610244576100d4565b8063181f5a771461011b578063349e8cca1461016d5780636cad5469146101b4575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e808015610114573d6000f35b3d6000fd5b005b6101576040518060400160405280601881526020017f4175746f6d6174696f6e526567697374727920322e322e30000000000000000081525081565b6040516101649190613dad565b60405180910390f35b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610164565b6101196101c23660046141dc565b610319565b610119611230565b61022160155460115463ffffffff780100000000000000000000000000000000000000000000000083048116937c01000000000000000000000000000000000000000000000000000000009093041691565b6040805163ffffffff948516815293909216602084015290820152606001610164565b60005473ffffffffffffffffffffffffffffffffffffffff1661018f565b6101196102703660046142f2565b611332565b61028861028336600461434e565b61154e565b604080519215158352602083019190915201610164565b601154601254604080516000815260208101939093527401000000000000000000000000000000000000000090910463ffffffff1690820152606001610164565b6101196102ee3660046143df565b6116c6565b610119610301366004614496565b61226f565b610119610314366004614525565b612298565b6103216122ac565b601f8651111561035d576040517f25d0209c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8360ff1660000361039a576040517fe77dba5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845186511415806103b957506103b1846003614571565b60ff16865111155b156103f0576040517f1d2d1c5800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601254600e546bffffffffffffffffffffffff9091169060005b816bffffffffffffffffffffffff168110156104725761045f600e82815481106104365761043661458d565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff16848461232f565b508061046a816145bc565b91505061040a565b5060008060005b836bffffffffffffffffffffffff1681101561057b57600d81815481106104a2576104a261458d565b600091825260209091200154600e805473ffffffffffffffffffffffffffffffffffffffff909216945090829081106104dd576104dd61458d565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff8681168452600c8352604080852080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001690559116808452600b90925290912080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055915080610573816145bc565b915050610479565b50610588600d6000613c82565b610594600e6000613c82565b604080516080810182526000808252602082018190529181018290526060810182905290805b8c518110156109fd57600c60008e83815181106105d9576105d961458d565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528101919091526040016000205460ff1615610644576040517f77cea0fa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168d828151811061066e5761066e61458d565b602002602001015173ffffffffffffffffffffffffffffffffffffffff16036106c3576040517f815e1d6400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405180604001604052806001151581526020018260ff16815250600c60008f84815181106106f4576106f461458d565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528181019290925260400160002082518154939092015160ff16610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909316929092171790558b518c908290811061079c5761079c61458d565b60200260200101519150600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160361080c576040517f58a70a0a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff82166000908152600b60209081526040918290208251608081018452905460ff80821615801584526101008304909116938301939093526bffffffffffffffffffffffff6201000082048116948301949094526e010000000000000000000000000000900490921660608301529093506108c7576040517f6a7281ad00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001835260ff80821660208086019182526bffffffffffffffffffffffff808b166060880190815273ffffffffffffffffffffffffffffffffffffffff87166000908152600b909352604092839020885181549551948a0151925184166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff939094166201000002929092167fffffffffffff000000000000000000000000000000000000000000000000ffff94909616610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000090951694909417179190911692909217919091179055806109f5816145bc565b9150506105ba565b50508a51610a139150600d9060208d0190613ca0565b508851610a2790600e9060208c0190613ca0565b50604051806101400160405280856bffffffffffffffffffffffff168152602001886000015163ffffffff168152602001886020015163ffffffff168152602001600063ffffffff168152602001886060015162ffffff168152602001886080015161ffff1681526020018960ff1681526020016012600001601e9054906101000a900460ff16151581526020016012600001601f9054906101000a900460ff1615158152602001886101e001511515815250601260008201518160000160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550602082015181600001600c6101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160106101000a81548163ffffffff021916908363ffffffff16021790555060608201518160000160146101000a81548163ffffffff021916908363ffffffff16021790555060808201518160000160186101000a81548162ffffff021916908362ffffff16021790555060a082015181600001601b6101000a81548161ffff021916908361ffff16021790555060c082015181600001601d6101000a81548160ff021916908360ff16021790555060e082015181600001601e6101000a81548160ff02191690831515021790555061010082015181600001601f6101000a81548160ff0219169083151502179055506101208201518160010160006101000a81548160ff0219169083151502179055509050506040518061018001604052808860a001516bffffffffffffffffffffffff16815260200188610180015173ffffffffffffffffffffffffffffffffffffffff168152602001601460010160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff168152602001886040015163ffffffff1681526020018860c0015163ffffffff168152602001601460010160149054906101000a900463ffffffff1663ffffffff168152602001601460010160189054906101000a900463ffffffff1663ffffffff1681526020016014600101601c9054906101000a900463ffffffff1663ffffffff1681526020018860e0015163ffffffff16815260200188610100015163ffffffff16815260200188610120015163ffffffff168152602001886101c0015173ffffffffffffffffffffffffffffffffffffffff16815250601460008201518160000160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550602082015181600001600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060408201518160010160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550606082015181600101600c6101000a81548163ffffffff021916908363ffffffff16021790555060808201518160010160106101000a81548163ffffffff021916908363ffffffff16021790555060a08201518160010160146101000a81548163ffffffff021916908363ffffffff16021790555060c08201518160010160186101000a81548163ffffffff021916908363ffffffff16021790555060e082015181600101601c6101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160020160006101000a81548163ffffffff021916908363ffffffff1602179055506101208201518160020160046101000a81548163ffffffff021916908363ffffffff1602179055506101408201518160020160086101000a81548163ffffffff021916908363ffffffff16021790555061016082015181600201600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505086610140015160178190555086610160015160188190555060006014600101601c9054906101000a900463ffffffff169050611017612537565b601580547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167c010000000000000000000000000000000000000000000000000000000063ffffffff938416021780825560019260189161109291859178010000000000000000000000000000000000000000000000009004166145f4565b92506101000a81548163ffffffff021916908363ffffffff1602179055506000886040516020016110c39190614662565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291905260155490915061112c90469030907801000000000000000000000000000000000000000000000000900463ffffffff168f8f8f878f8f6125ec565b60115560005b61113c6009612696565b81101561116c576111596111516009836126a6565b6009906126b9565b5080611164816145bc565b915050611132565b5060005b896101a00151518110156111c3576111b08a6101a0015182815181106111985761119861458d565b602002602001015160096126db90919063ffffffff16565b50806111bb816145bc565b915050611170565b507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0582601154601460010160189054906101000a900463ffffffff168f8f8f878f8f60405161121a999897969594939291906147dd565b60405180910390a1505050505050505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146112b6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146113a1576040517fc8bad78d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b602081146113db576040517fdfe9309000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006113e982840184614873565b60008181526004602052604090205490915065010000000000900463ffffffff90811614611443576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526004602052604090206001015461147e9085906c0100000000000000000000000090046bffffffffffffffffffffffff1661488c565b600082815260046020526040902060010180546bffffffffffffffffffffffff929092166c01000000000000000000000000027fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff9092169190911790556019546114e99085906148b1565b6019556040516bffffffffffffffffffffffff8516815273ffffffffffffffffffffffffffffffffffffffff86169082907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a35050505050565b6000806115596126fd565b6012547e01000000000000000000000000000000000000000000000000000000000000900460ff16156115b8576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600085815260046020908152604091829020825160e081018452815460ff811615158252610100810463ffffffff908116838601819052650100000000008304821684880152690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff16606084018190526001909401546bffffffffffffffffffffffff80821660808601526c0100000000000000000000000082041660a0850152780100000000000000000000000000000000000000000000000090041660c08301528451601f890185900485028101850190955287855290936116b793899089908190840183828082843760009201919091525061276c92505050565b9093509150505b935093915050565b60005a60408051610140810182526012546bffffffffffffffffffffffff8116825263ffffffff6c010000000000000000000000008204811660208401527001000000000000000000000000000000008204811693830193909352740100000000000000000000000000000000000000008104909216606082015262ffffff7801000000000000000000000000000000000000000000000000830416608082015261ffff7b0100000000000000000000000000000000000000000000000000000083041660a082015260ff7d0100000000000000000000000000000000000000000000000000000000008304811660c08301527e010000000000000000000000000000000000000000000000000000000000008304811615801560e08401527f010000000000000000000000000000000000000000000000000000000000000090930481161515610100830152601354161515610120820152919250611858576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600b602052604090205460ff166118a1576040517f1099ed7500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6011548a35146118dd576040517fdfdcf8e700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60c08101516118ed9060016148c4565b60ff16861415806118fe5750858414155b15611935576040517f0244f71a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6119458a8a8a8a8a8a8a8a612995565b60006119518a8a612bfe565b9050600081604001515167ffffffffffffffff81111561197357611973613dc0565b604051908082528060200260200182016040528015611a3757816020015b604080516101e0810182526000610100820181815261012083018290526101408301829052610160830182905261018083018290526101a083018290526101c0830182905282526020808301829052928201819052606082018190526080820181905260a0820181905260c0820181905260e082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816119915790505b5090506000805b836040015151811015611e81576004600085604001518381518110611a6557611a6561458d565b6020908102919091018101518252818101929092526040908101600020815160e081018352815460ff811615158252610100810463ffffffff90811695830195909552650100000000008104851693820193909352690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff166060830152600101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490911660c08201528351849083908110611b4a57611b4a61458d565b602002602001015160000181905250611b7f84604001518281518110611b7257611b7261458d565b6020026020010151612cb7565b838281518110611b9157611b9161458d565b6020026020010151608001906001811115611bae57611bae6148dd565b90816001811115611bc157611bc16148dd565b81525050611c3585848381518110611bdb57611bdb61458d565b60200260200101516080015186606001518481518110611bfd57611bfd61458d565b60200260200101518760a001518581518110611c1b57611c1b61458d565b602002602001015151886000015189602001516001612d62565b838281518110611c4757611c4761458d565b6020026020010151604001906bffffffffffffffffffffffff1690816bffffffffffffffffffffffff1681525050611cd484604001518281518110611c8e57611c8e61458d565b602002602001015185608001518381518110611cac57611cac61458d565b6020026020010151858481518110611cc657611cc661458d565b602002602001015188612dad565b848381518110611ce657611ce661458d565b6020026020010151602001858481518110611d0357611d0361458d565b602002602001015160e0018281525082151515158152505050828181518110611d2e57611d2e61458d565b60200260200101516020015115611d5157611d4a60018361490c565b9150611d56565b611e6f565b611dbc838281518110611d6b57611d6b61458d565b6020026020010151600001516060015185606001518381518110611d9157611d9161458d565b60200260200101518660a001518481518110611daf57611daf61458d565b602002602001015161276c565b848381518110611dce57611dce61458d565b6020026020010151606001858481518110611deb57611deb61458d565b602002602001015160a0018281525082151515158152505050828181518110611e1657611e1661458d565b602002602001015160a0015186611e2d9190614927565b9550611e6f84604001518281518110611e4857611e4861458d565b6020026020010151848381518110611e6257611e6261458d565b6020026020010151612f30565b80611e79816145bc565b915050611a3e565b508061ffff16600003611e98575050505050612265565b60c0840151611ea89060016148c4565b611eb79060ff1661044c61493a565b616b6c611ec58d601061493a565b5a611ed09089614927565b611eda91906148b1565b611ee491906148b1565b611eee91906148b1565b9450611b58611f0161ffff831687614980565b611f0b91906148b1565b945060008060008060005b87604001515181101561210c57868181518110611f3557611f3561458d565b602002602001015160200151156120fa57611f918a888381518110611f5c57611f5c61458d565b6020026020010151608001518a60a001518481518110611f7e57611f7e61458d565b6020026020010151518c60c00151613042565b878281518110611fa357611fa361458d565b602002602001015160c0018181525050611fff8989604001518381518110611fcd57611fcd61458d565b6020026020010151898481518110611fe757611fe761458d565b60200260200101518b600001518c602001518b613062565b909350915061200e828561488c565b935061201a838661488c565b945086818151811061202e5761202e61458d565b6020026020010151606001511515886040015182815181106120525761205261458d565b60200260200101517fad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b8486612087919061488c565b8a85815181106120995761209961458d565b602002602001015160a001518b86815181106120b7576120b761458d565b602002602001015160c001518d6080015187815181106120d9576120d961458d565b60200260200101516040516120f19493929190614994565b60405180910390a35b80612104816145bc565b915050611f16565b5050336000908152600b6020526040902080548492506002906121449084906201000090046bffffffffffffffffffffffff1661488c565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555080601260000160008282829054906101000a90046bffffffffffffffffffffffff1661219e919061488c565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060008f6001600381106121e1576121e161458d565b602002013560001c9050600060088264ffffffffff16901c9050876060015163ffffffff168163ffffffff16111561225b57601280547fffffffffffffffff00000000ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000063ffffffff8416021790555b5050505050505050505b5050505050505050565b612290868686868060200190518101906122899190614a77565b8686610319565b505050505050565b6122a06122ac565b6122a981613155565b50565b60005473ffffffffffffffffffffffffffffffffffffffff16331461232d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016112ad565b565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e010000000000000000000000000000900490911660608201529061252b5760008160600151856123c79190614be5565b905060006123d58583614c0a565b905080836040018181516123e9919061488c565b6bffffffffffffffffffffffff169052506124048582614c35565b83606001818151612415919061488c565b6bffffffffffffffffffffffff90811690915273ffffffffffffffffffffffffffffffffffffffff89166000908152600b602090815260409182902087518154928901519389015160608a015186166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff919096166201000002167fffffffffffff000000000000000000000000000000000000000000000000ffff60ff95909516610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909416939093171792909216179190911790555050505b60400151949350505050565b600060017f0000000000000000000000000000000000000000000000000000000000000000600281111561256d5761256d6148dd565b036125e757606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156125be573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906125e29190614c65565b905090565b504390565b6000808a8a8a8a8a8a8a8a8a60405160200161261099989796959493929190614c7e565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179b9a5050505050505050505050565b60006126a0825490565b92915050565b60006126b2838361324a565b9392505050565b60006126b28373ffffffffffffffffffffffffffffffffffffffff8416613274565b60006126b28373ffffffffffffffffffffffffffffffffffffffff841661336e565b3273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461232d576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60125460009081907f0100000000000000000000000000000000000000000000000000000000000000900460ff16156127d1576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601280547effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f01000000000000000000000000000000000000000000000000000000000000001790556040517f4585e33b000000000000000000000000000000000000000000000000000000009061284d908590602401613dad565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009094169390931790925290517f79188d1600000000000000000000000000000000000000000000000000000000815290935073ffffffffffffffffffffffffffffffffffffffff8616906379188d16906129209087908790600401614d13565b60408051808303816000875af115801561293e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906129629190614d2c565b601280547effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff16905590969095509350505050565b600087876040516129a7929190614d5a565b6040519081900381206129be918b90602001614d6a565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201208383019092526000808452908301819052909250906000805b88811015612b9557600185878360208110612a2a57612a2a61458d565b612a3791901a601b6148c4565b8c8c85818110612a4957612a4961458d565b905060200201358b8b86818110612a6257612a6261458d565b9050602002013560405160008152602001604052604051612a9f949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015612ac1573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff81166000908152600c602090815290849020838501909452925460ff8082161515808552610100909204169383019390935290955093509050612b6f576040517f0f4c073700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826020015160080260ff166001901b840193508080612b8d906145bc565b915050612a0d565b50827e01010101010101010101010101010101010101010101010101010101010101841614612bf0576040517fc103be2e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505050505050505050505050565b612c376040518060c001604052806000815260200160008152602001606081526020016060815260200160608152602001606081525090565b6000612c4583850185614e5b565b6040810151516060820151519192509081141580612c6857508082608001515114155b80612c785750808260a001515114155b15612caf576040517fb55ac75400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b509392505050565b6000818160045b600f811015612d44577fff000000000000000000000000000000000000000000000000000000000000008216838260208110612cfc57612cfc61458d565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614612d3257506000949350505050565b80612d3c816145bc565b915050612cbe565b5081600f1a6001811115612d5a57612d5a6148dd565b949350505050565b600080612d7488878b60c001516133bd565b9050600080612d8f8b8a63ffffffff16858a8a60018b613449565b9092509050612d9e818361488c565b9b9a5050505050505050505050565b600080808085608001516001811115612dc857612dc86148dd565b03612ded57612dd9878787876138a2565b612de857600092509050612f27565b612e64565b600185608001516001811115612e0557612e056148dd565b03612e32576000612e178888876139a4565b9250905080612e2c5750600092509050612f27565b50612e64565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612e6c612537565b85516040015163ffffffff1611612ec057867fc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd563687604051612ead9190613dad565b60405180910390a2600092509050612f27565b84604001516bffffffffffffffffffffffff16856000015160a001516bffffffffffffffffffffffff161015612f2057867f377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e0287604051612ead9190613dad565b6001925090505b94509492505050565b600081608001516001811115612f4857612f486148dd565b03612fba57612f55612537565b6000838152600460205260409020600101805463ffffffff929092167801000000000000000000000000000000000000000000000000027fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff9092169190911790555050565b600181608001516001811115612fd257612fd26148dd565b0361303e5760e08101805160009081526008602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055915191517fa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f29190a25b5050565b600061304f8484846133bd565b905080851015612d5a5750929392505050565b60008061307d888760a001518860c001518888886001613449565b9092509050600061308e828461488c565b600089815260046020526040902060010180549192508291600c906130d29084906c0100000000000000000000000090046bffffffffffffffffffffffff16614be5565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915560008a81526004602052604081206001018054859450909261311b9185911661488c565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555050965096945050505050565b3373ffffffffffffffffffffffffffffffffffffffff8216036131d4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016112ad565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008260000182815481106132615761326161458d565b9060005260206000200154905092915050565b6000818152600183016020526040812054801561335d576000613298600183614927565b85549091506000906132ac90600190614927565b90508181146133115760008660000182815481106132cc576132cc61458d565b90600052602060002001549050808760000184815481106132ef576132ef61458d565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061332257613322614f48565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506126a0565b60009150506126a0565b5092915050565b60008181526001830160205260408120546133b5575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556126a0565b5060006126a0565b600080808560018111156133d3576133d36148dd565b036133e2575062015f90613401565b60018560018111156133f6576133f66148dd565b03612e3257506201adb05b61341263ffffffff8516601461493a565b61341d8460016148c4565b61342c9060ff16611d4c61493a565b61343690836148b1565b61344091906148b1565b95945050505050565b60008060008960a0015161ffff1687613462919061493a565b90508380156134705750803a105b1561347857503a5b600060027f000000000000000000000000000000000000000000000000000000000000000060028111156134ae576134ae6148dd565b0361360d57604080516000815260208101909152851561350c576000366040518060800160405280604881526020016150e2604891396040516020016134f693929190614f77565b6040516020818303038152906040529050613574565b60165461352890640100000000900463ffffffff166004614f9e565b63ffffffff1667ffffffffffffffff81111561354657613546613dc0565b6040519080825280601f01601f191660200182016040528015613570576020820181803683370190505b5090505b6040517f49948e0e00000000000000000000000000000000000000000000000000000000815273420000000000000000000000000000000000000f906349948e0e906135c4908490600401613dad565b602060405180830381865afa1580156135e1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906136059190614c65565b915050613767565b60017f00000000000000000000000000000000000000000000000000000000000000006002811115613641576136416148dd565b036137675784156136c357606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613698573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906136bc9190614c65565b9050613767565b6000606c73ffffffffffffffffffffffffffffffffffffffff166341b247a86040518163ffffffff1660e01b815260040160c060405180830381865afa158015613711573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906137359190614fbe565b505060165492945061375893505050640100000000900463ffffffff168261493a565b61376390601061493a565b9150505b8461378357808b60a0015161ffff16613780919061493a565b90505b61379161ffff871682614980565b9050600087826137a18c8e6148b1565b6137ab908661493a565b6137b591906148b1565b6137c790670de0b6b3a764000061493a565b6137d19190614980565b905060008c6040015163ffffffff1664e8d4a510006137f0919061493a565b898e6020015163ffffffff16858f88613809919061493a565b61381391906148b1565b61382190633b9aca0061493a565b61382b919061493a565b6138359190614980565b61383f91906148b1565b90506b033b2e3c9fd0803ce800000061385882846148b1565b1115613890576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b600080848060200190518101906138b99190615008565b845160c00151815191925063ffffffff9081169116101561391657857f405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8866040516139049190613dad565b60405180910390a26000915050612d5a565b826101200151801561394a575060208101511580159061394a5750602081015181516139479063ffffffff16613b28565b14155b806139635750613958612537565b815163ffffffff1610155b1561399857857f6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301866040516139049190613dad565b50600195945050505050565b6000806000848060200190518101906139bd9190615060565b9050600086826000015183602001518460400151604051602001613a1f94939291909384526020840192909252604083015260e01b7fffffffff0000000000000000000000000000000000000000000000000000000016606082015260640190565b6040516020818303038152906040528051906020012090508461012001518015613a6d5750608082015115801590613a6d57508160800151613a6a836060015163ffffffff16613b28565b14155b80613a895750613a7b612537565b826060015163ffffffff1610155b15613ad357867f6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc30187604051613abe9190613dad565b60405180910390a26000935091506116be9050565b60008181526008602052604090205460ff1615613b1a57867f405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e887604051613abe9190613dad565b600197909650945050505050565b600060017f00000000000000000000000000000000000000000000000000000000000000006002811115613b5e57613b5e6148dd565b03613c78576000606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613bb1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613bd59190614c65565b90508083101580613bf05750610100613bee8483614927565b115b15613bfe5750600092915050565b6040517f2b407a8200000000000000000000000000000000000000000000000000000000815260048101849052606490632b407a8290602401602060405180830381865afa158015613c54573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906126b29190614c65565b504090565b919050565b50805460008255906000526020600020908101906122a99190613d2a565b828054828255906000526020600020908101928215613d1a579160200282015b82811115613d1a57825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190613cc0565b50613d26929150613d2a565b5090565b5b80821115613d265760008155600101613d2b565b60005b83811015613d5a578181015183820152602001613d42565b50506000910152565b60008151808452613d7b816020860160208601613d3f565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006126b26020830184613d63565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610200810167ffffffffffffffff81118282101715613e1357613e13613dc0565b60405290565b60405160c0810167ffffffffffffffff81118282101715613e1357613e13613dc0565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613e8357613e83613dc0565b604052919050565b600067ffffffffffffffff821115613ea557613ea5613dc0565b5060051b60200190565b73ffffffffffffffffffffffffffffffffffffffff811681146122a957600080fd5b8035613c7d81613eaf565b600082601f830112613eed57600080fd5b81356020613f02613efd83613e8b565b613e3c565b82815260059290921b84018101918181019086841115613f2157600080fd5b8286015b84811015613f45578035613f3881613eaf565b8352918301918301613f25565b509695505050505050565b803560ff81168114613c7d57600080fd5b63ffffffff811681146122a957600080fd5b8035613c7d81613f61565b62ffffff811681146122a957600080fd5b8035613c7d81613f7e565b61ffff811681146122a957600080fd5b8035613c7d81613f9a565b6bffffffffffffffffffffffff811681146122a957600080fd5b8035613c7d81613fb5565b80151581146122a957600080fd5b8035613c7d81613fda565b6000610200828403121561400657600080fd5b61400e613def565b905061401982613f73565b815261402760208301613f73565b602082015261403860408301613f73565b604082015261404960608301613f8f565b606082015261405a60808301613faa565b608082015261406b60a08301613fcf565b60a082015261407c60c08301613f73565b60c082015261408d60e08301613f73565b60e08201526101006140a0818401613f73565b908201526101206140b2838201613f73565b90820152610140828101359082015261016080830135908201526101806140da818401613ed1565b908201526101a08281013567ffffffffffffffff8111156140fa57600080fd5b61410685828601613edc565b8284015250506101c061411a818401613ed1565b908201526101e061412c838201613fe8565b9082015292915050565b803567ffffffffffffffff81168114613c7d57600080fd5b600082601f83011261415f57600080fd5b813567ffffffffffffffff81111561417957614179613dc0565b6141aa60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601613e3c565b8181528460208386010111156141bf57600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060c087890312156141f557600080fd5b863567ffffffffffffffff8082111561420d57600080fd5b6142198a838b01613edc565b9750602089013591508082111561422f57600080fd5b61423b8a838b01613edc565b965061424960408a01613f50565b9550606089013591508082111561425f57600080fd5b61426b8a838b01613ff3565b945061427960808a01614136565b935060a089013591508082111561428f57600080fd5b5061429c89828a0161414e565b9150509295509295509295565b60008083601f8401126142bb57600080fd5b50813567ffffffffffffffff8111156142d357600080fd5b6020830191508360208285010111156142eb57600080fd5b9250929050565b6000806000806060858703121561430857600080fd5b843561431381613eaf565b935060208501359250604085013567ffffffffffffffff81111561433657600080fd5b614342878288016142a9565b95989497509550505050565b60008060006040848603121561436357600080fd5b83359250602084013567ffffffffffffffff81111561438157600080fd5b61438d868287016142a9565b9497909650939450505050565b60008083601f8401126143ac57600080fd5b50813567ffffffffffffffff8111156143c457600080fd5b6020830191508360208260051b85010111156142eb57600080fd5b60008060008060008060008060e0898b0312156143fb57600080fd5b606089018a81111561440c57600080fd5b8998503567ffffffffffffffff8082111561442657600080fd5b6144328c838d016142a9565b909950975060808b013591508082111561444b57600080fd5b6144578c838d0161439a565b909750955060a08b013591508082111561447057600080fd5b5061447d8b828c0161439a565b999c989b50969995989497949560c00135949350505050565b60008060008060008060c087890312156144af57600080fd5b863567ffffffffffffffff808211156144c757600080fd5b6144d38a838b01613edc565b975060208901359150808211156144e957600080fd5b6144f58a838b01613edc565b965061450360408a01613f50565b9550606089013591508082111561451957600080fd5b61426b8a838b0161414e565b60006020828403121561453757600080fd5b81356126b281613eaf565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60ff818116838216029081169081811461336757613367614542565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036145ed576145ed614542565b5060010190565b63ffffffff81811683821601908082111561336757613367614542565b600081518084526020808501945080840160005b8381101561465757815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101614625565b509495945050505050565b6020815261467960208201835163ffffffff169052565b60006020830151614692604084018263ffffffff169052565b50604083015163ffffffff8116606084015250606083015162ffffff8116608084015250608083015161ffff811660a08401525060a08301516bffffffffffffffffffffffff811660c08401525060c083015163ffffffff811660e08401525060e083015161010061470b8185018363ffffffff169052565b84015190506101206147248482018363ffffffff169052565b840151905061014061473d8482018363ffffffff169052565b840151610160848101919091528401516101808085019190915284015190506101a06147808185018373ffffffffffffffffffffffffffffffffffffffff169052565b808501519150506102006101c081818601526147a0610220860184614611565b908601519092506101e06147cb8682018373ffffffffffffffffffffffffffffffffffffffff169052565b90950151151593019290925250919050565b600061012063ffffffff808d1684528b6020850152808b1660408501525080606084015261480d8184018a614611565b905082810360808401526148218189614611565b905060ff871660a084015282810360c084015261483e8187613d63565b905067ffffffffffffffff851660e08401528281036101008401526148638185613d63565b9c9b505050505050505050505050565b60006020828403121561488557600080fd5b5035919050565b6bffffffffffffffffffffffff81811683821601908082111561336757613367614542565b808201808211156126a0576126a0614542565b60ff81811683821601908111156126a0576126a0614542565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b61ffff81811683821601908082111561336757613367614542565b818103818111156126a0576126a0614542565b80820281158282048414176126a0576126a0614542565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60008261498f5761498f614951565b500490565b6bffffffffffffffffffffffff851681528360208201528260408201526080606082015260006149c76080830184613d63565b9695505050505050565b8051613c7d81613f61565b8051613c7d81613f7e565b8051613c7d81613f9a565b8051613c7d81613fb5565b8051613c7d81613eaf565b600082601f830112614a1957600080fd5b81516020614a29613efd83613e8b565b82815260059290921b84018101918181019086841115614a4857600080fd5b8286015b84811015613f45578051614a5f81613eaf565b8352918301918301614a4c565b8051613c7d81613fda565b600060208284031215614a8957600080fd5b815167ffffffffffffffff80821115614aa157600080fd5b908301906102008286031215614ab657600080fd5b614abe613def565b614ac7836149d1565b8152614ad5602084016149d1565b6020820152614ae6604084016149d1565b6040820152614af7606084016149dc565b6060820152614b08608084016149e7565b6080820152614b1960a084016149f2565b60a0820152614b2a60c084016149d1565b60c0820152614b3b60e084016149d1565b60e0820152610100614b4e8185016149d1565b90820152610120614b608482016149d1565b9082015261014083810151908201526101608084015190820152610180614b888185016149fd565b908201526101a08381015183811115614ba057600080fd5b614bac88828701614a08565b8284015250506101c09150614bc28284016149fd565b828201526101e09150614bd6828401614a6c565b91810191909152949350505050565b6bffffffffffffffffffffffff82811682821603908082111561336757613367614542565b60006bffffffffffffffffffffffff80841680614c2957614c29614951565b92169190910492915050565b6bffffffffffffffffffffffff818116838216028082169190828114614c5d57614c5d614542565b505092915050565b600060208284031215614c7757600080fd5b5051919050565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b166040850152816060850152614cc58285018b614611565b91508382036080850152614cd9828a614611565b915060ff881660a085015283820360c0850152614cf68288613d63565b90861660e085015283810361010085015290506148638185613d63565b828152604060208201526000612d5a6040830184613d63565b60008060408385031215614d3f57600080fd5b8251614d4a81613fda565b6020939093015192949293505050565b8183823760009101908152919050565b8281526080810160608360208401379392505050565b600082601f830112614d9157600080fd5b81356020614da1613efd83613e8b565b82815260059290921b84018101918181019086841115614dc057600080fd5b8286015b84811015613f455780358352918301918301614dc4565b600082601f830112614dec57600080fd5b81356020614dfc613efd83613e8b565b82815260059290921b84018101918181019086841115614e1b57600080fd5b8286015b84811015613f4557803567ffffffffffffffff811115614e3f5760008081fd5b614e4d8986838b010161414e565b845250918301918301614e1f565b600060208284031215614e6d57600080fd5b813567ffffffffffffffff80821115614e8557600080fd5b9083019060c08286031215614e9957600080fd5b614ea1613e19565b8235815260208301356020820152604083013582811115614ec157600080fd5b614ecd87828601614d80565b604083015250606083013582811115614ee557600080fd5b614ef187828601614d80565b606083015250608083013582811115614f0957600080fd5b614f1587828601614ddb565b60808301525060a083013582811115614f2d57600080fd5b614f3987828601614ddb565b60a08301525095945050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b828482376000838201600081528351614f94818360208801613d3f565b0195945050505050565b63ffffffff818116838216028082169190828114614c5d57614c5d614542565b60008060008060008060c08789031215614fd757600080fd5b865195506020870151945060408701519350606087015192506080870151915060a087015190509295509295509295565b60006040828403121561501a57600080fd5b6040516040810181811067ffffffffffffffff8211171561503d5761503d613dc0565b604052825161504b81613f61565b81526020928301519281019290925250919050565b600060a0828403121561507257600080fd5b60405160a0810181811067ffffffffffffffff8211171561509557615095613dc0565b8060405250825181526020830151602082015260408301516150b681613f61565b604082015260608301516150c981613f61565b6060820152608092830151928101929092525091905056fe307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000813000a",
}

var AutomationRegistryABI = AutomationRegistryMetaData.ABI

var AutomationRegistryBin = AutomationRegistryMetaData.Bin

func DeployAutomationRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, logicA common.Address) (common.Address, *types.Transaction, *AutomationRegistry, error) {
	parsed, err := AutomationRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AutomationRegistryBin), backend, logicA)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AutomationRegistry{address: address, abi: *parsed, AutomationRegistryCaller: AutomationRegistryCaller{contract: contract}, AutomationRegistryTransactor: AutomationRegistryTransactor{contract: contract}, AutomationRegistryFilterer: AutomationRegistryFilterer{contract: contract}}, nil
}

type AutomationRegistry struct {
	address common.Address
	abi     abi.ABI
	AutomationRegistryCaller
	AutomationRegistryTransactor
	AutomationRegistryFilterer
}

type AutomationRegistryCaller struct {
	contract *bind.BoundContract
}

type AutomationRegistryTransactor struct {
	contract *bind.BoundContract
}

type AutomationRegistryFilterer struct {
	contract *bind.BoundContract
}

type AutomationRegistrySession struct {
	Contract     *AutomationRegistry
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type AutomationRegistryCallerSession struct {
	Contract *AutomationRegistryCaller
	CallOpts bind.CallOpts
}

type AutomationRegistryTransactorSession struct {
	Contract     *AutomationRegistryTransactor
	TransactOpts bind.TransactOpts
}

type AutomationRegistryRaw struct {
	Contract *AutomationRegistry
}

type AutomationRegistryCallerRaw struct {
	Contract *AutomationRegistryCaller
}

type AutomationRegistryTransactorRaw struct {
	Contract *AutomationRegistryTransactor
}

func NewAutomationRegistry(address common.Address, backend bind.ContractBackend) (*AutomationRegistry, error) {
	abi, err := abi.JSON(strings.NewReader(AutomationRegistryABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindAutomationRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistry{address: address, abi: abi, AutomationRegistryCaller: AutomationRegistryCaller{contract: contract}, AutomationRegistryTransactor: AutomationRegistryTransactor{contract: contract}, AutomationRegistryFilterer: AutomationRegistryFilterer{contract: contract}}, nil
}

func NewAutomationRegistryCaller(address common.Address, caller bind.ContractCaller) (*AutomationRegistryCaller, error) {
	contract, err := bindAutomationRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryCaller{contract: contract}, nil
}

func NewAutomationRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*AutomationRegistryTransactor, error) {
	contract, err := bindAutomationRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryTransactor{contract: contract}, nil
}

func NewAutomationRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*AutomationRegistryFilterer, error) {
	contract, err := bindAutomationRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryFilterer{contract: contract}, nil
}

func bindAutomationRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AutomationRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_AutomationRegistry *AutomationRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationRegistry.Contract.AutomationRegistryCaller.contract.Call(opts, result, method, params...)
}

func (_AutomationRegistry *AutomationRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.AutomationRegistryTransactor.contract.Transfer(opts)
}

func (_AutomationRegistry *AutomationRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.AutomationRegistryTransactor.contract.Transact(opts, method, params...)
}

func (_AutomationRegistry *AutomationRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationRegistry.Contract.contract.Call(opts, result, method, params...)
}

func (_AutomationRegistry *AutomationRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.contract.Transfer(opts)
}

func (_AutomationRegistry *AutomationRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.contract.Transact(opts, method, params...)
}

func (_AutomationRegistry *AutomationRegistryCaller) FallbackTo(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistry.contract.Call(opts, &out, "fallbackTo")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistry *AutomationRegistrySession) FallbackTo() (common.Address, error) {
	return _AutomationRegistry.Contract.FallbackTo(&_AutomationRegistry.CallOpts)
}

func (_AutomationRegistry *AutomationRegistryCallerSession) FallbackTo() (common.Address, error) {
	return _AutomationRegistry.Contract.FallbackTo(&_AutomationRegistry.CallOpts)
}

func (_AutomationRegistry *AutomationRegistryCaller) LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _AutomationRegistry.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_AutomationRegistry *AutomationRegistrySession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _AutomationRegistry.Contract.LatestConfigDetails(&_AutomationRegistry.CallOpts)
}

func (_AutomationRegistry *AutomationRegistryCallerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _AutomationRegistry.Contract.LatestConfigDetails(&_AutomationRegistry.CallOpts)
}

func (_AutomationRegistry *AutomationRegistryCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

	error) {
	var out []interface{}
	err := _AutomationRegistry.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(LatestConfigDigestAndEpoch)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_AutomationRegistry *AutomationRegistrySession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _AutomationRegistry.Contract.LatestConfigDigestAndEpoch(&_AutomationRegistry.CallOpts)
}

func (_AutomationRegistry *AutomationRegistryCallerSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _AutomationRegistry.Contract.LatestConfigDigestAndEpoch(&_AutomationRegistry.CallOpts)
}

func (_AutomationRegistry *AutomationRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationRegistry *AutomationRegistrySession) Owner() (common.Address, error) {
	return _AutomationRegistry.Contract.Owner(&_AutomationRegistry.CallOpts)
}

func (_AutomationRegistry *AutomationRegistryCallerSession) Owner() (common.Address, error) {
	return _AutomationRegistry.Contract.Owner(&_AutomationRegistry.CallOpts)
}

func (_AutomationRegistry *AutomationRegistryCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AutomationRegistry.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_AutomationRegistry *AutomationRegistrySession) TypeAndVersion() (string, error) {
	return _AutomationRegistry.Contract.TypeAndVersion(&_AutomationRegistry.CallOpts)
}

func (_AutomationRegistry *AutomationRegistryCallerSession) TypeAndVersion() (string, error) {
	return _AutomationRegistry.Contract.TypeAndVersion(&_AutomationRegistry.CallOpts)
}

func (_AutomationRegistry *AutomationRegistryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationRegistry.contract.Transact(opts, "acceptOwnership")
}

func (_AutomationRegistry *AutomationRegistrySession) AcceptOwnership() (*types.Transaction, error) {
	return _AutomationRegistry.Contract.AcceptOwnership(&_AutomationRegistry.TransactOpts)
}

func (_AutomationRegistry *AutomationRegistryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _AutomationRegistry.Contract.AcceptOwnership(&_AutomationRegistry.TransactOpts)
}

func (_AutomationRegistry *AutomationRegistryTransactor) OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _AutomationRegistry.contract.Transact(opts, "onTokenTransfer", sender, amount, data)
}

func (_AutomationRegistry *AutomationRegistrySession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.OnTokenTransfer(&_AutomationRegistry.TransactOpts, sender, amount, data)
}

func (_AutomationRegistry *AutomationRegistryTransactorSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.OnTokenTransfer(&_AutomationRegistry.TransactOpts, sender, amount, data)
}

func (_AutomationRegistry *AutomationRegistryTransactor) SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistry.contract.Transact(opts, "setConfig", signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_AutomationRegistry *AutomationRegistrySession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.SetConfig(&_AutomationRegistry.TransactOpts, signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_AutomationRegistry *AutomationRegistryTransactorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.SetConfig(&_AutomationRegistry.TransactOpts, signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_AutomationRegistry *AutomationRegistryTransactor) SetConfigTypeSafe(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig AutomationRegistryBase22OnchainConfig, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistry.contract.Transact(opts, "setConfigTypeSafe", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_AutomationRegistry *AutomationRegistrySession) SetConfigTypeSafe(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig AutomationRegistryBase22OnchainConfig, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.SetConfigTypeSafe(&_AutomationRegistry.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_AutomationRegistry *AutomationRegistryTransactorSession) SetConfigTypeSafe(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig AutomationRegistryBase22OnchainConfig, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.SetConfigTypeSafe(&_AutomationRegistry.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_AutomationRegistry *AutomationRegistryTransactor) SimulatePerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error) {
	return _AutomationRegistry.contract.Transact(opts, "simulatePerformUpkeep", id, performData)
}

func (_AutomationRegistry *AutomationRegistrySession) SimulatePerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.SimulatePerformUpkeep(&_AutomationRegistry.TransactOpts, id, performData)
}

func (_AutomationRegistry *AutomationRegistryTransactorSession) SimulatePerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.SimulatePerformUpkeep(&_AutomationRegistry.TransactOpts, id, performData)
}

func (_AutomationRegistry *AutomationRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _AutomationRegistry.contract.Transact(opts, "transferOwnership", to)
}

func (_AutomationRegistry *AutomationRegistrySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.TransferOwnership(&_AutomationRegistry.TransactOpts, to)
}

func (_AutomationRegistry *AutomationRegistryTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.TransferOwnership(&_AutomationRegistry.TransactOpts, to)
}

func (_AutomationRegistry *AutomationRegistryTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _AutomationRegistry.contract.Transact(opts, "transmit", reportContext, rawReport, rs, ss, rawVs)
}

func (_AutomationRegistry *AutomationRegistrySession) Transmit(reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.Transmit(&_AutomationRegistry.TransactOpts, reportContext, rawReport, rs, ss, rawVs)
}

func (_AutomationRegistry *AutomationRegistryTransactorSession) Transmit(reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.Transmit(&_AutomationRegistry.TransactOpts, reportContext, rawReport, rs, ss, rawVs)
}

func (_AutomationRegistry *AutomationRegistryTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _AutomationRegistry.contract.RawTransact(opts, calldata)
}

func (_AutomationRegistry *AutomationRegistrySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.Fallback(&_AutomationRegistry.TransactOpts, calldata)
}

func (_AutomationRegistry *AutomationRegistryTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _AutomationRegistry.Contract.Fallback(&_AutomationRegistry.TransactOpts, calldata)
}

type AutomationRegistryAdminPrivilegeConfigSetIterator struct {
	Event *AutomationRegistryAdminPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryAdminPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryAdminPrivilegeConfigSet)
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
		it.Event = new(AutomationRegistryAdminPrivilegeConfigSet)
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

func (it *AutomationRegistryAdminPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryAdminPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryAdminPrivilegeConfigSet struct {
	Admin           common.Address
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*AutomationRegistryAdminPrivilegeConfigSetIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryAdminPrivilegeConfigSetIterator{contract: _AutomationRegistry.contract, event: "AdminPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryAdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryAdminPrivilegeConfigSet)
				if err := _AutomationRegistry.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseAdminPrivilegeConfigSet(log types.Log) (*AutomationRegistryAdminPrivilegeConfigSet, error) {
	event := new(AutomationRegistryAdminPrivilegeConfigSet)
	if err := _AutomationRegistry.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryCancelledUpkeepReportIterator struct {
	Event *AutomationRegistryCancelledUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryCancelledUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryCancelledUpkeepReport)
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
		it.Event = new(AutomationRegistryCancelledUpkeepReport)
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

func (it *AutomationRegistryCancelledUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryCancelledUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryCancelledUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryCancelledUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryCancelledUpkeepReportIterator{contract: _AutomationRegistry.contract, event: "CancelledUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryCancelledUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryCancelledUpkeepReport)
				if err := _AutomationRegistry.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseCancelledUpkeepReport(log types.Log) (*AutomationRegistryCancelledUpkeepReport, error) {
	event := new(AutomationRegistryCancelledUpkeepReport)
	if err := _AutomationRegistry.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryConfigSetIterator struct {
	Event *AutomationRegistryConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryConfigSet)
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
		it.Event = new(AutomationRegistryConfigSet)
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

func (it *AutomationRegistryConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	Raw                       types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterConfigSet(opts *bind.FilterOpts) (*AutomationRegistryConfigSetIterator, error) {

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryConfigSetIterator{contract: _AutomationRegistry.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryConfigSet) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryConfigSet)
				if err := _AutomationRegistry.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseConfigSet(log types.Log) (*AutomationRegistryConfigSet, error) {
	event := new(AutomationRegistryConfigSet)
	if err := _AutomationRegistry.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryDedupKeyAddedIterator struct {
	Event *AutomationRegistryDedupKeyAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryDedupKeyAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryDedupKeyAdded)
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
		it.Event = new(AutomationRegistryDedupKeyAdded)
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

func (it *AutomationRegistryDedupKeyAddedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryDedupKeyAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryDedupKeyAdded struct {
	DedupKey [32]byte
	Raw      types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*AutomationRegistryDedupKeyAddedIterator, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryDedupKeyAddedIterator{contract: _AutomationRegistry.contract, event: "DedupKeyAdded", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryDedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryDedupKeyAdded)
				if err := _AutomationRegistry.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseDedupKeyAdded(log types.Log) (*AutomationRegistryDedupKeyAdded, error) {
	event := new(AutomationRegistryDedupKeyAdded)
	if err := _AutomationRegistry.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryFundsAddedIterator struct {
	Event *AutomationRegistryFundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryFundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryFundsAdded)
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
		it.Event = new(AutomationRegistryFundsAdded)
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

func (it *AutomationRegistryFundsAddedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryFundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryFundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*AutomationRegistryFundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryFundsAddedIterator{contract: _AutomationRegistry.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryFundsAdded)
				if err := _AutomationRegistry.contract.UnpackLog(event, "FundsAdded", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseFundsAdded(log types.Log) (*AutomationRegistryFundsAdded, error) {
	event := new(AutomationRegistryFundsAdded)
	if err := _AutomationRegistry.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryFundsWithdrawnIterator struct {
	Event *AutomationRegistryFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryFundsWithdrawn)
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
		it.Event = new(AutomationRegistryFundsWithdrawn)
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

func (it *AutomationRegistryFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryFundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryFundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryFundsWithdrawnIterator{contract: _AutomationRegistry.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryFundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryFundsWithdrawn)
				if err := _AutomationRegistry.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseFundsWithdrawn(log types.Log) (*AutomationRegistryFundsWithdrawn, error) {
	event := new(AutomationRegistryFundsWithdrawn)
	if err := _AutomationRegistry.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryInsufficientFundsUpkeepReportIterator struct {
	Event *AutomationRegistryInsufficientFundsUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryInsufficientFundsUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryInsufficientFundsUpkeepReport)
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
		it.Event = new(AutomationRegistryInsufficientFundsUpkeepReport)
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

func (it *AutomationRegistryInsufficientFundsUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryInsufficientFundsUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryInsufficientFundsUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryInsufficientFundsUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryInsufficientFundsUpkeepReportIterator{contract: _AutomationRegistry.contract, event: "InsufficientFundsUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryInsufficientFundsUpkeepReport)
				if err := _AutomationRegistry.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseInsufficientFundsUpkeepReport(log types.Log) (*AutomationRegistryInsufficientFundsUpkeepReport, error) {
	event := new(AutomationRegistryInsufficientFundsUpkeepReport)
	if err := _AutomationRegistry.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryOwnerFundsWithdrawnIterator struct {
	Event *AutomationRegistryOwnerFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryOwnerFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryOwnerFundsWithdrawn)
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
		it.Event = new(AutomationRegistryOwnerFundsWithdrawn)
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

func (it *AutomationRegistryOwnerFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryOwnerFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryOwnerFundsWithdrawn struct {
	Amount *big.Int
	Raw    types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*AutomationRegistryOwnerFundsWithdrawnIterator, error) {

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryOwnerFundsWithdrawnIterator{contract: _AutomationRegistry.contract, event: "OwnerFundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryOwnerFundsWithdrawn) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryOwnerFundsWithdrawn)
				if err := _AutomationRegistry.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseOwnerFundsWithdrawn(log types.Log) (*AutomationRegistryOwnerFundsWithdrawn, error) {
	event := new(AutomationRegistryOwnerFundsWithdrawn)
	if err := _AutomationRegistry.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryOwnershipTransferRequestedIterator struct {
	Event *AutomationRegistryOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryOwnershipTransferRequested)
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
		it.Event = new(AutomationRegistryOwnershipTransferRequested)
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

func (it *AutomationRegistryOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryOwnershipTransferRequestedIterator{contract: _AutomationRegistry.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryOwnershipTransferRequested)
				if err := _AutomationRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseOwnershipTransferRequested(log types.Log) (*AutomationRegistryOwnershipTransferRequested, error) {
	event := new(AutomationRegistryOwnershipTransferRequested)
	if err := _AutomationRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryOwnershipTransferredIterator struct {
	Event *AutomationRegistryOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryOwnershipTransferred)
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
		it.Event = new(AutomationRegistryOwnershipTransferred)
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

func (it *AutomationRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryOwnershipTransferredIterator{contract: _AutomationRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryOwnershipTransferred)
				if err := _AutomationRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*AutomationRegistryOwnershipTransferred, error) {
	event := new(AutomationRegistryOwnershipTransferred)
	if err := _AutomationRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryPausedIterator struct {
	Event *AutomationRegistryPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryPaused)
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
		it.Event = new(AutomationRegistryPaused)
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

func (it *AutomationRegistryPausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterPaused(opts *bind.FilterOpts) (*AutomationRegistryPausedIterator, error) {

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryPausedIterator{contract: _AutomationRegistry.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryPaused) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryPaused)
				if err := _AutomationRegistry.contract.UnpackLog(event, "Paused", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParsePaused(log types.Log) (*AutomationRegistryPaused, error) {
	event := new(AutomationRegistryPaused)
	if err := _AutomationRegistry.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryPayeesUpdatedIterator struct {
	Event *AutomationRegistryPayeesUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryPayeesUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryPayeesUpdated)
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
		it.Event = new(AutomationRegistryPayeesUpdated)
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

func (it *AutomationRegistryPayeesUpdatedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryPayeesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryPayeesUpdated struct {
	Transmitters []common.Address
	Payees       []common.Address
	Raw          types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterPayeesUpdated(opts *bind.FilterOpts) (*AutomationRegistryPayeesUpdatedIterator, error) {

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryPayeesUpdatedIterator{contract: _AutomationRegistry.contract, event: "PayeesUpdated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryPayeesUpdated) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryPayeesUpdated)
				if err := _AutomationRegistry.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParsePayeesUpdated(log types.Log) (*AutomationRegistryPayeesUpdated, error) {
	event := new(AutomationRegistryPayeesUpdated)
	if err := _AutomationRegistry.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryPayeeshipTransferRequestedIterator struct {
	Event *AutomationRegistryPayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryPayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryPayeeshipTransferRequested)
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
		it.Event = new(AutomationRegistryPayeeshipTransferRequested)
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

func (it *AutomationRegistryPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryPayeeshipTransferRequested struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryPayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryPayeeshipTransferRequestedIterator{contract: _AutomationRegistry.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryPayeeshipTransferRequested)
				if err := _AutomationRegistry.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParsePayeeshipTransferRequested(log types.Log) (*AutomationRegistryPayeeshipTransferRequested, error) {
	event := new(AutomationRegistryPayeeshipTransferRequested)
	if err := _AutomationRegistry.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryPayeeshipTransferredIterator struct {
	Event *AutomationRegistryPayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryPayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryPayeeshipTransferred)
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
		it.Event = new(AutomationRegistryPayeeshipTransferred)
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

func (it *AutomationRegistryPayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryPayeeshipTransferred struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryPayeeshipTransferredIterator, error) {

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

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryPayeeshipTransferredIterator{contract: _AutomationRegistry.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryPayeeshipTransferred)
				if err := _AutomationRegistry.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParsePayeeshipTransferred(log types.Log) (*AutomationRegistryPayeeshipTransferred, error) {
	event := new(AutomationRegistryPayeeshipTransferred)
	if err := _AutomationRegistry.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryPaymentWithdrawnIterator struct {
	Event *AutomationRegistryPaymentWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryPaymentWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryPaymentWithdrawn)
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
		it.Event = new(AutomationRegistryPaymentWithdrawn)
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

func (it *AutomationRegistryPaymentWithdrawnIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryPaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryPaymentWithdrawn struct {
	Transmitter common.Address
	Amount      *big.Int
	To          common.Address
	Payee       common.Address
	Raw         types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*AutomationRegistryPaymentWithdrawnIterator, error) {

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

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryPaymentWithdrawnIterator{contract: _AutomationRegistry.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryPaymentWithdrawn)
				if err := _AutomationRegistry.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParsePaymentWithdrawn(log types.Log) (*AutomationRegistryPaymentWithdrawn, error) {
	event := new(AutomationRegistryPaymentWithdrawn)
	if err := _AutomationRegistry.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryReorgedUpkeepReportIterator struct {
	Event *AutomationRegistryReorgedUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryReorgedUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryReorgedUpkeepReport)
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
		it.Event = new(AutomationRegistryReorgedUpkeepReport)
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

func (it *AutomationRegistryReorgedUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryReorgedUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryReorgedUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryReorgedUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryReorgedUpkeepReportIterator{contract: _AutomationRegistry.contract, event: "ReorgedUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryReorgedUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryReorgedUpkeepReport)
				if err := _AutomationRegistry.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseReorgedUpkeepReport(log types.Log) (*AutomationRegistryReorgedUpkeepReport, error) {
	event := new(AutomationRegistryReorgedUpkeepReport)
	if err := _AutomationRegistry.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryStaleUpkeepReportIterator struct {
	Event *AutomationRegistryStaleUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryStaleUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryStaleUpkeepReport)
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
		it.Event = new(AutomationRegistryStaleUpkeepReport)
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

func (it *AutomationRegistryStaleUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryStaleUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryStaleUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryStaleUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryStaleUpkeepReportIterator{contract: _AutomationRegistry.contract, event: "StaleUpkeepReport", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryStaleUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryStaleUpkeepReport)
				if err := _AutomationRegistry.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseStaleUpkeepReport(log types.Log) (*AutomationRegistryStaleUpkeepReport, error) {
	event := new(AutomationRegistryStaleUpkeepReport)
	if err := _AutomationRegistry.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryTransmittedIterator struct {
	Event *AutomationRegistryTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryTransmitted)
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
		it.Event = new(AutomationRegistryTransmitted)
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

func (it *AutomationRegistryTransmittedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterTransmitted(opts *bind.FilterOpts) (*AutomationRegistryTransmittedIterator, error) {

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryTransmittedIterator{contract: _AutomationRegistry.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *AutomationRegistryTransmitted) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryTransmitted)
				if err := _AutomationRegistry.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseTransmitted(log types.Log) (*AutomationRegistryTransmitted, error) {
	event := new(AutomationRegistryTransmitted)
	if err := _AutomationRegistry.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUnpausedIterator struct {
	Event *AutomationRegistryUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUnpaused)
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
		it.Event = new(AutomationRegistryUnpaused)
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

func (it *AutomationRegistryUnpausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUnpaused(opts *bind.FilterOpts) (*AutomationRegistryUnpausedIterator, error) {

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUnpausedIterator{contract: _AutomationRegistry.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUnpaused) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUnpaused)
				if err := _AutomationRegistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUnpaused(log types.Log) (*AutomationRegistryUnpaused, error) {
	event := new(AutomationRegistryUnpaused)
	if err := _AutomationRegistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepAdminTransferRequestedIterator struct {
	Event *AutomationRegistryUpkeepAdminTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepAdminTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepAdminTransferRequested)
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
		it.Event = new(AutomationRegistryUpkeepAdminTransferRequested)
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

func (it *AutomationRegistryUpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryUpkeepAdminTransferRequestedIterator, error) {

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

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepAdminTransferRequestedIterator{contract: _AutomationRegistry.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepAdminTransferRequested)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepAdminTransferRequested(log types.Log) (*AutomationRegistryUpkeepAdminTransferRequested, error) {
	event := new(AutomationRegistryUpkeepAdminTransferRequested)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepAdminTransferredIterator struct {
	Event *AutomationRegistryUpkeepAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepAdminTransferred)
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
		it.Event = new(AutomationRegistryUpkeepAdminTransferred)
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

func (it *AutomationRegistryUpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryUpkeepAdminTransferredIterator, error) {

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

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepAdminTransferredIterator{contract: _AutomationRegistry.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepAdminTransferred)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepAdminTransferred(log types.Log) (*AutomationRegistryUpkeepAdminTransferred, error) {
	event := new(AutomationRegistryUpkeepAdminTransferred)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepCanceledIterator struct {
	Event *AutomationRegistryUpkeepCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepCanceled)
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
		it.Event = new(AutomationRegistryUpkeepCanceled)
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

func (it *AutomationRegistryUpkeepCanceledIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*AutomationRegistryUpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepCanceledIterator{contract: _AutomationRegistry.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepCanceled)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepCanceled(log types.Log) (*AutomationRegistryUpkeepCanceled, error) {
	event := new(AutomationRegistryUpkeepCanceled)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepCheckDataSetIterator struct {
	Event *AutomationRegistryUpkeepCheckDataSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepCheckDataSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepCheckDataSet)
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
		it.Event = new(AutomationRegistryUpkeepCheckDataSet)
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

func (it *AutomationRegistryUpkeepCheckDataSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepCheckDataSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepCheckDataSet struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepCheckDataSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepCheckDataSetIterator{contract: _AutomationRegistry.contract, event: "UpkeepCheckDataSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepCheckDataSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepCheckDataSet)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepCheckDataSet(log types.Log) (*AutomationRegistryUpkeepCheckDataSet, error) {
	event := new(AutomationRegistryUpkeepCheckDataSet)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepGasLimitSetIterator struct {
	Event *AutomationRegistryUpkeepGasLimitSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepGasLimitSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepGasLimitSet)
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
		it.Event = new(AutomationRegistryUpkeepGasLimitSet)
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

func (it *AutomationRegistryUpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepGasLimitSetIterator{contract: _AutomationRegistry.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepGasLimitSet)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepGasLimitSet(log types.Log) (*AutomationRegistryUpkeepGasLimitSet, error) {
	event := new(AutomationRegistryUpkeepGasLimitSet)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepMigratedIterator struct {
	Event *AutomationRegistryUpkeepMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepMigrated)
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
		it.Event = new(AutomationRegistryUpkeepMigrated)
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

func (it *AutomationRegistryUpkeepMigratedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepMigratedIterator{contract: _AutomationRegistry.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepMigrated)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepMigrated(log types.Log) (*AutomationRegistryUpkeepMigrated, error) {
	event := new(AutomationRegistryUpkeepMigrated)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepOffchainConfigSetIterator struct {
	Event *AutomationRegistryUpkeepOffchainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepOffchainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepOffchainConfigSet)
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
		it.Event = new(AutomationRegistryUpkeepOffchainConfigSet)
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

func (it *AutomationRegistryUpkeepOffchainConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepOffchainConfigSet struct {
	Id             *big.Int
	OffchainConfig []byte
	Raw            types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepOffchainConfigSetIterator{contract: _AutomationRegistry.contract, event: "UpkeepOffchainConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepOffchainConfigSet)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepOffchainConfigSet(log types.Log) (*AutomationRegistryUpkeepOffchainConfigSet, error) {
	event := new(AutomationRegistryUpkeepOffchainConfigSet)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepPausedIterator struct {
	Event *AutomationRegistryUpkeepPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepPaused)
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
		it.Event = new(AutomationRegistryUpkeepPaused)
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

func (it *AutomationRegistryUpkeepPausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepPaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepPausedIterator{contract: _AutomationRegistry.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepPaused)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepPaused(log types.Log) (*AutomationRegistryUpkeepPaused, error) {
	event := new(AutomationRegistryUpkeepPaused)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepPerformedIterator struct {
	Event *AutomationRegistryUpkeepPerformed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepPerformedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepPerformed)
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
		it.Event = new(AutomationRegistryUpkeepPerformed)
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

func (it *AutomationRegistryUpkeepPerformedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepPerformed struct {
	Id           *big.Int
	Success      bool
	TotalPayment *big.Int
	GasUsed      *big.Int
	GasOverhead  *big.Int
	Trigger      []byte
	Raw          types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*AutomationRegistryUpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepPerformedIterator{contract: _AutomationRegistry.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepPerformed)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepPerformed(log types.Log) (*AutomationRegistryUpkeepPerformed, error) {
	event := new(AutomationRegistryUpkeepPerformed)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepPrivilegeConfigSetIterator struct {
	Event *AutomationRegistryUpkeepPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepPrivilegeConfigSet)
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
		it.Event = new(AutomationRegistryUpkeepPrivilegeConfigSet)
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

func (it *AutomationRegistryUpkeepPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepPrivilegeConfigSet struct {
	Id              *big.Int
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepPrivilegeConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepPrivilegeConfigSetIterator{contract: _AutomationRegistry.contract, event: "UpkeepPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepPrivilegeConfigSet)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepPrivilegeConfigSet(log types.Log) (*AutomationRegistryUpkeepPrivilegeConfigSet, error) {
	event := new(AutomationRegistryUpkeepPrivilegeConfigSet)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepReceivedIterator struct {
	Event *AutomationRegistryUpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepReceived)
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
		it.Event = new(AutomationRegistryUpkeepReceived)
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

func (it *AutomationRegistryUpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepReceivedIterator{contract: _AutomationRegistry.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepReceived)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepReceived(log types.Log) (*AutomationRegistryUpkeepReceived, error) {
	event := new(AutomationRegistryUpkeepReceived)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepRegisteredIterator struct {
	Event *AutomationRegistryUpkeepRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepRegistered)
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
		it.Event = new(AutomationRegistryUpkeepRegistered)
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

func (it *AutomationRegistryUpkeepRegisteredIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepRegistered struct {
	Id         *big.Int
	PerformGas uint32
	Admin      common.Address
	Raw        types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepRegisteredIterator{contract: _AutomationRegistry.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepRegistered)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepRegistered(log types.Log) (*AutomationRegistryUpkeepRegistered, error) {
	event := new(AutomationRegistryUpkeepRegistered)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepTriggerConfigSetIterator struct {
	Event *AutomationRegistryUpkeepTriggerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepTriggerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepTriggerConfigSet)
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
		it.Event = new(AutomationRegistryUpkeepTriggerConfigSet)
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

func (it *AutomationRegistryUpkeepTriggerConfigSetIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepTriggerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepTriggerConfigSet struct {
	Id            *big.Int
	TriggerConfig []byte
	Raw           types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepTriggerConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepTriggerConfigSetIterator{contract: _AutomationRegistry.contract, event: "UpkeepTriggerConfigSet", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepTriggerConfigSet)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepTriggerConfigSet(log types.Log) (*AutomationRegistryUpkeepTriggerConfigSet, error) {
	event := new(AutomationRegistryUpkeepTriggerConfigSet)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AutomationRegistryUpkeepUnpausedIterator struct {
	Event *AutomationRegistryUpkeepUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryUpkeepUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryUpkeepUnpaused)
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
		it.Event = new(AutomationRegistryUpkeepUnpaused)
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

func (it *AutomationRegistryUpkeepUnpausedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryUpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryUpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryUpkeepUnpausedIterator{contract: _AutomationRegistry.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryUpkeepUnpaused)
				if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseUpkeepUnpaused(log types.Log) (*AutomationRegistryUpkeepUnpaused, error) {
	event := new(AutomationRegistryUpkeepUnpaused)
	if err := _AutomationRegistry.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LatestConfigDetails struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}
type LatestConfigDigestAndEpoch struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}

func (_AutomationRegistry *AutomationRegistry) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _AutomationRegistry.abi.Events["AdminPrivilegeConfigSet"].ID:
		return _AutomationRegistry.ParseAdminPrivilegeConfigSet(log)
	case _AutomationRegistry.abi.Events["CancelledUpkeepReport"].ID:
		return _AutomationRegistry.ParseCancelledUpkeepReport(log)
	case _AutomationRegistry.abi.Events["ConfigSet"].ID:
		return _AutomationRegistry.ParseConfigSet(log)
	case _AutomationRegistry.abi.Events["DedupKeyAdded"].ID:
		return _AutomationRegistry.ParseDedupKeyAdded(log)
	case _AutomationRegistry.abi.Events["FundsAdded"].ID:
		return _AutomationRegistry.ParseFundsAdded(log)
	case _AutomationRegistry.abi.Events["FundsWithdrawn"].ID:
		return _AutomationRegistry.ParseFundsWithdrawn(log)
	case _AutomationRegistry.abi.Events["InsufficientFundsUpkeepReport"].ID:
		return _AutomationRegistry.ParseInsufficientFundsUpkeepReport(log)
	case _AutomationRegistry.abi.Events["OwnerFundsWithdrawn"].ID:
		return _AutomationRegistry.ParseOwnerFundsWithdrawn(log)
	case _AutomationRegistry.abi.Events["OwnershipTransferRequested"].ID:
		return _AutomationRegistry.ParseOwnershipTransferRequested(log)
	case _AutomationRegistry.abi.Events["OwnershipTransferred"].ID:
		return _AutomationRegistry.ParseOwnershipTransferred(log)
	case _AutomationRegistry.abi.Events["Paused"].ID:
		return _AutomationRegistry.ParsePaused(log)
	case _AutomationRegistry.abi.Events["PayeesUpdated"].ID:
		return _AutomationRegistry.ParsePayeesUpdated(log)
	case _AutomationRegistry.abi.Events["PayeeshipTransferRequested"].ID:
		return _AutomationRegistry.ParsePayeeshipTransferRequested(log)
	case _AutomationRegistry.abi.Events["PayeeshipTransferred"].ID:
		return _AutomationRegistry.ParsePayeeshipTransferred(log)
	case _AutomationRegistry.abi.Events["PaymentWithdrawn"].ID:
		return _AutomationRegistry.ParsePaymentWithdrawn(log)
	case _AutomationRegistry.abi.Events["ReorgedUpkeepReport"].ID:
		return _AutomationRegistry.ParseReorgedUpkeepReport(log)
	case _AutomationRegistry.abi.Events["StaleUpkeepReport"].ID:
		return _AutomationRegistry.ParseStaleUpkeepReport(log)
	case _AutomationRegistry.abi.Events["Transmitted"].ID:
		return _AutomationRegistry.ParseTransmitted(log)
	case _AutomationRegistry.abi.Events["Unpaused"].ID:
		return _AutomationRegistry.ParseUnpaused(log)
	case _AutomationRegistry.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _AutomationRegistry.ParseUpkeepAdminTransferRequested(log)
	case _AutomationRegistry.abi.Events["UpkeepAdminTransferred"].ID:
		return _AutomationRegistry.ParseUpkeepAdminTransferred(log)
	case _AutomationRegistry.abi.Events["UpkeepCanceled"].ID:
		return _AutomationRegistry.ParseUpkeepCanceled(log)
	case _AutomationRegistry.abi.Events["UpkeepCheckDataSet"].ID:
		return _AutomationRegistry.ParseUpkeepCheckDataSet(log)
	case _AutomationRegistry.abi.Events["UpkeepGasLimitSet"].ID:
		return _AutomationRegistry.ParseUpkeepGasLimitSet(log)
	case _AutomationRegistry.abi.Events["UpkeepMigrated"].ID:
		return _AutomationRegistry.ParseUpkeepMigrated(log)
	case _AutomationRegistry.abi.Events["UpkeepOffchainConfigSet"].ID:
		return _AutomationRegistry.ParseUpkeepOffchainConfigSet(log)
	case _AutomationRegistry.abi.Events["UpkeepPaused"].ID:
		return _AutomationRegistry.ParseUpkeepPaused(log)
	case _AutomationRegistry.abi.Events["UpkeepPerformed"].ID:
		return _AutomationRegistry.ParseUpkeepPerformed(log)
	case _AutomationRegistry.abi.Events["UpkeepPrivilegeConfigSet"].ID:
		return _AutomationRegistry.ParseUpkeepPrivilegeConfigSet(log)
	case _AutomationRegistry.abi.Events["UpkeepReceived"].ID:
		return _AutomationRegistry.ParseUpkeepReceived(log)
	case _AutomationRegistry.abi.Events["UpkeepRegistered"].ID:
		return _AutomationRegistry.ParseUpkeepRegistered(log)
	case _AutomationRegistry.abi.Events["UpkeepTriggerConfigSet"].ID:
		return _AutomationRegistry.ParseUpkeepTriggerConfigSet(log)
	case _AutomationRegistry.abi.Events["UpkeepUnpaused"].ID:
		return _AutomationRegistry.ParseUpkeepUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (AutomationRegistryAdminPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x7c44b4eb59ee7873514e7e43e7718c269d872965938b288aa143befca62f99d2")
}

func (AutomationRegistryCancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd5636")
}

func (AutomationRegistryConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (AutomationRegistryDedupKeyAdded) Topic() common.Hash {
	return common.HexToHash("0xa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f2")
}

func (AutomationRegistryFundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (AutomationRegistryFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (AutomationRegistryInsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e02")
}

func (AutomationRegistryOwnerFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1")
}

func (AutomationRegistryOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (AutomationRegistryOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (AutomationRegistryPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (AutomationRegistryPayeesUpdated) Topic() common.Hash {
	return common.HexToHash("0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725")
}

func (AutomationRegistryPayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (AutomationRegistryPayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (AutomationRegistryPaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (AutomationRegistryReorgedUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301")
}

func (AutomationRegistryStaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8")
}

func (AutomationRegistryTransmitted) Topic() common.Hash {
	return common.HexToHash("0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62")
}

func (AutomationRegistryUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (AutomationRegistryUpkeepAdminTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35")
}

func (AutomationRegistryUpkeepAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c")
}

func (AutomationRegistryUpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (AutomationRegistryUpkeepCheckDataSet) Topic() common.Hash {
	return common.HexToHash("0xcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d")
}

func (AutomationRegistryUpkeepGasLimitSet) Topic() common.Hash {
	return common.HexToHash("0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c")
}

func (AutomationRegistryUpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (AutomationRegistryUpkeepOffchainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850")
}

func (AutomationRegistryUpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (AutomationRegistryUpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b")
}

func (AutomationRegistryUpkeepPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae7769")
}

func (AutomationRegistryUpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (AutomationRegistryUpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (AutomationRegistryUpkeepTriggerConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664")
}

func (AutomationRegistryUpkeepUnpaused) Topic() common.Hash {
	return common.HexToHash("0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456")
}

func (_AutomationRegistry *AutomationRegistry) Address() common.Address {
	return _AutomationRegistry.address
}

type AutomationRegistryInterface interface {
	FallbackTo(opts *bind.CallOpts) (common.Address, error)

	LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

		error)

	LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	SetConfigTypeSafe(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig AutomationRegistryBase22OnchainConfig, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	SimulatePerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error)

	FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*AutomationRegistryAdminPrivilegeConfigSetIterator, error)

	WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryAdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error)

	ParseAdminPrivilegeConfigSet(log types.Log) (*AutomationRegistryAdminPrivilegeConfigSet, error)

	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryCancelledUpkeepReportIterator, error)

	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryCancelledUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseCancelledUpkeepReport(log types.Log) (*AutomationRegistryCancelledUpkeepReport, error)

	FilterConfigSet(opts *bind.FilterOpts) (*AutomationRegistryConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*AutomationRegistryConfigSet, error)

	FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*AutomationRegistryDedupKeyAddedIterator, error)

	WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryDedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error)

	ParseDedupKeyAdded(log types.Log) (*AutomationRegistryDedupKeyAdded, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*AutomationRegistryFundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *AutomationRegistryFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*AutomationRegistryFundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryFundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryFundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*AutomationRegistryFundsWithdrawn, error)

	FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryInsufficientFundsUpkeepReportIterator, error)

	WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseInsufficientFundsUpkeepReport(log types.Log) (*AutomationRegistryInsufficientFundsUpkeepReport, error)

	FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*AutomationRegistryOwnerFundsWithdrawnIterator, error)

	WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryOwnerFundsWithdrawn) (event.Subscription, error)

	ParseOwnerFundsWithdrawn(log types.Log) (*AutomationRegistryOwnerFundsWithdrawn, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*AutomationRegistryOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutomationRegistryOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*AutomationRegistryOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*AutomationRegistryPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*AutomationRegistryPaused, error)

	FilterPayeesUpdated(opts *bind.FilterOpts) (*AutomationRegistryPayeesUpdatedIterator, error)

	WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryPayeesUpdated) (event.Subscription, error)

	ParsePayeesUpdated(log types.Log) (*AutomationRegistryPayeesUpdated, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryPayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*AutomationRegistryPayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*AutomationRegistryPayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*AutomationRegistryPayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*AutomationRegistryPaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *AutomationRegistryPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*AutomationRegistryPaymentWithdrawn, error)

	FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryReorgedUpkeepReportIterator, error)

	WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryReorgedUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseReorgedUpkeepReport(log types.Log) (*AutomationRegistryReorgedUpkeepReport, error)

	FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryStaleUpkeepReportIterator, error)

	WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *AutomationRegistryStaleUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseStaleUpkeepReport(log types.Log) (*AutomationRegistryStaleUpkeepReport, error)

	FilterTransmitted(opts *bind.FilterOpts) (*AutomationRegistryTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *AutomationRegistryTransmitted) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*AutomationRegistryTransmitted, error)

	FilterUnpaused(opts *bind.FilterOpts) (*AutomationRegistryUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*AutomationRegistryUnpaused, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryUpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*AutomationRegistryUpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*AutomationRegistryUpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*AutomationRegistryUpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*AutomationRegistryUpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*AutomationRegistryUpkeepCanceled, error)

	FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepCheckDataSetIterator, error)

	WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepCheckDataSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCheckDataSet(log types.Log) (*AutomationRegistryUpkeepCheckDataSet, error)

	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepGasLimitSetIterator, error)

	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepGasLimitSet(log types.Log) (*AutomationRegistryUpkeepGasLimitSet, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*AutomationRegistryUpkeepMigrated, error)

	FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepOffchainConfigSetIterator, error)

	WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepOffchainConfigSet(log types.Log) (*AutomationRegistryUpkeepOffchainConfigSet, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*AutomationRegistryUpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*AutomationRegistryUpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*AutomationRegistryUpkeepPerformed, error)

	FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepPrivilegeConfigSetIterator, error)

	WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPrivilegeConfigSet(log types.Log) (*AutomationRegistryUpkeepPrivilegeConfigSet, error)

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*AutomationRegistryUpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*AutomationRegistryUpkeepRegistered, error)

	FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepTriggerConfigSetIterator, error)

	WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepTriggerConfigSet(log types.Log) (*AutomationRegistryUpkeepTriggerConfigSet, error)

	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*AutomationRegistryUpkeepUnpausedIterator, error)

	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *AutomationRegistryUpkeepUnpaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepUnpaused(log types.Log) (*AutomationRegistryUpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
