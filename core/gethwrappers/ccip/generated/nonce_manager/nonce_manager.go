// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package nonce_manager

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

type AuthorizedCallersAuthorizedCallerArgs struct {
	AddedCallers   []common.Address
	RemovedCallers []common.Address
}

type NonceManagerPreviousRamps struct {
	PrevOnRamp  common.Address
	PrevOffRamp common.Address
}

type NonceManagerPreviousRampsArgs struct {
	RemoteChainSelector   uint64
	OverrideExistingRamps bool
	PrevRamps             NonceManagerPreviousRamps
}

var NonceManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"authorizedCallers\",\"type\":\"address[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PreviousRampAlreadySet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"UnauthorizedCaller\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"AuthorizedCallerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"AuthorizedCallerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"prevOnRamp\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"prevOffRamp\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structNonceManager.PreviousRamps\",\"name\":\"prevRamp\",\"type\":\"tuple\"}],\"name\":\"PreviousRampsUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"}],\"name\":\"SkippedIncorrectNonce\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address[]\",\"name\":\"addedCallers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"removedCallers\",\"type\":\"address[]\"}],\"internalType\":\"structAuthorizedCallers.AuthorizedCallerArgs\",\"name\":\"authorizedCallerArgs\",\"type\":\"tuple\"}],\"name\":\"applyAuthorizedCallerUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"overrideExistingRamps\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"prevOnRamp\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"prevOffRamp\",\"type\":\"address\"}],\"internalType\":\"structNonceManager.PreviousRamps\",\"name\":\"prevRamps\",\"type\":\"tuple\"}],\"internalType\":\"structNonceManager.PreviousRampsArgs[]\",\"name\":\"previousRampsArgs\",\"type\":\"tuple[]\"}],\"name\":\"applyPreviousRampsUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllAuthorizedCallers\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"}],\"name\":\"getInboundNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"getIncrementedOutboundNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"getOutboundNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"getPreviousRamps\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"prevOnRamp\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"prevOffRamp\",\"type\":\"address\"}],\"internalType\":\"structNonceManager.PreviousRamps\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"expectedNonce\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"}],\"name\":\"incrementInboundNonce\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001ad438038062001ad4833981016040819052620000349162000449565b80336000816200005757604051639b15e16f60e01b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b03848116919091179091558116156200008a576200008a81620000c0565b5050604080518082018252828152815160008152602080820190935291810191909152620000b8906200013a565b505062000569565b336001600160a01b03821603620000ea57604051636d6c4ee560e11b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b03838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b602081015160005b8151811015620001ca5760008282815181106200016357620001636200051b565b602090810291909101015190506200017d60028262000289565b15620001c0576040516001600160a01b03821681527fc3803387881faad271c47728894e3e36fac830ffc8602ca6fc07733cbda775809060200160405180910390a15b5060010162000142565b50815160005b815181101562000283576000828281518110620001f157620001f16200051b565b6020026020010151905060006001600160a01b0316816001600160a01b0316036200022f576040516342bcdf7f60e11b815260040160405180910390fd5b6200023c600282620002a9565b506040516001600160a01b03821681527feb1b9b92e50b7f88f9ff25d56765095ac6e91540eee214906f4036a908ffbdef9060200160405180910390a150600101620001d0565b50505050565b6000620002a0836001600160a01b038416620002c0565b90505b92915050565b6000620002a0836001600160a01b038416620003c4565b60008181526001830160205260408120548015620003b9576000620002e760018362000531565b8554909150600090620002fd9060019062000531565b9050818114620003695760008660000182815481106200032157620003216200051b565b90600052602060002001549050808760000184815481106200034757620003476200051b565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806200037d576200037d62000553565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050620002a3565b6000915050620002a3565b60008181526001830160205260408120546200040d57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155620002a3565b506000620002a3565b634e487b7160e01b600052604160045260246000fd5b80516001600160a01b03811681146200044457600080fd5b919050565b600060208083850312156200045d57600080fd5b82516001600160401b03808211156200047557600080fd5b818501915085601f8301126200048a57600080fd5b8151818111156200049f576200049f62000416565b8060051b604051601f19603f83011681018181108582111715620004c757620004c762000416565b604052918252848201925083810185019188831115620004e657600080fd5b938501935b828510156200050f57620004ff856200042c565b84529385019392850192620004eb565b98975050505050505050565b634e487b7160e01b600052603260045260246000fd5b81810381811115620002a357634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052603160045260246000fd5b61155b80620005796000396000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c806391a2749a11610081578063e0e03cae1161005b578063e0e03cae1461027c578063ea458c0c1461029f578063f2fde38b146102b257600080fd5b806391a2749a1461022a578063bf18402a1461023d578063c92236251461026957600080fd5b806379ba5097116100b257806379ba5097146101e55780637a75a094146101ef5780638da5cb5b1461020257600080fd5b8063181f5a77146100d95780632451a6271461012b578063294b563014610140575b600080fd5b6101156040518060400160405280601681526020017f4e6f6e63654d616e6167657220312e362e302d6465760000000000000000000081525081565b6040516101229190610f05565b60405180910390f35b6101336102c5565b6040516101229190610f72565b6101b161014e366004610fe2565b60408051808201909152600080825260208201525067ffffffffffffffff166000908152600460209081526040918290208251808401909352805473ffffffffffffffffffffffffffffffffffffffff9081168452600190910154169082015290565b60408051825173ffffffffffffffffffffffffffffffffffffffff9081168252602093840151169281019290925201610122565b6101ed6102d6565b005b6101ed6101fd366004610fff565b6103a4565b60015460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610122565b6101ed61023836600461118a565b610594565b61025061024b366004611231565b6105a8565b60405167ffffffffffffffff9091168152602001610122565b6102506102773660046112b3565b6105bd565b61028f61028a366004611308565b6105d4565b6040519015158152602001610122565b6102506102ad366004611231565b6106dd565b6101ed6102c036600461136d565b610771565b60606102d16002610782565b905090565b60005473ffffffffffffffffffffffffffffffffffffffff163314610327576040517f02b543c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560008054909116815560405173ffffffffffffffffffffffffffffffffffffffff909216929183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6103ac61078f565b60005b8181101561058f57368383838181106103ca576103ca61138a565b608002919091019150600090506004816103e76020850185610fe2565b67ffffffffffffffff1681526020810191909152604001600020805490915073ffffffffffffffffffffffffffffffffffffffff161515806104425750600181015473ffffffffffffffffffffffffffffffffffffffff1615155b1561048d5761045760408301602084016113b9565b61048d576040517fc6117ae200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61049d606083016040840161136d565b81547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff919091161781556104ed608083016060840161136d565b6001820180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff929092169190911790556105416020830183610fe2565b67ffffffffffffffff167fa2e43edcbc4fd175ae4bebbe3fd6139871ed1f1783cd4a1ace59b90d302c33198360400160405161057d91906113db565b60405180910390a250506001016103af565b505050565b61059c61078f565b6105a5816107e2565b50565b60006105b48383610974565b90505b92915050565b60006105ca848484610a91565b90505b9392505050565b60006105de610be2565b60006105eb868585610a91565b6105f6906001611452565b90508467ffffffffffffffff168167ffffffffffffffff161461065a577f606ff8179e5e3c059b82df931acc496b7b6053e8879042f8267f930e0595f69f868686866040516106489493929190611473565b60405180910390a160009150506106d5565b67ffffffffffffffff861660009081526006602052604090819020905182919061068790879087906114df565b908152604051908190036020019020805467ffffffffffffffff929092167fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000090921691909117905550600190505b949350505050565b60006106e7610be2565b60006106f38484610974565b6106fe906001611452565b67ffffffffffffffff808616600090815260056020908152604080832073ffffffffffffffffffffffffffffffffffffffff89168452909152902080549183167fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000090921691909117905591505092915050565b61077961078f565b6105a581610c29565b606060006105cd83610ced565b60015473ffffffffffffffffffffffffffffffffffffffff1633146107e0576040517f2b5c74de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b602081015160005b815181101561087d5760008282815181106108075761080761138a565b60200260200101519050610825816002610d4990919063ffffffff16565b156108745760405173ffffffffffffffffffffffffffffffffffffffff821681527fc3803387881faad271c47728894e3e36fac830ffc8602ca6fc07733cbda775809060200160405180910390a15b506001016107ea565b50815160005b815181101561096e5760008282815181106108a0576108a061138a565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610910576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61091b600282610d6b565b5060405173ffffffffffffffffffffffffffffffffffffffff821681527feb1b9b92e50b7f88f9ff25d56765095ac6e91540eee214906f4036a908ffbdef9060200160405180910390a150600101610883565b50505050565b67ffffffffffffffff808316600090815260056020908152604080832073ffffffffffffffffffffffffffffffffffffffff861684529091528120549091168082036105b45767ffffffffffffffff841660009081526004602052604090205473ffffffffffffffffffffffffffffffffffffffff168015610a89576040517f856c824700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff858116600483015282169063856c824790602401602060405180830381865afa158015610a5c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a8091906114ef565b925050506105b7565b509392505050565b67ffffffffffffffff83166000908152600660205260408082209051829190610abd90869086906114df565b9081526040519081900360200190205467ffffffffffffffff16905060008190036105ca5767ffffffffffffffff851660009081526004602052604090206001015473ffffffffffffffffffffffffffffffffffffffff168015610bd95773ffffffffffffffffffffffffffffffffffffffff811663856c8247610b438688018861136d565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401602060405180830381865afa158015610bac573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bd091906114ef565b925050506105cd565b50949350505050565b610bed600233610d8d565b6107e0576040517fd86ad9cf00000000000000000000000000000000000000000000000000000000815233600482015260240160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821603610c78576040517fdad89dca00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b606081600001805480602002602001604051908101604052809291908181526020018280548015610d3d57602002820191906000526020600020905b815481526020019060010190808311610d29575b50505050509050919050565b60006105b48373ffffffffffffffffffffffffffffffffffffffff8416610dbc565b60006105b48373ffffffffffffffffffffffffffffffffffffffff8416610eb6565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415156105b4565b60008181526001830160205260408120548015610ea5576000610de060018361150c565b8554909150600090610df49060019061150c565b9050818114610e59576000866000018281548110610e1457610e1461138a565b9060005260206000200154905080876000018481548110610e3757610e3761138a565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080610e6a57610e6a61151f565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506105b7565b60009150506105b7565b5092915050565b6000818152600183016020526040812054610efd575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556105b7565b5060006105b7565b60006020808352835180602085015260005b81811015610f3357858101830151858201604001528201610f17565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b6020808252825182820181905260009190848201906040850190845b81811015610fc057835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101610f8e565b50909695505050505050565b67ffffffffffffffff811681146105a557600080fd5b600060208284031215610ff457600080fd5b81356105b481610fcc565b6000806020838503121561101257600080fd5b823567ffffffffffffffff8082111561102a57600080fd5b818501915085601f83011261103e57600080fd5b81358181111561104d57600080fd5b8660208260071b850101111561106257600080fd5b60209290920196919550909350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff811681146105a557600080fd5b600082601f8301126110d657600080fd5b8135602067ffffffffffffffff808311156110f3576110f3611074565b8260051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f8301168101818110848211171561113657611136611074565b604052938452602081870181019490810192508785111561115657600080fd5b6020870191505b8482101561117f578135611170816110a3565b8352918301919083019061115d565b979650505050505050565b60006020828403121561119c57600080fd5b813567ffffffffffffffff808211156111b457600080fd5b90830190604082860312156111c857600080fd5b6040516040810181811083821117156111e3576111e3611074565b6040528235828111156111f557600080fd5b611201878286016110c5565b82525060208301358281111561121657600080fd5b611222878286016110c5565b60208301525095945050505050565b6000806040838503121561124457600080fd5b823561124f81610fcc565b9150602083013561125f816110a3565b809150509250929050565b60008083601f84011261127c57600080fd5b50813567ffffffffffffffff81111561129457600080fd5b6020830191508360208285010111156112ac57600080fd5b9250929050565b6000806000604084860312156112c857600080fd5b83356112d381610fcc565b9250602084013567ffffffffffffffff8111156112ef57600080fd5b6112fb8682870161126a565b9497909650939450505050565b6000806000806060858703121561131e57600080fd5b843561132981610fcc565b9350602085013561133981610fcc565b9250604085013567ffffffffffffffff81111561135557600080fd5b6113618782880161126a565b95989497509550505050565b60006020828403121561137f57600080fd5b81356105b4816110a3565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6000602082840312156113cb57600080fd5b813580151581146105b457600080fd5b6040810182356113ea816110a3565b73ffffffffffffffffffffffffffffffffffffffff9081168352602084013590611413826110a3565b8082166020850152505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b67ffffffffffffffff818116838216019080821115610eaf57610eaf611423565b600067ffffffffffffffff8087168352808616602084015250606060408301528260608301528284608084013760006080848401015260807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f850116830101905095945050505050565b8183823760009101908152919050565b60006020828403121561150157600080fd5b81516105b481610fcc565b818103818111156105b7576105b7611423565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000818000a",
}

