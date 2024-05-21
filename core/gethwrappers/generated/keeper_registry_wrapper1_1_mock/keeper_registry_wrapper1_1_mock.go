// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package keeper_registry_wrapper1_1_mock

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

var KeeperRegistryMockMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint24\",\"name\":\"blockCountPerTurn\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"}],\"name\":\"FlatFeeSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"keepers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"KeepersUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"RegistrarChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"maxLinkPayment\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"adjustedGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkEth\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"blockCountPerTurn\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"}],\"name\":\"emitConfigSet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"flatFeeMicroLink\",\"type\":\"uint32\"}],\"name\":\"emitFlatFeeSet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"emitFundsAdded\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitFundsWithdrawn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"keepers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"emitKeepersUpdated\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitOwnershipTransferRequested\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitOwnershipTransferred\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"emitPaused\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitPayeeshipTransferRequested\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitPayeeshipTransferred\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"emitPaymentWithdrawn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitRegistrarChanged\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"emitUnpaused\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"emitUpkeepCanceled\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"emitUpkeepPerformed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"emitUpkeepRegistered\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCanceledUpkeepList\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"blockCountPerTurn\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fallbackLinkPrice\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getKeeperList\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getMinBalanceForUpkeep\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getUpkeep\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"lastKeeper\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getUpkeepCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"_canceledUpkeepList\",\"type\":\"uint256[]\"}],\"name\":\"setCanceledUpkeepList\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"maxLinkPayment\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"adjustedGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkEth\",\"type\":\"uint256\"}],\"name\":\"setCheckUpkeepData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_paymentPremiumPPB\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_flatFeeMicroLink\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"_blockCountPerTurn\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"_checkGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"_stalenessSeconds\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"_gasCeilingMultiplier\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"_fallbackGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fallbackLinkPrice\",\"type\":\"uint256\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_keepers\",\"type\":\"address[]\"}],\"name\":\"setKeeperList\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"minBalance\",\"type\":\"uint96\"}],\"name\":\"setMinBalance\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"setPerformUpkeepSuccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"_executeGas\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"_balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"_maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"_lastKeeper\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_checkData\",\"type\":\"bytes\"}],\"name\":\"setUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_upkeepCount\",\"type\":\"uint256\"}],\"name\":\"setUpkeepCount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50611fb6806100206000396000f3fe608060405234801561001057600080fd5b50600436106101f05760003560e01c8063999a73bb1161010f578063c3f909d4116100a2578063db30a38611610071578063db30a38614610d76578063f7420bc214610dd7578063fecf27c914610e12578063ffc1d91c14610e2c576101f0565b8063c3f909d414610aed578063c41b813a14610b45578063c7c3a19a14610c12578063d5b16ded14610d31576101f0565b8063b019b4e8116100de578063b019b4e814610996578063b34362d7146109d1578063b657bc9c14610a0c578063c2030c8b14610a4a576101f0565b8063999a73bb146107b857806399e1a39b146108df5780639ec3ce4b14610924578063a6e95ed014610957576101f0565b806358e1e734116101875780637be5c756116101565780637be5c756146106b557806381d2c40c146106e8578063825bea391461073e5780638a8aa1651461076f576101f0565b806358e1e7341461056b57806367923e95146105b0578063749e9cc9146105fd5780637bbaf1ea1461062a576101f0565b80633e2d7056116101c35780633e2d7056146103815780634a16a9ad146103a65780634e6575e01461048a5780635181feaa146104ad576101f0565b806315a126ea146101f55780631ffe6c971461024d578063284403761461035c5780632cb6864d14610379575b600080fd5b6101fd610ecf565b60408051602080825283518183015283519192839290830191858101910280838360005b83811015610239578181015183820152602001610221565b505050509050019250505060405180910390f35b61035a600480360361010081101561026457600080fd5b81359173ffffffffffffffffffffffffffffffffffffffff602082013581169263ffffffff604084013516926bffffffffffffffffffffffff60608201351692608082013581169267ffffffffffffffff60a0840135169260c08101359092169190810190610100810160e08201356401000000008111156102e557600080fd5b8201836020820111156102f757600080fd5b8035906020019184600183028401116401000000008311171561031957600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610f3e945050505050565b005b61035a6004803603602081101561037257600080fd5b503561117b565b6101fd611180565b61035a6004803603604081101561039757600080fd5b508035906020013515156111d7565b61035a600480360360a08110156103bc57600080fd5b813591602081013515159173ffffffffffffffffffffffffffffffffffffffff604083013516916bffffffffffffffffffffffff6060820135169181019060a08101608082013564010000000081111561041557600080fd5b82018360208201111561042757600080fd5b8035906020019184600183028401116401000000008311171561044957600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550611215945050505050565b61035a600480360360208110156104a057600080fd5b503563ffffffff166112e6565b61035a600480360360c08110156104c357600080fd5b813591908101906040810160208201356401000000008111156104e557600080fd5b8201836020820111156104f757600080fd5b8035906020019184600183028401116401000000008311171561051957600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295505082359350505060208101359060408101359060600135611322565b61035a6004803603606081101561058157600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81358116916020810135821691604090910135166113a0565b61035a600480360360608110156105c657600080fd5b50803590602081013573ffffffffffffffffffffffffffffffffffffffff1690604001356bffffffffffffffffffffffff16611416565b61035a6004803603604081101561061357600080fd5b508035906020013567ffffffffffffffff16611476565b6106a16004803603604081101561064057600080fd5b8135919081019060408101602082013564010000000081111561066257600080fd5b82018360208201111561067457600080fd5b8035906020019184600183028401116401000000008311171561069657600080fd5b5090925090506114b1565b604080519115158252519081900360200190f35b61035a600480360360208110156106cb57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166114c8565b61035a600480360360e08110156106fe57600080fd5b5063ffffffff813581169162ffffff60208201358116926040830135169160608101359091169061ffff6080820135169060a08101359060c00135611514565b61035a6004803603604081101561075457600080fd5b50803590602001356bffffffffffffffffffffffff16611589565b61035a6004803603608081101561078557600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135811691602081013591604082013581169160600135166115d4565b61035a600480360360408110156107ce57600080fd5b8101906020810181356401000000008111156107e957600080fd5b8201836020820111156107fb57600080fd5b8035906020019184602083028401116401000000008311171561081d57600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600092019190915250929594936020810193503591505064010000000081111561086d57600080fd5b82018360208201111561087f57600080fd5b803590602001918460208302840111640100000000831117156108a157600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600092019190915250929550611656945050505050565b61035a600480360360608110156108f557600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135811691602081013582169160409091013516611715565b61035a6004803603602081101561093a57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff1661178b565b61035a6004803603606081101561096d57600080fd5b508035906020810135906040013573ffffffffffffffffffffffffffffffffffffffff166117d7565b61035a600480360360408110156109ac57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135811691602001351661182b565b61035a600480360360408110156109e757600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81358116916020013516611889565b610a2960048036036020811015610a2257600080fd5b50356118e7565b604080516bffffffffffffffffffffffff9092168252519081900360200190f35b61035a60048036036020811015610a6057600080fd5b810190602081018135640100000000811115610a7b57600080fd5b820183602082011115610a8d57600080fd5b80359060200191846020830284011164010000000083111715610aaf57600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600092019190915250929550611907945050505050565b610af561191e565b6040805163ffffffff988916815262ffffff9788166020820152959097168588015292909416606084015261ffff16608083015260a082019290925260c081019190915290519081900360e00190f35b610b7e60048036036040811015610b5b57600080fd5b508035906020013573ffffffffffffffffffffffffffffffffffffffff1661198e565b6040518080602001868152602001858152602001848152602001838152602001828103825287818151815260200191508051906020019080838360005b83811015610bd3578181015183820152602001610bbb565b50505050905090810190601f168015610c005780820380516001836020036101000a031916815260200191505b50965050505050505060405180910390f35b610c2f60048036036020811015610c2857600080fd5b5035611a76565b604051808873ffffffffffffffffffffffffffffffffffffffff1681526020018763ffffffff16815260200180602001866bffffffffffffffffffffffff1681526020018573ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1681526020018367ffffffffffffffff168152602001828103825287818151815260200191508051906020019080838360005b83811015610cf0578181015183820152602001610cd8565b50505050905090810190601f168015610d1d5780820380516001836020036101000a031916815260200191505b509850505050505050505060405180910390f35b61035a60048036036060811015610d4757600080fd5b50803590602081013563ffffffff16906040013573ffffffffffffffffffffffffffffffffffffffff16611c1f565b61035a6004803603610100811015610d8d57600080fd5b5063ffffffff8135811691602081013582169162ffffff604083013581169260608101359092169160808101359091169061ffff60a0820135169060c08101359060e00135611c79565b61035a60048036036040811015610ded57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81358116916020013516611ddd565b610e1a611e3b565b60408051918252519081900360200190f35b61035a60048036036020811015610e4257600080fd5b810190602081018135640100000000811115610e5d57600080fd5b820183602082011115610e6f57600080fd5b80359060200191846020830284011164010000000083111715610e9157600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600092019190915250929550611e41945050505050565b60606002805480602002602001604051908101604052809291908181526020018280548015610f3457602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610f09575b5050505050905090565b60006040518060c001604052808973ffffffffffffffffffffffffffffffffffffffff1681526020018863ffffffff168152602001876bffffffffffffffffffffffff1681526020018673ffffffffffffffffffffffffffffffffffffffff1681526020018567ffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff16815250905080600660008b815260200190815260200160002060008201518160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060208201518160000160146101000a81548163ffffffff021916908363ffffffff16021790555060408201518160010160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550606082015181600101600c6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060808201518160020160006101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555060a08201518160020160086101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505081600760008b8152602001908152602001600020908051906020019061116f929190611e54565b50505050505050505050565b600055565b60606001805480602002602001604051908101604052809291908181526020018280548015610f3457602002820191906000526020600020905b8154815260200190600101908083116111ba575050505050905090565b6000918252600a602052604090912080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055565b8273ffffffffffffffffffffffffffffffffffffffff16841515867fcaacad83e47cc45c280d487ec84184eee2fa3b54ebaa393bda7549f13da228f6858560405180836bffffffffffffffffffffffff16815260200180602001828103825283818151815260200191508051906020019080838360005b838110156112a457818101518382015260200161128c565b50505050905090810190601f1680156112d15780820380516001836020036101000a031916815260200191505b50935050505060405180910390a45050505050565b6040805163ffffffff8316815290517f17b46a44a823646eef686b7824df2962de896bc9a012a60b67694c5cbf184d8b9181900360200190a150565b6040805160a0810182528681526020808201879052818301869052606082018590526080820184905260008981526009825292909220815180519293919261136d9284920190611e54565b50602082015160018201556040820151600282015560608201516003820155608090910151600490910155505050505050565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a4505050565b604080516bffffffffffffffffffffffff83168152905173ffffffffffffffffffffffffffffffffffffffff84169185917fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039181900360200190a3505050565b60405167ffffffffffffffff82169083907f91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f79118190600090a35050565b50506000908152600a602052604090205460ff1690565b6040805173ffffffffffffffffffffffffffffffffffffffff8316815290517f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2589181900360200190a150565b6040805163ffffffff808a16825262ffffff808a166020840152908816828401528616606082015261ffff8516608082015260a0810184905260c0810183905290517feb3c06937e6595fd80ec1add18a195026d5cf65f122cc3ffedbfb18a9ed80b399181900360e00190a150505050505050565b60009182526008602052604090912080547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff909216919091179055565b8173ffffffffffffffffffffffffffffffffffffffff16838573ffffffffffffffffffffffffffffffffffffffff167f9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f4069884604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a450505050565b7f056264c94f28bb06c99d13f0446eb96c67c215d8d707bce2655a98ddf1c0b71f8282604051808060200180602001838103835285818151815260200191508051906020019060200280838360005b838110156116bd5781810151838201526020016116a5565b50505050905001838103825284818151815260200191508051906020019060200280838360005b838110156116fc5781810151838201526020016116e4565b5050505090500194505050505060405180910390a15050565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836760405160405180910390a4505050565b6040805173ffffffffffffffffffffffffffffffffffffffff8316815290517f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa9181900360200190a150565b6040805183815273ffffffffffffffffffffffffffffffffffffffff83166020820152815185927ff3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318928290030190a2505050565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167f9bf4a5b30267728df68663e14adb47e559863967c419dc6030638883408bed2e60405160405180910390a35050565b6000908152600860205260409020546bffffffffffffffffffffffff1690565b805161191a906001906020840190611ee0565b5050565b60035460045460055463ffffffff8084169468010000000000000000850462ffffff908116956b0100000000000000000000008104909316946f01000000000000000000000000000000840490911693720100000000000000000000000000000000000090930461ffff16929091565b60606000806000806000600960008981526020019081526020016000209050806000018160010154826002015483600301548460040154848054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015611a5a5780601f10611a2f57610100808354040283529160200191611a5a565b820191906000526020600020905b815481529060010190602001808311611a3d57829003601f168201915b5050505050945095509550955095509550509295509295909350565b6000818152600660209081526040808320815160c081018352815473ffffffffffffffffffffffffffffffffffffffff8082168084527401000000000000000000000000000000000000000090920463ffffffff168387018190526001808601546bffffffffffffffffffffffff81168689019081526c010000000000000000000000009091048416606080880191825260029889015467ffffffffffffffff811660808a019081526801000000000000000090910490961660a089019081528d8d5260078c528a8d20935190519251965184548c5161010097821615979097027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01169a909a04601f81018d90048d0286018d01909b528a85528c9b919a8c9a8b9a8b9a8b9a91999098909796949591939091879190830182828015611bfe5780601f10611bd357610100808354040283529160200191611bfe565b820191906000526020600020905b815481529060010190602001808311611be157829003601f168201915b50505050509450975097509750975097509750975050919395979092949650565b6040805163ffffffff8416815273ffffffffffffffffffffffffffffffffffffffff83166020820152815185927fbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012928290030190a2505050565b600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff998a16177fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff16640100000000988a1698909802979097177fffffffffffffffffffffffffffffffffffffffffff000000ffffffffffffffff166801000000000000000062ffffff97881602177fffffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffff166b0100000000000000000000009590981694909402969096177fffffffffffffffffffffffffffff000000ffffffffffffffffffffffffffffff166f010000000000000000000000000000009290941691909102929092177fffffffffffffffffffffffff0000ffffffffffffffffffffffffffffffffffff16720100000000000000000000000000000000000061ffff939093169290920291909117909155600491909155600555565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127860405160405180910390a35050565b60005490565b805161191a906002906020840190611f1a565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282611e8a5760008555611ed0565b82601f10611ea357805160ff1916838001178555611ed0565b82800160010185558215611ed0579182015b82811115611ed0578251825591602001919060010190611eb5565b50611edc929150611f94565b5090565b828054828255906000526020600020908101928215611ed05791602002820182811115611ed0578251825591602001919060010190611eb5565b828054828255906000526020600020908101928215611ed0579160200282015b82811115611ed057825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190611f3a565b5b80821115611edc5760008155600101611f9556fea164736f6c6343000706000a",
}

