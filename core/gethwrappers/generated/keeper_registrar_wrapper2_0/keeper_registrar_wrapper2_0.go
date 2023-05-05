// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package keeper_registrar_wrapper2_0

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

var KeeperRegistrarMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"LINKAddress\",\"type\":\"address\"},{\"internalType\":\"enumKeeperRegistrar2_0.AutoApproveType\",\"name\":\"autoApproveConfigType\",\"type\":\"uint8\"},{\"internalType\":\"uint16\",\"name\":\"autoApproveMaxAllowed\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"minLINKJuels\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AmountMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FunctionNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HashMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientPayment\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAdminAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"LinkTransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyAdminOrOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistrationRequestFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RequestNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SenderMismatch\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"AutoApproveAllowedSenderSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"enumKeeperRegistrar2_0.AutoApproveType\",\"name\":\"autoApproveConfigType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"autoApproveMaxAllowed\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"minLINKJuels\",\"type\":\"uint96\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"displayName\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"RegistrationApproved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"RegistrationRejected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encryptedEmail\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"upkeepContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"RegistrationRequested\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"upkeepContract\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"cancel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"}],\"name\":\"getAutoApproveAllowedSender\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"getPendingRequest\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRegistrationConfig\",\"outputs\":[{\"internalType\":\"enumKeeperRegistrar2_0.AutoApproveType\",\"name\":\"autoApproveConfigType\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"autoApproveMaxAllowed\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"approvedCount\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minLINKJuels\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"encryptedEmail\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"upkeepContract\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"encryptedEmail\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"upkeepContract\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"internalType\":\"structKeeperRegistrar2_0.RegistrationParams\",\"name\":\"requestParams\",\"type\":\"tuple\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"setAutoApproveAllowedSender\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumKeeperRegistrar2_0.AutoApproveType\",\"name\":\"autoApproveConfigType\",\"type\":\"uint8\"},{\"internalType\":\"uint16\",\"name\":\"autoApproveMaxAllowed\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"minLINKJuels\",\"type\":\"uint96\"}],\"name\":\"setRegistrationConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162002b0538038062002b05833981016040819052620000349162000394565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000ec565b5050506001600160601b0319606086901b16608052620000e18484848462000198565b50505050506200048d565b6001600160a01b038116331415620001475760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b620001a262000319565b6003546040805160a081019091526501000000000090910463ffffffff169080866002811115620001d757620001d762000477565b815261ffff8616602082015263ffffffff831660408201526001600160a01b03851660608201526001600160601b038416608090910152805160038054909190829060ff1916600183600281111562000234576200023462000477565b0217905550602082015181546040808501516060860151610100600160481b031990931661010063ffffffff9586160263ffffffff60281b19161765010000000000949091169390930292909217600160481b600160e81b03191669010000000000000000006001600160a01b0390921691909102178255608090920151600190910180546001600160601b0319166001600160601b03909216919091179055517f6293a703ec7145dfa23c5cde2e627d6a02e153fc2e9c03b14d1e22cbb4a7e9cd906200030a90879087908790879062000422565b60405180910390a15050505050565b6000546001600160a01b03163314620003755760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000082565b565b80516001600160a01b03811681146200038f57600080fd5b919050565b600080600080600060a08688031215620003ad57600080fd5b620003b88662000377565b9450602086015160038110620003cd57600080fd5b604087015190945061ffff81168114620003e657600080fd5b9250620003f66060870162000377565b60808701519092506001600160601b03811681146200041457600080fd5b809150509295509295909350565b60808101600386106200044557634e487b7160e01b600052602160045260246000fd5b94815261ffff9390931660208401526001600160a01b039190911660408301526001600160601b031660609091015290565b634e487b7160e01b600052602160045260246000fd5b60805160601c612636620004cf6000396000818161016e015281816103f30152818161099f01528181610cfd015281816111d9015261178301526126366000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c8063850af0cb11610097578063a611ea5611610066578063a611ea5614610325578063a793ab8b14610338578063c4d252f51461034b578063f2fde38b1461035e57600080fd5b8063850af0cb1461022e57806388b12d55146102475780638da5cb5b146102f4578063a4c0ed361461031257600080fd5b8063367b9b4f116100d3578063367b9b4f146101b557806362105854146101ca57806379ba5097146101dd5780637e776f7f146101e557600080fd5b806308b79da4146100fa578063181f5a77146101205780631b6b6d2314610169575b600080fd5b61010d610108366004612011565b610371565b6040519081526020015b60405180910390f35b61015c6040518060400160405280601581526020017f4b656570657252656769737472617220322e302e30000000000000000000000081525081565b6040516101179190612337565b6101907f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610117565b6101c86101c3366004611be0565b6104fe565b005b6101c86101d8366004611d10565b610590565b6101c86107a4565b61021e6101f3366004611bbc565b73ffffffffffffffffffffffffffffffffffffffff1660009081526005602052604090205460ff1690565b6040519015158152602001610117565b6102366108a6565b6040516101179594939291906122e7565b6102bb610255366004611c92565b60009081526002602090815260409182902082518084019093525473ffffffffffffffffffffffffffffffffffffffff8116808452740100000000000000000000000000000000000000009091046bffffffffffffffffffffffff169290910182905291565b6040805173ffffffffffffffffffffffffffffffffffffffff90931683526bffffffffffffffffffffffff909116602083015201610117565b60005473ffffffffffffffffffffffffffffffffffffffff16610190565b6101c8610320366004611c19565b610987565b6101c8610333366004611de7565b610ce5565b6101c8610346366004611cab565b610e78565b6101c8610359366004611c92565b61108d565b6101c861036c366004611bbc565b611326565b6004546000906bffffffffffffffffffffffff16610396610100840160e08501612066565b6bffffffffffffffffffffffff1610156103dc576040517fcd1c886700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166323b872dd333061042b610100870160e08801612066565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e086901b16815273ffffffffffffffffffffffffffffffffffffffff93841660048201529290911660248301526bffffffffffffffffffffffff166044820152606401602060405180830381600087803b1580156104ad57600080fd5b505af11580156104c1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104e59190611c75565b506104f86104f283612470565b3361133a565b92915050565b610506611622565b73ffffffffffffffffffffffffffffffffffffffff821660008181526005602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001685151590811790915591519182527f20c6237dac83526a849285a9f79d08a483291bdd3a056a0ef9ae94ecee1ad356910160405180910390a25050565b610598611622565b60008181526002602090815260409182902082518084019093525473ffffffffffffffffffffffffffffffffffffffff8116808452740100000000000000000000000000000000000000009091046bffffffffffffffffffffffff1691830191909152610631576040517f4b13b31e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000898989898989896040516020016106509796959493929190612180565b6040516020818303038152906040528051906020012090508083146106a1576040517f3f4d605300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008381526002602090815260408083208390558051610100810182528e8152815180840183529384528083019390935273ffffffffffffffffffffffffffffffffffffffff8d81168483015263ffffffff8d1660608501528b1660808401528051601f8a01839004830281018301909152888152610796929160a0830191908b908b9081908401838280828437600092019190915250505090825250604080516020601f8a01819004810282018101909252888152918101919089908990819084018382808284376000920191909152505050908252506020858101516bffffffffffffffffffffffff16910152826116a5565b505050505050505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461082a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6040805160a081019091526003805460009283928392839283928392829060ff1660028111156108d8576108d861259b565b60028111156108e9576108e961259b565b81528154610100810463ffffffff908116602080850191909152650100000000008304909116604080850191909152690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff166060808501919091526001909401546bffffffffffffffffffffffff90811660809485015285519186015192860151948601519590930151909b919a50929850929650169350915050565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146109f6576040517f018d10be00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81818080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050505060208101517fffffffff0000000000000000000000000000000000000000000000000000000081167fa611ea560000000000000000000000000000000000000000000000000000000014610aac576040517fe3d6792100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8484846000610abe82600481866123f7565b810190610acb9190611f10565b50975050505050505050806bffffffffffffffffffffffff168414610b1c576040517f55e97b0d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8988886000610b2e82600481866123f7565b810190610b3b9190611f10565b985050505050505050508073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1614610baa576040517ff8c5638e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6101248b1015610be6576040517fdfe9309000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6004546bffffffffffffffffffffffff168d1015610c30576040517fcd1c886700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60003073ffffffffffffffffffffffffffffffffffffffff168d8d604051610c59929190612170565b600060405180830381855af49150503d8060008114610c94576040519150601f19603f3d011682016040523d82523d6000602084013e610c99565b606091505b5050905080610cd4576040517f649bf81000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505050505050505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610d54576040517f018d10be00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610e696040518061010001604052808e81526020018d8d8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525073ffffffffffffffffffffffffffffffffffffffff808d1660208084019190915263ffffffff8d16604080850191909152918c1660608401528151601f8b018290048202810182019092528982526080909201918a908a9081908401838280828437600092019190915250505090825250604080516020601f8901819004810282018101909252878152918101919088908890819084018382808284376000920191909152505050908252506bffffffffffffffffffffffff85166020909101528261133a565b50505050505050505050505050565b610e80611622565b6003546040805160a081019091526501000000000090910463ffffffff169080866002811115610eb257610eb261259b565b815261ffff8616602082015263ffffffff8316604082015273ffffffffffffffffffffffffffffffffffffffff851660608201526bffffffffffffffffffffffff841660809091015280516003805490919082907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001836002811115610f3c57610f3c61259b565b02179055506020820151815460408085015160608601517fffffffffffffffffffffffffffffffffffffffffffffff0000000000000000ff90931661010063ffffffff958616027fffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffffff1617650100000000009490911693909302929092177fffffff0000000000000000000000000000000000000000ffffffffffffffffff16690100000000000000000073ffffffffffffffffffffffffffffffffffffffff90921691909102178255608090920151600190910180547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff909216919091179055517f6293a703ec7145dfa23c5cde2e627d6a02e153fc2e9c03b14d1e22cbb4a7e9cd9061107e908790879087908790612296565b60405180910390a15050505050565b60008181526002602090815260409182902082518084019093525473ffffffffffffffffffffffffffffffffffffffff8116808452740100000000000000000000000000000000000000009091046bffffffffffffffffffffffff1691830191909152331480611114575060005473ffffffffffffffffffffffffffffffffffffffff1633145b61114a576040517f61685c2b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805173ffffffffffffffffffffffffffffffffffffffff16611198576040517f4b13b31e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260026020908152604080832083905583519184015190517fa9059cbb0000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169263a9059cbb926112509260040173ffffffffffffffffffffffffffffffffffffffff9290921682526bffffffffffffffffffffffff16602082015260400190565b602060405180830381600087803b15801561126a57600080fd5b505af115801561127e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112a29190611c75565b9050806112f65781516040517fc2e4dce800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610821565b60405183907f3663fb28ebc87645eb972c9dad8521bf665c623f287e79f1c56f1eb374b82a2290600090a2505050565b61132e611622565b611337816118ec565b50565b608082015160009073ffffffffffffffffffffffffffffffffffffffff1661138e576040517f05bb467c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008360400151846060015185608001518660a001518760c001516040516020016113bd9594939291906121e7565b604051602081830303815290604052805190602001209050836040015173ffffffffffffffffffffffffffffffffffffffff16817f9b8456f925542af2c5fb15ff4be32cc8f209dda96c544766e301367df40f499886600001518760200151886060015189608001518a60a001518b60e001516040516114429695949392919061234a565b60405180910390a36040805160a081019091526003805460009283929091829060ff1660028111156114765761147661259b565b60028111156114875761148761259b565b8152815463ffffffff61010082048116602084015265010000000000820416604083015273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff1660809091015290506114fa81866119e2565b1561155f57604081015161150f906001612421565b6003805463ffffffff9290921665010000000000027fffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffffff90921691909117905561155886846116a5565b9150611619565b60e086015160008481526002602052604081205490916115a4917401000000000000000000000000000000000000000090046bffffffffffffffffffffffff16612449565b60408051808201825260808a015173ffffffffffffffffffffffffffffffffffffffff90811682526bffffffffffffffffffffffff938416602080840191825260008a815260029091529390932091519251909316740100000000000000000000000000000000000000000291909216179055505b50949350505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146116a3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610821565b565b6003546040838101516060850151608086015160a087015160c088015194517f6ded9eae0000000000000000000000000000000000000000000000000000000081526000966901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff169587958795636ded9eae9561172b95929491939092916004016121e7565b602060405180830381600087803b15801561174557600080fd5b505af1158015611759573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061177d919061204d565b905060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16634000aea0848860e00151856040516020016117d691815260200190565b6040516020818303038152906040526040518463ffffffff1660e01b81526004016118039392919061224a565b602060405180830381600087803b15801561181d57600080fd5b505af1158015611831573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906118559190611c75565b9050806118a6576040517fc2e4dce800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff84166004820152602401610821565b81857fb9a292fb7e3edd920cd2d2829a3615a640c43fd7de0a0820aa0668feb4c37d4b88600001516040516118db9190612337565b60405180910390a350949350505050565b73ffffffffffffffffffffffffffffffffffffffff811633141561196c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610821565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600080835160028111156119f8576119f861259b565b1415611a06575060006104f8565b600183516002811115611a1b57611a1b61259b565b148015611a4e575073ffffffffffffffffffffffffffffffffffffffff821660009081526005602052604090205460ff16155b15611a5b575060006104f8565b826020015163ffffffff16836040015163ffffffff161015611a7f575060016104f8565b50600092915050565b8035611a93816125f9565b919050565b60008083601f840112611aaa57600080fd5b50813567ffffffffffffffff811115611ac257600080fd5b602083019150836020828501011115611ada57600080fd5b9250929050565b600082601f830112611af257600080fd5b813567ffffffffffffffff80821115611b0d57611b0d6125ca565b604051601f83017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908282118183101715611b5357611b536125ca565b81604052838152866020858801011115611b6c57600080fd5b836020870160208301376000602085830101528094505050505092915050565b803563ffffffff81168114611a9357600080fd5b80356bffffffffffffffffffffffff81168114611a9357600080fd5b600060208284031215611bce57600080fd5b8135611bd9816125f9565b9392505050565b60008060408385031215611bf357600080fd5b8235611bfe816125f9565b91506020830135611c0e8161261b565b809150509250929050565b60008060008060608587031215611c2f57600080fd5b8435611c3a816125f9565b935060208501359250604085013567ffffffffffffffff811115611c5d57600080fd5b611c6987828801611a98565b95989497509550505050565b600060208284031215611c8757600080fd5b8151611bd98161261b565b600060208284031215611ca457600080fd5b5035919050565b60008060008060808587031215611cc157600080fd5b843560038110611cd057600080fd5b9350602085013561ffff81168114611ce757600080fd5b92506040850135611cf7816125f9565b9150611d0560608601611ba0565b905092959194509250565b600080600080600080600080600060e08a8c031215611d2e57600080fd5b893567ffffffffffffffff80821115611d4657600080fd5b611d528d838e01611ae1565b9a5060208c01359150611d64826125f9565b819950611d7360408d01611b8c565b985060608c01359150611d85826125f9565b90965060808b01359080821115611d9b57600080fd5b611da78d838e01611a98565b909750955060a08c0135915080821115611dc057600080fd5b50611dcd8c828d01611a98565b9a9d999c50979a9699959894979660c00135949350505050565b6000806000806000806000806000806000806101208d8f031215611e0a57600080fd5b67ffffffffffffffff8d351115611e2057600080fd5b611e2d8e8e358f01611ae1565b9b5067ffffffffffffffff60208e01351115611e4857600080fd5b611e588e60208f01358f01611a98565b909b509950611e6960408e01611a88565b9850611e7760608e01611b8c565b9750611e8560808e01611a88565b965067ffffffffffffffff60a08e01351115611ea057600080fd5b611eb08e60a08f01358f01611a98565b909650945067ffffffffffffffff60c08e01351115611ece57600080fd5b611ede8e60c08f01358f01611a98565b9094509250611eef60e08e01611ba0565b9150611efe6101008e01611a88565b90509295989b509295989b509295989b565b60008060008060008060008060006101208a8c031215611f2f57600080fd5b893567ffffffffffffffff80821115611f4757600080fd5b611f538d838e01611ae1565b9a5060208c0135915080821115611f6957600080fd5b611f758d838e01611ae1565b9950611f8360408d01611a88565b9850611f9160608d01611b8c565b9750611f9f60808d01611a88565b965060a08c0135915080821115611fb557600080fd5b611fc18d838e01611ae1565b955060c08c0135915080821115611fd757600080fd5b50611fe48c828d01611ae1565b935050611ff360e08b01611ba0565b91506120026101008b01611a88565b90509295985092959850929598565b60006020828403121561202357600080fd5b813567ffffffffffffffff81111561203a57600080fd5b82016101008185031215611bd957600080fd5b60006020828403121561205f57600080fd5b5051919050565b60006020828403121561207857600080fd5b611bd982611ba0565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b6000815180845260005b818110156120f0576020818501810151868301820152016120d4565b81811115612102576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6003811061216c577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b9052565b8183823760009101908152919050565b600073ffffffffffffffffffffffffffffffffffffffff808a16835263ffffffff8916602084015280881660408401525060a060608301526121c660a083018688612081565b82810360808401526121d9818587612081565b9a9950505050505050505050565b600073ffffffffffffffffffffffffffffffffffffffff808816835263ffffffff8716602084015280861660408401525060a0606083015261222c60a08301856120ca565b828103608084015261223e81856120ca565b98975050505050505050565b73ffffffffffffffffffffffffffffffffffffffff841681526bffffffffffffffffffffffff8316602082015260606040820152600061228d60608301846120ca565b95945050505050565b608081016122a48287612135565b61ffff8516602083015273ffffffffffffffffffffffffffffffffffffffff841660408301526bffffffffffffffffffffffff8316606083015295945050505050565b60a081016122f58288612135565b63ffffffff808716602084015280861660408401525073ffffffffffffffffffffffffffffffffffffffff841660608301528260808301529695505050505050565b602081526000611bd960208301846120ca565b60c08152600061235d60c08301896120ca565b828103602084015261236f81896120ca565b905063ffffffff8716604084015273ffffffffffffffffffffffffffffffffffffffff8616606084015282810360808401526123ab81866120ca565b9150506bffffffffffffffffffffffff831660a0830152979650505050505050565b604051610100810167ffffffffffffffff811182821017156123f1576123f16125ca565b60405290565b6000808585111561240757600080fd5b8386111561241457600080fd5b5050820193919092039150565b600063ffffffff8083168185168083038211156124405761244061256c565b01949350505050565b60006bffffffffffffffffffffffff8083168185168083038211156124405761244061256c565b6000610100823603121561248357600080fd5b61248b6123cd565b823567ffffffffffffffff808211156124a357600080fd5b6124af36838701611ae1565b835260208501359150808211156124c557600080fd5b6124d136838701611ae1565b60208401526124e260408601611a88565b60408401526124f360608601611b8c565b606084015261250460808601611a88565b608084015260a085013591508082111561251d57600080fd5b61252936838701611ae1565b60a084015260c085013591508082111561254257600080fd5b5061254f36828601611ae1565b60c08301525061256160e08401611ba0565b60e082015292915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff8116811461133757600080fd5b801515811461133757600080fdfea164736f6c6343000806000a",
}

