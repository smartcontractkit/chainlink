// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package keeper_registry_wrapper1_3

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
	Nonce               uint32
	OwnerLinkBalance    *big.Int
	ExpectedLinkBalance *big.Int
	NumUpkeeps          *big.Int
}

var KeeperRegistry13MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"keeperRegistryLogic\",\"type\":\"address\",\"internalType\":\"contractKeeperRegistryLogic1_3\"},{\"name\":\"config\",\"type\":\"tuple\",\"internalType\":\"structConfig\",\"components\":[{\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockCountPerTurn\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"checkGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"stalenessSeconds\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"minUpkeepSpend\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"maxPerformGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"fallbackGasPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"transcoder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"registrar\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"ARB_NITRO_ORACLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractArbGasInfo\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"FAST_GAS_FEED\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractAggregatorV3Interface\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"KEEPER_REGISTRY_LOGIC\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"LINK\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractLinkTokenInterface\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"LINK_ETH_FEED\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractAggregatorV3Interface\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"OPTIMISM_ORACLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractOVM_GasPriceOracle\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PAYMENT_MODEL\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumKeeperRegistryBase1_3.PaymentModel\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"REGISTRY_GAS_OVERHEAD\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptPayeeship\",\"inputs\":[{\"name\":\"keeper\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptUpkeepAdmin\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addFunds\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"checkUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"performData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"maxLinkPayment\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"adjustedGasWei\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"linkEth\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getActiveUpkeepIDs\",\"inputs\":[{\"name\":\"startIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getKeeperInfo\",\"inputs\":[{\"name\":\"query\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"payee\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"active\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"balance\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMaxPaymentForGas\",\"inputs\":[{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"maxPayment\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinBalanceForUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"minBalance\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPeerRegistryMigrationPermission\",\"inputs\":[{\"name\":\"peer\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumKeeperRegistryBase1_3.MigrationPermission\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getState\",\"inputs\":[],\"outputs\":[{\"name\":\"state\",\"type\":\"tuple\",\"internalType\":\"structState\",\"components\":[{\"name\":\"nonce\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"ownerLinkBalance\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"expectedLinkBalance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"numUpkeeps\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"config\",\"type\":\"tuple\",\"internalType\":\"structConfig\",\"components\":[{\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockCountPerTurn\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"checkGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"stalenessSeconds\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"minUpkeepSpend\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"maxPerformGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"fallbackGasPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"transcoder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"registrar\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"keepers\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"executeGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"checkData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"balance\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"lastKeeper\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"amountSpent\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"paused\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"migrateUpkeeps\",\"inputs\":[{\"name\":\"ids\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"destination\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onTokenTransfer\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pauseUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"performUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"performData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"success\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"receiveUpkeeps\",\"inputs\":[{\"name\":\"encodedUpkeeps\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"recoverFunds\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerUpkeep\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"checkData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setConfig\",\"inputs\":[{\"name\":\"config\",\"type\":\"tuple\",\"internalType\":\"structConfig\",\"components\":[{\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockCountPerTurn\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"checkGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"stalenessSeconds\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"minUpkeepSpend\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"maxPerformGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"fallbackGasPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"transcoder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"registrar\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setKeepers\",\"inputs\":[{\"name\":\"keepers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"payees\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPeerRegistryMigrationPermission\",\"inputs\":[{\"name\":\"peer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permission\",\"type\":\"uint8\",\"internalType\":\"enumKeeperRegistryBase1_3.MigrationPermission\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUpkeepGasLimit\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferPayeeship\",\"inputs\":[{\"name\":\"keeper\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"proposed\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferUpkeepAdmin\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proposed\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpauseUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateCheckData\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"newCheckData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upkeepTranscoderVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumUpkeepFormat\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdrawFunds\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawOwnerFunds\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawPayment\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"ConfigSet\",\"inputs\":[{\"name\":\"config\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structConfig\",\"components\":[{\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockCountPerTurn\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"checkGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"stalenessSeconds\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"minUpkeepSpend\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"maxPerformGas\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"fallbackGasPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"transcoder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"registrar\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsAdded\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsWithdrawn\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"KeepersUpdated\",\"inputs\":[{\"name\":\"keepers\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"payees\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnerFundsWithdrawn\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeeshipTransferRequested\",\"inputs\":[{\"name\":\"keeper\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeeshipTransferred\",\"inputs\":[{\"name\":\"keeper\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PaymentWithdrawn\",\"inputs\":[{\"name\":\"keeper\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"payee\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepAdminTransferRequested\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepAdminTransferred\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepCanceled\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"atBlockHeight\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepCheckDataUpdated\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"newCheckData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepGasLimitSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"gasLimit\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepMigrated\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"remainingBalance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"destination\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPaused\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPerformed\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"success\",\"type\":\"bool\",\"indexed\":true,\"internalType\":\"bool\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"payment\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"},{\"name\":\"performData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepReceived\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"startingBalance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"importedFrom\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepRegistered\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"executeGas\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"admin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepUnpaused\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ArrayHasNoEntries\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotCancel\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateEntry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EmptyAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GasLimitCanOnlyIncrease\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GasLimitOutsideRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IndexOutOfRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientFunds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidDataLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRecipient\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"KeepersMustTakeTurns\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MigrationNotPermitted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAContract\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyActiveKeepers\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByLINKToken\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwnerOrAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByProposedAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByProposedPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyPausedUpkeep\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlySimulatedBackend\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyUnpausedUpkeep\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterLengthError\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PaymentGreaterThanAllLINK\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TargetCheckReverted\",\"inputs\":[{\"name\":\"reason\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"TranscoderNotSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepCancelled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepNotCanceled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepNotNeeded\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ValueNotChanged\",\"inputs\":[]}]",
	Bin: "0x6101806040527f420000000000000000000000000000000000000f00000000000000000000000060e0526c6c000000000000000000000000610100523480156200004857600080fd5b50604051620048cd380380620048cd8339810160408190526200006b91620008a7565b816001600160a01b0316638811cbe86040518163ffffffff1660e01b815260040160206040518083038186803b158015620000a557600080fd5b505afa158015620000ba573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620000e09190620009ce565b826001600160a01b0316635077b2106040518163ffffffff1660e01b815260040160206040518083038186803b1580156200011a57600080fd5b505afa1580156200012f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001559190620009f1565b836001600160a01b0316631b6b6d236040518163ffffffff1660e01b815260040160206040518083038186803b1580156200018f57600080fd5b505afa158015620001a4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001ca919062000880565b846001600160a01b031663ad1783616040518163ffffffff1660e01b815260040160206040518083038186803b1580156200020457600080fd5b505afa15801562000219573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200023f919062000880565b856001600160a01b0316634584a4196040518163ffffffff1660e01b815260040160206040518083038186803b1580156200027957600080fd5b505afa1580156200028e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002b4919062000880565b33806000816200030b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200033e576200033e816200041b565b5050600160029081556003805460ff1916905586915081111562000366576200036662000b42565b6101208160028111156200037e576200037e62000b42565b60f81b9052506101408490526001600160a01b0383161580620003a857506001600160a01b038216155b80620003bb57506001600160a01b038116155b15620003da57604051637138356f60e01b815260040160405180910390fd5b6001600160601b0319606093841b811660805291831b821660a052821b811660c0529085901b16610160525062000413905081620004c7565b505062000b71565b6001600160a01b038116331415620004765760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000302565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b620004d1620007bc565b600e5460e082015163ffffffff918216911610156200050357604051630e6af04160e21b815260040160405180910390fd5b604051806101200160405280826000015163ffffffff168152602001826020015163ffffffff168152602001826040015162ffffff168152602001826060015163ffffffff168152602001826080015162ffffff1681526020018260a0015161ffff1681526020018260c001516001600160601b031681526020018260e0015163ffffffff168152602001600d60010160049054906101000a900463ffffffff1663ffffffff16815250600d60008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160086101000a81548162ffffff021916908362ffffff160217905550606082015181600001600b6101000a81548163ffffffff021916908363ffffffff160217905550608082015181600001600f6101000a81548162ffffff021916908362ffffff16021790555060a08201518160000160126101000a81548161ffff021916908361ffff16021790555060c08201518160000160146101000a8154816001600160601b0302191690836001600160601b0316021790555060e08201518160010160006101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160010160046101000a81548163ffffffff021916908363ffffffff160217905550905050806101000151600f81905550806101200151601081905550806101400151601360006101000a8154816001600160a01b0302191690836001600160a01b03160217905550806101600151601460006101000a8154816001600160a01b0302191690836001600160a01b031602179055507ffe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de32581604051620007b1919062000a0b565b60405180910390a150565b6000546001600160a01b03163314620008185760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000302565b565b8051620008278162000b58565b919050565b805161ffff811681146200082757600080fd5b805162ffffff811681146200082757600080fd5b805163ffffffff811681146200082757600080fd5b80516001600160601b03811681146200082757600080fd5b6000602082840312156200089357600080fd5b8151620008a08162000b58565b9392505050565b6000808284036101a0811215620008bd57600080fd5b8351620008ca8162000b58565b9250610180601f198201811315620008e157600080fd5b620008eb62000b0a565b9150620008fb6020860162000853565b82526200090b6040860162000853565b60208301526200091e606086016200083f565b6040830152620009316080860162000853565b60608301526200094460a086016200083f565b60808301526200095760c086016200082c565b60a08301526200096a60e0860162000868565b60c08301526101006200097f81870162000853565b60e084015261012080870151828501526101409150818701518185015250610160620009ad8188016200081a565b82850152620009be8388016200081a565b9084015250929590945092505050565b600060208284031215620009e157600080fd5b815160038110620008a057600080fd5b60006020828403121562000a0457600080fd5b5051919050565b815163ffffffff1681526101808101602083015162000a32602084018263ffffffff169052565b50604083015162000a4a604084018262ffffff169052565b50606083015162000a63606084018263ffffffff169052565b50608083015162000a7b608084018262ffffff169052565b5060a083015162000a9260a084018261ffff169052565b5060c083015162000aae60c08401826001600160601b03169052565b5060e083015162000ac760e084018263ffffffff169052565b5061010083810151908301526101208084015190830152610140808401516001600160a01b03908116918401919091526101609384015116929091019190915290565b60405161018081016001600160401b038111828210171562000b3c57634e487b7160e01b600052604160045260246000fd5b60405290565b634e487b7160e01b600052602160045260246000fd5b6001600160a01b038116811462000b6e57600080fd5b50565b60805160601c60a05160601c60c05160601c60e05160601c6101005160601c6101205160f81c610140516101605160601c613ca262000c2b600039600081816103600152610a4c0152600081816105e601526125460152600081816107490152818161259901526127150152600081816106a1015261274d0152600081816106d50152612684015260008181610590015261227a015260008181610891015261235b01526000818161046e015261144e0152613ca26000f3fe6080604052600436106103015760003560e01c80638811cbe81161018f578063b148ab6b116100e1578063c80480221161008a578063ef47a0ce11610064578063ef47a0ce146109b4578063f2fde38b146109d4578063faa3e996146109f457610310565b8063c8048022146108d3578063da5c674114610994578063eb5dcd6c1461084957610310565b8063b7fdb436116100bb578063b7fdb4361461090e578063c41b813a1461092e578063c7c3a19a1461095f57610310565b8063b148ab6b146108d3578063b657bc9c146108ee578063b79550be1461056957610310565b80639fab438611610143578063a72aa27e1161011d578063a72aa27e14610864578063ad1783611461087f578063b121e147146108b357610310565b80639fab438614610809578063a4c0ed3614610829578063a710b2211461084957610310565b80638e86139b116101745780638e86139b1461079657806393f0c1fc146107b1578063948108f7146107ee57610310565b80638811cbe8146107375780638da5cb5b1461076b57610310565b80635077b210116102535780637d9b97e0116101fc578063850cce34116101d6578063850cce34146106c357806385c1b0ba146106f75780638765ecbe1461071757610310565b80637d9b97e0146105695780637f37618e1461068f5780638456cb591461056957610310565b8063744bfe611161022d578063744bfe611461044157806379ba50971461065a5780637bbaf1ea1461066f57610310565b80635077b210146105d45780635165f2f5146106165780635c975abb1461063657610310565b80631a2af011116102b55780633f4ba83a1161028f5780633f4ba83a146105695780634584a4191461057e57806348013d7b146105b257610310565b80631a2af011146104415780631b6b6d231461045c5780631e12b8a51461049057610310565b8063181f5a77116102e6578063181f5a77146103a75780631865c57d146103fd578063187256e81461042157610310565b806306e3b632146103185780630852c7c91461034e57610310565b366103105761030e610a47565b005b61030e610a47565b34801561032457600080fd5b50610338610333366004613383565b610a72565b6040516103459190613655565b60405180910390f35b34801561035a57600080fd5b506103827f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610345565b3480156103b357600080fd5b506103f06040518060400160405280601481526020017f4b6565706572526567697374727920312e332e3000000000000000000000000081525081565b60405161034591906136e6565b34801561040957600080fd5b50610412610b6e565b60405161034593929190613766565b34801561042d57600080fd5b5061030e61043c366004613003565b610e26565b34801561044d57600080fd5b5061030e61043c366004613314565b34801561046857600080fd5b506103827f000000000000000000000000000000000000000000000000000000000000000081565b34801561049c57600080fd5b506105296104ab366004612fb5565b73ffffffffffffffffffffffffffffffffffffffff90811660009081526008602090815260409182902082516060810184528154948516808252740100000000000000000000000000000000000000009095046bffffffffffffffffffffffff1692810183905260019091015460ff16151592018290529192909190565b6040805173ffffffffffffffffffffffffffffffffffffffff909416845291151560208401526bffffffffffffffffffffffff1690820152606001610345565b34801561057557600080fd5b5061030e610e32565b34801561058a57600080fd5b506103827f000000000000000000000000000000000000000000000000000000000000000081565b3480156105be57600080fd5b506105c7600181565b604051610345919061374a565b3480156105e057600080fd5b506106087f000000000000000000000000000000000000000000000000000000000000000081565b604051908152602001610345565b34801561062257600080fd5b5061030e6106313660046132e2565b610e3a565b34801561064257600080fd5b5060035460ff165b6040519015158152602001610345565b34801561066657600080fd5b5061030e610fc6565b34801561067b57600080fd5b5061064a61068a366004613337565b6110c8565b34801561069b57600080fd5b506103827f000000000000000000000000000000000000000000000000000000000000000081565b3480156106cf57600080fd5b506103827f000000000000000000000000000000000000000000000000000000000000000081565b34801561070357600080fd5b5061030e61071236600461316e565b611126565b34801561072357600080fd5b5061030e6107323660046132e2565b611133565b34801561074357600080fd5b506105c77f000000000000000000000000000000000000000000000000000000000000000081565b34801561077757600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff16610382565b3480156107a257600080fd5b5061030e61043c3660046131c2565b3480156107bd57600080fd5b506107d16107cc3660046132e2565b6112df565b6040516bffffffffffffffffffffffff9091168152602001610345565b3480156107fa57600080fd5b5061030e61043c3660046133c8565b34801561081557600080fd5b5061030e610824366004613337565b6112fd565b34801561083557600080fd5b5061030e61084436600461303e565b611436565b34801561085557600080fd5b5061030e61043c366004612fd0565b34801561087057600080fd5b5061030e61043c3660046133a5565b34801561088b57600080fd5b506103827f000000000000000000000000000000000000000000000000000000000000000081565b3480156108bf57600080fd5b5061030e6108ce366004612fb5565b61162d565b3480156108df57600080fd5b5061030e6108ce3660046132e2565b3480156108fa57600080fd5b506107d16109093660046132e2565b611638565b34801561091a57600080fd5b5061030e61092936600461310e565b611659565b34801561093a57600080fd5b5061094e610949366004613314565b611667565b6040516103459594939291906136f9565b34801561096b57600080fd5b5061097f61097a3660046132e2565b611689565b604051610345999897969594939291906135c3565b3480156109a057600080fd5b506106086109af366004613098565b611859565b3480156109c057600080fd5b5061030e6109cf366004613204565b61186c565b3480156109e057600080fd5b5061030e6109ef366004612fb5565b611bb8565b348015610a0057600080fd5b50610a3a610a0f366004612fb5565b73ffffffffffffffffffffffffffffffffffffffff166000908152600c602052604090205460ff1690565b6040516103459190613730565b610a707f0000000000000000000000000000000000000000000000000000000000000000611bc9565b565b60606000610a806005611bed565b9050808410610abb576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b82610acd57610aca8482613945565b92505b60008367ffffffffffffffff811115610ae857610ae8613afe565b604051908082528060200260200182016040528015610b11578160200160208202803683370190505b50905060005b84811015610b6357610b34610b2c8288613885565b600590611bf7565b828281518110610b4657610b46613acf565b602090810291909101015280610b5b81613a09565b915050610b17565b509150505b92915050565b6040805160808101825260008082526020820181905291810182905260608101919091526040805161018081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290526101008101829052610120810182905261014081018290526101608101919091526040805161012081018252600d5463ffffffff8082168352640100000000808304821660208086019190915262ffffff6801000000000000000085048116868801526b010000000000000000000000850484166060878101919091526f010000000000000000000000000000008604909116608087015261ffff720100000000000000000000000000000000000086041660a08701526bffffffffffffffffffffffff74010000000000000000000000000000000000000000909504851660c0870152600e5480851660e088015292909204909216610100850181905287526011549092169086015260125492850192909252610cf06005611bed565b606080860191909152815163ffffffff908116855260208084015182168187015260408085015162ffffff90811682890152858501518416948801949094526080808601519094169387019390935260a08085015161ffff169087015260c0808501516bffffffffffffffffffffffff169087015260e08085015190921691860191909152600f5461010086015260105461012086015260135473ffffffffffffffffffffffffffffffffffffffff90811661014087015260145416610160860152600480548351818402810184019094528084528793879390918391830182828015610e1357602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610de8575b5050505050905093509350935050909192565b610e2e610a47565b5050565b610a70610a47565b60008181526007602090815260409182902082516101008101845281546bffffffffffffffffffffffff808216835273ffffffffffffffffffffffffffffffffffffffff6c0100000000000000000000000092839004811695840195909552600184015490811695830195909552909304821660608401526002015463ffffffff808216608085015264010000000082041660a084015268010000000000000000810490911660c083015260ff7c010000000000000000000000000000000000000000000000000000000090910416151560e0820152610f1981611c0a565b8060e00151610f54576040517f1b88a78400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260076020526040902060020180547fffffff00ffffffffffffffffffffffffffffffffffffffffffffffffffffffff169055610f96600583611cb7565b5060405182907f7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a4745690600090a25050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461104c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60006110d2611cc3565b61111e611119338686868080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525060019250611d30915050565b611e1a565b949350505050565b61112e610a47565b505050565b60008181526007602090815260409182902082516101008101845281546bffffffffffffffffffffffff808216835273ffffffffffffffffffffffffffffffffffffffff6c0100000000000000000000000092839004811695840195909552600184015490811695830195909552909304821660608401526002015463ffffffff808216608085015264010000000082041660a084015268010000000000000000810490911660c083015260ff7c010000000000000000000000000000000000000000000000000000000090910416151560e082015261121281611c0a565b8060e001511561124e576040517f514b6c2400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260076020526040902060020180547fffffff00ffffffffffffffffffffffffffffffffffffffffffffffffffffffff167c01000000000000000000000000000000000000000000000000000000001790556112af60058361223b565b5060405182907f8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f90600090a25050565b60008060006112ec612247565b9150915061111e8483836000612442565b60008381526007602090815260409182902082516101008101845281546bffffffffffffffffffffffff808216835273ffffffffffffffffffffffffffffffffffffffff6c0100000000000000000000000092839004811695840195909552600184015490811695830195909552909304821660608401526002015463ffffffff808216608085015264010000000082041660a084015268010000000000000000810490911660c083015260ff7c010000000000000000000000000000000000000000000000000000000090910416151560e08201526113dc81611c0a565b6000848152600b602052604090206113f5908484612ddd565b50837f7b778136e5211932b51a145badd01959415e79e051a933604b3d323f862dcabf8484604051611428929190613699565b60405180910390a250505050565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146114a5576040517fc8bad78d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b602081146114df576040517fdfe9309000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006114ed828401846132e2565b600081815260076020526040902060020154909150640100000000900463ffffffff90811614611549576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000818152600760205260409020546115719085906bffffffffffffffffffffffff1661389d565b600082815260076020526040902080547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff929092169190911790556012546115c8908590613885565b6012556040516bffffffffffffffffffffffff8516815273ffffffffffffffffffffffffffffffffffffffff86169082907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a35050505050565b611635610a47565b50565b600081815260076020526040812060020154610b689063ffffffff166112df565b611661610a47565b50505050565b60606000806000806116776128b3565b61167f610a47565b9295509295909350565b600081815260076020908152604080832081516101008101835281546bffffffffffffffffffffffff80821683526c010000000000000000000000009182900473ffffffffffffffffffffffffffffffffffffffff9081168488019081526001860154928316858801908152939092048116606080860191825260029096015463ffffffff80821660808801819052640100000000830490911660a0880190815268010000000000000000830490941660c088018190527c010000000000000000000000000000000000000000000000000000000090920460ff16151560e088019081528c8c52600b909a52978a2086519451925193519551995181548c9b999a8c9a8b9a8b9a8b9a8b9a8b9a93999895979596919593949093909187906117b0906139b5565b80601f01602080910402602001604051908101604052809291908181526020018280546117dc906139b5565b80156118295780601f106117fe57610100808354040283529160200191611829565b820191906000526020600020905b81548152906001019060200180831161180c57829003601f168201915b505050505096508263ffffffff169250995099509950995099509950995099509950509193959799909294969850565b6000611863610a47565b95945050505050565b6118746128eb565b600e5460e082015163ffffffff918216911610156118be576040517f39abc10400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604051806101200160405280826000015163ffffffff168152602001826020015163ffffffff168152602001826040015162ffffff168152602001826060015163ffffffff168152602001826080015162ffffff1681526020018260a0015161ffff1681526020018260c001516bffffffffffffffffffffffff1681526020018260e0015163ffffffff168152602001600d60010160049054906101000a900463ffffffff1663ffffffff16815250600d60008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160086101000a81548162ffffff021916908362ffffff160217905550606082015181600001600b6101000a81548163ffffffff021916908363ffffffff160217905550608082015181600001600f6101000a81548162ffffff021916908362ffffff16021790555060a08201518160000160126101000a81548161ffff021916908361ffff16021790555060c08201518160000160146101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060e08201518160010160006101000a81548163ffffffff021916908363ffffffff1602179055506101008201518160010160046101000a81548163ffffffff021916908363ffffffff160217905550905050806101000151600f81905550806101200151601081905550806101400151601360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550806101600151601460006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507ffe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de32581604051611bad9190613757565b60405180910390a150565b611bc06128eb565b6116358161296c565b3660008037600080366000845af43d6000803e808015611be8573d6000f35b3d6000fd5b6000610b68825490565b6000611c038383612a62565b9392505050565b806060015173ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611c73576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60a081015163ffffffff90811614611635576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000611c038383612a8c565b60035460ff1615610a70576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a20706175736564000000000000000000000000000000006044820152606401611043565b611d866040518060e00160405280600073ffffffffffffffffffffffffffffffffffffffff1681526020016000815260200160608152602001600081526020016000815260200160008152602001600081525090565b60008481526007602052604081206002015463ffffffff169080611da8612247565b915091506000611dba84848489612442565b6040805160e08101825273ffffffffffffffffffffffffffffffffffffffff909b168b5260208b0199909952978901969096526bffffffffffffffffffffffff9096166060880152608087019190915260a086015250505060c082015290565b6000611e24612adb565b60208083015160009081526007825260409081902081516101008101835281546bffffffffffffffffffffffff808216835273ffffffffffffffffffffffffffffffffffffffff6c0100000000000000000000000092839004811696840196909652600184015490811694830194909452909204831660608301526002015463ffffffff808216608084015264010000000082041660a0830181905268010000000000000000820490931660c083015260ff7c010000000000000000000000000000000000000000000000000000000090910416151560e0820152904310611f38576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611f4b8184600001518560600151612b4d565b60005a90506000634585e33b60e01b8560400151604051602401611f6f91906136e6565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050509050611fe185608001518460c0015183612c9e565b93505a611fee9083613945565b91506000612007838760a001518860c001516001612442565b6020808801516000908152600790915260409020549091506120389082906bffffffffffffffffffffffff1661395c565b6020878101805160009081526007909252604080832080547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff9586161790559051825290206001015461209b9183911661389d565b60208781018051600090815260078352604080822060010180547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff9687161790558a5192518252808220805486166c0100000000000000000000000073ffffffffffffffffffffffffffffffffffffffff958616021790558a51909216815260089092529020546121549183917401000000000000000000000000000000000000000090041661389d565b60086000886000015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160146101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550856000015173ffffffffffffffffffffffffffffffffffffffff1685151587602001517fcaacad83e47cc45c280d487ec84184eee2fa3b54ebaa393bda7549f13da228f6848a6040015160405161222092919061380d565b60405180910390a4505050506122366001600255565b919050565b6000611c038383612cea565b6000806000600d600001600f9054906101000a900462ffffff1662ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b1580156122de57600080fd5b505afa1580156122f2573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061231691906133eb565b50945090925084915050801561233a57506123318242613945565b8463ffffffff16105b80612346575060008113155b1561235557600f549550612359565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b1580156123bf57600080fd5b505afa1580156123d3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906123f791906133eb565b50945090925084915050801561241b57506124128242613945565b8463ffffffff16105b80612427575060008113155b1561243657601054945061243a565b8094505b505050509091565b6040805161012081018252600d5463ffffffff80821683526401000000008083048216602085015268010000000000000000830462ffffff908116958501959095526b0100000000000000000000008304821660608501526f01000000000000000000000000000000830490941660808401527201000000000000000000000000000000000000820461ffff1660a08401819052740100000000000000000000000000000000000000009092046bffffffffffffffffffffffff1660c0840152600e5480821660e08501529390930490921661010082015260009182906125299087613908565b90508380156125375750803a105b1561253f57503a5b600061256b7f000000000000000000000000000000000000000000000000000000000000000089613885565b6125759083613908565b83519091506000906125919063ffffffff16633b9aca00613885565b9050600060027f000000000000000000000000000000000000000000000000000000000000000060028111156125c9576125c9613a71565b141561271157604080516000815260208101909152871561262857600036604051806080016040528060488152602001613b3e604891396040516020016126129392919061359c565b6040516020818303038152906040529050612647565b6040518061014001604052806101108152602001613b86610110913990505b6040517f49948e0e00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016906349948e0e906126b99084906004016136e6565b60206040518083038186803b1580156126d157600080fd5b505afa1580156126e5573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061270991906132fb565b9150506127ec565b60017f0000000000000000000000000000000000000000000000000000000000000000600281111561274557612745613a71565b14156127ec577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b815260040160206040518083038186803b1580156127b157600080fd5b505afa1580156127c5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906127e991906132fb565b90505b8661280857808560a0015161ffff166128059190613908565b90505b6000856020015163ffffffff1664e8d4a510006128259190613908565b89846128318588613885565b61283f90633b9aca00613908565b6128499190613908565b61285391906138cd565b61285d9190613885565b90506b033b2e3c9fd0803ce80000008111156128a5576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b9a9950505050505050505050565b3215610a70576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005473ffffffffffffffffffffffffffffffffffffffff163314610a70576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401611043565b73ffffffffffffffffffffffffffffffffffffffff81163314156129ec576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401611043565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000826000018281548110612a7957612a79613acf565b9060005260206000200154905092915050565b6000818152600183016020526040812054612ad357508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610b68565b506000610b68565b600280541415612b47576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c006044820152606401611043565b60028055565b8260e0015115612b89576040517f514b6c2400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff821660009081526008602052604090206001015460ff16612beb576040517fcfbacfd800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b82516bffffffffffffffffffffffff16811115612c34576040517f356680b700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8173ffffffffffffffffffffffffffffffffffffffff16836020015173ffffffffffffffffffffffffffffffffffffffff16141561112e576040517f06bc104000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005a611388811015612cb057600080fd5b611388810390508460408204820311612cc857600080fd5b50823b612cd457600080fd5b60008083516020850160008789f1949350505050565b60008181526001830160205260408120548015612dd3576000612d0e600183613945565b8554909150600090612d2290600190613945565b9050818114612d87576000866000018281548110612d4257612d42613acf565b9060005260206000200154905080876000018481548110612d6557612d65613acf565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080612d9857612d98613aa0565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610b68565b6000915050610b68565b828054612de9906139b5565b90600052602060002090601f016020900481019282612e0b5760008555612e6f565b82601f10612e42578280017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00823516178555612e6f565b82800160010185558215612e6f579182015b82811115612e6f578235825591602001919060010190612e54565b50612e7b929150612e7f565b5090565b5b80821115612e7b5760008155600101612e80565b803573ffffffffffffffffffffffffffffffffffffffff8116811461223657600080fd5b60008083601f840112612eca57600080fd5b50813567ffffffffffffffff811115612ee257600080fd5b6020830191508360208260051b8501011115612efd57600080fd5b9250929050565b60008083601f840112612f1657600080fd5b50813567ffffffffffffffff811115612f2e57600080fd5b602083019150836020828501011115612efd57600080fd5b803561ffff8116811461223657600080fd5b803562ffffff8116811461223657600080fd5b803563ffffffff8116811461223657600080fd5b805169ffffffffffffffffffff8116811461223657600080fd5b80356bffffffffffffffffffffffff8116811461223657600080fd5b600060208284031215612fc757600080fd5b611c0382612e94565b60008060408385031215612fe357600080fd5b612fec83612e94565b9150612ffa60208401612e94565b90509250929050565b6000806040838503121561301657600080fd5b61301f83612e94565b915060208301356004811061303357600080fd5b809150509250929050565b6000806000806060858703121561305457600080fd5b61305d85612e94565b935060208501359250604085013567ffffffffffffffff81111561308057600080fd5b61308c87828801612f04565b95989497509550505050565b6000806000806000608086880312156130b057600080fd5b6130b986612e94565b94506130c760208701612f6b565b93506130d560408701612e94565b9250606086013567ffffffffffffffff8111156130f157600080fd5b6130fd88828901612f04565b969995985093965092949392505050565b6000806000806040858703121561312457600080fd5b843567ffffffffffffffff8082111561313c57600080fd5b61314888838901612eb8565b9096509450602087013591508082111561316157600080fd5b5061308c87828801612eb8565b60008060006040848603121561318357600080fd5b833567ffffffffffffffff81111561319a57600080fd5b6131a686828701612eb8565b90945092506131b9905060208501612e94565b90509250925092565b600080602083850312156131d557600080fd5b823567ffffffffffffffff8111156131ec57600080fd5b6131f885828601612f04565b90969095509350505050565b6000610180828403121561321757600080fd5b61321f613834565b61322883612f6b565b815261323660208401612f6b565b602082015261324760408401612f58565b604082015261325860608401612f6b565b606082015261326960808401612f58565b608082015261327a60a08401612f46565b60a082015261328b60c08401612f99565b60c082015261329c60e08401612f6b565b60e0820152610100838101359082015261012080840135908201526101406132c5818501612e94565b908201526101606132d7848201612e94565b908201529392505050565b6000602082840312156132f457600080fd5b5035919050565b60006020828403121561330d57600080fd5b5051919050565b6000806040838503121561332757600080fd5b82359150612ffa60208401612e94565b60008060006040848603121561334c57600080fd5b83359250602084013567ffffffffffffffff81111561336a57600080fd5b61337686828701612f04565b9497909650939450505050565b6000806040838503121561339657600080fd5b50508035926020909101359150565b600080604083850312156133b857600080fd5b82359150612ffa60208401612f6b565b600080604083850312156133db57600080fd5b82359150612ffa60208401612f99565b600080600080600060a0868803121561340357600080fd5b61340c86612f7f565b945060208601519350604086015192506060860151915061342f60808701612f7f565b90509295509295909350565b60008151808452613453816020860160208601613989565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b805163ffffffff16825260208101516134a6602084018263ffffffff169052565b5060408101516134bd604084018262ffffff169052565b5060608101516134d5606084018263ffffffff169052565b5060808101516134ec608084018262ffffff169052565b5060a081015161350260a084018261ffff169052565b5060c081015161352260c08401826bffffffffffffffffffffffff169052565b5060e081015161353a60e084018263ffffffff169052565b50610100818101519083015261012080820151908301526101408082015173ffffffffffffffffffffffffffffffffffffffff81168285015250506101608181015173ffffffffffffffffffffffffffffffffffffffff811684830152611661565b8284823760008382016000815283516135b9818360208801613989565b0195945050505050565b600061012073ffffffffffffffffffffffffffffffffffffffff808d16845263ffffffff8c1660208501528160408501526136008285018c61343b565b6bffffffffffffffffffffffff9a8b16606086015298811660808501529690961660a08301525067ffffffffffffffff9390931660c0840152941660e082015292151561010090930192909252949350505050565b6020808252825182820181905260009190848201906040850190845b8181101561368d57835183529284019291840191600101613671565b50909695505050505050565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b602081526000611c03602083018461343b565b60a08152600061370c60a083018861343b565b90508560208301528460408301528360608301528260808301529695505050505050565b602081016004831061374457613744613a71565b91905290565b6020810161374483613b2d565b6101808101610b688284613485565b600061022080830163ffffffff875116845260206bffffffffffffffffffffffff81890151168186015260408801516040860152606088015160608601526137b16080860188613485565b6102008501929092528451908190526102408401918086019160005b818110156137ff57835173ffffffffffffffffffffffffffffffffffffffff16855293820193928201926001016137cd565b509298975050505050505050565b6bffffffffffffffffffffffff8316815260406020820152600061111e604083018461343b565b604051610180810167ffffffffffffffff8111828210171561387f577f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405290565b6000821982111561389857613898613a42565b500190565b60006bffffffffffffffffffffffff8083168185168083038211156138c4576138c4613a42565b01949350505050565b600082613903577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561394057613940613a42565b500290565b60008282101561395757613957613a42565b500390565b60006bffffffffffffffffffffffff8381169083168181101561398157613981613a42565b039392505050565b60005b838110156139a457818101518382015260200161398c565b838111156116615750506000910152565b600181811c908216806139c957607f821691505b60208210811415613a03577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415613a3b57613a3b613a42565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6003811061163557611635613a7156fe3078666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666663078666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000806000a",
}

