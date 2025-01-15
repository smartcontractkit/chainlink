// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package multicall3

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
	_ = abi.ConvertType
)

type Multicall3Call struct {
	Target   common.Address
	CallData []byte
}

type Multicall3Call3 struct {
	Target       common.Address
	AllowFailure bool
	CallData     []byte
}

type Multicall3Call3Value struct {
	Target       common.Address
	AllowFailure bool
	Value        *big.Int
	CallData     []byte
}

type Multicall3Result struct {
	Success    bool
	ReturnData []byte
}

var Multicall3MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"aggregate\",\"inputs\":[{\"name\":\"calls\",\"type\":\"tuple[]\",\"internalType\":\"structMulticall3.Call[]\",\"components\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"callData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"returnData\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"aggregate3\",\"inputs\":[{\"name\":\"calls\",\"type\":\"tuple[]\",\"internalType\":\"structMulticall3.Call3[]\",\"components\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowFailure\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"callData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"returnData\",\"type\":\"tuple[]\",\"internalType\":\"structMulticall3.Result[]\",\"components\":[{\"name\":\"success\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"returnData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"aggregate3Value\",\"inputs\":[{\"name\":\"calls\",\"type\":\"tuple[]\",\"internalType\":\"structMulticall3.Call3Value[]\",\"components\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowFailure\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"callData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"returnData\",\"type\":\"tuple[]\",\"internalType\":\"structMulticall3.Result[]\",\"components\":[{\"name\":\"success\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"returnData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"blockAndAggregate\",\"inputs\":[{\"name\":\"calls\",\"type\":\"tuple[]\",\"internalType\":\"structMulticall3.Call[]\",\"components\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"callData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"returnData\",\"type\":\"tuple[]\",\"internalType\":\"structMulticall3.Result[]\",\"components\":[{\"name\":\"success\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"returnData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"getBasefee\",\"inputs\":[],\"outputs\":[{\"name\":\"basefee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBlockHash\",\"inputs\":[{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"blockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBlockNumber\",\"inputs\":[],\"outputs\":[{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getChainId\",\"inputs\":[],\"outputs\":[{\"name\":\"chainid\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentBlockCoinbase\",\"inputs\":[],\"outputs\":[{\"name\":\"coinbase\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentBlockDifficulty\",\"inputs\":[],\"outputs\":[{\"name\":\"difficulty\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentBlockGasLimit\",\"inputs\":[],\"outputs\":[{\"name\":\"gaslimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentBlockTimestamp\",\"inputs\":[],\"outputs\":[{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEthBalance\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLastBlockHash\",\"inputs\":[],\"outputs\":[{\"name\":\"blockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tryAggregate\",\"inputs\":[{\"name\":\"requireSuccess\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"calls\",\"type\":\"tuple[]\",\"internalType\":\"structMulticall3.Call[]\",\"components\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"callData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"returnData\",\"type\":\"tuple[]\",\"internalType\":\"structMulticall3.Result[]\",\"components\":[{\"name\":\"success\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"returnData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"tryBlockAndAggregate\",\"inputs\":[{\"name\":\"requireSuccess\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"calls\",\"type\":\"tuple[]\",\"internalType\":\"structMulticall3.Call[]\",\"components\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"callData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"returnData\",\"type\":\"tuple[]\",\"internalType\":\"structMulticall3.Result[]\",\"components\":[{\"name\":\"success\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"returnData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"payable\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610eb0806100206000396000f3fe6080604052600436106100f35760003560e01c80634d2301cc1161008a578063a8b0574e11610059578063a8b0574e1461025a578063bce38bd714610275578063c3077fa914610288578063ee82ac5e1461029b57600080fd5b80634d2301cc146101ec57806372425d9d1461022157806382ad56cb1461023457806386d516e81461024757600080fd5b80633408e470116100c65780633408e47014610191578063399542e9146101a45780633e64a696146101c657806342cbb15c146101d957600080fd5b80630f28c97d146100f8578063174dea711461011a578063252dba421461013a57806327e86d6e1461015b575b600080fd5b34801561010457600080fd5b50425b6040519081526020015b60405180910390f35b61012d610128366004610a85565b6102ba565b6040516101119190610bb7565b61014d610148366004610a85565b6104ef565b604051610111929190610bd1565b34801561016757600080fd5b50437fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0140610107565b34801561019d57600080fd5b5046610107565b6101b76101b2366004610c59565b610690565b60405161011193929190610cb3565b3480156101d257600080fd5b5048610107565b3480156101e557600080fd5b5043610107565b3480156101f857600080fd5b50610107610207366004610cdb565b73ffffffffffffffffffffffffffffffffffffffff163190565b34801561022d57600080fd5b5044610107565b61012d610242366004610a85565b6106ab565b34801561025357600080fd5b5045610107565b34801561026657600080fd5b50604051418152602001610111565b61012d610283366004610c59565b61085a565b6101b7610296366004610a85565b610a1a565b3480156102a757600080fd5b506101076102b6366004610d11565b4090565b60606000828067ffffffffffffffff8111156102d8576102d8610d2a565b60405190808252806020026020018201604052801561031e57816020015b6040805180820190915260008152606060208201528152602001906001900390816102f65790505b5092503660005b8281101561047757600085828151811061034157610341610d59565b6020026020010151905087878381811061035d5761035d610d59565b905060200281019061036f9190610d88565b6040810135958601959093506103886020850185610cdb565b73ffffffffffffffffffffffffffffffffffffffff16816103ac6060870187610dc6565b6040516103ba929190610e2b565b60006040518083038185875af1925050503d80600081146103f7576040519150601f19603f3d011682016040523d82523d6000602084013e6103fc565b606091505b50602080850191909152901515808452908501351761046d577f08c379a000000000000000000000000000000000000000000000000000000000600052602060045260176024527f4d756c746963616c6c333a2063616c6c206661696c656400000000000000000060445260846000fd5b5050600101610325565b508234146104e6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f4d756c746963616c6c333a2076616c7565206d69736d6174636800000000000060448201526064015b60405180910390fd5b50505092915050565b436060828067ffffffffffffffff81111561050c5761050c610d2a565b60405190808252806020026020018201604052801561053f57816020015b606081526020019060019003908161052a5790505b5091503660005b8281101561068657600087878381811061056257610562610d59565b90506020028101906105749190610e3b565b92506105836020840184610cdb565b73ffffffffffffffffffffffffffffffffffffffff166105a66020850185610dc6565b6040516105b4929190610e2b565b6000604051808303816000865af19150503d80600081146105f1576040519150601f19603f3d011682016040523d82523d6000602084013e6105f6565b606091505b5086848151811061060957610609610d59565b602090810291909101015290508061067d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f4d756c746963616c6c333a2063616c6c206661696c656400000000000000000060448201526064016104dd565b50600101610546565b5050509250929050565b43804060606106a086868661085a565b905093509350939050565b6060818067ffffffffffffffff8111156106c7576106c7610d2a565b60405190808252806020026020018201604052801561070d57816020015b6040805180820190915260008152606060208201528152602001906001900390816106e55790505b5091503660005b828110156104e657600084828151811061073057610730610d59565b6020026020010151905086868381811061074c5761074c610d59565b905060200281019061075e9190610e6f565b925061076d6020840184610cdb565b73ffffffffffffffffffffffffffffffffffffffff166107906040850185610dc6565b60405161079e929190610e2b565b6000604051808303816000865af19150503d80600081146107db576040519150601f19603f3d011682016040523d82523d6000602084013e6107e0565b606091505b506020808401919091529015158083529084013517610851577f08c379a000000000000000000000000000000000000000000000000000000000600052602060045260176024527f4d756c746963616c6c333a2063616c6c206661696c656400000000000000000060445260646000fd5b50600101610714565b6060818067ffffffffffffffff81111561087657610876610d2a565b6040519080825280602002602001820160405280156108bc57816020015b6040805180820190915260008152606060208201528152602001906001900390816108945790505b5091503660005b82811015610a105760008482815181106108df576108df610d59565b602002602001015190508686838181106108fb576108fb610d59565b905060200281019061090d9190610e3b565b925061091c6020840184610cdb565b73ffffffffffffffffffffffffffffffffffffffff1661093f6020850185610dc6565b60405161094d929190610e2b565b6000604051808303816000865af19150503d806000811461098a576040519150601f19603f3d011682016040523d82523d6000602084013e61098f565b606091505b506020830152151581528715610a07578051610a07576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f4d756c746963616c6c333a2063616c6c206661696c656400000000000000000060448201526064016104dd565b506001016108c3565b5050509392505050565b6000806060610a2b60018686610690565b919790965090945092505050565b60008083601f840112610a4b57600080fd5b50813567ffffffffffffffff811115610a6357600080fd5b6020830191508360208260051b8501011115610a7e57600080fd5b9250929050565b60008060208385031215610a9857600080fd5b823567ffffffffffffffff811115610aaf57600080fd5b610abb85828601610a39565b90969095509350505050565b6000815180845260005b81811015610aed57602081850181015186830182015201610ad1565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b600082825180855260208086019550808260051b84010181860160005b84811015610baa578583037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001895281518051151584528401516040858501819052610b9681860183610ac7565b9a86019a9450505090830190600101610b48565b5090979650505050505050565b602081526000610bca6020830184610b2b565b9392505050565b600060408201848352602060408185015281855180845260608601915060608160051b870101935082870160005b82811015610c4b577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0888703018452610c39868351610ac7565b95509284019290840190600101610bff565b509398975050505050505050565b600080600060408486031215610c6e57600080fd5b83358015158114610c7e57600080fd5b9250602084013567ffffffffffffffff811115610c9a57600080fd5b610ca686828701610a39565b9497909650939450505050565b838152826020820152606060408201526000610cd26060830184610b2b565b95945050505050565b600060208284031215610ced57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff81168114610bca57600080fd5b600060208284031215610d2357600080fd5b5035919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81833603018112610dbc57600080fd5b9190910192915050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610dfb57600080fd5b83018035915067ffffffffffffffff821115610e1657600080fd5b602001915036819003821315610a7e57600080fd5b8183823760009101908152919050565b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc1833603018112610dbc57600080fd5b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa1833603018112610dbc57600080fdfea164736f6c6343000813000a",
}

