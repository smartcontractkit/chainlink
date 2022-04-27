// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package keeper_registry_wrapper1_2

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
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
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
)

type Config struct {
	PaymentPremiumPPB    uint32
	FlatFeeMicroLink     uint32
	BlockCountPerTurn    *big.Int
	CheckGasLimit        uint32
	StalenessSeconds     *big.Int
	GasCeilingMultiplier uint16
	MinUpkeepSpend       *big.Int
	MaxPerformGas        uint32
	FallbackGasPrice     *big.Int
	FallbackLinkPrice    *big.Int
	Transcoder           common.Address
	Registrar            common.Address
}

type State struct {
	Nonce            uint32
	OwnerLinkBalance *big.Int
	NumUpkeeps       *big.Int
}

var KeeperRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkEthFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"fastGasFeed\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"blockCountPerTurn\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"registrar\",\"type\":\"address\"}],\"internalType\":\"structConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotChangePayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"KeepersMustTakeTurns\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveKeepers\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooFewKeepers\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotActive\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"blockCountPerTurn\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"registrar\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"keepers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"KeepersUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"FAST_GAS_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_ETH_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"maxLinkPayment\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"adjustedGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkEth\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"query\",\"type\":\"address\"}],\"name\":\"getKeeperInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"name\":\"getMaxPaymentForGas\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"maxPayment\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getMinBalanceForUpkeep\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"minBalance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"}],\"name\":\"getPeerRegistryMigrationPermission\",\"outputs\":[{\"internalType\":\"enumKeeperRegistry.MigrationPermission\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"ownerLinkBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"numUpkeeps\",\"type\":\"uint256\"}],\"internalType\":\"structState\",\"name\":\"state\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"blockCountPerTurn\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"registrar\",\"type\":\"address\"}],\"internalType\":\"structConfig\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"keepers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getUpkeep\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"lastKeeper\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"migrateUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedUpkeeps\",\"type\":\"bytes\"}],\"name\":\"receiveUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"blockCountPerTurn\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint96\",\"name\":\"minUpkeepSpend\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"maxPerformGas\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transcoder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"registrar\",\"type\":\"address\"}],\"internalType\":\"structConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"keepers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setKeepers\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"},{\"internalType\":\"enumKeeperRegistry.MigrationPermission\",\"name\":\"permission\",\"type\":\"uint8\"}],\"name\":\"setPeerRegistryMigrationPermission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTranscoderVersion\",\"outputs\":[{\"internalType\":\"enumUpkeepFormat\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawOwnerFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b5060405162006938380380620069388339810160408190526200003491620005a5565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000fe565b50506001600255506003805460ff191690556001600160a01b0380851660805283811660a052821660c052620000f481620001a9565b50505050620007f0565b336001600160a01b03821603620001585760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b620001b36200049e565b600d5460e082015163ffffffff91821691161015620001e557604051630e6af04160e21b815260040160405180910390fd5b604051806101200160405280826000015163ffffffff168152602001826020015163ffffffff168152602001826040015162ffffff168152602001826060015163ffffffff168152602001826080015162ffffff1681526020018260a0015161ffff1681526020018260c001516001600160601b031681526020018260e0015163ffffffff168152602001600c60010160049054906101000a900463ffffffff1663ffffffff16815250600c60008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160086101000a81548162ffffff021916908362ffffff160217905550606082015181600001600b6101000a81548163ffffffff021916908363ffffffff160217905550608082015181600001600f6101000a81548162ffffff021916908362ffffff16021790555060a08201518160000160126101000a81548161ffff021916908361ffff16021790555060c08201518160000160146101000a8154816001600160601b0302191690836001600160601b0316021790555060e08201518160010160006101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160010160046101000a81548163ffffffff021916908363ffffffff160217905550905050806101000151600e81905550806101200151600f81905550806101400151601260006101000a8154816001600160a01b0302191690836001600160a01b03160217905550806101600151601360006101000a8154816001600160a01b0302191690836001600160a01b031602179055507ffe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de32581604051620004939190620006f1565b60405180910390a150565b6000546001600160a01b03163314620004fa5760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000082565b565b80516001600160a01b03811681146200051457600080fd5b919050565b60405161018081016001600160401b03811182821017156200054b57634e487b7160e01b600052604160045260246000fd5b60405290565b805163ffffffff811681146200051457600080fd5b805162ffffff811681146200051457600080fd5b805161ffff811681146200051457600080fd5b80516001600160601b03811681146200051457600080fd5b6000806000808486036101e0811215620005be57600080fd5b620005c986620004fc565b9450620005d960208701620004fc565b9350620005e960408701620004fc565b925061018080605f19830112156200060057600080fd5b6200060a62000519565b91506200061a6060880162000551565b82526200062a6080880162000551565b60208301526200063d60a0880162000566565b60408301526200065060c0880162000551565b60608301526200066360e0880162000566565b6080830152610100620006788189016200057a565b60a08401526101206200068d818a016200058d565b60c0850152610140620006a2818b0162000551565b60e0860152610160808b015184870152848b015183870152620006c96101a08c01620004fc565b82870152620006dc6101c08c01620004fc565b90860152509699959850939650909450505050565b815163ffffffff1681526101808101602083015162000718602084018263ffffffff169052565b50604083015162000730604084018262ffffff169052565b50606083015162000749606084018263ffffffff169052565b50608083015162000761608084018262ffffff169052565b5060a08301516200077860a084018261ffff169052565b5060c08301516200079460c08401826001600160601b03169052565b5060e0830151620007ad60e084018263ffffffff169052565b5061010083810151908301526101208084015190830152610140808401516001600160a01b03908116918401919091526101609384015116929091019190915290565b60805160a05160c0516160d262000866600039600081816104240152614435015260008181610575015261450701526000818161030401528181610dff0152818161111b01528181611a2001528181611ed101528181611fb6015281816123ae0152818161271d01526127a101526160d26000f3fe608060405234801561001057600080fd5b506004361061025c5760003560e01c806393f0c1fc11610145578063b7fdb436116100bd578063da5c67411161008c578063ef47a0ce11610071578063ef47a0ce1461066a578063f2fde38b1461067d578063faa3e9961461069057600080fd5b8063da5c674114610636578063eb5dcd6c1461065757600080fd5b8063b7fdb436146105c5578063c41b813a146105d8578063c7c3a19a146105fc578063c80480221461062357600080fd5b8063a72aa27e11610114578063b121e147116100f9578063b121e14714610597578063b657bc9c146105aa578063b79550be146105bd57600080fd5b8063a72aa27e1461055d578063ad1783611461057057600080fd5b806393f0c1fc146104f4578063948108f714610524578063a4c0ed3614610537578063a710b2211461054a57600080fd5b80635c975abb116101d85780637d9b97e0116101a757806385c1b0ba1161018c57806385c1b0ba146104b05780638da5cb5b146104c35780638e86139b146104e157600080fd5b80637d9b97e0146104a05780638456cb59146104a857600080fd5b80635c975abb1461045b578063744bfe611461047257806379ba5097146104855780637bbaf1ea1461048d57600080fd5b80631b6b6d231161022f5780633f4ba83a116102145780633f4ba83a146104175780634584a4191461041f57806348013d7b1461044657600080fd5b80631b6b6d23146102ff5780631e12b8a51461034b57600080fd5b806306e3b63214610261578063181f5a771461028a5780631865c57d146102d3578063187256e8146102ea575b600080fd5b61027461026f366004614cb1565b6106d6565b6040516102819190614cd3565b60405180910390f35b6102c66040518060400160405280601481526020017f4b6565706572526567697374727920312e322e3000000000000000000000000081525081565b6040516102819190614d91565b6102db6107d5565b60405161028193929190614ebb565b6102fd6102f8366004614f81565b610a79565b005b6103267f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610281565b6103d7610359366004614fbc565b73ffffffffffffffffffffffffffffffffffffffff90811660009081526008602090815260409182902082516060810184528154948516808252740100000000000000000000000000000000000000009095046bffffffffffffffffffffffff1692810183905260019091015460ff16151592018290529192909190565b6040805173ffffffffffffffffffffffffffffffffffffffff909416845291151560208401526bffffffffffffffffffffffff1690820152606001610281565b6102fd610aea565b6103267f000000000000000000000000000000000000000000000000000000000000000081565b61044e600081565b604051610281919061501a565b60035460ff165b6040519015158152602001610281565b6102fd610480366004615028565b610afc565b6102fd610e78565b61046261049b36600461509d565b610f7a565b6102fd611043565b6102fd6111a2565b6102fd6104be36600461512e565b6111b2565b60005473ffffffffffffffffffffffffffffffffffffffff16610326565b6102fd6104ef366004615182565b611a9d565b6105076105023660046151c4565b611d0a565b6040516bffffffffffffffffffffffff9091168152602001610281565b6102fd6105323660046151f9565b611d3e565b6102fd61054536600461521c565b611f9e565b6102fd610558366004615276565b612199565b6102fd61056b3660046152b4565b612424565b6103267f000000000000000000000000000000000000000000000000000000000000000081565b6102fd6105a5366004614fbc565b6125cb565b6105076105b83660046151c4565b6126c3565b6102fd6126e4565b6102fd6105d33660046152d7565b612840565b6105eb6105e6366004615028565b612bfe565b604051610281959493929190615337565b61060f61060a3660046151c4565b612f26565b60405161028198979695949392919061536e565b6102fd6106313660046151c4565b6130b1565b6106496106443660046153f8565b6132a7565b604051908152602001610281565b6102fd610665366004615276565b61349e565b6102fd61067836600461555e565b6135fc565b6102fd61068b366004614fbc565b613948565b6106c961069e366004614fbc565b73ffffffffffffffffffffffffffffffffffffffff166000908152600b602052604090205460ff1690565b604051610281919061563c565b606060006106e4600561395c565b905080841061071f576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b82600003610734576107318482615685565b92505b60008367ffffffffffffffff81111561074f5761074f61546e565b604051908082528060200260200182016040528015610778578160200160208202803683370190505b50905060005b848110156107ca5761079b610793828861569c565b600590613966565b8282815181106107ad576107ad6156b4565b6020908102919091010152806107c2816156e3565b91505061077e565b509150505b92915050565b60408051606081018252600080825260208201819052918101919091526040805161018081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101829052610120810182905261014081018290526101608101919091526040805161012081018252600c5463ffffffff8082168352640100000000808304821660208086019190915262ffffff6801000000000000000085048116968601969096526b010000000000000000000000840483166060868101919091526f010000000000000000000000000000008504909616608086015261ffff720100000000000000000000000000000000000085041660a08601526bffffffffffffffffffffffff74010000000000000000000000000000000000000000909404841660c0860152600d5480841660e0870152919091049091166101008401819052865260105490911690850152610949600561395c565b604080860191909152815163ffffffff90811685526020808401518216818701528383015162ffffff908116878501526060808601518416908801526080808601519091169087015260a08085015161ffff169087015260c0808501516bffffffffffffffffffffffff169087015260e08085015190921691860191909152600e54610100860152600f5461012086015260125473ffffffffffffffffffffffffffffffffffffffff90811661014087015260135416610160860152600480548351818402810184019094528084528793879390918391830182828015610a6657602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610a3b575b5050505050905093509350935050909192565b610a81613979565b73ffffffffffffffffffffffffffffffffffffffff82166000908152600b6020526040902080548291907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001836003811115610ae157610ae1614fd7565b02179055505050565b610af2613979565b610afa6139fa565b565b8073ffffffffffffffffffffffffffffffffffffffff8116610b4a576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008381526007602052604090206002015483906c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff163314610bbc576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000848152600760205260409020600101544364010000000090910467ffffffffffffffff161115610c1a576040517fff84e5dd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600c54600085815260076020526040812080546002909101546bffffffffffffffffffffffff740100000000000000000000000000000000000000009094048416939182169291169083821015610c9e57610c75828561571b565b9050826bffffffffffffffffffffffff16816bffffffffffffffffffffffff161115610c9e5750815b6000610caa828561571b565b60008a815260076020526040902080547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000169055601054909150610cfd9083906bffffffffffffffffffffffff16615748565b601080547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff928316179055601154610d4591831690615685565b601155604080516bffffffffffffffffffffffff8316815273ffffffffffffffffffffffffffffffffffffffff8a1660208201528a917ff3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318910160405180910390a26040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff89811660048301526bffffffffffffffffffffffff831660248301527f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb906044016020604051808303816000875af1158015610e48573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e6c9190615788565b50505050505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610efe576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000610f8860035460ff1690565b15610fef576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a20706175736564000000000000000000000000000000006044820152606401610ef5565b61103b611036338686868080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525060019250613adb915050565b613bd5565b949350505050565b61104b613979565b6010546011546bffffffffffffffffffffffff9091169061106d908290615685565b601155601080547fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001690556040516bffffffffffffffffffffffff821681527f1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f19060200160405180910390a16040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526bffffffffffffffffffffffff821660248201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb906044015b6020604051808303816000875af115801561117a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061119e9190615788565b5050565b6111aa613979565b610afa614002565b600173ffffffffffffffffffffffffffffffffffffffff82166000908152600b602052604090205460ff1660038111156111ee576111ee614fd7565b141580156112365750600373ffffffffffffffffffffffffffffffffffffffff82166000908152600b602052604090205460ff16600381111561123357611233614fd7565b14155b1561126d576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60125473ffffffffffffffffffffffffffffffffffffffff166112bc576040517fd12d7d8d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008290036112f7576040517f2c2fc94100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006007600085856000818110611310576113106156b4565b905060200201358152602001908152602001600020600201600c9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050600061136c60005473ffffffffffffffffffffffffffffffffffffffff1690565b3373ffffffffffffffffffffffffffffffffffffffff9182168114925090831614801590611398575080155b156113cf576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160e081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c081018290526000808767ffffffffffffffff8111156114235761142361546e565b60405190808252806020026020018201604052801561145657816020015b60608152602001906001900390816114415790505b50905060008867ffffffffffffffff8111156114745761147461546e565b6040519080825280602002602001820160405280156114f957816020015b6040805160e08101825260008082526020808301829052928201819052606082018190526080820181905260a0820181905260c082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816114925790505b50905060005b898110156117fc578a8a82818110611519576115196156b4565b60209081029290920135600081815260078452604090819020815160e08101835281546bffffffffffffffffffffffff808216835273ffffffffffffffffffffffffffffffffffffffff6c0100000000000000000000000092839004811698840198909852600184015463ffffffff81169584019590955267ffffffffffffffff6401000000008604166060840152938190048716608083015260029092015492831660a08201529104841660c082018190529199509750918a16909114905061160f576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b606085015167ffffffffffffffff90811614611657576040517fd096219c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8482828151811061166a5761166a6156b4565b6020026020010181905250600a60008781526020019081526020016000208054611693906157a3565b80601f01602080910402602001604051908101604052809291908181526020018280546116bf906157a3565b801561170c5780601f106116e15761010080835404028352916020019161170c565b820191906000526020600020905b8154815290600101906020018083116116ef57829003601f168201915b5050505050838281518110611723576117236156b4565b60209081029190910101528451611748906bffffffffffffffffffffffff168561569c565b600087815260076020908152604080832083815560018101849055600201839055600a909152812091955061177d9190614b66565b6117886005876140c2565b508451604080516bffffffffffffffffffffffff909216825273ffffffffffffffffffffffffffffffffffffffff8b16602083015287917fb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff910160405180910390a2806117f4816156e3565b9150506114ff565b508260115461180b9190615685565b601155604051600090611828908c908c908590879060200161584b565b60405160208183030381529060405290508873ffffffffffffffffffffffffffffffffffffffff16638e86139b601260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663c71249ab60008d73ffffffffffffffffffffffffffffffffffffffff166348013d7b6040518163ffffffff1660e01b81526004016020604051808303816000875af11580156118e2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906119069190615999565b866040518463ffffffff1660e01b8152600401611925939291906159ba565b600060405180830381865afa158015611942573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526119889190810190615a77565b6040518263ffffffff1660e01b81526004016119a49190614d91565b600060405180830381600087803b1580156119be57600080fd5b505af11580156119d2573d6000803e3d6000fd5b50506040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8c81166004830152602482018890527f000000000000000000000000000000000000000000000000000000000000000016925063a9059cbb91506044016020604051808303816000875af1158015611a6b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611a8f9190615788565b505050505050505050505050565b60035460ff1615611b0a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a20706175736564000000000000000000000000000000006044820152606401610ef5565b6002336000908152600b602052604090205460ff166003811115611b3057611b30614fd7565b14158015611b6257506003336000908152600b602052604090205460ff166003811115611b5f57611b5f614fd7565b14155b15611b99576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008080611ba984860186615c97565b92509250925060005b8351811015611d0257611c6f848281518110611bd057611bd06156b4565b6020026020010151848381518110611bea57611bea6156b4565b602002602001015160800151858481518110611c0857611c086156b4565b602002602001015160400151868581518110611c2657611c266156b4565b602002602001015160c00151878681518110611c4457611c446156b4565b602002602001015160000151878781518110611c6257611c626156b4565b60200260200101516140ce565b838181518110611c8157611c816156b4565b60200260200101517f74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71848381518110611cbc57611cbc6156b4565b60209081029190910181015151604080516bffffffffffffffffffffffff909216825233928201929092520160405180910390a280611cfa816156e3565b915050611bb2565b505050505050565b6000806000611d17614402565b915091506000611d288360006145df565b9050611d35858284614624565b95945050505050565b60035460ff1615611dab576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a20706175736564000000000000000000000000000000006044820152606401610ef5565b6000828152600760205260409020600101548290640100000000900467ffffffffffffffff90811614611e0a576040517fd096219c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600083815260076020526040902054611e329083906bffffffffffffffffffffffff16615748565b600084815260076020526040902080547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff928316179055601154611e869184169061569c565b6011556040517f23b872dd0000000000000000000000000000000000000000000000000000000081523360048201523060248201526bffffffffffffffffffffffff831660448201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906323b872dd906064016020604051808303816000875af1158015611f2f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611f539190615788565b506040516bffffffffffffffffffffffff83168152339084907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a3505050565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461200d576040517fc8bad78d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208114612047576040517fdfe9309000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000612055828401846151c4565b600081815260076020526040902060010154909150640100000000900467ffffffffffffffff908116146120b5576040517fd096219c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000818152600760205260409020546120dd9085906bffffffffffffffffffffffff16615748565b600082815260076020526040902080547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff9290921691909117905560115461213490859061569c565b6011556040516bffffffffffffffffffffffff8516815273ffffffffffffffffffffffffffffffffffffffff86169082907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a35050505050565b8073ffffffffffffffffffffffffffffffffffffffff81166121e7576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff83811660009081526008602090815260409182902082516060810184528154948516808252740100000000000000000000000000000000000000009095046bffffffffffffffffffffffff16928101929092526001015460ff16151591810191909152903314612298576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff808516600090815260086020908152604090912080549092169091558101516011546122e7916bffffffffffffffffffffffff1690615685565b60115560208082015160405133815273ffffffffffffffffffffffffffffffffffffffff808716936bffffffffffffffffffffffff90931692908816917f9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698910160405180910390a460208101516040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff85811660048301526bffffffffffffffffffffffff90921660248201527f00000000000000000000000000000000000000000000000000000000000000009091169063a9059cbb906044016020604051808303816000875af11580156123f9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061241d9190615788565b5050505050565b6000828152600760205260409020600101548290640100000000900467ffffffffffffffff90811614612483576040517fd096219c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008381526007602052604090206002015483906c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1633146124f5576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108fc8363ffffffff1610806125165750600d5463ffffffff908116908416115b1561254d576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008481526007602090815260409182902060010180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff8716908117909155915191825285917fc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c910160405180910390a250505050565b73ffffffffffffffffffffffffffffffffffffffff81811660009081526009602052604090205416331461262b576040517f6752e7aa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff81811660008181526008602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556009909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b6000818152600760205260408120600101546107cf9063ffffffff16611d0a565b6126ec613979565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa158015612779573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061279d9190615d74565b90507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb33601154846127ea9190615685565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff9092166004830152602482015260440161115b565b612848613979565b828114612881576040517fb1c97b3900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60028310156128bc576040517febfd058f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b600454811015612948576000600482815481106128de576128de6156b4565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168252600890526040902060010180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690555080612940816156e3565b9150506128bf565b5060005b83811015612bad576000858583818110612968576129686156b4565b905060200201602081019061297d9190614fbc565b73ffffffffffffffffffffffffffffffffffffffff8082166000908152600860205260408120805493945092909116908686868181106129bf576129bf6156b4565b90506020020160208101906129d49190614fbc565b905073ffffffffffffffffffffffffffffffffffffffff8116612a23576040517fb387a23800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff821615801590612a7457508073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b8015612a96575073ffffffffffffffffffffffffffffffffffffffff81811614155b15612acd576040517f9f43d9bc00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600183015460ff1615612b0c576040517f357d0cc400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600183810180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016909117905573ffffffffffffffffffffffffffffffffffffffff81811614612b965782547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff82161783555b505050508080612ba5906156e3565b91505061294c565b50612bba60048585614ba0565b507f056264c94f28bb06c99d13f0446eb96c67c215d8d707bce2655a98ddf1c0b71f84848484604051612bf09493929190615de1565b60405180910390a150505050565b6060600080600080612c1260035460ff1690565b15612c79576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a20706175736564000000000000000000000000000000006044820152606401610ef5565b612c81614701565b6000878152600760209081526040808320815160e08101835281546bffffffffffffffffffffffff808216835273ffffffffffffffffffffffffffffffffffffffff6c0100000000000000000000000092839004811684880152600185015463ffffffff81168588015267ffffffffffffffff64010000000082041660608601528390048116608085015260029094015490811660a08401520490911660c08201528a8452600a90925280832090519192917f6e04ff0d0000000000000000000000000000000000000000000000000000000091612d6191602401615e13565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050509050600080836080015173ffffffffffffffffffffffffffffffffffffffff16600c600001600b9054906101000a900463ffffffff1663ffffffff1684604051612e089190615ef1565b60006040518083038160008787f1925050503d8060008114612e46576040519150601f19603f3d011682016040523d82523d6000602084013e612e4b565b606091505b509150915081612e8957806040517f96c36235000000000000000000000000000000000000000000000000000000008152600401610ef59190614d91565b80806020019051810190612e9d9190615f0d565b9950915081612ed8576040517f865676e300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000612ee78b8d8c6000613adb565b9050612efc8582600001518360600151614739565b6060810151608082015160a083015160c0909301519b9e919d509b50909998509650505050505050565b6000818152600760209081526040808320815160e08101835281546bffffffffffffffffffffffff808216835273ffffffffffffffffffffffffffffffffffffffff6c01000000000000000000000000928390048116848801908152600186015463ffffffff811686890181905267ffffffffffffffff64010000000083041660608881019182529287900485166080890181905260029099015495861660a089019081529690950490931660c087019081528b8b52600a9099529689208551915198519351945181548b9a8b998a998a998a998a9992989397929692959394939092908690613015906157a3565b80601f0160208091040260200160405190810160405280929190818152602001828054613041906157a3565b801561308e5780601f106130635761010080835404028352916020019161308e565b820191906000526020600020905b81548152906001019060200180831161307157829003601f168201915b505050505095509850985098509850985098509850985050919395975091939597565b60008181526007602052604081206001015467ffffffffffffffff64010000000090910481169190821415906130fc60005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614905081801561314c575080801561314a5750438367ffffffffffffffff16115b155b15613183576040517ffbc0357800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b801580156131c857506000848152600760205260409020600201546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff163314155b156131ff576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b43816132135761321060328261569c565b90505b600085815260076020526040902060010180547fffffffffffffffffffffffffffffffffffffffff0000000000000000ffffffff1664010000000067ffffffffffffffff8416021790556132686005866140c2565b5060405167ffffffffffffffff82169086907f91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f79118190600090a35050505050565b6000805473ffffffffffffffffffffffffffffffffffffffff1633148015906132e8575060135473ffffffffffffffffffffffffffffffffffffffff163314155b1561331f576040517fd48b678b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61332a600143615685565b600d5460408051924060208401523060601b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001690830152640100000000900460e01b7fffffffff000000000000000000000000000000000000000000000000000000001660548201526058016040516020818303038152906040528051906020012060001c90506133f781878787600088888080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506140ce92505050565b600d8054640100000000900463ffffffff1690600461341583615f5b565b91906101000a81548163ffffffff021916908363ffffffff16021790555050807fbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012868660405161348d92919063ffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a295945050505050565b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600860205260409020541633146134fe576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff82160361354d576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff82811660009081526009602052604090205481169082161461119e5773ffffffffffffffffffffffffffffffffffffffff82811660008181526009602052604080822080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169486169485179055513392917f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836791a45050565b613604613979565b600d5460e082015163ffffffff9182169116101561364e576040517f39abc10400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604051806101200160405280826000015163ffffffff168152602001826020015163ffffffff168152602001826040015162ffffff168152602001826060015163ffffffff168152602001826080015162ffffff1681526020018260a0015161ffff1681526020018260c001516bffffffffffffffffffffffff1681526020018260e0015163ffffffff168152602001600c60010160049054906101000a900463ffffffff1663ffffffff16815250600c60008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160086101000a81548162ffffff021916908362ffffff160217905550606082015181600001600b6101000a81548163ffffffff021916908363ffffffff160217905550608082015181600001600f6101000a81548162ffffff021916908362ffffff16021790555060a08201518160000160126101000a81548161ffff021916908361ffff16021790555060c08201518160000160146101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060e08201518160010160006101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160010160046101000a81548163ffffffff021916908363ffffffff160217905550905050806101000151600e81905550806101200151600f81905550806101400151601260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550806101600151601360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507ffe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de3258160405161393d9190615f7e565b60405180910390a150565b613950613979565b61395981614852565b50565b60006107cf825490565b60006139728383614947565b9392505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610afa576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610ef5565b60035460ff16613a66576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f5061757361626c653a206e6f74207061757365640000000000000000000000006044820152606401610ef5565b600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390a1565b613b316040518060e00160405280600073ffffffffffffffffffffffffffffffffffffffff1681526020016000815260200160608152602001600081526020016000815260200160008152602001600081525090565b60008481526007602052604081206001015463ffffffff169080613b53614402565b915091506000613b6383876145df565b90506000613b72858385614624565b6040805160e08101825273ffffffffffffffffffffffffffffffffffffffff8d168152602081018c90529081018a90526bffffffffffffffffffffffff909116606082015260808101959095525060a084015260c0830152509050949350505050565b60006002805403613c42576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c006044820152606401610ef5565b600280556020820151613c5481614971565b602080840151600090815260078252604090819020815160e08101835281546bffffffffffffffffffffffff808216835273ffffffffffffffffffffffffffffffffffffffff6c0100000000000000000000000092839004811696840196909652600184015463ffffffff81169584019590955267ffffffffffffffff640100000000860416606080850191909152948290048616608084015260029093015492831660a083015290910490921660c0830152845190850151613d18918391614739565b60005a90506000634585e33b60e01b8660400151604051602401613d3c9190614d91565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050509050613dae86608001518460800151836149ce565b94505a613dbb9083615685565b91506000613dd2838860a001518960c00151614624565b602080890151600090815260079091526040902054909150613e039082906bffffffffffffffffffffffff1661571b565b6020888101805160009081526007909252604080832080547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff95861617905590518252902060020154613e6691839116615748565b60208881018051600090815260078352604080822060020180547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff9687161790558b5192518252808220805486166c0100000000000000000000000073ffffffffffffffffffffffffffffffffffffffff958616021790558b5190921681526008909252902054613f1f91839174010000000000000000000000000000000000000000900416615748565b60086000896000015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160146101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550866000015173ffffffffffffffffffffffffffffffffffffffff1686151588602001517fcaacad83e47cc45c280d487ec84184eee2fa3b54ebaa393bda7549f13da228f6848b60400151604051613feb929190615f8d565b60405180910390a450505050506001600255919050565b60035460ff161561406f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a20706175736564000000000000000000000000000000006044820152606401610ef5565b600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258613ab13390565b60006139728383614a1a565b73ffffffffffffffffffffffffffffffffffffffff85163b61411c576040517f09ee12d500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108fc8463ffffffff16108061413d5750600d5463ffffffff908116908516115b15614174576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040518060e00160405280836bffffffffffffffffffffffff168152602001600073ffffffffffffffffffffffffffffffffffffffff1681526020018563ffffffff16815260200167ffffffffffffffff801681526020018673ffffffffffffffffffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff168152506007600088815260200190815260200160002060008201518160000160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550602082015181600001600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060408201518160010160006101000a81548163ffffffff021916908363ffffffff16021790555060608201518160010160046101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550608082015181600101600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060a08201518160020160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060c082015181600201600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505080600a600088815260200190815260200160002090805190602001906143ed929190614c28565b506143f9600587614b0d565b50505050505050565b6000806000600c600001600f9054906101000a900462ffffff1662ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa15801561449e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906144c29190615fce565b5094509092508491505080156144e657506144dd8242615685565b8463ffffffff16105b806144f2575060008113155b1561450157600e549550614505565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a060405180830381865afa158015614570573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906145949190615fce565b5094509092508491505080156145b857506145af8242615685565b8463ffffffff16105b806145c4575060008113155b156145d357600f5494506145d7565b8094505b505050509091565b600c54600090614609907201000000000000000000000000000000000000900461ffff168461601e565b90508180156146175750803a105b156107cf57503a92915050565b600080614634620138808661569c565b61463e908561601e565b600c5490915060009061465b9063ffffffff16633b9aca0061569c565b600c5490915060009061468190640100000000900463ffffffff1664e8d4a5100061601e565b858361469186633b9aca0061601e565b61469b919061601e565b6146a5919061605b565b6146af919061569c565b90506b033b2e3c9fd0803ce80000008111156146f7576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b9695505050505050565b3215610afa576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff821660009081526008602052604090206001015460ff1661479b576040517fcfbacfd800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b82516bffffffffffffffffffffffff168111156147e4576040517f356680b700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8173ffffffffffffffffffffffffffffffffffffffff16836020015173ffffffffffffffffffffffffffffffffffffffff160361484d576040517f06bc104000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505050565b3373ffffffffffffffffffffffffffffffffffffffff8216036148d1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610ef5565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600082600001828154811061495e5761495e6156b4565b9060005260206000200154905092915050565b6000818152600760205260409020600101544364010000000090910467ffffffffffffffff1611613959576040517fd096219c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005a6113888110156149e057600080fd5b6113888103905084604082048203116149f857600080fd5b50823b614a0457600080fd5b60008083516020850160008789f1949350505050565b60008181526001830160205260408120548015614b03576000614a3e600183615685565b8554909150600090614a5290600190615685565b9050818114614ab7576000866000018281548110614a7257614a726156b4565b9060005260206000200154905080876000018481548110614a9557614a956156b4565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080614ac857614ac8616096565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506107cf565b60009150506107cf565b600081815260018301602052604081205461397290849084908490614b5e575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556107cf565b5060006107cf565b508054614b72906157a3565b6000825580601f10614b82575050565b601f0160209004906000526020600020908101906139599190614c9c565b828054828255906000526020600020908101928215614c18579160200282015b82811115614c185781547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff843516178255602090920191600190910190614bc0565b50614c24929150614c9c565b5090565b828054614c34906157a3565b90600052602060002090601f016020900481019282614c565760008555614c18565b82601f10614c6f57805160ff1916838001178555614c18565b82800160010185558215614c18579182015b82811115614c18578251825591602001919060010190614c81565b5b80821115614c245760008155600101614c9d565b60008060408385031215614cc457600080fd5b50508035926020909101359150565b6020808252825182820181905260009190848201906040850190845b81811015614d0b57835183529284019291840191600101614cef565b50909695505050505050565b60005b83811015614d32578181015183820152602001614d1a565b83811115614d41576000848401525b50505050565b60008151808452614d5f816020860160208601614d17565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006139726020830184614d47565b805163ffffffff1682526020810151614dc5602084018263ffffffff169052565b506040810151614ddc604084018262ffffff169052565b506060810151614df4606084018263ffffffff169052565b506080810151614e0b608084018262ffffff169052565b5060a0810151614e2160a084018261ffff169052565b5060c0810151614e4160c08401826bffffffffffffffffffffffff169052565b5060e0810151614e5960e084018263ffffffff169052565b50610100818101519083015261012080820151908301526101408082015173ffffffffffffffffffffffffffffffffffffffff81168285015250506101608181015173ffffffffffffffffffffffffffffffffffffffff811684830152614d41565b600061020080830163ffffffff875116845260206bffffffffffffffffffffffff81890151168186015260408801516040860152614efc6060860188614da4565b6101e08501929092528451908190526102208401918086019160005b81811015614f4a57835173ffffffffffffffffffffffffffffffffffffffff1685529382019392820192600101614f18565b509298975050505050505050565b803573ffffffffffffffffffffffffffffffffffffffff81168114614f7c57600080fd5b919050565b60008060408385031215614f9457600080fd5b614f9d83614f58565b9150602083013560048110614fb157600080fd5b809150509250929050565b600060208284031215614fce57600080fd5b61397282614f58565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b6001811061501657615016614fd7565b9052565b602081016107cf8284615006565b6000806040838503121561503b57600080fd5b8235915061504b60208401614f58565b90509250929050565b60008083601f84011261506657600080fd5b50813567ffffffffffffffff81111561507e57600080fd5b60208301915083602082850101111561509657600080fd5b9250929050565b6000806000604084860312156150b257600080fd5b83359250602084013567ffffffffffffffff8111156150d057600080fd5b6150dc86828701615054565b9497909650939450505050565b60008083601f8401126150fb57600080fd5b50813567ffffffffffffffff81111561511357600080fd5b6020830191508360208260051b850101111561509657600080fd5b60008060006040848603121561514357600080fd5b833567ffffffffffffffff81111561515a57600080fd5b615166868287016150e9565b9094509250615179905060208501614f58565b90509250925092565b6000806020838503121561519557600080fd5b823567ffffffffffffffff8111156151ac57600080fd5b6151b885828601615054565b90969095509350505050565b6000602082840312156151d657600080fd5b5035919050565b80356bffffffffffffffffffffffff81168114614f7c57600080fd5b6000806040838503121561520c57600080fd5b8235915061504b602084016151dd565b6000806000806060858703121561523257600080fd5b61523b85614f58565b935060208501359250604085013567ffffffffffffffff81111561525e57600080fd5b61526a87828801615054565b95989497509550505050565b6000806040838503121561528957600080fd5b61529283614f58565b915061504b60208401614f58565b803563ffffffff81168114614f7c57600080fd5b600080604083850312156152c757600080fd5b8235915061504b602084016152a0565b600080600080604085870312156152ed57600080fd5b843567ffffffffffffffff8082111561530557600080fd5b615311888389016150e9565b9096509450602087013591508082111561532a57600080fd5b5061526a878288016150e9565b60a08152600061534a60a0830188614d47565b90508560208301528460408301528360608301528260808301529695505050505050565b600061010073ffffffffffffffffffffffffffffffffffffffff808c16845263ffffffff8b1660208501528160408501526153ab8285018b614d47565b6bffffffffffffffffffffffff998a16606086015297811660808501529590951660a08301525067ffffffffffffffff9290921660c083015290931660e090930192909252949350505050565b60008060008060006080868803121561541057600080fd5b61541986614f58565b9450615427602087016152a0565b935061543560408701614f58565b9250606086013567ffffffffffffffff81111561545157600080fd5b61545d88828901615054565b969995985093965092949392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610180810167ffffffffffffffff811182821017156154c1576154c161546e565b60405290565b60405160e0810167ffffffffffffffff811182821017156154c1576154c161546e565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156155315761553161546e565b604052919050565b803562ffffff81168114614f7c57600080fd5b803561ffff81168114614f7c57600080fd5b6000610180828403121561557157600080fd5b61557961549d565b615582836152a0565b8152615590602084016152a0565b60208201526155a160408401615539565b60408201526155b2606084016152a0565b60608201526155c360808401615539565b60808201526155d460a0840161554c565b60a08201526155e560c084016151dd565b60c08201526155f660e084016152a0565b60e08201526101008381013590820152610120808401359082015261014061561f818501614f58565b90820152610160615631848201614f58565b908201529392505050565b602081016004831061565057615650614fd7565b91905290565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60008282101561569757615697615656565b500390565b600082198211156156af576156af615656565b500190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361571457615714615656565b5060010190565b60006bffffffffffffffffffffffff8381169083168181101561574057615740615656565b039392505050565b60006bffffffffffffffffffffffff80831681851680830382111561576f5761576f615656565b01949350505050565b80518015158114614f7c57600080fd5b60006020828403121561579a57600080fd5b61397282615778565b600181811c908216806157b757607f821691505b6020821081036157f0577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b600081518084526020808501808196508360051b8101915082860160005b8581101561583e57828403895261582c848351614d47565b98850198935090840190600101615814565b5091979650505050505050565b60006060808352858184015260807f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff87111561588657600080fd5b8660051b808983870137808501905081810160008152602083878403018188015281895180845260a093508385019150828b01945060005b8181101561597557855180516bffffffffffffffffffffffff1684528481015173ffffffffffffffffffffffffffffffffffffffff9081168686015260408083015163ffffffff16908601528982015167ffffffffffffffff168a86015288820151168885015285810151615942878601826bffffffffffffffffffffffff169052565b5060c09081015173ffffffffffffffffffffffffffffffffffffffff16908401529483019460e0909201916001016158be565b50508781036040890152615989818a6157f6565b9c9b505050505050505050505050565b6000602082840312156159ab57600080fd5b81516001811061397257600080fd5b6159c48185615006565b6159d16020820184615006565b606060408201526000611d356060830184614d47565b600067ffffffffffffffff821115615a0157615a0161546e565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f830112615a3e57600080fd5b8151615a51615a4c826159e7565b6154ea565b818152846020838601011115615a6657600080fd5b61103b826020830160208701614d17565b600060208284031215615a8957600080fd5b815167ffffffffffffffff811115615aa057600080fd5b61103b84828501615a2d565b600067ffffffffffffffff821115615ac657615ac661546e565b5060051b60200190565b600082601f830112615ae157600080fd5b81356020615af1615a4c83615aac565b82815260e09283028501820192828201919087851115615b1057600080fd5b8387015b85811015615bc05781818a031215615b2c5760008081fd5b615b346154c7565b615b3d826151dd565b8152615b4a868301614f58565b868201526040615b5b8184016152a0565b9082015260608281013567ffffffffffffffff81168114615b7c5760008081fd5b908201526080615b8d838201614f58565b9082015260a0615b9e8382016151dd565b9082015260c0615baf838201614f58565b908201528452928401928101615b14565b5090979650505050505050565b600082601f830112615bde57600080fd5b81356020615bee615a4c83615aac565b82815260059290921b84018101918181019086841115615c0d57600080fd5b8286015b84811015615c8c57803567ffffffffffffffff811115615c315760008081fd5b8701603f81018913615c435760008081fd5b848101356040615c55615a4c836159e7565b8281528b82848601011115615c6a5760008081fd5b8282850189830137600092810188019290925250845250918301918301615c11565b509695505050505050565b600080600060608486031215615cac57600080fd5b833567ffffffffffffffff80821115615cc457600080fd5b818601915086601f830112615cd857600080fd5b81356020615ce8615a4c83615aac565b82815260059290921b8401810191818101908a841115615d0757600080fd5b948201945b83861015615d2557853582529482019490820190615d0c565b97505087013592505080821115615d3b57600080fd5b615d4787838801615ad0565b93506040860135915080821115615d5d57600080fd5b50615d6a86828701615bcd565b9150509250925092565b600060208284031215615d8657600080fd5b5051919050565b8183526000602080850194508260005b85811015615dd65773ffffffffffffffffffffffffffffffffffffffff615dc383614f58565b1687529582019590820190600101615d9d565b509495945050505050565b604081526000615df5604083018688615d8d565b8281036020840152615e08818587615d8d565b979650505050505050565b600060208083526000845481600182811c915080831680615e3557607f831692505b8583108103615e6b577f4e487b710000000000000000000000000000000000000000000000000000000085526022600452602485fd5b878601838152602001818015615e885760018114615eb757615ee2565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00861682528782019650615ee2565b60008b81526020902060005b86811015615edc57815484820152908501908901615ec3565b83019750505b50949998505050505050505050565b60008251615f03818460208701614d17565b9190910192915050565b60008060408385031215615f2057600080fd5b615f2983615778565b9150602083015167ffffffffffffffff811115615f4557600080fd5b615f5185828601615a2d565b9150509250929050565b600063ffffffff808316818103615f7457615f74615656565b6001019392505050565b61018081016107cf8284614da4565b6bffffffffffffffffffffffff8316815260406020820152600061103b6040830184614d47565b805169ffffffffffffffffffff81168114614f7c57600080fd5b600080600080600060a08688031215615fe657600080fd5b615fef86615fb4565b945060208601519350604086015192506060860151915061601260808701615fb4565b90509295509295909350565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561605657616056615656565b500290565b600082616091577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c634300080d000a",
}

