// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethereum

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

// Reference imports to suppress errors if they are not otherwise used.
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

// APIConsumerMetaData contains all meta data concerning the APIConsumer contract.
var APIConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ChainlinkCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ChainlinkFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ChainlinkRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"roundID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"PerfMetricsEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_payment\",\"type\":\"uint256\"},{\"internalType\":\"bytes4\",\"name\":\"_callbackFunctionId\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"_expiration\",\"type\":\"uint256\"}],\"name\":\"cancelRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_jobId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_payment\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_url\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_path\",\"type\":\"string\"},{\"internalType\":\"int256\",\"name\":\"_times\",\"type\":\"int256\"}],\"name\":\"createRequestTo\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentRoundID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"data\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_data\",\"type\":\"uint256\"}],\"name\":\"fulfill\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChainlinkToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"selector\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040526001600455600060075534801561001a57600080fd5b50604051620012d2380380620012d28339818101604052602081101561003f57600080fd5b5051600680546001600160a01b0319163317908190556040516001600160a01b0391909116906000907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a36001600160a01b0381166100b1576100ac6001600160e01b036100c916565b6100c3565b6100c3816001600160e01b0361015516565b50610177565b61015373c89bd4e1632d3a43cb03aaad5262cbe4038bc5716001600160a01b03166338cc48316040518163ffffffff1660e01b815260040160206040518083038186803b15801561011957600080fd5b505afa15801561012d573d6000803e3d6000fd5b505050506040513d602081101561014357600080fd5b50516001600160e01b0361015516565b565b600280546001600160a01b0319166001600160a01b0392909216919091179055565b61114b80620001876000396000f3fe608060405234801561001057600080fd5b50600436106100995760003560e01c8063165d35e11461009e57806316ef7f1a146100c25780634357855e1461021b57806373d4a13a146102405780638da5cb5b146102485780638dc654a2146102505780638f32d59b14610258578063a312c4f214610274578063ea3d508a1461027c578063ec65d0f8146102a1578063f2fde38b146102da575b600080fd5b6100a6610300565b604080516001600160a01b039092168252519081900360200190f35b610209600480360360c08110156100d857600080fd5b6001600160a01b038235169160208101359160408201359190810190608081016060820135600160201b81111561010e57600080fd5b82018360208201111561012057600080fd5b803590602001918460018302840111600160201b8311171561014157600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295949360208101935035915050600160201b81111561019357600080fd5b8201836020820111156101a557600080fd5b803590602001918460018302840111600160201b831117156101c657600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550509135925061030f915050565b60408051918252519081900360200190f35b61023e6004803603604081101561023157600080fd5b5080359060200135610426565b005b610209610478565b6100a661047e565b61023e61048d565b610260610621565b604080519115158252519081900360200190f35b610209610632565b610284610638565b604080516001600160e01b03199092168252519081900360200190f35b61023e600480360360808110156102b757600080fd5b508035906020810135906001600160e01b03196040820135169060600135610641565b61023e600480360360208110156102f057600080fd5b50356001600160a01b031661069a565b600061030a6106ea565b905090565b6000610319610621565b610358576040805162461bcd60e51b815260206004820181905260248201526000805160206110f6833981519152604482015290519081900360640190fd5b6009805463ffffffff1916634357855e17905561037361105d565b61038587306321abc2af60e11b6106f9565b60408051808201909152600381526219d95d60ea1b60208201529091506103b49082908763ffffffff61072416565b6040805180820190915260048152630e0c2e8d60e31b60208201526103e19082908663ffffffff61072416565b60408051808201909152600581526474696d657360d81b602082015261040f9082908563ffffffff61075316565b61041a88828861077d565b98975050505050505050565b6008819055600780546001019081905560408051918252602082018490524282820152517ffbaf68ee7b9032982942607eaea1859969ed8674797b5c2fc6fecaa7538519469181900360600190a15050565b60085481565b6006546001600160a01b031690565b610495610621565b6104d4576040805162461bcd60e51b815260206004820181905260248201526000805160206110f6833981519152604482015290519081900360640190fd5b60006104de6106ea565b604080516370a0823160e01b815230600482015290519192506001600160a01b0383169163a9059cbb91339184916370a08231916024808301926020929190829003018186803b15801561053157600080fd5b505afa158015610545573d6000803e3d6000fd5b505050506040513d602081101561055b57600080fd5b5051604080516001600160e01b031960e086901b1681526001600160a01b03909316600484015260248301919091525160448083019260209291908290030181600087803b1580156105ac57600080fd5b505af11580156105c0573d6000803e3d6000fd5b505050506040513d60208110156105d657600080fd5b505161061e576040805162461bcd60e51b81526020600482015260126024820152712ab730b13632903a37903a3930b739b332b960711b604482015290519081900360640190fd5b50565b6006546001600160a01b0316331490565b60075481565b60095460e01b81565b610649610621565b610688576040805162461bcd60e51b815260206004820181905260248201526000805160206110f6833981519152604482015290519081900360640190fd5b61069484848484610954565b50505050565b6106a2610621565b6106e1576040805162461bcd60e51b815260206004820181905260248201526000805160206110f6833981519152604482015290519081900360640190fd5b61061e81610a2c565b6002546001600160a01b031690565b61070161105d565b61070961105d565b61071b8186868663ffffffff610acd16565b95945050505050565b6080830151610739908363ffffffff610b0a16565b608083015161074e908263ffffffff610b0a16565b505050565b6080830151610768908363ffffffff610b0a16565b608083015161074e908263ffffffff610b2716565b6004546040805130606090811b60208084019190915260348084018690528451808503909101815260549093018452825192810192909220908601939093526000838152600590915281812080546001600160a01b0319166001600160a01b038816179055905182917fb5e6e01e79f91267dc17b4e6314d5d4d03593d2ceee0fbb452b750bd70ea5af991a26002546001600160a01b0316634000aea0858461082587610b86565b6040518463ffffffff1660e01b815260040180846001600160a01b03166001600160a01b0316815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b8381101561088f578181015183820152602001610877565b50505050905090810190601f1680156108bc5780820380516001836020036101000a031916815260200191505b50945050505050602060405180830381600087803b1580156108dd57600080fd5b505af11580156108f1573d6000803e3d6000fd5b505050506040513d602081101561090757600080fd5b50516109445760405162461bcd60e51b81526004018080602001828103825260238152602001806110d36023913960400191505060405180910390fd5b6004805460010190559392505050565b60008481526005602052604080822080546001600160a01b0319811690915590516001600160a01b039091169186917fe1fe3afa0f7f761ff0a8b89086790efd5140d2907ebd5b7ff6bfcb5e075fd4c59190a260408051636ee4d55360e01b815260048101879052602481018690526001600160e01b0319851660448201526064810184905290516001600160a01b03831691636ee4d55391608480830192600092919082900301818387803b158015610a0d57600080fd5b505af1158015610a21573d6000803e3d6000fd5b505050505050505050565b6001600160a01b038116610a715760405162461bcd60e51b81526004018080602001828103825260268152602001806110ad6026913960400191505060405180910390fd5b6006546040516001600160a01b038084169216907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a3600680546001600160a01b0319166001600160a01b0392909216919091179055565b610ad561105d565b610ae58560800151610100610caf565b50509183526001600160a01b031660208301526001600160e01b031916604082015290565b610b178260038351610cef565b61074e828263ffffffff610df916565b67ffffffffffffffff19811215610b4757610b428282610e1a565b610b82565b67ffffffffffffffff811315610b6157610b428282610e59565b60008112610b7557610b4282600083610cef565b610b828260018319610cef565b5050565b8051602080830151604080850151606086810151608088015151935160006024820181815260448301829052606483018a90526001600160a01b03881660848401526001600160e01b0319861660a484015260c48301849052600160e48401819052610100610104850190815288516101248601528851969b6320214ca360e11b9b949a8b9a91999098909796939591949361014401918501908083838e5b83811015610c3d578181015183820152602001610c25565b50505050905090810190601f168015610c6a5780820380516001836020036101000a031916815260200191505b5060408051601f198184030181529190526020810180516001600160e01b03166001600160e01b0319909d169c909c17909b5250989950505050505050505050919050565b610cb7611092565b6020820615610ccc5760208206602003820191505b506020808301829052604080518085526000815283019091019052815b92915050565b60178111610d1657610d108360e0600585901b16831763ffffffff610e9416565b5061074e565b60ff8111610d4c57610d39836018611fe0600586901b161763ffffffff610e9416565b50610d108382600163ffffffff610eac16565b61ffff8111610d8357610d70836019611fe0600586901b161763ffffffff610e9416565b50610d108382600263ffffffff610eac16565b63ffffffff8111610dbc57610da983601a611fe0600586901b161763ffffffff610e9416565b50610d108382600463ffffffff610eac16565b67ffffffffffffffff811161074e57610de683601b611fe0600586901b161763ffffffff610e9416565b506106948382600863ffffffff610eac16565b610e01611092565b610e1383846000015151848551610ecd565b9392505050565b610e2b8260c363ffffffff610e9416565b50610b8282826000190360405160200180828152602001915050604051602081830303815290604052610f79565b610e6a8260c263ffffffff610e9416565b50610b82828260405160200180828152602001915050604051602081830303815290604052610f79565b610e9c611092565b610e138384600001515184610f86565b610eb4611092565b610ec5848560000151518585610fd1565b949350505050565b610ed5611092565b8251821115610ee357600080fd5b84602001518285011115610f0d57610f0d85610f05876020015187860161102f565b600202611046565b600080865180518760208301019350808887011115610f2c5787860182525b505050602084015b60208410610f535780518252601f199093019260209182019101610f34565b51815160001960208690036101000a019081169019919091161790525083949350505050565b610b178260028351610cef565b610f8e611092565b83602001518310610faa57610faa848560200151600202611046565b835180516020858301018481535080851415610fc7576001810182525b5093949350505050565b610fd9611092565b84602001518483011115610ff657610ff685858401600202611046565b60006001836101000a0390508551838682010185831982511617815250805184870111156110245783860181525b509495945050505050565b600081831115611040575081610ce9565b50919050565b81516110528383610caf565b506106948382610df9565b6040805160a08101825260008082526020820181905291810182905260608101919091526080810161108d611092565b905290565b60405180604001604052806060815260200160008152509056fe4f776e61626c653a206e6577206f776e657220697320746865207a65726f2061646472657373756e61626c6520746f207472616e73666572416e6443616c6c20746f206f7261636c654f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572a26469706673582212203d68b6d8e7f8e263b1dd29932d411a9bbd292574f49c5e0d76abaece361fa56864736f6c63430006000033",
}

