// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package burn_mint_erc20

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

var BurnMintERC20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"maxSupply_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"preMint\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"supplyAfterMint\",\"type\":\"uint256\"}],\"name\":\"MaxSupplyExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotBurner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotMinter\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"burner\",\"type\":\"address\"}],\"name\":\"BurnAccessGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"burner\",\"type\":\"address\"}],\"name\":\"BurnAccessRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"CCIPAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"}],\"name\":\"MintAccessGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"}],\"name\":\"MintAccessRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burnFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseApproval\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBurners\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCCIPAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"burner\",\"type\":\"address\"}],\"name\":\"grantBurnRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"burnAndMinter\",\"type\":\"address\"}],\"name\":\"grantMintAndBurnRoles\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"}],\"name\":\"grantMintRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseApproval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"burner\",\"type\":\"address\"}],\"name\":\"isBurner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"}],\"name\":\"isMinter\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"burner\",\"type\":\"address\"}],\"name\":\"revokeBurnRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"}],\"name\":\"revokeMintRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"setCCIPAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c06040523480156200001157600080fd5b50604051620022f0380380620022f0833981016040819052620000349162000469565b336000868660036200004783826200058d565b5060046200005682826200058d565b5050506001600160a01b0382166200008157604051639b15e16f60e01b815260040160405180910390fd5b600680546001600160a01b0319166001600160a01b0384811691909117909155811615620000b457620000b48162000108565b505060ff831660805260a0829052600780546001600160a01b031916331790558015620000e757620000e7338262000184565b620000f2336200024a565b620000fd33620002a8565b50505050506200067b565b336001600160a01b038216036200013257604051636d6c4ee560e11b815260040160405180910390fd5b600580546001600160a01b0319166001600160a01b03838116918217909255600654604051919216907fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127890600090a350565b6001600160a01b038216620001df5760405162461bcd60e51b815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f206164647265737300604482015260640160405180910390fd5b8060026000828254620001f3919062000659565b90915550506001600160a01b038216600081815260208181526040808320805486019055518481527fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a35050565b6200025462000304565b6200026160088262000332565b15620002a5576040516001600160a01b03821681527fe46fef8bbff1389d9010703cf8ebb363fb3daf5bf56edc27080b67bc8d9251ea906020015b60405180910390a15b50565b620002b262000304565b620002bf600a8262000332565b15620002a5576040516001600160a01b03821681527f92308bb7573b2a3d17ddb868b39d8ebec433f3194421abc22d084f89658c9bad906020016200029c565b505050565b6006546001600160a01b0316331462000330576040516315ae3a6f60e11b815260040160405180910390fd5b565b600062000349836001600160a01b03841662000352565b90505b92915050565b60008181526001830160205260408120546200039b575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556200034c565b5060006200034c565b634e487b7160e01b600052604160045260246000fd5b600082601f830112620003cc57600080fd5b81516001600160401b0380821115620003e957620003e9620003a4565b604051601f8301601f19908116603f01168101908282118183101715620004145762000414620003a4565b816040528381526020925086838588010111156200043157600080fd5b600091505b8382101562000455578582018301518183018401529082019062000436565b600093810190920192909252949350505050565b600080600080600060a086880312156200048257600080fd5b85516001600160401b03808211156200049a57600080fd5b620004a889838a01620003ba565b96506020880151915080821115620004bf57600080fd5b50620004ce88828901620003ba565b945050604086015160ff81168114620004e657600080fd5b6060870151608090970151959894975095949392505050565b600181811c908216806200051457607f821691505b6020821081036200053557634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620002ff57600081815260208120601f850160051c81016020861015620005645750805b601f850160051c820191505b81811015620005855782815560010162000570565b505050505050565b81516001600160401b03811115620005a957620005a9620003a4565b620005c181620005ba8454620004ff565b846200053b565b602080601f831160018114620005f95760008415620005e05750858301515b600019600386901b1c1916600185901b17855562000585565b600085815260208120601f198616915b828110156200062a5788860151825594840194600190910190840162000609565b5085821015620006495787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b808201808211156200034c57634e487b7160e01b600052601160045260246000fd5b60805160a051611c41620006af600039600081816104970152818161083c01526108660152600061028c0152611c416000f3fe608060405234801561001057600080fd5b506004361061020b5760003560e01c806386fe8b431161012a578063aa271e1a116100bd578063d5abeb011161008c578063dd62ed3e11610071578063dd62ed3e146104ce578063f2fde38b14610514578063f81094f31461052757600080fd5b8063d5abeb0114610495578063d73dd623146104bb57600080fd5b8063aa271e1a14610449578063c2e3273d1461045c578063c630948d1461046f578063c64d0ebc1461048257600080fd5b80639dc29fac116100f95780639dc29fac146103fd578063a457c2d714610410578063a8fa343c14610423578063a9059cbb1461043657600080fd5b806386fe8b43146103905780638da5cb5b146103985780638fd6a6ac146103d757806395d89b41146103f557600080fd5b806342966c68116101a25780636b32810b116101715780636b32810b1461032a57806370a082311461033f57806379ba50971461037557806379cc67901461037d57600080fd5b806342966c68146102de5780634334614a146102f15780634f5632f814610304578063661884631461031757600080fd5b806323b872dd116101de57806323b872dd14610272578063313ce5671461028557806339509351146102b657806340c10f19146102c957600080fd5b806301ffc9a71461021057806306fdde0314610238578063095ea7b31461024d57806318160ddd14610260575b600080fd5b61022361021e366004611930565b61053a565b60405190151581526020015b60405180910390f35b6102406106b7565b60405161022f9190611972565b61022361025b366004611a07565b610749565b6002545b60405190815260200161022f565b610223610280366004611a31565b610761565b60405160ff7f000000000000000000000000000000000000000000000000000000000000000016815260200161022f565b6102236102c4366004611a07565b610785565b6102dc6102d7366004611a07565b6107d1565b005b6102dc6102ec366004611a6d565b6108f8565b6102236102ff366004611a86565b610945565b6102dc610312366004611a86565b610952565b610223610325366004611a07565b6109b7565b6103326109ca565b60405161022f9190611aa1565b61026461034d366004611a86565b73ffffffffffffffffffffffffffffffffffffffff1660009081526020819052604090205490565b6102dc6109db565b6102dc61038b366004611a07565b610aac565b610332610afb565b60065473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161022f565b60075473ffffffffffffffffffffffffffffffffffffffff166103b2565b610240610b07565b6102dc61040b366004611a07565b610b16565b61022361041e366004611a07565b610b20565b6102dc610431366004611a86565b610bf1565b610223610444366004611a07565b610c70565b610223610457366004611a86565b610c7e565b6102dc61046a366004611a86565b610c8b565b6102dc61047d366004611a86565b610ce9565b6102dc610490366004611a86565b610cf7565b7f0000000000000000000000000000000000000000000000000000000000000000610264565b6102dc6104c9366004611a07565b610d55565b6102646104dc366004611afb565b73ffffffffffffffffffffffffffffffffffffffff918216600090815260016020908152604080832093909416825291909152205490565b6102dc610522366004611a86565b610d5f565b6102dc610535366004611a86565b610d70565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f36372b070000000000000000000000000000000000000000000000000000000014806105cd57507fffffffff0000000000000000000000000000000000000000000000000000000082167fe6599b4d00000000000000000000000000000000000000000000000000000000145b8061061957507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b8061066557507fffffffff0000000000000000000000000000000000000000000000000000000082167f06e2784700000000000000000000000000000000000000000000000000000000145b806106b157507fffffffff0000000000000000000000000000000000000000000000000000000082167f8fd6a6ac00000000000000000000000000000000000000000000000000000000145b92915050565b6060600380546106c690611b2e565b80601f01602080910402602001604051908101604052809291908181526020018280546106f290611b2e565b801561073f5780601f106107145761010080835404028352916020019161073f565b820191906000526020600020905b81548152906001019060200180831161072257829003601f168201915b5050505050905090565b600033610757818585610dce565b5060019392505050565b60003361076f858285610e02565b61077a858585610ed3565b506001949350505050565b33600081815260016020908152604080832073ffffffffffffffffffffffffffffffffffffffff8716845290915281205490919061075790829086906107cc908790611bb0565b610dce565b6107da33610c7e565b610817576040517fe2c8c9d50000000000000000000000000000000000000000000000000000000081523360048201526024015b60405180910390fd5b813073ffffffffffffffffffffffffffffffffffffffff82160361083a57600080fd5b7f00000000000000000000000000000000000000000000000000000000000000001580159061089b57507f00000000000000000000000000000000000000000000000000000000000000008261088f60025490565b6108999190611bb0565b115b156108e957816108aa60025490565b6108b49190611bb0565b6040517fcbbf111300000000000000000000000000000000000000000000000000000000815260040161080e91815260200190565b6108f38383610f01565b505050565b61090133610945565b610939576040517fc820b10b00000000000000000000000000000000000000000000000000000000815233600482015260240161080e565b61094281610ff4565b50565b60006106b1600a83610ffe565b61095a61102d565b610965600a82611080565b156109425760405173ffffffffffffffffffffffffffffffffffffffff821681527f0a675452746933cefe3d74182e78db7afe57ba60eaa4234b5d85e9aa41b0610c906020015b60405180910390a150565b60006109c38383610b20565b9392505050565b60606109d660086110a2565b905090565b60055473ffffffffffffffffffffffffffffffffffffffff163314610a2c576040517f02b543c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600680547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560058054909116905560405173ffffffffffffffffffffffffffffffffffffffff909116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a350565b610ab533610945565b610aed576040517fc820b10b00000000000000000000000000000000000000000000000000000000815233600482015260240161080e565b610af782826110af565b5050565b60606109d6600a6110a2565b6060600480546106c690611b2e565b610af78282610aac565b33600081815260016020908152604080832073ffffffffffffffffffffffffffffffffffffffff8716845290915281205490919083811015610be4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760448201527f207a65726f000000000000000000000000000000000000000000000000000000606482015260840161080e565b61077a8286868403610dce565b610bf961102d565b6007805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681179093556040519116919082907f9524c9e4b0b61eb018dd58a1cd856e3e74009528328ab4a613b434fa631d724290600090a35050565b600033610757818585610ed3565b60006106b1600883610ffe565b610c9361102d565b610c9e6008826110c4565b156109425760405173ffffffffffffffffffffffffffffffffffffffff821681527fe46fef8bbff1389d9010703cf8ebb363fb3daf5bf56edc27080b67bc8d9251ea906020016109ac565b610cf281610c8b565b610942815b610cff61102d565b610d0a600a826110c4565b156109425760405173ffffffffffffffffffffffffffffffffffffffff821681527f92308bb7573b2a3d17ddb868b39d8ebec433f3194421abc22d084f89658c9bad906020016109ac565b6108f38282610785565b610d6761102d565b610942816110e6565b610d7861102d565b610d83600882611080565b156109425760405173ffffffffffffffffffffffffffffffffffffffff821681527fed998b960f6340d045f620c119730f7aa7995e7425c2401d3a5b64ff998a59e9906020016109ac565b813073ffffffffffffffffffffffffffffffffffffffff821603610df157600080fd5b610dfc8484846111ac565b50505050565b73ffffffffffffffffffffffffffffffffffffffff8381166000908152600160209081526040808320938616835292905220547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8114610dfc5781811015610ec6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f45524332303a20696e73756666696369656e7420616c6c6f77616e6365000000604482015260640161080e565b610dfc8484848403610dce565b813073ffffffffffffffffffffffffffffffffffffffff821603610ef657600080fd5b610dfc84848461135f565b73ffffffffffffffffffffffffffffffffffffffff8216610f7e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f206164647265737300604482015260640161080e565b8060026000828254610f909190611bb0565b909155505073ffffffffffffffffffffffffffffffffffffffff8216600081815260208181526040808320805486019055518481527fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a35050565b61094233826115ce565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415156109c3565b60065473ffffffffffffffffffffffffffffffffffffffff16331461107e576040517f2b5c74de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60006109c38373ffffffffffffffffffffffffffffffffffffffff8416611792565b606060006109c383611885565b6110ba823383610e02565b610af782826115ce565b60006109c38373ffffffffffffffffffffffffffffffffffffffff84166118e1565b3373ffffffffffffffffffffffffffffffffffffffff821603611135576040517fdad89dca00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600580547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217909255600654604051919216907fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127890600090a350565b73ffffffffffffffffffffffffffffffffffffffff831661124e576040517f08c379a0000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460448201527f7265737300000000000000000000000000000000000000000000000000000000606482015260840161080e565b73ffffffffffffffffffffffffffffffffffffffff82166112f1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f20616464726560448201527f7373000000000000000000000000000000000000000000000000000000000000606482015260840161080e565b73ffffffffffffffffffffffffffffffffffffffff83811660008181526001602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b73ffffffffffffffffffffffffffffffffffffffff8316611402576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f45524332303a207472616e736665722066726f6d20746865207a65726f20616460448201527f6472657373000000000000000000000000000000000000000000000000000000606482015260840161080e565b73ffffffffffffffffffffffffffffffffffffffff82166114a5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f45524332303a207472616e7366657220746f20746865207a65726f206164647260448201527f6573730000000000000000000000000000000000000000000000000000000000606482015260840161080e565b73ffffffffffffffffffffffffffffffffffffffff83166000908152602081905260409020548181101561155b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f45524332303a207472616e7366657220616d6f756e742065786365656473206260448201527f616c616e63650000000000000000000000000000000000000000000000000000606482015260840161080e565b73ffffffffffffffffffffffffffffffffffffffff848116600081815260208181526040808320878703905593871680835291849020805487019055925185815290927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a3610dfc565b73ffffffffffffffffffffffffffffffffffffffff8216611671576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f45524332303a206275726e2066726f6d20746865207a65726f2061646472657360448201527f7300000000000000000000000000000000000000000000000000000000000000606482015260840161080e565b73ffffffffffffffffffffffffffffffffffffffff821660009081526020819052604090205481811015611727576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f45524332303a206275726e20616d6f756e7420657863656564732062616c616e60448201527f6365000000000000000000000000000000000000000000000000000000000000606482015260840161080e565b73ffffffffffffffffffffffffffffffffffffffff83166000818152602081815260408083208686039055600280548790039055518581529192917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a3505050565b6000818152600183016020526040812054801561187b5760006117b6600183611bc3565b85549091506000906117ca90600190611bc3565b905081811461182f5760008660000182815481106117ea576117ea611bd6565b906000526020600020015490508087600001848154811061180d5761180d611bd6565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061184057611840611c05565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506106b1565b60009150506106b1565b6060816000018054806020026020016040519081016040528092919081815260200182805480156118d557602002820191906000526020600020905b8154815260200190600101908083116118c1575b50505050509050919050565b6000818152600183016020526040812054611928575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556106b1565b5060006106b1565b60006020828403121561194257600080fd5b81357fffffffff00000000000000000000000000000000000000000000000000000000811681146109c357600080fd5b600060208083528351808285015260005b8181101561199f57858101830151858201604001528201611983565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b803573ffffffffffffffffffffffffffffffffffffffff81168114611a0257600080fd5b919050565b60008060408385031215611a1a57600080fd5b611a23836119de565b946020939093013593505050565b600080600060608486031215611a4657600080fd5b611a4f846119de565b9250611a5d602085016119de565b9150604084013590509250925092565b600060208284031215611a7f57600080fd5b5035919050565b600060208284031215611a9857600080fd5b6109c3826119de565b6020808252825182820181905260009190848201906040850190845b81811015611aef57835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101611abd565b50909695505050505050565b60008060408385031215611b0e57600080fd5b611b17836119de565b9150611b25602084016119de565b90509250929050565b600181811c90821680611b4257607f821691505b602082108103611b7b577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b808201808211156106b1576106b1611b81565b818103818111156106b1576106b1611b81565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000813000a",
}

