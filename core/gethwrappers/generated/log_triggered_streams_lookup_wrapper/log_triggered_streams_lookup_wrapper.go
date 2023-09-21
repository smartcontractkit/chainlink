// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package log_triggered_streams_lookup_wrapper

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
	Timestamp   *big.Int
	TxHash      [32]byte
	BlockNumber *big.Int
	BlockHash   [32]byte
	Source      common.Address
	Topics      [][32]byte
	Data        []byte
}

var LogTriggeredStreamsLookupMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_useArbitrumBlockNum\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"_verify\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"StreamsLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"orderId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"exchange\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"blob\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"verified\",\"type\":\"bytes\"}],\"name\":\"PerformingLogTriggerUpkeep\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"topics\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLog\",\"name\":\"log\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"checkLog\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feedsHex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParam\",\"type\":\"string\"}],\"name\":\"setFeedParamKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"newFeeds\",\"type\":\"string[]\"}],\"name\":\"setFeedsHex\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"timeParam\",\"type\":\"string\"}],\"name\":\"setTimeParamKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbitrumBlockNum\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x610120604052604260a08181526080918291906200164e60c03990526200002a9060019081620000e5565b506040805180820190915260098152680cccacac892c890caf60bb1b60208201526002906200005a908262000261565b5060408051808201909152600b81526a313637b1b5a73ab6b132b960a91b60208201526003906200008c908262000261565b503480156200009a57600080fd5b506040516200169038038062001690833981016040819052620000bd9162000343565b6000805461ffff191692151561ff00191692909217610100911515919091021790556200037b565b82805482825590600052602060002090810192821562000130579160200282015b828111156200013057825182906200011f908262000261565b509160200191906001019062000106565b506200013e92915062000142565b5090565b808211156200013e57600062000159828262000163565b5060010162000142565b5080546200017190620001d2565b6000825580601f1062000182575050565b601f016020900490600052602060002090810190620001a29190620001a5565b50565b5b808211156200013e5760008155600101620001a6565b634e487b7160e01b600052604160045260246000fd5b600181811c90821680620001e757607f821691505b6020821081036200020857634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200025c57600081815260208120601f850160051c81016020861015620002375750805b601f850160051c820191505b81811015620002585782815560010162000243565b5050505b505050565b81516001600160401b038111156200027d576200027d620001bc565b62000295816200028e8454620001d2565b846200020e565b602080601f831160018114620002cd5760008415620002b45750858301515b600019600386901b1c1916600185901b17855562000258565b600085815260208120601f198616915b82811015620002fe57888601518255948401946001909101908401620002dd565b50858210156200031d5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b805180151581146200033e57600080fd5b919050565b600080604083850312156200035757600080fd5b62000362836200032d565b915062000372602084016200032d565b90509250929050565b6112c3806200038b6000396000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c8063642f6cef11610081578063afb28d1f1161005b578063afb28d1f14610196578063c98f10b01461019e578063fc735e99146101a657600080fd5b8063642f6cef146101465780639525d574146101635780639d6f1cc71461017657600080fd5b80634585e33b116100b25780634585e33b1461010d5780634b56a42e14610120578063601d5a711461013357600080fd5b806305e25131146100ce57806340691db4146100e3575b600080fd5b6100e16100dc3660046109cb565b6101b8565b005b6100f66100f1366004610a7c565b6101cf565b604051610104929190610b57565b60405180910390f35b6100e161011b366004610b7a565b61047f565b6100f661012e366004610bec565b61060e565b6100e1610141366004610ca9565b610664565b6000546101539060ff1681565b6040519015158152602001610104565b6100e1610171366004610ca9565b610670565b610189610184366004610cde565b61067c565b6040516101049190610cf7565b610189610728565b610189610735565b60005461015390610100900460ff1681565b80516101cb9060019060208401906107c8565b5050565b6000606060006101dd610742565b90507fd1ffe9e45581c11d7d9f2ed5f75217cd4be9f8b7eee6af0f6d03f46de53956cd61020d60c0870187610d11565b600081811061021e5761021e610d79565b90506020020135036103f757600061023960c0870187610d11565b600181811061024a5761024a610d79565b9050602002013560405160200161026391815260200190565b604051602081830303815290604052905060008180602001905181019061028a9190610da8565b9050600061029b60c0890189610d11565b60028181106102ac576102ac610d79565b905060200201356040516020016102c591815260200190565b60405160208183030381529060405290506000818060200190518101906102ec9190610da8565b905060006102fd60c08b018b610d11565b600381811061030e5761030e610d79565b9050602002013560405160200161032791815260200190565b604051602081830303815290604052905060008180602001905181019061034e9190610dea565b604080516020810188905290810185905273ffffffffffffffffffffffffffffffffffffffff821660608201529091506002906001906003908a90608001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527ff055e4a20000000000000000000000000000000000000000000000000000000082526103ee9594939291600401610ef3565b60405180910390fd5b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f636f756c64206e6f742066696e64206d61746368696e67206576656e7420736960448201527f670000000000000000000000000000000000000000000000000000000000000060648201526084016103ee565b60008061048e83850185610bec565b915091506000806000838060200190518101906104ab9190610fb6565b6040805160208101909152600080825254939650919450925090610100900460ff16156105a1577309dff56a4ff44e0f4436260a04f5cfa65636a48173ffffffffffffffffffffffffffffffffffffffff16638e760afe8760008151811061051557610515610d79565b60200260200101516040518263ffffffff1660e01b81526004016105399190610cf7565b6000604051808303816000875af1158015610558573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261059e9190810190610feb565b90505b327f299a03817e683a32b21e29e3ae3c31f1c9c773f7d532836d116b62a9281fbc9d8585856105ce610742565b8b6000815181106105e1576105e1610d79565b6020026020010151876040516105fc96959493929190611062565b60405180910390a25050505050505050565b60006060600084846040516020016106279291906110c2565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052600193509150505b9250929050565b60036101cb828261119c565b60026101cb828261119c565b6001818154811061068c57600080fd5b9060005260206000200160009150905080546106a790610e05565b80601f01602080910402602001604051908101604052809291908181526020018280546106d390610e05565b80156107205780601f106106f557610100808354040283529160200191610720565b820191906000526020600020905b81548152906001019060200180831161070357829003601f168201915b505050505081565b600280546106a790610e05565b600380546106a790610e05565b6000805460ff16156107c357606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561079a573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107be9190610da8565b905090565b504390565b82805482825590600052602060002090810192821561080e579160200282015b8281111561080e57825182906107fe908261119c565b50916020019190600101906107e8565b5061081a92915061081e565b5090565b8082111561081a576000610832828261083b565b5060010161081e565b50805461084790610e05565b6000825580601f10610857575050565b601f0160209004906000526020600020908101906108759190610878565b50565b5b8082111561081a5760008155600101610879565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156109035761090361088d565b604052919050565b600067ffffffffffffffff8211156109255761092561088d565b5060051b60200190565b600067ffffffffffffffff8211156109495761094961088d565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f83011261098657600080fd5b81356109996109948261092f565b6108bc565b8181528460208386010111156109ae57600080fd5b816020850160208301376000918101602001919091529392505050565b600060208083850312156109de57600080fd5b823567ffffffffffffffff808211156109f657600080fd5b818501915085601f830112610a0a57600080fd5b8135610a186109948261090b565b81815260059190911b83018401908481019088831115610a3757600080fd5b8585015b83811015610a6f57803585811115610a535760008081fd5b610a618b89838a0101610975565b845250918601918601610a3b565b5098975050505050505050565b60008060408385031215610a8f57600080fd5b823567ffffffffffffffff80821115610aa757600080fd5b908401906101008287031215610abc57600080fd5b90925060208401359080821115610ad257600080fd5b50610adf85828601610975565b9150509250929050565b60005b83811015610b04578181015183820152602001610aec565b50506000910152565b60008151808452610b25816020860160208601610ae9565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8215158152604060208201526000610b726040830184610b0d565b949350505050565b60008060208385031215610b8d57600080fd5b823567ffffffffffffffff80821115610ba557600080fd5b818501915085601f830112610bb957600080fd5b813581811115610bc857600080fd5b866020828501011115610bda57600080fd5b60209290920196919550909350505050565b60008060408385031215610bff57600080fd5b823567ffffffffffffffff80821115610c1757600080fd5b818501915085601f830112610c2b57600080fd5b81356020610c3b6109948361090b565b82815260059290921b84018101918181019089841115610c5a57600080fd5b8286015b84811015610c9257803586811115610c765760008081fd5b610c848c86838b0101610975565b845250918301918301610c5e565b5096505086013592505080821115610ad257600080fd5b600060208284031215610cbb57600080fd5b813567ffffffffffffffff811115610cd257600080fd5b610b7284828501610975565b600060208284031215610cf057600080fd5b5035919050565b602081526000610d0a6020830184610b0d565b9392505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610d4657600080fd5b83018035915067ffffffffffffffff821115610d6157600080fd5b6020019150600581901b360382131561065d57600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600060208284031215610dba57600080fd5b5051919050565b805173ffffffffffffffffffffffffffffffffffffffff81168114610de557600080fd5b919050565b600060208284031215610dfc57600080fd5b610d0a82610dc1565b600181811c90821680610e1957607f821691505b602082108103610e52577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60008154610e6581610e05565b808552602060018381168015610e825760018114610eba57610ee8565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b8901019550610ee8565b866000528260002060005b85811015610ee05781548a8201860152908301908401610ec5565b890184019650505b505050505092915050565b60a081526000610f0660a0830188610e58565b6020838203818501528188548084528284019150828160051b8501018a6000528360002060005b83811015610f78577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0878403018552610f668383610e58565b94860194925060019182019101610f2d565b50508681036040880152610f8c818b610e58565b9450505050508460608401528281036080840152610faa8185610b0d565b98975050505050505050565b600080600060608486031215610fcb57600080fd5b8351925060208401519150610fe260408501610dc1565b90509250925092565b600060208284031215610ffd57600080fd5b815167ffffffffffffffff81111561101457600080fd5b8201601f8101841361102557600080fd5b80516110336109948261092f565b81815285602083850101111561104857600080fd5b611059826020830160208601610ae9565b95945050505050565b86815285602082015273ffffffffffffffffffffffffffffffffffffffff8516604082015283606082015260c0608082015260006110a360c0830185610b0d565b82810360a08401526110b58185610b0d565b9998505050505050505050565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b83811015611137577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0888703018552611125868351610b0d565b955093820193908201906001016110eb565b5050858403818701525050506110598185610b0d565b601f82111561119757600081815260208120601f850160051c810160208610156111745750805b601f850160051c820191505b8181101561119357828155600101611180565b5050505b505050565b815167ffffffffffffffff8111156111b6576111b661088d565b6111ca816111c48454610e05565b8461114d565b602080601f83116001811461121d57600084156111e75750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555611193565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561126a5788860151825594840194600190910190840161124b565b50858210156112a657878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b0190555056fea164736f6c6343000810000a307834353534343832643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030",
}

