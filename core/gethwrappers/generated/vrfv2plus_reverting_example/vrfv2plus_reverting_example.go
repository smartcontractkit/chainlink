// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrfv2plus_reverting_example

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

var VRFV2PlusRevertingExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"}],\"name\":\"OnlyOwnerOrCoordinator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"createSubscriptionAndFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSubId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"minReqConfs\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"}],\"name\":\"requestRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_gasAvailable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"_subId\",\"type\":\"uint64\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"topUpSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"name\":\"updateSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50604051620013e6380380620013e68339810160408190526200003491620001cc565b8133806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf8162000103565b5050600280546001600160a01b03199081166001600160a01b0394851617909155600580548216958416959095179094555060068054909316911617905562000204565b6001600160a01b0381163314156200015e5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001c757600080fd5b919050565b60008060408385031215620001e057600080fd5b620001eb83620001af565b9150620001fb60208401620001af565b90509250929050565b6111d280620002146000396000f3fe608060405234801561001057600080fd5b50600436106100df5760003560e01c80638da5cb5b1161008c578063f08c5daa11610066578063f08c5daa146101a4578063f2fde38b146101ad578063f6eaffc8146101c0578063fbc6ba6f146101d357600080fd5b80638da5cb5b14610160578063cf62c8ab14610188578063e89e106a1461019b57600080fd5b80632fa4e442116100bd5780632fa4e4421461013257806336bfffed1461014557806379ba50971461015857600080fd5b8063177b9692146100e45780631fe543e31461010a5780632d6d99f31461011f575b600080fd5b6100f76100f2366004610dce565b610211565b6040519081526020015b60405180910390f35b61011d610118366004610e69565b61031a565b005b61011d61012d366004610cd1565b6103a0565b61011d610140366004610f2a565b6104d0565b61011d610153366004610d08565b610634565b61011d6107bc565b60005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610101565b61011d610196366004610f2a565b6108b9565b6100f760045481565b6100f760075481565b61011d6101bb366004610caf565b610ac4565b6100f76101ce366004610e37565b610ad8565b60025474010000000000000000000000000000000000000000900467ffffffffffffffff1660405167ffffffffffffffff9091168152602001610101565b6040805160c08101825286815267ffffffffffffffff861660208083019190915261ffff86168284015263ffffffff80861660608401528416608083015282519081018352600080825260a083019190915260055492517f596b8b88000000000000000000000000000000000000000000000000000000008152909273ffffffffffffffffffffffffffffffffffffffff169063596b8b88906102b890849060040161100f565b602060405180830381600087803b1580156102d257600080fd5b505af11580156102e6573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061030a9190610e50565b6004819055979650505050505050565b60025473ffffffffffffffffffffffffffffffffffffffff163314610392576002546040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff90911660248201526044015b60405180910390fd5b61039c8282600080fd5b5050565b60005473ffffffffffffffffffffffffffffffffffffffff1633148015906103e0575060025473ffffffffffffffffffffffffffffffffffffffff163314155b15610464573361040560005473ffffffffffffffffffffffffffffffffffffffff1690565b6002546040517f061db9c100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff93841660048201529183166024830152919091166044820152606401610389565b6002805467ffffffffffffffff90921674010000000000000000000000000000000000000000027fffffffff0000000000000000000000000000000000000000000000000000000090921673ffffffffffffffffffffffffffffffffffffffff90931692909217179055565b60025474010000000000000000000000000000000000000000900467ffffffffffffffff1661055b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f737562206e6f74207365740000000000000000000000000000000000000000006044820152606401610389565b600654600554600254604080517401000000000000000000000000000000000000000090920467ffffffffffffffff16602083015273ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918591015b6040516020818303038152906040526040518463ffffffff1660e01b81526004016105e293929190610fc3565b602060405180830381600087803b1580156105fc57600080fd5b505af1158015610610573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061039c9190610dac565b60025474010000000000000000000000000000000000000000900467ffffffffffffffff166106bf576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f7375624944206e6f7420736574000000000000000000000000000000000000006044820152606401610389565b60005b815181101561039c57600554600254835173ffffffffffffffffffffffffffffffffffffffff90921691637341c10c9174010000000000000000000000000000000000000000900467ffffffffffffffff169085908590811061072757610727611151565b60200260200101516040518363ffffffff1660e01b815260040161077792919067ffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b15801561079157600080fd5b505af11580156107a5573d6000803e3d6000fd5b5050505080806107b4906110f1565b9150506106c2565b60015473ffffffffffffffffffffffffffffffffffffffff16331461083d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610389565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60025474010000000000000000000000000000000000000000900467ffffffffffffffff1661055b57600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a21a23e46040518163ffffffff1660e01b8152600401602060405180830381600087803b15801561094c57600080fd5b505af1158015610960573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109849190610f0d565b600280547fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000067ffffffffffffffff938416810291909117918290556005546040517f7341c10c00000000000000000000000000000000000000000000000000000000815291909204909216600483015230602483015273ffffffffffffffffffffffffffffffffffffffff1690637341c10c90604401600060405180830381600087803b158015610a4b57600080fd5b505af1158015610a5f573d6000803e3d6000fd5b5050600654600554600254604080517401000000000000000000000000000000000000000090920467ffffffffffffffff16602083015273ffffffffffffffffffffffffffffffffffffffff9384169550634000aea0945092909116918591016105b5565b610acc610af9565b610ad581610b7c565b50565b60038181548110610ae857600080fd5b600091825260209091200154905081565b60005473ffffffffffffffffffffffffffffffffffffffff163314610b7a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610389565b565b73ffffffffffffffffffffffffffffffffffffffff8116331415610bfc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610389565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b803573ffffffffffffffffffffffffffffffffffffffff81168114610c9657600080fd5b919050565b803563ffffffff81168114610c9657600080fd5b600060208284031215610cc157600080fd5b610cca82610c72565b9392505050565b60008060408385031215610ce457600080fd5b610ced83610c72565b91506020830135610cfd816111af565b809150509250929050565b60006020808385031215610d1b57600080fd5b823567ffffffffffffffff811115610d3257600080fd5b8301601f81018513610d4357600080fd5b8035610d56610d51826110cd565b61107e565b80828252848201915084840188868560051b8701011115610d7657600080fd5b600094505b83851015610da057610d8c81610c72565b835260019490940193918501918501610d7b565b50979650505050505050565b600060208284031215610dbe57600080fd5b81518015158114610cca57600080fd5b600080600080600060a08688031215610de657600080fd5b853594506020860135610df8816111af565b9350604086013561ffff81168114610e0f57600080fd5b9250610e1d60608701610c9b565b9150610e2b60808701610c9b565b90509295509295909350565b600060208284031215610e4957600080fd5b5035919050565b600060208284031215610e6257600080fd5b5051919050565b60008060408385031215610e7c57600080fd5b8235915060208084013567ffffffffffffffff811115610e9b57600080fd5b8401601f81018613610eac57600080fd5b8035610eba610d51826110cd565b80828252848201915084840189868560051b8701011115610eda57600080fd5b600094505b83851015610efd578035835260019490940193918501918501610edf565b5080955050505050509250929050565b600060208284031215610f1f57600080fd5b8151610cca816111af565b600060208284031215610f3c57600080fd5b81356bffffffffffffffffffffffff81168114610cca57600080fd5b6000815180845260005b81811015610f7e57602081850181015186830182015201610f62565b81811115610f90576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b73ffffffffffffffffffffffffffffffffffffffff841681526bffffffffffffffffffffffff831660208201526060604082015260006110066060830184610f58565b95945050505050565b602081528151602082015267ffffffffffffffff602083015116604082015261ffff60408301511660608201526000606083015163ffffffff80821660808501528060808601511660a0850152505060a083015160c08084015261107660e0840182610f58565b949350505050565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156110c5576110c5611180565b604052919050565b600067ffffffffffffffff8211156110e7576110e7611180565b5060051b60200190565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82141561114a577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b67ffffffffffffffff81168114610ad557600080fdfea164736f6c6343000806000a",
}