var Multicall3ABI = Multicall3MetaData.ABI

var Multicall3Bin = Multicall3MetaData.Bin

func DeployMulticall3(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Multicall3, error) {
	parsed, err := Multicall3MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(Multicall3Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Multicall3{address: address, abi: *parsed, Multicall3Caller: Multicall3Caller{contract: contract}, Multicall3Transactor: Multicall3Transactor{contract: contract}, Multicall3Filterer: Multicall3Filterer{contract: contract}}, nil
}

type Multicall3 struct {
	address common.Address
	abi     abi.ABI
	Multicall3Caller
	Multicall3Transactor
	Multicall3Filterer
}

type Multicall3Caller struct {
	contract *bind.BoundContract
}

type Multicall3Transactor struct {
	contract *bind.BoundContract
}

type Multicall3Filterer struct {
	contract *bind.BoundContract
}

type Multicall3Session struct {
	Contract     *Multicall3
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type Multicall3CallerSession struct {
	Contract *Multicall3Caller
	CallOpts bind.CallOpts
}

type Multicall3TransactorSession struct {
	Contract     *Multicall3Transactor
	TransactOpts bind.TransactOpts
}

type Multicall3Raw struct {
	Contract *Multicall3
}

type Multicall3CallerRaw struct {
	Contract *Multicall3Caller
}

type Multicall3TransactorRaw struct {
	Contract *Multicall3Transactor
}

func NewMulticall3(address common.Address, backend bind.ContractBackend) (*Multicall3, error) {
	abi, err := abi.JSON(strings.NewReader(Multicall3ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMulticall3(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Multicall3{address: address, abi: abi, Multicall3Caller: Multicall3Caller{contract: contract}, Multicall3Transactor: Multicall3Transactor{contract: contract}, Multicall3Filterer: Multicall3Filterer{contract: contract}}, nil
}

func NewMulticall3Caller(address common.Address, caller bind.ContractCaller) (*Multicall3Caller, error) {
	contract, err := bindMulticall3(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Multicall3Caller{contract: contract}, nil
}

func NewMulticall3Transactor(address common.Address, transactor bind.ContractTransactor) (*Multicall3Transactor, error) {
	contract, err := bindMulticall3(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Multicall3Transactor{contract: contract}, nil
}

func NewMulticall3Filterer(address common.Address, filterer bind.ContractFilterer) (*Multicall3Filterer, error) {
	contract, err := bindMulticall3(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Multicall3Filterer{contract: contract}, nil
}

func bindMulticall3(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := Multicall3MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_Multicall3 *Multicall3Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Multicall3.Contract.Multicall3Caller.contract.Call(opts, result, method, params...)
}

func (_Multicall3 *Multicall3Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Multicall3.Contract.Multicall3Transactor.contract.Transfer(opts)
}

func (_Multicall3 *Multicall3Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Multicall3.Contract.Multicall3Transactor.contract.Transact(opts, method, params...)
}

func (_Multicall3 *Multicall3CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Multicall3.Contract.contract.Call(opts, result, method, params...)
}

func (_Multicall3 *Multicall3TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Multicall3.Contract.contract.Transfer(opts)
}

func (_Multicall3 *Multicall3TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Multicall3.Contract.contract.Transact(opts, method, params...)
}

func (_Multicall3 *Multicall3Caller) GetBasefee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Multicall3.contract.Call(opts, &out, "getBasefee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_Multicall3 *Multicall3Session) GetBasefee() (*big.Int, error) {
	return _Multicall3.Contract.GetBasefee(&_Multicall3.CallOpts)
}

func (_Multicall3 *Multicall3CallerSession) GetBasefee() (*big.Int, error) {
	return _Multicall3.Contract.GetBasefee(&_Multicall3.CallOpts)
}

func (_Multicall3 *Multicall3Caller) GetBlockHash(opts *bind.CallOpts, blockNumber *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Multicall3.contract.Call(opts, &out, "getBlockHash", blockNumber)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_Multicall3 *Multicall3Session) GetBlockHash(blockNumber *big.Int) ([32]byte, error) {
	return _Multicall3.Contract.GetBlockHash(&_Multicall3.CallOpts, blockNumber)
}

func (_Multicall3 *Multicall3CallerSession) GetBlockHash(blockNumber *big.Int) ([32]byte, error) {
	return _Multicall3.Contract.GetBlockHash(&_Multicall3.CallOpts, blockNumber)
}

func (_Multicall3 *Multicall3Caller) GetBlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Multicall3.contract.Call(opts, &out, "getBlockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_Multicall3 *Multicall3Session) GetBlockNumber() (*big.Int, error) {
	return _Multicall3.Contract.GetBlockNumber(&_Multicall3.CallOpts)
}

func (_Multicall3 *Multicall3CallerSession) GetBlockNumber() (*big.Int, error) {
	return _Multicall3.Contract.GetBlockNumber(&_Multicall3.CallOpts)
}

func (_Multicall3 *Multicall3Caller) GetChainId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Multicall3.contract.Call(opts, &out, "getChainId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_Multicall3 *Multicall3Session) GetChainId() (*big.Int, error) {
	return _Multicall3.Contract.GetChainId(&_Multicall3.CallOpts)
}

func (_Multicall3 *Multicall3CallerSession) GetChainId() (*big.Int, error) {
	return _Multicall3.Contract.GetChainId(&_Multicall3.CallOpts)
}

func (_Multicall3 *Multicall3Caller) GetCurrentBlockCoinbase(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Multicall3.contract.Call(opts, &out, "getCurrentBlockCoinbase")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_Multicall3 *Multicall3Session) GetCurrentBlockCoinbase() (common.Address, error) {
	return _Multicall3.Contract.GetCurrentBlockCoinbase(&_Multicall3.CallOpts)
}

func (_Multicall3 *Multicall3CallerSession) GetCurrentBlockCoinbase() (common.Address, error) {
	return _Multicall3.Contract.GetCurrentBlockCoinbase(&_Multicall3.CallOpts)
}

func (_Multicall3 *Multicall3Caller) GetCurrentBlockDifficulty(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Multicall3.contract.Call(opts, &out, "getCurrentBlockDifficulty")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_Multicall3 *Multicall3Session) GetCurrentBlockDifficulty() (*big.Int, error) {
	return _Multicall3.Contract.GetCurrentBlockDifficulty(&_Multicall3.CallOpts)
}

func (_Multicall3 *Multicall3CallerSession) GetCurrentBlockDifficulty() (*big.Int, error) {
	return _Multicall3.Contract.GetCurrentBlockDifficulty(&_Multicall3.CallOpts)
}

func (_Multicall3 *Multicall3Caller) GetCurrentBlockGasLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Multicall3.contract.Call(opts, &out, "getCurrentBlockGasLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_Multicall3 *Multicall3Session) GetCurrentBlockGasLimit() (*big.Int, error) {
	return _Multicall3.Contract.GetCurrentBlockGasLimit(&_Multicall3.CallOpts)
}

func (_Multicall3 *Multicall3CallerSession) GetCurrentBlockGasLimit() (*big.Int, error) {
	return _Multicall3.Contract.GetCurrentBlockGasLimit(&_Multicall3.CallOpts)
}

func (_Multicall3 *Multicall3Caller) GetCurrentBlockTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Multicall3.contract.Call(opts, &out, "getCurrentBlockTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_Multicall3 *Multicall3Session) GetCurrentBlockTimestamp() (*big.Int, error) {
	return _Multicall3.Contract.GetCurrentBlockTimestamp(&_Multicall3.CallOpts)
}

func (_Multicall3 *Multicall3CallerSession) GetCurrentBlockTimestamp() (*big.Int, error) {
	return _Multicall3.Contract.GetCurrentBlockTimestamp(&_Multicall3.CallOpts)
}

func (_Multicall3 *Multicall3Caller) GetEthBalance(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Multicall3.contract.Call(opts, &out, "getEthBalance", addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_Multicall3 *Multicall3Session) GetEthBalance(addr common.Address) (*big.Int, error) {
	return _Multicall3.Contract.GetEthBalance(&_Multicall3.CallOpts, addr)
}

func (_Multicall3 *Multicall3CallerSession) GetEthBalance(addr common.Address) (*big.Int, error) {
	return _Multicall3.Contract.GetEthBalance(&_Multicall3.CallOpts, addr)
}

func (_Multicall3 *Multicall3Caller) GetLastBlockHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Multicall3.contract.Call(opts, &out, "getLastBlockHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_Multicall3 *Multicall3Session) GetLastBlockHash() ([32]byte, error) {
	return _Multicall3.Contract.GetLastBlockHash(&_Multicall3.CallOpts)
}

func (_Multicall3 *Multicall3CallerSession) GetLastBlockHash() ([32]byte, error) {
	return _Multicall3.Contract.GetLastBlockHash(&_Multicall3.CallOpts)
}

func (_Multicall3 *Multicall3Transactor) Aggregate(opts *bind.TransactOpts, calls []Multicall3Call) (*types.Transaction, error) {
	return _Multicall3.contract.Transact(opts, "aggregate", calls)
}

func (_Multicall3 *Multicall3Session) Aggregate(calls []Multicall3Call) (*types.Transaction, error) {
	return _Multicall3.Contract.Aggregate(&_Multicall3.TransactOpts, calls)
}

func (_Multicall3 *Multicall3TransactorSession) Aggregate(calls []Multicall3Call) (*types.Transaction, error) {
	return _Multicall3.Contract.Aggregate(&_Multicall3.TransactOpts, calls)
}

func (_Multicall3 *Multicall3Transactor) Aggregate3(opts *bind.TransactOpts, calls []Multicall3Call3) (*types.Transaction, error) {
	return _Multicall3.contract.Transact(opts, "aggregate3", calls)
}

func (_Multicall3 *Multicall3Session) Aggregate3(calls []Multicall3Call3) (*types.Transaction, error) {
	return _Multicall3.Contract.Aggregate3(&_Multicall3.TransactOpts, calls)
}

func (_Multicall3 *Multicall3TransactorSession) Aggregate3(calls []Multicall3Call3) (*types.Transaction, error) {
	return _Multicall3.Contract.Aggregate3(&_Multicall3.TransactOpts, calls)
}

func (_Multicall3 *Multicall3Transactor) Aggregate3Value(opts *bind.TransactOpts, calls []Multicall3Call3Value) (*types.Transaction, error) {
	return _Multicall3.contract.Transact(opts, "aggregate3Value", calls)
}

func (_Multicall3 *Multicall3Session) Aggregate3Value(calls []Multicall3Call3Value) (*types.Transaction, error) {
	return _Multicall3.Contract.Aggregate3Value(&_Multicall3.TransactOpts, calls)
}

func (_Multicall3 *Multicall3TransactorSession) Aggregate3Value(calls []Multicall3Call3Value) (*types.Transaction, error) {
	return _Multicall3.Contract.Aggregate3Value(&_Multicall3.TransactOpts, calls)
}

func (_Multicall3 *Multicall3Transactor) BlockAndAggregate(opts *bind.TransactOpts, calls []Multicall3Call) (*types.Transaction, error) {
	return _Multicall3.contract.Transact(opts, "blockAndAggregate", calls)
}

func (_Multicall3 *Multicall3Session) BlockAndAggregate(calls []Multicall3Call) (*types.Transaction, error) {
	return _Multicall3.Contract.BlockAndAggregate(&_Multicall3.TransactOpts, calls)
}

func (_Multicall3 *Multicall3TransactorSession) BlockAndAggregate(calls []Multicall3Call) (*types.Transaction, error) {
	return _Multicall3.Contract.BlockAndAggregate(&_Multicall3.TransactOpts, calls)
}

func (_Multicall3 *Multicall3Transactor) TryAggregate(opts *bind.TransactOpts, requireSuccess bool, calls []Multicall3Call) (*types.Transaction, error) {
	return _Multicall3.contract.Transact(opts, "tryAggregate", requireSuccess, calls)
}

func (_Multicall3 *Multicall3Session) TryAggregate(requireSuccess bool, calls []Multicall3Call) (*types.Transaction, error) {
	return _Multicall3.Contract.TryAggregate(&_Multicall3.TransactOpts, requireSuccess, calls)
}

func (_Multicall3 *Multicall3TransactorSession) TryAggregate(requireSuccess bool, calls []Multicall3Call) (*types.Transaction, error) {
	return _Multicall3.Contract.TryAggregate(&_Multicall3.TransactOpts, requireSuccess, calls)
}

func (_Multicall3 *Multicall3Transactor) TryBlockAndAggregate(opts *bind.TransactOpts, requireSuccess bool, calls []Multicall3Call) (*types.Transaction, error) {
	return _Multicall3.contract.Transact(opts, "tryBlockAndAggregate", requireSuccess, calls)
}

func (_Multicall3 *Multicall3Session) TryBlockAndAggregate(requireSuccess bool, calls []Multicall3Call) (*types.Transaction, error) {
	return _Multicall3.Contract.TryBlockAndAggregate(&_Multicall3.TransactOpts, requireSuccess, calls)
}

func (_Multicall3 *Multicall3TransactorSession) TryBlockAndAggregate(requireSuccess bool, calls []Multicall3Call) (*types.Transaction, error) {
	return _Multicall3.Contract.TryBlockAndAggregate(&_Multicall3.TransactOpts, requireSuccess, calls)
}

func (_Multicall3 *Multicall3) Address() common.Address {
	return _Multicall3.address
}

type Multicall3Interface interface {
	GetBasefee(opts *bind.CallOpts) (*big.Int, error)

	GetBlockHash(opts *bind.CallOpts, blockNumber *big.Int) ([32]byte, error)

	GetBlockNumber(opts *bind.CallOpts) (*big.Int, error)

	GetChainId(opts *bind.CallOpts) (*big.Int, error)

	GetCurrentBlockCoinbase(opts *bind.CallOpts) (common.Address, error)

	GetCurrentBlockDifficulty(opts *bind.CallOpts) (*big.Int, error)

	GetCurrentBlockGasLimit(opts *bind.CallOpts) (*big.Int, error)

	GetCurrentBlockTimestamp(opts *bind.CallOpts) (*big.Int, error)

	GetEthBalance(opts *bind.CallOpts, addr common.Address) (*big.Int, error)

	GetLastBlockHash(opts *bind.CallOpts) ([32]byte, error)

	Aggregate(opts *bind.TransactOpts, calls []Multicall3Call) (*types.Transaction, error)

	Aggregate3(opts *bind.TransactOpts, calls []Multicall3Call3) (*types.Transaction, error)

	Aggregate3Value(opts *bind.TransactOpts, calls []Multicall3Call3Value) (*types.Transaction, error)

	BlockAndAggregate(opts *bind.TransactOpts, calls []Multicall3Call) (*types.Transaction, error)

	TryAggregate(opts *bind.TransactOpts, requireSuccess bool, calls []Multicall3Call) (*types.Transaction, error)

	TryBlockAndAggregate(opts *bind.TransactOpts, requireSuccess bool, calls []Multicall3Call) (*types.Transaction, error)

	Address() common.Address
}
