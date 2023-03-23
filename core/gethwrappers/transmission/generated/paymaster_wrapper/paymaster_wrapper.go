// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package paymaster_wrapper

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

type UserOperation struct {
	Sender               common.Address
	Nonce                *big.Int
	InitCode             []byte
	CallData             []byte
	CallGasLimit         *big.Int
	VerificationGasLimit *big.Int
	PreVerificationGas   *big.Int
	MaxFeePerGas         *big.Int
	MaxPriorityFeePerGas *big.Int
	PaymasterAndData     []byte
	Signature            []byte
}

var PaymasterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"linkToken\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumIPaymaster.PostOpMode\",\"name\":\"\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"context\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"actualGasCost\",\"type\":\"uint256\"}],\"name\":\"postOp\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"verificationGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"preVerificationGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxPriorityFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"paymasterAndData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structUserOperation\",\"name\":\"userOp\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"userOpHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"maxCost\",\"type\":\"uint256\"}],\"name\":\"validatePaymasterUserOp\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"context\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"validationData\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040526611c37937e0800060005534801561001b57600080fd5b50604051610c45380380610c4583398101604081905261003a9161004b565b6001600160a01b031660805261007b565b60006020828403121561005d57600080fd5b81516001600160a01b038116811461007457600080fd5b9392505050565b608051610b9b6100aa6000396000818160560152818161010301528181610555015261060d0152610b9b6000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80631b6b6d2314610051578063a4c0ed36146100a2578063a9a23409146100b7578063f465c77e146100ca575b600080fd5b6100787f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b6100b56100b036600461075b565b6100eb565b005b6100b56100c53660046107b7565b6101e8565b6100dd6100d8366004610817565b61026c565b60405161009992919061086b565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461015a576040517f44b0e3c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208114610194576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006101a2828401846108e6565b73ffffffffffffffffffffffffffffffffffffffff81166000908152600260205260408120805492935086929091906101dc908490610939565b90915550505050505050565b6000806101f784860186610951565b915091506000805484670de0b6b3a7640000610213919061097d565b61021d91906109ba565b90506102298282610939565b73ffffffffffffffffffffffffffffffffffffffff84166000908152600260205260408120805490919061025e9084906109f5565b909155505050505050505050565b6000828152600160205260408120546060919060ff16156102ee576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f616c72656164792074726965640000000000000000000000000000000000000060448201526064015b60405180910390fd5b60006102f986610479565b905060008160005486670de0b6b3a7640000610315919061097d565b61031f91906109ba565b6103299190610939565b9050806002600061033d60208b018b6108e6565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205410156103e0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e6473000000000000000000000000000060448201526064016102e5565b600086815260016020818152604090922080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169091179055610426908801886108e6565b6040805173ffffffffffffffffffffffffffffffffffffffff9092166020830152810183905260600160405160208183030381529060405261046b60008060006106b5565b935093505050935093915050565b6000610489610120830183610a0c565b905060140361049a57506000919050565b60006104aa610120840184610a0c565b60148181106104bb576104bb610a71565b919091013560f81c91508190506106af5760006104dc610120850185610a0c565b6104ea916015908290610aa0565b8101906104f79190610aca565b905080602001516000141580156105c45750602081015181516040517f70a0823100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201527f0000000000000000000000000000000000000000000000000000000000000000909116906370a0823190602401602060405180830381865afa15801561059e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105c29190610b53565b105b156106ad57805160408083015190517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169263a9059cbb926106619260040173ffffffffffffffffffffffffffffffffffffffff929092168252602082015260400190565b6020604051808303816000875af1158015610680573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106a49190610b6c565b50806040015192505b505b50919050565b600060d08265ffffffffffff16901b60a08465ffffffffffff16901b856106dd5760006106e0565b60015b60ff161717949350505050565b73ffffffffffffffffffffffffffffffffffffffff8116811461070f57600080fd5b50565b60008083601f84011261072457600080fd5b50813567ffffffffffffffff81111561073c57600080fd5b60208301915083602082850101111561075457600080fd5b9250929050565b6000806000806060858703121561077157600080fd5b843561077c816106ed565b935060208501359250604085013567ffffffffffffffff81111561079f57600080fd5b6107ab87828801610712565b95989497509550505050565b600080600080606085870312156107cd57600080fd5b8435600381106107dc57600080fd5b9350602085013567ffffffffffffffff8111156107f857600080fd5b61080487828801610712565b9598909750949560400135949350505050565b60008060006060848603121561082c57600080fd5b833567ffffffffffffffff81111561084357600080fd5b8401610160818703121561085657600080fd5b95602085013595506040909401359392505050565b604081526000835180604084015260005b81811015610899576020818701810151606086840101520161087c565b818111156108ab576000606083860101525b50602083019390935250601f919091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01601606001919050565b6000602082840312156108f857600080fd5b8135610903816106ed565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000821982111561094c5761094c61090a565b500190565b6000806040838503121561096457600080fd5b823561096f816106ed565b946020939093013593505050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156109b5576109b561090a565b500290565b6000826109f0577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b600082821015610a0757610a0761090a565b500390565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610a4157600080fd5b83018035915067ffffffffffffffff821115610a5c57600080fd5b60200191503681900382131561075457600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60008085851115610ab057600080fd5b83861115610abd57600080fd5b5050820193919092039150565b600060608284031215610adc57600080fd5b6040516060810181811067ffffffffffffffff82111715610b26577f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040528235610b34816106ed565b8152602083810135908201526040928301359281019290925250919050565b600060208284031215610b6557600080fd5b5051919050565b600060208284031215610b7e57600080fd5b8151801515811461090357600080fdfea164736f6c634300080f000a",
}

