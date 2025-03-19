// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ping_pong_demo

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

type CCIPBaseApprovedSenderUpdate struct {
	DestChainSelector uint64
	Sender            []byte
}

type CCIPBaseChainUpdate struct {
	ChainSelector  uint64
	Allowed        bool
	Recipient      []byte
	ExtraArgsBytes []byte
}

type ClientAny2EVMMessage struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	Sender              []byte
	Data                []byte
	DestTokenAmounts    []ClientEVMTokenAmount
}

type ClientEVMTokenAmount struct {
	Token  common.Address
	Amount *big.Int
}

var PingPongDemoMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeToken\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"abandonFailedMessage\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyChainUpdates\",\"inputs\":[{\"name\":\"chains\",\"type\":\"tuple[]\",\"internalType\":\"structCCIPBase.ChainUpdate[]\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"allowed\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"recipient\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"extraArgsBytes\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ccipReceive\",\"inputs\":[{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structClient.Any2EVMMessage\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ccipSend\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"getCounterpartAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCounterpartChainSelector\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFeeToken\",\"inputs\":[],\"outputs\":[{\"name\":\"feeToken\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMessageContents\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structClient.Any2EVMMessage\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOutOfOrderExecution\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRemoteChainConfig\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"recipient\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"extraArgsBytes\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRouter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isApprovedSender\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"senderAddr\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isFailedMessage\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isPaused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"processMessage\",\"inputs\":[{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structClient.Any2EVMMessage\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"retryFailedMessage\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"s_chainConfigs\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"recipient\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"extraArgsBytes\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setCounterpart\",\"inputs\":[{\"name\":\"counterpartChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"counterpartAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPaused\",\"inputs\":[{\"name\":\"pause\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"startPingPong\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"updateApprovedSenders\",\"inputs\":[{\"name\":\"adds\",\"type\":\"tuple[]\",\"internalType\":\"structCCIPBase.ApprovedSenderUpdate[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"removes\",\"type\":\"tuple[]\",\"internalType\":\"structCCIPBase.ApprovedSenderUpdate[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateFeeToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateRouter\",\"inputs\":[{\"name\":\"newRouter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"usePreFundedFeeTokens\",\"inputs\":[],\"outputs\":[{\"name\":\"isPreFunded\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdrawTokens\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"ApprovedSenderAdded\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"recipient\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ApprovedSenderRemoved\",\"inputs\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"recipient\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CCIPRouterModified\",\"inputs\":[{\"name\":\"oldRouter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newRouter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainAdded\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"recipient\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"},{\"name\":\"extraArgsBytes\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainRemoved\",\"inputs\":[{\"name\":\"removeChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeeTokenUpdated\",\"inputs\":[{\"name\":\"oldFeeToken\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newFeeToken\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MessageAbandoned\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"tokenReceiver\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MessageFailed\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"reason\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MessageRecovered\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MessageSent\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MessageSucceeded\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OutOfOrderExecutionChange\",\"inputs\":[{\"name\":\"isOutOfOrder\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Ping\",\"inputs\":[{\"name\":\"pingPongCount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Pong\",\"inputs\":[{\"name\":\"pingPongCount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokensWithdrawnByOwner\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InvalidChain\",\"inputs\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidRecipient\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidRouter\",\"inputs\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MessageNotFailed\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"OnlySelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddressNotAllowed\",\"inputs\":[]}]",
	Bin: "0x60a0806040523461022e57604081614177803803809161001f8285610355565b83398101031261022e5780516001600160a01b038116919082900361022e57602001516001600160a01b0381169081900361022e57331561031057600080546001600160a01b0319163317905581156102ff57600280546001600160a01b0319908116841790915560078054909116821790556001608052806100bc575b604051613d2e908161044982396080518181816119ee01526132ba0152f35b604051636eb1769f60e11b815230600482015260248101839052602081604481855afa9081156102f3576000916102c1575b5061025657604051602081019263095ea7b360e01b84526024820152600019604482015260448152610121606482610355565b6000806040948551936101348786610355565b602085527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65646020860152519082865af13d15610249573d906001600160401b0382116102335784516101a5949092610196601f8201601f191660200185610355565b83523d6000602085013e610378565b8051806101b3575b5061009d565b816020918101031261022e576020015180159081150361022e576101d85780806101ad565b5162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b6064820152608490fd5b600080fd5b634e487b7160e01b600052604160045260246000fd5b916101a592606091610378565b60405162461bcd60e51b815260206004820152603660248201527f5361666545524332303a20617070726f76652066726f6d206e6f6e2d7a65726f60448201527f20746f206e6f6e2d7a65726f20616c6c6f77616e6365000000000000000000006064820152608490fd5b90506020813d6020116102eb575b816102dc60209383610355565b8101031261022e5751386100ee565b3d91506102cf565b6040513d6000823e3d90fd5b6342bcdf7f60e11b60005260046000fd5b60405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f00000000000000006044820152606490fd5b601f909101601f19168101906001600160401b0382119082101761023357604052565b919290156103da575081511561038c575090565b3b156103955790565b60405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606490fd5b8251909150156103ed5750805190602001fd5b6040519062461bcd60e51b8252602060048301528181519182602483015260005b8381106104305750508160006044809484010152601f80199101168101030190fd5b6020828201810151604487840101528593500161040e56fe608080604052600436101561001d575b50361561001b57600080fd5b005b600090813560e01c90816301ffc9a714612169575080630a9094dc146121145780630e958d6b1461206357806316c38b3c14611fd6578063181f5a7714611f575780632874d8bf14611eb15780632b6e5d6314611e5f57806335f170ef14611dfd5780635e35359e14611bf15780636939cd9714611b695780636d62d63314611a1357806372be9eb8146119b857806379ba50971461189d5780638462a2b91461166d57806385572ffb146115ab578063898068fc146115295780638da5cb5b146114d85780639fe74e2614610f99578063ae90de5514610f55578063b0f479a114610f03578063b187bd2614610ebf578063b5a1101114610bf8578063bee518a414610baf578063c851cc3214610ac3578063c89245d514610796578063ca709a2514610744578063cf6730f81461053c578063e4ca87541461039c578063e89b4485146102965763f2fde38b0361000f57346102935760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102935773ffffffffffffffffffffffffffffffffffffffff6101bc6124d8565b6101c46137e2565b1633811461023557807fffffffffffffffffffffffff0000000000000000000000000000000000000000600154161760015573ffffffffffffffffffffffffffffffffffffffff8254167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12788380a380f35b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152fd5b80fd5b5060607ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610293576102c9612227565b906024359067ffffffffffffffff821161029357366023830112156102935781600401356102f6816126dd565b9261030460405194856122ba565b8184526024602085019260061b8201019036821161039857602401915b818310610360575050506044359067ffffffffffffffff82116102935760206103588585610352366004880161272c565b916131bd565b604051908152f35b604083360312610398576020604091825161037a8161229e565b6103838661251e565b81528286013583820152815201920191610321565b8380fd5b50346102935760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610293576004356103e6816000526006602052604060002054151590565b15610511578082526004602052816040812061044a60046040519261040a84612253565b8054845267ffffffffffffffff600182015416602085015261042e600282016123ee565b604085015261043f600382016123ee565b606085015201612760565b6080820152828252600460205261046360408320612844565b61046c83613a5f565b506104756137e2565b303b1561050257816104b491604051809381927fcf6730f800000000000000000000000000000000000000000000000000000000835260048301612595565b038183305af18015610506576104ed575b50807fef3bf8c64bc480286c4f3503b870ceb23e648d2d902e31fb7bb46680da6de8ad91a280f35b816104f7916122ba565b6105025781386104c5565b5080fd5b6040513d84823e3d90fd5b7fb6e78260000000000000000000000000000000000000000000000000000000008252600452602490fd5b50346102935761054b3661266e565b30330361071c5761055e60208201612985565b67ffffffffffffffff61057e610577604085018561299a565b36916126f5565b91168084526003602052610595604085205461239b565b159081156106f8575b506106b6575060ff60085460a01c16156105b6575080f35b6105c6816060602093019061299a565b9080929181010312610502573560018101809111610689576106569060018082160361065a577f48257dc961b6f792c2b78a080dacfed693b660960a702de21cee364e20270e2f6020604051838152a15b6040519060208201526020815261062f6040826122ba565b67ffffffffffffffff60075460a01c166040519061064e6020836122ba565b8482526131bd565b5080f35b7f58b69f57828e6962d216502094c54f6562f3bf082ba758966c3454f9e37b15256020604051838152a1610617565b6024827f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b6106f4906040519182917f5075bb38000000000000000000000000000000000000000000000000000000008352602060048401526024830190612358565b0390fd5b90508352600360205260ff6107136002604086200183613197565b5416153861059e565b6004827f14d4a4e8000000000000000000000000000000000000000000000000000000008152fd5b503461029357807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261029357602073ffffffffffffffffffffffffffffffffffffffff60075416604051908152f35b50346102935760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610293576107ce6124d8565b6107d66137e2565b73ffffffffffffffffffffffffffffffffffffffff6007541680610a58575b506007549073ffffffffffffffffffffffffffffffffffffffff8116807fffffffffffffffffffffffff000000000000000000000000000000000000000084161760075580610894575b506040805173ffffffffffffffffffffffffffffffffffffffff93841681529190921660208201527f91a03e1d689caf891fe531c01e290f7b718f9c6a3af6726d6d837d2b7bd82e6791819081015b0390a180f35b6002546040517fdd62ed3e00000000000000000000000000000000000000000000000000000000815230600482015273ffffffffffffffffffffffffffffffffffffffff90911660248201819052909190602083604481855afa928315610a4d578693610a14575b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83018093116109e7576040517f095ea7b300000000000000000000000000000000000000000000000000000000602082015273ffffffffffffffffffffffffffffffffffffffff918216602482015260448101939093527f91a03e1d689caf891fe531c01e290f7b718f9c6a3af6726d6d837d2b7bd82e67949093909261088e926109de91906109d982606481015b037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081018452836122ba565b613927565b9250509161083f565b6024867f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b9092506020813d602011610a45575b81610a30602093836122ba565b81010312610a40575191386108fc565b600080fd5b3d9150610a23565b6040513d88823e3d90fd5b91610abd9073ffffffffffffffffffffffffffffffffffffffff600254169093604051917f095ea7b3000000000000000000000000000000000000000000000000000000006020840152602483015260006044830152604482526109d96064836122ba565b386107f5565b50346102935760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102935773ffffffffffffffffffffffffffffffffffffffff610b106124d8565b610b186137e2565b168015610b875773ffffffffffffffffffffffffffffffffffffffff600254827fffffffffffffffffffffffff0000000000000000000000000000000000000000821617600255167f3672b589036f39ac008505b790fcb05d484d70b65680ec64c089a3c173fdc4c88380a380f35b6004827f8579befe000000000000000000000000000000000000000000000000000000008152fd5b503461029357807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261029357602067ffffffffffffffff60075460a01c16604051908152f35b50346102935760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261029357610c30612227565b9073ffffffffffffffffffffffffffffffffffffffff610c4e6124fb565b610c566137e2565b169182158015610ead575b610b875767ffffffffffffffff8116907fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff7bffffffffffffffff00000000000000000000000000000000000000006007549260a01b16911617600755827fffffffffffffffffffffffff000000000000000000000000000000000000000060085416176008558082526003602052610d166002604084200160405185602082015260208152610d116040826122ba565b613197565b60017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082541617905560405192602084015260208352610d576040846122ba565b8152600360205260408120825167ffffffffffffffff8111610e8057610d8781610d81845461239b565b84612aaa565b6020601f8211600114610de55781908495610dd5949592610dda575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b905580f35b015190503880610da3565b828452808420907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08316855b818110610e6857509583600195969710610e31575b505050811b01905580f35b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c19169055388080610e26565b9192602060018192868b015181550194019201610e11565b6024837f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b5067ffffffffffffffff811615610c61565b503461029357807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261029357602060ff60085460a01c166040519015158152f35b503461029357807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261029357602073ffffffffffffffffffffffffffffffffffffffff60025416604051908152f35b503461029357807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261029357602060ff60085460a81c166040519015158152f35b50346102935760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102935760043567ffffffffffffffff811161050257610fe990369060040161263d565b90610ff26137e2565b82907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8181360301915b838110156114d4578060051b820135838112156114d05782016080813603126114d057604051906080820182811067ffffffffffffffff8211176114a3576040526110658161223e565b8252602081013590811515820361149f5760208301918252604081013567ffffffffffffffff811161149b5761109e903690830161272c565b906040840191825260608101359067ffffffffffffffff8211611497576110c79136910161272c565b6060840190815291516111325750509067ffffffffffffffff82816001945116885260036020526111058460408a206110ff816127f5565b016127f5565b51167f5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d8599168780a25b0161101b565b90918151511561146f57815167ffffffffffffffff82511689526003602052604089209080519067ffffffffffffffff82116113a35761117c82611176855461239b565b85612aaa565b602090601f83116001146113d0576111c892918c9183610dda5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b90555b825167ffffffffffffffff82511689526003602052600160408a20019080519067ffffffffffffffff82116113a35761120882611176855461239b565b6020908b601f84116001146112cd57947f1ced5bcae649ed29cebfa0010298ad6794bf3822e8cb754a6eee5353a9a872129461128d856112c5966112ab9667ffffffffffffffff9660019e9d9c9b92610dda5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b90555b51169351945194602060405192828480945193849201612335565b810103902093604051918291602083526020830190612358565b0390a361112c565b50838c52818c2091907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe084168d5b81811061138b5750946001856112ab9567ffffffffffffffff95839d9c9b9a956112c5997f1ced5bcae649ed29cebfa0010298ad6794bf3822e8cb754a6eee5353a9a872129b10611354575b505050811b019055611290565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c19169055388080611347565b929360206001819287860151815501950193016112fb565b60248b7f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b838c52818c2091907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe084168d5b8181106114575750908460019594939210611420575b505050811b0190556111cb565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c19169055388080611413565b929360206001819287860151815501950193016113fd565b6004887f8579befe000000000000000000000000000000000000000000000000000000008152fd5b8980fd5b8880fd5b8780fd5b6024887f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b8580fd5b8480f35b503461029357807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102935773ffffffffffffffffffffffffffffffffffffffff6020915416604051908152f35b50346102935760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610293576115976115916001604067ffffffffffffffff611574612227565b1694858152600360205281812095815260036020522001926123ee565b916123ee565b906115a7604051928392836124b0565b0390f35b5034610293576115ba3661266e565b73ffffffffffffffffffffffffffffffffffffffff6002541633036116415767ffffffffffffffff6115ee60208301612985565b168083526003602052611604604084205461239b565b15611616575061161390612aef565b80f35b7fd79f2ea4000000000000000000000000000000000000000000000000000000008352600452602482fd5b6024827fd7f7333400000000000000000000000000000000000000000000000000000000815233600452fd5b50346102935760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102935760043567ffffffffffffffff8111610502576116bd90369060040161263d565b9060243567ffffffffffffffff8111610398576116de90369060040161263d565b906116e76137e2565b845b8281106117da57505050825b828110611700578380f35b8067ffffffffffffffff61171f61171a6001948787612945565b612985565b168552600360205261174f6002604087200161174961173f848888612945565b602081019061299a565b90612747565b827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082541617905561178561171a828686612945565b67ffffffffffffffff61179c61173f848888612945565b90816040519283928337810189815203902091167f72d9f73bb7cb11065e15df29d61e803a0eba356d509a7025a6f51ebdea07f9e78780a3016116f5565b8067ffffffffffffffff6117f461171a6001948787612945565b16875260036020526118146002604089200161174961173f848888612945565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00815416905561184861171a828686612945565b67ffffffffffffffff61185f61173f848888612945565b9081604051928392833781018b815203902091167f021290bab0d93f4d9a243bd430e45dd4bc8238451e9abbba6fab4463677dfce98980a3016116e9565b503461029357807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102935760015473ffffffffffffffffffffffffffffffffffffffff8116330361195a577fffffffffffffffffffffffff0000000000000000000000000000000000000000825491338284161784551660015573ffffffffffffffffffffffffffffffffffffffff3391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08380a380f35b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152fd5b503461029357807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102935760206040517f000000000000000000000000000000000000000000000000000000000000000015158152f35b50346102935760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261029357600435611a4e6124fb565b611a566137e2565b611a6d826000526006602052604060002054151590565b15611b3d578183526004602052611a8960046040852001612760565b928281526004602052611a9e60408220612844565b611aa783613a5f565b50805b8451811015611af75780611af173ffffffffffffffffffffffffffffffffffffffff611ad860019489612902565b515116856020611ae8858b612902565b51015191613861565b01611aaa565b50827fd5038100bd3dc9631d3c3f4f61a3e53e9d466f40c47af9897292c7b35e32a95760208473ffffffffffffffffffffffffffffffffffffffff60405191168152a280f35b602483837fb6e78260000000000000000000000000000000000000000000000000000000008252600452fd5b50346102935760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102935760406115a791606060808351611bae81612253565b8381528360208201528285820152828082015201526004358152600460205220611be060046040519261040a84612253565b608082015260405191829182612595565b50346102935760607ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261029357611c296124d8565b611c316124fb565b73ffffffffffffffffffffffffffffffffffffffff60443592611c526137e2565b169081611db457824710611d5657838080808673ffffffffffffffffffffffffffffffffffffffff86165af1611c86612a7a565b5015611cd257602073ffffffffffffffffffffffffffffffffffffffff7f6832d9be2410a86571981e1e60fd4c1f9ea2a1034d6102a2b7d6c5e480adf02e925b6040519586521693a380f35b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603a60248201527f416464726573733a20756e61626c6520746f2073656e642076616c75652c207260448201527f6563697069656e74206d617920686176652072657665727465640000000000006064820152fd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a20696e73756666696369656e742062616c616e63650000006044820152fd5b602073ffffffffffffffffffffffffffffffffffffffff82611df8867f6832d9be2410a86571981e1e60fd4c1f9ea2a1034d6102a2b7d6c5e480adf02e9587613861565b611cc6565b50346102935760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102935760409067ffffffffffffffff611e41612227565b1681526003602052206115976001611e58836123ee565b92016123ee565b503461029357807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261029357602073ffffffffffffffffffffffffffffffffffffffff60085416604051908152f35b503461029357807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261029357611ee86137e2565b7fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff600854166008557f48257dc961b6f792c2b78a080dacfed693b660960a702de21cee364e20270e2f602060405160018152a1610656604051600160208201526020815261062f6040826122ba565b503461029357807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261029357506115a7604051611f986040826122ba565b601681527f50696e67506f6e6744656d6f20312e362e302d646576000000000000000000006020820152604051918291602083526020830190612358565b50346102935760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610293576004358015158091036105025761201b6137e2565b7fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff74ff00000000000000000000000000000000000000006008549260a01b1691161760085580f35b50346102935760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102935761209b612227565b6024359067ffffffffffffffff821161211057366023830112156121105781600401359067ffffffffffffffff821161039857366024838501011161039857916121049160246002604060209767ffffffffffffffff60ff981681526003895220019201612747565b54166040519015158152f35b8280fd5b50346102935760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261029357602061215f6004356000526006602052604060002054151590565b6040519015158152f35b9050346105025760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610502576004357fffffffff00000000000000000000000000000000000000000000000000000000811680910361211057602092507f85572ffb0000000000000000000000000000000000000000000000000000000081149081156121fd575b5015158152f35b7f01ffc9a700000000000000000000000000000000000000000000000000000000915014386121f6565b6004359067ffffffffffffffff82168203610a4057565b359067ffffffffffffffff82168203610a4057565b60a0810190811067ffffffffffffffff82111761226f57604052565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040810190811067ffffffffffffffff82111761226f57604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff82111761226f57604052565b67ffffffffffffffff811161226f57601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b60005b8381106123485750506000910152565b8181015183820152602001612338565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f60209361239481518092818752878088019101612335565b0116010190565b90600182811c921680156123e4575b60208310146123b557565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b91607f16916123aa565b90604051918260008254926124028461239b565b80845293600181169081156124705750600114612429575b50612427925003836122ba565b565b90506000929192526020600020906000915b818310612454575050906020612427928201013861241a565b602091935080600191548385890101520191019091849261243b565b602093506124279592507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091501682840152151560051b8201013861241a565b90916124c76124d593604084526040840190612358565b916020818403910152612358565b90565b6004359073ffffffffffffffffffffffffffffffffffffffff82168203610a4057565b6024359073ffffffffffffffffffffffffffffffffffffffff82168203610a4057565b359073ffffffffffffffffffffffffffffffffffffffff82168203610a4057565b906020808351928381520192019060005b81811061255d5750505090565b8251805173ffffffffffffffffffffffffffffffffffffffff1685526020908101518186015260409094019390920191600101612550565b906124d591602081528151602082015267ffffffffffffffff6020830151166040820152608061260a6125d7604085015160a0606086015260c0850190612358565b60608501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08583030184860152612358565b9201519060a07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08285030191015261253f565b9181601f84011215610a405782359167ffffffffffffffff8311610a40576020808501948460051b010111610a4057565b60207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc820112610a40576004359067ffffffffffffffff8211610a40577ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8260a092030112610a405760040190565b67ffffffffffffffff811161226f5760051b60200190565b929192612701826122fb565b9161270f60405193846122ba565b829481845281830111610a40578281602093846000960137010152565b9080601f83011215610a40578160206124d5933591016126f5565b6020919283604051948593843782019081520301902090565b90815461276c816126dd565b9261277a60405194856122ba565b818452602084019060005260206000206000915b83831061279b5750505050565b600260206001926040516127ae8161229e565b73ffffffffffffffffffffffffffffffffffffffff8654168152848601548382015281520192019201919061278e565b8181106127e9575050565b600081556001016127de565b6127ff815461239b565b9081612809575050565b81601f6000931160011461281b575055565b8183526020832061283791601f0160051c8101906001016127de565b8082528160208120915555565b600490600081556000600182015561285e600282016127f5565b61286a600382016127f5565b01805490600081558161287b575050565b7f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821682036128d3576000908152602081209160011b8201915b8281106128c157505050565b808260029255826001820155016128b5565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b80518210156129165760209160051b010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b91908110156129165760051b810135907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc181360301821215610a40570190565b3567ffffffffffffffff81168103610a405790565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe181360301821215610a40570180359067ffffffffffffffff8211610a4057602001918136038313610a4057565b90357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe182360301811215610a4057016020813591019167ffffffffffffffff8211610a40578136038313610a4057565b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0938186528686013760008582860101520116010190565b3d15612aa5573d90612a8b826122fb565b91612a9960405193846122ba565b82523d6000602084013e565b606090565b9190601f8111612ab957505050565b612427926000526020600020906020601f840160051c83019310612ae5575b601f0160051c01906127de565b9091508190612ad8565b303b15610a40576040517fcf6730f8000000000000000000000000000000000000000000000000000000008152600060206004830152823592836024840152602081019267ffffffffffffffff612b458561223e565b1660448201526040820193612b6e612b5d86856129eb565b60a0606486015260c4850191612a3b565b612bb06060850191612b8083876129eb565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc878403016084880152612a3b565b926080850135937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe186360301851215613193578585016020813591019167ffffffffffffffff821161149b578160061b3603831361149b578381037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc0160a48501528181528891849160200190835b81811061314e57505081929350038183305af1908161313a575b5061310f57612c66612a7a565b95612c7088613bf5565b508786526004602052604086209288845567ffffffffffffffff612c976001860192612985565b167fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000825416179055612ccd60028401918661299a565b9067ffffffffffffffff82116114a357612ceb82611176855461239b565b8790601f831160011461306f57612d369291899183612f935750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b90555b612d4760038301918561299a565b9067ffffffffffffffff821161304257612d6582611176855461239b565b8690601f8311600114612f9e579180612db59260049695948a92612f935750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b90555b0191019081359167ffffffffffffffff8311610398576020018260061b3603811361039857680100000000000000008311612f66578154838355808410612eb4575b50908352602083209083905b838510612e50575050505050612e4b7f55bc02a9ef6f146737edeeb425738006f67f077e7138de3bf84a15bde1a5b56f91604051918291602083526020830190612358565b0390a2565b803573ffffffffffffffffffffffffffffffffffffffff81168091036121105760406001926002927fffffffffffffffffffffffff000000000000000000000000000000000000000087541617865560208101358487015501930194019391612e06565b7f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81168103612f39577f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff84168403612f3957828552602085209060011b8101908460011b015b818110612f275750612dfa565b80866002925586600182015501612f1a565b6024857f4e487b710000000000000000000000000000000000000000000000000000000081526011600452fd5b6024847f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b013590503880610da3565b83885260208820917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08416895b81811061302a575091600193918560049897969410612ff2575b505050811b019055612db8565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88560031b161c19910135169055388080612fe5565b91936020600181928787013581550195019201612fcb565b6024877f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b83895260208920917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe084168a5b8181106130f757509084600195949392106130bf575b505050811b019055612d39565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88560031b161c199101351690553880806130b2565b9193602060018192878701358155019501920161309c565b505050507fdf6958669026659bac75ba986685e11a7d271284989f565f2802522663e9a70f915080a2565b86613147919792976122ba565b9438612c59565b9250925060408060019273ffffffffffffffffffffffffffffffffffffffff6131768861251e565b168152602087013560208201520194019101928992859294612c3f565b8680fd5b6020906131b1928260405194838680955193849201612335565b82019081520301902090565b9192909267ffffffffffffffff8316908160005260036020526131e460406000205461239b565b156136c4576000918252600360205260409182902060075492519273ffffffffffffffffffffffffffffffffffffffff169161324a916001810191906132329061322d87612253565b6123ee565b855260208501528660408501528260608501526123ee565b608083015273ffffffffffffffffffffffffffffffffffffffff6002541690602060405180937f20487ded00000000000000000000000000000000000000000000000000000000825281806132a3888b600484016136f2565b03915afa91821561357157600092613690575b50817f00000000000000000000000000000000000000000000000000000000000000001580613687575b613673575b505060005b855181101561357d5761333173ffffffffffffffffffffffffffffffffffffffff6133158389612902565b5151166020613324848a612902565b51015190309033906138c0565b73ffffffffffffffffffffffffffffffffffffffff6133508288612902565b51511673ffffffffffffffffffffffffffffffffffffffff600754160361337a575b6001016132ea565b73ffffffffffffffffffffffffffffffffffffffff6133998288612902565b5151169073ffffffffffffffffffffffffffffffffffffffff6002541660206133c2838a612902565b510151801580156134be575b1561343a576040517f095ea7b300000000000000000000000000000000000000000000000000000000602082015273ffffffffffffffffffffffffffffffffffffffff9092166024830152604482015260019261343391906109d982606481016109ad565b9050613372565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603660248201527f5361666545524332303a20617070726f76652066726f6d206e6f6e2d7a65726f60448201527f20746f206e6f6e2d7a65726f20616c6c6f77616e6365000000000000000000006064820152fd5b506040517fdd62ed3e00000000000000000000000000000000000000000000000000000000815230600482015273ffffffffffffffffffffffffffffffffffffffff831660248201526020818060448101038173ffffffffffffffffffffffffffffffffffffffff89165afa90811561357157600091613540575b50156133ce565b906020823d8211613569575b81613559602093836122ba565b8101031261029357505138613539565b3d915061354c565b6040513d6000823e3d90fd5b5060209294506135f89373ffffffffffffffffffffffffffffffffffffffff600254169173ffffffffffffffffffffffffffffffffffffffff600754161560001461366b575b6040518096819582947f96f4e9f9000000000000000000000000000000000000000000000000000000008452600484016136f2565b03925af190811561357157600091613639575b507f54791b38f3859327992a1ca0590ad3c0f08feba98d1a4f56ab0dca74d203392a6020604051838152a190565b90506020813d602011613663575b81613654602093836122ba565b81010312610a4057513861360b565b3d9150613647565b5060006135c3565b61368091309033906138c0565b38816132e5565b508115156132e0565b90916020823d6020116136bc575b816136ab602093836122ba565b8101031261029357505190386132b6565b3d915061369e565b507fd79f2ea40000000000000000000000000000000000000000000000000000000060005260045260246000fd5b67ffffffffffffffff6124d5939216815260406020820152608061378f61375c613728855160a0604087015260e0860190612358565b60208601517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0868303016060870152612358565b60408501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0858303018486015261253f565b9273ffffffffffffffffffffffffffffffffffffffff60608201511660a084015201519060c07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc082850301910152612358565b73ffffffffffffffffffffffffffffffffffffffff60005416330361380357565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152fd5b6040517fa9059cbb00000000000000000000000000000000000000000000000000000000602082015273ffffffffffffffffffffffffffffffffffffffff90921660248301526044820192909252612427916109d982606481016109ad565b90919273ffffffffffffffffffffffffffffffffffffffff6124279481604051957f23b872dd0000000000000000000000000000000000000000000000000000000060208801521660248601521660448401526064830152606482526109d96084836122ba565b73ffffffffffffffffffffffffffffffffffffffff61399791169160409260008085519361395587866122ba565b602085527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564602086015260208151910182855af1613991612a7a565b91613c55565b8051806139a357505050565b8160209181010312610a405760200151801590811503610a40576139c45750565b608490517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152fd5b80548210156129165760005260206000200190600090565b6000818152600660205260409020548015613bee577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81018181116128d357600554907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82019182116128d357808203613b7f575b5050506005548015613b50577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01613b0d816005613a47565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82549160031b1b19169055600555600052600660205260006040812055600190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b613bd6613b90613ba1936005613a47565b90549060031b1c9283926005613a47565b81939154907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9060031b92831b921b19161790565b90556000526006602052604060002055388080613ad4565b5050600090565b80600052600660205260406000205415600014613c4f576005546801000000000000000081101561226f57613c36613ba18260018594016005556005613a47565b9055600554906000526006602052604060002055600190565b50600090565b91929015613cd05750815115613c69575090565b3b15613c725790565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152fd5b825190915015613ce35750805190602001fd5b6106f4906040519182917f08c379a000000000000000000000000000000000000000000000000000000000835260206004840152602483019061235856fea164736f6c634300081a000a",
}