var BurnMintERC20ABI = BurnMintERC20MetaData.ABI

var BurnMintERC20Bin = BurnMintERC20MetaData.Bin

func DeployBurnMintERC20(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string, decimals_ uint8, maxSupply_ *big.Int, preMint *big.Int) (common.Address, *types.Transaction, *BurnMintERC20, error) {
	parsed, err := BurnMintERC20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BurnMintERC20Bin), backend, name, symbol, decimals_, maxSupply_, preMint)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BurnMintERC20{address: address, abi: *parsed, BurnMintERC20Caller: BurnMintERC20Caller{contract: contract}, BurnMintERC20Transactor: BurnMintERC20Transactor{contract: contract}, BurnMintERC20Filterer: BurnMintERC20Filterer{contract: contract}}, nil
}

type BurnMintERC20 struct {
	address common.Address
	abi     abi.ABI
	BurnMintERC20Caller
	BurnMintERC20Transactor
	BurnMintERC20Filterer
}

type BurnMintERC20Caller struct {
	contract *bind.BoundContract
}

type BurnMintERC20Transactor struct {
	contract *bind.BoundContract
}

type BurnMintERC20Filterer struct {
	contract *bind.BoundContract
}

type BurnMintERC20Session struct {
	Contract     *BurnMintERC20
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type BurnMintERC20CallerSession struct {
	Contract *BurnMintERC20Caller
	CallOpts bind.CallOpts
}

type BurnMintERC20TransactorSession struct {
	Contract     *BurnMintERC20Transactor
	TransactOpts bind.TransactOpts
}

type BurnMintERC20Raw struct {
	Contract *BurnMintERC20
}

type BurnMintERC20CallerRaw struct {
	Contract *BurnMintERC20Caller
}

type BurnMintERC20TransactorRaw struct {
	Contract *BurnMintERC20Transactor
}

func NewBurnMintERC20(address common.Address, backend bind.ContractBackend) (*BurnMintERC20, error) {
	abi, err := abi.JSON(strings.NewReader(BurnMintERC20ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindBurnMintERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20{address: address, abi: abi, BurnMintERC20Caller: BurnMintERC20Caller{contract: contract}, BurnMintERC20Transactor: BurnMintERC20Transactor{contract: contract}, BurnMintERC20Filterer: BurnMintERC20Filterer{contract: contract}}, nil
}

func NewBurnMintERC20Caller(address common.Address, caller bind.ContractCaller) (*BurnMintERC20Caller, error) {
	contract, err := bindBurnMintERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20Caller{contract: contract}, nil
}

func NewBurnMintERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*BurnMintERC20Transactor, error) {
	contract, err := bindBurnMintERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20Transactor{contract: contract}, nil
}

func NewBurnMintERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*BurnMintERC20Filterer, error) {
	contract, err := bindBurnMintERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20Filterer{contract: contract}, nil
}

func bindBurnMintERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BurnMintERC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_BurnMintERC20 *BurnMintERC20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnMintERC20.Contract.BurnMintERC20Caller.contract.Call(opts, result, method, params...)
}

