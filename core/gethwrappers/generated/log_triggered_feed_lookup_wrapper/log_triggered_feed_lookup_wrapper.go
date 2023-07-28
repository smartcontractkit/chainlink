// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package log_triggered_feed_lookup_wrapper

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

type Log struct {
	Index       *big.Int
	TxIndex     *big.Int
	TxHash      [32]byte
	BlockNumber *big.Int
	BlockHash   [32]byte
	Source      common.Address
	Topics      [][32]byte
	Data        []byte
}

var LogTriggeredFeedLookupMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_useArbitrumBlockNum\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"FeedLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"orderId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"exchange\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"blob\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"verified\",\"type\":\"bytes\"}],\"name\":\"PerformingLogTriggerUpkeep\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"txIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"topics\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLog\",\"name\":\"log\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"checkLog\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feedsHex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParam\",\"type\":\"string\"}],\"name\":\"setFeedParamKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"newFeeds\",\"type\":\"string[]\"}],\"name\":\"setFeedsHex\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"timeParam\",\"type\":\"string\"}],\"name\":\"setTimeParamKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbitrumBlockNum\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x610120604052604260a0818152608091829190620015dd60c03990526200002a9060019081620000d4565b506040805180820190915260098152680cccacac892c890caf60bb1b60208201526002906200005a908262000250565b5060408051808201909152600b81526a313637b1b5a73ab6b132b960a91b60208201526003906200008c908262000250565b503480156200009a57600080fd5b506040516200161f3803806200161f833981016040819052620000bd916200031c565b6000805460ff191691151591909117905562000347565b8280548282559060005260206000209081019282156200011f579160200282015b828111156200011f57825182906200010e908262000250565b5091602001919060010190620000f5565b506200012d92915062000131565b5090565b808211156200012d57600062000148828262000152565b5060010162000131565b5080546200016090620001c1565b6000825580601f1062000171575050565b601f01602090049060005260206000209081019062000191919062000194565b50565b5b808211156200012d576000815560010162000195565b634e487b7160e01b600052604160045260246000fd5b600181811c90821680620001d657607f821691505b602082108103620001f757634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200024b57600081815260208120601f850160051c81016020861015620002265750805b601f850160051c820191505b81811015620002475782815560010162000232565b5050505b505050565b81516001600160401b038111156200026c576200026c620001ab565b62000284816200027d8454620001c1565b84620001fd565b602080601f831160018114620002bc5760008415620002a35750858301515b600019600386901b1c1916600185901b17855562000247565b600085815260208120601f198616915b82811015620002ed57888601518255948401946001909101908401620002cc565b50858210156200030c5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6000602082840312156200032f57600080fd5b815180151581146200034057600080fd5b9392505050565b61128680620003576000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c8063642f6cef116100765780639d6f1cc71161005b5780639d6f1cc71461016b578063afb28d1f1461018b578063c98f10b01461019357600080fd5b8063642f6cef1461013b5780639525d5741461015857600080fd5b80634585e33b116100a75780634585e33b146101025780634b56a42e14610115578063601d5a711461012857600080fd5b806305e25131146100c357806340691db4146100d8575b600080fd5b6100d66100d136600461098e565b61019b565b005b6100eb6100e6366004610a3f565b6101b2565b6040516100f9929190610b1a565b60405180910390f35b6100d6610110366004610b3d565b610462565b6100eb610123366004610baf565b6105d1565b6100d6610136366004610c6c565b610627565b6000546101489060ff1681565b60405190151581526020016100f9565b6100d6610166366004610c6c565b610633565b61017e610179366004610ca1565b61063f565b6040516100f99190610cba565b61017e6106eb565b61017e6106f8565b80516101ae90600190602084019061078b565b5050565b6000606060006101c0610705565b90507fd1ffe9e45581c11d7d9f2ed5f75217cd4be9f8b7eee6af0f6d03f46de53956cd6101f060c0870187610cd4565b600081811061020157610201610d3c565b90506020020135036103da57600061021c60c0870187610cd4565b600181811061022d5761022d610d3c565b9050602002013560405160200161024691815260200190565b604051602081830303815290604052905060008180602001905181019061026d9190610d6b565b9050600061027e60c0890189610cd4565b600281811061028f5761028f610d3c565b905060200201356040516020016102a891815260200190565b60405160208183030381529060405290506000818060200190518101906102cf9190610d6b565b905060006102e060c08b018b610cd4565b60038181106102f1576102f1610d3c565b9050602002013560405160200161030a91815260200190565b60405160208183030381529060405290506000818060200190518101906103319190610dad565b604080516020810188905290810185905273ffffffffffffffffffffffffffffffffffffffff821660608201529091506002906001906003908a90608001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527f7ddd933e0000000000000000000000000000000000000000000000000000000082526103d19594939291600401610eb6565b60405180910390fd5b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f636f756c64206e6f742066696e64206d61746368696e67206576656e7420736960448201527f670000000000000000000000000000000000000000000000000000000000000060648201526084016103d1565b60008061047183850185610baf565b9150915060008060008380602001905181019061048e9190610f79565b92509250925060007309dff56a4ff44e0f4436260a04f5cfa65636a48173ffffffffffffffffffffffffffffffffffffffff16638e760afe876000815181106104d9576104d9610d3c565b60200260200101516040518263ffffffff1660e01b81526004016104fd9190610cba565b6000604051808303816000875af115801561051c573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526105629190810190610fae565b9050327f299a03817e683a32b21e29e3ae3c31f1c9c773f7d532836d116b62a9281fbc9d858585610591610705565b8b6000815181106105a4576105a4610d3c565b6020026020010151876040516105bf96959493929190611025565b60405180910390a25050505050505050565b60006060600084846040516020016105ea929190611085565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052600193509150505b9250929050565b60036101ae828261115f565b60026101ae828261115f565b6001818154811061064f57600080fd5b90600052602060002001600091509050805461066a90610dc8565b80601f016020809104026020016040519081016040528092919081815260200182805461069690610dc8565b80156106e35780601f106106b8576101008083540402835291602001916106e3565b820191906000526020600020905b8154815290600101906020018083116106c657829003601f168201915b505050505081565b6002805461066a90610dc8565b6003805461066a90610dc8565b6000805460ff161561078657606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561075d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107819190610d6b565b905090565b504390565b8280548282559060005260206000209081019282156107d1579160200282015b828111156107d157825182906107c1908261115f565b50916020019190600101906107ab565b506107dd9291506107e1565b5090565b808211156107dd5760006107f582826107fe565b506001016107e1565b50805461080a90610dc8565b6000825580601f1061081a575050565b601f016020900490600052602060002090810190610838919061083b565b50565b5b808211156107dd576000815560010161083c565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156108c6576108c6610850565b604052919050565b600067ffffffffffffffff8211156108e8576108e8610850565b5060051b60200190565b600067ffffffffffffffff82111561090c5761090c610850565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f83011261094957600080fd5b813561095c610957826108f2565b61087f565b81815284602083860101111561097157600080fd5b816020850160208301376000918101602001919091529392505050565b600060208083850312156109a157600080fd5b823567ffffffffffffffff808211156109b957600080fd5b818501915085601f8301126109cd57600080fd5b81356109db610957826108ce565b81815260059190911b830184019084810190888311156109fa57600080fd5b8585015b83811015610a3257803585811115610a165760008081fd5b610a248b89838a0101610938565b8452509186019186016109fe565b5098975050505050505050565b60008060408385031215610a5257600080fd5b823567ffffffffffffffff80821115610a6a57600080fd5b908401906101008287031215610a7f57600080fd5b90925060208401359080821115610a9557600080fd5b50610aa285828601610938565b9150509250929050565b60005b83811015610ac7578181015183820152602001610aaf565b50506000910152565b60008151808452610ae8816020860160208601610aac565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8215158152604060208201526000610b356040830184610ad0565b949350505050565b60008060208385031215610b5057600080fd5b823567ffffffffffffffff80821115610b6857600080fd5b818501915085601f830112610b7c57600080fd5b813581811115610b8b57600080fd5b866020828501011115610b9d57600080fd5b60209290920196919550909350505050565b60008060408385031215610bc257600080fd5b823567ffffffffffffffff80821115610bda57600080fd5b818501915085601f830112610bee57600080fd5b81356020610bfe610957836108ce565b82815260059290921b84018101918181019089841115610c1d57600080fd5b8286015b84811015610c5557803586811115610c395760008081fd5b610c478c86838b0101610938565b845250918301918301610c21565b5096505086013592505080821115610a9557600080fd5b600060208284031215610c7e57600080fd5b813567ffffffffffffffff811115610c9557600080fd5b610b3584828501610938565b600060208284031215610cb357600080fd5b5035919050565b602081526000610ccd6020830184610ad0565b9392505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610d0957600080fd5b83018035915067ffffffffffffffff821115610d2457600080fd5b6020019150600581901b360382131561062057600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600060208284031215610d7d57600080fd5b5051919050565b805173ffffffffffffffffffffffffffffffffffffffff81168114610da857600080fd5b919050565b600060208284031215610dbf57600080fd5b610ccd82610d84565b600181811c90821680610ddc57607f821691505b602082108103610e15577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60008154610e2881610dc8565b808552602060018381168015610e455760018114610e7d57610eab565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b8901019550610eab565b866000528260002060005b85811015610ea35781548a8201860152908301908401610e88565b890184019650505b505050505092915050565b60a081526000610ec960a0830188610e1b565b6020838203818501528188548084528284019150828160051b8501018a6000528360002060005b83811015610f3b577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0878403018552610f298383610e1b565b94860194925060019182019101610ef0565b50508681036040880152610f4f818b610e1b565b9450505050508460608401528281036080840152610f6d8185610ad0565b98975050505050505050565b600080600060608486031215610f8e57600080fd5b8351925060208401519150610fa560408501610d84565b90509250925092565b600060208284031215610fc057600080fd5b815167ffffffffffffffff811115610fd757600080fd5b8201601f81018413610fe857600080fd5b8051610ff6610957826108f2565b81815285602083850101111561100b57600080fd5b61101c826020830160208601610aac565b95945050505050565b86815285602082015273ffffffffffffffffffffffffffffffffffffffff8516604082015283606082015260c06080820152600061106660c0830185610ad0565b82810360a08401526110788185610ad0565b9998505050505050505050565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b838110156110fa577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa08887030185526110e8868351610ad0565b955093820193908201906001016110ae565b50508584038187015250505061101c8185610ad0565b601f82111561115a57600081815260208120601f850160051c810160208610156111375750805b601f850160051c820191505b8181101561115657828155600101611143565b5050505b505050565b815167ffffffffffffffff81111561117957611179610850565b61118d816111878454610dc8565b84611110565b602080601f8311600181146111e057600084156111aa5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555611156565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561122d5788860151825594840194600190910190840161120e565b508582101561126957878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b0190555056fea164736f6c6343000810000a307834353534343832643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030",
}

