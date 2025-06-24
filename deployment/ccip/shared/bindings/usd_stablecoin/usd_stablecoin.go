// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package usd_stablecoin

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
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated"
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

var StablecoinMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"DOMAIN_SEPARATOR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"burn\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"drain\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"eip712Domain\",\"inputs\":[],\"outputs\":[{\"name\":\"fields\",\"type\":\"bytes1\",\"internalType\":\"bytes1\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"verifyingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"extensions\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"freeze\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"frozen\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_symbol\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nonces\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"permit\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"recoverERC20\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"test\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unfreeze\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Burn\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EIP712DomainChanged\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Freeze\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Mint\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unfreeze\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureLength\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureS\",\"inputs\":[{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC2612ExpiredSignature\",\"inputs\":[{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC2612InvalidSigner\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAccountNonce\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"currentNonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x60808060405234601557613273908161001b8239f35b600080fdfe6080604052600436101561001257600080fd5b60003560e01c806306fdde0314612494578063095ea7b3146124505780631171bda91461239d57806318160ddd1461234257806323b872dd14612144578063313ce5671461210a5780633644e515146120c95780633f4ba83a14611fcd57806342966c6814611f5057806345c8b1a614611e9f5780634cd88b76146112305780635c975abb146111d057806370a082311461114c578063715018a6146110bf57806379ba50971461101b5780637ecebe0014610f975780638456cb5914610ec357806384b0196e14610b035780638d1fdf2f14610a4f5780638da5cb5b146109de57806395d89b411461087f578063a0712d68146106f1578063a9059cbb146106a2578063d051665014610638578063d505accf14610464578063dd62ed3e146103af578063e30c39781461033e578063ece53132146102b8578063f2fde38b1461019c5763f8a8fd6d1461016657600080fd5b346101975760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261019757005b600080fd5b346101975760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101975773ffffffffffffffffffffffffffffffffffffffff6101e86125f7565b6101f0612a04565b16807fffffffffffffffffffffffff00000000000000000000000000000000000000007f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c005416177f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c005573ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930054167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700600080a3005b346101975760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101975761033c6102f26125f7565b73ffffffffffffffffffffffffffffffffffffffff81166000527f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0060205260406000205490612a44565b005b346101975760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261019757602073ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c005416604051908152f35b346101975760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610197576103e66125f7565b73ffffffffffffffffffffffffffffffffffffffff61044c61040661261a565b9273ffffffffffffffffffffffffffffffffffffffff166000527f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace01602052604060002090565b91166000526020526020604060002054604051908152f35b346101975760e07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101975761049b6125f7565b6104a361261a565b604435906064359260843560ff811681036101975784421161060a576105c56105bc73ffffffffffffffffffffffffffffffffffffffff9283851697886000527f5ab42ced628888259c08ac98db1eb0cf702fc1501344311d8b100cd1bfe4bb006020526040600020908154916001830190556040519060208201927f6e71edae12b1b97f4d1f60370fef10105fa2faae0126114a169c64845d6126c984528b6040840152878a1660608401528a608084015260a083015260c082015260c0815261056f60e0826126af565b51902061057a612b71565b90604051917f190100000000000000000000000000000000000000000000000000000000000083526002830152602282015260c43591604260a4359220612d96565b90929192612e2c565b168481036105d8575061033c93506127b8565b84907f4b800e460000000000000000000000000000000000000000000000000000000060005260045260245260446000fd5b847f627913020000000000000000000000000000000000000000000000000000000060005260045260246000fd5b346101975760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101975773ffffffffffffffffffffffffffffffffffffffff6106846125f7565b166000526000602052602060ff604060002054166040519015158152f35b346101975760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610197576106e66106dc6125f7565b6024359033612891565b602060405160018152f35b346101975760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101975760043561072b612a04565b331561085057610739612d41565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0254818101809111610821577f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0255336000527f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace00602052604060002081815401905560405181815260007fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef60203393a360405190815233907fab8530f87dc9b59234c4623bf917212bb2536d647574c8e7e5da92c2ede0c9f860203392a3602060405160018152f35b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7fec442f0500000000000000000000000000000000000000000000000000000000600052600060045260246000fd5b346101975760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101975760405160007f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace04546108de81612765565b808452906001811690811561099c575060011461091e575b61091a83610906818503826126af565b604051918291602083526020830190612598565b0390f35b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0460009081527f46a2803e59a4de4e7a4c574b1243f25977ac4c77d5a1a4a609b5394cebb4a2aa939250905b808210610982575090915081016020016109066108f6565b91926001816020925483858801015201910190929161096a565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660208086019190915291151560051b8401909101915061090690506108f6565b346101975760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261019757602073ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005416604051908152f35b346101975760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101975773ffffffffffffffffffffffffffffffffffffffff610a9b6125f7565b610aa3612a04565b16806000526000602052604060002060017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00825416179055337f51d18786e9cb144f87d46e7b796309ea84c7c687d91e09c97f051eacf59bc528600080a3005b346101975760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610197577fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d100541580610e9a575b15610e3c576040517fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d10254816000610b9183612765565b8083529260018116908115610dff5750600114610d80575b610bb5925003826126af565b6040517fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d10354816000610be683612765565b8083529260018116908115610d435750600114610cc4575b610c1191925092610c66949303826126af565b6020610c7460405192610c2483856126af565b6000845260003681376040519586957f0f00000000000000000000000000000000000000000000000000000000000000875260e08588015260e0870190612598565b908582036040870152612598565b466060850152306080850152600060a085015283810360c085015281808451928381520193019160005b828110610cad57505050500390f35b835185528695509381019392810192600101610c9e565b507fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d103600090815290917f5f9ce34815f8e11431c7bb75a8e6886a91478f7ffc1dbb0a98dc240fddd76b755b818310610d27575050906020610c1192820101610bfe565b6020919350806001915483858801015201910190918392610d0f565b60209250610c119491507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001682840152151560051b820101610bfe565b507fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d102600090815290917f42ad5d3e1f2e6e70edcf6d991b8a3023d3fca8047a131592f9edb9fd9b89d57d5b818310610de3575050906020610bb592820101610ba9565b6020919350806001915483858801015201910190918392610dcb565b60209250610bb59491507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001682840152151560051b820101610ba9565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f4549503731323a20556e696e697469616c697a656400000000000000000000006044820152fd5b507fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d1015415610b5b565b346101975760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261019757610efa612a04565b610f02612d41565b60017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff007fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005416177fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2586020604051338152a1005b346101975760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101975773ffffffffffffffffffffffffffffffffffffffff610fe36125f7565b166000527f5ab42ced628888259c08ac98db1eb0cf702fc1501344311d8b100cd1bfe4bb006020526020604060002054604051908152f35b346101975760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610197573373ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c005416036110915761033c33612c31565b7f118cdaa7000000000000000000000000000000000000000000000000000000006000523360045260246000fd5b346101975760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101975760646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f556e737570706f727465640000000000000000000000000000000000000000006044820152fd5b346101975760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101975773ffffffffffffffffffffffffffffffffffffffff6111986125f7565b166000527f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace006020526020604060002054604051908152f35b346101975760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261019757602060ff7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330054166040519015158152f35b346101975760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101975760043567ffffffffffffffff81116101975761127f9036906004016126f0565b60243567ffffffffffffffff81116101975761129f9036906004016126f0565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00549160ff8360401c16159267ffffffffffffffff811680159081611e97575b6001149081611e8d575b159081611e84575b50611e5a578360017fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000008316177ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0055611e05575b5061134d612bd8565b611355612bd8565b61135d612bd8565b80519267ffffffffffffffff8411611987576113997f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0354612765565b601f8111611d81575b50602093601f8111600114611ca1578091929394600091611c96575b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c1916177f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace03555b825167ffffffffffffffff81116119875761144a7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0454612765565b601f8111611c12575b506020601f8211600114611b315781929394600092611b26575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c1916177f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace04555b6114c7612bd8565b60408051916114d682846126af565b600183527f31000000000000000000000000000000000000000000000000000000000000006020840152611508612bd8565b835167ffffffffffffffff8111611987576115437fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d10254612765565b601f8111611aa2575b50602094601f82116001146119c1579481929394956000926119b6575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c1916177fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d102555b825167ffffffffffffffff8111611987576115f67fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d10354612765565b601f8111611903575b506020601f82116001146118225781929394600092611817575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c1916177fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d103555b60007fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d1005560007fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d101556116bb612bd8565b6116c3612bd8565b6116cb612bd8565b33156117e8576116da33612c31565b6116e2612bd8565b6116ea612bd8565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff007fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330054167fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005561175657005b60207fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2917fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0054167ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00555160018152a1005b7f1e4fbdf700000000000000000000000000000000000000000000000000000000600052600060045260246000fd5b015190508480611619565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08216907fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d103600052806000209160005b8181106118eb575095836001959697106118b4575b505050811b017fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d1035561166b565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c19169055848080611887565b9192602060018192868b015181550194019201611872565b7fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d1036000527f5f9ce34815f8e11431c7bb75a8e6886a91478f7ffc1dbb0a98dc240fddd76b75601f830160051c8101916020841061197d575b601f0160051c01905b81811061197157506115ff565b60008155600101611964565b909150819061195b565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b015190508580611569565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08216957fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d102600052806000209160005b888110611a8a57508360019596979810611a53575b505050811b017fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d102556115bb565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c19169055858080611a26565b91926020600181928685015181550194019201611a11565b7fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d1026000527f42ad5d3e1f2e6e70edcf6d991b8a3023d3fca8047a131592f9edb9fd9b89d57d601f830160051c81019160208410611b1c575b601f0160051c01905b818110611b10575061154c565b60008155600101611b03565b9091508190611afa565b01519050848061146d565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08216907f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace04600052806000209160005b818110611bfa57509583600195969710611bc3575b505050811b017f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace04556114bf565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c19169055848080611b96565b9192602060018192868b015181550194019201611b81565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace046000527f46a2803e59a4de4e7a4c574b1243f25977ac4c77d5a1a4a609b5394cebb4a2aa601f830160051c81019160208410611c8c575b601f0160051c01905b818110611c805750611453565b60008155600101611c73565b9091508190611c6a565b9050830151856113be565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08116947f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace03600052806000209060005b878110611d6957508260019495969710611d32575b5050811b017f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace035561140f565b8501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558580611d06565b90916020600181928589015181550193019101611cf1565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace036000527f2ae08a8e29253f69ac5d979a101956ab8f8d9d7ded63fa7a83b16fc47648eab0601f860160051c81019160208710611dfb575b601f0160051c01905b818110611def57506113a2565b60008155600101611de2565b9091508190611dd9565b7fffffffffffffffffffffffffffffffffffffffffffffff0000000000000000001668010000000000000001177ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005583611344565b7ff92ee8a90000000000000000000000000000000000000000000000000000000060005260046000fd5b905015856112f1565b303b1591506112e9565b8591506112df565b346101975760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101975773ffffffffffffffffffffffffffffffffffffffff611eeb6125f7565b611ef3612a04565b1680600052600060205260406000207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008154169055337f4f3ab9ff0cc4f039268532098e01239544b0420171876e36889d01c62c784c79600080a3005b346101975760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261019757600435611f8a612a04565b611f948133612a44565b60405190815233907fbac40739b0d4ca32fa2d82fc91630465ba3eddd1598da6fca393b26fb63b945360203392a3602060405160018152f35b346101975760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261019757612004612a04565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff81161561209f577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00167fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa6020604051338152a1005b7f8dfc202b0000000000000000000000000000000000000000000000000000000060005260046000fd5b346101975760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610197576020612102612b71565b604051908152f35b346101975760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261019757602060405160128152f35b34610197576121523661263d565b9061219c8373ffffffffffffffffffffffffffffffffffffffff166000527f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace01602052604060002090565b73ffffffffffffffffffffffffffffffffffffffff3316600052602052604060002054927fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff84106121f2575b6106e69350612891565b82841061230c57612201612d41565b61220a81612f18565b61221333612f18565b73ffffffffffffffffffffffffffffffffffffffff8116156122dd5733156122ae576106e6936122828273ffffffffffffffffffffffffffffffffffffffff166000527f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace01602052604060002090565b73ffffffffffffffffffffffffffffffffffffffff3316600052602052836040600020910390556121e8565b7f94280d6200000000000000000000000000000000000000000000000000000000600052600060045260246000fd5b7fe602df0500000000000000000000000000000000000000000000000000000000600052600060045260246000fd5b82847ffb8f41b2000000000000000000000000000000000000000000000000000000006000523360045260245260445260646000fd5b346101975760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101975760207f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0254604051908152f35b34610197576000602073ffffffffffffffffffffffffffffffffffffffff6044816123c73661263d565b6040517fa9059cbb0000000000000000000000000000000000000000000000000000000081529590911660048601526024850152929485938492165af180156124445761241057005b6020813d60201161243c575b81612429602093836126af565b8101031261019757518015150361019757005b3d915061241c565b6040513d6000823e3d90fd5b346101975760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610197576106e661248a6125f7565b60243590336127b8565b346101975760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101975760405160007f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace03546124f381612765565b808452906001811690811561099c575060011461251a5761091a83610906818503826126af565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0360009081527f2ae08a8e29253f69ac5d979a101956ab8f8d9d7ded63fa7a83b16fc47648eab0939250905b80821061257e575090915081016020016109066108f6565b919260018160209254838588010152019101909291612566565b919082519283825260005b8481106125e25750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8460006020809697860101520116010190565b806020809284010151828286010152016125a3565b6004359073ffffffffffffffffffffffffffffffffffffffff8216820361019757565b6024359073ffffffffffffffffffffffffffffffffffffffff8216820361019757565b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc60609101126101975760043573ffffffffffffffffffffffffffffffffffffffff81168103610197579060243573ffffffffffffffffffffffffffffffffffffffff81168103610197579060443590565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff82111761198757604052565b81601f820112156101975780359067ffffffffffffffff82116119875760405192612743601f84017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016602001856126af565b8284526020838301011161019757816000926020809301838601378301015290565b90600182811c921680156127ae575b602083101461277f57565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b91607f1691612774565b916127c1612d41565b6127ca83612f18565b6127d382612f18565b73ffffffffffffffffffffffffffffffffffffffff83169182156122dd5773ffffffffffffffffffffffffffffffffffffffff169283156122ae577f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259161287b60209273ffffffffffffffffffffffffffffffffffffffff166000527f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace01602052604060002090565b85600052825280604060002055604051908152a3565b909173ffffffffffffffffffffffffffffffffffffffff82169182156129d55773ffffffffffffffffffffffffffffffffffffffff8416938415610850576128de6128e3926128de612d41565b612f18565b60008281527f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0060205260408120548281106129a25791604082827fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9587602096527f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace00865203828220558681527f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace00845220818154019055604051908152a3565b6064937fe450d38c0000000000000000000000000000000000000000000000000000000083949352600452602452604452fd5b7f96c6fd1e00000000000000000000000000000000000000000000000000000000600052600060045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005416330361109157565b73ffffffffffffffffffffffffffffffffffffffff1680156129d557612a68612d41565b6000918183527f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace006020526040832054818110612b3f57817fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef926020928587527f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace008452036040862055807f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0254037f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0255604051908152a3565b83927fe450d38c0000000000000000000000000000000000000000000000000000000060649552600452602452604452fd5b612b79612fa3565b612b81613117565b6040519060208201927f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f8452604083015260608201524660808201523060a082015260a08152612bd260c0826126af565b51902090565b60ff7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005460401c1615612c0757565b7fd7e6bcf80000000000000000000000000000000000000000000000000000000060005260046000fd5b7fffffffffffffffffffffffff00000000000000000000000000000000000000007f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c0054167f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c005573ffffffffffffffffffffffffffffffffffffffff807f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930054921691827fffffffffffffffffffffffff00000000000000000000000000000000000000008216177f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930055167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3565b60ff7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005416612d6c57565b7fd93c06650000000000000000000000000000000000000000000000000000000060005260046000fd5b91907f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08411612e20579160209360809260ff60009560405194855216868401526040830152606082015282805260015afa156124445760005173ffffffffffffffffffffffffffffffffffffffff811615612e145790600090600090565b50600090600190600090565b50505060009160039190565b9190916004811015612ee95780612e4257509050565b600060018203612e76577ff645eedf0000000000000000000000000000000000000000000000000000000060005260046000fd5b5060028103612ead57827ffce698f70000000000000000000000000000000000000000000000000000000060005260045260246000fd5b9091600360009214612ebd575050565b602492507fd78bce0c000000000000000000000000000000000000000000000000000000008252600452fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff16600052600060205260ff60406000205416612f4557565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f4163636f756e742069732066726f7a656e0000000000000000000000000000006044820152fd5b6040517fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d1025490600081612fd584612765565b9182825260208201946001811690816000146130dd575060011461305e575b613000925003826126af565b5190811561300c572090565b50507fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d1005480156130395790565b507fc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a47090565b507fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d102600090815290917f42ad5d3e1f2e6e70edcf6d991b8a3023d3fca8047a131592f9edb9fd9b89d57d5b8183106130c157505090602061300092820101612ff4565b60209193508060019154838588010152019101909183926130a9565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001686525061300092151560051b82016020019050612ff4565b6040517fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d103549060008161314984612765565b91828252602082019460018116908160001461322c57506001146131ad575b613174925003826126af565b51908115613180572090565b50507fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d1015480156130395790565b507fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d103600090815290917f5f9ce34815f8e11431c7bb75a8e6886a91478f7ffc1dbb0a98dc240fddd76b755b81831061321057505090602061317492820101613168565b60209193508060019154838588010152019101909183926131f8565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001686525061317492151560051b8201602001905061316856fea164736f6c634300081a000a",
}