func (_BurnMintERC20 *BurnMintERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.BurnMintERC20Transactor.contract.Transfer(opts)
}

func (_BurnMintERC20 *BurnMintERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.BurnMintERC20Transactor.contract.Transact(opts, method, params...)
}

func (_BurnMintERC20 *BurnMintERC20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnMintERC20.Contract.contract.Call(opts, result, method, params...)
}

func (_BurnMintERC20 *BurnMintERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.contract.Transfer(opts)
}

func (_BurnMintERC20 *BurnMintERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.contract.Transact(opts, method, params...)
}

func (_BurnMintERC20 *BurnMintERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _BurnMintERC20.Contract.Allowance(&_BurnMintERC20.CallOpts, owner, spender)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _BurnMintERC20.Contract.Allowance(&_BurnMintERC20.CallOpts, owner, spender)
}

func (_BurnMintERC20 *BurnMintERC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _BurnMintERC20.Contract.BalanceOf(&_BurnMintERC20.CallOpts, account)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _BurnMintERC20.Contract.BalanceOf(&_BurnMintERC20.CallOpts, account)
}

func (_BurnMintERC20 *BurnMintERC20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) Decimals() (uint8, error) {
	return _BurnMintERC20.Contract.Decimals(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) Decimals() (uint8, error) {
	return _BurnMintERC20.Contract.Decimals(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20Caller) GetBurners(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "getBurners")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) GetBurners() ([]common.Address, error) {
	return _BurnMintERC20.Contract.GetBurners(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) GetBurners() ([]common.Address, error) {
	return _BurnMintERC20.Contract.GetBurners(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20Caller) GetCCIPAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "getCCIPAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) GetCCIPAdmin() (common.Address, error) {
	return _BurnMintERC20.Contract.GetCCIPAdmin(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) GetCCIPAdmin() (common.Address, error) {
	return _BurnMintERC20.Contract.GetCCIPAdmin(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20Caller) GetMinters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "getMinters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) GetMinters() ([]common.Address, error) {
	return _BurnMintERC20.Contract.GetMinters(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) GetMinters() ([]common.Address, error) {
	return _BurnMintERC20.Contract.GetMinters(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20Caller) IsBurner(opts *bind.CallOpts, burner common.Address) (bool, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "isBurner", burner)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) IsBurner(burner common.Address) (bool, error) {
	return _BurnMintERC20.Contract.IsBurner(&_BurnMintERC20.CallOpts, burner)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) IsBurner(burner common.Address) (bool, error) {
	return _BurnMintERC20.Contract.IsBurner(&_BurnMintERC20.CallOpts, burner)
}

func (_BurnMintERC20 *BurnMintERC20Caller) IsMinter(opts *bind.CallOpts, minter common.Address) (bool, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "isMinter", minter)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) IsMinter(minter common.Address) (bool, error) {
	return _BurnMintERC20.Contract.IsMinter(&_BurnMintERC20.CallOpts, minter)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) IsMinter(minter common.Address) (bool, error) {
	return _BurnMintERC20.Contract.IsMinter(&_BurnMintERC20.CallOpts, minter)
}

func (_BurnMintERC20 *BurnMintERC20Caller) MaxSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "maxSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) MaxSupply() (*big.Int, error) {
	return _BurnMintERC20.Contract.MaxSupply(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) MaxSupply() (*big.Int, error) {
	return _BurnMintERC20.Contract.MaxSupply(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) Name() (string, error) {
	return _BurnMintERC20.Contract.Name(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) Name() (string, error) {
	return _BurnMintERC20.Contract.Name(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) Owner() (common.Address, error) {
	return _BurnMintERC20.Contract.Owner(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) Owner() (common.Address, error) {
	return _BurnMintERC20.Contract.Owner(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20Caller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BurnMintERC20.Contract.SupportsInterface(&_BurnMintERC20.CallOpts, interfaceId)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BurnMintERC20.Contract.SupportsInterface(&_BurnMintERC20.CallOpts, interfaceId)
}

func (_BurnMintERC20 *BurnMintERC20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) Symbol() (string, error) {
	return _BurnMintERC20.Contract.Symbol(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) Symbol() (string, error) {
	return _BurnMintERC20.Contract.Symbol(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BurnMintERC20.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BurnMintERC20 *BurnMintERC20Session) TotalSupply() (*big.Int, error) {
	return _BurnMintERC20.Contract.TotalSupply(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _BurnMintERC20.Contract.TotalSupply(&_BurnMintERC20.CallOpts)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "acceptOwnership")
}

func (_BurnMintERC20 *BurnMintERC20Session) AcceptOwnership() (*types.Transaction, error) {
	return _BurnMintERC20.Contract.AcceptOwnership(&_BurnMintERC20.TransactOpts)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _BurnMintERC20.Contract.AcceptOwnership(&_BurnMintERC20.TransactOpts)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "approve", spender, amount)
}

func (_BurnMintERC20 *BurnMintERC20Session) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.Approve(&_BurnMintERC20.TransactOpts, spender, amount)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.Approve(&_BurnMintERC20.TransactOpts, spender, amount)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) Burn(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "burn", amount)
}

func (_BurnMintERC20 *BurnMintERC20Session) Burn(amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.Burn(&_BurnMintERC20.TransactOpts, amount)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.Burn(&_BurnMintERC20.TransactOpts, amount)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) Burn0(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "burn0", account, amount)
}

func (_BurnMintERC20 *BurnMintERC20Session) Burn0(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.Burn0(&_BurnMintERC20.TransactOpts, account, amount)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) Burn0(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.Burn0(&_BurnMintERC20.TransactOpts, account, amount)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) BurnFrom(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "burnFrom", account, amount)
}

func (_BurnMintERC20 *BurnMintERC20Session) BurnFrom(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.BurnFrom(&_BurnMintERC20.TransactOpts, account, amount)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) BurnFrom(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.BurnFrom(&_BurnMintERC20.TransactOpts, account, amount)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

func (_BurnMintERC20 *BurnMintERC20Session) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.DecreaseAllowance(&_BurnMintERC20.TransactOpts, spender, subtractedValue)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.DecreaseAllowance(&_BurnMintERC20.TransactOpts, spender, subtractedValue)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) DecreaseApproval(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "decreaseApproval", spender, subtractedValue)
}

func (_BurnMintERC20 *BurnMintERC20Session) DecreaseApproval(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.DecreaseApproval(&_BurnMintERC20.TransactOpts, spender, subtractedValue)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) DecreaseApproval(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.DecreaseApproval(&_BurnMintERC20.TransactOpts, spender, subtractedValue)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) GrantBurnRole(opts *bind.TransactOpts, burner common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "grantBurnRole", burner)
}

func (_BurnMintERC20 *BurnMintERC20Session) GrantBurnRole(burner common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.GrantBurnRole(&_BurnMintERC20.TransactOpts, burner)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) GrantBurnRole(burner common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.GrantBurnRole(&_BurnMintERC20.TransactOpts, burner)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) GrantMintAndBurnRoles(opts *bind.TransactOpts, burnAndMinter common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "grantMintAndBurnRoles", burnAndMinter)
}

func (_BurnMintERC20 *BurnMintERC20Session) GrantMintAndBurnRoles(burnAndMinter common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.GrantMintAndBurnRoles(&_BurnMintERC20.TransactOpts, burnAndMinter)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) GrantMintAndBurnRoles(burnAndMinter common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.GrantMintAndBurnRoles(&_BurnMintERC20.TransactOpts, burnAndMinter)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) GrantMintRole(opts *bind.TransactOpts, minter common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "grantMintRole", minter)
}

func (_BurnMintERC20 *BurnMintERC20Session) GrantMintRole(minter common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.GrantMintRole(&_BurnMintERC20.TransactOpts, minter)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) GrantMintRole(minter common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.GrantMintRole(&_BurnMintERC20.TransactOpts, minter)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

func (_BurnMintERC20 *BurnMintERC20Session) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.IncreaseAllowance(&_BurnMintERC20.TransactOpts, spender, addedValue)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.IncreaseAllowance(&_BurnMintERC20.TransactOpts, spender, addedValue)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) IncreaseApproval(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "increaseApproval", spender, addedValue)
}

