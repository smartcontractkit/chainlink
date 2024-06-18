// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package optimism_l2_bridge_adapter

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

var OptimismL2BridgeAdapterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIWrappedNative\",\"name\":\"wrappedNative\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BridgeAddressCannotBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"wanted\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"got\",\"type\":\"uint256\"}],\"name\":\"InsufficientEthValue\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"MsgShouldNotContainValue\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"msgValue\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"MsgValueDoesNotMatchAmount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"finalizeWithdrawERC20\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBridgeFeeInNative\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getL2Bridge\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getWrappedNative\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"sendERC20\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60c0604052734200000000000000000000000000000000000010608052600080546001600160401b031916905534801561003857600080fd5b50604051610c2e380380610c2e83398101604081905261005791610068565b6001600160a01b031660a052610098565b60006020828403121561007a57600080fd5b81516001600160a01b038116811461009157600080fd5b9392505050565b60805160a051610b4f6100df6000396000818161013e0152818161024201526102c2015260008181609a0152818161036e0152818161043c01526104f60152610b4f6000f3fe60806040526004361061005e5760003560e01c806338314bb21161004357806338314bb2146100df578063a71d98b71461010f578063e861e9071461012f57600080fd5b80632e4b1fc91461006a5780633429072c1461008b57600080fd5b3661006557005b600080fd5b34801561007657600080fd5b50604051600081526020015b60405180910390f35b34801561009757600080fd5b507f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610082565b3480156100eb57600080fd5b506100ff6100fa366004610903565b610162565b6040519015158152602001610082565b61012261011d366004610964565b61016d565b6040516100829190610a51565b34801561013b57600080fd5b507f00000000000000000000000000000000000000000000000000000000000000006100ba565b60015b949350505050565b606034156101ae576040517f2543d86e0000000000000000000000000000000000000000000000000000000081523460048201526024015b60405180910390fd5b6101d073ffffffffffffffffffffffffffffffffffffffff8816333087610574565b6000805467ffffffffffffffff1681806101e983610a6b565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555060405160200161022f919067ffffffffffffffff91909116815260200190565b60405160208183030381529060405290507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168873ffffffffffffffffffffffffffffffffffffffff16036103ff576040517f2e1a7d4d000000000000000000000000000000000000000000000000000000008152600481018690527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690632e1a7d4d90602401600060405180830381600087803b15801561031b57600080fd5b505af115801561032f573d6000803e3d6000fd5b50506040517fa3a7954800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016925063a3a79548915087906103c49073deaddeaddeaddeaddeaddeaddeaddeaddead0000908b9084906000908990600401610ab9565b6000604051808303818588803b1580156103dd57600080fd5b505af11580156103f1573d6000803e3d6000fd5b50505050508091505061056a565b6040517f095ea7b300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000811660048301526024820187905289169063095ea7b3906044016020604051808303816000875af1158015610494573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104b89190610b04565b506040517fa3a7954800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169063a3a7954890610534908b908a908a906000908890600401610ab9565b600060405180830381600087803b15801561054e57600080fd5b505af1158015610562573d6000803e3d6000fd5b509293505050505b9695505050505050565b6040805173ffffffffffffffffffffffffffffffffffffffff85811660248301528416604482015260648082018490528251808303909101815260849091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f23b872dd0000000000000000000000000000000000000000000000000000000017905261060990859061060f565b50505050565b6000610671826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166107209092919063ffffffff16565b80519091501561071b578080602001905181019061068f9190610b04565b61071b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f7420737563636565640000000000000000000000000000000000000000000060648201526084016101a5565b505050565b60606101658484600085856000808673ffffffffffffffffffffffffffffffffffffffff1685876040516107549190610b26565b60006040518083038185875af1925050503d8060008114610791576040519150601f19603f3d011682016040523d82523d6000602084013e610796565b606091505b50915091506107a7878383876107b2565b979650505050505050565b606083156108485782516000036108415773ffffffffffffffffffffffffffffffffffffffff85163b610841576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016101a5565b5081610165565b610165838381511561085d5781518083602001fd5b806040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101a59190610a51565b803573ffffffffffffffffffffffffffffffffffffffff811681146108b557600080fd5b919050565b60008083601f8401126108cc57600080fd5b50813567ffffffffffffffff8111156108e457600080fd5b6020830191508360208285010111156108fc57600080fd5b9250929050565b6000806000806060858703121561091957600080fd5b61092285610891565b935061093060208601610891565b9250604085013567ffffffffffffffff81111561094c57600080fd5b610958878288016108ba565b95989497509550505050565b60008060008060008060a0878903121561097d57600080fd5b61098687610891565b955061099460208801610891565b94506109a260408801610891565b935060608701359250608087013567ffffffffffffffff8111156109c557600080fd5b6109d189828a016108ba565b979a9699509497509295939492505050565b60005b838110156109fe5781810151838201526020016109e6565b50506000910152565b60008151808452610a1f8160208601602086016109e3565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000610a646020830184610a07565b9392505050565b600067ffffffffffffffff808316818103610aaf577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6001019392505050565b600073ffffffffffffffffffffffffffffffffffffffff808816835280871660208401525084604083015263ffffffff8416606083015260a060808301526107a760a0830184610a07565b600060208284031215610b1657600080fd5b81518015158114610a6457600080fd5b60008251610b388184602087016109e3565b919091019291505056fea164736f6c6343000818000a",
}

