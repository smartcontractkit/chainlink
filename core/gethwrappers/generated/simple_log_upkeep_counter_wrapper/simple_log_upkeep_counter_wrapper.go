// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package simple_log_upkeep_counter_wrapper

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

type CheckData struct {
	CheckBurnAmount   *big.Int
	PerformBurnAmount *big.Int
	EventSig          [32]byte
}

type Log struct {
	Index       *big.Int
	Timestamp   *big.Int
	TxHash      [32]byte
	BlockNumber *big.Int
	BlockHash   [32]byte
	Source      common.Address
	Topics      [][32]byte
	Data        []byte
}

var SimpleLogUpkeepCounterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_isStreamsLookup\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"StreamsLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lastBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timeToPerform\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isRecovered\",\"type\":\"bool\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"checkBurnAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"performBurnAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"eventSig\",\"type\":\"bytes32\"}],\"internalType\":\"structCheckData\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_checkDataConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"errCode\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkErrorHandler\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"topics\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLog\",\"name\":\"log\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkLog\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"dummyMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feedsHex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isStreamsLookup\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParam\",\"type\":\"string\"}],\"name\":\"setFeedParamKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"newFeeds\",\"type\":\"string[]\"}],\"name\":\"setFeedsHex\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"setShouldRetryOnErrorBool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"timeParam\",\"type\":\"string\"}],\"name\":\"setTimeParamKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"shouldRetryOnError\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeToPerform\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60e0604052600860a090815267030783030303230360c41b60c05260809081526200002f906007906001620000f3565b50604080518082019091526007808252666665656449447360c81b6020909201918252620000609160089162000157565b5060408051808201909152600980825268074696d657374616d760bc1b602090920191825262000091918162000157565b503480156200009f57600080fd5b5060405162001b9c38038062001b9c833981016040819052620000c2916200025c565b60006002819055436001556003819055600455600680549115156101000261ff0019909216919091179055620002c4565b82805482825590600052602060002090810192821562000145579160200282015b828111156200014557825180516200013491849160209091019062000157565b509160200191906001019062000114565b5062000153929150620001e2565b5090565b828054620001659062000287565b90600052602060002090601f016020900481019282620001895760008555620001d4565b82601f10620001a457805160ff1916838001178555620001d4565b82800160010185558215620001d4579182015b82811115620001d4578251825591602001919060010190620001b7565b506200015392915062000203565b8082111562000153576000620001f982826200021a565b50600101620001e2565b5b8082111562000153576000815560010162000204565b508054620002289062000287565b6000825580601f1062000239575050565b601f01602090049060005260206000209081019062000259919062000203565b50565b6000602082840312156200026f57600080fd5b815180151581146200028057600080fd5b9392505050565b600181811c908216806200029c57607f821691505b60208210811415620002be57634e487b7160e01b600052602260045260246000fd5b50919050565b6118c880620002d46000396000f3fe608060405234801561001057600080fd5b506004361061016c5760003560e01c806361bc221a116100cd5780639525d57411610081578063afb28d1f11610066578063afb28d1f14610310578063c6066f0d14610318578063c98f10b01461032157600080fd5b80639525d574146102dd5780639d6f1cc7146102f057600080fd5b80637145f11b116100b25780637145f11b146102a8578063806b984f146102cb578063917d895f146102d457600080fd5b806361bc221a1461028e578063697794731461029757600080fd5b806340691db4116101245780634585e33b116101095780634585e33b146102555780634b56a42e14610268578063601d5a711461027b57600080fd5b806340691db41461022f57806342eb3d921461024257600080fd5b806313fab5901161015557806313fab590146101b057806323148cee146101f65780632cb158641461021857600080fd5b806305e25131146101715780630fb172fb14610186575b600080fd5b61018461017f366004610eaa565b610329565b005b6101996101943660046110f1565b610340565b6040516101a79291906113f7565b60405180910390f35b6101846101be366004610f62565b6006805491151562010000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffff909216919091179055565b60065461020890610100900460ff1681565b60405190151581526020016101a7565b61022160035481565b6040519081526020016101a7565b61019961023d36600461107f565b610447565b6006546102089062010000900460ff1681565b610184610263366004610fa4565b6106cd565b610199610276366004610d1e565b610898565b610184610289366004610fe6565b6108ec565b61022160045481565b6101846102a5366004611023565b50565b6102086102b6366004610f8b565b60006020819052908152604090205460ff1681565b61022160015481565b61022160025481565b6101846102eb366004610fe6565b6108ff565b6103036102fe366004610f8b565b610912565b6040516101a79190611412565b6103036109be565b61022160055481565b6103036109cb565b805161033c9060079060208401906109d8565b5050565b6040805160028082526060828101909352600092918391816020015b606081526020019060019003908161035c575050604080516020810188905291925001604051602081830303815290604052816000815181106103a1576103a161185d565b6020026020010181905250836040516020016103bd9190611412565b604051602081830303815290604052816001815181106103df576103df61185d565b6020026020010181905250600081856040516020016103ff929190611363565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291905260065462010000900460ff169450925050505b9250929050565b6000606081808061045a8688018861112e565b92509250925060005a905060006104726001436117c3565b409050600085156104e1575b855a61048a90856117c3565b10156104e1578080156104ab575060008281526020819052604090205460ff165b6040805160208101859052309181019190915290915060600160405160208183030381529060405280519060200120915061047e565b60408051600280825260608201909252600091816020015b60608152602001906001900390816104f9579050506040805160006020820152919250016040516020818303038152906040528160008151811061053f5761053f61185d565b60200260200101819052506000604051602001610565919060ff91909116815260200190565b604051602081830303815290604052816001815181106105875761058761185d565b602002602001018190525060008c438d8d6040516020016105ab94939291906114e8565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190529050856105e960c08f018f6115db565b60028181106105fa576105fa61185d565b90506020020135141561069457600654610100900460ff161561065f5760086007600943846040517ff055e4a2000000000000000000000000000000000000000000000000000000008152600401610656959493929190611425565b60405180910390fd5b60018282604051602001610674929190611363565b6040516020818303038152906040529950995050505050505050506106c5565b600082826040516020016106a9929190611363565b6040516020818303038152906040529950995050505050505050505b935093915050565b6003546106d957436003555b4360019081556004546106eb916117ab565b600455600154600255600080808061070585870187610d82565b935093509350935082602001514261071d91906117c3565b600555600680547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556060830151821461077f57600680547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790555b600080600083806020019051810190610798919061115a565b92509250925060005a905060006107b06001436117c3565b4090506000841561081f575b845a6107c890856117c3565b101561081f578080156107e9575060008281526020819052604090205460ff165b604080516020810185905230918101919091529091506060016040516020818303038152906040528051906020012091506107bc565b600354600154600254600454600554600654604080519687526020870195909552938501929092526060840152608083015260ff16151560a082015232907f29eff4cb37911c3ea85db4630638cc5474fdd0631ec42215aef1d7ec96c8e63d9060c00160405180910390a2505050505050505050505050565b60006060600084846040516020016108b1929190611363565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529190526001969095509350505050565b805161033c906009906020840190610a35565b805161033c906008906020840190610a35565b6007818154811061092257600080fd5b90600052602060002001600091509050805461093d906117da565b80601f0160208091040260200160405190810160405280929190818152602001828054610969906117da565b80156109b65780601f1061098b576101008083540402835291602001916109b6565b820191906000526020600020905b81548152906001019060200180831161099957829003601f168201915b505050505081565b6008805461093d906117da565b6009805461093d906117da565b828054828255906000526020600020908101928215610a25579160200282015b82811115610a255782518051610a15918491602090910190610a35565b50916020019190600101906109f8565b50610a31929150610ab5565b5090565b828054610a41906117da565b90600052602060002090601f016020900481019282610a635760008555610aa9565b82601f10610a7c57805160ff1916838001178555610aa9565b82800160010185558215610aa9579182015b82811115610aa9578251825591602001919060010190610a8e565b50610a31929150610ad2565b80821115610a31576000610ac98282610ae7565b50600101610ab5565b5b80821115610a315760008155600101610ad3565b508054610af3906117da565b6000825580601f10610b03575050565b601f0160209004906000526020600020908101906102a59190610ad2565b803573ffffffffffffffffffffffffffffffffffffffff81168114610b4557600080fd5b919050565b600082601f830112610b5b57600080fd5b81356020610b70610b6b836116bc565b61166d565b80838252828201915082860187848660051b8901011115610b9057600080fd5b60005b85811015610baf57813584529284019290840190600101610b93565b5090979650505050505050565b600082601f830112610bcd57600080fd5b81356020610bdd610b6b836116bc565b80838252828201915082860187848660051b8901011115610bfd57600080fd5b6000805b86811015610c4057823567ffffffffffffffff811115610c1f578283fd5b610c2d8b88838d0101610c90565b8652509385019391850191600101610c01565b509198975050505050505050565b60008083601f840112610c6057600080fd5b50813567ffffffffffffffff811115610c7857600080fd5b60208301915083602082850101111561044057600080fd5b600082601f830112610ca157600080fd5b813567ffffffffffffffff811115610cbb57610cbb61188c565b610cec60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8401160161166d565b818152846020838601011115610d0157600080fd5b816020850160208301376000918101602001919091529392505050565b60008060408385031215610d3157600080fd5b823567ffffffffffffffff80821115610d4957600080fd5b610d5586838701610bbc565b93506020850135915080821115610d6b57600080fd5b50610d7885828601610c90565b9150509250929050565b60008060008060808587031215610d9857600080fd5b843567ffffffffffffffff80821115610db057600080fd5b610dbc88838901610bbc565b95506020870135915080821115610dd257600080fd5b908601906101008289031215610de757600080fd5b610def611643565b8235815260208301356020820152604083013560408201526060830135606082015260808301356080820152610e2760a08401610b21565b60a082015260c083013582811115610e3e57600080fd5b610e4a8a828601610b4a565b60c08301525060e083013582811115610e6257600080fd5b610e6e8a828601610c90565b60e0830152509450604087013593506060870135915080821115610e9157600080fd5b50610e9e87828801610c90565b91505092959194509250565b60006020808385031215610ebd57600080fd5b823567ffffffffffffffff80821115610ed557600080fd5b818501915085601f830112610ee957600080fd5b8135610ef7610b6b826116bc565b80828252858201915085850189878560051b8801011115610f1757600080fd5b6000805b85811015610f5257823587811115610f31578283fd5b610f3f8d8b838c0101610c90565b8652509388019391880191600101610f1b565b50919a9950505050505050505050565b600060208284031215610f7457600080fd5b81358015158114610f8457600080fd5b9392505050565b600060208284031215610f9d57600080fd5b5035919050565b60008060208385031215610fb757600080fd5b823567ffffffffffffffff811115610fce57600080fd5b610fda85828601610c4e565b90969095509350505050565b600060208284031215610ff857600080fd5b813567ffffffffffffffff81111561100f57600080fd5b61101b84828501610c90565b949350505050565b60006060828403121561103557600080fd5b6040516060810181811067ffffffffffffffff821117156110585761105861188c565b80604052508235815260208301356020820152604083013560408201528091505092915050565b60008060006040848603121561109457600080fd5b833567ffffffffffffffff808211156110ac57600080fd5b9085019061010082880312156110c157600080fd5b909350602085013590808211156110d757600080fd5b506110e486828701610c4e565b9497909650939450505050565b6000806040838503121561110457600080fd5b82359150602083013567ffffffffffffffff81111561112257600080fd5b610d7885828601610c90565b60008060006060848603121561114357600080fd5b505081359360208301359350604090920135919050565b60008060006060848603121561116f57600080fd5b8351925060208401519150604084015190509250925092565b81835260007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8311156111ba57600080fd5b8260051b8083602087013760009401602001938452509192915050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b6000815180845260005b818110156112465760208185018101518683018201520161122a565b81811115611258576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8054600090600181811c90808316806112a557607f831692505b60208084108214156112e0577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b8388528180156112f7576001811461132957611357565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008616828a0152604089019650611357565b876000528160002060005b8681101561134f5781548b8201850152908501908301611334565b8a0183019750505b50505050505092915050565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b838110156113d8577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa08887030185526113c6868351611220565b9550938201939082019060010161138c565b5050858403818701525050506113ee8185611220565b95945050505050565b821515815260406020820152600061101b6040830184611220565b602081526000610f846020830184611220565b60a08152600061143860a083018861128b565b6020838203818501528188548084528284019150828160051b8501018a6000528360002060005b838110156114aa577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0878403018552611498838361128b565b9486019492506001918201910161145f565b505086810360408801526114be818b61128b565b94505050505084606084015282810360808401526114dc8185611220565b98975050505050505050565b606081528435606082015260208501356080820152604085013560a0820152606085013560c0820152608085013560e0820152600061152960a08701610b21565b61010073ffffffffffffffffffffffffffffffffffffffff82168185015261155460c08901896116e0565b92508161012086015261156c61016086018483611188565b9250505061157d60e0880188611747565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0858403016101408601526115b38382846111d7565b9250505085602084015282810360408401526115d08185876111d7565b979650505050505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261161057600080fd5b83018035915067ffffffffffffffff82111561162b57600080fd5b6020019150600581901b360382131561044057600080fd5b604051610100810167ffffffffffffffff811182821017156116675761166761188c565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156116b4576116b461188c565b604052919050565b600067ffffffffffffffff8211156116d6576116d661188c565b5060051b60200190565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261171557600080fd5b830160208101925035905067ffffffffffffffff81111561173557600080fd5b8060051b360383131561044057600080fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261177c57600080fd5b830160208101925035905067ffffffffffffffff81111561179c57600080fd5b80360383131561044057600080fd5b600082198211156117be576117be61182e565b500190565b6000828210156117d5576117d561182e565b500390565b600181811c908216806117ee57607f821691505b60208210811415611828577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var SimpleLogUpkeepCounterABI = SimpleLogUpkeepCounterMetaData.ABI

