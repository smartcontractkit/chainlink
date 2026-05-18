// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_migratable_consumer_v2

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

var MigratableVRFConsumerV2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint16\",\"name\":\"minReqConfs\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"nativePayment\",\"type\":\"bool\"}],\"name\":\"requestRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"}],\"name\":\"setCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"setSubId\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50604051610ed2380380610ed283398101604081905261002f916101b7565b818133806000816100875760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100b7576100b78161010d565b5050600280546305d3b1d360e41b6001600160e01b036001600160401b03909516600160a01b026001600160e01b03199092166001600160a01b0390961695909517179290921692909217905550610209915050565b6001600160a01b0381163314156101665760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161007e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600080604083850312156101ca57600080fd5b82516001600160a01b03811681146101e157600080fd5b60208401519092506001600160401b03811681146101fe57600080fd5b809150509250929050565b610cba806102186000396000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c80638da5cb5b11610076578063e826149d1161005b578063e826149d14610139578063e89e106a1461014c578063f2fde38b1461015557600080fd5b80638da5cb5b146100fe5780638ea981171461012657600080fd5b80631fe543e3146100a85780633e2831fe146100bd57806379ba5097146100d05780637b95ba02146100d8575b600080fd5b6100bb6100b6366004610ac3565b610168565b005b6100bb6100cb366004610bd4565b6101ee565b6100bb610249565b6100eb6100e6366004610bb2565b610346565b6040519081526020015b60405180910390f35b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100f5565b6100bb610134366004610a16565b610377565b6100eb610147366004610a53565b6105ae565b6100eb60045481565b6100bb610163366004610a16565b61077e565b60025473ffffffffffffffffffffffffffffffffffffffff1633146101e0576002546040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff90911660248201526044015b60405180910390fd5b6101ea8282610792565b5050565b6101f6610824565b6002805467ffffffffffffffff90921674010000000000000000000000000000000000000000027fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff909216919091179055565b60015473ffffffffffffffffffffffffffffffffffffffff1633146102ca576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016101d7565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6003602052816000526040600020818154811061036257600080fd5b90600052602060002001600091509150505481565b61037f610824565b604080517f181f5a7700000000000000000000000000000000000000000000000000000000602082015260009101604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815290829052600254909250600091829173ffffffffffffffffffffffffffffffffffffffff1690610409908590610bfe565b6000604051808303816000865af19150503d8060008114610446576040519150601f19603f3d011682016040523d82523d6000602084013e61044b565b606091505b5091509150816104b7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f74797065416e6456657273696f6e206661696c6564000000000000000000000060448201526064016101d7565b6040516020016104f8906020808252601a908201527f565246436f6f7264696e61746f725632506c757320312e302e30000000000000604082015260600190565b604051602081830303815290604052805190602001208180519060200120141561056357600280547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fefcf1d94000000000000000000000000000000000000000000000000000000001790555b5050600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff939093169290921790915550565b600254604080516024810188905274010000000000000000000000000000000000000000830467ffffffffffffffff16604482015261ffff8716606482015263ffffffff8681166084830152851660a482015283151560c4808301919091528251808303909101815260e490910182526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167c0100000000000000000000000000000000000000000000000000000000850460e01b7fffffffff000000000000000000000000000000000000000000000000000000001617905290516000928391829173ffffffffffffffffffffffffffffffffffffffff16906106b5908590610bfe565b6000604051808303816000865af19150503d80600081146106f2576040519150601f19603f3d011682016040523d82523d6000602084013e6106f7565b606091505b509150915081610763576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f726571756573742072616e646f6d20776f726473206661696c6564000000000060448201526064016101d7565b61076c81610c39565b60048190559998505050505050505050565b610786610824565b61078f816108a7565b50565b60045482146107fd576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265717565737420494420697320696e636f727265637400000000000000000060448201526064016101d7565b6004546000908152600360209081526040909120825161081f9284019061099d565b505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146108a5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016101d7565b565b73ffffffffffffffffffffffffffffffffffffffff8116331415610927576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016101d7565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b8280548282559060005260206000209081019282156109d8579160200282015b828111156109d85782518255916020019190600101906109bd565b506109e49291506109e8565b5090565b5b808211156109e457600081556001016109e9565b803563ffffffff81168114610a1157600080fd5b919050565b600060208284031215610a2857600080fd5b813573ffffffffffffffffffffffffffffffffffffffff81168114610a4c57600080fd5b9392505050565b600080600080600060a08688031215610a6b57600080fd5b85359450602086013561ffff81168114610a8457600080fd5b9350610a92604087016109fd565b9250610aa0606087016109fd565b915060808601358015158114610ab557600080fd5b809150509295509295909350565b60008060408385031215610ad657600080fd5b8235915060208084013567ffffffffffffffff80821115610af657600080fd5b818601915086601f830112610b0a57600080fd5b813581811115610b1c57610b1c610c7e565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f83011681018181108582111715610b5f57610b5f610c7e565b604052828152858101935084860182860187018b1015610b7e57600080fd5b600095505b83861015610ba1578035855260019590950194938601938601610b83565b508096505050505050509250929050565b60008060408385031215610bc557600080fd5b50508035926020909101359150565b600060208284031215610be657600080fd5b813567ffffffffffffffff81168114610a4c57600080fd5b6000825160005b81811015610c1f5760208186018101518583015201610c05565b81811115610c2e576000828501525b509190910192915050565b80516020808301519190811015610c78577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8160200360031b1b821691505b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var MigratableVRFConsumerV2ABI = MigratableVRFConsumerV2MetaData.ABI

