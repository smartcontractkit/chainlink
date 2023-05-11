// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_router

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

var VRFRouterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"CoordinatorAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CoordinatorNotRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestID\",\"type\":\"uint256\"}],\"name\":\"RedemptionRouteNotFound\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"route\",\"type\":\"address\"}],\"name\":\"RouteNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnexpectedMigrationVersion\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"RouteSet\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"callWithExactGasEvenIfTargetIsNoContract\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"sufficientGas\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"deregisterCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCoordinators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"name\":\"getFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"name\":\"getFulfillmentFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"}],\"name\":\"getRoute\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"isCoordinatorRegistered\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestID\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"name\":\"redeemRandomness\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"randomness\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"registerCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"uint24\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"name\":\"requestRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"uint24\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"name\":\"requestRandomnessFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"}],\"name\":\"resetRoute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_redemptionRoutes\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"}],\"name\":\"setRoute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b611509806101576000396000f3fe608060405234801561001057600080fd5b50600436106101365760003560e01c80639e201036116100b2578063d9c4a44d11610081578063f2fde38b11610066578063f2fde38b146102f4578063f9c45ced14610307578063fa66358a1461031a57600080fd5b8063d9c4a44d146102ce578063db972c8b146102e157600080fd5b80639e20103614610265578063acfc6cdd14610278578063bb9b2b3414610298578063cef40aa6146102ab57600080fd5b80635d5d8d191161010957806379ba5097116100ee57806379ba5097146102395780638da5cb5b146102415780639479b74e1461025257600080fd5b80635d5d8d19146101e35780637612e884146101f857600080fd5b8063181f5a771461013b5780634b2407d4146101835780634ffac83a1461019857806355fe9763146101b9575b600080fd5b604080518082018252600f81527f565246526f7574657220312e302e3000000000000000000000000000000000006020820152905161017a9190610e82565b60405180910390f35b610196610191366004610eb1565b61032d565b005b6101ab6101a6366004610fa8565b610467565b60405190815260200161017a565b6101cc6101c7366004611010565b610525565b60408051921515835290151560208301520161017a565b6101eb61058e565b60405161017a9190611067565b6102216102063660046110b4565b6005602052600090815260409020546001600160a01b031681565b6040516001600160a01b03909116815260200161017a565b61019661059f565b6000546001600160a01b0316610221565b6101966102603660046110b4565b610655565b6101ab6102733660046110e1565b6106da565b61028b610286366004611153565b61077f565b60405161017a919061118d565b6101966102a63660046110b4565b610868565b6102be6102b9366004610eb1565b6108de565b604051901515815260200161017a565b6101966102dc366004610eb1565b6108f1565b6101ab6102ef3660046111c5565b610967565b610196610302366004610eb1565b610a16565b6101ab610315366004611263565b610a2a565b6102216103283660046110b4565b610ac9565b610335610b3d565b610340600382610b99565b15610377576040517fdcecb7bf00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000819050806001600160a01b031663294daa496040518163ffffffff1660e01b8152600401602060405180830381865afa1580156103ba573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103de91906112aa565b60ff1660000361041a576040517f6e81684800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610425600383610bbe565b506040516001600160a01b03831681527fb7cabbfc11e66731fc77de0444614282023bcbd41d16781c753a431d0af01625906020015b60405180910390a15050565b60008061047386610ac9565b90506000816001600160a01b03166362f8b62033898989896040518663ffffffff1660e01b81526004016104ab9594939291906112cd565b6020604051808303816000875af11580156104ca573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104ee9190611315565b600081815260056020526040902080546001600160a01b0319166001600160a01b0394909416939093179092555095945050505050565b60008033610534600382610b99565b610551576040516301fd70a160e51b815260040160405180910390fd5b5a611388811061058457611388810390508660408204820311156105845760008086516020880160008a8cf19350600192505b5050935093915050565b606061059a6003610bd3565b905090565b6001546001600160a01b031633146105fe5760405162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b33610661600382610b99565b61067e576040516301fd70a160e51b815260040160405180910390fd5b60008281526002602090815260409182902080546001600160a01b03191633908117909155915191825283917fd8bb190f51a471c0cc88b5df464ee10c29138cce0d70fb36c472d60b414f3b3a91015b60405180910390a25050565b6000806106e686610ac9565b6040517f9e2010360000000000000000000000000000000000000000000000000000000081529091506001600160a01b03821690639e2010369061073490899089908990899060040161132e565b602060405180830381865afa158015610751573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107759190611315565b9695505050505050565b6000828152600560205260409020546060906001600160a01b0316806107d4576040517f71e8b529000000000000000000000000000000000000000000000000000000008152600481018590526024016105f5565b6040517f6ae5fb3b00000000000000000000000000000000000000000000000000000000815281906001600160a01b03821690636ae5fb3b906108219033908a908a908a90600401611365565b6000604051808303816000875af1158015610840573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526107759190810190611393565b33610874600382610b99565b610891576040516301fd70a160e51b815260040160405180910390fd5b600082815260026020908152604080832080546001600160a01b03191690555191825283917fd8bb190f51a471c0cc88b5df464ee10c29138cce0d70fb36c472d60b414f3b3a91016106ce565b60006108eb600383610b99565b92915050565b6108f9610b3d565b80610905600382610b99565b610922576040516301fd70a160e51b815260040160405180910390fd5b61092d600383610be0565b506040516001600160a01b03831681527ff80a1a97fd42251f3c33cda98635e7399253033a6774fe37cd3f650b5282af379060200161045b565b60008061097388610ac9565b6040517f95009f080000000000000000000000000000000000000000000000000000000081529091506001600160a01b038216906395009f08906109c79033908c908c908c908c908c908c9060040161142d565b6020604051808303816000875af11580156109e6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a0a9190611315565b98975050505050505050565b610a1e610b3d565b610a2781610bf5565b50565b600080610a3684610ac9565b6040517ff9c45ced0000000000000000000000000000000000000000000000000000000081529091506001600160a01b0382169063f9c45ced90610a809087908790600401611496565b602060405180830381865afa158015610a9d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ac19190611315565b949350505050565b6000818152600260205260408120546001600160a01b031680610b0a576040516339910a6760e01b81526001600160a01b03821660048201526024016105f5565b610b15600382610b99565b6108eb576040516339910a6760e01b81526001600160a01b03821660048201526024016105f5565b6000546001600160a01b03163314610b975760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016105f5565b565b6001600160a01b038116600090815260018301602052604081205415155b9392505050565b6000610bb7836001600160a01b038416610c9e565b60606000610bb783610ced565b6000610bb7836001600160a01b038416610d49565b336001600160a01b03821603610c4d5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016105f5565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000818152600183016020526040812054610ce5575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556108eb565b5060006108eb565b606081600001805480602002602001604051908101604052809291908181526020018280548015610d3d57602002820191906000526020600020905b815481526020019060010190808311610d29575b50505050509050919050565b60008181526001830160205260408120548015610e32576000610d6d6001836114af565b8554909150600090610d81906001906114af565b9050818114610de6576000866000018281548110610da157610da16114d0565b9060005260206000200154905080876000018481548110610dc457610dc46114d0565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080610df757610df76114e6565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506108eb565b60009150506108eb565b6000815180845260005b81811015610e6257602081850181015186830182015201610e46565b506000602082860101526020601f19601f83011685010191505092915050565b602081526000610bb76020830184610e3c565b80356001600160a01b0381168114610eac57600080fd5b919050565b600060208284031215610ec357600080fd5b610bb782610e95565b803561ffff81168114610eac57600080fd5b803562ffffff81168114610eac57600080fd5b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff81118282101715610f3057610f30610ef1565b604052919050565b600082601f830112610f4957600080fd5b813567ffffffffffffffff811115610f6357610f63610ef1565b610f76601f8201601f1916602001610f07565b818152846020838601011115610f8b57600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060808587031215610fbe57600080fd5b84359350610fce60208601610ecc565b9250610fdc60408601610ede565b9150606085013567ffffffffffffffff811115610ff857600080fd5b61100487828801610f38565b91505092959194509250565b60008060006060848603121561102557600080fd5b8335925061103560208501610e95565b9150604084013567ffffffffffffffff81111561105157600080fd5b61105d86828701610f38565b9150509250925092565b6020808252825182820181905260009190848201906040850190845b818110156110a85783516001600160a01b031683529284019291840191600101611083565b50909695505050505050565b6000602082840312156110c657600080fd5b5035919050565b803563ffffffff81168114610eac57600080fd5b600080600080608085870312156110f757600080fd5b84359350611107602086016110cd565b9250604085013567ffffffffffffffff8082111561112457600080fd5b61113088838901610f38565b9350606087013591508082111561114657600080fd5b5061100487828801610f38565b60008060006060848603121561116857600080fd5b8335925060208401359150604084013567ffffffffffffffff81111561105157600080fd5b6020808252825182820181905260009190848201906040850190845b818110156110a8578351835292840192918401916001016111a9565b60008060008060008060c087890312156111de57600080fd5b863595506111ee60208801610ecc565b94506111fc60408801610ede565b935061120a606088016110cd565b9250608087013567ffffffffffffffff8082111561122757600080fd5b6112338a838b01610f38565b935060a089013591508082111561124957600080fd5b5061125689828a01610f38565b9150509295509295509295565b6000806040838503121561127657600080fd5b82359150602083013567ffffffffffffffff81111561129457600080fd5b6112a085828601610f38565b9150509250929050565b6000602082840312156112bc57600080fd5b815160ff81168114610bb757600080fd5b6001600160a01b038616815284602082015261ffff8416604082015262ffffff8316606082015260a06080820152600061130a60a0830184610e3c565b979650505050505050565b60006020828403121561132757600080fd5b5051919050565b84815263ffffffff841660208201526080604082015260006113536080830185610e3c565b828103606084015261130a8185610e3c565b6001600160a01b03851681528360208201528260408201526080606082015260006107756080830184610e3c565b600060208083850312156113a657600080fd5b825167ffffffffffffffff808211156113be57600080fd5b818501915085601f8301126113d257600080fd5b8151818111156113e4576113e4610ef1565b8060051b91506113f5848301610f07565b818152918301840191848101908884111561140f57600080fd5b938501935b83851015610a0a57845182529385019390850190611414565b6001600160a01b038816815286602082015261ffff8616604082015262ffffff8516606082015263ffffffff8416608082015260e060a0820152600061147660e0830185610e3c565b82810360c08401526114888185610e3c565b9a9950505050505050505050565b828152604060208201526000610ac16040830184610e3c565b818103818111156108eb57634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052603160045260246000fdfea164736f6c6343000813000a",
}

