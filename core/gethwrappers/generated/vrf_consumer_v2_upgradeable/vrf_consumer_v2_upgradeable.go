// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_consumer_v2_upgradeable

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

var VRFConsumerV2UpgradeableMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COORDINATOR\",\"outputs\":[{\"internalType\":\"contractVRFCoordinatorV2Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINKTOKEN\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"}],\"name\":\"__VRFConsumerBaseV2_init\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_gasAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_subId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"testCreateSubscriptionAndFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"minReqConfs\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"}],\"name\":\"testRequestRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"topUpSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"name\":\"updateSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50611262806100206000396000f3fe608060405234801561001057600080fd5b50600436106100df5760003560e01c8063485cc9551161008c578063706da1ca11610066578063706da1ca146101e3578063e89e106a14610228578063f08c5daa14610231578063f6eaffc81461023a57600080fd5b8063485cc9551461019d57806355380dfb146101b05780636802f726146101d057600080fd5b80632fa4e442116100bd5780632fa4e4421461013257806336bfffed146101455780633b2bcbf11461015857600080fd5b80631fe543e3146100e457806321cd1fc1146100f957806327784fad1461010c575b600080fd5b6100f76100f2366004610f76565b61024d565b005b6100f7610107366004610dc0565b6102de565b61011f61011a366004610edb565b6104ad565b6040519081526020015b60405180910390f35b6100f7610140366004611037565b61058f565b6100f7610153366004610e15565b6106ef565b6003546101789073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610129565b6100f76101ab366004610de2565b610877565b6004546101789073ffffffffffffffffffffffffffffffffffffffff1681565b6100f76101de366004611037565b610a61565b60045461020f9074010000000000000000000000000000000000000000900467ffffffffffffffff1681565b60405167ffffffffffffffff9091168152602001610129565b61011f60025481565b61011f60055481565b61011f610248366004610f44565b610c6b565b60005462010000900473ffffffffffffffffffffffffffffffffffffffff1633146102d0576000546040517f1cf993f40000000000000000000000000000000000000000000000000000000081523360048201526201000090910473ffffffffffffffffffffffffffffffffffffffff1660248201526044015b60405180910390fd5b6102da8282610c8c565b5050565b600054610100900460ff16158080156102fe5750600054600160ff909116105b806103185750303b158015610318575060005460ff166001145b6103a4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a656400000000000000000000000000000000000060648201526084016102c7565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055801561040257600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b600080547fffffffffffffffffffff0000000000000000000000000000000000000000ffff166201000073ffffffffffffffffffffffffffffffffffffffff85160217905580156102da57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15050565b6003546040517f5d3b1d300000000000000000000000000000000000000000000000000000000081526004810187905267ffffffffffffffff8616602482015261ffff8516604482015263ffffffff80851660648301528316608482015260009173ffffffffffffffffffffffffffffffffffffffff1690635d3b1d309060a401602060405180830381600087803b15801561054857600080fd5b505af115801561055c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105809190610f5d565b60028190559695505050505050565b60045474010000000000000000000000000000000000000000900467ffffffffffffffff1661061a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f737562206e6f742073657400000000000000000000000000000000000000000060448201526064016102c7565b6004546003546040805174010000000000000000000000000000000000000000840467ffffffffffffffff16602082015273ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918591015b6040516020818303038152906040526040518463ffffffff1660e01b815260040161069d93929190611065565b602060405180830381600087803b1580156106b757600080fd5b505af11580156106cb573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102da9190610eb9565b60045474010000000000000000000000000000000000000000900467ffffffffffffffff1661077a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f7375624944206e6f74207365740000000000000000000000000000000000000060448201526064016102c7565b60005b81518110156102da57600354600454835173ffffffffffffffffffffffffffffffffffffffff90921691637341c10c9174010000000000000000000000000000000000000000900467ffffffffffffffff16908590859081106107e2576107e26111de565b60200260200101516040518363ffffffff1660e01b815260040161083292919067ffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b15801561084c57600080fd5b505af1158015610860573d6000803e3d6000fd5b50505050808061086f9061117e565b91505061077d565b600054610100900460ff16158080156108975750600054600160ff909116105b806108b15750303b1580156108b1575060005460ff166001145b61093d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a656400000000000000000000000000000000000060648201526084016102c7565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055801561099b57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b6109a4836102de565b6003805473ffffffffffffffffffffffffffffffffffffffff8086167fffffffffffffffffffffffff00000000000000000000000000000000000000009283161790925560048054928516929091169190911790558015610a5c57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b505050565b60045474010000000000000000000000000000000000000000900467ffffffffffffffff1661061a57600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a21a23e46040518163ffffffff1660e01b8152600401602060405180830381600087803b158015610af457600080fd5b505af1158015610b08573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b2c919061101a565b600480547fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000067ffffffffffffffff9384168102919091178083556003546040517f7341c10c000000000000000000000000000000000000000000000000000000008152929091049093169181019190915230602482015273ffffffffffffffffffffffffffffffffffffffff90911690637341c10c90604401600060405180830381600087803b158015610bf657600080fd5b505af1158015610c0a573d6000803e3d6000fd5b50506004546003546040805174010000000000000000000000000000000000000000840467ffffffffffffffff16602082015273ffffffffffffffffffffffffffffffffffffffff9384169550634000aea094509290911691859101610670565b60018181548110610c7b57600080fd5b600091825260209091200154905081565b6002548214610cf7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265717565737420494420697320696e636f727265637400000000000000000060448201526064016102c7565b5a6005558051600180548282556000829052610a5c927fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf691820191602086018215610d5e579160200282015b82811115610d5e578251825591602001919060010190610d43565b50610d6a929150610d6e565b5090565b5b80821115610d6a5760008155600101610d6f565b803573ffffffffffffffffffffffffffffffffffffffff81168114610da757600080fd5b919050565b803563ffffffff81168114610da757600080fd5b600060208284031215610dd257600080fd5b610ddb82610d83565b9392505050565b60008060408385031215610df557600080fd5b610dfe83610d83565b9150610e0c60208401610d83565b90509250929050565b60006020808385031215610e2857600080fd5b823567ffffffffffffffff811115610e3f57600080fd5b8301601f81018513610e5057600080fd5b8035610e63610e5e8261115a565b61110b565b80828252848201915084840188868560051b8701011115610e8357600080fd5b600094505b83851015610ead57610e9981610d83565b835260019490940193918501918501610e88565b50979650505050505050565b600060208284031215610ecb57600080fd5b81518015158114610ddb57600080fd5b600080600080600060a08688031215610ef357600080fd5b853594506020860135610f058161123c565b9350604086013561ffff81168114610f1c57600080fd5b9250610f2a60608701610dac565b9150610f3860808701610dac565b90509295509295909350565b600060208284031215610f5657600080fd5b5035919050565b600060208284031215610f6f57600080fd5b5051919050565b60008060408385031215610f8957600080fd5b8235915060208084013567ffffffffffffffff811115610fa857600080fd5b8401601f81018613610fb957600080fd5b8035610fc7610e5e8261115a565b80828252848201915084840189868560051b8701011115610fe757600080fd5b600094505b8385101561100a578035835260019490940193918501918501610fec565b5080955050505050509250929050565b60006020828403121561102c57600080fd5b8151610ddb8161123c565b60006020828403121561104957600080fd5b81356bffffffffffffffffffffffff81168114610ddb57600080fd5b73ffffffffffffffffffffffffffffffffffffffff84168152600060206bffffffffffffffffffffffff85168184015260606040840152835180606085015260005b818110156110c3578581018301518582016080015282016110a7565b818111156110d5576000608083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160800195945050505050565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156111525761115261120d565b604052919050565b600067ffffffffffffffff8211156111745761117461120d565b5060051b60200190565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156111d7577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b67ffffffffffffffff8116811461125257600080fd5b5056fea164736f6c6343000806000a",
}

