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
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_useArbBlock\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"_staging\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"_verify\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"StreamsLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v0\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"verifiedV0\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"ed\",\"type\":\"bytes\"}],\"name\":\"MercuryPerformEvent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"callbackReturnBool\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"errCode\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkErrorHandler\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feeds\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"interval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"setCallbackReturnBool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_feeds\",\"type\":\"string[]\"}],\"name\":\"setFeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_timeParamKey\",\"type\":\"string\"}],\"name\":\"setParamKeys\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"setShouldRetryOnErrorBool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"setShouldRevertCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"shouldRetryOnError\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"shouldRevertCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"staging\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbBlock\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162001c3f38038062001c3f83398101604081905262000034916200020f565b60008581556001859055600281905560038190556004558215156080526040805180820190915260078152666665656449447360c81b60208201526006906200007e908262000312565b50604080518082019091526009815268074696d657374616d760bc1b6020820152600790620000ae908262000312565b50604051806020016040528060405180608001604052806042815260200162001bfd604291399052620000e690600590600162000122565b506008805461ffff191692151561ff00191692909217610100911515919091021764ffff000000191664010100000017905550620003de915050565b8280548282559060005260206000209081019282156200016d579160200282015b828111156200016d57825182906200015c908262000312565b509160200191906001019062000143565b506200017b9291506200017f565b5090565b808211156200017b576000620001968282620001a0565b506001016200017f565b508054620001ae9062000283565b6000825580601f10620001bf575050565b601f016020900490600052602060002090810190620001df9190620001e2565b50565b5b808211156200017b5760008155600101620001e3565b805180151581146200020a57600080fd5b919050565b600080600080600060a086880312156200022857600080fd5b85519450602086015193506200024160408701620001f9565b92506200025160608701620001f9565b91506200026160808701620001f9565b90509295509295909350565b634e487b7160e01b600052604160045260246000fd5b600181811c908216806200029857607f821691505b602082108103620002b957634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200030d57600081815260208120601f850160051c81016020861015620002e85750805b601f850160051c820191505b818110156200030957828155600101620002f4565b5050505b505050565b81516001600160401b038111156200032e576200032e6200026d565b62000346816200033f845462000283565b84620002bf565b602080601f8311600181146200037e5760008415620003655750858301515b600019600386901b1c1916600185901b17855562000309565b600085815260208120601f198616915b82811015620003af578886015182559484019460019091019084016200038e565b5085821015620003ce5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6080516117ee6200040f600039600081816103980152818161052a01528181610aa40152610c1301526117ee6000f3fe608060405234801561001057600080fd5b50600436106101a35760003560e01c806361bc221a116100ee578063917d895f11610097578063c98f10b011610071578063c98f10b0146103e7578063d826f88f146103ef578063d832d92f14610403578063fc735e991461040b57600080fd5b8063917d895f146103cd578063947a36fb146103d6578063afb28d1f146103df57600080fd5b80638340507c116100c85780638340507c1461038057806386b728e21461039357806386e330af146103ba57600080fd5b806361bc221a1461035b5780636250a13a146103645780636e04ff0d1461036d57600080fd5b806342eb3d92116101505780634b56a42e1161012a5780634b56a42e146102bb5780634bdb3862146102ce5780635b48391a1461031457600080fd5b806342eb3d92146102735780634585e33b146102885780634a5479f31461029b57600080fd5b806313fab5901161018157806313fab590146102055780631d1970b71461024f5780632cb158641461025c57600080fd5b806302be021f146101a85780630fb172fb146101d0578063102d538b146101f1575b600080fd5b6008546101bb9062010000900460ff1681565b60405190151581526020015b60405180910390f35b6101e36101de366004610ec1565b61041d565b6040516101c7929190610f76565b6008546101bb906301000000900460ff1681565b61024d610213366004610f99565b60088054911515640100000000027fffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffff909216919091179055565b005b6008546101bb9060ff1681565b61026560035481565b6040519081526020016101c7565b6008546101bb90640100000000900460ff1681565b61024d610296366004610fc2565b610526565b6102ae6102a9366004611034565b610853565b6040516101c7919061104d565b6101e36102c9366004611084565b6108ff565b61024d6102dc366004610f99565b6008805491151562010000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffff909216919091179055565b61024d610322366004610f99565b600880549115156301000000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffff909216919091179055565b61026560045481565b61026560005481565b6101e361037b366004610fc2565b6109d8565b61024d61038e36600461114e565b610bae565b6101bb7f000000000000000000000000000000000000000000000000000000000000000081565b61024d6103c836600461119b565b610bcc565b61026560025481565b61026560015481565b6102ae610be3565b6102ae610bf0565b61024d600060028190556003819055600455565b6101bb610bfd565b6008546101bb90610100900460ff1681565b6040805160028082526060828101909352600092918391816020015b60608152602001906001900390816104395750506040805160208101889052919250016040516020818303038152906040528160008151811061047e5761047e61124c565b60200260200101819052508360405160200161049a919061104d565b604051602081830303815290604052816001815181106104bc576104bc61124c565b6020026020010181905250600081856040516020016104dc92919061127b565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018152919052600854640100000000900460ff169450925050505b9250929050565b60007f0000000000000000000000000000000000000000000000000000000000000000156105c557606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561059a573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105be919061130f565b90506105c8565b50435b6003546000036105d85760038190555b6000806105e784860186611084565b600285905560045491935091506105ff906001611357565b600455604080516020808201835260008083528351918201909352918252600854909190610100900460ff16156107dd5760085460ff161561070e577360448b880c9f3b501af3f343da9284148bd7d77c73ffffffffffffffffffffffffffffffffffffffff16638e760afe8560008151811061067e5761067e61124c565b60200260200101516040518263ffffffff1660e01b81526004016106a2919061104d565b6000604051808303816000875af11580156106c1573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526107079190810190611370565b91506107dd565b7309dff56a4ff44e0f4436260a04f5cfa65636a48173ffffffffffffffffffffffffffffffffffffffff16638e760afe856000815181106107515761075161124c565b60200260200101516040518263ffffffff1660e01b8152600401610775919061104d565b6000604051808303816000875af1158015610794573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526107da9190810190611370565b91505b843373ffffffffffffffffffffffffffffffffffffffff167ff0f72c0b235fc8687d6a67c02ca543473a3cef8a18b48490f10e475a8dda1390866000815181106108295761082961124c565b60200260200101518587604051610842939291906113de565b60405180910390a350505050505050565b6005818154811061086357600080fd5b90600052602060002001600091509050805461087e90611421565b80601f01602080910402602001604051908101604052809291908181526020018280546108aa90611421565b80156108f75780601f106108cc576101008083540402835291602001916108f7565b820191906000526020600020905b8154815290600101906020018083116108da57829003601f168201915b505050505081565b60085460009060609062010000900460ff161561097d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f73686f756c6452657665727443616c6c6261636b20697320747275650000000060448201526064015b60405180910390fd5b6000848460405160200161099292919061127b565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529190526008546301000000900460ff16969095509350505050565b600060606109e4610bfd565b610a30576000848481818080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525095975091955061051f945050505050565b6040517f66656564496448657800000000000000000000000000000000000000000000006020820152600090602901604051602081830303815290604052805190602001206006604051602001610a879190611474565b6040516020818303038152906040528051906020012003610b46577f000000000000000000000000000000000000000000000000000000000000000015610b3f57606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610b14573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b38919061130f565b9050610b49565b5043610b49565b50425b604080516c6400000000000000000000000060208201528151601481830301815260348201928390527ff055e4a200000000000000000000000000000000000000000000000000000000909252610974916006916005916007918691906038016115a3565b6006610bba83826116b4565b506007610bc782826116b4565b505050565b8051610bdf906005906020840190610ce2565b5050565b6006805461087e90611421565b6007805461087e90611421565b6000600354600003610c0f5750600190565b60007f000000000000000000000000000000000000000000000000000000000000000015610cae57606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610c83573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ca7919061130f565b9050610cb1565b50435b600054600354610cc190836117ce565b108015610cdc5750600154600254610cd990836117ce565b10155b91505090565b828054828255906000526020600020908101928215610d28579160200282015b82811115610d285782518290610d1890826116b4565b5091602001919060010190610d02565b50610d34929150610d38565b5090565b80821115610d34576000610d4c8282610d55565b50600101610d38565b508054610d6190611421565b6000825580601f10610d71575050565b601f016020900490600052602060002090810190610d8f9190610d92565b50565b5b80821115610d345760008155600101610d93565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610e1d57610e1d610da7565b604052919050565b600067ffffffffffffffff821115610e3f57610e3f610da7565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f830112610e7c57600080fd5b8135610e8f610e8a82610e25565b610dd6565b818152846020838601011115610ea457600080fd5b816020850160208301376000918101602001919091529392505050565b60008060408385031215610ed457600080fd5b82359150602083013567ffffffffffffffff811115610ef257600080fd5b610efe85828601610e6b565b9150509250929050565b60005b83811015610f23578181015183820152602001610f0b565b50506000910152565b60008151808452610f44816020860160208601610f08565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8215158152604060208201526000610f916040830184610f2c565b949350505050565b600060208284031215610fab57600080fd5b81358015158114610fbb57600080fd5b9392505050565b60008060208385031215610fd557600080fd5b823567ffffffffffffffff80821115610fed57600080fd5b818501915085601f83011261100157600080fd5b81358181111561101057600080fd5b86602082850101111561102257600080fd5b60209290920196919550909350505050565b60006020828403121561104657600080fd5b5035919050565b602081526000610fbb6020830184610f2c565b600067ffffffffffffffff82111561107a5761107a610da7565b5060051b60200190565b6000806040838503121561109757600080fd5b823567ffffffffffffffff808211156110af57600080fd5b818501915085601f8301126110c357600080fd5b813560206110d3610e8a83611060565b82815260059290921b840181019181810190898411156110f257600080fd5b8286015b8481101561112a5780358681111561110e5760008081fd5b61111c8c86838b0101610e6b565b8452509183019183016110f6565b509650508601359250508082111561114157600080fd5b50610efe85828601610e6b565b6000806040838503121561116157600080fd5b823567ffffffffffffffff8082111561117957600080fd5b61118586838701610e6b565b9350602085013591508082111561114157600080fd5b600060208083850312156111ae57600080fd5b823567ffffffffffffffff808211156111c657600080fd5b818501915085601f8301126111da57600080fd5b81356111e8610e8a82611060565b81815260059190911b8301840190848101908883111561120757600080fd5b8585015b8381101561123f578035858111156112235760008081fd5b6112318b89838a0101610e6b565b84525091860191860161120b565b5098975050505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b838110156112f0577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa08887030185526112de868351610f2c565b955093820193908201906001016112a4565b5050858403818701525050506113068185610f2c565b95945050505050565b60006020828403121561132157600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082018082111561136a5761136a611328565b92915050565b60006020828403121561138257600080fd5b815167ffffffffffffffff81111561139957600080fd5b8201601f810184136113aa57600080fd5b80516113b8610e8a82610e25565b8181528560208385010111156113cd57600080fd5b611306826020830160208601610f08565b6060815260006113f16060830186610f2c565b82810360208401526114038186610f2c565b905082810360408401526114178185610f2c565b9695505050505050565b600181811c9082168061143557607f821691505b60208210810361146e577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b600080835461148281611421565b6001828116801561149a57600181146114cd576114fc565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00841687528215158302870194506114fc565b8760005260208060002060005b858110156114f35781548a8201529084019082016114da565b50505082870194505b50929695505050505050565b6000815461151581611421565b808552602060018381168015611532576001811461156a57611598565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b8901019550611598565b866000528260002060005b858110156115905781548a8201860152908301908401611575565b890184019650505b505050505092915050565b60a0815260006115b660a0830188611508565b6020838203818501528188548084528284019150828160051b8501018a6000528360002060005b83811015611628577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08784030185526116168383611508565b948601949250600191820191016115dd565b5050868103604088015261163c818b611508565b945050505050846060840152828103608084015261165a8185610f2c565b98975050505050505050565b601f821115610bc757600081815260208120601f850160051c8101602086101561168d5750805b601f850160051c820191505b818110156116ac57828155600101611699565b505050505050565b815167ffffffffffffffff8111156116ce576116ce610da7565b6116e2816116dc8454611421565b84611666565b602080601f83116001811461173557600084156116ff5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556116ac565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561178257888601518255948401946001909101908401611763565b50858210156117be57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b8181038181111561136a5761136a61132856fea164736f6c6343000810000a307830303032386339313564366166306664363662626132643066633934303532323662636138643638303633333331323161376439383332313033643135363363",
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

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCaller) ShouldRetryOnError(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _StreamsLookupUpkeep.contract.Call(opts, &out, "shouldRetryOnError")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) ShouldRetryOnError() (bool, error) {
	return _StreamsLookupUpkeep.Contract.ShouldRetryOnError(&_StreamsLookupUpkeep.CallOpts)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepCallerSession) ShouldRetryOnError() (bool, error) {
	return _StreamsLookupUpkeep.Contract.ShouldRetryOnError(&_StreamsLookupUpkeep.CallOpts)
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

func (_StreamsLookupUpkeep *StreamsLookupUpkeepTransactor) SetShouldRetryOnErrorBool(opts *bind.TransactOpts, value bool) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.contract.Transact(opts, "setShouldRetryOnErrorBool", value)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepSession) SetShouldRetryOnErrorBool(value bool) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.Contract.SetShouldRetryOnErrorBool(&_StreamsLookupUpkeep.TransactOpts, value)
}

func (_StreamsLookupUpkeep *StreamsLookupUpkeepTransactorSession) SetShouldRetryOnErrorBool(value bool) (*types.Transaction, error) {
	return _StreamsLookupUpkeep.Contract.SetShouldRetryOnErrorBool(&_StreamsLookupUpkeep.TransactOpts, value)
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

	ShouldRetryOnError(opts *bind.CallOpts) (bool, error)

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

	SetShouldRetryOnErrorBool(opts *bind.TransactOpts, value bool) (*types.Transaction, error)

	SetShouldRevertCallback(opts *bind.TransactOpts, value bool) (*types.Transaction, error)

	FilterMercuryPerformEvent(opts *bind.FilterOpts, sender []common.Address, blockNumber []*big.Int) (*StreamsLookupUpkeepMercuryPerformEventIterator, error)

	WatchMercuryPerformEvent(opts *bind.WatchOpts, sink chan<- *StreamsLookupUpkeepMercuryPerformEvent, sender []common.Address, blockNumber []*big.Int) (event.Subscription, error)

	ParseMercuryPerformEvent(log types.Log) (*StreamsLookupUpkeepMercuryPerformEvent, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
