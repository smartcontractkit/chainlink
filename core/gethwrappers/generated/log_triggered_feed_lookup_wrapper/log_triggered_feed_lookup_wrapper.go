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
	ABI: "[{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_useArbitrumBlockNum\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"FeedLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"orderId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"exchange\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"blob\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"verified\",\"type\":\"bytes\"}],\"name\":\"PerformingLogTriggerUpkeep\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"txIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"topics\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLog\",\"name\":\"log\",\"type\":\"tuple\"}],\"name\":\"checkLog\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feedsHex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"targetContract\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"}],\"name\":\"getAdvancedLogTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"logTrigger\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"targetContract\",\"type\":\"address\"}],\"name\":\"getBasicLogTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"logTrigger\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"newFeeds\",\"type\":\"string[]\"}],\"name\":\"setFeedsHex\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbitrumBlockNum\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x610120604052604260a08181526080918291906200151960c03990526200002a906001908162000072565b503480156200003857600080fd5b506040516200155b3803806200155b8339810160408190526200005b91620001db565b6000805460ff191691151591909117905562000243565b828054828255906000526020600020908101928215620000c4579160200282015b82811115620000c45782518051620000b3918491602090910190620000d6565b509160200191906001019062000093565b50620000d292915062000161565b5090565b828054620000e49062000206565b90600052602060002090601f01602090048101928262000108576000855562000153565b82601f106200012357805160ff191683800117855562000153565b8280016001018555821562000153579182015b828111156200015357825182559160200191906001019062000136565b50620000d292915062000182565b80821115620000d257600062000178828262000199565b5060010162000161565b5b80821115620000d2576000815560010162000183565b508054620001a79062000206565b6000825580601f10620001b8575050565b601f016020900490600052602060002090810190620001d8919062000182565b50565b600060208284031215620001ee57600080fd5b81518015158114620001ff57600080fd5b9392505050565b600181811c908216806200021b57607f821691505b602082108114156200023d57634e487b7160e01b600052602260045260246000fd5b50919050565b6112c680620002536000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c80639d6f1cc711610076578063be61b7751161005b578063be61b77514610235578063c98f10b014610248578063cde36fb91461028457600080fd5b80639d6f1cc7146101e6578063afb28d1f146101f957600080fd5b80634b56a42e116100a75780634b56a42e146100eb578063642f6cef146101155780638f3dba411461013257600080fd5b806305e25131146100c35780634585e33b146100d8575b600080fd5b6100d66100d1366004610c0d565b610340565b005b6100d66100e6366004610cd8565b610357565b6100fe6100f9366004610b2e565b6103f9565b60405161010c929190610ef0565b60405180910390f35b6000546101229060ff1681565b604051901515815260200161010c565b6101d9610140366004610ae2565b6040805160c0808201835273ffffffffffffffffffffffffffffffffffffffff9690961680825260ff958616602080840191825283850196875260608085019687526000608080870182815260a09788019283528851948501969096529351909916828701529651968101969096529351938501939093529151918301919091529151818401528151808203909301835260e001905290565b60405161010c9190610f13565b6101d96101f4366004610d86565b61044f565b6101d96040518060400160405280600981526020017f666565644944486578000000000000000000000000000000000000000000000081525081565b6100fe610243366004610d4a565b6104fb565b6101d96040518060400160405280600b81526020017f626c6f636b4e756d62657200000000000000000000000000000000000000000081525081565b6101d9610292366004610aa8565b6040805160c0808201835273ffffffffffffffffffffffffffffffffffffffff93909316808252600060208084018281527fd1ffe9e45581c11d7d9f2ed5f75217cd4be9f8b7eee6af0f6d03f46de53956cd8587019081526060808701858152608080890187815260a0998a019788528a5196870198909852935160ff16858a015291519084015251908201529151928201929092529051818401528151808203909301835260e001905290565b80516103539060019060208401906108bf565b5050565b60008061036683850185610b2e565b915091506000806000838060200190518101906103839190610db8565b919450925090506060327f299a03817e683a32b21e29e3ae3c31f1c9c773f7d532836d116b62a9281fbc9d8585856103b961082a565b8b6000815181106103cc576103cc611239565b6020026020010151876040516103e7969594939291906110aa565b60405180910390a25050505050505050565b6000606060008484604051602001610412929190610e5c565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052600193509150505b9250929050565b6001818154811061045f57600080fd5b90600052602060002001600091509050805461047a906111e5565b80601f01602080910402602001604051908101604052809291908181526020018280546104a6906111e5565b80156104f35780601f106104c8576101008083540402835291602001916104f3565b820191906000526020600020905b8154815290600101906020018083116104d657829003601f168201915b505050505081565b60006060600061050961082a565b90507fd1ffe9e45581c11d7d9f2ed5f75217cd4be9f8b7eee6af0f6d03f46de53956cd61053960c086018661110a565b600081811061054a5761054a611239565b9050602002013514156107a257600061056660c086018661110a565b600181811061057757610577611239565b9050602002013560405160200161059091815260200190565b60405160208183030381529060405290506000818060200190518101906105b79190610d9f565b905060006105c860c088018861110a565b60028181106105d9576105d9611239565b905060200201356040516020016105f291815260200190565b60405160208183030381529060405290506000818060200190518101906106199190610d9f565b9050600061062a60c08a018a61110a565b600381811061063b5761063b611239565b9050602002013560405160200161065491815260200190565b604051602081830303815290604052905060008180602001905181019061067b9190610ac5565b90506040518060400160405280600981526020017f666565644944486578000000000000000000000000000000000000000000000081525060016040518060400160405280600b81526020017f626c6f636b4e756d6265720000000000000000000000000000000000000000008152508988878660405160200161073793929190928352602083019190915260601b7fffffffffffffffffffffffffffffffffffffffff00000000000000000000000016604082015260540190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527f7ddd933e0000000000000000000000000000000000000000000000000000000082526107999594939291600401610f26565b60405180910390fd5b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f636f756c64206e6f742066696e64206d61746368696e67206576656e7420736960448201527f67000000000000000000000000000000000000000000000000000000000000006064820152608401610799565b6000805460ff16156108ba57606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b15801561087d57600080fd5b505afa158015610891573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108b59190610d9f565b905090565b504390565b82805482825590600052602060002090810192821561090c579160200282015b8281111561090c57825180516108fc91849160209091019061091c565b50916020019190600101906108df565b5061091892915061099c565b5090565b828054610928906111e5565b90600052602060002090601f01602090048101928261094a5760008555610990565b82601f1061096357805160ff1916838001178555610990565b82800160010185558215610990579182015b82811115610990578251825591602001919060010190610975565b506109189291506109b9565b808211156109185760006109b082826109ce565b5060010161099c565b5b8082111561091857600081556001016109ba565b5080546109da906111e5565b6000825580601f106109ea575050565b601f016020900490600052602060002090810190610a0891906109b9565b50565b600067ffffffffffffffff831115610a2557610a25611268565b610a5660207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f86011601611172565b9050828152838383011115610a6a57600080fd5b828260208301376000602084830101529392505050565b600082601f830112610a9257600080fd5b610aa183833560208501610a0b565b9392505050565b600060208284031215610aba57600080fd5b8135610aa181611297565b600060208284031215610ad757600080fd5b8151610aa181611297565b60008060008060808587031215610af857600080fd5b8435610b0381611297565b9350602085013560ff81168114610b1957600080fd5b93969395505050506040820135916060013590565b60008060408385031215610b4157600080fd5b823567ffffffffffffffff80821115610b5957600080fd5b818501915085601f830112610b6d57600080fd5b81356020610b82610b7d836111c1565b611172565b8083825282820191508286018a848660051b8901011115610ba257600080fd5b60005b85811015610bdd57813587811115610bbc57600080fd5b610bca8d87838c0101610a81565b8552509284019290840190600101610ba5565b50909750505086013592505080821115610bf657600080fd5b50610c0385828601610a81565b9150509250929050565b60006020808385031215610c2057600080fd5b823567ffffffffffffffff80821115610c3857600080fd5b818501915085601f830112610c4c57600080fd5b8135610c5a610b7d826111c1565b80828252858201915085850189878560051b8801011115610c7a57600080fd5b60005b84811015610cc957813586811115610c9457600080fd5b8701603f81018c13610ca557600080fd5b610cb68c8a83013560408401610a0b565b8552509287019290870190600101610c7d565b50909998505050505050505050565b60008060208385031215610ceb57600080fd5b823567ffffffffffffffff80821115610d0357600080fd5b818501915085601f830112610d1757600080fd5b813581811115610d2657600080fd5b866020828501011115610d3857600080fd5b60209290920196919550909350505050565b600060208284031215610d5c57600080fd5b813567ffffffffffffffff811115610d7357600080fd5b82016101008185031215610aa157600080fd5b600060208284031215610d9857600080fd5b5035919050565b600060208284031215610db157600080fd5b5051919050565b600080600060608486031215610dcd57600080fd5b83519250602084015191506040840151610de681611297565b809150509250925092565b6000815180845260005b81811015610e1757602081850181015186830182015201610dfb565b81811115610e29576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b83811015610ed1577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0888703018552610ebf868351610df1565b95509382019390820190600101610e85565b505085840381870152505050610ee78185610df1565b95945050505050565b8215158152604060208201526000610f0b6040830184610df1565b949350505050565b602081526000610aa16020830184610df1565b60a081526000610f3960a0830188610df1565b6020838203818501528188548084528284019150828160051b85010160008b8152848120815b8481101561106b578784037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001865281548390600181811c9080831680610fa757607f831692505b8b8310811415610fde577f4e487b710000000000000000000000000000000000000000000000000000000088526022600452602488fd5b82895260208901818015610ff9576001811461102857611052565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00861682528d82019650611052565b6000898152602090208a5b8681101561104c57815484820152908501908f01611033565b83019750505b505050988a019892965050509190910190600101610f5f565b5050508681036040880152611080818b610df1565b945050505050846060840152828103608084015261109e8185610df1565b98975050505050505050565b86815285602082015273ffffffffffffffffffffffffffffffffffffffff8516604082015283606082015260c0608082015260006110eb60c0830185610df1565b82810360a08401526110fd8185610df1565b9998505050505050505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261113f57600080fd5b83018035915067ffffffffffffffff82111561115a57600080fd5b6020019150600581901b360382131561044857600080fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156111b9576111b9611268565b604052919050565b600067ffffffffffffffff8211156111db576111db611268565b5060051b60200190565b600181811c908216806111f957607f821691505b60208210811415611233577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff81168114610a0857600080fdfea164736f6c6343000806000a307834353534343832643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030",
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

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupTransactor) GetAdvancedLogTriggerConfig(opts *bind.TransactOpts, targetContract common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.contract.Transact(opts, "getAdvancedLogTriggerConfig", targetContract, selector, topic0, topic1)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupSession) GetAdvancedLogTriggerConfig(targetContract common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.GetAdvancedLogTriggerConfig(&_LogTriggeredFeedLookup.TransactOpts, targetContract, selector, topic0, topic1)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupTransactorSession) GetAdvancedLogTriggerConfig(targetContract common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.GetAdvancedLogTriggerConfig(&_LogTriggeredFeedLookup.TransactOpts, targetContract, selector, topic0, topic1)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupTransactor) GetBasicLogTriggerConfig(opts *bind.TransactOpts, targetContract common.Address) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.contract.Transact(opts, "getBasicLogTriggerConfig", targetContract)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupSession) GetBasicLogTriggerConfig(targetContract common.Address) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.GetBasicLogTriggerConfig(&_LogTriggeredFeedLookup.TransactOpts, targetContract)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupTransactorSession) GetBasicLogTriggerConfig(targetContract common.Address) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.GetBasicLogTriggerConfig(&_LogTriggeredFeedLookup.TransactOpts, targetContract)
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

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupTransactor) SetFeedsHex(opts *bind.TransactOpts, newFeeds []string) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.contract.Transact(opts, "setFeedsHex", newFeeds)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupSession) SetFeedsHex(newFeeds []string) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.SetFeedsHex(&_LogTriggeredFeedLookup.TransactOpts, newFeeds)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupTransactorSession) SetFeedsHex(newFeeds []string) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.SetFeedsHex(&_LogTriggeredFeedLookup.TransactOpts, newFeeds)
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

	GetAdvancedLogTriggerConfig(opts *bind.TransactOpts, targetContract common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte) (*types.Transaction, error)

	GetBasicLogTriggerConfig(opts *bind.TransactOpts, targetContract common.Address) (*types.Transaction, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	SetFeedsHex(opts *bind.TransactOpts, newFeeds []string) (*types.Transaction, error)

	FilterPerformingLogTriggerUpkeep(opts *bind.FilterOpts, from []common.Address) (*LogTriggeredFeedLookupPerformingLogTriggerUpkeepIterator, error)

	WatchPerformingLogTriggerUpkeep(opts *bind.WatchOpts, sink chan<- *LogTriggeredFeedLookupPerformingLogTriggerUpkeep, from []common.Address) (event.Subscription, error)

	ParsePerformingLogTriggerUpkeep(log types.Log) (*LogTriggeredFeedLookupPerformingLogTriggerUpkeep, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
