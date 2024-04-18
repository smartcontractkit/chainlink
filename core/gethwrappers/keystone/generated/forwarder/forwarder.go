// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package forwarder

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
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

var KeystoneForwarderMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getTransmitter\",\"inputs\":[{\"name\":\"workflowExecutionId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"report\",\"inputs\":[{\"name\":\"targetAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signatures\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReportForwarded\",\"inputs\":[{\"name\":\"workflowExecutionId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"success\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InvalidData\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidSignature\",\"inputs\":[{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ReentrantCall\",\"inputs\":[]}]",
	Bin: "0x3630383036303430353233343830313536313030313035373630303038306664356235303333383036303030383136313030363735373630343035313632343631626364363065353162383135323630323036303034383230313532363031383630323438323031353237663433363136653665366637343230373336353734323036663737366536353732323037343666323037613635373236663030303030303030303030303030303036303434383230313532363036343031356236303430353138303931303339306664356236303030383035343630303136303031363061303162303331393136363030313630303136306130316230333834383131363931393039313137393039313535383131363135363130303937353736313030393738313631303039663536356235303530353036313031343835363562333336303031363030313630613031623033383231363033363130306637353736303430353136323436316263643630653531623831353236303230363030343832303135323630313736303234383230313532376634333631366536653666373432303734373236313665373336363635373232303734366632303733363536633636303030303030303030303030303030303030363034343832303135323630363430313631303035653536356236303031383035343630303136303031363061303162303331393136363030313630303136306130316230333833383131363931383231373930393235353630303038303534363034303531393239333136393137666564383838396635363033323665623133383932306438343231393266306562336464323262346631333963383761326335373533386530356261653132373839313930613335303536356236313063396338303631303135373630303033393630303066336665363038303630343035323334383031353631303031303537363030303830666435623530363030343336313036313030373235373630303033353630653031633830363363303936356463333131363130303530353738303633633039363564633331343631303130383537383036336536623731343538313436313031326235373830363366326664653338623134363130313631353736303030383066643562383036333138316635613737313436313030373735373830363337396261353039373134363130306266353738303633386461356362356231343631303063393537356236303030383066643562363034303830353138303832303138323532363031373831353237663462363537393733373436663665363534363666373237373631373236343635373232303331326533303265333030303030303030303030303030303030303036303230383230313532393035313631303062363931393036313038373835363562363034303531383039313033393066333562363130306337363130313734353635623030356236303030353437336666666666666666666666666666666666666666666666666666666666666666666666666666666631363562363034303531373366666666666666666666666666666666666666666666666666666666666666666666666666666666393039313136383135323630323030313631303062363536356236313031316236313031313633363630303436313038626235363562363130323736353635623630343035313930313531353831353236303230303136313030623635363562363130306533363130313339333636303034363130393937353635623630303039303831353236303033363032303532363034303930323035343733666666666666666666666666666666666666666666666666666666666666666666666666666666663136393035363562363130306337363130313666333636303034363130396230353635623631303632613536356236303031353437336666666666666666666666666666666666666666666666666666666666666666666666666666666631363333313436313031666135373630343035313766303863333739613030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303831353236303230363030343832303135323630313636303234383230313532376634643735373337343230363236353230373037323666373036663733363536343230366637373665363537323030303030303030303030303030303030303030363034343832303135323630363430313562363034303531383039313033393066643562363030303830353433333766666666666666666666666666666666666666666666666666303030303030303030303030303030303030303030303030303030303030303030303030303030303830383331363832313738343535363030313830353439303931313639303535363034303531373366666666666666666666666666666666666666666666666666666666666666666666666666666666393039323136393239303931383339313766386265303037396335333136353931343133343463643166643061346632383431393439376639373232613364616166653362343138366636623634353765303931613335303536356236303032353436303030393036306666313631353631303262363537363034303531376633376564333265383030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030383135323630303430313630343035313830393130333930666435623630303238303534376666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666663030313636303031313739303535363130326564363034303630303436313039666135363562383431303135363130333261353738343834363034303531376632613632363039623030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030383135323630303430313631303166313932393139303631306131333536356236303030363130333339383536303034383138393631306136303536356238313031393036313033343639313930363130616239353635623830353136303230383230313230393039313530363030303562383438313130313536313034333835373630303038303630303036313033633138393839383638313831313036313033373535373631303337353631306238383536356239303530363032303032383130313930363130333837393139303631306262373536356238303830363031663031363032303830393130343032363032303031363034303531393038313031363034303532383039333932393139303831383135323630323030313833383338303832383433373630303039323031393139303931353235303631303633653932353035303530353635623630343038303531363030303831353236303230383130313830383335323861393035323630666638333136393138313031393139303931353236303630383130313834393035323630383038313031383339303532393239353530393039333530393135303630303139303630613030313630323036303430353136303230383130333930383038343033393038353561666131353830313536313034316335373364363030303830336533643630303066643562353038353934353036313034333039333530383439323530363130633233393135303530353635623931353035303631303335333536356235303630343035313766373537333462653830303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303831353236303030393037335f5f2432386431363265353063393531373030346165623032666165616464363234386435245f5f39303633373537333462653839303631303438633930383639303630303430313631303837383536356236303430383035313830383330333831383635616634313538303135363130346138353733643630303038303365336436303030666435623530353035303530363034303531336436303166313936303166383230313136383230313830363034303532353038313031393036313034636339313930363130633562353635623630303038313831353236303033363032303532363034303930323035343930393235303733666666666666666666666666666666666666666666666666666666666666666666666666666666663136313539303530363130353037353736303030393335303530353035303631303566393536356236303030383937336666666666666666666666666666666666666666666666666666666666666666666666666666666631363839383936303430353136313035333039323931393036313063376635363562363030303630343035313830383330333831363030303836356166313931353035303364383036303030383131343631303536643537363034303531393135303630316631393630336633643031313638323031363034303532336438323532336436303030363032303834303133653631303537323536356236303630393135303562353035303630303038333831353236303033363032303532363034303930383139303230383035343766666666666666666666666666666666666666666666666666303030303030303030303030303030303030303030303030303030303030303030303030303030303136333339303831313739303931353539303531393139323530393038333930376631333236623337653330373133656366336361353738633865323364613632363137653034373134323635373837373733336639646166353932616432366239393036313035653839303835313531353831353236303230303139303536356236303430353138303931303339306133363030313934353035303530353035303562363030323830353437666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666303031363930353539353934353035303530353035303536356236313036333236313036396335363562363130363362383136313037316635363562353035363562363030303830363030303833353136303431313436313036383035373833363034303531376632616466646333303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030383135323630303430313631303166313931393036313038373835363562353035303530363032303831303135313630343038323031353136303630383330313531363030303161393139333930393235303536356236303030353437336666666666666666666666666666666666666666666666666666666666666666666666666666666631363333313436313037316435373630343035313766303863333739613030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303831353236303230363030343832303135323630313636303234383230313532376634663665366337393230363336313663366336313632366336353230363237393230366637373665363537323030303030303030303030303030303030303030363034343832303135323630363430313631303166313536356235363562333337336666666666666666666666666666666666666666666666666666666666666666666666666666666638323136303336313037396535373630343035313766303863333739613030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303831353236303230363030343832303135323630313736303234383230313532376634333631366536653666373432303734373236313665373336363635373232303734366632303733363536633636303030303030303030303030303030303030363034343832303135323630363430313631303166313536356236303031383035343766666666666666666666666666666666666666666666666666303030303030303030303030303030303030303030303030303030303030303030303030303030303136373366666666666666666666666666666666666666666666666666666666666666666666666666666666383338313136393138323137393039323535363030303830353436303430353139323933313639313766656438383839663536303332366562313338393230643834323139326630656233646432326234663133396338376132633537353338653035626165313237383931393061333530353635623630303038313531383038343532363030303562383138313130313536313038336135373630323038313835303138313031353138363833303138323031353230313631303831653536356235303630303036303230383238363031303135323630323037666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666653036303166383330313136383530313031393135303530393239313530353035363562363032303831353236303030363130383862363032303833303138343631303831343536356239333932353035303530353635623830333537336666666666666666666666666666666666666666666666666666666666666666666666666666666638313136383131343631303862363537363030303830666435623931393035303536356236303030383036303030383036303030363036303836383830333132313536313038643335373630303038306664356236313038646338363631303839323536356239343530363032303836303133353637666666666666666666666666666666663830383231313135363130386639353736303030383066643562383138383031393135303838363031663833303131323631303930643537363030303830666435623831333538313831313131353631303931633537363030303830666435623839363032303832383530313031313131353631303932653537363030303830666435623630323038333031393635303830393535303530363034303838303133353931353038303832313131353631303934633537363030303830666435623831383830313931353038383630316638333031313236313039363035373630303038306664356238313335383138313131313536313039366635373630303038306664356238393630323038323630303531623835303130313131313536313039383435373630303038306664356239363939393539383530393339363530363032303031393439333932353035303530353635623630303036303230383238343033313231353631303961393537363030303830666435623530333539313930353035363562363030303630323038323834303331323135363130396332353736303030383066643562363130383862383236313038393235363562376634653438376237313030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030363030303532363031313630303435323630323436303030666435623830383230313830383231313135363130613064353736313061306436313039636235363562393239313530353035363562363032303831353238313630323038323031353238313833363034303833303133373630303038313833303136303430393038313031393139303931353236303166393039323031376666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666530313630313031393139303530353635623630303038303835383531313135363130613730353736303030383066643562383338363131313536313061376435373630303038306664356235303530383230313933393139303932303339313530353635623766346534383762373130303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303630303035323630343136303034353236303234363030306664356236303030363032303832383430333132313536313061636235373630303038306664356238313335363766666666666666666666666666666666383038323131313536313061653335373630303038306664356238313834303139313530383436303166383330313132363130616637353736303030383066643562383133353831383131313135363130623039353736313062303936313061386135363562363034303531363031663832303137666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666653039303831313636303366303131363831303139303833383231313831383331303137313536313062346635373631306234663631306138613536356238313630343035323832383135323837363032303834383730313031313131353631306236383537363030303830666435623832363032303836303136303230383330313337363030303932383130313630323030313932393039323532353039353934353035303530353035303536356237663465343837623731303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303036303030353236303332363030343532363032343630303066643562363030303830383333353766666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666665313834333630333031383131323631306265633537363030303830666435623833303138303335393135303637666666666666666666666666666666663832313131353631306330373537363030303830666435623630323030313931353033363831393030333832313331353631306331633537363030303830666435623932353039323930353035363562363030303766666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666663832303336313063353435373631306335343631303963623536356235303630303130313930353635623630303038303630343038333835303331323135363130633665353736303030383066643562353035303830353136303230393039313031353139303932393039313530353635623831383338323337363030303931303139303831353239313930353035366665613136343733366636633633343330303038313330303061",
}

