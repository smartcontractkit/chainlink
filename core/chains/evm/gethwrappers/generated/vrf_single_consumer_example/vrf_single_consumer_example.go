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
	_ = abi.ConvertType
)

var VRFSingleConsumerExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"fundAndRequestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestConfig\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"subscribe\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"topUpSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"unsubscribe\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b50604051620013383803806200133883398101604081905262000034916200031a565b606086811b6001600160601b0319166080908152600080546001600160a01b03808b166001600160a01b0319928316178355600180548316918b16919091179055600680543392169190911790556040805160a08101825291825263ffffffff8881166020840181905261ffff891692840183905290871694830185905291909201849052600280546001600160701b0319166801000000000000000090920261ffff60601b1916919091176c010000000000000000000000009092029190911763ffffffff60701b1916600160701b90920291909117905560038190556200011c62000128565b505050505050620003e2565b6006546001600160a01b031633146200014057600080fd5b604080516001808252818301909252600091602080830190803683370190505090503081600081518110620001795762000179620003cc565b60200260200101906001600160a01b031690816001600160a01b03168152505060008054906101000a90046001600160a01b03166001600160a01b031663a21a23e46040518163ffffffff1660e01b8152600401602060405180830381600087803b158015620001e857600080fd5b505af1158015620001fd573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200022391906200039a565b600280546001600160401b0319166001600160401b039290921691821790556000805483516001600160a01b0390911692637341c10c9290918591906200026e576200026e620003cc565b60200260200101516040518363ffffffff1660e01b8152600401620002b19291906001600160401b039290921682526001600160a01b0316602082015260400190565b600060405180830381600087803b158015620002cc57600080fd5b505af1158015620002e1573d6000803e3d6000fd5b5050505050565b80516001600160a01b03811681146200030057600080fd5b919050565b805163ffffffff811681146200030057600080fd5b60008060008060008060c087890312156200033457600080fd5b6200033f87620002e8565b95506200034f60208801620002e8565b94506200035f6040880162000305565b9350606087015161ffff811681146200037757600080fd5b9250620003876080880162000305565b915060a087015190509295509295509295565b600060208284031215620003ad57600080fd5b81516001600160401b0381168114620003c557600080fd5b9392505050565b634e487b7160e01b600052603260045260246000fd5b60805160601c610f3062000408600039600081816102ea01526103520152610f306000f3fe608060405234801561001057600080fd5b50600436106100bd5760003560e01c806386850e9311610076578063e0c862891161005b578063e0c86289146101cb578063e89e106a146101d3578063f6eaffc8146101ea57600080fd5b806386850e93146101b05780638f449a05146101c357600080fd5b80636fd700bb116100a75780636fd700bb146100ea5780637262561c146100fd5780637db9263f1461011057600080fd5b8062f714ce146100c25780631fe543e3146100d7575b600080fd5b6100d56100d0366004610ce8565b6101fd565b005b6100d56100e5366004610d14565b6102d2565b6100d56100f8366004610cb6565b610392565b6100d561010b366004610c72565b6105e0565b6002546003546101699167ffffffffffffffff81169163ffffffff68010000000000000000830481169261ffff6c01000000000000000000000000820416926e0100000000000000000000000000009091049091169085565b6040805167ffffffffffffffff909616865263ffffffff948516602087015261ffff90931692850192909252919091166060830152608082015260a0015b60405180910390f35b6100d56101be366004610cb6565b6106c8565b6100d56107ad565b6100d56109d4565b6101dc60055481565b6040519081526020016101a7565b6101dc6101f8366004610cb6565b610b35565b60065473ffffffffffffffffffffffffffffffffffffffff16331461022157600080fd5b6001546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8381166004830152602482018590529091169063a9059cbb90604401602060405180830381600087803b15801561029557600080fd5b505af11580156102a9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102cd9190610c94565b505050565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610384576040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660248201526044015b60405180910390fd5b61038e8282610b56565b5050565b60065473ffffffffffffffffffffffffffffffffffffffff1633146103b657600080fd5b6040805160a08101825260025467ffffffffffffffff811680835263ffffffff680100000000000000008304811660208086019190915261ffff6c01000000000000000000000000850416858701526e010000000000000000000000000000909304166060840152600354608084015260015460005485518085019390935285518084039094018452828601958690527f4000aea000000000000000000000000000000000000000000000000000000000909552929373ffffffffffffffffffffffffffffffffffffffff93841693634000aea09361049e9391909216918791604401610e2d565b602060405180830381600087803b1580156104b857600080fd5b505af11580156104cc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104f09190610c94565b50600054608082015182516040808501516020860151606087015192517f5d3b1d30000000000000000000000000000000000000000000000000000000008152600481019590955267ffffffffffffffff909316602485015261ffff16604484015263ffffffff918216606484015216608482015273ffffffffffffffffffffffffffffffffffffffff90911690635d3b1d309060a401602060405180830381600087803b1580156105a157600080fd5b505af11580156105b5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105d99190610ccf565b6005555050565b60065473ffffffffffffffffffffffffffffffffffffffff16331461060457600080fd5b6000546002546040517fd7ae1d3000000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015273ffffffffffffffffffffffffffffffffffffffff83811660248301529091169063d7ae1d3090604401600060405180830381600087803b15801561068557600080fd5b505af1158015610699573d6000803e3d6000fd5b5050600280547fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000169055505050565b60065473ffffffffffffffffffffffffffffffffffffffff1633146106ec57600080fd5b6001546000546002546040805167ffffffffffffffff909216602083015273ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918591016040516020818303038152906040526040518463ffffffff1660e01b815260040161075b93929190610e2d565b602060405180830381600087803b15801561077557600080fd5b505af1158015610789573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061038e9190610c94565b60065473ffffffffffffffffffffffffffffffffffffffff1633146107d157600080fd5b60408051600180825281830190925260009160208083019080368337019050509050308160008151811061080757610807610ec5565b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff168152505060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a21a23e46040518163ffffffff1660e01b8152600401602060405180830381600087803b1580156108a957600080fd5b505af11580156108bd573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108e19190610e03565b600280547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff92909216918217905560008054835173ffffffffffffffffffffffffffffffffffffffff90911692637341c10c92909185919061094f5761094f610ec5565b60200260200101516040518363ffffffff1660e01b815260040161099f92919067ffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b1580156109b957600080fd5b505af11580156109cd573d6000803e3d6000fd5b5050505050565b60065473ffffffffffffffffffffffffffffffffffffffff1633146109f857600080fd5b6040805160a08101825260025467ffffffffffffffff811680835263ffffffff68010000000000000000830481166020850181905261ffff6c010000000000000000000000008504168587018190526e010000000000000000000000000000909404909116606085018190526003546080860181905260005496517f5d3b1d3000000000000000000000000000000000000000000000000000000000815260048101919091526024810193909352604483019390935260648201526084810191909152909173ffffffffffffffffffffffffffffffffffffffff1690635d3b1d309060a401602060405180830381600087803b158015610af757600080fd5b505af1158015610b0b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b2f9190610ccf565b60055550565b60048181548110610b4557600080fd5b600091825260209091200154905081565b6005548214610bc1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265717565737420494420697320696e636f7272656374000000000000000000604482015260640161037b565b80516004805482825560008290526102cd927f8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19b91820191602086018215610c24579160200282015b82811115610c24578251825591602001919060010190610c09565b50610c30929150610c34565b5090565b5b80821115610c305760008155600101610c35565b803573ffffffffffffffffffffffffffffffffffffffff81168114610c6d57600080fd5b919050565b600060208284031215610c8457600080fd5b610c8d82610c49565b9392505050565b600060208284031215610ca657600080fd5b81518015158114610c8d57600080fd5b600060208284031215610cc857600080fd5b5035919050565b600060208284031215610ce157600080fd5b5051919050565b60008060408385031215610cfb57600080fd5b82359150610d0b60208401610c49565b90509250929050565b60008060408385031215610d2757600080fd5b8235915060208084013567ffffffffffffffff80821115610d4757600080fd5b818601915086601f830112610d5b57600080fd5b813581811115610d6d57610d6d610ef4565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f83011681018181108582111715610db057610db0610ef4565b604052828152858101935084860182860187018b1015610dcf57600080fd5b600095505b83861015610df2578035855260019590950194938601938601610dd4565b508096505050505050509250929050565b600060208284031215610e1557600080fd5b815167ffffffffffffffff81168114610c8d57600080fd5b73ffffffffffffffffffffffffffffffffffffffff8416815260006020848184015260606040840152835180606085015260005b81811015610e7d57858101830151858201608001528201610e61565b81811115610e8f576000608083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160800195945050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
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
	return address, tx, &VRFSingleConsumerExample{address: address, abi: *parsed, VRFSingleConsumerExampleCaller: VRFSingleConsumerExampleCaller{contract: contract}, VRFSingleConsumerExampleTransactor: VRFSingleConsumerExampleTransactor{contract: contract}, VRFSingleConsumerExampleFilterer: VRFSingleConsumerExampleFilterer{contract: contract}}, nil
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
	parsed, err := VRFSingleConsumerExampleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactor) Unsubscribe(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.contract.Transact(opts, "unsubscribe", to)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleSession) Unsubscribe(to common.Address) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.Unsubscribe(&_VRFSingleConsumerExample.TransactOpts, to)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactorSession) Unsubscribe(to common.Address) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.Unsubscribe(&_VRFSingleConsumerExample.TransactOpts, to)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactor) Withdraw(opts *bind.TransactOpts, amount *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.contract.Transact(opts, "withdraw", amount, to)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleSession) Withdraw(amount *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.Withdraw(&_VRFSingleConsumerExample.TransactOpts, amount, to)
}

func (_VRFSingleConsumerExample *VRFSingleConsumerExampleTransactorSession) Withdraw(amount *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFSingleConsumerExample.Contract.Withdraw(&_VRFSingleConsumerExample.TransactOpts, amount, to)
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

	Unsubscribe(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Withdraw(opts *bind.TransactOpts, amount *big.Int, to common.Address) (*types.Transaction, error)

	Address() common.Address
}
