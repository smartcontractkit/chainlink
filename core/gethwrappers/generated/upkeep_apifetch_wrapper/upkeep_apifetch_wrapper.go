// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package upkeep_apifetch_wrapper

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
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated"
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
)

var UpkeepAPIFetchMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"url\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"string[]\",\"name\":\"jsonFields\",\"type\":\"string[]\"},{\"internalType\":\"bytes4\",\"name\":\"callbackSelector\",\"type\":\"bytes4\"}],\"name\":\"ChainlinkAPIFetch\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"id\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"abilities\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"types\",\"type\":\"string\"}],\"name\":\"PokemonUpkeep\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"abilities\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"string[]\",\"name\":\"values\",\"type\":\"string[]\"},{\"internalType\":\"uint256\",\"name\":\"statusCode\",\"type\":\"uint256\"}],\"name\":\"callback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"fields\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"id\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"interval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pokemon\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_testRange\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"input\",\"type\":\"string\"}],\"name\":\"setURLs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRange\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"types\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"url\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50604051620016c1380380620016c1833981016040819052620000349162000203565b60008281556001829055600381815543600255600482905560059182556040805160c08101825260808101928352620b9a5960ea1b60a082015291825280518082018252928352642e6e616d6560d81b6020848101919091528083019390935280518082018252601e81527f5b2e6162696c69746965735b5d207c202e6162696c6974792e6e616d655d00008185015282820152805160608181019092526021808252929391840192909162001677908301399052620000f990600b9060046200012c565b506040518060600160405280602981526020016200169860299139600690620001239082620002cd565b50505062000399565b82805482825590600052602060002090810192821562000177579160200282015b82811115620001775782518290620001669082620002cd565b50916020019190600101906200014d565b506200018592915062000189565b5090565b8082111562000185576000620001a08282620001aa565b5060010162000189565b508054620001b8906200023e565b6000825580601f10620001c9575050565b601f016020900490600052602060002090810190620001e99190620001ec565b50565b5b80821115620001855760008155600101620001ed565b600080604083850312156200021757600080fd5b505080516020909101519092909150565b634e487b7160e01b600052604160045260246000fd5b600181811c908216806200025357607f821691505b6020821081036200027457634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620002c857600081815260208120601f850160051c81016020861015620002a35750805b601f850160051c820191505b81811015620002c457828155600101620002af565b5050505b505050565b81516001600160401b03811115620002e957620002e962000228565b6200030181620002fa84546200023e565b846200027a565b602080601f831160018114620003395760008415620003205750858301515b600019600386901b1c1916600185901b178555620002c4565b600085815260208120601f198616915b828110156200036a5788860151825594840194600190910190840162000349565b5085821015620003895787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6112ce80620003a96000396000f3fe608060405234801561001057600080fd5b50600436106101365760003560e01c806371106628116100b2578063af640d0f11610081578063ccbff97911610066578063ccbff9791461023c578063d832d92f14610244578063e0a04d711461025c57600080fd5b8063af640d0f14610221578063b772d70a1461022957600080fd5b806371106628146101f3578063806b984f14610206578063917d895f1461020f578063947a36fb1461021857600080fd5b80634585e33b1161010957806361bc221a116100ee57806361bc221a146101c05780636250a13a146101c95780636e04ff0d146101d257600080fd5b80634585e33b146101a55780635600f04f146101b857600080fd5b8063191378561461013b5780631e34c585146101645780632cb1586414610186578063362ff8ac1461019d575b600080fd5b61014e6101493660046108b2565b610264565b60405161015b9190610945565b60405180910390f35b61018461017236600461095f565b60009182556001556004819055600555565b005b61018f60045481565b60405190815260200161015b565b61014e610310565b6101846101b33660046109c3565b61031d565b61014e6103f6565b61018f60055481565b61018f60005481565b6101e56101e03660046109c3565b610403565b60405161015b929190610a05565b610184610201366004610afa565b610545565b61018f60025481565b61018f60035481565b61018f60015481565b61014e610555565b6101e5610237366004610b2f565b610562565b61014e610719565b61024c610726565b604051901515815260200161015b565b61014e610768565b600b818154811061027457600080fd5b90600052602060002001600091509050805461028f90610bd1565b80601f01602080910402602001604051908101604052809291908181526020018280546102bb90610bd1565b80156103085780601f106102dd57610100808354040283529160200191610308565b820191906000526020600020905b8154815290600101906020018083116102eb57829003601f168201915b505050505081565b6008805461028f90610bd1565b60045460000361032c57436004555b4360025560055461033e906001610c53565b600555600080808061035285870187610c6b565b9296509094509250905060076103688582610d67565b5060086103758482610d67565b5060096103828382610d67565b50600a61038f8282610d67565b503273ffffffffffffffffffffffffffffffffffffffff167f7e09e0773d481e69887f4a2562b67d3ba3b5a4878177596081010882c3bd4038600760086009600a6040516103e09493929190610f1c565b60405180910390a2505060025460035550505050565b6006805461028f90610bd1565b6000606061040f610726565b61045b576000848481818080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525095975091955061053e945050505050565b6000610474600554600161046f9190610c53565b610775565b9050600060068260405160200161048c929190610f74565b60405160208183030381529060405290508086866040516020016104b1929190611019565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527fda8b32140000000000000000000000000000000000000000000000000000000082526105359291600b907fb772d70a0000000000000000000000000000000000000000000000000000000090600401611066565b60405180910390fd5b9250929050565b60066105518282610d67565b5050565b6007805461028f90610bd1565b6000606060008585600081811061057b5761057b61113a565b905060200281019061058d9190611169565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920182905250939450899250889150600190508181106105d9576105d961113a565b90506020028101906105eb9190611169565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201829052509394508a9250899150600290508181106106375761063761113a565b90506020028101906106499190611169565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201829052509394508b92508a9150600390508181106106955761069561113a565b90506020028101906106a79190611169565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250506040519293506001926106f892508791508690869086906020016111ce565b60405160208183030381529060405295509550505050509550959350505050565b600a805461028f90610bd1565b60006004546000036107385750600190565b600054600454610748904361121b565b1080156107635750600154600254610760904361121b565b10155b905090565b6009805461028f90610bd1565b6060816000036107b857505060408051808201909152600181527f3000000000000000000000000000000000000000000000000000000000000000602082015290565b8160005b81156107e257806107cc81611232565b91506107db9050600a83611299565b91506107bc565b60008167ffffffffffffffff8111156107fd576107fd610a20565b6040519080825280601f01601f191660200182016040528015610827576020820181803683370190505b5090505b84156108aa5761083c60018361121b565b9150610849600a866112ad565b610854906030610c53565b60f81b8183815181106108695761086961113a565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053506108a3600a86611299565b945061082b565b949350505050565b6000602082840312156108c457600080fd5b5035919050565b60005b838110156108e65781810151838201526020016108ce565b838111156108f5576000848401525b50505050565b600081518084526109138160208601602086016108cb565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60208152600061095860208301846108fb565b9392505050565b6000806040838503121561097257600080fd5b50508035926020909101359150565b60008083601f84011261099357600080fd5b50813567ffffffffffffffff8111156109ab57600080fd5b60208301915083602082850101111561053e57600080fd5b600080602083850312156109d657600080fd5b823567ffffffffffffffff8111156109ed57600080fd5b6109f985828601610981565b90969095509350505050565b82151581526040602082015260006108aa60408301846108fb565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600082601f830112610a6057600080fd5b813567ffffffffffffffff80821115610a7b57610a7b610a20565b604051601f83017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908282118183101715610ac157610ac1610a20565b81604052838152866020858801011115610ada57600080fd5b836020870160208301376000602085830101528094505050505092915050565b600060208284031215610b0c57600080fd5b813567ffffffffffffffff811115610b2357600080fd5b6108aa84828501610a4f565b600080600080600060608688031215610b4757600080fd5b853567ffffffffffffffff80821115610b5f57600080fd5b610b6b89838a01610981565b90975095506020880135915080821115610b8457600080fd5b818801915088601f830112610b9857600080fd5b813581811115610ba757600080fd5b8960208260051b8501011115610bbc57600080fd5b96999598505060200195604001359392505050565b600181811c90821680610be557607f821691505b602082108103610c1e577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60008219821115610c6657610c66610c24565b500190565b60008060008060808587031215610c8157600080fd5b843567ffffffffffffffff80821115610c9957600080fd5b610ca588838901610a4f565b95506020870135915080821115610cbb57600080fd5b610cc788838901610a4f565b94506040870135915080821115610cdd57600080fd5b610ce988838901610a4f565b93506060870135915080821115610cff57600080fd5b50610d0c87828801610a4f565b91505092959194509250565b601f821115610d6257600081815260208120601f850160051c81016020861015610d3f5750805b601f850160051c820191505b81811015610d5e57828155600101610d4b565b5050505b505050565b815167ffffffffffffffff811115610d8157610d81610a20565b610d9581610d8f8454610bd1565b84610d18565b602080601f831160018114610de85760008415610db25750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555610d5e565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015610e3557888601518255948401946001909101908401610e16565b5085821015610e7157878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b60008154610e8e81610bd1565b808552602060018381168015610eab5760018114610ee357610f11565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b8901019550610f11565b866000528260002060005b85811015610f095781548a8201860152908301908401610eee565b890184019650505b505050505092915050565b608081526000610f2f6080830187610e81565b8281036020840152610f418187610e81565b90508281036040840152610f558186610e81565b90508281036060840152610f698185610e81565b979650505050505050565b6000808454610f8281610bd1565b60018281168015610f9a5760018114610fcd57610ffc565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0084168752821515830287019450610ffc565b8860005260208060002060005b85811015610ff35781548a820152908401908201610fda565b50505082870194505b5050505083516110108183602088016108cb565b01949350505050565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b60808152600061107960808301876108fb565b60208382038185015261108c82886108fb565b915083820360408501528186548084528284019150828160051b850101886000528360002060005b838110156110ff577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08784030185526110ed8383610e81565b948601949250600191820191016110b4565b505080955050505050507fffffffff000000000000000000000000000000000000000000000000000000008316606083015295945050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261119e57600080fd5b83018035915067ffffffffffffffff8211156111b957600080fd5b60200191503681900382131561053e57600080fd5b6080815260006111e160808301876108fb565b82810360208401526111f381876108fb565b9050828103604084015261120781866108fb565b90508281036060840152610f6981856108fb565b60008282101561122d5761122d610c24565b500390565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361126357611263610c24565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000826112a8576112a861126a565b500490565b6000826112bc576112bc61126a565b50069056fea164736f6c634300080f000a5b2e74797065735b5d207c202e747970652e6e616d655d7c6a6f696e28222c222968747470733a2f2f706f6b656170692e636f2f6170692f7b76657273696f6e7d2f706f6b656d6f6e2f",
}

