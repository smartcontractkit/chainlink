// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package sca_wrapper

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

var SCAMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"entryPoint\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BadFormatOrOOG\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"currentNonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonceGiven\",\"type\":\"uint256\"}],\"name\":\"IncorrectNonce\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"operationHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"NotAuthorized\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"currentTimestamp\",\"type\":\"uint256\"}],\"name\":\"TransactionExpired\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint48\",\"name\":\"deadline\",\"type\":\"uint48\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"executeTransactionFromEntryPoint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_entryPoint\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_nonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"verificationGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"preVerificationGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxPriorityFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"paymasterAndData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structUserOperation\",\"name\":\"userOp\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"userOpHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validateUserOp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"validationData\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c060405234801561001057600080fd5b50604051610acb380380610acb83398101604081905261002f91610062565b6001600160a01b039182166080521660a052610095565b80516001600160a01b038116811461005d57600080fd5b919050565b6000806040838503121561007557600080fd5b61007e83610046565b915061008c60208401610046565b90509250929050565b60805160a051610a046100c760003960008181607101526102b801526000818161010101526101e30152610a046000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c80637eccf63e116100505780637eccf63e146100de57806389553be4146100e7578063dba6335f146100fc57600080fd5b8063140fcfb11461006c5780633a871cdd146100bd575b600080fd5b6100937f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b6100d06100cb366004610646565b610123565b6040519081526020016100b4565b6100d060005481565b6100fa6100f53660046106da565b6102a0565b005b6100937f000000000000000000000000000000000000000000000000000000000000000081565b60008054846020013514610179576000546040517f7ba633940000000000000000000000000000000000000000000000000000000081526004810191909152602085013560248201526044015b60405180910390fd5b60006101858430610439565b90506000610197610140870187610777565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509293505073ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016915061021190508284610550565b73ffffffffffffffffffffffffffffffffffffffff161461024257610239600160008061060e565b92505050610299565b60008054908061025183610812565b90915550600090506102666060880188610777565b61027491600490829061084a565b81019061028191906108a3565b5092505050610293600082600061060e565b93505050505b9392505050565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610311576040517f4a0bfec1000000000000000000000000000000000000000000000000000000008152336004820152602401610170565b65ffffffffffff83161580159061032f57508265ffffffffffff1642115b15610376576040517f300249d700000000000000000000000000000000000000000000000000000000815265ffffffffffff84166004820152426024820152604401610170565b6000808673ffffffffffffffffffffffffffffffffffffffff168685856040516103a192919061099f565b60006040518083038185875af1925050503d80600081146103de576040519150601f19603f3d011682016040523d82523d6000602084013e6103e3565b606091505b509150915081610430578051600003610428576040517f20e9b5d200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805181602001fd5b50505050505050565b604080517f4750045d47fce615521b32cee713ff8db50147e98aec5ca94926b52651ca3fa0602080830191909152818301859052825180830384018152606080840185528151918301919091207f190000000000000000000000000000000000000000000000000000000000000060808501527f010000000000000000000000000000000000000000000000000000000000000060818501527f1c7d3b72b37a35523e273aaadd7b4cd66f618bb81429ab053412d51f50ccea6160828501524660a28501529085901b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001660c284015260d6808401919091528351808403909101815260f690920190925280519101205b92915050565b602082015160408084015184516000939284918791908110610574576105746109af565b016020015160f81c905060018561058c83601b6109de565b6040805160008152602081018083529390935260ff90911690820152606081018590526080810184905260a0016020604051602081039080840390855afa1580156105db573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00151979650505050505050565b600060d08265ffffffffffff16901b60a08465ffffffffffff16901b85610636576000610639565b60015b60ff161717949350505050565b60008060006060848603121561065b57600080fd5b833567ffffffffffffffff81111561067257600080fd5b8401610160818703121561068557600080fd5b95602085013595506040909401359392505050565b73ffffffffffffffffffffffffffffffffffffffff811681146106bc57600080fd5b50565b803565ffffffffffff811681146106d557600080fd5b919050565b6000806000806000608086880312156106f257600080fd5b85356106fd8161069a565b945060208601359350610712604087016106bf565b9250606086013567ffffffffffffffff8082111561072f57600080fd5b818801915088601f83011261074357600080fd5b81358181111561075257600080fd5b89602082850101111561076457600080fd5b9699959850939650602001949392505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126107ac57600080fd5b83018035915067ffffffffffffffff8211156107c757600080fd5b6020019150368190038213156107dc57600080fd5b9250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203610843576108436107e3565b5060010190565b6000808585111561085a57600080fd5b8386111561086757600080fd5b5050820193919092039150565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080600080608085870312156108b957600080fd5b84356108c48161069a565b9350602085013592506108d9604086016106bf565b9150606085013567ffffffffffffffff808211156108f657600080fd5b818701915087601f83011261090a57600080fd5b81358181111561091c5761091c610874565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561096257610962610874565b816040528281528a602084870101111561097b57600080fd5b82602086016020830137600060208483010152809550505050505092959194509250565b8183823760009101908152919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60ff818116838216019081111561054a5761054a6107e356fea164736f6c6343000813000a",
}

