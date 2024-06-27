// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_coordinator_v2_plus_v2_example

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

type VRFV2PlusClientRandomWordsRequest struct {
	KeyHash              [32]byte
	SubId                *big.Int
	RequestConfirmations uint16
	CallbackGasLimit     uint32
	NumWords             uint32
	ExtraArgs            []byte
}

var VRFCoordinatorV2PlusV2ExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"prevCoordinator\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"transferredValue\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"expectedValue\",\"type\":\"uint96\"}],\"name\":\"InvalidNativeBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"requestVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"expectedVersion\",\"type\":\"uint8\"}],\"name\":\"InvalidVersion\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"previousCoordinator\",\"type\":\"address\"}],\"name\":\"MustBePreviousCoordinator\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestID\",\"type\":\"uint256\"}],\"name\":\"generateFakeRandomness\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"},{\"internalType\":\"uint96\",\"name\":\"linkBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"nativeBalance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedData\",\"type\":\"bytes\"}],\"name\":\"onMigration\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFV2PlusClient.RandomWordsRequest\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"requestRandomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_link\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_prevCoordinator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requestConsumerMapping\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_subscriptions\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"linkBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"nativeBalance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalLinkBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalNativeBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052600060035534801561001557600080fd5b50604051610f65380380610f6583398101604081905261003491610081565b600580546001600160a01b039384166001600160a01b031991821617909155600480549290931691161790556100b4565b80516001600160a01b038116811461007c57600080fd5b919050565b6000806040838503121561009457600080fd5b61009d83610065565b91506100ab60208401610065565b90509250929050565b610ea2806100c36000396000f3fe6080604052600436106100c75760003560e01c8063ce3f471911610074578063dc311dd31161004e578063dc311dd314610325578063e89e106a14610355578063ed8b558f1461036b57600080fd5b8063ce3f4719146102c3578063d6100d1c146102d8578063da4f5e6d146102f857600080fd5b806386175f58116100a557806386175f581461022557806393f3acb6146102685780639b1c385e1461029557600080fd5b80630495f265146100cc578063086597b31461018157806318e3dd27146101d3575b600080fd5b3480156100d857600080fd5b5061013b6100e7366004610c86565b6000602081905290815260409020805460029091015473ffffffffffffffffffffffffffffffffffffffff909116906bffffffffffffffffffffffff808216916c0100000000000000000000000090041683565b6040805173ffffffffffffffffffffffffffffffffffffffff90941684526bffffffffffffffffffffffff92831660208501529116908201526060015b60405180910390f35b34801561018d57600080fd5b506004546101ae9073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610178565b3480156101df57600080fd5b50600254610208906c0100000000000000000000000090046bffffffffffffffffffffffff1681565b6040516bffffffffffffffffffffffff9091168152602001610178565b34801561023157600080fd5b506101ae610240366004610c86565b60016020526000908152604090205473ffffffffffffffffffffffffffffffffffffffff1681565b34801561027457600080fd5b50610288610283366004610c86565b610390565b6040516101789190610d63565b3480156102a157600080fd5b506102b56102b0366004610b81565b61043b565b604051908152602001610178565b6102d66102d1366004610b0f565b61044c565b005b3480156102e457600080fd5b506102d66102f3366004610c86565b610746565b34801561030457600080fd5b506005546101ae9073ffffffffffffffffffffffffffffffffffffffff1681565b34801561033157600080fd5b50610345610340366004610c86565b6107ce565b6040516101789493929190610cda565b34801561036157600080fd5b506102b560035481565b34801561037757600080fd5b50600254610208906bffffffffffffffffffffffff1681565b6040805160018082528183019092526060916000919060208083019080368337019050509050826040516020016103fe918152604060208201819052600a908201527f6e6f742072616e646f6d00000000000000000000000000000000000000000000606082015260800190565b6040516020818303038152906040528051906020012060001c8160008151811061042a5761042a610e37565b602090810291909101015292915050565b6000610446336108f9565b92915050565b60045473ffffffffffffffffffffffffffffffffffffffff1633146104c557600480546040517ff5828f73000000000000000000000000000000000000000000000000000000008152339281019290925273ffffffffffffffffffffffffffffffffffffffff1660248201526044015b60405180910390fd5b60006104d382840184610bc3565b9050806000015160ff166001146105255780516040517f8df4607c00000000000000000000000000000000000000000000000000000000815260ff9091166004820152600160248201526044016104bc565b8060a001516bffffffffffffffffffffffff16341461058c5760a08101516040517f6acf13500000000000000000000000000000000000000000000000000000000081523460048201526bffffffffffffffffffffffff90911660248201526044016104bc565b60408051608080820183528383015173ffffffffffffffffffffffffffffffffffffffff90811683526060808601516020808601918252938701516bffffffffffffffffffffffff9081168688015260a0880151169185019190915282860151600090815280845294909420835181547fffffffffffffffffffffffff00000000000000000000000000000000000000001692169190911781559251805192939261063d9260018501920190610965565b506040820151600291820180546060909401516bffffffffffffffffffffffff9283167fffffffffffffffff000000000000000000000000000000000000000000000000909516949094176c01000000000000000000000000948316850217905560a084015182549093600c926106b8928692900416610dd8565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508060800151600260008282829054906101000a90046bffffffffffffffffffffffff166107139190610dd8565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550505050565b60008181526001602052604090205473ffffffffffffffffffffffffffffffffffffffff1680631fe543e38361077b81610390565b6040518363ffffffff1660e01b8152600401610798929190610d76565b600060405180830381600087803b1580156107b257600080fd5b505af11580156107c6573d6000803e3d6000fd5b505050505050565b6000818152602081905260408120546060908290819073ffffffffffffffffffffffffffffffffffffffff16610830576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008581526020818152604091829020805460028201546001909201805485518186028101860190965280865273ffffffffffffffffffffffffffffffffffffffff9092169490936bffffffffffffffffffffffff808516946c010000000000000000000000009004169285918301828280156108e357602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116108b8575b5050505050925093509350935093509193509193565b6000600354600161090a9190610dc0565b6003819055600081815260016020526040902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff94909416939093179092555090565b8280548282559060005260206000209081019282156109df579160200282015b828111156109df57825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190610985565b506109eb9291506109ef565b5090565b5b808211156109eb57600081556001016109f0565b803573ffffffffffffffffffffffffffffffffffffffff81168114610a2857600080fd5b919050565b600082601f830112610a3e57600080fd5b8135602067ffffffffffffffff80831115610a5b57610a5b610e66565b8260051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f83011681018181108482111715610a9e57610a9e610e66565b60405284815283810192508684018288018501891015610abd57600080fd5b600092505b85831015610ae757610ad381610a04565b845292840192600192909201918401610ac2565b50979650505050505050565b80356bffffffffffffffffffffffff81168114610a2857600080fd5b60008060208385031215610b2257600080fd5b823567ffffffffffffffff80821115610b3a57600080fd5b818501915085601f830112610b4e57600080fd5b813581811115610b5d57600080fd5b866020828501011115610b6f57600080fd5b60209290920196919550909350505050565b600060208284031215610b9357600080fd5b813567ffffffffffffffff811115610baa57600080fd5b820160c08185031215610bbc57600080fd5b9392505050565b600060208284031215610bd557600080fd5b813567ffffffffffffffff80821115610bed57600080fd5b9083019060c08286031215610c0157600080fd5b610c09610d97565b823560ff81168114610c1a57600080fd5b815260208381013590820152610c3260408401610a04565b6040820152606083013582811115610c4957600080fd5b610c5587828601610a2d565b606083015250610c6760808401610af3565b6080820152610c7860a08401610af3565b60a082015295945050505050565b600060208284031215610c9857600080fd5b5035919050565b600081518084526020808501945080840160005b83811015610ccf57815187529582019590820190600101610cb3565b509495945050505050565b60006080820173ffffffffffffffffffffffffffffffffffffffff8088168452602060808186015282885180855260a087019150828a01945060005b81811015610d34578551851683529483019491830191600101610d16565b50506bffffffffffffffffffffffff978816604087015295909616606090940193909352509195945050505050565b602081526000610bbc6020830184610c9f565b828152604060208201526000610d8f6040830184610c9f565b949350505050565b60405160c0810167ffffffffffffffff81118282101715610dba57610dba610e66565b60405290565b60008219821115610dd357610dd3610e08565b500190565b60006bffffffffffffffffffffffff808316818516808303821115610dff57610dff610e08565b01949350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFCoordinatorV2PlusV2ExampleABI = VRFCoordinatorV2PlusV2ExampleMetaData.ABI

var VRFCoordinatorV2PlusV2ExampleBin = VRFCoordinatorV2PlusV2ExampleMetaData.Bin

func DeployVRFCoordinatorV2PlusV2Example(auth *bind.TransactOpts, backend bind.ContractBackend, link common.Address, prevCoordinator common.Address) (common.Address, *types.Transaction, *VRFCoordinatorV2PlusV2Example, error) {
	parsed, err := VRFCoordinatorV2PlusV2ExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFCoordinatorV2PlusV2ExampleBin), backend, link, prevCoordinator)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFCoordinatorV2PlusV2Example{VRFCoordinatorV2PlusV2ExampleCaller: VRFCoordinatorV2PlusV2ExampleCaller{contract: contract}, VRFCoordinatorV2PlusV2ExampleTransactor: VRFCoordinatorV2PlusV2ExampleTransactor{contract: contract}, VRFCoordinatorV2PlusV2ExampleFilterer: VRFCoordinatorV2PlusV2ExampleFilterer{contract: contract}}, nil
}

