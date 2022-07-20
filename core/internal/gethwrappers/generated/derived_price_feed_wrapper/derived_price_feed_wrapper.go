// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package derived_price_feed_wrapper

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
)

var DerivedPriceFeedMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_base\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_quote\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BASE\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DECIMALS\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"QUOTE\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"\",\"type\":\"uint80\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60e060405234801561001057600080fd5b50604051610c14380380610c1483398101604081905261002f916100ec565b60ff8116158015906100455750601260ff821611155b6100895760405162461bcd60e51b8152602060048201526011602482015270496e76616c6964205f646563696d616c7360781b604482015260640160405180910390fd5b60f81b7fff000000000000000000000000000000000000000000000000000000000000001660c052606091821b6001600160601b0319908116608052911b1660a052610139565b80516001600160a01b03811681146100e757600080fd5b919050565b60008060006060848603121561010157600080fd5b61010a846100d0565b9250610118602085016100d0565b9150604084015160ff8116811461012e57600080fd5b809150509250925092565b60805160601c60a05160601c60c05160f81c610a6d6101a76000396000818160920152818160cd0152818161041a0152818161058f01526105bd0152600081816101950152818161044401526104ea0152600081816101e1015281816102cf01526103750152610a6d6000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c80639a6fc8f51161005b5780639a6fc8f5146101465780639c57983914610190578063ec342ad0146101dc578063feaf968c1461020357600080fd5b80632e0f26251461008d578063313ce567146100cb57806354fd4d50146100f15780637284e41614610107575b600080fd5b6100b47f000000000000000000000000000000000000000000000000000000000000000081565b60405160ff90911681526020015b60405180910390f35b7f00000000000000000000000000000000000000000000000000000000000000006100b4565b6100f9600081565b6040519081526020016100c2565b604080518082018252601481527f446572697665645072696365466565642e736f6c000000000000000000000000602082015290516100c2919061070c565b610159610154366004610674565b61020b565b6040805169ffffffffffffffffffff968716815260208101959095528401929092526060830152909116608082015260a0016100c2565b6101b77f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100c2565b6101b77f000000000000000000000000000000000000000000000000000000000000000081565b6101596102a6565b60008060008060006040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161029d9060208082526027908201527f6e6f7420696d706c656d656e746564202d20757365206c6174657374526f756e60408201527f6444617461282900000000000000000000000000000000000000000000000000606082015260800190565b60405180910390fd5b6000806000806000806102b76102ca565b9096909550429450849350600092509050565b6000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b15801561033357600080fd5b505afa158015610347573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061036b9190610691565b50505091505060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663313ce5676040518163ffffffff1660e01b815260040160206040518083038186803b1580156103d957600080fd5b505afa1580156103ed573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061041191906106e9565b905061043e82827f0000000000000000000000000000000000000000000000000000000000000000610601565b915060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b1580156104a857600080fd5b505afa1580156104bc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104e09190610691565b50505091505060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663313ce5676040518163ffffffff1660e01b815260040160206040518083038186803b15801561054e57600080fd5b505afa158015610562573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061058691906106e9565b90506105b382827f0000000000000000000000000000000000000000000000000000000000000000610601565b9150816105e460ff7f000000000000000000000000000000000000000000000000000000000000000016600a61086f565b6105ee9086610937565b6105f8919061077f565b94505050505090565b60008160ff168360ff16101561063a5761061b83836109f3565b6106299060ff16600a61086f565b6106339085610937565b905061066d565b8160ff168360ff16111561066a5761065282846109f3565b6106609060ff16600a61086f565b610633908561077f565b50825b9392505050565b60006020828403121561068657600080fd5b813561066d81610a45565b600080600080600060a086880312156106a957600080fd5b85516106b481610a45565b8095505060208601519350604086015192506060860151915060808601516106db81610a45565b809150509295509295909350565b6000602082840312156106fb57600080fd5b815160ff8116811461066d57600080fd5b600060208083528351808285015260005b818110156107395785810183015185820160400152820161071d565b8181111561074b576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b6000826107b5577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83147f80000000000000000000000000000000000000000000000000000000000000008314161561080957610809610a16565b500590565b600181815b8085111561086757817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0482111561084d5761084d610a16565b8085161561085a57918102915b93841c9390800290610813565b509250929050565b600061066d838360008261088557506001610931565b8161089257506000610931565b81600181146108a857600281146108b2576108ce565b6001915050610931565b60ff8411156108c3576108c3610a16565b50506001821b610931565b5060208310610133831016604e8410600b84101617156108f1575081810a610931565b6108fb838361080e565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0482111561092d5761092d610a16565b0290505b92915050565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60008413600084138583048511828216161561097857610978610a16565b7f800000000000000000000000000000000000000000000000000000000000000060008712868205881281841616156109b3576109b3610a16565b600087129250878205871284841616156109cf576109cf610a16565b878505871281841616156109e5576109e5610a16565b505050929093029392505050565b600060ff821660ff841680821015610a0d57610a0d610a16565b90039392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b69ffffffffffffffffffff81168114610a5d57600080fd5b5056fea164736f6c6343000806000a",
}

