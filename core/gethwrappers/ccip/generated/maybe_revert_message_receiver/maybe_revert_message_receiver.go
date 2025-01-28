// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package maybe_revert_message_receiver

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

type ClientAny2EVMMessage struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	Sender              []byte
	Data                []byte
	DestTokenAmounts    []ClientEVMTokenAmount
}

type ClientEVMTokenAmount struct {
	Token  common.Address
	Amount *big.Int
}

var MaybeRevertMessageReceiverMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"toRevert\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"balanceOfToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ccipReceive\",\"inputs\":[{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structClient.Any2EVMMessage\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"s_toRevert\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setErr\",\"inputs\":[{\"name\":\"err\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRevert\",\"inputs\":[{\"name\":\"toRevert\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"withdrawFunds\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawTokens\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"MessageReceived\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"sender\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\",\"indexed\":false,\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NativeFundsWithdrawn\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokensWithdrawn\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValueReceived\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CustomError\",\"inputs\":[{\"name\":\"err\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InsufficientBalance\",\"inputs\":[{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"required\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ReceiveRevert\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TransferFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[]}]",
	Bin: "0x60a034607657601f610dbb38819003918201601f19168301916001600160401b03831184841017607b57808492602094604052833981010312607657518015158091036076573360805260ff801960005416911617600055604051610d29908161009282396080518181816107e301526109300152f35b600080fd5b634e487b7160e01b600052604160045260246000fdfe608080604052600436101561007b575b50361561001b57600080fd5b60ff60005416610051577fe12e3b7047ff60a2dd763cf536a43597e5ce7fe7aa7476345bd4cd079912bcef6020604051348152a1005b7f3085b8db0000000000000000000000000000000000000000000000000000000060005260046000fd5b60003560e01c90816301ffc9a714610adf5750806306b091f9146108df57806324600fc3146107b25780635100fc211461077157806377f5b0e6146104db57806385572ffb1461022e5780638fb5f171146101c05763b99152d0146100e0573861000f565b346101a75760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101a7576024602073ffffffffffffffffffffffffffffffffffffffff610130610b9b565b16604051928380927f70a082310000000000000000000000000000000000000000000000000000000082523060048301525afa80156101b45760009061017c575b602090604051908152f35b506020813d6020116101ac575b8161019660209383610bbe565b810103126101a75760209051610171565b600080fd5b3d9150610189565b6040513d6000823e3d90fd5b346101a75760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101a7576004358015158091036101a75760ff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0060005416911617600055600080f35b346101a75760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101a75760043567ffffffffffffffff81116101a7578060040181360360a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8201126101a75760ff600054166103ef5760248301359267ffffffffffffffff84168094036101a7576102cf6044820184610c8c565b90916102de6064820186610c8c565b9490917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdd608482013592018212156101a757019360048501359467ffffffffffffffff86116101a757602401968560061b360388136101a75760209461036c9461035e9260405199358a52878a015260a060408a015260a0890191610cdd565b918683036060880152610cdd565b83810360808501528281520192906000905b8082106103ad577f707732b700184c0ab3b799f43f03de9b3606a144cfb367f98291044e71972cdc84860385a1005b90919384359073ffffffffffffffffffffffffffffffffffffffff82168092036101a7576040816001938293526020880135602082015201950192019061037e565b6040517f5a4ff6710000000000000000000000000000000000000000000000000000000081526020600482015280600060015461042b81610c39565b908160248501526001811690816000146104a3575060011461044c57500390fd5b6001600090815291507fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf65b81831061048957505081010360440190fd5b805460448487010152849350602090920191600101610477565b604493507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091501682840152151560051b8201010390fd5b346101a75760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101a75760043567ffffffffffffffff81116101a757366023820112156101a757806004013561053581610bff565b916105436040519384610bbe565b81835236602483830101116101a757816000926024602093018386013783010152805167ffffffffffffffff811161074257610580600154610c39565b601f811161069f575b50602091601f82116001146105e5579181926000926105da575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c191617600155600080f35b0151905082806105a3565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082169260016000527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf69160005b85811061068757508360019510610650575b505050811b01600155005b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c19169055828080610645565b91926020600181928685015181550194019201610633565b6001600052601f820160051c7fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf601906020831061071a575b601f0160051c7fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf601905b81811061070e5750610589565b60008155600101610701565b7fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf691506106d7565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b346101a75760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101a757602060ff600054166040519015158152f35b346101a75760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101a7577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff8116908133036108b55760008080804780955af13d156108b0573d61083b81610bff565b906108496040519283610bbe565b8152600060203d92013e5b156108865760207fd50b71a2790ecccf5881141fe9ae079e17c66aace5d50ba383d443ecd398ffa591604051908152a2005b7f90b8ec180000000000000000000000000000000000000000000000000000000060005260046000fd5b610854565b7f82b429000000000000000000000000000000000000000000000000000000000060005260046000fd5b346101a75760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101a757610916610b9b565b60243573ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016918233036108b55773ffffffffffffffffffffffffffffffffffffffff16906040517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152602081602481865afa9081156101b457600091610aad575b50818110610a7d57506040517fa9059cbb0000000000000000000000000000000000000000000000000000000081528360048201528160248201526020816044816000875af19081156101b457600091610a3b575b50156108865760207f6337ed398c0e8467698c581374fdce4db14922df487b5a39483079f5f59b60a491604051908152a3005b6020813d602011610a75575b81610a5460209383610bbe565b81010312610a715751908115158203610a6e575084610a08565b80fd5b5080fd5b3d9150610a47565b7fcf4791810000000000000000000000000000000000000000000000000000000060005260045260245260446000fd5b906020823d602011610ad7575b81610ac760209383610bbe565b81010312610a6e575051846109b3565b3d9150610aba565b346101a75760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101a757600435907fffffffff0000000000000000000000000000000000000000000000000000000082168092036101a757817f85572ffb0000000000000000000000000000000000000000000000000000000060209314908115610b71575b5015158152f35b7f01ffc9a70000000000000000000000000000000000000000000000000000000091501483610b6a565b6004359073ffffffffffffffffffffffffffffffffffffffff821682036101a757565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff82111761074257604052565b67ffffffffffffffff811161074257601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b90600182811c92168015610c82575b6020831014610c5357565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b91607f1691610c48565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1813603018212156101a7570180359067ffffffffffffffff82116101a7576020019181360383136101a757565b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe093818652868601376000858286010152011601019056fea164736f6c634300081a000a",
}

