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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"nativePayment\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"fundAndRequestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_randomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestConfig\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"nativePayment\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"subscribe\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"topUpSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"unsubscribe\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b50604051620013ab380380620013ab83398101604081905262000034916200032f565b606087811b6001600160601b0319166080908152600080546001600160a01b03808c166001600160a01b0319928316178355600180548316918c16919091179055600780543392169190911790556040805160c08101825291825263ffffffff8981166020840181905261ffff8a1692840183905290881694830185905292820186905284151560a0909201829052600280546001600160701b0319166801000000000000000090940261ffff60601b1916939093176c010000000000000000000000009091021763ffffffff60701b1916600160701b90930292909217905560038390556004805460ff19169091179055620001306200013d565b5050505050505062000413565b6007546001600160a01b031633146200015557600080fd5b6040805160018082528183019092526000916020808301908036833701905050905030816000815181106200018e576200018e620003fd565b60200260200101906001600160a01b031690816001600160a01b03168152505060008054906101000a90046001600160a01b03166001600160a01b031663a21a23e46040518163ffffffff1660e01b8152600401602060405180830381600087803b158015620001fd57600080fd5b505af115801562000212573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002389190620003cb565b600280546001600160401b0319166001600160401b039290921691821790556000805483516001600160a01b0390911692637341c10c929091859190620002835762000283620003fd565b60200260200101516040518363ffffffff1660e01b8152600401620002c69291906001600160401b039290921682526001600160a01b0316602082015260400190565b600060405180830381600087803b158015620002e157600080fd5b505af1158015620002f6573d6000803e3d6000fd5b5050505050565b80516001600160a01b03811681146200031557600080fd5b919050565b805163ffffffff811681146200031557600080fd5b600080600080600080600060e0888a0312156200034b57600080fd5b6200035688620002fd565b96506200036660208901620002fd565b955062000376604089016200031a565b9450606088015161ffff811681146200038e57600080fd5b93506200039e608089016200031a565b925060a0880151915060c08801518015158114620003bb57600080fd5b8091505092959891949750929550565b600060208284031215620003de57600080fd5b81516001600160401b0381168114620003f657600080fd5b9392505050565b634e487b7160e01b600052603260045260246000fd5b60805160601c610f7262000439600039600081816102f801526103600152610f726000f3fe608060405234801561001057600080fd5b50600436106100bd5760003560e01c806386850e9311610076578063e0c862891161005b578063e0c86289146101d9578063e89e106a146101e1578063f6eaffc8146101f857600080fd5b806386850e93146101be5780638f449a05146101d157600080fd5b80636fd700bb116100a75780636fd700bb146100ea5780637262561c146100fd5780637db9263f1461011057600080fd5b8062f714ce146100c25780631fe543e3146100d7575b600080fd5b6100d56100d0366004610d2a565b61020b565b005b6100d56100e5366004610d56565b6102e0565b6100d56100f8366004610cf8565b6103a0565b6100d561010b366004610cb4565b61060b565b60025460035460045461016f9267ffffffffffffffff81169263ffffffff68010000000000000000830481169361ffff6c01000000000000000000000000850416936e0100000000000000000000000000009004909116919060ff1686565b6040805167ffffffffffffffff909716875263ffffffff958616602088015261ffff90941693860193909352921660608401526080830191909152151560a082015260c0015b60405180910390f35b6100d56101cc366004610cf8565b6106f3565b6100d56107d8565b6100d56109ff565b6101ea60065481565b6040519081526020016101b5565b6101ea610206366004610cf8565b610b77565b60075473ffffffffffffffffffffffffffffffffffffffff16331461022f57600080fd5b6001546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8381166004830152602482018590529091169063a9059cbb90604401602060405180830381600087803b1580156102a357600080fd5b505af11580156102b7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102db9190610cd6565b505050565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610392576040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660248201526044015b60405180910390fd5b61039c8282610b98565b5050565b60075473ffffffffffffffffffffffffffffffffffffffff1633146103c457600080fd5b6040805160c08101825260025467ffffffffffffffff811680835263ffffffff680100000000000000008304811660208086019190915261ffff6c01000000000000000000000000850416858701526e010000000000000000000000000000909304166060840152600354608084015260045460ff16151560a084015260015460005485518085019390935285518084039094018452828601958690527f4000aea000000000000000000000000000000000000000000000000000000000909552929373ffffffffffffffffffffffffffffffffffffffff93841693634000aea0936104b99391909216918791604401610e6f565b602060405180830381600087803b1580156104d357600080fd5b505af11580156104e7573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061050b9190610cd6565b50600054608082015182516040808501516020860151606087015160a088015193517fefcf1d94000000000000000000000000000000000000000000000000000000008152600481019690965267ffffffffffffffff909416602486015261ffff909116604485015263ffffffff90811660648501529091166084830152151560a482015273ffffffffffffffffffffffffffffffffffffffff9091169063efcf1d949060c401602060405180830381600087803b1580156105cc57600080fd5b505af11580156105e0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106049190610d11565b6006555050565b60075473ffffffffffffffffffffffffffffffffffffffff16331461062f57600080fd5b6000546002546040517fd7ae1d3000000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015273ffffffffffffffffffffffffffffffffffffffff83811660248301529091169063d7ae1d3090604401600060405180830381600087803b1580156106b057600080fd5b505af11580156106c4573d6000803e3d6000fd5b5050600280547fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000169055505050565b60075473ffffffffffffffffffffffffffffffffffffffff16331461071757600080fd5b6001546000546002546040805167ffffffffffffffff909216602083015273ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918591016040516020818303038152906040526040518463ffffffff1660e01b815260040161078693929190610e6f565b602060405180830381600087803b1580156107a057600080fd5b505af11580156107b4573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061039c9190610cd6565b60075473ffffffffffffffffffffffffffffffffffffffff1633146107fc57600080fd5b60408051600180825281830190925260009160208083019080368337019050509050308160008151811061083257610832610f07565b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff168152505060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a21a23e46040518163ffffffff1660e01b8152600401602060405180830381600087803b1580156108d457600080fd5b505af11580156108e8573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061090c9190610e45565b600280547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff92909216918217905560008054835173ffffffffffffffffffffffffffffffffffffffff90911692637341c10c92909185919061097a5761097a610f07565b60200260200101516040518363ffffffff1660e01b81526004016109ca92919067ffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b1580156109e457600080fd5b505af11580156109f8573d6000803e3d6000fd5b5050505050565b60075473ffffffffffffffffffffffffffffffffffffffff163314610a2357600080fd5b6040805160c08101825260025467ffffffffffffffff811680835263ffffffff68010000000000000000830481166020850181905261ffff6c010000000000000000000000008504168587018190526e01000000000000000000000000000090940490911660608501819052600354608086018190526004805460ff16151560a0880181905260005498517fefcf1d9400000000000000000000000000000000000000000000000000000000815291820192909252602481019490945260448401949094526064830191909152608482015260a4810191909152909173ffffffffffffffffffffffffffffffffffffffff169063efcf1d949060c401602060405180830381600087803b158015610b3957600080fd5b505af1158015610b4d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b719190610d11565b60065550565b60058181548110610b8757600080fd5b600091825260209091200154905081565b6006548214610c03576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265717565737420494420697320696e636f72726563740000000000000000006044820152606401610389565b80516005805482825560008290526102db927f036b6384b5eca791c62761152d0c79bb0604c104a5fb6f4eb0703f3154bb3db091820191602086018215610c66579160200282015b82811115610c66578251825591602001919060010190610c4b565b50610c72929150610c76565b5090565b5b80821115610c725760008155600101610c77565b803573ffffffffffffffffffffffffffffffffffffffff81168114610caf57600080fd5b919050565b600060208284031215610cc657600080fd5b610ccf82610c8b565b9392505050565b600060208284031215610ce857600080fd5b81518015158114610ccf57600080fd5b600060208284031215610d0a57600080fd5b5035919050565b600060208284031215610d2357600080fd5b5051919050565b60008060408385031215610d3d57600080fd5b82359150610d4d60208401610c8b565b90509250929050565b60008060408385031215610d6957600080fd5b8235915060208084013567ffffffffffffffff80821115610d8957600080fd5b818601915086601f830112610d9d57600080fd5b813581811115610daf57610daf610f36565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f83011681018181108582111715610df257610df2610f36565b604052828152858101935084860182860187018b1015610e1157600080fd5b600095505b83861015610e34578035855260019590950194938601938601610e16565b508096505050505050509250929050565b600060208284031215610e5757600080fd5b815167ffffffffffffffff81168114610ccf57600080fd5b73ffffffffffffffffffffffffffffffffffffffff8416815260006020848184015260606040840152835180606085015260005b81811015610ebf57858101830151858201608001528201610ea3565b81811115610ed1576000608083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160800195945050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
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

	Subscribe(opts *bind.TransactOpts) (*types.Transaction, error)

	TopUpSubscription(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	Unsubscribe(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Withdraw(opts *bind.TransactOpts, amount *big.Int, to common.Address) (*types.Transaction, error)

	Address() common.Address
}
