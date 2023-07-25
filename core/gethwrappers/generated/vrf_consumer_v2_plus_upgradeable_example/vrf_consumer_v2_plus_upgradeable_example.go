// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_consumer_v2_plus_upgradeable_example

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

var VRFConsumerV2PlusUpgradeableExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COORDINATOR\",\"outputs\":[{\"internalType\":\"contractIVRFCoordinatorV2Plus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINKTOKEN\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"createSubscriptionAndFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"minReqConfs\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"}],\"name\":\"requestRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_gasAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_subId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"topUpSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"name\":\"updateSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506110c1806100206000396000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c806355380dfb11610081578063e89e106a1161005b578063e89e106a146101ce578063f08c5daa146101d7578063f6eaffc8146101e057600080fd5b806355380dfb14610192578063706da1ca146101b2578063cf62c8ab146101bb57600080fd5b806336bfffed116100b257806336bfffed146101275780633b2bcbf11461013a578063485cc9551461017f57600080fd5b80631fe543e3146100d95780632e75964e146100ee5780632fa4e44214610114575b600080fd5b6100ec6100e7366004610d95565b6101f3565b005b6101016100fc366004610d03565b610284565b6040519081526020015b60405180910390f35b6100ec610122366004610e39565b610381565b6100ec610135366004610c36565b6104a3565b60345461015a9073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161010b565b6100ec61018d366004610c03565b6105db565b60355461015a9073ffffffffffffffffffffffffffffffffffffffff1681565b61010160365481565b6100ec6101c9366004610e39565b6107c5565b61010160335481565b61010160375481565b6101016101ee366004610d63565b61093c565b60005462010000900473ffffffffffffffffffffffffffffffffffffffff163314610276576000546040517f1cf993f40000000000000000000000000000000000000000000000000000000081523360048201526201000090910473ffffffffffffffffffffffffffffffffffffffff1660248201526044015b60405180910390fd5b610280828261095d565b5050565b6040805160c081018252868152602080820187905261ffff86168284015263ffffffff80861660608401528416608083015282519081018352600080825260a083019190915260345492517f9b1c385e000000000000000000000000000000000000000000000000000000008152909273ffffffffffffffffffffffffffffffffffffffff1690639b1c385e9061031f908490600401610f1e565b602060405180830381600087803b15801561033957600080fd5b505af115801561034d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103719190610d7c565b6033819055979650505050505050565b6036546103ea576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f737562206e6f7420736574000000000000000000000000000000000000000000604482015260640161026d565b60355460345460365460408051602081019290925273ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918591015b6040516020818303038152906040526040518463ffffffff1660e01b815260040161045193929190610ed2565b602060405180830381600087803b15801561046b57600080fd5b505af115801561047f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102809190610cda565b60365461050c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f7375624944206e6f742073657400000000000000000000000000000000000000604482015260640161026d565b60005b815181101561028057603454603654835173ffffffffffffffffffffffffffffffffffffffff9092169163bec4c08c919085908590811061055257610552611056565b60200260200101516040518363ffffffff1660e01b815260040161059692919091825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b1580156105b057600080fd5b505af11580156105c4573d6000803e3d6000fd5b5050505080806105d390610ff6565b91505061050f565b600054610100900460ff16158080156105fb5750600054600160ff909116105b806106155750303b158015610615575060005460ff166001145b6106a1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a6564000000000000000000000000000000000000606482015260840161026d565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117905580156106ff57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b610708836109df565b6034805473ffffffffffffffffffffffffffffffffffffffff8086167fffffffffffffffffffffffff000000000000000000000000000000000000000092831617909255603580549285169290911691909117905580156107c057600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b505050565b6036546103ea57603460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a21a23e46040518163ffffffff1660e01b8152600401602060405180830381600087803b15801561083657600080fd5b505af115801561084a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061086e9190610d7c565b60368190556034546040517fbec4c08c000000000000000000000000000000000000000000000000000000008152600481019290925230602483015273ffffffffffffffffffffffffffffffffffffffff169063bec4c08c90604401600060405180830381600087803b1580156108e457600080fd5b505af11580156108f8573d6000803e3d6000fd5b5050505060355460345460365460405173ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918591610424919060200190815260200190565b6032818154811061094c57600080fd5b600091825260209091200154905081565b60335482146109c8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265717565737420494420697320696e636f7272656374000000000000000000604482015260640161026d565b5a60375580516107c0906032906020840190610b66565b600054610100900460ff16610a76576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e67000000000000000000000000000000000000000000606482015260840161026d565b73ffffffffffffffffffffffffffffffffffffffff8116610b19576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f6d75737420676976652076616c696420636f6f7264696e61746f72206164647260448201527f6573730000000000000000000000000000000000000000000000000000000000606482015260840161026d565b6000805473ffffffffffffffffffffffffffffffffffffffff90921662010000027fffffffffffffffffffff0000000000000000000000000000000000000000ffff909216919091179055565b828054828255906000526020600020908101928215610ba1579160200282015b82811115610ba1578251825591602001919060010190610b86565b50610bad929150610bb1565b5090565b5b80821115610bad5760008155600101610bb2565b803573ffffffffffffffffffffffffffffffffffffffff81168114610bea57600080fd5b919050565b803563ffffffff81168114610bea57600080fd5b60008060408385031215610c1657600080fd5b610c1f83610bc6565b9150610c2d60208401610bc6565b90509250929050565b60006020808385031215610c4957600080fd5b823567ffffffffffffffff811115610c6057600080fd5b8301601f81018513610c7157600080fd5b8035610c84610c7f82610fd2565b610f83565b80828252848201915084840188868560051b8701011115610ca457600080fd5b600094505b83851015610cce57610cba81610bc6565b835260019490940193918501918501610ca9565b50979650505050505050565b600060208284031215610cec57600080fd5b81518015158114610cfc57600080fd5b9392505050565b600080600080600060a08688031215610d1b57600080fd5b8535945060208601359350604086013561ffff81168114610d3b57600080fd5b9250610d4960608701610bef565b9150610d5760808701610bef565b90509295509295909350565b600060208284031215610d7557600080fd5b5035919050565b600060208284031215610d8e57600080fd5b5051919050565b60008060408385031215610da857600080fd5b8235915060208084013567ffffffffffffffff811115610dc757600080fd5b8401601f81018613610dd857600080fd5b8035610de6610c7f82610fd2565b80828252848201915084840189868560051b8701011115610e0657600080fd5b600094505b83851015610e29578035835260019490940193918501918501610e0b565b5080955050505050509250929050565b600060208284031215610e4b57600080fd5b81356bffffffffffffffffffffffff81168114610cfc57600080fd5b6000815180845260005b81811015610e8d57602081850181015186830182015201610e71565b81811115610e9f576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b73ffffffffffffffffffffffffffffffffffffffff841681526bffffffffffffffffffffffff83166020820152606060408201526000610f156060830184610e67565b95945050505050565b60208152815160208201526020820151604082015261ffff60408301511660608201526000606083015163ffffffff80821660808501528060808601511660a0850152505060a083015160c080840152610f7b60e0840182610e67565b949350505050565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610fca57610fca611085565b604052919050565b600067ffffffffffffffff821115610fec57610fec611085565b5060051b60200190565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82141561104f577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFConsumerV2PlusUpgradeableExampleABI = VRFConsumerV2PlusUpgradeableExampleMetaData.ABI

