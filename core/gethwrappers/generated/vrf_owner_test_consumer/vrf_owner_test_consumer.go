// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_owner_test_consumer

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

var VRFV2OwnerTestConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COORDINATOR\",\"outputs\":[{\"internalType\":\"contractVRFCoordinatorV2Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"}],\"name\":\"getRequestStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"requestTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestBlockNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentBlockNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"_requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bytes32\",\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_numWords\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"_requestCount\",\"type\":\"uint16\"}],\"name\":\"requestRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_averageFulfillmentInMillions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fastestFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastRequestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requests\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"requestTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestBlockNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentBlockNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_responseCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_slowestFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"subId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a0604052600060055560006006556103e76007553480156200002157600080fd5b50604051620013cc380380620013cc8339810160408190526200004491620002bf565b6001600160601b0319606082901b166080523380600081620000ad5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000e057620000e08162000213565b5050600280546001600160a01b0319166001600160a01b0384169081179091556040805163288688f960e21b8152905191925063a21a23e49160048083019260209291908290030181600087803b1580156200013b57600080fd5b505af115801562000150573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001769190620002f1565b600280546001600160401b03928316600160a01b908102600160a01b600160e01b03198316811793849055604051631cd0704360e21b81529190930490931660048401523060248401526001600160a01b0391821691161790637341c10c90604401600060405180830381600087803b158015620001f357600080fd5b505af115801562000208573d6000803e3d6000fd5b50505050506200031c565b6001600160a01b0381163314156200026e5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620000a4565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215620002d257600080fd5b81516001600160a01b0381168114620002ea57600080fd5b9392505050565b6000602082840312156200030457600080fd5b81516001600160401b0381168114620002ea57600080fd5b60805160601c61108a62000342600039600081816103000152610368015261108a6000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c8063a168fa8911610097578063d8a4676f11610066578063d8a4676f14610262578063dc1670db14610287578063eb1d28bb14610290578063f2fde38b146102d557600080fd5b8063a168fa89146101bc578063ad603ea214610227578063b1e217491461023a578063d826f88f1461024357600080fd5b8063737144bc116100d3578063737144bc1461018457806374dba1241461018d57806379ba5097146101965780638da5cb5b1461019e57600080fd5b80631757f11c146101055780631fe543e3146101215780633b2bcbf114610136578063557d2e921461017b575b600080fd5b61010e60065481565b6040519081526020015b60405180910390f35b61013461012f366004610dc2565b6102e8565b005b6002546101569073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610118565b61010e60045481565b61010e60055481565b61010e60075481565b6101346103a8565b60005473ffffffffffffffffffffffffffffffffffffffff16610156565b6101fd6101ca366004610d90565b600a602052600090815260409020805460028201546003830154600484015460059094015460ff90931693919290919085565b6040805195151586526020860194909452928401919091526060830152608082015260a001610118565b610134610235366004610d32565b6104a5565b61010e60085481565b6101346000600581905560068190556103e76007556004819055600355565b610275610270366004610d90565b610809565b60405161011896959493929190610eb1565b61010e60035481565b6002546102bc9074010000000000000000000000000000000000000000900467ffffffffffffffff1681565b60405167ffffffffffffffff9091168152602001610118565b6101346102e3366004610cf5565b6108ee565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461039a576040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660248201526044015b60405180910390fd5b6103a48282610902565b5050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610429576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610391565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6104ad610a2d565b60005b8161ffff168161ffff1610156106ad576002546040517f5d3b1d300000000000000000000000000000000000000000000000000000000081526004810187905274010000000000000000000000000000000000000000820467ffffffffffffffff16602482015261ffff8816604482015263ffffffff80871660648301528516608482015260009173ffffffffffffffffffffffffffffffffffffffff1690635d3b1d309060a401602060405180830381600087803b15801561057257600080fd5b505af1158015610586573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105aa9190610da9565b6008819055905060006105bb610ab0565b6040805160c08101825260008082528251818152602080820185528084019182524284860152606084018390526080840186905260a08401839052878352600a815293909120825181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169015151781559051805194955091939092610649926001850192910190610c6a565b506040820151600282015560608201516003820155608082015160048083019190915560a090920151600590910155805490600061068683610fe6565b909155505060009182526009602052604090912055806106a581610fc4565b9150506104b0565b506002546040517f9f87fad700000000000000000000000000000000000000000000000000000000815274010000000000000000000000000000000000000000820467ffffffffffffffff16600482015230602482015273ffffffffffffffffffffffffffffffffffffffff90911690639f87fad790604401600060405180830381600087803b15801561074057600080fd5b505af1158015610754573d6000803e3d6000fd5b50506002546040517fd7ae1d3000000000000000000000000000000000000000000000000000000000815274010000000000000000000000000000000000000000820467ffffffffffffffff16600482015233602482015273ffffffffffffffffffffffffffffffffffffffff909116925063d7ae1d309150604401600060405180830381600087803b1580156107ea57600080fd5b505af11580156107fe573d6000803e3d6000fd5b505050505050505050565b6000818152600a60209081526040808320815160c081018352815460ff161515815260018201805484518187028101870190955280855260609587958695869586958695919492938584019390929083018282801561088757602002820191906000526020600020905b815481526020019060010190808311610873575b505050505081526020016002820154815260200160038201548152602001600482015481526020016005820154815250509050806000015181602001518260400151836060015184608001518560a001519650965096509650965096505091939550919395565b6108f6610a2d565b6108ff81610b4d565b50565b600061090c610ab0565b600084815260096020526040812054919250906109299083610fad565b9050600061093a82620f4240610f70565b905060065482111561094c5760068290555b600754821061095d5760075461095f565b815b60075560035461096f57806109a2565b60035461097d906001610f1d565b8160035460055461098e9190610f70565b6109989190610f1d565b6109a29190610f35565b6005556000858152600a6020908152604090912080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001908117825586516109f4939290910191870190610c6a565b506000858152600a60205260408120426003808301919091556005909101859055805491610a2183610fe6565b91905055505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610aae576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610391565b565b600046610abc81610c43565b15610b4657606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b158015610b0857600080fd5b505afa158015610b1c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b409190610da9565b91505090565b4391505090565b73ffffffffffffffffffffffffffffffffffffffff8116331415610bcd576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610391565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600061a4b1821480610c57575062066eed82145b80610c64575062066eee82145b92915050565b828054828255906000526020600020908101928215610ca5579160200282015b82811115610ca5578251825591602001919060010190610c8a565b50610cb1929150610cb5565b5090565b5b80821115610cb15760008155600101610cb6565b803561ffff81168114610cdc57600080fd5b919050565b803563ffffffff81168114610cdc57600080fd5b600060208284031215610d0757600080fd5b813573ffffffffffffffffffffffffffffffffffffffff81168114610d2b57600080fd5b9392505050565b600080600080600060a08688031215610d4a57600080fd5b610d5386610cca565b945060208601359350610d6860408701610ce1565b9250610d7660608701610ce1565b9150610d8460808701610cca565b90509295509295909350565b600060208284031215610da257600080fd5b5035919050565b600060208284031215610dbb57600080fd5b5051919050565b60008060408385031215610dd557600080fd5b8235915060208084013567ffffffffffffffff80821115610df557600080fd5b818601915086601f830112610e0957600080fd5b813581811115610e1b57610e1b61104e565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f83011681018181108582111715610e5e57610e5e61104e565b604052828152858101935084860182860187018b1015610e7d57600080fd5b600095505b83861015610ea0578035855260019590950194938601938601610e82565b508096505050505050509250929050565b600060c082018815158352602060c08185015281895180845260e086019150828b01935060005b81811015610ef457845183529383019391830191600101610ed8565b505060408501989098525050506060810193909352608083019190915260a09091015292915050565b60008219821115610f3057610f3061101f565b500190565b600082610f6b577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615610fa857610fa861101f565b500290565b600082821015610fbf57610fbf61101f565b500390565b600061ffff80831681811415610fdc57610fdc61101f565b6001019392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156110185761101861101f565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFV2OwnerTestConsumerABI = VRFV2OwnerTestConsumerMetaData.ABI

