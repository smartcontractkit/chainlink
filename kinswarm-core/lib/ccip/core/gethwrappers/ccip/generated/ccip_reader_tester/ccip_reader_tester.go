// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ccip_reader_tester

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

type IRMNRemoteSignature struct {
	R [32]byte
	S [32]byte
}

type InternalEVM2AnyRampMessage struct {
	Header         InternalRampMessageHeader
	Sender         common.Address
	Data           []byte
	Receiver       []byte
	ExtraArgs      []byte
	FeeToken       common.Address
	FeeTokenAmount *big.Int
	FeeValueJuels  *big.Int
	TokenAmounts   []InternalEVM2AnyTokenTransfer
}

type InternalEVM2AnyTokenTransfer struct {
	SourcePoolAddress common.Address
	DestTokenAddress  []byte
	ExtraData         []byte
	Amount            *big.Int
	DestExecData      []byte
}

type InternalGasPriceUpdate struct {
	DestChainSelector uint64
	UsdPerUnitGas     *big.Int
}

type InternalMerkleRoot struct {
	SourceChainSelector uint64
	OnRampAddress       []byte
	MinSeqNr            uint64
	MaxSeqNr            uint64
	MerkleRoot          [32]byte
}

type InternalPriceUpdates struct {
	TokenPriceUpdates []InternalTokenPriceUpdate
	GasPriceUpdates   []InternalGasPriceUpdate
}

type InternalRampMessageHeader struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	DestChainSelector   uint64
	SequenceNumber      uint64
	Nonce               uint64
}

type InternalTokenPriceUpdate struct {
	SourceToken common.Address
	UsdPerToken *big.Int
}

type OffRampCommitReport struct {
	PriceUpdates  InternalPriceUpdates
	MerkleRoots   []InternalMerkleRoot
	RmnSignatures []IRMNRemoteSignature
	RmnRawVs      *big.Int
}

type OffRampSourceChainConfig struct {
	Router    common.Address
	IsEnabled bool
	MinSeqNr  uint64
	OnRamp    []byte
}

var CCIPReaderTesterMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeValueJuels\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"sourcePoolAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destExecData\",\"type\":\"bytes\"}],\"internalType\":\"structInternal.EVM2AnyTokenTransfer[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structInternal.EVM2AnyRampMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"CCIPSendRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRampAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"maxSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structInternal.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIRMNRemote.Signature[]\",\"name\":\"rmnSignatures\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"rmnRawVs\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structOffRamp.CommitReport\",\"name\":\"report\",\"type\":\"tuple\"}],\"name\":\"CommitReportAccepted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"state\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"name\":\"ExecutionStateChanged\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeValueJuels\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"sourcePoolAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destExecData\",\"type\":\"bytes\"}],\"internalType\":\"structInternal.EVM2AnyTokenTransfer[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.EVM2AnyRampMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"emitCCIPSendRequested\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRampAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"maxSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structInternal.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIRMNRemote.Signature[]\",\"name\":\"rmnSignatures\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"rmnRawVs\",\"type\":\"uint256\"}],\"internalType\":\"structOffRamp.CommitReport\",\"name\":\"report\",\"type\":\"tuple\"}],\"name\":\"emitCommitReportAccepted\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"state\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"name\":\"emitExecutionStateChanged\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"getExpectedNextSequenceNumber\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"getSourceChainConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"setDestChainSeqNr\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfig\",\"name\":\"sourceChainConfig\",\"type\":\"tuple\"}],\"name\":\"setSourceChainConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50611730806100206000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c8063c081a9621161005b578063c081a962146100db578063c1a5a355146100ee578063e83eabba14610143578063e9d68a8e1461015657600080fd5b8063198b821c146100825780634cf66e36146100975780639041be3d146100aa575b600080fd5b610095610090366004610919565b610176565b005b6100956100a5366004610ab1565b6101b0565b6100bd6100b8366004610b30565b610207565b60405167ffffffffffffffff90911681526020015b60405180910390f35b6100956100e9366004610d00565b610237565b6100956100fc366004610e3b565b67ffffffffffffffff918216600090815260016020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001691909216179055565b610095610151366004610e6e565b61027d565b610169610164366004610b30565b610367565b6040516100d29190610f8e565b7f062172a2ae1a84d7c8d18965d6267b71d579aeca033bbd2e0548947798c3d462816040516101a59190611159565b60405180910390a150565b828467ffffffffffffffff168667ffffffffffffffff167f8c324ce1367b83031769f6a813e3bb4c117aba2185789d66b98b791405be6df285856040516101f8929190611298565b60405180910390a45050505050565b67ffffffffffffffff808216600090815260016020819052604082205491926102319216906112e8565b92915050565b8167ffffffffffffffff167f273f5ec254c4821d9795f26fddd245b3658bf1126d45d9d1c13a1d18a4c7775082604051610271919061140d565b60405180910390a25050565b67ffffffffffffffff808316600090815260208181526040918290208451815492860151938601519094167501000000000000000000000000000000000000000000027fffffff0000000000000000ffffffffffffffffffffffffffffffffffffffffff93151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00000000000000000000000000000000000000000090931673ffffffffffffffffffffffffffffffffffffffff909516949094179190911791909116919091178155606082015182919060018201906103609082611609565b5050505050565b604080516080808201835260008083526020808401829052838501829052606080850181905267ffffffffffffffff87811684528383529286902086519485018752805473ffffffffffffffffffffffffffffffffffffffff8116865274010000000000000000000000000000000000000000810460ff16151593860193909352750100000000000000000000000000000000000000000090920490921694830194909452600184018054939492939184019161042390611565565b80601f016020809104026020016040519081016040528092919081815260200182805461044f90611565565b801561049c5780601f106104715761010080835404028352916020019161049c565b820191906000526020600020905b81548152906001019060200180831161047f57829003601f168201915b5050505050815250509050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff811182821017156104fe576104fe6104ac565b60405290565b60405160a0810167ffffffffffffffff811182821017156104fe576104fe6104ac565b6040516080810167ffffffffffffffff811182821017156104fe576104fe6104ac565b604051610120810167ffffffffffffffff811182821017156104fe576104fe6104ac565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156105b5576105b56104ac565b604052919050565b600067ffffffffffffffff8211156105d7576105d76104ac565b5060051b60200190565b73ffffffffffffffffffffffffffffffffffffffff8116811461060357600080fd5b50565b8035610611816105e1565b919050565b80357bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116811461061157600080fd5b803567ffffffffffffffff8116811461061157600080fd5b600082601f83011261066b57600080fd5b8135602061068061067b836105bd565b61056e565b82815260069290921b8401810191818101908684111561069f57600080fd5b8286015b848110156106ec57604081890312156106bc5760008081fd5b6106c46104db565b6106cd82610642565b81526106da858301610616565b818601528352918301916040016106a3565b509695505050505050565b600082601f83011261070857600080fd5b813567ffffffffffffffff811115610722576107226104ac565b61075360207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8401160161056e565b81815284602083860101111561076857600080fd5b816020850160208301376000918101602001919091529392505050565b600082601f83011261079657600080fd5b813560206107a661067b836105bd565b82815260059290921b840181019181810190868411156107c557600080fd5b8286015b848110156106ec57803567ffffffffffffffff808211156107ea5760008081fd5b818901915060a0807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848d030112156108235760008081fd5b61082b610504565b610836888501610642565b81526040808501358481111561084c5760008081fd5b61085a8e8b838901016106f7565b8a840152506060935061086e848601610642565b90820152608061087f858201610642565b938201939093529201359082015283529183019183016107c9565b600082601f8301126108ab57600080fd5b813560206108bb61067b836105bd565b82815260069290921b840181019181810190868411156108da57600080fd5b8286015b848110156106ec57604081890312156108f75760008081fd5b6108ff6104db565b8135815284820135858201528352918301916040016108de565b6000602080838503121561092c57600080fd5b823567ffffffffffffffff8082111561094457600080fd5b908401906080828703121561095857600080fd5b610960610527565b82358281111561096f57600080fd5b8301604081890381131561098257600080fd5b61098a6104db565b82358581111561099957600080fd5b8301601f81018b136109aa57600080fd5b80356109b861067b826105bd565b81815260069190911b8201890190898101908d8311156109d757600080fd5b928a01925b82841015610a275785848f0312156109f45760008081fd5b6109fc6104db565b8435610a07816105e1565b8152610a14858d01610616565b818d0152825292850192908a01906109dc565b845250505082870135915084821115610a3f57600080fd5b610a4b8a83850161065a565b81880152835250508284013582811115610a6457600080fd5b610a7088828601610785565b85830152506040830135935081841115610a8957600080fd5b610a958785850161089a565b6040820152606083013560608201528094505050505092915050565b600080600080600060a08688031215610ac957600080fd5b610ad286610642565b9450610ae060208701610642565b935060408601359250606086013560048110610afb57600080fd5b9150608086013567ffffffffffffffff811115610b1757600080fd5b610b23888289016106f7565b9150509295509295909350565b600060208284031215610b4257600080fd5b610b4b82610642565b9392505050565b600060a08284031215610b6457600080fd5b610b6c610504565b905081358152610b7e60208301610642565b6020820152610b8f60408301610642565b6040820152610ba060608301610642565b6060820152610bb160808301610642565b608082015292915050565b600082601f830112610bcd57600080fd5b81356020610bdd61067b836105bd565b82815260059290921b84018101918181019086841115610bfc57600080fd5b8286015b848110156106ec57803567ffffffffffffffff80821115610c215760008081fd5b818901915060a0807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848d03011215610c5a5760008081fd5b610c62610504565b610c6d888501610606565b815260408085013584811115610c835760008081fd5b610c918e8b838901016106f7565b8a8401525060608086013585811115610caa5760008081fd5b610cb88f8c838a01016106f7565b83850152506080915081860135818401525082850135925083831115610cde5760008081fd5b610cec8d8a858801016106f7565b908201528652505050918301918301610c00565b60008060408385031215610d1357600080fd5b610d1c83610642565b9150602083013567ffffffffffffffff80821115610d3957600080fd5b908401906101a08287031215610d4e57600080fd5b610d5661054a565b610d608784610b52565b8152610d6e60a08401610606565b602082015260c083013582811115610d8557600080fd5b610d91888286016106f7565b60408301525060e083013582811115610da957600080fd5b610db5888286016106f7565b6060830152506101008084013583811115610dcf57600080fd5b610ddb898287016106f7565b608084015250610dee6101208501610606565b60a083015261014084013560c083015261016084013560e083015261018084013583811115610e1c57600080fd5b610e2889828701610bbc565b8284015250508093505050509250929050565b60008060408385031215610e4e57600080fd5b610e5783610642565b9150610e6560208401610642565b90509250929050565b60008060408385031215610e8157600080fd5b610e8a83610642565b9150602083013567ffffffffffffffff80821115610ea757600080fd5b9084019060808287031215610ebb57600080fd5b610ec3610527565b8235610ece816105e1565b815260208301358015158114610ee357600080fd5b6020820152610ef460408401610642565b6040820152606083013582811115610f0b57600080fd5b610f17888286016106f7565b6060830152508093505050509250929050565b6000815180845260005b81811015610f5057602081850181015186830182015201610f34565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815273ffffffffffffffffffffffffffffffffffffffff825116602082015260208201511515604082015267ffffffffffffffff604083015116606082015260006060830151608080840152610fe960a0840182610f2a565b949350505050565b60008151808452602080850194506020840160005b83811015611055578151805167ffffffffffffffff1688528301517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff168388015260409096019590820190600101611006565b509495945050505050565b600082825180855260208086019550808260051b84010181860160005b84811015611110577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0868403018952815160a067ffffffffffffffff8083511686528683015182888801526110d483880182610f2a565b6040858101518416908901526060808601519093169288019290925250608092830151929095019190915250978301979083019060010161107d565b5090979650505050505050565b60008151808452602080850194506020840160005b83811015611055578151805188528301518388015260409096019590820190600101611132565b602080825282516080838301528051604060a08501819052815160e08601819052600094939284019185916101008801905b808410156111e8578451805173ffffffffffffffffffffffffffffffffffffffff1683528701517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff168783015293860193600193909301929082019061118b565b50938501518785037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600160c0890152936112228186610ff1565b9450505050508185015191507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0808583030160408601526112638284611060565b9250604086015191508085840301606086015250611281828261111d565b915050606084015160808401528091505092915050565b6000600484106112d1577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b83825260406020830152610fe96040830184610f2a565b67ffffffffffffffff818116838216019080821115611330577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5092915050565b600082825180855260208086019550808260051b84010181860160005b84811015611110577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0868403018952815160a073ffffffffffffffffffffffffffffffffffffffff82511685528582015181878701526113b682870182610f2a565b915050604080830151868303828801526113d08382610f2a565b925050506060808301518187015250608080830151925085820381870152506113f98183610f2a565b9a86019a9450505090830190600101611354565b6020815261145e60208201835180518252602081015167ffffffffffffffff808216602085015280604084015116604085015280606084015116606085015280608084015116608085015250505050565b6000602083015161148760c084018273ffffffffffffffffffffffffffffffffffffffff169052565b5060408301516101a08060e08501526114a46101c0850183610f2a565b915060608501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe06101008187860301818801526114e28584610f2a565b94506080880151925081878603016101208801526115008584610f2a565b945060a0880151925061152c61014088018473ffffffffffffffffffffffffffffffffffffffff169052565b60c088015161016088015260e088015161018088015287015186850390910183870152905061155b8382611337565b9695505050505050565b600181811c9082168061157957607f821691505b6020821081036115b2577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f821115611604576000816000526020600020601f850160051c810160208610156115e15750805b601f850160051c820191505b81811015611600578281556001016115ed565b5050505b505050565b815167ffffffffffffffff811115611623576116236104ac565b611637816116318454611565565b846115b8565b602080601f83116001811461168a57600084156116545750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555611600565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b828110156116d7578886015182559484019460019091019084016116b8565b508582101561171357878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b0190555056fea164736f6c6343000818000a",
}

