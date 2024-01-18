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
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIL2GatewayRouter\",\"name\":\"l2GatewayRouter\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BridgeAddressCannotBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"wanted\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"got\",\"type\":\"uint256\"}],\"name\":\"InsufficientEthValue\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"MsgShouldNotContainValue\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"msgValue\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"MsgValueDoesNotMatchAmount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"depositNativeToL1\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBridgeFeeInNative\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"l1Token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"l2Token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"sendERC20\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b5060405161098938038061098983398101604081905261002f91610067565b6001600160a01b03811661005657604051635e9c404d60e11b815260040160405180910390fd5b6001600160a01b0316608052610097565b60006020828403121561007957600080fd5b81516001600160a01b038116811461009057600080fd5b9392505050565b6080516108d76100b260003960006101c401526108d76000f3fe6080604052600436106100345760003560e01c80630ff98e31146100395780632e4b1fc91461004e57806379a35b4b1461006f575b600080fd5b61004c61004736600461064d565b610082565b005b34801561005a57600080fd5b50600060405190815260200160405180910390f35b61004c61007d36600461066f565b610119565b6040517f25e1606300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526064906325e1606390349060240160206040518083038185885af11580156100f0573d6000803e3d6000fd5b50505050506040513d601f19601f8201168201806040525081019061011591906106ba565b5050565b3415610158576040517f2543d86e0000000000000000000000000000000000000000000000000000000081523460048201526024015b60405180910390fd5b61017a73ffffffffffffffffffffffffffffffffffffffff8416333084610269565b604080516020810182526000815290517f7b3a3c8b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001691637b3a3c8b916101fd91889187918791600401610741565b6000604051808303816000875af115801561021c573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261026291908101906107b9565b5050505050565b6040805173ffffffffffffffffffffffffffffffffffffffff85811660248301528416604482015260648082018490528251808303909101815260849091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f23b872dd000000000000000000000000000000000000000000000000000000001790526102fe908590610304565b50505050565b6000610366826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166104159092919063ffffffff16565b80519091501561041057808060200190518101906103849190610879565b610410576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f74207375636365656400000000000000000000000000000000000000000000606482015260840161014f565b505050565b6060610424848460008561042c565b949350505050565b6060824710156104be576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c0000000000000000000000000000000000000000000000000000606482015260840161014f565b6000808673ffffffffffffffffffffffffffffffffffffffff1685876040516104e7919061089b565b60006040518083038185875af1925050503d8060008114610524576040519150601f19603f3d011682016040523d82523d6000602084013e610529565b606091505b509150915061053a87838387610545565b979650505050505050565b606083156105db5782516000036105d45773ffffffffffffffffffffffffffffffffffffffff85163b6105d4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640161014f565b5081610424565b61042483838151156105f05781518083602001fd5b806040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161014f91906108b7565b803573ffffffffffffffffffffffffffffffffffffffff8116811461064857600080fd5b919050565b60006020828403121561065f57600080fd5b61066882610624565b9392505050565b6000806000806080858703121561068557600080fd5b61068e85610624565b935061069c60208601610624565b92506106aa60408601610624565b9396929550929360600135925050565b6000602082840312156106cc57600080fd5b5051919050565b60005b838110156106ee5781810151838201526020016106d6565b50506000910152565b6000815180845261070f8160208601602086016106d3565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b600073ffffffffffffffffffffffffffffffffffffffff80871683528086166020840152508360408301526080606083015261078060808301846106f7565b9695505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000602082840312156107cb57600080fd5b815167ffffffffffffffff808211156107e357600080fd5b818401915084601f8301126107f757600080fd5b8151818111156108095761080961078a565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561084f5761084f61078a565b8160405282815287602084870101111561086857600080fd5b61053a8360208301602088016106d3565b60006020828403121561088b57600080fd5b8151801515811461066857600080fd5b600082516108ad8184602087016106d3565b9190910192915050565b60208152600061066860208301846106f756fea164736f6c6343000813000a",
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

func (_ArbitrumL2BridgeAdapter *ArbitrumL2BridgeAdapterTransactor) SendERC20(opts *bind.TransactOpts, l1Token common.Address, l2Token common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ArbitrumL2BridgeAdapter.contract.Transact(opts, "sendERC20", l1Token, l2Token, recipient, amount)
}

func (_ArbitrumL2BridgeAdapter *ArbitrumL2BridgeAdapterSession) SendERC20(l1Token common.Address, l2Token common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ArbitrumL2BridgeAdapter.Contract.SendERC20(&_ArbitrumL2BridgeAdapter.TransactOpts, l1Token, l2Token, recipient, amount)
}

func (_ArbitrumL2BridgeAdapter *ArbitrumL2BridgeAdapterTransactorSession) SendERC20(l1Token common.Address, l2Token common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ArbitrumL2BridgeAdapter.Contract.SendERC20(&_ArbitrumL2BridgeAdapter.TransactOpts, l1Token, l2Token, recipient, amount)
}

func (_ArbitrumL2BridgeAdapter *ArbitrumL2BridgeAdapter) Address() common.Address {
	return _ArbitrumL2BridgeAdapter.address
}

type ArbitrumL2BridgeAdapterInterface interface {
	GetBridgeFeeInNative(opts *bind.CallOpts) (*big.Int, error)

	DepositNativeToL1(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error)

	SendERC20(opts *bind.TransactOpts, l1Token common.Address, l2Token common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	Address() common.Address
}
