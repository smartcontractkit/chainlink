// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrfv2plus_reverting_example

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

var VRFV2PlusRevertingExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"}],\"name\":\"OnlyOwnerOrCoordinator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"}],\"name\":\"CoordinatorSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"createSubscriptionAndFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"minReqConfs\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"}],\"name\":\"requestRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_gasAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_subId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_vrfCoordinator\",\"outputs\":[{\"internalType\":\"contractIVRFCoordinatorV2Plus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"}],\"name\":\"setCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"topUpSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"name\":\"updateSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001237380380620012378339810160408190526200003491620001e8565b8133806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf8162000120565b5050506001600160a01b038116620000ea5760405163d92e233d60e01b815260040160405180910390fd5b600280546001600160a01b039283166001600160a01b031991821617909155600580549390921692169190911790555062000220565b336001600160a01b038216036200017a5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001e357600080fd5b919050565b60008060408385031215620001fc57600080fd5b6200020783620001cb565b91506200021760208401620001cb565b90509250929050565b61100780620002306000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c80638ea981171161008c578063e89e106a11610066578063e89e106a146101e6578063f08c5daa146101ef578063f2fde38b146101f8578063f6eaffc81461020b57600080fd5b80638ea98117146101a05780639eccacf6146101b3578063cf62c8ab146101d357600080fd5b806336bfffed116100c857806336bfffed1461013d578063706da1ca1461015057806379ba5097146101595780638da5cb5b1461016157600080fd5b80631fe543e3146100ef5780632e75964e146101045780632fa4e4421461012a575b600080fd5b6101026100fd366004610b85565b61021e565b005b610117610112366004610c1d565b6102a3565b6040519081526020015b60405180910390f35b610102610138366004610c7d565b610391565b61010261014b366004610d05565b6104ab565b61011760065481565b6101026105e6565b60005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610121565b6101026101ae366004610de8565b6106e3565b60025461017b9073ffffffffffffffffffffffffffffffffffffffff1681565b6101026101e1366004610c7d565b61086d565b61011760045481565b61011760075481565b610102610206366004610de8565b6109d8565b610117610219366004610e03565b6109ec565b60025473ffffffffffffffffffffffffffffffffffffffff163314610296576002546040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff90911660248201526044015b60405180910390fd5b61029e600080fd5b505050565b6040805160c081018252868152602080820187905261ffff86168284015263ffffffff80861660608401528416608083015282519081018352600080825260a083019190915260025492517f9b1c385e000000000000000000000000000000000000000000000000000000008152909273ffffffffffffffffffffffffffffffffffffffff1690639b1c385e9061033e908490600401610e80565b6020604051808303816000875af115801561035d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103819190610ee5565b6004819055979650505050505050565b6006546000036103fd576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f737562206e6f7420736574000000000000000000000000000000000000000000604482015260640161028d565b60055460025460065460408051602081019290925273ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918591015b6040516020818303038152906040526040518463ffffffff1660e01b815260040161046493929190610efe565b6020604051808303816000875af1158015610483573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104a79190610f4a565b5050565b600654600003610517576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f7375624944206e6f742073657400000000000000000000000000000000000000604482015260640161028d565b60005b81518110156104a757600254600654835173ffffffffffffffffffffffffffffffffffffffff9092169163bec4c08c919085908590811061055d5761055d610f6c565b60200260200101516040518363ffffffff1660e01b81526004016105a192919091825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b1580156105bb57600080fd5b505af11580156105cf573d6000803e3d6000fd5b5050505080806105de90610f9b565b91505061051a565b60015473ffffffffffffffffffffffffffffffffffffffff163314610667576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161028d565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60005473ffffffffffffffffffffffffffffffffffffffff163314801590610723575060025473ffffffffffffffffffffffffffffffffffffffff163314155b156107a7573361074860005473ffffffffffffffffffffffffffffffffffffffff1690565b6002546040517f061db9c100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9384166004820152918316602483015291909116604482015260640161028d565b73ffffffffffffffffffffffffffffffffffffffff81166107f4576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527fd1a6a14209a385a964d036e404cb5cfb71f4000cdb03c9366292430787261be69060200160405180910390a150565b6006546000036103fd57600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a21a23e46040518163ffffffff1660e01b81526004016020604051808303816000875af11580156108e6573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061090a9190610ee5565b60068190556002546040517fbec4c08c000000000000000000000000000000000000000000000000000000008152600481019290925230602483015273ffffffffffffffffffffffffffffffffffffffff169063bec4c08c90604401600060405180830381600087803b15801561098057600080fd5b505af1158015610994573d6000803e3d6000fd5b5050505060055460025460065460405173ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918591610437919060200190815260200190565b6109e0610a0d565b6109e981610a90565b50565b600381815481106109fc57600080fd5b600091825260209091200154905081565b60005473ffffffffffffffffffffffffffffffffffffffff163314610a8e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161028d565b565b3373ffffffffffffffffffffffffffffffffffffffff821603610b0f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161028d565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600080600060408486031215610b9a57600080fd5b83359250602084013567ffffffffffffffff80821115610bb957600080fd5b818601915086601f830112610bcd57600080fd5b813581811115610bdc57600080fd5b8760208260051b8501011115610bf157600080fd5b6020830194508093505050509250925092565b803563ffffffff81168114610c1857600080fd5b919050565b600080600080600060a08688031215610c3557600080fd5b8535945060208601359350604086013561ffff81168114610c5557600080fd5b9250610c6360608701610c04565b9150610c7160808701610c04565b90509295509295909350565b600060208284031215610c8f57600080fd5b81356bffffffffffffffffffffffff81168114610cab57600080fd5b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b803573ffffffffffffffffffffffffffffffffffffffff81168114610c1857600080fd5b60006020808385031215610d1857600080fd5b823567ffffffffffffffff80821115610d3057600080fd5b818501915085601f830112610d4457600080fd5b813581811115610d5657610d56610cb2565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f83011681018181108582111715610d9957610d99610cb2565b604052918252848201925083810185019188831115610db757600080fd5b938501935b82851015610ddc57610dcd85610ce1565b84529385019392850192610dbc565b98975050505050505050565b600060208284031215610dfa57600080fd5b610cab82610ce1565b600060208284031215610e1557600080fd5b5035919050565b6000815180845260005b81811015610e4257602081850181015186830182015201610e26565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b60208152815160208201526020820151604082015261ffff60408301511660608201526000606083015163ffffffff80821660808501528060808601511660a0850152505060a083015160c080840152610edd60e0840182610e1c565b949350505050565b600060208284031215610ef757600080fd5b5051919050565b73ffffffffffffffffffffffffffffffffffffffff841681526bffffffffffffffffffffffff83166020820152606060408201526000610f416060830184610e1c565b95945050505050565b600060208284031215610f5c57600080fd5b81518015158114610cab57600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203610ff3577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b506001019056fea164736f6c6343000813000a",
}

