// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package consumer_wrapper

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
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated"
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

var ConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_specId\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ChainlinkCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ChainlinkFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ChainlinkRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"price\",\"type\":\"bytes32\"}],\"name\":\"RequestFulfilled\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_requestId\",\"type\":\"bytes32\"}],\"name\":\"addExternalRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_payment\",\"type\":\"uint256\"},{\"internalType\":\"bytes4\",\"name\":\"_callbackFunctionId\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"_expiration\",\"type\":\"uint256\"}],\"name\":\"cancelRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentPrice\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentPriceInt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_price\",\"type\":\"bytes32\"}],\"name\":\"fulfill\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"}],\"name\":\"fulfillParametersWithCustomURLs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_currency\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_payment\",\"type\":\"uint256\"}],\"name\":\"requestEthereumPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_urlUSD\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_pathUSD\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_payment\",\"type\":\"uint256\"}],\"name\":\"requestMultipleParametersWithCustomURLs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_specId\",\"type\":\"bytes32\"}],\"name\":\"setSpecID\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052600160045534801561001557600080fd5b5060405161158b38038061158b8339818101604052606081101561003857600080fd5b508051602082015160409092015190919061005283610066565b61005b82610088565b600655506100aa9050565b600280546001600160a01b0319166001600160a01b0392909216919091179055565b600380546001600160a01b0319166001600160a01b0392909216919091179055565b6114d2806100b96000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c806383db5cbc116100765780639438a6011161005b5780639438a601146103745780639d1b464a1461038e578063e8d5359d14610396576100be565b806383db5cbc146102c45780638dc654a21461036c576100be565b80635591a608116100a75780635591a608146101055780635b8260051461017257806371c2002a14610195576100be565b8063042f2b65146100c3578063501fdd5d146100e8575b600080fd5b6100e6600480360360408110156100d957600080fd5b50803590602001356103cf565b005b6100e6600480360360208110156100fe57600080fd5b50356104dc565b6100e6600480360360a081101561011b57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813516906020810135906040810135907fffffffff0000000000000000000000000000000000000000000000000000000060608201351690608001356104e1565b6100e66004803603604081101561018857600080fd5b50803590602001356105a8565b6100e6600480360360608110156101ab57600080fd5b8101906020810181356401000000008111156101c657600080fd5b8201836020820111156101d857600080fd5b803590602001918460018302840111640100000000831117156101fa57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929594936020810193503591505064010000000081111561024d57600080fd5b82018360208201111561025f57600080fd5b8035906020019184600183028401116401000000008311171561028157600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092955050913592506106b5915050565b6100e6600480360360408110156102da57600080fd5b8101906020810181356401000000008111156102f557600080fd5b82018360208201111561030757600080fd5b8035906020019184600183028401116401000000008311171561032957600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550509135925061075e915050565b6100e661086a565b61037c610a34565b60408051918252519081900360200190f35b61037c610a3a565b6100e6600480360360408110156103ac57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135610a40565b600082815260056020526040902054829073ffffffffffffffffffffffffffffffffffffffff16331461044d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260288152602001806114576028913960400191505060405180910390fd5b60008181526005602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001690555182917f7cc135e0cebb02c3480ae5d74d377283180a2601f8f644edf7987b009316c63a91a2604051829084907f0c2366233f634048c0f0458060d1228fab36d00f7c0ecf6bdf2d9c458503631190600090a35060075550565b600655565b604080517f6ee4d55300000000000000000000000000000000000000000000000000000000815260048101869052602481018590527fffffffff0000000000000000000000000000000000000000000000000000000084166044820152606481018390529051869173ffffffffffffffffffffffffffffffffffffffff831691636ee4d5539160848082019260009290919082900301818387803b15801561058857600080fd5b505af115801561059c573d6000803e3d6000fd5b50505050505050505050565b600082815260056020526040902054829073ffffffffffffffffffffffffffffffffffffffff163314610626576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260288152602001806114576028913960400191505060405180910390fd5b60008181526005602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001690555182917f7cc135e0cebb02c3480ae5d74d377283180a2601f8f644edf7987b009316c63a91a2604051829084907f0c2366233f634048c0f0458060d1228fab36d00f7c0ecf6bdf2d9c458503631190600090a35060085550565b60006106ca600654635b82600560e01b610a4e565b60408051808201909152600681527f75726c5553440000000000000000000000000000000000000000000000000000602082015290915061070d90829086610a74565b60408051808201909152600781527f7061746855534400000000000000000000000000000000000000000000000000602082015261074d90829085610a74565b6107578183610a97565b5050505050565b600061077360065463042f2b6560e01b610a4e565b90506107cf6040518060400160405280600381526020017f676574000000000000000000000000000000000000000000000000000000000081525060405180608001604052806047815260200161147f60479139839190610a74565b604080516001808252818301909252600091816020015b60608152602001906001900390816107e6579050509050838160008151811061080b57fe5b60200260200101819052506108606040518060400160405280600481526020017f70617468000000000000000000000000000000000000000000000000000000008152508284610ac59092919063ffffffff16565b6107578284610a97565b6000610874610b2d565b90508073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb338373ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b1580156108fa57600080fd5b505afa15801561090e573d6000803e3d6000fd5b505050506040513d602081101561092457600080fd5b5051604080517fffffffff0000000000000000000000000000000000000000000000000000000060e086901b16815273ffffffffffffffffffffffffffffffffffffffff909316600484015260248301919091525160448083019260209291908290030181600087803b15801561099a57600080fd5b505af11580156109ae573d6000803e3d6000fd5b505050506040513d60208110156109c457600080fd5b5051610a3157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f556e61626c6520746f207472616e736665720000000000000000000000000000604482015290519081900360640190fd5b50565b60085481565b60075481565b610a4a8282610b49565b5050565b610a566113e4565b610a5e6113e4565b610a6a81853086610c30565b9150505b92915050565b6080830151610a839083610c9b565b6080830151610a929082610c9b565b505050565b600354600090610abe9073ffffffffffffffffffffffffffffffffffffffff168484610cb2565b9392505050565b6080830151610ad49083610c9b565b610ae18360800151610e4b565b60005b8151811015610b1f57610b17828281518110610afc57fe5b60200260200101518560800151610c9b90919063ffffffff16565b600101610ae4565b50610a928360800151610e56565b60025473ffffffffffffffffffffffffffffffffffffffff1690565b600081815260056020526040902054819073ffffffffffffffffffffffffffffffffffffffff1615610bdc57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f5265717565737420697320616c72656164792070656e64696e67000000000000604482015290519081900360640190fd5b50600090815260056020526040902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b610c386113e4565b610c488560800151610100610e61565b505082845273ffffffffffffffffffffffffffffffffffffffff821660208501527fffffffff0000000000000000000000000000000000000000000000000000000081166040850152835b949350505050565b610ca88260038351610e9b565b610a928282610fb6565b6000806004549050806001016004819055506000634042994660e01b60008087600001513089604001518760018c6080015160000151604051602401808973ffffffffffffffffffffffffffffffffffffffff1681526020018881526020018781526020018673ffffffffffffffffffffffffffffffffffffffff168152602001857bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200184815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b83811015610d9d578181015183820152602001610d85565b50505050905090810190601f168015610dca5780820380516001836020036101000a031916815260200191505b509950505050505050505050604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050509050610e4186838684610fd0565b9695505050505050565b610a318160046111e3565b610a318160076111e3565b610e69611419565b6020820615610e7e5760208206602003820191505b506020828101829052604080518085526000815290920101905290565b60178167ffffffffffffffff1611610ec657610ec08360e0600585901b1683176111f4565b50610a92565b60ff8167ffffffffffffffff1611610f0457610eed836018611fe0600586901b16176111f4565b50610ec08367ffffffffffffffff8316600161120c565b61ffff8167ffffffffffffffff1611610f4357610f2c836019611fe0600586901b16176111f4565b50610ec08367ffffffffffffffff8316600261120c565b63ffffffff8167ffffffffffffffff1611610f8457610f6d83601a611fe0600586901b16176111f4565b50610ec08367ffffffffffffffff8316600461120c565b610f9983601b611fe0600586901b16176111f4565b50610fb08367ffffffffffffffff8316600861120c565b50505050565b610fbe611419565b610abe83846000015151848551611225565b604080513060601b60208083019190915260348083018790528351808403909101815260549092018084528251928201929092206000818152600590925292812080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff891617905582917fb5e6e01e79f91267dc17b4e6314d5d4d03593d2ceee0fbb452b750bd70ea5af99190a26002546040517f4000aea000000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff87811660048301908152602483018790526060604484019081528651606485015286519290941693634000aea0938a938993899390929091608490910190602085019080838360005b838110156111145781810151838201526020016110fc565b50505050905090810190601f1680156111415780820380516001836020036101000a031916815260200191505b50945050505050602060405180830381600087803b15801561116257600080fd5b505af1158015611176573d6000803e3d6000fd5b505050506040513d602081101561118c57600080fd5b5051610c93576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260238152602001806114346023913960400191505060405180910390fd5b610a9282601f611fe0600585901b16175b6111fc611419565b610abe838460000151518461130d565b611214611419565b610c93848560000151518585611358565b61122d611419565b825182111561123b57600080fd5b84602001518285011115611265576112658561125d87602001518786016113b6565b6002026113cd565b6000808651805187602083010193508088870111156112845787860182525b505050602084015b602084106112c957805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0909301926020918201910161128c565b5181517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60208690036101000a019081169019919091161790525083949350505050565b611315611419565b83602001518310611331576113318485602001516002026113cd565b83518051602085830101848153508085141561134e576001810182525b5093949350505050565b611360611419565b8460200151848301111561137d5761137d858584016002026113cd565b60006001836101000a0390508551838682010185831982511617815250805184870111156113ab5783860181525b509495945050505050565b6000818311156113c7575081610a6e565b50919050565b81516113d98383610e61565b50610fb08382610fb6565b6040805160a081018252600080825260208201819052918101829052606081019190915260808101611414611419565b905290565b60405180604001604052806060815260200160008152509056fe756e61626c6520746f207472616e73666572416e6443616c6c20746f206f7261636c65536f75726365206d75737420626520746865206f7261636c65206f6620746865207265717565737468747470733a2f2f6d696e2d6170692e63727970746f636f6d706172652e636f6d2f646174612f70726963653f6673796d3d455448267473796d733d5553442c4555522c4a5059a164736f6c6343000706000a",
}

