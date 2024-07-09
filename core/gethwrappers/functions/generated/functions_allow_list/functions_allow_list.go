// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package functions_allow_list

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

type TermsOfServiceAllowListConfig struct {
	Enabled         bool
	SignerPublicKey common.Address
}

var TermsOfServiceAllowListMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"signerPublicKey\",\"type\":\"address\"}],\"internalType\":\"structTermsOfServiceAllowListConfig\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"initialAllowedSenders\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"initialBlockedSenders\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"previousToSContract\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidUsage\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RecipientIsBlocked\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"AddedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"BlockedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"signerPublicKey\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structTermsOfServiceAllowListConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"UnblockedAccess\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"acceptor\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"}],\"name\":\"acceptTermsOfService\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"blockSender\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllAllowedSenders\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowedSendersCount\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"allowedSenderIdxStart\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"allowedSenderIdxEnd\",\"type\":\"uint64\"}],\"name\":\"getAllowedSendersInRange\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"allowedSenders\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBlockedSendersCount\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"blockedSenderIdxStart\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"blockedSenderIdxEnd\",\"type\":\"uint64\"}],\"name\":\"getBlockedSendersInRange\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"blockedSenders\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"signerPublicKey\",\"type\":\"address\"}],\"internalType\":\"structTermsOfServiceAllowListConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"acceptor\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"getMessage\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"isBlockedSender\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_previousToSContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"unblockSender\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"signerPublicKey\",\"type\":\"address\"}],\"internalType\":\"structTermsOfServiceAllowListConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"updateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateFromPrevious\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001d3638038062001d36833981016040819052620000349162000525565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620001f5565b505050620000d284620002a060201b60201c565b60005b8351811015620001255762000111848281518110620000f857620000f8620005eb565b602002602001015160036200032760201b90919060201c565b506200011d8162000601565b9050620000d5565b5060005b8251811015620001ca57620001658382815181106200014c576200014c620005eb565b602002602001015160036200034760201b90919060201c565b156200018457604051638129bbcd60e01b815260040160405180910390fd5b620001b68382815181106200019d576200019d620005eb565b602002602001015160056200032760201b90919060201c565b50620001c28162000601565b905062000129565b50600280546001600160a01b0319166001600160a01b03929092169190911790555062000629915050565b336001600160a01b038216036200024f5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b620002aa6200036a565b805160078054602080850180516001600160a81b0319909316941515610100600160a81b03198116959095176101006001600160a01b039485160217909355604080519485529251909116908301527f0d22b8a99f411b3dd338c961284f608489ca0dab9cdad17366a343c361bcf80a910160405180910390a150565b60006200033e836001600160a01b038416620003c8565b90505b92915050565b6001600160a01b038116600090815260018301602052604081205415156200033e565b6000546001600160a01b03163314620003c65760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000082565b565b6000818152600183016020526040812054620004115750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915562000341565b50600062000341565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b03811182821017156200045557620004556200041a565b60405290565b80516001600160a01b03811681146200047357600080fd5b919050565b600082601f8301126200048a57600080fd5b815160206001600160401b0380831115620004a957620004a96200041a565b8260051b604051601f19603f83011681018181108482111715620004d157620004d16200041a565b604052938452858101830193838101925087851115620004f057600080fd5b83870191505b848210156200051a576200050a826200045b565b83529183019190830190620004f6565b979650505050505050565b60008060008084860360a08112156200053d57600080fd5b60408112156200054c57600080fd5b506200055762000430565b855180151581146200056857600080fd5b815262000578602087016200045b565b602082015260408601519094506001600160401b03808211156200059b57600080fd5b620005a98883890162000478565b94506060870151915080821115620005c057600080fd5b50620005cf8782880162000478565b925050620005e0608086016200045b565b905092959194509250565b634e487b7160e01b600052603260045260246000fd5b6000600182016200062257634e487b7160e01b600052601160045260246000fd5b5060010190565b6116fd80620006396000396000f3fe608060405234801561001057600080fd5b50600436106101515760003560e01c806382184c7b116100cd578063a5e1d61d11610081578063cc7ebf4911610066578063cc7ebf491461038f578063d2a6ee3714610397578063f2fde38b1461039f57600080fd5b8063a5e1d61d1461031d578063c3f909d41461033057600080fd5b80638da5cb5b116100b25780638da5cb5b1461024557806397231e8214610284578063a39b06e3146102a457600080fd5b806382184c7b1461021f57806389f9a2c41461023257600080fd5b80633908c4d4116101245780636b14daf8116101095780636b14daf8146101ec57806379ba50971461020f578063817ef62e1461021757600080fd5b80633908c4d4146101c457806347663acb146101d957600080fd5b806301a05958146101565780630a8c9c241461017c578063181f5a771461019c57806320229a86146101b1575b600080fd5b61015e6103b2565b60405167ffffffffffffffff90911681526020015b60405180910390f35b61018f61018a3660046111a8565b6103c3565b60405161017391906111db565b6101a461051f565b6040516101739190611235565b61018f6101bf3660046111a8565b61053b565b6101d76101d23660046112c3565b6106a1565b005b6101d76101e7366004611328565b61094c565b6101ff6101fa366004611345565b6109ad565b6040519015158152602001610173565b6101d76109d7565b61018f610ad9565b6101d761022d366004611328565b610ae5565b6101d7610240366004611448565b610b4b565b60005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610173565b60025461025f9073ffffffffffffffffffffffffffffffffffffffff1681565b61030f6102b23660046114ac565b6040517fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606084811b8216602084015283901b16603482015260009060480160405160208183030381529060405280519060200120905092915050565b604051908152602001610173565b6101ff61032b366004611328565b610c06565b60408051808201825260008082526020918201528151808301835260075460ff8116151580835273ffffffffffffffffffffffffffffffffffffffff610100909204821692840192835284519081529151169181019190915201610173565b61015e610c26565b6101d7610c32565b6101d76103ad366004611328565b610da1565b60006103be6005610db5565b905090565b60608167ffffffffffffffff168367ffffffffffffffff1611806103fa57506103ec6003610db5565b8267ffffffffffffffff1610155b15610431576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61043b8383611514565b610446906001611535565b67ffffffffffffffff1667ffffffffffffffff811115610468576104686113ca565b604051908082528060200260200182016040528015610491578160200160208202803683370190505b50905060005b6104a18484611514565b67ffffffffffffffff168111610517576104d06104c88267ffffffffffffffff8716611556565b600390610dbf565b8282815181106104e2576104e2611569565b73ffffffffffffffffffffffffffffffffffffffff9092166020928302919091019091015261051081611598565b9050610497565b505b92915050565b6040518060600160405280602c81526020016116c5602c913981565b60608167ffffffffffffffff168367ffffffffffffffff16118061057257506105646005610db5565b8267ffffffffffffffff1610155b8061058457506105826005610db5565b155b156105bb576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6105c58383611514565b6105d0906001611535565b67ffffffffffffffff1667ffffffffffffffff8111156105f2576105f26113ca565b60405190808252806020026020018201604052801561061b578160200160208202803683370190505b50905060005b61062b8484611514565b67ffffffffffffffff1681116105175761065a6106528267ffffffffffffffff8716611556565b600590610dbf565b82828151811061066c5761066c611569565b73ffffffffffffffffffffffffffffffffffffffff9092166020928302919091019091015261069a81611598565b9050610621565b6106ac600585610dcb565b156106e3576040517f62b7a34d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60408051606087811b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000009081166020808501919091529188901b16603483015282516028818403018152604890920190925280519101206000906040517f19457468657265756d205369676e6564204d6573736167653a0a3332000000006020820152603c810191909152605c01604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815282825280516020918201206007546000855291840180845281905260ff8616928401929092526060830187905260808301869052909250610100900473ffffffffffffffffffffffffffffffffffffffff169060019060a0016020604051602081039080840390855afa158015610817573d6000803e3d6000fd5b5050506020604051035173ffffffffffffffffffffffffffffffffffffffff161461086e576040517f8baa579f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff86161415806108b357503373ffffffffffffffffffffffffffffffffffffffff8716148015906108b35750333b155b156108ea576040517f381cfcbd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108f5600386610dfa565b156109445760405173ffffffffffffffffffffffffffffffffffffffff861681527f87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db49060200160405180910390a15b505050505050565b610954610e1c565b61095f600582610e9f565b5060405173ffffffffffffffffffffffffffffffffffffffff821681527f28bbd0761309a99e8fb5e5d02ada0b7b2db2e5357531ff5dbfc205c3f5b6592b906020015b60405180910390a150565b60075460009060ff166109c2575060016109d0565b6109cd600385610dcb565b90505b9392505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610a5d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60606103be6003610ec1565b610aed610e1c565b610af8600382610e9f565b50610b04600582610dfa565b5060405173ffffffffffffffffffffffffffffffffffffffff821681527f337cd0f3f594112b6d830afb510072d3b08556b446514f73b8109162fd1151e1906020016109a2565b610b53610e1c565b805160078054602080850180517fffffffffffffffffffffff0000000000000000000000000000000000000000009093169415157fffffffffffffffffffffff0000000000000000000000000000000000000000ff81169590951761010073ffffffffffffffffffffffffffffffffffffffff9485160217909355604080519485529251909116908301527f0d22b8a99f411b3dd338c961284f608489ca0dab9cdad17366a343c361bcf80a91016109a2565b60075460009060ff16610c1b57506000919050565b610519600583610dcb565b60006103be6003610db5565b600254604080517f817ef62e000000000000000000000000000000000000000000000000000000008152905160009273ffffffffffffffffffffffffffffffffffffffff169163817ef62e91600480830192869291908290030181865afa158015610ca1573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052610ce791908101906115d0565b905060005b8151811015610d9d57610d22828281518110610d0a57610d0a611569565b60200260200101516005610dcb90919063ffffffff16565b158015610d595750610d57828281518110610d3f57610d3f611569565b60200260200101516003610dcb90919063ffffffff16565b155b15610d8d57610d8b828281518110610d7357610d73611569565b60200260200101516003610dfa90919063ffffffff16565b505b610d9681611598565b9050610cec565b5050565b610da9610e1c565b610db281610ece565b50565b6000610519825490565b60006109d08383610fc3565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415156109d0565b60006109d08373ffffffffffffffffffffffffffffffffffffffff8416610fed565b60005473ffffffffffffffffffffffffffffffffffffffff163314610e9d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610a54565b565b60006109d08373ffffffffffffffffffffffffffffffffffffffff841661103c565b606060006109d08361112f565b3373ffffffffffffffffffffffffffffffffffffffff821603610f4d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610a54565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000826000018281548110610fda57610fda611569565b9060005260206000200154905092915050565b600081815260018301602052604081205461103457508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610519565b506000610519565b60008181526001830160205260408120548015611125576000611060600183611682565b855490915060009061107490600190611682565b90508181146110d957600086600001828154811061109457611094611569565b90600052602060002001549050808760000184815481106110b7576110b7611569565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806110ea576110ea611695565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610519565b6000915050610519565b60608160000180548060200260200160405190810160405280929190818152602001828054801561117f57602002820191906000526020600020905b81548152602001906001019080831161116b575b50505050509050919050565b803567ffffffffffffffff811681146111a357600080fd5b919050565b600080604083850312156111bb57600080fd5b6111c48361118b565b91506111d26020840161118b565b90509250929050565b6020808252825182820181905260009190848201906040850190845b8181101561122957835173ffffffffffffffffffffffffffffffffffffffff16835292840192918401916001016111f7565b50909695505050505050565b600060208083528351808285015260005b8181101561126257858101830151858201604001528201611246565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b73ffffffffffffffffffffffffffffffffffffffff81168114610db257600080fd5b600080600080600060a086880312156112db57600080fd5b85356112e6816112a1565b945060208601356112f6816112a1565b93506040860135925060608601359150608086013560ff8116811461131a57600080fd5b809150509295509295909350565b60006020828403121561133a57600080fd5b81356109d0816112a1565b60008060006040848603121561135a57600080fd5b8335611365816112a1565b9250602084013567ffffffffffffffff8082111561138257600080fd5b818601915086601f83011261139657600080fd5b8135818111156113a557600080fd5b8760208285010111156113b757600080fd5b6020830194508093505050509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611440576114406113ca565b604052919050565b60006040828403121561145a57600080fd5b6040516040810181811067ffffffffffffffff8211171561147d5761147d6113ca565b6040528235801515811461149057600080fd5b815260208301356114a0816112a1565b60208201529392505050565b600080604083850312156114bf57600080fd5b82356114ca816112a1565b915060208301356114da816112a1565b809150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b67ffffffffffffffff828116828216039080821115610517576105176114e5565b67ffffffffffffffff818116838216019080821115610517576105176114e5565b80820180821115610519576105196114e5565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036115c9576115c96114e5565b5060010190565b600060208083850312156115e357600080fd5b825167ffffffffffffffff808211156115fb57600080fd5b818501915085601f83011261160f57600080fd5b815181811115611621576116216113ca565b8060051b91506116328483016113f9565b818152918301840191848101908884111561164c57600080fd5b938501935b838510156116765784519250611666836112a1565b8282529385019390850190611651565b98975050505050505050565b81810381811115610519576105196114e5565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfe46756e6374696f6e73205465726d73206f66205365727669636520416c6c6f77204c6973742076312e312e30a164736f6c6343000813000a",
}

