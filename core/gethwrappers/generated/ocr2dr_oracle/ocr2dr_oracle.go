// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ocr2dr_oracle

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
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated"
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

var OCR2DROracleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EmptySendersList\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRequestID\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonceMustBeUnique\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAllowedToSetSenders\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSender\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"senders\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"changedBy\",\"type\":\"address\"}],\"name\":\"AuthorizedSendersChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"OracleRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"OracleResponse\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"fulfillRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAuthorizedSenders\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"isAuthorizedSender\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subscriptionId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"sendRequest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"senders\",\"type\":\"address[]\"}],\"name\":\"setAuthorizedSenders\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161116638038061116683398101604081905261002f91610172565b808060006001600160a01b03821661008e5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600380546001600160a01b0319166001600160a01b03848116919091179091558116156100be576100be816100c7565b505050506101a2565b6001600160a01b0381163314156101205760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610085565b600480546001600160a01b0319166001600160a01b03838116918217909255600354604051919216907fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127890600090a350565b60006020828403121561018457600080fd5b81516001600160a01b038116811461019b57600080fd5b9392505050565b610fb5806101b16000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c8063bb9fa3f51161005b578063bb9fa3f5146100f0578063ee56997b14610111578063f2fde38b14610124578063fa00763a1461013757600080fd5b80632408afaa1461008d57806339b05122146100ab57806379ba5097146100c05780638da5cb5b146100c8575b600080fd5b61009561015a565b6040516100a29190610e15565b60405180910390f35b6100be6100b9366004610c8e565b6101c9565b005b6100be610359565b60035460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100a2565b6101036100fe366004610d08565b61045f565b6040519081526020016100a2565b6100be61011f366004610c19565b610597565b6100be610132366004610bfe565b610706565b61014a610145366004610bfe565b61071a565b60405190151581526020016100a2565b606060028054806020026020016040519081016040528092919081815260200182805480156101bf57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610194575b5050505050905090565b600085815260066020526040902054859073ffffffffffffffffffffffffffffffffffffffff16610226576040517f803ed86300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61022e61072c565b60008681526006602090815260409182902054915188815273ffffffffffffffffffffffffffffffffffffffff909216917f9e9bc7616d42c2835d05ae617e508454e63b30b934be8aa932ebc125e0e58a64910160405180910390a16040517f0ca7617500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821690630ca76175906102e4908a908a908a908a908a90600401610e92565b600060405180830381600087803b1580156102fe57600080fd5b505af1158015610312573d6000803e3d6000fd5b50505060009788525050600660205250506040842080547fffffffffffffffffffffffff000000000000000000000000000000000000000016815560010193909355505050565b60045473ffffffffffffffffffffffffffffffffffffffff1633146103df576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b600380547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560048054909116905560405173ffffffffffffffffffffffffffffffffffffffff909116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a350565b600580546000918261047083610ee2565b90915550506005546040517fffffffffffffffffffffffffffffffffffffffff0000000000000000000000003360601b1660208201526034810191909152600090605401604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018152828252805160209182012083830183523384528184018981526000828152600690935291839020935184547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9091161784559051600190930192909255519091507f9dc59d1e6d6042f6c2d2af5a3f9f6502bf0abf7cf40f0eb80edf752ca396f2ba9061058790839087908790610e6f565b60405180910390a1949350505050565b61059f61076d565b6105d5576040517fad77f06100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8061060c576040517f75158c3b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b600254811015610668576106556002828154811061062f5761062f610f79565b600091825260208220015473ffffffffffffffffffffffffffffffffffffffff166107bb565b508061066081610ee2565b91505061060f565b5060005b818110156106b9576106a683838381811061068957610689610f79565b905060200201602081019061069e9190610bfe565b6000906107e4565b50806106b181610ee2565b91505061066c565b506106c660028383610aef565b507ff263cfb3e4298332e776194610cf9fdc09ccb3ada8b9aa39764d882e11fbf0a08282336040516106fa93929190610d9d565b60405180910390a15050565b61070e610806565b61071781610887565b50565b6000610726818361097e565b92915050565b6107353361071a565b61076b576040517f0809490800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60006107783361071a565b806107b657503361079e60035473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff16145b905090565b60006107dd8373ffffffffffffffffffffffffffffffffffffffff84166109ad565b9392505050565b60006107dd8373ffffffffffffffffffffffffffffffffffffffff8416610aa0565b60035473ffffffffffffffffffffffffffffffffffffffff16331461076b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016103d6565b73ffffffffffffffffffffffffffffffffffffffff8116331415610907576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016103d6565b600480547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217909255600354604051919216907fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127890600090a350565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415156107dd565b60008181526001830160205260408120548015610a965760006109d1600183610ecb565b85549091506000906109e590600190610ecb565b9050818114610a4a576000866000018281548110610a0557610a05610f79565b9060005260206000200154905080876000018481548110610a2857610a28610f79565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080610a5b57610a5b610f4a565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610726565b6000915050610726565b6000818152600183016020526040812054610ae757508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610726565b506000610726565b828054828255906000526020600020908101928215610b67579160200282015b82811115610b675781547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff843516178255602090920191600190910190610b0f565b50610b73929150610b77565b5090565b5b80821115610b735760008155600101610b78565b803573ffffffffffffffffffffffffffffffffffffffff81168114610bb057600080fd5b919050565b60008083601f840112610bc757600080fd5b50813567ffffffffffffffff811115610bdf57600080fd5b602083019150836020828501011115610bf757600080fd5b9250929050565b600060208284031215610c1057600080fd5b6107dd82610b8c565b60008060208385031215610c2c57600080fd5b823567ffffffffffffffff80821115610c4457600080fd5b818501915085601f830112610c5857600080fd5b813581811115610c6757600080fd5b8660208260051b8501011115610c7c57600080fd5b60209290920196919550909350505050565b600080600080600060608688031215610ca657600080fd5b85359450602086013567ffffffffffffffff80821115610cc557600080fd5b610cd189838a01610bb5565b90965094506040880135915080821115610cea57600080fd5b50610cf788828901610bb5565b969995985093965092949392505050565b600080600060408486031215610d1d57600080fd5b83359250602084013567ffffffffffffffff811115610d3b57600080fd5b610d4786828701610bb5565b9497909650939450505050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b6040808252810183905260008460608301825b86811015610deb5773ffffffffffffffffffffffffffffffffffffffff610dd684610b8c565b16825260209283019290910190600101610db0565b50809250505073ffffffffffffffffffffffffffffffffffffffff83166020830152949350505050565b6020808252825182820181905260009190848201906040850190845b81811015610e6357835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101610e31565b50909695505050505050565b838152604060208201526000610e89604083018486610d54565b95945050505050565b858152606060208201526000610eac606083018688610d54565b8281036040840152610ebf818587610d54565b98975050505050505050565b600082821015610edd57610edd610f1b565b500390565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415610f1457610f14610f1b565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fdfea164736f6c6343000806000a",
}

