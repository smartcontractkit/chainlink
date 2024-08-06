// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ccip_reader_tester

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

type EVM2EVMMultiOffRampCommitReport struct {
	PriceUpdates InternalPriceUpdates
	MerkleRoots  []EVM2EVMMultiOffRampMerkleRoot
}

type EVM2EVMMultiOffRampInterval struct {
	Min uint64
	Max uint64
}

type EVM2EVMMultiOffRampMerkleRoot struct {
	SourceChainSelector uint64
	Interval            EVM2EVMMultiOffRampInterval
	MerkleRoot          [32]byte
}

type EVM2EVMMultiOffRampSourceChainConfig struct {
	Router    common.Address
	IsEnabled bool
	MinSeqNr  uint64
	OnRamp    []byte
}

type InternalEVM2AnyRampMessage struct {
	Header         InternalRampMessageHeader
	Sender         common.Address
	Data           []byte
	Receiver       []byte
	ExtraArgs      []byte
	FeeToken       common.Address
	FeeTokenAmount *big.Int
	TokenAmounts   []InternalRampTokenAmount
}

type InternalGasPriceUpdate struct {
	DestChainSelector uint64
	UsdPerUnitGas     *big.Int
}

type InternalPriceUpdates struct {
	TokenPriceUpdates []InternalTokenPriceUpdate
	GasPriceUpdates   []InternalGasPriceUpdate
}

type InternalRampMessageHeader struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	DestChainSelector   uint64
	SequenceNumber      uint64
	Nonce               uint64
}

type InternalRampTokenAmount struct {
	SourcePoolAddress []byte
	DestTokenAddress  []byte
	ExtraData         []byte
	Amount            *big.Int
}

type InternalTokenPriceUpdate struct {
	SourceToken common.Address
	UsdPerToken *big.Int
}

var CCIPReaderTesterMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.RampTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structInternal.EVM2AnyRampMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"CCIPSendRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"min\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"max\",\"type\":\"uint64\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.Interval\",\"name\":\"interval\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structEVM2EVMMultiOffRamp.CommitReport\",\"name\":\"report\",\"type\":\"tuple\"}],\"name\":\"CommitReportAccepted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"state\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"name\":\"ExecutionStateChanged\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.RampTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.EVM2AnyRampMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"emitCCIPSendRequested\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"min\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"max\",\"type\":\"uint64\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.Interval\",\"name\":\"interval\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.CommitReport\",\"name\":\"report\",\"type\":\"tuple\"}],\"name\":\"emitCommitReportAccepted\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"state\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"name\":\"emitExecutionStateChanged\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"getSourceChainConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.SourceChainConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.SourceChainConfig\",\"name\":\"sourceChainConfig\",\"type\":\"tuple\"}],\"name\":\"setSourceChainConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50611149806100206000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c80634cf66e361461005c578063a65558f614610071578063e44302b714610084578063e83eabba14610097578063e9d68a8e146100aa575b600080fd5b61006f61006a3660046104a7565b6100d3565b005b61006f61007f366004610718565b610128565b61006f6100923660046109b6565b61016d565b61006f6100a5366004610b20565b6101a7565b6100bd6100b8366004610bdb565b610233565b6040516100ca9190610c43565b60405180910390f35b82846001600160401b0316866001600160401b03167f8c324ce1367b83031769f6a813e3bb4c117aba2185789d66b98b791405be6df28585604051610119929190610c9b565b60405180910390a45050505050565b816001600160401b03167f0f07cd31e53232da9125e517f09550fdde74bf43d6a0a76ebd41674dafe2ab29826040516101619190610d83565b60405180910390a25050565b7f3a3950e13dd607cc37980db0ef14266c40d2bba9c01b2e44bfe549808883095d8160405161019c9190610f3d565b60405180910390a150565b6001600160401b0380831660009081526020818152604091829020845181549286015193860151909416600160a81b02600160a81b600160e81b0319931515600160a01b026001600160a81b03199093166001600160a01b039095169490941791909117919091169190911781556060820151829190600182019061022c908261107d565b5050505050565b60408051608080820183526000808352602080840182905283850182905260608085018190526001600160401b038781168452838352928690208651948501875280546001600160a01b0381168652600160a01b810460ff16151593860193909352600160a81b9092049092169483019490945260018401805493949293918401916102be90610ff2565b80601f01602080910402602001604051908101604052809291908181526020018280546102ea90610ff2565b80156103375780601f1061030c57610100808354040283529160200191610337565b820191906000526020600020905b81548152906001019060200180831161031a57829003601f168201915b5050505050815250509050919050565b80356001600160401b038116811461035e57600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b604051608081016001600160401b038111828210171561039b5761039b610363565b60405290565b60405161010081016001600160401b038111828210171561039b5761039b610363565b604080519081016001600160401b038111828210171561039b5761039b610363565b604051606081016001600160401b038111828210171561039b5761039b610363565b604051601f8201601f191681016001600160401b038111828210171561043057610430610363565b604052919050565b600082601f83011261044957600080fd5b81356001600160401b0381111561046257610462610363565b610475601f8201601f1916602001610408565b81815284602083860101111561048a57600080fd5b816020850160208301376000918101602001919091529392505050565b600080600080600060a086880312156104bf57600080fd5b6104c886610347565b94506104d660208701610347565b9350604086013592506060860135600481106104f157600080fd5b915060808601356001600160401b0381111561050c57600080fd5b61051888828901610438565b9150509295509295909350565b600060a0828403121561053757600080fd5b60405160a081016001600160401b038111828210171561055957610559610363565b6040528235815290508061056f60208401610347565b602082015261058060408401610347565b604082015261059160608401610347565b60608201526105a260808401610347565b60808201525092915050565b6001600160a01b03811681146105c357600080fd5b50565b803561035e816105ae565b60006001600160401b038211156105ea576105ea610363565b5060051b60200190565b600082601f83011261060557600080fd5b8135602061061a610615836105d1565b610408565b82815260059290921b8401810191818101908684111561063957600080fd5b8286015b8481101561070d5780356001600160401b038082111561065d5760008081fd5b908801906080828b03601f19018113156106775760008081fd5b61067f610379565b87840135838111156106915760008081fd5b61069f8d8a83880101610438565b825250604080850135848111156106b65760008081fd5b6106c48e8b83890101610438565b8a84015250606080860135858111156106dd5760008081fd5b6106eb8f8c838a0101610438565b928401929092529490920135938101939093525050835291830191830161063d565b509695505050505050565b6000806040838503121561072b57600080fd5b61073483610347565b915060208301356001600160401b038082111561075057600080fd5b90840190610180828703121561076557600080fd5b61076d6103a1565b6107778784610525565b815261078560a084016105c6565b602082015260c08301358281111561079c57600080fd5b6107a888828601610438565b60408301525060e0830135828111156107c057600080fd5b6107cc88828601610438565b606083015250610100830135828111156107e557600080fd5b6107f188828601610438565b60808301525061080461012084016105c6565b60a082015261014083013560c08201526101608301358281111561082757600080fd5b610833888286016105f4565b60e0830152508093505050509250929050565b80356001600160e01b038116811461035e57600080fd5b600082601f83011261086e57600080fd5b8135602061087e610615836105d1565b82815260069290921b8401810191818101908684111561089d57600080fd5b8286015b8481101561070d57604081890312156108ba5760008081fd5b6108c26103c4565b6108cb82610347565b81526108d8858301610846565b818601528352918301916040016108a1565b600082601f8301126108fb57600080fd5b8135602061090b610615836105d1565b82815260079290921b8401810191818101908684111561092a57600080fd5b8286015b8481101561070d5780880360808112156109485760008081fd5b6109506103e6565b61095983610347565b8152604080601f198401121561096f5760008081fd5b6109776103c4565b9250610984878501610347565b8352610991818501610347565b838801528187019290925260608301359181019190915283529183019160800161092e565b600060208083850312156109c957600080fd5b82356001600160401b03808211156109e057600080fd5b818501915060408083880312156109f657600080fd5b6109fe6103c4565b833583811115610a0d57600080fd5b84016040818a031215610a1f57600080fd5b610a276103c4565b813585811115610a3657600080fd5b8201601f81018b13610a4757600080fd5b8035610a55610615826105d1565b81815260069190911b8201890190898101908d831115610a7457600080fd5b928a01925b82841015610ac45787848f031215610a915760008081fd5b610a996103c4565b8435610aa4816105ae565b8152610ab1858d01610846565b818d0152825292870192908a0190610a79565b845250505081870135935084841115610adc57600080fd5b610ae88a85840161085d565b8188015282525083850135915082821115610b0257600080fd5b610b0e888386016108ea565b85820152809550505050505092915050565b60008060408385031215610b3357600080fd5b610b3c83610347565b915060208301356001600160401b0380821115610b5857600080fd5b9084019060808287031215610b6c57600080fd5b610b74610379565b8235610b7f816105ae565b815260208301358015158114610b9457600080fd5b6020820152610ba560408401610347565b6040820152606083013582811115610bbc57600080fd5b610bc888828601610438565b6060830152508093505050509250929050565b600060208284031215610bed57600080fd5b610bf682610347565b9392505050565b6000815180845260005b81811015610c2357602081850181015186830182015201610c07565b506000602082860101526020601f19601f83011685010191505092915050565b602080825282516001600160a01b03168282015282015115156040808301919091528201516001600160401b0316606080830191909152820151608080830152600090610c9360a0840182610bfd565b949350505050565b600060048410610cbb57634e487b7160e01b600052602160045260246000fd5b83825260406020830152610c936040830184610bfd565b6001600160a01b03169052565b600082825180855260208086019550808260051b84010181860160005b84811015610d7657601f19868403018952815160808151818652610d2282870182610bfd565b9150508582015185820387870152610d3a8282610bfd565b91505060408083015186830382880152610d548382610bfd565b6060948501519790940196909652505098840198925090830190600101610cfc565b5090979650505050505050565b60208152610dd0602082018351805182526020808201516001600160401b039081169184019190915260408083015182169084015260608083015182169084015260809182015116910152565b60006020830151610de460c0840182610cd2565b5060408301516101808060e0850152610e016101a0850183610bfd565b91506060850151601f198086850301610100870152610e208483610bfd565b9350608087015191508086850301610120870152610e3e8483610bfd565b935060a08701519150610e55610140870183610cd2565b60c087015161016087015260e0870151915080868503018387015250610e7b8382610cdf565b9695505050505050565b60008151808452602080850194506020840160005b83811015610ed357815180516001600160401b031688528301516001600160e01b03168388015260409096019590820190600101610e9a565b509495945050505050565b600081518084526020808501945080840160005b83811015610ed357815180516001600160401b0390811689528482015180518216868b0152850151166040898101919091520151606088015260809096019590820190600101610ef2565b6000602080835283516040808386015260a0850182516040606088015281815180845260c0890191508683019350600092505b80831015610fab57835180516001600160a01b031683528701516001600160e01b031687830152928601926001929092019190840190610f70565b5093850151878503605f1901608089015293610fc78186610e85565b945050505050818501519150601f19848203016040850152610fe98183610ede565b95945050505050565b600181811c9082168061100657607f821691505b60208210810361102657634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115611078576000816000526020600020601f850160051c810160208610156110555750805b601f850160051c820191505b8181101561107457828155600101611061565b5050505b505050565b81516001600160401b0381111561109657611096610363565b6110aa816110a48454610ff2565b8461102c565b602080601f8311600181146110df57600084156110c75750858301515b600019600386901b1c1916600185901b178555611074565b600085815260208120601f198616915b8281101561110e578886015182559484019460019091019084016110ef565b508582101561112c5787850151600019600388901b60f8161c191681555b5050505050600190811b0190555056fea164736f6c6343000818000a",
}