var LogTriggeredFeedLookupABI = LogTriggeredFeedLookupMetaData.ABI

var LogTriggeredFeedLookupBin = LogTriggeredFeedLookupMetaData.Bin

func DeployLogTriggeredFeedLookup(auth *bind.TransactOpts, backend bind.ContractBackend, _useArbitrumBlockNum bool) (common.Address, *types.Transaction, *LogTriggeredFeedLookup, error) {
	parsed, err := LogTriggeredFeedLookupMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LogTriggeredFeedLookupBin), backend, _useArbitrumBlockNum)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LogTriggeredFeedLookup{LogTriggeredFeedLookupCaller: LogTriggeredFeedLookupCaller{contract: contract}, LogTriggeredFeedLookupTransactor: LogTriggeredFeedLookupTransactor{contract: contract}, LogTriggeredFeedLookupFilterer: LogTriggeredFeedLookupFilterer{contract: contract}}, nil
}

type LogTriggeredFeedLookup struct {
	address common.Address
	abi     abi.ABI
	LogTriggeredFeedLookupCaller
	LogTriggeredFeedLookupTransactor
	LogTriggeredFeedLookupFilterer
}

type LogTriggeredFeedLookupCaller struct {
	contract *bind.BoundContract
}

type LogTriggeredFeedLookupTransactor struct {
	contract *bind.BoundContract
}