var TermsOfServiceAllowListABI = TermsOfServiceAllowListMetaData.ABI

var TermsOfServiceAllowListBin = TermsOfServiceAllowListMetaData.Bin

func DeployTermsOfServiceAllowList(auth *bind.TransactOpts, backend bind.ContractBackend, config TermsOfServiceAllowListConfig, initialAllowedSenders []common.Address, initialBlockedSenders []common.Address, previousToSContract common.Address) (common.Address, *types.Transaction, *TermsOfServiceAllowList, error) {
	parsed, err := TermsOfServiceAllowListMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TermsOfServiceAllowListBin), backend, config, initialAllowedSenders, initialBlockedSenders, previousToSContract)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TermsOfServiceAllowList{address: address, abi: *parsed, TermsOfServiceAllowListCaller: TermsOfServiceAllowListCaller{contract: contract}, TermsOfServiceAllowListTransactor: TermsOfServiceAllowListTransactor{contract: contract}, TermsOfServiceAllowListFilterer: TermsOfServiceAllowListFilterer{contract: contract}}, nil
}

type TermsOfServiceAllowList struct {
	address common.Address
	abi     abi.ABI
	TermsOfServiceAllowListCaller
	TermsOfServiceAllowListTransactor
	TermsOfServiceAllowListFilterer
}

