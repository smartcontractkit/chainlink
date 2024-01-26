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

type IL1CrossDomainMessengerL2MessageInclusionProof struct {
	StateRoot            [32]byte
	StateRootBatchHeader LibOVMCodecChainBatchHeader
	StateRootProof       LibOVMCodecChainInclusionProof
	StateTrieWitness     []byte
	StorageTrieWitness   []byte
}

type LibOVMCodecChainBatchHeader struct {
	BatchIndex        *big.Int
	BatchRoot         [32]byte
	BatchSize         *big.Int
	PrevTotalElements *big.Int
	ExtraData         []byte
}

type LibOVMCodecChainInclusionProof struct {
	Index    *big.Int
	Siblings [][32]byte
}

var OptimismL1BridgeAdapterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIL1StandardBridge\",\"name\":\"l1Bridge\",\"type\":\"address\"},{\"internalType\":\"contractIWrappedNative\",\"name\":\"wrappedNative\",\"type\":\"address\"},{\"internalType\":\"contractIL1CrossDomainMessenger\",\"name\":\"l1CrossDomainMessenger\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BridgeAddressCannotBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"wanted\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"got\",\"type\":\"uint256\"}],\"name\":\"InsufficientEthValue\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"MsgShouldNotContainValue\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"msgValue\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"MsgValueDoesNotMatchAmount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"depositNativeToL2\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"remoteSender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"localReceiver\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"finalizeWithdrawERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"finalizeWithdrawNativeFromL2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBridgeFeeInNative\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getWrappedNative\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"messageNonce\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"batchIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"batchRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"batchSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"prevTotalElements\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"internalType\":\"structLib_OVMCodec.ChainBatchHeader\",\"name\":\"stateRootBatchHeader\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"siblings\",\"type\":\"bytes32[]\"}],\"internalType\":\"structLib_OVMCodec.ChainInclusionProof\",\"name\":\"stateRootProof\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"stateTrieWitness\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"storageTrieWitness\",\"type\":\"bytes\"}],\"internalType\":\"structIL1CrossDomainMessenger.L2MessageInclusionProof\",\"name\":\"proof\",\"type\":\"tuple\"}],\"name\":\"relayMessageFromL2ToL1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"sendERC20\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x60e0604052600080546001600160401b03191690553480156200002157600080fd5b50604051620015ab380380620015ab8339810160408190526200004491620000b7565b6001600160a01b03831615806200006257506001600160a01b038216155b156200008157604051635e9c404d60e11b815260040160405180910390fd5b6001600160a01b03928316608052821660a0521660c0526200010b565b6001600160a01b0381168114620000b457600080fd5b50565b600080600060608486031215620000cd57600080fd5b8351620000da816200009e565b6020850151909350620000ed816200009e565b604085015190925062000100816200009e565b809150509250925092565b60805160a05160c05161144662000165600039600081816101320152818161032d01526103ad015260006106a60152600081816101b9015281816102890152818161047a0152818161051201526107d801526114466000f3fe6080604052600436106100705760003560e01c806379a35b4b1161004e57806379a35b4b146100d85780638b2e4a2c146100f8578063e861e9071461010b578063f2bfa1e11461015c57600080fd5b806318b3050c146100755780632e4b1fc91461009757806338314bb2146100b8575b600080fd5b34801561008157600080fd5b50610095610090366004610b2e565b61017c565b005b3480156100a357600080fd5b50604051600081526020015b60405180910390f35b3480156100c457600080fd5b506100956100d3366004610b9d565b61022f565b6100eb6100e6366004610bfe565b6102c8565b6040516100af9190610cb7565b610095610106366004610cd1565b610619565b34801561011757600080fd5b5060405173ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001681526020016100af565b34801561016857600080fd5b50610095610177366004610f7f565b610669565b6040517f1532ec3400000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001690631532ec34906101f690889088908890889088906004016110fd565b600060405180830381600087803b15801561021057600080fd5b505af1158015610224573d6000803e3d6000fd5b505050505050505050565b600061023d8284018461113d565b8051602082015160408084015190517fa9f9e67500000000000000000000000000000000000000000000000000000000815293945073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169363a9f9e675936101f693909290918b918b918b908b906004016111a5565b60606102ec73ffffffffffffffffffffffffffffffffffffffff86163330856106e3565b341561032b576040517f2543d86e0000000000000000000000000000000000000000000000000000000081523460048201526024015b60405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff160361043d576040517f2e1a7d4d000000000000000000000000000000000000000000000000000000008152600481018390527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690632e1a7d4d90602401600060405180830381600087803b15801561040657600080fd5b505af115801561041a573d6000803e3d6000fd5b50505050610428838361077e565b50604080516020810190915260008152610611565b6040517f095ea7b300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000811660048301526024820184905286169063095ea7b3906044016020604051808303816000875af11580156104d2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104f69190611202565b506000805473ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169163838b252091889188918891889167ffffffffffffffff16818061055683611224565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555060405160200161059c919067ffffffffffffffff91909116815260200190565b6040516020818303038152906040526040518763ffffffff1660e01b81526004016105cc96959493929190611272565b600060405180830381600087803b1580156105e657600080fd5b505af11580156105fa573d6000803e3d6000fd5b505050506040518060200160405280600081525090505b949350505050565b80341461065b576040517f03da4d2300000000000000000000000000000000000000000000000000000000815234600482015260248101829052604401610322565b610665828261077e565b5050565b6040517fd7fd19dd00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169063d7fd19dd906101f69088908890889088908890600401611328565b6040805173ffffffffffffffffffffffffffffffffffffffff85811660248301528416604482015260648082018490528251808303909101815260849091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f23b872dd0000000000000000000000000000000000000000000000000000000017905261077890859061083a565b50505050565b6040517f9a2ac6d500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff83811660048301526000602483018190526060604484015260648301527f00000000000000000000000000000000000000000000000000000000000000001690639a2ac6d59083906084016000604051808303818588803b15801561081d57600080fd5b505af1158015610831573d6000803e3d6000fd5b50505050505050565b600061089c826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff1661094b9092919063ffffffff16565b80519091501561094657808060200190518101906108ba9190611202565b610946576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152608401610322565b505050565b60606106118484600085856000808673ffffffffffffffffffffffffffffffffffffffff16858760405161097f919061141d565b60006040518083038185875af1925050503d80600081146109bc576040519150601f19603f3d011682016040523d82523d6000602084013e6109c1565b606091505b50915091506109d2878383876109dd565b979650505050505050565b60608315610a73578251600003610a6c5773ffffffffffffffffffffffffffffffffffffffff85163b610a6c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610322565b5081610611565b6106118383815115610a885781518083602001fd5b806040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103229190610cb7565b803573ffffffffffffffffffffffffffffffffffffffff81168114610ae057600080fd5b919050565b60008083601f840112610af757600080fd5b50813567ffffffffffffffff811115610b0f57600080fd5b602083019150836020828501011115610b2757600080fd5b9250929050565b600080600080600060808688031215610b4657600080fd5b610b4f86610abc565b9450610b5d60208701610abc565b935060408601359250606086013567ffffffffffffffff811115610b8057600080fd5b610b8c88828901610ae5565b969995985093965092949392505050565b60008060008060608587031215610bb357600080fd5b610bbc85610abc565b9350610bca60208601610abc565b9250604085013567ffffffffffffffff811115610be657600080fd5b610bf287828801610ae5565b95989497509550505050565b60008060008060808587031215610c1457600080fd5b610c1d85610abc565b9350610c2b60208601610abc565b9250610c3960408601610abc565b9396929550929360600135925050565b60005b83811015610c64578181015183820152602001610c4c565b50506000910152565b60008151808452610c85816020860160208601610c49565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000610cca6020830184610c6d565b9392505050565b60008060408385031215610ce457600080fd5b610ced83610abc565b946020939093013593505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160a0810167ffffffffffffffff81118282101715610d4d57610d4d610cfb565b60405290565b6040805190810167ffffffffffffffff81118282101715610d4d57610d4d610cfb565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610dbd57610dbd610cfb565b604052919050565b600082601f830112610dd657600080fd5b813567ffffffffffffffff811115610df057610df0610cfb565b610e2160207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601610d76565b818152846020838601011115610e3657600080fd5b816020850160208301376000918101602001919091529392505050565b600060a08284031215610e6557600080fd5b610e6d610d2a565b905081358152602082013560208201526040820135604082015260608201356060820152608082013567ffffffffffffffff811115610eab57600080fd5b610eb784828501610dc5565b60808301525092915050565b600060408284031215610ed557600080fd5b610edd610d53565b90508135815260208083013567ffffffffffffffff80821115610eff57600080fd5b818501915085601f830112610f1357600080fd5b813581811115610f2557610f25610cfb565b8060051b9150610f36848301610d76565b8181529183018401918481019088841115610f5057600080fd5b938501935b83851015610f6e57843582529385019390850190610f55565b808688015250505050505092915050565b600080600080600060a08688031215610f9757600080fd5b610fa086610abc565b9450610fae60208701610abc565b9350604086013567ffffffffffffffff80821115610fcb57600080fd5b610fd789838a01610dc5565b9450606088013593506080880135915080821115610ff457600080fd5b9087019060a0828a03121561100857600080fd5b611010610d2a565b8235815260208301358281111561102657600080fd5b6110328b828601610e53565b60208301525060408301358281111561104a57600080fd5b6110568b828601610ec3565b60408301525060608301358281111561106e57600080fd5b61107a8b828601610dc5565b60608301525060808301358281111561109257600080fd5b61109e8b828601610dc5565b6080830152508093505050509295509295909350565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b600073ffffffffffffffffffffffffffffffffffffffff8088168352808716602084015250846040830152608060608301526109d26080830184866110b4565b60006060828403121561114f57600080fd5b6040516060810181811067ffffffffffffffff8211171561117257611172610cfb565b60405261117e83610abc565b815261118c60208401610abc565b6020820152604083013560408201528091505092915050565b600073ffffffffffffffffffffffffffffffffffffffff808a1683528089166020840152808816604084015280871660608401525084608083015260c060a08301526111f560c0830184866110b4565b9998505050505050505050565b60006020828403121561121457600080fd5b81518015158114610cca57600080fd5b600067ffffffffffffffff808316818103611268577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6001019392505050565b600073ffffffffffffffffffffffffffffffffffffffff8089168352808816602084015280871660408401525084606083015263ffffffff8416608083015260c060a08301526112c560c0830184610c6d565b98975050505050505050565b600060408301825184526020808401516040828701528281518085526060880191508383019450600092505b8083101561131d57845182529383019360019290920191908301906112fd565b509695505050505050565b600073ffffffffffffffffffffffffffffffffffffffff808816835280871660208401525060a0604083015261136160a0830186610c6d565b846060840152828103608084015283518152602084015160a06020830152805160a0830152602081015160c0830152604081015160e083015260608101516101008301526080810151905060a06101208301526113c2610140830182610c6d565b9050604085015182820360408401526113db82826112d1565b915050606085015182820360608401526113f58282610c6d565b9150506080850151828203608084015261140f8282610c6d565b9a9950505050505050505050565b6000825161142f818460208701610c49565b919091019291505056fea164736f6c6343000813000a",
}