// APIConsumerABI is the input ABI used to generate the binding from.
// Deprecated: Use APIConsumerMetaData.ABI instead.
var APIConsumerABI = APIConsumerMetaData.ABI

// APIConsumerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use APIConsumerMetaData.Bin instead.
var APIConsumerBin = APIConsumerMetaData.Bin

// DeployAPIConsumer deploys a new Ethereum contract, binding an instance of APIConsumer to it.
func DeployAPIConsumer(auth *bind.TransactOpts, backend bind.ContractBackend, _link common.Address) (common.Address, *types.Transaction, *APIConsumer, error) {
	parsed, err := APIConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(APIConsumerBin), backend, _link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &APIConsumer{APIConsumerCaller: APIConsumerCaller{contract: contract}, APIConsumerTransactor: APIConsumerTransactor{contract: contract}, APIConsumerFilterer: APIConsumerFilterer{contract: contract}}, nil
}

// APIConsumer is an auto generated Go binding around an Ethereum contract.
type APIConsumer struct {
	APIConsumerCaller     // Read-only binding to the contract
	APIConsumerTransactor // Write-only binding to the contract
	APIConsumerFilterer   // Log filterer for contract events
}

// APIConsumerCaller is an auto generated read-only Go binding around an Ethereum contract.
type APIConsumerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// APIConsumerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type APIConsumerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// APIConsumerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type APIConsumerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// APIConsumerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type APIConsumerSession struct {
	Contract     *APIConsumer      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// APIConsumerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type APIConsumerCallerSession struct {
	Contract *APIConsumerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// APIConsumerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type APIConsumerTransactorSession struct {
	Contract     *APIConsumerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// APIConsumerRaw is an auto generated low-level Go binding around an Ethereum contract.
