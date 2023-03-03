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

// KeeperRegistrar20RegistrationParams is an auto generated low-level Go binding around an user-defined struct.
type KeeperRegistrar20RegistrationParams struct {
	Name           string
	EncryptedEmail []byte
	UpkeepContract common.Address
	GasLimit       uint32
	AdminAddress   common.Address
	CheckData      []byte
	OffchainConfig []byte
	Amount         *big.Int
}

// KeeperRegistrar20MetaData contains all meta data concerning the KeeperRegistrar20 contract.
var KeeperRegistrar20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"LINKAddress\",\"type\":\"address\"},{\"internalType\":\"enumKeeperRegistrar2_0.AutoApproveType\",\"name\":\"autoApproveConfigType\",\"type\":\"uint8\"},{\"internalType\":\"uint16\",\"name\":\"autoApproveMaxAllowed\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"minLINKJuels\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AmountMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FunctionNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HashMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientPayment\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAdminAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"LinkTransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyAdminOrOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistrationRequestFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RequestNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SenderMismatch\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"AutoApproveAllowedSenderSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"enumKeeperRegistrar2_0.AutoApproveType\",\"name\":\"autoApproveConfigType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"autoApproveMaxAllowed\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"minLINKJuels\",\"type\":\"uint96\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"displayName\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"RegistrationApproved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"RegistrationRejected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encryptedEmail\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"upkeepContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"RegistrationRequested\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"upkeepContract\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"cancel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"}],\"name\":\"getAutoApproveAllowedSender\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"getPendingRequest\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRegistrationConfig\",\"outputs\":[{\"internalType\":\"enumKeeperRegistrar2_0.AutoApproveType\",\"name\":\"autoApproveConfigType\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"autoApproveMaxAllowed\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"approvedCount\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minLINKJuels\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"encryptedEmail\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"upkeepContract\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"encryptedEmail\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"upkeepContract\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"internalType\":\"structKeeperRegistrar2_0.RegistrationParams\",\"name\":\"requestParams\",\"type\":\"tuple\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"setAutoApproveAllowedSender\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumKeeperRegistrar2_0.AutoApproveType\",\"name\":\"autoApproveConfigType\",\"type\":\"uint8\"},{\"internalType\":\"uint16\",\"name\":\"autoApproveMaxAllowed\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"minLINKJuels\",\"type\":\"uint96\"}],\"name\":\"setRegistrationConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b506040516200238b3803806200238b833981016040819052620000349162000394565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000ec565b5050506001600160601b0319606086901b16608052620000e18484848462000198565b50505050506200048d565b6001600160a01b038116331415620001475760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b620001a262000319565b6003546040805160a081019091526501000000000090910463ffffffff169080866002811115620001d757620001d762000477565b815261ffff8616602082015263ffffffff831660408201526001600160a01b03851660608201526001600160601b038416608090910152805160038054909190829060ff1916600183600281111562000234576200023462000477565b0217905550602082015181546040808501516060860151610100600160481b031990931661010063ffffffff9586160263ffffffff60281b19161765010000000000949091169390930292909217600160481b600160e81b03191669010000000000000000006001600160a01b0390921691909102178255608090920151600190910180546001600160601b0319166001600160601b03909216919091179055517f6293a703ec7145dfa23c5cde2e627d6a02e153fc2e9c03b14d1e22cbb4a7e9cd906200030a90879087908790879062000422565b60405180910390a15050505050565b6000546001600160a01b03163314620003755760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000082565b565b80516001600160a01b03811681146200038f57600080fd5b919050565b600080600080600060a08688031215620003ad57600080fd5b620003b88662000377565b9450602086015160038110620003cd57600080fd5b604087015190945061ffff81168114620003e657600080fd5b9250620003f66060870162000377565b60808701519092506001600160601b03811681146200041457600080fd5b809150509295509295909350565b60808101600386106200044557634e487b7160e01b600052602160045260246000fd5b94815261ffff9390931660208401526001600160a01b039190911660408301526001600160601b031660609091015290565b634e487b7160e01b600052602160045260246000fd5b60805160601c611ebc620004cf600039600081816101360152818161030b01528181610781015281816109c001528181610d4101526111840152611ebc6000f3fe608060405234801561001057600080fd5b50600436106100c55760003560e01c806308b79da4146100ca578063181f5a77146100f05780631b6b6d2314610131578063367b9b4f14610165578063621058541461017a57806379ba50971461018d5780637e776f7f14610195578063850af0cb146101d157806388b12d55146101ea5780638da5cb5b14610249578063a4c0ed361461025a578063a611ea561461026d578063a793ab8b14610280578063c4d252f514610293578063f2fde38b146102a6575b600080fd5b6100dd6100d836600461194b565b6102b9565b6040519081526020015b60405180910390f35b6101246040518060400160405280601581526020017404b656570657252656769737472617220322e302e3605c1b81525081565b6040516100e79190611c06565b6101587f000000000000000000000000000000000000000000000000000000000000000081565b6040516100e79190611a62565b610178610173366004611521565b6103ec565b005b610178610188366004611650565b610453565b610178610600565b6101c16101a33660046114fd565b6001600160a01b031660009081526005602052604090205460ff1690565b60405190151581526020016100e7565b6101d96106af565b6040516100e7959493929190611bc3565b61023b6101f83660046115d2565b6000908152600260209081526040918290208251808401909352546001600160a01b038116808452600160a01b9091046001600160601b03169290910182905291565b6040516100e7929190611b26565b6000546001600160a01b0316610158565b61017861026836600461155a565b610776565b61017861027b366004611726565b6109b5565b61017861028e3660046115eb565b610b0f565b6101786102a13660046115d2565b610c7d565b6101786102b43660046114fd565b610e21565b6004546000906001600160601b03166102d9610100840160e0850161199f565b6001600160601b031610156103015760405163cd1c886760e01b815260040160405180910390fd5b6001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000166323b872dd3330610343610100870160e0880161199f565b6040516001600160e01b031960e086901b1681526001600160a01b0393841660048201529290911660248301526001600160601b03166044820152606401602060405180830381600087803b15801561039b57600080fd5b505af11580156103af573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103d391906115b5565b506103e66103e083611d26565b33610e35565b92915050565b6103f461107d565b6001600160a01b038216600081815260056020908152604091829020805460ff191685151590811790915591519182527f20c6237dac83526a849285a9f79d08a483291bdd3a056a0ef9ae94ecee1ad356910160405180910390a25050565b61045b61107d565b6000818152600260209081526040918290208251808401909352546001600160a01b038116808452600160a01b9091046001600160601b0316918301919091526104b857604051632589d98f60e11b815260040160405180910390fd5b6000898989898989896040516020016104d79796959493929190611a76565b60405160208183030381529060405280519060200120905080831461050f57604051633f4d605360e01b815260040160405180910390fd5b60008381526002602090815260408083208390558051610100810182528e815281518084018352938452808301939093526001600160a01b038d81168483015263ffffffff8d1660608501528b1660808401528051601f8a018390048302810183019091528881526105f2929160a0830191908b908b9081908401838280828437600092019190915250505090825250604080516020601f8a01819004810282018101909252888152918101919089908990819084018382808284376000920191909152505050908252506020858101516001600160601b0316910152826110d2565b505050505050505050505050565b6001546001600160a01b031633146106585760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b60448201526064015b60405180910390fd5b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6040805160a081019091526003805460009283928392839283928392829060ff1660028111156106e1576106e1611e37565b60028111156106f2576106f2611e37565b81528154610100810463ffffffff908116602080850191909152600160281b8304909116604080850191909152600160481b9092046001600160a01b03166060808501919091526001909401546001600160601b0390811660809485015285519186015192860151948601519590930151909b919a50929850929650169350915050565b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146107be5760405162c6885f60e11b815260040160405180910390fd5b81818080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050505060208101516001600160e01b03198116635308f52b60e11b1461082a5760405163e3d6792160e01b815260040160405180910390fd5b848484600061083c8260048186611cb2565b810190610849919061184b565b50975050505050505050806001600160601b0316841461087c576040516355e97b0d60e01b815260040160405180910390fd5b898888600061088e8260048186611cb2565b81019061089b919061184b565b98505050505050505050806001600160a01b0316846001600160a01b0316146108d757604051637c62b1c760e11b815260040160405180910390fd5b6101248b10156108fa57604051630dfe930960e41b815260040160405180910390fd5b6004546001600160601b03168d10156109265760405163cd1c886760e01b815260040160405180910390fd5b6000306001600160a01b03168d8d604051610942929190611a52565b600060405180830381855af49150503d806000811461097d576040519150601f19603f3d011682016040523d82523d6000602084013e610982565b606091505b50509050806109a457604051630649bf8160e41b815260040160405180910390fd5b505050505050505050505050505050565b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146109fd5760405162c6885f60e11b815260040160405180910390fd5b610b006040518061010001604052808e81526020018d8d8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152505050908252506001600160a01b03808d1660208084019190915263ffffffff8d16604080850191909152918c1660608401528151601f8b018290048202810182019092528982526080909201918a908a9081908401838280828437600092019190915250505090825250604080516020601f8901819004810282018101909252878152918101919088908890819084018382808284376000920191909152505050908252506001600160601b03851660209091015282610e35565b50505050505050505050505050565b610b1761107d565b6003546040805160a08101909152600160281b90910463ffffffff169080866002811115610b4757610b47611e37565b815261ffff8616602082015263ffffffff831660408201526001600160a01b03851660608201526001600160601b038416608090910152805160038054909190829060ff19166001836002811115610ba157610ba1611e37565b021790555060208201518154604080850151606086015168ffffffffffffffff001990931661010063ffffffff9586160263ffffffff60281b191617600160281b949091169390930292909217600160481b600160e81b031916600160481b6001600160a01b0390921691909102178255608090920151600190910180546001600160601b0319166001600160601b03909216919091179055517f6293a703ec7145dfa23c5cde2e627d6a02e153fc2e9c03b14d1e22cbb4a7e9cd90610c6e908790879087908790611b84565b60405180910390a15050505050565b6000818152600260209081526040918290208251808401909352546001600160a01b038116808452600160a01b9091046001600160601b031691830191909152331480610cd457506000546001600160a01b031633145b610cf1576040516361685c2b60e01b815260040160405180910390fd5b80516001600160a01b0316610d1957604051632589d98f60e11b815260040160405180910390fd5b6000828152600260209081526040808320839055835191840151905163a9059cbb60e01b81527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169263a9059cbb92610d7c92600401611b26565b602060405180830381600087803b158015610d9657600080fd5b505af1158015610daa573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610dce91906115b5565b905080610df157815160405163185c9b9d60e31b815261064f9190600401611a62565b60405183907f3663fb28ebc87645eb972c9dad8521bf665c623f287e79f1c56f1eb374b82a2290600090a2505050565b610e2961107d565b610e32816112b1565b50565b60808201516000906001600160a01b0316610e635760405163016ed19f60e21b815260040160405180910390fd5b60008360400151846060015185608001518660a001518760c00151604051602001610e92959493929190611ad0565b60405160208183030381529060405280519060200120905083604001516001600160a01b0316817f9b8456f925542af2c5fb15ff4be32cc8f209dda96c544766e301367df40f499886600001518760200151886060015189608001518a60a001518b60e00151604051610f0a96959493929190611c19565b60405180910390a36040805160a081019091526003805460009283929091829060ff166002811115610f3e57610f3e611e37565b6002811115610f4f57610f4f611e37565b8152815463ffffffff610100820481166020840152600160281b82041660408301526001600160a01b03600160481b9091041660608201526001909101546001600160601b03166080909101529050610fa88186611355565b15610ff3576040810151610fbd906001611cdc565b6003805463ffffffff92909216600160281b0263ffffffff60281b19909216919091179055610fec86846110d2565b9150611074565b60e0860151600084815260026020526040812054909161102291600160a01b90046001600160601b0316611d04565b60408051808201825260808a01516001600160a01b0390811682526001600160601b03938416602080840191825260008a815260029091529390932091519251909316600160a01b0291909216179055505b50949350505050565b6000546001600160a01b031633146110d05760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b604482015260640161064f565b565b6003546040838101516060850151608086015160a087015160c088015194516336f6cf5760e11b8152600096600160481b90046001600160a01b03169587958795636ded9eae9561112c9592949193909291600401611ad0565b602060405180830381600087803b15801561114657600080fd5b505af115801561115a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061117e9190611986565b905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316634000aea0848860e00151856040516020016111ca91815260200190565b6040516020818303038152906040526040518463ffffffff1660e01b81526004016111f793929190611b48565b602060405180830381600087803b15801561121157600080fd5b505af1158015611225573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061124991906115b5565b90508061126b578260405163185c9b9d60e31b815260040161064f9190611a62565b81857fb9a292fb7e3edd920cd2d2829a3615a640c43fd7de0a0820aa0668feb4c37d4b88600001516040516112a09190611c06565b60405180910390a350949350505050565b6001600160a01b0381163314156113045760405162461bcd60e51b815260206004820152601760248201527621b0b73737ba103a3930b739b332b9103a379039b2b63360491b604482015260640161064f565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000808351600281111561136b5761136b611e37565b1415611379575060006103e6565b60018351600281111561138e5761138e611e37565b1480156113b457506001600160a01b03821660009081526005602052604090205460ff16155b156113c1575060006103e6565b826020015163ffffffff16836040015163ffffffff1610156113e5575060016103e6565b50600092915050565b80356113f981611e63565b919050565b60008083601f84011261141057600080fd5b5081356001600160401b0381111561142757600080fd5b60208301915083602082850101111561143f57600080fd5b9250929050565b600082601f83011261145757600080fd5b81356001600160401b038082111561147157611471611e4d565b604051601f8301601f19908116603f0116810190828211818310171561149957611499611e4d565b816040528381528660208588010111156114b257600080fd5b836020870160208301376000602085830101528094505050505092915050565b803563ffffffff811681146113f957600080fd5b80356001600160601b03811681146113f957600080fd5b60006020828403121561150f57600080fd5b813561151a81611e63565b9392505050565b6000806040838503121561153457600080fd5b823561153f81611e63565b9150602083013561154f81611e78565b809150509250929050565b6000806000806060858703121561157057600080fd5b843561157b81611e63565b93506020850135925060408501356001600160401b0381111561159d57600080fd5b6115a9878288016113fe565b95989497509550505050565b6000602082840312156115c757600080fd5b815161151a81611e78565b6000602082840312156115e457600080fd5b5035919050565b6000806000806080858703121561160157600080fd5b84356003811061161057600080fd5b9350602085013561ffff8116811461162757600080fd5b9250604085013561163781611e63565b9150611645606086016114e6565b905092959194509250565b600080600080600080600080600060e08a8c03121561166e57600080fd5b89356001600160401b038082111561168557600080fd5b6116918d838e01611446565b9a5060208c013591506116a382611e63565b8199506116b260408d016114d2565b985060608c013591506116c482611e63565b90965060808b013590808211156116da57600080fd5b6116e68d838e016113fe565b909750955060a08c01359150808211156116ff57600080fd5b5061170c8c828d016113fe565b9a9d999c50979a9699959894979660c00135949350505050565b6000806000806000806000806000806000806101208d8f03121561174957600080fd5b6001600160401b038d35111561175e57600080fd5b61176b8e8e358f01611446565b9b506001600160401b0360208e0135111561178557600080fd5b6117958e60208f01358f016113fe565b909b5099506117a660408e016113ee565b98506117b460608e016114d2565b97506117c260808e016113ee565b96506001600160401b0360a08e013511156117dc57600080fd5b6117ec8e60a08f01358f016113fe565b90965094506001600160401b0360c08e0135111561180957600080fd5b6118198e60c08f01358f016113fe565b909450925061182a60e08e016114e6565b91506118396101008e016113ee565b90509295989b509295989b509295989b565b60008060008060008060008060006101208a8c03121561186a57600080fd5b89356001600160401b038082111561188157600080fd5b61188d8d838e01611446565b9a5060208c01359150808211156118a357600080fd5b6118af8d838e01611446565b99506118bd60408d016113ee565b98506118cb60608d016114d2565b97506118d960808d016113ee565b965060a08c01359150808211156118ef57600080fd5b6118fb8d838e01611446565b955060c08c013591508082111561191157600080fd5b5061191e8c828d01611446565b93505061192d60e08b016114e6565b915061193c6101008b016113ee565b90509295985092959850929598565b60006020828403121561195d57600080fd5b81356001600160401b0381111561197357600080fd5b8201610100818503121561151a57600080fd5b60006020828403121561199857600080fd5b5051919050565b6000602082840312156119b157600080fd5b61151a826114e6565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b6000815180845260005b81811015611a09576020818501810151868301820152016119ed565b81811115611a1b576000602083870101525b50601f01601f19169290920160200192915050565b60038110611a4e57634e487b7160e01b600052602160045260246000fd5b9052565b8183823760009101908152919050565b6001600160a01b0391909116815260200190565b6001600160a01b03888116825263ffffffff881660208301528616604082015260a060608201819052600090611aaf90830186886119ba565b8281036080840152611ac28185876119ba565b9a9950505050505050505050565b6001600160a01b03868116825263ffffffff861660208301528416604082015260a060608201819052600090611b08908301856119e3565b8281036080840152611b1a81856119e3565b98975050505050505050565b6001600160a01b039290921682526001600160601b0316602082015260400190565b6001600160a01b03841681526001600160601b0383166020820152606060408201819052600090611b7b908301846119e3565b95945050505050565b60808101611b928287611a30565b61ffff9490941660208201526001600160a01b039290921660408301526001600160601b0316606090910152919050565b60a08101611bd18288611a30565b63ffffffff95861660208301529390941660408501526001600160a01b03919091166060840152608090920191909152919050565b60208152600061151a60208301846119e3565b60c081526000611c2c60c08301896119e3565b8281036020840152611c3e81896119e3565b63ffffffff881660408501526001600160a01b038716606085015283810360808501529050611c6d81866119e3565b91505060018060601b03831660a0830152979650505050505050565b60405161010081016001600160401b0381118282101715611cac57611cac611e4d565b60405290565b60008085851115611cc257600080fd5b83861115611ccf57600080fd5b5050820193919092039150565b600063ffffffff808316818516808303821115611cfb57611cfb611e21565b01949350505050565b60006001600160601b03828116848216808303821115611cfb57611cfb611e21565b60006101008236031215611d3957600080fd5b611d41611c89565b82356001600160401b0380821115611d5857600080fd5b611d6436838701611446565b83526020850135915080821115611d7a57600080fd5b611d8636838701611446565b6020840152611d97604086016113ee565b6040840152611da8606086016114d2565b6060840152611db9608086016113ee565b608084015260a0850135915080821115611dd257600080fd5b611dde36838701611446565b60a084015260c0850135915080821115611df757600080fd5b50611e0436828601611446565b60c083015250611e1660e084016114e6565b60e082015292915050565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052602160045260246000fd5b634e487b7160e01b600052604160045260246000fd5b6001600160a01b0381168114610e3257600080fd5b8015158114610e3257600080fdfea26469706673582212209a04ff19e14c3c50e929bc0ce9127176d73bdbd904805d60285edd5a83970eab64736f6c63430008060033",
}

