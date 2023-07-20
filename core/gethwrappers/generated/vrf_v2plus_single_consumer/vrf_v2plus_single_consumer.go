// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_v2plus_single_consumer

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

var VRFV2PlusSingleConsumerExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"nativePayment\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlySubOwnerCanSetVRFCoordinator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"fundAndRequestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestConfig\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"nativePayment\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"}],\"name\":\"setVRFCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"subscribe\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"topUpSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"unsubscribe\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50604051620015cb380380620015cb833981016040819052620000349162000360565b866001600160a01b0381166200007f5760405162461bcd60e51b815260206004820152600c60248201526b7a65726f206164647265737360a01b604482015260640160405180910390fd5b600080546001600160a01b03199081166001600160a01b039384161782556002805482168b8516179055600380548216938a169390931790925560098054339316929092179091556040805160c08101825291825263ffffffff8781166020840181905261ffff8816928401839052908616606084018190526080840186905284151560a0909401849052600480546001600160701b0319166801000000000000000090930261ffff60601b1916929092176c010000000000000000000000009093029290921763ffffffff60701b1916600160701b90920291909117905560058390556006805460ff191690911790556200017a62000187565b5050505050505062000444565b6009546001600160a01b031633146200019f57600080fd5b604080516001808252818301909252600091602080830190803683370190505090503081600081518110620001d857620001d86200042e565b6001600160a01b039283166020918202929092018101919091526002546040805163288688f960e21b81529051919093169263a21a23e49260048083019391928290030181600087803b1580156200022f57600080fd5b505af115801562000244573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200026a9190620003fc565b600480546001600160401b0319166001600160401b0392909216918217905560025482516001600160a01b0390911691637341c10c918490600090620002b457620002b46200042e565b60200260200101516040518363ffffffff1660e01b8152600401620002f79291906001600160401b039290921682526001600160a01b0316602082015260400190565b600060405180830381600087803b1580156200031257600080fd5b505af115801562000327573d6000803e3d6000fd5b5050505050565b80516001600160a01b03811681146200034657600080fd5b919050565b805163ffffffff811681146200034657600080fd5b600080600080600080600060e0888a0312156200037c57600080fd5b62000387886200032e565b965062000397602089016200032e565b9550620003a7604089016200034b565b9450606088015161ffff81168114620003bf57600080fd5b9350620003cf608089016200034b565b925060a0880151915060c08801518015158114620003ec57600080fd5b8091505092959891949750929550565b6000602082840312156200040f57600080fd5b81516001600160401b03811681146200042757600080fd5b9392505050565b634e487b7160e01b600052603260045260246000fd5b61117780620004546000396000f3fe608060405234801561001057600080fd5b50600436106100c85760003560e01c80637db9263f11610081578063e0c862891161005b578063e0c86289146101f7578063e89e106a146101ff578063f6eaffc81461021657600080fd5b80637db9263f1461012e57806386850e93146101dc5780638f449a05146101ef57600080fd5b806344ff81ce116100b257806344ff81ce146100f55780636fd700bb146101085780637262561c1461011b57600080fd5b8062f714ce146100cd5780631fe543e3146100e2575b600080fd5b6100e06100db366004610eaf565b610229565b005b6100e06100f0366004610edb565b6102fe565b6100e0610103366004610e39565b610384565b6100e0610116366004610e7d565b61043e565b6100e0610129366004610e39565b6106d3565b60045460055460065461018d9267ffffffffffffffff81169263ffffffff68010000000000000000830481169361ffff6c01000000000000000000000000850416936e0100000000000000000000000000009004909116919060ff1686565b6040805167ffffffffffffffff909716875263ffffffff958616602088015261ffff90941693860193909352921660608401526080830191909152151560a082015260c0015b60405180910390f35b6100e06101ea366004610e7d565b6107be565b6100e06108a3565b6100e0610aa3565b61020860085481565b6040519081526020016101d3565b610208610224366004610e7d565b610c55565b60095473ffffffffffffffffffffffffffffffffffffffff16331461024d57600080fd5b6003546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8381166004830152602482018590529091169063a9059cbb90604401602060405180830381600087803b1580156102c157600080fd5b505af11580156102d5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102f99190610e5b565b505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610376576000546040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff90911660248201526044015b60405180910390fd5b6103808282610c76565b5050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146103f7576001546040517f4ae338ff00000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff909116602482015260440161036d565b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60095473ffffffffffffffffffffffffffffffffffffffff16331461046257600080fd5b6040805160c08101825260045467ffffffffffffffff811680835263ffffffff680100000000000000008304811660208086019190915261ffff6c01000000000000000000000000850416858701526e010000000000000000000000000000909304166060840152600554608084015260065460ff16151560a084015260035460025485518085019390935285518084039094018452828601958690527f4000aea000000000000000000000000000000000000000000000000000000000909552929373ffffffffffffffffffffffffffffffffffffffff93841693634000aea093610557939190921691879160440161105f565b602060405180830381600087803b15801561057157600080fd5b505af1158015610585573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105a99190610e5b565b5060006040518060c0016040528083608001518152602001836000015167ffffffffffffffff168152602001836040015161ffff168152602001836020015163ffffffff168152602001836060015163ffffffff16815260200161062060405180602001604052808660a001511515815250610cf4565b90526002546040517f596b8b8800000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff169063596b8b889061067990849060040161109d565b602060405180830381600087803b15801561069357600080fd5b505af11580156106a7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106cb9190610e96565b600855505050565b60095473ffffffffffffffffffffffffffffffffffffffff1633146106f757600080fd5b600254600480546040517fd7ae1d3000000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9091169181019190915273ffffffffffffffffffffffffffffffffffffffff83811660248301529091169063d7ae1d3090604401600060405180830381600087803b15801561077b57600080fd5b505af115801561078f573d6000803e3d6000fd5b5050600480547fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000169055505050565b60095473ffffffffffffffffffffffffffffffffffffffff1633146107e257600080fd5b6003546002546004546040805167ffffffffffffffff909216602083015273ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918591016040516020818303038152906040526040518463ffffffff1660e01b81526004016108519392919061105f565b602060405180830381600087803b15801561086b57600080fd5b505af115801561087f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103809190610e5b565b60095473ffffffffffffffffffffffffffffffffffffffff1633146108c757600080fd5b6040805160018082528183019092526000916020808301908036833701905050905030816000815181106108fd576108fd61110c565b73ffffffffffffffffffffffffffffffffffffffff928316602091820292909201810191909152600254604080517fa21a23e40000000000000000000000000000000000000000000000000000000081529051919093169263a21a23e49260048083019391928290030181600087803b15801561097957600080fd5b505af115801561098d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109b19190610fca565b600480547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff929092169182179055600254825173ffffffffffffffffffffffffffffffffffffffff90911691637341c10c918490600090610a1e57610a1e61110c565b60200260200101516040518363ffffffff1660e01b8152600401610a6e92919067ffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b158015610a8857600080fd5b505af1158015610a9c573d6000803e3d6000fd5b5050505050565b60095473ffffffffffffffffffffffffffffffffffffffff163314610ac757600080fd5b6040805160c0808201835260045467ffffffffffffffff808216845263ffffffff6801000000000000000083048116602080870191825261ffff6c0100000000000000000000000086048116888a019081526e01000000000000000000000000000090960484166060808a019182526005546080808c0191825260065460ff16151560a0808e019182528e519c8d018f5292518c528c519099168b8701529851909316898c01529351851693880193909352915190921693850193909352855190810190955251151584529192600092820190610ba390610cf4565b90526002546040517f596b8b8800000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff169063596b8b8890610bfc90849060040161109d565b602060405180830381600087803b158015610c1657600080fd5b505af1158015610c2a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c4e9190610e96565b6008555050565b60078181548110610c6557600080fd5b600091825260209091200154905081565b6008548214610ce1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265717565737420494420697320696e636f7272656374000000000000000000604482015260640161036d565b80516102f9906007906020840190610db0565b60607f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa82604051602401610d2d91511515815260200190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915292915050565b828054828255906000526020600020908101928215610deb579160200282015b82811115610deb578251825591602001919060010190610dd0565b50610df7929150610dfb565b5090565b5b80821115610df75760008155600101610dfc565b803573ffffffffffffffffffffffffffffffffffffffff81168114610e3457600080fd5b919050565b600060208284031215610e4b57600080fd5b610e5482610e10565b9392505050565b600060208284031215610e6d57600080fd5b81518015158114610e5457600080fd5b600060208284031215610e8f57600080fd5b5035919050565b600060208284031215610ea857600080fd5b5051919050565b60008060408385031215610ec257600080fd5b82359150610ed260208401610e10565b90509250929050565b60008060408385031215610eee57600080fd5b8235915060208084013567ffffffffffffffff80821115610f0e57600080fd5b818601915086601f830112610f2257600080fd5b813581811115610f3457610f3461113b565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f83011681018181108582111715610f7757610f7761113b565b604052828152858101935084860182860187018b1015610f9657600080fd5b600095505b83861015610fb9578035855260019590950194938601938601610f9b565b508096505050505050509250929050565b600060208284031215610fdc57600080fd5b815167ffffffffffffffff81168114610e5457600080fd5b6000815180845260005b8181101561101a57602081850181015186830182015201610ffe565b8181111561102c576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b73ffffffffffffffffffffffffffffffffffffffff841681528260208201526060604082015260006110946060830184610ff4565b95945050505050565b602081528151602082015267ffffffffffffffff602083015116604082015261ffff60408301511660608201526000606083015163ffffffff80821660808501528060808601511660a0850152505060a083015160c08084015261110460e0840182610ff4565b949350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFV2PlusSingleConsumerExampleABI = VRFV2PlusSingleConsumerExampleMetaData.ABI

