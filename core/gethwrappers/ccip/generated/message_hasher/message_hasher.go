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

type ClientSVMExtraArgsV1 struct {
	ComputeUnits             uint32
	AccountIsWritableBitmap  uint64
	AllowOutOfOrderExecution bool
	TokenReceiver            [32]byte
	Accounts                 [][32]byte
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
	ABI: "[{\"type\":\"function\",\"name\":\"decodeEVMExtraArgsV1\",\"inputs\":[{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structClient.EVMExtraArgsV1\",\"components\":[{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"decodeEVMExtraArgsV2\",\"inputs\":[{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"allowOutOfOrderExecution\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structClient.EVMExtraArgsV2\",\"components\":[{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"allowOutOfOrderExecution\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"decodeSVMExtraArgsV1\",\"inputs\":[{\"name\":\"computeUnits\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"accountIsWritableBitmap\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"allowOutOfOrderExecution\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"tokenReceiver\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"accounts\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structClient.SVMExtraArgsV1\",\"components\":[{\"name\":\"computeUnits\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"accountIsWritableBitmap\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"allowOutOfOrderExecution\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"tokenReceiver\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"accounts\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}]}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"encodeAny2EVMTokenAmountsHashPreimage\",\"inputs\":[{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destGasAmount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"encodeEVM2AnyTokenAmountsHashPreimage\",\"inputs\":[{\"name\":\"tokenAmount\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.EVM2AnyTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destTokenAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destExecData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"encodeEVMExtraArgsV1\",\"inputs\":[{\"name\":\"extraArgs\",\"type\":\"tuple\",\"internalType\":\"structClient.EVMExtraArgsV1\",\"components\":[{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"encodeEVMExtraArgsV2\",\"inputs\":[{\"name\":\"extraArgs\",\"type\":\"tuple\",\"internalType\":\"structClient.EVMExtraArgsV2\",\"components\":[{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"allowOutOfOrderExecution\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"encodeFinalHashPreimage\",\"inputs\":[{\"name\":\"leafDomainSeparator\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"metaDataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"fixedSizeFieldsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"senderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"tokenAmountsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"encodeFixedSizeFieldsHashPreimage\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"encodeMetadataHashPreimage\",\"inputs\":[{\"name\":\"any2EVMMessageHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRampHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"encodeSVMExtraArgsV1\",\"inputs\":[{\"name\":\"extraArgs\",\"type\":\"tuple\",\"internalType\":\"structClient.SVMExtraArgsV1\",\"components\":[{\"name\":\"computeUnits\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"accountIsWritableBitmap\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"allowOutOfOrderExecution\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"tokenReceiver\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"accounts\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"hash\",\"inputs\":[{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structInternal.Any2EVMRampMessage\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structInternal.RampMessageHeader\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destGasAmount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"onRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"}]",
	Bin: "0x60808060405234601557611196908161001b8239f35b600080fdfe608080604052600436101561001357600080fd5b60003560e01c9081631914fbd214610b96575080633ec7c37714610acf57806357ac4454146109ed5780638503839d146106ce57806394b6624b14610442578063ae5663d7146103db578063b17df71414610382578063bf0619ad1461030e578063c63641bd14610287578063c7ca9a18146101bd578063e04767b8146101425763e733d209146100a357600080fd5b3461013d5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013d576101396040516100e181610ccd565b6004358152604051907f97a657c900000000000000000000000000000000000000000000000000000000602083015251602482015260248152610125604482610d05565b604051918291602083526020830190610e1e565b0390f35b600080fd5b3461013d5760807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013d5761013961017c610d6e565b67ffffffffffffffff61018d610d57565b8160405193600435602086015216604084015216606082015260643560808201526080815261012560a082610d05565b3461013d5760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013d5761013961025b61012560405161020181610ce9565b600435815261020e610d9a565b60208201526040519283917f181dcf100000000000000000000000000000000000000000000000000000000060208401526024830191909160208060408301948051845201511515910152565b037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08101835282610d05565b3461013d5760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013d576101396102c1610d9a565b600060206040516102d181610ce9565b8281520152604051906102e382610ce9565b6004358252151560208201526040519182918291909160208060408301948051845201511515910152565b3461013d5760c07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013d57610139604051600435602082015260243560408201526044356060820152606435608082015260843560a082015260a43560c082015260c0815261012560e082610d05565b3461013d5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013d57602060043560006040516103c481610ccd565b52806040516103d281610ccd565b52604051908152f35b3461013d5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013d5760043567ffffffffffffffff811161013d5761025b610125610433610139933690600401610f8f565b604051928391602083016110b3565b3461013d5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013d5760043567ffffffffffffffff811161013d573660238201121561013d5780600401359061049d82610da9565b906104ab6040519283610d05565b82825260208201906024829460051b8201019036821161013d5760248101925b8284106105df57858560405190604082019060208084015251809152606082019060608160051b84010193916000905b82821061053757610139856101258189037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08101835282610d05565b909192946020806105d1837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0896001960301865289519073ffffffffffffffffffffffffffffffffffffffff825116815260806105b66105a48685015160a08886015260a0850190610e1e565b60408501518482036040860152610e1e565b92606081015160608401520151906080818403910152610e1e565b9701920192019092916104fb565b833567ffffffffffffffff811161013d57820160a07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc823603011261013d576040519161062b83610cb1565b61063760248301610e7d565b8352604482013567ffffffffffffffff811161013d5761065d9060243691850101610f1a565b6020840152606482013567ffffffffffffffff811161013d576106869060243691850101610f1a565b60408401526084820135606084015260a48201359267ffffffffffffffff841161013d576106be602094936024869536920101610f1a565b60808201528152019301926104cb565b3461013d5760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013d5760043567ffffffffffffffff811161013d577ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc81360301610140811261013d576040519060c082019082821067ffffffffffffffff8311176109be5760a0916040521261013d5760405161077081610cb1565b8260040135815261078360248401610d85565b602082015261079460448401610d85565b60408201526107a560648401610d85565b60608201526107b660848401610d85565b6080820152815260a482013567ffffffffffffffff811161013d576107e19060043691850101610f1a565b6020820190815260c483013567ffffffffffffffff811161013d5761080c9060043691860101610f1a565b6040830190815261081f60e48501610e7d565b60608401908152608084019461010481013586526101248101359067ffffffffffffffff821161013d5760046108589236920101610f8f565b9060a085019182526024359567ffffffffffffffff871161013d576109526108866020983690600401610f1a565b87519067ffffffffffffffff6040818c8501511693015116908a8151910120604051918b8301937f2425b0b9f9054c76ff151b0a175b18f37a4a4e82013a72e9f15c9caa095ed21f8552604084015260608301526080820152608081526108ee60a082610d05565b5190209651805193516060808301519451608093840151604080518e8101998a5273ffffffffffffffffffffffffffffffffffffffff9590951660208a015267ffffffffffffffff9788169089015291870152909316908401528160a0840161025b565b51902092518581519101209151858151910120905160405161097b8161025b89820194856110b3565b5190209160405193868501956000875260408601526060850152608084015260a083015260c082015260c081526109b360e082610d05565b519020604051908152f35b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b3461013d5760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013d5760043563ffffffff811680910361013d57610a35610d6e565b906044359182151580930361013d576084359267ffffffffffffffff841161013d5767ffffffffffffffff610a71610139953690600401610dc1565b9260606080604051610a8281610cb1565b600081526000602082015260006040820152600083820152015260405194610aa986610cb1565b855216602084015260408301526064356060830152608082015260405191829182610e9e565b3461013d5760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013d5760243573ffffffffffffffffffffffffffffffffffffffff8116810361013d57610b26610d57565b9060843567ffffffffffffffff8116810361013d5760408051600435602082015273ffffffffffffffffffffffffffffffffffffffff9093169083015267ffffffffffffffff928316606083015260643560808301529190911660a0820152610139906101258160c0810161025b565b3461013d5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261013d5760043567ffffffffffffffff811161013d5760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc823603011261013d57610c0c82610cb1565b610c1881600401610d46565b8252610c2660248201610d85565b60208301526044810135801515810361013d5760408301526064810135606083015260848101359067ffffffffffffffff821161013d5761013992610c7761012592600461025b9536920101610dc1565b60808201526040519283917f1f3b3aba00000000000000000000000000000000000000000000000000000000602084015260248301610e9e565b60a0810190811067ffffffffffffffff8211176109be57604052565b6020810190811067ffffffffffffffff8211176109be57604052565b6040810190811067ffffffffffffffff8211176109be57604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff8211176109be57604052565b359063ffffffff8216820361013d57565b6044359067ffffffffffffffff8216820361013d57565b6024359067ffffffffffffffff8216820361013d57565b359067ffffffffffffffff8216820361013d57565b60243590811515820361013d57565b67ffffffffffffffff81116109be5760051b60200190565b9080601f8301121561013d578135610dd881610da9565b92610de66040519485610d05565b81845260208085019260051b82010192831161013d57602001905b828210610e0e5750505090565b8135815260209182019101610e01565b919082519283825260005b848110610e685750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8460006020809697860101520116010190565b80602080928401015182828601015201610e29565b359073ffffffffffffffffffffffffffffffffffffffff8216820361013d57565b9060208252602060e0608060c085019363ffffffff8151168487015267ffffffffffffffff8482015116604087015260408101511515606087015260608101518287015201519360a08082015284518094520192019060005b818110610f045750505090565b8251845260209384019390920191600101610ef7565b81601f8201121561013d5780359067ffffffffffffffff82116109be5760405192610f6d60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8601160185610d05565b8284526020838301011161013d57816000926020809301838601378301015290565b81601f8201121561013d57803590610fa682610da9565b92610fb46040519485610d05565b82845260208085019360051b8301019181831161013d5760208101935b838510610fe057505050505090565b843567ffffffffffffffff811161013d57820160a07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0828603011261013d576040519161102c83610cb1565b602082013567ffffffffffffffff811161013d5785602061104f92850101610f1a565b835261105d60408301610e7d565b602084015261106e60608301610d46565b604084015260808201359267ffffffffffffffff841161013d5760a08361109c886020809881980101610f1a565b606084015201356080820152815201940193610fd1565b602081016020825282518091526040820191602060408360051b8301019401926000915b8383106110e657505050505090565b9091929394602080827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08560019503018652885190608080611174611134855160a0865260a0860190610e1e565b73ffffffffffffffffffffffffffffffffffffffff87870151168786015263ffffffff604087015116604086015260608601518582036060870152610e1e565b930151910152970193019301919392906110d756fea164736f6c634300081a000a",
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

func (_MessageHasher *MessageHasherCaller) DecodeSVMExtraArgsV1(opts *bind.CallOpts, computeUnits uint32, accountIsWritableBitmap uint64, allowOutOfOrderExecution bool, tokenReceiver [32]byte, accounts [][32]byte) (ClientSVMExtraArgsV1, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "decodeSVMExtraArgsV1", computeUnits, accountIsWritableBitmap, allowOutOfOrderExecution, tokenReceiver, accounts)

	if err != nil {
		return *new(ClientSVMExtraArgsV1), err
	}

	out0 := *abi.ConvertType(out[0], new(ClientSVMExtraArgsV1)).(*ClientSVMExtraArgsV1)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) DecodeSVMExtraArgsV1(computeUnits uint32, accountIsWritableBitmap uint64, allowOutOfOrderExecution bool, tokenReceiver [32]byte, accounts [][32]byte) (ClientSVMExtraArgsV1, error) {
	return _MessageHasher.Contract.DecodeSVMExtraArgsV1(&_MessageHasher.CallOpts, computeUnits, accountIsWritableBitmap, allowOutOfOrderExecution, tokenReceiver, accounts)
}