var KeeperRegistry13ABI = KeeperRegistry13MetaData.ABI

var KeeperRegistry13Bin = KeeperRegistry13MetaData.Bin

func DeployKeeperRegistry13(auth *bind.TransactOpts, backend bind.ContractBackend, keeperRegistryLogic common.Address, config Config) (common.Address, *types.Transaction, *KeeperRegistry13, error) {
	parsed, err := KeeperRegistry13MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistry13Bin), backend, keeperRegistryLogic, config)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperRegistry13{address: address, abi: *parsed, KeeperRegistry13Caller: KeeperRegistry13Caller{contract: contract}, KeeperRegistry13Transactor: KeeperRegistry13Transactor{contract: contract}, KeeperRegistry13Filterer: KeeperRegistry13Filterer{contract: contract}}, nil
}

type KeeperRegistry13 struct {
	address common.Address
	abi     abi.ABI
	KeeperRegistry13Caller
	KeeperRegistry13Transactor
	KeeperRegistry13Filterer
}

type KeeperRegistry13Caller struct {
	contract *bind.BoundContract
}

type KeeperRegistry13Transactor struct {
	contract *bind.BoundContract
}

type KeeperRegistry13Filterer struct {
	contract *bind.BoundContract
}

type KeeperRegistry13Session struct {
	Contract     *KeeperRegistry13
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeeperRegistry13CallerSession struct {
	Contract *KeeperRegistry13Caller
	CallOpts bind.CallOpts
}

type KeeperRegistry13TransactorSession struct {
	Contract     *KeeperRegistry13Transactor
	TransactOpts bind.TransactOpts
}

type KeeperRegistry13Raw struct {
	Contract *KeeperRegistry13
}

type KeeperRegistry13CallerRaw struct {
	Contract *KeeperRegistry13Caller
}

type KeeperRegistry13TransactorRaw struct {
	Contract *KeeperRegistry13Transactor
}

func NewKeeperRegistry13(address common.Address, backend bind.ContractBackend) (*KeeperRegistry13, error) {
	abi, err := abi.JSON(strings.NewReader(KeeperRegistry13ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeeperRegistry13(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13{address: address, abi: abi, KeeperRegistry13Caller: KeeperRegistry13Caller{contract: contract}, KeeperRegistry13Transactor: KeeperRegistry13Transactor{contract: contract}, KeeperRegistry13Filterer: KeeperRegistry13Filterer{contract: contract}}, nil
}

func NewKeeperRegistry13Caller(address common.Address, caller bind.ContractCaller) (*KeeperRegistry13Caller, error) {
	contract, err := bindKeeperRegistry13(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13Caller{contract: contract}, nil
}

func NewKeeperRegistry13Transactor(address common.Address, transactor bind.ContractTransactor) (*KeeperRegistry13Transactor, error) {
	contract, err := bindKeeperRegistry13(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13Transactor{contract: contract}, nil
}

func NewKeeperRegistry13Filterer(address common.Address, filterer bind.ContractFilterer) (*KeeperRegistry13Filterer, error) {
	contract, err := bindKeeperRegistry13(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13Filterer{contract: contract}, nil
}

func bindKeeperRegistry13(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeeperRegistry13MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_KeeperRegistry13 *KeeperRegistry13Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistry13.Contract.KeeperRegistry13Caller.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistry13 *KeeperRegistry13Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.KeeperRegistry13Transactor.contract.Transfer(opts)
}

func (_KeeperRegistry13 *KeeperRegistry13Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.KeeperRegistry13Transactor.contract.Transact(opts, method, params...)
}

func (_KeeperRegistry13 *KeeperRegistry13CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistry13.Contract.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.contract.Transfer(opts)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.contract.Transact(opts, method, params...)
}

func (_KeeperRegistry13 *KeeperRegistry13Caller) ARBNITROORACLE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "ARB_NITRO_ORACLE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistry13 *KeeperRegistry13Session) ARBNITROORACLE() (common.Address, error) {
	return _KeeperRegistry13.Contract.ARBNITROORACLE(&_KeeperRegistry13.CallOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13CallerSession) ARBNITROORACLE() (common.Address, error) {
	return _KeeperRegistry13.Contract.ARBNITROORACLE(&_KeeperRegistry13.CallOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13Caller) FASTGASFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "FAST_GAS_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistry13 *KeeperRegistry13Session) FASTGASFEED() (common.Address, error) {
	return _KeeperRegistry13.Contract.FASTGASFEED(&_KeeperRegistry13.CallOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13CallerSession) FASTGASFEED() (common.Address, error) {
	return _KeeperRegistry13.Contract.FASTGASFEED(&_KeeperRegistry13.CallOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13Caller) KEEPERREGISTRYLOGIC(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "KEEPER_REGISTRY_LOGIC")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistry13 *KeeperRegistry13Session) KEEPERREGISTRYLOGIC() (common.Address, error) {
	return _KeeperRegistry13.Contract.KEEPERREGISTRYLOGIC(&_KeeperRegistry13.CallOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13CallerSession) KEEPERREGISTRYLOGIC() (common.Address, error) {
	return _KeeperRegistry13.Contract.KEEPERREGISTRYLOGIC(&_KeeperRegistry13.CallOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13Caller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistry13 *KeeperRegistry13Session) LINK() (common.Address, error) {
	return _KeeperRegistry13.Contract.LINK(&_KeeperRegistry13.CallOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13CallerSession) LINK() (common.Address, error) {
	return _KeeperRegistry13.Contract.LINK(&_KeeperRegistry13.CallOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13Caller) LINKETHFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "LINK_ETH_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistry13 *KeeperRegistry13Session) LINKETHFEED() (common.Address, error) {
	return _KeeperRegistry13.Contract.LINKETHFEED(&_KeeperRegistry13.CallOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13CallerSession) LINKETHFEED() (common.Address, error) {
	return _KeeperRegistry13.Contract.LINKETHFEED(&_KeeperRegistry13.CallOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13Caller) OPTIMISMORACLE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "OPTIMISM_ORACLE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistry13 *KeeperRegistry13Session) OPTIMISMORACLE() (common.Address, error) {
	return _KeeperRegistry13.Contract.OPTIMISMORACLE(&_KeeperRegistry13.CallOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13CallerSession) OPTIMISMORACLE() (common.Address, error) {
	return _KeeperRegistry13.Contract.OPTIMISMORACLE(&_KeeperRegistry13.CallOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13Caller) PAYMENTMODEL(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "PAYMENT_MODEL")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistry13 *KeeperRegistry13Session) PAYMENTMODEL() (uint8, error) {
	return _KeeperRegistry13.Contract.PAYMENTMODEL(&_KeeperRegistry13.CallOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13CallerSession) PAYMENTMODEL() (uint8, error) {
	return _KeeperRegistry13.Contract.PAYMENTMODEL(&_KeeperRegistry13.CallOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13Caller) REGISTRYGASOVERHEAD(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "REGISTRY_GAS_OVERHEAD")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistry13 *KeeperRegistry13Session) REGISTRYGASOVERHEAD() (*big.Int, error) {
	return _KeeperRegistry13.Contract.REGISTRYGASOVERHEAD(&_KeeperRegistry13.CallOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13CallerSession) REGISTRYGASOVERHEAD() (*big.Int, error) {
	return _KeeperRegistry13.Contract.REGISTRYGASOVERHEAD(&_KeeperRegistry13.CallOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13Caller) GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "getActiveUpkeepIDs", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_KeeperRegistry13 *KeeperRegistry13Session) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _KeeperRegistry13.Contract.GetActiveUpkeepIDs(&_KeeperRegistry13.CallOpts, startIndex, maxCount)
}

func (_KeeperRegistry13 *KeeperRegistry13CallerSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _KeeperRegistry13.Contract.GetActiveUpkeepIDs(&_KeeperRegistry13.CallOpts, startIndex, maxCount)
}

func (_KeeperRegistry13 *KeeperRegistry13Caller) GetKeeperInfo(opts *bind.CallOpts, query common.Address) (GetKeeperInfo,

	error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "getKeeperInfo", query)

	outstruct := new(GetKeeperInfo)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Payee = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Active = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.Balance = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_KeeperRegistry13 *KeeperRegistry13Session) GetKeeperInfo(query common.Address) (GetKeeperInfo,

	error) {
	return _KeeperRegistry13.Contract.GetKeeperInfo(&_KeeperRegistry13.CallOpts, query)
}

func (_KeeperRegistry13 *KeeperRegistry13CallerSession) GetKeeperInfo(query common.Address) (GetKeeperInfo,

	error) {
	return _KeeperRegistry13.Contract.GetKeeperInfo(&_KeeperRegistry13.CallOpts, query)
}

func (_KeeperRegistry13 *KeeperRegistry13Caller) GetMaxPaymentForGas(opts *bind.CallOpts, gasLimit *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "getMaxPaymentForGas", gasLimit)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistry13 *KeeperRegistry13Session) GetMaxPaymentForGas(gasLimit *big.Int) (*big.Int, error) {
	return _KeeperRegistry13.Contract.GetMaxPaymentForGas(&_KeeperRegistry13.CallOpts, gasLimit)
}

func (_KeeperRegistry13 *KeeperRegistry13CallerSession) GetMaxPaymentForGas(gasLimit *big.Int) (*big.Int, error) {
	return _KeeperRegistry13.Contract.GetMaxPaymentForGas(&_KeeperRegistry13.CallOpts, gasLimit)
}

func (_KeeperRegistry13 *KeeperRegistry13Caller) GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "getMinBalanceForUpkeep", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistry13 *KeeperRegistry13Session) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _KeeperRegistry13.Contract.GetMinBalanceForUpkeep(&_KeeperRegistry13.CallOpts, id)
}

func (_KeeperRegistry13 *KeeperRegistry13CallerSession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _KeeperRegistry13.Contract.GetMinBalanceForUpkeep(&_KeeperRegistry13.CallOpts, id)
}

func (_KeeperRegistry13 *KeeperRegistry13Caller) GetPeerRegistryMigrationPermission(opts *bind.CallOpts, peer common.Address) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "getPeerRegistryMigrationPermission", peer)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistry13 *KeeperRegistry13Session) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _KeeperRegistry13.Contract.GetPeerRegistryMigrationPermission(&_KeeperRegistry13.CallOpts, peer)
}

func (_KeeperRegistry13 *KeeperRegistry13CallerSession) GetPeerRegistryMigrationPermission(peer common.Address) (uint8, error) {
	return _KeeperRegistry13.Contract.GetPeerRegistryMigrationPermission(&_KeeperRegistry13.CallOpts, peer)
}

func (_KeeperRegistry13 *KeeperRegistry13Caller) GetState(opts *bind.CallOpts) (GetState,

	error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "getState")

	outstruct := new(GetState)
	if err != nil {
		return *outstruct, err
	}

	outstruct.State = *abi.ConvertType(out[0], new(State)).(*State)
	outstruct.Config = *abi.ConvertType(out[1], new(Config)).(*Config)
	outstruct.Keepers = *abi.ConvertType(out[2], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

func (_KeeperRegistry13 *KeeperRegistry13Session) GetState() (GetState,

	error) {
	return _KeeperRegistry13.Contract.GetState(&_KeeperRegistry13.CallOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13CallerSession) GetState() (GetState,

	error) {
	return _KeeperRegistry13.Contract.GetState(&_KeeperRegistry13.CallOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13Caller) GetUpkeep(opts *bind.CallOpts, id *big.Int) (GetUpkeep,

	error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "getUpkeep", id)

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
	outstruct.Paused = *abi.ConvertType(out[8], new(bool)).(*bool)

	return *outstruct, err

}

func (_KeeperRegistry13 *KeeperRegistry13Session) GetUpkeep(id *big.Int) (GetUpkeep,

	error) {
	return _KeeperRegistry13.Contract.GetUpkeep(&_KeeperRegistry13.CallOpts, id)
}

func (_KeeperRegistry13 *KeeperRegistry13CallerSession) GetUpkeep(id *big.Int) (GetUpkeep,

	error) {
	return _KeeperRegistry13.Contract.GetUpkeep(&_KeeperRegistry13.CallOpts, id)
}

func (_KeeperRegistry13 *KeeperRegistry13Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistry13 *KeeperRegistry13Session) Owner() (common.Address, error) {
	return _KeeperRegistry13.Contract.Owner(&_KeeperRegistry13.CallOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13CallerSession) Owner() (common.Address, error) {
	return _KeeperRegistry13.Contract.Owner(&_KeeperRegistry13.CallOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13Caller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_KeeperRegistry13 *KeeperRegistry13Session) Paused() (bool, error) {
	return _KeeperRegistry13.Contract.Paused(&_KeeperRegistry13.CallOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13CallerSession) Paused() (bool, error) {
	return _KeeperRegistry13.Contract.Paused(&_KeeperRegistry13.CallOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13Caller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_KeeperRegistry13 *KeeperRegistry13Session) TypeAndVersion() (string, error) {
	return _KeeperRegistry13.Contract.TypeAndVersion(&_KeeperRegistry13.CallOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13CallerSession) TypeAndVersion() (string, error) {
	return _KeeperRegistry13.Contract.TypeAndVersion(&_KeeperRegistry13.CallOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13Caller) UpkeepTranscoderVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistry13.contract.Call(opts, &out, "upkeepTranscoderVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistry13 *KeeperRegistry13Session) UpkeepTranscoderVersion() (uint8, error) {
	return _KeeperRegistry13.Contract.UpkeepTranscoderVersion(&_KeeperRegistry13.CallOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13CallerSession) UpkeepTranscoderVersion() (uint8, error) {
	return _KeeperRegistry13.Contract.UpkeepTranscoderVersion(&_KeeperRegistry13.CallOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "acceptOwnership")
}

func (_KeeperRegistry13 *KeeperRegistry13Session) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.AcceptOwnership(&_KeeperRegistry13.TransactOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.AcceptOwnership(&_KeeperRegistry13.TransactOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) AcceptPayeeship(opts *bind.TransactOpts, keeper common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "acceptPayeeship", keeper)
}

func (_KeeperRegistry13 *KeeperRegistry13Session) AcceptPayeeship(keeper common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.AcceptPayeeship(&_KeeperRegistry13.TransactOpts, keeper)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) AcceptPayeeship(keeper common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.AcceptPayeeship(&_KeeperRegistry13.TransactOpts, keeper)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "acceptUpkeepAdmin", id)
}

func (_KeeperRegistry13 *KeeperRegistry13Session) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.AcceptUpkeepAdmin(&_KeeperRegistry13.TransactOpts, id)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.AcceptUpkeepAdmin(&_KeeperRegistry13.TransactOpts, id)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "addFunds", id, amount)
}

func (_KeeperRegistry13 *KeeperRegistry13Session) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.AddFunds(&_KeeperRegistry13.TransactOpts, id, amount)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.AddFunds(&_KeeperRegistry13.TransactOpts, id, amount)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "cancelUpkeep", id)
}

func (_KeeperRegistry13 *KeeperRegistry13Session) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.CancelUpkeep(&_KeeperRegistry13.TransactOpts, id)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.CancelUpkeep(&_KeeperRegistry13.TransactOpts, id)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) CheckUpkeep(opts *bind.TransactOpts, id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "checkUpkeep", id, from)
}

func (_KeeperRegistry13 *KeeperRegistry13Session) CheckUpkeep(id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.CheckUpkeep(&_KeeperRegistry13.TransactOpts, id, from)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) CheckUpkeep(id *big.Int, from common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.CheckUpkeep(&_KeeperRegistry13.TransactOpts, id, from)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "migrateUpkeeps", ids, destination)
}

func (_KeeperRegistry13 *KeeperRegistry13Session) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.MigrateUpkeeps(&_KeeperRegistry13.TransactOpts, ids, destination)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.MigrateUpkeeps(&_KeeperRegistry13.TransactOpts, ids, destination)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "onTokenTransfer", sender, amount, data)
}

func (_KeeperRegistry13 *KeeperRegistry13Session) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.OnTokenTransfer(&_KeeperRegistry13.TransactOpts, sender, amount, data)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.OnTokenTransfer(&_KeeperRegistry13.TransactOpts, sender, amount, data)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "pause")
}

func (_KeeperRegistry13 *KeeperRegistry13Session) Pause() (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.Pause(&_KeeperRegistry13.TransactOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) Pause() (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.Pause(&_KeeperRegistry13.TransactOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "pauseUpkeep", id)
}

func (_KeeperRegistry13 *KeeperRegistry13Session) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.PauseUpkeep(&_KeeperRegistry13.TransactOpts, id)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.PauseUpkeep(&_KeeperRegistry13.TransactOpts, id)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) PerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "performUpkeep", id, performData)
}

func (_KeeperRegistry13 *KeeperRegistry13Session) PerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.PerformUpkeep(&_KeeperRegistry13.TransactOpts, id, performData)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) PerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.PerformUpkeep(&_KeeperRegistry13.TransactOpts, id, performData)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "receiveUpkeeps", encodedUpkeeps)
}

