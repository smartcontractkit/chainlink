// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package keeper_registry_wrapper_2_1

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

type KeeperRegistryBase21OnchainConfig struct {
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
}

var KeeperRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractKeeperRegistryLogicB2_1\",\"name\":\"logicA\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTrigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTriggerType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"AdminPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"dedupKey\",\"type\":\"bytes32\"}],\"name\":\"DedupKeyAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepTriggerConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fallbackTo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfigBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCheckDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxRevertDataSize\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"registrars\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"upkeepPrivilegeManager\",\"type\":\"address\"}],\"internalType\":\"structKeeperRegistryBase2_1.OnchainConfig\",\"name\":\"onchainConfig\",\"type\":\"tuple\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfigTypeSafe\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"simulatePerformUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"rawReport\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6101406040523480156200001257600080fd5b50604051620054d0380380620054d08339810160408190526200003591620003df565b80816001600160a01b0316634b4fd03b6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200009b919062000406565b826001600160a01b031663ca30e6036040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000da573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001009190620003df565b836001600160a01b031663b10b673c6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200013f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001659190620003df565b846001600160a01b0316636709d0e56040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca9190620003df565b856001600160a01b0316635425d8ac6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000209573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200022f9190620003df565b3380600081620002865760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620002b957620002b9816200031b565b505050846002811115620002d157620002d162000429565b60e0816002811115620002e857620002e862000429565b9052506001600160a01b0393841660805291831660a052821660c0528116610100529190911661012052506200043f9050565b336001600160a01b03821603620003755760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200027d565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b0381168114620003dc57600080fd5b50565b600060208284031215620003f257600080fd5b8151620003ff81620003c6565b9392505050565b6000602082840312156200041957600080fd5b815160038110620003ff57600080fd5b634e487b7160e01b600052602160045260246000fd5b60805160a05160c05160e051610100516101205161502f620004a16000396000818160d6015261016f01526000505060008181612eb701528181613220015281816133b30152613a4901526000505060005050600061043b015261502f6000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c8063aed2e92911610081578063e29b753c1161005b578063e29b753c146102e8578063e3d0e712146102fb578063f2fde38b1461030e576100d4565b8063aed2e92914610262578063afcb95d71461028c578063b1dc65a4146102d5576100d4565b806381ff7048116100b257806381ff7048146101bc5780638da5cb5b14610231578063a4c0ed361461024f576100d4565b8063181f5a771461011b578063349e8cca1461016d57806379ba5097146101b4575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e808015610114573d6000f35b3d6000fd5b005b6101576040518060400160405280601481526020017f4b6565706572526567697374727920322e312e3000000000000000000000000081525081565b6040516101649190613cc8565b60405180910390f35b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610164565b610119610321565b61020e60145460115463ffffffff780100000000000000000000000000000000000000000000000083048116937c01000000000000000000000000000000000000000000000000000000009093041691565b6040805163ffffffff948516815293909216602084015290820152606001610164565b60005473ffffffffffffffffffffffffffffffffffffffff1661018f565b61011961025d366004613d51565b610423565b610275610270366004613dad565b61063f565b604080519215158352602083019190915201610164565b601154601254604080516000815260208101939093527c010000000000000000000000000000000000000000000000000000000090910463ffffffff1690820152606001610164565b6101196102e3366004613e3e565b6107a7565b6101196102f63660046142b9565b6112ea565b610119610309366004614386565b6121e6565b61011961031c366004614415565b61220f565b60015473ffffffffffffffffffffffffffffffffffffffff1633146103a7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610492576040517fc8bad78d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b602081146104cc576040517fdfe9309000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006104da82840184614432565b60008181526004602052604090205490915065010000000000900463ffffffff90811614610534576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526004602052604090206001015461056f9085906c0100000000000000000000000090046bffffffffffffffffffffffff1661447a565b600082815260046020526040902060010180546bffffffffffffffffffffffff929092166c01000000000000000000000000027fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff9092169190911790556018546105da90859061449f565b6018556040516bffffffffffffffffffffffff8516815273ffffffffffffffffffffffffffffffffffffffff86169082907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a35050505050565b60008061064a612223565b6012546e010000000000000000000000000000900460ff1615610699576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600085815260046020908152604091829020825160e081018452815460ff811615158252610100810463ffffffff908116838601819052650100000000008304821684880152690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff16606084018190526001909401546bffffffffffffffffffffffff80821660808601526c0100000000000000000000000082041660a0850152780100000000000000000000000000000000000000000000000090041660c08301528451601f8901859004850281018501909552878552909361079893899089908190840183828082843760009201919091525061225d92505050565b9093509150505b935093915050565b60005a604080516101208101825260125460ff808216835261010080830463ffffffff90811660208601526501000000000084048116958501959095526901000000000000000000830462ffffff1660608501526c01000000000000000000000000830461ffff1660808501526e0100000000000000000000000000008304821615801560a08601526f010000000000000000000000000000008404909216151560c085015270010000000000000000000000000000000083046bffffffffffffffffffffffff1660e08501527c0100000000000000000000000000000000000000000000000000000000909204909316908201529192506108d5576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600b602052604090205460ff1661091e576040517f1099ed7500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6011548a351461095a576040517fdfdcf8e700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80516109679060016144e1565b60ff16861415806109785750858414155b156109af576040517f0244f71a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6109bf8a8a8a8a8a8a8a8a612468565b60006109cb8a8a6126d1565b9050600081604001515167ffffffffffffffff8111156109ed576109ed613ef5565b604051908082528060200260200182016040528015610ab157816020015b604080516101e0810182526000610100820181815261012083018290526101408301829052610160830182905261018083018290526101a083018290526101c0830182905282526020808301829052928201819052606082018190526080820181905260a0820181905260c0820181905260e082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181610a0b5790505b5090506000805b836040015151811015610efa576004600085604001518381518110610adf57610adf6144b2565b6020908102919091018101518252818101929092526040908101600020815160e081018352815460ff811615158252610100810463ffffffff90811695830195909552650100000000008104851693820193909352690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff166060830152600101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490911660c08201528351849083908110610bc457610bc46144b2565b602002602001015160000181905250610bf984604001518281518110610bec57610bec6144b2565b602002602001015161278c565b838281518110610c0b57610c0b6144b2565b6020026020010151608001906001811115610c2857610c286144fa565b90816001811115610c3b57610c3b6144fa565b81525050610caf85848381518110610c5557610c556144b2565b60200260200101516080015186606001518481518110610c7757610c776144b2565b60200260200101518760a001518581518110610c9557610c956144b2565b602002602001015151886000015189602001516001612837565b838281518110610cc157610cc16144b2565b6020026020010151604001906bffffffffffffffffffffffff1690816bffffffffffffffffffffffff1681525050610d4d84604001518281518110610d0857610d086144b2565b602002602001015185608001518381518110610d2657610d266144b2565b6020026020010151858481518110610d4057610d406144b2565b6020026020010151612882565b848381518110610d5f57610d5f6144b2565b6020026020010151602001858481518110610d7c57610d7c6144b2565b602002602001015160e0018281525082151515158152505050828181518110610da757610da76144b2565b60200260200101516020015115610dca57610dc3600183614529565b9150610dcf565b610ee8565b610e35838281518110610de457610de46144b2565b6020026020010151600001516060015185606001518381518110610e0a57610e0a6144b2565b60200260200101518660a001518481518110610e2857610e286144b2565b602002602001015161225d565b848381518110610e4757610e476144b2565b6020026020010151606001858481518110610e6457610e646144b2565b602002602001015160a0018281525082151515158152505050828181518110610e8f57610e8f6144b2565b602002602001015160a0015186610ea69190614544565b9550610ee884604001518281518110610ec157610ec16144b2565b6020026020010151848381518110610edb57610edb6144b2565b6020026020010151612a01565b80610ef281614557565b915050610ab8565b508061ffff16600003610f115750505050506112e0565b8351610f1e9060016144e1565b610f2d9060ff1661044c61458f565b616b6c610f3b8d601061458f565b5a610f469089614544565b610f50919061449f565b610f5a919061449f565b610f64919061449f565b9450611b58610f7761ffff8316876145fb565b610f81919061449f565b945060008060008060005b87604001515181101561118257868181518110610fab57610fab6144b2565b60200260200101516020015115611170576110078a888381518110610fd257610fd26144b2565b6020026020010151608001518a60a001518481518110610ff457610ff46144b2565b6020026020010151518c60000151612b13565b878281518110611019576110196144b2565b602002602001015160c00181815250506110758989604001518381518110611043576110436144b2565b602002602001015189848151811061105d5761105d6144b2565b60200260200101518b600001518c602001518b612b33565b9093509150611084828561447a565b9350611090838661447a565b94508681815181106110a4576110a46144b2565b6020026020010151606001511515886040015182815181106110c8576110c86144b2565b60200260200101517fad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b84866110fd919061447a565b8a858151811061110f5761110f6144b2565b602002602001015160a001518b868151811061112d5761112d6144b2565b602002602001015160c001518d60800151878151811061114f5761114f6144b2565b6020026020010151604051611167949392919061460f565b60405180910390a35b8061117a81614557565b915050610f8c565b5050336000908152600b6020526040902080548492506002906111ba9084906201000090046bffffffffffffffffffffffff1661447a565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555080601260000160108282829054906101000a90046bffffffffffffffffffffffff16611214919061447a565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060008f600160038110611257576112576144b2565b602002013560001c9050600060088264ffffffffff16901c905087610100015163ffffffff168163ffffffff1611156112d657601280547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167c010000000000000000000000000000000000000000000000000000000063ffffffff8416021790555b5050505050505050505b5050505050505050565b6112f2612c26565b601f8651111561132e576040517f25d0209c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8360ff1660000361136b576040517fe77dba5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8451865114158061138a575061138284600361464c565b60ff16865111155b156113c1576040517f1d2d1c5800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601254600e547001000000000000000000000000000000009091046bffffffffffffffffffffffff169060005b816bffffffffffffffffffffffff1681101561145657611443600e828154811061141a5761141a6144b2565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff168484612ca7565b508061144e81614557565b9150506113ee565b5060008060005b836bffffffffffffffffffffffff1681101561155f57600d8181548110611486576114866144b2565b600091825260209091200154600e805473ffffffffffffffffffffffffffffffffffffffff909216945090829081106114c1576114c16144b2565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff8681168452600c8352604080852080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001690559116808452600b90925290912080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905591508061155781614557565b91505061145d565b5061156c600d6000613b9d565b611578600e6000613b9d565b604080516080810182526000808252602082018190529181018290526060810182905290805b8c518110156119e157600c60008e83815181106115bd576115bd6144b2565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528101919091526040016000205460ff1615611628576040517f77cea0fa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168d8281518110611652576116526144b2565b602002602001015173ffffffffffffffffffffffffffffffffffffffff16036116a7576040517f815e1d6400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405180604001604052806001151581526020018260ff16815250600c60008f84815181106116d8576116d86144b2565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528181019290925260400160002082518154939092015160ff16610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909316929092171790558b518c9082908110611780576117806144b2565b60200260200101519150600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036117f0576040517f58a70a0a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff82166000908152600b60209081526040918290208251608081018452905460ff80821615801584526101008304909116938301939093526bffffffffffffffffffffffff6201000082048116948301949094526e010000000000000000000000000000900490921660608301529093506118ab576040517f6a7281ad00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001835260ff80821660208086019182526bffffffffffffffffffffffff808b166060880190815273ffffffffffffffffffffffffffffffffffffffff87166000908152600b909352604092839020885181549551948a0151925184166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff939094166201000002929092167fffffffffffff000000000000000000000000000000000000000000000000ffff94909616610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000090951694909417179190911692909217919091179055806119d981614557565b91505061159e565b50508a516119f79150600d9060208d0190613bbb565b508851611a0b90600e9060208c0190613bbb565b506040518061012001604052808960ff168152602001886000015163ffffffff168152602001886020015163ffffffff168152602001886060015162ffffff168152602001886080015161ffff1681526020016012600001600e9054906101000a900460ff16151581526020016012600001600f9054906101000a900460ff1615158152602001856bffffffffffffffffffffffff168152602001600063ffffffff16815250601260008201518160000160006101000a81548160ff021916908360ff16021790555060208201518160000160016101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160056101000a81548163ffffffff021916908363ffffffff16021790555060608201518160000160096101000a81548162ffffff021916908362ffffff160217905550608082015181600001600c6101000a81548161ffff021916908361ffff16021790555060a082015181600001600e6101000a81548160ff02191690831515021790555060c082015181600001600f6101000a81548160ff02191690831515021790555060e08201518160000160106101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555061010082015181600001601c6101000a81548163ffffffff021916908363ffffffff1602179055509050506040518061018001604052808860a001516bffffffffffffffffffffffff16815260200188610180015173ffffffffffffffffffffffffffffffffffffffff168152602001601360010160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff168152602001886040015163ffffffff1681526020018860c0015163ffffffff168152602001601360010160149054906101000a900463ffffffff1663ffffffff168152602001601360010160189054906101000a900463ffffffff1663ffffffff1681526020016013600101601c9054906101000a900463ffffffff1663ffffffff1681526020018860e0015163ffffffff16815260200188610100015163ffffffff16815260200188610120015163ffffffff168152602001886101c0015173ffffffffffffffffffffffffffffffffffffffff16815250601360008201518160000160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550602082015181600001600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060408201518160010160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550606082015181600101600c6101000a81548163ffffffff021916908363ffffffff16021790555060808201518160010160106101000a81548163ffffffff021916908363ffffffff16021790555060a08201518160010160146101000a81548163ffffffff021916908363ffffffff16021790555060c08201518160010160186101000a81548163ffffffff021916908363ffffffff16021790555060e082015181600101601c6101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160020160006101000a81548163ffffffff021916908363ffffffff1602179055506101208201518160020160046101000a81548163ffffffff021916908363ffffffff1602179055506101408201518160020160086101000a81548163ffffffff021916908363ffffffff16021790555061016082015181600201600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505086610140015160168190555086610160015160178190555060006013600101601c9054906101000a900463ffffffff169050611fcd612eb1565b601480547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167c010000000000000000000000000000000000000000000000000000000063ffffffff93841602178082556001926018916120489185917801000000000000000000000000000000000000000000000000900416614675565b92506101000a81548163ffffffff021916908363ffffffff16021790555060008860405160200161207991906146e3565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529190526014549091506120e290469030907801000000000000000000000000000000000000000000000000900463ffffffff168f8f8f878f8f612f66565b60115560005b6120f26009613010565b8110156121225761210f61210760098361301a565b600990613026565b508061211a81614557565b9150506120e8565b5060005b896101a0015151811015612179576121668a6101a00151828151811061214e5761214e6144b2565b6020026020010151600961304890919063ffffffff16565b508061217181614557565b915050612126565b507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0582601154601360010160189054906101000a900463ffffffff168f8f8f878f8f6040516121d099989796959493929190614847565b60405180910390a1505050505050505050505050565b612207868686868060200190518101906122009190614978565b86866112ea565b505050505050565b612217612c26565b6122208161306a565b50565b321561225b576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60125460009081906f01000000000000000000000000000000900460ff16156122b2576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601280547fffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffff166f010000000000000000000000000000001790556040517f4585e33b000000000000000000000000000000000000000000000000000000009061231f908590602401613cc8565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009094169390931790925290517f79188d1600000000000000000000000000000000000000000000000000000000815290935073ffffffffffffffffffffffffffffffffffffffff8616906379188d16906123f29087908790600401614ad2565b60408051808303816000875af1158015612410573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906124349190614aeb565b601280547fffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffff16905590969095509350505050565b6000878760405161247a929190614b1e565b604051908190038120612491918b90602001614b2e565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201208383019092526000808452908301819052909250906000805b88811015612668576001858783602081106124fd576124fd6144b2565b61250a91901a601b6144e1565b8c8c8581811061251c5761251c6144b2565b905060200201358b8b86818110612535576125356144b2565b9050602002013560405160008152602001604052604051612572949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015612594573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff81166000908152600c602090815290849020838501909452925460ff8082161515808552610100909204169383019390935290955093509050612642576040517f0f4c073700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826020015160080260ff166001901b84019350808061266090614557565b9150506124e0565b50827e010101010101010101010101010101010101010101010101010101010101018416146126c3576040517fc103be2e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505050505050505050505050565b61270a6040518060c001604052806000815260200160008152602001606081526020016060815260200160608152602001606081525090565b600061271883850185614c1f565b604081015151606082015151919250908114158061273b57508082608001515114155b8061274b5750808260a001515114155b15612782576040517fb55ac75400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5090505b92915050565b6000818160045b600f811015612819577fff0000000000000000000000000000000000000000000000000000000000000082168382602081106127d1576127d16144b2565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161461280757506000949350505050565b8061281181614557565b915050612793565b5081600f1a600181111561282f5761282f6144fa565b949350505050565b60008061284988878b6000015161315f565b90506000806128648b8a63ffffffff16858a8a60018b6131eb565b9092509050612873818361447a565b9b9a5050505050505050505050565b60008080808460800151600181111561289d5761289d6144fa565b036128c1576128ad868686613644565b6128bc5760009250905061079f565b612938565b6001846080015160018111156128d9576128d96144fa565b036129065760006128eb878787613738565b9250905080612900575060009250905061079f565b50612938565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612940612eb1565b84516040015163ffffffff161161299457857fc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd5636866040516129819190613cc8565b60405180910390a260009250905061079f565b83604001516bffffffffffffffffffffffff16846000015160a001516bffffffffffffffffffffffff1610156129f457857f377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e02866040516129819190613cc8565b6001969095509350505050565b600081608001516001811115612a1957612a196144fa565b03612a8b57612a26612eb1565b6000838152600460205260409020600101805463ffffffff929092167801000000000000000000000000000000000000000000000000027fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff9092169190911790555050565b600181608001516001811115612aa357612aa36144fa565b03612b0f5760e08101805160009081526008602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055915191517fa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f29190a25b5050565b6000612b2084848461315f565b90508085101561282f5750929392505050565b600080612b4e888760a001518860c0015188888860016131eb565b90925090506000612b5f828461447a565b600089815260046020526040902060010180549192508291600c90612ba39084906c0100000000000000000000000090046bffffffffffffffffffffffff16614d0c565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915560008a815260046020526040812060010180548594509092612bec9185911661447a565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555050965096945050505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461225b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161039e565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e0100000000000000000000000000009004909116606082015290612ea3576000816060015185612d3f9190614d0c565b90506000612d4d8583614d31565b90508083604001818151612d61919061447a565b6bffffffffffffffffffffffff16905250612d7c8582614d5c565b83606001818151612d8d919061447a565b6bffffffffffffffffffffffff90811690915273ffffffffffffffffffffffffffffffffffffffff89166000908152600b602090815260409182902087518154928901519389015160608a015186166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff919096166201000002167fffffffffffff000000000000000000000000000000000000000000000000ffff60ff95909516610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909416939093171792909216179190911790555050505b6040015190505b9392505050565b600060017f00000000000000000000000000000000000000000000000000000000000000006002811115612ee757612ee76144fa565b03612f6157606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015612f38573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612f5c9190614d90565b905090565b504390565b6000808a8a8a8a8a8a8a8a8a604051602001612f8a99989796959493929190614da9565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179b9a5050505050505050505050565b6000612786825490565b6000612eaa83836138d0565b6000612eaa8373ffffffffffffffffffffffffffffffffffffffff84166138fa565b6000612eaa8373ffffffffffffffffffffffffffffffffffffffff84166139f4565b3373ffffffffffffffffffffffffffffffffffffffff8216036130e9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161039e565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008080856001811115613175576131756144fa565b03613184575062015f906131a3565b6001856001811115613198576131986144fa565b0361290657506201adb05b6131b463ffffffff8516601461458f565b6131bf8460016144e1565b6131ce9060ff16611d4c61458f565b6131d8908361449f565b6131e2919061449f565b95945050505050565b6000806000896080015161ffff1687613204919061458f565b90508380156132125750803a105b1561321a57503a5b600060027f00000000000000000000000000000000000000000000000000000000000000006002811115613250576132506144fa565b036133af5760408051600081526020810190915285156132ae57600036604051806080016040528060488152602001614fdb6048913960405160200161329893929190614e3e565b6040516020818303038152906040529050613316565b6015546132ca90640100000000900463ffffffff166004614e65565b63ffffffff1667ffffffffffffffff8111156132e8576132e8613ef5565b6040519080825280601f01601f191660200182016040528015613312576020820181803683370190505b5090505b6040517f49948e0e00000000000000000000000000000000000000000000000000000000815273420000000000000000000000000000000000000f906349948e0e90613366908490600401613cc8565b602060405180830381865afa158015613383573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906133a79190614d90565b915050613509565b60017f000000000000000000000000000000000000000000000000000000000000000060028111156133e3576133e36144fa565b0361350957841561346557606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561343a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061345e9190614d90565b9050613509565b6000606c73ffffffffffffffffffffffffffffffffffffffff166341b247a86040518163ffffffff1660e01b815260040160c060405180830381865afa1580156134b3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906134d79190614e88565b50506015549294506134fa93505050640100000000900463ffffffff168261458f565b61350590601061458f565b9150505b8461352557808b6080015161ffff16613522919061458f565b90505b61353361ffff8716826145fb565b9050600087826135438c8e61449f565b61354d908661458f565b613557919061449f565b61356990670de0b6b3a764000061458f565b61357391906145fb565b905060008c6040015163ffffffff1664e8d4a51000613592919061458f565b898e6020015163ffffffff16858f886135ab919061458f565b6135b5919061449f565b6135c390633b9aca0061458f565b6135cd919061458f565b6135d791906145fb565b6135e1919061449f565b90506b033b2e3c9fd0803ce80000006135fa828461449f565b1115613632576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b6000808380602001905181019061365b9190614ed2565b835160c00151815191925063ffffffff908116911610156136b857847f405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8856040516136a69190613cc8565b60405180910390a26000915050612eaa565b6020810151158015906136df5750602081015181516136dc9063ffffffff16613a43565b14155b806136f857506136ed612eb1565b815163ffffffff1610155b1561372d57847f6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301856040516136a69190613cc8565b506001949350505050565b6000806000848060200190518101906137519190614f2a565b90506000868260000151836020015184604001516040516020016137b394939291909384526020840192909252604083015260e01b7fffffffff0000000000000000000000000000000000000000000000000000000016606082015260640190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052805160209091012060808301519091501580159061381557508160800151613812836060015163ffffffff16613a43565b14155b806138315750613823612eb1565b826060015163ffffffff1610155b1561387b57867f6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301876040516138669190613cc8565b60405180910390a260009350915061079f9050565b60008181526008602052604090205460ff16156138c257867f405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8876040516138669190613cc8565b600197909650945050505050565b60008260000182815481106138e7576138e76144b2565b9060005260206000200154905092915050565b600081815260018301602052604081205480156139e357600061391e600183614544565b855490915060009061393290600190614544565b9050818114613997576000866000018281548110613952576139526144b2565b9060005260206000200154905080876000018481548110613975576139756144b2565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806139a8576139a8614fab565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050612786565b6000915050612786565b5092915050565b6000818152600183016020526040812054613a3b57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155612786565b506000612786565b600060017f00000000000000000000000000000000000000000000000000000000000000006002811115613a7957613a796144fa565b03613b93576000606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613acc573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613af09190614d90565b90508083101580613b0b5750610100613b098483614544565b115b15613b195750600092915050565b6040517f2b407a8200000000000000000000000000000000000000000000000000000000815260048101849052606490632b407a8290602401602060405180830381865afa158015613b6f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612eaa9190614d90565b504090565b919050565b50805460008255906000526020600020908101906122209190613c45565b828054828255906000526020600020908101928215613c35579160200282015b82811115613c3557825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190613bdb565b50613c41929150613c45565b5090565b5b80821115613c415760008155600101613c46565b60005b83811015613c75578181015183820152602001613c5d565b50506000910152565b60008151808452613c96816020860160208601613c5a565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000612eaa6020830184613c7e565b73ffffffffffffffffffffffffffffffffffffffff8116811461222057600080fd5b8035613b9881613cdb565b60008083601f840112613d1a57600080fd5b50813567ffffffffffffffff811115613d3257600080fd5b602083019150836020828501011115613d4a57600080fd5b9250929050565b60008060008060608587031215613d6757600080fd5b8435613d7281613cdb565b935060208501359250604085013567ffffffffffffffff811115613d9557600080fd5b613da187828801613d08565b95989497509550505050565b600080600060408486031215613dc257600080fd5b83359250602084013567ffffffffffffffff811115613de057600080fd5b613dec86828701613d08565b9497909650939450505050565b60008083601f840112613e0b57600080fd5b50813567ffffffffffffffff811115613e2357600080fd5b6020830191508360208260051b8501011115613d4a57600080fd5b60008060008060008060008060e0898b031215613e5a57600080fd5b606089018a811115613e6b57600080fd5b8998503567ffffffffffffffff80821115613e8557600080fd5b613e918c838d01613d08565b909950975060808b0135915080821115613eaa57600080fd5b613eb68c838d01613df9565b909750955060a08b0135915080821115613ecf57600080fd5b50613edc8b828c01613df9565b999c989b50969995989497949560c00135949350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516101e0810167ffffffffffffffff81118282101715613f4857613f48613ef5565b60405290565b60405160c0810167ffffffffffffffff81118282101715613f4857613f48613ef5565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613fb857613fb8613ef5565b604052919050565b600067ffffffffffffffff821115613fda57613fda613ef5565b5060051b60200190565b600082601f830112613ff557600080fd5b8135602061400a61400583613fc0565b613f71565b82815260059290921b8401810191818101908684111561402957600080fd5b8286015b8481101561404d57803561404081613cdb565b835291830191830161402d565b509695505050505050565b803560ff81168114613b9857600080fd5b63ffffffff8116811461222057600080fd5b8035613b9881614069565b62ffffff8116811461222057600080fd5b8035613b9881614086565b61ffff8116811461222057600080fd5b8035613b98816140a2565b6bffffffffffffffffffffffff8116811461222057600080fd5b8035613b98816140bd565b60006101e082840312156140f557600080fd5b6140fd613f24565b90506141088261407b565b81526141166020830161407b565b60208201526141276040830161407b565b604082015261413860608301614097565b6060820152614149608083016140b2565b608082015261415a60a083016140d7565b60a082015261416b60c0830161407b565b60c082015261417c60e0830161407b565b60e082015261010061418f81840161407b565b908201526101206141a183820161407b565b90820152610140828101359082015261016080830135908201526101806141c9818401613cfd565b908201526101a08281013567ffffffffffffffff8111156141e957600080fd5b6141f585828601613fe4565b8284015250506101c0614209818401613cfd565b9082015292915050565b803567ffffffffffffffff81168114613b9857600080fd5b600082601f83011261423c57600080fd5b813567ffffffffffffffff81111561425657614256613ef5565b61428760207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601613f71565b81815284602083860101111561429c57600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060c087890312156142d257600080fd5b863567ffffffffffffffff808211156142ea57600080fd5b6142f68a838b01613fe4565b9750602089013591508082111561430c57600080fd5b6143188a838b01613fe4565b965061432660408a01614058565b9550606089013591508082111561433c57600080fd5b6143488a838b016140e2565b945061435660808a01614213565b935060a089013591508082111561436c57600080fd5b5061437989828a0161422b565b9150509295509295509295565b60008060008060008060c0878903121561439f57600080fd5b863567ffffffffffffffff808211156143b757600080fd5b6143c38a838b01613fe4565b975060208901359150808211156143d957600080fd5b6143e58a838b01613fe4565b96506143f360408a01614058565b9550606089013591508082111561440957600080fd5b6143488a838b0161422b565b60006020828403121561442757600080fd5b8135612eaa81613cdb565b60006020828403121561444457600080fd5b5035919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6bffffffffffffffffffffffff8181168382160190808211156139ed576139ed61444b565b808201808211156127865761278661444b565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60ff81811683821601908111156127865761278661444b565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b61ffff8181168382160190808211156139ed576139ed61444b565b818103818111156127865761278661444b565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036145885761458861444b565b5060010190565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156145c7576145c761444b565b500290565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60008261460a5761460a6145cc565b500490565b6bffffffffffffffffffffffff851681528360208201528260408201526080606082015260006146426080830184613c7e565b9695505050505050565b600060ff821660ff84168160ff048111821515161561466d5761466d61444b565b029392505050565b63ffffffff8181168382160190808211156139ed576139ed61444b565b600081518084526020808501945080840160005b838110156146d857815173ffffffffffffffffffffffffffffffffffffffff16875295820195908201906001016146a6565b509495945050505050565b602081526146fa60208201835163ffffffff169052565b60006020830151614713604084018263ffffffff169052565b50604083015163ffffffff8116606084015250606083015162ffffff8116608084015250608083015161ffff811660a08401525060a08301516bffffffffffffffffffffffff811660c08401525060c083015163ffffffff811660e08401525060e083015161010061478c8185018363ffffffff169052565b84015190506101206147a58482018363ffffffff169052565b84015190506101406147be8482018363ffffffff169052565b840151610160848101919091528401516101808085019190915284015190506101a06148018185018373ffffffffffffffffffffffffffffffffffffffff169052565b808501519150506101e06101c08181860152614821610200860184614692565b95015173ffffffffffffffffffffffffffffffffffffffff169301929092525090919050565b600061012063ffffffff808d1684528b6020850152808b166040850152508060608401526148778184018a614692565b9050828103608084015261488b8189614692565b905060ff871660a084015282810360c08401526148a88187613c7e565b905067ffffffffffffffff851660e08401528281036101008401526148cd8185613c7e565b9c9b505050505050505050505050565b8051613b9881614069565b8051613b9881614086565b8051613b98816140a2565b8051613b98816140bd565b8051613b9881613cdb565b600082601f83011261492557600080fd5b8151602061493561400583613fc0565b82815260059290921b8401810191818101908684111561495457600080fd5b8286015b8481101561404d57805161496b81613cdb565b8352918301918301614958565b60006020828403121561498a57600080fd5b815167ffffffffffffffff808211156149a257600080fd5b908301906101e082860312156149b757600080fd5b6149bf613f24565b6149c8836148dd565b81526149d6602084016148dd565b60208201526149e7604084016148dd565b60408201526149f8606084016148e8565b6060820152614a09608084016148f3565b6080820152614a1a60a084016148fe565b60a0820152614a2b60c084016148dd565b60c0820152614a3c60e084016148dd565b60e0820152610100614a4f8185016148dd565b90820152610120614a618482016148dd565b9082015261014083810151908201526101608084015190820152610180614a89818501614909565b908201526101a08381015183811115614aa157600080fd5b614aad88828701614914565b8284015250506101c09150614ac3828401614909565b91810191909152949350505050565b82815260406020820152600061282f6040830184613c7e565b60008060408385031215614afe57600080fd5b82518015158114614b0e57600080fd5b6020939093015192949293505050565b8183823760009101908152919050565b8281526080810160608360208401379392505050565b600082601f830112614b5557600080fd5b81356020614b6561400583613fc0565b82815260059290921b84018101918181019086841115614b8457600080fd5b8286015b8481101561404d5780358352918301918301614b88565b600082601f830112614bb057600080fd5b81356020614bc061400583613fc0565b82815260059290921b84018101918181019086841115614bdf57600080fd5b8286015b8481101561404d57803567ffffffffffffffff811115614c035760008081fd5b614c118986838b010161422b565b845250918301918301614be3565b600060208284031215614c3157600080fd5b813567ffffffffffffffff80821115614c4957600080fd5b9083019060c08286031215614c5d57600080fd5b614c65613f4e565b8235815260208301356020820152604083013582811115614c8557600080fd5b614c9187828601614b44565b604083015250606083013582811115614ca957600080fd5b614cb587828601614b44565b606083015250608083013582811115614ccd57600080fd5b614cd987828601614b9f565b60808301525060a083013582811115614cf157600080fd5b614cfd87828601614b9f565b60a08301525095945050505050565b6bffffffffffffffffffffffff8281168282160390808211156139ed576139ed61444b565b60006bffffffffffffffffffffffff80841680614d5057614d506145cc565b92169190910492915050565b60006bffffffffffffffffffffffff80831681851681830481118215151615614d8757614d8761444b565b02949350505050565b600060208284031215614da257600080fd5b5051919050565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b166040850152816060850152614df08285018b614692565b91508382036080850152614e04828a614692565b915060ff881660a085015283820360c0850152614e218288613c7e565b90861660e085015283810361010085015290506148cd8185613c7e565b828482376000838201600081528351614e5b818360208801613c5a565b0195945050505050565b600063ffffffff80831681851681830481118215151615614d8757614d8761444b565b60008060008060008060c08789031215614ea157600080fd5b865195506020870151945060408701519350606087015192506080870151915060a087015190509295509295509295565b600060408284031215614ee457600080fd5b6040516040810181811067ffffffffffffffff82111715614f0757614f07613ef5565b6040528251614f1581614069565b81526020928301519281019290925250919050565b600060a08284031215614f3c57600080fd5b60405160a0810181811067ffffffffffffffff82111715614f5f57614f5f613ef5565b806040525082518152602083015160208201526040830151614f8081614069565b60408201526060830151614f9381614069565b60608201526080928301519281019290925250919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfe307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000810000a",
}