var VRFRouterABI = VRFRouterMetaData.ABI

var VRFRouterBin = VRFRouterMetaData.Bin

func DeployVRFRouter(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *VRFRouter, error) {
	parsed, err := VRFRouterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFRouterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFRouter{VRFRouterCaller: VRFRouterCaller{contract: contract}, VRFRouterTransactor: VRFRouterTransactor{contract: contract}, VRFRouterFilterer: VRFRouterFilterer{contract: contract}}, nil
}

type VRFRouter struct {
	address common.Address
	abi     abi.ABI
	VRFRouterCaller
	VRFRouterTransactor
	VRFRouterFilterer
}

type VRFRouterCaller struct {
	contract *bind.BoundContract
}

type VRFRouterTransactor struct {
	contract *bind.BoundContract
}

type VRFRouterFilterer struct {
	contract *bind.BoundContract
}

type VRFRouterSession struct {
	Contract     *VRFRouter
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFRouterCallerSession struct {
	Contract *VRFRouterCaller
	CallOpts bind.CallOpts
}

type VRFRouterTransactorSession struct {
	Contract     *VRFRouterTransactor
	TransactOpts bind.TransactOpts
}

type VRFRouterRaw struct {
	Contract *VRFRouter
}

type VRFRouterCallerRaw struct {
	Contract *VRFRouterCaller
}

type VRFRouterTransactorRaw struct {
	Contract *VRFRouterTransactor
}

func NewVRFRouter(address common.Address, backend bind.ContractBackend) (*VRFRouter, error) {
	abi, err := abi.JSON(strings.NewReader(VRFRouterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFRouter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFRouter{address: address, abi: abi, VRFRouterCaller: VRFRouterCaller{contract: contract}, VRFRouterTransactor: VRFRouterTransactor{contract: contract}, VRFRouterFilterer: VRFRouterFilterer{contract: contract}}, nil
}

func NewVRFRouterCaller(address common.Address, caller bind.ContractCaller) (*VRFRouterCaller, error) {
	contract, err := bindVRFRouter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFRouterCaller{contract: contract}, nil
}

func NewVRFRouterTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFRouterTransactor, error) {
	contract, err := bindVRFRouter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFRouterTransactor{contract: contract}, nil
}

func NewVRFRouterFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFRouterFilterer, error) {
	contract, err := bindVRFRouter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFRouterFilterer{contract: contract}, nil
}

func bindVRFRouter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFRouterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFRouter *VRFRouterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFRouter.Contract.VRFRouterCaller.contract.Call(opts, result, method, params...)
}

