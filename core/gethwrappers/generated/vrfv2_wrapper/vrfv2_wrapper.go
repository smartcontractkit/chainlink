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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_linkEthFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_coordinator\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"WrapperFulfillmentFailed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COORDINATOR\",\"outputs\":[{\"internalType\":\"contractExtendedVRFCoordinatorV2Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FULFILLMENT_COMPRESSED_TX_DATA_SIZE_BYTES\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_ETH_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SUBSCRIPTION_ID\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"}],\"name\":\"calculateRequestPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disable\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enable\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"_requestGasPriceWei\",\"type\":\"uint256\"}],\"name\":\"estimateRequestPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPM\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"wrapperGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"coordinatorGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"wrapperPremiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"maxNumWords\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRequestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_callbacks\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"callbackAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"requestGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"requestWeiPerUnitLink\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"juelsPaid\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_configured\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_disabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_wrapperGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_coordinatorGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"_wrapperPremiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"_keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"_maxNumWords\",\"type\":\"uint8\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101206040523480156200001257600080fd5b506040516200256a3803806200256a8339810160408190526200003591620002c2565b8033806000816200008d5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c057620000c081620001f9565b5050506001600160601b0319606091821b811660805284821b811660a05283821b811660c0529082901b1660e0526040805163288688f960e21b815290516000916001600160a01b0384169163a21a23e49160048082019260209290919082900301818787803b1580156200013457600080fd5b505af115801562000149573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200016f91906200030c565b60c081901b6001600160c01b03191661010052604051631cd0704360e21b81526001600160401b03821660048201523060248201529091506001600160a01b03831690637341c10c90604401600060405180830381600087803b158015620001d657600080fd5b505af1158015620001eb573d6000803e3d6000fd5b50505050505050506200033e565b6001600160a01b038116331415620002545760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000084565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620002bd57600080fd5b919050565b600080600060608486031215620002d857600080fd5b620002e384620002a5565b9250620002f360208501620002a5565b91506200030360408501620002a5565b90509250925092565b6000602082840312156200031f57600080fd5b81516001600160401b03811681146200033757600080fd5b9392505050565b60805160601c60a05160601c60c05160601c60e05160601c6101005160c01c612199620003d16000396000818161018c0152610c1801526000818161027901528181610bd901528181610f580152818161103901526110d4015260008181610410015261158801526000818161021001528181610a2b015261122601526000818161051b015261058301526121996000f3fe608060405234801561001057600080fd5b50600436106101825760003560e01c80637fb5d19d116100d8578063ad1783611161008c578063f2fde38b11610066578063f2fde38b146104d4578063f3fef3a3146104e7578063fc2a88c3146104fa57600080fd5b8063ad1783611461040b578063c15ce4d714610432578063c3f909d41461044557600080fd5b8063a3907d71116100bd578063a3907d71146103de578063a4c0ed36146103e6578063a608a1e1146103f957600080fd5b80637fb5d19d146103ad5780638da5cb5b146103c057600080fd5b80633b2bcbf11161013a57806357a8070a1161011457806357a8070a1461036a578063607f7ec01461038757806379ba5097146103a557600080fd5b80633b2bcbf1146102745780634306d3541461029b57806348baa1c5146102bc57600080fd5b80631b6b6d231161016b5780631b6b6d231461020b5780631fe543e3146102575780632f2770db1461026c57600080fd5b8063030932bb14610187578063181f5a77146101cc575b600080fd5b6101ae7f000000000000000000000000000000000000000000000000000000000000000081565b60405167ffffffffffffffff90911681526020015b60405180910390f35b604080518082018252601281527f56524656325772617070657220312e302e300000000000000000000000000000602082015290516101c39190611f2c565b6102327f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101c3565b61026a610265366004611c11565b610503565b005b61026a6105c3565b6102327f000000000000000000000000000000000000000000000000000000000000000081565b6102ae6102a9366004611d4a565b6105f9565b6040519081526020016101c3565b6103266102ca366004611bf8565b600860205260009081526040902080546001820154600283015460039093015473ffffffffffffffffffffffffffffffffffffffff8316937401000000000000000000000000000000000000000090930463ffffffff16929085565b6040805173ffffffffffffffffffffffffffffffffffffffff909616865263ffffffff9094166020860152928401919091526060830152608082015260a0016101c3565b6003546103779060ff1681565b60405190151581526020016101c3565b6103906101f481565b60405163ffffffff90911681526020016101c3565b61026a610700565b6102ae6103bb366004611db2565b6107fd565b60005473ffffffffffffffffffffffffffffffffffffffff16610232565b61026a610903565b61026a6103f4366004611ad7565b610935565b60035461037790610100900460ff1681565b6102327f000000000000000000000000000000000000000000000000000000000000000081565b61026a610440366004611e86565b610e13565b6004546005546006546007546040805194855263ffffffff80851660208701526401000000008504811691860191909152680100000000000000008404811660608601526c01000000000000000000000000840416608085015260ff700100000000000000000000000000000000909304831660a085015260c08401919091521660e0820152610100016101c3565b61026a6104e2366004611a92565b6111be565b61026a6104f5366004611aad565b6111d2565b6102ae60025481565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146105b5576040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660248201526044015b60405180910390fd5b6105bf82826112a7565b5050565b6105cb6114b2565b600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff16610100179055565b60035460009060ff16610668576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f77726170706572206973206e6f7420636f6e666967757265640000000000000060448201526064016105ac565b600354610100900460ff16156106da576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f777261707065722069732064697361626c65640000000000000000000000000060448201526064016105ac565b60006106e4611535565b90506106f78363ffffffff163a836116a9565b9150505b919050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610781576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016105ac565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60035460009060ff1661086c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f77726170706572206973206e6f7420636f6e666967757265640000000000000060448201526064016105ac565b600354610100900460ff16156108de576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f777261707065722069732064697361626c65640000000000000000000000000060448201526064016105ac565b60006108e8611535565b90506108fb8463ffffffff1684836116a9565b949350505050565b61090b6114b2565b600380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055565b60035460ff166109a1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f77726170706572206973206e6f7420636f6e666967757265640000000000000060448201526064016105ac565b600354610100900460ff1615610a13576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f777261707065722069732064697361626c65640000000000000000000000000060448201526064016105ac565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610ab2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f6f6e6c792063616c6c61626c652066726f6d204c494e4b00000000000000000060448201526064016105ac565b60008080610ac284860186611d67565b9250925092506000610ad3846117d7565b90506000610adf611535565b90506000610af48663ffffffff163a846116a9565b905080891015610b60576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f66656520746f6f206c6f7700000000000000000000000000000000000000000060448201526064016105ac565b60075460ff1663ffffffff85161115610bd5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f6e756d576f72647320746f6f206869676800000000000000000000000000000060448201526064016105ac565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16635d3b1d306006547f000000000000000000000000000000000000000000000000000000000000000089600560089054906101000a900463ffffffff16898d610c579190612005565b610c619190612005565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e087901b168152600481019490945267ffffffffffffffff909216602484015261ffff16604483015263ffffffff90811660648301528816608482015260a401602060405180830381600087803b158015610ce057600080fd5b505af1158015610cf4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d189190611b80565b90506040518060a001604052808c73ffffffffffffffffffffffffffffffffffffffff1681526020018863ffffffff1681526020013a81526020018481526020018b8152506008600083815260200190815260200160002060008201518160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060208201518160000160146101000a81548163ffffffff021916908363ffffffff160217905550604082015181600101556060820151816002015560808201518160030155905050806002819055505050505050505050505050565b610e1b6114b2565b6005805460ff808616700100000000000000000000000000000000027fffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffffff63ffffffff8981166c01000000000000000000000000027fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff918c166801000000000000000002919091167fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff909516949094179390931792909216919091179091556006839055600780549183167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00928316179055600380549091166001179055604080517fc3f909d4000000000000000000000000000000000000000000000000000000008152905173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169163c3f909d4916004828101926080929190829003018186803b158015610f9e57600080fd5b505afa158015610fb2573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610fd69190611b99565b50600580547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff929092169190911790555050604080517f356dac7100000000000000000000000000000000000000000000000000000000815290517f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169163356dac71916004808301926020929190829003018186803b15801561109457600080fd5b505afa1580156110a8573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110cc9190611b80565b6004819055507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16635fbbc0d26040518163ffffffff1660e01b81526004016101206040518083038186803b15801561113957600080fd5b505afa15801561114d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111719190611dd0565b50506005805463ffffffff909816640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff909816979097179096555050505050505050505050565b6111c66114b2565b6111cf816117f5565b50565b6111da6114b2565b6040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8381166004830152602482018390527f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb90604401602060405180830381600087803b15801561126a57600080fd5b505af115801561127e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112a29190611b5e565b505050565b6000828152600860208181526040808420815160a081018352815473ffffffffffffffffffffffffffffffffffffffff808216835263ffffffff740100000000000000000000000000000000000000008304168387015260018401805495840195909552600284018054606085015260038501805460808601528b8a52979096527fffffffffffffffff0000000000000000000000000000000000000000000000009091169092559185905591849055929091558151166113c4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f72657175657374206e6f7420666f756e6400000000000000000000000000000060448201526064016105ac565b600080631fe543e360e01b85856040516024016113e2929190611f9f565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050509050600061145c846020015163ffffffff168560000151846118eb565b9050806114aa57835160405173ffffffffffffffffffffffffffffffffffffffff9091169087907fc551b83c151f2d1c7eeb938ac59008e0409f1c1dc1e2f112449d4d79b458902290600090a35b505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314611533576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016105ac565b565b600554604080517ffeaf968c000000000000000000000000000000000000000000000000000000008152905160009263ffffffff161515918391829173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169163feaf968c9160048082019260a092909190829003018186803b1580156115cf57600080fd5b505afa1580156115e3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906116079190611ee8565b50945090925084915050801561162d575061162282426120c6565b60055463ffffffff16105b1561163757506004545b60008112156116a2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f496e76616c6964204c494e4b207765692070726963650000000000000000000060448201526064016105ac565b9392505050565b6000806116b4611937565b5050505091505060006116c5611a06565b6116d1906101f4611fed565b6116db9083612089565b600554909150600090829063ffffffff6c0100000000000000000000000082048116916117169168010000000000000000909104168a611fed565b6117209190611fed565b61172a9088612089565b6117349190611fed565b905060008561174b83670de0b6b3a7640000612089565b6117559190612052565b60055490915060009060649061178290700100000000000000000000000000000000900460ff168261202d565b61178f9060ff1684612089565b6117999190612052565b6005549091506000906117bf90640100000000900463ffffffff1664e8d4a51000612089565b6117c99083611fed565b9a9950505050505050505050565b60006117e4603f83612066565b6117ef906001612005565b92915050565b73ffffffffffffffffffffffffffffffffffffffff8116331415611875576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016105ac565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60005a6113888110156118fd57600080fd5b61138881039050846040820482031161191557600080fd5b50823b61192157600080fd5b60008083516020850160008789f1949350505050565b600080808080804661a4b1811480611951575062066eed81145b156119e757606c73ffffffffffffffffffffffffffffffffffffffff166341b247a86040518163ffffffff1660e01b815260040160c06040518083038186803b15801561199d57600080fd5b505afa1580156119b1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906119d59190611d00565b965096509650965096509650506119fe565b600080600080600080965096509650965096509650505b909192939495565b60004661a4b1811480611a1b575062066eed81145b15611a2857608c91505090565b600091505090565b803573ffffffffffffffffffffffffffffffffffffffff811681146106fb57600080fd5b805162ffffff811681146106fb57600080fd5b803560ff811681146106fb57600080fd5b805169ffffffffffffffffffff811681146106fb57600080fd5b600060208284031215611aa457600080fd5b6116a282611a30565b60008060408385031215611ac057600080fd5b611ac983611a30565b946020939093013593505050565b60008060008060608587031215611aed57600080fd5b611af685611a30565b935060208501359250604085013567ffffffffffffffff80821115611b1a57600080fd5b818701915087601f830112611b2e57600080fd5b813581811115611b3d57600080fd5b886020828501011115611b4f57600080fd5b95989497505060200194505050565b600060208284031215611b7057600080fd5b815180151581146116a257600080fd5b600060208284031215611b9257600080fd5b5051919050565b60008060008060808587031215611baf57600080fd5b8451611bba8161216a565b6020860151909450611bcb8161217a565b6040860151909350611bdc8161217a565b6060860151909250611bed8161217a565b939692955090935050565b600060208284031215611c0a57600080fd5b5035919050565b60008060408385031215611c2457600080fd5b8235915060208084013567ffffffffffffffff80821115611c4457600080fd5b818601915086601f830112611c5857600080fd5b813581811115611c6a57611c6a61213b565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f83011681018181108582111715611cad57611cad61213b565b604052828152858101935084860182860187018b1015611ccc57600080fd5b600095505b83861015611cef578035855260019590950194938601938601611cd1565b508096505050505050509250929050565b60008060008060008060c08789031215611d1957600080fd5b865195506020870151945060408701519350606087015192506080870151915060a087015190509295509295509295565b600060208284031215611d5c57600080fd5b81356116a28161217a565b600080600060608486031215611d7c57600080fd5b8335611d878161217a565b92506020840135611d978161216a565b91506040840135611da78161217a565b809150509250925092565b60008060408385031215611dc557600080fd5b8235611ac98161217a565b60008060008060008060008060006101208a8c031215611def57600080fd5b8951611dfa8161217a565b60208b0151909950611e0b8161217a565b60408b0151909850611e1c8161217a565b60608b0151909750611e2d8161217a565b60808b0151909650611e3e8161217a565b9450611e4c60a08b01611a54565b9350611e5a60c08b01611a54565b9250611e6860e08b01611a54565b9150611e776101008b01611a54565b90509295985092959850929598565b600080600080600060a08688031215611e9e57600080fd5b8535611ea98161217a565b94506020860135611eb98161217a565b9350611ec760408701611a67565b925060608601359150611edc60808701611a67565b90509295509295909350565b600080600080600060a08688031215611f0057600080fd5b611f0986611a78565b9450602086015193506040860151925060608601519150611edc60808701611a78565b600060208083528351808285015260005b81811015611f5957858101830151858201604001528201611f3d565b81811115611f6b576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b6000604082018483526020604081850152818551808452606086019150828701935060005b81811015611fe057845183529383019391830191600101611fc4565b5090979650505050505050565b60008219821115612000576120006120dd565b500190565b600063ffffffff808316818516808303821115612024576120246120dd565b01949350505050565b600060ff821660ff84168060ff0382111561204a5761204a6120dd565b019392505050565b6000826120615761206161210c565b500490565b600063ffffffff8084168061207d5761207d61210c565b92169190910492915050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156120c1576120c16120dd565b500290565b6000828210156120d8576120d86120dd565b500390565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b61ffff811681146111cf57600080fd5b63ffffffff811681146111cf57600080fdfea164736f6c6343000806000a",
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
	return address, tx, &VRFV2Wrapper{VRFV2WrapperCaller: VRFV2WrapperCaller{contract: contract}, VRFV2WrapperTransactor: VRFV2WrapperTransactor{contract: contract}, VRFV2WrapperFilterer: VRFV2WrapperFilterer{contract: contract}}, nil
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