var PaymasterABI = PaymasterMetaData.ABI

var PaymasterBin = PaymasterMetaData.Bin

func DeployPaymaster(auth *bind.TransactOpts, backend bind.ContractBackend, linkToken common.Address) (common.Address, *types.Transaction, *Paymaster, error) {
	parsed, err := PaymasterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PaymasterBin), backend, linkToken)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Paymaster{PaymasterCaller: PaymasterCaller{contract: contract}, PaymasterTransactor: PaymasterTransactor{contract: contract}, PaymasterFilterer: PaymasterFilterer{contract: contract}}, nil
}

type Paymaster struct {
	address common.Address
	abi     abi.ABI
	PaymasterCaller
	PaymasterTransactor
	PaymasterFilterer
}

type PaymasterCaller struct {
	contract *bind.BoundContract
}

type PaymasterTransactor struct {
	contract *bind.BoundContract
}

type PaymasterFilterer struct {
	contract *bind.BoundContract
}

type PaymasterSession struct {
	Contract     *Paymaster
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type PaymasterCallerSession struct {
	Contract *PaymasterCaller
	CallOpts bind.CallOpts
}

type PaymasterTransactorSession struct {
	Contract     *PaymasterTransactor
	TransactOpts bind.TransactOpts
}

type PaymasterRaw struct {
	Contract *Paymaster
}

type PaymasterCallerRaw struct {
	Contract *PaymasterCaller
}

type PaymasterTransactorRaw struct {
	Contract *PaymasterTransactor
}

func NewPaymaster(address common.Address, backend bind.ContractBackend) (*Paymaster, error) {
	abi, err := abi.JSON(strings.NewReader(PaymasterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindPaymaster(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Paymaster{address: address, abi: abi, PaymasterCaller: PaymasterCaller{contract: contract}, PaymasterTransactor: PaymasterTransactor{contract: contract}, PaymasterFilterer: PaymasterFilterer{contract: contract}}, nil
}

func NewPaymasterCaller(address common.Address, caller bind.ContractCaller) (*PaymasterCaller, error) {
	contract, err := bindPaymaster(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PaymasterCaller{contract: contract}, nil
}

func NewPaymasterTransactor(address common.Address, transactor bind.ContractTransactor) (*PaymasterTransactor, error) {
	contract, err := bindPaymaster(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PaymasterTransactor{contract: contract}, nil
}

func NewPaymasterFilterer(address common.Address, filterer bind.ContractFilterer) (*PaymasterFilterer, error) {
	contract, err := bindPaymaster(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PaymasterFilterer{contract: contract}, nil
}

func bindPaymaster(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PaymasterABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_Paymaster *PaymasterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Paymaster.Contract.PaymasterCaller.contract.Call(opts, result, method, params...)
}

func (_Paymaster *PaymasterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Paymaster.Contract.PaymasterTransactor.contract.Transfer(opts)
}

func (_Paymaster *PaymasterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Paymaster.Contract.PaymasterTransactor.contract.Transact(opts, method, params...)
}

func (_Paymaster *PaymasterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Paymaster.Contract.contract.Call(opts, result, method, params...)
}

func (_Paymaster *PaymasterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Paymaster.Contract.contract.Transfer(opts)
}

func (_Paymaster *PaymasterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Paymaster.Contract.contract.Transact(opts, method, params...)
}

func (_Paymaster *PaymasterCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Paymaster.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_Paymaster *PaymasterSession) LINK() (common.Address, error) {
	return _Paymaster.Contract.LINK(&_Paymaster.CallOpts)
}

func (_Paymaster *PaymasterCallerSession) LINK() (common.Address, error) {
	return _Paymaster.Contract.LINK(&_Paymaster.CallOpts)
}

func (_Paymaster *PaymasterTransactor) OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _Paymaster.contract.Transact(opts, "onTokenTransfer", arg0, _amount, _data)
}

func (_Paymaster *PaymasterSession) OnTokenTransfer(arg0 common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _Paymaster.Contract.OnTokenTransfer(&_Paymaster.TransactOpts, arg0, _amount, _data)
}

func (_Paymaster *PaymasterTransactorSession) OnTokenTransfer(arg0 common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _Paymaster.Contract.OnTokenTransfer(&_Paymaster.TransactOpts, arg0, _amount, _data)
}

func (_Paymaster *PaymasterTransactor) PostOp(opts *bind.TransactOpts, arg0 uint8, context []byte, actualGasCost *big.Int) (*types.Transaction, error) {
	return _Paymaster.contract.Transact(opts, "postOp", arg0, context, actualGasCost)
}

func (_Paymaster *PaymasterSession) PostOp(arg0 uint8, context []byte, actualGasCost *big.Int) (*types.Transaction, error) {
	return _Paymaster.Contract.PostOp(&_Paymaster.TransactOpts, arg0, context, actualGasCost)
}

func (_Paymaster *PaymasterTransactorSession) PostOp(arg0 uint8, context []byte, actualGasCost *big.Int) (*types.Transaction, error) {
	return _Paymaster.Contract.PostOp(&_Paymaster.TransactOpts, arg0, context, actualGasCost)
}

func (_Paymaster *PaymasterTransactor) ValidatePaymasterUserOp(opts *bind.TransactOpts, userOp UserOperation, userOpHash [32]byte, maxCost *big.Int) (*types.Transaction, error) {
	return _Paymaster.contract.Transact(opts, "validatePaymasterUserOp", userOp, userOpHash, maxCost)
}

func (_Paymaster *PaymasterSession) ValidatePaymasterUserOp(userOp UserOperation, userOpHash [32]byte, maxCost *big.Int) (*types.Transaction, error) {
	return _Paymaster.Contract.ValidatePaymasterUserOp(&_Paymaster.TransactOpts, userOp, userOpHash, maxCost)
}

func (_Paymaster *PaymasterTransactorSession) ValidatePaymasterUserOp(userOp UserOperation, userOpHash [32]byte, maxCost *big.Int) (*types.Transaction, error) {
	return _Paymaster.Contract.ValidatePaymasterUserOp(&_Paymaster.TransactOpts, userOp, userOpHash, maxCost)
}

func (_Paymaster *Paymaster) Address() common.Address {
	return _Paymaster.address
}

type PaymasterInterface interface {
	LINK(opts *bind.CallOpts) (common.Address, error)

	OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error)

	PostOp(opts *bind.TransactOpts, arg0 uint8, context []byte, actualGasCost *big.Int) (*types.Transaction, error)

	ValidatePaymasterUserOp(opts *bind.TransactOpts, userOp UserOperation, userOpHash [32]byte, maxCost *big.Int) (*types.Transaction, error)

	Address() common.Address
}