var VRFV2PlusRevertingExampleABI = VRFV2PlusRevertingExampleMetaData.ABI

var VRFV2PlusRevertingExampleBin = VRFV2PlusRevertingExampleMetaData.Bin

func DeployVRFV2PlusRevertingExample(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, link common.Address) (common.Address, *types.Transaction, *VRFV2PlusRevertingExample, error) {
	parsed, err := VRFV2PlusRevertingExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2PlusRevertingExampleBin), backend, vrfCoordinator, link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2PlusRevertingExample{VRFV2PlusRevertingExampleCaller: VRFV2PlusRevertingExampleCaller{contract: contract}, VRFV2PlusRevertingExampleTransactor: VRFV2PlusRevertingExampleTransactor{contract: contract}, VRFV2PlusRevertingExampleFilterer: VRFV2PlusRevertingExampleFilterer{contract: contract}}, nil
}

type VRFV2PlusRevertingExample struct {
	address common.Address
	abi     abi.ABI
	VRFV2PlusRevertingExampleCaller
	VRFV2PlusRevertingExampleTransactor
	VRFV2PlusRevertingExampleFilterer
}

type VRFV2PlusRevertingExampleCaller struct {
	contract *bind.BoundContract
}

type VRFV2PlusRevertingExampleTransactor struct {
	contract *bind.BoundContract
}

