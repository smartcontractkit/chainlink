// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethereum

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// KeeperRegistrarMetaData contains all meta data concerning the KeeperRegistrar contract.
var KeeperRegistrarMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"LINKAddress\",\"type\":\"address\"},{\"internalType\":\"enumKeeperRegistrar.AutoApproveType\",\"name\":\"autoApproveConfigType\",\"type\":\"uint8\"},{\"internalType\":\"uint16\",\"name\":\"autoApproveMaxAllowed\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"minLINKJuels\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AmountMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FunctionNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HashMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientPayment\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAdminAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"LinkTransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyAdminOrOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistrationRequestFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RequestNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SenderMismatch\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"AutoApproveAllowedSenderSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"enumKeeperRegistrar.AutoApproveType\",\"name\":\"autoApproveConfigType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"autoApproveMaxAllowed\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"minLINKJuels\",\"type\":\"uint96\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"displayName\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"RegistrationApproved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"RegistrationRejected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encryptedEmail\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"upkeepContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"source\",\"type\":\"uint8\"}],\"name\":\"RegistrationRequested\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"upkeepContract\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"cancel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"}],\"name\":\"getAutoApproveAllowedSender\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"getPendingRequest\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRegistrationConfig\",\"outputs\":[{\"internalType\":\"enumKeeperRegistrar.AutoApproveType\",\"name\":\"autoApproveConfigType\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"autoApproveMaxAllowed\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"approvedCount\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minLINKJuels\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"encryptedEmail\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"upkeepContract\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"internalType\":\"uint8\",\"name\":\"source\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"setAutoApproveAllowedSender\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumKeeperRegistrar.AutoApproveType\",\"name\":\"autoApproveConfigType\",\"type\":\"uint8\"},{\"internalType\":\"uint16\",\"name\":\"autoApproveMaxAllowed\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"minLINKJuels\",\"type\":\"uint96\"}],\"name\":\"setRegistrationConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162001c7338038062001c7383398101604081905262000034916200038e565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000e7565b5050506001600160a01b038516608052620000dc8484848462000192565b505050505062000487565b336001600160a01b03821603620001415760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6200019c62000313565b6003546040805160a081019091526501000000000090910463ffffffff169080866002811115620001d157620001d16200041c565b815261ffff8616602082015263ffffffff831660408201526001600160a01b03851660608201526001600160601b038416608090910152805160038054909190829060ff191660018360028111156200022e576200022e6200041c565b0217905550602082015181546040808501516060860151610100600160481b031990931661010063ffffffff9586160263ffffffff60281b19161765010000000000949091169390930292909217600160481b600160e81b03191669010000000000000000006001600160a01b0390921691909102178255608090920151600190910180546001600160601b0319166001600160601b03909216919091179055517f6293a703ec7145dfa23c5cde2e627d6a02e153fc2e9c03b14d1e22cbb4a7e9cd906200030490879087908790879062000432565b60405180910390a15050505050565b6000546001600160a01b031633146200036f5760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000082565b565b80516001600160a01b03811681146200038957600080fd5b919050565b600080600080600060a08688031215620003a757600080fd5b620003b28662000371565b9450602086015160038110620003c757600080fd5b604087015190945061ffff81168114620003e057600080fd5b9250620003f06060870162000371565b60808701519092506001600160601b03811681146200040e57600080fd5b809150509295509295909350565b634e487b7160e01b600052602160045260246000fd5b60808101600386106200045557634e487b7160e01b600052602160045260246000fd5b94815261ffff9390931660208401526001600160a01b039190911660408301526001600160601b031660609091015290565b6080516117b4620004bf600039600081816101230152818161038c015281816107c801528181610c330152610dfa01526117b46000f3fe608060405234801561001057600080fd5b50600436106100ba5760003560e01c8063181f5a77146100bf578063183310b3146101095780631b6b6d231461011e5780633659d66614610152578063367b9b4f1461016557806379ba5097146101785780637e776f7f14610180578063850af0cb146101bc57806388b12d55146101d55780638da5cb5b14610234578063a4c0ed3614610245578063a793ab8b14610258578063c4d252f51461026b578063f2fde38b1461027e575b600080fd5b6100f36040518060400160405280601581526020017404b656570657252656769737472617220312e312e3605c1b81525081565b604051610100919061109d565b60405180910390f35b61011c6101173660046111d1565b610291565b005b6101457f000000000000000000000000000000000000000000000000000000000000000081565b6040516101009190611275565b61011c6101603660046112b1565b610381565b61011c6101733660046113b8565b6105e0565b61011c610647565b6101ac61018e3660046113ef565b6001600160a01b031660009081526005602052604090205460ff1690565b6040519015158152602001610100565b6101c46106f6565b604051610100959493929190611442565b6102266101e3366004611485565b6000908152600260209081526040918290208251808401909352546001600160a01b038116808452600160a01b9091046001600160601b03169290910182905291565b60405161010092919061149e565b6000546001600160a01b0316610145565b61011c6102533660046114c0565b6107bd565b61011c610266366004611519565b610a03565b61011c610279366004611485565b610b71565b61011c61028c3660046113ef565b610d05565b610299610d19565b6000818152600260209081526040918290208251808401909352546001600160a01b038116808452600160a01b9091046001600160601b0316918301919091526102f657604051632589d98f60e11b815260040160405180910390fd5b600087878787876040516020016103119594939291906115a5565b60405160208183030381529060405280519060200120905080831461034957604051633f4d605360e01b815260040160405180910390fd5b6000838152600260209081526040822091909155820151610376908a908a908a908a908a908a908a610d6e565b505050505050505050565b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146103c95760405162c6885f60e11b815260040160405180910390fd5b6001600160a01b0386166103f05760405163016ed19f60e21b815260040160405180910390fd5b6000888888888860405160200161040b9594939291906115a5565b6040516020818303038152906040528051906020012090508260ff16896001600160a01b0316827fc3f5df4aefec026f610a3fcb08f19476492d69d2cb78b1c2eba259a8820e6a788f8f8f8e8e8e8e8e60405161046f9897969594939291906115e9565b60405180910390a46040805160a08101909152600380546000929190829060ff1660028111156104a1576104a161140a565b60028111156104b2576104b261140a565b8152815463ffffffff610100820481166020840152600160281b82041660408301526001600160a01b03600160481b9091041660608201526001909101546001600160601b0316608090910152905061050b8184610f14565b1561055a576040810151610520906001611673565b6003805463ffffffff92909216600160281b0263ffffffff60281b199092169190911790556105558d8b8b8b8b8b8b89610d6e565b6105d1565b600082815260026020526040812054610584908790600160a01b90046001600160601b031661169b565b6040805180820182526001600160a01b03808d1682526001600160601b039384166020808401918252600089815260029091529390932091519251909316600160a01b0291909216179055505b50505050505050505050505050565b6105e8610d19565b6001600160a01b038216600081815260056020908152604091829020805460ff191685151590811790915591519182527f20c6237dac83526a849285a9f79d08a483291bdd3a056a0ef9ae94ecee1ad356910160405180910390a25050565b6001546001600160a01b0316331461069f5760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b60448201526064015b60405180910390fd5b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6040805160a081019091526003805460009283928392839283928392829060ff1660028111156107285761072861140a565b60028111156107395761073961140a565b81528154610100810463ffffffff908116602080850191909152600160281b8304909116604080850191909152600160481b9092046001600160a01b03166060808501919091526001909401546001600160601b0390811660809485015285519186015192860151948601519590930151909b919a50929850929650169350915050565b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146108055760405162c6885f60e11b815260040160405180910390fd5b81818080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050505060208101516001600160e01b03198116631b2ceb3360e11b146108715760405163e3d6792160e01b815260040160405180910390fd5b8484848080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050505060e48101518281146108cd576040516355e97b0d60e01b815260040160405180910390fd5b8887878080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152505050506101248101516001600160a01b038381169082161461093657604051637c62b1c760e11b815260040160405180910390fd5b61012489101561095957604051630dfe930960e41b815260040160405180910390fd5b6004546001600160601b03168b10156109855760405163cd1c886760e01b815260040160405180910390fd5b6000306001600160a01b03168b8b6040516109a19291906116bd565b600060405180830381855af49150503d80600081146109dc576040519150601f19603f3d011682016040523d82523d6000602084013e6109e1565b606091505b50509050806105d157604051630649bf8160e41b815260040160405180910390fd5b610a0b610d19565b6003546040805160a08101909152600160281b90910463ffffffff169080866002811115610a3b57610a3b61140a565b815261ffff8616602082015263ffffffff831660408201526001600160a01b03851660608201526001600160601b038416608090910152805160038054909190829060ff19166001836002811115610a9557610a9561140a565b021790555060208201518154604080850151606086015168ffffffffffffffff001990931661010063ffffffff9586160263ffffffff60281b191617600160281b949091169390930292909217600160481b600160e81b031916600160481b6001600160a01b0390921691909102178255608090920151600190910180546001600160601b0319166001600160601b03909216919091179055517f6293a703ec7145dfa23c5cde2e627d6a02e153fc2e9c03b14d1e22cbb4a7e9cd90610b629087908790879087906116cd565b60405180910390a15050505050565b6000818152600260209081526040918290208251808401909352546001600160a01b038116808452600160a01b9091046001600160601b031691830191909152331480610bc857506000546001600160a01b031633145b610be5576040516361685c2b60e01b815260040160405180910390fd5b80516001600160a01b0316610c0d57604051632589d98f60e11b815260040160405180910390fd5b600082815260026020908152604080832083905590830151905163a9059cbb60e01b81527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169163a9059cbb91610c7091339160040161149e565b6020604051808303816000875af1158015610c8f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610cb3919061170c565b905080610cd5573360405163185c9b9d60e31b81526004016106969190611275565b60405183907f3663fb28ebc87645eb972c9dad8521bf665c623f287e79f1c56f1eb374b82a2290600090a2505050565b610d0d610d19565b610d1681610fad565b50565b6000546001600160a01b03163314610d6c5760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b6044820152606401610696565b565b60035460405163da5c674160e01b8152600160481b9091046001600160a01b031690600090829063da5c674190610db1908c908c908c908c908c906004016115a5565b6020604051808303816000875af1158015610dd0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610df49190611729565b905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316634000aea0848785604051602001610e3c91815260200190565b6040516020818303038152906040526040518463ffffffff1660e01b8152600401610e6993929190611742565b6020604051808303816000875af1158015610e88573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610eac919061170c565b905080610ece578260405163185c9b9d60e31b81526004016106969190611275565b81847fb9a292fb7e3edd920cd2d2829a3615a640c43fd7de0a0820aa0668feb4c37d4b8d604051610eff919061109d565b60405180910390a35050505050505050505050565b60008083516002811115610f2a57610f2a61140a565b03610f3757506000610fa7565b600183516002811115610f4c57610f4c61140a565b148015610f7257506001600160a01b03821660009081526005602052604090205460ff16155b15610f7f57506000610fa7565b826020015163ffffffff16836040015163ffffffff161015610fa357506001610fa7565b5060005b92915050565b336001600160a01b03821603610fff5760405162461bcd60e51b815260206004820152601760248201527621b0b73737ba103a3930b739b332b9103a379039b2b63360491b6044820152606401610696565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000815180845260005b818110156110765760208185018101518683018201520161105a565b81811115611088576000602083870101525b50601f01601f19169290920160200192915050565b6020815260006110b06020830184611050565b9392505050565b634e487b7160e01b600052604160045260246000fd5b600082601f8301126110de57600080fd5b81356001600160401b03808211156110f8576110f86110b7565b604051601f8301601f19908116603f01168101908282118183101715611120576111206110b7565b8160405283815286602085880101111561113957600080fd5b836020870160208301376000602085830101528094505050505092915050565b80356001600160a01b038116811461117057600080fd5b919050565b803563ffffffff8116811461117057600080fd5b60008083601f84011261119b57600080fd5b5081356001600160401b038111156111b257600080fd5b6020830191508360208285010111156111ca57600080fd5b9250929050565b600080600080600080600060c0888a0312156111ec57600080fd5b87356001600160401b038082111561120357600080fd5b61120f8b838c016110cd565b985061121d60208b01611159565b975061122b60408b01611175565b965061123960608b01611159565b955060808a013591508082111561124f57600080fd5b5061125c8a828b01611189565b989b979a5095989497959660a090950135949350505050565b6001600160a01b0391909116815260200190565b80356001600160601b038116811461117057600080fd5b803560ff8116811461117057600080fd5b60008060008060008060008060008060006101208c8e0312156112d357600080fd5b6001600160401b03808d3511156112e957600080fd5b6112f68e8e358f016110cd565b9b508060208e0135111561130957600080fd5b6113198e60208f01358f01611189565b909b50995061132a60408e01611159565b985061133860608e01611175565b975061134660808e01611159565b96508060a08e0135111561135957600080fd5b5061136a8d60a08e01358e01611189565b909550935061137b60c08d01611289565b925061138960e08d016112a0565b91506113986101008d01611159565b90509295989b509295989b9093969950565b8015158114610d1657600080fd5b600080604083850312156113cb57600080fd5b6113d483611159565b915060208301356113e4816113aa565b809150509250929050565b60006020828403121561140157600080fd5b6110b082611159565b634e487b7160e01b600052602160045260246000fd5b6003811061143e57634e487b7160e01b600052602160045260246000fd5b9052565b60a081016114508288611420565b63ffffffff95861660208301529390941660408501526001600160a01b03919091166060840152608090920191909152919050565b60006020828403121561149757600080fd5b5035919050565b6001600160a01b039290921682526001600160601b0316602082015260400190565b600080600080606085870312156114d657600080fd5b6114df85611159565b93506020850135925060408501356001600160401b0381111561150157600080fd5b61150d87828801611189565b95989497509550505050565b6000806000806080858703121561152f57600080fd5b84356003811061153e57600080fd5b9350602085013561ffff8116811461155557600080fd5b925061156360408601611159565b915061157160608601611289565b905092959194509250565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b6001600160a01b03868116825263ffffffff86166020830152841660408201526080606082018190526000906115de908301848661157c565b979650505050505050565b60c0815260006115fc60c083018b611050565b828103602084015261160f818a8c61157c565b63ffffffff891660408501526001600160a01b03881660608501528381036080850152905061163f81868861157c565b91505060018060601b03831660a08301529998505050505050505050565b634e487b7160e01b600052601160045260246000fd5b600063ffffffff8083168185168083038211156116925761169261165d565b01949350505050565b60006001600160601b038281168482168083038211156116925761169261165d565b8183823760009101908152919050565b608081016116db8287611420565b61ffff9490941660208201526001600160a01b039290921660408301526001600160601b0316606090910152919050565b60006020828403121561171e57600080fd5b81516110b0816113aa565b60006020828403121561173b57600080fd5b5051919050565b6001600160a01b03841681526001600160601b038316602082015260606040820181905260009061177590830184611050565b9594505050505056fea26469706673582212203bb77b89ff16f1e3e7a2275da40c66a13535ba065d32c0993af15e25297561f464736f6c634300080d0033",
}

