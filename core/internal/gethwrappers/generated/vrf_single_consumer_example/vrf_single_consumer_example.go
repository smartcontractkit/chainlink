// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_single_consumer_example

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

var VRFSingleConsumerExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"fundAndRequestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestConfig\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"subscribe\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"topUpSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unsubscribe\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162001061380380620010618339810160408190526200003491620002f3565b606086811b6001600160601b0319166080908152600080546001600160a01b03808b166001600160a01b031992831617835560018054909216908a161790556040805160a08101825291825263ffffffff8881166020840181905261ffff891692840183905290871694830185905291909201849052600280546001600160701b0319166801000000000000000090920261ffff60601b1916919091176c010000000000000000000000009092029190911763ffffffff60701b1916600160701b90920291909117905560038190556200010d62000119565b505050505050620003bb565b604080516001808252818301909252600091602080830190803683370190505090503081600081518110620001525762000152620003a5565b60200260200101906001600160a01b031690816001600160a01b03168152505060008054906101000a90046001600160a01b03166001600160a01b031663a21a23e46040518163ffffffff1660e01b8152600401602060405180830381600087803b158015620001c157600080fd5b505af1158015620001d6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001fc919062000373565b600280546001600160401b0319166001600160401b039290921691821790556000805483516001600160a01b0390911692637341c10c929091859190620002475762000247620003a5565b60200260200101516040518363ffffffff1660e01b81526004016200028a9291906001600160401b039290921682526001600160a01b0316602082015260400190565b600060405180830381600087803b158015620002a557600080fd5b505af1158015620002ba573d6000803e3d6000fd5b5050505050565b80516001600160a01b0381168114620002d957600080fd5b919050565b805163ffffffff81168114620002d957600080fd5b60008060008060008060c087890312156200030d57600080fd5b6200031887620002c1565b95506200032860208801620002c1565b94506200033860408801620002de565b9350606087015161ffff811681146200035057600080fd5b92506200036060808801620002de565b915060a087015190509295509295509295565b6000602082840312156200038657600080fd5b81516001600160401b03811681146200039e57600080fd5b9392505050565b634e487b7160e01b600052603260045260246000fd5b60805160601c610c80620003e1600039600081816101dd01526102450152610c806000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c80638f449a0511610076578063e89e106a1161005b578063e89e106a14610193578063f6eaffc8146101aa578063fcae4484146101bd57600080fd5b80638f449a0514610183578063e0c862891461018b57600080fd5b80631fe543e3146100a85780636fd700bb146100bd5780637db9263f146100d057806386850e9314610170575b600080fd5b6100bb6100b6366004610a64565b6101c5565b005b6100bb6100cb366004610a32565b610284565b6002546003546101299167ffffffffffffffff81169163ffffffff68010000000000000000830481169261ffff6c01000000000000000000000000820416926e0100000000000000000000000000009091049091169085565b6040805167ffffffffffffffff909616865263ffffffff948516602087015261ffff90931692850192909252919091166060830152608082015260a0015b60405180910390f35b6100bb61017e366004610a32565b6104ae565b6100bb61056f565b6100bb610772565b61019c60055481565b604051908152602001610167565b61019c6101b8366004610a32565b6108af565b6100bb6108d0565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610276576040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016602482015260440160405180910390fd5b6102808282610991565b5050565b6040805160a08101825260025467ffffffffffffffff811680835263ffffffff680100000000000000008304811660208086019190915261ffff6c01000000000000000000000000850416858701526e010000000000000000000000000000909304166060840152600354608084015260015460005485518085019390935285518084039094018452828601958690527f4000aea000000000000000000000000000000000000000000000000000000000909552929373ffffffffffffffffffffffffffffffffffffffff93841693634000aea09361036c9391909216918791604401610b7d565b602060405180830381600087803b15801561038657600080fd5b505af115801561039a573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103be9190610a09565b50600054608082015182516040808501516020860151606087015192517f5d3b1d30000000000000000000000000000000000000000000000000000000008152600481019590955267ffffffffffffffff909316602485015261ffff16604484015263ffffffff918216606484015216608482015273ffffffffffffffffffffffffffffffffffffffff90911690635d3b1d309060a401602060405180830381600087803b15801561046f57600080fd5b505af1158015610483573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104a79190610a4b565b6005555050565b6001546000546002546040805167ffffffffffffffff909216602083015273ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918591016040516020818303038152906040526040518463ffffffff1660e01b815260040161051d93929190610b7d565b602060405180830381600087803b15801561053757600080fd5b505af115801561054b573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102809190610a09565b6040805160018082528183019092526000916020808301908036833701905050905030816000815181106105a5576105a5610c15565b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff168152505060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a21a23e46040518163ffffffff1660e01b8152600401602060405180830381600087803b15801561064757600080fd5b505af115801561065b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061067f9190610b53565b600280547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff92909216918217905560008054835173ffffffffffffffffffffffffffffffffffffffff90911692637341c10c9290918591906106ed576106ed610c15565b60200260200101516040518363ffffffff1660e01b815260040161073d92919067ffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b15801561075757600080fd5b505af115801561076b573d6000803e3d6000fd5b5050505050565b6040805160a08101825260025467ffffffffffffffff811680835263ffffffff68010000000000000000830481166020850181905261ffff6c010000000000000000000000008504168587018190526e010000000000000000000000000000909404909116606085018190526003546080860181905260005496517f5d3b1d3000000000000000000000000000000000000000000000000000000000815260048101919091526024810193909352604483019390935260648201526084810191909152909173ffffffffffffffffffffffffffffffffffffffff1690635d3b1d309060a401602060405180830381600087803b15801561087157600080fd5b505af1158015610885573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108a99190610a4b565b60055550565b600481815481106108bf57600080fd5b600091825260209091200154905081565b6000546002546040517fd7ae1d3000000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015230602482015273ffffffffffffffffffffffffffffffffffffffff9091169063d7ae1d3090604401600060405180830381600087803b15801561094f57600080fd5b505af1158015610963573d6000803e3d6000fd5b5050600280547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001690555050565b80516109a49060049060208401906109a9565b505050565b8280548282559060005260206000209081019282156109e4579160200282015b828111156109e45782518255916020019190600101906109c9565b506109f09291506109f4565b5090565b5b808211156109f057600081556001016109f5565b600060208284031215610a1b57600080fd5b81518015158114610a2b57600080fd5b9392505050565b600060208284031215610a4457600080fd5b5035919050565b600060208284031215610a5d57600080fd5b5051919050565b60008060408385031215610a7757600080fd5b8235915060208084013567ffffffffffffffff80821115610a9757600080fd5b818601915086601f830112610aab57600080fd5b813581811115610abd57610abd610c44565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f83011681018181108582111715610b0057610b00610c44565b604052828152858101935084860182860187018b1015610b1f57600080fd5b600095505b83861015610b42578035855260019590950194938601938601610b24565b508096505050505050509250929050565b600060208284031215610b6557600080fd5b815167ffffffffffffffff81168114610a2b57600080fd5b73ffffffffffffffffffffffffffffffffffffffff8416815260006020848184015260606040840152835180606085015260005b81811015610bcd57858101830151858201608001528201610bb1565b81811115610bdf576000608083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160800195945050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFSingleConsumerExampleABI = VRFSingleConsumerExampleMetaData.ABI