var UpkeepAPIFetchABI = UpkeepAPIFetchMetaData.ABI

var UpkeepAPIFetchBin = UpkeepAPIFetchMetaData.Bin

func DeployUpkeepAPIFetch(auth *bind.TransactOpts, backend bind.ContractBackend, _testRange *big.Int, _interval *big.Int) (common.Address, *types.Transaction, *UpkeepAPIFetch, error) {
	parsed, err := UpkeepAPIFetchMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UpkeepAPIFetchBin), backend, _testRange, _interval)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UpkeepAPIFetch{UpkeepAPIFetchCaller: UpkeepAPIFetchCaller{contract: contract}, UpkeepAPIFetchTransactor: UpkeepAPIFetchTransactor{contract: contract}, UpkeepAPIFetchFilterer: UpkeepAPIFetchFilterer{contract: contract}}, nil
}

type UpkeepAPIFetch struct {
	address common.Address
	abi     abi.ABI
	UpkeepAPIFetchCaller
	UpkeepAPIFetchTransactor
	UpkeepAPIFetchFilterer
}

type UpkeepAPIFetchCaller struct {
	contract *bind.BoundContract
}

type UpkeepAPIFetchTransactor struct {
	contract *bind.BoundContract
}

type UpkeepAPIFetchFilterer struct {
	contract *bind.BoundContract
}

