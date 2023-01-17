// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package lottery_consumer

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

type LotteryConsumerLotteryOutcome struct {
	ClientRequestId      [32]byte
	LotteryType          uint8
	VrfExternalRequestId *big.Int
	WinningNumbers       []uint8
}

type LotteryConsumerLotteryRequest struct {
	ClientRequestId      [32]byte
	LotteryType          uint8
	VrfExternalRequestId *big.Int
}

var LotteryConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"UnallowedCaller\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"AllowedCallerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"AllowedCallerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"vrfRequestId\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"clientRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"lotteryType\",\"type\":\"uint8\"},{\"internalType\":\"uint128\",\"name\":\"vrfExternalRequestId\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structLotteryConsumer.LotteryRequest\",\"name\":\"request\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"clientRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"lotteryType\",\"type\":\"uint8\"},{\"internalType\":\"uint128\",\"name\":\"vrfExternalRequestId\",\"type\":\"uint128\"},{\"internalType\":\"uint8[]\",\"name\":\"winningNumbers\",\"type\":\"uint8[]\"}],\"indexed\":false,\"internalType\":\"structLotteryConsumer.LotteryOutcome\",\"name\":\"outcome\",\"type\":\"tuple\"}],\"name\":\"LotterySettled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"vrfRequestId\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"clientRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"lotteryType\",\"type\":\"uint8\"},{\"internalType\":\"uint128\",\"name\":\"vrfExternalRequestId\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structLotteryConsumer.LotteryRequest\",\"name\":\"request\",\"type\":\"tuple\"}],\"name\":\"LotteryStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMostRecentVrfRequestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRequestConfig\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"minRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"clientRequestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"lotteryType\",\"type\":\"uint8\"},{\"internalType\":\"uint128\",\"name\":\"vrfExternalRequestId\",\"type\":\"uint128\"}],\"name\":\"requestRandomness\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"minRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"}],\"name\":\"setRequestConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"randomness\",\"type\":\"uint256[]\"}],\"name\":\"shuffle35\",\"outputs\":[{\"internalType\":\"uint8[]\",\"name\":\"\",\"type\":\"uint8[]\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b506040516200140f3803806200140f8339810160408190526200003491620001fc565b8033806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf8162000150565b50505060601b6001600160601b0319166080526001600160a01b0381166200012a5760405162461bcd60e51b815260206004820181905260248201527f76726620636f6f7264696e61746f72206d757374206265206e6f6e2d7a65726f604482015260640162000083565b600880546001600160a01b0319166001600160a01b03929092169190911790556200022e565b6001600160a01b038116331415620001ab5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156200020f57600080fd5b81516001600160a01b03811681146200022757600080fd5b9392505050565b60805160601c6111bb6200025460003960008181610424015261048c01526111bb6000f3fe608060405234801561001057600080fd5b50600436106100a25760003560e01c806379ba509711610076578063b575d2af1161005b578063b575d2af146101b3578063f2fde38b146101d3578063f35c71ae146101e657600080fd5b806379ba5097146101835780638da5cb5b1461018b57600080fd5b8062012291146100a757806315ef3ef2146101485780631fe543e31461015d5780634f5b9f0d14610170575b600080fd5b6040805160a0808201835260025480835260035467ffffffffffffffff8116602080860182905261ffff6801000000000000000084041686880181905263ffffffff6a0100000000000000000000850481166060808a018290526e0100000000000000000000000000009096049091166080988901819052895196875292860193909352968401969096529082015291820192909252015b60405180910390f35b61015b610156366004610ee7565b6101f7565b005b61015b61016b366004610f5d565b61040c565b61015b61017e366004610e71565b6104cc565b61015b6105a3565b60005460405173ffffffffffffffffffffffffffffffffffffffff909116815260200161013f565b6101c66101c1366004610e34565b6106a0565b60405161013f9190610fe2565b61015b6101e1366004610df7565b61086a565b60075460405190815260200161013f565b6040805160a08101825260025480825260035467ffffffffffffffff81166020840181905261ffff6801000000000000000083041684860181905263ffffffff6a010000000000000000000084048116606087018190526e010000000000000000000000000000909404166080860181905260085496517f5d3b1d3000000000000000000000000000000000000000000000000000000000815260048101959095526024850192909252604484015260648301919091526084820152909160009173ffffffffffffffffffffffffffffffffffffffff90911690635d3b1d309060a401602060405180830381600087803b1580156102f457600080fd5b505af1158015610308573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061032c9190610f44565b604080516060808201835288825260ff88811660208085019182526fffffffffffffffffffffffffffffffff8a8116868801908152600089815260048452889020875180825585516001909201805484518616610100027fffffffffffffffffffffffffffffff0000000000000000000000000000000000909116938916939093179290921790915588519081529351909416918301919091529151909116938101939093529293509183917f1577f444780fa640e5541a52439f4d0ce916960cf3738f33b35316c6399621b491015b60405180910390a2505050505050565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146104be576040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660248201526044015b60405180910390fd5b6104c8828261087e565b5050565b6104d4610aeb565b6002949094556003805467ffffffffffffffff949094167fffffffffffffffffffffffffffffffffffffffffffff00000000000000000000909416939093176801000000000000000061ffff9390931692909202919091177fffffffffffffffffffffffffffff0000000000000000ffffffffffffffffffff166a010000000000000000000063ffffffff928316027fffffffffffffffffffffffffffff00000000ffffffffffffffffffffffffffff16176e0100000000000000000000000000009190931602919091179055565b60015473ffffffffffffffffffffffffffffffffffffffff163314610624576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016104b5565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6060815160231461070d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f6d757374206265206f66206c656e67746820333500000000000000000000000060448201526064016104b5565b6040805160238082526104808201909252600091602082016104608036833701905050905060005b815181101561077f5761074981600161107e565b82828151811061075b5761075b611150565b60ff9092166020928302919091019091015280610777816110ad565b915050610735565b5060005b8151811015610863576000610799826023611096565b8583815181106107ab576107ab611150565b60200260200101516107bd91906110e6565b6107c7908361107e565b905060008382815181106107dd576107dd611150565b602002602001015190508383815181106107f9576107f9611150565b602002602001015184838151811061081357610813611150565b602002602001019060ff16908160ff16815250508084848151811061083a5761083a611150565b602002602001019060ff16908160ff16815250505050808061085b906110ad565b915050610783565b5092915050565b610872610aeb565b61087b81610b6e565b50565b60008281526004602090815260409182902082516060810184528154815260019091015460ff811692820183905261010090046fffffffffffffffffffffffffffffffff1692810192909252610930576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f7265717565737420756e7265636f676e697a656400000000000000000000000060448201526064016104b5565b600061093b836106a0565b60408051600580825260c08201909252919250600091906020820160a08036833701905050905060005b81518110156109be5782818151811061098057610980611150565b602002602001015182828151811061099a5761099a611150565b60ff90921660209283029190910190910152806109b6816110ad565b915050610965565b5060408051608081018252600080825260208083018281528385018381526060808601908152895186528984015160ff90811684528a8801516fffffffffffffffffffffffffffffffff90811684528983528d875260058652979095208651815592516001840180549351909816610100027fffffffffffffffffffffffffffffff00000000000000000000000000000000009093169516949094171790945590518051929384939092610a79926002850192910190610c64565b505050600086815260046020526040808220918255600190910180547fffffffffffffffffffffffffffffff00000000000000000000000000000000001690555186907fc7fd21de2615469409ad07bcbf39e226e8a48ccb612c074e2b758f81b8689594906103fc9087908590610ff5565b60005473ffffffffffffffffffffffffffffffffffffffff163314610b6c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016104b5565b565b73ffffffffffffffffffffffffffffffffffffffff8116331415610bee576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016104b5565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b82805482825590600052602060002090601f01602090048101928215610cfa5791602002820160005b83821115610ccb57835183826101000a81548160ff021916908360ff1602179055509260200192600101602081600001049283019260010302610c8d565b8015610cf85782816101000a81549060ff0219169055600101602081600001049283019260010302610ccb565b505b50610d06929150610d0a565b5090565b5b80821115610d065760008155600101610d0b565b600082601f830112610d3057600080fd5b8135602067ffffffffffffffff80831115610d4d57610d4d61117f565b8260051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f83011681018181108482111715610d9057610d9061117f565b60405284815283810192508684018288018501891015610daf57600080fd5b600092505b85831015610dd2578035845292840192600192909201918401610db4565b50979650505050505050565b803563ffffffff81168114610df257600080fd5b919050565b600060208284031215610e0957600080fd5b813573ffffffffffffffffffffffffffffffffffffffff81168114610e2d57600080fd5b9392505050565b600060208284031215610e4657600080fd5b813567ffffffffffffffff811115610e5d57600080fd5b610e6984828501610d1f565b949350505050565b600080600080600060a08688031215610e8957600080fd5b85359450602086013567ffffffffffffffff81168114610ea857600080fd5b9350604086013561ffff81168114610ebf57600080fd5b9250610ecd60608701610dde565b9150610edb60808701610dde565b90509295509295909350565b600080600060608486031215610efc57600080fd5b83359250602084013560ff81168114610f1457600080fd5b915060408401356fffffffffffffffffffffffffffffffff81168114610f3957600080fd5b809150509250925092565b600060208284031215610f5657600080fd5b5051919050565b60008060408385031215610f7057600080fd5b82359150602083013567ffffffffffffffff811115610f8e57600080fd5b610f9a85828601610d1f565b9150509250929050565b600081518084526020808501945080840160005b83811015610fd757815160ff1687529582019590820190600101610fb8565b509495945050505050565b602081526000610e2d6020830184610fa4565b8251815260208084015160ff16908201526040808401516fffffffffffffffffffffffffffffffff1690820152608060608201528151608082015260ff60208301511660a08201526fffffffffffffffffffffffffffffffff60408301511660c082015260006060830151608060e0840152611075610100840182610fa4565b95945050505050565b6000821982111561109157611091611121565b500190565b6000828210156110a8576110a8611121565b500390565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156110df576110df611121565b5060010190565b60008261111c577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500690565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var LotteryConsumerABI = LotteryConsumerMetaData.ABI