// KeeperRegistrarABI is the input ABI used to generate the binding from.
// Deprecated: Use KeeperRegistrarMetaData.ABI instead.
var KeeperRegistrarABI = KeeperRegistrarMetaData.ABI

// KeeperRegistrarBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use KeeperRegistrarMetaData.Bin instead.
var KeeperRegistrarBin = KeeperRegistrarMetaData.Bin

// DeployKeeperRegistrar deploys a new Ethereum contract, binding an instance of KeeperRegistrar to it.
func DeployKeeperRegistrar(auth *bind.TransactOpts, backend bind.ContractBackend, LINKAddress common.Address, autoApproveConfigType uint8, autoApproveMaxAllowed uint16, keeperRegistry common.Address, minLINKJuels *big.Int) (common.Address, *types.Transaction, *KeeperRegistrar, error) {
	parsed, err := KeeperRegistrarMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistrarBin), backend, LINKAddress, autoApproveConfigType, autoApproveMaxAllowed, keeperRegistry, minLINKJuels)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperRegistrar{KeeperRegistrarCaller: KeeperRegistrarCaller{contract: contract}, KeeperRegistrarTransactor: KeeperRegistrarTransactor{contract: contract}, KeeperRegistrarFilterer: KeeperRegistrarFilterer{contract: contract}}, nil
}

