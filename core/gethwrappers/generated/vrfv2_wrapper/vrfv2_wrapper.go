// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrfv2_wrapper

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

var VRFV2WrapperMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_linkEthFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_coordinator\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"WrapperFulfillmentFailed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COORDINATOR\",\"outputs\":[{\"internalType\":\"contractExtendedVRFCoordinatorV2Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_ETH_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SUBSCRIPTION_ID\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"}],\"name\":\"calculateRequestPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disable\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enable\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"_requestGasPriceWei\",\"type\":\"uint256\"}],\"name\":\"estimateRequestPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPM\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"wrapperGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"coordinatorGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"wrapperPremiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"maxNumWords\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRequestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_callbacks\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"callbackAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"requestGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"requestWeiPerUnitLink\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"juelsPaid\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_configured\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_disabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fulfillmentTxSizeBytes\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_wrapperGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_coordinatorGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"_wrapperPremiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"_maxNumWords\",\"type\":\"uint8\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"size\",\"type\":\"uint32\"}],\"name\":\"setFulfillmentTxSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101206040526001805463ffffffff60a01b1916609160a21b1790553480156200002857600080fd5b5060405162002a3c38038062002a3c8339810160408190526200004b91620002d8565b803380600081620000a35760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000d657620000d6816200020f565b5050506001600160601b0319606091821b811660805284821b811660a05283821b811660c0529082901b1660e0526040805163288688f960e21b815290516000916001600160a01b0384169163a21a23e49160048082019260209290919082900301818787803b1580156200014a57600080fd5b505af11580156200015f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000185919062000322565b60c081901b6001600160c01b03191661010052604051631cd0704360e21b81526001600160401b03821660048201523060248201529091506001600160a01b03831690637341c10c90604401600060405180830381600087803b158015620001ec57600080fd5b505af115801562000201573d6000803e3d6000fd5b505050505050505062000354565b6001600160a01b0381163314156200026a5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200009a565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620002d357600080fd5b919050565b600080600060608486031215620002ee57600080fd5b620002f984620002bb565b92506200030960208501620002bb565b91506200031960408501620002bb565b90509250925092565b6000602082840312156200033557600080fd5b81516001600160401b03811681146200034d57600080fd5b9392505050565b60805160601c60a05160601c60c05160601c60e05160601c6101005160c01c612655620003e7600039600081816101970152610c5701526000818161028401528181610c1801528181610fee015281816110cf015261116a0152600081816103fd015261161e01526000818161021b01528181610a6a01526112bc01526000818161055801526105c001526126556000f3fe608060405234801561001057600080fd5b506004361061018d5760003560e01c80638da5cb5b116100e3578063c15ce4d71161008c578063f2fde38b11610066578063f2fde38b14610511578063f3fef3a314610524578063fc2a88c31461053757600080fd5b8063c15ce4d714610432578063c3f909d414610445578063cdd8d885146104d457600080fd5b8063a608a1e1116100bd578063a608a1e1146103e6578063ad178361146103f8578063bf17e5591461041f57600080fd5b80638da5cb5b146103ad578063a3907d71146103cb578063a4c0ed36146103d357600080fd5b80633b2bcbf11161014557806357a8070a1161011f57806357a8070a1461037557806379ba5097146103925780637fb5d19d1461039a57600080fd5b80633b2bcbf11461027f5780634306d354146102a657806348baa1c5146102c757600080fd5b80631b6b6d23116101765780631b6b6d23146102165780631fe543e3146102625780632f2770db1461027757600080fd5b8063030932bb14610192578063181f5a77146101d7575b600080fd5b6101b97f000000000000000000000000000000000000000000000000000000000000000081565b60405167ffffffffffffffff90911681526020015b60405180910390f35b604080518082018252601281527f56524656325772617070657220312e302e300000000000000000000000000000602082015290516101ce91906122c1565b61023d7f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101ce565b610275610270366004611fa6565b610540565b005b610275610600565b61023d7f000000000000000000000000000000000000000000000000000000000000000081565b6102b96102b43660046120df565b610636565b6040519081526020016101ce565b6103316102d5366004611f8d565b600860205260009081526040902080546001820154600283015460039093015473ffffffffffffffffffffffffffffffffffffffff8316937401000000000000000000000000000000000000000090930463ffffffff16929085565b6040805173ffffffffffffffffffffffffffffffffffffffff909616865263ffffffff9094166020860152928401919091526060830152608082015260a0016101ce565b6003546103829060ff1681565b60405190151581526020016101ce565b61027561073d565b6102b96103a8366004612147565b61083a565b60005473ffffffffffffffffffffffffffffffffffffffff1661023d565b610275610942565b6102756103e1366004611e6c565b610974565b60035461038290610100900460ff1681565b61023d7f000000000000000000000000000000000000000000000000000000000000000081565b61027561042d3660046120df565b610e52565b61027561044036600461221b565b610ea9565b6004546005546006546007546040805194855263ffffffff80851660208701526401000000008504811691860191909152680100000000000000008404811660608601526c01000000000000000000000000840416608085015260ff700100000000000000000000000000000000909304831660a085015260c08401919091521660e0820152610100016101ce565b6001546104fc9074010000000000000000000000000000000000000000900463ffffffff1681565b60405163ffffffff90911681526020016101ce565b61027561051f366004611e27565b611254565b610275610532366004611e42565b611268565b6102b960025481565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146105f2576040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660248201526044015b60405180910390fd5b6105fc828261133d565b5050565b610608611548565b600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff16610100179055565b60035460009060ff166106a5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f77726170706572206973206e6f7420636f6e666967757265640000000000000060448201526064016105e9565b600354610100900460ff1615610717576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f777261707065722069732064697361626c65640000000000000000000000000060448201526064016105e9565b60006107216115cb565b90506107348363ffffffff163a8361173f565b9150505b919050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146107be576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016105e9565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60035460009060ff166108a9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f77726170706572206973206e6f7420636f6e666967757265640000000000000060448201526064016105e9565b600354610100900460ff161561091b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f777261707065722069732064697361626c65640000000000000000000000000060448201526064016105e9565b60006109256115cb565b90506109388463ffffffff16848361173f565b9150505b92915050565b61094a611548565b600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055565b60035460ff166109e0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f77726170706572206973206e6f7420636f6e666967757265640000000000000060448201526064016105e9565b600354610100900460ff1615610a52576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f777261707065722069732064697361626c65640000000000000000000000000060448201526064016105e9565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610af1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f6f6e6c792063616c6c61626c652066726f6d204c494e4b00000000000000000060448201526064016105e9565b60008080610b01848601866120fc565b9250925092506000610b1284611860565b90506000610b1e6115cb565b90506000610b338663ffffffff163a8461173f565b905080891015610b9f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f66656520746f6f206c6f7700000000000000000000000000000000000000000060448201526064016105e9565b60075460ff1663ffffffff85161115610c14576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f6e756d576f72647320746f6f206869676800000000000000000000000000000060448201526064016105e9565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16635d3b1d306006547f000000000000000000000000000000000000000000000000000000000000000089600560089054906101000a900463ffffffff16898d610c96919061239a565b610ca0919061239a565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e087901b168152600481019490945267ffffffffffffffff909216602484015261ffff16604483015263ffffffff90811660648301528816608482015260a401602060405180830381600087803b158015610d1f57600080fd5b505af1158015610d33573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d579190611f15565b90506040518060a001604052808c73ffffffffffffffffffffffffffffffffffffffff1681526020018863ffffffff1681526020013a81526020018481526020018b8152506008600083815260200190815260200160002060008201518160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060208201518160000160146101000a81548163ffffffff021916908363ffffffff160217905550604082015181600101556060820151816002015560808201518160030155905050806002819055505050505050505050505050565b610e5a611548565b6001805463ffffffff90921674010000000000000000000000000000000000000000027fffffffffffffffff00000000ffffffffffffffffffffffffffffffffffffffff909216919091179055565b610eb1611548565b6005805460ff808616700100000000000000000000000000000000027fffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffffff63ffffffff8981166c01000000000000000000000000027fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff918c166801000000000000000002919091167fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff909516949094179390931792909216919091179091556006839055600780549183167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00928316179055600380549091166001179055604080517fc3f909d4000000000000000000000000000000000000000000000000000000008152905173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169163c3f909d4916004828101926080929190829003018186803b15801561103457600080fd5b505afa158015611048573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061106c9190611f2e565b50600580547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff929092169190911790555050604080517f356dac7100000000000000000000000000000000000000000000000000000000815290517f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169163356dac71916004808301926020929190829003018186803b15801561112a57600080fd5b505afa15801561113e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111629190611f15565b6004819055507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16635fbbc0d26040518163ffffffff1660e01b81526004016101206040518083038186803b1580156111cf57600080fd5b505afa1580156111e3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112079190612165565b50506005805463ffffffff909816640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff909816979097179096555050505050505050505050565b61125c611548565b61126581611878565b50565b611270611548565b6040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8381166004830152602482018390527f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb90604401602060405180830381600087803b15801561130057600080fd5b505af1158015611314573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113389190611ef3565b505050565b6000828152600860208181526040808420815160a081018352815473ffffffffffffffffffffffffffffffffffffffff808216835263ffffffff740100000000000000000000000000000000000000008304168387015260018401805495840195909552600284018054606085015260038501805460808601528b8a52979096527fffffffffffffffff00000000000000000000000000000000000000000000000090911690925591859055918490559290915581511661145a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f72657175657374206e6f7420666f756e6400000000000000000000000000000060448201526064016105e9565b600080631fe543e360e01b8585604051602401611478929190612334565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050905060006114f2846020015163ffffffff1685600001518461196e565b90508061154057835160405173ffffffffffffffffffffffffffffffffffffffff9091169087907fc551b83c151f2d1c7eeb938ac59008e0409f1c1dc1e2f112449d4d79b458902290600090a35b505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146115c9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016105e9565b565b600554604080517ffeaf968c000000000000000000000000000000000000000000000000000000008152905160009263ffffffff161515918391829173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169163feaf968c9160048082019260a092909190829003018186803b15801561166557600080fd5b505afa158015611679573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061169d919061227d565b5094509092508491505080156116c357506116b88242612582565b60055463ffffffff16105b156116cd57506004545b6000811215611738576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f496e76616c6964204c494e4b207765692070726963650000000000000000000060448201526064016105e9565b9392505050565b600154600090819061176e9074010000000000000000000000000000000000000000900463ffffffff166119ba565b60055463ffffffff6c0100000000000000000000000082048116916117a191680100000000000000009091041688612382565b6117ab9190612382565b6117b59086612545565b6117bf9190612382565b90506000836117d683670de0b6b3a7640000612545565b6117e091906123e7565b60055490915060009060649061180d90700100000000000000000000000000000000900460ff16826123c2565b61181a9060ff1684612545565b61182491906123e7565b60055490915060009061184a90640100000000900463ffffffff1664e8d4a51000612545565b6118549083612382565b98975050505050505050565b600061186d603f836123fb565b61093c90600161239a565b73ffffffffffffffffffffffffffffffffffffffff81163314156118f8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016105e9565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60005a61138881101561198057600080fd5b61138881039050846040820482031161199857600080fd5b50823b6119a457600080fd5b60008083516020850160008789f1949350505050565b6000466119c681611a92565b15611a72576000606c73ffffffffffffffffffffffffffffffffffffffff166341b247a86040518163ffffffff1660e01b815260040160c06040518083038186803b158015611a1457600080fd5b505afa158015611a28573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611a4c9190612095565b5050505091505083608c611a609190612382565b611a6a9082612545565b949350505050565b611a7b81611ab5565b15611a895761073483611aef565b50600092915050565b600061a4b1821480611aa6575062066eed82145b8061093c57505062066eee1490565b6000600a821480611ac757506101a482145b80611ad4575062aa37dc82145b80611ae0575061210582145b8061093c57505062014a331490565b60008073420000000000000000000000000000000000000f73ffffffffffffffffffffffffffffffffffffffff1663519b4bd36040518163ffffffff1660e01b815260040160206040518083038186803b158015611b4c57600080fd5b505afa158015611b60573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611b849190611f15565b9050600080611b938186612582565b90506000611ba2826010612545565b611bad846004612545565b611bb79190612382565b9050600073420000000000000000000000000000000000000f73ffffffffffffffffffffffffffffffffffffffff16630c18c1626040518163ffffffff1660e01b815260040160206040518083038186803b158015611c1557600080fd5b505afa158015611c29573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611c4d9190611f15565b9050600073420000000000000000000000000000000000000f73ffffffffffffffffffffffffffffffffffffffff1663f45e65d86040518163ffffffff1660e01b815260040160206040518083038186803b158015611cab57600080fd5b505afa158015611cbf573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611ce39190611f15565b9050600073420000000000000000000000000000000000000f73ffffffffffffffffffffffffffffffffffffffff1663313ce5676040518163ffffffff1660e01b815260040160206040518083038186803b158015611d4157600080fd5b505afa158015611d55573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611d799190611f15565b90506000611d8882600a61247f565b905060008184611d988789612382565b611da2908c612545565b611dac9190612545565b611db691906123e7565b9b9a5050505050505050505050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461073857600080fd5b805162ffffff8116811461073857600080fd5b803560ff8116811461073857600080fd5b805169ffffffffffffffffffff8116811461073857600080fd5b600060208284031215611e3957600080fd5b61173882611dc5565b60008060408385031215611e5557600080fd5b611e5e83611dc5565b946020939093013593505050565b60008060008060608587031215611e8257600080fd5b611e8b85611dc5565b935060208501359250604085013567ffffffffffffffff80821115611eaf57600080fd5b818701915087601f830112611ec357600080fd5b813581811115611ed257600080fd5b886020828501011115611ee457600080fd5b95989497505060200194505050565b600060208284031215611f0557600080fd5b8151801515811461173857600080fd5b600060208284031215611f2757600080fd5b5051919050565b60008060008060808587031215611f4457600080fd5b8451611f4f81612626565b6020860151909450611f6081612636565b6040860151909350611f7181612636565b6060860151909250611f8281612636565b939692955090935050565b600060208284031215611f9f57600080fd5b5035919050565b60008060408385031215611fb957600080fd5b8235915060208084013567ffffffffffffffff80821115611fd957600080fd5b818601915086601f830112611fed57600080fd5b813581811115611fff57611fff6125f7565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f83011681018181108582111715612042576120426125f7565b604052828152858101935084860182860187018b101561206157600080fd5b600095505b83861015612084578035855260019590950194938601938601612066565b508096505050505050509250929050565b60008060008060008060c087890312156120ae57600080fd5b865195506020870151945060408701519350606087015192506080870151915060a087015190509295509295509295565b6000602082840312156120f157600080fd5b813561173881612636565b60008060006060848603121561211157600080fd5b833561211c81612636565b9250602084013561212c81612626565b9150604084013561213c81612636565b809150509250925092565b6000806040838503121561215a57600080fd5b8235611e5e81612636565b60008060008060008060008060006101208a8c03121561218457600080fd5b895161218f81612636565b60208b01519099506121a081612636565b60408b01519098506121b181612636565b60608b01519097506121c281612636565b60808b01519096506121d381612636565b94506121e160a08b01611de9565b93506121ef60c08b01611de9565b92506121fd60e08b01611de9565b915061220c6101008b01611de9565b90509295985092959850929598565b600080600080600060a0868803121561223357600080fd5b853561223e81612636565b9450602086013561224e81612636565b935061225c60408701611dfc565b92506060860135915061227160808701611dfc565b90509295509295909350565b600080600080600060a0868803121561229557600080fd5b61229e86611e0d565b945060208601519350604086015192506060860151915061227160808701611e0d565b600060208083528351808285015260005b818110156122ee578581018301518582016040015282016122d2565b81811115612300576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b6000604082018483526020604081850152818551808452606086019150828701935060005b8181101561237557845183529383019391830191600101612359565b5090979650505050505050565b6000821982111561239557612395612599565b500190565b600063ffffffff8083168185168083038211156123b9576123b9612599565b01949350505050565b600060ff821660ff84168060ff038211156123df576123df612599565b019392505050565b6000826123f6576123f66125c8565b500490565b600063ffffffff80841680612412576124126125c8565b92169190910492915050565b600181815b8085111561247757817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0482111561245d5761245d612599565b8085161561246a57918102915b93841c9390800290612423565b509250929050565b600061173883836000826124955750600161093c565b816124a25750600061093c565b81600181146124b857600281146124c2576124de565b600191505061093c565b60ff8411156124d3576124d3612599565b50506001821b61093c565b5060208310610133831016604e8410600b8410161715612501575081810a61093c565b61250b838361241e565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0482111561253d5761253d612599565b029392505050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561257d5761257d612599565b500290565b60008282101561259457612594612599565b500390565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b61ffff8116811461126557600080fd5b63ffffffff8116811461126557600080fdfea164736f6c6343000806000a",
}

