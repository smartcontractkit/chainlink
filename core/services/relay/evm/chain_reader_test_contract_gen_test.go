// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package evm_test

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

// EvmTestMetaData contains all meta data concerning the EvmTest contract.
var EvmTestMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"bytes32\",\"name\":\"account\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"accounts\",\"type\":\"bytes32[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"AddTestStruct\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"i\",\"type\":\"uint256\"}],\"name\":\"GetElementAtIndex\",\"outputs\":[{\"components\":[{\"internalType\":\"int32\",\"name\":\"Field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"DifferentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"OracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"OracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"bytes32\",\"name\":\"Account\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"Accounts\",\"type\":\"bytes32[]\"},{\"internalType\":\"int192\",\"name\":\"BigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"NestedStruct\",\"type\":\"tuple\"}],\"internalType\":\"structTestStruct\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610d0a806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80637dd6af5b1461003b5780639ca04f6714610050575b600080fd5b61004e610049366004610851565b610079565b005b61006361005e36600461093c565b6102a8565b6040516100709190610a53565b60405180910390f35b60006040518061010001604052808c60030b81526020018b8b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060ff8a166020808301919091526040805161040081810183529190930192918b918390839080828437600092019190915250505081526020808201899052604080518883028181018401835289825291909301929189918991829190850190849080828437600092019190915250505090825250601785900b602082015260400161015384610b6f565b905281546001818101845560009384526020938490208351600a90930201805460039390930b63ffffffff1663ffffffff19909316929092178255838301518051939492936101aa9392850192919091019061054e565b50604082015160028201805460ff191660ff90921691909117905560608201516101da90600383019060206105d2565b506080820151600482015560a08201518051610200916005840191602090910190610660565b5060c082015160068201805460179290920b6001600160c01b03166001600160c01b031990921691909117905560e082015180516007808401805460f09390931c61ffff1990931692909217825560208084015180516008870180549190940b67ffffffffffffffff1667ffffffffffffffff199091161783558082015180519193926102959260098901929091019061054e565b5050505050505050505050505050505050565b6102b061069a565b600082815481106102c3576102c3610ca8565b90600052602060002090600a0201604051806101000160405290816000820160009054906101000a900460030b60030b60030b815260200160018201805461030a90610c73565b80601f016020809104026020016040519081016040528092919081815260200182805461033690610c73565b80156103835780601f1061035857610100808354040283529160200191610383565b820191906000526020600020905b81548152906001019060200180831161036657829003601f168201915b5050509183525050600282015460ff166020808301919091526040805161040081018083529190930192916003850191826000855b825461010083900a900460ff168152602060019283018181049485019490930390920291018084116103b857905050505050508152602001600482015481526020016005820180548060200260200160405190810160405280929190818152602001828054801561044857602002820191906000526020600020905b815481526020019060010190808311610434575b50505091835250506006820154601790810b810b900b6020808301919091526040805180820182526007808601805460f01b6001600160f01b031916835283518085018552600888018054840b840b90930b815260098801805495909701969395919486830194919392840191906104bf90610c73565b80601f01602080910402602001604051908101604052809291908181526020018280546104eb90610c73565b80156105385780601f1061050d57610100808354040283529160200191610538565b820191906000526020600020905b81548152906001019060200180831161051b57829003601f168201915b5050509190925250505090525090525092915050565b82805461055a90610c73565b90600052602060002090601f01602090048101928261057c57600085556105c2565b82601f1061059557805160ff19168380011785556105c2565b828001600101855582156105c2579182015b828111156105c25782518255916020019190600101906105a7565b506105ce9291506106e9565b5090565b6001830191839082156105c25791602002820160005b8382111561062657835183826101000a81548160ff021916908360ff16021790555092602001926001016020816000010492830192600103026105e8565b80156106535782816101000a81549060ff0219169055600101602081600001049283019260010302610626565b50506105ce9291506106e9565b8280548282559060005260206000209081019282156105c257916020028201828111156105c25782518255916020019190600101906105a7565b60408051610100810182526000808252606060208301819052928201529081016106c26106fe565b81526000602082018190526060604083018190528201526080016106e461071d565b905290565b5b808211156105ce57600081556001016106ea565b6040518061040001604052806020906020820280368337509192915050565b604051806040016040528060006001600160f01b03191681526020016106e46040518060400160405280600060070b8152602001606081525090565b60008083601f84011261076b57600080fd5b50813567ffffffffffffffff81111561078357600080fd5b6020830191508360208260051b850101111561079e57600080fd5b9250929050565b8061040081018310156107b757600080fd5b92915050565b8035601781900b81146107cf57600080fd5b919050565b8035600381900b81146107cf57600080fd5b60008083601f8401126107f857600080fd5b50813567ffffffffffffffff81111561081057600080fd5b60208301915083602082850101111561079e57600080fd5b60006040828403121561083a57600080fd5b50919050565b803560ff811681146107cf57600080fd5b6000806000806000806000806000806104e08b8d03121561087157600080fd5b61087a8b6107d4565b995060208b013567ffffffffffffffff8082111561089757600080fd5b6108a38e838f016107e6565b909b5099508991506108b760408e01610840565b98506108c68e60608f016107a5565b97506104608d013596506104808d01359150808211156108e557600080fd5b6108f18e838f01610759565b90965094508491506109066104a08e016107bd565b93506104c08d013591508082111561091d57600080fd5b5061092a8d828e01610828565b9150509295989b9194979a5092959850565b60006020828403121561094e57600080fd5b5035919050565b600081518084526020808501945080840160005b8381101561098557815187529582019590820190600101610969565b509495945050505050565b8060005b60208082106109a357506109ba565b825160ff1685529384019390910190600101610994565b50505050565b6000815180845260005b818110156109e6576020818501810151868301820152016109ca565b818111156109f8576000602083870101525b50601f01601f19169290920160200192915050565b61ffff60f01b81511682526000602082015160406020850152805160070b60408501526020810151905060406060850152610a4b60808501826109c0565b949350505050565b60208152610a6760208201835160030b9052565b600060208301516104e0806040850152610a856105008501836109c0565b91506040850151610a9b606086018260ff169052565b506060850151610aae6080860182610990565b50608085015161048085015260a0850151601f1980868503016104a0870152610ad78483610955565b935060c08701519150610af06104c087018360170b9052565b60e0870151915080868503018387015250610b0b8382610a0d565b9695505050505050565b6040805190810167ffffffffffffffff81118282101715610b3857610b38610cbe565b60405290565b604051601f8201601f1916810167ffffffffffffffff81118282101715610b6757610b67610cbe565b604052919050565b600060408236031215610b8157600080fd5b610b89610b15565b82356001600160f01b031981168114610ba157600080fd5b815260208381013567ffffffffffffffff80821115610bbf57600080fd5b818601915060408236031215610bd457600080fd5b610bdc610b15565b82358060070b8114610bed57600080fd5b81528284013582811115610c0057600080fd5b929092019136601f840112610c1457600080fd5b823582811115610c2657610c26610cbe565b610c38601f8201601f19168601610b3e565b92508083523685828601011115610c4e57600080fd5b8085850186850137600090830185015280840191909152918301919091525092915050565b600181811c90821680610c8757607f821691505b6020821081141561083a57634e487b7160e01b600052602260045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052604160045260246000fdfea2646970667358221220198b9d5324ad4bde80005411dbaa95c987f2056c252dbf4f8948dc03653dbac864736f6c63430008060033",
}

