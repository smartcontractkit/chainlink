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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"name\":\"InvalidRouter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"gotToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"expectedToken\",\"type\":\"address\"}],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gotAmounts\",\"type\":\"uint256\"}],\"name\":\"InvalidTokenAmounts\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gotAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"msgValue\",\"type\":\"uint256\"}],\"name\":\"TokenAmountNotEqualToMsgValue\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structClient.Any2EVMMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"ccipReceive\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destinationChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structClient.EVM2AnyMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"ccipSend\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destinationChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structClient.EVM2AnyMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"getFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRouter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_weth\",\"outputs\":[{\"internalType\":\"contractIWrappedNative\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60c060409080825234620001275762000032906200154e803803809162000027828562000196565b8339810190620001d0565b6001600160a01b03908082169081156200017e57608052825163e861e90760e01b81526020918282600481845afa9182156200017357600094849360449287916200013f575b5016918260a0528651958693849263095ea7b360e01b84526004840152811960248401525af180156200013457620000f5575b825161135c9081620001f2823960805181818160c80152818161017b015281816105c00152610cda015260a0518181816102610152818161053001528181610c6001526111f10152f35b81813d83116200012c575b6200010c818362000196565b810103126200012757518015150362000127573880620000ab565b600080fd5b503d62000100565b83513d6000823e3d90fd5b620001649150853d87116200016b575b6200015b818362000196565b810190620001d0565b3862000078565b503d6200014f565b85513d6000823e3d90fd5b83516335fdcccd60e21b815260006004820152602490fd5b601f909101601f19168101906001600160401b03821190821017620001ba57604052565b634e487b7160e01b600052604160045260246000fd5b908160209103126200012757516001600160a01b038116810362000127579056fe6080604081815260049182361015610022575b505050361561002057600080fd5b005b600092833560e01c91826301ffc9a7146106bb57508163181f5a771461063c57816320487ded146105545781634dbe7e92146104e557816385572ffb1461010c5750806396f4e9f9146100f05763b0f479a11461007f5780610012565b346100ec57817ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100ec576020905173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b5080fd5b6020826101056100ff366108cb565b90610bc1565b9051908152f35b919050346104e1576020907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc9082823601126104dd5783359167ffffffffffffffff908184116104135760a09084360301126104035773ffffffffffffffffffffffffffffffffffffffff90817f00000000000000000000000000000000000000000000000000000000000000001633036104ae578251936101ad85610777565b808701358552602481013582811681036104aa578686015260448101358281116104aa576101e090883691840101610a38565b8486015260648101358281116104aa576101ff90883691840101610a38565b916060860192835260848201359081116104aa57608091886102249236920101610aa0565b9401908482525185818051810103126104a657850151938285168095036104a657516001810361047757508161025a8251610b6d565b51511691807f0000000000000000000000000000000000000000000000000000000000000000168093036104345750610294859151610b6d565b51015190803b156104135782517f2e1a7d4d0000000000000000000000000000000000000000000000000000000081528287820152878160248183865af1801561042a57610417575b508680808085885af16102ee61124f565b50156102f8578680f35b803b15610413578251937fd0e30db0000000000000000000000000000000000000000000000000000000008552878086898187875af1978815610407578798979596976103eb575b61039b97508651978895869485937fa9059cbb00000000000000000000000000000000000000000000000000000000855284016020909392919373ffffffffffffffffffffffffffffffffffffffff60408201951681520152565b03925af19081156103e257506103b4575b808080808680f35b816103d392903d106103db575b6103cb81836107f2565b810190610ba9565b5038806103ac565b503d6103c1565b513d85823e3d90fd5b94506103f781976107c2565b61040357858794610340565b8580fd5b508451903d90823e3d90fd5b8680fd5b610423909791976107c2565b95386102dd565b84513d8a823e3d90fd5b6044945061044487949251610b6d565b5151169051927f0fc746a10000000000000000000000000000000000000000000000000000000084528301526024820152fd5b866024918551917f83b9f0ae000000000000000000000000000000000000000000000000000000008352820152fd5b8780fd5b8880fd5b82517fd7f733340000000000000000000000000000000000000000000000000000000081523381880152602490fd5b8480fd5b8280fd5b5050346100ec57817ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100ec576020905173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b919050346104e157602061057261056a366108cb565b9190916110e6565b936105a7845195869384937f20487ded000000000000000000000000000000000000000000000000000000008552840161092f565b038173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa9182156106325783926105fa575b6020838351908152f35b9091506020813d60201161062a575b81610616602093836107f2565b810103126104e157602092505190386105f0565b3d9150610609565b81513d85823e3d90fd5b5050346100ec57817ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100ec5780516106b79161067b826107d6565b601982527f457468657253656e646572526563656976657220312e352e300000000000000060208301525191829160208352602083019061086d565b0390f35b8491346104e15760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126104e157357fffffffff0000000000000000000000000000000000000000000000000000000081168091036104e157602092507f85572ffb00000000000000000000000000000000000000000000000000000000811490811561074d575b5015158152f35b7f01ffc9a70000000000000000000000000000000000000000000000000000000091501483610746565b60a0810190811067ffffffffffffffff82111761079357604052565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b67ffffffffffffffff811161079357604052565b6040810190811067ffffffffffffffff82111761079357604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff82111761079357604052565b67ffffffffffffffff811161079357601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b919082519283825260005b8481106108b75750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8460006020809697860101520116010190565b602081830181015184830182015201610878565b907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc9160408382011261092a5767ffffffffffffffff90600435828116810361092a579360243592831161092a578260a09203011261092a5760040190565b600080fd5b9092919267ffffffffffffffff604091168252602091604083820152610961855160a0604084015260e083019061086d565b9161099b84870151937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0948585830301606086015261086d565b604087015194848483030160808501528080875193848152019601926000905b838210610a0257505050505060c060808673ffffffffffffffffffffffffffffffffffffffff60606109ff989901511660a08501520151928285030191015261086d565b90565b8451805173ffffffffffffffffffffffffffffffffffffffff1689528301518884015296870196938201936001909101906109bb565b81601f8201121561092a57803590610a4f82610833565b92610a5d60405194856107f2565b8284526020838301011161092a57816000926020809301838601378301015290565b359073ffffffffffffffffffffffffffffffffffffffff8216820361092a57565b81601f8201121561092a5767ffffffffffffffff906020813583811161079357604093845195610ad5848460051b01886107f2565b828752838088019360061b8601019481861161092a578401925b858410610b00575050505050505090565b868483031261092a578651908782019082821085831117610b3f57889287928452610b2a87610a7f565b81528287013583820152815201930192610aef565b602460007f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b805115610b7a5760200190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b9081602091031261092a5751801515810361092a5790565b90604091600090838301357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112156104e15783019283359367ffffffffffffffff85116110e2576020948060061b3603868301136104dd57156110b5578501359060608101359173ffffffffffffffffffffffffffffffffffffffff9283811680910361040357611078575b50610c5c906110e6565b91817f0000000000000000000000000000000000000000000000000000000000000000169185610c8e88860151610b6d565b510151833b15610403579085918851928380927fd0e30db00000000000000000000000000000000000000000000000000000000082526004958691895af1801561106e5761105b575b507f000000000000000000000000000000000000000000000000000000000000000093818516948951927f20487ded000000000000000000000000000000000000000000000000000000008452898480610d348b8a8a840161092f565b03818a5afa938415611051578994611022575b5060608801818151168c81610de657505050505050509185949391610d9c969347908951988995869485937f96f4e9f9000000000000000000000000000000000000000000000000000000008552840161092f565b03925af1938415610ddc575092610db257505090565b90809250813d8311610dd5575b610dc981836107f2565b8101031261092a575190565b503d610dbf565b51903d90823e3d90fd5b610e7c918c808f9380517f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564868201927f23b872dd0000000000000000000000000000000000000000000000000000000084523360248401523060448401528d606484015260648352610e5783610777565b5196610e62886107d6565b808852870152519082855af1610e7661124f565b9161127f565b8051908c82610f84575b5050505116918203610ed4575b505050918493918695610d9c97948951988995869485937f96f4e9f9000000000000000000000000000000000000000000000000000000008552840161092f565b918891610f3293898c518096819582947f095ea7b30000000000000000000000000000000000000000000000000000000084528a84016020909392919373ffffffffffffffffffffffffffffffffffffffff60408201951681520152565b03925af18015610f7a579587969281969492610d9c999694610f5d575b509294978294969750610e93565b610f7390873d89116103db576103cb81836107f2565b5038610f4f565b88513d88823e3d90fd5b80610f93938301019101610ba9565b15610fa05738808c610e86565b6084868c8e51917f08c379a0000000000000000000000000000000000000000000000000000000008352820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152fd5b9093508981813d831161104a575b61103a81836107f2565b8101031261092a57519238610d47565b503d611030565b8b513d8b823e3d90fd5b611067909691966107c2565b9438610cd7565b89513d89823e3d90fd5b348114610c52576044908751907fba2f74670000000000000000000000000000000000000000000000000000000082526004820152346024820152fd5b6024847f4e487b710000000000000000000000000000000000000000000000000000000081526032600452fd5b8380fd5b604080516110f381610777565b60609182825282602083015282818301528260806000938483820152015260a0843603126100ec5780519361112785610777565b67ffffffffffffffff81358181116104dd576111469036908401610a38565b865260208201358181116104dd576111619036908401610a38565b9360208701948552838301358281116100ec576111819036908501610aa0565b95848801968752611193818501610a7f565b90880152608083013591821161124c57506111b091369101610a38565b60808501528251516001810361121d57505190336020830152602082526111d6826107d6565b5261121873ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169151610b6d565b515290565b60249151907f83b9f0ae0000000000000000000000000000000000000000000000000000000082526004820152fd5b80fd5b3d1561127a573d9061126082610833565b9161126e60405193846107f2565b82523d6000602084013e565b606090565b919290156112fa5750815115611293575090565b3b1561129c5790565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152fd5b82519091501561130d5750805190602001fd5b61134b906040519182917f08c379a000000000000000000000000000000000000000000000000000000000835260206004840152602483019061086d565b0390fdfea164736f6c6343000818000a",
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