var LogTriggeredStreamsLookupABI = LogTriggeredStreamsLookupMetaData.ABI

var LogTriggeredStreamsLookupBin = LogTriggeredStreamsLookupMetaData.Bin

func DeployLogTriggeredStreamsLookup(auth *bind.TransactOpts, backend bind.ContractBackend, _useArbitrumBlockNum bool, _verify bool) (common.Address, *types.Transaction, *LogTriggeredStreamsLookup, error) {
	parsed, err := LogTriggeredStreamsLookupMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LogTriggeredStreamsLookupBin), backend, _useArbitrumBlockNum, _verify)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LogTriggeredStreamsLookup{LogTriggeredStreamsLookupCaller: LogTriggeredStreamsLookupCaller{contract: contract}, LogTriggeredStreamsLookupTransactor: LogTriggeredStreamsLookupTransactor{contract: contract}, LogTriggeredStreamsLookupFilterer: LogTriggeredStreamsLookupFilterer{contract: contract}}, nil
}

type LogTriggeredStreamsLookup struct {
	address common.Address
	abi     abi.ABI
	LogTriggeredStreamsLookupCaller
	LogTriggeredStreamsLookupTransactor
	LogTriggeredStreamsLookupFilterer
}

type LogTriggeredStreamsLookupCaller struct {
	contract *bind.BoundContract
}