var VRFV2PlusSingleConsumerExampleBin = VRFV2PlusSingleConsumerExampleMetaData.Bin

func DeployVRFV2PlusSingleConsumerExample(auth *bind.TransactOpts, backend bind.ContractBackend, vrfCoordinator common.Address, link common.Address, callbackGasLimit uint32, requestConfirmations uint16, numWords uint32, keyHash [32]byte, nativePayment bool) (common.Address, *types.Transaction, *VRFV2PlusSingleConsumerExample, error) {
	parsed, err := VRFV2PlusSingleConsumerExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2PlusSingleConsumerExampleBin), backend, vrfCoordinator, link, callbackGasLimit, requestConfirmations, numWords, keyHash, nativePayment)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2PlusSingleConsumerExample{VRFV2PlusSingleConsumerExampleCaller: VRFV2PlusSingleConsumerExampleCaller{contract: contract}, VRFV2PlusSingleConsumerExampleTransactor: VRFV2PlusSingleConsumerExampleTransactor{contract: contract}, VRFV2PlusSingleConsumerExampleFilterer: VRFV2PlusSingleConsumerExampleFilterer{contract: contract}}, nil
}

type VRFV2PlusSingleConsumerExample struct {
	address common.Address
	abi     abi.ABI
	VRFV2PlusSingleConsumerExampleCaller
	VRFV2PlusSingleConsumerExampleTransactor
	VRFV2PlusSingleConsumerExampleFilterer
}