var KeeperRegistryABI = KeeperRegistryMetaData.ABI

var KeeperRegistryBin = KeeperRegistryMetaData.Bin

func DeployKeeperRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, logicA common.Address) (common.Address, *types.Transaction, *KeeperRegistry, error) {
	parsed, err := KeeperRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistryBin), backend, logicA)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperRegistry{address: address, abi: *parsed, KeeperRegistryCaller: KeeperRegistryCaller{contract: contract}, KeeperRegistryTransactor: KeeperRegistryTransactor{contract: contract}, KeeperRegistryFilterer: KeeperRegistryFilterer{contract: contract}}, nil
}

type KeeperRegistry struct {
	address common.Address
	abi     abi.ABI
	KeeperRegistryCaller
	KeeperRegistryTransactor
	KeeperRegistryFilterer
}

type KeeperRegistryCaller struct {
	contract *bind.BoundContract
}

type KeeperRegistryTransactor struct {
	contract *bind.BoundContract
}

type KeeperRegistryFilterer struct {
	contract *bind.BoundContract
}

type KeeperRegistrySession struct {
	Contract     *KeeperRegistry
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeeperRegistryCallerSession struct {
	Contract *KeeperRegistryCaller
	CallOpts bind.CallOpts
}

type KeeperRegistryTransactorSession struct {
	Contract     *KeeperRegistryTransactor
	TransactOpts bind.TransactOpts
}

type KeeperRegistryRaw struct {
	Contract *KeeperRegistry
}

type KeeperRegistryCallerRaw struct {
	Contract *KeeperRegistryCaller
}

type KeeperRegistryTransactorRaw struct {
	Contract *KeeperRegistryTransactor
}

func NewKeeperRegistry(address common.Address, backend bind.ContractBackend) (*KeeperRegistry, error) {
	abi, err := abi.JSON(strings.NewReader(KeeperRegistryABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeeperRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry{address: address, abi: abi, KeeperRegistryCaller: KeeperRegistryCaller{contract: contract}, KeeperRegistryTransactor: KeeperRegistryTransactor{contract: contract}, KeeperRegistryFilterer: KeeperRegistryFilterer{contract: contract}}, nil
}

func NewKeeperRegistryCaller(address common.Address, caller bind.ContractCaller) (*KeeperRegistryCaller, error) {
	contract, err := bindKeeperRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCaller{contract: contract}, nil
}

func NewKeeperRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*KeeperRegistryTransactor, error) {
	contract, err := bindKeeperRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryTransactor{contract: contract}, nil
}

func NewKeeperRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*KeeperRegistryFilterer, error) {
	contract, err := bindKeeperRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryFilterer{contract: contract}, nil
}

func bindKeeperRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeeperRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_KeeperRegistry *KeeperRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistry.Contract.KeeperRegistryCaller.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistry *KeeperRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.KeeperRegistryTransactor.contract.Transfer(opts)
}

