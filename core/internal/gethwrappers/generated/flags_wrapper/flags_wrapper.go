// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package flags_wrapper

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// FlagsABI is the input ABI used to generate the binding from.
const FlagsABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"racAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"AddedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"subject\",\"type\":\"address\"}],\"name\":\"FlagLowered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"subject\",\"type\":\"address\"}],\"name\":\"FlagRaised\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"RaisingAccessControllerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"RemovedAccess\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"addAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"subject\",\"type\":\"address\"}],\"name\":\"getFlag\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"subjects\",\"type\":\"address[]\"}],\"name\":\"getFlags\",\"outputs\":[{\"internalType\":\"bool[]\",\"name\":\"\",\"type\":\"bool[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_calldata\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"subjects\",\"type\":\"address[]\"}],\"name\":\"lowerFlags\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"subject\",\"type\":\"address\"}],\"name\":\"raiseFlag\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"subjects\",\"type\":\"address[]\"}],\"name\":\"raiseFlags\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"raisingAccessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"removeAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"racAddress\",\"type\":\"address\"}],\"name\":\"setRaisingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// FlagsBin is the compiled bytecode used for deploying new contracts.
var FlagsBin = "0x608060405234801561001057600080fd5b506040516115cd3803806115cd8339818101604052602081101561003357600080fd5b5051600080546001600160a01b031916331790556001805460ff60a01b1916600160a01b17905561006c816001600160e01b0361007216565b5061013a565b6000546001600160a01b031633146100d1576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b6003546001600160a01b03908116908216811461013657600380546001600160a01b0319166001600160a01b0384811691821790925560405190918316907fbaf9ea078655a4fffefd08f9435677bbc91e457a6d015fe7de1d0e68b8802cac90600090a35b5050565b611484806101496000396000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c80637d723cac11610097578063a118f24911610066578063a118f24914610468578063d74af2631461049b578063dc7f0124146104ce578063f2fde38b146104d657610100565b80637d723cac146103655780638038e4a1146104255780638823da6c1461042d5780638da5cb5b1461046057610100565b8063517e89fe116100d3578063517e89fe146101f75780636b14daf81461022a578063760bc82d146102ed57806379ba50971461035d57610100565b80630a75698314610105578063282865961461010f5780632e1d859c1461017f578063357e47fe146101b0575b600080fd5b61010d610509565b005b61010d6004803603602081101561012557600080fd5b81019060208101813564010000000081111561014057600080fd5b82018360208201111561015257600080fd5b8035906020019184602083028401116401000000008311171561017457600080fd5b509092509050610606565b610187610761565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b6101e3600480360360208110156101c657600080fd5b503573ffffffffffffffffffffffffffffffffffffffff1661077d565b604080519115158252519081900360200190f35b61010d6004803603602081101561020d57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610857565b6101e36004803603604081101561024057600080fd5b73ffffffffffffffffffffffffffffffffffffffff823516919081019060408101602082013564010000000081111561027857600080fd5b82018360208201111561028a57600080fd5b803590602001918460018302840111640100000000831117156102ac57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610978945050505050565b61010d6004803603602081101561030357600080fd5b81019060208101813564010000000081111561031e57600080fd5b82018360208201111561033057600080fd5b8035906020019184602083028401116401000000008311171561035257600080fd5b5090925090506109ab565b61010d610a62565b6103d56004803603602081101561037b57600080fd5b81019060208101813564010000000081111561039657600080fd5b8201836020820111156103a857600080fd5b803590602001918460208302840111640100000000831117156103ca57600080fd5b509092509050610b64565b60408051602080825283518183015283519192839290830191858101910280838360005b838110156104115781810151838201526020016103f9565b505050509050019250505060405180910390f35b61010d610d04565b61010d6004803603602081101561044357600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610e16565b610187610f4e565b61010d6004803603602081101561047e57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610f6a565b61010d600480360360208110156104b157600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166110a3565b6101e361111f565b61010d600480360360208110156104ec57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16611140565b60005473ffffffffffffffffffffffffffffffffffffffff16331461058f57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b60015474010000000000000000000000000000000000000000900460ff161561060457600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1690556040517f3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f53963890600090a15b565b60005473ffffffffffffffffffffffffffffffffffffffff16331461068c57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b60005b8181101561075c5760008383838181106106a557fe5b6020908102929092013573ffffffffffffffffffffffffffffffffffffffff16600081815260049093526040909220549192505060ff16156107535773ffffffffffffffffffffffffffffffffffffffff811660008181526004602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055517fd86728e2e5cbaa28c1d357b5fbccc9c1ab0add09950eb7cac42df9acb24c4bc89190a25b5060010161068f565b505050565b60035473ffffffffffffffffffffffffffffffffffffffff1681565b60006107c0336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061097892505050565b61082b57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b5073ffffffffffffffffffffffffffffffffffffffff1660009081526004602052604090205460ff1690565b60005473ffffffffffffffffffffffffffffffffffffffff1633146108dd57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b60035473ffffffffffffffffffffffffffffffffffffffff908116908216811461097457600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff84811691821790925560405190918316907fbaf9ea078655a4fffefd08f9435677bbc91e457a6d015fe7de1d0e68b8802cac90600090a35b5050565b6000610984838361123c565b806109a4575073ffffffffffffffffffffffffffffffffffffffff831632145b9392505050565b6109b3611291565b610a1e57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f4e6f7420616c6c6f77656420746f20726169736520666c616773000000000000604482015290519081900360640190fd5b60005b8181101561075c57610a5a838383818110610a3857fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff166113aa565b600101610a21565b60015473ffffffffffffffffffffffffffffffffffffffff163314610ae857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6060610ba7336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061097892505050565b610c1257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b60608267ffffffffffffffff81118015610c2b57600080fd5b50604051908082528060200260200182016040528015610c55578160200160208202803683370190505b50905060005b83811015610cfc5760046000868684818110610c7357fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16828281518110610ce457fe5b91151560209283029190910190910152600101610c5b565b509392505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610d8a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b60015474010000000000000000000000000000000000000000900460ff1661060457600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff16740100000000000000000000000000000000000000001790556040517faebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c348090600090a1565b60005473ffffffffffffffffffffffffffffffffffffffff163314610e9c57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff811660009081526002602052604090205460ff1615610f4b5773ffffffffffffffffffffffffffffffffffffffff811660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055815192835290517f3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d19281900390910190a15b50565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b60005473ffffffffffffffffffffffffffffffffffffffff163314610ff057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff811660009081526002602052604090205460ff16610f4b5773ffffffffffffffffffffffffffffffffffffffff811660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055815192835290517f87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db49281900390910190a150565b6110ab611291565b61111657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f4e6f7420616c6c6f77656420746f20726169736520666c616773000000000000604482015290519081900360640190fd5b610f4b816113aa565b60015474010000000000000000000000000000000000000000900460ff1681565b60005473ffffffffffffffffffffffffffffffffffffffff1633146111c657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b73ffffffffffffffffffffffffffffffffffffffff821660009081526002602052604081205460ff16806109a457505060015474010000000000000000000000000000000000000000900460ff161592915050565b6000805473ffffffffffffffffffffffffffffffffffffffff163314806113a55750600354604080517f6b14daf8000000000000000000000000000000000000000000000000000000008152336004820181815260248301938452366044840181905273ffffffffffffffffffffffffffffffffffffffff90951694636b14daf894929360009391929190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b15801561137857600080fd5b505afa15801561138c573d6000803e3d6000fd5b505050506040513d60208110156113a257600080fd5b50515b905090565b73ffffffffffffffffffffffffffffffffffffffff811660009081526004602052604090205460ff16610f4b5773ffffffffffffffffffffffffffffffffffffffff811660008181526004602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055517f881febd4cd194dd4ace637642862aef1fb59a65c7e5551a5d9208f268d11c0069190a25056fea264697066735822beefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef64736f6c6343decafe0033"

