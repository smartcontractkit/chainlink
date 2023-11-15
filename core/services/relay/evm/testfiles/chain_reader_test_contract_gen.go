// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package testfiles

import (
	"errors"
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

// InnerTestStruct is an auto generated low-level Go binding around an user-defined struct.
type InnerTestStruct struct {
	I int64
	S string
}

// MidLevelTestStruct is an auto generated low-level Go binding around an user-defined struct.
type MidLevelTestStruct struct {
	FixedBytes [2]byte
	Inner      InnerTestStruct
}

// TestStruct is an auto generated low-level Go binding around an user-defined struct.
type TestStruct struct {
	Field          int32
	DifferentField string
	OracleId       uint8
	OracleIds      [32]uint8
	Account        [32]byte
	Accounts       [][32]byte
	BigField       *big.Int
	NestedStruct   MidLevelTestStruct
}

// TestfilesMetaData contains all meta data concerning the Testfiles contract.
var TestfilesMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"bytes32\",\"name\":\"account\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"accounts\",\"type\":\"bytes32[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"AddTestStruct\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"i\",\"type\":\"uint256\"}],\"name\":\"GetElementAtIndex\",\"outputs\":[{\"components\":[{\"internalType\":\"int32\",\"name\":\"Field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"DifferentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"OracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"OracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"bytes32\",\"name\":\"Account\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"Accounts\",\"type\":\"bytes32[]\"},{\"internalType\":\"int192\",\"name\":\"BigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"NestedStruct\",\"type\":\"tuple\"}],\"internalType\":\"structTestStruct\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetPrimitiveValue\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetSliceValue\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50600180548082018255600082905260048082047fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6908101805460086003958616810261010090810a8088026001600160401b0391820219909416939093179093558654808801909755848704909301805496909516909202900a918202910219909216919091179055610e4d806100a96000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80637dd6af5b146100515780639ca04f6714610066578063bdb37c901461008f578063da8e7a82146100a4575b600080fd5b61006461005f366004610921565b6100b3565b005b610079610074366004610a0c565b6102e2565b6040516100869190610b71565b60405180910390f35b610097610592565b6040516100869190610b23565b60405160038152602001610086565b60006040518061010001604052808c60030b81526020018b8b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060ff8a166020808301919091526040805161040081810183529190930192918b918390839080828437600092019190915250505081526020808201899052604080518883028181018401835289825291909301929189918991829190850190849080828437600092019190915250505090825250601785900b602082015260400161018d84610cb2565b905281546001818101845560009384526020938490208351600a90930201805460039390930b63ffffffff1663ffffffff19909316929092178255838301518051939492936101e49392850192919091019061061e565b50604082015160028201805460ff191660ff909216919091179055606082015161021490600383019060206106a2565b506080820151600482015560a0820151805161023a916005840191602090910190610730565b5060c082015160068201805460179290920b6001600160c01b03166001600160c01b031990921691909117905560e082015180516007808401805460f09390931c61ffff1990931692909217825560208084015180516008870180549190940b67ffffffffffffffff1667ffffffffffffffff199091161783558082015180519193926102cf9260098901929091019061061e565b5050505050505050505050505050505050565b6102ea61076a565b60006102f7600184610c8d565b8154811061030757610307610deb565b90600052602060002090600a0201604051806101000160405290816000820160009054906101000a900460030b60030b60030b815260200160018201805461034e90610db6565b80601f016020809104026020016040519081016040528092919081815260200182805461037a90610db6565b80156103c75780601f1061039c576101008083540402835291602001916103c7565b820191906000526020600020905b8154815290600101906020018083116103aa57829003601f168201915b5050509183525050600282015460ff166020808301919091526040805161040081018083529190930192916003850191826000855b825461010083900a900460ff168152602060019283018181049485019490930390920291018084116103fc57905050505050508152602001600482015481526020016005820180548060200260200160405190810160405280929190818152602001828054801561048c57602002820191906000526020600020905b815481526020019060010190808311610478575b50505091835250506006820154601790810b810b900b6020808301919091526040805180820182526007808601805460f01b6001600160f01b031916835283518085018552600888018054840b840b90930b8152600988018054959097019693959194868301949193928401919061050390610db6565b80601f016020809104026020016040519081016040528092919081815260200182805461052f90610db6565b801561057c5780601f106105515761010080835404028352916020019161057c565b820191906000526020600020905b81548152906001019060200180831161055f57829003601f168201915b5050509190925250505090525090525092915050565b6060600180548060200260200160405190810160405280929190818152602001828054801561061457602002820191906000526020600020906000905b82829054906101000a900467ffffffffffffffff1667ffffffffffffffff16815260200190600801906020826007010492830192600103820291508084116105cf5790505b5050505050905090565b82805461062a90610db6565b90600052602060002090601f01602090048101928261064c5760008555610692565b82601f1061066557805160ff1916838001178555610692565b82800160010185558215610692579182015b82811115610692578251825591602001919060010190610677565b5061069e9291506107b9565b5090565b6001830191839082156106925791602002820160005b838211156106f657835183826101000a81548160ff021916908360ff16021790555092602001926001016020816000010492830192600103026106b8565b80156107235782816101000a81549060ff02191690556001016020816000010492830192600103026106f6565b505061069e9291506107b9565b8280548282559060005260206000209081019282156106925791602002820182811115610692578251825591602001919060010190610677565b60408051610100810182526000808252606060208301819052928201529081016107926107ce565b81526000602082018190526060604083018190528201526080016107b46107ed565b905290565b5b8082111561069e57600081556001016107ba565b6040518061040001604052806020906020820280368337509192915050565b604051806040016040528060006001600160f01b03191681526020016107b46040518060400160405280600060070b8152602001606081525090565b60008083601f84011261083b57600080fd5b50813567ffffffffffffffff81111561085357600080fd5b6020830191508360208260051b850101111561086e57600080fd5b9250929050565b80610400810183101561088757600080fd5b92915050565b8035601781900b811461089f57600080fd5b919050565b8035600381900b811461089f57600080fd5b60008083601f8401126108c857600080fd5b50813567ffffffffffffffff8111156108e057600080fd5b60208301915083602082850101111561086e57600080fd5b60006040828403121561090a57600080fd5b50919050565b803560ff8116811461089f57600080fd5b6000806000806000806000806000806104e08b8d03121561094157600080fd5b61094a8b6108a4565b995060208b013567ffffffffffffffff8082111561096757600080fd5b6109738e838f016108b6565b909b50995089915061098760408e01610910565b98506109968e60608f01610875565b97506104608d013596506104808d01359150808211156109b557600080fd5b6109c18e838f01610829565b90965094508491506109d66104a08e0161088d565b93506104c08d01359150808211156109ed57600080fd5b506109fa8d828e016108f8565b9150509295989b9194979a5092959850565b600060208284031215610a1e57600080fd5b5035919050565b600081518084526020808501945080840160005b83811015610a5557815187529582019590820190600101610a39565b509495945050505050565b8060005b6020808210610a735750610a8a565b825160ff1685529384019390910190600101610a64565b50505050565b6000815180845260005b81811015610ab657602081850181015186830182015201610a9a565b81811115610ac8576000602083870101525b50601f01601f19169290920160200192915050565b61ffff60f01b81511682526000602082015160406020850152805160070b60408501526020810151905060406060850152610b1b6080850182610a90565b949350505050565b6020808252825182820181905260009190848201906040850190845b81811015610b6557835167ffffffffffffffff1683529284019291840191600101610b3f565b50909695505050505050565b60208152610b8560208201835160030b9052565b600060208301516104e0806040850152610ba3610500850183610a90565b91506040850151610bb9606086018260ff169052565b506060850151610bcc6080860182610a60565b50608085015161048085015260a0850151601f1980868503016104a0870152610bf58483610a25565b935060c08701519150610c0e6104c087018360170b9052565b60e0870151915080868503018387015250610c298382610add565b9695505050505050565b6040805190810167ffffffffffffffff81118282101715610c5657610c56610e01565b60405290565b604051601f8201601f1916810167ffffffffffffffff81118282101715610c8557610c85610e01565b604052919050565b600082821015610cad57634e487b7160e01b600052601160045260246000fd5b500390565b600060408236031215610cc457600080fd5b610ccc610c33565b82356001600160f01b031981168114610ce457600080fd5b815260208381013567ffffffffffffffff80821115610d0257600080fd5b818601915060408236031215610d1757600080fd5b610d1f610c33565b82358060070b8114610d3057600080fd5b81528284013582811115610d4357600080fd5b929092019136601f840112610d5757600080fd5b823582811115610d6957610d69610e01565b610d7b601f8201601f19168601610c5c565b92508083523685828601011115610d9157600080fd5b8085850186850137600090830185015280840191909152918301919091525092915050565b600181811c90821680610dca57607f821691505b6020821081141561090a57634e487b7160e01b600052602260045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052604160045260246000fdfea26469706673582212206d3fcd0d6af66016f39be3d9948d2703697b0dfdbab6f47b494ab3d1f3d0d4c264736f6c63430008060033",
}

