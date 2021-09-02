// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_malicious_consumer_v2

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

var VRFMaliciousConsumerV2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_gasAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_subId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"setKeyHash\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"testCreateSubscriptionAndFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRequestRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"name\":\"updateSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b50604051610d89380380610d8983398101604081905261002f9161008e565b6001600160601b0319606083901b16608052600280546001600160a01b03199081166001600160a01b0394851617909155600380549290931691161790556100c1565b80516001600160a01b038116811461008957600080fd5b919050565b600080604083850312156100a157600080fd5b6100aa83610072565b91506100b860208401610072565b90509250929050565b60805160601c610ca36100e66000396000818161019301526101fb0152610ca36000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c8063706da1ca11610076578063e89e106a1161005b578063e89e106a14610156578063f08c5daa1461015f578063f6eaffc81461016857600080fd5b8063706da1ca146100fe578063985447101461014357600080fd5b80631fe543e3146100a857806336bfffed146100bd5780636802f726146100d0578063702a2ec9146100e3575b600080fd5b6100bb6100b6366004610969565b61017b565b005b6100bb6100cb36600461084f565b61023b565b6100bb6100de366004610a37565b6103c3565b6100eb6105d5565b6040519081526020015b60405180910390f35b60035461012a9074010000000000000000000000000000000000000000900467ffffffffffffffff1681565b60405167ffffffffffffffff90911681526020016100f5565b6100bb610151366004610937565b600555565b6100eb60015481565b6100eb60045481565b6100eb610176366004610937565b6106c4565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461022d576040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660248201526044015b60405180910390fd5b61023782826106e5565b5050565b60035474010000000000000000000000000000000000000000900467ffffffffffffffff166102c6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f7375624944206e6f7420736574000000000000000000000000000000000000006044820152606401610224565b60005b815181101561023757600254600354835173ffffffffffffffffffffffffffffffffffffffff90921691637341c10c9174010000000000000000000000000000000000000000900467ffffffffffffffff169085908590811061032e5761032e610c38565b60200260200101516040518363ffffffff1660e01b815260040161037e92919067ffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b15801561039857600080fd5b505af11580156103ac573d6000803e3d6000fd5b5050505080806103bb90610bd8565b9150506102c9565b60035474010000000000000000000000000000000000000000900467ffffffffffffffff166105015760408051600180825281830190925260009160208083019080368337019050509050308160008151811061042257610422610c38565b73ffffffffffffffffffffffffffffffffffffffff92831660209182029290920101526002546040517f6b9f7d38000000000000000000000000000000000000000000000000000000008152911690636b9f7d3890610485908490600401610b0b565b602060405180830381600087803b15801561049f57600080fd5b505af11580156104b3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104d79190610a0d565b600360146101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550505b6003546002546040805174010000000000000000000000000000000000000000840467ffffffffffffffff16602082015273ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918591016040516020818303038152906040526040518463ffffffff1660e01b815260040161058393929190610a65565b602060405180830381600087803b15801561059d57600080fd5b505af11580156105b1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610237919061090e565b6002546005546003546040517f5d3b1d30000000000000000000000000000000000000000000000000000000008152600481019290925274010000000000000000000000000000000000000000900467ffffffffffffffff1660248201526001604482018190526207a1206064830152608482015260009173ffffffffffffffffffffffffffffffffffffffff1690635d3b1d309060a401602060405180830381600087803b15801561068757600080fd5b505af115801561069b573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106bf9190610950565b905090565b600081815481106106d457600080fd5b600091825260209091200154905081565b5a60045580516106fc9060009060208401906107ef565b5060018281556002546005546003546040517f5d3b1d30000000000000000000000000000000000000000000000000000000008152600481019290925274010000000000000000000000000000000000000000900467ffffffffffffffff1660248201526044810183905262030d406064820152608481019290925273ffffffffffffffffffffffffffffffffffffffff1690635d3b1d309060a401602060405180830381600087803b1580156107b257600080fd5b505af11580156107c6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107ea9190610950565b505050565b82805482825590600052602060002090810192821561082a579160200282015b8281111561082a57825182559160200191906001019061080f565b5061083692915061083a565b5090565b5b80821115610836576000815560010161083b565b6000602080838503121561086257600080fd5b823567ffffffffffffffff81111561087957600080fd5b8301601f8101851361088a57600080fd5b803561089d61089882610bb4565b610b65565b80828252848201915084840188868560051b87010111156108bd57600080fd5b60009450845b8481101561090057813573ffffffffffffffffffffffffffffffffffffffff811681146108ee578687fd5b845292860192908601906001016108c3565b509098975050505050505050565b60006020828403121561092057600080fd5b8151801515811461093057600080fd5b9392505050565b60006020828403121561094957600080fd5b5035919050565b60006020828403121561096257600080fd5b5051919050565b6000806040838503121561097c57600080fd5b8235915060208084013567ffffffffffffffff81111561099b57600080fd5b8401601f810186136109ac57600080fd5b80356109ba61089882610bb4565b80828252848201915084840189868560051b87010111156109da57600080fd5b600094505b838510156109fd5780358352600194909401939185019185016109df565b5080955050505050509250929050565b600060208284031215610a1f57600080fd5b815167ffffffffffffffff8116811461093057600080fd5b600060208284031215610a4957600080fd5b81356bffffffffffffffffffffffff8116811461093057600080fd5b73ffffffffffffffffffffffffffffffffffffffff84168152600060206bffffffffffffffffffffffff85168184015260606040840152835180606085015260005b81811015610ac357858101830151858201608001528201610aa7565b81811115610ad5576000608083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160800195945050505050565b6020808252825182820181905260009190848201906040850190845b81811015610b5957835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101610b27565b50909695505050505050565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610bac57610bac610c67565b604052919050565b600067ffffffffffffffff821115610bce57610bce610c67565b5060051b60200190565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415610c31577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFMaliciousConsumerV2ABI = VRFMaliciousConsumerV2MetaData.ABI

