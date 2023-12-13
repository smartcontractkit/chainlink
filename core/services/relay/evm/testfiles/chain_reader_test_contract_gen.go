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
	Account        common.Address
	Accounts       []common.Address
	BigField       *big.Int
	NestedStruct   MidLevelTestStruct
}

// TestfilesMetaData contains all meta data concerning the Testfiles contract.
var TestfilesMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"Account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"Accounts\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"indexed\":false,\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"Triggered\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"AddTestStruct\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetDifferentPrimitiveValue\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"i\",\"type\":\"uint256\"}],\"name\":\"GetElementAtIndex\",\"outputs\":[{\"components\":[{\"internalType\":\"int32\",\"name\":\"Field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"DifferentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"OracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"OracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"Account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"Accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"BigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"NestedStruct\",\"type\":\"tuple\"}],\"internalType\":\"structTestStruct\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetPrimitiveValue\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetSliceValue\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"ReturnSeen\",\"outputs\":[{\"components\":[{\"internalType\":\"int32\",\"name\":\"Field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"DifferentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"OracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"OracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"Account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"Accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"BigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"NestedStruct\",\"type\":\"tuple\"}],\"internalType\":\"structTestStruct\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"TriggerEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50600180548082018255600082905260048082047fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6908101805460086003958616810261010090810a8088026001600160401b0391820219909416939093179093558654808801909755848704909301805496909516909202900a91820291021990921691909117905561125b806100a96000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c80639ca04f671161005b5780639ca04f67146100cc578063b9dad6b0146100ec578063bdb37c90146100ff578063da8e7a821461011457600080fd5b8063030d3ca2146100825780636c7cf955146100a45780638b659d6e146100b9575b600080fd5b6107c65b60405167ffffffffffffffff90911681526020015b60405180910390f35b6100b76100b2366004610b6c565b61011b565b005b6100b76100c7366004610b6c565b610371565b6100df6100da366004610c5e565b6103c6565b60405161009b9190610f97565b6100df6100fa366004610b6c565b610685565b610107610781565b60405161009b9190610e8c565b6003610086565b60006040518061010001604052808c60030b81526020018b8b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060ff8a166020808301919091526040805161040081810183529190930192918b918390839080828437600092019190915250505081526001600160a01b038816602080830191909152604080518883028181018401835289825291909301929189918991829190850190849080828437600092019190915250505090825250601785900b6020820152604001610200846110d7565b905281546001818101845560009384526020938490208351600a90930201805460039390930b63ffffffff1663ffffffff19909316929092178255838301518051939492936102579392850192919091019061080d565b50604082015160028201805460ff191660ff90921691909117905560608201516102879060038301906020610891565b5060808201516004820180546001600160a01b0319166001600160a01b0390921691909117905560a082015180516102c991600584019160209091019061091f565b5060c082015160068201805460179290920b6001600160c01b03166001600160c01b031990921691909117905560e082015180516007808401805460f09390931c61ffff1990931692909217825560208084015180516008870180549190940b67ffffffffffffffff1667ffffffffffffffff1990911617835580820151805191939261035e9260098901929091019061080d565b5050505050505050505050505050505050565b7f7188419dcd8b51877b71766f075f3626586c0ff190e7d056aa65ce9acb649a3d8a8a8a8a8a8a8a8a8a8a6040516103b29a99989796959493929190610eda565b60405180910390a150505050505050505050565b6103ce610974565b60006103db6001846110b2565b815481106103eb576103eb6111f9565b90600052602060002090600a0201604051806101000160405290816000820160009054906101000a900460030b60030b60030b8152602001600182018054610432906111c4565b80601f016020809104026020016040519081016040528092919081815260200182805461045e906111c4565b80156104ab5780601f10610480576101008083540402835291602001916104ab565b820191906000526020600020905b81548152906001019060200180831161048e57829003601f168201915b5050509183525050600282015460ff166020808301919091526040805161040081018083529190930192916003850191826000855b825461010083900a900460ff168152602060019283018181049485019490930390920291018084116104e05750505092845250505060048201546001600160a01b0316602080830191909152600583018054604080518285028101850182528281529401939283018282801561057f57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610561575b50505091835250506006820154601790810b810b900b6020808301919091526040805180820182526007808601805460f01b6001600160f01b031916835283518085018552600888018054840b840b90930b815260098801805495909701969395919486830194919392840191906105f6906111c4565b80601f0160208091040260200160405190810160405280929190818152602001828054610622906111c4565b801561066f5780601f106106445761010080835404028352916020019161066f565b820191906000526020600020905b81548152906001019060200180831161065257829003601f168201915b5050509190925250505090525090525092915050565b61068d610974565b6040518061010001604052808c60030b81526020018b8b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060ff8a166020808301919091526040805161040081810183529190930192918b918390839080828437600092019190915250505081526001600160a01b038816602080830191909152604080518883028181018401835289825291909301929189918991829190850190849080828437600092019190915250505090825250601785900b6020820152604001610770846110d7565b90529b9a5050505050505050505050565b6060600180548060200260200160405190810160405280929190818152602001828054801561080357602002820191906000526020600020906000905b82829054906101000a900467ffffffffffffffff1667ffffffffffffffff16815260200190600801906020826007010492830192600103820291508084116107be5790505b5050505050905090565b828054610819906111c4565b90600052602060002090601f01602090048101928261083b5760008555610881565b82601f1061085457805160ff1916838001178555610881565b82800160010185558215610881579182015b82811115610881578251825591602001919060010190610866565b5061088d9291506109c3565b5090565b6001830191839082156108815791602002820160005b838211156108e557835183826101000a81548160ff021916908360ff16021790555092602001926001016020816000010492830192600103026108a7565b80156109125782816101000a81549060ff02191690556001016020816000010492830192600103026108e5565b505061088d9291506109c3565b828054828255906000526020600020908101928215610881579160200282015b8281111561088157825182546001600160a01b0319166001600160a01b0390911617825560209092019160019091019061093f565b604080516101008101825260008082526060602083018190529282015290810161099c6109d8565b81526000602082018190526060604083018190528201526080016109be6109f7565b905290565b5b8082111561088d57600081556001016109c4565b6040518061040001604052806020906020820280368337509192915050565b604051806040016040528060006001600160f01b03191681526020016109be6040518060400160405280600060070b8152602001606081525090565b80356001600160a01b0381168114610a4a57600080fd5b919050565b60008083601f840112610a6157600080fd5b50813567ffffffffffffffff811115610a7957600080fd5b6020830191508360208260051b8501011115610a9457600080fd5b9250929050565b806104008101831015610aad57600080fd5b92915050565b80356001600160f01b031981168114610a4a57600080fd5b8035601781900b8114610a4a57600080fd5b8035600381900b8114610a4a57600080fd5b8035600781900b8114610a4a57600080fd5b60008083601f840112610b1357600080fd5b50813567ffffffffffffffff811115610b2b57600080fd5b602083019150836020828501011115610a9457600080fd5b600060408284031215610b5557600080fd5b50919050565b803560ff81168114610a4a57600080fd5b6000806000806000806000806000806104e08b8d031215610b8c57600080fd5b610b958b610add565b995060208b013567ffffffffffffffff80821115610bb257600080fd5b610bbe8e838f01610b01565b909b509950899150610bd260408e01610b5b565b9850610be18e60608f01610a9b565b9750610bf06104608e01610a33565b96506104808d0135915080821115610c0757600080fd5b610c138e838f01610a4f565b9096509450849150610c286104a08e01610acb565b93506104c08d0135915080821115610c3f57600080fd5b50610c4c8d828e01610b43565b9150509295989b9194979a5092959850565b600060208284031215610c7057600080fd5b5035919050565b8183526000602080850194508260005b85811015610cb3576001600160a01b03610ca083610a33565b1687529582019590820190600101610c87565b509495945050505050565b600081518084526020808501945080840160005b83811015610cb35781516001600160a01b031687529582019590820190600101610cd2565b8060005b6020808210610d0a5750610d21565b825160ff1685529384019390910190600101610cfb565b50505050565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b6000815180845260005b81811015610d7657602081850181015186830182015201610d5a565b81811115610d88576000602083870101525b50601f01601f19169290920160200192915050565b6001600160f01b0319610daf82610ab3565b16825260006020820135603e19833603018112610dcb57600080fd5b604060208501528201610ddd81610aef565b60070b60408501526020810135601e19823603018112610dfc57600080fd5b8101803567ffffffffffffffff811115610e1557600080fd5b803603831315610e2457600080fd5b60406060870152610e3c608087018260208501610d27565b9695505050505050565b61ffff60f01b81511682526000602082015160406020850152805160070b60408501526020810151905060406060850152610e846080850182610d50565b949350505050565b6020808252825182820181905260009190848201906040850190845b81811015610ece57835167ffffffffffffffff1683529284019291840191600101610ea8565b50909695505050505050565b60006104e08c60030b835260208181850152610ef98285018d8f610d27565b915060ff808c166040860152606085018b60005b84811015610f325783610f1f83610b5b565b1683529184019190840190600101610f0d565b5050505050610f4d6104608401896001600160a01b03169052565b828103610480840152610f61818789610c77565b9050610f736104a084018660170b9052565b8281036104c0840152610f868185610d9d565b9d9c50505050505050505050505050565b60208152610fab60208201835160030b9052565b600060208301516104e0806040850152610fc9610500850183610d50565b91506040850151610fdf606086018260ff169052565b506060850151610ff26080860182610cf7565b5060808501516001600160a01b031661048085015260a0850151601f1985840381016104a08701526110248483610cbe565b935060c0870151915061103d6104c087018360170b9052565b60e0870151915080868503018387015250610e3c8382610e46565b6040805190810167ffffffffffffffff8111828210171561107b5761107b61120f565b60405290565b604051601f8201601f1916810167ffffffffffffffff811182821017156110aa576110aa61120f565b604052919050565b6000828210156110d257634e487b7160e01b600052601160045260246000fd5b500390565b6000604082360312156110e957600080fd5b6110f1611058565b6110fa83610ab3565b815260208084013567ffffffffffffffff8082111561111857600080fd5b81860191506040823603121561112d57600080fd5b611135611058565b61113e83610aef565b8152838301358281111561115157600080fd5b929092019136601f84011261116557600080fd5b8235828111156111775761117761120f565b611189601f8201601f19168601611081565b9250808352368582860101111561119f57600080fd5b8085850186850137600090830185015280840191909152918301919091525092915050565b600181811c908216806111d857607f821691505b60208210811415610b5557634e487b7160e01b600052602260045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052604160045260246000fdfea2646970667358221220dce1950ddad4c207dfaef7c48111c80f3a73cf85e5711bf7a3afe15473519a1664736f6c63430008060033",
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

