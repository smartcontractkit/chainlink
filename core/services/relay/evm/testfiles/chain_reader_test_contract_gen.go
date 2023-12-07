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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"account\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"accounts\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"indexed\":false,\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"Triggered\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"bytes32\",\"name\":\"account\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"accounts\",\"type\":\"bytes32[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"AddTestStruct\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"i\",\"type\":\"uint256\"}],\"name\":\"GetElementAtIndex\",\"outputs\":[{\"components\":[{\"internalType\":\"int32\",\"name\":\"Field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"DifferentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"OracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"OracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"bytes32\",\"name\":\"Account\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"Accounts\",\"type\":\"bytes32[]\"},{\"internalType\":\"int192\",\"name\":\"BigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"NestedStruct\",\"type\":\"tuple\"}],\"internalType\":\"structTestStruct\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetPrimitiveValue\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetSliceValue\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"bytes32\",\"name\":\"account\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"accounts\",\"type\":\"bytes32[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"ReturnSeen\",\"outputs\":[{\"components\":[{\"internalType\":\"int32\",\"name\":\"Field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"DifferentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"OracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"OracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"bytes32\",\"name\":\"Account\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"Accounts\",\"type\":\"bytes32[]\"},{\"internalType\":\"int192\",\"name\":\"BigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"NestedStruct\",\"type\":\"tuple\"}],\"internalType\":\"structTestStruct\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"bytes32\",\"name\":\"account\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"accounts\",\"type\":\"bytes32[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"TriggerEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50600180548082018255600082905260048082047fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6908101805460086003958616810261010090810a8088026001600160401b0391820219909416939093179093558654808801909755848704909301805496909516909202900a91820291021990921691909117905561118e806100a96000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c80637dd6af5b146100675780639ca04f671461007c578063b95ad411146100a5578063bdb37c90146100b8578063da8e7a82146100cd578063e669831d146100dc575b600080fd5b61007a610075366004610acd565b6100ef565b005b61008f61008a366004610bb8565b61031e565b60405161009c9190610ed3565b60405180910390f35b61008f6100b3366004610acd565b6105ce565b6100c06106bf565b60405161009c9190610dd7565b6040516003815260200161009c565b61007a6100ea366004610acd565b61074b565b60006040518061010001604052808c60030b81526020018b8b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060ff8a166020808301919091526040805161040081810183529190930192918b918390839080828437600092019190915250505081526020808201899052604080518883028181018401835289825291909301929189918991829190850190849080828437600092019190915250505090825250601785900b60208201526040016101c98461100a565b905281546001818101845560009384526020938490208351600a90930201805460039390930b63ffffffff1663ffffffff1990931692909217825583830151805193949293610220939285019291909101906107a0565b50604082015160028201805460ff191660ff90921691909117905560608201516102509060038301906020610824565b506080820151600482015560a082015180516102769160058401916020909101906108b2565b5060c082015160068201805460179290920b6001600160c01b03166001600160c01b031990921691909117905560e082015180516007808401805460f09390931c61ffff1990931692909217825560208084015180516008870180549190940b67ffffffffffffffff1667ffffffffffffffff1990911617835580820151805191939261030b926009890192909101906107a0565b5050505050505050505050505050505050565b6103266108ec565b6000610333600184610fe5565b815481106103435761034361112c565b90600052602060002090600a0201604051806101000160405290816000820160009054906101000a900460030b60030b60030b815260200160018201805461038a906110f7565b80601f01602080910402602001604051908101604052809291908181526020018280546103b6906110f7565b80156104035780601f106103d857610100808354040283529160200191610403565b820191906000526020600020905b8154815290600101906020018083116103e657829003601f168201915b5050509183525050600282015460ff166020808301919091526040805161040081018083529190930192916003850191826000855b825461010083900a900460ff168152602060019283018181049485019490930390920291018084116104385790505050505050815260200160048201548152602001600582018054806020026020016040519081016040528092919081815260200182805480156104c857602002820191906000526020600020905b8154815260200190600101908083116104b4575b50505091835250506006820154601790810b810b900b6020808301919091526040805180820182526007808601805460f01b6001600160f01b031916835283518085018552600888018054840b840b90930b8152600988018054959097019693959194868301949193928401919061053f906110f7565b80601f016020809104026020016040519081016040528092919081815260200182805461056b906110f7565b80156105b85780601f1061058d576101008083540402835291602001916105b8565b820191906000526020600020905b81548152906001019060200180831161059b57829003601f168201915b5050509190925250505090525090525092915050565b6105d66108ec565b6040518061010001604052808c60030b81526020018b8b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060ff8a166020808301919091526040805161040081810183529190930192918b918390839080828437600092019190915250505081526020808201899052604080518883028181018401835289825291909301929189918991829190850190849080828437600092019190915250505090825250601785900b60208201526040016106ae8461100a565b90529b9a5050505050505050505050565b6060600180548060200260200160405190810160405280929190818152602001828054801561074157602002820191906000526020600020906000905b82829054906101000a900467ffffffffffffffff1667ffffffffffffffff16815260200190600801906020826007010492830192600103820291508084116106fc5790505b5050505050905090565b7f7d2abe6109e46b893ac1835c9287d6ad5c5ccf3d0254d2ca72225873795e0f538a8a8a8a8a8a8a8a8a8a60405161078c9a99989796959493929190610e25565b60405180910390a150505050505050505050565b8280546107ac906110f7565b90600052602060002090601f0160209004810192826107ce5760008555610814565b82601f106107e757805160ff1916838001178555610814565b82800160010185558215610814579182015b828111156108145782518255916020019190600101906107f9565b5061082092915061093b565b5090565b6001830191839082156108145791602002820160005b8382111561087857835183826101000a81548160ff021916908360ff160217905550926020019260010160208160000104928301926001030261083a565b80156108a55782816101000a81549060ff0219169055600101602081600001049283019260010302610878565b505061082092915061093b565b82805482825590600052602060002090810192821561081457916020028201828111156108145782518255916020019190600101906107f9565b6040805161010081018252600080825260606020830181905292820152908101610914610950565b815260006020820181905260606040830181905282015260800161093661096f565b905290565b5b80821115610820576000815560010161093c565b6040518061040001604052806020906020820280368337509192915050565b604051806040016040528060006001600160f01b03191681526020016109366040518060400160405280600060070b8152602001606081525090565b60008083601f8401126109bd57600080fd5b50813567ffffffffffffffff8111156109d557600080fd5b6020830191508360208260051b85010111156109f057600080fd5b9250929050565b806104008101831015610a0957600080fd5b92915050565b80356001600160f01b031981168114610a2757600080fd5b919050565b8035601781900b8114610a2757600080fd5b8035600381900b8114610a2757600080fd5b8035600781900b8114610a2757600080fd5b60008083601f840112610a7457600080fd5b50813567ffffffffffffffff811115610a8c57600080fd5b6020830191508360208285010111156109f057600080fd5b600060408284031215610ab657600080fd5b50919050565b803560ff81168114610a2757600080fd5b6000806000806000806000806000806104e08b8d031215610aed57600080fd5b610af68b610a3e565b995060208b013567ffffffffffffffff80821115610b1357600080fd5b610b1f8e838f01610a62565b909b509950899150610b3360408e01610abc565b9850610b428e60608f016109f7565b97506104608d013596506104808d0135915080821115610b6157600080fd5b610b6d8e838f016109ab565b9096509450849150610b826104a08e01610a2c565b93506104c08d0135915080821115610b9957600080fd5b50610ba68d828e01610aa4565b9150509295989b9194979a5092959850565b600060208284031215610bca57600080fd5b5035919050565b81835260006001600160fb1b03831115610bea57600080fd5b8260051b8083602087013760009401602001938452509192915050565b600081518084526020808501945080840160005b83811015610c3757815187529582019590820190600101610c1b565b509495945050505050565b8060005b6020808210610c555750610c6c565b825160ff1685529384019390910190600101610c46565b50505050565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b6000815180845260005b81811015610cc157602081850181015186830182015201610ca5565b81811115610cd3576000602083870101525b50601f01601f19169290920160200192915050565b6001600160f01b0319610cfa82610a0f565b16825260006020820135603e19833603018112610d1657600080fd5b604060208501528201610d2881610a50565b60070b60408501526020810135601e19823603018112610d4757600080fd5b8101803567ffffffffffffffff811115610d6057600080fd5b803603831315610d6f57600080fd5b60406060870152610d87608087018260208501610c72565b9695505050505050565b61ffff60f01b81511682526000602082015160406020850152805160070b60408501526020810151905060406060850152610dcf6080850182610c9b565b949350505050565b6020808252825182820181905260009190848201906040850190845b81811015610e1957835167ffffffffffffffff1683529284019291840191600101610df3565b50909695505050505050565b60006104e08c60030b835260208181850152610e448285018d8f610c72565b915060ff808c166040860152606085018b60005b84811015610e7d5783610e6a83610abc565b1683529184019190840190600101610e58565b505050505087610460840152828103610480840152610e9d818789610bd1565b9050610eaf6104a084018660170b9052565b8281036104c0840152610ec28185610ce8565b9d9c50505050505050505050505050565b60208152610ee760208201835160030b9052565b600060208301516104e0806040850152610f05610500850183610c9b565b91506040850151610f1b606086018260ff169052565b506060850151610f2e6080860182610c42565b50608085015161048085015260a0850151601f1980868503016104a0870152610f578483610c07565b935060c08701519150610f706104c087018360170b9052565b60e0870151915080868503018387015250610d878382610d91565b6040805190810167ffffffffffffffff81118282101715610fae57610fae611142565b60405290565b604051601f8201601f1916810167ffffffffffffffff81118282101715610fdd57610fdd611142565b604052919050565b60008282101561100557634e487b7160e01b600052601160045260246000fd5b500390565b60006040823603121561101c57600080fd5b611024610f8b565b61102d83610a0f565b815260208084013567ffffffffffffffff8082111561104b57600080fd5b81860191506040823603121561106057600080fd5b611068610f8b565b61107183610a50565b8152838301358281111561108457600080fd5b929092019136601f84011261109857600080fd5b8235828111156110aa576110aa611142565b6110bc601f8201601f19168601610fb4565b925080835236858286010111156110d257600080fd5b8085850186850137600090830185015280840191909152918301919091525092915050565b600181811c9082168061110b57607f821691505b60208210811415610ab657634e487b7160e01b600052602260045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052604160045260246000fdfea2646970667358221220dc7c8830d1691fb07419d1e88774f8489193b9b32a6a2406f298335883e9593b64736f6c63430008060033",
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