var KeeperRegistryABI = KeeperRegistryMetaData.ABI

var KeeperRegistryBin = KeeperRegistryMetaData.Bin

func DeployKeeperRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, link common.Address, linkEthFeed common.Address, fastGasFeed common.Address, config Config) (common.Address, *types.Transaction, *KeeperRegistry, error) {
	parsed, err := KeeperRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistryBin), backend, link, linkEthFeed, fastGasFeed, config)
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
	parsed, err := abi.JSON(strings.NewReader(KeeperRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
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

func (_KeeperRegistry *KeeperRegistryCaller) FASTGASFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "FAST_GAS_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) FASTGASFEED() (common.Address, error) {
	return _KeeperRegistry.Contract.FASTGASFEED(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) FASTGASFEED() (common.Address, error) {
	return _KeeperRegistry.Contract.FASTGASFEED(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) LINK() (common.Address, error) {
	return _KeeperRegistry.Contract.LINK(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) LINK() (common.Address, error) {
	return _KeeperRegistry.Contract.LINK(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCaller) LINKETHFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "LINK_ETH_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) LINKETHFEED() (common.Address, error) {
	return _KeeperRegistry.Contract.LINKETHFEED(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) LINKETHFEED() (common.Address, error) {
	return _KeeperRegistry.Contract.LINKETHFEED(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCaller) GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getActiveUpkeepIDs", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _KeeperRegistry.Contract.GetActiveUpkeepIDs(&_KeeperRegistry.CallOpts, startIndex, maxCount)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _KeeperRegistry.Contract.GetActiveUpkeepIDs(&_KeeperRegistry.CallOpts, startIndex, maxCount)
}

func (_KeeperRegistry *KeeperRegistryCaller) GetKeeperInfo(opts *bind.CallOpts, query common.Address) (GetKeeperInfo,

	error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getKeeperInfo", query)

	outstruct := new(GetKeeperInfo)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Payee = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Active = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.Balance = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_KeeperRegistry *KeeperRegistrySession) GetKeeperInfo(query common.Address) (GetKeeperInfo,

	error) {
	return _KeeperRegistry.Contract.GetKeeperInfo(&_KeeperRegistry.CallOpts, query)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) GetKeeperInfo(query common.Address) (GetKeeperInfo,

	error) {
	return _KeeperRegistry.Contract.GetKeeperInfo(&_KeeperRegistry.CallOpts, query)
}

func (_KeeperRegistry *KeeperRegistryCaller) GetMaxPaymentForGas(opts *bind.CallOpts, gasLimit *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getMaxPaymentForGas", gasLimit)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) GetMaxPaymentForGas(gasLimit *big.Int) (*big.Int, error) {
	return _KeeperRegistry.Contract.GetMaxPaymentForGas(&_KeeperRegistry.CallOpts, gasLimit)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) GetMaxPaymentForGas(gasLimit *big.Int) (*big.Int, error) {
	return _KeeperRegistry.Contract.GetMaxPaymentForGas(&_KeeperRegistry.CallOpts, gasLimit)
}

func (_KeeperRegistry *KeeperRegistryCaller) GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getMinBalanceForUpkeep", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _KeeperRegistry.Contract.GetMinBalanceForUpkeep(&_KeeperRegistry.CallOpts, id)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _KeeperRegistry.Contract.GetMinBalanceForUpkeep(&_KeeperRegistry.CallOpts, id)
}

func (_KeeperRegistry *KeeperRegistryCaller) GetPeerRegistryMigrationPermission(opts *bind.CallOpts, peer common.Address) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getPeerRegistryMigrationPermission", peer)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _KeeperRegistry.Contract.GetPeerRegistryMigrationPermission(&_KeeperRegistry.CallOpts, peer)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _KeeperRegistry.Contract.GetPeerRegistryMigrationPermission(&_KeeperRegistry.CallOpts, peer)
}

func (_KeeperRegistry *KeeperRegistryCaller) GetState(opts *bind.CallOpts) (GetState,

	error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getState")

	outstruct := new(GetState)
	if err != nil {
		return *outstruct, err
	}

	outstruct.State = *abi.ConvertType(out[0], new(State)).(*State)
	outstruct.Config = *abi.ConvertType(out[1], new(Config)).(*Config)
	outstruct.Keepers = *abi.ConvertType(out[2], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

func (_KeeperRegistry *KeeperRegistrySession) GetState() (GetState,

	error) {
	return _KeeperRegistry.Contract.GetState(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) GetState() (GetState,

	error) {
	return _KeeperRegistry.Contract.GetState(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCaller) GetUpkeep(opts *bind.CallOpts, id *big.Int) (GetUpkeep,

	error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "getUpkeep", id)

	outstruct := new(GetUpkeep)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Target = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.ExecuteGas = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.CheckData = *abi.ConvertType(out[2], new([]byte)).(*[]byte)
	outstruct.Balance = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.LastKeeper = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Admin = *abi.ConvertType(out[5], new(common.Address)).(*common.Address)
	outstruct.MaxValidBlocknumber = *abi.ConvertType(out[6], new(uint64)).(*uint64)
	outstruct.AmountSpent = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_KeeperRegistry *KeeperRegistrySession) GetUpkeep(id *big.Int) (GetUpkeep,

	error) {
	return _KeeperRegistry.Contract.GetUpkeep(&_KeeperRegistry.CallOpts, id)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) GetUpkeep(id *big.Int) (GetUpkeep,

	error) {
	return _KeeperRegistry.Contract.GetUpkeep(&_KeeperRegistry.CallOpts, id)
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

func (_KeeperRegistry *KeeperRegistryCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) Paused() (bool, error) {
	return _KeeperRegistry.Contract.Paused(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) Paused() (bool, error) {
	return _KeeperRegistry.Contract.Paused(&_KeeperRegistry.CallOpts)
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

func (_KeeperRegistry *KeeperRegistryCaller) UpkeepTranscoderVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistry.contract.Call(opts, &out, "upkeepTranscoderVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistry *KeeperRegistrySession) UpkeepTranscoderVersion() (uint8, error) {
	return _KeeperRegistry.Contract.UpkeepTranscoderVersion(&_KeeperRegistry.CallOpts)
}

func (_KeeperRegistry *KeeperRegistryCallerSession) UpkeepTranscoderVersion() (uint8, error) {
	return _KeeperRegistry.Contract.UpkeepTranscoderVersion(&_KeeperRegistry.CallOpts)
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

func (_KeeperRegistry *KeeperRegistryTransactor) AcceptPayeeship(opts *bind.TransactOpts, keeper common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "acceptPayeeship", keeper)
}

func (_KeeperRegistry *KeeperRegistrySession) AcceptPayeeship(keeper common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.AcceptPayeeship(&_KeeperRegistry.TransactOpts, keeper)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) AcceptPayeeship(keeper common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.AcceptPayeeship(&_KeeperRegistry.TransactOpts, keeper)
}

func (_KeeperRegistry *KeeperRegistryTransactor) AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "addFunds", id, amount)
}

func (_KeeperRegistry *KeeperRegistrySession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.AddFunds(&_KeeperRegistry.TransactOpts, id, amount)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.AddFunds(&_KeeperRegistry.TransactOpts, id, amount)
}

func (_KeeperRegistry *KeeperRegistryTransactor) CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "cancelUpkeep", id)
}

func (_KeeperRegistry *KeeperRegistrySession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.CancelUpkeep(&_KeeperRegistry.TransactOpts, id)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.CancelUpkeep(&_KeeperRegistry.TransactOpts, id)
}

func (_KeeperRegistry *KeeperRegistryTransactor) CheckUpkeep(opts *bind.TransactOpts, id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "checkUpkeep", id, from)
}

func (_KeeperRegistry *KeeperRegistrySession) CheckUpkeep(id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.CheckUpkeep(&_KeeperRegistry.TransactOpts, id, from)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) CheckUpkeep(id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.CheckUpkeep(&_KeeperRegistry.TransactOpts, id, from)
}

func (_KeeperRegistry *KeeperRegistryTransactor) MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "migrateUpkeeps", ids, destination)
}

func (_KeeperRegistry *KeeperRegistrySession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.MigrateUpkeeps(&_KeeperRegistry.TransactOpts, ids, destination)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.MigrateUpkeeps(&_KeeperRegistry.TransactOpts, ids, destination)
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

func (_KeeperRegistry *KeeperRegistryTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "pause")
}

func (_KeeperRegistry *KeeperRegistrySession) Pause() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.Pause(&_KeeperRegistry.TransactOpts)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) Pause() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.Pause(&_KeeperRegistry.TransactOpts)
}

