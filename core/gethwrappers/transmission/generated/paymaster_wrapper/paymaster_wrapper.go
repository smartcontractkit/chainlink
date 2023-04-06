// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package paymaster_wrapper

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

type UserOperation struct {
	Sender               common.Address
	Nonce                *big.Int
	InitCode             []byte
	CallData             []byte
	CallGasLimit         *big.Int
	VerificationGasLimit *big.Int
	PreVerificationGas   *big.Int
	MaxFeePerGas         *big.Int
	MaxPriorityFeePerGas *big.Int
	PaymasterAndData     []byte
	Signature            []byte
}

var PaymasterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"linkEthFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"entryPoint\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"juelsNeeded\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"subscriptionBalance\",\"type\":\"uint256\"}],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"}],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"userOpHash\",\"type\":\"bytes32\"}],\"name\":\"UserOperationAlreadyTried\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_entryPoint\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_linkEthFeed\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_linkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumIPaymaster.PostOpMode\",\"name\":\"\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"context\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"actualGasCost\",\"type\":\"uint256\"}],\"name\":\"postOp\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_config\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"verificationGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"preVerificationGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxPriorityFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"paymasterAndData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structUserOperation\",\"name\":\"userOp\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"userOpHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"maxCost\",\"type\":\"uint256\"}],\"name\":\"validatePaymasterUserOp\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"context\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"validationData\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b50604051620014fb380380620014fb8339810160408190526200003491620001a3565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000df565b5050506001600160a01b0392831660805290821660a0521660c052620001f7565b336001600160a01b03821603620001395760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b0381168114620001a057600080fd5b50565b600080600060608486031215620001b957600080fd5b8351620001c6816200018a565b6020850151909350620001d9816200018a565b6040850151909250620001ec816200018a565b809150509250925092565b60805160a05160c05161129c6200025f600039600081816101080152818161049f01528181610507015281816105cd015261063501526000818161018f0152610cb60152600081816101dc015281816103a201528181610ac90152610b81015261129c6000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c80639b9bd4de11610081578063db37983b1161005b578063db37983b146101d7578063f2fde38b146101fe578063f465c77e1461021157600080fd5b80639b9bd4de1461018a578063a4c0ed36146101b1578063a9a23409146101c457600080fd5b806379ba5097116100b257806379ba50971461014f5780638a38f365146101595780638da5cb5b1461016c57600080fd5b8063088070f5146100ce578063140fcfb114610103575b600080fd5b6002546003546100e29163ffffffff169082565b6040805163ffffffff90931683526020830191909152015b60405180910390f35b61012a7f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100fa565b610157610232565b005b610157610167366004610d5b565b610334565b60005473ffffffffffffffffffffffffffffffffffffffff1661012a565b61012a7f000000000000000000000000000000000000000000000000000000000000000081565b6101576101bf366004610dfb565b61038a565b6101576101d2366004610e57565b610487565b61012a7f000000000000000000000000000000000000000000000000000000000000000081565b61015761020c366004610eb7565b61059d565b61022461021f366004610edb565b6105b1565b6040516100fa929190610f2f565b60015473ffffffffffffffffffffffffffffffffffffffff1633146102b8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b61033c610849565b6040805180820190915263ffffffff9092168083526020909201819052600280547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000016909217909155600355565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146103f9576040517f44b0e3c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208114610433576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061044182840184610eb7565b73ffffffffffffffffffffffffffffffffffffffff811660009081526005602052604081208054929350869290919061047b908490610fd9565b90915550505050505050565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610534576040517f295a81c100000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660248201526044016102af565b60008061054384860186610ff1565b9150915080610551846108cc565b61055b9190610fd9565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600560205260408120805490919061059090849061100f565b9091555050505050505050565b6105a5610849565b6105ae816108f8565b50565b606060003373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610662576040517f295a81c100000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660248201526044016102af565b60008481526004602052604090205460ff16156106ae576040517f7413dcf8000000000000000000000000000000000000000000000000000000008152600481018590526024016102af565b60006106b9866109ed565b90506000816106c7866108cc565b6106d19190610fd9565b905080600560006106e560208b018b610eb7565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205410156107b057806005600061073860208b018b610eb7565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546040517f03eb8b540000000000000000000000000000000000000000000000000000000081526004016102af929190918252602082015260400190565b600086815260046020908152604090912080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790556107f690880188610eb7565b6040805173ffffffffffffffffffffffffffffffffffffffff9092166020830152810183905260600160405160208183030381529060405261083b6000806000610c29565b935093505050935093915050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146108ca576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016102af565b565b60006108d6610c61565b6108e883670de0b6b3a7640000611026565b6108f29190611063565b92915050565b3373ffffffffffffffffffffffffffffffffffffffff821603610977576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016102af565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006109fd61012083018361109e565b9050601403610a0e57506000919050565b6000610a1e61012084018461109e565b6014818110610a2f57610a2f611103565b919091013560f81c9150819050610c23576000610a5061012085018561109e565b610a5e916015908290611132565b810190610a6b919061115c565b90508060200151600014158015610b385750602081015181516040517f70a0823100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201527f0000000000000000000000000000000000000000000000000000000000000000909116906370a0823190602401602060405180830381865afa158015610b12573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b3691906111e5565b105b15610c2157805160408083015190517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169263a9059cbb92610bd59260040173ffffffffffffffffffffffffffffffffffffffff929092168252602082015260400190565b6020604051808303816000875af1158015610bf4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c1891906111fe565b50806040015192505b505b50919050565b600060d08265ffffffffffff16901b60a08465ffffffffffff16901b85610c51576000610c54565b60015b60ff161717949350505050565b600254604080517ffeaf968c000000000000000000000000000000000000000000000000000000008152905160009263ffffffff1691821515918491829173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169163feaf968c9160048083019260a09291908290030181865afa158015610d01573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d25919061123f565b509450909250849150508015610d495750610d40824261100f565b8463ffffffff16105b15610d5357506003545b949350505050565b60008060408385031215610d6e57600080fd5b823563ffffffff81168114610d8257600080fd5b946020939093013593505050565b73ffffffffffffffffffffffffffffffffffffffff811681146105ae57600080fd5b60008083601f840112610dc457600080fd5b50813567ffffffffffffffff811115610ddc57600080fd5b602083019150836020828501011115610df457600080fd5b9250929050565b60008060008060608587031215610e1157600080fd5b8435610e1c81610d90565b935060208501359250604085013567ffffffffffffffff811115610e3f57600080fd5b610e4b87828801610db2565b95989497509550505050565b60008060008060608587031215610e6d57600080fd5b843560038110610e7c57600080fd5b9350602085013567ffffffffffffffff811115610e9857600080fd5b610ea487828801610db2565b9598909750949560400135949350505050565b600060208284031215610ec957600080fd5b8135610ed481610d90565b9392505050565b600080600060608486031215610ef057600080fd5b833567ffffffffffffffff811115610f0757600080fd5b84016101608187031215610f1a57600080fd5b95602085013595506040909401359392505050565b604081526000835180604084015260005b81811015610f5d5760208187018101516060868401015201610f40565b81811115610f6f576000606083860101525b50602083019390935250601f919091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01601606001919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60008219821115610fec57610fec610faa565b500190565b6000806040838503121561100457600080fd5b8235610d8281610d90565b60008282101561102157611021610faa565b500390565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561105e5761105e610faa565b500290565b600082611099577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126110d357600080fd5b83018035915067ffffffffffffffff8211156110ee57600080fd5b602001915036819003821315610df457600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6000808585111561114257600080fd5b8386111561114f57600080fd5b5050820193919092039150565b60006060828403121561116e57600080fd5b6040516060810181811067ffffffffffffffff821117156111b8577f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405282356111c681610d90565b8152602083810135908201526040928301359281019290925250919050565b6000602082840312156111f757600080fd5b5051919050565b60006020828403121561121057600080fd5b81518015158114610ed457600080fd5b805169ffffffffffffffffffff8116811461123a57600080fd5b919050565b600080600080600060a0868803121561125757600080fd5b61126086611220565b945060208601519350604086015192506060860151915061128360808701611220565b9050929550929590935056fea164736f6c634300080f000a",
}