var SimpleLogUpkeepCounterBin = SimpleLogUpkeepCounterMetaData.Bin

func DeploySimpleLogUpkeepCounter(auth *bind.TransactOpts, backend bind.ContractBackend, _isStreamsLookup bool) (common.Address, *types.Transaction, *SimpleLogUpkeepCounter, error) {
	parsed, err := SimpleLogUpkeepCounterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SimpleLogUpkeepCounterBin), backend, _isStreamsLookup)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SimpleLogUpkeepCounter{address: address, abi: *parsed, SimpleLogUpkeepCounterCaller: SimpleLogUpkeepCounterCaller{contract: contract}, SimpleLogUpkeepCounterTransactor: SimpleLogUpkeepCounterTransactor{contract: contract}, SimpleLogUpkeepCounterFilterer: SimpleLogUpkeepCounterFilterer{contract: contract}}, nil
}

type SimpleLogUpkeepCounter struct {
	address common.Address
	abi     abi.ABI
	SimpleLogUpkeepCounterCaller
	SimpleLogUpkeepCounterTransactor
	SimpleLogUpkeepCounterFilterer
}

type SimpleLogUpkeepCounterCaller struct {
	contract *bind.BoundContract
}

type SimpleLogUpkeepCounterTransactor struct {
	contract *bind.BoundContract
}

