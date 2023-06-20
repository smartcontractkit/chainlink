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
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_useArbitrumBlockNum\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"FeedLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"bn\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"ts\",\"type\":\"uint256\"}],\"name\":\"DoNotTriggerMercury\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"logBlockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"PerformingLogTriggerUpkeep\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"logBlockNumber\",\"type\":\"uint256\"}],\"name\":\"TriggerMercury\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"txIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"topics\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLog\",\"name\":\"log\",\"type\":\"tuple\"}],\"name\":\"checkLog\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feedsHex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"interval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"logTrigger\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"newFeeds\",\"type\":\"string[]\"}],\"name\":\"setFeedsHex\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setSpread\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"startLogs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbitrumBlockNum\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x610120604052604260a0818152608091829190620015c560c03990526200002b9060009060016200008b565b503480156200003957600080fd5b5060405162001607380380620016078339810160408190526200005c91620001f4565b6001929092556002556006805460ff191691151591909117905560006003819055600481905560055562000271565b828054828255906000526020600020908101928215620000dd579160200282015b82811115620000dd5782518051620000cc918491602090910190620000ef565b5091602001919060010190620000ac565b50620000eb9291506200017a565b5090565b828054620000fd9062000234565b90600052602060002090601f0160209004810192826200012157600085556200016c565b82601f106200013c57805160ff19168380011785556200016c565b828001600101855582156200016c579182015b828111156200016c5782518255916020019190600101906200014f565b50620000eb9291506200019b565b80821115620000eb576000620001918282620001b2565b506001016200017a565b5b80821115620000eb57600081556001016200019c565b508054620001c09062000234565b6000825580601f10620001d1575050565b601f016020900490600052602060002090810190620001f191906200019b565b50565b6000806000606084860312156200020a57600080fd5b8351925060208401519150604084015180151581146200022957600080fd5b809150509250925092565b600181811c908216806200024957607f821691505b602082108114156200026b57634e487b7160e01b600052602260045260246000fd5b50919050565b61134480620002816000396000f3fe608060405234801561001057600080fd5b506004361061011b5760003560e01c80637f407edf116100b2578063afb28d1f11610081578063c98f10b011610066578063c98f10b0146102f3578063d832d92f1461032f578063eadaa41f1461033757600080fd5b8063afb28d1f146102a4578063be61b775146102e057600080fd5b80637f407edf1461025c578063917d895f1461027f578063947a36fb146102885780639d6f1cc71461029157600080fd5b80634b56a42e116100ee5780634b56a42e1461020c57806361bc221a1461022d5780636250a13a14610236578063642f6cef1461023f57600080fd5b806305e2513114610120578063214d536b146101355780632cb15864146101e25780634585e33b146101f9575b600080fd5b61013361012e366004610ca2565b61034a565b005b6040805160c0808201835230808352600060208085018281527fcd89a1cdede3e128a8e92d77495b16cc12f0fc7564a712113f006adaf640a4a686880190815260608088018581526080808a0187815260a09a8b019788528b5196870198909852935160ff16858b015291519084015251908201529151938201939093529151828201528251808303909101815260e09091019091525b6040516101d99190610fb5565b60405180910390f35b6101eb60045481565b6040519081526020016101d9565b610133610207366004610d6d565b610361565b61021f61021a366004610bc3565b610473565b6040516101d9929190610f92565b6101eb60055481565b6101eb60015481565b60065461024c9060ff1681565b60405190151581526020016101d9565b61013361026a366004610e4d565b60019190915560025560006004819055600555565b6101eb60035481565b6101eb60025481565b6101cc61029f366004610e1b565b6104c9565b6101cc6040518060400160405280600981526020017f666565644944486578000000000000000000000000000000000000000000000081525081565b61021f6102ee366004610ddf565b610575565b6101cc6040518060400160405280600b81526020017f626c6f636b4e756d62657200000000000000000000000000000000000000000081525081565b61024c610887565b610133610345366004610e1b565b6108d8565b805161035d9060009060208401906109da565b5050565b600061036b610943565b90506004546000141561037e5760048190555b60055461038c906001611227565b60055560038190556000806103a384860186610bc3565b91509150600080828060200190518101906103be9190610e6f565b9150915084827fcd89a1cdede3e128a8e92d77495b16cc12f0fc7564a712113f006adaf640a4a660405160405180910390a3604051429086907f3338104ccab396f091b767e17a2a863f70d777261982306eeb72053c94b9cd4790600090a360055460408051848152602081019290925281018290526060810186905232907ff694f423fadb03d2b520837cc15d973963b694ca830f618df6b0f105fc37dafd9060800160405180910390a250505050505050565b600060606000848460405160200161048c929190610efe565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052600193509150505b9250929050565b600081815481106104d957600080fd5b9060005260206000200160009150905080546104f490611256565b80601f016020809104026020016040519081016040528092919081815260200182805461052090611256565b801561056d5780601f106105425761010080835404028352916020019161056d565b820191906000526020600020905b81548152906001019060200180831161055057829003601f168201915b505050505081565b600060606000610583610943565b90507fcd89a1cdede3e128a8e92d77495b16cc12f0fc7564a712113f006adaf640a4a66105b360c086018661114c565b60008181106105c4576105c46112d9565b9050602002013514156107ff5760006105e060c086018661114c565b60018181106105f1576105f16112d9565b9050602002013560405160200161060a91815260200190565b60405160208183030381529060405290506000818060200190518101906106319190610e34565b9050600061064260c088018861114c565b6002818110610653576106536112d9565b9050602002013560405160200161066c91815260200190565b60405160208183030381529060405290506000818060200190518101906106939190610e34565b9050600061069f610887565b9050801561079d576040518060400160405280600981526020017f666565644944486578000000000000000000000000000000000000000000000081525060006040518060400160405280600b81526020017f626c6f636b4e756d626572000000000000000000000000000000000000000000815250888786604051602001610732929190918252602082015260400190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527f7ddd933e0000000000000000000000000000000000000000000000000000000082526107949594939291600401610fc8565b60405180910390fd5b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f75706b656570206e6f74206e65656465640000000000000000000000000000006044820152606401610794565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f636f756c64206e6f742066696e64206d61746368696e67206576656e7420736960448201527f67000000000000000000000000000000000000000000000000000000000000006064820152608401610794565b60006004546000141561089a5750600190565b60006108a4610943565b9050600154600454826108b7919061123f565b1080156108d257506002546003546108cf908361123f565b10155b91505090565b60006108e2610943565b905080827fcd89a1cdede3e128a8e92d77495b16cc12f0fc7564a712113f006adaf640a4a660405160405180910390a3604051429082907f3338104ccab396f091b767e17a2a863f70d777261982306eeb72053c94b9cd4790600090a35050565b60065460009060ff16156109d557606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b15801561099857600080fd5b505afa1580156109ac573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109d09190610e34565b905090565b504390565b828054828255906000526020600020908101928215610a27579160200282015b82811115610a275782518051610a17918491602090910190610a37565b50916020019190600101906109fa565b50610a33929150610ab7565b5090565b828054610a4390611256565b90600052602060002090601f016020900481019282610a655760008555610aab565b82601f10610a7e57805160ff1916838001178555610aab565b82800160010185558215610aab579182015b82811115610aab578251825591602001919060010190610a90565b50610a33929150610ad4565b80821115610a33576000610acb8282610ae9565b50600101610ab7565b5b80821115610a335760008155600101610ad5565b508054610af590611256565b6000825580601f10610b05575050565b601f016020900490600052602060002090810190610b239190610ad4565b50565b600067ffffffffffffffff831115610b4057610b40611308565b610b7160207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f860116016111b4565b9050828152838383011115610b8557600080fd5b828260208301376000602084830101529392505050565b600082601f830112610bad57600080fd5b610bbc83833560208501610b26565b9392505050565b60008060408385031215610bd657600080fd5b823567ffffffffffffffff80821115610bee57600080fd5b818501915085601f830112610c0257600080fd5b81356020610c17610c1283611203565b6111b4565b8083825282820191508286018a848660051b8901011115610c3757600080fd5b60005b85811015610c7257813587811115610c5157600080fd5b610c5f8d87838c0101610b9c565b8552509284019290840190600101610c3a565b50909750505086013592505080821115610c8b57600080fd5b50610c9885828601610b9c565b9150509250929050565b60006020808385031215610cb557600080fd5b823567ffffffffffffffff80821115610ccd57600080fd5b818501915085601f830112610ce157600080fd5b8135610cef610c1282611203565b80828252858201915085850189878560051b8801011115610d0f57600080fd5b60005b84811015610d5e57813586811115610d2957600080fd5b8701603f81018c13610d3a57600080fd5b610d4b8c8a83013560408401610b26565b8552509287019290870190600101610d12565b50909998505050505050505050565b60008060208385031215610d8057600080fd5b823567ffffffffffffffff80821115610d9857600080fd5b818501915085601f830112610dac57600080fd5b813581811115610dbb57600080fd5b866020828501011115610dcd57600080fd5b60209290920196919550909350505050565b600060208284031215610df157600080fd5b813567ffffffffffffffff811115610e0857600080fd5b82016101008185031215610bbc57600080fd5b600060208284031215610e2d57600080fd5b5035919050565b600060208284031215610e4657600080fd5b5051919050565b60008060408385031215610e6057600080fd5b50508035926020909101359150565b60008060408385031215610e8257600080fd5b505080516020909101519092909150565b6000815180845260005b81811015610eb957602081850181015186830182015201610e9d565b81811115610ecb576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b83811015610f73577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0888703018552610f61868351610e93565b95509382019390820190600101610f27565b505085840381870152505050610f898185610e93565b95945050505050565b8215158152604060208201526000610fad6040830184610e93565b949350505050565b602081526000610bbc6020830184610e93565b60a081526000610fdb60a0830188610e93565b6020838203818501528188548084528284019150828160051b85010160008b8152848120815b8481101561110d578784037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001865281548390600181811c908083168061104957607f831692505b8b8310811415611080577f4e487b710000000000000000000000000000000000000000000000000000000088526022600452602488fd5b8289526020890181801561109b57600181146110ca576110f4565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00861682528d820196506110f4565b6000898152602090208a5b868110156110ee57815484820152908501908f016110d5565b83019750505b505050988a019892965050509190910190600101611001565b5050508681036040880152611122818b610e93565b94505050505084606084015282810360808401526111408185610e93565b98975050505050505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261118157600080fd5b83018035915067ffffffffffffffff82111561119c57600080fd5b6020019150600581901b36038213156104c257600080fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156111fb576111fb611308565b604052919050565b600067ffffffffffffffff82111561121d5761121d611308565b5060051b60200190565b6000821982111561123a5761123a6112aa565b500190565b600082821015611251576112516112aa565b500390565b600181811c9082168061126a57607f821691505b602082108114156112a4577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a307834353534343832643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030",
}

