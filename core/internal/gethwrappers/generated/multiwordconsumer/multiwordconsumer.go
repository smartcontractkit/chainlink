// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package multiwordconsumer

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// MultiwordConsumerABI is the input ABI used to generate the binding from.
const MultiwordConsumerABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_specId\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ChainlinkCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ChainlinkFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ChainlinkRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes\",\"name\":\"price\",\"type\":\"bytes\"}],\"name\":\"RequestFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"first\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"second\",\"type\":\"bytes32\"}],\"name\":\"RequestMultipleFulfilled\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_requestId\",\"type\":\"bytes32\"}],\"name\":\"addExternalRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_payment\",\"type\":\"uint256\"},{\"internalType\":\"bytes4\",\"name\":\"_callbackFunctionId\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"_expiration\",\"type\":\"uint256\"}],\"name\":\"cancelRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentPrice\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"first\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"_price\",\"type\":\"bytes\"}],\"name\":\"fulfillBytes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_first\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_second\",\"type\":\"bytes32\"}],\"name\":\"fulfillMultipleParameters\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_currency\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_payment\",\"type\":\"uint256\"}],\"name\":\"requestEthereumPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_currency\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_payment\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_callback\",\"type\":\"address\"}],\"name\":\"requestEthereumPriceByCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_currency\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_payment\",\"type\":\"uint256\"}],\"name\":\"requestMultipleParameters\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"second\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// MultiwordConsumerBin is the compiled bytecode used for deploying new contracts.
var MultiwordConsumerBin = "0x6080604052600160045534801561001557600080fd5b50604051611a37380380611a378339818101604052606081101561003857600080fd5b508051602082015160409092015190919061005b836001600160e01b0361007816565b61006d826001600160e01b0361009a16565b600655506100bc9050565b600280546001600160a01b0319166001600160a01b0392909216919091179055565b600380546001600160a01b0319166001600160a01b0392909216919091179055565b61196c806100cb6000396000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c806383db5cbc11610081578063c2fb85231161005b578063c2fb852314610376578063e89855ba14610423578063e8d5359d146104cb576100c9565b806383db5cbc146102495780638dc654a2146102f15780639d1b464a146102f9576100c9565b80635591a608116100b25780635591a608146101135780635a8ac02d1461018057806374961d4d14610188576100c9565b80633df4ddf4146100ce57806353389072146100e8575b600080fd5b6100d6610504565b60408051918252519081900360200190f35b610111600480360360608110156100fe57600080fd5b508035906020810135906040013561050a565b005b610111600480360360a081101561012957600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813516906020810135906040810135907fffffffff00000000000000000000000000000000000000000000000000000000606082013516906080013561061f565b6100d66106e6565b6101116004803603606081101561019e57600080fd5b8101906020810181356401000000008111156101b957600080fd5b8201836020820111156101cb57600080fd5b803590602001918460018302840111640100000000831117156101ed57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550508235935050506020013573ffffffffffffffffffffffffffffffffffffffff166106ec565b6101116004803603604081101561025f57600080fd5b81019060208101813564010000000081111561027a57600080fd5b82018360208201111561028c57600080fd5b803590602001918460018302840111640100000000831117156102ae57600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295505091359250610827915050565b610111610836565b6103016109f3565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561033b578181015183820152602001610323565b50505050905090810190601f1680156103685780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6101116004803603604081101561038c57600080fd5b813591908101906040810160208201356401000000008111156103ae57600080fd5b8201836020820111156103c057600080fd5b803590602001918460018302840111640100000000831117156103e257600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610a9f945050505050565b6101116004803603604081101561043957600080fd5b81019060208101813564010000000081111561045457600080fd5b82018360208201111561046657600080fd5b8035906020019184600183028401116401000000008311171561048857600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295505091359250610c51915050565b610111600480360360408110156104e157600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135610d8b565b60085481565b600083815260056020526040902054839073ffffffffffffffffffffffffffffffffffffffff163314610588576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260288152602001806118c86028913960400191505060405180910390fd5b60008181526005602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001690555182917f7cc135e0cebb02c3480ae5d74d377283180a2601f8f644edf7987b009316c63a91a28183857fd368a628c6f427add4c36c69828a9be4d937a803adfda79c1dbf7eb26cdf4bc460405160405180910390a45060089190915560095550565b604080517f6ee4d55300000000000000000000000000000000000000000000000000000000815260048101869052602481018590527fffffffff0000000000000000000000000000000000000000000000000000000084166044820152606481018390529051869173ffffffffffffffffffffffffffffffffffffffff831691636ee4d5539160848082019260009290919082900301818387803b1580156106c657600080fd5b505af11580156106da573d6000803e3d6000fd5b50505050505050505050565b60095481565b6106f46117bd565b60065461072290837fc2fb852300000000000000000000000000000000000000000000000000000000610d95565b90506107846040518060400160405280600381526020017f67657400000000000000000000000000000000000000000000000000000000008152506040518060800160405280604781526020016118f06047913983919063ffffffff610dc016565b604080516001808252818301909252606091816020015b606081526020019060019003908161079b57905050905084816000815181106107c057fe5b60200260200101819052506108156040518060400160405280600481526020017f70617468000000000000000000000000000000000000000000000000000000008152508284610def9092919063ffffffff16565b61081f8285610e5d565b505050505050565b6108328282306106ec565b5050565b6000610840610e8d565b604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905191925073ffffffffffffffffffffffffffffffffffffffff83169163a9059cbb91339184916370a08231916024808301926020929190829003018186803b1580156108b957600080fd5b505afa1580156108cd573d6000803e3d6000fd5b505050506040513d60208110156108e357600080fd5b5051604080517fffffffff0000000000000000000000000000000000000000000000000000000060e086901b16815273ffffffffffffffffffffffffffffffffffffffff909316600484015260248301919091525160448083019260209291908290030181600087803b15801561095957600080fd5b505af115801561096d573d6000803e3d6000fd5b505050506040513d602081101561098357600080fd5b50516109f057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f556e61626c6520746f207472616e736665720000000000000000000000000000604482015290519081900360640190fd5b50565b6007805460408051602060026001851615610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190941693909304601f81018490048402820184019092528181529291830182828015610a975780601f10610a6c57610100808354040283529160200191610a97565b820191906000526020600020905b815481529060010190602001808311610a7a57829003601f168201915b505050505081565b600082815260056020526040902054829073ffffffffffffffffffffffffffffffffffffffff163314610b1d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260288152602001806118c86028913960400191505060405180910390fd5b60008181526005602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001690555182917f7cc135e0cebb02c3480ae5d74d377283180a2601f8f644edf7987b009316c63a91a2816040518082805190602001908083835b60208310610bc657805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610b89565b5181516020939093036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990911692169190911790526040519201829003822093508692507f1a111c5dcf9a71088bd5e1797fdfaf399fec2afbb24aca247e4e3e9f4b61df919160009150a38151610c4b9060079060208501906117f2565b50505050565b610c596117bd565b600654610c8790307f5338907200000000000000000000000000000000000000000000000000000000610d95565b9050610ce96040518060400160405280600381526020017f67657400000000000000000000000000000000000000000000000000000000008152506040518060800160405280604781526020016118f06047913983919063ffffffff610dc016565b604080516001808252818301909252606091816020015b6060815260200190600190039081610d005790505090508381600081518110610d2557fe5b6020026020010181905250610d7a6040518060400160405280600481526020017f70617468000000000000000000000000000000000000000000000000000000008152508284610def9092919063ffffffff16565b610d848284610e5d565b5050505050565b6108328282610eaa565b610d9d6117bd565b610da56117bd565b610db78186868663ffffffff610f9116565b95945050505050565b6080830151610dd5908363ffffffff610ff316565b6080830151610dea908263ffffffff610ff316565b505050565b6080830151610e04908363ffffffff610ff316565b610e118360800151611010565b60005b8151811015610e4f57610e47828281518110610e2c57fe5b60200260200101518560800151610ff390919063ffffffff16565b600101610e14565b50610dea836080015161101b565b600354600090610e849073ffffffffffffffffffffffffffffffffffffffff168484611026565b90505b92915050565b60025473ffffffffffffffffffffffffffffffffffffffff165b90565b600081815260056020526040902054819073ffffffffffffffffffffffffffffffffffffffff1615610f3d57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f5265717565737420697320616c72656164792070656e64696e67000000000000604482015290519081900360640190fd5b50600090815260056020526040902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b610f996117bd565b610fa98560800151610100611263565b505091835273ffffffffffffffffffffffffffffffffffffffff1660208301527fffffffff0000000000000000000000000000000000000000000000000000000016604082015290565b611000826003835161129d565b610dea828263ffffffff6113a716565b6109f08160046113c1565b6109f08160076113c1565b6004546040805130606090811b60208084019190915260348084018690528451808503909101815260549093018452825192810192909220908601939093526000838152600590915281812080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8816179055905182917fb5e6e01e79f91267dc17b4e6314d5d4d03593d2ceee0fbb452b750bd70ea5af991a260025473ffffffffffffffffffffffffffffffffffffffff16634000aea08584611100876113dc565b6040518463ffffffff1660e01b8152600401808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b8381101561118457818101518382015260200161116c565b50505050905090810190601f1680156111b15780820380516001836020036101000a031916815260200191505b50945050505050602060405180830381600087803b1580156111d257600080fd5b505af11580156111e6573d6000803e3d6000fd5b505050506040513d60208110156111fc57600080fd5b5051611253576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260238152602001806118a56023913960400191505060405180910390fd5b6004805460010190559392505050565b61126b611870565b60208206156112805760208206602003820191505b506020828101829052604080518085526000815290920101905290565b601781116112c4576112be8360e0600585901b16831763ffffffff6115c516565b50610dea565b60ff81116112fa576112e7836018611fe0600586901b161763ffffffff6115c516565b506112be8382600163ffffffff6115dd16565b61ffff81116113315761131e836019611fe0600586901b161763ffffffff6115c516565b506112be8382600263ffffffff6115dd16565b63ffffffff811161136a5761135783601a611fe0600586901b161763ffffffff6115c516565b506112be8382600463ffffffff6115dd16565b67ffffffffffffffff8111610dea5761139483601b611fe0600586901b161763ffffffff6115c516565b50610c4b8382600863ffffffff6115dd16565b6113af611870565b610e84838460000151518485516115fe565b610dea82601f611fe0600585901b161763ffffffff6115c516565b6060634042994660e01b60008084600001518560200151866040015187606001516001896080015160000151604051602401808973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018881526020018781526020018673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001857bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200184815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b838110156115085781810151838201526020016114f0565b50505050905090810190601f1680156115355780820380516001836020036101000a031916815260200191505b50604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909d169c909c17909b5250989950505050505050505050919050565b6115cd611870565b610e8483846000015151846116e6565b6115e5611870565b6115f6848560000151518585611731565b949350505050565b611606611870565b825182111561161457600080fd5b8460200151828501111561163e5761163e85611636876020015187860161178f565b6002026117a6565b60008086518051876020830101935080888701111561165d5787860182525b505050602084015b602084106116a257805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09093019260209182019101611665565b5181517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60208690036101000a019081169019919091161790525083949350505050565b6116ee611870565b8360200151831061170a5761170a8485602001516002026117a6565b835180516020858301018481535080851415611727576001810182525b5093949350505050565b611739611870565b8460200151848301111561175657611756858584016002026117a6565b60006001836101000a0390508551838682010185831982511617815250805184870111156117845783860181525b509495945050505050565b6000818311156117a0575081610e87565b50919050565b81516117b28383611263565b50610c4b83826113a7565b6040805160a0810182526000808252602082018190529181018290526060810191909152608081016117ed611870565b905290565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061183357805160ff1916838001178555611860565b82800160010185558215611860579182015b82811115611860578251825591602001919060010190611845565b5061186c92915061188a565b5090565b604051806040016040528060608152602001600081525090565b610ea791905b8082111561186c576000815560010161189056fe756e61626c6520746f207472616e73666572416e6443616c6c20746f206f7261636c65536f75726365206d75737420626520746865206f7261636c65206f6620746865207265717565737468747470733a2f2f6d696e2d6170692e63727970746f636f6d706172652e636f6d2f646174612f70726963653f6673796d3d455448267473796d733d5553442c4555522c4a5059a264697066735822beefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef64736f6c6343decafe0033"