func (_VRFRouter *VRFRouterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFRouter.Contract.VRFRouterTransactor.contract.Transfer(opts)
}

func (_VRFRouter *VRFRouterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFRouter.Contract.VRFRouterTransactor.contract.Transact(opts, method, params...)
}

func (_VRFRouter *VRFRouterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFRouter.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFRouter *VRFRouterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFRouter.Contract.contract.Transfer(opts)
}

func (_VRFRouter *VRFRouterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFRouter.Contract.contract.Transact(opts, method, params...)
}

func (_VRFRouter *VRFRouterCaller) GetCoordinators(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _VRFRouter.contract.Call(opts, &out, "getCoordinators")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_VRFRouter *VRFRouterSession) GetCoordinators() ([]common.Address, error) {
	return _VRFRouter.Contract.GetCoordinators(&_VRFRouter.CallOpts)
}

func (_VRFRouter *VRFRouterCallerSession) GetCoordinators() ([]common.Address, error) {
	return _VRFRouter.Contract.GetCoordinators(&_VRFRouter.CallOpts)
}

func (_VRFRouter *VRFRouterCaller) GetFee(opts *bind.CallOpts, subID *big.Int, extraArgs []byte) (*big.Int, error) {
	var out []interface{}
	err := _VRFRouter.contract.Call(opts, &out, "getFee", subID, extraArgs)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFRouter *VRFRouterSession) GetFee(subID *big.Int, extraArgs []byte) (*big.Int, error) {
	return _VRFRouter.Contract.GetFee(&_VRFRouter.CallOpts, subID, extraArgs)
}

