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
	Bin: "0x60c0604052600060045560006005556103e76006553480156200002157600080fd5b506040516200204638038062002046833981016040819052620000449162000205565b33806000836000819050806001600160a01b0316631c4695f46040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200008d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620000b3919062000205565b6001600160a01b0390811660805290811660a052831690506200011d5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620001505762000150816200015a565b5050505062000237565b336001600160a01b03821603620001b45760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000114565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156200021857600080fd5b81516001600160a01b03811681146200023057600080fd5b9392505050565b60805160a051611db1620002956000396000818161033a015281816105f8015281816112de0152818161138901528181611433015281816115af015261161e0152600081816104940152818161079f015261134d0152611db16000f3fe6080604052600436106101795760003560e01c8063958cccb7116100cb578063d826f88f1161007f578063e76d516811610059578063e76d516814610485578063f1765962146104b8578063f2fde38b146104d857600080fd5b8063d826f88f14610427578063d8a4676f1461043c578063dc1670db1461046f57600080fd5b8063a168fa89116100b0578063a168fa891461035c578063afacbf9c146103f1578063b1e217491461041157600080fd5b8063958cccb7146102f35780639ed0868d1461032857600080fd5b8063737144bc1161012d5780637a8042bd116101075780637a8042bd1461026757806384276d81146102875780638da5cb5b146102a757600080fd5b8063737144bc1461022657806374dba1241461023c57806379ba50971461025257600080fd5b80631757f11c1161015e5780631757f11c146101d85780631fe543e3146101ee578063557d2e921461021057600080fd5b80630b2634861461018557806312065fe0146101bb57600080fd5b3661018057005b600080fd5b34801561019157600080fd5b506101a56101a0366004611858565b6104f8565b6040516101b2919061187a565b60405180910390f35b3480156101c757600080fd5b50475b6040519081526020016101b2565b3480156101e457600080fd5b506101ca60055481565b3480156101fa57600080fd5b5061020e6102093660046118f3565b6105f6565b005b34801561021c57600080fd5b506101ca60035481565b34801561023257600080fd5b506101ca60045481565b34801561024857600080fd5b506101ca60065481565b34801561025e57600080fd5b5061020e610698565b34801561027357600080fd5b5061020e6102823660046119db565b610795565b34801561029357600080fd5b5061020e6102a23660046119db565b610892565b3480156102b357600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101b2565b3480156102ff57600080fd5b5061031361030e3660046119db565b610964565b60405163ffffffff90911681526020016101b2565b34801561033457600080fd5b506102ce7f000000000000000000000000000000000000000000000000000000000000000081565b34801561036857600080fd5b506103ba6103773660046119db565b600a602052600090815260409020805460018201546003830154600484015460058501546006860154600790960154949560ff9485169593949293919290911687565b604080519788529515156020880152948601939093526060850191909152608084015260a0830152151560c082015260e0016101b2565b3480156103fd57600080fd5b5061020e61040c366004611a1f565b61099e565b34801561041d57600080fd5b506101ca60075481565b34801561043357600080fd5b5061020e610b7d565b34801561044857600080fd5b5061045c6104573660046119db565b610ba7565b6040516101b29796959493929190611aae565b34801561047b57600080fd5b506101ca60025481565b34801561049157600080fd5b507f00000000000000000000000000000000000000000000000000000000000000006102ce565b3480156104c457600080fd5b5061020e6104d3366004611a1f565b610d2a565b3480156104e457600080fd5b5061020e6104f3366004611af5565b610f01565b606060006105068385611b61565b60085490915081111561051857506008545b60006105248583611b74565b67ffffffffffffffff81111561053c5761053c6118c4565b604051908082528060200260200182016040528015610565578160200160208202803683370190505b509050845b828110156105eb576008818154811061058557610585611b87565b6000918252602090912060088204015460079091166004026101000a900463ffffffff16826105b48884611b74565b815181106105c4576105c4611b87565b63ffffffff90921660209283029190910190910152806105e381611bb6565b91505061056a565b509150505b92915050565b7f00000000000000000000000000000000000000000000000000000000000000003373ffffffffffffffffffffffffffffffffffffffff821614610689576040517f8ba9316e00000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff821660248201526044015b60405180910390fd5b6106938383610f15565b505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610719576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610680565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b61079d61114b565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb6107f860005473ffffffffffffffffffffffffffffffffffffffff1690565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602481018490526044016020604051808303816000875af115801561086a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061088e9190611bee565b5050565b61089a61114b565b6000805460405173ffffffffffffffffffffffffffffffffffffffff9091169083908381818185875af1925050503d80600081146108f4576040519150601f19603f3d011682016040523d82523d6000602084013e6108f9565b606091505b505090508061088e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f77697468647261774e6174697665206661696c656400000000000000000000006044820152606401610680565b6008818154811061097457600080fd5b9060005260206000209060089182820401919006600402915054906101000a900463ffffffff1681565b60005b8161ffff168161ffff161015610b765760006109cd6040518060200160405280600015158152506111cc565b90506000806109de88888886611288565b6007829055909250905060006109f26114cb565b604080516101008101825284815260006020808301828152845183815280830186528486019081524260608601526080850184905260a0850187905260c0850184905260e08501849052898452600a8352949092208351815591516001830180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169115159190911790559251805194955091939092610a9a9260028501929101906117d7565b5060608201516003828101919091556080830151600483015560a0830151600583015560c0830151600683015560e090920151600790910180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169115159190911790558054906000610b0d83611bb6565b9091555050600083815260096020526040908190208290555183907f5f56b4c20db9f5b294cbf6f681368de4a992a27e2de2ee702dcf2cbbfa791ec490610b579085815260200190565b60405180910390a2505050508080610b6e90611c10565b9150506109a1565b5050505050565b6000600481905560058190556103e760065560038190556002819055610ba590600890611822565b565b6000818152600a602052604081205481906060908290819081908190610c29576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f72657175657374206e6f7420666f756e640000000000000000000000000000006044820152606401610680565b6000888152600a6020908152604080832081516101008101835281548152600182015460ff16151581850152600282018054845181870281018701865281815292959394860193830182828015610c9f57602002820191906000526020600020905b815481526020019060010190808311610c8b575b50505050508152602001600382015481526020016004820154815260200160058201548152602001600682015481526020016007820160009054906101000a900460ff1615151515815250509050806000015181602001518260400151836060015184608001518560a001518660c00151975097509750975097509750975050919395979092949650565b60005b8161ffff168161ffff161015610b76576000610d596040518060200160405280600115158152506111cc565b9050600080610d6a88888886611559565b600782905590925090506000610d7e6114cb565b604080516101008101825284815260006020808301828152845183815280830186528486019081524260608601526080850184905260a0850187905260c08501849052600160e086018190528a8552600a84529590932084518155905194810180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001695151595909517909455905180519495509193610e2592600285019201906117d7565b5060608201516003828101919091556080830151600483015560a0830151600583015560c0830151600683015560e090920151600790910180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169115159190911790558054906000610e9883611bb6565b9091555050600083815260096020526040908190208290555183907f5f56b4c20db9f5b294cbf6f681368de4a992a27e2de2ee702dcf2cbbfa791ec490610ee29085815260200190565b60405180910390a2505050508080610ef990611c10565b915050610d2d565b610f0961114b565b610f12816116bf565b50565b6000828152600a6020526040902054610f8a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f72657175657374206e6f7420666f756e640000000000000000000000000000006044820152606401610680565b6000610f946114cb565b60008481526009602052604081205491925090610fb19083611b74565b90506000610fc282620f4240611c31565b9050600554821115610fd45760058290555b6006548210610fe557600654610fe7565b815b600655600254610ff7578061102a565b600254611005906001611b61565b816002546004546110169190611c31565b6110209190611b61565b61102a9190611c48565b6004556002805490600061103d83611bb6565b90915550506000858152600a60209081526040909120600181810180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690911790558551611095926002909201918701906117d7565b506000858152600a6020526040908190204260048083019190915560068201869055600880546001810182557ff3f7a9fe364faab93b216da50a3214154f22a0a2b415b23a84c8169e8b636ee391810491909101805460079092169092026101000a63ffffffff81810219909216918716021790555490517f6c84e12b4c188e61f1b4727024a5cf05c025fa58467e5eedf763c0744c89da7b9161113c9188918891611c83565b60405180910390a15050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610ba5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610680565b60607f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa8260405160240161120591511515815260200190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915292915050565b6040517f27e5c50a00000000000000000000000000000000000000000000000000000000815263ffffffff808616600483015283166024820152600090819073ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016906327e5c50a90604401602060405180830381865afa158015611325573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113499190611cac565b90507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16634000aea07f000000000000000000000000000000000000000000000000000000000000000083898989896040516020016113c09493929190611d29565b6040516020818303038152906040526040518463ffffffff1660e01b81526004016113ed93929190611d66565b6020604051808303816000875af115801561140c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114309190611bee565b507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663fc2a88c36040518163ffffffff1660e01b8152600401602060405180830381865afa15801561149c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114c09190611cac565b915094509492505050565b6000466114d7816117b4565b1561155257606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015611528573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061154c9190611cac565b91505090565b4391505090565b6040517f13c34b7f00000000000000000000000000000000000000000000000000000000815263ffffffff808616600483015283166024820152600090819073ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016906313c34b7f90604401602060405180830381865afa1580156115f6573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061161a9190611cac565b90507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16639cfc058e82888888886040518663ffffffff1660e01b815260040161167c9493929190611d29565b60206040518083038185885af115801561169a573d6000803e3d6000fd5b50505050506040513d601f19601f820116820180604052508101906114c09190611cac565b3373ffffffffffffffffffffffffffffffffffffffff82160361173e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610680565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600061a4b18214806117c8575062066eed82145b806105f057505062066eee1490565b828054828255906000526020600020908101928215611812579160200282015b828111156118125782518255916020019190600101906117f7565b5061181e929150611843565b5090565b508054600082556007016008900490600052602060002090810190610f1291905b5b8082111561181e5760008155600101611844565b6000806040838503121561186b57600080fd5b50508035926020909101359150565b6020808252825182820181905260009190848201906040850190845b818110156118b857835163ffffffff1683529284019291840191600101611896565b50909695505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000806040838503121561190657600080fd5b8235915060208084013567ffffffffffffffff8082111561192657600080fd5b818601915086601f83011261193a57600080fd5b81358181111561194c5761194c6118c4565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f8301168101818110858211171561198f5761198f6118c4565b6040529182528482019250838101850191898311156119ad57600080fd5b938501935b828510156119cb578435845293850193928501926119b2565b8096505050505050509250929050565b6000602082840312156119ed57600080fd5b5035919050565b803563ffffffff81168114611a0857600080fd5b919050565b803561ffff81168114611a0857600080fd5b60008060008060808587031215611a3557600080fd5b611a3e856119f4565b9350611a4c60208601611a0d565b9250611a5a604086016119f4565b9150611a6860608601611a0d565b905092959194509250565b600081518084526020808501945080840160005b83811015611aa357815187529582019590820190600101611a87565b509495945050505050565b878152861515602082015260e060408201526000611acf60e0830188611a73565b90508560608301528460808301528360a08301528260c083015298975050505050505050565b600060208284031215611b0757600080fd5b813573ffffffffffffffffffffffffffffffffffffffff81168114611b2b57600080fd5b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b808201808211156105f0576105f0611b32565b818103818111156105f0576105f0611b32565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611be757611be7611b32565b5060010190565b600060208284031215611c0057600080fd5b81518015158114611b2b57600080fd5b600061ffff808316818103611c2757611c27611b32565b6001019392505050565b80820281158282048414176105f0576105f0611b32565b600082611c7e577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b838152606060208201526000611c9c6060830185611a73565b9050826040830152949350505050565b600060208284031215611cbe57600080fd5b5051919050565b6000815180845260005b81811015611ceb57602081850181015186830182015201611ccf565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b600063ffffffff808716835261ffff8616602084015280851660408401525060806060830152611d5c6080830184611cc5565b9695505050505050565b73ffffffffffffffffffffffffffffffffffffffff84168152826020820152606060408201526000611d9b6060830184611cc5565b9594505050505056fea164736f6c6343000813000a",
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