// DeployFlags deploys a new Ethereum contract, binding an instance of Flags to it.
func DeployFlags(auth *bind.TransactOpts, backend bind.ContractBackend, racAddress common.Address) (common.Address, *types.Transaction, *Flags, error) {
	parsed, err := abi.JSON(strings.NewReader(FlagsABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(FlagsBin), backend, racAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Flags{FlagsCaller: FlagsCaller{contract: contract}, FlagsTransactor: FlagsTransactor{contract: contract}, FlagsFilterer: FlagsFilterer{contract: contract}}, nil
}

// Flags is an auto generated Go binding around an Ethereum contract.
type Flags struct {
	FlagsCaller     // Read-only binding to the contract
	FlagsTransactor // Write-only binding to the contract
	FlagsFilterer   // Log filterer for contract events
}

// FlagsCaller is an auto generated read-only Go binding around an Ethereum contract.
type FlagsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlagsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FlagsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlagsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FlagsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlagsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FlagsSession struct {
	Contract     *Flags            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FlagsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FlagsCallerSession struct {
	Contract *FlagsCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// FlagsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FlagsTransactorSession struct {
	Contract     *FlagsTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FlagsRaw is an auto generated low-level Go binding around an Ethereum contract.
type FlagsRaw struct {
	Contract *Flags // Generic contract binding to access the raw methods on
}

// FlagsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FlagsCallerRaw struct {
	Contract *FlagsCaller // Generic read-only contract binding to access the raw methods on
}

// FlagsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FlagsTransactorRaw struct {
	Contract *FlagsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFlags creates a new instance of Flags, bound to a specific deployed contract.
func NewFlags(address common.Address, backend bind.ContractBackend) (*Flags, error) {
	contract, err := bindFlags(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Flags{FlagsCaller: FlagsCaller{contract: contract}, FlagsTransactor: FlagsTransactor{contract: contract}, FlagsFilterer: FlagsFilterer{contract: contract}}, nil
}

// NewFlagsCaller creates a new read-only instance of Flags, bound to a specific deployed contract.
func NewFlagsCaller(address common.Address, caller bind.ContractCaller) (*FlagsCaller, error) {
	contract, err := bindFlags(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FlagsCaller{contract: contract}, nil
}

// NewFlagsTransactor creates a new write-only instance of Flags, bound to a specific deployed contract.
func NewFlagsTransactor(address common.Address, transactor bind.ContractTransactor) (*FlagsTransactor, error) {
	contract, err := bindFlags(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FlagsTransactor{contract: contract}, nil
}

// NewFlagsFilterer creates a new log filterer instance of Flags, bound to a specific deployed contract.
func NewFlagsFilterer(address common.Address, filterer bind.ContractFilterer) (*FlagsFilterer, error) {
	contract, err := bindFlags(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FlagsFilterer{contract: contract}, nil
}

// bindFlags binds a generic wrapper to an already deployed contract.
func bindFlags(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(FlagsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Flags *FlagsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Flags.Contract.FlagsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Flags *FlagsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Flags.Contract.FlagsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Flags *FlagsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Flags.Contract.FlagsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Flags *FlagsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Flags.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Flags *FlagsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Flags.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Flags *FlagsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Flags.Contract.contract.Transact(opts, method, params...)
}

// CheckEnabled is a free data retrieval call binding the contract method 0xdc7f0124.
//
// Solidity: function checkEnabled() view returns(bool)
func (_Flags *FlagsCaller) CheckEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Flags.contract.Call(opts, &out, "checkEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckEnabled is a free data retrieval call binding the contract method 0xdc7f0124.
//
// Solidity: function checkEnabled() view returns(bool)
func (_Flags *FlagsSession) CheckEnabled() (bool, error) {
	return _Flags.Contract.CheckEnabled(&_Flags.CallOpts)
}

// CheckEnabled is a free data retrieval call binding the contract method 0xdc7f0124.
//
// Solidity: function checkEnabled() view returns(bool)
func (_Flags *FlagsCallerSession) CheckEnabled() (bool, error) {
	return _Flags.Contract.CheckEnabled(&_Flags.CallOpts)
}

// GetFlag is a free data retrieval call binding the contract method 0x357e47fe.
//
// Solidity: function getFlag(address subject) view returns(bool)
func (_Flags *FlagsCaller) GetFlag(opts *bind.CallOpts, subject common.Address) (bool, error) {
	var out []interface{}
	err := _Flags.contract.Call(opts, &out, "getFlag", subject)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetFlag is a free data retrieval call binding the contract method 0x357e47fe.
//
// Solidity: function getFlag(address subject) view returns(bool)
func (_Flags *FlagsSession) GetFlag(subject common.Address) (bool, error) {
	return _Flags.Contract.GetFlag(&_Flags.CallOpts, subject)
}

// GetFlag is a free data retrieval call binding the contract method 0x357e47fe.
//
// Solidity: function getFlag(address subject) view returns(bool)
func (_Flags *FlagsCallerSession) GetFlag(subject common.Address) (bool, error) {
	return _Flags.Contract.GetFlag(&_Flags.CallOpts, subject)
}

// GetFlags is a free data retrieval call binding the contract method 0x7d723cac.
//
// Solidity: function getFlags(address[] subjects) view returns(bool[])
func (_Flags *FlagsCaller) GetFlags(opts *bind.CallOpts, subjects []common.Address) ([]bool, error) {
	var out []interface{}
	err := _Flags.contract.Call(opts, &out, "getFlags", subjects)

	if err != nil {
		return *new([]bool), err
	}

	out0 := *abi.ConvertType(out[0], new([]bool)).(*[]bool)

	return out0, err

}

// GetFlags is a free data retrieval call binding the contract method 0x7d723cac.
//
// Solidity: function getFlags(address[] subjects) view returns(bool[])
func (_Flags *FlagsSession) GetFlags(subjects []common.Address) ([]bool, error) {
	return _Flags.Contract.GetFlags(&_Flags.CallOpts, subjects)
}

// GetFlags is a free data retrieval call binding the contract method 0x7d723cac.
//
// Solidity: function getFlags(address[] subjects) view returns(bool[])
func (_Flags *FlagsCallerSession) GetFlags(subjects []common.Address) ([]bool, error) {
	return _Flags.Contract.GetFlags(&_Flags.CallOpts, subjects)
}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address _user, bytes _calldata) view returns(bool)
func (_Flags *FlagsCaller) HasAccess(opts *bind.CallOpts, _user common.Address, _calldata []byte) (bool, error) {
	var out []interface{}
	err := _Flags.contract.Call(opts, &out, "hasAccess", _user, _calldata)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address _user, bytes _calldata) view returns(bool)
func (_Flags *FlagsSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _Flags.Contract.HasAccess(&_Flags.CallOpts, _user, _calldata)
}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address _user, bytes _calldata) view returns(bool)
func (_Flags *FlagsCallerSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _Flags.Contract.HasAccess(&_Flags.CallOpts, _user, _calldata)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Flags *FlagsCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Flags.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Flags *FlagsSession) Owner() (common.Address, error) {
	return _Flags.Contract.Owner(&_Flags.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Flags *FlagsCallerSession) Owner() (common.Address, error) {
	return _Flags.Contract.Owner(&_Flags.CallOpts)
}

// RaisingAccessController is a free data retrieval call binding the contract method 0x2e1d859c.
//
// Solidity: function raisingAccessController() view returns(address)
func (_Flags *FlagsCaller) RaisingAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Flags.contract.Call(opts, &out, "raisingAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RaisingAccessController is a free data retrieval call binding the contract method 0x2e1d859c.
//
// Solidity: function raisingAccessController() view returns(address)
func (_Flags *FlagsSession) RaisingAccessController() (common.Address, error) {
	return _Flags.Contract.RaisingAccessController(&_Flags.CallOpts)
}

// RaisingAccessController is a free data retrieval call binding the contract method 0x2e1d859c.
//
// Solidity: function raisingAccessController() view returns(address)
func (_Flags *FlagsCallerSession) RaisingAccessController() (common.Address, error) {
	return _Flags.Contract.RaisingAccessController(&_Flags.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Flags *FlagsTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Flags *FlagsSession) AcceptOwnership() (*types.Transaction, error) {
	return _Flags.Contract.AcceptOwnership(&_Flags.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Flags *FlagsTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Flags.Contract.AcceptOwnership(&_Flags.TransactOpts)
}

// AddAccess is a paid mutator transaction binding the contract method 0xa118f249.
//
// Solidity: function addAccess(address _user) returns()
func (_Flags *FlagsTransactor) AddAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "addAccess", _user)
}

// AddAccess is a paid mutator transaction binding the contract method 0xa118f249.
//
// Solidity: function addAccess(address _user) returns()
func (_Flags *FlagsSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _Flags.Contract.AddAccess(&_Flags.TransactOpts, _user)
}

// AddAccess is a paid mutator transaction binding the contract method 0xa118f249.
//
// Solidity: function addAccess(address _user) returns()
func (_Flags *FlagsTransactorSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _Flags.Contract.AddAccess(&_Flags.TransactOpts, _user)
}

// DisableAccessCheck is a paid mutator transaction binding the contract method 0x0a756983.
//
// Solidity: function disableAccessCheck() returns()
func (_Flags *FlagsTransactor) DisableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "disableAccessCheck")
}

// DisableAccessCheck is a paid mutator transaction binding the contract method 0x0a756983.
//
// Solidity: function disableAccessCheck() returns()
func (_Flags *FlagsSession) DisableAccessCheck() (*types.Transaction, error) {
	return _Flags.Contract.DisableAccessCheck(&_Flags.TransactOpts)
}

// DisableAccessCheck is a paid mutator transaction binding the contract method 0x0a756983.
//
// Solidity: function disableAccessCheck() returns()
func (_Flags *FlagsTransactorSession) DisableAccessCheck() (*types.Transaction, error) {
	return _Flags.Contract.DisableAccessCheck(&_Flags.TransactOpts)
}

// EnableAccessCheck is a paid mutator transaction binding the contract method 0x8038e4a1.
//
// Solidity: function enableAccessCheck() returns()
func (_Flags *FlagsTransactor) EnableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "enableAccessCheck")
}

// EnableAccessCheck is a paid mutator transaction binding the contract method 0x8038e4a1.
//
// Solidity: function enableAccessCheck() returns()
func (_Flags *FlagsSession) EnableAccessCheck() (*types.Transaction, error) {
	return _Flags.Contract.EnableAccessCheck(&_Flags.TransactOpts)
}

// EnableAccessCheck is a paid mutator transaction binding the contract method 0x8038e4a1.
//
// Solidity: function enableAccessCheck() returns()
func (_Flags *FlagsTransactorSession) EnableAccessCheck() (*types.Transaction, error) {
	return _Flags.Contract.EnableAccessCheck(&_Flags.TransactOpts)
}

// LowerFlags is a paid mutator transaction binding the contract method 0x28286596.
//
// Solidity: function lowerFlags(address[] subjects) returns()
func (_Flags *FlagsTransactor) LowerFlags(opts *bind.TransactOpts, subjects []common.Address) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "lowerFlags", subjects)
}

// LowerFlags is a paid mutator transaction binding the contract method 0x28286596.
//
// Solidity: function lowerFlags(address[] subjects) returns()
func (_Flags *FlagsSession) LowerFlags(subjects []common.Address) (*types.Transaction, error) {
	return _Flags.Contract.LowerFlags(&_Flags.TransactOpts, subjects)
}

// LowerFlags is a paid mutator transaction binding the contract method 0x28286596.
//
// Solidity: function lowerFlags(address[] subjects) returns()
func (_Flags *FlagsTransactorSession) LowerFlags(subjects []common.Address) (*types.Transaction, error) {
	return _Flags.Contract.LowerFlags(&_Flags.TransactOpts, subjects)
}

// RaiseFlag is a paid mutator transaction binding the contract method 0xd74af263.
//
// Solidity: function raiseFlag(address subject) returns()
func (_Flags *FlagsTransactor) RaiseFlag(opts *bind.TransactOpts, subject common.Address) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "raiseFlag", subject)
}

// RaiseFlag is a paid mutator transaction binding the contract method 0xd74af263.
//
// Solidity: function raiseFlag(address subject) returns()
func (_Flags *FlagsSession) RaiseFlag(subject common.Address) (*types.Transaction, error) {
	return _Flags.Contract.RaiseFlag(&_Flags.TransactOpts, subject)
}

// RaiseFlag is a paid mutator transaction binding the contract method 0xd74af263.
//
// Solidity: function raiseFlag(address subject) returns()
func (_Flags *FlagsTransactorSession) RaiseFlag(subject common.Address) (*types.Transaction, error) {
	return _Flags.Contract.RaiseFlag(&_Flags.TransactOpts, subject)
}

// RaiseFlags is a paid mutator transaction binding the contract method 0x760bc82d.
//
// Solidity: function raiseFlags(address[] subjects) returns()
func (_Flags *FlagsTransactor) RaiseFlags(opts *bind.TransactOpts, subjects []common.Address) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "raiseFlags", subjects)
}