func (_KeeperRegistry *KeeperRegistryTransactor) PerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "performUpkeep", id, performData)
}

func (_KeeperRegistry *KeeperRegistrySession) PerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.PerformUpkeep(&_KeeperRegistry.TransactOpts, id, performData)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) PerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.PerformUpkeep(&_KeeperRegistry.TransactOpts, id, performData)
}

func (_KeeperRegistry *KeeperRegistryTransactor) ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "receiveUpkeeps", encodedUpkeeps)
}

func (_KeeperRegistry *KeeperRegistrySession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.ReceiveUpkeeps(&_KeeperRegistry.TransactOpts, encodedUpkeeps)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.ReceiveUpkeeps(&_KeeperRegistry.TransactOpts, encodedUpkeeps)
}

func (_KeeperRegistry *KeeperRegistryTransactor) RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "recoverFunds")
}

func (_KeeperRegistry *KeeperRegistrySession) RecoverFunds() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.RecoverFunds(&_KeeperRegistry.TransactOpts)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) RecoverFunds() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.RecoverFunds(&_KeeperRegistry.TransactOpts)
}

func (_KeeperRegistry *KeeperRegistryTransactor) RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "registerUpkeep", target, gasLimit, admin, checkData)
}

func (_KeeperRegistry *KeeperRegistrySession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.RegisterUpkeep(&_KeeperRegistry.TransactOpts, target, gasLimit, admin, checkData)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.RegisterUpkeep(&_KeeperRegistry.TransactOpts, target, gasLimit, admin, checkData)
}