// EvmTestABI is the input ABI used to generate the binding from.
// Deprecated: Use EvmTestMetaData.ABI instead.
var EvmTestABI = EvmTestMetaData.ABI

// EvmTestBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use EvmTestMetaData.Bin instead.
var EvmTestBin = EvmTestMetaData.Bin

// DeployEvmTest deploys a new Ethereum contract, binding an instance of EvmTest to it.
func DeployEvmTest(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *EvmTest, error) {
	parsed, err := EvmTestMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EvmTestBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EvmTest{EvmTestCaller: EvmTestCaller{contract: contract}, EvmTestTransactor: EvmTestTransactor{contract: contract}, EvmTestFilterer: EvmTestFilterer{contract: contract}}, nil
}

// EvmTest is an auto generated Go binding around an Ethereum contract.
type EvmTest struct {
	EvmTestCaller     // Read-only binding to the contract
	EvmTestTransactor // Write-only binding to the contract
	EvmTestFilterer   // Log filterer for contract events
}

// EvmTestCaller is an auto generated read-only Go binding around an Ethereum contract.
type EvmTestCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EvmTestTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EvmTestTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EvmTestFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EvmTestFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EvmTestSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EvmTestSession struct {
	Contract     *EvmTest          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EvmTestCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EvmTestCallerSession struct {
	Contract *EvmTestCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// EvmTestTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EvmTestTransactorSession struct {
	Contract     *EvmTestTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// EvmTestRaw is an auto generated low-level Go binding around an Ethereum contract.
type EvmTestRaw struct {
	Contract *EvmTest // Generic contract binding to access the raw methods on
}

// EvmTestCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EvmTestCallerRaw struct {
	Contract *EvmTestCaller // Generic read-only contract binding to access the raw methods on
}

// EvmTestTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EvmTestTransactorRaw struct {
	Contract *EvmTestTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEvmTest creates a new instance of EvmTest, bound to a specific deployed contract.
func NewEvmTest(address common.Address, backend bind.ContractBackend) (*EvmTest, error) {
	contract, err := bindEvmTest(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EvmTest{EvmTestCaller: EvmTestCaller{contract: contract}, EvmTestTransactor: EvmTestTransactor{contract: contract}, EvmTestFilterer: EvmTestFilterer{contract: contract}}, nil
}

// NewEvmTestCaller creates a new read-only instance of EvmTest, bound to a specific deployed contract.
func NewEvmTestCaller(address common.Address, caller bind.ContractCaller) (*EvmTestCaller, error) {
	contract, err := bindEvmTest(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EvmTestCaller{contract: contract}, nil
}

// NewEvmTestTransactor creates a new write-only instance of EvmTest, bound to a specific deployed contract.
func NewEvmTestTransactor(address common.Address, transactor bind.ContractTransactor) (*EvmTestTransactor, error) {
	contract, err := bindEvmTest(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EvmTestTransactor{contract: contract}, nil
}

// NewEvmTestFilterer creates a new log filterer instance of EvmTest, bound to a specific deployed contract.
func NewEvmTestFilterer(address common.Address, filterer bind.ContractFilterer) (*EvmTestFilterer, error) {
	contract, err := bindEvmTest(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EvmTestFilterer{contract: contract}, nil
}

// bindEvmTest binds a generic wrapper to an already deployed contract.
func bindEvmTest(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EvmTestMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EvmTest *EvmTestRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EvmTest.Contract.EvmTestCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EvmTest *EvmTestRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EvmTest.Contract.EvmTestTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EvmTest *EvmTestRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EvmTest.Contract.EvmTestTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EvmTest *EvmTestCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EvmTest.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EvmTest *EvmTestTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EvmTest.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EvmTest *EvmTestTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EvmTest.Contract.contract.Transact(opts, method, params...)
}

// GetElementAtIndex is a free data retrieval call binding the contract method 0x9ca04f67.
//
// Solidity: function GetElementAtIndex(uint256 i) view returns((int32,string,uint8,uint8[32],bytes32,bytes32[],int192,(bytes2,(int64,string))))
func (_EvmTest *EvmTestCaller) GetElementAtIndex(opts *bind.CallOpts, i *big.Int) (TestStruct, error) {
	var out []interface{}
	err := _EvmTest.contract.Call(opts, &out, "GetElementAtIndex", i)

	if err != nil {
		return *new(TestStruct), err
	}

	out0 := *abi.ConvertType(out[0], new(TestStruct)).(*TestStruct)

	return out0, err

}

// GetElementAtIndex is a free data retrieval call binding the contract method 0x9ca04f67.
//
// Solidity: function GetElementAtIndex(uint256 i) view returns((int32,string,uint8,uint8[32],bytes32,bytes32[],int192,(bytes2,(int64,string))))
func (_EvmTest *EvmTestSession) GetElementAtIndex(i *big.Int) (TestStruct, error) {
	return _EvmTest.Contract.GetElementAtIndex(&_EvmTest.CallOpts, i)
}

// GetElementAtIndex is a free data retrieval call binding the contract method 0x9ca04f67.
//
// Solidity: function GetElementAtIndex(uint256 i) view returns((int32,string,uint8,uint8[32],bytes32,bytes32[],int192,(bytes2,(int64,string))))
func (_EvmTest *EvmTestCallerSession) GetElementAtIndex(i *big.Int) (TestStruct, error) {
	return _EvmTest.Contract.GetElementAtIndex(&_EvmTest.CallOpts, i)
}

// AddTestStruct is a paid mutator transaction binding the contract method 0x7dd6af5b.
//
// Solidity: function AddTestStruct(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, bytes32 account, bytes32[] accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct) returns()
func (_EvmTest *EvmTestTransactor) AddTestStruct(opts *bind.TransactOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account [32]byte, accounts [][32]byte, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _EvmTest.contract.Transact(opts, "AddTestStruct", field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

// AddTestStruct is a paid mutator transaction binding the contract method 0x7dd6af5b.
//
// Solidity: function AddTestStruct(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, bytes32 account, bytes32[] accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct) returns()
func (_EvmTest *EvmTestSession) AddTestStruct(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account [32]byte, accounts [][32]byte, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _EvmTest.Contract.AddTestStruct(&_EvmTest.TransactOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

// AddTestStruct is a paid mutator transaction binding the contract method 0x7dd6af5b.
//
// Solidity: function AddTestStruct(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, bytes32 account, bytes32[] accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct) returns()
func (_EvmTest *EvmTestTransactorSession) AddTestStruct(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account [32]byte, accounts [][32]byte, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _EvmTest.Contract.AddTestStruct(&_EvmTest.TransactOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}