func (_KeeperRegistry13 *KeeperRegistry13Session) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.ReceiveUpkeeps(&_KeeperRegistry13.TransactOpts, encodedUpkeeps)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.ReceiveUpkeeps(&_KeeperRegistry13.TransactOpts, encodedUpkeeps)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "recoverFunds")
}

func (_KeeperRegistry13 *KeeperRegistry13Session) RecoverFunds() (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.RecoverFunds(&_KeeperRegistry13.TransactOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) RecoverFunds() (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.RecoverFunds(&_KeeperRegistry13.TransactOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "registerUpkeep", target, gasLimit, admin, checkData)
}

func (_KeeperRegistry13 *KeeperRegistry13Session) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.RegisterUpkeep(&_KeeperRegistry13.TransactOpts, target, gasLimit, admin, checkData)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.RegisterUpkeep(&_KeeperRegistry13.TransactOpts, target, gasLimit, admin, checkData)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) SetConfig(opts *bind.TransactOpts, config Config) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "setConfig", config)
}

func (_KeeperRegistry13 *KeeperRegistry13Session) SetConfig(config Config) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.SetConfig(&_KeeperRegistry13.TransactOpts, config)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) SetConfig(config Config) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.SetConfig(&_KeeperRegistry13.TransactOpts, config)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) SetKeepers(opts *bind.TransactOpts, keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "setKeepers", keepers, payees)
}

