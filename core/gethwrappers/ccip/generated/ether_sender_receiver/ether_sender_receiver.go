// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ether_sender_receiver

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

type ClientAny2EVMMessage struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	Sender              []byte
	Data                []byte
	DestTokenAmounts    []ClientEVMTokenAmount
}

type ClientEVM2AnyMessage struct {
	Receiver     []byte
	Data         []byte
	TokenAmounts []ClientEVMTokenAmount
	FeeToken     common.Address
	ExtraArgs    []byte
}

type ClientEVMTokenAmount struct {
	Token  common.Address
	Amount *big.Int
}

var EtherSenderReceiverMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"ccipReceive\",\"inputs\":[{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structClient.Any2EVMMessage\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ccipSend\",\"inputs\":[{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structClient.EVM2AnyMessage\",\"components\":[{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"feeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"getFee\",\"inputs\":[{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structClient.EVM2AnyMessage\",\"components\":[{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"feeToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"extraArgs\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"fee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRouter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"i_weth\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIWrappedNative\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"error\",\"name\":\"InvalidRouter\",\"inputs\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidToken\",\"inputs\":[{\"name\":\"gotToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"expectedToken\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidTokenAmounts\",\"inputs\":[{\"name\":\"gotAmounts\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TokenAmountNotEqualToMsgValue\",\"inputs\":[{\"name\":\"gotAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"msgValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]",
	Bin: "0x60c080604052346101245761002a9061161380380380916100208285610180565b83398101906101b9565b6001600160a01b03811690811561016a5760805260405163e861e90760e01b815290602082600481845afa908115610131576044602092600094859161013d575b506001600160a01b031660a081905260405163095ea7b360e01b8152600481019390935284196024840152919384928391905af18015610131576100f4575b60405161143a90816101d9823960805181818160bd01528181610185015281816106170152610d44015260a0518181816102cf0152818161058601528181610ccc01526112d40152f35b6020813d602011610129575b8161010d60209383610180565b8101031261012457518015150361012457386100aa565b600080fd5b3d9150610100565b6040513d6000823e3d90fd5b61015d9150843d8611610163575b6101558183610180565b8101906101b9565b3861006b565b503d61014b565b6335fdcccd60e21b600052600060045260246000fd5b601f909101601f19168101906001600160401b038211908210176101a357604052565b634e487b7160e01b600052604160045260246000fd5b9081602091031261012457516001600160a01b0381168103610124579056fe608080604052600436101561001d575b50361561001b57600080fd5b005b600090813560e01c90816301ffc9a71461070b57508063181f5a771461068857806320487ded146105aa5780634dbe7e921461053b57806385572ffb1461010057806396f4e9f9146100e45763b0f479a10361000f57346100e157807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100e157602060405173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b80fd5b60206100f86100f2366108f2565b90610c23565b604051908152f35b50346100e15760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100e15760043567ffffffffffffffff81116105375760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc82360301126105375773ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016330361050b576040516101b7816107cd565b81600401358152602482013567ffffffffffffffff81168103610468576020820152604482013567ffffffffffffffff8111610468576101fd9060043691850101610ac0565b6040820152606482013567ffffffffffffffff8111610468576102269060043691850101610ac0565b9160608201928352608481013567ffffffffffffffff81116105075760809160046102549236920101610b28565b9101918183525160208180518101031261046857602001519073ffffffffffffffffffffffffffffffffffffffff82168092036104685751600181036104dc575073ffffffffffffffffffffffffffffffffffffffff6102b48351610bcf565b5151169173ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001680930361048b5761030260209151610bcf565b51015190823b15610468576040517f2e1a7d4d000000000000000000000000000000000000000000000000000000008152826004820152848160248183885af180156104805761046c575b508380808085855af161035e61132d565b5015610368578380f35b823b1561046857836040517fd0e30db0000000000000000000000000000000000000000000000000000000008152818160048187895af180156104445761044f575b506040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff929092166004830152602482019290925291602091839160449183915af1801561044457610415575b80808380f35b6104369060203d60201161043d575b61042e8183610818565b810190610c0b565b503861040f565b503d610424565b6040513d84823e3d90fd5b8161045c91949394610818565b610468579083386103aa565b8380fd5b8461047991959295610818565b923861034d565b6040513d87823e3d90fd5b838373ffffffffffffffffffffffffffffffffffffffff6104ae60449451610bcf565b5151167f0fc746a1000000000000000000000000000000000000000000000000000000008352600452602452fd5b7f83b9f0ae000000000000000000000000000000000000000000000000000000008452600452602483fd5b8480fd5b6024827fd7f7333400000000000000000000000000000000000000000000000000000000815233600452fd5b5080fd5b50346100e157807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100e157602060405173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b50346100e1576105fe60206105c96105c1366108f2565b9190916111ae565b9060405193849283927f20487ded0000000000000000000000000000000000000000000000000000000084526004840161097c565b038173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa908115610444578291610652575b602082604051908152f35b90506020813d602011610680575b8161066d60209383610818565b8101031261053757602091505138610647565b3d9150610660565b50346100e157807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100e157506107076040516106c9604082610818565b601981527f457468657253656e646572526563656976657220312e352e30000000000000006020820152604051918291602083526020830190610893565b0390f35b9050346105375760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610537576004357fffffffff0000000000000000000000000000000000000000000000000000000081168091036107c957602092507f85572ffb00000000000000000000000000000000000000000000000000000000811490811561079f575b5015158152f35b7f01ffc9a70000000000000000000000000000000000000000000000000000000091501438610798565b8280fd5b60a0810190811067ffffffffffffffff8211176107e957604052565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff8211176107e957604052565b67ffffffffffffffff81116107e957601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b919082519283825260005b8481106108dd5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8460006020809697860101520116010190565b8060208092840101518282860101520161089e565b9060407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8301126109775760043567ffffffffffffffff8116810361097757916024359067ffffffffffffffff8211610977577ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8260a0920301126109775760040190565b600080fd5b9067ffffffffffffffff90939293168152604060208201526109e16109ad845160a0604085015260e0840190610893565b60208501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0848303016060850152610893565b906040840151917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08282030160808301526020808451928381520193019060005b818110610a885750505060808473ffffffffffffffffffffffffffffffffffffffff6060610a85969701511660a084015201519060c07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc082850301910152610893565b90565b8251805173ffffffffffffffffffffffffffffffffffffffff1686526020908101518187015260409095019490920191600101610a22565b81601f8201121561097757803590610ad782610859565b92610ae56040519485610818565b8284526020838301011161097757816000926020809301838601378301015290565b359073ffffffffffffffffffffffffffffffffffffffff8216820361097757565b81601f820112156109775780359067ffffffffffffffff82116107e95760405192610b5960208460051b0185610818565b82845260208085019360061b8301019181831161097757602001925b828410610b83575050505090565b6040848303126109775760405190604082019082821067ffffffffffffffff8311176107e9576040926020928452610bba87610b07565b81528287013583820152815201930192610b75565b805115610bdc5760200190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b90816020910312610977575180151581036109775790565b9060009060408101357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1823603018112156107c9578101803567ffffffffffffffff8111610468578060061b3603602083011361046857156111815760400135606082013573ffffffffffffffffffffffffffffffffffffffff81168091036104685761114b575b50610cb5906111ae565b73ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166020610cfb6040840151610bcf565b510151813b156104685783600491604051928380927fd0e30db0000000000000000000000000000000000000000000000000000000008252865af180156111405761112c575b507f00000000000000000000000000000000000000000000000000000000000000009373ffffffffffffffffffffffffffffffffffffffff851691604051907f20487ded00000000000000000000000000000000000000000000000000000000825260208280610db588876004840161097c565b0381875afa9182156111215786926110ed575b50606085019673ffffffffffffffffffffffffffffffffffffffff8851169788610e805750505050610e319394509060209291476040518096819582947f96f4e9f90000000000000000000000000000000000000000000000000000000084526004840161097c565b03925af1918215610e745791610e45575090565b90506020813d602011610e6c575b81610e6060209383610818565b81010312610977575190565b3d9150610e53565b604051903d90823e3d90fd5b610f20604051602081019a7f23b872dd000000000000000000000000000000000000000000000000000000008c5233602483015230604483015286606483015260648252610ecf608483610818565b8a8060409d8e94610ee286519687610818565b602086527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65646020870152519082855af1610f1a61132d565b9161135d565b80518061104c575b50505173ffffffffffffffffffffffffffffffffffffffff16918203610fa9575b505050916020918493610f8a9587518097819582947f96f4e9f90000000000000000000000000000000000000000000000000000000084526004840161097c565b03925af1928315610f9f575091610e45575090565b51903d90823e3d90fd5b87517f095ea7b300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91909116600482015260248101929092526020908290604490829089905af18015611042579160209391610f8a969593611025575b8193959650829450610f49565b61103b90853d871161043d5761042e8183610818565b5038611018565b86513d87823e3d90fd5b9060208061105e938301019101610c0b565b1561106a573880610f28565b608489517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152fd5b9091506020813d602011611119575b8161110960209383610818565b8101031261097757519038610dc8565b3d91506110fc565b6040513d88823e3d90fd5b8361113991949294610818565b9138610d41565b6040513d86823e3d90fd5b348114610cab577fba2f746700000000000000000000000000000000000000000000000000000000835260045234602452604482fd5b6024837f4e487b710000000000000000000000000000000000000000000000000000000081526032600452fd5b606060806040516111be816107cd565b828152826020820152826040820152600083820152015260a08136031261097757604051906111ec826107cd565b803567ffffffffffffffff81116109775761120a9036908301610ac0565b8252602081013567ffffffffffffffff81116109775761122d9036908301610ac0565b60208301908152604082013567ffffffffffffffff8111610977576112559036908401610b28565b916040840192835261126960608201610b07565b606085015260808101359067ffffffffffffffff82116109775761128f91369101610ac0565b608084015281515160018103611300575060405190336020830152602082526112b9604083610818565b526112fb73ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169151610bcf565b515290565b7f83b9f0ae0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b3d15611358573d9061133e82610859565b9161134c6040519384610818565b82523d6000602084013e565b606090565b919290156113d85750815115611371575090565b3b1561137a5790565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152fd5b8251909150156113eb5750805190602001fd5b611429906040519182917f08c379a0000000000000000000000000000000000000000000000000000000008352602060048401526024830190610893565b0390fdfea164736f6c634300081a000a",
}