var StablecoinABI = StablecoinMetaData.ABI

var StablecoinBin = StablecoinMetaData.Bin

func DeployStablecoin(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Stablecoin, error) {
	parsed, err := StablecoinMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StablecoinBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Stablecoin{address: address, abi: *parsed, StablecoinCaller: StablecoinCaller{contract: contract}, StablecoinTransactor: StablecoinTransactor{contract: contract}, StablecoinFilterer: StablecoinFilterer{contract: contract}}, nil
}

type Stablecoin struct {
	address common.Address
	abi     abi.ABI
	StablecoinCaller
	StablecoinTransactor
	StablecoinFilterer
}

type StablecoinCaller struct {
	contract *bind.BoundContract
}

type StablecoinTransactor struct {
	contract *bind.BoundContract
}

type StablecoinFilterer struct {
	contract *bind.BoundContract
}

type StablecoinSession struct {
	Contract     *Stablecoin
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type StablecoinCallerSession struct {
	Contract *StablecoinCaller
	CallOpts bind.CallOpts
}

type StablecoinTransactorSession struct {
	Contract     *StablecoinTransactor
	TransactOpts bind.TransactOpts
}

type StablecoinRaw struct {
	Contract *Stablecoin
}

type StablecoinCallerRaw struct {
	Contract *StablecoinCaller
}

type StablecoinTransactorRaw struct {
	Contract *StablecoinTransactor
}

func NewStablecoin(address common.Address, backend bind.ContractBackend) (*Stablecoin, error) {
	abi, err := abi.JSON(strings.NewReader(StablecoinABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindStablecoin(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Stablecoin{address: address, abi: abi, StablecoinCaller: StablecoinCaller{contract: contract}, StablecoinTransactor: StablecoinTransactor{contract: contract}, StablecoinFilterer: StablecoinFilterer{contract: contract}}, nil
}

func NewStablecoinCaller(address common.Address, caller bind.ContractCaller) (*StablecoinCaller, error) {
	contract, err := bindStablecoin(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StablecoinCaller{contract: contract}, nil
}

func NewStablecoinTransactor(address common.Address, transactor bind.ContractTransactor) (*StablecoinTransactor, error) {
	contract, err := bindStablecoin(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StablecoinTransactor{contract: contract}, nil
}

func NewStablecoinFilterer(address common.Address, filterer bind.ContractFilterer) (*StablecoinFilterer, error) {
	contract, err := bindStablecoin(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StablecoinFilterer{contract: contract}, nil
}

func bindStablecoin(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StablecoinMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_Stablecoin *StablecoinRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Stablecoin.Contract.StablecoinCaller.contract.Call(opts, result, method, params...)
}

func (_Stablecoin *StablecoinRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stablecoin.Contract.StablecoinTransactor.contract.Transfer(opts)
}

func (_Stablecoin *StablecoinRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Stablecoin.Contract.StablecoinTransactor.contract.Transact(opts, method, params...)
}

func (_Stablecoin *StablecoinCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Stablecoin.Contract.contract.Call(opts, result, method, params...)
}

func (_Stablecoin *StablecoinTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stablecoin.Contract.contract.Transfer(opts)
}

func (_Stablecoin *StablecoinTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Stablecoin.Contract.contract.Transact(opts, method, params...)
}

func (_Stablecoin *StablecoinCaller) DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Stablecoin.contract.Call(opts, &out, "DOMAIN_SEPARATOR")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_Stablecoin *StablecoinSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _Stablecoin.Contract.DOMAINSEPARATOR(&_Stablecoin.CallOpts)
}

func (_Stablecoin *StablecoinCallerSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _Stablecoin.Contract.DOMAINSEPARATOR(&_Stablecoin.CallOpts)
}

func (_Stablecoin *StablecoinCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Stablecoin.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_Stablecoin *StablecoinSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Stablecoin.Contract.Allowance(&_Stablecoin.CallOpts, owner, spender)
}

func (_Stablecoin *StablecoinCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Stablecoin.Contract.Allowance(&_Stablecoin.CallOpts, owner, spender)
}

func (_Stablecoin *StablecoinCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Stablecoin.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_Stablecoin *StablecoinSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _Stablecoin.Contract.BalanceOf(&_Stablecoin.CallOpts, account)
}

func (_Stablecoin *StablecoinCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _Stablecoin.Contract.BalanceOf(&_Stablecoin.CallOpts, account)
}

func (_Stablecoin *StablecoinCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Stablecoin.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_Stablecoin *StablecoinSession) Decimals() (uint8, error) {
	return _Stablecoin.Contract.Decimals(&_Stablecoin.CallOpts)
}

func (_Stablecoin *StablecoinCallerSession) Decimals() (uint8, error) {
	return _Stablecoin.Contract.Decimals(&_Stablecoin.CallOpts)
}

func (_Stablecoin *StablecoinCaller) Eip712Domain(opts *bind.CallOpts) (Eip712Domain,

	error) {
	var out []interface{}
	err := _Stablecoin.contract.Call(opts, &out, "eip712Domain")

	outstruct := new(Eip712Domain)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fields = *abi.ConvertType(out[0], new([1]byte)).(*[1]byte)
	outstruct.Name = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Version = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.ChainId = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.VerifyingContract = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Salt = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.Extensions = *abi.ConvertType(out[6], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

func (_Stablecoin *StablecoinSession) Eip712Domain() (Eip712Domain,

	error) {
	return _Stablecoin.Contract.Eip712Domain(&_Stablecoin.CallOpts)
}

func (_Stablecoin *StablecoinCallerSession) Eip712Domain() (Eip712Domain,

	error) {
	return _Stablecoin.Contract.Eip712Domain(&_Stablecoin.CallOpts)
}

func (_Stablecoin *StablecoinCaller) Frozen(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Stablecoin.contract.Call(opts, &out, "frozen", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_Stablecoin *StablecoinSession) Frozen(arg0 common.Address) (bool, error) {
	return _Stablecoin.Contract.Frozen(&_Stablecoin.CallOpts, arg0)
}

func (_Stablecoin *StablecoinCallerSession) Frozen(arg0 common.Address) (bool, error) {
	return _Stablecoin.Contract.Frozen(&_Stablecoin.CallOpts, arg0)
}

func (_Stablecoin *StablecoinCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Stablecoin.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_Stablecoin *StablecoinSession) Name() (string, error) {
	return _Stablecoin.Contract.Name(&_Stablecoin.CallOpts)
}

func (_Stablecoin *StablecoinCallerSession) Name() (string, error) {
	return _Stablecoin.Contract.Name(&_Stablecoin.CallOpts)
}

func (_Stablecoin *StablecoinCaller) Nonces(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Stablecoin.contract.Call(opts, &out, "nonces", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_Stablecoin *StablecoinSession) Nonces(owner common.Address) (*big.Int, error) {
	return _Stablecoin.Contract.Nonces(&_Stablecoin.CallOpts, owner)
}

func (_Stablecoin *StablecoinCallerSession) Nonces(owner common.Address) (*big.Int, error) {
	return _Stablecoin.Contract.Nonces(&_Stablecoin.CallOpts, owner)
}

func (_Stablecoin *StablecoinCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Stablecoin.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_Stablecoin *StablecoinSession) Owner() (common.Address, error) {
	return _Stablecoin.Contract.Owner(&_Stablecoin.CallOpts)
}

func (_Stablecoin *StablecoinCallerSession) Owner() (common.Address, error) {
	return _Stablecoin.Contract.Owner(&_Stablecoin.CallOpts)
}

func (_Stablecoin *StablecoinCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Stablecoin.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_Stablecoin *StablecoinSession) Paused() (bool, error) {
	return _Stablecoin.Contract.Paused(&_Stablecoin.CallOpts)
}

func (_Stablecoin *StablecoinCallerSession) Paused() (bool, error) {
	return _Stablecoin.Contract.Paused(&_Stablecoin.CallOpts)
}

func (_Stablecoin *StablecoinCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Stablecoin.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_Stablecoin *StablecoinSession) PendingOwner() (common.Address, error) {
	return _Stablecoin.Contract.PendingOwner(&_Stablecoin.CallOpts)
}

func (_Stablecoin *StablecoinCallerSession) PendingOwner() (common.Address, error) {
	return _Stablecoin.Contract.PendingOwner(&_Stablecoin.CallOpts)
}

func (_Stablecoin *StablecoinCaller) RenounceOwnership(opts *bind.CallOpts) error {
	var out []interface{}
	err := _Stablecoin.contract.Call(opts, &out, "renounceOwnership")

	if err != nil {
		return err
	}

	return err

}

func (_Stablecoin *StablecoinSession) RenounceOwnership() error {
	return _Stablecoin.Contract.RenounceOwnership(&_Stablecoin.CallOpts)
}

func (_Stablecoin *StablecoinCallerSession) RenounceOwnership() error {
	return _Stablecoin.Contract.RenounceOwnership(&_Stablecoin.CallOpts)
}

func (_Stablecoin *StablecoinCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Stablecoin.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_Stablecoin *StablecoinSession) Symbol() (string, error) {
	return _Stablecoin.Contract.Symbol(&_Stablecoin.CallOpts)
}

func (_Stablecoin *StablecoinCallerSession) Symbol() (string, error) {
	return _Stablecoin.Contract.Symbol(&_Stablecoin.CallOpts)
}

func (_Stablecoin *StablecoinCaller) Test(opts *bind.CallOpts) error {
	var out []interface{}
	err := _Stablecoin.contract.Call(opts, &out, "test")

	if err != nil {
		return err
	}

	return err

}

func (_Stablecoin *StablecoinSession) Test() error {
	return _Stablecoin.Contract.Test(&_Stablecoin.CallOpts)
}

func (_Stablecoin *StablecoinCallerSession) Test() error {
	return _Stablecoin.Contract.Test(&_Stablecoin.CallOpts)
}

func (_Stablecoin *StablecoinCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Stablecoin.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_Stablecoin *StablecoinSession) TotalSupply() (*big.Int, error) {
	return _Stablecoin.Contract.TotalSupply(&_Stablecoin.CallOpts)
}

func (_Stablecoin *StablecoinCallerSession) TotalSupply() (*big.Int, error) {
	return _Stablecoin.Contract.TotalSupply(&_Stablecoin.CallOpts)
}

func (_Stablecoin *StablecoinTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stablecoin.contract.Transact(opts, "acceptOwnership")
}

func (_Stablecoin *StablecoinSession) AcceptOwnership() (*types.Transaction, error) {
	return _Stablecoin.Contract.AcceptOwnership(&_Stablecoin.TransactOpts)
}

func (_Stablecoin *StablecoinTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Stablecoin.Contract.AcceptOwnership(&_Stablecoin.TransactOpts)
}

func (_Stablecoin *StablecoinTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Stablecoin.contract.Transact(opts, "approve", spender, value)
}

func (_Stablecoin *StablecoinSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Stablecoin.Contract.Approve(&_Stablecoin.TransactOpts, spender, value)
}

func (_Stablecoin *StablecoinTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Stablecoin.Contract.Approve(&_Stablecoin.TransactOpts, spender, value)
}

func (_Stablecoin *StablecoinTransactor) Burn(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _Stablecoin.contract.Transact(opts, "burn", amount)
}

func (_Stablecoin *StablecoinSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _Stablecoin.Contract.Burn(&_Stablecoin.TransactOpts, amount)
}

func (_Stablecoin *StablecoinTransactorSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _Stablecoin.Contract.Burn(&_Stablecoin.TransactOpts, amount)
}

func (_Stablecoin *StablecoinTransactor) Drain(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _Stablecoin.contract.Transact(opts, "drain", account)
}

func (_Stablecoin *StablecoinSession) Drain(account common.Address) (*types.Transaction, error) {
	return _Stablecoin.Contract.Drain(&_Stablecoin.TransactOpts, account)
}

func (_Stablecoin *StablecoinTransactorSession) Drain(account common.Address) (*types.Transaction, error) {
	return _Stablecoin.Contract.Drain(&_Stablecoin.TransactOpts, account)
}

func (_Stablecoin *StablecoinTransactor) Freeze(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _Stablecoin.contract.Transact(opts, "freeze", account)
}

func (_Stablecoin *StablecoinSession) Freeze(account common.Address) (*types.Transaction, error) {
	return _Stablecoin.Contract.Freeze(&_Stablecoin.TransactOpts, account)
}

func (_Stablecoin *StablecoinTransactorSession) Freeze(account common.Address) (*types.Transaction, error) {
	return _Stablecoin.Contract.Freeze(&_Stablecoin.TransactOpts, account)
}

func (_Stablecoin *StablecoinTransactor) Initialize(opts *bind.TransactOpts, _name string, _symbol string) (*types.Transaction, error) {
	return _Stablecoin.contract.Transact(opts, "initialize", _name, _symbol)
}

func (_Stablecoin *StablecoinSession) Initialize(_name string, _symbol string) (*types.Transaction, error) {
	return _Stablecoin.Contract.Initialize(&_Stablecoin.TransactOpts, _name, _symbol)
}

func (_Stablecoin *StablecoinTransactorSession) Initialize(_name string, _symbol string) (*types.Transaction, error) {
	return _Stablecoin.Contract.Initialize(&_Stablecoin.TransactOpts, _name, _symbol)
}

func (_Stablecoin *StablecoinTransactor) Mint(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _Stablecoin.contract.Transact(opts, "mint", amount)
}

func (_Stablecoin *StablecoinSession) Mint(amount *big.Int) (*types.Transaction, error) {
	return _Stablecoin.Contract.Mint(&_Stablecoin.TransactOpts, amount)
}

func (_Stablecoin *StablecoinTransactorSession) Mint(amount *big.Int) (*types.Transaction, error) {
	return _Stablecoin.Contract.Mint(&_Stablecoin.TransactOpts, amount)
}

func (_Stablecoin *StablecoinTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stablecoin.contract.Transact(opts, "pause")
}

func (_Stablecoin *StablecoinSession) Pause() (*types.Transaction, error) {
	return _Stablecoin.Contract.Pause(&_Stablecoin.TransactOpts)
}

func (_Stablecoin *StablecoinTransactorSession) Pause() (*types.Transaction, error) {
	return _Stablecoin.Contract.Pause(&_Stablecoin.TransactOpts)
}

func (_Stablecoin *StablecoinTransactor) Permit(opts *bind.TransactOpts, owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Stablecoin.contract.Transact(opts, "permit", owner, spender, value, deadline, v, r, s)
}

func (_Stablecoin *StablecoinSession) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Stablecoin.Contract.Permit(&_Stablecoin.TransactOpts, owner, spender, value, deadline, v, r, s)
}

func (_Stablecoin *StablecoinTransactorSession) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Stablecoin.Contract.Permit(&_Stablecoin.TransactOpts, owner, spender, value, deadline, v, r, s)
}

func (_Stablecoin *StablecoinTransactor) RecoverERC20(opts *bind.TransactOpts, token common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Stablecoin.contract.Transact(opts, "recoverERC20", token, recipient, amount)
}

func (_Stablecoin *StablecoinSession) RecoverERC20(token common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Stablecoin.Contract.RecoverERC20(&_Stablecoin.TransactOpts, token, recipient, amount)
}

func (_Stablecoin *StablecoinTransactorSession) RecoverERC20(token common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Stablecoin.Contract.RecoverERC20(&_Stablecoin.TransactOpts, token, recipient, amount)
}

func (_Stablecoin *StablecoinTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Stablecoin.contract.Transact(opts, "transfer", to, value)
}

func (_Stablecoin *StablecoinSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Stablecoin.Contract.Transfer(&_Stablecoin.TransactOpts, to, value)
}

func (_Stablecoin *StablecoinTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Stablecoin.Contract.Transfer(&_Stablecoin.TransactOpts, to, value)
}

func (_Stablecoin *StablecoinTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Stablecoin.contract.Transact(opts, "transferFrom", from, to, value)
}

func (_Stablecoin *StablecoinSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Stablecoin.Contract.TransferFrom(&_Stablecoin.TransactOpts, from, to, value)
}

func (_Stablecoin *StablecoinTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Stablecoin.Contract.TransferFrom(&_Stablecoin.TransactOpts, from, to, value)
}

func (_Stablecoin *StablecoinTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Stablecoin.contract.Transact(opts, "transferOwnership", newOwner)
}

func (_Stablecoin *StablecoinSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Stablecoin.Contract.TransferOwnership(&_Stablecoin.TransactOpts, newOwner)
}

func (_Stablecoin *StablecoinTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Stablecoin.Contract.TransferOwnership(&_Stablecoin.TransactOpts, newOwner)
}

func (_Stablecoin *StablecoinTransactor) Unfreeze(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _Stablecoin.contract.Transact(opts, "unfreeze", account)
}

func (_Stablecoin *StablecoinSession) Unfreeze(account common.Address) (*types.Transaction, error) {
	return _Stablecoin.Contract.Unfreeze(&_Stablecoin.TransactOpts, account)
}

func (_Stablecoin *StablecoinTransactorSession) Unfreeze(account common.Address) (*types.Transaction, error) {
	return _Stablecoin.Contract.Unfreeze(&_Stablecoin.TransactOpts, account)
}

func (_Stablecoin *StablecoinTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stablecoin.contract.Transact(opts, "unpause")
}

func (_Stablecoin *StablecoinSession) Unpause() (*types.Transaction, error) {
	return _Stablecoin.Contract.Unpause(&_Stablecoin.TransactOpts)
}

func (_Stablecoin *StablecoinTransactorSession) Unpause() (*types.Transaction, error) {
	return _Stablecoin.Contract.Unpause(&_Stablecoin.TransactOpts)
}

type StablecoinApprovalIterator struct {
	Event *StablecoinApproval

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StablecoinApprovalIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StablecoinApproval)
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
		it.Event = new(StablecoinApproval)
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

func (it *StablecoinApprovalIterator) Error() error {
	return it.fail
}

func (it *StablecoinApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StablecoinApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log
}

func (_Stablecoin *StablecoinFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*StablecoinApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Stablecoin.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &StablecoinApprovalIterator{contract: _Stablecoin.contract, event: "Approval", logs: logs, sub: sub}, nil
}

func (_Stablecoin *StablecoinFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *StablecoinApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Stablecoin.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StablecoinApproval)
				if err := _Stablecoin.contract.UnpackLog(event, "Approval", log); err != nil {
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

func (_Stablecoin *StablecoinFilterer) ParseApproval(log types.Log) (*StablecoinApproval, error) {
	event := new(StablecoinApproval)
	if err := _Stablecoin.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StablecoinBurnIterator struct {
	Event *StablecoinBurn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StablecoinBurnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StablecoinBurn)
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
		it.Event = new(StablecoinBurn)
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

func (it *StablecoinBurnIterator) Error() error {
	return it.fail
}

func (it *StablecoinBurnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StablecoinBurn struct {
	Caller common.Address
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_Stablecoin *StablecoinFilterer) FilterBurn(opts *bind.FilterOpts, caller []common.Address, from []common.Address) (*StablecoinBurnIterator, error) {

	var callerRule []interface{}
	for _, callerItem := range caller {
		callerRule = append(callerRule, callerItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _Stablecoin.contract.FilterLogs(opts, "Burn", callerRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &StablecoinBurnIterator{contract: _Stablecoin.contract, event: "Burn", logs: logs, sub: sub}, nil
}

func (_Stablecoin *StablecoinFilterer) WatchBurn(opts *bind.WatchOpts, sink chan<- *StablecoinBurn, caller []common.Address, from []common.Address) (event.Subscription, error) {

	var callerRule []interface{}
	for _, callerItem := range caller {
		callerRule = append(callerRule, callerItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _Stablecoin.contract.WatchLogs(opts, "Burn", callerRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StablecoinBurn)
				if err := _Stablecoin.contract.UnpackLog(event, "Burn", log); err != nil {
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

func (_Stablecoin *StablecoinFilterer) ParseBurn(log types.Log) (*StablecoinBurn, error) {
	event := new(StablecoinBurn)
	if err := _Stablecoin.contract.UnpackLog(event, "Burn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StablecoinEIP712DomainChangedIterator struct {
	Event *StablecoinEIP712DomainChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StablecoinEIP712DomainChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StablecoinEIP712DomainChanged)
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
		it.Event = new(StablecoinEIP712DomainChanged)
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

func (it *StablecoinEIP712DomainChangedIterator) Error() error {
	return it.fail
}

func (it *StablecoinEIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StablecoinEIP712DomainChanged struct {
	Raw types.Log
}

func (_Stablecoin *StablecoinFilterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*StablecoinEIP712DomainChangedIterator, error) {

	logs, sub, err := _Stablecoin.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &StablecoinEIP712DomainChangedIterator{contract: _Stablecoin.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

func (_Stablecoin *StablecoinFilterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *StablecoinEIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _Stablecoin.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StablecoinEIP712DomainChanged)
				if err := _Stablecoin.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
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

func (_Stablecoin *StablecoinFilterer) ParseEIP712DomainChanged(log types.Log) (*StablecoinEIP712DomainChanged, error) {
	event := new(StablecoinEIP712DomainChanged)
	if err := _Stablecoin.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StablecoinFreezeIterator struct {
	Event *StablecoinFreeze

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StablecoinFreezeIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StablecoinFreeze)
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
		it.Event = new(StablecoinFreeze)
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

func (it *StablecoinFreezeIterator) Error() error {
	return it.fail
}

func (it *StablecoinFreezeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StablecoinFreeze struct {
	Caller  common.Address
	Account common.Address
	Raw     types.Log
}

func (_Stablecoin *StablecoinFilterer) FilterFreeze(opts *bind.FilterOpts, caller []common.Address, account []common.Address) (*StablecoinFreezeIterator, error) {

	var callerRule []interface{}
	for _, callerItem := range caller {
		callerRule = append(callerRule, callerItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Stablecoin.contract.FilterLogs(opts, "Freeze", callerRule, accountRule)
	if err != nil {
		return nil, err
	}
	return &StablecoinFreezeIterator{contract: _Stablecoin.contract, event: "Freeze", logs: logs, sub: sub}, nil
}

func (_Stablecoin *StablecoinFilterer) WatchFreeze(opts *bind.WatchOpts, sink chan<- *StablecoinFreeze, caller []common.Address, account []common.Address) (event.Subscription, error) {

	var callerRule []interface{}
	for _, callerItem := range caller {
		callerRule = append(callerRule, callerItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Stablecoin.contract.WatchLogs(opts, "Freeze", callerRule, accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StablecoinFreeze)
				if err := _Stablecoin.contract.UnpackLog(event, "Freeze", log); err != nil {
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

func (_Stablecoin *StablecoinFilterer) ParseFreeze(log types.Log) (*StablecoinFreeze, error) {
	event := new(StablecoinFreeze)
	if err := _Stablecoin.contract.UnpackLog(event, "Freeze", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StablecoinInitializedIterator struct {
	Event *StablecoinInitialized

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StablecoinInitializedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StablecoinInitialized)
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
		it.Event = new(StablecoinInitialized)
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

func (it *StablecoinInitializedIterator) Error() error {
	return it.fail
}

func (it *StablecoinInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StablecoinInitialized struct {
	Version uint64
	Raw     types.Log
}

func (_Stablecoin *StablecoinFilterer) FilterInitialized(opts *bind.FilterOpts) (*StablecoinInitializedIterator, error) {

	logs, sub, err := _Stablecoin.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &StablecoinInitializedIterator{contract: _Stablecoin.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

func (_Stablecoin *StablecoinFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *StablecoinInitialized) (event.Subscription, error) {

	logs, sub, err := _Stablecoin.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StablecoinInitialized)
				if err := _Stablecoin.contract.UnpackLog(event, "Initialized", log); err != nil {
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

func (_Stablecoin *StablecoinFilterer) ParseInitialized(log types.Log) (*StablecoinInitialized, error) {
	event := new(StablecoinInitialized)
	if err := _Stablecoin.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StablecoinMintIterator struct {
	Event *StablecoinMint

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StablecoinMintIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StablecoinMint)
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
		it.Event = new(StablecoinMint)
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

func (it *StablecoinMintIterator) Error() error {
	return it.fail
}

func (it *StablecoinMintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StablecoinMint struct {
	Caller common.Address
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_Stablecoin *StablecoinFilterer) FilterMint(opts *bind.FilterOpts, caller []common.Address, to []common.Address) (*StablecoinMintIterator, error) {

	var callerRule []interface{}
	for _, callerItem := range caller {
		callerRule = append(callerRule, callerItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Stablecoin.contract.FilterLogs(opts, "Mint", callerRule, toRule)
	if err != nil {
		return nil, err
	}
	return &StablecoinMintIterator{contract: _Stablecoin.contract, event: "Mint", logs: logs, sub: sub}, nil
}

func (_Stablecoin *StablecoinFilterer) WatchMint(opts *bind.WatchOpts, sink chan<- *StablecoinMint, caller []common.Address, to []common.Address) (event.Subscription, error) {

	var callerRule []interface{}
	for _, callerItem := range caller {
		callerRule = append(callerRule, callerItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Stablecoin.contract.WatchLogs(opts, "Mint", callerRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StablecoinMint)
				if err := _Stablecoin.contract.UnpackLog(event, "Mint", log); err != nil {
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

func (_Stablecoin *StablecoinFilterer) ParseMint(log types.Log) (*StablecoinMint, error) {
	event := new(StablecoinMint)
	if err := _Stablecoin.contract.UnpackLog(event, "Mint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StablecoinOwnershipTransferStartedIterator struct {
	Event *StablecoinOwnershipTransferStarted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StablecoinOwnershipTransferStartedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StablecoinOwnershipTransferStarted)
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
		it.Event = new(StablecoinOwnershipTransferStarted)
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

func (it *StablecoinOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

func (it *StablecoinOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StablecoinOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log
}

func (_Stablecoin *StablecoinFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*StablecoinOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Stablecoin.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &StablecoinOwnershipTransferStartedIterator{contract: _Stablecoin.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

func (_Stablecoin *StablecoinFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *StablecoinOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Stablecoin.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StablecoinOwnershipTransferStarted)
				if err := _Stablecoin.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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

func (_Stablecoin *StablecoinFilterer) ParseOwnershipTransferStarted(log types.Log) (*StablecoinOwnershipTransferStarted, error) {
	event := new(StablecoinOwnershipTransferStarted)
	if err := _Stablecoin.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StablecoinOwnershipTransferredIterator struct {
	Event *StablecoinOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StablecoinOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StablecoinOwnershipTransferred)
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
		it.Event = new(StablecoinOwnershipTransferred)
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

func (it *StablecoinOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *StablecoinOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StablecoinOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log
}

func (_Stablecoin *StablecoinFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*StablecoinOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Stablecoin.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &StablecoinOwnershipTransferredIterator{contract: _Stablecoin.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_Stablecoin *StablecoinFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *StablecoinOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Stablecoin.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StablecoinOwnershipTransferred)
				if err := _Stablecoin.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_Stablecoin *StablecoinFilterer) ParseOwnershipTransferred(log types.Log) (*StablecoinOwnershipTransferred, error) {
	event := new(StablecoinOwnershipTransferred)
	if err := _Stablecoin.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StablecoinPausedIterator struct {
	Event *StablecoinPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StablecoinPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StablecoinPaused)
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
		it.Event = new(StablecoinPaused)
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

func (it *StablecoinPausedIterator) Error() error {
	return it.fail
}

func (it *StablecoinPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StablecoinPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_Stablecoin *StablecoinFilterer) FilterPaused(opts *bind.FilterOpts) (*StablecoinPausedIterator, error) {

	logs, sub, err := _Stablecoin.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &StablecoinPausedIterator{contract: _Stablecoin.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_Stablecoin *StablecoinFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *StablecoinPaused) (event.Subscription, error) {

	logs, sub, err := _Stablecoin.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StablecoinPaused)
				if err := _Stablecoin.contract.UnpackLog(event, "Paused", log); err != nil {
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

func (_Stablecoin *StablecoinFilterer) ParsePaused(log types.Log) (*StablecoinPaused, error) {
	event := new(StablecoinPaused)
	if err := _Stablecoin.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StablecoinTransferIterator struct {
	Event *StablecoinTransfer

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StablecoinTransferIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StablecoinTransfer)
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
		it.Event = new(StablecoinTransfer)
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

func (it *StablecoinTransferIterator) Error() error {
	return it.fail
}

func (it *StablecoinTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StablecoinTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log
}

func (_Stablecoin *StablecoinFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*StablecoinTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Stablecoin.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &StablecoinTransferIterator{contract: _Stablecoin.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

func (_Stablecoin *StablecoinFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *StablecoinTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Stablecoin.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StablecoinTransfer)
				if err := _Stablecoin.contract.UnpackLog(event, "Transfer", log); err != nil {
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

func (_Stablecoin *StablecoinFilterer) ParseTransfer(log types.Log) (*StablecoinTransfer, error) {
	event := new(StablecoinTransfer)
	if err := _Stablecoin.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StablecoinUnfreezeIterator struct {
	Event *StablecoinUnfreeze

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StablecoinUnfreezeIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StablecoinUnfreeze)
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
		it.Event = new(StablecoinUnfreeze)
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

func (it *StablecoinUnfreezeIterator) Error() error {
	return it.fail
}

func (it *StablecoinUnfreezeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StablecoinUnfreeze struct {
	Caller  common.Address
	Account common.Address
	Raw     types.Log
}

func (_Stablecoin *StablecoinFilterer) FilterUnfreeze(opts *bind.FilterOpts, caller []common.Address, account []common.Address) (*StablecoinUnfreezeIterator, error) {

	var callerRule []interface{}
	for _, callerItem := range caller {
		callerRule = append(callerRule, callerItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Stablecoin.contract.FilterLogs(opts, "Unfreeze", callerRule, accountRule)
	if err != nil {
		return nil, err
	}
	return &StablecoinUnfreezeIterator{contract: _Stablecoin.contract, event: "Unfreeze", logs: logs, sub: sub}, nil
}

func (_Stablecoin *StablecoinFilterer) WatchUnfreeze(opts *bind.WatchOpts, sink chan<- *StablecoinUnfreeze, caller []common.Address, account []common.Address) (event.Subscription, error) {

	var callerRule []interface{}
	for _, callerItem := range caller {
		callerRule = append(callerRule, callerItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Stablecoin.contract.WatchLogs(opts, "Unfreeze", callerRule, accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StablecoinUnfreeze)
				if err := _Stablecoin.contract.UnpackLog(event, "Unfreeze", log); err != nil {
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

func (_Stablecoin *StablecoinFilterer) ParseUnfreeze(log types.Log) (*StablecoinUnfreeze, error) {
	event := new(StablecoinUnfreeze)
	if err := _Stablecoin.contract.UnpackLog(event, "Unfreeze", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type StablecoinUnpausedIterator struct {
	Event *StablecoinUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *StablecoinUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StablecoinUnpaused)
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
		it.Event = new(StablecoinUnpaused)
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

func (it *StablecoinUnpausedIterator) Error() error {
	return it.fail
}

func (it *StablecoinUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type StablecoinUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_Stablecoin *StablecoinFilterer) FilterUnpaused(opts *bind.FilterOpts) (*StablecoinUnpausedIterator, error) {

	logs, sub, err := _Stablecoin.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &StablecoinUnpausedIterator{contract: _Stablecoin.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_Stablecoin *StablecoinFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *StablecoinUnpaused) (event.Subscription, error) {

	logs, sub, err := _Stablecoin.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(StablecoinUnpaused)
				if err := _Stablecoin.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

func (_Stablecoin *StablecoinFilterer) ParseUnpaused(log types.Log) (*StablecoinUnpaused, error) {
	event := new(StablecoinUnpaused)
	if err := _Stablecoin.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type Eip712Domain struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}

func (_Stablecoin *Stablecoin) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _Stablecoin.abi.Events["Approval"].ID:
		return _Stablecoin.ParseApproval(log)
	case _Stablecoin.abi.Events["Burn"].ID:
		return _Stablecoin.ParseBurn(log)
	case _Stablecoin.abi.Events["EIP712DomainChanged"].ID:
		return _Stablecoin.ParseEIP712DomainChanged(log)
	case _Stablecoin.abi.Events["Freeze"].ID:
		return _Stablecoin.ParseFreeze(log)
	case _Stablecoin.abi.Events["Initialized"].ID:
		return _Stablecoin.ParseInitialized(log)
	case _Stablecoin.abi.Events["Mint"].ID:
		return _Stablecoin.ParseMint(log)
	case _Stablecoin.abi.Events["OwnershipTransferStarted"].ID:
		return _Stablecoin.ParseOwnershipTransferStarted(log)
	case _Stablecoin.abi.Events["OwnershipTransferred"].ID:
		return _Stablecoin.ParseOwnershipTransferred(log)
	case _Stablecoin.abi.Events["Paused"].ID:
		return _Stablecoin.ParsePaused(log)
	case _Stablecoin.abi.Events["Transfer"].ID:
		return _Stablecoin.ParseTransfer(log)
	case _Stablecoin.abi.Events["Unfreeze"].ID:
		return _Stablecoin.ParseUnfreeze(log)
	case _Stablecoin.abi.Events["Unpaused"].ID:
		return _Stablecoin.ParseUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (StablecoinApproval) Topic() common.Hash {
	return common.HexToHash("0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925")
}

func (StablecoinBurn) Topic() common.Hash {
	return common.HexToHash("0xbac40739b0d4ca32fa2d82fc91630465ba3eddd1598da6fca393b26fb63b9453")
}

func (StablecoinEIP712DomainChanged) Topic() common.Hash {
	return common.HexToHash("0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31")
}

func (StablecoinFreeze) Topic() common.Hash {
	return common.HexToHash("0x51d18786e9cb144f87d46e7b796309ea84c7c687d91e09c97f051eacf59bc528")
}

func (StablecoinInitialized) Topic() common.Hash {
	return common.HexToHash("0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2")
}

func (StablecoinMint) Topic() common.Hash {
	return common.HexToHash("0xab8530f87dc9b59234c4623bf917212bb2536d647574c8e7e5da92c2ede0c9f8")
}

func (StablecoinOwnershipTransferStarted) Topic() common.Hash {
	return common.HexToHash("0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700")
}

func (StablecoinOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (StablecoinPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (StablecoinTransfer) Topic() common.Hash {
	return common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
}

func (StablecoinUnfreeze) Topic() common.Hash {
	return common.HexToHash("0x4f3ab9ff0cc4f039268532098e01239544b0420171876e36889d01c62c784c79")
}

func (StablecoinUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (_Stablecoin *Stablecoin) Address() common.Address {
	return _Stablecoin.address
}

type StablecoinInterface interface {
	DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error)

	Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error)

	BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error)

	Decimals(opts *bind.CallOpts) (uint8, error)

	Eip712Domain(opts *bind.CallOpts) (Eip712Domain,

		error)

	Frozen(opts *bind.CallOpts, arg0 common.Address) (bool, error)

	Name(opts *bind.CallOpts) (string, error)

	Nonces(opts *bind.CallOpts, owner common.Address) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	Paused(opts *bind.CallOpts) (bool, error)

	PendingOwner(opts *bind.CallOpts) (common.Address, error)

	RenounceOwnership(opts *bind.CallOpts) error

	Symbol(opts *bind.CallOpts) (string, error)

	Test(opts *bind.CallOpts) error

	TotalSupply(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error)

	Burn(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	Drain(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error)

	Freeze(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error)

	Initialize(opts *bind.TransactOpts, _name string, _symbol string) (*types.Transaction, error)

	Mint(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	Pause(opts *bind.TransactOpts) (*types.Transaction, error)

	Permit(opts *bind.TransactOpts, owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error)

	RecoverERC20(opts *bind.TransactOpts, token common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error)

	TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error)

	Unfreeze(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error)

	Unpause(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*StablecoinApprovalIterator, error)

	WatchApproval(opts *bind.WatchOpts, sink chan<- *StablecoinApproval, owner []common.Address, spender []common.Address) (event.Subscription, error)

	ParseApproval(log types.Log) (*StablecoinApproval, error)

	FilterBurn(opts *bind.FilterOpts, caller []common.Address, from []common.Address) (*StablecoinBurnIterator, error)

	WatchBurn(opts *bind.WatchOpts, sink chan<- *StablecoinBurn, caller []common.Address, from []common.Address) (event.Subscription, error)

	ParseBurn(log types.Log) (*StablecoinBurn, error)

	FilterEIP712DomainChanged(opts *bind.FilterOpts) (*StablecoinEIP712DomainChangedIterator, error)

	WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *StablecoinEIP712DomainChanged) (event.Subscription, error)

	ParseEIP712DomainChanged(log types.Log) (*StablecoinEIP712DomainChanged, error)

	FilterFreeze(opts *bind.FilterOpts, caller []common.Address, account []common.Address) (*StablecoinFreezeIterator, error)

	WatchFreeze(opts *bind.WatchOpts, sink chan<- *StablecoinFreeze, caller []common.Address, account []common.Address) (event.Subscription, error)

	ParseFreeze(log types.Log) (*StablecoinFreeze, error)

	FilterInitialized(opts *bind.FilterOpts) (*StablecoinInitializedIterator, error)

	WatchInitialized(opts *bind.WatchOpts, sink chan<- *StablecoinInitialized) (event.Subscription, error)

	ParseInitialized(log types.Log) (*StablecoinInitialized, error)

	FilterMint(opts *bind.FilterOpts, caller []common.Address, to []common.Address) (*StablecoinMintIterator, error)

	WatchMint(opts *bind.WatchOpts, sink chan<- *StablecoinMint, caller []common.Address, to []common.Address) (event.Subscription, error)

	ParseMint(log types.Log) (*StablecoinMint, error)

	FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*StablecoinOwnershipTransferStartedIterator, error)

	WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *StablecoinOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error)

	ParseOwnershipTransferStarted(log types.Log) (*StablecoinOwnershipTransferStarted, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*StablecoinOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *StablecoinOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*StablecoinOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*StablecoinPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *StablecoinPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*StablecoinPaused, error)

	FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*StablecoinTransferIterator, error)

	WatchTransfer(opts *bind.WatchOpts, sink chan<- *StablecoinTransfer, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseTransfer(log types.Log) (*StablecoinTransfer, error)

	FilterUnfreeze(opts *bind.FilterOpts, caller []common.Address, account []common.Address) (*StablecoinUnfreezeIterator, error)

	WatchUnfreeze(opts *bind.WatchOpts, sink chan<- *StablecoinUnfreeze, caller []common.Address, account []common.Address) (event.Subscription, error)

	ParseUnfreeze(log types.Log) (*StablecoinUnfreeze, error)

	FilterUnpaused(opts *bind.FilterOpts) (*StablecoinUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *StablecoinUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*StablecoinUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
