// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package test_api_consumer_wrapper

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

var TestAPIConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ChainlinkCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ChainlinkFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ChainlinkRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"roundID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"PerfMetricsEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_payment\",\"type\":\"uint256\"},{\"internalType\":\"bytes4\",\"name\":\"_callbackFunctionId\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"_expiration\",\"type\":\"uint256\"}],\"name\":\"cancelRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_jobId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_payment\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_url\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_path\",\"type\":\"string\"},{\"internalType\":\"int256\",\"name\":\"_times\",\"type\":\"int256\"}],\"name\":\"createRequestTo\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentRoundID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"data\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_data\",\"type\":\"uint256\"}],\"name\":\"fulfill\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChainlinkToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"selector\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052600160045560006007553480156200001b57600080fd5b50604051620017a1380380620017a1833981810160405260208110156200004157600080fd5b5051600680546001600160a01b0319163317908190556040516001600160a01b0391909116906000907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a36001600160a01b038116620000b757620000b16001600160e01b03620000d216565b620000cb565b620000cb816001600160e01b036200016316565b5062000185565b6200016173c89bd4e1632d3a43cb03aaad5262cbe4038bc5716001600160a01b03166338cc48316040518163ffffffff1660e01b815260040160206040518083038186803b1580156200012457600080fd5b505afa15801562000139573d6000803e3d6000fd5b505050506040513d60208110156200015057600080fd5b50516001600160e01b036200016316565b565b600280546001600160a01b0319166001600160a01b0392909216919091179055565b61160c80620001956000396000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c80638dc654a211610081578063ea3d508a1161005b578063ea3d508a146102ca578063ec65d0f814610307578063f2fde38b14610358576100c9565b80638dc654a21461029e5780638f32d59b146102a6578063a312c4f2146102c2576100c9565b80634357855e116100b25780634357855e1461026957806373d4a13a1461028e5780638da5cb5b14610296576100c9565b8063165d35e1146100ce57806316ef7f1a146100ff575b600080fd5b6100d661038b565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b610257600480360360c081101561011557600080fd5b73ffffffffffffffffffffffffffffffffffffffff823516916020810135916040820135919081019060808101606082013564010000000081111561015957600080fd5b82018360208201111561016b57600080fd5b8035906020019184600183028401116401000000008311171561018d57600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092959493602081019350359150506401000000008111156101e057600080fd5b8201836020820111156101f257600080fd5b8035906020019184600183028401116401000000008311171561021457600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550509135925061039a915050565b60408051918252519081900360200190f35b61028c6004803603604081101561027f57600080fd5b508035906020013561055c565b005b6102576105ae565b6100d66105b4565b61028c6105d0565b6102ae610800565b604080519115158252519081900360200190f35b61025761081e565b6102d2610824565b604080517fffffffff000000000000000000000000000000000000000000000000000000009092168252519081900360200190f35b61028c6004803603608081101561031d57600080fd5b508035906020810135907fffffffff00000000000000000000000000000000000000000000000000000000604082013516906060013561082d565b61028c6004803603602081101561036e57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166108b2565b600061039561092e565b905090565b60006103a4610800565b61040f57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b600980547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000016634357855e179055610445611567565b61047087307f4357855e0000000000000000000000000000000000000000000000000000000061094a565b60408051808201909152600381527f676574000000000000000000000000000000000000000000000000000000000060208201529091506104b99082908763ffffffff61097516565b60408051808201909152600481527f706174680000000000000000000000000000000000000000000000000000000060208201526104ff9082908663ffffffff61097516565b60408051808201909152600581527f74696d657300000000000000000000000000000000000000000000000000000060208201526105459082908563ffffffff6109a416565b6105508882886109ce565b98975050505050505050565b6008819055600780546001019081905560408051918252602082018490524282820152517ffbaf68ee7b9032982942607eaea1859969ed8674797b5c2fc6fecaa7538519469181900360600190a15050565b60085481565b60065473ffffffffffffffffffffffffffffffffffffffff1690565b6105d8610800565b61064357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b600061064d61092e565b604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905191925073ffffffffffffffffffffffffffffffffffffffff83169163a9059cbb91339184916370a08231916024808301926020929190829003018186803b1580156106c657600080fd5b505afa1580156106da573d6000803e3d6000fd5b505050506040513d60208110156106f057600080fd5b5051604080517fffffffff0000000000000000000000000000000000000000000000000000000060e086901b16815273ffffffffffffffffffffffffffffffffffffffff909316600484015260248301919091525160448083019260209291908290030181600087803b15801561076657600080fd5b505af115801561077a573d6000803e3d6000fd5b505050506040513d602081101561079057600080fd5b50516107fd57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f556e61626c6520746f207472616e736665720000000000000000000000000000604482015290519081900360640190fd5b50565b60065473ffffffffffffffffffffffffffffffffffffffff16331490565b60075481565b60095460e01b81565b610835610800565b6108a057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b6108ac84848484610c0b565b50505050565b6108ba610800565b61092557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b6107fd81610d46565b60025473ffffffffffffffffffffffffffffffffffffffff1690565b610952611567565b61095a611567565b61096c8186868663ffffffff610e4016565b95945050505050565b608083015161098a908363ffffffff610ea216565b608083015161099f908263ffffffff610ea216565b505050565b60808301516109b9908363ffffffff610ea216565b608083015161099f908263ffffffff610ebf16565b6004546040805130606090811b60208084019190915260348084018690528451808503909101815260549093018452825192810192909220908601939093526000838152600590915281812080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8816179055905182917fb5e6e01e79f91267dc17b4e6314d5d4d03593d2ceee0fbb452b750bd70ea5af991a260025473ffffffffffffffffffffffffffffffffffffffff16634000aea08584610aa887610f35565b6040518463ffffffff1660e01b8152600401808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b83811015610b2c578181015183820152602001610b14565b50505050905090810190601f168015610b595780820380516001836020036101000a031916815260200191505b50945050505050602060405180830381600087803b158015610b7a57600080fd5b505af1158015610b8e573d6000803e3d6000fd5b505050506040513d6020811015610ba457600080fd5b5051610bfb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260238152602001806115dd6023913960400191505060405180910390fd5b6004805460010190559392505050565b60008481526005602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000008116909155905173ffffffffffffffffffffffffffffffffffffffff9091169186917fe1fe3afa0f7f761ff0a8b89086790efd5140d2907ebd5b7ff6bfcb5e075fd4c59190a2604080517f6ee4d55300000000000000000000000000000000000000000000000000000000815260048101879052602481018690527fffffffff000000000000000000000000000000000000000000000000000000008516604482015260648101849052905173ffffffffffffffffffffffffffffffffffffffff831691636ee4d55391608480830192600092919082900301818387803b158015610d2757600080fd5b505af1158015610d3b573d6000803e3d6000fd5b505050505050505050565b73ffffffffffffffffffffffffffffffffffffffff8116610db2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260268152602001806115b76026913960400191505060405180910390fd5b60065460405173ffffffffffffffffffffffffffffffffffffffff8084169216907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a3600680547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b610e48611567565b610e58856080015161010061111e565b505091835273ffffffffffffffffffffffffffffffffffffffff1660208301527fffffffff0000000000000000000000000000000000000000000000000000000016604082015290565b610eaf826003835161115e565b61099f828263ffffffff6112a916565b7fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000811215610ef657610ef182826112ca565b610f31565b67ffffffffffffffff811315610f1057610ef18282611327565b60008112610f2457610ef18260008361115e565b610f31826001831961115e565b5050565b6060634042994660e01b60008084600001518560200151866040015187606001516001896080015160000151604051602401808973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018881526020018781526020018673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001857bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200184815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b83811015611061578181015183820152602001611049565b50505050905090810190601f16801561108e5780820380516001836020036101000a031916815260200191505b50604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909d169c909c17909b5250989950505050505050505050919050565b61112661159c565b602082061561113b5760208206602003820191505b506020808301829052604080518085526000815283019091019052815b92915050565b60178167ffffffffffffffff161161118f576111898360e0600585901b16831763ffffffff61136216565b5061099f565b60ff8167ffffffffffffffff16116111d9576111bc836018611fe0600586901b161763ffffffff61136216565b506111898367ffffffffffffffff8316600163ffffffff61137a16565b61ffff8167ffffffffffffffff161161122457611207836019611fe0600586901b161763ffffffff61136216565b506111898367ffffffffffffffff8316600263ffffffff61137a16565b63ffffffff8167ffffffffffffffff16116112715761125483601a611fe0600586901b161763ffffffff61136216565b506111898367ffffffffffffffff8316600463ffffffff61137a16565b61128c83601b611fe0600586901b161763ffffffff61136216565b506108ac8367ffffffffffffffff8316600863ffffffff61137a16565b6112b161159c565b6112c38384600001515184855161139b565b9392505050565b6112db8260c363ffffffff61136216565b50610f3182827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0360405160200180828152602001915050604051602081830303815290604052611483565b6113388260c263ffffffff61136216565b50610f31828260405160200180828152602001915050604051602081830303815290604052611483565b61136a61159c565b6112c38384600001515184611490565b61138261159c565b6113938485600001515185856114db565b949350505050565b6113a361159c565b82518211156113b157600080fd5b846020015182850111156113db576113db856113d38760200151878601611539565b600202611550565b6000808651805187602083010193508088870111156113fa5787860182525b505050602084015b6020841061143f57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09093019260209182019101611402565b5181517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60208690036101000a019081169019919091161790525083949350505050565b610eaf826002835161115e565b61149861159c565b836020015183106114b4576114b4848560200151600202611550565b8351805160208583010184815350808514156114d1576001810182525b5093949350505050565b6114e361159c565b846020015184830111156115005761150085858401600202611550565b60006001836101000a03905085518386820101858319825116178152508051848701111561152e5783860181525b509495945050505050565b60008183111561154a575081611158565b50919050565b815161155c838361111e565b506108ac83826112a9565b6040805160a08101825260008082526020820181905291810182905260608101919091526080810161159761159c565b905290565b60405180604001604052806060815260200160008152509056fe4f776e61626c653a206e6577206f776e657220697320746865207a65726f2061646472657373756e61626c6520746f207472616e73666572416e6443616c6c20746f206f7261636c65a164736f6c6343000606000a",
}