func (_BurnMintERC20 *BurnMintERC20Session) IncreaseApproval(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.IncreaseApproval(&_BurnMintERC20.TransactOpts, spender, addedValue)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) IncreaseApproval(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.IncreaseApproval(&_BurnMintERC20.TransactOpts, spender, addedValue)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) Mint(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "mint", account, amount)
}

func (_BurnMintERC20 *BurnMintERC20Session) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.Mint(&_BurnMintERC20.TransactOpts, account, amount)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.Mint(&_BurnMintERC20.TransactOpts, account, amount)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) RevokeBurnRole(opts *bind.TransactOpts, burner common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "revokeBurnRole", burner)
}

func (_BurnMintERC20 *BurnMintERC20Session) RevokeBurnRole(burner common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.RevokeBurnRole(&_BurnMintERC20.TransactOpts, burner)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) RevokeBurnRole(burner common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.RevokeBurnRole(&_BurnMintERC20.TransactOpts, burner)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) RevokeMintRole(opts *bind.TransactOpts, minter common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "revokeMintRole", minter)
}

func (_BurnMintERC20 *BurnMintERC20Session) RevokeMintRole(minter common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.RevokeMintRole(&_BurnMintERC20.TransactOpts, minter)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) RevokeMintRole(minter common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.RevokeMintRole(&_BurnMintERC20.TransactOpts, minter)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) SetCCIPAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "setCCIPAdmin", newAdmin)
}

