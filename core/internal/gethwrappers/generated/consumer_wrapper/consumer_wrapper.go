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
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_specId\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ChainlinkCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ChainlinkFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ChainlinkRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"price\",\"type\":\"bytes32\"}],\"name\":\"RequestFulfilled\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_requestId\",\"type\":\"bytes32\"}],\"name\":\"addExternalRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_payment\",\"type\":\"uint256\"},{\"internalType\":\"bytes4\",\"name\":\"_callbackFunctionId\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"_expiration\",\"type\":\"uint256\"}],\"name\":\"cancelRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentPrice\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentPriceInt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_price\",\"type\":\"bytes32\"}],\"name\":\"fulfill\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"}],\"name\":\"fulfillParametersWithCustomURLs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_currency\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_payment\",\"type\":\"uint256\"}],\"name\":\"requestEthereumPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_currency\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_payment\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_callback\",\"type\":\"address\"}],\"name\":\"requestEthereumPriceByCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_urlUSD\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_pathUSD\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_payment\",\"type\":\"uint256\"}],\"name\":\"requestMultipleParametersWithCustomURLs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_specId\",\"type\":\"bytes32\"}],\"name\":\"setSpecID\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052600160045534801561001557600080fd5b506040516116623803806116628339818101604052606081101561003857600080fd5b508051602082015160409092015190919061005283610066565b61005b82610088565b600655506100aa9050565b600280546001600160a01b0319166001600160a01b0392909216919091179055565b600380546001600160a01b0319166001600160a01b0392909216919091179055565b6115a9806100b96000396000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c806374961d4d116100815780639438a6011161005b5780639438a601146104405780639d1b464a1461045a578063e8d5359d14610462576100c9565b806374961d4d146102cf57806383db5cbc146103905780638dc654a214610438576100c9565b80635591a608116100b25780635591a608146101105780635b8260051461017d57806371c2002a146101a0576100c9565b8063042f2b65146100ce578063501fdd5d146100f3575b600080fd5b6100f1600480360360408110156100e457600080fd5b508035906020013561049b565b005b6100f16004803603602081101561010957600080fd5b50356105a8565b6100f1600480360360a081101561012657600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813516906020810135906040810135907fffffffff0000000000000000000000000000000000000000000000000000000060608201351690608001356105ad565b6100f16004803603604081101561019357600080fd5b5080359060200135610674565b6100f1600480360360608110156101b657600080fd5b8101906020810181356401000000008111156101d157600080fd5b8201836020820111156101e357600080fd5b8035906020019184600183028401116401000000008311171561020557600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929594936020810193503591505064010000000081111561025857600080fd5b82018360208201111561026a57600080fd5b8035906020019184600183028401116401000000008311171561028c57600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295505091359250610781915050565b6100f1600480360360608110156102e557600080fd5b81019060208101813564010000000081111561030057600080fd5b82018360208201111561031257600080fd5b8035906020019184600183028401116401000000008311171561033457600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550508235935050506020013573ffffffffffffffffffffffffffffffffffffffff1661082b565b6100f1600480360360408110156103a657600080fd5b8101906020810181356401000000008111156103c157600080fd5b8201836020820111156103d357600080fd5b803590602001918460018302840111640100000000831117156103f557600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295505091359250610940915050565b6100f161094f565b610448610b19565b60408051918252519081900360200190f35b610448610b1f565b6100f16004803603604081101561047857600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135610b25565b600082815260056020526040902054829073ffffffffffffffffffffffffffffffffffffffff163314610519576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602881526020018061152e6028913960400191505060405180910390fd5b60008181526005602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001690555182917f7cc135e0cebb02c3480ae5d74d377283180a2601f8f644edf7987b009316c63a91a2604051829084907f0c2366233f634048c0f0458060d1228fab36d00f7c0ecf6bdf2d9c458503631190600090a35060075550565b600655565b604080517f6ee4d55300000000000000000000000000000000000000000000000000000000815260048101869052602481018590527fffffffff0000000000000000000000000000000000000000000000000000000084166044820152606481018390529051869173ffffffffffffffffffffffffffffffffffffffff831691636ee4d5539160848082019260009290919082900301818387803b15801561065457600080fd5b505af1158015610668573d6000803e3d6000fd5b50505050505050505050565b600082815260056020526040902054829073ffffffffffffffffffffffffffffffffffffffff1633146106f2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602881526020018061152e6028913960400191505060405180910390fd5b60008181526005602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001690555182917f7cc135e0cebb02c3480ae5d74d377283180a2601f8f644edf7987b009316c63a91a2604051829084907f0c2366233f634048c0f0458060d1228fab36d00f7c0ecf6bdf2d9c458503631190600090a35060085550565b600061079760065430635b82600560e01b610b2f565b60408051808201909152600681527f75726c555344000000000000000000000000000000000000000000000000000060208201529091506107da90829086610b54565b60408051808201909152600781527f7061746855534400000000000000000000000000000000000000000000000000602082015261081a90829085610b54565b6108248183610b77565b5050505050565b60006108416006548363042f2b6560e01b610b2f565b905061089d6040518060400160405280600381526020017f676574000000000000000000000000000000000000000000000000000000000081525060405180608001604052806047815260200161155660479139839190610b54565b604080516001808252818301909252600091816020015b60608152602001906001900390816108b457905050905084816000815181106108d957fe5b602002602001018190525061092e6040518060400160405280600481526020017f70617468000000000000000000000000000000000000000000000000000000008152508284610ba79092919063ffffffff16565b6109388285610b77565b505050505050565b61094b82823061082b565b5050565b6000610959610c0f565b90508073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb338373ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b1580156109df57600080fd5b505afa1580156109f3573d6000803e3d6000fd5b505050506040513d6020811015610a0957600080fd5b5051604080517fffffffff0000000000000000000000000000000000000000000000000000000060e086901b16815273ffffffffffffffffffffffffffffffffffffffff909316600484015260248301919091525160448083019260209291908290030181600087803b158015610a7f57600080fd5b505af1158015610a93573d6000803e3d6000fd5b505050506040513d6020811015610aa957600080fd5b5051610b1657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f556e61626c6520746f207472616e736665720000000000000000000000000000604482015290519081900360640190fd5b50565b60085481565b60075481565b61094b8282610c2b565b610b376114bb565b610b3f6114bb565b610b4b81868686610d12565b95945050505050565b6080830151610b639083610d74565b6080830151610b729082610d74565b505050565b600354600090610b9e9073ffffffffffffffffffffffffffffffffffffffff168484610d8b565b90505b92915050565b6080830151610bb69083610d74565b610bc38360800151610dc3565b60005b8151811015610c0157610bf9828281518110610bde57fe5b60200260200101518560800151610d7490919063ffffffff16565b600101610bc6565b50610b728360800151610dce565b60025473ffffffffffffffffffffffffffffffffffffffff1690565b600081815260056020526040902054819073ffffffffffffffffffffffffffffffffffffffff1615610cbe57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f5265717565737420697320616c72656164792070656e64696e67000000000000604482015290519081900360640190fd5b50600090815260056020526040902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b610d1a6114bb565b610d2a8560800151610100610dd9565b505091835273ffffffffffffffffffffffffffffffffffffffff1660208301527fffffffff0000000000000000000000000000000000000000000000000000000016604082015290565b610d818260038351610e13565b610b728282610eed565b6000610dbb84848460017f4042994600000000000000000000000000000000000000000000000000000000610f07565b949350505050565b610b168160046112ba565b610b168160076112ba565b610de16114f0565b6020820615610df65760208206602003820191505b506020828101829052604080518085526000815290920101905290565b60178111610e3457610e2e8360e0600585901b1683176112cb565b50610b72565b60ff8111610e5e57610e51836018611fe0600586901b16176112cb565b50610e2e838260016112e3565b61ffff8111610e8957610e7c836019611fe0600586901b16176112cb565b50610e2e838260026112e3565b63ffffffff8111610eb657610ea983601a611fe0600586901b16176112cb565b50610e2e838260046112e3565b67ffffffffffffffff8111610b7257610eda83601b611fe0600586901b16176112cb565b50610ee7838260086112e3565b50505050565b610ef56114f0565b610b9e838460000151518485516112fc565b6004546040805130606090811b60208084019190915260348084018690528451808503909101815260549093018452825192810192909220908801939093526000838152600590915281812080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8a16179055905182917fb5e6e01e79f91267dc17b4e6314d5d4d03593d2ceee0fbb452b750bd70ea5af991a2600082600080886000015189602001518a604001518b606001518a8d6080015160000151604051602401808973ffffffffffffffffffffffffffffffffffffffff1681526020018881526020018781526020018673ffffffffffffffffffffffffffffffffffffffff168152602001857bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200184815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b83811015611091578181015183820152602001611079565b50505050905090810190601f1680156110be5780820380516001836020036101000a031916815260200191505b509950505050505050505050604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050509050600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16634000aea08887846040518463ffffffff1660e01b8152600401808473ffffffffffffffffffffffffffffffffffffffff16815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b838110156111d85781810151838201526020016111c0565b50505050905090810190601f1680156112055780820380516001836020036101000a031916815260200191505b50945050505050602060405180830381600087803b15801561122657600080fd5b505af115801561123a573d6000803e3d6000fd5b505050506040513d602081101561125057600080fd5b50516112a7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602381526020018061150b6023913960400191505060405180910390fd5b5060048054600101905595945050505050565b610b7282601f611fe0600585901b16175b6112d36114f0565b610b9e83846000015151846113e4565b6112eb6114f0565b610dbb84856000015151858561142f565b6113046114f0565b825182111561131257600080fd5b8460200151828501111561133c5761133c85611334876020015187860161148d565b6002026114a4565b60008086518051876020830101935080888701111561135b5787860182525b505050602084015b602084106113a057805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09093019260209182019101611363565b5181517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60208690036101000a019081169019919091161790525083949350505050565b6113ec6114f0565b83602001518310611408576114088485602001516002026114a4565b835180516020858301018481535080851415611425576001810182525b5093949350505050565b6114376114f0565b8460200151848301111561145457611454858584016002026114a4565b60006001836101000a0390508551838682010185831982511617815250805184870111156114825783860181525b509495945050505050565b60008183111561149e575081610ba1565b50919050565b81516114b08383610dd9565b50610ee78382610eed565b6040805160a0810182526000808252602082018190529181018290526060810191909152608081016114eb6114f0565b905290565b60405180604001604052806060815260200160008152509056fe756e61626c6520746f207472616e73666572416e6443616c6c20746f206f7261636c65536f75726365206d75737420626520746865206f7261636c65206f6620746865207265717565737468747470733a2f2f6d696e2d6170692e63727970746f636f6d706172652e636f6d2f646174612f70726963653f6673796d3d455448267473796d733d5553442c4555522c4a5059a164736f6c6343000706000a",
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

func (_Consumer *ConsumerTransactor) RequestEthereumPriceByCallback(opts *bind.TransactOpts, _currency string, _payment *big.Int, _callback common.Address) (*types.Transaction, error) {
	return _Consumer.contract.Transact(opts, "requestEthereumPriceByCallback", _currency, _payment, _callback)
}

func (_Consumer *ConsumerSession) RequestEthereumPriceByCallback(_currency string, _payment *big.Int, _callback common.Address) (*types.Transaction, error) {
	return _Consumer.Contract.RequestEthereumPriceByCallback(&_Consumer.TransactOpts, _currency, _payment, _callback)
}

func (_Consumer *ConsumerTransactorSession) RequestEthereumPriceByCallback(_currency string, _payment *big.Int, _callback common.Address) (*types.Transaction, error) {
	return _Consumer.Contract.RequestEthereumPriceByCallback(&_Consumer.TransactOpts, _currency, _payment, _callback)
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

	RequestEthereumPriceByCallback(opts *bind.TransactOpts, _currency string, _payment *big.Int, _callback common.Address) (*types.Transaction, error)

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