var VRFV2PlusRevertingExampleABI = VRFV2PlusRevertingExampleMetaData.ABI

var VRFV2PlusRevertingExampleBin = VRFV2PlusRevertingExampleMetaData.Bin

func DeployVRFV2PlusRevertingExample(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, link common.Address) (common.Address, *types.Transaction, *VRFV2PlusRevertingExample, error) {
	parsed, err := VRFV2PlusRevertingExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2PlusRevertingExampleBin), backend, vrfCoordinator, link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2PlusRevertingExample{address: address, abi: *parsed, VRFV2PlusRevertingExampleCaller: VRFV2PlusRevertingExampleCaller{contract: contract}, VRFV2PlusRevertingExampleTransactor: VRFV2PlusRevertingExampleTransactor{contract: contract}, VRFV2PlusRevertingExampleFilterer: VRFV2PlusRevertingExampleFilterer{contract: contract}}, nil
}

type VRFV2PlusRevertingExample struct {
	address common.Address
	abi     abi.ABI
	VRFV2PlusRevertingExampleCaller
	VRFV2PlusRevertingExampleTransactor
	VRFV2PlusRevertingExampleFilterer
}

type VRFV2PlusRevertingExampleCaller struct {
	contract *bind.BoundContract
}

type VRFV2PlusRevertingExampleTransactor struct {
	contract *bind.BoundContract
}