func (_BurnMintERC20 *BurnMintERC20Session) SetCCIPAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.SetCCIPAdmin(&_BurnMintERC20.TransactOpts, newAdmin)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) SetCCIPAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.SetCCIPAdmin(&_BurnMintERC20.TransactOpts, newAdmin)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "transfer", to, amount)
}

func (_BurnMintERC20 *BurnMintERC20Session) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.Transfer(&_BurnMintERC20.TransactOpts, to, amount)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.Transfer(&_BurnMintERC20.TransactOpts, to, amount)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "transferFrom", from, to, amount)
}

func (_BurnMintERC20 *BurnMintERC20Session) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.TransferFrom(&_BurnMintERC20.TransactOpts, from, to, amount)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.TransferFrom(&_BurnMintERC20.TransactOpts, from, to, amount)
}

func (_BurnMintERC20 *BurnMintERC20Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.contract.Transact(opts, "transferOwnership", to)
}

func (_BurnMintERC20 *BurnMintERC20Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.TransferOwnership(&_BurnMintERC20.TransactOpts, to)
}

func (_BurnMintERC20 *BurnMintERC20TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BurnMintERC20.Contract.TransferOwnership(&_BurnMintERC20.TransactOpts, to)
}

type BurnMintERC20ApprovalIterator struct {
	Event *BurnMintERC20Approval

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC20ApprovalIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC20Approval)
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
		it.Event = new(BurnMintERC20Approval)
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

func (it *BurnMintERC20ApprovalIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log
}