var KeeperRegistrarABI = KeeperRegistrarMetaData.ABI

var KeeperRegistrarBin = KeeperRegistrarMetaData.Bin

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

type KeeperRegistrar struct {
	address common.Address
	abi     abi.ABI
	KeeperRegistrarCaller
	KeeperRegistrarTransactor
	KeeperRegistrarFilterer
}

type KeeperRegistrarCaller struct {
	contract *bind.BoundContract
}

type KeeperRegistrarTransactor struct {
	contract *bind.BoundContract
}

type KeeperRegistrarFilterer struct {
	contract *bind.BoundContract
}

type KeeperRegistrarSession struct {
	Contract     *KeeperRegistrar
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeeperRegistrarCallerSession struct {
	Contract *KeeperRegistrarCaller
	CallOpts bind.CallOpts
}

type KeeperRegistrarTransactorSession struct {
	Contract     *KeeperRegistrarTransactor
	TransactOpts bind.TransactOpts
}

type KeeperRegistrarRaw struct {
	Contract *KeeperRegistrar
}

type KeeperRegistrarCallerRaw struct {
	Contract *KeeperRegistrarCaller
}

type KeeperRegistrarTransactorRaw struct {
	Contract *KeeperRegistrarTransactor
}

func NewKeeperRegistrar(address common.Address, backend bind.ContractBackend) (*KeeperRegistrar, error) {
	abi, err := abi.JSON(strings.NewReader(KeeperRegistrarABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeeperRegistrar(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrar{address: address, abi: abi, KeeperRegistrarCaller: KeeperRegistrarCaller{contract: contract}, KeeperRegistrarTransactor: KeeperRegistrarTransactor{contract: contract}, KeeperRegistrarFilterer: KeeperRegistrarFilterer{contract: contract}}, nil
}

func NewKeeperRegistrarCaller(address common.Address, caller bind.ContractCaller) (*KeeperRegistrarCaller, error) {
	contract, err := bindKeeperRegistrar(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarCaller{contract: contract}, nil
}

func NewKeeperRegistrarTransactor(address common.Address, transactor bind.ContractTransactor) (*KeeperRegistrarTransactor, error) {
	contract, err := bindKeeperRegistrar(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarTransactor{contract: contract}, nil
}

func NewKeeperRegistrarFilterer(address common.Address, filterer bind.ContractFilterer) (*KeeperRegistrarFilterer, error) {
	contract, err := bindKeeperRegistrar(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarFilterer{contract: contract}, nil
}

func bindKeeperRegistrar(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeeperRegistrarMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_KeeperRegistrar *KeeperRegistrarRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistrar.Contract.KeeperRegistrarCaller.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistrar *KeeperRegistrarRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.KeeperRegistrarTransactor.contract.Transfer(opts)
}

func (_KeeperRegistrar *KeeperRegistrarRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.KeeperRegistrarTransactor.contract.Transact(opts, method, params...)
}

func (_KeeperRegistrar *KeeperRegistrarCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistrar.Contract.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistrar *KeeperRegistrarTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.contract.Transfer(opts)
}

func (_KeeperRegistrar *KeeperRegistrarTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.contract.Transact(opts, method, params...)
}

func (_KeeperRegistrar *KeeperRegistrarCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistrar.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistrar *KeeperRegistrarSession) LINK() (common.Address, error) {
	return _KeeperRegistrar.Contract.LINK(&_KeeperRegistrar.CallOpts)
}

func (_KeeperRegistrar *KeeperRegistrarCallerSession) LINK() (common.Address, error) {
	return _KeeperRegistrar.Contract.LINK(&_KeeperRegistrar.CallOpts)
}

func (_KeeperRegistrar *KeeperRegistrarCaller) GetAutoApproveAllowedSender(opts *bind.CallOpts, senderAddress common.Address) (bool, error) {
	var out []interface{}
	err := _KeeperRegistrar.contract.Call(opts, &out, "getAutoApproveAllowedSender", senderAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_KeeperRegistrar *KeeperRegistrarSession) GetAutoApproveAllowedSender(senderAddress common.Address) (bool, error) {
	return _KeeperRegistrar.Contract.GetAutoApproveAllowedSender(&_KeeperRegistrar.CallOpts, senderAddress)
}

func (_KeeperRegistrar *KeeperRegistrarCallerSession) GetAutoApproveAllowedSender(senderAddress common.Address) (bool, error) {
	return _KeeperRegistrar.Contract.GetAutoApproveAllowedSender(&_KeeperRegistrar.CallOpts, senderAddress)
}

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

func (_KeeperRegistrar *KeeperRegistrarSession) GetPendingRequest(hash [32]byte) (common.Address, *big.Int, error) {
	return _KeeperRegistrar.Contract.GetPendingRequest(&_KeeperRegistrar.CallOpts, hash)
}

func (_KeeperRegistrar *KeeperRegistrarCallerSession) GetPendingRequest(hash [32]byte) (common.Address, *big.Int, error) {
	return _KeeperRegistrar.Contract.GetPendingRequest(&_KeeperRegistrar.CallOpts, hash)
}

func (_KeeperRegistrar *KeeperRegistrarCaller) GetRegistrationConfig(opts *bind.CallOpts) (GetRegistrationConfig,

	error) {
	var out []interface{}
	err := _KeeperRegistrar.contract.Call(opts, &out, "getRegistrationConfig")

	outstruct := new(GetRegistrationConfig)
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

func (_KeeperRegistrar *KeeperRegistrarSession) GetRegistrationConfig() (GetRegistrationConfig,

	error) {
	return _KeeperRegistrar.Contract.GetRegistrationConfig(&_KeeperRegistrar.CallOpts)
}

func (_KeeperRegistrar *KeeperRegistrarCallerSession) GetRegistrationConfig() (GetRegistrationConfig,

	error) {
	return _KeeperRegistrar.Contract.GetRegistrationConfig(&_KeeperRegistrar.CallOpts)
}

func (_KeeperRegistrar *KeeperRegistrarCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistrar.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistrar *KeeperRegistrarSession) Owner() (common.Address, error) {
	return _KeeperRegistrar.Contract.Owner(&_KeeperRegistrar.CallOpts)
}

func (_KeeperRegistrar *KeeperRegistrarCallerSession) Owner() (common.Address, error) {
	return _KeeperRegistrar.Contract.Owner(&_KeeperRegistrar.CallOpts)
}

func (_KeeperRegistrar *KeeperRegistrarCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _KeeperRegistrar.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_KeeperRegistrar *KeeperRegistrarSession) TypeAndVersion() (string, error) {
	return _KeeperRegistrar.Contract.TypeAndVersion(&_KeeperRegistrar.CallOpts)
}

func (_KeeperRegistrar *KeeperRegistrarCallerSession) TypeAndVersion() (string, error) {
	return _KeeperRegistrar.Contract.TypeAndVersion(&_KeeperRegistrar.CallOpts)
}

func (_KeeperRegistrar *KeeperRegistrarTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "acceptOwnership")
}

func (_KeeperRegistrar *KeeperRegistrarSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.AcceptOwnership(&_KeeperRegistrar.TransactOpts)
}

func (_KeeperRegistrar *KeeperRegistrarTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.AcceptOwnership(&_KeeperRegistrar.TransactOpts)
}

func (_KeeperRegistrar *KeeperRegistrarTransactor) Approve(opts *bind.TransactOpts, name string, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, offchainConfig []byte, hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "approve", name, upkeepContract, gasLimit, adminAddress, checkData, offchainConfig, hash)
}

func (_KeeperRegistrar *KeeperRegistrarSession) Approve(name string, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, offchainConfig []byte, hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.Approve(&_KeeperRegistrar.TransactOpts, name, upkeepContract, gasLimit, adminAddress, checkData, offchainConfig, hash)
}

func (_KeeperRegistrar *KeeperRegistrarTransactorSession) Approve(name string, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, offchainConfig []byte, hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.Approve(&_KeeperRegistrar.TransactOpts, name, upkeepContract, gasLimit, adminAddress, checkData, offchainConfig, hash)
}

func (_KeeperRegistrar *KeeperRegistrarTransactor) Cancel(opts *bind.TransactOpts, hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "cancel", hash)
}

func (_KeeperRegistrar *KeeperRegistrarSession) Cancel(hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.Cancel(&_KeeperRegistrar.TransactOpts, hash)
}

func (_KeeperRegistrar *KeeperRegistrarTransactorSession) Cancel(hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.Cancel(&_KeeperRegistrar.TransactOpts, hash)
}

func (_KeeperRegistrar *KeeperRegistrarTransactor) OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "onTokenTransfer", sender, amount, data)
}

func (_KeeperRegistrar *KeeperRegistrarSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.OnTokenTransfer(&_KeeperRegistrar.TransactOpts, sender, amount, data)
}

func (_KeeperRegistrar *KeeperRegistrarTransactorSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.OnTokenTransfer(&_KeeperRegistrar.TransactOpts, sender, amount, data)
}

func (_KeeperRegistrar *KeeperRegistrarTransactor) Register(opts *bind.TransactOpts, name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, offchainConfig []byte, amount *big.Int, sender common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "register", name, encryptedEmail, upkeepContract, gasLimit, adminAddress, checkData, offchainConfig, amount, sender)
}

func (_KeeperRegistrar *KeeperRegistrarSession) Register(name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, offchainConfig []byte, amount *big.Int, sender common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.Register(&_KeeperRegistrar.TransactOpts, name, encryptedEmail, upkeepContract, gasLimit, adminAddress, checkData, offchainConfig, amount, sender)
}

func (_KeeperRegistrar *KeeperRegistrarTransactorSession) Register(name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, offchainConfig []byte, amount *big.Int, sender common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.Register(&_KeeperRegistrar.TransactOpts, name, encryptedEmail, upkeepContract, gasLimit, adminAddress, checkData, offchainConfig, amount, sender)
}

func (_KeeperRegistrar *KeeperRegistrarTransactor) RegisterUpkeep(opts *bind.TransactOpts, requestParams KeeperRegistrar20RegistrationParams) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "registerUpkeep", requestParams)
}

func (_KeeperRegistrar *KeeperRegistrarSession) RegisterUpkeep(requestParams KeeperRegistrar20RegistrationParams) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.RegisterUpkeep(&_KeeperRegistrar.TransactOpts, requestParams)
}

