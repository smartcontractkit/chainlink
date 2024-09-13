// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package message_hasher

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

type ClientEVMExtraArgsV1 struct {
	GasLimit *big.Int
}

type ClientEVMExtraArgsV2 struct {
	GasLimit                 *big.Int
	AllowOutOfOrderExecution bool
}

type InternalAny2EVMRampMessage struct {
	Header       InternalRampMessageHeader
	Sender       []byte
	Data         []byte
	Receiver     common.Address
	GasLimit     *big.Int
	TokenAmounts []InternalAny2EVMTokenTransfer
}

type InternalAny2EVMTokenTransfer struct {
	SourcePoolAddress []byte
	DestTokenAddress  common.Address
	DestGasAmount     uint32
	ExtraData         []byte
	Amount            *big.Int
}

type InternalEVM2AnyTokenTransfer struct {
	SourcePoolAddress common.Address
	DestTokenAddress  []byte
	ExtraData         []byte
	Amount            *big.Int
	DestExecData      []byte
}

type InternalRampMessageHeader struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	DestChainSelector   uint64
	SequenceNumber      uint64
	Nonce               uint64
}

var MessageHasherMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"name\":\"decodeEVMExtraArgsV1\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMExtraArgsV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"allowOutOfOrderExecution\",\"type\":\"bool\"}],\"name\":\"decodeEVMExtraArgsV2\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"allowOutOfOrderExecution\",\"type\":\"bool\"}],\"internalType\":\"structClient.EVMExtraArgsV2\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMExtraArgsV1\",\"name\":\"extraArgs\",\"type\":\"tuple\"}],\"name\":\"encodeEVMExtraArgsV1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"allowOutOfOrderExecution\",\"type\":\"bool\"}],\"internalType\":\"structClient.EVMExtraArgsV2\",\"name\":\"extraArgs\",\"type\":\"tuple\"}],\"name\":\"encodeEVMExtraArgsV2\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"leafDomainSeparator\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"implicitMetadataHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"fixedSizeFieldsHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"tokenAmountsHash\",\"type\":\"bytes32\"}],\"name\":\"encodeFinalHashPreimage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"name\":\"encodeFixedSizeFieldsHashPreimage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"any2EVMMessageHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"name\":\"encodeMetadataHashPreimage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourcePoolAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destExecData\",\"type\":\"bytes\"}],\"internalType\":\"structInternal.EVM2AnyTokenTransfer[]\",\"name\":\"tokenAmount\",\"type\":\"tuple[]\"}],\"name\":\"encodeTokenAmountsHashPreimage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"destTokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"destGasAmount\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"name\":\"encodeTokenAmountsHashPreimage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"destTokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"destGasAmount\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.Any2EVMRampMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"name\":\"hash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506110c2806100206000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c8063a91d3aeb11610076578063c63641bd1161005b578063c63641bd146101cc578063c7ca9a1814610223578063e733d2091461023657600080fd5b8063a91d3aeb1461017e578063b17df7141461019157600080fd5b80638d3191a0116100a75780638d3191a01461010d57806399df8d0514610120578063a1e747df1461016b57600080fd5b8063750e2ea8146100c35780638503839d146100ec575b600080fd5b6100d66100d136600461074b565b610249565b6040516100e39190610924565b60405180910390f35b6100ff6100fa366004610b00565b610272565b6040519081526020016100e3565b6100d661011b366004610c0a565b6102e1565b6100d661012e366004610c47565b604080516020810196909652858101949094526060850192909252608084015260a0808401919091528151808403909101815260c0909201905290565b6100d6610179366004610c82565b6102f4565b6100d661018c366004610cea565b610326565b6101bd61019f366004610d6b565b60408051602080820183526000909152815190810190915290815290565b604051905181526020016100e3565b6102066101da366004610d94565b604080518082019091526000808252602082015250604080518082019091529182521515602082015290565b6040805182518152602092830151151592810192909252016100e3565b6100d6610231366004610dc0565b61035e565b6100d6610244366004610e14565b61036f565b60608160405160200161025c9190610e56565b6040516020818303038152906040529050919050565b60006102da837f2425b0b9f9054c76ff151b0a175b18f37a4a4e82013a72e9f15c9caa095ed21f856000015160200151866000015160400151866040516020016102bf9493929190610f43565b6040516020818303038152906040528051906020012061037a565b9392505050565b60608160405160200161025c9190610f80565b60608484848460405160200161030d9493929190610f43565b6040516020818303038152906040529050949350505050565b606086868686868660405160200161034396959493929190611055565b60405160208183030381529060405290509695505050505050565b6060610369826104ad565b92915050565b60606103698261056f565b8151805160608085015190830151608080870151940151604051600095869588956103ec959194909391929160200194855273ffffffffffffffffffffffffffffffffffffffff93909316602085015267ffffffffffffffff9182166040850152606084015216608082015260a00190565b604051602081830303815290604052805190602001208560200151805190602001208660400151805190602001208760a0015160405160200161042f9190610f80565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181528282528051602091820120908301979097528101949094526060840192909252608083015260a082015260c081019190915260e00160405160208183030381529060405280519060200120905092915050565b604051815160248201526020820151151560448201526060907f181dcf1000000000000000000000000000000000000000000000000000000000906064015b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915292915050565b604051815160248201526060907f97a657c900000000000000000000000000000000000000000000000000000000906044016104ec565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160a0810167ffffffffffffffff811182821017156105f8576105f86105a6565b60405290565b60405160c0810167ffffffffffffffff811182821017156105f8576105f86105a6565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610668576106686105a6565b604052919050565b600067ffffffffffffffff82111561068a5761068a6105a6565b5060051b60200190565b803573ffffffffffffffffffffffffffffffffffffffff811681146106b857600080fd5b919050565b600082601f8301126106ce57600080fd5b813567ffffffffffffffff8111156106e8576106e86105a6565b61071960207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601610621565b81815284602083860101111561072e57600080fd5b816020850160208301376000918101602001919091529392505050565b6000602080838503121561075e57600080fd5b823567ffffffffffffffff8082111561077657600080fd5b818501915085601f83011261078a57600080fd5b813561079d61079882610670565b610621565b81815260059190911b830184019084810190888311156107bc57600080fd5b8585015b838110156108b3578035858111156107d757600080fd5b860160a0818c037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001121561080c5760008081fd5b6108146105d5565b61081f898301610694565b8152604080830135888111156108355760008081fd5b6108438e8c838701016106bd565b8b840152506060808401358981111561085c5760008081fd5b61086a8f8d838801016106bd565b83850152506080915081840135818401525060a08301358881111561088f5760008081fd5b61089d8e8c838701016106bd565b91830191909152508452509186019186016107c0565b5098975050505050505050565b6000815180845260005b818110156108e6576020818501810151868301820152016108ca565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815260006102da60208301846108c0565b803567ffffffffffffffff811681146106b857600080fd5b600060a0828403121561096157600080fd5b6109696105d5565b90508135815261097b60208301610937565b602082015261098c60408301610937565b604082015261099d60608301610937565b60608201526109ae60808301610937565b608082015292915050565b600082601f8301126109ca57600080fd5b813560206109da61079883610670565b82815260059290921b840181019181810190868411156109f957600080fd5b8286015b84811015610af557803567ffffffffffffffff80821115610a1e5760008081fd5b818901915060a0807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848d03011215610a575760008081fd5b610a5f6105d5565b8784013583811115610a715760008081fd5b610a7f8d8a838801016106bd565b8252506040610a8f818601610694565b8983015260608086013563ffffffff81168114610aac5760008081fd5b808385015250608091508186013585811115610ac85760008081fd5b610ad68f8c838a01016106bd565b91840191909152509190930135908301525083529183019183016109fd565b509695505050505050565b60008060408385031215610b1357600080fd5b823567ffffffffffffffff80821115610b2b57600080fd5b908401906101408287031215610b4057600080fd5b610b486105fe565b610b52878461094f565b815260a083013582811115610b6657600080fd5b610b72888286016106bd565b60208301525060c083013582811115610b8a57600080fd5b610b96888286016106bd565b604083015250610ba860e08401610694565b6060820152610100830135608082015261012083013582811115610bcb57600080fd5b610bd7888286016109b9565b60a08301525093506020850135915080821115610bf357600080fd5b50610c00858286016106bd565b9150509250929050565b600060208284031215610c1c57600080fd5b813567ffffffffffffffff811115610c3357600080fd5b610c3f848285016109b9565b949350505050565b600080600080600060a08688031215610c5f57600080fd5b505083359560208501359550604085013594606081013594506080013592509050565b60008060008060808587031215610c9857600080fd5b84359350610ca860208601610937565b9250610cb660408601610937565b9150606085013567ffffffffffffffff811115610cd257600080fd5b610cde878288016106bd565b91505092959194509250565b60008060008060008060c08789031215610d0357600080fd5b86359550602087013567ffffffffffffffff811115610d2157600080fd5b610d2d89828a016106bd565b955050610d3c60408801610694565b9350610d4a60608801610937565b925060808701359150610d5f60a08801610937565b90509295509295509295565b600060208284031215610d7d57600080fd5b5035919050565b803580151581146106b857600080fd5b60008060408385031215610da757600080fd5b82359150610db760208401610d84565b90509250929050565b600060408284031215610dd257600080fd5b6040516040810181811067ffffffffffffffff82111715610df557610df56105a6565b60405282358152610e0860208401610d84565b60208201529392505050565b600060208284031215610e2657600080fd5b6040516020810181811067ffffffffffffffff82111715610e4957610e496105a6565b6040529135825250919050565b600060208083018184528085518083526040925060408601915060408160051b87010184880160005b83811015610f35577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0898403018552815160a073ffffffffffffffffffffffffffffffffffffffff825116855288820151818a870152610ee1828701826108c0565b9150508782015185820389870152610ef982826108c0565b915050606080830151818701525060808083015192508582038187015250610f2181836108c0565b968901969450505090860190600101610e7f565b509098975050505050505050565b848152600067ffffffffffffffff808616602084015280851660408401525060806060830152610f7660808301846108c0565b9695505050505050565b600060208083018184528085518083526040925060408601915060408160051b87010184880160005b83811015610f35577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0898403018552815160a08151818652610fed828701826108c0565b91505073ffffffffffffffffffffffffffffffffffffffff89830151168986015263ffffffff8883015116888601526060808301518683038288015261103383826108c0565b6080948501519790940196909652505094870194925090860190600101610fa9565b86815260c06020820152600061106e60c08301886108c0565b73ffffffffffffffffffffffffffffffffffffffff9690961660408301525067ffffffffffffffff9384166060820152608081019290925290911660a0909101529291505056fea164736f6c6343000818000a",
}

