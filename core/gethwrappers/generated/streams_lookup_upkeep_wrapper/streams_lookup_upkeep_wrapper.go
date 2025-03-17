// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package streams_lookup_upkeep_wrapper

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

var StreamsLookupUpkeepMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_useArbBlock\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"_staging\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"_verify\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"StreamsLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v0\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"verifiedV0\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"ed\",\"type\":\"bytes\"}],\"name\":\"MercuryPerformEvent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"callbackReturnBool\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"errCode\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkErrorHandler\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feeds\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"interval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"setCallbackReturnBool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_feeds\",\"type\":\"string[]\"}],\"name\":\"setFeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_timeParamKey\",\"type\":\"string\"}],\"name\":\"setParamKeys\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"setShouldRevertCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"shouldRevertCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"staging\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbBlock\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162001adb38038062001adb83398101604081905262000034916200020f565b60008581556001859055600281905560038190556004558215156080526040805180820190915260078152666665656449447360c81b60208201526006906200007e908262000312565b50604080518082019091526009815268074696d657374616d760bc1b6020820152600790620000ae908262000312565b50604051806020016040528060405180608001604052806042815260200162001a99604291399052620000e690600590600162000122565b506008805463ff000000199215156101000261ff00199415159490941661ffff19909116179290921716630100000017905550620003de915050565b8280548282559060005260206000209081019282156200016d579160200282015b828111156200016d57825182906200015c908262000312565b509160200191906001019062000143565b506200017b9291506200017f565b5090565b808211156200017b576000620001968282620001a0565b506001016200017f565b508054620001ae9062000283565b6000825580601f10620001bf575050565b601f016020900490600052602060002090810190620001df9190620001e2565b50565b5b808211156200017b5760008155600101620001e3565b805180151581146200020a57600080fd5b919050565b600080600080600060a086880312156200022857600080fd5b85519450602086015193506200024160408701620001f9565b92506200025160608701620001f9565b91506200026160808701620001f9565b90509295509295909350565b634e487b7160e01b600052604160045260246000fd5b600181811c908216806200029857607f821691505b602082108103620002b957634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200030d57600081815260208120601f850160051c81016020861015620002e85750805b601f850160051c820191505b818110156200030957828155600101620002f4565b5050505b505050565b81516001600160401b038111156200032e576200032e6200026d565b62000346816200033f845462000283565b84620002bf565b602080601f8311600181146200037e5760008415620003655750858301515b600019600386901b1c1916600185901b17855562000309565b600085815260208120601f198616915b82811015620003af578886015182559484019460019091019084016200038e565b5085821015620003ce5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b60805161168a6200040f60003960008181610325015281816103c6015281816109400152610aaf015261168a6000f3fe608060405234801561001057600080fd5b506004361061018d5760003560e01c80636250a13a116100e3578063947a36fb1161008c578063d826f88f11610066578063d826f88f1461037c578063d832d92f14610390578063fc735e991461039857600080fd5b8063947a36fb14610363578063afb28d1f1461036c578063c98f10b01461037457600080fd5b806386b728e2116100bd57806386b728e21461032057806386e330af14610347578063917d895f1461035a57600080fd5b80636250a13a146102f15780636e04ff0d146102fa5780638340507c1461030d57600080fd5b80634585e33b116101455780634bdb38621161011f5780634bdb38621461025b5780635b48391a146102a157806361bc221a146102e857600080fd5b80634585e33b146102135780634a5479f3146102285780634b56a42e1461024857600080fd5b8063102d538b11610176578063102d538b146101db5780631d1970b7146101ef5780632cb15864146101fc57600080fd5b806302be021f146101925780630fb172fb146101ba575b600080fd5b6008546101a59062010000900460ff1681565b60405190151581526020015b60405180910390f35b6101cd6101c8366004610d5d565b6103aa565b6040516101b1929190610e12565b6008546101a5906301000000900460ff1681565b6008546101a59060ff1681565b61020560035481565b6040519081526020016101b1565b610226610221366004610e35565b6103c2565b005b61023b610236366004610ea7565b6106ef565b6040516101b19190610ec0565b6101cd610256366004610efe565b61079b565b610226610269366004610fc8565b6008805491151562010000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffff909216919091179055565b6102266102af366004610fc8565b600880549115156301000000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffff909216919091179055565b61020560045481565b61020560005481565b6101cd610308366004610e35565b610874565b61022661031b366004610fea565b610a4a565b6101a57f000000000000000000000000000000000000000000000000000000000000000081565b610226610355366004611037565b610a68565b61020560025481565b61020560015481565b61023b610a7f565b61023b610a8c565b610226600060028190556003819055600455565b6101a5610a99565b6008546101a590610100900460ff1681565b604080516000808252602082019092525b9250929050565b60007f00000000000000000000000000000000000000000000000000000000000000001561046157606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610436573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061045a91906110e8565b9050610464565b50435b6003546000036104745760038190555b60008061048384860186610efe565b6002859055600454919350915061049b906001611130565b600455604080516020808201835260008083528351918201909352918252600854909190610100900460ff16156106795760085460ff16156105aa577360448b880c9f3b501af3f343da9284148bd7d77c73ffffffffffffffffffffffffffffffffffffffff16638e760afe8560008151811061051a5761051a611149565b60200260200101516040518263ffffffff1660e01b815260040161053e9190610ec0565b6000604051808303816000875af115801561055d573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526105a39190810190611178565b9150610679565b7309dff56a4ff44e0f4436260a04f5cfa65636a48173ffffffffffffffffffffffffffffffffffffffff16638e760afe856000815181106105ed576105ed611149565b60200260200101516040518263ffffffff1660e01b81526004016106119190610ec0565b6000604051808303816000875af1158015610630573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526106769190810190611178565b91505b843373ffffffffffffffffffffffffffffffffffffffff167ff0f72c0b235fc8687d6a67c02ca543473a3cef8a18b48490f10e475a8dda1390866000815181106106c5576106c5611149565b602002602001015185876040516106de939291906111ef565b60405180910390a350505050505050565b600581815481106106ff57600080fd5b90600052602060002001600091509050805461071a90611232565b80601f016020809104026020016040519081016040528092919081815260200182805461074690611232565b80156107935780601f1061076857610100808354040283529160200191610793565b820191906000526020600020905b81548152906001019060200180831161077657829003601f168201915b505050505081565b60085460009060609062010000900460ff1615610819576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f73686f756c6452657665727443616c6c6261636b20697320747275650000000060448201526064015b60405180910390fd5b6000848460405160200161082e929190611285565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529190526008546301000000900460ff16969095509350505050565b60006060610880610a99565b6108cc576000848481818080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509597509195506103bb945050505050565b6040517f666565644964486578000000000000000000000000000000000000000000000060208201526000906029016040516020818303038152906040528051906020012060066040516020016109239190611310565b60405160208183030381529060405280519060200120036109e2577f0000000000000000000000000000000000000000000000000000000000000000156109db57606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156109b0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109d491906110e8565b90506109e5565b50436109e5565b50425b604080516c6400000000000000000000000060208201528151601481830301815260348201928390527ff055e4a2000000000000000000000000000000000000000000000000000000009092526108109160069160059160079186919060380161143f565b6006610a568382611550565b506007610a638282611550565b505050565b8051610a7b906005906020840190610b7e565b5050565b6006805461071a90611232565b6007805461071a90611232565b6000600354600003610aab5750600190565b60007f000000000000000000000000000000000000000000000000000000000000000015610b4a57606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610b1f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b4391906110e8565b9050610b4d565b50435b600054600354610b5d908361166a565b108015610b785750600154600254610b75908361166a565b10155b91505090565b828054828255906000526020600020908101928215610bc4579160200282015b82811115610bc45782518290610bb49082611550565b5091602001919060010190610b9e565b50610bd0929150610bd4565b5090565b80821115610bd0576000610be88282610bf1565b50600101610bd4565b508054610bfd90611232565b6000825580601f10610c0d575050565b601f016020900490600052602060002090810190610c2b9190610c2e565b50565b5b80821115610bd05760008155600101610c2f565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610cb957610cb9610c43565b604052919050565b600067ffffffffffffffff821115610cdb57610cdb610c43565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f830112610d1857600080fd5b8135610d2b610d2682610cc1565b610c72565b818152846020838601011115610d4057600080fd5b816020850160208301376000918101602001919091529392505050565b60008060408385031215610d7057600080fd5b82359150602083013567ffffffffffffffff811115610d8e57600080fd5b610d9a85828601610d07565b9150509250929050565b60005b83811015610dbf578181015183820152602001610da7565b50506000910152565b60008151808452610de0816020860160208601610da4565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8215158152604060208201526000610e2d6040830184610dc8565b949350505050565b60008060208385031215610e4857600080fd5b823567ffffffffffffffff80821115610e6057600080fd5b818501915085601f830112610e7457600080fd5b813581811115610e8357600080fd5b866020828501011115610e9557600080fd5b60209290920196919550909350505050565b600060208284031215610eb957600080fd5b5035919050565b602081526000610ed36020830184610dc8565b9392505050565b600067ffffffffffffffff821115610ef457610ef4610c43565b5060051b60200190565b60008060408385031215610f1157600080fd5b823567ffffffffffffffff80821115610f2957600080fd5b818501915085601f830112610f3d57600080fd5b81356020610f4d610d2683610eda565b82815260059290921b84018101918181019089841115610f6c57600080fd5b8286015b84811015610fa457803586811115610f885760008081fd5b610f968c86838b0101610d07565b845250918301918301610f70565b5096505086013592505080821115610fbb57600080fd5b50610d9a85828601610d07565b600060208284031215610fda57600080fd5b81358015158114610ed357600080fd5b60008060408385031215610ffd57600080fd5b823567ffffffffffffffff8082111561101557600080fd5b61102186838701610d07565b93506020850135915080821115610fbb57600080fd5b6000602080838503121561104a57600080fd5b823567ffffffffffffffff8082111561106257600080fd5b818501915085601f83011261107657600080fd5b8135611084610d2682610eda565b81815260059190911b830184019084810190888311156110a357600080fd5b8585015b838110156110db578035858111156110bf5760008081fd5b6110cd8b89838a0101610d07565b8452509186019186016110a7565b5098975050505050505050565b6000602082840312156110fa57600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082018082111561114357611143611101565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60006020828403121561118a57600080fd5b815167ffffffffffffffff8111156111a157600080fd5b8201601f810184136111b257600080fd5b80516111c0610d2682610cc1565b8181528560208385010111156111d557600080fd5b6111e6826020830160208601610da4565b95945050505050565b6060815260006112026060830186610dc8565b82810360208401526112148186610dc8565b905082810360408401526112288185610dc8565b9695505050505050565b600181811c9082168061124657607f821691505b60208210810361127f577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b838110156112fa577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa08887030185526112e8868351610dc8565b955093820193908201906001016112ae565b5050858403818701525050506111e68185610dc8565b600080835461131e81611232565b60018281168015611336576001811461136957611398565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0084168752821515830287019450611398565b8760005260208060002060005b8581101561138f5781548a820152908401908201611376565b50505082870194505b50929695505050505050565b600081546113b181611232565b8085526020600183811680156113ce576001811461140657611434565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b8901019550611434565b866000528260002060005b8581101561142c5781548a8201860152908301908401611411565b890184019650505b505050505092915050565b60a08152600061145260a08301886113a4565b6020838203818501528188548084528284019150828160051b8501018a6000528360002060005b838110156114c4577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08784030185526114b283836113a4565b94860194925060019182019101611479565b505086810360408801526114d8818b6113a4565b94505050505084606084015282810360808401526114f68185610dc8565b98975050505050505050565b601f821115610a6357600081815260208120601f850160051c810160208610156115295750805b601f850160051c820191505b8181101561154857828155600101611535565b505050505050565b815167ffffffffffffffff81111561156a5761156a610c43565b61157e816115788454611232565b84611502565b602080601f8311600181146115d1576000841561159b5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555611548565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561161e578886015182559484019460019091019084016115ff565b508582101561165a57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b818103818111156111435761114361110156fea164736f6c6343000810000a307830303032386339313564366166306664363662626132643066633934303532323662636138643638303633333331323161376439383332313033643135363363",
}

