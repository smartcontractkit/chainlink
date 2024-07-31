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
	TokenAmounts []InternalRampTokenAmount
}

type InternalRampMessageHeader struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	DestChainSelector   uint64
	SequenceNumber      uint64
	Nonce               uint64
}

type InternalRampTokenAmount struct {
	SourcePoolAddress []byte
	DestTokenAddress  []byte
	ExtraData         []byte
	Amount            *big.Int
}

var MessageHasherMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"name\":\"decodeEVMExtraArgsV1\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMExtraArgsV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"allowOutOfOrderExecution\",\"type\":\"bool\"}],\"name\":\"decodeEVMExtraArgsV2\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"allowOutOfOrderExecution\",\"type\":\"bool\"}],\"internalType\":\"structClient.EVMExtraArgsV2\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMExtraArgsV1\",\"name\":\"extraArgs\",\"type\":\"tuple\"}],\"name\":\"encodeEVMExtraArgsV1\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"allowOutOfOrderExecution\",\"type\":\"bool\"}],\"internalType\":\"structClient.EVMExtraArgsV2\",\"name\":\"extraArgs\",\"type\":\"tuple\"}],\"name\":\"encodeEVMExtraArgsV2\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"leafDomainSeparator\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"implicitMetadataHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"fixedSizeFieldsHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"tokenAmountsHash\",\"type\":\"bytes32\"}],\"name\":\"encodeFinalHashPreimage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"name\":\"encodeFixedSizeFieldsHashPreimage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"any2EVMMessageHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"name\":\"encodeMetadataHashPreimage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.RampTokenAmount[]\",\"name\":\"rampTokenAmounts\",\"type\":\"tuple[]\"}],\"name\":\"encodeTokenAmountsHashPreimage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.RampTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.Any2EVMRampMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"name\":\"hash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610de7806100206000396000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c8063a91d3aeb11610076578063c63641bd1161005b578063c63641bd1461019e578063c7ca9a18146101f5578063e733d2091461020857600080fd5b8063a91d3aeb14610150578063b17df7141461016357600080fd5b8063902e94a0146100a85780639511afaa146100d157806399df8d05146100f2578063a1e747df1461013d575b600080fd5b6100bb6100b63660046107d9565b61021b565b6040516100c8919061087a565b60405180910390f35b6100e46100df366004610958565b610244565b6040519081526020016100c8565b6100bb610100366004610a62565b604080516020810196909652858101949094526060850192909252608084015260a0808401919091528151808403909101815260c0909201905290565b6100bb61014b366004610a9d565b610257565b6100bb61015e366004610b05565b610289565b61018f610171366004610b86565b60408051602080820183526000909152815190810190915290815290565b604051905181526020016100c8565b6101d86101ac366004610baf565b604080518082019091526000808252602082015250604080518082019091529182521515602082015290565b6040805182518152602092830151151592810192909252016100c8565b6100bb610203366004610bdb565b6102c1565b6100bb610216366004610c2f565b6102d2565b60608160405160200161022e9190610c71565b6040516020818303038152906040529050919050565b600061025083836102dd565b9392505050565b6060848484846040516020016102709493929190610d3d565b6040516020818303038152906040529050949350505050565b60608686868686866040516020016102a696959493929190610d7a565b60405160208183030381529060405290509695505050505050565b60606102cc8261043a565b92915050565b60606102cc826104fc565b815160208082015160409283015192516000938493610323937f2425b0b9f9054c76ff151b0a175b18f37a4a4e82013a72e9f15c9caa095ed21f93909291889101610d3d565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815290829052805160209182012086518051888401516060808b0151908401516080808d0151950151959761038a9794969395929491939101610d7a565b604051602081830303815290604052805190602001208560400151805190602001208660a001516040516020016103c19190610c71565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181528282528051602091820120908301969096528101939093526060830191909152608082015260a081019190915260c00160405160208183030381529060405280519060200120905092915050565b604051815160248201526020820151151560448201526060907f181dcf1000000000000000000000000000000000000000000000000000000000906064015b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915292915050565b604051815160248201526060907f97a657c90000000000000000000000000000000000000000000000000000000090604401610479565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516080810167ffffffffffffffff8111828210171561058557610585610533565b60405290565b60405160c0810167ffffffffffffffff8111828210171561058557610585610533565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156105f5576105f5610533565b604052919050565b600082601f83011261060e57600080fd5b813567ffffffffffffffff81111561062857610628610533565b61065960207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016105ae565b81815284602083860101111561066e57600080fd5b816020850160208301376000918101602001919091529392505050565b600082601f83011261069c57600080fd5b8135602067ffffffffffffffff808311156106b9576106b9610533565b8260051b6106c88382016105ae565b93845285810183019383810190888611156106e257600080fd5b84880192505b858310156107cd578235848111156107005760008081fd5b88016080818b037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018113156107365760008081fd5b61073e610562565b87830135878111156107505760008081fd5b61075e8d8a838701016105fd565b825250604080840135888111156107755760008081fd5b6107838e8b838801016105fd565b8a840152506060808501358981111561079c5760008081fd5b6107aa8f8c838901016105fd565b9284019290925293909201359281019290925250825291840191908401906106e8565b98975050505050505050565b6000602082840312156107eb57600080fd5b813567ffffffffffffffff81111561080257600080fd5b61080e8482850161068b565b949350505050565b6000815180845260005b8181101561083c57602081850181015186830182015201610820565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815260006102506020830184610816565b803567ffffffffffffffff811681146108a557600080fd5b919050565b600060a082840312156108bc57600080fd5b60405160a0810181811067ffffffffffffffff821117156108df576108df610533565b604052823581529050806108f56020840161088d565b60208201526109066040840161088d565b60408201526109176060840161088d565b60608201526109286080840161088d565b60808201525092915050565b803573ffffffffffffffffffffffffffffffffffffffff811681146108a557600080fd5b6000806040838503121561096b57600080fd5b823567ffffffffffffffff8082111561098357600080fd5b90840190610140828703121561099857600080fd5b6109a061058b565b6109aa87846108aa565b815260a0830135828111156109be57600080fd5b6109ca888286016105fd565b60208301525060c0830135828111156109e257600080fd5b6109ee888286016105fd565b604083015250610a0060e08401610934565b6060820152610100830135608082015261012083013582811115610a2357600080fd5b610a2f8882860161068b565b60a08301525093506020850135915080821115610a4b57600080fd5b50610a58858286016105fd565b9150509250929050565b600080600080600060a08688031215610a7a57600080fd5b505083359560208501359550604085013594606081013594506080013592509050565b60008060008060808587031215610ab357600080fd5b84359350610ac36020860161088d565b9250610ad16040860161088d565b9150606085013567ffffffffffffffff811115610aed57600080fd5b610af9878288016105fd565b91505092959194509250565b60008060008060008060c08789031215610b1e57600080fd5b86359550602087013567ffffffffffffffff811115610b3c57600080fd5b610b4889828a016105fd565b955050610b5760408801610934565b9350610b656060880161088d565b925060808701359150610b7a60a0880161088d565b90509295509295509295565b600060208284031215610b9857600080fd5b5035919050565b803580151581146108a557600080fd5b60008060408385031215610bc257600080fd5b82359150610bd260208401610b9f565b90509250929050565b600060408284031215610bed57600080fd5b6040516040810181811067ffffffffffffffff82111715610c1057610c10610533565b60405282358152610c2360208401610b9f565b60208201529392505050565b600060208284031215610c4157600080fd5b6040516020810181811067ffffffffffffffff82111715610c6457610c64610533565b6040529135825250919050565b600060208083018184528085518083526040925060408601915060408160051b87010184880160005b83811015610d2f577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0898403018552815160808151818652610cde82870182610816565b915050888201518582038a870152610cf68282610816565b9150508782015185820389870152610d0e8282610816565b60609384015196909301959095525094870194925090860190600101610c9a565b509098975050505050505050565b848152600067ffffffffffffffff808616602084015280851660408401525060806060830152610d706080830184610816565b9695505050505050565b86815260c060208201526000610d9360c0830188610816565b73ffffffffffffffffffffffffffffffffffffffff9690961660408301525067ffffffffffffffff9384166060820152608081019290925290911660a0909101529291505056fea164736f6c6343000818000a",
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

func (_MessageHasher *MessageHasherCaller) EncodeTokenAmountsHashPreimage(opts *bind.CallOpts, rampTokenAmounts []InternalRampTokenAmount) ([]byte, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "encodeTokenAmountsHashPreimage", rampTokenAmounts)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) EncodeTokenAmountsHashPreimage(rampTokenAmounts []InternalRampTokenAmount) ([]byte, error) {
	return _MessageHasher.Contract.EncodeTokenAmountsHashPreimage(&_MessageHasher.CallOpts, rampTokenAmounts)
}

func (_MessageHasher *MessageHasherCallerSession) EncodeTokenAmountsHashPreimage(rampTokenAmounts []InternalRampTokenAmount) ([]byte, error) {
	return _MessageHasher.Contract.EncodeTokenAmountsHashPreimage(&_MessageHasher.CallOpts, rampTokenAmounts)
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

	EncodeTokenAmountsHashPreimage(opts *bind.CallOpts, rampTokenAmounts []InternalRampTokenAmount) ([]byte, error)

	Hash(opts *bind.CallOpts, message InternalAny2EVMRampMessage, onRamp []byte) ([32]byte, error)

	Address() common.Address
}
