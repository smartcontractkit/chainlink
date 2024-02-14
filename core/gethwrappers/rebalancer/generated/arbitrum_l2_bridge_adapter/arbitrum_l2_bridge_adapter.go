// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package arbitrum_l2_bridge_adapter

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

var ArbitrumL2BridgeAdapterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIL2GatewayRouter\",\"name\":\"l2GatewayRouter\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BridgeAddressCannotBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"wanted\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"got\",\"type\":\"uint256\"}],\"name\":\"InsufficientEthValue\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"MsgShouldNotContainValue\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"msgValue\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"MsgValueDoesNotMatchAmount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"depositNativeToL1\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"finalizeWithdrawERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBridgeFeeInNative\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"sendERC20\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b50604051610a6c380380610a6c83398101604081905261002f91610067565b6001600160a01b03811661005657604051635e9c404d60e11b815260040160405180910390fd5b6001600160a01b0316608052610097565b60006020828403121561007957600080fd5b81516001600160a01b038116811461009057600080fd5b9392505050565b6080516109ba6100b260003960006101ff01526109ba6000f3fe60806040526004361061003f5760003560e01c80630ff98e31146100445780632e4b1fc91461005957806338314bb21461007a578063a71d98b71461009b575b600080fd5b610057610052366004610652565b6100bb565b005b34801561006557600080fd5b50604051600081526020015b60405180910390f35b34801561008657600080fd5b506100576100953660046106bd565b50505050565b6100ae6100a936600461071e565b610152565b604051610071919061080b565b6040517f25e1606300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526064906325e1606390349060240160206040518083038185885af1158015610129573d6000803e3d6000fd5b50505050506040513d601f19601f8201168201806040525081019061014e919061081e565b5050565b60603415610193576040517f2543d86e0000000000000000000000000000000000000000000000000000000081523460048201526024015b60405180910390fd5b6101b573ffffffffffffffffffffffffffffffffffffffff88163330876102a8565b604080516020810182526000815290517f7b3a3c8b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001691637b3a3c8b91610238918a918a918a91600401610837565b6000604051808303816000875af1158015610257573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261029d91908101906108af565b979650505050505050565b6040805173ffffffffffffffffffffffffffffffffffffffff8581166024830152848116604483015260648083018590528351808403909101815260849092018352602080830180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f23b872dd0000000000000000000000000000000000000000000000000000000017905283518085019094528084527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564908401526100959287929160009161037b91851690849061042a565b8051909150156104255780806020019051810190610399919061096f565b610425576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f74207375636365656400000000000000000000000000000000000000000000606482015260840161018a565b505050565b60606104398484600085610441565b949350505050565b6060824710156104d3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c0000000000000000000000000000000000000000000000000000606482015260840161018a565b6000808673ffffffffffffffffffffffffffffffffffffffff1685876040516104fc9190610991565b60006040518083038185875af1925050503d8060008114610539576040519150601f19603f3d011682016040523d82523d6000602084013e61053e565b606091505b509150915061029d87838387606083156105e05782516000036105d95773ffffffffffffffffffffffffffffffffffffffff85163b6105d9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640161018a565b5081610439565b61043983838151156105f55781518083602001fd5b806040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161018a919061080b565b803573ffffffffffffffffffffffffffffffffffffffff8116811461064d57600080fd5b919050565b60006020828403121561066457600080fd5b61066d82610629565b9392505050565b60008083601f84011261068657600080fd5b50813567ffffffffffffffff81111561069e57600080fd5b6020830191508360208285010111156106b657600080fd5b9250929050565b600080600080606085870312156106d357600080fd5b6106dc85610629565b93506106ea60208601610629565b9250604085013567ffffffffffffffff81111561070657600080fd5b61071287828801610674565b95989497509550505050565b60008060008060008060a0878903121561073757600080fd5b61074087610629565b955061074e60208801610629565b945061075c60408801610629565b935060608701359250608087013567ffffffffffffffff81111561077f57600080fd5b61078b89828a01610674565b979a9699509497509295939492505050565b60005b838110156107b85781810151838201526020016107a0565b50506000910152565b600081518084526107d981602086016020860161079d565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60208152600061066d60208301846107c1565b60006020828403121561083057600080fd5b5051919050565b600073ffffffffffffffffffffffffffffffffffffffff80871683528086166020840152508360408301526080606083015261087660808301846107c1565b9695505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000602082840312156108c157600080fd5b815167ffffffffffffffff808211156108d957600080fd5b818401915084601f8301126108ed57600080fd5b8151818111156108ff576108ff610880565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561094557610945610880565b8160405282815287602084870101111561095e57600080fd5b61029d83602083016020880161079d565b60006020828403121561098157600080fd5b8151801515811461066d57600080fd5b600082516109a381846020870161079d565b919091019291505056fea164736f6c6343000813000a",
}

var ArbitrumL2BridgeAdapterABI = ArbitrumL2BridgeAdapterMetaData.ABI

var ArbitrumL2BridgeAdapterBin = ArbitrumL2BridgeAdapterMetaData.Bin