var VRFMaliciousConsumerV2Bin = VRFMaliciousConsumerV2MetaData.Bin

func DeployVRFMaliciousConsumerV2(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, link common.Address) (common.Address, *types.Transaction, *VRFMaliciousConsumerV2, error) {
	parsed, err := VRFMaliciousConsumerV2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFMaliciousConsumerV2Bin), backend, vrfCoordinator, link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFMaliciousConsumerV2{VRFMaliciousConsumerV2Caller: VRFMaliciousConsumerV2Caller{contract: contract}, VRFMaliciousConsumerV2Transactor: VRFMaliciousConsumerV2Transactor{contract: contract}, VRFMaliciousConsumerV2Filterer: VRFMaliciousConsumerV2Filterer{contract: contract}}, nil
}

type VRFMaliciousConsumerV2 struct {
	address common.Address
	abi     abi.ABI
	VRFMaliciousConsumerV2Caller
	VRFMaliciousConsumerV2Transactor
	VRFMaliciousConsumerV2Filterer
}

type VRFMaliciousConsumerV2Caller struct {
	contract *bind.BoundContract
}

type VRFMaliciousConsumerV2Transactor struct {
	contract *bind.BoundContract
}

type VRFMaliciousConsumerV2Filterer struct {
	contract *bind.BoundContract
}