var MaybeRevertMessageReceiverABI = MaybeRevertMessageReceiverMetaData.ABI

var MaybeRevertMessageReceiverBin = MaybeRevertMessageReceiverMetaData.Bin

func DeployMaybeRevertMessageReceiver(auth *bind.TransactOpts, backend bind.ContractBackend, toRevert bool) (common.Address, *types.Transaction, *MaybeRevertMessageReceiver, error) {
	parsed, err := MaybeRevertMessageReceiverMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MaybeRevertMessageReceiverBin), backend, toRevert)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MaybeRevertMessageReceiver{address: address, abi: *parsed, MaybeRevertMessageReceiverCaller: MaybeRevertMessageReceiverCaller{contract: contract}, MaybeRevertMessageReceiverTransactor: MaybeRevertMessageReceiverTransactor{contract: contract}, MaybeRevertMessageReceiverFilterer: MaybeRevertMessageReceiverFilterer{contract: contract}}, nil
}

type MaybeRevertMessageReceiver struct {
	address common.Address
	abi     abi.ABI
	MaybeRevertMessageReceiverCaller
	MaybeRevertMessageReceiverTransactor
	MaybeRevertMessageReceiverFilterer
}

type MaybeRevertMessageReceiverCaller struct {
	contract *bind.BoundContract
}

type MaybeRevertMessageReceiverTransactor struct {
	contract *bind.BoundContract
}

type MaybeRevertMessageReceiverFilterer struct {
	contract *bind.BoundContract
}