var OCR2DROracleABI = OCR2DROracleMetaData.ABI

var OCR2DROracleBin = OCR2DROracleMetaData.Bin

func DeployOCR2DROracle(auth *bind.TransactOpts, backend bind.ContractBackend, owner common.Address) (common.Address, *types.Transaction, *OCR2DROracle, error) {
	parsed, err := OCR2DROracleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OCR2DROracleBin), backend, owner)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OCR2DROracle{OCR2DROracleCaller: OCR2DROracleCaller{contract: contract}, OCR2DROracleTransactor: OCR2DROracleTransactor{contract: contract}, OCR2DROracleFilterer: OCR2DROracleFilterer{contract: contract}}, nil
}

type OCR2DROracle struct {
	address common.Address
	abi     abi.ABI
	OCR2DROracleCaller
	OCR2DROracleTransactor
	OCR2DROracleFilterer
}

type OCR2DROracleCaller struct {
	contract *bind.BoundContract
}

type OCR2DROracleTransactor struct {
	contract *bind.BoundContract
}

type OCR2DROracleFilterer struct {
	contract *bind.BoundContract
}

type OCR2DROracleSession struct {
	Contract     *OCR2DROracle
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OCR2DROracleCallerSession struct {
	Contract *OCR2DROracleCaller
	CallOpts bind.CallOpts
}

type OCR2DROracleTransactorSession struct {
	Contract     *OCR2DROracleTransactor
	TransactOpts bind.TransactOpts
}

type OCR2DROracleRaw struct {
	Contract *OCR2DROracle
}

type OCR2DROracleCallerRaw struct {
	Contract *OCR2DROracleCaller
}

type OCR2DROracleTransactorRaw struct {
	Contract *OCR2DROracleTransactor
}

func NewOCR2DROracle(address common.Address, backend bind.ContractBackend) (*OCR2DROracle, error) {
	abi, err := abi.JSON(strings.NewReader(OCR2DROracleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOCR2DROracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCR2DROracle{address: address, abi: abi, OCR2DROracleCaller: OCR2DROracleCaller{contract: contract}, OCR2DROracleTransactor: OCR2DROracleTransactor{contract: contract}, OCR2DROracleFilterer: OCR2DROracleFilterer{contract: contract}}, nil
}

func NewOCR2DROracleCaller(address common.Address, caller bind.ContractCaller) (*OCR2DROracleCaller, error) {
	contract, err := bindOCR2DROracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2DROracleCaller{contract: contract}, nil
}

func NewOCR2DROracleTransactor(address common.Address, transactor bind.ContractTransactor) (*OCR2DROracleTransactor, error) {
	contract, err := bindOCR2DROracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2DROracleTransactor{contract: contract}, nil
}

func NewOCR2DROracleFilterer(address common.Address, filterer bind.ContractFilterer) (*OCR2DROracleFilterer, error) {
	contract, err := bindOCR2DROracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCR2DROracleFilterer{contract: contract}, nil
}

func bindOCR2DROracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OCR2DROracleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_OCR2DROracle *OCR2DROracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2DROracle.Contract.OCR2DROracleCaller.contract.Call(opts, result, method, params...)
}

func (_OCR2DROracle *OCR2DROracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2DROracle.Contract.OCR2DROracleTransactor.contract.Transfer(opts)
}

func (_OCR2DROracle *OCR2DROracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2DROracle.Contract.OCR2DROracleTransactor.contract.Transact(opts, method, params...)
}

func (_OCR2DROracle *OCR2DROracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2DROracle.Contract.contract.Call(opts, result, method, params...)
}

func (_OCR2DROracle *OCR2DROracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2DROracle.Contract.contract.Transfer(opts)
}

func (_OCR2DROracle *OCR2DROracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2DROracle.Contract.contract.Transact(opts, method, params...)
}

func (_OCR2DROracle *OCR2DROracleCaller) GetAuthorizedSenders(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _OCR2DROracle.contract.Call(opts, &out, "getAuthorizedSenders")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_OCR2DROracle *OCR2DROracleSession) GetAuthorizedSenders() ([]common.Address, error) {
	return _OCR2DROracle.Contract.GetAuthorizedSenders(&_OCR2DROracle.CallOpts)
}

func (_OCR2DROracle *OCR2DROracleCallerSession) GetAuthorizedSenders() ([]common.Address, error) {
	return _OCR2DROracle.Contract.GetAuthorizedSenders(&_OCR2DROracle.CallOpts)
}

func (_OCR2DROracle *OCR2DROracleCaller) IsAuthorizedSender(opts *bind.CallOpts, sender common.Address) (bool, error) {
	var out []interface{}
	err := _OCR2DROracle.contract.Call(opts, &out, "isAuthorizedSender", sender)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_OCR2DROracle *OCR2DROracleSession) IsAuthorizedSender(sender common.Address) (bool, error) {
	return _OCR2DROracle.Contract.IsAuthorizedSender(&_OCR2DROracle.CallOpts, sender)
}

func (_OCR2DROracle *OCR2DROracleCallerSession) IsAuthorizedSender(sender common.Address) (bool, error) {
	return _OCR2DROracle.Contract.IsAuthorizedSender(&_OCR2DROracle.CallOpts, sender)
}

func (_OCR2DROracle *OCR2DROracleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OCR2DROracle.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OCR2DROracle *OCR2DROracleSession) Owner() (common.Address, error) {
	return _OCR2DROracle.Contract.Owner(&_OCR2DROracle.CallOpts)
}

func (_OCR2DROracle *OCR2DROracleCallerSession) Owner() (common.Address, error) {
	return _OCR2DROracle.Contract.Owner(&_OCR2DROracle.CallOpts)
}

func (_OCR2DROracle *OCR2DROracleTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2DROracle.contract.Transact(opts, "acceptOwnership")
}

func (_OCR2DROracle *OCR2DROracleSession) AcceptOwnership() (*types.Transaction, error) {
	return _OCR2DROracle.Contract.AcceptOwnership(&_OCR2DROracle.TransactOpts)
}

func (_OCR2DROracle *OCR2DROracleTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OCR2DROracle.Contract.AcceptOwnership(&_OCR2DROracle.TransactOpts)
}

func (_OCR2DROracle *OCR2DROracleTransactor) FulfillRequest(opts *bind.TransactOpts, requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _OCR2DROracle.contract.Transact(opts, "fulfillRequest", requestId, response, err)
}

func (_OCR2DROracle *OCR2DROracleSession) FulfillRequest(requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _OCR2DROracle.Contract.FulfillRequest(&_OCR2DROracle.TransactOpts, requestId, response, err)
}

func (_OCR2DROracle *OCR2DROracleTransactorSession) FulfillRequest(requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _OCR2DROracle.Contract.FulfillRequest(&_OCR2DROracle.TransactOpts, requestId, response, err)
}

func (_OCR2DROracle *OCR2DROracleTransactor) SendRequest(opts *bind.TransactOpts, subscriptionId *big.Int, data []byte) (*types.Transaction, error) {
	return _OCR2DROracle.contract.Transact(opts, "sendRequest", subscriptionId, data)
}

func (_OCR2DROracle *OCR2DROracleSession) SendRequest(subscriptionId *big.Int, data []byte) (*types.Transaction, error) {
	return _OCR2DROracle.Contract.SendRequest(&_OCR2DROracle.TransactOpts, subscriptionId, data)
}

func (_OCR2DROracle *OCR2DROracleTransactorSession) SendRequest(subscriptionId *big.Int, data []byte) (*types.Transaction, error) {
	return _OCR2DROracle.Contract.SendRequest(&_OCR2DROracle.TransactOpts, subscriptionId, data)
}

func (_OCR2DROracle *OCR2DROracleTransactor) SetAuthorizedSenders(opts *bind.TransactOpts, senders []common.Address) (*types.Transaction, error) {
	return _OCR2DROracle.contract.Transact(opts, "setAuthorizedSenders", senders)
}

func (_OCR2DROracle *OCR2DROracleSession) SetAuthorizedSenders(senders []common.Address) (*types.Transaction, error) {
	return _OCR2DROracle.Contract.SetAuthorizedSenders(&_OCR2DROracle.TransactOpts, senders)
}

func (_OCR2DROracle *OCR2DROracleTransactorSession) SetAuthorizedSenders(senders []common.Address) (*types.Transaction, error) {
	return _OCR2DROracle.Contract.SetAuthorizedSenders(&_OCR2DROracle.TransactOpts, senders)
}

func (_OCR2DROracle *OCR2DROracleTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _OCR2DROracle.contract.Transact(opts, "transferOwnership", to)
}

func (_OCR2DROracle *OCR2DROracleSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OCR2DROracle.Contract.TransferOwnership(&_OCR2DROracle.TransactOpts, to)
}

func (_OCR2DROracle *OCR2DROracleTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OCR2DROracle.Contract.TransferOwnership(&_OCR2DROracle.TransactOpts, to)
}

type OCR2DROracleAuthorizedSendersChangedIterator struct {
	Event *OCR2DROracleAuthorizedSendersChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DROracleAuthorizedSendersChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DROracleAuthorizedSendersChanged)
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
		it.Event = new(OCR2DROracleAuthorizedSendersChanged)
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

func (it *OCR2DROracleAuthorizedSendersChangedIterator) Error() error {
	return it.fail
}

func (it *OCR2DROracleAuthorizedSendersChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DROracleAuthorizedSendersChanged struct {
	Senders   []common.Address
	ChangedBy common.Address
	Raw       types.Log
}

func (_OCR2DROracle *OCR2DROracleFilterer) FilterAuthorizedSendersChanged(opts *bind.FilterOpts) (*OCR2DROracleAuthorizedSendersChangedIterator, error) {

	logs, sub, err := _OCR2DROracle.contract.FilterLogs(opts, "AuthorizedSendersChanged")
	if err != nil {
		return nil, err
	}
	return &OCR2DROracleAuthorizedSendersChangedIterator{contract: _OCR2DROracle.contract, event: "AuthorizedSendersChanged", logs: logs, sub: sub}, nil
}

func (_OCR2DROracle *OCR2DROracleFilterer) WatchAuthorizedSendersChanged(opts *bind.WatchOpts, sink chan<- *OCR2DROracleAuthorizedSendersChanged) (event.Subscription, error) {

	logs, sub, err := _OCR2DROracle.contract.WatchLogs(opts, "AuthorizedSendersChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DROracleAuthorizedSendersChanged)
				if err := _OCR2DROracle.contract.UnpackLog(event, "AuthorizedSendersChanged", log); err != nil {
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

func (_OCR2DROracle *OCR2DROracleFilterer) ParseAuthorizedSendersChanged(log types.Log) (*OCR2DROracleAuthorizedSendersChanged, error) {
	event := new(OCR2DROracleAuthorizedSendersChanged)
	if err := _OCR2DROracle.contract.UnpackLog(event, "AuthorizedSendersChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DROracleOracleRequestIterator struct {
	Event *OCR2DROracleOracleRequest

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DROracleOracleRequestIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DROracleOracleRequest)
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
		it.Event = new(OCR2DROracleOracleRequest)
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

func (it *OCR2DROracleOracleRequestIterator) Error() error {
	return it.fail
}

func (it *OCR2DROracleOracleRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DROracleOracleRequest struct {
	RequestId [32]byte
	Data      []byte
	Raw       types.Log
}

func (_OCR2DROracle *OCR2DROracleFilterer) FilterOracleRequest(opts *bind.FilterOpts) (*OCR2DROracleOracleRequestIterator, error) {

	logs, sub, err := _OCR2DROracle.contract.FilterLogs(opts, "OracleRequest")
	if err != nil {
		return nil, err
	}
	return &OCR2DROracleOracleRequestIterator{contract: _OCR2DROracle.contract, event: "OracleRequest", logs: logs, sub: sub}, nil
}

func (_OCR2DROracle *OCR2DROracleFilterer) WatchOracleRequest(opts *bind.WatchOpts, sink chan<- *OCR2DROracleOracleRequest) (event.Subscription, error) {

	logs, sub, err := _OCR2DROracle.contract.WatchLogs(opts, "OracleRequest")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DROracleOracleRequest)
				if err := _OCR2DROracle.contract.UnpackLog(event, "OracleRequest", log); err != nil {
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

func (_OCR2DROracle *OCR2DROracleFilterer) ParseOracleRequest(log types.Log) (*OCR2DROracleOracleRequest, error) {
	event := new(OCR2DROracleOracleRequest)
	if err := _OCR2DROracle.contract.UnpackLog(event, "OracleRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DROracleOracleResponseIterator struct {
	Event *OCR2DROracleOracleResponse

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DROracleOracleResponseIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DROracleOracleResponse)
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
		it.Event = new(OCR2DROracleOracleResponse)
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

func (it *OCR2DROracleOracleResponseIterator) Error() error {
	return it.fail
}

func (it *OCR2DROracleOracleResponseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DROracleOracleResponse struct {
	RequestId [32]byte
	Raw       types.Log
}

func (_OCR2DROracle *OCR2DROracleFilterer) FilterOracleResponse(opts *bind.FilterOpts) (*OCR2DROracleOracleResponseIterator, error) {

	logs, sub, err := _OCR2DROracle.contract.FilterLogs(opts, "OracleResponse")
	if err != nil {
		return nil, err
	}
	return &OCR2DROracleOracleResponseIterator{contract: _OCR2DROracle.contract, event: "OracleResponse", logs: logs, sub: sub}, nil
}

func (_OCR2DROracle *OCR2DROracleFilterer) WatchOracleResponse(opts *bind.WatchOpts, sink chan<- *OCR2DROracleOracleResponse) (event.Subscription, error) {

	logs, sub, err := _OCR2DROracle.contract.WatchLogs(opts, "OracleResponse")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DROracleOracleResponse)
				if err := _OCR2DROracle.contract.UnpackLog(event, "OracleResponse", log); err != nil {
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

func (_OCR2DROracle *OCR2DROracleFilterer) ParseOracleResponse(log types.Log) (*OCR2DROracleOracleResponse, error) {
	event := new(OCR2DROracleOracleResponse)
	if err := _OCR2DROracle.contract.UnpackLog(event, "OracleResponse", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DROracleOwnershipTransferRequestedIterator struct {
	Event *OCR2DROracleOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DROracleOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DROracleOwnershipTransferRequested)
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
		it.Event = new(OCR2DROracleOwnershipTransferRequested)
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

func (it *OCR2DROracleOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *OCR2DROracleOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DROracleOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OCR2DROracle *OCR2DROracleFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2DROracleOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2DROracle.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OCR2DROracleOwnershipTransferRequestedIterator{contract: _OCR2DROracle.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_OCR2DROracle *OCR2DROracleFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OCR2DROracleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2DROracle.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DROracleOwnershipTransferRequested)
				if err := _OCR2DROracle.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_OCR2DROracle *OCR2DROracleFilterer) ParseOwnershipTransferRequested(log types.Log) (*OCR2DROracleOwnershipTransferRequested, error) {
	event := new(OCR2DROracleOwnershipTransferRequested)
	if err := _OCR2DROracle.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2DROracleOwnershipTransferredIterator struct {
	Event *OCR2DROracleOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2DROracleOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2DROracleOwnershipTransferred)
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
		it.Event = new(OCR2DROracleOwnershipTransferred)
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

func (it *OCR2DROracleOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *OCR2DROracleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2DROracleOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OCR2DROracle *OCR2DROracleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2DROracleOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2DROracle.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OCR2DROracleOwnershipTransferredIterator{contract: _OCR2DROracle.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_OCR2DROracle *OCR2DROracleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OCR2DROracleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2DROracle.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2DROracleOwnershipTransferred)
				if err := _OCR2DROracle.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_OCR2DROracle *OCR2DROracleFilterer) ParseOwnershipTransferred(log types.Log) (*OCR2DROracleOwnershipTransferred, error) {
	event := new(OCR2DROracleOwnershipTransferred)
	if err := _OCR2DROracle.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_OCR2DROracle *OCR2DROracle) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _OCR2DROracle.abi.Events["AuthorizedSendersChanged"].ID:
		return _OCR2DROracle.ParseAuthorizedSendersChanged(log)
	case _OCR2DROracle.abi.Events["OracleRequest"].ID:
		return _OCR2DROracle.ParseOracleRequest(log)
	case _OCR2DROracle.abi.Events["OracleResponse"].ID:
		return _OCR2DROracle.ParseOracleResponse(log)
	case _OCR2DROracle.abi.Events["OwnershipTransferRequested"].ID:
		return _OCR2DROracle.ParseOwnershipTransferRequested(log)
	case _OCR2DROracle.abi.Events["OwnershipTransferred"].ID:
		return _OCR2DROracle.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (OCR2DROracleAuthorizedSendersChanged) Topic() common.Hash {
	return common.HexToHash("0xf263cfb3e4298332e776194610cf9fdc09ccb3ada8b9aa39764d882e11fbf0a0")
}

func (OCR2DROracleOracleRequest) Topic() common.Hash {
	return common.HexToHash("0x9dc59d1e6d6042f6c2d2af5a3f9f6502bf0abf7cf40f0eb80edf752ca396f2ba")
}

func (OCR2DROracleOracleResponse) Topic() common.Hash {
	return common.HexToHash("0x9e9bc7616d42c2835d05ae617e508454e63b30b934be8aa932ebc125e0e58a64")
}

func (OCR2DROracleOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (OCR2DROracleOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_OCR2DROracle *OCR2DROracle) Address() common.Address {
	return _OCR2DROracle.address
}

type OCR2DROracleInterface interface {
	GetAuthorizedSenders(opts *bind.CallOpts) ([]common.Address, error)

	IsAuthorizedSender(opts *bind.CallOpts, sender common.Address) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	FulfillRequest(opts *bind.TransactOpts, requestId [32]byte, response []byte, err []byte) (*types.Transaction, error)

	SendRequest(opts *bind.TransactOpts, subscriptionId *big.Int, data []byte) (*types.Transaction, error)

	SetAuthorizedSenders(opts *bind.TransactOpts, senders []common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterAuthorizedSendersChanged(opts *bind.FilterOpts) (*OCR2DROracleAuthorizedSendersChangedIterator, error)

	WatchAuthorizedSendersChanged(opts *bind.WatchOpts, sink chan<- *OCR2DROracleAuthorizedSendersChanged) (event.Subscription, error)

	ParseAuthorizedSendersChanged(log types.Log) (*OCR2DROracleAuthorizedSendersChanged, error)

	FilterOracleRequest(opts *bind.FilterOpts) (*OCR2DROracleOracleRequestIterator, error)

	WatchOracleRequest(opts *bind.WatchOpts, sink chan<- *OCR2DROracleOracleRequest) (event.Subscription, error)

	ParseOracleRequest(log types.Log) (*OCR2DROracleOracleRequest, error)

	FilterOracleResponse(opts *bind.FilterOpts) (*OCR2DROracleOracleResponseIterator, error)

	WatchOracleResponse(opts *bind.WatchOpts, sink chan<- *OCR2DROracleOracleResponse) (event.Subscription, error)

	ParseOracleResponse(log types.Log) (*OCR2DROracleOracleResponse, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2DROracleOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OCR2DROracleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*OCR2DROracleOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2DROracleOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OCR2DROracleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*OCR2DROracleOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