type APIConsumerRaw struct {
	Contract *APIConsumer // Generic contract binding to access the raw methods on
}

// APIConsumerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type APIConsumerCallerRaw struct {
	Contract *APIConsumerCaller // Generic read-only contract binding to access the raw methods on
}

// APIConsumerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type APIConsumerTransactorRaw struct {
	Contract *APIConsumerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAPIConsumer creates a new instance of APIConsumer, bound to a specific deployed contract.
func NewAPIConsumer(address common.Address, backend bind.ContractBackend) (*APIConsumer, error) {
	contract, err := bindAPIConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &APIConsumer{APIConsumerCaller: APIConsumerCaller{contract: contract}, APIConsumerTransactor: APIConsumerTransactor{contract: contract}, APIConsumerFilterer: APIConsumerFilterer{contract: contract}}, nil
}

// NewAPIConsumerCaller creates a new read-only instance of APIConsumer, bound to a specific deployed contract.
func NewAPIConsumerCaller(address common.Address, caller bind.ContractCaller) (*APIConsumerCaller, error) {
	contract, err := bindAPIConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &APIConsumerCaller{contract: contract}, nil
}

// NewAPIConsumerTransactor creates a new write-only instance of APIConsumer, bound to a specific deployed contract.
func NewAPIConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*APIConsumerTransactor, error) {
	contract, err := bindAPIConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &APIConsumerTransactor{contract: contract}, nil
}

// NewAPIConsumerFilterer creates a new log filterer instance of APIConsumer, bound to a specific deployed contract.
func NewAPIConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*APIConsumerFilterer, error) {
	contract, err := bindAPIConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &APIConsumerFilterer{contract: contract}, nil
}