type MaybeRevertMessageReceiverSession struct {
	Contract     *MaybeRevertMessageReceiver
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MaybeRevertMessageReceiverCallerSession struct {
	Contract *MaybeRevertMessageReceiverCaller
	CallOpts bind.CallOpts
}

type MaybeRevertMessageReceiverTransactorSession struct {
	Contract     *MaybeRevertMessageReceiverTransactor
	TransactOpts bind.TransactOpts
}

type MaybeRevertMessageReceiverRaw struct {
	Contract *MaybeRevertMessageReceiver
}

type MaybeRevertMessageReceiverCallerRaw struct {
	Contract *MaybeRevertMessageReceiverCaller
}

type MaybeRevertMessageReceiverTransactorRaw struct {
	Contract *MaybeRevertMessageReceiverTransactor
}

func NewMaybeRevertMessageReceiver(address common.Address, backend bind.ContractBackend) (*MaybeRevertMessageReceiver, error) {
	abi, err := abi.JSON(strings.NewReader(MaybeRevertMessageReceiverABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMaybeRevertMessageReceiver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MaybeRevertMessageReceiver{address: address, abi: abi, MaybeRevertMessageReceiverCaller: MaybeRevertMessageReceiverCaller{contract: contract}, MaybeRevertMessageReceiverTransactor: MaybeRevertMessageReceiverTransactor{contract: contract}, MaybeRevertMessageReceiverFilterer: MaybeRevertMessageReceiverFilterer{contract: contract}}, nil
}

func NewMaybeRevertMessageReceiverCaller(address common.Address, caller bind.ContractCaller) (*MaybeRevertMessageReceiverCaller, error) {
	contract, err := bindMaybeRevertMessageReceiver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MaybeRevertMessageReceiverCaller{contract: contract}, nil
}

func NewMaybeRevertMessageReceiverTransactor(address common.Address, transactor bind.ContractTransactor) (*MaybeRevertMessageReceiverTransactor, error) {
	contract, err := bindMaybeRevertMessageReceiver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MaybeRevertMessageReceiverTransactor{contract: contract}, nil
}

func NewMaybeRevertMessageReceiverFilterer(address common.Address, filterer bind.ContractFilterer) (*MaybeRevertMessageReceiverFilterer, error) {
	contract, err := bindMaybeRevertMessageReceiver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MaybeRevertMessageReceiverFilterer{contract: contract}, nil
}

func bindMaybeRevertMessageReceiver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MaybeRevertMessageReceiverMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MaybeRevertMessageReceiver.Contract.MaybeRevertMessageReceiverCaller.contract.Call(opts, result, method, params...)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.MaybeRevertMessageReceiverTransactor.contract.Transfer(opts)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.MaybeRevertMessageReceiverTransactor.contract.Transact(opts, method, params...)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MaybeRevertMessageReceiver.Contract.contract.Call(opts, result, method, params...)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.contract.Transfer(opts)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.contract.Transact(opts, method, params...)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverCaller) BalanceOfToken(opts *bind.CallOpts, token common.Address) (*big.Int, error) {
	var out []interface{}
	err := _MaybeRevertMessageReceiver.contract.Call(opts, &out, "balanceOfToken", token)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverSession) BalanceOfToken(token common.Address) (*big.Int, error) {
	return _MaybeRevertMessageReceiver.Contract.BalanceOfToken(&_MaybeRevertMessageReceiver.CallOpts, token)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverCallerSession) BalanceOfToken(token common.Address) (*big.Int, error) {
	return _MaybeRevertMessageReceiver.Contract.BalanceOfToken(&_MaybeRevertMessageReceiver.CallOpts, token)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverCaller) SToRevert(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _MaybeRevertMessageReceiver.contract.Call(opts, &out, "s_toRevert")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverSession) SToRevert() (bool, error) {
	return _MaybeRevertMessageReceiver.Contract.SToRevert(&_MaybeRevertMessageReceiver.CallOpts)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverCallerSession) SToRevert() (bool, error) {
	return _MaybeRevertMessageReceiver.Contract.SToRevert(&_MaybeRevertMessageReceiver.CallOpts)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _MaybeRevertMessageReceiver.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _MaybeRevertMessageReceiver.Contract.SupportsInterface(&_MaybeRevertMessageReceiver.CallOpts, interfaceId)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _MaybeRevertMessageReceiver.Contract.SupportsInterface(&_MaybeRevertMessageReceiver.CallOpts, interfaceId)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverTransactor) CcipReceive(opts *bind.TransactOpts, message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.contract.Transact(opts, "ccipReceive", message)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverSession) CcipReceive(message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.CcipReceive(&_MaybeRevertMessageReceiver.TransactOpts, message)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverTransactorSession) CcipReceive(message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.CcipReceive(&_MaybeRevertMessageReceiver.TransactOpts, message)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverTransactor) SetErr(opts *bind.TransactOpts, err []byte) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.contract.Transact(opts, "setErr", err)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverSession) SetErr(err []byte) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.SetErr(&_MaybeRevertMessageReceiver.TransactOpts, err)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverTransactorSession) SetErr(err []byte) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.SetErr(&_MaybeRevertMessageReceiver.TransactOpts, err)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverTransactor) SetRevert(opts *bind.TransactOpts, toRevert bool) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.contract.Transact(opts, "setRevert", toRevert)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverSession) SetRevert(toRevert bool) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.SetRevert(&_MaybeRevertMessageReceiver.TransactOpts, toRevert)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverTransactorSession) SetRevert(toRevert bool) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.SetRevert(&_MaybeRevertMessageReceiver.TransactOpts, toRevert)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverTransactor) WithdrawFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.contract.Transact(opts, "withdrawFunds")
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverSession) WithdrawFunds() (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.WithdrawFunds(&_MaybeRevertMessageReceiver.TransactOpts)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverTransactorSession) WithdrawFunds() (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.WithdrawFunds(&_MaybeRevertMessageReceiver.TransactOpts)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverTransactor) WithdrawTokens(opts *bind.TransactOpts, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.contract.Transact(opts, "withdrawTokens", token, amount)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverSession) WithdrawTokens(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.WithdrawTokens(&_MaybeRevertMessageReceiver.TransactOpts, token, amount)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverTransactorSession) WithdrawTokens(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.WithdrawTokens(&_MaybeRevertMessageReceiver.TransactOpts, token, amount)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.contract.RawTransact(opts, nil)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverSession) Receive() (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.Receive(&_MaybeRevertMessageReceiver.TransactOpts)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverTransactorSession) Receive() (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.Receive(&_MaybeRevertMessageReceiver.TransactOpts)
}