func (_KeeperRegistry13 *KeeperRegistry13Session) SetKeepers(keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.SetKeepers(&_KeeperRegistry13.TransactOpts, keepers, payees)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) SetKeepers(keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.SetKeepers(&_KeeperRegistry13.TransactOpts, keepers, payees)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "setPeerRegistryMigrationPermission", peer, permission)
}

func (_KeeperRegistry13 *KeeperRegistry13Session) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.SetPeerRegistryMigrationPermission(&_KeeperRegistry13.TransactOpts, peer, permission)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.SetPeerRegistryMigrationPermission(&_KeeperRegistry13.TransactOpts, peer, permission)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "setUpkeepGasLimit", id, gasLimit)
}

func (_KeeperRegistry13 *KeeperRegistry13Session) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.SetUpkeepGasLimit(&_KeeperRegistry13.TransactOpts, id, gasLimit)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.SetUpkeepGasLimit(&_KeeperRegistry13.TransactOpts, id, gasLimit)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "transferOwnership", to)
}

func (_KeeperRegistry13 *KeeperRegistry13Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.TransferOwnership(&_KeeperRegistry13.TransactOpts, to)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.TransferOwnership(&_KeeperRegistry13.TransactOpts, to)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) TransferPayeeship(opts *bind.TransactOpts, keeper common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "transferPayeeship", keeper, proposed)
}

