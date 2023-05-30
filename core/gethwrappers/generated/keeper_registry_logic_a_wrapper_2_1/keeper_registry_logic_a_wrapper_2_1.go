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
<<<<<<< HEAD
<<<<<<< HEAD
	ABI: "[{\"inputs\":[{\"internalType\":\"contractKeeperRegistryLogicB2_1\",\"name\":\"logicB\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnchainConfigNonEmpty\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"acceptUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumUpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumUpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fallbackTo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFastGasFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkNativeFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMode\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase2_1.Mode\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTriggerType\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase2_1.Trigger\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_next\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"migrateUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"pauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedUpkeeps\",\"type\":\"bytes\"}],\"name\":\"receiveUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"name\":\"setUpkeepOffchainConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"unpauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"updateCheckData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTranscoderVersion\",\"outputs\":[{\"internalType\":\"enumUpkeepFormat\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepVersion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101206040523480156200001257600080fd5b5060405162006102380380620061028339810160408190526200003591620003aa565b80816001600160a01b0316634b4fd03b6040518163ffffffff1660e01b815260040160206040518083038186803b1580156200007057600080fd5b505afa15801562000085573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620000ab9190620003d1565b826001600160a01b031663ca30e6036040518163ffffffff1660e01b815260040160206040518083038186803b158015620000e557600080fd5b505afa158015620000fa573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001209190620003aa565b836001600160a01b031663b10b673c6040518163ffffffff1660e01b815260040160206040518083038186803b1580156200015a57600080fd5b505afa1580156200016f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001959190620003aa565b846001600160a01b0316636709d0e56040518163ffffffff1660e01b815260040160206040518083038186803b158015620001cf57600080fd5b505afa158015620001e4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200020a9190620003aa565b3380600081620002615760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b038481169190911790915581161562000294576200029481620002fe565b505050836002811115620002ac57620002ac620003f4565b60e0816002811115620002c357620002c3620003f4565b60f81b9052506001600160601b0319606093841b811660805291831b821660a052821b811660c05292901b9091166101005250620004239050565b6001600160a01b038116331415620003595760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000258565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215620003bd57600080fd5b8151620003ca816200040a565b9392505050565b600060208284031215620003e457600080fd5b815160038110620003ca57600080fd5b634e487b7160e01b600052602160045260246000fd5b6001600160a01b03811681146200042057600080fd5b50565b60805160601c60a05160601c60c05160601c60e05160f81c6101005160601c615c45620004bd600039600081816101f20152818161025101526102e00152600081816102b50152818161352701528181613f50015261411d01526000818161034101526132e301526000818161049001526133d70152600081816104fc015281816111ab01528181611df501526126ce0152615c456000f3fe60806040523480156200001157600080fd5b5060043610620001f05760003560e01c80638dcf0fe71162000111578063b121e14711620000a5578063ca30e603116200007b578063ca30e60314620004fa578063eb5dcd6c1462000521578063f2fde38b1462000538578063f7d334ba146200054f57620001f0565b8063b121e14714620004b5578063b148ab6b14620004cc578063c804802214620004e357620001f0565b8063a710b22111620000e7578063a710b2211462000444578063a72aa27e146200045b578063aab9edd61462000472578063b10b673c146200048e57620001f0565b80638dcf0fe714620003ff5780638e86139b14620004165780639fab4386146200042d57620001f0565b80636709d0e5116200018957806379ba5097116200015f57806379ba509714620003a857806385c1b0ba14620003b25780638765ecbe14620003c95780638da5cb5b14620003e057620001f0565b80636709d0e5146200033f57806371791aa01462000366578063744bfe61146200039157620001f0565b80634b4fd03b11620001cb5780634b4fd03b14620002b35780634ff597c114620002da5780635147cd5914620003025780635165f2f5146200032857620001f0565b80631a2af0111462000236578063349e8cca146200024f57806348013d7b146200029b575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e80801562000231573d6000f35b3d6000fd5b6200024d6200024736600462004cba565b62000566565b005b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b620002a4600081565b6040516200029291906200516a565b7f0000000000000000000000000000000000000000000000000000000000000000620002a4565b620002717f000000000000000000000000000000000000000000000000000000000000000081565b620003196200031336600462004c86565b50600090565b6040516200029291906200517f565b6200024d6200033936600462004c86565b620006bf565b7f000000000000000000000000000000000000000000000000000000000000000062000271565b6200037d6200037736600462004d32565b6200084e565b6040516200029296959493929190620050b3565b6200024d620003a236600462004cba565b62000e91565b6200024d620012ba565b6200024d620003c336600462004a21565b620013bd565b6200024d620003da36600462004c86565b62001e85565b60005473ffffffffffffffffffffffffffffffffffffffff1662000271565b6200024d6200041036600462004ce2565b62002021565b6200024d6200042736600462004c07565b62002089565b6200024d6200043e36600462004ce2565b62002468565b6200024d62000455366004620049e3565b6200251e565b6200024d6200046c36600462004d73565b620027b3565b6200047b600281565b60405160ff909116815260200162000292565b7f000000000000000000000000000000000000000000000000000000000000000062000271565b6200024d620004c6366004620049bc565b62002899565b6200024d620004dd36600462004c86565b62002992565b6200024d620004f436600462004c86565b62002b97565b7f000000000000000000000000000000000000000000000000000000000000000062000271565b6200024d62000532366004620049e3565b62002f72565b6200024d62000549366004620049bc565b620030d4565b6200037d6200056036600462004c86565b620030ec565b6200057182620031b9565b73ffffffffffffffffffffffffffffffffffffffff8116331415620005c2576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff811662000610576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526006602052604090205473ffffffffffffffffffffffffffffffffffffffff828116911614620006bb5760008281526006602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff851690811790915590519091339185917fb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b3591a45b5050565b620006ca81620031b9565b600081815260046020908152604091829020825161010081018452815463ffffffff8082168352640100000000820481169483019490945260ff68010000000000000000820416151594820185905273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091048116606083015260018301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a08401527801000000000000000000000000000000000000000000000000900490931660c082015260029091015490911660e082015290620007dd576040517f1b88a78400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff1690556200081e6002836200326e565b5060405182907f7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a4745690600090a25050565b600060606000806000806200086262003285565b6000600f604051806101200160405290816000820160009054906101000a900460ff1660ff1660ff1681526020016000820160019054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160059054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160099054906101000a900462ffffff1662ffffff1662ffffff16815260200160008201600c9054906101000a900461ffff1661ffff1661ffff16815260200160008201600e9054906101000a900460ff1615151515815260200160008201600f9054906101000a900460ff161515151581526020016000820160109054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160008201601c9054906101000a900463ffffffff1663ffffffff1663ffffffff168152505090506000600460008b8152602001908152602001600020604051806101000160405290816000820160009054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160049054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160089054906101000a900460ff161515151581526020016000820160099054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160018201600c9054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff1681526020016001820160189054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016002820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681525050905063ffffffff8016816020015163ffffffff161462000bd057505060408051602081019091526000808252965094506001935085915081905062000e87565b80604001511562000c0157505060408051602081019091526000808252965094506002935085915081905062000e87565b62000c0c82620032c0565b825160125492965090945060009162000c4c9185917801000000000000000000000000000000000000000000000000900463ffffffff16888886620034d2565b9050806bffffffffffffffffffffffff168260a001516bffffffffffffffffffffffff16101562000c9a57600060405180602001604052806000815250600698509850985050505062000e87565b5a95506000636e04ff0d60e01b8b60405160240162000cba919062005155565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009094169390931790925260e08501516012549251919350600092839273ffffffffffffffffffffffffffffffffffffffff9092169163ffffffff9091169062000d7890869062004f2d565b60006040518083038160008787f1925050503d806000811462000db8576040519150601f19603f3d011682016040523d82523d6000602084013e62000dbd565b606091505b50915091505a62000dcf908a6200540d565b98508162000de1576003995062000e7d565b8080602001905181019062000df7919062004bb0565b909c5090508b62000e285760006040518060200160405280600081525060049b509b509b5050505050505062000e87565b6012548151780100000000000000000000000000000000000000000000000090910463ffffffff16101562000e7d5760006040518060200160405280600081525060059b509b509b5050505050505062000e87565b909a509850505050505b9295509295509295565b600f546f01000000000000000000000000000000900460ff161562000ee2576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600f80547fffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffff166f0100000000000000000000000000000017905573ffffffffffffffffffffffffffffffffffffffff811662000f6a576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000828152600460209081526040808320815161010081018352815463ffffffff8082168352640100000000820481168387015260ff6801000000000000000083041615158386015273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009092048216606084015260018401546bffffffffffffffffffffffff80821660808601526c0100000000000000000000000082041660a0850152780100000000000000000000000000000000000000000000000090041660c0830152600290920154821660e08201528685526005909352922054909116331462001083576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6200108d62003521565b816020015163ffffffff161115620010d1576040517fff84e5dd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000838152600460205260409020600101546015546c010000000000000000000000009091046bffffffffffffffffffffffff1690620011139082906200540d565b60155560008481526004602081905260409182902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff16905590517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff858116928201929092526bffffffffffffffffffffffff831660248201527f00000000000000000000000000000000000000000000000000000000000000009091169063a9059cbb90604401602060405180830381600087803b158015620011f257600080fd5b505af115801562001207573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200122d919062004b90565b50604080516bffffffffffffffffffffffff8316815273ffffffffffffffffffffffffffffffffffffffff8516602082015285917ff3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318910160405180910390a25050600f80547fffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffff1690555050565b60015473ffffffffffffffffffffffffffffffffffffffff16331462001341576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600173ffffffffffffffffffffffffffffffffffffffff821660009081526016602052604090205460ff166003811115620013fc57620013fc6200557a565b14158015620014485750600373ffffffffffffffffffffffffffffffffffffffff821660009081526016602052604090205460ff1660038111156200144557620014456200557a565b14155b1562001480576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6010546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff16620014e0576040517fd12d7d8d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8162001518576040517f2c2fc94100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805161010081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526000808567ffffffffffffffff81111562001577576200157762005607565b604051908082528060200260200182016040528015620015ac57816020015b6060815260200190600190039081620015965790505b50905060008667ffffffffffffffff811115620015cd57620015cd62005607565b604051908082528060200260200182016040528015620015f7578160200160208202803683370190505b50905060008767ffffffffffffffff81111562001618576200161862005607565b604051908082528060200260200182016040528015620016a757816020015b604080516101008101825260008082526020808301829052928201819052606082018190526080820181905260a0820181905260c0820181905260e082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909201910181620016375790505b50905060008867ffffffffffffffff811115620016c857620016c862005607565b604051908082528060200260200182016040528015620016fd57816020015b6060815260200190600190039081620016e75790505b50905060005b8981101562001b9e578a8a82818110620017215762001721620055d8565b60209081029290920135600081815260048452604090819020815161010081018352815463ffffffff8082168352640100000000820481169783019790975260ff6801000000000000000082041615159382019390935273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009093048316606082015260018201546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490951660c0860152600201541660e0840152995090975062001812905088620031b9565b86838281518110620018285762001828620055d8565b602002602001018190525060076000898152602001908152602001600020805462001853906200548a565b80601f016020809104026020016040519081016040528092919081815260200182805462001881906200548a565b8015620018d25780601f10620018a657610100808354040283529160200191620018d2565b820191906000526020600020905b815481529060010190602001808311620018b457829003601f168201915b5050505050858281518110620018ec57620018ec620055d8565b60200260200101819052506005600089815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16848281518110620019405762001940620055d8565b73ffffffffffffffffffffffffffffffffffffffff909216602092830291909101820152600089815260179091526040902080546200197f906200548a565b80601f0160208091040260200160405190810160405280929190818152602001828054620019ad906200548a565b8015620019fe5780601f10620019d257610100808354040283529160200191620019fe565b820191906000526020600020905b815481529060010190602001808311620019e057829003601f168201915b505050505082828151811062001a185762001a18620055d8565b60200260200101819052508660a001516bffffffffffffffffffffffff168662001a439190620052b5565b600089815260046020908152604080832080547fffffff00000000000000000000000000000000000000000000000000000000001681556001810180547fffffffff0000000000000000000000000000000000000000000000000000000016905560020180547fffffffffffffffffffffffff00000000000000000000000000000000000000001690556007909152812091975062001ae391906200443c565b600088815260066020526040902080547fffffffffffffffffffffffff000000000000000000000000000000000000000016905562001b24600289620035ee565b5060a0870151604080516bffffffffffffffffffffffff909216825273ffffffffffffffffffffffffffffffffffffffff8b16602083015289917fb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff910160405180910390a28062001b9581620054e0565b91505062001703565b508460155462001baf91906200540d565b60155560405160009062001bd2908c908c90869089908990889060200162004f4b565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905260105490915073ffffffffffffffffffffffffffffffffffffffff808b1691638e86139b916c010000000000000000000000009091041663c71249ab60028d73ffffffffffffffffffffffffffffffffffffffff1663aab9edd66040518163ffffffff1660e01b8152600401602060405180830381600087803b15801562001c8757600080fd5b505af115801562001c9c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001cc2919062004df7565b866040518463ffffffff1660e01b815260040162001ce39392919062005196565b60006040518083038186803b15801562001cfc57600080fd5b505afa15801562001d11573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405262001d59919081019062004c4d565b6040518263ffffffff1660e01b815260040162001d77919062005155565b600060405180830381600087803b15801562001d9257600080fd5b505af115801562001da7573d6000803e3d6000fd5b50506040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8c81166004830152602482018a90527f000000000000000000000000000000000000000000000000000000000000000016925063a9059cbb9150604401602060405180830381600087803b15801562001e3c57600080fd5b505af115801562001e51573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001e77919062004b90565b505050505050505050505050565b62001e9081620031b9565b600081815260046020908152604091829020825161010081018452815463ffffffff8082168352640100000000820481169483019490945260ff680100000000000000008204161580159583019590955273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091048116606083015260018301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a08401527801000000000000000000000000000000000000000000000000900490931660c082015260029091015490911660e08201529062001fa5576040517f514b6c2400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff166801000000000000000017905562001ff1600283620035ee565b5060405182907f8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f90600090a25050565b6200202c83620031b9565b6000838152601760205260409020620020479083836200447b565b50827f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf485083836040516200207c92919062005108565b60405180910390a2505050565b60023360009081526016602052604090205460ff166003811115620020b257620020b26200557a565b14158015620020e8575060033360009081526016602052604090205460ff166003811115620020e557620020e56200557a565b14155b1562002120576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080808080620021348688018862004aae565b9450945094509450945060005b85518110156200245e57600073ffffffffffffffffffffffffffffffffffffffff16858281518110620021785762002178620055d8565b60200260200101516060015173ffffffffffffffffffffffffffffffffffffffff1614156200228057858181518110620021b657620021b6620055d8565b6020026020010151858281518110620021d357620021d3620055d8565b602002602001015160e00151604051620021ed9062004528565b91825273ffffffffffffffffffffffffffffffffffffffff166020820152604001604051809103906000f0801580156200222b573d6000803e3d6000fd5b50858281518110620022415762002241620055d8565b60200260200101516060019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250505b620023a3868281518110620022995762002299620055d8565b6020026020010151868381518110620022b657620022b6620055d8565b602002602001015160e00151878481518110620022d757620022d7620055d8565b602002602001015160000151868581518110620022f857620022f8620055d8565b6020026020010151898681518110620023155762002315620055d8565b602002602001015160a00151898781518110620023365762002336620055d8565b60200260200101518b8881518110620023535762002353620055d8565b602002602001015160400151898981518110620023745762002374620055d8565b60200260200101518d8a81518110620023915762002391620055d8565b602002602001015160600151620035fc565b858181518110620023b857620023b8620055d8565b60200260200101517f74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71868381518110620023f657620023f6620055d8565b602002602001015160a0015133604051620024419291906bffffffffffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a2806200245581620054e0565b91505062002141565b5050505050505050565b6200247383620031b9565b60125474010000000000000000000000000000000000000000900463ffffffff16811115620024ce576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000838152600760205260409020620024e99083836200447b565b50827f7b778136e5211932b51a145badd01959415e79e051a933604b3d323f862dcabf83836040516200207c92919062005108565b73ffffffffffffffffffffffffffffffffffffffff81166200256c576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600c6020526040902054163314620025cd576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600f54600b546000916200260691859170010000000000000000000000000000000090046bffffffffffffffffffffffff169062003ac9565b73ffffffffffffffffffffffffffffffffffffffff8416600090815260086020526040902080547fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff16905560155490915062002672906bffffffffffffffffffffffff8316906200540d565b6015556040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff83811660048301526bffffffffffffffffffffffff831660248301527f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb90604401602060405180830381600087803b1580156200271357600080fd5b505af115801562002728573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200274e919062004b90565b5060405133815273ffffffffffffffffffffffffffffffffffffffff808416916bffffffffffffffffffffffff8416918616907f9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f406989060200160405180910390a4505050565b6108fc8163ffffffff161080620027dd575060125463ffffffff6401000000009091048116908216115b1562002815576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6200282082620031b9565b60008281526004602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff8516908117909155915191825283917fc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c910160405180910390a25050565b73ffffffffffffffffffffffffffffffffffffffff8181166000908152600d6020526040902054163314620028fa576040517f6752e7aa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8181166000818152600c602090815260408083208054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217909355600d909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b600081815260046020908152604091829020825161010081018452815463ffffffff80821683526401000000008204811694830185905260ff6801000000000000000083041615159583019590955273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091048116606083015260018301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a084015278010000000000000000000000000000000000000000000000009004851660c083015260029092015490911660e0820152911462002aa5576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526006602052604090205473ffffffffffffffffffffffffffffffffffffffff16331462002b03576040517f6352a85300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526005602090815260408083208054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821790935560069094528285208054909216909155905173ffffffffffffffffffffffffffffffffffffffff90911692839186917f5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c91a4505050565b6000818152600460209081526040808320815161010081018352815463ffffffff80821683526401000000008204811695830186905260ff6801000000000000000083041615159483019490945273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091048116606083015260018301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a084015278010000000000000000000000000000000000000000000000009004841660c083015260029092015490911660e0820152929114159062002c9560005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614905081801562002cf0575080801562002cee575062002ce162003521565b836020015163ffffffff16115b155b1562002d28576040517ffbc0357800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8015801562002d5b575060008481526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314155b1562002d93576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600062002d9f62003521565b90508162002db75762002db4603282620052b5565b90505b6000858152600460205260409020805463ffffffff808416640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff9092169190911790915562002e12906002908790620035ee16565b5060105460808501516bffffffffffffffffffffffff918216916000911682111562002e7b57608086015162002e49908362005427565b90508560a001516bffffffffffffffffffffffff16816bffffffffffffffffffffffff16111562002e7b575060a08501515b808660a0015162002e8d919062005427565b600088815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff9384160217905560115462002ef591839116620052f8565b601180547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff9290921691909117905560405167ffffffffffffffff84169088907f91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f79118190600090a350505050505050565b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600c602052604090205416331462002fd3576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff811633141562003024576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600d6020526040902054811690821614620006bb5773ffffffffffffffffffffffffffffffffffffffff8281166000818152600d602052604080822080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169486169485179055513392917f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836791a45050565b620030de62003cfc565b620030e98162003d7f565b50565b60006060600080600080620031a587600760008a815260200190815260200160002080546200311b906200548a565b80601f016020809104026020016040519081016040528092919081815260200182805462003149906200548a565b80156200319a5780601f106200316e576101008083540402835291602001916200319a565b820191906000526020600020905b8154815290600101906020018083116200317c57829003601f168201915b50505050506200084e565b949c939b5091995097509550909350915050565b60008181526005602052604090205473ffffffffffffffffffffffffffffffffffffffff16331462003217576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600081815260046020526040902054640100000000900463ffffffff90811614620030e9576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006200327c838362003e77565b90505b92915050565b3215620032be576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b6000806000836060015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b1580156200334857600080fd5b505afa1580156200335d573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003383919062004da2565b50945090925050506000811315806200339b57508142105b80620033c05750828015620033c05750620033b782426200540d565b8463ffffffff16105b15620033d1576013549550620033d5565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b1580156200343c57600080fd5b505afa15801562003451573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003477919062004da2565b50945090925050506000811315806200348f57508142105b80620034b45750828015620034b45750620034ab82426200540d565b8463ffffffff16105b15620034c5576014549450620034c9565b8094505b50505050915091565b600080620034e586896000015162003ec9565b9050600080620035028a8a63ffffffff16858a8a60018b62003f17565b9092509050620035138183620052f8565b9a9950505050505050505050565b600060017f000000000000000000000000000000000000000000000000000000000000000060028111156200355a576200355a6200557a565b1415620035e957606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b158015620035a957600080fd5b505afa158015620035be573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620035e4919062004ca0565b905090565b504390565b60006200327c838362004338565b600f546e010000000000000000000000000000900460ff16156200364c576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff88163b6200369b576040517f09ee12d500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60125484517401000000000000000000000000000000000000000090910463ffffffff161015620036f8576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108fc8763ffffffff16108062003722575060125463ffffffff6401000000009091048116908816115b156200375a576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008981526004602052604090206002015473ffffffffffffffffffffffffffffffffffffffff1615620037ba576040517f6e3b930b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040518061010001604052808863ffffffff16815260200163ffffffff8016815260200184151581526020018273ffffffffffffffffffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff168152602001866bffffffffffffffffffffffff168152602001600063ffffffff1681526020018973ffffffffffffffffffffffffffffffffffffffff16815250600460008b815260200190815260200160002060008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160086101000a81548160ff02191690831515021790555060608201518160000160096101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060808201518160010160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060a082015181600101600c6101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060c08201518160010160186101000a81548163ffffffff021916908363ffffffff16021790555060e08201518160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505085600560008b815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550846bffffffffffffffffffffffff1660155462003a699190620052b5565b6015556000898152600760209081526040909120855162003a8d9287019062004536565b506000898152601760209081526040909120835162003aaf9285019062004536565b5062003abd60028a6200326e565b50505050505050505050565b73ffffffffffffffffffffffffffffffffffffffff831660009081526008602090815260408083208151608081018352905460ff80821615158352610100820416938201939093526bffffffffffffffffffffffff6201000084048116928201929092526e0100000000000000000000000000009092041660608201819052829062003b56908662005427565b9050600062003b66858362005342565b9050808360400181815162003b7c9190620052f8565b6bffffffffffffffffffffffff908116909152871660608501525062003ba38582620053df565b62003baf908362005427565b6011805460009062003bd19084906bffffffffffffffffffffffff16620052f8565b825461010092830a6bffffffffffffffffffffffff81810219909216928216029190911790925573ffffffffffffffffffffffffffffffffffffffff999099166000908152600860209081526040918290208751815492890151938901516060909901517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009093169015157fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff161760ff909316909b02919091177fffffffffffff000000000000000000000000000000000000000000000000ffff1662010000878416027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff16176e010000000000000000000000000000919092160217909755509095945050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314620032be576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162001338565b73ffffffffffffffffffffffffffffffffffffffff811633141562003e01576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162001338565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600081815260018301602052604081205462003ec0575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556200327f565b5060006200327f565b600062003ede63ffffffff8416601462005370565b62003eeb836001620052d0565b62003efc9060ff16611d4c62005370565b62003f0b90620124f8620052b5565b6200327c9190620052b5565b6000806000896080015161ffff168762003f32919062005370565b905083801562003f415750803a105b1562003f4a57503a5b600060027f0000000000000000000000000000000000000000000000000000000000000000600281111562003f835762003f836200557a565b14156200411957604080516000815260208101909152851562003fe85760003660405180608001604052806048815260200162005bf16048913960405160200162003fd19392919062004f04565b60405160208183030381529060405290506200406a565b6012546200401a907801000000000000000000000000000000000000000000000000900463ffffffff166004620053b0565b63ffffffff1667ffffffffffffffff8111156200403b576200403b62005607565b6040519080825280601f01601f19166020018201604052801562004066576020820181803683370190505b5090505b6040517f49948e0e00000000000000000000000000000000000000000000000000000000815273420000000000000000000000000000000000000f906349948e0e90620040bc90849060040162005155565b60206040518083038186803b158015620040d557600080fd5b505afa158015620040ea573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062004110919062004ca0565b915050620041dd565b60017f000000000000000000000000000000000000000000000000000000000000000060028111156200415057620041506200557a565b1415620041dd57606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b815260040160206040518083038186803b1580156200419f57600080fd5b505afa158015620041b4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620041da919062004ca0565b90505b84620041fc57808b6080015161ffff16620041f9919062005370565b90505b6200420c61ffff8716826200532b565b9050600087826200421e8c8e620052b5565b6200422a908662005370565b620042369190620052b5565b6200424a90670de0b6b3a764000062005370565b6200425691906200532b565b905060008c6040015163ffffffff1664e8d4a5100062004277919062005370565b898e6020015163ffffffff16858f8862004292919062005370565b6200429e9190620052b5565b620042ae90633b9aca0062005370565b620042ba919062005370565b620042c691906200532b565b620042d29190620052b5565b90506b033b2e3c9fd0803ce8000000620042ed8284620052b5565b111562004326576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b60008181526001830160205260408120548015620044315760006200435f6001836200540d565b855490915060009062004375906001906200540d565b9050818114620043e1576000866000018281548110620043995762004399620055d8565b9060005260206000200154905080876000018481548110620043bf57620043bf620055d8565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080620043f557620043f5620055a9565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506200327f565b60009150506200327f565b5080546200444a906200548a565b6000825580601f106200445b575050565b601f016020900490600052602060002090810190620030e99190620045b3565b82805462004489906200548a565b90600052602060002090601f016020900481019282620044ad576000855562004516565b82601f10620044e6578280017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082351617855562004516565b8280016001018555821562004516579182015b8281111562004516578235825591602001919060010190620044f9565b5062004524929150620045b3565b5090565b610575806200567c83390190565b82805462004544906200548a565b90600052602060002090601f01602090048101928262004568576000855562004516565b82601f106200458357805160ff191683800117855562004516565b8280016001018555821562004516579182015b828111156200451657825182559160200191906001019062004596565b5b80821115620045245760008155600101620045b4565b8035620045d78162005649565b919050565b600082601f830112620045ee57600080fd5b8135602062004607620046018362005245565b620051f3565b80838252828201915082860187848660051b89010111156200462857600080fd5b60005b8581101562004654578135620046418162005649565b845292840192908401906001016200462b565b5090979650505050505050565b600082601f8301126200467357600080fd5b8135602062004686620046018362005245565b80838252828201915082860187848660051b8901011115620046a757600080fd5b6000805b86811015620046ef57823567ffffffffffffffff811115620046cb578283fd5b620046db8b88838d0101620048c5565b8652509385019391850191600101620046ab565b509198975050505050505050565b600082601f8301126200470f57600080fd5b8135602062004722620046018362005245565b80838252828201915082860187848660081b89010111156200474357600080fd5b6000805b86811015620046ef5761010080848c03121562004762578283fd5b6200476c620051c6565b62004777856200496f565b8152620047868886016200496f565b888201526040808601356200479b816200566c565b908201526060620047ae868201620045ca565b908201526080620047c18682016200499f565b9082015260a0620047d48682016200499f565b9082015260c0620047e78682016200496f565b9082015260e0620047fa868201620045ca565b90820152865294860194929092019160010162004747565b600082601f8301126200482457600080fd5b8135602062004837620046018362005245565b80838252828201915082860187848660051b89010111156200485857600080fd5b60005b8581101562004654578135845292840192908401906001016200485b565b60008083601f8401126200488c57600080fd5b50813567ffffffffffffffff811115620048a557600080fd5b602083019150836020828501011115620048be57600080fd5b9250929050565b600082601f830112620048d757600080fd5b8135620048e862004601826200526c565b818152846020838601011115620048fe57600080fd5b816020850160208301376000918101602001919091529392505050565b600082601f8301126200492d57600080fd5b81516200493e62004601826200526c565b8181528460208386010111156200495457600080fd5b6200496782602083016020870162005457565b949350505050565b803563ffffffff81168114620045d757600080fd5b805169ffffffffffffffffffff81168114620045d757600080fd5b80356bffffffffffffffffffffffff81168114620045d757600080fd5b600060208284031215620049cf57600080fd5b8135620049dc8162005649565b9392505050565b60008060408385031215620049f757600080fd5b823562004a048162005649565b9150602083013562004a168162005649565b809150509250929050565b60008060006040848603121562004a3757600080fd5b833567ffffffffffffffff8082111562004a5057600080fd5b818601915086601f83011262004a6557600080fd5b81358181111562004a7557600080fd5b8760208260051b850101111562004a8b57600080fd5b6020928301955093505084013562004aa38162005649565b809150509250925092565b600080600080600060a0868803121562004ac757600080fd5b853567ffffffffffffffff8082111562004ae057600080fd5b62004aee89838a0162004812565b9650602088013591508082111562004b0557600080fd5b62004b1389838a01620046fd565b9550604088013591508082111562004b2a57600080fd5b62004b3889838a0162004661565b9450606088013591508082111562004b4f57600080fd5b62004b5d89838a01620045dc565b9350608088013591508082111562004b7457600080fd5b5062004b838882890162004661565b9150509295509295909350565b60006020828403121562004ba357600080fd5b8151620049dc816200566c565b6000806040838503121562004bc457600080fd5b825162004bd1816200566c565b602084015190925067ffffffffffffffff81111562004bef57600080fd5b62004bfd858286016200491b565b9150509250929050565b6000806020838503121562004c1b57600080fd5b823567ffffffffffffffff81111562004c3357600080fd5b62004c418582860162004879565b90969095509350505050565b60006020828403121562004c6057600080fd5b815167ffffffffffffffff81111562004c7857600080fd5b62004967848285016200491b565b60006020828403121562004c9957600080fd5b5035919050565b60006020828403121562004cb357600080fd5b5051919050565b6000806040838503121562004cce57600080fd5b82359150602083013562004a168162005649565b60008060006040848603121562004cf857600080fd5b83359250602084013567ffffffffffffffff81111562004d1757600080fd5b62004d258682870162004879565b9497909650939450505050565b6000806040838503121562004d4657600080fd5b82359150602083013567ffffffffffffffff81111562004d6557600080fd5b62004bfd85828601620048c5565b6000806040838503121562004d8757600080fd5b8235915062004d99602084016200496f565b90509250929050565b600080600080600060a0868803121562004dbb57600080fd5b62004dc68662004984565b945060208601519350604086015192506060860151915062004deb6080870162004984565b90509295509295909350565b60006020828403121562004e0a57600080fd5b815160ff81168114620049dc57600080fd5b600081518084526020808501945080840160005b8381101562004e6457815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010162004e30565b509495945050505050565b6000815180845260208085019450848260051b860182860160005b858110156200465457838303895262004ea583835162004eb8565b9885019892509084019060010162004e8a565b6000815180845262004ed281602086016020860162005457565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b82848237600083820160008152835162004f2381836020880162005457565b0195945050505050565b6000825162004f4181846020870162005457565b9190910192915050565b600060a0808352878184015260c07f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff89111562004f8757600080fd5b8860051b808b838701378085019050818101600081526020838784030181880152818b5180845260e093508385019150828d01945060005b818110156200505f578551805163ffffffff90811685528582015181168686015260408083015115159086015260608083015173ffffffffffffffffffffffffffffffffffffffff908116918701919091526080808401516bffffffffffffffffffffffff908116918801919091528b840151168b8701528983015190911689860152908601511685840152948301946101009092019160010162004fbf565b5050878103604089015262005075818c62004e6f565b95505050505050828103606084015262005090818662004e1c565b90508281036080840152620050a6818562004e6f565b9998505050505050505050565b861515815260c060208201526000620050d060c083018862004eb8565b905060078610620050e557620050e56200557a565b8560408301528460608301528360808301528260a0830152979650505050505050565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b6020815260006200327c602083018462004eb8565b60208101620051798362005636565b91905290565b60208101600483106200517957620051796200557a565b60ff8416815260ff83166020820152606060408201526000620051bd606083018462004eb8565b95945050505050565b604051610100810167ffffffffffffffff81118282101715620051ed57620051ed62005607565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156200523d576200523d62005607565b604052919050565b600067ffffffffffffffff82111562005262576200526262005607565b5060051b60200190565b600067ffffffffffffffff82111562005289576200528962005607565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b60008219821115620052cb57620052cb6200551c565b500190565b600060ff821660ff84168060ff03821115620052f057620052f06200551c565b019392505050565b60006bffffffffffffffffffffffff8083168185168083038211156200532257620053226200551c565b01949350505050565b6000826200533d576200533d6200554b565b500490565b60006bffffffffffffffffffffffff808416806200536457620053646200554b565b92169190910492915050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615620053ab57620053ab6200551c565b500290565b600063ffffffff80831681851681830481118215151615620053d657620053d66200551c565b02949350505050565b60006bffffffffffffffffffffffff80831681851681830481118215151615620053d657620053d66200551c565b6000828210156200542257620054226200551c565b500390565b60006bffffffffffffffffffffffff838116908316818110156200544f576200544f6200551c565b039392505050565b60005b83811015620054745781810151838201526020016200545a565b8381111562005484576000848401525b50505050565b600181811c908216806200549f57607f821691505b60208210811415620054da577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156200551557620055156200551c565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60038110620030e957620030e96200557a565b73ffffffffffffffffffffffffffffffffffffffff81168114620030e957600080fd5b8015158114620030e957600080fdfe60c060405234801561001057600080fd5b5060405161057538038061057583398101604081905261002f9161005b565b60008054336001600160a01b031990911617905560601b6001600160601b03191660805260a052610098565b6000806040838503121561006e57600080fd5b825160208401519092506001600160a01b038116811461008d57600080fd5b809150509250929050565b60805160601c60a0516104ae6100c76000396000610145015260008181610170015261028001526104ae6000f3fe608060405234801561001057600080fd5b50600436106100725760003560e01c806379188d161161005057806379188d161461011d5780638ee489b214610140578063f00e6a2a1461016e57600080fd5b8063181f5a77146100775780631a5da6c8146100c95780635ab1bd53146100de575b600080fd5b6100b36040518060400160405280601981526020017f4175746f6d6174696f6e466f7277617264657220312e302e300000000000000081525081565b6040516100c091906103ff565b60405180910390f35b6100dc6100d73660046102e9565b610194565b005b60005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100c0565b61013061012b366004610326565b61022c565b60405190151581526020016100c0565b6040517f000000000000000000000000000000000000000000000000000000000000000081526020016100c0565b7f00000000000000000000000000000000000000000000000000000000000000006100f8565b60005473ffffffffffffffffffffffffffffffffffffffff1633146101e5576040517fea8e4eb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b6000805473ffffffffffffffffffffffffffffffffffffffff16331461027e576040517fea8e4eb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f00000000000000000000000000000000000000000000000000000000000000005a6113888110156102af57600080fd5b6113888103905084604082048203116102c757600080fd5b50803b6102d357600080fd5b60008084516020860160008589f1949350505050565b6000602082840312156102fb57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461031f57600080fd5b9392505050565b6000806040838503121561033957600080fd5b82359150602083013567ffffffffffffffff8082111561035857600080fd5b818501915085601f83011261036c57600080fd5b81358181111561037e5761037e610472565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019083821181831017156103c4576103c4610472565b816040528281528860208487010111156103dd57600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b600060208083528351808285015260005b8181101561042c57858101830151858201604001528201610410565b8181111561043e576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000806000a",