var TestAPIConsumerABI = TestAPIConsumerMetaData.ABI

var TestAPIConsumerBin = TestAPIConsumerMetaData.Bin

func DeployTestAPIConsumer(auth *bind.TransactOpts, backend bind.ContractBackend, _link common.Address) (common.Address, *types.Transaction, *TestAPIConsumer, error) {
	parsed, err := TestAPIConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestAPIConsumerBin), backend, _link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestAPIConsumer{address: address, abi: *parsed, TestAPIConsumerCaller: TestAPIConsumerCaller{contract: contract}, TestAPIConsumerTransactor: TestAPIConsumerTransactor{contract: contract}, TestAPIConsumerFilterer: TestAPIConsumerFilterer{contract: contract}}, nil
}

type TestAPIConsumer struct {
	address common.Address
	abi     abi.ABI
	TestAPIConsumerCaller
	TestAPIConsumerTransactor
	TestAPIConsumerFilterer
}

type TestAPIConsumerCaller struct {
	contract *bind.BoundContract
}

type TestAPIConsumerTransactor struct {
	contract *bind.BoundContract
}

type TestAPIConsumerFilterer struct {
	contract *bind.BoundContract
}

type TestAPIConsumerSession struct {
	Contract     *TestAPIConsumer
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type TestAPIConsumerCallerSession struct {
	Contract *TestAPIConsumerCaller
	CallOpts bind.CallOpts
}

type TestAPIConsumerTransactorSession struct {
	Contract     *TestAPIConsumerTransactor
	TransactOpts bind.TransactOpts
}

type TestAPIConsumerRaw struct {
	Contract *TestAPIConsumer
}

type TestAPIConsumerCallerRaw struct {
	Contract *TestAPIConsumerCaller
}

type TestAPIConsumerTransactorRaw struct {
	Contract *TestAPIConsumerTransactor
}

func NewTestAPIConsumer(address common.Address, backend bind.ContractBackend) (*TestAPIConsumer, error) {
	abi, err := abi.JSON(strings.NewReader(TestAPIConsumerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindTestAPIConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestAPIConsumer{address: address, abi: abi, TestAPIConsumerCaller: TestAPIConsumerCaller{contract: contract}, TestAPIConsumerTransactor: TestAPIConsumerTransactor{contract: contract}, TestAPIConsumerFilterer: TestAPIConsumerFilterer{contract: contract}}, nil
}

func NewTestAPIConsumerCaller(address common.Address, caller bind.ContractCaller) (*TestAPIConsumerCaller, error) {
	contract, err := bindTestAPIConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestAPIConsumerCaller{contract: contract}, nil
}

func NewTestAPIConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*TestAPIConsumerTransactor, error) {
	contract, err := bindTestAPIConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestAPIConsumerTransactor{contract: contract}, nil
}

func NewTestAPIConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*TestAPIConsumerFilterer, error) {
	contract, err := bindTestAPIConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestAPIConsumerFilterer{contract: contract}, nil
}

func bindTestAPIConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestAPIConsumerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_TestAPIConsumer *TestAPIConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestAPIConsumer.Contract.TestAPIConsumerCaller.contract.Call(opts, result, method, params...)
}