var VRFV2WrapperABI = VRFV2WrapperMetaData.ABI

var VRFV2WrapperBin = VRFV2WrapperMetaData.Bin

func DeployVRFV2Wrapper(auth *bind.TransactOpts, backend bind.ContractBackend, _link common.Address, _linkEthFeed common.Address, _coordinator common.Address) (common.Address, *types.Transaction, *VRFV2Wrapper, error) {
	parsed, err := VRFV2WrapperMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2WrapperBin), backend, _link, _linkEthFeed, _coordinator)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2Wrapper{address: address, abi: *parsed, VRFV2WrapperCaller: VRFV2WrapperCaller{contract: contract}, VRFV2WrapperTransactor: VRFV2WrapperTransactor{contract: contract}, VRFV2WrapperFilterer: VRFV2WrapperFilterer{contract: contract}}, nil
}

type VRFV2Wrapper struct {
	address common.Address
	abi     abi.ABI
	VRFV2WrapperCaller
	VRFV2WrapperTransactor
	VRFV2WrapperFilterer
}

type VRFV2WrapperCaller struct {
	contract *bind.BoundContract
}

type VRFV2WrapperTransactor struct {
	contract *bind.BoundContract
}

type VRFV2WrapperFilterer struct {
	contract *bind.BoundContract
}

