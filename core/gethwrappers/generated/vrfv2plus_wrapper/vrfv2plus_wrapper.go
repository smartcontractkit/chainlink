// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrfv2plus_wrapper

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

var VRFV2PlusWrapperMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_coordinator\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"LinkAlreadySet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"}],\"name\":\"OnlyOwnerOrCoordinator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"WrapperFulfillmentFailed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COORDINATOR\",\"outputs\":[{\"internalType\":\"contractExtendedVRFCoordinatorV2PlusInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PRICE_REGISTRY\",\"outputs\":[{\"internalType\":\"contractIVRFV2PlusPriceRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SUBSCRIPTION_ID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disable\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enable\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"maxNumWords\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPriceRegistry\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRequestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"_requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"_numWords\",\"type\":\"uint32\"}],\"name\":\"requestRandomWordsInNative\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_callbacks\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"callbackAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"requestGasPrice\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_configured\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_disabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_link\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_wrapperGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"_maxNumWords\",\"type\":\"uint8\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"}],\"name\":\"setCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"name\":\"setLINK\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawNative\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b506040516200220c3803806200220c83398101604081905262000034916200038a565b8033806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf816200029c565b5050600280546001600160a01b0319166001600160a01b03938416179055508216156200010257600380546001600160a01b0319166001600160a01b0384161790555b806001600160a01b03166080816001600160a01b031660601b81525050806001600160a01b0316631ed257266040518163ffffffff1660e01b815260040160206040518083038186803b1580156200015957600080fd5b505afa1580156200016e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000194919062000365565b6001600160a01b031660a0816001600160a01b031660601b815250506000816001600160a01b031663a21a23e46040518163ffffffff1660e01b8152600401602060405180830381600087803b158015620001ee57600080fd5b505af115801562000203573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002299190620003c2565b60c0819052604051632fb1302360e21b8152600481018290523060248201529091506001600160a01b0383169063bec4c08c90604401600060405180830381600087803b1580156200027a57600080fd5b505af11580156200028f573d6000803e3d6000fd5b50505050505050620003dc565b6001600160a01b038116331415620002f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03811681146200036057600080fd5b919050565b6000602082840312156200037857600080fd5b620003838262000348565b9392505050565b600080604083850312156200039e57600080fd5b620003a98362000348565b9150620003b96020840162000348565b90509250929050565b600060208284031215620003d557600080fd5b5051919050565b60805160601c60a05160601c60c051611dc962000443600039600081816101a20152818161099201526111a6015260008181610208015281816102ab01528181610820015261103401526000818161033401528181610a5c015261122d0152611dc96000f3fe60806040526004361061018b5760003560e01c806379ba5097116100d6578063a608a1e11161007f578063f2fde38b11610059578063f2fde38b14610560578063f3fef3a314610580578063fc2a88c3146105a057600080fd5b8063a608a1e1146104ea578063c3f909d414610509578063da4f5e6d1461053357600080fd5b8063a02e0616116100b0578063a02e061614610495578063a3907d71146104b5578063a4c0ed36146104ca57600080fd5b806379ba5097146104355780638da5cb5b1461044a5780638ea981171461047557600080fd5b80632f2770db1161013857806348baa1c51161011257806348baa1c51461035657806357a8070a146103f857806362a504fc1461042257600080fd5b80632f2770db146102ed5780633549456e146103025780633b2bcbf11461032257600080fd5b8063181f5a7711610169578063181f5a771461024d5780631ed25726146102995780631fe543e3146102cd57600080fd5b8063030932bb1461019057806307b18bde146101d75780630d6c107e146101f9575b600080fd5b34801561019c57600080fd5b506101c47f000000000000000000000000000000000000000000000000000000000000000081565b6040519081526020015b60405180910390f35b3480156101e357600080fd5b506101f76101f236600461193d565b6105b6565b005b34801561020557600080fd5b507f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101ce565b34801561025957600080fd5b50604080518082018252601281527f56524656325772617070657220312e302e300000000000000000000000000000602082015290516101ce9190611c2e565b3480156102a557600080fd5b506102287f000000000000000000000000000000000000000000000000000000000000000081565b3480156102d957600080fd5b506101f76102e8366004611a42565b610692565b3480156102f957600080fd5b506101f7610713565b34801561030e57600080fd5b506101f761031d366004611b31565b610749565b34801561032e57600080fd5b506102287f000000000000000000000000000000000000000000000000000000000000000081565b34801561036257600080fd5b506103c1610371366004611a10565b6008602052600090815260409020805460019091015473ffffffffffffffffffffffffffffffffffffffff82169174010000000000000000000000000000000000000000900463ffffffff169083565b6040805173ffffffffffffffffffffffffffffffffffffffff909416845263ffffffff9092166020840152908201526060016101ce565b34801561040457600080fd5b506005546104129060ff1681565b60405190151581526020016101ce565b6101c4610430366004611b77565b6107c5565b34801561044157600080fd5b506101f7610b8c565b34801561045657600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff16610228565b34801561048157600080fd5b506101f761049036600461191b565b610c89565b3480156104a157600080fd5b506101f76104b036600461191b565b610d94565b3480156104c157600080fd5b506101f7610e33565b3480156104d657600080fd5b506101f76104e5366004611967565b610e65565b3480156104f657600080fd5b5060055461041290610100900460ff1681565b34801561051557600080fd5b506006546007546040805192835260ff9091166020830152016101ce565b34801561053f57600080fd5b506003546102289073ffffffffffffffffffffffffffffffffffffffff1681565b34801561056c57600080fd5b506101f761057b36600461191b565b61138b565b34801561058c57600080fd5b506101f761059b36600461193d565b61139f565b3480156105ac57600080fd5b506101c460045481565b6105be611453565b60008273ffffffffffffffffffffffffffffffffffffffff168260405160006040518083038185875af1925050503d8060008114610618576040519150601f19603f3d011682016040523d82523d6000602084013e61061d565b606091505b505090508061068d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f6661696c656420746f207769746864726177206e61746976650000000000000060448201526064015b60405180910390fd5b505050565b60025473ffffffffffffffffffffffffffffffffffffffff163314610705576002546040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff9091166024820152604401610684565b61070f82826114d6565b5050565b61071b611453565b600580547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff16610100179055565b610751611453565b600580546006939093556007805460ff9093167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0093841617905563ffffffff9093166201000002167fffffffffffffffffffffffffffffffffffffffffffffffffffff00000000ff00909116176001179055565b6000806107d1856116c2565b6040517f907d645900000000000000000000000000000000000000000000000000000000815263ffffffff8716600482015290915060009073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169063907d64599060240160206040518083038186803b15801561086257600080fd5b505afa158015610876573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061089a9190611a29565b905080341015610906576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f66656520746f6f206c6f770000000000000000000000000000000000000000006044820152606401610684565b60075460ff1663ffffffff8516111561097b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f6e756d576f72647320746f6f20686967680000000000000000000000000000006044820152606401610684565b60006040518060c0016040528060065481526020017f000000000000000000000000000000000000000000000000000000000000000081526020018761ffff168152602001600560029054906101000a900463ffffffff16858a6109df9190611cf4565b6109e99190611cf4565b63ffffffff1681526020018663ffffffff168152602001610a1a6040518060200160405280600115158152506116e0565b90526040517f9b1c385e00000000000000000000000000000000000000000000000000000000815290915073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001690639b1c385e90610a91908490600401611c41565b602060405180830381600087803b158015610aab57600080fd5b505af1158015610abf573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ae39190611a29565b6040805160608101825233815263ffffffff808b1660208084019182523a8486019081526000878152600890925294902092518354915190921674010000000000000000000000000000000000000000027fffffffffffffffff00000000000000000000000000000000000000000000000090911673ffffffffffffffffffffffffffffffffffffffff9290921691909117178155905160019091015593505050509392505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610c0d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610684565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60005473ffffffffffffffffffffffffffffffffffffffff163314801590610cc9575060025473ffffffffffffffffffffffffffffffffffffffff163314155b15610d4d5733610cee60005473ffffffffffffffffffffffffffffffffffffffff1690565b6002546040517f061db9c100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff93841660048201529183166024830152919091166044820152606401610684565b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b610d9c611453565b60035473ffffffffffffffffffffffffffffffffffffffff1615610dec576040517f2d118a6e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b610e3b611453565b600580547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055565b60055460ff16610ed1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f77726170706572206973206e6f7420636f6e66696775726564000000000000006044820152606401610684565b600554610100900460ff1615610f43576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f777261707065722069732064697361626c6564000000000000000000000000006044820152606401610684565b60035473ffffffffffffffffffffffffffffffffffffffff163314610fc4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f6f6e6c792063616c6c61626c652066726f6d204c494e4b0000000000000000006044820152606401610684565b60008080610fd484860186611b77565b9250925092506000610fe5846116c2565b6040517fe16ad7cf00000000000000000000000000000000000000000000000000000000815263ffffffff8616600482015290915060009073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169063e16ad7cf9060240160206040518083038186803b15801561107657600080fd5b505afa15801561108a573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110ae9190611a29565b90508088101561111a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f66656520746f6f206c6f770000000000000000000000000000000000000000006044820152606401610684565b60075460ff1663ffffffff8416111561118f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f6e756d576f72647320746f6f20686967680000000000000000000000000000006044820152606401610684565b60006040518060c0016040528060065481526020017f000000000000000000000000000000000000000000000000000000000000000081526020018661ffff168152602001600560029054906101000a900463ffffffff1685896111f39190611cf4565b6111fd9190611cf4565b63ffffffff1681526020018563ffffffff16815260200160405180602001604052806000815250815250905060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16639b1c385e836040518263ffffffff1660e01b81526004016112849190611c41565b602060405180830381600087803b15801561129e57600080fd5b505af11580156112b2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112d69190611a29565b6040805160608101825273ffffffffffffffffffffffffffffffffffffffff9d8e16815263ffffffff998a1660208083019182523a838501908152600086815260089092529390209151825491519f167fffffffffffffffff00000000000000000000000000000000000000000000000090911617740100000000000000000000000000000000000000009e909a169d909d02989098178c5596516001909b019a909a55505050600492909255505050505050565b611393611453565b61139c8161179c565b50565b6113a7611453565b6003546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8481166004830152602482018490529091169063a9059cbb90604401602060405180830381600087803b15801561141b57600080fd5b505af115801561142f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061068d91906119ee565b60005473ffffffffffffffffffffffffffffffffffffffff1633146114d4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610684565b565b60008281526008602081815260408084208151606081018352815473ffffffffffffffffffffffffffffffffffffffff808216835263ffffffff740100000000000000000000000000000000000000008304168387015260018401805495840195909552898852959094527fffffffffffffffff0000000000000000000000000000000000000000000000009093169055929092558151166115d4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f72657175657374206e6f7420666f756e640000000000000000000000000000006044820152606401610684565b600080631fe543e360e01b85856040516024016115f2929190611ca6565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050509050600061166c846020015163ffffffff16856000015184611892565b9050806116ba57835160405173ffffffffffffffffffffffffffffffffffffffff9091169087907fc551b83c151f2d1c7eeb938ac59008e0409f1c1dc1e2f112449d4d79b458902290600090a35b505050505050565b60006116cf603f83611d43565b6116da906001611cf4565b92915050565b60607f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa8260405160240161171991511515815260200190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915292915050565b73ffffffffffffffffffffffffffffffffffffffff811633141561181c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610684565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60005a6113888110156118a457600080fd5b6113888103905084604082048203116118bc57600080fd5b50823b6118c857600080fd5b60008083516020850160008789f1949350505050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461190257600080fd5b919050565b803563ffffffff8116811461190257600080fd5b60006020828403121561192d57600080fd5b611936826118de565b9392505050565b6000806040838503121561195057600080fd5b611959836118de565b946020939093013593505050565b6000806000806060858703121561197d57600080fd5b611986856118de565b935060208501359250604085013567ffffffffffffffff808211156119aa57600080fd5b818701915087601f8301126119be57600080fd5b8135818111156119cd57600080fd5b8860208285010111156119df57600080fd5b95989497505060200194505050565b600060208284031215611a0057600080fd5b8151801515811461193657600080fd5b600060208284031215611a2257600080fd5b5035919050565b600060208284031215611a3b57600080fd5b5051919050565b60008060408385031215611a5557600080fd5b8235915060208084013567ffffffffffffffff80821115611a7557600080fd5b818601915086601f830112611a8957600080fd5b813581811115611a9b57611a9b611d8d565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f83011681018181108582111715611ade57611ade611d8d565b604052828152858101935084860182860187018b1015611afd57600080fd5b600095505b83861015611b20578035855260019590950194938601938601611b02565b508096505050505050509250929050565b600080600060608486031215611b4657600080fd5b611b4f84611907565b925060208401359150604084013560ff81168114611b6c57600080fd5b809150509250925092565b600080600060608486031215611b8c57600080fd5b611b9584611907565b9250602084013561ffff81168114611bac57600080fd5b9150611bba60408501611907565b90509250925092565b6000815180845260005b81811015611be957602081850181015186830182015201611bcd565b81811115611bfb576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006119366020830184611bc3565b60208152815160208201526020820151604082015261ffff60408301511660608201526000606083015163ffffffff80821660808501528060808601511660a0850152505060a083015160c080840152611c9e60e0840182611bc3565b949350505050565b6000604082018483526020604081850152818551808452606086019150828701935060005b81811015611ce757845183529383019391830191600101611ccb565b5090979650505050505050565b600063ffffffff808316818516808303821115611d3a577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b01949350505050565b600063ffffffff80841680611d81577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b92169190910492915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFV2PlusWrapperABI = VRFV2PlusWrapperMetaData.ABI

