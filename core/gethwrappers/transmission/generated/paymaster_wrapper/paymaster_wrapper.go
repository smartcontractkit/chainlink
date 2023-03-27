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
	ABI: "[{\"inputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"linkToken\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"juelsNeeded\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"subscriptionBalance\",\"type\":\"uint256\"}],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"userOpHash\",\"type\":\"bytes32\"}],\"name\":\"UserOperationAlreadyTried\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumIPaymaster.PostOpMode\",\"name\":\"\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"context\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"actualGasCost\",\"type\":\"uint256\"}],\"name\":\"postOp\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"verificationGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"preVerificationGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxPriorityFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"paymasterAndData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structUserOperation\",\"name\":\"userOp\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"userOpHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"maxCost\",\"type\":\"uint256\"}],\"name\":\"validatePaymasterUserOp\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"context\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"validationData\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040526611c37937e0800060005534801561001b57600080fd5b50604051610c40380380610c4083398101604081905261003a9161004b565b6001600160a01b031660805261007b565b60006020828403121561005d57600080fd5b81516001600160a01b038116811461007457600080fd5b9392505050565b608051610b966100aa600039600081816056015281816101030152818161055001526106080152610b966000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80631b6b6d2314610051578063a4c0ed36146100a2578063a9a23409146100b7578063f465c77e146100ca575b600080fd5b6100787f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b6100b56100b0366004610756565b6100eb565b005b6100b56100c53660046107b2565b6101e8565b6100dd6100d8366004610812565b61026c565b604051610099929190610866565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461015a576040517f44b0e3c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208114610194576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006101a2828401846108e1565b73ffffffffffffffffffffffffffffffffffffffff81166000908152600260205260408120805492935086929091906101dc908490610934565b90915550505050505050565b6000806101f78486018661094c565b915091506000805484670de0b6b3a76400006102139190610978565b61021d91906109b5565b90506102298282610934565b73ffffffffffffffffffffffffffffffffffffffff84166000908152600260205260408120805490919061025e9084906109f0565b909155505050505050505050565b6000828152600160205260408120546060919060ff16156102c1576040517f7413dcf8000000000000000000000000000000000000000000000000000000008152600481018590526024015b60405180910390fd5b60006102cc86610474565b905060008160005486670de0b6b3a76400006102e89190610978565b6102f291906109b5565b6102fc9190610934565b9050806002600061031060208b018b6108e1565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205410156103db57806002600061036360208b018b6108e1565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546040517f03eb8b540000000000000000000000000000000000000000000000000000000081526004016102b8929190918252602082015260400190565b600086815260016020818152604090922080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169091179055610421908801886108e1565b6040805173ffffffffffffffffffffffffffffffffffffffff9092166020830152810183905260600160405160208183030381529060405261046660008060006106b0565b935093505050935093915050565b6000610484610120830183610a07565b905060140361049557506000919050565b60006104a5610120840184610a07565b60148181106104b6576104b6610a6c565b919091013560f81c91508190506106aa5760006104d7610120850185610a07565b6104e5916015908290610a9b565b8101906104f29190610ac5565b905080602001516000141580156105bf5750602081015181516040517f70a0823100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201527f0000000000000000000000000000000000000000000000000000000000000000909116906370a0823190602401602060405180830381865afa158015610599573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105bd9190610b4e565b105b156106a857805160408083015190517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169263a9059cbb9261065c9260040173ffffffffffffffffffffffffffffffffffffffff929092168252602082015260400190565b6020604051808303816000875af115801561067b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061069f9190610b67565b50806040015192505b505b50919050565b600060d08265ffffffffffff16901b60a08465ffffffffffff16901b856106d85760006106db565b60015b60ff161717949350505050565b73ffffffffffffffffffffffffffffffffffffffff8116811461070a57600080fd5b50565b60008083601f84011261071f57600080fd5b50813567ffffffffffffffff81111561073757600080fd5b60208301915083602082850101111561074f57600080fd5b9250929050565b6000806000806060858703121561076c57600080fd5b8435610777816106e8565b935060208501359250604085013567ffffffffffffffff81111561079a57600080fd5b6107a68782880161070d565b95989497509550505050565b600080600080606085870312156107c857600080fd5b8435600381106107d757600080fd5b9350602085013567ffffffffffffffff8111156107f357600080fd5b6107ff8782880161070d565b9598909750949560400135949350505050565b60008060006060848603121561082757600080fd5b833567ffffffffffffffff81111561083e57600080fd5b8401610160818703121561085157600080fd5b95602085013595506040909401359392505050565b604081526000835180604084015260005b818110156108945760208187018101516060868401015201610877565b818111156108a6576000606083860101525b50602083019390935250601f919091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01601606001919050565b6000602082840312156108f357600080fd5b81356108fe816106e8565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000821982111561094757610947610905565b500190565b6000806040838503121561095f57600080fd5b823561096a816106e8565b946020939093013593505050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156109b0576109b0610905565b500290565b6000826109eb577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b600082821015610a0257610a02610905565b500390565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610a3c57600080fd5b83018035915067ffffffffffffffff821115610a5757600080fd5b60200191503681900382131561074f57600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60008085851115610aab57600080fd5b83861115610ab857600080fd5b5050820193919092039150565b600060608284031215610ad757600080fd5b6040516060810181811067ffffffffffffffff82111715610b21577f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040528235610b2f816106e8565b8152602083810135908201526040928301359281019290925250919050565b600060208284031215610b6057600080fd5b5051919050565b600060208284031215610b7957600080fd5b815180151581146108fe57600080fdfea164736f6c634300080f000a",
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