var DerivedPriceFeedABI = DerivedPriceFeedMetaData.ABI

var DerivedPriceFeedBin = DerivedPriceFeedMetaData.Bin

func DeployDerivedPriceFeed(auth *bind.TransactOpts, backend bind.ContractBackend, _base common.Address, _quote common.Address, _decimals uint8) (common.Address, *types.Transaction, *DerivedPriceFeed, error) {
	parsed, err := DerivedPriceFeedMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DerivedPriceFeedBin), backend, _base, _quote, _decimals)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DerivedPriceFeed{DerivedPriceFeedCaller: DerivedPriceFeedCaller{contract: contract}, DerivedPriceFeedTransactor: DerivedPriceFeedTransactor{contract: contract}, DerivedPriceFeedFilterer: DerivedPriceFeedFilterer{contract: contract}}, nil
}

type DerivedPriceFeed struct {
	address common.Address
	abi     abi.ABI
	DerivedPriceFeedCaller
	DerivedPriceFeedTransactor
	DerivedPriceFeedFilterer
}

type DerivedPriceFeedCaller struct {
	contract *bind.BoundContract
}

type DerivedPriceFeedTransactor struct {
	contract *bind.BoundContract
}

type DerivedPriceFeedFilterer struct {
	contract *bind.BoundContract
}