func (_KeeperRegistry *KeeperRegistryTransactor) SetConfig(opts *bind.TransactOpts, config Config) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "setConfig", config)
}

func (_KeeperRegistry *KeeperRegistrySession) SetConfig(config Config) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetConfig(&_KeeperRegistry.TransactOpts, config)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) SetConfig(config Config) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetConfig(&_KeeperRegistry.TransactOpts, config)
}

func (_KeeperRegistry *KeeperRegistryTransactor) SetKeepers(opts *bind.TransactOpts, keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "setKeepers", keepers, payees)
}

func (_KeeperRegistry *KeeperRegistrySession) SetKeepers(keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetKeepers(&_KeeperRegistry.TransactOpts, keepers, payees)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) SetKeepers(keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetKeepers(&_KeeperRegistry.TransactOpts, keepers, payees)
}

func (_KeeperRegistry *KeeperRegistryTransactor) SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "setPeerRegistryMigrationPermission", peer, permission)
}

func (_KeeperRegistry *KeeperRegistrySession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetPeerRegistryMigrationPermission(&_KeeperRegistry.TransactOpts, peer, permission)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetPeerRegistryMigrationPermission(&_KeeperRegistry.TransactOpts, peer, permission)
}

func (_KeeperRegistry *KeeperRegistryTransactor) SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "setUpkeepGasLimit", id, gasLimit)
}

