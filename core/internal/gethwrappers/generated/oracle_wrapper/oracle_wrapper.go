// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package oracle_wrapper

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
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
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

var OracleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"CancelOracleRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"specId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"callbackAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"callbackFunctionId\",\"type\":\"bytes4\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"cancelExpiration\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"dataVersion\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"OracleRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"EXPIRY_TIME\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_payment\",\"type\":\"uint256\"},{\"internalType\":\"bytes4\",\"name\":\"_callbackFunc\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"_expiration\",\"type\":\"uint256\"}],\"name\":\"cancelOracleRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_payment\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_callbackAddress\",\"type\":\"address\"},{\"internalType\":\"bytes4\",\"name\":\"_callbackFunctionId\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"_expiration\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_data\",\"type\":\"bytes32\"}],\"name\":\"fulfillOracleRequest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_node\",\"type\":\"address\"}],\"name\":\"getAuthorizationStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChainlinkToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_payment\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_specId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_callbackAddress\",\"type\":\"address\"},{\"internalType\":\"bytes4\",\"name\":\"_callbackFunctionId\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"_nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_dataVersion\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"oracleRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_node\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_allowed\",\"type\":\"bool\"}],\"name\":\"setFulfillmentPermission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052600160045534801561001557600080fd5b506040516119513803806119518339818101604052602081101561003857600080fd5b5051600080546001600160a01b03191633178082556040516001600160a01b039190911691907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a3600180546001600160a01b0319166001600160a01b039290921691909117905561189f806100b26000396000f3fe608060405234801561001057600080fd5b50600436106100df5760003560e01c80637fcd56db1161008c578063a4c0ed3611610066578063a4c0ed3614610332578063d3e9c314146103fa578063f2fde38b1461042d578063f3fef3a314610460576100df565b80637fcd56db146102e75780638da5cb5b146103225780638f32d59b1461032a576100df565b80634b602282116100bd5780634b60228214610274578063501883011461028e5780636ee4d55314610296576100df565b8063165d35e1146100e457806340429946146101155780634ab0d190146101ed575b600080fd5b6100ec610499565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b6101eb600480360361010081101561012c57600080fd5b73ffffffffffffffffffffffffffffffffffffffff8235811692602081013592604082013592606083013516917fffffffff000000000000000000000000000000000000000000000000000000006080820135169160a08201359160c081013591810190610100810160e08201356401000000008111156101ac57600080fd5b8201836020820111156101be57600080fd5b803590602001918460018302840111640100000000831117156101e057600080fd5b5090925090506104b5565b005b610260600480360360c081101561020357600080fd5b5080359060208101359073ffffffffffffffffffffffffffffffffffffffff604082013516907fffffffff000000000000000000000000000000000000000000000000000000006060820135169060808101359060a001356108e6565b604080519115158252519081900360200190f35b61027c610ce5565b60408051918252519081900360200190f35b61027c610ceb565b6101eb600480360360808110156102ac57600080fd5b508035906020810135907fffffffff000000000000000000000000000000000000000000000000000000006040820135169060600135610d79565b6101eb600480360360408110156102fd57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81351690602001351515610fac565b6100ec611075565b610260611091565b6101eb6004803603606081101561034857600080fd5b73ffffffffffffffffffffffffffffffffffffffff8235169160208101359181019060608101604082013564010000000081111561038557600080fd5b82018360208201111561039757600080fd5b803590602001918460018302840111640100000000831117156103b957600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295506110af945050505050565b6102606004803603602081101561041057600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166113cb565b6101eb6004803603602081101561044357600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166113f6565b6101eb6004803603604081101561047657600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135611475565b60015473ffffffffffffffffffffffffffffffffffffffff1690565b6104bd610499565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461055657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f4d75737420757365204c494e4b20746f6b656e00000000000000000000000000604482015290519081900360640190fd5b600154869073ffffffffffffffffffffffffffffffffffffffff808316911614156105e257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f742063616c6c6261636b20746f204c494e4b000000000000000000604482015290519081900360640190fd5b604080517fffffffffffffffffffffffffffffffffffffffff00000000000000000000000060608d901b16602080830191909152603480830189905283518084039091018152605490920183528151918101919091206000818152600290925291902054156106b257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f4d75737420757365206120756e69717565204944000000000000000000000000604482015290519081900360640190fd5b60006106c64261012c63ffffffff61162216565b90508a898983604051602001808581526020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1660601b8152601401837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152600401828152602001945050505050604051602081830303815290604052805190602001206002600084815260200190815260200160002081905550897fd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c658d848e8d8d878d8d8d604051808a73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018981526020018881526020018773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001867bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001858152602001848152602001806020018281038252848482818152602001925080828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169092018290039c50909a5050505050505050505050a2505050505050505050505050565b3360009081526003602052604081205460ff16806109365750610907611075565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b61098b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a815260200180611869602a913960400191505060405180910390fd5b6000878152600260205260409020548790610a0757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f4d757374206861766520612076616c6964207265717565737449640000000000604482015290519081900360640190fd5b6040805160208082018a90527fffffffffffffffffffffffffffffffffffffffff00000000000000000000000060608a901b16828401527fffffffff00000000000000000000000000000000000000000000000000000000881660548301526058808301889052835180840390910181526078909201835281519181019190912060008b81526002909252919020548114610b0357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f506172616d7320646f206e6f74206d6174636820726571756573742049440000604482015290519081900360640190fd5b600454610b16908963ffffffff61162216565b60045560008981526002602052604081205562061a805a1015610b9a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4d7573742070726f7669646520636f6e73756d657220656e6f75676820676173604482015290519081900360640190fd5b60408051602481018b9052604480820187905282518083039091018152606490910182526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000008a161781529151815160009373ffffffffffffffffffffffffffffffffffffffff8c169392918291908083835b60208310610c6d57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610c30565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114610ccf576040519150601f19603f3d011682016040523d82523d6000602084013e610cd4565b606091505b50909b9a5050505050505050505050565b61012c81565b6000610cf5611091565b610d6057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b600454610d7490600163ffffffff61169d16565b905090565b6040805160208082018690523360601b828401527fffffffff00000000000000000000000000000000000000000000000000000000851660548301526058808301859052835180840390910181526078909201835281519181019190912060008781526002909252919020548114610e5257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f506172616d7320646f206e6f74206d6174636820726571756573742049440000604482015290519081900360640190fd5b42821115610ec157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f52657175657374206973206e6f74206578706972656400000000000000000000604482015290519081900360640190fd5b6000858152600260205260408082208290555186917fa7842b9ec549398102c0d91b1b9919b2f20558aefdadf57528a95c6cd3292e9391a2600154604080517fa9059cbb00000000000000000000000000000000000000000000000000000000815233600482015260248101879052905173ffffffffffffffffffffffffffffffffffffffff9092169163a9059cbb916044808201926020929091908290030181600087803b158015610f7357600080fd5b505af1158015610f87573d6000803e3d6000fd5b505050506040513d6020811015610f9d57600080fd5b5051610fa557fe5b5050505050565b610fb4611091565b61101f57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff91909116600090815260036020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055565b60005473ffffffffffffffffffffffffffffffffffffffff1690565b60005473ffffffffffffffffffffffffffffffffffffffff16331490565b6110b7610499565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461115057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f4d75737420757365204c494e4b20746f6b656e00000000000000000000000000604482015290519081900360640190fd5b80518190604411156111c357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f496e76616c69642072657175657374206c656e67746800000000000000000000604482015290519081900360640190fd5b602082015182907fffffffff0000000000000000000000000000000000000000000000000000000081167f40429946000000000000000000000000000000000000000000000000000000001461127a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f4d757374207573652077686974656c69737465642066756e6374696f6e730000604482015290519081900360640190fd5b85602485015284604485015260003073ffffffffffffffffffffffffffffffffffffffff16856040518082805190602001908083835b602083106112ed57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe090920191602091820191016112b0565b6001836020036101000a038019825116818451168082178552505050505050905001915050600060405180830381855af49150503d806000811461134d576040519150601f19603f3d011682016040523d82523d6000602084013e611352565b606091505b50509050806113c257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f556e61626c6520746f2063726561746520726571756573740000000000000000604482015290519081900360640190fd5b50505050505050565b73ffffffffffffffffffffffffffffffffffffffff1660009081526003602052604090205460ff1690565b6113fe611091565b61146957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b61147281611714565b50565b61147d611091565b6114e857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b806114fa81600163ffffffff61162216565b6004541015611554576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260358152602001806118346035913960400191505060405180910390fd5b600454611567908363ffffffff61169d16565b6004908155600154604080517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff87811694820194909452602481018690529051929091169163a9059cbb916044808201926020929091908290030181600087803b1580156115eb57600080fd5b505af11580156115ff573d6000803e3d6000fd5b505050506040513d602081101561161557600080fd5b505161161d57fe5b505050565b60008282018381101561169657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b9392505050565b60008282111561170e57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b73ffffffffffffffffffffffffffffffffffffffff8116611780576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602681526020018061180e6026913960400191505060405180910390fd5b6000805460405173ffffffffffffffffffffffffffffffffffffffff808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9290921691909117905556fe4f776e61626c653a206e6577206f776e657220697320746865207a65726f2061646472657373416d6f756e74207265717565737465642069732067726561746572207468616e20776974686472617761626c652062616c616e63654e6f7420616e20617574686f72697a6564206e6f646520746f2066756c66696c6c207265717565737473a164736f6c6343000606000a",
}