func (_MessageHasher *MessageHasherCallerSession) DecodeSVMExtraArgsV1(computeUnits uint32, accountIsWritableBitmap uint64, allowOutOfOrderExecution bool, tokenReceiver [32]byte, accounts [][32]byte) (ClientSVMExtraArgsV1, error) {
	return _MessageHasher.Contract.DecodeSVMExtraArgsV1(&_MessageHasher.CallOpts, computeUnits, accountIsWritableBitmap, allowOutOfOrderExecution, tokenReceiver, accounts)
}

func (_MessageHasher *MessageHasherCaller) EncodeAny2EVMTokenAmountsHashPreimage(opts *bind.CallOpts, tokenAmounts []InternalAny2EVMTokenTransfer) ([]byte, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "encodeAny2EVMTokenAmountsHashPreimage", tokenAmounts)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) EncodeAny2EVMTokenAmountsHashPreimage(tokenAmounts []InternalAny2EVMTokenTransfer) ([]byte, error) {
	return _MessageHasher.Contract.EncodeAny2EVMTokenAmountsHashPreimage(&_MessageHasher.CallOpts, tokenAmounts)
}

func (_MessageHasher *MessageHasherCallerSession) EncodeAny2EVMTokenAmountsHashPreimage(tokenAmounts []InternalAny2EVMTokenTransfer) ([]byte, error) {
	return _MessageHasher.Contract.EncodeAny2EVMTokenAmountsHashPreimage(&_MessageHasher.CallOpts, tokenAmounts)
}