var KeeperRegistryMockABI = KeeperRegistryMockMetaData.ABI

var KeeperRegistryMockBin = KeeperRegistryMockMetaData.Bin

func DeployKeeperRegistryMock(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *KeeperRegistryMock, error) {
	parsed, err := KeeperRegistryMockMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistryMockBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperRegistryMock{address: address, abi: *parsed, KeeperRegistryMockCaller: KeeperRegistryMockCaller{contract: contract}, KeeperRegistryMockTransactor: KeeperRegistryMockTransactor{contract: contract}, KeeperRegistryMockFilterer: KeeperRegistryMockFilterer{contract: contract}}, nil
}

type KeeperRegistryMock struct {
	address common.Address
	abi     abi.ABI
	KeeperRegistryMockCaller
	KeeperRegistryMockTransactor
	KeeperRegistryMockFilterer
}

type KeeperRegistryMockCaller struct {
	contract *bind.BoundContract
}

type KeeperRegistryMockTransactor struct {
	contract *bind.BoundContract
}

type KeeperRegistryMockFilterer struct {
	contract *bind.BoundContract
}

type KeeperRegistryMockSession struct {
	Contract     *KeeperRegistryMock
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeeperRegistryMockCallerSession struct {
	Contract *KeeperRegistryMockCaller
	CallOpts bind.CallOpts
}

type KeeperRegistryMockTransactorSession struct {
	Contract     *KeeperRegistryMockTransactor
	TransactOpts bind.TransactOpts
}

type KeeperRegistryMockRaw struct {
	Contract *KeeperRegistryMock
}

type KeeperRegistryMockCallerRaw struct {
	Contract *KeeperRegistryMockCaller
}

type KeeperRegistryMockTransactorRaw struct {
	Contract *KeeperRegistryMockTransactor
}

func NewKeeperRegistryMock(address common.Address, backend bind.ContractBackend) (*KeeperRegistryMock, error) {
	abi, err := abi.JSON(strings.NewReader(KeeperRegistryMockABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeeperRegistryMock(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryMock{address: address, abi: abi, KeeperRegistryMockCaller: KeeperRegistryMockCaller{contract: contract}, KeeperRegistryMockTransactor: KeeperRegistryMockTransactor{contract: contract}, KeeperRegistryMockFilterer: KeeperRegistryMockFilterer{contract: contract}}, nil
}

func NewKeeperRegistryMockCaller(address common.Address, caller bind.ContractCaller) (*KeeperRegistryMockCaller, error) {
	contract, err := bindKeeperRegistryMock(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryMockCaller{contract: contract}, nil
}

func NewKeeperRegistryMockTransactor(address common.Address, transactor bind.ContractTransactor) (*KeeperRegistryMockTransactor, error) {
	contract, err := bindKeeperRegistryMock(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryMockTransactor{contract: contract}, nil
}

func NewKeeperRegistryMockFilterer(address common.Address, filterer bind.ContractFilterer) (*KeeperRegistryMockFilterer, error) {
	contract, err := bindKeeperRegistryMock(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryMockFilterer{contract: contract}, nil
}

func bindKeeperRegistryMock(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeeperRegistryMockMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_KeeperRegistryMock *KeeperRegistryMockRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistryMock.Contract.KeeperRegistryMockCaller.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistryMock *KeeperRegistryMockRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.KeeperRegistryMockTransactor.contract.Transfer(opts)
}

func (_KeeperRegistryMock *KeeperRegistryMockRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.KeeperRegistryMockTransactor.contract.Transact(opts, method, params...)
}

func (_KeeperRegistryMock *KeeperRegistryMockCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistryMock.Contract.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.contract.Transfer(opts)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.contract.Transact(opts, method, params...)
}

func (_KeeperRegistryMock *KeeperRegistryMockCaller) CheckUpkeep(opts *bind.CallOpts, id *big.Int, from common.Address) (CheckUpkeep,

	error) {
	var out []interface{}
	err := _KeeperRegistryMock.contract.Call(opts, &out, "checkUpkeep", id, from)

	outstruct := new(CheckUpkeep)
	if err != nil {
		return *outstruct, err
	}

	outstruct.PerformData = *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	outstruct.MaxLinkPayment = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.GasLimit = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.AdjustedGasWei = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.LinkEth = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_KeeperRegistryMock *KeeperRegistryMockSession) CheckUpkeep(id *big.Int, from common.Address) (CheckUpkeep,

	error) {
	return _KeeperRegistryMock.Contract.CheckUpkeep(&_KeeperRegistryMock.CallOpts, id, from)
}

func (_KeeperRegistryMock *KeeperRegistryMockCallerSession) CheckUpkeep(id *big.Int, from common.Address) (CheckUpkeep,

	error) {
	return _KeeperRegistryMock.Contract.CheckUpkeep(&_KeeperRegistryMock.CallOpts, id, from)
}

func (_KeeperRegistryMock *KeeperRegistryMockCaller) GetCanceledUpkeepList(opts *bind.CallOpts) ([]*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistryMock.contract.Call(opts, &out, "getCanceledUpkeepList")

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_KeeperRegistryMock *KeeperRegistryMockSession) GetCanceledUpkeepList() ([]*big.Int, error) {
	return _KeeperRegistryMock.Contract.GetCanceledUpkeepList(&_KeeperRegistryMock.CallOpts)
}

func (_KeeperRegistryMock *KeeperRegistryMockCallerSession) GetCanceledUpkeepList() ([]*big.Int, error) {
	return _KeeperRegistryMock.Contract.GetCanceledUpkeepList(&_KeeperRegistryMock.CallOpts)
}

func (_KeeperRegistryMock *KeeperRegistryMockCaller) GetConfig(opts *bind.CallOpts) (GetConfig,

	error) {
	var out []interface{}
	err := _KeeperRegistryMock.contract.Call(opts, &out, "getConfig")

	outstruct := new(GetConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.PaymentPremiumPPB = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockCountPerTurn = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.CheckGasLimit = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.StalenessSeconds = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.GasCeilingMultiplier = *abi.ConvertType(out[4], new(uint16)).(*uint16)
	outstruct.FallbackGasPrice = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.FallbackLinkPrice = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_KeeperRegistryMock *KeeperRegistryMockSession) GetConfig() (GetConfig,

	error) {
	return _KeeperRegistryMock.Contract.GetConfig(&_KeeperRegistryMock.CallOpts)
}

func (_KeeperRegistryMock *KeeperRegistryMockCallerSession) GetConfig() (GetConfig,

	error) {
	return _KeeperRegistryMock.Contract.GetConfig(&_KeeperRegistryMock.CallOpts)
}

func (_KeeperRegistryMock *KeeperRegistryMockCaller) GetKeeperList(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryMock.contract.Call(opts, &out, "getKeeperList")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_KeeperRegistryMock *KeeperRegistryMockSession) GetKeeperList() ([]common.Address, error) {
	return _KeeperRegistryMock.Contract.GetKeeperList(&_KeeperRegistryMock.CallOpts)
}

func (_KeeperRegistryMock *KeeperRegistryMockCallerSession) GetKeeperList() ([]common.Address, error) {
	return _KeeperRegistryMock.Contract.GetKeeperList(&_KeeperRegistryMock.CallOpts)
}

func (_KeeperRegistryMock *KeeperRegistryMockCaller) GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistryMock.contract.Call(opts, &out, "getMinBalanceForUpkeep", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistryMock *KeeperRegistryMockSession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _KeeperRegistryMock.Contract.GetMinBalanceForUpkeep(&_KeeperRegistryMock.CallOpts, id)
}

func (_KeeperRegistryMock *KeeperRegistryMockCallerSession) GetMinBalanceForUpkeep(id *big.Int) (*big.Int, error) {
	return _KeeperRegistryMock.Contract.GetMinBalanceForUpkeep(&_KeeperRegistryMock.CallOpts, id)
}

func (_KeeperRegistryMock *KeeperRegistryMockCaller) GetUpkeep(opts *bind.CallOpts, id *big.Int) (GetUpkeep,

	error) {
	var out []interface{}
	err := _KeeperRegistryMock.contract.Call(opts, &out, "getUpkeep", id)

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

	return *outstruct, err

}

func (_KeeperRegistryMock *KeeperRegistryMockSession) GetUpkeep(id *big.Int) (GetUpkeep,

	error) {
	return _KeeperRegistryMock.Contract.GetUpkeep(&_KeeperRegistryMock.CallOpts, id)
}

func (_KeeperRegistryMock *KeeperRegistryMockCallerSession) GetUpkeep(id *big.Int) (GetUpkeep,

	error) {
	return _KeeperRegistryMock.Contract.GetUpkeep(&_KeeperRegistryMock.CallOpts, id)
}

func (_KeeperRegistryMock *KeeperRegistryMockCaller) GetUpkeepCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KeeperRegistryMock.contract.Call(opts, &out, "getUpkeepCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_KeeperRegistryMock *KeeperRegistryMockSession) GetUpkeepCount() (*big.Int, error) {
	return _KeeperRegistryMock.Contract.GetUpkeepCount(&_KeeperRegistryMock.CallOpts)
}

func (_KeeperRegistryMock *KeeperRegistryMockCallerSession) GetUpkeepCount() (*big.Int, error) {
	return _KeeperRegistryMock.Contract.GetUpkeepCount(&_KeeperRegistryMock.CallOpts)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactor) EmitConfigSet(opts *bind.TransactOpts, paymentPremiumPPB uint32, blockCountPerTurn *big.Int, checkGasLimit uint32, stalenessSeconds *big.Int, gasCeilingMultiplier uint16, fallbackGasPrice *big.Int, fallbackLinkPrice *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryMock.contract.Transact(opts, "emitConfigSet", paymentPremiumPPB, blockCountPerTurn, checkGasLimit, stalenessSeconds, gasCeilingMultiplier, fallbackGasPrice, fallbackLinkPrice)
}

func (_KeeperRegistryMock *KeeperRegistryMockSession) EmitConfigSet(paymentPremiumPPB uint32, blockCountPerTurn *big.Int, checkGasLimit uint32, stalenessSeconds *big.Int, gasCeilingMultiplier uint16, fallbackGasPrice *big.Int, fallbackLinkPrice *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitConfigSet(&_KeeperRegistryMock.TransactOpts, paymentPremiumPPB, blockCountPerTurn, checkGasLimit, stalenessSeconds, gasCeilingMultiplier, fallbackGasPrice, fallbackLinkPrice)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorSession) EmitConfigSet(paymentPremiumPPB uint32, blockCountPerTurn *big.Int, checkGasLimit uint32, stalenessSeconds *big.Int, gasCeilingMultiplier uint16, fallbackGasPrice *big.Int, fallbackLinkPrice *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitConfigSet(&_KeeperRegistryMock.TransactOpts, paymentPremiumPPB, blockCountPerTurn, checkGasLimit, stalenessSeconds, gasCeilingMultiplier, fallbackGasPrice, fallbackLinkPrice)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactor) EmitFlatFeeSet(opts *bind.TransactOpts, flatFeeMicroLink uint32) (*types.Transaction, error) {
	return _KeeperRegistryMock.contract.Transact(opts, "emitFlatFeeSet", flatFeeMicroLink)
}

func (_KeeperRegistryMock *KeeperRegistryMockSession) EmitFlatFeeSet(flatFeeMicroLink uint32) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitFlatFeeSet(&_KeeperRegistryMock.TransactOpts, flatFeeMicroLink)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorSession) EmitFlatFeeSet(flatFeeMicroLink uint32) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitFlatFeeSet(&_KeeperRegistryMock.TransactOpts, flatFeeMicroLink)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactor) EmitFundsAdded(opts *bind.TransactOpts, id *big.Int, from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryMock.contract.Transact(opts, "emitFundsAdded", id, from, amount)
}

func (_KeeperRegistryMock *KeeperRegistryMockSession) EmitFundsAdded(id *big.Int, from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitFundsAdded(&_KeeperRegistryMock.TransactOpts, id, from, amount)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorSession) EmitFundsAdded(id *big.Int, from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitFundsAdded(&_KeeperRegistryMock.TransactOpts, id, from, amount)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactor) EmitFundsWithdrawn(opts *bind.TransactOpts, id *big.Int, amount *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.contract.Transact(opts, "emitFundsWithdrawn", id, amount, to)
}

func (_KeeperRegistryMock *KeeperRegistryMockSession) EmitFundsWithdrawn(id *big.Int, amount *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitFundsWithdrawn(&_KeeperRegistryMock.TransactOpts, id, amount, to)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorSession) EmitFundsWithdrawn(id *big.Int, amount *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitFundsWithdrawn(&_KeeperRegistryMock.TransactOpts, id, amount, to)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactor) EmitKeepersUpdated(opts *bind.TransactOpts, keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.contract.Transact(opts, "emitKeepersUpdated", keepers, payees)
}

func (_KeeperRegistryMock *KeeperRegistryMockSession) EmitKeepersUpdated(keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitKeepersUpdated(&_KeeperRegistryMock.TransactOpts, keepers, payees)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorSession) EmitKeepersUpdated(keepers []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitKeepersUpdated(&_KeeperRegistryMock.TransactOpts, keepers, payees)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactor) EmitOwnershipTransferRequested(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.contract.Transact(opts, "emitOwnershipTransferRequested", from, to)
}

func (_KeeperRegistryMock *KeeperRegistryMockSession) EmitOwnershipTransferRequested(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitOwnershipTransferRequested(&_KeeperRegistryMock.TransactOpts, from, to)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorSession) EmitOwnershipTransferRequested(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitOwnershipTransferRequested(&_KeeperRegistryMock.TransactOpts, from, to)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactor) EmitOwnershipTransferred(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.contract.Transact(opts, "emitOwnershipTransferred", from, to)
}

func (_KeeperRegistryMock *KeeperRegistryMockSession) EmitOwnershipTransferred(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitOwnershipTransferred(&_KeeperRegistryMock.TransactOpts, from, to)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorSession) EmitOwnershipTransferred(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitOwnershipTransferred(&_KeeperRegistryMock.TransactOpts, from, to)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactor) EmitPaused(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.contract.Transact(opts, "emitPaused", account)
}

func (_KeeperRegistryMock *KeeperRegistryMockSession) EmitPaused(account common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitPaused(&_KeeperRegistryMock.TransactOpts, account)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorSession) EmitPaused(account common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitPaused(&_KeeperRegistryMock.TransactOpts, account)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactor) EmitPayeeshipTransferRequested(opts *bind.TransactOpts, keeper common.Address, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.contract.Transact(opts, "emitPayeeshipTransferRequested", keeper, from, to)
}

func (_KeeperRegistryMock *KeeperRegistryMockSession) EmitPayeeshipTransferRequested(keeper common.Address, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitPayeeshipTransferRequested(&_KeeperRegistryMock.TransactOpts, keeper, from, to)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorSession) EmitPayeeshipTransferRequested(keeper common.Address, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitPayeeshipTransferRequested(&_KeeperRegistryMock.TransactOpts, keeper, from, to)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactor) EmitPayeeshipTransferred(opts *bind.TransactOpts, keeper common.Address, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.contract.Transact(opts, "emitPayeeshipTransferred", keeper, from, to)
}

func (_KeeperRegistryMock *KeeperRegistryMockSession) EmitPayeeshipTransferred(keeper common.Address, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitPayeeshipTransferred(&_KeeperRegistryMock.TransactOpts, keeper, from, to)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorSession) EmitPayeeshipTransferred(keeper common.Address, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitPayeeshipTransferred(&_KeeperRegistryMock.TransactOpts, keeper, from, to)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactor) EmitPaymentWithdrawn(opts *bind.TransactOpts, keeper common.Address, amount *big.Int, to common.Address, payee common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.contract.Transact(opts, "emitPaymentWithdrawn", keeper, amount, to, payee)
}

func (_KeeperRegistryMock *KeeperRegistryMockSession) EmitPaymentWithdrawn(keeper common.Address, amount *big.Int, to common.Address, payee common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitPaymentWithdrawn(&_KeeperRegistryMock.TransactOpts, keeper, amount, to, payee)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorSession) EmitPaymentWithdrawn(keeper common.Address, amount *big.Int, to common.Address, payee common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitPaymentWithdrawn(&_KeeperRegistryMock.TransactOpts, keeper, amount, to, payee)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactor) EmitRegistrarChanged(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.contract.Transact(opts, "emitRegistrarChanged", from, to)
}

func (_KeeperRegistryMock *KeeperRegistryMockSession) EmitRegistrarChanged(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitRegistrarChanged(&_KeeperRegistryMock.TransactOpts, from, to)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorSession) EmitRegistrarChanged(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitRegistrarChanged(&_KeeperRegistryMock.TransactOpts, from, to)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactor) EmitUnpaused(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.contract.Transact(opts, "emitUnpaused", account)
}

func (_KeeperRegistryMock *KeeperRegistryMockSession) EmitUnpaused(account common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitUnpaused(&_KeeperRegistryMock.TransactOpts, account)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorSession) EmitUnpaused(account common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitUnpaused(&_KeeperRegistryMock.TransactOpts, account)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactor) EmitUpkeepCanceled(opts *bind.TransactOpts, id *big.Int, atBlockHeight uint64) (*types.Transaction, error) {
	return _KeeperRegistryMock.contract.Transact(opts, "emitUpkeepCanceled", id, atBlockHeight)
}

func (_KeeperRegistryMock *KeeperRegistryMockSession) EmitUpkeepCanceled(id *big.Int, atBlockHeight uint64) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitUpkeepCanceled(&_KeeperRegistryMock.TransactOpts, id, atBlockHeight)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorSession) EmitUpkeepCanceled(id *big.Int, atBlockHeight uint64) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitUpkeepCanceled(&_KeeperRegistryMock.TransactOpts, id, atBlockHeight)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactor) EmitUpkeepPerformed(opts *bind.TransactOpts, id *big.Int, success bool, from common.Address, payment *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistryMock.contract.Transact(opts, "emitUpkeepPerformed", id, success, from, payment, performData)
}

func (_KeeperRegistryMock *KeeperRegistryMockSession) EmitUpkeepPerformed(id *big.Int, success bool, from common.Address, payment *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitUpkeepPerformed(&_KeeperRegistryMock.TransactOpts, id, success, from, payment, performData)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorSession) EmitUpkeepPerformed(id *big.Int, success bool, from common.Address, payment *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitUpkeepPerformed(&_KeeperRegistryMock.TransactOpts, id, success, from, payment, performData)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactor) EmitUpkeepRegistered(opts *bind.TransactOpts, id *big.Int, executeGas uint32, admin common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.contract.Transact(opts, "emitUpkeepRegistered", id, executeGas, admin)
}

func (_KeeperRegistryMock *KeeperRegistryMockSession) EmitUpkeepRegistered(id *big.Int, executeGas uint32, admin common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitUpkeepRegistered(&_KeeperRegistryMock.TransactOpts, id, executeGas, admin)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorSession) EmitUpkeepRegistered(id *big.Int, executeGas uint32, admin common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.EmitUpkeepRegistered(&_KeeperRegistryMock.TransactOpts, id, executeGas, admin)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactor) PerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistryMock.contract.Transact(opts, "performUpkeep", id, performData)
}

func (_KeeperRegistryMock *KeeperRegistryMockSession) PerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.PerformUpkeep(&_KeeperRegistryMock.TransactOpts, id, performData)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorSession) PerformUpkeep(id *big.Int, performData []byte) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.PerformUpkeep(&_KeeperRegistryMock.TransactOpts, id, performData)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactor) SetCanceledUpkeepList(opts *bind.TransactOpts, _canceledUpkeepList []*big.Int) (*types.Transaction, error) {
	return _KeeperRegistryMock.contract.Transact(opts, "setCanceledUpkeepList", _canceledUpkeepList)
}

func (_KeeperRegistryMock *KeeperRegistryMockSession) SetCanceledUpkeepList(_canceledUpkeepList []*big.Int) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.SetCanceledUpkeepList(&_KeeperRegistryMock.TransactOpts, _canceledUpkeepList)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorSession) SetCanceledUpkeepList(_canceledUpkeepList []*big.Int) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.SetCanceledUpkeepList(&_KeeperRegistryMock.TransactOpts, _canceledUpkeepList)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactor) SetCheckUpkeepData(opts *bind.TransactOpts, id *big.Int, performData []byte, maxLinkPayment *big.Int, gasLimit *big.Int, adjustedGasWei *big.Int, linkEth *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryMock.contract.Transact(opts, "setCheckUpkeepData", id, performData, maxLinkPayment, gasLimit, adjustedGasWei, linkEth)
}