var NonceManagerABI = NonceManagerMetaData.ABI

var NonceManagerBin = NonceManagerMetaData.Bin

func DeployNonceManager(auth *bind.TransactOpts, backend bind.ContractBackend, authorizedCallers []common.Address) (common.Address, *types.Transaction, *NonceManager, error) {
	parsed, err := NonceManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NonceManagerBin), backend, authorizedCallers)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NonceManager{address: address, abi: *parsed, NonceManagerCaller: NonceManagerCaller{contract: contract}, NonceManagerTransactor: NonceManagerTransactor{contract: contract}, NonceManagerFilterer: NonceManagerFilterer{contract: contract}}, nil
}

type NonceManager struct {
	address common.Address
	abi     abi.ABI
	NonceManagerCaller
	NonceManagerTransactor
	NonceManagerFilterer
}

type NonceManagerCaller struct {
	contract *bind.BoundContract
}

type NonceManagerTransactor struct {
	contract *bind.BoundContract
}

type NonceManagerFilterer struct {
	contract *bind.BoundContract
}

type NonceManagerSession struct {
	Contract     *NonceManager
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type NonceManagerCallerSession struct {
	Contract *NonceManagerCaller
	CallOpts bind.CallOpts
}

type NonceManagerTransactorSession struct {
	Contract     *NonceManagerTransactor
	TransactOpts bind.TransactOpts
}

type NonceManagerRaw struct {
	Contract *NonceManager
}

type NonceManagerCallerRaw struct {
	Contract *NonceManagerCaller
}

type NonceManagerTransactorRaw struct {
	Contract *NonceManagerTransactor
}

func NewNonceManager(address common.Address, backend bind.ContractBackend) (*NonceManager, error) {
	abi, err := abi.JSON(strings.NewReader(NonceManagerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindNonceManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NonceManager{address: address, abi: abi, NonceManagerCaller: NonceManagerCaller{contract: contract}, NonceManagerTransactor: NonceManagerTransactor{contract: contract}, NonceManagerFilterer: NonceManagerFilterer{contract: contract}}, nil
}

func NewNonceManagerCaller(address common.Address, caller bind.ContractCaller) (*NonceManagerCaller, error) {
	contract, err := bindNonceManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NonceManagerCaller{contract: contract}, nil
}

func NewNonceManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*NonceManagerTransactor, error) {
	contract, err := bindNonceManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NonceManagerTransactor{contract: contract}, nil
}

func NewNonceManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*NonceManagerFilterer, error) {
	contract, err := bindNonceManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NonceManagerFilterer{contract: contract}, nil
}

func bindNonceManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NonceManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_NonceManager *NonceManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NonceManager.Contract.NonceManagerCaller.contract.Call(opts, result, method, params...)
}

