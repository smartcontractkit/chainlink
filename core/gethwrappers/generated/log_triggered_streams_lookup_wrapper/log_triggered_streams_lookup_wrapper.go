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
	ABI: "[{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_useArbitrumBlockNum\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"_verify\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"_checkErrReturnBool\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"StreamsLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"IgnoringErrorHandlerData\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"orderId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"exchange\",\"type\":\"address\"}],\"name\":\"LimitOrderExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"orderId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"exchange\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"blob\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"verified\",\"type\":\"bytes\"}],\"name\":\"PerformingLogTriggerUpkeep\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkErrReturnBool\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"errCode\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkErrorHandler\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"topics\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLog\",\"name\":\"log\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"checkLog\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feedsHex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParam\",\"type\":\"string\"}],\"name\":\"setFeedParamKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"newFeeds\",\"type\":\"string[]\"}],\"name\":\"setFeedsHex\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"timeParam\",\"type\":\"string\"}],\"name\":\"setTimeParamKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"start\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbitrumBlockNum\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x610120604052604260a08181526080918291906200188360c03990526200002a9060019081620000fc565b506040805180820190915260098152680cccacac892c890caf60bb1b60208201526002906200005a908262000278565b5060408051808201909152600b81526a313637b1b5a73ab6b132b960a91b60208201526003906200008c908262000278565b503480156200009a57600080fd5b50604051620018c5380380620018c5833981016040819052620000bd916200035a565b6000805461ffff191693151561ff00191693909317610100921515929092029190911782556005805460ff1916911515919091179055600455620003a4565b82805482825590600052602060002090810192821562000147579160200282015b8281111562000147578251829062000136908262000278565b50916020019190600101906200011d565b506200015592915062000159565b5090565b80821115620001555760006200017082826200017a565b5060010162000159565b5080546200018890620001e9565b6000825580601f1062000199575050565b601f016020900490600052602060002090810190620001b99190620001bc565b50565b5b80821115620001555760008155600101620001bd565b634e487b7160e01b600052604160045260246000fd5b600181811c90821680620001fe57607f821691505b6020821081036200021f57634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200027357600081815260208120601f850160051c810160208610156200024e5750805b601f850160051c820191505b818110156200026f578281556001016200025a565b5050505b505050565b81516001600160401b03811115620002945762000294620001d3565b620002ac81620002a58454620001e9565b8462000225565b602080601f831160018114620002e45760008415620002cb5750858301515b600019600386901b1c1916600185901b1785556200026f565b600085815260208120601f198616915b828110156200031557888601518255948401946001909101908401620002f4565b5085821015620003345787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b805180151581146200035557600080fd5b919050565b6000806000606084860312156200037057600080fd5b6200037b8462000344565b92506200038b6020850162000344565b91506200039b6040850162000344565b90509250925092565b6114cf80620003b46000396000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c806361bc221a11610097578063afb28d1f11610066578063afb28d1f146101f9578063be9a655514610201578063c98f10b014610209578063fc735e991461021157600080fd5b806361bc221a146101a2578063642f6cef146101b95780639525d574146101c65780639d6f1cc7146101d957600080fd5b806340691db4116100d357806340691db4146101565780634585e33b146101695780634b56a42e1461017c578063601d5a711461018f57600080fd5b806305e25131146100fa5780630fb172fb1461010f5780631d1477b714610139575b600080fd5b61010d610108366004610b52565b610223565b005b61012261011d366004610c03565b61023a565b604051610130929190610cb8565b60405180910390f35b6005546101469060ff1681565b6040519015158152602001610130565b610122610164366004610cdb565b61025a565b61010d610177366004610d3e565b610530565b61012261018a366004610db0565b610764565b61010d61019d366004610e6d565b6107b8565b6101ab60045481565b604051908152602001610130565b6000546101469060ff1681565b61010d6101d4366004610e6d565b6107c4565b6101ec6101e7366004610ea2565b6107d0565b6040516101309190610ebb565b6101ec61087c565b61010d610889565b6101ec6108bc565b60005461014690610100900460ff1681565b805161023690600190602084019061094f565b5050565b60055460408051600081526020810190915260ff909116905b9250929050565b6000606060006102686108c9565b90507fd1ffe9e45581c11d7d9f2ed5f75217cd4be9f8b7eee6af0f6d03f46de53956cd61029860c0870187610ed5565b60008181106102a9576102a9610f3d565b90506020020135036104a85760006102c460c0870187610ed5565b60018181106102d5576102d5610f3d565b905060200201356040516020016102ee91815260200190565b60405160208183030381529060405290506000818060200190518101906103159190610f6c565b9050600061032660c0890189610ed5565b600281811061033757610337610f3d565b9050602002013560405160200161035091815260200190565b60405160208183030381529060405290506000818060200190518101906103779190610f6c565b9050600061038860c08b018b610ed5565b600381811061039957610399610f3d565b905060200201356040516020016103b291815260200190565b60405160208183030381529060405290506000818060200190518101906103d99190610fae565b604080516020810188905290810185905273ffffffffffffffffffffffffffffffffffffffff821660608201527fd1ffe9e45581c11d7d9f2ed5f75217cd4be9f8b7eee6af0f6d03f46de53956cd60808201529091506002906001906003908a9060a001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527ff055e4a200000000000000000000000000000000000000000000000000000000825261049f95949392916004016110b7565b60405180910390fd5b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f636f756c64206e6f742066696e64206d61746368696e67206576656e7420736960448201527f6700000000000000000000000000000000000000000000000000000000000000606482015260840161049f565b6000819003610566576040517f4f11cf8b97a2a25842c70eb74501fde8edea68ca2884190c1e3d26c8850441fe90600090a15050565b60008061057583850185610db0565b9150915060008060008084806020019051810190610593919061117a565b6040805160208101909152600080825254949850929650909450925090610100900460ff161561068c577309dff56a4ff44e0f4436260a04f5cfa65636a48173ffffffffffffffffffffffffffffffffffffffff16638e760afe8860008151811061060057610600610f3d565b60200260200101516040518263ffffffff1660e01b81526004016106249190610ebb565b6000604051808303816000875af1158015610643573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261068991908101906111b7565b90505b60045461069a90600161122e565b6004557f2e00161baa7e3ee28260d12a08ade832b4160748111950f092fc0b921ac6a93382016106f6576040516000906064906001907fd1ffe9e45581c11d7d9f2ed5f75217cd4be9f8b7eee6af0f6d03f46de53956cd908490a45b327f299a03817e683a32b21e29e3ae3c31f1c9c773f7d532836d116b62a9281fbc9d8686866107236108c9565b8c60008151811061073657610736610f3d565b6020026020010151876040516107519695949392919061126e565b60405180910390a2505050505050505050565b600060606000848460405160200161077d9291906112ce565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529190526001969095509350505050565b600361023682826113a8565b600261023682826113a8565b600181815481106107e057600080fd5b9060005260206000200160009150905080546107fb90610fc9565b80601f016020809104026020016040519081016040528092919081815260200182805461082790610fc9565b80156108745780601f1061084957610100808354040283529160200191610874565b820191906000526020600020905b81548152906001019060200180831161085757829003601f168201915b505050505081565b600280546107fb90610fc9565b6040516000906064906001907fd1ffe9e45581c11d7d9f2ed5f75217cd4be9f8b7eee6af0f6d03f46de53956cd908490a4565b600380546107fb90610fc9565b6000805460ff161561094a57606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610921573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109459190610f6c565b905090565b504390565b828054828255906000526020600020908101928215610995579160200282015b82811115610995578251829061098590826113a8565b509160200191906001019061096f565b506109a19291506109a5565b5090565b808211156109a15760006109b982826109c2565b506001016109a5565b5080546109ce90610fc9565b6000825580601f106109de575050565b601f0160209004906000526020600020908101906109fc91906109ff565b50565b5b808211156109a15760008155600101610a00565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610a8a57610a8a610a14565b604052919050565b600067ffffffffffffffff821115610aac57610aac610a14565b5060051b60200190565b600067ffffffffffffffff821115610ad057610ad0610a14565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f830112610b0d57600080fd5b8135610b20610b1b82610ab6565b610a43565b818152846020838601011115610b3557600080fd5b816020850160208301376000918101602001919091529392505050565b60006020808385031215610b6557600080fd5b823567ffffffffffffffff80821115610b7d57600080fd5b818501915085601f830112610b9157600080fd5b8135610b9f610b1b82610a92565b81815260059190911b83018401908481019088831115610bbe57600080fd5b8585015b83811015610bf657803585811115610bda5760008081fd5b610be88b89838a0101610afc565b845250918601918601610bc2565b5098975050505050505050565b60008060408385031215610c1657600080fd5b82359150602083013567ffffffffffffffff811115610c3457600080fd5b610c4085828601610afc565b9150509250929050565b60005b83811015610c65578181015183820152602001610c4d565b50506000910152565b60008151808452610c86816020860160208601610c4a565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8215158152604060208201526000610cd36040830184610c6e565b949350505050565b60008060408385031215610cee57600080fd5b823567ffffffffffffffff80821115610d0657600080fd5b908401906101008287031215610d1b57600080fd5b90925060208401359080821115610d3157600080fd5b50610c4085828601610afc565b60008060208385031215610d5157600080fd5b823567ffffffffffffffff80821115610d6957600080fd5b818501915085601f830112610d7d57600080fd5b813581811115610d8c57600080fd5b866020828501011115610d9e57600080fd5b60209290920196919550909350505050565b60008060408385031215610dc357600080fd5b823567ffffffffffffffff80821115610ddb57600080fd5b818501915085601f830112610def57600080fd5b81356020610dff610b1b83610a92565b82815260059290921b84018101918181019089841115610e1e57600080fd5b8286015b84811015610e5657803586811115610e3a5760008081fd5b610e488c86838b0101610afc565b845250918301918301610e22565b5096505086013592505080821115610d3157600080fd5b600060208284031215610e7f57600080fd5b813567ffffffffffffffff811115610e9657600080fd5b610cd384828501610afc565b600060208284031215610eb457600080fd5b5035919050565b602081526000610ece6020830184610c6e565b9392505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610f0a57600080fd5b83018035915067ffffffffffffffff821115610f2557600080fd5b6020019150600581901b360382131561025357600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600060208284031215610f7e57600080fd5b5051919050565b805173ffffffffffffffffffffffffffffffffffffffff81168114610fa957600080fd5b919050565b600060208284031215610fc057600080fd5b610ece82610f85565b600181811c90821680610fdd57607f821691505b602082108103611016577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b6000815461102981610fc9565b808552602060018381168015611046576001811461107e576110ac565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b89010195506110ac565b866000528260002060005b858110156110a45781548a8201860152908301908401611089565b890184019650505b505050505092915050565b60a0815260006110ca60a083018861101c565b6020838203818501528188548084528284019150828160051b8501018a6000528360002060005b8381101561113c577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe087840301855261112a838361101c565b948601949250600191820191016110f1565b50508681036040880152611150818b61101c565b945050505050846060840152828103608084015261116e8185610c6e565b98975050505050505050565b6000806000806080858703121561119057600080fd5b84519350602085015192506111a760408601610f85565b6060959095015193969295505050565b6000602082840312156111c957600080fd5b815167ffffffffffffffff8111156111e057600080fd5b8201601f810184136111f157600080fd5b80516111ff610b1b82610ab6565b81815285602083850101111561121457600080fd5b611225826020830160208601610c4a565b95945050505050565b80820180821115611268577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b92915050565b86815285602082015273ffffffffffffffffffffffffffffffffffffffff8516604082015283606082015260c0608082015260006112af60c0830185610c6e565b82810360a08401526112c18185610c6e565b9998505050505050505050565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b83811015611343577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0888703018552611331868351610c6e565b955093820193908201906001016112f7565b5050858403818701525050506112258185610c6e565b601f8211156113a357600081815260208120601f850160051c810160208610156113805750805b601f850160051c820191505b8181101561139f5782815560010161138c565b5050505b505050565b815167ffffffffffffffff8111156113c2576113c2610a14565b6113d6816113d08454610fc9565b84611359565b602080601f83116001811461142957600084156113f35750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b17855561139f565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561147657888601518255948401946001909101908401611457565b50858210156114b257878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b0190555056fea164736f6c6343000810000a307834353534343832643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030",
}