var VRFConsumerV2PlusUpgradeableExampleBin = VRFConsumerV2PlusUpgradeableExampleMetaData.Bin

func DeployVRFConsumerV2PlusUpgradeableExample(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *VRFConsumerV2PlusUpgradeableExample, error) {
	parsed, err := VRFConsumerV2PlusUpgradeableExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFConsumerV2PlusUpgradeableExampleBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFConsumerV2PlusUpgradeableExample{VRFConsumerV2PlusUpgradeableExampleCaller: VRFConsumerV2PlusUpgradeableExampleCaller{contract: contract}, VRFConsumerV2PlusUpgradeableExampleTransactor: VRFConsumerV2PlusUpgradeableExampleTransactor{contract: contract}, VRFConsumerV2PlusUpgradeableExampleFilterer: VRFConsumerV2PlusUpgradeableExampleFilterer{contract: contract}}, nil
}

type VRFConsumerV2PlusUpgradeableExample struct {
	address common.Address
	abi     abi.ABI
	VRFConsumerV2PlusUpgradeableExampleCaller
	VRFConsumerV2PlusUpgradeableExampleTransactor
	VRFConsumerV2PlusUpgradeableExampleFilterer
}

type VRFConsumerV2PlusUpgradeableExampleCaller struct {
	contract *bind.BoundContract
}

