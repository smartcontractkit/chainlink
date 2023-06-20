// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrfv2plus_consumer_example

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

var VRFV2PlusConsumerExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"subOwner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlySubOwnerCanSetVRFCoordinator\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"name\":\"redeemRandomness\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"nativePayment\",\"type\":\"bool\"}],\"name\":\"requestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_linkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requests\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"cost\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_vrfCoordinator\",\"outputs\":[{\"internalType\":\"contractIVRFCoordinatorV2Plus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"}],\"name\":\"setVRFCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b506040516200136d3803806200136d833981016040819052620000349162000270565b828133806000816200008d5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c057620000c081620001a7565b5050506001600160a01b0382166200010a5760405162461bcd60e51b815260206004820152600c60248201526b7a65726f206164647265737360a01b604482015260640162000084565b6001600160a01b038116620001515760405162461bcd60e51b815260206004820152600c60248201526b7a65726f206164647265737360a01b604482015260640162000084565b600280546001600160a01b03199081166001600160a01b039485161790915560038054821692841692909217909155600480548216958316959095179094556005805490941692169190911790915550620002ba565b6001600160a01b038116331415620002025760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000084565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03811681146200026b57600080fd5b919050565b6000806000606084860312156200028657600080fd5b620002918462000253565b9250620002a16020850162000253565b9150620002b16040850162000253565b90509250925092565b6110a380620002ca6000396000f3fe6080604052600436106100b15760003560e01c80639c7450ff11610069578063a168fa891161004e578063a168fa89146101dc578063f2fde38b14610275578063f793a70e1461029557600080fd5b80639c7450ff1461018f5780639eccacf6146101af57600080fd5b80637725135b1161009a5780637725135b146100f857806379ba50971461014f5780638da5cb5b1461016457600080fd5b80631fe543e3146100b657806344ff81ce146100d8575b600080fd5b3480156100c257600080fd5b506100d66100d1366004610e78565b6102b5565b005b3480156100e457600080fd5b506100d66100f3366004610e09565b61033b565b34801561010457600080fd5b506005546101259073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b34801561015b57600080fd5b506100d66103f5565b34801561017057600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff16610125565b6101a261019d366004610e46565b6104f2565b6040516101469190611023565b3480156101bb57600080fd5b506004546101259073ffffffffffffffffffffffffffffffffffffffff1681565b3480156101e857600080fd5b506102396101f7366004610e46565b600660205260009081526040902080546002820154600390920154909160ff81169161010090910473ffffffffffffffffffffffffffffffffffffffff169084565b6040516101469493929190938452911515602084015273ffffffffffffffffffffffffffffffffffffffff166040830152606082015260800190565b34801561028157600080fd5b506100d6610290366004610e09565b61079c565b3480156102a157600080fd5b506100d66102b0366004610f67565b6107b0565b60025473ffffffffffffffffffffffffffffffffffffffff16331461032d576002546040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff90911660248201526044015b60405180910390fd5b61033782826109c5565b5050565b60035473ffffffffffffffffffffffffffffffffffffffff1633146103ae576003546040517f4ae338ff00000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff9091166024820152604401610324565b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60015473ffffffffffffffffffffffffffffffffffffffff163314610476576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610324565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000818152600660209081526040808320815160a08101835281548152600182018054845181870281018701909552808552606096959294858401939092919083018282801561056157602002820191906000526020600020905b81548152602001906001019080831161054d575b5050509183525050600282015460ff811615156020830152610100900473ffffffffffffffffffffffffffffffffffffffff1660408201526003909101546060909101528051909150610610576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265717565737420494420697320696e636f72726563740000000000000000006044820152606401610324565b806040015161067b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f72657175657374206e6f742066756c66696c6c656420796574000000000000006044820152606401610324565b606081015173ffffffffffffffffffffffffffffffffffffffff163314610724576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f6f6e6c792063616c6c61626c652062792072657175657374696e67206164647260448201527f65737300000000000000000000000000000000000000000000000000000000006064820152608401610324565b8060800151341015610792576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e647300000000000000000000000000006044820152606401610324565b6020015192915050565b6107a4610c17565b6107ad81610c9a565b50565b600480546040517fefcf1d9400000000000000000000000000000000000000000000000000000000815291820184905267ffffffffffffffff8816602483015261ffff8616604483015263ffffffff80881660648401528516608483015282151560a483015260009173ffffffffffffffffffffffffffffffffffffffff9091169063efcf1d949060c401602060405180830381600087803b15801561085557600080fd5b505af1158015610869573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061088d9190610e5f565b905060006040518060a00160405280838152602001600067ffffffffffffffff8111156108bc576108bc611067565b6040519080825280602002602001820160405280156108e5578160200160208202803683370190505b5081526000602080830182905233604080850191909152606090930182905285825260068152919020825181558282015180519394508493919261093192600185019290910190610d90565b506040820151600282018054606085015173ffffffffffffffffffffffffffffffffffffffff16610100027fffffffffffffffffffffff0000000000000000000000000000000000000000ff931515939093167fffffffffffffffffffffff000000000000000000000000000000000000000000909116179190911790556080909101516003909101555050505050505050565b6000828152600660209081526040808320815160a081018352815481526001820180548451818702810187019095528085529194929385840193909290830182828015610a3157602002820191906000526020600020905b815481526020019060010190808311610a1d575b5050509183525050600282015460ff811615156020830152610100900473ffffffffffffffffffffffffffffffffffffffff1660408201526003909101546060909101528051909150610ae0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265717565737420494420697320696e636f72726563740000000000000000006044820152606401610324565b60008381526006602090815260409091208351610b0592600190920191850190610d90565b506000838152600660205260409081902060020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790556004805491517ffd26ba4b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9092169163fd26ba4b91610b9a9187910190815260200190565b60206040518083038186803b158015610bb257600080fd5b505afa158015610bc6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bea9190610ff5565b6bffffffffffffffffffffffff166006600085815260200190815260200160002060030181905550505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610c98576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610324565b565b73ffffffffffffffffffffffffffffffffffffffff8116331415610d1a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610324565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215610dcb579160200282015b82811115610dcb578251825591602001919060010190610db0565b50610dd7929150610ddb565b5090565b5b80821115610dd75760008155600101610ddc565b803563ffffffff81168114610e0457600080fd5b919050565b600060208284031215610e1b57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff81168114610e3f57600080fd5b9392505050565b600060208284031215610e5857600080fd5b5035919050565b600060208284031215610e7157600080fd5b5051919050565b60008060408385031215610e8b57600080fd5b8235915060208084013567ffffffffffffffff80821115610eab57600080fd5b818601915086601f830112610ebf57600080fd5b813581811115610ed157610ed1611067565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f83011681018181108582111715610f1457610f14611067565b604052828152858101935084860182860187018b1015610f3357600080fd5b600095505b83861015610f56578035855260019590950194938601938601610f38565b508096505050505050509250929050565b60008060008060008060c08789031215610f8057600080fd5b863567ffffffffffffffff81168114610f9857600080fd5b9550610fa660208801610df0565b9450604087013561ffff81168114610fbd57600080fd5b9350610fcb60608801610df0565b92506080870135915060a08701358015158114610fe757600080fd5b809150509295509295509295565b60006020828403121561100757600080fd5b81516bffffffffffffffffffffffff81168114610e3f57600080fd5b6020808252825182820181905260009190848201906040850190845b8181101561105b5783518352928401929184019160010161103f565b50909695505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFV2PlusConsumerExampleABI = VRFV2PlusConsumerExampleMetaData.ABI