func DeployArbitrumL2BridgeAdapter(auth *bind.TransactOpts, backend bind.ContractBackend, l2GatewayRouter common.Address) (common.Address, *types.Transaction, *ArbitrumL2BridgeAdapter, error) {
	parsed, err := ArbitrumL2BridgeAdapterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ArbitrumL2BridgeAdapterBin), backend, l2GatewayRouter)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ArbitrumL2BridgeAdapter{address: address, abi: *parsed, ArbitrumL2BridgeAdapterCaller: ArbitrumL2BridgeAdapterCaller{contract: contract}, ArbitrumL2BridgeAdapterTransactor: ArbitrumL2BridgeAdapterTransactor{contract: contract}, ArbitrumL2BridgeAdapterFilterer: ArbitrumL2BridgeAdapterFilterer{contract: contract}}, nil
}

type ArbitrumL2BridgeAdapter struct {
	address common.Address
	abi     abi.ABI
	ArbitrumL2BridgeAdapterCaller
	ArbitrumL2BridgeAdapterTransactor
	ArbitrumL2BridgeAdapterFilterer
}

type ArbitrumL2BridgeAdapterCaller struct {
	contract *bind.BoundContract
}

type ArbitrumL2BridgeAdapterTransactor struct {
	contract *bind.BoundContract
}

type ArbitrumL2BridgeAdapterFilterer struct {
	contract *bind.BoundContract
}