var ConsumerABI = ConsumerMetaData.ABI

var ConsumerBin = ConsumerMetaData.Bin

func DeployConsumer(auth *bind.TransactOpts, backend bind.ContractBackend, _link common.Address, _oracle common.Address, _specId [32]byte) (common.Address, *types.Transaction, *Consumer, error) {
	parsed, err := ConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ConsumerBin), backend, _link, _oracle, _specId)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Consumer{ConsumerCaller: ConsumerCaller{contract: contract}, ConsumerTransactor: ConsumerTransactor{contract: contract}, ConsumerFilterer: ConsumerFilterer{contract: contract}}, nil
}

type Consumer struct {
	address common.Address
	abi     abi.ABI
	ConsumerCaller
	ConsumerTransactor
	ConsumerFilterer
}

type ConsumerCaller struct {
	contract *bind.BoundContract
}

type ConsumerTransactor struct {
	contract *bind.BoundContract
}

type ConsumerFilterer struct {
	contract *bind.BoundContract
}

type ConsumerSession struct {
	Contract     *Consumer
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ConsumerCallerSession struct {
	Contract *ConsumerCaller
	CallOpts bind.CallOpts
}

type ConsumerTransactorSession struct {
	Contract     *ConsumerTransactor
	TransactOpts bind.TransactOpts
}

type ConsumerRaw struct {
	Contract *Consumer
}

type ConsumerCallerRaw struct {
	Contract *ConsumerCaller
}

type ConsumerTransactorRaw struct {
	Contract *ConsumerTransactor
}

func NewConsumer(address common.Address, backend bind.ContractBackend) (*Consumer, error) {
	abi, err := abi.JSON(strings.NewReader(ConsumerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Consumer{address: address, abi: abi, ConsumerCaller: ConsumerCaller{contract: contract}, ConsumerTransactor: ConsumerTransactor{contract: contract}, ConsumerFilterer: ConsumerFilterer{contract: contract}}, nil
}

func NewConsumerCaller(address common.Address, caller bind.ContractCaller) (*ConsumerCaller, error) {
	contract, err := bindConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConsumerCaller{contract: contract}, nil
}

func NewConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*ConsumerTransactor, error) {
	contract, err := bindConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConsumerTransactor{contract: contract}, nil
}

func NewConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*ConsumerFilterer, error) {
	contract, err := bindConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConsumerFilterer{contract: contract}, nil
}

func bindConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ConsumerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_Consumer *ConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Consumer.Contract.ConsumerCaller.contract.Call(opts, result, method, params...)
}