type MaybeRevertMessageReceiverMessageReceivedIterator struct {
	Event *MaybeRevertMessageReceiverMessageReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MaybeRevertMessageReceiverMessageReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MaybeRevertMessageReceiverMessageReceived)
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
		it.Event = new(MaybeRevertMessageReceiverMessageReceived)
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

func (it *MaybeRevertMessageReceiverMessageReceivedIterator) Error() error {
	return it.fail
}

func (it *MaybeRevertMessageReceiverMessageReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MaybeRevertMessageReceiverMessageReceived struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	Sender              []byte
	Data                []byte
	DestTokenAmounts    []ClientEVMTokenAmount
	Raw                 types.Log
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverFilterer) FilterMessageReceived(opts *bind.FilterOpts) (*MaybeRevertMessageReceiverMessageReceivedIterator, error) {

	logs, sub, err := _MaybeRevertMessageReceiver.contract.FilterLogs(opts, "MessageReceived")
	if err != nil {
		return nil, err
	}
	return &MaybeRevertMessageReceiverMessageReceivedIterator{contract: _MaybeRevertMessageReceiver.contract, event: "MessageReceived", logs: logs, sub: sub}, nil
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverFilterer) WatchMessageReceived(opts *bind.WatchOpts, sink chan<- *MaybeRevertMessageReceiverMessageReceived) (event.Subscription, error) {

	logs, sub, err := _MaybeRevertMessageReceiver.contract.WatchLogs(opts, "MessageReceived")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MaybeRevertMessageReceiverMessageReceived)
				if err := _MaybeRevertMessageReceiver.contract.UnpackLog(event, "MessageReceived", log); err != nil {
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

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverFilterer) ParseMessageReceived(log types.Log) (*MaybeRevertMessageReceiverMessageReceived, error) {
	event := new(MaybeRevertMessageReceiverMessageReceived)
	if err := _MaybeRevertMessageReceiver.contract.UnpackLog(event, "MessageReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MaybeRevertMessageReceiverNativeFundsWithdrawnIterator struct {
	Event *MaybeRevertMessageReceiverNativeFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MaybeRevertMessageReceiverNativeFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MaybeRevertMessageReceiverNativeFundsWithdrawn)
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
		it.Event = new(MaybeRevertMessageReceiverNativeFundsWithdrawn)
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

func (it *MaybeRevertMessageReceiverNativeFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *MaybeRevertMessageReceiverNativeFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MaybeRevertMessageReceiverNativeFundsWithdrawn struct {
	Owner  common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverFilterer) FilterNativeFundsWithdrawn(opts *bind.FilterOpts, owner []common.Address) (*MaybeRevertMessageReceiverNativeFundsWithdrawnIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _MaybeRevertMessageReceiver.contract.FilterLogs(opts, "NativeFundsWithdrawn", ownerRule)
	if err != nil {
		return nil, err
	}
	return &MaybeRevertMessageReceiverNativeFundsWithdrawnIterator{contract: _MaybeRevertMessageReceiver.contract, event: "NativeFundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverFilterer) WatchNativeFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *MaybeRevertMessageReceiverNativeFundsWithdrawn, owner []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _MaybeRevertMessageReceiver.contract.WatchLogs(opts, "NativeFundsWithdrawn", ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MaybeRevertMessageReceiverNativeFundsWithdrawn)
				if err := _MaybeRevertMessageReceiver.contract.UnpackLog(event, "NativeFundsWithdrawn", log); err != nil {
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

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverFilterer) ParseNativeFundsWithdrawn(log types.Log) (*MaybeRevertMessageReceiverNativeFundsWithdrawn, error) {
	event := new(MaybeRevertMessageReceiverNativeFundsWithdrawn)
	if err := _MaybeRevertMessageReceiver.contract.UnpackLog(event, "NativeFundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MaybeRevertMessageReceiverTokensWithdrawnIterator struct {
	Event *MaybeRevertMessageReceiverTokensWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MaybeRevertMessageReceiverTokensWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MaybeRevertMessageReceiverTokensWithdrawn)
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
		it.Event = new(MaybeRevertMessageReceiverTokensWithdrawn)
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

func (it *MaybeRevertMessageReceiverTokensWithdrawnIterator) Error() error {
	return it.fail
}

func (it *MaybeRevertMessageReceiverTokensWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MaybeRevertMessageReceiverTokensWithdrawn struct {
	Token  common.Address
	Owner  common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverFilterer) FilterTokensWithdrawn(opts *bind.FilterOpts, token []common.Address, owner []common.Address) (*MaybeRevertMessageReceiverTokensWithdrawnIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _MaybeRevertMessageReceiver.contract.FilterLogs(opts, "TokensWithdrawn", tokenRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &MaybeRevertMessageReceiverTokensWithdrawnIterator{contract: _MaybeRevertMessageReceiver.contract, event: "TokensWithdrawn", logs: logs, sub: sub}, nil
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverFilterer) WatchTokensWithdrawn(opts *bind.WatchOpts, sink chan<- *MaybeRevertMessageReceiverTokensWithdrawn, token []common.Address, owner []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _MaybeRevertMessageReceiver.contract.WatchLogs(opts, "TokensWithdrawn", tokenRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MaybeRevertMessageReceiverTokensWithdrawn)
				if err := _MaybeRevertMessageReceiver.contract.UnpackLog(event, "TokensWithdrawn", log); err != nil {
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

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverFilterer) ParseTokensWithdrawn(log types.Log) (*MaybeRevertMessageReceiverTokensWithdrawn, error) {
	event := new(MaybeRevertMessageReceiverTokensWithdrawn)
	if err := _MaybeRevertMessageReceiver.contract.UnpackLog(event, "TokensWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MaybeRevertMessageReceiverValueReceivedIterator struct {
	Event *MaybeRevertMessageReceiverValueReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MaybeRevertMessageReceiverValueReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MaybeRevertMessageReceiverValueReceived)
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
		it.Event = new(MaybeRevertMessageReceiverValueReceived)
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

func (it *MaybeRevertMessageReceiverValueReceivedIterator) Error() error {
	return it.fail
}

func (it *MaybeRevertMessageReceiverValueReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MaybeRevertMessageReceiverValueReceived struct {
	Amount *big.Int
	Raw    types.Log
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverFilterer) FilterValueReceived(opts *bind.FilterOpts) (*MaybeRevertMessageReceiverValueReceivedIterator, error) {

	logs, sub, err := _MaybeRevertMessageReceiver.contract.FilterLogs(opts, "ValueReceived")
	if err != nil {
		return nil, err
	}
	return &MaybeRevertMessageReceiverValueReceivedIterator{contract: _MaybeRevertMessageReceiver.contract, event: "ValueReceived", logs: logs, sub: sub}, nil
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverFilterer) WatchValueReceived(opts *bind.WatchOpts, sink chan<- *MaybeRevertMessageReceiverValueReceived) (event.Subscription, error) {

	logs, sub, err := _MaybeRevertMessageReceiver.contract.WatchLogs(opts, "ValueReceived")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MaybeRevertMessageReceiverValueReceived)
				if err := _MaybeRevertMessageReceiver.contract.UnpackLog(event, "ValueReceived", log); err != nil {
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

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverFilterer) ParseValueReceived(log types.Log) (*MaybeRevertMessageReceiverValueReceived, error) {
	event := new(MaybeRevertMessageReceiverValueReceived)
	if err := _MaybeRevertMessageReceiver.contract.UnpackLog(event, "ValueReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiver) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _MaybeRevertMessageReceiver.abi.Events["MessageReceived"].ID:
		return _MaybeRevertMessageReceiver.ParseMessageReceived(log)
	case _MaybeRevertMessageReceiver.abi.Events["NativeFundsWithdrawn"].ID:
		return _MaybeRevertMessageReceiver.ParseNativeFundsWithdrawn(log)
	case _MaybeRevertMessageReceiver.abi.Events["TokensWithdrawn"].ID:
		return _MaybeRevertMessageReceiver.ParseTokensWithdrawn(log)
	case _MaybeRevertMessageReceiver.abi.Events["ValueReceived"].ID:
		return _MaybeRevertMessageReceiver.ParseValueReceived(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (MaybeRevertMessageReceiverMessageReceived) Topic() common.Hash {
	return common.HexToHash("0x707732b700184c0ab3b799f43f03de9b3606a144cfb367f98291044e71972cdc")
}

func (MaybeRevertMessageReceiverNativeFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xd50b71a2790ecccf5881141fe9ae079e17c66aace5d50ba383d443ecd398ffa5")
}

func (MaybeRevertMessageReceiverTokensWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x6337ed398c0e8467698c581374fdce4db14922df487b5a39483079f5f59b60a4")
}

func (MaybeRevertMessageReceiverValueReceived) Topic() common.Hash {
	return common.HexToHash("0xe12e3b7047ff60a2dd763cf536a43597e5ce7fe7aa7476345bd4cd079912bcef")
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiver) Address() common.Address {
	return _MaybeRevertMessageReceiver.address
}

type MaybeRevertMessageReceiverInterface interface {
	BalanceOfToken(opts *bind.CallOpts, token common.Address) (*big.Int, error)

	SToRevert(opts *bind.CallOpts) (bool, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	CcipReceive(opts *bind.TransactOpts, message ClientAny2EVMMessage) (*types.Transaction, error)

	SetErr(opts *bind.TransactOpts, err []byte) (*types.Transaction, error)

	SetRevert(opts *bind.TransactOpts, toRevert bool) (*types.Transaction, error)

	WithdrawFunds(opts *bind.TransactOpts) (*types.Transaction, error)

	WithdrawTokens(opts *bind.TransactOpts, token common.Address, amount *big.Int) (*types.Transaction, error)

	Receive(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterMessageReceived(opts *bind.FilterOpts) (*MaybeRevertMessageReceiverMessageReceivedIterator, error)

	WatchMessageReceived(opts *bind.WatchOpts, sink chan<- *MaybeRevertMessageReceiverMessageReceived) (event.Subscription, error)

	ParseMessageReceived(log types.Log) (*MaybeRevertMessageReceiverMessageReceived, error)

	FilterNativeFundsWithdrawn(opts *bind.FilterOpts, owner []common.Address) (*MaybeRevertMessageReceiverNativeFundsWithdrawnIterator, error)

	WatchNativeFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *MaybeRevertMessageReceiverNativeFundsWithdrawn, owner []common.Address) (event.Subscription, error)

	ParseNativeFundsWithdrawn(log types.Log) (*MaybeRevertMessageReceiverNativeFundsWithdrawn, error)

	FilterTokensWithdrawn(opts *bind.FilterOpts, token []common.Address, owner []common.Address) (*MaybeRevertMessageReceiverTokensWithdrawnIterator, error)

	WatchTokensWithdrawn(opts *bind.WatchOpts, sink chan<- *MaybeRevertMessageReceiverTokensWithdrawn, token []common.Address, owner []common.Address) (event.Subscription, error)

	ParseTokensWithdrawn(log types.Log) (*MaybeRevertMessageReceiverTokensWithdrawn, error)

	FilterValueReceived(opts *bind.FilterOpts) (*MaybeRevertMessageReceiverValueReceivedIterator, error)

	WatchValueReceived(opts *bind.WatchOpts, sink chan<- *MaybeRevertMessageReceiverValueReceived) (event.Subscription, error)

	ParseValueReceived(log types.Log) (*MaybeRevertMessageReceiverValueReceived, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