type VRFV2PlusRevertingExampleFilterer struct {
	contract *bind.BoundContract
}

type VRFV2PlusRevertingExampleSession struct {
	Contract     *VRFV2PlusRevertingExample
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2PlusRevertingExampleCallerSession struct {
	Contract *VRFV2PlusRevertingExampleCaller
	CallOpts bind.CallOpts
}

type VRFV2PlusRevertingExampleTransactorSession struct {
	Contract     *VRFV2PlusRevertingExampleTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2PlusRevertingExampleRaw struct {
	Contract *VRFV2PlusRevertingExample
}

type VRFV2PlusRevertingExampleCallerRaw struct {
	Contract *VRFV2PlusRevertingExampleCaller
}

type VRFV2PlusRevertingExampleTransactorRaw struct {
	Contract *VRFV2PlusRevertingExampleTransactor
}

func NewVRFV2PlusRevertingExample(address common.Address, backend bind.ContractBackend) (*VRFV2PlusRevertingExample, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2PlusRevertingExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2PlusRevertingExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusRevertingExample{address: address, abi: abi, VRFV2PlusRevertingExampleCaller: VRFV2PlusRevertingExampleCaller{contract: contract}, VRFV2PlusRevertingExampleTransactor: VRFV2PlusRevertingExampleTransactor{contract: contract}, VRFV2PlusRevertingExampleFilterer: VRFV2PlusRevertingExampleFilterer{contract: contract}}, nil
}

func NewVRFV2PlusRevertingExampleCaller(address common.Address, caller bind.ContractCaller) (*VRFV2PlusRevertingExampleCaller, error) {
	contract, err := bindVRFV2PlusRevertingExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusRevertingExampleCaller{contract: contract}, nil
}

func NewVRFV2PlusRevertingExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2PlusRevertingExampleTransactor, error) {
	contract, err := bindVRFV2PlusRevertingExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusRevertingExampleTransactor{contract: contract}, nil
}

func NewVRFV2PlusRevertingExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2PlusRevertingExampleFilterer, error) {
	contract, err := bindVRFV2PlusRevertingExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusRevertingExampleFilterer{contract: contract}, nil
}