var PaymasterABI = PaymasterMetaData.ABI

var PaymasterBin = PaymasterMetaData.Bin

func DeployPaymaster(auth *bind.TransactOpts, backend bind.ContractBackend, linkToken common.Address, linkEthFeed common.Address, entryPoint common.Address) (common.Address, *types.Transaction, *Paymaster, error) {
	parsed, err := PaymasterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PaymasterBin), backend, linkToken, linkEthFeed, entryPoint)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Paymaster{PaymasterCaller: PaymasterCaller{contract: contract}, PaymasterTransactor: PaymasterTransactor{contract: contract}, PaymasterFilterer: PaymasterFilterer{contract: contract}}, nil
}

type Paymaster struct {
	address common.Address
	abi     abi.ABI
	PaymasterCaller
	PaymasterTransactor
	PaymasterFilterer
}

type PaymasterCaller struct {
	contract *bind.BoundContract
}

type PaymasterTransactor struct {
	contract *bind.BoundContract
}

type PaymasterFilterer struct {
	contract *bind.BoundContract
}

type PaymasterSession struct {
	Contract     *Paymaster
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type PaymasterCallerSession struct {
	Contract *PaymasterCaller
	CallOpts bind.CallOpts
}

type PaymasterTransactorSession struct {
	Contract     *PaymasterTransactor
	TransactOpts bind.TransactOpts
}

type PaymasterRaw struct {
	Contract *Paymaster
}

type PaymasterCallerRaw struct {
	Contract *PaymasterCaller
}

type PaymasterTransactorRaw struct {
	Contract *PaymasterTransactor
}

func NewPaymaster(address common.Address, backend bind.ContractBackend) (*Paymaster, error) {
	abi, err := abi.JSON(strings.NewReader(PaymasterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindPaymaster(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Paymaster{address: address, abi: abi, PaymasterCaller: PaymasterCaller{contract: contract}, PaymasterTransactor: PaymasterTransactor{contract: contract}, PaymasterFilterer: PaymasterFilterer{contract: contract}}, nil
}

func NewPaymasterCaller(address common.Address, caller bind.ContractCaller) (*PaymasterCaller, error) {
	contract, err := bindPaymaster(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PaymasterCaller{contract: contract}, nil
}

func NewPaymasterTransactor(address common.Address, transactor bind.ContractTransactor) (*PaymasterTransactor, error) {
	contract, err := bindPaymaster(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PaymasterTransactor{contract: contract}, nil
}

func NewPaymasterFilterer(address common.Address, filterer bind.ContractFilterer) (*PaymasterFilterer, error) {
	contract, err := bindPaymaster(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PaymasterFilterer{contract: contract}, nil
}

func bindPaymaster(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PaymasterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_Paymaster *PaymasterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Paymaster.Contract.PaymasterCaller.contract.Call(opts, result, method, params...)
}

func (_Paymaster *PaymasterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Paymaster.Contract.PaymasterTransactor.contract.Transfer(opts)
}

func (_Paymaster *PaymasterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Paymaster.Contract.PaymasterTransactor.contract.Transact(opts, method, params...)
}

func (_Paymaster *PaymasterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Paymaster.Contract.contract.Call(opts, result, method, params...)
}

func (_Paymaster *PaymasterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Paymaster.Contract.contract.Transfer(opts)
}

func (_Paymaster *PaymasterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Paymaster.Contract.contract.Transact(opts, method, params...)
}

func (_Paymaster *PaymasterCaller) IEntryPoint(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Paymaster.contract.Call(opts, &out, "i_entryPoint")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_Paymaster *PaymasterSession) IEntryPoint() (common.Address, error) {
	return _Paymaster.Contract.IEntryPoint(&_Paymaster.CallOpts)
}

func (_Paymaster *PaymasterCallerSession) IEntryPoint() (common.Address, error) {
	return _Paymaster.Contract.IEntryPoint(&_Paymaster.CallOpts)
}

func (_Paymaster *PaymasterCaller) ILinkEthFeed(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Paymaster.contract.Call(opts, &out, "i_linkEthFeed")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_Paymaster *PaymasterSession) ILinkEthFeed() (common.Address, error) {
	return _Paymaster.Contract.ILinkEthFeed(&_Paymaster.CallOpts)
}

func (_Paymaster *PaymasterCallerSession) ILinkEthFeed() (common.Address, error) {
	return _Paymaster.Contract.ILinkEthFeed(&_Paymaster.CallOpts)
}

func (_Paymaster *PaymasterCaller) ILinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Paymaster.contract.Call(opts, &out, "i_linkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_Paymaster *PaymasterSession) ILinkToken() (common.Address, error) {
	return _Paymaster.Contract.ILinkToken(&_Paymaster.CallOpts)
}

func (_Paymaster *PaymasterCallerSession) ILinkToken() (common.Address, error) {
	return _Paymaster.Contract.ILinkToken(&_Paymaster.CallOpts)
}

func (_Paymaster *PaymasterCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Paymaster.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_Paymaster *PaymasterSession) Owner() (common.Address, error) {
	return _Paymaster.Contract.Owner(&_Paymaster.CallOpts)
}

func (_Paymaster *PaymasterCallerSession) Owner() (common.Address, error) {
	return _Paymaster.Contract.Owner(&_Paymaster.CallOpts)
}

func (_Paymaster *PaymasterCaller) SConfig(opts *bind.CallOpts) (SConfig,

	error) {
	var out []interface{}
	err := _Paymaster.contract.Call(opts, &out, "s_config")

	outstruct := new(SConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.StalenessSeconds = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.FallbackWeiPerUnitLink = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_Paymaster *PaymasterSession) SConfig() (SConfig,

	error) {
	return _Paymaster.Contract.SConfig(&_Paymaster.CallOpts)
}

func (_Paymaster *PaymasterCallerSession) SConfig() (SConfig,

	error) {
	return _Paymaster.Contract.SConfig(&_Paymaster.CallOpts)
}

func (_Paymaster *PaymasterTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Paymaster.contract.Transact(opts, "acceptOwnership")
}

func (_Paymaster *PaymasterSession) AcceptOwnership() (*types.Transaction, error) {
	return _Paymaster.Contract.AcceptOwnership(&_Paymaster.TransactOpts)
}

func (_Paymaster *PaymasterTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Paymaster.Contract.AcceptOwnership(&_Paymaster.TransactOpts)
}

func (_Paymaster *PaymasterTransactor) OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _Paymaster.contract.Transact(opts, "onTokenTransfer", arg0, _amount, _data)
}

func (_Paymaster *PaymasterSession) OnTokenTransfer(arg0 common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _Paymaster.Contract.OnTokenTransfer(&_Paymaster.TransactOpts, arg0, _amount, _data)
}

func (_Paymaster *PaymasterTransactorSession) OnTokenTransfer(arg0 common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _Paymaster.Contract.OnTokenTransfer(&_Paymaster.TransactOpts, arg0, _amount, _data)
}

func (_Paymaster *PaymasterTransactor) PostOp(opts *bind.TransactOpts, arg0 uint8, context []byte, actualGasCost *big.Int) (*types.Transaction, error) {
	return _Paymaster.contract.Transact(opts, "postOp", arg0, context, actualGasCost)
}

func (_Paymaster *PaymasterSession) PostOp(arg0 uint8, context []byte, actualGasCost *big.Int) (*types.Transaction, error) {
	return _Paymaster.Contract.PostOp(&_Paymaster.TransactOpts, arg0, context, actualGasCost)
}

func (_Paymaster *PaymasterTransactorSession) PostOp(arg0 uint8, context []byte, actualGasCost *big.Int) (*types.Transaction, error) {
	return _Paymaster.Contract.PostOp(&_Paymaster.TransactOpts, arg0, context, actualGasCost)
}

func (_Paymaster *PaymasterTransactor) SetConfig(opts *bind.TransactOpts, stalenessSeconds uint32, fallbackWeiPerUnitLink *big.Int) (*types.Transaction, error) {
	return _Paymaster.contract.Transact(opts, "setConfig", stalenessSeconds, fallbackWeiPerUnitLink)
}

func (_Paymaster *PaymasterSession) SetConfig(stalenessSeconds uint32, fallbackWeiPerUnitLink *big.Int) (*types.Transaction, error) {
	return _Paymaster.Contract.SetConfig(&_Paymaster.TransactOpts, stalenessSeconds, fallbackWeiPerUnitLink)
}

func (_Paymaster *PaymasterTransactorSession) SetConfig(stalenessSeconds uint32, fallbackWeiPerUnitLink *big.Int) (*types.Transaction, error) {
	return _Paymaster.Contract.SetConfig(&_Paymaster.TransactOpts, stalenessSeconds, fallbackWeiPerUnitLink)
}

func (_Paymaster *PaymasterTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _Paymaster.contract.Transact(opts, "transferOwnership", to)
}

func (_Paymaster *PaymasterSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _Paymaster.Contract.TransferOwnership(&_Paymaster.TransactOpts, to)
}

func (_Paymaster *PaymasterTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _Paymaster.Contract.TransferOwnership(&_Paymaster.TransactOpts, to)
}

func (_Paymaster *PaymasterTransactor) ValidatePaymasterUserOp(opts *bind.TransactOpts, userOp UserOperation, userOpHash [32]byte, maxCost *big.Int) (*types.Transaction, error) {
	return _Paymaster.contract.Transact(opts, "validatePaymasterUserOp", userOp, userOpHash, maxCost)
}

func (_Paymaster *PaymasterSession) ValidatePaymasterUserOp(userOp UserOperation, userOpHash [32]byte, maxCost *big.Int) (*types.Transaction, error) {
	return _Paymaster.Contract.ValidatePaymasterUserOp(&_Paymaster.TransactOpts, userOp, userOpHash, maxCost)
}

func (_Paymaster *PaymasterTransactorSession) ValidatePaymasterUserOp(userOp UserOperation, userOpHash [32]byte, maxCost *big.Int) (*types.Transaction, error) {
	return _Paymaster.Contract.ValidatePaymasterUserOp(&_Paymaster.TransactOpts, userOp, userOpHash, maxCost)
}

type PaymasterOwnershipTransferRequestedIterator struct {
	Event *PaymasterOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PaymasterOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymasterOwnershipTransferRequested)
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
		it.Event = new(PaymasterOwnershipTransferRequested)
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

func (it *PaymasterOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *PaymasterOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PaymasterOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_Paymaster *PaymasterFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PaymasterOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Paymaster.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &PaymasterOwnershipTransferRequestedIterator{contract: _Paymaster.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_Paymaster *PaymasterFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *PaymasterOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Paymaster.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PaymasterOwnershipTransferRequested)
				if err := _Paymaster.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_Paymaster *PaymasterFilterer) ParseOwnershipTransferRequested(log types.Log) (*PaymasterOwnershipTransferRequested, error) {
	event := new(PaymasterOwnershipTransferRequested)
	if err := _Paymaster.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PaymasterOwnershipTransferredIterator struct {
	Event *PaymasterOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PaymasterOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymasterOwnershipTransferred)
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
		it.Event = new(PaymasterOwnershipTransferred)
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

func (it *PaymasterOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *PaymasterOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PaymasterOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_Paymaster *PaymasterFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PaymasterOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Paymaster.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &PaymasterOwnershipTransferredIterator{contract: _Paymaster.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_Paymaster *PaymasterFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PaymasterOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Paymaster.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PaymasterOwnershipTransferred)
				if err := _Paymaster.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_Paymaster *PaymasterFilterer) ParseOwnershipTransferred(log types.Log) (*PaymasterOwnershipTransferred, error) {
	event := new(PaymasterOwnershipTransferred)
	if err := _Paymaster.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SConfig struct {
	StalenessSeconds       uint32
	FallbackWeiPerUnitLink *big.Int
}

func (_Paymaster *Paymaster) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _Paymaster.abi.Events["OwnershipTransferRequested"].ID:
		return _Paymaster.ParseOwnershipTransferRequested(log)
	case _Paymaster.abi.Events["OwnershipTransferred"].ID:
		return _Paymaster.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (PaymasterOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (PaymasterOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_Paymaster *Paymaster) Address() common.Address {
	return _Paymaster.address
}

type PaymasterInterface interface {
	IEntryPoint(opts *bind.CallOpts) (common.Address, error)

	ILinkEthFeed(opts *bind.CallOpts) (common.Address, error)

	ILinkToken(opts *bind.CallOpts) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SConfig(opts *bind.CallOpts) (SConfig,

		error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error)

	PostOp(opts *bind.TransactOpts, arg0 uint8, context []byte, actualGasCost *big.Int) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, stalenessSeconds uint32, fallbackWeiPerUnitLink *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	ValidatePaymasterUserOp(opts *bind.TransactOpts, userOp UserOperation, userOpHash [32]byte, maxCost *big.Int) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PaymasterOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *PaymasterOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*PaymasterOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PaymasterOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PaymasterOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*PaymasterOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
