// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_consumer_v2_upgradeable_example

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
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated"
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

var VRFConsumerV2UpgradeableExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COORDINATOR\",\"outputs\":[{\"internalType\":\"contractVRFCoordinatorV2Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINKTOKEN\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"createSubscriptionAndFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"minReqConfs\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"}],\"name\":\"requestRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_gasAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_subId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"topUpSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"name\":\"updateSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506111c9806100206000396000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c806355380dfb11610081578063e89e106a1161005b578063e89e106a1461020a578063f08c5daa14610213578063f6eaffc81461021c57600080fd5b806355380dfb14610192578063706da1ca146101b2578063cf62c8ab146101f757600080fd5b806336bfffed116100b257806336bfffed146101275780633b2bcbf11461013a578063485cc9551461017f57600080fd5b8063177b9692146100d95780631fe543e3146100ff5780632fa4e44214610114575b600080fd5b6100ec6100e7366004610e42565b61022f565b6040519081526020015b60405180910390f35b61011261010d366004610edd565b610311565b005b610112610122366004610f9e565b6103a2565b610112610135366004610d75565b610502565b60345461015a9073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100f6565b61011261018d366004610d42565b61068a565b60355461015a9073ffffffffffffffffffffffffffffffffffffffff1681565b6035546101de9074010000000000000000000000000000000000000000900467ffffffffffffffff1681565b60405167ffffffffffffffff90911681526020016100f6565b610112610205366004610f9e565b610874565b6100ec60335481565b6100ec60365481565b6100ec61022a366004610eab565b610a7b565b6034546040517f5d3b1d300000000000000000000000000000000000000000000000000000000081526004810187905267ffffffffffffffff8616602482015261ffff8516604482015263ffffffff80851660648301528316608482015260009173ffffffffffffffffffffffffffffffffffffffff1690635d3b1d309060a401602060405180830381600087803b1580156102ca57600080fd5b505af11580156102de573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103029190610ec4565b60338190559695505050505050565b60005462010000900473ffffffffffffffffffffffffffffffffffffffff163314610394576000546040517f1cf993f40000000000000000000000000000000000000000000000000000000081523360048201526201000090910473ffffffffffffffffffffffffffffffffffffffff1660248201526044015b60405180910390fd5b61039e8282610a9c565b5050565b60355474010000000000000000000000000000000000000000900467ffffffffffffffff1661042d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f737562206e6f7420736574000000000000000000000000000000000000000000604482015260640161038b565b6035546034546040805174010000000000000000000000000000000000000000840467ffffffffffffffff16602082015273ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918591015b6040516020818303038152906040526040518463ffffffff1660e01b81526004016104b093929190610fcc565b602060405180830381600087803b1580156104ca57600080fd5b505af11580156104de573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061039e9190610e19565b60355474010000000000000000000000000000000000000000900467ffffffffffffffff1661058d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f7375624944206e6f742073657400000000000000000000000000000000000000604482015260640161038b565b60005b815181101561039e57603454603554835173ffffffffffffffffffffffffffffffffffffffff90921691637341c10c9174010000000000000000000000000000000000000000900467ffffffffffffffff16908590859081106105f5576105f5611145565b60200260200101516040518363ffffffff1660e01b815260040161064592919067ffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b15801561065f57600080fd5b505af1158015610673573d6000803e3d6000fd5b505050508080610682906110e5565b915050610590565b600054610100900460ff16158080156106aa5750600054600160ff909116105b806106c45750303b1580156106c4575060005460ff166001145b610750576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a6564000000000000000000000000000000000000606482015260840161038b565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117905580156107ae57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b6107b783610b1e565b6034805473ffffffffffffffffffffffffffffffffffffffff8086167fffffffffffffffffffffffff0000000000000000000000000000000000000000928316179092556035805492851692909116919091179055801561086f57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b505050565b60355474010000000000000000000000000000000000000000900467ffffffffffffffff1661042d57603460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a21a23e46040518163ffffffff1660e01b8152600401602060405180830381600087803b15801561090757600080fd5b505af115801561091b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061093f9190610f81565b603580547fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000067ffffffffffffffff938416810291909117918290556034546040517f7341c10c00000000000000000000000000000000000000000000000000000000815291909204909216600483015230602483015273ffffffffffffffffffffffffffffffffffffffff1690637341c10c90604401600060405180830381600087803b158015610a0657600080fd5b505af1158015610a1a573d6000803e3d6000fd5b50506035546034546040805174010000000000000000000000000000000000000000840467ffffffffffffffff16602082015273ffffffffffffffffffffffffffffffffffffffff9384169550634000aea094509290911691859101610483565b60328181548110610a8b57600080fd5b600091825260209091200154905081565b6033548214610b07576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265717565737420494420697320696e636f7272656374000000000000000000604482015260640161038b565b5a603655805161086f906032906020840190610ca5565b600054610100900460ff16610bb5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e67000000000000000000000000000000000000000000606482015260840161038b565b73ffffffffffffffffffffffffffffffffffffffff8116610c58576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f6d75737420676976652076616c696420636f6f7264696e61746f72206164647260448201527f6573730000000000000000000000000000000000000000000000000000000000606482015260840161038b565b6000805473ffffffffffffffffffffffffffffffffffffffff90921662010000027fffffffffffffffffffff0000000000000000000000000000000000000000ffff909216919091179055565b828054828255906000526020600020908101928215610ce0579160200282015b82811115610ce0578251825591602001919060010190610cc5565b50610cec929150610cf0565b5090565b5b80821115610cec5760008155600101610cf1565b803573ffffffffffffffffffffffffffffffffffffffff81168114610d2957600080fd5b919050565b803563ffffffff81168114610d2957600080fd5b60008060408385031215610d5557600080fd5b610d5e83610d05565b9150610d6c60208401610d05565b90509250929050565b60006020808385031215610d8857600080fd5b823567ffffffffffffffff811115610d9f57600080fd5b8301601f81018513610db057600080fd5b8035610dc3610dbe826110c1565b611072565b80828252848201915084840188868560051b8701011115610de357600080fd5b600094505b83851015610e0d57610df981610d05565b835260019490940193918501918501610de8565b50979650505050505050565b600060208284031215610e2b57600080fd5b81518015158114610e3b57600080fd5b9392505050565b600080600080600060a08688031215610e5a57600080fd5b853594506020860135610e6c816111a3565b9350604086013561ffff81168114610e8357600080fd5b9250610e9160608701610d2e565b9150610e9f60808701610d2e565b90509295509295909350565b600060208284031215610ebd57600080fd5b5035919050565b600060208284031215610ed657600080fd5b5051919050565b60008060408385031215610ef057600080fd5b8235915060208084013567ffffffffffffffff811115610f0f57600080fd5b8401601f81018613610f2057600080fd5b8035610f2e610dbe826110c1565b80828252848201915084840189868560051b8701011115610f4e57600080fd5b600094505b83851015610f71578035835260019490940193918501918501610f53565b5080955050505050509250929050565b600060208284031215610f9357600080fd5b8151610e3b816111a3565b600060208284031215610fb057600080fd5b81356bffffffffffffffffffffffff81168114610e3b57600080fd5b73ffffffffffffffffffffffffffffffffffffffff84168152600060206bffffffffffffffffffffffff85168184015260606040840152835180606085015260005b8181101561102a5785810183015185820160800152820161100e565b8181111561103c576000608083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160800195945050505050565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156110b9576110b9611174565b604052919050565b600067ffffffffffffffff8211156110db576110db611174565b5060051b60200190565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82141561113e577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b67ffffffffffffffff811681146111b957600080fd5b5056fea164736f6c6343000806000a",
}

