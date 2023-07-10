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

var KeeperRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractKeeperRegistryLogicA2_1\",\"name\":\"logicA\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTrigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTriggerType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByUpkeepPrivilegeManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"triggerID\",\"type\":\"bytes32\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"triggerID\",\"type\":\"bytes32\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"triggerID\",\"type\":\"bytes32\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"triggerID\",\"type\":\"bytes32\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"triggerID\",\"type\":\"bytes32\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"privilegeConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepPrivilegeConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepTriggerConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fallbackTo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFastGasFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkNativeFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMode\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase2_1.Mode\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTriggerType\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase2_1.Trigger\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfigBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"simulatePerformUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"rawReport\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6101206040523480156200001257600080fd5b506040516200501f3803806200501f833981016040819052620000359162000374565b80816001600160a01b0316634b4fd03b6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200009b91906200039b565b826001600160a01b031663ca30e6036040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000da573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000100919062000374565b836001600160a01b031663b10b673c6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200013f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000165919062000374565b846001600160a01b0316636709d0e56040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca919062000374565b3380600081620002215760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b038481169190911790915581161562000254576200025481620002b0565b5050508360028111156200026c576200026c620003be565b60e0816002811115620002835762000283620003be565b9052506001600160a01b0392831660805290821660a052811660c052919091166101005250620003d49050565b336001600160a01b038216036200030a5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000218565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200037157600080fd5b50565b6000602082840312156200038757600080fd5b815162000394816200035b565b9392505050565b600060208284031215620003ae57600080fd5b8151600381106200039457600080fd5b634e487b7160e01b600052602160045260246000fd5b60805160a05160c05160e05161010051614bdc6200044360003960008181610102015261019b0152600081816101e201528181612ea40152818161320d015281816133a001526139c901526000610230015260006103770152600081816103b001526105bf0152614bdc6000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c8063a4c0ed3611610097578063b1dc65a411610066578063b1dc65a41461039b578063ca30e603146103ae578063e3d0e712146103d4578063f2fde38b146103e757610100565b8063a4c0ed36146102ef578063aed2e92914610302578063afcb95d71461032c578063b10b673c1461037557610100565b80636709d0e5116100d35780636709d0e51461022e57806379ba50971461025457806381ff70481461025c5780638da5cb5b146102d157610100565b8063181f5a7714610147578063349e8cca146101995780634b4fd03b146101e05780635147cd591461020e575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e808015610140573d6000f35b3d6000fd5b005b6101836040518060400160405280601481526020017f4b6565706572526567697374727920322e312e3000000000000000000000000081525081565b6040516101909190613c48565b60405180910390f35b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610190565b7f00000000000000000000000000000000000000000000000000000000000000006040516101909190613c8a565b61022161021c366004613ca4565b6103fa565b6040516101909190613cbd565b7f00000000000000000000000000000000000000000000000000000000000000006101bb565b6101456104a5565b6102ae60145460115463ffffffff780100000000000000000000000000000000000000000000000083048116937c01000000000000000000000000000000000000000000000000000000009093041691565b6040805163ffffffff948516815293909216602084015290820152606001610190565b60005473ffffffffffffffffffffffffffffffffffffffff166101bb565b6101456102fd366004613d3c565b6105a7565b610315610310366004613d98565b6107c3565b604080519215158352602083019190915201610190565b601154601254604080516000815260208101939093527c010000000000000000000000000000000000000000000000000000000090910463ffffffff1690820152606001610190565b7f00000000000000000000000000000000000000000000000000000000000000006101bb565b6101456103a9366004613e29565b610939565b7f00000000000000000000000000000000000000000000000000000000000000006101bb565b6101456103e23660046140fa565b611559565b6101456103f53660046141c7565b612338565b6000818160045b600f811015610487577fff00000000000000000000000000000000000000000000000000000000000000821683826020811061043f5761043f6141e4565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161461047557506000949350505050565b8061047f81614242565b915050610401565b5081600f1a600181111561049d5761049d613c5b565b949350505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461052b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610616576040517fc8bad78d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208114610650576040517fdfe9309000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061065e82840184613ca4565b60008181526004602052604090205490915065010000000000900463ffffffff908116146106b8576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000818152600460205260409020600101546106f39085906c0100000000000000000000000090046bffffffffffffffffffffffff1661427a565b600082815260046020526040902060010180546bffffffffffffffffffffffff929092166c01000000000000000000000000027fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff90921691909117905560185461075e90859061429f565b6018556040516bffffffffffffffffffffffff8516815273ffffffffffffffffffffffffffffffffffffffff86169082907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a35050505050565b6000806107ce61234c565b6012546e010000000000000000000000000000900460ff161561081d576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008581526004602090815260409182902082516101008082018552825460ff81161515835263ffffffff918104821683860181905265010000000000820483168488015273ffffffffffffffffffffffffffffffffffffffff690100000000000000000090920482166060850181905260018601546bffffffffffffffffffffffff80821660808801526c0100000000000000000000000082041660a08701527801000000000000000000000000000000000000000000000000900490931660c08501526002909401541660e08301528451601f8901859004850281018501909552878552909361092c93919291899089908190840183828082843760009201919091525061238692505050565b9097909650945050505050565b60005a604080516101208101825260125460ff808216835261010080830463ffffffff90811660208601526501000000000084048116958501959095526901000000000000000000830462ffffff1660608501526c01000000000000000000000000830461ffff1660808501526e0100000000000000000000000000008304821615801560a08601526f010000000000000000000000000000008404909216151560c085015270010000000000000000000000000000000083046bffffffffffffffffffffffff1660e08501527c010000000000000000000000000000000000000000000000000000000090920490931690820152919250610a67576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600b602052604090205460ff16610ab0576040517f1099ed7500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6011548a3514610aec576040517fdfdcf8e700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8051610af99060016142b2565b60ff1686141580610b0a5750858414155b15610b41576040517f0244f71a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610b518a8a8a8a8a8a8a8a612588565b6000610b5d8a8a6127f1565b9050600081604001515167ffffffffffffffff811115610b7f57610b7f613ee0565b604051908082528060200260200182016040528015610c4b57816020015b60408051610200810182526000610100820181815261012083018290526101408301829052610160830182905261018083018290526101a083018290526101c083018290526101e0830182905282526020808301829052928201819052606082018190526080820181905260a0820181905260c0820181905260e082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181610b9d5790505b5090506000805b83604001515181101561114d576004600085604001518381518110610c7957610c796141e4565b602090810291909101810151825281810192909252604090810160002081516101008082018452825460ff81161515835263ffffffff91810482169583019590955265010000000000850481169382019390935273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009094048416606082015260018201546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490921660c08301526002015490911660e08201528351849083908110610d6b57610d6b6141e4565b602002602001015160000181905250610da084604001518281518110610d9357610d936141e4565b60200260200101516103fa565b838281518110610db257610db26141e4565b6020026020010151608001906001811115610dcf57610dcf613c5b565b90816001811115610de257610de2613c5b565b81525050610e5a85848381518110610dfc57610dfc6141e4565b602002602001015160800151858481518110610e1a57610e1a6141e4565b602002602001015160000151602001518760a001518581518110610e4057610e406141e4565b6020026020010151518860000151896020015160016128ac565b838281518110610e6c57610e6c6141e4565b6020026020010151604001906bffffffffffffffffffffffff1690816bffffffffffffffffffffffff1681525050610ede84604001518281518110610eb357610eb36141e4565b602002602001015185608001518381518110610ed157610ed16141e4565b60200260200101516128f7565b838281518110610ef057610ef06141e4565b602002602001015160e0018181525050610fbc84604001518281518110610f1957610f196141e4565b6020026020010151848381518110610f3357610f336141e4565b602002602001015160800151858481518110610f5157610f516141e4565b602002602001015160e0015187608001518581518110610f7357610f736141e4565b6020026020010151878681518110610f8d57610f8d6141e4565b602002602001015160000151888781518110610fab57610fab6141e4565b60200260200101516040015161292a565b838281518110610fce57610fce6141e4565b60200260200101516020019015159081151581525050828181518110610ff657610ff66141e4565b60200260200101516020015115611019576110126001836142cb565b915061101e565b61113b565b611084838281518110611033576110336141e4565b6020026020010151600001516060015185606001518381518110611059576110596141e4565b60200260200101518660a001518481518110611077576110776141e4565b6020026020010151612386565b848381518110611096576110966141e4565b60200260200101516060018584815181106110b3576110b36141e4565b602002602001015160a00182815250821515151581525050508281815181106110de576110de6141e4565b602002602001015160a00151866110f591906142e6565b955061113b84604001518281518110611110576111106141e4565b602002602001015184838151811061112a5761112a6141e4565b602002602001015160800151612a79565b8061114581614242565b915050610c52565b508061ffff1660000361116457505050505061154f565b83516111719060016142b2565b6111809060ff1661044c6142f9565b616b6c61118e8d60106142f9565b5a61119990896142e6565b6111a3919061429f565b6111ad919061429f565b6111b7919061429f565b9450611b586111ca61ffff831687614365565b6111d4919061429f565b945060008060008060005b8760400151518110156113f1578681815181106111fe576111fe6141e4565b602002602001015160200151156113df5761125a8a888381518110611225576112256141e4565b6020026020010151608001518a60a001518481518110611247576112476141e4565b6020026020010151518c60000151612b00565b87828151811061126c5761126c6141e4565b602002602001015160c00181815250506112c88989604001518381518110611296576112966141e4565b60200260200101518984815181106112b0576112b06141e4565b60200260200101518b600001518c602001518b612b20565b90935091506112d7828561427a565b93506112e3838661427a565b94508681815181106112f7576112f76141e4565b60200260200101516060015115158860400151828151811061131b5761131b6141e4565b60200260200101517fdacc0b3d36199a5f78afb27c76cd01bf100f605edf5ee0c0e23055b0553b68c68486611350919061427a565b8a8581518110611362576113626141e4565b602002602001015160a001518b8681518110611380576113806141e4565b602002602001015160c001518c878151811061139e5761139e6141e4565b60209081029190910181015160e00151604080516bffffffffffffffffffffffff909616865291850193909352830152606082015260800160405180910390a35b806113e981614242565b9150506111df565b5050336000908152600b6020526040902080548492506002906114299084906201000090046bffffffffffffffffffffffff1661427a565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555080601260000160108282829054906101000a90046bffffffffffffffffffffffff16611483919061427a565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060008f6001600381106114c6576114c66141e4565b602002013560001c9050600060088264ffffffffff16901c905087610100015163ffffffff168163ffffffff16111561154557601280547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167c010000000000000000000000000000000000000000000000000000000063ffffffff8416021790555b5050505050505050505b5050505050505050565b611561612c13565b601f8651111561159d576040517f25d0209c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8360ff166000036115da576040517fe77dba5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845186511415806115f957506115f1846003614379565b60ff16865111155b15611630576040517f1d2d1c5800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601254600e547001000000000000000000000000000000009091046bffffffffffffffffffffffff169060005b816bffffffffffffffffffffffff168110156116c5576116b2600e8281548110611689576116896141e4565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff168484612c94565b50806116bd81614242565b91505061165d565b5060008060005b836bffffffffffffffffffffffff168110156117ce57600d81815481106116f5576116f56141e4565b600091825260209091200154600e805473ffffffffffffffffffffffffffffffffffffffff90921694509082908110611730576117306141e4565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff8681168452600c8352604080852080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001690559116808452600b90925290912080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690559150806117c681614242565b9150506116cc565b506117db600d6000613b1d565b6117e7600e6000613b1d565b604080516080810182526000808252602082018190529181018290526060810182905290805b8c51811015611b6b57600c60008e838151811061182c5761182c6141e4565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528101919091526040016000205460ff1615611897576040517f77cea0fa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405180604001604052806001151581526020018260ff16815250600c60008f84815181106118c8576118c86141e4565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528181019290925260400160002082518154939092015160ff16610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909316929092171790558b518c9082908110611970576119706141e4565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff81166000908152600b83526040908190208151608081018352905460ff80821615801584526101008304909116958301959095526bffffffffffffffffffffffff6201000082048116938301939093526e0100000000000000000000000000009004909116606082015294509250611a35576040517f6a7281ad00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001835260ff80821660208086019182526bffffffffffffffffffffffff808b166060880190815273ffffffffffffffffffffffffffffffffffffffff87166000908152600b909352604092839020885181549551948a0151925184166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff939094166201000002929092167fffffffffffff000000000000000000000000000000000000000000000000ffff94909616610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009095169490941717919091169290921791909117905580611b6381614242565b91505061180d565b50508a51611b819150600d9060208d0190613b3b565b508851611b9590600e9060208c0190613b3b565b50600087806020019051810190611bac9190614466565b90506040518061012001604052808a60ff168152602001826000015163ffffffff168152602001826020015163ffffffff168152602001826060015162ffffff168152602001826080015161ffff1681526020016012600001600e9054906101000a900460ff16151581526020016012600001600f9054906101000a900460ff1615158152602001866bffffffffffffffffffffffff168152602001600063ffffffff16815250601260008201518160000160006101000a81548160ff021916908360ff16021790555060208201518160000160016101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160056101000a81548163ffffffff021916908363ffffffff16021790555060608201518160000160096101000a81548162ffffff021916908362ffffff160217905550608082015181600001600c6101000a81548161ffff021916908361ffff16021790555060a082015181600001600e6101000a81548160ff02191690831515021790555060c082015181600001600f6101000a81548160ff02191690831515021790555060e08201518160000160106101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555061010082015181600001601c6101000a81548163ffffffff021916908363ffffffff1602179055509050506040518061018001604052808260a001516bffffffffffffffffffffffff16815260200182610180015173ffffffffffffffffffffffffffffffffffffffff168152602001601360010160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff168152602001826040015163ffffffff1681526020018260c0015163ffffffff168152602001601360010160149054906101000a900463ffffffff1663ffffffff168152602001601360010160189054906101000a900463ffffffff1663ffffffff1681526020016013600101601c9054906101000a900463ffffffff1663ffffffff1681526020018260e0015163ffffffff16815260200182610100015163ffffffff16815260200182610120015163ffffffff168152602001826101c0015173ffffffffffffffffffffffffffffffffffffffff16815250601360008201518160000160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550602082015181600001600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060408201518160010160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550606082015181600101600c6101000a81548163ffffffff021916908363ffffffff16021790555060808201518160010160106101000a81548163ffffffff021916908363ffffffff16021790555060a08201518160010160146101000a81548163ffffffff021916908363ffffffff16021790555060c08201518160010160186101000a81548163ffffffff021916908363ffffffff16021790555060e082015181600101601c6101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160020160006101000a81548163ffffffff021916908363ffffffff1602179055506101208201518160020160046101000a81548163ffffffff021916908363ffffffff1602179055506101408201518160020160086101000a81548163ffffffff021916908363ffffffff16021790555061016082015181600201600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505080610140015160168190555080610160015160178190555060006013600101601c9054906101000a900463ffffffff16905061216f612e9e565b601480547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167c010000000000000000000000000000000000000000000000000000000063ffffffff93841602178082556001926018916121ea91859178010000000000000000000000000000000000000000000000009004166145c0565b92506101000a81548163ffffffff021916908363ffffffff1602179055506122344630601360010160189054906101000a900463ffffffff1663ffffffff168f8f8f8f8f8f612f53565b60115560005b6122446009612ffd565b81101561227457612261612259600983613007565b600990613013565b508061226c81614242565b91505061223a565b5060005b826101a00151518110156122cb576122b8836101a0015182815181106122a0576122a06141e4565b6020026020010151600961303590919063ffffffff16565b50806122c381614242565b915050612278565b507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0581601154601360010160189054906101000a900463ffffffff168f8f8f8f8f8f6040516123229998979695949392919061462e565b60405180910390a1505050505050505050505050565b612340612c13565b61234981613057565b50565b3215612384576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60125460009081906f01000000000000000000000000000000900460ff16156123db576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b601280547fffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffff166f010000000000000000000000000000001790555a9050634585e33b60e01b836040516024016124319190613c48565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009094169390931790925290517f79188d1600000000000000000000000000000000000000000000000000000000815290935073ffffffffffffffffffffffffffffffffffffffff8616906379188d169061250490879087906004016146c4565b6020604051808303816000875af1158015612523573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061254791906146dd565b91505a61255490826142e6565b9050601280547fffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffff1690559094909350915050565b6000878760405161259a9291906146ff565b6040519081900381206125b1918b9060200161470f565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201208383019092526000808452908301819052909250906000805b888110156127885760018587836020811061261d5761261d6141e4565b61262a91901a601b6142b2565b8c8c8581811061263c5761263c6141e4565b905060200201358b8b86818110612655576126556141e4565b9050602002013560405160008152602001604052604051612692949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa1580156126b4573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff81166000908152600c602090815290849020838501909452925460ff8082161515808552610100909204169383019390935290955093509050612762576040517f0f4c073700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826020015160080260ff166001901b84019350808061278090614242565b915050612600565b50827e010101010101010101010101010101010101010101010101010101010101018416146127e3576040517fc103be2e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505050505050505050505050565b61282a6040518060c001604052806000815260200160008152602001606081526020016060815260200160608152602001606081525090565b600061283883850185614800565b604081015151606082015151919250908114158061285b57508082608001515114155b8061286b5750808260a001515114155b156128a2576040517fb55ac75400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5090505b92915050565b6000806128be88878b6000015161314c565b90506000806128d98b8a63ffffffff16858a8a60018b6131d8565b90925090506128e8818361427a565b9b9a5050505050505050505050565b6000828260405160200161290c9291906148ed565b60405160208183030381529060405280519060200120905092915050565b60008086600181111561293f5761293f613c5b565b036129615761295087868686613583565b61295c57506000612a6f565b6129b7565b600186600181111561297557612975613c5b565b036129855761295087868661367f565b6040517ff2b2d41200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6129bf612e9e565b836040015163ffffffff1611612a1157867f846af3fdde418e2757c00c4cf8c6827e65ac4179753128c454922ccd7963886d86604051612a0191815260200190565b60405180910390a2506000612a6f565b816bffffffffffffffffffffffff168360a001516bffffffffffffffffffffffff161015612a6b57867ffcb5fa059ae0ddeeac53240ccaa0a8ad80993b33e18e8821f1e1cbc53481240e86604051612a0191815260200190565b5060015b9695505050505050565b6000816001811115612a8d57612a8d613c5b565b03612afc57612a9a612e9e565b6000838152600460205260409020600101805463ffffffff929092167801000000000000000000000000000000000000000000000000027fffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffff9092169190911790555b5050565b6000612b0d84848461314c565b90508085101561049d5750929392505050565b600080612b3b888760a001518860c0015188888860016131d8565b90925090506000612b4c828461427a565b600089815260046020526040902060010180549192508291600c90612b909084906c0100000000000000000000000090046bffffffffffffffffffffffff16614913565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915560008a815260046020526040812060010180548594509092612bd99185911661427a565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555050965096945050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314612384576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610522565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600b602090815260408083208151608081018352905460ff80821615801584526101008304909116948301949094526bffffffffffffffffffffffff6201000082048116938301939093526e0100000000000000000000000000009004909116606082015290612e90576000816060015185612d2c9190614913565b90506000612d3a8583614938565b90508083604001818151612d4e919061427a565b6bffffffffffffffffffffffff16905250612d698582614963565b83606001818151612d7a919061427a565b6bffffffffffffffffffffffff90811690915273ffffffffffffffffffffffffffffffffffffffff89166000908152600b602090815260409182902087518154928901519389015160608a015186166e010000000000000000000000000000027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff919096166201000002167fffffffffffff000000000000000000000000000000000000000000000000ffff60ff95909516610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909416939093171792909216179190911790555050505b6040015190505b9392505050565b600060017f00000000000000000000000000000000000000000000000000000000000000006002811115612ed457612ed4613c5b565b03612f4e57606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015612f25573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612f499190614997565b905090565b504390565b6000808a8a8a8a8a8a8a8a8a604051602001612f77999897969594939291906149b0565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179b9a5050505050505050505050565b60006128a6825490565b6000612e978383613850565b6000612e978373ffffffffffffffffffffffffffffffffffffffff841661387a565b6000612e978373ffffffffffffffffffffffffffffffffffffffff8416613974565b3373ffffffffffffffffffffffffffffffffffffffff8216036130d6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610522565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000808085600181111561316257613162613c5b565b03613171575062013880613190565b600185600181111561318557613185613c5b565b036129855750620186a05b6131a163ffffffff851660146142f9565b6131ac8460016142b2565b6131bb9060ff16611d4c6142f9565b6131c5908361429f565b6131cf919061429f565b95945050505050565b6000806000896080015161ffff16876131f191906142f9565b90508380156131ff5750803a105b1561320757503a5b600060027f0000000000000000000000000000000000000000000000000000000000000000600281111561323d5761323d613c5b565b0361339c57604080516000815260208101909152851561329b57600036604051806080016040528060488152602001614b886048913960405160200161328593929190614a45565b6040516020818303038152906040529050613303565b6015546132b790640100000000900463ffffffff166004614a6c565b63ffffffff1667ffffffffffffffff8111156132d5576132d5613ee0565b6040519080825280601f01601f1916602001820160405280156132ff576020820181803683370190505b5090505b6040517f49948e0e00000000000000000000000000000000000000000000000000000000815273420000000000000000000000000000000000000f906349948e0e90613353908490600401613c48565b602060405180830381865afa158015613370573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906133949190614997565b915050613448565b60017f000000000000000000000000000000000000000000000000000000000000000060028111156133d0576133d0613c5b565b0361344857606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613421573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906134459190614997565b90505b8461346457808b6080015161ffff1661346191906142f9565b90505b61347261ffff871682614365565b9050600087826134828c8e61429f565b61348c90866142f9565b613496919061429f565b6134a890670de0b6b3a76400006142f9565b6134b29190614365565b905060008c6040015163ffffffff1664e8d4a510006134d191906142f9565b898e6020015163ffffffff16858f886134ea91906142f9565b6134f4919061429f565b61350290633b9aca006142f9565b61350c91906142f9565b6135169190614365565b613520919061429f565b90506b033b2e3c9fd0803ce8000000613539828461429f565b1115613571576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b6000808380602001905181019061359a9190614a8f565b90508260c0015163ffffffff16816000015163ffffffff1610156135fc57857f909d8c7e883c6d2d53a5fb2dbb9ec57875a53c3769c099230027312ffb2a6b11866040516135ea91815260200190565b60405180910390a2600091505061049d565b6020810151158015906136235750602081015181516136209063ffffffff166139c3565b14155b8061363c5750613631612e9e565b815163ffffffff1610155b1561367357857f07c4f19521a82d67071422fca888431cac8267935db8abe73e354791a3d38868866040516135ea91815260200190565b50600195945050505050565b600080828060200190518101906136969190614ae6565b6060810151909150158015906136c2575080606001516136bf826040015163ffffffff166139c3565b14155b806136de57506136d0612e9e565b816040015163ffffffff1610155b1561372757847f07c4f19521a82d67071422fca888431cac8267935db8abe73e354791a3d388688560405161371591815260200190565b60405180910390a26000915050612e97565b8051602080830151604051600093613779938a93919201928352602083019190915260e01b7fffffffff0000000000000000000000000000000000000000000000000000000016604082015260440190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291815281516020928301206000818152600890935291205490915060ff161561380c57857f909d8c7e883c6d2d53a5fb2dbb9ec57875a53c3769c099230027312ffb2a6b11866040516137f991815260200190565b60405180910390a2600092505050612e97565b600090815260086020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660019081179091559150509392505050565b6000826000018281548110613867576138676141e4565b9060005260206000200154905092915050565b6000818152600183016020526040812054801561396357600061389e6001836142e6565b85549091506000906138b2906001906142e6565b90508181146139175760008660000182815481106138d2576138d26141e4565b90600052602060002001549050808760000184815481106138f5576138f56141e4565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061392857613928614b58565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506128a6565b60009150506128a6565b5092915050565b60008181526001830160205260408120546139bb575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556128a6565b5060006128a6565b600060017f000000000000000000000000000000000000000000000000000000000000000060028111156139f9576139f9613c5b565b03613b13576000606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613a4c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613a709190614997565b90508083101580613a8b5750610100613a8984836142e6565b115b15613a995750600092915050565b6040517f2b407a8200000000000000000000000000000000000000000000000000000000815260048101849052606490632b407a8290602401602060405180830381865afa158015613aef573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612e979190614997565b504090565b919050565b50805460008255906000526020600020908101906123499190613bc5565b828054828255906000526020600020908101928215613bb5579160200282015b82811115613bb557825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190613b5b565b50613bc1929150613bc5565b5090565b5b80821115613bc15760008155600101613bc6565b60005b83811015613bf5578181015183820152602001613bdd565b50506000910152565b60008151808452613c16816020860160208601613bda565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000612e976020830184613bfe565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b6020810160038310613c9e57613c9e613c5b565b91905290565b600060208284031215613cb657600080fd5b5035919050565b6020810160028310613c9e57613c9e613c5b565b73ffffffffffffffffffffffffffffffffffffffff8116811461234957600080fd5b60008083601f840112613d0557600080fd5b50813567ffffffffffffffff811115613d1d57600080fd5b602083019150836020828501011115613d3557600080fd5b9250929050565b60008060008060608587031215613d5257600080fd5b8435613d5d81613cd1565b935060208501359250604085013567ffffffffffffffff811115613d8057600080fd5b613d8c87828801613cf3565b95989497509550505050565b600080600060408486031215613dad57600080fd5b83359250602084013567ffffffffffffffff811115613dcb57600080fd5b613dd786828701613cf3565b9497909650939450505050565b60008083601f840112613df657600080fd5b50813567ffffffffffffffff811115613e0e57600080fd5b6020830191508360208260051b8501011115613d3557600080fd5b60008060008060008060008060e0898b031215613e4557600080fd5b606089018a811115613e5657600080fd5b8998503567ffffffffffffffff80821115613e7057600080fd5b613e7c8c838d01613cf3565b909950975060808b0135915080821115613e9557600080fd5b613ea18c838d01613de4565b909750955060a08b0135915080821115613eba57600080fd5b50613ec78b828c01613de4565b999c989b50969995989497949560c00135949350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516101e0810167ffffffffffffffff81118282101715613f3357613f33613ee0565b60405290565b60405160c0810167ffffffffffffffff81118282101715613f3357613f33613ee0565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613fa357613fa3613ee0565b604052919050565b600067ffffffffffffffff821115613fc557613fc5613ee0565b5060051b60200190565b600082601f830112613fe057600080fd5b81356020613ff5613ff083613fab565b613f5c565b82815260059290921b8401810191818101908684111561401457600080fd5b8286015b8481101561403857803561402b81613cd1565b8352918301918301614018565b509695505050505050565b803560ff81168114613b1857600080fd5b600082601f83011261406557600080fd5b813567ffffffffffffffff81111561407f5761407f613ee0565b6140b060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601613f5c565b8181528460208386010111156140c557600080fd5b816020850160208301376000918101602001919091529392505050565b803567ffffffffffffffff81168114613b1857600080fd5b60008060008060008060c0878903121561411357600080fd5b863567ffffffffffffffff8082111561412b57600080fd5b6141378a838b01613fcf565b9750602089013591508082111561414d57600080fd5b6141598a838b01613fcf565b965061416760408a01614043565b9550606089013591508082111561417d57600080fd5b6141898a838b01614054565b945061419760808a016140e2565b935060a08901359150808211156141ad57600080fd5b506141ba89828a01614054565b9150509295509295509295565b6000602082840312156141d957600080fd5b8135612e9781613cd1565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361427357614273614213565b5060010190565b6bffffffffffffffffffffffff81811683821601908082111561396d5761396d614213565b808201808211156128a6576128a6614213565b60ff81811683821601908111156128a6576128a6614213565b61ffff81811683821601908082111561396d5761396d614213565b818103818111156128a6576128a6614213565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561433157614331614213565b500290565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60008261437457614374614336565b500490565b600060ff821660ff84168160ff048111821515161561439a5761439a614213565b029392505050565b805163ffffffff81168114613b1857600080fd5b805162ffffff81168114613b1857600080fd5b805161ffff81168114613b1857600080fd5b80516bffffffffffffffffffffffff81168114613b1857600080fd5b8051613b1881613cd1565b600082601f83011261441357600080fd5b81516020614423613ff083613fab565b82815260059290921b8401810191818101908684111561444257600080fd5b8286015b8481101561403857805161445981613cd1565b8352918301918301614446565b60006020828403121561447857600080fd5b815167ffffffffffffffff8082111561449057600080fd5b908301906101e082860312156144a557600080fd5b6144ad613f0f565b6144b6836143a2565b81526144c4602084016143a2565b60208201526144d5604084016143a2565b60408201526144e6606084016143b6565b60608201526144f7608084016143c9565b608082015261450860a084016143db565b60a082015261451960c084016143a2565b60c082015261452a60e084016143a2565b60e082015261010061453d8185016143a2565b9082015261012061454f8482016143a2565b90820152610140838101519082015261016080840151908201526101806145778185016143f7565b908201526101a0838101518381111561458f57600080fd5b61459b88828701614402565b8284015250506101c091506145b18284016143f7565b91810191909152949350505050565b63ffffffff81811683821601908082111561396d5761396d614213565b600081518084526020808501945080840160005b8381101561462357815173ffffffffffffffffffffffffffffffffffffffff16875295820195908201906001016145f1565b509495945050505050565b600061012063ffffffff808d1684528b6020850152808b1660408501525080606084015261465e8184018a6145dd565b9050828103608084015261467281896145dd565b905060ff871660a084015282810360c084015261468f8187613bfe565b905067ffffffffffffffff851660e08401528281036101008401526146b48185613bfe565b9c9b505050505050505050505050565b82815260406020820152600061049d6040830184613bfe565b6000602082840312156146ef57600080fd5b81518015158114612e9757600080fd5b8183823760009101908152919050565b8281526080810160608360208401379392505050565b600082601f83011261473657600080fd5b81356020614746613ff083613fab565b82815260059290921b8401810191818101908684111561476557600080fd5b8286015b848110156140385780358352918301918301614769565b600082601f83011261479157600080fd5b813560206147a1613ff083613fab565b82815260059290921b840181019181810190868411156147c057600080fd5b8286015b8481101561403857803567ffffffffffffffff8111156147e45760008081fd5b6147f28986838b0101614054565b8452509183019183016147c4565b60006020828403121561481257600080fd5b813567ffffffffffffffff8082111561482a57600080fd5b9083019060c0828603121561483e57600080fd5b614846613f39565b823581526020830135602082015260408301358281111561486657600080fd5b61487287828601614725565b60408301525060608301358281111561488a57600080fd5b61489687828601614725565b6060830152506080830135828111156148ae57600080fd5b6148ba87828601614780565b60808301525060a0830135828111156148d257600080fd5b6148de87828601614780565b60a08301525095945050505050565b82815260008251614905816020850160208701613bda565b919091016020019392505050565b6bffffffffffffffffffffffff82811682821603908082111561396d5761396d614213565b60006bffffffffffffffffffffffff8084168061495757614957614336565b92169190910492915050565b60006bffffffffffffffffffffffff8083168185168183048111821515161561498e5761498e614213565b02949350505050565b6000602082840312156149a957600080fd5b5051919050565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b1660408501528160608501526149f78285018b6145dd565b91508382036080850152614a0b828a6145dd565b915060ff881660a085015283820360c0850152614a288288613bfe565b90861660e085015283810361010085015290506146b48185613bfe565b828482376000838201600081528351614a62818360208801613bda565b0195945050505050565b600063ffffffff8083168185168183048111821515161561498e5761498e614213565b600060408284031215614aa157600080fd5b6040516040810181811067ffffffffffffffff82111715614ac457614ac4613ee0565b604052614ad0836143a2565b8152602083015160208201528091505092915050565b600060808284031215614af857600080fd5b6040516080810181811067ffffffffffffffff82111715614b1b57614b1b613ee0565b60405282518152614b2e602084016143a2565b6020820152614b3f604084016143a2565b6040820152606083015160608201528091505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfe307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000810000a",
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
	return address, tx, &KeeperRegistry{KeeperRegistryCaller: KeeperRegistryCaller{contract: contract}, KeeperRegistryTransactor: KeeperRegistryTransactor{contract: contract}, KeeperRegistryFilterer: KeeperRegistryFilterer{contract: contract}}, nil
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