// RaiseFlags is a paid mutator transaction binding the contract method 0x760bc82d.
//
// Solidity: function raiseFlags(address[] subjects) returns()
func (_Flags *FlagsSession) RaiseFlags(subjects []common.Address) (*types.Transaction, error) {
	return _Flags.Contract.RaiseFlags(&_Flags.TransactOpts, subjects)
}

// RaiseFlags is a paid mutator transaction binding the contract method 0x760bc82d.
//
// Solidity: function raiseFlags(address[] subjects) returns()
func (_Flags *FlagsTransactorSession) RaiseFlags(subjects []common.Address) (*types.Transaction, error) {
	return _Flags.Contract.RaiseFlags(&_Flags.TransactOpts, subjects)
}

// RemoveAccess is a paid mutator transaction binding the contract method 0x8823da6c.
//
// Solidity: function removeAccess(address _user) returns()
func (_Flags *FlagsTransactor) RemoveAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "removeAccess", _user)
}

// RemoveAccess is a paid mutator transaction binding the contract method 0x8823da6c.
//
// Solidity: function removeAccess(address _user) returns()
func (_Flags *FlagsSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _Flags.Contract.RemoveAccess(&_Flags.TransactOpts, _user)
}

// RemoveAccess is a paid mutator transaction binding the contract method 0x8823da6c.
//
// Solidity: function removeAccess(address _user) returns()
func (_Flags *FlagsTransactorSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _Flags.Contract.RemoveAccess(&_Flags.TransactOpts, _user)
}