var VRFConsumerV2UpgradeableExampleABI = VRFConsumerV2UpgradeableExampleMetaData.ABI

var VRFConsumerV2UpgradeableExampleBin = VRFConsumerV2UpgradeableExampleMetaData.Bin

func DeployVRFConsumerV2UpgradeableExample(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *VRFConsumerV2UpgradeableExample, error) {
	parsed, err := VRFConsumerV2UpgradeableExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFConsumerV2UpgradeableExampleBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFConsumerV2UpgradeableExample{VRFConsumerV2UpgradeableExampleCaller: VRFConsumerV2UpgradeableExampleCaller{contract: contract}, VRFConsumerV2UpgradeableExampleTransactor: VRFConsumerV2UpgradeableExampleTransactor{contract: contract}, VRFConsumerV2UpgradeableExampleFilterer: VRFConsumerV2UpgradeableExampleFilterer{contract: contract}}, nil
}

type VRFConsumerV2UpgradeableExample struct {
	address common.Address
	abi     abi.ABI
	VRFConsumerV2UpgradeableExampleCaller
	VRFConsumerV2UpgradeableExampleTransactor
	VRFConsumerV2UpgradeableExampleFilterer
}

type VRFConsumerV2UpgradeableExampleCaller struct {
	contract *bind.BoundContract
}

