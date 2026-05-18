// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package optimism_l1_bridge_adapter

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

var OptimismL1BridgeAdapterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIL1StandardBridge\",\"name\":\"l1Bridge\",\"type\":\"address\"},{\"internalType\":\"contractIWrappedNative\",\"name\":\"wrappedNative\",\"type\":\"address\"},{\"internalType\":\"contractIOptimismPortal\",\"name\":\"optimismPortal\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BridgeAddressCannotBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"wanted\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"got\",\"type\":\"uint256\"}],\"name\":\"InsufficientEthValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidFinalizationAction\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"MsgShouldNotContainValue\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"msgValue\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"MsgValueDoesNotMatchAmount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"finalizeWithdrawERC20\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBridgeFeeInNative\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getL1Bridge\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOptimismPortal\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getWrappedNative\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"sendERC20\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60e0604052600080546001600160401b03191690553480156200002157600080fd5b50604051620016e7380380620016e78339810160408190526200004491620000cb565b6001600160a01b03831615806200006257506001600160a01b038216155b806200007557506001600160a01b038116155b156200009457604051635e9c404d60e11b815260040160405180910390fd5b6001600160a01b0392831660805290821660a0521660c0526200011f565b6001600160a01b0381168114620000c857600080fd5b50565b600080600060608486031215620000e157600080fd5b8351620000ee81620000b2565b60208501519093506200010181620000b2565b60408501519092506200011481620000b2565b809150509250925092565b60805160a05160c0516115686200017f6000396000818160d5015281816105f901526106da01526000818161017c0152818161035401526103d40152600081816101490152818161048001528181610515015261057701526115686000f3fe6080604052600436106100695760003560e01c8063a71d98b711610043578063a71d98b71461011a578063c86d5bdd1461013a578063e861e9071461016d57600080fd5b80632e4b1fc91461007557806338314bb21461009657806354fd969f146100c657600080fd5b3661007057005b600080fd5b34801561008157600080fd5b50604051600081526020015b60405180910390f35b3480156100a257600080fd5b506100b66100b1366004610c62565b6101a0565b604051901515815260200161008d565b3480156100d257600080fd5b507f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161008d565b61012d610128366004610cc7565b61027f565b60405161008d9190610dba565b34801561014657600080fd5b507f00000000000000000000000000000000000000000000000000000000000000006100f5565b34801561017957600080fd5b507f00000000000000000000000000000000000000000000000000000000000000006100f5565b6000806101af83850185610ee4565b90506000815160018111156101c6576101c6610faa565b036101fb57600081602001518060200190518101906101e5919061116c565b90506101f0816105f7565b600092505050610277565b60018151600181111561021057610210610faa565b03610245576000816020015180602001905181019061022f919061126a565b905061023a8161069b565b600192505050610277565b6040517fee2ef09800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b949350505050565b60606102a373ffffffffffffffffffffffffffffffffffffffff881633308761070e565b34156102e2576040517f2543d86e0000000000000000000000000000000000000000000000000000000081523460048201526024015b60405180910390fd5b6000805467ffffffffffffffff1681806102fb836112ed565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550604051602001610341919067ffffffffffffffff91909116815260200190565b60405160208183030381529060405290507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168873ffffffffffffffffffffffffffffffffffffffff16036104f9576040517f2e1a7d4d000000000000000000000000000000000000000000000000000000008152600481018690527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690632e1a7d4d90602401600060405180830381600087803b15801561042d57600080fd5b505af1158015610441573d6000803e3d6000fd5b50506040517f9a2ac6d500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169250639a2ac6d5915087906104be908a90600090879060040161133b565b6000604051808303818588803b1580156104d757600080fd5b505af11580156104eb573d6000803e3d6000fd5b5050505050809150506105ed565b61053a73ffffffffffffffffffffffffffffffffffffffff89167f0000000000000000000000000000000000000000000000000000000000000000876107f0565b6040517f838b252000000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169063838b2520906105b7908b908b908b908b90600090899060040161137f565b600060405180830381600087803b1580156105d157600080fd5b505af11580156105e5573d6000803e3d6000fd5b509293505050505b9695505050505050565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16634870496f82600001518360200151846040015185606001516040518563ffffffff1660e01b81526004016106669493929190611435565b600060405180830381600087803b15801561068057600080fd5b505af1158015610694573d6000803e3d6000fd5b5050505050565b80516040517f8c3152e900000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001691638c3152e99161066691906004016114f1565b60405173ffffffffffffffffffffffffffffffffffffffff808516602483015283166044820152606481018290526107ea9085907f23b872dd00000000000000000000000000000000000000000000000000000000906084015b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152610977565b50505050565b80158061089057506040517fdd62ed3e00000000000000000000000000000000000000000000000000000000815230600482015273ffffffffffffffffffffffffffffffffffffffff838116602483015284169063dd62ed3e90604401602060405180830381865afa15801561086a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061088e9190611504565b155b61091c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603660248201527f5361666545524332303a20617070726f76652066726f6d206e6f6e2d7a65726f60448201527f20746f206e6f6e2d7a65726f20616c6c6f77616e63650000000000000000000060648201526084016102d9565b60405173ffffffffffffffffffffffffffffffffffffffff83166024820152604481018290526109729084907f095ea7b30000000000000000000000000000000000000000000000000000000090606401610768565b505050565b60006109d9826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff16610a839092919063ffffffff16565b80519091501561097257808060200190518101906109f7919061151d565b610972576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f7420737563636565640000000000000000000000000000000000000000000060648201526084016102d9565b60606102778484600085856000808673ffffffffffffffffffffffffffffffffffffffff168587604051610ab7919061153f565b60006040518083038185875af1925050503d8060008114610af4576040519150601f19603f3d011682016040523d82523d6000602084013e610af9565b606091505b5091509150610b0a87838387610b15565b979650505050505050565b60608315610bab578251600003610ba45773ffffffffffffffffffffffffffffffffffffffff85163b610ba4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016102d9565b5081610277565b6102778383815115610bc05781518083602001fd5b806040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102d99190610dba565b73ffffffffffffffffffffffffffffffffffffffff81168114610c1657600080fd5b50565b60008083601f840112610c2b57600080fd5b50813567ffffffffffffffff811115610c4357600080fd5b602083019150836020828501011115610c5b57600080fd5b9250929050565b60008060008060608587031215610c7857600080fd5b8435610c8381610bf4565b93506020850135610c9381610bf4565b9250604085013567ffffffffffffffff811115610caf57600080fd5b610cbb87828801610c19565b95989497509550505050565b60008060008060008060a08789031215610ce057600080fd5b8635610ceb81610bf4565b95506020870135610cfb81610bf4565b94506040870135610d0b81610bf4565b935060608701359250608087013567ffffffffffffffff811115610d2e57600080fd5b610d3a89828a01610c19565b979a9699509497509295939492505050565b60005b83811015610d67578181015183820152602001610d4f565b50506000910152565b60008151808452610d88816020860160208601610d4c565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000610dcd6020830184610d70565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff81118282101715610e2657610e26610dd4565b60405290565b6040516080810167ffffffffffffffff81118282101715610e2657610e26610dd4565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610e9657610e96610dd4565b604052919050565b600067ffffffffffffffff821115610eb857610eb8610dd4565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b60006020808385031215610ef757600080fd5b823567ffffffffffffffff80821115610f0f57600080fd5b9084019060408287031215610f2357600080fd5b610f2b610e03565b823560028110610f3a57600080fd5b81528284013582811115610f4d57600080fd5b80840193505086601f840112610f6257600080fd5b82359150610f77610f7283610e9e565b610e4f565b8281528785848601011115610f8b57600080fd5b8285850186830137600092810185019290925292830152509392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600082601f830112610fea57600080fd5b8151610ff8610f7282610e9e565b81815284602083860101111561100d57600080fd5b610277826020830160208701610d4c565b600060c0828403121561103057600080fd5b60405160c0810167ffffffffffffffff828210818311171561105457611054610dd4565b81604052829350845183526020850151915061106f82610bf4565b8160208401526040850151915061108582610bf4565b816040840152606085015160608401526080850151608084015260a08501519150808211156110b357600080fd5b506110c085828601610fd9565b60a0830152505092915050565b600082601f8301126110de57600080fd5b8151602067ffffffffffffffff808311156110fb576110fb610dd4565b8260051b61110a838201610e4f565b938452858101830193838101908886111561112457600080fd5b84880192505b85831015611160578251848111156111425760008081fd5b6111508a87838c0101610fd9565b835250918401919084019061112a565b98975050505050505050565b60006020828403121561117e57600080fd5b815167ffffffffffffffff8082111561119657600080fd5b9083019081850360e08112156111ab57600080fd5b6111b3610e2c565b8351838111156111c257600080fd5b6111ce8882870161101e565b8252506020840151602082015260807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08301121561120b57600080fd5b611213610e2c565b6040858101518252606080870151602084015260808701518284015260a08701519083015282015260c084015191508282111561124f57600080fd5b61125b878386016110cd565b60608201529695505050505050565b60006020828403121561127c57600080fd5b815167ffffffffffffffff8082111561129457600080fd5b90830190602082860312156112a857600080fd5b6040516020810181811083821117156112c3576112c3610dd4565b6040528251828111156112d557600080fd5b6112e18782860161101e565b82525095945050505050565b600067ffffffffffffffff808316818103611331577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6001019392505050565b73ffffffffffffffffffffffffffffffffffffffff8416815263ffffffff831660208201526060604082015260006113766060830184610d70565b95945050505050565b600073ffffffffffffffffffffffffffffffffffffffff8089168352808816602084015280871660408401525084606083015263ffffffff8416608083015260c060a083015261116060c0830184610d70565b805182526000602082015173ffffffffffffffffffffffffffffffffffffffff80821660208601528060408501511660408601525050606082015160608401526080820151608084015260a082015160c060a085015261027760c0850182610d70565b60e08152600061144860e08301876113d2565b602086818501528551604085015280860151606085015260408601516080850152606086015160a085015283820360c08501528185518084528284019150828160051b85010183880160005b838110156114e0577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08784030185526114ce838351610d70565b94860194925090850190600101611494565b50909b9a5050505050505050505050565b602081526000610dcd60208301846113d2565b60006020828403121561151657600080fd5b5051919050565b60006020828403121561152f57600080fd5b81518015158114610dcd57600080fd5b60008251611551818460208701610d4c565b919091019291505056fea164736f6c6343000818000a",
}