func (_Consumer *ConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Consumer.Contract.ConsumerTransactor.contract.Transfer(opts)
}

func (_Consumer *ConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Consumer.Contract.ConsumerTransactor.contract.Transact(opts, method, params...)
}

func (_Consumer *ConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Consumer.Contract.contract.Call(opts, result, method, params...)
}

func (_Consumer *ConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Consumer.Contract.contract.Transfer(opts)
}

func (_Consumer *ConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Consumer.Contract.contract.Transact(opts, method, params...)
}

func (_Consumer *ConsumerCaller) CurrentPrice(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Consumer.contract.Call(opts, &out, "currentPrice")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_Consumer *ConsumerSession) CurrentPrice() ([32]byte, error) {
	return _Consumer.Contract.CurrentPrice(&_Consumer.CallOpts)
}

func (_Consumer *ConsumerCallerSession) CurrentPrice() ([32]byte, error) {
	return _Consumer.Contract.CurrentPrice(&_Consumer.CallOpts)
}

func (_Consumer *ConsumerCaller) CurrentPriceInt(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Consumer.contract.Call(opts, &out, "currentPriceInt")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_Consumer *ConsumerSession) CurrentPriceInt() (*big.Int, error) {
	return _Consumer.Contract.CurrentPriceInt(&_Consumer.CallOpts)
}