type LogTriggeredFeedLookupFilterer struct {
	contract *bind.BoundContract
}

type LogTriggeredFeedLookupSession struct {
	Contract     *LogTriggeredFeedLookup
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type LogTriggeredFeedLookupCallerSession struct {
	Contract *LogTriggeredFeedLookupCaller
	CallOpts bind.CallOpts
}

type LogTriggeredFeedLookupTransactorSession struct {
	Contract     *LogTriggeredFeedLookupTransactor
	TransactOpts bind.TransactOpts
}

type LogTriggeredFeedLookupRaw struct {
	Contract *LogTriggeredFeedLookup
}

type LogTriggeredFeedLookupCallerRaw struct {
	Contract *LogTriggeredFeedLookupCaller
}

type LogTriggeredFeedLookupTransactorRaw struct {
	Contract *LogTriggeredFeedLookupTransactor
}

func NewLogTriggeredFeedLookup(address common.Address, backend bind.ContractBackend) (*LogTriggeredFeedLookup, error) {
	abi, err := abi.JSON(strings.NewReader(LogTriggeredFeedLookupABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindLogTriggeredFeedLookup(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LogTriggeredFeedLookup{address: address, abi: abi, LogTriggeredFeedLookupCaller: LogTriggeredFeedLookupCaller{contract: contract}, LogTriggeredFeedLookupTransactor: LogTriggeredFeedLookupTransactor{contract: contract}, LogTriggeredFeedLookupFilterer: LogTriggeredFeedLookupFilterer{contract: contract}}, nil
}

func NewLogTriggeredFeedLookupCaller(address common.Address, caller bind.ContractCaller) (*LogTriggeredFeedLookupCaller, error) {
	contract, err := bindLogTriggeredFeedLookup(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LogTriggeredFeedLookupCaller{contract: contract}, nil
}

func NewLogTriggeredFeedLookupTransactor(address common.Address, transactor bind.ContractTransactor) (*LogTriggeredFeedLookupTransactor, error) {
	contract, err := bindLogTriggeredFeedLookup(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LogTriggeredFeedLookupTransactor{contract: contract}, nil
}

func NewLogTriggeredFeedLookupFilterer(address common.Address, filterer bind.ContractFilterer) (*LogTriggeredFeedLookupFilterer, error) {
	contract, err := bindLogTriggeredFeedLookup(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LogTriggeredFeedLookupFilterer{contract: contract}, nil
}

func bindLogTriggeredFeedLookup(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LogTriggeredFeedLookupMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LogTriggeredFeedLookup.Contract.LogTriggeredFeedLookupCaller.contract.Call(opts, result, method, params...)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.LogTriggeredFeedLookupTransactor.contract.Transfer(opts)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.LogTriggeredFeedLookupTransactor.contract.Transact(opts, method, params...)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LogTriggeredFeedLookup.Contract.contract.Call(opts, result, method, params...)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.contract.Transfer(opts)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.contract.Transact(opts, method, params...)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupCaller) CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (bool, []byte, error) {
	var out []interface{}
	err := _LogTriggeredFeedLookup.contract.Call(opts, &out, "checkCallback", values, extraData)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupSession) CheckCallback(values [][]byte, extraData []byte) (bool, []byte, error) {
	return _LogTriggeredFeedLookup.Contract.CheckCallback(&_LogTriggeredFeedLookup.CallOpts, values, extraData)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupCallerSession) CheckCallback(values [][]byte, extraData []byte) (bool, []byte, error) {
	return _LogTriggeredFeedLookup.Contract.CheckCallback(&_LogTriggeredFeedLookup.CallOpts, values, extraData)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupCaller) FeedParamKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LogTriggeredFeedLookup.contract.Call(opts, &out, "feedParamKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupSession) FeedParamKey() (string, error) {
	return _LogTriggeredFeedLookup.Contract.FeedParamKey(&_LogTriggeredFeedLookup.CallOpts)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupCallerSession) FeedParamKey() (string, error) {
	return _LogTriggeredFeedLookup.Contract.FeedParamKey(&_LogTriggeredFeedLookup.CallOpts)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupCaller) FeedsHex(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _LogTriggeredFeedLookup.contract.Call(opts, &out, "feedsHex", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupSession) FeedsHex(arg0 *big.Int) (string, error) {
	return _LogTriggeredFeedLookup.Contract.FeedsHex(&_LogTriggeredFeedLookup.CallOpts, arg0)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupCallerSession) FeedsHex(arg0 *big.Int) (string, error) {
	return _LogTriggeredFeedLookup.Contract.FeedsHex(&_LogTriggeredFeedLookup.CallOpts, arg0)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupCaller) TimeParamKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LogTriggeredFeedLookup.contract.Call(opts, &out, "timeParamKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupSession) TimeParamKey() (string, error) {
	return _LogTriggeredFeedLookup.Contract.TimeParamKey(&_LogTriggeredFeedLookup.CallOpts)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupCallerSession) TimeParamKey() (string, error) {
	return _LogTriggeredFeedLookup.Contract.TimeParamKey(&_LogTriggeredFeedLookup.CallOpts)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupCaller) UseArbitrumBlockNum(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _LogTriggeredFeedLookup.contract.Call(opts, &out, "useArbitrumBlockNum")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupSession) UseArbitrumBlockNum() (bool, error) {
	return _LogTriggeredFeedLookup.Contract.UseArbitrumBlockNum(&_LogTriggeredFeedLookup.CallOpts)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupCallerSession) UseArbitrumBlockNum() (bool, error) {
	return _LogTriggeredFeedLookup.Contract.UseArbitrumBlockNum(&_LogTriggeredFeedLookup.CallOpts)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupTransactor) CheckLog(opts *bind.TransactOpts, log Log, arg1 []byte) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.contract.Transact(opts, "checkLog", log, arg1)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupSession) CheckLog(log Log, arg1 []byte) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.CheckLog(&_LogTriggeredFeedLookup.TransactOpts, log, arg1)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupTransactorSession) CheckLog(log Log, arg1 []byte) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.CheckLog(&_LogTriggeredFeedLookup.TransactOpts, log, arg1)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.contract.Transact(opts, "performUpkeep", performData)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.PerformUpkeep(&_LogTriggeredFeedLookup.TransactOpts, performData)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.PerformUpkeep(&_LogTriggeredFeedLookup.TransactOpts, performData)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupTransactor) SetFeedParamKey(opts *bind.TransactOpts, feedParam string) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.contract.Transact(opts, "setFeedParamKey", feedParam)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupSession) SetFeedParamKey(feedParam string) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.SetFeedParamKey(&_LogTriggeredFeedLookup.TransactOpts, feedParam)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupTransactorSession) SetFeedParamKey(feedParam string) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.SetFeedParamKey(&_LogTriggeredFeedLookup.TransactOpts, feedParam)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupTransactor) SetFeedsHex(opts *bind.TransactOpts, newFeeds []string) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.contract.Transact(opts, "setFeedsHex", newFeeds)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupSession) SetFeedsHex(newFeeds []string) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.SetFeedsHex(&_LogTriggeredFeedLookup.TransactOpts, newFeeds)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupTransactorSession) SetFeedsHex(newFeeds []string) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.SetFeedsHex(&_LogTriggeredFeedLookup.TransactOpts, newFeeds)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupTransactor) SetTimeParamKey(opts *bind.TransactOpts, timeParam string) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.contract.Transact(opts, "setTimeParamKey", timeParam)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupSession) SetTimeParamKey(timeParam string) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.SetTimeParamKey(&_LogTriggeredFeedLookup.TransactOpts, timeParam)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupTransactorSession) SetTimeParamKey(timeParam string) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.SetTimeParamKey(&_LogTriggeredFeedLookup.TransactOpts, timeParam)
}