func (_VRFRouter *VRFRouterCallerSession) GetFee(subID *big.Int, extraArgs []byte) (*big.Int, error) {
	return _VRFRouter.Contract.GetFee(&_VRFRouter.CallOpts, subID, extraArgs)
}

func (_VRFRouter *VRFRouterCaller) GetFulfillmentFee(opts *bind.CallOpts, subID *big.Int, callbackGasLimit uint32, arguments []byte, extraArgs []byte) (*big.Int, error) {
	var out []interface{}
	err := _VRFRouter.contract.Call(opts, &out, "getFulfillmentFee", subID, callbackGasLimit, arguments, extraArgs)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFRouter *VRFRouterSession) GetFulfillmentFee(subID *big.Int, callbackGasLimit uint32, arguments []byte, extraArgs []byte) (*big.Int, error) {
	return _VRFRouter.Contract.GetFulfillmentFee(&_VRFRouter.CallOpts, subID, callbackGasLimit, arguments, extraArgs)
}

func (_VRFRouter *VRFRouterCallerSession) GetFulfillmentFee(subID *big.Int, callbackGasLimit uint32, arguments []byte, extraArgs []byte) (*big.Int, error) {
	return _VRFRouter.Contract.GetFulfillmentFee(&_VRFRouter.CallOpts, subID, callbackGasLimit, arguments, extraArgs)
}

func (_VRFRouter *VRFRouterCaller) GetRoute(opts *bind.CallOpts, subID *big.Int) (common.Address, error) {
	var out []interface{}
	err := _VRFRouter.contract.Call(opts, &out, "getRoute", subID)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFRouter *VRFRouterSession) GetRoute(subID *big.Int) (common.Address, error) {
	return _VRFRouter.Contract.GetRoute(&_VRFRouter.CallOpts, subID)
}

func (_VRFRouter *VRFRouterCallerSession) GetRoute(subID *big.Int) (common.Address, error) {
	return _VRFRouter.Contract.GetRoute(&_VRFRouter.CallOpts, subID)
}