func (_TestAPIConsumer *TestAPIConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestAPIConsumer.Contract.TestAPIConsumerTransactor.contract.Transfer(opts)
}

func (_TestAPIConsumer *TestAPIConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestAPIConsumer.Contract.TestAPIConsumerTransactor.contract.Transact(opts, method, params...)
}

func (_TestAPIConsumer *TestAPIConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestAPIConsumer.Contract.contract.Call(opts, result, method, params...)
}

func (_TestAPIConsumer *TestAPIConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestAPIConsumer.Contract.contract.Transfer(opts)
}

func (_TestAPIConsumer *TestAPIConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestAPIConsumer.Contract.contract.Transact(opts, method, params...)
}

func (_TestAPIConsumer *TestAPIConsumerCaller) CurrentRoundID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestAPIConsumer.contract.Call(opts, &out, "currentRoundID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_TestAPIConsumer *TestAPIConsumerSession) CurrentRoundID() (*big.Int, error) {
	return _TestAPIConsumer.Contract.CurrentRoundID(&_TestAPIConsumer.CallOpts)
}

func (_TestAPIConsumer *TestAPIConsumerCallerSession) CurrentRoundID() (*big.Int, error) {
	return _TestAPIConsumer.Contract.CurrentRoundID(&_TestAPIConsumer.CallOpts)
}