var CCIPReaderTesterABI = CCIPReaderTesterMetaData.ABI

var CCIPReaderTesterBin = CCIPReaderTesterMetaData.Bin

func DeployCCIPReaderTester(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CCIPReaderTester, error) {
	parsed, err := CCIPReaderTesterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CCIPReaderTesterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CCIPReaderTester{address: address, abi: *parsed, CCIPReaderTesterCaller: CCIPReaderTesterCaller{contract: contract}, CCIPReaderTesterTransactor: CCIPReaderTesterTransactor{contract: contract}, CCIPReaderTesterFilterer: CCIPReaderTesterFilterer{contract: contract}}, nil
}

type CCIPReaderTester struct {
	address common.Address
	abi     abi.ABI
	CCIPReaderTesterCaller
	CCIPReaderTesterTransactor
	CCIPReaderTesterFilterer
}

type CCIPReaderTesterCaller struct {
	contract *bind.BoundContract
}

type CCIPReaderTesterTransactor struct {
	contract *bind.BoundContract
}

type CCIPReaderTesterFilterer struct {
	contract *bind.BoundContract
}

type CCIPReaderTesterSession struct {
	Contract     *CCIPReaderTester
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type CCIPReaderTesterCallerSession struct {
	Contract *CCIPReaderTesterCaller
	CallOpts bind.CallOpts
}

type CCIPReaderTesterTransactorSession struct {
	Contract     *CCIPReaderTesterTransactor
	TransactOpts bind.TransactOpts
}

type CCIPReaderTesterRaw struct {
	Contract *CCIPReaderTester
}

type CCIPReaderTesterCallerRaw struct {
	Contract *CCIPReaderTesterCaller
}

type CCIPReaderTesterTransactorRaw struct {
	Contract *CCIPReaderTesterTransactor
}

func NewCCIPReaderTester(address common.Address, backend bind.ContractBackend) (*CCIPReaderTester, error) {
	abi, err := abi.JSON(strings.NewReader(CCIPReaderTesterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindCCIPReaderTester(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CCIPReaderTester{address: address, abi: abi, CCIPReaderTesterCaller: CCIPReaderTesterCaller{contract: contract}, CCIPReaderTesterTransactor: CCIPReaderTesterTransactor{contract: contract}, CCIPReaderTesterFilterer: CCIPReaderTesterFilterer{contract: contract}}, nil
}

func NewCCIPReaderTesterCaller(address common.Address, caller bind.ContractCaller) (*CCIPReaderTesterCaller, error) {
	contract, err := bindCCIPReaderTester(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CCIPReaderTesterCaller{contract: contract}, nil
}

func NewCCIPReaderTesterTransactor(address common.Address, transactor bind.ContractTransactor) (*CCIPReaderTesterTransactor, error) {
	contract, err := bindCCIPReaderTester(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CCIPReaderTesterTransactor{contract: contract}, nil
}

func NewCCIPReaderTesterFilterer(address common.Address, filterer bind.ContractFilterer) (*CCIPReaderTesterFilterer, error) {
	contract, err := bindCCIPReaderTester(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CCIPReaderTesterFilterer{contract: contract}, nil
}

func bindCCIPReaderTester(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CCIPReaderTesterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_CCIPReaderTester *CCIPReaderTesterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CCIPReaderTester.Contract.CCIPReaderTesterCaller.contract.Call(opts, result, method, params...)
}

func (_CCIPReaderTester *CCIPReaderTesterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.CCIPReaderTesterTransactor.contract.Transfer(opts)
}

func (_CCIPReaderTester *CCIPReaderTesterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.CCIPReaderTesterTransactor.contract.Transact(opts, method, params...)
}

func (_CCIPReaderTester *CCIPReaderTesterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CCIPReaderTester.Contract.contract.Call(opts, result, method, params...)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.contract.Transfer(opts)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.contract.Transact(opts, method, params...)
}

func (_CCIPReaderTester *CCIPReaderTesterCaller) GetSourceChainConfig(opts *bind.CallOpts, sourceChainSelector uint64) (EVM2EVMMultiOffRampSourceChainConfig, error) {
	var out []interface{}
	err := _CCIPReaderTester.contract.Call(opts, &out, "getSourceChainConfig", sourceChainSelector)

	if err != nil {
		return *new(EVM2EVMMultiOffRampSourceChainConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(EVM2EVMMultiOffRampSourceChainConfig)).(*EVM2EVMMultiOffRampSourceChainConfig)

	return out0, err

}

func (_CCIPReaderTester *CCIPReaderTesterSession) GetSourceChainConfig(sourceChainSelector uint64) (EVM2EVMMultiOffRampSourceChainConfig, error) {
	return _CCIPReaderTester.Contract.GetSourceChainConfig(&_CCIPReaderTester.CallOpts, sourceChainSelector)
}

func (_CCIPReaderTester *CCIPReaderTesterCallerSession) GetSourceChainConfig(sourceChainSelector uint64) (EVM2EVMMultiOffRampSourceChainConfig, error) {
	return _CCIPReaderTester.Contract.GetSourceChainConfig(&_CCIPReaderTester.CallOpts, sourceChainSelector)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactor) EmitCCIPSendRequested(opts *bind.TransactOpts, destChainSelector uint64, message InternalEVM2AnyRampMessage) (*types.Transaction, error) {
	return _CCIPReaderTester.contract.Transact(opts, "emitCCIPSendRequested", destChainSelector, message)
}

func (_CCIPReaderTester *CCIPReaderTesterSession) EmitCCIPSendRequested(destChainSelector uint64, message InternalEVM2AnyRampMessage) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.EmitCCIPSendRequested(&_CCIPReaderTester.TransactOpts, destChainSelector, message)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactorSession) EmitCCIPSendRequested(destChainSelector uint64, message InternalEVM2AnyRampMessage) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.EmitCCIPSendRequested(&_CCIPReaderTester.TransactOpts, destChainSelector, message)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactor) EmitCommitReportAccepted(opts *bind.TransactOpts, report EVM2EVMMultiOffRampCommitReport) (*types.Transaction, error) {
	return _CCIPReaderTester.contract.Transact(opts, "emitCommitReportAccepted", report)
}

func (_CCIPReaderTester *CCIPReaderTesterSession) EmitCommitReportAccepted(report EVM2EVMMultiOffRampCommitReport) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.EmitCommitReportAccepted(&_CCIPReaderTester.TransactOpts, report)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactorSession) EmitCommitReportAccepted(report EVM2EVMMultiOffRampCommitReport) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.EmitCommitReportAccepted(&_CCIPReaderTester.TransactOpts, report)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactor) EmitExecutionStateChanged(opts *bind.TransactOpts, sourceChainSelector uint64, sequenceNumber uint64, messageId [32]byte, state uint8, returnData []byte) (*types.Transaction, error) {
	return _CCIPReaderTester.contract.Transact(opts, "emitExecutionStateChanged", sourceChainSelector, sequenceNumber, messageId, state, returnData)
}

func (_CCIPReaderTester *CCIPReaderTesterSession) EmitExecutionStateChanged(sourceChainSelector uint64, sequenceNumber uint64, messageId [32]byte, state uint8, returnData []byte) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.EmitExecutionStateChanged(&_CCIPReaderTester.TransactOpts, sourceChainSelector, sequenceNumber, messageId, state, returnData)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactorSession) EmitExecutionStateChanged(sourceChainSelector uint64, sequenceNumber uint64, messageId [32]byte, state uint8, returnData []byte) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.EmitExecutionStateChanged(&_CCIPReaderTester.TransactOpts, sourceChainSelector, sequenceNumber, messageId, state, returnData)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactor) SetSourceChainConfig(opts *bind.TransactOpts, sourceChainSelector uint64, sourceChainConfig EVM2EVMMultiOffRampSourceChainConfig) (*types.Transaction, error) {
	return _CCIPReaderTester.contract.Transact(opts, "setSourceChainConfig", sourceChainSelector, sourceChainConfig)
}

func (_CCIPReaderTester *CCIPReaderTesterSession) SetSourceChainConfig(sourceChainSelector uint64, sourceChainConfig EVM2EVMMultiOffRampSourceChainConfig) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.SetSourceChainConfig(&_CCIPReaderTester.TransactOpts, sourceChainSelector, sourceChainConfig)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactorSession) SetSourceChainConfig(sourceChainSelector uint64, sourceChainConfig EVM2EVMMultiOffRampSourceChainConfig) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.SetSourceChainConfig(&_CCIPReaderTester.TransactOpts, sourceChainSelector, sourceChainConfig)
}