type SimpleLogUpkeepCounterFilterer struct {
	contract *bind.BoundContract
}

type SimpleLogUpkeepCounterSession struct {
	Contract     *SimpleLogUpkeepCounter
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type SimpleLogUpkeepCounterCallerSession struct {
	Contract *SimpleLogUpkeepCounterCaller
	CallOpts bind.CallOpts
}

type SimpleLogUpkeepCounterTransactorSession struct {
	Contract     *SimpleLogUpkeepCounterTransactor
	TransactOpts bind.TransactOpts
}

type SimpleLogUpkeepCounterRaw struct {
	Contract *SimpleLogUpkeepCounter
}

type SimpleLogUpkeepCounterCallerRaw struct {
	Contract *SimpleLogUpkeepCounterCaller
}

type SimpleLogUpkeepCounterTransactorRaw struct {
	Contract *SimpleLogUpkeepCounterTransactor
}

func NewSimpleLogUpkeepCounter(address common.Address, backend bind.ContractBackend) (*SimpleLogUpkeepCounter, error) {
	abi, err := abi.JSON(strings.NewReader(SimpleLogUpkeepCounterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindSimpleLogUpkeepCounter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SimpleLogUpkeepCounter{address: address, abi: abi, SimpleLogUpkeepCounterCaller: SimpleLogUpkeepCounterCaller{contract: contract}, SimpleLogUpkeepCounterTransactor: SimpleLogUpkeepCounterTransactor{contract: contract}, SimpleLogUpkeepCounterFilterer: SimpleLogUpkeepCounterFilterer{contract: contract}}, nil
}

func NewSimpleLogUpkeepCounterCaller(address common.Address, caller bind.ContractCaller) (*SimpleLogUpkeepCounterCaller, error) {
	contract, err := bindSimpleLogUpkeepCounter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleLogUpkeepCounterCaller{contract: contract}, nil
}

func NewSimpleLogUpkeepCounterTransactor(address common.Address, transactor bind.ContractTransactor) (*SimpleLogUpkeepCounterTransactor, error) {
	contract, err := bindSimpleLogUpkeepCounter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleLogUpkeepCounterTransactor{contract: contract}, nil
}

func NewSimpleLogUpkeepCounterFilterer(address common.Address, filterer bind.ContractFilterer) (*SimpleLogUpkeepCounterFilterer, error) {
	contract, err := bindSimpleLogUpkeepCounter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimpleLogUpkeepCounterFilterer{contract: contract}, nil
}

func bindSimpleLogUpkeepCounter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SimpleLogUpkeepCounterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleLogUpkeepCounter.Contract.SimpleLogUpkeepCounterCaller.contract.Call(opts, result, method, params...)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.SimpleLogUpkeepCounterTransactor.contract.Transfer(opts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.SimpleLogUpkeepCounterTransactor.contract.Transact(opts, method, params...)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleLogUpkeepCounter.Contract.contract.Call(opts, result, method, params...)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.contract.Transfer(opts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.contract.Transact(opts, method, params...)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (bool, []byte, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "checkCallback", values, extraData)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) CheckCallback(values [][]byte, extraData []byte) (bool, []byte, error) {
	return _SimpleLogUpkeepCounter.Contract.CheckCallback(&_SimpleLogUpkeepCounter.CallOpts, values, extraData)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) CheckCallback(values [][]byte, extraData []byte) (bool, []byte, error) {
	return _SimpleLogUpkeepCounter.Contract.CheckCallback(&_SimpleLogUpkeepCounter.CallOpts, values, extraData)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) CheckErrorHandler(opts *bind.CallOpts, errCode *big.Int, extraData []byte) (CheckErrorHandler,

	error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "checkErrorHandler", errCode, extraData)

	outstruct := new(CheckErrorHandler)
	if err != nil {
		return *outstruct, err
	}

	outstruct.UpkeepNeeded = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.PerformData = *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return *outstruct, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) CheckErrorHandler(errCode *big.Int, extraData []byte) (CheckErrorHandler,

	error) {
	return _SimpleLogUpkeepCounter.Contract.CheckErrorHandler(&_SimpleLogUpkeepCounter.CallOpts, errCode, extraData)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) CheckErrorHandler(errCode *big.Int, extraData []byte) (CheckErrorHandler,

	error) {
	return _SimpleLogUpkeepCounter.Contract.CheckErrorHandler(&_SimpleLogUpkeepCounter.CallOpts, errCode, extraData)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) CheckLog(opts *bind.CallOpts, log Log, checkData []byte) (bool, []byte, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "checkLog", log, checkData)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) CheckLog(log Log, checkData []byte) (bool, []byte, error) {
	return _SimpleLogUpkeepCounter.Contract.CheckLog(&_SimpleLogUpkeepCounter.CallOpts, log, checkData)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) CheckLog(log Log, checkData []byte) (bool, []byte, error) {
	return _SimpleLogUpkeepCounter.Contract.CheckLog(&_SimpleLogUpkeepCounter.CallOpts, log, checkData)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) Counter() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.Counter(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) Counter() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.Counter(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) DummyMap(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "dummyMap", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) DummyMap(arg0 [32]byte) (bool, error) {
	return _SimpleLogUpkeepCounter.Contract.DummyMap(&_SimpleLogUpkeepCounter.CallOpts, arg0)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) DummyMap(arg0 [32]byte) (bool, error) {
	return _SimpleLogUpkeepCounter.Contract.DummyMap(&_SimpleLogUpkeepCounter.CallOpts, arg0)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) FeedParamKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "feedParamKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) FeedParamKey() (string, error) {
	return _SimpleLogUpkeepCounter.Contract.FeedParamKey(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) FeedParamKey() (string, error) {
	return _SimpleLogUpkeepCounter.Contract.FeedParamKey(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) FeedsHex(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "feedsHex", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) FeedsHex(arg0 *big.Int) (string, error) {
	return _SimpleLogUpkeepCounter.Contract.FeedsHex(&_SimpleLogUpkeepCounter.CallOpts, arg0)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) FeedsHex(arg0 *big.Int) (string, error) {
	return _SimpleLogUpkeepCounter.Contract.FeedsHex(&_SimpleLogUpkeepCounter.CallOpts, arg0)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) InitialBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "initialBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) InitialBlock() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.InitialBlock(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) InitialBlock() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.InitialBlock(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) IsStreamsLookup(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "isStreamsLookup")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) IsStreamsLookup() (bool, error) {
	return _SimpleLogUpkeepCounter.Contract.IsStreamsLookup(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) IsStreamsLookup() (bool, error) {
	return _SimpleLogUpkeepCounter.Contract.IsStreamsLookup(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) LastBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "lastBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) LastBlock() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.LastBlock(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) LastBlock() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.LastBlock(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "previousPerformBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) PreviousPerformBlock() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.PreviousPerformBlock(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) PreviousPerformBlock() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.PreviousPerformBlock(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) ShouldRetryOnError(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "shouldRetryOnError")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) ShouldRetryOnError() (bool, error) {
	return _SimpleLogUpkeepCounter.Contract.ShouldRetryOnError(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) ShouldRetryOnError() (bool, error) {
	return _SimpleLogUpkeepCounter.Contract.ShouldRetryOnError(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) TimeParamKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "timeParamKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) TimeParamKey() (string, error) {
	return _SimpleLogUpkeepCounter.Contract.TimeParamKey(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) TimeParamKey() (string, error) {
	return _SimpleLogUpkeepCounter.Contract.TimeParamKey(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCaller) TimeToPerform(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SimpleLogUpkeepCounter.contract.Call(opts, &out, "timeToPerform")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) TimeToPerform() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.TimeToPerform(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterCallerSession) TimeToPerform() (*big.Int, error) {
	return _SimpleLogUpkeepCounter.Contract.TimeToPerform(&_SimpleLogUpkeepCounter.CallOpts)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactor) CheckDataConfig(opts *bind.TransactOpts, arg0 CheckData) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.contract.Transact(opts, "_checkDataConfig", arg0)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) CheckDataConfig(arg0 CheckData) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.CheckDataConfig(&_SimpleLogUpkeepCounter.TransactOpts, arg0)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactorSession) CheckDataConfig(arg0 CheckData) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.CheckDataConfig(&_SimpleLogUpkeepCounter.TransactOpts, arg0)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.contract.Transact(opts, "performUpkeep", performData)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.PerformUpkeep(&_SimpleLogUpkeepCounter.TransactOpts, performData)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.PerformUpkeep(&_SimpleLogUpkeepCounter.TransactOpts, performData)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactor) SetFeedParamKey(opts *bind.TransactOpts, feedParam string) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.contract.Transact(opts, "setFeedParamKey", feedParam)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) SetFeedParamKey(feedParam string) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.SetFeedParamKey(&_SimpleLogUpkeepCounter.TransactOpts, feedParam)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactorSession) SetFeedParamKey(feedParam string) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.SetFeedParamKey(&_SimpleLogUpkeepCounter.TransactOpts, feedParam)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactor) SetFeedsHex(opts *bind.TransactOpts, newFeeds []string) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.contract.Transact(opts, "setFeedsHex", newFeeds)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) SetFeedsHex(newFeeds []string) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.SetFeedsHex(&_SimpleLogUpkeepCounter.TransactOpts, newFeeds)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactorSession) SetFeedsHex(newFeeds []string) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.SetFeedsHex(&_SimpleLogUpkeepCounter.TransactOpts, newFeeds)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactor) SetShouldRetryOnErrorBool(opts *bind.TransactOpts, value bool) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.contract.Transact(opts, "setShouldRetryOnErrorBool", value)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) SetShouldRetryOnErrorBool(value bool) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.SetShouldRetryOnErrorBool(&_SimpleLogUpkeepCounter.TransactOpts, value)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactorSession) SetShouldRetryOnErrorBool(value bool) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.SetShouldRetryOnErrorBool(&_SimpleLogUpkeepCounter.TransactOpts, value)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactor) SetTimeParamKey(opts *bind.TransactOpts, timeParam string) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.contract.Transact(opts, "setTimeParamKey", timeParam)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterSession) SetTimeParamKey(timeParam string) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.SetTimeParamKey(&_SimpleLogUpkeepCounter.TransactOpts, timeParam)
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterTransactorSession) SetTimeParamKey(timeParam string) (*types.Transaction, error) {
	return _SimpleLogUpkeepCounter.Contract.SetTimeParamKey(&_SimpleLogUpkeepCounter.TransactOpts, timeParam)
}

