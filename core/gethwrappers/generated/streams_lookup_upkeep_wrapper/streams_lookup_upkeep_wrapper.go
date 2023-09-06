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
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_useArbBlock\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"_staging\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"_verify\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"StreamsLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v0\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v1\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"verifiedV0\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"verifiedV1\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"ed\",\"type\":\"bytes\"}],\"name\":\"MercuryPerformEvent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"callbackReturnBool\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feeds\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"interval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"setCallbackReturnBool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_feeds\",\"type\":\"string[]\"}],\"name\":\"setFeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_timeParamKey\",\"type\":\"string\"}],\"name\":\"setParamKeys\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"setShouldRevertCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"shouldRevertCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"staging\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbBlock\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162001ca738038062001ca7833981016040819052620000349162000232565b60008581556001859055600281905560038190556004558215156080526040805180820190915260098152680cccacac892c890caf60bb1b602082015260069062000080908262000335565b5060408051808201909152600b81526a313637b1b5a73ab6b132b960a91b6020820152600790620000b2908262000335565b50604051806040016040528060405180608001604052806042815260200162001c2360429139815260200160405180608001604052806042815260200162001c656042913990526200010990600590600262000145565b506008805463ff000000199215156101000261ff00199415159490941661ffff1990911617929092171663010000001790555062000401915050565b82805482825590600052602060002090810192821562000190579160200282015b828111156200019057825182906200017f908262000335565b509160200191906001019062000166565b506200019e929150620001a2565b5090565b808211156200019e576000620001b98282620001c3565b50600101620001a2565b508054620001d190620002a6565b6000825580601f10620001e2575050565b601f01602090049060005260206000209081019062000202919062000205565b50565b5b808211156200019e576000815560010162000206565b805180151581146200022d57600080fd5b919050565b600080600080600060a086880312156200024b57600080fd5b855194506020860151935062000264604087016200021c565b925062000274606087016200021c565b915062000284608087016200021c565b90509295509295909350565b634e487b7160e01b600052604160045260246000fd5b600181811c90821680620002bb57607f821691505b602082108103620002dc57634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200033057600081815260208120601f850160051c810160208610156200030b5750805b601f850160051c820191505b818110156200032c5782815560010162000317565b5050505b505050565b81516001600160401b0381111562000351576200035162000290565b6200036981620003628454620002a6565b84620002e2565b602080601f831160018114620003a15760008415620003885750858301515b600019600386901b1c1916600185901b1785556200032c565b600085815260208120601f198616915b82811015620003d257888601518255948401946001909101908401620003b1565b5085821015620003f15787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6080516117f162000432600039600081816103070152818161039001528181610ac60152610c3501526117f16000f3fe608060405234801561001057600080fd5b50600436106101825760003560e01c80636e04ff0d116100d8578063947a36fb1161008c578063d826f88f11610066578063d826f88f1461035e578063d832d92f14610372578063fc735e991461037a57600080fd5b8063947a36fb14610345578063afb28d1f1461034e578063c98f10b01461035657600080fd5b806386b728e2116100bd57806386b728e21461030257806386e330af14610329578063917d895f1461033c57600080fd5b80636e04ff0d146102dc5780638340507c146102ef57600080fd5b80634a5479f31161013a5780635b48391a116101145780635b48391a1461028357806361bc221a146102ca5780636250a13a146102d357600080fd5b80634a5479f3146101fc5780634b56a42e1461021c5780634bdb38621461023d57600080fd5b80631d1970b71161016b5780631d1970b7146101c35780632cb15864146101d05780634585e33b146101e757600080fd5b806302be021f14610187578063102d538b146101af575b600080fd5b60085461019a9062010000900460ff1681565b60405190151581526020015b60405180910390f35b60085461019a906301000000900460ff1681565b60085461019a9060ff1681565b6101d960035481565b6040519081526020016101a6565b6101fa6101f5366004610dc9565b61038c565b005b61020f61020a366004610e3b565b610873565b6040516101a69190610ec2565b61022f61022a36600461101a565b61091f565b6040516101a69291906110ee565b6101fa61024b366004611111565b6008805491151562010000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffff909216919091179055565b6101fa610291366004611111565b600880549115156301000000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffff909216919091179055565b6101d960045481565b6101d960005481565b61022f6102ea366004610dc9565b6109fa565b6101fa6102fd366004611133565b610bd0565b61019a7f000000000000000000000000000000000000000000000000000000000000000081565b6101fa610337366004611180565b610bee565b6101d960025481565b6101d960015481565b61020f610c05565b61020f610c12565b6101fa600060028190556003819055600455565b61019a610c1f565b60085461019a90610100900460ff1681565b60007f00000000000000000000000000000000000000000000000000000000000000001561042b57606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610400573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104249190611231565b905061042e565b50435b60035460000361043e5760038190555b60008061044d8486018661101a565b60028590556004549193509150610465906001611279565b600455604080516020808201835260008083528351918201909352918252600854909190610100900460ff16156107df5760085460ff1615610642577360448b880c9f3b501af3f343da9284148bd7d77c73ffffffffffffffffffffffffffffffffffffffff16638e760afe856000815181106104e4576104e4611292565b60200260200101516040518263ffffffff1660e01b81526004016105089190610ec2565b6000604051808303816000875af1158015610527573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261056d91908101906112c1565b91507360448b880c9f3b501af3f343da9284148bd7d77c73ffffffffffffffffffffffffffffffffffffffff16638e760afe856001815181106105b2576105b2611292565b60200260200101516040518263ffffffff1660e01b81526004016105d69190610ec2565b6000604051808303816000875af11580156105f5573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261063b91908101906112c1565b90506107df565b7309dff56a4ff44e0f4436260a04f5cfa65636a48173ffffffffffffffffffffffffffffffffffffffff16638e760afe8560008151811061068557610685611292565b60200260200101516040518263ffffffff1660e01b81526004016106a99190610ec2565b6000604051808303816000875af11580156106c8573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261070e91908101906112c1565b91507309dff56a4ff44e0f4436260a04f5cfa65636a48173ffffffffffffffffffffffffffffffffffffffff16638e760afe8560018151811061075357610753611292565b60200260200101516040518263ffffffff1660e01b81526004016107779190610ec2565b6000604051808303816000875af1158015610796573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526107dc91908101906112c1565b90505b843373ffffffffffffffffffffffffffffffffffffffff167f1c85d6186f024e964616014c8247533455ec5129a5095711202292f8a7ea1d548660008151811061082b5761082b611292565b60200260200101518760018151811061084657610846611292565b6020026020010151868689604051610862959493929190611338565b60405180910390a350505050505050565b6005818154811061088357600080fd5b90600052602060002001600091509050805461089e906113a5565b80601f01602080910402602001604051908101604052809291908181526020018280546108ca906113a5565b80156109175780601f106108ec57610100808354040283529160200191610917565b820191906000526020600020905b8154815290600101906020018083116108fa57829003601f168201915b505050505081565b60085460009060609062010000900460ff161561099d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f73686f756c6452657665727443616c6c6261636b20697320747275650000000060448201526064015b60405180910390fd5b600084846040516020016109b29291906113f8565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529190526008546301000000900460ff1693509150505b9250929050565b60006060610a06610c1f565b610a52576000848481818080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509597509195506109f3945050505050565b6040517f66656564496448657800000000000000000000000000000000000000000000006020820152600090602901604051602081830303815290604052805190602001206007604051602001610aa99190611483565b6040516020818303038152906040528051906020012003610b68577f000000000000000000000000000000000000000000000000000000000000000015610b6157606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610b36573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b5a9190611231565b9050610b6b565b5043610b6b565b50425b604080516c6400000000000000000000000060208201528151601481830301815260348201928390527ff055e4a200000000000000000000000000000000000000000000000000000000909252610994916006916005916007918691906038016115b2565b6006610bdc83826116b7565b506007610be982826116b7565b505050565b8051610c01906005906020840190610d04565b5050565b6006805461089e906113a5565b6007805461089e906113a5565b6000600354600003610c315750600190565b60007f000000000000000000000000000000000000000000000000000000000000000015610cd057606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610ca5573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610cc99190611231565b9050610cd3565b50435b600054600354610ce390836117d1565b108015610cfe5750600154600254610cfb90836117d1565b10155b91505090565b828054828255906000526020600020908101928215610d4a579160200282015b82811115610d4a5782518290610d3a90826116b7565b5091602001919060010190610d24565b50610d56929150610d5a565b5090565b80821115610d56576000610d6e8282610d77565b50600101610d5a565b508054610d83906113a5565b6000825580601f10610d93575050565b601f016020900490600052602060002090810190610db19190610db4565b50565b5b80821115610d565760008155600101610db5565b60008060208385031215610ddc57600080fd5b823567ffffffffffffffff80821115610df457600080fd5b818501915085601f830112610e0857600080fd5b813581811115610e1757600080fd5b866020828501011115610e2957600080fd5b60209290920196919550909350505050565b600060208284031215610e4d57600080fd5b5035919050565b60005b83811015610e6f578181015183820152602001610e57565b50506000910152565b60008151808452610e90816020860160208601610e54565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000610ed56020830184610e78565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610f5257610f52610edc565b604052919050565b600067ffffffffffffffff821115610f7457610f74610edc565b5060051b60200190565b600067ffffffffffffffff821115610f9857610f98610edc565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f830112610fd557600080fd5b8135610fe8610fe382610f7e565b610f0b565b818152846020838601011115610ffd57600080fd5b816020850160208301376000918101602001919091529392505050565b6000806040838503121561102d57600080fd5b823567ffffffffffffffff8082111561104557600080fd5b818501915085601f83011261105957600080fd5b81356020611069610fe383610f5a565b82815260059290921b8401810191818101908984111561108857600080fd5b8286015b848110156110c0578035868111156110a45760008081fd5b6110b28c86838b0101610fc4565b84525091830191830161108c565b50965050860135925050808211156110d757600080fd5b506110e485828601610fc4565b9150509250929050565b82151581526040602082015260006111096040830184610e78565b949350505050565b60006020828403121561112357600080fd5b81358015158114610ed557600080fd5b6000806040838503121561114657600080fd5b823567ffffffffffffffff8082111561115e57600080fd5b61116a86838701610fc4565b935060208501359150808211156110d757600080fd5b6000602080838503121561119357600080fd5b823567ffffffffffffffff808211156111ab57600080fd5b818501915085601f8301126111bf57600080fd5b81356111cd610fe382610f5a565b81815260059190911b830184019084810190888311156111ec57600080fd5b8585015b83811015611224578035858111156112085760008081fd5b6112168b89838a0101610fc4565b8452509186019186016111f0565b5098975050505050505050565b60006020828403121561124357600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082018082111561128c5761128c61124a565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6000602082840312156112d357600080fd5b815167ffffffffffffffff8111156112ea57600080fd5b8201601f810184136112fb57600080fd5b8051611309610fe382610f7e565b81815285602083850101111561131e57600080fd5b61132f826020830160208601610e54565b95945050505050565b60a08152600061134b60a0830188610e78565b828103602084015261135d8188610e78565b905082810360408401526113718187610e78565b905082810360608401526113858186610e78565b905082810360808401526113998185610e78565b98975050505050505050565b600181811c908216806113b957607f821691505b6020821081036113f2577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b8381101561146d577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa088870301855261145b868351610e78565b95509382019390820190600101611421565b50508584038187015250505061132f8185610e78565b6000808354611491816113a5565b600182811680156114a957600181146114dc5761150b565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008416875282151583028701945061150b565b8760005260208060002060005b858110156115025781548a8201529084019082016114e9565b50505082870194505b50929695505050505050565b60008154611524816113a5565b8085526020600183811680156115415760018114611579576115a7565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b89010195506115a7565b866000528260002060005b8581101561159f5781548a8201860152908301908401611584565b890184019650505b505050505092915050565b60a0815260006115c560a0830188611517565b6020838203818501528188548084528284019150828160051b8501018a6000528360002060005b83811015611637577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08784030185526116258383611517565b948601949250600191820191016115ec565b5050868103604088015261164b818b611517565b94505050505084606084015282810360808401526113998185610e78565b601f821115610be957600081815260208120601f850160051c810160208610156116905750805b601f850160051c820191505b818110156116af5782815560010161169c565b505050505050565b815167ffffffffffffffff8111156116d1576116d1610edc565b6116e5816116df84546113a5565b84611669565b602080601f83116001811461173857600084156117025750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556116af565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561178557888601518255948401946001909101908401611766565b50858210156117c157878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b8181038181111561128c5761128c61124a56fea164736f6c6343000810000a307834353534343832643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030307834323534343332643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030",
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
	return address, tx, &StreamsLookupUpkeep{StreamsLookupUpkeepCaller: StreamsLookupUpkeepCaller{contract: contract}, StreamsLookupUpkeepTransactor: StreamsLookupUpkeepTransactor{contract: contract}, StreamsLookupUpkeepFilterer: StreamsLookupUpkeepFilterer{contract: contract}}, nil
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
	V1          []byte
	VerifiedV0  []byte
	VerifiedV1  []byte
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

func (_StreamsLookupUpkeep *StreamsLookupUpkeep) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _StreamsLookupUpkeep.abi.Events["MercuryPerformEvent"].ID:
		return _StreamsLookupUpkeep.ParseMercuryPerformEvent(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (StreamsLookupUpkeepMercuryPerformEvent) Topic() common.Hash {
	return common.HexToHash("0x1c85d6186f024e964616014c8247533455ec5129a5095711202292f8a7ea1d54")
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeep) Address() common.Address {
	return _StreamsLookupUpkeep.address
}

type StreamsLookupUpkeepInterface interface {
	CallbackReturnBool(opts *bind.CallOpts) (bool, error)

	CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (bool, []byte, error)

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