type VRFV2PlusSingleConsumerExampleCaller struct {
	contract *bind.BoundContract
}

type VRFV2PlusSingleConsumerExampleTransactor struct {
	contract *bind.BoundContract
}

type VRFV2PlusSingleConsumerExampleFilterer struct {
	contract *bind.BoundContract
}

type VRFV2PlusSingleConsumerExampleSession struct {
	Contract     *VRFV2PlusSingleConsumerExample
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2PlusSingleConsumerExampleCallerSession struct {
	Contract *VRFV2PlusSingleConsumerExampleCaller
	CallOpts bind.CallOpts
}

type VRFV2PlusSingleConsumerExampleTransactorSession struct {
	Contract     *VRFV2PlusSingleConsumerExampleTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2PlusSingleConsumerExampleRaw struct {
	Contract *VRFV2PlusSingleConsumerExample
}

type VRFV2PlusSingleConsumerExampleCallerRaw struct {
	Contract *VRFV2PlusSingleConsumerExampleCaller
}

type VRFV2PlusSingleConsumerExampleTransactorRaw struct {
	Contract *VRFV2PlusSingleConsumerExampleTransactor
}

func NewVRFV2PlusSingleConsumerExample(address common.Address, backend bind.ContractBackend) (*VRFV2PlusSingleConsumerExample, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2PlusSingleConsumerExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2PlusSingleConsumerExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusSingleConsumerExample{address: address, abi: abi, VRFV2PlusSingleConsumerExampleCaller: VRFV2PlusSingleConsumerExampleCaller{contract: contract}, VRFV2PlusSingleConsumerExampleTransactor: VRFV2PlusSingleConsumerExampleTransactor{contract: contract}, VRFV2PlusSingleConsumerExampleFilterer: VRFV2PlusSingleConsumerExampleFilterer{contract: contract}}, nil
}

func NewVRFV2PlusSingleConsumerExampleCaller(address common.Address, caller bind.ContractCaller) (*VRFV2PlusSingleConsumerExampleCaller, error) {
	contract, err := bindVRFV2PlusSingleConsumerExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusSingleConsumerExampleCaller{contract: contract}, nil
}

func NewVRFV2PlusSingleConsumerExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2PlusSingleConsumerExampleTransactor, error) {
	contract, err := bindVRFV2PlusSingleConsumerExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusSingleConsumerExampleTransactor{contract: contract}, nil
}

func NewVRFV2PlusSingleConsumerExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2PlusSingleConsumerExampleFilterer, error) {
	contract, err := bindVRFV2PlusSingleConsumerExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusSingleConsumerExampleFilterer{contract: contract}, nil
}

