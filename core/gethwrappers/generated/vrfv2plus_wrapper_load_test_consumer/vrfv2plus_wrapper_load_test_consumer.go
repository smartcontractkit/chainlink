// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrfv2plus_wrapper_load_test_consumer

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

var VRFV2PlusWrapperLoadTestConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfV2PlusWrapper\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyVRFWrapperCanFulfill\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"}],\"name\":\"WrappedRequestFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"paid\",\"type\":\"uint256\"}],\"name\":\"WrapperRequestMade\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"offset\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"quantity\",\"type\":\"uint256\"}],\"name\":\"getRequestBlockTimes\",\"outputs\":[{\"internalType\":\"uint32[]\",\"name\":\"\",\"type\":\"uint32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"}],\"name\":\"getRequestStatus\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"paid\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"requestTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestBlockNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentBlockNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_vrfV2PlusWrapper\",\"outputs\":[{\"internalType\":\"contractIVRFV2PlusWrapper\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"_requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"_numWords\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"_requestCount\",\"type\":\"uint16\"}],\"name\":\"makeRequests\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"_requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"_numWords\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"_requestCount\",\"type\":\"uint16\"}],\"name\":\"makeRequestsNative\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"_randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_averageFulfillmentInMillions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fastestFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastRequestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requestBlockTimes\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_requestCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requests\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"paid\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"requestTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestBlockNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fulfilmentBlockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"native\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_responseCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_slowestFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawNative\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60c0604052600060045560006005556103e76006553480156200002157600080fd5b506040516200211738038062002117833981016040819052620000449162000225565b33806000836000819050806001600160a01b0316631c4695f46040518163ffffffff1660e01b8152600401602060405180830381600087803b1580156200008a57600080fd5b505af11580156200009f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620000c5919062000225565b6001600160601b0319606091821b811660805291901b1660a052506001600160a01b0382166200013c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200016f576200016f8162000179565b5050505062000257565b6001600160a01b038116331415620001d45760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000133565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156200023857600080fd5b81516001600160a01b03811681146200025057600080fd5b9392505050565b60805160601c60a05160601c611e5c620002bb6000396000818161033a015281816105f6015281816112f3015281816113ad01528181611466015281816115f801526116760152600081816104940152818161079d01526113710152611e5c6000f3fe6080604052600436106101795760003560e01c8063958cccb7116100cb578063d826f88f1161007f578063e76d516811610059578063e76d516814610485578063f1765962146104b8578063f2fde38b146104d857600080fd5b8063d826f88f14610427578063d8a4676f1461043c578063dc1670db1461046f57600080fd5b8063a168fa89116100b0578063a168fa891461035c578063afacbf9c146103f1578063b1e217491461041157600080fd5b8063958cccb7146102f35780639ed0868d1461032857600080fd5b8063737144bc1161012d5780637a8042bd116101075780637a8042bd1461026757806384276d81146102875780638da5cb5b146102a757600080fd5b8063737144bc1461022657806374dba1241461023c57806379ba50971461025257600080fd5b80631757f11c1161015e5780631757f11c146101d85780631fe543e3146101ee578063557d2e921461021057600080fd5b80630b2634861461018557806312065fe0146101bb57600080fd5b3661018057005b600080fd5b34801561019157600080fd5b506101a56101a0366004611a6f565b6104f8565b6040516101b29190611bc9565b60405180910390f35b3480156101c757600080fd5b50475b6040519081526020016101b2565b3480156101e457600080fd5b506101ca60055481565b3480156101fa57600080fd5b5061020e610209366004611980565b6105f4565b005b34801561021c57600080fd5b506101ca60035481565b34801561023257600080fd5b506101ca60045481565b34801561024857600080fd5b506101ca60065481565b34801561025e57600080fd5b5061020e610696565b34801561027357600080fd5b5061020e61028236600461194e565b610793565b34801561029357600080fd5b5061020e6102a236600461194e565b61089f565b3480156102b357600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101b2565b3480156102ff57600080fd5b5061031361030e36600461194e565b610971565b60405163ffffffff90911681526020016101b2565b34801561033457600080fd5b506102ce7f000000000000000000000000000000000000000000000000000000000000000081565b34801561036857600080fd5b506103ba61037736600461194e565b600a602052600090815260409020805460018201546003830154600484015460058501546006860154600790960154949560ff9485169593949293919290911687565b604080519788529515156020880152948601939093526060850191909152608084015260a0830152151560c082015260e0016101b2565b3480156103fd57600080fd5b5061020e61040c366004611a91565b6109ab565b34801561041d57600080fd5b506101ca60075481565b34801561043357600080fd5b5061020e610b92565b34801561044857600080fd5b5061045c61045736600461194e565b610bbc565b6040516101b29796959493929190611c3c565b34801561047b57600080fd5b506101ca60025481565b34801561049157600080fd5b507f00000000000000000000000000000000000000000000000000000000000000006102ce565b3480156104c457600080fd5b5061020e6104d3366004611a91565b610d3f565b3480156104e457600080fd5b5061020e6104f33660046118ef565b610f1e565b606060006105068385611cc0565b60085490915081111561051857506008545b60006105248583611d50565b67ffffffffffffffff81111561053c5761053c611e20565b604051908082528060200260200182016040528015610565578160200160208202803683370190505b509050845b828110156105eb576008818154811061058557610585611df1565b6000918252602090912060088204015460079091166004026101000a900463ffffffff16826105b48884611d50565b815181106105c4576105c4611df1565b63ffffffff90921660209283029190910190910152806105e381611d89565b91505061056a565b50949350505050565b7f00000000000000000000000000000000000000000000000000000000000000003373ffffffffffffffffffffffffffffffffffffffff821614610687576040517f8ba9316e00000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff821660248201526044015b60405180910390fd5b6106918383610f32565b505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610717576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161067e565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b61079b611168565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb6107f660005473ffffffffffffffffffffffffffffffffffffffff1690565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff909116600482015260248101849052604401602060405180830381600087803b15801561086357600080fd5b505af1158015610877573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061089b919061192c565b5050565b6108a7611168565b6000805460405173ffffffffffffffffffffffffffffffffffffffff9091169083908381818185875af1925050503d8060008114610901576040519150601f19603f3d011682016040523d82523d6000602084013e610906565b606091505b505090508061089b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f77697468647261774e6174697665206661696c65640000000000000000000000604482015260640161067e565b6008818154811061098157600080fd5b9060005260206000209060089182820401919006600402915054906101000a900463ffffffff1681565b6109b3611168565b60005b8161ffff168161ffff161015610b8b5760006109e26040518060200160405280600015158152506111e9565b90506000806109f3888888866112a5565b600782905590925090506000610a0761150d565b604080516101008101825284815260006020808301828152845183815280830186528486019081524260608601526080850184905260a0850187905260c0850184905260e08501849052898452600a8352949092208351815591516001830180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169115159190911790559251805194955091939092610aaf926002850192910190611843565b5060608201516003828101919091556080830151600483015560a0830151600583015560c0830151600683015560e090920151600790910180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169115159190911790558054906000610b2283611d89565b9091555050600083815260096020526040908190208290555183907f5f56b4c20db9f5b294cbf6f681368de4a992a27e2de2ee702dcf2cbbfa791ec490610b6c9085815260200190565b60405180910390a2505050508080610b8390611d67565b9150506109b6565b5050505050565b6000600481905560058190556103e760065560038190556002819055610bba9060089061188e565b565b6000818152600a602052604081205481906060908290819081908190610c3e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f72657175657374206e6f7420666f756e64000000000000000000000000000000604482015260640161067e565b6000888152600a6020908152604080832081516101008101835281548152600182015460ff16151581850152600282018054845181870281018701865281815292959394860193830182828015610cb457602002820191906000526020600020905b815481526020019060010190808311610ca0575b50505050508152602001600382015481526020016004820154815260200160058201548152602001600682015481526020016007820160009054906101000a900460ff1615151515815250509050806000015181602001518260400151836060015184608001518560a001518660c00151975097509750975097509750975050919395979092949650565b610d47611168565b60005b8161ffff168161ffff161015610b8b576000610d766040518060200160405280600115158152506111e9565b9050600080610d87888888866115aa565b600782905590925090506000610d9b61150d565b604080516101008101825284815260006020808301828152845183815280830186528486019081524260608601526080850184905260a0850187905260c08501849052600160e086018190528a8552600a84529590932084518155905194810180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001695151595909517909455905180519495509193610e429260028501920190611843565b5060608201516003828101919091556080830151600483015560a0830151600583015560c0830151600683015560e090920151600790910180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169115159190911790558054906000610eb583611d89565b9091555050600083815260096020526040908190208290555183907f5f56b4c20db9f5b294cbf6f681368de4a992a27e2de2ee702dcf2cbbfa791ec490610eff9085815260200190565b60405180910390a2505050508080610f1690611d67565b915050610d4a565b610f26611168565b610f2f81611726565b50565b6000828152600a6020526040902054610fa7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f72657175657374206e6f7420666f756e64000000000000000000000000000000604482015260640161067e565b6000610fb161150d565b60008481526009602052604081205491925090610fce9083611d50565b90506000610fdf82620f4240611d13565b9050600554821115610ff15760058290555b600654821061100257600654611004565b815b6006556002546110145780611047565b600254611022906001611cc0565b816002546004546110339190611d13565b61103d9190611cc0565b6110479190611cd8565b6004556002805490600061105a83611d89565b90915550506000858152600a60209081526040909120600181810180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016909117905585516110b292600290920191870190611843565b506000858152600a6020526040908190204260048083019190915560068201869055600880546001810182557ff3f7a9fe364faab93b216da50a3214154f22a0a2b415b23a84c8169e8b636ee391810491909101805460079092169092026101000a63ffffffff81810219909216918716021790555490517f6c84e12b4c188e61f1b4727024a5cf05c025fa58467e5eedf763c0744c89da7b916111599188918891611c13565b60405180910390a15050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610bba576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161067e565b60607f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa8260405160240161122291511515815260200190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915292915050565b6040517f4306d35400000000000000000000000000000000000000000000000000000000815263ffffffff85166004820152600090819073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001690634306d3549060240160206040518083038186803b15801561133557600080fd5b505afa158015611349573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061136d9190611967565b90507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16634000aea07f000000000000000000000000000000000000000000000000000000000000000083898989896040516020016113e49493929190611c83565b6040516020818303038152906040526040518463ffffffff1660e01b815260040161141193929190611b8b565b602060405180830381600087803b15801561142b57600080fd5b505af115801561143f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611463919061192c565b507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663fc2a88c36040518163ffffffff1660e01b815260040160206040518083038186803b1580156114ca57600080fd5b505afa1580156114de573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115029190611967565b915094509492505050565b6000466115198161181c565b156115a357606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b15801561156557600080fd5b505afa158015611579573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061159d9190611967565b91505090565b4391505090565b6040517f4b16093500000000000000000000000000000000000000000000000000000000815263ffffffff85166004820152600090819073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001690634b1609359060240160206040518083038186803b15801561163a57600080fd5b505afa15801561164e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906116729190611967565b90507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16639cfc058e82888888886040518663ffffffff1660e01b81526004016116d49493929190611c83565b6020604051808303818588803b1580156116ed57600080fd5b505af1158015611701573d6000803e3d6000fd5b50505050506040513d601f19601f820116820180604052508101906115029190611967565b73ffffffffffffffffffffffffffffffffffffffff81163314156117a6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161067e565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600061a4b1821480611830575062066eed82145b8061183d575062066eee82145b92915050565b82805482825590600052602060002090810192821561187e579160200282015b8281111561187e578251825591602001919060010190611863565b5061188a9291506118af565b5090565b508054600082556007016008900490600052602060002090810190610f2f91905b5b8082111561188a57600081556001016118b0565b803561ffff811681146118d657600080fd5b919050565b803563ffffffff811681146118d657600080fd5b60006020828403121561190157600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461192557600080fd5b9392505050565b60006020828403121561193e57600080fd5b8151801515811461192557600080fd5b60006020828403121561196057600080fd5b5035919050565b60006020828403121561197957600080fd5b5051919050565b6000806040838503121561199357600080fd5b8235915060208084013567ffffffffffffffff808211156119b357600080fd5b818601915086601f8301126119c757600080fd5b8135818111156119d9576119d9611e20565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f83011681018181108582111715611a1c57611a1c611e20565b604052828152858101935084860182860187018b1015611a3b57600080fd5b600095505b83861015611a5e578035855260019590950194938601938601611a40565b508096505050505050509250929050565b60008060408385031215611a8257600080fd5b50508035926020909101359150565b60008060008060808587031215611aa757600080fd5b611ab0856118db565b9350611abe602086016118c4565b9250611acc604086016118db565b9150611ada606086016118c4565b905092959194509250565b600081518084526020808501945080840160005b83811015611b1557815187529582019590820190600101611af9565b509495945050505050565b6000815180845260005b81811015611b4657602081850181015186830182015201611b2a565b81811115611b58576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b73ffffffffffffffffffffffffffffffffffffffff84168152826020820152606060408201526000611bc06060830184611b20565b95945050505050565b6020808252825182820181905260009190848201906040850190845b81811015611c0757835163ffffffff1683529284019291840191600101611be5565b50909695505050505050565b838152606060208201526000611c2c6060830185611ae5565b9050826040830152949350505050565b878152861515602082015260e060408201526000611c5d60e0830188611ae5565b90508560608301528460808301528360a08301528260c083015298975050505050505050565b600063ffffffff808716835261ffff8616602084015280851660408401525060806060830152611cb66080830184611b20565b9695505050505050565b60008219821115611cd357611cd3611dc2565b500190565b600082611d0e577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615611d4b57611d4b611dc2565b500290565b600082821015611d6257611d62611dc2565b500390565b600061ffff80831681811415611d7f57611d7f611dc2565b6001019392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415611dbb57611dbb611dc2565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFV2PlusWrapperLoadTestConsumerABI = VRFV2PlusWrapperLoadTestConsumerMetaData.ABI