var MigratableVRFConsumerV2Bin = MigratableVRFConsumerV2MetaData.Bin

func DeployMigratableVRFConsumerV2(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, subId uint64) (common.Address, *types.Transaction, *MigratableVRFConsumerV2, error) {
	parsed, err := MigratableVRFConsumerV2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MigratableVRFConsumerV2Bin), backend, vrfCoordinator, subId)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MigratableVRFConsumerV2{MigratableVRFConsumerV2Caller: MigratableVRFConsumerV2Caller{contract: contract}, MigratableVRFConsumerV2Transactor: MigratableVRFConsumerV2Transactor{contract: contract}, MigratableVRFConsumerV2Filterer: MigratableVRFConsumerV2Filterer{contract: contract}}, nil
}

type MigratableVRFConsumerV2 struct {
	address common.Address
	abi     abi.ABI
	MigratableVRFConsumerV2Caller
	MigratableVRFConsumerV2Transactor
	MigratableVRFConsumerV2Filterer
}

type MigratableVRFConsumerV2Caller struct {
	contract *bind.BoundContract
}

type MigratableVRFConsumerV2Transactor struct {
	contract *bind.BoundContract
}

type MigratableVRFConsumerV2Filterer struct {
	contract *bind.BoundContract
}

type MigratableVRFConsumerV2Session struct {
	Contract     *MigratableVRFConsumerV2
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MigratableVRFConsumerV2CallerSession struct {
	Contract *MigratableVRFConsumerV2Caller
	CallOpts bind.CallOpts
}

type MigratableVRFConsumerV2TransactorSession struct {
	Contract     *MigratableVRFConsumerV2Transactor
	TransactOpts bind.TransactOpts
}

type MigratableVRFConsumerV2Raw struct {
	Contract *MigratableVRFConsumerV2
}

type MigratableVRFConsumerV2CallerRaw struct {
	Contract *MigratableVRFConsumerV2Caller
}

type MigratableVRFConsumerV2TransactorRaw struct {
	Contract *MigratableVRFConsumerV2Transactor
}

func NewMigratableVRFConsumerV2(address common.Address, backend bind.ContractBackend) (*MigratableVRFConsumerV2, error) {
	abi, err := abi.JSON(strings.NewReader(MigratableVRFConsumerV2ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMigratableVRFConsumerV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MigratableVRFConsumerV2{address: address, abi: abi, MigratableVRFConsumerV2Caller: MigratableVRFConsumerV2Caller{contract: contract}, MigratableVRFConsumerV2Transactor: MigratableVRFConsumerV2Transactor{contract: contract}, MigratableVRFConsumerV2Filterer: MigratableVRFConsumerV2Filterer{contract: contract}}, nil
}

func NewMigratableVRFConsumerV2Caller(address common.Address, caller bind.ContractCaller) (*MigratableVRFConsumerV2Caller, error) {
	contract, err := bindMigratableVRFConsumerV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MigratableVRFConsumerV2Caller{contract: contract}, nil
}

func NewMigratableVRFConsumerV2Transactor(address common.Address, transactor bind.ContractTransactor) (*MigratableVRFConsumerV2Transactor, error) {
	contract, err := bindMigratableVRFConsumerV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MigratableVRFConsumerV2Transactor{contract: contract}, nil
}

func NewMigratableVRFConsumerV2Filterer(address common.Address, filterer bind.ContractFilterer) (*MigratableVRFConsumerV2Filterer, error) {
	contract, err := bindMigratableVRFConsumerV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MigratableVRFConsumerV2Filterer{contract: contract}, nil
}

func bindMigratableVRFConsumerV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MigratableVRFConsumerV2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MigratableVRFConsumerV2.Contract.MigratableVRFConsumerV2Caller.contract.Call(opts, result, method, params...)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MigratableVRFConsumerV2.Contract.MigratableVRFConsumerV2Transactor.contract.Transfer(opts)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MigratableVRFConsumerV2.Contract.MigratableVRFConsumerV2Transactor.contract.Transact(opts, method, params...)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MigratableVRFConsumerV2.Contract.contract.Call(opts, result, method, params...)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MigratableVRFConsumerV2.Contract.contract.Transfer(opts)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MigratableVRFConsumerV2.Contract.contract.Transact(opts, method, params...)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MigratableVRFConsumerV2.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Session) Owner() (common.Address, error) {
	return _MigratableVRFConsumerV2.Contract.Owner(&_MigratableVRFConsumerV2.CallOpts)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2CallerSession) Owner() (common.Address, error) {
	return _MigratableVRFConsumerV2.Contract.Owner(&_MigratableVRFConsumerV2.CallOpts)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Caller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _MigratableVRFConsumerV2.contract.Call(opts, &out, "s_randomWords", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Session) SRandomWords(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _MigratableVRFConsumerV2.Contract.SRandomWords(&_MigratableVRFConsumerV2.CallOpts, arg0, arg1)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2CallerSession) SRandomWords(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _MigratableVRFConsumerV2.Contract.SRandomWords(&_MigratableVRFConsumerV2.CallOpts, arg0, arg1)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Caller) SRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MigratableVRFConsumerV2.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Session) SRequestId() (*big.Int, error) {
	return _MigratableVRFConsumerV2.Contract.SRequestId(&_MigratableVRFConsumerV2.CallOpts)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2CallerSession) SRequestId() (*big.Int, error) {
	return _MigratableVRFConsumerV2.Contract.SRequestId(&_MigratableVRFConsumerV2.CallOpts)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MigratableVRFConsumerV2.contract.Transact(opts, "acceptOwnership")
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Session) AcceptOwnership() (*types.Transaction, error) {
	return _MigratableVRFConsumerV2.Contract.AcceptOwnership(&_MigratableVRFConsumerV2.TransactOpts)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _MigratableVRFConsumerV2.Contract.AcceptOwnership(&_MigratableVRFConsumerV2.TransactOpts)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Transactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _MigratableVRFConsumerV2.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Session) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _MigratableVRFConsumerV2.Contract.RawFulfillRandomWords(&_MigratableVRFConsumerV2.TransactOpts, requestId, randomWords)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2TransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _MigratableVRFConsumerV2.Contract.RawFulfillRandomWords(&_MigratableVRFConsumerV2.TransactOpts, requestId, randomWords)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Transactor) RequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, minReqConfs uint16, callbackGasLimit uint32, numWords uint32, nativePayment bool) (*types.Transaction, error) {
	return _MigratableVRFConsumerV2.contract.Transact(opts, "requestRandomness", keyHash, minReqConfs, callbackGasLimit, numWords, nativePayment)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Session) RequestRandomness(keyHash [32]byte, minReqConfs uint16, callbackGasLimit uint32, numWords uint32, nativePayment bool) (*types.Transaction, error) {
	return _MigratableVRFConsumerV2.Contract.RequestRandomness(&_MigratableVRFConsumerV2.TransactOpts, keyHash, minReqConfs, callbackGasLimit, numWords, nativePayment)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2TransactorSession) RequestRandomness(keyHash [32]byte, minReqConfs uint16, callbackGasLimit uint32, numWords uint32, nativePayment bool) (*types.Transaction, error) {
	return _MigratableVRFConsumerV2.Contract.RequestRandomness(&_MigratableVRFConsumerV2.TransactOpts, keyHash, minReqConfs, callbackGasLimit, numWords, nativePayment)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Transactor) SetCoordinator(opts *bind.TransactOpts, coordinator common.Address) (*types.Transaction, error) {
	return _MigratableVRFConsumerV2.contract.Transact(opts, "setCoordinator", coordinator)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Session) SetCoordinator(coordinator common.Address) (*types.Transaction, error) {
	return _MigratableVRFConsumerV2.Contract.SetCoordinator(&_MigratableVRFConsumerV2.TransactOpts, coordinator)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2TransactorSession) SetCoordinator(coordinator common.Address) (*types.Transaction, error) {
	return _MigratableVRFConsumerV2.Contract.SetCoordinator(&_MigratableVRFConsumerV2.TransactOpts, coordinator)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Transactor) SetSubId(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error) {
	return _MigratableVRFConsumerV2.contract.Transact(opts, "setSubId", subId)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Session) SetSubId(subId uint64) (*types.Transaction, error) {
	return _MigratableVRFConsumerV2.Contract.SetSubId(&_MigratableVRFConsumerV2.TransactOpts, subId)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2TransactorSession) SetSubId(subId uint64) (*types.Transaction, error) {
	return _MigratableVRFConsumerV2.Contract.SetSubId(&_MigratableVRFConsumerV2.TransactOpts, subId)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _MigratableVRFConsumerV2.contract.Transact(opts, "transferOwnership", to)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MigratableVRFConsumerV2.Contract.TransferOwnership(&_MigratableVRFConsumerV2.TransactOpts, to)
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _MigratableVRFConsumerV2.Contract.TransferOwnership(&_MigratableVRFConsumerV2.TransactOpts, to)
}