func (_TestAPIConsumer *TestAPIConsumerCaller) Data(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestAPIConsumer.contract.Call(opts, &out, "data")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_TestAPIConsumer *TestAPIConsumerSession) Data() (*big.Int, error) {
	return _TestAPIConsumer.Contract.Data(&_TestAPIConsumer.CallOpts)
}

func (_TestAPIConsumer *TestAPIConsumerCallerSession) Data() (*big.Int, error) {
	return _TestAPIConsumer.Contract.Data(&_TestAPIConsumer.CallOpts)
}

func (_TestAPIConsumer *TestAPIConsumerCaller) GetChainlinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestAPIConsumer.contract.Call(opts, &out, "getChainlinkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_TestAPIConsumer *TestAPIConsumerSession) GetChainlinkToken() (common.Address, error) {
	return _TestAPIConsumer.Contract.GetChainlinkToken(&_TestAPIConsumer.CallOpts)
}

func (_TestAPIConsumer *TestAPIConsumerCallerSession) GetChainlinkToken() (common.Address, error) {
	return _TestAPIConsumer.Contract.GetChainlinkToken(&_TestAPIConsumer.CallOpts)
}

func (_TestAPIConsumer *TestAPIConsumerCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TestAPIConsumer.contract.Call(opts, &out, "isOwner")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_TestAPIConsumer *TestAPIConsumerSession) IsOwner() (bool, error) {
	return _TestAPIConsumer.Contract.IsOwner(&_TestAPIConsumer.CallOpts)
}

func (_TestAPIConsumer *TestAPIConsumerCallerSession) IsOwner() (bool, error) {
	return _TestAPIConsumer.Contract.IsOwner(&_TestAPIConsumer.CallOpts)
}

func (_TestAPIConsumer *TestAPIConsumerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestAPIConsumer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_TestAPIConsumer *TestAPIConsumerSession) Owner() (common.Address, error) {
	return _TestAPIConsumer.Contract.Owner(&_TestAPIConsumer.CallOpts)
}

func (_TestAPIConsumer *TestAPIConsumerCallerSession) Owner() (common.Address, error) {
	return _TestAPIConsumer.Contract.Owner(&_TestAPIConsumer.CallOpts)
}