var VRFV2PlusConsumerExampleBin = VRFV2PlusConsumerExampleMetaData.Bin

func DeployVRFV2PlusConsumerExample(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, link common.Address, subOwner common.Address) (common.Address, *types.Transaction, *VRFV2PlusConsumerExample, error) {
	parsed, err := VRFV2PlusConsumerExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2PlusConsumerExampleBin), backend, vrfCoordinator, link, subOwner)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2PlusConsumerExample{VRFV2PlusConsumerExampleCaller: VRFV2PlusConsumerExampleCaller{contract: contract}, VRFV2PlusConsumerExampleTransactor: VRFV2PlusConsumerExampleTransactor{contract: contract}, VRFV2PlusConsumerExampleFilterer: VRFV2PlusConsumerExampleFilterer{contract: contract}}, nil
}

type VRFV2PlusConsumerExample struct {
	address common.Address
	abi     abi.ABI
	VRFV2PlusConsumerExampleCaller
	VRFV2PlusConsumerExampleTransactor
	VRFV2PlusConsumerExampleFilterer
}

type VRFV2PlusConsumerExampleCaller struct {
	contract *bind.BoundContract
}

type VRFV2PlusConsumerExampleTransactor struct {
	contract *bind.BoundContract
}

type VRFV2PlusConsumerExampleFilterer struct {
	contract *bind.BoundContract
}