func (_KeeperRegistry13 *KeeperRegistry13Session) TransferPayeeship(keeper common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.TransferPayeeship(&_KeeperRegistry13.TransactOpts, keeper, proposed)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) TransferPayeeship(keeper common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.TransferPayeeship(&_KeeperRegistry13.TransactOpts, keeper, proposed)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "transferUpkeepAdmin", id, proposed)
}

func (_KeeperRegistry13 *KeeperRegistry13Session) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.TransferUpkeepAdmin(&_KeeperRegistry13.TransactOpts, id, proposed)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.TransferUpkeepAdmin(&_KeeperRegistry13.TransactOpts, id, proposed)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "unpause")
}

func (_KeeperRegistry13 *KeeperRegistry13Session) Unpause() (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.Unpause(&_KeeperRegistry13.TransactOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) Unpause() (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.Unpause(&_KeeperRegistry13.TransactOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "unpauseUpkeep", id)
}

func (_KeeperRegistry13 *KeeperRegistry13Session) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.UnpauseUpkeep(&_KeeperRegistry13.TransactOpts, id)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.UnpauseUpkeep(&_KeeperRegistry13.TransactOpts, id)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) UpdateCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "updateCheckData", id, newCheckData)
}

func (_KeeperRegistry13 *KeeperRegistry13Session) UpdateCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.UpdateCheckData(&_KeeperRegistry13.TransactOpts, id, newCheckData)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) UpdateCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.UpdateCheckData(&_KeeperRegistry13.TransactOpts, id, newCheckData)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "withdrawFunds", id, to)
}