func (_VRFV2Wrapper *VRFV2WrapperCaller) FULFILLMENTCOMPRESSEDTXDATASIZEBYTES(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _VRFV2Wrapper.contract.Call(opts, &out, "FULFILLMENT_COMPRESSED_TX_DATA_SIZE_BYTES")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_VRFV2Wrapper *VRFV2WrapperSession) FULFILLMENTCOMPRESSEDTXDATASIZEBYTES() (uint32, error) {
	return _VRFV2Wrapper.Contract.FULFILLMENTCOMPRESSEDTXDATASIZEBYTES(&_VRFV2Wrapper.CallOpts)
}

func (_VRFV2Wrapper *VRFV2WrapperCallerSession) FULFILLMENTCOMPRESSEDTXDATASIZEBYTES() (uint32, error) {
	return _VRFV2Wrapper.Contract.FULFILLMENTCOMPRESSEDTXDATASIZEBYTES(&_VRFV2Wrapper.CallOpts)
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

	FULFILLMENTCOMPRESSEDTXDATASIZEBYTES(opts *bind.CallOpts) (uint32, error)

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

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	Disable(opts *bind.TransactOpts) (*types.Transaction, error)

	Enable(opts *bind.TransactOpts) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, _sender common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error)

	RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, _wrapperGasOverhead uint32, _coordinatorGasOverhead uint32, _wrapperPremiumPercentage uint8, _keyHash [32]byte, _maxNumWords uint8) (*types.Transaction, error)

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