var CCIPReaderTesterABI = CCIPReaderTesterMetaData.ABI

var CCIPReaderTesterBin = CCIPReaderTesterMetaData.Bin

func DeployCCIPReaderTester(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CCIPReaderTester, error) {
	parsed, err := CCIPReaderTesterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CCIPReaderTesterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CCIPReaderTester{address: address, abi: *parsed, CCIPReaderTesterCaller: CCIPReaderTesterCaller{contract: contract}, CCIPReaderTesterTransactor: CCIPReaderTesterTransactor{contract: contract}, CCIPReaderTesterFilterer: CCIPReaderTesterFilterer{contract: contract}}, nil
}

type CCIPReaderTester struct {
	address common.Address
	abi     abi.ABI
	CCIPReaderTesterCaller
	CCIPReaderTesterTransactor
	CCIPReaderTesterFilterer
}

type CCIPReaderTesterCaller struct {
	contract *bind.BoundContract
}

type CCIPReaderTesterTransactor struct {
	contract *bind.BoundContract
}

type CCIPReaderTesterFilterer struct {
	contract *bind.BoundContract
}

type CCIPReaderTesterSession struct {
	Contract     *CCIPReaderTester
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type CCIPReaderTesterCallerSession struct {
	Contract *CCIPReaderTesterCaller
	CallOpts bind.CallOpts
}

type CCIPReaderTesterTransactorSession struct {
	Contract     *CCIPReaderTesterTransactor
	TransactOpts bind.TransactOpts
}

type CCIPReaderTesterRaw struct {
	Contract *CCIPReaderTester
}

type CCIPReaderTesterCallerRaw struct {
	Contract *CCIPReaderTesterCaller
}

type CCIPReaderTesterTransactorRaw struct {
	Contract *CCIPReaderTesterTransactor
}

func NewCCIPReaderTester(address common.Address, backend bind.ContractBackend) (*CCIPReaderTester, error) {
	abi, err := abi.JSON(strings.NewReader(CCIPReaderTesterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindCCIPReaderTester(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CCIPReaderTester{address: address, abi: abi, CCIPReaderTesterCaller: CCIPReaderTesterCaller{contract: contract}, CCIPReaderTesterTransactor: CCIPReaderTesterTransactor{contract: contract}, CCIPReaderTesterFilterer: CCIPReaderTesterFilterer{contract: contract}}, nil
}

func NewCCIPReaderTesterCaller(address common.Address, caller bind.ContractCaller) (*CCIPReaderTesterCaller, error) {
	contract, err := bindCCIPReaderTester(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CCIPReaderTesterCaller{contract: contract}, nil
}

func NewCCIPReaderTesterTransactor(address common.Address, transactor bind.ContractTransactor) (*CCIPReaderTesterTransactor, error) {
	contract, err := bindCCIPReaderTester(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CCIPReaderTesterTransactor{contract: contract}, nil
}

func NewCCIPReaderTesterFilterer(address common.Address, filterer bind.ContractFilterer) (*CCIPReaderTesterFilterer, error) {
	contract, err := bindCCIPReaderTester(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CCIPReaderTesterFilterer{contract: contract}, nil
}

func bindCCIPReaderTester(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CCIPReaderTesterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_CCIPReaderTester *CCIPReaderTesterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CCIPReaderTester.Contract.CCIPReaderTesterCaller.contract.Call(opts, result, method, params...)
}

func (_CCIPReaderTester *CCIPReaderTesterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.CCIPReaderTesterTransactor.contract.Transfer(opts)
}

func (_CCIPReaderTester *CCIPReaderTesterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.CCIPReaderTesterTransactor.contract.Transact(opts, method, params...)
}

func (_CCIPReaderTester *CCIPReaderTesterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CCIPReaderTester.Contract.contract.Call(opts, result, method, params...)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.contract.Transfer(opts)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.contract.Transact(opts, method, params...)
}

func (_CCIPReaderTester *CCIPReaderTesterCaller) GetExpectedNextSequenceNumber(opts *bind.CallOpts, destChainSelector uint64) (uint64, error) {
	var out []interface{}
	err := _CCIPReaderTester.contract.Call(opts, &out, "getExpectedNextSequenceNumber", destChainSelector)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_CCIPReaderTester *CCIPReaderTesterSession) GetExpectedNextSequenceNumber(destChainSelector uint64) (uint64, error) {
	return _CCIPReaderTester.Contract.GetExpectedNextSequenceNumber(&_CCIPReaderTester.CallOpts, destChainSelector)
}

func (_CCIPReaderTester *CCIPReaderTesterCallerSession) GetExpectedNextSequenceNumber(destChainSelector uint64) (uint64, error) {
	return _CCIPReaderTester.Contract.GetExpectedNextSequenceNumber(&_CCIPReaderTester.CallOpts, destChainSelector)
}

func (_CCIPReaderTester *CCIPReaderTesterCaller) GetSourceChainConfig(opts *bind.CallOpts, sourceChainSelector uint64) (OffRampSourceChainConfig, error) {
	var out []interface{}
	err := _CCIPReaderTester.contract.Call(opts, &out, "getSourceChainConfig", sourceChainSelector)

	if err != nil {
		return *new(OffRampSourceChainConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OffRampSourceChainConfig)).(*OffRampSourceChainConfig)

	return out0, err

}

func (_CCIPReaderTester *CCIPReaderTesterSession) GetSourceChainConfig(sourceChainSelector uint64) (OffRampSourceChainConfig, error) {
	return _CCIPReaderTester.Contract.GetSourceChainConfig(&_CCIPReaderTester.CallOpts, sourceChainSelector)
}

func (_CCIPReaderTester *CCIPReaderTesterCallerSession) GetSourceChainConfig(sourceChainSelector uint64) (OffRampSourceChainConfig, error) {
	return _CCIPReaderTester.Contract.GetSourceChainConfig(&_CCIPReaderTester.CallOpts, sourceChainSelector)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactor) EmitCCIPSendRequested(opts *bind.TransactOpts, destChainSelector uint64, message InternalEVM2AnyRampMessage) (*types.Transaction, error) {
	return _CCIPReaderTester.contract.Transact(opts, "emitCCIPSendRequested", destChainSelector, message)
}

func (_CCIPReaderTester *CCIPReaderTesterSession) EmitCCIPSendRequested(destChainSelector uint64, message InternalEVM2AnyRampMessage) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.EmitCCIPSendRequested(&_CCIPReaderTester.TransactOpts, destChainSelector, message)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactorSession) EmitCCIPSendRequested(destChainSelector uint64, message InternalEVM2AnyRampMessage) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.EmitCCIPSendRequested(&_CCIPReaderTester.TransactOpts, destChainSelector, message)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactor) EmitCommitReportAccepted(opts *bind.TransactOpts, report OffRampCommitReport) (*types.Transaction, error) {
	return _CCIPReaderTester.contract.Transact(opts, "emitCommitReportAccepted", report)
}

func (_CCIPReaderTester *CCIPReaderTesterSession) EmitCommitReportAccepted(report OffRampCommitReport) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.EmitCommitReportAccepted(&_CCIPReaderTester.TransactOpts, report)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactorSession) EmitCommitReportAccepted(report OffRampCommitReport) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.EmitCommitReportAccepted(&_CCIPReaderTester.TransactOpts, report)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactor) EmitExecutionStateChanged(opts *bind.TransactOpts, sourceChainSelector uint64, sequenceNumber uint64, messageId [32]byte, state uint8, returnData []byte) (*types.Transaction, error) {
	return _CCIPReaderTester.contract.Transact(opts, "emitExecutionStateChanged", sourceChainSelector, sequenceNumber, messageId, state, returnData)
}

func (_CCIPReaderTester *CCIPReaderTesterSession) EmitExecutionStateChanged(sourceChainSelector uint64, sequenceNumber uint64, messageId [32]byte, state uint8, returnData []byte) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.EmitExecutionStateChanged(&_CCIPReaderTester.TransactOpts, sourceChainSelector, sequenceNumber, messageId, state, returnData)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactorSession) EmitExecutionStateChanged(sourceChainSelector uint64, sequenceNumber uint64, messageId [32]byte, state uint8, returnData []byte) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.EmitExecutionStateChanged(&_CCIPReaderTester.TransactOpts, sourceChainSelector, sequenceNumber, messageId, state, returnData)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactor) SetDestChainSeqNr(opts *bind.TransactOpts, destChainSelector uint64, sequenceNumber uint64) (*types.Transaction, error) {
	return _CCIPReaderTester.contract.Transact(opts, "setDestChainSeqNr", destChainSelector, sequenceNumber)
}

func (_CCIPReaderTester *CCIPReaderTesterSession) SetDestChainSeqNr(destChainSelector uint64, sequenceNumber uint64) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.SetDestChainSeqNr(&_CCIPReaderTester.TransactOpts, destChainSelector, sequenceNumber)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactorSession) SetDestChainSeqNr(destChainSelector uint64, sequenceNumber uint64) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.SetDestChainSeqNr(&_CCIPReaderTester.TransactOpts, destChainSelector, sequenceNumber)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactor) SetSourceChainConfig(opts *bind.TransactOpts, sourceChainSelector uint64, sourceChainConfig OffRampSourceChainConfig) (*types.Transaction, error) {
	return _CCIPReaderTester.contract.Transact(opts, "setSourceChainConfig", sourceChainSelector, sourceChainConfig)
}