var LogTriggeredFeedLookupABI = LogTriggeredFeedLookupMetaData.ABI

var LogTriggeredFeedLookupBin = LogTriggeredFeedLookupMetaData.Bin

func DeployLogTriggeredFeedLookup(auth *bind.TransactOpts, backend bind.ContractBackend, _testRange *big.Int, _interval *big.Int, _useArbitrumBlockNum bool) (common.Address, *types.Transaction, *LogTriggeredFeedLookup, error) {
	parsed, err := LogTriggeredFeedLookupMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LogTriggeredFeedLookupBin), backend, _testRange, _interval, _useArbitrumBlockNum)
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

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LogTriggeredFeedLookup.contract.Call(opts, &out, "counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupSession) Counter() (*big.Int, error) {
	return _LogTriggeredFeedLookup.Contract.Counter(&_LogTriggeredFeedLookup.CallOpts)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupCallerSession) Counter() (*big.Int, error) {
	return _LogTriggeredFeedLookup.Contract.Counter(&_LogTriggeredFeedLookup.CallOpts)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupCaller) Eligible(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _LogTriggeredFeedLookup.contract.Call(opts, &out, "eligible")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupSession) Eligible() (bool, error) {
	return _LogTriggeredFeedLookup.Contract.Eligible(&_LogTriggeredFeedLookup.CallOpts)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupCallerSession) Eligible() (bool, error) {
	return _LogTriggeredFeedLookup.Contract.Eligible(&_LogTriggeredFeedLookup.CallOpts)
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

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupCaller) InitialBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LogTriggeredFeedLookup.contract.Call(opts, &out, "initialBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupSession) InitialBlock() (*big.Int, error) {
	return _LogTriggeredFeedLookup.Contract.InitialBlock(&_LogTriggeredFeedLookup.CallOpts)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupCallerSession) InitialBlock() (*big.Int, error) {
	return _LogTriggeredFeedLookup.Contract.InitialBlock(&_LogTriggeredFeedLookup.CallOpts)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupCaller) Interval(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LogTriggeredFeedLookup.contract.Call(opts, &out, "interval")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupSession) Interval() (*big.Int, error) {
	return _LogTriggeredFeedLookup.Contract.Interval(&_LogTriggeredFeedLookup.CallOpts)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupCallerSession) Interval() (*big.Int, error) {
	return _LogTriggeredFeedLookup.Contract.Interval(&_LogTriggeredFeedLookup.CallOpts)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupCaller) PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LogTriggeredFeedLookup.contract.Call(opts, &out, "previousPerformBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupSession) PreviousPerformBlock() (*big.Int, error) {
	return _LogTriggeredFeedLookup.Contract.PreviousPerformBlock(&_LogTriggeredFeedLookup.CallOpts)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupCallerSession) PreviousPerformBlock() (*big.Int, error) {
	return _LogTriggeredFeedLookup.Contract.PreviousPerformBlock(&_LogTriggeredFeedLookup.CallOpts)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupCaller) TestRange(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LogTriggeredFeedLookup.contract.Call(opts, &out, "testRange")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupSession) TestRange() (*big.Int, error) {
	return _LogTriggeredFeedLookup.Contract.TestRange(&_LogTriggeredFeedLookup.CallOpts)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupCallerSession) TestRange() (*big.Int, error) {
	return _LogTriggeredFeedLookup.Contract.TestRange(&_LogTriggeredFeedLookup.CallOpts)
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

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupTransactor) RegisterUpkeep(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.contract.Transact(opts, "registerUpkeep")
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupSession) RegisterUpkeep() (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.RegisterUpkeep(&_LogTriggeredFeedLookup.TransactOpts)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupTransactorSession) RegisterUpkeep() (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.RegisterUpkeep(&_LogTriggeredFeedLookup.TransactOpts)
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

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupTransactor) SetSpread(opts *bind.TransactOpts, _testRange *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.contract.Transact(opts, "setSpread", _testRange, _interval)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupSession) SetSpread(_testRange *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.SetSpread(&_LogTriggeredFeedLookup.TransactOpts, _testRange, _interval)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupTransactorSession) SetSpread(_testRange *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.SetSpread(&_LogTriggeredFeedLookup.TransactOpts, _testRange, _interval)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupTransactor) StartLogs(opts *bind.TransactOpts, upkeepId *big.Int) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.contract.Transact(opts, "startLogs", upkeepId)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupSession) StartLogs(upkeepId *big.Int) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.StartLogs(&_LogTriggeredFeedLookup.TransactOpts, upkeepId)
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupTransactorSession) StartLogs(upkeepId *big.Int) (*types.Transaction, error) {
	return _LogTriggeredFeedLookup.Contract.StartLogs(&_LogTriggeredFeedLookup.TransactOpts, upkeepId)
}