var VRFConsumerV2UpgradeableABI = VRFConsumerV2UpgradeableMetaData.ABI

var VRFConsumerV2UpgradeableBin = VRFConsumerV2UpgradeableMetaData.Bin

func DeployVRFConsumerV2Upgradeable(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *VRFConsumerV2Upgradeable, error) {
	parsed, err := VRFConsumerV2UpgradeableMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFConsumerV2UpgradeableBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFConsumerV2Upgradeable{VRFConsumerV2UpgradeableCaller: VRFConsumerV2UpgradeableCaller{contract: contract}, VRFConsumerV2UpgradeableTransactor: VRFConsumerV2UpgradeableTransactor{contract: contract}, VRFConsumerV2UpgradeableFilterer: VRFConsumerV2UpgradeableFilterer{contract: contract}}, nil
}

type VRFConsumerV2Upgradeable struct {
	address common.Address
	abi     abi.ABI
	VRFConsumerV2UpgradeableCaller
	VRFConsumerV2UpgradeableTransactor
	VRFConsumerV2UpgradeableFilterer
}

type VRFConsumerV2UpgradeableCaller struct {
	contract *bind.BoundContract
}

type VRFConsumerV2UpgradeableTransactor struct {
	contract *bind.BoundContract
}

type VRFConsumerV2UpgradeableFilterer struct {
	contract *bind.BoundContract
}