var VRFV2PlusWrapperLoadTestConsumerBin = VRFV2PlusWrapperLoadTestConsumerMetaData.Bin

func DeployVRFV2PlusWrapperLoadTestConsumer(auth *bind.TransactOpts, backend bind.ContractBackend, _vrfV2PlusWrapper common.Address) (common.Address, *types.Transaction, *VRFV2PlusWrapperLoadTestConsumer, error) {
	parsed, err := VRFV2PlusWrapperLoadTestConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2PlusWrapperLoadTestConsumerBin), backend, _vrfV2PlusWrapper)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2PlusWrapperLoadTestConsumer{address: address, abi: *parsed, VRFV2PlusWrapperLoadTestConsumerCaller: VRFV2PlusWrapperLoadTestConsumerCaller{contract: contract}, VRFV2PlusWrapperLoadTestConsumerTransactor: VRFV2PlusWrapperLoadTestConsumerTransactor{contract: contract}, VRFV2PlusWrapperLoadTestConsumerFilterer: VRFV2PlusWrapperLoadTestConsumerFilterer{contract: contract}}, nil
}

type VRFV2PlusWrapperLoadTestConsumer struct {
	address common.Address
	abi     abi.ABI
	VRFV2PlusWrapperLoadTestConsumerCaller
	VRFV2PlusWrapperLoadTestConsumerTransactor
	VRFV2PlusWrapperLoadTestConsumerFilterer
}