func bindVRFV2PlusSingleConsumerExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFV2PlusSingleConsumerExampleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusSingleConsumerExample.Contract.VRFV2PlusSingleConsumerExampleCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.VRFV2PlusSingleConsumerExampleTransactor.contract.Transfer(opts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.VRFV2PlusSingleConsumerExampleTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusSingleConsumerExample.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.contract.Transfer(opts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleCaller) SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusSingleConsumerExample.contract.Call(opts, &out, "s_randomWords", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.SRandomWords(&_VRFV2PlusSingleConsumerExample.CallOpts, arg0)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleCallerSession) SRandomWords(arg0 *big.Int) (*big.Int, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.SRandomWords(&_VRFV2PlusSingleConsumerExample.CallOpts, arg0)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleCaller) SRequestConfig(opts *bind.CallOpts) (SRequestConfig,

	error) {
	var out []interface{}
	err := _VRFV2PlusSingleConsumerExample.contract.Call(opts, &out, "s_requestConfig")

	outstruct := new(SRequestConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.SubId = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.CallbackGasLimit = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.RequestConfirmations = *abi.ConvertType(out[2], new(uint16)).(*uint16)
	outstruct.NumWords = *abi.ConvertType(out[3], new(uint32)).(*uint32)
	outstruct.KeyHash = *abi.ConvertType(out[4], new([32]byte)).(*[32]byte)
	outstruct.NativePayment = *abi.ConvertType(out[5], new(bool)).(*bool)

	return *outstruct, err

}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) SRequestConfig() (SRequestConfig,

	error) {
	return _VRFV2PlusSingleConsumerExample.Contract.SRequestConfig(&_VRFV2PlusSingleConsumerExample.CallOpts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleCallerSession) SRequestConfig() (SRequestConfig,

	error) {
	return _VRFV2PlusSingleConsumerExample.Contract.SRequestConfig(&_VRFV2PlusSingleConsumerExample.CallOpts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleCaller) SRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusSingleConsumerExample.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) SRequestId() (*big.Int, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.SRequestId(&_VRFV2PlusSingleConsumerExample.CallOpts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleCallerSession) SRequestId() (*big.Int, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.SRequestId(&_VRFV2PlusSingleConsumerExample.CallOpts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactor) FundAndRequestRandomWords(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.contract.Transact(opts, "fundAndRequestRandomWords", amount)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) FundAndRequestRandomWords(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.FundAndRequestRandomWords(&_VRFV2PlusSingleConsumerExample.TransactOpts, amount)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactorSession) FundAndRequestRandomWords(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.FundAndRequestRandomWords(&_VRFV2PlusSingleConsumerExample.TransactOpts, amount)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.RawFulfillRandomWords(&_VRFV2PlusSingleConsumerExample.TransactOpts, requestId, randomWords)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.RawFulfillRandomWords(&_VRFV2PlusSingleConsumerExample.TransactOpts, requestId, randomWords)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactor) RequestRandomWords(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.contract.Transact(opts, "requestRandomWords")
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) RequestRandomWords() (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.RequestRandomWords(&_VRFV2PlusSingleConsumerExample.TransactOpts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactorSession) RequestRandomWords() (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.RequestRandomWords(&_VRFV2PlusSingleConsumerExample.TransactOpts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactor) SetVRFCoordinator(opts *bind.TransactOpts, _vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.contract.Transact(opts, "setVRFCoordinator", _vrfCoordinator)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) SetVRFCoordinator(_vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.SetVRFCoordinator(&_VRFV2PlusSingleConsumerExample.TransactOpts, _vrfCoordinator)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactorSession) SetVRFCoordinator(_vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.SetVRFCoordinator(&_VRFV2PlusSingleConsumerExample.TransactOpts, _vrfCoordinator)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactor) Subscribe(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.contract.Transact(opts, "subscribe")
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) Subscribe() (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.Subscribe(&_VRFV2PlusSingleConsumerExample.TransactOpts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactorSession) Subscribe() (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.Subscribe(&_VRFV2PlusSingleConsumerExample.TransactOpts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactor) TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.contract.Transact(opts, "topUpSubscription", amount)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.TopUpSubscription(&_VRFV2PlusSingleConsumerExample.TransactOpts, amount)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactorSession) TopUpSubscription(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.TopUpSubscription(&_VRFV2PlusSingleConsumerExample.TransactOpts, amount)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactor) Unsubscribe(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.contract.Transact(opts, "unsubscribe", to)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) Unsubscribe(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.Unsubscribe(&_VRFV2PlusSingleConsumerExample.TransactOpts, to)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactorSession) Unsubscribe(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.Unsubscribe(&_VRFV2PlusSingleConsumerExample.TransactOpts, to)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactor) Withdraw(opts *bind.TransactOpts, amount *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.contract.Transact(opts, "withdraw", amount, to)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) Withdraw(amount *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.Withdraw(&_VRFV2PlusSingleConsumerExample.TransactOpts, amount, to)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactorSession) Withdraw(amount *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.Withdraw(&_VRFV2PlusSingleConsumerExample.TransactOpts, amount, to)
}

type SRequestConfig struct {
	SubId                uint64
	CallbackGasLimit     uint32
	RequestConfirmations uint16
	NumWords             uint32
	KeyHash              [32]byte
	NativePayment        bool
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExample) Address() common.Address {
	return _VRFV2PlusSingleConsumerExample.address
}

type VRFV2PlusSingleConsumerExampleInterface interface {
	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestConfig(opts *bind.CallOpts) (SRequestConfig,

		error)

	SRequestId(opts *bind.CallOpts) (*big.Int, error)

	FundAndRequestRandomWords(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts) (*types.Transaction, error)

	SetVRFCoordinator(opts *bind.TransactOpts, _vrfCoordinator common.Address) (*types.Transaction, error)

	Subscribe(opts *bind.TransactOpts) (*types.Transaction, error)

	TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	Unsubscribe(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Withdraw(opts *bind.TransactOpts, amount *big.Int, to common.Address) (*types.Transaction, error)

	Address() common.Address
}