func (_MessageHasher *MessageHasherCaller) EncodeEVM2AnyTokenAmountsHashPreimage(opts *bind.CallOpts, tokenAmount []InternalEVM2AnyTokenTransfer) ([]byte, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "encodeEVM2AnyTokenAmountsHashPreimage", tokenAmount)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) EncodeEVM2AnyTokenAmountsHashPreimage(tokenAmount []InternalEVM2AnyTokenTransfer) ([]byte, error) {
	return _MessageHasher.Contract.EncodeEVM2AnyTokenAmountsHashPreimage(&_MessageHasher.CallOpts, tokenAmount)
}

func (_MessageHasher *MessageHasherCallerSession) EncodeEVM2AnyTokenAmountsHashPreimage(tokenAmount []InternalEVM2AnyTokenTransfer) ([]byte, error) {
	return _MessageHasher.Contract.EncodeEVM2AnyTokenAmountsHashPreimage(&_MessageHasher.CallOpts, tokenAmount)
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

func (_MessageHasher *MessageHasherCaller) EncodeFinalHashPreimage(opts *bind.CallOpts, leafDomainSeparator [32]byte, metaDataHash [32]byte, fixedSizeFieldsHash [32]byte, senderHash [32]byte, dataHash [32]byte, tokenAmountsHash [32]byte) ([]byte, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "encodeFinalHashPreimage", leafDomainSeparator, metaDataHash, fixedSizeFieldsHash, senderHash, dataHash, tokenAmountsHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) EncodeFinalHashPreimage(leafDomainSeparator [32]byte, metaDataHash [32]byte, fixedSizeFieldsHash [32]byte, senderHash [32]byte, dataHash [32]byte, tokenAmountsHash [32]byte) ([]byte, error) {
	return _MessageHasher.Contract.EncodeFinalHashPreimage(&_MessageHasher.CallOpts, leafDomainSeparator, metaDataHash, fixedSizeFieldsHash, senderHash, dataHash, tokenAmountsHash)
}

func (_MessageHasher *MessageHasherCallerSession) EncodeFinalHashPreimage(leafDomainSeparator [32]byte, metaDataHash [32]byte, fixedSizeFieldsHash [32]byte, senderHash [32]byte, dataHash [32]byte, tokenAmountsHash [32]byte) ([]byte, error) {
	return _MessageHasher.Contract.EncodeFinalHashPreimage(&_MessageHasher.CallOpts, leafDomainSeparator, metaDataHash, fixedSizeFieldsHash, senderHash, dataHash, tokenAmountsHash)
}

func (_MessageHasher *MessageHasherCaller) EncodeFixedSizeFieldsHashPreimage(opts *bind.CallOpts, messageId [32]byte, receiver common.Address, sequenceNumber uint64, gasLimit *big.Int, nonce uint64) ([]byte, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "encodeFixedSizeFieldsHashPreimage", messageId, receiver, sequenceNumber, gasLimit, nonce)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) EncodeFixedSizeFieldsHashPreimage(messageId [32]byte, receiver common.Address, sequenceNumber uint64, gasLimit *big.Int, nonce uint64) ([]byte, error) {
	return _MessageHasher.Contract.EncodeFixedSizeFieldsHashPreimage(&_MessageHasher.CallOpts, messageId, receiver, sequenceNumber, gasLimit, nonce)
}

