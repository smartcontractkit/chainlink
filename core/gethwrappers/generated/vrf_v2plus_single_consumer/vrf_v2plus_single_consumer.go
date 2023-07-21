// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_v2plus_single_consumer

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

var VRFV2PlusSingleConsumerExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"nativePayment\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"}],\"name\":\"OnlyOwnerOrCoordinator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"fundAndRequestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestConfig\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"nativePayment\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"_subId\",\"type\":\"uint64\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"subscribe\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"topUpSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"unsubscribe\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001a0738038062001a078339810160408190526200003491620004a1565b8633806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf81620001cc565b5050600280546001600160a01b03199081166001600160a01b03948516179091556003805482168b8516179055600480548216938a169390931790925550600a80543392169190911790556040805160c0810182526000815263ffffffff8781166020830181905261ffff8816938301849052908616606083018190526080830186905284151560a0909301839052600580546001600160701b0319166801000000000000000090930261ffff60601b1916929092176c010000000000000000000000009094029390931763ffffffff60701b1916600160701b9093029290921790915560068390556007805460ff19169091179055620001bf62000278565b5050505050505062000585565b6001600160a01b038116331415620002275760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6200028262000411565b604080516001808252818301909252600091602080830190803683370190505090503081600081518110620002bb57620002bb6200056f565b6001600160a01b039283166020918202929092018101919091526003546040805163288688f960e21b81529051919093169263a21a23e49260048083019391928290030181600087803b1580156200031257600080fd5b505af115801562000327573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200034d91906200053d565b600580546001600160401b0319166001600160401b0392909216918217905560035482516001600160a01b0390911691637341c10c9184906000906200039757620003976200056f565b60200260200101516040518363ffffffff1660e01b8152600401620003da9291906001600160401b039290921682526001600160a01b0316602082015260400190565b600060405180830381600087803b158015620003f557600080fd5b505af11580156200040a573d6000803e3d6000fd5b5050505050565b6000546001600160a01b031633146200046d5760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000083565b565b80516001600160a01b03811681146200048757600080fd5b919050565b805163ffffffff811681146200048757600080fd5b600080600080600080600060e0888a031215620004bd57600080fd5b620004c8886200046f565b9650620004d8602089016200046f565b9550620004e8604089016200048c565b9450606088015161ffff811681146200050057600080fd5b935062000510608089016200048c565b925060a0880151915060c088015180151581146200052d57600080fd5b8091505092959891949750929550565b6000602082840312156200055057600080fd5b81516001600160401b03811681146200056857600080fd5b9392505050565b634e487b7160e01b600052603260045260246000fd5b61147280620005956000396000f3fe608060405234801561001057600080fd5b50600436106100e95760003560e01c806386850e931161008c578063e0c8628911610066578063e0c8628914610248578063e89e106a14610250578063f2fde38b14610267578063f6eaffc81461027a57600080fd5b806386850e93146102055780638da5cb5b146102185780638f449a051461024057600080fd5b80636fd700bb116100c85780636fd700bb146101295780637262561c1461013c57806379ba50971461014f5780637db9263f1461015757600080fd5b8062f714ce146100ee5780631fe543e3146101035780632d6d99f314610116575b600080fd5b6101016100fc3660046111a1565b61028d565b005b6101016101113660046111cd565b610348565b610101610124366004611116565b6103ce565b61010161013736600461116f565b6104fe565b61010161014a3660046110f4565b610777565b610101610843565b6005546006546007546101b69267ffffffffffffffff81169263ffffffff68010000000000000000830481169361ffff6c01000000000000000000000000850416936e0100000000000000000000000000009004909116919060ff1686565b6040805167ffffffffffffffff909716875263ffffffff958616602088015261ffff90941693860193909352921660608401526080830191909152151560a082015260c0015b60405180910390f35b61010161021336600461116f565b610940565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101fc565b610101610a09565b610101610bed565b61025960095481565b6040519081526020016101fc565b6101016102753660046110f4565b610d83565b61025961028836600461116f565b610d97565b610295610db8565b600480546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff848116938201939093526024810185905291169063a9059cbb90604401602060405180830381600087803b15801561030b57600080fd5b505af115801561031f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610343919061114d565b505050565b60025473ffffffffffffffffffffffffffffffffffffffff1633146103c0576002546040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff90911660248201526044015b60405180910390fd5b6103ca8282610e3b565b5050565b60005473ffffffffffffffffffffffffffffffffffffffff16331480159061040e575060025473ffffffffffffffffffffffffffffffffffffffff163314155b15610492573361043360005473ffffffffffffffffffffffffffffffffffffffff1690565b6002546040517f061db9c100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff938416600482015291831660248301529190911660448201526064016103b7565b6002805467ffffffffffffffff90921674010000000000000000000000000000000000000000027fffffffff0000000000000000000000000000000000000000000000000000000090921673ffffffffffffffffffffffffffffffffffffffff90931692909217179055565b610506610db8565b6040805160c08101825260055467ffffffffffffffff811680835263ffffffff680100000000000000008304811660208086019190915261ffff6c01000000000000000000000000850416858701526e010000000000000000000000000000909304166060840152600654608084015260075460ff16151560a084015260045460035485518085019390935285518084039094018452828601958690527f4000aea000000000000000000000000000000000000000000000000000000000909552929373ffffffffffffffffffffffffffffffffffffffff93841693634000aea0936105fb9391909216918791604401611344565b602060405180830381600087803b15801561061557600080fd5b505af1158015610629573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061064d919061114d565b5060006040518060c0016040528083608001518152602001836000015167ffffffffffffffff168152602001836040015161ffff168152602001836020015163ffffffff168152602001836060015163ffffffff1681526020016106c460405180602001604052808660a001511515815250610eb9565b90526003546040517f596b8b8800000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff169063596b8b889061071d908490600401611382565b602060405180830381600087803b15801561073757600080fd5b505af115801561074b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061076f9190611188565b600955505050565b61077f610db8565b6003546005546040517fd7ae1d3000000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015273ffffffffffffffffffffffffffffffffffffffff83811660248301529091169063d7ae1d3090604401600060405180830381600087803b15801561080057600080fd5b505af1158015610814573d6000803e3d6000fd5b5050600580547fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000169055505050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146108c4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016103b7565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610948610db8565b6004546003546005546040805167ffffffffffffffff909216602083015273ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918591016040516020818303038152906040526040518463ffffffff1660e01b81526004016109b793929190611344565b602060405180830381600087803b1580156109d157600080fd5b505af11580156109e5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103ca919061114d565b610a11610db8565b604080516001808252818301909252600091602080830190803683370190505090503081600081518110610a4757610a476113f1565b73ffffffffffffffffffffffffffffffffffffffff928316602091820292909201810191909152600354604080517fa21a23e40000000000000000000000000000000000000000000000000000000081529051919093169263a21a23e49260048083019391928290030181600087803b158015610ac357600080fd5b505af1158015610ad7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610afb91906112bc565b600580547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff929092169182179055600354825173ffffffffffffffffffffffffffffffffffffffff90911691637341c10c918490600090610b6857610b686113f1565b60200260200101516040518363ffffffff1660e01b8152600401610bb892919067ffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b158015610bd257600080fd5b505af1158015610be6573d6000803e3d6000fd5b5050505050565b610bf5610db8565b6040805160c0808201835260055467ffffffffffffffff808216845263ffffffff6801000000000000000083048116602080870191825261ffff6c0100000000000000000000000086048116888a019081526e01000000000000000000000000000090960484166060808a019182526006546080808c0191825260075460ff16151560a0808e019182528e519c8d018f5292518c528c519099168b8701529851909316898c01529351851693880193909352915190921693850193909352855190810190955251151584529192600092820190610cd190610eb9565b90526003546040517f596b8b8800000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff169063596b8b8890610d2a908490600401611382565b602060405180830381600087803b158015610d4457600080fd5b505af1158015610d58573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d7c9190611188565b6009555050565b610d8b610db8565b610d9481610f75565b50565b60088181548110610da757600080fd5b600091825260209091200154905081565b60005473ffffffffffffffffffffffffffffffffffffffff163314610e39576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016103b7565b565b6009548214610ea6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265717565737420494420697320696e636f727265637400000000000000000060448201526064016103b7565b805161034390600890602084019061106b565b60607f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa82604051602401610ef291511515815260200190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915292915050565b73ffffffffffffffffffffffffffffffffffffffff8116331415610ff5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016103b7565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b8280548282559060005260206000209081019282156110a6579160200282015b828111156110a657825182559160200191906001019061108b565b506110b29291506110b6565b5090565b5b808211156110b257600081556001016110b7565b803573ffffffffffffffffffffffffffffffffffffffff811681146110ef57600080fd5b919050565b60006020828403121561110657600080fd5b61110f826110cb565b9392505050565b6000806040838503121561112957600080fd5b611132836110cb565b915060208301356111428161144f565b809150509250929050565b60006020828403121561115f57600080fd5b8151801515811461110f57600080fd5b60006020828403121561118157600080fd5b5035919050565b60006020828403121561119a57600080fd5b5051919050565b600080604083850312156111b457600080fd5b823591506111c4602084016110cb565b90509250929050565b600080604083850312156111e057600080fd5b8235915060208084013567ffffffffffffffff8082111561120057600080fd5b818601915086601f83011261121457600080fd5b81358181111561122657611226611420565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f8301168101818110858211171561126957611269611420565b604052828152858101935084860182860187018b101561128857600080fd5b600095505b838610156112ab57803585526001959095019493860193860161128d565b508096505050505050509250929050565b6000602082840312156112ce57600080fd5b815161110f8161144f565b6000815180845260005b818110156112ff576020818501810151868301820152016112e3565b81811115611311576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b73ffffffffffffffffffffffffffffffffffffffff8416815282602082015260606040820152600061137960608301846112d9565b95945050505050565b602081528151602082015267ffffffffffffffff602083015116604082015261ffff60408301511660608201526000606083015163ffffffff80821660808501528060808601511660a0850152505060a083015160c0808401526113e960e08401826112d9565b949350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b67ffffffffffffffff81168114610d9457600080fdfea164736f6c6343000806000a",
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

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusSingleConsumerExample.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) Owner() (common.Address, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.Owner(&_VRFV2PlusSingleConsumerExample.CallOpts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleCallerSession) Owner() (common.Address, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.Owner(&_VRFV2PlusSingleConsumerExample.CallOpts)
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

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.contract.Transact(opts, "acceptOwnership")
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.AcceptOwnership(&_VRFV2PlusSingleConsumerExample.TransactOpts)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.AcceptOwnership(&_VRFV2PlusSingleConsumerExample.TransactOpts)
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

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactor) SetConfig(opts *bind.TransactOpts, _vrfCoordinator common.Address, _subId uint64) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.contract.Transact(opts, "setConfig", _vrfCoordinator, _subId)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) SetConfig(_vrfCoordinator common.Address, _subId uint64) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.SetConfig(&_VRFV2PlusSingleConsumerExample.TransactOpts, _vrfCoordinator, _subId)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactorSession) SetConfig(_vrfCoordinator common.Address, _subId uint64) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.SetConfig(&_VRFV2PlusSingleConsumerExample.TransactOpts, _vrfCoordinator, _subId)
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

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.TransferOwnership(&_VRFV2PlusSingleConsumerExample.TransactOpts, to)
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusSingleConsumerExample.Contract.TransferOwnership(&_VRFV2PlusSingleConsumerExample.TransactOpts, to)
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