type UpkeepAPIFetchSession struct {
	Contract     *UpkeepAPIFetch
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type UpkeepAPIFetchCallerSession struct {
	Contract *UpkeepAPIFetchCaller
	CallOpts bind.CallOpts
}

type UpkeepAPIFetchTransactorSession struct {
	Contract     *UpkeepAPIFetchTransactor
	TransactOpts bind.TransactOpts
}

type UpkeepAPIFetchRaw struct {
	Contract *UpkeepAPIFetch
}

type UpkeepAPIFetchCallerRaw struct {
	Contract *UpkeepAPIFetchCaller
}

type UpkeepAPIFetchTransactorRaw struct {
	Contract *UpkeepAPIFetchTransactor
}

func NewUpkeepAPIFetch(address common.Address, backend bind.ContractBackend) (*UpkeepAPIFetch, error) {
	abi, err := abi.JSON(strings.NewReader(UpkeepAPIFetchABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindUpkeepAPIFetch(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UpkeepAPIFetch{address: address, abi: abi, UpkeepAPIFetchCaller: UpkeepAPIFetchCaller{contract: contract}, UpkeepAPIFetchTransactor: UpkeepAPIFetchTransactor{contract: contract}, UpkeepAPIFetchFilterer: UpkeepAPIFetchFilterer{contract: contract}}, nil
}

func NewUpkeepAPIFetchCaller(address common.Address, caller bind.ContractCaller) (*UpkeepAPIFetchCaller, error) {
	contract, err := bindUpkeepAPIFetch(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UpkeepAPIFetchCaller{contract: contract}, nil
}

func NewUpkeepAPIFetchTransactor(address common.Address, transactor bind.ContractTransactor) (*UpkeepAPIFetchTransactor, error) {
	contract, err := bindUpkeepAPIFetch(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UpkeepAPIFetchTransactor{contract: contract}, nil
}

func NewUpkeepAPIFetchFilterer(address common.Address, filterer bind.ContractFilterer) (*UpkeepAPIFetchFilterer, error) {
	contract, err := bindUpkeepAPIFetch(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UpkeepAPIFetchFilterer{contract: contract}, nil
}

func bindUpkeepAPIFetch(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UpkeepAPIFetchABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_UpkeepAPIFetch *UpkeepAPIFetchRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpkeepAPIFetch.Contract.UpkeepAPIFetchCaller.contract.Call(opts, result, method, params...)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepAPIFetch.Contract.UpkeepAPIFetchTransactor.contract.Transfer(opts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpkeepAPIFetch.Contract.UpkeepAPIFetchTransactor.contract.Transact(opts, method, params...)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpkeepAPIFetch.Contract.contract.Call(opts, result, method, params...)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpkeepAPIFetch.Contract.contract.Transfer(opts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpkeepAPIFetch.Contract.contract.Transact(opts, method, params...)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) Abilities(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "abilities")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) Abilities() (string, error) {
	return _UpkeepAPIFetch.Contract.Abilities(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) Abilities() (string, error) {
	return _UpkeepAPIFetch.Contract.Abilities(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) Callback(opts *bind.CallOpts, extraData []byte, values []string, statusCode *big.Int) (bool, []byte, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "callback", extraData, values, statusCode)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) Callback(extraData []byte, values []string, statusCode *big.Int) (bool, []byte, error) {
	return _UpkeepAPIFetch.Contract.Callback(&_UpkeepAPIFetch.CallOpts, extraData, values, statusCode)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) Callback(extraData []byte, values []string, statusCode *big.Int) (bool, []byte, error) {
	return _UpkeepAPIFetch.Contract.Callback(&_UpkeepAPIFetch.CallOpts, extraData, values, statusCode)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) CheckUpkeep(opts *bind.CallOpts, data []byte) (bool, []byte, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "checkUpkeep", data)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) CheckUpkeep(data []byte) (bool, []byte, error) {
	return _UpkeepAPIFetch.Contract.CheckUpkeep(&_UpkeepAPIFetch.CallOpts, data)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) CheckUpkeep(data []byte) (bool, []byte, error) {
	return _UpkeepAPIFetch.Contract.CheckUpkeep(&_UpkeepAPIFetch.CallOpts, data)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) Counter() (*big.Int, error) {
	return _UpkeepAPIFetch.Contract.Counter(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) Counter() (*big.Int, error) {
	return _UpkeepAPIFetch.Contract.Counter(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) Eligible(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "eligible")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) Eligible() (bool, error) {
	return _UpkeepAPIFetch.Contract.Eligible(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) Eligible() (bool, error) {
	return _UpkeepAPIFetch.Contract.Eligible(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) Fields(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "fields", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) Fields(arg0 *big.Int) (string, error) {
	return _UpkeepAPIFetch.Contract.Fields(&_UpkeepAPIFetch.CallOpts, arg0)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) Fields(arg0 *big.Int) (string, error) {
	return _UpkeepAPIFetch.Contract.Fields(&_UpkeepAPIFetch.CallOpts, arg0)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) Id(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "id")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) Id() (string, error) {
	return _UpkeepAPIFetch.Contract.Id(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) Id() (string, error) {
	return _UpkeepAPIFetch.Contract.Id(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) InitialBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "initialBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) InitialBlock() (*big.Int, error) {
	return _UpkeepAPIFetch.Contract.InitialBlock(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) InitialBlock() (*big.Int, error) {
	return _UpkeepAPIFetch.Contract.InitialBlock(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) Interval(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "interval")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) Interval() (*big.Int, error) {
	return _UpkeepAPIFetch.Contract.Interval(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) Interval() (*big.Int, error) {
	return _UpkeepAPIFetch.Contract.Interval(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) LastBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "lastBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) LastBlock() (*big.Int, error) {
	return _UpkeepAPIFetch.Contract.LastBlock(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) LastBlock() (*big.Int, error) {
	return _UpkeepAPIFetch.Contract.LastBlock(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) Pokemon(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "pokemon")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) Pokemon() (string, error) {
	return _UpkeepAPIFetch.Contract.Pokemon(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) Pokemon() (string, error) {
	return _UpkeepAPIFetch.Contract.Pokemon(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "previousPerformBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) PreviousPerformBlock() (*big.Int, error) {
	return _UpkeepAPIFetch.Contract.PreviousPerformBlock(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) PreviousPerformBlock() (*big.Int, error) {
	return _UpkeepAPIFetch.Contract.PreviousPerformBlock(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) TestRange(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "testRange")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) TestRange() (*big.Int, error) {
	return _UpkeepAPIFetch.Contract.TestRange(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) TestRange() (*big.Int, error) {
	return _UpkeepAPIFetch.Contract.TestRange(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) Types(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "types")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) Types() (string, error) {
	return _UpkeepAPIFetch.Contract.Types(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) Types() (string, error) {
	return _UpkeepAPIFetch.Contract.Types(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCaller) Url(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _UpkeepAPIFetch.contract.Call(opts, &out, "url")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) Url() (string, error) {
	return _UpkeepAPIFetch.Contract.Url(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchCallerSession) Url() (string, error) {
	return _UpkeepAPIFetch.Contract.Url(&_UpkeepAPIFetch.CallOpts)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _UpkeepAPIFetch.contract.Transact(opts, "performUpkeep", performData)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _UpkeepAPIFetch.Contract.PerformUpkeep(&_UpkeepAPIFetch.TransactOpts, performData)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _UpkeepAPIFetch.Contract.PerformUpkeep(&_UpkeepAPIFetch.TransactOpts, performData)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchTransactor) SetConfig(opts *bind.TransactOpts, _testRange *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _UpkeepAPIFetch.contract.Transact(opts, "setConfig", _testRange, _interval)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) SetConfig(_testRange *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _UpkeepAPIFetch.Contract.SetConfig(&_UpkeepAPIFetch.TransactOpts, _testRange, _interval)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchTransactorSession) SetConfig(_testRange *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _UpkeepAPIFetch.Contract.SetConfig(&_UpkeepAPIFetch.TransactOpts, _testRange, _interval)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchTransactor) SetURLs(opts *bind.TransactOpts, input string) (*types.Transaction, error) {
	return _UpkeepAPIFetch.contract.Transact(opts, "setURLs", input)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchSession) SetURLs(input string) (*types.Transaction, error) {
	return _UpkeepAPIFetch.Contract.SetURLs(&_UpkeepAPIFetch.TransactOpts, input)
}

func (_UpkeepAPIFetch *UpkeepAPIFetchTransactorSession) SetURLs(input string) (*types.Transaction, error) {
	return _UpkeepAPIFetch.Contract.SetURLs(&_UpkeepAPIFetch.TransactOpts, input)
}

type UpkeepAPIFetchPokemonUpkeepIterator struct {
	Event *UpkeepAPIFetchPokemonUpkeep

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *UpkeepAPIFetchPokemonUpkeepIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpkeepAPIFetchPokemonUpkeep)
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
		it.Event = new(UpkeepAPIFetchPokemonUpkeep)
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

func (it *UpkeepAPIFetchPokemonUpkeepIterator) Error() error {
	return it.fail
}

func (it *UpkeepAPIFetchPokemonUpkeepIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type UpkeepAPIFetchPokemonUpkeep struct {
	From      common.Address
	Id        string
	Name      string
	Abilities string
	Types     string
	Raw       types.Log
}

func (_UpkeepAPIFetch *UpkeepAPIFetchFilterer) FilterPokemonUpkeep(opts *bind.FilterOpts, from []common.Address) (*UpkeepAPIFetchPokemonUpkeepIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _UpkeepAPIFetch.contract.FilterLogs(opts, "PokemonUpkeep", fromRule)
	if err != nil {
		return nil, err
	}
	return &UpkeepAPIFetchPokemonUpkeepIterator{contract: _UpkeepAPIFetch.contract, event: "PokemonUpkeep", logs: logs, sub: sub}, nil
}

func (_UpkeepAPIFetch *UpkeepAPIFetchFilterer) WatchPokemonUpkeep(opts *bind.WatchOpts, sink chan<- *UpkeepAPIFetchPokemonUpkeep, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _UpkeepAPIFetch.contract.WatchLogs(opts, "PokemonUpkeep", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(UpkeepAPIFetchPokemonUpkeep)
				if err := _UpkeepAPIFetch.contract.UnpackLog(event, "PokemonUpkeep", log); err != nil {
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

func (_UpkeepAPIFetch *UpkeepAPIFetchFilterer) ParsePokemonUpkeep(log types.Log) (*UpkeepAPIFetchPokemonUpkeep, error) {
	event := new(UpkeepAPIFetchPokemonUpkeep)
	if err := _UpkeepAPIFetch.contract.UnpackLog(event, "PokemonUpkeep", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_UpkeepAPIFetch *UpkeepAPIFetch) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _UpkeepAPIFetch.abi.Events["PokemonUpkeep"].ID:
		return _UpkeepAPIFetch.ParsePokemonUpkeep(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (UpkeepAPIFetchPokemonUpkeep) Topic() common.Hash {
	return common.HexToHash("0x7e09e0773d481e69887f4a2562b67d3ba3b5a4878177596081010882c3bd4038")
}

func (_UpkeepAPIFetch *UpkeepAPIFetch) Address() common.Address {
	return _UpkeepAPIFetch.address
}

type UpkeepAPIFetchInterface interface {
	Abilities(opts *bind.CallOpts) (string, error)

	Callback(opts *bind.CallOpts, extraData []byte, values []string, statusCode *big.Int) (bool, []byte, error)

	CheckUpkeep(opts *bind.CallOpts, data []byte) (bool, []byte, error)

	Counter(opts *bind.CallOpts) (*big.Int, error)

	Eligible(opts *bind.CallOpts) (bool, error)

	Fields(opts *bind.CallOpts, arg0 *big.Int) (string, error)

	Id(opts *bind.CallOpts) (string, error)

	InitialBlock(opts *bind.CallOpts) (*big.Int, error)

	Interval(opts *bind.CallOpts) (*big.Int, error)

	LastBlock(opts *bind.CallOpts) (*big.Int, error)

	Pokemon(opts *bind.CallOpts) (string, error)

	PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error)

	TestRange(opts *bind.CallOpts) (*big.Int, error)

	Types(opts *bind.CallOpts) (string, error)

	Url(opts *bind.CallOpts) (string, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, _testRange *big.Int, _interval *big.Int) (*types.Transaction, error)

	SetURLs(opts *bind.TransactOpts, input string) (*types.Transaction, error)

	FilterPokemonUpkeep(opts *bind.FilterOpts, from []common.Address) (*UpkeepAPIFetchPokemonUpkeepIterator, error)

	WatchPokemonUpkeep(opts *bind.WatchOpts, sink chan<- *UpkeepAPIFetchPokemonUpkeep, from []common.Address) (event.Subscription, error)

	ParsePokemonUpkeep(log types.Log) (*UpkeepAPIFetchPokemonUpkeep, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