type VRFConsumerV2PlusUpgradeableExampleTransactor struct {
	contract *bind.BoundContract
}

type VRFConsumerV2PlusUpgradeableExampleFilterer struct {
	contract *bind.BoundContract
}

type VRFConsumerV2PlusUpgradeableExampleSession struct {
	Contract     *VRFConsumerV2PlusUpgradeableExample
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFConsumerV2PlusUpgradeableExampleCallerSession struct {
	Contract *VRFConsumerV2PlusUpgradeableExampleCaller
	CallOpts bind.CallOpts
}

type VRFConsumerV2PlusUpgradeableExampleTransactorSession struct {
	Contract     *VRFConsumerV2PlusUpgradeableExampleTransactor
	TransactOpts bind.TransactOpts
}

type VRFConsumerV2PlusUpgradeableExampleRaw struct {
	Contract *VRFConsumerV2PlusUpgradeableExample
}

type VRFConsumerV2PlusUpgradeableExampleCallerRaw struct {
	Contract *VRFConsumerV2PlusUpgradeableExampleCaller
}

type VRFConsumerV2PlusUpgradeableExampleTransactorRaw struct {
	Contract *VRFConsumerV2PlusUpgradeableExampleTransactor
}

func NewVRFConsumerV2PlusUpgradeableExample(address common.Address, backend bind.ContractBackend) (*VRFConsumerV2PlusUpgradeableExample, error) {
	abi, err := abi.JSON(strings.NewReader(VRFConsumerV2PlusUpgradeableExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFConsumerV2PlusUpgradeableExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2PlusUpgradeableExample{address: address, abi: abi, VRFConsumerV2PlusUpgradeableExampleCaller: VRFConsumerV2PlusUpgradeableExampleCaller{contract: contract}, VRFConsumerV2PlusUpgradeableExampleTransactor: VRFConsumerV2PlusUpgradeableExampleTransactor{contract: contract}, VRFConsumerV2PlusUpgradeableExampleFilterer: VRFConsumerV2PlusUpgradeableExampleFilterer{contract: contract}}, nil
}

func NewVRFConsumerV2PlusUpgradeableExampleCaller(address common.Address, caller bind.ContractCaller) (*VRFConsumerV2PlusUpgradeableExampleCaller, error) {
	contract, err := bindVRFConsumerV2PlusUpgradeableExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2PlusUpgradeableExampleCaller{contract: contract}, nil
}

func NewVRFConsumerV2PlusUpgradeableExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFConsumerV2PlusUpgradeableExampleTransactor, error) {
	contract, err := bindVRFConsumerV2PlusUpgradeableExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2PlusUpgradeableExampleTransactor{contract: contract}, nil
}

func NewVRFConsumerV2PlusUpgradeableExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFConsumerV2PlusUpgradeableExampleFilterer, error) {
	contract, err := bindVRFConsumerV2PlusUpgradeableExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2PlusUpgradeableExampleFilterer{contract: contract}, nil
}

func bindVRFConsumerV2PlusUpgradeableExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFConsumerV2PlusUpgradeableExampleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.VRFConsumerV2PlusUpgradeableExampleCaller.contract.Call(opts, result, method, params...)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.VRFConsumerV2PlusUpgradeableExampleTransactor.contract.Transfer(opts)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.VRFConsumerV2PlusUpgradeableExampleTransactor.contract.Transact(opts, method, params...)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.contract.Transfer(opts)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.contract.Transact(opts, method, params...)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleCaller) COORDINATOR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFConsumerV2PlusUpgradeableExample.contract.Call(opts, &out, "COORDINATOR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleSession) COORDINATOR() (common.Address, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.COORDINATOR(&_VRFConsumerV2PlusUpgradeableExample.CallOpts)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleCallerSession) COORDINATOR() (common.Address, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.COORDINATOR(&_VRFConsumerV2PlusUpgradeableExample.CallOpts)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleCaller) LINKTOKEN(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFConsumerV2PlusUpgradeableExample.contract.Call(opts, &out, "LINKTOKEN")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleSession) LINKTOKEN() (common.Address, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.LINKTOKEN(&_VRFConsumerV2PlusUpgradeableExample.CallOpts)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleCallerSession) LINKTOKEN() (common.Address, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.LINKTOKEN(&_VRFConsumerV2PlusUpgradeableExample.CallOpts)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleCaller) SGasAvailable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFConsumerV2PlusUpgradeableExample.contract.Call(opts, &out, "s_gasAvailable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleSession) SGasAvailable() (*big.Int, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.SGasAvailable(&_VRFConsumerV2PlusUpgradeableExample.CallOpts)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleCallerSession) SGasAvailable() (*big.Int, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.SGasAvailable(&_VRFConsumerV2PlusUpgradeableExample.CallOpts)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleCaller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFConsumerV2PlusUpgradeableExample.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.SRandomWords(&_VRFConsumerV2PlusUpgradeableExample.CallOpts, arg0)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleCallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.SRandomWords(&_VRFConsumerV2PlusUpgradeableExample.CallOpts, arg0)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleCaller) SRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFConsumerV2PlusUpgradeableExample.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleSession) SRequestId() (*big.Int, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.SRequestId(&_VRFConsumerV2PlusUpgradeableExample.CallOpts)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleCallerSession) SRequestId() (*big.Int, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.SRequestId(&_VRFConsumerV2PlusUpgradeableExample.CallOpts)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleCaller) SSubId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFConsumerV2PlusUpgradeableExample.contract.Call(opts, &out, "s_subId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleSession) SSubId() (*big.Int, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.SSubId(&_VRFConsumerV2PlusUpgradeableExample.CallOpts)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleCallerSession) SSubId() (*big.Int, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.SSubId(&_VRFConsumerV2PlusUpgradeableExample.CallOpts)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleTransactor) CreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2PlusUpgradeableExample.contract.Transact(opts, "createSubscriptionAndFund", amount)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleSession) CreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.CreateSubscriptionAndFund(&_VRFConsumerV2PlusUpgradeableExample.TransactOpts, amount)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleTransactorSession) CreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.CreateSubscriptionAndFund(&_VRFConsumerV2PlusUpgradeableExample.TransactOpts, amount)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleTransactor) Initialize(opts *bind.TransactOpts, _vrfCoordinator common.Address, _link common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2PlusUpgradeableExample.contract.Transact(opts, "initialize", _vrfCoordinator, _link)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleSession) Initialize(_vrfCoordinator common.Address, _link common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.Initialize(&_VRFConsumerV2PlusUpgradeableExample.TransactOpts, _vrfCoordinator, _link)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleTransactorSession) Initialize(_vrfCoordinator common.Address, _link common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.Initialize(&_VRFConsumerV2PlusUpgradeableExample.TransactOpts, _vrfCoordinator, _link)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2PlusUpgradeableExample.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.RawFulfillRandomWords(&_VRFConsumerV2PlusUpgradeableExample.TransactOpts, requestId, randomWords)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.RawFulfillRandomWords(&_VRFConsumerV2PlusUpgradeableExample.TransactOpts, requestId, randomWords)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleTransactor) RequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subId *big.Int, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFConsumerV2PlusUpgradeableExample.contract.Transact(opts, "requestRandomness", keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleSession) RequestRandomness(keyHash [32]byte, subId *big.Int, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.RequestRandomness(&_VRFConsumerV2PlusUpgradeableExample.TransactOpts, keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleTransactorSession) RequestRandomness(keyHash [32]byte, subId *big.Int, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.RequestRandomness(&_VRFConsumerV2PlusUpgradeableExample.TransactOpts, keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleTransactor) TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2PlusUpgradeableExample.contract.Transact(opts, "topUpSubscription", amount)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.TopUpSubscription(&_VRFConsumerV2PlusUpgradeableExample.TransactOpts, amount)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleTransactorSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.TopUpSubscription(&_VRFConsumerV2PlusUpgradeableExample.TransactOpts, amount)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleTransactor) UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2PlusUpgradeableExample.contract.Transact(opts, "updateSubscription", consumers)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.UpdateSubscription(&_VRFConsumerV2PlusUpgradeableExample.TransactOpts, consumers)
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleTransactorSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2PlusUpgradeableExample.Contract.UpdateSubscription(&_VRFConsumerV2PlusUpgradeableExample.TransactOpts, consumers)
}