var VRFSingleConsumerExampleBin = VRFSingleConsumerExampleMetaData.Bin

func DeployVRFSingleConsumerExample(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, link common.Address, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte) (common.Address, *types.Transaction, *VRFSingleConsumerExample, error) {
	parsed, err := VRFSingleConsumerExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFSingleConsumerExampleBin), backend, vrfCoordinator, link, callbackGasLimit, requestConfirmations, numWords, keyHash)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFSingleConsumerExample{VRFSingleConsumerExampleCaller: VRFSingleConsumerExampleCaller{contract: contract}, VRFSingleConsumerExampleTransactor: VRFSingleConsumerExampleTransactor{contract: contract}, VRFSingleConsumerExampleFilterer: VRFSingleConsumerExampleFilterer{contract: contract}}, nil
}

type VRFSingleConsumerExample struct {
	address common.Address
	abi     abi.ABI
	VRFSingleConsumerExampleCaller
	VRFSingleConsumerExampleTransactor
	VRFSingleConsumerExampleFilterer
}

type VRFSingleConsumerExampleCaller struct {
	contract *bind.BoundContract
}

type VRFSingleConsumerExampleTransactor struct {
	contract *bind.BoundContract
}

type VRFSingleConsumerExampleFilterer struct {
	contract *bind.BoundContract
}