// DeployMultiwordConsumer deploys a new Ethereum contract, binding an instance of MultiwordConsumer to it.
func DeployMultiwordConsumer(auth *bind.TransactOpts, backend bind.ContractBackend, _link common.Address, _oracle common.Address, _specId [32]byte) (common.Address, *types.Transaction, *MultiwordConsumer, error) {
	parsed, err := abi.JSON(strings.NewReader(MultiwordConsumerABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(MultiwordConsumerBin), backend, _link, _oracle, _specId)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MultiwordConsumer{MultiwordConsumerCaller: MultiwordConsumerCaller{contract: contract}, MultiwordConsumerTransactor: MultiwordConsumerTransactor{contract: contract}, MultiwordConsumerFilterer: MultiwordConsumerFilterer{contract: contract}}, nil
}

// MultiwordConsumer is an auto generated Go binding around an Ethereum contract.
type MultiwordConsumer struct {
	MultiwordConsumerCaller     // Read-only binding to the contract
	MultiwordConsumerTransactor // Write-only binding to the contract
	MultiwordConsumerFilterer   // Log filterer for contract events
}

// MultiwordConsumerCaller is an auto generated read-only Go binding around an Ethereum contract.
type MultiwordConsumerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiwordConsumerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MultiwordConsumerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiwordConsumerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MultiwordConsumerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiwordConsumerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MultiwordConsumerSession struct {
	Contract     *MultiwordConsumer // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// MultiwordConsumerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MultiwordConsumerCallerSession struct {
	Contract *MultiwordConsumerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// MultiwordConsumerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MultiwordConsumerTransactorSession struct {
	Contract     *MultiwordConsumerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// MultiwordConsumerRaw is an auto generated low-level Go binding around an Ethereum contract.
type MultiwordConsumerRaw struct {
	Contract *MultiwordConsumer // Generic contract binding to access the raw methods on
}

// MultiwordConsumerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MultiwordConsumerCallerRaw struct {
	Contract *MultiwordConsumerCaller // Generic read-only contract binding to access the raw methods on
}

// MultiwordConsumerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MultiwordConsumerTransactorRaw struct {
	Contract *MultiwordConsumerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMultiwordConsumer creates a new instance of MultiwordConsumer, bound to a specific deployed contract.
func NewMultiwordConsumer(address common.Address, backend bind.ContractBackend) (*MultiwordConsumer, error) {
	contract, err := bindMultiwordConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MultiwordConsumer{MultiwordConsumerCaller: MultiwordConsumerCaller{contract: contract}, MultiwordConsumerTransactor: MultiwordConsumerTransactor{contract: contract}, MultiwordConsumerFilterer: MultiwordConsumerFilterer{contract: contract}}, nil
}

// NewMultiwordConsumerCaller creates a new read-only instance of MultiwordConsumer, bound to a specific deployed contract.
func NewMultiwordConsumerCaller(address common.Address, caller bind.ContractCaller) (*MultiwordConsumerCaller, error) {
	contract, err := bindMultiwordConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MultiwordConsumerCaller{contract: contract}, nil
}

// NewMultiwordConsumerTransactor creates a new write-only instance of MultiwordConsumer, bound to a specific deployed contract.
func NewMultiwordConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*MultiwordConsumerTransactor, error) {
	contract, err := bindMultiwordConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MultiwordConsumerTransactor{contract: contract}, nil
}

