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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"prevCoordinator\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"transferredValue\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"expectedValue\",\"type\":\"uint96\"}],\"name\":\"InvalidNativeBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"requestVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"expectedVersion\",\"type\":\"uint8\"}],\"name\":\"InvalidVersion\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"previousCoordinator\",\"type\":\"address\"}],\"name\":\"MustBePreviousCoordinator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SubscriptionIDCollisionFound\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestID\",\"type\":\"uint256\"}],\"name\":\"generateFakeRandomness\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"linkBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"nativeBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedData\",\"type\":\"bytes\"}],\"name\":\"onMigration\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFV2PlusClient.RandomWordsRequest\",\"name\":\"req\",\"type\":\"tuple\"}],\"name\":\"requestRandomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_link\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_prevCoordinator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requestConsumerMapping\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_subscriptions\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"linkBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"nativeBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalLinkBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalNativeBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052600060035534801561001557600080fd5b506040516111d73803806111d783398101604081905261003491610081565b600580546001600160a01b039384166001600160a01b031991821617909155600480549290931691161790556100b4565b80516001600160a01b038116811461007c57600080fd5b919050565b6000806040838503121561009457600080fd5b61009d83610065565b91506100ab60208401610065565b90509250929050565b611114806100c36000396000f3fe6080604052600436106100c75760003560e01c8063ce3f471911610074578063dc311dd31161004e578063dc311dd314610361578063e89e106a14610392578063ed8b558f146103a857600080fd5b8063ce3f4719146102ff578063d6100d1c14610314578063da4f5e6d1461033457600080fd5b806386175f58116100a557806386175f581461026157806393f3acb6146102a45780639b1c385e146102d157600080fd5b80630495f265146100cc578063086597b3146101bd57806318e3dd271461020f575b600080fd5b3480156100d857600080fd5b506101636100e7366004610ec8565b600060208190529081526040902080546001909101546bffffffffffffffffffffffff808316926c01000000000000000000000000810490911691780100000000000000000000000000000000000000000000000090910467ffffffffffffffff169073ffffffffffffffffffffffffffffffffffffffff1684565b604080516bffffffffffffffffffffffff958616815294909316602085015267ffffffffffffffff9091169183019190915273ffffffffffffffffffffffffffffffffffffffff1660608201526080015b60405180910390f35b3480156101c957600080fd5b506004546101ea9073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101b4565b34801561021b57600080fd5b50600254610244906c0100000000000000000000000090046bffffffffffffffffffffffff1681565b6040516bffffffffffffffffffffffff90911681526020016101b4565b34801561026d57600080fd5b506101ea61027c366004610ec8565b60016020526000908152604090205473ffffffffffffffffffffffffffffffffffffffff1681565b3480156102b057600080fd5b506102c46102bf366004610ec8565b6103cd565b6040516101b49190610f1c565b3480156102dd57600080fd5b506102f16102ec366004610dca565b610478565b6040519081526020016101b4565b61031261030d366004610d58565b6105aa565b005b34801561032057600080fd5b5061031261032f366004610ec8565b61095e565b34801561034057600080fd5b506005546101ea9073ffffffffffffffffffffffffffffffffffffffff1681565b34801561036d57600080fd5b5061038161037c366004610ec8565b6109e6565b6040516101b4959493929190610f50565b34801561039e57600080fd5b506102f160035481565b3480156103b457600080fd5b50600254610244906bffffffffffffffffffffffff1681565b60408051600180825281830190925260609160009190602080830190803683370190505090508260405160200161043b918152604060208201819052600a908201527f6e6f742072616e646f6d00000000000000000000000000000000000000000000606082015260800190565b6040516020818303038152906040528051906020012060001c81600081518110610467576104676110a9565b602090810291909101015292915050565b60208181013560009081528082526040808220815160a08101835281546bffffffffffffffffffffffff80821683526c01000000000000000000000000820416828701527801000000000000000000000000000000000000000000000000900467ffffffffffffffff1681840152600182015473ffffffffffffffffffffffffffffffffffffffff166060820152600282018054845181880281018801909552808552949586959294608086019390929183018282801561056f57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610544575b50505050508152505090508060400151600161058b919061102b565b67ffffffffffffffff1660408201526105a333610b42565b9392505050565b60045473ffffffffffffffffffffffffffffffffffffffff16331461062357600480546040517ff5828f73000000000000000000000000000000000000000000000000000000008152339281019290925273ffffffffffffffffffffffffffffffffffffffff1660248201526044015b60405180910390fd5b600061063182840184610e05565b9050806000015160ff166001146106835780516040517f8df4607c00000000000000000000000000000000000000000000000000000000815260ff90911660048201526001602482015260440161061a565b8060a001516bffffffffffffffffffffffff1634146106ea5760a08101516040517f6acf13500000000000000000000000000000000000000000000000000000000081523460048201526bffffffffffffffffffffffff909116602482015260440161061a565b602080820151600090815290819052604090206001015473ffffffffffffffffffffffffffffffffffffffff161561074e576040517f4d5f486a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160a080820183526080808501516bffffffffffffffffffffffff9081168452918501518216602080850191825260008587018181528888015173ffffffffffffffffffffffffffffffffffffffff9081166060808a019182528b0151968901968752848b0151845283855298909220875181549551925167ffffffffffffffff1678010000000000000000000000000000000000000000000000000277ffffffffffffffffffffffffffffffffffffffffffffffff9389166c01000000000000000000000000027fffffffffffffffff000000000000000000000000000000000000000000000000909716919098161794909417169490941782559451600182018054919094167fffffffffffffffffffffffff000000000000000000000000000000000000000090911617909255518051929391926108989260028501920190610bae565b50505060a081015160028054600c906108d09084906c0100000000000000000000000090046bffffffffffffffffffffffff16611057565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508060800151600260008282829054906101000a90046bffffffffffffffffffffffff1661092b9190611057565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550505050565b60008181526001602052604090205473ffffffffffffffffffffffffffffffffffffffff1680631fe543e383610993816103cd565b6040518363ffffffff1660e01b81526004016109b0929190610f2f565b600060405180830381600087803b1580156109ca57600080fd5b505af11580156109de573d6000803e3d6000fd5b505050505050565b60008181526020819052604081206001015481908190819060609073ffffffffffffffffffffffffffffffffffffffff16610a4d576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000868152602081815260409182902080546001820154600290920180548551818602810186019096528086526bffffffffffffffffffffffff808416966c01000000000000000000000000850490911695780100000000000000000000000000000000000000000000000090940467ffffffffffffffff169473ffffffffffffffffffffffffffffffffffffffff169390918391830182828015610b2857602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610afd575b505050505090509450945094509450945091939590929450565b60006003546001610b539190611013565b6003819055600081815260016020526040902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff94909416939093179092555090565b828054828255906000526020600020908101928215610c28579160200282015b82811115610c2857825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190610bce565b50610c34929150610c38565b5090565b5b80821115610c345760008155600101610c39565b803573ffffffffffffffffffffffffffffffffffffffff81168114610c7157600080fd5b919050565b600082601f830112610c8757600080fd5b8135602067ffffffffffffffff80831115610ca457610ca46110d8565b8260051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f83011681018181108482111715610ce757610ce76110d8565b60405284815283810192508684018288018501891015610d0657600080fd5b600092505b85831015610d3057610d1c81610c4d565b845292840192600192909201918401610d0b565b50979650505050505050565b80356bffffffffffffffffffffffff81168114610c7157600080fd5b60008060208385031215610d6b57600080fd5b823567ffffffffffffffff80821115610d8357600080fd5b818501915085601f830112610d9757600080fd5b813581811115610da657600080fd5b866020828501011115610db857600080fd5b60209290920196919550909350505050565b600060208284031215610ddc57600080fd5b813567ffffffffffffffff811115610df357600080fd5b820160c081850312156105a357600080fd5b600060208284031215610e1757600080fd5b813567ffffffffffffffff80821115610e2f57600080fd5b9083019060c08286031215610e4357600080fd5b610e4b610fea565b823560ff81168114610e5c57600080fd5b815260208381013590820152610e7460408401610c4d565b6040820152606083013582811115610e8b57600080fd5b610e9787828601610c76565b606083015250610ea960808401610d3c565b6080820152610eba60a08401610d3c565b60a082015295945050505050565b600060208284031215610eda57600080fd5b5035919050565b600081518084526020808501945080840160005b83811015610f1157815187529582019590820190600101610ef5565b509495945050505050565b6020815260006105a36020830184610ee1565b828152604060208201526000610f486040830184610ee1565b949350505050565b600060a082016bffffffffffffffffffffffff808916845260208189168186015267ffffffffffffffff8816604086015273ffffffffffffffffffffffffffffffffffffffff9150818716606086015260a0608086015282865180855260c087019150828801945060005b81811015610fd9578551851683529483019491830191600101610fbb565b50909b9a5050505050505050505050565b60405160c0810167ffffffffffffffff8111828210171561100d5761100d6110d8565b60405290565b600082198211156110265761102661107a565b500190565b600067ffffffffffffffff80831681851680830382111561104e5761104e61107a565b01949350505050565b60006bffffffffffffffffffffffff80831681851680830382111561104e5761104e5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
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

	outstruct.LinkBalance = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.NativeBalance = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.ReqCount = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.Owner = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.Consumers = *abi.ConvertType(out[4], new([]common.Address)).(*[]common.Address)

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

	outstruct.LinkBalance = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.NativeBalance = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.ReqCount = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.Owner = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)

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

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleTransactor) RequestRandomWords(opts *bind.TransactOpts, req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusV2Example.contract.Transact(opts, "requestRandomWords", req)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleSession) RequestRandomWords(req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.RequestRandomWords(&_VRFCoordinatorV2PlusV2Example.TransactOpts, req)
}

func (_VRFCoordinatorV2PlusV2Example *VRFCoordinatorV2PlusV2ExampleTransactorSession) RequestRandomWords(req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusV2Example.Contract.RequestRandomWords(&_VRFCoordinatorV2PlusV2Example.TransactOpts, req)
}

type GetSubscription struct {
	LinkBalance   *big.Int
	NativeBalance *big.Int
	ReqCount      uint64
	Owner         common.Address
	Consumers     []common.Address
}
type SSubscriptions struct {
	LinkBalance   *big.Int
	NativeBalance *big.Int
	ReqCount      uint64
	Owner         common.Address
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

	RequestRandomWords(opts *bind.TransactOpts, req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error)

	Address() common.Address
}