func (_NonceManager *NonceManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NonceManager.Contract.NonceManagerTransactor.contract.Transfer(opts)
}

func (_NonceManager *NonceManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NonceManager.Contract.NonceManagerTransactor.contract.Transact(opts, method, params...)
}

func (_NonceManager *NonceManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NonceManager.Contract.contract.Call(opts, result, method, params...)
}

func (_NonceManager *NonceManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NonceManager.Contract.contract.Transfer(opts)
}

func (_NonceManager *NonceManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NonceManager.Contract.contract.Transact(opts, method, params...)
}

func (_NonceManager *NonceManagerCaller) GetAllAuthorizedCallers(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _NonceManager.contract.Call(opts, &out, "getAllAuthorizedCallers")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_NonceManager *NonceManagerSession) GetAllAuthorizedCallers() ([]common.Address, error) {
	return _NonceManager.Contract.GetAllAuthorizedCallers(&_NonceManager.CallOpts)
}

func (_NonceManager *NonceManagerCallerSession) GetAllAuthorizedCallers() ([]common.Address, error) {
	return _NonceManager.Contract.GetAllAuthorizedCallers(&_NonceManager.CallOpts)
}

func (_NonceManager *NonceManagerCaller) GetInboundNonce(opts *bind.CallOpts, sourceChainSelector uint64, sender []byte) (uint64, error) {
	var out []interface{}
	err := _NonceManager.contract.Call(opts, &out, "getInboundNonce", sourceChainSelector, sender)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_NonceManager *NonceManagerSession) GetInboundNonce(sourceChainSelector uint64, sender []byte) (uint64, error) {
	return _NonceManager.Contract.GetInboundNonce(&_NonceManager.CallOpts, sourceChainSelector, sender)
}