var OptimismL2BridgeAdapterABI = OptimismL2BridgeAdapterMetaData.ABI

var OptimismL2BridgeAdapterBin = OptimismL2BridgeAdapterMetaData.Bin

func DeployOptimismL2BridgeAdapter(auth *bind.TransactOpts, backend bind.ContractBackend, wrappedNative common.Address) (common.Address, *types.Transaction, *OptimismL2BridgeAdapter, error) {
	parsed, err := OptimismL2BridgeAdapterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OptimismL2BridgeAdapterBin), backend, wrappedNative)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OptimismL2BridgeAdapter{address: address, abi: *parsed, OptimismL2BridgeAdapterCaller: OptimismL2BridgeAdapterCaller{contract: contract}, OptimismL2BridgeAdapterTransactor: OptimismL2BridgeAdapterTransactor{contract: contract}, OptimismL2BridgeAdapterFilterer: OptimismL2BridgeAdapterFilterer{contract: contract}}, nil
}

type OptimismL2BridgeAdapter struct {
	address common.Address
	abi     abi.ABI
	OptimismL2BridgeAdapterCaller
	OptimismL2BridgeAdapterTransactor
	OptimismL2BridgeAdapterFilterer
}

type OptimismL2BridgeAdapterCaller struct {
	contract *bind.BoundContract
}

type OptimismL2BridgeAdapterTransactor struct {
	contract *bind.BoundContract
}

type OptimismL2BridgeAdapterFilterer struct {
	contract *bind.BoundContract
}