var VRFV2OwnerTestConsumerBin = VRFV2OwnerTestConsumerMetaData.Bin

func DeployVRFV2OwnerTestConsumer(auth *bind.TransactOpts, backend bind.ContractBackend, _vrfCoordinator common.Address) (common.Address, *types.Transaction, *VRFV2OwnerTestConsumer, error) {
	parsed, err := VRFV2OwnerTestConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2OwnerTestConsumerBin), backend, _vrfCoordinator)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2OwnerTestConsumer{VRFV2OwnerTestConsumerCaller: VRFV2OwnerTestConsumerCaller{contract: contract}, VRFV2OwnerTestConsumerTransactor: VRFV2OwnerTestConsumerTransactor{contract: contract}, VRFV2OwnerTestConsumerFilterer: VRFV2OwnerTestConsumerFilterer{contract: contract}}, nil
}

type VRFV2OwnerTestConsumer struct {
	address common.Address
	abi     abi.ABI
	VRFV2OwnerTestConsumerCaller
	VRFV2OwnerTestConsumerTransactor
	VRFV2OwnerTestConsumerFilterer
}

type VRFV2OwnerTestConsumerCaller struct {
	contract *bind.BoundContract
}

type VRFV2OwnerTestConsumerTransactor struct {
	contract *bind.BoundContract
}