var OracleABI = OracleMetaData.ABI

var OracleBin = OracleMetaData.Bin

func DeployOracle(auth *bind.TransactOpts, backend bind.ContractBackend, _link common.Address) (common.Address, *types.Transaction, *Oracle, error) {
	parsed, err := OracleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OracleBin), backend, _link)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Oracle{OracleCaller: OracleCaller{contract: contract}, OracleTransactor: OracleTransactor{contract: contract}, OracleFilterer: OracleFilterer{contract: contract}}, nil
}

type Oracle struct {
	address common.Address
	abi     abi.ABI
	OracleCaller
	OracleTransactor
	OracleFilterer
}

type OracleCaller struct {
	contract *bind.BoundContract
}

type OracleTransactor struct {
	contract *bind.BoundContract
}

type OracleFilterer struct {
	contract *bind.BoundContract
}

type OracleSession struct {
	Contract     *Oracle
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OracleCallerSession struct {
	Contract *OracleCaller
	CallOpts bind.CallOpts
}

type OracleTransactorSession struct {
	Contract     *OracleTransactor
	TransactOpts bind.TransactOpts
}

type OracleRaw struct {
	Contract *Oracle
}

type OracleCallerRaw struct {
	Contract *OracleCaller
}

type OracleTransactorRaw struct {
	Contract *OracleTransactor
}

func NewOracle(address common.Address, backend bind.ContractBackend) (*Oracle, error) {
	abi, err := abi.JSON(strings.NewReader(OracleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Oracle{address: address, abi: abi, OracleCaller: OracleCaller{contract: contract}, OracleTransactor: OracleTransactor{contract: contract}, OracleFilterer: OracleFilterer{contract: contract}}, nil
}

func NewOracleCaller(address common.Address, caller bind.ContractCaller) (*OracleCaller, error) {
	contract, err := bindOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OracleCaller{contract: contract}, nil
}

func NewOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*OracleTransactor, error) {
	contract, err := bindOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OracleTransactor{contract: contract}, nil
}

func NewOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*OracleFilterer, error) {
	contract, err := bindOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OracleFilterer{contract: contract}, nil
}

func bindOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OracleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_Oracle *OracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Oracle.Contract.OracleCaller.contract.Call(opts, result, method, params...)
}