var LotteryConsumerBin = LotteryConsumerMetaData.Bin

func DeployLotteryConsumer(auth *bind.TransactOpts, backend bind.ContractBackend, _vrfCoordinator common.Address) (common.Address, *types.Transaction, *LotteryConsumer, error) {
	parsed, err := LotteryConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LotteryConsumerBin), backend, _vrfCoordinator)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LotteryConsumer{LotteryConsumerCaller: LotteryConsumerCaller{contract: contract}, LotteryConsumerTransactor: LotteryConsumerTransactor{contract: contract}, LotteryConsumerFilterer: LotteryConsumerFilterer{contract: contract}}, nil
}

type LotteryConsumer struct {
	address common.Address
	abi     abi.ABI
	LotteryConsumerCaller
	LotteryConsumerTransactor
	LotteryConsumerFilterer
}

type LotteryConsumerCaller struct {
	contract *bind.BoundContract
}

type LotteryConsumerTransactor struct {
	contract *bind.BoundContract
}

type LotteryConsumerFilterer struct {
	contract *bind.BoundContract
}

type LotteryConsumerSession struct {
	Contract     *LotteryConsumer
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type LotteryConsumerCallerSession struct {
	Contract *LotteryConsumerCaller
	CallOpts bind.CallOpts
}

type LotteryConsumerTransactorSession struct {
	Contract     *LotteryConsumerTransactor
	TransactOpts bind.TransactOpts
}

type LotteryConsumerRaw struct {
	Contract *LotteryConsumer
}

type LotteryConsumerCallerRaw struct {
	Contract *LotteryConsumerCaller
}

type LotteryConsumerTransactorRaw struct {
	Contract *LotteryConsumerTransactor
}

func NewLotteryConsumer(address common.Address, backend bind.ContractBackend) (*LotteryConsumer, error) {
	abi, err := abi.JSON(strings.NewReader(LotteryConsumerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindLotteryConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LotteryConsumer{address: address, abi: abi, LotteryConsumerCaller: LotteryConsumerCaller{contract: contract}, LotteryConsumerTransactor: LotteryConsumerTransactor{contract: contract}, LotteryConsumerFilterer: LotteryConsumerFilterer{contract: contract}}, nil
}

func NewLotteryConsumerCaller(address common.Address, caller bind.ContractCaller) (*LotteryConsumerCaller, error) {
	contract, err := bindLotteryConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LotteryConsumerCaller{contract: contract}, nil
}

func NewLotteryConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*LotteryConsumerTransactor, error) {
	contract, err := bindLotteryConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LotteryConsumerTransactor{contract: contract}, nil
}

func NewLotteryConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*LotteryConsumerFilterer, error) {
	contract, err := bindLotteryConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LotteryConsumerFilterer{contract: contract}, nil
}

func bindLotteryConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LotteryConsumerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_LotteryConsumer *LotteryConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LotteryConsumer.Contract.LotteryConsumerCaller.contract.Call(opts, result, method, params...)
}