var KeystoneForwarderABI = KeystoneForwarderMetaData.ABI

var KeystoneForwarderBin = KeystoneForwarderMetaData.Bin

func DeployKeystoneForwarder(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *KeystoneForwarder, error) {
	parsed, err := KeystoneForwarderMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeystoneForwarderBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeystoneForwarder{address: address, abi: *parsed, KeystoneForwarderCaller: KeystoneForwarderCaller{contract: contract}, KeystoneForwarderTransactor: KeystoneForwarderTransactor{contract: contract}, KeystoneForwarderFilterer: KeystoneForwarderFilterer{contract: contract}}, nil
}

type KeystoneForwarder struct {
	address common.Address
	abi     abi.ABI
	KeystoneForwarderCaller
	KeystoneForwarderTransactor
	KeystoneForwarderFilterer
}

type KeystoneForwarderCaller struct {
	contract *bind.BoundContract
}

type KeystoneForwarderTransactor struct {
	contract *bind.BoundContract
}

type KeystoneForwarderFilterer struct {
	contract *bind.BoundContract
}

type KeystoneForwarderSession struct {
	Contract     *KeystoneForwarder
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeystoneForwarderCallerSession struct {
	Contract *KeystoneForwarderCaller
	CallOpts bind.CallOpts
}

type KeystoneForwarderTransactorSession struct {
	Contract     *KeystoneForwarderTransactor
	TransactOpts bind.TransactOpts
}

type KeystoneForwarderRaw struct {
	Contract *KeystoneForwarder
}

type KeystoneForwarderCallerRaw struct {
	Contract *KeystoneForwarderCaller
}

type KeystoneForwarderTransactorRaw struct {
	Contract *KeystoneForwarderTransactor
}

func NewKeystoneForwarder(address common.Address, backend bind.ContractBackend) (*KeystoneForwarder, error) {
	abi, err := abi.JSON(strings.NewReader(KeystoneForwarderABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeystoneForwarder(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarder{address: address, abi: abi, KeystoneForwarderCaller: KeystoneForwarderCaller{contract: contract}, KeystoneForwarderTransactor: KeystoneForwarderTransactor{contract: contract}, KeystoneForwarderFilterer: KeystoneForwarderFilterer{contract: contract}}, nil
}

func NewKeystoneForwarderCaller(address common.Address, caller bind.ContractCaller) (*KeystoneForwarderCaller, error) {
	contract, err := bindKeystoneForwarder(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderCaller{contract: contract}, nil
}

func NewKeystoneForwarderTransactor(address common.Address, transactor bind.ContractTransactor) (*KeystoneForwarderTransactor, error) {
	contract, err := bindKeystoneForwarder(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderTransactor{contract: contract}, nil
}

func NewKeystoneForwarderFilterer(address common.Address, filterer bind.ContractFilterer) (*KeystoneForwarderFilterer, error) {
	contract, err := bindKeystoneForwarder(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderFilterer{contract: contract}, nil
}

func bindKeystoneForwarder(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeystoneForwarderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_KeystoneForwarder *KeystoneForwarderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeystoneForwarder.Contract.KeystoneForwarderCaller.contract.Call(opts, result, method, params...)
}

func (_KeystoneForwarder *KeystoneForwarderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.KeystoneForwarderTransactor.contract.Transfer(opts)
}

func (_KeystoneForwarder *KeystoneForwarderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.KeystoneForwarderTransactor.contract.Transact(opts, method, params...)
}

func (_KeystoneForwarder *KeystoneForwarderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeystoneForwarder.Contract.contract.Call(opts, result, method, params...)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.contract.Transfer(opts)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.contract.Transact(opts, method, params...)
}

func (_KeystoneForwarder *KeystoneForwarderCaller) GetTransmitter(opts *bind.CallOpts, workflowExecutionId [32]byte) (common.Address, error) {
	var out []interface{}
	err := _KeystoneForwarder.contract.Call(opts, &out, "getTransmitter", workflowExecutionId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeystoneForwarder *KeystoneForwarderSession) GetTransmitter(workflowExecutionId [32]byte) (common.Address, error) {
	return _KeystoneForwarder.Contract.GetTransmitter(&_KeystoneForwarder.CallOpts, workflowExecutionId)
}

func (_KeystoneForwarder *KeystoneForwarderCallerSession) GetTransmitter(workflowExecutionId [32]byte) (common.Address, error) {
	return _KeystoneForwarder.Contract.GetTransmitter(&_KeystoneForwarder.CallOpts, workflowExecutionId)
}

func (_KeystoneForwarder *KeystoneForwarderCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeystoneForwarder.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeystoneForwarder *KeystoneForwarderSession) Owner() (common.Address, error) {
	return _KeystoneForwarder.Contract.Owner(&_KeystoneForwarder.CallOpts)
}

func (_KeystoneForwarder *KeystoneForwarderCallerSession) Owner() (common.Address, error) {
	return _KeystoneForwarder.Contract.Owner(&_KeystoneForwarder.CallOpts)
}

func (_KeystoneForwarder *KeystoneForwarderCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _KeystoneForwarder.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_KeystoneForwarder *KeystoneForwarderSession) TypeAndVersion() (string, error) {
	return _KeystoneForwarder.Contract.TypeAndVersion(&_KeystoneForwarder.CallOpts)
}

func (_KeystoneForwarder *KeystoneForwarderCallerSession) TypeAndVersion() (string, error) {
	return _KeystoneForwarder.Contract.TypeAndVersion(&_KeystoneForwarder.CallOpts)
}

func (_KeystoneForwarder *KeystoneForwarderTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeystoneForwarder.contract.Transact(opts, "acceptOwnership")
}

func (_KeystoneForwarder *KeystoneForwarderSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.AcceptOwnership(&_KeystoneForwarder.TransactOpts)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.AcceptOwnership(&_KeystoneForwarder.TransactOpts)
}

func (_KeystoneForwarder *KeystoneForwarderTransactor) Report(opts *bind.TransactOpts, targetAddress common.Address, data []byte, signatures [][]byte) (*types.Transaction, error) {
	return _KeystoneForwarder.contract.Transact(opts, "report", targetAddress, data, signatures)
}

func (_KeystoneForwarder *KeystoneForwarderSession) Report(targetAddress common.Address, data []byte, signatures [][]byte) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.Report(&_KeystoneForwarder.TransactOpts, targetAddress, data, signatures)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorSession) Report(targetAddress common.Address, data []byte, signatures [][]byte) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.Report(&_KeystoneForwarder.TransactOpts, targetAddress, data, signatures)
}

func (_KeystoneForwarder *KeystoneForwarderTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.contract.Transact(opts, "transferOwnership", to)
}

func (_KeystoneForwarder *KeystoneForwarderSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.TransferOwnership(&_KeystoneForwarder.TransactOpts, to)
}

func (_KeystoneForwarder *KeystoneForwarderTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeystoneForwarder.Contract.TransferOwnership(&_KeystoneForwarder.TransactOpts, to)
}

type KeystoneForwarderOwnershipTransferRequestedIterator struct {
	Event *KeystoneForwarderOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneForwarderOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneForwarderOwnershipTransferRequested)
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

	select {
	case log := <-it.logs:
		it.Event = new(KeystoneForwarderOwnershipTransferRequested)
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

func (it *KeystoneForwarderOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeystoneForwarderOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneForwarderOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneForwarderOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderOwnershipTransferRequestedIterator{contract: _KeystoneForwarder.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneForwarderOwnershipTransferRequested)
				if err := _KeystoneForwarder.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_KeystoneForwarder *KeystoneForwarderFilterer) ParseOwnershipTransferRequested(log types.Log) (*KeystoneForwarderOwnershipTransferRequested, error) {
	event := new(KeystoneForwarderOwnershipTransferRequested)
	if err := _KeystoneForwarder.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeystoneForwarderOwnershipTransferredIterator struct {
	Event *KeystoneForwarderOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneForwarderOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneForwarderOwnershipTransferred)
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

	select {
	case log := <-it.logs:
		it.Event = new(KeystoneForwarderOwnershipTransferred)
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

func (it *KeystoneForwarderOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeystoneForwarderOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneForwarderOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneForwarderOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderOwnershipTransferredIterator{contract: _KeystoneForwarder.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneForwarderOwnershipTransferred)
				if err := _KeystoneForwarder.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_KeystoneForwarder *KeystoneForwarderFilterer) ParseOwnershipTransferred(log types.Log) (*KeystoneForwarderOwnershipTransferred, error) {
	event := new(KeystoneForwarderOwnershipTransferred)
	if err := _KeystoneForwarder.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeystoneForwarderReportForwardedIterator struct {
	Event *KeystoneForwarderReportForwarded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneForwarderReportForwardedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneForwarderReportForwarded)
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

	select {
	case log := <-it.logs:
		it.Event = new(KeystoneForwarderReportForwarded)
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

func (it *KeystoneForwarderReportForwardedIterator) Error() error {
	return it.fail
}

func (it *KeystoneForwarderReportForwardedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneForwarderReportForwarded struct {
	WorkflowExecutionId [32]byte
	Transmitter         common.Address
	Success             bool
	Raw                 types.Log
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) FilterReportForwarded(opts *bind.FilterOpts, workflowExecutionId [][32]byte, transmitter []common.Address) (*KeystoneForwarderReportForwardedIterator, error) {

	var workflowExecutionIdRule []interface{}
	for _, workflowExecutionIdItem := range workflowExecutionId {
		workflowExecutionIdRule = append(workflowExecutionIdRule, workflowExecutionIdItem)
	}
	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.FilterLogs(opts, "ReportForwarded", workflowExecutionIdRule, transmitterRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneForwarderReportForwardedIterator{contract: _KeystoneForwarder.contract, event: "ReportForwarded", logs: logs, sub: sub}, nil
}

func (_KeystoneForwarder *KeystoneForwarderFilterer) WatchReportForwarded(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderReportForwarded, workflowExecutionId [][32]byte, transmitter []common.Address) (event.Subscription, error) {

	var workflowExecutionIdRule []interface{}
	for _, workflowExecutionIdItem := range workflowExecutionId {
		workflowExecutionIdRule = append(workflowExecutionIdRule, workflowExecutionIdItem)
	}
	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}

	logs, sub, err := _KeystoneForwarder.contract.WatchLogs(opts, "ReportForwarded", workflowExecutionIdRule, transmitterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneForwarderReportForwarded)
				if err := _KeystoneForwarder.contract.UnpackLog(event, "ReportForwarded", log); err != nil {
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

func (_KeystoneForwarder *KeystoneForwarderFilterer) ParseReportForwarded(log types.Log) (*KeystoneForwarderReportForwarded, error) {
	event := new(KeystoneForwarderReportForwarded)
	if err := _KeystoneForwarder.contract.UnpackLog(event, "ReportForwarded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_KeystoneForwarder *KeystoneForwarder) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _KeystoneForwarder.abi.Events["OwnershipTransferRequested"].ID:
		return _KeystoneForwarder.ParseOwnershipTransferRequested(log)
	case _KeystoneForwarder.abi.Events["OwnershipTransferred"].ID:
		return _KeystoneForwarder.ParseOwnershipTransferred(log)
	case _KeystoneForwarder.abi.Events["ReportForwarded"].ID:
		return _KeystoneForwarder.ParseReportForwarded(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeystoneForwarderOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (KeystoneForwarderOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (KeystoneForwarderReportForwarded) Topic() common.Hash {
	return common.HexToHash("0x1326b37e30713ecf3ca578c8e23da62617e047142657877733f9daf592ad26b9")
}

func (_KeystoneForwarder *KeystoneForwarder) Address() common.Address {
	return _KeystoneForwarder.address
}

type KeystoneForwarderInterface interface {
	GetTransmitter(opts *bind.CallOpts, workflowExecutionId [32]byte) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	Report(opts *bind.TransactOpts, targetAddress common.Address, data []byte, signatures [][]byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneForwarderOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*KeystoneForwarderOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneForwarderOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*KeystoneForwarderOwnershipTransferred, error)

	FilterReportForwarded(opts *bind.FilterOpts, workflowExecutionId [][32]byte, transmitter []common.Address) (*KeystoneForwarderReportForwardedIterator, error)

	WatchReportForwarded(opts *bind.WatchOpts, sink chan<- *KeystoneForwarderReportForwarded, workflowExecutionId [][32]byte, transmitter []common.Address) (event.Subscription, error)

	ParseReportForwarded(log types.Log) (*KeystoneForwarderReportForwarded, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