// ReturnSeen is a free data retrieval call binding the contract method 0xb95ad411.
//
// Solidity: function ReturnSeen(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, bytes32 account, bytes32[] accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct) pure returns((int32,string,uint8,uint8[32],bytes32,bytes32[],int192,(bytes2,(int64,string))))
func (_Testfiles *TestfilesCaller) ReturnSeen(opts *bind.CallOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account [32]byte, accounts [][32]byte, bigField *big.Int, nestedStruct MidLevelTestStruct) (TestStruct, error) {
	var out []interface{}
	err := _Testfiles.contract.Call(opts, &out, "ReturnSeen", field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)

	if err != nil {
		return *new(TestStruct), err
	}

	out0 := *abi.ConvertType(out[0], new(TestStruct)).(*TestStruct)

	return out0, err

}

// ReturnSeen is a free data retrieval call binding the contract method 0xb95ad411.
//
// Solidity: function ReturnSeen(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, bytes32 account, bytes32[] accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct) pure returns((int32,string,uint8,uint8[32],bytes32,bytes32[],int192,(bytes2,(int64,string))))
func (_Testfiles *TestfilesSession) ReturnSeen(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account [32]byte, accounts [][32]byte, bigField *big.Int, nestedStruct MidLevelTestStruct) (TestStruct, error) {
	return _Testfiles.Contract.ReturnSeen(&_Testfiles.CallOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

// ReturnSeen is a free data retrieval call binding the contract method 0xb95ad411.
//
// Solidity: function ReturnSeen(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, bytes32 account, bytes32[] accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct) pure returns((int32,string,uint8,uint8[32],bytes32,bytes32[],int192,(bytes2,(int64,string))))
func (_Testfiles *TestfilesCallerSession) ReturnSeen(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account [32]byte, accounts [][32]byte, bigField *big.Int, nestedStruct MidLevelTestStruct) (TestStruct, error) {
	return _Testfiles.Contract.ReturnSeen(&_Testfiles.CallOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
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

// TriggerEvent is a paid mutator transaction binding the contract method 0xe669831d.
//
// Solidity: function TriggerEvent(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, bytes32 account, bytes32[] accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct) returns()
func (_Testfiles *TestfilesTransactor) TriggerEvent(opts *bind.TransactOpts, field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account [32]byte, accounts [][32]byte, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _Testfiles.contract.Transact(opts, "TriggerEvent", field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

// TriggerEvent is a paid mutator transaction binding the contract method 0xe669831d.
//
// Solidity: function TriggerEvent(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, bytes32 account, bytes32[] accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct) returns()
func (_Testfiles *TestfilesSession) TriggerEvent(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account [32]byte, accounts [][32]byte, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
	return _Testfiles.Contract.TriggerEvent(&_Testfiles.TransactOpts, field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct)
}

// TriggerEvent is a paid mutator transaction binding the contract method 0xe669831d.
//
// Solidity: function TriggerEvent(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, bytes32 account, bytes32[] accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct) returns()
func (_Testfiles *TestfilesTransactorSession) TriggerEvent(field int32, differentField string, oracleId uint8, oracleIds [32]uint8, account [32]byte, accounts [][32]byte, bigField *big.Int, nestedStruct MidLevelTestStruct) (*types.Transaction, error) {
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
	Account        [32]byte
	Accounts       [][32]byte
	BigField       *big.Int
	NestedStruct   MidLevelTestStruct
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterTriggered is a free log retrieval operation binding the contract event 0x7d2abe6109e46b893ac1835c9287d6ad5c5ccf3d0254d2ca72225873795e0f53.
//
// Solidity: event Triggered(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, bytes32 account, bytes32[] accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct)
func (_Testfiles *TestfilesFilterer) FilterTriggered(opts *bind.FilterOpts) (*TestfilesTriggeredIterator, error) {

	logs, sub, err := _Testfiles.contract.FilterLogs(opts, "Triggered")
	if err != nil {
		return nil, err
	}
	return &TestfilesTriggeredIterator{contract: _Testfiles.contract, event: "Triggered", logs: logs, sub: sub}, nil
}

// WatchTriggered is a free log subscription operation binding the contract event 0x7d2abe6109e46b893ac1835c9287d6ad5c5ccf3d0254d2ca72225873795e0f53.
//
// Solidity: event Triggered(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, bytes32 account, bytes32[] accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct)
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

// ParseTriggered is a log parse operation binding the contract event 0x7d2abe6109e46b893ac1835c9287d6ad5c5ccf3d0254d2ca72225873795e0f53.
//
// Solidity: event Triggered(int32 field, string differentField, uint8 oracleId, uint8[32] oracleIds, bytes32 account, bytes32[] accounts, int192 bigField, (bytes2,(int64,string)) nestedStruct)
func (_Testfiles *TestfilesFilterer) ParseTriggered(log types.Log) (*TestfilesTriggered, error) {
	event := new(TestfilesTriggered)
	if err := _Testfiles.contract.UnpackLog(event, "Triggered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