func (_LotteryConsumer *LotteryConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LotteryConsumer.Contract.LotteryConsumerTransactor.contract.Transfer(opts)
}

func (_LotteryConsumer *LotteryConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LotteryConsumer.Contract.LotteryConsumerTransactor.contract.Transact(opts, method, params...)
}

func (_LotteryConsumer *LotteryConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LotteryConsumer.Contract.contract.Call(opts, result, method, params...)
}

func (_LotteryConsumer *LotteryConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LotteryConsumer.Contract.contract.Transfer(opts)
}

func (_LotteryConsumer *LotteryConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LotteryConsumer.Contract.contract.Transact(opts, method, params...)
}

func (_LotteryConsumer *LotteryConsumerCaller) GetMostRecentVrfRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LotteryConsumer.contract.Call(opts, &out, "getMostRecentVrfRequestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LotteryConsumer *LotteryConsumerSession) GetMostRecentVrfRequestId() (*big.Int, error) {
	return _LotteryConsumer.Contract.GetMostRecentVrfRequestId(&_LotteryConsumer.CallOpts)
}

func (_LotteryConsumer *LotteryConsumerCallerSession) GetMostRecentVrfRequestId() (*big.Int, error) {
	return _LotteryConsumer.Contract.GetMostRecentVrfRequestId(&_LotteryConsumer.CallOpts)
}

