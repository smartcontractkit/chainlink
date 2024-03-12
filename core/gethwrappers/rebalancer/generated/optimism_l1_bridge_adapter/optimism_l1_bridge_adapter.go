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
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIL1StandardBridge\",\"name\":\"l1Bridge\",\"type\":\"address\"},{\"internalType\":\"contractIWrappedNative\",\"name\":\"wrappedNative\",\"type\":\"address\"},{\"internalType\":\"contractIL1CrossDomainMessenger\",\"name\":\"l1CrossDomainMessenger\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BridgeAddressCannotBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"wanted\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"got\",\"type\":\"uint256\"}],\"name\":\"InsufficientEthValue\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"MsgShouldNotContainValue\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"msgValue\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"MsgValueDoesNotMatchAmount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"remoteSender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"localReceiver\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"finalizeWithdrawERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"finalizeWithdrawNativeFromL2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBridgeFeeInNative\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getWrappedNative\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"messageNonce\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"batchIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"batchRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"batchSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"prevTotalElements\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"internalType\":\"structLib_OVMCodec.ChainBatchHeader\",\"name\":\"stateRootBatchHeader\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"siblings\",\"type\":\"bytes32[]\"}],\"internalType\":\"structLib_OVMCodec.ChainInclusionProof\",\"name\":\"stateRootProof\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"stateTrieWitness\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"storageTrieWitness\",\"type\":\"bytes\"}],\"internalType\":\"structIL1CrossDomainMessenger.L2MessageInclusionProof\",\"name\":\"proof\",\"type\":\"tuple\"}],\"name\":\"relayMessageFromL2ToL1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"sendERC20\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x60e0604052600080546001600160401b03191690553480156200002157600080fd5b5060405162001776380380620017768339810160408190526200004491620000b7565b6001600160a01b03831615806200006257506001600160a01b038216155b156200008157604051635e9c404d60e11b815260040160405180910390fd5b6001600160a01b03928316608052821660a0521660c0526200010b565b6001600160a01b0381168114620000b457600080fd5b50565b600080600060608486031215620000cd57600080fd5b8351620000da816200009e565b6020850151909350620000ed816200009e565b604085015190925062000100816200009e565b809150509250925092565b60805160a05160c05161161162000165600039600081816101140152818161037f01526103ff0152600061065f01526000818161019b0152818161026b015281816104ab0152818161054001526105a201526116116000f3fe6080604052600436106100655760003560e01c8063a71d98b711610043578063a71d98b7146100cd578063e861e907146100ed578063f2bfa1e11461013e57600080fd5b806318b3050c1461006a5780632e4b1fc91461008c57806338314bb2146100ad575b600080fd5b34801561007657600080fd5b5061008a610085366004610c92565b61015e565b005b34801561009857600080fd5b50604051600081526020015b60405180910390f35b3480156100b957600080fd5b5061008a6100c8366004610d01565b610211565b6100e06100db366004610d62565b6102aa565b6040516100a49190610e4f565b3480156100f957600080fd5b5060405173ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001681526020016100a4565b34801561014a57600080fd5b5061008a6101593660046110ed565b610622565b6040517f1532ec3400000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001690631532ec34906101d8908890889088908890889060040161126b565b600060405180830381600087803b1580156101f257600080fd5b505af1158015610206573d6000803e3d6000fd5b505050505050505050565b600061021f828401846112ab565b8051602082015160408084015190517fa9f9e67500000000000000000000000000000000000000000000000000000000815293945073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169363a9f9e675936101d893909290918b918b918b908b90600401611313565b60606102ce73ffffffffffffffffffffffffffffffffffffffff881633308761069c565b341561030d576040517f2543d86e0000000000000000000000000000000000000000000000000000000081523460048201526024015b60405180910390fd5b6000805467ffffffffffffffff16818061032683611370565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555060405160200161036c919067ffffffffffffffff91909116815260200190565b60405160208183030381529060405290507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168873ffffffffffffffffffffffffffffffffffffffff1603610524576040517f2e1a7d4d000000000000000000000000000000000000000000000000000000008152600481018690527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690632e1a7d4d90602401600060405180830381600087803b15801561045857600080fd5b505af115801561046c573d6000803e3d6000fd5b50506040517f9a2ac6d500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169250639a2ac6d5915087906104e9908a9060009087906004016113be565b6000604051808303818588803b15801561050257600080fd5b505af1158015610516573d6000803e3d6000fd5b505050505080915050610618565b61056573ffffffffffffffffffffffffffffffffffffffff89167f00000000000000000000000000000000000000000000000000000000000000008761077e565b6040517f838b252000000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169063838b2520906105e2908b908b908b908b906000908990600401611402565b600060405180830381600087803b1580156105fc57600080fd5b505af1158015610610573d6000803e3d6000fd5b509293505050505b9695505050505050565b6040517fd7fd19dd00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169063d7fd19dd906101d890889088908890889088906004016114b8565b60405173ffffffffffffffffffffffffffffffffffffffff808516602483015283166044820152606481018290526107789085907f23b872dd00000000000000000000000000000000000000000000000000000000906084015b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152610905565b50505050565b80158061081e57506040517fdd62ed3e00000000000000000000000000000000000000000000000000000000815230600482015273ffffffffffffffffffffffffffffffffffffffff838116602483015284169063dd62ed3e90604401602060405180830381865afa1580156107f8573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061081c91906115ad565b155b6108aa576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603660248201527f5361666545524332303a20617070726f76652066726f6d206e6f6e2d7a65726f60448201527f20746f206e6f6e2d7a65726f20616c6c6f77616e6365000000000000000000006064820152608401610304565b60405173ffffffffffffffffffffffffffffffffffffffff83166024820152604481018290526109009084907f095ea7b300000000000000000000000000000000000000000000000000000000906064016106f6565b505050565b6000610967826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff16610a119092919063ffffffff16565b805190915015610900578080602001905181019061098591906115c6565b610900576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152608401610304565b6060610a208484600085610a28565b949350505050565b606082471015610aba576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c00000000000000000000000000000000000000000000000000006064820152608401610304565b6000808673ffffffffffffffffffffffffffffffffffffffff168587604051610ae391906115e8565b60006040518083038185875af1925050503d8060008114610b20576040519150601f19603f3d011682016040523d82523d6000602084013e610b25565b606091505b5091509150610b3687838387610b41565b979650505050505050565b60608315610bd7578251600003610bd05773ffffffffffffffffffffffffffffffffffffffff85163b610bd0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610304565b5081610a20565b610a208383815115610bec5781518083602001fd5b806040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103049190610e4f565b803573ffffffffffffffffffffffffffffffffffffffff81168114610c4457600080fd5b919050565b60008083601f840112610c5b57600080fd5b50813567ffffffffffffffff811115610c7357600080fd5b602083019150836020828501011115610c8b57600080fd5b9250929050565b600080600080600060808688031215610caa57600080fd5b610cb386610c20565b9450610cc160208701610c20565b935060408601359250606086013567ffffffffffffffff811115610ce457600080fd5b610cf088828901610c49565b969995985093965092949392505050565b60008060008060608587031215610d1757600080fd5b610d2085610c20565b9350610d2e60208601610c20565b9250604085013567ffffffffffffffff811115610d4a57600080fd5b610d5687828801610c49565b95989497509550505050565b60008060008060008060a08789031215610d7b57600080fd5b610d8487610c20565b9550610d9260208801610c20565b9450610da060408801610c20565b935060608701359250608087013567ffffffffffffffff811115610dc357600080fd5b610dcf89828a01610c49565b979a9699509497509295939492505050565b60005b83811015610dfc578181015183820152602001610de4565b50506000910152565b60008151808452610e1d816020860160208601610de1565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000610e626020830184610e05565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160a0810167ffffffffffffffff81118282101715610ebb57610ebb610e69565b60405290565b6040805190810167ffffffffffffffff81118282101715610ebb57610ebb610e69565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610f2b57610f2b610e69565b604052919050565b600082601f830112610f4457600080fd5b813567ffffffffffffffff811115610f5e57610f5e610e69565b610f8f60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601610ee4565b818152846020838601011115610fa457600080fd5b816020850160208301376000918101602001919091529392505050565b600060a08284031215610fd357600080fd5b610fdb610e98565b905081358152602082013560208201526040820135604082015260608201356060820152608082013567ffffffffffffffff81111561101957600080fd5b61102584828501610f33565b60808301525092915050565b60006040828403121561104357600080fd5b61104b610ec1565b90508135815260208083013567ffffffffffffffff8082111561106d57600080fd5b818501915085601f83011261108157600080fd5b81358181111561109357611093610e69565b8060051b91506110a4848301610ee4565b81815291830184019184810190888411156110be57600080fd5b938501935b838510156110dc578435825293850193908501906110c3565b808688015250505050505092915050565b600080600080600060a0868803121561110557600080fd5b61110e86610c20565b945061111c60208701610c20565b9350604086013567ffffffffffffffff8082111561113957600080fd5b61114589838a01610f33565b945060608801359350608088013591508082111561116257600080fd5b9087019060a0828a03121561117657600080fd5b61117e610e98565b8235815260208301358281111561119457600080fd5b6111a08b828601610fc1565b6020830152506040830135828111156111b857600080fd5b6111c48b828601611031565b6040830152506060830135828111156111dc57600080fd5b6111e88b828601610f33565b60608301525060808301358281111561120057600080fd5b61120c8b828601610f33565b6080830152508093505050509295509295909350565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b600073ffffffffffffffffffffffffffffffffffffffff808816835280871660208401525084604083015260806060830152610b36608083018486611222565b6000606082840312156112bd57600080fd5b6040516060810181811067ffffffffffffffff821117156112e0576112e0610e69565b6040526112ec83610c20565b81526112fa60208401610c20565b6020820152604083013560408201528091505092915050565b600073ffffffffffffffffffffffffffffffffffffffff808a1683528089166020840152808816604084015280871660608401525084608083015260c060a083015261136360c083018486611222565b9998505050505050505050565b600067ffffffffffffffff8083168181036113b4577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6001019392505050565b73ffffffffffffffffffffffffffffffffffffffff8416815263ffffffff831660208201526060604082015260006113f96060830184610e05565b95945050505050565b600073ffffffffffffffffffffffffffffffffffffffff8089168352808816602084015280871660408401525084606083015263ffffffff8416608083015260c060a083015261145560c0830184610e05565b98975050505050505050565b600060408301825184526020808401516040828701528281518085526060880191508383019450600092505b808310156114ad578451825293830193600192909201919083019061148d565b509695505050505050565b600073ffffffffffffffffffffffffffffffffffffffff808816835280871660208401525060a060408301526114f160a0830186610e05565b846060840152828103608084015283518152602084015160a06020830152805160a0830152602081015160c0830152604081015160e083015260608101516101008301526080810151905060a0610120830152611552610140830182610e05565b90506040850151828203604084015261156b8282611461565b915050606085015182820360608401526115858282610e05565b9150506080850151828203608084015261159f8282610e05565b9a9950505050505050505050565b6000602082840312156115bf57600080fd5b5051919050565b6000602082840312156115d857600080fd5b81518015158114610e6257600080fd5b600082516115fa818460208701610de1565b919091019291505056fea164736f6c6343000813000a",
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

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterTransactor) SendERC20(opts *bind.TransactOpts, localToken common.Address, remoteToken common.Address, recipient common.Address, amount *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.contract.Transact(opts, "sendERC20", localToken, remoteToken, recipient, amount, arg4)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterSession) SendERC20(localToken common.Address, remoteToken common.Address, recipient common.Address, amount *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.Contract.SendERC20(&_OptimismL1BridgeAdapter.TransactOpts, localToken, remoteToken, recipient, amount, arg4)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapterTransactorSession) SendERC20(localToken common.Address, remoteToken common.Address, recipient common.Address, amount *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _OptimismL1BridgeAdapter.Contract.SendERC20(&_OptimismL1BridgeAdapter.TransactOpts, localToken, remoteToken, recipient, amount, arg4)
}

func (_OptimismL1BridgeAdapter *OptimismL1BridgeAdapter) Address() common.Address {
	return _OptimismL1BridgeAdapter.address
}

type OptimismL1BridgeAdapterInterface interface {
	GetBridgeFeeInNative(opts *bind.CallOpts) (*big.Int, error)

	GetWrappedNative(opts *bind.CallOpts) (common.Address, error)

	FinalizeWithdrawERC20(opts *bind.TransactOpts, remoteSender common.Address, localReceiver common.Address, data []byte) (*types.Transaction, error)

	FinalizeWithdrawNativeFromL2(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	RelayMessageFromL2ToL1(opts *bind.TransactOpts, target common.Address, sender common.Address, message []byte, messageNonce *big.Int, proof IL1CrossDomainMessengerL2MessageInclusionProof) (*types.Transaction, error)

	SendERC20(opts *bind.TransactOpts, localToken common.Address, remoteToken common.Address, recipient common.Address, amount *big.Int, arg4 []byte) (*types.Transaction, error)

	Address() common.Address
}