type VRFV2PlusWrapperLoadTestConsumerCaller struct {
	contract *bind.BoundContract
}

type VRFV2PlusWrapperLoadTestConsumerTransactor struct {
	contract *bind.BoundContract
}

type VRFV2PlusWrapperLoadTestConsumerFilterer struct {
	contract *bind.BoundContract
}

type VRFV2PlusWrapperLoadTestConsumerSession struct {
	Contract     *VRFV2PlusWrapperLoadTestConsumer
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2PlusWrapperLoadTestConsumerCallerSession struct {
	Contract *VRFV2PlusWrapperLoadTestConsumerCaller
	CallOpts bind.CallOpts
}

type VRFV2PlusWrapperLoadTestConsumerTransactorSession struct {
	Contract     *VRFV2PlusWrapperLoadTestConsumerTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2PlusWrapperLoadTestConsumerRaw struct {
	Contract *VRFV2PlusWrapperLoadTestConsumer
}

type VRFV2PlusWrapperLoadTestConsumerCallerRaw struct {
	Contract *VRFV2PlusWrapperLoadTestConsumerCaller
}

type VRFV2PlusWrapperLoadTestConsumerTransactorRaw struct {
	Contract *VRFV2PlusWrapperLoadTestConsumerTransactor
}

func NewVRFV2PlusWrapperLoadTestConsumer(address common.Address, backend bind.ContractBackend) (*VRFV2PlusWrapperLoadTestConsumer, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2PlusWrapperLoadTestConsumerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2PlusWrapperLoadTestConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusWrapperLoadTestConsumer{address: address, abi: abi, VRFV2PlusWrapperLoadTestConsumerCaller: VRFV2PlusWrapperLoadTestConsumerCaller{contract: contract}, VRFV2PlusWrapperLoadTestConsumerTransactor: VRFV2PlusWrapperLoadTestConsumerTransactor{contract: contract}, VRFV2PlusWrapperLoadTestConsumerFilterer: VRFV2PlusWrapperLoadTestConsumerFilterer{contract: contract}}, nil
}

func NewVRFV2PlusWrapperLoadTestConsumerCaller(address common.Address, caller bind.ContractCaller) (*VRFV2PlusWrapperLoadTestConsumerCaller, error) {
	contract, err := bindVRFV2PlusWrapperLoadTestConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusWrapperLoadTestConsumerCaller{contract: contract}, nil
}

func NewVRFV2PlusWrapperLoadTestConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2PlusWrapperLoadTestConsumerTransactor, error) {
	contract, err := bindVRFV2PlusWrapperLoadTestConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusWrapperLoadTestConsumerTransactor{contract: contract}, nil
}

func NewVRFV2PlusWrapperLoadTestConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2PlusWrapperLoadTestConsumerFilterer, error) {
	contract, err := bindVRFV2PlusWrapperLoadTestConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusWrapperLoadTestConsumerFilterer{contract: contract}, nil
}

func bindVRFV2PlusWrapperLoadTestConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFV2PlusWrapperLoadTestConsumerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.VRFV2PlusWrapperLoadTestConsumerCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.VRFV2PlusWrapperLoadTestConsumerTransactor.contract.Transfer(opts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.VRFV2PlusWrapperLoadTestConsumerTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.contract.Transfer(opts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCaller) GetBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusWrapperLoadTestConsumer.contract.Call(opts, &out, "getBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerSession) GetBalance() (*big.Int, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.GetBalance(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCallerSession) GetBalance() (*big.Int, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.GetBalance(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCaller) GetLinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusWrapperLoadTestConsumer.contract.Call(opts, &out, "getLinkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerSession) GetLinkToken() (common.Address, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.GetLinkToken(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCallerSession) GetLinkToken() (common.Address, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.GetLinkToken(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCaller) GetRequestBlockTimes(opts *bind.CallOpts, offset *big.Int, quantity *big.Int) ([]uint32, error) {
	var out []interface{}
	err := _VRFV2PlusWrapperLoadTestConsumer.contract.Call(opts, &out, "getRequestBlockTimes", offset, quantity)

	if err != nil {
		return *new([]uint32), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint32)).(*[]uint32)

	return out0, err

}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerSession) GetRequestBlockTimes(offset *big.Int, quantity *big.Int) ([]uint32, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.GetRequestBlockTimes(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts, offset, quantity)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCallerSession) GetRequestBlockTimes(offset *big.Int, quantity *big.Int) ([]uint32, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.GetRequestBlockTimes(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts, offset, quantity)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCaller) GetRequestStatus(opts *bind.CallOpts, _requestId *big.Int) (GetRequestStatus,

	error) {
	var out []interface{}
	err := _VRFV2PlusWrapperLoadTestConsumer.contract.Call(opts, &out, "getRequestStatus", _requestId)

	outstruct := new(GetRequestStatus)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Paid = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Fulfilled = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.RandomWords = *abi.ConvertType(out[2], new([]*big.Int)).(*[]*big.Int)
	outstruct.RequestTimestamp = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.FulfilmentTimestamp = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.RequestBlockNumber = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.FulfilmentBlockNumber = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerSession) GetRequestStatus(_requestId *big.Int) (GetRequestStatus,

	error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.GetRequestStatus(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts, _requestId)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCallerSession) GetRequestStatus(_requestId *big.Int) (GetRequestStatus,

	error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.GetRequestStatus(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts, _requestId)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCaller) IVrfV2PlusWrapper(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusWrapperLoadTestConsumer.contract.Call(opts, &out, "i_vrfV2PlusWrapper")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerSession) IVrfV2PlusWrapper() (common.Address, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.IVrfV2PlusWrapper(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCallerSession) IVrfV2PlusWrapper() (common.Address, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.IVrfV2PlusWrapper(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusWrapperLoadTestConsumer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerSession) Owner() (common.Address, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.Owner(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCallerSession) Owner() (common.Address, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.Owner(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCaller) SAverageFulfillmentInMillions(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusWrapperLoadTestConsumer.contract.Call(opts, &out, "s_averageFulfillmentInMillions")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerSession) SAverageFulfillmentInMillions() (*big.Int, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.SAverageFulfillmentInMillions(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCallerSession) SAverageFulfillmentInMillions() (*big.Int, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.SAverageFulfillmentInMillions(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCaller) SFastestFulfillment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusWrapperLoadTestConsumer.contract.Call(opts, &out, "s_fastestFulfillment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerSession) SFastestFulfillment() (*big.Int, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.SFastestFulfillment(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCallerSession) SFastestFulfillment() (*big.Int, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.SFastestFulfillment(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCaller) SLastRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusWrapperLoadTestConsumer.contract.Call(opts, &out, "s_lastRequestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerSession) SLastRequestId() (*big.Int, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.SLastRequestId(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCallerSession) SLastRequestId() (*big.Int, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.SLastRequestId(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCaller) SRequestBlockTimes(opts *bind.CallOpts, arg0 *big.Int) (uint32, error) {
	var out []interface{}
	err := _VRFV2PlusWrapperLoadTestConsumer.contract.Call(opts, &out, "s_requestBlockTimes", arg0)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerSession) SRequestBlockTimes(arg0 *big.Int) (uint32, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.SRequestBlockTimes(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts, arg0)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCallerSession) SRequestBlockTimes(arg0 *big.Int) (uint32, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.SRequestBlockTimes(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts, arg0)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCaller) SRequestCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusWrapperLoadTestConsumer.contract.Call(opts, &out, "s_requestCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerSession) SRequestCount() (*big.Int, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.SRequestCount(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCallerSession) SRequestCount() (*big.Int, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.SRequestCount(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCaller) SRequests(opts *bind.CallOpts, arg0 *big.Int) (SRequests,

	error) {
	var out []interface{}
	err := _VRFV2PlusWrapperLoadTestConsumer.contract.Call(opts, &out, "s_requests", arg0)

	outstruct := new(SRequests)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Paid = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Fulfilled = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.RequestTimestamp = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.FulfilmentTimestamp = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.RequestBlockNumber = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.FulfilmentBlockNumber = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.Native = *abi.ConvertType(out[6], new(bool)).(*bool)

	return *outstruct, err

}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerSession) SRequests(arg0 *big.Int) (SRequests,

	error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.SRequests(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts, arg0)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCallerSession) SRequests(arg0 *big.Int) (SRequests,

	error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.SRequests(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts, arg0)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCaller) SResponseCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusWrapperLoadTestConsumer.contract.Call(opts, &out, "s_responseCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerSession) SResponseCount() (*big.Int, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.SResponseCount(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCallerSession) SResponseCount() (*big.Int, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.SResponseCount(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCaller) SSlowestFulfillment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusWrapperLoadTestConsumer.contract.Call(opts, &out, "s_slowestFulfillment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerSession) SSlowestFulfillment() (*big.Int, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.SSlowestFulfillment(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerCallerSession) SSlowestFulfillment() (*big.Int, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.SSlowestFulfillment(&_VRFV2PlusWrapperLoadTestConsumer.CallOpts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.contract.Transact(opts, "acceptOwnership")
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.AcceptOwnership(&_VRFV2PlusWrapperLoadTestConsumer.TransactOpts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.AcceptOwnership(&_VRFV2PlusWrapperLoadTestConsumer.TransactOpts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerTransactor) MakeRequests(opts *bind.TransactOpts, _callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32, _requestCount uint16) (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.contract.Transact(opts, "makeRequests", _callbackGasLimit, _requestConfirmations, _numWords, _requestCount)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerSession) MakeRequests(_callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32, _requestCount uint16) (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.MakeRequests(&_VRFV2PlusWrapperLoadTestConsumer.TransactOpts, _callbackGasLimit, _requestConfirmations, _numWords, _requestCount)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerTransactorSession) MakeRequests(_callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32, _requestCount uint16) (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.MakeRequests(&_VRFV2PlusWrapperLoadTestConsumer.TransactOpts, _callbackGasLimit, _requestConfirmations, _numWords, _requestCount)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerTransactor) MakeRequestsNative(opts *bind.TransactOpts, _callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32, _requestCount uint16) (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.contract.Transact(opts, "makeRequestsNative", _callbackGasLimit, _requestConfirmations, _numWords, _requestCount)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerSession) MakeRequestsNative(_callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32, _requestCount uint16) (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.MakeRequestsNative(&_VRFV2PlusWrapperLoadTestConsumer.TransactOpts, _callbackGasLimit, _requestConfirmations, _numWords, _requestCount)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerTransactorSession) MakeRequestsNative(_callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32, _requestCount uint16) (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.MakeRequestsNative(&_VRFV2PlusWrapperLoadTestConsumer.TransactOpts, _callbackGasLimit, _requestConfirmations, _numWords, _requestCount)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, _requestId *big.Int, _randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.contract.Transact(opts, "rawFulfillRandomWords", _requestId, _randomWords)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerSession) RawFulfillRandomWords(_requestId *big.Int, _randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.RawFulfillRandomWords(&_VRFV2PlusWrapperLoadTestConsumer.TransactOpts, _requestId, _randomWords)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerTransactorSession) RawFulfillRandomWords(_requestId *big.Int, _randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.RawFulfillRandomWords(&_VRFV2PlusWrapperLoadTestConsumer.TransactOpts, _requestId, _randomWords)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerTransactor) Reset(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.contract.Transact(opts, "reset")
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerSession) Reset() (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.Reset(&_VRFV2PlusWrapperLoadTestConsumer.TransactOpts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerTransactorSession) Reset() (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.Reset(&_VRFV2PlusWrapperLoadTestConsumer.TransactOpts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.TransferOwnership(&_VRFV2PlusWrapperLoadTestConsumer.TransactOpts, to)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.TransferOwnership(&_VRFV2PlusWrapperLoadTestConsumer.TransactOpts, to)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerTransactor) WithdrawLink(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.contract.Transact(opts, "withdrawLink", amount)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerSession) WithdrawLink(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.WithdrawLink(&_VRFV2PlusWrapperLoadTestConsumer.TransactOpts, amount)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerTransactorSession) WithdrawLink(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.WithdrawLink(&_VRFV2PlusWrapperLoadTestConsumer.TransactOpts, amount)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerTransactor) WithdrawNative(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.contract.Transact(opts, "withdrawNative", amount)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerSession) WithdrawNative(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.WithdrawNative(&_VRFV2PlusWrapperLoadTestConsumer.TransactOpts, amount)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerTransactorSession) WithdrawNative(amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.WithdrawNative(&_VRFV2PlusWrapperLoadTestConsumer.TransactOpts, amount)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.contract.RawTransact(opts, nil)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerSession) Receive() (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.Receive(&_VRFV2PlusWrapperLoadTestConsumer.TransactOpts)
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerTransactorSession) Receive() (*types.Transaction, error) {
	return _VRFV2PlusWrapperLoadTestConsumer.Contract.Receive(&_VRFV2PlusWrapperLoadTestConsumer.TransactOpts)
}

type VRFV2PlusWrapperLoadTestConsumerOwnershipTransferRequestedIterator struct {
	Event *VRFV2PlusWrapperLoadTestConsumerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusWrapperLoadTestConsumerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusWrapperLoadTestConsumerOwnershipTransferRequested)
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
		it.Event = new(VRFV2PlusWrapperLoadTestConsumerOwnershipTransferRequested)
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

func (it *VRFV2PlusWrapperLoadTestConsumerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusWrapperLoadTestConsumerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusWrapperLoadTestConsumerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusWrapperLoadTestConsumerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusWrapperLoadTestConsumer.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusWrapperLoadTestConsumerOwnershipTransferRequestedIterator{contract: _VRFV2PlusWrapperLoadTestConsumer.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2PlusWrapperLoadTestConsumerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusWrapperLoadTestConsumer.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusWrapperLoadTestConsumerOwnershipTransferRequested)
				if err := _VRFV2PlusWrapperLoadTestConsumer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFV2PlusWrapperLoadTestConsumerOwnershipTransferRequested, error) {
	event := new(VRFV2PlusWrapperLoadTestConsumerOwnershipTransferRequested)
	if err := _VRFV2PlusWrapperLoadTestConsumer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2PlusWrapperLoadTestConsumerOwnershipTransferredIterator struct {
	Event *VRFV2PlusWrapperLoadTestConsumerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusWrapperLoadTestConsumerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusWrapperLoadTestConsumerOwnershipTransferred)
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
		it.Event = new(VRFV2PlusWrapperLoadTestConsumerOwnershipTransferred)
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

func (it *VRFV2PlusWrapperLoadTestConsumerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusWrapperLoadTestConsumerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusWrapperLoadTestConsumerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusWrapperLoadTestConsumerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusWrapperLoadTestConsumer.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusWrapperLoadTestConsumerOwnershipTransferredIterator{contract: _VRFV2PlusWrapperLoadTestConsumer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2PlusWrapperLoadTestConsumerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusWrapperLoadTestConsumer.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusWrapperLoadTestConsumerOwnershipTransferred)
				if err := _VRFV2PlusWrapperLoadTestConsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerFilterer) ParseOwnershipTransferred(log types.Log) (*VRFV2PlusWrapperLoadTestConsumerOwnershipTransferred, error) {
	event := new(VRFV2PlusWrapperLoadTestConsumerOwnershipTransferred)
	if err := _VRFV2PlusWrapperLoadTestConsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2PlusWrapperLoadTestConsumerWrappedRequestFulfilledIterator struct {
	Event *VRFV2PlusWrapperLoadTestConsumerWrappedRequestFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusWrapperLoadTestConsumerWrappedRequestFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusWrapperLoadTestConsumerWrappedRequestFulfilled)
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
		it.Event = new(VRFV2PlusWrapperLoadTestConsumerWrappedRequestFulfilled)
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

func (it *VRFV2PlusWrapperLoadTestConsumerWrappedRequestFulfilledIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusWrapperLoadTestConsumerWrappedRequestFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusWrapperLoadTestConsumerWrappedRequestFulfilled struct {
	RequestId   *big.Int
	RandomWords []*big.Int
	Payment     *big.Int
	Raw         types.Log
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerFilterer) FilterWrappedRequestFulfilled(opts *bind.FilterOpts) (*VRFV2PlusWrapperLoadTestConsumerWrappedRequestFulfilledIterator, error) {

	logs, sub, err := _VRFV2PlusWrapperLoadTestConsumer.contract.FilterLogs(opts, "WrappedRequestFulfilled")
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusWrapperLoadTestConsumerWrappedRequestFulfilledIterator{contract: _VRFV2PlusWrapperLoadTestConsumer.contract, event: "WrappedRequestFulfilled", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerFilterer) WatchWrappedRequestFulfilled(opts *bind.WatchOpts, sink chan<- *VRFV2PlusWrapperLoadTestConsumerWrappedRequestFulfilled) (event.Subscription, error) {

	logs, sub, err := _VRFV2PlusWrapperLoadTestConsumer.contract.WatchLogs(opts, "WrappedRequestFulfilled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusWrapperLoadTestConsumerWrappedRequestFulfilled)
				if err := _VRFV2PlusWrapperLoadTestConsumer.contract.UnpackLog(event, "WrappedRequestFulfilled", log); err != nil {
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

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerFilterer) ParseWrappedRequestFulfilled(log types.Log) (*VRFV2PlusWrapperLoadTestConsumerWrappedRequestFulfilled, error) {
	event := new(VRFV2PlusWrapperLoadTestConsumerWrappedRequestFulfilled)
	if err := _VRFV2PlusWrapperLoadTestConsumer.contract.UnpackLog(event, "WrappedRequestFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2PlusWrapperLoadTestConsumerWrapperRequestMadeIterator struct {
	Event *VRFV2PlusWrapperLoadTestConsumerWrapperRequestMade

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusWrapperLoadTestConsumerWrapperRequestMadeIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusWrapperLoadTestConsumerWrapperRequestMade)
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
		it.Event = new(VRFV2PlusWrapperLoadTestConsumerWrapperRequestMade)
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

func (it *VRFV2PlusWrapperLoadTestConsumerWrapperRequestMadeIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusWrapperLoadTestConsumerWrapperRequestMadeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusWrapperLoadTestConsumerWrapperRequestMade struct {
	RequestId *big.Int
	Paid      *big.Int
	Raw       types.Log
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerFilterer) FilterWrapperRequestMade(opts *bind.FilterOpts, requestId []*big.Int) (*VRFV2PlusWrapperLoadTestConsumerWrapperRequestMadeIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _VRFV2PlusWrapperLoadTestConsumer.contract.FilterLogs(opts, "WrapperRequestMade", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusWrapperLoadTestConsumerWrapperRequestMadeIterator{contract: _VRFV2PlusWrapperLoadTestConsumer.contract, event: "WrapperRequestMade", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerFilterer) WatchWrapperRequestMade(opts *bind.WatchOpts, sink chan<- *VRFV2PlusWrapperLoadTestConsumerWrapperRequestMade, requestId []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _VRFV2PlusWrapperLoadTestConsumer.contract.WatchLogs(opts, "WrapperRequestMade", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusWrapperLoadTestConsumerWrapperRequestMade)
				if err := _VRFV2PlusWrapperLoadTestConsumer.contract.UnpackLog(event, "WrapperRequestMade", log); err != nil {
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

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumerFilterer) ParseWrapperRequestMade(log types.Log) (*VRFV2PlusWrapperLoadTestConsumerWrapperRequestMade, error) {
	event := new(VRFV2PlusWrapperLoadTestConsumerWrapperRequestMade)
	if err := _VRFV2PlusWrapperLoadTestConsumer.contract.UnpackLog(event, "WrapperRequestMade", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetRequestStatus struct {
	Paid                  *big.Int
	Fulfilled             bool
	RandomWords           []*big.Int
	RequestTimestamp      *big.Int
	FulfilmentTimestamp   *big.Int
	RequestBlockNumber    *big.Int
	FulfilmentBlockNumber *big.Int
}
type SRequests struct {
	Paid                  *big.Int
	Fulfilled             bool
	RequestTimestamp      *big.Int
	FulfilmentTimestamp   *big.Int
	RequestBlockNumber    *big.Int
	FulfilmentBlockNumber *big.Int
	Native                bool
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumer) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFV2PlusWrapperLoadTestConsumer.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFV2PlusWrapperLoadTestConsumer.ParseOwnershipTransferRequested(log)
	case _VRFV2PlusWrapperLoadTestConsumer.abi.Events["OwnershipTransferred"].ID:
		return _VRFV2PlusWrapperLoadTestConsumer.ParseOwnershipTransferred(log)
	case _VRFV2PlusWrapperLoadTestConsumer.abi.Events["WrappedRequestFulfilled"].ID:
		return _VRFV2PlusWrapperLoadTestConsumer.ParseWrappedRequestFulfilled(log)
	case _VRFV2PlusWrapperLoadTestConsumer.abi.Events["WrapperRequestMade"].ID:
		return _VRFV2PlusWrapperLoadTestConsumer.ParseWrapperRequestMade(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFV2PlusWrapperLoadTestConsumerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFV2PlusWrapperLoadTestConsumerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VRFV2PlusWrapperLoadTestConsumerWrappedRequestFulfilled) Topic() common.Hash {
	return common.HexToHash("0x6c84e12b4c188e61f1b4727024a5cf05c025fa58467e5eedf763c0744c89da7b")
}

func (VRFV2PlusWrapperLoadTestConsumerWrapperRequestMade) Topic() common.Hash {
	return common.HexToHash("0x5f56b4c20db9f5b294cbf6f681368de4a992a27e2de2ee702dcf2cbbfa791ec4")
}

func (_VRFV2PlusWrapperLoadTestConsumer *VRFV2PlusWrapperLoadTestConsumer) Address() common.Address {
	return _VRFV2PlusWrapperLoadTestConsumer.address
}

type VRFV2PlusWrapperLoadTestConsumerInterface interface {
	GetBalance(opts *bind.CallOpts) (*big.Int, error)

	GetLinkToken(opts *bind.CallOpts) (common.Address, error)

	GetRequestBlockTimes(opts *bind.CallOpts, offset *big.Int, quantity *big.Int) ([]uint32, error)

	GetRequestStatus(opts *bind.CallOpts, _requestId *big.Int) (GetRequestStatus,

		error)

	IVrfV2PlusWrapper(opts *bind.CallOpts) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SAverageFulfillmentInMillions(opts *bind.CallOpts) (*big.Int, error)

	SFastestFulfillment(opts *bind.CallOpts) (*big.Int, error)

	SLastRequestId(opts *bind.CallOpts) (*big.Int, error)

	SRequestBlockTimes(opts *bind.CallOpts, arg0 *big.Int) (uint32, error)

	SRequestCount(opts *bind.CallOpts) (*big.Int, error)

	SRequests(opts *bind.CallOpts, arg0 *big.Int) (SRequests,

		error)

	SResponseCount(opts *bind.CallOpts) (*big.Int, error)

	SSlowestFulfillment(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	MakeRequests(opts *bind.TransactOpts, _callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32, _requestCount uint16) (*types.Transaction, error)

	MakeRequestsNative(opts *bind.TransactOpts, _callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32, _requestCount uint16) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, _requestId *big.Int, _randomWords []*big.Int) (*types.Transaction, error)

	Reset(opts *bind.TransactOpts) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	WithdrawLink(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	WithdrawNative(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	Receive(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusWrapperLoadTestConsumerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2PlusWrapperLoadTestConsumerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFV2PlusWrapperLoadTestConsumerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusWrapperLoadTestConsumerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2PlusWrapperLoadTestConsumerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFV2PlusWrapperLoadTestConsumerOwnershipTransferred, error)

	FilterWrappedRequestFulfilled(opts *bind.FilterOpts) (*VRFV2PlusWrapperLoadTestConsumerWrappedRequestFulfilledIterator, error)

	WatchWrappedRequestFulfilled(opts *bind.WatchOpts, sink chan<- *VRFV2PlusWrapperLoadTestConsumerWrappedRequestFulfilled) (event.Subscription, error)

	ParseWrappedRequestFulfilled(log types.Log) (*VRFV2PlusWrapperLoadTestConsumerWrappedRequestFulfilled, error)

	FilterWrapperRequestMade(opts *bind.FilterOpts, requestId []*big.Int) (*VRFV2PlusWrapperLoadTestConsumerWrapperRequestMadeIterator, error)

	WatchWrapperRequestMade(opts *bind.WatchOpts, sink chan<- *VRFV2PlusWrapperLoadTestConsumerWrapperRequestMade, requestId []*big.Int) (event.Subscription, error)

	ParseWrapperRequestMade(log types.Log) (*VRFV2PlusWrapperLoadTestConsumerWrapperRequestMade, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
