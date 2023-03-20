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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumIPaymaster.PostOpMode\",\"name\":\"mode\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"context\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"actualGasCost\",\"type\":\"uint256\"}],\"name\":\"postOp\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"verificationGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"preVerificationGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxPriorityFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"paymasterAndData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structUserOperation\",\"name\":\"userOp\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"userOpHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"maxCost\",\"type\":\"uint256\"}],\"name\":\"validatePaymasterUserOp\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"context\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"validationData\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040526611c37937e0800060015534801561001b57600080fd5b5060405161078438038061078483398101604081905261003a9161005f565b600080546001600160a01b0319166001600160a01b039290921691909117905561008f565b60006020828403121561007157600080fd5b81516001600160a01b038116811461008857600080fd5b9392505050565b6106e68061009e6000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80631b6b6d2314610051578063a4c0ed361461009b578063a9a23409146100b0578063f465c77e146100c3575b600080fd5b6000546100719073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b6100ae6100a9366004610454565b6100e4565b005b6100ae6100be3660046104b0565b6101c3565b6100d66100d1366004610510565b61023d565b604051610092929190610564565b60005473ffffffffffffffffffffffffffffffffffffffff163314610135576040517f44b0e3c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6020811461016f576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061017d828401846105df565b73ffffffffffffffffffffffffffffffffffffffff81166000908152600360205260408120805492935086929091906101b7908490610632565b90915550505050505050565b60006101d1838501856105df565b9050600060015483670de0b6b3a76400006101ec919061064a565b6101f69190610687565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600360205260408120805492935083929091906102309084906106c2565b9091555050505050505050565b6000828152600260205260408120546060919060ff161561027d5761026560016000806103ae565b604080516020810190915260008152925090506103a6565b60015460009061029585670de0b6b3a764000061064a565b61029f9190610687565b905080600360006102b360208a018a6105df565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054101561031c5761030160016000806103ae565b604051806020016040528060008152509092509250506103a6565b600085815260026020908152604090912080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055610362908701876105df565b6040805173ffffffffffffffffffffffffffffffffffffffff9092166020830152016040516020818303038152906040526103a060008060006103ae565b92509250505b935093915050565b600060d08265ffffffffffff16901b60a08465ffffffffffff16901b856103d65760006103d9565b60015b60ff161717949350505050565b73ffffffffffffffffffffffffffffffffffffffff8116811461040857600080fd5b50565b60008083601f84011261041d57600080fd5b50813567ffffffffffffffff81111561043557600080fd5b60208301915083602082850101111561044d57600080fd5b9250929050565b6000806000806060858703121561046a57600080fd5b8435610475816103e6565b935060208501359250604085013567ffffffffffffffff81111561049857600080fd5b6104a48782880161040b565b95989497509550505050565b600080600080606085870312156104c657600080fd5b8435600381106104d557600080fd5b9350602085013567ffffffffffffffff8111156104f157600080fd5b6104fd8782880161040b565b9598909750949560400135949350505050565b60008060006060848603121561052557600080fd5b833567ffffffffffffffff81111561053c57600080fd5b8401610160818703121561054f57600080fd5b95602085013595506040909401359392505050565b604081526000835180604084015260005b818110156105925760208187018101516060868401015201610575565b818111156105a4576000606083860101525b50602083019390935250601f919091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01601606001919050565b6000602082840312156105f157600080fd5b81356105fc816103e6565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000821982111561064557610645610603565b500190565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561068257610682610603565b500290565b6000826106bd577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b6000828210156106d4576106d4610603565b50039056fea164736f6c634300080f000a",
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