// KeeperRegistrar is an auto generated Go binding around an Ethereum contract.
type KeeperRegistrar struct {
	KeeperRegistrarCaller     // Read-only binding to the contract
	KeeperRegistrarTransactor // Write-only binding to the contract
	KeeperRegistrarFilterer   // Log filterer for contract events
}

// KeeperRegistrarCaller is an auto generated read-only Go binding around an Ethereum contract.
type KeeperRegistrarCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperRegistrarTransactor is an auto generated write-only Go binding around an Ethereum contract.
type KeeperRegistrarTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperRegistrarFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KeeperRegistrarFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperRegistrarSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KeeperRegistrarSession struct {
	Contract     *KeeperRegistrar  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KeeperRegistrarCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KeeperRegistrarCallerSession struct {
	Contract *KeeperRegistrarCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// KeeperRegistrarTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KeeperRegistrarTransactorSession struct {
	Contract     *KeeperRegistrarTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// KeeperRegistrarRaw is an auto generated low-level Go binding around an Ethereum contract.
type KeeperRegistrarRaw struct {
	Contract *KeeperRegistrar // Generic contract binding to access the raw methods on
}

// KeeperRegistrarCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KeeperRegistrarCallerRaw struct {
	Contract *KeeperRegistrarCaller // Generic read-only contract binding to access the raw methods on
}

// KeeperRegistrarTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KeeperRegistrarTransactorRaw struct {
	Contract *KeeperRegistrarTransactor // Generic write-only contract binding to access the raw methods on
}