var OptimismL1BridgeAdapterABI = OptimismL1BridgeAdapterMetaData.ABI

var OptimismL1BridgeAdapterBin = OptimismL1BridgeAdapterMetaData.Bin

func DeployOptimismL1BridgeAdapter(auth *bind.TransactOpts, backend bind.ContractBackend, l1Bridge common.Address, wrappedNative common.Address, optimismPortal common.Address) (common.Address, *types.Transaction, *OptimismL1BridgeAdapter, error) {
	parsed, err := OptimismL1BridgeAdapterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OptimismL1BridgeAdapterBin), backend, l1Bridge, wrappedNative, optimismPortal)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OptimismL1BridgeAdapter{address: address, abi: *parsed, OptimismL1BridgeAdapterCaller: OptimismL1BridgeAdapterCaller{contract: contract}, OptimismL1BridgeAdapterTransactor: OptimismL1BridgeAdapterTransactor{contract: contract}, OptimismL1BridgeAdapterFilterer: OptimismL1BridgeAdapterFilterer{contract: contract}}, nil
}

type OptimismL1BridgeAdapter struct {
	address common.Address
	abi     abi.ABI
	OptimismL1BridgeAdapterCaller
	OptimismL1BridgeAdapterTransactor
	OptimismL1BridgeAdapterFilterer
}