func (_MessageHasher *MessageHasherCallerSession) EncodeFixedSizeFieldsHashPreimage(messageId [32]byte, receiver common.Address, sequenceNumber uint64, gasLimit *big.Int, nonce uint64) ([]byte, error) {
	return _MessageHasher.Contract.EncodeFixedSizeFieldsHashPreimage(&_MessageHasher.CallOpts, messageId, receiver, sequenceNumber, gasLimit, nonce)
}

func (_MessageHasher *MessageHasherCaller) EncodeMetadataHashPreimage(opts *bind.CallOpts, any2EVMMessageHash [32]byte, sourceChainSelector uint64, destChainSelector uint64, onRampHash [32]byte) ([]byte, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "encodeMetadataHashPreimage", any2EVMMessageHash, sourceChainSelector, destChainSelector, onRampHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) EncodeMetadataHashPreimage(any2EVMMessageHash [32]byte, sourceChainSelector uint64, destChainSelector uint64, onRampHash [32]byte) ([]byte, error) {
	return _MessageHasher.Contract.EncodeMetadataHashPreimage(&_MessageHasher.CallOpts, any2EVMMessageHash, sourceChainSelector, destChainSelector, onRampHash)
}

func (_MessageHasher *MessageHasherCallerSession) EncodeMetadataHashPreimage(any2EVMMessageHash [32]byte, sourceChainSelector uint64, destChainSelector uint64, onRampHash [32]byte) ([]byte, error) {
	return _MessageHasher.Contract.EncodeMetadataHashPreimage(&_MessageHasher.CallOpts, any2EVMMessageHash, sourceChainSelector, destChainSelector, onRampHash)
}

