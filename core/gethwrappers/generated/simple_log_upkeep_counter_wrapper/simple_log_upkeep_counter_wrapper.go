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
	Feeds             []string
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
	ABI: "[{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_isStreamsLookup\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"StreamsLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lastBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timeToPerform\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isRecovered\",\"type\":\"bool\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"checkBurnAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"performBurnAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"eventSig\",\"type\":\"bytes32\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"}],\"internalType\":\"structCheckData\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"_checkDataConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"errCode\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkErrorHandler\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"topics\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLog\",\"name\":\"log\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkLog\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"dummyMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isStreamsLookup\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"previousPerformBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParam\",\"type\":\"string\"}],\"name\":\"setFeedParamKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"setShouldRetryOnErrorBool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"timeParam\",\"type\":\"string\"}],\"name\":\"setTimeParamKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"shouldRetryOnError\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeToPerform\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60c060405260076080819052666665656449447360c81b60a090815262000028919081620000bd565b5060408051808201909152600980825268074696d657374616d760bc1b60209092019182526200005b91600891620000bd565b503480156200006957600080fd5b5060405162001b3238038062001b328339810160408190526200008c9162000163565b60006002819055436001556003819055600455600680549115156101000261ff0019909216919091179055620001cb565b828054620000cb906200018e565b90600052602060002090601f016020900481019282620000ef57600085556200013a565b82601f106200010a57805160ff19168380011785556200013a565b828001600101855582156200013a579182015b828111156200013a5782518255916020019190600101906200011d565b50620001489291506200014c565b5090565b5b808211156200014857600081556001016200014d565b6000602082840312156200017657600080fd5b815180151581146200018757600080fd5b9392505050565b600181811c90821680620001a357607f821691505b60208210811415620001c557634e487b7160e01b600052602260045260246000fd5b50919050565b61195780620001db6000396000f3fe608060405234801561001057600080fd5b50600436106101365760003560e01c8063601d5a71116100b2578063917d895f11610081578063afb28d1f11610066578063afb28d1f146102a7578063c6066f0d146102bc578063c98f10b0146102c557600080fd5b8063917d895f1461028b5780639525d5741461029457600080fd5b8063601d5a711461024357806361bc221a146102565780637145f11b1461025f578063806b984f1461028257600080fd5b806340691db41161010957806342eb3d92116100ee57806342eb3d921461020a5780634585e33b1461021d5780634b56a42e1461023057600080fd5b806340691db4146101e657806342b0fe9e146101f957600080fd5b80630fb172fb1461013b57806313fab5901461016557806323148cee146101ad5780632cb15864146101cf575b600080fd5b61014e61014936600461101d565b6102cd565b60405161015c9291906113ea565b60405180910390f35b6101ab610173366004610d4a565b6006805491151562010000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffff909216919091179055565b005b6006546101bf90610100900460ff1681565b604051901515815260200161015c565b6101d860035481565b60405190815260200161015c565b61014e6101f4366004610ea7565b6103d4565b6101ab610207366004610e04565b50565b6006546101bf9062010000900460ff1681565b6101ab61022b366004610d85565b61065a565b61014e61023e366004610c70565b6108c0565b6101ab610251366004610dc7565b610914565b6101d860045481565b6101bf61026d366004610d6c565b60006020819052908152604090205460ff1681565b6101d860015481565b6101d860025481565b6101ab6102a2366004610dc7565b61092b565b6102af61093e565b60405161015c9190611405565b6101d860055481565b6102af6109cc565b6040805160028082526060828101909352600092918391816020015b60608152602001906001900390816102e95750506040805160208101889052919250016040516020818303038152906040528160008151811061032e5761032e6118ca565b60200260200101819052508360405160200161034a9190611405565b6040516020818303038152906040528160018151811061036c5761036c6118ca565b60200260200101819052506000818560405160200161038c929190611356565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291905260065462010000900460ff169450925050505b9250929050565b600060608180806103e78688018861105a565b9350935050925060005a90506000610400600143611800565b4090506000851561046f575b855a6104189085611800565b101561046f57808015610439575060008281526020819052604090205460ff165b6040805160208101859052309181019190915290915060600160405160208183030381529060405280519060200120915061040c565b60408051600280825260608201909252600091816020015b606081526020019060019003908161048757905050604080516000602082015291925001604051602081830303815290604052816000815181106104cd576104cd6118ca565b602002602001018190525060006040516020016104f3919060ff91909116815260200190565b60405160208183030381529060405281600181518110610515576105156118ca565b602002602001018190525060008c438d8d60405160200161053994939291906114d7565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905290508661057760c08f018f6115d2565b6002818110610588576105886118ca565b90506020020135141561062157600654610100900460ff16156105ec57600786600842846040517ff055e4a20000000000000000000000000000000000000000000000000000000081526004016105e3959493929190611418565b60405180910390fd5b60018282604051602001610601929190611356565b604051602081830303815290604052995099505050505050505050610652565b60008282604051602001610636929190611356565b6040516020818303038152906040529950995050505050505050505b935093915050565b60035461066657436003555b436001908155600454610678916117e8565b600455600154600255600061068f82840184610c70565b9150506000806000838060200190518101906106ab9190610f19565b9250925092508260200151426106c19190611800565b600555600680547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556060830151821461072357600680547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790555b6000808280602001905181019061073a91906110b4565b50925092505060005a90506000610752600143611800565b4090506000838860c0015160028151811061076f5761076f6118ca565b6020026020010151146107de576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f496e76616c6964206576656e74207369676e617475726500000000000000000060448201526064016105e3565b8415610848575b845a6107f19085611800565b101561084857808015610812575060008281526020819052604090205460ff165b604080516020810185905230918101919091529091506060016040516020818303038152906040528051906020012091506107e5565b600354600154600254600454600554600654604080519687526020870195909552938501929092526060840152608083015260ff16151560a082015232907f29eff4cb37911c3ea85db4630638cc5474fdd0631ec42215aef1d7ec96c8e63d9060c00160405180910390a25050505050505050505050565b60006060600084846040516020016108d9929190611356565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529190526001969095509350505050565b80516109279060089060208401906109d9565b5050565b80516109279060079060208401906109d9565b6007805461094b90611847565b80601f016020809104026020016040519081016040528092919081815260200182805461097790611847565b80156109c45780601f10610999576101008083540402835291602001916109c4565b820191906000526020600020905b8154815290600101906020018083116109a757829003601f168201915b505050505081565b6008805461094b90611847565b8280546109e590611847565b90600052602060002090601f016020900481019282610a075760008555610a4d565b82601f10610a2057805160ff1916838001178555610a4d565b82800160010185558215610a4d579182015b82811115610a4d578251825591602001919060010190610a32565b50610a59929150610a5d565b5090565b5b80821115610a595760008155600101610a5e565b6000610a85610a80846116d7565b611664565b9050828152838383011115610a9957600080fd5b610aa7836020830184611817565b9392505050565b8051610ab981611928565b919050565b600082601f830112610acf57600080fd5b81516020610adf610a80836116b3565b80838252828201915082860187848660051b8901011115610aff57600080fd5b60005b85811015610b1e57815184529284019290840190600101610b02565b5090979650505050505050565b600082601f830112610b3c57600080fd5b81356020610b4c610a80836116b3565b80838252828201915082860187848660051b8901011115610b6c57600080fd5b6000805b86811015610baf57823567ffffffffffffffff811115610b8e578283fd5b610b9c8b88838d0101610bff565b8652509385019391850191600101610b70565b509198975050505050505050565b60008083601f840112610bcf57600080fd5b50813567ffffffffffffffff811115610be757600080fd5b6020830191508360208285010111156103cd57600080fd5b600082601f830112610c1057600080fd5b8135610c1e610a80826116d7565b818152846020838601011115610c3357600080fd5b816020850160208301376000918101602001919091529392505050565b600082601f830112610c6157600080fd5b610aa783835160208501610a72565b60008060408385031215610c8357600080fd5b823567ffffffffffffffff80821115610c9b57600080fd5b818501915085601f830112610caf57600080fd5b81356020610cbf610a80836116b3565b8083825282820191508286018a848660051b8901011115610cdf57600080fd5b60005b85811015610d1a57813587811115610cf957600080fd5b610d078d87838c0101610bff565b8552509284019290840190600101610ce2565b50909750505086013592505080821115610d3357600080fd5b50610d4085828601610bff565b9150509250929050565b600060208284031215610d5c57600080fd5b81358015158114610aa757600080fd5b600060208284031215610d7e57600080fd5b5035919050565b60008060208385031215610d9857600080fd5b823567ffffffffffffffff811115610daf57600080fd5b610dbb85828601610bbd565b90969095509350505050565b600060208284031215610dd957600080fd5b813567ffffffffffffffff811115610df057600080fd5b610dfc84828501610bff565b949350505050565b600060208284031215610e1657600080fd5b813567ffffffffffffffff80821115610e2e57600080fd5b9083019060808286031215610e4257600080fd5b604051608081018181108382111715610e5d57610e5d6118f9565b8060405250823581526020830135602082015260408301356040820152606083013582811115610e8c57600080fd5b610e9887828601610b2b565b60608301525095945050505050565b600080600060408486031215610ebc57600080fd5b833567ffffffffffffffff80821115610ed457600080fd5b908501906101008288031215610ee957600080fd5b90935060208501359080821115610eff57600080fd5b50610f0c86828701610bbd565b9497909650939450505050565b600080600060608486031215610f2e57600080fd5b835167ffffffffffffffff80821115610f4657600080fd5b908501906101008288031215610f5b57600080fd5b610f6361163a565b8251815260208301516020820152604083015160408201526060830151606082015260808301516080820152610f9b60a08401610aae565b60a082015260c083015182811115610fb257600080fd5b610fbe89828601610abe565b60c08301525060e083015182811115610fd657600080fd5b610fe289828601610c50565b60e08301525060208701516040880151919650945091508082111561100657600080fd5b5061101386828701610c50565b9150509250925092565b6000806040838503121561103057600080fd5b82359150602083013567ffffffffffffffff81111561104e57600080fd5b610d4085828601610bff565b6000806000806080858703121561107057600080fd5b843593506020850135925060408501359150606085013567ffffffffffffffff81111561109c57600080fd5b6110a887828801610b2b565b91505092959194509250565b600080600080608085870312156110ca57600080fd5b84519350602080860151935060408601519250606086015167ffffffffffffffff808211156110f857600080fd5b818801915088601f83011261110c57600080fd5b815161111a610a80826116b3565b8082825285820191508585018c878560051b880101111561113a57600080fd5b60005b848110156111895781518681111561115457600080fd5b8701603f81018f1361116557600080fd5b6111768f8a83015160408401610a72565b855250928701929087019060010161113d565b505080965050505050505092959194509250565b81835260007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8311156111cf57600080fd5b8260051b8083602087013760009401602001938452509192915050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b6000815180845261124d816020860160208601611817565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8054600090600181811c908083168061129957607f831692505b60208084108214156112d4577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b838852602088018280156112ef576001811461131e57611349565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00871682528282019750611349565b60008981526020902060005b878110156113435781548482015290860190840161132a565b83019850505b5050505050505092915050565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b838110156113cb577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa08887030185526113b9868351611235565b9550938201939082019060010161137f565b5050858403818701525050506113e18185611235565b95945050505050565b8215158152604060208201526000610dfc6040830184611235565b602081526000610aa76020830184611235565b60a08152600061142b60a083018861127f565b6020838203818501528188518084528284019150828160051b850101838b0160005b83811015611499577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0878403018552611487838351611235565b9486019492509085019060010161144d565b505086810360408801526114ad818b61127f565b94505050505084606084015282810360808401526114cb8185611235565b98975050505050505050565b606081528435606082015260208501356080820152604085013560a0820152606085013560c0820152608085013560e0820152600060a086013561151a81611928565b61010061153e8185018373ffffffffffffffffffffffffffffffffffffffff169052565b61154b60c089018961171d565b9250816101208601526115636101608601848361119d565b9250505061157460e0880188611784565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0858403016101408601526115aa8382846111ec565b9250505085602084015282810360408401526115c78185876111ec565b979650505050505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261160757600080fd5b83018035915067ffffffffffffffff82111561162257600080fd5b6020019150600581901b36038213156103cd57600080fd5b604051610100810167ffffffffffffffff8111828210171561165e5761165e6118f9565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156116ab576116ab6118f9565b604052919050565b600067ffffffffffffffff8211156116cd576116cd6118f9565b5060051b60200190565b600067ffffffffffffffff8211156116f1576116f16118f9565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261175257600080fd5b830160208101925035905067ffffffffffffffff81111561177257600080fd5b8060051b36038313156103cd57600080fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126117b957600080fd5b830160208101925035905067ffffffffffffffff8111156117d957600080fd5b8036038313156103cd57600080fd5b600082198211156117fb576117fb61189b565b500190565b6000828210156118125761181261189b565b500390565b60005b8381101561183257818101518382015260200161181a565b83811115611841576000848401525b50505050565b600181811c9082168061185b57607f821691505b60208210811415611895577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff8116811461020757600080fdfea164736f6c6343000806000a",
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

	SetShouldRetryOnErrorBool(opts *bind.TransactOpts, value bool) (*types.Transaction, error)

	SetTimeParamKey(opts *bind.TransactOpts, timeParam string) (*types.Transaction, error)

	FilterPerformingUpkeep(opts *bind.FilterOpts, from []common.Address) (*SimpleLogUpkeepCounterPerformingUpkeepIterator, error)

	WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *SimpleLogUpkeepCounterPerformingUpkeep, from []common.Address) (event.Subscription, error)

	ParsePerformingUpkeep(log types.Log) (*SimpleLogUpkeepCounterPerformingUpkeep, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