func (_TestAPIConsumer *TestAPIConsumerCaller) Selector(opts *bind.CallOpts) ([4]byte, error) {
	var out []interface{}
	err := _TestAPIConsumer.contract.Call(opts, &out, "selector")

	if err != nil {
		return *new([4]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([4]byte)).(*[4]byte)

	return out0, err

}

func (_TestAPIConsumer *TestAPIConsumerSession) Selector() ([4]byte, error) {
	return _TestAPIConsumer.Contract.Selector(&_TestAPIConsumer.CallOpts)
}

func (_TestAPIConsumer *TestAPIConsumerCallerSession) Selector() ([4]byte, error) {
	return _TestAPIConsumer.Contract.Selector(&_TestAPIConsumer.CallOpts)
}

func (_TestAPIConsumer *TestAPIConsumerTransactor) CancelRequest(opts *bind.TransactOpts, _requestId [32]byte, _payment *big.Int, _callbackFunctionId [4]byte, _expiration *big.Int) (*types.Transaction, error) {
	return _TestAPIConsumer.contract.Transact(opts, "cancelRequest", _requestId, _payment, _callbackFunctionId, _expiration)
}

func (_TestAPIConsumer *TestAPIConsumerSession) CancelRequest(_requestId [32]byte, _payment *big.Int, _callbackFunctionId [4]byte, _expiration *big.Int) (*types.Transaction, error) {
	return _TestAPIConsumer.Contract.CancelRequest(&_TestAPIConsumer.TransactOpts, _requestId, _payment, _callbackFunctionId, _expiration)
}

func (_TestAPIConsumer *TestAPIConsumerTransactorSession) CancelRequest(_requestId [32]byte, _payment *big.Int, _callbackFunctionId [4]byte, _expiration *big.Int) (*types.Transaction, error) {
	return _TestAPIConsumer.Contract.CancelRequest(&_TestAPIConsumer.TransactOpts, _requestId, _payment, _callbackFunctionId, _expiration)
}

func (_TestAPIConsumer *TestAPIConsumerTransactor) CreateRequestTo(opts *bind.TransactOpts, _oracle common.Address, _jobId [32]byte, _payment *big.Int, _url string, _path string, _times *big.Int) (*types.Transaction, error) {
	return _TestAPIConsumer.contract.Transact(opts, "createRequestTo", _oracle, _jobId, _payment, _url, _path, _times)
}

func (_TestAPIConsumer *TestAPIConsumerSession) CreateRequestTo(_oracle common.Address, _jobId [32]byte, _payment *big.Int, _url string, _path string, _times *big.Int) (*types.Transaction, error) {
	return _TestAPIConsumer.Contract.CreateRequestTo(&_TestAPIConsumer.TransactOpts, _oracle, _jobId, _payment, _url, _path, _times)
}

func (_TestAPIConsumer *TestAPIConsumerTransactorSession) CreateRequestTo(_oracle common.Address, _jobId [32]byte, _payment *big.Int, _url string, _path string, _times *big.Int) (*types.Transaction, error) {
	return _TestAPIConsumer.Contract.CreateRequestTo(&_TestAPIConsumer.TransactOpts, _oracle, _jobId, _payment, _url, _path, _times)
}

func (_TestAPIConsumer *TestAPIConsumerTransactor) Fulfill(opts *bind.TransactOpts, _requestId [32]byte, _data *big.Int) (*types.Transaction, error) {
	return _TestAPIConsumer.contract.Transact(opts, "fulfill", _requestId, _data)
}

func (_TestAPIConsumer *TestAPIConsumerSession) Fulfill(_requestId [32]byte, _data *big.Int) (*types.Transaction, error) {
	return _TestAPIConsumer.Contract.Fulfill(&_TestAPIConsumer.TransactOpts, _requestId, _data)
}

func (_TestAPIConsumer *TestAPIConsumerTransactorSession) Fulfill(_requestId [32]byte, _data *big.Int) (*types.Transaction, error) {
	return _TestAPIConsumer.Contract.Fulfill(&_TestAPIConsumer.TransactOpts, _requestId, _data)
}

func (_TestAPIConsumer *TestAPIConsumerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TestAPIConsumer.contract.Transact(opts, "transferOwnership", newOwner)
}

func (_TestAPIConsumer *TestAPIConsumerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TestAPIConsumer.Contract.TransferOwnership(&_TestAPIConsumer.TransactOpts, newOwner)
}

func (_TestAPIConsumer *TestAPIConsumerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TestAPIConsumer.Contract.TransferOwnership(&_TestAPIConsumer.TransactOpts, newOwner)
}

func (_TestAPIConsumer *TestAPIConsumerTransactor) WithdrawLink(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestAPIConsumer.contract.Transact(opts, "withdrawLink")
}

func (_TestAPIConsumer *TestAPIConsumerSession) WithdrawLink() (*types.Transaction, error) {
	return _TestAPIConsumer.Contract.WithdrawLink(&_TestAPIConsumer.TransactOpts)
}

func (_TestAPIConsumer *TestAPIConsumerTransactorSession) WithdrawLink() (*types.Transaction, error) {
	return _TestAPIConsumer.Contract.WithdrawLink(&_TestAPIConsumer.TransactOpts)
}

type TestAPIConsumerChainlinkCancelledIterator struct {
	Event *TestAPIConsumerChainlinkCancelled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TestAPIConsumerChainlinkCancelledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestAPIConsumerChainlinkCancelled)
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
		it.Event = new(TestAPIConsumerChainlinkCancelled)
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

func (it *TestAPIConsumerChainlinkCancelledIterator) Error() error {
	return it.fail
}

func (it *TestAPIConsumerChainlinkCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TestAPIConsumerChainlinkCancelled struct {
	Id  [32]byte
	Raw types.Log
}

func (_TestAPIConsumer *TestAPIConsumerFilterer) FilterChainlinkCancelled(opts *bind.FilterOpts, id [][32]byte) (*TestAPIConsumerChainlinkCancelledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _TestAPIConsumer.contract.FilterLogs(opts, "ChainlinkCancelled", idRule)
	if err != nil {
		return nil, err
	}
	return &TestAPIConsumerChainlinkCancelledIterator{contract: _TestAPIConsumer.contract, event: "ChainlinkCancelled", logs: logs, sub: sub}, nil
}

func (_TestAPIConsumer *TestAPIConsumerFilterer) WatchChainlinkCancelled(opts *bind.WatchOpts, sink chan<- *TestAPIConsumerChainlinkCancelled, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _TestAPIConsumer.contract.WatchLogs(opts, "ChainlinkCancelled", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TestAPIConsumerChainlinkCancelled)
				if err := _TestAPIConsumer.contract.UnpackLog(event, "ChainlinkCancelled", log); err != nil {
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

func (_TestAPIConsumer *TestAPIConsumerFilterer) ParseChainlinkCancelled(log types.Log) (*TestAPIConsumerChainlinkCancelled, error) {
	event := new(TestAPIConsumerChainlinkCancelled)
	if err := _TestAPIConsumer.contract.UnpackLog(event, "ChainlinkCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TestAPIConsumerChainlinkFulfilledIterator struct {
	Event *TestAPIConsumerChainlinkFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TestAPIConsumerChainlinkFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestAPIConsumerChainlinkFulfilled)
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
		it.Event = new(TestAPIConsumerChainlinkFulfilled)
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

func (it *TestAPIConsumerChainlinkFulfilledIterator) Error() error {
	return it.fail
}

func (it *TestAPIConsumerChainlinkFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TestAPIConsumerChainlinkFulfilled struct {
	Id  [32]byte
	Raw types.Log
}

func (_TestAPIConsumer *TestAPIConsumerFilterer) FilterChainlinkFulfilled(opts *bind.FilterOpts, id [][32]byte) (*TestAPIConsumerChainlinkFulfilledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _TestAPIConsumer.contract.FilterLogs(opts, "ChainlinkFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return &TestAPIConsumerChainlinkFulfilledIterator{contract: _TestAPIConsumer.contract, event: "ChainlinkFulfilled", logs: logs, sub: sub}, nil
}

func (_TestAPIConsumer *TestAPIConsumerFilterer) WatchChainlinkFulfilled(opts *bind.WatchOpts, sink chan<- *TestAPIConsumerChainlinkFulfilled, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _TestAPIConsumer.contract.WatchLogs(opts, "ChainlinkFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TestAPIConsumerChainlinkFulfilled)
				if err := _TestAPIConsumer.contract.UnpackLog(event, "ChainlinkFulfilled", log); err != nil {
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

func (_TestAPIConsumer *TestAPIConsumerFilterer) ParseChainlinkFulfilled(log types.Log) (*TestAPIConsumerChainlinkFulfilled, error) {
	event := new(TestAPIConsumerChainlinkFulfilled)
	if err := _TestAPIConsumer.contract.UnpackLog(event, "ChainlinkFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TestAPIConsumerChainlinkRequestedIterator struct {
	Event *TestAPIConsumerChainlinkRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TestAPIConsumerChainlinkRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestAPIConsumerChainlinkRequested)
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
		it.Event = new(TestAPIConsumerChainlinkRequested)
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

func (it *TestAPIConsumerChainlinkRequestedIterator) Error() error {
	return it.fail
}

func (it *TestAPIConsumerChainlinkRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TestAPIConsumerChainlinkRequested struct {
	Id  [32]byte
	Raw types.Log
}

func (_TestAPIConsumer *TestAPIConsumerFilterer) FilterChainlinkRequested(opts *bind.FilterOpts, id [][32]byte) (*TestAPIConsumerChainlinkRequestedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _TestAPIConsumer.contract.FilterLogs(opts, "ChainlinkRequested", idRule)
	if err != nil {
		return nil, err
	}
	return &TestAPIConsumerChainlinkRequestedIterator{contract: _TestAPIConsumer.contract, event: "ChainlinkRequested", logs: logs, sub: sub}, nil
}

func (_TestAPIConsumer *TestAPIConsumerFilterer) WatchChainlinkRequested(opts *bind.WatchOpts, sink chan<- *TestAPIConsumerChainlinkRequested, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _TestAPIConsumer.contract.WatchLogs(opts, "ChainlinkRequested", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TestAPIConsumerChainlinkRequested)
				if err := _TestAPIConsumer.contract.UnpackLog(event, "ChainlinkRequested", log); err != nil {
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

func (_TestAPIConsumer *TestAPIConsumerFilterer) ParseChainlinkRequested(log types.Log) (*TestAPIConsumerChainlinkRequested, error) {
	event := new(TestAPIConsumerChainlinkRequested)
	if err := _TestAPIConsumer.contract.UnpackLog(event, "ChainlinkRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TestAPIConsumerOwnershipTransferredIterator struct {
	Event *TestAPIConsumerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TestAPIConsumerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestAPIConsumerOwnershipTransferred)
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
		it.Event = new(TestAPIConsumerOwnershipTransferred)
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

func (it *TestAPIConsumerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *TestAPIConsumerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TestAPIConsumerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log
}

func (_TestAPIConsumer *TestAPIConsumerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TestAPIConsumerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TestAPIConsumer.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TestAPIConsumerOwnershipTransferredIterator{contract: _TestAPIConsumer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_TestAPIConsumer *TestAPIConsumerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TestAPIConsumerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TestAPIConsumer.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TestAPIConsumerOwnershipTransferred)
				if err := _TestAPIConsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_TestAPIConsumer *TestAPIConsumerFilterer) ParseOwnershipTransferred(log types.Log) (*TestAPIConsumerOwnershipTransferred, error) {
	event := new(TestAPIConsumerOwnershipTransferred)
	if err := _TestAPIConsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TestAPIConsumerPerfMetricsEventIterator struct {
	Event *TestAPIConsumerPerfMetricsEvent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TestAPIConsumerPerfMetricsEventIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestAPIConsumerPerfMetricsEvent)
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
		it.Event = new(TestAPIConsumerPerfMetricsEvent)
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

func (it *TestAPIConsumerPerfMetricsEventIterator) Error() error {
	return it.fail
}

func (it *TestAPIConsumerPerfMetricsEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TestAPIConsumerPerfMetricsEvent struct {
	RoundID   *big.Int
	RequestId [32]byte
	Timestamp *big.Int
	Raw       types.Log
}

func (_TestAPIConsumer *TestAPIConsumerFilterer) FilterPerfMetricsEvent(opts *bind.FilterOpts) (*TestAPIConsumerPerfMetricsEventIterator, error) {

	logs, sub, err := _TestAPIConsumer.contract.FilterLogs(opts, "PerfMetricsEvent")
	if err != nil {
		return nil, err
	}
	return &TestAPIConsumerPerfMetricsEventIterator{contract: _TestAPIConsumer.contract, event: "PerfMetricsEvent", logs: logs, sub: sub}, nil
}

func (_TestAPIConsumer *TestAPIConsumerFilterer) WatchPerfMetricsEvent(opts *bind.WatchOpts, sink chan<- *TestAPIConsumerPerfMetricsEvent) (event.Subscription, error) {

	logs, sub, err := _TestAPIConsumer.contract.WatchLogs(opts, "PerfMetricsEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TestAPIConsumerPerfMetricsEvent)
				if err := _TestAPIConsumer.contract.UnpackLog(event, "PerfMetricsEvent", log); err != nil {
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

func (_TestAPIConsumer *TestAPIConsumerFilterer) ParsePerfMetricsEvent(log types.Log) (*TestAPIConsumerPerfMetricsEvent, error) {
	event := new(TestAPIConsumerPerfMetricsEvent)
	if err := _TestAPIConsumer.contract.UnpackLog(event, "PerfMetricsEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_TestAPIConsumer *TestAPIConsumer) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _TestAPIConsumer.abi.Events["ChainlinkCancelled"].ID:
		return _TestAPIConsumer.ParseChainlinkCancelled(log)
	case _TestAPIConsumer.abi.Events["ChainlinkFulfilled"].ID:
		return _TestAPIConsumer.ParseChainlinkFulfilled(log)
	case _TestAPIConsumer.abi.Events["ChainlinkRequested"].ID:
		return _TestAPIConsumer.ParseChainlinkRequested(log)
	case _TestAPIConsumer.abi.Events["OwnershipTransferred"].ID:
		return _TestAPIConsumer.ParseOwnershipTransferred(log)
	case _TestAPIConsumer.abi.Events["PerfMetricsEvent"].ID:
		return _TestAPIConsumer.ParsePerfMetricsEvent(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (TestAPIConsumerChainlinkCancelled) Topic() common.Hash {
	return common.HexToHash("0xe1fe3afa0f7f761ff0a8b89086790efd5140d2907ebd5b7ff6bfcb5e075fd4c5")
}

func (TestAPIConsumerChainlinkFulfilled) Topic() common.Hash {
	return common.HexToHash("0x7cc135e0cebb02c3480ae5d74d377283180a2601f8f644edf7987b009316c63a")
}

func (TestAPIConsumerChainlinkRequested) Topic() common.Hash {
	return common.HexToHash("0xb5e6e01e79f91267dc17b4e6314d5d4d03593d2ceee0fbb452b750bd70ea5af9")
}

func (TestAPIConsumerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (TestAPIConsumerPerfMetricsEvent) Topic() common.Hash {
	return common.HexToHash("0xfbaf68ee7b9032982942607eaea1859969ed8674797b5c2fc6fecaa753851946")
}

func (_TestAPIConsumer *TestAPIConsumer) Address() common.Address {
	return _TestAPIConsumer.address
}

type TestAPIConsumerInterface interface {
	CurrentRoundID(opts *bind.CallOpts) (*big.Int, error)

	Data(opts *bind.CallOpts) (*big.Int, error)

	GetChainlinkToken(opts *bind.CallOpts) (common.Address, error)

	IsOwner(opts *bind.CallOpts) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	Selector(opts *bind.CallOpts) ([4]byte, error)

	CancelRequest(opts *bind.TransactOpts, _requestId [32]byte, _payment *big.Int, _callbackFunctionId [4]byte, _expiration *big.Int) (*types.Transaction, error)

	CreateRequestTo(opts *bind.TransactOpts, _oracle common.Address, _jobId [32]byte, _payment *big.Int, _url string, _path string, _times *big.Int) (*types.Transaction, error)

	Fulfill(opts *bind.TransactOpts, _requestId [32]byte, _data *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error)

	WithdrawLink(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterChainlinkCancelled(opts *bind.FilterOpts, id [][32]byte) (*TestAPIConsumerChainlinkCancelledIterator, error)

	WatchChainlinkCancelled(opts *bind.WatchOpts, sink chan<- *TestAPIConsumerChainlinkCancelled, id [][32]byte) (event.Subscription, error)

	ParseChainlinkCancelled(log types.Log) (*TestAPIConsumerChainlinkCancelled, error)

	FilterChainlinkFulfilled(opts *bind.FilterOpts, id [][32]byte) (*TestAPIConsumerChainlinkFulfilledIterator, error)

	WatchChainlinkFulfilled(opts *bind.WatchOpts, sink chan<- *TestAPIConsumerChainlinkFulfilled, id [][32]byte) (event.Subscription, error)

	ParseChainlinkFulfilled(log types.Log) (*TestAPIConsumerChainlinkFulfilled, error)

	FilterChainlinkRequested(opts *bind.FilterOpts, id [][32]byte) (*TestAPIConsumerChainlinkRequestedIterator, error)

	WatchChainlinkRequested(opts *bind.WatchOpts, sink chan<- *TestAPIConsumerChainlinkRequested, id [][32]byte) (event.Subscription, error)

	ParseChainlinkRequested(log types.Log) (*TestAPIConsumerChainlinkRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TestAPIConsumerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TestAPIConsumerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*TestAPIConsumerOwnershipTransferred, error)

	FilterPerfMetricsEvent(opts *bind.FilterOpts) (*TestAPIConsumerPerfMetricsEventIterator, error)

	WatchPerfMetricsEvent(opts *bind.WatchOpts, sink chan<- *TestAPIConsumerPerfMetricsEvent) (event.Subscription, error)

	ParsePerfMetricsEvent(log types.Log) (*TestAPIConsumerPerfMetricsEvent, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
