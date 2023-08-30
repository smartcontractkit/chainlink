// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mercury_upkeep_wrapper

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

var MercuryUpkeepMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_useArbBlock\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"_staging\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"_verify\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"FeedLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v0\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v1\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"verifiedV0\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"verifiedV1\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"ed\",\"type\":\"bytes\"}],\"name\":\"MercuryPerformEvent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"callbackReturnBool\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feeds\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"interval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"setCallbackReturnBool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_feeds\",\"type\":\"string[]\"}],\"name\":\"setFeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_timeParamKey\",\"type\":\"string\"}],\"name\":\"setParamKeys\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"setShouldRevertCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"shouldRevertCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"staging\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbBlock\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162001b9c38038062001b9c833981016040819052620000349162000232565b60008581556001859055600281905560038190556004558215156080526040805180820190915260098152680cccacac892c890caf60bb1b602082015260069062000080908262000335565b5060408051808201909152600b81526a313637b1b5a73ab6b132b960a91b6020820152600790620000b2908262000335565b50604051806040016040528060405180608001604052806042815260200162001b1860429139815260200160405180608001604052806042815260200162001b5a6042913990526200010990600590600262000145565b506008805463ff000000199215156101000261ff00199415159490941661ffff1990911617929092171663010000001790555062000401915050565b82805482825590600052602060002090810192821562000190579160200282015b828111156200019057825182906200017f908262000335565b509160200191906001019062000166565b506200019e929150620001a2565b5090565b808211156200019e576000620001b98282620001c3565b50600101620001a2565b508054620001d190620002a6565b6000825580601f10620001e2575050565b601f01602090049060005260206000209081019062000202919062000205565b50565b5b808211156200019e576000815560010162000206565b805180151581146200022d57600080fd5b919050565b600080600080600060a086880312156200024b57600080fd5b855194506020860151935062000264604087016200021c565b925062000274606087016200021c565b915062000284608087016200021c565b90509295509295909350565b634e487b7160e01b600052604160045260246000fd5b600181811c90821680620002bb57607f821691505b602082108103620002dc57634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200033057600081815260208120601f850160051c810160208610156200030b5750805b601f850160051c820191505b818110156200032c5782815560010162000317565b5050505b505050565b81516001600160401b0381111562000351576200035162000290565b6200036981620003628454620002a6565b84620002e2565b602080601f831160018114620003a15760008415620003885750858301515b600019600386901b1c1916600185901b1785556200032c565b600085815260208120601f198616915b82811015620003d257888601518255948401946001909101908401620003b1565b5085821015620003f15787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6080516116e662000432600039600081816103070152818161039001528181610a560152610bbe01526116e66000f3fe608060405234801561001057600080fd5b50600436106101825760003560e01c80636e04ff0d116100d8578063947a36fb1161008c578063d826f88f11610066578063d826f88f1461035e578063d832d92f14610372578063fc735e991461037a57600080fd5b8063947a36fb14610345578063afb28d1f1461034e578063c98f10b01461035657600080fd5b806386b728e2116100bd57806386b728e21461030257806386e330af14610329578063917d895f1461033c57600080fd5b80636e04ff0d146102dc5780638340507c146102ef57600080fd5b80634a5479f31161013a5780635b48391a116101145780635b48391a1461028357806361bc221a146102ca5780636250a13a146102d357600080fd5b80634a5479f3146101fc5780634b56a42e1461021c5780634bdb38621461023d57600080fd5b80631d1970b71161016b5780631d1970b7146101c35780632cb15864146101d05780634585e33b146101e757600080fd5b806302be021f14610187578063102d538b146101af575b600080fd5b60085461019a9062010000900460ff1681565b60405190151581526020015b60405180910390f35b60085461019a906301000000900460ff1681565b60085461019a9060ff1681565b6101d960035481565b6040519081526020016101a6565b6101fa6101f5366004610d52565b61038c565b005b61020f61020a366004610dc4565b610873565b6040516101a69190610e4b565b61022f61022a366004610fa3565b61091f565b6040516101a6929190611077565b6101fa61024b36600461109a565b6008805491151562010000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffff909216919091179055565b6101fa61029136600461109a565b600880549115156301000000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffff909216919091179055565b6101d960045481565b6101d960005481565b61022f6102ea366004610d52565b6109fa565b6101fa6102fd3660046110bc565b610b59565b61019a7f000000000000000000000000000000000000000000000000000000000000000081565b6101fa610337366004611109565b610b77565b6101d960025481565b6101d960015481565b61020f610b8e565b61020f610b9b565b6101fa600060028190556003819055600455565b61019a610ba8565b60085461019a90610100900460ff1681565b60007f00000000000000000000000000000000000000000000000000000000000000001561042b57606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610400573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061042491906111ba565b905061042e565b50435b60035460000361043e5760038190555b60008061044d84860186610fa3565b60028590556004549193509150610465906001611202565b600455604080516020808201835260008083528351918201909352918252600854909190610100900460ff16156107df5760085460ff1615610642577360448b880c9f3b501af3f343da9284148bd7d77c73ffffffffffffffffffffffffffffffffffffffff16638e760afe856000815181106104e4576104e461121b565b60200260200101516040518263ffffffff1660e01b81526004016105089190610e4b565b6000604051808303816000875af1158015610527573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261056d919081019061124a565b91507360448b880c9f3b501af3f343da9284148bd7d77c73ffffffffffffffffffffffffffffffffffffffff16638e760afe856001815181106105b2576105b261121b565b60200260200101516040518263ffffffff1660e01b81526004016105d69190610e4b565b6000604051808303816000875af11580156105f5573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261063b919081019061124a565b90506107df565b7309dff56a4ff44e0f4436260a04f5cfa65636a48173ffffffffffffffffffffffffffffffffffffffff16638e760afe856000815181106106855761068561121b565b60200260200101516040518263ffffffff1660e01b81526004016106a99190610e4b565b6000604051808303816000875af11580156106c8573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261070e919081019061124a565b91507309dff56a4ff44e0f4436260a04f5cfa65636a48173ffffffffffffffffffffffffffffffffffffffff16638e760afe856001815181106107535761075361121b565b60200260200101516040518263ffffffff1660e01b81526004016107779190610e4b565b6000604051808303816000875af1158015610796573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526107dc919081019061124a565b90505b843373ffffffffffffffffffffffffffffffffffffffff167f1c85d6186f024e964616014c8247533455ec5129a5095711202292f8a7ea1d548660008151811061082b5761082b61121b565b6020026020010151876001815181106108465761084661121b565b60200260200101518686896040516108629594939291906112c1565b60405180910390a350505050505050565b6005818154811061088357600080fd5b90600052602060002001600091509050805461089e9061132e565b80601f01602080910402602001604051908101604052809291908181526020018280546108ca9061132e565b80156109175780601f106108ec57610100808354040283529160200191610917565b820191906000526020600020905b8154815290600101906020018083116108fa57829003601f168201915b505050505081565b60085460009060609062010000900460ff161561099d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f73686f756c6452657665727443616c6c6261636b20697320747275650000000060448201526064015b60405180910390fd5b600084846040516020016109b2929190611381565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529190526008546301000000900460ff1693509150505b9250929050565b60006060610a06610ba8565b610a52576000848481818080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509597509195506109f3945050505050565b60007f000000000000000000000000000000000000000000000000000000000000000015610af157606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610ac6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610aea91906111ba565b9050610af4565b50435b604080516c6400000000000000000000000060208201528151601481830301815260348201928390527f7ddd933e00000000000000000000000000000000000000000000000000000000909252610994916006916005916007918691906038016114a7565b6006610b6583826115ac565b506007610b7282826115ac565b505050565b8051610b8a906005906020840190610c8d565b5050565b6006805461089e9061132e565b6007805461089e9061132e565b6000600354600003610bba5750600190565b60007f000000000000000000000000000000000000000000000000000000000000000015610c5957606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610c2e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c5291906111ba565b9050610c5c565b50435b600054600354610c6c90836116c6565b108015610c875750600154600254610c8490836116c6565b10155b91505090565b828054828255906000526020600020908101928215610cd3579160200282015b82811115610cd35782518290610cc390826115ac565b5091602001919060010190610cad565b50610cdf929150610ce3565b5090565b80821115610cdf576000610cf78282610d00565b50600101610ce3565b508054610d0c9061132e565b6000825580601f10610d1c575050565b601f016020900490600052602060002090810190610d3a9190610d3d565b50565b5b80821115610cdf5760008155600101610d3e565b60008060208385031215610d6557600080fd5b823567ffffffffffffffff80821115610d7d57600080fd5b818501915085601f830112610d9157600080fd5b813581811115610da057600080fd5b866020828501011115610db257600080fd5b60209290920196919550909350505050565b600060208284031215610dd657600080fd5b5035919050565b60005b83811015610df8578181015183820152602001610de0565b50506000910152565b60008151808452610e19816020860160208601610ddd565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000610e5e6020830184610e01565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610edb57610edb610e65565b604052919050565b600067ffffffffffffffff821115610efd57610efd610e65565b5060051b60200190565b600067ffffffffffffffff821115610f2157610f21610e65565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f830112610f5e57600080fd5b8135610f71610f6c82610f07565b610e94565b818152846020838601011115610f8657600080fd5b816020850160208301376000918101602001919091529392505050565b60008060408385031215610fb657600080fd5b823567ffffffffffffffff80821115610fce57600080fd5b818501915085601f830112610fe257600080fd5b81356020610ff2610f6c83610ee3565b82815260059290921b8401810191818101908984111561101157600080fd5b8286015b848110156110495780358681111561102d5760008081fd5b61103b8c86838b0101610f4d565b845250918301918301611015565b509650508601359250508082111561106057600080fd5b5061106d85828601610f4d565b9150509250929050565b82151581526040602082015260006110926040830184610e01565b949350505050565b6000602082840312156110ac57600080fd5b81358015158114610e5e57600080fd5b600080604083850312156110cf57600080fd5b823567ffffffffffffffff808211156110e757600080fd5b6110f386838701610f4d565b9350602085013591508082111561106057600080fd5b6000602080838503121561111c57600080fd5b823567ffffffffffffffff8082111561113457600080fd5b818501915085601f83011261114857600080fd5b8135611156610f6c82610ee3565b81815260059190911b8301840190848101908883111561117557600080fd5b8585015b838110156111ad578035858111156111915760008081fd5b61119f8b89838a0101610f4d565b845250918601918601611179565b5098975050505050505050565b6000602082840312156111cc57600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b80820180821115611215576112156111d3565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60006020828403121561125c57600080fd5b815167ffffffffffffffff81111561127357600080fd5b8201601f8101841361128457600080fd5b8051611292610f6c82610f07565b8181528560208385010111156112a757600080fd5b6112b8826020830160208601610ddd565b95945050505050565b60a0815260006112d460a0830188610e01565b82810360208401526112e68188610e01565b905082810360408401526112fa8187610e01565b9050828103606084015261130e8186610e01565b905082810360808401526113228185610e01565b98975050505050505050565b600181811c9082168061134257607f821691505b60208210810361137b577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b838110156113f6577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa08887030185526113e4868351610e01565b955093820193908201906001016113aa565b5050858403818701525050506112b88185610e01565b600081546114198161132e565b808552602060018381168015611436576001811461146e5761149c565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b890101955061149c565b866000528260002060005b858110156114945781548a8201860152908301908401611479565b890184019650505b505050505092915050565b60a0815260006114ba60a083018861140c565b6020838203818501528188548084528284019150828160051b8501018a6000528360002060005b8381101561152c577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe087840301855261151a838361140c565b948601949250600191820191016114e1565b50508681036040880152611540818b61140c565b94505050505084606084015282810360808401526113228185610e01565b601f821115610b7257600081815260208120601f850160051c810160208610156115855750805b601f850160051c820191505b818110156115a457828155600101611591565b505050505050565b815167ffffffffffffffff8111156115c6576115c6610e65565b6115da816115d4845461132e565b8461155e565b602080601f83116001811461162d57600084156115f75750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556115a4565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561167a5788860151825594840194600190910190840161165b565b50858210156116b657878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b81810381811115611215576112156111d356fea164736f6c6343000810000a307834353534343832643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030307834323534343332643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030",
}