var PingPongDemoABI = PingPongDemoMetaData.ABI

var PingPongDemoBin = PingPongDemoMetaData.Bin

func DeployPingPongDemo(auth *bind.TransactOpts, backend bind.ContractBackend, router common.Address, feeToken common.Address) (common.Address, *types.Transaction, *PingPongDemo, error) {
	parsed, err := PingPongDemoMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PingPongDemoBin), backend, router, feeToken)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PingPongDemo{address: address, abi: *parsed, PingPongDemoCaller: PingPongDemoCaller{contract: contract}, PingPongDemoTransactor: PingPongDemoTransactor{contract: contract}, PingPongDemoFilterer: PingPongDemoFilterer{contract: contract}}, nil
}

type PingPongDemo struct {
	address common.Address
	abi     abi.ABI
	PingPongDemoCaller
	PingPongDemoTransactor
	PingPongDemoFilterer
}

type PingPongDemoCaller struct {
	contract *bind.BoundContract
}

type PingPongDemoTransactor struct {
	contract *bind.BoundContract
}

type PingPongDemoFilterer struct {
	contract *bind.BoundContract
}

type PingPongDemoSession struct {
	Contract     *PingPongDemo
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type PingPongDemoCallerSession struct {
	Contract *PingPongDemoCaller
	CallOpts bind.CallOpts
}

type PingPongDemoTransactorSession struct {
	Contract     *PingPongDemoTransactor
	TransactOpts bind.TransactOpts
}

type PingPongDemoRaw struct {
	Contract *PingPongDemo
}

type PingPongDemoCallerRaw struct {
	Contract *PingPongDemoCaller
}

type PingPongDemoTransactorRaw struct {
	Contract *PingPongDemoTransactor
}

func NewPingPongDemo(address common.Address, backend bind.ContractBackend) (*PingPongDemo, error) {
	abi, err := abi.JSON(strings.NewReader(PingPongDemoABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindPingPongDemo(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PingPongDemo{address: address, abi: abi, PingPongDemoCaller: PingPongDemoCaller{contract: contract}, PingPongDemoTransactor: PingPongDemoTransactor{contract: contract}, PingPongDemoFilterer: PingPongDemoFilterer{contract: contract}}, nil
}

func NewPingPongDemoCaller(address common.Address, caller bind.ContractCaller) (*PingPongDemoCaller, error) {
	contract, err := bindPingPongDemo(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PingPongDemoCaller{contract: contract}, nil
}

func NewPingPongDemoTransactor(address common.Address, transactor bind.ContractTransactor) (*PingPongDemoTransactor, error) {
	contract, err := bindPingPongDemo(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PingPongDemoTransactor{contract: contract}, nil
}

func NewPingPongDemoFilterer(address common.Address, filterer bind.ContractFilterer) (*PingPongDemoFilterer, error) {
	contract, err := bindPingPongDemo(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PingPongDemoFilterer{contract: contract}, nil
}

func bindPingPongDemo(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PingPongDemoMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_PingPongDemo *PingPongDemoRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PingPongDemo.Contract.PingPongDemoCaller.contract.Call(opts, result, method, params...)
}

func (_PingPongDemo *PingPongDemoRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PingPongDemo.Contract.PingPongDemoTransactor.contract.Transfer(opts)
}

func (_PingPongDemo *PingPongDemoRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PingPongDemo.Contract.PingPongDemoTransactor.contract.Transact(opts, method, params...)
}

func (_PingPongDemo *PingPongDemoCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PingPongDemo.Contract.contract.Call(opts, result, method, params...)
}

func (_PingPongDemo *PingPongDemoTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PingPongDemo.Contract.contract.Transfer(opts)
}

func (_PingPongDemo *PingPongDemoTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PingPongDemo.Contract.contract.Transact(opts, method, params...)
}

func (_PingPongDemo *PingPongDemoCaller) GetCounterpartAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PingPongDemo.contract.Call(opts, &out, "getCounterpartAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_PingPongDemo *PingPongDemoSession) GetCounterpartAddress() (common.Address, error) {
	return _PingPongDemo.Contract.GetCounterpartAddress(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCallerSession) GetCounterpartAddress() (common.Address, error) {
	return _PingPongDemo.Contract.GetCounterpartAddress(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCaller) GetCounterpartChainSelector(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _PingPongDemo.contract.Call(opts, &out, "getCounterpartChainSelector")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_PingPongDemo *PingPongDemoSession) GetCounterpartChainSelector() (uint64, error) {
	return _PingPongDemo.Contract.GetCounterpartChainSelector(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCallerSession) GetCounterpartChainSelector() (uint64, error) {
	return _PingPongDemo.Contract.GetCounterpartChainSelector(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCaller) GetFeeToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PingPongDemo.contract.Call(opts, &out, "getFeeToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_PingPongDemo *PingPongDemoSession) GetFeeToken() (common.Address, error) {
	return _PingPongDemo.Contract.GetFeeToken(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCallerSession) GetFeeToken() (common.Address, error) {
	return _PingPongDemo.Contract.GetFeeToken(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCaller) GetMessageContents(opts *bind.CallOpts, messageId [32]byte) (ClientAny2EVMMessage, error) {
	var out []interface{}
	err := _PingPongDemo.contract.Call(opts, &out, "getMessageContents", messageId)

	if err != nil {
		return *new(ClientAny2EVMMessage), err
	}

	out0 := *abi.ConvertType(out[0], new(ClientAny2EVMMessage)).(*ClientAny2EVMMessage)

	return out0, err

}

func (_PingPongDemo *PingPongDemoSession) GetMessageContents(messageId [32]byte) (ClientAny2EVMMessage, error) {
	return _PingPongDemo.Contract.GetMessageContents(&_PingPongDemo.CallOpts, messageId)
}

func (_PingPongDemo *PingPongDemoCallerSession) GetMessageContents(messageId [32]byte) (ClientAny2EVMMessage, error) {
	return _PingPongDemo.Contract.GetMessageContents(&_PingPongDemo.CallOpts, messageId)
}

func (_PingPongDemo *PingPongDemoCaller) GetOutOfOrderExecution(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _PingPongDemo.contract.Call(opts, &out, "getOutOfOrderExecution")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_PingPongDemo *PingPongDemoSession) GetOutOfOrderExecution() (bool, error) {
	return _PingPongDemo.Contract.GetOutOfOrderExecution(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCallerSession) GetOutOfOrderExecution() (bool, error) {
	return _PingPongDemo.Contract.GetOutOfOrderExecution(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCaller) GetRemoteChainConfig(opts *bind.CallOpts, remoteChainSelector uint64) (GetRemoteChainConfig,

	error) {
	var out []interface{}
	err := _PingPongDemo.contract.Call(opts, &out, "getRemoteChainConfig", remoteChainSelector)

	outstruct := new(GetRemoteChainConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Recipient = *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	outstruct.ExtraArgsBytes = *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return *outstruct, err

}

func (_PingPongDemo *PingPongDemoSession) GetRemoteChainConfig(remoteChainSelector uint64) (GetRemoteChainConfig,

	error) {
	return _PingPongDemo.Contract.GetRemoteChainConfig(&_PingPongDemo.CallOpts, remoteChainSelector)
}

func (_PingPongDemo *PingPongDemoCallerSession) GetRemoteChainConfig(remoteChainSelector uint64) (GetRemoteChainConfig,

	error) {
	return _PingPongDemo.Contract.GetRemoteChainConfig(&_PingPongDemo.CallOpts, remoteChainSelector)
}

func (_PingPongDemo *PingPongDemoCaller) GetRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PingPongDemo.contract.Call(opts, &out, "getRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_PingPongDemo *PingPongDemoSession) GetRouter() (common.Address, error) {
	return _PingPongDemo.Contract.GetRouter(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCallerSession) GetRouter() (common.Address, error) {
	return _PingPongDemo.Contract.GetRouter(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCaller) IsApprovedSender(opts *bind.CallOpts, sourceChainSelector uint64, senderAddr []byte) (bool, error) {
	var out []interface{}
	err := _PingPongDemo.contract.Call(opts, &out, "isApprovedSender", sourceChainSelector, senderAddr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_PingPongDemo *PingPongDemoSession) IsApprovedSender(sourceChainSelector uint64, senderAddr []byte) (bool, error) {
	return _PingPongDemo.Contract.IsApprovedSender(&_PingPongDemo.CallOpts, sourceChainSelector, senderAddr)
}

func (_PingPongDemo *PingPongDemoCallerSession) IsApprovedSender(sourceChainSelector uint64, senderAddr []byte) (bool, error) {
	return _PingPongDemo.Contract.IsApprovedSender(&_PingPongDemo.CallOpts, sourceChainSelector, senderAddr)
}

func (_PingPongDemo *PingPongDemoCaller) IsFailedMessage(opts *bind.CallOpts, messageId [32]byte) (bool, error) {
	var out []interface{}
	err := _PingPongDemo.contract.Call(opts, &out, "isFailedMessage", messageId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_PingPongDemo *PingPongDemoSession) IsFailedMessage(messageId [32]byte) (bool, error) {
	return _PingPongDemo.Contract.IsFailedMessage(&_PingPongDemo.CallOpts, messageId)
}

func (_PingPongDemo *PingPongDemoCallerSession) IsFailedMessage(messageId [32]byte) (bool, error) {
	return _PingPongDemo.Contract.IsFailedMessage(&_PingPongDemo.CallOpts, messageId)
}

func (_PingPongDemo *PingPongDemoCaller) IsPaused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _PingPongDemo.contract.Call(opts, &out, "isPaused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_PingPongDemo *PingPongDemoSession) IsPaused() (bool, error) {
	return _PingPongDemo.Contract.IsPaused(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCallerSession) IsPaused() (bool, error) {
	return _PingPongDemo.Contract.IsPaused(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PingPongDemo.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_PingPongDemo *PingPongDemoSession) Owner() (common.Address, error) {
	return _PingPongDemo.Contract.Owner(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCallerSession) Owner() (common.Address, error) {
	return _PingPongDemo.Contract.Owner(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCaller) SChainConfigs(opts *bind.CallOpts, destChainSelector uint64) (SChainConfigs,

	error) {
	var out []interface{}
	err := _PingPongDemo.contract.Call(opts, &out, "s_chainConfigs", destChainSelector)

	outstruct := new(SChainConfigs)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Recipient = *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	outstruct.ExtraArgsBytes = *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return *outstruct, err

}

func (_PingPongDemo *PingPongDemoSession) SChainConfigs(destChainSelector uint64) (SChainConfigs,

	error) {
	return _PingPongDemo.Contract.SChainConfigs(&_PingPongDemo.CallOpts, destChainSelector)
}

func (_PingPongDemo *PingPongDemoCallerSession) SChainConfigs(destChainSelector uint64) (SChainConfigs,

	error) {
	return _PingPongDemo.Contract.SChainConfigs(&_PingPongDemo.CallOpts, destChainSelector)
}

func (_PingPongDemo *PingPongDemoCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _PingPongDemo.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_PingPongDemo *PingPongDemoSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _PingPongDemo.Contract.SupportsInterface(&_PingPongDemo.CallOpts, interfaceId)
}

func (_PingPongDemo *PingPongDemoCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _PingPongDemo.Contract.SupportsInterface(&_PingPongDemo.CallOpts, interfaceId)
}

func (_PingPongDemo *PingPongDemoCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _PingPongDemo.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_PingPongDemo *PingPongDemoSession) TypeAndVersion() (string, error) {
	return _PingPongDemo.Contract.TypeAndVersion(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCallerSession) TypeAndVersion() (string, error) {
	return _PingPongDemo.Contract.TypeAndVersion(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCaller) UsePreFundedFeeTokens(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _PingPongDemo.contract.Call(opts, &out, "usePreFundedFeeTokens")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_PingPongDemo *PingPongDemoSession) UsePreFundedFeeTokens() (bool, error) {
	return _PingPongDemo.Contract.UsePreFundedFeeTokens(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoCallerSession) UsePreFundedFeeTokens() (bool, error) {
	return _PingPongDemo.Contract.UsePreFundedFeeTokens(&_PingPongDemo.CallOpts)
}

func (_PingPongDemo *PingPongDemoTransactor) AbandonFailedMessage(opts *bind.TransactOpts, messageId [32]byte, receiver common.Address) (*types.Transaction, error) {
	return _PingPongDemo.contract.Transact(opts, "abandonFailedMessage", messageId, receiver)
}

func (_PingPongDemo *PingPongDemoSession) AbandonFailedMessage(messageId [32]byte, receiver common.Address) (*types.Transaction, error) {
	return _PingPongDemo.Contract.AbandonFailedMessage(&_PingPongDemo.TransactOpts, messageId, receiver)
}

func (_PingPongDemo *PingPongDemoTransactorSession) AbandonFailedMessage(messageId [32]byte, receiver common.Address) (*types.Transaction, error) {
	return _PingPongDemo.Contract.AbandonFailedMessage(&_PingPongDemo.TransactOpts, messageId, receiver)
}

func (_PingPongDemo *PingPongDemoTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PingPongDemo.contract.Transact(opts, "acceptOwnership")
}

func (_PingPongDemo *PingPongDemoSession) AcceptOwnership() (*types.Transaction, error) {
	return _PingPongDemo.Contract.AcceptOwnership(&_PingPongDemo.TransactOpts)
}

func (_PingPongDemo *PingPongDemoTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _PingPongDemo.Contract.AcceptOwnership(&_PingPongDemo.TransactOpts)
}

func (_PingPongDemo *PingPongDemoTransactor) ApplyChainUpdates(opts *bind.TransactOpts, chains []CCIPBaseChainUpdate) (*types.Transaction, error) {
	return _PingPongDemo.contract.Transact(opts, "applyChainUpdates", chains)
}

func (_PingPongDemo *PingPongDemoSession) ApplyChainUpdates(chains []CCIPBaseChainUpdate) (*types.Transaction, error) {
	return _PingPongDemo.Contract.ApplyChainUpdates(&_PingPongDemo.TransactOpts, chains)
}

func (_PingPongDemo *PingPongDemoTransactorSession) ApplyChainUpdates(chains []CCIPBaseChainUpdate) (*types.Transaction, error) {
	return _PingPongDemo.Contract.ApplyChainUpdates(&_PingPongDemo.TransactOpts, chains)
}

func (_PingPongDemo *PingPongDemoTransactor) CcipReceive(opts *bind.TransactOpts, message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _PingPongDemo.contract.Transact(opts, "ccipReceive", message)
}

func (_PingPongDemo *PingPongDemoSession) CcipReceive(message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _PingPongDemo.Contract.CcipReceive(&_PingPongDemo.TransactOpts, message)
}

func (_PingPongDemo *PingPongDemoTransactorSession) CcipReceive(message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _PingPongDemo.Contract.CcipReceive(&_PingPongDemo.TransactOpts, message)
}

func (_PingPongDemo *PingPongDemoTransactor) CcipSend(opts *bind.TransactOpts, destChainSelector uint64, tokenAmounts []ClientEVMTokenAmount, data []byte) (*types.Transaction, error) {
	return _PingPongDemo.contract.Transact(opts, "ccipSend", destChainSelector, tokenAmounts, data)
}

func (_PingPongDemo *PingPongDemoSession) CcipSend(destChainSelector uint64, tokenAmounts []ClientEVMTokenAmount, data []byte) (*types.Transaction, error) {
	return _PingPongDemo.Contract.CcipSend(&_PingPongDemo.TransactOpts, destChainSelector, tokenAmounts, data)
}

func (_PingPongDemo *PingPongDemoTransactorSession) CcipSend(destChainSelector uint64, tokenAmounts []ClientEVMTokenAmount, data []byte) (*types.Transaction, error) {
	return _PingPongDemo.Contract.CcipSend(&_PingPongDemo.TransactOpts, destChainSelector, tokenAmounts, data)
}

func (_PingPongDemo *PingPongDemoTransactor) ProcessMessage(opts *bind.TransactOpts, message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _PingPongDemo.contract.Transact(opts, "processMessage", message)
}

func (_PingPongDemo *PingPongDemoSession) ProcessMessage(message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _PingPongDemo.Contract.ProcessMessage(&_PingPongDemo.TransactOpts, message)
}

func (_PingPongDemo *PingPongDemoTransactorSession) ProcessMessage(message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _PingPongDemo.Contract.ProcessMessage(&_PingPongDemo.TransactOpts, message)
}

func (_PingPongDemo *PingPongDemoTransactor) RetryFailedMessage(opts *bind.TransactOpts, messageId [32]byte) (*types.Transaction, error) {
	return _PingPongDemo.contract.Transact(opts, "retryFailedMessage", messageId)
}

func (_PingPongDemo *PingPongDemoSession) RetryFailedMessage(messageId [32]byte) (*types.Transaction, error) {
	return _PingPongDemo.Contract.RetryFailedMessage(&_PingPongDemo.TransactOpts, messageId)
}

func (_PingPongDemo *PingPongDemoTransactorSession) RetryFailedMessage(messageId [32]byte) (*types.Transaction, error) {
	return _PingPongDemo.Contract.RetryFailedMessage(&_PingPongDemo.TransactOpts, messageId)
}

func (_PingPongDemo *PingPongDemoTransactor) SetCounterpart(opts *bind.TransactOpts, counterpartChainSelector uint64, counterpartAddress common.Address) (*types.Transaction, error) {
	return _PingPongDemo.contract.Transact(opts, "setCounterpart", counterpartChainSelector, counterpartAddress)
}

func (_PingPongDemo *PingPongDemoSession) SetCounterpart(counterpartChainSelector uint64, counterpartAddress common.Address) (*types.Transaction, error) {
	return _PingPongDemo.Contract.SetCounterpart(&_PingPongDemo.TransactOpts, counterpartChainSelector, counterpartAddress)
}

func (_PingPongDemo *PingPongDemoTransactorSession) SetCounterpart(counterpartChainSelector uint64, counterpartAddress common.Address) (*types.Transaction, error) {
	return _PingPongDemo.Contract.SetCounterpart(&_PingPongDemo.TransactOpts, counterpartChainSelector, counterpartAddress)
}

func (_PingPongDemo *PingPongDemoTransactor) SetPaused(opts *bind.TransactOpts, pause bool) (*types.Transaction, error) {
	return _PingPongDemo.contract.Transact(opts, "setPaused", pause)
}

func (_PingPongDemo *PingPongDemoSession) SetPaused(pause bool) (*types.Transaction, error) {
	return _PingPongDemo.Contract.SetPaused(&_PingPongDemo.TransactOpts, pause)
}

func (_PingPongDemo *PingPongDemoTransactorSession) SetPaused(pause bool) (*types.Transaction, error) {
	return _PingPongDemo.Contract.SetPaused(&_PingPongDemo.TransactOpts, pause)
}

func (_PingPongDemo *PingPongDemoTransactor) StartPingPong(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PingPongDemo.contract.Transact(opts, "startPingPong")
}

func (_PingPongDemo *PingPongDemoSession) StartPingPong() (*types.Transaction, error) {
	return _PingPongDemo.Contract.StartPingPong(&_PingPongDemo.TransactOpts)
}

func (_PingPongDemo *PingPongDemoTransactorSession) StartPingPong() (*types.Transaction, error) {
	return _PingPongDemo.Contract.StartPingPong(&_PingPongDemo.TransactOpts)
}

func (_PingPongDemo *PingPongDemoTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _PingPongDemo.contract.Transact(opts, "transferOwnership", to)
}

func (_PingPongDemo *PingPongDemoSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _PingPongDemo.Contract.TransferOwnership(&_PingPongDemo.TransactOpts, to)
}

func (_PingPongDemo *PingPongDemoTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _PingPongDemo.Contract.TransferOwnership(&_PingPongDemo.TransactOpts, to)
}

func (_PingPongDemo *PingPongDemoTransactor) UpdateApprovedSenders(opts *bind.TransactOpts, adds []CCIPBaseApprovedSenderUpdate, removes []CCIPBaseApprovedSenderUpdate) (*types.Transaction, error) {
	return _PingPongDemo.contract.Transact(opts, "updateApprovedSenders", adds, removes)
}

func (_PingPongDemo *PingPongDemoSession) UpdateApprovedSenders(adds []CCIPBaseApprovedSenderUpdate, removes []CCIPBaseApprovedSenderUpdate) (*types.Transaction, error) {
	return _PingPongDemo.Contract.UpdateApprovedSenders(&_PingPongDemo.TransactOpts, adds, removes)
}

func (_PingPongDemo *PingPongDemoTransactorSession) UpdateApprovedSenders(adds []CCIPBaseApprovedSenderUpdate, removes []CCIPBaseApprovedSenderUpdate) (*types.Transaction, error) {
	return _PingPongDemo.Contract.UpdateApprovedSenders(&_PingPongDemo.TransactOpts, adds, removes)
}

func (_PingPongDemo *PingPongDemoTransactor) UpdateFeeToken(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _PingPongDemo.contract.Transact(opts, "updateFeeToken", token)
}

func (_PingPongDemo *PingPongDemoSession) UpdateFeeToken(token common.Address) (*types.Transaction, error) {
	return _PingPongDemo.Contract.UpdateFeeToken(&_PingPongDemo.TransactOpts, token)
}

func (_PingPongDemo *PingPongDemoTransactorSession) UpdateFeeToken(token common.Address) (*types.Transaction, error) {
	return _PingPongDemo.Contract.UpdateFeeToken(&_PingPongDemo.TransactOpts, token)
}

func (_PingPongDemo *PingPongDemoTransactor) UpdateRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error) {
	return _PingPongDemo.contract.Transact(opts, "updateRouter", newRouter)
}

func (_PingPongDemo *PingPongDemoSession) UpdateRouter(newRouter common.Address) (*types.Transaction, error) {
	return _PingPongDemo.Contract.UpdateRouter(&_PingPongDemo.TransactOpts, newRouter)
}

func (_PingPongDemo *PingPongDemoTransactorSession) UpdateRouter(newRouter common.Address) (*types.Transaction, error) {
	return _PingPongDemo.Contract.UpdateRouter(&_PingPongDemo.TransactOpts, newRouter)
}

func (_PingPongDemo *PingPongDemoTransactor) WithdrawTokens(opts *bind.TransactOpts, token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PingPongDemo.contract.Transact(opts, "withdrawTokens", token, to, amount)
}

func (_PingPongDemo *PingPongDemoSession) WithdrawTokens(token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PingPongDemo.Contract.WithdrawTokens(&_PingPongDemo.TransactOpts, token, to, amount)
}

func (_PingPongDemo *PingPongDemoTransactorSession) WithdrawTokens(token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PingPongDemo.Contract.WithdrawTokens(&_PingPongDemo.TransactOpts, token, to, amount)
}

func (_PingPongDemo *PingPongDemoTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PingPongDemo.contract.RawTransact(opts, nil)
}

func (_PingPongDemo *PingPongDemoSession) Receive() (*types.Transaction, error) {
	return _PingPongDemo.Contract.Receive(&_PingPongDemo.TransactOpts)
}

func (_PingPongDemo *PingPongDemoTransactorSession) Receive() (*types.Transaction, error) {
	return _PingPongDemo.Contract.Receive(&_PingPongDemo.TransactOpts)
}

type PingPongDemoApprovedSenderAddedIterator struct {
	Event *PingPongDemoApprovedSenderAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PingPongDemoApprovedSenderAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PingPongDemoApprovedSenderAdded)
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
		it.Event = new(PingPongDemoApprovedSenderAdded)
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

func (it *PingPongDemoApprovedSenderAddedIterator) Error() error {
	return it.fail
}

func (it *PingPongDemoApprovedSenderAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PingPongDemoApprovedSenderAdded struct {
	DestChainSelector uint64
	Recipient         common.Hash
	Raw               types.Log
}

func (_PingPongDemo *PingPongDemoFilterer) FilterApprovedSenderAdded(opts *bind.FilterOpts, destChainSelector []uint64, recipient [][]byte) (*PingPongDemoApprovedSenderAddedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _PingPongDemo.contract.FilterLogs(opts, "ApprovedSenderAdded", destChainSelectorRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &PingPongDemoApprovedSenderAddedIterator{contract: _PingPongDemo.contract, event: "ApprovedSenderAdded", logs: logs, sub: sub}, nil
}

func (_PingPongDemo *PingPongDemoFilterer) WatchApprovedSenderAdded(opts *bind.WatchOpts, sink chan<- *PingPongDemoApprovedSenderAdded, destChainSelector []uint64, recipient [][]byte) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _PingPongDemo.contract.WatchLogs(opts, "ApprovedSenderAdded", destChainSelectorRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PingPongDemoApprovedSenderAdded)
				if err := _PingPongDemo.contract.UnpackLog(event, "ApprovedSenderAdded", log); err != nil {
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

func (_PingPongDemo *PingPongDemoFilterer) ParseApprovedSenderAdded(log types.Log) (*PingPongDemoApprovedSenderAdded, error) {
	event := new(PingPongDemoApprovedSenderAdded)
	if err := _PingPongDemo.contract.UnpackLog(event, "ApprovedSenderAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PingPongDemoApprovedSenderRemovedIterator struct {
	Event *PingPongDemoApprovedSenderRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PingPongDemoApprovedSenderRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PingPongDemoApprovedSenderRemoved)
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
		it.Event = new(PingPongDemoApprovedSenderRemoved)
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

func (it *PingPongDemoApprovedSenderRemovedIterator) Error() error {
	return it.fail
}

func (it *PingPongDemoApprovedSenderRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PingPongDemoApprovedSenderRemoved struct {
	DestChainSelector uint64
	Recipient         common.Hash
	Raw               types.Log
}

func (_PingPongDemo *PingPongDemoFilterer) FilterApprovedSenderRemoved(opts *bind.FilterOpts, destChainSelector []uint64, recipient [][]byte) (*PingPongDemoApprovedSenderRemovedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _PingPongDemo.contract.FilterLogs(opts, "ApprovedSenderRemoved", destChainSelectorRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &PingPongDemoApprovedSenderRemovedIterator{contract: _PingPongDemo.contract, event: "ApprovedSenderRemoved", logs: logs, sub: sub}, nil
}

func (_PingPongDemo *PingPongDemoFilterer) WatchApprovedSenderRemoved(opts *bind.WatchOpts, sink chan<- *PingPongDemoApprovedSenderRemoved, destChainSelector []uint64, recipient [][]byte) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _PingPongDemo.contract.WatchLogs(opts, "ApprovedSenderRemoved", destChainSelectorRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PingPongDemoApprovedSenderRemoved)
				if err := _PingPongDemo.contract.UnpackLog(event, "ApprovedSenderRemoved", log); err != nil {
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

func (_PingPongDemo *PingPongDemoFilterer) ParseApprovedSenderRemoved(log types.Log) (*PingPongDemoApprovedSenderRemoved, error) {
	event := new(PingPongDemoApprovedSenderRemoved)
	if err := _PingPongDemo.contract.UnpackLog(event, "ApprovedSenderRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PingPongDemoCCIPRouterModifiedIterator struct {
	Event *PingPongDemoCCIPRouterModified

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PingPongDemoCCIPRouterModifiedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PingPongDemoCCIPRouterModified)
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
		it.Event = new(PingPongDemoCCIPRouterModified)
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

func (it *PingPongDemoCCIPRouterModifiedIterator) Error() error {
	return it.fail
}

func (it *PingPongDemoCCIPRouterModifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PingPongDemoCCIPRouterModified struct {
	OldRouter common.Address
	NewRouter common.Address
	Raw       types.Log
}

func (_PingPongDemo *PingPongDemoFilterer) FilterCCIPRouterModified(opts *bind.FilterOpts, oldRouter []common.Address, newRouter []common.Address) (*PingPongDemoCCIPRouterModifiedIterator, error) {

	var oldRouterRule []interface{}
	for _, oldRouterItem := range oldRouter {
		oldRouterRule = append(oldRouterRule, oldRouterItem)
	}
	var newRouterRule []interface{}
	for _, newRouterItem := range newRouter {
		newRouterRule = append(newRouterRule, newRouterItem)
	}

	logs, sub, err := _PingPongDemo.contract.FilterLogs(opts, "CCIPRouterModified", oldRouterRule, newRouterRule)
	if err != nil {
		return nil, err
	}
	return &PingPongDemoCCIPRouterModifiedIterator{contract: _PingPongDemo.contract, event: "CCIPRouterModified", logs: logs, sub: sub}, nil
}

func (_PingPongDemo *PingPongDemoFilterer) WatchCCIPRouterModified(opts *bind.WatchOpts, sink chan<- *PingPongDemoCCIPRouterModified, oldRouter []common.Address, newRouter []common.Address) (event.Subscription, error) {

	var oldRouterRule []interface{}
	for _, oldRouterItem := range oldRouter {
		oldRouterRule = append(oldRouterRule, oldRouterItem)
	}
	var newRouterRule []interface{}
	for _, newRouterItem := range newRouter {
		newRouterRule = append(newRouterRule, newRouterItem)
	}

	logs, sub, err := _PingPongDemo.contract.WatchLogs(opts, "CCIPRouterModified", oldRouterRule, newRouterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PingPongDemoCCIPRouterModified)
				if err := _PingPongDemo.contract.UnpackLog(event, "CCIPRouterModified", log); err != nil {
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

func (_PingPongDemo *PingPongDemoFilterer) ParseCCIPRouterModified(log types.Log) (*PingPongDemoCCIPRouterModified, error) {
	event := new(PingPongDemoCCIPRouterModified)
	if err := _PingPongDemo.contract.UnpackLog(event, "CCIPRouterModified", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PingPongDemoChainAddedIterator struct {
	Event *PingPongDemoChainAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PingPongDemoChainAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PingPongDemoChainAdded)
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
		it.Event = new(PingPongDemoChainAdded)
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

func (it *PingPongDemoChainAddedIterator) Error() error {
	return it.fail
}

func (it *PingPongDemoChainAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PingPongDemoChainAdded struct {
	RemoteChainSelector uint64
	Recipient           common.Hash
	ExtraArgsBytes      []byte
	Raw                 types.Log
}

func (_PingPongDemo *PingPongDemoFilterer) FilterChainAdded(opts *bind.FilterOpts, remoteChainSelector []uint64, recipient [][]byte) (*PingPongDemoChainAddedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _PingPongDemo.contract.FilterLogs(opts, "ChainAdded", remoteChainSelectorRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &PingPongDemoChainAddedIterator{contract: _PingPongDemo.contract, event: "ChainAdded", logs: logs, sub: sub}, nil
}

func (_PingPongDemo *PingPongDemoFilterer) WatchChainAdded(opts *bind.WatchOpts, sink chan<- *PingPongDemoChainAdded, remoteChainSelector []uint64, recipient [][]byte) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _PingPongDemo.contract.WatchLogs(opts, "ChainAdded", remoteChainSelectorRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PingPongDemoChainAdded)
				if err := _PingPongDemo.contract.UnpackLog(event, "ChainAdded", log); err != nil {
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

func (_PingPongDemo *PingPongDemoFilterer) ParseChainAdded(log types.Log) (*PingPongDemoChainAdded, error) {
	event := new(PingPongDemoChainAdded)
	if err := _PingPongDemo.contract.UnpackLog(event, "ChainAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PingPongDemoChainRemovedIterator struct {
	Event *PingPongDemoChainRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PingPongDemoChainRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PingPongDemoChainRemoved)
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
		it.Event = new(PingPongDemoChainRemoved)
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

func (it *PingPongDemoChainRemovedIterator) Error() error {
	return it.fail
}

func (it *PingPongDemoChainRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PingPongDemoChainRemoved struct {
	RemoveChainSelector uint64
	Raw                 types.Log
}

func (_PingPongDemo *PingPongDemoFilterer) FilterChainRemoved(opts *bind.FilterOpts, removeChainSelector []uint64) (*PingPongDemoChainRemovedIterator, error) {

	var removeChainSelectorRule []interface{}
	for _, removeChainSelectorItem := range removeChainSelector {
		removeChainSelectorRule = append(removeChainSelectorRule, removeChainSelectorItem)
	}

	logs, sub, err := _PingPongDemo.contract.FilterLogs(opts, "ChainRemoved", removeChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &PingPongDemoChainRemovedIterator{contract: _PingPongDemo.contract, event: "ChainRemoved", logs: logs, sub: sub}, nil
}

func (_PingPongDemo *PingPongDemoFilterer) WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *PingPongDemoChainRemoved, removeChainSelector []uint64) (event.Subscription, error) {

	var removeChainSelectorRule []interface{}
	for _, removeChainSelectorItem := range removeChainSelector {
		removeChainSelectorRule = append(removeChainSelectorRule, removeChainSelectorItem)
	}

	logs, sub, err := _PingPongDemo.contract.WatchLogs(opts, "ChainRemoved", removeChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PingPongDemoChainRemoved)
				if err := _PingPongDemo.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
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

func (_PingPongDemo *PingPongDemoFilterer) ParseChainRemoved(log types.Log) (*PingPongDemoChainRemoved, error) {
	event := new(PingPongDemoChainRemoved)
	if err := _PingPongDemo.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PingPongDemoFeeTokenUpdatedIterator struct {
	Event *PingPongDemoFeeTokenUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PingPongDemoFeeTokenUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PingPongDemoFeeTokenUpdated)
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
		it.Event = new(PingPongDemoFeeTokenUpdated)
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

func (it *PingPongDemoFeeTokenUpdatedIterator) Error() error {
	return it.fail
}

func (it *PingPongDemoFeeTokenUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PingPongDemoFeeTokenUpdated struct {
	OldFeeToken common.Address
	NewFeeToken common.Address
	Raw         types.Log
}

func (_PingPongDemo *PingPongDemoFilterer) FilterFeeTokenUpdated(opts *bind.FilterOpts) (*PingPongDemoFeeTokenUpdatedIterator, error) {

	logs, sub, err := _PingPongDemo.contract.FilterLogs(opts, "FeeTokenUpdated")
	if err != nil {
		return nil, err
	}
	return &PingPongDemoFeeTokenUpdatedIterator{contract: _PingPongDemo.contract, event: "FeeTokenUpdated", logs: logs, sub: sub}, nil
}

func (_PingPongDemo *PingPongDemoFilterer) WatchFeeTokenUpdated(opts *bind.WatchOpts, sink chan<- *PingPongDemoFeeTokenUpdated) (event.Subscription, error) {

	logs, sub, err := _PingPongDemo.contract.WatchLogs(opts, "FeeTokenUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PingPongDemoFeeTokenUpdated)
				if err := _PingPongDemo.contract.UnpackLog(event, "FeeTokenUpdated", log); err != nil {
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

func (_PingPongDemo *PingPongDemoFilterer) ParseFeeTokenUpdated(log types.Log) (*PingPongDemoFeeTokenUpdated, error) {
	event := new(PingPongDemoFeeTokenUpdated)
	if err := _PingPongDemo.contract.UnpackLog(event, "FeeTokenUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PingPongDemoMessageAbandonedIterator struct {
	Event *PingPongDemoMessageAbandoned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PingPongDemoMessageAbandonedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PingPongDemoMessageAbandoned)
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
		it.Event = new(PingPongDemoMessageAbandoned)
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

func (it *PingPongDemoMessageAbandonedIterator) Error() error {
	return it.fail
}

func (it *PingPongDemoMessageAbandonedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PingPongDemoMessageAbandoned struct {
	MessageId     [32]byte
	TokenReceiver common.Address
	Raw           types.Log
}

func (_PingPongDemo *PingPongDemoFilterer) FilterMessageAbandoned(opts *bind.FilterOpts, messageId [][32]byte) (*PingPongDemoMessageAbandonedIterator, error) {

	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}

	logs, sub, err := _PingPongDemo.contract.FilterLogs(opts, "MessageAbandoned", messageIdRule)
	if err != nil {
		return nil, err
	}
	return &PingPongDemoMessageAbandonedIterator{contract: _PingPongDemo.contract, event: "MessageAbandoned", logs: logs, sub: sub}, nil
}

func (_PingPongDemo *PingPongDemoFilterer) WatchMessageAbandoned(opts *bind.WatchOpts, sink chan<- *PingPongDemoMessageAbandoned, messageId [][32]byte) (event.Subscription, error) {

	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}

	logs, sub, err := _PingPongDemo.contract.WatchLogs(opts, "MessageAbandoned", messageIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PingPongDemoMessageAbandoned)
				if err := _PingPongDemo.contract.UnpackLog(event, "MessageAbandoned", log); err != nil {
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

func (_PingPongDemo *PingPongDemoFilterer) ParseMessageAbandoned(log types.Log) (*PingPongDemoMessageAbandoned, error) {
	event := new(PingPongDemoMessageAbandoned)
	if err := _PingPongDemo.contract.UnpackLog(event, "MessageAbandoned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PingPongDemoMessageFailedIterator struct {
	Event *PingPongDemoMessageFailed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PingPongDemoMessageFailedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PingPongDemoMessageFailed)
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
		it.Event = new(PingPongDemoMessageFailed)
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

func (it *PingPongDemoMessageFailedIterator) Error() error {
	return it.fail
}

func (it *PingPongDemoMessageFailedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PingPongDemoMessageFailed struct {
	MessageId [32]byte
	Reason    []byte
	Raw       types.Log
}

func (_PingPongDemo *PingPongDemoFilterer) FilterMessageFailed(opts *bind.FilterOpts, messageId [][32]byte) (*PingPongDemoMessageFailedIterator, error) {

	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}

	logs, sub, err := _PingPongDemo.contract.FilterLogs(opts, "MessageFailed", messageIdRule)
	if err != nil {
		return nil, err
	}
	return &PingPongDemoMessageFailedIterator{contract: _PingPongDemo.contract, event: "MessageFailed", logs: logs, sub: sub}, nil
}

func (_PingPongDemo *PingPongDemoFilterer) WatchMessageFailed(opts *bind.WatchOpts, sink chan<- *PingPongDemoMessageFailed, messageId [][32]byte) (event.Subscription, error) {

	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}

	logs, sub, err := _PingPongDemo.contract.WatchLogs(opts, "MessageFailed", messageIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PingPongDemoMessageFailed)
				if err := _PingPongDemo.contract.UnpackLog(event, "MessageFailed", log); err != nil {
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

func (_PingPongDemo *PingPongDemoFilterer) ParseMessageFailed(log types.Log) (*PingPongDemoMessageFailed, error) {
	event := new(PingPongDemoMessageFailed)
	if err := _PingPongDemo.contract.UnpackLog(event, "MessageFailed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PingPongDemoMessageRecoveredIterator struct {
	Event *PingPongDemoMessageRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PingPongDemoMessageRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PingPongDemoMessageRecovered)
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
		it.Event = new(PingPongDemoMessageRecovered)
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

func (it *PingPongDemoMessageRecoveredIterator) Error() error {
	return it.fail
}

func (it *PingPongDemoMessageRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PingPongDemoMessageRecovered struct {
	MessageId [32]byte
	Raw       types.Log
}

func (_PingPongDemo *PingPongDemoFilterer) FilterMessageRecovered(opts *bind.FilterOpts, messageId [][32]byte) (*PingPongDemoMessageRecoveredIterator, error) {

	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}

	logs, sub, err := _PingPongDemo.contract.FilterLogs(opts, "MessageRecovered", messageIdRule)
	if err != nil {
		return nil, err
	}
	return &PingPongDemoMessageRecoveredIterator{contract: _PingPongDemo.contract, event: "MessageRecovered", logs: logs, sub: sub}, nil
}

func (_PingPongDemo *PingPongDemoFilterer) WatchMessageRecovered(opts *bind.WatchOpts, sink chan<- *PingPongDemoMessageRecovered, messageId [][32]byte) (event.Subscription, error) {

	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}

	logs, sub, err := _PingPongDemo.contract.WatchLogs(opts, "MessageRecovered", messageIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PingPongDemoMessageRecovered)
				if err := _PingPongDemo.contract.UnpackLog(event, "MessageRecovered", log); err != nil {
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

func (_PingPongDemo *PingPongDemoFilterer) ParseMessageRecovered(log types.Log) (*PingPongDemoMessageRecovered, error) {
	event := new(PingPongDemoMessageRecovered)
	if err := _PingPongDemo.contract.UnpackLog(event, "MessageRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PingPongDemoMessageSentIterator struct {
	Event *PingPongDemoMessageSent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PingPongDemoMessageSentIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PingPongDemoMessageSent)
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
		it.Event = new(PingPongDemoMessageSent)
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

func (it *PingPongDemoMessageSentIterator) Error() error {
	return it.fail
}

func (it *PingPongDemoMessageSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PingPongDemoMessageSent struct {
	MessageId [32]byte
	Raw       types.Log
}

func (_PingPongDemo *PingPongDemoFilterer) FilterMessageSent(opts *bind.FilterOpts) (*PingPongDemoMessageSentIterator, error) {

	logs, sub, err := _PingPongDemo.contract.FilterLogs(opts, "MessageSent")
	if err != nil {
		return nil, err
	}
	return &PingPongDemoMessageSentIterator{contract: _PingPongDemo.contract, event: "MessageSent", logs: logs, sub: sub}, nil
}

func (_PingPongDemo *PingPongDemoFilterer) WatchMessageSent(opts *bind.WatchOpts, sink chan<- *PingPongDemoMessageSent) (event.Subscription, error) {

	logs, sub, err := _PingPongDemo.contract.WatchLogs(opts, "MessageSent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PingPongDemoMessageSent)
				if err := _PingPongDemo.contract.UnpackLog(event, "MessageSent", log); err != nil {
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

func (_PingPongDemo *PingPongDemoFilterer) ParseMessageSent(log types.Log) (*PingPongDemoMessageSent, error) {
	event := new(PingPongDemoMessageSent)
	if err := _PingPongDemo.contract.UnpackLog(event, "MessageSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PingPongDemoMessageSucceededIterator struct {
	Event *PingPongDemoMessageSucceeded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PingPongDemoMessageSucceededIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PingPongDemoMessageSucceeded)
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
		it.Event = new(PingPongDemoMessageSucceeded)
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

func (it *PingPongDemoMessageSucceededIterator) Error() error {
	return it.fail
}

func (it *PingPongDemoMessageSucceededIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PingPongDemoMessageSucceeded struct {
	MessageId [32]byte
	Raw       types.Log
}

func (_PingPongDemo *PingPongDemoFilterer) FilterMessageSucceeded(opts *bind.FilterOpts, messageId [][32]byte) (*PingPongDemoMessageSucceededIterator, error) {

	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}

	logs, sub, err := _PingPongDemo.contract.FilterLogs(opts, "MessageSucceeded", messageIdRule)
	if err != nil {
		return nil, err
	}
	return &PingPongDemoMessageSucceededIterator{contract: _PingPongDemo.contract, event: "MessageSucceeded", logs: logs, sub: sub}, nil
}

func (_PingPongDemo *PingPongDemoFilterer) WatchMessageSucceeded(opts *bind.WatchOpts, sink chan<- *PingPongDemoMessageSucceeded, messageId [][32]byte) (event.Subscription, error) {

	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}

	logs, sub, err := _PingPongDemo.contract.WatchLogs(opts, "MessageSucceeded", messageIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PingPongDemoMessageSucceeded)
				if err := _PingPongDemo.contract.UnpackLog(event, "MessageSucceeded", log); err != nil {
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

func (_PingPongDemo *PingPongDemoFilterer) ParseMessageSucceeded(log types.Log) (*PingPongDemoMessageSucceeded, error) {
	event := new(PingPongDemoMessageSucceeded)
	if err := _PingPongDemo.contract.UnpackLog(event, "MessageSucceeded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PingPongDemoOutOfOrderExecutionChangeIterator struct {
	Event *PingPongDemoOutOfOrderExecutionChange

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PingPongDemoOutOfOrderExecutionChangeIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PingPongDemoOutOfOrderExecutionChange)
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
		it.Event = new(PingPongDemoOutOfOrderExecutionChange)
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

func (it *PingPongDemoOutOfOrderExecutionChangeIterator) Error() error {
	return it.fail
}

func (it *PingPongDemoOutOfOrderExecutionChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PingPongDemoOutOfOrderExecutionChange struct {
	IsOutOfOrder bool
	Raw          types.Log
}

func (_PingPongDemo *PingPongDemoFilterer) FilterOutOfOrderExecutionChange(opts *bind.FilterOpts) (*PingPongDemoOutOfOrderExecutionChangeIterator, error) {

	logs, sub, err := _PingPongDemo.contract.FilterLogs(opts, "OutOfOrderExecutionChange")
	if err != nil {
		return nil, err
	}
	return &PingPongDemoOutOfOrderExecutionChangeIterator{contract: _PingPongDemo.contract, event: "OutOfOrderExecutionChange", logs: logs, sub: sub}, nil
}

func (_PingPongDemo *PingPongDemoFilterer) WatchOutOfOrderExecutionChange(opts *bind.WatchOpts, sink chan<- *PingPongDemoOutOfOrderExecutionChange) (event.Subscription, error) {

	logs, sub, err := _PingPongDemo.contract.WatchLogs(opts, "OutOfOrderExecutionChange")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PingPongDemoOutOfOrderExecutionChange)
				if err := _PingPongDemo.contract.UnpackLog(event, "OutOfOrderExecutionChange", log); err != nil {
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

func (_PingPongDemo *PingPongDemoFilterer) ParseOutOfOrderExecutionChange(log types.Log) (*PingPongDemoOutOfOrderExecutionChange, error) {
	event := new(PingPongDemoOutOfOrderExecutionChange)
	if err := _PingPongDemo.contract.UnpackLog(event, "OutOfOrderExecutionChange", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PingPongDemoOwnershipTransferRequestedIterator struct {
	Event *PingPongDemoOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PingPongDemoOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PingPongDemoOwnershipTransferRequested)
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
		it.Event = new(PingPongDemoOwnershipTransferRequested)
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

func (it *PingPongDemoOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *PingPongDemoOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PingPongDemoOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_PingPongDemo *PingPongDemoFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PingPongDemoOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PingPongDemo.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &PingPongDemoOwnershipTransferRequestedIterator{contract: _PingPongDemo.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_PingPongDemo *PingPongDemoFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *PingPongDemoOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PingPongDemo.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PingPongDemoOwnershipTransferRequested)
				if err := _PingPongDemo.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_PingPongDemo *PingPongDemoFilterer) ParseOwnershipTransferRequested(log types.Log) (*PingPongDemoOwnershipTransferRequested, error) {
	event := new(PingPongDemoOwnershipTransferRequested)
	if err := _PingPongDemo.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PingPongDemoOwnershipTransferredIterator struct {
	Event *PingPongDemoOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PingPongDemoOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PingPongDemoOwnershipTransferred)
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
		it.Event = new(PingPongDemoOwnershipTransferred)
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

func (it *PingPongDemoOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *PingPongDemoOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PingPongDemoOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_PingPongDemo *PingPongDemoFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PingPongDemoOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PingPongDemo.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &PingPongDemoOwnershipTransferredIterator{contract: _PingPongDemo.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_PingPongDemo *PingPongDemoFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PingPongDemoOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PingPongDemo.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PingPongDemoOwnershipTransferred)
				if err := _PingPongDemo.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_PingPongDemo *PingPongDemoFilterer) ParseOwnershipTransferred(log types.Log) (*PingPongDemoOwnershipTransferred, error) {
	event := new(PingPongDemoOwnershipTransferred)
	if err := _PingPongDemo.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PingPongDemoPingIterator struct {
	Event *PingPongDemoPing

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PingPongDemoPingIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PingPongDemoPing)
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
		it.Event = new(PingPongDemoPing)
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

func (it *PingPongDemoPingIterator) Error() error {
	return it.fail
}

func (it *PingPongDemoPingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PingPongDemoPing struct {
	PingPongCount *big.Int
	Raw           types.Log
}

func (_PingPongDemo *PingPongDemoFilterer) FilterPing(opts *bind.FilterOpts) (*PingPongDemoPingIterator, error) {

	logs, sub, err := _PingPongDemo.contract.FilterLogs(opts, "Ping")
	if err != nil {
		return nil, err
	}
	return &PingPongDemoPingIterator{contract: _PingPongDemo.contract, event: "Ping", logs: logs, sub: sub}, nil
}

func (_PingPongDemo *PingPongDemoFilterer) WatchPing(opts *bind.WatchOpts, sink chan<- *PingPongDemoPing) (event.Subscription, error) {

	logs, sub, err := _PingPongDemo.contract.WatchLogs(opts, "Ping")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PingPongDemoPing)
				if err := _PingPongDemo.contract.UnpackLog(event, "Ping", log); err != nil {
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

func (_PingPongDemo *PingPongDemoFilterer) ParsePing(log types.Log) (*PingPongDemoPing, error) {
	event := new(PingPongDemoPing)
	if err := _PingPongDemo.contract.UnpackLog(event, "Ping", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PingPongDemoPongIterator struct {
	Event *PingPongDemoPong

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PingPongDemoPongIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PingPongDemoPong)
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
		it.Event = new(PingPongDemoPong)
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

func (it *PingPongDemoPongIterator) Error() error {
	return it.fail
}

func (it *PingPongDemoPongIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PingPongDemoPong struct {
	PingPongCount *big.Int
	Raw           types.Log
}

func (_PingPongDemo *PingPongDemoFilterer) FilterPong(opts *bind.FilterOpts) (*PingPongDemoPongIterator, error) {

	logs, sub, err := _PingPongDemo.contract.FilterLogs(opts, "Pong")
	if err != nil {
		return nil, err
	}
	return &PingPongDemoPongIterator{contract: _PingPongDemo.contract, event: "Pong", logs: logs, sub: sub}, nil
}

func (_PingPongDemo *PingPongDemoFilterer) WatchPong(opts *bind.WatchOpts, sink chan<- *PingPongDemoPong) (event.Subscription, error) {

	logs, sub, err := _PingPongDemo.contract.WatchLogs(opts, "Pong")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PingPongDemoPong)
				if err := _PingPongDemo.contract.UnpackLog(event, "Pong", log); err != nil {
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

func (_PingPongDemo *PingPongDemoFilterer) ParsePong(log types.Log) (*PingPongDemoPong, error) {
	event := new(PingPongDemoPong)
	if err := _PingPongDemo.contract.UnpackLog(event, "Pong", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type PingPongDemoTokensWithdrawnByOwnerIterator struct {
	Event *PingPongDemoTokensWithdrawnByOwner

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *PingPongDemoTokensWithdrawnByOwnerIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PingPongDemoTokensWithdrawnByOwner)
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
		it.Event = new(PingPongDemoTokensWithdrawnByOwner)
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

func (it *PingPongDemoTokensWithdrawnByOwnerIterator) Error() error {
	return it.fail
}

func (it *PingPongDemoTokensWithdrawnByOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type PingPongDemoTokensWithdrawnByOwner struct {
	Token  common.Address
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_PingPongDemo *PingPongDemoFilterer) FilterTokensWithdrawnByOwner(opts *bind.FilterOpts, token []common.Address, to []common.Address) (*PingPongDemoTokensWithdrawnByOwnerIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PingPongDemo.contract.FilterLogs(opts, "TokensWithdrawnByOwner", tokenRule, toRule)
	if err != nil {
		return nil, err
	}
	return &PingPongDemoTokensWithdrawnByOwnerIterator{contract: _PingPongDemo.contract, event: "TokensWithdrawnByOwner", logs: logs, sub: sub}, nil
}

func (_PingPongDemo *PingPongDemoFilterer) WatchTokensWithdrawnByOwner(opts *bind.WatchOpts, sink chan<- *PingPongDemoTokensWithdrawnByOwner, token []common.Address, to []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PingPongDemo.contract.WatchLogs(opts, "TokensWithdrawnByOwner", tokenRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(PingPongDemoTokensWithdrawnByOwner)
				if err := _PingPongDemo.contract.UnpackLog(event, "TokensWithdrawnByOwner", log); err != nil {
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

func (_PingPongDemo *PingPongDemoFilterer) ParseTokensWithdrawnByOwner(log types.Log) (*PingPongDemoTokensWithdrawnByOwner, error) {
	event := new(PingPongDemoTokensWithdrawnByOwner)
	if err := _PingPongDemo.contract.UnpackLog(event, "TokensWithdrawnByOwner", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetRemoteChainConfig struct {
	Recipient      []byte
	ExtraArgsBytes []byte
}
type SChainConfigs struct {
	Recipient      []byte
	ExtraArgsBytes []byte
}

func (_PingPongDemo *PingPongDemo) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _PingPongDemo.abi.Events["ApprovedSenderAdded"].ID:
		return _PingPongDemo.ParseApprovedSenderAdded(log)
	case _PingPongDemo.abi.Events["ApprovedSenderRemoved"].ID:
		return _PingPongDemo.ParseApprovedSenderRemoved(log)
	case _PingPongDemo.abi.Events["CCIPRouterModified"].ID:
		return _PingPongDemo.ParseCCIPRouterModified(log)
	case _PingPongDemo.abi.Events["ChainAdded"].ID:
		return _PingPongDemo.ParseChainAdded(log)
	case _PingPongDemo.abi.Events["ChainRemoved"].ID:
		return _PingPongDemo.ParseChainRemoved(log)
	case _PingPongDemo.abi.Events["FeeTokenUpdated"].ID:
		return _PingPongDemo.ParseFeeTokenUpdated(log)
	case _PingPongDemo.abi.Events["MessageAbandoned"].ID:
		return _PingPongDemo.ParseMessageAbandoned(log)
	case _PingPongDemo.abi.Events["MessageFailed"].ID:
		return _PingPongDemo.ParseMessageFailed(log)
	case _PingPongDemo.abi.Events["MessageRecovered"].ID:
		return _PingPongDemo.ParseMessageRecovered(log)
	case _PingPongDemo.abi.Events["MessageSent"].ID:
		return _PingPongDemo.ParseMessageSent(log)
	case _PingPongDemo.abi.Events["MessageSucceeded"].ID:
		return _PingPongDemo.ParseMessageSucceeded(log)
	case _PingPongDemo.abi.Events["OutOfOrderExecutionChange"].ID:
		return _PingPongDemo.ParseOutOfOrderExecutionChange(log)
	case _PingPongDemo.abi.Events["OwnershipTransferRequested"].ID:
		return _PingPongDemo.ParseOwnershipTransferRequested(log)
	case _PingPongDemo.abi.Events["OwnershipTransferred"].ID:
		return _PingPongDemo.ParseOwnershipTransferred(log)
	case _PingPongDemo.abi.Events["Ping"].ID:
		return _PingPongDemo.ParsePing(log)
	case _PingPongDemo.abi.Events["Pong"].ID:
		return _PingPongDemo.ParsePong(log)
	case _PingPongDemo.abi.Events["TokensWithdrawnByOwner"].ID:
		return _PingPongDemo.ParseTokensWithdrawnByOwner(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (PingPongDemoApprovedSenderAdded) Topic() common.Hash {
	return common.HexToHash("0x72d9f73bb7cb11065e15df29d61e803a0eba356d509a7025a6f51ebdea07f9e7")
}

func (PingPongDemoApprovedSenderRemoved) Topic() common.Hash {
	return common.HexToHash("0x021290bab0d93f4d9a243bd430e45dd4bc8238451e9abbba6fab4463677dfce9")
}

func (PingPongDemoCCIPRouterModified) Topic() common.Hash {
	return common.HexToHash("0x3672b589036f39ac008505b790fcb05d484d70b65680ec64c089a3c173fdc4c8")
}

func (PingPongDemoChainAdded) Topic() common.Hash {
	return common.HexToHash("0x1ced5bcae649ed29cebfa0010298ad6794bf3822e8cb754a6eee5353a9a87212")
}

func (PingPongDemoChainRemoved) Topic() common.Hash {
	return common.HexToHash("0x5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d859916")
}

func (PingPongDemoFeeTokenUpdated) Topic() common.Hash {
	return common.HexToHash("0x91a03e1d689caf891fe531c01e290f7b718f9c6a3af6726d6d837d2b7bd82e67")
}

func (PingPongDemoMessageAbandoned) Topic() common.Hash {
	return common.HexToHash("0xd5038100bd3dc9631d3c3f4f61a3e53e9d466f40c47af9897292c7b35e32a957")
}

func (PingPongDemoMessageFailed) Topic() common.Hash {
	return common.HexToHash("0x55bc02a9ef6f146737edeeb425738006f67f077e7138de3bf84a15bde1a5b56f")
}

func (PingPongDemoMessageRecovered) Topic() common.Hash {
	return common.HexToHash("0xef3bf8c64bc480286c4f3503b870ceb23e648d2d902e31fb7bb46680da6de8ad")
}

func (PingPongDemoMessageSent) Topic() common.Hash {
	return common.HexToHash("0x54791b38f3859327992a1ca0590ad3c0f08feba98d1a4f56ab0dca74d203392a")
}

func (PingPongDemoMessageSucceeded) Topic() common.Hash {
	return common.HexToHash("0xdf6958669026659bac75ba986685e11a7d271284989f565f2802522663e9a70f")
}

func (PingPongDemoOutOfOrderExecutionChange) Topic() common.Hash {
	return common.HexToHash("0x05a3fef9935c9013a24c6193df2240d34fcf6b0ebf8786b85efe8401d696cdd9")
}

func (PingPongDemoOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (PingPongDemoOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (PingPongDemoPing) Topic() common.Hash {
	return common.HexToHash("0x48257dc961b6f792c2b78a080dacfed693b660960a702de21cee364e20270e2f")
}

func (PingPongDemoPong) Topic() common.Hash {
	return common.HexToHash("0x58b69f57828e6962d216502094c54f6562f3bf082ba758966c3454f9e37b1525")
}

func (PingPongDemoTokensWithdrawnByOwner) Topic() common.Hash {
	return common.HexToHash("0x6832d9be2410a86571981e1e60fd4c1f9ea2a1034d6102a2b7d6c5e480adf02e")
}

func (_PingPongDemo *PingPongDemo) Address() common.Address {
	return _PingPongDemo.address
}

type PingPongDemoInterface interface {
	GetCounterpartAddress(opts *bind.CallOpts) (common.Address, error)

	GetCounterpartChainSelector(opts *bind.CallOpts) (uint64, error)

	GetFeeToken(opts *bind.CallOpts) (common.Address, error)

	GetMessageContents(opts *bind.CallOpts, messageId [32]byte) (ClientAny2EVMMessage, error)

	GetOutOfOrderExecution(opts *bind.CallOpts) (bool, error)

	GetRemoteChainConfig(opts *bind.CallOpts, remoteChainSelector uint64) (GetRemoteChainConfig,

		error)

	GetRouter(opts *bind.CallOpts) (common.Address, error)

	IsApprovedSender(opts *bind.CallOpts, sourceChainSelector uint64, senderAddr []byte) (bool, error)

	IsFailedMessage(opts *bind.CallOpts, messageId [32]byte) (bool, error)

	IsPaused(opts *bind.CallOpts) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SChainConfigs(opts *bind.CallOpts, destChainSelector uint64) (SChainConfigs,

		error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	UsePreFundedFeeTokens(opts *bind.CallOpts) (bool, error)

	AbandonFailedMessage(opts *bind.TransactOpts, messageId [32]byte, receiver common.Address) (*types.Transaction, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ApplyChainUpdates(opts *bind.TransactOpts, chains []CCIPBaseChainUpdate) (*types.Transaction, error)

	CcipReceive(opts *bind.TransactOpts, message ClientAny2EVMMessage) (*types.Transaction, error)

	CcipSend(opts *bind.TransactOpts, destChainSelector uint64, tokenAmounts []ClientEVMTokenAmount, data []byte) (*types.Transaction, error)

	ProcessMessage(opts *bind.TransactOpts, message ClientAny2EVMMessage) (*types.Transaction, error)

	RetryFailedMessage(opts *bind.TransactOpts, messageId [32]byte) (*types.Transaction, error)

	SetCounterpart(opts *bind.TransactOpts, counterpartChainSelector uint64, counterpartAddress common.Address) (*types.Transaction, error)

	SetPaused(opts *bind.TransactOpts, pause bool) (*types.Transaction, error)

	StartPingPong(opts *bind.TransactOpts) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateApprovedSenders(opts *bind.TransactOpts, adds []CCIPBaseApprovedSenderUpdate, removes []CCIPBaseApprovedSenderUpdate) (*types.Transaction, error)

	UpdateFeeToken(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error)

	UpdateRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error)

	WithdrawTokens(opts *bind.TransactOpts, token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error)

	Receive(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterApprovedSenderAdded(opts *bind.FilterOpts, destChainSelector []uint64, recipient [][]byte) (*PingPongDemoApprovedSenderAddedIterator, error)

	WatchApprovedSenderAdded(opts *bind.WatchOpts, sink chan<- *PingPongDemoApprovedSenderAdded, destChainSelector []uint64, recipient [][]byte) (event.Subscription, error)

	ParseApprovedSenderAdded(log types.Log) (*PingPongDemoApprovedSenderAdded, error)

	FilterApprovedSenderRemoved(opts *bind.FilterOpts, destChainSelector []uint64, recipient [][]byte) (*PingPongDemoApprovedSenderRemovedIterator, error)

	WatchApprovedSenderRemoved(opts *bind.WatchOpts, sink chan<- *PingPongDemoApprovedSenderRemoved, destChainSelector []uint64, recipient [][]byte) (event.Subscription, error)

	ParseApprovedSenderRemoved(log types.Log) (*PingPongDemoApprovedSenderRemoved, error)

	FilterCCIPRouterModified(opts *bind.FilterOpts, oldRouter []common.Address, newRouter []common.Address) (*PingPongDemoCCIPRouterModifiedIterator, error)

	WatchCCIPRouterModified(opts *bind.WatchOpts, sink chan<- *PingPongDemoCCIPRouterModified, oldRouter []common.Address, newRouter []common.Address) (event.Subscription, error)

	ParseCCIPRouterModified(log types.Log) (*PingPongDemoCCIPRouterModified, error)

	FilterChainAdded(opts *bind.FilterOpts, remoteChainSelector []uint64, recipient [][]byte) (*PingPongDemoChainAddedIterator, error)

	WatchChainAdded(opts *bind.WatchOpts, sink chan<- *PingPongDemoChainAdded, remoteChainSelector []uint64, recipient [][]byte) (event.Subscription, error)

	ParseChainAdded(log types.Log) (*PingPongDemoChainAdded, error)

	FilterChainRemoved(opts *bind.FilterOpts, removeChainSelector []uint64) (*PingPongDemoChainRemovedIterator, error)

	WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *PingPongDemoChainRemoved, removeChainSelector []uint64) (event.Subscription, error)

	ParseChainRemoved(log types.Log) (*PingPongDemoChainRemoved, error)

	FilterFeeTokenUpdated(opts *bind.FilterOpts) (*PingPongDemoFeeTokenUpdatedIterator, error)

	WatchFeeTokenUpdated(opts *bind.WatchOpts, sink chan<- *PingPongDemoFeeTokenUpdated) (event.Subscription, error)

	ParseFeeTokenUpdated(log types.Log) (*PingPongDemoFeeTokenUpdated, error)

	FilterMessageAbandoned(opts *bind.FilterOpts, messageId [][32]byte) (*PingPongDemoMessageAbandonedIterator, error)

	WatchMessageAbandoned(opts *bind.WatchOpts, sink chan<- *PingPongDemoMessageAbandoned, messageId [][32]byte) (event.Subscription, error)

	ParseMessageAbandoned(log types.Log) (*PingPongDemoMessageAbandoned, error)

	FilterMessageFailed(opts *bind.FilterOpts, messageId [][32]byte) (*PingPongDemoMessageFailedIterator, error)

	WatchMessageFailed(opts *bind.WatchOpts, sink chan<- *PingPongDemoMessageFailed, messageId [][32]byte) (event.Subscription, error)

	ParseMessageFailed(log types.Log) (*PingPongDemoMessageFailed, error)

	FilterMessageRecovered(opts *bind.FilterOpts, messageId [][32]byte) (*PingPongDemoMessageRecoveredIterator, error)

	WatchMessageRecovered(opts *bind.WatchOpts, sink chan<- *PingPongDemoMessageRecovered, messageId [][32]byte) (event.Subscription, error)

	ParseMessageRecovered(log types.Log) (*PingPongDemoMessageRecovered, error)

	FilterMessageSent(opts *bind.FilterOpts) (*PingPongDemoMessageSentIterator, error)

	WatchMessageSent(opts *bind.WatchOpts, sink chan<- *PingPongDemoMessageSent) (event.Subscription, error)

	ParseMessageSent(log types.Log) (*PingPongDemoMessageSent, error)

	FilterMessageSucceeded(opts *bind.FilterOpts, messageId [][32]byte) (*PingPongDemoMessageSucceededIterator, error)

	WatchMessageSucceeded(opts *bind.WatchOpts, sink chan<- *PingPongDemoMessageSucceeded, messageId [][32]byte) (event.Subscription, error)

	ParseMessageSucceeded(log types.Log) (*PingPongDemoMessageSucceeded, error)

	FilterOutOfOrderExecutionChange(opts *bind.FilterOpts) (*PingPongDemoOutOfOrderExecutionChangeIterator, error)

	WatchOutOfOrderExecutionChange(opts *bind.WatchOpts, sink chan<- *PingPongDemoOutOfOrderExecutionChange) (event.Subscription, error)

	ParseOutOfOrderExecutionChange(log types.Log) (*PingPongDemoOutOfOrderExecutionChange, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PingPongDemoOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *PingPongDemoOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*PingPongDemoOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PingPongDemoOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PingPongDemoOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*PingPongDemoOwnershipTransferred, error)

	FilterPing(opts *bind.FilterOpts) (*PingPongDemoPingIterator, error)

	WatchPing(opts *bind.WatchOpts, sink chan<- *PingPongDemoPing) (event.Subscription, error)

	ParsePing(log types.Log) (*PingPongDemoPing, error)

	FilterPong(opts *bind.FilterOpts) (*PingPongDemoPongIterator, error)

	WatchPong(opts *bind.WatchOpts, sink chan<- *PingPongDemoPong) (event.Subscription, error)

	ParsePong(log types.Log) (*PingPongDemoPong, error)

	FilterTokensWithdrawnByOwner(opts *bind.FilterOpts, token []common.Address, to []common.Address) (*PingPongDemoTokensWithdrawnByOwnerIterator, error)

	WatchTokensWithdrawnByOwner(opts *bind.WatchOpts, sink chan<- *PingPongDemoTokensWithdrawnByOwner, token []common.Address, to []common.Address) (event.Subscription, error)

	ParseTokensWithdrawnByOwner(log types.Log) (*PingPongDemoTokensWithdrawnByOwner, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