type VRFConsumerV2PlusUpgradeableExampleInitializedIterator struct {
	Event *VRFConsumerV2PlusUpgradeableExampleInitialized

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFConsumerV2PlusUpgradeableExampleInitializedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFConsumerV2PlusUpgradeableExampleInitialized)
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
		it.Event = new(VRFConsumerV2PlusUpgradeableExampleInitialized)
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

func (it *VRFConsumerV2PlusUpgradeableExampleInitializedIterator) Error() error {
	return it.fail
}

func (it *VRFConsumerV2PlusUpgradeableExampleInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFConsumerV2PlusUpgradeableExampleInitialized struct {
	Version uint8
	Raw     types.Log
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleFilterer) FilterInitialized(opts *bind.FilterOpts) (*VRFConsumerV2PlusUpgradeableExampleInitializedIterator, error) {

	logs, sub, err := _VRFConsumerV2PlusUpgradeableExample.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2PlusUpgradeableExampleInitializedIterator{contract: _VRFConsumerV2PlusUpgradeableExample.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *VRFConsumerV2PlusUpgradeableExampleInitialized) (event.Subscription, error) {

	logs, sub, err := _VRFConsumerV2PlusUpgradeableExample.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFConsumerV2PlusUpgradeableExampleInitialized)
				if err := _VRFConsumerV2PlusUpgradeableExample.contract.UnpackLog(event, "Initialized", log); err != nil {
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

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExampleFilterer) ParseInitialized(log types.Log) (*VRFConsumerV2PlusUpgradeableExampleInitialized, error) {
	event := new(VRFConsumerV2PlusUpgradeableExampleInitialized)
	if err := _VRFConsumerV2PlusUpgradeableExample.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExample) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFConsumerV2PlusUpgradeableExample.abi.Events["Initialized"].ID:
		return _VRFConsumerV2PlusUpgradeableExample.ParseInitialized(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFConsumerV2PlusUpgradeableExampleInitialized) Topic() common.Hash {
	return common.HexToHash("0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498")
}

func (_VRFConsumerV2PlusUpgradeableExample *VRFConsumerV2PlusUpgradeableExample) Address() common.Address {
	return _VRFConsumerV2PlusUpgradeableExample.address
}

type VRFConsumerV2PlusUpgradeableExampleInterface interface {
	COORDINATOR(opts *bind.CallOpts) (common.Address, error)

	LINKTOKEN(opts *bind.CallOpts) (common.Address, error)

	SGasAvailable(opts *bind.CallOpts) (*big.Int, error)

	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestId(opts *bind.CallOpts) (*big.Int, error)

	SSubId(opts *bind.CallOpts) (*big.Int, error)

	CreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	Initialize(opts *bind.TransactOpts, _vrfCoordinator common.Address, _link common.Address) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subId *big.Int, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error)

	TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error)

	FilterInitialized(opts *bind.FilterOpts) (*VRFConsumerV2PlusUpgradeableExampleInitializedIterator, error)

	WatchInitialized(opts *bind.WatchOpts, sink chan<- *VRFConsumerV2PlusUpgradeableExampleInitialized) (event.Subscription, error)

	ParseInitialized(log types.Log) (*VRFConsumerV2PlusUpgradeableExampleInitialized, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