func bindVRFV2PlusRevertingExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFV2PlusRevertingExampleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusRevertingExample.Contract.VRFV2PlusRevertingExampleCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.VRFV2PlusRevertingExampleTransactor.contract.Transfer(opts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.VRFV2PlusRevertingExampleTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusRevertingExample.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.contract.Transfer(opts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusRevertingExample.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) Owner() (common.Address, error) {
	return _VRFV2PlusRevertingExample.Contract.Owner(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCallerSession) Owner() (common.Address, error) {
	return _VRFV2PlusRevertingExample.Contract.Owner(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCaller) SGasAvailable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusRevertingExample.contract.Call(opts, &out, "s_gasAvailable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) SGasAvailable() (*big.Int, error) {
	return _VRFV2PlusRevertingExample.Contract.SGasAvailable(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCallerSession) SGasAvailable() (*big.Int, error) {
	return _VRFV2PlusRevertingExample.Contract.SGasAvailable(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCaller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusRevertingExample.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFV2PlusRevertingExample.Contract.SRandomWords(&_VRFV2PlusRevertingExample.CallOpts, arg0)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFV2PlusRevertingExample.Contract.SRandomWords(&_VRFV2PlusRevertingExample.CallOpts, arg0)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCaller) SRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusRevertingExample.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) SRequestId() (*big.Int, error) {
	return _VRFV2PlusRevertingExample.Contract.SRequestId(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCallerSession) SRequestId() (*big.Int, error) {
	return _VRFV2PlusRevertingExample.Contract.SRequestId(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCaller) SSubId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusRevertingExample.contract.Call(opts, &out, "s_subId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) SSubId() (*big.Int, error) {
	return _VRFV2PlusRevertingExample.Contract.SSubId(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCallerSession) SSubId() (*big.Int, error) {
	return _VRFV2PlusRevertingExample.Contract.SSubId(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCaller) SVrfCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusRevertingExample.contract.Call(opts, &out, "s_vrfCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) SVrfCoordinator() (common.Address, error) {
	return _VRFV2PlusRevertingExample.Contract.SVrfCoordinator(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCallerSession) SVrfCoordinator() (common.Address, error) {
	return _VRFV2PlusRevertingExample.Contract.SVrfCoordinator(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "acceptOwnership")
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.AcceptOwnership(&_VRFV2PlusRevertingExample.TransactOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.AcceptOwnership(&_VRFV2PlusRevertingExample.TransactOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) CreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "createSubscriptionAndFund", amount)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) CreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.CreateSubscriptionAndFund(&_VRFV2PlusRevertingExample.TransactOpts, amount)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) CreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.CreateSubscriptionAndFund(&_VRFV2PlusRevertingExample.TransactOpts, amount)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.RawFulfillRandomWords(&_VRFV2PlusRevertingExample.TransactOpts, requestId, randomWords)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.RawFulfillRandomWords(&_VRFV2PlusRevertingExample.TransactOpts, requestId, randomWords)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) RequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subId *big.Int, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "requestRandomness", keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) RequestRandomness(keyHash [32]byte, subId *big.Int, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.RequestRandomness(&_VRFV2PlusRevertingExample.TransactOpts, keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) RequestRandomness(keyHash [32]byte, subId *big.Int, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.RequestRandomness(&_VRFV2PlusRevertingExample.TransactOpts, keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) SetCoordinator(opts *bind.TransactOpts, _vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "setCoordinator", _vrfCoordinator)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) SetCoordinator(_vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.SetCoordinator(&_VRFV2PlusRevertingExample.TransactOpts, _vrfCoordinator)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) SetCoordinator(_vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.SetCoordinator(&_VRFV2PlusRevertingExample.TransactOpts, _vrfCoordinator)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "topUpSubscription", amount)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.TopUpSubscription(&_VRFV2PlusRevertingExample.TransactOpts, amount)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.TopUpSubscription(&_VRFV2PlusRevertingExample.TransactOpts, amount)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.TransferOwnership(&_VRFV2PlusRevertingExample.TransactOpts, to)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.TransferOwnership(&_VRFV2PlusRevertingExample.TransactOpts, to)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "updateSubscription", consumers)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.UpdateSubscription(&_VRFV2PlusRevertingExample.TransactOpts, consumers)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.UpdateSubscription(&_VRFV2PlusRevertingExample.TransactOpts, consumers)
}

type VRFV2PlusRevertingExampleCoordinatorSetIterator struct {
	Event *VRFV2PlusRevertingExampleCoordinatorSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusRevertingExampleCoordinatorSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusRevertingExampleCoordinatorSet)
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
		it.Event = new(VRFV2PlusRevertingExampleCoordinatorSet)
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

func (it *VRFV2PlusRevertingExampleCoordinatorSetIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusRevertingExampleCoordinatorSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusRevertingExampleCoordinatorSet struct {
	VrfCoordinator common.Address
	Raw            types.Log
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleFilterer) FilterCoordinatorSet(opts *bind.FilterOpts) (*VRFV2PlusRevertingExampleCoordinatorSetIterator, error) {

	logs, sub, err := _VRFV2PlusRevertingExample.contract.FilterLogs(opts, "CoordinatorSet")
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusRevertingExampleCoordinatorSetIterator{contract: _VRFV2PlusRevertingExample.contract, event: "CoordinatorSet", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleFilterer) WatchCoordinatorSet(opts *bind.WatchOpts, sink chan<- *VRFV2PlusRevertingExampleCoordinatorSet) (event.Subscription, error) {

	logs, sub, err := _VRFV2PlusRevertingExample.contract.WatchLogs(opts, "CoordinatorSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusRevertingExampleCoordinatorSet)
				if err := _VRFV2PlusRevertingExample.contract.UnpackLog(event, "CoordinatorSet", log); err != nil {
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

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleFilterer) ParseCoordinatorSet(log types.Log) (*VRFV2PlusRevertingExampleCoordinatorSet, error) {
	event := new(VRFV2PlusRevertingExampleCoordinatorSet)
	if err := _VRFV2PlusRevertingExample.contract.UnpackLog(event, "CoordinatorSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2PlusRevertingExampleOwnershipTransferRequestedIterator struct {
	Event *VRFV2PlusRevertingExampleOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusRevertingExampleOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusRevertingExampleOwnershipTransferRequested)
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
		it.Event = new(VRFV2PlusRevertingExampleOwnershipTransferRequested)
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

func (it *VRFV2PlusRevertingExampleOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusRevertingExampleOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusRevertingExampleOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusRevertingExampleOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusRevertingExample.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusRevertingExampleOwnershipTransferRequestedIterator{contract: _VRFV2PlusRevertingExample.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2PlusRevertingExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusRevertingExample.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusRevertingExampleOwnershipTransferRequested)
				if err := _VRFV2PlusRevertingExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFV2PlusRevertingExampleOwnershipTransferRequested, error) {
	event := new(VRFV2PlusRevertingExampleOwnershipTransferRequested)
	if err := _VRFV2PlusRevertingExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2PlusRevertingExampleOwnershipTransferredIterator struct {
	Event *VRFV2PlusRevertingExampleOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusRevertingExampleOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusRevertingExampleOwnershipTransferred)
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
		it.Event = new(VRFV2PlusRevertingExampleOwnershipTransferred)
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

func (it *VRFV2PlusRevertingExampleOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusRevertingExampleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusRevertingExampleOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusRevertingExampleOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusRevertingExample.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusRevertingExampleOwnershipTransferredIterator{contract: _VRFV2PlusRevertingExample.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2PlusRevertingExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusRevertingExample.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusRevertingExampleOwnershipTransferred)
				if err := _VRFV2PlusRevertingExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleFilterer) ParseOwnershipTransferred(log types.Log) (*VRFV2PlusRevertingExampleOwnershipTransferred, error) {
	event := new(VRFV2PlusRevertingExampleOwnershipTransferred)
	if err := _VRFV2PlusRevertingExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExample) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFV2PlusRevertingExample.abi.Events["CoordinatorSet"].ID:
		return _VRFV2PlusRevertingExample.ParseCoordinatorSet(log)
	case _VRFV2PlusRevertingExample.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFV2PlusRevertingExample.ParseOwnershipTransferRequested(log)
	case _VRFV2PlusRevertingExample.abi.Events["OwnershipTransferred"].ID:
		return _VRFV2PlusRevertingExample.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFV2PlusRevertingExampleCoordinatorSet) Topic() common.Hash {
	return common.HexToHash("0xd1a6a14209a385a964d036e404cb5cfb71f4000cdb03c9366292430787261be6")
}

func (VRFV2PlusRevertingExampleOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFV2PlusRevertingExampleOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExample) Address() common.Address {
	return _VRFV2PlusRevertingExample.address
}

type VRFV2PlusRevertingExampleInterface interface {
	Owner(opts *bind.CallOpts) (common.Address, error)

	SGasAvailable(opts *bind.CallOpts) (*big.Int, error)

	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestId(opts *bind.CallOpts) (*big.Int, error)

	SSubId(opts *bind.CallOpts) (*big.Int, error)

	SVrfCoordinator(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	CreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subId *big.Int, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error)

	SetCoordinator(opts *bind.TransactOpts, _vrfCoordinator common.Address) (*types.Transaction, error)

	TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error)

	FilterCoordinatorSet(opts *bind.FilterOpts) (*VRFV2PlusRevertingExampleCoordinatorSetIterator, error)

	WatchCoordinatorSet(opts *bind.WatchOpts, sink chan<- *VRFV2PlusRevertingExampleCoordinatorSet) (event.Subscription, error)

	ParseCoordinatorSet(log types.Log) (*VRFV2PlusRevertingExampleCoordinatorSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusRevertingExampleOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2PlusRevertingExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFV2PlusRevertingExampleOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusRevertingExampleOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2PlusRevertingExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFV2PlusRevertingExampleOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