func (_VRFRouter *VRFRouterCaller) IsCoordinatorRegistered(opts *bind.CallOpts, coordinatorAddress common.Address) (bool, error) {
	var out []interface{}
	err := _VRFRouter.contract.Call(opts, &out, "isCoordinatorRegistered", coordinatorAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VRFRouter *VRFRouterSession) IsCoordinatorRegistered(coordinatorAddress common.Address) (bool, error) {
	return _VRFRouter.Contract.IsCoordinatorRegistered(&_VRFRouter.CallOpts, coordinatorAddress)
}

func (_VRFRouter *VRFRouterCallerSession) IsCoordinatorRegistered(coordinatorAddress common.Address) (bool, error) {
	return _VRFRouter.Contract.IsCoordinatorRegistered(&_VRFRouter.CallOpts, coordinatorAddress)
}

func (_VRFRouter *VRFRouterCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFRouter.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFRouter *VRFRouterSession) Owner() (common.Address, error) {
	return _VRFRouter.Contract.Owner(&_VRFRouter.CallOpts)
}

func (_VRFRouter *VRFRouterCallerSession) Owner() (common.Address, error) {
	return _VRFRouter.Contract.Owner(&_VRFRouter.CallOpts)
}

func (_VRFRouter *VRFRouterCaller) SRedemptionRoutes(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _VRFRouter.contract.Call(opts, &out, "s_redemptionRoutes", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFRouter *VRFRouterSession) SRedemptionRoutes(arg0 *big.Int) (common.Address, error) {
	return _VRFRouter.Contract.SRedemptionRoutes(&_VRFRouter.CallOpts, arg0)
}

func (_VRFRouter *VRFRouterCallerSession) SRedemptionRoutes(arg0 *big.Int) (common.Address, error) {
	return _VRFRouter.Contract.SRedemptionRoutes(&_VRFRouter.CallOpts, arg0)
}

func (_VRFRouter *VRFRouterCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VRFRouter.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VRFRouter *VRFRouterSession) TypeAndVersion() (string, error) {
	return _VRFRouter.Contract.TypeAndVersion(&_VRFRouter.CallOpts)
}

func (_VRFRouter *VRFRouterCallerSession) TypeAndVersion() (string, error) {
	return _VRFRouter.Contract.TypeAndVersion(&_VRFRouter.CallOpts)
}

func (_VRFRouter *VRFRouterTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFRouter.contract.Transact(opts, "acceptOwnership")
}

func (_VRFRouter *VRFRouterSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFRouter.Contract.AcceptOwnership(&_VRFRouter.TransactOpts)
}

func (_VRFRouter *VRFRouterTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFRouter.Contract.AcceptOwnership(&_VRFRouter.TransactOpts)
}

func (_VRFRouter *VRFRouterTransactor) CallWithExactGasEvenIfTargetIsNoContract(opts *bind.TransactOpts, gasAmount *big.Int, target common.Address, data []byte) (*types.Transaction, error) {
	return _VRFRouter.contract.Transact(opts, "callWithExactGasEvenIfTargetIsNoContract", gasAmount, target, data)
}

func (_VRFRouter *VRFRouterSession) CallWithExactGasEvenIfTargetIsNoContract(gasAmount *big.Int, target common.Address, data []byte) (*types.Transaction, error) {
	return _VRFRouter.Contract.CallWithExactGasEvenIfTargetIsNoContract(&_VRFRouter.TransactOpts, gasAmount, target, data)
}

func (_VRFRouter *VRFRouterTransactorSession) CallWithExactGasEvenIfTargetIsNoContract(gasAmount *big.Int, target common.Address, data []byte) (*types.Transaction, error) {
	return _VRFRouter.Contract.CallWithExactGasEvenIfTargetIsNoContract(&_VRFRouter.TransactOpts, gasAmount, target, data)
}

func (_VRFRouter *VRFRouterTransactor) DeregisterCoordinator(opts *bind.TransactOpts, coordinatorAddress common.Address) (*types.Transaction, error) {
	return _VRFRouter.contract.Transact(opts, "deregisterCoordinator", coordinatorAddress)
}

func (_VRFRouter *VRFRouterSession) DeregisterCoordinator(coordinatorAddress common.Address) (*types.Transaction, error) {
	return _VRFRouter.Contract.DeregisterCoordinator(&_VRFRouter.TransactOpts, coordinatorAddress)
}

func (_VRFRouter *VRFRouterTransactorSession) DeregisterCoordinator(coordinatorAddress common.Address) (*types.Transaction, error) {
	return _VRFRouter.Contract.DeregisterCoordinator(&_VRFRouter.TransactOpts, coordinatorAddress)
}

func (_VRFRouter *VRFRouterTransactor) RedeemRandomness(opts *bind.TransactOpts, subID *big.Int, requestID *big.Int, extraArgs []byte) (*types.Transaction, error) {
	return _VRFRouter.contract.Transact(opts, "redeemRandomness", subID, requestID, extraArgs)
}

func (_VRFRouter *VRFRouterSession) RedeemRandomness(subID *big.Int, requestID *big.Int, extraArgs []byte) (*types.Transaction, error) {
	return _VRFRouter.Contract.RedeemRandomness(&_VRFRouter.TransactOpts, subID, requestID, extraArgs)
}

func (_VRFRouter *VRFRouterTransactorSession) RedeemRandomness(subID *big.Int, requestID *big.Int, extraArgs []byte) (*types.Transaction, error) {
	return _VRFRouter.Contract.RedeemRandomness(&_VRFRouter.TransactOpts, subID, requestID, extraArgs)
}

func (_VRFRouter *VRFRouterTransactor) RegisterCoordinator(opts *bind.TransactOpts, coordinatorAddress common.Address) (*types.Transaction, error) {
	return _VRFRouter.contract.Transact(opts, "registerCoordinator", coordinatorAddress)
}

func (_VRFRouter *VRFRouterSession) RegisterCoordinator(coordinatorAddress common.Address) (*types.Transaction, error) {
	return _VRFRouter.Contract.RegisterCoordinator(&_VRFRouter.TransactOpts, coordinatorAddress)
}

func (_VRFRouter *VRFRouterTransactorSession) RegisterCoordinator(coordinatorAddress common.Address) (*types.Transaction, error) {
	return _VRFRouter.Contract.RegisterCoordinator(&_VRFRouter.TransactOpts, coordinatorAddress)
}

func (_VRFRouter *VRFRouterTransactor) RequestRandomness(opts *bind.TransactOpts, subID *big.Int, numWords uint16, confDelay *big.Int, extraArgs []byte) (*types.Transaction, error) {
	return _VRFRouter.contract.Transact(opts, "requestRandomness", subID, numWords, confDelay, extraArgs)
}

func (_VRFRouter *VRFRouterSession) RequestRandomness(subID *big.Int, numWords uint16, confDelay *big.Int, extraArgs []byte) (*types.Transaction, error) {
	return _VRFRouter.Contract.RequestRandomness(&_VRFRouter.TransactOpts, subID, numWords, confDelay, extraArgs)
}

func (_VRFRouter *VRFRouterTransactorSession) RequestRandomness(subID *big.Int, numWords uint16, confDelay *big.Int, extraArgs []byte) (*types.Transaction, error) {
	return _VRFRouter.Contract.RequestRandomness(&_VRFRouter.TransactOpts, subID, numWords, confDelay, extraArgs)
}

func (_VRFRouter *VRFRouterTransactor) RequestRandomnessFulfillment(opts *bind.TransactOpts, subID *big.Int, numWords uint16, confDelay *big.Int, callbackGasLimit uint32, arguments []byte, extraArgs []byte) (*types.Transaction, error) {
	return _VRFRouter.contract.Transact(opts, "requestRandomnessFulfillment", subID, numWords, confDelay, callbackGasLimit, arguments, extraArgs)
}

func (_VRFRouter *VRFRouterSession) RequestRandomnessFulfillment(subID *big.Int, numWords uint16, confDelay *big.Int, callbackGasLimit uint32, arguments []byte, extraArgs []byte) (*types.Transaction, error) {
	return _VRFRouter.Contract.RequestRandomnessFulfillment(&_VRFRouter.TransactOpts, subID, numWords, confDelay, callbackGasLimit, arguments, extraArgs)
}

func (_VRFRouter *VRFRouterTransactorSession) RequestRandomnessFulfillment(subID *big.Int, numWords uint16, confDelay *big.Int, callbackGasLimit uint32, arguments []byte, extraArgs []byte) (*types.Transaction, error) {
	return _VRFRouter.Contract.RequestRandomnessFulfillment(&_VRFRouter.TransactOpts, subID, numWords, confDelay, callbackGasLimit, arguments, extraArgs)
}

func (_VRFRouter *VRFRouterTransactor) ResetRoute(opts *bind.TransactOpts, subID *big.Int) (*types.Transaction, error) {
	return _VRFRouter.contract.Transact(opts, "resetRoute", subID)
}

func (_VRFRouter *VRFRouterSession) ResetRoute(subID *big.Int) (*types.Transaction, error) {
	return _VRFRouter.Contract.ResetRoute(&_VRFRouter.TransactOpts, subID)
}

func (_VRFRouter *VRFRouterTransactorSession) ResetRoute(subID *big.Int) (*types.Transaction, error) {
	return _VRFRouter.Contract.ResetRoute(&_VRFRouter.TransactOpts, subID)
}

func (_VRFRouter *VRFRouterTransactor) SetRoute(opts *bind.TransactOpts, subID *big.Int) (*types.Transaction, error) {
	return _VRFRouter.contract.Transact(opts, "setRoute", subID)
}

func (_VRFRouter *VRFRouterSession) SetRoute(subID *big.Int) (*types.Transaction, error) {
	return _VRFRouter.Contract.SetRoute(&_VRFRouter.TransactOpts, subID)
}

func (_VRFRouter *VRFRouterTransactorSession) SetRoute(subID *big.Int) (*types.Transaction, error) {
	return _VRFRouter.Contract.SetRoute(&_VRFRouter.TransactOpts, subID)
}

func (_VRFRouter *VRFRouterTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFRouter.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFRouter *VRFRouterSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFRouter.Contract.TransferOwnership(&_VRFRouter.TransactOpts, to)
}

func (_VRFRouter *VRFRouterTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFRouter.Contract.TransferOwnership(&_VRFRouter.TransactOpts, to)
}

type VRFRouterCoordinatorDeregisteredIterator struct {
	Event *VRFRouterCoordinatorDeregistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFRouterCoordinatorDeregisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFRouterCoordinatorDeregistered)
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
		it.Event = new(VRFRouterCoordinatorDeregistered)
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

func (it *VRFRouterCoordinatorDeregisteredIterator) Error() error {
	return it.fail
}

func (it *VRFRouterCoordinatorDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFRouterCoordinatorDeregistered struct {
	CoordinatorAddress common.Address
	Raw                types.Log
}

func (_VRFRouter *VRFRouterFilterer) FilterCoordinatorDeregistered(opts *bind.FilterOpts) (*VRFRouterCoordinatorDeregisteredIterator, error) {

	logs, sub, err := _VRFRouter.contract.FilterLogs(opts, "CoordinatorDeregistered")
	if err != nil {
		return nil, err
	}
	return &VRFRouterCoordinatorDeregisteredIterator{contract: _VRFRouter.contract, event: "CoordinatorDeregistered", logs: logs, sub: sub}, nil
}

func (_VRFRouter *VRFRouterFilterer) WatchCoordinatorDeregistered(opts *bind.WatchOpts, sink chan<- *VRFRouterCoordinatorDeregistered) (event.Subscription, error) {

	logs, sub, err := _VRFRouter.contract.WatchLogs(opts, "CoordinatorDeregistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFRouterCoordinatorDeregistered)
				if err := _VRFRouter.contract.UnpackLog(event, "CoordinatorDeregistered", log); err != nil {
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

func (_VRFRouter *VRFRouterFilterer) ParseCoordinatorDeregistered(log types.Log) (*VRFRouterCoordinatorDeregistered, error) {
	event := new(VRFRouterCoordinatorDeregistered)
	if err := _VRFRouter.contract.UnpackLog(event, "CoordinatorDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFRouterCoordinatorRegisteredIterator struct {
	Event *VRFRouterCoordinatorRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFRouterCoordinatorRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFRouterCoordinatorRegistered)
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
		it.Event = new(VRFRouterCoordinatorRegistered)
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

func (it *VRFRouterCoordinatorRegisteredIterator) Error() error {
	return it.fail
}

func (it *VRFRouterCoordinatorRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFRouterCoordinatorRegistered struct {
	CoordinatorAddress common.Address
	Raw                types.Log
}

func (_VRFRouter *VRFRouterFilterer) FilterCoordinatorRegistered(opts *bind.FilterOpts) (*VRFRouterCoordinatorRegisteredIterator, error) {

	logs, sub, err := _VRFRouter.contract.FilterLogs(opts, "CoordinatorRegistered")
	if err != nil {
		return nil, err
	}
	return &VRFRouterCoordinatorRegisteredIterator{contract: _VRFRouter.contract, event: "CoordinatorRegistered", logs: logs, sub: sub}, nil
}

func (_VRFRouter *VRFRouterFilterer) WatchCoordinatorRegistered(opts *bind.WatchOpts, sink chan<- *VRFRouterCoordinatorRegistered) (event.Subscription, error) {

	logs, sub, err := _VRFRouter.contract.WatchLogs(opts, "CoordinatorRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFRouterCoordinatorRegistered)
				if err := _VRFRouter.contract.UnpackLog(event, "CoordinatorRegistered", log); err != nil {
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

func (_VRFRouter *VRFRouterFilterer) ParseCoordinatorRegistered(log types.Log) (*VRFRouterCoordinatorRegistered, error) {
	event := new(VRFRouterCoordinatorRegistered)
	if err := _VRFRouter.contract.UnpackLog(event, "CoordinatorRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFRouterOwnershipTransferRequestedIterator struct {
	Event *VRFRouterOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFRouterOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFRouterOwnershipTransferRequested)
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
		it.Event = new(VRFRouterOwnershipTransferRequested)
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

func (it *VRFRouterOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFRouterOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFRouterOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFRouter *VRFRouterFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFRouterOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFRouter.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFRouterOwnershipTransferRequestedIterator{contract: _VRFRouter.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFRouter *VRFRouterFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFRouterOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFRouter.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFRouterOwnershipTransferRequested)
				if err := _VRFRouter.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFRouter *VRFRouterFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFRouterOwnershipTransferRequested, error) {
	event := new(VRFRouterOwnershipTransferRequested)
	if err := _VRFRouter.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFRouterOwnershipTransferredIterator struct {
	Event *VRFRouterOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFRouterOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFRouterOwnershipTransferred)
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
		it.Event = new(VRFRouterOwnershipTransferred)
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

func (it *VRFRouterOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFRouterOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFRouterOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFRouter *VRFRouterFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFRouterOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFRouter.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFRouterOwnershipTransferredIterator{contract: _VRFRouter.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFRouter *VRFRouterFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFRouterOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFRouter.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFRouterOwnershipTransferred)
				if err := _VRFRouter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFRouter *VRFRouterFilterer) ParseOwnershipTransferred(log types.Log) (*VRFRouterOwnershipTransferred, error) {
	event := new(VRFRouterOwnershipTransferred)
	if err := _VRFRouter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFRouterRouteSetIterator struct {
	Event *VRFRouterRouteSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFRouterRouteSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFRouterRouteSet)
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
		it.Event = new(VRFRouterRouteSet)
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

func (it *VRFRouterRouteSetIterator) Error() error {
	return it.fail
}

func (it *VRFRouterRouteSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFRouterRouteSet struct {
	SubID              *big.Int
	CoordinatorAddress common.Address
	Raw                types.Log
}

func (_VRFRouter *VRFRouterFilterer) FilterRouteSet(opts *bind.FilterOpts, subID []*big.Int) (*VRFRouterRouteSetIterator, error) {

	var subIDRule []interface{}
	for _, subIDItem := range subID {
		subIDRule = append(subIDRule, subIDItem)
	}

	logs, sub, err := _VRFRouter.contract.FilterLogs(opts, "RouteSet", subIDRule)
	if err != nil {
		return nil, err
	}
	return &VRFRouterRouteSetIterator{contract: _VRFRouter.contract, event: "RouteSet", logs: logs, sub: sub}, nil
}

func (_VRFRouter *VRFRouterFilterer) WatchRouteSet(opts *bind.WatchOpts, sink chan<- *VRFRouterRouteSet, subID []*big.Int) (event.Subscription, error) {

	var subIDRule []interface{}
	for _, subIDItem := range subID {
		subIDRule = append(subIDRule, subIDItem)
	}

	logs, sub, err := _VRFRouter.contract.WatchLogs(opts, "RouteSet", subIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFRouterRouteSet)
				if err := _VRFRouter.contract.UnpackLog(event, "RouteSet", log); err != nil {
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

func (_VRFRouter *VRFRouterFilterer) ParseRouteSet(log types.Log) (*VRFRouterRouteSet, error) {
	event := new(VRFRouterRouteSet)
	if err := _VRFRouter.contract.UnpackLog(event, "RouteSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_VRFRouter *VRFRouter) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFRouter.abi.Events["CoordinatorDeregistered"].ID:
		return _VRFRouter.ParseCoordinatorDeregistered(log)
	case _VRFRouter.abi.Events["CoordinatorRegistered"].ID:
		return _VRFRouter.ParseCoordinatorRegistered(log)
	case _VRFRouter.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFRouter.ParseOwnershipTransferRequested(log)
	case _VRFRouter.abi.Events["OwnershipTransferred"].ID:
		return _VRFRouter.ParseOwnershipTransferred(log)
	case _VRFRouter.abi.Events["RouteSet"].ID:
		return _VRFRouter.ParseRouteSet(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFRouterCoordinatorDeregistered) Topic() common.Hash {
	return common.HexToHash("0xf80a1a97fd42251f3c33cda98635e7399253033a6774fe37cd3f650b5282af37")
}

func (VRFRouterCoordinatorRegistered) Topic() common.Hash {
	return common.HexToHash("0xb7cabbfc11e66731fc77de0444614282023bcbd41d16781c753a431d0af01625")
}

func (VRFRouterOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFRouterOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VRFRouterRouteSet) Topic() common.Hash {
	return common.HexToHash("0xd8bb190f51a471c0cc88b5df464ee10c29138cce0d70fb36c472d60b414f3b3a")
}

func (_VRFRouter *VRFRouter) Address() common.Address {
	return _VRFRouter.address
}

type VRFRouterInterface interface {
	GetCoordinators(opts *bind.CallOpts) ([]common.Address, error)

	GetFee(opts *bind.CallOpts, subID *big.Int, extraArgs []byte) (*big.Int, error)

	GetFulfillmentFee(opts *bind.CallOpts, subID *big.Int, callbackGasLimit uint32, arguments []byte, extraArgs []byte) (*big.Int, error)

	GetRoute(opts *bind.CallOpts, subID *big.Int) (common.Address, error)

	IsCoordinatorRegistered(opts *bind.CallOpts, coordinatorAddress common.Address) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SRedemptionRoutes(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	CallWithExactGasEvenIfTargetIsNoContract(opts *bind.TransactOpts, gasAmount *big.Int, target common.Address, data []byte) (*types.Transaction, error)

	DeregisterCoordinator(opts *bind.TransactOpts, coordinatorAddress common.Address) (*types.Transaction, error)

	RedeemRandomness(opts *bind.TransactOpts, subID *big.Int, requestID *big.Int, extraArgs []byte) (*types.Transaction, error)

	RegisterCoordinator(opts *bind.TransactOpts, coordinatorAddress common.Address) (*types.Transaction, error)

	RequestRandomness(opts *bind.TransactOpts, subID *big.Int, numWords uint16, confDelay *big.Int, extraArgs []byte) (*types.Transaction, error)

	RequestRandomnessFulfillment(opts *bind.TransactOpts, subID *big.Int, numWords uint16, confDelay *big.Int, callbackGasLimit uint32, arguments []byte, extraArgs []byte) (*types.Transaction, error)

	ResetRoute(opts *bind.TransactOpts, subID *big.Int) (*types.Transaction, error)

	SetRoute(opts *bind.TransactOpts, subID *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterCoordinatorDeregistered(opts *bind.FilterOpts) (*VRFRouterCoordinatorDeregisteredIterator, error)

	WatchCoordinatorDeregistered(opts *bind.WatchOpts, sink chan<- *VRFRouterCoordinatorDeregistered) (event.Subscription, error)

	ParseCoordinatorDeregistered(log types.Log) (*VRFRouterCoordinatorDeregistered, error)

	FilterCoordinatorRegistered(opts *bind.FilterOpts) (*VRFRouterCoordinatorRegisteredIterator, error)

	WatchCoordinatorRegistered(opts *bind.WatchOpts, sink chan<- *VRFRouterCoordinatorRegistered) (event.Subscription, error)

	ParseCoordinatorRegistered(log types.Log) (*VRFRouterCoordinatorRegistered, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFRouterOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFRouterOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFRouterOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFRouterOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFRouterOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFRouterOwnershipTransferred, error)

	FilterRouteSet(opts *bind.FilterOpts, subID []*big.Int) (*VRFRouterRouteSetIterator, error)

	WatchRouteSet(opts *bind.WatchOpts, sink chan<- *VRFRouterRouteSet, subID []*big.Int) (event.Subscription, error)

	ParseRouteSet(log types.Log) (*VRFRouterRouteSet, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
