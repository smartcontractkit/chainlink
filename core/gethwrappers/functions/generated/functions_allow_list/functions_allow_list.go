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
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"signerPublicKey\",\"type\":\"address\"}],\"internalType\":\"structTermsOfServiceAllowListConfig\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"initialAllowedSenders\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"initialBlockedSenders\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"previousToSContract\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidUsage\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RecipientIsBlocked\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"AddedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"BlockedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"signerPublicKey\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structTermsOfServiceAllowListConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"UnblockedAccess\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"acceptor\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"}],\"name\":\"acceptTermsOfService\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"blockSender\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllAllowedSenders\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowedSendersCount\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"allowedSenderIdxStart\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"allowedSenderIdxEnd\",\"type\":\"uint64\"}],\"name\":\"getAllowedSendersInRange\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"allowedSenders\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBlockedSendersCount\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"blockedSenderIdxStart\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"blockedSenderIdxEnd\",\"type\":\"uint64\"}],\"name\":\"getBlockedSendersInRange\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"blockedSenders\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"signerPublicKey\",\"type\":\"address\"}],\"internalType\":\"structTermsOfServiceAllowListConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"acceptor\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"getMessage\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"isBlockedSender\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"previousSendersToAdd\",\"type\":\"address[]\"}],\"name\":\"migratePreviouslyAllowedSenders\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"unblockSender\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"signerPublicKey\",\"type\":\"address\"}],\"internalType\":\"structTermsOfServiceAllowListConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"updateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001d8838038062001d88833981016040819052620000349162000525565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620001f5565b505050620000d284620002a060201b60201c565b60005b8351811015620001255762000111848281518110620000f857620000f8620005eb565b602002602001015160036200032760201b90919060201c565b506200011d8162000601565b9050620000d5565b5060005b8251811015620001ca57620001658382815181106200014c576200014c620005eb565b602002602001015160036200034760201b90919060201c565b156200018457604051638129bbcd60e01b815260040160405180910390fd5b620001b68382815181106200019d576200019d620005eb565b602002602001015160056200032760201b90919060201c565b50620001c28162000601565b905062000129565b50600280546001600160a01b0319166001600160a01b03929092169190911790555062000629915050565b336001600160a01b038216036200024f5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b620002aa6200036a565b805160078054602080850180516001600160a81b0319909316941515610100600160a81b03198116959095176101006001600160a01b039485160217909355604080519485529251909116908301527f0d22b8a99f411b3dd338c961284f608489ca0dab9cdad17366a343c361bcf80a910160405180910390a150565b60006200033e836001600160a01b038416620003c8565b90505b92915050565b6001600160a01b038116600090815260018301602052604081205415156200033e565b6000546001600160a01b03163314620003c65760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000082565b565b6000818152600183016020526040812054620004115750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915562000341565b50600062000341565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b03811182821017156200045557620004556200041a565b60405290565b80516001600160a01b03811681146200047357600080fd5b919050565b600082601f8301126200048a57600080fd5b815160206001600160401b0380831115620004a957620004a96200041a565b8260051b604051601f19603f83011681018181108482111715620004d157620004d16200041a565b604052938452858101830193838101925087851115620004f057600080fd5b83870191505b848210156200051a576200050a826200045b565b83529183019190830190620004f6565b979650505050505050565b60008060008084860360a08112156200053d57600080fd5b60408112156200054c57600080fd5b506200055762000430565b855180151581146200056857600080fd5b815262000578602087016200045b565b602082015260408601519094506001600160401b03808211156200059b57600080fd5b620005a98883890162000478565b94506060870151915080821115620005c057600080fd5b50620005cf8782880162000478565b925050620005e0608086016200045b565b905092959194509250565b634e487b7160e01b600052603260045260246000fd5b6000600182016200062257634e487b7160e01b600052601160045260246000fd5b5060010190565b61174f80620006396000396000f3fe608060405234801561001057600080fd5b50600436106101365760003560e01c8063817ef62e116100b2578063a39b06e311610081578063c3f909d411610066578063c3f909d4146102f1578063cc7ebf4914610350578063f2fde38b1461035857600080fd5b8063a39b06e314610265578063a5e1d61d146102de57600080fd5b8063817ef62e1461020f57806382184c7b1461021757806389f9a2c41461022a5780638da5cb5b1461023d57600080fd5b80633908c4d4116101095780634e10a5b3116100ee5780634e10a5b3146101d15780636b14daf8146101e457806379ba50971461020757600080fd5b80633908c4d4146101a957806347663acb146101be57600080fd5b806301a059581461013b5780630a8c9c2414610161578063181f5a771461018157806320229a8614610196575b600080fd5b61014361036b565b60405167ffffffffffffffff90911681526020015b60405180910390f35b61017461016f3660046111f0565b61037c565b6040516101589190611223565b6101896104d8565b604051610158919061127d565b6101746101a43660046111f0565b6104f4565b6101bc6101b736600461130d565b61065a565b005b6101bc6101cc36600461136e565b610905565b6101bc6101df366004611407565b610966565b6101f76101f23660046114b4565b610b69565b6040519015158152602001610158565b6101bc610b93565b610174610c90565b6101bc61022536600461136e565b610c9c565b6101bc610238366004611545565b610d02565b60005460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610158565b6102d06102733660046115a2565b6040517fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606084811b8216602084015283901b16603482015260009060480160405160208183030381529060405280519060200120905092915050565b604051908152602001610158565b6101f76102ec36600461136e565b610dbd565b60408051808201825260008082526020918201528151808301835260075460ff8116151580835273ffffffffffffffffffffffffffffffffffffffff610100909204821692840192835284519081529151169181019190915201610158565b610143610ddd565b6101bc61036636600461136e565b610de9565b60006103776005610dfd565b905090565b60608167ffffffffffffffff168367ffffffffffffffff1611806103b357506103a56003610dfd565b8267ffffffffffffffff1610155b156103ea576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6103f483836115fb565b6103ff90600161161c565b67ffffffffffffffff1667ffffffffffffffff81111561042157610421611389565b60405190808252806020026020018201604052801561044a578160200160208202803683370190505b50905060005b61045a84846115fb565b67ffffffffffffffff1681116104d0576104896104818267ffffffffffffffff871661163d565b600390610e07565b82828151811061049b5761049b611650565b73ffffffffffffffffffffffffffffffffffffffff909216602092830291909101909101526104c98161167f565b9050610450565b505b92915050565b6040518060600160405280602c8152602001611717602c913981565b60608167ffffffffffffffff168367ffffffffffffffff16118061052b575061051d6005610dfd565b8267ffffffffffffffff1610155b8061053d575061053b6005610dfd565b155b15610574576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61057e83836115fb565b61058990600161161c565b67ffffffffffffffff1667ffffffffffffffff8111156105ab576105ab611389565b6040519080825280602002602001820160405280156105d4578160200160208202803683370190505b50905060005b6105e484846115fb565b67ffffffffffffffff1681116104d05761061361060b8267ffffffffffffffff871661163d565b600590610e07565b82828151811061062557610625611650565b73ffffffffffffffffffffffffffffffffffffffff909216602092830291909101909101526106538161167f565b90506105da565b610665600585610e13565b1561069c576040517f62b7a34d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60408051606087811b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000009081166020808501919091529188901b16603483015282516028818403018152604890920190925280519101206000906040517f19457468657265756d205369676e6564204d6573736167653a0a3332000000006020820152603c810191909152605c01604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815282825280516020918201206007546000855291840180845281905260ff8616928401929092526060830187905260808301869052909250610100900473ffffffffffffffffffffffffffffffffffffffff169060019060a0016020604051602081039080840390855afa1580156107d0573d6000803e3d6000fd5b5050506020604051035173ffffffffffffffffffffffffffffffffffffffff1614610827576040517f8baa579f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff861614158061086c57503373ffffffffffffffffffffffffffffffffffffffff87161480159061086c5750333b155b156108a3576040517f381cfcbd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108ae600386610e42565b156108fd5760405173ffffffffffffffffffffffffffffffffffffffff861681527f87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db49060200160405180910390a15b505050505050565b61090d610e64565b610918600582610ee7565b5060405173ffffffffffffffffffffffffffffffffffffffff821681527f28bbd0761309a99e8fb5e5d02ada0b7b2db2e5357531ff5dbfc205c3f5b6592b906020015b60405180910390a150565b61096e610e64565b60025473ffffffffffffffffffffffffffffffffffffffff166109f2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f4e6f2070726576696f757320546f5320636f6e7472616374000000000000000060448201526064015b60405180910390fd5b60025473ffffffffffffffffffffffffffffffffffffffff1660005b8251811015610b64578173ffffffffffffffffffffffffffffffffffffffff16636b14daf8848381518110610a4557610a45611650565b6020908102919091010151604080517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff9092166004830152602482015260006044820152606401602060405180830381865afa158015610ac6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610aea91906116b7565b8015610b205750610b1e838281518110610b0657610b06611650565b60200260200101516005610e1390919063ffffffff16565b155b15610b5457610b52838281518110610b3a57610b3a611650565b60200260200101516003610e4290919063ffffffff16565b505b610b5d8161167f565b9050610a0e565b505050565b60075460009060ff16610b7e57506001610b8c565b610b89600385610e13565b90505b9392505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610c14576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016109e9565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60606103776003610f09565b610ca4610e64565b610caf600382610ee7565b50610cbb600582610e42565b5060405173ffffffffffffffffffffffffffffffffffffffff821681527f337cd0f3f594112b6d830afb510072d3b08556b446514f73b8109162fd1151e19060200161095b565b610d0a610e64565b805160078054602080850180517fffffffffffffffffffffff0000000000000000000000000000000000000000009093169415157fffffffffffffffffffffff0000000000000000000000000000000000000000ff81169590951761010073ffffffffffffffffffffffffffffffffffffffff9485160217909355604080519485529251909116908301527f0d22b8a99f411b3dd338c961284f608489ca0dab9cdad17366a343c361bcf80a910161095b565b60075460009060ff16610dd257506000919050565b6104d2600583610e13565b60006103776003610dfd565b610df1610e64565b610dfa81610f16565b50565b60006104d2825490565b6000610b8c838361100b565b73ffffffffffffffffffffffffffffffffffffffff811660009081526001830160205260408120541515610b8c565b6000610b8c8373ffffffffffffffffffffffffffffffffffffffff8416611035565b60005473ffffffffffffffffffffffffffffffffffffffff163314610ee5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016109e9565b565b6000610b8c8373ffffffffffffffffffffffffffffffffffffffff8416611084565b60606000610b8c83611177565b3373ffffffffffffffffffffffffffffffffffffffff821603610f95576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016109e9565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600082600001828154811061102257611022611650565b9060005260206000200154905092915050565b600081815260018301602052604081205461107c575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556104d2565b5060006104d2565b6000818152600183016020526040812054801561116d5760006110a86001836116d4565b85549091506000906110bc906001906116d4565b90508181146111215760008660000182815481106110dc576110dc611650565b90600052602060002001549050808760000184815481106110ff576110ff611650565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080611132576111326116e7565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506104d2565b60009150506104d2565b6060816000018054806020026020016040519081016040528092919081815260200182805480156111c757602002820191906000526020600020905b8154815260200190600101908083116111b3575b50505050509050919050565b803567ffffffffffffffff811681146111eb57600080fd5b919050565b6000806040838503121561120357600080fd5b61120c836111d3565b915061121a602084016111d3565b90509250929050565b6020808252825182820181905260009190848201906040850190845b8181101561127157835173ffffffffffffffffffffffffffffffffffffffff168352928401929184019160010161123f565b50909695505050505050565b600060208083528351808285015260005b818110156112aa5785810183015185820160400152820161128e565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b803573ffffffffffffffffffffffffffffffffffffffff811681146111eb57600080fd5b600080600080600060a0868803121561132557600080fd5b61132e866112e9565b945061133c602087016112e9565b93506040860135925060608601359150608086013560ff8116811461136057600080fd5b809150509295509295909350565b60006020828403121561138057600080fd5b610b8c826112e9565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156113ff576113ff611389565b604052919050565b6000602080838503121561141a57600080fd5b823567ffffffffffffffff8082111561143257600080fd5b818501915085601f83011261144657600080fd5b81358181111561145857611458611389565b8060051b91506114698483016113b8565b818152918301840191848101908884111561148357600080fd5b938501935b838510156114a857611499856112e9565b82529385019390850190611488565b98975050505050505050565b6000806000604084860312156114c957600080fd5b6114d2846112e9565b9250602084013567ffffffffffffffff808211156114ef57600080fd5b818601915086601f83011261150357600080fd5b81358181111561151257600080fd5b87602082850101111561152457600080fd5b6020830194508093505050509250925092565b8015158114610dfa57600080fd5b60006040828403121561155757600080fd5b6040516040810181811067ffffffffffffffff8211171561157a5761157a611389565b604052823561158881611537565b8152611596602084016112e9565b60208201529392505050565b600080604083850312156115b557600080fd5b6115be836112e9565b915061121a602084016112e9565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b67ffffffffffffffff8281168282160390808211156104d0576104d06115cc565b67ffffffffffffffff8181168382160190808211156104d0576104d06115cc565b808201808211156104d2576104d26115cc565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036116b0576116b06115cc565b5060010190565b6000602082840312156116c957600080fd5b8151610b8c81611537565b818103818111156104d2576104d26115cc565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfe46756e6374696f6e73205465726d73206f66205365727669636520416c6c6f77204c6973742076312e312e31a164736f6c6343000813000a",
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

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactor) MigratePreviouslyAllowedSenders(opts *bind.TransactOpts, previousSendersToAdd []common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.contract.Transact(opts, "migratePreviouslyAllowedSenders", previousSendersToAdd)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) MigratePreviouslyAllowedSenders(previousSendersToAdd []common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.MigratePreviouslyAllowedSenders(&_TermsOfServiceAllowList.TransactOpts, previousSendersToAdd)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactorSession) MigratePreviouslyAllowedSenders(previousSendersToAdd []common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.MigratePreviouslyAllowedSenders(&_TermsOfServiceAllowList.TransactOpts, previousSendersToAdd)
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

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptTermsOfService(opts *bind.TransactOpts, acceptor common.Address, recipient common.Address, r [32]byte, s [32]byte, v uint8) (*types.Transaction, error)

	BlockSender(opts *bind.TransactOpts, sender common.Address) (*types.Transaction, error)

	MigratePreviouslyAllowedSenders(opts *bind.TransactOpts, previousSendersToAdd []common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UnblockSender(opts *bind.TransactOpts, sender common.Address) (*types.Transaction, error)

	UpdateConfig(opts *bind.TransactOpts, config TermsOfServiceAllowListConfig) (*types.Transaction, error)

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