func (_KeeperRegistrar *KeeperRegistrarTransactorSession) RegisterUpkeep(requestParams KeeperRegistrar20RegistrationParams) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.RegisterUpkeep(&_KeeperRegistrar.TransactOpts, requestParams)
}

func (_KeeperRegistrar *KeeperRegistrarTransactor) SetAutoApproveAllowedSender(opts *bind.TransactOpts, senderAddress common.Address, allowed bool) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "setAutoApproveAllowedSender", senderAddress, allowed)
}

func (_KeeperRegistrar *KeeperRegistrarSession) SetAutoApproveAllowedSender(senderAddress common.Address, allowed bool) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.SetAutoApproveAllowedSender(&_KeeperRegistrar.TransactOpts, senderAddress, allowed)
}

func (_KeeperRegistrar *KeeperRegistrarTransactorSession) SetAutoApproveAllowedSender(senderAddress common.Address, allowed bool) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.SetAutoApproveAllowedSender(&_KeeperRegistrar.TransactOpts, senderAddress, allowed)
}

func (_KeeperRegistrar *KeeperRegistrarTransactor) SetRegistrationConfig(opts *bind.TransactOpts, autoApproveConfigType uint8, autoApproveMaxAllowed uint16, keeperRegistry common.Address, minLINKJuels *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "setRegistrationConfig", autoApproveConfigType, autoApproveMaxAllowed, keeperRegistry, minLINKJuels)
}