type OptimismL1BridgeAdapterCaller struct {
	contract *bind.BoundContract
}

type OptimismL1BridgeAdapterTransactor struct {
	contract *bind.BoundContract
}

type OptimismL1BridgeAdapterFilterer struct {
	contract *bind.BoundContract
}

type OptimismL1BridgeAdapterSession struct {
	Contract     *OptimismL1BridgeAdapter
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OptimismL1BridgeAdapterCallerSession struct {
	Contract *OptimismL1BridgeAdapterCaller
	CallOpts bind.CallOpts
}

type OptimismL1BridgeAdapterTransactorSession struct {
	Contract     *OptimismL1BridgeAdapterTransactor
	TransactOpts bind.TransactOpts
}

type OptimismL1BridgeAdapterRaw struct {
	Contract *OptimismL1BridgeAdapter
}

type OptimismL1BridgeAdapterCallerRaw struct {
	Contract *OptimismL1BridgeAdapterCaller
}

type OptimismL1BridgeAdapterTransactorRaw struct {
	Contract *OptimismL1BridgeAdapterTransactor
}

func NewOptimismL1BridgeAdapter(address common.Address, backend bind.ContractBackend) (*OptimismL1BridgeAdapter, error) {
	abi, err := abi.JSON(strings.NewReader(OptimismL1BridgeAdapterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOptimismL1BridgeAdapter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OptimismL1BridgeAdapter{address: address, abi: abi, OptimismL1BridgeAdapterCaller: OptimismL1BridgeAdapterCaller{contract: contract}, OptimismL1BridgeAdapterTransactor: OptimismL1BridgeAdapterTransactor{contract: contract}, OptimismL1BridgeAdapterFilterer: OptimismL1BridgeAdapterFilterer{contract: contract}}, nil
}

func NewOptimismL1BridgeAdapterCaller(address common.Address, caller bind.ContractCaller) (*OptimismL1BridgeAdapterCaller, error) {
	contract, err := bindOptimismL1BridgeAdapter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismL1BridgeAdapterCaller{contract: contract}, nil
}

func NewOptimismL1BridgeAdapterTransactor(address common.Address, transactor bind.ContractTransactor) (*OptimismL1BridgeAdapterTransactor, error) {
	contract, err := bindOptimismL1BridgeAdapter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismL1BridgeAdapterTransactor{contract: contract}, nil
}

func NewOptimismL1BridgeAdapterFilterer(address common.Address, filterer bind.ContractFilterer) (*OptimismL1BridgeAdapterFilterer, error) {
	contract, err := bindOptimismL1BridgeAdapter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OptimismL1BridgeAdapterFilterer{contract: contract}, nil
}

func bindOptimismL1BridgeAdapter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OptimismL1BridgeAdapterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismL1BridgeAdapter.Contract.OptimismL1BridgeAdapterCaller.contract.Call(opts, result, method, params...)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.Contract.OptimismL1BridgeAdapterTransactor.contract.Transfer(opts)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.Contract.OptimismL1BridgeAdapterTransactor.contract.Transact(opts, method, params...)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismL1BridgeAdapter.Contract.contract.Call(opts, result, method, params...)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.Contract.contract.Transfer(opts)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.Contract.contract.Transact(opts, method, params...)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterCaller) GetBridgeFeeInNative(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OptimismL1BridgeAdapter.contract.Call(opts, &out, "getBridgeFeeInNative")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterSession) GetBridgeFeeInNative() (*big.Int, error) {
	return _OptimismL1BridgeAdapter.Contract.GetBridgeFeeInNative(&_OptimismL1BridgeAdapter.CallOpts)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterCallerSession) GetBridgeFeeInNative() (*big.Int, error) {
	return _OptimismL1BridgeAdapter.Contract.GetBridgeFeeInNative(&_OptimismL1BridgeAdapter.CallOpts)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterCaller) GetL1Bridge(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OptimismL1BridgeAdapter.contract.Call(opts, &out, "getL1Bridge")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterSession) GetL1Bridge() (common.Address, error) {
	return _OptimismL1BridgeAdapter.Contract.GetL1Bridge(&_OptimismL1BridgeAdapter.CallOpts)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterCallerSession) GetL1Bridge() (common.Address, error) {
	return _OptimismL1BridgeAdapter.Contract.GetL1Bridge(&_OptimismL1BridgeAdapter.CallOpts)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterCaller) GetOptimismPortal(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OptimismL1BridgeAdapter.contract.Call(opts, &out, "getOptimismPortal")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterSession) GetOptimismPortal() (common.Address, error) {
	return _OptimismL1BridgeAdapter.Contract.GetOptimismPortal(&_OptimismL1BridgeAdapter.CallOpts)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterCallerSession) GetOptimismPortal() (common.Address, error) {
	return _OptimismL1BridgeAdapter.Contract.GetOptimismPortal(&_OptimismL1BridgeAdapter.CallOpts)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterCaller) GetWrappedNative(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OptimismL1BridgeAdapter.contract.Call(opts, &out, "getWrappedNative")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterSession) GetWrappedNative() (common.Address, error) {
	return _OptimismL1BridgeAdapter.Contract.GetWrappedNative(&_OptimismL1BridgeAdapter.CallOpts)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterCallerSession) GetWrappedNative() (common.Address, error) {
	return _OptimismL1BridgeAdapter.Contract.GetWrappedNative(&_OptimismL1BridgeAdapter.CallOpts)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterTransactor) FinalizeWithdrawERC20(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, data []byte) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.contract.Transact(opts, "finalizeWithdrawERC20", arg0, arg1, data)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterSession) FinalizeWithdrawERC20(arg0 common.Address, arg1 common.Address, data []byte) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.Contract.FinalizeWithdrawERC20(&_OptimismL1BridgeAdapter.TransactOpts, arg0, arg1, data)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterTransactorSession) FinalizeWithdrawERC20(arg0 common.Address, arg1 common.Address, data []byte) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.Contract.FinalizeWithdrawERC20(&_OptimismL1BridgeAdapter.TransactOpts, arg0, arg1, data)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterTransactor) SendERC20(opts *bind.TransactOpts, localToken common.Address, remoteToken common.Address, recipient common.Address, amount *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.contract.Transact(opts, "sendERC20", localToken, remoteToken, recipient, amount, arg4)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterSession) SendERC20(localToken common.Address, remoteToken common.Address, recipient common.Address, amount *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.Contract.SendERC20(&_OptimismL1BridgeAdapter.TransactOpts, localToken, remoteToken, recipient, amount, arg4)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterTransactorSession) SendERC20(localToken common.Address, remoteToken common.Address, recipient common.Address, amount *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.Contract.SendERC20(&_OptimismL1BridgeAdapter.TransactOpts, localToken, remoteToken, recipient, amount, arg4)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.contract.RawTransact(opts, nil)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterSession) Receive() (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.Contract.Receive(&_OptimismL1BridgeAdapter.TransactOpts)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterTransactorSession) Receive() (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.Contract.Receive(&_OptimismL1BridgeAdapter.TransactOpts)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapter) Address() common.Address {
	return _OptimismL1BridgeAdapter.address
}

type OptimismL1BridgeAdapterInterface interface {
	GetBridgeFeeInNative(opts *bind.CallOpts) (*big.Int, error)

	GetL1Bridge(opts *bind.CallOpts) (common.Address, error)

	GetOptimismPortal(opts *bind.CallOpts) (common.Address, error)

	GetWrappedNative(opts *bind.CallOpts) (common.Address, error)

	FinalizeWithdrawERC20(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, data []byte) (*types.Transaction, error)

	SendERC20(opts *bind.TransactOpts, localToken common.Address, remoteToken common.Address, recipient common.Address, amount *big.Int, arg4 []byte) (*types.Transaction, error)

	Receive(opts *bind.TransactOpts) (*types.Transaction, error)

	Address() common.Address
}