type VRFV2WrapperSession struct {
	Contract     *VRFV2Wrapper
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2WrapperCallerSession struct {
	Contract *VRFV2WrapperCaller
	CallOpts bind.CallOpts
}

type VRFV2WrapperTransactorSession struct {
	Contract     *VRFV2WrapperTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2WrapperRaw struct {
	Contract *VRFV2Wrapper
}

type VRFV2WrapperCallerRaw struct {
	Contract *VRFV2WrapperCaller
}

type VRFV2WrapperTransactorRaw struct {
	Contract *VRFV2WrapperTransactor
}

func NewVRFV2Wrapper(address common.Address, backend bind.ContractBackend) (*VRFV2Wrapper, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2WrapperABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2Wrapper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2Wrapper{address: address, abi: abi, VRFV2WrapperCaller: VRFV2WrapperCaller{contract: contract}, VRFV2WrapperTransactor: VRFV2WrapperTransactor{contract: contract}, VRFV2WrapperFilterer: VRFV2WrapperFilterer{contract: contract}}, nil
}

func NewVRFV2WrapperCaller(address common.Address, caller bind.ContractCaller) (*VRFV2WrapperCaller, error) {
	contract, err := bindVRFV2Wrapper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2WrapperCaller{contract: contract}, nil
}

func NewVRFV2WrapperTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2WrapperTransactor, error) {
	contract, err := bindVRFV2Wrapper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2WrapperTransactor{contract: contract}, nil
}

func NewVRFV2WrapperFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2WrapperFilterer, error) {
	contract, err := bindVRFV2Wrapper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2WrapperFilterer{contract: contract}, nil
}

