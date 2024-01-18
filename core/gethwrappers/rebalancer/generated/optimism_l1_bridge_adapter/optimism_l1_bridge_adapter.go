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
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIL1StandardBridge\",\"name\":\"l1Bridge\",\"type\":\"address\"},{\"internalType\":\"contractIWrappedNative\",\"name\":\"wrappedNative\",\"type\":\"address\"},{\"internalType\":\"contractIL1CrossDomainMessenger\",\"name\":\"l1CrossDomainMessenger\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BridgeAddressCannotBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"wanted\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"got\",\"type\":\"uint256\"}],\"name\":\"InsufficientEthValue\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"MsgShouldNotContainValue\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"msgValue\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"MsgValueDoesNotMatchAmount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"depositNativeToL2\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"finalizeWithdrawERC20FromL2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"finalizeWithdrawNativeFromL2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBridgeFeeInNative\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getWrappedNative\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"messageNonce\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"batchIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"batchRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"batchSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"prevTotalElements\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"internalType\":\"structLib_OVMCodec.ChainBatchHeader\",\"name\":\"stateRootBatchHeader\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"siblings\",\"type\":\"bytes32[]\"}],\"internalType\":\"structLib_OVMCodec.ChainInclusionProof\",\"name\":\"stateRootProof\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"stateTrieWitness\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"storageTrieWitness\",\"type\":\"bytes\"}],\"internalType\":\"structIL1CrossDomainMessenger.L2MessageInclusionProof\",\"name\":\"proof\",\"type\":\"tuple\"}],\"name\":\"relayMessageFromL2ToL1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"l1Token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"l2Token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"sendERC20\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x60e0604052600080546001600160401b03191690553480156200002157600080fd5b506040516200166a3803806200166a8339810160408190526200004491620000b7565b6001600160a01b03831615806200006257506001600160a01b038216155b156200008157604051635e9c404d60e11b815260040160405180910390fd5b6001600160a01b03928316608052821660a0521660c0526200010b565b6001600160a01b0381168114620000b457600080fd5b50565b600080600060608486031215620000cd57600080fd5b8351620000da816200009e565b6020850151909350620000ed816200009e565b604085015190925062000100816200009e565b809150509250925092565b60805160a05160c0516115056200016560003960008181610125015281816102850152610305015260006107690152600081816101ac015281816103c20152818161045a015281816105aa01526106ed01526115056000f3fe6080604052600436106100705760003560e01c80638b2e4a2c1161004e5780638b2e4a2c146100cb578063cb7a6e16146100de578063e861e907146100fe578063f2bfa1e11461014f57600080fd5b806318b3050c146100755780632e4b1fc91461009757806379a35b4b146100b8575b600080fd5b34801561008157600080fd5b50610095610090366004610ba9565b61016f565b005b3480156100a357600080fd5b50604051600081526020015b60405180910390f35b6100956100c6366004610c18565b610222565b6100956100d9366004610c63565b61054d565b3480156100ea57600080fd5b506100956100f9366004610c8d565b610693565b34801561010a57600080fd5b5060405173ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001681526020016100af565b34801561015b57600080fd5b5061009561016a366004610f72565b61072c565b6040517f1532ec3400000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001690631532ec34906101e990889088908890889088906004016110f0565b600060405180830381600087803b15801561020357600080fd5b505af1158015610217573d6000803e3d6000fd5b505050505050505050565b61024473ffffffffffffffffffffffffffffffffffffffff85163330846107a6565b3415610283576040517f2543d86e0000000000000000000000000000000000000000000000000000000081523460048201526024015b60405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1603610385576040517f2e1a7d4d000000000000000000000000000000000000000000000000000000008152600481018290527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690632e1a7d4d90602401600060405180830381600087803b15801561035e57600080fd5b505af1158015610372573d6000803e3d6000fd5b50505050610380828261054d565b610547565b6040517f095ea7b300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000811660048301526024820183905285169063095ea7b3906044016020604051808303816000875af115801561041a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061043e9190611130565b506000805473ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169163838b252091879187918791879167ffffffffffffffff16818061049e83611159565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055506040516020016104e4919067ffffffffffffffff91909116815260200190565b6040516020818303038152906040526040518763ffffffff1660e01b815260040161051496959493929190611215565b600060405180830381600087803b15801561052e57600080fd5b505af1158015610542573d6000803e3d6000fd5b505050505b50505050565b80341461058f576040517f03da4d230000000000000000000000000000000000000000000000000000000081523460048201526024810182905260440161027a565b6000805473ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001691639a2ac6d5913491869167ffffffffffffffff1681806105ea83611159565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550604051602001610630919067ffffffffffffffff91909116815260200190565b6040516020818303038152906040526040518563ffffffff1660e01b815260040161065d93929190611274565b6000604051808303818588803b15801561067657600080fd5b505af115801561068a573d6000803e3d6000fd5b50505050505050565b60006106a1828401846112b8565b8051602082015160408084015190517fa9f9e67500000000000000000000000000000000000000000000000000000000815293945073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169363a9f9e675936101e993909290918b918b918b908b90600401611320565b6040517fd7fd19dd00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169063d7fd19dd906101e990889088908890889088906004016113d4565b6040805173ffffffffffffffffffffffffffffffffffffffff8581166024830152848116604483015260648083018590528351808403909101815260849092018352602080830180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f23b872dd0000000000000000000000000000000000000000000000000000000017905283518085019094528084527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65649084015261054792879291600091610879918516908490610928565b80519091501561092357808060200190518101906108979190611130565b610923576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f74207375636365656400000000000000000000000000000000000000000000606482015260840161027a565b505050565b6060610937848460008561093f565b949350505050565b6060824710156109d1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c0000000000000000000000000000000000000000000000000000606482015260840161027a565b6000808673ffffffffffffffffffffffffffffffffffffffff1685876040516109fa91906114c9565b60006040518083038185875af1925050503d8060008114610a37576040519150601f19603f3d011682016040523d82523d6000602084013e610a3c565b606091505b5091509150610a4d87838387610a58565b979650505050505050565b60608315610aee578251600003610ae75773ffffffffffffffffffffffffffffffffffffffff85163b610ae7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640161027a565b5081610937565b6109378383815115610b035781518083602001fd5b806040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161027a91906114e5565b803573ffffffffffffffffffffffffffffffffffffffff81168114610b5b57600080fd5b919050565b60008083601f840112610b7257600080fd5b50813567ffffffffffffffff811115610b8a57600080fd5b602083019150836020828501011115610ba257600080fd5b9250929050565b600080600080600060808688031215610bc157600080fd5b610bca86610b37565b9450610bd860208701610b37565b935060408601359250606086013567ffffffffffffffff811115610bfb57600080fd5b610c0788828901610b60565b969995985093965092949392505050565b60008060008060808587031215610c2e57600080fd5b610c3785610b37565b9350610c4560208601610b37565b9250610c5360408601610b37565b9396929550929360600135925050565b60008060408385031215610c7657600080fd5b610c7f83610b37565b946020939093013593505050565b60008060008060608587031215610ca357600080fd5b610cac85610b37565b9350610cba60208601610b37565b9250604085013567ffffffffffffffff811115610cd657600080fd5b610ce287828801610b60565b95989497509550505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160a0810167ffffffffffffffff81118282101715610d4057610d40610cee565b60405290565b6040805190810167ffffffffffffffff81118282101715610d4057610d40610cee565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610db057610db0610cee565b604052919050565b600082601f830112610dc957600080fd5b813567ffffffffffffffff811115610de357610de3610cee565b610e1460207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601610d69565b818152846020838601011115610e2957600080fd5b816020850160208301376000918101602001919091529392505050565b600060a08284031215610e5857600080fd5b610e60610d1d565b905081358152602082013560208201526040820135604082015260608201356060820152608082013567ffffffffffffffff811115610e9e57600080fd5b610eaa84828501610db8565b60808301525092915050565b600060408284031215610ec857600080fd5b610ed0610d46565b90508135815260208083013567ffffffffffffffff80821115610ef257600080fd5b818501915085601f830112610f0657600080fd5b813581811115610f1857610f18610cee565b8060051b9150610f29848301610d69565b8181529183018401918481019088841115610f4357600080fd5b938501935b83851015610f6157843582529385019390850190610f48565b808688015250505050505092915050565b600080600080600060a08688031215610f8a57600080fd5b610f9386610b37565b9450610fa160208701610b37565b9350604086013567ffffffffffffffff80821115610fbe57600080fd5b610fca89838a01610db8565b9450606088013593506080880135915080821115610fe757600080fd5b9087019060a0828a031215610ffb57600080fd5b611003610d1d565b8235815260208301358281111561101957600080fd5b6110258b828601610e46565b60208301525060408301358281111561103d57600080fd5b6110498b828601610eb6565b60408301525060608301358281111561106157600080fd5b61106d8b828601610db8565b60608301525060808301358281111561108557600080fd5b6110918b828601610db8565b6080830152508093505050509295509295909350565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b600073ffffffffffffffffffffffffffffffffffffffff808816835280871660208401525084604083015260806060830152610a4d6080830184866110a7565b60006020828403121561114257600080fd5b8151801515811461115257600080fd5b9392505050565b600067ffffffffffffffff80831681810361119d577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6001019392505050565b60005b838110156111c25781810151838201526020016111aa565b50506000910152565b600081518084526111e38160208601602086016111a7565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b600073ffffffffffffffffffffffffffffffffffffffff8089168352808816602084015280871660408401525084606083015263ffffffff8416608083015260c060a083015261126860c08301846111cb565b98975050505050505050565b73ffffffffffffffffffffffffffffffffffffffff8416815263ffffffff831660208201526060604082015260006112af60608301846111cb565b95945050505050565b6000606082840312156112ca57600080fd5b6040516060810181811067ffffffffffffffff821117156112ed576112ed610cee565b6040526112f983610b37565b815261130760208401610b37565b6020820152604083013560408201528091505092915050565b600073ffffffffffffffffffffffffffffffffffffffff808a1683528089166020840152808816604084015280871660608401525084608083015260c060a083015261137060c0830184866110a7565b9998505050505050505050565b600060408301825184526020808401516040828701528281518085526060880191508383019450600092505b808310156113c957845182529383019360019290920191908301906113a9565b509695505050505050565b600073ffffffffffffffffffffffffffffffffffffffff808816835280871660208401525060a0604083015261140d60a08301866111cb565b846060840152828103608084015283518152602084015160a06020830152805160a0830152602081015160c0830152604081015160e083015260608101516101008301526080810151905060a061012083015261146e6101408301826111cb565b905060408501518282036040840152611487828261137d565b915050606085015182820360608401526114a182826111cb565b915050608085015182820360808401526114bb82826111cb565b9a9950505050505050505050565b600082516114db8184602087016111a7565b9190910192915050565b60208152600061115260208301846111cb56fea164736f6c6343000813000a",
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

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterTransactor) FinalizeWithdrawERC20FromL2(opts *bind.TransactOpts, from common.Address, to common.Address, data []byte) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.contract.Transact(opts, "finalizeWithdrawERC20FromL2", from, to, data)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterSession) FinalizeWithdrawERC20FromL2(from common.Address, to common.Address, data []byte) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.Contract.FinalizeWithdrawERC20FromL2(&_OptimismL1BridgeAdapter.TransactOpts, from, to, data)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterTransactorSession) FinalizeWithdrawERC20FromL2(from common.Address, to common.Address, data []byte) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.Contract.FinalizeWithdrawERC20FromL2(&_OptimismL1BridgeAdapter.TransactOpts, from, to, data)
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

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterTransactor) SendERC20(opts *bind.TransactOpts, l1Token common.Address, l2Token common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.contract.Transact(opts, "sendERC20", l1Token, l2Token, recipient, amount)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterSession) SendERC20(l1Token common.Address, l2Token common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.Contract.SendERC20(&_OptimismL1BridgeAdapter.TransactOpts, l1Token, l2Token, recipient, amount)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterTransactorSession) SendERC20(l1Token common.Address, l2Token common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.Contract.SendERC20(&_OptimismL1BridgeAdapter.TransactOpts, l1Token, l2Token, recipient, amount)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapter) Address() common.Address {
	return _OptimismL1BridgeAdapter.address
}

type OptimismL1BridgeAdapterInterface interface {
	GetBridgeFeeInNative(opts *bind.CallOpts) (*big.Int, error)

	GetWrappedNative(opts *bind.CallOpts) (common.Address, error)

	DepositNativeToL2(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	FinalizeWithdrawERC20FromL2(opts *bind.TransactOpts, from common.Address, to common.Address, data []byte) (*types.Transaction, error)

	FinalizeWithdrawNativeFromL2(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	RelayMessageFromL2ToL1(opts *bind.TransactOpts, target common.Address, sender common.Address, message []byte, messageNonce *big.Int, proof IL1CrossDomainMessengerL2MessageInclusionProof) (*types.Transaction, error)

	SendERC20(opts *bind.TransactOpts, l1Token common.Address, l2Token common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	Address() common.Address
}
