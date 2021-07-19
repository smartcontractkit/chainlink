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

const VRFConsumerV2ABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_gasAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_subId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"testCreateSubscriptionAndFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"minReqConfs\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"consumerID\",\"type\":\"uint32\"}],\"name\":\"testRequestRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"name\":\"updateSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

var VRFConsumerV2Bin = "0x608060405234801561001057600080fd5b50604051610bb8380380610bb883398101604081905261002f9161007c565b600280546001600160a01b039384166001600160a01b031991821617909155600380549290931691161790556100af565b80516001600160a01b038116811461007757600080fd5b919050565b6000806040838503121561008f57600080fd5b61009883610060565b91506100a660208401610060565b90509250929050565b610afa806100be6000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c80638324698a1161005b5780638324698a14610112578063e89e106a14610133578063f08c5daa1461013c578063f6eaffc81461014557600080fd5b806336bfffed1461008d57806338ba4614146100a25780636802f726146100b5578063706da1ca146100c8575b600080fd5b6100a061009b36600461064d565b610158565b005b6100a06100b03660046107df565b610299565b6100a06100c33660046108a0565b6102b7565b6003546100f49074010000000000000000000000000000000000000000900467ffffffffffffffff1681565b60405167ffffffffffffffff90911681526020015b60405180910390f35b610125610120366004610735565b6104cd565b604051908152602001610109565b61012560015481565b61012560045481565b6101256101533660046107ad565b6105b3565b60035474010000000000000000000000000000000000000000900467ffffffffffffffff166101e7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f7375624944206e6f742073657400000000000000000000000000000000000000604482015260640160405180910390fd5b6002546003546040517f15ae7e3800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff909216916315ae7e3891610264917401000000000000000000000000000000000000000090910467ffffffffffffffff169085906004016109d8565b600060405180830381600087803b15801561027e57600080fd5b505af1158015610292573d6000803e3d6000fd5b5050505050565b5a60045580516102b09060009060208401906105d4565b5050600155565b60035474010000000000000000000000000000000000000000900467ffffffffffffffff166103f55760408051600180825281830190925260009160208083019080368337019050509050308160008151811061031657610316610a76565b73ffffffffffffffffffffffffffffffffffffffff92831660209182029290920101526002546040517f6b9f7d38000000000000000000000000000000000000000000000000000000008152911690636b9f7d38906103799084906004016109c5565b602060405180830381600087803b15801561039357600080fd5b505af11580156103a7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103cb9190610883565b600360146101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550505b6003546002546040805174010000000000000000000000000000000000000000840467ffffffffffffffff16602082015273ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918591016040516020818303038152906040526040518463ffffffff1660e01b81526004016104779392919061091f565b602060405180830381600087803b15801561049157600080fd5b505af11580156104a5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104c9919061070c565b5050565b6002546040517f243460600000000000000000000000000000000000000000000000000000000081526004810188905267ffffffffffffffff8716602482015261ffff8616604482015263ffffffff80861660648301528085166084830152831660a482015260009173ffffffffffffffffffffffffffffffffffffffff169063243460609060c401602060405180830381600087803b15801561057057600080fd5b505af1158015610584573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105a891906107c6565b979650505050505050565b600081815481106105c357600080fd5b600091825260209091200154905081565b82805482825590600052602060002090810192821561060f579160200282015b8281111561060f5782518255916020019190600101906105f4565b5061061b92915061061f565b5090565b5b8082111561061b5760008155600101610620565b803563ffffffff8116811461064857600080fd5b919050565b6000602080838503121561066057600080fd5b823567ffffffffffffffff81111561067757600080fd5b8301601f8101851361068857600080fd5b803561069b61069682610a52565b610a03565b80828252848201915084840188868560051b87010111156106bb57600080fd5b60009450845b848110156106fe57813573ffffffffffffffffffffffffffffffffffffffff811681146106ec578687fd5b845292860192908601906001016106c1565b509098975050505050505050565b60006020828403121561071e57600080fd5b8151801515811461072e57600080fd5b9392505050565b60008060008060008060c0878903121561074e57600080fd5b86359550602087013561076081610ad4565b9450604087013561ffff8116811461077757600080fd5b935061078560608801610634565b925061079360808801610634565b91506107a160a08801610634565b90509295509295509295565b6000602082840312156107bf57600080fd5b5035919050565b6000602082840312156107d857600080fd5b5051919050565b600080604083850312156107f257600080fd5b8235915060208084013567ffffffffffffffff81111561081157600080fd5b8401601f8101861361082257600080fd5b803561083061069682610a52565b80828252848201915084840189868560051b870101111561085057600080fd5b600094505b83851015610873578035835260019490940193918501918501610855565b5080955050505050509250929050565b60006020828403121561089557600080fd5b815161072e81610ad4565b6000602082840312156108b257600080fd5b81356bffffffffffffffffffffffff8116811461072e57600080fd5b600081518084526020808501945080840160005b8381101561091457815173ffffffffffffffffffffffffffffffffffffffff16875295820195908201906001016108e2565b509495945050505050565b73ffffffffffffffffffffffffffffffffffffffff84168152600060206bffffffffffffffffffffffff85168184015260606040840152835180606085015260005b8181101561097d57858101830151858201608001528201610961565b8181111561098f576000608083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160800195945050505050565b60208152600061072e60208301846108ce565b67ffffffffffffffff831681526040602082015260006109fb60408301846108ce565b949350505050565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610a4a57610a4a610aa5565b604052919050565b600067ffffffffffffffff821115610a6c57610a6c610aa5565b5060051b60200190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b67ffffffffffffffff81168114610aea57600080fd5b5056fea164736f6c6343000806000a"

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

func (_VRFConsumerV2 *VRFConsumerV2Transactor) TestRequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32, consumerID uint32) (*types.Transaction, error) {
	return _VRFConsumerV2.contract.Transact(opts, "testRequestRandomness", keyHash, subId, minReqConfs, callbackGasLimit, numWords, consumerID)
}

func (_VRFConsumerV2 *VRFConsumerV2Session) TestRequestRandomness(keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32, consumerID uint32) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.TestRequestRandomness(&_VRFConsumerV2.TransactOpts, keyHash, subId, minReqConfs, callbackGasLimit, numWords, consumerID)
}

func (_VRFConsumerV2 *VRFConsumerV2TransactorSession) TestRequestRandomness(keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32, consumerID uint32) (*types.Transaction, error) {
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

	TestRequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32, consumerID uint32) (*types.Transaction, error)

	UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error)

	Address() common.Address
}