type VRFConsumerV2UpgradeableExampleTransactor struct {
	contract *bind.BoundContract
}

type VRFConsumerV2UpgradeableExampleFilterer struct {
	contract *bind.BoundContract
}

type VRFConsumerV2UpgradeableExampleSession struct {
	Contract     *VRFConsumerV2UpgradeableExample
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFConsumerV2UpgradeableExampleCallerSession struct {
	Contract *VRFConsumerV2UpgradeableExampleCaller
	CallOpts bind.CallOpts
}

type VRFConsumerV2UpgradeableExampleTransactorSession struct {
	Contract     *VRFConsumerV2UpgradeableExampleTransactor
	TransactOpts bind.TransactOpts
}

type VRFConsumerV2UpgradeableExampleRaw struct {
	Contract *VRFConsumerV2UpgradeableExample
}

type VRFConsumerV2UpgradeableExampleCallerRaw struct {
	Contract *VRFConsumerV2UpgradeableExampleCaller
}

type VRFConsumerV2UpgradeableExampleTransactorRaw struct {
	Contract *VRFConsumerV2UpgradeableExampleTransactor
}

func NewVRFConsumerV2UpgradeableExample(address common.Address, backend bind.ContractBackend) (*VRFConsumerV2UpgradeableExample, error) {
	abi, err := abi.JSON(strings.NewReader(VRFConsumerV2UpgradeableExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFConsumerV2UpgradeableExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2UpgradeableExample{address: address, abi: abi, VRFConsumerV2UpgradeableExampleCaller: VRFConsumerV2UpgradeableExampleCaller{contract: contract}, VRFConsumerV2UpgradeableExampleTransactor: VRFConsumerV2UpgradeableExampleTransactor{contract: contract}, VRFConsumerV2UpgradeableExampleFilterer: VRFConsumerV2UpgradeableExampleFilterer{contract: contract}}, nil
}

func NewVRFConsumerV2UpgradeableExampleCaller(address common.Address, caller bind.ContractCaller) (*VRFConsumerV2UpgradeableExampleCaller, error) {
	contract, err := bindVRFConsumerV2UpgradeableExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2UpgradeableExampleCaller{contract: contract}, nil
}

func NewVRFConsumerV2UpgradeableExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFConsumerV2UpgradeableExampleTransactor, error) {
	contract, err := bindVRFConsumerV2UpgradeableExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2UpgradeableExampleTransactor{contract: contract}, nil
}

func NewVRFConsumerV2UpgradeableExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFConsumerV2UpgradeableExampleFilterer, error) {
	contract, err := bindVRFConsumerV2UpgradeableExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2UpgradeableExampleFilterer{contract: contract}, nil
}

func bindVRFConsumerV2UpgradeableExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFConsumerV2UpgradeableExampleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFConsumerV2UpgradeableExample.Contract.VRFConsumerV2UpgradeableExampleCaller.contract.Call(opts, result, method, params...)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.VRFConsumerV2UpgradeableExampleTransactor.contract.Transfer(opts)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.VRFConsumerV2UpgradeableExampleTransactor.contract.Transact(opts, method, params...)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFConsumerV2UpgradeableExample.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.contract.Transfer(opts)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.contract.Transact(opts, method, params...)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleCaller) COORDINATOR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFConsumerV2UpgradeableExample.contract.Call(opts, &out, "COORDINATOR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleSession) COORDINATOR() (common.Address, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.COORDINATOR(&_VRFConsumerV2UpgradeableExample.CallOpts)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleCallerSession) COORDINATOR() (common.Address, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.COORDINATOR(&_VRFConsumerV2UpgradeableExample.CallOpts)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleCaller) LINKTOKEN(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFConsumerV2UpgradeableExample.contract.Call(opts, &out, "LINKTOKEN")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleSession) LINKTOKEN() (common.Address, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.LINKTOKEN(&_VRFConsumerV2UpgradeableExample.CallOpts)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleCallerSession) LINKTOKEN() (common.Address, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.LINKTOKEN(&_VRFConsumerV2UpgradeableExample.CallOpts)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleCaller) SGasAvailable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFConsumerV2UpgradeableExample.contract.Call(opts, &out, "s_gasAvailable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleSession) SGasAvailable() (*big.Int, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.SGasAvailable(&_VRFConsumerV2UpgradeableExample.CallOpts)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleCallerSession) SGasAvailable() (*big.Int, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.SGasAvailable(&_VRFConsumerV2UpgradeableExample.CallOpts)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleCaller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFConsumerV2UpgradeableExample.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.SRandomWords(&_VRFConsumerV2UpgradeableExample.CallOpts, arg0)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleCallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.SRandomWords(&_VRFConsumerV2UpgradeableExample.CallOpts, arg0)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleCaller) SRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFConsumerV2UpgradeableExample.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleSession) SRequestId() (*big.Int, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.SRequestId(&_VRFConsumerV2UpgradeableExample.CallOpts)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleCallerSession) SRequestId() (*big.Int, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.SRequestId(&_VRFConsumerV2UpgradeableExample.CallOpts)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleCaller) SSubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFConsumerV2UpgradeableExample.contract.Call(opts, &out, "s_subId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleSession) SSubId() (uint64, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.SSubId(&_VRFConsumerV2UpgradeableExample.CallOpts)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleCallerSession) SSubId() (uint64, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.SSubId(&_VRFConsumerV2UpgradeableExample.CallOpts)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleTransactor) CreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2UpgradeableExample.contract.Transact(opts, "createSubscriptionAndFund", amount)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleSession) CreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.CreateSubscriptionAndFund(&_VRFConsumerV2UpgradeableExample.TransactOpts, amount)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleTransactorSession) CreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.CreateSubscriptionAndFund(&_VRFConsumerV2UpgradeableExample.TransactOpts, amount)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleTransactor) Initialize(opts *bind.TransactOpts, _vrfCoordinator common.Address, _link common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2UpgradeableExample.contract.Transact(opts, "initialize", _vrfCoordinator, _link)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleSession) Initialize(_vrfCoordinator common.Address, _link common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.Initialize(&_VRFConsumerV2UpgradeableExample.TransactOpts, _vrfCoordinator, _link)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleTransactorSession) Initialize(_vrfCoordinator common.Address, _link common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.Initialize(&_VRFConsumerV2UpgradeableExample.TransactOpts, _vrfCoordinator, _link)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2UpgradeableExample.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.RawFulfillRandomWords(&_VRFConsumerV2UpgradeableExample.TransactOpts, requestId, randomWords)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.RawFulfillRandomWords(&_VRFConsumerV2UpgradeableExample.TransactOpts, requestId, randomWords)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleTransactor) RequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFConsumerV2UpgradeableExample.contract.Transact(opts, "requestRandomness", keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleSession) RequestRandomness(keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.RequestRandomness(&_VRFConsumerV2UpgradeableExample.TransactOpts, keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleTransactorSession) RequestRandomness(keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.RequestRandomness(&_VRFConsumerV2UpgradeableExample.TransactOpts, keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleTransactor) TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2UpgradeableExample.contract.Transact(opts, "topUpSubscription", amount)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.TopUpSubscription(&_VRFConsumerV2UpgradeableExample.TransactOpts, amount)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleTransactorSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.TopUpSubscription(&_VRFConsumerV2UpgradeableExample.TransactOpts, amount)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleTransactor) UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2UpgradeableExample.contract.Transact(opts, "updateSubscription", consumers)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.UpdateSubscription(&_VRFConsumerV2UpgradeableExample.TransactOpts, consumers)
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleTransactorSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2UpgradeableExample.Contract.UpdateSubscription(&_VRFConsumerV2UpgradeableExample.TransactOpts, consumers)
}