func (_Oracle *OracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.Contract.OracleTransactor.contract.Transfer(opts)
}

func (_Oracle *OracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Oracle.Contract.OracleTransactor.contract.Transact(opts, method, params...)
}

func (_Oracle *OracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Oracle.Contract.contract.Call(opts, result, method, params...)
}

func (_Oracle *OracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.Contract.contract.Transfer(opts)
}

func (_Oracle *OracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Oracle.Contract.contract.Transact(opts, method, params...)
}

func (_Oracle *OracleCaller) EXPIRYTIME(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "EXPIRY_TIME")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_Oracle *OracleSession) EXPIRYTIME() (*big.Int, error) {
	return _Oracle.Contract.EXPIRYTIME(&_Oracle.CallOpts)
}

func (_Oracle *OracleCallerSession) EXPIRYTIME() (*big.Int, error) {
	return _Oracle.Contract.EXPIRYTIME(&_Oracle.CallOpts)
}

func (_Oracle *OracleCaller) GetAuthorizationStatus(opts *bind.CallOpts, _node common.Address) (bool, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "getAuthorizationStatus", _node)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_Oracle *OracleSession) GetAuthorizationStatus(_node common.Address) (bool, error) {
	return _Oracle.Contract.GetAuthorizationStatus(&_Oracle.CallOpts, _node)
}

