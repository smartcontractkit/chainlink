// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package keeper_registrar_wrapper1_2

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

var KeeperRegistrarMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"LINKAddress\",\"type\":\"address\"},{\"internalType\":\"enumKeeperRegistrar.AutoApproveType\",\"name\":\"autoApproveConfigType\",\"type\":\"uint8\"},{\"internalType\":\"uint16\",\"name\":\"autoApproveMaxAllowed\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"minLINKJuels\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AmountMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FunctionNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HashMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientPayment\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAdminAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"LinkTransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyAdminOrOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistrationRequestFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RequestNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SenderMismatch\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"AutoApproveAllowedSenderSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"enumKeeperRegistrar.AutoApproveType\",\"name\":\"autoApproveConfigType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"autoApproveMaxAllowed\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"minLINKJuels\",\"type\":\"uint96\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"displayName\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"RegistrationApproved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"RegistrationRejected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encryptedEmail\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"upkeepContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"source\",\"type\":\"uint8\"}],\"name\":\"RegistrationRequested\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"upkeepContract\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"cancel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"}],\"name\":\"getAutoApproveAllowedSender\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"getPendingRequest\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRegistrationConfig\",\"outputs\":[{\"internalType\":\"enumKeeperRegistrar.AutoApproveType\",\"name\":\"autoApproveConfigType\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"autoApproveMaxAllowed\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"approvedCount\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minLINKJuels\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"encryptedEmail\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"upkeepContract\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"internalType\":\"uint8\",\"name\":\"source\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"setAutoApproveAllowedSender\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumKeeperRegistrar.AutoApproveType\",\"name\":\"autoApproveConfigType\",\"type\":\"uint8\"},{\"internalType\":\"uint16\",\"name\":\"autoApproveMaxAllowed\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"keeperRegistry\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"minLINKJuels\",\"type\":\"uint96\"}],\"name\":\"setRegistrationConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b506040516200234638038062002346833981016040819052620000349162000394565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000ec565b5050506001600160601b0319606086901b16608052620000e18484848462000198565b50505050506200048d565b6001600160a01b038116331415620001475760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b620001a262000319565b6003546040805160a081019091526501000000000090910463ffffffff169080866002811115620001d757620001d762000477565b815261ffff8616602082015263ffffffff831660408201526001600160a01b03851660608201526001600160601b038416608090910152805160038054909190829060ff1916600183600281111562000234576200023462000477565b0217905550602082015181546040808501516060860151610100600160481b031990931661010063ffffffff9586160263ffffffff60281b19161765010000000000949091169390930292909217600160481b600160e81b03191669010000000000000000006001600160a01b0390921691909102178255608090920151600190910180546001600160601b0319166001600160601b03909216919091179055517f6293a703ec7145dfa23c5cde2e627d6a02e153fc2e9c03b14d1e22cbb4a7e9cd906200030a90879087908790879062000422565b60405180910390a15050505050565b6000546001600160a01b03163314620003755760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000082565b565b80516001600160a01b03811681146200038f57600080fd5b919050565b600080600080600060a08688031215620003ad57600080fd5b620003b88662000377565b9450602086015160038110620003cd57600080fd5b604087015190945061ffff81168114620003e657600080fd5b9250620003f66060870162000377565b60808701519092506001600160601b03811681146200041457600080fd5b809150509295509295909350565b60808101600386106200044557634e487b7160e01b600052602160045260246000fd5b94815261ffff9390931660208401526001600160a01b039190911660408301526001600160601b031660609091015290565b634e487b7160e01b600052602160045260246000fd5b60805160601c611e7e620004c86000396000818161015b015281816104a601528181610a410152818161110b01526113500152611e7e6000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c8063850af0cb1161008c578063a4c0ed3611610066578063a4c0ed36146102fd578063a793ab8b14610310578063c4d252f514610323578063f2fde38b1461033657600080fd5b8063850af0cb1461021957806388b12d55146102325780638da5cb5b146102df57600080fd5b80633659d666116100c85780633659d666146101a2578063367b9b4f146101b557806379ba5097146101c85780637e776f7f146101d057600080fd5b8063181f5a77146100ef578063183310b3146101415780631b6b6d2314610156575b600080fd5b61012b6040518060400160405280601581526020017f4b656570657252656769737472617220312e312e30000000000000000000000081525081565b6040516101389190611cf1565b60405180910390f35b61015461014f3660046118fc565b610349565b005b61017d7f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610138565b6101546101b03660046119a1565b61048e565b6101546101c33660046117d2565b6107b4565b610154610846565b6102096101de3660046117b0565b73ffffffffffffffffffffffffffffffffffffffff1660009081526005602052604090205460ff1690565b6040519015158152602001610138565b610221610948565b604051610138959493929190611ca1565b6102a6610240366004611880565b60009081526002602090815260409182902082518084019093525473ffffffffffffffffffffffffffffffffffffffff8116808452740100000000000000000000000000000000000000009091046bffffffffffffffffffffffff169290910182905291565b6040805173ffffffffffffffffffffffffffffffffffffffff90931683526bffffffffffffffffffffffff909116602083015201610138565b60005473ffffffffffffffffffffffffffffffffffffffff1661017d565b61015461030b366004611809565b610a29565b61015461031e366004611899565b610d7c565b610154610331366004611880565b610f91565b6101546103443660046117b0565b6111f2565b610351611206565b60008181526002602090815260409182902082518084019093525473ffffffffffffffffffffffffffffffffffffffff8116808452740100000000000000000000000000000000000000009091046bffffffffffffffffffffffff16918301919091526103ea576040517f4b13b31e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008787878787604051602001610405959493929190611bb3565b604051602081830303815290604052805190602001209050808314610456576040517f3f4d605300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000838152600260209081526040822091909155820151610483908a908a908a908a908a908a908a611289565b505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146104fd576040517f018d10be00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff861661054a576040517f05bb467c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008888888888604051602001610565959493929190611bb3565b6040516020818303038152906040528051906020012090508260ff168973ffffffffffffffffffffffffffffffffffffffff16827fc3f5df4aefec026f610a3fcb08f19476492d69d2cb78b1c2eba259a8820e6a788f8f8f8e8e8e8e8e6040516105d6989796959493929190611d04565b60405180910390a46040805160a08101909152600380546000929190829060ff16600281111561060857610608611e05565b600281111561061957610619611e05565b8152815463ffffffff61010082048116602084015265010000000000820416604083015273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff16608090910152905061068c81846114b5565b156106f55760408101516106a1906001611d8b565b6003805463ffffffff9290921665010000000000027fffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffffff9092169190911790556106f08d8b8b8b8b8b8b89611289565b6107a5565b6000828152600260205260408120546107359087907401000000000000000000000000000000000000000090046bffffffffffffffffffffffff16611db3565b60408051808201825273ffffffffffffffffffffffffffffffffffffffff808d1682526bffffffffffffffffffffffff9384166020808401918252600089815260029091529390932091519251909316740100000000000000000000000000000000000000000291909216179055505b50505050505050505050505050565b6107bc611206565b73ffffffffffffffffffffffffffffffffffffffff821660008181526005602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001685151590811790915591519182527f20c6237dac83526a849285a9f79d08a483291bdd3a056a0ef9ae94ecee1ad356910160405180910390a25050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146108cc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6040805160a081019091526003805460009283928392839283928392829060ff16600281111561097a5761097a611e05565b600281111561098b5761098b611e05565b81528154610100810463ffffffff908116602080850191909152650100000000008304909116604080850191909152690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff166060808501919091526001909401546bffffffffffffffffffffffff90811660809485015285519186015192860151948601519590930151909b919a50929850929650169350915050565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610a98576040517f018d10be00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81818080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050505060208101517fffffffff0000000000000000000000000000000000000000000000000000000081167f3659d6660000000000000000000000000000000000000000000000000000000014610b4e576040517fe3d6792100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8484848080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050505060e4810151828114610bc3576040517f55e97b0d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8887878080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050505061012481015173ffffffffffffffffffffffffffffffffffffffff83811690821614610c52576040517ff8c5638e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610124891015610c8e576040517fdfe9309000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6004546bffffffffffffffffffffffff168b1015610cd8576040517fcd1c886700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60003073ffffffffffffffffffffffffffffffffffffffff168b8b604051610d01929190611ba3565b600060405180830381855af49150503d8060008114610d3c576040519150601f19603f3d011682016040523d82523d6000602084013e610d41565b606091505b50509050806107a5576040517f649bf81000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610d84611206565b6003546040805160a081019091526501000000000090910463ffffffff169080866002811115610db657610db6611e05565b815261ffff8616602082015263ffffffff8316604082015273ffffffffffffffffffffffffffffffffffffffff851660608201526bffffffffffffffffffffffff841660809091015280516003805490919082907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001836002811115610e4057610e40611e05565b02179055506020820151815460408085015160608601517fffffffffffffffffffffffffffffffffffffffffffffff0000000000000000ff90931661010063ffffffff958616027fffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffffff1617650100000000009490911693909302929092177fffffff0000000000000000000000000000000000000000ffffffffffffffffff16690100000000000000000073ffffffffffffffffffffffffffffffffffffffff90921691909102178255608090920151600190910180547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff909216919091179055517f6293a703ec7145dfa23c5cde2e627d6a02e153fc2e9c03b14d1e22cbb4a7e9cd90610f82908790879087908790611c50565b60405180910390a15050505050565b60008181526002602090815260409182902082518084019093525473ffffffffffffffffffffffffffffffffffffffff8116808452740100000000000000000000000000000000000000009091046bffffffffffffffffffffffff1691830191909152331480611018575060005473ffffffffffffffffffffffffffffffffffffffff1633145b61104e576040517f61685c2b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805173ffffffffffffffffffffffffffffffffffffffff1661109c576040517f4b13b31e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526002602090815260408083208390559083015190517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526bffffffffffffffffffffffff909116602482015273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb90604401602060405180830381600087803b15801561114f57600080fd5b505af1158015611163573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111879190611863565b9050806111c2576040517fc2e4dce80000000000000000000000000000000000000000000000000000000081523360048201526024016108c3565b60405183907f3663fb28ebc87645eb972c9dad8521bf665c623f287e79f1c56f1eb374b82a2290600090a2505050565b6111fa611206565b6112038161155c565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314611287576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016108c3565b565b6003546040517fda5c6741000000000000000000000000000000000000000000000000000000008152690100000000000000000090910473ffffffffffffffffffffffffffffffffffffffff1690600090829063da5c6741906112f8908c908c908c908c908c90600401611bb3565b602060405180830381600087803b15801561131257600080fd5b505af1158015611326573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061134a9190611a9b565b905060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16634000aea084878560405160200161139f91815260200190565b6040516020818303038152906040526040518463ffffffff1660e01b81526004016113cc93929190611c04565b602060405180830381600087803b1580156113e657600080fd5b505af11580156113fa573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061141e9190611863565b90508061146f576040517fc2e4dce800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff841660048201526024016108c3565b81847fb9a292fb7e3edd920cd2d2829a3615a640c43fd7de0a0820aa0668feb4c37d4b8d6040516114a09190611cf1565b60405180910390a35050505050505050505050565b600080835160028111156114cb576114cb611e05565b14156114d957506000611556565b6001835160028111156114ee576114ee611e05565b148015611521575073ffffffffffffffffffffffffffffffffffffffff821660009081526005602052604090205460ff16155b1561152e57506000611556565b826020015163ffffffff16836040015163ffffffff16101561155257506001611556565b5060005b92915050565b73ffffffffffffffffffffffffffffffffffffffff81163314156115dc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016108c3565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b803573ffffffffffffffffffffffffffffffffffffffff8116811461167657600080fd5b919050565b60008083601f84011261168d57600080fd5b50813567ffffffffffffffff8111156116a557600080fd5b6020830191508360208285010111156116bd57600080fd5b9250929050565b600082601f8301126116d557600080fd5b813567ffffffffffffffff808211156116f0576116f0611e34565b604051601f83017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190828211818310171561173657611736611e34565b8160405283815286602085880101111561174f57600080fd5b836020870160208301376000602085830101528094505050505092915050565b803563ffffffff8116811461167657600080fd5b803560ff8116811461167657600080fd5b80356bffffffffffffffffffffffff8116811461167657600080fd5b6000602082840312156117c257600080fd5b6117cb82611652565b9392505050565b600080604083850312156117e557600080fd5b6117ee83611652565b915060208301356117fe81611e63565b809150509250929050565b6000806000806060858703121561181f57600080fd5b61182885611652565b935060208501359250604085013567ffffffffffffffff81111561184b57600080fd5b6118578782880161167b565b95989497509550505050565b60006020828403121561187557600080fd5b81516117cb81611e63565b60006020828403121561189257600080fd5b5035919050565b600080600080608085870312156118af57600080fd5b8435600381106118be57600080fd5b9350602085013561ffff811681146118d557600080fd5b92506118e360408601611652565b91506118f160608601611794565b905092959194509250565b600080600080600080600060c0888a03121561191757600080fd5b873567ffffffffffffffff8082111561192f57600080fd5b61193b8b838c016116c4565b985061194960208b01611652565b975061195760408b0161176f565b965061196560608b01611652565b955060808a013591508082111561197b57600080fd5b506119888a828b0161167b565b989b979a5095989497959660a090950135949350505050565b60008060008060008060008060008060006101208c8e0312156119c357600080fd5b67ffffffffffffffff808d3511156119da57600080fd5b6119e78e8e358f016116c4565b9b508060208e013511156119fa57600080fd5b611a0a8e60208f01358f0161167b565b909b509950611a1b60408e01611652565b9850611a2960608e0161176f565b9750611a3760808e01611652565b96508060a08e01351115611a4a57600080fd5b50611a5b8d60a08e01358e0161167b565b9095509350611a6c60c08d01611794565b9250611a7a60e08d01611783565b9150611a896101008d01611652565b90509295989b509295989b9093969950565b600060208284031215611aad57600080fd5b5051919050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b6000815180845260005b81811015611b2357602081850181015186830182015201611b07565b81811115611b35576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60038110611b9f577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b9052565b8183823760009101908152919050565b600073ffffffffffffffffffffffffffffffffffffffff808816835263ffffffff8716602084015280861660408401525060806060830152611bf9608083018486611ab4565b979650505050505050565b73ffffffffffffffffffffffffffffffffffffffff841681526bffffffffffffffffffffffff83166020820152606060408201526000611c476060830184611afd565b95945050505050565b60808101611c5e8287611b68565b61ffff8516602083015273ffffffffffffffffffffffffffffffffffffffff841660408301526bffffffffffffffffffffffff8316606083015295945050505050565b60a08101611caf8288611b68565b63ffffffff808716602084015280861660408401525073ffffffffffffffffffffffffffffffffffffffff841660608301528260808301529695505050505050565b6020815260006117cb6020830184611afd565b60c081526000611d1760c083018b611afd565b8281036020840152611d2a818a8c611ab4565b905063ffffffff8816604084015273ffffffffffffffffffffffffffffffffffffffff871660608401528281036080840152611d67818688611ab4565b9150506bffffffffffffffffffffffff831660a08301529998505050505050505050565b600063ffffffff808316818516808303821115611daa57611daa611dd6565b01949350505050565b60006bffffffffffffffffffffffff808316818516808303821115611daa57611daa5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b801515811461120357600080fdfea164736f6c6343000806000a",
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