var MercuryUpkeepABI = MercuryUpkeepMetaData.ABI

var MercuryUpkeepBin = MercuryUpkeepMetaData.Bin

func DeployMercuryUpkeep(auth *bind.TransactOpts, backend bind.ContractBackend, _testRange *big.Int, _interval *big.Int, _useArbBlock bool, _staging bool, _verify bool) (common.Address, *types.Transaction, *MercuryUpkeep, error) {
	parsed, err := MercuryUpkeepMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MercuryUpkeepBin), backend, _testRange, _interval, _useArbBlock, _staging, _verify)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MercuryUpkeep{MercuryUpkeepCaller: MercuryUpkeepCaller{contract: contract}, MercuryUpkeepTransactor: MercuryUpkeepTransactor{contract: contract}, MercuryUpkeepFilterer: MercuryUpkeepFilterer{contract: contract}}, nil
}

type MercuryUpkeep struct {
	address common.Address
	abi     abi.ABI
	MercuryUpkeepCaller
	MercuryUpkeepTransactor
	MercuryUpkeepFilterer
}

type MercuryUpkeepCaller struct {
	contract *bind.BoundContract
}

type MercuryUpkeepTransactor struct {
	contract *bind.BoundContract
}

type MercuryUpkeepFilterer struct {
	contract *bind.BoundContract
}