func (_Oracle *OracleCallerSession) GetAuthorizationStatus(_node common.Address) (bool, error) {
	return _Oracle.Contract.GetAuthorizationStatus(&_Oracle.CallOpts, _node)
}

func (_Oracle *OracleCaller) GetChainlinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "getChainlinkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_Oracle *OracleSession) GetChainlinkToken() (common.Address, error) {
	return _Oracle.Contract.GetChainlinkToken(&_Oracle.CallOpts)
}

func (_Oracle *OracleCallerSession) GetChainlinkToken() (common.Address, error) {
	return _Oracle.Contract.GetChainlinkToken(&_Oracle.CallOpts)
}

func (_Oracle *OracleCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "isOwner")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_Oracle *OracleSession) IsOwner() (bool, error) {
	return _Oracle.Contract.IsOwner(&_Oracle.CallOpts)
}

func (_Oracle *OracleCallerSession) IsOwner() (bool, error) {
	return _Oracle.Contract.IsOwner(&_Oracle.CallOpts)
}

func (_Oracle *OracleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_Oracle *OracleSession) Owner() (common.Address, error) {
	return _Oracle.Contract.Owner(&_Oracle.CallOpts)
}

func (_Oracle *OracleCallerSession) Owner() (common.Address, error) {
	return _Oracle.Contract.Owner(&_Oracle.CallOpts)
}

func (_Oracle *OracleCaller) Withdrawable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "withdrawable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_Oracle *OracleSession) Withdrawable() (*big.Int, error) {
	return _Oracle.Contract.Withdrawable(&_Oracle.CallOpts)
}

func (_Oracle *OracleCallerSession) Withdrawable() (*big.Int, error) {
	return _Oracle.Contract.Withdrawable(&_Oracle.CallOpts)
}

func (_Oracle *OracleTransactor) CancelOracleRequest(opts *bind.TransactOpts, _requestId [32]byte, _payment *big.Int, _callbackFunc [4]byte, _expiration *big.Int) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "cancelOracleRequest", _requestId, _payment, _callbackFunc, _expiration)
}

func (_Oracle *OracleSession) CancelOracleRequest(_requestId [32]byte, _payment *big.Int, _callbackFunc [4]byte, _expiration *big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.CancelOracleRequest(&_Oracle.TransactOpts, _requestId, _payment, _callbackFunc, _expiration)
}