=======
	ABI: "[{\"inputs\":[{\"internalType\":\"contractKeeperRegistryLogicB2_1\",\"name\":\"logicB\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTrigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepTriggerConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumUpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumUpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fallbackTo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFastGasFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkNativeFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMode\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase2_1.Mode\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTriggerType\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase2_1.Trigger\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_next\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"migrateUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedUpkeeps\",\"type\":\"bytes\"}],\"name\":\"receiveUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"setUpkeepTriggerConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTranscoderVersion\",\"outputs\":[{\"internalType\":\"enumUpkeepFormat\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepVersion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumKeeperRegistryBase2_1.Trigger\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"validateTriggerConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101206040523480156200001257600080fd5b5060405162005dcc38038062005dcc8339810160408190526200003591620003aa565b80816001600160a01b0316634b4fd03b6040518163ffffffff1660e01b815260040160206040518083038186803b1580156200007057600080fd5b505afa15801562000085573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620000ab9190620003d1565b826001600160a01b031663ca30e6036040518163ffffffff1660e01b815260040160206040518083038186803b158015620000e557600080fd5b505afa158015620000fa573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001209190620003aa565b836001600160a01b031663b10b673c6040518163ffffffff1660e01b815260040160206040518083038186803b1580156200015a57600080fd5b505afa1580156200016f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001959190620003aa565b846001600160a01b0316636709d0e56040518163ffffffff1660e01b815260040160206040518083038186803b158015620001cf57600080fd5b505afa158015620001e4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200020a9190620003aa565b3380600081620002615760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b038481169190911790915581161562000294576200029481620002fe565b505050836002811115620002ac57620002ac620003f4565b60e0816002811115620002c357620002c3620003f4565b60f81b9052506001600160601b0319606093841b811660805291831b821660a052821b811660c05292901b9091166101005250620004239050565b6001600160a01b038116331415620003595760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000258565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215620003bd57600080fd5b8151620003ca816200040a565b9392505050565b600060208284031215620003e457600080fd5b815160038110620003ca57600080fd5b634e487b7160e01b600052602160045260246000fd5b6001600160a01b03811681146200042057600080fd5b50565b60805160601c60a05160601c60c05160601c60e05160f81c6101005160601c61590f620004bd6000396000818161019e015281816101e4015261028c01526000818161024801528181613385015281816135cd015281816137ec01526139b90152600081816102d6015261312a0152600081816103d8015261321e01526000818161041601528181611b8601526121ef015261590f6000f3fe60806040523480156200001157600080fd5b50600436106200019c5760003560e01c806385c1b0ba11620000ed578063b10b673c1162000099578063f2fde38b116200006f578063f2fde38b146200043b578063f7d334ba1462000452578063fd9541d71462000469576200019c565b8063b10b673c14620003d6578063c804802214620003fd578063ca30e6031462000414576200019c565b80638e86139b11620000cf5780638e86139b146200038c578063948108f714620003a3578063aab9edd614620003ba576200019c565b806385c1b0ba14620003565780638da5cb5b146200036d576200019c565b80635147cd59116200014d5780636ded9eae116200012f5780636ded9eae14620002fb57806371791aa0146200032157806379ba5097146200034c576200019c565b80635147cd5914620002ae5780636709d0e514620002d4576200019c565b80634b4fd03b11620001835780634b4fd03b14620002465780634ee88d35146200026d5780634ff597c11462000286576200019c565b8063349e8cca14620001e257806348013d7b146200022e575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e808015620001dd573d6000f35b3d6000fd5b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b62000237600081565b60405162000225919062004dcf565b7f000000000000000000000000000000000000000000000000000000000000000062000237565b620002846200027e366004620048f2565b62000480565b005b620002047f000000000000000000000000000000000000000000000000000000000000000081565b620002c5620002bf366004620048d8565b6200053a565b60405162000225919062004de4565b7f000000000000000000000000000000000000000000000000000000000000000062000204565b620003126200030c366004620042ec565b620005f0565b60405190815260200162000225565b620003386200033236600462004942565b62000871565b604051620002259695949392919062004d18565b6200028462000efd565b62000284620003673660046200439f565b62001000565b60005473ffffffffffffffffffffffffffffffffffffffff1662000204565b620002846200039d366004620045c5565b62001c17565b62000284620003b436600462004975565b62002017565b620003c3600281565b60405160ff909116815260200162000225565b7f000000000000000000000000000000000000000000000000000000000000000062000204565b620002846200040e366004620048d8565b620022cf565b7f000000000000000000000000000000000000000000000000000000000000000062000204565b620002846200044c366004620042cc565b620026aa565b6200033862000463366004620048d8565b620026c2565b620002846200047a36600462004644565b6200278f565b6200048b8362002997565b600062000498846200053a565b9050620004dc8184848080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506200278f92505050565b6000848152601760205260409020620004f790848462003d2a565b50837f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d566484846040516200052c92919062004d6d565b60405180910390a250505050565b6000818160045b600f811015620005cf577fff000000000000000000000000000000000000000000000000000000000000008216838260208110620005835762000583620052a2565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614620005ba57506000949350505050565b80620005c6816200516c565b91505062000541565b5081600f1a6003811115620005e857620005e862005244565b949350505050565b6000805473ffffffffffffffffffffffffffffffffffffffff1633148015906200064257506011546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff163314155b156200067a576040517fd48b678b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080806200068c858701876200468e565b9250925092506200069e83836200278f565b620006a98362002a4c565b93506000848c604051620006bd9062003dd7565b91825273ffffffffffffffffffffffffffffffffffffffff166020820152604001604051809103906000f080158015620006fb573d6000803e3d6000fd5b50905062000749858d8d8d60008e8e8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092018290525092508c91508b90508a62002bdc565b6012805468010000000000000000900463ffffffff169060086200076d83620051a8565b91906101000a81548163ffffffff021916908363ffffffff16021790555050847fbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d0128c8c604051620007e692919063ffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a2847f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d56648460405162000820919062004dba565b60405180910390a2847f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850836040516200085a919062004dba565b60405180910390a250505050979650505050505050565b6000606060008060008062000885620030cc565b600062000892896200053a565b90506000600f604051806101200160405290816000820160009054906101000a900460ff1660ff1660ff1681526020016000820160019054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160059054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160099054906101000a900462ffffff1662ffffff1662ffffff16815260200160008201600c9054906101000a900461ffff1661ffff1661ffff16815260200160008201600e9054906101000a900460ff1615151515815260200160008201600f9054906101000a900460ff161515151581526020016000820160109054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160008201601c9054906101000a900463ffffffff1663ffffffff1663ffffffff168152505090506000600460008c8152602001908152602001600020604051806101000160405290816000820160009054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160049054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160089054906101000a900460ff161515151581526020016000820160099054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160018201600c9054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff1681526020016001820160189054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016002820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681525050905063ffffffff8016816020015163ffffffff161462000c0457505060408051602081019091526000808252975095506001945086925082915062000ef39050565b80604001511562000c3757505060408051602081019091526000808252975095506002945086925082915062000ef39050565b62000c428262003107565b825160125492975090955060009162000c829185917801000000000000000000000000000000000000000000000000900463ffffffff1689898662003319565b9050806bffffffffffffffffffffffff168260a001516bffffffffffffffffffffffff16101562000cd15760006040518060200160405280600081525060069950995099505050505062000ef3565b60019950600084600381111562000cec5762000cec62005244565b148062000d0d5750600184600381111562000d0b5762000d0b62005244565b145b1562000eee575a96506000636e04ff0d60e01b8c60405160240162000d33919062004dba565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009094169390931790925260e0850151601254925191935073ffffffffffffffffffffffffffffffffffffffff169163ffffffff169062000de890849062004b63565b60006040518083038160008787f1925050503d806000811462000e28576040519150601f19603f3d011682016040523d82523d6000602084013e62000e2d565b606091505b50909b5099505a62000e40908962005059565b97508a62000e52576003985062000eec565b8980602001905181019062000e68919062004554565b909b5099508a62000e985760006040518060200160405280600081525060049a509a509a50505050505062000ef3565b6012548a51780100000000000000000000000000000000000000000000000090910463ffffffff16101562000eec5760006040518060200160405280600081525060059a509a509a50505050505062000ef3565b505b505050505b9295509295509295565b60015473ffffffffffffffffffffffffffffffffffffffff16331462000f84576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600173ffffffffffffffffffffffffffffffffffffffff821660009081526016602052604090205460ff1660038111156200103f576200103f62005244565b141580156200108b5750600373ffffffffffffffffffffffffffffffffffffffff821660009081526016602052604090205460ff16600381111562001088576200108862005244565b14155b15620010c3576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6010546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1662001123576040517fd12d7d8d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b816200115b576040517f2c2fc94100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805161010081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526000808567ffffffffffffffff811115620011ba57620011ba620052d1565b604051908082528060200260200182016040528015620011ef57816020015b6060815260200190600190039081620011d95790505b50905060008667ffffffffffffffff811115620012105762001210620052d1565b6040519080825280602002602001820160405280156200123a578160200160208202803683370190505b50905060008767ffffffffffffffff8111156200125b576200125b620052d1565b604051908082528060200260200182016040528015620012ea57816020015b604080516101008101825260008082526020808301829052928201819052606082018190526080820181905260a0820181905260c0820181905260e082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816200127a5790505b50905060008867ffffffffffffffff8111156200130b576200130b620052d1565b6040519080825280602002602001820160405280156200134057816020015b60608152602001906001900390816200132a5790505b50905060008967ffffffffffffffff811115620013615762001361620052d1565b6040519080825280602002602001820160405280156200139657816020015b6060815260200190600190039081620013805790505b50905060005b8a8110156200192d578b8b82818110620013ba57620013ba620052a2565b60209081029290920135600081815260048452604090819020815161010081018352815463ffffffff8082168352640100000000820481169783019790975260ff6801000000000000000082041615159382019390935273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009093048316606082015260018201546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490951660c0860152600201541660e08401529a50909850620014ab90508962002997565b87848281518110620014c157620014c1620052a2565b6020026020010181905250600760008a81526020019081526020016000208054620014ec906200511c565b80601f01602080910402602001604051908101604052809291908181526020018280546200151a906200511c565b80156200156b5780601f106200153f576101008083540402835291602001916200156b565b820191906000526020600020905b8154815290600101906020018083116200154d57829003601f168201915b5050505050868281518110620015855762001585620052a2565b6020026020010181905250600560008a815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16858281518110620015d957620015d9620052a2565b73ffffffffffffffffffffffffffffffffffffffff90921660209283029190910182015260008a8152601790915260409020805462001618906200511c565b80601f016020809104026020016040519081016040528092919081815260200182805462001646906200511c565b8015620016975780601f106200166b5761010080835404028352916020019162001697565b820191906000526020600020905b8154815290600101906020018083116200167957829003601f168201915b5050505050838281518110620016b157620016b1620052a2565b6020026020010181905250601860008a81526020019081526020016000208054620016dc906200511c565b80601f01602080910402602001604051908101604052809291908181526020018280546200170a906200511c565b80156200175b5780601f106200172f576101008083540402835291602001916200175b565b820191906000526020600020905b8154815290600101906020018083116200173d57829003601f168201915b5050505050828281518110620017755762001775620052a2565b60200260200101819052508760a001516bffffffffffffffffffffffff1687620017a0919062004f5d565b60008a815260046020908152604080832080547fffffff00000000000000000000000000000000000000000000000000000000001681556001810180547fffffffff0000000000000000000000000000000000000000000000000000000016905560020180547fffffffffffffffffffffffff00000000000000000000000000000000000000001690556007909152812091985062001840919062003de5565b6000898152601760205260408120620018599162003de5565b6000898152601860205260408120620018729162003de5565b600089815260066020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055620018b360028a62003368565b5060a0880151604080516bffffffffffffffffffffffff909216825273ffffffffffffffffffffffffffffffffffffffff8c1660208301528a917fb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff910160405180910390a28062001924816200516c565b9150506200139c565b50856015546200193e919062005059565b60155560405160009062001963908d908d9087908a908a908990899060200162004b81565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905260105490915073ffffffffffffffffffffffffffffffffffffffff808c1691638e86139b916c010000000000000000000000009091041663c71249ab60028e73ffffffffffffffffffffffffffffffffffffffff1663aab9edd66040518163ffffffff1660e01b8152600401602060405180830381600087803b15801562001a1857600080fd5b505af115801562001a2d573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001a539190620049f9565b866040518463ffffffff1660e01b815260040162001a749392919062004e18565b60006040518083038186803b15801562001a8d57600080fd5b505afa15801562001aa2573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405262001aea91908101906200460b565b6040518263ffffffff1660e01b815260040162001b08919062004dba565b600060405180830381600087803b15801562001b2357600080fd5b505af115801562001b38573d6000803e3d6000fd5b50506040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8d81166004830152602482018b90527f000000000000000000000000000000000000000000000000000000000000000016925063a9059cbb9150604401602060405180830381600087803b15801562001bcd57600080fd5b505af115801562001be2573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001c08919062004534565b50505050505050505050505050565b60023360009081526016602052604090205460ff16600381111562001c405762001c4062005244565b1415801562001c76575060033360009081526016602052604090205460ff16600381111562001c735762001c7362005244565b14155b1562001cae576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000808080808062001cc3878901896200442c565b95509550955095509550955060005b86518110156200200c57600073ffffffffffffffffffffffffffffffffffffffff1686828151811062001d095762001d09620052a2565b60200260200101516060015173ffffffffffffffffffffffffffffffffffffffff16141562001e115786818151811062001d475762001d47620052a2565b602002602001015186828151811062001d645762001d64620052a2565b602002602001015160e0015160405162001d7e9062003dd7565b91825273ffffffffffffffffffffffffffffffffffffffff166020820152604001604051809103906000f08015801562001dbc573d6000803e3d6000fd5b5086828151811062001dd25762001dd2620052a2565b60200260200101516060019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250505b62001f5187828151811062001e2a5762001e2a620052a2565b602002602001015187838151811062001e475762001e47620052a2565b602002602001015160e0015188848151811062001e685762001e68620052a2565b60200260200101516000015187858151811062001e895762001e89620052a2565b60200260200101518a868151811062001ea65762001ea6620052a2565b602002602001015160a001518a878151811062001ec75762001ec7620052a2565b60200260200101518c888151811062001ee45762001ee4620052a2565b6020026020010151604001518a898151811062001f055762001f05620052a2565b60200260200101518a8a8151811062001f225762001f22620052a2565b60200260200101518f8b8151811062001f3f5762001f3f620052a2565b60200260200101516060015162002bdc565b86818151811062001f665762001f66620052a2565b60200260200101517f74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a7187838151811062001fa45762001fa4620052a2565b602002602001015160a001513360405162001fef9291906bffffffffffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a28062002003816200516c565b91505062001cd2565b505050505050505050565b600082815260046020908152604091829020825161010081018452815463ffffffff80821683526401000000008204811694830185905260ff6801000000000000000083041615159583019590955273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091048116606083015260018301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a084015278010000000000000000000000000000000000000000000000009004851660c083015260029092015490911660e082015291146200212a576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818160a001516200213c919062004fa0565b600084815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff93841602179055601554620021a49184169062004f5d565b6015556040517f23b872dd0000000000000000000000000000000000000000000000000000000081523360048201523060248201526bffffffffffffffffffffffff831660448201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906323b872dd90606401602060405180830381600087803b1580156200224957600080fd5b505af11580156200225e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062002284919062004534565b506040516bffffffffffffffffffffffff83168152339084907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a3505050565b6000818152600460209081526040808320815161010081018352815463ffffffff80821683526401000000008204811695830186905260ff6801000000000000000083041615159483019490945273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091048116606083015260018301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a084015278010000000000000000000000000000000000000000000000009004841660c083015260029092015490911660e08201529291141590620023cd60005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16149050818015620024285750808015620024265750620024196200337f565b836020015163ffffffff16115b155b1562002460576040517ffbc0357800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8015801562002493575060008481526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314155b15620024cb576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000620024d76200337f565b905081620024ef57620024ec60328262004f5d565b90505b6000858152600460205260409020805463ffffffff808416640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff909216919091179091556200254a9060029087906200336816565b5060105460808501516bffffffffffffffffffffffff9182169160009116821115620025b357608086015162002581908362005073565b90508560a001516bffffffffffffffffffffffff16816bffffffffffffffffffffffff161115620025b3575060a08501515b808660a00151620025c5919062005073565b600088815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff938416021790556011546200262d9183911662004fa0565b601180547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff9290921691909117905560405167ffffffffffffffff84169088907f91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f79118190600090a350505050505050565b620026b46200344c565b620026bf81620034cf565b50565b600060606000806000806200277b87600760008a81526020019081526020016000208054620026f1906200511c565b80601f01602080910402602001604051908101604052809291908181526020018280546200271f906200511c565b8015620027705780601f10620027445761010080835404028352916020019162002770565b820191906000526020600020905b8154815290600101906020018083116200275257829003601f168201915b505050505062000871565b949c939b5091995097509550909350915050565b33301415620028f3576000826003811115620027af57620027af62005244565b1415620027c757805115620027c357600080fd5b5050565b6001826003811115620027de57620027de62005244565b14156200284b57600081806020019051810190620027fd9190620047bc565b805190915073ffffffffffffffffffffffffffffffffffffffff166200282257600080fd5b60408101516200283157600080fd5b6008816020015160ff16106200284657600080fd5b505050565b600282600381111562002862576200286262005244565b1415620028a4576000818060200190518101906200288191906200470b565b90506020816020015151620028979190620051cf565b6004146200284657600080fd5b6003826003811115620028bb57620028bb62005244565b1415620028ee57600081806020019051810190620028da91906200484c565b8051519091506200289790602090620051cf565b600080fd5b6040517ffd9541d7000000000000000000000000000000000000000000000000000000008152309063fd9541d79062002933908590859060040162004df4565b600060405180830381600087803b1580156200294e57600080fd5b505af192505050801562002960575060015b620027c3576040517fa768d7fd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314620029f5576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600081815260046020526040902054640100000000900463ffffffff90811614620026bf576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080600062002a73600162002a616200337f565b62002a6d919062005059565b620035c7565b601254604080516020810193909352309083015268010000000000000000900463ffffffff166060820152608001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201209083015201604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052905060045b600f81101562002b73578282828151811062002b2f5762002b2f620052a2565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508062002b6a816200516c565b91505062002b0f565b5083600381111562002b895762002b8962005244565b60f81b81600f8151811062002ba25762002ba2620052a2565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350620005e881620050a3565b600f546e010000000000000000000000000000900460ff161562002c2c576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff89163b62002c7b576040517f09ee12d500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60125485517401000000000000000000000000000000000000000090910463ffffffff16101562002cd8576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108fc8863ffffffff16108062002d02575060125463ffffffff6401000000009091048116908916115b1562002d3a576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008a81526004602052604090206002015473ffffffffffffffffffffffffffffffffffffffff161562002d9a576040517f6e3b930b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040518061010001604052808963ffffffff16815260200163ffffffff8016815260200185151581526020018273ffffffffffffffffffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff168152602001876bffffffffffffffffffffffff168152602001600063ffffffff1681526020018a73ffffffffffffffffffffffffffffffffffffffff16815250600460008c815260200190815260200160002060008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160086101000a81548160ff02191690831515021790555060608201518160000160096101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060808201518160010160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060a082015181600101600c6101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060c08201518160010160186101000a81548163ffffffff021916908363ffffffff16021790555060e08201518160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505086600560008c815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550856bffffffffffffffffffffffff1660155462003049919062004f5d565b60155560008a815260076020908152604090912086516200306d9288019062003e24565b5060008a815260176020908152604090912084516200308f9286019062003e24565b5060008a81526018602090815260409091208351620030b19285019062003e24565b50620030bf60028b62003757565b5050505050505050505050565b321562003105576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b6000806000836060015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b1580156200318f57600080fd5b505afa158015620031a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620031ca9190620049a4565b5094509092505050600081131580620031e257508142105b80620032075750828015620032075750620031fe824262005059565b8463ffffffff16105b15620032185760135495506200321c565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b1580156200328357600080fd5b505afa15801562003298573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620032be9190620049a4565b5094509092505050600081131580620032d657508142105b80620032fb5750828015620032fb5750620032f2824262005059565b8463ffffffff16105b156200330c57601454945062003310565b8094505b50505050915091565b6000806200332c86896000015162003765565b9050600080620033498a8a63ffffffff16858a8a60018b620037b3565b90925090506200335a818362004fa0565b9a9950505050505050505050565b600062003376838362003bd4565b90505b92915050565b600060017f00000000000000000000000000000000000000000000000000000000000000006002811115620033b857620033b862005244565b14156200344757606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b1580156200340757600080fd5b505afa1580156200341c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620034429190620045ab565b905090565b504390565b60005473ffffffffffffffffffffffffffffffffffffffff16331462003105576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000f7b565b73ffffffffffffffffffffffffffffffffffffffff811633141562003551576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000f7b565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060017f0000000000000000000000000000000000000000000000000000000000000000600281111562003600576200360062005244565b14156200374d576000606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b1580156200365157600080fd5b505afa15801562003666573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200368c9190620045ab565b90508083101580620036aa5750610100620036a8848362005059565b115b15620036b95750600092915050565b6040517f2b407a8200000000000000000000000000000000000000000000000000000000815260048101849052606490632b407a829060240160206040518083038186803b1580156200370b57600080fd5b505afa15801562003720573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620037469190620045ab565b9392505050565b504090565b919050565b600062003376838362003cd8565b60006200377a63ffffffff8416601462004fea565b6200378783600162004f78565b620037989060ff16611d4c62004fea565b620037a790620124f862004f5d565b62003376919062004f5d565b6000806000896080015161ffff1687620037ce919062004fea565b9050838015620037dd5750803a105b15620037e657503a5b600060027f000000000000000000000000000000000000000000000000000000000000000060028111156200381f576200381f62005244565b1415620039b55760408051600081526020810190915285156200388457600036604051806080016040528060488152602001620058bb604891396040516020016200386d9392919062004b3a565b604051602081830303815290604052905062003906565b601254620038b6907801000000000000000000000000000000000000000000000000900463ffffffff1660046200502a565b63ffffffff1667ffffffffffffffff811115620038d757620038d7620052d1565b6040519080825280601f01601f19166020018201604052801562003902576020820181803683370190505b5090505b6040517f49948e0e00000000000000000000000000000000000000000000000000000000815273420000000000000000000000000000000000000f906349948e0e906200395890849060040162004dba565b60206040518083038186803b1580156200397157600080fd5b505afa15801562003986573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620039ac9190620045ab565b91505062003a79565b60017f00000000000000000000000000000000000000000000000000000000000000006002811115620039ec57620039ec62005244565b141562003a7957606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b815260040160206040518083038186803b15801562003a3b57600080fd5b505afa15801562003a50573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003a769190620045ab565b90505b8462003a9857808b6080015161ffff1662003a95919062004fea565b90505b62003aa861ffff87168262004fd3565b90506000878262003aba8c8e62004f5d565b62003ac6908662004fea565b62003ad2919062004f5d565b62003ae690670de0b6b3a764000062004fea565b62003af2919062004fd3565b905060008c6040015163ffffffff1664e8d4a5100062003b13919062004fea565b898e6020015163ffffffff16858f8862003b2e919062004fea565b62003b3a919062004f5d565b62003b4a90633b9aca0062004fea565b62003b56919062004fea565b62003b62919062004fd3565b62003b6e919062004f5d565b90506b033b2e3c9fd0803ce800000062003b89828462004f5d565b111562003bc2576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b6000818152600183016020526040812054801562003ccd57600062003bfb60018362005059565b855490915060009062003c119060019062005059565b905081811462003c7d57600086600001828154811062003c355762003c35620052a2565b906000526020600020015490508087600001848154811062003c5b5762003c5b620052a2565b6000918252602080832090910192909255918252600188019052604090208390555b855486908062003c915762003c9162005273565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505062003379565b600091505062003379565b600081815260018301602052604081205462003d215750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915562003379565b50600062003379565b82805462003d38906200511c565b90600052602060002090601f01602090048101928262003d5c576000855562003dc5565b82601f1062003d95578280017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082351617855562003dc5565b8280016001018555821562003dc5579182015b8281111562003dc557823582559160200191906001019062003da8565b5062003dd392915062003ea1565b5090565b610575806200534683390190565b50805462003df3906200511c565b6000825580601f1062003e04575050565b601f016020900490600052602060002090810190620026bf919062003ea1565b82805462003e32906200511c565b90600052602060002090601f01602090048101928262003e56576000855562003dc5565b82601f1062003e7157805160ff191683800117855562003dc5565b8280016001018555821562003dc5579182015b8281111562003dc557825182559160200191906001019062003e84565b5b8082111562003dd3576000815560010162003ea2565b600062003ecf62003ec98462004f14565b62004e9b565b905082815283838301111562003ee457600080fd5b62003746836020830184620050e9565b8035620037528162005313565b600082601f83011262003f1357600080fd5b8135602062003f2662003ec98362004eed565b80838252828201915082860187848660051b890101111562003f4757600080fd5b60005b8581101562003f7357813562003f608162005313565b8452928401929084019060010162003f4a565b5090979650505050505050565b600082601f83011262003f9257600080fd5b8135602062003fa562003ec98362004eed565b80838252828201915082860187848660051b890101111562003fc657600080fd5b6000805b868110156200400e57823567ffffffffffffffff81111562003fea578283fd5b62003ffa8b88838d0101620041e4565b865250938501939185019160010162003fca565b509198975050505050505050565b600082601f8301126200402e57600080fd5b813560206200404162003ec98362004eed565b80838252828201915082860187848660081b89010111156200406257600080fd5b6000805b868110156200400e5761010080848c03121562004081578283fd5b6200408b62004e48565b62004096856200426d565b8152620040a58886016200426d565b88820152604080860135620040ba8162005336565b908201526060620040cd86820162003ef4565b908201526080620040e0868201620042af565b9082015260a0620040f3868201620042af565b9082015260c0620041068682016200426d565b9082015260e06200411986820162003ef4565b90820152865294860194929092019160010162004066565b600082601f8301126200414357600080fd5b813560206200415662003ec98362004eed565b80838252828201915082860187848660051b89010111156200417757600080fd5b60005b8581101562003f73578135845292840192908401906001016200417a565b60008083601f840112620041ab57600080fd5b50813567ffffffffffffffff811115620041c457600080fd5b602083019150836020828501011115620041dd57600080fd5b9250929050565b600082601f830112620041f657600080fd5b81356200420762003ec98262004f14565b8181528460208386010111156200421d57600080fd5b816020850160208301376000918101602001919091529392505050565b600082601f8301126200424c57600080fd5b620033768383516020850162003eb8565b8035600481106200375257600080fd5b803563ffffffff811681146200375257600080fd5b805169ffffffffffffffffffff811681146200375257600080fd5b805160ff811681146200375257600080fd5b80356bffffffffffffffffffffffff811681146200375257600080fd5b600060208284031215620042df57600080fd5b8135620037468162005313565b600080600080600080600060a0888a0312156200430857600080fd5b8735620043158162005313565b965062004325602089016200426d565b95506040880135620043378162005313565b9450606088013567ffffffffffffffff808211156200435557600080fd5b620043638b838c0162004198565b909650945060808a01359150808211156200437d57600080fd5b506200438c8a828b0162004198565b989b979a50959850939692959293505050565b600080600060408486031215620043b557600080fd5b833567ffffffffffffffff80821115620043ce57600080fd5b818601915086601f830112620043e357600080fd5b813581811115620043f357600080fd5b8760208260051b85010111156200440957600080fd5b60209283019550935050840135620044218162005313565b809150509250925092565b60008060008060008060c087890312156200444657600080fd5b863567ffffffffffffffff808211156200445f57600080fd5b6200446d8a838b0162004131565b975060208901359150808211156200448457600080fd5b620044928a838b016200401c565b96506040890135915080821115620044a957600080fd5b620044b78a838b0162003f80565b95506060890135915080821115620044ce57600080fd5b620044dc8a838b0162003f01565b94506080890135915080821115620044f357600080fd5b620045018a838b0162003f80565b935060a08901359150808211156200451857600080fd5b506200452789828a0162003f80565b9150509295509295509295565b6000602082840312156200454757600080fd5b8151620037468162005336565b600080604083850312156200456857600080fd5b8251620045758162005336565b602084015190925067ffffffffffffffff8111156200459357600080fd5b620045a1858286016200423a565b9150509250929050565b600060208284031215620045be57600080fd5b5051919050565b60008060208385031215620045d957600080fd5b823567ffffffffffffffff811115620045f157600080fd5b620045ff8582860162004198565b90969095509350505050565b6000602082840312156200461e57600080fd5b815167ffffffffffffffff8111156200463657600080fd5b620005e8848285016200423a565b600080604083850312156200465857600080fd5b62004663836200425d565b9150602083013567ffffffffffffffff8111156200468057600080fd5b620045a185828601620041e4565b600080600060608486031215620046a457600080fd5b620046af846200425d565b9250602084013567ffffffffffffffff80821115620046cd57600080fd5b620046db87838801620041e4565b93506040860135915080821115620046f257600080fd5b506200470186828701620041e4565b9150509250925092565b6000602082840312156200471e57600080fd5b815167ffffffffffffffff808211156200473757600080fd5b90830190604082860312156200474c57600080fd5b6200475662004e75565b8251828111156200476657600080fd5b8301601f810187136200477857600080fd5b620047898782516020840162003eb8565b8252506020830151828111156200479f57600080fd5b620047ad878286016200423a565b60208301525095945050505050565b600060c08284031215620047cf57600080fd5b60405160c0810181811067ffffffffffffffff82111715620047f557620047f5620052d1565b6040528251620048058162005313565b815262004815602084016200429d565b602082015260408301516040820152606083015160608201526080830151608082015260a083015160a08201528091505092915050565b6000602082840312156200485f57600080fd5b815167ffffffffffffffff808211156200487857600080fd5b90830190602082860312156200488d57600080fd5b604051602081018181108382111715620048ab57620048ab620052d1565b604052825182811115620048be57600080fd5b620048cc878286016200423a565b82525095945050505050565b600060208284031215620048eb57600080fd5b5035919050565b6000806000604084860312156200490857600080fd5b83359250602084013567ffffffffffffffff8111156200492757600080fd5b620049358682870162004198565b9497909650939450505050565b600080604083850312156200495657600080fd5b82359150602083013567ffffffffffffffff8111156200468057600080fd5b600080604083850312156200498957600080fd5b823591506200499b60208401620042af565b90509250929050565b600080600080600060a08688031215620049bd57600080fd5b620049c88662004282565b9450602086015193506040860151925060608601519150620049ed6080870162004282565b90509295509295909350565b60006020828403121562004a0c57600080fd5b62003376826200429d565b600081518084526020808501945080840160005b8381101562004a5f57815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010162004a2b565b509495945050505050565b600082825180855260208086019550808260051b84010181860160005b8481101562003f73577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe086840301895262004ac483835162004ad7565b9884019892509083019060010162004a87565b6000815180845262004af1816020860160208601620050e9565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6004811062004b365762004b3662005244565b9052565b82848237600083820160008152835162004b59818360208801620050e9565b0195945050505050565b6000825162004b77818460208701620050e9565b9190910192915050565b600060c0808352888184015260e07f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8a111562004bbd57600080fd5b8960051b808c838701378085019050818101600081526020838784030181880152818c5180845261010093508385019150828e01945060005b8181101562004cbb578551805163ffffffff908116855285820151168585015260408082015115159085015260608082015173ffffffffffffffffffffffffffffffffffffffff81168287015250506080818101516bffffffffffffffffffffffff811686830152505060a0818101516bffffffffffffffffffffffff81168683015250508089015163ffffffff8116858b01525087015173ffffffffffffffffffffffffffffffffffffffff81168489015250948301949184019160010162004bf6565b5050878103604089015262004cd1818d62004a6a565b95505050505050828103606084015262004cec818762004a17565b9050828103608084015262004d02818662004a6a565b905082810360a08401526200335a818562004a6a565b861515815260c06020820152600062004d3560c083018862004ad7565b90506007861062004d4a5762004d4a62005244565b8560408301528460608301528360808301528260a0830152979650505050505050565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b60208152600062003376602083018462004ad7565b6020810162004dde8362005300565b91905290565b6020810162003379828462004b23565b62004e00818462004b23565b604060208201526000620005e8604083018462004ad7565b60ff8416815260ff8316602082015260606040820152600062004e3f606083018462004ad7565b95945050505050565b604051610100810167ffffffffffffffff8111828210171562004e6f5762004e6f620052d1565b60405290565b6040805190810167ffffffffffffffff8111828210171562004e6f5762004e6f620052d1565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171562004ee55762004ee5620052d1565b604052919050565b600067ffffffffffffffff82111562004f0a5762004f0a620052d1565b5060051b60200190565b600067ffffffffffffffff82111562004f315762004f31620052d1565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b6000821982111562004f735762004f73620051e6565b500190565b600060ff821660ff84168060ff0382111562004f985762004f98620051e6565b019392505050565b60006bffffffffffffffffffffffff80831681851680830382111562004fca5762004fca620051e6565b01949350505050565b60008262004fe55762004fe562005215565b500490565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615620050255762005025620051e6565b500290565b600063ffffffff80831681851681830481118215151615620050505762005050620051e6565b02949350505050565b6000828210156200506e576200506e620051e6565b500390565b60006bffffffffffffffffffffffff838116908316818110156200509b576200509b620051e6565b039392505050565b80516020808301519190811015620050e3577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8160200360031b1b821691505b50919050565b60005b8381101562005106578181015183820152602001620050ec565b8381111562005116576000848401525b50505050565b600181811c908216806200513157607f821691505b60208210811415620050e3577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415620051a157620051a1620051e6565b5060010190565b600063ffffffff80831681811415620051c557620051c5620051e6565b6001019392505050565b600082620051e157620051e162005215565b500690565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60038110620026bf57620026bf62005244565b73ffffffffffffffffffffffffffffffffffffffff81168114620026bf57600080fd5b8015158114620026bf57600080fdfe60c060405234801561001057600080fd5b5060405161057538038061057583398101604081905261002f9161005b565b60008054336001600160a01b031990911617905560601b6001600160601b03191660805260a052610098565b6000806040838503121561006e57600080fd5b825160208401519092506001600160a01b038116811461008d57600080fd5b809150509250929050565b60805160601c60a0516104ae6100c76000396000610145015260008181610170015261028001526104ae6000f3fe608060405234801561001057600080fd5b50600436106100725760003560e01c806379188d161161005057806379188d161461011d5780638ee489b214610140578063f00e6a2a1461016e57600080fd5b8063181f5a77146100775780631a5da6c8146100c95780635ab1bd53146100de575b600080fd5b6100b36040518060400160405280601981526020017f4175746f6d6174696f6e466f7277617264657220312e302e300000000000000081525081565b6040516100c091906103ff565b60405180910390f35b6100dc6100d73660046102e9565b610194565b005b60005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100c0565b61013061012b366004610326565b61022c565b60405190151581526020016100c0565b6040517f000000000000000000000000000000000000000000000000000000000000000081526020016100c0565b7f00000000000000000000000000000000000000000000000000000000000000006100f8565b60005473ffffffffffffffffffffffffffffffffffffffff1633146101e5576040517fea8e4eb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b6000805473ffffffffffffffffffffffffffffffffffffffff16331461027e576040517fea8e4eb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f00000000000000000000000000000000000000000000000000000000000000005a6113888110156102af57600080fd5b6113888103905084604082048203116102c757600080fd5b50803b6102d357600080fd5b60008084516020860160008589f1949350505050565b6000602082840312156102fb57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461031f57600080fd5b9392505050565b6000806040838503121561033957600080fd5b82359150602083013567ffffffffffffffff8082111561035857600080fd5b818501915085601f83011261036c57600080fd5b81358181111561037e5761037e610472565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019083821181831017156103c4576103c4610472565b816040528281528860208487010111156103dd57600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b600060208083528351808285015260005b8181101561042c57858101830151858201604001528201610410565b8181111561043e576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000806000a",