func (_Consumer *ConsumerCallerSession) CurrentPriceInt() (*big.Int, error) {
	return _Consumer.Contract.CurrentPriceInt(&_Consumer.CallOpts)
}

func (_Consumer *ConsumerTransactor) AddExternalRequest(opts *bind.TransactOpts, _oracle common.Address, _requestId [32]byte) (*types.Transaction, error) {
	return _Consumer.contract.Transact(opts, "addExternalRequest", _oracle, _requestId)
}

func (_Consumer *ConsumerSession) AddExternalRequest(_oracle common.Address, _requestId [32]byte) (*types.Transaction, error) {
	return _Consumer.Contract.AddExternalRequest(&_Consumer.TransactOpts, _oracle, _requestId)
}

func (_Consumer *ConsumerTransactorSession) AddExternalRequest(_oracle common.Address, _requestId [32]byte) (*types.Transaction, error) {
	return _Consumer.Contract.AddExternalRequest(&_Consumer.TransactOpts, _oracle, _requestId)
}

func (_Consumer *ConsumerTransactor) CancelRequest(opts *bind.TransactOpts, _oracle common.Address, _requestId [32]byte, _payment *big.Int, _callbackFunctionId [4]byte, _expiration *big.Int) (*types.Transaction, error) {
	return _Consumer.contract.Transact(opts, "cancelRequest", _oracle, _requestId, _payment, _callbackFunctionId, _expiration)
}

func (_Consumer *ConsumerSession) CancelRequest(_oracle common.Address, _requestId [32]byte, _payment *big.Int, _callbackFunctionId [4]byte, _expiration *big.Int) (*types.Transaction, error) {
	return _Consumer.Contract.CancelRequest(&_Consumer.TransactOpts, _oracle, _requestId, _payment, _callbackFunctionId, _expiration)
}