func (_MessageHasher *MessageHasherCaller) EncodeSVMExtraArgsV1(opts *bind.CallOpts, extraArgs ClientSVMExtraArgsV1) ([]byte, error) {
	var out []interface{}
	err := _MessageHasher.contract.Call(opts, &out, "encodeSVMExtraArgsV1", extraArgs)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_MessageHasher *MessageHasherSession) EncodeSVMExtraArgsV1(extraArgs ClientSVMExtraArgsV1) ([]byte, error) {
	return _MessageHasher.Contract.EncodeSVMExtraArgsV1(&_MessageHasher.CallOpts, extraArgs)
}

func (_MessageHasher *MessageHasherCallerSession) EncodeSVMExtraArgsV1(extraArgs ClientSVMExtraArgsV1) ([]byte, error) {
	return _MessageHasher.Contract.EncodeSVMExtraArgsV1(&_MessageHasher.CallOpts, extraArgs)
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

	DecodeSVMExtraArgsV1(opts *bind.CallOpts, computeUnits uint32, accountIsWritableBitmap uint64, allowOutOfOrderExecution bool, tokenReceiver [32]byte, accounts [][32]byte) (ClientSVMExtraArgsV1, error)

	EncodeAny2EVMTokenAmountsHashPreimage(opts *bind.CallOpts, tokenAmounts []InternalAny2EVMTokenTransfer) ([]byte, error)

	EncodeEVM2AnyTokenAmountsHashPreimage(opts *bind.CallOpts, tokenAmount []InternalEVM2AnyTokenTransfer) ([]byte, error)

	EncodeEVMExtraArgsV1(opts *bind.CallOpts, extraArgs ClientEVMExtraArgsV1) ([]byte, error)

	EncodeEVMExtraArgsV2(opts *bind.CallOpts, extraArgs ClientEVMExtraArgsV2) ([]byte, error)

	EncodeFinalHashPreimage(opts *bind.CallOpts, leafDomainSeparator [32]byte, metaDataHash [32]byte, fixedSizeFieldsHash [32]byte, senderHash [32]byte, dataHash [32]byte, tokenAmountsHash [32]byte) ([]byte, error)

	EncodeFixedSizeFieldsHashPreimage(opts *bind.CallOpts, messageId [32]byte, receiver common.Address, sequenceNumber uint64, gasLimit *big.Int, nonce uint64) ([]byte, error)

	EncodeMetadataHashPreimage(opts *bind.CallOpts, any2EVMMessageHash [32]byte, sourceChainSelector uint64, destChainSelector uint64, onRampHash [32]byte) ([]byte, error)

	EncodeSVMExtraArgsV1(opts *bind.CallOpts, extraArgs ClientSVMExtraArgsV1) ([]byte, error)

	Hash(opts *bind.CallOpts, message InternalAny2EVMRampMessage, onRamp []byte) ([32]byte, error)

	Address() common.Address
}