func (_NonceManager *NonceManagerCallerSession) GetInboundNonce(sourceChainSelector uint64, sender []byte) (uint64, error) {
	return _NonceManager.Contract.GetInboundNonce(&_NonceManager.CallOpts, sourceChainSelector, sender)
}

func (_NonceManager *NonceManagerCaller) GetOutboundNonce(opts *bind.CallOpts, destChainSelector uint64, sender common.Address) (uint64, error) {
	var out []interface{}
	err := _NonceManager.contract.Call(opts, &out, "getOutboundNonce", destChainSelector, sender)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_NonceManager *NonceManagerSession) GetOutboundNonce(destChainSelector uint64, sender common.Address) (uint64, error) {
	return _NonceManager.Contract.GetOutboundNonce(&_NonceManager.CallOpts, destChainSelector, sender)
}

func (_NonceManager *NonceManagerCallerSession) GetOutboundNonce(destChainSelector uint64, sender common.Address) (uint64, error) {
	return _NonceManager.Contract.GetOutboundNonce(&_NonceManager.CallOpts, destChainSelector, sender)
}

func (_NonceManager *NonceManagerCaller) GetPreviousRamps(opts *bind.CallOpts, chainSelector uint64) (NonceManagerPreviousRamps, error) {
	var out []interface{}
	err := _NonceManager.contract.Call(opts, &out, "getPreviousRamps", chainSelector)

	if err != nil {
		return *new(NonceManagerPreviousRamps), err
	}

	out0 := *abi.ConvertType(out[0], new(NonceManagerPreviousRamps)).(*NonceManagerPreviousRamps)

	return out0, err

}

func (_NonceManager *NonceManagerSession) GetPreviousRamps(chainSelector uint64) (NonceManagerPreviousRamps, error) {
	return _NonceManager.Contract.GetPreviousRamps(&_NonceManager.CallOpts, chainSelector)
}

func (_NonceManager *NonceManagerCallerSession) GetPreviousRamps(chainSelector uint64) (NonceManagerPreviousRamps, error) {
	return _NonceManager.Contract.GetPreviousRamps(&_NonceManager.CallOpts, chainSelector)
}

