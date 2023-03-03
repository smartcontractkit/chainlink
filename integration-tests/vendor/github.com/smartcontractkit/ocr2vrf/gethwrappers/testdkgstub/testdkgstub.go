package testdkgstub

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

type KeyDataStructKeyData struct {
	PublicKey []byte
	Hashes    [][32]byte
}

var DKGClientMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structKeyDataStruct.KeyData\",\"name\":\"kd\",\"type\":\"tuple\"}],\"name\":\"keyGenerated\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"newKeyRequested\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

var DKGClientABI = DKGClientMetaData.ABI

type DKGClient struct {
	DKGClientCaller
	DKGClientTransactor
	DKGClientFilterer
}

type DKGClientCaller struct {
	contract *bind.BoundContract
}

type DKGClientTransactor struct {
	contract *bind.BoundContract
}

type DKGClientFilterer struct {
	contract *bind.BoundContract
}

type DKGClientSession struct {
	Contract     *DKGClient
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type DKGClientCallerSession struct {
	Contract *DKGClientCaller
	CallOpts bind.CallOpts
}

type DKGClientTransactorSession struct {
	Contract     *DKGClientTransactor
	TransactOpts bind.TransactOpts
}

type DKGClientRaw struct {
	Contract *DKGClient
}

type DKGClientCallerRaw struct {
	Contract *DKGClientCaller
}

type DKGClientTransactorRaw struct {
	Contract *DKGClientTransactor
}

func NewDKGClient(address common.Address, backend bind.ContractBackend) (*DKGClient, error) {
	contract, err := bindDKGClient(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DKGClient{DKGClientCaller: DKGClientCaller{contract: contract}, DKGClientTransactor: DKGClientTransactor{contract: contract}, DKGClientFilterer: DKGClientFilterer{contract: contract}}, nil
}

func NewDKGClientCaller(address common.Address, caller bind.ContractCaller) (*DKGClientCaller, error) {
	contract, err := bindDKGClient(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DKGClientCaller{contract: contract}, nil
}

func NewDKGClientTransactor(address common.Address, transactor bind.ContractTransactor) (*DKGClientTransactor, error) {
	contract, err := bindDKGClient(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DKGClientTransactor{contract: contract}, nil
}

func NewDKGClientFilterer(address common.Address, filterer bind.ContractFilterer) (*DKGClientFilterer, error) {
	contract, err := bindDKGClient(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DKGClientFilterer{contract: contract}, nil
}

func bindDKGClient(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DKGClientABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_DKGClient *DKGClientRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DKGClient.Contract.DKGClientCaller.contract.Call(opts, result, method, params...)
}

func (_DKGClient *DKGClientRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DKGClient.Contract.DKGClientTransactor.contract.Transfer(opts)
}

func (_DKGClient *DKGClientRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DKGClient.Contract.DKGClientTransactor.contract.Transact(opts, method, params...)
}

func (_DKGClient *DKGClientCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DKGClient.Contract.contract.Call(opts, result, method, params...)
}

func (_DKGClient *DKGClientTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DKGClient.Contract.contract.Transfer(opts)
}

func (_DKGClient *DKGClientTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DKGClient.Contract.contract.Transact(opts, method, params...)
}

func (_DKGClient *DKGClientTransactor) KeyGenerated(opts *bind.TransactOpts, kd KeyDataStructKeyData) (*types.Transaction, error) {
	return _DKGClient.contract.Transact(opts, "keyGenerated", kd)
}

func (_DKGClient *DKGClientSession) KeyGenerated(kd KeyDataStructKeyData) (*types.Transaction, error) {
	return _DKGClient.Contract.KeyGenerated(&_DKGClient.TransactOpts, kd)
}

func (_DKGClient *DKGClientTransactorSession) KeyGenerated(kd KeyDataStructKeyData) (*types.Transaction, error) {
	return _DKGClient.Contract.KeyGenerated(&_DKGClient.TransactOpts, kd)
}

func (_DKGClient *DKGClientTransactor) NewKeyRequested(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DKGClient.contract.Transact(opts, "newKeyRequested")
}

func (_DKGClient *DKGClientSession) NewKeyRequested() (*types.Transaction, error) {
	return _DKGClient.Contract.NewKeyRequested(&_DKGClient.TransactOpts)
}

func (_DKGClient *DKGClientTransactorSession) NewKeyRequested() (*types.Transaction, error) {
	return _DKGClient.Contract.NewKeyRequested(&_DKGClient.TransactOpts)
}

var KeyDataStructMetaData = &bind.MetaData{
	ABI: "[]",
}

var KeyDataStructABI = KeyDataStructMetaData.ABI

type KeyDataStruct struct {
	KeyDataStructCaller
	KeyDataStructTransactor
	KeyDataStructFilterer
}

type KeyDataStructCaller struct {
	contract *bind.BoundContract
}

type KeyDataStructTransactor struct {
	contract *bind.BoundContract
}

type KeyDataStructFilterer struct {
	contract *bind.BoundContract
}

type KeyDataStructSession struct {
	Contract     *KeyDataStruct
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeyDataStructCallerSession struct {
	Contract *KeyDataStructCaller
	CallOpts bind.CallOpts
}

type KeyDataStructTransactorSession struct {
	Contract     *KeyDataStructTransactor
	TransactOpts bind.TransactOpts
}

type KeyDataStructRaw struct {
	Contract *KeyDataStruct
}

type KeyDataStructCallerRaw struct {
	Contract *KeyDataStructCaller
}

type KeyDataStructTransactorRaw struct {
	Contract *KeyDataStructTransactor
}

func NewKeyDataStruct(address common.Address, backend bind.ContractBackend) (*KeyDataStruct, error) {
	contract, err := bindKeyDataStruct(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeyDataStruct{KeyDataStructCaller: KeyDataStructCaller{contract: contract}, KeyDataStructTransactor: KeyDataStructTransactor{contract: contract}, KeyDataStructFilterer: KeyDataStructFilterer{contract: contract}}, nil
}

func NewKeyDataStructCaller(address common.Address, caller bind.ContractCaller) (*KeyDataStructCaller, error) {
	contract, err := bindKeyDataStruct(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeyDataStructCaller{contract: contract}, nil
}

func NewKeyDataStructTransactor(address common.Address, transactor bind.ContractTransactor) (*KeyDataStructTransactor, error) {
	contract, err := bindKeyDataStruct(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeyDataStructTransactor{contract: contract}, nil
}

func NewKeyDataStructFilterer(address common.Address, filterer bind.ContractFilterer) (*KeyDataStructFilterer, error) {
	contract, err := bindKeyDataStruct(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeyDataStructFilterer{contract: contract}, nil
}

func bindKeyDataStruct(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KeyDataStructABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_KeyDataStruct *KeyDataStructRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeyDataStruct.Contract.KeyDataStructCaller.contract.Call(opts, result, method, params...)
}

func (_KeyDataStruct *KeyDataStructRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeyDataStruct.Contract.KeyDataStructTransactor.contract.Transfer(opts)
}

func (_KeyDataStruct *KeyDataStructRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeyDataStruct.Contract.KeyDataStructTransactor.contract.Transact(opts, method, params...)
}

func (_KeyDataStruct *KeyDataStructCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeyDataStruct.Contract.contract.Call(opts, result, method, params...)
}

func (_KeyDataStruct *KeyDataStructTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeyDataStruct.Contract.contract.Transfer(opts)
}

func (_KeyDataStruct *KeyDataStructTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeyDataStruct.Contract.contract.Transact(opts, method, params...)
}

var TestDKGStubMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"key\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"keyID\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"contractDKGClient\",\"name\":\"clientAddress\",\"type\":\"address\"}],\"name\":\"addClient\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structKeyDataStruct.KeyData\",\"name\":\"kd\",\"type\":\"tuple\"}],\"name\":\"keyGenerated\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161058b38038061058b83398101604081905261002f9161005b565b600061003b83826101bc565b506001555061027b565b634e487b7160e01b600052604160045260246000fd5b6000806040838503121561006e57600080fd5b82516001600160401b038082111561008557600080fd5b818501915085601f83011261009957600080fd5b8151818111156100ab576100ab610045565b604051601f8201601f19908116603f011681019083821181831017156100d3576100d3610045565b816040528281526020935088848487010111156100ef57600080fd5b600091505b8282101561011157848201840151818301850152908301906100f4565b828211156101225760008484830101525b969092015195979596505050505050565b600181811c9082168061014757607f821691505b60208210810361016757634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156101b757600081815260208120601f850160051c810160208610156101945750805b601f850160051c820191505b818110156101b3578281556001016101a0565b5050505b505050565b81516001600160401b038111156101d5576101d5610045565b6101e9816101e38454610133565b8461016d565b602080601f83116001811461021e57600084156102065750858301515b600019600386901b1c1916600185901b1785556101b3565b600085815260208120601f198616915b8281101561024d5788860151825594840194600190910190840161022e565b508582101561026b5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6103018061028a6000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80637bf1ffc51461003b578063bf2732c714610093575b600080fd5b61009161004936600461012c565b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9290921691909117905550565b005b6100916100a1366004610175565b6002546040517fbf2732c700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091169063bf2732c7906100f7908490600401610256565b600060405180830381600087803b15801561011157600080fd5b505af1158015610125573d6000803e3d6000fd5b5050505050565b6000806040838503121561013f57600080fd5b82359150602083013573ffffffffffffffffffffffffffffffffffffffff8116811461016a57600080fd5b809150509250929050565b60006020828403121561018757600080fd5b813567ffffffffffffffff81111561019e57600080fd5b8201604081850312156101b057600080fd5b9392505050565b6000808335601e198436030181126101ce57600080fd5b830160208101925035905067ffffffffffffffff8111156101ee57600080fd5b8060051b360382131561020057600080fd5b9250929050565b81835260007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83111561023957600080fd5b8260051b8083602087013760009401602001938452509192915050565b6020815260008235601e1984360301811261027057600080fd5b830160208101903567ffffffffffffffff81111561028d57600080fd5b80360382131561029c57600080fd5b6040602085015280606085015280826080860137600084820160800152601f01601f1916830190506102d160208501856101b7565b60608584030160408601526102ea608084018284610207565b969550505050505056fea164736f6c634300080f000a",
}

var TestDKGStubABI = TestDKGStubMetaData.ABI

var TestDKGStubBin = TestDKGStubMetaData.Bin

func DeployTestDKGStub(auth *bind.TransactOpts, backend bind.ContractBackend, key []byte, keyID [32]byte) (common.Address, *types.Transaction, *TestDKGStub, error) {
	parsed, err := TestDKGStubMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestDKGStubBin), backend, key, keyID)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestDKGStub{TestDKGStubCaller: TestDKGStubCaller{contract: contract}, TestDKGStubTransactor: TestDKGStubTransactor{contract: contract}, TestDKGStubFilterer: TestDKGStubFilterer{contract: contract}}, nil
}

type TestDKGStub struct {
	TestDKGStubCaller
	TestDKGStubTransactor
	TestDKGStubFilterer
}

type TestDKGStubCaller struct {
	contract *bind.BoundContract
}

type TestDKGStubTransactor struct {
	contract *bind.BoundContract
}

type TestDKGStubFilterer struct {
	contract *bind.BoundContract
}

type TestDKGStubSession struct {
	Contract     *TestDKGStub
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type TestDKGStubCallerSession struct {
	Contract *TestDKGStubCaller
	CallOpts bind.CallOpts
}

type TestDKGStubTransactorSession struct {
	Contract     *TestDKGStubTransactor
	TransactOpts bind.TransactOpts
}

type TestDKGStubRaw struct {
	Contract *TestDKGStub
}

type TestDKGStubCallerRaw struct {
	Contract *TestDKGStubCaller
}

type TestDKGStubTransactorRaw struct {
	Contract *TestDKGStubTransactor
}

func NewTestDKGStub(address common.Address, backend bind.ContractBackend) (*TestDKGStub, error) {
	contract, err := bindTestDKGStub(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestDKGStub{TestDKGStubCaller: TestDKGStubCaller{contract: contract}, TestDKGStubTransactor: TestDKGStubTransactor{contract: contract}, TestDKGStubFilterer: TestDKGStubFilterer{contract: contract}}, nil
}

func NewTestDKGStubCaller(address common.Address, caller bind.ContractCaller) (*TestDKGStubCaller, error) {
	contract, err := bindTestDKGStub(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestDKGStubCaller{contract: contract}, nil
}

func NewTestDKGStubTransactor(address common.Address, transactor bind.ContractTransactor) (*TestDKGStubTransactor, error) {
	contract, err := bindTestDKGStub(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestDKGStubTransactor{contract: contract}, nil
}

func NewTestDKGStubFilterer(address common.Address, filterer bind.ContractFilterer) (*TestDKGStubFilterer, error) {
	contract, err := bindTestDKGStub(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestDKGStubFilterer{contract: contract}, nil
}

func bindTestDKGStub(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TestDKGStubABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_TestDKGStub *TestDKGStubRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestDKGStub.Contract.TestDKGStubCaller.contract.Call(opts, result, method, params...)
}

func (_TestDKGStub *TestDKGStubRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestDKGStub.Contract.TestDKGStubTransactor.contract.Transfer(opts)
}

func (_TestDKGStub *TestDKGStubRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestDKGStub.Contract.TestDKGStubTransactor.contract.Transact(opts, method, params...)
}

func (_TestDKGStub *TestDKGStubCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestDKGStub.Contract.contract.Call(opts, result, method, params...)
}

func (_TestDKGStub *TestDKGStubTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestDKGStub.Contract.contract.Transfer(opts)
}

func (_TestDKGStub *TestDKGStubTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestDKGStub.Contract.contract.Transact(opts, method, params...)
}

func (_TestDKGStub *TestDKGStubTransactor) AddClient(opts *bind.TransactOpts, arg0 [32]byte, clientAddress common.Address) (*types.Transaction, error) {
	return _TestDKGStub.contract.Transact(opts, "addClient", arg0, clientAddress)
}

func (_TestDKGStub *TestDKGStubSession) AddClient(arg0 [32]byte, clientAddress common.Address) (*types.Transaction, error) {
	return _TestDKGStub.Contract.AddClient(&_TestDKGStub.TransactOpts, arg0, clientAddress)
}

func (_TestDKGStub *TestDKGStubTransactorSession) AddClient(arg0 [32]byte, clientAddress common.Address) (*types.Transaction, error) {
	return _TestDKGStub.Contract.AddClient(&_TestDKGStub.TransactOpts, arg0, clientAddress)
}

func (_TestDKGStub *TestDKGStubTransactor) KeyGenerated(opts *bind.TransactOpts, kd KeyDataStructKeyData) (*types.Transaction, error) {
	return _TestDKGStub.contract.Transact(opts, "keyGenerated", kd)
}

func (_TestDKGStub *TestDKGStubSession) KeyGenerated(kd KeyDataStructKeyData) (*types.Transaction, error) {
	return _TestDKGStub.Contract.KeyGenerated(&_TestDKGStub.TransactOpts, kd)
}

func (_TestDKGStub *TestDKGStubTransactorSession) KeyGenerated(kd KeyDataStructKeyData) (*types.Transaction, error) {
	return _TestDKGStub.Contract.KeyGenerated(&_TestDKGStub.TransactOpts, kd)
}
