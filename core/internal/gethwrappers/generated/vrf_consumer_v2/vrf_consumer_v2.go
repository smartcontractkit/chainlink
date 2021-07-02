// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_consumer_v2

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

const VRFConsumerV2ABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_gasAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_subId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"testCreateSubscriptionAndFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"minReqConfs\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"callbackGasLimit\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"numWords\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"consumerID\",\"type\":\"uint16\"}],\"name\":\"testRequestRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"name\":\"updateSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

var VRFConsumerV2Bin = "0x608060405234801561001057600080fd5b50604051610b70380380610b7083398101604081905261002f9161007c565b600280546001600160a01b039384166001600160a01b031991821617909155600380549290931691161790556100af565b80516001600160a01b038116811461007757600080fd5b919050565b6000806040838503121561008f57600080fd5b61009883610060565b91506100a660208401610060565b90509250929050565b610ab2806100be6000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c8063706da1ca1161005b578063706da1ca146100ee578063e89e106a14610133578063f08c5daa1461013c578063f6eaffc81461014557600080fd5b806336bfffed1461008d57806338ba4614146100a25780635d8013a9146100b55780636802f726146100db575b600080fd5b6100a061009b3660046106a3565b610158565b005b6100a06100b036600461083d565b610299565b6100c86100c336600461078b565b6102b7565b6040519081526020015b60405180910390f35b6100a06100e93660046108fe565b610399565b60035461011a9074010000000000000000000000000000000000000000900467ffffffffffffffff1681565b60405167ffffffffffffffff90911681526020016100d2565b6100c860015481565b6100c860045481565b6100c861015336600461080b565b610622565b60035474010000000000000000000000000000000000000000900467ffffffffffffffff166101e7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f7375624944206e6f742073657400000000000000000000000000000000000000604482015260640160405180910390fd5b6002546003546040517f15ae7e3800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff909216916315ae7e3891610264917401000000000000000000000000000000000000000090910467ffffffffffffffff16908590600401610990565b600060405180830381600087803b15801561027e57600080fd5b505af1158015610292573d6000803e3d6000fd5b5050505050565b5a60045580516102b0906000906020840190610643565b5050600155565b6002546040517f74fe4c550000000000000000000000000000000000000000000000000000000081526004810188905267ffffffffffffffff8088166024830152808716604483015280861660648301528416608482015261ffff831660a482015260009173ffffffffffffffffffffffffffffffffffffffff16906374fe4c559060c401602060405180830381600087803b15801561035657600080fd5b505af115801561036a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061038e9190610824565b979650505050505050565b60035474010000000000000000000000000000000000000000900467ffffffffffffffff166104d7576040805160018082528183019092526000916020808301908036833701905050905030816000815181106103f8576103f8610a2e565b73ffffffffffffffffffffffffffffffffffffffff92831660209182029290920101526002546040517f6b9f7d38000000000000000000000000000000000000000000000000000000008152911690636b9f7d389061045b90849060040161097d565b602060405180830381600087803b15801561047557600080fd5b505af1158015610489573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104ad91906108e1565b600360146101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550505b6003546002546040517f095ea7b300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201526bffffffffffffffffffffffff8416602482015291169063095ea7b390604401602060405180830381600087803b15801561055a57600080fd5b505af115801561056e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105929190610762565b506002546003546040517fafc69b530000000000000000000000000000000000000000000000000000000081527401000000000000000000000000000000000000000090910467ffffffffffffffff1660048201526bffffffffffffffffffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063afc69b5390604401610264565b6000818154811061063257600080fd5b600091825260209091200154905081565b82805482825590600052602060002090810192821561067e579160200282015b8281111561067e578251825591602001919060010190610663565b5061068a92915061068e565b5090565b5b8082111561068a576000815560010161068f565b600060208083850312156106b657600080fd5b823567ffffffffffffffff8111156106cd57600080fd5b8301601f810185136106de57600080fd5b80356106f16106ec82610a0a565b6109bb565b80828252848201915084840188868560051b870101111561071157600080fd5b60009450845b8481101561075457813573ffffffffffffffffffffffffffffffffffffffff81168114610742578687fd5b84529286019290860190600101610717565b509098975050505050505050565b60006020828403121561077457600080fd5b8151801515811461078457600080fd5b9392505050565b60008060008060008060c087890312156107a457600080fd5b8635955060208701356107b681610a8c565b945060408701356107c681610a8c565b935060608701356107d681610a8c565b925060808701356107e681610a8c565b915060a087013561ffff811681146107fd57600080fd5b809150509295509295509295565b60006020828403121561081d57600080fd5b5035919050565b60006020828403121561083657600080fd5b5051919050565b6000806040838503121561085057600080fd5b8235915060208084013567ffffffffffffffff81111561086f57600080fd5b8401601f8101861361088057600080fd5b803561088e6106ec82610a0a565b80828252848201915084840189868560051b87010111156108ae57600080fd5b600094505b838510156108d15780358352600194909401939185019185016108b3565b5080955050505050509250929050565b6000602082840312156108f357600080fd5b815161078481610a8c565b60006020828403121561091057600080fd5b81356bffffffffffffffffffffffff8116811461078457600080fd5b600081518084526020808501945080840160005b8381101561097257815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101610940565b509495945050505050565b602081526000610784602083018461092c565b67ffffffffffffffff831681526040602082015260006109b3604083018461092c565b949350505050565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610a0257610a02610a5d565b604052919050565b600067ffffffffffffffff821115610a2457610a24610a5d565b5060051b60200190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b67ffffffffffffffff81168114610aa257600080fd5b5056fea164736f6c6343000806000a"