>>>>>>> automation-log-trigger-support
=======
	ABI: "[{\"inputs\":[{\"internalType\":\"contractKeeperRegistryLogicB2_1\",\"name\":\"logicB\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTrigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"trigger\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepTriggerConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumUpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumUpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fallbackTo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFastGasFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkNativeFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMode\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase2_1.Mode\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTriggerType\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase2_1.Trigger\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_next\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"migrateUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedUpkeeps\",\"type\":\"bytes\"}],\"name\":\"receiveUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"setUpkeepTriggerConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTranscoderVersion\",\"outputs\":[{\"internalType\":\"enumUpkeepFormat\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepVersion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumKeeperRegistryBase2_1.Trigger\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"}],\"name\":\"validateTriggerConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101206040523480156200001257600080fd5b5060405162005dcc38038062005dcc8339810160408190526200003591620003aa565b80816001600160a01b0316634b4fd03b6040518163ffffffff1660e01b815260040160206040518083038186803b1580156200007057600080fd5b505afa15801562000085573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620000ab9190620003d1565b826001600160a01b031663ca30e6036040518163ffffffff1660e01b815260040160206040518083038186803b158015620000e557600080fd5b505afa158015620000fa573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001209190620003aa565b836001600160a01b031663b10b673c6040518163ffffffff1660e01b815260040160206040518083038186803b1580156200015a57600080fd5b505afa1580156200016f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001959190620003aa565b846001600160a01b0316636709d0e56040518163ffffffff1660e01b815260040160206040518083038186803b158015620001cf57600080fd5b505afa158015620001e4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200020a9190620003aa565b3380600081620002615760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b038481169190911790915581161562000294576200029481620002fe565b505050836002811115620002ac57620002ac620003f4565b60e0816002811115620002c357620002c3620003f4565b60f81b9052506001600160601b0319606093841b811660805291831b821660a052821b811660c05292901b9091166101005250620004239050565b6001600160a01b038116331415620003595760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000258565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215620003bd57600080fd5b8151620003ca816200040a565b9392505050565b600060208284031215620003e457600080fd5b815160038110620003ca57600080fd5b634e487b7160e01b600052602160045260246000fd5b6001600160a01b03811681146200042057600080fd5b50565b60805160601c60a05160601c60c05160601c60e05160f81c6101005160601c61590f620004bd6000396000818161019e015281816101e4015261028c01526000818161024801528181613385015281816135cd015281816137ec01526139b90152600081816102d6015261312a0152600081816103d8015261321e01526000818161041601528181611b8601526121ef015261590f6000f3fe60806040523480156200001157600080fd5b50600436106200019c5760003560e01c806385c1b0ba11620000ed578063b10b673c1162000099578063f2fde38b116200006f578063f2fde38b146200043b578063f7d334ba1462000452578063fd9541d71462000469576200019c565b8063b10b673c14620003d6578063c804802214620003fd578063ca30e6031462000414576200019c565b80638e86139b11620000cf5780638e86139b146200038c578063948108f714620003a3578063aab9edd614620003ba576200019c565b806385c1b0ba14620003565780638da5cb5b146200036d576200019c565b80635147cd59116200014d5780636ded9eae116200012f5780636ded9eae14620002fb57806371791aa0146200032157806379ba5097146200034c576200019c565b80635147cd5914620002ae5780636709d0e514620002d4576200019c565b80634b4fd03b11620001835780634b4fd03b14620002465780634ee88d35146200026d5780634ff597c11462000286576200019c565b8063349e8cca14620001e257806348013d7b146200022e575b7f00000000000000000000000000000000000000000000000000000000000000003660008037600080366000845af43d6000803e808015620001dd573d6000f35b3d6000fd5b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b62000237600081565b60405162000225919062004dcf565b7f000000000000000000000000000000000000000000000000000000000000000062000237565b620002846200027e366004620048f2565b62000480565b005b620002047f000000000000000000000000000000000000000000000000000000000000000081565b620002c5620002bf366004620048d8565b6200053a565b60405162000225919062004de4565b7f000000000000000000000000000000000000000000000000000000000000000062000204565b620003126200030c366004620042ec565b620005f0565b60405190815260200162000225565b620003386200033236600462004942565b62000871565b604051620002259695949392919062004d18565b6200028462000efd565b62000284620003673660046200439f565b62001000565b60005473ffffffffffffffffffffffffffffffffffffffff1662000204565b620002846200039d366004620045c5565b62001c17565b62000284620003b436600462004975565b62002017565b620003c3600281565b60405160ff909116815260200162000225565b7f000000000000000000000000000000000000000000000000000000000000000062000204565b620002846200040e366004620048d8565b620022cf565b7f000000000000000000000000000000000000000000000000000000000000000062000204565b620002846200044c366004620042cc565b620026aa565b6200033862000463366004620048d8565b620026c2565b620002846200047a36600462004644565b6200278f565b6200048b8362002997565b600062000498846200053a565b9050620004dc8184848080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506200278f92505050565b6000848152601760205260409020620004f790848462003d2a565b50837f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d566484846040516200052c92919062004d6d565b60405180910390a250505050565b6000818160045b600f811015620005cf577fff000000000000000000000000000000000000000000000000000000000000008216838260208110620005835762000583620052a2565b1a60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614620005ba57506000949350505050565b80620005c6816200516c565b91505062000541565b5081600f1a6003811115620005e857620005e862005244565b949350505050565b6000805473ffffffffffffffffffffffffffffffffffffffff1633148015906200064257506011546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff163314155b156200067a576040517fd48b678b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080806200068c858701876200468e565b9250925092506200069e83836200278f565b620006a98362002a4c565b93506000848c604051620006bd9062003dd7565b91825273ffffffffffffffffffffffffffffffffffffffff166020820152604001604051809103906000f080158015620006fb573d6000803e3d6000fd5b50905062000749858d8d8d60008e8e8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092018290525092508c91508b90508a62002bdc565b6012805468010000000000000000900463ffffffff169060086200076d83620051a8565b91906101000a81548163ffffffff021916908363ffffffff16021790555050847fbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d0128c8c604051620007e692919063ffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a2847f2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d56648460405162000820919062004dba565b60405180910390a2847f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850836040516200085a919062004dba565b60405180910390a250505050979650505050505050565b6000606060008060008062000885620030cc565b600062000892896200053a565b90506000600f604051806101200160405290816000820160009054906101000a900460ff1660ff1660ff1681526020016000820160019054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160059054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160099054906101000a900462ffffff1662ffffff1662ffffff16815260200160008201600c9054906101000a900461ffff1661ffff1661ffff16815260200160008201600e9054906101000a900460ff1615151515815260200160008201600f9054906101000a900460ff161515151581526020016000820160109054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160008201601c9054906101000a900463ffffffff1663ffffffff1663ffffffff168152505090506000600460008c8152602001908152602001600020604051806101000160405290816000820160009054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160049054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160089054906101000a900460ff161515151581526020016000820160099054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160018201600c9054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff1681526020016001820160189054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016002820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681525050905063ffffffff8016816020015163ffffffff161462000c0457505060408051602081019091526000808252975095506001945086925082915062000ef39050565b80604001511562000c3757505060408051602081019091526000808252975095506002945086925082915062000ef39050565b62000c428262003107565b825160125492975090955060009162000c829185917801000000000000000000000000000000000000000000000000900463ffffffff1689898662003319565b9050806bffffffffffffffffffffffff168260a001516bffffffffffffffffffffffff16101562000cd15760006040518060200160405280600081525060069950995099505050505062000ef3565b60019950600084600381111562000cec5762000cec62005244565b148062000d0d5750600184600381111562000d0b5762000d0b62005244565b145b1562000eee575a96506000636e04ff0d60e01b8c60405160240162000d33919062004dba565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009094169390931790925260e0850151601254925191935073ffffffffffffffffffffffffffffffffffffffff169163ffffffff169062000de890849062004b63565b60006040518083038160008787f1925050503d806000811462000e28576040519150601f19603f3d011682016040523d82523d6000602084013e62000e2d565b606091505b50909b5099505a62000e40908962005059565b97508a62000e52576003985062000eec565b8980602001905181019062000e68919062004554565b909b5099508a62000e985760006040518060200160405280600081525060049a509a509a50505050505062000ef3565b6012548a51780100000000000000000000000000000000000000000000000090910463ffffffff16101562000eec5760006040518060200160405280600081525060059a509a509a50505050505062000ef3565b505b505050505b9295509295509295565b60015473ffffffffffffffffffffffffffffffffffffffff16331462000f84576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600173ffffffffffffffffffffffffffffffffffffffff821660009081526016602052604090205460ff1660038111156200103f576200103f62005244565b141580156200108b5750600373ffffffffffffffffffffffffffffffffffffffff821660009081526016602052604090205460ff16600381111562001088576200108862005244565b14155b15620010c3576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6010546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1662001123576040517fd12d7d8d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b816200115b576040517f2c2fc94100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805161010081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526000808567ffffffffffffffff811115620011ba57620011ba620052d1565b604051908082528060200260200182016040528015620011ef57816020015b6060815260200190600190039081620011d95790505b50905060008667ffffffffffffffff811115620012105762001210620052d1565b6040519080825280602002602001820160405280156200123a578160200160208202803683370190505b50905060008767ffffffffffffffff8111156200125b576200125b620052d1565b604051908082528060200260200182016040528015620012ea57816020015b604080516101008101825260008082526020808301829052928201819052606082018190526080820181905260a0820181905260c0820181905260e082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816200127a5790505b50905060008867ffffffffffffffff8111156200130b576200130b620052d1565b6040519080825280602002602001820160405280156200134057816020015b60608152602001906001900390816200132a5790505b50905060008967ffffffffffffffff811115620013615762001361620052d1565b6040519080825280602002602001820160405280156200139657816020015b6060815260200190600190039081620013805790505b50905060005b8a8110156200192d578b8b82818110620013ba57620013ba620052a2565b60209081029290920135600081815260048452604090819020815161010081018352815463ffffffff8082168352640100000000820481169783019790975260ff6801000000000000000082041615159382019390935273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009093048316606082015260018201546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490951660c0860152600201541660e08401529a50909850620014ab90508962002997565b87848281518110620014c157620014c1620052a2565b6020026020010181905250600760008a81526020019081526020016000208054620014ec906200511c565b80601f01602080910402602001604051908101604052809291908181526020018280546200151a906200511c565b80156200156b5780601f106200153f576101008083540402835291602001916200156b565b820191906000526020600020905b8154815290600101906020018083116200154d57829003601f168201915b5050505050868281518110620015855762001585620052a2565b6020026020010181905250600560008a815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16858281518110620015d957620015d9620052a2565b73ffffffffffffffffffffffffffffffffffffffff90921660209283029190910182015260008a8152601790915260409020805462001618906200511c565b80601f016020809104026020016040519081016040528092919081815260200182805462001646906200511c565b8015620016975780601f106200166b5761010080835404028352916020019162001697565b820191906000526020600020905b8154815290600101906020018083116200167957829003601f168201915b5050505050838281518110620016b157620016b1620052a2565b6020026020010181905250601860008a81526020019081526020016000208054620016dc906200511c565b80601f01602080910402602001604051908101604052809291908181526020018280546200170a906200511c565b80156200175b5780601f106200172f576101008083540402835291602001916200175b565b820191906000526020600020905b8154815290600101906020018083116200173d57829003601f168201915b5050505050828281518110620017755762001775620052a2565b60200260200101819052508760a001516bffffffffffffffffffffffff1687620017a0919062004f5d565b60008a815260046020908152604080832080547fffffff00000000000000000000000000000000000000000000000000000000001681556001810180547fffffffff0000000000000000000000000000000000000000000000000000000016905560020180547fffffffffffffffffffffffff00000000000000000000000000000000000000001690556007909152812091985062001840919062003de5565b6000898152601760205260408120620018599162003de5565b6000898152601860205260408120620018729162003de5565b600089815260066020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055620018b360028a62003368565b5060a0880151604080516bffffffffffffffffffffffff909216825273ffffffffffffffffffffffffffffffffffffffff8c1660208301528a917fb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff910160405180910390a28062001924816200516c565b9150506200139c565b50856015546200193e919062005059565b60155560405160009062001963908d908d9087908a908a908990899060200162004b81565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905260105490915073ffffffffffffffffffffffffffffffffffffffff808c1691638e86139b916c010000000000000000000000009091041663c71249ab60028e73ffffffffffffffffffffffffffffffffffffffff1663aab9edd66040518163ffffffff1660e01b8152600401602060405180830381600087803b15801562001a1857600080fd5b505af115801562001a2d573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001a539190620049f9565b866040518463ffffffff1660e01b815260040162001a749392919062004e18565b60006040518083038186803b15801562001a8d57600080fd5b505afa15801562001aa2573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405262001aea91908101906200460b565b6040518263ffffffff1660e01b815260040162001b08919062004dba565b600060405180830381600087803b15801562001b2357600080fd5b505af115801562001b38573d6000803e3d6000fd5b50506040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8d81166004830152602482018b90527f000000000000000000000000000000000000000000000000000000000000000016925063a9059cbb9150604401602060405180830381600087803b15801562001bcd57600080fd5b505af115801562001be2573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001c08919062004534565b50505050505050505050505050565b60023360009081526016602052604090205460ff16600381111562001c405762001c4062005244565b1415801562001c76575060033360009081526016602052604090205460ff16600381111562001c735762001c7362005244565b14155b1562001cae576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000808080808062001cc3878901896200442c565b95509550955095509550955060005b86518110156200200c57600073ffffffffffffffffffffffffffffffffffffffff1686828151811062001d095762001d09620052a2565b60200260200101516060015173ffffffffffffffffffffffffffffffffffffffff16141562001e115786818151811062001d475762001d47620052a2565b602002602001015186828151811062001d645762001d64620052a2565b602002602001015160e0015160405162001d7e9062003dd7565b91825273ffffffffffffffffffffffffffffffffffffffff166020820152604001604051809103906000f08015801562001dbc573d6000803e3d6000fd5b5086828151811062001dd25762001dd2620052a2565b60200260200101516060019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250505b62001f5187828151811062001e2a5762001e2a620052a2565b602002602001015187838151811062001e475762001e47620052a2565b602002602001015160e0015188848151811062001e685762001e68620052a2565b60200260200101516000015187858151811062001e895762001e89620052a2565b60200260200101518a868151811062001ea65762001ea6620052a2565b602002602001015160a001518a878151811062001ec75762001ec7620052a2565b60200260200101518c888151811062001ee45762001ee4620052a2565b6020026020010151604001518a898151811062001f055762001f05620052a2565b60200260200101518a8a8151811062001f225762001f22620052a2565b60200260200101518f8b8151811062001f3f5762001f3f620052a2565b60200260200101516060015162002bdc565b86818151811062001f665762001f66620052a2565b60200260200101517f74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a7187838151811062001fa45762001fa4620052a2565b602002602001015160a001513360405162001fef9291906bffffffffffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a28062002003816200516c565b91505062001cd2565b505050505050505050565b600082815260046020908152604091829020825161010081018452815463ffffffff80821683526401000000008204811694830185905260ff6801000000000000000083041615159583019590955273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091048116606083015260018301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a084015278010000000000000000000000000000000000000000000000009004851660c083015260029092015490911660e082015291146200212a576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818160a001516200213c919062004fa0565b600084815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff93841602179055601554620021a49184169062004f5d565b6015556040517f23b872dd0000000000000000000000000000000000000000000000000000000081523360048201523060248201526bffffffffffffffffffffffff831660448201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906323b872dd90606401602060405180830381600087803b1580156200224957600080fd5b505af11580156200225e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062002284919062004534565b506040516bffffffffffffffffffffffff83168152339084907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a3505050565b6000818152600460209081526040808320815161010081018352815463ffffffff80821683526401000000008204811695830186905260ff6801000000000000000083041615159483019490945273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091048116606083015260018301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a084015278010000000000000000000000000000000000000000000000009004841660c083015260029092015490911660e08201529291141590620023cd60005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16149050818015620024285750808015620024265750620024196200337f565b836020015163ffffffff16115b155b1562002460576040517ffbc0357800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8015801562002493575060008481526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314155b15620024cb576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000620024d76200337f565b905081620024ef57620024ec60328262004f5d565b90505b6000858152600460205260409020805463ffffffff808416640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff909216919091179091556200254a9060029087906200336816565b5060105460808501516bffffffffffffffffffffffff9182169160009116821115620025b357608086015162002581908362005073565b90508560a001516bffffffffffffffffffffffff16816bffffffffffffffffffffffff161115620025b3575060a08501515b808660a00151620025c5919062005073565b600088815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff938416021790556011546200262d9183911662004fa0565b601180547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff9290921691909117905560405167ffffffffffffffff84169088907f91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f79118190600090a350505050505050565b620026b46200344c565b620026bf81620034cf565b50565b600060606000806000806200277b87600760008a81526020019081526020016000208054620026f1906200511c565b80601f01602080910402602001604051908101604052809291908181526020018280546200271f906200511c565b8015620027705780601f10620027445761010080835404028352916020019162002770565b820191906000526020600020905b8154815290600101906020018083116200275257829003601f168201915b505050505062000871565b949c939b5091995097509550909350915050565b33301415620028f3576000826003811115620027af57620027af62005244565b1415620027c757805115620027c357600080fd5b5050565b6001826003811115620027de57620027de62005244565b14156200284b57600081806020019051810190620027fd9190620047bc565b805190915073ffffffffffffffffffffffffffffffffffffffff166200282257600080fd5b60408101516200283157600080fd5b6008816020015160ff16106200284657600080fd5b505050565b600282600381111562002862576200286262005244565b1415620028a4576000818060200190518101906200288191906200470b565b90506020816020015151620028979190620051cf565b6004146200284657600080fd5b6003826003811115620028bb57620028bb62005244565b1415620028ee57600081806020019051810190620028da91906200484c565b8051519091506200289790602090620051cf565b600080fd5b6040517ffd9541d7000000000000000000000000000000000000000000000000000000008152309063fd9541d79062002933908590859060040162004df4565b600060405180830381600087803b1580156200294e57600080fd5b505af192505050801562002960575060015b620027c3576040517fa768d7fd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008181526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314620029f5576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600081815260046020526040902054640100000000900463ffffffff90811614620026bf576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080600062002a73600162002a616200337f565b62002a6d919062005059565b620035c7565b601254604080516020810193909352309083015268010000000000000000900463ffffffff166060820152608001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201209083015201604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052905060045b600f81101562002b73578282828151811062002b2f5762002b2f620052a2565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508062002b6a816200516c565b91505062002b0f565b5083600381111562002b895762002b8962005244565b60f81b81600f8151811062002ba25762002ba2620052a2565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350620005e881620050a3565b600f546e010000000000000000000000000000900460ff161562002c2c576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff89163b62002c7b576040517f09ee12d500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60125485517401000000000000000000000000000000000000000090910463ffffffff16101562002cd8576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108fc8863ffffffff16108062002d02575060125463ffffffff6401000000009091048116908916115b1562002d3a576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008a81526004602052604090206002015473ffffffffffffffffffffffffffffffffffffffff161562002d9a576040517f6e3b930b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040518061010001604052808963ffffffff16815260200163ffffffff8016815260200185151581526020018273ffffffffffffffffffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff168152602001876bffffffffffffffffffffffff168152602001600063ffffffff1681526020018a73ffffffffffffffffffffffffffffffffffffffff16815250600460008c815260200190815260200160002060008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160086101000a81548160ff02191690831515021790555060608201518160000160096101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060808201518160010160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060a082015181600101600c6101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060c08201518160010160186101000a81548163ffffffff021916908363ffffffff16021790555060e08201518160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505086600560008c815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550856bffffffffffffffffffffffff1660155462003049919062004f5d565b60155560008a815260076020908152604090912086516200306d9288019062003e24565b5060008a815260176020908152604090912084516200308f9286019062003e24565b5060008a81526018602090815260409091208351620030b19285019062003e24565b50620030bf60028b62003757565b5050505050505050505050565b321562003105576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b6000806000836060015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b1580156200318f57600080fd5b505afa158015620031a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620031ca9190620049a4565b5094509092505050600081131580620031e257508142105b80620032075750828015620032075750620031fe824262005059565b8463ffffffff16105b15620032185760135495506200321c565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b1580156200328357600080fd5b505afa15801562003298573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620032be9190620049a4565b5094509092505050600081131580620032d657508142105b80620032fb5750828015620032fb5750620032f2824262005059565b8463ffffffff16105b156200330c57601454945062003310565b8094505b50505050915091565b6000806200332c86896000015162003765565b9050600080620033498a8a63ffffffff16858a8a60018b620037b3565b90925090506200335a818362004fa0565b9a9950505050505050505050565b600062003376838362003bd4565b90505b92915050565b600060017f00000000000000000000000000000000000000000000000000000000000000006002811115620033b857620033b862005244565b14156200344757606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b1580156200340757600080fd5b505afa1580156200341c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620034429190620045ab565b905090565b504390565b60005473ffffffffffffffffffffffffffffffffffffffff16331462003105576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000f7b565b73ffffffffffffffffffffffffffffffffffffffff811633141562003551576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000f7b565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060017f0000000000000000000000000000000000000000000000000000000000000000600281111562003600576200360062005244565b14156200374d576000606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b1580156200365157600080fd5b505afa15801562003666573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200368c9190620045ab565b90508083101580620036aa5750610100620036a8848362005059565b115b15620036b95750600092915050565b6040517f2b407a8200000000000000000000000000000000000000000000000000000000815260048101849052606490632b407a829060240160206040518083038186803b1580156200370b57600080fd5b505afa15801562003720573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620037469190620045ab565b9392505050565b504090565b919050565b600062003376838362003cd8565b60006200377a63ffffffff8416601462004fea565b6200378783600162004f78565b620037989060ff16611d4c62004fea565b620037a790620124f862004f5d565b62003376919062004f5d565b6000806000896080015161ffff1687620037ce919062004fea565b9050838015620037dd5750803a105b15620037e657503a5b600060027f000000000000000000000000000000000000000000000000000000000000000060028111156200381f576200381f62005244565b1415620039b55760408051600081526020810190915285156200388457600036604051806080016040528060488152602001620058bb604891396040516020016200386d9392919062004b3a565b604051602081830303815290604052905062003906565b601254620038b6907801000000000000000000000000000000000000000000000000900463ffffffff1660046200502a565b63ffffffff1667ffffffffffffffff811115620038d757620038d7620052d1565b6040519080825280601f01601f19166020018201604052801562003902576020820181803683370190505b5090505b6040517f49948e0e00000000000000000000000000000000000000000000000000000000815273420000000000000000000000000000000000000f906349948e0e906200395890849060040162004dba565b60206040518083038186803b1580156200397157600080fd5b505afa15801562003986573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620039ac9190620045ab565b91505062003a79565b60017f00000000000000000000000000000000000000000000000000000000000000006002811115620039ec57620039ec62005244565b141562003a7957606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b815260040160206040518083038186803b15801562003a3b57600080fd5b505afa15801562003a50573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003a769190620045ab565b90505b8462003a9857808b6080015161ffff1662003a95919062004fea565b90505b62003aa861ffff87168262004fd3565b90506000878262003aba8c8e62004f5d565b62003ac6908662004fea565b62003ad2919062004f5d565b62003ae690670de0b6b3a764000062004fea565b62003af2919062004fd3565b905060008c6040015163ffffffff1664e8d4a5100062003b13919062004fea565b898e6020015163ffffffff16858f8862003b2e919062004fea565b62003b3a919062004f5d565b62003b4a90633b9aca0062004fea565b62003b56919062004fea565b62003b62919062004fd3565b62003b6e919062004f5d565b90506b033b2e3c9fd0803ce800000062003b89828462004f5d565b111562003bc2576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b6000818152600183016020526040812054801562003ccd57600062003bfb60018362005059565b855490915060009062003c119060019062005059565b905081811462003c7d57600086600001828154811062003c355762003c35620052a2565b906000526020600020015490508087600001848154811062003c5b5762003c5b620052a2565b6000918252602080832090910192909255918252600188019052604090208390555b855486908062003c915762003c9162005273565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505062003379565b600091505062003379565b600081815260018301602052604081205462003d215750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915562003379565b50600062003379565b82805462003d38906200511c565b90600052602060002090601f01602090048101928262003d5c576000855562003dc5565b82601f1062003d95578280017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082351617855562003dc5565b8280016001018555821562003dc5579182015b8281111562003dc557823582559160200191906001019062003da8565b5062003dd392915062003ea1565b5090565b610575806200534683390190565b50805462003df3906200511c565b6000825580601f1062003e04575050565b601f016020900490600052602060002090810190620026bf919062003ea1565b82805462003e32906200511c565b90600052602060002090601f01602090048101928262003e56576000855562003dc5565b82601f1062003e7157805160ff191683800117855562003dc5565b8280016001018555821562003dc5579182015b8281111562003dc557825182559160200191906001019062003e84565b5b8082111562003dd3576000815560010162003ea2565b600062003ecf62003ec98462004f14565b62004e9b565b905082815283838301111562003ee457600080fd5b62003746836020830184620050e9565b8035620037528162005313565b600082601f83011262003f1357600080fd5b8135602062003f2662003ec98362004eed565b80838252828201915082860187848660051b890101111562003f4757600080fd5b60005b8581101562003f7357813562003f608162005313565b8452928401929084019060010162003f4a565b5090979650505050505050565b600082601f83011262003f9257600080fd5b8135602062003fa562003ec98362004eed565b80838252828201915082860187848660051b890101111562003fc657600080fd5b6000805b868110156200400e57823567ffffffffffffffff81111562003fea578283fd5b62003ffa8b88838d0101620041e4565b865250938501939185019160010162003fca565b509198975050505050505050565b600082601f8301126200402e57600080fd5b813560206200404162003ec98362004eed565b80838252828201915082860187848660081b89010111156200406257600080fd5b6000805b868110156200400e5761010080848c03121562004081578283fd5b6200408b62004e48565b62004096856200426d565b8152620040a58886016200426d565b88820152604080860135620040ba8162005336565b908201526060620040cd86820162003ef4565b908201526080620040e0868201620042af565b9082015260a0620040f3868201620042af565b9082015260c0620041068682016200426d565b9082015260e06200411986820162003ef4565b90820152865294860194929092019160010162004066565b600082601f8301126200414357600080fd5b813560206200415662003ec98362004eed565b80838252828201915082860187848660051b89010111156200417757600080fd5b60005b8581101562003f73578135845292840192908401906001016200417a565b60008083601f840112620041ab57600080fd5b50813567ffffffffffffffff811115620041c457600080fd5b602083019150836020828501011115620041dd57600080fd5b9250929050565b600082601f830112620041f657600080fd5b81356200420762003ec98262004f14565b8181528460208386010111156200421d57600080fd5b816020850160208301376000918101602001919091529392505050565b600082601f8301126200424c57600080fd5b620033768383516020850162003eb8565b8035600481106200375257600080fd5b803563ffffffff811681146200375257600080fd5b805169ffffffffffffffffffff811681146200375257600080fd5b805160ff811681146200375257600080fd5b80356bffffffffffffffffffffffff811681146200375257600080fd5b600060208284031215620042df57600080fd5b8135620037468162005313565b600080600080600080600060a0888a0312156200430857600080fd5b8735620043158162005313565b965062004325602089016200426d565b95506040880135620043378162005313565b9450606088013567ffffffffffffffff808211156200435557600080fd5b620043638b838c0162004198565b909650945060808a01359150808211156200437d57600080fd5b506200438c8a828b0162004198565b989b979a50959850939692959293505050565b600080600060408486031215620043b557600080fd5b833567ffffffffffffffff80821115620043ce57600080fd5b818601915086601f830112620043e357600080fd5b813581811115620043f357600080fd5b8760208260051b85010111156200440957600080fd5b60209283019550935050840135620044218162005313565b809150509250925092565b60008060008060008060c087890312156200444657600080fd5b863567ffffffffffffffff808211156200445f57600080fd5b6200446d8a838b0162004131565b975060208901359150808211156200448457600080fd5b620044928a838b016200401c565b96506040890135915080821115620044a957600080fd5b620044b78a838b0162003f80565b95506060890135915080821115620044ce57600080fd5b620044dc8a838b0162003f01565b94506080890135915080821115620044f357600080fd5b620045018a838b0162003f80565b935060a08901359150808211156200451857600080fd5b506200452789828a0162003f80565b9150509295509295509295565b6000602082840312156200454757600080fd5b8151620037468162005336565b600080604083850312156200456857600080fd5b8251620045758162005336565b602084015190925067ffffffffffffffff8111156200459357600080fd5b620045a1858286016200423a565b9150509250929050565b600060208284031215620045be57600080fd5b5051919050565b60008060208385031215620045d957600080fd5b823567ffffffffffffffff811115620045f157600080fd5b620045ff8582860162004198565b90969095509350505050565b6000602082840312156200461e57600080fd5b815167ffffffffffffffff8111156200463657600080fd5b620005e8848285016200423a565b600080604083850312156200465857600080fd5b62004663836200425d565b9150602083013567ffffffffffffffff8111156200468057600080fd5b620045a185828601620041e4565b600080600060608486031215620046a457600080fd5b620046af846200425d565b9250602084013567ffffffffffffffff80821115620046cd57600080fd5b620046db87838801620041e4565b93506040860135915080821115620046f257600080fd5b506200470186828701620041e4565b9150509250925092565b6000602082840312156200471e57600080fd5b815167ffffffffffffffff808211156200473757600080fd5b90830190604082860312156200474c57600080fd5b6200475662004e75565b8251828111156200476657600080fd5b8301601f810187136200477857600080fd5b620047898782516020840162003eb8565b8252506020830151828111156200479f57600080fd5b620047ad878286016200423a565b60208301525095945050505050565b600060c08284031215620047cf57600080fd5b60405160c0810181811067ffffffffffffffff82111715620047f557620047f5620052d1565b6040528251620048058162005313565b815262004815602084016200429d565b602082015260408301516040820152606083015160608201526080830151608082015260a083015160a08201528091505092915050565b6000602082840312156200485f57600080fd5b815167ffffffffffffffff808211156200487857600080fd5b90830190602082860312156200488d57600080fd5b604051602081018181108382111715620048ab57620048ab620052d1565b604052825182811115620048be57600080fd5b620048cc878286016200423a565b82525095945050505050565b600060208284031215620048eb57600080fd5b5035919050565b6000806000604084860312156200490857600080fd5b83359250602084013567ffffffffffffffff8111156200492757600080fd5b620049358682870162004198565b9497909650939450505050565b600080604083850312156200495657600080fd5b82359150602083013567ffffffffffffffff8111156200468057600080fd5b600080604083850312156200498957600080fd5b823591506200499b60208401620042af565b90509250929050565b600080600080600060a08688031215620049bd57600080fd5b620049c88662004282565b9450602086015193506040860151925060608601519150620049ed6080870162004282565b90509295509295909350565b60006020828403121562004a0c57600080fd5b62003376826200429d565b600081518084526020808501945080840160005b8381101562004a5f57815173ffffffffffffffffffffffffffffffffffffffff168752958201959082019060010162004a2b565b509495945050505050565b600082825180855260208086019550808260051b84010181860160005b8481101562003f73577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe086840301895262004ac483835162004ad7565b9884019892509083019060010162004a87565b6000815180845262004af1816020860160208601620050e9565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6004811062004b365762004b3662005244565b9052565b82848237600083820160008152835162004b59818360208801620050e9565b0195945050505050565b6000825162004b77818460208701620050e9565b9190910192915050565b600060c0808352888184015260e07f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8a111562004bbd57600080fd5b8960051b808c838701378085019050818101600081526020838784030181880152818c5180845261010093508385019150828e01945060005b8181101562004cbb578551805163ffffffff908116855285820151168585015260408082015115159085015260608082015173ffffffffffffffffffffffffffffffffffffffff81168287015250506080818101516bffffffffffffffffffffffff811686830152505060a0818101516bffffffffffffffffffffffff81168683015250508089015163ffffffff8116858b01525087015173ffffffffffffffffffffffffffffffffffffffff81168489015250948301949184019160010162004bf6565b5050878103604089015262004cd1818d62004a6a565b95505050505050828103606084015262004cec818762004a17565b9050828103608084015262004d02818662004a6a565b905082810360a08401526200335a818562004a6a565b861515815260c06020820152600062004d3560c083018862004ad7565b90506007861062004d4a5762004d4a62005244565b8560408301528460608301528360808301528260a0830152979650505050505050565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b60208152600062003376602083018462004ad7565b6020810162004dde8362005300565b91905290565b6020810162003379828462004b23565b62004e00818462004b23565b604060208201526000620005e8604083018462004ad7565b60ff8416815260ff8316602082015260606040820152600062004e3f606083018462004ad7565b95945050505050565b604051610100810167ffffffffffffffff8111828210171562004e6f5762004e6f620052d1565b60405290565b6040805190810167ffffffffffffffff8111828210171562004e6f5762004e6f620052d1565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171562004ee55762004ee5620052d1565b604052919050565b600067ffffffffffffffff82111562004f0a5762004f0a620052d1565b5060051b60200190565b600067ffffffffffffffff82111562004f315762004f31620052d1565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b6000821982111562004f735762004f73620051e6565b500190565b600060ff821660ff84168060ff0382111562004f985762004f98620051e6565b019392505050565b60006bffffffffffffffffffffffff80831681851680830382111562004fca5762004fca620051e6565b01949350505050565b60008262004fe55762004fe562005215565b500490565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615620050255762005025620051e6565b500290565b600063ffffffff80831681851681830481118215151615620050505762005050620051e6565b02949350505050565b6000828210156200506e576200506e620051e6565b500390565b60006bffffffffffffffffffffffff838116908316818110156200509b576200509b620051e6565b039392505050565b80516020808301519190811015620050e3577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8160200360031b1b821691505b50919050565b60005b8381101562005106578181015183820152602001620050ec565b8381111562005116576000848401525b50505050565b600181811c908216806200513157607f821691505b60208210811415620050e3577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415620051a157620051a1620051e6565b5060010190565b600063ffffffff80831681811415620051c557620051c5620051e6565b6001019392505050565b600082620051e157620051e162005215565b500690565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60038110620026bf57620026bf62005244565b73ffffffffffffffffffffffffffffffffffffffff81168114620026bf57600080fd5b8015158114620026bf57600080fdfe60c060405234801561001057600080fd5b5060405161057538038061057583398101604081905261002f9161005b565b60008054336001600160a01b031990911617905560601b6001600160601b03191660805260a052610098565b6000806040838503121561006e57600080fd5b825160208401519092506001600160a01b038116811461008d57600080fd5b809150509250929050565b60805160601c60a0516104ae6100c76000396000610145015260008181610170015261028001526104ae6000f3fe608060405234801561001057600080fd5b50600436106100725760003560e01c806379188d161161005057806379188d161461011d5780638ee489b214610140578063f00e6a2a1461016e57600080fd5b8063181f5a77146100775780631a5da6c8146100c95780635ab1bd53146100de575b600080fd5b6100b36040518060400160405280601981526020017f4175746f6d6174696f6e466f7277617264657220312e302e300000000000000081525081565b6040516100c091906103ff565b60405180910390f35b6100dc6100d73660046102e9565b610194565b005b60005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100c0565b61013061012b366004610326565b61022c565b60405190151581526020016100c0565b6040517f000000000000000000000000000000000000000000000000000000000000000081526020016100c0565b7f00000000000000000000000000000000000000000000000000000000000000006100f8565b60005473ffffffffffffffffffffffffffffffffffffffff1633146101e5576040517fea8e4eb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b6000805473ffffffffffffffffffffffffffffffffffffffff16331461027e576040517fea8e4eb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f00000000000000000000000000000000000000000000000000000000000000005a6113888110156102af57600080fd5b6113888103905084604082048203116102c757600080fd5b50803b6102d357600080fd5b60008084516020860160008589f1949350505050565b6000602082840312156102fb57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461031f57600080fd5b9392505050565b6000806040838503121561033957600080fd5b82359150602083013567ffffffffffffffff8082111561035857600080fd5b818501915085601f83011261036c57600080fd5b81358181111561037e5761037e610472565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019083821181831017156103c4576103c4610472565b816040528281528860208487010111156103dd57600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b600060208083528351808285015260005b8181101561042c57858101830151858201604001528201610410565b8181111561043e576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000806000a",
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicACaller) INext(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogicA.contract.Call(opts, &out, "i_next")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) INext() (common.Address, error) {
	return _KeeperRegistryLogicA.Contract.INext(&_KeeperRegistryLogicA.CallOpts)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicACallerSession) INext() (common.Address, error) {
	return _KeeperRegistryLogicA.Contract.INext(&_KeeperRegistryLogicA.CallOpts)
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

<<<<<<< HEAD
<<<<<<< HEAD
func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "acceptPayeeship", transmitter)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.AcceptPayeeship(&_KeeperRegistryLogicA.TransactOpts, transmitter)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.AcceptPayeeship(&_KeeperRegistryLogicA.TransactOpts, transmitter)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "acceptUpkeepAdmin", id)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.AcceptUpkeepAdmin(&_KeeperRegistryLogicA.TransactOpts, id)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.AcceptUpkeepAdmin(&_KeeperRegistryLogicA.TransactOpts, id)
=======
=======
>>>>>>> develop
func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "addFunds", id, amount)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.AddFunds(&_KeeperRegistryLogicA.TransactOpts, id, amount)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.AddFunds(&_KeeperRegistryLogicA.TransactOpts, id, amount)
<<<<<<< HEAD
>>>>>>> automation-log-trigger-support
=======
>>>>>>> develop
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

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "migrateUpkeeps", ids, destination)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.MigrateUpkeeps(&_KeeperRegistryLogicA.TransactOpts, ids, destination)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.MigrateUpkeeps(&_KeeperRegistryLogicA.TransactOpts, ids, destination)
}