type SimpleLogUpkeepCounterPerformingUpkeepIterator struct {
	Event *SimpleLogUpkeepCounterPerformingUpkeep

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SimpleLogUpkeepCounterPerformingUpkeepIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleLogUpkeepCounterPerformingUpkeep)
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
		it.Event = new(SimpleLogUpkeepCounterPerformingUpkeep)
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

func (it *SimpleLogUpkeepCounterPerformingUpkeepIterator) Error() error {
	return it.fail
}

func (it *SimpleLogUpkeepCounterPerformingUpkeepIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SimpleLogUpkeepCounterPerformingUpkeep struct {
	From          common.Address
	InitialBlock  *big.Int
	LastBlock     *big.Int
	PreviousBlock *big.Int
	Counter       *big.Int
	TimeToPerform *big.Int
	IsRecovered   bool
	Raw           types.Log
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterFilterer) FilterPerformingUpkeep(opts *bind.FilterOpts, from []common.Address) (*SimpleLogUpkeepCounterPerformingUpkeepIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _SimpleLogUpkeepCounter.contract.FilterLogs(opts, "PerformingUpkeep", fromRule)
	if err != nil {
		return nil, err
	}
	return &SimpleLogUpkeepCounterPerformingUpkeepIterator{contract: _SimpleLogUpkeepCounter.contract, event: "PerformingUpkeep", logs: logs, sub: sub}, nil
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterFilterer) WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *SimpleLogUpkeepCounterPerformingUpkeep, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _SimpleLogUpkeepCounter.contract.WatchLogs(opts, "PerformingUpkeep", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(SimpleLogUpkeepCounterPerformingUpkeep)
				if err := _SimpleLogUpkeepCounter.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
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

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounterFilterer) ParsePerformingUpkeep(log types.Log) (*SimpleLogUpkeepCounterPerformingUpkeep, error) {
	event := new(SimpleLogUpkeepCounterPerformingUpkeep)
	if err := _SimpleLogUpkeepCounter.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CheckErrorHandler struct {
	UpkeepNeeded bool
	PerformData  []byte
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounter) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _SimpleLogUpkeepCounter.abi.Events["PerformingUpkeep"].ID:
		return _SimpleLogUpkeepCounter.ParsePerformingUpkeep(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (SimpleLogUpkeepCounterPerformingUpkeep) Topic() common.Hash {
	return common.HexToHash("0x29eff4cb37911c3ea85db4630638cc5474fdd0631ec42215aef1d7ec96c8e63d")
}

func (_SimpleLogUpkeepCounter *SimpleLogUpkeepCounter) Address() common.Address {
	return _SimpleLogUpkeepCounter.address
}

type SimpleLogUpkeepCounterInterface interface {
	CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (bool, []byte, error)

	CheckErrorHandler(opts *bind.CallOpts, errCode *big.Int, extraData []byte) (CheckErrorHandler,

		error)

	CheckLog(opts *bind.CallOpts, log Log, checkData []byte) (bool, []byte, error)

	Counter(opts *bind.CallOpts) (*big.Int, error)

	DummyMap(opts *bind.CallOpts, arg0 [32]byte) (bool, error)

	FeedParamKey(opts *bind.CallOpts) (string, error)

	FeedsHex(opts *bind.CallOpts, arg0 *big.Int) (string, error)

	InitialBlock(opts *bind.CallOpts) (*big.Int, error)

	IsStreamsLookup(opts *bind.CallOpts) (bool, error)

	LastBlock(opts *bind.CallOpts) (*big.Int, error)

	PreviousPerformBlock(opts *bind.CallOpts) (*big.Int, error)

	ShouldRetryOnError(opts *bind.CallOpts) (bool, error)

	TimeParamKey(opts *bind.CallOpts) (string, error)

	TimeToPerform(opts *bind.CallOpts) (*big.Int, error)

	CheckDataConfig(opts *bind.TransactOpts, arg0 CheckData) (*types.Transaction, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	SetFeedParamKey(opts *bind.TransactOpts, feedParam string) (*types.Transaction, error)

	SetFeedsHex(opts *bind.TransactOpts, newFeeds []string) (*types.Transaction, error)

	SetShouldRetryOnErrorBool(opts *bind.TransactOpts, value bool) (*types.Transaction, error)

	SetTimeParamKey(opts *bind.TransactOpts, timeParam string) (*types.Transaction, error)

	FilterPerformingUpkeep(opts *bind.FilterOpts, from []common.Address) (*SimpleLogUpkeepCounterPerformingUpkeepIterator, error)

	WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *SimpleLogUpkeepCounterPerformingUpkeep, from []common.Address) (event.Subscription, error)

	ParsePerformingUpkeep(log types.Log) (*SimpleLogUpkeepCounterPerformingUpkeep, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
