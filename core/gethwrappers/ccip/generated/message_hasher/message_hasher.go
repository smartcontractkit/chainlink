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

type ClientEVMTokenAmount struct {
	Token  common.Address
	Amount *big.Int
}

type InternalAny2EVMRampMessage struct {
	Header          InternalRampMessageHeader
	Sender          []byte
	Data            []byte
	Receiver        common.Address
	GasLimit        *big.Int
	TokenAmounts    []ClientEVMTokenAmount
	SourceTokenData [][]byte
}

type InternalRampMessageHeader struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	DestChainSelector   uint64
	SequenceNumber      uint64
	Nonce               uint64
}

var MessageHasherMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"leafDomainSeparator\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"implicitMetadataHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"fixedSizeFieldsHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"tokenAmountsHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"sourceTokenDataHash\",\"type\":\"bytes32\"}],\"name\":\"encodeFinalHashPreimage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"name\":\"encodeFixedSizeFieldsHashPreimage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"any2EVMMessageHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"name\":\"encodeMetadataHashPreimage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"sourceTokenData\",\"type\":\"bytes[]\"}],\"name\":\"encodeSourceTokenDataHashPreimage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"name\":\"encodeTokenAmountsHashPreimage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[]\",\"name\":\"sourceTokenData\",\"type\":\"bytes[]\"}],\"internalType\":\"structInternal.Any2EVMRampMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"name\":\"hash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610bb8806100206000396000f3fe608060405234801561001057600080fd5b50600436106100725760003560e01c8063bf0619ad11610050578063bf0619ad146100d0578063c2cfa165146100e3578063cd34ba82146100f657600080fd5b80633b61b52e14610077578063a1e747df1461009d578063a91d3aeb146100bd575b600080fd5b61008a6100853660046106e3565b610109565b6040519081526020015b60405180910390f35b6100b06100ab366004610812565b61011c565b60405161009491906108de565b6100b06100cb3660046108f1565b61014e565b6100b06100de366004610972565b610186565b6100b06100f13660046109b5565b6101bd565b6100b06101043660046109f2565b6101e6565b600061011583836101f9565b9392505050565b6060848484846040516020016101359493929190610a27565b6040516020818303038152906040529050949350505050565b606086868686868660405160200161016b96959493929190610a64565b60405160208183030381529060405290509695505050505050565b604080516020810188905290810186905260608181018690526080820185905260a0820184905260c082018390529060e00161016b565b6060816040516020016101d09190610ac4565b6040516020818303038152906040529050919050565b6060816040516020016101d09190610b46565b81516020808201516040928301519251600093849361023f937f2425b0b9f9054c76ff151b0a175b18f37a4a4e82013a72e9f15c9caa095ed21f93909291889101610a27565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815290829052805160209182012086518051888401516060808b0151908401516080808d015195015195976102a69794969395929491939101610a64565b604051602081830303815290604052805190602001208560400151805190602001208660a001516040516020016102dd9190610b46565b604051602081830303815290604052805190602001208760c001516040516020016103089190610ac4565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181528282528051602091820120908301979097528101949094526060840192909252608083015260a082015260c081019190915260e00160405160208183030381529060405280519060200120905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff811182821017156103d8576103d8610386565b60405290565b60405160e0810167ffffffffffffffff811182821017156103d8576103d8610386565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561044857610448610386565b604052919050565b803567ffffffffffffffff8116811461046857600080fd5b919050565b600060a0828403121561047f57600080fd5b60405160a0810181811067ffffffffffffffff821117156104a2576104a2610386565b604052823581529050806104b860208401610450565b60208201526104c960408401610450565b60408201526104da60608401610450565b60608201526104eb60808401610450565b60808201525092915050565b600082601f83011261050857600080fd5b813567ffffffffffffffff81111561052257610522610386565b61055360207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601610401565b81815284602083860101111561056857600080fd5b816020850160208301376000918101602001919091529392505050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461046857600080fd5b600067ffffffffffffffff8211156105c3576105c3610386565b5060051b60200190565b600082601f8301126105de57600080fd5b813560206105f36105ee836105a9565b610401565b82815260069290921b8401810191818101908684111561061257600080fd5b8286015b84811015610658576040818903121561062f5760008081fd5b6106376103b5565b61064082610585565b81528185013585820152835291830191604001610616565b509695505050505050565b600082601f83011261067457600080fd5b813560206106846105ee836105a9565b82815260059290921b840181019181810190868411156106a357600080fd5b8286015b8481101561065857803567ffffffffffffffff8111156106c75760008081fd5b6106d58986838b01016104f7565b8452509183019183016106a7565b600080604083850312156106f657600080fd5b823567ffffffffffffffff8082111561070e57600080fd5b90840190610160828703121561072357600080fd5b61072b6103de565b610735878461046d565b815260a08301358281111561074957600080fd5b610755888286016104f7565b60208301525060c08301358281111561076d57600080fd5b610779888286016104f7565b60408301525061078b60e08401610585565b60608201526101008301356080820152610120830135828111156107ae57600080fd5b6107ba888286016105cd565b60a083015250610140830135828111156107d357600080fd5b6107df88828601610663565b60c083015250935060208501359150808211156107fb57600080fd5b50610808858286016104f7565b9150509250929050565b6000806000806080858703121561082857600080fd5b8435935061083860208601610450565b925061084660408601610450565b9150606085013567ffffffffffffffff81111561086257600080fd5b61086e878288016104f7565b91505092959194509250565b6000815180845260005b818110156108a057602081850181015186830182015201610884565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000610115602083018461087a565b60008060008060008060c0878903121561090a57600080fd5b86359550602087013567ffffffffffffffff81111561092857600080fd5b61093489828a016104f7565b95505061094360408801610585565b935061095160608801610450565b92506080870135915061096660a08801610450565b90509295509295509295565b60008060008060008060c0878903121561098b57600080fd5b505084359660208601359650604086013595606081013595506080810135945060a0013592509050565b6000602082840312156109c757600080fd5b813567ffffffffffffffff8111156109de57600080fd5b6109ea84828501610663565b949350505050565b600060208284031215610a0457600080fd5b813567ffffffffffffffff811115610a1b57600080fd5b6109ea848285016105cd565b848152600067ffffffffffffffff808616602084015280851660408401525060806060830152610a5a608083018461087a565b9695505050505050565b86815260c060208201526000610a7d60c083018861087a565b73ffffffffffffffffffffffffffffffffffffffff9690961660408301525067ffffffffffffffff9384166060820152608081019290925290911660a09091015292915050565b600060208083016020845280855180835260408601915060408160051b87010192506020870160005b82811015610b39577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0888603018452610b2785835161087a565b94509285019290850190600101610aed565b5092979650505050505050565b602080825282518282018190526000919060409081850190868401855b82811015610b9e578151805173ffffffffffffffffffffffffffffffffffffffff168552860151868501529284019290850190600101610b63565b509197965050505050505056fea164736f6c6343000818000a",
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