var StreamsLookupUpkeepABI = StreamsLookupUpkeepMetaData.ABI

var StreamsLookupUpkeepBin = StreamsLookupUpkeepMetaData.Bin

func DeployStreamsLookupUpkeep(auth *bind.TransactOpts, backend bind.ContractBackend, _testRange *big.Int, _interval *big.Int, _useArbBlock bool, _staging bool, _verify bool) (common.Address, *types.Transaction, *StreamsLookupUpkeep, error) {
	parsed, err := StreamsLookupUpkeepMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StreamsLookupUpkeepBin), backend, _testRange, _interval, _useArbBlock, _staging, _verify)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StreamsLookupUpkeep{address: address, abi: *parsed, StreamsLookupUpkeepCaller: StreamsLookupUpkeepCaller{contract: contract}, StreamsLookupUpkeepTransactor: StreamsLookupUpkeepTransactor{contract: contract}, StreamsLookupUpkeepFilterer: StreamsLookupUpkeepFilterer{contract: contract}}, nil
}

type StreamsLookupUpkeep struct {
	address common.Address
	abi     abi.ABI
	StreamsLookupUpkeepCaller
	StreamsLookupUpkeepTransactor
	StreamsLookupUpkeepFilterer
}

type StreamsLookupUpkeepCaller struct {
	contract *bind.BoundContract
}