func (_KeeperRegistryMock *KeeperRegistryMockSession) SetCheckUpkeepData(id *big.Int, performData []byte, maxLinkPayment *big.Int, gasLimit *big.Int, adjustedGasWei *big.Int, linkEth *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.SetCheckUpkeepData(&_KeeperRegistryMock.TransactOpts, id, performData, maxLinkPayment, gasLimit, adjustedGasWei, linkEth)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorSession) SetCheckUpkeepData(id *big.Int, performData []byte, maxLinkPayment *big.Int, gasLimit *big.Int, adjustedGasWei *big.Int, linkEth *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.SetCheckUpkeepData(&_KeeperRegistryMock.TransactOpts, id, performData, maxLinkPayment, gasLimit, adjustedGasWei, linkEth)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactor) SetConfig(opts *bind.TransactOpts, _paymentPremiumPPB uint32, _flatFeeMicroLink uint32, _blockCountPerTurn *big.Int, _checkGasLimit uint32, _stalenessSeconds *big.Int, _gasCeilingMultiplier uint16, _fallbackGasPrice *big.Int, _fallbackLinkPrice *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryMock.contract.Transact(opts, "setConfig", _paymentPremiumPPB, _flatFeeMicroLink, _blockCountPerTurn, _checkGasLimit, _stalenessSeconds, _gasCeilingMultiplier, _fallbackGasPrice, _fallbackLinkPrice)
}