func (_KeeperRegistry *KeeperRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.KeeperRegistryTransactor.contract.Transact(opts, method, params...)
}

func (_KeeperRegistry *KeeperRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistry.Contract.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistry *KeeperRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.contract.Transfer(opts)
}

func (_KeeperRegistry *KeeperRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.contract.Transact(opts, method, params...)
}

func (_KeeperRegistry *KeeperRegistryCaller) FallbackTo(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "fallbackTo")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) FallbackTo() (common.Address, error) {
	return _KeeperRegistry.Contract.FallbackTo(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) FallbackTo() (common.Address, error) {
	return _KeeperRegistry.Contract.FallbackTo(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCaller) LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_KeeperRegistry *KeeperRegistrySession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _KeeperRegistry.Contract.LatestConfigDetails(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _KeeperRegistry.Contract.LatestConfigDetails(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

	error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(LatestConfigDigestAndEpoch)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_KeeperRegistry *KeeperRegistrySession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _KeeperRegistry.Contract.LatestConfigDigestAndEpoch(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _KeeperRegistry.Contract.LatestConfigDigestAndEpoch(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) Owner() (common.Address, error) {
	return _KeeperRegistry.Contract.Owner(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) Owner() (common.Address, error) {
	return _KeeperRegistry.Contract.Owner(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) TypeAndVersion() (string, error) {
	return _KeeperRegistry.Contract.TypeAndVersion(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) TypeAndVersion() (string, error) {
	return _KeeperRegistry.Contract.TypeAndVersion(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "acceptOwnership")
}

func (_KeeperRegistry *KeeperRegistrySession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.AcceptOwnership(&_KeeperRegistry.TransactOpts)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.AcceptOwnership(&_KeeperRegistry.TransactOpts)
}

func (_KeeperRegistry *KeeperRegistryTransactor) OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "onTokenTransfer", sender, amount, data)
}

func (_KeeperRegistry *KeeperRegistrySession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.OnTokenTransfer(&_KeeperRegistry.TransactOpts, sender, amount, data)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.OnTokenTransfer(&_KeeperRegistry.TransactOpts, sender, amount, data)
}

func (_KeeperRegistry *KeeperRegistryTransactor) SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "setConfig", signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_KeeperRegistry *KeeperRegistrySession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetConfig(&_KeeperRegistry.TransactOpts, signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetConfig(&_KeeperRegistry.TransactOpts, signers, transmitters, f, onchainConfigBytes, offchainConfigVersion, offchainConfig)
}

func (_KeeperRegistry *KeeperRegistryTransactor) SetConfigTypeSafe(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig KeeperRegistryBase21OnchainConfig, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "setConfigTypeSafe", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_KeeperRegistry *KeeperRegistrySession) SetConfigTypeSafe(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig KeeperRegistryBase21OnchainConfig, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetConfigTypeSafe(&_KeeperRegistry.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) SetConfigTypeSafe(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig KeeperRegistryBase21OnchainConfig, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetConfigTypeSafe(&_KeeperRegistry.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_KeeperRegistry *KeeperRegistryTransactor) SimulatePerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "simulatePerformUpkeep", id, performData)
}

func (_KeeperRegistry *KeeperRegistrySession) SimulatePerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SimulatePerformUpkeep(&_KeeperRegistry.TransactOpts, id, performData)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) SimulatePerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SimulatePerformUpkeep(&_KeeperRegistry.TransactOpts, id, performData)
}

func (_KeeperRegistry *KeeperRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "transferOwnership", to)
}

func (_KeeperRegistry *KeeperRegistrySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.TransferOwnership(&_KeeperRegistry.TransactOpts, to)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.TransferOwnership(&_KeeperRegistry.TransactOpts, to)
}

func (_KeeperRegistry *KeeperRegistryTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "transmit", reportContext, rawReport, rs, ss, rawVs)
}

func (_KeeperRegistry *KeeperRegistrySession) Transmit(reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.Transmit(&_KeeperRegistry.TransactOpts, reportContext, rawReport, rs, ss, rawVs)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) Transmit(reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.Transmit(&_KeeperRegistry.TransactOpts, reportContext, rawReport, rs, ss, rawVs)
}

func (_KeeperRegistry *KeeperRegistryTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _KeeperRegistry.contract.RawTransact(opts, calldata)
}

func (_KeeperRegistry *KeeperRegistrySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.Fallback(&_KeeperRegistry.TransactOpts, calldata)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.Fallback(&_KeeperRegistry.TransactOpts, calldata)
}

type KeeperRegistryAdminPrivilegeConfigSetIterator struct {
	Event *KeeperRegistryAdminPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryAdminPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryAdminPrivilegeConfigSet)
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
		it.Event = new(KeeperRegistryAdminPrivilegeConfigSet)
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

func (it *KeeperRegistryAdminPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryAdminPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryAdminPrivilegeConfigSet struct {
	Admin           common.Address
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*KeeperRegistryAdminPrivilegeConfigSetIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryAdminPrivilegeConfigSetIterator{contract: _KeeperRegistry.contract, event: "AdminPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryAdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "AdminPrivilegeConfigSet", adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryAdminPrivilegeConfigSet)
				if err := _KeeperRegistry.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseAdminPrivilegeConfigSet(log types.Log) (*KeeperRegistryAdminPrivilegeConfigSet, error) {
	event := new(KeeperRegistryAdminPrivilegeConfigSet)
	if err := _KeeperRegistry.contract.UnpackLog(event, "AdminPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryCancelledUpkeepReportIterator struct {
	Event *KeeperRegistryCancelledUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryCancelledUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryCancelledUpkeepReport)
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
		it.Event = new(KeeperRegistryCancelledUpkeepReport)
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

func (it *KeeperRegistryCancelledUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryCancelledUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryCancelledUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryCancelledUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryCancelledUpkeepReportIterator{contract: _KeeperRegistry.contract, event: "CancelledUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryCancelledUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryCancelledUpkeepReport)
				if err := _KeeperRegistry.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseCancelledUpkeepReport(log types.Log) (*KeeperRegistryCancelledUpkeepReport, error) {
	event := new(KeeperRegistryCancelledUpkeepReport)
	if err := _KeeperRegistry.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryConfigSetIterator struct {
	Event *KeeperRegistryConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryConfigSet)
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
		it.Event = new(KeeperRegistryConfigSet)
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

func (it *KeeperRegistryConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryConfigSet struct {
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

func (_KeeperRegistry *KeeperRegistryFilterer) FilterConfigSet(opts *bind.FilterOpts) (*KeeperRegistryConfigSetIterator, error) {

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryConfigSetIterator{contract: _KeeperRegistry.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryConfigSet) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryConfigSet)
				if err := _KeeperRegistry.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseConfigSet(log types.Log) (*KeeperRegistryConfigSet, error) {
	event := new(KeeperRegistryConfigSet)
	if err := _KeeperRegistry.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryDedupKeyAddedIterator struct {
	Event *KeeperRegistryDedupKeyAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryDedupKeyAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryDedupKeyAdded)
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
		it.Event = new(KeeperRegistryDedupKeyAdded)
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

func (it *KeeperRegistryDedupKeyAddedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryDedupKeyAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryDedupKeyAdded struct {
	DedupKey [32]byte
	Raw      types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*KeeperRegistryDedupKeyAddedIterator, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryDedupKeyAddedIterator{contract: _KeeperRegistry.contract, event: "DedupKeyAdded", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryDedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error) {

	var dedupKeyRule []interface{}
	for _, dedupKeyItem := range dedupKey {
		dedupKeyRule = append(dedupKeyRule, dedupKeyItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "DedupKeyAdded", dedupKeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryDedupKeyAdded)
				if err := _KeeperRegistry.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseDedupKeyAdded(log types.Log) (*KeeperRegistryDedupKeyAdded, error) {
	event := new(KeeperRegistryDedupKeyAdded)
	if err := _KeeperRegistry.contract.UnpackLog(event, "DedupKeyAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryFundsAddedIterator struct {
	Event *KeeperRegistryFundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryFundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryFundsAdded)
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
		it.Event = new(KeeperRegistryFundsAdded)
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

func (it *KeeperRegistryFundsAddedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryFundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryFundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistryFundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryFundsAddedIterator{contract: _KeeperRegistry.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryFundsAdded)
				if err := _KeeperRegistry.contract.UnpackLog(event, "FundsAdded", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseFundsAdded(log types.Log) (*KeeperRegistryFundsAdded, error) {
	event := new(KeeperRegistryFundsAdded)
	if err := _KeeperRegistry.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryFundsWithdrawnIterator struct {
	Event *KeeperRegistryFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryFundsWithdrawn)
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
		it.Event = new(KeeperRegistryFundsWithdrawn)
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

func (it *KeeperRegistryFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryFundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryFundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryFundsWithdrawnIterator{contract: _KeeperRegistry.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryFundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryFundsWithdrawn)
				if err := _KeeperRegistry.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseFundsWithdrawn(log types.Log) (*KeeperRegistryFundsWithdrawn, error) {
	event := new(KeeperRegistryFundsWithdrawn)
	if err := _KeeperRegistry.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryInsufficientFundsUpkeepReportIterator struct {
	Event *KeeperRegistryInsufficientFundsUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryInsufficientFundsUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryInsufficientFundsUpkeepReport)
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
		it.Event = new(KeeperRegistryInsufficientFundsUpkeepReport)
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

func (it *KeeperRegistryInsufficientFundsUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryInsufficientFundsUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryInsufficientFundsUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryInsufficientFundsUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryInsufficientFundsUpkeepReportIterator{contract: _KeeperRegistry.contract, event: "InsufficientFundsUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryInsufficientFundsUpkeepReport)
				if err := _KeeperRegistry.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseInsufficientFundsUpkeepReport(log types.Log) (*KeeperRegistryInsufficientFundsUpkeepReport, error) {
	event := new(KeeperRegistryInsufficientFundsUpkeepReport)
	if err := _KeeperRegistry.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryOwnerFundsWithdrawnIterator struct {
	Event *KeeperRegistryOwnerFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryOwnerFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryOwnerFundsWithdrawn)
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
		it.Event = new(KeeperRegistryOwnerFundsWithdrawn)
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

func (it *KeeperRegistryOwnerFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryOwnerFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryOwnerFundsWithdrawn struct {
	Amount *big.Int
	Raw    types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*KeeperRegistryOwnerFundsWithdrawnIterator, error) {

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryOwnerFundsWithdrawnIterator{contract: _KeeperRegistry.contract, event: "OwnerFundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryOwnerFundsWithdrawn) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryOwnerFundsWithdrawn)
				if err := _KeeperRegistry.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseOwnerFundsWithdrawn(log types.Log) (*KeeperRegistryOwnerFundsWithdrawn, error) {
	event := new(KeeperRegistryOwnerFundsWithdrawn)
	if err := _KeeperRegistry.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryOwnershipTransferRequestedIterator struct {
	Event *KeeperRegistryOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryOwnershipTransferRequested)
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
		it.Event = new(KeeperRegistryOwnershipTransferRequested)
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

func (it *KeeperRegistryOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryOwnershipTransferRequestedIterator{contract: _KeeperRegistry.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryOwnershipTransferRequested)
				if err := _KeeperRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryOwnershipTransferRequested, error) {
	event := new(KeeperRegistryOwnershipTransferRequested)
	if err := _KeeperRegistry.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryOwnershipTransferredIterator struct {
	Event *KeeperRegistryOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryOwnershipTransferred)
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
		it.Event = new(KeeperRegistryOwnershipTransferred)
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

func (it *KeeperRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryOwnershipTransferredIterator{contract: _KeeperRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryOwnershipTransferred)
				if err := _KeeperRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*KeeperRegistryOwnershipTransferred, error) {
	event := new(KeeperRegistryOwnershipTransferred)
	if err := _KeeperRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryPausedIterator struct {
	Event *KeeperRegistryPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryPaused)
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
		it.Event = new(KeeperRegistryPaused)
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

func (it *KeeperRegistryPausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterPaused(opts *bind.FilterOpts) (*KeeperRegistryPausedIterator, error) {

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryPausedIterator{contract: _KeeperRegistry.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPaused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryPaused)
				if err := _KeeperRegistry.contract.UnpackLog(event, "Paused", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParsePaused(log types.Log) (*KeeperRegistryPaused, error) {
	event := new(KeeperRegistryPaused)
	if err := _KeeperRegistry.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryPayeesUpdatedIterator struct {
	Event *KeeperRegistryPayeesUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryPayeesUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryPayeesUpdated)
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
		it.Event = new(KeeperRegistryPayeesUpdated)
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

func (it *KeeperRegistryPayeesUpdatedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryPayeesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryPayeesUpdated struct {
	Transmitters []common.Address
	Payees       []common.Address
	Raw          types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterPayeesUpdated(opts *bind.FilterOpts) (*KeeperRegistryPayeesUpdatedIterator, error) {

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryPayeesUpdatedIterator{contract: _KeeperRegistry.contract, event: "PayeesUpdated", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPayeesUpdated) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryPayeesUpdated)
				if err := _KeeperRegistry.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParsePayeesUpdated(log types.Log) (*KeeperRegistryPayeesUpdated, error) {
	event := new(KeeperRegistryPayeesUpdated)
	if err := _KeeperRegistry.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryPayeeshipTransferRequestedIterator struct {
	Event *KeeperRegistryPayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryPayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryPayeeshipTransferRequested)
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
		it.Event = new(KeeperRegistryPayeeshipTransferRequested)
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

func (it *KeeperRegistryPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryPayeeshipTransferRequested struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryPayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryPayeeshipTransferRequestedIterator{contract: _KeeperRegistry.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryPayeeshipTransferRequested)
				if err := _KeeperRegistry.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistryPayeeshipTransferRequested, error) {
	event := new(KeeperRegistryPayeeshipTransferRequested)
	if err := _KeeperRegistry.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryPayeeshipTransferredIterator struct {
	Event *KeeperRegistryPayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryPayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryPayeeshipTransferred)
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
		it.Event = new(KeeperRegistryPayeeshipTransferred)
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

func (it *KeeperRegistryPayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryPayeeshipTransferred struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryPayeeshipTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryPayeeshipTransferredIterator{contract: _KeeperRegistry.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryPayeeshipTransferred)
				if err := _KeeperRegistry.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParsePayeeshipTransferred(log types.Log) (*KeeperRegistryPayeeshipTransferred, error) {
	event := new(KeeperRegistryPayeeshipTransferred)
	if err := _KeeperRegistry.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryPaymentWithdrawnIterator struct {
	Event *KeeperRegistryPaymentWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryPaymentWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryPaymentWithdrawn)
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
		it.Event = new(KeeperRegistryPaymentWithdrawn)
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

func (it *KeeperRegistryPaymentWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryPaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryPaymentWithdrawn struct {
	Transmitter common.Address
	Amount      *big.Int
	To          common.Address
	Payee       common.Address
	Raw         types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistryPaymentWithdrawnIterator, error) {

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

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryPaymentWithdrawnIterator{contract: _KeeperRegistry.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryPaymentWithdrawn)
				if err := _KeeperRegistry.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParsePaymentWithdrawn(log types.Log) (*KeeperRegistryPaymentWithdrawn, error) {
	event := new(KeeperRegistryPaymentWithdrawn)
	if err := _KeeperRegistry.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryReorgedUpkeepReportIterator struct {
	Event *KeeperRegistryReorgedUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryReorgedUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryReorgedUpkeepReport)
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
		it.Event = new(KeeperRegistryReorgedUpkeepReport)
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

func (it *KeeperRegistryReorgedUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryReorgedUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryReorgedUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryReorgedUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryReorgedUpkeepReportIterator{contract: _KeeperRegistry.contract, event: "ReorgedUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryReorgedUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryReorgedUpkeepReport)
				if err := _KeeperRegistry.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseReorgedUpkeepReport(log types.Log) (*KeeperRegistryReorgedUpkeepReport, error) {
	event := new(KeeperRegistryReorgedUpkeepReport)
	if err := _KeeperRegistry.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryStaleUpkeepReportIterator struct {
	Event *KeeperRegistryStaleUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryStaleUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryStaleUpkeepReport)
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
		it.Event = new(KeeperRegistryStaleUpkeepReport)
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

func (it *KeeperRegistryStaleUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryStaleUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryStaleUpkeepReport struct {
	Id      *big.Int
	Trigger []byte
	Raw     types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryStaleUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryStaleUpkeepReportIterator{contract: _KeeperRegistry.contract, event: "StaleUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryStaleUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryStaleUpkeepReport)
				if err := _KeeperRegistry.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseStaleUpkeepReport(log types.Log) (*KeeperRegistryStaleUpkeepReport, error) {
	event := new(KeeperRegistryStaleUpkeepReport)
	if err := _KeeperRegistry.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryTransmittedIterator struct {
	Event *KeeperRegistryTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryTransmitted)
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
		it.Event = new(KeeperRegistryTransmitted)
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

func (it *KeeperRegistryTransmittedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterTransmitted(opts *bind.FilterOpts) (*KeeperRegistryTransmittedIterator, error) {

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryTransmittedIterator{contract: _KeeperRegistry.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *KeeperRegistryTransmitted) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryTransmitted)
				if err := _KeeperRegistry.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseTransmitted(log types.Log) (*KeeperRegistryTransmitted, error) {
	event := new(KeeperRegistryTransmitted)
	if err := _KeeperRegistry.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUnpausedIterator struct {
	Event *KeeperRegistryUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUnpaused)
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
		it.Event = new(KeeperRegistryUnpaused)
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

func (it *KeeperRegistryUnpausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistryUnpausedIterator, error) {

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUnpausedIterator{contract: _KeeperRegistry.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUnpaused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUnpaused)
				if err := _KeeperRegistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUnpaused(log types.Log) (*KeeperRegistryUnpaused, error) {
	event := new(KeeperRegistryUnpaused)
	if err := _KeeperRegistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUpkeepAdminTransferRequestedIterator struct {
	Event *KeeperRegistryUpkeepAdminTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUpkeepAdminTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepAdminTransferRequested)
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
		it.Event = new(KeeperRegistryUpkeepAdminTransferRequested)
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

func (it *KeeperRegistryUpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryUpkeepAdminTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepAdminTransferRequestedIterator{contract: _KeeperRegistry.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUpkeepAdminTransferRequested)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepAdminTransferRequested(log types.Log) (*KeeperRegistryUpkeepAdminTransferRequested, error) {
	event := new(KeeperRegistryUpkeepAdminTransferRequested)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUpkeepAdminTransferredIterator struct {
	Event *KeeperRegistryUpkeepAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUpkeepAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepAdminTransferred)
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
		it.Event = new(KeeperRegistryUpkeepAdminTransferred)
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

func (it *KeeperRegistryUpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryUpkeepAdminTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepAdminTransferredIterator{contract: _KeeperRegistry.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUpkeepAdminTransferred)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepAdminTransferred(log types.Log) (*KeeperRegistryUpkeepAdminTransferred, error) {
	event := new(KeeperRegistryUpkeepAdminTransferred)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUpkeepCanceledIterator struct {
	Event *KeeperRegistryUpkeepCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUpkeepCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepCanceled)
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
		it.Event = new(KeeperRegistryUpkeepCanceled)
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

func (it *KeeperRegistryUpkeepCanceledIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistryUpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepCanceledIterator{contract: _KeeperRegistry.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUpkeepCanceled)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepCanceled(log types.Log) (*KeeperRegistryUpkeepCanceled, error) {
	event := new(KeeperRegistryUpkeepCanceled)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUpkeepCheckDataSetIterator struct {
	Event *KeeperRegistryUpkeepCheckDataSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUpkeepCheckDataSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepCheckDataSet)
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
		it.Event = new(KeeperRegistryUpkeepCheckDataSet)
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

func (it *KeeperRegistryUpkeepCheckDataSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUpkeepCheckDataSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUpkeepCheckDataSet struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepCheckDataSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepCheckDataSetIterator{contract: _KeeperRegistry.contract, event: "UpkeepCheckDataSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepCheckDataSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepCheckDataSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUpkeepCheckDataSet)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepCheckDataSet(log types.Log) (*KeeperRegistryUpkeepCheckDataSet, error) {
	event := new(KeeperRegistryUpkeepCheckDataSet)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepCheckDataSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUpkeepGasLimitSetIterator struct {
	Event *KeeperRegistryUpkeepGasLimitSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUpkeepGasLimitSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepGasLimitSet)
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
		it.Event = new(KeeperRegistryUpkeepGasLimitSet)
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

func (it *KeeperRegistryUpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepGasLimitSetIterator{contract: _KeeperRegistry.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUpkeepGasLimitSet)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepGasLimitSet(log types.Log) (*KeeperRegistryUpkeepGasLimitSet, error) {
	event := new(KeeperRegistryUpkeepGasLimitSet)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUpkeepMigratedIterator struct {
	Event *KeeperRegistryUpkeepMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUpkeepMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepMigrated)
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
		it.Event = new(KeeperRegistryUpkeepMigrated)
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

func (it *KeeperRegistryUpkeepMigratedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepMigratedIterator{contract: _KeeperRegistry.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUpkeepMigrated)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepMigrated(log types.Log) (*KeeperRegistryUpkeepMigrated, error) {
	event := new(KeeperRegistryUpkeepMigrated)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUpkeepOffchainConfigSetIterator struct {
	Event *KeeperRegistryUpkeepOffchainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUpkeepOffchainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepOffchainConfigSet)
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
		it.Event = new(KeeperRegistryUpkeepOffchainConfigSet)
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

func (it *KeeperRegistryUpkeepOffchainConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUpkeepOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUpkeepOffchainConfigSet struct {
	Id             *big.Int
	OffchainConfig []byte
	Raw            types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepOffchainConfigSetIterator{contract: _KeeperRegistry.contract, event: "UpkeepOffchainConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUpkeepOffchainConfigSet)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepOffchainConfigSet(log types.Log) (*KeeperRegistryUpkeepOffchainConfigSet, error) {
	event := new(KeeperRegistryUpkeepOffchainConfigSet)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUpkeepPausedIterator struct {
	Event *KeeperRegistryUpkeepPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUpkeepPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepPaused)
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
		it.Event = new(KeeperRegistryUpkeepPaused)
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

func (it *KeeperRegistryUpkeepPausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUpkeepPaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepPausedIterator{contract: _KeeperRegistry.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUpkeepPaused)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepPaused(log types.Log) (*KeeperRegistryUpkeepPaused, error) {
	event := new(KeeperRegistryUpkeepPaused)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUpkeepPerformedIterator struct {
	Event *KeeperRegistryUpkeepPerformed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUpkeepPerformedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepPerformed)
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
		it.Event = new(KeeperRegistryUpkeepPerformed)
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

func (it *KeeperRegistryUpkeepPerformedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUpkeepPerformed struct {
	Id           *big.Int
	Success      bool
	TotalPayment *big.Int
	GasUsed      *big.Int
	GasOverhead  *big.Int
	Trigger      []byte
	Raw          types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*KeeperRegistryUpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepPerformedIterator{contract: _KeeperRegistry.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUpkeepPerformed)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepPerformed(log types.Log) (*KeeperRegistryUpkeepPerformed, error) {
	event := new(KeeperRegistryUpkeepPerformed)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUpkeepPrivilegeConfigSetIterator struct {
	Event *KeeperRegistryUpkeepPrivilegeConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUpkeepPrivilegeConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepPrivilegeConfigSet)
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
		it.Event = new(KeeperRegistryUpkeepPrivilegeConfigSet)
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

func (it *KeeperRegistryUpkeepPrivilegeConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUpkeepPrivilegeConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUpkeepPrivilegeConfigSet struct {
	Id              *big.Int
	PrivilegeConfig []byte
	Raw             types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepPrivilegeConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepPrivilegeConfigSetIterator{contract: _KeeperRegistry.contract, event: "UpkeepPrivilegeConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepPrivilegeConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUpkeepPrivilegeConfigSet)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepPrivilegeConfigSet(log types.Log) (*KeeperRegistryUpkeepPrivilegeConfigSet, error) {
	event := new(KeeperRegistryUpkeepPrivilegeConfigSet)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepPrivilegeConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUpkeepReceivedIterator struct {
	Event *KeeperRegistryUpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepReceived)
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
		it.Event = new(KeeperRegistryUpkeepReceived)
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

func (it *KeeperRegistryUpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepReceivedIterator{contract: _KeeperRegistry.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUpkeepReceived)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepReceived(log types.Log) (*KeeperRegistryUpkeepReceived, error) {
	event := new(KeeperRegistryUpkeepReceived)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUpkeepRegisteredIterator struct {
	Event *KeeperRegistryUpkeepRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUpkeepRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepRegistered)
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
		it.Event = new(KeeperRegistryUpkeepRegistered)
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

func (it *KeeperRegistryUpkeepRegisteredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUpkeepRegistered struct {
	Id         *big.Int
	PerformGas uint32
	Admin      common.Address
	Raw        types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepRegisteredIterator{contract: _KeeperRegistry.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUpkeepRegistered)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepRegistered(log types.Log) (*KeeperRegistryUpkeepRegistered, error) {
	event := new(KeeperRegistryUpkeepRegistered)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUpkeepTriggerConfigSetIterator struct {
	Event *KeeperRegistryUpkeepTriggerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUpkeepTriggerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepTriggerConfigSet)
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
		it.Event = new(KeeperRegistryUpkeepTriggerConfigSet)
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

func (it *KeeperRegistryUpkeepTriggerConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUpkeepTriggerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUpkeepTriggerConfigSet struct {
	Id            *big.Int
	TriggerConfig []byte
	Raw           types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepTriggerConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepTriggerConfigSetIterator{contract: _KeeperRegistry.contract, event: "UpkeepTriggerConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepTriggerConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUpkeepTriggerConfigSet)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepTriggerConfigSet(log types.Log) (*KeeperRegistryUpkeepTriggerConfigSet, error) {
	event := new(KeeperRegistryUpkeepTriggerConfigSet)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepTriggerConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryUpkeepUnpausedIterator struct {
	Event *KeeperRegistryUpkeepUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryUpkeepUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryUpkeepUnpaused)
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
		it.Event = new(KeeperRegistryUpkeepUnpaused)
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

func (it *KeeperRegistryUpkeepUnpausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryUpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryUpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepUnpausedIterator{contract: _KeeperRegistry.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryUpkeepUnpaused)
				if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistryFilterer) ParseUpkeepUnpaused(log types.Log) (*KeeperRegistryUpkeepUnpaused, error) {
	event := new(KeeperRegistryUpkeepUnpaused)
	if err := _KeeperRegistry.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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

func (_KeeperRegistry *KeeperRegistry) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _KeeperRegistry.abi.Events["AdminPrivilegeConfigSet"].ID:
		return _KeeperRegistry.ParseAdminPrivilegeConfigSet(log)
	case _KeeperRegistry.abi.Events["CancelledUpkeepReport"].ID:
		return _KeeperRegistry.ParseCancelledUpkeepReport(log)
	case _KeeperRegistry.abi.Events["ConfigSet"].ID:
		return _KeeperRegistry.ParseConfigSet(log)
	case _KeeperRegistry.abi.Events["DedupKeyAdded"].ID:
		return _KeeperRegistry.ParseDedupKeyAdded(log)
	case _KeeperRegistry.abi.Events["FundsAdded"].ID:
		return _KeeperRegistry.ParseFundsAdded(log)
	case _KeeperRegistry.abi.Events["FundsWithdrawn"].ID:
		return _KeeperRegistry.ParseFundsWithdrawn(log)
	case _KeeperRegistry.abi.Events["InsufficientFundsUpkeepReport"].ID:
		return _KeeperRegistry.ParseInsufficientFundsUpkeepReport(log)
	case _KeeperRegistry.abi.Events["OwnerFundsWithdrawn"].ID:
		return _KeeperRegistry.ParseOwnerFundsWithdrawn(log)
	case _KeeperRegistry.abi.Events["OwnershipTransferRequested"].ID:
		return _KeeperRegistry.ParseOwnershipTransferRequested(log)
	case _KeeperRegistry.abi.Events["OwnershipTransferred"].ID:
		return _KeeperRegistry.ParseOwnershipTransferred(log)
	case _KeeperRegistry.abi.Events["Paused"].ID:
		return _KeeperRegistry.ParsePaused(log)
	case _KeeperRegistry.abi.Events["PayeesUpdated"].ID:
		return _KeeperRegistry.ParsePayeesUpdated(log)
	case _KeeperRegistry.abi.Events["PayeeshipTransferRequested"].ID:
		return _KeeperRegistry.ParsePayeeshipTransferRequested(log)
	case _KeeperRegistry.abi.Events["PayeeshipTransferred"].ID:
		return _KeeperRegistry.ParsePayeeshipTransferred(log)
	case _KeeperRegistry.abi.Events["PaymentWithdrawn"].ID:
		return _KeeperRegistry.ParsePaymentWithdrawn(log)
	case _KeeperRegistry.abi.Events["ReorgedUpkeepReport"].ID:
		return _KeeperRegistry.ParseReorgedUpkeepReport(log)
	case _KeeperRegistry.abi.Events["StaleUpkeepReport"].ID:
		return _KeeperRegistry.ParseStaleUpkeepReport(log)
	case _KeeperRegistry.abi.Events["Transmitted"].ID:
		return _KeeperRegistry.ParseTransmitted(log)
	case _KeeperRegistry.abi.Events["Unpaused"].ID:
		return _KeeperRegistry.ParseUnpaused(log)
	case _KeeperRegistry.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _KeeperRegistry.ParseUpkeepAdminTransferRequested(log)
	case _KeeperRegistry.abi.Events["UpkeepAdminTransferred"].ID:
		return _KeeperRegistry.ParseUpkeepAdminTransferred(log)
	case _KeeperRegistry.abi.Events["UpkeepCanceled"].ID:
		return _KeeperRegistry.ParseUpkeepCanceled(log)
	case _KeeperRegistry.abi.Events["UpkeepCheckDataSet"].ID:
		return _KeeperRegistry.ParseUpkeepCheckDataSet(log)
	case _KeeperRegistry.abi.Events["UpkeepGasLimitSet"].ID:
		return _KeeperRegistry.ParseUpkeepGasLimitSet(log)
	case _KeeperRegistry.abi.Events["UpkeepMigrated"].ID:
		return _KeeperRegistry.ParseUpkeepMigrated(log)
	case _KeeperRegistry.abi.Events["UpkeepOffchainConfigSet"].ID:
		return _KeeperRegistry.ParseUpkeepOffchainConfigSet(log)
	case _KeeperRegistry.abi.Events["UpkeepPaused"].ID:
		return _KeeperRegistry.ParseUpkeepPaused(log)
	case _KeeperRegistry.abi.Events["UpkeepPerformed"].ID:
		return _KeeperRegistry.ParseUpkeepPerformed(log)
	case _KeeperRegistry.abi.Events["UpkeepPrivilegeConfigSet"].ID:
		return _KeeperRegistry.ParseUpkeepPrivilegeConfigSet(log)
	case _KeeperRegistry.abi.Events["UpkeepReceived"].ID:
		return _KeeperRegistry.ParseUpkeepReceived(log)
	case _KeeperRegistry.abi.Events["UpkeepRegistered"].ID:
		return _KeeperRegistry.ParseUpkeepRegistered(log)
	case _KeeperRegistry.abi.Events["UpkeepTriggerConfigSet"].ID:
		return _KeeperRegistry.ParseUpkeepTriggerConfigSet(log)
	case _KeeperRegistry.abi.Events["UpkeepUnpaused"].ID:
		return _KeeperRegistry.ParseUpkeepUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeeperRegistryAdminPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x7c44b4eb59ee7873514e7e43e7718c269d872965938b288aa143befca62f99d2")
}

func (KeeperRegistryCancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xc3237c8807c467c1b39b8d0395eff077313e691bf0a7388106792564ebfd5636")
}

func (KeeperRegistryConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (KeeperRegistryDedupKeyAdded) Topic() common.Hash {
	return common.HexToHash("0xa4a4e334c0e330143f9437484fe516c13bc560b86b5b0daf58e7084aaac228f2")
}

func (KeeperRegistryFundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (KeeperRegistryFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (KeeperRegistryInsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x377c8b0c126ae5248d27aca1c76fac4608aff85673ee3caf09747e1044549e02")
}

func (KeeperRegistryOwnerFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1")
}

func (KeeperRegistryOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (KeeperRegistryOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (KeeperRegistryPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (KeeperRegistryPayeesUpdated) Topic() common.Hash {
	return common.HexToHash("0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725")
}

func (KeeperRegistryPayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (KeeperRegistryPayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (KeeperRegistryPaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (KeeperRegistryReorgedUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x6aa7f60c176da7af894b384daea2249497448137f5943c1237ada8bc92bdc301")
}

func (KeeperRegistryStaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x405288ea7be309e16cfdf481367f90a413e1d4634fcdaf8966546db9b93012e8")
}

func (KeeperRegistryTransmitted) Topic() common.Hash {
	return common.HexToHash("0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62")
}

func (KeeperRegistryUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (KeeperRegistryUpkeepAdminTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35")
}

func (KeeperRegistryUpkeepAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c")
}

func (KeeperRegistryUpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (KeeperRegistryUpkeepCheckDataSet) Topic() common.Hash {
	return common.HexToHash("0xcba2d5723b2ee59e53a8e8a82a4a7caf4fdfe70e9f7c582950bf7e7a5c24e83d")
}

func (KeeperRegistryUpkeepGasLimitSet) Topic() common.Hash {
	return common.HexToHash("0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c")
}

func (KeeperRegistryUpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (KeeperRegistryUpkeepOffchainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850")
}

func (KeeperRegistryUpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (KeeperRegistryUpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xad8cc9579b21dfe2c2f6ea35ba15b656e46b4f5b0cb424f52739b8ce5cac9c5b")
}

func (KeeperRegistryUpkeepPrivilegeConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2fd8d70753a007014349d4591843cc031c2dd7a260d7dd82eca8253686ae7769")
}

func (KeeperRegistryUpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (KeeperRegistryUpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (KeeperRegistryUpkeepTriggerConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664")
}

func (KeeperRegistryUpkeepUnpaused) Topic() common.Hash {
	return common.HexToHash("0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456")
}

func (_KeeperRegistry *KeeperRegistry) Address() common.Address {
	return _KeeperRegistry.address
}

type KeeperRegistryInterface interface {
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

	SetConfigTypeSafe(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig KeeperRegistryBase21OnchainConfig, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	SimulatePerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error)

	FilterAdminPrivilegeConfigSet(opts *bind.FilterOpts, admin []common.Address) (*KeeperRegistryAdminPrivilegeConfigSetIterator, error)

	WatchAdminPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryAdminPrivilegeConfigSet, admin []common.Address) (event.Subscription, error)

	ParseAdminPrivilegeConfigSet(log types.Log) (*KeeperRegistryAdminPrivilegeConfigSet, error)

	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryCancelledUpkeepReportIterator, error)

	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryCancelledUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseCancelledUpkeepReport(log types.Log) (*KeeperRegistryCancelledUpkeepReport, error)

	FilterConfigSet(opts *bind.FilterOpts) (*KeeperRegistryConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*KeeperRegistryConfigSet, error)

	FilterDedupKeyAdded(opts *bind.FilterOpts, dedupKey [][32]byte) (*KeeperRegistryDedupKeyAddedIterator, error)

	WatchDedupKeyAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryDedupKeyAdded, dedupKey [][32]byte) (event.Subscription, error)

	ParseDedupKeyAdded(log types.Log) (*KeeperRegistryDedupKeyAdded, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistryFundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*KeeperRegistryFundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryFundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryFundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*KeeperRegistryFundsWithdrawn, error)

	FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryInsufficientFundsUpkeepReportIterator, error)

	WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseInsufficientFundsUpkeepReport(log types.Log) (*KeeperRegistryInsufficientFundsUpkeepReport, error)

	FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*KeeperRegistryOwnerFundsWithdrawnIterator, error)

	WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryOwnerFundsWithdrawn) (event.Subscription, error)

	ParseOwnerFundsWithdrawn(log types.Log) (*KeeperRegistryOwnerFundsWithdrawn, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*KeeperRegistryOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*KeeperRegistryPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*KeeperRegistryPaused, error)

	FilterPayeesUpdated(opts *bind.FilterOpts) (*KeeperRegistryPayeesUpdatedIterator, error)

	WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPayeesUpdated) (event.Subscription, error)

	ParsePayeesUpdated(log types.Log) (*KeeperRegistryPayeesUpdated, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryPayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistryPayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryPayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*KeeperRegistryPayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistryPaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*KeeperRegistryPaymentWithdrawn, error)

	FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryReorgedUpkeepReportIterator, error)

	WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryReorgedUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseReorgedUpkeepReport(log types.Log) (*KeeperRegistryReorgedUpkeepReport, error)

	FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryStaleUpkeepReportIterator, error)

	WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryStaleUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseStaleUpkeepReport(log types.Log) (*KeeperRegistryStaleUpkeepReport, error)

	FilterTransmitted(opts *bind.FilterOpts) (*KeeperRegistryTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *KeeperRegistryTransmitted) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*KeeperRegistryTransmitted, error)

	FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistryUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*KeeperRegistryUnpaused, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryUpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*KeeperRegistryUpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryUpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*KeeperRegistryUpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistryUpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*KeeperRegistryUpkeepCanceled, error)

	FilterUpkeepCheckDataSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepCheckDataSetIterator, error)

	WatchUpkeepCheckDataSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepCheckDataSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCheckDataSet(log types.Log) (*KeeperRegistryUpkeepCheckDataSet, error)

	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepGasLimitSetIterator, error)

	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepGasLimitSet(log types.Log) (*KeeperRegistryUpkeepGasLimitSet, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*KeeperRegistryUpkeepMigrated, error)

	FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepOffchainConfigSetIterator, error)

	WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepOffchainConfigSet(log types.Log) (*KeeperRegistryUpkeepOffchainConfigSet, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*KeeperRegistryUpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*KeeperRegistryUpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*KeeperRegistryUpkeepPerformed, error)

	FilterUpkeepPrivilegeConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepPrivilegeConfigSetIterator, error)

	WatchUpkeepPrivilegeConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepPrivilegeConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPrivilegeConfigSet(log types.Log) (*KeeperRegistryUpkeepPrivilegeConfigSet, error)

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*KeeperRegistryUpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*KeeperRegistryUpkeepRegistered, error)

	FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepTriggerConfigSetIterator, error)

	WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepTriggerConfigSet(log types.Log) (*KeeperRegistryUpkeepTriggerConfigSet, error)

	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepUnpausedIterator, error)

	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepUnpaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepUnpaused(log types.Log) (*KeeperRegistryUpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
