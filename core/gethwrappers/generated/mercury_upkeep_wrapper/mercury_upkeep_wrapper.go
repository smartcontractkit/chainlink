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
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_useArbBlock\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"_isV02\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"_staging\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"_verify\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"FeedLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v0\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v1\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"verifiedV0\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"verifiedV1\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"ed\",\"type\":\"bytes\"}],\"name\":\"MercuryPerformEvent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"callbackReturnBool\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feeds\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"interval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"setCallbackReturnBool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"setShouldRevertCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"shouldRevertCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"staging\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbBlock\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b506040516200196838038062001968833981016040819052620000349162000304565b60008681556001869055600281905560038190556004558315156080528215620000c0576040805180820190915260098152680cccacac892c890caf60bb1b602082015260069062000087908262000418565b5060408051808201909152600b81526a313637b1b5a73ab6b132b960a91b6020820152600790620000b9908262000418565b506200011f565b6040805180820190915260078152666665656449447360c81b6020820152600690620000ed908262000418565b50604080518082019091526009815268074696d657374616d760bc1b60208201526007906200011d908262000418565b505b81156200018357604051806040016040528060405180608001604052806042815260200162001860604291398152602001604051806080016040528060428152602001620019266042913990526200017c90600590600262000217565b50620001db565b6040518060400160405280604051806080016040528060428152602001620018a2604291398152602001604051806080016040528060428152602001620018e4604291399052620001d990600590600262000217565b505b6008805463ff000000199215156101000261ff00199415159490941661ffff19909116179290921716630100000017905550620004e492505050565b82805482825590600052602060002090810192821562000262579160200282015b8281111562000262578251829062000251908262000418565b509160200191906001019062000238565b506200027092915062000274565b5090565b80821115620002705760006200028b828262000295565b5060010162000274565b508054620002a39062000389565b6000825580601f10620002b4575050565b601f016020900490600052602060002090810190620002d49190620002d7565b50565b5b80821115620002705760008155600101620002d8565b80518015158114620002ff57600080fd5b919050565b60008060008060008060c087890312156200031e57600080fd5b86519550602087015194506200033760408801620002ee565b93506200034760608801620002ee565b92506200035760808801620002ee565b91506200036760a08801620002ee565b90509295509295509295565b634e487b7160e01b600052604160045260246000fd5b600181811c908216806200039e57607f821691505b602082108103620003bf57634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200041357600081815260208120601f850160051c81016020861015620003ee5750805b601f850160051c820191505b818110156200040f57828155600101620003fa565b5050505b505050565b81516001600160401b0381111562000434576200043462000373565b6200044c8162000445845462000389565b84620003c5565b602080601f8311600181146200048457600084156200046b5750858301515b600019600386901b1c1916600185901b1785556200040f565b600085815260208120601f198616915b82811015620004b55788860151825594840194600190910190840162000494565b5085821015620004d45787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b60805161134b62000515600039600081816102c301528181610325015281816109eb0152610b1e015261134b6000f3fe608060405234801561001057600080fd5b50600436106101515760003560e01c806361bc221a116100cd578063947a36fb11610081578063c98f10b011610066578063c98f10b0146102ff578063d832d92f14610307578063fc735e991461030f57600080fd5b8063947a36fb146102ee578063afb28d1f146102f757600080fd5b80636e04ff0d116100b25780636e04ff0d146102ab57806386b728e2146102be578063917d895f146102e557600080fd5b806361bc221a146102995780636250a13a146102a257600080fd5b80634585e33b116101245780634b56a42e116101095780634b56a42e146101eb5780634bdb38621461020c5780635b48391a1461025257600080fd5b80634585e33b146101b65780634a5479f3146101cb57600080fd5b806302be021f14610156578063102d538b1461017e5780631d1970b7146101925780632cb158641461019f575b600080fd5b6008546101699062010000900460ff1681565b60405190151581526020015b60405180910390f35b600854610169906301000000900460ff1681565b6008546101699060ff1681565b6101a860035481565b604051908152602001610175565b6101c96101c4366004610bed565b610321565b005b6101de6101d9366004610c5f565b610808565b6040516101759190610ce6565b6101fe6101f9366004610e1a565b6108b4565b604051610175929190610f00565b6101c961021a366004610f23565b6008805491151562010000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffff909216919091179055565b6101c9610260366004610f23565b600880549115156301000000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffff909216919091179055565b6101a860045481565b6101a860005481565b6101fe6102b9366004610bed565b61098f565b6101697f000000000000000000000000000000000000000000000000000000000000000081565b6101a860025481565b6101a860015481565b6101de610aee565b6101de610afb565b610169610b08565b60085461016990610100900460ff1681565b60007f0000000000000000000000000000000000000000000000000000000000000000156103c057606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610395573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103b99190610f45565b90506103c3565b50435b6003546000036103d35760038190555b6000806103e284860186610e1a565b600285905560045491935091506103fa906001610f8d565b600455604080516020808201835260008083528351918201909352918252600854909190610100900460ff16156107745760085460ff16156105d7577360448b880c9f3b501af3f343da9284148bd7d77c73ffffffffffffffffffffffffffffffffffffffff16638e760afe8560008151811061047957610479610fa6565b60200260200101516040518263ffffffff1660e01b815260040161049d9190610ce6565b6000604051808303816000875af11580156104bc573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526105029190810190610fd5565b91507360448b880c9f3b501af3f343da9284148bd7d77c73ffffffffffffffffffffffffffffffffffffffff16638e760afe8560018151811061054757610547610fa6565b60200260200101516040518263ffffffff1660e01b815260040161056b9190610ce6565b6000604051808303816000875af115801561058a573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526105d09190810190610fd5565b9050610774565b7309dff56a4ff44e0f4436260a04f5cfa65636a48173ffffffffffffffffffffffffffffffffffffffff16638e760afe8560008151811061061a5761061a610fa6565b60200260200101516040518263ffffffff1660e01b815260040161063e9190610ce6565b6000604051808303816000875af115801561065d573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526106a39190810190610fd5565b91507309dff56a4ff44e0f4436260a04f5cfa65636a48173ffffffffffffffffffffffffffffffffffffffff16638e760afe856001815181106106e8576106e8610fa6565b60200260200101516040518263ffffffff1660e01b815260040161070c9190610ce6565b6000604051808303816000875af115801561072b573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526107719190810190610fd5565b90505b843373ffffffffffffffffffffffffffffffffffffffff167f1c85d6186f024e964616014c8247533455ec5129a5095711202292f8a7ea1d54866000815181106107c0576107c0610fa6565b6020026020010151876001815181106107db576107db610fa6565b60200260200101518686896040516107f795949392919061104c565b60405180910390a350505050505050565b6005818154811061081857600080fd5b906000526020600020016000915090508054610833906110b9565b80601f016020809104026020016040519081016040528092919081815260200182805461085f906110b9565b80156108ac5780601f10610881576101008083540402835291602001916108ac565b820191906000526020600020905b81548152906001019060200180831161088f57829003601f168201915b505050505081565b60085460009060609062010000900460ff1615610932576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f73686f756c6452657665727443616c6c6261636b20697320747275650000000060448201526064015b60405180910390fd5b6000848460405160200161094792919061110c565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529190526008546301000000900460ff1693509150505b9250929050565b6000606061099b610b08565b6109e7576000848481818080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250959750919550610988945050505050565b60007f000000000000000000000000000000000000000000000000000000000000000015610a8657606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610a5b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a7f9190610f45565b9050610a89565b50435b604080516c6400000000000000000000000060208201528151601481830301815260348201928390527f7ddd933e0000000000000000000000000000000000000000000000000000000090925261092991600691600591600791869190603801611274565b60068054610833906110b9565b60078054610833906110b9565b6000600354600003610b1a5750600190565b60007f000000000000000000000000000000000000000000000000000000000000000015610bb957606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610b8e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bb29190610f45565b9050610bbc565b50435b600054600354610bcc908361132b565b108015610be75750600154600254610be4908361132b565b10155b91505090565b60008060208385031215610c0057600080fd5b823567ffffffffffffffff80821115610c1857600080fd5b818501915085601f830112610c2c57600080fd5b813581811115610c3b57600080fd5b866020828501011115610c4d57600080fd5b60209290920196919550909350505050565b600060208284031215610c7157600080fd5b5035919050565b60005b83811015610c93578181015183820152602001610c7b565b50506000910152565b60008151808452610cb4816020860160208601610c78565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000610cf96020830184610c9c565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610d7657610d76610d00565b604052919050565b600067ffffffffffffffff821115610d9857610d98610d00565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f830112610dd557600080fd5b8135610de8610de382610d7e565b610d2f565b818152846020838601011115610dfd57600080fd5b816020850160208301376000918101602001919091529392505050565b60008060408385031215610e2d57600080fd5b823567ffffffffffffffff80821115610e4557600080fd5b818501915085601f830112610e5957600080fd5b8135602082821115610e6d57610e6d610d00565b8160051b610e7c828201610d2f565b928352848101820192828101908a851115610e9657600080fd5b83870192505b84831015610ed257823586811115610eb45760008081fd5b610ec28c86838b0101610dc4565b8352509183019190830190610e9c565b9750505086013592505080821115610ee957600080fd5b50610ef685828601610dc4565b9150509250929050565b8215158152604060208201526000610f1b6040830184610c9c565b949350505050565b600060208284031215610f3557600080fd5b81358015158114610cf957600080fd5b600060208284031215610f5757600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b80820180821115610fa057610fa0610f5e565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600060208284031215610fe757600080fd5b815167ffffffffffffffff811115610ffe57600080fd5b8201601f8101841361100f57600080fd5b805161101d610de382610d7e565b81815285602083850101111561103257600080fd5b611043826020830160208601610c78565b95945050505050565b60a08152600061105f60a0830188610c9c565b82810360208401526110718188610c9c565b905082810360408401526110858187610c9c565b905082810360608401526110998186610c9c565b905082810360808401526110ad8185610c9c565b98975050505050505050565b600181811c908216806110cd57607f821691505b602082108103611106577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b83811015611181577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa088870301855261116f868351610c9c565b95509382019390820190600101611135565b5050858403818701525050506110438185610c9c565b8054600090600181811c90808316806111b157607f831692505b602080841082036111eb577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b838852818015611202576001811461123a57611268565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008616828a01528185151560051b8a01019650611268565b876000528160002060005b868110156112605781548b8201850152908501908301611245565b8a0183019750505b50505050505092915050565b60a08152600061128760a0830188611197565b6020838203818501528188548084528284019150828160051b8501018a6000528360002060005b838110156112f9577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08784030185526112e78383611197565b948601949250600191820191016112ae565b5050868103604088015261130d818b611197565b94505050505084606084015282810360808401526110ad8185610c9c565b81810381811115610fa057610fa0610f5e56fea164736f6c6343000810000a307866373533653132303164353461633934646664393333346335343235363266663765343239393334313961363631323631643031306166306362666434653334307834353534343832643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030307834323534343332643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030307836393632653632396333613066356237653365393239346230633238336339623230663934663163383963386261386331656534363530373338663230666232",
}