func (_BurnMintERC20 *BurnMintERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*BurnMintERC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _BurnMintERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20ApprovalIterator{contract: _BurnMintERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

func (_BurnMintERC20 *BurnMintERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *BurnMintERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _BurnMintERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC20Approval)
				if err := _BurnMintERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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

func (_BurnMintERC20 *BurnMintERC20Filterer) ParseApproval(log types.Log) (*BurnMintERC20Approval, error) {
	event := new(BurnMintERC20Approval)
	if err := _BurnMintERC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintERC20BurnAccessGrantedIterator struct {
	Event *BurnMintERC20BurnAccessGranted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC20BurnAccessGrantedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC20BurnAccessGranted)
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
		it.Event = new(BurnMintERC20BurnAccessGranted)
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

func (it *BurnMintERC20BurnAccessGrantedIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC20BurnAccessGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC20BurnAccessGranted struct {
	Burner common.Address
	Raw    types.Log
}

func (_BurnMintERC20 *BurnMintERC20Filterer) FilterBurnAccessGranted(opts *bind.FilterOpts) (*BurnMintERC20BurnAccessGrantedIterator, error) {

	logs, sub, err := _BurnMintERC20.contract.FilterLogs(opts, "BurnAccessGranted")
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20BurnAccessGrantedIterator{contract: _BurnMintERC20.contract, event: "BurnAccessGranted", logs: logs, sub: sub}, nil
}

func (_BurnMintERC20 *BurnMintERC20Filterer) WatchBurnAccessGranted(opts *bind.WatchOpts, sink chan<- *BurnMintERC20BurnAccessGranted) (event.Subscription, error) {

	logs, sub, err := _BurnMintERC20.contract.WatchLogs(opts, "BurnAccessGranted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC20BurnAccessGranted)
				if err := _BurnMintERC20.contract.UnpackLog(event, "BurnAccessGranted", log); err != nil {
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

func (_BurnMintERC20 *BurnMintERC20Filterer) ParseBurnAccessGranted(log types.Log) (*BurnMintERC20BurnAccessGranted, error) {
	event := new(BurnMintERC20BurnAccessGranted)
	if err := _BurnMintERC20.contract.UnpackLog(event, "BurnAccessGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintERC20BurnAccessRevokedIterator struct {
	Event *BurnMintERC20BurnAccessRevoked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC20BurnAccessRevokedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC20BurnAccessRevoked)
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
		it.Event = new(BurnMintERC20BurnAccessRevoked)
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

func (it *BurnMintERC20BurnAccessRevokedIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC20BurnAccessRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC20BurnAccessRevoked struct {
	Burner common.Address
	Raw    types.Log
}

func (_BurnMintERC20 *BurnMintERC20Filterer) FilterBurnAccessRevoked(opts *bind.FilterOpts) (*BurnMintERC20BurnAccessRevokedIterator, error) {

	logs, sub, err := _BurnMintERC20.contract.FilterLogs(opts, "BurnAccessRevoked")
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20BurnAccessRevokedIterator{contract: _BurnMintERC20.contract, event: "BurnAccessRevoked", logs: logs, sub: sub}, nil
}

func (_BurnMintERC20 *BurnMintERC20Filterer) WatchBurnAccessRevoked(opts *bind.WatchOpts, sink chan<- *BurnMintERC20BurnAccessRevoked) (event.Subscription, error) {

	logs, sub, err := _BurnMintERC20.contract.WatchLogs(opts, "BurnAccessRevoked")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC20BurnAccessRevoked)
				if err := _BurnMintERC20.contract.UnpackLog(event, "BurnAccessRevoked", log); err != nil {
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

func (_BurnMintERC20 *BurnMintERC20Filterer) ParseBurnAccessRevoked(log types.Log) (*BurnMintERC20BurnAccessRevoked, error) {
	event := new(BurnMintERC20BurnAccessRevoked)
	if err := _BurnMintERC20.contract.UnpackLog(event, "BurnAccessRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintERC20CCIPAdminTransferredIterator struct {
	Event *BurnMintERC20CCIPAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC20CCIPAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC20CCIPAdminTransferred)
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
		it.Event = new(BurnMintERC20CCIPAdminTransferred)
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

func (it *BurnMintERC20CCIPAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC20CCIPAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC20CCIPAdminTransferred struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log
}

func (_BurnMintERC20 *BurnMintERC20Filterer) FilterCCIPAdminTransferred(opts *bind.FilterOpts, previousAdmin []common.Address, newAdmin []common.Address) (*BurnMintERC20CCIPAdminTransferredIterator, error) {

	var previousAdminRule []interface{}
	for _, previousAdminItem := range previousAdmin {
		previousAdminRule = append(previousAdminRule, previousAdminItem)
	}
	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _BurnMintERC20.contract.FilterLogs(opts, "CCIPAdminTransferred", previousAdminRule, newAdminRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20CCIPAdminTransferredIterator{contract: _BurnMintERC20.contract, event: "CCIPAdminTransferred", logs: logs, sub: sub}, nil
}

func (_BurnMintERC20 *BurnMintERC20Filterer) WatchCCIPAdminTransferred(opts *bind.WatchOpts, sink chan<- *BurnMintERC20CCIPAdminTransferred, previousAdmin []common.Address, newAdmin []common.Address) (event.Subscription, error) {

	var previousAdminRule []interface{}
	for _, previousAdminItem := range previousAdmin {
		previousAdminRule = append(previousAdminRule, previousAdminItem)
	}
	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _BurnMintERC20.contract.WatchLogs(opts, "CCIPAdminTransferred", previousAdminRule, newAdminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC20CCIPAdminTransferred)
				if err := _BurnMintERC20.contract.UnpackLog(event, "CCIPAdminTransferred", log); err != nil {
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

func (_BurnMintERC20 *BurnMintERC20Filterer) ParseCCIPAdminTransferred(log types.Log) (*BurnMintERC20CCIPAdminTransferred, error) {
	event := new(BurnMintERC20CCIPAdminTransferred)
	if err := _BurnMintERC20.contract.UnpackLog(event, "CCIPAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintERC20MintAccessGrantedIterator struct {
	Event *BurnMintERC20MintAccessGranted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC20MintAccessGrantedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC20MintAccessGranted)
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
		it.Event = new(BurnMintERC20MintAccessGranted)
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

func (it *BurnMintERC20MintAccessGrantedIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC20MintAccessGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC20MintAccessGranted struct {
	Minter common.Address
	Raw    types.Log
}

func (_BurnMintERC20 *BurnMintERC20Filterer) FilterMintAccessGranted(opts *bind.FilterOpts) (*BurnMintERC20MintAccessGrantedIterator, error) {

	logs, sub, err := _BurnMintERC20.contract.FilterLogs(opts, "MintAccessGranted")
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20MintAccessGrantedIterator{contract: _BurnMintERC20.contract, event: "MintAccessGranted", logs: logs, sub: sub}, nil
}

func (_BurnMintERC20 *BurnMintERC20Filterer) WatchMintAccessGranted(opts *bind.WatchOpts, sink chan<- *BurnMintERC20MintAccessGranted) (event.Subscription, error) {

	logs, sub, err := _BurnMintERC20.contract.WatchLogs(opts, "MintAccessGranted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC20MintAccessGranted)
				if err := _BurnMintERC20.contract.UnpackLog(event, "MintAccessGranted", log); err != nil {
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

func (_BurnMintERC20 *BurnMintERC20Filterer) ParseMintAccessGranted(log types.Log) (*BurnMintERC20MintAccessGranted, error) {
	event := new(BurnMintERC20MintAccessGranted)
	if err := _BurnMintERC20.contract.UnpackLog(event, "MintAccessGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintERC20MintAccessRevokedIterator struct {
	Event *BurnMintERC20MintAccessRevoked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC20MintAccessRevokedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC20MintAccessRevoked)
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
		it.Event = new(BurnMintERC20MintAccessRevoked)
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

func (it *BurnMintERC20MintAccessRevokedIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC20MintAccessRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC20MintAccessRevoked struct {
	Minter common.Address
	Raw    types.Log
}

func (_BurnMintERC20 *BurnMintERC20Filterer) FilterMintAccessRevoked(opts *bind.FilterOpts) (*BurnMintERC20MintAccessRevokedIterator, error) {

	logs, sub, err := _BurnMintERC20.contract.FilterLogs(opts, "MintAccessRevoked")
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20MintAccessRevokedIterator{contract: _BurnMintERC20.contract, event: "MintAccessRevoked", logs: logs, sub: sub}, nil
}

func (_BurnMintERC20 *BurnMintERC20Filterer) WatchMintAccessRevoked(opts *bind.WatchOpts, sink chan<- *BurnMintERC20MintAccessRevoked) (event.Subscription, error) {

	logs, sub, err := _BurnMintERC20.contract.WatchLogs(opts, "MintAccessRevoked")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC20MintAccessRevoked)
				if err := _BurnMintERC20.contract.UnpackLog(event, "MintAccessRevoked", log); err != nil {
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

func (_BurnMintERC20 *BurnMintERC20Filterer) ParseMintAccessRevoked(log types.Log) (*BurnMintERC20MintAccessRevoked, error) {
	event := new(BurnMintERC20MintAccessRevoked)
	if err := _BurnMintERC20.contract.UnpackLog(event, "MintAccessRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintERC20OwnershipTransferRequestedIterator struct {
	Event *BurnMintERC20OwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC20OwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC20OwnershipTransferRequested)
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
		it.Event = new(BurnMintERC20OwnershipTransferRequested)
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

func (it *BurnMintERC20OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC20OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC20OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BurnMintERC20 *BurnMintERC20Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintERC20OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintERC20.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20OwnershipTransferRequestedIterator{contract: _BurnMintERC20.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_BurnMintERC20 *BurnMintERC20Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BurnMintERC20OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintERC20.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC20OwnershipTransferRequested)
				if err := _BurnMintERC20.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_BurnMintERC20 *BurnMintERC20Filterer) ParseOwnershipTransferRequested(log types.Log) (*BurnMintERC20OwnershipTransferRequested, error) {
	event := new(BurnMintERC20OwnershipTransferRequested)
	if err := _BurnMintERC20.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintERC20OwnershipTransferredIterator struct {
	Event *BurnMintERC20OwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC20OwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC20OwnershipTransferred)
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
		it.Event = new(BurnMintERC20OwnershipTransferred)
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

func (it *BurnMintERC20OwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC20OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC20OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BurnMintERC20 *BurnMintERC20Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintERC20OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintERC20.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20OwnershipTransferredIterator{contract: _BurnMintERC20.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_BurnMintERC20 *BurnMintERC20Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BurnMintERC20OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintERC20.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC20OwnershipTransferred)
				if err := _BurnMintERC20.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_BurnMintERC20 *BurnMintERC20Filterer) ParseOwnershipTransferred(log types.Log) (*BurnMintERC20OwnershipTransferred, error) {
	event := new(BurnMintERC20OwnershipTransferred)
	if err := _BurnMintERC20.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintERC20TransferIterator struct {
	Event *BurnMintERC20Transfer

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintERC20TransferIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintERC20Transfer)
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
		it.Event = new(BurnMintERC20Transfer)
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

func (it *BurnMintERC20TransferIterator) Error() error {
	return it.fail
}

func (it *BurnMintERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log
}

func (_BurnMintERC20 *BurnMintERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintERC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintERC20TransferIterator{contract: _BurnMintERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

func (_BurnMintERC20 *BurnMintERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *BurnMintERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintERC20Transfer)
				if err := _BurnMintERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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

func (_BurnMintERC20 *BurnMintERC20Filterer) ParseTransfer(log types.Log) (*BurnMintERC20Transfer, error) {
	event := new(BurnMintERC20Transfer)
	if err := _BurnMintERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_BurnMintERC20 *BurnMintERC20) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _BurnMintERC20.abi.Events["Approval"].ID:
		return _BurnMintERC20.ParseApproval(log)
	case _BurnMintERC20.abi.Events["BurnAccessGranted"].ID:
		return _BurnMintERC20.ParseBurnAccessGranted(log)
	case _BurnMintERC20.abi.Events["BurnAccessRevoked"].ID:
		return _BurnMintERC20.ParseBurnAccessRevoked(log)
	case _BurnMintERC20.abi.Events["CCIPAdminTransferred"].ID:
		return _BurnMintERC20.ParseCCIPAdminTransferred(log)
	case _BurnMintERC20.abi.Events["MintAccessGranted"].ID:
		return _BurnMintERC20.ParseMintAccessGranted(log)
	case _BurnMintERC20.abi.Events["MintAccessRevoked"].ID:
		return _BurnMintERC20.ParseMintAccessRevoked(log)
	case _BurnMintERC20.abi.Events["OwnershipTransferRequested"].ID:
		return _BurnMintERC20.ParseOwnershipTransferRequested(log)
	case _BurnMintERC20.abi.Events["OwnershipTransferred"].ID:
		return _BurnMintERC20.ParseOwnershipTransferred(log)
	case _BurnMintERC20.abi.Events["Transfer"].ID:
		return _BurnMintERC20.ParseTransfer(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (BurnMintERC20Approval) Topic() common.Hash {
	return common.HexToHash("0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925")
}

func (BurnMintERC20BurnAccessGranted) Topic() common.Hash {
	return common.HexToHash("0x92308bb7573b2a3d17ddb868b39d8ebec433f3194421abc22d084f89658c9bad")
}

func (BurnMintERC20BurnAccessRevoked) Topic() common.Hash {
	return common.HexToHash("0x0a675452746933cefe3d74182e78db7afe57ba60eaa4234b5d85e9aa41b0610c")
}

func (BurnMintERC20CCIPAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x9524c9e4b0b61eb018dd58a1cd856e3e74009528328ab4a613b434fa631d7242")
}

func (BurnMintERC20MintAccessGranted) Topic() common.Hash {
	return common.HexToHash("0xe46fef8bbff1389d9010703cf8ebb363fb3daf5bf56edc27080b67bc8d9251ea")
}

func (BurnMintERC20MintAccessRevoked) Topic() common.Hash {
	return common.HexToHash("0xed998b960f6340d045f620c119730f7aa7995e7425c2401d3a5b64ff998a59e9")
}

func (BurnMintERC20OwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (BurnMintERC20OwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (BurnMintERC20Transfer) Topic() common.Hash {
	return common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
}

func (_BurnMintERC20 *BurnMintERC20) Address() common.Address {
	return _BurnMintERC20.address
}

type BurnMintERC20Interface interface {
	Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error)

	BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error)

	Decimals(opts *bind.CallOpts) (uint8, error)

	GetBurners(opts *bind.CallOpts) ([]common.Address, error)

	GetCCIPAdmin(opts *bind.CallOpts) (common.Address, error)

	GetMinters(opts *bind.CallOpts) ([]common.Address, error)

	IsBurner(opts *bind.CallOpts, burner common.Address) (bool, error)

	IsMinter(opts *bind.CallOpts, minter common.Address) (bool, error)

	MaxSupply(opts *bind.CallOpts) (*big.Int, error)

	Name(opts *bind.CallOpts) (string, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	Symbol(opts *bind.CallOpts) (string, error)

	TotalSupply(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error)

	Burn(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	Burn0(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error)

	BurnFrom(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error)

	DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error)

	DecreaseApproval(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error)

	GrantBurnRole(opts *bind.TransactOpts, burner common.Address) (*types.Transaction, error)

	GrantMintAndBurnRoles(opts *bind.TransactOpts, burnAndMinter common.Address) (*types.Transaction, error)

	GrantMintRole(opts *bind.TransactOpts, minter common.Address) (*types.Transaction, error)

	IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error)

	IncreaseApproval(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error)

	Mint(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error)

	RevokeBurnRole(opts *bind.TransactOpts, burner common.Address) (*types.Transaction, error)

	RevokeMintRole(opts *bind.TransactOpts, minter common.Address) (*types.Transaction, error)

	SetCCIPAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error)

	Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error)

	TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*BurnMintERC20ApprovalIterator, error)

	WatchApproval(opts *bind.WatchOpts, sink chan<- *BurnMintERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error)

	ParseApproval(log types.Log) (*BurnMintERC20Approval, error)

	FilterBurnAccessGranted(opts *bind.FilterOpts) (*BurnMintERC20BurnAccessGrantedIterator, error)

	WatchBurnAccessGranted(opts *bind.WatchOpts, sink chan<- *BurnMintERC20BurnAccessGranted) (event.Subscription, error)

	ParseBurnAccessGranted(log types.Log) (*BurnMintERC20BurnAccessGranted, error)

	FilterBurnAccessRevoked(opts *bind.FilterOpts) (*BurnMintERC20BurnAccessRevokedIterator, error)

	WatchBurnAccessRevoked(opts *bind.WatchOpts, sink chan<- *BurnMintERC20BurnAccessRevoked) (event.Subscription, error)

	ParseBurnAccessRevoked(log types.Log) (*BurnMintERC20BurnAccessRevoked, error)

	FilterCCIPAdminTransferred(opts *bind.FilterOpts, previousAdmin []common.Address, newAdmin []common.Address) (*BurnMintERC20CCIPAdminTransferredIterator, error)

	WatchCCIPAdminTransferred(opts *bind.WatchOpts, sink chan<- *BurnMintERC20CCIPAdminTransferred, previousAdmin []common.Address, newAdmin []common.Address) (event.Subscription, error)

	ParseCCIPAdminTransferred(log types.Log) (*BurnMintERC20CCIPAdminTransferred, error)

	FilterMintAccessGranted(opts *bind.FilterOpts) (*BurnMintERC20MintAccessGrantedIterator, error)

	WatchMintAccessGranted(opts *bind.WatchOpts, sink chan<- *BurnMintERC20MintAccessGranted) (event.Subscription, error)

	ParseMintAccessGranted(log types.Log) (*BurnMintERC20MintAccessGranted, error)

	FilterMintAccessRevoked(opts *bind.FilterOpts) (*BurnMintERC20MintAccessRevokedIterator, error)

	WatchMintAccessRevoked(opts *bind.WatchOpts, sink chan<- *BurnMintERC20MintAccessRevoked) (event.Subscription, error)

	ParseMintAccessRevoked(log types.Log) (*BurnMintERC20MintAccessRevoked, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintERC20OwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BurnMintERC20OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*BurnMintERC20OwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintERC20OwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BurnMintERC20OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*BurnMintERC20OwnershipTransferred, error)

	FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintERC20TransferIterator, error)

	WatchTransfer(opts *bind.WatchOpts, sink chan<- *BurnMintERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseTransfer(log types.Log) (*BurnMintERC20Transfer, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