func (_MessageHasher *MessageHasherCaller) EncodeFinalHashPreimage(opts *bind.CallOpts, leafDomainSeparator [32]byte, implicitMetadataHash [32]byte, fixedSizeFieldsHash [32]byte, dataHash [32]byte, tokenAmountsHash [32]byte, sourceTokenDataHash [32]byte) ([]byte, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "encodeFinalHashPreimage", leafDomainSeparator, implicitMetadataHash, fixedSizeFieldsHash, dataHash, tokenAmountsHash, sourceTokenDataHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) EncodeFinalHashPreimage(leafDomainSeparator [32]byte, implicitMetadataHash [32]byte, fixedSizeFieldsHash [32]byte, dataHash [32]byte, tokenAmountsHash [32]byte, sourceTokenDataHash [32]byte) ([]byte, error) {
	return _MessageHasher.Contract.EncodeFinalHashPreimage(&_MessageHasher.CallOpts, leafDomainSeparator, implicitMetadataHash, fixedSizeFieldsHash, dataHash, tokenAmountsHash, sourceTokenDataHash)
}

func (_MessageHasher *MessageHasherCallerSession) EncodeFinalHashPreimage(leafDomainSeparator [32]byte, implicitMetadataHash [32]byte, fixedSizeFieldsHash [32]byte, dataHash [32]byte, tokenAmountsHash [32]byte, sourceTokenDataHash [32]byte) ([]byte, error) {
	return _MessageHasher.Contract.EncodeFinalHashPreimage(&_MessageHasher.CallOpts, leafDomainSeparator, implicitMetadataHash, fixedSizeFieldsHash, dataHash, tokenAmountsHash, sourceTokenDataHash)
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

func (_MessageHasher *MessageHasherCaller) EncodeSourceTokenDataHashPreimage(opts *bind.CallOpts, sourceTokenData [][]byte) ([]byte, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "encodeSourceTokenDataHashPreimage", sourceTokenData)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) EncodeSourceTokenDataHashPreimage(sourceTokenData [][]byte) ([]byte, error) {
	return _MessageHasher.Contract.EncodeSourceTokenDataHashPreimage(&_MessageHasher.CallOpts, sourceTokenData)
}