func (_Oracle *OracleTransactorSession) CancelOracleRequest(_requestId [32]byte, _payment *big.Int, _callbackFunc [4]byte, _expiration *big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.CancelOracleRequest(&_Oracle.TransactOpts, _requestId, _payment, _callbackFunc, _expiration)
}

func (_Oracle *OracleTransactor) FulfillOracleRequest(opts *bind.TransactOpts, _requestId [32]byte, _payment *big.Int, _callbackAddress common.Address, _callbackFunctionId [4]byte, _expiration *big.Int, _data [32]byte) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "fulfillOracleRequest", _requestId, _payment, _callbackAddress, _callbackFunctionId, _expiration, _data)
}

func (_Oracle *OracleSession) FulfillOracleRequest(_requestId [32]byte, _payment *big.Int, _callbackAddress common.Address, _callbackFunctionId [4]byte, _expiration *big.Int, _data [32]byte) (*types.Transaction, error) {
	return _Oracle.Contract.FulfillOracleRequest(&_Oracle.TransactOpts, _requestId, _payment, _callbackAddress, _callbackFunctionId, _expiration, _data)
}

func (_Oracle *OracleTransactorSession) FulfillOracleRequest(_requestId [32]byte, _payment *big.Int, _callbackAddress common.Address, _callbackFunctionId [4]byte, _expiration *big.Int, _data [32]byte) (*types.Transaction, error) {
	return _Oracle.Contract.FulfillOracleRequest(&_Oracle.TransactOpts, _requestId, _payment, _callbackAddress, _callbackFunctionId, _expiration, _data)
}

func (_Oracle *OracleTransactor) OnTokenTransfer(opts *bind.TransactOpts, _sender common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "onTokenTransfer", _sender, _amount, _data)
}

func (_Oracle *OracleSession) OnTokenTransfer(_sender common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _Oracle.Contract.OnTokenTransfer(&_Oracle.TransactOpts, _sender, _amount, _data)
}

func (_Oracle *OracleTransactorSession) OnTokenTransfer(_sender common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _Oracle.Contract.OnTokenTransfer(&_Oracle.TransactOpts, _sender, _amount, _data)
}

func (_Oracle *OracleTransactor) OracleRequest(opts *bind.TransactOpts, _sender common.Address, _payment *big.Int, _specId [32]byte, _callbackAddress common.Address, _callbackFunctionId [4]byte, _nonce *big.Int, _dataVersion *big.Int, _data []byte) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "oracleRequest", _sender, _payment, _specId, _callbackAddress, _callbackFunctionId, _nonce, _dataVersion, _data)
}

func (_Oracle *OracleSession) OracleRequest(_sender common.Address, _payment *big.Int, _specId [32]byte, _callbackAddress common.Address, _callbackFunctionId [4]byte, _nonce *big.Int, _dataVersion *big.Int, _data []byte) (*types.Transaction, error) {
	return _Oracle.Contract.OracleRequest(&_Oracle.TransactOpts, _sender, _payment, _specId, _callbackAddress, _callbackFunctionId, _nonce, _dataVersion, _data)
}

func (_Oracle *OracleTransactorSession) OracleRequest(_sender common.Address, _payment *big.Int, _specId [32]byte, _callbackAddress common.Address, _callbackFunctionId [4]byte, _nonce *big.Int, _dataVersion *big.Int, _data []byte) (*types.Transaction, error) {
	return _Oracle.Contract.OracleRequest(&_Oracle.TransactOpts, _sender, _payment, _specId, _callbackAddress, _callbackFunctionId, _nonce, _dataVersion, _data)
}

func (_Oracle *OracleTransactor) SetFulfillmentPermission(opts *bind.TransactOpts, _node common.Address, _allowed bool) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "setFulfillmentPermission", _node, _allowed)
}

func (_Oracle *OracleSession) SetFulfillmentPermission(_node common.Address, _allowed bool) (*types.Transaction, error) {
	return _Oracle.Contract.SetFulfillmentPermission(&_Oracle.TransactOpts, _node, _allowed)
}