type VRFV2PlusSingleConsumerExampleOwnershipTransferRequestedIterator struct {
	Event *VRFV2PlusSingleConsumerExampleOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusSingleConsumerExampleOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusSingleConsumerExampleOwnershipTransferRequested)
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
		it.Event = new(VRFV2PlusSingleConsumerExampleOwnershipTransferRequested)
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

func (it *VRFV2PlusSingleConsumerExampleOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusSingleConsumerExampleOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusSingleConsumerExampleOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusSingleConsumerExampleOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusSingleConsumerExample.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusSingleConsumerExampleOwnershipTransferRequestedIterator{contract: _VRFV2PlusSingleConsumerExample.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2PlusSingleConsumerExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusSingleConsumerExample.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusSingleConsumerExampleOwnershipTransferRequested)
				if err := _VRFV2PlusSingleConsumerExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFV2PlusSingleConsumerExampleOwnershipTransferRequested, error) {
	event := new(VRFV2PlusSingleConsumerExampleOwnershipTransferRequested)
	if err := _VRFV2PlusSingleConsumerExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2PlusSingleConsumerExampleOwnershipTransferredIterator struct {
	Event *VRFV2PlusSingleConsumerExampleOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusSingleConsumerExampleOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusSingleConsumerExampleOwnershipTransferred)
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
		it.Event = new(VRFV2PlusSingleConsumerExampleOwnershipTransferred)
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

func (it *VRFV2PlusSingleConsumerExampleOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusSingleConsumerExampleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusSingleConsumerExampleOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusSingleConsumerExampleOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusSingleConsumerExample.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusSingleConsumerExampleOwnershipTransferredIterator{contract: _VRFV2PlusSingleConsumerExample.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2PlusSingleConsumerExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusSingleConsumerExample.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusSingleConsumerExampleOwnershipTransferred)
				if err := _VRFV2PlusSingleConsumerExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExampleFilterer) ParseOwnershipTransferred(log types.Log) (*VRFV2PlusSingleConsumerExampleOwnershipTransferred, error) {
	event := new(VRFV2PlusSingleConsumerExampleOwnershipTransferred)
	if err := _VRFV2PlusSingleConsumerExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SRequestConfig struct {
	SubId                uint64
	CallbackGasLimit     uint32
	RequestConfirmations uint16
	NumWords             uint32
	KeyHash              [32]byte
	NativePayment        bool
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExample) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFV2PlusSingleConsumerExample.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFV2PlusSingleConsumerExample.ParseOwnershipTransferRequested(log)
	case _VRFV2PlusSingleConsumerExample.abi.Events["OwnershipTransferred"].ID:
		return _VRFV2PlusSingleConsumerExample.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFV2PlusSingleConsumerExampleOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFV2PlusSingleConsumerExampleOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_VRFV2PlusSingleConsumerExample *VRFV2PlusSingleConsumerExample) Address() common.Address {
	return _VRFV2PlusSingleConsumerExample.address
}

type VRFV2PlusSingleConsumerExampleInterface interface {
	Owner(opts *bind.CallOpts) (common.Address, error)

	SRandomWords(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	SRequestConfig(opts *bind.CallOpts) (SRequestConfig,

		error)

	SRequestId(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	FundAndRequestRandomWords(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, _vrfCoordinator common.Address, _subId uint64) (*types.Transaction, error)

	Subscribe(opts *bind.TransactOpts) (*types.Transaction, error)

	TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Unsubscribe(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Withdraw(opts *bind.TransactOpts, amount *big.Int, to common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusSingleConsumerExampleOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2PlusSingleConsumerExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFV2PlusSingleConsumerExampleOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusSingleConsumerExampleOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2PlusSingleConsumerExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFV2PlusSingleConsumerExampleOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