func (_KeeperRegistry13 *KeeperRegistry13Session) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.WithdrawFunds(&_KeeperRegistry13.TransactOpts, id, to)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.WithdrawFunds(&_KeeperRegistry13.TransactOpts, id, to)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) WithdrawOwnerFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "withdrawOwnerFunds")
}

func (_KeeperRegistry13 *KeeperRegistry13Session) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.WithdrawOwnerFunds(&_KeeperRegistry13.TransactOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.WithdrawOwnerFunds(&_KeeperRegistry13.TransactOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.Transact(opts, "withdrawPayment", from, to)
}

func (_KeeperRegistry13 *KeeperRegistry13Session) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.WithdrawPayment(&_KeeperRegistry13.TransactOpts, from, to)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.WithdrawPayment(&_KeeperRegistry13.TransactOpts, from, to)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.RawTransact(opts, calldata)
}

func (_KeeperRegistry13 *KeeperRegistry13Session) Fallback(calldata []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.Fallback(&_KeeperRegistry13.TransactOpts, calldata)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.Fallback(&_KeeperRegistry13.TransactOpts, calldata)
}

func (_KeeperRegistry13 *KeeperRegistry13Transactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistry13.contract.RawTransact(opts, nil)
}

func (_KeeperRegistry13 *KeeperRegistry13Session) Receive() (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.Receive(&_KeeperRegistry13.TransactOpts)
}

func (_KeeperRegistry13 *KeeperRegistry13TransactorSession) Receive() (*types.Transaction, error) {
	return _KeeperRegistry13.Contract.Receive(&_KeeperRegistry13.TransactOpts)
}