type VRFV2PlusRevertingExampleFilterer struct {
	contract *bind.BoundContract
}

type VRFV2PlusRevertingExampleSession struct {
	Contract     *VRFV2PlusRevertingExample
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2PlusRevertingExampleCallerSession struct {
	Contract *VRFV2PlusRevertingExampleCaller
	CallOpts bind.CallOpts
}

type VRFV2PlusRevertingExampleTransactorSession struct {
	Contract     *VRFV2PlusRevertingExampleTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2PlusRevertingExampleRaw struct {
	Contract *VRFV2PlusRevertingExample
}

type VRFV2PlusRevertingExampleCallerRaw struct {
	Contract *VRFV2PlusRevertingExampleCaller
}

type VRFV2PlusRevertingExampleTransactorRaw struct {
	Contract *VRFV2PlusRevertingExampleTransactor
}

func NewVRFV2PlusRevertingExample(address common.Address, backend bind.ContractBackend) (*VRFV2PlusRevertingExample, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2PlusRevertingExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2PlusRevertingExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusRevertingExample{address: address, abi: abi, VRFV2PlusRevertingExampleCaller: VRFV2PlusRevertingExampleCaller{contract: contract}, VRFV2PlusRevertingExampleTransactor: VRFV2PlusRevertingExampleTransactor{contract: contract}, VRFV2PlusRevertingExampleFilterer: VRFV2PlusRevertingExampleFilterer{contract: contract}}, nil
}

func NewVRFV2PlusRevertingExampleCaller(address common.Address, caller bind.ContractCaller) (*VRFV2PlusRevertingExampleCaller, error) {
	contract, err := bindVRFV2PlusRevertingExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusRevertingExampleCaller{contract: contract}, nil
}

func NewVRFV2PlusRevertingExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2PlusRevertingExampleTransactor, error) {
	contract, err := bindVRFV2PlusRevertingExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusRevertingExampleTransactor{contract: contract}, nil
}

func NewVRFV2PlusRevertingExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2PlusRevertingExampleFilterer, error) {
	contract, err := bindVRFV2PlusRevertingExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusRevertingExampleFilterer{contract: contract}, nil
}

func bindVRFV2PlusRevertingExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFV2PlusRevertingExampleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusRevertingExample.Contract.VRFV2PlusRevertingExampleCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.VRFV2PlusRevertingExampleTransactor.contract.Transfer(opts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.VRFV2PlusRevertingExampleTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusRevertingExample.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.contract.Transfer(opts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCaller) GetSubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFV2PlusRevertingExample.contract.Call(opts, &out, "getSubId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) GetSubId() (uint64, error) {
	return _VRFV2PlusRevertingExample.Contract.GetSubId(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCallerSession) GetSubId() (uint64, error) {
	return _VRFV2PlusRevertingExample.Contract.GetSubId(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusRevertingExample.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) Owner() (common.Address, error) {
	return _VRFV2PlusRevertingExample.Contract.Owner(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCallerSession) Owner() (common.Address, error) {
	return _VRFV2PlusRevertingExample.Contract.Owner(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCaller) SGasAvailable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusRevertingExample.contract.Call(opts, &out, "s_gasAvailable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) SGasAvailable() (*big.Int, error) {
	return _VRFV2PlusRevertingExample.Contract.SGasAvailable(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCallerSession) SGasAvailable() (*big.Int, error) {
	return _VRFV2PlusRevertingExample.Contract.SGasAvailable(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCaller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusRevertingExample.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFV2PlusRevertingExample.Contract.SRandomWords(&_VRFV2PlusRevertingExample.CallOpts, arg0)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFV2PlusRevertingExample.Contract.SRandomWords(&_VRFV2PlusRevertingExample.CallOpts, arg0)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCaller) SRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusRevertingExample.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) SRequestId() (*big.Int, error) {
	return _VRFV2PlusRevertingExample.Contract.SRequestId(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleCallerSession) SRequestId() (*big.Int, error) {
	return _VRFV2PlusRevertingExample.Contract.SRequestId(&_VRFV2PlusRevertingExample.CallOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "acceptOwnership")
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.AcceptOwnership(&_VRFV2PlusRevertingExample.TransactOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.AcceptOwnership(&_VRFV2PlusRevertingExample.TransactOpts)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) CreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "createSubscriptionAndFund", amount)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) CreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.CreateSubscriptionAndFund(&_VRFV2PlusRevertingExample.TransactOpts, amount)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) CreateSubscriptionAndFund(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.CreateSubscriptionAndFund(&_VRFV2PlusRevertingExample.TransactOpts, amount)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.RawFulfillRandomWords(&_VRFV2PlusRevertingExample.TransactOpts, requestId, randomWords)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.RawFulfillRandomWords(&_VRFV2PlusRevertingExample.TransactOpts, requestId, randomWords)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) RequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "requestRandomness", keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) RequestRandomness(keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.RequestRandomness(&_VRFV2PlusRevertingExample.TransactOpts, keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) RequestRandomness(keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.RequestRandomness(&_VRFV2PlusRevertingExample.TransactOpts, keyHash, subId, minReqConfs, callbackGasLimit, numWords)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) SetConfig(opts *bind.TransactOpts, _vrfCoordinator common.Address, _subId uint64) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "setConfig", _vrfCoordinator, _subId)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) SetConfig(_vrfCoordinator common.Address, _subId uint64) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.SetConfig(&_VRFV2PlusRevertingExample.TransactOpts, _vrfCoordinator, _subId)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) SetConfig(_vrfCoordinator common.Address, _subId uint64) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.SetConfig(&_VRFV2PlusRevertingExample.TransactOpts, _vrfCoordinator, _subId)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "topUpSubscription", amount)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.TopUpSubscription(&_VRFV2PlusRevertingExample.TransactOpts, amount)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.TopUpSubscription(&_VRFV2PlusRevertingExample.TransactOpts, amount)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.TransferOwnership(&_VRFV2PlusRevertingExample.TransactOpts, to)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.TransferOwnership(&_VRFV2PlusRevertingExample.TransactOpts, to)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactor) UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.contract.Transact(opts, "updateSubscription", consumers)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.UpdateSubscription(&_VRFV2PlusRevertingExample.TransactOpts, consumers)
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleTransactorSession) UpdateSubscription(consumers []common.Address) (*types.Transaction, error) {
	return _VRFV2PlusRevertingExample.Contract.UpdateSubscription(&_VRFV2PlusRevertingExample.TransactOpts, consumers)
}