var MessageHasherABI = MessageHasherMetaData.ABI

var MessageHasherBin = MessageHasherMetaData.Bin

func DeployMessageHasher(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MessageHasher, error) {
	parsed, err := MessageHasherMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MessageHasherBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MessageHasher{address: address, abi: *parsed, MessageHasherCaller: MessageHasherCaller{contract: contract}, MessageHasherTransactor: MessageHasherTransactor{contract: contract}, MessageHasherFilterer: MessageHasherFilterer{contract: contract}}, nil
}

type MessageHasher struct {
	address common.Address
	abi     abi.ABI
	MessageHasherCaller
	MessageHasherTransactor
	MessageHasherFilterer
}

type MessageHasherCaller struct {
	contract *bind.BoundContract
}

type MessageHasherTransactor struct {
	contract *bind.BoundContract
}

type MessageHasherFilterer struct {
	contract *bind.BoundContract
}

type MessageHasherSession struct {
	Contract     *MessageHasher
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MessageHasherCallerSession struct {
	Contract *MessageHasherCaller
	CallOpts bind.CallOpts
}

type MessageHasherTransactorSession struct {
	Contract     *MessageHasherTransactor
	TransactOpts bind.TransactOpts
}

type MessageHasherRaw struct {
	Contract *MessageHasher
}

type MessageHasherCallerRaw struct {
	Contract *MessageHasherCaller
}

type MessageHasherTransactorRaw struct {
	Contract *MessageHasherTransactor
}

func NewMessageHasher(address common.Address, backend bind.ContractBackend) (*MessageHasher, error) {
	abi, err := abi.JSON(strings.NewReader(MessageHasherABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMessageHasher(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MessageHasher{address: address, abi: abi, MessageHasherCaller: MessageHasherCaller{contract: contract}, MessageHasherTransactor: MessageHasherTransactor{contract: contract}, MessageHasherFilterer: MessageHasherFilterer{contract: contract}}, nil
}

func NewMessageHasherCaller(address common.Address, caller bind.ContractCaller) (*MessageHasherCaller, error) {
	contract, err := bindMessageHasher(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MessageHasherCaller{contract: contract}, nil
}

func NewMessageHasherTransactor(address common.Address, transactor bind.ContractTransactor) (*MessageHasherTransactor, error) {
	contract, err := bindMessageHasher(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MessageHasherTransactor{contract: contract}, nil
}

func NewMessageHasherFilterer(address common.Address, filterer bind.ContractFilterer) (*MessageHasherFilterer, error) {
	contract, err := bindMessageHasher(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MessageHasherFilterer{contract: contract}, nil
}

func bindMessageHasher(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MessageHasherMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MessageHasher *MessageHasherRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MessageHasher.Contract.MessageHasherCaller.contract.Call(opts, result, method, params...)
}

func (_MessageHasher *MessageHasherRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MessageHasher.Contract.MessageHasherTransactor.contract.Transfer(opts)
}

func (_MessageHasher *MessageHasherRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MessageHasher.Contract.MessageHasherTransactor.contract.Transact(opts, method, params...)
}

func (_MessageHasher *MessageHasherCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MessageHasher.Contract.contract.Call(opts, result, method, params...)
}

func (_MessageHasher *MessageHasherTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MessageHasher.Contract.contract.Transfer(opts)
}

func (_MessageHasher *MessageHasherTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MessageHasher.Contract.contract.Transact(opts, method, params...)
}

func (_MessageHasher *MessageHasherCaller) DecodeEVMExtraArgsV1(opts *bind.CallOpts, gasLimit *big.Int) (ClientEVMExtraArgsV1, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "decodeEVMExtraArgsV1", gasLimit)

	if err != nil {
		return *new(ClientEVMExtraArgsV1), err
	}

	out0 := *abi.ConvertType(out[0], new(ClientEVMExtraArgsV1)).(*ClientEVMExtraArgsV1)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) DecodeEVMExtraArgsV1(gasLimit *big.Int) (ClientEVMExtraArgsV1, error) {
	return _MessageHasher.Contract.DecodeEVMExtraArgsV1(&_MessageHasher.CallOpts, gasLimit)
}

func (_MessageHasher *MessageHasherCallerSession) DecodeEVMExtraArgsV1(gasLimit *big.Int) (ClientEVMExtraArgsV1, error) {
	return _MessageHasher.Contract.DecodeEVMExtraArgsV1(&_MessageHasher.CallOpts, gasLimit)
}

func (_MessageHasher *MessageHasherCaller) DecodeEVMExtraArgsV2(opts *bind.CallOpts, gasLimit *big.Int, allowOutOfOrderExecution bool) (ClientEVMExtraArgsV2, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "decodeEVMExtraArgsV2", gasLimit, allowOutOfOrderExecution)

	if err != nil {
		return *new(ClientEVMExtraArgsV2), err
	}

	out0 := *abi.ConvertType(out[0], new(ClientEVMExtraArgsV2)).(*ClientEVMExtraArgsV2)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) DecodeEVMExtraArgsV2(gasLimit *big.Int, allowOutOfOrderExecution bool) (ClientEVMExtraArgsV2, error) {
	return _MessageHasher.Contract.DecodeEVMExtraArgsV2(&_MessageHasher.CallOpts, gasLimit, allowOutOfOrderExecution)
}

func (_MessageHasher *MessageHasherCallerSession) DecodeEVMExtraArgsV2(gasLimit *big.Int, allowOutOfOrderExecution bool) (ClientEVMExtraArgsV2, error) {
	return _MessageHasher.Contract.DecodeEVMExtraArgsV2(&_MessageHasher.CallOpts, gasLimit, allowOutOfOrderExecution)
}

func (_MessageHasher *MessageHasherCaller) EncodeEVMExtraArgsV1(opts *bind.CallOpts, extraArgs ClientEVMExtraArgsV1) ([]byte, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "encodeEVMExtraArgsV1", extraArgs)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) EncodeEVMExtraArgsV1(extraArgs ClientEVMExtraArgsV1) ([]byte, error) {
	return _MessageHasher.Contract.EncodeEVMExtraArgsV1(&_MessageHasher.CallOpts, extraArgs)
}

func (_MessageHasher *MessageHasherCallerSession) EncodeEVMExtraArgsV1(extraArgs ClientEVMExtraArgsV1) ([]byte, error) {
	return _MessageHasher.Contract.EncodeEVMExtraArgsV1(&_MessageHasher.CallOpts, extraArgs)
}

func (_MessageHasher *MessageHasherCaller) EncodeEVMExtraArgsV2(opts *bind.CallOpts, extraArgs ClientEVMExtraArgsV2) ([]byte, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "encodeEVMExtraArgsV2", extraArgs)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) EncodeEVMExtraArgsV2(extraArgs ClientEVMExtraArgsV2) ([]byte, error) {
	return _MessageHasher.Contract.EncodeEVMExtraArgsV2(&_MessageHasher.CallOpts, extraArgs)
}

func (_MessageHasher *MessageHasherCallerSession) EncodeEVMExtraArgsV2(extraArgs ClientEVMExtraArgsV2) ([]byte, error) {
	return _MessageHasher.Contract.EncodeEVMExtraArgsV2(&_MessageHasher.CallOpts, extraArgs)
}

func (_MessageHasher *MessageHasherCaller) EncodeFinalHashPreimage(opts *bind.CallOpts, leafDomainSeparator [32]byte, implicitMetadataHash [32]byte, fixedSizeFieldsHash [32]byte, dataHash [32]byte, tokenAmountsHash [32]byte) ([]byte, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "encodeFinalHashPreimage", leafDomainSeparator, implicitMetadataHash, fixedSizeFieldsHash, dataHash, tokenAmountsHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) EncodeFinalHashPreimage(leafDomainSeparator [32]byte, implicitMetadataHash [32]byte, fixedSizeFieldsHash [32]byte, dataHash [32]byte, tokenAmountsHash [32]byte) ([]byte, error) {
	return _MessageHasher.Contract.EncodeFinalHashPreimage(&_MessageHasher.CallOpts, leafDomainSeparator, implicitMetadataHash, fixedSizeFieldsHash, dataHash, tokenAmountsHash)
}

func (_MessageHasher *MessageHasherCallerSession) EncodeFinalHashPreimage(leafDomainSeparator [32]byte, implicitMetadataHash [32]byte, fixedSizeFieldsHash [32]byte, dataHash [32]byte, tokenAmountsHash [32]byte) ([]byte, error) {
	return _MessageHasher.Contract.EncodeFinalHashPreimage(&_MessageHasher.CallOpts, leafDomainSeparator, implicitMetadataHash, fixedSizeFieldsHash, dataHash, tokenAmountsHash)
}

func (_MessageHasher *MessageHasherCaller) EncodeFixedSizeFieldsHashPreimage(opts *bind.CallOpts, messageId [32]byte, sender []byte, receiver common.Address, sequenceNumber uint64, gasLimit *big.Int, nonce uint64) ([]byte, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "encodeFixedSizeFieldsHashPreimage", messageId, sender, receiver, sequenceNumber, gasLimit, nonce)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) EncodeFixedSizeFieldsHashPreimage(messageId [32]byte, sender []byte, receiver common.Address, sequenceNumber uint64, gasLimit *big.Int, nonce uint64) ([]byte, error) {
	return _MessageHasher.Contract.EncodeFixedSizeFieldsHashPreimage(&_MessageHasher.CallOpts, messageId, sender, receiver, sequenceNumber, gasLimit, nonce)
}

func (_MessageHasher *MessageHasherCallerSession) EncodeFixedSizeFieldsHashPreimage(messageId [32]byte, sender []byte, receiver common.Address, sequenceNumber uint64, gasLimit *big.Int, nonce uint64) ([]byte, error) {
	return _MessageHasher.Contract.EncodeFixedSizeFieldsHashPreimage(&_MessageHasher.CallOpts, messageId, sender, receiver, sequenceNumber, gasLimit, nonce)
}

func (_MessageHasher *MessageHasherCaller) EncodeMetadataHashPreimage(opts *bind.CallOpts, any2EVMMessageHash [32]byte, sourceChainSelector uint64, destChainSelector uint64, onRamp []byte) ([]byte, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "encodeMetadataHashPreimage", any2EVMMessageHash, sourceChainSelector, destChainSelector, onRamp)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) EncodeMetadataHashPreimage(any2EVMMessageHash [32]byte, sourceChainSelector uint64, destChainSelector uint64, onRamp []byte) ([]byte, error) {
	return _MessageHasher.Contract.EncodeMetadataHashPreimage(&_MessageHasher.CallOpts, any2EVMMessageHash, sourceChainSelector, destChainSelector, onRamp)
}

func (_MessageHasher *MessageHasherCallerSession) EncodeMetadataHashPreimage(any2EVMMessageHash [32]byte, sourceChainSelector uint64, destChainSelector uint64, onRamp []byte) ([]byte, error) {
	return _MessageHasher.Contract.EncodeMetadataHashPreimage(&_MessageHasher.CallOpts, any2EVMMessageHash, sourceChainSelector, destChainSelector, onRamp)
}

func (_MessageHasher *MessageHasherCaller) EncodeTokenAmountsHashPreimage(opts *bind.CallOpts, tokenAmount []InternalEVM2AnyTokenTransfer) ([]byte, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "encodeTokenAmountsHashPreimage", tokenAmount)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) EncodeTokenAmountsHashPreimage(tokenAmount []InternalEVM2AnyTokenTransfer) ([]byte, error) {
	return _MessageHasher.Contract.EncodeTokenAmountsHashPreimage(&_MessageHasher.CallOpts, tokenAmount)
}

