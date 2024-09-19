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
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_useArbBlock\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"_staging\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"_verify\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"StreamsLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v0\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"verifiedV0\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"ed\",\"type\":\"bytes\"}],\"name\":\"MercuryPerformEvent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"callbackReturnBool\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"errCode\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkErrorHandler\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feeds\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"interval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"production_testnet_verifier_proxy\",\"outputs\":[{\"internalType\":\"contractIVerifierProxy\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"setCallbackReturnBool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_feeds\",\"type\":\"string[]\"}],\"name\":\"setFeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_timeParamKey\",\"type\":\"string\"}],\"name\":\"setParamKeys\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIVerifierProxy\",\"name\":\"proxy\",\"type\":\"address\"}],\"name\":\"setProductionTestnetVerifierProxy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"setShouldRevertCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_staging\",\"type\":\"bool\"}],\"name\":\"setStaging\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIVerifierProxy\",\"name\":\"proxy\",\"type\":\"address\"}],\"name\":\"setStagingTestnetVerifierProxy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_n\",\"type\":\"uint256\"}],\"name\":\"setVerifyNthReport\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"shouldRevertCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"staging\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"staging_testnet_verifier_proxy\",\"outputs\":[{\"internalType\":\"contractIVerifierProxy\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbBlock\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"verifyNthReport\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a060405260008054732ff010debc1297f19579b4246cad07bd24f2488a6001600160a01b031991821681179092556001805490911690911790553480156200004757600080fd5b5060405162001b7b38038062001b7b8339810160408190526200006a916200024a565b600285905560038490556000600481905560058190556006558215156080526040805180820190915260078152666665656449447360c81b6020820152600890620000b690826200034d565b5060408051808201909152600980825268074696d657374616d760bc1b602083015290620000e590826200034d565b5060405180602001604052806040518060c001604052806084815260200162001af76084913990526200011d9060079060016200015d565b50600a805463ff000000199215156101000261ff00199415159490941661ffff19909116179290921716630100000017905550506000600b555062000419565b828054828255906000526020600020908101928215620001a8579160200282015b82811115620001a857825182906200019790826200034d565b50916020019190600101906200017e565b50620001b6929150620001ba565b5090565b80821115620001b6576000620001d18282620001db565b50600101620001ba565b508054620001e990620002be565b6000825580601f10620001fa575050565b601f0160209004906000526020600020908101906200021a91906200021d565b50565b5b80821115620001b657600081556001016200021e565b805180151581146200024557600080fd5b919050565b600080600080600060a086880312156200026357600080fd5b85519450602086015193506200027c6040870162000234565b92506200028c6060870162000234565b91506200029c6080870162000234565b90509295509295909350565b634e487b7160e01b600052604160045260246000fd5b600181811c90821680620002d357607f821691505b602082108103620002f457634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200034857600081815260208120601f850160051c81016020861015620003235750805b601f850160051c820191505b8181101562000344578281556001016200032f565b5050505b505050565b81516001600160401b03811115620003695762000369620002a8565b62000381816200037a8454620002be565b84620002fa565b602080601f831160018114620003b95760008415620003a05750858301515b600019600386901b1c1916600185901b17855562000344565b600085815260208120601f198616915b82811015620003ea57888601518255948401946001909101908401620003c9565b5085821015620004095787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6080516116b462000443600039600081816104250152818161057d0152610b3701526116b46000f3fe608060405234801561001057600080fd5b50600436106101da5760003560e01c80636e04ff0d11610104578063ac500e79116100a2578063d826f88f11610071578063d826f88f14610520578063d832d92f14610534578063daf84a4f1461053c578063fc735e991461054f57600080fd5b8063ac500e79146104ab578063afb28d1f146104f0578063c297b86a146104f8578063c98f10b01461051857600080fd5b806386e330af116100de57806386e330af14610447578063917d895f1461045a578063947a36fb14610463578063a4441d971461046c57600080fd5b80636e04ff0d146103fa5780638340507c1461040d57806386b728e21461042057600080fd5b80633652d87c1161017c5780634bdb38621161014b5780634bdb38621461035b5780635b48391a146103a157806361bc221a146103e85780636250a13a146103f157600080fd5b80633652d87c1461030c5780634585e33b146103155780634a5479f3146103285780634b56a42e1461034857600080fd5b8063111554bd116101b8578063111554bd1461023c5780631d1970b7146102935780632ac64d75146102a05780632cb15864146102f557600080fd5b806302be021f146101df5780630fb172fb14610207578063102d538b14610228575b600080fd5b600a546101f29062010000900460ff1681565b60405190151581526020015b60405180910390f35b61021a610215366004610de5565b610561565b6040516101fe929190610e9a565b600a546101f2906301000000900460ff1681565b61029161024a366004610ebd565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b005b600a546101f29060ff1681565b6102916102ae366004610ebd565b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b6102fe60055481565b6040519081526020016101fe565b6102fe600b5481565b610291610323366004610efa565b610579565b61033b610336366004610f6c565b61088f565b6040516101fe9190610f85565b61021a610356366004610fbc565b61093b565b610291610369366004611086565b600a805491151562010000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffff909216919091179055565b6102916103af366004611086565b600a80549115156301000000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffff909216919091179055565b6102fe60065481565b6102fe60025481565b61021a610408366004610efa565b610a14565b61029161041b3660046110a8565b610ad2565b6101f27f000000000000000000000000000000000000000000000000000000000000000081565b6102916104553660046110f5565b610af0565b6102fe60045481565b6102fe60035481565b61029161047a366004611086565b600a80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055565b6001546104cb9073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101fe565b61033b610b07565b6000546104cb9073ffffffffffffffffffffffffffffffffffffffff1681565b61033b610b14565b610291600060048190556005819055600655565b6101f2610b21565b61029161054a366004610f6c565b600b55565b600a546101f290610100900460ff1681565b604080516000808252602082019092525b9250929050565b60007f00000000000000000000000000000000000000000000000000000000000000001561061857606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156105ed573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061061191906111a6565b905061061b565b50435b60055460000361062b5760058190555b60008061063a84860186610fbc565b600485905560065491935091506106529060016111ee565b600655604080516020808201835260008083528351918201909352918252600a54909190610100900460ff161561081857600a5460ff161561075557600154600b54855173ffffffffffffffffffffffffffffffffffffffff90921691638e760afe91879181106106c5576106c5611207565b60200260200101516040518263ffffffff1660e01b81526004016106e99190610f85565b6000604051808303816000875af1158015610708573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261074e9190810190611236565b9150610818565b600054600b54855173ffffffffffffffffffffffffffffffffffffffff90921691638e760afe918791811061078c5761078c611207565b60200260200101516040518263ffffffff1660e01b81526004016107b09190610f85565b6000604051808303816000875af11580156107cf573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526108159190810190611236565b91505b843373ffffffffffffffffffffffffffffffffffffffff167ff0f72c0b235fc8687d6a67c02ca543473a3cef8a18b48490f10e475a8dda139086600b548151811061086557610865611207565b6020026020010151858760405161087e939291906112ad565b60405180910390a350505050505050565b6007818154811061089f57600080fd5b9060005260206000200160009150905080546108ba906112f0565b80601f01602080910402602001604051908101604052809291908181526020018280546108e6906112f0565b80156109335780601f1061090857610100808354040283529160200191610933565b820191906000526020600020905b81548152906001019060200180831161091657829003601f168201915b505050505081565b600a5460009060609062010000900460ff16156109b9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f73686f756c6452657665727443616c6c6261636b20697320747275650000000060448201526064015b60405180910390fd5b600084846040516020016109ce929190611343565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018152919052600a546301000000900460ff16969095509350505050565b60006060610a20610b21565b610a6c576000848481818080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250959750919550610572945050505050565b604080516c6400000000000000000000000060208201528151601481830301815260348201928390527ff055e4a20000000000000000000000000000000000000000000000000000000090925242916109b0916008916007916009918691603801611469565b6008610ade838261157a565b506009610aeb828261157a565b505050565b8051610b03906007906020840190610c06565b5050565b600880546108ba906112f0565b600980546108ba906112f0565b6000600554600003610b335750600190565b60007f000000000000000000000000000000000000000000000000000000000000000015610bd257606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610ba7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bcb91906111a6565b9050610bd5565b50435b600254600554610be59083611694565b108015610c005750600354600454610bfd9083611694565b10155b91505090565b828054828255906000526020600020908101928215610c4c579160200282015b82811115610c4c5782518290610c3c908261157a565b5091602001919060010190610c26565b50610c58929150610c5c565b5090565b80821115610c58576000610c708282610c79565b50600101610c5c565b508054610c85906112f0565b6000825580601f10610c95575050565b601f016020900490600052602060002090810190610cb39190610cb6565b50565b5b80821115610c585760008155600101610cb7565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610d4157610d41610ccb565b604052919050565b600067ffffffffffffffff821115610d6357610d63610ccb565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f830112610da057600080fd5b8135610db3610dae82610d49565b610cfa565b818152846020838601011115610dc857600080fd5b816020850160208301376000918101602001919091529392505050565b60008060408385031215610df857600080fd5b82359150602083013567ffffffffffffffff811115610e1657600080fd5b610e2285828601610d8f565b9150509250929050565b60005b83811015610e47578181015183820152602001610e2f565b50506000910152565b60008151808452610e68816020860160208601610e2c565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8215158152604060208201526000610eb56040830184610e50565b949350505050565b600060208284031215610ecf57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff81168114610ef357600080fd5b9392505050565b60008060208385031215610f0d57600080fd5b823567ffffffffffffffff80821115610f2557600080fd5b818501915085601f830112610f3957600080fd5b813581811115610f4857600080fd5b866020828501011115610f5a57600080fd5b60209290920196919550909350505050565b600060208284031215610f7e57600080fd5b5035919050565b602081526000610ef36020830184610e50565b600067ffffffffffffffff821115610fb257610fb2610ccb565b5060051b60200190565b60008060408385031215610fcf57600080fd5b823567ffffffffffffffff80821115610fe757600080fd5b818501915085601f830112610ffb57600080fd5b8135602061100b610dae83610f98565b82815260059290921b8401810191818101908984111561102a57600080fd5b8286015b84811015611062578035868111156110465760008081fd5b6110548c86838b0101610d8f565b84525091830191830161102e565b509650508601359250508082111561107957600080fd5b50610e2285828601610d8f565b60006020828403121561109857600080fd5b81358015158114610ef357600080fd5b600080604083850312156110bb57600080fd5b823567ffffffffffffffff808211156110d357600080fd5b6110df86838701610d8f565b9350602085013591508082111561107957600080fd5b6000602080838503121561110857600080fd5b823567ffffffffffffffff8082111561112057600080fd5b818501915085601f83011261113457600080fd5b8135611142610dae82610f98565b81815260059190911b8301840190848101908883111561116157600080fd5b8585015b838110156111995780358581111561117d5760008081fd5b61118b8b89838a0101610d8f565b845250918601918601611165565b5098975050505050505050565b6000602082840312156111b857600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b80820180821115611201576112016111bf565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60006020828403121561124857600080fd5b815167ffffffffffffffff81111561125f57600080fd5b8201601f8101841361127057600080fd5b805161127e610dae82610d49565b81815285602083850101111561129357600080fd5b6112a4826020830160208601610e2c565b95945050505050565b6060815260006112c06060830186610e50565b82810360208401526112d28186610e50565b905082810360408401526112e68185610e50565b9695505050505050565b600181811c9082168061130457607f821691505b60208210810361133d577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b838110156113b8577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa08887030185526113a6868351610e50565b9550938201939082019060010161136c565b5050858403818701525050506112a48185610e50565b600081546113db816112f0565b8085526020600183811680156113f857600181146114305761145e565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b890101955061145e565b866000528260002060005b858110156114565781548a820186015290830190840161143b565b890184019650505b505050505092915050565b60a08152600061147c60a08301886113ce565b6020838203818501528188548084528284019150828160051b8501018a6000528360002060005b838110156114ee577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08784030185526114dc83836113ce565b948601949250600191820191016114a3565b50508681036040880152611502818b6113ce565b94505050505084606084015282810360808401526115208185610e50565b98975050505050505050565b601f821115610aeb57600081815260208120601f850160051c810160208610156115535750805b601f850160051c820191505b818110156115725782815560010161155f565b505050505050565b815167ffffffffffffffff81111561159457611594610ccb565b6115a8816115a284546112f0565b8461152c565b602080601f8311600181146115fb57600084156115c55750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555611572565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561164857888601518255948401946001909101908401611629565b508582101561168457878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b81810381811115611201576112016111bf56fea164736f6c6343000810000a307830303033376461303664353664303833666535393933393761343736396130343264363361613733646334656635373730396433316539393731613562343339307830303033353938343361353433656532666534313464633134633765373932306566313066343337323939306237396436333631636463306464316261373832",
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

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCaller) ProductionTestnetVerifierProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StreamsLookupUpkeep.contract.Call(opts, &out, "production_testnet_verifier_proxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) ProductionTestnetVerifierProxy() (common.Address, error) {
	return _StreamsLookupUpkeep.Contract.ProductionTestnetVerifierProxy(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCallerSession) ProductionTestnetVerifierProxy() (common.Address, error) {
	return _StreamsLookupUpkeep.Contract.ProductionTestnetVerifierProxy(&_StreamsLookupUpkeep.CallOpts)
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

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCaller) StagingTestnetVerifierProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StreamsLookupUpkeep.contract.Call(opts, &out, "staging_testnet_verifier_proxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) StagingTestnetVerifierProxy() (common.Address, error) {
	return _StreamsLookupUpkeep.Contract.StagingTestnetVerifierProxy(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCallerSession) StagingTestnetVerifierProxy() (common.Address, error) {
	return _StreamsLookupUpkeep.Contract.StagingTestnetVerifierProxy(&_StreamsLookupUpkeep.CallOpts)
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

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCaller) VerifyNthReport(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StreamsLookupUpkeep.contract.Call(opts, &out, "verifyNthReport")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) VerifyNthReport() (*big.Int, error) {
	return _StreamsLookupUpkeep.Contract.VerifyNthReport(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCallerSession) VerifyNthReport() (*big.Int, error) {
	return _StreamsLookupUpkeep.Contract.VerifyNthReport(&_StreamsLookupUpkeep.CallOpts)
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

func (_StreamsLookupUpkeep *StreamsLookupUpkeepTransactor) SetProductionTestnetVerifierProxy(opts *bind.TransactOpts, proxy common.Address) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.contract.Transact(opts, "setProductionTestnetVerifierProxy", proxy)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) SetProductionTestnetVerifierProxy(proxy common.Address) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.Contract.SetProductionTestnetVerifierProxy(&_StreamsLookupUpkeep.TransactOpts, proxy)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepTransactorSession) SetProductionTestnetVerifierProxy(proxy common.Address) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.Contract.SetProductionTestnetVerifierProxy(&_StreamsLookupUpkeep.TransactOpts, proxy)
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

func (_StreamsLookupUpkeep *StreamsLookupUpkeepTransactor) SetStaging(opts *bind.TransactOpts, _staging bool) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.contract.Transact(opts, "setStaging", _staging)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) SetStaging(_staging bool) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.Contract.SetStaging(&_StreamsLookupUpkeep.TransactOpts, _staging)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepTransactorSession) SetStaging(_staging bool) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.Contract.SetStaging(&_StreamsLookupUpkeep.TransactOpts, _staging)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepTransactor) SetStagingTestnetVerifierProxy(opts *bind.TransactOpts, proxy common.Address) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.contract.Transact(opts, "setStagingTestnetVerifierProxy", proxy)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) SetStagingTestnetVerifierProxy(proxy common.Address) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.Contract.SetStagingTestnetVerifierProxy(&_StreamsLookupUpkeep.TransactOpts, proxy)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepTransactorSession) SetStagingTestnetVerifierProxy(proxy common.Address) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.Contract.SetStagingTestnetVerifierProxy(&_StreamsLookupUpkeep.TransactOpts, proxy)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepTransactor) SetVerifyNthReport(opts *bind.TransactOpts, _n *big.Int) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.contract.Transact(opts, "setVerifyNthReport", _n)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) SetVerifyNthReport(_n *big.Int) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.Contract.SetVerifyNthReport(&_StreamsLookupUpkeep.TransactOpts, _n)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepTransactorSession) SetVerifyNthReport(_n *big.Int) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.Contract.SetVerifyNthReport(&_StreamsLookupUpkeep.TransactOpts, _n)
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

	ProductionTestnetVerifierProxy(opts *bind.CallOpts) (common.Address, error)

	ShouldRevertCallback(opts *bind.CallOpts) (bool, error)

	Staging(opts *bind.CallOpts) (bool, error)

	StagingTestnetVerifierProxy(opts *bind.CallOpts) (common.Address, error)

	TestRange(opts *bind.CallOpts) (*big.Int, error)

	TimeParamKey(opts *bind.CallOpts) (string, error)

	UseArbBlock(opts *bind.CallOpts) (bool, error)

	Verify(opts *bind.CallOpts) (bool, error)

	VerifyNthReport(opts *bind.CallOpts) (*big.Int, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	Reset(opts *bind.TransactOpts) (*types.Transaction, error)

	SetCallbackReturnBool(opts *bind.TransactOpts, value bool) (*types.Transaction, error)

	SetFeeds(opts *bind.TransactOpts, _feeds []string) (*types.Transaction, error)

	SetParamKeys(opts *bind.TransactOpts, _feedParamKey string, _timeParamKey string) (*types.Transaction, error)

	SetProductionTestnetVerifierProxy(opts *bind.TransactOpts, proxy common.Address) (*types.Transaction, error)

	SetShouldRevertCallback(opts *bind.TransactOpts, value bool) (*types.Transaction, error)

	SetStaging(opts *bind.TransactOpts, _staging bool) (*types.Transaction, error)

	SetStagingTestnetVerifierProxy(opts *bind.TransactOpts, proxy common.Address) (*types.Transaction, error)

	SetVerifyNthReport(opts *bind.TransactOpts, _n *big.Int) (*types.Transaction, error)

	FilterMercuryPerformEvent(opts *bind.FilterOpts, sender []common.Address, blockNumber []*big.Int) (*StreamsLookupUpkeepMercuryPerformEventIterator, error)

	WatchMercuryPerformEvent(opts *bind.WatchOpts, sink chan<- *StreamsLookupUpkeepMercuryPerformEvent, sender []common.Address, blockNumber []*big.Int) (event.Subscription, error)

	ParseMercuryPerformEvent(log types.Log) (*StreamsLookupUpkeepMercuryPerformEvent, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