func bindVRFV2Wrapper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFV2WrapperMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFV2Wrapper *VRFV2WrapperRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2Wrapper.Contract.VRFV2WrapperCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2Wrapper *VRFV2WrapperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2Wrapper.Contract.VRFV2WrapperTransactor.contract.Transfer(opts)
}

func (_VRFV2Wrapper *VRFV2WrapperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2Wrapper.Contract.VRFV2WrapperTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2Wrapper *VRFV2WrapperCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2Wrapper.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2Wrapper *VRFV2WrapperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2Wrapper.Contract.contract.Transfer(opts)
}

func (_VRFV2Wrapper *VRFV2WrapperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2Wrapper.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2Wrapper *VRFV2WrapperCaller) COORDINATOR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2Wrapper.contract.Call(opts, &out, "COORDINATOR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2Wrapper *VRFV2WrapperSession) COORDINATOR() (common.Address, error) {
	return _VRFV2Wrapper.Contract.COORDINATOR(&_VRFV2Wrapper.CallOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperCallerSession) COORDINATOR() (common.Address, error) {
	return _VRFV2Wrapper.Contract.COORDINATOR(&_VRFV2Wrapper.CallOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2Wrapper.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2Wrapper *VRFV2WrapperSession) LINK() (common.Address, error) {
	return _VRFV2Wrapper.Contract.LINK(&_VRFV2Wrapper.CallOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperCallerSession) LINK() (common.Address, error) {
	return _VRFV2Wrapper.Contract.LINK(&_VRFV2Wrapper.CallOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperCaller) LINKETHFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2Wrapper.contract.Call(opts, &out, "LINK_ETH_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2Wrapper *VRFV2WrapperSession) LINKETHFEED() (common.Address, error) {
	return _VRFV2Wrapper.Contract.LINKETHFEED(&_VRFV2Wrapper.CallOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperCallerSession) LINKETHFEED() (common.Address, error) {
	return _VRFV2Wrapper.Contract.LINKETHFEED(&_VRFV2Wrapper.CallOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperCaller) SUBSCRIPTIONID(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFV2Wrapper.contract.Call(opts, &out, "SUBSCRIPTION_ID")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFV2Wrapper *VRFV2WrapperSession) SUBSCRIPTIONID() (uint64, error) {
	return _VRFV2Wrapper.Contract.SUBSCRIPTIONID(&_VRFV2Wrapper.CallOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperCallerSession) SUBSCRIPTIONID() (uint64, error) {
	return _VRFV2Wrapper.Contract.SUBSCRIPTIONID(&_VRFV2Wrapper.CallOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperCaller) CalculateRequestPrice(opts *bind.CallOpts, _callbackGasLimit uint32) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2Wrapper.contract.Call(opts, &out, "calculateRequestPrice", _callbackGasLimit)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2Wrapper *VRFV2WrapperSession) CalculateRequestPrice(_callbackGasLimit uint32) (*big.Int, error) {
	return _VRFV2Wrapper.Contract.CalculateRequestPrice(&_VRFV2Wrapper.CallOpts, _callbackGasLimit)
}

func (_VRFV2Wrapper *VRFV2WrapperCallerSession) CalculateRequestPrice(_callbackGasLimit uint32) (*big.Int, error) {
	return _VRFV2Wrapper.Contract.CalculateRequestPrice(&_VRFV2Wrapper.CallOpts, _callbackGasLimit)
}

func (_VRFV2Wrapper *VRFV2WrapperCaller) EstimateRequestPrice(opts *bind.CallOpts, _callbackGasLimit uint32, _requestGasPriceWei *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2Wrapper.contract.Call(opts, &out, "estimateRequestPrice", _callbackGasLimit, _requestGasPriceWei)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2Wrapper *VRFV2WrapperSession) EstimateRequestPrice(_callbackGasLimit uint32, _requestGasPriceWei *big.Int) (*big.Int, error) {
	return _VRFV2Wrapper.Contract.EstimateRequestPrice(&_VRFV2Wrapper.CallOpts, _callbackGasLimit, _requestGasPriceWei)
}

func (_VRFV2Wrapper *VRFV2WrapperCallerSession) EstimateRequestPrice(_callbackGasLimit uint32, _requestGasPriceWei *big.Int) (*big.Int, error) {
	return _VRFV2Wrapper.Contract.EstimateRequestPrice(&_VRFV2Wrapper.CallOpts, _callbackGasLimit, _requestGasPriceWei)
}

func (_VRFV2Wrapper *VRFV2WrapperCaller) GetConfig(opts *bind.CallOpts) (GetConfig,

	error) {
	var out []interface{}
	err := _VRFV2Wrapper.contract.Call(opts, &out, "getConfig")

	outstruct := new(GetConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.FallbackWeiPerUnitLink = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.StalenessSeconds = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.FulfillmentFlatFeeLinkPPM = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.WrapperGasOverhead = *abi.ConvertType(out[3], new(uint32)).(*uint32)
	outstruct.CoordinatorGasOverhead = *abi.ConvertType(out[4], new(uint32)).(*uint32)
	outstruct.WrapperPremiumPercentage = *abi.ConvertType(out[5], new(uint8)).(*uint8)
	outstruct.KeyHash = *abi.ConvertType(out[6], new([32]byte)).(*[32]byte)
	outstruct.MaxNumWords = *abi.ConvertType(out[7], new(uint8)).(*uint8)

	return *outstruct, err

}

func (_VRFV2Wrapper *VRFV2WrapperSession) GetConfig() (GetConfig,

	error) {
	return _VRFV2Wrapper.Contract.GetConfig(&_VRFV2Wrapper.CallOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperCallerSession) GetConfig() (GetConfig,

	error) {
	return _VRFV2Wrapper.Contract.GetConfig(&_VRFV2Wrapper.CallOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperCaller) LastRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFV2Wrapper.contract.Call(opts, &out, "lastRequestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFV2Wrapper *VRFV2WrapperSession) LastRequestId() (*big.Int, error) {
	return _VRFV2Wrapper.Contract.LastRequestId(&_VRFV2Wrapper.CallOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperCallerSession) LastRequestId() (*big.Int, error) {
	return _VRFV2Wrapper.Contract.LastRequestId(&_VRFV2Wrapper.CallOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFV2Wrapper.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFV2Wrapper *VRFV2WrapperSession) Owner() (common.Address, error) {
	return _VRFV2Wrapper.Contract.Owner(&_VRFV2Wrapper.CallOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperCallerSession) Owner() (common.Address, error) {
	return _VRFV2Wrapper.Contract.Owner(&_VRFV2Wrapper.CallOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperCaller) SCallbacks(opts *bind.CallOpts, arg0 *big.Int) (SCallbacks,

	error) {
	var out []interface{}
	err := _VRFV2Wrapper.contract.Call(opts, &out, "s_callbacks", arg0)

	outstruct := new(SCallbacks)
	if err != nil {
		return *outstruct, err
	}

	outstruct.CallbackAddress = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.CallbackGasLimit = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.RequestGasPrice = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.RequestWeiPerUnitLink = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.JuelsPaid = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_VRFV2Wrapper *VRFV2WrapperSession) SCallbacks(arg0 *big.Int) (SCallbacks,

	error) {
	return _VRFV2Wrapper.Contract.SCallbacks(&_VRFV2Wrapper.CallOpts, arg0)
}

func (_VRFV2Wrapper *VRFV2WrapperCallerSession) SCallbacks(arg0 *big.Int) (SCallbacks,

	error) {
	return _VRFV2Wrapper.Contract.SCallbacks(&_VRFV2Wrapper.CallOpts, arg0)
}

func (_VRFV2Wrapper *VRFV2WrapperCaller) SConfigured(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _VRFV2Wrapper.contract.Call(opts, &out, "s_configured")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VRFV2Wrapper *VRFV2WrapperSession) SConfigured() (bool, error) {
	return _VRFV2Wrapper.Contract.SConfigured(&_VRFV2Wrapper.CallOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperCallerSession) SConfigured() (bool, error) {
	return _VRFV2Wrapper.Contract.SConfigured(&_VRFV2Wrapper.CallOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperCaller) SDisabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _VRFV2Wrapper.contract.Call(opts, &out, "s_disabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VRFV2Wrapper *VRFV2WrapperSession) SDisabled() (bool, error) {
	return _VRFV2Wrapper.Contract.SDisabled(&_VRFV2Wrapper.CallOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperCallerSession) SDisabled() (bool, error) {
	return _VRFV2Wrapper.Contract.SDisabled(&_VRFV2Wrapper.CallOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperCaller) SFulfillmentTxSizeBytes(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _VRFV2Wrapper.contract.Call(opts, &out, "s_fulfillmentTxSizeBytes")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_VRFV2Wrapper *VRFV2WrapperSession) SFulfillmentTxSizeBytes() (uint32, error) {
	return _VRFV2Wrapper.Contract.SFulfillmentTxSizeBytes(&_VRFV2Wrapper.CallOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperCallerSession) SFulfillmentTxSizeBytes() (uint32, error) {
	return _VRFV2Wrapper.Contract.SFulfillmentTxSizeBytes(&_VRFV2Wrapper.CallOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VRFV2Wrapper.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VRFV2Wrapper *VRFV2WrapperSession) TypeAndVersion() (string, error) {
	return _VRFV2Wrapper.Contract.TypeAndVersion(&_VRFV2Wrapper.CallOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperCallerSession) TypeAndVersion() (string, error) {
	return _VRFV2Wrapper.Contract.TypeAndVersion(&_VRFV2Wrapper.CallOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2Wrapper.contract.Transact(opts, "acceptOwnership")
}

func (_VRFV2Wrapper *VRFV2WrapperSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2Wrapper.Contract.AcceptOwnership(&_VRFV2Wrapper.TransactOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFV2Wrapper.Contract.AcceptOwnership(&_VRFV2Wrapper.TransactOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperTransactor) Disable(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2Wrapper.contract.Transact(opts, "disable")
}

func (_VRFV2Wrapper *VRFV2WrapperSession) Disable() (*types.Transaction, error) {
	return _VRFV2Wrapper.Contract.Disable(&_VRFV2Wrapper.TransactOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperTransactorSession) Disable() (*types.Transaction, error) {
	return _VRFV2Wrapper.Contract.Disable(&_VRFV2Wrapper.TransactOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperTransactor) Enable(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2Wrapper.contract.Transact(opts, "enable")
}

func (_VRFV2Wrapper *VRFV2WrapperSession) Enable() (*types.Transaction, error) {
	return _VRFV2Wrapper.Contract.Enable(&_VRFV2Wrapper.TransactOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperTransactorSession) Enable() (*types.Transaction, error) {
	return _VRFV2Wrapper.Contract.Enable(&_VRFV2Wrapper.TransactOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperTransactor) OnTokenTransfer(opts *bind.TransactOpts, _sender common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _VRFV2Wrapper.contract.Transact(opts, "onTokenTransfer", _sender, _amount, _data)
}

func (_VRFV2Wrapper *VRFV2WrapperSession) OnTokenTransfer(_sender common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _VRFV2Wrapper.Contract.OnTokenTransfer(&_VRFV2Wrapper.TransactOpts, _sender, _amount, _data)
}

func (_VRFV2Wrapper *VRFV2WrapperTransactorSession) OnTokenTransfer(_sender common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _VRFV2Wrapper.Contract.OnTokenTransfer(&_VRFV2Wrapper.TransactOpts, _sender, _amount, _data)
}

func (_VRFV2Wrapper *VRFV2WrapperTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2Wrapper.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

func (_VRFV2Wrapper *VRFV2WrapperSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2Wrapper.Contract.RawFulfillRandomWords(&_VRFV2Wrapper.TransactOpts, requestId, randomWords)
}

func (_VRFV2Wrapper *VRFV2WrapperTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFV2Wrapper.Contract.RawFulfillRandomWords(&_VRFV2Wrapper.TransactOpts, requestId, randomWords)
}

func (_VRFV2Wrapper *VRFV2WrapperTransactor) SetConfig(opts *bind.TransactOpts, _wrapperGasOverhead uint32, _coordinatorGasOverhead uint32, _wrapperPremiumPercentage uint8, _keyHash [32]byte, _maxNumWords uint8) (*types.Transaction, error) {
	return _VRFV2Wrapper.contract.Transact(opts, "setConfig", _wrapperGasOverhead, _coordinatorGasOverhead, _wrapperPremiumPercentage, _keyHash, _maxNumWords)
}

func (_VRFV2Wrapper *VRFV2WrapperSession) SetConfig(_wrapperGasOverhead uint32, _coordinatorGasOverhead uint32, _wrapperPremiumPercentage uint8, _keyHash [32]byte, _maxNumWords uint8) (*types.Transaction, error) {
	return _VRFV2Wrapper.Contract.SetConfig(&_VRFV2Wrapper.TransactOpts, _wrapperGasOverhead, _coordinatorGasOverhead, _wrapperPremiumPercentage, _keyHash, _maxNumWords)
}

func (_VRFV2Wrapper *VRFV2WrapperTransactorSession) SetConfig(_wrapperGasOverhead uint32, _coordinatorGasOverhead uint32, _wrapperPremiumPercentage uint8, _keyHash [32]byte, _maxNumWords uint8) (*types.Transaction, error) {
	return _VRFV2Wrapper.Contract.SetConfig(&_VRFV2Wrapper.TransactOpts, _wrapperGasOverhead, _coordinatorGasOverhead, _wrapperPremiumPercentage, _keyHash, _maxNumWords)
}

func (_VRFV2Wrapper *VRFV2WrapperTransactor) SetFulfillmentTxSize(opts *bind.TransactOpts, size uint32) (*types.Transaction, error) {
	return _VRFV2Wrapper.contract.Transact(opts, "setFulfillmentTxSize", size)
}

func (_VRFV2Wrapper *VRFV2WrapperSession) SetFulfillmentTxSize(size uint32) (*types.Transaction, error) {
	return _VRFV2Wrapper.Contract.SetFulfillmentTxSize(&_VRFV2Wrapper.TransactOpts, size)
}

func (_VRFV2Wrapper *VRFV2WrapperTransactorSession) SetFulfillmentTxSize(size uint32) (*types.Transaction, error) {
	return _VRFV2Wrapper.Contract.SetFulfillmentTxSize(&_VRFV2Wrapper.TransactOpts, size)
}

func (_VRFV2Wrapper *VRFV2WrapperTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFV2Wrapper.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFV2Wrapper *VRFV2WrapperSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2Wrapper.Contract.TransferOwnership(&_VRFV2Wrapper.TransactOpts, to)
}

func (_VRFV2Wrapper *VRFV2WrapperTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFV2Wrapper.Contract.TransferOwnership(&_VRFV2Wrapper.TransactOpts, to)
}

func (_VRFV2Wrapper *VRFV2WrapperTransactor) Withdraw(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _VRFV2Wrapper.contract.Transact(opts, "withdraw", _recipient, _amount)
}

func (_VRFV2Wrapper *VRFV2WrapperSession) Withdraw(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _VRFV2Wrapper.Contract.Withdraw(&_VRFV2Wrapper.TransactOpts, _recipient, _amount)
}

func (_VRFV2Wrapper *VRFV2WrapperTransactorSession) Withdraw(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _VRFV2Wrapper.Contract.Withdraw(&_VRFV2Wrapper.TransactOpts, _recipient, _amount)
}

type VRFV2WrapperOwnershipTransferRequestedIterator struct {
	Event *VRFV2WrapperOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2WrapperOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2WrapperOwnershipTransferRequested)
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
		it.Event = new(VRFV2WrapperOwnershipTransferRequested)
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

func (it *VRFV2WrapperOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFV2WrapperOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2WrapperOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2Wrapper *VRFV2WrapperFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2WrapperOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2Wrapper.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2WrapperOwnershipTransferRequestedIterator{contract: _VRFV2Wrapper.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFV2Wrapper *VRFV2WrapperFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2WrapperOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2Wrapper.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2WrapperOwnershipTransferRequested)
				if err := _VRFV2Wrapper.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFV2Wrapper *VRFV2WrapperFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFV2WrapperOwnershipTransferRequested, error) {
	event := new(VRFV2WrapperOwnershipTransferRequested)
	if err := _VRFV2Wrapper.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2WrapperOwnershipTransferredIterator struct {
	Event *VRFV2WrapperOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2WrapperOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2WrapperOwnershipTransferred)
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
		it.Event = new(VRFV2WrapperOwnershipTransferred)
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

func (it *VRFV2WrapperOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFV2WrapperOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2WrapperOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFV2Wrapper *VRFV2WrapperFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2WrapperOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2Wrapper.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2WrapperOwnershipTransferredIterator{contract: _VRFV2Wrapper.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFV2Wrapper *VRFV2WrapperFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2WrapperOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFV2Wrapper.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2WrapperOwnershipTransferred)
				if err := _VRFV2Wrapper.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFV2Wrapper *VRFV2WrapperFilterer) ParseOwnershipTransferred(log types.Log) (*VRFV2WrapperOwnershipTransferred, error) {
	event := new(VRFV2WrapperOwnershipTransferred)
	if err := _VRFV2Wrapper.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFV2WrapperWrapperFulfillmentFailedIterator struct {
	Event *VRFV2WrapperWrapperFulfillmentFailed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFV2WrapperWrapperFulfillmentFailedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFV2WrapperWrapperFulfillmentFailed)
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
		it.Event = new(VRFV2WrapperWrapperFulfillmentFailed)
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

func (it *VRFV2WrapperWrapperFulfillmentFailedIterator) Error() error {
	return it.fail
}

func (it *VRFV2WrapperWrapperFulfillmentFailedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFV2WrapperWrapperFulfillmentFailed struct {
	RequestId *big.Int
	Consumer  common.Address
	Raw       types.Log
}

func (_VRFV2Wrapper *VRFV2WrapperFilterer) FilterWrapperFulfillmentFailed(opts *bind.FilterOpts, requestId []*big.Int, consumer []common.Address) (*VRFV2WrapperWrapperFulfillmentFailedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var consumerRule []interface{}
	for _, consumerItem := range consumer {
		consumerRule = append(consumerRule, consumerItem)
	}

	logs, sub, err := _VRFV2Wrapper.contract.FilterLogs(opts, "WrapperFulfillmentFailed", requestIdRule, consumerRule)
	if err != nil {
		return nil, err
	}
	return &VRFV2WrapperWrapperFulfillmentFailedIterator{contract: _VRFV2Wrapper.contract, event: "WrapperFulfillmentFailed", logs: logs, sub: sub}, nil
}

func (_VRFV2Wrapper *VRFV2WrapperFilterer) WatchWrapperFulfillmentFailed(opts *bind.WatchOpts, sink chan<- *VRFV2WrapperWrapperFulfillmentFailed, requestId []*big.Int, consumer []common.Address) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var consumerRule []interface{}
	for _, consumerItem := range consumer {
		consumerRule = append(consumerRule, consumerItem)
	}

	logs, sub, err := _VRFV2Wrapper.contract.WatchLogs(opts, "WrapperFulfillmentFailed", requestIdRule, consumerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFV2WrapperWrapperFulfillmentFailed)
				if err := _VRFV2Wrapper.contract.UnpackLog(event, "WrapperFulfillmentFailed", log); err != nil {
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

func (_VRFV2Wrapper *VRFV2WrapperFilterer) ParseWrapperFulfillmentFailed(log types.Log) (*VRFV2WrapperWrapperFulfillmentFailed, error) {
	event := new(VRFV2WrapperWrapperFulfillmentFailed)
	if err := _VRFV2Wrapper.contract.UnpackLog(event, "WrapperFulfillmentFailed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetConfig struct {
	FallbackWeiPerUnitLink    *big.Int
	StalenessSeconds          uint32
	FulfillmentFlatFeeLinkPPM uint32
	WrapperGasOverhead        uint32
	CoordinatorGasOverhead    uint32
	WrapperPremiumPercentage  uint8
	KeyHash                   [32]byte
	MaxNumWords               uint8
}
type SCallbacks struct {
	CallbackAddress       common.Address
	CallbackGasLimit      uint32
	RequestGasPrice       *big.Int
	RequestWeiPerUnitLink *big.Int
	JuelsPaid             *big.Int
}

func (_VRFV2Wrapper *VRFV2Wrapper) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFV2Wrapper.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFV2Wrapper.ParseOwnershipTransferRequested(log)
	case _VRFV2Wrapper.abi.Events["OwnershipTransferred"].ID:
		return _VRFV2Wrapper.ParseOwnershipTransferred(log)
	case _VRFV2Wrapper.abi.Events["WrapperFulfillmentFailed"].ID:
		return _VRFV2Wrapper.ParseWrapperFulfillmentFailed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFV2WrapperOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFV2WrapperOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VRFV2WrapperWrapperFulfillmentFailed) Topic() common.Hash {
	return common.HexToHash("0xc551b83c151f2d1c7eeb938ac59008e0409f1c1dc1e2f112449d4d79b4589022")
}

func (_VRFV2Wrapper *VRFV2Wrapper) Address() common.Address {
	return _VRFV2Wrapper.address
}

type VRFV2WrapperInterface interface {
	COORDINATOR(opts *bind.CallOpts) (common.Address, error)

	LINK(opts *bind.CallOpts) (common.Address, error)

	LINKETHFEED(opts *bind.CallOpts) (common.Address, error)

	SUBSCRIPTIONID(opts *bind.CallOpts) (uint64, error)

	CalculateRequestPrice(opts *bind.CallOpts, _callbackGasLimit uint32) (*big.Int, error)

	EstimateRequestPrice(opts *bind.CallOpts, _callbackGasLimit uint32, _requestGasPriceWei *big.Int) (*big.Int, error)

	GetConfig(opts *bind.CallOpts) (GetConfig,

		error)

	LastRequestId(opts *bind.CallOpts) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SCallbacks(opts *bind.CallOpts, arg0 *big.Int) (SCallbacks,

		error)

	SConfigured(opts *bind.CallOpts) (bool, error)

	SDisabled(opts *bind.CallOpts) (bool, error)

	SFulfillmentTxSizeBytes(opts *bind.CallOpts) (uint32, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	Disable(opts *bind.TransactOpts) (*types.Transaction, error)

	Enable(opts *bind.TransactOpts) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, _sender common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, _wrapperGasOverhead uint32, _coordinatorGasOverhead uint32, _wrapperPremiumPercentage uint8, _keyHash [32]byte, _maxNumWords uint8) (*types.Transaction, error)

	SetFulfillmentTxSize(opts *bind.TransactOpts, size uint32) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Withdraw(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2WrapperOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFV2WrapperOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFV2WrapperOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFV2WrapperOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFV2WrapperOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFV2WrapperOwnershipTransferred, error)

	FilterWrapperFulfillmentFailed(opts *bind.FilterOpts, requestId []*big.Int, consumer []common.Address) (*VRFV2WrapperWrapperFulfillmentFailedIterator, error)

	WatchWrapperFulfillmentFailed(opts *bind.WatchOpts, sink chan<- *VRFV2WrapperWrapperFulfillmentFailed, requestId []*big.Int, consumer []common.Address) (event.Subscription, error)

	ParseWrapperFulfillmentFailed(log types.Log) (*VRFV2WrapperWrapperFulfillmentFailed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