func (_Oracle *OracleTransactorSession) SetFulfillmentPermission(_node common.Address, _allowed bool) (*types.Transaction, error) {
	return _Oracle.Contract.SetFulfillmentPermission(&_Oracle.TransactOpts, _node, _allowed)
}

func (_Oracle *OracleTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "transferOwnership", newOwner)
}

func (_Oracle *OracleSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.TransferOwnership(&_Oracle.TransactOpts, newOwner)
}

func (_Oracle *OracleTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.TransferOwnership(&_Oracle.TransactOpts, newOwner)
}

func (_Oracle *OracleTransactor) Withdraw(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "withdraw", _recipient, _amount)
}

func (_Oracle *OracleSession) Withdraw(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.Withdraw(&_Oracle.TransactOpts, _recipient, _amount)
}

func (_Oracle *OracleTransactorSession) Withdraw(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.Withdraw(&_Oracle.TransactOpts, _recipient, _amount)
}

type OracleCancelOracleRequestIterator struct {
	Event *OracleCancelOracleRequest

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OracleCancelOracleRequestIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleCancelOracleRequest)
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
		it.Event = new(OracleCancelOracleRequest)
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

func (it *OracleCancelOracleRequestIterator) Error() error {
	return it.fail
}

func (it *OracleCancelOracleRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OracleCancelOracleRequest struct {
	RequestId [32]byte
	Raw       types.Log
}

func (_Oracle *OracleFilterer) FilterCancelOracleRequest(opts *bind.FilterOpts, requestId [][32]byte) (*OracleCancelOracleRequestIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "CancelOracleRequest", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &OracleCancelOracleRequestIterator{contract: _Oracle.contract, event: "CancelOracleRequest", logs: logs, sub: sub}, nil
}

func (_Oracle *OracleFilterer) WatchCancelOracleRequest(opts *bind.WatchOpts, sink chan<- *OracleCancelOracleRequest, requestId [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "CancelOracleRequest", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OracleCancelOracleRequest)
				if err := _Oracle.contract.UnpackLog(event, "CancelOracleRequest", log); err != nil {
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

func (_Oracle *OracleFilterer) ParseCancelOracleRequest(log types.Log) (*OracleCancelOracleRequest, error) {
	event := new(OracleCancelOracleRequest)
	if err := _Oracle.contract.UnpackLog(event, "CancelOracleRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OracleOracleRequestIterator struct {
	Event *OracleOracleRequest

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OracleOracleRequestIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleOracleRequest)
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
		it.Event = new(OracleOracleRequest)
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

func (it *OracleOracleRequestIterator) Error() error {
	return it.fail
}

func (it *OracleOracleRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OracleOracleRequest struct {
	SpecId             [32]byte
	Requester          common.Address
	RequestId          [32]byte
	Payment            *big.Int
	CallbackAddr       common.Address
	CallbackFunctionId [4]byte
	CancelExpiration   *big.Int
	DataVersion        *big.Int
	Data               []byte
	Raw                types.Log
}

func (_Oracle *OracleFilterer) FilterOracleRequest(opts *bind.FilterOpts, specId [][32]byte) (*OracleOracleRequestIterator, error) {

	var specIdRule []interface{}
	for _, specIdItem := range specId {
		specIdRule = append(specIdRule, specIdItem)
	}

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "OracleRequest", specIdRule)
	if err != nil {
		return nil, err
	}
	return &OracleOracleRequestIterator{contract: _Oracle.contract, event: "OracleRequest", logs: logs, sub: sub}, nil
}

func (_Oracle *OracleFilterer) WatchOracleRequest(opts *bind.WatchOpts, sink chan<- *OracleOracleRequest, specId [][32]byte) (event.Subscription, error) {

	var specIdRule []interface{}
	for _, specIdItem := range specId {
		specIdRule = append(specIdRule, specIdItem)
	}

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "OracleRequest", specIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OracleOracleRequest)
				if err := _Oracle.contract.UnpackLog(event, "OracleRequest", log); err != nil {
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

func (_Oracle *OracleFilterer) ParseOracleRequest(log types.Log) (*OracleOracleRequest, error) {
	event := new(OracleOracleRequest)
	if err := _Oracle.contract.UnpackLog(event, "OracleRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OracleOwnershipTransferredIterator struct {
	Event *OracleOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OracleOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleOwnershipTransferred)
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
		it.Event = new(OracleOwnershipTransferred)
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

func (it *OracleOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *OracleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OracleOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log
}

func (_Oracle *OracleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OracleOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OracleOwnershipTransferredIterator{contract: _Oracle.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_Oracle *OracleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OracleOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OracleOwnershipTransferred)
				if err := _Oracle.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_Oracle *OracleFilterer) ParseOwnershipTransferred(log types.Log) (*OracleOwnershipTransferred, error) {
	event := new(OracleOwnershipTransferred)
	if err := _Oracle.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_Oracle *Oracle) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _Oracle.abi.Events["CancelOracleRequest"].ID:
		return _Oracle.ParseCancelOracleRequest(log)
	case _Oracle.abi.Events["OracleRequest"].ID:
		return _Oracle.ParseOracleRequest(log)
	case _Oracle.abi.Events["OwnershipTransferred"].ID:
		return _Oracle.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (OracleCancelOracleRequest) Topic() common.Hash {
	return common.HexToHash("0xa7842b9ec549398102c0d91b1b9919b2f20558aefdadf57528a95c6cd3292e93")
}

func (OracleOracleRequest) Topic() common.Hash {
	return common.HexToHash("0xd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c65")
}

func (OracleOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_Oracle *Oracle) Address() common.Address {
	return _Oracle.address
}

type OracleInterface interface {
	EXPIRYTIME(opts *bind.CallOpts) (*big.Int, error)

	GetAuthorizationStatus(opts *bind.CallOpts, _node common.Address) (bool, error)

	GetChainlinkToken(opts *bind.CallOpts) (common.Address, error)

	IsOwner(opts *bind.CallOpts) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	Withdrawable(opts *bind.CallOpts) (*big.Int, error)

	CancelOracleRequest(opts *bind.TransactOpts, _requestId [32]byte, _payment *big.Int, _callbackFunc [4]byte, _expiration *big.Int) (*types.Transaction, error)

	FulfillOracleRequest(opts *bind.TransactOpts, _requestId [32]byte, _payment *big.Int, _callbackAddress common.Address, _callbackFunctionId [4]byte, _expiration *big.Int, _data [32]byte) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, _sender common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error)

	OracleRequest(opts *bind.TransactOpts, _sender common.Address, _payment *big.Int, _specId [32]byte, _callbackAddress common.Address, _callbackFunctionId [4]byte, _nonce *big.Int, _dataVersion *big.Int, _data []byte) (*types.Transaction, error)

	SetFulfillmentPermission(opts *bind.TransactOpts, _node common.Address, _allowed bool) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error)

	Withdraw(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error)

	FilterCancelOracleRequest(opts *bind.FilterOpts, requestId [][32]byte) (*OracleCancelOracleRequestIterator, error)

	WatchCancelOracleRequest(opts *bind.WatchOpts, sink chan<- *OracleCancelOracleRequest, requestId [][32]byte) (event.Subscription, error)

	ParseCancelOracleRequest(log types.Log) (*OracleCancelOracleRequest, error)

	FilterOracleRequest(opts *bind.FilterOpts, specId [][32]byte) (*OracleOracleRequestIterator, error)

	WatchOracleRequest(opts *bind.WatchOpts, sink chan<- *OracleOracleRequest, specId [][32]byte) (event.Subscription, error)

	ParseOracleRequest(log types.Log) (*OracleOracleRequest, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OracleOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OracleOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*OracleOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