var EtherSenderReceiverABI = EtherSenderReceiverMetaData.ABI

var EtherSenderReceiverBin = EtherSenderReceiverMetaData.Bin

func DeployEtherSenderReceiver(auth *bind.TransactOpts, backend bind.ContractBackend, router common.Address) (common.Address, *types.Transaction, *EtherSenderReceiver, error) {
	parsed, err := EtherSenderReceiverMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EtherSenderReceiverBin), backend, router)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EtherSenderReceiver{address: address, abi: *parsed, EtherSenderReceiverCaller: EtherSenderReceiverCaller{contract: contract}, EtherSenderReceiverTransactor: EtherSenderReceiverTransactor{contract: contract}, EtherSenderReceiverFilterer: EtherSenderReceiverFilterer{contract: contract}}, nil
}

type EtherSenderReceiver struct {
	address common.Address
	abi     abi.ABI
	EtherSenderReceiverCaller
	EtherSenderReceiverTransactor
	EtherSenderReceiverFilterer
}

type EtherSenderReceiverCaller struct {
	contract *bind.BoundContract
}

type EtherSenderReceiverTransactor struct {
	contract *bind.BoundContract
}

type EtherSenderReceiverFilterer struct {
	contract *bind.BoundContract
}

type EtherSenderReceiverSession struct {
	Contract     *EtherSenderReceiver
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type EtherSenderReceiverCallerSession struct {
	Contract *EtherSenderReceiverCaller
	CallOpts bind.CallOpts
}

type EtherSenderReceiverTransactorSession struct {
	Contract     *EtherSenderReceiverTransactor
	TransactOpts bind.TransactOpts
}

type EtherSenderReceiverRaw struct {
	Contract *EtherSenderReceiver
}

type EtherSenderReceiverCallerRaw struct {
	Contract *EtherSenderReceiverCaller
}

type EtherSenderReceiverTransactorRaw struct {
	Contract *EtherSenderReceiverTransactor
}

func NewEtherSenderReceiver(address common.Address, backend bind.ContractBackend) (*EtherSenderReceiver, error) {
	abi, err := abi.JSON(strings.NewReader(EtherSenderReceiverABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindEtherSenderReceiver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EtherSenderReceiver{address: address, abi: abi, EtherSenderReceiverCaller: EtherSenderReceiverCaller{contract: contract}, EtherSenderReceiverTransactor: EtherSenderReceiverTransactor{contract: contract}, EtherSenderReceiverFilterer: EtherSenderReceiverFilterer{contract: contract}}, nil
}

func NewEtherSenderReceiverCaller(address common.Address, caller bind.ContractCaller) (*EtherSenderReceiverCaller, error) {
	contract, err := bindEtherSenderReceiver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EtherSenderReceiverCaller{contract: contract}, nil
}

func NewEtherSenderReceiverTransactor(address common.Address, transactor bind.ContractTransactor) (*EtherSenderReceiverTransactor, error) {
	contract, err := bindEtherSenderReceiver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EtherSenderReceiverTransactor{contract: contract}, nil
}

func NewEtherSenderReceiverFilterer(address common.Address, filterer bind.ContractFilterer) (*EtherSenderReceiverFilterer, error) {
	contract, err := bindEtherSenderReceiver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EtherSenderReceiverFilterer{contract: contract}, nil
}

func bindEtherSenderReceiver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EtherSenderReceiverMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_EtherSenderReceiver *EtherSenderReceiverRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EtherSenderReceiver.Contract.EtherSenderReceiverCaller.contract.Call(opts, result, method, params...)
}

func (_EtherSenderReceiver *EtherSenderReceiverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EtherSenderReceiver.Contract.EtherSenderReceiverTransactor.contract.Transfer(opts)
}

func (_EtherSenderReceiver *EtherSenderReceiverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EtherSenderReceiver.Contract.EtherSenderReceiverTransactor.contract.Transact(opts, method, params...)
}

func (_EtherSenderReceiver *EtherSenderReceiverCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EtherSenderReceiver.Contract.contract.Call(opts, result, method, params...)
}

func (_EtherSenderReceiver *EtherSenderReceiverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EtherSenderReceiver.Contract.contract.Transfer(opts)
}

func (_EtherSenderReceiver *EtherSenderReceiverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EtherSenderReceiver.Contract.contract.Transact(opts, method, params...)
}

func (_EtherSenderReceiver *EtherSenderReceiverCaller) GetFee(opts *bind.CallOpts, destinationChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	var out []interface{}
	err := _EtherSenderReceiver.contract.Call(opts, &out, "getFee", destinationChainSelector, message)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_EtherSenderReceiver *EtherSenderReceiverSession) GetFee(destinationChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	return _EtherSenderReceiver.Contract.GetFee(&_EtherSenderReceiver.CallOpts, destinationChainSelector, message)
}

func (_EtherSenderReceiver *EtherSenderReceiverCallerSession) GetFee(destinationChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error) {
	return _EtherSenderReceiver.Contract.GetFee(&_EtherSenderReceiver.CallOpts, destinationChainSelector, message)
}

func (_EtherSenderReceiver *EtherSenderReceiverCaller) GetRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EtherSenderReceiver.contract.Call(opts, &out, "getRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_EtherSenderReceiver *EtherSenderReceiverSession) GetRouter() (common.Address, error) {
	return _EtherSenderReceiver.Contract.GetRouter(&_EtherSenderReceiver.CallOpts)
}

func (_EtherSenderReceiver *EtherSenderReceiverCallerSession) GetRouter() (common.Address, error) {
	return _EtherSenderReceiver.Contract.GetRouter(&_EtherSenderReceiver.CallOpts)
}

func (_EtherSenderReceiver *EtherSenderReceiverCaller) IWeth(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EtherSenderReceiver.contract.Call(opts, &out, "i_weth")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_EtherSenderReceiver *EtherSenderReceiverSession) IWeth() (common.Address, error) {
	return _EtherSenderReceiver.Contract.IWeth(&_EtherSenderReceiver.CallOpts)
}

func (_EtherSenderReceiver *EtherSenderReceiverCallerSession) IWeth() (common.Address, error) {
	return _EtherSenderReceiver.Contract.IWeth(&_EtherSenderReceiver.CallOpts)
}

func (_EtherSenderReceiver *EtherSenderReceiverCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _EtherSenderReceiver.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_EtherSenderReceiver *EtherSenderReceiverSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _EtherSenderReceiver.Contract.SupportsInterface(&_EtherSenderReceiver.CallOpts, interfaceId)
}

func (_EtherSenderReceiver *EtherSenderReceiverCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _EtherSenderReceiver.Contract.SupportsInterface(&_EtherSenderReceiver.CallOpts, interfaceId)
}

func (_EtherSenderReceiver *EtherSenderReceiverCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _EtherSenderReceiver.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_EtherSenderReceiver *EtherSenderReceiverSession) TypeAndVersion() (string, error) {
	return _EtherSenderReceiver.Contract.TypeAndVersion(&_EtherSenderReceiver.CallOpts)
}

func (_EtherSenderReceiver *EtherSenderReceiverCallerSession) TypeAndVersion() (string, error) {
	return _EtherSenderReceiver.Contract.TypeAndVersion(&_EtherSenderReceiver.CallOpts)
}

func (_EtherSenderReceiver *EtherSenderReceiverTransactor) CcipReceive(opts *bind.TransactOpts, message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _EtherSenderReceiver.contract.Transact(opts, "ccipReceive", message)
}

func (_EtherSenderReceiver *EtherSenderReceiverSession) CcipReceive(message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _EtherSenderReceiver.Contract.CcipReceive(&_EtherSenderReceiver.TransactOpts, message)
}

func (_EtherSenderReceiver *EtherSenderReceiverTransactorSession) CcipReceive(message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _EtherSenderReceiver.Contract.CcipReceive(&_EtherSenderReceiver.TransactOpts, message)
}

func (_EtherSenderReceiver *EtherSenderReceiverTransactor) CcipSend(opts *bind.TransactOpts, destinationChainSelector uint64, message ClientEVM2AnyMessage) (*types.Transaction, error) {
	return _EtherSenderReceiver.contract.Transact(opts, "ccipSend", destinationChainSelector, message)
}

func (_EtherSenderReceiver *EtherSenderReceiverSession) CcipSend(destinationChainSelector uint64, message ClientEVM2AnyMessage) (*types.Transaction, error) {
	return _EtherSenderReceiver.Contract.CcipSend(&_EtherSenderReceiver.TransactOpts, destinationChainSelector, message)
}

func (_EtherSenderReceiver *EtherSenderReceiverTransactorSession) CcipSend(destinationChainSelector uint64, message ClientEVM2AnyMessage) (*types.Transaction, error) {
	return _EtherSenderReceiver.Contract.CcipSend(&_EtherSenderReceiver.TransactOpts, destinationChainSelector, message)
}

func (_EtherSenderReceiver *EtherSenderReceiverTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EtherSenderReceiver.contract.RawTransact(opts, nil)
}

func (_EtherSenderReceiver *EtherSenderReceiverSession) Receive() (*types.Transaction, error) {
	return _EtherSenderReceiver.Contract.Receive(&_EtherSenderReceiver.TransactOpts)
}

func (_EtherSenderReceiver *EtherSenderReceiverTransactorSession) Receive() (*types.Transaction, error) {
	return _EtherSenderReceiver.Contract.Receive(&_EtherSenderReceiver.TransactOpts)
}

func (_EtherSenderReceiver *EtherSenderReceiver) Address() common.Address {
	return _EtherSenderReceiver.address
}

type EtherSenderReceiverInterface interface {
	GetFee(opts *bind.CallOpts, destinationChainSelector uint64, message ClientEVM2AnyMessage) (*big.Int, error)

	GetRouter(opts *bind.CallOpts) (common.Address, error)

	IWeth(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	CcipReceive(opts *bind.TransactOpts, message ClientAny2EVMMessage) (*types.Transaction, error)

	CcipSend(opts *bind.TransactOpts, destinationChainSelector uint64, message ClientEVM2AnyMessage) (*types.Transaction, error)

	Receive(opts *bind.TransactOpts) (*types.Transaction, error)

	Address() common.Address
}