func (_MessageHasher *MessageHasherCallerSession) EncodeTokenAmountsHashPreimage(tokenAmount []InternalEVM2AnyTokenTransfer) ([]byte, error) {
	return _MessageHasher.Contract.EncodeTokenAmountsHashPreimage(&_MessageHasher.CallOpts, tokenAmount)
}

func (_MessageHasher *MessageHasherCaller) EncodeTokenAmountsHashPreimage0(opts *bind.CallOpts, tokenAmounts []InternalAny2EVMTokenTransfer) ([]byte, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "encodeTokenAmountsHashPreimage0", tokenAmounts)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) EncodeTokenAmountsHashPreimage0(tokenAmounts []InternalAny2EVMTokenTransfer) ([]byte, error) {
	return _MessageHasher.Contract.EncodeTokenAmountsHashPreimage0(&_MessageHasher.CallOpts, tokenAmounts)
}

func (_MessageHasher *MessageHasherCallerSession) EncodeTokenAmountsHashPreimage0(tokenAmounts []InternalAny2EVMTokenTransfer) ([]byte, error) {
	return _MessageHasher.Contract.EncodeTokenAmountsHashPreimage0(&_MessageHasher.CallOpts, tokenAmounts)
}

func (_MessageHasher *MessageHasherCaller) Hash(opts *bind.CallOpts, message InternalAny2EVMRampMessage, onRamp []byte) ([32]byte, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "hash", message, onRamp)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) Hash(message InternalAny2EVMRampMessage, onRamp []byte) ([32]byte, error) {
	return _MessageHasher.Contract.Hash(&_MessageHasher.CallOpts, message, onRamp)
}