type VRFCoordinatorV2PlusV2Example struct {
	address common.Address
	abi     abi.ABI
	VRFCoordinatorV2PlusV2ExampleCaller
	VRFCoordinatorV2PlusV2ExampleTransactor
	VRFCoordinatorV2PlusV2ExampleFilterer
}

type VRFCoordinatorV2PlusV2ExampleCaller struct {
	contract *bind.BoundContract
}

type VRFCoordinatorV2PlusV2ExampleTransactor struct {
	contract *bind.BoundContract
}

type VRFCoordinatorV2PlusV2ExampleFilterer struct {
	contract *bind.BoundContract
}

type VRFCoordinatorV2PlusV2ExampleSession struct {
	Contract     *VRFCoordinatorV2PlusV2Example
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorV2PlusV2ExampleCallerSession struct {
	Contract *VRFCoordinatorV2PlusV2ExampleCaller
	CallOpts bind.CallOpts
}

type VRFCoordinatorV2PlusV2ExampleTransactorSession struct {
	Contract     *VRFCoordinatorV2PlusV2ExampleTransactor
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorV2PlusV2ExampleRaw struct {
	Contract *VRFCoordinatorV2PlusV2Example
}

type VRFCoordinatorV2PlusV2ExampleCallerRaw struct {
	Contract *VRFCoordinatorV2PlusV2ExampleCaller
}

type VRFCoordinatorV2PlusV2ExampleTransactorRaw struct {
	Contract *VRFCoordinatorV2PlusV2ExampleTransactor
}

func NewVRFCoordinatorV2PlusV2Example(address common.Address, backend bind.ContractBackend) (*VRFCoordinatorV2PlusV2Example, error) {
	abi, err := abi.JSON(strings.NewReader(VRFCoordinatorV2PlusV2ExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFCoordinatorV2PlusV2Example(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusV2Example{address: address, abi: abi, VRFCoordinatorV2PlusV2ExampleCaller: VRFCoordinatorV2PlusV2ExampleCaller{contract: contract}, VRFCoordinatorV2PlusV2ExampleTransactor: VRFCoordinatorV2PlusV2ExampleTransactor{contract: contract}, VRFCoordinatorV2PlusV2ExampleFilterer: VRFCoordinatorV2PlusV2ExampleFilterer{contract: contract}}, nil
}

func NewVRFCoordinatorV2PlusV2ExampleCaller(address common.Address, caller bind.ContractCaller) (*VRFCoordinatorV2PlusV2ExampleCaller, error) {
	contract, err := bindVRFCoordinatorV2PlusV2Example(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusV2ExampleCaller{contract: contract}, nil
}

func NewVRFCoordinatorV2PlusV2ExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFCoordinatorV2PlusV2ExampleTransactor, error) {
	contract, err := bindVRFCoordinatorV2PlusV2Example(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusV2ExampleTransactor{contract: contract}, nil
}

func NewVRFCoordinatorV2PlusV2ExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFCoordinatorV2PlusV2ExampleFilterer, error) {
	contract, err := bindVRFCoordinatorV2PlusV2Example(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusV2ExampleFilterer{contract: contract}, nil
}

func bindVRFCoordinatorV2PlusV2Example(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFCoordinatorV2PlusV2ExampleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorV2PlusV2Example.Contract.VRFCoordinatorV2PlusV2ExampleCaller.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.VRFCoordinatorV2PlusV2ExampleTransactor.contract.Transfer(opts)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.VRFCoordinatorV2PlusV2ExampleTransactor.contract.Transact(opts, method, params...)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorV2PlusV2Example.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.contract.Transfer(opts)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.contract.Transact(opts, method, params...)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleCaller) GenerateFakeRandomness(opts *bind.CallOpts, requestID *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusV2Example.contract.Call(opts, &out, "generateFakeRandomness", requestID)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleSession) GenerateFakeRandomness(requestID *big.Int) ([]*big.Int, error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.GenerateFakeRandomness(&_VRFCoordinatorV2PlusV2Example.CallOpts, requestID)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleCallerSession) GenerateFakeRandomness(requestID *big.Int) ([]*big.Int, error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.GenerateFakeRandomness(&_VRFCoordinatorV2PlusV2Example.CallOpts, requestID)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleCaller) GetSubscription(opts *bind.CallOpts, subId *big.Int) (GetSubscription,

	error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusV2Example.contract.Call(opts, &out, "getSubscription", subId)

	outstruct := new(GetSubscription)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Owner = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Consumers = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)
	outstruct.LinkBalance = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.NativeBalance = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleSession) GetSubscription(subId *big.Int) (GetSubscription,

	error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.GetSubscription(&_VRFCoordinatorV2PlusV2Example.CallOpts, subId)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleCallerSession) GetSubscription(subId *big.Int) (GetSubscription,

	error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.GetSubscription(&_VRFCoordinatorV2PlusV2Example.CallOpts, subId)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleCaller) SLink(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusV2Example.contract.Call(opts, &out, "s_link")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleSession) SLink() (common.Address, error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.SLink(&_VRFCoordinatorV2PlusV2Example.CallOpts)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleCallerSession) SLink() (common.Address, error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.SLink(&_VRFCoordinatorV2PlusV2Example.CallOpts)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleCaller) SPrevCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusV2Example.contract.Call(opts, &out, "s_prevCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleSession) SPrevCoordinator() (common.Address, error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.SPrevCoordinator(&_VRFCoordinatorV2PlusV2Example.CallOpts)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleCallerSession) SPrevCoordinator() (common.Address, error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.SPrevCoordinator(&_VRFCoordinatorV2PlusV2Example.CallOpts)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleCaller) SRequestConsumerMapping(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusV2Example.contract.Call(opts, &out, "s_requestConsumerMapping", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleSession) SRequestConsumerMapping(arg0 *big.Int) (common.Address, error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.SRequestConsumerMapping(&_VRFCoordinatorV2PlusV2Example.CallOpts, arg0)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleCallerSession) SRequestConsumerMapping(arg0 *big.Int) (common.Address, error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.SRequestConsumerMapping(&_VRFCoordinatorV2PlusV2Example.CallOpts, arg0)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleCaller) SRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusV2Example.contract.Call(opts, &out, "s_requestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleSession) SRequestId() (*big.Int, error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.SRequestId(&_VRFCoordinatorV2PlusV2Example.CallOpts)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleCallerSession) SRequestId() (*big.Int, error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.SRequestId(&_VRFCoordinatorV2PlusV2Example.CallOpts)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleCaller) SSubscriptions(opts *bind.CallOpts, arg0 *big.Int) (SSubscriptions,

	error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusV2Example.contract.Call(opts, &out, "s_subscriptions", arg0)

	outstruct := new(SSubscriptions)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Owner = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.LinkBalance = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.NativeBalance = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleSession) SSubscriptions(arg0 *big.Int) (SSubscriptions,

	error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.SSubscriptions(&_VRFCoordinatorV2PlusV2Example.CallOpts, arg0)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleCallerSession) SSubscriptions(arg0 *big.Int) (SSubscriptions,

	error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.SSubscriptions(&_VRFCoordinatorV2PlusV2Example.CallOpts, arg0)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleCaller) STotalLinkBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusV2Example.contract.Call(opts, &out, "s_totalLinkBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleSession) STotalLinkBalance() (*big.Int, error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.STotalLinkBalance(&_VRFCoordinatorV2PlusV2Example.CallOpts)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleCallerSession) STotalLinkBalance() (*big.Int, error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.STotalLinkBalance(&_VRFCoordinatorV2PlusV2Example.CallOpts)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleCaller) STotalNativeBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusV2Example.contract.Call(opts, &out, "s_totalNativeBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleSession) STotalNativeBalance() (*big.Int, error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.STotalNativeBalance(&_VRFCoordinatorV2PlusV2Example.CallOpts)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleCallerSession) STotalNativeBalance() (*big.Int, error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.STotalNativeBalance(&_VRFCoordinatorV2PlusV2Example.CallOpts)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleTransactor) FulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusV2Example.contract.Transact(opts, "fulfillRandomWords", requestId)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleSession) FulfillRandomWords(requestId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.FulfillRandomWords(&_VRFCoordinatorV2PlusV2Example.TransactOpts, requestId)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleTransactorSession) FulfillRandomWords(requestId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.FulfillRandomWords(&_VRFCoordinatorV2PlusV2Example.TransactOpts, requestId)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleTransactor) OnMigration(opts *bind.TransactOpts, encodedData []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusV2Example.contract.Transact(opts, "onMigration", encodedData)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleSession) OnMigration(encodedData []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.OnMigration(&_VRFCoordinatorV2PlusV2Example.TransactOpts, encodedData)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleTransactorSession) OnMigration(encodedData []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.OnMigration(&_VRFCoordinatorV2PlusV2Example.TransactOpts, encodedData)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleTransactor) RequestRandomWords(opts *bind.TransactOpts, arg0 VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusV2Example.contract.Transact(opts, "requestRandomWords", arg0)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleSession) RequestRandomWords(arg0 VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.RequestRandomWords(&_VRFCoordinatorV2PlusV2Example.TransactOpts, arg0)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleTransactorSession) RequestRandomWords(arg0 VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.RequestRandomWords(&_VRFCoordinatorV2PlusV2Example.TransactOpts, arg0)
}

type GetSubscription struct {
	Owner         common.Address
	Consumers     []common.Address
	LinkBalance   *big.Int
	NativeBalance *big.Int
}
type SSubscriptions struct {
	Owner         common.Address
	LinkBalance   *big.Int
	NativeBalance *big.Int
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2Example) Address() common.Address {
	return _VRFCoordinatorV2PlusV2Example.address
}

type VRFCoordinatorV2PlusV2ExampleInterface interface {
	GenerateFakeRandomness(opts *bind.CallOpts, requestID *big.Int) ([]*big.Int, error)

	GetSubscription(opts *bind.CallOpts, subId *big.Int) (GetSubscription,

		error)

	SLink(opts *bind.CallOpts) (common.Address, error)

	SPrevCoordinator(opts *bind.CallOpts) (common.Address, error)

	SRequestConsumerMapping(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error)

	SRequestId(opts *bind.CallOpts) (*big.Int, error)

	SSubscriptions(opts *bind.CallOpts, arg0 *big.Int) (SSubscriptions,

		error)

	STotalLinkBalance(opts *bind.CallOpts) (*big.Int, error)

	STotalNativeBalance(opts *bind.CallOpts) (*big.Int, error)

	FulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int) (*types.Transaction, error)

	OnMigration(opts *bind.TransactOpts, encodedData []byte) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts, arg0 VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error)

	Address() common.Address
}
