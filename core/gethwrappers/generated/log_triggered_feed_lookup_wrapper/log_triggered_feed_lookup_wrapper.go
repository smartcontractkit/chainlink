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
	ABI: "[{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_useArbitrumBlockNum\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"FeedLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"orderId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"exchange\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"blob\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"verified\",\"type\":\"bytes\"}],\"name\":\"PerformingLogTriggerUpkeep\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"txIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"topics\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLog\",\"name\":\"log\",\"type\":\"tuple\"}],\"name\":\"checkLog\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feedsHex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParam\",\"type\":\"string\"}],\"name\":\"setFeedParamKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"newFeeds\",\"type\":\"string[]\"}],\"name\":\"setFeedsHex\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"timeParam\",\"type\":\"string\"}],\"name\":\"setTimeParamKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbitrumBlockNum\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x610120604052604260a08181526080918291906200147560c03990526200002a9060019081620000d4565b506040805180820190915260098152680cccacac892c890caf60bb1b60208201526002906200005a908262000250565b5060408051808201909152600b81526a313637b1b5a73ab6b132b960a91b60208201526003906200008c908262000250565b503480156200009a57600080fd5b50604051620014b7380380620014b7833981016040819052620000bd916200031c565b6000805460ff191691151591909117905562000347565b8280548282559060005260206000209081019282156200011f579160200282015b828111156200011f57825182906200010e908262000250565b5091602001919060010190620000f5565b506200012d92915062000131565b5090565b808211156200012d57600062000148828262000152565b5060010162000131565b5080546200016090620001c1565b6000825580601f1062000171575050565b601f01602090049060005260206000209081019062000191919062000194565b50565b5b808211156200012d576000815560010162000195565b634e487b7160e01b600052604160045260246000fd5b600181811c90821680620001d657607f821691505b602082108103620001f757634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200024b57600081815260208120601f850160051c81016020861015620002265750805b601f850160051c820191505b81811015620002475782815560010162000232565b5050505b505050565b81516001600160401b038111156200026c576200026c620001ab565b62000284816200027d8454620001c1565b84620001fd565b602080601f831160018114620002bc5760008415620002a35750858301515b600019600386901b1c1916600185901b17855562000247565b600085815260208120601f198616915b82811015620002ed57888601518255948401946001909101908401620002cc565b50858210156200030c5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6000602082840312156200032f57600080fd5b815180151581146200034057600080fd5b9392505050565b61111e80620003576000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c80639525d57411610076578063afb28d1f1161005b578063afb28d1f14610178578063be61b77514610180578063c98f10b01461019357600080fd5b80639525d574146101455780639d6f1cc71461015857600080fd5b80634b56a42e116100a75780634b56a42e146100eb578063601d5a7114610115578063642f6cef1461012857600080fd5b806305e25131146100c35780634585e33b146100d8575b600080fd5b6100d66100d13660046108b3565b61019b565b005b6100d66100e6366004610969565b6101b2565b6100fe6100f93660046109db565b610254565b60405161010c929190610b13565b60405180910390f35b6100d6610123366004610b36565b6102aa565b6000546101359060ff1681565b604051901515815260200161010c565b6100d6610153366004610b36565b6102b6565b61016b610166366004610b6b565b6102c2565b60405161010c9190610b84565b61016b61036e565b6100fe61018e366004610b9e565b61037b565b61016b61062b565b80516101ae9060019060208401906106be565b5050565b6000806101c1838501856109db565b915091506000806000838060200190518101906101de9190610c03565b919450925090506060327f299a03817e683a32b21e29e3ae3c31f1c9c773f7d532836d116b62a9281fbc9d858585610214610638565b8b60008151811061022757610227610c38565b60200260200101518760405161024296959493929190610c67565b60405180910390a25050505050505050565b600060606000848460405160200161026d929190610cc7565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052600193509150505b9250929050565b60036101ae8282610dfd565b60026101ae8282610dfd565b600181815481106102d257600080fd5b9060005260206000200160009150905080546102ed90610d5b565b80601f016020809104026020016040519081016040528092919081815260200182805461031990610d5b565b80156103665780601f1061033b57610100808354040283529160200191610366565b820191906000526020600020905b81548152906001019060200180831161034957829003601f168201915b505050505081565b600280546102ed90610d5b565b600060606000610389610638565b90507fd1ffe9e45581c11d7d9f2ed5f75217cd4be9f8b7eee6af0f6d03f46de53956cd6103b960c0860186610f17565b60008181106103ca576103ca610c38565b90506020020135036105a35760006103e560c0860186610f17565b60018181106103f6576103f6610c38565b9050602002013560405160200161040f91815260200190565b60405160208183030381529060405290506000818060200190518101906104369190610f7f565b9050600061044760c0880188610f17565b600281811061045857610458610c38565b9050602002013560405160200161047191815260200190565b60405160208183030381529060405290506000818060200190518101906104989190610f7f565b905060006104a960c08a018a610f17565b60038181106104ba576104ba610c38565b905060200201356040516020016104d391815260200190565b60405160208183030381529060405290506000818060200190518101906104fa9190610f98565b604080516020810188905290810185905273ffffffffffffffffffffffffffffffffffffffff821660608201529091506002906001906003908a90608001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527f7ddd933e00000000000000000000000000000000000000000000000000000000825261059a959493929160040161104e565b60405180910390fd5b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f636f756c64206e6f742066696e64206d61746368696e67206576656e7420736960448201527f6700000000000000000000000000000000000000000000000000000000000000606482015260840161059a565b600380546102ed90610d5b565b6000805460ff16156106b957606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610690573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106b49190610f7f565b905090565b504390565b828054828255906000526020600020908101928215610704579160200282015b8281111561070457825182906106f49082610dfd565b50916020019190600101906106de565b50610710929150610714565b5090565b808211156107105760006107288282610731565b50600101610714565b50805461073d90610d5b565b6000825580601f1061074d575050565b601f01602090049060005260206000209081019061076b919061076e565b50565b5b80821115610710576000815560010161076f565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156107f9576107f9610783565b604052919050565b600067ffffffffffffffff82111561081b5761081b610783565b5060051b60200190565b600082601f83011261083657600080fd5b813567ffffffffffffffff81111561085057610850610783565b61088160207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016107b2565b81815284602083860101111561089657600080fd5b816020850160208301376000918101602001919091529392505050565b600060208083850312156108c657600080fd5b823567ffffffffffffffff808211156108de57600080fd5b818501915085601f8301126108f257600080fd5b813561090561090082610801565b6107b2565b81815260059190911b8301840190848101908883111561092457600080fd5b8585015b8381101561095c578035858111156109405760008081fd5b61094e8b89838a0101610825565b845250918601918601610928565b5098975050505050505050565b6000806020838503121561097c57600080fd5b823567ffffffffffffffff8082111561099457600080fd5b818501915085601f8301126109a857600080fd5b8135818111156109b757600080fd5b8660208285010111156109c957600080fd5b60209290920196919550909350505050565b600080604083850312156109ee57600080fd5b823567ffffffffffffffff80821115610a0657600080fd5b818501915085601f830112610a1a57600080fd5b81356020610a2a61090083610801565b82815260059290921b84018101918181019089841115610a4957600080fd5b8286015b84811015610a8157803586811115610a655760008081fd5b610a738c86838b0101610825565b845250918301918301610a4d565b5096505086013592505080821115610a9857600080fd5b50610aa585828601610825565b9150509250929050565b6000815180845260005b81811015610ad557602081850181015186830182015201610ab9565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b8215158152604060208201526000610b2e6040830184610aaf565b949350505050565b600060208284031215610b4857600080fd5b813567ffffffffffffffff811115610b5f57600080fd5b610b2e84828501610825565b600060208284031215610b7d57600080fd5b5035919050565b602081526000610b976020830184610aaf565b9392505050565b600060208284031215610bb057600080fd5b813567ffffffffffffffff811115610bc757600080fd5b82016101008185031215610b9757600080fd5b805173ffffffffffffffffffffffffffffffffffffffff81168114610bfe57600080fd5b919050565b600080600060608486031215610c1857600080fd5b8351925060208401519150610c2f60408501610bda565b90509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b86815285602082015273ffffffffffffffffffffffffffffffffffffffff8516604082015283606082015260c060808201526000610ca860c0830185610aaf565b82810360a0840152610cba8185610aaf565b9998505050505050505050565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b83811015610d3c577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0888703018552610d2a868351610aaf565b95509382019390820190600101610cf0565b505085840381870152505050610d528185610aaf565b95945050505050565b600181811c90821680610d6f57607f821691505b602082108103610da8577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f821115610df857600081815260208120601f850160051c81016020861015610dd55750805b601f850160051c820191505b81811015610df457828155600101610de1565b5050505b505050565b815167ffffffffffffffff811115610e1757610e17610783565b610e2b81610e258454610d5b565b84610dae565b602080601f831160018114610e7e5760008415610e485750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555610df4565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015610ecb57888601518255948401946001909101908401610eac565b5085821015610f0757878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610f4c57600080fd5b83018035915067ffffffffffffffff821115610f6757600080fd5b6020019150600581901b36038213156102a357600080fd5b600060208284031215610f9157600080fd5b5051919050565b600060208284031215610faa57600080fd5b610b9782610bda565b60008154610fc081610d5b565b808552602060018381168015610fdd576001811461101557611043565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b8901019550611043565b866000528260002060005b8581101561103b5781548a8201860152908301908401611020565b890184019650505b505050505092915050565b60a08152600061106160a0830188610fb3565b6020838203818501528188548084528284019150828160051b8501018a6000528360002060005b838110156110d3577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08784030185526110c18383610fb3565b94860194925060019182019101611088565b505086810360408801526110e7818b610fb3565b94505050505084606084015282810360808401526111058185610aaf565b9897505050505050505056fea164736f6c6343000810000a307834353534343832643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030",
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

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupCaller) CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (CheckCallback,

	error) {
	var out []interface{}
	err := _LogTriggeredFeedLookup.contract.Call(opts, &out, "checkCallback", values, extraData)

	outstruct := new(CheckCallback)
	if err != nil {
		return *outstruct, err
	}

	outstruct.UpkeepNeeded = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.PerformData = *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return *outstruct, err

}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupSession) CheckCallback(values [][]byte, extraData []byte) (CheckCallback,

	error) {
	return _LogTriggeredFeedLookup.Contract.CheckCallback(&_LogTriggeredFeedLookup.CallOpts, values, extraData)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupCallerSession) CheckCallback(values [][]byte, extraData []byte) (CheckCallback,

	error) {
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

type CheckCallback struct {
	UpkeepNeeded bool
	PerformData  []byte
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
	CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (CheckCallback,

		error)

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