func (_Consumer *ConsumerTransactorSession) CancelRequest(_oracle common.Address, _requestId [32]byte, _payment *big.Int, _callbackFunctionId [4]byte, _expiration *big.Int) (*types.Transaction, error) {
	return _Consumer.Contract.CancelRequest(&_Consumer.TransactOpts, _oracle, _requestId, _payment, _callbackFunctionId, _expiration)
}

func (_Consumer *ConsumerTransactor) Fulfill(opts *bind.TransactOpts, _requestId [32]byte, _price [32]byte) (*types.Transaction, error) {
	return _Consumer.contract.Transact(opts, "fulfill", _requestId, _price)
}

func (_Consumer *ConsumerSession) Fulfill(_requestId [32]byte, _price [32]byte) (*types.Transaction, error) {
	return _Consumer.Contract.Fulfill(&_Consumer.TransactOpts, _requestId, _price)
}

func (_Consumer *ConsumerTransactorSession) Fulfill(_requestId [32]byte, _price [32]byte) (*types.Transaction, error) {
	return _Consumer.Contract.Fulfill(&_Consumer.TransactOpts, _requestId, _price)
}

func (_Consumer *ConsumerTransactor) FulfillParametersWithCustomURLs(opts *bind.TransactOpts, _requestId [32]byte, _price *big.Int) (*types.Transaction, error) {
	return _Consumer.contract.Transact(opts, "fulfillParametersWithCustomURLs", _requestId, _price)
}

func (_Consumer *ConsumerSession) FulfillParametersWithCustomURLs(_requestId [32]byte, _price *big.Int) (*types.Transaction, error) {
	return _Consumer.Contract.FulfillParametersWithCustomURLs(&_Consumer.TransactOpts, _requestId, _price)
}

func (_Consumer *ConsumerTransactorSession) FulfillParametersWithCustomURLs(_requestId [32]byte, _price *big.Int) (*types.Transaction, error) {
	return _Consumer.Contract.FulfillParametersWithCustomURLs(&_Consumer.TransactOpts, _requestId, _price)
}

func (_Consumer *ConsumerTransactor) RequestEthereumPrice(opts *bind.TransactOpts, _currency string, _payment *big.Int) (*types.Transaction, error) {
	return _Consumer.contract.Transact(opts, "requestEthereumPrice", _currency, _payment)
}

func (_Consumer *ConsumerSession) RequestEthereumPrice(_currency string, _payment *big.Int) (*types.Transaction, error) {
	return _Consumer.Contract.RequestEthereumPrice(&_Consumer.TransactOpts, _currency, _payment)
}

func (_Consumer *ConsumerTransactorSession) RequestEthereumPrice(_currency string, _payment *big.Int) (*types.Transaction, error) {
	return _Consumer.Contract.RequestEthereumPrice(&_Consumer.TransactOpts, _currency, _payment)
}

func (_Consumer *ConsumerTransactor) RequestMultipleParametersWithCustomURLs(opts *bind.TransactOpts, _urlUSD string, _pathUSD string, _payment *big.Int) (*types.Transaction, error) {
	return _Consumer.contract.Transact(opts, "requestMultipleParametersWithCustomURLs", _urlUSD, _pathUSD, _payment)
}

func (_Consumer *ConsumerSession) RequestMultipleParametersWithCustomURLs(_urlUSD string, _pathUSD string, _payment *big.Int) (*types.Transaction, error) {
	return _Consumer.Contract.RequestMultipleParametersWithCustomURLs(&_Consumer.TransactOpts, _urlUSD, _pathUSD, _payment)
}

func (_Consumer *ConsumerTransactorSession) RequestMultipleParametersWithCustomURLs(_urlUSD string, _pathUSD string, _payment *big.Int) (*types.Transaction, error) {
	return _Consumer.Contract.RequestMultipleParametersWithCustomURLs(&_Consumer.TransactOpts, _urlUSD, _pathUSD, _payment)
}

