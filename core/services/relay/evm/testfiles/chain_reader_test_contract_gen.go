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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"account\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"accounts\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"indexed\":false,\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"Triggered\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"bytes32\",\"name\":\"account\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"accounts\",\"type\":\"bytes32[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"AddTestStruct\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetDifferentPrimitiveValue\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"i\",\"type\":\"uint256\"}],\"name\":\"GetElementAtIndex\",\"outputs\":[{\"components\":[{\"internalType\":\"int32\",\"name\":\"Field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"DifferentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"OracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"OracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"bytes32\",\"name\":\"Account\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"Accounts\",\"type\":\"bytes32[]\"},{\"internalType\":\"int192\",\"name\":\"BigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"NestedStruct\",\"type\":\"tuple\"}],\"internalType\":\"structTestStruct\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetPrimitiveValue\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetSliceValue\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"bytes32\",\"name\":\"account\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"accounts\",\"type\":\"bytes32[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"ReturnSeen\",\"outputs\":[{\"components\":[{\"internalType\":\"int32\",\"name\":\"Field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"DifferentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"OracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"OracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"bytes32\",\"name\":\"Account\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"Accounts\",\"type\":\"bytes32[]\"},{\"internalType\":\"int192\",\"name\":\"BigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"NestedStruct\",\"type\":\"tuple\"}],\"internalType\":\"structTestStruct\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int32\",\"name\":\"field\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"differentField\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"oracleId\",\"type\":\"uint8\"},{\"internalType\":\"uint8[32]\",\"name\":\"oracleIds\",\"type\":\"uint8[32]\"},{\"internalType\":\"bytes32\",\"name\":\"account\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"accounts\",\"type\":\"bytes32[]\"},{\"internalType\":\"int192\",\"name\":\"bigField\",\"type\":\"int192\"},{\"components\":[{\"internalType\":\"bytes2\",\"name\":\"FixedBytes\",\"type\":\"bytes2\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"I\",\"type\":\"int64\"},{\"internalType\":\"string\",\"name\":\"S\",\"type\":\"string\"}],\"internalType\":\"structInnerTestStruct\",\"name\":\"Inner\",\"type\":\"tuple\"}],\"internalType\":\"structMidLevelTestStruct\",\"name\":\"nestedStruct\",\"type\":\"tuple\"}],\"name\":\"TriggerEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50600180548082018255600082905260048082047fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6908101805460086003958616810261010090810a8088026001600160401b0391820219909416939093179093558654808801909755848704909301805496909516909202900a9182029102199092169190911790556111ba806100a96000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c8063b95ad4111161005b578063b95ad411146100d9578063bdb37c90146100ec578063da8e7a8214610101578063e669831d1461010857600080fd5b8063030d3ca2146100825780637dd6af5b146100a45780639ca04f67146100b9575b600080fd5b6107c65b60405167ffffffffffffffff90911681526020015b60405180910390f35b6100b76100b2366004610af9565b61011b565b005b6100cc6100c7366004610be4565b61034a565b60405161009b9190610eff565b6100cc6100e7366004610af9565b6105fa565b6100f46106eb565b60405161009b9190610e03565b6003610086565b6100b7610116366004610af9565b610777565b60006040518061010001604052808c60030b81526020018b8b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060ff8a166020808301919091526040805161040081810183529190930192918b918390839080828437600092019190915250505081526020808201899052604080518883028181018401835289825291909301929189918991829190850190849080828437600092019190915250505090825250601785900b60208201526040016101f584611036565b905281546001818101845560009384526020938490208351600a90930201805460039390930b63ffffffff1663ffffffff199093169290921782558383015180519394929361024c939285019291909101906107cc565b50604082015160028201805460ff191660ff909216919091179055606082015161027c9060038301906020610850565b506080820151600482015560a082015180516102a29160058401916020909101906108de565b5060c082015160068201805460179290920b6001600160c01b03166001600160c01b031990921691909117905560e082015180516007808401805460f09390931c61ffff1990931692909217825560208084015180516008870180549190940b67ffffffffffffffff1667ffffffffffffffff19909116178355808201518051919392610337926009890192909101906107cc565b5050505050505050505050505050505050565b610352610918565b600061035f600184611011565b8154811061036f5761036f611158565b90600052602060002090600a0201604051806101000160405290816000820160009054906101000a900460030b60030b60030b81526020016001820180546103b690611123565b80601f01602080910402602001604051908101604052809291908181526020018280546103e290611123565b801561042f5780601f106104045761010080835404028352916020019161042f565b820191906000526020600020905b81548152906001019060200180831161041257829003601f168201915b5050509183525050600282015460ff166020808301919091526040805161040081018083529190930192916003850191826000855b825461010083900a900460ff168152602060019283018181049485019490930390920291018084116104645790505050505050815260200160048201548152602001600582018054806020026020016040519081016040528092919081815260200182805480156104f457602002820191906000526020600020905b8154815260200190600101908083116104e0575b50505091835250506006820154601790810b810b900b6020808301919091526040805180820182526007808601805460f01b6001600160f01b031916835283518085018552600888018054840b840b90930b8152600988018054959097019693959194868301949193928401919061056b90611123565b80601f016020809104026020016040519081016040528092919081815260200182805461059790611123565b80156105e45780601f106105b9576101008083540402835291602001916105e4565b820191906000526020600020905b8154815290600101906020018083116105c757829003601f168201915b5050509190925250505090525090525092915050565b610602610918565b6040518061010001604052808c60030b81526020018b8b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060ff8a166020808301919091526040805161040081810183529190930192918b918390839080828437600092019190915250505081526020808201899052604080518883028181018401835289825291909301929189918991829190850190849080828437600092019190915250505090825250601785900b60208201526040016106da84611036565b90529b9a5050505050505050505050565b6060600180548060200260200160405190810160405280929190818152602001828054801561076d57602002820191906000526020600020906000905b82829054906101000a900467ffffffffffffffff1667ffffffffffffffff16815260200190600801906020826007010492830192600103820291508084116107285790505b5050505050905090565b7f7d2abe6109e46b893ac1835c9287d6ad5c5ccf3d0254d2ca72225873795e0f538a8a8a8a8a8a8a8a8a8a6040516107b89a99989796959493929190610e51565b60405180910390a150505050505050505050565b8280546107d890611123565b90600052602060002090601f0160209004810192826107fa5760008555610840565b82601f1061081357805160ff1916838001178555610840565b82800160010185558215610840579182015b82811115610840578251825591602001919060010190610825565b5061084c929150610967565b5090565b6001830191839082156108405791602002820160005b838211156108a457835183826101000a81548160ff021916908360ff1602179055509260200192600101602081600001049283019260010302610866565b80156108d15782816101000a81549060ff02191690556001016020816000010492830192600103026108a4565b505061084c929150610967565b8280548282559060005260206000209081019282156108405791602002820182811115610840578251825591602001919060010190610825565b604080516101008101825260008082526060602083018190529282015290810161094061097c565b815260006020820181905260606040830181905282015260800161096261099b565b905290565b5b8082111561084c5760008155600101610968565b6040518061040001604052806020906020820280368337509192915050565b604051806040016040528060006001600160f01b03191681526020016109626040518060400160405280600060070b8152602001606081525090565b60008083601f8401126109e957600080fd5b50813567ffffffffffffffff811115610a0157600080fd5b6020830191508360208260051b8501011115610a1c57600080fd5b9250929050565b806104008101831015610a3557600080fd5b92915050565b80356001600160f01b031981168114610a5357600080fd5b919050565b8035601781900b8114610a5357600080fd5b8035600381900b8114610a5357600080fd5b8035600781900b8114610a5357600080fd5b60008083601f840112610aa057600080fd5b50813567ffffffffffffffff811115610ab857600080fd5b602083019150836020828501011115610a1c57600080fd5b600060408284031215610ae257600080fd5b50919050565b803560ff81168114610a5357600080fd5b6000806000806000806000806000806104e08b8d031215610b1957600080fd5b610b228b610a6a565b995060208b013567ffffffffffffffff80821115610b3f57600080fd5b610b4b8e838f01610a8e565b909b509950899150610b5f60408e01610ae8565b9850610b6e8e60608f01610a23565b97506104608d013596506104808d0135915080821115610b8d57600080fd5b610b998e838f016109d7565b9096509450849150610bae6104a08e01610a58565b93506104c08d0135915080821115610bc557600080fd5b50610bd28d828e01610ad0565b9150509295989b9194979a5092959850565b600060208284031215610bf657600080fd5b5035919050565b81835260006001600160fb1b03831115610c1657600080fd5b8260051b8083602087013760009401602001938452509192915050565b600081518084526020808501945080840160005b83811015610c6357815187529582019590820190600101610c47565b509495945050505050565b8060005b6020808210610c815750610c98565b825160ff1685529384019390910190600101610c72565b50505050565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b6000815180845260005b81811015610ced57602081850181015186830182015201610cd1565b81811115610cff576000602083870101525b50601f01601f19169290920160200192915050565b6001600160f01b0319610d2682610a3b565b16825260006020820135603e19833603018112610d4257600080fd5b604060208501528201610d5481610a7c565b60070b60408501526020810135601e19823603018112610d7357600080fd5b8101803567ffffffffffffffff811115610d8c57600080fd5b803603831315610d9b57600080fd5b60406060870152610db3608087018260208501610c9e565b9695505050505050565b61ffff60f01b81511682526000602082015160406020850152805160070b60408501526020810151905060406060850152610dfb6080850182610cc7565b949350505050565b6020808252825182820181905260009190848201906040850190845b81811015610e4557835167ffffffffffffffff1683529284019291840191600101610e1f565b50909695505050505050565b60006104e08c60030b835260208181850152610e708285018d8f610c9e565b915060ff808c166040860152606085018b60005b84811015610ea95783610e9683610ae8565b1683529184019190840190600101610e84565b505050505087610460840152828103610480840152610ec9818789610bfd565b9050610edb6104a084018660170b9052565b8281036104c0840152610eee8185610d14565b9d9c50505050505050505050505050565b60208152610f1360208201835160030b9052565b600060208301516104e0806040850152610f31610500850183610cc7565b91506040850151610f47606086018260ff169052565b506060850151610f5a6080860182610c6e565b50608085015161048085015260a0850151601f1980868503016104a0870152610f838483610c33565b935060c08701519150610f9c6104c087018360170b9052565b60e0870151915080868503018387015250610db38382610dbd565b6040805190810167ffffffffffffffff81118282101715610fda57610fda61116e565b60405290565b604051601f8201601f1916810167ffffffffffffffff811182821017156110095761100961116e565b604052919050565b60008282101561103157634e487b7160e01b600052601160045260246000fd5b500390565b60006040823603121561104857600080fd5b611050610fb7565b61105983610a3b565b815260208084013567ffffffffffffffff8082111561107757600080fd5b81860191506040823603121561108c57600080fd5b611094610fb7565b61109d83610a7c565b815283830135828111156110b057600080fd5b929092019136601f8401126110c457600080fd5b8235828111156110d6576110d661116e565b6110e8601f8201601f19168601610fe0565b925080835236858286010111156110fe57600080fd5b8085850186850137600090830185015280840191909152918301919091525092915050565b600181811c9082168061113757607f821691505b60208210811415610ae257634e487b7160e01b600052602260045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052604160045260246000fdfea26469706673582212207f1ef47c84498df82a83ce9c4eed62d431c29fbd05ad8899d4ef1dfcb22f147e64736f6c63430008060033",
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