func (_KeeperRegistrar *KeeperRegistrarSession) SetRegistrationConfig(autoApproveConfigType uint8, autoApproveMaxAllowed uint16, keeperRegistry common.Address, minLINKJuels *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.SetRegistrationConfig(&_KeeperRegistrar.TransactOpts, autoApproveConfigType, autoApproveMaxAllowed, keeperRegistry, minLINKJuels)
}

func (_KeeperRegistrar *KeeperRegistrarTransactorSession) SetRegistrationConfig(autoApproveConfigType uint8, autoApproveMaxAllowed uint16, keeperRegistry common.Address, minLINKJuels *big.Int) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.SetRegistrationConfig(&_KeeperRegistrar.TransactOpts, autoApproveConfigType, autoApproveMaxAllowed, keeperRegistry, minLINKJuels)
}

func (_KeeperRegistrar *KeeperRegistrarTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "transferOwnership", to)
}

func (_KeeperRegistrar *KeeperRegistrarSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.TransferOwnership(&_KeeperRegistrar.TransactOpts, to)
}

func (_KeeperRegistrar *KeeperRegistrarTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.TransferOwnership(&_KeeperRegistrar.TransactOpts, to)
}

type KeeperRegistrarAutoApproveAllowedSenderSetIterator struct {
	Event *KeeperRegistrarAutoApproveAllowedSenderSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistrarAutoApproveAllowedSenderSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

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

func (it *KeeperRegistrarAutoApproveAllowedSenderSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistrarAutoApproveAllowedSenderSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistrarAutoApproveAllowedSenderSet struct {
	SenderAddress common.Address
	Allowed       bool
	Raw           types.Log
}

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

func (_KeeperRegistrar *KeeperRegistrarFilterer) ParseAutoApproveAllowedSenderSet(log types.Log) (*KeeperRegistrarAutoApproveAllowedSenderSet, error) {
	event := new(KeeperRegistrarAutoApproveAllowedSenderSet)
	if err := _KeeperRegistrar.contract.UnpackLog(event, "AutoApproveAllowedSenderSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistrarConfigChangedIterator struct {
	Event *KeeperRegistrarConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistrarConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

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

func (it *KeeperRegistrarConfigChangedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistrarConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistrarConfigChanged struct {
	AutoApproveConfigType uint8
	AutoApproveMaxAllowed uint32
	KeeperRegistry        common.Address
	MinLINKJuels          *big.Int
	Raw                   types.Log
}

func (_KeeperRegistrar *KeeperRegistrarFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*KeeperRegistrarConfigChangedIterator, error) {

	logs, sub, err := _KeeperRegistrar.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarConfigChangedIterator{contract: _KeeperRegistrar.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

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

func (_KeeperRegistrar *KeeperRegistrarFilterer) ParseConfigChanged(log types.Log) (*KeeperRegistrarConfigChanged, error) {
	event := new(KeeperRegistrarConfigChanged)
	if err := _KeeperRegistrar.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistrarOwnershipTransferRequestedIterator struct {
	Event *KeeperRegistrarOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistrarOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

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

func (it *KeeperRegistrarOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistrarOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistrarOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

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

func (_KeeperRegistrar *KeeperRegistrarFilterer) ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistrarOwnershipTransferRequested, error) {
	event := new(KeeperRegistrarOwnershipTransferRequested)
	if err := _KeeperRegistrar.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistrarOwnershipTransferredIterator struct {
	Event *KeeperRegistrarOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistrarOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

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

func (it *KeeperRegistrarOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistrarOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistrarOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

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

func (_KeeperRegistrar *KeeperRegistrarFilterer) ParseOwnershipTransferred(log types.Log) (*KeeperRegistrarOwnershipTransferred, error) {
	event := new(KeeperRegistrarOwnershipTransferred)
	if err := _KeeperRegistrar.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistrarRegistrationApprovedIterator struct {
	Event *KeeperRegistrarRegistrationApproved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistrarRegistrationApprovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

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

func (it *KeeperRegistrarRegistrationApprovedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistrarRegistrationApprovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistrarRegistrationApproved struct {
	Hash        [32]byte
	DisplayName string
	UpkeepId    *big.Int
	Raw         types.Log
}

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

func (_KeeperRegistrar *KeeperRegistrarFilterer) ParseRegistrationApproved(log types.Log) (*KeeperRegistrarRegistrationApproved, error) {
	event := new(KeeperRegistrarRegistrationApproved)
	if err := _KeeperRegistrar.contract.UnpackLog(event, "RegistrationApproved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistrarRegistrationRejectedIterator struct {
	Event *KeeperRegistrarRegistrationRejected

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistrarRegistrationRejectedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

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

func (it *KeeperRegistrarRegistrationRejectedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistrarRegistrationRejectedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistrarRegistrationRejected struct {
	Hash [32]byte
	Raw  types.Log
}

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

func (_KeeperRegistrar *KeeperRegistrarFilterer) ParseRegistrationRejected(log types.Log) (*KeeperRegistrarRegistrationRejected, error) {
	event := new(KeeperRegistrarRegistrationRejected)
	if err := _KeeperRegistrar.contract.UnpackLog(event, "RegistrationRejected", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistrarRegistrationRequestedIterator struct {
	Event *KeeperRegistrarRegistrationRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistrarRegistrationRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

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

func (it *KeeperRegistrarRegistrationRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistrarRegistrationRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistrarRegistrationRequested struct {
	Hash           [32]byte
	Name           string
	EncryptedEmail []byte
	UpkeepContract common.Address
	GasLimit       uint32
	AdminAddress   common.Address
	CheckData      []byte
	Amount         *big.Int
	Raw            types.Log
}

func (_KeeperRegistrar *KeeperRegistrarFilterer) FilterRegistrationRequested(opts *bind.FilterOpts, hash [][32]byte, upkeepContract []common.Address) (*KeeperRegistrarRegistrationRequestedIterator, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	var upkeepContractRule []interface{}
	for _, upkeepContractItem := range upkeepContract {
		upkeepContractRule = append(upkeepContractRule, upkeepContractItem)
	}

	logs, sub, err := _KeeperRegistrar.contract.FilterLogs(opts, "RegistrationRequested", hashRule, upkeepContractRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistrarRegistrationRequestedIterator{contract: _KeeperRegistrar.contract, event: "RegistrationRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistrar *KeeperRegistrarFilterer) WatchRegistrationRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarRegistrationRequested, hash [][32]byte, upkeepContract []common.Address) (event.Subscription, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	var upkeepContractRule []interface{}
	for _, upkeepContractItem := range upkeepContract {
		upkeepContractRule = append(upkeepContractRule, upkeepContractItem)
	}

	logs, sub, err := _KeeperRegistrar.contract.WatchLogs(opts, "RegistrationRequested", hashRule, upkeepContractRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

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

func (_KeeperRegistrar *KeeperRegistrarFilterer) ParseRegistrationRequested(log types.Log) (*KeeperRegistrarRegistrationRequested, error) {
	event := new(KeeperRegistrarRegistrationRequested)
	if err := _KeeperRegistrar.contract.UnpackLog(event, "RegistrationRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetRegistrationConfig struct {
	AutoApproveConfigType uint8
	AutoApproveMaxAllowed uint32
	ApprovedCount         uint32
	KeeperRegistry        common.Address
	MinLINKJuels          *big.Int
}

func (_KeeperRegistrar *KeeperRegistrar) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _KeeperRegistrar.abi.Events["AutoApproveAllowedSenderSet"].ID:
		return _KeeperRegistrar.ParseAutoApproveAllowedSenderSet(log)
	case _KeeperRegistrar.abi.Events["ConfigChanged"].ID:
		return _KeeperRegistrar.ParseConfigChanged(log)
	case _KeeperRegistrar.abi.Events["OwnershipTransferRequested"].ID:
		return _KeeperRegistrar.ParseOwnershipTransferRequested(log)
	case _KeeperRegistrar.abi.Events["OwnershipTransferred"].ID:
		return _KeeperRegistrar.ParseOwnershipTransferred(log)
	case _KeeperRegistrar.abi.Events["RegistrationApproved"].ID:
		return _KeeperRegistrar.ParseRegistrationApproved(log)
	case _KeeperRegistrar.abi.Events["RegistrationRejected"].ID:
		return _KeeperRegistrar.ParseRegistrationRejected(log)
	case _KeeperRegistrar.abi.Events["RegistrationRequested"].ID:
		return _KeeperRegistrar.ParseRegistrationRequested(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeeperRegistrarAutoApproveAllowedSenderSet) Topic() common.Hash {
	return common.HexToHash("0x20c6237dac83526a849285a9f79d08a483291bdd3a056a0ef9ae94ecee1ad356")
}

func (KeeperRegistrarConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x6293a703ec7145dfa23c5cde2e627d6a02e153fc2e9c03b14d1e22cbb4a7e9cd")
}

func (KeeperRegistrarOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (KeeperRegistrarOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (KeeperRegistrarRegistrationApproved) Topic() common.Hash {
	return common.HexToHash("0xb9a292fb7e3edd920cd2d2829a3615a640c43fd7de0a0820aa0668feb4c37d4b")
}

func (KeeperRegistrarRegistrationRejected) Topic() common.Hash {
	return common.HexToHash("0x3663fb28ebc87645eb972c9dad8521bf665c623f287e79f1c56f1eb374b82a22")
}

func (KeeperRegistrarRegistrationRequested) Topic() common.Hash {
	return common.HexToHash("0x9b8456f925542af2c5fb15ff4be32cc8f209dda96c544766e301367df40f4998")
}

func (_KeeperRegistrar *KeeperRegistrar) Address() common.Address {
	return _KeeperRegistrar.address
}

type KeeperRegistrarInterface interface {
	LINK(opts *bind.CallOpts) (common.Address, error)

	GetAutoApproveAllowedSender(opts *bind.CallOpts, senderAddress common.Address) (bool, error)

	GetPendingRequest(opts *bind.CallOpts, hash [32]byte) (common.Address, *big.Int, error)

	GetRegistrationConfig(opts *bind.CallOpts) (GetRegistrationConfig,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	Approve(opts *bind.TransactOpts, name string, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, offchainConfig []byte, hash [32]byte) (*types.Transaction, error)

	Cancel(opts *bind.TransactOpts, hash [32]byte) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	Register(opts *bind.TransactOpts, name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, offchainConfig []byte, amount *big.Int, sender common.Address) (*types.Transaction, error)

	RegisterUpkeep(opts *bind.TransactOpts, requestParams KeeperRegistrar20RegistrationParams) (*types.Transaction, error)

	SetAutoApproveAllowedSender(opts *bind.TransactOpts, senderAddress common.Address, allowed bool) (*types.Transaction, error)

	SetRegistrationConfig(opts *bind.TransactOpts, autoApproveConfigType uint8, autoApproveMaxAllowed uint16, keeperRegistry common.Address, minLINKJuels *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterAutoApproveAllowedSenderSet(opts *bind.FilterOpts, senderAddress []common.Address) (*KeeperRegistrarAutoApproveAllowedSenderSetIterator, error)

	WatchAutoApproveAllowedSenderSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarAutoApproveAllowedSenderSet, senderAddress []common.Address) (event.Subscription, error)

	ParseAutoApproveAllowedSenderSet(log types.Log) (*KeeperRegistrarAutoApproveAllowedSenderSet, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*KeeperRegistrarConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*KeeperRegistrarConfigChanged, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistrarOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistrarOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistrarOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*KeeperRegistrarOwnershipTransferred, error)

	FilterRegistrationApproved(opts *bind.FilterOpts, hash [][32]byte, upkeepId []*big.Int) (*KeeperRegistrarRegistrationApprovedIterator, error)

	WatchRegistrationApproved(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarRegistrationApproved, hash [][32]byte, upkeepId []*big.Int) (event.Subscription, error)

	ParseRegistrationApproved(log types.Log) (*KeeperRegistrarRegistrationApproved, error)

	FilterRegistrationRejected(opts *bind.FilterOpts, hash [][32]byte) (*KeeperRegistrarRegistrationRejectedIterator, error)

	WatchRegistrationRejected(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarRegistrationRejected, hash [][32]byte) (event.Subscription, error)

	ParseRegistrationRejected(log types.Log) (*KeeperRegistrarRegistrationRejected, error)

	FilterRegistrationRequested(opts *bind.FilterOpts, hash [][32]byte, upkeepContract []common.Address) (*KeeperRegistrarRegistrationRequestedIterator, error)

	WatchRegistrationRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarRegistrationRequested, hash [][32]byte, upkeepContract []common.Address) (event.Subscription, error)

	ParseRegistrationRequested(log types.Log) (*KeeperRegistrarRegistrationRequested, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