func (_KeeperRegistry *KeeperRegistryCaller) GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getFastGasFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) GetFastGasFeedAddress() (common.Address, error) {
	return _KeeperRegistry.Contract.GetFastGasFeedAddress(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) GetFastGasFeedAddress() (common.Address, error) {
	return _KeeperRegistry.Contract.GetFastGasFeedAddress(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCaller) GetLinkAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getLinkAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) GetLinkAddress() (common.Address, error) {
	return _KeeperRegistry.Contract.GetLinkAddress(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) GetLinkAddress() (common.Address, error) {
	return _KeeperRegistry.Contract.GetLinkAddress(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCaller) GetLinkNativeFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getLinkNativeFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) GetLinkNativeFeedAddress() (common.Address, error) {
	return _KeeperRegistry.Contract.GetLinkNativeFeedAddress(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) GetLinkNativeFeedAddress() (common.Address, error) {
	return _KeeperRegistry.Contract.GetLinkNativeFeedAddress(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCaller) GetMode(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getMode")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) GetMode() (uint8, error) {
	return _KeeperRegistry.Contract.GetMode(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) GetMode() (uint8, error) {
	return _KeeperRegistry.Contract.GetMode(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCaller) GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getTriggerType", upkeepId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _KeeperRegistry.Contract.GetTriggerType(&_KeeperRegistry.CallOpts, upkeepId)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _KeeperRegistry.Contract.GetTriggerType(&_KeeperRegistry.CallOpts, upkeepId)
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
	Id        *big.Int
	TriggerID [32]byte
	Raw       types.Log
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
	Id        *big.Int
	TriggerID [32]byte
	Raw       types.Log
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
	Id        *big.Int
	TriggerID [32]byte
	Raw       types.Log
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
	Id        *big.Int
	TriggerID [32]byte
	Raw       types.Log
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
	TriggerID    [32]byte
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
	ExecuteGas uint32
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
	case _KeeperRegistry.abi.Events["CancelledUpkeepReport"].ID:
		return _KeeperRegistry.ParseCancelledUpkeepReport(log)
	case _KeeperRegistry.abi.Events["ConfigSet"].ID:
		return _KeeperRegistry.ParseConfigSet(log)
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

func (KeeperRegistryCancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x846af3fdde418e2757c00c4cf8c6827e65ac4179753128c454922ccd7963886d")
}

func (KeeperRegistryConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (KeeperRegistryFundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (KeeperRegistryFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (KeeperRegistryInsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xfcb5fa059ae0ddeeac53240ccaa0a8ad80993b33e18e8821f1e1cbc53481240e")
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
	return common.HexToHash("0x07c4f19521a82d67071422fca888431cac8267935db8abe73e354791a3d38868")
}

func (KeeperRegistryStaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x909d8c7e883c6d2d53a5fb2dbb9ec57875a53c3769c099230027312ffb2a6b11")
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
	return common.HexToHash("0xdacc0b3d36199a5f78afb27c76cd01bf100f605edf5ee0c0e23055b0553b68c6")
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

	GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error)

	GetLinkAddress(opts *bind.CallOpts) (common.Address, error)

	GetLinkNativeFeedAddress(opts *bind.CallOpts) (common.Address, error)

	GetMode(opts *bind.CallOpts) (uint8, error)

	GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error)

	LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

		error)

	LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfigBytes []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	SimulatePerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, rawReport []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error)

	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryCancelledUpkeepReportIterator, error)

	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryCancelledUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseCancelledUpkeepReport(log types.Log) (*KeeperRegistryCancelledUpkeepReport, error)

	FilterConfigSet(opts *bind.FilterOpts) (*KeeperRegistryConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*KeeperRegistryConfigSet, error)

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