// NewMultiwordConsumerFilterer creates a new log filterer instance of MultiwordConsumer, bound to a specific deployed contract.
func NewMultiwordConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*MultiwordConsumerFilterer, error) {
	contract, err := bindMultiwordConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MultiwordConsumerFilterer{contract: contract}, nil
}

// bindMultiwordConsumer binds a generic wrapper to an already deployed contract.
func bindMultiwordConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MultiwordConsumerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MultiwordConsumer *MultiwordConsumerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _MultiwordConsumer.Contract.MultiwordConsumerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MultiwordConsumer *MultiwordConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiwordConsumer.Contract.MultiwordConsumerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MultiwordConsumer *MultiwordConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultiwordConsumer.Contract.MultiwordConsumerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MultiwordConsumer *MultiwordConsumerCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _MultiwordConsumer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MultiwordConsumer *MultiwordConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiwordConsumer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MultiwordConsumer *MultiwordConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultiwordConsumer.Contract.contract.Transact(opts, method, params...)
}

// CurrentPrice is a free data retrieval call binding the contract method 0x9d1b464a.
//
// Solidity: function currentPrice() view returns(bytes)
func (_MultiwordConsumer *MultiwordConsumerCaller) CurrentPrice(opts *bind.CallOpts) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _MultiwordConsumer.contract.Call(opts, out, "currentPrice")
	return *ret0, err
}

// CurrentPrice is a free data retrieval call binding the contract method 0x9d1b464a.
//
// Solidity: function currentPrice() view returns(bytes)
func (_MultiwordConsumer *MultiwordConsumerSession) CurrentPrice() ([]byte, error) {
	return _MultiwordConsumer.Contract.CurrentPrice(&_MultiwordConsumer.CallOpts)
}

// CurrentPrice is a free data retrieval call binding the contract method 0x9d1b464a.
//
// Solidity: function currentPrice() view returns(bytes)
func (_MultiwordConsumer *MultiwordConsumerCallerSession) CurrentPrice() ([]byte, error) {
	return _MultiwordConsumer.Contract.CurrentPrice(&_MultiwordConsumer.CallOpts)
}