// GetDifferentPrimitiveValue is a free data retrieval call binding the contract method 0x030d3ca2.
//
// Solidity: function GetDifferentPrimitiveValue() pure returns(uint64)
func (_Testfiles *TestfilesCaller) GetDifferentPrimitiveValue(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Testfiles.contract.Call(opts, &out, "GetDifferentPrimitiveValue")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetDifferentPrimitiveValue is a free data retrieval call binding the contract method 0x030d3ca2.
//
// Solidity: function GetDifferentPrimitiveValue() pure returns(uint64)
func (_Testfiles *TestfilesSession) GetDifferentPrimitiveValue() (uint64, error) {
	return _Testfiles.Contract.GetDifferentPrimitiveValue(&_Testfiles.CallOpts)
}

// GetDifferentPrimitiveValue is a free data retrieval call binding the contract method 0x030d3ca2.
//
// Solidity: function GetDifferentPrimitiveValue() pure returns(uint64)
func (_Testfiles *TestfilesCallerSession) GetDifferentPrimitiveValue() (uint64, error) {
	return _Testfiles.Contract.GetDifferentPrimitiveValue(&_Testfiles.CallOpts)
}

// GetElementAtIndex is a free data retrieval call binding the contract method 0x9ca04f67.
//
// Solidity: function GetElementAtIndex(uint256 i) view returns((int32,string,uint8,uint8[32],address,address[],int192,(bytes2,(int64,string))))
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
// Solidity: function GetElementAtIndex(uint256 i) view returns((int32,string,uint8,uint8[32],address,address[],int192,(bytes2,(int64,string))))
func (_Testfiles *TestfilesSession) GetElementAtIndex(i *big.Int) (TestStruct, error) {
	return _Testfiles.Contract.GetElementAtIndex(&_Testfiles.CallOpts, i)
}

// GetElementAtIndex is a free data retrieval call binding the contract method 0x9ca04f67.
//
// Solidity: function GetElementAtIndex(uint256 i) view returns((int32,string,uint8,uint8[32],address,address[],int192,(bytes2,(int64,string))))
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

// ReturnSeen is a free data retrieval call binding the contract method 0xb9dad6b0.
//
// Solidity: function ReturnSeen(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, address account, address[] accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct) pure returns((int32,string,uint8,uint8[32],address,address[],int192,(bytes2,(int64,string))))
func (_Testfiles *TestfilesCaller) ReturnSeen(opts *bind.CallOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (TestStruct, error) {
	var out []interface{}
	err := _Testfiles.contract.Call(opts, &out, "ReturnSeen", field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)

	if err != nil {
		return *new(TestStruct), err
	}

	out0 := *abi.ConvertType(out[0], new(TestStruct)).(*TestStruct)

	return out0, err

}

// ReturnSeen is a free data retrieval call binding the contract method 0xb9dad6b0.
//
// Solidity: function ReturnSeen(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, address account, address[] accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct) pure returns((int32,string,uint8,uint8[32],address,address[],int192,(bytes2,(int64,string))))
func (_Testfiles *TestfilesSession) ReturnSeen(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (TestStruct, error) {
	return _Testfiles.Contract.ReturnSeen(&_Testfiles.CallOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

// ReturnSeen is a free data retrieval call binding the contract method 0xb9dad6b0.
//
// Solidity: function ReturnSeen(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, address account, address[] accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct) pure returns((int32,string,uint8,uint8[32],address,address[],int192,(bytes2,(int64,string))))
func (_Testfiles *TestfilesCallerSession) ReturnSeen(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (TestStruct, error) {
	return _Testfiles.Contract.ReturnSeen(&_Testfiles.CallOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

// AddTestStruct is a paid mutator transaction binding the contract method 0x6c7cf955.
//
// Solidity: function AddTestStruct(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, address account, address[] accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct) returns()
func (_Testfiles *TestfilesTransactor) AddTestStruct(opts *bind.TransactOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _Testfiles.contract.Transact(opts, "AddTestStruct", field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

// AddTestStruct is a paid mutator transaction binding the contract method 0x6c7cf955.
//
// Solidity: function AddTestStruct(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, address account, address[] accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct) returns()
func (_Testfiles *TestfilesSession) AddTestStruct(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _Testfiles.Contract.AddTestStruct(&_Testfiles.TransactOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

// AddTestStruct is a paid mutator transaction binding the contract method 0x6c7cf955.
//
// Solidity: function AddTestStruct(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, address account, address[] accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct) returns()
func (_Testfiles *TestfilesTransactorSession) AddTestStruct(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _Testfiles.Contract.AddTestStruct(&_Testfiles.TransactOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

// TriggerEvent is a paid mutator transaction binding the contract method 0x8b659d6e.
//
// Solidity: function TriggerEvent(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, address account, address[] accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct) returns()
func (_Testfiles *TestfilesTransactor) TriggerEvent(opts *bind.TransactOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _Testfiles.contract.Transact(opts, "TriggerEvent", field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

// TriggerEvent is a paid mutator transaction binding the contract method 0x8b659d6e.
//
// Solidity: function TriggerEvent(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, address account, address[] accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct) returns()
func (_Testfiles *TestfilesSession) TriggerEvent(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _Testfiles.Contract.TriggerEvent(&_Testfiles.TransactOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

// TriggerEvent is a paid mutator transaction binding the contract method 0x8b659d6e.
//
// Solidity: function TriggerEvent(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, address account, address[] accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct) returns()
func (_Testfiles *TestfilesTransactorSession) TriggerEvent(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account common.Address, accounts []common.Address, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _Testfiles.Contract.TriggerEvent(&_Testfiles.TransactOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

// TestfilesTriggeredIterator is returned from FilterTriggered and is used to iterate over the raw logs and unpacked data for Triggered events raised by the Testfiles contract.
type TestfilesTriggeredIterator struct {
	Event *TestfilesTriggered // Event containing the contract specifics and raw log

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
func (it *TestfilesTriggeredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestfilesTriggered)
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
		it.Event = new(TestfilesTriggered)
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
func (it *TestfilesTriggeredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestfilesTriggeredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestfilesTriggered represents a Triggered event raised by the Testfiles contract.
type TestfilesTriggered struct {
	Field          int32
	DifferentField string
	OracleId       uint8
	OracleIds      [32]uint8
	Account        common.Address
	Accounts       []common.Address
	BigField       *big.Int
	NestedStruct   MidLevelTestStruct
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterTriggered is a free log retrieval operation binding the contract event 0x7188419dcd8b51877b71766f075f3626586c0ff190e7d056aa65ce9acb649a3d.
//
// Solidity: event Triggered(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, address Account, address[] Accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct)
func (_Testfiles *TestfilesFilterer) FilterTriggered(opts *bind.FilterOpts) (*TestfilesTriggeredIterator, error) {

	logs, sub, err := _Testfiles.contract.FilterLogs(opts, "Triggered")
	if err != nil {
		return nil, err
	}
	return &TestfilesTriggeredIterator{contract: _Testfiles.contract, event: "Triggered", logs: logs, sub: sub}, nil
}

// WatchTriggered is a free log subscription operation binding the contract event 0x7188419dcd8b51877b71766f075f3626586c0ff190e7d056aa65ce9acb649a3d.
//
// Solidity: event Triggered(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, address Account, address[] Accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct)
func (_Testfiles *TestfilesFilterer) WatchTriggered(opts *bind.WatchOpts, sink chan<- *TestfilesTriggered) (event.Subscription, error) {

	logs, sub, err := _Testfiles.contract.WatchLogs(opts, "Triggered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestfilesTriggered)
				if err := _Testfiles.contract.UnpackLog(event, "Triggered", log); err != nil {
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

// ParseTriggered is a log parse operation binding the contract event 0x7188419dcd8b51877b71766f075f3626586c0ff190e7d056aa65ce9acb649a3d.
//
// Solidity: event Triggered(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, address Account, address[] Accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct)
func (_Testfiles *TestfilesFilterer) ParseTriggered(log types.Log) (*TestfilesTriggered, error) {
	event := new(TestfilesTriggered)
	if err := _Testfiles.contract.UnpackLog(event, "Triggered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