func (_CCIPReaderTester *CCIPReaderTesterSession) SetSourceChainConfig(sourceChainSelector uint64, sourceChainConfig OffRampSourceChainConfig) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.SetSourceChainConfig(&_CCIPReaderTester.TransactOpts, sourceChainSelector, sourceChainConfig)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactorSession) SetSourceChainConfig(sourceChainSelector uint64, sourceChainConfig OffRampSourceChainConfig) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.SetSourceChainConfig(&_CCIPReaderTester.TransactOpts, sourceChainSelector, sourceChainConfig)
}

type CCIPReaderTesterCCIPSendRequestedIterator struct {
	Event *CCIPReaderTesterCCIPSendRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CCIPReaderTesterCCIPSendRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPReaderTesterCCIPSendRequested)
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
		it.Event = new(CCIPReaderTesterCCIPSendRequested)
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

func (it *CCIPReaderTesterCCIPSendRequestedIterator) Error() error {
	return it.fail
}

func (it *CCIPReaderTesterCCIPSendRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CCIPReaderTesterCCIPSendRequested struct {
	DestChainSelector uint64
	Message           InternalEVM2AnyRampMessage
	Raw               types.Log
}

func (_CCIPReaderTester *CCIPReaderTesterFilterer) FilterCCIPSendRequested(opts *bind.FilterOpts, destChainSelector []uint64) (*CCIPReaderTesterCCIPSendRequestedIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _CCIPReaderTester.contract.FilterLogs(opts, "CCIPSendRequested", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &CCIPReaderTesterCCIPSendRequestedIterator{contract: _CCIPReaderTester.contract, event: "CCIPSendRequested", logs: logs, sub: sub}, nil
}

func (_CCIPReaderTester *CCIPReaderTesterFilterer) WatchCCIPSendRequested(opts *bind.WatchOpts, sink chan<- *CCIPReaderTesterCCIPSendRequested, destChainSelector []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}

	logs, sub, err := _CCIPReaderTester.contract.WatchLogs(opts, "CCIPSendRequested", destChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CCIPReaderTesterCCIPSendRequested)
				if err := _CCIPReaderTester.contract.UnpackLog(event, "CCIPSendRequested", log); err != nil {
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

func (_CCIPReaderTester *CCIPReaderTesterFilterer) ParseCCIPSendRequested(log types.Log) (*CCIPReaderTesterCCIPSendRequested, error) {
	event := new(CCIPReaderTesterCCIPSendRequested)
	if err := _CCIPReaderTester.contract.UnpackLog(event, "CCIPSendRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CCIPReaderTesterCommitReportAcceptedIterator struct {
	Event *CCIPReaderTesterCommitReportAccepted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CCIPReaderTesterCommitReportAcceptedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPReaderTesterCommitReportAccepted)
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
		it.Event = new(CCIPReaderTesterCommitReportAccepted)
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

func (it *CCIPReaderTesterCommitReportAcceptedIterator) Error() error {
	return it.fail
}

func (it *CCIPReaderTesterCommitReportAcceptedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CCIPReaderTesterCommitReportAccepted struct {
	Report OffRampCommitReport
	Raw    types.Log
}

func (_CCIPReaderTester *CCIPReaderTesterFilterer) FilterCommitReportAccepted(opts *bind.FilterOpts) (*CCIPReaderTesterCommitReportAcceptedIterator, error) {

	logs, sub, err := _CCIPReaderTester.contract.FilterLogs(opts, "CommitReportAccepted")
	if err != nil {
		return nil, err
	}
	return &CCIPReaderTesterCommitReportAcceptedIterator{contract: _CCIPReaderTester.contract, event: "CommitReportAccepted", logs: logs, sub: sub}, nil
}

func (_CCIPReaderTester *CCIPReaderTesterFilterer) WatchCommitReportAccepted(opts *bind.WatchOpts, sink chan<- *CCIPReaderTesterCommitReportAccepted) (event.Subscription, error) {

	logs, sub, err := _CCIPReaderTester.contract.WatchLogs(opts, "CommitReportAccepted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CCIPReaderTesterCommitReportAccepted)
				if err := _CCIPReaderTester.contract.UnpackLog(event, "CommitReportAccepted", log); err != nil {
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

func (_CCIPReaderTester *CCIPReaderTesterFilterer) ParseCommitReportAccepted(log types.Log) (*CCIPReaderTesterCommitReportAccepted, error) {
	event := new(CCIPReaderTesterCommitReportAccepted)
	if err := _CCIPReaderTester.contract.UnpackLog(event, "CommitReportAccepted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CCIPReaderTesterExecutionStateChangedIterator struct {
	Event *CCIPReaderTesterExecutionStateChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CCIPReaderTesterExecutionStateChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPReaderTesterExecutionStateChanged)
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
		it.Event = new(CCIPReaderTesterExecutionStateChanged)
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

func (it *CCIPReaderTesterExecutionStateChangedIterator) Error() error {
	return it.fail
}

func (it *CCIPReaderTesterExecutionStateChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CCIPReaderTesterExecutionStateChanged struct {
	SourceChainSelector uint64
	SequenceNumber      uint64
	MessageId           [32]byte
	State               uint8
	ReturnData          []byte
	Raw                 types.Log
}

func (_CCIPReaderTester *CCIPReaderTesterFilterer) FilterExecutionStateChanged(opts *bind.FilterOpts, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (*CCIPReaderTesterExecutionStateChangedIterator, error) {

	var sourceChainSelectorRule []interface{}
	for _, sourceChainSelectorItem := range sourceChainSelector {
		sourceChainSelectorRule = append(sourceChainSelectorRule, sourceChainSelectorItem)
	}
	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}
	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}

	logs, sub, err := _CCIPReaderTester.contract.FilterLogs(opts, "ExecutionStateChanged", sourceChainSelectorRule, sequenceNumberRule, messageIdRule)
	if err != nil {
		return nil, err
	}
	return &CCIPReaderTesterExecutionStateChangedIterator{contract: _CCIPReaderTester.contract, event: "ExecutionStateChanged", logs: logs, sub: sub}, nil
}

func (_CCIPReaderTester *CCIPReaderTesterFilterer) WatchExecutionStateChanged(opts *bind.WatchOpts, sink chan<- *CCIPReaderTesterExecutionStateChanged, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (event.Subscription, error) {

	var sourceChainSelectorRule []interface{}
	for _, sourceChainSelectorItem := range sourceChainSelector {
		sourceChainSelectorRule = append(sourceChainSelectorRule, sourceChainSelectorItem)
	}
	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}
	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}

	logs, sub, err := _CCIPReaderTester.contract.WatchLogs(opts, "ExecutionStateChanged", sourceChainSelectorRule, sequenceNumberRule, messageIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CCIPReaderTesterExecutionStateChanged)
				if err := _CCIPReaderTester.contract.UnpackLog(event, "ExecutionStateChanged", log); err != nil {
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

func (_CCIPReaderTester *CCIPReaderTesterFilterer) ParseExecutionStateChanged(log types.Log) (*CCIPReaderTesterExecutionStateChanged, error) {
	event := new(CCIPReaderTesterExecutionStateChanged)
	if err := _CCIPReaderTester.contract.UnpackLog(event, "ExecutionStateChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_CCIPReaderTester *CCIPReaderTester) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _CCIPReaderTester.abi.Events["CCIPSendRequested"].ID:
		return _CCIPReaderTester.ParseCCIPSendRequested(log)
	case _CCIPReaderTester.abi.Events["CommitReportAccepted"].ID:
		return _CCIPReaderTester.ParseCommitReportAccepted(log)
	case _CCIPReaderTester.abi.Events["ExecutionStateChanged"].ID:
		return _CCIPReaderTester.ParseExecutionStateChanged(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (CCIPReaderTesterCCIPSendRequested) Topic() common.Hash {
	return common.HexToHash("0x273f5ec254c4821d9795f26fddd245b3658bf1126d45d9d1c13a1d18a4c77750")
}

func (CCIPReaderTesterCommitReportAccepted) Topic() common.Hash {
	return common.HexToHash("0x062172a2ae1a84d7c8d18965d6267b71d579aeca033bbd2e0548947798c3d462")
}

func (CCIPReaderTesterExecutionStateChanged) Topic() common.Hash {
	return common.HexToHash("0x8c324ce1367b83031769f6a813e3bb4c117aba2185789d66b98b791405be6df2")
}

func (_CCIPReaderTester *CCIPReaderTester) Address() common.Address {
	return _CCIPReaderTester.address
}

type CCIPReaderTesterInterface interface {
	GetExpectedNextSequenceNumber(opts *bind.CallOpts, destChainSelector uint64) (uint64, error)

	GetSourceChainConfig(opts *bind.CallOpts, sourceChainSelector uint64) (OffRampSourceChainConfig, error)

	EmitCCIPSendRequested(opts *bind.TransactOpts, destChainSelector uint64, message InternalEVM2AnyRampMessage) (*types.Transaction, error)

	EmitCommitReportAccepted(opts *bind.TransactOpts, report OffRampCommitReport) (*types.Transaction, error)

	EmitExecutionStateChanged(opts *bind.TransactOpts, sourceChainSelector uint64, sequenceNumber uint64, messageId [32]byte, state uint8, returnData []byte) (*types.Transaction, error)

	SetDestChainSeqNr(opts *bind.TransactOpts, destChainSelector uint64, sequenceNumber uint64) (*types.Transaction, error)

	SetSourceChainConfig(opts *bind.TransactOpts, sourceChainSelector uint64, sourceChainConfig OffRampSourceChainConfig) (*types.Transaction, error)

	FilterCCIPSendRequested(opts *bind.FilterOpts, destChainSelector []uint64) (*CCIPReaderTesterCCIPSendRequestedIterator, error)

	WatchCCIPSendRequested(opts *bind.WatchOpts, sink chan<- *CCIPReaderTesterCCIPSendRequested, destChainSelector []uint64) (event.Subscription, error)

	ParseCCIPSendRequested(log types.Log) (*CCIPReaderTesterCCIPSendRequested, error)

	FilterCommitReportAccepted(opts *bind.FilterOpts) (*CCIPReaderTesterCommitReportAcceptedIterator, error)

	WatchCommitReportAccepted(opts *bind.WatchOpts, sink chan<- *CCIPReaderTesterCommitReportAccepted) (event.Subscription, error)

	ParseCommitReportAccepted(log types.Log) (*CCIPReaderTesterCommitReportAccepted, error)

	FilterExecutionStateChanged(opts *bind.FilterOpts, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (*CCIPReaderTesterExecutionStateChangedIterator, error)

	WatchExecutionStateChanged(opts *bind.WatchOpts, sink chan<- *CCIPReaderTesterExecutionStateChanged, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (event.Subscription, error)

	ParseExecutionStateChanged(log types.Log) (*CCIPReaderTesterExecutionStateChanged, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