func (_Consumer *ConsumerTransactor) SetSpecID(opts *bind.TransactOpts, _specId [32]byte) (*types.Transaction, error) {
	return _Consumer.contract.Transact(opts, "setSpecID", _specId)
}

func (_Consumer *ConsumerSession) SetSpecID(_specId [32]byte) (*types.Transaction, error) {
	return _Consumer.Contract.SetSpecID(&_Consumer.TransactOpts, _specId)
}

func (_Consumer *ConsumerTransactorSession) SetSpecID(_specId [32]byte) (*types.Transaction, error) {
	return _Consumer.Contract.SetSpecID(&_Consumer.TransactOpts, _specId)
}

func (_Consumer *ConsumerTransactor) WithdrawLink(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Consumer.contract.Transact(opts, "withdrawLink")
}

func (_Consumer *ConsumerSession) WithdrawLink() (*types.Transaction, error) {
	return _Consumer.Contract.WithdrawLink(&_Consumer.TransactOpts)
}

func (_Consumer *ConsumerTransactorSession) WithdrawLink() (*types.Transaction, error) {
	return _Consumer.Contract.WithdrawLink(&_Consumer.TransactOpts)
}

type ConsumerChainlinkCancelledIterator struct {
	Event *ConsumerChainlinkCancelled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ConsumerChainlinkCancelledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsumerChainlinkCancelled)
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
		it.Event = new(ConsumerChainlinkCancelled)
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

func (it *ConsumerChainlinkCancelledIterator) Error() error {
	return it.fail
}