// KeeperRegistrar20ABI is the input ABI used to generate the binding from.
// Deprecated: Use KeeperRegistrar20MetaData.ABI instead.
var KeeperRegistrar20ABI = KeeperRegistrar20MetaData.ABI

// KeeperRegistrar20Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use KeeperRegistrar20MetaData.Bin instead.
var KeeperRegistrar20Bin = KeeperRegistrar20MetaData.Bin

// DeployKeeperRegistrar20 deploys a new Ethereum contract, binding an instance of KeeperRegistrar20 to it.
func DeployKeeperRegistrar20(auth *bind.TransactOpts, backend bind.ContractBackend, LINKAddress common.Address, autoApproveConfigType uint8, autoApproveMaxAllowed uint16, keeperRegistry common.Address, minLINKJuels *big.Int) (common.Address, *types.Transaction, *KeeperRegistrar20, error) {
	parsed, err := KeeperRegistrar20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistrar20Bin), backend, LINKAddress, autoApproveConfigType, autoApproveMaxAllowed, keeperRegistry, minLINKJuels)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperRegistrar20{KeeperRegistrar20Caller: KeeperRegistrar20Caller{contract: contract}, KeeperRegistrar20Transactor: KeeperRegistrar20Transactor{contract: contract}, KeeperRegistrar20Filterer: KeeperRegistrar20Filterer{contract: contract}}, nil
}

