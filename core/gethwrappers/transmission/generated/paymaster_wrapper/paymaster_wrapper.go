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
	ABI: "[{\"inputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"linkEthFeed\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"juelsNeeded\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"subscriptionBalance\",\"type\":\"uint256\"}],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"userOpHash\",\"type\":\"bytes32\"}],\"name\":\"UserOperationAlreadyTried\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_ETH_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumIPaymaster.PostOpMode\",\"name\":\"\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"context\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"actualGasCost\",\"type\":\"uint256\"}],\"name\":\"postOp\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_config\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"verificationGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"preVerificationGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxPriorityFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"paymasterAndData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structUserOperation\",\"name\":\"userOp\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"userOpHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"maxCost\",\"type\":\"uint256\"}],\"name\":\"validatePaymasterUserOp\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"context\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"validationData\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c060405234801561001057600080fd5b50604051610edf380380610edf83398101604081905261002f9161005e565b6001600160a01b039182166080521660a052610098565b6001600160a01b038116811461005b57600080fd5b50565b6000806040838503121561007157600080fd5b825161007c81610046565b602084015190925061008d81610046565b809150509250929050565b60805160a051610e076100d86000396000818161018c015261081d01526000818160bc015281816101e70152818161063001526106e80152610e076000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c8063a4c0ed361161005b578063a4c0ed3614610161578063a9a2340914610174578063ad17836114610187578063f465c77e146101ae57600080fd5b8063088070f5146100825780631b6b6d23146100b75780638a38f36514610103575b600080fd5b6000546001546100969163ffffffff169082565b6040805163ffffffff90931683526020830191909152015b60405180910390f35b6100de7f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100ae565b61015f6101113660046108c3565b6040805180820190915263ffffffff9092168083526020909201819052600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000016909217909155600155565b005b61015f61016f366004610966565b6101cf565b61015f6101823660046109c2565b6102cc565b6100de7f000000000000000000000000000000000000000000000000000000000000000081565b6101c16101bc366004610a22565b610335565b6040516100ae929190610a76565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461023e576040517f44b0e3c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208114610278576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061028682840184610af1565b73ffffffffffffffffffffffffffffffffffffffff81166000908152600360205260408120805492935086929091906102c0908490610b44565b90915550505050505050565b6000806102db84860186610b5c565b91509150806102e984610525565b6102f39190610b44565b73ffffffffffffffffffffffffffffffffffffffff831660009081526003602052604081208054909190610328908490610b7a565b9091555050505050505050565b6000828152600260205260408120546060919060ff161561038a576040517f7413dcf8000000000000000000000000000000000000000000000000000000008152600481018590526024015b60405180910390fd5b600061039586610554565b90506000816103a386610525565b6103ad9190610b44565b905080600360006103c160208b018b610af1565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054101561048c57806003600061041460208b018b610af1565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546040517f03eb8b54000000000000000000000000000000000000000000000000000000008152600401610381929190918252602082015260400190565b600086815260026020908152604090912080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790556104d290880188610af1565b6040805173ffffffffffffffffffffffffffffffffffffffff909216602083015281018390526060016040516020818303038152906040526105176000806000610790565b935093505050935093915050565b6000806105306107c8565b61054284670de0b6b3a7640000610b91565b61054c9190610bce565b905050919050565b6000610564610120830183610c09565b905060140361057557506000919050565b6000610585610120840184610c09565b601481811061059657610596610c6e565b919091013560f81c915081905061078a5760006105b7610120850185610c09565b6105c5916015908290610c9d565b8101906105d29190610cc7565b9050806020015160001415801561069f5750602081015181516040517f70a0823100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201527f0000000000000000000000000000000000000000000000000000000000000000909116906370a0823190602401602060405180830381865afa158015610679573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061069d9190610d50565b105b1561078857805160408083015190517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169263a9059cbb9261073c9260040173ffffffffffffffffffffffffffffffffffffffff929092168252602082015260400190565b6020604051808303816000875af115801561075b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061077f9190610d69565b50806040015192505b505b50919050565b600060d08265ffffffffffff16901b60a08465ffffffffffff16901b856107b85760006107bb565b60015b60ff161717949350505050565b60008054604080517ffeaf968c000000000000000000000000000000000000000000000000000000008152905163ffffffff90921691821515918491829173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169163feaf968c9160048082019260a0929091908290030181865afa158015610869573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061088d9190610daa565b5094509092508491505080156108b157506108a88242610b7a565b8463ffffffff16105b156108bb57506001545b949350505050565b600080604083850312156108d657600080fd5b823563ffffffff811681146108ea57600080fd5b946020939093013593505050565b73ffffffffffffffffffffffffffffffffffffffff8116811461091a57600080fd5b50565b60008083601f84011261092f57600080fd5b50813567ffffffffffffffff81111561094757600080fd5b60208301915083602082850101111561095f57600080fd5b9250929050565b6000806000806060858703121561097c57600080fd5b8435610987816108f8565b935060208501359250604085013567ffffffffffffffff8111156109aa57600080fd5b6109b68782880161091d565b95989497509550505050565b600080600080606085870312156109d857600080fd5b8435600381106109e757600080fd5b9350602085013567ffffffffffffffff811115610a0357600080fd5b610a0f8782880161091d565b9598909750949560400135949350505050565b600080600060608486031215610a3757600080fd5b833567ffffffffffffffff811115610a4e57600080fd5b84016101608187031215610a6157600080fd5b95602085013595506040909401359392505050565b604081526000835180604084015260005b81811015610aa45760208187018101516060868401015201610a87565b81811115610ab6576000606083860101525b50602083019390935250601f919091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01601606001919050565b600060208284031215610b0357600080fd5b8135610b0e816108f8565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60008219821115610b5757610b57610b15565b500190565b60008060408385031215610b6f57600080fd5b82356108ea816108f8565b600082821015610b8c57610b8c610b15565b500390565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615610bc957610bc9610b15565b500290565b600082610c04577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610c3e57600080fd5b83018035915067ffffffffffffffff821115610c5957600080fd5b60200191503681900382131561095f57600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60008085851115610cad57600080fd5b83861115610cba57600080fd5b5050820193919092039150565b600060608284031215610cd957600080fd5b6040516060810181811067ffffffffffffffff82111715610d23577f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040528235610d31816108f8565b8152602083810135908201526040928301359281019290925250919050565b600060208284031215610d6257600080fd5b5051919050565b600060208284031215610d7b57600080fd5b81518015158114610b0e57600080fd5b805169ffffffffffffffffffff81168114610da557600080fd5b919050565b600080600080600060a08688031215610dc257600080fd5b610dcb86610d8b565b9450602086015193506040860151925060608601519150610dee60808701610d8b565b9050929550929590935056fea164736f6c634300080f000a",
}