var SCAABI = SCAMetaData.ABI

var SCABin = SCAMetaData.Bin

func DeploySCA(auth *bind.TransactOpts, backend bind.ContractBackend, owner common.Address, entryPoint common.Address) (common.Address, *types.Transaction, *SCA, error) {
	parsed, err := SCAMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SCABin), backend, owner, entryPoint)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SCA{address: address, abi: *parsed, SCACaller: SCACaller{contract: contract}, SCATransactor: SCATransactor{contract: contract}, SCAFilterer: SCAFilterer{contract: contract}}, nil
}

type SCA struct {
	address common.Address
	abi     abi.ABI
	SCACaller
	SCATransactor
	SCAFilterer
}

type SCACaller struct {
	contract *bind.BoundContract
}

type SCATransactor struct {
	contract *bind.BoundContract
}

type SCAFilterer struct {
	contract *bind.BoundContract
}

type SCASession struct {
	Contract     *SCA
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type SCACallerSession struct {
	Contract *SCACaller
	CallOpts bind.CallOpts
}

type SCATransactorSession struct {
	Contract     *SCATransactor
	TransactOpts bind.TransactOpts
}

type SCARaw struct {
	Contract *SCA
}

type SCACallerRaw struct {
	Contract *SCACaller
}

type SCATransactorRaw struct {
	Contract *SCATransactor
}

func NewSCA(address common.Address, backend bind.ContractBackend) (*SCA, error) {
	abi, err := abi.JSON(strings.NewReader(SCAABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindSCA(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SCA{address: address, abi: abi, SCACaller: SCACaller{contract: contract}, SCATransactor: SCATransactor{contract: contract}, SCAFilterer: SCAFilterer{contract: contract}}, nil
}

func NewSCACaller(address common.Address, caller bind.ContractCaller) (*SCACaller, error) {
	contract, err := bindSCA(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SCACaller{contract: contract}, nil
}

func NewSCATransactor(address common.Address, transactor bind.ContractTransactor) (*SCATransactor, error) {
	contract, err := bindSCA(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SCATransactor{contract: contract}, nil
}

func NewSCAFilterer(address common.Address, filterer bind.ContractFilterer) (*SCAFilterer, error) {
	contract, err := bindSCA(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SCAFilterer{contract: contract}, nil
}

func bindSCA(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SCAMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_SCA *SCARaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SCA.Contract.SCACaller.contract.Call(opts, result, method, params...)
}

func (_SCA *SCARaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SCA.Contract.SCATransactor.contract.Transfer(opts)
}

func (_SCA *SCARaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SCA.Contract.SCATransactor.contract.Transact(opts, method, params...)
}

func (_SCA *SCACallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SCA.Contract.contract.Call(opts, result, method, params...)
}

func (_SCA *SCATransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SCA.Contract.contract.Transfer(opts)
}

func (_SCA *SCATransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SCA.Contract.contract.Transact(opts, method, params...)
}

func (_SCA *SCACaller) IEntryPoint(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SCA.contract.Call(opts, &out, "i_entryPoint")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_SCA *SCASession) IEntryPoint() (common.Address, error) {
	return _SCA.Contract.IEntryPoint(&_SCA.CallOpts)
}

func (_SCA *SCACallerSession) IEntryPoint() (common.Address, error) {
	return _SCA.Contract.IEntryPoint(&_SCA.CallOpts)
}

func (_SCA *SCACaller) IOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SCA.contract.Call(opts, &out, "i_owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_SCA *SCASession) IOwner() (common.Address, error) {
	return _SCA.Contract.IOwner(&_SCA.CallOpts)
}

func (_SCA *SCACallerSession) IOwner() (common.Address, error) {
	return _SCA.Contract.IOwner(&_SCA.CallOpts)
}

func (_SCA *SCACaller) SNonce(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SCA.contract.Call(opts, &out, "s_nonce")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_SCA *SCASession) SNonce() (*big.Int, error) {
	return _SCA.Contract.SNonce(&_SCA.CallOpts)
}

func (_SCA *SCACallerSession) SNonce() (*big.Int, error) {
	return _SCA.Contract.SNonce(&_SCA.CallOpts)
}

func (_SCA *SCATransactor) ExecuteTransactionFromEntryPoint(opts *bind.TransactOpts, to common.Address, value *big.Int, deadline *big.Int, data []byte) (*types.Transaction, error) {
	return _SCA.contract.Transact(opts, "executeTransactionFromEntryPoint", to, value, deadline, data)
}

func (_SCA *SCASession) ExecuteTransactionFromEntryPoint(to common.Address, value *big.Int, deadline *big.Int, data []byte) (*types.Transaction, error) {
	return _SCA.Contract.ExecuteTransactionFromEntryPoint(&_SCA.TransactOpts, to, value, deadline, data)
}

func (_SCA *SCATransactorSession) ExecuteTransactionFromEntryPoint(to common.Address, value *big.Int, deadline *big.Int, data []byte) (*types.Transaction, error) {
	return _SCA.Contract.ExecuteTransactionFromEntryPoint(&_SCA.TransactOpts, to, value, deadline, data)
}

func (_SCA *SCATransactor) ValidateUserOp(opts *bind.TransactOpts, userOp UserOperation, userOpHash [32]byte, arg2 *big.Int) (*types.Transaction, error) {
	return _SCA.contract.Transact(opts, "validateUserOp", userOp, userOpHash, arg2)
}

func (_SCA *SCASession) ValidateUserOp(userOp UserOperation, userOpHash [32]byte, arg2 *big.Int) (*types.Transaction, error) {
	return _SCA.Contract.ValidateUserOp(&_SCA.TransactOpts, userOp, userOpHash, arg2)
}

func (_SCA *SCATransactorSession) ValidateUserOp(userOp UserOperation, userOpHash [32]byte, arg2 *big.Int) (*types.Transaction, error) {
	return _SCA.Contract.ValidateUserOp(&_SCA.TransactOpts, userOp, userOpHash, arg2)
}

func (_SCA *SCA) Address() common.Address {
	return _SCA.address
}

type SCAInterface interface {
	IEntryPoint(opts *bind.CallOpts) (common.Address, error)

	IOwner(opts *bind.CallOpts) (common.Address, error)

	SNonce(opts *bind.CallOpts) (*big.Int, error)

	ExecuteTransactionFromEntryPoint(opts *bind.TransactOpts, to common.Address, value *big.Int, deadline *big.Int, data []byte) (*types.Transaction, error)

	ValidateUserOp(opts *bind.TransactOpts, userOp UserOperation, userOpHash [32]byte, arg2 *big.Int) (*types.Transaction, error)

	Address() common.Address
}
