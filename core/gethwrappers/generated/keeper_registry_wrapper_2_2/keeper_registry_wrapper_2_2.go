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
	ChainModule            common.Address
	ReorgProtectionEnabled bool
}

var AutomationRegistryMetaData = &bind.MetaData{
<<<<<<< HEAD
<<<<<<< HEAD
	ABI: "[{\"inputs\":[{\"internalType\":\"contractAutomationRegistryLogicB2_2\",\"name\":\"logicA\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTrigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTriggerType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"AdminPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"DedupKeyAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepTriggerConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fallbackTo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfigBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"reorgProtectionEnabled\",\"type\":\"bool\"}],\"internalType\":\"structAutomationRegistryBase2_2.OnchainConfig\",\"name\":\"onchainConfig\",\"type\":\"tuple\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfigTypeSafe\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"simulatePerformUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"rawReport\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6101606040523480156200001257600080fd5b506040516200563e3803806200563e83398101604081905262000035916200044b565b80816001600160a01b0316634b4fd03b6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200009b919062000472565b826001600160a01b031663ca30e6036040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000da573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200010091906200044b565b836001600160a01b031663b10b673c6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200013f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200016591906200044b565b846001600160a01b0316636709d0e56040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca91906200044b565b856001600160a01b0316635425d8ac6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000209573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200022f91906200044b565b866001600160a01b031663a08714c06040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200026e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200029491906200044b565b3380600081620002eb5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200031e576200031e8162000387565b50505085600281111562000336576200033662000495565b60e08160028111156200034d576200034d62000495565b9052506001600160a01b0394851660805292841660a05290831660c052821661010052811661012052919091166101405250620004ab9050565b336001600160a01b03821603620003e15760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620002e2565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200044857600080fd5b50565b6000602082840312156200045e57600080fd5b81516200046b8162000432565b9392505050565b6000602082840312156200048557600080fd5b8151600381106200046b57600080fd5b634e487b7160e01b600052602160045260246000fd5b60805160a05160c05160e051610100516101205161014051615126620005186000396000818160d6015261016f01526000611e2401526000505060008181611c4c0152818161346e015281816136010152613b1e01526000505060005050600061134a01526151266000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c8063a4c0ed3611610081578063b1dc65a41161005b578063b1dc65a4146102e0578063e3d0e712146102f3578063f2fde38b14610306576100d4565b8063a4c0ed3614610262578063aed2e92914610275578063afcb95d71461029f576100d4565b806379ba5097116100b257806379ba5097146101c757806381ff7048146101cf5780638da5cb5b14610244576100d4565b8063181f5a771461011b578063349e8cca1461016d5780636cad5469146101b4575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e808015610114573d6000f35b3d6000fd5b005b6101576040518060400160405280601881526020017f4175746f6d6174696f6e526567697374727920322e322e30000000000000000081525081565b6040516101649190613d9d565b60405180910390f35b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610164565b6101196101c23660046141cc565b610319565b610119611230565b61022160155460115463ffffffff780100000000000000000000000000000000000000000000000083048116937c01000000000000000000000000000000000000000000000000000000009093041691565b6040805163ffffffff948516815293909216602084015290820152606001610164565b60005473ffffffffffffffffffffffffffffffffffffffff1661018f565b6101196102703660046142e2565b611332565b61028861028336600461433e565b61154e565b604080519215158352602083019190915201610164565b601154601254604080516000815260208101939093527401000000000000000000000000000000000000000090910463ffffffff1690820152606001610164565b6101196102ee3660046143cf565b6116c6565b610119610301366004614486565b61197e565b610119610314366004614515565b6119a7565b6103216119bb565b601f8651111561035d576040517f25d0209c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8360ff1660000361039a576040517fe77dba5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845186511415806103b957506103b1846003614561565b60ff16865111155b156103f0576040517f1d2d1c5800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601254600e546bffffffffffffffffffffffff9091169060005b816bffffffffffffffffffffffff168110156104725761045f600e82815481106104365761043661457d565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff168484611a3e565b508061046a816145ac565b91505061040a565b5060008060005b836bffffffffffffffffffffffff1681101561057b57600d81815481106104a2576104a261457d565b600091825260209091200154600e805473ffffffffffffffffffffffffffffffffffffffff909216945090829081106104dd576104dd61457d565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff8681168452600c8352604080852080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001690559116808452600b90925290912080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055915080610573816145ac565b915050610479565b50610588600d6000613c72565b610594600e6000613c72565b604080516080810182526000808252602082018190529181018290526060810182905290805b8c518110156109fd57600c60008e83815181106105d9576105d961457d565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528101919091526040016000205460ff1615610644576040517f77cea0fa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168d828151811061066e5761066e61457d565b602002602001015173ffffffffffffffffffffffffffffffffffffffff16036106c3576040517f815e1d6400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405180604001604052806001151581526020018260ff16815250600c60008f84815181106106f4576106f461457d565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528181019290925260400160002082518154939092015160ff16610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909316929092171790558b518c908290811061079c5761079c61457d565b60200260200101519150600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160361080c576040517f58a70a0a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff82166000908152600b60209081526040918290208251608081018452905460ff80821615801584526101008304909116938301939093526bffffffffffffffffffffffff6201000082048116948301949094526e010000000000000000000000000000900490921660608301529093506108c7576040517f6a7281ad00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001835260ff80821660208086019182526bffffffffffffffffffffffff808b166060880190815273ffffffffffffffffffffffffffffffffffffffff87166000908152600b909352604092839020885181549551948a0151925184166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff939094166201000002929092167fffffffffffff000000000000000000000000000000000000000000000000ffff94909616610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000090951694909417179190911692909217919091179055806109f5816145ac565b9150506105ba565b50508a51610a139150600d9060208d0190613c90565b508851610a2790600e9060208c0190613c90565b50604051806101400160405280856bffffffffffffffffffffffff168152602001886000015163ffffffff168152602001886020015163ffffffff168152602001600063ffffffff168152602001886060015162ffffff168152602001886080015161ffff1681526020018960ff1681526020016012600001601e9054906101000a900460ff16151581526020016012600001601f9054906101000a900460ff1615158152602001886101e001511515815250601260008201518160000160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550602082015181600001600c6101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160106101000a81548163ffffffff021916908363ffffffff16021790555060608201518160000160146101000a81548163ffffffff021916908363ffffffff16021790555060808201518160000160186101000a81548162ffffff021916908362ffffff16021790555060a082015181600001601b6101000a81548161ffff021916908361ffff16021790555060c082015181600001601d6101000a81548160ff021916908360ff16021790555060e082015181600001601e6101000a81548160ff02191690831515021790555061010082015181600001601f6101000a81548160ff0219169083151502179055506101208201518160010160006101000a81548160ff0219169083151502179055509050506040518061018001604052808860a001516bffffffffffffffffffffffff16815260200188610180015173ffffffffffffffffffffffffffffffffffffffff168152602001601460010160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff168152602001886040015163ffffffff1681526020018860c0015163ffffffff168152602001601460010160149054906101000a900463ffffffff1663ffffffff168152602001601460010160189054906101000a900463ffffffff1663ffffffff1681526020016014600101601c9054906101000a900463ffffffff1663ffffffff1681526020018860e0015163ffffffff16815260200188610100015163ffffffff16815260200188610120015163ffffffff168152602001886101c0015173ffffffffffffffffffffffffffffffffffffffff16815250601460008201518160000160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550602082015181600001600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060408201518160010160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550606082015181600101600c6101000a81548163ffffffff021916908363ffffffff16021790555060808201518160010160106101000a81548163ffffffff021916908363ffffffff16021790555060a08201518160010160146101000a81548163ffffffff021916908363ffffffff16021790555060c08201518160010160186101000a81548163ffffffff021916908363ffffffff16021790555060e082015181600101601c6101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160020160006101000a81548163ffffffff021916908363ffffffff1602179055506101208201518160020160046101000a81548163ffffffff021916908363ffffffff1602179055506101408201518160020160086101000a81548163ffffffff021916908363ffffffff16021790555061016082015181600201600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505086610140015160178190555086610160015160188190555060006014600101601c9054906101000a900463ffffffff169050611017611c46565b601580547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167c010000000000000000000000000000000000000000000000000000000063ffffffff938416021780825560019260189161109291859178010000000000000000000000000000000000000000000000009004166145e4565b92506101000a81548163ffffffff021916908363ffffffff1602179055506000886040516020016110c39190614652565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291905260155490915061112c90469030907801000000000000000000000000000000000000000000000000900463ffffffff168f8f8f878f8f611cfb565b60115560005b61113c6009611da5565b81101561116c57611159611151600983611db5565b600990611dc8565b5080611164816145ac565b915050611132565b5060005b896101a00151518110156111c3576111b08a6101a0015182815181106111985761119861457d565b60200260200101516009611dea90919063ffffffff16565b50806111bb816145ac565b915050611170565b507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0582601154601460010160189054906101000a900463ffffffff168f8f8f878f8f60405161121a999897969594939291906147cd565b60405180910390a1505050505050505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146112b6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146113a1576040517fc8bad78d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b602081146113db576040517fdfe9309000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006113e982840184614863565b60008181526004602052604090205490915065010000000000900463ffffffff90811614611443576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526004602052604090206001015461147e9085906c0100000000000000000000000090046bffffffffffffffffffffffff1661487c565b600082815260046020526040902060010180546bffffffffffffffffffffffff929092166c01000000000000000000000000027fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff9092169190911790556019546114e99085906148a1565b6019556040516bffffffffffffffffffffffff8516815273ffffffffffffffffffffffffffffffffffffffff86169082907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a35050505050565b600080611559611e0c565b6012547e01000000000000000000000000000000000000000000000000000000000000900460ff16156115b8576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600085815260046020908152604091829020825160e081018452815460ff811615158252610100810463ffffffff908116838601819052650100000000008304821684880152690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff16606084018190526001909401546bffffffffffffffffffffffff80821660808601526c0100000000000000000000000082041660a0850152780100000000000000000000000000000000000000000000000090041660c08301528451601f890185900485028101850190955287855290936116b7938990899081908401838280828437600092019190915250611e7b92505050565b9093509150505b935093915050565b60005a60408051610140810182526012546bffffffffffffffffffffffff8116825263ffffffff6c010000000000000000000000008204811660208401527001000000000000000000000000000000008204811693830193909352740100000000000000000000000000000000000000008104909216606082015262ffffff7801000000000000000000000000000000000000000000000000830416608082015261ffff7b0100000000000000000000000000000000000000000000000000000083041660a082015260ff7d0100000000000000000000000000000000000000000000000000000000008304811660c08301527e010000000000000000000000000000000000000000000000000000000000008304811615801560e08401527f010000000000000000000000000000000000000000000000000000000000000090930481161515610100830152601354161515610120820152919250611858576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600b602052604090205460ff166118a1576040517f1099ed7500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6011548a35146118dd576040517fdfdcf8e700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60c08101516118ed9060016148b4565b60ff16861415806118fe5750858414155b15611935576040517f0244f71a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6119458a8a8a8a8a8a8a8a6120a4565b60006119518a8a61230d565b905060208b0135600881901c63ffffffff1661196f848487846123c6565b50505050505050505050505050565b61199f868686868060200190518101906119989190614973565b8686610319565b505050505050565b6119af6119bb565b6119b881612ca7565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314611a3c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016112ad565b565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e0100000000000000000000000000009004909116606082015290611c3a576000816060015185611ad69190614ae1565b90506000611ae48583614b35565b90508083604001818151611af8919061487c565b6bffffffffffffffffffffffff16905250611b138582614b60565b83606001818151611b24919061487c565b6bffffffffffffffffffffffff90811690915273ffffffffffffffffffffffffffffffffffffffff89166000908152600b602090815260409182902087518154928901519389015160608a015186166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff919096166201000002167fffffffffffff000000000000000000000000000000000000000000000000ffff60ff95909516610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909416939093171792909216179190911790555050505b60400151949350505050565b600060017f00000000000000000000000000000000000000000000000000000000000000006002811115611c7c57611c7c614b90565b03611cf657606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015611ccd573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611cf19190614bbf565b905090565b504390565b6000808a8a8a8a8a8a8a8a8a604051602001611d1f99989796959493929190614bd8565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179b9a5050505050505050505050565b6000611daf825490565b92915050565b6000611dc18383612d9c565b9392505050565b6000611dc18373ffffffffffffffffffffffffffffffffffffffff8416612dc6565b6000611dc18373ffffffffffffffffffffffffffffffffffffffff8416612ec0565b3273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614611a3c576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60125460009081907f0100000000000000000000000000000000000000000000000000000000000000900460ff1615611ee0576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601280547effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f01000000000000000000000000000000000000000000000000000000000000001790556040517f4585e33b0000000000000000000000000000000000000000000000000000000090611f5c908590602401613d9d565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009094169390931790925290517f79188d1600000000000000000000000000000000000000000000000000000000815290935073ffffffffffffffffffffffffffffffffffffffff8616906379188d169061202f9087908790600401614c6d565b60408051808303816000875af115801561204d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906120719190614c86565b601280547effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff16905590969095509350505050565b600087876040516120b6929190614cb4565b6040519081900381206120cd918b90602001614cc4565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201208383019092526000808452908301819052909250906000805b888110156122a4576001858783602081106121395761213961457d565b61214691901a601b6148b4565b8c8c858181106121585761215861457d565b905060200201358b8b868181106121715761217161457d565b90506020020135604051600081526020016040526040516121ae949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa1580156121d0573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff81166000908152600c602090815290849020838501909452925460ff808216151580855261010090920416938301939093529095509350905061227e576040517f0f4c073700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826020015160080260ff166001901b84019350808061229c906145ac565b91505061211c565b50827e010101010101010101010101010101010101010101010101010101010101018416146122ff576040517fc103be2e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505050505050505050505050565b6123466040518060c001604052806000815260200160008152602001606081526020016060815260200160608152602001606081525090565b600061235483850185614db5565b604081015151606082015151919250908114158061237757508082608001515114155b806123875750808260a001515114155b156123be576040517fb55ac75400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b509392505050565b600083604001515167ffffffffffffffff8111156123e6576123e6613db0565b6040519080825280602002602001820160405280156124aa57816020015b604080516101e0810182526000610100820181815261012083018290526101408301829052610160830182905261018083018290526101a083018290526101c0830182905282526020808301829052928201819052606082018190526080820181905260a0820181905260c0820181905260e082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816124045790505b5090506000805b8560400151518110156128f45760046000876040015183815181106124d8576124d861457d565b6020908102919091018101518252818101929092526040908101600020815160e081018352815460ff811615158252610100810463ffffffff90811695830195909552650100000000008104851693820193909352690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff166060830152600101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490911660c082015283518490839081106125bd576125bd61457d565b6020026020010151600001819052506125f2866040015182815181106125e5576125e561457d565b6020026020010151612f0f565b8382815181106126045761260461457d565b602002602001015160800190600181111561262157612621614b90565b9081600181111561263457612634614b90565b815250506126a88784838151811061264e5761264e61457d565b602002602001015160800151886060015184815181106126705761267061457d565b60200260200101518960a00151858151811061268e5761268e61457d565b6020026020010151518a600001518b602001516001612fba565b8382815181106126ba576126ba61457d565b6020026020010151604001906bffffffffffffffffffffffff1690816bffffffffffffffffffffffff1681525050612747866040015182815181106127015761270161457d565b60200260200101518760800151838151811061271f5761271f61457d565b60200260200101518584815181106127395761273961457d565b60200260200101518a613005565b8483815181106127595761275961457d565b60200260200101516020018584815181106127765761277661457d565b602002602001015160e00182815250821515151581525050508281815181106127a1576127a161457d565b602002602001015160200151156127c4576127bd600183614ea2565b91506127c9565b6128e2565b61282f8382815181106127de576127de61457d565b60200260200101516000015160600151876060015183815181106128045761280461457d565b60200260200101518860a0015184815181106128225761282261457d565b6020026020010151611e7b565b8483815181106128415761284161457d565b602002602001015160600185848151811061285e5761285e61457d565b602002602001015160a00182815250821515151581525050508281815181106128895761288961457d565b602002602001015160a00151856128a09190614ebd565b94506128e2866040015182815181106128bb576128bb61457d565b60200260200101518483815181106128d5576128d561457d565b6020026020010151613188565b806128ec816145ac565b9150506124b1565b508061ffff16600003612908575050612ca1565b60c08601516129189060016148b4565b6129279060ff1661044c614ed0565b616b6c612935366010614ed0565b5a6129409088614ebd565b61294a91906148a1565b61295491906148a1565b61295e91906148a1565b9350611b5861297161ffff831686614ee7565b61297b91906148a1565b935060008060008060005b896040015151811015612b7c578681815181106129a5576129a561457d565b60200260200101516020015115612b6a57612a01898883815181106129cc576129cc61457d565b6020026020010151608001518c60a0015184815181106129ee576129ee61457d565b6020026020010151518e60c0015161329a565b878281518110612a1357612a1361457d565b602002602001015160c0018181525050612a6f8b8b604001518381518110612a3d57612a3d61457d565b6020026020010151898481518110612a5757612a5761457d565b60200260200101518d600001518e602001518b6132ba565b9093509150612a7e828561487c565b9350612a8a838661487c565b9450868181518110612a9e57612a9e61457d565b60200260200101516060015115158a604001518281518110612ac257612ac261457d565b60200260200101517fad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b8486612af7919061487c565b8a8581518110612b0957612b0961457d565b602002602001015160a001518b8681518110612b2757612b2761457d565b602002602001015160c001518f608001518781518110612b4957612b4961457d565b6020026020010151604051612b619493929190614efb565b60405180910390a35b80612b74816145ac565b915050612986565b5050336000908152600b602052604090208054849250600290612bb49084906201000090046bffffffffffffffffffffffff1661487c565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555080601260000160008282829054906101000a90046bffffffffffffffffffffffff16612c0e919061487c565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550876060015163ffffffff168563ffffffff161115612c9c57601280547fffffffffffffffff00000000ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000063ffffffff8816021790555b505050505b50505050565b3373ffffffffffffffffffffffffffffffffffffffff821603612d26576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016112ad565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000826000018281548110612db357612db361457d565b9060005260206000200154905092915050565b60008181526001830160205260408120548015612eaf576000612dea600183614ebd565b8554909150600090612dfe90600190614ebd565b9050818114612e63576000866000018281548110612e1e57612e1e61457d565b9060005260206000200154905080876000018481548110612e4157612e4161457d565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080612e7457612e74614f38565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050611daf565b6000915050611daf565b5092915050565b6000818152600183016020526040812054612f0757508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155611daf565b506000611daf565b6000818160045b600f811015612f9c577fff000000000000000000000000000000000000000000000000000000000000008216838260208110612f5457612f5461457d565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614612f8a57506000949350505050565b80612f94816145ac565b915050612f16565b5081600f1a6001811115612fb257612fb2614b90565b949350505050565b600080612fcc88878b60c001516133ad565b9050600080612fe78b8a63ffffffff16858a8a60018b613439565b9092509050612ff6818361487c565b9b9a5050505050505050505050565b60008080808560800151600181111561302057613020614b90565b036130455761303187878787613892565b6130405760009250905061317f565b6130bc565b60018560800151600181111561305d5761305d614b90565b0361308a57600061306f888887613994565b9250905080613084575060009250905061317f565b506130bc565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6130c4611c46565b85516040015163ffffffff161161311857867fc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd5636876040516131059190613d9d565b60405180910390a260009250905061317f565b84604001516bffffffffffffffffffffffff16856000015160a001516bffffffffffffffffffffffff16101561317857867f377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e02876040516131059190613d9d565b6001925090505b94509492505050565b6000816080015160018111156131a0576131a0614b90565b03613212576131ad611c46565b6000838152600460205260409020600101805463ffffffff929092167801000000000000000000000000000000000000000000000000027fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff9092169190911790555050565b60018160800151600181111561322a5761322a614b90565b036132965760e08101805160009081526008602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055915191517fa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f29190a25b5050565b60006132a78484846133ad565b905080851015612fb25750929392505050565b6000806132d5888760a001518860c001518888886001613439565b909250905060006132e6828461487c565b600089815260046020526040902060010180549192508291600c9061332a9084906c0100000000000000000000000090046bffffffffffffffffffffffff16614ae1565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915560008a8152600460205260408120600101805485945090926133739185911661487c565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555050965096945050505050565b600080808560018111156133c3576133c3614b90565b036133d2575062015f906133f1565b60018560018111156133e6576133e6614b90565b0361308a57506201adb05b61340263ffffffff85166014614ed0565b61340d8460016148b4565b61341c9060ff16611d4c614ed0565b61342690836148a1565b61343091906148a1565b95945050505050565b60008060008960a0015161ffff16876134529190614ed0565b90508380156134605750803a105b1561346857503a5b600060027f0000000000000000000000000000000000000000000000000000000000000000600281111561349e5761349e614b90565b036135fd5760408051600081526020810190915285156134fc576000366040518060800160405280604881526020016150d2604891396040516020016134e693929190614f67565b6040516020818303038152906040529050613564565b60165461351890640100000000900463ffffffff166004614f8e565b63ffffffff1667ffffffffffffffff81111561353657613536613db0565b6040519080825280601f01601f191660200182016040528015613560576020820181803683370190505b5090505b6040517f49948e0e00000000000000000000000000000000000000000000000000000000815273420000000000000000000000000000000000000f906349948e0e906135b4908490600401613d9d565b602060405180830381865afa1580156135d1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906135f59190614bbf565b915050613757565b60017f0000000000000000000000000000000000000000000000000000000000000000600281111561363157613631614b90565b036137575784156136b357606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613688573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906136ac9190614bbf565b9050613757565b6000606c73ffffffffffffffffffffffffffffffffffffffff166341b247a86040518163ffffffff1660e01b815260040160c060405180830381865afa158015613701573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906137259190614fae565b505060165492945061374893505050640100000000900463ffffffff1682614ed0565b613753906010614ed0565b9150505b8461377357808b60a0015161ffff166137709190614ed0565b90505b61378161ffff871682614ee7565b9050600087826137918c8e6148a1565b61379b9086614ed0565b6137a591906148a1565b6137b790670de0b6b3a7640000614ed0565b6137c19190614ee7565b905060008c6040015163ffffffff1664e8d4a510006137e09190614ed0565b898e6020015163ffffffff16858f886137f99190614ed0565b61380391906148a1565b61381190633b9aca00614ed0565b61381b9190614ed0565b6138259190614ee7565b61382f91906148a1565b90506b033b2e3c9fd0803ce800000061384882846148a1565b1115613880576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b600080848060200190518101906138a99190614ff8565b845160c00151815191925063ffffffff9081169116101561390657857f405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8866040516138f49190613d9d565b60405180910390a26000915050612fb2565b826101200151801561393a575060208101511580159061393a5750602081015181516139379063ffffffff16613b18565b14155b806139535750613948611c46565b815163ffffffff1610155b1561398857857f6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301866040516138f49190613d9d565b50600195945050505050565b6000806000848060200190518101906139ad9190615050565b9050600086826000015183602001518460400151604051602001613a0f94939291909384526020840192909252604083015260e01b7fffffffff0000000000000000000000000000000000000000000000000000000016606082015260640190565b6040516020818303038152906040528051906020012090508461012001518015613a5d5750608082015115801590613a5d57508160800151613a5a836060015163ffffffff16613b18565b14155b80613a795750613a6b611c46565b826060015163ffffffff1610155b15613ac357867f6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc30187604051613aae9190613d9d565b60405180910390a26000935091506116be9050565b60008181526008602052604090205460ff1615613b0a57867f405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e887604051613aae9190613d9d565b600197909650945050505050565b600060017f00000000000000000000000000000000000000000000000000000000000000006002811115613b4e57613b4e614b90565b03613c68576000606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613ba1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613bc59190614bbf565b90508083101580613be05750610100613bde8483614ebd565b115b15613bee5750600092915050565b6040517f2b407a8200000000000000000000000000000000000000000000000000000000815260048101849052606490632b407a8290602401602060405180830381865afa158015613c44573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611dc19190614bbf565b504090565b919050565b50805460008255906000526020600020908101906119b89190613d1a565b828054828255906000526020600020908101928215613d0a579160200282015b82811115613d0a57825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190613cb0565b50613d16929150613d1a565b5090565b5b80821115613d165760008155600101613d1b565b60005b83811015613d4a578181015183820152602001613d32565b50506000910152565b60008151808452613d6b816020860160208601613d2f565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000611dc16020830184613d53565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610200810167ffffffffffffffff81118282101715613e0357613e03613db0565b60405290565b60405160c0810167ffffffffffffffff81118282101715613e0357613e03613db0565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613e7357613e73613db0565b604052919050565b600067ffffffffffffffff821115613e9557613e95613db0565b5060051b60200190565b73ffffffffffffffffffffffffffffffffffffffff811681146119b857600080fd5b8035613c6d81613e9f565b600082601f830112613edd57600080fd5b81356020613ef2613eed83613e7b565b613e2c565b82815260059290921b84018101918181019086841115613f1157600080fd5b8286015b84811015613f35578035613f2881613e9f565b8352918301918301613f15565b509695505050505050565b803560ff81168114613c6d57600080fd5b63ffffffff811681146119b857600080fd5b8035613c6d81613f51565b62ffffff811681146119b857600080fd5b8035613c6d81613f6e565b61ffff811681146119b857600080fd5b8035613c6d81613f8a565b6bffffffffffffffffffffffff811681146119b857600080fd5b8035613c6d81613fa5565b80151581146119b857600080fd5b8035613c6d81613fca565b60006102008284031215613ff657600080fd5b613ffe613ddf565b905061400982613f63565b815261401760208301613f63565b602082015261402860408301613f63565b604082015261403960608301613f7f565b606082015261404a60808301613f9a565b608082015261405b60a08301613fbf565b60a082015261406c60c08301613f63565b60c082015261407d60e08301613f63565b60e0820152610100614090818401613f63565b908201526101206140a2838201613f63565b90820152610140828101359082015261016080830135908201526101806140ca818401613ec1565b908201526101a08281013567ffffffffffffffff8111156140ea57600080fd5b6140f685828601613ecc565b8284015250506101c061410a818401613ec1565b908201526101e061411c838201613fd8565b9082015292915050565b803567ffffffffffffffff81168114613c6d57600080fd5b600082601f83011261414f57600080fd5b813567ffffffffffffffff81111561416957614169613db0565b61419a60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601613e2c565b8181528460208386010111156141af57600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060c087890312156141e557600080fd5b863567ffffffffffffffff808211156141fd57600080fd5b6142098a838b01613ecc565b9750602089013591508082111561421f57600080fd5b61422b8a838b01613ecc565b965061423960408a01613f40565b9550606089013591508082111561424f57600080fd5b61425b8a838b01613fe3565b945061426960808a01614126565b935060a089013591508082111561427f57600080fd5b5061428c89828a0161413e565b9150509295509295509295565b60008083601f8401126142ab57600080fd5b50813567ffffffffffffffff8111156142c357600080fd5b6020830191508360208285010111156142db57600080fd5b9250929050565b600080600080606085870312156142f857600080fd5b843561430381613e9f565b935060208501359250604085013567ffffffffffffffff81111561432657600080fd5b61433287828801614299565b95989497509550505050565b60008060006040848603121561435357600080fd5b83359250602084013567ffffffffffffffff81111561437157600080fd5b61437d86828701614299565b9497909650939450505050565b60008083601f84011261439c57600080fd5b50813567ffffffffffffffff8111156143b457600080fd5b6020830191508360208260051b85010111156142db57600080fd5b60008060008060008060008060e0898b0312156143eb57600080fd5b606089018a8111156143fc57600080fd5b8998503567ffffffffffffffff8082111561441657600080fd5b6144228c838d01614299565b909950975060808b013591508082111561443b57600080fd5b6144478c838d0161438a565b909750955060a08b013591508082111561446057600080fd5b5061446d8b828c0161438a565b999c989b50969995989497949560c00135949350505050565b60008060008060008060c0878903121561449f57600080fd5b863567ffffffffffffffff808211156144b757600080fd5b6144c38a838b01613ecc565b975060208901359150808211156144d957600080fd5b6144e58a838b01613ecc565b96506144f360408a01613f40565b9550606089013591508082111561450957600080fd5b61425b8a838b0161413e565b60006020828403121561452757600080fd5b8135611dc181613e9f565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60ff8181168382160290811690818114612eb957612eb9614532565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036145dd576145dd614532565b5060010190565b63ffffffff818116838216019080821115612eb957612eb9614532565b600081518084526020808501945080840160005b8381101561464757815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101614615565b509495945050505050565b6020815261466960208201835163ffffffff169052565b60006020830151614682604084018263ffffffff169052565b50604083015163ffffffff8116606084015250606083015162ffffff8116608084015250608083015161ffff811660a08401525060a08301516bffffffffffffffffffffffff811660c08401525060c083015163ffffffff811660e08401525060e08301516101006146fb8185018363ffffffff169052565b84015190506101206147148482018363ffffffff169052565b840151905061014061472d8482018363ffffffff169052565b840151610160848101919091528401516101808085019190915284015190506101a06147708185018373ffffffffffffffffffffffffffffffffffffffff169052565b808501519150506102006101c08181860152614790610220860184614601565b908601519092506101e06147bb8682018373ffffffffffffffffffffffffffffffffffffffff169052565b90950151151593019290925250919050565b600061012063ffffffff808d1684528b6020850152808b166040850152508060608401526147fd8184018a614601565b905082810360808401526148118189614601565b905060ff871660a084015282810360c084015261482e8187613d53565b905067ffffffffffffffff851660e08401528281036101008401526148538185613d53565b9c9b505050505050505050505050565b60006020828403121561487557600080fd5b5035919050565b6bffffffffffffffffffffffff818116838216019080821115612eb957612eb9614532565b80820180821115611daf57611daf614532565b60ff8181168382160190811115611daf57611daf614532565b8051613c6d81613f51565b8051613c6d81613f6e565b8051613c6d81613f8a565b8051613c6d81613fa5565b8051613c6d81613e9f565b600082601f83011261491557600080fd5b81516020614925613eed83613e7b565b82815260059290921b8401810191818101908684111561494457600080fd5b8286015b84811015613f3557805161495b81613e9f565b8352918301918301614948565b8051613c6d81613fca565b60006020828403121561498557600080fd5b815167ffffffffffffffff8082111561499d57600080fd5b9083019061020082860312156149b257600080fd5b6149ba613ddf565b6149c3836148cd565b81526149d1602084016148cd565b60208201526149e2604084016148cd565b60408201526149f3606084016148d8565b6060820152614a04608084016148e3565b6080820152614a1560a084016148ee565b60a0820152614a2660c084016148cd565b60c0820152614a3760e084016148cd565b60e0820152610100614a4a8185016148cd565b90820152610120614a5c8482016148cd565b9082015261014083810151908201526101608084015190820152610180614a848185016148f9565b908201526101a08381015183811115614a9c57600080fd5b614aa888828701614904565b8284015250506101c09150614abe8284016148f9565b828201526101e09150614ad2828401614968565b91810191909152949350505050565b6bffffffffffffffffffffffff828116828216039080821115612eb957612eb9614532565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60006bffffffffffffffffffffffff80841680614b5457614b54614b06565b92169190910492915050565b6bffffffffffffffffffffffff818116838216028082169190828114614b8857614b88614532565b505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600060208284031215614bd157600080fd5b5051919050565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b166040850152816060850152614c1f8285018b614601565b91508382036080850152614c33828a614601565b915060ff881660a085015283820360c0850152614c508288613d53565b90861660e085015283810361010085015290506148538185613d53565b828152604060208201526000612fb26040830184613d53565b60008060408385031215614c9957600080fd5b8251614ca481613fca565b6020939093015192949293505050565b8183823760009101908152919050565b8281526080810160608360208401379392505050565b600082601f830112614ceb57600080fd5b81356020614cfb613eed83613e7b565b82815260059290921b84018101918181019086841115614d1a57600080fd5b8286015b84811015613f355780358352918301918301614d1e565b600082601f830112614d4657600080fd5b81356020614d56613eed83613e7b565b82815260059290921b84018101918181019086841115614d7557600080fd5b8286015b84811015613f3557803567ffffffffffffffff811115614d995760008081fd5b614da78986838b010161413e565b845250918301918301614d79565b600060208284031215614dc757600080fd5b813567ffffffffffffffff80821115614ddf57600080fd5b9083019060c08286031215614df357600080fd5b614dfb613e09565b8235815260208301356020820152604083013582811115614e1b57600080fd5b614e2787828601614cda565b604083015250606083013582811115614e3f57600080fd5b614e4b87828601614cda565b606083015250608083013582811115614e6357600080fd5b614e6f87828601614d35565b60808301525060a083013582811115614e8757600080fd5b614e9387828601614d35565b60a08301525095945050505050565b61ffff818116838216019080821115612eb957612eb9614532565b81810381811115611daf57611daf614532565b8082028115828204841417611daf57611daf614532565b600082614ef657614ef6614b06565b500490565b6bffffffffffffffffffffffff85168152836020820152826040820152608060608201526000614f2e6080830184613d53565b9695505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b828482376000838201600081528351614f84818360208801613d2f565b0195945050505050565b63ffffffff818116838216028082169190828114614b8857614b88614532565b60008060008060008060c08789031215614fc757600080fd5b865195506020870151945060408701519350606087015192506080870151915060a087015190509295509295509295565b60006040828403121561500a57600080fd5b6040516040810181811067ffffffffffffffff8211171561502d5761502d613db0565b604052825161503b81613f51565b81526020928301519281019290925250919050565b600060a0828403121561506257600080fd5b60405160a0810181811067ffffffffffffffff8211171561508557615085613db0565b8060405250825181526020830151602082015260408301516150a681613f51565b604082015260608301516150b981613f51565b6060820152608092830151928101929092525091905056fe307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000813000a",
=======
	ABI: "[{\"inputs\":[{\"internalType\":\"contractAutomationRegistryLogicB2_2\",\"name\":\"logicA\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTrigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTriggerType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"AdminPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newModule\",\"type\":\"address\"}],\"name\":\"ChainSpecificModuleUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"DedupKeyAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepTriggerConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fallbackTo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfigBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"chainSpecificModule\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"reorgProtectionEnabled\",\"type\":\"bool\"}],\"internalType\":\"structAutomationRegistryBase2_2.OnchainConfig\",\"name\":\"onchainConfig\",\"type\":\"tuple\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfigTypeSafe\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"simulatePerformUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"rawReport\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6101206040523480156200001257600080fd5b50604051620054a6380380620054a6833981016040819052620000359162000345565b80816001600160a01b031663ca30e6036040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200009b919062000345565b826001600160a01b031663b10b673c6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000da573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000100919062000345565b836001600160a01b0316636709d0e56040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200013f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000165919062000345565b846001600160a01b0316635425d8ac6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca919062000345565b3380600081620002215760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200025457620002548162000281565b5050506001600160a01b0393841660805291831660a052821660c052811660e0521661010052506200036c565b336001600160a01b03821603620002db5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000218565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200034257600080fd5b50565b6000602082840312156200035857600080fd5b815162000365816200032c565b9392505050565b60805160a05160c05160e051610100516150f8620003ae6000396000818160d6015261016f0152600050506000505060005050600061043301526150f86000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c8063aed2e92911610081578063e3d0e7121161005b578063e3d0e712146102e0578063f2fde38b146102f3578063f75f6b1114610306576100d4565b8063aed2e92914610262578063afcb95d71461028c578063b1dc65a4146102cd576100d4565b806381ff7048116100b257806381ff7048146101bc5780638da5cb5b14610231578063a4c0ed361461024f576100d4565b8063181f5a771461011b578063349e8cca1461016d57806379ba5097146101b4575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e808015610114573d6000f35b3d6000fd5b005b6101576040518060400160405280601881526020017f4175746f6d6174696f6e526567697374727920322e322e30000000000000000081525081565b6040516101649190613d70565b60405180910390f35b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610164565b610119610319565b61020e60155460115463ffffffff780100000000000000000000000000000000000000000000000083048116937c01000000000000000000000000000000000000000000000000000000009093041691565b6040805163ffffffff948516815293909216602084015290820152606001610164565b60005473ffffffffffffffffffffffffffffffffffffffff1661018f565b61011961025d366004613dfe565b61041b565b610275610270366004613e5a565b610637565b604080519215158352602083019190915201610164565b601154601254604080516000815260208101939093527401000000000000000000000000000000000000000090910463ffffffff1690820152606001610164565b6101196102db366004613eeb565b6107af565b6101196102ee3660046141bc565b61137d565b610119610301366004614289565b6113a6565b61011961031436600461448d565b6113ba565b60015473ffffffffffffffffffffffffffffffffffffffff16331461039f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461048a576040517fc8bad78d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b602081146104c4576040517fdfe9309000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006104d28284018461451c565b60008181526004602052604090205490915065010000000000900463ffffffff9081161461052c576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000818152600460205260409020600101546105679085906c0100000000000000000000000090046bffffffffffffffffffffffff16614564565b600082815260046020526040902060010180546bffffffffffffffffffffffff929092166c01000000000000000000000000027fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff9092169190911790556019546105d2908590614589565b6019556040516bffffffffffffffffffffffff8516815273ffffffffffffffffffffffffffffffffffffffff86169082907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a35050505050565b6000806106426123a6565b6012547e01000000000000000000000000000000000000000000000000000000000000900460ff16156106a1576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600085815260046020908152604091829020825160e081018452815460ff811615158252610100810463ffffffff908116838601819052650100000000008304821684880152690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff16606084018190526001909401546bffffffffffffffffffffffff80821660808601526c0100000000000000000000000082041660a0850152780100000000000000000000000000000000000000000000000090041660c08301528451601f890185900485028101850190955287855290936107a09389908990819084018382808284376000920191909152506123e092505050565b9093509150505b935093915050565b60005a60408051610160810182526012546bffffffffffffffffffffffff8116825263ffffffff6c010000000000000000000000008204811660208401527001000000000000000000000000000000008204811693830193909352740100000000000000000000000000000000000000008104909216606082015262ffffff7801000000000000000000000000000000000000000000000000830416608082015261ffff7b0100000000000000000000000000000000000000000000000000000083041660a082015260ff7d0100000000000000000000000000000000000000000000000000000000008304811660c08301527e010000000000000000000000000000000000000000000000000000000000008304811615801560e08401527f01000000000000000000000000000000000000000000000000000000000000009093048116151561010080840191909152601354918216151561012084015273ffffffffffffffffffffffffffffffffffffffff910416610140820152919250610965576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600b602052604090205460ff166109ae576040517f1099ed7500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6011548a35146109ea576040517fdfdcf8e700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60c08101516109fa9060016145cb565b60ff1686141580610a0b5750858414155b15610a42576040517f0244f71a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610a528a8a8a8a8a8a8a8a612609565b6000610a5e8a8a612872565b9050600081604001515167ffffffffffffffff811115610a8057610a80613fa2565b604051908082528060200260200182016040528015610b4457816020015b604080516101e0810182526000610100820181815261012083018290526101408301829052610160830182905261018083018290526101a083018290526101c0830182905282526020808301829052928201819052606082018190526080820181905260a0820181905260c0820181905260e082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181610a9e5790505b5090506000805b836040015151811015610f8f576004600085604001518381518110610b7257610b7261459c565b6020908102919091018101518252818101929092526040908101600020815160e081018352815460ff811615158252610100810463ffffffff90811695830195909552650100000000008104851693820193909352690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff166060830152600101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490911660c08201528351849083908110610c5757610c5761459c565b602002602001015160000181905250610c8c84604001518281518110610c7f57610c7f61459c565b602002602001015161292d565b838281518110610c9e57610c9e61459c565b6020026020010151608001906001811115610cbb57610cbb6145e4565b90816001811115610cce57610cce6145e4565b81525050610d4285848381518110610ce857610ce861459c565b60200260200101516080015186606001518481518110610d0a57610d0a61459c565b60200260200101518760a001518581518110610d2857610d2861459c565b6020026020010151518860000151896020015160016129d8565b838281518110610d5457610d5461459c565b6020026020010151604001906bffffffffffffffffffffffff1690816bffffffffffffffffffffffff1681525050610de184604001518281518110610d9b57610d9b61459c565b602002602001015185608001518381518110610db957610db961459c565b6020026020010151858481518110610dd357610dd361459c565b602002602001015188612a23565b848381518110610df357610df361459c565b6020026020010151602001858481518110610e1057610e1061459c565b602002602001015160e0018281525082151515158152505050828181518110610e3b57610e3b61459c565b60200260200101516020015115610e5e57610e57600183614613565b9150610e63565b610f7d565b610ec9838281518110610e7857610e7861459c565b6020026020010151600001516060015185606001518381518110610e9e57610e9e61459c565b60200260200101518660a001518481518110610ebc57610ebc61459c565b60200260200101516123e0565b848381518110610edb57610edb61459c565b6020026020010151606001858481518110610ef857610ef861459c565b602002602001015160a0018281525082151515158152505050828181518110610f2357610f2361459c565b602002602001015160a0015186610f3a919061462e565b9550610f7d84604001518281518110610f5557610f5561459c565b6020026020010151848381518110610f6f57610f6f61459c565b602002602001015187612c12565b80610f8781614641565b915050610b4b565b508061ffff16600003610fa6575050505050611373565b60c0840151610fb69060016145cb565b610fc59060ff1661044c614679565b616b6c610fd38d6010614679565b5a610fde908961462e565b610fe89190614589565b610ff29190614589565b610ffc9190614589565b9450611b5861100f61ffff8316876146e5565b6110199190614589565b945060008060008060005b87604001515181101561121a578681815181106110435761104361459c565b602002602001015160200151156112085761109f8a88838151811061106a5761106a61459c565b6020026020010151608001518a60a00151848151811061108c5761108c61459c565b6020026020010151518c60c00151612d92565b8782815181106110b1576110b161459c565b602002602001015160c001818152505061110d89896040015183815181106110db576110db61459c565b60200260200101518984815181106110f5576110f561459c565b60200260200101518b600001518c602001518b612db2565b909350915061111c8285614564565b93506111288386614564565b945086818151811061113c5761113c61459c565b6020026020010151606001511515886040015182815181106111605761116061459c565b60200260200101517fad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b84866111959190614564565b8a85815181106111a7576111a761459c565b602002602001015160a001518b86815181106111c5576111c561459c565b602002602001015160c001518d6080015187815181106111e7576111e761459c565b60200260200101516040516111ff94939291906146f9565b60405180910390a35b8061121281614641565b915050611024565b5050336000908152600b6020526040902080548492506002906112529084906201000090046bffffffffffffffffffffffff16614564565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555080601260000160008282829054906101000a90046bffffffffffffffffffffffff166112ac9190614564565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060008f6001600381106112ef576112ef61459c565b602002013560001c9050600060088264ffffffffff16901c9050876060015163ffffffff168163ffffffff16111561136957601280547fffffffffffffffff00000000ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000063ffffffff8416021790555b5050505050505050505b5050505050505050565b61139e8686868680602001905181019061139791906147dc565b86866113ba565b505050505050565b6113ae612ea5565b6113b781612f26565b50565b6113c2612ea5565b601f865111156113fe576040517f25d0209c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8360ff1660000361143b576040517fe77dba5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8451865114158061145a575061145284600361495e565b60ff16865111155b15611491576040517f1d2d1c5800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601254600e546bffffffffffffffffffffffff9091169060005b816bffffffffffffffffffffffff1681101561151357611500600e82815481106114d7576114d761459c565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff16848461301b565b508061150b81614641565b9150506114ab565b5060008060005b836bffffffffffffffffffffffff1681101561161c57600d81815481106115435761154361459c565b600091825260209091200154600e805473ffffffffffffffffffffffffffffffffffffffff9092169450908290811061157e5761157e61459c565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff8681168452600c8352604080852080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001690559116808452600b90925290912080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905591508061161481614641565b91505061151a565b50611629600d6000613c4f565b611635600e6000613c4f565b604080516080810182526000808252602082018190529181018290526060810182905290805b8c51811015611a9e57600c60008e838151811061167a5761167a61459c565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528101919091526040016000205460ff16156116e5576040517f77cea0fa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168d828151811061170f5761170f61459c565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1603611764576040517f815e1d6400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405180604001604052806001151581526020018260ff16815250600c60008f84815181106117955761179561459c565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528181019290925260400160002082518154939092015160ff16610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909316929092171790558b518c908290811061183d5761183d61459c565b60200260200101519150600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036118ad576040517f58a70a0a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff82166000908152600b60209081526040918290208251608081018452905460ff80821615801584526101008304909116938301939093526bffffffffffffffffffffffff6201000082048116948301949094526e01000000000000000000000000000090049092166060830152909350611968576040517f6a7281ad00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001835260ff80821660208086019182526bffffffffffffffffffffffff808b166060880190815273ffffffffffffffffffffffffffffffffffffffff87166000908152600b909352604092839020885181549551948a0151925184166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff939094166201000002929092167fffffffffffff000000000000000000000000000000000000000000000000ffff94909616610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009095169490941717919091169290921791909117905580611a9681614641565b91505061165b565b50508a51611ab49150600d9060208d0190613c6d565b508851611ac890600e9060208c0190613c6d565b50604051806101600160405280856bffffffffffffffffffffffff168152602001886000015163ffffffff168152602001886020015163ffffffff168152602001600063ffffffff168152602001886060015162ffffff168152602001886080015161ffff1681526020018960ff1681526020016012600001601e9054906101000a900460ff16151581526020016012600001601f9054906101000a900460ff161515815260200188610200015115158152602001886101e0015173ffffffffffffffffffffffffffffffffffffffff16815250601260008201518160000160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550602082015181600001600c6101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160106101000a81548163ffffffff021916908363ffffffff16021790555060608201518160000160146101000a81548163ffffffff021916908363ffffffff16021790555060808201518160000160186101000a81548162ffffff021916908362ffffff16021790555060a082015181600001601b6101000a81548161ffff021916908361ffff16021790555060c082015181600001601d6101000a81548160ff021916908360ff16021790555060e082015181600001601e6101000a81548160ff02191690831515021790555061010082015181600001601f6101000a81548160ff0219169083151502179055506101208201518160010160006101000a81548160ff0219169083151502179055506101408201518160010160016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055509050506040518061018001604052808860a001516bffffffffffffffffffffffff16815260200188610180015173ffffffffffffffffffffffffffffffffffffffff168152602001601460010160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff168152602001886040015163ffffffff1681526020018860c0015163ffffffff168152602001601460010160149054906101000a900463ffffffff1663ffffffff168152602001601460010160189054906101000a900463ffffffff1663ffffffff1681526020016014600101601c9054906101000a900463ffffffff1663ffffffff1681526020018860e0015163ffffffff16815260200188610100015163ffffffff16815260200188610120015163ffffffff168152602001886101c0015173ffffffffffffffffffffffffffffffffffffffff16815250601460008201518160000160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550602082015181600001600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060408201518160010160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550606082015181600101600c6101000a81548163ffffffff021916908363ffffffff16021790555060808201518160010160106101000a81548163ffffffff021916908363ffffffff16021790555060a08201518160010160146101000a81548163ffffffff021916908363ffffffff16021790555060c08201518160010160186101000a81548163ffffffff021916908363ffffffff16021790555060e082015181600101601c6101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160020160006101000a81548163ffffffff021916908363ffffffff1602179055506101208201518160020160046101000a81548163ffffffff021916908363ffffffff1602179055506101408201518160020160086101000a81548163ffffffff021916908363ffffffff16021790555061016082015181600201600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505086610140015160178190555086610160015160188190555060006014600101601c9054906101000a900463ffffffff169050876101e0015173ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa158015612169573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061218d9190614987565b601580547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167c010000000000000000000000000000000000000000000000000000000063ffffffff938416021780825560019260189161220891859178010000000000000000000000000000000000000000000000009004166149a0565b92506101000a81548163ffffffff021916908363ffffffff1602179055506000886040516020016122399190614a0e565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529190526015549091506122a290469030907801000000000000000000000000000000000000000000000000900463ffffffff168f8f8f878f8f613223565b60115560005b6122b260096132cd565b8110156122e2576122cf6122c76009836132d7565b6009906132ea565b50806122da81614641565b9150506122a8565b5060005b896101a0015151811015612339576123268a6101a00151828151811061230e5761230e61459c565b6020026020010151600961330c90919063ffffffff16565b508061233181614641565b9150506122e6565b507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0582601154601460010160189054906101000a900463ffffffff168f8f8f878f8f60405161239099989796959493929190614bb2565b60405180910390a1505050505050505050505050565b32156123de576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60125460009081907f0100000000000000000000000000000000000000000000000000000000000000900460ff1615612445576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601280547effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f01000000000000000000000000000000000000000000000000000000000000001790556040517f4585e33b00000000000000000000000000000000000000000000000000000000906124c1908590602401613d70565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009094169390931790925290517f79188d1600000000000000000000000000000000000000000000000000000000815290935073ffffffffffffffffffffffffffffffffffffffff8616906379188d16906125949087908790600401614c48565b60408051808303816000875af11580156125b2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906125d69190614c61565b601280547effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff16905590969095509350505050565b6000878760405161261b929190614c8f565b604051908190038120612632918b90602001614c9f565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201208383019092526000808452908301819052909250906000805b888110156128095760018587836020811061269e5761269e61459c565b6126ab91901a601b6145cb565b8c8c858181106126bd576126bd61459c565b905060200201358b8b868181106126d6576126d661459c565b9050602002013560405160008152602001604052604051612713949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015612735573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff81166000908152600c602090815290849020838501909452925460ff80821615158085526101009092041693830193909352909550935090506127e3576040517f0f4c073700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826020015160080260ff166001901b84019350808061280190614641565b915050612681565b50827e01010101010101010101010101010101010101010101010101010101010101841614612864576040517fc103be2e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505050505050505050505050565b6128ab6040518060c001604052806000815260200160008152602001606081526020016060815260200160608152602001606081525090565b60006128b983850185614d90565b60408101515160608201515191925090811415806128dc57508082608001515114155b806128ec5750808260a001515114155b15612923576040517fb55ac75400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5090505b92915050565b6000818160045b600f8110156129ba577fff0000000000000000000000000000000000000000000000000000000000000082168382602081106129725761297261459c565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916146129a857506000949350505050565b806129b281614641565b915050612934565b5081600f1a60018111156129d0576129d06145e4565b949350505050565b6000806129ea88878b60c0015161332e565b9050600080612a058b8a63ffffffff16858a8a60018b6133ba565b9092509050612a148183614564565b9b9a5050505050505050505050565b600080808085608001516001811115612a3e57612a3e6145e4565b03612a6357612a4f87878787613655565b612a5e57600092509050612c09565b612ada565b600185608001516001811115612a7b57612a7b6145e4565b03612aa8576000612a8d88888761385e565b9250905080612aa25750600092509050612c09565b50612ada565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b83610140015173ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa158015612b2a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612b4e9190614987565b85516040015163ffffffff1611612ba257867fc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd563687604051612b8f9190613d70565b60405180910390a2600092509050612c09565b84604001516bffffffffffffffffffffffff16856000015160a001516bffffffffffffffffffffffff161015612c0257867f377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e0287604051612b8f9190613d70565b6001925090505b94509492505050565b600082608001516001811115612c2a57612c2a6145e4565b03612d095780610140015173ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa158015612c7f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612ca39190614987565b6000848152600460205260409020600101805463ffffffff929092167801000000000000000000000000000000000000000000000000027fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff909216919091179055505050565b600182608001516001811115612d2157612d216145e4565b03612d8d5760e08201805160009081526008602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055915191517fa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f29190a25b505050565b6000612d9f84848461332e565b9050808510156129d05750929392505050565b600080612dcd888760a001518860c0015188888860016133ba565b90925090506000612dde8284614564565b600089815260046020526040902060010180549192508291600c90612e229084906c0100000000000000000000000090046bffffffffffffffffffffffff16614e7d565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915560008a815260046020526040812060010180548594509092612e6b91859116614564565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555050965096945050505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146123de576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610396565b3373ffffffffffffffffffffffffffffffffffffffff821603612fa5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610396565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e01000000000000000000000000000090049091166060820152906132175760008160600151856130b39190614e7d565b905060006130c18583614ea2565b905080836040018181516130d59190614564565b6bffffffffffffffffffffffff169052506130f08582614ecd565b836060018181516131019190614564565b6bffffffffffffffffffffffff90811690915273ffffffffffffffffffffffffffffffffffffffff89166000908152600b602090815260409182902087518154928901519389015160608a015186166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff919096166201000002167fffffffffffff000000000000000000000000000000000000000000000000ffff60ff95909516610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909416939093171792909216179190911790555050505b60400151949350505050565b6000808a8a8a8a8a8a8a8a8a60405160200161324799989796959493929190614f01565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179b9a5050505050505050505050565b6000612927825490565b60006132e38383613adc565b9392505050565b60006132e38373ffffffffffffffffffffffffffffffffffffffff8416613b06565b60006132e38373ffffffffffffffffffffffffffffffffffffffff8416613c00565b60008080856001811115613344576133446145e4565b03613353575062015f90613372565b6001856001811115613367576133676145e4565b03612aa857506201adb05b61338363ffffffff85166014614679565b61338e8460016145cb565b61339d9060ff16611d4c614679565b6133a79083614589565b6133b19190614589565b95945050505050565b60008060008960a0015161ffff16876133d39190614679565b90508380156133e15750803a105b156133e957503a5b6000841561347a578a610140015173ffffffffffffffffffffffffffffffffffffffff166349948e0e6000366040518363ffffffff1660e01b8152600401613432929190614f96565b602060405180830381865afa15801561344f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906134739190614987565b9050613536565b6101408b01516016546040517f1254414000000000000000000000000000000000000000000000000000000000815264010000000090910463ffffffff16600482015273ffffffffffffffffffffffffffffffffffffffff90911690631254414090602401602060405180830381865afa1580156134fc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906135209190614987565b8b60a0015161ffff166135339190614679565b90505b61354461ffff8716826146e5565b9050600087826135548c8e614589565b61355e9086614679565b6135689190614589565b61357a90670de0b6b3a7640000614679565b61358491906146e5565b905060008c6040015163ffffffff1664e8d4a510006135a39190614679565b898e6020015163ffffffff16858f886135bc9190614679565b6135c69190614589565b6135d490633b9aca00614679565b6135de9190614679565b6135e891906146e5565b6135f29190614589565b90506b033b2e3c9fd0803ce800000061360b8284614589565b1115613643576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b6000808480602001905181019061366c9190614fe3565b845160c00151815191925063ffffffff908116911610156136c957857f405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8866040516136b79190613d70565b60405180910390a260009150506129d0565b610140830151610120840151801561378957506020820151158015906137895750602082015182516040517f85df51fd00000000000000000000000000000000000000000000000000000000815263ffffffff909116600482015273ffffffffffffffffffffffffffffffffffffffff8316906385df51fd90602401602060405180830381865afa158015613762573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906137869190614987565b14155b8061380957508073ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa1580156137da573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906137fe9190614987565b825163ffffffff1610155b1561385157867f6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc3018760405161383e9190613d70565b60405180910390a26000925050506129d0565b5060019695505050505050565b600080600084806020019051810190613877919061503b565b90506000868260000151836020015184604001516040516020016138d994939291909384526020840192909252604083015260e01b7fffffffff0000000000000000000000000000000000000000000000000000000016606082015260640190565b6040516020818303038152906040528051906020012090506000856101400151905085610120015180156139b857506080830151158015906139b85750608083015160608401516040517f85df51fd00000000000000000000000000000000000000000000000000000000815263ffffffff909116600482015273ffffffffffffffffffffffffffffffffffffffff8316906385df51fd90602401602060405180830381865afa158015613991573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906139b59190614987565b14155b80613a3b57508073ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa158015613a09573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613a2d9190614987565b836060015163ffffffff1610155b15613a8657877f6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc30188604051613a709190613d70565b60405180910390a2506000935091506107a79050565b60008281526008602052604090205460ff1615613acd57877f405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e888604051613a709190613d70565b50600197909650945050505050565b6000826000018281548110613af357613af361459c565b9060005260206000200154905092915050565b60008181526001830160205260408120548015613bef576000613b2a60018361462e565b8554909150600090613b3e9060019061462e565b9050818114613ba3576000866000018281548110613b5e57613b5e61459c565b9060005260206000200154905080876000018481548110613b8157613b8161459c565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080613bb457613bb46150bc565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050612927565b6000915050612927565b5092915050565b6000818152600183016020526040812054613c4757508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155612927565b506000612927565b50805460008255906000526020600020908101906113b79190613cf7565b828054828255906000526020600020908101928215613ce7579160200282015b82811115613ce757825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190613c8d565b50613cf3929150613cf7565b5090565b5b80821115613cf35760008155600101613cf8565b6000815180845260005b81811015613d3257602081850181015186830182015201613d16565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815260006132e36020830184613d0c565b73ffffffffffffffffffffffffffffffffffffffff811681146113b757600080fd5b8035613db081613d83565b919050565b60008083601f840112613dc757600080fd5b50813567ffffffffffffffff811115613ddf57600080fd5b602083019150836020828501011115613df757600080fd5b9250929050565b60008060008060608587031215613e1457600080fd5b8435613e1f81613d83565b935060208501359250604085013567ffffffffffffffff811115613e4257600080fd5b613e4e87828801613db5565b95989497509550505050565b600080600060408486031215613e6f57600080fd5b83359250602084013567ffffffffffffffff811115613e8d57600080fd5b613e9986828701613db5565b9497909650939450505050565b60008083601f840112613eb857600080fd5b50813567ffffffffffffffff811115613ed057600080fd5b6020830191508360208260051b8501011115613df757600080fd5b60008060008060008060008060e0898b031215613f0757600080fd5b606089018a811115613f1857600080fd5b8998503567ffffffffffffffff80821115613f3257600080fd5b613f3e8c838d01613db5565b909950975060808b0135915080821115613f5757600080fd5b613f638c838d01613ea6565b909750955060a08b0135915080821115613f7c57600080fd5b50613f898b828c01613ea6565b999c989b50969995989497949560c00135949350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610220810167ffffffffffffffff81118282101715613ff557613ff5613fa2565b60405290565b60405160c0810167ffffffffffffffff81118282101715613ff557613ff5613fa2565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561406557614065613fa2565b604052919050565b600067ffffffffffffffff82111561408757614087613fa2565b5060051b60200190565b600082601f8301126140a257600080fd5b813560206140b76140b28361406d565b61401e565b82815260059290921b840181019181810190868411156140d657600080fd5b8286015b848110156140fa5780356140ed81613d83565b83529183019183016140da565b509695505050505050565b803560ff81168114613db057600080fd5b600082601f83011261412757600080fd5b813567ffffffffffffffff81111561414157614141613fa2565b61417260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8401160161401e565b81815284602083860101111561418757600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff81168114613db057600080fd5b60008060008060008060c087890312156141d557600080fd5b863567ffffffffffffffff808211156141ed57600080fd5b6141f98a838b01614091565b9750602089013591508082111561420f57600080fd5b61421b8a838b01614091565b965061422960408a01614105565b9550606089013591508082111561423f57600080fd5b61424b8a838b01614116565b945061425960808a016141a4565b935060a089013591508082111561426f57600080fd5b5061427c89828a01614116565b9150509295509295509295565b60006020828403121561429b57600080fd5b81356132e381613d83565b63ffffffff811681146113b757600080fd5b8035613db0816142a6565b62ffffff811681146113b757600080fd5b8035613db0816142c3565b61ffff811681146113b757600080fd5b8035613db0816142df565b6bffffffffffffffffffffffff811681146113b757600080fd5b8035613db0816142fa565b80151581146113b757600080fd5b8035613db08161431f565b6000610220828403121561434b57600080fd5b614353613fd1565b905061435e826142b8565b815261436c602083016142b8565b602082015261437d604083016142b8565b604082015261438e606083016142d4565b606082015261439f608083016142ef565b60808201526143b060a08301614314565b60a08201526143c160c083016142b8565b60c08201526143d260e083016142b8565b60e08201526101006143e58184016142b8565b908201526101206143f78382016142b8565b908201526101408281013590820152610160808301359082015261018061441f818401613da5565b908201526101a08281013567ffffffffffffffff81111561443f57600080fd5b61444b85828601614091565b8284015250506101c061445f818401613da5565b908201526101e0614471838201613da5565b9082015261020061448383820161432d565b9082015292915050565b60008060008060008060c087890312156144a657600080fd5b863567ffffffffffffffff808211156144be57600080fd5b6144ca8a838b01614091565b975060208901359150808211156144e057600080fd5b6144ec8a838b01614091565b96506144fa60408a01614105565b9550606089013591508082111561451057600080fd5b61424b8a838b01614338565b60006020828403121561452e57600080fd5b5035919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6bffffffffffffffffffffffff818116838216019080821115613bf957613bf9614535565b8082018082111561292757612927614535565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60ff818116838216019081111561292757612927614535565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b61ffff818116838216019080821115613bf957613bf9614535565b8181038181111561292757612927614535565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361467257614672614535565b5060010190565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156146b1576146b1614535565b500290565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000826146f4576146f46146b6565b500490565b6bffffffffffffffffffffffff8516815283602082015282604082015260806060820152600061472c6080830184613d0c565b9695505050505050565b8051613db0816142a6565b8051613db0816142c3565b8051613db0816142df565b8051613db0816142fa565b8051613db081613d83565b600082601f83011261477e57600080fd5b8151602061478e6140b28361406d565b82815260059290921b840181019181810190868411156147ad57600080fd5b8286015b848110156140fa5780516147c481613d83565b83529183019183016147b1565b8051613db08161431f565b6000602082840312156147ee57600080fd5b815167ffffffffffffffff8082111561480657600080fd5b90830190610220828603121561481b57600080fd5b614823613fd1565b61482c83614736565b815261483a60208401614736565b602082015261484b60408401614736565b604082015261485c60608401614741565b606082015261486d6080840161474c565b608082015261487e60a08401614757565b60a082015261488f60c08401614736565b60c08201526148a060e08401614736565b60e08201526101006148b3818501614736565b908201526101206148c5848201614736565b90820152610140838101519082015261016080840151908201526101806148ed818501614762565b908201526101a0838101518381111561490557600080fd5b6149118882870161476d565b8284015250506101c09150614927828401614762565b828201526101e0915061493b828401614762565b82820152610200915061494f8284016147d1565b91810191909152949350505050565b600060ff821660ff84168160ff048111821515161561497f5761497f614535565b029392505050565b60006020828403121561499957600080fd5b5051919050565b63ffffffff818116838216019080821115613bf957613bf9614535565b600081518084526020808501945080840160005b83811015614a0357815173ffffffffffffffffffffffffffffffffffffffff16875295820195908201906001016149d1565b509495945050505050565b60208152614a2560208201835163ffffffff169052565b60006020830151614a3e604084018263ffffffff169052565b50604083015163ffffffff8116606084015250606083015162ffffff8116608084015250608083015161ffff811660a08401525060a08301516bffffffffffffffffffffffff811660c08401525060c083015163ffffffff811660e08401525060e0830151610100614ab78185018363ffffffff169052565b8401519050610120614ad08482018363ffffffff169052565b8401519050610140614ae98482018363ffffffff169052565b840151610160848101919091528401516101808085019190915284015190506101a0614b2c8185018373ffffffffffffffffffffffffffffffffffffffff169052565b808501519150506102206101c08181860152614b4c6102408601846149bd565b908601519092506101e0614b778682018373ffffffffffffffffffffffffffffffffffffffff169052565b8601519050610200614ba08682018373ffffffffffffffffffffffffffffffffffffffff169052565b90950151151593019290925250919050565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152614be28184018a6149bd565b90508281036080840152614bf681896149bd565b905060ff871660a084015282810360c0840152614c138187613d0c565b905067ffffffffffffffff851660e0840152828103610100840152614c388185613d0c565b9c9b505050505050505050505050565b8281526040602082015260006129d06040830184613d0c565b60008060408385031215614c7457600080fd5b8251614c7f8161431f565b6020939093015192949293505050565b8183823760009101908152919050565b8281526080810160608360208401379392505050565b600082601f830112614cc657600080fd5b81356020614cd66140b28361406d565b82815260059290921b84018101918181019086841115614cf557600080fd5b8286015b848110156140fa5780358352918301918301614cf9565b600082601f830112614d2157600080fd5b81356020614d316140b28361406d565b82815260059290921b84018101918181019086841115614d5057600080fd5b8286015b848110156140fa57803567ffffffffffffffff811115614d745760008081fd5b614d828986838b0101614116565b845250918301918301614d54565b600060208284031215614da257600080fd5b813567ffffffffffffffff80821115614dba57600080fd5b9083019060c08286031215614dce57600080fd5b614dd6613ffb565b8235815260208301356020820152604083013582811115614df657600080fd5b614e0287828601614cb5565b604083015250606083013582811115614e1a57600080fd5b614e2687828601614cb5565b606083015250608083013582811115614e3e57600080fd5b614e4a87828601614d10565b60808301525060a083013582811115614e6257600080fd5b614e6e87828601614d10565b60a08301525095945050505050565b6bffffffffffffffffffffffff828116828216039080821115613bf957613bf9614535565b60006bffffffffffffffffffffffff80841680614ec157614ec16146b6565b92169190910492915050565b60006bffffffffffffffffffffffff80831681851681830481118215151615614ef857614ef8614535565b02949350505050565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b166040850152816060850152614f488285018b6149bd565b91508382036080850152614f5c828a6149bd565b915060ff881660a085015283820360c0850152614f798288613d0c565b90861660e08501528381036101008501529050614c388185613d0c565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b600060408284031215614ff557600080fd5b6040516040810181811067ffffffffffffffff8211171561501857615018613fa2565b6040528251615026816142a6565b81526020928301519281019290925250919050565b600060a0828403121561504d57600080fd5b60405160a0810181811067ffffffffffffffff8211171561507057615070613fa2565b806040525082518152602083015160208201526040830151615091816142a6565b604082015260608301516150a4816142a6565b60608201526080928301519281019290925250919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000810000a",
>>>>>>> ab4212059a (implement modules)
=======
	ABI: "[{\"inputs\":[{\"internalType\":\"contractAutomationRegistryLogicB2_2\",\"name\":\"logicA\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTrigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTriggerType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"AdminPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newModule\",\"type\":\"address\"}],\"name\":\"ChainSpecificModuleUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"DedupKeyAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepTriggerConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fallbackTo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfigBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"},{\"internalType\":\"contractIChainModule\",\"name\":\"chainModule\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"reorgProtectionEnabled\",\"type\":\"bool\"}],\"internalType\":\"structAutomationRegistryBase2_2.OnchainConfig\",\"name\":\"onchainConfig\",\"type\":\"tuple\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfigTypeSafe\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"simulatePerformUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"rawReport\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	Bin: "0x6101206040523480156200001257600080fd5b506040516200532b3803806200532b833981016040819052620000359162000345565b80816001600160a01b031663ca30e6036040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200009b919062000345565b826001600160a01b031663b10b673c6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000da573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000100919062000345565b836001600160a01b0316636709d0e56040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200013f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000165919062000345565b846001600160a01b0316635425d8ac6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca919062000345565b3380600081620002215760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200025457620002548162000281565b5050506001600160a01b0393841660805291831660a052821660c052811660e0521661010052506200036c565b336001600160a01b03821603620002db5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000218565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200034257600080fd5b50565b6000602082840312156200035857600080fd5b815162000365816200032c565b9392505050565b60805160a05160c05160e05161010051614f7d620003ae6000396000818160d6015261016f015260005050600050506000505060006104330152614f7d6000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c8063aed2e92911610081578063e3d0e7121161005b578063e3d0e712146102e0578063f2fde38b146102f3578063f75f6b1114610306576100d4565b8063aed2e92914610262578063afcb95d71461028c578063b1dc65a4146102cd576100d4565b806381ff7048116100b257806381ff7048146101bc5780638da5cb5b14610231578063a4c0ed361461024f576100d4565b8063181f5a771461011b578063349e8cca1461016d57806379ba5097146101b4575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e808015610114573d6000f35b3d6000fd5b005b6101576040518060400160405280601881526020017f4175746f6d6174696f6e526567697374727920322e322e30000000000000000081525081565b6040516101649190613c0d565b60405180910390f35b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610164565b610119610319565b61020e60155460115463ffffffff780100000000000000000000000000000000000000000000000083048116937c01000000000000000000000000000000000000000000000000000000009093041691565b6040805163ffffffff948516815293909216602084015290820152606001610164565b60005473ffffffffffffffffffffffffffffffffffffffff1661018f565b61011961025d366004613c9b565b61041b565b610275610270366004613cf7565b610637565b604080519215158352602083019190915201610164565b601154601254604080516000815260208101939093527401000000000000000000000000000000000000000090910463ffffffff1690820152606001610164565b6101196102db366004613e9d565b6107ad565b6101196102ee36600461409c565b6113f5565b610119610301366004614169565b61141e565b61011961031436600461436d565b611432565b60015473ffffffffffffffffffffffffffffffffffffffff16331461039f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461048a576040517fc8bad78d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b602081146104c4576040517fdfe9309000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006104d2828401846143fc565b60008181526004602052604090205490915065010000000000900463ffffffff9081161461052c576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000818152600460205260409020600101546105679085906c0100000000000000000000000090046bffffffffffffffffffffffff16614444565b600082815260046020526040902060010180546bffffffffffffffffffffffff929092166c01000000000000000000000000027fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff9092169190911790556019546105d2908590614469565b6019556040516bffffffffffffffffffffffff8516815273ffffffffffffffffffffffffffffffffffffffff86169082907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a35050505050565b60008061064261241e565b6012547e01000000000000000000000000000000000000000000000000000000000000900460ff16156106a1576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600085815260046020908152604091829020825160e081018452815460ff811615158252610100810463ffffffff908116838601819052650100000000008304821684880152690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff16606084018190526001909401546bffffffffffffffffffffffff80821660808601526c0100000000000000000000000082041660a0850152780100000000000000000000000000000000000000000000000090041660c08301528451601f890185900485028101850190955287855290936107a093899089908190840183828082843760009201919091525061245892505050565b9097909650945050505050565b60005a60408051610160810182526012546bffffffffffffffffffffffff8116825263ffffffff6c010000000000000000000000008204811660208401527001000000000000000000000000000000008204811693830193909352740100000000000000000000000000000000000000008104909216606082015262ffffff7801000000000000000000000000000000000000000000000000830416608082015261ffff7b0100000000000000000000000000000000000000000000000000000083041660a082015260ff7d0100000000000000000000000000000000000000000000000000000000008304811660c08301527e010000000000000000000000000000000000000000000000000000000000008304811615801560e08401527f01000000000000000000000000000000000000000000000000000000000000009093048116151561010080840191909152601354918216151561012084015273ffffffffffffffffffffffffffffffffffffffff910416610140820152919250610963576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600b602052604090205460ff166109ac576040517f1099ed7500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6011548935146109e8576040517fdfdcf8e700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60c08101516109f89060016144ab565b60ff1685141580610a0a575083518514155b15610a41576040517f0244f71a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610a5089898989898989612681565b6000610a5c89896128ea565b9050600081604001515167ffffffffffffffff811115610a7e57610a7e613d43565b604051908082528060200260200182016040528015610b4257816020015b604080516101e0810182526000610100820181815261012083018290526101408301829052610160830182905261018083018290526101a083018290526101c0830182905282526020808301829052928201819052606082018190526080820181905260a0820181905260c0820181905260e082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181610a9c5790505b50905060008084610140015173ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa158015610b98573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bbc91906144c4565b905060005b846040015151811015611006576004600086604001518381518110610be857610be861447c565b6020908102919091018101518252818101929092526040908101600020815160e081018352815460ff811615158252610100810463ffffffff90811695830195909552650100000000008104851693820193909352690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff166060830152600101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490911660c08201528451859083908110610ccd57610ccd61447c565b602002602001015160000181905250610d0285604001518281518110610cf557610cf561447c565b60200260200101516129a5565b848281518110610d1457610d1461447c565b6020026020010151608001906001811115610d3157610d316144dd565b90816001811115610d4457610d446144dd565b81525050610db886858381518110610d5e57610d5e61447c565b60200260200101516080015187606001518481518110610d8057610d8061447c565b60200260200101518860a001518581518110610d9e57610d9e61447c565b60200260200101515189600001518a602001516001612a50565b848281518110610dca57610dca61447c565b6020026020010151604001906bffffffffffffffffffffffff1690816bffffffffffffffffffffffff1681525050610e5885604001518281518110610e1157610e1161447c565b60200260200101518387608001518481518110610e3057610e3061447c565b6020026020010151878581518110610e4a57610e4a61447c565b60200260200101518a612a9b565b858381518110610e6a57610e6a61447c565b6020026020010151602001868481518110610e8757610e8761447c565b602002602001015160e0018281525082151515158152505050838181518110610eb257610eb261447c565b60200260200101516020015115610ed557610ece60018461450c565b9250610eda565b610ff4565b610f40848281518110610eef57610eef61447c565b6020026020010151600001516060015186606001518381518110610f1557610f1561447c565b60200260200101518760a001518481518110610f3357610f3361447c565b6020026020010151612458565b858381518110610f5257610f5261447c565b6020026020010151606001868481518110610f6f57610f6f61447c565b602002602001015160a0018281525082151515158152505050838181518110610f9a57610f9a61447c565b602002602001015160a0015187610fb19190614527565b9650610ff485604001518281518110610fcc57610fcc61447c565b602002602001015183868481518110610fe757610fe761447c565b6020026020010151612c1a565b80610ffe8161453a565b915050610bc1565b508161ffff1660000361101e575050505050506113ec565b60c085015161102e9060016144ab565b61103d9060ff16610546614572565b61762a61104b8d6010614572565b5a611056908a614527565b6110609190614469565b61106a9190614469565b6110749190614469565b9550611d4c61108761ffff8416886145de565b6110919190614469565b955060008060008060005b886040015151811015611292578781815181106110bb576110bb61447c565b60200260200101516020015115611280576111178b8983815181106110e2576110e261447c565b6020026020010151608001518b60a0015184815181106111045761110461447c565b6020026020010151518d60c00151612d1f565b8882815181106111295761112961447c565b602002602001015160c00181815250506111858a8a6040015183815181106111535761115361447c565b60200260200101518a848151811061116d5761116d61447c565b60200260200101518c600001518d602001518c612d3f565b90935091506111948285614444565b93506111a08386614444565b94508781815181106111b4576111b461447c565b6020026020010151606001511515896040015182815181106111d8576111d861447c565b60200260200101517fad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b848661120d9190614444565b8b858151811061121f5761121f61447c565b602002602001015160a001518c868151811061123d5761123d61447c565b602002602001015160c001518e60800151878151811061125f5761125f61447c565b602002602001015160405161127794939291906145f2565b60405180910390a35b8061128a8161453a565b91505061109c565b5050336000908152600b6020526040902080548492506002906112ca9084906201000090046bffffffffffffffffffffffff16614444565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555080601260000160008282829054906101000a90046bffffffffffffffffffffffff166113249190614444565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060008f6001600381106113675761136761447c565b602002013560001c9050600060088264ffffffffff16901c9050886060015163ffffffff168163ffffffff1611156113e157601280547fffffffffffffffff00000000ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000063ffffffff8416021790555b505050505050505050505b50505050505050565b6114168686868680602001905181019061140f91906146d5565b8686611432565b505050505050565b611426612e32565b61142f81612eb3565b50565b61143a612e32565b601f86511115611476576040517f25d0209c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8360ff166000036114b3576040517fe77dba5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845186511415806114d257506114ca846003614857565b60ff16865111155b15611509576040517f1d2d1c5800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601254600e546bffffffffffffffffffffffff9091169060005b816bffffffffffffffffffffffff1681101561158b57611578600e828154811061154f5761154f61447c565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff168484612fa8565b50806115838161453a565b915050611523565b5060008060005b836bffffffffffffffffffffffff1681101561169457600d81815481106115bb576115bb61447c565b600091825260209091200154600e805473ffffffffffffffffffffffffffffffffffffffff909216945090829081106115f6576115f661447c565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff8681168452600c8352604080852080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001690559116808452600b90925290912080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905591508061168c8161453a565b915050611592565b506116a1600d6000613aec565b6116ad600e6000613aec565b604080516080810182526000808252602082018190529181018290526060810182905290805b8c51811015611b1657600c60008e83815181106116f2576116f261447c565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528101919091526040016000205460ff161561175d576040517f77cea0fa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168d82815181106117875761178761447c565b602002602001015173ffffffffffffffffffffffffffffffffffffffff16036117dc576040517f815e1d6400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405180604001604052806001151581526020018260ff16815250600c60008f848151811061180d5761180d61447c565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528181019290925260400160002082518154939092015160ff16610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909316929092171790558b518c90829081106118b5576118b561447c565b60200260200101519150600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603611925576040517f58a70a0a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff82166000908152600b60209081526040918290208251608081018452905460ff80821615801584526101008304909116938301939093526bffffffffffffffffffffffff6201000082048116948301949094526e010000000000000000000000000000900490921660608301529093506119e0576040517f6a7281ad00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001835260ff80821660208086019182526bffffffffffffffffffffffff808b166060880190815273ffffffffffffffffffffffffffffffffffffffff87166000908152600b909352604092839020885181549551948a0151925184166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff939094166201000002929092167fffffffffffff000000000000000000000000000000000000000000000000ffff94909616610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009095169490941717919091169290921791909117905580611b0e8161453a565b9150506116d3565b50508a51611b2c9150600d9060208d0190613b0a565b508851611b4090600e9060208c0190613b0a565b50604051806101600160405280856bffffffffffffffffffffffff168152602001886000015163ffffffff168152602001886020015163ffffffff168152602001600063ffffffff168152602001886060015162ffffff168152602001886080015161ffff1681526020018960ff1681526020016012600001601e9054906101000a900460ff16151581526020016012600001601f9054906101000a900460ff161515815260200188610200015115158152602001886101e0015173ffffffffffffffffffffffffffffffffffffffff16815250601260008201518160000160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550602082015181600001600c6101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160106101000a81548163ffffffff021916908363ffffffff16021790555060608201518160000160146101000a81548163ffffffff021916908363ffffffff16021790555060808201518160000160186101000a81548162ffffff021916908362ffffff16021790555060a082015181600001601b6101000a81548161ffff021916908361ffff16021790555060c082015181600001601d6101000a81548160ff021916908360ff16021790555060e082015181600001601e6101000a81548160ff02191690831515021790555061010082015181600001601f6101000a81548160ff0219169083151502179055506101208201518160010160006101000a81548160ff0219169083151502179055506101408201518160010160016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055509050506040518061018001604052808860a001516bffffffffffffffffffffffff16815260200188610180015173ffffffffffffffffffffffffffffffffffffffff168152602001601460010160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff168152602001886040015163ffffffff1681526020018860c0015163ffffffff168152602001601460010160149054906101000a900463ffffffff1663ffffffff168152602001601460010160189054906101000a900463ffffffff1663ffffffff1681526020016014600101601c9054906101000a900463ffffffff1663ffffffff1681526020018860e0015163ffffffff16815260200188610100015163ffffffff16815260200188610120015163ffffffff168152602001886101c0015173ffffffffffffffffffffffffffffffffffffffff16815250601460008201518160000160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550602082015181600001600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060408201518160010160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550606082015181600101600c6101000a81548163ffffffff021916908363ffffffff16021790555060808201518160010160106101000a81548163ffffffff021916908363ffffffff16021790555060a08201518160010160146101000a81548163ffffffff021916908363ffffffff16021790555060c08201518160010160186101000a81548163ffffffff021916908363ffffffff16021790555060e082015181600101601c6101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160020160006101000a81548163ffffffff021916908363ffffffff1602179055506101208201518160020160046101000a81548163ffffffff021916908363ffffffff1602179055506101408201518160020160086101000a81548163ffffffff021916908363ffffffff16021790555061016082015181600201600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505086610140015160178190555086610160015160188190555060006014600101601c9054906101000a900463ffffffff169050876101e0015173ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa1580156121e1573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061220591906144c4565b601580547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167c010000000000000000000000000000000000000000000000000000000063ffffffff93841602178082556001926018916122809185917801000000000000000000000000000000000000000000000000900416614880565b92506101000a81548163ffffffff021916908363ffffffff1602179055506000886040516020016122b191906148ee565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291905260155490915061231a90469030907801000000000000000000000000000000000000000000000000900463ffffffff168f8f8f878f8f6131b0565b60115560005b61232a600961325a565b81101561235a5761234761233f600983613264565b600990613277565b50806123528161453a565b915050612320565b5060005b896101a00151518110156123b15761239e8a6101a0015182815181106123865761238661447c565b6020026020010151600961329990919063ffffffff16565b50806123a98161453a565b91505061235e565b507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0582601154601460010160189054906101000a900463ffffffff168f8f8f878f8f60405161240899989796959493929190614a92565b60405180910390a1505050505050505050505050565b3215612456576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60125460009081907f0100000000000000000000000000000000000000000000000000000000000000900460ff16156124bd576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601280547effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f01000000000000000000000000000000000000000000000000000000000000001790556040517f4585e33b0000000000000000000000000000000000000000000000000000000090612539908590602401613c0d565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009094169390931790925290517f79188d1600000000000000000000000000000000000000000000000000000000815290935073ffffffffffffffffffffffffffffffffffffffff8616906379188d169061260c9087908790600401614b28565b60408051808303816000875af115801561262a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061264e9190614b41565b601280547effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff16905590969095509350505050565b60008686604051612693929190614b6f565b6040519081900381206126aa918a90602001614b7f565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201208383019092526000808452908301819052909250906000805b87811015612882576001858783602081106127165761271661447c565b61272391901a601b6144ab565b8b8b858181106127355761273561447c565b905060200201358a858151811061274e5761274e61447c565b60200260200101516040516000815260200160405260405161278c949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa1580156127ae573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff81166000908152600c602090815290849020838501909452925460ff808216151580855261010090920416938301939093529095509350905061285c576040517f0f4c073700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826020015160080260ff166001901b84019350808061287a9061453a565b9150506126f9565b50827e010101010101010101010101010101010101010101010101010101010101018416146128dd576040517fc103be2e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050505050505050505050565b6129236040518060c001604052806000815260200160008152602001606081526020016060815260200160608152602001606081525090565b600061293183850185614c15565b604081015151606082015151919250908114158061295457508082608001515114155b806129645750808260a001515114155b1561299b576040517fb55ac75400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5090505b92915050565b6000818160045b600f811015612a32577fff0000000000000000000000000000000000000000000000000000000000000082168382602081106129ea576129ea61447c565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614612a2057506000949350505050565b80612a2a8161453a565b9150506129ac565b5081600f1a6001811115612a4857612a486144dd565b949350505050565b600080612a6288878b60c001516132bb565b9050600080612a7d8b8a63ffffffff16858a8a60018b613347565b9092509050612a8c8183614444565b9b9a5050505050505050505050565b600080808085608001516001811115612ab657612ab66144dd565b03612adc57612ac888888888886135e2565b612ad757600092509050612c10565b612b54565b600185608001516001811115612af457612af46144dd565b03612b22576000612b078989898861376b565b9250905080612b1c5750600092509050612c10565b50612b54565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b84516040015163ffffffff168710612ba957877fc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd563687604051612b969190613c0d565b60405180910390a2600092509050612c10565b84604001516bffffffffffffffffffffffff16856000015160a001516bffffffffffffffffffffffff161015612c0957877f377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e0287604051612b969190613c0d565b6001925090505b9550959350505050565b600081608001516001811115612c3257612c326144dd565b03612c9657600083815260046020526040902060010180547fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff16780100000000000000000000000000000000000000000000000063ffffffff851602179055505050565b600181608001516001811115612cae57612cae6144dd565b03612d1a5760e08101805160009081526008602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055915191517fa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f29190a25b505050565b6000612d2c8484846132bb565b905080851015612a485750929392505050565b600080612d5a888760a001518860c001518888886001613347565b90925090506000612d6b8284614444565b600089815260046020526040902060010180549192508291600c90612daf9084906c0100000000000000000000000090046bffffffffffffffffffffffff16614d02565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915560008a815260046020526040812060010180548594509092612df891859116614444565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555050965096945050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314612456576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610396565b3373ffffffffffffffffffffffffffffffffffffffff821603612f32576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610396565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e01000000000000000000000000000090049091166060820152906131a45760008160600151856130409190614d02565b9050600061304e8583614d27565b905080836040018181516130629190614444565b6bffffffffffffffffffffffff1690525061307d8582614d52565b8360600181815161308e9190614444565b6bffffffffffffffffffffffff90811690915273ffffffffffffffffffffffffffffffffffffffff89166000908152600b602090815260409182902087518154928901519389015160608a015186166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff919096166201000002167fffffffffffff000000000000000000000000000000000000000000000000ffff60ff95909516610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909416939093171792909216179190911790555050505b60400151949350505050565b6000808a8a8a8a8a8a8a8a8a6040516020016131d499989796959493929190614d86565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179b9a5050505050505050505050565b600061299f825490565b60006132708383613979565b9392505050565b60006132708373ffffffffffffffffffffffffffffffffffffffff84166139a3565b60006132708373ffffffffffffffffffffffffffffffffffffffff8416613a9d565b600080808560018111156132d1576132d16144dd565b036132e0575062015f906132ff565b60018560018111156132f4576132f46144dd565b03612b2257506201adb05b61331063ffffffff85166014614572565b61331b8460016144ab565b61332a9060ff16611d4c614572565b6133349083614469565b61333e9190614469565b95945050505050565b60008060008960a0015161ffff16876133609190614572565b905083801561336e5750803a105b1561337657503a5b60008415613407578a610140015173ffffffffffffffffffffffffffffffffffffffff166349948e0e6000366040518363ffffffff1660e01b81526004016133bf929190614e1b565b602060405180830381865afa1580156133dc573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061340091906144c4565b90506134c3565b6101408b01516016546040517f1254414000000000000000000000000000000000000000000000000000000000815264010000000090910463ffffffff16600482015273ffffffffffffffffffffffffffffffffffffffff90911690631254414090602401602060405180830381865afa158015613489573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906134ad91906144c4565b8b60a0015161ffff166134c09190614572565b90505b6134d161ffff8716826145de565b9050600087826134e18c8e614469565b6134eb9086614572565b6134f59190614469565b61350790670de0b6b3a7640000614572565b61351191906145de565b905060008c6040015163ffffffff1664e8d4a510006135309190614572565b898e6020015163ffffffff16858f886135499190614572565b6135539190614469565b61356190633b9aca00614572565b61356b9190614572565b61357591906145de565b61357f9190614469565b90506b033b2e3c9fd0803ce80000006135988284614469565b11156135d0576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b600080848060200190518101906135f99190614e68565b845160c00151815191925063ffffffff9081169116101561365657867f405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8866040516136449190613c0d565b60405180910390a2600091505061333e565b826101200151801561371757506020810151158015906137175750602081015161014084015182516040517f85df51fd00000000000000000000000000000000000000000000000000000000815263ffffffff909116600482015273ffffffffffffffffffffffffffffffffffffffff909116906385df51fd90602401602060405180830381865afa1580156136f0573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061371491906144c4565b14155b806137295750805163ffffffff168611155b1561375e57867f6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301866040516136449190613c0d565b5060019695505050505050565b6000806000848060200190518101906137849190614ec0565b90506000878260000151836020015184604001516040516020016137e694939291909384526020840192909252604083015260e01b7fffffffff0000000000000000000000000000000000000000000000000000000016606082015260640190565b60405160208183030381529060405280519060200120905084610120015180156138c257506080820151158015906138c25750608082015161014086015160608401516040517f85df51fd00000000000000000000000000000000000000000000000000000000815263ffffffff909116600482015273ffffffffffffffffffffffffffffffffffffffff909116906385df51fd90602401602060405180830381865afa15801561389b573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906138bf91906144c4565b14155b806138d7575086826060015163ffffffff1610155b1561392157877f6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc3018760405161390c9190613c0d565b60405180910390a26000935091506139709050565b60008181526008602052604090205460ff161561396857877f405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e88760405161390c9190613c0d565b600193509150505b94509492505050565b60008260000182815481106139905761399061447c565b9060005260206000200154905092915050565b60008181526001830160205260408120548015613a8c5760006139c7600183614527565b85549091506000906139db90600190614527565b9050818114613a405760008660000182815481106139fb576139fb61447c565b9060005260206000200154905080876000018481548110613a1e57613a1e61447c565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080613a5157613a51614f41565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505061299f565b600091505061299f565b5092915050565b6000818152600183016020526040812054613ae45750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915561299f565b50600061299f565b508054600082559060005260206000209081019061142f9190613b94565b828054828255906000526020600020908101928215613b84579160200282015b82811115613b8457825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190613b2a565b50613b90929150613b94565b5090565b5b80821115613b905760008155600101613b95565b6000815180845260005b81811015613bcf57602081850181015186830182015201613bb3565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815260006132706020830184613ba9565b73ffffffffffffffffffffffffffffffffffffffff8116811461142f57600080fd5b8035613c4d81613c20565b919050565b60008083601f840112613c6457600080fd5b50813567ffffffffffffffff811115613c7c57600080fd5b602083019150836020828501011115613c9457600080fd5b9250929050565b60008060008060608587031215613cb157600080fd5b8435613cbc81613c20565b935060208501359250604085013567ffffffffffffffff811115613cdf57600080fd5b613ceb87828801613c52565b95989497509550505050565b600080600060408486031215613d0c57600080fd5b83359250602084013567ffffffffffffffff811115613d2a57600080fd5b613d3686828701613c52565b9497909650939450505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610220810167ffffffffffffffff81118282101715613d9657613d96613d43565b60405290565b60405160c0810167ffffffffffffffff81118282101715613d9657613d96613d43565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613e0657613e06613d43565b604052919050565b600067ffffffffffffffff821115613e2857613e28613d43565b5060051b60200190565b600082601f830112613e4357600080fd5b81356020613e58613e5383613e0e565b613dbf565b82815260059290921b84018101918181019086841115613e7757600080fd5b8286015b84811015613e925780358352918301918301613e7b565b509695505050505050565b600080600080600080600060e0888a031215613eb857600080fd5b6060880189811115613ec957600080fd5b8897503567ffffffffffffffff80821115613ee357600080fd5b613eef8b838c01613c52565b909850965060808a0135915080821115613f0857600080fd5b818a0191508a601f830112613f1c57600080fd5b813581811115613f2b57600080fd5b8b60208260051b8501011115613f4057600080fd5b6020830196508095505060a08a0135915080821115613f5e57600080fd5b50613f6b8a828b01613e32565b92505060c0880135905092959891949750929550565b600082601f830112613f9257600080fd5b81356020613fa2613e5383613e0e565b82815260059290921b84018101918181019086841115613fc157600080fd5b8286015b84811015613e92578035613fd881613c20565b8352918301918301613fc5565b803560ff81168114613c4d57600080fd5b600082601f83011261400757600080fd5b813567ffffffffffffffff81111561402157614021613d43565b61405260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601613dbf565b81815284602083860101111561406757600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff81168114613c4d57600080fd5b60008060008060008060c087890312156140b557600080fd5b863567ffffffffffffffff808211156140cd57600080fd5b6140d98a838b01613f81565b975060208901359150808211156140ef57600080fd5b6140fb8a838b01613f81565b965061410960408a01613fe5565b9550606089013591508082111561411f57600080fd5b61412b8a838b01613ff6565b945061413960808a01614084565b935060a089013591508082111561414f57600080fd5b5061415c89828a01613ff6565b9150509295509295509295565b60006020828403121561417b57600080fd5b813561327081613c20565b63ffffffff8116811461142f57600080fd5b8035613c4d81614186565b62ffffff8116811461142f57600080fd5b8035613c4d816141a3565b61ffff8116811461142f57600080fd5b8035613c4d816141bf565b6bffffffffffffffffffffffff8116811461142f57600080fd5b8035613c4d816141da565b801515811461142f57600080fd5b8035613c4d816141ff565b6000610220828403121561422b57600080fd5b614233613d72565b905061423e82614198565b815261424c60208301614198565b602082015261425d60408301614198565b604082015261426e606083016141b4565b606082015261427f608083016141cf565b608082015261429060a083016141f4565b60a08201526142a160c08301614198565b60c08201526142b260e08301614198565b60e08201526101006142c5818401614198565b908201526101206142d7838201614198565b90820152610140828101359082015261016080830135908201526101806142ff818401613c42565b908201526101a08281013567ffffffffffffffff81111561431f57600080fd5b61432b85828601613f81565b8284015250506101c061433f818401613c42565b908201526101e0614351838201613c42565b9082015261020061436383820161420d565b9082015292915050565b60008060008060008060c0878903121561438657600080fd5b863567ffffffffffffffff8082111561439e57600080fd5b6143aa8a838b01613f81565b975060208901359150808211156143c057600080fd5b6143cc8a838b01613f81565b96506143da60408a01613fe5565b955060608901359150808211156143f057600080fd5b61412b8a838b01614218565b60006020828403121561440e57600080fd5b5035919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6bffffffffffffffffffffffff818116838216019080821115613a9657613a96614415565b8082018082111561299f5761299f614415565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60ff818116838216019081111561299f5761299f614415565b6000602082840312156144d657600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b61ffff818116838216019080821115613a9657613a96614415565b8181038181111561299f5761299f614415565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361456b5761456b614415565b5060010190565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156145aa576145aa614415565b500290565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000826145ed576145ed6145af565b500490565b6bffffffffffffffffffffffff851681528360208201528260408201526080606082015260006146256080830184613ba9565b9695505050505050565b8051613c4d81614186565b8051613c4d816141a3565b8051613c4d816141bf565b8051613c4d816141da565b8051613c4d81613c20565b600082601f83011261467757600080fd5b81516020614687613e5383613e0e565b82815260059290921b840181019181810190868411156146a657600080fd5b8286015b84811015613e925780516146bd81613c20565b83529183019183016146aa565b8051613c4d816141ff565b6000602082840312156146e757600080fd5b815167ffffffffffffffff808211156146ff57600080fd5b90830190610220828603121561471457600080fd5b61471c613d72565b6147258361462f565b81526147336020840161462f565b60208201526147446040840161462f565b60408201526147556060840161463a565b606082015261476660808401614645565b608082015261477760a08401614650565b60a082015261478860c0840161462f565b60c082015261479960e0840161462f565b60e08201526101006147ac81850161462f565b908201526101206147be84820161462f565b90820152610140838101519082015261016080840151908201526101806147e681850161465b565b908201526101a083810151838111156147fe57600080fd5b61480a88828701614666565b8284015250506101c0915061482082840161465b565b828201526101e0915061483482840161465b565b8282015261020091506148488284016146ca565b91810191909152949350505050565b600060ff821660ff84168160ff048111821515161561487857614878614415565b029392505050565b63ffffffff818116838216019080821115613a9657613a96614415565b600081518084526020808501945080840160005b838110156148e357815173ffffffffffffffffffffffffffffffffffffffff16875295820195908201906001016148b1565b509495945050505050565b6020815261490560208201835163ffffffff169052565b6000602083015161491e604084018263ffffffff169052565b50604083015163ffffffff8116606084015250606083015162ffffff8116608084015250608083015161ffff811660a08401525060a08301516bffffffffffffffffffffffff811660c08401525060c083015163ffffffff811660e08401525060e08301516101006149978185018363ffffffff169052565b84015190506101206149b08482018363ffffffff169052565b84015190506101406149c98482018363ffffffff169052565b840151610160848101919091528401516101808085019190915284015190506101a0614a0c8185018373ffffffffffffffffffffffffffffffffffffffff169052565b808501519150506102206101c08181860152614a2c61024086018461489d565b908601519092506101e0614a578682018373ffffffffffffffffffffffffffffffffffffffff169052565b8601519050610200614a808682018373ffffffffffffffffffffffffffffffffffffffff169052565b90950151151593019290925250919050565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152614ac28184018a61489d565b90508281036080840152614ad6818961489d565b905060ff871660a084015282810360c0840152614af38187613ba9565b905067ffffffffffffffff851660e0840152828103610100840152614b188185613ba9565b9c9b505050505050505050505050565b828152604060208201526000612a486040830184613ba9565b60008060408385031215614b5457600080fd5b8251614b5f816141ff565b6020939093015192949293505050565b8183823760009101908152919050565b8281526080810160608360208401379392505050565b600082601f830112614ba657600080fd5b81356020614bb6613e5383613e0e565b82815260059290921b84018101918181019086841115614bd557600080fd5b8286015b84811015613e9257803567ffffffffffffffff811115614bf95760008081fd5b614c078986838b0101613ff6565b845250918301918301614bd9565b600060208284031215614c2757600080fd5b813567ffffffffffffffff80821115614c3f57600080fd5b9083019060c08286031215614c5357600080fd5b614c5b613d9c565b8235815260208301356020820152604083013582811115614c7b57600080fd5b614c8787828601613e32565b604083015250606083013582811115614c9f57600080fd5b614cab87828601613e32565b606083015250608083013582811115614cc357600080fd5b614ccf87828601614b95565b60808301525060a083013582811115614ce757600080fd5b614cf387828601614b95565b60a08301525095945050505050565b6bffffffffffffffffffffffff828116828216039080821115613a9657613a96614415565b60006bffffffffffffffffffffffff80841680614d4657614d466145af565b92169190910492915050565b60006bffffffffffffffffffffffff80831681851681830481118215151615614d7d57614d7d614415565b02949350505050565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b166040850152816060850152614dcd8285018b61489d565b91508382036080850152614de1828a61489d565b915060ff881660a085015283820360c0850152614dfe8288613ba9565b90861660e08501528381036101008501529050614b188185613ba9565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b600060408284031215614e7a57600080fd5b6040516040810181811067ffffffffffffffff82111715614e9d57614e9d613d43565b6040528251614eab81614186565b81526020928301519281019290925250919050565b600060a08284031215614ed257600080fd5b60405160a0810181811067ffffffffffffffff82111715614ef557614ef5613d43565b806040525082518152602083015160208201526040830151614f1681614186565b60408201526060830151614f2981614186565b60608201526080928301519281019290925250919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000810000a",
>>>>>>> d1f3a03254 (generate wrappers)
=======
	Bin: "0x6101206040523480156200001257600080fd5b50604051620052f4380380620052f4833981016040819052620000359162000345565b80816001600160a01b031663ca30e6036040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200009b919062000345565b826001600160a01b031663b10b673c6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000da573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000100919062000345565b836001600160a01b0316636709d0e56040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200013f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000165919062000345565b846001600160a01b0316635425d8ac6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca919062000345565b3380600081620002215760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200025457620002548162000281565b5050506001600160a01b0393841660805291831660a052821660c052811660e0521661010052506200036c565b336001600160a01b03821603620002db5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000218565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200034257600080fd5b50565b6000602082840312156200035857600080fd5b815162000365816200032c565b9392505050565b60805160a05160c05160e05161010051614f46620003ae6000396000818160d6015261016f015260005050600050506000505060006104330152614f466000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c8063aed2e92911610081578063e3d0e7121161005b578063e3d0e712146102e0578063f2fde38b146102f3578063f75f6b1114610306576100d4565b8063aed2e92914610262578063afcb95d71461028c578063b1dc65a4146102cd576100d4565b806381ff7048116100b257806381ff7048146101bc5780638da5cb5b14610231578063a4c0ed361461024f576100d4565b8063181f5a771461011b578063349e8cca1461016d57806379ba5097146101b4575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e808015610114573d6000f35b3d6000fd5b005b6101576040518060400160405280601881526020017f4175746f6d6174696f6e526567697374727920322e322e30000000000000000081525081565b6040516101649190613c0d565b60405180910390f35b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610164565b610119610319565b61020e60155460115463ffffffff780100000000000000000000000000000000000000000000000083048116937c01000000000000000000000000000000000000000000000000000000009093041691565b6040805163ffffffff948516815293909216602084015290820152606001610164565b60005473ffffffffffffffffffffffffffffffffffffffff1661018f565b61011961025d366004613c9b565b61041b565b610275610270366004613cf7565b610637565b604080519215158352602083019190915201610164565b601154601254604080516000815260208101939093527401000000000000000000000000000000000000000090910463ffffffff1690820152606001610164565b6101196102db366004613e9d565b6107ad565b6101196102ee36600461409c565b6113f5565b610119610301366004614169565b61141e565b61011961031436600461436d565b611432565b60015473ffffffffffffffffffffffffffffffffffffffff16331461039f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461048a576040517fc8bad78d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b602081146104c4576040517fdfe9309000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006104d2828401846143fc565b60008181526004602052604090205490915065010000000000900463ffffffff9081161461052c576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000818152600460205260409020600101546105679085906c0100000000000000000000000090046bffffffffffffffffffffffff16614444565b600082815260046020526040902060010180546bffffffffffffffffffffffff929092166c01000000000000000000000000027fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff9092169190911790556019546105d2908590614469565b6019556040516bffffffffffffffffffffffff8516815273ffffffffffffffffffffffffffffffffffffffff86169082907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a35050505050565b60008061064261241e565b6012547e01000000000000000000000000000000000000000000000000000000000000900460ff16156106a1576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600085815260046020908152604091829020825160e081018452815460ff811615158252610100810463ffffffff908116838601819052650100000000008304821684880152690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff16606084018190526001909401546bffffffffffffffffffffffff80821660808601526c0100000000000000000000000082041660a0850152780100000000000000000000000000000000000000000000000090041660c08301528451601f890185900485028101850190955287855290936107a093899089908190840183828082843760009201919091525061245892505050565b9097909650945050505050565b60005a60408051610160810182526012546bffffffffffffffffffffffff8116825263ffffffff6c010000000000000000000000008204811660208401527001000000000000000000000000000000008204811693830193909352740100000000000000000000000000000000000000008104909216606082015262ffffff7801000000000000000000000000000000000000000000000000830416608082015261ffff7b0100000000000000000000000000000000000000000000000000000083041660a082015260ff7d0100000000000000000000000000000000000000000000000000000000008304811660c08301527e010000000000000000000000000000000000000000000000000000000000008304811615801560e08401527f01000000000000000000000000000000000000000000000000000000000000009093048116151561010080840191909152601354918216151561012084015273ffffffffffffffffffffffffffffffffffffffff910416610140820152919250610963576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600b602052604090205460ff166109ac576040517f1099ed7500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6011548935146109e8576040517fdfdcf8e700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60c08101516109f89060016144ab565b60ff1685141580610a0a575083518514155b15610a41576040517f0244f71a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610a5089898989898989612681565b6000610a5c89896128ea565b9050600081604001515167ffffffffffffffff811115610a7e57610a7e613d43565b604051908082528060200260200182016040528015610b4257816020015b604080516101e0810182526000610100820181815261012083018290526101408301829052610160830182905261018083018290526101a083018290526101c0830182905282526020808301829052928201819052606082018190526080820181905260a0820181905260c0820181905260e082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181610a9c5790505b50905060008084610140015173ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa158015610b98573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bbc91906144c4565b905060005b846040015151811015611006576004600086604001518381518110610be857610be861447c565b6020908102919091018101518252818101929092526040908101600020815160e081018352815460ff811615158252610100810463ffffffff90811695830195909552650100000000008104851693820193909352690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff166060830152600101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490911660c08201528451859083908110610ccd57610ccd61447c565b602002602001015160000181905250610d0285604001518281518110610cf557610cf561447c565b60200260200101516129a5565b848281518110610d1457610d1461447c565b6020026020010151608001906001811115610d3157610d316144dd565b90816001811115610d4457610d446144dd565b81525050610db886858381518110610d5e57610d5e61447c565b60200260200101516080015187606001518481518110610d8057610d8061447c565b60200260200101518860a001518581518110610d9e57610d9e61447c565b60200260200101515189600001518a602001516001612a50565b848281518110610dca57610dca61447c565b6020026020010151604001906bffffffffffffffffffffffff1690816bffffffffffffffffffffffff1681525050610e5885604001518281518110610e1157610e1161447c565b60200260200101518387608001518481518110610e3057610e3061447c565b6020026020010151878581518110610e4a57610e4a61447c565b60200260200101518a612a9b565b858381518110610e6a57610e6a61447c565b6020026020010151602001868481518110610e8757610e8761447c565b602002602001015160e0018281525082151515158152505050838181518110610eb257610eb261447c565b60200260200101516020015115610ed557610ece60018461450c565b9250610eda565b610ff4565b610f40848281518110610eef57610eef61447c565b6020026020010151600001516060015186606001518381518110610f1557610f1561447c565b60200260200101518760a001518481518110610f3357610f3361447c565b6020026020010151612458565b858381518110610f5257610f5261447c565b6020026020010151606001868481518110610f6f57610f6f61447c565b602002602001015160a0018281525082151515158152505050838181518110610f9a57610f9a61447c565b602002602001015160a0015187610fb19190614527565b9650610ff485604001518281518110610fcc57610fcc61447c565b602002602001015183868481518110610fe757610fe761447c565b6020026020010151612c1a565b80610ffe8161453a565b915050610bc1565b508161ffff1660000361101e575050505050506113ec565b60c085015161102e9060016144ab565b61103d9060ff1661047e614572565b616fe061104b8d6010614572565b5a611056908a614527565b6110609190614469565b61106a9190614469565b6110749190614469565b9550611e7861108761ffff8416886145b8565b6110919190614469565b955060008060008060005b886040015151811015611292578781815181106110bb576110bb61447c565b60200260200101516020015115611280576111178b8983815181106110e2576110e261447c565b6020026020010151608001518b60a0015184815181106111045761110461447c565b6020026020010151518d60c00151612d1f565b8882815181106111295761112961447c565b602002602001015160c00181815250506111858a8a6040015183815181106111535761115361447c565b60200260200101518a848151811061116d5761116d61447c565b60200260200101518c600001518d602001518c612d3f565b90935091506111948285614444565b93506111a08386614444565b94508781815181106111b4576111b461447c565b6020026020010151606001511515896040015182815181106111d8576111d861447c565b60200260200101517fad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b848661120d9190614444565b8b858151811061121f5761121f61447c565b602002602001015160a001518c868151811061123d5761123d61447c565b602002602001015160c001518e60800151878151811061125f5761125f61447c565b602002602001015160405161127794939291906145cc565b60405180910390a35b8061128a8161453a565b91505061109c565b5050336000908152600b6020526040902080548492506002906112ca9084906201000090046bffffffffffffffffffffffff16614444565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555080601260000160008282829054906101000a90046bffffffffffffffffffffffff166113249190614444565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060008f6001600381106113675761136761447c565b602002013560001c9050600060088264ffffffffff16901c9050886060015163ffffffff168163ffffffff1611156113e157601280547fffffffffffffffff00000000ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000063ffffffff8416021790555b505050505050505050505b50505050505050565b6114168686868680602001905181019061140f91906146af565b8686611432565b505050505050565b611426612e32565b61142f81612eb3565b50565b61143a612e32565b601f86511115611476576040517f25d0209c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8360ff166000036114b3576040517fe77dba5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845186511415806114d257506114ca846003614831565b60ff16865111155b15611509576040517f1d2d1c5800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601254600e546bffffffffffffffffffffffff9091169060005b816bffffffffffffffffffffffff1681101561158b57611578600e828154811061154f5761154f61447c565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff168484612fa8565b50806115838161453a565b915050611523565b5060008060005b836bffffffffffffffffffffffff1681101561169457600d81815481106115bb576115bb61447c565b600091825260209091200154600e805473ffffffffffffffffffffffffffffffffffffffff909216945090829081106115f6576115f661447c565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff8681168452600c8352604080852080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001690559116808452600b90925290912080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905591508061168c8161453a565b915050611592565b506116a1600d6000613aec565b6116ad600e6000613aec565b604080516080810182526000808252602082018190529181018290526060810182905290805b8c51811015611b1657600c60008e83815181106116f2576116f261447c565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528101919091526040016000205460ff161561175d576040517f77cea0fa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168d82815181106117875761178761447c565b602002602001015173ffffffffffffffffffffffffffffffffffffffff16036117dc576040517f815e1d6400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405180604001604052806001151581526020018260ff16815250600c60008f848151811061180d5761180d61447c565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528181019290925260400160002082518154939092015160ff16610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909316929092171790558b518c90829081106118b5576118b561447c565b60200260200101519150600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603611925576040517f58a70a0a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff82166000908152600b60209081526040918290208251608081018452905460ff80821615801584526101008304909116938301939093526bffffffffffffffffffffffff6201000082048116948301949094526e010000000000000000000000000000900490921660608301529093506119e0576040517f6a7281ad00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001835260ff80821660208086019182526bffffffffffffffffffffffff808b166060880190815273ffffffffffffffffffffffffffffffffffffffff87166000908152600b909352604092839020885181549551948a0151925184166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff939094166201000002929092167fffffffffffff000000000000000000000000000000000000000000000000ffff94909616610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009095169490941717919091169290921791909117905580611b0e8161453a565b9150506116d3565b50508a51611b2c9150600d9060208d0190613b0a565b508851611b4090600e9060208c0190613b0a565b50604051806101600160405280856bffffffffffffffffffffffff168152602001886000015163ffffffff168152602001886020015163ffffffff168152602001600063ffffffff168152602001886060015162ffffff168152602001886080015161ffff1681526020018960ff1681526020016012600001601e9054906101000a900460ff16151581526020016012600001601f9054906101000a900460ff161515815260200188610200015115158152602001886101e0015173ffffffffffffffffffffffffffffffffffffffff16815250601260008201518160000160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550602082015181600001600c6101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160106101000a81548163ffffffff021916908363ffffffff16021790555060608201518160000160146101000a81548163ffffffff021916908363ffffffff16021790555060808201518160000160186101000a81548162ffffff021916908362ffffff16021790555060a082015181600001601b6101000a81548161ffff021916908361ffff16021790555060c082015181600001601d6101000a81548160ff021916908360ff16021790555060e082015181600001601e6101000a81548160ff02191690831515021790555061010082015181600001601f6101000a81548160ff0219169083151502179055506101208201518160010160006101000a81548160ff0219169083151502179055506101408201518160010160016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055509050506040518061018001604052808860a001516bffffffffffffffffffffffff16815260200188610180015173ffffffffffffffffffffffffffffffffffffffff168152602001601460010160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff168152602001886040015163ffffffff1681526020018860c0015163ffffffff168152602001601460010160149054906101000a900463ffffffff1663ffffffff168152602001601460010160189054906101000a900463ffffffff1663ffffffff1681526020016014600101601c9054906101000a900463ffffffff1663ffffffff1681526020018860e0015163ffffffff16815260200188610100015163ffffffff16815260200188610120015163ffffffff168152602001886101c0015173ffffffffffffffffffffffffffffffffffffffff16815250601460008201518160000160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550602082015181600001600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060408201518160010160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550606082015181600101600c6101000a81548163ffffffff021916908363ffffffff16021790555060808201518160010160106101000a81548163ffffffff021916908363ffffffff16021790555060a08201518160010160146101000a81548163ffffffff021916908363ffffffff16021790555060c08201518160010160186101000a81548163ffffffff021916908363ffffffff16021790555060e082015181600101601c6101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160020160006101000a81548163ffffffff021916908363ffffffff1602179055506101208201518160020160046101000a81548163ffffffff021916908363ffffffff1602179055506101408201518160020160086101000a81548163ffffffff021916908363ffffffff16021790555061016082015181600201600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505086610140015160178190555086610160015160188190555060006014600101601c9054906101000a900463ffffffff169050876101e0015173ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa1580156121e1573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061220591906144c4565b601580547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167c010000000000000000000000000000000000000000000000000000000063ffffffff9384160217808255600192601891612280918591780100000000000000000000000000000000000000000000000090041661484d565b92506101000a81548163ffffffff021916908363ffffffff1602179055506000886040516020016122b191906148bb565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291905260155490915061231a90469030907801000000000000000000000000000000000000000000000000900463ffffffff168f8f8f878f8f6131b0565b60115560005b61232a600961325a565b81101561235a5761234761233f600983613264565b600990613277565b50806123528161453a565b915050612320565b5060005b896101a00151518110156123b15761239e8a6101a0015182815181106123865761238661447c565b6020026020010151600961329990919063ffffffff16565b50806123a98161453a565b91505061235e565b507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0582601154601460010160189054906101000a900463ffffffff168f8f8f878f8f60405161240899989796959493929190614a5f565b60405180910390a1505050505050505050505050565b3215612456576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60125460009081907f0100000000000000000000000000000000000000000000000000000000000000900460ff16156124bd576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601280547effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f01000000000000000000000000000000000000000000000000000000000000001790556040517f4585e33b0000000000000000000000000000000000000000000000000000000090612539908590602401613c0d565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009094169390931790925290517f79188d1600000000000000000000000000000000000000000000000000000000815290935073ffffffffffffffffffffffffffffffffffffffff8616906379188d169061260c9087908790600401614af5565b60408051808303816000875af115801561262a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061264e9190614b0e565b601280547effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff16905590969095509350505050565b60008686604051612693929190614b3c565b6040519081900381206126aa918a90602001614b4c565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201208383019092526000808452908301819052909250906000805b87811015612882576001858783602081106127165761271661447c565b61272391901a601b6144ab565b8b8b858181106127355761273561447c565b905060200201358a858151811061274e5761274e61447c565b60200260200101516040516000815260200160405260405161278c949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa1580156127ae573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff81166000908152600c602090815290849020838501909452925460ff808216151580855261010090920416938301939093529095509350905061285c576040517f0f4c073700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826020015160080260ff166001901b84019350808061287a9061453a565b9150506126f9565b50827e010101010101010101010101010101010101010101010101010101010101018416146128dd576040517fc103be2e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050505050505050505050565b6129236040518060c001604052806000815260200160008152602001606081526020016060815260200160608152602001606081525090565b600061293183850185614be2565b604081015151606082015151919250908114158061295457508082608001515114155b806129645750808260a001515114155b1561299b576040517fb55ac75400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5090505b92915050565b6000818160045b600f811015612a32577fff0000000000000000000000000000000000000000000000000000000000000082168382602081106129ea576129ea61447c565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614612a2057506000949350505050565b80612a2a8161453a565b9150506129ac565b5081600f1a6001811115612a4857612a486144dd565b949350505050565b600080612a6288878b60c001516132bb565b9050600080612a7d8b8a63ffffffff16858a8a60018b613347565b9092509050612a8c8183614444565b9b9a5050505050505050505050565b600080808085608001516001811115612ab657612ab66144dd565b03612adc57612ac888888888886135e2565b612ad757600092509050612c10565b612b54565b600185608001516001811115612af457612af46144dd565b03612b22576000612b078989898861376b565b9250905080612b1c5750600092509050612c10565b50612b54565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b84516040015163ffffffff168710612ba957877fc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd563687604051612b969190613c0d565b60405180910390a2600092509050612c10565b84604001516bffffffffffffffffffffffff16856000015160a001516bffffffffffffffffffffffff161015612c0957877f377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e0287604051612b969190613c0d565b6001925090505b9550959350505050565b600081608001516001811115612c3257612c326144dd565b03612c9657600083815260046020526040902060010180547fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff16780100000000000000000000000000000000000000000000000063ffffffff851602179055505050565b600181608001516001811115612cae57612cae6144dd565b03612d1a5760e08101805160009081526008602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055915191517fa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f29190a25b505050565b6000612d2c8484846132bb565b905080851015612a485750929392505050565b600080612d5a888760a001518860c001518888886001613347565b90925090506000612d6b8284614444565b600089815260046020526040902060010180549192508291600c90612daf9084906c0100000000000000000000000090046bffffffffffffffffffffffff16614ccf565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915560008a815260046020526040812060010180548594509092612df891859116614444565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555050965096945050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314612456576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610396565b3373ffffffffffffffffffffffffffffffffffffffff821603612f32576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610396565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e01000000000000000000000000000090049091166060820152906131a45760008160600151856130409190614ccf565b9050600061304e8583614cf4565b905080836040018181516130629190614444565b6bffffffffffffffffffffffff1690525061307d8582614d1f565b8360600181815161308e9190614444565b6bffffffffffffffffffffffff90811690915273ffffffffffffffffffffffffffffffffffffffff89166000908152600b602090815260409182902087518154928901519389015160608a015186166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff919096166201000002167fffffffffffff000000000000000000000000000000000000000000000000ffff60ff95909516610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909416939093171792909216179190911790555050505b60400151949350505050565b6000808a8a8a8a8a8a8a8a8a6040516020016131d499989796959493929190614d4f565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179b9a5050505050505050505050565b600061299f825490565b60006132708383613979565b9392505050565b60006132708373ffffffffffffffffffffffffffffffffffffffff84166139a3565b60006132708373ffffffffffffffffffffffffffffffffffffffff8416613a9d565b600080808560018111156132d1576132d16144dd565b036132e0575062015f906132ff565b60018560018111156132f4576132f46144dd565b03612b2257506201b5e45b61331063ffffffff85166014614572565b61331b8460016144ab565b61332a9060ff16611d4c614572565b6133349083614469565b61333e9190614469565b95945050505050565b60008060008960a0015161ffff16876133609190614572565b905083801561336e5750803a105b1561337657503a5b60008415613407578a610140015173ffffffffffffffffffffffffffffffffffffffff166349948e0e6000366040518363ffffffff1660e01b81526004016133bf929190614de4565b602060405180830381865afa1580156133dc573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061340091906144c4565b90506134c3565b6101408b01516016546040517f1254414000000000000000000000000000000000000000000000000000000000815264010000000090910463ffffffff16600482015273ffffffffffffffffffffffffffffffffffffffff90911690631254414090602401602060405180830381865afa158015613489573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906134ad91906144c4565b8b60a0015161ffff166134c09190614572565b90505b6134d161ffff8716826145b8565b9050600087826134e18c8e614469565b6134eb9086614572565b6134f59190614469565b61350790670de0b6b3a7640000614572565b61351191906145b8565b905060008c6040015163ffffffff1664e8d4a510006135309190614572565b898e6020015163ffffffff16858f886135499190614572565b6135539190614469565b61356190633b9aca00614572565b61356b9190614572565b61357591906145b8565b61357f9190614469565b90506b033b2e3c9fd0803ce80000006135988284614469565b11156135d0576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b600080848060200190518101906135f99190614e31565b845160c00151815191925063ffffffff9081169116101561365657867f405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8866040516136449190613c0d565b60405180910390a2600091505061333e565b826101200151801561371757506020810151158015906137175750602081015161014084015182516040517f85df51fd00000000000000000000000000000000000000000000000000000000815263ffffffff909116600482015273ffffffffffffffffffffffffffffffffffffffff909116906385df51fd90602401602060405180830381865afa1580156136f0573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061371491906144c4565b14155b806137295750805163ffffffff168611155b1561375e57867f6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301866040516136449190613c0d565b5060019695505050505050565b6000806000848060200190518101906137849190614e89565b90506000878260000151836020015184604001516040516020016137e694939291909384526020840192909252604083015260e01b7fffffffff0000000000000000000000000000000000000000000000000000000016606082015260640190565b60405160208183030381529060405280519060200120905084610120015180156138c257506080820151158015906138c25750608082015161014086015160608401516040517f85df51fd00000000000000000000000000000000000000000000000000000000815263ffffffff909116600482015273ffffffffffffffffffffffffffffffffffffffff909116906385df51fd90602401602060405180830381865afa15801561389b573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906138bf91906144c4565b14155b806138d7575086826060015163ffffffff1610155b1561392157877f6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc3018760405161390c9190613c0d565b60405180910390a26000935091506139709050565b60008181526008602052604090205460ff161561396857877f405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e88760405161390c9190613c0d565b600193509150505b94509492505050565b60008260000182815481106139905761399061447c565b9060005260206000200154905092915050565b60008181526001830160205260408120548015613a8c5760006139c7600183614527565b85549091506000906139db90600190614527565b9050818114613a405760008660000182815481106139fb576139fb61447c565b9060005260206000200154905080876000018481548110613a1e57613a1e61447c565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080613a5157613a51614f0a565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505061299f565b600091505061299f565b5092915050565b6000818152600183016020526040812054613ae45750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915561299f565b50600061299f565b508054600082559060005260206000209081019061142f9190613b94565b828054828255906000526020600020908101928215613b84579160200282015b82811115613b8457825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190613b2a565b50613b90929150613b94565b5090565b5b80821115613b905760008155600101613b95565b6000815180845260005b81811015613bcf57602081850181015186830182015201613bb3565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815260006132706020830184613ba9565b73ffffffffffffffffffffffffffffffffffffffff8116811461142f57600080fd5b8035613c4d81613c20565b919050565b60008083601f840112613c6457600080fd5b50813567ffffffffffffffff811115613c7c57600080fd5b602083019150836020828501011115613c9457600080fd5b9250929050565b60008060008060608587031215613cb157600080fd5b8435613cbc81613c20565b935060208501359250604085013567ffffffffffffffff811115613cdf57600080fd5b613ceb87828801613c52565b95989497509550505050565b600080600060408486031215613d0c57600080fd5b83359250602084013567ffffffffffffffff811115613d2a57600080fd5b613d3686828701613c52565b9497909650939450505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610220810167ffffffffffffffff81118282101715613d9657613d96613d43565b60405290565b60405160c0810167ffffffffffffffff81118282101715613d9657613d96613d43565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613e0657613e06613d43565b604052919050565b600067ffffffffffffffff821115613e2857613e28613d43565b5060051b60200190565b600082601f830112613e4357600080fd5b81356020613e58613e5383613e0e565b613dbf565b82815260059290921b84018101918181019086841115613e7757600080fd5b8286015b84811015613e925780358352918301918301613e7b565b509695505050505050565b600080600080600080600060e0888a031215613eb857600080fd5b6060880189811115613ec957600080fd5b8897503567ffffffffffffffff80821115613ee357600080fd5b613eef8b838c01613c52565b909850965060808a0135915080821115613f0857600080fd5b818a0191508a601f830112613f1c57600080fd5b813581811115613f2b57600080fd5b8b60208260051b8501011115613f4057600080fd5b6020830196508095505060a08a0135915080821115613f5e57600080fd5b50613f6b8a828b01613e32565b92505060c0880135905092959891949750929550565b600082601f830112613f9257600080fd5b81356020613fa2613e5383613e0e565b82815260059290921b84018101918181019086841115613fc157600080fd5b8286015b84811015613e92578035613fd881613c20565b8352918301918301613fc5565b803560ff81168114613c4d57600080fd5b600082601f83011261400757600080fd5b813567ffffffffffffffff81111561402157614021613d43565b61405260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601613dbf565b81815284602083860101111561406757600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff81168114613c4d57600080fd5b60008060008060008060c087890312156140b557600080fd5b863567ffffffffffffffff808211156140cd57600080fd5b6140d98a838b01613f81565b975060208901359150808211156140ef57600080fd5b6140fb8a838b01613f81565b965061410960408a01613fe5565b9550606089013591508082111561411f57600080fd5b61412b8a838b01613ff6565b945061413960808a01614084565b935060a089013591508082111561414f57600080fd5b5061415c89828a01613ff6565b9150509295509295509295565b60006020828403121561417b57600080fd5b813561327081613c20565b63ffffffff8116811461142f57600080fd5b8035613c4d81614186565b62ffffff8116811461142f57600080fd5b8035613c4d816141a3565b61ffff8116811461142f57600080fd5b8035613c4d816141bf565b6bffffffffffffffffffffffff8116811461142f57600080fd5b8035613c4d816141da565b801515811461142f57600080fd5b8035613c4d816141ff565b6000610220828403121561422b57600080fd5b614233613d72565b905061423e82614198565b815261424c60208301614198565b602082015261425d60408301614198565b604082015261426e606083016141b4565b606082015261427f608083016141cf565b608082015261429060a083016141f4565b60a08201526142a160c08301614198565b60c08201526142b260e08301614198565b60e08201526101006142c5818401614198565b908201526101206142d7838201614198565b90820152610140828101359082015261016080830135908201526101806142ff818401613c42565b908201526101a08281013567ffffffffffffffff81111561431f57600080fd5b61432b85828601613f81565b8284015250506101c061433f818401613c42565b908201526101e0614351838201613c42565b9082015261020061436383820161420d565b9082015292915050565b60008060008060008060c0878903121561438657600080fd5b863567ffffffffffffffff8082111561439e57600080fd5b6143aa8a838b01613f81565b975060208901359150808211156143c057600080fd5b6143cc8a838b01613f81565b96506143da60408a01613fe5565b955060608901359150808211156143f057600080fd5b61412b8a838b01614218565b60006020828403121561440e57600080fd5b5035919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6bffffffffffffffffffffffff818116838216019080821115613a9657613a96614415565b8082018082111561299f5761299f614415565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60ff818116838216019081111561299f5761299f614415565b6000602082840312156144d657600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b61ffff818116838216019080821115613a9657613a96614415565b8181038181111561299f5761299f614415565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361456b5761456b614415565b5060010190565b808202811582820484141761299f5761299f614415565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000826145c7576145c7614589565b500490565b6bffffffffffffffffffffffff851681528360208201528260408201526080606082015260006145ff6080830184613ba9565b9695505050505050565b8051613c4d81614186565b8051613c4d816141a3565b8051613c4d816141bf565b8051613c4d816141da565b8051613c4d81613c20565b600082601f83011261465157600080fd5b81516020614661613e5383613e0e565b82815260059290921b8401810191818101908684111561468057600080fd5b8286015b84811015613e9257805161469781613c20565b8352918301918301614684565b8051613c4d816141ff565b6000602082840312156146c157600080fd5b815167ffffffffffffffff808211156146d957600080fd5b9083019061022082860312156146ee57600080fd5b6146f6613d72565b6146ff83614609565b815261470d60208401614609565b602082015261471e60408401614609565b604082015261472f60608401614614565b60608201526147406080840161461f565b608082015261475160a0840161462a565b60a082015261476260c08401614609565b60c082015261477360e08401614609565b60e0820152610100614786818501614609565b90820152610120614798848201614609565b90820152610140838101519082015261016080840151908201526101806147c0818501614635565b908201526101a083810151838111156147d857600080fd5b6147e488828701614640565b8284015250506101c091506147fa828401614635565b828201526101e0915061480e828401614635565b8282015261020091506148228284016146a4565b91810191909152949350505050565b60ff8181168382160290811690818114613a9657613a96614415565b63ffffffff818116838216019080821115613a9657613a96614415565b600081518084526020808501945080840160005b838110156148b057815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010161487e565b509495945050505050565b602081526148d260208201835163ffffffff169052565b600060208301516148eb604084018263ffffffff169052565b50604083015163ffffffff8116606084015250606083015162ffffff8116608084015250608083015161ffff811660a08401525060a08301516bffffffffffffffffffffffff811660c08401525060c083015163ffffffff811660e08401525060e08301516101006149648185018363ffffffff169052565b840151905061012061497d8482018363ffffffff169052565b84015190506101406149968482018363ffffffff169052565b840151610160848101919091528401516101808085019190915284015190506101a06149d98185018373ffffffffffffffffffffffffffffffffffffffff169052565b808501519150506102206101c081818601526149f961024086018461486a565b908601519092506101e0614a248682018373ffffffffffffffffffffffffffffffffffffffff169052565b8601519050610200614a4d8682018373ffffffffffffffffffffffffffffffffffffffff169052565b90950151151593019290925250919050565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152614a8f8184018a61486a565b90508281036080840152614aa3818961486a565b905060ff871660a084015282810360c0840152614ac08187613ba9565b905067ffffffffffffffff851660e0840152828103610100840152614ae58185613ba9565b9c9b505050505050505050505050565b828152604060208201526000612a486040830184613ba9565b60008060408385031215614b2157600080fd5b8251614b2c816141ff565b6020939093015192949293505050565b8183823760009101908152919050565b8281526080810160608360208401379392505050565b600082601f830112614b7357600080fd5b81356020614b83613e5383613e0e565b82815260059290921b84018101918181019086841115614ba257600080fd5b8286015b84811015613e9257803567ffffffffffffffff811115614bc65760008081fd5b614bd48986838b0101613ff6565b845250918301918301614ba6565b600060208284031215614bf457600080fd5b813567ffffffffffffffff80821115614c0c57600080fd5b9083019060c08286031215614c2057600080fd5b614c28613d9c565b8235815260208301356020820152604083013582811115614c4857600080fd5b614c5487828601613e32565b604083015250606083013582811115614c6c57600080fd5b614c7887828601613e32565b606083015250608083013582811115614c9057600080fd5b614c9c87828601614b62565b60808301525060a083013582811115614cb457600080fd5b614cc087828601614b62565b60a08301525095945050505050565b6bffffffffffffffffffffffff828116828216039080821115613a9657613a96614415565b60006bffffffffffffffffffffffff80841680614d1357614d13614589565b92169190910492915050565b6bffffffffffffffffffffffff818116838216028082169190828114614d4757614d47614415565b505092915050565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b166040850152816060850152614d968285018b61486a565b91508382036080850152614daa828a61486a565b915060ff881660a085015283820360c0850152614dc78288613ba9565b90861660e08501528381036101008501529050614ae58185613ba9565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b600060408284031215614e4357600080fd5b6040516040810181811067ffffffffffffffff82111715614e6657614e66613d43565b6040528251614e7481614186565b81526020928301519281019290925250919050565b600060a08284031215614e9b57600080fd5b60405160a0810181811067ffffffffffffffff82111715614ebe57614ebe613d43565b806040525082518152602083015160208201526040830151614edf81614186565b60408201526060830151614ef281614186565b60608201526080928301519281019290925250919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000813000a",
>>>>>>> e15f3cf45d (regen wrappers)
=======
	Bin: "0x6101206040523480156200001257600080fd5b506040516200529938038062005299833981016040819052620000359162000345565b80816001600160a01b031663ca30e6036040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200009b919062000345565b826001600160a01b031663b10b673c6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000da573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000100919062000345565b836001600160a01b0316636709d0e56040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200013f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000165919062000345565b846001600160a01b0316635425d8ac6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca919062000345565b3380600081620002215760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200025457620002548162000281565b5050506001600160a01b0393841660805291831660a052821660c052811660e0521661010052506200036c565b336001600160a01b03821603620002db5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000218565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200034257600080fd5b50565b6000602082840312156200035857600080fd5b815162000365816200032c565b9392505050565b60805160a05160c05160e05161010051614eeb620003ae6000396000818160d6015261016f015260005050600050506000505060006104330152614eeb6000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c8063aed2e92911610081578063e3d0e7121161005b578063e3d0e712146102e0578063f2fde38b146102f3578063f75f6b1114610306576100d4565b8063aed2e92914610262578063afcb95d71461028c578063b1dc65a4146102cd576100d4565b806381ff7048116100b257806381ff7048146101bc5780638da5cb5b14610231578063a4c0ed361461024f576100d4565b8063181f5a771461011b578063349e8cca1461016d57806379ba5097146101b4575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e808015610114573d6000f35b3d6000fd5b005b6101576040518060400160405280601881526020017f4175746f6d6174696f6e526567697374727920322e322e30000000000000000081525081565b6040516101649190613bff565b60405180910390f35b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610164565b610119610319565b61020e60155460115463ffffffff780100000000000000000000000000000000000000000000000083048116937c01000000000000000000000000000000000000000000000000000000009093041691565b6040805163ffffffff948516815293909216602084015290820152606001610164565b60005473ffffffffffffffffffffffffffffffffffffffff1661018f565b61011961025d366004613c8d565b61041b565b610275610270366004613ce9565b610637565b604080519215158352602083019190915201610164565b601154601254604080516000815260208101939093527401000000000000000000000000000000000000000090910463ffffffff1690820152606001610164565b6101196102db366004613e8f565b6107ad565b6101196102ee36600461408e565b6113f5565b61011961030136600461415b565b61141e565b61011961031436600461435f565b611432565b60015473ffffffffffffffffffffffffffffffffffffffff16331461039f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461048a576040517fc8bad78d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b602081146104c4576040517fdfe9309000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006104d2828401846143ee565b60008181526004602052604090205490915065010000000000900463ffffffff9081161461052c576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000818152600460205260409020600101546105679085906c0100000000000000000000000090046bffffffffffffffffffffffff16614436565b600082815260046020526040902060010180546bffffffffffffffffffffffff929092166c01000000000000000000000000027fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff9092169190911790556019546105d290859061445b565b6019556040516bffffffffffffffffffffffff8516815273ffffffffffffffffffffffffffffffffffffffff86169082907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a35050505050565b60008061064261241e565b6012547e01000000000000000000000000000000000000000000000000000000000000900460ff16156106a1576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600085815260046020908152604091829020825160e081018452815460ff811615158252610100810463ffffffff908116838601819052650100000000008304821684880152690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff16606084018190526001909401546bffffffffffffffffffffffff80821660808601526c0100000000000000000000000082041660a0850152780100000000000000000000000000000000000000000000000090041660c08301528451601f890185900485028101850190955287855290936107a093899089908190840183828082843760009201919091525061245892505050565b9097909650945050505050565b60005a60408051610160810182526012546bffffffffffffffffffffffff8116825263ffffffff6c010000000000000000000000008204811660208401527001000000000000000000000000000000008204811693830193909352740100000000000000000000000000000000000000008104909216606082015262ffffff7801000000000000000000000000000000000000000000000000830416608082015261ffff7b0100000000000000000000000000000000000000000000000000000083041660a082015260ff7d0100000000000000000000000000000000000000000000000000000000008304811660c08301527e010000000000000000000000000000000000000000000000000000000000008304811615801560e08401527f01000000000000000000000000000000000000000000000000000000000000009093048116151561010080840191909152601354918216151561012084015273ffffffffffffffffffffffffffffffffffffffff910416610140820152919250610963576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600b602052604090205460ff166109ac576040517f1099ed7500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6011548935146109e8576040517fdfdcf8e700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60c08101516109f890600161449d565b60ff1685141580610a0a575083518514155b15610a41576040517f0244f71a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610a5089898989898989612681565b6000610a5c89896128ea565b9050600081604001515167ffffffffffffffff811115610a7e57610a7e613d35565b604051908082528060200260200182016040528015610b4257816020015b604080516101e0810182526000610100820181815261012083018290526101408301829052610160830182905261018083018290526101a083018290526101c0830182905282526020808301829052928201819052606082018190526080820181905260a0820181905260c0820181905260e082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181610a9c5790505b50905060008084610140015173ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa158015610b98573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bbc91906144b6565b905060005b846040015151811015611006576004600086604001518381518110610be857610be861446e565b6020908102919091018101518252818101929092526040908101600020815160e081018352815460ff811615158252610100810463ffffffff90811695830195909552650100000000008104851693820193909352690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff166060830152600101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490911660c08201528451859083908110610ccd57610ccd61446e565b602002602001015160000181905250610d0285604001518281518110610cf557610cf561446e565b60200260200101516129a5565b848281518110610d1457610d1461446e565b6020026020010151608001906001811115610d3157610d316144cf565b90816001811115610d4457610d446144cf565b81525050610db886858381518110610d5e57610d5e61446e565b60200260200101516080015187606001518481518110610d8057610d8061446e565b60200260200101518860a001518581518110610d9e57610d9e61446e565b60200260200101515189600001518a602001516001612a50565b848281518110610dca57610dca61446e565b6020026020010151604001906bffffffffffffffffffffffff1690816bffffffffffffffffffffffff1681525050610e5885604001518281518110610e1157610e1161446e565b60200260200101518387608001518481518110610e3057610e3061446e565b6020026020010151878581518110610e4a57610e4a61446e565b60200260200101518a612a9b565b858381518110610e6a57610e6a61446e565b6020026020010151602001868481518110610e8757610e8761446e565b602002602001015160e0018281525082151515158152505050838181518110610eb257610eb261446e565b60200260200101516020015115610ed557610ece6001846144fe565b9250610eda565b610ff4565b610f40848281518110610eef57610eef61446e565b6020026020010151600001516060015186606001518381518110610f1557610f1561446e565b60200260200101518760a001518481518110610f3357610f3361446e565b6020026020010151612458565b858381518110610f5257610f5261446e565b6020026020010151606001868481518110610f6f57610f6f61446e565b602002602001015160a0018281525082151515158152505050838181518110610f9a57610f9a61446e565b602002602001015160a0015187610fb19190614519565b9650610ff485604001518281518110610fcc57610fcc61446e565b602002602001015183868481518110610fe757610fe761446e565b6020026020010151612c1a565b80610ffe8161452c565b915050610bc1565b508161ffff1660000361101e575050505050506113ec565b60c085015161102e90600161449d565b61103d9060ff1661044c614564565b616fe061104b8d6010614564565b5a611056908a614519565b611060919061445b565b61106a919061445b565b611074919061445b565b9550611d4c61108761ffff8416886145aa565b611091919061445b565b955060008060008060005b886040015151811015611292578781815181106110bb576110bb61446e565b60200260200101516020015115611280576111178b8983815181106110e2576110e261446e565b6020026020010151608001518b60a0015184815181106111045761110461446e565b6020026020010151518d60c00151612d1f565b8882815181106111295761112961446e565b602002602001015160c00181815250506111858a8a6040015183815181106111535761115361446e565b60200260200101518a848151811061116d5761116d61446e565b60200260200101518c600001518d602001518c612d3f565b90935091506111948285614436565b93506111a08386614436565b94508781815181106111b4576111b461446e565b6020026020010151606001511515896040015182815181106111d8576111d861446e565b60200260200101517fad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b848661120d9190614436565b8b858151811061121f5761121f61446e565b602002602001015160a001518c868151811061123d5761123d61446e565b602002602001015160c001518e60800151878151811061125f5761125f61446e565b602002602001015160405161127794939291906145be565b60405180910390a35b8061128a8161452c565b91505061109c565b5050336000908152600b6020526040902080548492506002906112ca9084906201000090046bffffffffffffffffffffffff16614436565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555080601260000160008282829054906101000a90046bffffffffffffffffffffffff166113249190614436565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060008f6001600381106113675761136761446e565b602002013560001c9050600060088264ffffffffff16901c9050886060015163ffffffff168163ffffffff1611156113e157601280547fffffffffffffffff00000000ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000063ffffffff8416021790555b505050505050505050505b50505050505050565b6114168686868680602001905181019061140f91906146a1565b8686611432565b505050505050565b611426612e32565b61142f81612eb3565b50565b61143a612e32565b601f86511115611476576040517f25d0209c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8360ff166000036114b3576040517fe77dba5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845186511415806114d257506114ca846003614823565b60ff16865111155b15611509576040517f1d2d1c5800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601254600e546bffffffffffffffffffffffff9091169060005b816bffffffffffffffffffffffff1681101561158b57611578600e828154811061154f5761154f61446e565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff168484612fa8565b50806115838161452c565b915050611523565b5060008060005b836bffffffffffffffffffffffff1681101561169457600d81815481106115bb576115bb61446e565b600091825260209091200154600e805473ffffffffffffffffffffffffffffffffffffffff909216945090829081106115f6576115f661446e565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff8681168452600c8352604080852080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001690559116808452600b90925290912080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905591508061168c8161452c565b915050611592565b506116a1600d6000613ade565b6116ad600e6000613ade565b604080516080810182526000808252602082018190529181018290526060810182905290805b8c51811015611b1657600c60008e83815181106116f2576116f261446e565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528101919091526040016000205460ff161561175d576040517f77cea0fa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168d82815181106117875761178761446e565b602002602001015173ffffffffffffffffffffffffffffffffffffffff16036117dc576040517f815e1d6400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405180604001604052806001151581526020018260ff16815250600c60008f848151811061180d5761180d61446e565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528181019290925260400160002082518154939092015160ff16610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909316929092171790558b518c90829081106118b5576118b561446e565b60200260200101519150600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603611925576040517f58a70a0a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff82166000908152600b60209081526040918290208251608081018452905460ff80821615801584526101008304909116938301939093526bffffffffffffffffffffffff6201000082048116948301949094526e010000000000000000000000000000900490921660608301529093506119e0576040517f6a7281ad00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001835260ff80821660208086019182526bffffffffffffffffffffffff808b166060880190815273ffffffffffffffffffffffffffffffffffffffff87166000908152600b909352604092839020885181549551948a0151925184166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff939094166201000002929092167fffffffffffff000000000000000000000000000000000000000000000000ffff94909616610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009095169490941717919091169290921791909117905580611b0e8161452c565b9150506116d3565b50508a51611b2c9150600d9060208d0190613afc565b508851611b4090600e9060208c0190613afc565b50604051806101600160405280856bffffffffffffffffffffffff168152602001886000015163ffffffff168152602001886020015163ffffffff168152602001600063ffffffff168152602001886060015162ffffff168152602001886080015161ffff1681526020018960ff1681526020016012600001601e9054906101000a900460ff16151581526020016012600001601f9054906101000a900460ff161515815260200188610200015115158152602001886101e0015173ffffffffffffffffffffffffffffffffffffffff16815250601260008201518160000160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550602082015181600001600c6101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160106101000a81548163ffffffff021916908363ffffffff16021790555060608201518160000160146101000a81548163ffffffff021916908363ffffffff16021790555060808201518160000160186101000a81548162ffffff021916908362ffffff16021790555060a082015181600001601b6101000a81548161ffff021916908361ffff16021790555060c082015181600001601d6101000a81548160ff021916908360ff16021790555060e082015181600001601e6101000a81548160ff02191690831515021790555061010082015181600001601f6101000a81548160ff0219169083151502179055506101208201518160010160006101000a81548160ff0219169083151502179055506101408201518160010160016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055509050506040518061018001604052808860a001516bffffffffffffffffffffffff16815260200188610180015173ffffffffffffffffffffffffffffffffffffffff168152602001601460010160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff168152602001886040015163ffffffff1681526020018860c0015163ffffffff168152602001601460010160149054906101000a900463ffffffff1663ffffffff168152602001601460010160189054906101000a900463ffffffff1663ffffffff1681526020016014600101601c9054906101000a900463ffffffff1663ffffffff1681526020018860e0015163ffffffff16815260200188610100015163ffffffff16815260200188610120015163ffffffff168152602001886101c0015173ffffffffffffffffffffffffffffffffffffffff16815250601460008201518160000160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550602082015181600001600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060408201518160010160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550606082015181600101600c6101000a81548163ffffffff021916908363ffffffff16021790555060808201518160010160106101000a81548163ffffffff021916908363ffffffff16021790555060a08201518160010160146101000a81548163ffffffff021916908363ffffffff16021790555060c08201518160010160186101000a81548163ffffffff021916908363ffffffff16021790555060e082015181600101601c6101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160020160006101000a81548163ffffffff021916908363ffffffff1602179055506101208201518160020160046101000a81548163ffffffff021916908363ffffffff1602179055506101408201518160020160086101000a81548163ffffffff021916908363ffffffff16021790555061016082015181600201600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505086610140015160178190555086610160015160188190555060006014600101601c9054906101000a900463ffffffff169050876101e0015173ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa1580156121e1573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061220591906144b6565b601580547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167c010000000000000000000000000000000000000000000000000000000063ffffffff9384160217808255600192601891612280918591780100000000000000000000000000000000000000000000000090041661483f565b92506101000a81548163ffffffff021916908363ffffffff1602179055506000886040516020016122b191906148ad565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291905260155490915061231a90469030907801000000000000000000000000000000000000000000000000900463ffffffff168f8f8f878f8f6131b0565b60115560005b61232a600961325a565b81101561235a5761234761233f600983613264565b600990613277565b50806123528161452c565b915050612320565b5060005b896101a00151518110156123b15761239e8a6101a0015182815181106123865761238661446e565b6020026020010151600961329990919063ffffffff16565b50806123a98161452c565b91505061235e565b507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0582601154601460010160189054906101000a900463ffffffff168f8f8f878f8f60405161240899989796959493929190614a51565b60405180910390a1505050505050505050505050565b3215612456576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60125460009081907f0100000000000000000000000000000000000000000000000000000000000000900460ff16156124bd576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601280547effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f01000000000000000000000000000000000000000000000000000000000000001790556040517f4585e33b0000000000000000000000000000000000000000000000000000000090612539908590602401613bff565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009094169390931790925290517f79188d1600000000000000000000000000000000000000000000000000000000815290935073ffffffffffffffffffffffffffffffffffffffff8616906379188d169061260c9087908790600401614ae7565b60408051808303816000875af115801561262a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061264e9190614b00565b601280547effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff16905590969095509350505050565b60008686604051612693929190614b2e565b6040519081900381206126aa918a90602001614b3e565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201208383019092526000808452908301819052909250906000805b87811015612882576001858783602081106127165761271661446e565b61272391901a601b61449d565b8b8b858181106127355761273561446e565b905060200201358a858151811061274e5761274e61446e565b60200260200101516040516000815260200160405260405161278c949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa1580156127ae573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff81166000908152600c602090815290849020838501909452925460ff808216151580855261010090920416938301939093529095509350905061285c576040517f0f4c073700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826020015160080260ff166001901b84019350808061287a9061452c565b9150506126f9565b50827e010101010101010101010101010101010101010101010101010101010101018416146128dd576040517fc103be2e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050505050505050505050565b6129236040518060c001604052806000815260200160008152602001606081526020016060815260200160608152602001606081525090565b600061293183850185614bd4565b604081015151606082015151919250908114158061295457508082608001515114155b806129645750808260a001515114155b1561299b576040517fb55ac75400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5090505b92915050565b6000818160045b600f811015612a32577fff0000000000000000000000000000000000000000000000000000000000000082168382602081106129ea576129ea61446e565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614612a2057506000949350505050565b80612a2a8161452c565b9150506129ac565b5081600f1a6001811115612a4857612a486144cf565b949350505050565b600080612a6288878b60c001516132bb565b9050600080612a7d8b8a63ffffffff16858a8a60018b613347565b9092509050612a8c8183614436565b9b9a5050505050505050505050565b600080808085608001516001811115612ab657612ab66144cf565b03612adc57612ac888888888886135d4565b612ad757600092509050612c10565b612b54565b600185608001516001811115612af457612af46144cf565b03612b22576000612b078989898861375d565b9250905080612b1c5750600092509050612c10565b50612b54565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b84516040015163ffffffff168710612ba957877fc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd563687604051612b969190613bff565b60405180910390a2600092509050612c10565b84604001516bffffffffffffffffffffffff16856000015160a001516bffffffffffffffffffffffff161015612c0957877f377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e0287604051612b969190613bff565b6001925090505b9550959350505050565b600081608001516001811115612c3257612c326144cf565b03612c9657600083815260046020526040902060010180547fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff16780100000000000000000000000000000000000000000000000063ffffffff851602179055505050565b600181608001516001811115612cae57612cae6144cf565b03612d1a5760e08101805160009081526008602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055915191517fa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f29190a25b505050565b6000612d2c8484846132bb565b905080851015612a485750929392505050565b600080612d5a888760a001518860c001518888886001613347565b90925090506000612d6b8284614436565b600089815260046020526040902060010180549192508291600c90612daf9084906c0100000000000000000000000090046bffffffffffffffffffffffff16614cc1565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915560008a815260046020526040812060010180548594509092612df891859116614436565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555050965096945050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314612456576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610396565b3373ffffffffffffffffffffffffffffffffffffffff821603612f32576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610396565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e01000000000000000000000000000090049091166060820152906131a45760008160600151856130409190614cc1565b9050600061304e8583614ce6565b905080836040018181516130629190614436565b6bffffffffffffffffffffffff1690525061307d8582614d11565b8360600181815161308e9190614436565b6bffffffffffffffffffffffff90811690915273ffffffffffffffffffffffffffffffffffffffff89166000908152600b602090815260409182902087518154928901519389015160608a015186166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff919096166201000002167fffffffffffff000000000000000000000000000000000000000000000000ffff60ff95909516610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909416939093171792909216179190911790555050505b60400151949350505050565b6000808a8a8a8a8a8a8a8a8a6040516020016131d499989796959493929190614d41565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179b9a5050505050505050505050565b600061299f825490565b6000613270838361396b565b9392505050565b60006132708373ffffffffffffffffffffffffffffffffffffffff8416613995565b60006132708373ffffffffffffffffffffffffffffffffffffffff8416613a8f565b600080808560018111156132d1576132d16144cf565b036132e0575062015f906132ff565b60018560018111156132f4576132f46144cf565b03612b2257506201b38c5b61331063ffffffff85166014614564565b61331b84600161449d565b61332a9060ff16611d4c614564565b613334908361445b565b61333e919061445b565b95945050505050565b60008060008960a0015161ffff16876133609190614564565b905083801561336e5750803a105b1561337657503a5b600084156133f9578a610140015173ffffffffffffffffffffffffffffffffffffffff166318b8f6136040518163ffffffff1660e01b8152600401602060405180830381865afa1580156133ce573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906133f291906144b6565b90506134b5565b6101408b01516016546040517f1254414000000000000000000000000000000000000000000000000000000000815264010000000090910463ffffffff16600482015273ffffffffffffffffffffffffffffffffffffffff90911690631254414090602401602060405180830381865afa15801561347b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061349f91906144b6565b8b60a0015161ffff166134b29190614564565b90505b6134c361ffff8716826145aa565b9050600087826134d38c8e61445b565b6134dd9086614564565b6134e7919061445b565b6134f990670de0b6b3a7640000614564565b61350391906145aa565b905060008c6040015163ffffffff1664e8d4a510006135229190614564565b898e6020015163ffffffff16858f8861353b9190614564565b613545919061445b565b61355390633b9aca00614564565b61355d9190614564565b61356791906145aa565b613571919061445b565b90506b033b2e3c9fd0803ce800000061358a828461445b565b11156135c2576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b600080848060200190518101906135eb9190614dd6565b845160c00151815191925063ffffffff9081169116101561364857867f405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8866040516136369190613bff565b60405180910390a2600091505061333e565b826101200151801561370957506020810151158015906137095750602081015161014084015182516040517f85df51fd00000000000000000000000000000000000000000000000000000000815263ffffffff909116600482015273ffffffffffffffffffffffffffffffffffffffff909116906385df51fd90602401602060405180830381865afa1580156136e2573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061370691906144b6565b14155b8061371b5750805163ffffffff168611155b1561375057867f6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301866040516136369190613bff565b5060019695505050505050565b6000806000848060200190518101906137769190614e2e565b90506000878260000151836020015184604001516040516020016137d894939291909384526020840192909252604083015260e01b7fffffffff0000000000000000000000000000000000000000000000000000000016606082015260640190565b60405160208183030381529060405280519060200120905084610120015180156138b457506080820151158015906138b45750608082015161014086015160608401516040517f85df51fd00000000000000000000000000000000000000000000000000000000815263ffffffff909116600482015273ffffffffffffffffffffffffffffffffffffffff909116906385df51fd90602401602060405180830381865afa15801561388d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906138b191906144b6565b14155b806138c9575086826060015163ffffffff1610155b1561391357877f6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301876040516138fe9190613bff565b60405180910390a26000935091506139629050565b60008181526008602052604090205460ff161561395a57877f405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8876040516138fe9190613bff565b600193509150505b94509492505050565b60008260000182815481106139825761398261446e565b9060005260206000200154905092915050565b60008181526001830160205260408120548015613a7e5760006139b9600183614519565b85549091506000906139cd90600190614519565b9050818114613a325760008660000182815481106139ed576139ed61446e565b9060005260206000200154905080876000018481548110613a1057613a1061446e565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080613a4357613a43614eaf565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505061299f565b600091505061299f565b5092915050565b6000818152600183016020526040812054613ad65750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915561299f565b50600061299f565b508054600082559060005260206000209081019061142f9190613b86565b828054828255906000526020600020908101928215613b76579160200282015b82811115613b7657825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190613b1c565b50613b82929150613b86565b5090565b5b80821115613b825760008155600101613b87565b6000815180845260005b81811015613bc157602081850181015186830182015201613ba5565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815260006132706020830184613b9b565b73ffffffffffffffffffffffffffffffffffffffff8116811461142f57600080fd5b8035613c3f81613c12565b919050565b60008083601f840112613c5657600080fd5b50813567ffffffffffffffff811115613c6e57600080fd5b602083019150836020828501011115613c8657600080fd5b9250929050565b60008060008060608587031215613ca357600080fd5b8435613cae81613c12565b935060208501359250604085013567ffffffffffffffff811115613cd157600080fd5b613cdd87828801613c44565b95989497509550505050565b600080600060408486031215613cfe57600080fd5b83359250602084013567ffffffffffffffff811115613d1c57600080fd5b613d2886828701613c44565b9497909650939450505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610220810167ffffffffffffffff81118282101715613d8857613d88613d35565b60405290565b60405160c0810167ffffffffffffffff81118282101715613d8857613d88613d35565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613df857613df8613d35565b604052919050565b600067ffffffffffffffff821115613e1a57613e1a613d35565b5060051b60200190565b600082601f830112613e3557600080fd5b81356020613e4a613e4583613e00565b613db1565b82815260059290921b84018101918181019086841115613e6957600080fd5b8286015b84811015613e845780358352918301918301613e6d565b509695505050505050565b600080600080600080600060e0888a031215613eaa57600080fd5b6060880189811115613ebb57600080fd5b8897503567ffffffffffffffff80821115613ed557600080fd5b613ee18b838c01613c44565b909850965060808a0135915080821115613efa57600080fd5b818a0191508a601f830112613f0e57600080fd5b813581811115613f1d57600080fd5b8b60208260051b8501011115613f3257600080fd5b6020830196508095505060a08a0135915080821115613f5057600080fd5b50613f5d8a828b01613e24565b92505060c0880135905092959891949750929550565b600082601f830112613f8457600080fd5b81356020613f94613e4583613e00565b82815260059290921b84018101918181019086841115613fb357600080fd5b8286015b84811015613e84578035613fca81613c12565b8352918301918301613fb7565b803560ff81168114613c3f57600080fd5b600082601f830112613ff957600080fd5b813567ffffffffffffffff81111561401357614013613d35565b61404460207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601613db1565b81815284602083860101111561405957600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff81168114613c3f57600080fd5b60008060008060008060c087890312156140a757600080fd5b863567ffffffffffffffff808211156140bf57600080fd5b6140cb8a838b01613f73565b975060208901359150808211156140e157600080fd5b6140ed8a838b01613f73565b96506140fb60408a01613fd7565b9550606089013591508082111561411157600080fd5b61411d8a838b01613fe8565b945061412b60808a01614076565b935060a089013591508082111561414157600080fd5b5061414e89828a01613fe8565b9150509295509295509295565b60006020828403121561416d57600080fd5b813561327081613c12565b63ffffffff8116811461142f57600080fd5b8035613c3f81614178565b62ffffff8116811461142f57600080fd5b8035613c3f81614195565b61ffff8116811461142f57600080fd5b8035613c3f816141b1565b6bffffffffffffffffffffffff8116811461142f57600080fd5b8035613c3f816141cc565b801515811461142f57600080fd5b8035613c3f816141f1565b6000610220828403121561421d57600080fd5b614225613d64565b90506142308261418a565b815261423e6020830161418a565b602082015261424f6040830161418a565b6040820152614260606083016141a6565b6060820152614271608083016141c1565b608082015261428260a083016141e6565b60a082015261429360c0830161418a565b60c08201526142a460e0830161418a565b60e08201526101006142b781840161418a565b908201526101206142c983820161418a565b90820152610140828101359082015261016080830135908201526101806142f1818401613c34565b908201526101a08281013567ffffffffffffffff81111561431157600080fd5b61431d85828601613f73565b8284015250506101c0614331818401613c34565b908201526101e0614343838201613c34565b908201526102006143558382016141ff565b9082015292915050565b60008060008060008060c0878903121561437857600080fd5b863567ffffffffffffffff8082111561439057600080fd5b61439c8a838b01613f73565b975060208901359150808211156143b257600080fd5b6143be8a838b01613f73565b96506143cc60408a01613fd7565b955060608901359150808211156143e257600080fd5b61411d8a838b0161420a565b60006020828403121561440057600080fd5b5035919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6bffffffffffffffffffffffff818116838216019080821115613a8857613a88614407565b8082018082111561299f5761299f614407565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60ff818116838216019081111561299f5761299f614407565b6000602082840312156144c857600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b61ffff818116838216019080821115613a8857613a88614407565b8181038181111561299f5761299f614407565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361455d5761455d614407565b5060010190565b808202811582820484141761299f5761299f614407565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000826145b9576145b961457b565b500490565b6bffffffffffffffffffffffff851681528360208201528260408201526080606082015260006145f16080830184613b9b565b9695505050505050565b8051613c3f81614178565b8051613c3f81614195565b8051613c3f816141b1565b8051613c3f816141cc565b8051613c3f81613c12565b600082601f83011261464357600080fd5b81516020614653613e4583613e00565b82815260059290921b8401810191818101908684111561467257600080fd5b8286015b84811015613e8457805161468981613c12565b8352918301918301614676565b8051613c3f816141f1565b6000602082840312156146b357600080fd5b815167ffffffffffffffff808211156146cb57600080fd5b9083019061022082860312156146e057600080fd5b6146e8613d64565b6146f1836145fb565b81526146ff602084016145fb565b6020820152614710604084016145fb565b604082015261472160608401614606565b606082015261473260808401614611565b608082015261474360a0840161461c565b60a082015261475460c084016145fb565b60c082015261476560e084016145fb565b60e08201526101006147788185016145fb565b9082015261012061478a8482016145fb565b90820152610140838101519082015261016080840151908201526101806147b2818501614627565b908201526101a083810151838111156147ca57600080fd5b6147d688828701614632565b8284015250506101c091506147ec828401614627565b828201526101e09150614800828401614627565b828201526102009150614814828401614696565b91810191909152949350505050565b60ff8181168382160290811690818114613a8857613a88614407565b63ffffffff818116838216019080821115613a8857613a88614407565b600081518084526020808501945080840160005b838110156148a257815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101614870565b509495945050505050565b602081526148c460208201835163ffffffff169052565b600060208301516148dd604084018263ffffffff169052565b50604083015163ffffffff8116606084015250606083015162ffffff8116608084015250608083015161ffff811660a08401525060a08301516bffffffffffffffffffffffff811660c08401525060c083015163ffffffff811660e08401525060e08301516101006149568185018363ffffffff169052565b840151905061012061496f8482018363ffffffff169052565b84015190506101406149888482018363ffffffff169052565b840151610160848101919091528401516101808085019190915284015190506101a06149cb8185018373ffffffffffffffffffffffffffffffffffffffff169052565b808501519150506102206101c081818601526149eb61024086018461485c565b908601519092506101e0614a168682018373ffffffffffffffffffffffffffffffffffffffff169052565b8601519050610200614a3f8682018373ffffffffffffffffffffffffffffffffffffffff169052565b90950151151593019290925250919050565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152614a818184018a61485c565b90508281036080840152614a95818961485c565b905060ff871660a084015282810360c0840152614ab28187613b9b565b905067ffffffffffffffff851660e0840152828103610100840152614ad78185613b9b565b9c9b505050505050505050505050565b828152604060208201526000612a486040830184613b9b565b60008060408385031215614b1357600080fd5b8251614b1e816141f1565b6020939093015192949293505050565b8183823760009101908152919050565b8281526080810160608360208401379392505050565b600082601f830112614b6557600080fd5b81356020614b75613e4583613e00565b82815260059290921b84018101918181019086841115614b9457600080fd5b8286015b84811015613e8457803567ffffffffffffffff811115614bb85760008081fd5b614bc68986838b0101613fe8565b845250918301918301614b98565b600060208284031215614be657600080fd5b813567ffffffffffffffff80821115614bfe57600080fd5b9083019060c08286031215614c1257600080fd5b614c1a613d8e565b8235815260208301356020820152604083013582811115614c3a57600080fd5b614c4687828601613e24565b604083015250606083013582811115614c5e57600080fd5b614c6a87828601613e24565b606083015250608083013582811115614c8257600080fd5b614c8e87828601614b54565b60808301525060a083013582811115614ca657600080fd5b614cb287828601614b54565b60a08301525095945050505050565b6bffffffffffffffffffffffff828116828216039080821115613a8857613a88614407565b60006bffffffffffffffffffffffff80841680614d0557614d0561457b565b92169190910492915050565b6bffffffffffffffffffffffff818116838216028082169190828114614d3957614d39614407565b505092915050565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b166040850152816060850152614d888285018b61485c565b91508382036080850152614d9c828a61485c565b915060ff881660a085015283820360c0850152614db98288613b9b565b90861660e08501528381036101008501529050614ad78185613b9b565b600060408284031215614de857600080fd5b6040516040810181811067ffffffffffffffff82111715614e0b57614e0b613d35565b6040528251614e1981614178565b81526020928301519281019290925250919050565b600060a08284031215614e4057600080fd5b60405160a0810181811067ffffffffffffffff82111715614e6357614e63613d35565b806040525082518152602083015160208201526040830151614e8481614178565b60408201526060830151614e9781614178565b60608201526080928301519281019290925250919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000813000a",
>>>>>>> 6d733192b6 (address some comments)
=======
	Bin: "0x6101206040523480156200001257600080fd5b506040516200529938038062005299833981016040819052620000359162000345565b80816001600160a01b031663ca30e6036040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200009b919062000345565b826001600160a01b031663b10b673c6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000da573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000100919062000345565b836001600160a01b0316636709d0e56040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200013f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000165919062000345565b846001600160a01b0316635425d8ac6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca919062000345565b3380600081620002215760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200025457620002548162000281565b5050506001600160a01b0393841660805291831660a052821660c052811660e0521661010052506200036c565b336001600160a01b03821603620002db5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000218565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200034257600080fd5b50565b6000602082840312156200035857600080fd5b815162000365816200032c565b9392505050565b60805160a05160c05160e05161010051614eeb620003ae6000396000818160d6015261016f015260005050600050506000505060006104330152614eeb6000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c8063aed2e92911610081578063e3d0e7121161005b578063e3d0e712146102e0578063f2fde38b146102f3578063f75f6b1114610306576100d4565b8063aed2e92914610262578063afcb95d71461028c578063b1dc65a4146102cd576100d4565b806381ff7048116100b257806381ff7048146101bc5780638da5cb5b14610231578063a4c0ed361461024f576100d4565b8063181f5a771461011b578063349e8cca1461016d57806379ba5097146101b4575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e808015610114573d6000f35b3d6000fd5b005b6101576040518060400160405280601881526020017f4175746f6d6174696f6e526567697374727920322e322e30000000000000000081525081565b6040516101649190613bff565b60405180910390f35b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610164565b610119610319565b61020e60155460115463ffffffff780100000000000000000000000000000000000000000000000083048116937c01000000000000000000000000000000000000000000000000000000009093041691565b6040805163ffffffff948516815293909216602084015290820152606001610164565b60005473ffffffffffffffffffffffffffffffffffffffff1661018f565b61011961025d366004613c8d565b61041b565b610275610270366004613ce9565b610637565b604080519215158352602083019190915201610164565b601154601254604080516000815260208101939093527401000000000000000000000000000000000000000090910463ffffffff1690820152606001610164565b6101196102db366004613e8f565b6107ad565b6101196102ee36600461408e565b6113f5565b61011961030136600461415b565b61141e565b61011961031436600461435f565b611432565b60015473ffffffffffffffffffffffffffffffffffffffff16331461039f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461048a576040517fc8bad78d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b602081146104c4576040517fdfe9309000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006104d2828401846143ee565b60008181526004602052604090205490915065010000000000900463ffffffff9081161461052c576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000818152600460205260409020600101546105679085906c0100000000000000000000000090046bffffffffffffffffffffffff16614436565b600082815260046020526040902060010180546bffffffffffffffffffffffff929092166c01000000000000000000000000027fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff9092169190911790556019546105d290859061445b565b6019556040516bffffffffffffffffffffffff8516815273ffffffffffffffffffffffffffffffffffffffff86169082907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a35050505050565b60008061064261241e565b6012547e01000000000000000000000000000000000000000000000000000000000000900460ff16156106a1576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600085815260046020908152604091829020825160e081018452815460ff811615158252610100810463ffffffff908116838601819052650100000000008304821684880152690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff16606084018190526001909401546bffffffffffffffffffffffff80821660808601526c0100000000000000000000000082041660a0850152780100000000000000000000000000000000000000000000000090041660c08301528451601f890185900485028101850190955287855290936107a093899089908190840183828082843760009201919091525061245892505050565b9097909650945050505050565b60005a60408051610160810182526012546bffffffffffffffffffffffff8116825263ffffffff6c010000000000000000000000008204811660208401527001000000000000000000000000000000008204811693830193909352740100000000000000000000000000000000000000008104909216606082015262ffffff7801000000000000000000000000000000000000000000000000830416608082015261ffff7b0100000000000000000000000000000000000000000000000000000083041660a082015260ff7d0100000000000000000000000000000000000000000000000000000000008304811660c08301527e010000000000000000000000000000000000000000000000000000000000008304811615801560e08401527f01000000000000000000000000000000000000000000000000000000000000009093048116151561010080840191909152601354918216151561012084015273ffffffffffffffffffffffffffffffffffffffff910416610140820152919250610963576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600b602052604090205460ff166109ac576040517f1099ed7500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6011548935146109e8576040517fdfdcf8e700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60c08101516109f890600161449d565b60ff1685141580610a0a575083518514155b15610a41576040517f0244f71a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610a5089898989898989612681565b6000610a5c89896128ea565b9050600081604001515167ffffffffffffffff811115610a7e57610a7e613d35565b604051908082528060200260200182016040528015610b4257816020015b604080516101e0810182526000610100820181815261012083018290526101408301829052610160830182905261018083018290526101a083018290526101c0830182905282526020808301829052928201819052606082018190526080820181905260a0820181905260c0820181905260e082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181610a9c5790505b50905060008084610140015173ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa158015610b98573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bbc91906144b6565b905060005b846040015151811015611006576004600086604001518381518110610be857610be861446e565b6020908102919091018101518252818101929092526040908101600020815160e081018352815460ff811615158252610100810463ffffffff90811695830195909552650100000000008104851693820193909352690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff166060830152600101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490911660c08201528451859083908110610ccd57610ccd61446e565b602002602001015160000181905250610d0285604001518281518110610cf557610cf561446e565b60200260200101516129a5565b848281518110610d1457610d1461446e565b6020026020010151608001906001811115610d3157610d316144cf565b90816001811115610d4457610d446144cf565b81525050610db886858381518110610d5e57610d5e61446e565b60200260200101516080015187606001518481518110610d8057610d8061446e565b60200260200101518860a001518581518110610d9e57610d9e61446e565b60200260200101515189600001518a602001516001612a50565b848281518110610dca57610dca61446e565b6020026020010151604001906bffffffffffffffffffffffff1690816bffffffffffffffffffffffff1681525050610e5885604001518281518110610e1157610e1161446e565b60200260200101518387608001518481518110610e3057610e3061446e565b6020026020010151878581518110610e4a57610e4a61446e565b60200260200101518a612a9b565b858381518110610e6a57610e6a61446e565b6020026020010151602001868481518110610e8757610e8761446e565b602002602001015160e0018281525082151515158152505050838181518110610eb257610eb261446e565b60200260200101516020015115610ed557610ece6001846144fe565b9250610eda565b610ff4565b610f40848281518110610eef57610eef61446e565b6020026020010151600001516060015186606001518381518110610f1557610f1561446e565b60200260200101518760a001518481518110610f3357610f3361446e565b6020026020010151612458565b858381518110610f5257610f5261446e565b6020026020010151606001868481518110610f6f57610f6f61446e565b602002602001015160a0018281525082151515158152505050838181518110610f9a57610f9a61446e565b602002602001015160a0015187610fb19190614519565b9650610ff485604001518281518110610fcc57610fcc61446e565b602002602001015183868481518110610fe757610fe761446e565b6020026020010151612c1a565b80610ffe8161452c565b915050610bc1565b508161ffff1660000361101e575050505050506113ec565b60c085015161102e90600161449d565b61103d9060ff1661044c614564565b616dc461104b8d6010614564565b5a611056908a614519565b611060919061445b565b61106a919061445b565b611074919061445b565b9550611b5861108761ffff8416886145aa565b611091919061445b565b955060008060008060005b886040015151811015611292578781815181106110bb576110bb61446e565b60200260200101516020015115611280576111178b8983815181106110e2576110e261446e565b6020026020010151608001518b60a0015184815181106111045761110461446e565b6020026020010151518d60c00151612d1f565b8882815181106111295761112961446e565b602002602001015160c00181815250506111858a8a6040015183815181106111535761115361446e565b60200260200101518a848151811061116d5761116d61446e565b60200260200101518c600001518d602001518c612d3f565b90935091506111948285614436565b93506111a08386614436565b94508781815181106111b4576111b461446e565b6020026020010151606001511515896040015182815181106111d8576111d861446e565b60200260200101517fad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b848661120d9190614436565b8b858151811061121f5761121f61446e565b602002602001015160a001518c868151811061123d5761123d61446e565b602002602001015160c001518e60800151878151811061125f5761125f61446e565b602002602001015160405161127794939291906145be565b60405180910390a35b8061128a8161452c565b91505061109c565b5050336000908152600b6020526040902080548492506002906112ca9084906201000090046bffffffffffffffffffffffff16614436565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555080601260000160008282829054906101000a90046bffffffffffffffffffffffff166113249190614436565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060008f6001600381106113675761136761446e565b602002013560001c9050600060088264ffffffffff16901c9050886060015163ffffffff168163ffffffff1611156113e157601280547fffffffffffffffff00000000ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000063ffffffff8416021790555b505050505050505050505b50505050505050565b6114168686868680602001905181019061140f91906146a1565b8686611432565b505050505050565b611426612e32565b61142f81612eb3565b50565b61143a612e32565b601f86511115611476576040517f25d0209c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8360ff166000036114b3576040517fe77dba5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845186511415806114d257506114ca846003614823565b60ff16865111155b15611509576040517f1d2d1c5800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601254600e546bffffffffffffffffffffffff9091169060005b816bffffffffffffffffffffffff1681101561158b57611578600e828154811061154f5761154f61446e565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff168484612fa8565b50806115838161452c565b915050611523565b5060008060005b836bffffffffffffffffffffffff1681101561169457600d81815481106115bb576115bb61446e565b600091825260209091200154600e805473ffffffffffffffffffffffffffffffffffffffff909216945090829081106115f6576115f661446e565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff8681168452600c8352604080852080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001690559116808452600b90925290912080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905591508061168c8161452c565b915050611592565b506116a1600d6000613ade565b6116ad600e6000613ade565b604080516080810182526000808252602082018190529181018290526060810182905290805b8c51811015611b1657600c60008e83815181106116f2576116f261446e565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528101919091526040016000205460ff161561175d576040517f77cea0fa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168d82815181106117875761178761446e565b602002602001015173ffffffffffffffffffffffffffffffffffffffff16036117dc576040517f815e1d6400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405180604001604052806001151581526020018260ff16815250600c60008f848151811061180d5761180d61446e565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528181019290925260400160002082518154939092015160ff16610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909316929092171790558b518c90829081106118b5576118b561446e565b60200260200101519150600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603611925576040517f58a70a0a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff82166000908152600b60209081526040918290208251608081018452905460ff80821615801584526101008304909116938301939093526bffffffffffffffffffffffff6201000082048116948301949094526e010000000000000000000000000000900490921660608301529093506119e0576040517f6a7281ad00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001835260ff80821660208086019182526bffffffffffffffffffffffff808b166060880190815273ffffffffffffffffffffffffffffffffffffffff87166000908152600b909352604092839020885181549551948a0151925184166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff939094166201000002929092167fffffffffffff000000000000000000000000000000000000000000000000ffff94909616610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009095169490941717919091169290921791909117905580611b0e8161452c565b9150506116d3565b50508a51611b2c9150600d9060208d0190613afc565b508851611b4090600e9060208c0190613afc565b50604051806101600160405280856bffffffffffffffffffffffff168152602001886000015163ffffffff168152602001886020015163ffffffff168152602001600063ffffffff168152602001886060015162ffffff168152602001886080015161ffff1681526020018960ff1681526020016012600001601e9054906101000a900460ff16151581526020016012600001601f9054906101000a900460ff161515815260200188610200015115158152602001886101e0015173ffffffffffffffffffffffffffffffffffffffff16815250601260008201518160000160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550602082015181600001600c6101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160106101000a81548163ffffffff021916908363ffffffff16021790555060608201518160000160146101000a81548163ffffffff021916908363ffffffff16021790555060808201518160000160186101000a81548162ffffff021916908362ffffff16021790555060a082015181600001601b6101000a81548161ffff021916908361ffff16021790555060c082015181600001601d6101000a81548160ff021916908360ff16021790555060e082015181600001601e6101000a81548160ff02191690831515021790555061010082015181600001601f6101000a81548160ff0219169083151502179055506101208201518160010160006101000a81548160ff0219169083151502179055506101408201518160010160016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055509050506040518061018001604052808860a001516bffffffffffffffffffffffff16815260200188610180015173ffffffffffffffffffffffffffffffffffffffff168152602001601460010160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff168152602001886040015163ffffffff1681526020018860c0015163ffffffff168152602001601460010160149054906101000a900463ffffffff1663ffffffff168152602001601460010160189054906101000a900463ffffffff1663ffffffff1681526020016014600101601c9054906101000a900463ffffffff1663ffffffff1681526020018860e0015163ffffffff16815260200188610100015163ffffffff16815260200188610120015163ffffffff168152602001886101c0015173ffffffffffffffffffffffffffffffffffffffff16815250601460008201518160000160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550602082015181600001600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060408201518160010160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550606082015181600101600c6101000a81548163ffffffff021916908363ffffffff16021790555060808201518160010160106101000a81548163ffffffff021916908363ffffffff16021790555060a08201518160010160146101000a81548163ffffffff021916908363ffffffff16021790555060c08201518160010160186101000a81548163ffffffff021916908363ffffffff16021790555060e082015181600101601c6101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160020160006101000a81548163ffffffff021916908363ffffffff1602179055506101208201518160020160046101000a81548163ffffffff021916908363ffffffff1602179055506101408201518160020160086101000a81548163ffffffff021916908363ffffffff16021790555061016082015181600201600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505086610140015160178190555086610160015160188190555060006014600101601c9054906101000a900463ffffffff169050876101e0015173ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa1580156121e1573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061220591906144b6565b601580547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167c010000000000000000000000000000000000000000000000000000000063ffffffff9384160217808255600192601891612280918591780100000000000000000000000000000000000000000000000090041661483f565b92506101000a81548163ffffffff021916908363ffffffff1602179055506000886040516020016122b191906148ad565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291905260155490915061231a90469030907801000000000000000000000000000000000000000000000000900463ffffffff168f8f8f878f8f6131b0565b60115560005b61232a600961325a565b81101561235a5761234761233f600983613264565b600990613277565b50806123528161452c565b915050612320565b5060005b896101a00151518110156123b15761239e8a6101a0015182815181106123865761238661446e565b6020026020010151600961329990919063ffffffff16565b50806123a98161452c565b91505061235e565b507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0582601154601460010160189054906101000a900463ffffffff168f8f8f878f8f60405161240899989796959493929190614a51565b60405180910390a1505050505050505050505050565b3215612456576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60125460009081907f0100000000000000000000000000000000000000000000000000000000000000900460ff16156124bd576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601280547effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f01000000000000000000000000000000000000000000000000000000000000001790556040517f4585e33b0000000000000000000000000000000000000000000000000000000090612539908590602401613bff565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009094169390931790925290517f79188d1600000000000000000000000000000000000000000000000000000000815290935073ffffffffffffffffffffffffffffffffffffffff8616906379188d169061260c9087908790600401614ae7565b60408051808303816000875af115801561262a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061264e9190614b00565b601280547effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff16905590969095509350505050565b60008686604051612693929190614b2e565b6040519081900381206126aa918a90602001614b3e565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201208383019092526000808452908301819052909250906000805b87811015612882576001858783602081106127165761271661446e565b61272391901a601b61449d565b8b8b858181106127355761273561446e565b905060200201358a858151811061274e5761274e61446e565b60200260200101516040516000815260200160405260405161278c949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa1580156127ae573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff81166000908152600c602090815290849020838501909452925460ff808216151580855261010090920416938301939093529095509350905061285c576040517f0f4c073700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826020015160080260ff166001901b84019350808061287a9061452c565b9150506126f9565b50827e010101010101010101010101010101010101010101010101010101010101018416146128dd576040517fc103be2e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050505050505050505050565b6129236040518060c001604052806000815260200160008152602001606081526020016060815260200160608152602001606081525090565b600061293183850185614bd4565b604081015151606082015151919250908114158061295457508082608001515114155b806129645750808260a001515114155b1561299b576040517fb55ac75400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5090505b92915050565b6000818160045b600f811015612a32577fff0000000000000000000000000000000000000000000000000000000000000082168382602081106129ea576129ea61446e565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614612a2057506000949350505050565b80612a2a8161452c565b9150506129ac565b5081600f1a6001811115612a4857612a486144cf565b949350505050565b600080612a6288878b60c001516132bb565b9050600080612a7d8b8a63ffffffff16858a8a60018b613347565b9092509050612a8c8183614436565b9b9a5050505050505050505050565b600080808085608001516001811115612ab657612ab66144cf565b03612adc57612ac888888888886135d4565b612ad757600092509050612c10565b612b54565b600185608001516001811115612af457612af46144cf565b03612b22576000612b078989898861375d565b9250905080612b1c5750600092509050612c10565b50612b54565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b84516040015163ffffffff168710612ba957877fc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd563687604051612b969190613bff565b60405180910390a2600092509050612c10565b84604001516bffffffffffffffffffffffff16856000015160a001516bffffffffffffffffffffffff161015612c0957877f377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e0287604051612b969190613bff565b6001925090505b9550959350505050565b600081608001516001811115612c3257612c326144cf565b03612c9657600083815260046020526040902060010180547fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff16780100000000000000000000000000000000000000000000000063ffffffff851602179055505050565b600181608001516001811115612cae57612cae6144cf565b03612d1a5760e08101805160009081526008602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055915191517fa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f29190a25b505050565b6000612d2c8484846132bb565b905080851015612a485750929392505050565b600080612d5a888760a001518860c001518888886001613347565b90925090506000612d6b8284614436565b600089815260046020526040902060010180549192508291600c90612daf9084906c0100000000000000000000000090046bffffffffffffffffffffffff16614cc1565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915560008a815260046020526040812060010180548594509092612df891859116614436565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555050965096945050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314612456576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610396565b3373ffffffffffffffffffffffffffffffffffffffff821603612f32576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610396565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e01000000000000000000000000000090049091166060820152906131a45760008160600151856130409190614cc1565b9050600061304e8583614ce6565b905080836040018181516130629190614436565b6bffffffffffffffffffffffff1690525061307d8582614d11565b8360600181815161308e9190614436565b6bffffffffffffffffffffffff90811690915273ffffffffffffffffffffffffffffffffffffffff89166000908152600b602090815260409182902087518154928901519389015160608a015186166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff919096166201000002167fffffffffffff000000000000000000000000000000000000000000000000ffff60ff95909516610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909416939093171792909216179190911790555050505b60400151949350505050565b6000808a8a8a8a8a8a8a8a8a6040516020016131d499989796959493929190614d41565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179b9a5050505050505050505050565b600061299f825490565b6000613270838361396b565b9392505050565b60006132708373ffffffffffffffffffffffffffffffffffffffff8416613995565b60006132708373ffffffffffffffffffffffffffffffffffffffff8416613a8f565b600080808560018111156132d1576132d16144cf565b036132e0575062015f906132ff565b60018560018111156132f4576132f46144cf565b03612b2257506201ae145b61331063ffffffff85166014614564565b61331b84600161449d565b61332a9060ff16611d4c614564565b613334908361445b565b61333e919061445b565b95945050505050565b60008060008960a0015161ffff16876133609190614564565b905083801561336e5750803a105b1561337657503a5b600084156133f9578a610140015173ffffffffffffffffffffffffffffffffffffffff166318b8f6136040518163ffffffff1660e01b8152600401602060405180830381865afa1580156133ce573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906133f291906144b6565b90506134b5565b6101408b01516016546040517f1254414000000000000000000000000000000000000000000000000000000000815264010000000090910463ffffffff16600482015273ffffffffffffffffffffffffffffffffffffffff90911690631254414090602401602060405180830381865afa15801561347b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061349f91906144b6565b8b60a0015161ffff166134b29190614564565b90505b6134c361ffff8716826145aa565b9050600087826134d38c8e61445b565b6134dd9086614564565b6134e7919061445b565b6134f990670de0b6b3a7640000614564565b61350391906145aa565b905060008c6040015163ffffffff1664e8d4a510006135229190614564565b898e6020015163ffffffff16858f8861353b9190614564565b613545919061445b565b61355390633b9aca00614564565b61355d9190614564565b61356791906145aa565b613571919061445b565b90506b033b2e3c9fd0803ce800000061358a828461445b565b11156135c2576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b600080848060200190518101906135eb9190614dd6565b845160c00151815191925063ffffffff9081169116101561364857867f405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8866040516136369190613bff565b60405180910390a2600091505061333e565b826101200151801561370957506020810151158015906137095750602081015161014084015182516040517f85df51fd00000000000000000000000000000000000000000000000000000000815263ffffffff909116600482015273ffffffffffffffffffffffffffffffffffffffff909116906385df51fd90602401602060405180830381865afa1580156136e2573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061370691906144b6565b14155b8061371b5750805163ffffffff168611155b1561375057867f6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301866040516136369190613bff565b5060019695505050505050565b6000806000848060200190518101906137769190614e2e565b90506000878260000151836020015184604001516040516020016137d894939291909384526020840192909252604083015260e01b7fffffffff0000000000000000000000000000000000000000000000000000000016606082015260640190565b60405160208183030381529060405280519060200120905084610120015180156138b457506080820151158015906138b45750608082015161014086015160608401516040517f85df51fd00000000000000000000000000000000000000000000000000000000815263ffffffff909116600482015273ffffffffffffffffffffffffffffffffffffffff909116906385df51fd90602401602060405180830381865afa15801561388d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906138b191906144b6565b14155b806138c9575086826060015163ffffffff1610155b1561391357877f6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301876040516138fe9190613bff565b60405180910390a26000935091506139629050565b60008181526008602052604090205460ff161561395a57877f405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8876040516138fe9190613bff565b600193509150505b94509492505050565b60008260000182815481106139825761398261446e565b9060005260206000200154905092915050565b60008181526001830160205260408120548015613a7e5760006139b9600183614519565b85549091506000906139cd90600190614519565b9050818114613a325760008660000182815481106139ed576139ed61446e565b9060005260206000200154905080876000018481548110613a1057613a1061446e565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080613a4357613a43614eaf565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505061299f565b600091505061299f565b5092915050565b6000818152600183016020526040812054613ad65750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915561299f565b50600061299f565b508054600082559060005260206000209081019061142f9190613b86565b828054828255906000526020600020908101928215613b76579160200282015b82811115613b7657825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190613b1c565b50613b82929150613b86565b5090565b5b80821115613b825760008155600101613b87565b6000815180845260005b81811015613bc157602081850181015186830182015201613ba5565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815260006132706020830184613b9b565b73ffffffffffffffffffffffffffffffffffffffff8116811461142f57600080fd5b8035613c3f81613c12565b919050565b60008083601f840112613c5657600080fd5b50813567ffffffffffffffff811115613c6e57600080fd5b602083019150836020828501011115613c8657600080fd5b9250929050565b60008060008060608587031215613ca357600080fd5b8435613cae81613c12565b935060208501359250604085013567ffffffffffffffff811115613cd157600080fd5b613cdd87828801613c44565b95989497509550505050565b600080600060408486031215613cfe57600080fd5b83359250602084013567ffffffffffffffff811115613d1c57600080fd5b613d2886828701613c44565b9497909650939450505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610220810167ffffffffffffffff81118282101715613d8857613d88613d35565b60405290565b60405160c0810167ffffffffffffffff81118282101715613d8857613d88613d35565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613df857613df8613d35565b604052919050565b600067ffffffffffffffff821115613e1a57613e1a613d35565b5060051b60200190565b600082601f830112613e3557600080fd5b81356020613e4a613e4583613e00565b613db1565b82815260059290921b84018101918181019086841115613e6957600080fd5b8286015b84811015613e845780358352918301918301613e6d565b509695505050505050565b600080600080600080600060e0888a031215613eaa57600080fd5b6060880189811115613ebb57600080fd5b8897503567ffffffffffffffff80821115613ed557600080fd5b613ee18b838c01613c44565b909850965060808a0135915080821115613efa57600080fd5b818a0191508a601f830112613f0e57600080fd5b813581811115613f1d57600080fd5b8b60208260051b8501011115613f3257600080fd5b6020830196508095505060a08a0135915080821115613f5057600080fd5b50613f5d8a828b01613e24565b92505060c0880135905092959891949750929550565b600082601f830112613f8457600080fd5b81356020613f94613e4583613e00565b82815260059290921b84018101918181019086841115613fb357600080fd5b8286015b84811015613e84578035613fca81613c12565b8352918301918301613fb7565b803560ff81168114613c3f57600080fd5b600082601f830112613ff957600080fd5b813567ffffffffffffffff81111561401357614013613d35565b61404460207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601613db1565b81815284602083860101111561405957600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff81168114613c3f57600080fd5b60008060008060008060c087890312156140a757600080fd5b863567ffffffffffffffff808211156140bf57600080fd5b6140cb8a838b01613f73565b975060208901359150808211156140e157600080fd5b6140ed8a838b01613f73565b96506140fb60408a01613fd7565b9550606089013591508082111561411157600080fd5b61411d8a838b01613fe8565b945061412b60808a01614076565b935060a089013591508082111561414157600080fd5b5061414e89828a01613fe8565b9150509295509295509295565b60006020828403121561416d57600080fd5b813561327081613c12565b63ffffffff8116811461142f57600080fd5b8035613c3f81614178565b62ffffff8116811461142f57600080fd5b8035613c3f81614195565b61ffff8116811461142f57600080fd5b8035613c3f816141b1565b6bffffffffffffffffffffffff8116811461142f57600080fd5b8035613c3f816141cc565b801515811461142f57600080fd5b8035613c3f816141f1565b6000610220828403121561421d57600080fd5b614225613d64565b90506142308261418a565b815261423e6020830161418a565b602082015261424f6040830161418a565b6040820152614260606083016141a6565b6060820152614271608083016141c1565b608082015261428260a083016141e6565b60a082015261429360c0830161418a565b60c08201526142a460e0830161418a565b60e08201526101006142b781840161418a565b908201526101206142c983820161418a565b90820152610140828101359082015261016080830135908201526101806142f1818401613c34565b908201526101a08281013567ffffffffffffffff81111561431157600080fd5b61431d85828601613f73565b8284015250506101c0614331818401613c34565b908201526101e0614343838201613c34565b908201526102006143558382016141ff565b9082015292915050565b60008060008060008060c0878903121561437857600080fd5b863567ffffffffffffffff8082111561439057600080fd5b61439c8a838b01613f73565b975060208901359150808211156143b257600080fd5b6143be8a838b01613f73565b96506143cc60408a01613fd7565b955060608901359150808211156143e257600080fd5b61411d8a838b0161420a565b60006020828403121561440057600080fd5b5035919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6bffffffffffffffffffffffff818116838216019080821115613a8857613a88614407565b8082018082111561299f5761299f614407565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60ff818116838216019081111561299f5761299f614407565b6000602082840312156144c857600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b61ffff818116838216019080821115613a8857613a88614407565b8181038181111561299f5761299f614407565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361455d5761455d614407565b5060010190565b808202811582820484141761299f5761299f614407565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000826145b9576145b961457b565b500490565b6bffffffffffffffffffffffff851681528360208201528260408201526080606082015260006145f16080830184613b9b565b9695505050505050565b8051613c3f81614178565b8051613c3f81614195565b8051613c3f816141b1565b8051613c3f816141cc565b8051613c3f81613c12565b600082601f83011261464357600080fd5b81516020614653613e4583613e00565b82815260059290921b8401810191818101908684111561467257600080fd5b8286015b84811015613e8457805161468981613c12565b8352918301918301614676565b8051613c3f816141f1565b6000602082840312156146b357600080fd5b815167ffffffffffffffff808211156146cb57600080fd5b9083019061022082860312156146e057600080fd5b6146e8613d64565b6146f1836145fb565b81526146ff602084016145fb565b6020820152614710604084016145fb565b604082015261472160608401614606565b606082015261473260808401614611565b608082015261474360a0840161461c565b60a082015261475460c084016145fb565b60c082015261476560e084016145fb565b60e08201526101006147788185016145fb565b9082015261012061478a8482016145fb565b90820152610140838101519082015261016080840151908201526101806147b2818501614627565b908201526101a083810151838111156147ca57600080fd5b6147d688828701614632565b8284015250506101c091506147ec828401614627565b828201526101e09150614800828401614627565b828201526102009150614814828401614696565b91810191909152949350505050565b60ff8181168382160290811690818114613a8857613a88614407565b63ffffffff818116838216019080821115613a8857613a88614407565b600081518084526020808501945080840160005b838110156148a257815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101614870565b509495945050505050565b602081526148c460208201835163ffffffff169052565b600060208301516148dd604084018263ffffffff169052565b50604083015163ffffffff8116606084015250606083015162ffffff8116608084015250608083015161ffff811660a08401525060a08301516bffffffffffffffffffffffff811660c08401525060c083015163ffffffff811660e08401525060e08301516101006149568185018363ffffffff169052565b840151905061012061496f8482018363ffffffff169052565b84015190506101406149888482018363ffffffff169052565b840151610160848101919091528401516101808085019190915284015190506101a06149cb8185018373ffffffffffffffffffffffffffffffffffffffff169052565b808501519150506102206101c081818601526149eb61024086018461485c565b908601519092506101e0614a168682018373ffffffffffffffffffffffffffffffffffffffff169052565b8601519050610200614a3f8682018373ffffffffffffffffffffffffffffffffffffffff169052565b90950151151593019290925250919050565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152614a818184018a61485c565b90508281036080840152614a95818961485c565b905060ff871660a084015282810360c0840152614ab28187613b9b565b905067ffffffffffffffff851660e0840152828103610100840152614ad78185613b9b565b9c9b505050505050505050505050565b828152604060208201526000612a486040830184613b9b565b60008060408385031215614b1357600080fd5b8251614b1e816141f1565b6020939093015192949293505050565b8183823760009101908152919050565b8281526080810160608360208401379392505050565b600082601f830112614b6557600080fd5b81356020614b75613e4583613e00565b82815260059290921b84018101918181019086841115614b9457600080fd5b8286015b84811015613e8457803567ffffffffffffffff811115614bb85760008081fd5b614bc68986838b0101613fe8565b845250918301918301614b98565b600060208284031215614be657600080fd5b813567ffffffffffffffff80821115614bfe57600080fd5b9083019060c08286031215614c1257600080fd5b614c1a613d8e565b8235815260208301356020820152604083013582811115614c3a57600080fd5b614c4687828601613e24565b604083015250606083013582811115614c5e57600080fd5b614c6a87828601613e24565b606083015250608083013582811115614c8257600080fd5b614c8e87828601614b54565b60808301525060a083013582811115614ca657600080fd5b614cb287828601614b54565b60a08301525095945050505050565b6bffffffffffffffffffffffff828116828216039080821115613a8857613a88614407565b60006bffffffffffffffffffffffff80841680614d0557614d0561457b565b92169190910492915050565b6bffffffffffffffffffffffff818116838216028082169190828114614d3957614d39614407565b505092915050565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b166040850152816060850152614d888285018b61485c565b91508382036080850152614d9c828a61485c565b915060ff881660a085015283820360c0850152614db98288613b9b565b90861660e08501528381036101008501529050614ad78185613b9b565b600060408284031215614de857600080fd5b6040516040810181811067ffffffffffffffff82111715614e0b57614e0b613d35565b6040528251614e1981614178565b81526020928301519281019290925250919050565b600060a08284031215614e4057600080fd5b60405160a0810181811067ffffffffffffffff82111715614e6357614e63613d35565b806040525082518152602083015160208201526040830151614e8481614178565b60408201526060830151614e9781614178565b60608201526080928301519281019290925250919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000813000a",
>>>>>>> 9507a3be49 (adjust gas overheads)
=======
	Bin: "0x6101206040523480156200001257600080fd5b506040516200529938038062005299833981016040819052620000359162000345565b80816001600160a01b031663ca30e6036040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200009b919062000345565b826001600160a01b031663b10b673c6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000da573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000100919062000345565b836001600160a01b0316636709d0e56040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200013f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000165919062000345565b846001600160a01b0316635425d8ac6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca919062000345565b3380600081620002215760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200025457620002548162000281565b5050506001600160a01b0393841660805291831660a052821660c052811660e0521661010052506200036c565b336001600160a01b03821603620002db5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000218565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200034257600080fd5b50565b6000602082840312156200035857600080fd5b815162000365816200032c565b9392505050565b60805160a05160c05160e05161010051614eeb620003ae6000396000818160d6015261016f015260005050600050506000505060006104330152614eeb6000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c8063aed2e92911610081578063e3d0e7121161005b578063e3d0e712146102e0578063f2fde38b146102f3578063f75f6b1114610306576100d4565b8063aed2e92914610262578063afcb95d71461028c578063b1dc65a4146102cd576100d4565b806381ff7048116100b257806381ff7048146101bc5780638da5cb5b14610231578063a4c0ed361461024f576100d4565b8063181f5a771461011b578063349e8cca1461016d57806379ba5097146101b4575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e808015610114573d6000f35b3d6000fd5b005b6101576040518060400160405280601881526020017f4175746f6d6174696f6e526567697374727920322e322e30000000000000000081525081565b6040516101649190613bff565b60405180910390f35b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610164565b610119610319565b61020e60155460115463ffffffff780100000000000000000000000000000000000000000000000083048116937c01000000000000000000000000000000000000000000000000000000009093041691565b6040805163ffffffff948516815293909216602084015290820152606001610164565b60005473ffffffffffffffffffffffffffffffffffffffff1661018f565b61011961025d366004613c8d565b61041b565b610275610270366004613ce9565b610637565b604080519215158352602083019190915201610164565b601154601254604080516000815260208101939093527401000000000000000000000000000000000000000090910463ffffffff1690820152606001610164565b6101196102db366004613e8f565b6107ad565b6101196102ee36600461408e565b6113f5565b61011961030136600461415b565b61141e565b61011961031436600461435f565b611432565b60015473ffffffffffffffffffffffffffffffffffffffff16331461039f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461048a576040517fc8bad78d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b602081146104c4576040517fdfe9309000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006104d2828401846143ee565b60008181526004602052604090205490915065010000000000900463ffffffff9081161461052c576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000818152600460205260409020600101546105679085906c0100000000000000000000000090046bffffffffffffffffffffffff16614436565b600082815260046020526040902060010180546bffffffffffffffffffffffff929092166c01000000000000000000000000027fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff9092169190911790556019546105d290859061445b565b6019556040516bffffffffffffffffffffffff8516815273ffffffffffffffffffffffffffffffffffffffff86169082907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a35050505050565b60008061064261241e565b6012547e01000000000000000000000000000000000000000000000000000000000000900460ff16156106a1576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600085815260046020908152604091829020825160e081018452815460ff811615158252610100810463ffffffff908116838601819052650100000000008304821684880152690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff16606084018190526001909401546bffffffffffffffffffffffff80821660808601526c0100000000000000000000000082041660a0850152780100000000000000000000000000000000000000000000000090041660c08301528451601f890185900485028101850190955287855290936107a093899089908190840183828082843760009201919091525061245892505050565b9097909650945050505050565b60005a60408051610160810182526012546bffffffffffffffffffffffff8116825263ffffffff6c010000000000000000000000008204811660208401527001000000000000000000000000000000008204811693830193909352740100000000000000000000000000000000000000008104909216606082015262ffffff7801000000000000000000000000000000000000000000000000830416608082015261ffff7b0100000000000000000000000000000000000000000000000000000083041660a082015260ff7d0100000000000000000000000000000000000000000000000000000000008304811660c08301527e010000000000000000000000000000000000000000000000000000000000008304811615801560e08401527f01000000000000000000000000000000000000000000000000000000000000009093048116151561010080840191909152601354918216151561012084015273ffffffffffffffffffffffffffffffffffffffff910416610140820152919250610963576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600b602052604090205460ff166109ac576040517f1099ed7500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6011548935146109e8576040517fdfdcf8e700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60c08101516109f890600161449d565b60ff1685141580610a0a575083518514155b15610a41576040517f0244f71a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610a5089898989898989612681565b6000610a5c89896128ea565b9050600081604001515167ffffffffffffffff811115610a7e57610a7e613d35565b604051908082528060200260200182016040528015610b4257816020015b604080516101e0810182526000610100820181815261012083018290526101408301829052610160830182905261018083018290526101a083018290526101c0830182905282526020808301829052928201819052606082018190526080820181905260a0820181905260c0820181905260e082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181610a9c5790505b50905060008084610140015173ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa158015610b98573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bbc91906144b6565b905060005b846040015151811015611006576004600086604001518381518110610be857610be861446e565b6020908102919091018101518252818101929092526040908101600020815160e081018352815460ff811615158252610100810463ffffffff90811695830195909552650100000000008104851693820193909352690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff166060830152600101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490911660c08201528451859083908110610ccd57610ccd61446e565b602002602001015160000181905250610d0285604001518281518110610cf557610cf561446e565b60200260200101516129a5565b848281518110610d1457610d1461446e565b6020026020010151608001906001811115610d3157610d316144cf565b90816001811115610d4457610d446144cf565b81525050610db886858381518110610d5e57610d5e61446e565b60200260200101516080015187606001518481518110610d8057610d8061446e565b60200260200101518860a001518581518110610d9e57610d9e61446e565b60200260200101515189600001518a602001516001612a50565b848281518110610dca57610dca61446e565b6020026020010151604001906bffffffffffffffffffffffff1690816bffffffffffffffffffffffff1681525050610e5885604001518281518110610e1157610e1161446e565b60200260200101518387608001518481518110610e3057610e3061446e565b6020026020010151878581518110610e4a57610e4a61446e565b60200260200101518a612a9b565b858381518110610e6a57610e6a61446e565b6020026020010151602001868481518110610e8757610e8761446e565b602002602001015160e0018281525082151515158152505050838181518110610eb257610eb261446e565b60200260200101516020015115610ed557610ece6001846144fe565b9250610eda565b610ff4565b610f40848281518110610eef57610eef61446e565b6020026020010151600001516060015186606001518381518110610f1557610f1561446e565b60200260200101518760a001518481518110610f3357610f3361446e565b6020026020010151612458565b858381518110610f5257610f5261446e565b6020026020010151606001868481518110610f6f57610f6f61446e565b602002602001015160a0018281525082151515158152505050838181518110610f9a57610f9a61446e565b602002602001015160a0015187610fb19190614519565b9650610ff485604001518281518110610fcc57610fcc61446e565b602002602001015183868481518110610fe757610fe761446e565b6020026020010151612c1a565b80610ffe8161452c565b915050610bc1565b508161ffff1660000361101e575050505050506113ec565b60c085015161102e90600161449d565b61103d9060ff1661044c614564565b616dc461104b8d6010614564565b5a611056908a614519565b611060919061445b565b61106a919061445b565b611074919061445b565b9550611c2061108761ffff8416886145aa565b611091919061445b565b955060008060008060005b886040015151811015611292578781815181106110bb576110bb61446e565b60200260200101516020015115611280576111178b8983815181106110e2576110e261446e565b6020026020010151608001518b60a0015184815181106111045761110461446e565b6020026020010151518d60c00151612d1f565b8882815181106111295761112961446e565b602002602001015160c00181815250506111858a8a6040015183815181106111535761115361446e565b60200260200101518a848151811061116d5761116d61446e565b60200260200101518c600001518d602001518c612d3f565b90935091506111948285614436565b93506111a08386614436565b94508781815181106111b4576111b461446e565b6020026020010151606001511515896040015182815181106111d8576111d861446e565b60200260200101517fad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b848661120d9190614436565b8b858151811061121f5761121f61446e565b602002602001015160a001518c868151811061123d5761123d61446e565b602002602001015160c001518e60800151878151811061125f5761125f61446e565b602002602001015160405161127794939291906145be565b60405180910390a35b8061128a8161452c565b91505061109c565b5050336000908152600b6020526040902080548492506002906112ca9084906201000090046bffffffffffffffffffffffff16614436565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555080601260000160008282829054906101000a90046bffffffffffffffffffffffff166113249190614436565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060008f6001600381106113675761136761446e565b602002013560001c9050600060088264ffffffffff16901c9050886060015163ffffffff168163ffffffff1611156113e157601280547fffffffffffffffff00000000ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000063ffffffff8416021790555b505050505050505050505b50505050505050565b6114168686868680602001905181019061140f91906146a1565b8686611432565b505050505050565b611426612e32565b61142f81612eb3565b50565b61143a612e32565b601f86511115611476576040517f25d0209c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8360ff166000036114b3576040517fe77dba5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845186511415806114d257506114ca846003614823565b60ff16865111155b15611509576040517f1d2d1c5800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601254600e546bffffffffffffffffffffffff9091169060005b816bffffffffffffffffffffffff1681101561158b57611578600e828154811061154f5761154f61446e565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff168484612fa8565b50806115838161452c565b915050611523565b5060008060005b836bffffffffffffffffffffffff1681101561169457600d81815481106115bb576115bb61446e565b600091825260209091200154600e805473ffffffffffffffffffffffffffffffffffffffff909216945090829081106115f6576115f661446e565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff8681168452600c8352604080852080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001690559116808452600b90925290912080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905591508061168c8161452c565b915050611592565b506116a1600d6000613ade565b6116ad600e6000613ade565b604080516080810182526000808252602082018190529181018290526060810182905290805b8c51811015611b1657600c60008e83815181106116f2576116f261446e565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528101919091526040016000205460ff161561175d576040517f77cea0fa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168d82815181106117875761178761446e565b602002602001015173ffffffffffffffffffffffffffffffffffffffff16036117dc576040517f815e1d6400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405180604001604052806001151581526020018260ff16815250600c60008f848151811061180d5761180d61446e565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528181019290925260400160002082518154939092015160ff16610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909316929092171790558b518c90829081106118b5576118b561446e565b60200260200101519150600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603611925576040517f58a70a0a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff82166000908152600b60209081526040918290208251608081018452905460ff80821615801584526101008304909116938301939093526bffffffffffffffffffffffff6201000082048116948301949094526e010000000000000000000000000000900490921660608301529093506119e0576040517f6a7281ad00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001835260ff80821660208086019182526bffffffffffffffffffffffff808b166060880190815273ffffffffffffffffffffffffffffffffffffffff87166000908152600b909352604092839020885181549551948a0151925184166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff939094166201000002929092167fffffffffffff000000000000000000000000000000000000000000000000ffff94909616610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009095169490941717919091169290921791909117905580611b0e8161452c565b9150506116d3565b50508a51611b2c9150600d9060208d0190613afc565b508851611b4090600e9060208c0190613afc565b50604051806101600160405280856bffffffffffffffffffffffff168152602001886000015163ffffffff168152602001886020015163ffffffff168152602001600063ffffffff168152602001886060015162ffffff168152602001886080015161ffff1681526020018960ff1681526020016012600001601e9054906101000a900460ff16151581526020016012600001601f9054906101000a900460ff161515815260200188610200015115158152602001886101e0015173ffffffffffffffffffffffffffffffffffffffff16815250601260008201518160000160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550602082015181600001600c6101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160106101000a81548163ffffffff021916908363ffffffff16021790555060608201518160000160146101000a81548163ffffffff021916908363ffffffff16021790555060808201518160000160186101000a81548162ffffff021916908362ffffff16021790555060a082015181600001601b6101000a81548161ffff021916908361ffff16021790555060c082015181600001601d6101000a81548160ff021916908360ff16021790555060e082015181600001601e6101000a81548160ff02191690831515021790555061010082015181600001601f6101000a81548160ff0219169083151502179055506101208201518160010160006101000a81548160ff0219169083151502179055506101408201518160010160016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055509050506040518061018001604052808860a001516bffffffffffffffffffffffff16815260200188610180015173ffffffffffffffffffffffffffffffffffffffff168152602001601460010160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff168152602001886040015163ffffffff1681526020018860c0015163ffffffff168152602001601460010160149054906101000a900463ffffffff1663ffffffff168152602001601460010160189054906101000a900463ffffffff1663ffffffff1681526020016014600101601c9054906101000a900463ffffffff1663ffffffff1681526020018860e0015163ffffffff16815260200188610100015163ffffffff16815260200188610120015163ffffffff168152602001886101c0015173ffffffffffffffffffffffffffffffffffffffff16815250601460008201518160000160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550602082015181600001600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060408201518160010160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550606082015181600101600c6101000a81548163ffffffff021916908363ffffffff16021790555060808201518160010160106101000a81548163ffffffff021916908363ffffffff16021790555060a08201518160010160146101000a81548163ffffffff021916908363ffffffff16021790555060c08201518160010160186101000a81548163ffffffff021916908363ffffffff16021790555060e082015181600101601c6101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160020160006101000a81548163ffffffff021916908363ffffffff1602179055506101208201518160020160046101000a81548163ffffffff021916908363ffffffff1602179055506101408201518160020160086101000a81548163ffffffff021916908363ffffffff16021790555061016082015181600201600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505086610140015160178190555086610160015160188190555060006014600101601c9054906101000a900463ffffffff169050876101e0015173ffffffffffffffffffffffffffffffffffffffff166357e871e76040518163ffffffff1660e01b8152600401602060405180830381865afa1580156121e1573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061220591906144b6565b601580547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167c010000000000000000000000000000000000000000000000000000000063ffffffff9384160217808255600192601891612280918591780100000000000000000000000000000000000000000000000090041661483f565b92506101000a81548163ffffffff021916908363ffffffff1602179055506000886040516020016122b191906148ad565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291905260155490915061231a90469030907801000000000000000000000000000000000000000000000000900463ffffffff168f8f8f878f8f6131b0565b60115560005b61232a600961325a565b81101561235a5761234761233f600983613264565b600990613277565b50806123528161452c565b915050612320565b5060005b896101a00151518110156123b15761239e8a6101a0015182815181106123865761238661446e565b6020026020010151600961329990919063ffffffff16565b50806123a98161452c565b91505061235e565b507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0582601154601460010160189054906101000a900463ffffffff168f8f8f878f8f60405161240899989796959493929190614a51565b60405180910390a1505050505050505050505050565b3215612456576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60125460009081907f0100000000000000000000000000000000000000000000000000000000000000900460ff16156124bd576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601280547effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f01000000000000000000000000000000000000000000000000000000000000001790556040517f4585e33b0000000000000000000000000000000000000000000000000000000090612539908590602401613bff565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009094169390931790925290517f79188d1600000000000000000000000000000000000000000000000000000000815290935073ffffffffffffffffffffffffffffffffffffffff8616906379188d169061260c9087908790600401614ae7565b60408051808303816000875af115801561262a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061264e9190614b00565b601280547effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff16905590969095509350505050565b60008686604051612693929190614b2e565b6040519081900381206126aa918a90602001614b3e565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201208383019092526000808452908301819052909250906000805b87811015612882576001858783602081106127165761271661446e565b61272391901a601b61449d565b8b8b858181106127355761273561446e565b905060200201358a858151811061274e5761274e61446e565b60200260200101516040516000815260200160405260405161278c949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa1580156127ae573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff81166000908152600c602090815290849020838501909452925460ff808216151580855261010090920416938301939093529095509350905061285c576040517f0f4c073700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826020015160080260ff166001901b84019350808061287a9061452c565b9150506126f9565b50827e010101010101010101010101010101010101010101010101010101010101018416146128dd576040517fc103be2e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050505050505050505050565b6129236040518060c001604052806000815260200160008152602001606081526020016060815260200160608152602001606081525090565b600061293183850185614bd4565b604081015151606082015151919250908114158061295457508082608001515114155b806129645750808260a001515114155b1561299b576040517fb55ac75400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5090505b92915050565b6000818160045b600f811015612a32577fff0000000000000000000000000000000000000000000000000000000000000082168382602081106129ea576129ea61446e565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614612a2057506000949350505050565b80612a2a8161452c565b9150506129ac565b5081600f1a6001811115612a4857612a486144cf565b949350505050565b600080612a6288878b60c001516132bb565b9050600080612a7d8b8a63ffffffff16858a8a60018b613347565b9092509050612a8c8183614436565b9b9a5050505050505050505050565b600080808085608001516001811115612ab657612ab66144cf565b03612adc57612ac888888888886135d4565b612ad757600092509050612c10565b612b54565b600185608001516001811115612af457612af46144cf565b03612b22576000612b078989898861375d565b9250905080612b1c5750600092509050612c10565b50612b54565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b84516040015163ffffffff168710612ba957877fc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd563687604051612b969190613bff565b60405180910390a2600092509050612c10565b84604001516bffffffffffffffffffffffff16856000015160a001516bffffffffffffffffffffffff161015612c0957877f377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e0287604051612b969190613bff565b6001925090505b9550959350505050565b600081608001516001811115612c3257612c326144cf565b03612c9657600083815260046020526040902060010180547fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff16780100000000000000000000000000000000000000000000000063ffffffff851602179055505050565b600181608001516001811115612cae57612cae6144cf565b03612d1a5760e08101805160009081526008602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055915191517fa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f29190a25b505050565b6000612d2c8484846132bb565b905080851015612a485750929392505050565b600080612d5a888760a001518860c001518888886001613347565b90925090506000612d6b8284614436565b600089815260046020526040902060010180549192508291600c90612daf9084906c0100000000000000000000000090046bffffffffffffffffffffffff16614cc1565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915560008a815260046020526040812060010180548594509092612df891859116614436565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555050965096945050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314612456576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610396565b3373ffffffffffffffffffffffffffffffffffffffff821603612f32576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610396565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e01000000000000000000000000000090049091166060820152906131a45760008160600151856130409190614cc1565b9050600061304e8583614ce6565b905080836040018181516130629190614436565b6bffffffffffffffffffffffff1690525061307d8582614d11565b8360600181815161308e9190614436565b6bffffffffffffffffffffffff90811690915273ffffffffffffffffffffffffffffffffffffffff89166000908152600b602090815260409182902087518154928901519389015160608a015186166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff919096166201000002167fffffffffffff000000000000000000000000000000000000000000000000ffff60ff95909516610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909416939093171792909216179190911790555050505b60400151949350505050565b6000808a8a8a8a8a8a8a8a8a6040516020016131d499989796959493929190614d41565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179b9a5050505050505050505050565b600061299f825490565b6000613270838361396b565b9392505050565b60006132708373ffffffffffffffffffffffffffffffffffffffff8416613995565b60006132708373ffffffffffffffffffffffffffffffffffffffff8416613a8f565b600080808560018111156132d1576132d16144cf565b036132e0575062015f906132ff565b60018560018111156132f4576132f46144cf565b03612b2257506201af405b61331063ffffffff85166014614564565b61331b84600161449d565b61332a9060ff16611d4c614564565b613334908361445b565b61333e919061445b565b95945050505050565b60008060008960a0015161ffff16876133609190614564565b905083801561336e5750803a105b1561337657503a5b600084156133f9578a610140015173ffffffffffffffffffffffffffffffffffffffff166318b8f6136040518163ffffffff1660e01b8152600401602060405180830381865afa1580156133ce573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906133f291906144b6565b90506134b5565b6101408b01516016546040517f1254414000000000000000000000000000000000000000000000000000000000815264010000000090910463ffffffff16600482015273ffffffffffffffffffffffffffffffffffffffff90911690631254414090602401602060405180830381865afa15801561347b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061349f91906144b6565b8b60a0015161ffff166134b29190614564565b90505b6134c361ffff8716826145aa565b9050600087826134d38c8e61445b565b6134dd9086614564565b6134e7919061445b565b6134f990670de0b6b3a7640000614564565b61350391906145aa565b905060008c6040015163ffffffff1664e8d4a510006135229190614564565b898e6020015163ffffffff16858f8861353b9190614564565b613545919061445b565b61355390633b9aca00614564565b61355d9190614564565b61356791906145aa565b613571919061445b565b90506b033b2e3c9fd0803ce800000061358a828461445b565b11156135c2576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b600080848060200190518101906135eb9190614dd6565b845160c00151815191925063ffffffff9081169116101561364857867f405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8866040516136369190613bff565b60405180910390a2600091505061333e565b826101200151801561370957506020810151158015906137095750602081015161014084015182516040517f85df51fd00000000000000000000000000000000000000000000000000000000815263ffffffff909116600482015273ffffffffffffffffffffffffffffffffffffffff909116906385df51fd90602401602060405180830381865afa1580156136e2573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061370691906144b6565b14155b8061371b5750805163ffffffff168611155b1561375057867f6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301866040516136369190613bff565b5060019695505050505050565b6000806000848060200190518101906137769190614e2e565b90506000878260000151836020015184604001516040516020016137d894939291909384526020840192909252604083015260e01b7fffffffff0000000000000000000000000000000000000000000000000000000016606082015260640190565b60405160208183030381529060405280519060200120905084610120015180156138b457506080820151158015906138b45750608082015161014086015160608401516040517f85df51fd00000000000000000000000000000000000000000000000000000000815263ffffffff909116600482015273ffffffffffffffffffffffffffffffffffffffff909116906385df51fd90602401602060405180830381865afa15801561388d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906138b191906144b6565b14155b806138c9575086826060015163ffffffff1610155b1561391357877f6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301876040516138fe9190613bff565b60405180910390a26000935091506139629050565b60008181526008602052604090205460ff161561395a57877f405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8876040516138fe9190613bff565b600193509150505b94509492505050565b60008260000182815481106139825761398261446e565b9060005260206000200154905092915050565b60008181526001830160205260408120548015613a7e5760006139b9600183614519565b85549091506000906139cd90600190614519565b9050818114613a325760008660000182815481106139ed576139ed61446e565b9060005260206000200154905080876000018481548110613a1057613a1061446e565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080613a4357613a43614eaf565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505061299f565b600091505061299f565b5092915050565b6000818152600183016020526040812054613ad65750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915561299f565b50600061299f565b508054600082559060005260206000209081019061142f9190613b86565b828054828255906000526020600020908101928215613b76579160200282015b82811115613b7657825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190613b1c565b50613b82929150613b86565b5090565b5b80821115613b825760008155600101613b87565b6000815180845260005b81811015613bc157602081850181015186830182015201613ba5565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815260006132706020830184613b9b565b73ffffffffffffffffffffffffffffffffffffffff8116811461142f57600080fd5b8035613c3f81613c12565b919050565b60008083601f840112613c5657600080fd5b50813567ffffffffffffffff811115613c6e57600080fd5b602083019150836020828501011115613c8657600080fd5b9250929050565b60008060008060608587031215613ca357600080fd5b8435613cae81613c12565b935060208501359250604085013567ffffffffffffffff811115613cd157600080fd5b613cdd87828801613c44565b95989497509550505050565b600080600060408486031215613cfe57600080fd5b83359250602084013567ffffffffffffffff811115613d1c57600080fd5b613d2886828701613c44565b9497909650939450505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610220810167ffffffffffffffff81118282101715613d8857613d88613d35565b60405290565b60405160c0810167ffffffffffffffff81118282101715613d8857613d88613d35565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613df857613df8613d35565b604052919050565b600067ffffffffffffffff821115613e1a57613e1a613d35565b5060051b60200190565b600082601f830112613e3557600080fd5b81356020613e4a613e4583613e00565b613db1565b82815260059290921b84018101918181019086841115613e6957600080fd5b8286015b84811015613e845780358352918301918301613e6d565b509695505050505050565b600080600080600080600060e0888a031215613eaa57600080fd5b6060880189811115613ebb57600080fd5b8897503567ffffffffffffffff80821115613ed557600080fd5b613ee18b838c01613c44565b909850965060808a0135915080821115613efa57600080fd5b818a0191508a601f830112613f0e57600080fd5b813581811115613f1d57600080fd5b8b60208260051b8501011115613f3257600080fd5b6020830196508095505060a08a0135915080821115613f5057600080fd5b50613f5d8a828b01613e24565b92505060c0880135905092959891949750929550565b600082601f830112613f8457600080fd5b81356020613f94613e4583613e00565b82815260059290921b84018101918181019086841115613fb357600080fd5b8286015b84811015613e84578035613fca81613c12565b8352918301918301613fb7565b803560ff81168114613c3f57600080fd5b600082601f830112613ff957600080fd5b813567ffffffffffffffff81111561401357614013613d35565b61404460207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601613db1565b81815284602083860101111561405957600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff81168114613c3f57600080fd5b60008060008060008060c087890312156140a757600080fd5b863567ffffffffffffffff808211156140bf57600080fd5b6140cb8a838b01613f73565b975060208901359150808211156140e157600080fd5b6140ed8a838b01613f73565b96506140fb60408a01613fd7565b9550606089013591508082111561411157600080fd5b61411d8a838b01613fe8565b945061412b60808a01614076565b935060a089013591508082111561414157600080fd5b5061414e89828a01613fe8565b9150509295509295509295565b60006020828403121561416d57600080fd5b813561327081613c12565b63ffffffff8116811461142f57600080fd5b8035613c3f81614178565b62ffffff8116811461142f57600080fd5b8035613c3f81614195565b61ffff8116811461142f57600080fd5b8035613c3f816141b1565b6bffffffffffffffffffffffff8116811461142f57600080fd5b8035613c3f816141cc565b801515811461142f57600080fd5b8035613c3f816141f1565b6000610220828403121561421d57600080fd5b614225613d64565b90506142308261418a565b815261423e6020830161418a565b602082015261424f6040830161418a565b6040820152614260606083016141a6565b6060820152614271608083016141c1565b608082015261428260a083016141e6565b60a082015261429360c0830161418a565b60c08201526142a460e0830161418a565b60e08201526101006142b781840161418a565b908201526101206142c983820161418a565b90820152610140828101359082015261016080830135908201526101806142f1818401613c34565b908201526101a08281013567ffffffffffffffff81111561431157600080fd5b61431d85828601613f73565b8284015250506101c0614331818401613c34565b908201526101e0614343838201613c34565b908201526102006143558382016141ff565b9082015292915050565b60008060008060008060c0878903121561437857600080fd5b863567ffffffffffffffff8082111561439057600080fd5b61439c8a838b01613f73565b975060208901359150808211156143b257600080fd5b6143be8a838b01613f73565b96506143cc60408a01613fd7565b955060608901359150808211156143e257600080fd5b61411d8a838b0161420a565b60006020828403121561440057600080fd5b5035919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6bffffffffffffffffffffffff818116838216019080821115613a8857613a88614407565b8082018082111561299f5761299f614407565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60ff818116838216019081111561299f5761299f614407565b6000602082840312156144c857600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b61ffff818116838216019080821115613a8857613a88614407565b8181038181111561299f5761299f614407565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361455d5761455d614407565b5060010190565b808202811582820484141761299f5761299f614407565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000826145b9576145b961457b565b500490565b6bffffffffffffffffffffffff851681528360208201528260408201526080606082015260006145f16080830184613b9b565b9695505050505050565b8051613c3f81614178565b8051613c3f81614195565b8051613c3f816141b1565b8051613c3f816141cc565b8051613c3f81613c12565b600082601f83011261464357600080fd5b81516020614653613e4583613e00565b82815260059290921b8401810191818101908684111561467257600080fd5b8286015b84811015613e8457805161468981613c12565b8352918301918301614676565b8051613c3f816141f1565b6000602082840312156146b357600080fd5b815167ffffffffffffffff808211156146cb57600080fd5b9083019061022082860312156146e057600080fd5b6146e8613d64565b6146f1836145fb565b81526146ff602084016145fb565b6020820152614710604084016145fb565b604082015261472160608401614606565b606082015261473260808401614611565b608082015261474360a0840161461c565b60a082015261475460c084016145fb565b60c082015261476560e084016145fb565b60e08201526101006147788185016145fb565b9082015261012061478a8482016145fb565b90820152610140838101519082015261016080840151908201526101806147b2818501614627565b908201526101a083810151838111156147ca57600080fd5b6147d688828701614632565b8284015250506101c091506147ec828401614627565b828201526101e09150614800828401614627565b828201526102009150614814828401614696565b91810191909152949350505050565b60ff8181168382160290811690818114613a8857613a88614407565b63ffffffff818116838216019080821115613a8857613a88614407565b600081518084526020808501945080840160005b838110156148a257815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101614870565b509495945050505050565b602081526148c460208201835163ffffffff169052565b600060208301516148dd604084018263ffffffff169052565b50604083015163ffffffff8116606084015250606083015162ffffff8116608084015250608083015161ffff811660a08401525060a08301516bffffffffffffffffffffffff811660c08401525060c083015163ffffffff811660e08401525060e08301516101006149568185018363ffffffff169052565b840151905061012061496f8482018363ffffffff169052565b84015190506101406149888482018363ffffffff169052565b840151610160848101919091528401516101808085019190915284015190506101a06149cb8185018373ffffffffffffffffffffffffffffffffffffffff169052565b808501519150506102206101c081818601526149eb61024086018461485c565b908601519092506101e0614a168682018373ffffffffffffffffffffffffffffffffffffffff169052565b8601519050610200614a3f8682018373ffffffffffffffffffffffffffffffffffffffff169052565b90950151151593019290925250919050565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152614a818184018a61485c565b90508281036080840152614a95818961485c565b905060ff871660a084015282810360c0840152614ab28187613b9b565b905067ffffffffffffffff851660e0840152828103610100840152614ad78185613b9b565b9c9b505050505050505050505050565b828152604060208201526000612a486040830184613b9b565b60008060408385031215614b1357600080fd5b8251614b1e816141f1565b6020939093015192949293505050565b8183823760009101908152919050565b8281526080810160608360208401379392505050565b600082601f830112614b6557600080fd5b81356020614b75613e4583613e00565b82815260059290921b84018101918181019086841115614b9457600080fd5b8286015b84811015613e8457803567ffffffffffffffff811115614bb85760008081fd5b614bc68986838b0101613fe8565b845250918301918301614b98565b600060208284031215614be657600080fd5b813567ffffffffffffffff80821115614bfe57600080fd5b9083019060c08286031215614c1257600080fd5b614c1a613d8e565b8235815260208301356020820152604083013582811115614c3a57600080fd5b614c4687828601613e24565b604083015250606083013582811115614c5e57600080fd5b614c6a87828601613e24565b606083015250608083013582811115614c8257600080fd5b614c8e87828601614b54565b60808301525060a083013582811115614ca657600080fd5b614cb287828601614b54565b60a08301525095945050505050565b6bffffffffffffffffffffffff828116828216039080821115613a8857613a88614407565b60006bffffffffffffffffffffffff80841680614d0557614d0561457b565b92169190910492915050565b6bffffffffffffffffffffffff818116838216028082169190828114614d3957614d39614407565b505092915050565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b166040850152816060850152614d888285018b61485c565b91508382036080850152614d9c828a61485c565b915060ff881660a085015283820360c0850152614db98288613b9b565b90861660e08501528381036101008501529050614ad78185613b9b565b600060408284031215614de857600080fd5b6040516040810181811067ffffffffffffffff82111715614e0b57614e0b613d35565b6040528251614e1981614178565b81526020928301519281019290925250919050565b600060a08284031215614e4057600080fd5b60405160a0810181811067ffffffffffffffff82111715614e6357614e63613d35565b806040525082518152602083015160208201526040830151614e8481614178565b60408201526060830151614e9781614178565b60608201526080928301519281019290925250919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000813000a",
>>>>>>> 911f05e7a3 (adjust gas overhead again)
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

type AutomationRegistryChainSpecificModuleUpdatedIterator struct {
	Event *AutomationRegistryChainSpecificModuleUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AutomationRegistryChainSpecificModuleUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AutomationRegistryChainSpecificModuleUpdated)
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
		it.Event = new(AutomationRegistryChainSpecificModuleUpdated)
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

func (it *AutomationRegistryChainSpecificModuleUpdatedIterator) Error() error {
	return it.fail
}

func (it *AutomationRegistryChainSpecificModuleUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AutomationRegistryChainSpecificModuleUpdated struct {
	NewModule common.Address
	Raw       types.Log
}

func (_AutomationRegistry *AutomationRegistryFilterer) FilterChainSpecificModuleUpdated(opts *bind.FilterOpts) (*AutomationRegistryChainSpecificModuleUpdatedIterator, error) {

	logs, sub, err := _AutomationRegistry.contract.FilterLogs(opts, "ChainSpecificModuleUpdated")
	if err != nil {
		return nil, err
	}
	return &AutomationRegistryChainSpecificModuleUpdatedIterator{contract: _AutomationRegistry.contract, event: "ChainSpecificModuleUpdated", logs: logs, sub: sub}, nil
}

func (_AutomationRegistry *AutomationRegistryFilterer) WatchChainSpecificModuleUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryChainSpecificModuleUpdated) (event.Subscription, error) {

	logs, sub, err := _AutomationRegistry.contract.WatchLogs(opts, "ChainSpecificModuleUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AutomationRegistryChainSpecificModuleUpdated)
				if err := _AutomationRegistry.contract.UnpackLog(event, "ChainSpecificModuleUpdated", log); err != nil {
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

func (_AutomationRegistry *AutomationRegistryFilterer) ParseChainSpecificModuleUpdated(log types.Log) (*AutomationRegistryChainSpecificModuleUpdated, error) {
	event := new(AutomationRegistryChainSpecificModuleUpdated)
	if err := _AutomationRegistry.contract.UnpackLog(event, "ChainSpecificModuleUpdated", log); err != nil {
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
	case _AutomationRegistry.abi.Events["ChainSpecificModuleUpdated"].ID:
		return _AutomationRegistry.ParseChainSpecificModuleUpdated(log)
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

func (AutomationRegistryChainSpecificModuleUpdated) Topic() common.Hash {
	return common.HexToHash("0xdefc28b11a7980dbe0c49dbbd7055a1584bc8075097d1e8b3b57fb7283df2ad7")
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

	FilterChainSpecificModuleUpdated(opts *bind.FilterOpts) (*AutomationRegistryChainSpecificModuleUpdatedIterator, error)

	WatchChainSpecificModuleUpdated(opts *bind.WatchOpts, sink chan<- *AutomationRegistryChainSpecificModuleUpdated) (event.Subscription, error)

	ParseChainSpecificModuleUpdated(log types.Log) (*AutomationRegistryChainSpecificModuleUpdated, error)

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