// TestfilesABI is the input ABI used to generate the binding from.
// Deprecated: Use TestfilesMetaData.ABI instead.
var TestfilesABI = TestfilesMetaData.ABI

// TestfilesBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestfilesMetaData.Bin instead.
var TestfilesBin = TestfilesMetaData.Bin

// DeployTestfiles deploys a new Ethereum contract, binding an instance of Testfiles to it.
func DeployTestfiles(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Testfiles, error) {
	parsed, err := TestfilesMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestfilesBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Testfiles{TestfilesCaller: TestfilesCaller{contract: contract}, TestfilesTransactor: TestfilesTransactor{contract: contract}, TestfilesFilterer: TestfilesFilterer{contract: contract}}, nil
}

// Testfiles is an auto generated Go binding around an Ethereum contract.
type Testfiles struct {
	TestfilesCaller     // Read-only binding to the contract
	TestfilesTransactor // Write-only binding to the contract
	TestfilesFilterer   // Log filterer for contract events
}

// TestfilesCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestfilesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestfilesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestfilesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestfilesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestfilesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestfilesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestfilesSession struct {
	Contract     *Testfiles        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestfilesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestfilesCallerSession struct {
	Contract *TestfilesCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// TestfilesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestfilesTransactorSession struct {
	Contract     *TestfilesTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// TestfilesRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestfilesRaw struct {
	Contract *Testfiles // Generic contract binding to access the raw methods on
}

// TestfilesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestfilesCallerRaw struct {
	Contract *TestfilesCaller // Generic read-only contract binding to access the raw methods on
}

// TestfilesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestfilesTransactorRaw struct {
	Contract *TestfilesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestfiles creates a new instance of Testfiles, bound to a specific deployed contract.
func NewTestfiles(address common.Address, backend bind.ContractBackend) (*Testfiles, error) {
	contract, err := bindTestfiles(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Testfiles{TestfilesCaller: TestfilesCaller{contract: contract}, TestfilesTransactor: TestfilesTransactor{contract: contract}, TestfilesFilterer: TestfilesFilterer{contract: contract}}, nil
}

// NewTestfilesCaller creates a new read-only instance of Testfiles, bound to a specific deployed contract.
func NewTestfilesCaller(address common.Address, caller bind.ContractCaller) (*TestfilesCaller, error) {
	contract, err := bindTestfiles(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestfilesCaller{contract: contract}, nil
}

// NewTestfilesTransactor creates a new write-only instance of Testfiles, bound to a specific deployed contract.
func NewTestfilesTransactor(address common.Address, transactor bind.ContractTransactor) (*TestfilesTransactor, error) {
	contract, err := bindTestfiles(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestfilesTransactor{contract: contract}, nil
}

// NewTestfilesFilterer creates a new log filterer instance of Testfiles, bound to a specific deployed contract.
func NewTestfilesFilterer(address common.Address, filterer bind.ContractFilterer) (*TestfilesFilterer, error) {
	contract, err := bindTestfiles(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestfilesFilterer{contract: contract}, nil
}

// bindTestfiles binds a generic wrapper to an already deployed contract.
func bindTestfiles(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestfilesMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Testfiles *TestfilesRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Testfiles.Contract.TestfilesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Testfiles *TestfilesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Testfiles.Contract.TestfilesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Testfiles *TestfilesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Testfiles.Contract.TestfilesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Testfiles *TestfilesCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Testfiles.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Testfiles *TestfilesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Testfiles.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Testfiles *TestfilesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Testfiles.Contract.contract.Transact(opts, method, params...)
}

// GetElementAtIndex is a free data retrieval call binding the contract method 0x9ca04f67.
//
// Solidity: function GetElementAtIndex(uint256 i) view returns((int32,string,uint8,uint8[32],bytes32,bytes32[],int192,(bytes2,(int64,string))))
func (_Testfiles *TestfilesCaller) GetElementAtIndex(opts *bind.CallOpts, i *big.Int) (TestStruct, error) {
	var out []interface{}
	err := _Testfiles.contract.Call(opts, &out, "GetElementAtIndex", i)

	if err != nil {
		return *new(TestStruct), err
	}

	out0 := *abi.ConvertType(out[0], new(TestStruct)).(*TestStruct)

	return out0, err

}

// GetElementAtIndex is a free data retrieval call binding the contract method 0x9ca04f67.
//
// Solidity: function GetElementAtIndex(uint256 i) view returns((int32,string,uint8,uint8[32],bytes32,bytes32[],int192,(bytes2,(int64,string))))
func (_Testfiles *TestfilesSession) GetElementAtIndex(i *big.Int) (TestStruct, error) {
	return _Testfiles.Contract.GetElementAtIndex(&_Testfiles.CallOpts, i)
}

// GetElementAtIndex is a free data retrieval call binding the contract method 0x9ca04f67.
//
// Solidity: function GetElementAtIndex(uint256 i) view returns((int32,string,uint8,uint8[32],bytes32,bytes32[],int192,(bytes2,(int64,string))))
func (_Testfiles *TestfilesCallerSession) GetElementAtIndex(i *big.Int) (TestStruct, error) {
	return _Testfiles.Contract.GetElementAtIndex(&_Testfiles.CallOpts, i)
}

// GetPrimitiveValue is a free data retrieval call binding the contract method 0xda8e7a82.
//
// Solidity: function GetPrimitiveValue() pure returns(uint64)
func (_Testfiles *TestfilesCaller) GetPrimitiveValue(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Testfiles.contract.Call(opts, &out, "GetPrimitiveValue")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetPrimitiveValue is a free data retrieval call binding the contract method 0xda8e7a82.
//
// Solidity: function GetPrimitiveValue() pure returns(uint64)
func (_Testfiles *TestfilesSession) GetPrimitiveValue() (uint64, error) {
	return _Testfiles.Contract.GetPrimitiveValue(&_Testfiles.CallOpts)
}

// GetPrimitiveValue is a free data retrieval call binding the contract method 0xda8e7a82.
//
// Solidity: function GetPrimitiveValue() pure returns(uint64)
func (_Testfiles *TestfilesCallerSession) GetPrimitiveValue() (uint64, error) {
	return _Testfiles.Contract.GetPrimitiveValue(&_Testfiles.CallOpts)
}

// GetSliceValue is a free data retrieval call binding the contract method 0xbdb37c90.
//
// Solidity: function GetSliceValue() view returns(uint64[])
func (_Testfiles *TestfilesCaller) GetSliceValue(opts *bind.CallOpts) ([]uint64, error) {
	var out []interface{}
	err := _Testfiles.contract.Call(opts, &out, "GetSliceValue")

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

// GetSliceValue is a free data retrieval call binding the contract method 0xbdb37c90.
//
// Solidity: function GetSliceValue() view returns(uint64[])
func (_Testfiles *TestfilesSession) GetSliceValue() ([]uint64, error) {
	return _Testfiles.Contract.GetSliceValue(&_Testfiles.CallOpts)
}

// GetSliceValue is a free data retrieval call binding the contract method 0xbdb37c90.
//
// Solidity: function GetSliceValue() view returns(uint64[])
func (_Testfiles *TestfilesCallerSession) GetSliceValue() ([]uint64, error) {
	return _Testfiles.Contract.GetSliceValue(&_Testfiles.CallOpts)
}

// AddTestStruct is a paid mutator transaction binding the contract method 0x7dd6af5b.
//
// Solidity: function AddTestStruct(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, bytes32 account, bytes32[] accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct) returns()
func (_Testfiles *TestfilesTransactor) AddTestStruct(opts *bind.TransactOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account [32]byte, accounts [][32]byte, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _Testfiles.contract.Transact(opts, "AddTestStruct", field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

// AddTestStruct is a paid mutator transaction binding the contract method 0x7dd6af5b.
//
// Solidity: function AddTestStruct(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, bytes32 account, bytes32[] accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct) returns()
func (_Testfiles *TestfilesSession) AddTestStruct(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account [32]byte, accounts [][32]byte, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _Testfiles.Contract.AddTestStruct(&_Testfiles.TransactOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

// AddTestStruct is a paid mutator transaction binding the contract method 0x7dd6af5b.
//
// Solidity: function AddTestStruct(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, bytes32 account, bytes32[] accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct) returns()
func (_Testfiles *TestfilesTransactorSession) AddTestStruct(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account [32]byte, accounts [][32]byte, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _Testfiles.Contract.AddTestStruct(&_Testfiles.TransactOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}