<<<<<<< HEAD
<<<<<<< HEAD
func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "pauseUpkeep", id)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.PauseUpkeep(&_KeeperRegistryLogicA.TransactOpts, id)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.PauseUpkeep(&_KeeperRegistryLogicA.TransactOpts, id)
}

=======
>>>>>>> automation-log-trigger-support
=======
>>>>>>> develop
func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "receiveUpkeeps", encodedUpkeeps)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.ReceiveUpkeeps(&_KeeperRegistryLogicA.TransactOpts, encodedUpkeeps)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.ReceiveUpkeeps(&_KeeperRegistryLogicA.TransactOpts, encodedUpkeeps)
}

<<<<<<< HEAD
<<<<<<< HEAD
func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "setUpkeepGasLimit", id, gasLimit)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.SetUpkeepGasLimit(&_KeeperRegistryLogicA.TransactOpts, id, gasLimit)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.SetUpkeepGasLimit(&_KeeperRegistryLogicA.TransactOpts, id, gasLimit)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) SetUpkeepOffchainConfig(opts *bind.TransactOpts, id *big.Int, config []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "setUpkeepOffchainConfig", id, config)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.SetUpkeepOffchainConfig(&_KeeperRegistryLogicA.TransactOpts, id, config)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.SetUpkeepOffchainConfig(&_KeeperRegistryLogicA.TransactOpts, id, config)
=======
=======
>>>>>>> develop
func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte, extraData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "registerUpkeep", target, gasLimit, admin, checkData, extraData)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, checkData []byte, extraData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.RegisterUpkeep(&_KeeperRegistryLogicA.TransactOpts, target, gasLimit, admin, checkData, extraData)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, checkData []byte, extraData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.RegisterUpkeep(&_KeeperRegistryLogicA.TransactOpts, target, gasLimit, admin, checkData, extraData)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) SetUpkeepTriggerConfig(opts *bind.TransactOpts, id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "setUpkeepTriggerConfig", id, triggerConfig)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.SetUpkeepTriggerConfig(&_KeeperRegistryLogicA.TransactOpts, id, triggerConfig)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) SetUpkeepTriggerConfig(id *big.Int, triggerConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.SetUpkeepTriggerConfig(&_KeeperRegistryLogicA.TransactOpts, id, triggerConfig)
<<<<<<< HEAD
>>>>>>> automation-log-trigger-support
=======
>>>>>>> develop
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

<<<<<<< HEAD
<<<<<<< HEAD
func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "transferPayeeship", transmitter, proposed)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.TransferPayeeship(&_KeeperRegistryLogicA.TransactOpts, transmitter, proposed)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.TransferPayeeship(&_KeeperRegistryLogicA.TransactOpts, transmitter, proposed)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "transferUpkeepAdmin", id, proposed)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.TransferUpkeepAdmin(&_KeeperRegistryLogicA.TransactOpts, id, proposed)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.TransferUpkeepAdmin(&_KeeperRegistryLogicA.TransactOpts, id, proposed)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "unpauseUpkeep", id)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.UnpauseUpkeep(&_KeeperRegistryLogicA.TransactOpts, id)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.UnpauseUpkeep(&_KeeperRegistryLogicA.TransactOpts, id)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) UpdateCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "updateCheckData", id, newCheckData)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) UpdateCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.UpdateCheckData(&_KeeperRegistryLogicA.TransactOpts, id, newCheckData)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) UpdateCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.UpdateCheckData(&_KeeperRegistryLogicA.TransactOpts, id, newCheckData)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "withdrawFunds", id, to)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.WithdrawFunds(&_KeeperRegistryLogicA.TransactOpts, id, to)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.WithdrawFunds(&_KeeperRegistryLogicA.TransactOpts, id, to)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "withdrawPayment", from, to)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.WithdrawPayment(&_KeeperRegistryLogicA.TransactOpts, from, to)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.WithdrawPayment(&_KeeperRegistryLogicA.TransactOpts, from, to)
=======
=======
>>>>>>> develop
func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactor) ValidateTriggerConfig(opts *bind.TransactOpts, triggerType uint8, triggerConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.contract.Transact(opts, "validateTriggerConfig", triggerType, triggerConfig)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicASession) ValidateTriggerConfig(triggerType uint8, triggerConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.ValidateTriggerConfig(&_KeeperRegistryLogicA.TransactOpts, triggerType, triggerConfig)
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicATransactorSession) ValidateTriggerConfig(triggerType uint8, triggerConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogicA.Contract.ValidateTriggerConfig(&_KeeperRegistryLogicA.TransactOpts, triggerType, triggerConfig)
<<<<<<< HEAD
>>>>>>> automation-log-trigger-support
=======
>>>>>>> develop
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
	Id  *big.Int
	Raw types.Log
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
	Id  *big.Int
	Raw types.Log
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
	Id  *big.Int
	Raw types.Log
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
	Id  *big.Int
	Raw types.Log
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

type KeeperRegistryLogicAUpkeepCheckDataUpdatedIterator struct {
	Event *KeeperRegistryLogicAUpkeepCheckDataUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicAUpkeepCheckDataUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicAUpkeepCheckDataUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryLogicAUpkeepCheckDataUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicAUpkeepCheckDataUpdatedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicAUpkeepCheckDataUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicAUpkeepCheckDataUpdated struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) FilterUpkeepCheckDataUpdated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepCheckDataUpdatedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.FilterLogs(opts, "UpkeepCheckDataUpdated", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicAUpkeepCheckDataUpdatedIterator{contract: _KeeperRegistryLogicA.contract, event: "UpkeepCheckDataUpdated", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) WatchUpkeepCheckDataUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepCheckDataUpdated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogicA.contract.WatchLogs(opts, "UpkeepCheckDataUpdated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicAUpkeepCheckDataUpdated)
				if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepCheckDataUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogicA *KeeperRegistryLogicAFilterer) ParseUpkeepCheckDataUpdated(log types.Log) (*KeeperRegistryLogicAUpkeepCheckDataUpdated, error) {
	event := new(KeeperRegistryLogicAUpkeepCheckDataUpdated)
	if err := _KeeperRegistryLogicA.contract.UnpackLog(event, "UpkeepCheckDataUpdated", log); err != nil {
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

<<<<<<< HEAD
<<<<<<< HEAD
=======
=======
>>>>>>> develop
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

<<<<<<< HEAD
>>>>>>> automation-log-trigger-support
=======
>>>>>>> develop
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
	case _KeeperRegistryLogicA.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepAdminTransferRequested(log)
	case _KeeperRegistryLogicA.abi.Events["UpkeepAdminTransferred"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepAdminTransferred(log)
	case _KeeperRegistryLogicA.abi.Events["UpkeepCanceled"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepCanceled(log)
	case _KeeperRegistryLogicA.abi.Events["UpkeepCheckDataUpdated"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepCheckDataUpdated(log)
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
	case _KeeperRegistryLogicA.abi.Events["UpkeepReceived"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepReceived(log)
	case _KeeperRegistryLogicA.abi.Events["UpkeepRegistered"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepRegistered(log)
<<<<<<< HEAD
<<<<<<< HEAD
=======
	case _KeeperRegistryLogicA.abi.Events["UpkeepTriggerConfigSet"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepTriggerConfigSet(log)
>>>>>>> automation-log-trigger-support
=======
	case _KeeperRegistryLogicA.abi.Events["UpkeepTriggerConfigSet"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepTriggerConfigSet(log)
>>>>>>> develop
	case _KeeperRegistryLogicA.abi.Events["UpkeepUnpaused"].ID:
		return _KeeperRegistryLogicA.ParseUpkeepUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeeperRegistryLogicACancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xd84831b6a3a7fbd333f42fe7f9104a139da6cca4cc1507aef4ddad79b31d017f")
}

func (KeeperRegistryLogicAFundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (KeeperRegistryLogicAFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (KeeperRegistryLogicAInsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x7895fdfe292beab0842d5beccd078e85296b9e17a30eaee4c261a2696b84eb96")
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
	return common.HexToHash("0x561ff77e59394941a01a456497a9418dea82e2a39abb3ecebfb1cef7e0bfdc13")
}

func (KeeperRegistryLogicAStaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x5aa44821f7938098502bff537fbbdc9aaaa2fa655c10740646fce27e54987a89")
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

func (KeeperRegistryLogicAUpkeepCheckDataUpdated) Topic() common.Hash {
	return common.HexToHash("0x7b778136e5211932b51a145badd01959415e79e051a933604b3d323f862dcabf")
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

func (KeeperRegistryLogicAUpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (KeeperRegistryLogicAUpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

<<<<<<< HEAD
<<<<<<< HEAD
=======
=======
>>>>>>> develop
func (KeeperRegistryLogicAUpkeepTriggerConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2b72ac786c97e68dbab71023ed6f2bdbfc80ad9bb7808941929229d71b7d5664")
}

<<<<<<< HEAD
>>>>>>> automation-log-trigger-support
=======
>>>>>>> develop
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

	INext(opts *bind.CallOpts) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	UpkeepTranscoderVersion(opts *bind.CallOpts) (uint8, error)

	UpkeepVersion(opts *bind.CallOpts) (uint8, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

<<<<<<< HEAD
<<<<<<< HEAD
	AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error)

	AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)
=======
	AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error)
>>>>>>> automation-log-trigger-support
=======
	AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error)
>>>>>>> develop

	CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	CheckUpkeep(opts *bind.TransactOpts, id *big.Int, checkData []byte) (*types.Transaction, error)

	CheckUpkeep0(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error)

<<<<<<< HEAD
<<<<<<< HEAD
	PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error)

	SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error)

	SetUpkeepOffchainConfig(opts *bind.TransactOpts, id *big.Int, config []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error)

	TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error)

	UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	UpdateCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error)

	WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error)

	WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)
=======
=======
>>>>>>> develop
	ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error)

	RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte, extraData []byte) (*types.Transaction, error)

	SetUpkeepTriggerConfig(opts *bind.TransactOpts, id *big.Int, triggerConfig []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	ValidateTriggerConfig(opts *bind.TransactOpts, triggerType uint8, triggerConfig []byte) (*types.Transaction, error)
<<<<<<< HEAD
>>>>>>> automation-log-trigger-support
=======
>>>>>>> develop

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

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogicAUpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*KeeperRegistryLogicAUpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogicAUpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*KeeperRegistryLogicAUpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistryLogicAUpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*KeeperRegistryLogicAUpkeepCanceled, error)

	FilterUpkeepCheckDataUpdated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepCheckDataUpdatedIterator, error)

	WatchUpkeepCheckDataUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepCheckDataUpdated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCheckDataUpdated(log types.Log) (*KeeperRegistryLogicAUpkeepCheckDataUpdated, error)

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

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*KeeperRegistryLogicAUpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*KeeperRegistryLogicAUpkeepRegistered, error)

<<<<<<< HEAD
<<<<<<< HEAD
=======
=======
>>>>>>> develop
	FilterUpkeepTriggerConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepTriggerConfigSetIterator, error)

	WatchUpkeepTriggerConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepTriggerConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepTriggerConfigSet(log types.Log) (*KeeperRegistryLogicAUpkeepTriggerConfigSet, error)

<<<<<<< HEAD
>>>>>>> automation-log-trigger-support
=======
>>>>>>> develop
	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicAUpkeepUnpausedIterator, error)

	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicAUpkeepUnpaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepUnpaused(log types.Log) (*KeeperRegistryLogicAUpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
