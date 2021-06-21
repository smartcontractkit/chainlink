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

const VRFConsumerV2ABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"_randomWords\",\"type\":\"uint256[]\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"subId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"testCreateSubscriptionAndFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"_subId\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"minReqConfs\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"callbackGasLimit\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"numWords\",\"type\":\"uint64\"}],\"name\":\"testRequestRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

var VRFConsumerV2Bin = "0x608060405234801561001057600080fd5b5060405161088d38038061088d83398101604081905261002f9161007c565b600280546001600160a01b039384166001600160a01b031991821617909155600380549290931691161790556100ae565b80516001600160a01b038116811461007757600080fd5b919050565b6000806040838503121561008e578182fd5b61009783610060565b91506100a560208401610060565b90509250929050565b6107d0806100bd6000396000f3fe608060405234801561001057600080fd5b50600436106100715760003560e01c8063beff730f11610050578063beff730f146100ba578063eb1d28bb146100cd578063ec6074491461011257600080fd5b80626d6cae1461007657806338ba4614146100925780638ed8d353146100a7575b600080fd5b61007f60015481565b6040519081526020015b60405180910390f35b6100a56100a036600461061b565b610125565b005b61007f6100b5366004610584565b61013f565b61007f6100c83660046105eb565b610216565b6003546100f99074010000000000000000000000000000000000000000900467ffffffffffffffff1681565b60405167ffffffffffffffff9091168152602001610089565b6100a56101203660046105eb565b610237565b80516101389060009060208401906104fd565b5050600155565b6002546040517f9cb1298b0000000000000000000000000000000000000000000000000000000081526004810187905267ffffffffffffffff8087166024830152808616604483015280851660648301528316608482015260009173ffffffffffffffffffffffffffffffffffffffff1690639cb1298b9060a401602060405180830381600087803b1580156101d457600080fd5b505af11580156101e8573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061020c9190610603565b9695505050505050565b6000818154811061022657600080fd5b600091825260209091200154905081565b60035474010000000000000000000000000000000000000000900467ffffffffffffffff1661039c576040805160018082528183019092526000916020808301908036833701905050905030816000815181106102bd577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff92831660209182029290920101526002546040517f6b9f7d38000000000000000000000000000000000000000000000000000000008152911690636b9f7d3890610320908490600401610721565b602060405180830381600087803b15801561033a57600080fd5b505af115801561034e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103729190610705565b600360146101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550505b6003546002546040517f095ea7b300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201526024810184905291169063095ea7b390604401602060405180830381600087803b15801561041257600080fd5b505af1158015610426573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061044a919061055d565b506002546003546040517fdd2819280000000000000000000000000000000000000000000000000000000081527401000000000000000000000000000000000000000090910467ffffffffffffffff1660048201526024810183905273ffffffffffffffffffffffffffffffffffffffff9091169063dd28192890604401600060405180830381600087803b1580156104e257600080fd5b505af11580156104f6573d6000803e3d6000fd5b5050505050565b828054828255906000526020600020908101928215610538579160200282015b8281111561053857825182559160200191906001019061051d565b50610544929150610548565b5090565b5b808211156105445760008155600101610549565b60006020828403121561056e578081fd5b8151801515811461057d578182fd5b9392505050565b600080600080600060a0868803121561059b578081fd5b8535945060208601356105ad816107aa565b935060408601356105bd816107aa565b925060608601356105cd816107aa565b915060808601356105dd816107aa565b809150509295509295909350565b6000602082840312156105fc578081fd5b5035919050565b600060208284031215610614578081fd5b5051919050565b6000806040838503121561062d578182fd5b8235915060208084013567ffffffffffffffff8082111561064c578384fd5b818601915086601f83011261065f578384fd5b8135818111156106715761067161077b565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f830116810181811085821117156106b4576106b461077b565b604052828152858101935084860182860187018b10156106d2578788fd5b8795505b838610156106f45780358552600195909501949386019386016106d6565b508096505050505050509250929050565b600060208284031215610716578081fd5b815161057d816107aa565b6020808252825182820181905260009190848201906040850190845b8181101561076f57835173ffffffffffffffffffffffffffffffffffffffff168352928401929184019160010161073d565b50909695505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b67ffffffffffffffff811681146107c057600080fd5b5056fea164736f6c6343000804000a"

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