type StreamsLookupUpkeepTransactor struct {
	contract *bind.BoundContract
}

type StreamsLookupUpkeepFilterer struct {
	contract *bind.BoundContract
}

type StreamsLookupUpkeepSession struct {
	Contract     *StreamsLookupUpkeep
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type StreamsLookupUpkeepCallerSession struct {
	Contract *StreamsLookupUpkeepCaller
	CallOpts bind.CallOpts
}

type StreamsLookupUpkeepTransactorSession struct {
	Contract     *StreamsLookupUpkeepTransactor
	TransactOpts bind.TransactOpts
}

type StreamsLookupUpkeepRaw struct {
	Contract *StreamsLookupUpkeep
}

type StreamsLookupUpkeepCallerRaw struct {
	Contract *StreamsLookupUpkeepCaller
}

type StreamsLookupUpkeepTransactorRaw struct {
	Contract *StreamsLookupUpkeepTransactor
}

func NewStreamsLookupUpkeep(address common.Address, backend bind.ContractBackend) (*StreamsLookupUpkeep, error) {
	abi, err := abi.JSON(strings.NewReader(StreamsLookupUpkeepABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindStreamsLookupUpkeep(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StreamsLookupUpkeep{address: address, abi: abi, StreamsLookupUpkeepCaller: StreamsLookupUpkeepCaller{contract: contract}, StreamsLookupUpkeepTransactor: StreamsLookupUpkeepTransactor{contract: contract}, StreamsLookupUpkeepFilterer: StreamsLookupUpkeepFilterer{contract: contract}}, nil
}

func NewStreamsLookupUpkeepCaller(address common.Address, caller bind.ContractCaller) (*StreamsLookupUpkeepCaller, error) {
	contract, err := bindStreamsLookupUpkeep(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StreamsLookupUpkeepCaller{contract: contract}, nil
}

func NewStreamsLookupUpkeepTransactor(address common.Address, transactor bind.ContractTransactor) (*StreamsLookupUpkeepTransactor, error) {
	contract, err := bindStreamsLookupUpkeep(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StreamsLookupUpkeepTransactor{contract: contract}, nil
}

func NewStreamsLookupUpkeepFilterer(address common.Address, filterer bind.ContractFilterer) (*StreamsLookupUpkeepFilterer, error) {
	contract, err := bindStreamsLookupUpkeep(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StreamsLookupUpkeepFilterer{contract: contract}, nil
}

func bindStreamsLookupUpkeep(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StreamsLookupUpkeepMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StreamsLookupUpkeep.Contract.StreamsLookupUpkeepCaller.contract.Call(opts, result, method, params...)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.Contract.StreamsLookupUpkeepTransactor.contract.Transfer(opts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.Contract.StreamsLookupUpkeepTransactor.contract.Transact(opts, method, params...)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StreamsLookupUpkeep.Contract.contract.Call(opts, result, method, params...)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.Contract.contract.Transfer(opts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.Contract.contract.Transact(opts, method, params...)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCaller) CallbackReturnBool(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _StreamsLookupUpkeep.contract.Call(opts, &out, "callbackReturnBool")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) CallbackReturnBool() (bool, error) {
	return _StreamsLookupUpkeep.Contract.CallbackReturnBool(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCallerSession) CallbackReturnBool() (bool, error) {
	return _StreamsLookupUpkeep.Contract.CallbackReturnBool(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCaller) CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (bool, []byte, error) {
	var out []interface{}
	err := _StreamsLookupUpkeep.contract.Call(opts, &out, "checkCallback", values, extraData)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) CheckCallback(values [][]byte, extraData []byte) (bool, []byte, error) {
	return _StreamsLookupUpkeep.Contract.CheckCallback(&_StreamsLookupUpkeep.CallOpts, values, extraData)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCallerSession) CheckCallback(values [][]byte, extraData []byte) (bool, []byte, error) {
	return _StreamsLookupUpkeep.Contract.CheckCallback(&_StreamsLookupUpkeep.CallOpts, values, extraData)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCaller) CheckErrorHandler(opts *bind.CallOpts, errCode *big.Int, extraData []byte) (CheckErrorHandler,

	error) {
	var out []interface{}
	err := _StreamsLookupUpkeep.contract.Call(opts, &out, "checkErrorHandler", errCode, extraData)

	outstruct := new(CheckErrorHandler)
	if err != nil {
		return *outstruct, err
	}

	outstruct.UpkeepNeeded = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.PerformData = *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return *outstruct, err

}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) CheckErrorHandler(errCode *big.Int, extraData []byte) (CheckErrorHandler,

	error) {
	return _StreamsLookupUpkeep.Contract.CheckErrorHandler(&_StreamsLookupUpkeep.CallOpts, errCode, extraData)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCallerSession) CheckErrorHandler(errCode *big.Int, extraData []byte) (CheckErrorHandler,

	error) {
	return _StreamsLookupUpkeep.Contract.CheckErrorHandler(&_StreamsLookupUpkeep.CallOpts, errCode, extraData)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCaller) CheckUpkeep(opts *bind.CallOpts, data []byte) (bool, []byte, error) {
	var out []interface{}
	err := _StreamsLookupUpkeep.contract.Call(opts, &out, "checkUpkeep", data)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) CheckUpkeep(data []byte) (bool, []byte, error) {
	return _StreamsLookupUpkeep.Contract.CheckUpkeep(&_StreamsLookupUpkeep.CallOpts, data)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCallerSession) CheckUpkeep(data []byte) (bool, []byte, error) {
	return _StreamsLookupUpkeep.Contract.CheckUpkeep(&_StreamsLookupUpkeep.CallOpts, data)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StreamsLookupUpkeep.contract.Call(opts, &out, "counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) Counter() (*big.Int, error) {
	return _StreamsLookupUpkeep.Contract.Counter(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCallerSession) Counter() (*big.Int, error) {
	return _StreamsLookupUpkeep.Contract.Counter(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCaller) Eligible(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _StreamsLookupUpkeep.contract.Call(opts, &out, "eligible")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) Eligible() (bool, error) {
	return _StreamsLookupUpkeep.Contract.Eligible(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCallerSession) Eligible() (bool, error) {
	return _StreamsLookupUpkeep.Contract.Eligible(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCaller) FeedParamKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _StreamsLookupUpkeep.contract.Call(opts, &out, "feedParamKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) FeedParamKey() (string, error) {
	return _StreamsLookupUpkeep.Contract.FeedParamKey(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCallerSession) FeedParamKey() (string, error) {
	return _StreamsLookupUpkeep.Contract.FeedParamKey(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCaller) Feeds(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _StreamsLookupUpkeep.contract.Call(opts, &out, "feeds", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) Feeds(arg0 *big.Int) (string, error) {
	return _StreamsLookupUpkeep.Contract.Feeds(&_StreamsLookupUpkeep.CallOpts, arg0)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCallerSession) Feeds(arg0 *big.Int) (string, error) {
	return _StreamsLookupUpkeep.Contract.Feeds(&_StreamsLookupUpkeep.CallOpts, arg0)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCaller) InitialBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StreamsLookupUpkeep.contract.Call(opts, &out, "initialBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) InitialBlock() (*big.Int, error) {
	return _StreamsLookupUpkeep.Contract.InitialBlock(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCallerSession) InitialBlock() (*big.Int, error) {
	return _StreamsLookupUpkeep.Contract.InitialBlock(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCaller) Interval(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StreamsLookupUpkeep.contract.Call(opts, &out, "interval")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) Interval() (*big.Int, error) {
	return _StreamsLookupUpkeep.Contract.Interval(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCallerSession) Interval() (*big.Int, error) {
	return _StreamsLookupUpkeep.Contract.Interval(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCaller) PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StreamsLookupUpkeep.contract.Call(opts, &out, "previousPerformBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) PreviousPerformBlock() (*big.Int, error) {
	return _StreamsLookupUpkeep.Contract.PreviousPerformBlock(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCallerSession) PreviousPerformBlock() (*big.Int, error) {
	return _StreamsLookupUpkeep.Contract.PreviousPerformBlock(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCaller) ShouldRevertCallback(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _StreamsLookupUpkeep.contract.Call(opts, &out, "shouldRevertCallback")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) ShouldRevertCallback() (bool, error) {
	return _StreamsLookupUpkeep.Contract.ShouldRevertCallback(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCallerSession) ShouldRevertCallback() (bool, error) {
	return _StreamsLookupUpkeep.Contract.ShouldRevertCallback(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCaller) Staging(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _StreamsLookupUpkeep.contract.Call(opts, &out, "staging")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) Staging() (bool, error) {
	return _StreamsLookupUpkeep.Contract.Staging(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCallerSession) Staging() (bool, error) {
	return _StreamsLookupUpkeep.Contract.Staging(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCaller) TestRange(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StreamsLookupUpkeep.contract.Call(opts, &out, "testRange")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) TestRange() (*big.Int, error) {
	return _StreamsLookupUpkeep.Contract.TestRange(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCallerSession) TestRange() (*big.Int, error) {
	return _StreamsLookupUpkeep.Contract.TestRange(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCaller) TimeParamKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _StreamsLookupUpkeep.contract.Call(opts, &out, "timeParamKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) TimeParamKey() (string, error) {
	return _StreamsLookupUpkeep.Contract.TimeParamKey(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCallerSession) TimeParamKey() (string, error) {
	return _StreamsLookupUpkeep.Contract.TimeParamKey(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCaller) UseArbBlock(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _StreamsLookupUpkeep.contract.Call(opts, &out, "useArbBlock")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) UseArbBlock() (bool, error) {
	return _StreamsLookupUpkeep.Contract.UseArbBlock(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCallerSession) UseArbBlock() (bool, error) {
	return _StreamsLookupUpkeep.Contract.UseArbBlock(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCaller) Verify(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _StreamsLookupUpkeep.contract.Call(opts, &out, "verify")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) Verify() (bool, error) {
	return _StreamsLookupUpkeep.Contract.Verify(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCallerSession) Verify() (bool, error) {
	return _StreamsLookupUpkeep.Contract.Verify(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.contract.Transact(opts, "performUpkeep", performData)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.Contract.PerformUpkeep(&_StreamsLookupUpkeep.TransactOpts, performData)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.Contract.PerformUpkeep(&_StreamsLookupUpkeep.TransactOpts, performData)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepTransactor) Reset(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.contract.Transact(opts, "reset")
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) Reset() (*types.Transaction, error) {
	return _StreamsLookupUpkeep.Contract.Reset(&_StreamsLookupUpkeep.TransactOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepTransactorSession) Reset() (*types.Transaction, error) {
	return _StreamsLookupUpkeep.Contract.Reset(&_StreamsLookupUpkeep.TransactOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepTransactor) SetCallbackReturnBool(opts *bind.TransactOpts, value bool) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.contract.Transact(opts, "setCallbackReturnBool", value)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) SetCallbackReturnBool(value bool) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.Contract.SetCallbackReturnBool(&_StreamsLookupUpkeep.TransactOpts, value)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepTransactorSession) SetCallbackReturnBool(value bool) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.Contract.SetCallbackReturnBool(&_StreamsLookupUpkeep.TransactOpts, value)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepTransactor) SetFeeds(opts *bind.TransactOpts, _feeds []string) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.contract.Transact(opts, "setFeeds", _feeds)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) SetFeeds(_feeds []string) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.Contract.SetFeeds(&_StreamsLookupUpkeep.TransactOpts, _feeds)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepTransactorSession) SetFeeds(_feeds []string) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.Contract.SetFeeds(&_StreamsLookupUpkeep.TransactOpts, _feeds)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepTransactor) SetParamKeys(opts *bind.TransactOpts, _feedParamKey string, _timeParamKey string) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.contract.Transact(opts, "setParamKeys", _feedParamKey, _timeParamKey)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) SetParamKeys(_feedParamKey string, _timeParamKey string) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.Contract.SetParamKeys(&_StreamsLookupUpkeep.TransactOpts, _feedParamKey, _timeParamKey)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepTransactorSession) SetParamKeys(_feedParamKey string, _timeParamKey string) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.Contract.SetParamKeys(&_StreamsLookupUpkeep.TransactOpts, _feedParamKey, _timeParamKey)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepTransactor) SetShouldRevertCallback(opts *bind.TransactOpts, value bool) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.contract.Transact(opts, "setShouldRevertCallback", value)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) SetShouldRevertCallback(value bool) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.Contract.SetShouldRevertCallback(&_StreamsLookupUpkeep.TransactOpts, value)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepTransactorSession) SetShouldRevertCallback(value bool) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.Contract.SetShouldRevertCallback(&_StreamsLookupUpkeep.TransactOpts, value)
}

type StreamsLookupUpkeepMercuryPerformEventIterator struct {
	Event *StreamsLookupUpkeepMercuryPerformEvent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StreamsLookupUpkeepMercuryPerformEventIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StreamsLookupUpkeepMercuryPerformEvent)
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
		it.Event = new(StreamsLookupUpkeepMercuryPerformEvent)
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

func (it *StreamsLookupUpkeepMercuryPerformEventIterator) Error() error {
	return it.fail
}

func (it *StreamsLookupUpkeepMercuryPerformEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StreamsLookupUpkeepMercuryPerformEvent struct {
	Sender      common.Address
	BlockNumber *big.Int
	V0          []byte
	VerifiedV0  []byte
	Ed          []byte
	Raw         types.Log
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepFilterer) FilterMercuryPerformEvent(opts *bind.FilterOpts, sender []common.Address, blockNumber []*big.Int) (*StreamsLookupUpkeepMercuryPerformEventIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var blockNumberRule []interface{}
	for _, blockNumberItem := range blockNumber {
		blockNumberRule = append(blockNumberRule, blockNumberItem)
	}

	logs, sub, err := _StreamsLookupUpkeep.contract.FilterLogs(opts, "MercuryPerformEvent", senderRule, blockNumberRule)
	if err != nil {
		return nil, err
	}
	return &StreamsLookupUpkeepMercuryPerformEventIterator{contract: _StreamsLookupUpkeep.contract, event: "MercuryPerformEvent", logs: logs, sub: sub}, nil
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepFilterer) WatchMercuryPerformEvent(opts *bind.WatchOpts, sink chan<- *StreamsLookupUpkeepMercuryPerformEvent, sender []common.Address, blockNumber []*big.Int) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var blockNumberRule []interface{}
	for _, blockNumberItem := range blockNumber {
		blockNumberRule = append(blockNumberRule, blockNumberItem)
	}

	logs, sub, err := _StreamsLookupUpkeep.contract.WatchLogs(opts, "MercuryPerformEvent", senderRule, blockNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StreamsLookupUpkeepMercuryPerformEvent)
				if err := _StreamsLookupUpkeep.contract.UnpackLog(event, "MercuryPerformEvent", log); err != nil {
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

func (_StreamsLookupUpkeep *StreamsLookupUpkeepFilterer) ParseMercuryPerformEvent(log types.Log) (*StreamsLookupUpkeepMercuryPerformEvent, error) {
	event := new(StreamsLookupUpkeepMercuryPerformEvent)
	if err := _StreamsLookupUpkeep.contract.UnpackLog(event, "MercuryPerformEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CheckErrorHandler struct {
	UpkeepNeeded bool
	PerformData  []byte
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeep) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _StreamsLookupUpkeep.abi.Events["MercuryPerformEvent"].ID:
		return _StreamsLookupUpkeep.ParseMercuryPerformEvent(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (StreamsLookupUpkeepMercuryPerformEvent) Topic() common.Hash {
	return common.HexToHash("0xf0f72c0b235fc8687d6a67c02ca543473a3cef8a18b48490f10e475a8dda1390")
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeep) Address() common.Address {
	return _StreamsLookupUpkeep.address
}

type StreamsLookupUpkeepInterface interface {
	CallbackReturnBool(opts *bind.CallOpts) (bool, error)

	CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (bool, []byte, error)

	CheckErrorHandler(opts *bind.CallOpts, errCode *big.Int, extraData []byte) (CheckErrorHandler,

		error)

	CheckUpkeep(opts *bind.CallOpts, data []byte) (bool, []byte, error)

	Counter(opts *bind.CallOpts) (*big.Int, error)

	Eligible(opts *bind.CallOpts) (bool, error)

	FeedParamKey(opts *bind.CallOpts) (string, error)

	Feeds(opts *bind.CallOpts, arg0 *big.Int) (string, error)

	InitialBlock(opts *bind.CallOpts) (*big.Int, error)

	Interval(opts *bind.CallOpts) (*big.Int, error)

	PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error)

	ShouldRevertCallback(opts *bind.CallOpts) (bool, error)

	Staging(opts *bind.CallOpts) (bool, error)

	TestRange(opts *bind.CallOpts) (*big.Int, error)

	TimeParamKey(opts *bind.CallOpts) (string, error)

	UseArbBlock(opts *bind.CallOpts) (bool, error)

	Verify(opts *bind.CallOpts) (bool, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	Reset(opts *bind.TransactOpts) (*types.Transaction, error)

	SetCallbackReturnBool(opts *bind.TransactOpts, value bool) (*types.Transaction, error)

	SetFeeds(opts *bind.TransactOpts, _feeds []string) (*types.Transaction, error)

	SetParamKeys(opts *bind.TransactOpts, _feedParamKey string, _timeParamKey string) (*types.Transaction, error)

	SetShouldRevertCallback(opts *bind.TransactOpts, value bool) (*types.Transaction, error)

	FilterMercuryPerformEvent(opts *bind.FilterOpts, sender []common.Address, blockNumber []*big.Int) (*StreamsLookupUpkeepMercuryPerformEventIterator, error)

	WatchMercuryPerformEvent(opts *bind.WatchOpts, sink chan<- *StreamsLookupUpkeepMercuryPerformEvent, sender []common.Address, blockNumber []*big.Int) (event.Subscription, error)

	ParseMercuryPerformEvent(log types.Log) (*StreamsLookupUpkeepMercuryPerformEvent, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