// KeeperRegistrar20 is an auto generated Go binding around an Ethereum contract.
type KeeperRegistrar20 struct {
	KeeperRegistrar20Caller     // Read-only binding to the contract
	KeeperRegistrar20Transactor // Write-only binding to the contract
	KeeperRegistrar20Filterer   // Log filterer for contract events
}

// KeeperRegistrar20Caller is an auto generated read-only Go binding around an Ethereum contract.
type KeeperRegistrar20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperRegistrar20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type KeeperRegistrar20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperRegistrar20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KeeperRegistrar20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeeperRegistrar20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KeeperRegistrar20Session struct {
	Contract     *KeeperRegistrar20 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// KeeperRegistrar20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KeeperRegistrar20CallerSession struct {
	Contract *KeeperRegistrar20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// KeeperRegistrar20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KeeperRegistrar20TransactorSession struct {
	Contract     *KeeperRegistrar20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// KeeperRegistrar20Raw is an auto generated low-level Go binding around an Ethereum contract.
type KeeperRegistrar20Raw struct {
	Contract *KeeperRegistrar20 // Generic contract binding to access the raw methods on
}

// KeeperRegistrar20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KeeperRegistrar20CallerRaw struct {
	Contract *KeeperRegistrar20Caller // Generic read-only contract binding to access the raw methods on
}

// KeeperRegistrar20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KeeperRegistrar20TransactorRaw struct {
	Contract *KeeperRegistrar20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewKeeperRegistrar20 creates a new instance of KeeperRegistrar20, bound to a specific deployed contract.