// SetRaisingAccessController is a paid mutator transaction binding the contract method 0x517e89fe.
//
// Solidity: function setRaisingAccessController(address racAddress) returns()
func (_Flags *FlagsTransactor) SetRaisingAccessController(opts *bind.TransactOpts, racAddress common.Address) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "setRaisingAccessController", racAddress)
}

// SetRaisingAccessController is a paid mutator transaction binding the contract method 0x517e89fe.
//
// Solidity: function setRaisingAccessController(address racAddress) returns()
func (_Flags *FlagsSession) SetRaisingAccessController(racAddress common.Address) (*types.Transaction, error) {
	return _Flags.Contract.SetRaisingAccessController(&_Flags.TransactOpts, racAddress)
}

// SetRaisingAccessController is a paid mutator transaction binding the contract method 0x517e89fe.
//
// Solidity: function setRaisingAccessController(address racAddress) returns()
func (_Flags *FlagsTransactorSession) SetRaisingAccessController(racAddress common.Address) (*types.Transaction, error) {
	return _Flags.Contract.SetRaisingAccessController(&_Flags.TransactOpts, racAddress)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_Flags *FlagsTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _Flags.contract.Transact(opts, "transferOwnership", _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_Flags *FlagsSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _Flags.Contract.TransferOwnership(&_Flags.TransactOpts, _to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _to) returns()
func (_Flags *FlagsTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _Flags.Contract.TransferOwnership(&_Flags.TransactOpts, _to)
}

// FlagsAddedAccessIterator is returned from FilterAddedAccess and is used to iterate over the raw logs and unpacked data for AddedAccess events raised by the Flags contract.
type FlagsAddedAccessIterator struct {
	Event *FlagsAddedAccess // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlagsAddedAccessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsAddedAccess)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlagsAddedAccess)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlagsAddedAccessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsAddedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsAddedAccess represents a AddedAccess event raised by the Flags contract.
type FlagsAddedAccess struct {
	User common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterAddedAccess is a free log retrieval operation binding the contract event 0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4.
//
// Solidity: event AddedAccess(address user)
func (_Flags *FlagsFilterer) FilterAddedAccess(opts *bind.FilterOpts) (*FlagsAddedAccessIterator, error) {

	logs, sub, err := _Flags.contract.FilterLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return &FlagsAddedAccessIterator{contract: _Flags.contract, event: "AddedAccess", logs: logs, sub: sub}, nil
}

// WatchAddedAccess is a free log subscription operation binding the contract event 0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4.
//
// Solidity: event AddedAccess(address user)
func (_Flags *FlagsFilterer) WatchAddedAccess(opts *bind.WatchOpts, sink chan<- *FlagsAddedAccess) (event.Subscription, error) {

	logs, sub, err := _Flags.contract.WatchLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsAddedAccess)
				if err := _Flags.contract.UnpackLog(event, "AddedAccess", log); err != nil {
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

// ParseAddedAccess is a log parse operation binding the contract event 0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4.
//
// Solidity: event AddedAccess(address user)
func (_Flags *FlagsFilterer) ParseAddedAccess(log types.Log) (*FlagsAddedAccess, error) {
	event := new(FlagsAddedAccess)
	if err := _Flags.contract.UnpackLog(event, "AddedAccess", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FlagsCheckAccessDisabledIterator is returned from FilterCheckAccessDisabled and is used to iterate over the raw logs and unpacked data for CheckAccessDisabled events raised by the Flags contract.
type FlagsCheckAccessDisabledIterator struct {
	Event *FlagsCheckAccessDisabled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlagsCheckAccessDisabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsCheckAccessDisabled)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlagsCheckAccessDisabled)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlagsCheckAccessDisabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsCheckAccessDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsCheckAccessDisabled represents a CheckAccessDisabled event raised by the Flags contract.
type FlagsCheckAccessDisabled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterCheckAccessDisabled is a free log retrieval operation binding the contract event 0x3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f539638.
//
// Solidity: event CheckAccessDisabled()
func (_Flags *FlagsFilterer) FilterCheckAccessDisabled(opts *bind.FilterOpts) (*FlagsCheckAccessDisabledIterator, error) {

	logs, sub, err := _Flags.contract.FilterLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return &FlagsCheckAccessDisabledIterator{contract: _Flags.contract, event: "CheckAccessDisabled", logs: logs, sub: sub}, nil
}

// WatchCheckAccessDisabled is a free log subscription operation binding the contract event 0x3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f539638.
//
// Solidity: event CheckAccessDisabled()
func (_Flags *FlagsFilterer) WatchCheckAccessDisabled(opts *bind.WatchOpts, sink chan<- *FlagsCheckAccessDisabled) (event.Subscription, error) {

	logs, sub, err := _Flags.contract.WatchLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsCheckAccessDisabled)
				if err := _Flags.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
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

// ParseCheckAccessDisabled is a log parse operation binding the contract event 0x3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f539638.
//
// Solidity: event CheckAccessDisabled()
func (_Flags *FlagsFilterer) ParseCheckAccessDisabled(log types.Log) (*FlagsCheckAccessDisabled, error) {
	event := new(FlagsCheckAccessDisabled)
	if err := _Flags.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FlagsCheckAccessEnabledIterator is returned from FilterCheckAccessEnabled and is used to iterate over the raw logs and unpacked data for CheckAccessEnabled events raised by the Flags contract.
type FlagsCheckAccessEnabledIterator struct {
	Event *FlagsCheckAccessEnabled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlagsCheckAccessEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsCheckAccessEnabled)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlagsCheckAccessEnabled)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlagsCheckAccessEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsCheckAccessEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsCheckAccessEnabled represents a CheckAccessEnabled event raised by the Flags contract.
type FlagsCheckAccessEnabled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterCheckAccessEnabled is a free log retrieval operation binding the contract event 0xaebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c3480.
//
// Solidity: event CheckAccessEnabled()
func (_Flags *FlagsFilterer) FilterCheckAccessEnabled(opts *bind.FilterOpts) (*FlagsCheckAccessEnabledIterator, error) {

	logs, sub, err := _Flags.contract.FilterLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return &FlagsCheckAccessEnabledIterator{contract: _Flags.contract, event: "CheckAccessEnabled", logs: logs, sub: sub}, nil
}

// WatchCheckAccessEnabled is a free log subscription operation binding the contract event 0xaebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c3480.
//
// Solidity: event CheckAccessEnabled()
func (_Flags *FlagsFilterer) WatchCheckAccessEnabled(opts *bind.WatchOpts, sink chan<- *FlagsCheckAccessEnabled) (event.Subscription, error) {

	logs, sub, err := _Flags.contract.WatchLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsCheckAccessEnabled)
				if err := _Flags.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
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

// ParseCheckAccessEnabled is a log parse operation binding the contract event 0xaebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c3480.
//
// Solidity: event CheckAccessEnabled()
func (_Flags *FlagsFilterer) ParseCheckAccessEnabled(log types.Log) (*FlagsCheckAccessEnabled, error) {
	event := new(FlagsCheckAccessEnabled)
	if err := _Flags.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FlagsFlagLoweredIterator is returned from FilterFlagLowered and is used to iterate over the raw logs and unpacked data for FlagLowered events raised by the Flags contract.
type FlagsFlagLoweredIterator struct {
	Event *FlagsFlagLowered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlagsFlagLoweredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsFlagLowered)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlagsFlagLowered)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlagsFlagLoweredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsFlagLoweredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsFlagLowered represents a FlagLowered event raised by the Flags contract.
type FlagsFlagLowered struct {
	Subject common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterFlagLowered is a free log retrieval operation binding the contract event 0xd86728e2e5cbaa28c1d357b5fbccc9c1ab0add09950eb7cac42df9acb24c4bc8.
//
// Solidity: event FlagLowered(address indexed subject)
func (_Flags *FlagsFilterer) FilterFlagLowered(opts *bind.FilterOpts, subject []common.Address) (*FlagsFlagLoweredIterator, error) {

	var subjectRule []interface{}
	for _, subjectItem := range subject {
		subjectRule = append(subjectRule, subjectItem)
	}

	logs, sub, err := _Flags.contract.FilterLogs(opts, "FlagLowered", subjectRule)
	if err != nil {
		return nil, err
	}
	return &FlagsFlagLoweredIterator{contract: _Flags.contract, event: "FlagLowered", logs: logs, sub: sub}, nil
}

// WatchFlagLowered is a free log subscription operation binding the contract event 0xd86728e2e5cbaa28c1d357b5fbccc9c1ab0add09950eb7cac42df9acb24c4bc8.
//
// Solidity: event FlagLowered(address indexed subject)
func (_Flags *FlagsFilterer) WatchFlagLowered(opts *bind.WatchOpts, sink chan<- *FlagsFlagLowered, subject []common.Address) (event.Subscription, error) {

	var subjectRule []interface{}
	for _, subjectItem := range subject {
		subjectRule = append(subjectRule, subjectItem)
	}

	logs, sub, err := _Flags.contract.WatchLogs(opts, "FlagLowered", subjectRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsFlagLowered)
				if err := _Flags.contract.UnpackLog(event, "FlagLowered", log); err != nil {
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

// ParseFlagLowered is a log parse operation binding the contract event 0xd86728e2e5cbaa28c1d357b5fbccc9c1ab0add09950eb7cac42df9acb24c4bc8.
//
// Solidity: event FlagLowered(address indexed subject)
func (_Flags *FlagsFilterer) ParseFlagLowered(log types.Log) (*FlagsFlagLowered, error) {
	event := new(FlagsFlagLowered)
	if err := _Flags.contract.UnpackLog(event, "FlagLowered", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FlagsFlagRaisedIterator is returned from FilterFlagRaised and is used to iterate over the raw logs and unpacked data for FlagRaised events raised by the Flags contract.
type FlagsFlagRaisedIterator struct {
	Event *FlagsFlagRaised // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlagsFlagRaisedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsFlagRaised)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlagsFlagRaised)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlagsFlagRaisedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsFlagRaisedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsFlagRaised represents a FlagRaised event raised by the Flags contract.
type FlagsFlagRaised struct {
	Subject common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterFlagRaised is a free log retrieval operation binding the contract event 0x881febd4cd194dd4ace637642862aef1fb59a65c7e5551a5d9208f268d11c006.
//
// Solidity: event FlagRaised(address indexed subject)
func (_Flags *FlagsFilterer) FilterFlagRaised(opts *bind.FilterOpts, subject []common.Address) (*FlagsFlagRaisedIterator, error) {

	var subjectRule []interface{}
	for _, subjectItem := range subject {
		subjectRule = append(subjectRule, subjectItem)
	}

	logs, sub, err := _Flags.contract.FilterLogs(opts, "FlagRaised", subjectRule)
	if err != nil {
		return nil, err
	}
	return &FlagsFlagRaisedIterator{contract: _Flags.contract, event: "FlagRaised", logs: logs, sub: sub}, nil
}

// WatchFlagRaised is a free log subscription operation binding the contract event 0x881febd4cd194dd4ace637642862aef1fb59a65c7e5551a5d9208f268d11c006.
//
// Solidity: event FlagRaised(address indexed subject)
func (_Flags *FlagsFilterer) WatchFlagRaised(opts *bind.WatchOpts, sink chan<- *FlagsFlagRaised, subject []common.Address) (event.Subscription, error) {

	var subjectRule []interface{}
	for _, subjectItem := range subject {
		subjectRule = append(subjectRule, subjectItem)
	}

	logs, sub, err := _Flags.contract.WatchLogs(opts, "FlagRaised", subjectRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsFlagRaised)
				if err := _Flags.contract.UnpackLog(event, "FlagRaised", log); err != nil {
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

// ParseFlagRaised is a log parse operation binding the contract event 0x881febd4cd194dd4ace637642862aef1fb59a65c7e5551a5d9208f268d11c006.
//
// Solidity: event FlagRaised(address indexed subject)
func (_Flags *FlagsFilterer) ParseFlagRaised(log types.Log) (*FlagsFlagRaised, error) {
	event := new(FlagsFlagRaised)
	if err := _Flags.contract.UnpackLog(event, "FlagRaised", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FlagsOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the Flags contract.
type FlagsOwnershipTransferRequestedIterator struct {
	Event *FlagsOwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlagsOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsOwnershipTransferRequested)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlagsOwnershipTransferRequested)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlagsOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the Flags contract.
type FlagsOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_Flags *FlagsFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FlagsOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Flags.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FlagsOwnershipTransferRequestedIterator{contract: _Flags.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_Flags *FlagsFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FlagsOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Flags.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsOwnershipTransferRequested)
				if err := _Flags.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_Flags *FlagsFilterer) ParseOwnershipTransferRequested(log types.Log) (*FlagsOwnershipTransferRequested, error) {
	event := new(FlagsOwnershipTransferRequested)
	if err := _Flags.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FlagsOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Flags contract.
type FlagsOwnershipTransferredIterator struct {
	Event *FlagsOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlagsOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsOwnershipTransferred)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlagsOwnershipTransferred)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlagsOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsOwnershipTransferred represents a OwnershipTransferred event raised by the Flags contract.
type FlagsOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_Flags *FlagsFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FlagsOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Flags.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FlagsOwnershipTransferredIterator{contract: _Flags.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_Flags *FlagsFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FlagsOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Flags.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsOwnershipTransferred)
				if err := _Flags.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_Flags *FlagsFilterer) ParseOwnershipTransferred(log types.Log) (*FlagsOwnershipTransferred, error) {
	event := new(FlagsOwnershipTransferred)
	if err := _Flags.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FlagsRaisingAccessControllerUpdatedIterator is returned from FilterRaisingAccessControllerUpdated and is used to iterate over the raw logs and unpacked data for RaisingAccessControllerUpdated events raised by the Flags contract.
type FlagsRaisingAccessControllerUpdatedIterator struct {
	Event *FlagsRaisingAccessControllerUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlagsRaisingAccessControllerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsRaisingAccessControllerUpdated)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlagsRaisingAccessControllerUpdated)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlagsRaisingAccessControllerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsRaisingAccessControllerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsRaisingAccessControllerUpdated represents a RaisingAccessControllerUpdated event raised by the Flags contract.
type FlagsRaisingAccessControllerUpdated struct {
	Previous common.Address
	Current  common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterRaisingAccessControllerUpdated is a free log retrieval operation binding the contract event 0xbaf9ea078655a4fffefd08f9435677bbc91e457a6d015fe7de1d0e68b8802cac.
//
// Solidity: event RaisingAccessControllerUpdated(address indexed previous, address indexed current)
func (_Flags *FlagsFilterer) FilterRaisingAccessControllerUpdated(opts *bind.FilterOpts, previous []common.Address, current []common.Address) (*FlagsRaisingAccessControllerUpdatedIterator, error) {

	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _Flags.contract.FilterLogs(opts, "RaisingAccessControllerUpdated", previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &FlagsRaisingAccessControllerUpdatedIterator{contract: _Flags.contract, event: "RaisingAccessControllerUpdated", logs: logs, sub: sub}, nil
}

// WatchRaisingAccessControllerUpdated is a free log subscription operation binding the contract event 0xbaf9ea078655a4fffefd08f9435677bbc91e457a6d015fe7de1d0e68b8802cac.
//
// Solidity: event RaisingAccessControllerUpdated(address indexed previous, address indexed current)
func (_Flags *FlagsFilterer) WatchRaisingAccessControllerUpdated(opts *bind.WatchOpts, sink chan<- *FlagsRaisingAccessControllerUpdated, previous []common.Address, current []common.Address) (event.Subscription, error) {

	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _Flags.contract.WatchLogs(opts, "RaisingAccessControllerUpdated", previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsRaisingAccessControllerUpdated)
				if err := _Flags.contract.UnpackLog(event, "RaisingAccessControllerUpdated", log); err != nil {
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

// ParseRaisingAccessControllerUpdated is a log parse operation binding the contract event 0xbaf9ea078655a4fffefd08f9435677bbc91e457a6d015fe7de1d0e68b8802cac.
//
// Solidity: event RaisingAccessControllerUpdated(address indexed previous, address indexed current)
func (_Flags *FlagsFilterer) ParseRaisingAccessControllerUpdated(log types.Log) (*FlagsRaisingAccessControllerUpdated, error) {
	event := new(FlagsRaisingAccessControllerUpdated)
	if err := _Flags.contract.UnpackLog(event, "RaisingAccessControllerUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// FlagsRemovedAccessIterator is returned from FilterRemovedAccess and is used to iterate over the raw logs and unpacked data for RemovedAccess events raised by the Flags contract.
type FlagsRemovedAccessIterator struct {
	Event *FlagsRemovedAccess // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlagsRemovedAccessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlagsRemovedAccess)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlagsRemovedAccess)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlagsRemovedAccessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlagsRemovedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlagsRemovedAccess represents a RemovedAccess event raised by the Flags contract.
type FlagsRemovedAccess struct {
	User common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterRemovedAccess is a free log retrieval operation binding the contract event 0x3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1.
//
// Solidity: event RemovedAccess(address user)
func (_Flags *FlagsFilterer) FilterRemovedAccess(opts *bind.FilterOpts) (*FlagsRemovedAccessIterator, error) {

	logs, sub, err := _Flags.contract.FilterLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return &FlagsRemovedAccessIterator{contract: _Flags.contract, event: "RemovedAccess", logs: logs, sub: sub}, nil
}

// WatchRemovedAccess is a free log subscription operation binding the contract event 0x3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1.
//
// Solidity: event RemovedAccess(address user)
func (_Flags *FlagsFilterer) WatchRemovedAccess(opts *bind.WatchOpts, sink chan<- *FlagsRemovedAccess) (event.Subscription, error) {

	logs, sub, err := _Flags.contract.WatchLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlagsRemovedAccess)
				if err := _Flags.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
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

// ParseRemovedAccess is a log parse operation binding the contract event 0x3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1.
//
// Solidity: event RemovedAccess(address user)
func (_Flags *FlagsFilterer) ParseRemovedAccess(log types.Log) (*FlagsRemovedAccess, error) {
	event := new(FlagsRemovedAccess)
	if err := _Flags.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
		return nil, err
	}
	return event, nil
}
