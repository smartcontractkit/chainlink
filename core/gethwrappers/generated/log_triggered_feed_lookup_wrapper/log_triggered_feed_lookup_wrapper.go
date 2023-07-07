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
	ABI: "[{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_useArbitrumBlockNum\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"FeedLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"orderId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"exchange\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"blob\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"verified\",\"type\":\"bytes\"}],\"name\":\"PerformingLogTriggerUpkeep\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"txIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"topics\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLog\",\"name\":\"log\",\"type\":\"tuple\"}],\"name\":\"checkLog\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feedsHex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParam\",\"type\":\"string\"}],\"name\":\"setFeedParamKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"newFeeds\",\"type\":\"string[]\"}],\"name\":\"setFeedsHex\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"timeParam\",\"type\":\"string\"}],\"name\":\"setTimeParamKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbitrumBlockNum\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x610120604052604260a0818152608091829190620015c360c03990526200002a9060019081620000d4565b506040805180820190915260098152680cccacac892c890caf60bb1b60208201526002906200005a908262000250565b5060408051808201909152600b81526a313637b1b5a73ab6b132b960a91b60208201526003906200008c908262000250565b503480156200009a57600080fd5b506040516200160538038062001605833981016040819052620000bd916200031c565b6000805460ff191691151591909117905562000347565b8280548282559060005260206000209081019282156200011f579160200282015b828111156200011f57825182906200010e908262000250565b5091602001919060010190620000f5565b506200012d92915062000131565b5090565b808211156200012d57600062000148828262000152565b5060010162000131565b5080546200016090620001c1565b6000825580601f1062000171575050565b601f01602090049060005260206000209081019062000191919062000194565b50565b5b808211156200012d576000815560010162000195565b634e487b7160e01b600052604160045260246000fd5b600181811c90821680620001d657607f821691505b602082108103620001f757634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200024b57600081815260208120601f850160051c81016020861015620002265750805b601f850160051c820191505b81811015620002475782815560010162000232565b5050505b505050565b81516001600160401b038111156200026c576200026c620001ab565b62000284816200027d8454620001c1565b84620001fd565b602080601f831160018114620002bc5760008415620002a35750858301515b600019600386901b1c1916600185901b17855562000247565b600085815260208120601f198616915b82811015620002ed57888601518255948401946001909101908401620002cc565b50858210156200030c5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6000602082840312156200032f57600080fd5b815180151581146200034057600080fd5b9392505050565b61126c80620003576000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c80639525d57411610076578063afb28d1f1161005b578063afb28d1f14610178578063be61b77514610180578063c98f10b01461019357600080fd5b80639525d574146101455780639d6f1cc71461015857600080fd5b80634b56a42e116100a75780634b56a42e146100eb578063601d5a7114610115578063642f6cef1461012857600080fd5b806305e25131146100c35780634585e33b146100d8575b600080fd5b6100d66100d136600461098e565b61019b565b005b6100d66100e6366004610a3f565b6101b2565b6100fe6100f9366004610ab1565b610321565b60405161010c929190610bf3565b60405180910390f35b6100d6610123366004610c16565b610377565b6000546101359060ff1681565b604051901515815260200161010c565b6100d6610153366004610c16565b610383565b61016b610166366004610c4b565b61038f565b60405161010c9190610c64565b61016b61043b565b6100fe61018e366004610c7e565b610448565b61016b6106f8565b80516101ae90600190602084019061078b565b5050565b6000806101c183850185610ab1565b915091506000806000838060200190518101906101de9190610ce3565b92509250925060007309dff56a4ff44e0f4436260a04f5cfa65636a48173ffffffffffffffffffffffffffffffffffffffff16638e760afe8760008151811061022957610229610d18565b60200260200101516040518263ffffffff1660e01b815260040161024d9190610c64565b6000604051808303816000875af115801561026c573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526102b29190810190610d47565b9050327f299a03817e683a32b21e29e3ae3c31f1c9c773f7d532836d116b62a9281fbc9d8585856102e1610705565b8b6000815181106102f4576102f4610d18565b60200260200101518760405161030f96959493929190610dbe565b60405180910390a25050505050505050565b600060606000848460405160200161033a929190610e1e565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052600193509150505b9250929050565b60036101ae8282610f4b565b60026101ae8282610f4b565b6001818154811061039f57600080fd5b9060005260206000200160009150905080546103ba90610ea9565b80601f01602080910402602001604051908101604052809291908181526020018280546103e690610ea9565b80156104335780601f1061040857610100808354040283529160200191610433565b820191906000526020600020905b81548152906001019060200180831161041657829003601f168201915b505050505081565b600280546103ba90610ea9565b600060606000610456610705565b90507fd1ffe9e45581c11d7d9f2ed5f75217cd4be9f8b7eee6af0f6d03f46de53956cd61048660c0860186611065565b600081811061049757610497610d18565b90506020020135036106705760006104b260c0860186611065565b60018181106104c3576104c3610d18565b905060200201356040516020016104dc91815260200190565b604051602081830303815290604052905060008180602001905181019061050391906110cd565b9050600061051460c0880188611065565b600281811061052557610525610d18565b9050602002013560405160200161053e91815260200190565b604051602081830303815290604052905060008180602001905181019061056591906110cd565b9050600061057660c08a018a611065565b600381811061058757610587610d18565b905060200201356040516020016105a091815260200190565b60405160208183030381529060405290506000818060200190518101906105c791906110e6565b604080516020810188905290810185905273ffffffffffffffffffffffffffffffffffffffff821660608201529091506002906001906003908a90608001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527f7ddd933e000000000000000000000000000000000000000000000000000000008252610667959493929160040161119c565b60405180910390fd5b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f636f756c64206e6f742066696e64206d61746368696e67206576656e7420736960448201527f67000000000000000000000000000000000000000000000000000000000000006064820152608401610667565b600380546103ba90610ea9565b6000805460ff161561078657606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561075d573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061078191906110cd565b905090565b504390565b8280548282559060005260206000209081019282156107d1579160200282015b828111156107d157825182906107c19082610f4b565b50916020019190600101906107ab565b506107dd9291506107e1565b5090565b808211156107dd5760006107f582826107fe565b506001016107e1565b50805461080a90610ea9565b6000825580601f1061081a575050565b601f016020900490600052602060002090810190610838919061083b565b50565b5b808211156107dd576000815560010161083c565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156108c6576108c6610850565b604052919050565b600067ffffffffffffffff8211156108e8576108e8610850565b5060051b60200190565b600067ffffffffffffffff82111561090c5761090c610850565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f83011261094957600080fd5b813561095c610957826108f2565b61087f565b81815284602083860101111561097157600080fd5b816020850160208301376000918101602001919091529392505050565b600060208083850312156109a157600080fd5b823567ffffffffffffffff808211156109b957600080fd5b818501915085601f8301126109cd57600080fd5b81356109db610957826108ce565b81815260059190911b830184019084810190888311156109fa57600080fd5b8585015b83811015610a3257803585811115610a165760008081fd5b610a248b89838a0101610938565b8452509186019186016109fe565b5098975050505050505050565b60008060208385031215610a5257600080fd5b823567ffffffffffffffff80821115610a6a57600080fd5b818501915085601f830112610a7e57600080fd5b813581811115610a8d57600080fd5b866020828501011115610a9f57600080fd5b60209290920196919550909350505050565b60008060408385031215610ac457600080fd5b823567ffffffffffffffff80821115610adc57600080fd5b818501915085601f830112610af057600080fd5b81356020610b00610957836108ce565b82815260059290921b84018101918181019089841115610b1f57600080fd5b8286015b84811015610b5757803586811115610b3b5760008081fd5b610b498c86838b0101610938565b845250918301918301610b23565b5096505086013592505080821115610b6e57600080fd5b50610b7b85828601610938565b9150509250929050565b60005b83811015610ba0578181015183820152602001610b88565b50506000910152565b60008151808452610bc1816020860160208601610b85565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8215158152604060208201526000610c0e6040830184610ba9565b949350505050565b600060208284031215610c2857600080fd5b813567ffffffffffffffff811115610c3f57600080fd5b610c0e84828501610938565b600060208284031215610c5d57600080fd5b5035919050565b602081526000610c776020830184610ba9565b9392505050565b600060208284031215610c9057600080fd5b813567ffffffffffffffff811115610ca757600080fd5b82016101008185031215610c7757600080fd5b805173ffffffffffffffffffffffffffffffffffffffff81168114610cde57600080fd5b919050565b600080600060608486031215610cf857600080fd5b8351925060208401519150610d0f60408501610cba565b90509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600060208284031215610d5957600080fd5b815167ffffffffffffffff811115610d7057600080fd5b8201601f81018413610d8157600080fd5b8051610d8f610957826108f2565b818152856020838501011115610da457600080fd5b610db5826020830160208601610b85565b95945050505050565b86815285602082015273ffffffffffffffffffffffffffffffffffffffff8516604082015283606082015260c060808201526000610dff60c0830185610ba9565b82810360a0840152610e118185610ba9565b9998505050505050505050565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b83811015610e93577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0888703018552610e81868351610ba9565b95509382019390820190600101610e47565b505085840381870152505050610db58185610ba9565b600181811c90821680610ebd57607f821691505b602082108103610ef6577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f821115610f4657600081815260208120601f850160051c81016020861015610f235750805b601f850160051c820191505b81811015610f4257828155600101610f2f565b5050505b505050565b815167ffffffffffffffff811115610f6557610f65610850565b610f7981610f738454610ea9565b84610efc565b602080601f831160018114610fcc5760008415610f965750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555610f42565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561101957888601518255948401946001909101908401610ffa565b508582101561105557878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261109a57600080fd5b83018035915067ffffffffffffffff8211156110b557600080fd5b6020019150600581901b360382131561037057600080fd5b6000602082840312156110df57600080fd5b5051919050565b6000602082840312156110f857600080fd5b610c7782610cba565b6000815461110e81610ea9565b80855260206001838116801561112b576001811461116357611191565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b8901019550611191565b866000528260002060005b858110156111895781548a820186015290830190840161116e565b890184019650505b505050505092915050565b60a0815260006111af60a0830188611101565b6020838203818501528188548084528284019150828160051b8501018a6000528360002060005b83811015611221577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe087840301855261120f8383611101565b948601949250600191820191016111d6565b50508681036040880152611235818b611101565b94505050505084606084015282810360808401526112538185610ba9565b9897505050505050505056fea164736f6c6343000810000a307834353534343832643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030",
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

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupTransactor) CheckLog(opts *bind.TransactOpts, log Log) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.contract.Transact(opts, "checkLog", log)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupSession) CheckLog(log Log) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.CheckLog(&_LogTriggeredFeedLookup.TransactOpts, log)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupTransactorSession) CheckLog(log Log) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.CheckLog(&_LogTriggeredFeedLookup.TransactOpts, log)
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

	CheckLog(opts *bind.TransactOpts, log Log) (*types.Transaction, error)

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