type VRFConsumerV2UpgradeableExampleInitializedIterator struct {
	Event *VRFConsumerV2UpgradeableExampleInitialized

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFConsumerV2UpgradeableExampleInitializedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFConsumerV2UpgradeableExampleInitialized)
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
		it.Event = new(VRFConsumerV2UpgradeableExampleInitialized)
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

func (it *VRFConsumerV2UpgradeableExampleInitializedIterator) Error() error {
	return it.fail
}

func (it *VRFConsumerV2UpgradeableExampleInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFConsumerV2UpgradeableExampleInitialized struct {
	Version uint8
	Raw     types.Log
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleFilterer) FilterInitialized(opts *bind.FilterOpts) (*VRFConsumerV2UpgradeableExampleInitializedIterator, error) {

	logs, sub, err := _VRFConsumerV2UpgradeableExample.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2UpgradeableExampleInitializedIterator{contract: _VRFConsumerV2UpgradeableExample.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *VRFConsumerV2UpgradeableExampleInitialized) (event.Subscription, error) {

	logs, sub, err := _VRFConsumerV2UpgradeableExample.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFConsumerV2UpgradeableExampleInitialized)
				if err := _VRFConsumerV2UpgradeableExample.contract.UnpackLog(event, "Initialized", log); err != nil {
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

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExampleFilterer) ParseInitialized(log types.Log) (*VRFConsumerV2UpgradeableExampleInitialized, error) {
	event := new(VRFConsumerV2UpgradeableExampleInitialized)
	if err := _VRFConsumerV2UpgradeableExample.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExample) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFConsumerV2UpgradeableExample.abi.Events["Initialized"].ID:
		return _VRFConsumerV2UpgradeableExample.ParseInitialized(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFConsumerV2UpgradeableExampleInitialized) Topic() common.Hash {
	return common.HexToHash("0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498")
}

func (_VRFConsumerV2UpgradeableExample *VRFConsumerV2UpgradeableExample) Address() common.Address {
	return _VRFConsumerV2UpgradeableExample.address
}

type VRFConsumerV2UpgradeableExampleInterface interface {
	COORDINATOR(opts *bind.CallOpts) (common.Address, error)

	LINKTOKEN(opts *bind.CallOpts) (common.Address, error)

	SGasAvailable(opts *bind.CallOpts) (*big.Int, error)

	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestId(opts *bind.CallOpts) (*big.Int, error)

	SSubId(opts *bind.CallOpts) (uint64, error)

	CreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	Initialize(opts *bind.TransactOpts, _vrfCoordinator common.Address, _link common.Address) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error)

	TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error)

	FilterInitialized(opts *bind.FilterOpts) (*VRFConsumerV2UpgradeableExampleInitializedIterator, error)

	WatchInitialized(opts *bind.WatchOpts, sink chan<- *VRFConsumerV2UpgradeableExampleInitialized) (event.Subscription, error)

	ParseInitialized(log types.Log) (*VRFConsumerV2UpgradeableExampleInitialized, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