type LogTriggeredFeedLookupDoNotTriggerMercuryIterator struct {
	Event *LogTriggeredFeedLookupDoNotTriggerMercury

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LogTriggeredFeedLookupDoNotTriggerMercuryIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LogTriggeredFeedLookupDoNotTriggerMercury)
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
		it.Event = new(LogTriggeredFeedLookupDoNotTriggerMercury)
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

func (it *LogTriggeredFeedLookupDoNotTriggerMercuryIterator) Error() error {
	return it.fail
}

func (it *LogTriggeredFeedLookupDoNotTriggerMercuryIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LogTriggeredFeedLookupDoNotTriggerMercury struct {
	Bn  *big.Int
	Ts  *big.Int
	Raw types.Log
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupFilterer) FilterDoNotTriggerMercury(opts *bind.FilterOpts, bn []*big.Int, ts []*big.Int) (*LogTriggeredFeedLookupDoNotTriggerMercuryIterator, error) {

	var bnRule []interface{}
	for _, bnItem := range bn {
		bnRule = append(bnRule, bnItem)
	}
	var tsRule []interface{}
	for _, tsItem := range ts {
		tsRule = append(tsRule, tsItem)
	}

	logs, sub, err := _LogTriggeredFeedLookup.contract.FilterLogs(opts, "DoNotTriggerMercury", bnRule, tsRule)
	if err != nil {
		return nil, err
	}
	return &LogTriggeredFeedLookupDoNotTriggerMercuryIterator{contract: _LogTriggeredFeedLookup.contract, event: "DoNotTriggerMercury", logs: logs, sub: sub}, nil
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupFilterer) WatchDoNotTriggerMercury(opts *bind.WatchOpts, sink chan<- *LogTriggeredFeedLookupDoNotTriggerMercury, bn []*big.Int, ts []*big.Int) (event.Subscription, error) {

	var bnRule []interface{}
	for _, bnItem := range bn {
		bnRule = append(bnRule, bnItem)
	}
	var tsRule []interface{}
	for _, tsItem := range ts {
		tsRule = append(tsRule, tsItem)
	}

	logs, sub, err := _LogTriggeredFeedLookup.contract.WatchLogs(opts, "DoNotTriggerMercury", bnRule, tsRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LogTriggeredFeedLookupDoNotTriggerMercury)
				if err := _LogTriggeredFeedLookup.contract.UnpackLog(event, "DoNotTriggerMercury", log); err != nil {
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

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupFilterer) ParseDoNotTriggerMercury(log types.Log) (*LogTriggeredFeedLookupDoNotTriggerMercury, error) {
	event := new(LogTriggeredFeedLookupDoNotTriggerMercury)
	if err := _LogTriggeredFeedLookup.contract.UnpackLog(event, "DoNotTriggerMercury", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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
	From           common.Address
	UpkeepId       *big.Int
	Counter        *big.Int
	LogBlockNumber *big.Int
	BlockNumber    *big.Int
	Raw            types.Log
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

type LogTriggeredFeedLookupTriggerMercuryIterator struct {
	Event *LogTriggeredFeedLookupTriggerMercury

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LogTriggeredFeedLookupTriggerMercuryIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LogTriggeredFeedLookupTriggerMercury)
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
		it.Event = new(LogTriggeredFeedLookupTriggerMercury)
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

func (it *LogTriggeredFeedLookupTriggerMercuryIterator) Error() error {
	return it.fail
}

func (it *LogTriggeredFeedLookupTriggerMercuryIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LogTriggeredFeedLookupTriggerMercury struct {
	UpkeepId       *big.Int
	LogBlockNumber *big.Int
	Raw            types.Log
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupFilterer) FilterTriggerMercury(opts *bind.FilterOpts, upkeepId []*big.Int, logBlockNumber []*big.Int) (*LogTriggeredFeedLookupTriggerMercuryIterator, error) {

	var upkeepIdRule []interface{}
	for _, upkeepIdItem := range upkeepId {
		upkeepIdRule = append(upkeepIdRule, upkeepIdItem)
	}
	var logBlockNumberRule []interface{}
	for _, logBlockNumberItem := range logBlockNumber {
		logBlockNumberRule = append(logBlockNumberRule, logBlockNumberItem)
	}

	logs, sub, err := _LogTriggeredFeedLookup.contract.FilterLogs(opts, "TriggerMercury", upkeepIdRule, logBlockNumberRule)
	if err != nil {
		return nil, err
	}
	return &LogTriggeredFeedLookupTriggerMercuryIterator{contract: _LogTriggeredFeedLookup.contract, event: "TriggerMercury", logs: logs, sub: sub}, nil
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupFilterer) WatchTriggerMercury(opts *bind.WatchOpts, sink chan<- *LogTriggeredFeedLookupTriggerMercury, upkeepId []*big.Int, logBlockNumber []*big.Int) (event.Subscription, error) {

	var upkeepIdRule []interface{}
	for _, upkeepIdItem := range upkeepId {
		upkeepIdRule = append(upkeepIdRule, upkeepIdItem)
	}
	var logBlockNumberRule []interface{}
	for _, logBlockNumberItem := range logBlockNumber {
		logBlockNumberRule = append(logBlockNumberRule, logBlockNumberItem)
	}

	logs, sub, err := _LogTriggeredFeedLookup.contract.WatchLogs(opts, "TriggerMercury", upkeepIdRule, logBlockNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LogTriggeredFeedLookupTriggerMercury)
				if err := _LogTriggeredFeedLookup.contract.UnpackLog(event, "TriggerMercury", log); err != nil {
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

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookupFilterer) ParseTriggerMercury(log types.Log) (*LogTriggeredFeedLookupTriggerMercury, error) {
	event := new(LogTriggeredFeedLookupTriggerMercury)
	if err := _LogTriggeredFeedLookup.contract.UnpackLog(event, "TriggerMercury", log); err != nil {
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
	case _LogTriggeredFeedLookup.abi.Events["DoNotTriggerMercury"].ID:
		return _LogTriggeredFeedLookup.ParseDoNotTriggerMercury(log)
	case _LogTriggeredFeedLookup.abi.Events["PerformingLogTriggerUpkeep"].ID:
		return _LogTriggeredFeedLookup.ParsePerformingLogTriggerUpkeep(log)
	case _LogTriggeredFeedLookup.abi.Events["TriggerMercury"].ID:
		return _LogTriggeredFeedLookup.ParseTriggerMercury(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (LogTriggeredFeedLookupDoNotTriggerMercury) Topic() common.Hash {
	return common.HexToHash("0x3338104ccab396f091b767e17a2a863f70d777261982306eeb72053c94b9cd47")
}

func (LogTriggeredFeedLookupPerformingLogTriggerUpkeep) Topic() common.Hash {
	return common.HexToHash("0xf694f423fadb03d2b520837cc15d973963b694ca830f618df6b0f105fc37dafd")
}

func (LogTriggeredFeedLookupTriggerMercury) Topic() common.Hash {
	return common.HexToHash("0xcd89a1cdede3e128a8e92d77495b16cc12f0fc7564a712113f006adaf640a4a6")
}

func (_LogTriggeredFeedLookup *LogTriggeredFeedLookup) Address() common.Address {
	return _LogTriggeredFeedLookup.address
}

type LogTriggeredFeedLookupInterface interface {
	CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (CheckCallback,

		error)

	Counter(opts *bind.CallOpts) (*big.Int, error)

	Eligible(opts *bind.CallOpts) (bool, error)

	FeedParamKey(opts *bind.CallOpts) (string, error)

	FeedsHex(opts *bind.CallOpts, arg0 *big.Int) (string, error)

	InitialBlock(opts *bind.CallOpts) (*big.Int, error)

	Interval(opts *bind.CallOpts) (*big.Int, error)

	PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error)

	TestRange(opts *bind.CallOpts) (*big.Int, error)

	TimeParamKey(opts *bind.CallOpts) (string, error)

	UseArbitrumBlockNum(opts *bind.CallOpts) (bool, error)

	CheckLog(opts *bind.TransactOpts, log Log) (*types.Transaction, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	RegisterUpkeep(opts *bind.TransactOpts) (*types.Transaction, error)

	SetFeedsHex(opts *bind.TransactOpts, newFeeds []string) (*types.Transaction, error)

	SetSpread(opts *bind.TransactOpts, _testRange *big.Int, _interval *big.Int) (*types.Transaction, error)

	StartLogs(opts *bind.TransactOpts, upkeepId *big.Int) (*types.Transaction, error)

	FilterDoNotTriggerMercury(opts *bind.FilterOpts, bn []*big.Int, ts []*big.Int) (*LogTriggeredFeedLookupDoNotTriggerMercuryIterator, error)

	WatchDoNotTriggerMercury(opts *bind.WatchOpts, sink chan<- *LogTriggeredFeedLookupDoNotTriggerMercury, bn []*big.Int, ts []*big.Int) (event.Subscription, error)

	ParseDoNotTriggerMercury(log types.Log) (*LogTriggeredFeedLookupDoNotTriggerMercury, error)

	FilterPerformingLogTriggerUpkeep(opts *bind.FilterOpts, from []common.Address) (*LogTriggeredFeedLookupPerformingLogTriggerUpkeepIterator, error)

	WatchPerformingLogTriggerUpkeep(opts *bind.WatchOpts, sink chan<- *LogTriggeredFeedLookupPerformingLogTriggerUpkeep, from []common.Address) (event.Subscription, error)

	ParsePerformingLogTriggerUpkeep(log types.Log) (*LogTriggeredFeedLookupPerformingLogTriggerUpkeep, error)

	FilterTriggerMercury(opts *bind.FilterOpts, upkeepId []*big.Int, logBlockNumber []*big.Int) (*LogTriggeredFeedLookupTriggerMercuryIterator, error)

	WatchTriggerMercury(opts *bind.WatchOpts, sink chan<- *LogTriggeredFeedLookupTriggerMercury, upkeepId []*big.Int, logBlockNumber []*big.Int) (event.Subscription, error)

	ParseTriggerMercury(log types.Log) (*LogTriggeredFeedLookupTriggerMercury, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