type VRFSingleConsumerExampleSession struct {
	Contract     *VRFSingleConsumerExample
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFSingleConsumerExampleCallerSession struct {
	Contract *VRFSingleConsumerExampleCaller
	CallOpts bind.CallOpts
}

type VRFSingleConsumerExampleTransactorSession struct {
	Contract     *VRFSingleConsumerExampleTransactor
	TransactOpts bind.TransactOpts
}

type VRFSingleConsumerExampleRaw struct {
	Contract *VRFSingleConsumerExample
}

type VRFSingleConsumerExampleCallerRaw struct {
	Contract *VRFSingleConsumerExampleCaller
}

type VRFSingleConsumerExampleTransactorRaw struct {
	Contract *VRFSingleConsumerExampleTransactor
}

func NewVRFSingleConsumerExample(address common.Address, backend bind.ContractBackend) (*VRFSingleConsumerExample, error) {
	abi, err := abi.JSON(strings.NewReader(VRFSingleConsumerExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFSingleConsumerExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFSingleConsumerExample{address: address, abi: abi, VRFSingleConsumerExampleCaller: VRFSingleConsumerExampleCaller{contract: contract}, VRFSingleConsumerExampleTransactor: VRFSingleConsumerExampleTransactor{contract: contract}, VRFSingleConsumerExampleFilterer: VRFSingleConsumerExampleFilterer{contract: contract}}, nil
}

func NewVRFSingleConsumerExampleCaller(address common.Address, caller bind.ContractCaller) (*VRFSingleConsumerExampleCaller, error) {
	contract, err := bindVRFSingleConsumerExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFSingleConsumerExampleCaller{contract: contract}, nil
}

func NewVRFSingleConsumerExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFSingleConsumerExampleTransactor, error) {
	contract, err := bindVRFSingleConsumerExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFSingleConsumerExampleTransactor{contract: contract}, nil
}

func NewVRFSingleConsumerExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFSingleConsumerExampleFilterer, error) {
	contract, err := bindVRFSingleConsumerExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFSingleConsumerExampleFilterer{contract: contract}, nil
}

func bindVRFSingleConsumerExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFSingleConsumerExampleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFSingleConsumerExample.Contract.VRFSingleConsumerExampleCaller.contract.Call(opts, result, method, params...)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.VRFSingleConsumerExampleTransactor.contract.Transfer(opts)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.VRFSingleConsumerExampleTransactor.contract.Transact(opts, method, params...)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFSingleConsumerExample.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.contract.Transfer(opts)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.contract.Transact(opts, method, params...)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleCaller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFSingleConsumerExample.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFSingleConsumerExample.Contract.SRandomWords(&_VRFSingleConsumerExample.CallOpts, arg0)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleCallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFSingleConsumerExample.Contract.SRandomWords(&_VRFSingleConsumerExample.CallOpts, arg0)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleCaller) SRequestConfig(opts *bind.CallOpts) (SRequestConfig,

	error) {
	var out []interface{}
	err := _VRFSingleConsumerExample.contract.Call(opts, &out, "s_requestConfig")

	outstruct := new(SRequestConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.SubId = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.CallbackGasLimit = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.RequestConfirmations = *abi.ConvertType(out[2], new(uint16)).(*uint16)
	outstruct.NumWords = *abi.ConvertType(out[3], new(uint32)).(*uint32)
	outstruct.KeyHash = *abi.ConvertType(out[4], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleSession) SRequestConfig() (SRequestConfig,

	error) {
	return _VRFSingleConsumerExample.Contract.SRequestConfig(&_VRFSingleConsumerExample.CallOpts)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleCallerSession) SRequestConfig() (SRequestConfig,

	error) {
	return _VRFSingleConsumerExample.Contract.SRequestConfig(&_VRFSingleConsumerExample.CallOpts)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleCaller) SRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFSingleConsumerExample.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleSession) SRequestId() (*big.Int, error) {
	return _VRFSingleConsumerExample.Contract.SRequestId(&_VRFSingleConsumerExample.CallOpts)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleCallerSession) SRequestId() (*big.Int, error) {
	return _VRFSingleConsumerExample.Contract.SRequestId(&_VRFSingleConsumerExample.CallOpts)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactor) FundAndRequestRandomWords(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.contract.Transact(opts, "fundAndRequestRandomWords", amount)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleSession) FundAndRequestRandomWords(amount *big.Int) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.FundAndRequestRandomWords(&_VRFSingleConsumerExample.TransactOpts, amount)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactorSession) FundAndRequestRandomWords(amount *big.Int) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.FundAndRequestRandomWords(&_VRFSingleConsumerExample.TransactOpts, amount)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.RawFulfillRandomWords(&_VRFSingleConsumerExample.TransactOpts, requestId, randomWords)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.RawFulfillRandomWords(&_VRFSingleConsumerExample.TransactOpts, requestId, randomWords)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactor) RequestRandomWords(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.contract.Transact(opts, "requestRandomWords")
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleSession) RequestRandomWords() (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.RequestRandomWords(&_VRFSingleConsumerExample.TransactOpts)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactorSession) RequestRandomWords() (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.RequestRandomWords(&_VRFSingleConsumerExample.TransactOpts)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactor) Subscribe(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.contract.Transact(opts, "subscribe")
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleSession) Subscribe() (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.Subscribe(&_VRFSingleConsumerExample.TransactOpts)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactorSession) Subscribe() (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.Subscribe(&_VRFSingleConsumerExample.TransactOpts)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactor) TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.contract.Transact(opts, "topUpSubscription", amount)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.TopUpSubscription(&_VRFSingleConsumerExample.TransactOpts, amount)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactorSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.TopUpSubscription(&_VRFSingleConsumerExample.TransactOpts, amount)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactor) Unsubscribe(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.contract.Transact(opts, "unsubscribe")
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleSession) Unsubscribe() (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.Unsubscribe(&_VRFSingleConsumerExample.TransactOpts)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactorSession) Unsubscribe() (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.Unsubscribe(&_VRFSingleConsumerExample.TransactOpts)
}

type SRequestConfig struct {
	SubId                uint64
	CallbackGasLimit     uint32
	RequestConfirmations uint16
	NumWords             uint32
	KeyHash              [32]byte
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExample) Address() common.Address {
	return _VRFSingleConsumerExample.address
}

type VRFSingleConsumerExampleInterface interface {
	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestConfig(opts *bind.CallOpts) (SRequestConfig,

		error)

	SRequestId(opts *bind.CallOpts) (*big.Int, error)

	FundAndRequestRandomWords(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts) (*types.Transaction, error)

	Subscribe(opts *bind.TransactOpts) (*types.Transaction, error)

	TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	Unsubscribe(opts *bind.TransactOpts) (*types.Transaction, error)

	Address() common.Address
}