func NewKeeperRegistrar20(address common.Address, backend bind.ContractBackend) (*KeeperRegistrar20, error) {
	contract, err := bindKeeperRegistrar20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrar20{KeeperRegistrar20Caller: KeeperRegistrar20Caller{contract: contract}, KeeperRegistrar20Transactor: KeeperRegistrar20Transactor{contract: contract}, KeeperRegistrar20Filterer: KeeperRegistrar20Filterer{contract: contract}}, nil
}

// NewKeeperRegistrar20Caller creates a new read-only instance of KeeperRegistrar20, bound to a specific deployed contract.
func NewKeeperRegistrar20Caller(address common.Address, caller bind.ContractCaller) (*KeeperRegistrar20Caller, error) {
	contract, err := bindKeeperRegistrar20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrar20Caller{contract: contract}, nil
}

// NewKeeperRegistrar20Transactor creates a new write-only instance of KeeperRegistrar20, bound to a specific deployed contract.
func NewKeeperRegistrar20Transactor(address common.Address, transactor bind.ContractTransactor) (*KeeperRegistrar20Transactor, error) {
	contract, err := bindKeeperRegistrar20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrar20Transactor{contract: contract}, nil
}

// NewKeeperRegistrar20Filterer creates a new log filterer instance of KeeperRegistrar20, bound to a specific deployed contract.
func NewKeeperRegistrar20Filterer(address common.Address, filterer bind.ContractFilterer) (*KeeperRegistrar20Filterer, error) {
	contract, err := bindKeeperRegistrar20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrar20Filterer{contract: contract}, nil
}

// bindKeeperRegistrar20 binds a generic wrapper to an already deployed contract.
func bindKeeperRegistrar20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KeeperRegistrar20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeeperRegistrar20 *KeeperRegistrar20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistrar20.Contract.KeeperRegistrar20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeeperRegistrar20 *KeeperRegistrar20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistrar20.Contract.KeeperRegistrar20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeeperRegistrar20 *KeeperRegistrar20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistrar20.Contract.KeeperRegistrar20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeeperRegistrar20 *KeeperRegistrar20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistrar20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeeperRegistrar20 *KeeperRegistrar20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistrar20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeeperRegistrar20 *KeeperRegistrar20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistrar20.Contract.contract.Transact(opts, method, params...)
}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_KeeperRegistrar20 *KeeperRegistrar20Caller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistrar20.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_KeeperRegistrar20 *KeeperRegistrar20Session) LINK() (common.Address, error) {
	return _KeeperRegistrar20.Contract.LINK(&_KeeperRegistrar20.CallOpts)
}

// LINK is a free data retrieval call binding the contract method 0x1b6b6d23.
//
// Solidity: function LINK() view returns(address)
func (_KeeperRegistrar20 *KeeperRegistrar20CallerSession) LINK() (common.Address, error) {
	return _KeeperRegistrar20.Contract.LINK(&_KeeperRegistrar20.CallOpts)
}