func DeployVRFConsumerV2(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, link common.Address) (common.Address, *types.Transaction, *VRFConsumerV2, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFConsumerV2ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(VRFConsumerV2Bin), backend, vrfCoordinator, link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFConsumerV2{VRFConsumerV2Caller: VRFConsumerV2Caller{contract: contract}, VRFConsumerV2Transactor: VRFConsumerV2Transactor{contract: contract}, VRFConsumerV2Filterer: VRFConsumerV2Filterer{contract: contract}}, nil
}

type VRFConsumerV2 struct {
	address common.Address
	abi     abi.ABI
	VRFConsumerV2Caller
	VRFConsumerV2Transactor
	VRFConsumerV2Filterer
}

type VRFConsumerV2Caller struct {
	contract *bind.BoundContract
}

type VRFConsumerV2Transactor struct {
	contract *bind.BoundContract
}

type VRFConsumerV2Filterer struct {
	contract *bind.BoundContract
}

type VRFConsumerV2Session struct {
	Contract     *VRFConsumerV2
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFConsumerV2CallerSession struct {
	Contract *VRFConsumerV2Caller
	CallOpts bind.CallOpts
}

type VRFConsumerV2TransactorSession struct {
	Contract     *VRFConsumerV2Transactor
	TransactOpts bind.TransactOpts
}

type VRFConsumerV2Raw struct {
	Contract *VRFConsumerV2
}

type VRFConsumerV2CallerRaw struct {
	Contract *VRFConsumerV2Caller
}

type VRFConsumerV2TransactorRaw struct {
	Contract *VRFConsumerV2Transactor
}

func NewVRFConsumerV2(address common.Address, backend bind.ContractBackend) (*VRFConsumerV2, error) {
	abi, err := abi.JSON(strings.NewReader(VRFConsumerV2ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFConsumerV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2{address: address, abi: abi, VRFConsumerV2Caller: VRFConsumerV2Caller{contract: contract}, VRFConsumerV2Transactor: VRFConsumerV2Transactor{contract: contract}, VRFConsumerV2Filterer: VRFConsumerV2Filterer{contract: contract}}, nil
}

func NewVRFConsumerV2Caller(address common.Address, caller bind.ContractCaller) (*VRFConsumerV2Caller, error) {
	contract, err := bindVRFConsumerV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2Caller{contract: contract}, nil
}

func NewVRFConsumerV2Transactor(address common.Address, transactor bind.ContractTransactor) (*VRFConsumerV2Transactor, error) {
	contract, err := bindVRFConsumerV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2Transactor{contract: contract}, nil
}

func NewVRFConsumerV2Filterer(address common.Address, filterer bind.ContractFilterer) (*VRFConsumerV2Filterer, error) {
	contract, err := bindVRFConsumerV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2Filterer{contract: contract}, nil
}

func bindVRFConsumerV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFConsumerV2ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_VRFConsumerV2 *VRFConsumerV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFConsumerV2.Contract.VRFConsumerV2Caller.contract.Call(opts, result, method, params...)
}

func (_VRFConsumerV2 *VRFConsumerV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.VRFConsumerV2Transactor.contract.Transfer(opts)
}

func (_VRFConsumerV2 *VRFConsumerV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.VRFConsumerV2Transactor.contract.Transact(opts, method, params...)
}

func (_VRFConsumerV2 *VRFConsumerV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFConsumerV2.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFConsumerV2 *VRFConsumerV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.contract.Transfer(opts)
}

func (_VRFConsumerV2 *VRFConsumerV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.contract.Transact(opts, method, params...)
}

func (_VRFConsumerV2 *VRFConsumerV2Caller) SGasAvailable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFConsumerV2.contract.Call(opts, &out, "s_gasAvailable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFConsumerV2 *VRFConsumerV2Session) SGasAvailable() (*big.Int, error) {
	return _VRFConsumerV2.Contract.SGasAvailable(&_VRFConsumerV2.CallOpts)
}

func (_VRFConsumerV2 *VRFConsumerV2CallerSession) SGasAvailable() (*big.Int, error) {
	return _VRFConsumerV2.Contract.SGasAvailable(&_VRFConsumerV2.CallOpts)
}

func (_VRFConsumerV2 *VRFConsumerV2Caller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFConsumerV2.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFConsumerV2 *VRFConsumerV2Session) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFConsumerV2.Contract.SRandomWords(&_VRFConsumerV2.CallOpts, arg0)
}

func (_VRFConsumerV2 *VRFConsumerV2CallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFConsumerV2.Contract.SRandomWords(&_VRFConsumerV2.CallOpts, arg0)
}

func (_VRFConsumerV2 *VRFConsumerV2Caller) SRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFConsumerV2.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFConsumerV2 *VRFConsumerV2Session) SRequestId() (*big.Int, error) {
	return _VRFConsumerV2.Contract.SRequestId(&_VRFConsumerV2.CallOpts)
}