func (_KeeperRegistryMock *KeeperRegistryMockSession) SetConfig(_paymentPremiumPPB uint32, _flatFeeMicroLink uint32, _blockCountPerTurn *big.Int, _checkGasLimit uint32, _stalenessSeconds *big.Int, _gasCeilingMultiplier uint16, _fallbackGasPrice *big.Int, _fallbackLinkPrice *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.SetConfig(&_KeeperRegistryMock.TransactOpts, _paymentPremiumPPB, _flatFeeMicroLink, _blockCountPerTurn, _checkGasLimit, _stalenessSeconds, _gasCeilingMultiplier, _fallbackGasPrice, _fallbackLinkPrice)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorSession) SetConfig(_paymentPremiumPPB uint32, _flatFeeMicroLink uint32, _blockCountPerTurn *big.Int, _checkGasLimit uint32, _stalenessSeconds *big.Int, _gasCeilingMultiplier uint16, _fallbackGasPrice *big.Int, _fallbackLinkPrice *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.SetConfig(&_KeeperRegistryMock.TransactOpts, _paymentPremiumPPB, _flatFeeMicroLink, _blockCountPerTurn, _checkGasLimit, _stalenessSeconds, _gasCeilingMultiplier, _fallbackGasPrice, _fallbackLinkPrice)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactor) SetKeeperList(opts *bind.TransactOpts, _keepers []common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.contract.Transact(opts, "setKeeperList", _keepers)
}