// First is a free data retrieval call binding the contract method 0x3df4ddf4.
//
// Solidity: function first() view returns(bytes32)
func (_MultiwordConsumer *MultiwordConsumerCaller) First(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _MultiwordConsumer.contract.Call(opts, out, "first")
	return *ret0, err
}

// First is a free data retrieval call binding the contract method 0x3df4ddf4.
//
// Solidity: function first() view returns(bytes32)
func (_MultiwordConsumer *MultiwordConsumerSession) First() ([32]byte, error) {
	return _MultiwordConsumer.Contract.First(&_MultiwordConsumer.CallOpts)
}

// First is a free data retrieval call binding the contract method 0x3df4ddf4.
//
// Solidity: function first() view returns(bytes32)
func (_MultiwordConsumer *MultiwordConsumerCallerSession) First() ([32]byte, error) {
	return _MultiwordConsumer.Contract.First(&_MultiwordConsumer.CallOpts)
}

// Second is a free data retrieval call binding the contract method 0x5a8ac02d.
//
// Solidity: function second() view returns(bytes32)
func (_MultiwordConsumer *MultiwordConsumerCaller) Second(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _MultiwordConsumer.contract.Call(opts, out, "second")
	return *ret0, err
}

// Second is a free data retrieval call binding the contract method 0x5a8ac02d.
//
// Solidity: function second() view returns(bytes32)
func (_MultiwordConsumer *MultiwordConsumerSession) Second() ([32]byte, error) {
	return _MultiwordConsumer.Contract.Second(&_MultiwordConsumer.CallOpts)
}

// Second is a free data retrieval call binding the contract method 0x5a8ac02d.
//
// Solidity: function second() view returns(bytes32)
func (_MultiwordConsumer *MultiwordConsumerCallerSession) Second() ([32]byte, error) {
	return _MultiwordConsumer.Contract.Second(&_MultiwordConsumer.CallOpts)
}

// AddExternalRequest is a paid mutator transaction binding the contract method 0xe8d5359d.
//
// Solidity: function addExternalRequest(address _oracle, bytes32 _requestId) returns()
func (_MultiwordConsumer *MultiwordConsumerTransactor) AddExternalRequest(opts *bind.TransactOpts, _oracle common.Address, _requestId [32]byte) (*types.Transaction, error) {
	return _MultiwordConsumer.contract.Transact(opts, "addExternalRequest", _oracle, _requestId)
}

// AddExternalRequest is a paid mutator transaction binding the contract method 0xe8d5359d.
//
// Solidity: function addExternalRequest(address _oracle, bytes32 _requestId) returns()
func (_MultiwordConsumer *MultiwordConsumerSession) AddExternalRequest(_oracle common.Address, _requestId [32]byte) (*types.Transaction, error) {
	return _MultiwordConsumer.Contract.AddExternalRequest(&_MultiwordConsumer.TransactOpts, _oracle, _requestId)
}

// AddExternalRequest is a paid mutator transaction binding the contract method 0xe8d5359d.
//
// Solidity: function addExternalRequest(address _oracle, bytes32 _requestId) returns()
func (_MultiwordConsumer *MultiwordConsumerTransactorSession) AddExternalRequest(_oracle common.Address, _requestId [32]byte) (*types.Transaction, error) {
	return _MultiwordConsumer.Contract.AddExternalRequest(&_MultiwordConsumer.TransactOpts, _oracle, _requestId)
}

// CancelRequest is a paid mutator transaction binding the contract method 0x5591a608.
//
// Solidity: function cancelRequest(address _oracle, bytes32 _requestId, uint256 _payment, bytes4 _callbackFunctionId, uint256 _expiration) returns()
func (_MultiwordConsumer *MultiwordConsumerTransactor) CancelRequest(opts *bind.TransactOpts, _oracle common.Address, _requestId [32]byte, _payment *big.Int, _callbackFunctionId [4]byte, _expiration *big.Int) (*types.Transaction, error) {
	return _MultiwordConsumer.contract.Transact(opts, "cancelRequest", _oracle, _requestId, _payment, _callbackFunctionId, _expiration)
}

// CancelRequest is a paid mutator transaction binding the contract method 0x5591a608.
//
// Solidity: function cancelRequest(address _oracle, bytes32 _requestId, uint256 _payment, bytes4 _callbackFunctionId, uint256 _expiration) returns()
func (_MultiwordConsumer *MultiwordConsumerSession) CancelRequest(_oracle common.Address, _requestId [32]byte, _payment *big.Int, _callbackFunctionId [4]byte, _expiration *big.Int) (*types.Transaction, error) {
	return _MultiwordConsumer.Contract.CancelRequest(&_MultiwordConsumer.TransactOpts, _oracle, _requestId, _payment, _callbackFunctionId, _expiration)
}

// CancelRequest is a paid mutator transaction binding the contract method 0x5591a608.
//
// Solidity: function cancelRequest(address _oracle, bytes32 _requestId, uint256 _payment, bytes4 _callbackFunctionId, uint256 _expiration) returns()
func (_MultiwordConsumer *MultiwordConsumerTransactorSession) CancelRequest(_oracle common.Address, _requestId [32]byte, _payment *big.Int, _callbackFunctionId [4]byte, _expiration *big.Int) (*types.Transaction, error) {
	return _MultiwordConsumer.Contract.CancelRequest(&_MultiwordConsumer.TransactOpts, _oracle, _requestId, _payment, _callbackFunctionId, _expiration)
}

// FulfillBytes is a paid mutator transaction binding the contract method 0xc2fb8523.
//
// Solidity: function fulfillBytes(bytes32 _requestId, bytes _price) returns()
func (_MultiwordConsumer *MultiwordConsumerTransactor) FulfillBytes(opts *bind.TransactOpts, _requestId [32]byte, _price []byte) (*types.Transaction, error) {
	return _MultiwordConsumer.contract.Transact(opts, "fulfillBytes", _requestId, _price)
}