type CCIPReaderTesterCCIPSendRequestedIterator struct {
	Event *CCIPReaderTesterCCIPSendRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CCIPReaderTesterCCIPSendRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPReaderTesterCCIPSendRequested)
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
		it.Event = new(CCIPReaderTesterCCIPSendRequested)
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

func (it *CCIPReaderTesterCCIPSendRequestedIterator) Error() error {
	return it.fail
}

func (it *CCIPReaderTesterCCIPSendRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CCIPReaderTesterCCIPSendRequested struct {
	DestChainSelector uint64
	Message           InternalEVM2AnyRampMessage
	Raw               types.Log
}

func (_CCIPReaderTester *CCIPReaderTesterFilterer) FilterCCIPSendRequested(opts *bind.FilterOpts, destChainSelector []uint64) (*CCIPReaderTesterCCIPSendRequestedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _CCIPReaderTester.contract.FilterLogs(opts, "CCIPSendRequested", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &CCIPReaderTesterCCIPSendRequestedIterator{contract: _CCIPReaderTester.contract, event: "CCIPSendRequested", logs: logs, sub: sub}, nil
}

func (_CCIPReaderTester *CCIPReaderTesterFilterer) WatchCCIPSendRequested(opts *bind.WatchOpts, sink chan<- *CCIPReaderTesterCCIPSendRequested, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _CCIPReaderTester.contract.WatchLogs(opts, "CCIPSendRequested", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CCIPReaderTesterCCIPSendRequested)
				if err := _CCIPReaderTester.contract.UnpackLog(event, "CCIPSendRequested", log); err != nil {
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

func (_CCIPReaderTester *CCIPReaderTesterFilterer) ParseCCIPSendRequested(log types.Log) (*CCIPReaderTesterCCIPSendRequested, error) {
	event := new(CCIPReaderTesterCCIPSendRequested)
	if err := _CCIPReaderTester.contract.UnpackLog(event, "CCIPSendRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CCIPReaderTesterCommitReportAcceptedIterator struct {
	Event *CCIPReaderTesterCommitReportAccepted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CCIPReaderTesterCommitReportAcceptedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPReaderTesterCommitReportAccepted)
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
		it.Event = new(CCIPReaderTesterCommitReportAccepted)
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

func (it *CCIPReaderTesterCommitReportAcceptedIterator) Error() error {
	return it.fail
}

func (it *CCIPReaderTesterCommitReportAcceptedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CCIPReaderTesterCommitReportAccepted struct {
	Report EVM2EVMMultiOffRampCommitReport
	Raw    types.Log
}

func (_CCIPReaderTester *CCIPReaderTesterFilterer) FilterCommitReportAccepted(opts *bind.FilterOpts) (*CCIPReaderTesterCommitReportAcceptedIterator, error) {

	logs, sub, err := _CCIPReaderTester.contract.FilterLogs(opts, "CommitReportAccepted")
	if err != nil {
		return nil, err
	}
	return &CCIPReaderTesterCommitReportAcceptedIterator{contract: _CCIPReaderTester.contract, event: "CommitReportAccepted", logs: logs, sub: sub}, nil
}

func (_CCIPReaderTester *CCIPReaderTesterFilterer) WatchCommitReportAccepted(opts *bind.WatchOpts, sink chan<- *CCIPReaderTesterCommitReportAccepted) (event.Subscription, error) {

	logs, sub, err := _CCIPReaderTester.contract.WatchLogs(opts, "CommitReportAccepted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CCIPReaderTesterCommitReportAccepted)
				if err := _CCIPReaderTester.contract.UnpackLog(event, "CommitReportAccepted", log); err != nil {
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

func (_CCIPReaderTester *CCIPReaderTesterFilterer) ParseCommitReportAccepted(log types.Log) (*CCIPReaderTesterCommitReportAccepted, error) {
	event := new(CCIPReaderTesterCommitReportAccepted)
	if err := _CCIPReaderTester.contract.UnpackLog(event, "CommitReportAccepted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CCIPReaderTesterExecutionStateChangedIterator struct {
	Event *CCIPReaderTesterExecutionStateChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CCIPReaderTesterExecutionStateChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPReaderTesterExecutionStateChanged)
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
		it.Event = new(CCIPReaderTesterExecutionStateChanged)
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

func (it *CCIPReaderTesterExecutionStateChangedIterator) Error() error {
	return it.fail
}

func (it *CCIPReaderTesterExecutionStateChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CCIPReaderTesterExecutionStateChanged struct {
	SourceChainSelector uint64
	SequenceNumber      uint64
	MessageId           [32]byte
	State               uint8
	ReturnData          []byte
	Raw                 types.Log
}

func (_CCIPReaderTester *CCIPReaderTesterFilterer) FilterExecutionStateChanged(opts *bind.FilterOpts, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (*CCIPReaderTesterExecutionStateChangedIterator, error) {

	var sourceChainSelectorRule []interface{}
	for _, sourceChainSelectorItem := range sourceChainSelector {
		sourceChainSelectorRule = append(sourceChainSelectorRule, sourceChainSelectorItem)
	}
	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}
	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}

	logs, sub, err := _CCIPReaderTester.contract.FilterLogs(opts, "ExecutionStateChanged", sourceChainSelectorRule, sequenceNumberRule, messageIdRule)
	if err != nil {
		return nil, err
	}
	return &CCIPReaderTesterExecutionStateChangedIterator{contract: _CCIPReaderTester.contract, event: "ExecutionStateChanged", logs: logs, sub: sub}, nil
}

func (_CCIPReaderTester *CCIPReaderTesterFilterer) WatchExecutionStateChanged(opts *bind.WatchOpts, sink chan<- *CCIPReaderTesterExecutionStateChanged, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (event.Subscription, error) {

	var sourceChainSelectorRule []interface{}
	for _, sourceChainSelectorItem := range sourceChainSelector {
		sourceChainSelectorRule = append(sourceChainSelectorRule, sourceChainSelectorItem)
	}
	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}
	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}

	logs, sub, err := _CCIPReaderTester.contract.WatchLogs(opts, "ExecutionStateChanged", sourceChainSelectorRule, sequenceNumberRule, messageIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CCIPReaderTesterExecutionStateChanged)
				if err := _CCIPReaderTester.contract.UnpackLog(event, "ExecutionStateChanged", log); err != nil {
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

func (_CCIPReaderTester *CCIPReaderTesterFilterer) ParseExecutionStateChanged(log types.Log) (*CCIPReaderTesterExecutionStateChanged, error) {
	event := new(CCIPReaderTesterExecutionStateChanged)
	if err := _CCIPReaderTester.contract.UnpackLog(event, "ExecutionStateChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_CCIPReaderTester *CCIPReaderTester) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _CCIPReaderTester.abi.Events["CCIPSendRequested"].ID:
		return _CCIPReaderTester.ParseCCIPSendRequested(log)
	case _CCIPReaderTester.abi.Events["CommitReportAccepted"].ID:
		return _CCIPReaderTester.ParseCommitReportAccepted(log)
	case _CCIPReaderTester.abi.Events["ExecutionStateChanged"].ID:
		return _CCIPReaderTester.ParseExecutionStateChanged(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (CCIPReaderTesterCCIPSendRequested) Topic() common.Hash {
	return common.HexToHash("0x0f07cd31e53232da9125e517f09550fdde74bf43d6a0a76ebd41674dafe2ab29")
}

func (CCIPReaderTesterCommitReportAccepted) Topic() common.Hash {
	return common.HexToHash("0x3a3950e13dd607cc37980db0ef14266c40d2bba9c01b2e44bfe549808883095d")
}

func (CCIPReaderTesterExecutionStateChanged) Topic() common.Hash {
	return common.HexToHash("0x8c324ce1367b83031769f6a813e3bb4c117aba2185789d66b98b791405be6df2")
}

func (_CCIPReaderTester *CCIPReaderTester) Address() common.Address {
	return _CCIPReaderTester.address
}

type CCIPReaderTesterInterface interface {
	GetSourceChainConfig(opts *bind.CallOpts, sourceChainSelector uint64) (EVM2EVMMultiOffRampSourceChainConfig, error)

	EmitCCIPSendRequested(opts *bind.TransactOpts, destChainSelector uint64, message InternalEVM2AnyRampMessage) (*types.Transaction, error)

	EmitCommitReportAccepted(opts *bind.TransactOpts, report EVM2EVMMultiOffRampCommitReport) (*types.Transaction, error)

	EmitExecutionStateChanged(opts *bind.TransactOpts, sourceChainSelector uint64, sequenceNumber uint64, messageId [32]byte, state uint8, returnData []byte) (*types.Transaction, error)

	SetSourceChainConfig(opts *bind.TransactOpts, sourceChainSelector uint64, sourceChainConfig EVM2EVMMultiOffRampSourceChainConfig) (*types.Transaction, error)

	FilterCCIPSendRequested(opts *bind.FilterOpts, destChainSelector []uint64) (*CCIPReaderTesterCCIPSendRequestedIterator, error)

	WatchCCIPSendRequested(opts *bind.WatchOpts, sink chan<- *CCIPReaderTesterCCIPSendRequested, destChainSelector []uint64) (event.Subscription, error)

	ParseCCIPSendRequested(log types.Log) (*CCIPReaderTesterCCIPSendRequested, error)

	FilterCommitReportAccepted(opts *bind.FilterOpts) (*CCIPReaderTesterCommitReportAcceptedIterator, error)

	WatchCommitReportAccepted(opts *bind.WatchOpts, sink chan<- *CCIPReaderTesterCommitReportAccepted) (event.Subscription, error)

	ParseCommitReportAccepted(log types.Log) (*CCIPReaderTesterCommitReportAccepted, error)

	FilterExecutionStateChanged(opts *bind.FilterOpts, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (*CCIPReaderTesterExecutionStateChangedIterator, error)

	WatchExecutionStateChanged(opts *bind.WatchOpts, sink chan<- *CCIPReaderTesterExecutionStateChanged, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (event.Subscription, error)

	ParseExecutionStateChanged(log types.Log) (*CCIPReaderTesterExecutionStateChanged, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