func (_KeeperRegistryMock *KeeperRegistryMockSession) SetKeeperList(_keepers []common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.SetKeeperList(&_KeeperRegistryMock.TransactOpts, _keepers)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorSession) SetKeeperList(_keepers []common.Address) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.SetKeeperList(&_KeeperRegistryMock.TransactOpts, _keepers)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactor) SetMinBalance(opts *bind.TransactOpts, id *big.Int, minBalance *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryMock.contract.Transact(opts, "setMinBalance", id, minBalance)
}

func (_KeeperRegistryMock *KeeperRegistryMockSession) SetMinBalance(id *big.Int, minBalance *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.SetMinBalance(&_KeeperRegistryMock.TransactOpts, id, minBalance)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorSession) SetMinBalance(id *big.Int, minBalance *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.SetMinBalance(&_KeeperRegistryMock.TransactOpts, id, minBalance)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactor) SetPerformUpkeepSuccess(opts *bind.TransactOpts, id *big.Int, success bool) (*types.Transaction, error) {
	return _KeeperRegistryMock.contract.Transact(opts, "setPerformUpkeepSuccess", id, success)
}

func (_KeeperRegistryMock *KeeperRegistryMockSession) SetPerformUpkeepSuccess(id *big.Int, success bool) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.SetPerformUpkeepSuccess(&_KeeperRegistryMock.TransactOpts, id, success)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorSession) SetPerformUpkeepSuccess(id *big.Int, success bool) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.SetPerformUpkeepSuccess(&_KeeperRegistryMock.TransactOpts, id, success)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactor) SetUpkeep(opts *bind.TransactOpts, id *big.Int, _target common.Address, _executeGas uint32, _balance *big.Int, _admin common.Address, _maxValidBlocknumber uint64, _lastKeeper common.Address, _checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistryMock.contract.Transact(opts, "setUpkeep", id, _target, _executeGas, _balance, _admin, _maxValidBlocknumber, _lastKeeper, _checkData)
}

func (_KeeperRegistryMock *KeeperRegistryMockSession) SetUpkeep(id *big.Int, _target common.Address, _executeGas uint32, _balance *big.Int, _admin common.Address, _maxValidBlocknumber uint64, _lastKeeper common.Address, _checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.SetUpkeep(&_KeeperRegistryMock.TransactOpts, id, _target, _executeGas, _balance, _admin, _maxValidBlocknumber, _lastKeeper, _checkData)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorSession) SetUpkeep(id *big.Int, _target common.Address, _executeGas uint32, _balance *big.Int, _admin common.Address, _maxValidBlocknumber uint64, _lastKeeper common.Address, _checkData []byte) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.SetUpkeep(&_KeeperRegistryMock.TransactOpts, id, _target, _executeGas, _balance, _admin, _maxValidBlocknumber, _lastKeeper, _checkData)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactor) SetUpkeepCount(opts *bind.TransactOpts, _upkeepCount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryMock.contract.Transact(opts, "setUpkeepCount", _upkeepCount)
}

func (_KeeperRegistryMock *KeeperRegistryMockSession) SetUpkeepCount(_upkeepCount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.SetUpkeepCount(&_KeeperRegistryMock.TransactOpts, _upkeepCount)
}