func (_VRFConsumerV2 *VRFConsumerV2CallerSession) SRequestId() (*big.Int, error) {
	return _VRFConsumerV2.Contract.SRequestId(&_VRFConsumerV2.CallOpts)
}

func (_VRFConsumerV2 *VRFConsumerV2Caller) SSubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFConsumerV2.contract.Call(opts, &out, "s_subId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFConsumerV2 *VRFConsumerV2Session) SSubId() (uint64, error) {
	return _VRFConsumerV2.Contract.SSubId(&_VRFConsumerV2.CallOpts)
}

func (_VRFConsumerV2 *VRFConsumerV2CallerSession) SSubId() (uint64, error) {
	return _VRFConsumerV2.Contract.SSubId(&_VRFConsumerV2.CallOpts)
}

func (_VRFConsumerV2 *VRFConsumerV2Transactor) FulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.contract.Transact(opts, "fulfillRandomWords", requestId, randomWords)
}

func (_VRFConsumerV2 *VRFConsumerV2Session) FulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.FulfillRandomWords(&_VRFConsumerV2.TransactOpts, requestId, randomWords)
}

func (_VRFConsumerV2 *VRFConsumerV2TransactorSession) FulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.FulfillRandomWords(&_VRFConsumerV2.TransactOpts, requestId, randomWords)
}

func (_VRFConsumerV2 *VRFConsumerV2Transactor) TestCreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.contract.Transact(opts, "testCreateSubscriptionAndFund", amount)
}

func (_VRFConsumerV2 *VRFConsumerV2Session) TestCreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.TestCreateSubscriptionAndFund(&_VRFConsumerV2.TransactOpts, amount)
}

func (_VRFConsumerV2 *VRFConsumerV2TransactorSession) TestCreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.TestCreateSubscriptionAndFund(&_VRFConsumerV2.TransactOpts, amount)
}

func (_VRFConsumerV2 *VRFConsumerV2Transactor) TestRequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, minReqConfs uint64, callbackGasLimit uint64, numWords uint64, consumerID uint16) (*types.Transaction, error) {
	return _VRFConsumerV2.contract.Transact(opts, "testRequestRandomness", keyHash, subId, minReqConfs, callbackGasLimit, numWords, consumerID)
}

func (_VRFConsumerV2 *VRFConsumerV2Session) TestRequestRandomness(keyHash [32]byte, subId uint64, minReqConfs uint64, callbackGasLimit uint64, numWords uint64, consumerID uint16) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.TestRequestRandomness(&_VRFConsumerV2.TransactOpts, keyHash, subId, minReqConfs, callbackGasLimit, numWords, consumerID)
}

func (_VRFConsumerV2 *VRFConsumerV2TransactorSession) TestRequestRandomness(keyHash [32]byte, subId uint64, minReqConfs uint64, callbackGasLimit uint64, numWords uint64, consumerID uint16) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.TestRequestRandomness(&_VRFConsumerV2.TransactOpts, keyHash, subId, minReqConfs, callbackGasLimit, numWords, consumerID)
}

func (_VRFConsumerV2 *VRFConsumerV2Transactor) UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2.contract.Transact(opts, "updateSubscription", consumers)
}

func (_VRFConsumerV2 *VRFConsumerV2Session) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.UpdateSubscription(&_VRFConsumerV2.TransactOpts, consumers)
}

func (_VRFConsumerV2 *VRFConsumerV2TransactorSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.UpdateSubscription(&_VRFConsumerV2.TransactOpts, consumers)
}

func (_VRFConsumerV2 *VRFConsumerV2) Address() common.Address {
	return _VRFConsumerV2.address
}

type VRFConsumerV2Interface interface {
	SGasAvailable(opts *bind.CallOpts) (*big.Int, error)

	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestId(opts *bind.CallOpts) (*big.Int, error)

	SSubId(opts *bind.CallOpts) (uint64, error)

	FulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	TestCreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	TestRequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, minReqConfs uint64, callbackGasLimit uint64, numWords uint64, consumerID uint16) (*types.Transaction, error)

	UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error)

	Address() common.Address
}