var VRFV2PlusWrapperBin = VRFV2PlusWrapperMetaData.Bin

func DeployVRFV2PlusWrapper(auth *bind.TransactOpts, backend bind.ContractBackend, _link common.Address, _coordinator common.Address) (common.Address, *types.Transaction, *VRFV2PlusWrapper, error) {
	parsed, err := VRFV2PlusWrapperMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2PlusWrapperBin), backend, _link, _coordinator)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2PlusWrapper{VRFV2PlusWrapperCaller: VRFV2PlusWrapperCaller{contract: contract}, VRFV2PlusWrapperTransactor: VRFV2PlusWrapperTransactor{contract: contract}, VRFV2PlusWrapperFilterer: VRFV2PlusWrapperFilterer{contract: contract}}, nil
}

type VRFV2PlusWrapper struct {
	address common.Address
	abi     abi.ABI
	VRFV2PlusWrapperCaller
	VRFV2PlusWrapperTransactor
	VRFV2PlusWrapperFilterer
}

type VRFV2PlusWrapperCaller struct {
	contract *bind.BoundContract
}

type VRFV2PlusWrapperTransactor struct {
	contract *bind.BoundContract
}

type VRFV2PlusWrapperFilterer struct {
	contract *bind.BoundContract
}

type VRFV2PlusWrapperSession struct {
	Contract     *VRFV2PlusWrapper
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2PlusWrapperCallerSession struct {
	Contract *VRFV2PlusWrapperCaller
	CallOpts bind.CallOpts
}

type VRFV2PlusWrapperTransactorSession struct {
	Contract     *VRFV2PlusWrapperTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2PlusWrapperRaw struct {
	Contract *VRFV2PlusWrapper
}

type VRFV2PlusWrapperCallerRaw struct {
	Contract *VRFV2PlusWrapperCaller
}

type VRFV2PlusWrapperTransactorRaw struct {
	Contract *VRFV2PlusWrapperTransactor
}

func NewVRFV2PlusWrapper(address common.Address, backend bind.ContractBackend) (*VRFV2PlusWrapper, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2PlusWrapperABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2PlusWrapper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusWrapper{address: address, abi: abi, VRFV2PlusWrapperCaller: VRFV2PlusWrapperCaller{contract: contract}, VRFV2PlusWrapperTransactor: VRFV2PlusWrapperTransactor{contract: contract}, VRFV2PlusWrapperFilterer: VRFV2PlusWrapperFilterer{contract: contract}}, nil
}

func NewVRFV2PlusWrapperCaller(address common.Address, caller bind.ContractCaller) (*VRFV2PlusWrapperCaller, error) {
	contract, err := bindVRFV2PlusWrapper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusWrapperCaller{contract: contract}, nil
}

func NewVRFV2PlusWrapperTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2PlusWrapperTransactor, error) {
	contract, err := bindVRFV2PlusWrapper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusWrapperTransactor{contract: contract}, nil
}

func NewVRFV2PlusWrapperFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2PlusWrapperFilterer, error) {
	contract, err := bindVRFV2PlusWrapper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusWrapperFilterer{contract: contract}, nil
}