func (_MessageHasher *MessageHasherCallerSession) EncodeSourceTokenDataHashPreimage(sourceTokenData [][]byte) ([]byte, error) {
	return _MessageHasher.Contract.EncodeSourceTokenDataHashPreimage(&_MessageHasher.CallOpts, sourceTokenData)
}

func (_MessageHasher *MessageHasherCaller) EncodeTokenAmountsHashPreimage(opts *bind.CallOpts, tokenAmounts []ClientEVMTokenAmount) ([]byte, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "encodeTokenAmountsHashPreimage", tokenAmounts)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) EncodeTokenAmountsHashPreimage(tokenAmounts []ClientEVMTokenAmount) ([]byte, error) {
	return _MessageHasher.Contract.EncodeTokenAmountsHashPreimage(&_MessageHasher.CallOpts, tokenAmounts)
}

func (_MessageHasher *MessageHasherCallerSession) EncodeTokenAmountsHashPreimage(tokenAmounts []ClientEVMTokenAmount) ([]byte, error) {
	return _MessageHasher.Contract.EncodeTokenAmountsHashPreimage(&_MessageHasher.CallOpts, tokenAmounts)
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
	EncodeFinalHashPreimage(opts *bind.CallOpts, leafDomainSeparator [32]byte, implicitMetadataHash [32]byte, fixedSizeFieldsHash [32]byte, dataHash [32]byte, tokenAmountsHash [32]byte, sourceTokenDataHash [32]byte) ([]byte, error)

	EncodeFixedSizeFieldsHashPreimage(opts *bind.CallOpts, messageId [32]byte, sender []byte, receiver common.Address, sequenceNumber uint64, gasLimit *big.Int, nonce uint64) ([]byte, error)

	EncodeMetadataHashPreimage(opts *bind.CallOpts, any2EVMMessageHash [32]byte, sourceChainSelector uint64, destChainSelector uint64, onRamp []byte) ([]byte, error)

	EncodeSourceTokenDataHashPreimage(opts *bind.CallOpts, sourceTokenData [][]byte) ([]byte, error)

	EncodeTokenAmountsHashPreimage(opts *bind.CallOpts, tokenAmounts []ClientEVMTokenAmount) ([]byte, error)

	Hash(opts *bind.CallOpts, message InternalAny2EVMRampMessage, onRamp []byte) ([32]byte, error)

	Address() common.Address
}