type VRFV2PlusConsumerExampleSession struct {
	Contract     *VRFV2PlusConsumerExample
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2PlusConsumerExampleCallerSession struct {
	Contract *VRFV2PlusConsumerExampleCaller
	CallOpts bind.CallOpts
}

type VRFV2PlusConsumerExampleTransactorSession struct {
	Contract     *VRFV2PlusConsumerExampleTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2PlusConsumerExampleRaw struct {
	Contract *VRFV2PlusConsumerExample
}

type VRFV2PlusConsumerExampleCallerRaw struct {
	Contract *VRFV2PlusConsumerExampleCaller
}

type VRFV2PlusConsumerExampleTransactorRaw struct {
	Contract *VRFV2PlusConsumerExampleTransactor
}

func NewVRFV2PlusConsumerExample(address common.Address, backend bind.ContractBackend) (*VRFV2PlusConsumerExample, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2PlusConsumerExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2PlusConsumerExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusConsumerExample{address: address, abi: abi, VRFV2PlusConsumerExampleCaller: VRFV2PlusConsumerExampleCaller{contract: contract}, VRFV2PlusConsumerExampleTransactor: VRFV2PlusConsumerExampleTransactor{contract: contract}, VRFV2PlusConsumerExampleFilterer: VRFV2PlusConsumerExampleFilterer{contract: contract}}, nil
}

func NewVRFV2PlusConsumerExampleCaller(address common.Address, caller bind.ContractCaller) (*VRFV2PlusConsumerExampleCaller, error) {
	contract, err := bindVRFV2PlusConsumerExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusConsumerExampleCaller{contract: contract}, nil
}

func NewVRFV2PlusConsumerExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2PlusConsumerExampleTransactor, error) {
	contract, err := bindVRFV2PlusConsumerExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusConsumerExampleTransactor{contract: contract}, nil
}

func NewVRFV2PlusConsumerExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2PlusConsumerExampleFilterer, error) {
	contract, err := bindVRFV2PlusConsumerExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusConsumerExampleFilterer{contract: contract}, nil
}