var MercuryUpkeepABI = MercuryUpkeepMetaData.ABI

var MercuryUpkeepBin = MercuryUpkeepMetaData.Bin

func DeployMercuryUpkeep(auth *bind.TransactOpts, backend bind.ContractBackend, _testRange *big.Int, _interval *big.Int, _useArbBlock bool, _isV02 bool, _staging bool, _verify bool) (common.Address, *types.Transaction, *MercuryUpkeep, error) {
	parsed, err := MercuryUpkeepMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MercuryUpkeepBin), backend, _testRange, _interval, _useArbBlock, _isV02, _staging, _verify)
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

func (_MercuryUpkeep *MercuryUpkeepTransactor) SetCallbackReturnBool(opts *bind.TransactOpts, value bool) (*types.Transaction, error) {
	return _MercuryUpkeep.contract.Transact(opts, "setCallbackReturnBool", value)
}

func (_MercuryUpkeep *MercuryUpkeepSession) SetCallbackReturnBool(value bool) (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.SetCallbackReturnBool(&_MercuryUpkeep.TransactOpts, value)
}

func (_MercuryUpkeep *MercuryUpkeepTransactorSession) SetCallbackReturnBool(value bool) (*types.Transaction, error) {
	return _MercuryUpkeep.Contract.SetCallbackReturnBool(&_MercuryUpkeep.TransactOpts, value)
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

	SetCallbackReturnBool(opts *bind.TransactOpts, value bool) (*types.Transaction, error)

	SetShouldRevertCallback(opts *bind.TransactOpts, value bool) (*types.Transaction, error)

	FilterMercuryPerformEvent(opts *bind.FilterOpts, sender []common.Address, blockNumber []*big.Int) (*MercuryUpkeepMercuryPerformEventIterator, error)

	WatchMercuryPerformEvent(opts *bind.WatchOpts, sink chan<- *MercuryUpkeepMercuryPerformEvent, sender []common.Address, blockNumber []*big.Int) (event.Subscription, error)

	ParseMercuryPerformEvent(log types.Log) (*MercuryUpkeepMercuryPerformEvent, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
