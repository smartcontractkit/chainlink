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
	ABI: "[{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_useArbitrumBlockNum\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"_verify\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"StreamsLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"orderId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"exchange\",\"type\":\"address\"}],\"name\":\"LimitOrderExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"orderId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"exchange\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"blob\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"verified\",\"type\":\"bytes\"}],\"name\":\"PerformingLogTriggerUpkeep\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"errCode\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkErrorHandler\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"topics\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLog\",\"name\":\"log\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"checkLog\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feedsHex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParam\",\"type\":\"string\"}],\"name\":\"setFeedParamKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"newFeeds\",\"type\":\"string[]\"}],\"name\":\"setFeedsHex\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"timeParam\",\"type\":\"string\"}],\"name\":\"setTimeParamKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"start\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbitrumBlockNum\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x610120604052604260a08181526080918291906200180660c03990526200002a9060019081620000e8565b506040805180820190915260098152680cccacac892c890caf60bb1b60208201526002906200005a908262000264565b5060408051808201909152600b81526a313637b1b5a73ab6b132b960a91b60208201526003906200008c908262000264565b503480156200009a57600080fd5b506040516200184838038062001848833981016040819052620000bd9162000346565b6000805461ffff191692151561ff00191692909217610100911515919091021781556004556200037e565b82805482825590600052602060002090810192821562000133579160200282015b8281111562000133578251829062000122908262000264565b509160200191906001019062000109565b506200014192915062000145565b5090565b80821115620001415760006200015c828262000166565b5060010162000145565b5080546200017490620001d5565b6000825580601f1062000185575050565b601f016020900490600052602060002090810190620001a59190620001a8565b50565b5b80821115620001415760008155600101620001a9565b634e487b7160e01b600052604160045260246000fd5b600181811c90821680620001ea57607f821691505b6020821081036200020b57634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200025f57600081815260208120601f850160051c810160208610156200023a5750805b601f850160051c820191505b818110156200025b5782815560010162000246565b5050505b505050565b81516001600160401b03811115620002805762000280620001bf565b6200029881620002918454620001d5565b8462000211565b602080601f831160018114620002d05760008415620002b75750858301515b600019600386901b1c1916600185901b1785556200025b565b600085815260208120601f198616915b828110156200030157888601518255948401946001909101908401620002e0565b5085821015620003205787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b805180151581146200034157600080fd5b919050565b600080604083850312156200035a57600080fd5b620003658362000330565b9150620003756020840162000330565b90509250929050565b611478806200038e6000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c80639525d5741161008c578063afb28d1f11610066578063afb28d1f146101e4578063be9a6555146101ec578063c98f10b0146101f4578063fc735e99146101fc57600080fd5b80639525d5741461019b5780639d6f1cc7146101ae578063a4dc5e33146101ce57600080fd5b80634b56a42e116100c85780634b56a42e14610141578063601d5a711461015457806361bc221a14610167578063642f6cef1461017e57600080fd5b806305e25131146100ef57806340691db4146101045780634585e33b1461012e575b600080fd5b6101026100fd366004610ae9565b61020e565b005b610117610112366004610b9a565b610225565b604051610125929190610c75565b60405180910390f35b61010261013c366004610c98565b6104fb565b61011761014f366004610d0a565b6106f9565b610102610162366004610dc7565b61074f565b61017060045481565b604051908152602001610125565b60005461018b9060ff1681565b6040519015158152602001610125565b6101026101a9366004610dc7565b61075b565b6101c16101bc366004610dfc565b610767565b6040516101259190610e15565b6101176101dc366004610e2f565b600192909150565b6101c1610813565b610102610820565b6101c1610853565b60005461018b90610100900460ff1681565b80516102219060019060208401906108e6565b5050565b600060606000610233610860565b90507fd1ffe9e45581c11d7d9f2ed5f75217cd4be9f8b7eee6af0f6d03f46de53956cd61026360c0870187610e7e565b600081811061027457610274610ee6565b905060200201350361047357600061028f60c0870187610e7e565b60018181106102a0576102a0610ee6565b905060200201356040516020016102b991815260200190565b60405160208183030381529060405290506000818060200190518101906102e09190610f15565b905060006102f160c0890189610e7e565b600281811061030257610302610ee6565b9050602002013560405160200161031b91815260200190565b60405160208183030381529060405290506000818060200190518101906103429190610f15565b9050600061035360c08b018b610e7e565b600381811061036457610364610ee6565b9050602002013560405160200161037d91815260200190565b60405160208183030381529060405290506000818060200190518101906103a49190610f57565b604080516020810188905290810185905273ffffffffffffffffffffffffffffffffffffffff821660608201527fd1ffe9e45581c11d7d9f2ed5f75217cd4be9f8b7eee6af0f6d03f46de53956cd60808201529091506002906001906003908a9060a001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527ff055e4a200000000000000000000000000000000000000000000000000000000825261046a9594939291600401611060565b60405180910390fd5b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f636f756c64206e6f742066696e64206d61746368696e67206576656e7420736960448201527f6700000000000000000000000000000000000000000000000000000000000000606482015260840161046a565b60008061050a83850185610d0a565b91509150600080600080848060200190518101906105289190611123565b6040805160208101909152600080825254949850929650909450925090610100900460ff1615610621577309dff56a4ff44e0f4436260a04f5cfa65636a48173ffffffffffffffffffffffffffffffffffffffff16638e760afe8860008151811061059557610595610ee6565b60200260200101516040518263ffffffff1660e01b81526004016105b99190610e15565b6000604051808303816000875af11580156105d8573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261061e9190810190611160565b90505b60045461062f9060016111d7565b6004557f2e00161baa7e3ee28260d12a08ade832b4160748111950f092fc0b921ac6a933820161068b576040516000906064906001907fd1ffe9e45581c11d7d9f2ed5f75217cd4be9f8b7eee6af0f6d03f46de53956cd908490a45b327f299a03817e683a32b21e29e3ae3c31f1c9c773f7d532836d116b62a9281fbc9d8686866106b8610860565b8c6000815181106106cb576106cb610ee6565b6020026020010151876040516106e696959493929190611217565b60405180910390a2505050505050505050565b6000606060008484604051602001610712929190611277565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052600193509150505b9250929050565b60036102218282611351565b60026102218282611351565b6001818154811061077757600080fd5b90600052602060002001600091509050805461079290610f72565b80601f01602080910402602001604051908101604052809291908181526020018280546107be90610f72565b801561080b5780601f106107e05761010080835404028352916020019161080b565b820191906000526020600020905b8154815290600101906020018083116107ee57829003601f168201915b505050505081565b6002805461079290610f72565b6040516000906064906001907fd1ffe9e45581c11d7d9f2ed5f75217cd4be9f8b7eee6af0f6d03f46de53956cd908490a4565b6003805461079290610f72565b6000805460ff16156108e157606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156108b8573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108dc9190610f15565b905090565b504390565b82805482825590600052602060002090810192821561092c579160200282015b8281111561092c578251829061091c9082611351565b5091602001919060010190610906565b5061093892915061093c565b5090565b808211156109385760006109508282610959565b5060010161093c565b50805461096590610f72565b6000825580601f10610975575050565b601f0160209004906000526020600020908101906109939190610996565b50565b5b808211156109385760008155600101610997565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610a2157610a216109ab565b604052919050565b600067ffffffffffffffff821115610a4357610a436109ab565b5060051b60200190565b600067ffffffffffffffff821115610a6757610a676109ab565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f830112610aa457600080fd5b8135610ab7610ab282610a4d565b6109da565b818152846020838601011115610acc57600080fd5b816020850160208301376000918101602001919091529392505050565b60006020808385031215610afc57600080fd5b823567ffffffffffffffff80821115610b1457600080fd5b818501915085601f830112610b2857600080fd5b8135610b36610ab282610a29565b81815260059190911b83018401908481019088831115610b5557600080fd5b8585015b83811015610b8d57803585811115610b715760008081fd5b610b7f8b89838a0101610a93565b845250918601918601610b59565b5098975050505050505050565b60008060408385031215610bad57600080fd5b823567ffffffffffffffff80821115610bc557600080fd5b908401906101008287031215610bda57600080fd5b90925060208401359080821115610bf057600080fd5b50610bfd85828601610a93565b9150509250929050565b60005b83811015610c22578181015183820152602001610c0a565b50506000910152565b60008151808452610c43816020860160208601610c07565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8215158152604060208201526000610c906040830184610c2b565b949350505050565b60008060208385031215610cab57600080fd5b823567ffffffffffffffff80821115610cc357600080fd5b818501915085601f830112610cd757600080fd5b813581811115610ce657600080fd5b866020828501011115610cf857600080fd5b60209290920196919550909350505050565b60008060408385031215610d1d57600080fd5b823567ffffffffffffffff80821115610d3557600080fd5b818501915085601f830112610d4957600080fd5b81356020610d59610ab283610a29565b82815260059290921b84018101918181019089841115610d7857600080fd5b8286015b84811015610db057803586811115610d945760008081fd5b610da28c86838b0101610a93565b845250918301918301610d7c565b5096505086013592505080821115610bf057600080fd5b600060208284031215610dd957600080fd5b813567ffffffffffffffff811115610df057600080fd5b610c9084828501610a93565b600060208284031215610e0e57600080fd5b5035919050565b602081526000610e286020830184610c2b565b9392505050565b60008060408385031215610e4257600080fd5b823563ffffffff81168114610e5657600080fd5b9150602083013567ffffffffffffffff811115610e7257600080fd5b610bfd85828601610a93565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610eb357600080fd5b83018035915067ffffffffffffffff821115610ece57600080fd5b6020019150600581901b360382131561074857600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600060208284031215610f2757600080fd5b5051919050565b805173ffffffffffffffffffffffffffffffffffffffff81168114610f5257600080fd5b919050565b600060208284031215610f6957600080fd5b610e2882610f2e565b600181811c90821680610f8657607f821691505b602082108103610fbf577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60008154610fd281610f72565b808552602060018381168015610fef576001811461102757611055565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b8901019550611055565b866000528260002060005b8581101561104d5781548a8201860152908301908401611032565b890184019650505b505050505092915050565b60a08152600061107360a0830188610fc5565b6020838203818501528188548084528284019150828160051b8501018a6000528360002060005b838110156110e5577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08784030185526110d38383610fc5565b9486019492506001918201910161109a565b505086810360408801526110f9818b610fc5565b94505050505084606084015282810360808401526111178185610c2b565b98975050505050505050565b6000806000806080858703121561113957600080fd5b845193506020850151925061115060408601610f2e565b6060959095015193969295505050565b60006020828403121561117257600080fd5b815167ffffffffffffffff81111561118957600080fd5b8201601f8101841361119a57600080fd5b80516111a8610ab282610a4d565b8181528560208385010111156111bd57600080fd5b6111ce826020830160208601610c07565b95945050505050565b80820180821115611211577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b92915050565b86815285602082015273ffffffffffffffffffffffffffffffffffffffff8516604082015283606082015260c06080820152600061125860c0830185610c2b565b82810360a084015261126a8185610c2b565b9998505050505050505050565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b838110156112ec577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa08887030185526112da868351610c2b565b955093820193908201906001016112a0565b5050858403818701525050506111ce8185610c2b565b601f82111561134c57600081815260208120601f850160051c810160208610156113295750805b601f850160051c820191505b8181101561134857828155600101611335565b5050505b505050565b815167ffffffffffffffff81111561136b5761136b6109ab565b61137f816113798454610f72565b84611302565b602080601f8311600181146113d2576000841561139c5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555611348565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561141f57888601518255948401946001909101908401611400565b508582101561145b57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b0190555056fea164736f6c6343000810000a307834353534343832643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030",
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
	return address, tx, &LogTriggeredStreamsLookup{address: address, abi: *parsed, LogTriggeredStreamsLookupCaller: LogTriggeredStreamsLookupCaller{contract: contract}, LogTriggeredStreamsLookupTransactor: LogTriggeredStreamsLookupTransactor{contract: contract}, LogTriggeredStreamsLookupFilterer: LogTriggeredStreamsLookupFilterer{contract: contract}}, nil
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

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupCaller) CheckErrorHandler(opts *bind.CallOpts, errCode uint32, extraData []byte) (CheckErrorHandler,

	error) {
	var out []interface{}
	err := _LogTriggeredStreamsLookup.contract.Call(opts, &out, "checkErrorHandler", errCode, extraData)

	outstruct := new(CheckErrorHandler)
	if err != nil {
		return *outstruct, err
	}

	outstruct.UpkeepNeeded = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.PerformData = *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return *outstruct, err

}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupSession) CheckErrorHandler(errCode uint32, extraData []byte) (CheckErrorHandler,

	error) {
	return _LogTriggeredStreamsLookup.Contract.CheckErrorHandler(&_LogTriggeredStreamsLookup.CallOpts, errCode, extraData)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupCallerSession) CheckErrorHandler(errCode uint32, extraData []byte) (CheckErrorHandler,

	error) {
	return _LogTriggeredStreamsLookup.Contract.CheckErrorHandler(&_LogTriggeredStreamsLookup.CallOpts, errCode, extraData)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LogTriggeredStreamsLookup.contract.Call(opts, &out, "counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupSession) Counter() (*big.Int, error) {
	return _LogTriggeredStreamsLookup.Contract.Counter(&_LogTriggeredStreamsLookup.CallOpts)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupCallerSession) Counter() (*big.Int, error) {
	return _LogTriggeredStreamsLookup.Contract.Counter(&_LogTriggeredStreamsLookup.CallOpts)
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

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupTransactor) Start(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LogTriggeredStreamsLookup.contract.Transact(opts, "start")
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupSession) Start() (*types.Transaction, error) {
	return _LogTriggeredStreamsLookup.Contract.Start(&_LogTriggeredStreamsLookup.TransactOpts)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupTransactorSession) Start() (*types.Transaction, error) {
	return _LogTriggeredStreamsLookup.Contract.Start(&_LogTriggeredStreamsLookup.TransactOpts)
}

type LogTriggeredStreamsLookupLimitOrderExecutedIterator struct {
	Event *LogTriggeredStreamsLookupLimitOrderExecuted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LogTriggeredStreamsLookupLimitOrderExecutedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LogTriggeredStreamsLookupLimitOrderExecuted)
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
		it.Event = new(LogTriggeredStreamsLookupLimitOrderExecuted)
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

func (it *LogTriggeredStreamsLookupLimitOrderExecutedIterator) Error() error {
	return it.fail
}

func (it *LogTriggeredStreamsLookupLimitOrderExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LogTriggeredStreamsLookupLimitOrderExecuted struct {
	OrderId  *big.Int
	Amount   *big.Int
	Exchange common.Address
	Raw      types.Log
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupFilterer) FilterLimitOrderExecuted(opts *bind.FilterOpts, orderId []*big.Int, amount []*big.Int, exchange []common.Address) (*LogTriggeredStreamsLookupLimitOrderExecutedIterator, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var exchangeRule []interface{}
	for _, exchangeItem := range exchange {
		exchangeRule = append(exchangeRule, exchangeItem)
	}

	logs, sub, err := _LogTriggeredStreamsLookup.contract.FilterLogs(opts, "LimitOrderExecuted", orderIdRule, amountRule, exchangeRule)
	if err != nil {
		return nil, err
	}
	return &LogTriggeredStreamsLookupLimitOrderExecutedIterator{contract: _LogTriggeredStreamsLookup.contract, event: "LimitOrderExecuted", logs: logs, sub: sub}, nil
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupFilterer) WatchLimitOrderExecuted(opts *bind.WatchOpts, sink chan<- *LogTriggeredStreamsLookupLimitOrderExecuted, orderId []*big.Int, amount []*big.Int, exchange []common.Address) (event.Subscription, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var exchangeRule []interface{}
	for _, exchangeItem := range exchange {
		exchangeRule = append(exchangeRule, exchangeItem)
	}

	logs, sub, err := _LogTriggeredStreamsLookup.contract.WatchLogs(opts, "LimitOrderExecuted", orderIdRule, amountRule, exchangeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LogTriggeredStreamsLookupLimitOrderExecuted)
				if err := _LogTriggeredStreamsLookup.contract.UnpackLog(event, "LimitOrderExecuted", log); err != nil {
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

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupFilterer) ParseLimitOrderExecuted(log types.Log) (*LogTriggeredStreamsLookupLimitOrderExecuted, error) {
	event := new(LogTriggeredStreamsLookupLimitOrderExecuted)
	if err := _LogTriggeredStreamsLookup.contract.UnpackLog(event, "LimitOrderExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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

type CheckErrorHandler struct {
	UpkeepNeeded bool
	PerformData  []byte
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookup) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _LogTriggeredStreamsLookup.abi.Events["LimitOrderExecuted"].ID:
		return _LogTriggeredStreamsLookup.ParseLimitOrderExecuted(log)
	case _LogTriggeredStreamsLookup.abi.Events["PerformingLogTriggerUpkeep"].ID:
		return _LogTriggeredStreamsLookup.ParsePerformingLogTriggerUpkeep(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (LogTriggeredStreamsLookupLimitOrderExecuted) Topic() common.Hash {
	return common.HexToHash("0xd1ffe9e45581c11d7d9f2ed5f75217cd4be9f8b7eee6af0f6d03f46de53956cd")
}

func (LogTriggeredStreamsLookupPerformingLogTriggerUpkeep) Topic() common.Hash {
	return common.HexToHash("0x299a03817e683a32b21e29e3ae3c31f1c9c773f7d532836d116b62a9281fbc9d")
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookup) Address() common.Address {
	return _LogTriggeredStreamsLookup.address
}

type LogTriggeredStreamsLookupInterface interface {
	CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (bool, []byte, error)

	CheckErrorHandler(opts *bind.CallOpts, errCode uint32, extraData []byte) (CheckErrorHandler,

		error)

	Counter(opts *bind.CallOpts) (*big.Int, error)

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

	Start(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterLimitOrderExecuted(opts *bind.FilterOpts, orderId []*big.Int, amount []*big.Int, exchange []common.Address) (*LogTriggeredStreamsLookupLimitOrderExecutedIterator, error)

	WatchLimitOrderExecuted(opts *bind.WatchOpts, sink chan<- *LogTriggeredStreamsLookupLimitOrderExecuted, orderId []*big.Int, amount []*big.Int, exchange []common.Address) (event.Subscription, error)

	ParseLimitOrderExecuted(log types.Log) (*LogTriggeredStreamsLookupLimitOrderExecuted, error)

	FilterPerformingLogTriggerUpkeep(opts *bind.FilterOpts, from []common.Address) (*LogTriggeredStreamsLookupPerformingLogTriggerUpkeepIterator, error)

	WatchPerformingLogTriggerUpkeep(opts *bind.WatchOpts, sink chan<- *LogTriggeredStreamsLookupPerformingLogTriggerUpkeep, from []common.Address) (event.Subscription, error)

	ParsePerformingLogTriggerUpkeep(log types.Log) (*LogTriggeredStreamsLookupPerformingLogTriggerUpkeep, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