func bindVRFV2PlusWrapper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFV2PlusWrapperMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusWrapper.Contract.VRFV2PlusWrapperCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.VRFV2PlusWrapperTransactor.contract.Transfer(opts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.VRFV2PlusWrapperTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusWrapper.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.contract.Transfer(opts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperCaller) COORDINATOR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusWrapper.contract.Call(opts, &out, "COORDINATOR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperSession) COORDINATOR() (common.Address, error) {
	return _VRFV2PlusWrapper.Contract.COORDINATOR(&_VRFV2PlusWrapper.CallOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperCallerSession) COORDINATOR() (common.Address, error) {
	return _VRFV2PlusWrapper.Contract.COORDINATOR(&_VRFV2PlusWrapper.CallOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperCaller) PRICEREGISTRY(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusWrapper.contract.Call(opts, &out, "PRICE_REGISTRY")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperSession) PRICEREGISTRY() (common.Address, error) {
	return _VRFV2PlusWrapper.Contract.PRICEREGISTRY(&_VRFV2PlusWrapper.CallOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperCallerSession) PRICEREGISTRY() (common.Address, error) {
	return _VRFV2PlusWrapper.Contract.PRICEREGISTRY(&_VRFV2PlusWrapper.CallOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperCaller) SUBSCRIPTIONID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusWrapper.contract.Call(opts, &out, "SUBSCRIPTION_ID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperSession) SUBSCRIPTIONID() (*big.Int, error) {
	return _VRFV2PlusWrapper.Contract.SUBSCRIPTIONID(&_VRFV2PlusWrapper.CallOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperCallerSession) SUBSCRIPTIONID() (*big.Int, error) {
	return _VRFV2PlusWrapper.Contract.SUBSCRIPTIONID(&_VRFV2PlusWrapper.CallOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperCaller) GetConfig(opts *bind.CallOpts) (GetConfig,

	error) {
	var out []interface{}
	err := _VRFV2PlusWrapper.contract.Call(opts, &out, "getConfig")

	outstruct := new(GetConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.KeyHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.MaxNumWords = *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return *outstruct, err

}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperSession) GetConfig() (GetConfig,

	error) {
	return _VRFV2PlusWrapper.Contract.GetConfig(&_VRFV2PlusWrapper.CallOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperCallerSession) GetConfig() (GetConfig,

	error) {
	return _VRFV2PlusWrapper.Contract.GetConfig(&_VRFV2PlusWrapper.CallOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperCaller) GetPriceRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusWrapper.contract.Call(opts, &out, "getPriceRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperSession) GetPriceRegistry() (common.Address, error) {
	return _VRFV2PlusWrapper.Contract.GetPriceRegistry(&_VRFV2PlusWrapper.CallOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperCallerSession) GetPriceRegistry() (common.Address, error) {
	return _VRFV2PlusWrapper.Contract.GetPriceRegistry(&_VRFV2PlusWrapper.CallOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperCaller) LastRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2PlusWrapper.contract.Call(opts, &out, "lastRequestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperSession) LastRequestId() (*big.Int, error) {
	return _VRFV2PlusWrapper.Contract.LastRequestId(&_VRFV2PlusWrapper.CallOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperCallerSession) LastRequestId() (*big.Int, error) {
	return _VRFV2PlusWrapper.Contract.LastRequestId(&_VRFV2PlusWrapper.CallOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusWrapper.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperSession) Owner() (common.Address, error) {
	return _VRFV2PlusWrapper.Contract.Owner(&_VRFV2PlusWrapper.CallOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperCallerSession) Owner() (common.Address, error) {
	return _VRFV2PlusWrapper.Contract.Owner(&_VRFV2PlusWrapper.CallOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperCaller) SCallbacks(opts *bind.CallOpts, arg0 *big.Int) (SCallbacks,

	error) {
	var out []interface{}
	err := _VRFV2PlusWrapper.contract.Call(opts, &out, "s_callbacks", arg0)

	outstruct := new(SCallbacks)
	if err != nil {
		return *outstruct, err
	}

	outstruct.CallbackAddress = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.CallbackGasLimit = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.RequestGasPrice = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperSession) SCallbacks(arg0 *big.Int) (SCallbacks,

	error) {
	return _VRFV2PlusWrapper.Contract.SCallbacks(&_VRFV2PlusWrapper.CallOpts, arg0)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperCallerSession) SCallbacks(arg0 *big.Int) (SCallbacks,

	error) {
	return _VRFV2PlusWrapper.Contract.SCallbacks(&_VRFV2PlusWrapper.CallOpts, arg0)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperCaller) SConfigured(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _VRFV2PlusWrapper.contract.Call(opts, &out, "s_configured")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperSession) SConfigured() (bool, error) {
	return _VRFV2PlusWrapper.Contract.SConfigured(&_VRFV2PlusWrapper.CallOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperCallerSession) SConfigured() (bool, error) {
	return _VRFV2PlusWrapper.Contract.SConfigured(&_VRFV2PlusWrapper.CallOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperCaller) SDisabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _VRFV2PlusWrapper.contract.Call(opts, &out, "s_disabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperSession) SDisabled() (bool, error) {
	return _VRFV2PlusWrapper.Contract.SDisabled(&_VRFV2PlusWrapper.CallOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperCallerSession) SDisabled() (bool, error) {
	return _VRFV2PlusWrapper.Contract.SDisabled(&_VRFV2PlusWrapper.CallOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperCaller) SLink(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2PlusWrapper.contract.Call(opts, &out, "s_link")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperSession) SLink() (common.Address, error) {
	return _VRFV2PlusWrapper.Contract.SLink(&_VRFV2PlusWrapper.CallOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperCallerSession) SLink() (common.Address, error) {
	return _VRFV2PlusWrapper.Contract.SLink(&_VRFV2PlusWrapper.CallOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VRFV2PlusWrapper.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperSession) TypeAndVersion() (string, error) {
	return _VRFV2PlusWrapper.Contract.TypeAndVersion(&_VRFV2PlusWrapper.CallOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperCallerSession) TypeAndVersion() (string, error) {
	return _VRFV2PlusWrapper.Contract.TypeAndVersion(&_VRFV2PlusWrapper.CallOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.contract.Transact(opts, "acceptOwnership")
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.AcceptOwnership(&_VRFV2PlusWrapper.TransactOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.AcceptOwnership(&_VRFV2PlusWrapper.TransactOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperTransactor) Disable(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.contract.Transact(opts, "disable")
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperSession) Disable() (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.Disable(&_VRFV2PlusWrapper.TransactOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperTransactorSession) Disable() (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.Disable(&_VRFV2PlusWrapper.TransactOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperTransactor) Enable(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.contract.Transact(opts, "enable")
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperSession) Enable() (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.Enable(&_VRFV2PlusWrapper.TransactOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperTransactorSession) Enable() (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.Enable(&_VRFV2PlusWrapper.TransactOpts)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperTransactor) OnTokenTransfer(opts *bind.TransactOpts, _sender common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.contract.Transact(opts, "onTokenTransfer", _sender, _amount, _data)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperSession) OnTokenTransfer(_sender common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.OnTokenTransfer(&_VRFV2PlusWrapper.TransactOpts, _sender, _amount, _data)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperTransactorSession) OnTokenTransfer(_sender common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.OnTokenTransfer(&_VRFV2PlusWrapper.TransactOpts, _sender, _amount, _data)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.RawFulfillRandomWords(&_VRFV2PlusWrapper.TransactOpts, requestId, randomWords)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.RawFulfillRandomWords(&_VRFV2PlusWrapper.TransactOpts, requestId, randomWords)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperTransactor) RequestRandomWordsInNative(opts *bind.TransactOpts, _callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.contract.Transact(opts, "requestRandomWordsInNative", _callbackGasLimit, _requestConfirmations, _numWords)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperSession) RequestRandomWordsInNative(_callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.RequestRandomWordsInNative(&_VRFV2PlusWrapper.TransactOpts, _callbackGasLimit, _requestConfirmations, _numWords)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperTransactorSession) RequestRandomWordsInNative(_callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.RequestRandomWordsInNative(&_VRFV2PlusWrapper.TransactOpts, _callbackGasLimit, _requestConfirmations, _numWords)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperTransactor) SetConfig(opts *bind.TransactOpts, _wrapperGasOverhead uint32, _keyHash [32]byte, _maxNumWords uint8) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.contract.Transact(opts, "setConfig", _wrapperGasOverhead, _keyHash, _maxNumWords)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperSession) SetConfig(_wrapperGasOverhead uint32, _keyHash [32]byte, _maxNumWords uint8) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.SetConfig(&_VRFV2PlusWrapper.TransactOpts, _wrapperGasOverhead, _keyHash, _maxNumWords)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperTransactorSession) SetConfig(_wrapperGasOverhead uint32, _keyHash [32]byte, _maxNumWords uint8) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.SetConfig(&_VRFV2PlusWrapper.TransactOpts, _wrapperGasOverhead, _keyHash, _maxNumWords)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperTransactor) SetCoordinator(opts *bind.TransactOpts, _vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.contract.Transact(opts, "setCoordinator", _vrfCoordinator)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperSession) SetCoordinator(_vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.SetCoordinator(&_VRFV2PlusWrapper.TransactOpts, _vrfCoordinator)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperTransactorSession) SetCoordinator(_vrfCoordinator common.Address) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.SetCoordinator(&_VRFV2PlusWrapper.TransactOpts, _vrfCoordinator)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperTransactor) SetLINK(opts *bind.TransactOpts, link common.Address) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.contract.Transact(opts, "setLINK", link)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperSession) SetLINK(link common.Address) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.SetLINK(&_VRFV2PlusWrapper.TransactOpts, link)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperTransactorSession) SetLINK(link common.Address) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.SetLINK(&_VRFV2PlusWrapper.TransactOpts, link)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.TransferOwnership(&_VRFV2PlusWrapper.TransactOpts, to)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.TransferOwnership(&_VRFV2PlusWrapper.TransactOpts, to)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperTransactor) Withdraw(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.contract.Transact(opts, "withdraw", _recipient, _amount)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperSession) Withdraw(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.Withdraw(&_VRFV2PlusWrapper.TransactOpts, _recipient, _amount)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperTransactorSession) Withdraw(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.Withdraw(&_VRFV2PlusWrapper.TransactOpts, _recipient, _amount)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperTransactor) WithdrawNative(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.contract.Transact(opts, "withdrawNative", _recipient, _amount)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperSession) WithdrawNative(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.WithdrawNative(&_VRFV2PlusWrapper.TransactOpts, _recipient, _amount)
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperTransactorSession) WithdrawNative(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _VRFV2PlusWrapper.Contract.WithdrawNative(&_VRFV2PlusWrapper.TransactOpts, _recipient, _amount)
}

type VRFV2PlusWrapperOwnershipTransferRequestedIterator struct {
	Event *VRFV2PlusWrapperOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusWrapperOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusWrapperOwnershipTransferRequested)
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
		it.Event = new(VRFV2PlusWrapperOwnershipTransferRequested)
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

func (it *VRFV2PlusWrapperOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusWrapperOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusWrapperOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusWrapperOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusWrapper.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusWrapperOwnershipTransferRequestedIterator{contract: _VRFV2PlusWrapper.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2PlusWrapperOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusWrapper.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusWrapperOwnershipTransferRequested)
				if err := _VRFV2PlusWrapper.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFV2PlusWrapper *VRFV2PlusWrapperFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFV2PlusWrapperOwnershipTransferRequested, error) {
	event := new(VRFV2PlusWrapperOwnershipTransferRequested)
	if err := _VRFV2PlusWrapper.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2PlusWrapperOwnershipTransferredIterator struct {
	Event *VRFV2PlusWrapperOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusWrapperOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusWrapperOwnershipTransferred)
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
		it.Event = new(VRFV2PlusWrapperOwnershipTransferred)
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

func (it *VRFV2PlusWrapperOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusWrapperOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusWrapperOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusWrapperOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusWrapper.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusWrapperOwnershipTransferredIterator{contract: _VRFV2PlusWrapper.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2PlusWrapperOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2PlusWrapper.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusWrapperOwnershipTransferred)
				if err := _VRFV2PlusWrapper.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFV2PlusWrapper *VRFV2PlusWrapperFilterer) ParseOwnershipTransferred(log types.Log) (*VRFV2PlusWrapperOwnershipTransferred, error) {
	event := new(VRFV2PlusWrapperOwnershipTransferred)
	if err := _VRFV2PlusWrapper.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2PlusWrapperWrapperFulfillmentFailedIterator struct {
	Event *VRFV2PlusWrapperWrapperFulfillmentFailed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2PlusWrapperWrapperFulfillmentFailedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2PlusWrapperWrapperFulfillmentFailed)
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
		it.Event = new(VRFV2PlusWrapperWrapperFulfillmentFailed)
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

func (it *VRFV2PlusWrapperWrapperFulfillmentFailedIterator) Error() error {
	return it.fail
}

func (it *VRFV2PlusWrapperWrapperFulfillmentFailedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2PlusWrapperWrapperFulfillmentFailed struct {
	RequestId *big.Int
	Consumer  common.Address
	Raw       types.Log
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperFilterer) FilterWrapperFulfillmentFailed(opts *bind.FilterOpts, requestId []*big.Int, consumer []common.Address) (*VRFV2PlusWrapperWrapperFulfillmentFailedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var consumerRule []interface{}
	for _, consumerItem := range consumer {
		consumerRule = append(consumerRule, consumerItem)
	}

	logs, sub, err := _VRFV2PlusWrapper.contract.FilterLogs(opts, "WrapperFulfillmentFailed", requestIdRule, consumerRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusWrapperWrapperFulfillmentFailedIterator{contract: _VRFV2PlusWrapper.contract, event: "WrapperFulfillmentFailed", logs: logs, sub: sub}, nil
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapperFilterer) WatchWrapperFulfillmentFailed(opts *bind.WatchOpts, sink chan<- *VRFV2PlusWrapperWrapperFulfillmentFailed, requestId []*big.Int, consumer []common.Address) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var consumerRule []interface{}
	for _, consumerItem := range consumer {
		consumerRule = append(consumerRule, consumerItem)
	}

	logs, sub, err := _VRFV2PlusWrapper.contract.WatchLogs(opts, "WrapperFulfillmentFailed", requestIdRule, consumerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2PlusWrapperWrapperFulfillmentFailed)
				if err := _VRFV2PlusWrapper.contract.UnpackLog(event, "WrapperFulfillmentFailed", log); err != nil {
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

func (_VRFV2PlusWrapper *VRFV2PlusWrapperFilterer) ParseWrapperFulfillmentFailed(log types.Log) (*VRFV2PlusWrapperWrapperFulfillmentFailed, error) {
	event := new(VRFV2PlusWrapperWrapperFulfillmentFailed)
	if err := _VRFV2PlusWrapper.contract.UnpackLog(event, "WrapperFulfillmentFailed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetConfig struct {
	KeyHash     [32]byte
	MaxNumWords uint8
}
type SCallbacks struct {
	CallbackAddress  common.Address
	CallbackGasLimit uint32
	RequestGasPrice  *big.Int
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapper) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFV2PlusWrapper.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFV2PlusWrapper.ParseOwnershipTransferRequested(log)
	case _VRFV2PlusWrapper.abi.Events["OwnershipTransferred"].ID:
		return _VRFV2PlusWrapper.ParseOwnershipTransferred(log)
	case _VRFV2PlusWrapper.abi.Events["WrapperFulfillmentFailed"].ID:
		return _VRFV2PlusWrapper.ParseWrapperFulfillmentFailed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFV2PlusWrapperOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFV2PlusWrapperOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VRFV2PlusWrapperWrapperFulfillmentFailed) Topic() common.Hash {
	return common.HexToHash("0xc551b83c151f2d1c7eeb938ac59008e0409f1c1dc1e2f112449d4d79b4589022")
}

func (_VRFV2PlusWrapper *VRFV2PlusWrapper) Address() common.Address {
	return _VRFV2PlusWrapper.address
}

type VRFV2PlusWrapperInterface interface {
	COORDINATOR(opts *bind.CallOpts) (common.Address, error)

	PRICEREGISTRY(opts *bind.CallOpts) (common.Address, error)

	SUBSCRIPTIONID(opts *bind.CallOpts) (*big.Int, error)

	GetConfig(opts *bind.CallOpts) (GetConfig,

		error)

	GetPriceRegistry(opts *bind.CallOpts) (common.Address, error)

	LastRequestId(opts *bind.CallOpts) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SCallbacks(opts *bind.CallOpts, arg0 *big.Int) (SCallbacks,

		error)

	SConfigured(opts *bind.CallOpts) (bool, error)

	SDisabled(opts *bind.CallOpts) (bool, error)

	SLink(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	Disable(opts *bind.TransactOpts) (*types.Transaction, error)

	Enable(opts *bind.TransactOpts) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, _sender common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	RequestRandomWordsInNative(opts *bind.TransactOpts, _callbackGasLimit uint32, _requestConfirmations uint16, _numWords uint32) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, _wrapperGasOverhead uint32, _keyHash [32]byte, _maxNumWords uint8) (*types.Transaction, error)

	SetCoordinator(opts *bind.TransactOpts, _vrfCoordinator common.Address) (*types.Transaction, error)

	SetLINK(opts *bind.TransactOpts, link common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Withdraw(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error)

	WithdrawNative(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusWrapperOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2PlusWrapperOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFV2PlusWrapperOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2PlusWrapperOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2PlusWrapperOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFV2PlusWrapperOwnershipTransferred, error)

	FilterWrapperFulfillmentFailed(opts *bind.FilterOpts, requestId []*big.Int, consumer []common.Address) (*VRFV2PlusWrapperWrapperFulfillmentFailedIterator, error)

	WatchWrapperFulfillmentFailed(opts *bind.WatchOpts, sink chan<- *VRFV2PlusWrapperWrapperFulfillmentFailed, requestId []*big.Int, consumer []common.Address) (event.Subscription, error)

	ParseWrapperFulfillmentFailed(log types.Log) (*VRFV2PlusWrapperWrapperFulfillmentFailed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