// NewKeeperRegistrar creates a new instance of KeeperRegistrar, bound to a specific deployed contract.
func NewKeeperRegistrar(address common.Address, backend bind.ContractBackend) (*KeeperRegistrar, error) {
	contract, err := bindKeeperRegistrar(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrar{KeeperRegistrarCaller: KeeperRegistrarCaller{contract: contract}, KeeperRegistrarTransactor: KeeperRegistrarTransactor{contract: contract}, KeeperRegistrarFilterer: KeeperRegistrarFilterer{contract: contract}}, nil
}

// NewKeeperRegistrarCaller creates a new read-only instance of KeeperRegistrar, bound to a specific deployed contract.
func NewKeeperRegistrarCaller(address common.Address, caller bind.ContractCaller) (*KeeperRegistrarCaller, error) {
	contract, err := bindKeeperRegistrar(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarCaller{contract: contract}, nil
}

// NewKeeperRegistrarTransactor creates a new write-only instance of KeeperRegistrar, bound to a specific deployed contract.
func NewKeeperRegistrarTransactor(address common.Address, transactor bind.ContractTransactor) (*KeeperRegistrarTransactor, error) {
	contract, err := bindKeeperRegistrar(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarTransactor{contract: contract}, nil
}

// NewKeeperRegistrarFilterer creates a new log filterer instance of KeeperRegistrar, bound to a specific deployed contract.
func NewKeeperRegistrarFilterer(address common.Address, filterer bind.ContractFilterer) (*KeeperRegistrarFilterer, error) {
	contract, err := bindKeeperRegistrar(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarFilterer{contract: contract}, nil
}

// bindKeeperRegistrar binds a generic wrapper to an already deployed contract.
func bindKeeperRegistrar(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KeeperRegistrarABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeeperRegistrar *KeeperRegistrarRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistrar.Contract.KeeperRegistrarCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeeperRegistrar *KeeperRegistrarRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.KeeperRegistrarTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeeperRegistrar *KeeperRegistrarRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.KeeperRegistrarTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeeperRegistrar *KeeperRegistrarCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistrar.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeeperRegistrar *KeeperRegistrarTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeeperRegistrar *KeeperRegistrarTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.contract.Transact(opts, method, params...)
}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_KeeperRegistrar *KeeperRegistrarCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistrar.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_KeeperRegistrar *KeeperRegistrarSession) LINK() (common.Address, error) {
	return _KeeperRegistrar.Contract.LINK(&_KeeperRegistrar.CallOpts)
}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_KeeperRegistrar *KeeperRegistrarCallerSession) LINK() (common.Address, error) {
	return _KeeperRegistrar.Contract.LINK(&_KeeperRegistrar.CallOpts)
}

// GetAutoApproveAllowedSender is a free data retrieval call binding the contract method 0x7e776f7f.
//
// Solidity: function getAutoApproveAllowedSender(address senderAddress) view returns(bool)
func (_KeeperRegistrar *KeeperRegistrarCaller) GetAutoApproveAllowedSender(opts *bind.CallOpts, senderAddress common.Address) (bool, error) {
	var out []interface{}
	err := _KeeperRegistrar.contract.Call(opts, &out, "getAutoApproveAllowedSender", senderAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetAutoApproveAllowedSender is a free data retrieval call binding the contract method 0x7e776f7f.
//
// Solidity: function getAutoApproveAllowedSender(address senderAddress) view returns(bool)
func (_KeeperRegistrar *KeeperRegistrarSession) GetAutoApproveAllowedSender(senderAddress common.Address) (bool, error) {
	return _KeeperRegistrar.Contract.GetAutoApproveAllowedSender(&_KeeperRegistrar.CallOpts, senderAddress)
}

// GetAutoApproveAllowedSender is a free data retrieval call binding the contract method 0x7e776f7f.
//
// Solidity: function getAutoApproveAllowedSender(address senderAddress) view returns(bool)
func (_KeeperRegistrar *KeeperRegistrarCallerSession) GetAutoApproveAllowedSender(senderAddress common.Address) (bool, error) {
	return _KeeperRegistrar.Contract.GetAutoApproveAllowedSender(&_KeeperRegistrar.CallOpts, senderAddress)
}

// GetPendingRequest is a free data retrieval call binding the contract method 0x88b12d55.
//
// Solidity: function getPendingRequest(bytes32 hash) view returns(address, uint96)
func (_KeeperRegistrar *KeeperRegistrarCaller) GetPendingRequest(opts *bind.CallOpts, hash [32]byte) (common.Address, *big.Int, error) {
	var out []interface{}
	err := _KeeperRegistrar.contract.Call(opts, &out, "getPendingRequest", hash)

	if err != nil {
		return *new(common.Address), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetPendingRequest is a free data retrieval call binding the contract method 0x88b12d55.
//
// Solidity: function getPendingRequest(bytes32 hash) view returns(address, uint96)
func (_KeeperRegistrar *KeeperRegistrarSession) GetPendingRequest(hash [32]byte) (common.Address, *big.Int, error) {
	return _KeeperRegistrar.Contract.GetPendingRequest(&_KeeperRegistrar.CallOpts, hash)
}

// GetPendingRequest is a free data retrieval call binding the contract method 0x88b12d55.
//
// Solidity: function getPendingRequest(bytes32 hash) view returns(address, uint96)
func (_KeeperRegistrar *KeeperRegistrarCallerSession) GetPendingRequest(hash [32]byte) (common.Address, *big.Int, error) {
	return _KeeperRegistrar.Contract.GetPendingRequest(&_KeeperRegistrar.CallOpts, hash)
}

// GetRegistrationConfig is a free data retrieval call binding the contract method 0x850af0cb.
//
// Solidity: function getRegistrationConfig() view returns(uint8 autoApproveConfigType, uint32 autoApproveMaxAllowed, uint32 approvedCount, address keeperRegistry, uint256 minLINKJuels)
func (_KeeperRegistrar *KeeperRegistrarCaller) GetRegistrationConfig(opts *bind.CallOpts) (struct {
	AutoApproveConfigType uint8
	AutoApproveMaxAllowed uint32
	ApprovedCount         uint32
	KeeperRegistry        common.Address
	MinLINKJuels          *big.Int
}, error) {
	var out []interface{}
	err := _KeeperRegistrar.contract.Call(opts, &out, "getRegistrationConfig")

	outstruct := new(struct {
		AutoApproveConfigType uint8
		AutoApproveMaxAllowed uint32
		ApprovedCount         uint32
		KeeperRegistry        common.Address
		MinLINKJuels          *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.AutoApproveConfigType = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.AutoApproveMaxAllowed = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ApprovedCount = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.KeeperRegistry = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.MinLINKJuels = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetRegistrationConfig is a free data retrieval call binding the contract method 0x850af0cb.
//
// Solidity: function getRegistrationConfig() view returns(uint8 autoApproveConfigType, uint32 autoApproveMaxAllowed, uint32 approvedCount, address keeperRegistry, uint256 minLINKJuels)
func (_KeeperRegistrar *KeeperRegistrarSession) GetRegistrationConfig() (struct {
	AutoApproveConfigType uint8
	AutoApproveMaxAllowed uint32
	ApprovedCount         uint32
	KeeperRegistry        common.Address
	MinLINKJuels          *big.Int
}, error) {
	return _KeeperRegistrar.Contract.GetRegistrationConfig(&_KeeperRegistrar.CallOpts)
}

// GetRegistrationConfig is a free data retrieval call binding the contract method 0x850af0cb.
//
// Solidity: function getRegistrationConfig() view returns(uint8 autoApproveConfigType, uint32 autoApproveMaxAllowed, uint32 approvedCount, address keeperRegistry, uint256 minLINKJuels)
func (_KeeperRegistrar *KeeperRegistrarCallerSession) GetRegistrationConfig() (struct {
	AutoApproveConfigType uint8
	AutoApproveMaxAllowed uint32
	ApprovedCount         uint32
	KeeperRegistry        common.Address
	MinLINKJuels          *big.Int
}, error) {
	return _KeeperRegistrar.Contract.GetRegistrationConfig(&_KeeperRegistrar.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KeeperRegistrar *KeeperRegistrarCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistrar.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KeeperRegistrar *KeeperRegistrarSession) Owner() (common.Address, error) {
	return _KeeperRegistrar.Contract.Owner(&_KeeperRegistrar.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KeeperRegistrar *KeeperRegistrarCallerSession) Owner() (common.Address, error) {
	return _KeeperRegistrar.Contract.Owner(&_KeeperRegistrar.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_KeeperRegistrar *KeeperRegistrarCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _KeeperRegistrar.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_KeeperRegistrar *KeeperRegistrarSession) TypeAndVersion() (string, error) {
	return _KeeperRegistrar.Contract.TypeAndVersion(&_KeeperRegistrar.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_KeeperRegistrar *KeeperRegistrarCallerSession) TypeAndVersion() (string, error) {
	return _KeeperRegistrar.Contract.TypeAndVersion(&_KeeperRegistrar.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_KeeperRegistrar *KeeperRegistrarTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_KeeperRegistrar *KeeperRegistrarSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.AcceptOwnership(&_KeeperRegistrar.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_KeeperRegistrar *KeeperRegistrarTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.AcceptOwnership(&_KeeperRegistrar.TransactOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x183310b3.
//
// Solidity: function approve(string name, address upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, bytes32 hash) returns()
func (_KeeperRegistrar *KeeperRegistrarTransactor) Approve(opts *bind.TransactOpts, name string, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "approve", name, upkeepContract, gasLimit, adminAddress, checkData, hash)
}

// Approve is a paid mutator transaction binding the contract method 0x183310b3.
//
// Solidity: function approve(string name, address upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, bytes32 hash) returns()
func (_KeeperRegistrar *KeeperRegistrarSession) Approve(name string, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.Approve(&_KeeperRegistrar.TransactOpts, name, upkeepContract, gasLimit, adminAddress, checkData, hash)
}

// Approve is a paid mutator transaction binding the contract method 0x183310b3.
//
// Solidity: function approve(string name, address upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, bytes32 hash) returns()
func (_KeeperRegistrar *KeeperRegistrarTransactorSession) Approve(name string, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.Approve(&_KeeperRegistrar.TransactOpts, name, upkeepContract, gasLimit, adminAddress, checkData, hash)
}

// Cancel is a paid mutator transaction binding the contract method 0xc4d252f5.
//
// Solidity: function cancel(bytes32 hash) returns()
func (_KeeperRegistrar *KeeperRegistrarTransactor) Cancel(opts *bind.TransactOpts, hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "cancel", hash)
}

// Cancel is a paid mutator transaction binding the contract method 0xc4d252f5.
//
// Solidity: function cancel(bytes32 hash) returns()
func (_KeeperRegistrar *KeeperRegistrarSession) Cancel(hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.Cancel(&_KeeperRegistrar.TransactOpts, hash)
}

// Cancel is a paid mutator transaction binding the contract method 0xc4d252f5.
//
// Solidity: function cancel(bytes32 hash) returns()
func (_KeeperRegistrar *KeeperRegistrarTransactorSession) Cancel(hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.Cancel(&_KeeperRegistrar.TransactOpts, hash)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_KeeperRegistrar *KeeperRegistrarTransactor) OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "onTokenTransfer", sender, amount, data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_KeeperRegistrar *KeeperRegistrarSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.OnTokenTransfer(&_KeeperRegistrar.TransactOpts, sender, amount, data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_KeeperRegistrar *KeeperRegistrarTransactorSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.OnTokenTransfer(&_KeeperRegistrar.TransactOpts, sender, amount, data)
}

// Register is a paid mutator transaction binding the contract method 0x3659d666.
//
// Solidity: function register(string name, bytes encryptedEmail, address upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, uint96 amount, uint8 source, address sender) returns()
func (_KeeperRegistrar *KeeperRegistrarTransactor) Register(opts *bind.TransactOpts, name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, amount *big.Int, source uint8, sender common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "register", name, encryptedEmail, upkeepContract, gasLimit, adminAddress, checkData, amount, source, sender)
}

// Register is a paid mutator transaction binding the contract method 0x3659d666.
//
// Solidity: function register(string name, bytes encryptedEmail, address upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, uint96 amount, uint8 source, address sender) returns()
func (_KeeperRegistrar *KeeperRegistrarSession) Register(name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, amount *big.Int, source uint8, sender common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.Register(&_KeeperRegistrar.TransactOpts, name, encryptedEmail, upkeepContract, gasLimit, adminAddress, checkData, amount, source, sender)
}

// Register is a paid mutator transaction binding the contract method 0x3659d666.
//
// Solidity: function register(string name, bytes encryptedEmail, address upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, uint96 amount, uint8 source, address sender) returns()
func (_KeeperRegistrar *KeeperRegistrarTransactorSession) Register(name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, amount *big.Int, source uint8, sender common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.Register(&_KeeperRegistrar.TransactOpts, name, encryptedEmail, upkeepContract, gasLimit, adminAddress, checkData, amount, source, sender)
}

// SetAutoApproveAllowedSender is a paid mutator transaction binding the contract method 0x367b9b4f.
//
// Solidity: function setAutoApproveAllowedSender(address senderAddress, bool allowed) returns()
func (_KeeperRegistrar *KeeperRegistrarTransactor) SetAutoApproveAllowedSender(opts *bind.TransactOpts, senderAddress common.Address, allowed bool) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "setAutoApproveAllowedSender", senderAddress, allowed)
}

// SetAutoApproveAllowedSender is a paid mutator transaction binding the contract method 0x367b9b4f.
//
// Solidity: function setAutoApproveAllowedSender(address senderAddress, bool allowed) returns()
func (_KeeperRegistrar *KeeperRegistrarSession) SetAutoApproveAllowedSender(senderAddress common.Address, allowed bool) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.SetAutoApproveAllowedSender(&_KeeperRegistrar.TransactOpts, senderAddress, allowed)
}

// SetAutoApproveAllowedSender is a paid mutator transaction binding the contract method 0x367b9b4f.
//
// Solidity: function setAutoApproveAllowedSender(address senderAddress, bool allowed) returns()
func (_KeeperRegistrar *KeeperRegistrarTransactorSession) SetAutoApproveAllowedSender(senderAddress common.Address, allowed bool) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.SetAutoApproveAllowedSender(&_KeeperRegistrar.TransactOpts, senderAddress, allowed)
}

// SetRegistrationConfig is a paid mutator transaction binding the contract method 0xa793ab8b.
//
// Solidity: function setRegistrationConfig(uint8 autoApproveConfigType, uint16 autoApproveMaxAllowed, address keeperRegistry, uint96 minLINKJuels) returns()
func (_KeeperRegistrar *KeeperRegistrarTransactor) SetRegistrationConfig(opts *bind.TransactOpts, autoApproveConfigType uint8, autoApproveMaxAllowed uint16, keeperRegistry common.Address, minLINKJuels *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "setRegistrationConfig", autoApproveConfigType, autoApproveMaxAllowed, keeperRegistry, minLINKJuels)
}

// SetRegistrationConfig is a paid mutator transaction binding the contract method 0xa793ab8b.
//
// Solidity: function setRegistrationConfig(uint8 autoApproveConfigType, uint16 autoApproveMaxAllowed, address keeperRegistry, uint96 minLINKJuels) returns()
func (_KeeperRegistrar *KeeperRegistrarSession) SetRegistrationConfig(autoApproveConfigType uint8, autoApproveMaxAllowed uint16, keeperRegistry common.Address, minLINKJuels *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.SetRegistrationConfig(&_KeeperRegistrar.TransactOpts, autoApproveConfigType, autoApproveMaxAllowed, keeperRegistry, minLINKJuels)
}

// SetRegistrationConfig is a paid mutator transaction binding the contract method 0xa793ab8b.
//
// Solidity: function setRegistrationConfig(uint8 autoApproveConfigType, uint16 autoApproveMaxAllowed, address keeperRegistry, uint96 minLINKJuels) returns()
func (_KeeperRegistrar *KeeperRegistrarTransactorSession) SetRegistrationConfig(autoApproveConfigType uint8, autoApproveMaxAllowed uint16, keeperRegistry common.Address, minLINKJuels *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.SetRegistrationConfig(&_KeeperRegistrar.TransactOpts, autoApproveConfigType, autoApproveMaxAllowed, keeperRegistry, minLINKJuels)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_KeeperRegistrar *KeeperRegistrarTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_KeeperRegistrar *KeeperRegistrarSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.TransferOwnership(&_KeeperRegistrar.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_KeeperRegistrar *KeeperRegistrarTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.TransferOwnership(&_KeeperRegistrar.TransactOpts, to)
}

// KeeperRegistrarAutoApproveAllowedSenderSetIterator is returned from FilterAutoApproveAllowedSenderSet and is used to iterate over the raw logs and unpacked data for AutoApproveAllowedSenderSet events raised by the KeeperRegistrar contract.
type KeeperRegistrarAutoApproveAllowedSenderSetIterator struct {
	Event *KeeperRegistrarAutoApproveAllowedSenderSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistrarAutoApproveAllowedSenderSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrarAutoApproveAllowedSenderSet)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistrarAutoApproveAllowedSenderSet)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistrarAutoApproveAllowedSenderSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistrarAutoApproveAllowedSenderSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistrarAutoApproveAllowedSenderSet represents a AutoApproveAllowedSenderSet event raised by the KeeperRegistrar contract.
type KeeperRegistrarAutoApproveAllowedSenderSet struct {
	SenderAddress common.Address
	Allowed       bool
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAutoApproveAllowedSenderSet is a free log retrieval operation binding the contract event 0x20c6237dac83526a849285a9f79d08a483291bdd3a056a0ef9ae94ecee1ad356.
//
// Solidity: event AutoApproveAllowedSenderSet(address indexed senderAddress, bool allowed)
func (_KeeperRegistrar *KeeperRegistrarFilterer) FilterAutoApproveAllowedSenderSet(opts *bind.FilterOpts, senderAddress []common.Address) (*KeeperRegistrarAutoApproveAllowedSenderSetIterator, error) {

	var senderAddressRule []interface{}
	for _, senderAddressItem := range senderAddress {
		senderAddressRule = append(senderAddressRule, senderAddressItem)
	}

	logs, sub, err := _KeeperRegistrar.contract.FilterLogs(opts, "AutoApproveAllowedSenderSet", senderAddressRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarAutoApproveAllowedSenderSetIterator{contract: _KeeperRegistrar.contract, event: "AutoApproveAllowedSenderSet", logs: logs, sub: sub}, nil
}

// WatchAutoApproveAllowedSenderSet is a free log subscription operation binding the contract event 0x20c6237dac83526a849285a9f79d08a483291bdd3a056a0ef9ae94ecee1ad356.
//
// Solidity: event AutoApproveAllowedSenderSet(address indexed senderAddress, bool allowed)
func (_KeeperRegistrar *KeeperRegistrarFilterer) WatchAutoApproveAllowedSenderSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarAutoApproveAllowedSenderSet, senderAddress []common.Address) (event.Subscription, error) {

	var senderAddressRule []interface{}
	for _, senderAddressItem := range senderAddress {
		senderAddressRule = append(senderAddressRule, senderAddressItem)
	}

	logs, sub, err := _KeeperRegistrar.contract.WatchLogs(opts, "AutoApproveAllowedSenderSet", senderAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistrarAutoApproveAllowedSenderSet)
				if err := _KeeperRegistrar.contract.UnpackLog(event, "AutoApproveAllowedSenderSet", log); err != nil {
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

// ParseAutoApproveAllowedSenderSet is a log parse operation binding the contract event 0x20c6237dac83526a849285a9f79d08a483291bdd3a056a0ef9ae94ecee1ad356.
//
// Solidity: event AutoApproveAllowedSenderSet(address indexed senderAddress, bool allowed)
func (_KeeperRegistrar *KeeperRegistrarFilterer) ParseAutoApproveAllowedSenderSet(log types.Log) (*KeeperRegistrarAutoApproveAllowedSenderSet, error) {
	event := new(KeeperRegistrarAutoApproveAllowedSenderSet)
	if err := _KeeperRegistrar.contract.UnpackLog(event, "AutoApproveAllowedSenderSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistrarConfigChangedIterator is returned from FilterConfigChanged and is used to iterate over the raw logs and unpacked data for ConfigChanged events raised by the KeeperRegistrar contract.
type KeeperRegistrarConfigChangedIterator struct {
	Event *KeeperRegistrarConfigChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistrarConfigChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrarConfigChanged)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistrarConfigChanged)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistrarConfigChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistrarConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistrarConfigChanged represents a ConfigChanged event raised by the KeeperRegistrar contract.
type KeeperRegistrarConfigChanged struct {
	AutoApproveConfigType uint8
	AutoApproveMaxAllowed uint32
	KeeperRegistry        common.Address
	MinLINKJuels          *big.Int
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterConfigChanged is a free log retrieval operation binding the contract event 0x6293a703ec7145dfa23c5cde2e627d6a02e153fc2e9c03b14d1e22cbb4a7e9cd.
//
// Solidity: event ConfigChanged(uint8 autoApproveConfigType, uint32 autoApproveMaxAllowed, address keeperRegistry, uint96 minLINKJuels)
func (_KeeperRegistrar *KeeperRegistrarFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*KeeperRegistrarConfigChangedIterator, error) {

	logs, sub, err := _KeeperRegistrar.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarConfigChangedIterator{contract: _KeeperRegistrar.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

// WatchConfigChanged is a free log subscription operation binding the contract event 0x6293a703ec7145dfa23c5cde2e627d6a02e153fc2e9c03b14d1e22cbb4a7e9cd.
//
// Solidity: event ConfigChanged(uint8 autoApproveConfigType, uint32 autoApproveMaxAllowed, address keeperRegistry, uint96 minLINKJuels)
func (_KeeperRegistrar *KeeperRegistrarFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarConfigChanged) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistrar.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistrarConfigChanged)
				if err := _KeeperRegistrar.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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

// ParseConfigChanged is a log parse operation binding the contract event 0x6293a703ec7145dfa23c5cde2e627d6a02e153fc2e9c03b14d1e22cbb4a7e9cd.
//
// Solidity: event ConfigChanged(uint8 autoApproveConfigType, uint32 autoApproveMaxAllowed, address keeperRegistry, uint96 minLINKJuels)
func (_KeeperRegistrar *KeeperRegistrarFilterer) ParseConfigChanged(log types.Log) (*KeeperRegistrarConfigChanged, error) {
	event := new(KeeperRegistrarConfigChanged)
	if err := _KeeperRegistrar.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistrarOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the KeeperRegistrar contract.
type KeeperRegistrarOwnershipTransferRequestedIterator struct {
	Event *KeeperRegistrarOwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistrarOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrarOwnershipTransferRequested)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistrarOwnershipTransferRequested)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistrarOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistrarOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistrarOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the KeeperRegistrar contract.
type KeeperRegistrarOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_KeeperRegistrar *KeeperRegistrarFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistrarOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistrar.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarOwnershipTransferRequestedIterator{contract: _KeeperRegistrar.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_KeeperRegistrar *KeeperRegistrarFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistrar.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistrarOwnershipTransferRequested)
				if err := _KeeperRegistrar.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_KeeperRegistrar *KeeperRegistrarFilterer) ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistrarOwnershipTransferRequested, error) {
	event := new(KeeperRegistrarOwnershipTransferRequested)
	if err := _KeeperRegistrar.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistrarOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the KeeperRegistrar contract.
type KeeperRegistrarOwnershipTransferredIterator struct {
	Event *KeeperRegistrarOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistrarOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrarOwnershipTransferred)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistrarOwnershipTransferred)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistrarOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistrarOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistrarOwnershipTransferred represents a OwnershipTransferred event raised by the KeeperRegistrar contract.
type KeeperRegistrarOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_KeeperRegistrar *KeeperRegistrarFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistrarOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistrar.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarOwnershipTransferredIterator{contract: _KeeperRegistrar.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_KeeperRegistrar *KeeperRegistrarFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistrar.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistrarOwnershipTransferred)
				if err := _KeeperRegistrar.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_KeeperRegistrar *KeeperRegistrarFilterer) ParseOwnershipTransferred(log types.Log) (*KeeperRegistrarOwnershipTransferred, error) {
	event := new(KeeperRegistrarOwnershipTransferred)
	if err := _KeeperRegistrar.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistrarRegistrationApprovedIterator is returned from FilterRegistrationApproved and is used to iterate over the raw logs and unpacked data for RegistrationApproved events raised by the KeeperRegistrar contract.
type KeeperRegistrarRegistrationApprovedIterator struct {
	Event *KeeperRegistrarRegistrationApproved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistrarRegistrationApprovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrarRegistrationApproved)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistrarRegistrationApproved)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistrarRegistrationApprovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistrarRegistrationApprovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistrarRegistrationApproved represents a RegistrationApproved event raised by the KeeperRegistrar contract.
type KeeperRegistrarRegistrationApproved struct {
	Hash        [32]byte
	DisplayName string
	UpkeepId    *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRegistrationApproved is a free log retrieval operation binding the contract event 0xb9a292fb7e3edd920cd2d2829a3615a640c43fd7de0a0820aa0668feb4c37d4b.
//
// Solidity: event RegistrationApproved(bytes32 indexed hash, string displayName, uint256 indexed upkeepId)
func (_KeeperRegistrar *KeeperRegistrarFilterer) FilterRegistrationApproved(opts *bind.FilterOpts, hash [][32]byte, upkeepId []*big.Int) (*KeeperRegistrarRegistrationApprovedIterator, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	var upkeepIdRule []interface{}
	for _, upkeepIdItem := range upkeepId {
		upkeepIdRule = append(upkeepIdRule, upkeepIdItem)
	}

	logs, sub, err := _KeeperRegistrar.contract.FilterLogs(opts, "RegistrationApproved", hashRule, upkeepIdRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarRegistrationApprovedIterator{contract: _KeeperRegistrar.contract, event: "RegistrationApproved", logs: logs, sub: sub}, nil
}

// WatchRegistrationApproved is a free log subscription operation binding the contract event 0xb9a292fb7e3edd920cd2d2829a3615a640c43fd7de0a0820aa0668feb4c37d4b.
//
// Solidity: event RegistrationApproved(bytes32 indexed hash, string displayName, uint256 indexed upkeepId)
func (_KeeperRegistrar *KeeperRegistrarFilterer) WatchRegistrationApproved(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarRegistrationApproved, hash [][32]byte, upkeepId []*big.Int) (event.Subscription, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	var upkeepIdRule []interface{}
	for _, upkeepIdItem := range upkeepId {
		upkeepIdRule = append(upkeepIdRule, upkeepIdItem)
	}

	logs, sub, err := _KeeperRegistrar.contract.WatchLogs(opts, "RegistrationApproved", hashRule, upkeepIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistrarRegistrationApproved)
				if err := _KeeperRegistrar.contract.UnpackLog(event, "RegistrationApproved", log); err != nil {
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

// ParseRegistrationApproved is a log parse operation binding the contract event 0xb9a292fb7e3edd920cd2d2829a3615a640c43fd7de0a0820aa0668feb4c37d4b.
//
// Solidity: event RegistrationApproved(bytes32 indexed hash, string displayName, uint256 indexed upkeepId)
func (_KeeperRegistrar *KeeperRegistrarFilterer) ParseRegistrationApproved(log types.Log) (*KeeperRegistrarRegistrationApproved, error) {
	event := new(KeeperRegistrarRegistrationApproved)
	if err := _KeeperRegistrar.contract.UnpackLog(event, "RegistrationApproved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistrarRegistrationRejectedIterator is returned from FilterRegistrationRejected and is used to iterate over the raw logs and unpacked data for RegistrationRejected events raised by the KeeperRegistrar contract.
type KeeperRegistrarRegistrationRejectedIterator struct {
	Event *KeeperRegistrarRegistrationRejected // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistrarRegistrationRejectedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrarRegistrationRejected)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistrarRegistrationRejected)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistrarRegistrationRejectedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistrarRegistrationRejectedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistrarRegistrationRejected represents a RegistrationRejected event raised by the KeeperRegistrar contract.
type KeeperRegistrarRegistrationRejected struct {
	Hash [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterRegistrationRejected is a free log retrieval operation binding the contract event 0x3663fb28ebc87645eb972c9dad8521bf665c623f287e79f1c56f1eb374b82a22.
//
// Solidity: event RegistrationRejected(bytes32 indexed hash)
func (_KeeperRegistrar *KeeperRegistrarFilterer) FilterRegistrationRejected(opts *bind.FilterOpts, hash [][32]byte) (*KeeperRegistrarRegistrationRejectedIterator, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	logs, sub, err := _KeeperRegistrar.contract.FilterLogs(opts, "RegistrationRejected", hashRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarRegistrationRejectedIterator{contract: _KeeperRegistrar.contract, event: "RegistrationRejected", logs: logs, sub: sub}, nil
}

// WatchRegistrationRejected is a free log subscription operation binding the contract event 0x3663fb28ebc87645eb972c9dad8521bf665c623f287e79f1c56f1eb374b82a22.
//
// Solidity: event RegistrationRejected(bytes32 indexed hash)
func (_KeeperRegistrar *KeeperRegistrarFilterer) WatchRegistrationRejected(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarRegistrationRejected, hash [][32]byte) (event.Subscription, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	logs, sub, err := _KeeperRegistrar.contract.WatchLogs(opts, "RegistrationRejected", hashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistrarRegistrationRejected)
				if err := _KeeperRegistrar.contract.UnpackLog(event, "RegistrationRejected", log); err != nil {
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

// ParseRegistrationRejected is a log parse operation binding the contract event 0x3663fb28ebc87645eb972c9dad8521bf665c623f287e79f1c56f1eb374b82a22.
//
// Solidity: event RegistrationRejected(bytes32 indexed hash)
func (_KeeperRegistrar *KeeperRegistrarFilterer) ParseRegistrationRejected(log types.Log) (*KeeperRegistrarRegistrationRejected, error) {
	event := new(KeeperRegistrarRegistrationRejected)
	if err := _KeeperRegistrar.contract.UnpackLog(event, "RegistrationRejected", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistrarRegistrationRequestedIterator is returned from FilterRegistrationRequested and is used to iterate over the raw logs and unpacked data for RegistrationRequested events raised by the KeeperRegistrar contract.
type KeeperRegistrarRegistrationRequestedIterator struct {
	Event *KeeperRegistrarRegistrationRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeeperRegistrarRegistrationRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrarRegistrationRequested)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeeperRegistrarRegistrationRequested)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeeperRegistrarRegistrationRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistrarRegistrationRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistrarRegistrationRequested represents a RegistrationRequested event raised by the KeeperRegistrar contract.
type KeeperRegistrarRegistrationRequested struct {
	Hash           [32]byte
	Name           string
	EncryptedEmail []byte
	UpkeepContract common.Address
	GasLimit       uint32
	AdminAddress   common.Address
	CheckData      []byte
	Amount         *big.Int
	Source         uint8
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterRegistrationRequested is a free log retrieval operation binding the contract event 0xc3f5df4aefec026f610a3fcb08f19476492d69d2cb78b1c2eba259a8820e6a78.
//
// Solidity: event RegistrationRequested(bytes32 indexed hash, string name, bytes encryptedEmail, address indexed upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, uint96 amount, uint8 indexed source)
func (_KeeperRegistrar *KeeperRegistrarFilterer) FilterRegistrationRequested(opts *bind.FilterOpts, hash [][32]byte, upkeepContract []common.Address, source []uint8) (*KeeperRegistrarRegistrationRequestedIterator, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	var upkeepContractRule []interface{}
	for _, upkeepContractItem := range upkeepContract {
		upkeepContractRule = append(upkeepContractRule, upkeepContractItem)
	}

	var sourceRule []interface{}
	for _, sourceItem := range source {
		sourceRule = append(sourceRule, sourceItem)
	}

	logs, sub, err := _KeeperRegistrar.contract.FilterLogs(opts, "RegistrationRequested", hashRule, upkeepContractRule, sourceRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarRegistrationRequestedIterator{contract: _KeeperRegistrar.contract, event: "RegistrationRequested", logs: logs, sub: sub}, nil
}

// WatchRegistrationRequested is a free log subscription operation binding the contract event 0xc3f5df4aefec026f610a3fcb08f19476492d69d2cb78b1c2eba259a8820e6a78.
//
// Solidity: event RegistrationRequested(bytes32 indexed hash, string name, bytes encryptedEmail, address indexed upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, uint96 amount, uint8 indexed source)
func (_KeeperRegistrar *KeeperRegistrarFilterer) WatchRegistrationRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarRegistrationRequested, hash [][32]byte, upkeepContract []common.Address, source []uint8) (event.Subscription, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	var upkeepContractRule []interface{}
	for _, upkeepContractItem := range upkeepContract {
		upkeepContractRule = append(upkeepContractRule, upkeepContractItem)
	}

	var sourceRule []interface{}
	for _, sourceItem := range source {
		sourceRule = append(sourceRule, sourceItem)
	}

	logs, sub, err := _KeeperRegistrar.contract.WatchLogs(opts, "RegistrationRequested", hashRule, upkeepContractRule, sourceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistrarRegistrationRequested)
				if err := _KeeperRegistrar.contract.UnpackLog(event, "RegistrationRequested", log); err != nil {
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

// ParseRegistrationRequested is a log parse operation binding the contract event 0xc3f5df4aefec026f610a3fcb08f19476492d69d2cb78b1c2eba259a8820e6a78.
//
// Solidity: event RegistrationRequested(bytes32 indexed hash, string name, bytes encryptedEmail, address indexed upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, uint96 amount, uint8 indexed source)
func (_KeeperRegistrar *KeeperRegistrarFilterer) ParseRegistrationRequested(log types.Log) (*KeeperRegistrarRegistrationRequested, error) {
	event := new(KeeperRegistrarRegistrationRequested)
	if err := _KeeperRegistrar.contract.UnpackLog(event, "RegistrationRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
