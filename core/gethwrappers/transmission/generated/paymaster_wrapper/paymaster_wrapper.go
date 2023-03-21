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
	ABI: "[{\"inputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"linkToken\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumIPaymaster.PostOpMode\",\"name\":\"mode\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"context\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"actualGasCost\",\"type\":\"uint256\"}],\"name\":\"postOp\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"verificationGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"preVerificationGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxPriorityFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"paymasterAndData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structUserOperation\",\"name\":\"userOp\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"userOpHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"maxCost\",\"type\":\"uint256\"}],\"name\":\"validatePaymasterUserOp\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"context\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"validationData\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040526611c37937e0800060005534801561001b57600080fd5b50604051610b84380380610b8483398101604081905261003a9161004b565b6001600160a01b031660805261007b565b60006020828403121561005d57600080fd5b81516001600160a01b038116811461007457600080fd5b9392505050565b608051610ada6100aa60003960008181605601528181610103015281816104f801526105b00152610ada6000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80631b6b6d2314610051578063a4c0ed36146100a2578063a9a23409146100b7578063f465c77e146100ca575b600080fd5b6100787f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b6100b56100b03660046106c6565b6100eb565b005b6100b56100c5366004610722565b6101e8565b6100dd6100d8366004610782565b610261565b6040516100999291906107d6565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461015a576040517f44b0e3c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208114610194576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006101a282840184610851565b73ffffffffffffffffffffffffffffffffffffffff81166000908152600260205260408120805492935086929091906101dc9084906108a4565b90915550505050505050565b60006101f683850185610851565b90506000805483670de0b6b3a764000061021091906108bc565b61021a91906108f9565b73ffffffffffffffffffffffffffffffffffffffff8316600090815260026020526040812080549293508392909190610254908490610934565b9091555050505050505050565b6000828152600160205260408120546060919060ff16156102a15761028960016000806103e4565b604080516020810190915260008152925090506103dc565b60006102ac8661041c565b6000546102c186670de0b6b3a76400006108bc565b6102cb91906108f9565b6102d591906108a4565b905080600260006102e960208a018a610851565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205410156103525761033760016000806103e4565b604051806020016040528060008152509092509250506103dc565b600085815260016020818152604090922080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016909117905561039890870187610851565b6040805173ffffffffffffffffffffffffffffffffffffffff9092166020830152016040516020818303038152906040526103d660008060006103e4565b92509250505b935093915050565b600060d08265ffffffffffff16901b60a08465ffffffffffff16901b8561040c57600061040f565b60015b60ff161717949350505050565b600061042c61012083018361094b565b905060140361043d57506000919050565b600061044d61012084018461094b565b601481811061045e5761045e6109b0565b919091013560f81c915081905061065257600061047f61012085018561094b565b61048d9160159082906109df565b81019061049a9190610a09565b905080602001516000141580156105675750602081015181516040517f70a0823100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201527f0000000000000000000000000000000000000000000000000000000000000000909116906370a0823190602401602060405180830381865afa158015610541573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105659190610a92565b105b1561065057805160408083015190517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169263a9059cbb926106049260040173ffffffffffffffffffffffffffffffffffffffff929092168252602082015260400190565b6020604051808303816000875af1158015610623573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106479190610aab565b50806040015192505b505b50919050565b73ffffffffffffffffffffffffffffffffffffffff8116811461067a57600080fd5b50565b60008083601f84011261068f57600080fd5b50813567ffffffffffffffff8111156106a757600080fd5b6020830191508360208285010111156106bf57600080fd5b9250929050565b600080600080606085870312156106dc57600080fd5b84356106e781610658565b935060208501359250604085013567ffffffffffffffff81111561070a57600080fd5b6107168782880161067d565b95989497509550505050565b6000806000806060858703121561073857600080fd5b84356003811061074757600080fd5b9350602085013567ffffffffffffffff81111561076357600080fd5b61076f8782880161067d565b9598909750949560400135949350505050565b60008060006060848603121561079757600080fd5b833567ffffffffffffffff8111156107ae57600080fd5b840161016081870312156107c157600080fd5b95602085013595506040909401359392505050565b604081526000835180604084015260005b8181101561080457602081870181015160608684010152016107e7565b81811115610816576000606083860101525b50602083019390935250601f919091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01601606001919050565b60006020828403121561086357600080fd5b813561086e81610658565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600082198211156108b7576108b7610875565b500190565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156108f4576108f4610875565b500290565b60008261092f577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b60008282101561094657610946610875565b500390565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261098057600080fd5b83018035915067ffffffffffffffff82111561099b57600080fd5b6020019150368190038213156106bf57600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600080858511156109ef57600080fd5b838611156109fc57600080fd5b5050820193919092039150565b600060608284031215610a1b57600080fd5b6040516060810181811067ffffffffffffffff82111715610a65577f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040528235610a7381610658565b8152602083810135908201526040928301359281019290925250919050565b600060208284031215610aa457600080fd5b5051919050565b600060208284031215610abd57600080fd5b8151801515811461086e57600080fdfea164736f6c634300080f000a",
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

func (_Paymaster *PaymasterTransactor) PostOp(opts *bind.TransactOpts, mode uint8, context []byte, actualGasCost *big.Int) (*types.Transaction, error) {
	return _Paymaster.contract.Transact(opts, "postOp", mode, context, actualGasCost)
}

func (_Paymaster *PaymasterSession) PostOp(mode uint8, context []byte, actualGasCost *big.Int) (*types.Transaction, error) {
	return _Paymaster.Contract.PostOp(&_Paymaster.TransactOpts, mode, context, actualGasCost)
}

func (_Paymaster *PaymasterTransactorSession) PostOp(mode uint8, context []byte, actualGasCost *big.Int) (*types.Transaction, error) {
	return _Paymaster.Contract.PostOp(&_Paymaster.TransactOpts, mode, context, actualGasCost)
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

	PostOp(opts *bind.TransactOpts, mode uint8, context []byte, actualGasCost *big.Int) (*types.Transaction, error)

	ValidatePaymasterUserOp(opts *bind.TransactOpts, userOp UserOperation, userOpHash [32]byte, maxCost *big.Int) (*types.Transaction, error)

	Address() common.Address
}