func (_VRFConsumerV2 *VRFConsumerV2Caller) RandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFConsumerV2.contract.Call(opts, &out, "randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFConsumerV2 *VRFConsumerV2Session) RandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFConsumerV2.Contract.RandomWords(&_VRFConsumerV2.CallOpts, arg0)
}

func (_VRFConsumerV2 *VRFConsumerV2CallerSession) RandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFConsumerV2.Contract.RandomWords(&_VRFConsumerV2.CallOpts, arg0)
}

func (_VRFConsumerV2 *VRFConsumerV2Caller) RequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFConsumerV2.contract.Call(opts, &out, "requestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFConsumerV2 *VRFConsumerV2Session) RequestId() (*big.Int, error) {
	return _VRFConsumerV2.Contract.RequestId(&_VRFConsumerV2.CallOpts)
}

func (_VRFConsumerV2 *VRFConsumerV2CallerSession) RequestId() (*big.Int, error) {
	return _VRFConsumerV2.Contract.RequestId(&_VRFConsumerV2.CallOpts)
}

func (_VRFConsumerV2 *VRFConsumerV2Caller) SubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFConsumerV2.contract.Call(opts, &out, "subId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFConsumerV2 *VRFConsumerV2Session) SubId() (uint64, error) {
	return _VRFConsumerV2.Contract.SubId(&_VRFConsumerV2.CallOpts)
}

func (_VRFConsumerV2 *VRFConsumerV2CallerSession) SubId() (uint64, error) {
	return _VRFConsumerV2.Contract.SubId(&_VRFConsumerV2.CallOpts)
}

func (_VRFConsumerV2 *VRFConsumerV2Transactor) FulfillRandomWords(opts *bind.TransactOpts, _requestId *big.Int, _randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.contract.Transact(opts, "fulfillRandomWords", _requestId, _randomWords)
}

func (_VRFConsumerV2 *VRFConsumerV2Session) FulfillRandomWords(_requestId *big.Int, _randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.FulfillRandomWords(&_VRFConsumerV2.TransactOpts, _requestId, _randomWords)
}

func (_VRFConsumerV2 *VRFConsumerV2TransactorSession) FulfillRandomWords(_requestId *big.Int, _randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.FulfillRandomWords(&_VRFConsumerV2.TransactOpts, _requestId, _randomWords)
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

func (_VRFConsumerV2 *VRFConsumerV2Transactor) TestRequestRandomness(opts *bind.TransactOpts, _keyHash [32]byte, _subId uint64, minReqConfs uint64, callbackGasLimit uint64, numWords uint64) (*types.Transaction, error) {
	return _VRFConsumerV2.contract.Transact(opts, "testRequestRandomness", _keyHash, _subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFConsumerV2 *VRFConsumerV2Session) TestRequestRandomness(_keyHash [32]byte, _subId uint64, minReqConfs uint64, callbackGasLimit uint64, numWords uint64) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.TestRequestRandomness(&_VRFConsumerV2.TransactOpts, _keyHash, _subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFConsumerV2 *VRFConsumerV2TransactorSession) TestRequestRandomness(_keyHash [32]byte, _subId uint64, minReqConfs uint64, callbackGasLimit uint64, numWords uint64) (*types.Transaction, error) {
	return _VRFConsumerV2.Contract.TestRequestRandomness(&_VRFConsumerV2.TransactOpts, _keyHash, _subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFConsumerV2 *VRFConsumerV2) Address() common.Address {
	return _VRFConsumerV2.address
}

type VRFConsumerV2Interface interface {
	RandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	RequestId(opts *bind.CallOpts) (*big.Int, error)

	SubId(opts *bind.CallOpts) (uint64, error)

	FulfillRandomWords(opts *bind.TransactOpts, _requestId *big.Int, _randomWords []*big.Int) (*types.Transaction, error)

	TestCreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	TestRequestRandomness(opts *bind.TransactOpts, _keyHash [32]byte, _subId uint64, minReqConfs uint64, callbackGasLimit uint64, numWords uint64) (*types.Transaction, error)

	Address() common.Address
}