type TermsOfServiceAllowListCaller struct {
	contract *bind.BoundContract
}

type TermsOfServiceAllowListTransactor struct {
	contract *bind.BoundContract
}

type TermsOfServiceAllowListFilterer struct {
	contract *bind.BoundContract
}

type TermsOfServiceAllowListSession struct {
	Contract     *TermsOfServiceAllowList
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type TermsOfServiceAllowListCallerSession struct {
	Contract *TermsOfServiceAllowListCaller
	CallOpts bind.CallOpts
}

type TermsOfServiceAllowListTransactorSession struct {
	Contract     *TermsOfServiceAllowListTransactor
	TransactOpts bind.TransactOpts
}

type TermsOfServiceAllowListRaw struct {
	Contract *TermsOfServiceAllowList
}

type TermsOfServiceAllowListCallerRaw struct {
	Contract *TermsOfServiceAllowListCaller
}

type TermsOfServiceAllowListTransactorRaw struct {
	Contract *TermsOfServiceAllowListTransactor
}

func NewTermsOfServiceAllowList(address common.Address, backend bind.ContractBackend) (*TermsOfServiceAllowList, error) {
	abi, err := abi.JSON(strings.NewReader(TermsOfServiceAllowListABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindTermsOfServiceAllowList(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TermsOfServiceAllowList{address: address, abi: abi, TermsOfServiceAllowListCaller: TermsOfServiceAllowListCaller{contract: contract}, TermsOfServiceAllowListTransactor: TermsOfServiceAllowListTransactor{contract: contract}, TermsOfServiceAllowListFilterer: TermsOfServiceAllowListFilterer{contract: contract}}, nil
}

func NewTermsOfServiceAllowListCaller(address common.Address, caller bind.ContractCaller) (*TermsOfServiceAllowListCaller, error) {
	contract, err := bindTermsOfServiceAllowList(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TermsOfServiceAllowListCaller{contract: contract}, nil
}

func NewTermsOfServiceAllowListTransactor(address common.Address, transactor bind.ContractTransactor) (*TermsOfServiceAllowListTransactor, error) {
	contract, err := bindTermsOfServiceAllowList(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TermsOfServiceAllowListTransactor{contract: contract}, nil
}

func NewTermsOfServiceAllowListFilterer(address common.Address, filterer bind.ContractFilterer) (*TermsOfServiceAllowListFilterer, error) {
	contract, err := bindTermsOfServiceAllowList(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TermsOfServiceAllowListFilterer{contract: contract}, nil
}

func bindTermsOfServiceAllowList(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TermsOfServiceAllowListMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TermsOfServiceAllowList.Contract.TermsOfServiceAllowListCaller.contract.Call(opts, result, method, params...)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.TermsOfServiceAllowListTransactor.contract.Transfer(opts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.TermsOfServiceAllowListTransactor.contract.Transact(opts, method, params...)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TermsOfServiceAllowList.Contract.contract.Call(opts, result, method, params...)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.contract.Transfer(opts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.contract.Transact(opts, method, params...)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) GetAllAllowedSenders(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "getAllAllowedSenders")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) GetAllAllowedSenders() ([]common.Address, error) {
	return _TermsOfServiceAllowList.Contract.GetAllAllowedSenders(&_TermsOfServiceAllowList.CallOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) GetAllAllowedSenders() ([]common.Address, error) {
	return _TermsOfServiceAllowList.Contract.GetAllAllowedSenders(&_TermsOfServiceAllowList.CallOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) GetAllowedSendersCount(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "getAllowedSendersCount")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) GetAllowedSendersCount() (uint64, error) {
	return _TermsOfServiceAllowList.Contract.GetAllowedSendersCount(&_TermsOfServiceAllowList.CallOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) GetAllowedSendersCount() (uint64, error) {
	return _TermsOfServiceAllowList.Contract.GetAllowedSendersCount(&_TermsOfServiceAllowList.CallOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) GetAllowedSendersInRange(opts *bind.CallOpts, allowedSenderIdxStart uint64, allowedSenderIdxEnd uint64) ([]common.Address, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "getAllowedSendersInRange", allowedSenderIdxStart, allowedSenderIdxEnd)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) GetAllowedSendersInRange(allowedSenderIdxStart uint64, allowedSenderIdxEnd uint64) ([]common.Address, error) {
	return _TermsOfServiceAllowList.Contract.GetAllowedSendersInRange(&_TermsOfServiceAllowList.CallOpts, allowedSenderIdxStart, allowedSenderIdxEnd)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) GetAllowedSendersInRange(allowedSenderIdxStart uint64, allowedSenderIdxEnd uint64) ([]common.Address, error) {
	return _TermsOfServiceAllowList.Contract.GetAllowedSendersInRange(&_TermsOfServiceAllowList.CallOpts, allowedSenderIdxStart, allowedSenderIdxEnd)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) GetBlockedSendersCount(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "getBlockedSendersCount")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) GetBlockedSendersCount() (uint64, error) {
	return _TermsOfServiceAllowList.Contract.GetBlockedSendersCount(&_TermsOfServiceAllowList.CallOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) GetBlockedSendersCount() (uint64, error) {
	return _TermsOfServiceAllowList.Contract.GetBlockedSendersCount(&_TermsOfServiceAllowList.CallOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) GetBlockedSendersInRange(opts *bind.CallOpts, blockedSenderIdxStart uint64, blockedSenderIdxEnd uint64) ([]common.Address, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "getBlockedSendersInRange", blockedSenderIdxStart, blockedSenderIdxEnd)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) GetBlockedSendersInRange(blockedSenderIdxStart uint64, blockedSenderIdxEnd uint64) ([]common.Address, error) {
	return _TermsOfServiceAllowList.Contract.GetBlockedSendersInRange(&_TermsOfServiceAllowList.CallOpts, blockedSenderIdxStart, blockedSenderIdxEnd)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) GetBlockedSendersInRange(blockedSenderIdxStart uint64, blockedSenderIdxEnd uint64) ([]common.Address, error) {
	return _TermsOfServiceAllowList.Contract.GetBlockedSendersInRange(&_TermsOfServiceAllowList.CallOpts, blockedSenderIdxStart, blockedSenderIdxEnd)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) GetConfig(opts *bind.CallOpts) (TermsOfServiceAllowListConfig, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "getConfig")

	if err != nil {
		return *new(TermsOfServiceAllowListConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(TermsOfServiceAllowListConfig)).(*TermsOfServiceAllowListConfig)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) GetConfig() (TermsOfServiceAllowListConfig, error) {
	return _TermsOfServiceAllowList.Contract.GetConfig(&_TermsOfServiceAllowList.CallOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) GetConfig() (TermsOfServiceAllowListConfig, error) {
	return _TermsOfServiceAllowList.Contract.GetConfig(&_TermsOfServiceAllowList.CallOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) GetMessage(opts *bind.CallOpts, acceptor common.Address, recipient common.Address) ([32]byte, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "getMessage", acceptor, recipient)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) GetMessage(acceptor common.Address, recipient common.Address) ([32]byte, error) {
	return _TermsOfServiceAllowList.Contract.GetMessage(&_TermsOfServiceAllowList.CallOpts, acceptor, recipient)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) GetMessage(acceptor common.Address, recipient common.Address) ([32]byte, error) {
	return _TermsOfServiceAllowList.Contract.GetMessage(&_TermsOfServiceAllowList.CallOpts, acceptor, recipient)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) HasAccess(opts *bind.CallOpts, user common.Address, arg1 []byte) (bool, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "hasAccess", user, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) HasAccess(user common.Address, arg1 []byte) (bool, error) {
	return _TermsOfServiceAllowList.Contract.HasAccess(&_TermsOfServiceAllowList.CallOpts, user, arg1)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) HasAccess(user common.Address, arg1 []byte) (bool, error) {
	return _TermsOfServiceAllowList.Contract.HasAccess(&_TermsOfServiceAllowList.CallOpts, user, arg1)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) IsBlockedSender(opts *bind.CallOpts, sender common.Address) (bool, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "isBlockedSender", sender)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) IsBlockedSender(sender common.Address) (bool, error) {
	return _TermsOfServiceAllowList.Contract.IsBlockedSender(&_TermsOfServiceAllowList.CallOpts, sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) IsBlockedSender(sender common.Address) (bool, error) {
	return _TermsOfServiceAllowList.Contract.IsBlockedSender(&_TermsOfServiceAllowList.CallOpts, sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) Owner() (common.Address, error) {
	return _TermsOfServiceAllowList.Contract.Owner(&_TermsOfServiceAllowList.CallOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) Owner() (common.Address, error) {
	return _TermsOfServiceAllowList.Contract.Owner(&_TermsOfServiceAllowList.CallOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) SPreviousToSContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "s_previousToSContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) SPreviousToSContract() (common.Address, error) {
	return _TermsOfServiceAllowList.Contract.SPreviousToSContract(&_TermsOfServiceAllowList.CallOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) SPreviousToSContract() (common.Address, error) {
	return _TermsOfServiceAllowList.Contract.SPreviousToSContract(&_TermsOfServiceAllowList.CallOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) TypeAndVersion() (string, error) {
	return _TermsOfServiceAllowList.Contract.TypeAndVersion(&_TermsOfServiceAllowList.CallOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) TypeAndVersion() (string, error) {
	return _TermsOfServiceAllowList.Contract.TypeAndVersion(&_TermsOfServiceAllowList.CallOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.contract.Transact(opts, "acceptOwnership")
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) AcceptOwnership() (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.AcceptOwnership(&_TermsOfServiceAllowList.TransactOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.AcceptOwnership(&_TermsOfServiceAllowList.TransactOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactor) AcceptTermsOfService(opts *bind.TransactOpts, acceptor common.Address, recipient common.Address, r [32]byte, s [32]byte, v uint8) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.contract.Transact(opts, "acceptTermsOfService", acceptor, recipient, r, s, v)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) AcceptTermsOfService(acceptor common.Address, recipient common.Address, r [32]byte, s [32]byte, v uint8) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.AcceptTermsOfService(&_TermsOfServiceAllowList.TransactOpts, acceptor, recipient, r, s, v)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactorSession) AcceptTermsOfService(acceptor common.Address, recipient common.Address, r [32]byte, s [32]byte, v uint8) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.AcceptTermsOfService(&_TermsOfServiceAllowList.TransactOpts, acceptor, recipient, r, s, v)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactor) BlockSender(opts *bind.TransactOpts, sender common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.contract.Transact(opts, "blockSender", sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) BlockSender(sender common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.BlockSender(&_TermsOfServiceAllowList.TransactOpts, sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactorSession) BlockSender(sender common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.BlockSender(&_TermsOfServiceAllowList.TransactOpts, sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.contract.Transact(opts, "transferOwnership", to)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.TransferOwnership(&_TermsOfServiceAllowList.TransactOpts, to)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.TransferOwnership(&_TermsOfServiceAllowList.TransactOpts, to)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactor) UnblockSender(opts *bind.TransactOpts, sender common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.contract.Transact(opts, "unblockSender", sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) UnblockSender(sender common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.UnblockSender(&_TermsOfServiceAllowList.TransactOpts, sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactorSession) UnblockSender(sender common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.UnblockSender(&_TermsOfServiceAllowList.TransactOpts, sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactor) UpdateConfig(opts *bind.TransactOpts, config TermsOfServiceAllowListConfig) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.contract.Transact(opts, "updateConfig", config)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) UpdateConfig(config TermsOfServiceAllowListConfig) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.UpdateConfig(&_TermsOfServiceAllowList.TransactOpts, config)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactorSession) UpdateConfig(config TermsOfServiceAllowListConfig) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.UpdateConfig(&_TermsOfServiceAllowList.TransactOpts, config)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactor) UpdateFromPrevious(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.contract.Transact(opts, "updateFromPrevious")
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) UpdateFromPrevious() (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.UpdateFromPrevious(&_TermsOfServiceAllowList.TransactOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactorSession) UpdateFromPrevious() (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.UpdateFromPrevious(&_TermsOfServiceAllowList.TransactOpts)
}

type TermsOfServiceAllowListAddedAccessIterator struct {
	Event *TermsOfServiceAllowListAddedAccess

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TermsOfServiceAllowListAddedAccessIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TermsOfServiceAllowListAddedAccess)
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
		it.Event = new(TermsOfServiceAllowListAddedAccess)
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

func (it *TermsOfServiceAllowListAddedAccessIterator) Error() error {
	return it.fail
}

func (it *TermsOfServiceAllowListAddedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TermsOfServiceAllowListAddedAccess struct {
	User common.Address
	Raw  types.Log
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) FilterAddedAccess(opts *bind.FilterOpts) (*TermsOfServiceAllowListAddedAccessIterator, error) {

	logs, sub, err := _TermsOfServiceAllowList.contract.FilterLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return &TermsOfServiceAllowListAddedAccessIterator{contract: _TermsOfServiceAllowList.contract, event: "AddedAccess", logs: logs, sub: sub}, nil
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) WatchAddedAccess(opts *bind.WatchOpts, sink chan<- *TermsOfServiceAllowListAddedAccess) (event.Subscription, error) {

	logs, sub, err := _TermsOfServiceAllowList.contract.WatchLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TermsOfServiceAllowListAddedAccess)
				if err := _TermsOfServiceAllowList.contract.UnpackLog(event, "AddedAccess", log); err != nil {
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

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) ParseAddedAccess(log types.Log) (*TermsOfServiceAllowListAddedAccess, error) {
	event := new(TermsOfServiceAllowListAddedAccess)
	if err := _TermsOfServiceAllowList.contract.UnpackLog(event, "AddedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TermsOfServiceAllowListBlockedAccessIterator struct {
	Event *TermsOfServiceAllowListBlockedAccess

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TermsOfServiceAllowListBlockedAccessIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TermsOfServiceAllowListBlockedAccess)
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
		it.Event = new(TermsOfServiceAllowListBlockedAccess)
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

func (it *TermsOfServiceAllowListBlockedAccessIterator) Error() error {
	return it.fail
}

func (it *TermsOfServiceAllowListBlockedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TermsOfServiceAllowListBlockedAccess struct {
	User common.Address
	Raw  types.Log
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) FilterBlockedAccess(opts *bind.FilterOpts) (*TermsOfServiceAllowListBlockedAccessIterator, error) {

	logs, sub, err := _TermsOfServiceAllowList.contract.FilterLogs(opts, "BlockedAccess")
	if err != nil {
		return nil, err
	}
	return &TermsOfServiceAllowListBlockedAccessIterator{contract: _TermsOfServiceAllowList.contract, event: "BlockedAccess", logs: logs, sub: sub}, nil
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) WatchBlockedAccess(opts *bind.WatchOpts, sink chan<- *TermsOfServiceAllowListBlockedAccess) (event.Subscription, error) {

	logs, sub, err := _TermsOfServiceAllowList.contract.WatchLogs(opts, "BlockedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TermsOfServiceAllowListBlockedAccess)
				if err := _TermsOfServiceAllowList.contract.UnpackLog(event, "BlockedAccess", log); err != nil {
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

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) ParseBlockedAccess(log types.Log) (*TermsOfServiceAllowListBlockedAccess, error) {
	event := new(TermsOfServiceAllowListBlockedAccess)
	if err := _TermsOfServiceAllowList.contract.UnpackLog(event, "BlockedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TermsOfServiceAllowListConfigUpdatedIterator struct {
	Event *TermsOfServiceAllowListConfigUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TermsOfServiceAllowListConfigUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TermsOfServiceAllowListConfigUpdated)
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
		it.Event = new(TermsOfServiceAllowListConfigUpdated)
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

func (it *TermsOfServiceAllowListConfigUpdatedIterator) Error() error {
	return it.fail
}

func (it *TermsOfServiceAllowListConfigUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TermsOfServiceAllowListConfigUpdated struct {
	Config TermsOfServiceAllowListConfig
	Raw    types.Log
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) FilterConfigUpdated(opts *bind.FilterOpts) (*TermsOfServiceAllowListConfigUpdatedIterator, error) {

	logs, sub, err := _TermsOfServiceAllowList.contract.FilterLogs(opts, "ConfigUpdated")
	if err != nil {
		return nil, err
	}
	return &TermsOfServiceAllowListConfigUpdatedIterator{contract: _TermsOfServiceAllowList.contract, event: "ConfigUpdated", logs: logs, sub: sub}, nil
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) WatchConfigUpdated(opts *bind.WatchOpts, sink chan<- *TermsOfServiceAllowListConfigUpdated) (event.Subscription, error) {

	logs, sub, err := _TermsOfServiceAllowList.contract.WatchLogs(opts, "ConfigUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TermsOfServiceAllowListConfigUpdated)
				if err := _TermsOfServiceAllowList.contract.UnpackLog(event, "ConfigUpdated", log); err != nil {
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

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) ParseConfigUpdated(log types.Log) (*TermsOfServiceAllowListConfigUpdated, error) {
	event := new(TermsOfServiceAllowListConfigUpdated)
	if err := _TermsOfServiceAllowList.contract.UnpackLog(event, "ConfigUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TermsOfServiceAllowListOwnershipTransferRequestedIterator struct {
	Event *TermsOfServiceAllowListOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TermsOfServiceAllowListOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TermsOfServiceAllowListOwnershipTransferRequested)
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
		it.Event = new(TermsOfServiceAllowListOwnershipTransferRequested)
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

func (it *TermsOfServiceAllowListOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *TermsOfServiceAllowListOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TermsOfServiceAllowListOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TermsOfServiceAllowListOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TermsOfServiceAllowList.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TermsOfServiceAllowListOwnershipTransferRequestedIterator{contract: _TermsOfServiceAllowList.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *TermsOfServiceAllowListOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TermsOfServiceAllowList.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TermsOfServiceAllowListOwnershipTransferRequested)
				if err := _TermsOfServiceAllowList.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) ParseOwnershipTransferRequested(log types.Log) (*TermsOfServiceAllowListOwnershipTransferRequested, error) {
	event := new(TermsOfServiceAllowListOwnershipTransferRequested)
	if err := _TermsOfServiceAllowList.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TermsOfServiceAllowListOwnershipTransferredIterator struct {
	Event *TermsOfServiceAllowListOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TermsOfServiceAllowListOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TermsOfServiceAllowListOwnershipTransferred)
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
		it.Event = new(TermsOfServiceAllowListOwnershipTransferred)
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

func (it *TermsOfServiceAllowListOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *TermsOfServiceAllowListOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TermsOfServiceAllowListOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TermsOfServiceAllowListOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TermsOfServiceAllowList.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TermsOfServiceAllowListOwnershipTransferredIterator{contract: _TermsOfServiceAllowList.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TermsOfServiceAllowListOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TermsOfServiceAllowList.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TermsOfServiceAllowListOwnershipTransferred)
				if err := _TermsOfServiceAllowList.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) ParseOwnershipTransferred(log types.Log) (*TermsOfServiceAllowListOwnershipTransferred, error) {
	event := new(TermsOfServiceAllowListOwnershipTransferred)
	if err := _TermsOfServiceAllowList.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TermsOfServiceAllowListUnblockedAccessIterator struct {
	Event *TermsOfServiceAllowListUnblockedAccess

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TermsOfServiceAllowListUnblockedAccessIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TermsOfServiceAllowListUnblockedAccess)
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
		it.Event = new(TermsOfServiceAllowListUnblockedAccess)
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

func (it *TermsOfServiceAllowListUnblockedAccessIterator) Error() error {
	return it.fail
}

func (it *TermsOfServiceAllowListUnblockedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TermsOfServiceAllowListUnblockedAccess struct {
	User common.Address
	Raw  types.Log
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) FilterUnblockedAccess(opts *bind.FilterOpts) (*TermsOfServiceAllowListUnblockedAccessIterator, error) {

	logs, sub, err := _TermsOfServiceAllowList.contract.FilterLogs(opts, "UnblockedAccess")
	if err != nil {
		return nil, err
	}
	return &TermsOfServiceAllowListUnblockedAccessIterator{contract: _TermsOfServiceAllowList.contract, event: "UnblockedAccess", logs: logs, sub: sub}, nil
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) WatchUnblockedAccess(opts *bind.WatchOpts, sink chan<- *TermsOfServiceAllowListUnblockedAccess) (event.Subscription, error) {

	logs, sub, err := _TermsOfServiceAllowList.contract.WatchLogs(opts, "UnblockedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TermsOfServiceAllowListUnblockedAccess)
				if err := _TermsOfServiceAllowList.contract.UnpackLog(event, "UnblockedAccess", log); err != nil {
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

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) ParseUnblockedAccess(log types.Log) (*TermsOfServiceAllowListUnblockedAccess, error) {
	event := new(TermsOfServiceAllowListUnblockedAccess)
	if err := _TermsOfServiceAllowList.contract.UnpackLog(event, "UnblockedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowList) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _TermsOfServiceAllowList.abi.Events["AddedAccess"].ID:
		return _TermsOfServiceAllowList.ParseAddedAccess(log)
	case _TermsOfServiceAllowList.abi.Events["BlockedAccess"].ID:
		return _TermsOfServiceAllowList.ParseBlockedAccess(log)
	case _TermsOfServiceAllowList.abi.Events["ConfigUpdated"].ID:
		return _TermsOfServiceAllowList.ParseConfigUpdated(log)
	case _TermsOfServiceAllowList.abi.Events["OwnershipTransferRequested"].ID:
		return _TermsOfServiceAllowList.ParseOwnershipTransferRequested(log)
	case _TermsOfServiceAllowList.abi.Events["OwnershipTransferred"].ID:
		return _TermsOfServiceAllowList.ParseOwnershipTransferred(log)
	case _TermsOfServiceAllowList.abi.Events["UnblockedAccess"].ID:
		return _TermsOfServiceAllowList.ParseUnblockedAccess(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (TermsOfServiceAllowListAddedAccess) Topic() common.Hash {
	return common.HexToHash("0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4")
}

func (TermsOfServiceAllowListBlockedAccess) Topic() common.Hash {
	return common.HexToHash("0x337cd0f3f594112b6d830afb510072d3b08556b446514f73b8109162fd1151e1")
}

func (TermsOfServiceAllowListConfigUpdated) Topic() common.Hash {
	return common.HexToHash("0x0d22b8a99f411b3dd338c961284f608489ca0dab9cdad17366a343c361bcf80a")
}

func (TermsOfServiceAllowListOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (TermsOfServiceAllowListOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (TermsOfServiceAllowListUnblockedAccess) Topic() common.Hash {
	return common.HexToHash("0x28bbd0761309a99e8fb5e5d02ada0b7b2db2e5357531ff5dbfc205c3f5b6592b")
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowList) Address() common.Address {
	return _TermsOfServiceAllowList.address
}

type TermsOfServiceAllowListInterface interface {
	GetAllAllowedSenders(opts *bind.CallOpts) ([]common.Address, error)

	GetAllowedSendersCount(opts *bind.CallOpts) (uint64, error)

	GetAllowedSendersInRange(opts *bind.CallOpts, allowedSenderIdxStart uint64, allowedSenderIdxEnd uint64) ([]common.Address, error)

	GetBlockedSendersCount(opts *bind.CallOpts) (uint64, error)

	GetBlockedSendersInRange(opts *bind.CallOpts, blockedSenderIdxStart uint64, blockedSenderIdxEnd uint64) ([]common.Address, error)

	GetConfig(opts *bind.CallOpts) (TermsOfServiceAllowListConfig, error)

	GetMessage(opts *bind.CallOpts, acceptor common.Address, recipient common.Address) ([32]byte, error)

	HasAccess(opts *bind.CallOpts, user common.Address, arg1 []byte) (bool, error)

	IsBlockedSender(opts *bind.CallOpts, sender common.Address) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SPreviousToSContract(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptTermsOfService(opts *bind.TransactOpts, acceptor common.Address, recipient common.Address, r [32]byte, s [32]byte, v uint8) (*types.Transaction, error)

	BlockSender(opts *bind.TransactOpts, sender common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UnblockSender(opts *bind.TransactOpts, sender common.Address) (*types.Transaction, error)

	UpdateConfig(opts *bind.TransactOpts, config TermsOfServiceAllowListConfig) (*types.Transaction, error)

	UpdateFromPrevious(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterAddedAccess(opts *bind.FilterOpts) (*TermsOfServiceAllowListAddedAccessIterator, error)

	WatchAddedAccess(opts *bind.WatchOpts, sink chan<- *TermsOfServiceAllowListAddedAccess) (event.Subscription, error)

	ParseAddedAccess(log types.Log) (*TermsOfServiceAllowListAddedAccess, error)

	FilterBlockedAccess(opts *bind.FilterOpts) (*TermsOfServiceAllowListBlockedAccessIterator, error)

	WatchBlockedAccess(opts *bind.WatchOpts, sink chan<- *TermsOfServiceAllowListBlockedAccess) (event.Subscription, error)

	ParseBlockedAccess(log types.Log) (*TermsOfServiceAllowListBlockedAccess, error)

	FilterConfigUpdated(opts *bind.FilterOpts) (*TermsOfServiceAllowListConfigUpdatedIterator, error)

	WatchConfigUpdated(opts *bind.WatchOpts, sink chan<- *TermsOfServiceAllowListConfigUpdated) (event.Subscription, error)

	ParseConfigUpdated(log types.Log) (*TermsOfServiceAllowListConfigUpdated, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TermsOfServiceAllowListOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *TermsOfServiceAllowListOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*TermsOfServiceAllowListOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TermsOfServiceAllowListOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TermsOfServiceAllowListOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*TermsOfServiceAllowListOwnershipTransferred, error)

	FilterUnblockedAccess(opts *bind.FilterOpts) (*TermsOfServiceAllowListUnblockedAccessIterator, error)

	WatchUnblockedAccess(opts *bind.WatchOpts, sink chan<- *TermsOfServiceAllowListUnblockedAccess) (event.Subscription, error)

	ParseUnblockedAccess(log types.Log) (*TermsOfServiceAllowListUnblockedAccess, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