func (it *ConsumerChainlinkCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ConsumerChainlinkCancelled struct {
	Id  [32]byte
	Raw types.Log
}

func (_Consumer *ConsumerFilterer) FilterChainlinkCancelled(opts *bind.FilterOpts, id [][32]byte) (*ConsumerChainlinkCancelledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Consumer.contract.FilterLogs(opts, "ChainlinkCancelled", idRule)
	if err != nil {
		return nil, err
	}
	return &ConsumerChainlinkCancelledIterator{contract: _Consumer.contract, event: "ChainlinkCancelled", logs: logs, sub: sub}, nil
}

func (_Consumer *ConsumerFilterer) WatchChainlinkCancelled(opts *bind.WatchOpts, sink chan<- *ConsumerChainlinkCancelled, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Consumer.contract.WatchLogs(opts, "ChainlinkCancelled", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ConsumerChainlinkCancelled)
				if err := _Consumer.contract.UnpackLog(event, "ChainlinkCancelled", log); err != nil {
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

func (_Consumer *ConsumerFilterer) ParseChainlinkCancelled(log types.Log) (*ConsumerChainlinkCancelled, error) {
	event := new(ConsumerChainlinkCancelled)
	if err := _Consumer.contract.UnpackLog(event, "ChainlinkCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ConsumerChainlinkFulfilledIterator struct {
	Event *ConsumerChainlinkFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ConsumerChainlinkFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsumerChainlinkFulfilled)
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
		it.Event = new(ConsumerChainlinkFulfilled)
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

func (it *ConsumerChainlinkFulfilledIterator) Error() error {
	return it.fail
}

func (it *ConsumerChainlinkFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ConsumerChainlinkFulfilled struct {
	Id  [32]byte
	Raw types.Log
}

func (_Consumer *ConsumerFilterer) FilterChainlinkFulfilled(opts *bind.FilterOpts, id [][32]byte) (*ConsumerChainlinkFulfilledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Consumer.contract.FilterLogs(opts, "ChainlinkFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return &ConsumerChainlinkFulfilledIterator{contract: _Consumer.contract, event: "ChainlinkFulfilled", logs: logs, sub: sub}, nil
}

func (_Consumer *ConsumerFilterer) WatchChainlinkFulfilled(opts *bind.WatchOpts, sink chan<- *ConsumerChainlinkFulfilled, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Consumer.contract.WatchLogs(opts, "ChainlinkFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ConsumerChainlinkFulfilled)
				if err := _Consumer.contract.UnpackLog(event, "ChainlinkFulfilled", log); err != nil {
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

func (_Consumer *ConsumerFilterer) ParseChainlinkFulfilled(log types.Log) (*ConsumerChainlinkFulfilled, error) {
	event := new(ConsumerChainlinkFulfilled)
	if err := _Consumer.contract.UnpackLog(event, "ChainlinkFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ConsumerChainlinkRequestedIterator struct {
	Event *ConsumerChainlinkRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ConsumerChainlinkRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsumerChainlinkRequested)
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
		it.Event = new(ConsumerChainlinkRequested)
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

func (it *ConsumerChainlinkRequestedIterator) Error() error {
	return it.fail
}

func (it *ConsumerChainlinkRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ConsumerChainlinkRequested struct {
	Id  [32]byte
	Raw types.Log
}

func (_Consumer *ConsumerFilterer) FilterChainlinkRequested(opts *bind.FilterOpts, id [][32]byte) (*ConsumerChainlinkRequestedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Consumer.contract.FilterLogs(opts, "ChainlinkRequested", idRule)
	if err != nil {
		return nil, err
	}
	return &ConsumerChainlinkRequestedIterator{contract: _Consumer.contract, event: "ChainlinkRequested", logs: logs, sub: sub}, nil
}

func (_Consumer *ConsumerFilterer) WatchChainlinkRequested(opts *bind.WatchOpts, sink chan<- *ConsumerChainlinkRequested, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Consumer.contract.WatchLogs(opts, "ChainlinkRequested", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ConsumerChainlinkRequested)
				if err := _Consumer.contract.UnpackLog(event, "ChainlinkRequested", log); err != nil {
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

func (_Consumer *ConsumerFilterer) ParseChainlinkRequested(log types.Log) (*ConsumerChainlinkRequested, error) {
	event := new(ConsumerChainlinkRequested)
	if err := _Consumer.contract.UnpackLog(event, "ChainlinkRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ConsumerRequestFulfilledIterator struct {
	Event *ConsumerRequestFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ConsumerRequestFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsumerRequestFulfilled)
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
		it.Event = new(ConsumerRequestFulfilled)
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

func (it *ConsumerRequestFulfilledIterator) Error() error {
	return it.fail
}

func (it *ConsumerRequestFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ConsumerRequestFulfilled struct {
	RequestId [32]byte
	Price     [32]byte
	Raw       types.Log
}

func (_Consumer *ConsumerFilterer) FilterRequestFulfilled(opts *bind.FilterOpts, requestId [][32]byte, price [][32]byte) (*ConsumerRequestFulfilledIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var priceRule []interface{}
	for _, priceItem := range price {
		priceRule = append(priceRule, priceItem)
	}

	logs, sub, err := _Consumer.contract.FilterLogs(opts, "RequestFulfilled", requestIdRule, priceRule)
	if err != nil {
		return nil, err
	}
	return &ConsumerRequestFulfilledIterator{contract: _Consumer.contract, event: "RequestFulfilled", logs: logs, sub: sub}, nil
}

func (_Consumer *ConsumerFilterer) WatchRequestFulfilled(opts *bind.WatchOpts, sink chan<- *ConsumerRequestFulfilled, requestId [][32]byte, price [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var priceRule []interface{}
	for _, priceItem := range price {
		priceRule = append(priceRule, priceItem)
	}

	logs, sub, err := _Consumer.contract.WatchLogs(opts, "RequestFulfilled", requestIdRule, priceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ConsumerRequestFulfilled)
				if err := _Consumer.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
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

func (_Consumer *ConsumerFilterer) ParseRequestFulfilled(log types.Log) (*ConsumerRequestFulfilled, error) {
	event := new(ConsumerRequestFulfilled)
	if err := _Consumer.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_Consumer *Consumer) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _Consumer.abi.Events["ChainlinkCancelled"].ID:
		return _Consumer.ParseChainlinkCancelled(log)
	case _Consumer.abi.Events["ChainlinkFulfilled"].ID:
		return _Consumer.ParseChainlinkFulfilled(log)
	case _Consumer.abi.Events["ChainlinkRequested"].ID:
		return _Consumer.ParseChainlinkRequested(log)
	case _Consumer.abi.Events["RequestFulfilled"].ID:
		return _Consumer.ParseRequestFulfilled(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (ConsumerChainlinkCancelled) Topic() common.Hash {
	return common.HexToHash("0xe1fe3afa0f7f761ff0a8b89086790efd5140d2907ebd5b7ff6bfcb5e075fd4c5")
}

func (ConsumerChainlinkFulfilled) Topic() common.Hash {
	return common.HexToHash("0x7cc135e0cebb02c3480ae5d74d377283180a2601f8f644edf7987b009316c63a")
}

func (ConsumerChainlinkRequested) Topic() common.Hash {
	return common.HexToHash("0xb5e6e01e79f91267dc17b4e6314d5d4d03593d2ceee0fbb452b750bd70ea5af9")
}

func (ConsumerRequestFulfilled) Topic() common.Hash {
	return common.HexToHash("0x0c2366233f634048c0f0458060d1228fab36d00f7c0ecf6bdf2d9c4585036311")
}

func (_Consumer *Consumer) Address() common.Address {
	return _Consumer.address
}

type ConsumerInterface interface {
	CurrentPrice(opts *bind.CallOpts) ([32]byte, error)

	CurrentPriceInt(opts *bind.CallOpts) (*big.Int, error)

	AddExternalRequest(opts *bind.TransactOpts, _oracle common.Address, _requestId [32]byte) (*types.Transaction, error)

	CancelRequest(opts *bind.TransactOpts, _oracle common.Address, _requestId [32]byte, _payment *big.Int, _callbackFunctionId [4]byte, _expiration *big.Int) (*types.Transaction, error)

	Fulfill(opts *bind.TransactOpts, _requestId [32]byte, _price [32]byte) (*types.Transaction, error)

	FulfillParametersWithCustomURLs(opts *bind.TransactOpts, _requestId [32]byte, _price *big.Int) (*types.Transaction, error)

	RequestEthereumPrice(opts *bind.TransactOpts, _currency string, _payment *big.Int) (*types.Transaction, error)

	RequestMultipleParametersWithCustomURLs(opts *bind.TransactOpts, _urlUSD string, _pathUSD string, _payment *big.Int) (*types.Transaction, error)

	SetSpecID(opts *bind.TransactOpts, _specId [32]byte) (*types.Transaction, error)

	WithdrawLink(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterChainlinkCancelled(opts *bind.FilterOpts, id [][32]byte) (*ConsumerChainlinkCancelledIterator, error)

	WatchChainlinkCancelled(opts *bind.WatchOpts, sink chan<- *ConsumerChainlinkCancelled, id [][32]byte) (event.Subscription, error)

	ParseChainlinkCancelled(log types.Log) (*ConsumerChainlinkCancelled, error)

	FilterChainlinkFulfilled(opts *bind.FilterOpts, id [][32]byte) (*ConsumerChainlinkFulfilledIterator, error)

	WatchChainlinkFulfilled(opts *bind.WatchOpts, sink chan<- *ConsumerChainlinkFulfilled, id [][32]byte) (event.Subscription, error)

	ParseChainlinkFulfilled(log types.Log) (*ConsumerChainlinkFulfilled, error)

	FilterChainlinkRequested(opts *bind.FilterOpts, id [][32]byte) (*ConsumerChainlinkRequestedIterator, error)

	WatchChainlinkRequested(opts *bind.WatchOpts, sink chan<- *ConsumerChainlinkRequested, id [][32]byte) (event.Subscription, error)

	ParseChainlinkRequested(log types.Log) (*ConsumerChainlinkRequested, error)

	FilterRequestFulfilled(opts *bind.FilterOpts, requestId [][32]byte, price [][32]byte) (*ConsumerRequestFulfilledIterator, error)

	WatchRequestFulfilled(opts *bind.WatchOpts, sink chan<- *ConsumerRequestFulfilled, requestId [][32]byte, price [][32]byte) (event.Subscription, error)

	ParseRequestFulfilled(log types.Log) (*ConsumerRequestFulfilled, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