type VRFConsumerV2UpgradeableSession struct {
	Contract     *VRFConsumerV2Upgradeable
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFConsumerV2UpgradeableCallerSession struct {
	Contract *VRFConsumerV2UpgradeableCaller
	CallOpts bind.CallOpts
}

type VRFConsumerV2UpgradeableTransactorSession struct {
	Contract     *VRFConsumerV2UpgradeableTransactor
	TransactOpts bind.TransactOpts
}

type VRFConsumerV2UpgradeableRaw struct {
	Contract *VRFConsumerV2Upgradeable
}

type VRFConsumerV2UpgradeableCallerRaw struct {
	Contract *VRFConsumerV2UpgradeableCaller
}

type VRFConsumerV2UpgradeableTransactorRaw struct {
	Contract *VRFConsumerV2UpgradeableTransactor
}

func NewVRFConsumerV2Upgradeable(address common.Address, backend bind.ContractBackend) (*VRFConsumerV2Upgradeable, error) {
	abi, err := abi.JSON(strings.NewReader(VRFConsumerV2UpgradeableABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFConsumerV2Upgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2Upgradeable{address: address, abi: abi, VRFConsumerV2UpgradeableCaller: VRFConsumerV2UpgradeableCaller{contract: contract}, VRFConsumerV2UpgradeableTransactor: VRFConsumerV2UpgradeableTransactor{contract: contract}, VRFConsumerV2UpgradeableFilterer: VRFConsumerV2UpgradeableFilterer{contract: contract}}, nil
}

func NewVRFConsumerV2UpgradeableCaller(address common.Address, caller bind.ContractCaller) (*VRFConsumerV2UpgradeableCaller, error) {
	contract, err := bindVRFConsumerV2Upgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2UpgradeableCaller{contract: contract}, nil
}

func NewVRFConsumerV2UpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFConsumerV2UpgradeableTransactor, error) {
	contract, err := bindVRFConsumerV2Upgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2UpgradeableTransactor{contract: contract}, nil
}

func NewVRFConsumerV2UpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFConsumerV2UpgradeableFilterer, error) {
	contract, err := bindVRFConsumerV2Upgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2UpgradeableFilterer{contract: contract}, nil
}

func bindVRFConsumerV2Upgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFConsumerV2UpgradeableABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFConsumerV2Upgradeable.Contract.VRFConsumerV2UpgradeableCaller.contract.Call(opts, result, method, params...)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFConsumerV2Upgradeable.Contract.VRFConsumerV2UpgradeableTransactor.contract.Transfer(opts)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFConsumerV2Upgradeable.Contract.VRFConsumerV2UpgradeableTransactor.contract.Transact(opts, method, params...)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFConsumerV2Upgradeable.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFConsumerV2Upgradeable.Contract.contract.Transfer(opts)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFConsumerV2Upgradeable.Contract.contract.Transact(opts, method, params...)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableCaller) COORDINATOR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFConsumerV2Upgradeable.contract.Call(opts, &out, "COORDINATOR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableSession) COORDINATOR() (common.Address, error) {
	return _VRFConsumerV2Upgradeable.Contract.COORDINATOR(&_VRFConsumerV2Upgradeable.CallOpts)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableCallerSession) COORDINATOR() (common.Address, error) {
	return _VRFConsumerV2Upgradeable.Contract.COORDINATOR(&_VRFConsumerV2Upgradeable.CallOpts)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableCaller) LINKTOKEN(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFConsumerV2Upgradeable.contract.Call(opts, &out, "LINKTOKEN")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableSession) LINKTOKEN() (common.Address, error) {
	return _VRFConsumerV2Upgradeable.Contract.LINKTOKEN(&_VRFConsumerV2Upgradeable.CallOpts)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableCallerSession) LINKTOKEN() (common.Address, error) {
	return _VRFConsumerV2Upgradeable.Contract.LINKTOKEN(&_VRFConsumerV2Upgradeable.CallOpts)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableCaller) SGasAvailable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFConsumerV2Upgradeable.contract.Call(opts, &out, "s_gasAvailable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableSession) SGasAvailable() (*big.Int, error) {
	return _VRFConsumerV2Upgradeable.Contract.SGasAvailable(&_VRFConsumerV2Upgradeable.CallOpts)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableCallerSession) SGasAvailable() (*big.Int, error) {
	return _VRFConsumerV2Upgradeable.Contract.SGasAvailable(&_VRFConsumerV2Upgradeable.CallOpts)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableCaller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFConsumerV2Upgradeable.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFConsumerV2Upgradeable.Contract.SRandomWords(&_VRFConsumerV2Upgradeable.CallOpts, arg0)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableCallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFConsumerV2Upgradeable.Contract.SRandomWords(&_VRFConsumerV2Upgradeable.CallOpts, arg0)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableCaller) SRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFConsumerV2Upgradeable.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableSession) SRequestId() (*big.Int, error) {
	return _VRFConsumerV2Upgradeable.Contract.SRequestId(&_VRFConsumerV2Upgradeable.CallOpts)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableCallerSession) SRequestId() (*big.Int, error) {
	return _VRFConsumerV2Upgradeable.Contract.SRequestId(&_VRFConsumerV2Upgradeable.CallOpts)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableCaller) SSubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFConsumerV2Upgradeable.contract.Call(opts, &out, "s_subId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableSession) SSubId() (uint64, error) {
	return _VRFConsumerV2Upgradeable.Contract.SSubId(&_VRFConsumerV2Upgradeable.CallOpts)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableCallerSession) SSubId() (uint64, error) {
	return _VRFConsumerV2Upgradeable.Contract.SSubId(&_VRFConsumerV2Upgradeable.CallOpts)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableTransactor) VRFConsumerBaseV2Init(opts *bind.TransactOpts, _vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2Upgradeable.contract.Transact(opts, "__VRFConsumerBaseV2_init", _vrfCoordinator)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableSession) VRFConsumerBaseV2Init(_vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2Upgradeable.Contract.VRFConsumerBaseV2Init(&_VRFConsumerV2Upgradeable.TransactOpts, _vrfCoordinator)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableTransactorSession) VRFConsumerBaseV2Init(_vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2Upgradeable.Contract.VRFConsumerBaseV2Init(&_VRFConsumerV2Upgradeable.TransactOpts, _vrfCoordinator)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableTransactor) Initialize(opts *bind.TransactOpts, _vrfCoordinator common.Address, _link common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2Upgradeable.contract.Transact(opts, "initialize", _vrfCoordinator, _link)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableSession) Initialize(_vrfCoordinator common.Address, _link common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2Upgradeable.Contract.Initialize(&_VRFConsumerV2Upgradeable.TransactOpts, _vrfCoordinator, _link)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableTransactorSession) Initialize(_vrfCoordinator common.Address, _link common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2Upgradeable.Contract.Initialize(&_VRFConsumerV2Upgradeable.TransactOpts, _vrfCoordinator, _link)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2Upgradeable.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2Upgradeable.Contract.RawFulfillRandomWords(&_VRFConsumerV2Upgradeable.TransactOpts, requestId, randomWords)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2Upgradeable.Contract.RawFulfillRandomWords(&_VRFConsumerV2Upgradeable.TransactOpts, requestId, randomWords)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableTransactor) TestCreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2Upgradeable.contract.Transact(opts, "testCreateSubscriptionAndFund", amount)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableSession) TestCreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2Upgradeable.Contract.TestCreateSubscriptionAndFund(&_VRFConsumerV2Upgradeable.TransactOpts, amount)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableTransactorSession) TestCreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2Upgradeable.Contract.TestCreateSubscriptionAndFund(&_VRFConsumerV2Upgradeable.TransactOpts, amount)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableTransactor) TestRequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFConsumerV2Upgradeable.contract.Transact(opts, "testRequestRandomness", keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableSession) TestRequestRandomness(keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFConsumerV2Upgradeable.Contract.TestRequestRandomness(&_VRFConsumerV2Upgradeable.TransactOpts, keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableTransactorSession) TestRequestRandomness(keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFConsumerV2Upgradeable.Contract.TestRequestRandomness(&_VRFConsumerV2Upgradeable.TransactOpts, keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableTransactor) TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2Upgradeable.contract.Transact(opts, "topUpSubscription", amount)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2Upgradeable.Contract.TopUpSubscription(&_VRFConsumerV2Upgradeable.TransactOpts, amount)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableTransactorSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFConsumerV2Upgradeable.Contract.TopUpSubscription(&_VRFConsumerV2Upgradeable.TransactOpts, amount)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableTransactor) UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2Upgradeable.contract.Transact(opts, "updateSubscription", consumers)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2Upgradeable.Contract.UpdateSubscription(&_VRFConsumerV2Upgradeable.TransactOpts, consumers)
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableTransactorSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFConsumerV2Upgradeable.Contract.UpdateSubscription(&_VRFConsumerV2Upgradeable.TransactOpts, consumers)
}