var OptimismL1BridgeAdapterABI = OptimismL1BridgeAdapterMetaData.ABI

var OptimismL1BridgeAdapterBin = OptimismL1BridgeAdapterMetaData.Bin

func DeployOptimismL1BridgeAdapter(auth *bind.TransactOpts, backend bind.ContractBackend, l1Bridge common.Address, wrappedNative common.Address, l1CrossDomainMessenger common.Address) (common.Address, *types.Transaction, *OptimismL1BridgeAdapter, error) {
	parsed, err := OptimismL1BridgeAdapterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OptimismL1BridgeAdapterBin), backend, l1Bridge, wrappedNative, l1CrossDomainMessenger)
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

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterTransactor) DepositNativeToL2(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.contract.Transact(opts, "depositNativeToL2", recipient, amount)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterSession) DepositNativeToL2(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.Contract.DepositNativeToL2(&_OptimismL1BridgeAdapter.TransactOpts, recipient, amount)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterTransactorSession) DepositNativeToL2(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.Contract.DepositNativeToL2(&_OptimismL1BridgeAdapter.TransactOpts, recipient, amount)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterTransactor) FinalizeWithdrawERC20(opts *bind.TransactOpts, remoteSender common.Address, localReceiver common.Address, data []byte) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.contract.Transact(opts, "finalizeWithdrawERC20", remoteSender, localReceiver, data)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterSession) FinalizeWithdrawERC20(remoteSender common.Address, localReceiver common.Address, data []byte) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.Contract.FinalizeWithdrawERC20(&_OptimismL1BridgeAdapter.TransactOpts, remoteSender, localReceiver, data)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterTransactorSession) FinalizeWithdrawERC20(remoteSender common.Address, localReceiver common.Address, data []byte) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.Contract.FinalizeWithdrawERC20(&_OptimismL1BridgeAdapter.TransactOpts, remoteSender, localReceiver, data)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterTransactor) FinalizeWithdrawNativeFromL2(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.contract.Transact(opts, "finalizeWithdrawNativeFromL2", from, to, amount, data)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterSession) FinalizeWithdrawNativeFromL2(from common.Address, to common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.Contract.FinalizeWithdrawNativeFromL2(&_OptimismL1BridgeAdapter.TransactOpts, from, to, amount, data)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterTransactorSession) FinalizeWithdrawNativeFromL2(from common.Address, to common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.Contract.FinalizeWithdrawNativeFromL2(&_OptimismL1BridgeAdapter.TransactOpts, from, to, amount, data)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterTransactor) RelayMessageFromL2ToL1(opts *bind.TransactOpts, target common.Address, sender common.Address, message []byte, messageNonce *big.Int, proof IL1CrossDomainMessengerL2MessageInclusionProof) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.contract.Transact(opts, "relayMessageFromL2ToL1", target, sender, message, messageNonce, proof)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterSession) RelayMessageFromL2ToL1(target common.Address, sender common.Address, message []byte, messageNonce *big.Int, proof IL1CrossDomainMessengerL2MessageInclusionProof) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.Contract.RelayMessageFromL2ToL1(&_OptimismL1BridgeAdapter.TransactOpts, target, sender, message, messageNonce, proof)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterTransactorSession) RelayMessageFromL2ToL1(target common.Address, sender common.Address, message []byte, messageNonce *big.Int, proof IL1CrossDomainMessengerL2MessageInclusionProof) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.Contract.RelayMessageFromL2ToL1(&_OptimismL1BridgeAdapter.TransactOpts, target, sender, message, messageNonce, proof)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterTransactor) SendERC20(opts *bind.TransactOpts, localToken common.Address, remoteToken common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.contract.Transact(opts, "sendERC20", localToken, remoteToken, recipient, amount)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterSession) SendERC20(localToken common.Address, remoteToken common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.Contract.SendERC20(&_OptimismL1BridgeAdapter.TransactOpts, localToken, remoteToken, recipient, amount)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterTransactorSession) SendERC20(localToken common.Address, remoteToken common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.Contract.SendERC20(&_OptimismL1BridgeAdapter.TransactOpts, localToken, remoteToken, recipient, amount)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapter) Address() common.Address {
	return _OptimismL1BridgeAdapter.address
}

type OptimismL1BridgeAdapterInterface interface {
	GetBridgeFeeInNative(opts *bind.CallOpts) (*big.Int, error)

	GetWrappedNative(opts *bind.CallOpts) (common.Address, error)

	DepositNativeToL2(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	FinalizeWithdrawERC20(opts *bind.TransactOpts, remoteSender common.Address, localReceiver common.Address, data []byte) (*types.Transaction, error)

	FinalizeWithdrawNativeFromL2(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	RelayMessageFromL2ToL1(opts *bind.TransactOpts, target common.Address, sender common.Address, message []byte, messageNonce *big.Int, proof IL1CrossDomainMessengerL2MessageInclusionProof) (*types.Transaction, error)

	SendERC20(opts *bind.TransactOpts, localToken common.Address, remoteToken common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	Address() common.Address
}