type LogTriggeredFeedLookupPerformingLogTriggerUpkeepIterator struct {
	Event *LogTriggeredFeedLookupPerformingLogTriggerUpkeep

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LogTriggeredFeedLookupPerformingLogTriggerUpkeepIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LogTriggeredFeedLookupPerformingLogTriggerUpkeep)
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
		it.Event = new(LogTriggeredFeedLookupPerformingLogTriggerUpkeep)
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

func (it *LogTriggeredFeedLookupPerformingLogTriggerUpkeepIterator) Error() error {
	return it.fail
}

func (it *LogTriggeredFeedLookupPerformingLogTriggerUpkeepIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LogTriggeredFeedLookupPerformingLogTriggerUpkeep struct {
	From        common.Address
	OrderId     *big.Int
	Amount      *big.Int
	Exchange    common.Address
	BlockNumber *big.Int
	Blob        []byte
	Verified    []byte
	Raw         types.Log
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupFilterer) FilterPerformingLogTriggerUpkeep(opts *bind.FilterOpts, from []common.Address) (*LogTriggeredFeedLookupPerformingLogTriggerUpkeepIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _LogTriggeredFeedLookup.contract.FilterLogs(opts, "PerformingLogTriggerUpkeep", fromRule)
	if err != nil {
		return nil, err
	}
	return &LogTriggeredFeedLookupPerformingLogTriggerUpkeepIterator{contract: _LogTriggeredFeedLookup.contract, event: "PerformingLogTriggerUpkeep", logs: logs, sub: sub}, nil
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupFilterer) WatchPerformingLogTriggerUpkeep(opts *bind.WatchOpts, sink chan<- *LogTriggeredFeedLookupPerformingLogTriggerUpkeep, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _LogTriggeredFeedLookup.contract.WatchLogs(opts, "PerformingLogTriggerUpkeep", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LogTriggeredFeedLookupPerformingLogTriggerUpkeep)
				if err := _LogTriggeredFeedLookup.contract.UnpackLog(event, "PerformingLogTriggerUpkeep", log); err != nil {
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

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupFilterer) ParsePerformingLogTriggerUpkeep(log types.Log) (*LogTriggeredFeedLookupPerformingLogTriggerUpkeep, error) {
	event := new(LogTriggeredFeedLookupPerformingLogTriggerUpkeep)
	if err := _LogTriggeredFeedLookup.contract.UnpackLog(event, "PerformingLogTriggerUpkeep", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookup) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _LogTriggeredFeedLookup.abi.Events["PerformingLogTriggerUpkeep"].ID:
		return _LogTriggeredFeedLookup.ParsePerformingLogTriggerUpkeep(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (LogTriggeredFeedLookupPerformingLogTriggerUpkeep) Topic() common.Hash {
	return common.HexToHash("0x299a03817e683a32b21e29e3ae3c31f1c9c773f7d532836d116b62a9281fbc9d")
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookup) Address() common.Address {
	return _LogTriggeredFeedLookup.address
}

type LogTriggeredFeedLookupInterface interface {
	CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (bool, []byte, error)

	FeedParamKey(opts *bind.CallOpts) (string, error)

	FeedsHex(opts *bind.CallOpts, arg0 *big.Int) (string, error)

	TimeParamKey(opts *bind.CallOpts) (string, error)

	UseArbitrumBlockNum(opts *bind.CallOpts) (bool, error)

	CheckLog(opts *bind.TransactOpts, log Log, arg1 []byte) (*types.Transaction, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	SetFeedParamKey(opts *bind.TransactOpts, feedParam string) (*types.Transaction, error)

	SetFeedsHex(opts *bind.TransactOpts, newFeeds []string) (*types.Transaction, error)

	SetTimeParamKey(opts *bind.TransactOpts, timeParam string) (*types.Transaction, error)

	FilterPerformingLogTriggerUpkeep(opts *bind.FilterOpts, from []common.Address) (*LogTriggeredFeedLookupPerformingLogTriggerUpkeepIterator, error)

	WatchPerformingLogTriggerUpkeep(opts *bind.WatchOpts, sink chan<- *LogTriggeredFeedLookupPerformingLogTriggerUpkeep, from []common.Address) (event.Subscription, error)

	ParsePerformingLogTriggerUpkeep(log types.Log) (*LogTriggeredFeedLookupPerformingLogTriggerUpkeep, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