func (_LotteryConsumer *LotteryConsumerCaller) GetRequestConfig(opts *bind.CallOpts) (GetRequestConfig,

	error) {
	var out []interface{}
	err := _LotteryConsumer.contract.Call(opts, &out, "getRequestConfig")

	outstruct := new(GetRequestConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.KeyHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.SubscriptionId = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.MinRequestConfirmations = *abi.ConvertType(out[2], new(uint16)).(*uint16)
	outstruct.CallbackGasLimit = *abi.ConvertType(out[3], new(uint32)).(*uint32)
	outstruct.NumWords = *abi.ConvertType(out[4], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_LotteryConsumer *LotteryConsumerSession) GetRequestConfig() (GetRequestConfig,

	error) {
	return _LotteryConsumer.Contract.GetRequestConfig(&_LotteryConsumer.CallOpts)
}

func (_LotteryConsumer *LotteryConsumerCallerSession) GetRequestConfig() (GetRequestConfig,

	error) {
	return _LotteryConsumer.Contract.GetRequestConfig(&_LotteryConsumer.CallOpts)
}

func (_LotteryConsumer *LotteryConsumerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LotteryConsumer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LotteryConsumer *LotteryConsumerSession) Owner() (common.Address, error) {
	return _LotteryConsumer.Contract.Owner(&_LotteryConsumer.CallOpts)
}

func (_LotteryConsumer *LotteryConsumerCallerSession) Owner() (common.Address, error) {
	return _LotteryConsumer.Contract.Owner(&_LotteryConsumer.CallOpts)
}

func (_LotteryConsumer *LotteryConsumerCaller) Shuffle35(opts *bind.CallOpts, randomness []*big.Int) ([]uint8, error) {
	var out []interface{}
	err := _LotteryConsumer.contract.Call(opts, &out, "shuffle35", randomness)

	if err != nil {
		return *new([]uint8), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint8)).(*[]uint8)

	return out0, err

}

func (_LotteryConsumer *LotteryConsumerSession) Shuffle35(randomness []*big.Int) ([]uint8, error) {
	return _LotteryConsumer.Contract.Shuffle35(&_LotteryConsumer.CallOpts, randomness)
}

func (_LotteryConsumer *LotteryConsumerCallerSession) Shuffle35(randomness []*big.Int) ([]uint8, error) {
	return _LotteryConsumer.Contract.Shuffle35(&_LotteryConsumer.CallOpts, randomness)
}

func (_LotteryConsumer *LotteryConsumerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LotteryConsumer.contract.Transact(opts, "acceptOwnership")
}

func (_LotteryConsumer *LotteryConsumerSession) AcceptOwnership() (*types.Transaction, error) {
	return _LotteryConsumer.Contract.AcceptOwnership(&_LotteryConsumer.TransactOpts)
}

func (_LotteryConsumer *LotteryConsumerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _LotteryConsumer.Contract.AcceptOwnership(&_LotteryConsumer.TransactOpts)
}

func (_LotteryConsumer *LotteryConsumerTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _LotteryConsumer.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_LotteryConsumer *LotteryConsumerSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _LotteryConsumer.Contract.RawFulfillRandomWords(&_LotteryConsumer.TransactOpts, requestId, randomWords)
}

func (_LotteryConsumer *LotteryConsumerTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _LotteryConsumer.Contract.RawFulfillRandomWords(&_LotteryConsumer.TransactOpts, requestId, randomWords)
}

func (_LotteryConsumer *LotteryConsumerTransactor) RequestRandomness(opts *bind.TransactOpts, clientRequestId [32]byte, lotteryType uint8, vrfExternalRequestId *big.Int) (*types.Transaction, error) {
	return _LotteryConsumer.contract.Transact(opts, "requestRandomness", clientRequestId, lotteryType, vrfExternalRequestId)
}

func (_LotteryConsumer *LotteryConsumerSession) RequestRandomness(clientRequestId [32]byte, lotteryType uint8, vrfExternalRequestId *big.Int) (*types.Transaction, error) {
	return _LotteryConsumer.Contract.RequestRandomness(&_LotteryConsumer.TransactOpts, clientRequestId, lotteryType, vrfExternalRequestId)
}

func (_LotteryConsumer *LotteryConsumerTransactorSession) RequestRandomness(clientRequestId [32]byte, lotteryType uint8, vrfExternalRequestId *big.Int) (*types.Transaction, error) {
	return _LotteryConsumer.Contract.RequestRandomness(&_LotteryConsumer.TransactOpts, clientRequestId, lotteryType, vrfExternalRequestId)
}

func (_LotteryConsumer *LotteryConsumerTransactor) SetRequestConfig(opts *bind.TransactOpts, keyHash [32]byte, subscriptionId uint64, minRequestConfirmations uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _LotteryConsumer.contract.Transact(opts, "setRequestConfig", keyHash, subscriptionId, minRequestConfirmations, callbackGasLimit, numWords)
}

func (_LotteryConsumer *LotteryConsumerSession) SetRequestConfig(keyHash [32]byte, subscriptionId uint64, minRequestConfirmations uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _LotteryConsumer.Contract.SetRequestConfig(&_LotteryConsumer.TransactOpts, keyHash, subscriptionId, minRequestConfirmations, callbackGasLimit, numWords)
}

func (_LotteryConsumer *LotteryConsumerTransactorSession) SetRequestConfig(keyHash [32]byte, subscriptionId uint64, minRequestConfirmations uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _LotteryConsumer.Contract.SetRequestConfig(&_LotteryConsumer.TransactOpts, keyHash, subscriptionId, minRequestConfirmations, callbackGasLimit, numWords)
}

func (_LotteryConsumer *LotteryConsumerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _LotteryConsumer.contract.Transact(opts, "transferOwnership", to)
}

func (_LotteryConsumer *LotteryConsumerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _LotteryConsumer.Contract.TransferOwnership(&_LotteryConsumer.TransactOpts, to)
}

func (_LotteryConsumer *LotteryConsumerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _LotteryConsumer.Contract.TransferOwnership(&_LotteryConsumer.TransactOpts, to)
}

type LotteryConsumerAllowedCallerAddedIterator struct {
	Event *LotteryConsumerAllowedCallerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LotteryConsumerAllowedCallerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LotteryConsumerAllowedCallerAdded)
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
		it.Event = new(LotteryConsumerAllowedCallerAdded)
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

func (it *LotteryConsumerAllowedCallerAddedIterator) Error() error {
	return it.fail
}

func (it *LotteryConsumerAllowedCallerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LotteryConsumerAllowedCallerAdded struct {
	Caller common.Address
	Raw    types.Log
}

func (_LotteryConsumer *LotteryConsumerFilterer) FilterAllowedCallerAdded(opts *bind.FilterOpts) (*LotteryConsumerAllowedCallerAddedIterator, error) {

	logs, sub, err := _LotteryConsumer.contract.FilterLogs(opts, "AllowedCallerAdded")
	if err != nil {
		return nil, err
	}
	return &LotteryConsumerAllowedCallerAddedIterator{contract: _LotteryConsumer.contract, event: "AllowedCallerAdded", logs: logs, sub: sub}, nil
}

func (_LotteryConsumer *LotteryConsumerFilterer) WatchAllowedCallerAdded(opts *bind.WatchOpts, sink chan<- *LotteryConsumerAllowedCallerAdded) (event.Subscription, error) {

	logs, sub, err := _LotteryConsumer.contract.WatchLogs(opts, "AllowedCallerAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LotteryConsumerAllowedCallerAdded)
				if err := _LotteryConsumer.contract.UnpackLog(event, "AllowedCallerAdded", log); err != nil {
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

func (_LotteryConsumer *LotteryConsumerFilterer) ParseAllowedCallerAdded(log types.Log) (*LotteryConsumerAllowedCallerAdded, error) {
	event := new(LotteryConsumerAllowedCallerAdded)
	if err := _LotteryConsumer.contract.UnpackLog(event, "AllowedCallerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LotteryConsumerAllowedCallerRemovedIterator struct {
	Event *LotteryConsumerAllowedCallerRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LotteryConsumerAllowedCallerRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LotteryConsumerAllowedCallerRemoved)
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
		it.Event = new(LotteryConsumerAllowedCallerRemoved)
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

func (it *LotteryConsumerAllowedCallerRemovedIterator) Error() error {
	return it.fail
}

func (it *LotteryConsumerAllowedCallerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LotteryConsumerAllowedCallerRemoved struct {
	Caller common.Address
	Raw    types.Log
}

func (_LotteryConsumer *LotteryConsumerFilterer) FilterAllowedCallerRemoved(opts *bind.FilterOpts) (*LotteryConsumerAllowedCallerRemovedIterator, error) {

	logs, sub, err := _LotteryConsumer.contract.FilterLogs(opts, "AllowedCallerRemoved")
	if err != nil {
		return nil, err
	}
	return &LotteryConsumerAllowedCallerRemovedIterator{contract: _LotteryConsumer.contract, event: "AllowedCallerRemoved", logs: logs, sub: sub}, nil
}

func (_LotteryConsumer *LotteryConsumerFilterer) WatchAllowedCallerRemoved(opts *bind.WatchOpts, sink chan<- *LotteryConsumerAllowedCallerRemoved) (event.Subscription, error) {

	logs, sub, err := _LotteryConsumer.contract.WatchLogs(opts, "AllowedCallerRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LotteryConsumerAllowedCallerRemoved)
				if err := _LotteryConsumer.contract.UnpackLog(event, "AllowedCallerRemoved", log); err != nil {
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

func (_LotteryConsumer *LotteryConsumerFilterer) ParseAllowedCallerRemoved(log types.Log) (*LotteryConsumerAllowedCallerRemoved, error) {
	event := new(LotteryConsumerAllowedCallerRemoved)
	if err := _LotteryConsumer.contract.UnpackLog(event, "AllowedCallerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LotteryConsumerLotterySettledIterator struct {
	Event *LotteryConsumerLotterySettled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LotteryConsumerLotterySettledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LotteryConsumerLotterySettled)
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
		it.Event = new(LotteryConsumerLotterySettled)
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

func (it *LotteryConsumerLotterySettledIterator) Error() error {
	return it.fail
}

func (it *LotteryConsumerLotterySettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LotteryConsumerLotterySettled struct {
	VrfRequestId *big.Int
	Request      LotteryConsumerLotteryRequest
	Outcome      LotteryConsumerLotteryOutcome
	Raw          types.Log
}

func (_LotteryConsumer *LotteryConsumerFilterer) FilterLotterySettled(opts *bind.FilterOpts, vrfRequestId []*big.Int) (*LotteryConsumerLotterySettledIterator, error) {

	var vrfRequestIdRule []interface{}
	for _, vrfRequestIdItem := range vrfRequestId {
		vrfRequestIdRule = append(vrfRequestIdRule, vrfRequestIdItem)
	}

	logs, sub, err := _LotteryConsumer.contract.FilterLogs(opts, "LotterySettled", vrfRequestIdRule)
	if err != nil {
		return nil, err
	}
	return &LotteryConsumerLotterySettledIterator{contract: _LotteryConsumer.contract, event: "LotterySettled", logs: logs, sub: sub}, nil
}

func (_LotteryConsumer *LotteryConsumerFilterer) WatchLotterySettled(opts *bind.WatchOpts, sink chan<- *LotteryConsumerLotterySettled, vrfRequestId []*big.Int) (event.Subscription, error) {

	var vrfRequestIdRule []interface{}
	for _, vrfRequestIdItem := range vrfRequestId {
		vrfRequestIdRule = append(vrfRequestIdRule, vrfRequestIdItem)
	}

	logs, sub, err := _LotteryConsumer.contract.WatchLogs(opts, "LotterySettled", vrfRequestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LotteryConsumerLotterySettled)
				if err := _LotteryConsumer.contract.UnpackLog(event, "LotterySettled", log); err != nil {
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

func (_LotteryConsumer *LotteryConsumerFilterer) ParseLotterySettled(log types.Log) (*LotteryConsumerLotterySettled, error) {
	event := new(LotteryConsumerLotterySettled)
	if err := _LotteryConsumer.contract.UnpackLog(event, "LotterySettled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LotteryConsumerLotteryStartedIterator struct {
	Event *LotteryConsumerLotteryStarted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LotteryConsumerLotteryStartedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LotteryConsumerLotteryStarted)
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
		it.Event = new(LotteryConsumerLotteryStarted)
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

func (it *LotteryConsumerLotteryStartedIterator) Error() error {
	return it.fail
}

func (it *LotteryConsumerLotteryStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LotteryConsumerLotteryStarted struct {
	VrfRequestId *big.Int
	Request      LotteryConsumerLotteryRequest
	Raw          types.Log
}

func (_LotteryConsumer *LotteryConsumerFilterer) FilterLotteryStarted(opts *bind.FilterOpts, vrfRequestId []*big.Int) (*LotteryConsumerLotteryStartedIterator, error) {

	var vrfRequestIdRule []interface{}
	for _, vrfRequestIdItem := range vrfRequestId {
		vrfRequestIdRule = append(vrfRequestIdRule, vrfRequestIdItem)
	}

	logs, sub, err := _LotteryConsumer.contract.FilterLogs(opts, "LotteryStarted", vrfRequestIdRule)
	if err != nil {
		return nil, err
	}
	return &LotteryConsumerLotteryStartedIterator{contract: _LotteryConsumer.contract, event: "LotteryStarted", logs: logs, sub: sub}, nil
}

func (_LotteryConsumer *LotteryConsumerFilterer) WatchLotteryStarted(opts *bind.WatchOpts, sink chan<- *LotteryConsumerLotteryStarted, vrfRequestId []*big.Int) (event.Subscription, error) {

	var vrfRequestIdRule []interface{}
	for _, vrfRequestIdItem := range vrfRequestId {
		vrfRequestIdRule = append(vrfRequestIdRule, vrfRequestIdItem)
	}

	logs, sub, err := _LotteryConsumer.contract.WatchLogs(opts, "LotteryStarted", vrfRequestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LotteryConsumerLotteryStarted)
				if err := _LotteryConsumer.contract.UnpackLog(event, "LotteryStarted", log); err != nil {
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

func (_LotteryConsumer *LotteryConsumerFilterer) ParseLotteryStarted(log types.Log) (*LotteryConsumerLotteryStarted, error) {
	event := new(LotteryConsumerLotteryStarted)
	if err := _LotteryConsumer.contract.UnpackLog(event, "LotteryStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LotteryConsumerOwnershipTransferRequestedIterator struct {
	Event *LotteryConsumerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LotteryConsumerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LotteryConsumerOwnershipTransferRequested)
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
		it.Event = new(LotteryConsumerOwnershipTransferRequested)
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

func (it *LotteryConsumerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *LotteryConsumerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LotteryConsumerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_LotteryConsumer *LotteryConsumerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LotteryConsumerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LotteryConsumer.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LotteryConsumerOwnershipTransferRequestedIterator{contract: _LotteryConsumer.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_LotteryConsumer *LotteryConsumerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *LotteryConsumerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LotteryConsumer.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LotteryConsumerOwnershipTransferRequested)
				if err := _LotteryConsumer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_LotteryConsumer *LotteryConsumerFilterer) ParseOwnershipTransferRequested(log types.Log) (*LotteryConsumerOwnershipTransferRequested, error) {
	event := new(LotteryConsumerOwnershipTransferRequested)
	if err := _LotteryConsumer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LotteryConsumerOwnershipTransferredIterator struct {
	Event *LotteryConsumerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LotteryConsumerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LotteryConsumerOwnershipTransferred)
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
		it.Event = new(LotteryConsumerOwnershipTransferred)
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

func (it *LotteryConsumerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *LotteryConsumerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LotteryConsumerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_LotteryConsumer *LotteryConsumerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LotteryConsumerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LotteryConsumer.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LotteryConsumerOwnershipTransferredIterator{contract: _LotteryConsumer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_LotteryConsumer *LotteryConsumerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *LotteryConsumerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LotteryConsumer.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LotteryConsumerOwnershipTransferred)
				if err := _LotteryConsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_LotteryConsumer *LotteryConsumerFilterer) ParseOwnershipTransferred(log types.Log) (*LotteryConsumerOwnershipTransferred, error) {
	event := new(LotteryConsumerOwnershipTransferred)
	if err := _LotteryConsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetRequestConfig struct {
	KeyHash                 [32]byte
	SubscriptionId          uint64
	MinRequestConfirmations uint16
	CallbackGasLimit        uint32
	NumWords                uint32
}

func (_LotteryConsumer *LotteryConsumer) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _LotteryConsumer.abi.Events["AllowedCallerAdded"].ID:
		return _LotteryConsumer.ParseAllowedCallerAdded(log)
	case _LotteryConsumer.abi.Events["AllowedCallerRemoved"].ID:
		return _LotteryConsumer.ParseAllowedCallerRemoved(log)
	case _LotteryConsumer.abi.Events["LotterySettled"].ID:
		return _LotteryConsumer.ParseLotterySettled(log)
	case _LotteryConsumer.abi.Events["LotteryStarted"].ID:
		return _LotteryConsumer.ParseLotteryStarted(log)
	case _LotteryConsumer.abi.Events["OwnershipTransferRequested"].ID:
		return _LotteryConsumer.ParseOwnershipTransferRequested(log)
	case _LotteryConsumer.abi.Events["OwnershipTransferred"].ID:
		return _LotteryConsumer.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (LotteryConsumerAllowedCallerAdded) Topic() common.Hash {
	return common.HexToHash("0x663c7e9ed36d9138863ef4306bbfcf01f60e1e7ca69b370c53d3094369e2cb02")
}

func (LotteryConsumerAllowedCallerRemoved) Topic() common.Hash {
	return common.HexToHash("0xbc0a6e072a312bde289d32bc84e5b758d7c617f734ecc0d69f995b2d7e69be36")
}

func (LotteryConsumerLotterySettled) Topic() common.Hash {
	return common.HexToHash("0xc7fd21de2615469409ad07bcbf39e226e8a48ccb612c074e2b758f81b8689594")
}

func (LotteryConsumerLotteryStarted) Topic() common.Hash {
	return common.HexToHash("0x1577f444780fa640e5541a52439f4d0ce916960cf3738f33b35316c6399621b4")
}

func (LotteryConsumerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (LotteryConsumerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_LotteryConsumer *LotteryConsumer) Address() common.Address {
	return _LotteryConsumer.address
}

type LotteryConsumerInterface interface {
	GetMostRecentVrfRequestId(opts *bind.CallOpts) (*big.Int, error)

	GetRequestConfig(opts *bind.CallOpts) (GetRequestConfig,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	Shuffle35(opts *bind.CallOpts, randomness []*big.Int) ([]uint8, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomness(opts *bind.TransactOpts, clientRequestId [32]byte, lotteryType uint8, vrfExternalRequestId *big.Int) (*types.Transaction, error)

	SetRequestConfig(opts *bind.TransactOpts, keyHash [32]byte, subscriptionId uint64, minRequestConfirmations uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterAllowedCallerAdded(opts *bind.FilterOpts) (*LotteryConsumerAllowedCallerAddedIterator, error)

	WatchAllowedCallerAdded(opts *bind.WatchOpts, sink chan<- *LotteryConsumerAllowedCallerAdded) (event.Subscription, error)

	ParseAllowedCallerAdded(log types.Log) (*LotteryConsumerAllowedCallerAdded, error)

	FilterAllowedCallerRemoved(opts *bind.FilterOpts) (*LotteryConsumerAllowedCallerRemovedIterator, error)

	WatchAllowedCallerRemoved(opts *bind.WatchOpts, sink chan<- *LotteryConsumerAllowedCallerRemoved) (event.Subscription, error)

	ParseAllowedCallerRemoved(log types.Log) (*LotteryConsumerAllowedCallerRemoved, error)

	FilterLotterySettled(opts *bind.FilterOpts, vrfRequestId []*big.Int) (*LotteryConsumerLotterySettledIterator, error)

	WatchLotterySettled(opts *bind.WatchOpts, sink chan<- *LotteryConsumerLotterySettled, vrfRequestId []*big.Int) (event.Subscription, error)

	ParseLotterySettled(log types.Log) (*LotteryConsumerLotterySettled, error)

	FilterLotteryStarted(opts *bind.FilterOpts, vrfRequestId []*big.Int) (*LotteryConsumerLotteryStartedIterator, error)

	WatchLotteryStarted(opts *bind.WatchOpts, sink chan<- *LotteryConsumerLotteryStarted, vrfRequestId []*big.Int) (event.Subscription, error)

	ParseLotteryStarted(log types.Log) (*LotteryConsumerLotteryStarted, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LotteryConsumerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *LotteryConsumerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*LotteryConsumerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LotteryConsumerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *LotteryConsumerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*LotteryConsumerOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