type VRFMaliciousConsumerV2Session struct {
	Contract     *VRFMaliciousConsumerV2
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFMaliciousConsumerV2CallerSession struct {
	Contract *VRFMaliciousConsumerV2Caller
	CallOpts bind.CallOpts
}

type VRFMaliciousConsumerV2TransactorSession struct {
	Contract     *VRFMaliciousConsumerV2Transactor
	TransactOpts bind.TransactOpts
}

type VRFMaliciousConsumerV2Raw struct {
	Contract *VRFMaliciousConsumerV2
}

type VRFMaliciousConsumerV2CallerRaw struct {
	Contract *VRFMaliciousConsumerV2Caller
}

type VRFMaliciousConsumerV2TransactorRaw struct {
	Contract *VRFMaliciousConsumerV2Transactor
}

func NewVRFMaliciousConsumerV2(address common.Address, backend bind.ContractBackend) (*VRFMaliciousConsumerV2, error) {
	abi, err := abi.JSON(strings.NewReader(VRFMaliciousConsumerV2ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFMaliciousConsumerV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFMaliciousConsumerV2{address: address, abi: abi, VRFMaliciousConsumerV2Caller: VRFMaliciousConsumerV2Caller{contract: contract}, VRFMaliciousConsumerV2Transactor: VRFMaliciousConsumerV2Transactor{contract: contract}, VRFMaliciousConsumerV2Filterer: VRFMaliciousConsumerV2Filterer{contract: contract}}, nil
}

func NewVRFMaliciousConsumerV2Caller(address common.Address, caller bind.ContractCaller) (*VRFMaliciousConsumerV2Caller, error) {
	contract, err := bindVRFMaliciousConsumerV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFMaliciousConsumerV2Caller{contract: contract}, nil
}

func NewVRFMaliciousConsumerV2Transactor(address common.Address, transactor bind.ContractTransactor) (*VRFMaliciousConsumerV2Transactor, error) {
	contract, err := bindVRFMaliciousConsumerV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFMaliciousConsumerV2Transactor{contract: contract}, nil
}

func NewVRFMaliciousConsumerV2Filterer(address common.Address, filterer bind.ContractFilterer) (*VRFMaliciousConsumerV2Filterer, error) {
	contract, err := bindVRFMaliciousConsumerV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFMaliciousConsumerV2Filterer{contract: contract}, nil
}

func bindVRFMaliciousConsumerV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFMaliciousConsumerV2ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFMaliciousConsumerV2.Contract.VRFMaliciousConsumerV2Caller.contract.Call(opts, result, method, params...)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.VRFMaliciousConsumerV2Transactor.contract.Transfer(opts)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.VRFMaliciousConsumerV2Transactor.contract.Transact(opts, method, params...)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFMaliciousConsumerV2.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.contract.Transfer(opts)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.contract.Transact(opts, method, params...)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Caller) SGasAvailable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFMaliciousConsumerV2.contract.Call(opts, &out, "s_gasAvailable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Session) SGasAvailable() (*big.Int, error) {
	return _VRFMaliciousConsumerV2.Contract.SGasAvailable(&_VRFMaliciousConsumerV2.CallOpts)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2CallerSession) SGasAvailable() (*big.Int, error) {
	return _VRFMaliciousConsumerV2.Contract.SGasAvailable(&_VRFMaliciousConsumerV2.CallOpts)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Caller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFMaliciousConsumerV2.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Session) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFMaliciousConsumerV2.Contract.SRandomWords(&_VRFMaliciousConsumerV2.CallOpts, arg0)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2CallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFMaliciousConsumerV2.Contract.SRandomWords(&_VRFMaliciousConsumerV2.CallOpts, arg0)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Caller) SRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFMaliciousConsumerV2.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Session) SRequestId() (*big.Int, error) {
	return _VRFMaliciousConsumerV2.Contract.SRequestId(&_VRFMaliciousConsumerV2.CallOpts)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2CallerSession) SRequestId() (*big.Int, error) {
	return _VRFMaliciousConsumerV2.Contract.SRequestId(&_VRFMaliciousConsumerV2.CallOpts)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Caller) SSubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFMaliciousConsumerV2.contract.Call(opts, &out, "s_subId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Session) SSubId() (uint64, error) {
	return _VRFMaliciousConsumerV2.Contract.SSubId(&_VRFMaliciousConsumerV2.CallOpts)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2CallerSession) SSubId() (uint64, error) {
	return _VRFMaliciousConsumerV2.Contract.SSubId(&_VRFMaliciousConsumerV2.CallOpts)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Transactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Session) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.RawFulfillRandomWords(&_VRFMaliciousConsumerV2.TransactOpts, requestId, randomWords)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2TransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.RawFulfillRandomWords(&_VRFMaliciousConsumerV2.TransactOpts, requestId, randomWords)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Transactor) SetKeyHash(opts *bind.TransactOpts, keyHash [32]byte) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.contract.Transact(opts, "setKeyHash", keyHash)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Session) SetKeyHash(keyHash [32]byte) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.SetKeyHash(&_VRFMaliciousConsumerV2.TransactOpts, keyHash)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2TransactorSession) SetKeyHash(keyHash [32]byte) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.SetKeyHash(&_VRFMaliciousConsumerV2.TransactOpts, keyHash)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Transactor) TestCreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.contract.Transact(opts, "testCreateSubscriptionAndFund", amount)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Session) TestCreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.TestCreateSubscriptionAndFund(&_VRFMaliciousConsumerV2.TransactOpts, amount)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2TransactorSession) TestCreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.TestCreateSubscriptionAndFund(&_VRFMaliciousConsumerV2.TransactOpts, amount)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Transactor) TestRequestRandomness(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.contract.Transact(opts, "testRequestRandomness")
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Session) TestRequestRandomness() (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.TestRequestRandomness(&_VRFMaliciousConsumerV2.TransactOpts)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2TransactorSession) TestRequestRandomness() (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.TestRequestRandomness(&_VRFMaliciousConsumerV2.TransactOpts)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Transactor) UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.contract.Transact(opts, "updateSubscription", consumers)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2Session) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.UpdateSubscription(&_VRFMaliciousConsumerV2.TransactOpts, consumers)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2TransactorSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFMaliciousConsumerV2.Contract.UpdateSubscription(&_VRFMaliciousConsumerV2.TransactOpts, consumers)
}

func (_VRFMaliciousConsumerV2 *VRFMaliciousConsumerV2) Address() common.Address {
	return _VRFMaliciousConsumerV2.address
}

type VRFMaliciousConsumerV2Interface interface {
	SGasAvailable(opts *bind.CallOpts) (*big.Int, error)

	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestId(opts *bind.CallOpts) (*big.Int, error)

	SSubId(opts *bind.CallOpts) (uint64, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	SetKeyHash(opts *bind.TransactOpts, keyHash [32]byte) (*types.Transaction, error)

	TestCreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	TestRequestRandomness(opts *bind.TransactOpts) (*types.Transaction, error)

	UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error)

	Address() common.Address
}