// bindAPIConsumer binds a generic wrapper to an already deployed contract.
func bindAPIConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(APIConsumerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_APIConsumer *APIConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _APIConsumer.Contract.APIConsumerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_APIConsumer *APIConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _APIConsumer.Contract.APIConsumerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_APIConsumer *APIConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _APIConsumer.Contract.APIConsumerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_APIConsumer *APIConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _APIConsumer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_APIConsumer *APIConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _APIConsumer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_APIConsumer *APIConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _APIConsumer.Contract.contract.Transact(opts, method, params...)
}

// CurrentRoundID is a free data retrieval call binding the contract method 0xa312c4f2.
//
// Solidity: function currentRoundID() view returns(uint256)
func (_APIConsumer *APIConsumerCaller) CurrentRoundID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _APIConsumer.contract.Call(opts, &out, "currentRoundID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentRoundID is a free data retrieval call binding the contract method 0xa312c4f2.
//
// Solidity: function currentRoundID() view returns(uint256)
func (_APIConsumer *APIConsumerSession) CurrentRoundID() (*big.Int, error) {
	return _APIConsumer.Contract.CurrentRoundID(&_APIConsumer.CallOpts)
}

// CurrentRoundID is a free data retrieval call binding the contract method 0xa312c4f2.
//
// Solidity: function currentRoundID() view returns(uint256)
func (_APIConsumer *APIConsumerCallerSession) CurrentRoundID() (*big.Int, error) {
	return _APIConsumer.Contract.CurrentRoundID(&_APIConsumer.CallOpts)
}

// Data is a free data retrieval call binding the contract method 0x73d4a13a.
//
// Solidity: function data() view returns(uint256)
func (_APIConsumer *APIConsumerCaller) Data(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _APIConsumer.contract.Call(opts, &out, "data")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Data is a free data retrieval call binding the contract method 0x73d4a13a.
//
// Solidity: function data() view returns(uint256)
func (_APIConsumer *APIConsumerSession) Data() (*big.Int, error) {
	return _APIConsumer.Contract.Data(&_APIConsumer.CallOpts)
}

// Data is a free data retrieval call binding the contract method 0x73d4a13a.
//
// Solidity: function data() view returns(uint256)
func (_APIConsumer *APIConsumerCallerSession) Data() (*big.Int, error) {
	return _APIConsumer.Contract.Data(&_APIConsumer.CallOpts)
}

// GetChainlinkToken is a free data retrieval call binding the contract method 0x165d35e1.
//
// Solidity: function getChainlinkToken() view returns(address)
func (_APIConsumer *APIConsumerCaller) GetChainlinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _APIConsumer.contract.Call(opts, &out, "getChainlinkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetChainlinkToken is a free data retrieval call binding the contract method 0x165d35e1.
//
// Solidity: function getChainlinkToken() view returns(address)
func (_APIConsumer *APIConsumerSession) GetChainlinkToken() (common.Address, error) {
	return _APIConsumer.Contract.GetChainlinkToken(&_APIConsumer.CallOpts)
}

// GetChainlinkToken is a free data retrieval call binding the contract method 0x165d35e1.
//
// Solidity: function getChainlinkToken() view returns(address)
func (_APIConsumer *APIConsumerCallerSession) GetChainlinkToken() (common.Address, error) {
	return _APIConsumer.Contract.GetChainlinkToken(&_APIConsumer.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_APIConsumer *APIConsumerCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _APIConsumer.contract.Call(opts, &out, "isOwner")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_APIConsumer *APIConsumerSession) IsOwner() (bool, error) {
	return _APIConsumer.Contract.IsOwner(&_APIConsumer.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_APIConsumer *APIConsumerCallerSession) IsOwner() (bool, error) {
	return _APIConsumer.Contract.IsOwner(&_APIConsumer.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_APIConsumer *APIConsumerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _APIConsumer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_APIConsumer *APIConsumerSession) Owner() (common.Address, error) {
	return _APIConsumer.Contract.Owner(&_APIConsumer.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_APIConsumer *APIConsumerCallerSession) Owner() (common.Address, error) {
	return _APIConsumer.Contract.Owner(&_APIConsumer.CallOpts)
}

// Selector is a free data retrieval call binding the contract method 0xea3d508a.
//
// Solidity: function selector() view returns(bytes4)
func (_APIConsumer *APIConsumerCaller) Selector(opts *bind.CallOpts) ([4]byte, error) {
	var out []interface{}
	err := _APIConsumer.contract.Call(opts, &out, "selector")

	if err != nil {
		return *new([4]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([4]byte)).(*[4]byte)

	return out0, err

}

// Selector is a free data retrieval call binding the contract method 0xea3d508a.
//
// Solidity: function selector() view returns(bytes4)
func (_APIConsumer *APIConsumerSession) Selector() ([4]byte, error) {
	return _APIConsumer.Contract.Selector(&_APIConsumer.CallOpts)
}

// Selector is a free data retrieval call binding the contract method 0xea3d508a.
//
// Solidity: function selector() view returns(bytes4)
func (_APIConsumer *APIConsumerCallerSession) Selector() ([4]byte, error) {
	return _APIConsumer.Contract.Selector(&_APIConsumer.CallOpts)
}

// CancelRequest is a paid mutator transaction binding the contract method 0xec65d0f8.
//
// Solidity: function cancelRequest(bytes32 _requestId, uint256 _payment, bytes4 _callbackFunctionId, uint256 _expiration) returns()
func (_APIConsumer *APIConsumerTransactor) CancelRequest(opts *bind.TransactOpts, _requestId [32]byte, _payment *big.Int, _callbackFunctionId [4]byte, _expiration *big.Int) (*types.Transaction, error) {
	return _APIConsumer.contract.Transact(opts, "cancelRequest", _requestId, _payment, _callbackFunctionId, _expiration)
}

// CancelRequest is a paid mutator transaction binding the contract method 0xec65d0f8.
//
// Solidity: function cancelRequest(bytes32 _requestId, uint256 _payment, bytes4 _callbackFunctionId, uint256 _expiration) returns()
func (_APIConsumer *APIConsumerSession) CancelRequest(_requestId [32]byte, _payment *big.Int, _callbackFunctionId [4]byte, _expiration *big.Int) (*types.Transaction, error) {
	return _APIConsumer.Contract.CancelRequest(&_APIConsumer.TransactOpts, _requestId, _payment, _callbackFunctionId, _expiration)
}

// CancelRequest is a paid mutator transaction binding the contract method 0xec65d0f8.
//
// Solidity: function cancelRequest(bytes32 _requestId, uint256 _payment, bytes4 _callbackFunctionId, uint256 _expiration) returns()
func (_APIConsumer *APIConsumerTransactorSession) CancelRequest(_requestId [32]byte, _payment *big.Int, _callbackFunctionId [4]byte, _expiration *big.Int) (*types.Transaction, error) {
	return _APIConsumer.Contract.CancelRequest(&_APIConsumer.TransactOpts, _requestId, _payment, _callbackFunctionId, _expiration)
}

// CreateRequestTo is a paid mutator transaction binding the contract method 0x16ef7f1a.
//
// Solidity: function createRequestTo(address _oracle, bytes32 _jobId, uint256 _payment, string _url, string _path, int256 _times) returns(bytes32 requestId)
func (_APIConsumer *APIConsumerTransactor) CreateRequestTo(opts *bind.TransactOpts, _oracle common.Address, _jobId [32]byte, _payment *big.Int, _url string, _path string, _times *big.Int) (*types.Transaction, error) {
	return _APIConsumer.contract.Transact(opts, "createRequestTo", _oracle, _jobId, _payment, _url, _path, _times)
}

// CreateRequestTo is a paid mutator transaction binding the contract method 0x16ef7f1a.
//
// Solidity: function createRequestTo(address _oracle, bytes32 _jobId, uint256 _payment, string _url, string _path, int256 _times) returns(bytes32 requestId)
func (_APIConsumer *APIConsumerSession) CreateRequestTo(_oracle common.Address, _jobId [32]byte, _payment *big.Int, _url string, _path string, _times *big.Int) (*types.Transaction, error) {
	return _APIConsumer.Contract.CreateRequestTo(&_APIConsumer.TransactOpts, _oracle, _jobId, _payment, _url, _path, _times)
}

// CreateRequestTo is a paid mutator transaction binding the contract method 0x16ef7f1a.
//
// Solidity: function createRequestTo(address _oracle, bytes32 _jobId, uint256 _payment, string _url, string _path, int256 _times) returns(bytes32 requestId)
func (_APIConsumer *APIConsumerTransactorSession) CreateRequestTo(_oracle common.Address, _jobId [32]byte, _payment *big.Int, _url string, _path string, _times *big.Int) (*types.Transaction, error) {
	return _APIConsumer.Contract.CreateRequestTo(&_APIConsumer.TransactOpts, _oracle, _jobId, _payment, _url, _path, _times)
}

// Fulfill is a paid mutator transaction binding the contract method 0x4357855e.
//
// Solidity: function fulfill(bytes32 _requestId, uint256 _data) returns()
func (_APIConsumer *APIConsumerTransactor) Fulfill(opts *bind.TransactOpts, _requestId [32]byte, _data *big.Int) (*types.Transaction, error) {
	return _APIConsumer.contract.Transact(opts, "fulfill", _requestId, _data)
}

// Fulfill is a paid mutator transaction binding the contract method 0x4357855e.
//
// Solidity: function fulfill(bytes32 _requestId, uint256 _data) returns()
func (_APIConsumer *APIConsumerSession) Fulfill(_requestId [32]byte, _data *big.Int) (*types.Transaction, error) {
	return _APIConsumer.Contract.Fulfill(&_APIConsumer.TransactOpts, _requestId, _data)
}

// Fulfill is a paid mutator transaction binding the contract method 0x4357855e.
//
// Solidity: function fulfill(bytes32 _requestId, uint256 _data) returns()
func (_APIConsumer *APIConsumerTransactorSession) Fulfill(_requestId [32]byte, _data *big.Int) (*types.Transaction, error) {
	return _APIConsumer.Contract.Fulfill(&_APIConsumer.TransactOpts, _requestId, _data)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_APIConsumer *APIConsumerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _APIConsumer.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_APIConsumer *APIConsumerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _APIConsumer.Contract.TransferOwnership(&_APIConsumer.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_APIConsumer *APIConsumerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _APIConsumer.Contract.TransferOwnership(&_APIConsumer.TransactOpts, newOwner)
}

// WithdrawLink is a paid mutator transaction binding the contract method 0x8dc654a2.
//
// Solidity: function withdrawLink() returns()
func (_APIConsumer *APIConsumerTransactor) WithdrawLink(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _APIConsumer.contract.Transact(opts, "withdrawLink")
}

// WithdrawLink is a paid mutator transaction binding the contract method 0x8dc654a2.
//
// Solidity: function withdrawLink() returns()
func (_APIConsumer *APIConsumerSession) WithdrawLink() (*types.Transaction, error) {
	return _APIConsumer.Contract.WithdrawLink(&_APIConsumer.TransactOpts)
}

// WithdrawLink is a paid mutator transaction binding the contract method 0x8dc654a2.
//
// Solidity: function withdrawLink() returns()
func (_APIConsumer *APIConsumerTransactorSession) WithdrawLink() (*types.Transaction, error) {
	return _APIConsumer.Contract.WithdrawLink(&_APIConsumer.TransactOpts)
}

// APIConsumerChainlinkCancelledIterator is returned from FilterChainlinkCancelled and is used to iterate over the raw logs and unpacked data for ChainlinkCancelled events raised by the APIConsumer contract.
type APIConsumerChainlinkCancelledIterator struct {
	Event *APIConsumerChainlinkCancelled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *APIConsumerChainlinkCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(APIConsumerChainlinkCancelled)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(APIConsumerChainlinkCancelled)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *APIConsumerChainlinkCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *APIConsumerChainlinkCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// APIConsumerChainlinkCancelled represents a ChainlinkCancelled event raised by the APIConsumer contract.
type APIConsumerChainlinkCancelled struct {
	Id  [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterChainlinkCancelled is a free log retrieval operation binding the contract event 0xe1fe3afa0f7f761ff0a8b89086790efd5140d2907ebd5b7ff6bfcb5e075fd4c5.
//
// Solidity: event ChainlinkCancelled(bytes32 indexed id)
func (_APIConsumer *APIConsumerFilterer) FilterChainlinkCancelled(opts *bind.FilterOpts, id [][32]byte) (*APIConsumerChainlinkCancelledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _APIConsumer.contract.FilterLogs(opts, "ChainlinkCancelled", idRule)
	if err != nil {
		return nil, err
	}
	return &APIConsumerChainlinkCancelledIterator{contract: _APIConsumer.contract, event: "ChainlinkCancelled", logs: logs, sub: sub}, nil
}

// WatchChainlinkCancelled is a free log subscription operation binding the contract event 0xe1fe3afa0f7f761ff0a8b89086790efd5140d2907ebd5b7ff6bfcb5e075fd4c5.
//
// Solidity: event ChainlinkCancelled(bytes32 indexed id)
func (_APIConsumer *APIConsumerFilterer) WatchChainlinkCancelled(opts *bind.WatchOpts, sink chan<- *APIConsumerChainlinkCancelled, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _APIConsumer.contract.WatchLogs(opts, "ChainlinkCancelled", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(APIConsumerChainlinkCancelled)
				if err := _APIConsumer.contract.UnpackLog(event, "ChainlinkCancelled", log); err != nil {
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

// ParseChainlinkCancelled is a log parse operation binding the contract event 0xe1fe3afa0f7f761ff0a8b89086790efd5140d2907ebd5b7ff6bfcb5e075fd4c5.
//
// Solidity: event ChainlinkCancelled(bytes32 indexed id)
func (_APIConsumer *APIConsumerFilterer) ParseChainlinkCancelled(log types.Log) (*APIConsumerChainlinkCancelled, error) {
	event := new(APIConsumerChainlinkCancelled)
	if err := _APIConsumer.contract.UnpackLog(event, "ChainlinkCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// APIConsumerChainlinkFulfilledIterator is returned from FilterChainlinkFulfilled and is used to iterate over the raw logs and unpacked data for ChainlinkFulfilled events raised by the APIConsumer contract.
type APIConsumerChainlinkFulfilledIterator struct {
	Event *APIConsumerChainlinkFulfilled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *APIConsumerChainlinkFulfilledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(APIConsumerChainlinkFulfilled)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(APIConsumerChainlinkFulfilled)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *APIConsumerChainlinkFulfilledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *APIConsumerChainlinkFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// APIConsumerChainlinkFulfilled represents a ChainlinkFulfilled event raised by the APIConsumer contract.
type APIConsumerChainlinkFulfilled struct {
	Id  [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterChainlinkFulfilled is a free log retrieval operation binding the contract event 0x7cc135e0cebb02c3480ae5d74d377283180a2601f8f644edf7987b009316c63a.
//
// Solidity: event ChainlinkFulfilled(bytes32 indexed id)
func (_APIConsumer *APIConsumerFilterer) FilterChainlinkFulfilled(opts *bind.FilterOpts, id [][32]byte) (*APIConsumerChainlinkFulfilledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _APIConsumer.contract.FilterLogs(opts, "ChainlinkFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return &APIConsumerChainlinkFulfilledIterator{contract: _APIConsumer.contract, event: "ChainlinkFulfilled", logs: logs, sub: sub}, nil
}

// WatchChainlinkFulfilled is a free log subscription operation binding the contract event 0x7cc135e0cebb02c3480ae5d74d377283180a2601f8f644edf7987b009316c63a.
//
// Solidity: event ChainlinkFulfilled(bytes32 indexed id)
func (_APIConsumer *APIConsumerFilterer) WatchChainlinkFulfilled(opts *bind.WatchOpts, sink chan<- *APIConsumerChainlinkFulfilled, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _APIConsumer.contract.WatchLogs(opts, "ChainlinkFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(APIConsumerChainlinkFulfilled)
				if err := _APIConsumer.contract.UnpackLog(event, "ChainlinkFulfilled", log); err != nil {
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

// ParseChainlinkFulfilled is a log parse operation binding the contract event 0x7cc135e0cebb02c3480ae5d74d377283180a2601f8f644edf7987b009316c63a.
//
// Solidity: event ChainlinkFulfilled(bytes32 indexed id)
func (_APIConsumer *APIConsumerFilterer) ParseChainlinkFulfilled(log types.Log) (*APIConsumerChainlinkFulfilled, error) {
	event := new(APIConsumerChainlinkFulfilled)
	if err := _APIConsumer.contract.UnpackLog(event, "ChainlinkFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// APIConsumerChainlinkRequestedIterator is returned from FilterChainlinkRequested and is used to iterate over the raw logs and unpacked data for ChainlinkRequested events raised by the APIConsumer contract.
type APIConsumerChainlinkRequestedIterator struct {
	Event *APIConsumerChainlinkRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *APIConsumerChainlinkRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(APIConsumerChainlinkRequested)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(APIConsumerChainlinkRequested)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *APIConsumerChainlinkRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *APIConsumerChainlinkRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// APIConsumerChainlinkRequested represents a ChainlinkRequested event raised by the APIConsumer contract.
type APIConsumerChainlinkRequested struct {
	Id  [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterChainlinkRequested is a free log retrieval operation binding the contract event 0xb5e6e01e79f91267dc17b4e6314d5d4d03593d2ceee0fbb452b750bd70ea5af9.
//
// Solidity: event ChainlinkRequested(bytes32 indexed id)
func (_APIConsumer *APIConsumerFilterer) FilterChainlinkRequested(opts *bind.FilterOpts, id [][32]byte) (*APIConsumerChainlinkRequestedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _APIConsumer.contract.FilterLogs(opts, "ChainlinkRequested", idRule)
	if err != nil {
		return nil, err
	}
	return &APIConsumerChainlinkRequestedIterator{contract: _APIConsumer.contract, event: "ChainlinkRequested", logs: logs, sub: sub}, nil
}

// WatchChainlinkRequested is a free log subscription operation binding the contract event 0xb5e6e01e79f91267dc17b4e6314d5d4d03593d2ceee0fbb452b750bd70ea5af9.
//
// Solidity: event ChainlinkRequested(bytes32 indexed id)
func (_APIConsumer *APIConsumerFilterer) WatchChainlinkRequested(opts *bind.WatchOpts, sink chan<- *APIConsumerChainlinkRequested, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _APIConsumer.contract.WatchLogs(opts, "ChainlinkRequested", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(APIConsumerChainlinkRequested)
				if err := _APIConsumer.contract.UnpackLog(event, "ChainlinkRequested", log); err != nil {
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

// ParseChainlinkRequested is a log parse operation binding the contract event 0xb5e6e01e79f91267dc17b4e6314d5d4d03593d2ceee0fbb452b750bd70ea5af9.
//
// Solidity: event ChainlinkRequested(bytes32 indexed id)
func (_APIConsumer *APIConsumerFilterer) ParseChainlinkRequested(log types.Log) (*APIConsumerChainlinkRequested, error) {
	event := new(APIConsumerChainlinkRequested)
	if err := _APIConsumer.contract.UnpackLog(event, "ChainlinkRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// APIConsumerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the APIConsumer contract.
type APIConsumerOwnershipTransferredIterator struct {
	Event *APIConsumerOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *APIConsumerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(APIConsumerOwnershipTransferred)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(APIConsumerOwnershipTransferred)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *APIConsumerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *APIConsumerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// APIConsumerOwnershipTransferred represents a OwnershipTransferred event raised by the APIConsumer contract.
type APIConsumerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_APIConsumer *APIConsumerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*APIConsumerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _APIConsumer.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &APIConsumerOwnershipTransferredIterator{contract: _APIConsumer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_APIConsumer *APIConsumerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *APIConsumerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _APIConsumer.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(APIConsumerOwnershipTransferred)
				if err := _APIConsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_APIConsumer *APIConsumerFilterer) ParseOwnershipTransferred(log types.Log) (*APIConsumerOwnershipTransferred, error) {
	event := new(APIConsumerOwnershipTransferred)
	if err := _APIConsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// APIConsumerPerfMetricsEventIterator is returned from FilterPerfMetricsEvent and is used to iterate over the raw logs and unpacked data for PerfMetricsEvent events raised by the APIConsumer contract.
type APIConsumerPerfMetricsEventIterator struct {
	Event *APIConsumerPerfMetricsEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *APIConsumerPerfMetricsEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(APIConsumerPerfMetricsEvent)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(APIConsumerPerfMetricsEvent)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *APIConsumerPerfMetricsEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *APIConsumerPerfMetricsEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// APIConsumerPerfMetricsEvent represents a PerfMetricsEvent event raised by the APIConsumer contract.
type APIConsumerPerfMetricsEvent struct {
	RoundID   *big.Int
	RequestId [32]byte
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterPerfMetricsEvent is a free log retrieval operation binding the contract event 0xfbaf68ee7b9032982942607eaea1859969ed8674797b5c2fc6fecaa753851946.
//
// Solidity: event PerfMetricsEvent(uint256 roundID, bytes32 requestId, uint256 timestamp)
func (_APIConsumer *APIConsumerFilterer) FilterPerfMetricsEvent(opts *bind.FilterOpts) (*APIConsumerPerfMetricsEventIterator, error) {

	logs, sub, err := _APIConsumer.contract.FilterLogs(opts, "PerfMetricsEvent")
	if err != nil {
		return nil, err
	}
	return &APIConsumerPerfMetricsEventIterator{contract: _APIConsumer.contract, event: "PerfMetricsEvent", logs: logs, sub: sub}, nil
}

// WatchPerfMetricsEvent is a free log subscription operation binding the contract event 0xfbaf68ee7b9032982942607eaea1859969ed8674797b5c2fc6fecaa753851946.
//
// Solidity: event PerfMetricsEvent(uint256 roundID, bytes32 requestId, uint256 timestamp)
func (_APIConsumer *APIConsumerFilterer) WatchPerfMetricsEvent(opts *bind.WatchOpts, sink chan<- *APIConsumerPerfMetricsEvent) (event.Subscription, error) {

	logs, sub, err := _APIConsumer.contract.WatchLogs(opts, "PerfMetricsEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(APIConsumerPerfMetricsEvent)
				if err := _APIConsumer.contract.UnpackLog(event, "PerfMetricsEvent", log); err != nil {
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

// ParsePerfMetricsEvent is a log parse operation binding the contract event 0xfbaf68ee7b9032982942607eaea1859969ed8674797b5c2fc6fecaa753851946.
//
// Solidity: event PerfMetricsEvent(uint256 roundID, bytes32 requestId, uint256 timestamp)
func (_APIConsumer *APIConsumerFilterer) ParsePerfMetricsEvent(log types.Log) (*APIConsumerPerfMetricsEvent, error) {
	event := new(APIConsumerPerfMetricsEvent)
	if err := _APIConsumer.contract.UnpackLog(event, "PerfMetricsEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
