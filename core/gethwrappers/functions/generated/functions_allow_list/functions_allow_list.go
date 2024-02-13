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
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"signerPublicKey\",\"type\":\"address\"}],\"internalType\":\"structTermsOfServiceAllowListConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidUsage\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RecipientIsBlocked\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"AddedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"BlockedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"signerPublicKey\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structTermsOfServiceAllowListConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"UnblockedAccess\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"acceptor\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"}],\"name\":\"acceptTermsOfService\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"blockSender\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllAllowedSenders\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowedSendersCount\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"allowedSenderIdxStart\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"allowedSenderIdxEnd\",\"type\":\"uint64\"}],\"name\":\"getAllowedSendersInRange\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"allowedSenders\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBlockedSendersCount\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"blockedSenderIdxStart\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"blockedSenderIdxEnd\",\"type\":\"uint64\"}],\"name\":\"getBlockedSendersInRange\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"blockedSenders\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"signerPublicKey\",\"type\":\"address\"}],\"internalType\":\"structTermsOfServiceAllowListConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"acceptor\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"getMessage\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"isBlockedSender\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"unblockSender\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"signerPublicKey\",\"type\":\"address\"}],\"internalType\":\"structTermsOfServiceAllowListConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"updateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b506040516200173438038062001734833981016040819052620000349162000269565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000d9565b505050620000d2816200018460201b60201c565b50620002ea565b336001600160a01b03821603620001335760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6200018e6200020b565b805160068054602080850180516001600160a81b0319909316941515610100600160a81b03198116959095176101006001600160a01b039485160217909355604080519485529251909116908301527f0d22b8a99f411b3dd338c961284f608489ca0dab9cdad17366a343c361bcf80a910160405180910390a150565b6000546001600160a01b03163314620002675760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000082565b565b6000604082840312156200027c57600080fd5b604080519081016001600160401b0381118282101715620002ad57634e487b7160e01b600052604160045260246000fd5b60405282518015158114620002c157600080fd5b815260208301516001600160a01b0381168114620002de57600080fd5b60208201529392505050565b61143a80620002fa6000396000f3fe608060405234801561001057600080fd5b506004361061011b5760003560e01c8063817ef62e116100b2578063a39b06e311610081578063c3f909d411610066578063c3f909d4146102c3578063cc7ebf4914610322578063f2fde38b1461032a57600080fd5b8063a39b06e314610237578063a5e1d61d146102b057600080fd5b8063817ef62e146101e157806382184c7b146101e957806389f9a2c4146101fc5780638da5cb5b1461020f57600080fd5b80633908c4d4116100ee5780633908c4d41461018e57806347663acb146101a35780636b14daf8146101b657806379ba5097146101d957600080fd5b806301a05958146101205780630a8c9c2414610146578063181f5a771461016657806320229a861461017b575b600080fd5b61012861033d565b60405167ffffffffffffffff90911681526020015b60405180910390f35b610159610154366004610fd6565b61034e565b60405161013d9190611009565b61016e6104bc565b60405161013d9190611063565b610159610189366004610fd6565b6104d8565b6101a161019c3660046110f3565b61063e565b005b6101a16101b1366004611154565b6108e9565b6101c96101c436600461116f565b61094a565b604051901515815260200161013d565b6101a1610974565b610159610a76565b6101a16101f7366004611154565b610a82565b6101a161020a366004611221565b610ae8565b60005460405173ffffffffffffffffffffffffffffffffffffffff909116815260200161013d565b6102a26102453660046112aa565b6040517fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606084811b8216602084015283901b16603482015260009060480160405160208183030381529060405280519060200120905092915050565b60405190815260200161013d565b6101c96102be366004611154565b610ba3565b60408051808201825260008082526020918201528151808301835260065460ff8116151580835273ffffffffffffffffffffffffffffffffffffffff61010090920482169284019283528451908152915116918101919091520161013d565b610128610bc3565b6101a1610338366004611154565b610bcf565b60006103496004610be3565b905090565b60608167ffffffffffffffff168367ffffffffffffffff16118061038557506103776002610be3565b8267ffffffffffffffff1610155b8061039757506103956002610be3565b155b156103ce576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6103d88383611303565b6103e3906001611324565b67ffffffffffffffff1667ffffffffffffffff811115610405576104056111f2565b60405190808252806020026020018201604052801561042e578160200160208202803683370190505b50905060005b61043e8484611303565b67ffffffffffffffff1681116104b45761046d6104658267ffffffffffffffff8716611345565b600290610bed565b82828151811061047f5761047f611358565b73ffffffffffffffffffffffffffffffffffffffff909216602092830291909101909101526104ad81611387565b9050610434565b505b92915050565b6040518060600160405280602c8152602001611402602c913981565b60608167ffffffffffffffff168367ffffffffffffffff16118061050f57506105016004610be3565b8267ffffffffffffffff1610155b80610521575061051f6004610be3565b155b15610558576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6105628383611303565b61056d906001611324565b67ffffffffffffffff1667ffffffffffffffff81111561058f5761058f6111f2565b6040519080825280602002602001820160405280156105b8578160200160208202803683370190505b50905060005b6105c88484611303565b67ffffffffffffffff1681116104b4576105f76105ef8267ffffffffffffffff8716611345565b600490610bed565b82828151811061060957610609611358565b73ffffffffffffffffffffffffffffffffffffffff9092166020928302919091019091015261063781611387565b90506105be565b610649600485610bf9565b15610680576040517f62b7a34d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60408051606087811b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000009081166020808501919091529188901b16603483015282516028818403018152604890920190925280519101206000906040517f19457468657265756d205369676e6564204d6573736167653a0a3332000000006020820152603c810191909152605c01604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815282825280516020918201206006546000855291840180845281905260ff8616928401929092526060830187905260808301869052909250610100900473ffffffffffffffffffffffffffffffffffffffff169060019060a0016020604051602081039080840390855afa1580156107b4573d6000803e3d6000fd5b5050506020604051035173ffffffffffffffffffffffffffffffffffffffff161461080b576040517f8baa579f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff861614158061085057503373ffffffffffffffffffffffffffffffffffffffff8716148015906108505750333b155b15610887576040517f381cfcbd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610892600286610c28565b156108e15760405173ffffffffffffffffffffffffffffffffffffffff861681527f87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db49060200160405180910390a15b505050505050565b6108f1610c4a565b6108fc600482610ccd565b5060405173ffffffffffffffffffffffffffffffffffffffff821681527f28bbd0761309a99e8fb5e5d02ada0b7b2db2e5357531ff5dbfc205c3f5b6592b906020015b60405180910390a150565b60065460009060ff1661095f5750600161096d565b61096a600285610bf9565b90505b9392505050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146109fa576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b60606103496002610cef565b610a8a610c4a565b610a95600282610ccd565b50610aa1600482610c28565b5060405173ffffffffffffffffffffffffffffffffffffffff821681527f337cd0f3f594112b6d830afb510072d3b08556b446514f73b8109162fd1151e19060200161093f565b610af0610c4a565b805160068054602080850180517fffffffffffffffffffffff0000000000000000000000000000000000000000009093169415157fffffffffffffffffffffff0000000000000000000000000000000000000000ff81169590951761010073ffffffffffffffffffffffffffffffffffffffff9485160217909355604080519485529251909116908301527f0d22b8a99f411b3dd338c961284f608489ca0dab9cdad17366a343c361bcf80a910161093f565b60065460009060ff16610bb857506000919050565b6104b6600483610bf9565b60006103496002610be3565b610bd7610c4a565b610be081610cfc565b50565b60006104b6825490565b600061096d8383610df1565b73ffffffffffffffffffffffffffffffffffffffff81166000908152600183016020526040812054151561096d565b600061096d8373ffffffffffffffffffffffffffffffffffffffff8416610e1b565b60005473ffffffffffffffffffffffffffffffffffffffff163314610ccb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016109f1565b565b600061096d8373ffffffffffffffffffffffffffffffffffffffff8416610e6a565b6060600061096d83610f5d565b3373ffffffffffffffffffffffffffffffffffffffff821603610d7b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016109f1565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000826000018281548110610e0857610e08611358565b9060005260206000200154905092915050565b6000818152600183016020526040812054610e62575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556104b6565b5060006104b6565b60008181526001830160205260408120548015610f53576000610e8e6001836113bf565b8554909150600090610ea2906001906113bf565b9050818114610f07576000866000018281548110610ec257610ec2611358565b9060005260206000200154905080876000018481548110610ee557610ee5611358565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080610f1857610f186113d2565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506104b6565b60009150506104b6565b606081600001805480602002602001604051908101604052809291908181526020018280548015610fad57602002820191906000526020600020905b815481526020019060010190808311610f99575b50505050509050919050565b803567ffffffffffffffff81168114610fd157600080fd5b919050565b60008060408385031215610fe957600080fd5b610ff283610fb9565b915061100060208401610fb9565b90509250929050565b6020808252825182820181905260009190848201906040850190845b8181101561105757835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101611025565b50909695505050505050565b600060208083528351808285015260005b8181101561109057858101830151858201604001528201611074565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b803573ffffffffffffffffffffffffffffffffffffffff81168114610fd157600080fd5b600080600080600060a0868803121561110b57600080fd5b611114866110cf565b9450611122602087016110cf565b93506040860135925060608601359150608086013560ff8116811461114657600080fd5b809150509295509295909350565b60006020828403121561116657600080fd5b61096d826110cf565b60008060006040848603121561118457600080fd5b61118d846110cf565b9250602084013567ffffffffffffffff808211156111aa57600080fd5b818601915086601f8301126111be57600080fd5b8135818111156111cd57600080fd5b8760208285010111156111df57600080fd5b6020830194508093505050509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60006040828403121561123357600080fd5b6040516040810181811067ffffffffffffffff8211171561127d577f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040528235801515811461129057600080fd5b815261129e602084016110cf565b60208201529392505050565b600080604083850312156112bd57600080fd5b6112c6836110cf565b9150611000602084016110cf565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b67ffffffffffffffff8281168282160390808211156104b4576104b46112d4565b67ffffffffffffffff8181168382160190808211156104b4576104b46112d4565b808201808211156104b6576104b66112d4565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036113b8576113b86112d4565b5060010190565b818103818111156104b6576104b66112d4565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfe46756e6374696f6e73205465726d73206f66205365727669636520416c6c6f77204c6973742076312e312e30a164736f6c6343000813000a",
}

var TermsOfServiceAllowListABI = TermsOfServiceAllowListMetaData.ABI

var TermsOfServiceAllowListBin = TermsOfServiceAllowListMetaData.Bin

func DeployTermsOfServiceAllowList(auth *bind.TransactOpts, backend bind.ContractBackend, config TermsOfServiceAllowListConfig) (common.Address, *types.Transaction, *TermsOfServiceAllowList, error) {
	parsed, err := TermsOfServiceAllowListMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TermsOfServiceAllowListBin), backend, config)
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