type VRFConsumerV2UpgradeableInitializedIterator struct {
	Event *VRFConsumerV2UpgradeableInitialized

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFConsumerV2UpgradeableInitializedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFConsumerV2UpgradeableInitialized)
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
		it.Event = new(VRFConsumerV2UpgradeableInitialized)
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

func (it *VRFConsumerV2UpgradeableInitializedIterator) Error() error {
	return it.fail
}

func (it *VRFConsumerV2UpgradeableInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFConsumerV2UpgradeableInitialized struct {
	Version uint8
	Raw     types.Log
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableFilterer) FilterInitialized(opts *bind.FilterOpts) (*VRFConsumerV2UpgradeableInitializedIterator, error) {

	logs, sub, err := _VRFConsumerV2Upgradeable.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &VRFConsumerV2UpgradeableInitializedIterator{contract: _VRFConsumerV2Upgradeable.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *VRFConsumerV2UpgradeableInitialized) (event.Subscription, error) {

	logs, sub, err := _VRFConsumerV2Upgradeable.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFConsumerV2UpgradeableInitialized)
				if err := _VRFConsumerV2Upgradeable.contract.UnpackLog(event, "Initialized", log); err != nil {
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

func (_VRFConsumerV2Upgradeable *VRFConsumerV2UpgradeableFilterer) ParseInitialized(log types.Log) (*VRFConsumerV2UpgradeableInitialized, error) {
	event := new(VRFConsumerV2UpgradeableInitialized)
	if err := _VRFConsumerV2Upgradeable.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2Upgradeable) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFConsumerV2Upgradeable.abi.Events["Initialized"].ID:
		return _VRFConsumerV2Upgradeable.ParseInitialized(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFConsumerV2UpgradeableInitialized) Topic() common.Hash {
	return common.HexToHash("0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498")
}

func (_VRFConsumerV2Upgradeable *VRFConsumerV2Upgradeable) Address() common.Address {
	return _VRFConsumerV2Upgradeable.address
}

type VRFConsumerV2UpgradeableInterface interface {
	COORDINATOR(opts *bind.CallOpts) (common.Address, error)

	LINKTOKEN(opts *bind.CallOpts) (common.Address, error)

	SGasAvailable(opts *bind.CallOpts) (*big.Int, error)

	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestId(opts *bind.CallOpts) (*big.Int, error)

	SSubId(opts *bind.CallOpts) (uint64, error)

	VRFConsumerBaseV2Init(opts *bind.TransactOpts, _vrfCoordinator common.Address) (*types.Transaction, error)

	Initialize(opts *bind.TransactOpts, _vrfCoordinator common.Address, _link common.Address) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	TestCreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	TestRequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error)

	TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error)

	FilterInitialized(opts *bind.FilterOpts) (*VRFConsumerV2UpgradeableInitializedIterator, error)

	WatchInitialized(opts *bind.WatchOpts, sink chan<- *VRFConsumerV2UpgradeableInitialized) (event.Subscription, error)

	ParseInitialized(log types.Log) (*VRFConsumerV2UpgradeableInitialized, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