type VRFV2PlusRevertingExampleOwnershipTransferRequestedIterator struct {
	Event *VRFV2PlusRevertingExampleOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusRevertingExampleOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusRevertingExampleOwnershipTransferRequested)
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
		it.Event = new(VRFV2PlusRevertingExampleOwnershipTransferRequested)
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

func (it *VRFV2PlusRevertingExampleOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusRevertingExampleOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusRevertingExampleOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusRevertingExampleOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusRevertingExample.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusRevertingExampleOwnershipTransferRequestedIterator{contract: _VRFV2PlusRevertingExample.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2PlusRevertingExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusRevertingExample.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusRevertingExampleOwnershipTransferRequested)
				if err := _VRFV2PlusRevertingExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFV2PlusRevertingExampleOwnershipTransferRequested, error) {
	event := new(VRFV2PlusRevertingExampleOwnershipTransferRequested)
	if err := _VRFV2PlusRevertingExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2PlusRevertingExampleOwnershipTransferredIterator struct {
	Event *VRFV2PlusRevertingExampleOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusRevertingExampleOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusRevertingExampleOwnershipTransferred)
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
		it.Event = new(VRFV2PlusRevertingExampleOwnershipTransferred)
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

func (it *VRFV2PlusRevertingExampleOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusRevertingExampleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusRevertingExampleOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusRevertingExampleOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusRevertingExample.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusRevertingExampleOwnershipTransferredIterator{contract: _VRFV2PlusRevertingExample.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2PlusRevertingExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusRevertingExample.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusRevertingExampleOwnershipTransferred)
				if err := _VRFV2PlusRevertingExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExampleFilterer) ParseOwnershipTransferred(log types.Log) (*VRFV2PlusRevertingExampleOwnershipTransferred, error) {
	event := new(VRFV2PlusRevertingExampleOwnershipTransferred)
	if err := _VRFV2PlusRevertingExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExample) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFV2PlusRevertingExample.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFV2PlusRevertingExample.ParseOwnershipTransferRequested(log)
	case _VRFV2PlusRevertingExample.abi.Events["OwnershipTransferred"].ID:
		return _VRFV2PlusRevertingExample.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFV2PlusRevertingExampleOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFV2PlusRevertingExampleOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_VRFV2PlusRevertingExample *VRFV2PlusRevertingExample) Address() common.Address {
	return _VRFV2PlusRevertingExample.address
}

type VRFV2PlusRevertingExampleInterface interface {
	GetSubId(opts *bind.CallOpts) (uint64, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SGasAvailable(opts *bind.CallOpts) (*big.Int, error)

	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestId(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	CreateSubscriptionAndFund(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomness(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, minReqConfs uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, _vrfCoordinator common.Address, _subId uint64) (*types.Transaction, error)

	TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateSubscription(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusRevertingExampleOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2PlusRevertingExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFV2PlusRevertingExampleOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusRevertingExampleOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2PlusRevertingExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFV2PlusRevertingExampleOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