func (_KeeperRegistry *KeeperRegistrySession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetUpkeepGasLimit(&_KeeperRegistry.TransactOpts, id, gasLimit)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.SetUpkeepGasLimit(&_KeeperRegistry.TransactOpts, id, gasLimit)
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

func (_KeeperRegistry *KeeperRegistryTransactor) TransferPayeeship(opts *bind.TransactOpts, keeper common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "transferPayeeship", keeper, proposed)
}

func (_KeeperRegistry *KeeperRegistrySession) TransferPayeeship(keeper common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.TransferPayeeship(&_KeeperRegistry.TransactOpts, keeper, proposed)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) TransferPayeeship(keeper common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.TransferPayeeship(&_KeeperRegistry.TransactOpts, keeper, proposed)
}

func (_KeeperRegistry *KeeperRegistryTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "unpause")
}

func (_KeeperRegistry *KeeperRegistrySession) Unpause() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.Unpause(&_KeeperRegistry.TransactOpts)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) Unpause() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.Unpause(&_KeeperRegistry.TransactOpts)
}

func (_KeeperRegistry *KeeperRegistryTransactor) WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "withdrawFunds", id, to)
}

func (_KeeperRegistry *KeeperRegistrySession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.WithdrawFunds(&_KeeperRegistry.TransactOpts, id, to)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.WithdrawFunds(&_KeeperRegistry.TransactOpts, id, to)
}