func (_NonceManager *NonceManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NonceManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_NonceManager *NonceManagerSession) Owner() (common.Address, error) {
	return _NonceManager.Contract.Owner(&_NonceManager.CallOpts)
}

func (_NonceManager *NonceManagerCallerSession) Owner() (common.Address, error) {
	return _NonceManager.Contract.Owner(&_NonceManager.CallOpts)
}

func (_NonceManager *NonceManagerCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _NonceManager.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_NonceManager *NonceManagerSession) TypeAndVersion() (string, error) {
	return _NonceManager.Contract.TypeAndVersion(&_NonceManager.CallOpts)
}

func (_NonceManager *NonceManagerCallerSession) TypeAndVersion() (string, error) {
	return _NonceManager.Contract.TypeAndVersion(&_NonceManager.CallOpts)
}

func (_NonceManager *NonceManagerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NonceManager.contract.Transact(opts, "acceptOwnership")
}

func (_NonceManager *NonceManagerSession) AcceptOwnership() (*types.Transaction, error) {
	return _NonceManager.Contract.AcceptOwnership(&_NonceManager.TransactOpts)
}

func (_NonceManager *NonceManagerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _NonceManager.Contract.AcceptOwnership(&_NonceManager.TransactOpts)
}

func (_NonceManager *NonceManagerTransactor) ApplyAuthorizedCallerUpdates(opts *bind.TransactOpts, authorizedCallerArgs AuthorizedCallersAuthorizedCallerArgs) (*types.Transaction, error) {
	return _NonceManager.contract.Transact(opts, "applyAuthorizedCallerUpdates", authorizedCallerArgs)
}

func (_NonceManager *NonceManagerSession) ApplyAuthorizedCallerUpdates(authorizedCallerArgs AuthorizedCallersAuthorizedCallerArgs) (*types.Transaction, error) {
	return _NonceManager.Contract.ApplyAuthorizedCallerUpdates(&_NonceManager.TransactOpts, authorizedCallerArgs)
}

func (_NonceManager *NonceManagerTransactorSession) ApplyAuthorizedCallerUpdates(authorizedCallerArgs AuthorizedCallersAuthorizedCallerArgs) (*types.Transaction, error) {
	return _NonceManager.Contract.ApplyAuthorizedCallerUpdates(&_NonceManager.TransactOpts, authorizedCallerArgs)
}

func (_NonceManager *NonceManagerTransactor) ApplyPreviousRampsUpdates(opts *bind.TransactOpts, previousRampsArgs []NonceManagerPreviousRampsArgs) (*types.Transaction, error) {
	return _NonceManager.contract.Transact(opts, "applyPreviousRampsUpdates", previousRampsArgs)
}

func (_NonceManager *NonceManagerSession) ApplyPreviousRampsUpdates(previousRampsArgs []NonceManagerPreviousRampsArgs) (*types.Transaction, error) {
	return _NonceManager.Contract.ApplyPreviousRampsUpdates(&_NonceManager.TransactOpts, previousRampsArgs)
}

func (_NonceManager *NonceManagerTransactorSession) ApplyPreviousRampsUpdates(previousRampsArgs []NonceManagerPreviousRampsArgs) (*types.Transaction, error) {
	return _NonceManager.Contract.ApplyPreviousRampsUpdates(&_NonceManager.TransactOpts, previousRampsArgs)
}

func (_NonceManager *NonceManagerTransactor) GetIncrementedOutboundNonce(opts *bind.TransactOpts, destChainSelector uint64, sender common.Address) (*types.Transaction, error) {
	return _NonceManager.contract.Transact(opts, "getIncrementedOutboundNonce", destChainSelector, sender)
}

func (_NonceManager *NonceManagerSession) GetIncrementedOutboundNonce(destChainSelector uint64, sender common.Address) (*types.Transaction, error) {
	return _NonceManager.Contract.GetIncrementedOutboundNonce(&_NonceManager.TransactOpts, destChainSelector, sender)
}

func (_NonceManager *NonceManagerTransactorSession) GetIncrementedOutboundNonce(destChainSelector uint64, sender common.Address) (*types.Transaction, error) {
	return _NonceManager.Contract.GetIncrementedOutboundNonce(&_NonceManager.TransactOpts, destChainSelector, sender)
}

func (_NonceManager *NonceManagerTransactor) IncrementInboundNonce(opts *bind.TransactOpts, sourceChainSelector uint64, expectedNonce uint64, sender []byte) (*types.Transaction, error) {
	return _NonceManager.contract.Transact(opts, "incrementInboundNonce", sourceChainSelector, expectedNonce, sender)
}

func (_NonceManager *NonceManagerSession) IncrementInboundNonce(sourceChainSelector uint64, expectedNonce uint64, sender []byte) (*types.Transaction, error) {
	return _NonceManager.Contract.IncrementInboundNonce(&_NonceManager.TransactOpts, sourceChainSelector, expectedNonce, sender)
}