func bindVRFV2PlusConsumerExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFV2PlusConsumerExampleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusConsumerExample.Contract.VRFV2PlusConsumerExampleCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.VRFV2PlusConsumerExampleTransactor.contract.Transfer(opts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.VRFV2PlusConsumerExampleTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusConsumerExample.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.contract.Transfer(opts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusConsumerExample.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) Owner() (common.Address, error) {
	return _VRFV2PlusConsumerExample.Contract.Owner(&_VRFV2PlusConsumerExample.CallOpts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCallerSession) Owner() (common.Address, error) {
	return _VRFV2PlusConsumerExample.Contract.Owner(&_VRFV2PlusConsumerExample.CallOpts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCaller) SLinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusConsumerExample.contract.Call(opts, &out, "s_linkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) SLinkToken() (common.Address, error) {
	return _VRFV2PlusConsumerExample.Contract.SLinkToken(&_VRFV2PlusConsumerExample.CallOpts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCallerSession) SLinkToken() (common.Address, error) {
	return _VRFV2PlusConsumerExample.Contract.SLinkToken(&_VRFV2PlusConsumerExample.CallOpts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCaller) SRequests(opts *bind.CallOpts, arg0 *big.Int) (SRequests,

	error) {
	var out []interface{}
	err := _VRFV2PlusConsumerExample.contract.Call(opts, &out, "s_requests", arg0)

	outstruct := new(SRequests)
	if err != nil {
		return *outstruct, err
	}

	outstruct.RequestId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Fulfilled = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.Requester = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.Cost = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) SRequests(arg0 *big.Int) (SRequests,

	error) {
	return _VRFV2PlusConsumerExample.Contract.SRequests(&_VRFV2PlusConsumerExample.CallOpts, arg0)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCallerSession) SRequests(arg0 *big.Int) (SRequests,

	error) {
	return _VRFV2PlusConsumerExample.Contract.SRequests(&_VRFV2PlusConsumerExample.CallOpts, arg0)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCaller) SVrfCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusConsumerExample.contract.Call(opts, &out, "s_vrfCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) SVrfCoordinator() (common.Address, error) {
	return _VRFV2PlusConsumerExample.Contract.SVrfCoordinator(&_VRFV2PlusConsumerExample.CallOpts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleCallerSession) SVrfCoordinator() (common.Address, error) {
	return _VRFV2PlusConsumerExample.Contract.SVrfCoordinator(&_VRFV2PlusConsumerExample.CallOpts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.contract.Transact(opts, "acceptOwnership")
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.AcceptOwnership(&_VRFV2PlusConsumerExample.TransactOpts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.AcceptOwnership(&_VRFV2PlusConsumerExample.TransactOpts)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.RawFulfillRandomWords(&_VRFV2PlusConsumerExample.TransactOpts, requestId, randomWords)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.RawFulfillRandomWords(&_VRFV2PlusConsumerExample.TransactOpts, requestId, randomWords)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactor) RedeemRandomness(opts *bind.TransactOpts, requestId *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.contract.Transact(opts, "redeemRandomness", requestId)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) RedeemRandomness(requestId *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.RedeemRandomness(&_VRFV2PlusConsumerExample.TransactOpts, requestId)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactorSession) RedeemRandomness(requestId *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.RedeemRandomness(&_VRFV2PlusConsumerExample.TransactOpts, requestId)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactor) RequestRandomWords(opts *bind.TransactOpts, subId uint64, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte, nativePayment bool) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.contract.Transact(opts, "requestRandomWords", subId, callbackGasLimit, requestConfirmations, numWords, keyHash, nativePayment)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) RequestRandomWords(subId uint64, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte, nativePayment bool) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.RequestRandomWords(&_VRFV2PlusConsumerExample.TransactOpts, subId, callbackGasLimit, requestConfirmations, numWords, keyHash, nativePayment)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactorSession) RequestRandomWords(subId uint64, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte, nativePayment bool) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.RequestRandomWords(&_VRFV2PlusConsumerExample.TransactOpts, subId, callbackGasLimit, requestConfirmations, numWords, keyHash, nativePayment)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactor) SetVRFCoordinator(opts *bind.TransactOpts, _vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.contract.Transact(opts, "setVRFCoordinator", _vrfCoordinator)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) SetVRFCoordinator(_vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.SetVRFCoordinator(&_VRFV2PlusConsumerExample.TransactOpts, _vrfCoordinator)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactorSession) SetVRFCoordinator(_vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.SetVRFCoordinator(&_VRFV2PlusConsumerExample.TransactOpts, _vrfCoordinator)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.TransferOwnership(&_VRFV2PlusConsumerExample.TransactOpts, to)
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusConsumerExample.Contract.TransferOwnership(&_VRFV2PlusConsumerExample.TransactOpts, to)
}

type VRFV2PlusConsumerExampleOwnershipTransferRequestedIterator struct {
	Event *VRFV2PlusConsumerExampleOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusConsumerExampleOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusConsumerExampleOwnershipTransferRequested)
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
		it.Event = new(VRFV2PlusConsumerExampleOwnershipTransferRequested)
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

func (it *VRFV2PlusConsumerExampleOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusConsumerExampleOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusConsumerExampleOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusConsumerExampleOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusConsumerExample.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusConsumerExampleOwnershipTransferRequestedIterator{contract: _VRFV2PlusConsumerExample.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2PlusConsumerExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusConsumerExample.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusConsumerExampleOwnershipTransferRequested)
				if err := _VRFV2PlusConsumerExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFV2PlusConsumerExampleOwnershipTransferRequested, error) {
	event := new(VRFV2PlusConsumerExampleOwnershipTransferRequested)
	if err := _VRFV2PlusConsumerExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2PlusConsumerExampleOwnershipTransferredIterator struct {
	Event *VRFV2PlusConsumerExampleOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusConsumerExampleOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusConsumerExampleOwnershipTransferred)
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
		it.Event = new(VRFV2PlusConsumerExampleOwnershipTransferred)
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

func (it *VRFV2PlusConsumerExampleOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusConsumerExampleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusConsumerExampleOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusConsumerExampleOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusConsumerExample.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusConsumerExampleOwnershipTransferredIterator{contract: _VRFV2PlusConsumerExample.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2PlusConsumerExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusConsumerExample.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusConsumerExampleOwnershipTransferred)
				if err := _VRFV2PlusConsumerExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExampleFilterer) ParseOwnershipTransferred(log types.Log) (*VRFV2PlusConsumerExampleOwnershipTransferred, error) {
	event := new(VRFV2PlusConsumerExampleOwnershipTransferred)
	if err := _VRFV2PlusConsumerExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SRequests struct {
	RequestId *big.Int
	Fulfilled bool
	Requester common.Address
	Cost      *big.Int
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExample) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFV2PlusConsumerExample.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFV2PlusConsumerExample.ParseOwnershipTransferRequested(log)
	case _VRFV2PlusConsumerExample.abi.Events["OwnershipTransferred"].ID:
		return _VRFV2PlusConsumerExample.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFV2PlusConsumerExampleOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFV2PlusConsumerExampleOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_VRFV2PlusConsumerExample *VRFV2PlusConsumerExample) Address() common.Address {
	return _VRFV2PlusConsumerExample.address
}

type VRFV2PlusConsumerExampleInterface interface {
	Owner(opts *bind.CallOpts) (common.Address, error)

	SLinkToken(opts *bind.CallOpts) (common.Address, error)

	SRequests(opts *bind.CallOpts, arg0 *big.Int) (SRequests,

		error)

	SVrfCoordinator(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RedeemRandomness(opts *bind.TransactOpts, requestId *big.Int) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts, subId uint64, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte, nativePayment bool) (*types.Transaction, error)

	SetVRFCoordinator(opts *bind.TransactOpts, _vrfCoordinator common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusConsumerExampleOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2PlusConsumerExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFV2PlusConsumerExampleOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusConsumerExampleOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2PlusConsumerExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFV2PlusConsumerExampleOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