// GetAutoApproveAllowedSender is a free data retrieval call binding the contract method 0x7e776f7f.
//
// Solidity: function getAutoApproveAllowedSender(address senderAddress) view returns(bool)
func (_KeeperRegistrar20 *KeeperRegistrar20Caller) GetAutoApproveAllowedSender(opts *bind.CallOpts, senderAddress common.Address) (bool, error) {
	var out []interface{}
	err := _KeeperRegistrar20.contract.Call(opts, &out, "getAutoApproveAllowedSender", senderAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetAutoApproveAllowedSender is a free data retrieval call binding the contract method 0x7e776f7f.
//
// Solidity: function getAutoApproveAllowedSender(address senderAddress) view returns(bool)
func (_KeeperRegistrar20 *KeeperRegistrar20Session) GetAutoApproveAllowedSender(senderAddress common.Address) (bool, error) {
	return _KeeperRegistrar20.Contract.GetAutoApproveAllowedSender(&_KeeperRegistrar20.CallOpts, senderAddress)
}

// GetAutoApproveAllowedSender is a free data retrieval call binding the contract method 0x7e776f7f.
//
// Solidity: function getAutoApproveAllowedSender(address senderAddress) view returns(bool)
func (_KeeperRegistrar20 *KeeperRegistrar20CallerSession) GetAutoApproveAllowedSender(senderAddress common.Address) (bool, error) {
	return _KeeperRegistrar20.Contract.GetAutoApproveAllowedSender(&_KeeperRegistrar20.CallOpts, senderAddress)
}

// GetPendingRequest is a free data retrieval call binding the contract method 0x88b12d55.
//
// Solidity: function getPendingRequest(bytes32 hash) view returns(address, uint96)
func (_KeeperRegistrar20 *KeeperRegistrar20Caller) GetPendingRequest(opts *bind.CallOpts, hash [32]byte) (common.Address, *big.Int, error) {
	var out []interface{}
	err := _KeeperRegistrar20.contract.Call(opts, &out, "getPendingRequest", hash)

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
func (_KeeperRegistrar20 *KeeperRegistrar20Session) GetPendingRequest(hash [32]byte) (common.Address, *big.Int, error) {
	return _KeeperRegistrar20.Contract.GetPendingRequest(&_KeeperRegistrar20.CallOpts, hash)
}

// GetPendingRequest is a free data retrieval call binding the contract method 0x88b12d55.
//
// Solidity: function getPendingRequest(bytes32 hash) view returns(address, uint96)
func (_KeeperRegistrar20 *KeeperRegistrar20CallerSession) GetPendingRequest(hash [32]byte) (common.Address, *big.Int, error) {
	return _KeeperRegistrar20.Contract.GetPendingRequest(&_KeeperRegistrar20.CallOpts, hash)
}

// GetRegistrationConfig is a free data retrieval call binding the contract method 0x850af0cb.
//
// Solidity: function getRegistrationConfig() view returns(uint8 autoApproveConfigType, uint32 autoApproveMaxAllowed, uint32 approvedCount, address keeperRegistry, uint256 minLINKJuels)
func (_KeeperRegistrar20 *KeeperRegistrar20Caller) GetRegistrationConfig(opts *bind.CallOpts) (struct {
	AutoApproveConfigType uint8
	AutoApproveMaxAllowed uint32
	ApprovedCount         uint32
	KeeperRegistry        common.Address
	MinLINKJuels          *big.Int
}, error) {
	var out []interface{}
	err := _KeeperRegistrar20.contract.Call(opts, &out, "getRegistrationConfig")

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
func (_KeeperRegistrar20 *KeeperRegistrar20Session) GetRegistrationConfig() (struct {
	AutoApproveConfigType uint8
	AutoApproveMaxAllowed uint32
	ApprovedCount         uint32
	KeeperRegistry        common.Address
	MinLINKJuels          *big.Int
}, error) {
	return _KeeperRegistrar20.Contract.GetRegistrationConfig(&_KeeperRegistrar20.CallOpts)
}

// GetRegistrationConfig is a free data retrieval call binding the contract method 0x850af0cb.
//
// Solidity: function getRegistrationConfig() view returns(uint8 autoApproveConfigType, uint32 autoApproveMaxAllowed, uint32 approvedCount, address keeperRegistry, uint256 minLINKJuels)
func (_KeeperRegistrar20 *KeeperRegistrar20CallerSession) GetRegistrationConfig() (struct {
	AutoApproveConfigType uint8
	AutoApproveMaxAllowed uint32
	ApprovedCount         uint32
	KeeperRegistry        common.Address
	MinLINKJuels          *big.Int
}, error) {
	return _KeeperRegistrar20.Contract.GetRegistrationConfig(&_KeeperRegistrar20.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KeeperRegistrar20 *KeeperRegistrar20Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistrar20.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KeeperRegistrar20 *KeeperRegistrar20Session) Owner() (common.Address, error) {
	return _KeeperRegistrar20.Contract.Owner(&_KeeperRegistrar20.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_KeeperRegistrar20 *KeeperRegistrar20CallerSession) Owner() (common.Address, error) {
	return _KeeperRegistrar20.Contract.Owner(&_KeeperRegistrar20.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_KeeperRegistrar20 *KeeperRegistrar20Caller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _KeeperRegistrar20.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_KeeperRegistrar20 *KeeperRegistrar20Session) TypeAndVersion() (string, error) {
	return _KeeperRegistrar20.Contract.TypeAndVersion(&_KeeperRegistrar20.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_KeeperRegistrar20 *KeeperRegistrar20CallerSession) TypeAndVersion() (string, error) {
	return _KeeperRegistrar20.Contract.TypeAndVersion(&_KeeperRegistrar20.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_KeeperRegistrar20 *KeeperRegistrar20Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistrar20.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_KeeperRegistrar20 *KeeperRegistrar20Session) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistrar20.Contract.AcceptOwnership(&_KeeperRegistrar20.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_KeeperRegistrar20 *KeeperRegistrar20TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistrar20.Contract.AcceptOwnership(&_KeeperRegistrar20.TransactOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x62105854.
//
// Solidity: function approve(string name, address upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, bytes offchainConfig, bytes32 hash) returns()
func (_KeeperRegistrar20 *KeeperRegistrar20Transactor) Approve(opts *bind.TransactOpts, name string, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, offchainConfig []byte, hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrar20.contract.Transact(opts, "approve", name, upkeepContract, gasLimit, adminAddress, checkData, offchainConfig, hash)
}

// Approve is a paid mutator transaction binding the contract method 0x62105854.
//
// Solidity: function approve(string name, address upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, bytes offchainConfig, bytes32 hash) returns()
func (_KeeperRegistrar20 *KeeperRegistrar20Session) Approve(name string, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, offchainConfig []byte, hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrar20.Contract.Approve(&_KeeperRegistrar20.TransactOpts, name, upkeepContract, gasLimit, adminAddress, checkData, offchainConfig, hash)
}

// Approve is a paid mutator transaction binding the contract method 0x62105854.
//
// Solidity: function approve(string name, address upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, bytes offchainConfig, bytes32 hash) returns()
func (_KeeperRegistrar20 *KeeperRegistrar20TransactorSession) Approve(name string, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, offchainConfig []byte, hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrar20.Contract.Approve(&_KeeperRegistrar20.TransactOpts, name, upkeepContract, gasLimit, adminAddress, checkData, offchainConfig, hash)
}

// Cancel is a paid mutator transaction binding the contract method 0xc4d252f5.
//
// Solidity: function cancel(bytes32 hash) returns()
func (_KeeperRegistrar20 *KeeperRegistrar20Transactor) Cancel(opts *bind.TransactOpts, hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrar20.contract.Transact(opts, "cancel", hash)
}

// Cancel is a paid mutator transaction binding the contract method 0xc4d252f5.
//
// Solidity: function cancel(bytes32 hash) returns()
func (_KeeperRegistrar20 *KeeperRegistrar20Session) Cancel(hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrar20.Contract.Cancel(&_KeeperRegistrar20.TransactOpts, hash)
}

// Cancel is a paid mutator transaction binding the contract method 0xc4d252f5.
//
// Solidity: function cancel(bytes32 hash) returns()
func (_KeeperRegistrar20 *KeeperRegistrar20TransactorSession) Cancel(hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrar20.Contract.Cancel(&_KeeperRegistrar20.TransactOpts, hash)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_KeeperRegistrar20 *KeeperRegistrar20Transactor) OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistrar20.contract.Transact(opts, "onTokenTransfer", sender, amount, data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_KeeperRegistrar20 *KeeperRegistrar20Session) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistrar20.Contract.OnTokenTransfer(&_KeeperRegistrar20.TransactOpts, sender, amount, data)
}

// OnTokenTransfer is a paid mutator transaction binding the contract method 0xa4c0ed36.
//
// Solidity: function onTokenTransfer(address sender, uint256 amount, bytes data) returns()
func (_KeeperRegistrar20 *KeeperRegistrar20TransactorSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistrar20.Contract.OnTokenTransfer(&_KeeperRegistrar20.TransactOpts, sender, amount, data)
}

// Register is a paid mutator transaction binding the contract method 0xa611ea56.
//
// Solidity: function register(string name, bytes encryptedEmail, address upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, bytes offchainConfig, uint96 amount, address sender) returns()
func (_KeeperRegistrar20 *KeeperRegistrar20Transactor) Register(opts *bind.TransactOpts, name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, offchainConfig []byte, amount *big.Int, sender common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar20.contract.Transact(opts, "register", name, encryptedEmail, upkeepContract, gasLimit, adminAddress, checkData, offchainConfig, amount, sender)
}

// Register is a paid mutator transaction binding the contract method 0xa611ea56.
//
// Solidity: function register(string name, bytes encryptedEmail, address upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, bytes offchainConfig, uint96 amount, address sender) returns()
func (_KeeperRegistrar20 *KeeperRegistrar20Session) Register(name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, offchainConfig []byte, amount *big.Int, sender common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar20.Contract.Register(&_KeeperRegistrar20.TransactOpts, name, encryptedEmail, upkeepContract, gasLimit, adminAddress, checkData, offchainConfig, amount, sender)
}

// Register is a paid mutator transaction binding the contract method 0xa611ea56.
//
// Solidity: function register(string name, bytes encryptedEmail, address upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, bytes offchainConfig, uint96 amount, address sender) returns()
func (_KeeperRegistrar20 *KeeperRegistrar20TransactorSession) Register(name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, offchainConfig []byte, amount *big.Int, sender common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar20.Contract.Register(&_KeeperRegistrar20.TransactOpts, name, encryptedEmail, upkeepContract, gasLimit, adminAddress, checkData, offchainConfig, amount, sender)
}

// RegisterUpkeep is a paid mutator transaction binding the contract method 0x08b79da4.
//
// Solidity: function registerUpkeep((string,bytes,address,uint32,address,bytes,bytes,uint96) requestParams) returns(uint256)
func (_KeeperRegistrar20 *KeeperRegistrar20Transactor) RegisterUpkeep(opts *bind.TransactOpts, requestParams KeeperRegistrar20RegistrationParams) (*types.Transaction, error) {
	return _KeeperRegistrar20.contract.Transact(opts, "registerUpkeep", requestParams)
}

// RegisterUpkeep is a paid mutator transaction binding the contract method 0x08b79da4.
//
// Solidity: function registerUpkeep((string,bytes,address,uint32,address,bytes,bytes,uint96) requestParams) returns(uint256)
func (_KeeperRegistrar20 *KeeperRegistrar20Session) RegisterUpkeep(requestParams KeeperRegistrar20RegistrationParams) (*types.Transaction, error) {
	return _KeeperRegistrar20.Contract.RegisterUpkeep(&_KeeperRegistrar20.TransactOpts, requestParams)
}

// RegisterUpkeep is a paid mutator transaction binding the contract method 0x08b79da4.
//
// Solidity: function registerUpkeep((string,bytes,address,uint32,address,bytes,bytes,uint96) requestParams) returns(uint256)
func (_KeeperRegistrar20 *KeeperRegistrar20TransactorSession) RegisterUpkeep(requestParams KeeperRegistrar20RegistrationParams) (*types.Transaction, error) {
	return _KeeperRegistrar20.Contract.RegisterUpkeep(&_KeeperRegistrar20.TransactOpts, requestParams)
}

// SetAutoApproveAllowedSender is a paid mutator transaction binding the contract method 0x367b9b4f.
//
// Solidity: function setAutoApproveAllowedSender(address senderAddress, bool allowed) returns()
func (_KeeperRegistrar20 *KeeperRegistrar20Transactor) SetAutoApproveAllowedSender(opts *bind.TransactOpts, senderAddress common.Address, allowed bool) (*types.Transaction, error) {
	return _KeeperRegistrar20.contract.Transact(opts, "setAutoApproveAllowedSender", senderAddress, allowed)
}

// SetAutoApproveAllowedSender is a paid mutator transaction binding the contract method 0x367b9b4f.
//
// Solidity: function setAutoApproveAllowedSender(address senderAddress, bool allowed) returns()
func (_KeeperRegistrar20 *KeeperRegistrar20Session) SetAutoApproveAllowedSender(senderAddress common.Address, allowed bool) (*types.Transaction, error) {
	return _KeeperRegistrar20.Contract.SetAutoApproveAllowedSender(&_KeeperRegistrar20.TransactOpts, senderAddress, allowed)
}

// SetAutoApproveAllowedSender is a paid mutator transaction binding the contract method 0x367b9b4f.
//
// Solidity: function setAutoApproveAllowedSender(address senderAddress, bool allowed) returns()
func (_KeeperRegistrar20 *KeeperRegistrar20TransactorSession) SetAutoApproveAllowedSender(senderAddress common.Address, allowed bool) (*types.Transaction, error) {
	return _KeeperRegistrar20.Contract.SetAutoApproveAllowedSender(&_KeeperRegistrar20.TransactOpts, senderAddress, allowed)
}

// SetRegistrationConfig is a paid mutator transaction binding the contract method 0xa793ab8b.
//
// Solidity: function setRegistrationConfig(uint8 autoApproveConfigType, uint16 autoApproveMaxAllowed, address keeperRegistry, uint96 minLINKJuels) returns()
func (_KeeperRegistrar20 *KeeperRegistrar20Transactor) SetRegistrationConfig(opts *bind.TransactOpts, autoApproveConfigType uint8, autoApproveMaxAllowed uint16, keeperRegistry common.Address, minLINKJuels *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrar20.contract.Transact(opts, "setRegistrationConfig", autoApproveConfigType, autoApproveMaxAllowed, keeperRegistry, minLINKJuels)
}

// SetRegistrationConfig is a paid mutator transaction binding the contract method 0xa793ab8b.
//
// Solidity: function setRegistrationConfig(uint8 autoApproveConfigType, uint16 autoApproveMaxAllowed, address keeperRegistry, uint96 minLINKJuels) returns()
func (_KeeperRegistrar20 *KeeperRegistrar20Session) SetRegistrationConfig(autoApproveConfigType uint8, autoApproveMaxAllowed uint16, keeperRegistry common.Address, minLINKJuels *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrar20.Contract.SetRegistrationConfig(&_KeeperRegistrar20.TransactOpts, autoApproveConfigType, autoApproveMaxAllowed, keeperRegistry, minLINKJuels)
}

// SetRegistrationConfig is a paid mutator transaction binding the contract method 0xa793ab8b.
//
// Solidity: function setRegistrationConfig(uint8 autoApproveConfigType, uint16 autoApproveMaxAllowed, address keeperRegistry, uint96 minLINKJuels) returns()
func (_KeeperRegistrar20 *KeeperRegistrar20TransactorSession) SetRegistrationConfig(autoApproveConfigType uint8, autoApproveMaxAllowed uint16, keeperRegistry common.Address, minLINKJuels *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrar20.Contract.SetRegistrationConfig(&_KeeperRegistrar20.TransactOpts, autoApproveConfigType, autoApproveMaxAllowed, keeperRegistry, minLINKJuels)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_KeeperRegistrar20 *KeeperRegistrar20Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar20.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_KeeperRegistrar20 *KeeperRegistrar20Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar20.Contract.TransferOwnership(&_KeeperRegistrar20.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_KeeperRegistrar20 *KeeperRegistrar20TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar20.Contract.TransferOwnership(&_KeeperRegistrar20.TransactOpts, to)
}

// KeeperRegistrar20AutoApproveAllowedSenderSetIterator is returned from FilterAutoApproveAllowedSenderSet and is used to iterate over the raw logs and unpacked data for AutoApproveAllowedSenderSet events raised by the KeeperRegistrar20 contract.
type KeeperRegistrar20AutoApproveAllowedSenderSetIterator struct {
	Event *KeeperRegistrar20AutoApproveAllowedSenderSet // Event containing the contract specifics and raw log

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
func (it *KeeperRegistrar20AutoApproveAllowedSenderSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrar20AutoApproveAllowedSenderSet)
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
		it.Event = new(KeeperRegistrar20AutoApproveAllowedSenderSet)
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
func (it *KeeperRegistrar20AutoApproveAllowedSenderSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistrar20AutoApproveAllowedSenderSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistrar20AutoApproveAllowedSenderSet represents a AutoApproveAllowedSenderSet event raised by the KeeperRegistrar20 contract.
type KeeperRegistrar20AutoApproveAllowedSenderSet struct {
	SenderAddress common.Address
	Allowed       bool
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAutoApproveAllowedSenderSet is a free log retrieval operation binding the contract event 0x20c6237dac83526a849285a9f79d08a483291bdd3a056a0ef9ae94ecee1ad356.
//
// Solidity: event AutoApproveAllowedSenderSet(address indexed senderAddress, bool allowed)
func (_KeeperRegistrar20 *KeeperRegistrar20Filterer) FilterAutoApproveAllowedSenderSet(opts *bind.FilterOpts, senderAddress []common.Address) (*KeeperRegistrar20AutoApproveAllowedSenderSetIterator, error) {

	var senderAddressRule []interface{}
	for _, senderAddressItem := range senderAddress {
		senderAddressRule = append(senderAddressRule, senderAddressItem)
	}

	logs, sub, err := _KeeperRegistrar20.contract.FilterLogs(opts, "AutoApproveAllowedSenderSet", senderAddressRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrar20AutoApproveAllowedSenderSetIterator{contract: _KeeperRegistrar20.contract, event: "AutoApproveAllowedSenderSet", logs: logs, sub: sub}, nil
}

// WatchAutoApproveAllowedSenderSet is a free log subscription operation binding the contract event 0x20c6237dac83526a849285a9f79d08a483291bdd3a056a0ef9ae94ecee1ad356.
//
// Solidity: event AutoApproveAllowedSenderSet(address indexed senderAddress, bool allowed)
func (_KeeperRegistrar20 *KeeperRegistrar20Filterer) WatchAutoApproveAllowedSenderSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistrar20AutoApproveAllowedSenderSet, senderAddress []common.Address) (event.Subscription, error) {

	var senderAddressRule []interface{}
	for _, senderAddressItem := range senderAddress {
		senderAddressRule = append(senderAddressRule, senderAddressItem)
	}

	logs, sub, err := _KeeperRegistrar20.contract.WatchLogs(opts, "AutoApproveAllowedSenderSet", senderAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistrar20AutoApproveAllowedSenderSet)
				if err := _KeeperRegistrar20.contract.UnpackLog(event, "AutoApproveAllowedSenderSet", log); err != nil {
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
func (_KeeperRegistrar20 *KeeperRegistrar20Filterer) ParseAutoApproveAllowedSenderSet(log types.Log) (*KeeperRegistrar20AutoApproveAllowedSenderSet, error) {
	event := new(KeeperRegistrar20AutoApproveAllowedSenderSet)
	if err := _KeeperRegistrar20.contract.UnpackLog(event, "AutoApproveAllowedSenderSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistrar20ConfigChangedIterator is returned from FilterConfigChanged and is used to iterate over the raw logs and unpacked data for ConfigChanged events raised by the KeeperRegistrar20 contract.
type KeeperRegistrar20ConfigChangedIterator struct {
	Event *KeeperRegistrar20ConfigChanged // Event containing the contract specifics and raw log

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
func (it *KeeperRegistrar20ConfigChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrar20ConfigChanged)
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
		it.Event = new(KeeperRegistrar20ConfigChanged)
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
func (it *KeeperRegistrar20ConfigChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistrar20ConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistrar20ConfigChanged represents a ConfigChanged event raised by the KeeperRegistrar20 contract.
type KeeperRegistrar20ConfigChanged struct {
	AutoApproveConfigType uint8
	AutoApproveMaxAllowed uint32
	KeeperRegistry        common.Address
	MinLINKJuels          *big.Int
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterConfigChanged is a free log retrieval operation binding the contract event 0x6293a703ec7145dfa23c5cde2e627d6a02e153fc2e9c03b14d1e22cbb4a7e9cd.
//
// Solidity: event ConfigChanged(uint8 autoApproveConfigType, uint32 autoApproveMaxAllowed, address keeperRegistry, uint96 minLINKJuels)
func (_KeeperRegistrar20 *KeeperRegistrar20Filterer) FilterConfigChanged(opts *bind.FilterOpts) (*KeeperRegistrar20ConfigChangedIterator, error) {

	logs, sub, err := _KeeperRegistrar20.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrar20ConfigChangedIterator{contract: _KeeperRegistrar20.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

// WatchConfigChanged is a free log subscription operation binding the contract event 0x6293a703ec7145dfa23c5cde2e627d6a02e153fc2e9c03b14d1e22cbb4a7e9cd.
//
// Solidity: event ConfigChanged(uint8 autoApproveConfigType, uint32 autoApproveMaxAllowed, address keeperRegistry, uint96 minLINKJuels)
func (_KeeperRegistrar20 *KeeperRegistrar20Filterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *KeeperRegistrar20ConfigChanged) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistrar20.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistrar20ConfigChanged)
				if err := _KeeperRegistrar20.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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
func (_KeeperRegistrar20 *KeeperRegistrar20Filterer) ParseConfigChanged(log types.Log) (*KeeperRegistrar20ConfigChanged, error) {
	event := new(KeeperRegistrar20ConfigChanged)
	if err := _KeeperRegistrar20.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistrar20OwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the KeeperRegistrar20 contract.
type KeeperRegistrar20OwnershipTransferRequestedIterator struct {
	Event *KeeperRegistrar20OwnershipTransferRequested // Event containing the contract specifics and raw log

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
func (it *KeeperRegistrar20OwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrar20OwnershipTransferRequested)
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
		it.Event = new(KeeperRegistrar20OwnershipTransferRequested)
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
func (it *KeeperRegistrar20OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistrar20OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistrar20OwnershipTransferRequested represents a OwnershipTransferRequested event raised by the KeeperRegistrar20 contract.
type KeeperRegistrar20OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_KeeperRegistrar20 *KeeperRegistrar20Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistrar20OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistrar20.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrar20OwnershipTransferRequestedIterator{contract: _KeeperRegistrar20.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_KeeperRegistrar20 *KeeperRegistrar20Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistrar20OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistrar20.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistrar20OwnershipTransferRequested)
				if err := _KeeperRegistrar20.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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
func (_KeeperRegistrar20 *KeeperRegistrar20Filterer) ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistrar20OwnershipTransferRequested, error) {
	event := new(KeeperRegistrar20OwnershipTransferRequested)
	if err := _KeeperRegistrar20.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistrar20OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the KeeperRegistrar20 contract.
type KeeperRegistrar20OwnershipTransferredIterator struct {
	Event *KeeperRegistrar20OwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *KeeperRegistrar20OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrar20OwnershipTransferred)
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
		it.Event = new(KeeperRegistrar20OwnershipTransferred)
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
func (it *KeeperRegistrar20OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistrar20OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistrar20OwnershipTransferred represents a OwnershipTransferred event raised by the KeeperRegistrar20 contract.
type KeeperRegistrar20OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_KeeperRegistrar20 *KeeperRegistrar20Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistrar20OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistrar20.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrar20OwnershipTransferredIterator{contract: _KeeperRegistrar20.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_KeeperRegistrar20 *KeeperRegistrar20Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistrar20OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistrar20.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistrar20OwnershipTransferred)
				if err := _KeeperRegistrar20.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_KeeperRegistrar20 *KeeperRegistrar20Filterer) ParseOwnershipTransferred(log types.Log) (*KeeperRegistrar20OwnershipTransferred, error) {
	event := new(KeeperRegistrar20OwnershipTransferred)
	if err := _KeeperRegistrar20.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistrar20RegistrationApprovedIterator is returned from FilterRegistrationApproved and is used to iterate over the raw logs and unpacked data for RegistrationApproved events raised by the KeeperRegistrar20 contract.
type KeeperRegistrar20RegistrationApprovedIterator struct {
	Event *KeeperRegistrar20RegistrationApproved // Event containing the contract specifics and raw log

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
func (it *KeeperRegistrar20RegistrationApprovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrar20RegistrationApproved)
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
		it.Event = new(KeeperRegistrar20RegistrationApproved)
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
func (it *KeeperRegistrar20RegistrationApprovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistrar20RegistrationApprovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistrar20RegistrationApproved represents a RegistrationApproved event raised by the KeeperRegistrar20 contract.
type KeeperRegistrar20RegistrationApproved struct {
	Hash        [32]byte
	DisplayName string
	UpkeepId    *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRegistrationApproved is a free log retrieval operation binding the contract event 0xb9a292fb7e3edd920cd2d2829a3615a640c43fd7de0a0820aa0668feb4c37d4b.
//
// Solidity: event RegistrationApproved(bytes32 indexed hash, string displayName, uint256 indexed upkeepId)
func (_KeeperRegistrar20 *KeeperRegistrar20Filterer) FilterRegistrationApproved(opts *bind.FilterOpts, hash [][32]byte, upkeepId []*big.Int) (*KeeperRegistrar20RegistrationApprovedIterator, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	var upkeepIdRule []interface{}
	for _, upkeepIdItem := range upkeepId {
		upkeepIdRule = append(upkeepIdRule, upkeepIdItem)
	}

	logs, sub, err := _KeeperRegistrar20.contract.FilterLogs(opts, "RegistrationApproved", hashRule, upkeepIdRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrar20RegistrationApprovedIterator{contract: _KeeperRegistrar20.contract, event: "RegistrationApproved", logs: logs, sub: sub}, nil
}

// WatchRegistrationApproved is a free log subscription operation binding the contract event 0xb9a292fb7e3edd920cd2d2829a3615a640c43fd7de0a0820aa0668feb4c37d4b.
//
// Solidity: event RegistrationApproved(bytes32 indexed hash, string displayName, uint256 indexed upkeepId)
func (_KeeperRegistrar20 *KeeperRegistrar20Filterer) WatchRegistrationApproved(opts *bind.WatchOpts, sink chan<- *KeeperRegistrar20RegistrationApproved, hash [][32]byte, upkeepId []*big.Int) (event.Subscription, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	var upkeepIdRule []interface{}
	for _, upkeepIdItem := range upkeepId {
		upkeepIdRule = append(upkeepIdRule, upkeepIdItem)
	}

	logs, sub, err := _KeeperRegistrar20.contract.WatchLogs(opts, "RegistrationApproved", hashRule, upkeepIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistrar20RegistrationApproved)
				if err := _KeeperRegistrar20.contract.UnpackLog(event, "RegistrationApproved", log); err != nil {
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
func (_KeeperRegistrar20 *KeeperRegistrar20Filterer) ParseRegistrationApproved(log types.Log) (*KeeperRegistrar20RegistrationApproved, error) {
	event := new(KeeperRegistrar20RegistrationApproved)
	if err := _KeeperRegistrar20.contract.UnpackLog(event, "RegistrationApproved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistrar20RegistrationRejectedIterator is returned from FilterRegistrationRejected and is used to iterate over the raw logs and unpacked data for RegistrationRejected events raised by the KeeperRegistrar20 contract.
type KeeperRegistrar20RegistrationRejectedIterator struct {
	Event *KeeperRegistrar20RegistrationRejected // Event containing the contract specifics and raw log

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
func (it *KeeperRegistrar20RegistrationRejectedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrar20RegistrationRejected)
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
		it.Event = new(KeeperRegistrar20RegistrationRejected)
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
func (it *KeeperRegistrar20RegistrationRejectedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistrar20RegistrationRejectedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistrar20RegistrationRejected represents a RegistrationRejected event raised by the KeeperRegistrar20 contract.
type KeeperRegistrar20RegistrationRejected struct {
	Hash [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterRegistrationRejected is a free log retrieval operation binding the contract event 0x3663fb28ebc87645eb972c9dad8521bf665c623f287e79f1c56f1eb374b82a22.
//
// Solidity: event RegistrationRejected(bytes32 indexed hash)
func (_KeeperRegistrar20 *KeeperRegistrar20Filterer) FilterRegistrationRejected(opts *bind.FilterOpts, hash [][32]byte) (*KeeperRegistrar20RegistrationRejectedIterator, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	logs, sub, err := _KeeperRegistrar20.contract.FilterLogs(opts, "RegistrationRejected", hashRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrar20RegistrationRejectedIterator{contract: _KeeperRegistrar20.contract, event: "RegistrationRejected", logs: logs, sub: sub}, nil
}

// WatchRegistrationRejected is a free log subscription operation binding the contract event 0x3663fb28ebc87645eb972c9dad8521bf665c623f287e79f1c56f1eb374b82a22.
//
// Solidity: event RegistrationRejected(bytes32 indexed hash)
func (_KeeperRegistrar20 *KeeperRegistrar20Filterer) WatchRegistrationRejected(opts *bind.WatchOpts, sink chan<- *KeeperRegistrar20RegistrationRejected, hash [][32]byte) (event.Subscription, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	logs, sub, err := _KeeperRegistrar20.contract.WatchLogs(opts, "RegistrationRejected", hashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistrar20RegistrationRejected)
				if err := _KeeperRegistrar20.contract.UnpackLog(event, "RegistrationRejected", log); err != nil {
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
func (_KeeperRegistrar20 *KeeperRegistrar20Filterer) ParseRegistrationRejected(log types.Log) (*KeeperRegistrar20RegistrationRejected, error) {
	event := new(KeeperRegistrar20RegistrationRejected)
	if err := _KeeperRegistrar20.contract.UnpackLog(event, "RegistrationRejected", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KeeperRegistrar20RegistrationRequestedIterator is returned from FilterRegistrationRequested and is used to iterate over the raw logs and unpacked data for RegistrationRequested events raised by the KeeperRegistrar20 contract.
type KeeperRegistrar20RegistrationRequestedIterator struct {
	Event *KeeperRegistrar20RegistrationRequested // Event containing the contract specifics and raw log

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
func (it *KeeperRegistrar20RegistrationRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistrar20RegistrationRequested)
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
		it.Event = new(KeeperRegistrar20RegistrationRequested)
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
func (it *KeeperRegistrar20RegistrationRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeeperRegistrar20RegistrationRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeeperRegistrar20RegistrationRequested represents a RegistrationRequested event raised by the KeeperRegistrar20 contract.
type KeeperRegistrar20RegistrationRequested struct {
	Hash           [32]byte
	Name           string
	EncryptedEmail []byte
	UpkeepContract common.Address
	GasLimit       uint32
	AdminAddress   common.Address
	CheckData      []byte
	Amount         *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterRegistrationRequested is a free log retrieval operation binding the contract event 0x9b8456f925542af2c5fb15ff4be32cc8f209dda96c544766e301367df40f4998.
//
// Solidity: event RegistrationRequested(bytes32 indexed hash, string name, bytes encryptedEmail, address indexed upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, uint96 amount)
func (_KeeperRegistrar20 *KeeperRegistrar20Filterer) FilterRegistrationRequested(opts *bind.FilterOpts, hash [][32]byte, upkeepContract []common.Address) (*KeeperRegistrar20RegistrationRequestedIterator, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	var upkeepContractRule []interface{}
	for _, upkeepContractItem := range upkeepContract {
		upkeepContractRule = append(upkeepContractRule, upkeepContractItem)
	}

	logs, sub, err := _KeeperRegistrar20.contract.FilterLogs(opts, "RegistrationRequested", hashRule, upkeepContractRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrar20RegistrationRequestedIterator{contract: _KeeperRegistrar20.contract, event: "RegistrationRequested", logs: logs, sub: sub}, nil
}

// WatchRegistrationRequested is a free log subscription operation binding the contract event 0x9b8456f925542af2c5fb15ff4be32cc8f209dda96c544766e301367df40f4998.
//
// Solidity: event RegistrationRequested(bytes32 indexed hash, string name, bytes encryptedEmail, address indexed upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, uint96 amount)
func (_KeeperRegistrar20 *KeeperRegistrar20Filterer) WatchRegistrationRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistrar20RegistrationRequested, hash [][32]byte, upkeepContract []common.Address) (event.Subscription, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	var upkeepContractRule []interface{}
	for _, upkeepContractItem := range upkeepContract {
		upkeepContractRule = append(upkeepContractRule, upkeepContractItem)
	}

	logs, sub, err := _KeeperRegistrar20.contract.WatchLogs(opts, "RegistrationRequested", hashRule, upkeepContractRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeeperRegistrar20RegistrationRequested)
				if err := _KeeperRegistrar20.contract.UnpackLog(event, "RegistrationRequested", log); err != nil {
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

// ParseRegistrationRequested is a log parse operation binding the contract event 0x9b8456f925542af2c5fb15ff4be32cc8f209dda96c544766e301367df40f4998.
//
// Solidity: event RegistrationRequested(bytes32 indexed hash, string name, bytes encryptedEmail, address indexed upkeepContract, uint32 gasLimit, address adminAddress, bytes checkData, uint96 amount)
func (_KeeperRegistrar20 *KeeperRegistrar20Filterer) ParseRegistrationRequested(log types.Log) (*KeeperRegistrar20RegistrationRequested, error) {
	event := new(KeeperRegistrar20RegistrationRequested)
	if err := _KeeperRegistrar20.contract.UnpackLog(event, "RegistrationRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