type OptimismL2BridgeAdapterSession struct {
	Contract     *OptimismL2BridgeAdapter
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OptimismL2BridgeAdapterCallerSession struct {
	Contract *OptimismL2BridgeAdapterCaller
	CallOpts bind.CallOpts
}

type OptimismL2BridgeAdapterTransactorSession struct {
	Contract     *OptimismL2BridgeAdapterTransactor
	TransactOpts bind.TransactOpts
}

type OptimismL2BridgeAdapterRaw struct {
	Contract *OptimismL2BridgeAdapter
}

type OptimismL2BridgeAdapterCallerRaw struct {
	Contract *OptimismL2BridgeAdapterCaller
}

type OptimismL2BridgeAdapterTransactorRaw struct {
	Contract *OptimismL2BridgeAdapterTransactor
}

func NewOptimismL2BridgeAdapter(address common.Address, backend bind.ContractBackend) (*OptimismL2BridgeAdapter, error) {
	abi, err := abi.JSON(strings.NewReader(OptimismL2BridgeAdapterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOptimismL2BridgeAdapter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OptimismL2BridgeAdapter{address: address, abi: abi, OptimismL2BridgeAdapterCaller: OptimismL2BridgeAdapterCaller{contract: contract}, OptimismL2BridgeAdapterTransactor: OptimismL2BridgeAdapterTransactor{contract: contract}, OptimismL2BridgeAdapterFilterer: OptimismL2BridgeAdapterFilterer{contract: contract}}, nil
}

func NewOptimismL2BridgeAdapterCaller(address common.Address, caller bind.ContractCaller) (*OptimismL2BridgeAdapterCaller, error) {
	contract, err := bindOptimismL2BridgeAdapter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismL2BridgeAdapterCaller{contract: contract}, nil
}

func NewOptimismL2BridgeAdapterTransactor(address common.Address, transactor bind.ContractTransactor) (*OptimismL2BridgeAdapterTransactor, error) {
	contract, err := bindOptimismL2BridgeAdapter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismL2BridgeAdapterTransactor{contract: contract}, nil
}

func NewOptimismL2BridgeAdapterFilterer(address common.Address, filterer bind.ContractFilterer) (*OptimismL2BridgeAdapterFilterer, error) {
	contract, err := bindOptimismL2BridgeAdapter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OptimismL2BridgeAdapterFilterer{contract: contract}, nil
}

func bindOptimismL2BridgeAdapter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OptimismL2BridgeAdapterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_OptimismL2BridgeAdapter *OptimismL2BridgeAdapterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismL2BridgeAdapter.Contract.OptimismL2BridgeAdapterCaller.contract.Call(opts, result, method, params...)
}

func (_OptimismL2BridgeAdapter *OptimismL2BridgeAdapterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismL2BridgeAdapter.Contract.OptimismL2BridgeAdapterTransactor.contract.Transfer(opts)
}

func (_OptimismL2BridgeAdapter *OptimismL2BridgeAdapterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismL2BridgeAdapter.Contract.OptimismL2BridgeAdapterTransactor.contract.Transact(opts, method, params...)
}

func (_OptimismL2BridgeAdapter *OptimismL2BridgeAdapterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismL2BridgeAdapter.Contract.contract.Call(opts, result, method, params...)
}

func (_OptimismL2BridgeAdapter *OptimismL2BridgeAdapterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismL2BridgeAdapter.Contract.contract.Transfer(opts)
}

func (_OptimismL2BridgeAdapter *OptimismL2BridgeAdapterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismL2BridgeAdapter.Contract.contract.Transact(opts, method, params...)
}

func (_OptimismL2BridgeAdapter *OptimismL2BridgeAdapterCaller) FinalizeWithdrawERC20(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address, arg2 []byte) (bool, error) {
	var out []interface{}
	err := _OptimismL2BridgeAdapter.contract.Call(opts, &out, "finalizeWithdrawERC20", arg0, arg1, arg2)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_OptimismL2BridgeAdapter *OptimismL2BridgeAdapterSession) FinalizeWithdrawERC20(arg0 common.Address, arg1 common.Address, arg2 []byte) (bool, error) {
	return _OptimismL2BridgeAdapter.Contract.FinalizeWithdrawERC20(&_OptimismL2BridgeAdapter.CallOpts, arg0, arg1, arg2)
}

func (_OptimismL2BridgeAdapter *OptimismL2BridgeAdapterCallerSession) FinalizeWithdrawERC20(arg0 common.Address, arg1 common.Address, arg2 []byte) (bool, error) {
	return _OptimismL2BridgeAdapter.Contract.FinalizeWithdrawERC20(&_OptimismL2BridgeAdapter.CallOpts, arg0, arg1, arg2)
}

func (_OptimismL2BridgeAdapter *OptimismL2BridgeAdapterCaller) GetBridgeFeeInNative(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OptimismL2BridgeAdapter.contract.Call(opts, &out, "getBridgeFeeInNative")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OptimismL2BridgeAdapter *OptimismL2BridgeAdapterSession) GetBridgeFeeInNative() (*big.Int, error) {
	return _OptimismL2BridgeAdapter.Contract.GetBridgeFeeInNative(&_OptimismL2BridgeAdapter.CallOpts)
}

func (_OptimismL2BridgeAdapter *OptimismL2BridgeAdapterCallerSession) GetBridgeFeeInNative() (*big.Int, error) {
	return _OptimismL2BridgeAdapter.Contract.GetBridgeFeeInNative(&_OptimismL2BridgeAdapter.CallOpts)
}

func (_OptimismL2BridgeAdapter *OptimismL2BridgeAdapterCaller) GetL2Bridge(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OptimismL2BridgeAdapter.contract.Call(opts, &out, "getL2Bridge")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OptimismL2BridgeAdapter *OptimismL2BridgeAdapterSession) GetL2Bridge() (common.Address, error) {
	return _OptimismL2BridgeAdapter.Contract.GetL2Bridge(&_OptimismL2BridgeAdapter.CallOpts)
}

func (_OptimismL2BridgeAdapter *OptimismL2BridgeAdapterCallerSession) GetL2Bridge() (common.Address, error) {
	return _OptimismL2BridgeAdapter.Contract.GetL2Bridge(&_OptimismL2BridgeAdapter.CallOpts)
}

func (_OptimismL2BridgeAdapter *OptimismL2BridgeAdapterCaller) GetWrappedNative(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OptimismL2BridgeAdapter.contract.Call(opts, &out, "getWrappedNative")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OptimismL2BridgeAdapter *OptimismL2BridgeAdapterSession) GetWrappedNative() (common.Address, error) {
	return _OptimismL2BridgeAdapter.Contract.GetWrappedNative(&_OptimismL2BridgeAdapter.CallOpts)
}

func (_OptimismL2BridgeAdapter *OptimismL2BridgeAdapterCallerSession) GetWrappedNative() (common.Address, error) {
	return _OptimismL2BridgeAdapter.Contract.GetWrappedNative(&_OptimismL2BridgeAdapter.CallOpts)
}

func (_OptimismL2BridgeAdapter *OptimismL2BridgeAdapterTransactor) SendERC20(opts *bind.TransactOpts, localToken common.Address, arg1 common.Address, recipient common.Address, amount *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _OptimismL2BridgeAdapter.contract.Transact(opts, "sendERC20", localToken, arg1, recipient, amount, arg4)
}

func (_OptimismL2BridgeAdapter *OptimismL2BridgeAdapterSession) SendERC20(localToken common.Address, arg1 common.Address, recipient common.Address, amount *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _OptimismL2BridgeAdapter.Contract.SendERC20(&_OptimismL2BridgeAdapter.TransactOpts, localToken, arg1, recipient, amount, arg4)
}

func (_OptimismL2BridgeAdapter *OptimismL2BridgeAdapterTransactorSession) SendERC20(localToken common.Address, arg1 common.Address, recipient common.Address, amount *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _OptimismL2BridgeAdapter.Contract.SendERC20(&_OptimismL2BridgeAdapter.TransactOpts, localToken, arg1, recipient, amount, arg4)
}

func (_OptimismL2BridgeAdapter *OptimismL2BridgeAdapterTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismL2BridgeAdapter.contract.RawTransact(opts, nil)
}

func (_OptimismL2BridgeAdapter *OptimismL2BridgeAdapterSession) Receive() (*types.Transaction, error) {
	return _OptimismL2BridgeAdapter.Contract.Receive(&_OptimismL2BridgeAdapter.TransactOpts)
}

func (_OptimismL2BridgeAdapter *OptimismL2BridgeAdapterTransactorSession) Receive() (*types.Transaction, error) {
	return _OptimismL2BridgeAdapter.Contract.Receive(&_OptimismL2BridgeAdapter.TransactOpts)
}

func (_OptimismL2BridgeAdapter *OptimismL2BridgeAdapter) Address() common.Address {
	return _OptimismL2BridgeAdapter.address
}

type OptimismL2BridgeAdapterInterface interface {
	FinalizeWithdrawERC20(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address, arg2 []byte) (bool, error)

	GetBridgeFeeInNative(opts *bind.CallOpts) (*big.Int, error)

	GetL2Bridge(opts *bind.CallOpts) (common.Address, error)

	GetWrappedNative(opts *bind.CallOpts) (common.Address, error)

	SendERC20(opts *bind.TransactOpts, localToken common.Address, arg1 common.Address, recipient common.Address, amount *big.Int, arg4 []byte) (*types.Transaction, error)

	Receive(opts *bind.TransactOpts) (*types.Transaction, error)

	Address() common.Address
}