func (_KeeperRegistryMock *KeeperRegistryMockTransactorSession) SetUpkeepCount(_upkeepCount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryMock.Contract.SetUpkeepCount(&_KeeperRegistryMock.TransactOpts, _upkeepCount)
}

type KeeperRegistryMockConfigSetIterator struct {
	Event *KeeperRegistryMockConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryMockConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryMockConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryMockConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryMockConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryMockConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryMockConfigSet struct {
	PaymentPremiumPPB    uint32
	BlockCountPerTurn    *big.Int
	CheckGasLimit        uint32
	StalenessSeconds     *big.Int
	GasCeilingMultiplier uint16
	FallbackGasPrice     *big.Int
	FallbackLinkPrice    *big.Int
	Raw                  types.Log
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) FilterConfigSet(opts *bind.FilterOpts) (*KeeperRegistryMockConfigSetIterator, error) {

	logs, sub, err := _KeeperRegistryMock.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryMockConfigSetIterator{contract: _KeeperRegistryMock.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockConfigSet) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryMock.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryMockConfigSet)
				if err := _KeeperRegistryMock.contract.UnpackLog(event, "ConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) ParseConfigSet(log types.Log) (*KeeperRegistryMockConfigSet, error) {
	event := new(KeeperRegistryMockConfigSet)
	if err := _KeeperRegistryMock.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryMockFlatFeeSetIterator struct {
	Event *KeeperRegistryMockFlatFeeSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryMockFlatFeeSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryMockFlatFeeSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryMockFlatFeeSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryMockFlatFeeSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryMockFlatFeeSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryMockFlatFeeSet struct {
	FlatFeeMicroLink uint32
	Raw              types.Log
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) FilterFlatFeeSet(opts *bind.FilterOpts) (*KeeperRegistryMockFlatFeeSetIterator, error) {

	logs, sub, err := _KeeperRegistryMock.contract.FilterLogs(opts, "FlatFeeSet")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryMockFlatFeeSetIterator{contract: _KeeperRegistryMock.contract, event: "FlatFeeSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) WatchFlatFeeSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockFlatFeeSet) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryMock.contract.WatchLogs(opts, "FlatFeeSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryMockFlatFeeSet)
				if err := _KeeperRegistryMock.contract.UnpackLog(event, "FlatFeeSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) ParseFlatFeeSet(log types.Log) (*KeeperRegistryMockFlatFeeSet, error) {
	event := new(KeeperRegistryMockFlatFeeSet)
	if err := _KeeperRegistryMock.contract.UnpackLog(event, "FlatFeeSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryMockFundsAddedIterator struct {
	Event *KeeperRegistryMockFundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryMockFundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryMockFundsAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryMockFundsAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryMockFundsAddedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryMockFundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryMockFundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistryMockFundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistryMock.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryMockFundsAddedIterator{contract: _KeeperRegistryMock.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistryMock.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryMockFundsAdded)
				if err := _KeeperRegistryMock.contract.UnpackLog(event, "FundsAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) ParseFundsAdded(log types.Log) (*KeeperRegistryMockFundsAdded, error) {
	event := new(KeeperRegistryMockFundsAdded)
	if err := _KeeperRegistryMock.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryMockFundsWithdrawnIterator struct {
	Event *KeeperRegistryMockFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryMockFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryMockFundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryMockFundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryMockFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryMockFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryMockFundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryMockFundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryMock.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryMockFundsWithdrawnIterator{contract: _KeeperRegistryMock.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockFundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryMock.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryMockFundsWithdrawn)
				if err := _KeeperRegistryMock.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) ParseFundsWithdrawn(log types.Log) (*KeeperRegistryMockFundsWithdrawn, error) {
	event := new(KeeperRegistryMockFundsWithdrawn)
	if err := _KeeperRegistryMock.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryMockKeepersUpdatedIterator struct {
	Event *KeeperRegistryMockKeepersUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryMockKeepersUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryMockKeepersUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryMockKeepersUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryMockKeepersUpdatedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryMockKeepersUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryMockKeepersUpdated struct {
	Keepers []common.Address
	Payees  []common.Address
	Raw     types.Log
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) FilterKeepersUpdated(opts *bind.FilterOpts) (*KeeperRegistryMockKeepersUpdatedIterator, error) {

	logs, sub, err := _KeeperRegistryMock.contract.FilterLogs(opts, "KeepersUpdated")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryMockKeepersUpdatedIterator{contract: _KeeperRegistryMock.contract, event: "KeepersUpdated", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) WatchKeepersUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockKeepersUpdated) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryMock.contract.WatchLogs(opts, "KeepersUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryMockKeepersUpdated)
				if err := _KeeperRegistryMock.contract.UnpackLog(event, "KeepersUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) ParseKeepersUpdated(log types.Log) (*KeeperRegistryMockKeepersUpdated, error) {
	event := new(KeeperRegistryMockKeepersUpdated)
	if err := _KeeperRegistryMock.contract.UnpackLog(event, "KeepersUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryMockOwnershipTransferRequestedIterator struct {
	Event *KeeperRegistryMockOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryMockOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryMockOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryMockOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryMockOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryMockOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryMockOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryMockOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryMock.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryMockOwnershipTransferRequestedIterator{contract: _KeeperRegistryMock.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryMock.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryMockOwnershipTransferRequested)
				if err := _KeeperRegistryMock.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryMockOwnershipTransferRequested, error) {
	event := new(KeeperRegistryMockOwnershipTransferRequested)
	if err := _KeeperRegistryMock.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryMockOwnershipTransferredIterator struct {
	Event *KeeperRegistryMockOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryMockOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryMockOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryMockOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryMockOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryMockOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryMockOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryMockOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryMock.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryMockOwnershipTransferredIterator{contract: _KeeperRegistryMock.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryMock.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryMockOwnershipTransferred)
				if err := _KeeperRegistryMock.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) ParseOwnershipTransferred(log types.Log) (*KeeperRegistryMockOwnershipTransferred, error) {
	event := new(KeeperRegistryMockOwnershipTransferred)
	if err := _KeeperRegistryMock.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryMockPausedIterator struct {
	Event *KeeperRegistryMockPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryMockPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryMockPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryMockPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryMockPausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryMockPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryMockPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) FilterPaused(opts *bind.FilterOpts) (*KeeperRegistryMockPausedIterator, error) {

	logs, sub, err := _KeeperRegistryMock.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryMockPausedIterator{contract: _KeeperRegistryMock.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockPaused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryMock.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryMockPaused)
				if err := _KeeperRegistryMock.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) ParsePaused(log types.Log) (*KeeperRegistryMockPaused, error) {
	event := new(KeeperRegistryMockPaused)
	if err := _KeeperRegistryMock.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryMockPayeeshipTransferRequestedIterator struct {
	Event *KeeperRegistryMockPayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryMockPayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryMockPayeeshipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryMockPayeeshipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryMockPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryMockPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryMockPayeeshipTransferRequested struct {
	Keeper common.Address
	From   common.Address
	To     common.Address
	Raw    types.Log
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, keeper []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryMockPayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistryMock.contract.FilterLogs(opts, "PayeeshipTransferRequested", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryMockPayeeshipTransferRequestedIterator{contract: _KeeperRegistryMock.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockPayeeshipTransferRequested, keeper []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryMock.contract.WatchLogs(opts, "PayeeshipTransferRequested", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryMockPayeeshipTransferRequested)
				if err := _KeeperRegistryMock.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistryMockPayeeshipTransferRequested, error) {
	event := new(KeeperRegistryMockPayeeshipTransferRequested)
	if err := _KeeperRegistryMock.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryMockPayeeshipTransferredIterator struct {
	Event *KeeperRegistryMockPayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryMockPayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryMockPayeeshipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryMockPayeeshipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryMockPayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryMockPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryMockPayeeshipTransferred struct {
	Keeper common.Address
	From   common.Address
	To     common.Address
	Raw    types.Log
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, keeper []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryMockPayeeshipTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistryMock.contract.FilterLogs(opts, "PayeeshipTransferred", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryMockPayeeshipTransferredIterator{contract: _KeeperRegistryMock.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockPayeeshipTransferred, keeper []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryMock.contract.WatchLogs(opts, "PayeeshipTransferred", keeperRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryMockPayeeshipTransferred)
				if err := _KeeperRegistryMock.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) ParsePayeeshipTransferred(log types.Log) (*KeeperRegistryMockPayeeshipTransferred, error) {
	event := new(KeeperRegistryMockPayeeshipTransferred)
	if err := _KeeperRegistryMock.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryMockPaymentWithdrawnIterator struct {
	Event *KeeperRegistryMockPaymentWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryMockPaymentWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryMockPaymentWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryMockPaymentWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryMockPaymentWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryMockPaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryMockPaymentWithdrawn struct {
	Keeper common.Address
	Amount *big.Int
	To     common.Address
	Payee  common.Address
	Raw    types.Log
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, keeper []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistryMockPaymentWithdrawnIterator, error) {

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

	logs, sub, err := _KeeperRegistryMock.contract.FilterLogs(opts, "PaymentWithdrawn", keeperRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryMockPaymentWithdrawnIterator{contract: _KeeperRegistryMock.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockPaymentWithdrawn, keeper []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryMock.contract.WatchLogs(opts, "PaymentWithdrawn", keeperRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryMockPaymentWithdrawn)
				if err := _KeeperRegistryMock.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) ParsePaymentWithdrawn(log types.Log) (*KeeperRegistryMockPaymentWithdrawn, error) {
	event := new(KeeperRegistryMockPaymentWithdrawn)
	if err := _KeeperRegistryMock.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryMockRegistrarChangedIterator struct {
	Event *KeeperRegistryMockRegistrarChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryMockRegistrarChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryMockRegistrarChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryMockRegistrarChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryMockRegistrarChangedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryMockRegistrarChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryMockRegistrarChanged struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) FilterRegistrarChanged(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryMockRegistrarChangedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryMock.contract.FilterLogs(opts, "RegistrarChanged", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryMockRegistrarChangedIterator{contract: _KeeperRegistryMock.contract, event: "RegistrarChanged", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) WatchRegistrarChanged(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockRegistrarChanged, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryMock.contract.WatchLogs(opts, "RegistrarChanged", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryMockRegistrarChanged)
				if err := _KeeperRegistryMock.contract.UnpackLog(event, "RegistrarChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) ParseRegistrarChanged(log types.Log) (*KeeperRegistryMockRegistrarChanged, error) {
	event := new(KeeperRegistryMockRegistrarChanged)
	if err := _KeeperRegistryMock.contract.UnpackLog(event, "RegistrarChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryMockUnpausedIterator struct {
	Event *KeeperRegistryMockUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryMockUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryMockUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryMockUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryMockUnpausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryMockUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryMockUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistryMockUnpausedIterator, error) {

	logs, sub, err := _KeeperRegistryMock.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryMockUnpausedIterator{contract: _KeeperRegistryMock.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockUnpaused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryMock.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryMockUnpaused)
				if err := _KeeperRegistryMock.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) ParseUnpaused(log types.Log) (*KeeperRegistryMockUnpaused, error) {
	event := new(KeeperRegistryMockUnpaused)
	if err := _KeeperRegistryMock.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryMockUpkeepCanceledIterator struct {
	Event *KeeperRegistryMockUpkeepCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryMockUpkeepCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryMockUpkeepCanceled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryMockUpkeepCanceled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryMockUpkeepCanceledIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryMockUpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryMockUpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistryMockUpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistryMock.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryMockUpkeepCanceledIterator{contract: _KeeperRegistryMock.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistryMock.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryMockUpkeepCanceled)
				if err := _KeeperRegistryMock.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) ParseUpkeepCanceled(log types.Log) (*KeeperRegistryMockUpkeepCanceled, error) {
	event := new(KeeperRegistryMockUpkeepCanceled)
	if err := _KeeperRegistryMock.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryMockUpkeepPerformedIterator struct {
	Event *KeeperRegistryMockUpkeepPerformed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryMockUpkeepPerformedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryMockUpkeepPerformed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryMockUpkeepPerformed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryMockUpkeepPerformedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryMockUpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryMockUpkeepPerformed struct {
	Id          *big.Int
	Success     bool
	From        common.Address
	Payment     *big.Int
	PerformData []byte
	Raw         types.Log
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool, from []common.Address) (*KeeperRegistryMockUpkeepPerformedIterator, error) {

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

	logs, sub, err := _KeeperRegistryMock.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryMockUpkeepPerformedIterator{contract: _KeeperRegistryMock.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockUpkeepPerformed, id []*big.Int, success []bool, from []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryMock.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryMockUpkeepPerformed)
				if err := _KeeperRegistryMock.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) ParseUpkeepPerformed(log types.Log) (*KeeperRegistryMockUpkeepPerformed, error) {
	event := new(KeeperRegistryMockUpkeepPerformed)
	if err := _KeeperRegistryMock.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryMockUpkeepRegisteredIterator struct {
	Event *KeeperRegistryMockUpkeepRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryMockUpkeepRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryMockUpkeepRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistryMockUpkeepRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryMockUpkeepRegisteredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryMockUpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryMockUpkeepRegistered struct {
	Id         *big.Int
	ExecuteGas uint32
	Admin      common.Address
	Raw        types.Log
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryMockUpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryMock.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryMockUpkeepRegisteredIterator{contract: _KeeperRegistryMock.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockUpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryMock.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryMockUpkeepRegistered)
				if err := _KeeperRegistryMock.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryMock *KeeperRegistryMockFilterer) ParseUpkeepRegistered(log types.Log) (*KeeperRegistryMockUpkeepRegistered, error) {
	event := new(KeeperRegistryMockUpkeepRegistered)
	if err := _KeeperRegistryMock.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CheckUpkeep struct {
	PerformData    []byte
	MaxLinkPayment *big.Int
	GasLimit       *big.Int
	AdjustedGasWei *big.Int
	LinkEth        *big.Int
}
type GetConfig struct {
	PaymentPremiumPPB    uint32
	BlockCountPerTurn    *big.Int
	CheckGasLimit        uint32
	StalenessSeconds     *big.Int
	GasCeilingMultiplier uint16
	FallbackGasPrice     *big.Int
	FallbackLinkPrice    *big.Int
}
type GetUpkeep struct {
	Target              common.Address
	ExecuteGas          uint32
	CheckData           []byte
	Balance             *big.Int
	LastKeeper          common.Address
	Admin               common.Address
	MaxValidBlocknumber uint64
}

func (_KeeperRegistryMock *KeeperRegistryMock) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _KeeperRegistryMock.abi.Events["ConfigSet"].ID:
		return _KeeperRegistryMock.ParseConfigSet(log)
	case _KeeperRegistryMock.abi.Events["FlatFeeSet"].ID:
		return _KeeperRegistryMock.ParseFlatFeeSet(log)
	case _KeeperRegistryMock.abi.Events["FundsAdded"].ID:
		return _KeeperRegistryMock.ParseFundsAdded(log)
	case _KeeperRegistryMock.abi.Events["FundsWithdrawn"].ID:
		return _KeeperRegistryMock.ParseFundsWithdrawn(log)
	case _KeeperRegistryMock.abi.Events["KeepersUpdated"].ID:
		return _KeeperRegistryMock.ParseKeepersUpdated(log)
	case _KeeperRegistryMock.abi.Events["OwnershipTransferRequested"].ID:
		return _KeeperRegistryMock.ParseOwnershipTransferRequested(log)
	case _KeeperRegistryMock.abi.Events["OwnershipTransferred"].ID:
		return _KeeperRegistryMock.ParseOwnershipTransferred(log)
	case _KeeperRegistryMock.abi.Events["Paused"].ID:
		return _KeeperRegistryMock.ParsePaused(log)
	case _KeeperRegistryMock.abi.Events["PayeeshipTransferRequested"].ID:
		return _KeeperRegistryMock.ParsePayeeshipTransferRequested(log)
	case _KeeperRegistryMock.abi.Events["PayeeshipTransferred"].ID:
		return _KeeperRegistryMock.ParsePayeeshipTransferred(log)
	case _KeeperRegistryMock.abi.Events["PaymentWithdrawn"].ID:
		return _KeeperRegistryMock.ParsePaymentWithdrawn(log)
	case _KeeperRegistryMock.abi.Events["RegistrarChanged"].ID:
		return _KeeperRegistryMock.ParseRegistrarChanged(log)
	case _KeeperRegistryMock.abi.Events["Unpaused"].ID:
		return _KeeperRegistryMock.ParseUnpaused(log)
	case _KeeperRegistryMock.abi.Events["UpkeepCanceled"].ID:
		return _KeeperRegistryMock.ParseUpkeepCanceled(log)
	case _KeeperRegistryMock.abi.Events["UpkeepPerformed"].ID:
		return _KeeperRegistryMock.ParseUpkeepPerformed(log)
	case _KeeperRegistryMock.abi.Events["UpkeepRegistered"].ID:
		return _KeeperRegistryMock.ParseUpkeepRegistered(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeeperRegistryMockConfigSet) Topic() common.Hash {
	return common.HexToHash("0xeb3c06937e6595fd80ec1add18a195026d5cf65f122cc3ffedbfb18a9ed80b39")
}

func (KeeperRegistryMockFlatFeeSet) Topic() common.Hash {
	return common.HexToHash("0x17b46a44a823646eef686b7824df2962de896bc9a012a60b67694c5cbf184d8b")
}

func (KeeperRegistryMockFundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (KeeperRegistryMockFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (KeeperRegistryMockKeepersUpdated) Topic() common.Hash {
	return common.HexToHash("0x056264c94f28bb06c99d13f0446eb96c67c215d8d707bce2655a98ddf1c0b71f")
}

func (KeeperRegistryMockOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (KeeperRegistryMockOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (KeeperRegistryMockPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (KeeperRegistryMockPayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (KeeperRegistryMockPayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (KeeperRegistryMockPaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (KeeperRegistryMockRegistrarChanged) Topic() common.Hash {
	return common.HexToHash("0x9bf4a5b30267728df68663e14adb47e559863967c419dc6030638883408bed2e")
}

func (KeeperRegistryMockUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (KeeperRegistryMockUpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (KeeperRegistryMockUpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0xcaacad83e47cc45c280d487ec84184eee2fa3b54ebaa393bda7549f13da228f6")
}

func (KeeperRegistryMockUpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (_KeeperRegistryMock *KeeperRegistryMock) Address() common.Address {
	return _KeeperRegistryMock.address
}

type KeeperRegistryMockInterface interface {
	CheckUpkeep(opts *bind.CallOpts, id *big.Int, from common.Address) (CheckUpkeep,

		error)

	GetCanceledUpkeepList(opts *bind.CallOpts) ([]*big.Int, error)

	GetConfig(opts *bind.CallOpts) (GetConfig,

		error)

	GetKeeperList(opts *bind.CallOpts) ([]common.Address, error)

	GetMinBalanceForUpkeep(opts *bind.CallOpts, id *big.Int) (*big.Int, error)

	GetUpkeep(opts *bind.CallOpts, id *big.Int) (GetUpkeep,

		error)

	GetUpkeepCount(opts *bind.CallOpts) (*big.Int, error)

	EmitConfigSet(opts *bind.TransactOpts, paymentPremiumPPB uint32, blockCountPerTurn *big.Int, checkGasLimit uint32, stalenessSeconds *big.Int, gasCeilingMultiplier uint16, fallbackGasPrice *big.Int, fallbackLinkPrice *big.Int) (*types.Transaction, error)

	EmitFlatFeeSet(opts *bind.TransactOpts, flatFeeMicroLink uint32) (*types.Transaction, error)

	EmitFundsAdded(opts *bind.TransactOpts, id *big.Int, from common.Address, amount *big.Int) (*types.Transaction, error)

	EmitFundsWithdrawn(opts *bind.TransactOpts, id *big.Int, amount *big.Int, to common.Address) (*types.Transaction, error)

	EmitKeepersUpdated(opts *bind.TransactOpts, keepers []common.Address, payees []common.Address) (*types.Transaction, error)

	EmitOwnershipTransferRequested(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	EmitOwnershipTransferred(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	EmitPaused(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error)

	EmitPayeeshipTransferRequested(opts *bind.TransactOpts, keeper common.Address, from common.Address, to common.Address) (*types.Transaction, error)

	EmitPayeeshipTransferred(opts *bind.TransactOpts, keeper common.Address, from common.Address, to common.Address) (*types.Transaction, error)

	EmitPaymentWithdrawn(opts *bind.TransactOpts, keeper common.Address, amount *big.Int, to common.Address, payee common.Address) (*types.Transaction, error)

	EmitRegistrarChanged(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	EmitUnpaused(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error)

	EmitUpkeepCanceled(opts *bind.TransactOpts, id *big.Int, atBlockHeight uint64) (*types.Transaction, error)

	EmitUpkeepPerformed(opts *bind.TransactOpts, id *big.Int, success bool, from common.Address, payment *big.Int, performData []byte) (*types.Transaction, error)

	EmitUpkeepRegistered(opts *bind.TransactOpts, id *big.Int, executeGas uint32, admin common.Address) (*types.Transaction, error)

	PerformUpkeep(opts *bind.TransactOpts, id *big.Int, performData []byte) (*types.Transaction, error)

	SetCanceledUpkeepList(opts *bind.TransactOpts, _canceledUpkeepList []*big.Int) (*types.Transaction, error)

	SetCheckUpkeepData(opts *bind.TransactOpts, id *big.Int, performData []byte, maxLinkPayment *big.Int, gasLimit *big.Int, adjustedGasWei *big.Int, linkEth *big.Int) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, _paymentPremiumPPB uint32, _flatFeeMicroLink uint32, _blockCountPerTurn *big.Int, _checkGasLimit uint32, _stalenessSeconds *big.Int, _gasCeilingMultiplier uint16, _fallbackGasPrice *big.Int, _fallbackLinkPrice *big.Int) (*types.Transaction, error)

	SetKeeperList(opts *bind.TransactOpts, _keepers []common.Address) (*types.Transaction, error)

	SetMinBalance(opts *bind.TransactOpts, id *big.Int, minBalance *big.Int) (*types.Transaction, error)

	SetPerformUpkeepSuccess(opts *bind.TransactOpts, id *big.Int, success bool) (*types.Transaction, error)

	SetUpkeep(opts *bind.TransactOpts, id *big.Int, _target common.Address, _executeGas uint32, _balance *big.Int, _admin common.Address, _maxValidBlocknumber uint64, _lastKeeper common.Address, _checkData []byte) (*types.Transaction, error)

	SetUpkeepCount(opts *bind.TransactOpts, _upkeepCount *big.Int) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*KeeperRegistryMockConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*KeeperRegistryMockConfigSet, error)

	FilterFlatFeeSet(opts *bind.FilterOpts) (*KeeperRegistryMockFlatFeeSetIterator, error)

	WatchFlatFeeSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockFlatFeeSet) (event.Subscription, error)

	ParseFlatFeeSet(log types.Log) (*KeeperRegistryMockFlatFeeSet, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistryMockFundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*KeeperRegistryMockFundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryMockFundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockFundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*KeeperRegistryMockFundsWithdrawn, error)

	FilterKeepersUpdated(opts *bind.FilterOpts) (*KeeperRegistryMockKeepersUpdatedIterator, error)

	WatchKeepersUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockKeepersUpdated) (event.Subscription, error)

	ParseKeepersUpdated(log types.Log) (*KeeperRegistryMockKeepersUpdated, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryMockOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryMockOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryMockOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*KeeperRegistryMockOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*KeeperRegistryMockPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*KeeperRegistryMockPaused, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, keeper []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryMockPayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockPayeeshipTransferRequested, keeper []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistryMockPayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, keeper []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryMockPayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockPayeeshipTransferred, keeper []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*KeeperRegistryMockPayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, keeper []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistryMockPaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockPaymentWithdrawn, keeper []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*KeeperRegistryMockPaymentWithdrawn, error)

	FilterRegistrarChanged(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryMockRegistrarChangedIterator, error)

	WatchRegistrarChanged(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockRegistrarChanged, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseRegistrarChanged(log types.Log) (*KeeperRegistryMockRegistrarChanged, error)

	FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistryMockUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*KeeperRegistryMockUnpaused, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistryMockUpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockUpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*KeeperRegistryMockUpkeepCanceled, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool, from []common.Address) (*KeeperRegistryMockUpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockUpkeepPerformed, id []*big.Int, success []bool, from []common.Address) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*KeeperRegistryMockUpkeepPerformed, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryMockUpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistryMockUpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*KeeperRegistryMockUpkeepRegistered, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