// FulfillBytes is a paid mutator transaction binding the contract method 0xc2fb8523.
//
// Solidity: function fulfillBytes(bytes32 _requestId, bytes _price) returns()
func (_MultiwordConsumer *MultiwordConsumerSession) FulfillBytes(_requestId [32]byte, _price []byte) (*types.Transaction, error) {
	return _MultiwordConsumer.Contract.FulfillBytes(&_MultiwordConsumer.TransactOpts, _requestId, _price)
}

// FulfillBytes is a paid mutator transaction binding the contract method 0xc2fb8523.
//
// Solidity: function fulfillBytes(bytes32 _requestId, bytes _price) returns()
func (_MultiwordConsumer *MultiwordConsumerTransactorSession) FulfillBytes(_requestId [32]byte, _price []byte) (*types.Transaction, error) {
	return _MultiwordConsumer.Contract.FulfillBytes(&_MultiwordConsumer.TransactOpts, _requestId, _price)
}

// FulfillMultipleParameters is a paid mutator transaction binding the contract method 0x53389072.
//
// Solidity: function fulfillMultipleParameters(bytes32 _requestId, bytes32 _first, bytes32 _second) returns()
func (_MultiwordConsumer *MultiwordConsumerTransactor) FulfillMultipleParameters(opts *bind.TransactOpts, _requestId [32]byte, _first [32]byte, _second [32]byte) (*types.Transaction, error) {
	return _MultiwordConsumer.contract.Transact(opts, "fulfillMultipleParameters", _requestId, _first, _second)
}

// FulfillMultipleParameters is a paid mutator transaction binding the contract method 0x53389072.
//
// Solidity: function fulfillMultipleParameters(bytes32 _requestId, bytes32 _first, bytes32 _second) returns()
func (_MultiwordConsumer *MultiwordConsumerSession) FulfillMultipleParameters(_requestId [32]byte, _first [32]byte, _second [32]byte) (*types.Transaction, error) {
	return _MultiwordConsumer.Contract.FulfillMultipleParameters(&_MultiwordConsumer.TransactOpts, _requestId, _first, _second)
}

// FulfillMultipleParameters is a paid mutator transaction binding the contract method 0x53389072.
//
// Solidity: function fulfillMultipleParameters(bytes32 _requestId, bytes32 _first, bytes32 _second) returns()
func (_MultiwordConsumer *MultiwordConsumerTransactorSession) FulfillMultipleParameters(_requestId [32]byte, _first [32]byte, _second [32]byte) (*types.Transaction, error) {
	return _MultiwordConsumer.Contract.FulfillMultipleParameters(&_MultiwordConsumer.TransactOpts, _requestId, _first, _second)
}

// RequestEthereumPrice is a paid mutator transaction binding the contract method 0x83db5cbc.
//
// Solidity: function requestEthereumPrice(string _currency, uint256 _payment) returns()
func (_MultiwordConsumer *MultiwordConsumerTransactor) RequestEthereumPrice(opts *bind.TransactOpts, _currency string, _payment *big.Int) (*types.Transaction, error) {
	return _MultiwordConsumer.contract.Transact(opts, "requestEthereumPrice", _currency, _payment)
}

// RequestEthereumPrice is a paid mutator transaction binding the contract method 0x83db5cbc.
//
// Solidity: function requestEthereumPrice(string _currency, uint256 _payment) returns()
func (_MultiwordConsumer *MultiwordConsumerSession) RequestEthereumPrice(_currency string, _payment *big.Int) (*types.Transaction, error) {
	return _MultiwordConsumer.Contract.RequestEthereumPrice(&_MultiwordConsumer.TransactOpts, _currency, _payment)
}

// RequestEthereumPrice is a paid mutator transaction binding the contract method 0x83db5cbc.
//
// Solidity: function requestEthereumPrice(string _currency, uint256 _payment) returns()
func (_MultiwordConsumer *MultiwordConsumerTransactorSession) RequestEthereumPrice(_currency string, _payment *big.Int) (*types.Transaction, error) {
	return _MultiwordConsumer.Contract.RequestEthereumPrice(&_MultiwordConsumer.TransactOpts, _currency, _payment)
}

// RequestEthereumPriceByCallback is a paid mutator transaction binding the contract method 0x74961d4d.
//
// Solidity: function requestEthereumPriceByCallback(string _currency, uint256 _payment, address _callback) returns()
func (_MultiwordConsumer *MultiwordConsumerTransactor) RequestEthereumPriceByCallback(opts *bind.TransactOpts, _currency string, _payment *big.Int, _callback common.Address) (*types.Transaction, error) {
	return _MultiwordConsumer.contract.Transact(opts, "requestEthereumPriceByCallback", _currency, _payment, _callback)
}

// RequestEthereumPriceByCallback is a paid mutator transaction binding the contract method 0x74961d4d.
//
// Solidity: function requestEthereumPriceByCallback(string _currency, uint256 _payment, address _callback) returns()
func (_MultiwordConsumer *MultiwordConsumerSession) RequestEthereumPriceByCallback(_currency string, _payment *big.Int, _callback common.Address) (*types.Transaction, error) {
	return _MultiwordConsumer.Contract.RequestEthereumPriceByCallback(&_MultiwordConsumer.TransactOpts, _currency, _payment, _callback)
}

// RequestEthereumPriceByCallback is a paid mutator transaction binding the contract method 0x74961d4d.
//
// Solidity: function requestEthereumPriceByCallback(string _currency, uint256 _payment, address _callback) returns()
func (_MultiwordConsumer *MultiwordConsumerTransactorSession) RequestEthereumPriceByCallback(_currency string, _payment *big.Int, _callback common.Address) (*types.Transaction, error) {
	return _MultiwordConsumer.Contract.RequestEthereumPriceByCallback(&_MultiwordConsumer.TransactOpts, _currency, _payment, _callback)
}