type VRFV2OwnerTestConsumerFilterer struct {
	contract *bind.BoundContract
}

type VRFV2OwnerTestConsumerSession struct {
	Contract     *VRFV2OwnerTestConsumer
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2OwnerTestConsumerCallerSession struct {
	Contract *VRFV2OwnerTestConsumerCaller
	CallOpts bind.CallOpts
}

type VRFV2OwnerTestConsumerTransactorSession struct {
	Contract     *VRFV2OwnerTestConsumerTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2OwnerTestConsumerRaw struct {
	Contract *VRFV2OwnerTestConsumer
}

type VRFV2OwnerTestConsumerCallerRaw struct {
	Contract *VRFV2OwnerTestConsumerCaller
}

type VRFV2OwnerTestConsumerTransactorRaw struct {
	Contract *VRFV2OwnerTestConsumerTransactor
}

func NewVRFV2OwnerTestConsumer(address common.Address, backend bind.ContractBackend) (*VRFV2OwnerTestConsumer, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2OwnerTestConsumerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2OwnerTestConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2OwnerTestConsumer{address: address, abi: abi, VRFV2OwnerTestConsumerCaller: VRFV2OwnerTestConsumerCaller{contract: contract}, VRFV2OwnerTestConsumerTransactor: VRFV2OwnerTestConsumerTransactor{contract: contract}, VRFV2OwnerTestConsumerFilterer: VRFV2OwnerTestConsumerFilterer{contract: contract}}, nil
}

func NewVRFV2OwnerTestConsumerCaller(address common.Address, caller bind.ContractCaller) (*VRFV2OwnerTestConsumerCaller, error) {
	contract, err := bindVRFV2OwnerTestConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2OwnerTestConsumerCaller{contract: contract}, nil
}

func NewVRFV2OwnerTestConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2OwnerTestConsumerTransactor, error) {
	contract, err := bindVRFV2OwnerTestConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2OwnerTestConsumerTransactor{contract: contract}, nil
}

func NewVRFV2OwnerTestConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2OwnerTestConsumerFilterer, error) {
	contract, err := bindVRFV2OwnerTestConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2OwnerTestConsumerFilterer{contract: contract}, nil
}