type MigratableVRFConsumerV2OwnershipTransferRequestedIterator struct {
	Event *MigratableVRFConsumerV2OwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MigratableVRFConsumerV2OwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MigratableVRFConsumerV2OwnershipTransferRequested)
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
		it.Event = new(MigratableVRFConsumerV2OwnershipTransferRequested)
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

func (it *MigratableVRFConsumerV2OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *MigratableVRFConsumerV2OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MigratableVRFConsumerV2OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MigratableVRFConsumerV2OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MigratableVRFConsumerV2.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MigratableVRFConsumerV2OwnershipTransferRequestedIterator{contract: _MigratableVRFConsumerV2.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MigratableVRFConsumerV2OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MigratableVRFConsumerV2.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MigratableVRFConsumerV2OwnershipTransferRequested)
				if err := _MigratableVRFConsumerV2.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Filterer) ParseOwnershipTransferRequested(log types.Log) (*MigratableVRFConsumerV2OwnershipTransferRequested, error) {
	event := new(MigratableVRFConsumerV2OwnershipTransferRequested)
	if err := _MigratableVRFConsumerV2.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MigratableVRFConsumerV2OwnershipTransferredIterator struct {
	Event *MigratableVRFConsumerV2OwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MigratableVRFConsumerV2OwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MigratableVRFConsumerV2OwnershipTransferred)
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
		it.Event = new(MigratableVRFConsumerV2OwnershipTransferred)
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

func (it *MigratableVRFConsumerV2OwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *MigratableVRFConsumerV2OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MigratableVRFConsumerV2OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MigratableVRFConsumerV2OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MigratableVRFConsumerV2.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MigratableVRFConsumerV2OwnershipTransferredIterator{contract: _MigratableVRFConsumerV2.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MigratableVRFConsumerV2OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MigratableVRFConsumerV2.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MigratableVRFConsumerV2OwnershipTransferred)
				if err := _MigratableVRFConsumerV2.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2Filterer) ParseOwnershipTransferred(log types.Log) (*MigratableVRFConsumerV2OwnershipTransferred, error) {
	event := new(MigratableVRFConsumerV2OwnershipTransferred)
	if err := _MigratableVRFConsumerV2.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _MigratableVRFConsumerV2.abi.Events["OwnershipTransferRequested"].ID:
		return _MigratableVRFConsumerV2.ParseOwnershipTransferRequested(log)
	case _MigratableVRFConsumerV2.abi.Events["OwnershipTransferred"].ID:
		return _MigratableVRFConsumerV2.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (MigratableVRFConsumerV2OwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (MigratableVRFConsumerV2OwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_MigratableVRFConsumerV2 *MigratableVRFConsumerV2) Address() common.Address {
	return _MigratableVRFConsumerV2.address
}

type MigratableVRFConsumerV2Interface interface {
	Owner(opts *bind.CallOpts) (common.Address, error)

	SRandomWords(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error)

	SRequestId(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, minReqConfs uint16, callbackGasLimit uint32, numWords uint32, nativePayment bool) (*types.Transaction, error)

	SetCoordinator(opts *bind.TransactOpts, coordinator common.Address) (*types.Transaction, error)

	SetSubId(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MigratableVRFConsumerV2OwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *MigratableVRFConsumerV2OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*MigratableVRFConsumerV2OwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MigratableVRFConsumerV2OwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MigratableVRFConsumerV2OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*MigratableVRFConsumerV2OwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