// RequestMultipleParameters is a paid mutator transaction binding the contract method 0xe89855ba.
//
// Solidity: function requestMultipleParameters(string _currency, uint256 _payment) returns()
func (_MultiwordConsumer *MultiwordConsumerTransactor) RequestMultipleParameters(opts *bind.TransactOpts, _currency string, _payment *big.Int) (*types.Transaction, error) {
	return _MultiwordConsumer.contract.Transact(opts, "requestMultipleParameters", _currency, _payment)
}

// RequestMultipleParameters is a paid mutator transaction binding the contract method 0xe89855ba.
//
// Solidity: function requestMultipleParameters(string _currency, uint256 _payment) returns()
func (_MultiwordConsumer *MultiwordConsumerSession) RequestMultipleParameters(_currency string, _payment *big.Int) (*types.Transaction, error) {
	return _MultiwordConsumer.Contract.RequestMultipleParameters(&_MultiwordConsumer.TransactOpts, _currency, _payment)
}

// RequestMultipleParameters is a paid mutator transaction binding the contract method 0xe89855ba.
//
// Solidity: function requestMultipleParameters(string _currency, uint256 _payment) returns()
func (_MultiwordConsumer *MultiwordConsumerTransactorSession) RequestMultipleParameters(_currency string, _payment *big.Int) (*types.Transaction, error) {
	return _MultiwordConsumer.Contract.RequestMultipleParameters(&_MultiwordConsumer.TransactOpts, _currency, _payment)
}

// WithdrawLink is a paid mutator transaction binding the contract method 0x8dc654a2.
//
// Solidity: function withdrawLink() returns()
func (_MultiwordConsumer *MultiwordConsumerTransactor) WithdrawLink(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiwordConsumer.contract.Transact(opts, "withdrawLink")
}

// WithdrawLink is a paid mutator transaction binding the contract method 0x8dc654a2.
//
// Solidity: function withdrawLink() returns()
func (_MultiwordConsumer *MultiwordConsumerSession) WithdrawLink() (*types.Transaction, error) {
	return _MultiwordConsumer.Contract.WithdrawLink(&_MultiwordConsumer.TransactOpts)
}

// WithdrawLink is a paid mutator transaction binding the contract method 0x8dc654a2.
//
// Solidity: function withdrawLink() returns()
func (_MultiwordConsumer *MultiwordConsumerTransactorSession) WithdrawLink() (*types.Transaction, error) {
	return _MultiwordConsumer.Contract.WithdrawLink(&_MultiwordConsumer.TransactOpts)
}