func (_KeeperRegistrar *KeeperRegistrarTransactor) Approve(opts *bind.TransactOpts, name string, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "approve", name, upkeepContract, gasLimit, adminAddress, checkData, hash)
}

func (_KeeperRegistrar *KeeperRegistrarSession) Approve(name string, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.Approve(&_KeeperRegistrar.TransactOpts, name, upkeepContract, gasLimit, adminAddress, checkData, hash)
}

func (_KeeperRegistrar *KeeperRegistrarTransactorSession) Approve(name string, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, hash [32]byte) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.Approve(&_KeeperRegistrar.TransactOpts, name, upkeepContract, gasLimit, adminAddress, checkData, hash)
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

func (_KeeperRegistrar *KeeperRegistrarTransactor) Register(opts *bind.TransactOpts, name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, amount *big.Int, source uint8, sender common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar.contract.Transact(opts, "register", name, encryptedEmail, upkeepContract, gasLimit, adminAddress, checkData, amount, source, sender)
}

func (_KeeperRegistrar *KeeperRegistrarSession) Register(name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, amount *big.Int, source uint8, sender common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.Register(&_KeeperRegistrar.TransactOpts, name, encryptedEmail, upkeepContract, gasLimit, adminAddress, checkData, amount, source, sender)
}

func (_KeeperRegistrar *KeeperRegistrarTransactorSession) Register(name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, amount *big.Int, source uint8, sender common.Address) (*types.Transaction, error) {
	return _KeeperRegistrar.Contract.Register(&_KeeperRegistrar.TransactOpts, name, encryptedEmail, upkeepContract, gasLimit, adminAddress, checkData, amount, source, sender)
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
	Source         uint8
	Raw            types.Log
}

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
	return common.HexToHash("0xc3f5df4aefec026f610a3fcb08f19476492d69d2cb78b1c2eba259a8820e6a78")
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

	Approve(opts *bind.TransactOpts, name string, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, hash [32]byte) (*types.Transaction, error)

	Cancel(opts *bind.TransactOpts, hash [32]byte) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	Register(opts *bind.TransactOpts, name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, amount *big.Int, source uint8, sender common.Address) (*types.Transaction, error)

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

	FilterRegistrationRequested(opts *bind.FilterOpts, hash [][32]byte, upkeepContract []common.Address, source []uint8) (*KeeperRegistrarRegistrationRequestedIterator, error)

	WatchRegistrationRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistrarRegistrationRequested, hash [][32]byte, upkeepContract []common.Address, source []uint8) (event.Subscription, error)

	ParseRegistrationRequested(log types.Log) (*KeeperRegistrarRegistrationRequested, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