func (_MessageHasher *MessageHasherCallerSession) Hash(message InternalAny2EVMRampMessage, onRamp []byte) ([32]byte, error) {
	return _MessageHasher.Contract.Hash(&_MessageHasher.CallOpts, message, onRamp)
}

func (_MessageHasher *MessageHasher) Address() common.Address {
	return _MessageHasher.address
}

type MessageHasherInterface interface {
	DecodeEVMExtraArgsV1(opts *bind.CallOpts, gasLimit *big.Int) (ClientEVMExtraArgsV1, error)

	DecodeEVMExtraArgsV2(opts *bind.CallOpts, gasLimit *big.Int, allowOutOfOrderExecution bool) (ClientEVMExtraArgsV2, error)

	EncodeEVMExtraArgsV1(opts *bind.CallOpts, extraArgs ClientEVMExtraArgsV1) ([]byte, error)

	EncodeEVMExtraArgsV2(opts *bind.CallOpts, extraArgs ClientEVMExtraArgsV2) ([]byte, error)

	EncodeFinalHashPreimage(opts *bind.CallOpts, leafDomainSeparator [32]byte, implicitMetadataHash [32]byte, fixedSizeFieldsHash [32]byte, dataHash [32]byte, tokenAmountsHash [32]byte) ([]byte, error)

	EncodeFixedSizeFieldsHashPreimage(opts *bind.CallOpts, messageId [32]byte, sender []byte, receiver common.Address, sequenceNumber uint64, gasLimit *big.Int, nonce uint64) ([]byte, error)

	EncodeMetadataHashPreimage(opts *bind.CallOpts, any2EVMMessageHash [32]byte, sourceChainSelector uint64, destChainSelector uint64, onRamp []byte) ([]byte, error)

	EncodeTokenAmountsHashPreimage(opts *bind.CallOpts, tokenAmount []InternalEVM2AnyTokenTransfer) ([]byte, error)

	EncodeTokenAmountsHashPreimage0(opts *bind.CallOpts, tokenAmounts []InternalAny2EVMTokenTransfer) ([]byte, error)

	Hash(opts *bind.CallOpts, message InternalAny2EVMRampMessage, onRamp []byte) ([32]byte, error)

	Address() common.Address
}