type LogTriggeredStreamsLookupTransactor struct {
	contract *bind.BoundContract
}

type LogTriggeredStreamsLookupFilterer struct {
	contract *bind.BoundContract
}

type LogTriggeredStreamsLookupSession struct {
	Contract     *LogTriggeredStreamsLookup
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type LogTriggeredStreamsLookupCallerSession struct {
	Contract *LogTriggeredStreamsLookupCaller
	CallOpts bind.CallOpts
}

type LogTriggeredStreamsLookupTransactorSession struct {
	Contract     *LogTriggeredStreamsLookupTransactor
	TransactOpts bind.TransactOpts
}

type LogTriggeredStreamsLookupRaw struct {
	Contract *LogTriggeredStreamsLookup
}

type LogTriggeredStreamsLookupCallerRaw struct {
	Contract *LogTriggeredStreamsLookupCaller
}

type LogTriggeredStreamsLookupTransactorRaw struct {
	Contract *LogTriggeredStreamsLookupTransactor
}

func NewLogTriggeredStreamsLookup(address common.Address, backend bind.ContractBackend) (*LogTriggeredStreamsLookup, error) {
	abi, err := abi.JSON(strings.NewReader(LogTriggeredStreamsLookupABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindLogTriggeredStreamsLookup(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LogTriggeredStreamsLookup{address: address, abi: abi, LogTriggeredStreamsLookupCaller: LogTriggeredStreamsLookupCaller{contract: contract}, LogTriggeredStreamsLookupTransactor: LogTriggeredStreamsLookupTransactor{contract: contract}, LogTriggeredStreamsLookupFilterer: LogTriggeredStreamsLookupFilterer{contract: contract}}, nil
}

func NewLogTriggeredStreamsLookupCaller(address common.Address, caller bind.ContractCaller) (*LogTriggeredStreamsLookupCaller, error) {
	contract, err := bindLogTriggeredStreamsLookup(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LogTriggeredStreamsLookupCaller{contract: contract}, nil
}

func NewLogTriggeredStreamsLookupTransactor(address common.Address, transactor bind.ContractTransactor) (*LogTriggeredStreamsLookupTransactor, error) {
	contract, err := bindLogTriggeredStreamsLookup(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LogTriggeredStreamsLookupTransactor{contract: contract}, nil
}

func NewLogTriggeredStreamsLookupFilterer(address common.Address, filterer bind.ContractFilterer) (*LogTriggeredStreamsLookupFilterer, error) {
	contract, err := bindLogTriggeredStreamsLookup(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LogTriggeredStreamsLookupFilterer{contract: contract}, nil
}

func bindLogTriggeredStreamsLookup(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LogTriggeredStreamsLookupMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LogTriggeredStreamsLookup.Contract.LogTriggeredStreamsLookupCaller.contract.Call(opts, result, method, params...)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LogTriggeredStreamsLookup.Contract.LogTriggeredStreamsLookupTransactor.contract.Transfer(opts)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LogTriggeredStreamsLookup.Contract.LogTriggeredStreamsLookupTransactor.contract.Transact(opts, method, params...)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LogTriggeredStreamsLookup.Contract.contract.Call(opts, result, method, params...)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LogTriggeredStreamsLookup.Contract.contract.Transfer(opts)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LogTriggeredStreamsLookup.Contract.contract.Transact(opts, method, params...)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupCaller) CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (bool, []byte, error) {
	var out []interface{}
	err := _LogTriggeredStreamsLookup.contract.Call(opts, &out, "checkCallback", values, extraData)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupSession) CheckCallback(values [][]byte, extraData []byte) (bool, []byte, error) {
	return _LogTriggeredStreamsLookup.Contract.CheckCallback(&_LogTriggeredStreamsLookup.CallOpts, values, extraData)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupCallerSession) CheckCallback(values [][]byte, extraData []byte) (bool, []byte, error) {
	return _LogTriggeredStreamsLookup.Contract.CheckCallback(&_LogTriggeredStreamsLookup.CallOpts, values, extraData)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupCaller) FeedParamKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LogTriggeredStreamsLookup.contract.Call(opts, &out, "feedParamKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupSession) FeedParamKey() (string, error) {
	return _LogTriggeredStreamsLookup.Contract.FeedParamKey(&_LogTriggeredStreamsLookup.CallOpts)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupCallerSession) FeedParamKey() (string, error) {
	return _LogTriggeredStreamsLookup.Contract.FeedParamKey(&_LogTriggeredStreamsLookup.CallOpts)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupCaller) FeedsHex(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _LogTriggeredStreamsLookup.contract.Call(opts, &out, "feedsHex", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupSession) FeedsHex(arg0 *big.Int) (string, error) {
	return _LogTriggeredStreamsLookup.Contract.FeedsHex(&_LogTriggeredStreamsLookup.CallOpts, arg0)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupCallerSession) FeedsHex(arg0 *big.Int) (string, error) {
	return _LogTriggeredStreamsLookup.Contract.FeedsHex(&_LogTriggeredStreamsLookup.CallOpts, arg0)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupCaller) TimeParamKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LogTriggeredStreamsLookup.contract.Call(opts, &out, "timeParamKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupSession) TimeParamKey() (string, error) {
	return _LogTriggeredStreamsLookup.Contract.TimeParamKey(&_LogTriggeredStreamsLookup.CallOpts)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupCallerSession) TimeParamKey() (string, error) {
	return _LogTriggeredStreamsLookup.Contract.TimeParamKey(&_LogTriggeredStreamsLookup.CallOpts)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupCaller) UseArbitrumBlockNum(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _LogTriggeredStreamsLookup.contract.Call(opts, &out, "useArbitrumBlockNum")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupSession) UseArbitrumBlockNum() (bool, error) {
	return _LogTriggeredStreamsLookup.Contract.UseArbitrumBlockNum(&_LogTriggeredStreamsLookup.CallOpts)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupCallerSession) UseArbitrumBlockNum() (bool, error) {
	return _LogTriggeredStreamsLookup.Contract.UseArbitrumBlockNum(&_LogTriggeredStreamsLookup.CallOpts)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupCaller) Verify(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _LogTriggeredStreamsLookup.contract.Call(opts, &out, "verify")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupSession) Verify() (bool, error) {
	return _LogTriggeredStreamsLookup.Contract.Verify(&_LogTriggeredStreamsLookup.CallOpts)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupCallerSession) Verify() (bool, error) {
	return _LogTriggeredStreamsLookup.Contract.Verify(&_LogTriggeredStreamsLookup.CallOpts)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupTransactor) CheckLog(opts *bind.TransactOpts, log Log, arg1 []byte) (*types.Transaction, error) {
	return _LogTriggeredStreamsLookup.contract.Transact(opts, "checkLog", log, arg1)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupSession) CheckLog(log Log, arg1 []byte) (*types.Transaction, error) {
	return _LogTriggeredStreamsLookup.Contract.CheckLog(&_LogTriggeredStreamsLookup.TransactOpts, log, arg1)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupTransactorSession) CheckLog(log Log, arg1 []byte) (*types.Transaction, error) {
	return _LogTriggeredStreamsLookup.Contract.CheckLog(&_LogTriggeredStreamsLookup.TransactOpts, log, arg1)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _LogTriggeredStreamsLookup.contract.Transact(opts, "performUpkeep", performData)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _LogTriggeredStreamsLookup.Contract.PerformUpkeep(&_LogTriggeredStreamsLookup.TransactOpts, performData)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _LogTriggeredStreamsLookup.Contract.PerformUpkeep(&_LogTriggeredStreamsLookup.TransactOpts, performData)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupTransactor) SetFeedParamKey(opts *bind.TransactOpts, feedParam string) (*types.Transaction, error) {
	return _LogTriggeredStreamsLookup.contract.Transact(opts, "setFeedParamKey", feedParam)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupSession) SetFeedParamKey(feedParam string) (*types.Transaction, error) {
	return _LogTriggeredStreamsLookup.Contract.SetFeedParamKey(&_LogTriggeredStreamsLookup.TransactOpts, feedParam)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupTransactorSession) SetFeedParamKey(feedParam string) (*types.Transaction, error) {
	return _LogTriggeredStreamsLookup.Contract.SetFeedParamKey(&_LogTriggeredStreamsLookup.TransactOpts, feedParam)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupTransactor) SetFeedsHex(opts *bind.TransactOpts, newFeeds []string) (*types.Transaction, error) {
	return _LogTriggeredStreamsLookup.contract.Transact(opts, "setFeedsHex", newFeeds)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupSession) SetFeedsHex(newFeeds []string) (*types.Transaction, error) {
	return _LogTriggeredStreamsLookup.Contract.SetFeedsHex(&_LogTriggeredStreamsLookup.TransactOpts, newFeeds)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupTransactorSession) SetFeedsHex(newFeeds []string) (*types.Transaction, error) {
	return _LogTriggeredStreamsLookup.Contract.SetFeedsHex(&_LogTriggeredStreamsLookup.TransactOpts, newFeeds)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupTransactor) SetTimeParamKey(opts *bind.TransactOpts, timeParam string) (*types.Transaction, error) {
	return _LogTriggeredStreamsLookup.contract.Transact(opts, "setTimeParamKey", timeParam)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupSession) SetTimeParamKey(timeParam string) (*types.Transaction, error) {
	return _LogTriggeredStreamsLookup.Contract.SetTimeParamKey(&_LogTriggeredStreamsLookup.TransactOpts, timeParam)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupTransactorSession) SetTimeParamKey(timeParam string) (*types.Transaction, error) {
	return _LogTriggeredStreamsLookup.Contract.SetTimeParamKey(&_LogTriggeredStreamsLookup.TransactOpts, timeParam)
}

type LogTriggeredStreamsLookupPerformingLogTriggerUpkeepIterator struct {
	Event *LogTriggeredStreamsLookupPerformingLogTriggerUpkeep

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LogTriggeredStreamsLookupPerformingLogTriggerUpkeepIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LogTriggeredStreamsLookupPerformingLogTriggerUpkeep)
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
		it.Event = new(LogTriggeredStreamsLookupPerformingLogTriggerUpkeep)
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

func (it *LogTriggeredStreamsLookupPerformingLogTriggerUpkeepIterator) Error() error {
	return it.fail
}

func (it *LogTriggeredStreamsLookupPerformingLogTriggerUpkeepIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LogTriggeredStreamsLookupPerformingLogTriggerUpkeep struct {
	From        common.Address
	OrderId     *big.Int
	Amount      *big.Int
	Exchange    common.Address
	BlockNumber *big.Int
	Blob        []byte
	Verified    []byte
	Raw         types.Log
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupFilterer) FilterPerformingLogTriggerUpkeep(opts *bind.FilterOpts, from []common.Address) (*LogTriggeredStreamsLookupPerformingLogTriggerUpkeepIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _LogTriggeredStreamsLookup.contract.FilterLogs(opts, "PerformingLogTriggerUpkeep", fromRule)
	if err != nil {
		return nil, err
	}
	return &LogTriggeredStreamsLookupPerformingLogTriggerUpkeepIterator{contract: _LogTriggeredStreamsLookup.contract, event: "PerformingLogTriggerUpkeep", logs: logs, sub: sub}, nil
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupFilterer) WatchPerformingLogTriggerUpkeep(opts *bind.WatchOpts, sink chan<- *LogTriggeredStreamsLookupPerformingLogTriggerUpkeep, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _LogTriggeredStreamsLookup.contract.WatchLogs(opts, "PerformingLogTriggerUpkeep", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LogTriggeredStreamsLookupPerformingLogTriggerUpkeep)
				if err := _LogTriggeredStreamsLookup.contract.UnpackLog(event, "PerformingLogTriggerUpkeep", log); err != nil {
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

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupFilterer) ParsePerformingLogTriggerUpkeep(log types.Log) (*LogTriggeredStreamsLookupPerformingLogTriggerUpkeep, error) {
	event := new(LogTriggeredStreamsLookupPerformingLogTriggerUpkeep)
	if err := _LogTriggeredStreamsLookup.contract.UnpackLog(event, "PerformingLogTriggerUpkeep", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookup) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _LogTriggeredStreamsLookup.abi.Events["PerformingLogTriggerUpkeep"].ID:
		return _LogTriggeredStreamsLookup.ParsePerformingLogTriggerUpkeep(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (LogTriggeredStreamsLookupPerformingLogTriggerUpkeep) Topic() common.Hash {
	return common.HexToHash("0x299a03817e683a32b21e29e3ae3c31f1c9c773f7d532836d116b62a9281fbc9d")
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookup) Address() common.Address {
	return _LogTriggeredStreamsLookup.address
}

type LogTriggeredStreamsLookupInterface interface {
	CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (bool, []byte, error)

	FeedParamKey(opts *bind.CallOpts) (string, error)

	FeedsHex(opts *bind.CallOpts, arg0 *big.Int) (string, error)

	TimeParamKey(opts *bind.CallOpts) (string, error)

	UseArbitrumBlockNum(opts *bind.CallOpts) (bool, error)

	Verify(opts *bind.CallOpts) (bool, error)

	CheckLog(opts *bind.TransactOpts, log Log, arg1 []byte) (*types.Transaction, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	SetFeedParamKey(opts *bind.TransactOpts, feedParam string) (*types.Transaction, error)

	SetFeedsHex(opts *bind.TransactOpts, newFeeds []string) (*types.Transaction, error)

	SetTimeParamKey(opts *bind.TransactOpts, timeParam string) (*types.Transaction, error)

	FilterPerformingLogTriggerUpkeep(opts *bind.FilterOpts, from []common.Address) (*LogTriggeredStreamsLookupPerformingLogTriggerUpkeepIterator, error)

	WatchPerformingLogTriggerUpkeep(opts *bind.WatchOpts, sink chan<- *LogTriggeredStreamsLookupPerformingLogTriggerUpkeep, from []common.Address) (event.Subscription, error)

	ParsePerformingLogTriggerUpkeep(log types.Log) (*LogTriggeredStreamsLookupPerformingLogTriggerUpkeep, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