type KeeperRegistry13ConfigSetIterator struct {
	Event *KeeperRegistry13ConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistry13ConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13ConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry13ConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistry13ConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistry13ConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistry13ConfigSet struct {
	Config Config
	Raw    types.Log
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterConfigSet(opts *bind.FilterOpts) (*KeeperRegistry13ConfigSetIterator, error) {

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13ConfigSetIterator{contract: _KeeperRegistry13.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13ConfigSet) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistry13ConfigSet)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "ConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseConfigSet(log types.Log) (*KeeperRegistry13ConfigSet, error) {
	event := new(KeeperRegistry13ConfigSet)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistry13FundsAddedIterator struct {
	Event *KeeperRegistry13FundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistry13FundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13FundsAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry13FundsAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistry13FundsAddedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistry13FundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistry13FundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistry13FundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13FundsAddedIterator{contract: _KeeperRegistry13.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13FundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistry13FundsAdded)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "FundsAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseFundsAdded(log types.Log) (*KeeperRegistry13FundsAdded, error) {
	event := new(KeeperRegistry13FundsAdded)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistry13FundsWithdrawnIterator struct {
	Event *KeeperRegistry13FundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistry13FundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13FundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry13FundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistry13FundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistry13FundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistry13FundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry13FundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13FundsWithdrawnIterator{contract: _KeeperRegistry13.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13FundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistry13FundsWithdrawn)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseFundsWithdrawn(log types.Log) (*KeeperRegistry13FundsWithdrawn, error) {
	event := new(KeeperRegistry13FundsWithdrawn)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistry13KeepersUpdatedIterator struct {
	Event *KeeperRegistry13KeepersUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistry13KeepersUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13KeepersUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry13KeepersUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistry13KeepersUpdatedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistry13KeepersUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistry13KeepersUpdated struct {
	Keepers []common.Address
	Payees  []common.Address
	Raw     types.Log
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterKeepersUpdated(opts *bind.FilterOpts) (*KeeperRegistry13KeepersUpdatedIterator, error) {

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "KeepersUpdated")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13KeepersUpdatedIterator{contract: _KeeperRegistry13.contract, event: "KeepersUpdated", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchKeepersUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13KeepersUpdated) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "KeepersUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistry13KeepersUpdated)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "KeepersUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseKeepersUpdated(log types.Log) (*KeeperRegistry13KeepersUpdated, error) {
	event := new(KeeperRegistry13KeepersUpdated)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "KeepersUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistry13OwnerFundsWithdrawnIterator struct {
	Event *KeeperRegistry13OwnerFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistry13OwnerFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13OwnerFundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry13OwnerFundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistry13OwnerFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistry13OwnerFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistry13OwnerFundsWithdrawn struct {
	Amount *big.Int
	Raw    types.Log
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*KeeperRegistry13OwnerFundsWithdrawnIterator, error) {

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13OwnerFundsWithdrawnIterator{contract: _KeeperRegistry13.contract, event: "OwnerFundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13OwnerFundsWithdrawn) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistry13OwnerFundsWithdrawn)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseOwnerFundsWithdrawn(log types.Log) (*KeeperRegistry13OwnerFundsWithdrawn, error) {
	event := new(KeeperRegistry13OwnerFundsWithdrawn)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistry13OwnershipTransferRequestedIterator struct {
	Event *KeeperRegistry13OwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistry13OwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13OwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry13OwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistry13OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistry13OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistry13OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistry13OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13OwnershipTransferRequestedIterator{contract: _KeeperRegistry13.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistry13OwnershipTransferRequested)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistry13OwnershipTransferRequested, error) {
	event := new(KeeperRegistry13OwnershipTransferRequested)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistry13OwnershipTransferredIterator struct {
	Event *KeeperRegistry13OwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistry13OwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13OwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry13OwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistry13OwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistry13OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistry13OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistry13OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13OwnershipTransferredIterator{contract: _KeeperRegistry13.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistry13OwnershipTransferred)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseOwnershipTransferred(log types.Log) (*KeeperRegistry13OwnershipTransferred, error) {
	event := new(KeeperRegistry13OwnershipTransferred)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistry13PausedIterator struct {
	Event *KeeperRegistry13Paused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistry13PausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13Paused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry13Paused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistry13PausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistry13PausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistry13Paused struct {
	Account common.Address
	Raw     types.Log
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterPaused(opts *bind.FilterOpts) (*KeeperRegistry13PausedIterator, error) {

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13PausedIterator{contract: _KeeperRegistry13.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13Paused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistry13Paused)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParsePaused(log types.Log) (*KeeperRegistry13Paused, error) {
	event := new(KeeperRegistry13Paused)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistry13PayeeshipTransferRequestedIterator struct {
	Event *KeeperRegistry13PayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistry13PayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13PayeeshipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry13PayeeshipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistry13PayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistry13PayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistry13PayeeshipTransferRequested struct {
	Keeper common.Address
	From   common.Address
	To     common.Address
	Raw    types.Log
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, keeper []common.Address, from []common.Address, to []common.Address) (*KeeperRegistry13PayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "PayeeshipTransferRequested", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13PayeeshipTransferRequestedIterator{contract: _KeeperRegistry13.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13PayeeshipTransferRequested, keeper []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "PayeeshipTransferRequested", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistry13PayeeshipTransferRequested)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistry13PayeeshipTransferRequested, error) {
	event := new(KeeperRegistry13PayeeshipTransferRequested)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistry13PayeeshipTransferredIterator struct {
	Event *KeeperRegistry13PayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistry13PayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13PayeeshipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry13PayeeshipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistry13PayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistry13PayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistry13PayeeshipTransferred struct {
	Keeper common.Address
	From   common.Address
	To     common.Address
	Raw    types.Log
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, keeper []common.Address, from []common.Address, to []common.Address) (*KeeperRegistry13PayeeshipTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "PayeeshipTransferred", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13PayeeshipTransferredIterator{contract: _KeeperRegistry13.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13PayeeshipTransferred, keeper []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "PayeeshipTransferred", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistry13PayeeshipTransferred)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParsePayeeshipTransferred(log types.Log) (*KeeperRegistry13PayeeshipTransferred, error) {
	event := new(KeeperRegistry13PayeeshipTransferred)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistry13PaymentWithdrawnIterator struct {
	Event *KeeperRegistry13PaymentWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistry13PaymentWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13PaymentWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry13PaymentWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistry13PaymentWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistry13PaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistry13PaymentWithdrawn struct {
	Keeper common.Address
	Amount *big.Int
	To     common.Address
	Payee  common.Address
	Raw    types.Log
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, keeper []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistry13PaymentWithdrawnIterator, error) {

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

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "PaymentWithdrawn", keeperRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13PaymentWithdrawnIterator{contract: _KeeperRegistry13.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13PaymentWithdrawn, keeper []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "PaymentWithdrawn", keeperRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistry13PaymentWithdrawn)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParsePaymentWithdrawn(log types.Log) (*KeeperRegistry13PaymentWithdrawn, error) {
	event := new(KeeperRegistry13PaymentWithdrawn)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistry13UnpausedIterator struct {
	Event *KeeperRegistry13Unpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistry13UnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13Unpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry13Unpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistry13UnpausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistry13UnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistry13Unpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistry13UnpausedIterator, error) {

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13UnpausedIterator{contract: _KeeperRegistry13.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13Unpaused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistry13Unpaused)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseUnpaused(log types.Log) (*KeeperRegistry13Unpaused, error) {
	event := new(KeeperRegistry13Unpaused)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistry13UpkeepAdminTransferRequestedIterator struct {
	Event *KeeperRegistry13UpkeepAdminTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistry13UpkeepAdminTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13UpkeepAdminTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry13UpkeepAdminTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistry13UpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistry13UpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistry13UpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistry13UpkeepAdminTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13UpkeepAdminTransferRequestedIterator{contract: _KeeperRegistry13.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistry13UpkeepAdminTransferRequested)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseUpkeepAdminTransferRequested(log types.Log) (*KeeperRegistry13UpkeepAdminTransferRequested, error) {
	event := new(KeeperRegistry13UpkeepAdminTransferRequested)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistry13UpkeepAdminTransferredIterator struct {
	Event *KeeperRegistry13UpkeepAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistry13UpkeepAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13UpkeepAdminTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry13UpkeepAdminTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistry13UpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistry13UpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistry13UpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistry13UpkeepAdminTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13UpkeepAdminTransferredIterator{contract: _KeeperRegistry13.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistry13UpkeepAdminTransferred)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseUpkeepAdminTransferred(log types.Log) (*KeeperRegistry13UpkeepAdminTransferred, error) {
	event := new(KeeperRegistry13UpkeepAdminTransferred)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistry13UpkeepCanceledIterator struct {
	Event *KeeperRegistry13UpkeepCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistry13UpkeepCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13UpkeepCanceled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry13UpkeepCanceled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistry13UpkeepCanceledIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistry13UpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistry13UpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistry13UpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13UpkeepCanceledIterator{contract: _KeeperRegistry13.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistry13UpkeepCanceled)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseUpkeepCanceled(log types.Log) (*KeeperRegistry13UpkeepCanceled, error) {
	event := new(KeeperRegistry13UpkeepCanceled)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistry13UpkeepCheckDataUpdatedIterator struct {
	Event *KeeperRegistry13UpkeepCheckDataUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistry13UpkeepCheckDataUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13UpkeepCheckDataUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry13UpkeepCheckDataUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistry13UpkeepCheckDataUpdatedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistry13UpkeepCheckDataUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistry13UpkeepCheckDataUpdated struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterUpkeepCheckDataUpdated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry13UpkeepCheckDataUpdatedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "UpkeepCheckDataUpdated", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13UpkeepCheckDataUpdatedIterator{contract: _KeeperRegistry13.contract, event: "UpkeepCheckDataUpdated", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchUpkeepCheckDataUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepCheckDataUpdated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "UpkeepCheckDataUpdated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistry13UpkeepCheckDataUpdated)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepCheckDataUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseUpkeepCheckDataUpdated(log types.Log) (*KeeperRegistry13UpkeepCheckDataUpdated, error) {
	event := new(KeeperRegistry13UpkeepCheckDataUpdated)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepCheckDataUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistry13UpkeepGasLimitSetIterator struct {
	Event *KeeperRegistry13UpkeepGasLimitSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistry13UpkeepGasLimitSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13UpkeepGasLimitSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry13UpkeepGasLimitSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistry13UpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistry13UpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistry13UpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry13UpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13UpkeepGasLimitSetIterator{contract: _KeeperRegistry13.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistry13UpkeepGasLimitSet)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseUpkeepGasLimitSet(log types.Log) (*KeeperRegistry13UpkeepGasLimitSet, error) {
	event := new(KeeperRegistry13UpkeepGasLimitSet)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistry13UpkeepMigratedIterator struct {
	Event *KeeperRegistry13UpkeepMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistry13UpkeepMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13UpkeepMigrated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry13UpkeepMigrated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistry13UpkeepMigratedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistry13UpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistry13UpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry13UpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13UpkeepMigratedIterator{contract: _KeeperRegistry13.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistry13UpkeepMigrated)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseUpkeepMigrated(log types.Log) (*KeeperRegistry13UpkeepMigrated, error) {
	event := new(KeeperRegistry13UpkeepMigrated)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistry13UpkeepPausedIterator struct {
	Event *KeeperRegistry13UpkeepPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistry13UpkeepPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13UpkeepPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry13UpkeepPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistry13UpkeepPausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistry13UpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistry13UpkeepPaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry13UpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13UpkeepPausedIterator{contract: _KeeperRegistry13.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistry13UpkeepPaused)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseUpkeepPaused(log types.Log) (*KeeperRegistry13UpkeepPaused, error) {
	event := new(KeeperRegistry13UpkeepPaused)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistry13UpkeepPerformedIterator struct {
	Event *KeeperRegistry13UpkeepPerformed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistry13UpkeepPerformedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13UpkeepPerformed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry13UpkeepPerformed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistry13UpkeepPerformedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistry13UpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistry13UpkeepPerformed struct {
	Id          *big.Int
	Success     bool
	From        common.Address
	Payment     *big.Int
	PerformData []byte
	Raw         types.Log
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool, from []common.Address) (*KeeperRegistry13UpkeepPerformedIterator, error) {

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

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13UpkeepPerformedIterator{contract: _KeeperRegistry13.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepPerformed, id []*big.Int, success []bool, from []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistry13UpkeepPerformed)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseUpkeepPerformed(log types.Log) (*KeeperRegistry13UpkeepPerformed, error) {
	event := new(KeeperRegistry13UpkeepPerformed)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistry13UpkeepReceivedIterator struct {
	Event *KeeperRegistry13UpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistry13UpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13UpkeepReceived)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry13UpkeepReceived)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistry13UpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistry13UpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistry13UpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry13UpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13UpkeepReceivedIterator{contract: _KeeperRegistry13.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistry13UpkeepReceived)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseUpkeepReceived(log types.Log) (*KeeperRegistry13UpkeepReceived, error) {
	event := new(KeeperRegistry13UpkeepReceived)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistry13UpkeepRegisteredIterator struct {
	Event *KeeperRegistry13UpkeepRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistry13UpkeepRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13UpkeepRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry13UpkeepRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistry13UpkeepRegisteredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistry13UpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistry13UpkeepRegistered struct {
	Id         *big.Int
	ExecuteGas uint32
	Admin      common.Address
	Raw        types.Log
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry13UpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13UpkeepRegisteredIterator{contract: _KeeperRegistry13.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistry13UpkeepRegistered)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseUpkeepRegistered(log types.Log) (*KeeperRegistry13UpkeepRegistered, error) {
	event := new(KeeperRegistry13UpkeepRegistered)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistry13UpkeepUnpausedIterator struct {
	Event *KeeperRegistry13UpkeepUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistry13UpkeepUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistry13UpkeepUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistry13UpkeepUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistry13UpkeepUnpausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistry13UpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistry13UpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry13UpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistry13UpkeepUnpausedIterator{contract: _KeeperRegistry13.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistry13.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistry13UpkeepUnpaused)
				if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistry13 *KeeperRegistry13Filterer) ParseUpkeepUnpaused(log types.Log) (*KeeperRegistry13UpkeepUnpaused, error) {
	event := new(KeeperRegistry13UpkeepUnpaused)
	if err := _KeeperRegistry13.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
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
	Paused              bool
}

func (_KeeperRegistry13 *KeeperRegistry13) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _KeeperRegistry13.abi.Events["ConfigSet"].ID:
		return _KeeperRegistry13.ParseConfigSet(log)
	case _KeeperRegistry13.abi.Events["FundsAdded"].ID:
		return _KeeperRegistry13.ParseFundsAdded(log)
	case _KeeperRegistry13.abi.Events["FundsWithdrawn"].ID:
		return _KeeperRegistry13.ParseFundsWithdrawn(log)
	case _KeeperRegistry13.abi.Events["KeepersUpdated"].ID:
		return _KeeperRegistry13.ParseKeepersUpdated(log)
	case _KeeperRegistry13.abi.Events["OwnerFundsWithdrawn"].ID:
		return _KeeperRegistry13.ParseOwnerFundsWithdrawn(log)
	case _KeeperRegistry13.abi.Events["OwnershipTransferRequested"].ID:
		return _KeeperRegistry13.ParseOwnershipTransferRequested(log)
	case _KeeperRegistry13.abi.Events["OwnershipTransferred"].ID:
		return _KeeperRegistry13.ParseOwnershipTransferred(log)
	case _KeeperRegistry13.abi.Events["Paused"].ID:
		return _KeeperRegistry13.ParsePaused(log)
	case _KeeperRegistry13.abi.Events["PayeeshipTransferRequested"].ID:
		return _KeeperRegistry13.ParsePayeeshipTransferRequested(log)
	case _KeeperRegistry13.abi.Events["PayeeshipTransferred"].ID:
		return _KeeperRegistry13.ParsePayeeshipTransferred(log)
	case _KeeperRegistry13.abi.Events["PaymentWithdrawn"].ID:
		return _KeeperRegistry13.ParsePaymentWithdrawn(log)
	case _KeeperRegistry13.abi.Events["Unpaused"].ID:
		return _KeeperRegistry13.ParseUnpaused(log)
	case _KeeperRegistry13.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _KeeperRegistry13.ParseUpkeepAdminTransferRequested(log)
	case _KeeperRegistry13.abi.Events["UpkeepAdminTransferred"].ID:
		return _KeeperRegistry13.ParseUpkeepAdminTransferred(log)
	case _KeeperRegistry13.abi.Events["UpkeepCanceled"].ID:
		return _KeeperRegistry13.ParseUpkeepCanceled(log)
	case _KeeperRegistry13.abi.Events["UpkeepCheckDataUpdated"].ID:
		return _KeeperRegistry13.ParseUpkeepCheckDataUpdated(log)
	case _KeeperRegistry13.abi.Events["UpkeepGasLimitSet"].ID:
		return _KeeperRegistry13.ParseUpkeepGasLimitSet(log)
	case _KeeperRegistry13.abi.Events["UpkeepMigrated"].ID:
		return _KeeperRegistry13.ParseUpkeepMigrated(log)
	case _KeeperRegistry13.abi.Events["UpkeepPaused"].ID:
		return _KeeperRegistry13.ParseUpkeepPaused(log)
	case _KeeperRegistry13.abi.Events["UpkeepPerformed"].ID:
		return _KeeperRegistry13.ParseUpkeepPerformed(log)
	case _KeeperRegistry13.abi.Events["UpkeepReceived"].ID:
		return _KeeperRegistry13.ParseUpkeepReceived(log)
	case _KeeperRegistry13.abi.Events["UpkeepRegistered"].ID:
		return _KeeperRegistry13.ParseUpkeepRegistered(log)
	case _KeeperRegistry13.abi.Events["UpkeepUnpaused"].ID:
		return _KeeperRegistry13.ParseUpkeepUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeeperRegistry13ConfigSet) Topic() common.Hash {
	return common.HexToHash("0xfe125a41957477226ba20f85ef30a4024ea3bb8d066521ddc16df3f2944de325")
}

func (KeeperRegistry13FundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (KeeperRegistry13FundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (KeeperRegistry13KeepersUpdated) Topic() common.Hash {
	return common.HexToHash("0x056264c94f28bb06c99d13f0446eb96c67c215d8d707bce2655a98ddf1c0b71f")
}

func (KeeperRegistry13OwnerFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1")
}

func (KeeperRegistry13OwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (KeeperRegistry13OwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (KeeperRegistry13Paused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (KeeperRegistry13PayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (KeeperRegistry13PayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (KeeperRegistry13PaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (KeeperRegistry13Unpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (KeeperRegistry13UpkeepAdminTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35")
}

func (KeeperRegistry13UpkeepAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c")
}

func (KeeperRegistry13UpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (KeeperRegistry13UpkeepCheckDataUpdated) Topic() common.Hash {
	return common.HexToHash("0x7b778136e5211932b51a145badd01959415e79e051a933604b3d323f862dcabf")
}

func (KeeperRegistry13UpkeepGasLimitSet) Topic() common.Hash {
	return common.HexToHash("0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c")
}

func (KeeperRegistry13UpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (KeeperRegistry13UpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (KeeperRegistry13UpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xcaacad83e47cc45c280d487ec84184eee2fa3b54ebaa393bda7549f13da228f6")
}

func (KeeperRegistry13UpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (KeeperRegistry13UpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (KeeperRegistry13UpkeepUnpaused) Topic() common.Hash {
	return common.HexToHash("0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456")
}

func (_KeeperRegistry13 *KeeperRegistry13) Address() common.Address {
	return _KeeperRegistry13.address
}

type KeeperRegistry13Interface interface {
	ARBNITROORACLE(opts *bind.CallOpts) (common.Address, error)

	FASTGASFEED(opts *bind.CallOpts) (common.Address, error)

	KEEPERREGISTRYLOGIC(opts *bind.CallOpts) (common.Address, error)

	LINK(opts *bind.CallOpts) (common.Address, error)

	LINKETHFEED(opts *bind.CallOpts) (common.Address, error)

	OPTIMISMORACLE(opts *bind.CallOpts) (common.Address, error)

	PAYMENTMODEL(opts *bind.CallOpts) (uint8, error)

	REGISTRYGASOVERHEAD(opts *bind.CallOpts) (*big.Int, error)

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

	AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error)

	CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	CheckUpkeep(opts *bind.TransactOpts, id *big.Int, from common.Address) (*types.Transaction, error)

	MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	Pause(opts *bind.TransactOpts) (*types.Transaction, error)

	PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

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

	TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error)

	Unpause(opts *bind.TransactOpts) (*types.Transaction, error)

	UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	UpdateCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error)

	WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error)

	WithdrawOwnerFunds(opts *bind.TransactOpts) (*types.Transaction, error)

	WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error)

	Receive(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*KeeperRegistry13ConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13ConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*KeeperRegistry13ConfigSet, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistry13FundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13FundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*KeeperRegistry13FundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry13FundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13FundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*KeeperRegistry13FundsWithdrawn, error)

	FilterKeepersUpdated(opts *bind.FilterOpts) (*KeeperRegistry13KeepersUpdatedIterator, error)

	WatchKeepersUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13KeepersUpdated) (event.Subscription, error)

	ParseKeepersUpdated(log types.Log) (*KeeperRegistry13KeepersUpdated, error)

	FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*KeeperRegistry13OwnerFundsWithdrawnIterator, error)

	WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13OwnerFundsWithdrawn) (event.Subscription, error)

	ParseOwnerFundsWithdrawn(log types.Log) (*KeeperRegistry13OwnerFundsWithdrawn, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistry13OwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistry13OwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistry13OwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*KeeperRegistry13OwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*KeeperRegistry13PausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13Paused) (event.Subscription, error)

	ParsePaused(log types.Log) (*KeeperRegistry13Paused, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, keeper []common.Address, from []common.Address, to []common.Address) (*KeeperRegistry13PayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13PayeeshipTransferRequested, keeper []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistry13PayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, keeper []common.Address, from []common.Address, to []common.Address) (*KeeperRegistry13PayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13PayeeshipTransferred, keeper []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*KeeperRegistry13PayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, keeper []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistry13PaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13PaymentWithdrawn, keeper []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*KeeperRegistry13PaymentWithdrawn, error)

	FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistry13UnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13Unpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*KeeperRegistry13Unpaused, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistry13UpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*KeeperRegistry13UpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistry13UpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*KeeperRegistry13UpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistry13UpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*KeeperRegistry13UpkeepCanceled, error)

	FilterUpkeepCheckDataUpdated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry13UpkeepCheckDataUpdatedIterator, error)

	WatchUpkeepCheckDataUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepCheckDataUpdated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCheckDataUpdated(log types.Log) (*KeeperRegistry13UpkeepCheckDataUpdated, error)

	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry13UpkeepGasLimitSetIterator, error)

	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepGasLimitSet(log types.Log) (*KeeperRegistry13UpkeepGasLimitSet, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry13UpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*KeeperRegistry13UpkeepMigrated, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry13UpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*KeeperRegistry13UpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool, from []common.Address) (*KeeperRegistry13UpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepPerformed, id []*big.Int, success []bool, from []common.Address) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*KeeperRegistry13UpkeepPerformed, error)

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry13UpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*KeeperRegistry13UpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry13UpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*KeeperRegistry13UpkeepRegistered, error)

	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistry13UpkeepUnpausedIterator, error)

	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistry13UpkeepUnpaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepUnpaused(log types.Log) (*KeeperRegistry13UpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