type DerivedPriceFeedSession struct {
	Contract     *DerivedPriceFeed
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type DerivedPriceFeedCallerSession struct {
	Contract *DerivedPriceFeedCaller
	CallOpts bind.CallOpts
}

type DerivedPriceFeedTransactorSession struct {
	Contract     *DerivedPriceFeedTransactor
	TransactOpts bind.TransactOpts
}

type DerivedPriceFeedRaw struct {
	Contract *DerivedPriceFeed
}

type DerivedPriceFeedCallerRaw struct {
	Contract *DerivedPriceFeedCaller
}

type DerivedPriceFeedTransactorRaw struct {
	Contract *DerivedPriceFeedTransactor
}

func NewDerivedPriceFeed(address common.Address, backend bind.ContractBackend) (*DerivedPriceFeed, error) {
	abi, err := abi.JSON(strings.NewReader(DerivedPriceFeedABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindDerivedPriceFeed(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DerivedPriceFeed{address: address, abi: abi, DerivedPriceFeedCaller: DerivedPriceFeedCaller{contract: contract}, DerivedPriceFeedTransactor: DerivedPriceFeedTransactor{contract: contract}, DerivedPriceFeedFilterer: DerivedPriceFeedFilterer{contract: contract}}, nil
}

func NewDerivedPriceFeedCaller(address common.Address, caller bind.ContractCaller) (*DerivedPriceFeedCaller, error) {
	contract, err := bindDerivedPriceFeed(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DerivedPriceFeedCaller{contract: contract}, nil
}

func NewDerivedPriceFeedTransactor(address common.Address, transactor bind.ContractTransactor) (*DerivedPriceFeedTransactor, error) {
	contract, err := bindDerivedPriceFeed(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DerivedPriceFeedTransactor{contract: contract}, nil
}

func NewDerivedPriceFeedFilterer(address common.Address, filterer bind.ContractFilterer) (*DerivedPriceFeedFilterer, error) {
	contract, err := bindDerivedPriceFeed(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DerivedPriceFeedFilterer{contract: contract}, nil
}

func bindDerivedPriceFeed(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DerivedPriceFeedABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_DerivedPriceFeed *DerivedPriceFeedRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DerivedPriceFeed.Contract.DerivedPriceFeedCaller.contract.Call(opts, result, method, params...)
}

func (_DerivedPriceFeed *DerivedPriceFeedRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DerivedPriceFeed.Contract.DerivedPriceFeedTransactor.contract.Transfer(opts)
}

func (_DerivedPriceFeed *DerivedPriceFeedRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DerivedPriceFeed.Contract.DerivedPriceFeedTransactor.contract.Transact(opts, method, params...)
}

func (_DerivedPriceFeed *DerivedPriceFeedCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DerivedPriceFeed.Contract.contract.Call(opts, result, method, params...)
}

func (_DerivedPriceFeed *DerivedPriceFeedTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DerivedPriceFeed.Contract.contract.Transfer(opts)
}

func (_DerivedPriceFeed *DerivedPriceFeedTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DerivedPriceFeed.Contract.contract.Transact(opts, method, params...)
}

func (_DerivedPriceFeed *DerivedPriceFeedCaller) BASE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DerivedPriceFeed.contract.Call(opts, &out, "BASE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_DerivedPriceFeed *DerivedPriceFeedSession) BASE() (common.Address, error) {
	return _DerivedPriceFeed.Contract.BASE(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCallerSession) BASE() (common.Address, error) {
	return _DerivedPriceFeed.Contract.BASE(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCaller) DECIMALS(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _DerivedPriceFeed.contract.Call(opts, &out, "DECIMALS")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_DerivedPriceFeed *DerivedPriceFeedSession) DECIMALS() (uint8, error) {
	return _DerivedPriceFeed.Contract.DECIMALS(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCallerSession) DECIMALS() (uint8, error) {
	return _DerivedPriceFeed.Contract.DECIMALS(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCaller) QUOTE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DerivedPriceFeed.contract.Call(opts, &out, "QUOTE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_DerivedPriceFeed *DerivedPriceFeedSession) QUOTE() (common.Address, error) {
	return _DerivedPriceFeed.Contract.QUOTE(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCallerSession) QUOTE() (common.Address, error) {
	return _DerivedPriceFeed.Contract.QUOTE(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _DerivedPriceFeed.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_DerivedPriceFeed *DerivedPriceFeedSession) Decimals() (uint8, error) {
	return _DerivedPriceFeed.Contract.Decimals(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCallerSession) Decimals() (uint8, error) {
	return _DerivedPriceFeed.Contract.Decimals(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _DerivedPriceFeed.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_DerivedPriceFeed *DerivedPriceFeedSession) Description() (string, error) {
	return _DerivedPriceFeed.Contract.Description(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCallerSession) Description() (string, error) {
	return _DerivedPriceFeed.Contract.Description(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCaller) GetRoundData(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error) {
	var out []interface{}
	err := _DerivedPriceFeed.contract.Call(opts, &out, "getRoundData", arg0)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), *new(*big.Int), *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	out3 := *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	out4 := *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return out0, out1, out2, out3, out4, err

}

func (_DerivedPriceFeed *DerivedPriceFeedSession) GetRoundData(arg0 *big.Int) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error) {
	return _DerivedPriceFeed.Contract.GetRoundData(&_DerivedPriceFeed.CallOpts, arg0)
}

func (_DerivedPriceFeed *DerivedPriceFeedCallerSession) GetRoundData(arg0 *big.Int) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error) {
	return _DerivedPriceFeed.Contract.GetRoundData(&_DerivedPriceFeed.CallOpts, arg0)
}

func (_DerivedPriceFeed *DerivedPriceFeedCaller) LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

	error) {
	var out []interface{}
	err := _DerivedPriceFeed.contract.Call(opts, &out, "latestRoundData")

	outstruct := new(LatestRoundData)
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Answer = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_DerivedPriceFeed *DerivedPriceFeedSession) LatestRoundData() (LatestRoundData,

	error) {
	return _DerivedPriceFeed.Contract.LatestRoundData(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCallerSession) LatestRoundData() (LatestRoundData,

	error) {
	return _DerivedPriceFeed.Contract.LatestRoundData(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DerivedPriceFeed.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DerivedPriceFeed *DerivedPriceFeedSession) Version() (*big.Int, error) {
	return _DerivedPriceFeed.Contract.Version(&_DerivedPriceFeed.CallOpts)
}

func (_DerivedPriceFeed *DerivedPriceFeedCallerSession) Version() (*big.Int, error) {
	return _DerivedPriceFeed.Contract.Version(&_DerivedPriceFeed.CallOpts)
}

type LatestRoundData struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}

func (_DerivedPriceFeed *DerivedPriceFeed) Address() common.Address {
	return _DerivedPriceFeed.address
}

type DerivedPriceFeedInterface interface {
	BASE(opts *bind.CallOpts) (common.Address, error)

	DECIMALS(opts *bind.CallOpts) (uint8, error)

	QUOTE(opts *bind.CallOpts) (common.Address, error)

	Decimals(opts *bind.CallOpts) (uint8, error)

	Description(opts *bind.CallOpts) (string, error)

	GetRoundData(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error)

	LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

		error)

	Version(opts *bind.CallOpts) (*big.Int, error)

	Address() common.Address
}