type MercuryUpkeepSession struct {
	Contract     *MercuryUpkeep
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MercuryUpkeepCallerSession struct {
	Contract *MercuryUpkeepCaller
	CallOpts bind.CallOpts
}

type MercuryUpkeepTransactorSession struct {
	Contract     *MercuryUpkeepTransactor
	TransactOpts bind.TransactOpts
}

type MercuryUpkeepRaw struct {
	Contract *MercuryUpkeep
}

type MercuryUpkeepCallerRaw struct {
	Contract *MercuryUpkeepCaller
}

type MercuryUpkeepTransactorRaw struct {
	Contract *MercuryUpkeepTransactor
}

func NewMercuryUpkeep(address common.Address, backend bind.ContractBackend) (*MercuryUpkeep, error) {
	abi, err := abi.JSON(strings.NewReader(MercuryUpkeepABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMercuryUpkeep(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MercuryUpkeep{address: address, abi: abi, MercuryUpkeepCaller: MercuryUpkeepCaller{contract: contract}, MercuryUpkeepTransactor: MercuryUpkeepTransactor{contract: contract}, MercuryUpkeepFilterer: MercuryUpkeepFilterer{contract: contract}}, nil
}

func NewMercuryUpkeepCaller(address common.Address, caller bind.ContractCaller) (*MercuryUpkeepCaller, error) {
	contract, err := bindMercuryUpkeep(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MercuryUpkeepCaller{contract: contract}, nil
}

func NewMercuryUpkeepTransactor(address common.Address, transactor bind.ContractTransactor) (*MercuryUpkeepTransactor, error) {
	contract, err := bindMercuryUpkeep(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MercuryUpkeepTransactor{contract: contract}, nil
}

func NewMercuryUpkeepFilterer(address common.Address, filterer bind.ContractFilterer) (*MercuryUpkeepFilterer, error) {
	contract, err := bindMercuryUpkeep(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MercuryUpkeepFilterer{contract: contract}, nil
}

func bindMercuryUpkeep(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MercuryUpkeepMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MercuryUpkeep *MercuryUpkeepRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MercuryUpkeep.Contract.MercuryUpkeepCaller.contract.Call(opts, result, method, params...)
}

func (_MercuryUpkeep *MercuryUpkeepRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.MercuryUpkeepTransactor.contract.Transfer(opts)
}

func (_MercuryUpkeep *MercuryUpkeepRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.MercuryUpkeepTransactor.contract.Transact(opts, method, params...)
}

func (_MercuryUpkeep *MercuryUpkeepCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MercuryUpkeep.Contract.contract.Call(opts, result, method, params...)
}

func (_MercuryUpkeep *MercuryUpkeepTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.contract.Transfer(opts)
}

func (_MercuryUpkeep *MercuryUpkeepTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.contract.Transact(opts, method, params...)
}

func (_MercuryUpkeep *MercuryUpkeepCaller) CallbackReturnBool(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "callbackReturnBool")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) CallbackReturnBool() (bool, error) {
	return _MercuryUpkeep.Contract.CallbackReturnBool(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) CallbackReturnBool() (bool, error) {
	return _MercuryUpkeep.Contract.CallbackReturnBool(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCaller) CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (bool, []byte, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "checkCallback", values, extraData)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) CheckCallback(values [][]byte, extraData []byte) (bool, []byte, error) {
	return _MercuryUpkeep.Contract.CheckCallback(&_MercuryUpkeep.CallOpts, values, extraData)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) CheckCallback(values [][]byte, extraData []byte) (bool, []byte, error) {
	return _MercuryUpkeep.Contract.CheckCallback(&_MercuryUpkeep.CallOpts, values, extraData)
}

func (_MercuryUpkeep *MercuryUpkeepCaller) CheckUpkeep(opts *bind.CallOpts, data []byte) (bool, []byte, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "checkUpkeep", data)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) CheckUpkeep(data []byte) (bool, []byte, error) {
	return _MercuryUpkeep.Contract.CheckUpkeep(&_MercuryUpkeep.CallOpts, data)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) CheckUpkeep(data []byte) (bool, []byte, error) {
	return _MercuryUpkeep.Contract.CheckUpkeep(&_MercuryUpkeep.CallOpts, data)
}

func (_MercuryUpkeep *MercuryUpkeepCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) Counter() (*big.Int, error) {
	return _MercuryUpkeep.Contract.Counter(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) Counter() (*big.Int, error) {
	return _MercuryUpkeep.Contract.Counter(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCaller) Eligible(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "eligible")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) Eligible() (bool, error) {
	return _MercuryUpkeep.Contract.Eligible(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) Eligible() (bool, error) {
	return _MercuryUpkeep.Contract.Eligible(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCaller) FeedParamKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "feedParamKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) FeedParamKey() (string, error) {
	return _MercuryUpkeep.Contract.FeedParamKey(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) FeedParamKey() (string, error) {
	return _MercuryUpkeep.Contract.FeedParamKey(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCaller) Feeds(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "feeds", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) Feeds(arg0 *big.Int) (string, error) {
	return _MercuryUpkeep.Contract.Feeds(&_MercuryUpkeep.CallOpts, arg0)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) Feeds(arg0 *big.Int) (string, error) {
	return _MercuryUpkeep.Contract.Feeds(&_MercuryUpkeep.CallOpts, arg0)
}

func (_MercuryUpkeep *MercuryUpkeepCaller) InitialBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "initialBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) InitialBlock() (*big.Int, error) {
	return _MercuryUpkeep.Contract.InitialBlock(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) InitialBlock() (*big.Int, error) {
	return _MercuryUpkeep.Contract.InitialBlock(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCaller) Interval(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "interval")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) Interval() (*big.Int, error) {
	return _MercuryUpkeep.Contract.Interval(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) Interval() (*big.Int, error) {
	return _MercuryUpkeep.Contract.Interval(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCaller) PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "previousPerformBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) PreviousPerformBlock() (*big.Int, error) {
	return _MercuryUpkeep.Contract.PreviousPerformBlock(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) PreviousPerformBlock() (*big.Int, error) {
	return _MercuryUpkeep.Contract.PreviousPerformBlock(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCaller) ShouldRevertCallback(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "shouldRevertCallback")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) ShouldRevertCallback() (bool, error) {
	return _MercuryUpkeep.Contract.ShouldRevertCallback(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) ShouldRevertCallback() (bool, error) {
	return _MercuryUpkeep.Contract.ShouldRevertCallback(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCaller) Staging(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "staging")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) Staging() (bool, error) {
	return _MercuryUpkeep.Contract.Staging(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) Staging() (bool, error) {
	return _MercuryUpkeep.Contract.Staging(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCaller) TestRange(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "testRange")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) TestRange() (*big.Int, error) {
	return _MercuryUpkeep.Contract.TestRange(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) TestRange() (*big.Int, error) {
	return _MercuryUpkeep.Contract.TestRange(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCaller) TimeParamKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "timeParamKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) TimeParamKey() (string, error) {
	return _MercuryUpkeep.Contract.TimeParamKey(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) TimeParamKey() (string, error) {
	return _MercuryUpkeep.Contract.TimeParamKey(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCaller) UseArbBlock(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "useArbBlock")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) UseArbBlock() (bool, error) {
	return _MercuryUpkeep.Contract.UseArbBlock(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) UseArbBlock() (bool, error) {
	return _MercuryUpkeep.Contract.UseArbBlock(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCaller) Verify(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _MercuryUpkeep.contract.Call(opts, &out, "verify")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MercuryUpkeep *MercuryUpkeepSession) Verify() (bool, error) {
	return _MercuryUpkeep.Contract.Verify(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepCallerSession) Verify() (bool, error) {
	return _MercuryUpkeep.Contract.Verify(&_MercuryUpkeep.CallOpts)
}

func (_MercuryUpkeep *MercuryUpkeepTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _MercuryUpkeep.contract.Transact(opts, "performUpkeep", performData)
}

func (_MercuryUpkeep *MercuryUpkeepSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.PerformUpkeep(&_MercuryUpkeep.TransactOpts, performData)
}

func (_MercuryUpkeep *MercuryUpkeepTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.PerformUpkeep(&_MercuryUpkeep.TransactOpts, performData)
}

func (_MercuryUpkeep *MercuryUpkeepTransactor) Reset(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MercuryUpkeep.contract.Transact(opts, "reset")
}

func (_MercuryUpkeep *MercuryUpkeepSession) Reset() (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.Reset(&_MercuryUpkeep.TransactOpts)
}

func (_MercuryUpkeep *MercuryUpkeepTransactorSession) Reset() (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.Reset(&_MercuryUpkeep.TransactOpts)
}

func (_MercuryUpkeep *MercuryUpkeepTransactor) SetCallbackReturnBool(opts *bind.TransactOpts, value bool) (*types.Transaction, error) {
	return _MercuryUpkeep.contract.Transact(opts, "setCallbackReturnBool", value)
}

func (_MercuryUpkeep *MercuryUpkeepSession) SetCallbackReturnBool(value bool) (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.SetCallbackReturnBool(&_MercuryUpkeep.TransactOpts, value)
}

func (_MercuryUpkeep *MercuryUpkeepTransactorSession) SetCallbackReturnBool(value bool) (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.SetCallbackReturnBool(&_MercuryUpkeep.TransactOpts, value)
}

func (_MercuryUpkeep *MercuryUpkeepTransactor) SetFeeds(opts *bind.TransactOpts, _feeds []string) (*types.Transaction, error) {
	return _MercuryUpkeep.contract.Transact(opts, "setFeeds", _feeds)
}

func (_MercuryUpkeep *MercuryUpkeepSession) SetFeeds(_feeds []string) (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.SetFeeds(&_MercuryUpkeep.TransactOpts, _feeds)
}

func (_MercuryUpkeep *MercuryUpkeepTransactorSession) SetFeeds(_feeds []string) (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.SetFeeds(&_MercuryUpkeep.TransactOpts, _feeds)
}

func (_MercuryUpkeep *MercuryUpkeepTransactor) SetParamKeys(opts *bind.TransactOpts, _feedParamKey string, _timeParamKey string) (*types.Transaction, error) {
	return _MercuryUpkeep.contract.Transact(opts, "setParamKeys", _feedParamKey, _timeParamKey)
}

func (_MercuryUpkeep *MercuryUpkeepSession) SetParamKeys(_feedParamKey string, _timeParamKey string) (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.SetParamKeys(&_MercuryUpkeep.TransactOpts, _feedParamKey, _timeParamKey)
}

func (_MercuryUpkeep *MercuryUpkeepTransactorSession) SetParamKeys(_feedParamKey string, _timeParamKey string) (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.SetParamKeys(&_MercuryUpkeep.TransactOpts, _feedParamKey, _timeParamKey)
}

func (_MercuryUpkeep *MercuryUpkeepTransactor) SetShouldRevertCallback(opts *bind.TransactOpts, value bool) (*types.Transaction, error) {
	return _MercuryUpkeep.contract.Transact(opts, "setShouldRevertCallback", value)
}

func (_MercuryUpkeep *MercuryUpkeepSession) SetShouldRevertCallback(value bool) (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.SetShouldRevertCallback(&_MercuryUpkeep.TransactOpts, value)
}

func (_MercuryUpkeep *MercuryUpkeepTransactorSession) SetShouldRevertCallback(value bool) (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.SetShouldRevertCallback(&_MercuryUpkeep.TransactOpts, value)
}

type MercuryUpkeepMercuryPerformEventIterator struct {
	Event *MercuryUpkeepMercuryPerformEvent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MercuryUpkeepMercuryPerformEventIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MercuryUpkeepMercuryPerformEvent)
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
		it.Event = new(MercuryUpkeepMercuryPerformEvent)
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

func (it *MercuryUpkeepMercuryPerformEventIterator) Error() error {
	return it.fail
}

func (it *MercuryUpkeepMercuryPerformEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MercuryUpkeepMercuryPerformEvent struct {
	Sender      common.Address
	BlockNumber *big.Int
	V0          []byte
	V1          []byte
	VerifiedV0  []byte
	VerifiedV1  []byte
	Ed          []byte
	Raw         types.Log
}

func (_MercuryUpkeep *MercuryUpkeepFilterer) FilterMercuryPerformEvent(opts *bind.FilterOpts, sender []common.Address, blockNumber []*big.Int) (*MercuryUpkeepMercuryPerformEventIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var blockNumberRule []interface{}
	for _, blockNumberItem := range blockNumber {
		blockNumberRule = append(blockNumberRule, blockNumberItem)
	}

	logs, sub, err := _MercuryUpkeep.contract.FilterLogs(opts, "MercuryPerformEvent", senderRule, blockNumberRule)
	if err != nil {
		return nil, err
	}
	return &MercuryUpkeepMercuryPerformEventIterator{contract: _MercuryUpkeep.contract, event: "MercuryPerformEvent", logs: logs, sub: sub}, nil
}

func (_MercuryUpkeep *MercuryUpkeepFilterer) WatchMercuryPerformEvent(opts *bind.WatchOpts, sink chan<- *MercuryUpkeepMercuryPerformEvent, sender []common.Address, blockNumber []*big.Int) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var blockNumberRule []interface{}
	for _, blockNumberItem := range blockNumber {
		blockNumberRule = append(blockNumberRule, blockNumberItem)
	}

	logs, sub, err := _MercuryUpkeep.contract.WatchLogs(opts, "MercuryPerformEvent", senderRule, blockNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MercuryUpkeepMercuryPerformEvent)
				if err := _MercuryUpkeep.contract.UnpackLog(event, "MercuryPerformEvent", log); err != nil {
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

func (_MercuryUpkeep *MercuryUpkeepFilterer) ParseMercuryPerformEvent(log types.Log) (*MercuryUpkeepMercuryPerformEvent, error) {
	event := new(MercuryUpkeepMercuryPerformEvent)
	if err := _MercuryUpkeep.contract.UnpackLog(event, "MercuryPerformEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_MercuryUpkeep *MercuryUpkeep) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _MercuryUpkeep.abi.Events["MercuryPerformEvent"].ID:
		return _MercuryUpkeep.ParseMercuryPerformEvent(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (MercuryUpkeepMercuryPerformEvent) Topic() common.Hash {
	return common.HexToHash("0x1c85d6186f024e964616014c8247533455ec5129a5095711202292f8a7ea1d54")
}

func (_MercuryUpkeep *MercuryUpkeep) Address() common.Address {
	return _MercuryUpkeep.address
}

type MercuryUpkeepInterface interface {
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

	FilterMercuryPerformEvent(opts *bind.FilterOpts, sender []common.Address, blockNumber []*big.Int) (*MercuryUpkeepMercuryPerformEventIterator, error)

	WatchMercuryPerformEvent(opts *bind.WatchOpts, sink chan<- *MercuryUpkeepMercuryPerformEvent, sender []common.Address, blockNumber []*big.Int) (event.Subscription, error)

	ParseMercuryPerformEvent(log types.Log) (*MercuryUpkeepMercuryPerformEvent, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