var PaymasterABI = PaymasterMetaData.ABI

var PaymasterBin = PaymasterMetaData.Bin

func DeployPaymaster(auth *bind.TransactOpts, backend bind.ContractBackend, linkToken common.Address, linkEthFeed common.Address) (common.Address, *types.Transaction, *Paymaster, error) {
	parsed, err := PaymasterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PaymasterBin), backend, linkToken, linkEthFeed)
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

func (_Paymaster *PaymasterCaller) LINKETHFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Paymaster.contract.Call(opts, &out, "LINK_ETH_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_Paymaster *PaymasterSession) LINKETHFEED() (common.Address, error) {
	return _Paymaster.Contract.LINKETHFEED(&_Paymaster.CallOpts)
}

func (_Paymaster *PaymasterCallerSession) LINKETHFEED() (common.Address, error) {
	return _Paymaster.Contract.LINKETHFEED(&_Paymaster.CallOpts)
}

func (_Paymaster *PaymasterCaller) SConfig(opts *bind.CallOpts) (SConfig,

	error) {
	var out []interface{}
	err := _Paymaster.contract.Call(opts, &out, "s_config")

	outstruct := new(SConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.StalenessSeconds = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.FallbackWeiPerUnitLink = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_Paymaster *PaymasterSession) SConfig() (SConfig,

	error) {
	return _Paymaster.Contract.SConfig(&_Paymaster.CallOpts)
}

func (_Paymaster *PaymasterCallerSession) SConfig() (SConfig,

	error) {
	return _Paymaster.Contract.SConfig(&_Paymaster.CallOpts)
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

func (_Paymaster *PaymasterTransactor) SetConfig(opts *bind.TransactOpts, stalenessSeconds uint32, fallbackWeiPerUnitLink *big.Int) (*types.Transaction, error) {
	return _Paymaster.contract.Transact(opts, "setConfig", stalenessSeconds, fallbackWeiPerUnitLink)
}

func (_Paymaster *PaymasterSession) SetConfig(stalenessSeconds uint32, fallbackWeiPerUnitLink *big.Int) (*types.Transaction, error) {
	return _Paymaster.Contract.SetConfig(&_Paymaster.TransactOpts, stalenessSeconds, fallbackWeiPerUnitLink)
}

func (_Paymaster *PaymasterTransactorSession) SetConfig(stalenessSeconds uint32, fallbackWeiPerUnitLink *big.Int) (*types.Transaction, error) {
	return _Paymaster.Contract.SetConfig(&_Paymaster.TransactOpts, stalenessSeconds, fallbackWeiPerUnitLink)
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

type SConfig struct {
	StalenessSeconds       uint32
	FallbackWeiPerUnitLink *big.Int
}

func (_Paymaster *Paymaster) Address() common.Address {
	return _Paymaster.address
}

type PaymasterInterface interface {
	LINK(opts *bind.CallOpts) (common.Address, error)

	LINKETHFEED(opts *bind.CallOpts) (common.Address, error)

	SConfig(opts *bind.CallOpts) (SConfig,

		error)

	OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error)

	PostOp(opts *bind.TransactOpts, arg0 uint8, context []byte, actualGasCost *big.Int) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, stalenessSeconds uint32, fallbackWeiPerUnitLink *big.Int) (*types.Transaction, error)

	ValidatePaymasterUserOp(opts *bind.TransactOpts, userOp UserOperation, userOpHash [32]byte, maxCost *big.Int) (*types.Transaction, error)

	Address() common.Address
}