// MultiwordConsumerChainlinkCancelledIterator is returned from FilterChainlinkCancelled and is used to iterate over the raw logs and unpacked data for ChainlinkCancelled events raised by the MultiwordConsumer contract.
type MultiwordConsumerChainlinkCancelledIterator struct {
	Event *MultiwordConsumerChainlinkCancelled // Event containing the contract specifics and raw log

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
func (it *MultiwordConsumerChainlinkCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiwordConsumerChainlinkCancelled)
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
		it.Event = new(MultiwordConsumerChainlinkCancelled)
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
func (it *MultiwordConsumerChainlinkCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiwordConsumerChainlinkCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiwordConsumerChainlinkCancelled represents a ChainlinkCancelled event raised by the MultiwordConsumer contract.
type MultiwordConsumerChainlinkCancelled struct {
	Id  [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterChainlinkCancelled is a free log retrieval operation binding the contract event 0xe1fe3afa0f7f761ff0a8b89086790efd5140d2907ebd5b7ff6bfcb5e075fd4c5.
//
// Solidity: event ChainlinkCancelled(bytes32 indexed id)
func (_MultiwordConsumer *MultiwordConsumerFilterer) FilterChainlinkCancelled(opts *bind.FilterOpts, id [][32]byte) (*MultiwordConsumerChainlinkCancelledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _MultiwordConsumer.contract.FilterLogs(opts, "ChainlinkCancelled", idRule)
	if err != nil {
		return nil, err
	}
	return &MultiwordConsumerChainlinkCancelledIterator{contract: _MultiwordConsumer.contract, event: "ChainlinkCancelled", logs: logs, sub: sub}, nil
}

// WatchChainlinkCancelled is a free log subscription operation binding the contract event 0xe1fe3afa0f7f761ff0a8b89086790efd5140d2907ebd5b7ff6bfcb5e075fd4c5.
//
// Solidity: event ChainlinkCancelled(bytes32 indexed id)
func (_MultiwordConsumer *MultiwordConsumerFilterer) WatchChainlinkCancelled(opts *bind.WatchOpts, sink chan<- *MultiwordConsumerChainlinkCancelled, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _MultiwordConsumer.contract.WatchLogs(opts, "ChainlinkCancelled", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiwordConsumerChainlinkCancelled)
				if err := _MultiwordConsumer.contract.UnpackLog(event, "ChainlinkCancelled", log); err != nil {
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
func (_MultiwordConsumer *MultiwordConsumerFilterer) ParseChainlinkCancelled(log types.Log) (*MultiwordConsumerChainlinkCancelled, error) {
	event := new(MultiwordConsumerChainlinkCancelled)
	if err := _MultiwordConsumer.contract.UnpackLog(event, "ChainlinkCancelled", log); err != nil {
		return nil, err
	}
	return event, nil
}

// MultiwordConsumerChainlinkFulfilledIterator is returned from FilterChainlinkFulfilled and is used to iterate over the raw logs and unpacked data for ChainlinkFulfilled events raised by the MultiwordConsumer contract.
type MultiwordConsumerChainlinkFulfilledIterator struct {
	Event *MultiwordConsumerChainlinkFulfilled // Event containing the contract specifics and raw log

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
func (it *MultiwordConsumerChainlinkFulfilledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiwordConsumerChainlinkFulfilled)
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
		it.Event = new(MultiwordConsumerChainlinkFulfilled)
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
func (it *MultiwordConsumerChainlinkFulfilledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiwordConsumerChainlinkFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiwordConsumerChainlinkFulfilled represents a ChainlinkFulfilled event raised by the MultiwordConsumer contract.
type MultiwordConsumerChainlinkFulfilled struct {
	Id  [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterChainlinkFulfilled is a free log retrieval operation binding the contract event 0x7cc135e0cebb02c3480ae5d74d377283180a2601f8f644edf7987b009316c63a.
//
// Solidity: event ChainlinkFulfilled(bytes32 indexed id)
func (_MultiwordConsumer *MultiwordConsumerFilterer) FilterChainlinkFulfilled(opts *bind.FilterOpts, id [][32]byte) (*MultiwordConsumerChainlinkFulfilledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _MultiwordConsumer.contract.FilterLogs(opts, "ChainlinkFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return &MultiwordConsumerChainlinkFulfilledIterator{contract: _MultiwordConsumer.contract, event: "ChainlinkFulfilled", logs: logs, sub: sub}, nil
}

// WatchChainlinkFulfilled is a free log subscription operation binding the contract event 0x7cc135e0cebb02c3480ae5d74d377283180a2601f8f644edf7987b009316c63a.
//
// Solidity: event ChainlinkFulfilled(bytes32 indexed id)
func (_MultiwordConsumer *MultiwordConsumerFilterer) WatchChainlinkFulfilled(opts *bind.WatchOpts, sink chan<- *MultiwordConsumerChainlinkFulfilled, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _MultiwordConsumer.contract.WatchLogs(opts, "ChainlinkFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiwordConsumerChainlinkFulfilled)
				if err := _MultiwordConsumer.contract.UnpackLog(event, "ChainlinkFulfilled", log); err != nil {
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
func (_MultiwordConsumer *MultiwordConsumerFilterer) ParseChainlinkFulfilled(log types.Log) (*MultiwordConsumerChainlinkFulfilled, error) {
	event := new(MultiwordConsumerChainlinkFulfilled)
	if err := _MultiwordConsumer.contract.UnpackLog(event, "ChainlinkFulfilled", log); err != nil {
		return nil, err
	}
	return event, nil
}

// MultiwordConsumerChainlinkRequestedIterator is returned from FilterChainlinkRequested and is used to iterate over the raw logs and unpacked data for ChainlinkRequested events raised by the MultiwordConsumer contract.
type MultiwordConsumerChainlinkRequestedIterator struct {
	Event *MultiwordConsumerChainlinkRequested // Event containing the contract specifics and raw log

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
func (it *MultiwordConsumerChainlinkRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiwordConsumerChainlinkRequested)
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
		it.Event = new(MultiwordConsumerChainlinkRequested)
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
func (it *MultiwordConsumerChainlinkRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiwordConsumerChainlinkRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiwordConsumerChainlinkRequested represents a ChainlinkRequested event raised by the MultiwordConsumer contract.
type MultiwordConsumerChainlinkRequested struct {
	Id  [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterChainlinkRequested is a free log retrieval operation binding the contract event 0xb5e6e01e79f91267dc17b4e6314d5d4d03593d2ceee0fbb452b750bd70ea5af9.
//
// Solidity: event ChainlinkRequested(bytes32 indexed id)
func (_MultiwordConsumer *MultiwordConsumerFilterer) FilterChainlinkRequested(opts *bind.FilterOpts, id [][32]byte) (*MultiwordConsumerChainlinkRequestedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _MultiwordConsumer.contract.FilterLogs(opts, "ChainlinkRequested", idRule)
	if err != nil {
		return nil, err
	}
	return &MultiwordConsumerChainlinkRequestedIterator{contract: _MultiwordConsumer.contract, event: "ChainlinkRequested", logs: logs, sub: sub}, nil
}

// WatchChainlinkRequested is a free log subscription operation binding the contract event 0xb5e6e01e79f91267dc17b4e6314d5d4d03593d2ceee0fbb452b750bd70ea5af9.
//
// Solidity: event ChainlinkRequested(bytes32 indexed id)
func (_MultiwordConsumer *MultiwordConsumerFilterer) WatchChainlinkRequested(opts *bind.WatchOpts, sink chan<- *MultiwordConsumerChainlinkRequested, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _MultiwordConsumer.contract.WatchLogs(opts, "ChainlinkRequested", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiwordConsumerChainlinkRequested)
				if err := _MultiwordConsumer.contract.UnpackLog(event, "ChainlinkRequested", log); err != nil {
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
func (_MultiwordConsumer *MultiwordConsumerFilterer) ParseChainlinkRequested(log types.Log) (*MultiwordConsumerChainlinkRequested, error) {
	event := new(MultiwordConsumerChainlinkRequested)
	if err := _MultiwordConsumer.contract.UnpackLog(event, "ChainlinkRequested", log); err != nil {
		return nil, err
	}
	return event, nil
}

// MultiwordConsumerRequestFulfilledIterator is returned from FilterRequestFulfilled and is used to iterate over the raw logs and unpacked data for RequestFulfilled events raised by the MultiwordConsumer contract.
type MultiwordConsumerRequestFulfilledIterator struct {
	Event *MultiwordConsumerRequestFulfilled // Event containing the contract specifics and raw log

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
func (it *MultiwordConsumerRequestFulfilledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiwordConsumerRequestFulfilled)
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
		it.Event = new(MultiwordConsumerRequestFulfilled)
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
func (it *MultiwordConsumerRequestFulfilledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiwordConsumerRequestFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiwordConsumerRequestFulfilled represents a RequestFulfilled event raised by the MultiwordConsumer contract.
type MultiwordConsumerRequestFulfilled struct {
	RequestId [32]byte
	Price     common.Hash
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRequestFulfilled is a free log retrieval operation binding the contract event 0x1a111c5dcf9a71088bd5e1797fdfaf399fec2afbb24aca247e4e3e9f4b61df91.
//
// Solidity: event RequestFulfilled(bytes32 indexed requestId, bytes indexed price)
func (_MultiwordConsumer *MultiwordConsumerFilterer) FilterRequestFulfilled(opts *bind.FilterOpts, requestId [][32]byte, price [][]byte) (*MultiwordConsumerRequestFulfilledIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var priceRule []interface{}
	for _, priceItem := range price {
		priceRule = append(priceRule, priceItem)
	}

	logs, sub, err := _MultiwordConsumer.contract.FilterLogs(opts, "RequestFulfilled", requestIdRule, priceRule)
	if err != nil {
		return nil, err
	}
	return &MultiwordConsumerRequestFulfilledIterator{contract: _MultiwordConsumer.contract, event: "RequestFulfilled", logs: logs, sub: sub}, nil
}

// WatchRequestFulfilled is a free log subscription operation binding the contract event 0x1a111c5dcf9a71088bd5e1797fdfaf399fec2afbb24aca247e4e3e9f4b61df91.
//
// Solidity: event RequestFulfilled(bytes32 indexed requestId, bytes indexed price)
func (_MultiwordConsumer *MultiwordConsumerFilterer) WatchRequestFulfilled(opts *bind.WatchOpts, sink chan<- *MultiwordConsumerRequestFulfilled, requestId [][32]byte, price [][]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var priceRule []interface{}
	for _, priceItem := range price {
		priceRule = append(priceRule, priceItem)
	}

	logs, sub, err := _MultiwordConsumer.contract.WatchLogs(opts, "RequestFulfilled", requestIdRule, priceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiwordConsumerRequestFulfilled)
				if err := _MultiwordConsumer.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
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

// ParseRequestFulfilled is a log parse operation binding the contract event 0x1a111c5dcf9a71088bd5e1797fdfaf399fec2afbb24aca247e4e3e9f4b61df91.
//
// Solidity: event RequestFulfilled(bytes32 indexed requestId, bytes indexed price)
func (_MultiwordConsumer *MultiwordConsumerFilterer) ParseRequestFulfilled(log types.Log) (*MultiwordConsumerRequestFulfilled, error) {
	event := new(MultiwordConsumerRequestFulfilled)
	if err := _MultiwordConsumer.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
		return nil, err
	}
	return event, nil
}

// MultiwordConsumerRequestMultipleFulfilledIterator is returned from FilterRequestMultipleFulfilled and is used to iterate over the raw logs and unpacked data for RequestMultipleFulfilled events raised by the MultiwordConsumer contract.
type MultiwordConsumerRequestMultipleFulfilledIterator struct {
	Event *MultiwordConsumerRequestMultipleFulfilled // Event containing the contract specifics and raw log

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
func (it *MultiwordConsumerRequestMultipleFulfilledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiwordConsumerRequestMultipleFulfilled)
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
		it.Event = new(MultiwordConsumerRequestMultipleFulfilled)
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
func (it *MultiwordConsumerRequestMultipleFulfilledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiwordConsumerRequestMultipleFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiwordConsumerRequestMultipleFulfilled represents a RequestMultipleFulfilled event raised by the MultiwordConsumer contract.
type MultiwordConsumerRequestMultipleFulfilled struct {
	RequestId [32]byte
	First     [32]byte
	Second    [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRequestMultipleFulfilled is a free log retrieval operation binding the contract event 0xd368a628c6f427add4c36c69828a9be4d937a803adfda79c1dbf7eb26cdf4bc4.
//
// Solidity: event RequestMultipleFulfilled(bytes32 indexed requestId, bytes32 indexed first, bytes32 indexed second)
func (_MultiwordConsumer *MultiwordConsumerFilterer) FilterRequestMultipleFulfilled(opts *bind.FilterOpts, requestId [][32]byte, first [][32]byte, second [][32]byte) (*MultiwordConsumerRequestMultipleFulfilledIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var firstRule []interface{}
	for _, firstItem := range first {
		firstRule = append(firstRule, firstItem)
	}
	var secondRule []interface{}
	for _, secondItem := range second {
		secondRule = append(secondRule, secondItem)
	}

	logs, sub, err := _MultiwordConsumer.contract.FilterLogs(opts, "RequestMultipleFulfilled", requestIdRule, firstRule, secondRule)
	if err != nil {
		return nil, err
	}
	return &MultiwordConsumerRequestMultipleFulfilledIterator{contract: _MultiwordConsumer.contract, event: "RequestMultipleFulfilled", logs: logs, sub: sub}, nil
}

// WatchRequestMultipleFulfilled is a free log subscription operation binding the contract event 0xd368a628c6f427add4c36c69828a9be4d937a803adfda79c1dbf7eb26cdf4bc4.
//
// Solidity: event RequestMultipleFulfilled(bytes32 indexed requestId, bytes32 indexed first, bytes32 indexed second)
func (_MultiwordConsumer *MultiwordConsumerFilterer) WatchRequestMultipleFulfilled(opts *bind.WatchOpts, sink chan<- *MultiwordConsumerRequestMultipleFulfilled, requestId [][32]byte, first [][32]byte, second [][32]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var firstRule []interface{}
	for _, firstItem := range first {
		firstRule = append(firstRule, firstItem)
	}
	var secondRule []interface{}
	for _, secondItem := range second {
		secondRule = append(secondRule, secondItem)
	}

	logs, sub, err := _MultiwordConsumer.contract.WatchLogs(opts, "RequestMultipleFulfilled", requestIdRule, firstRule, secondRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiwordConsumerRequestMultipleFulfilled)
				if err := _MultiwordConsumer.contract.UnpackLog(event, "RequestMultipleFulfilled", log); err != nil {
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

// ParseRequestMultipleFulfilled is a log parse operation binding the contract event 0xd368a628c6f427add4c36c69828a9be4d937a803adfda79c1dbf7eb26cdf4bc4.
//
// Solidity: event RequestMultipleFulfilled(bytes32 indexed requestId, bytes32 indexed first, bytes32 indexed second)
func (_MultiwordConsumer *MultiwordConsumerFilterer) ParseRequestMultipleFulfilled(log types.Log) (*MultiwordConsumerRequestMultipleFulfilled, error) {
	event := new(MultiwordConsumerRequestMultipleFulfilled)
	if err := _MultiwordConsumer.contract.UnpackLog(event, "RequestMultipleFulfilled", log); err != nil {
		return nil, err
	}
	return event, nil
}
