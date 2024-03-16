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
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIL1StandardBridge\",\"name\":\"l1Bridge\",\"type\":\"address\"},{\"internalType\":\"contractIWrappedNative\",\"name\":\"wrappedNative\",\"type\":\"address\"},{\"internalType\":\"contractIOptimismPortal\",\"name\":\"optimismPortal\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BridgeAddressCannotBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"wanted\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"got\",\"type\":\"uint256\"}],\"name\":\"InsufficientEthValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidFinalizationAction\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"MsgShouldNotContainValue\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"msgValue\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"MsgValueDoesNotMatchAmount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"finalizeWithdrawERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBridgeFeeInNative\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getL1Bridge\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOptimismPortal\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getWrappedNative\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"sendERC20\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60e0604052600080546001600160401b03191690553480156200002157600080fd5b506040516200175e3803806200175e8339810160408190526200004491620000cb565b6001600160a01b03831615806200006257506001600160a01b038216155b806200007557506001600160a01b038116155b156200009457604051635e9c404d60e11b815260040160405180910390fd5b6001600160a01b0392831660805290821660a0521660c0526200011f565b6001600160a01b0381168114620000c857600080fd5b50565b600080600060608486031215620000e157600080fd5b8351620000ee81620000b2565b60208501519093506200010181620000b2565b60408501519092506200011481620000b2565b809150509250925092565b60805160a05160c0516115df6200017f6000396000818160c7015281816105d901526106b301526000818161016e0152818161033401526103b401526000818161013b01528181610460015281816104f5015261055701526115df6000f3fe6080604052600436106100695760003560e01c8063a71d98b711610043578063a71d98b71461010c578063c86d5bdd1461012c578063e861e9071461015f57600080fd5b80632e4b1fc91461007557806338314bb21461009657806354fd969f146100b857600080fd5b3661007057005b600080fd5b34801561008157600080fd5b50604051600081526020015b60405180910390f35b3480156100a257600080fd5b506100b66100b1366004610cd9565b610192565b005b3480156100c457600080fd5b507f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161008d565b61011f61011a366004610d3e565b61025f565b60405161008d9190610e31565b34801561013857600080fd5b507f00000000000000000000000000000000000000000000000000000000000000006100e7565b34801561016b57600080fd5b507f00000000000000000000000000000000000000000000000000000000000000006100e7565b60006101a082840184610f5b565b90506000815160018111156101b7576101b7611021565b036101e757600081602001518060200190518101906101d691906111e3565b90506101e1816105d7565b50610258565b6001815160018111156101fc576101fc611021565b03610226576000816020015180602001905181019061021b91906112e1565b90506101e181610674565b6040517fee2ef09800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050505050565b606061028373ffffffffffffffffffffffffffffffffffffffff88163330876106e7565b34156102c2576040517f2543d86e0000000000000000000000000000000000000000000000000000000081523460048201526024015b60405180910390fd5b6000805467ffffffffffffffff1681806102db83611364565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550604051602001610321919067ffffffffffffffff91909116815260200190565b60405160208183030381529060405290507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168873ffffffffffffffffffffffffffffffffffffffff16036104d9576040517f2e1a7d4d000000000000000000000000000000000000000000000000000000008152600481018690527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690632e1a7d4d90602401600060405180830381600087803b15801561040d57600080fd5b505af1158015610421573d6000803e3d6000fd5b50506040517f9a2ac6d500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169250639a2ac6d59150879061049e908a9060009087906004016113b2565b6000604051808303818588803b1580156104b757600080fd5b505af11580156104cb573d6000803e3d6000fd5b5050505050809150506105cd565b61051a73ffffffffffffffffffffffffffffffffffffffff89167f0000000000000000000000000000000000000000000000000000000000000000876107c9565b6040517f838b252000000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169063838b252090610597908b908b908b908b9060009089906004016113f6565b600060405180830381600087803b1580156105b157600080fd5b505af11580156105c5573d6000803e3d6000fd5b509293505050505b9695505050505050565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16634870496f82600001518360200151846040015185606001516040518563ffffffff1660e01b815260040161064694939291906114ac565b600060405180830381600087803b15801561066057600080fd5b505af1158015610258573d6000803e3d6000fd5b80516040517f8c3152e900000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001691638c3152e9916106469190600401611568565b60405173ffffffffffffffffffffffffffffffffffffffff808516602483015283166044820152606481018290526107c39085907f23b872dd00000000000000000000000000000000000000000000000000000000906084015b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152610950565b50505050565b80158061086957506040517fdd62ed3e00000000000000000000000000000000000000000000000000000000815230600482015273ffffffffffffffffffffffffffffffffffffffff838116602483015284169063dd62ed3e90604401602060405180830381865afa158015610843573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610867919061157b565b155b6108f5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603660248201527f5361666545524332303a20617070726f76652066726f6d206e6f6e2d7a65726f60448201527f20746f206e6f6e2d7a65726f20616c6c6f77616e63650000000000000000000060648201526084016102b9565b60405173ffffffffffffffffffffffffffffffffffffffff831660248201526044810182905261094b9084907f095ea7b30000000000000000000000000000000000000000000000000000000090606401610741565b505050565b60006109b2826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff16610a5c9092919063ffffffff16565b80519091501561094b57808060200190518101906109d09190611594565b61094b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f7420737563636565640000000000000000000000000000000000000000000060648201526084016102b9565b6060610a6b8484600085610a73565b949350505050565b606082471015610b05576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c000000000000000000000000000000000000000000000000000060648201526084016102b9565b6000808673ffffffffffffffffffffffffffffffffffffffff168587604051610b2e91906115b6565b60006040518083038185875af1925050503d8060008114610b6b576040519150601f19603f3d011682016040523d82523d6000602084013e610b70565b606091505b5091509150610b8187838387610b8c565b979650505050505050565b60608315610c22578251600003610c1b5773ffffffffffffffffffffffffffffffffffffffff85163b610c1b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016102b9565b5081610a6b565b610a6b8383815115610c375781518083602001fd5b806040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102b99190610e31565b73ffffffffffffffffffffffffffffffffffffffff81168114610c8d57600080fd5b50565b60008083601f840112610ca257600080fd5b50813567ffffffffffffffff811115610cba57600080fd5b602083019150836020828501011115610cd257600080fd5b9250929050565b60008060008060608587031215610cef57600080fd5b8435610cfa81610c6b565b93506020850135610d0a81610c6b565b9250604085013567ffffffffffffffff811115610d2657600080fd5b610d3287828801610c90565b95989497509550505050565b60008060008060008060a08789031215610d5757600080fd5b8635610d6281610c6b565b95506020870135610d7281610c6b565b94506040870135610d8281610c6b565b935060608701359250608087013567ffffffffffffffff811115610da557600080fd5b610db189828a01610c90565b979a9699509497509295939492505050565b60005b83811015610dde578181015183820152602001610dc6565b50506000910152565b60008151808452610dff816020860160208601610dc3565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000610e446020830184610de7565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff81118282101715610e9d57610e9d610e4b565b60405290565b6040516080810167ffffffffffffffff81118282101715610e9d57610e9d610e4b565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610f0d57610f0d610e4b565b604052919050565b600067ffffffffffffffff821115610f2f57610f2f610e4b565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b60006020808385031215610f6e57600080fd5b823567ffffffffffffffff80821115610f8657600080fd5b9084019060408287031215610f9a57600080fd5b610fa2610e7a565b823560028110610fb157600080fd5b81528284013582811115610fc457600080fd5b80840193505086601f840112610fd957600080fd5b82359150610fee610fe983610f15565b610ec6565b828152878584860101111561100257600080fd5b8285850186830137600092810185019290925292830152509392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600082601f83011261106157600080fd5b815161106f610fe982610f15565b81815284602083860101111561108457600080fd5b610a6b826020830160208701610dc3565b600060c082840312156110a757600080fd5b60405160c0810167ffffffffffffffff82821081831117156110cb576110cb610e4b565b8160405282935084518352602085015191506110e682610c6b565b816020840152604085015191506110fc82610c6b565b816040840152606085015160608401526080850151608084015260a085015191508082111561112a57600080fd5b5061113785828601611050565b60a0830152505092915050565b600082601f83011261115557600080fd5b8151602067ffffffffffffffff8083111561117257611172610e4b565b8260051b611181838201610ec6565b938452858101830193838101908886111561119b57600080fd5b84880192505b858310156111d7578251848111156111b95760008081fd5b6111c78a87838c0101611050565b83525091840191908401906111a1565b98975050505050505050565b6000602082840312156111f557600080fd5b815167ffffffffffffffff8082111561120d57600080fd5b9083019081850360e081121561122257600080fd5b61122a610ea3565b83518381111561123957600080fd5b61124588828701611095565b8252506020840151602082015260807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08301121561128257600080fd5b61128a610ea3565b6040858101518252606080870151602084015260808701518284015260a08701519083015282015260c08401519150828211156112c657600080fd5b6112d287838601611144565b60608201529695505050505050565b6000602082840312156112f357600080fd5b815167ffffffffffffffff8082111561130b57600080fd5b908301906020828603121561131f57600080fd5b60405160208101818110838211171561133a5761133a610e4b565b60405282518281111561134c57600080fd5b61135887828601611095565b82525095945050505050565b600067ffffffffffffffff8083168181036113a8577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6001019392505050565b73ffffffffffffffffffffffffffffffffffffffff8416815263ffffffff831660208201526060604082015260006113ed6060830184610de7565b95945050505050565b600073ffffffffffffffffffffffffffffffffffffffff8089168352808816602084015280871660408401525084606083015263ffffffff8416608083015260c060a08301526111d760c0830184610de7565b805182526000602082015173ffffffffffffffffffffffffffffffffffffffff80821660208601528060408501511660408601525050606082015160608401526080820151608084015260a082015160c060a0850152610a6b60c0850182610de7565b60e0815260006114bf60e0830187611449565b602086818501528551604085015280860151606085015260408601516080850152606086015160a085015283820360c08501528185518084528284019150828160051b85010183880160005b83811015611557577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0878403018552611545838351610de7565b9486019492509085019060010161150b565b50909b9a5050505050505050505050565b602081526000610e446020830184611449565b60006020828403121561158d57600080fd5b5051919050565b6000602082840312156115a657600080fd5b81518015158114610e4457600080fd5b600082516115c8818460208701610dc3565b919091019291505056fea164736f6c6343000813000a",
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