type ArbitrumL2BridgeAdapterSession struct {
	Contract     *ArbitrumL2BridgeAdapter
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ArbitrumL2BridgeAdapterCallerSession struct {
	Contract *ArbitrumL2BridgeAdapterCaller
	CallOpts bind.CallOpts
}

type ArbitrumL2BridgeAdapterTransactorSession struct {
	Contract     *ArbitrumL2BridgeAdapterTransactor
	TransactOpts bind.TransactOpts
}

type ArbitrumL2BridgeAdapterRaw struct {
	Contract *ArbitrumL2BridgeAdapter
}

type ArbitrumL2BridgeAdapterCallerRaw struct {
	Contract *ArbitrumL2BridgeAdapterCaller
}

type ArbitrumL2BridgeAdapterTransactorRaw struct {
	Contract *ArbitrumL2BridgeAdapterTransactor
}

func NewArbitrumL2BridgeAdapter(address common.Address, backend bind.ContractBackend) (*ArbitrumL2BridgeAdapter, error) {
	abi, err := abi.JSON(strings.NewReader(ArbitrumL2BridgeAdapterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindArbitrumL2BridgeAdapter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ArbitrumL2BridgeAdapter{address: address, abi: abi, ArbitrumL2BridgeAdapterCaller: ArbitrumL2BridgeAdapterCaller{contract: contract}, ArbitrumL2BridgeAdapterTransactor: ArbitrumL2BridgeAdapterTransactor{contract: contract}, ArbitrumL2BridgeAdapterFilterer: ArbitrumL2BridgeAdapterFilterer{contract: contract}}, nil
}

func NewArbitrumL2BridgeAdapterCaller(address common.Address, caller bind.ContractCaller) (*ArbitrumL2BridgeAdapterCaller, error) {
	contract, err := bindArbitrumL2BridgeAdapter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ArbitrumL2BridgeAdapterCaller{contract: contract}, nil
}

func NewArbitrumL2BridgeAdapterTransactor(address common.Address, transactor bind.ContractTransactor) (*ArbitrumL2BridgeAdapterTransactor, error) {
	contract, err := bindArbitrumL2BridgeAdapter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ArbitrumL2BridgeAdapterTransactor{contract: contract}, nil
}

func NewArbitrumL2BridgeAdapterFilterer(address common.Address, filterer bind.ContractFilterer) (*ArbitrumL2BridgeAdapterFilterer, error) {
	contract, err := bindArbitrumL2BridgeAdapter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ArbitrumL2BridgeAdapterFilterer{contract: contract}, nil
}

func bindArbitrumL2BridgeAdapter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ArbitrumL2BridgeAdapterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ArbitrumL2BridgeAdapter *ArbitrumL2BridgeAdapterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ArbitrumL2BridgeAdapter.Contract.ArbitrumL2BridgeAdapterCaller.contract.Call(opts, result, method, params...)
}

func (_ArbitrumL2BridgeAdapter *ArbitrumL2BridgeAdapterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArbitrumL2BridgeAdapter.Contract.ArbitrumL2BridgeAdapterTransactor.contract.Transfer(opts)
}

func (_ArbitrumL2BridgeAdapter *ArbitrumL2BridgeAdapterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ArbitrumL2BridgeAdapter.Contract.ArbitrumL2BridgeAdapterTransactor.contract.Transact(opts, method, params...)
}

func (_ArbitrumL2BridgeAdapter *ArbitrumL2BridgeAdapterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ArbitrumL2BridgeAdapter.Contract.contract.Call(opts, result, method, params...)
}

func (_ArbitrumL2BridgeAdapter *ArbitrumL2BridgeAdapterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArbitrumL2BridgeAdapter.Contract.contract.Transfer(opts)
}

func (_ArbitrumL2BridgeAdapter *ArbitrumL2BridgeAdapterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ArbitrumL2BridgeAdapter.Contract.contract.Transact(opts, method, params...)
}

func (_ArbitrumL2BridgeAdapter *ArbitrumL2BridgeAdapterCaller) GetBridgeFeeInNative(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ArbitrumL2BridgeAdapter.contract.Call(opts, &out, "getBridgeFeeInNative")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ArbitrumL2BridgeAdapter *ArbitrumL2BridgeAdapterSession) GetBridgeFeeInNative() (*big.Int, error) {
	return _ArbitrumL2BridgeAdapter.Contract.GetBridgeFeeInNative(&_ArbitrumL2BridgeAdapter.CallOpts)
}

func (_ArbitrumL2BridgeAdapter *ArbitrumL2BridgeAdapterCallerSession) GetBridgeFeeInNative() (*big.Int, error) {
	return _ArbitrumL2BridgeAdapter.Contract.GetBridgeFeeInNative(&_ArbitrumL2BridgeAdapter.CallOpts)
}

func (_ArbitrumL2BridgeAdapter *ArbitrumL2BridgeAdapterTransactor) DepositNativeToL1(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _ArbitrumL2BridgeAdapter.contract.Transact(opts, "depositNativeToL1", recipient)
}

func (_ArbitrumL2BridgeAdapter *ArbitrumL2BridgeAdapterSession) DepositNativeToL1(recipient common.Address) (*types.Transaction, error) {
	return _ArbitrumL2BridgeAdapter.Contract.DepositNativeToL1(&_ArbitrumL2BridgeAdapter.TransactOpts, recipient)
}

func (_ArbitrumL2BridgeAdapter *ArbitrumL2BridgeAdapterTransactorSession) DepositNativeToL1(recipient common.Address) (*types.Transaction, error) {
	return _ArbitrumL2BridgeAdapter.Contract.DepositNativeToL1(&_ArbitrumL2BridgeAdapter.TransactOpts, recipient)
}

func (_ArbitrumL2BridgeAdapter *ArbitrumL2BridgeAdapterTransactor) FinalizeWithdrawERC20(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 []byte) (*types.Transaction, error) {
	return _ArbitrumL2BridgeAdapter.contract.Transact(opts, "finalizeWithdrawERC20", arg0, arg1, arg2)
}

func (_ArbitrumL2BridgeAdapter *ArbitrumL2BridgeAdapterSession) FinalizeWithdrawERC20(arg0 common.Address, arg1 common.Address, arg2 []byte) (*types.Transaction, error) {
	return _ArbitrumL2BridgeAdapter.Contract.FinalizeWithdrawERC20(&_ArbitrumL2BridgeAdapter.TransactOpts, arg0, arg1, arg2)
}

func (_ArbitrumL2BridgeAdapter *ArbitrumL2BridgeAdapterTransactorSession) FinalizeWithdrawERC20(arg0 common.Address, arg1 common.Address, arg2 []byte) (*types.Transaction, error) {
	return _ArbitrumL2BridgeAdapter.Contract.FinalizeWithdrawERC20(&_ArbitrumL2BridgeAdapter.TransactOpts, arg0, arg1, arg2)
}

func (_ArbitrumL2BridgeAdapter *ArbitrumL2BridgeAdapterTransactor) SendERC20(opts *bind.TransactOpts, localToken common.Address, remoteToken common.Address, recipient common.Address, amount *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _ArbitrumL2BridgeAdapter.contract.Transact(opts, "sendERC20", localToken, remoteToken, recipient, amount, arg4)
}

func (_ArbitrumL2BridgeAdapter *ArbitrumL2BridgeAdapterSession) SendERC20(localToken common.Address, remoteToken common.Address, recipient common.Address, amount *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _ArbitrumL2BridgeAdapter.Contract.SendERC20(&_ArbitrumL2BridgeAdapter.TransactOpts, localToken, remoteToken, recipient, amount, arg4)
}

func (_ArbitrumL2BridgeAdapter *ArbitrumL2BridgeAdapterTransactorSession) SendERC20(localToken common.Address, remoteToken common.Address, recipient common.Address, amount *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _ArbitrumL2BridgeAdapter.Contract.SendERC20(&_ArbitrumL2BridgeAdapter.TransactOpts, localToken, remoteToken, recipient, amount, arg4)
}

func (_ArbitrumL2BridgeAdapter *ArbitrumL2BridgeAdapter) Address() common.Address {
	return _ArbitrumL2BridgeAdapter.address
}

type ArbitrumL2BridgeAdapterInterface interface {
	GetBridgeFeeInNative(opts *bind.CallOpts) (*big.Int, error)

	DepositNativeToL1(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error)

	FinalizeWithdrawERC20(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 []byte) (*types.Transaction, error)

	SendERC20(opts *bind.TransactOpts, localToken common.Address, remoteToken common.Address, recipient common.Address, amount *big.Int, arg4 []byte) (*types.Transaction, error)

	Address() common.Address
}