func bindVRFV2OwnerTestConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFV2OwnerTestConsumerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2OwnerTestConsumer.Contract.VRFV2OwnerTestConsumerCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.VRFV2OwnerTestConsumerTransactor.contract.Transfer(opts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.VRFV2OwnerTestConsumerTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2OwnerTestConsumer.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.contract.Transfer(opts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCaller) COORDINATOR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2OwnerTestConsumer.contract.Call(opts, &out, "COORDINATOR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) COORDINATOR() (common.Address, error) {
	return _VRFV2OwnerTestConsumer.Contract.COORDINATOR(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCallerSession) COORDINATOR() (common.Address, error) {
	return _VRFV2OwnerTestConsumer.Contract.COORDINATOR(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCaller) GetRequestStatus(opts *bind.CallOpts, _requestId *big.Int) (GetRequestStatus,

	error) {
	var out []interface{}
	err := _VRFV2OwnerTestConsumer.contract.Call(opts, &out, "getRequestStatus", _requestId)

	outstruct := new(GetRequestStatus)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fulfilled = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.RandomWords = *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)
	outstruct.RequestTimestamp = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.FulfilmentTimestamp = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.RequestBlockNumber = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.FulfilmentBlockNumber = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) GetRequestStatus(_requestId *big.Int) (GetRequestStatus,

	error) {
	return _VRFV2OwnerTestConsumer.Contract.GetRequestStatus(&_VRFV2OwnerTestConsumer.CallOpts, _requestId)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCallerSession) GetRequestStatus(_requestId *big.Int) (GetRequestStatus,

	error) {
	return _VRFV2OwnerTestConsumer.Contract.GetRequestStatus(&_VRFV2OwnerTestConsumer.CallOpts, _requestId)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2OwnerTestConsumer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) Owner() (common.Address, error) {
	return _VRFV2OwnerTestConsumer.Contract.Owner(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCallerSession) Owner() (common.Address, error) {
	return _VRFV2OwnerTestConsumer.Contract.Owner(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCaller) SAverageFulfillmentInMillions(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2OwnerTestConsumer.contract.Call(opts, &out, "s_averageFulfillmentInMillions")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) SAverageFulfillmentInMillions() (*big.Int, error) {
	return _VRFV2OwnerTestConsumer.Contract.SAverageFulfillmentInMillions(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCallerSession) SAverageFulfillmentInMillions() (*big.Int, error) {
	return _VRFV2OwnerTestConsumer.Contract.SAverageFulfillmentInMillions(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCaller) SFastestFulfillment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2OwnerTestConsumer.contract.Call(opts, &out, "s_fastestFulfillment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) SFastestFulfillment() (*big.Int, error) {
	return _VRFV2OwnerTestConsumer.Contract.SFastestFulfillment(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCallerSession) SFastestFulfillment() (*big.Int, error) {
	return _VRFV2OwnerTestConsumer.Contract.SFastestFulfillment(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCaller) SLastRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2OwnerTestConsumer.contract.Call(opts, &out, "s_lastRequestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) SLastRequestId() (*big.Int, error) {
	return _VRFV2OwnerTestConsumer.Contract.SLastRequestId(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCallerSession) SLastRequestId() (*big.Int, error) {
	return _VRFV2OwnerTestConsumer.Contract.SLastRequestId(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCaller) SRequestCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2OwnerTestConsumer.contract.Call(opts, &out, "s_requestCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) SRequestCount() (*big.Int, error) {
	return _VRFV2OwnerTestConsumer.Contract.SRequestCount(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCallerSession) SRequestCount() (*big.Int, error) {
	return _VRFV2OwnerTestConsumer.Contract.SRequestCount(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCaller) SRequests(opts *bind.CallOpts, arg0 *big.Int) (SRequests,

	error) {
	var out []interface{}
	err := _VRFV2OwnerTestConsumer.contract.Call(opts, &out, "s_requests", arg0)

	outstruct := new(SRequests)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fulfilled = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.RequestTimestamp = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.FulfilmentTimestamp = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.RequestBlockNumber = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.FulfilmentBlockNumber = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) SRequests(arg0 *big.Int) (SRequests,

	error) {
	return _VRFV2OwnerTestConsumer.Contract.SRequests(&_VRFV2OwnerTestConsumer.CallOpts, arg0)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCallerSession) SRequests(arg0 *big.Int) (SRequests,

	error) {
	return _VRFV2OwnerTestConsumer.Contract.SRequests(&_VRFV2OwnerTestConsumer.CallOpts, arg0)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCaller) SResponseCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2OwnerTestConsumer.contract.Call(opts, &out, "s_responseCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) SResponseCount() (*big.Int, error) {
	return _VRFV2OwnerTestConsumer.Contract.SResponseCount(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCallerSession) SResponseCount() (*big.Int, error) {
	return _VRFV2OwnerTestConsumer.Contract.SResponseCount(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCaller) SSlowestFulfillment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2OwnerTestConsumer.contract.Call(opts, &out, "s_slowestFulfillment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) SSlowestFulfillment() (*big.Int, error) {
	return _VRFV2OwnerTestConsumer.Contract.SSlowestFulfillment(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCallerSession) SSlowestFulfillment() (*big.Int, error) {
	return _VRFV2OwnerTestConsumer.Contract.SSlowestFulfillment(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCaller) SubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFV2OwnerTestConsumer.contract.Call(opts, &out, "subId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) SubId() (uint64, error) {
	return _VRFV2OwnerTestConsumer.Contract.SubId(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerCallerSession) SubId() (uint64, error) {
	return _VRFV2OwnerTestConsumer.Contract.SubId(&_VRFV2OwnerTestConsumer.CallOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.contract.Transact(opts, "acceptOwnership")
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.AcceptOwnership(&_VRFV2OwnerTestConsumer.TransactOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.AcceptOwnership(&_VRFV2OwnerTestConsumer.TransactOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.RawFulfillRandomWords(&_VRFV2OwnerTestConsumer.TransactOpts, requestId, randomWords)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.RawFulfillRandomWords(&_VRFV2OwnerTestConsumer.TransactOpts, requestId, randomWords)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerTransactor) RequestRandomWords(opts *bind.TransactOpts, _requestConfirmations uint16, _keyHash [32]byte, _callbackGasLimit uint32, _numWords uint32, _requestCount uint16) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.contract.Transact(opts, "requestRandomWords", _requestConfirmations, _keyHash, _callbackGasLimit, _numWords, _requestCount)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) RequestRandomWords(_requestConfirmations uint16, _keyHash [32]byte, _callbackGasLimit uint32, _numWords uint32, _requestCount uint16) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.RequestRandomWords(&_VRFV2OwnerTestConsumer.TransactOpts, _requestConfirmations, _keyHash, _callbackGasLimit, _numWords, _requestCount)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerTransactorSession) RequestRandomWords(_requestConfirmations uint16, _keyHash [32]byte, _callbackGasLimit uint32, _numWords uint32, _requestCount uint16) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.RequestRandomWords(&_VRFV2OwnerTestConsumer.TransactOpts, _requestConfirmations, _keyHash, _callbackGasLimit, _numWords, _requestCount)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerTransactor) Reset(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.contract.Transact(opts, "reset")
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) Reset() (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.Reset(&_VRFV2OwnerTestConsumer.TransactOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerTransactorSession) Reset() (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.Reset(&_VRFV2OwnerTestConsumer.TransactOpts)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.TransferOwnership(&_VRFV2OwnerTestConsumer.TransactOpts, to)
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2OwnerTestConsumer.Contract.TransferOwnership(&_VRFV2OwnerTestConsumer.TransactOpts, to)
}

type VRFV2OwnerTestConsumerOwnershipTransferRequestedIterator struct {
	Event *VRFV2OwnerTestConsumerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2OwnerTestConsumerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2OwnerTestConsumerOwnershipTransferRequested)
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
		it.Event = new(VRFV2OwnerTestConsumerOwnershipTransferRequested)
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

func (it *VRFV2OwnerTestConsumerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFV2OwnerTestConsumerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2OwnerTestConsumerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2OwnerTestConsumerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2OwnerTestConsumer.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2OwnerTestConsumerOwnershipTransferRequestedIterator{contract: _VRFV2OwnerTestConsumer.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2OwnerTestConsumerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2OwnerTestConsumer.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2OwnerTestConsumerOwnershipTransferRequested)
				if err := _VRFV2OwnerTestConsumer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFV2OwnerTestConsumerOwnershipTransferRequested, error) {
	event := new(VRFV2OwnerTestConsumerOwnershipTransferRequested)
	if err := _VRFV2OwnerTestConsumer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2OwnerTestConsumerOwnershipTransferredIterator struct {
	Event *VRFV2OwnerTestConsumerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2OwnerTestConsumerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2OwnerTestConsumerOwnershipTransferred)
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
		it.Event = new(VRFV2OwnerTestConsumerOwnershipTransferred)
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

func (it *VRFV2OwnerTestConsumerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFV2OwnerTestConsumerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2OwnerTestConsumerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2OwnerTestConsumerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2OwnerTestConsumer.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2OwnerTestConsumerOwnershipTransferredIterator{contract: _VRFV2OwnerTestConsumer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2OwnerTestConsumerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2OwnerTestConsumer.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2OwnerTestConsumerOwnershipTransferred)
				if err := _VRFV2OwnerTestConsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumerFilterer) ParseOwnershipTransferred(log types.Log) (*VRFV2OwnerTestConsumerOwnershipTransferred, error) {
	event := new(VRFV2OwnerTestConsumerOwnershipTransferred)
	if err := _VRFV2OwnerTestConsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetRequestStatus struct {
	Fulfilled             bool
	RandomWords           []*big.Int
	RequestTimestamp      *big.Int
	FulfilmentTimestamp   *big.Int
	RequestBlockNumber    *big.Int
	FulfilmentBlockNumber *big.Int
}
type SRequests struct {
	Fulfilled             bool
	RequestTimestamp      *big.Int
	FulfilmentTimestamp   *big.Int
	RequestBlockNumber    *big.Int
	FulfilmentBlockNumber *big.Int
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumer) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFV2OwnerTestConsumer.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFV2OwnerTestConsumer.ParseOwnershipTransferRequested(log)
	case _VRFV2OwnerTestConsumer.abi.Events["OwnershipTransferred"].ID:
		return _VRFV2OwnerTestConsumer.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFV2OwnerTestConsumerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFV2OwnerTestConsumerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_VRFV2OwnerTestConsumer *VRFV2OwnerTestConsumer) Address() common.Address {
	return _VRFV2OwnerTestConsumer.address
}

type VRFV2OwnerTestConsumerInterface interface {
	COORDINATOR(opts *bind.CallOpts) (common.Address, error)

	GetRequestStatus(opts *bind.CallOpts, _requestId *big.Int) (GetRequestStatus,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SAverageFulfillmentInMillions(opts *bind.CallOpts) (*big.Int, error)

	SFastestFulfillment(opts *bind.CallOpts) (*big.Int, error)

	SLastRequestId(opts *bind.CallOpts) (*big.Int, error)

	SRequestCount(opts *bind.CallOpts) (*big.Int, error)

	SRequests(opts *bind.CallOpts, arg0 *big.Int) (SRequests,

		error)

	SResponseCount(opts *bind.CallOpts) (*big.Int, error)

	SSlowestFulfillment(opts *bind.CallOpts) (*big.Int, error)

	SubId(opts *bind.CallOpts) (uint64, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts, _requestConfirmations uint16, _keyHash [32]byte, _callbackGasLimit uint32, _numWords uint32, _requestCount uint16) (*types.Transaction, error)

	Reset(opts *bind.TransactOpts) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2OwnerTestConsumerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2OwnerTestConsumerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFV2OwnerTestConsumerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2OwnerTestConsumerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2OwnerTestConsumerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFV2OwnerTestConsumerOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