func (_KeeperRegistry *KeeperRegistryTransactor) WithdrawOwnerFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "withdrawOwnerFunds")
}

func (_KeeperRegistry *KeeperRegistrySession) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.WithdrawOwnerFunds(&_KeeperRegistry.TransactOpts)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _KeeperRegistry.Contract.WithdrawOwnerFunds(&_KeeperRegistry.TransactOpts)
}

func (_KeeperRegistry *KeeperRegistryTransactor) WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.contract.Transact(opts, "withdrawPayment", from, to)
}

func (_KeeperRegistry *KeeperRegistrySession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.WithdrawPayment(&_KeeperRegistry.TransactOpts, from, to)
}

func (_KeeperRegistry *KeeperRegistryTransactorSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry.Contract.WithdrawPayment(&_KeeperRegistry.TransactOpts, from, to)
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
	Config Config
	Raw    types.Log
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

type KeeperRegistryKeepersUpdatedIterator struct {
	Event *KeeperRegistryKeepersUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryKeepersUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryKeepersUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryKeepersUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryKeepersUpdatedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryKeepersUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryKeepersUpdated struct {
	Keepers []common.Address
	Payees  []common.Address
	Raw     types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterKeepersUpdated(opts *bind.FilterOpts) (*KeeperRegistryKeepersUpdatedIterator, error) {

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "KeepersUpdated")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryKeepersUpdatedIterator{contract: _KeeperRegistry.contract, event: "KeepersUpdated", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchKeepersUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryKeepersUpdated) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "KeepersUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryKeepersUpdated)
				if err := _KeeperRegistry.contract.UnpackLog(event, "KeepersUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) ParseKeepersUpdated(log types.Log) (*KeeperRegistryKeepersUpdated, error) {
	event := new(KeeperRegistryKeepersUpdated)
	if err := _KeeperRegistry.contract.UnpackLog(event, "KeepersUpdated", log); err != nil {
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
	Keeper common.Address
	From   common.Address
	To     common.Address
	Raw    types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, keeper []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryPayeeshipTransferRequestedIterator, error) {

	var keeperRule []interface{}
	for _, keeperItem := range keeper {
		keeperRule = append(keeperRule, keeperItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "PayeeshipTransferRequested", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryPayeeshipTransferRequestedIterator{contract: _KeeperRegistry.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPayeeshipTransferRequested, keeper []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var keeperRule []interface{}
	for _, keeperItem := range keeper {
		keeperRule = append(keeperRule, keeperItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "PayeeshipTransferRequested", keeperRule, fromRule, toRule)
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
	Keeper common.Address
	From   common.Address
	To     common.Address
	Raw    types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, keeper []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryPayeeshipTransferredIterator, error) {

	var keeperRule []interface{}
	for _, keeperItem := range keeper {
		keeperRule = append(keeperRule, keeperItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "PayeeshipTransferred", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryPayeeshipTransferredIterator{contract: _KeeperRegistry.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPayeeshipTransferred, keeper []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var keeperRule []interface{}
	for _, keeperItem := range keeper {
		keeperRule = append(keeperRule, keeperItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "PayeeshipTransferred", keeperRule, fromRule, toRule)
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
	Keeper common.Address
	Amount *big.Int
	To     common.Address
	Payee  common.Address
	Raw    types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, keeper []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistryPaymentWithdrawnIterator, error) {

	var keeperRule []interface{}
	for _, keeperItem := range keeper {
		keeperRule = append(keeperRule, keeperItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "PaymentWithdrawn", keeperRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryPaymentWithdrawnIterator{contract: _KeeperRegistry.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPaymentWithdrawn, keeper []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

	var keeperRule []interface{}
	for _, keeperItem := range keeper {
		keeperRule = append(keeperRule, keeperItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "PaymentWithdrawn", keeperRule, amountRule, toRule)
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
	Id          *big.Int
	Success     bool
	From        common.Address
	Payment     *big.Int
	PerformData []byte
	Raw         types.Log
}

func (_KeeperRegistry *KeeperRegistryFilterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool, from []common.Address) (*KeeperRegistryUpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistry.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryUpkeepPerformedIterator{contract: _KeeperRegistry.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry *KeeperRegistryFilterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepPerformed, id []*big.Int, success []bool, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistry.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule, fromRule)
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

type GetKeeperInfo struct {
	Payee   common.Address
	Active  bool
	Balance *big.Int
}
type GetState struct {
	State   State
	Config  Config
	Keepers []common.Address
}
type GetUpkeep struct {
	Target              common.Address
	ExecuteGas          uint32
	CheckData           []byte
	Balance             *big.Int
	LastKeeper          common.Address
	Admin               common.Address
	MaxValidBlocknumber uint64
	AmountSpent         *big.Int
}

func (_KeeperRegistry *KeeperRegistry) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _KeeperRegistry.abi.Events["ConfigSet"].ID:
		return _KeeperRegistry.ParseConfigSet(log)
	case _KeeperRegistry.abi.Events["FundsAdded"].ID:
		return _KeeperRegistry.ParseFundsAdded(log)
	case _KeeperRegistry.abi.Events["FundsWithdrawn"].ID:
		return _KeeperRegistry.ParseFundsWithdrawn(log)
	case _KeeperRegistry.abi.Events["KeepersUpdated"].ID:
		return _KeeperRegistry.ParseKeepersUpdated(log)
	case _KeeperRegistry.abi.Events["OwnerFundsWithdrawn"].ID:
		return _KeeperRegistry.ParseOwnerFundsWithdrawn(log)
	case _KeeperRegistry.abi.Events["OwnershipTransferRequested"].ID:
		return _KeeperRegistry.ParseOwnershipTransferRequested(log)
	case _KeeperRegistry.abi.Events["OwnershipTransferred"].ID:
		return _KeeperRegistry.ParseOwnershipTransferred(log)
	case _KeeperRegistry.abi.Events["Paused"].ID:
		return _KeeperRegistry.ParsePaused(log)
	case _KeeperRegistry.abi.Events["PayeeshipTransferRequested"].ID:
		return _KeeperRegistry.ParsePayeeshipTransferRequested(log)
	case _KeeperRegistry.abi.Events["PayeeshipTransferred"].ID:
		return _KeeperRegistry.ParsePayeeshipTransferred(log)
	case _KeeperRegistry.abi.Events["PaymentWithdrawn"].ID:
		return _KeeperRegistry.ParsePaymentWithdrawn(log)
	case _KeeperRegistry.abi.Events["Unpaused"].ID:
		return _KeeperRegistry.ParseUnpaused(log)
	case _KeeperRegistry.abi.Events["UpkeepCanceled"].ID:
		return _KeeperRegistry.ParseUpkeepCanceled(log)
	case _KeeperRegistry.abi.Events["UpkeepGasLimitSet"].ID:
		return _KeeperRegistry.ParseUpkeepGasLimitSet(log)
	case _KeeperRegistry.abi.Events["UpkeepMigrated"].ID:
		return _KeeperRegistry.ParseUpkeepMigrated(log)
	case _KeeperRegistry.abi.Events["UpkeepPerformed"].ID:
		return _KeeperRegistry.ParseUpkeepPerformed(log)
	case _KeeperRegistry.abi.Events["UpkeepReceived"].ID:
		return _KeeperRegistry.ParseUpkeepReceived(log)
	case _KeeperRegistry.abi.Events["UpkeepRegistered"].ID:
		return _KeeperRegistry.ParseUpkeepRegistered(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeeperRegistryConfigSet) Topic() common.Hash {
	return common.HexToHash("0xfe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de325")
}

func (KeeperRegistryFundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (KeeperRegistryFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (KeeperRegistryKeepersUpdated) Topic() common.Hash {
	return common.HexToHash("0x056264c94f28bb06c99d13f0446eb96c67c215d8d707bce2655a98ddf1c0b71f")
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

func (KeeperRegistryPayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (KeeperRegistryPayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (KeeperRegistryPaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (KeeperRegistryUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (KeeperRegistryUpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (KeeperRegistryUpkeepGasLimitSet) Topic() common.Hash {
	return common.HexToHash("0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c")
}

func (KeeperRegistryUpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (KeeperRegistryUpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xcaacad83e47cc45c280d487ec84184eee2fa3b54ebaa393bda7549f13da228f6")
}

func (KeeperRegistryUpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (KeeperRegistryUpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (_KeeperRegistry *KeeperRegistry) Address() common.Address {
	return _KeeperRegistry.address
}

type KeeperRegistryInterface interface {
	FASTGASFEED(opts *bind.CallOpts) (common.Address, error)

	LINK(opts *bind.CallOpts) (common.Address, error)

	LINKETHFEED(opts *bind.CallOpts) (common.Address, error)

	GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)

	GetKeeperInfo(opts *bind.CallOpts, query common.Address) (GetKeeperInfo,

		error)

	GetMaxPaymentForGas(opts *bind.CallOpts, gasLimit *big.Int) (*big.Int, error)

	GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error)

	GetPeerRegistryMigrationPermission(opts *bind.CallOpts, peer common.Address) (uint8, error)

	GetState(opts *bind.CallOpts) (GetState,

		error)

	GetUpkeep(opts *bind.CallOpts, id *big.Int) (GetUpkeep,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	Paused(opts *bind.CallOpts) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	UpkeepTranscoderVersion(opts *bind.CallOpts) (uint8, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptPayeeship(opts *bind.TransactOpts, keeper common.Address) (*types.Transaction, error)

	AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error)

	CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	CheckUpkeep(opts *bind.TransactOpts, id *big.Int, from common.Address) (*types.Transaction, error)

	MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	Pause(opts *bind.TransactOpts) (*types.Transaction, error)

	PerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error)

	ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error)

	RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error)

	RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, config Config) (*types.Transaction, error)

	SetKeepers(opts *bind.TransactOpts, keepers []common.Address, payees []common.Address) (*types.Transaction, error)

	SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error)

	SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	TransferPayeeship(opts *bind.TransactOpts, keeper common.Address, proposed common.Address) (*types.Transaction, error)

	Unpause(opts *bind.TransactOpts) (*types.Transaction, error)

	WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error)

	WithdrawOwnerFunds(opts *bind.TransactOpts) (*types.Transaction, error)

	WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*KeeperRegistryConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*KeeperRegistryConfigSet, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistryFundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*KeeperRegistryFundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryFundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryFundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*KeeperRegistryFundsWithdrawn, error)

	FilterKeepersUpdated(opts *bind.FilterOpts) (*KeeperRegistryKeepersUpdatedIterator, error)

	WatchKeepersUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryKeepersUpdated) (event.Subscription, error)

	ParseKeepersUpdated(log types.Log) (*KeeperRegistryKeepersUpdated, error)

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

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, keeper []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryPayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPayeeshipTransferRequested, keeper []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistryPayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, keeper []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryPayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPayeeshipTransferred, keeper []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*KeeperRegistryPayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, keeper []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistryPaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryPaymentWithdrawn, keeper []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*KeeperRegistryPaymentWithdrawn, error)

	FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistryUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*KeeperRegistryUnpaused, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistryUpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*KeeperRegistryUpkeepCanceled, error)

	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepGasLimitSetIterator, error)

	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepGasLimitSet(log types.Log) (*KeeperRegistryUpkeepGasLimitSet, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*KeeperRegistryUpkeepMigrated, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool, from []common.Address) (*KeeperRegistryUpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepPerformed, id []*big.Int, success []bool, from []common.Address) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*KeeperRegistryUpkeepPerformed, error)

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*KeeperRegistryUpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryUpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistryUpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*KeeperRegistryUpkeepRegistered, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