var LogTriggeredStreamsLookupABI = LogTriggeredStreamsLookupMetaData.ABI

var LogTriggeredStreamsLookupBin = LogTriggeredStreamsLookupMetaData.Bin

func DeployLogTriggeredStreamsLookup(auth *bind.TransactOpts, backend bind.ContractBackend, _useArbitrumBlockNum bool, _verify bool, _checkErrReturnBool bool) (common.Address, *types.Transaction, *LogTriggeredStreamsLookup, error) {
	parsed, err := LogTriggeredStreamsLookupMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LogTriggeredStreamsLookupBin), backend, _useArbitrumBlockNum, _verify, _checkErrReturnBool)
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

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupCaller) CheckErrReturnBool(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _LogTriggeredStreamsLookup.contract.Call(opts, &out, "checkErrReturnBool")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupSession) CheckErrReturnBool() (bool, error) {
	return _LogTriggeredStreamsLookup.Contract.CheckErrReturnBool(&_LogTriggeredStreamsLookup.CallOpts)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupCallerSession) CheckErrReturnBool() (bool, error) {
	return _LogTriggeredStreamsLookup.Contract.CheckErrReturnBool(&_LogTriggeredStreamsLookup.CallOpts)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupCaller) CheckErrorHandler(opts *bind.CallOpts, errCode *big.Int, extraData []byte) (CheckErrorHandler,

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

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupSession) CheckErrorHandler(errCode *big.Int, extraData []byte) (CheckErrorHandler,

	error) {
	return _LogTriggeredStreamsLookup.Contract.CheckErrorHandler(&_LogTriggeredStreamsLookup.CallOpts, errCode, extraData)
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupCallerSession) CheckErrorHandler(errCode *big.Int, extraData []byte) (CheckErrorHandler,

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

type LogTriggeredStreamsLookupIgnoringErrorHandlerDataIterator struct {
	Event *LogTriggeredStreamsLookupIgnoringErrorHandlerData

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LogTriggeredStreamsLookupIgnoringErrorHandlerDataIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LogTriggeredStreamsLookupIgnoringErrorHandlerData)
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
		it.Event = new(LogTriggeredStreamsLookupIgnoringErrorHandlerData)
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

func (it *LogTriggeredStreamsLookupIgnoringErrorHandlerDataIterator) Error() error {
	return it.fail
}

func (it *LogTriggeredStreamsLookupIgnoringErrorHandlerDataIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LogTriggeredStreamsLookupIgnoringErrorHandlerData struct {
	Raw types.Log
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupFilterer) FilterIgnoringErrorHandlerData(opts *bind.FilterOpts) (*LogTriggeredStreamsLookupIgnoringErrorHandlerDataIterator, error) {

	logs, sub, err := _LogTriggeredStreamsLookup.contract.FilterLogs(opts, "IgnoringErrorHandlerData")
	if err != nil {
		return nil, err
	}
	return &LogTriggeredStreamsLookupIgnoringErrorHandlerDataIterator{contract: _LogTriggeredStreamsLookup.contract, event: "IgnoringErrorHandlerData", logs: logs, sub: sub}, nil
}

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupFilterer) WatchIgnoringErrorHandlerData(opts *bind.WatchOpts, sink chan<- *LogTriggeredStreamsLookupIgnoringErrorHandlerData) (event.Subscription, error) {

	logs, sub, err := _LogTriggeredStreamsLookup.contract.WatchLogs(opts, "IgnoringErrorHandlerData")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LogTriggeredStreamsLookupIgnoringErrorHandlerData)
				if err := _LogTriggeredStreamsLookup.contract.UnpackLog(event, "IgnoringErrorHandlerData", log); err != nil {
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

func (_LogTriggeredStreamsLookup *LogTriggeredStreamsLookupFilterer) ParseIgnoringErrorHandlerData(log types.Log) (*LogTriggeredStreamsLookupIgnoringErrorHandlerData, error) {
	event := new(LogTriggeredStreamsLookupIgnoringErrorHandlerData)
	if err := _LogTriggeredStreamsLookup.contract.UnpackLog(event, "IgnoringErrorHandlerData", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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
	case _LogTriggeredStreamsLookup.abi.Events["IgnoringErrorHandlerData"].ID:
		return _LogTriggeredStreamsLookup.ParseIgnoringErrorHandlerData(log)
	case _LogTriggeredStreamsLookup.abi.Events["LimitOrderExecuted"].ID:
		return _LogTriggeredStreamsLookup.ParseLimitOrderExecuted(log)
	case _LogTriggeredStreamsLookup.abi.Events["PerformingLogTriggerUpkeep"].ID:
		return _LogTriggeredStreamsLookup.ParsePerformingLogTriggerUpkeep(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (LogTriggeredStreamsLookupIgnoringErrorHandlerData) Topic() common.Hash {
	return common.HexToHash("0x4f11cf8b97a2a25842c70eb74501fde8edea68ca2884190c1e3d26c8850441fe")
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

	CheckErrReturnBool(opts *bind.CallOpts) (bool, error)

	CheckErrorHandler(opts *bind.CallOpts, errCode *big.Int, extraData []byte) (CheckErrorHandler,

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

	FilterIgnoringErrorHandlerData(opts *bind.FilterOpts) (*LogTriggeredStreamsLookupIgnoringErrorHandlerDataIterator, error)

	WatchIgnoringErrorHandlerData(opts *bind.WatchOpts, sink chan<- *LogTriggeredStreamsLookupIgnoringErrorHandlerData) (event.Subscription, error)

	ParseIgnoringErrorHandlerData(log types.Log) (*LogTriggeredStreamsLookupIgnoringErrorHandlerData, error)

	FilterLimitOrderExecuted(opts *bind.FilterOpts, orderId []*big.Int, amount []*big.Int, exchange []common.Address) (*LogTriggeredStreamsLookupLimitOrderExecutedIterator, error)

	WatchLimitOrderExecuted(opts *bind.WatchOpts, sink chan<- *LogTriggeredStreamsLookupLimitOrderExecuted, orderId []*big.Int, amount []*big.Int, exchange []common.Address) (event.Subscription, error)

	ParseLimitOrderExecuted(log types.Log) (*LogTriggeredStreamsLookupLimitOrderExecuted, error)

	FilterPerformingLogTriggerUpkeep(opts *bind.FilterOpts, from []common.Address) (*LogTriggeredStreamsLookupPerformingLogTriggerUpkeepIterator, error)

	WatchPerformingLogTriggerUpkeep(opts *bind.WatchOpts, sink chan<- *LogTriggeredStreamsLookupPerformingLogTriggerUpkeep, from []common.Address) (event.Subscription, error)

	ParsePerformingLogTriggerUpkeep(log types.Log) (*LogTriggeredStreamsLookupPerformingLogTriggerUpkeep, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