func (_NonceManager *NonceManagerTransactorSession) IncrementInboundNonce(sourceChainSelector uint64, expectedNonce uint64, sender []byte) (*types.Transaction, error) {
	return _NonceManager.Contract.IncrementInboundNonce(&_NonceManager.TransactOpts, sourceChainSelector, expectedNonce, sender)
}

func (_NonceManager *NonceManagerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _NonceManager.contract.Transact(opts, "transferOwnership", to)
}

func (_NonceManager *NonceManagerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _NonceManager.Contract.TransferOwnership(&_NonceManager.TransactOpts, to)
}

func (_NonceManager *NonceManagerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _NonceManager.Contract.TransferOwnership(&_NonceManager.TransactOpts, to)
}

type NonceManagerAuthorizedCallerAddedIterator struct {
	Event *NonceManagerAuthorizedCallerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NonceManagerAuthorizedCallerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NonceManagerAuthorizedCallerAdded)
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
		it.Event = new(NonceManagerAuthorizedCallerAdded)
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

func (it *NonceManagerAuthorizedCallerAddedIterator) Error() error {
	return it.fail
}

func (it *NonceManagerAuthorizedCallerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NonceManagerAuthorizedCallerAdded struct {
	Caller common.Address
	Raw    types.Log
}

func (_NonceManager *NonceManagerFilterer) FilterAuthorizedCallerAdded(opts *bind.FilterOpts) (*NonceManagerAuthorizedCallerAddedIterator, error) {

	logs, sub, err := _NonceManager.contract.FilterLogs(opts, "AuthorizedCallerAdded")
	if err != nil {
		return nil, err
	}
	return &NonceManagerAuthorizedCallerAddedIterator{contract: _NonceManager.contract, event: "AuthorizedCallerAdded", logs: logs, sub: sub}, nil
}

func (_NonceManager *NonceManagerFilterer) WatchAuthorizedCallerAdded(opts *bind.WatchOpts, sink chan<- *NonceManagerAuthorizedCallerAdded) (event.Subscription, error) {

	logs, sub, err := _NonceManager.contract.WatchLogs(opts, "AuthorizedCallerAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NonceManagerAuthorizedCallerAdded)
				if err := _NonceManager.contract.UnpackLog(event, "AuthorizedCallerAdded", log); err != nil {
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

func (_NonceManager *NonceManagerFilterer) ParseAuthorizedCallerAdded(log types.Log) (*NonceManagerAuthorizedCallerAdded, error) {
	event := new(NonceManagerAuthorizedCallerAdded)
	if err := _NonceManager.contract.UnpackLog(event, "AuthorizedCallerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NonceManagerAuthorizedCallerRemovedIterator struct {
	Event *NonceManagerAuthorizedCallerRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NonceManagerAuthorizedCallerRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NonceManagerAuthorizedCallerRemoved)
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
		it.Event = new(NonceManagerAuthorizedCallerRemoved)
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

func (it *NonceManagerAuthorizedCallerRemovedIterator) Error() error {
	return it.fail
}

func (it *NonceManagerAuthorizedCallerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NonceManagerAuthorizedCallerRemoved struct {
	Caller common.Address
	Raw    types.Log
}

func (_NonceManager *NonceManagerFilterer) FilterAuthorizedCallerRemoved(opts *bind.FilterOpts) (*NonceManagerAuthorizedCallerRemovedIterator, error) {

	logs, sub, err := _NonceManager.contract.FilterLogs(opts, "AuthorizedCallerRemoved")
	if err != nil {
		return nil, err
	}
	return &NonceManagerAuthorizedCallerRemovedIterator{contract: _NonceManager.contract, event: "AuthorizedCallerRemoved", logs: logs, sub: sub}, nil
}

func (_NonceManager *NonceManagerFilterer) WatchAuthorizedCallerRemoved(opts *bind.WatchOpts, sink chan<- *NonceManagerAuthorizedCallerRemoved) (event.Subscription, error) {

	logs, sub, err := _NonceManager.contract.WatchLogs(opts, "AuthorizedCallerRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NonceManagerAuthorizedCallerRemoved)
				if err := _NonceManager.contract.UnpackLog(event, "AuthorizedCallerRemoved", log); err != nil {
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

func (_NonceManager *NonceManagerFilterer) ParseAuthorizedCallerRemoved(log types.Log) (*NonceManagerAuthorizedCallerRemoved, error) {
	event := new(NonceManagerAuthorizedCallerRemoved)
	if err := _NonceManager.contract.UnpackLog(event, "AuthorizedCallerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NonceManagerOwnershipTransferRequestedIterator struct {
	Event *NonceManagerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NonceManagerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NonceManagerOwnershipTransferRequested)
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
		it.Event = new(NonceManagerOwnershipTransferRequested)
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

func (it *NonceManagerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *NonceManagerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NonceManagerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_NonceManager *NonceManagerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*NonceManagerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _NonceManager.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &NonceManagerOwnershipTransferRequestedIterator{contract: _NonceManager.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_NonceManager *NonceManagerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *NonceManagerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _NonceManager.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NonceManagerOwnershipTransferRequested)
				if err := _NonceManager.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_NonceManager *NonceManagerFilterer) ParseOwnershipTransferRequested(log types.Log) (*NonceManagerOwnershipTransferRequested, error) {
	event := new(NonceManagerOwnershipTransferRequested)
	if err := _NonceManager.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NonceManagerOwnershipTransferredIterator struct {
	Event *NonceManagerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NonceManagerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NonceManagerOwnershipTransferred)
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
		it.Event = new(NonceManagerOwnershipTransferred)
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

func (it *NonceManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *NonceManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NonceManagerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_NonceManager *NonceManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*NonceManagerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _NonceManager.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &NonceManagerOwnershipTransferredIterator{contract: _NonceManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_NonceManager *NonceManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *NonceManagerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _NonceManager.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NonceManagerOwnershipTransferred)
				if err := _NonceManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_NonceManager *NonceManagerFilterer) ParseOwnershipTransferred(log types.Log) (*NonceManagerOwnershipTransferred, error) {
	event := new(NonceManagerOwnershipTransferred)
	if err := _NonceManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NonceManagerPreviousRampsUpdatedIterator struct {
	Event *NonceManagerPreviousRampsUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NonceManagerPreviousRampsUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NonceManagerPreviousRampsUpdated)
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
		it.Event = new(NonceManagerPreviousRampsUpdated)
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

func (it *NonceManagerPreviousRampsUpdatedIterator) Error() error {
	return it.fail
}

func (it *NonceManagerPreviousRampsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NonceManagerPreviousRampsUpdated struct {
	RemoteChainSelector uint64
	PrevRamp            NonceManagerPreviousRamps
	Raw                 types.Log
}

func (_NonceManager *NonceManagerFilterer) FilterPreviousRampsUpdated(opts *bind.FilterOpts, remoteChainSelector []uint64) (*NonceManagerPreviousRampsUpdatedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _NonceManager.contract.FilterLogs(opts, "PreviousRampsUpdated", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &NonceManagerPreviousRampsUpdatedIterator{contract: _NonceManager.contract, event: "PreviousRampsUpdated", logs: logs, sub: sub}, nil
}

func (_NonceManager *NonceManagerFilterer) WatchPreviousRampsUpdated(opts *bind.WatchOpts, sink chan<- *NonceManagerPreviousRampsUpdated, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _NonceManager.contract.WatchLogs(opts, "PreviousRampsUpdated", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NonceManagerPreviousRampsUpdated)
				if err := _NonceManager.contract.UnpackLog(event, "PreviousRampsUpdated", log); err != nil {
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

func (_NonceManager *NonceManagerFilterer) ParsePreviousRampsUpdated(log types.Log) (*NonceManagerPreviousRampsUpdated, error) {
	event := new(NonceManagerPreviousRampsUpdated)
	if err := _NonceManager.contract.UnpackLog(event, "PreviousRampsUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NonceManagerSkippedIncorrectNonceIterator struct {
	Event *NonceManagerSkippedIncorrectNonce

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NonceManagerSkippedIncorrectNonceIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NonceManagerSkippedIncorrectNonce)
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
		it.Event = new(NonceManagerSkippedIncorrectNonce)
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

func (it *NonceManagerSkippedIncorrectNonceIterator) Error() error {
	return it.fail
}

func (it *NonceManagerSkippedIncorrectNonceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NonceManagerSkippedIncorrectNonce struct {
	SourceChainSelector uint64
	Nonce               uint64
	Sender              []byte
	Raw                 types.Log
}

func (_NonceManager *NonceManagerFilterer) FilterSkippedIncorrectNonce(opts *bind.FilterOpts) (*NonceManagerSkippedIncorrectNonceIterator, error) {

	logs, sub, err := _NonceManager.contract.FilterLogs(opts, "SkippedIncorrectNonce")
	if err != nil {
		return nil, err
	}
	return &NonceManagerSkippedIncorrectNonceIterator{contract: _NonceManager.contract, event: "SkippedIncorrectNonce", logs: logs, sub: sub}, nil
}

func (_NonceManager *NonceManagerFilterer) WatchSkippedIncorrectNonce(opts *bind.WatchOpts, sink chan<- *NonceManagerSkippedIncorrectNonce) (event.Subscription, error) {

	logs, sub, err := _NonceManager.contract.WatchLogs(opts, "SkippedIncorrectNonce")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NonceManagerSkippedIncorrectNonce)
				if err := _NonceManager.contract.UnpackLog(event, "SkippedIncorrectNonce", log); err != nil {
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

func (_NonceManager *NonceManagerFilterer) ParseSkippedIncorrectNonce(log types.Log) (*NonceManagerSkippedIncorrectNonce, error) {
	event := new(NonceManagerSkippedIncorrectNonce)
	if err := _NonceManager.contract.UnpackLog(event, "SkippedIncorrectNonce", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_NonceManager *NonceManager) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _NonceManager.abi.Events["AuthorizedCallerAdded"].ID:
		return _NonceManager.ParseAuthorizedCallerAdded(log)
	case _NonceManager.abi.Events["AuthorizedCallerRemoved"].ID:
		return _NonceManager.ParseAuthorizedCallerRemoved(log)
	case _NonceManager.abi.Events["OwnershipTransferRequested"].ID:
		return _NonceManager.ParseOwnershipTransferRequested(log)
	case _NonceManager.abi.Events["OwnershipTransferred"].ID:
		return _NonceManager.ParseOwnershipTransferred(log)
	case _NonceManager.abi.Events["PreviousRampsUpdated"].ID:
		return _NonceManager.ParsePreviousRampsUpdated(log)
	case _NonceManager.abi.Events["SkippedIncorrectNonce"].ID:
		return _NonceManager.ParseSkippedIncorrectNonce(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (NonceManagerAuthorizedCallerAdded) Topic() common.Hash {
	return common.HexToHash("0xeb1b9b92e50b7f88f9ff25d56765095ac6e91540eee214906f4036a908ffbdef")
}

func (NonceManagerAuthorizedCallerRemoved) Topic() common.Hash {
	return common.HexToHash("0xc3803387881faad271c47728894e3e36fac830ffc8602ca6fc07733cbda77580")
}

func (NonceManagerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (NonceManagerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (NonceManagerPreviousRampsUpdated) Topic() common.Hash {
	return common.HexToHash("0xa2e43edcbc4fd175ae4bebbe3fd6139871ed1f1783cd4a1ace59b90d302c3319")
}

func (NonceManagerSkippedIncorrectNonce) Topic() common.Hash {
	return common.HexToHash("0x606ff8179e5e3c059b82df931acc496b7b6053e8879042f8267f930e0595f69f")
}

func (_NonceManager *NonceManager) Address() common.Address {
	return _NonceManager.address
}

type NonceManagerInterface interface {
	GetAllAuthorizedCallers(opts *bind.CallOpts) ([]common.Address, error)

	GetInboundNonce(opts *bind.CallOpts, sourceChainSelector uint64, sender []byte) (uint64, error)

	GetOutboundNonce(opts *bind.CallOpts, destChainSelector uint64, sender common.Address) (uint64, error)

	GetPreviousRamps(opts *bind.CallOpts, chainSelector uint64) (NonceManagerPreviousRamps, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ApplyAuthorizedCallerUpdates(opts *bind.TransactOpts, authorizedCallerArgs AuthorizedCallersAuthorizedCallerArgs) (*types.Transaction, error)

	ApplyPreviousRampsUpdates(opts *bind.TransactOpts, previousRampsArgs []NonceManagerPreviousRampsArgs) (*types.Transaction, error)

	GetIncrementedOutboundNonce(opts *bind.TransactOpts, destChainSelector uint64, sender common.Address) (*types.Transaction, error)

	IncrementInboundNonce(opts *bind.TransactOpts, sourceChainSelector uint64, expectedNonce uint64, sender []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterAuthorizedCallerAdded(opts *bind.FilterOpts) (*NonceManagerAuthorizedCallerAddedIterator, error)

	WatchAuthorizedCallerAdded(opts *bind.WatchOpts, sink chan<- *NonceManagerAuthorizedCallerAdded) (event.Subscription, error)

	ParseAuthorizedCallerAdded(log types.Log) (*NonceManagerAuthorizedCallerAdded, error)

	FilterAuthorizedCallerRemoved(opts *bind.FilterOpts) (*NonceManagerAuthorizedCallerRemovedIterator, error)

	WatchAuthorizedCallerRemoved(opts *bind.WatchOpts, sink chan<- *NonceManagerAuthorizedCallerRemoved) (event.Subscription, error)

	ParseAuthorizedCallerRemoved(log types.Log) (*NonceManagerAuthorizedCallerRemoved, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*NonceManagerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *NonceManagerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*NonceManagerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*NonceManagerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *NonceManagerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*NonceManagerOwnershipTransferred, error)

	FilterPreviousRampsUpdated(opts *bind.FilterOpts, remoteChainSelector []uint64) (*NonceManagerPreviousRampsUpdatedIterator, error)

	WatchPreviousRampsUpdated(opts *bind.WatchOpts, sink chan<- *NonceManagerPreviousRampsUpdated, remoteChainSelector []uint64) (event.Subscription, error)

	ParsePreviousRampsUpdated(log types.Log) (*NonceManagerPreviousRampsUpdated, error)

	FilterSkippedIncorrectNonce(opts *bind.FilterOpts) (*NonceManagerSkippedIncorrectNonceIterator, error)

	WatchSkippedIncorrectNonce(opts *bind.WatchOpts, sink chan<- *NonceManagerSkippedIncorrectNonce) (event.Subscription, error)

	ParseSkippedIncorrectNonce(log types.Log) (*NonceManagerSkippedIncorrectNonce, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
