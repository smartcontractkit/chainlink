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
}

type OffRampSourceChainConfig struct {
	Router    common.Address
	IsEnabled bool
	MinSeqNr  uint64
	OnRamp    []byte
}

var CCIPReaderTesterMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeValueJuels\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"sourcePoolAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destExecData\",\"type\":\"bytes\"}],\"internalType\":\"structInternal.EVM2AnyTokenTransfer[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structInternal.EVM2AnyRampMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"CCIPMessageSent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRampAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"maxSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"indexed\":false,\"internalType\":\"structInternal.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"}],\"name\":\"CommitReportAccepted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"messageHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"state\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"name\":\"ExecutionStateChanged\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeValueJuels\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"sourcePoolAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destExecData\",\"type\":\"bytes\"}],\"internalType\":\"structInternal.EVM2AnyTokenTransfer[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.EVM2AnyRampMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"emitCCIPMessageSent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRampAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"maxSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structInternal.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIRMNRemote.Signature[]\",\"name\":\"rmnSignatures\",\"type\":\"tuple[]\"}],\"internalType\":\"structOffRamp.CommitReport\",\"name\":\"report\",\"type\":\"tuple\"}],\"name\":\"emitCommitReportAccepted\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"messageHash\",\"type\":\"bytes32\"},{\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"state\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"name\":\"emitExecutionStateChanged\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"getExpectedNextSequenceNumber\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"}],\"name\":\"getInboundNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"getSourceChainConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"setDestChainSeqNr\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"testNonce\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"}],\"name\":\"setInboundNonce\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfig\",\"name\":\"sourceChainConfig\",\"type\":\"tuple\"}],\"name\":\"setSourceChainConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506118e3806100206000396000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c8063c1a5a35511610076578063c92236251161005b578063c92236251461017c578063e83eabba1461018f578063e9d68a8e146101a257600080fd5b8063c1a5a35514610114578063c7c1cba11461016957600080fd5b80634bf78697146100a85780639041be3d146100bd57806393df2867146100ee578063bfc9b78914610101575b600080fd5b6100bb6100b63660046109d1565b6101c2565b005b6100d06100cb366004610b0c565b61021b565b60405167ffffffffffffffff90911681526020015b60405180910390f35b6100bb6100fc366004610b77565b61024b565b6100bb61010f366004610e25565b6102c6565b6100bb610122366004610fb0565b67ffffffffffffffff918216600090815260016020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001691909216179055565b6100bb610177366004610fe3565b610308565b6100d061018a366004611075565b610365565b6100bb61019d3660046110c8565b6103b1565b6101b56101b0366004610b0c565b61049b565b6040516100e591906111e8565b80600001516060015167ffffffffffffffff168267ffffffffffffffff167f192442a2b2adb6a7948f097023cb6b57d29d3a7a5dd33e6666d33c39cc456f328360405161020f919061132e565b60405180910390a35050565b67ffffffffffffffff80821660009081526001602081905260408220549192610245921690611486565b92915050565b67ffffffffffffffff841660009081526002602052604090819020905184919061027890859085906114d5565b908152604051908190036020019020805467ffffffffffffffff929092167fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000090921691909117905550505050565b602081015181516040517f35c02761bcd3ef995c6a601a1981f4ed3934dcbe5041e24e286c89f5531d17e4926102fd9290916115d9565b60405180910390a150565b848667ffffffffffffffff168867ffffffffffffffff167f05665fe9ad095383d018353f4cbcba77e84db27dd215081bbf7cdf9ae6fbe48b8787878760405161035494939291906116b1565b60405180910390a450505050505050565b67ffffffffffffffff8316600090815260026020526040808220905161038e90859085906114d5565b9081526040519081900360200190205467ffffffffffffffff1690509392505050565b67ffffffffffffffff808316600090815260208181526040918290208451815492860151938601519094167501000000000000000000000000000000000000000000027fffffff0000000000000000ffffffffffffffffffffffffffffffffffffffffff93151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00000000000000000000000000000000000000000090931673ffffffffffffffffffffffffffffffffffffffff9095169490941791909117919091169190911781556060820151829190600182019061049490826117bc565b5050505050565b604080516080808201835260008083526020808401829052838501829052606080850181905267ffffffffffffffff87811684528383529286902086519485018752805473ffffffffffffffffffffffffffffffffffffffff8116865274010000000000000000000000000000000000000000810460ff16151593860193909352750100000000000000000000000000000000000000000090920490921694830194909452600184018054939492939184019161055790611718565b80601f016020809104026020016040519081016040528092919081815260200182805461058390611718565b80156105d05780601f106105a5576101008083540402835291602001916105d0565b820191906000526020600020905b8154815290600101906020018083116105b357829003601f168201915b5050505050815250509050919050565b803567ffffffffffffffff811681146105f857600080fd5b919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160a0810167ffffffffffffffff8111828210171561064f5761064f6105fd565b60405290565b604051610120810167ffffffffffffffff8111828210171561064f5761064f6105fd565b6040805190810167ffffffffffffffff8111828210171561064f5761064f6105fd565b6040516060810167ffffffffffffffff8111828210171561064f5761064f6105fd565b6040516080810167ffffffffffffffff8111828210171561064f5761064f6105fd565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610729576107296105fd565b604052919050565b600060a0828403121561074357600080fd5b61074b61062c565b90508135815261075d602083016105e0565b602082015261076e604083016105e0565b604082015261077f606083016105e0565b6060820152610790608083016105e0565b608082015292915050565b73ffffffffffffffffffffffffffffffffffffffff811681146107bd57600080fd5b50565b80356105f88161079b565b600082601f8301126107dc57600080fd5b813567ffffffffffffffff8111156107f6576107f66105fd565b61082760207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016106e2565b81815284602083860101111561083c57600080fd5b816020850160208301376000918101602001919091529392505050565b600067ffffffffffffffff821115610873576108736105fd565b5060051b60200190565b600082601f83011261088e57600080fd5b813560206108a361089e83610859565b6106e2565b82815260059290921b840181019181810190868411156108c257600080fd5b8286015b848110156109c657803567ffffffffffffffff808211156108e75760008081fd5b818901915060a0807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848d030112156109205760008081fd5b61092861062c565b6109338885016107c0565b8152604080850135848111156109495760008081fd5b6109578e8b838901016107cb565b8a84015250606080860135858111156109705760008081fd5b61097e8f8c838a01016107cb565b838501525060809150818601358184015250828501359250838311156109a45760008081fd5b6109b28d8a858801016107cb565b9082015286525050509183019183016108c6565b509695505050505050565b600080604083850312156109e457600080fd5b6109ed836105e0565b9150602083013567ffffffffffffffff80821115610a0a57600080fd5b908401906101a08287031215610a1f57600080fd5b610a27610655565b610a318784610731565b8152610a3f60a084016107c0565b602082015260c083013582811115610a5657600080fd5b610a62888286016107cb565b60408301525060e083013582811115610a7a57600080fd5b610a86888286016107cb565b6060830152506101008084013583811115610aa057600080fd5b610aac898287016107cb565b608084015250610abf61012085016107c0565b60a083015261014084013560c083015261016084013560e083015261018084013583811115610aed57600080fd5b610af98982870161087d565b8284015250508093505050509250929050565b600060208284031215610b1e57600080fd5b610b27826105e0565b9392505050565b60008083601f840112610b4057600080fd5b50813567ffffffffffffffff811115610b5857600080fd5b602083019150836020828501011115610b7057600080fd5b9250929050565b60008060008060608587031215610b8d57600080fd5b610b96856105e0565b9350610ba4602086016105e0565b9250604085013567ffffffffffffffff811115610bc057600080fd5b610bcc87828801610b2e565b95989497509550505050565b80357bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811681146105f857600080fd5b600082601f830112610c1557600080fd5b81356020610c2561089e83610859565b82815260069290921b84018101918181019086841115610c4457600080fd5b8286015b848110156109c65760408189031215610c615760008081fd5b610c69610679565b610c72826105e0565b8152610c7f858301610bd8565b81860152835291830191604001610c48565b600082601f830112610ca257600080fd5b81356020610cb261089e83610859565b82815260059290921b84018101918181019086841115610cd157600080fd5b8286015b848110156109c657803567ffffffffffffffff80821115610cf65760008081fd5b818901915060a0807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848d03011215610d2f5760008081fd5b610d3761062c565b610d428885016105e0565b815260408085013584811115610d585760008081fd5b610d668e8b838901016107cb565b8a8401525060609350610d7a8486016105e0565b908201526080610d8b8582016105e0565b93820193909352920135908201528352918301918301610cd5565b600082601f830112610db757600080fd5b81356020610dc761089e83610859565b82815260069290921b84018101918181019086841115610de657600080fd5b8286015b848110156109c65760408189031215610e035760008081fd5b610e0b610679565b813581528482013585820152835291830191604001610dea565b60006020808385031215610e3857600080fd5b823567ffffffffffffffff80821115610e5057600080fd5b9084019060608287031215610e6457600080fd5b610e6c61069c565b823582811115610e7b57600080fd5b83016040818903811315610e8e57600080fd5b610e96610679565b823585811115610ea557600080fd5b8301601f81018b13610eb657600080fd5b8035610ec461089e82610859565b81815260069190911b8201890190898101908d831115610ee357600080fd5b928a01925b82841015610f335785848f031215610f005760008081fd5b610f08610679565b8435610f138161079b565b8152610f20858d01610bd8565b818d0152825292850192908a0190610ee8565b845250505082870135915084821115610f4b57600080fd5b610f578a838501610c04565b81880152835250508284013582811115610f7057600080fd5b610f7c88828601610c91565b85830152506040830135935081841115610f9557600080fd5b610fa187858501610da6565b60408201529695505050505050565b60008060408385031215610fc357600080fd5b610fcc836105e0565b9150610fda602084016105e0565b90509250929050565b600080600080600080600060e0888a031215610ffe57600080fd5b611007886105e0565b9650611015602089016105e0565b9550604088013594506060880135935060808801356004811061103757600080fd5b925060a088013567ffffffffffffffff81111561105357600080fd5b61105f8a828b016107cb565b92505060c0880135905092959891949750929550565b60008060006040848603121561108a57600080fd5b611093846105e0565b9250602084013567ffffffffffffffff8111156110af57600080fd5b6110bb86828701610b2e565b9497909650939450505050565b600080604083850312156110db57600080fd5b6110e4836105e0565b9150602083013567ffffffffffffffff8082111561110157600080fd5b908401906080828703121561111557600080fd5b61111d6106bf565b82356111288161079b565b81526020830135801515811461113d57600080fd5b602082015261114e604084016105e0565b604082015260608301358281111561116557600080fd5b611171888286016107cb565b6060830152508093505050509250929050565b6000815180845260005b818110156111aa5760208185018101518683018201520161118e565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b6020815273ffffffffffffffffffffffffffffffffffffffff825116602082015260208201511515604082015267ffffffffffffffff60408301511660608201526000606083015160808084015261124360a0840182611184565b949350505050565b600082825180855260208086019550808260051b84010181860160005b84811015611321577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0868403018952815160a073ffffffffffffffffffffffffffffffffffffffff82511685528582015181878701526112ca82870182611184565b915050604080830151868303828801526112e48382611184565b9250505060608083015181870152506080808301519250858203818701525061130d8183611184565b9a86019a9450505090830190600101611268565b5090979650505050505050565b6020815261137f60208201835180518252602081015167ffffffffffffffff808216602085015280604084015116604085015280606084015116606085015280608084015116608085015250505050565b600060208301516113a860c084018273ffffffffffffffffffffffffffffffffffffffff169052565b5060408301516101a08060e08501526113c56101c0850183611184565b915060608501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe06101008187860301818801526114038584611184565b94506080880151925081878603016101208801526114218584611184565b945060a0880151925061144d61014088018473ffffffffffffffffffffffffffffffffffffffff169052565b60c088015161016088015260e088015161018088015287015186850390910183870152905061147c838261124b565b9695505050505050565b67ffffffffffffffff8181168382160190808211156114ce577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5092915050565b8183823760009101908152919050565b805160408084528151848201819052600092602091908201906060870190855b8181101561155e578351805173ffffffffffffffffffffffffffffffffffffffff1684528501517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16858401529284019291850191600101611505565b50508583015187820388850152805180835290840192506000918401905b808310156115cd578351805167ffffffffffffffff1683528501517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff168583015292840192600192909201919085019061157c565b50979650505050505050565b60006040808301604084528086518083526060925060608601915060608160051b8701016020808a0160005b84811015611691577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa08a8503018652815160a067ffffffffffffffff80835116875285830151828789015261165c83890182611184565b848d01518316898e01528b8501519092168b890152506080928301519290960191909152509482019490820190600101611605565b5050878203908801526116a481896114e5565b9998505050505050505050565b8481526000600485106116ed577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b846020830152608060408301526117076080830185611184565b905082606083015295945050505050565b600181811c9082168061172c57607f821691505b602082108103611765577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f8211156117b7576000816000526020600020601f850160051c810160208610156117945750805b601f850160051c820191505b818110156117b3578281556001016117a0565b5050505b505050565b815167ffffffffffffffff8111156117d6576117d66105fd565b6117ea816117e48454611718565b8461176b565b602080601f83116001811461183d57600084156118075750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556117b3565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561188a5788860151825594840194600190910190840161186b565b50858210156118c657878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b0190555056fea164736f6c6343000818000a",
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

func (_CCIPReaderTester *CCIPReaderTesterCaller) GetInboundNonce(opts *bind.CallOpts, sourceChainSelector uint64, sender []byte) (uint64, error) {
	var out []interface{}
	err := _CCIPReaderTester.contract.Call(opts, &out, "getInboundNonce", sourceChainSelector, sender)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_CCIPReaderTester *CCIPReaderTesterSession) GetInboundNonce(sourceChainSelector uint64, sender []byte) (uint64, error) {
	return _CCIPReaderTester.Contract.GetInboundNonce(&_CCIPReaderTester.CallOpts, sourceChainSelector, sender)
}

func (_CCIPReaderTester *CCIPReaderTesterCallerSession) GetInboundNonce(sourceChainSelector uint64, sender []byte) (uint64, error) {
	return _CCIPReaderTester.Contract.GetInboundNonce(&_CCIPReaderTester.CallOpts, sourceChainSelector, sender)
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

func (_CCIPReaderTester *CCIPReaderTesterTransactor) EmitCCIPMessageSent(opts *bind.TransactOpts, destChainSelector uint64, message InternalEVM2AnyRampMessage) (*types.Transaction, error) {
	return _CCIPReaderTester.contract.Transact(opts, "emitCCIPMessageSent", destChainSelector, message)
}

func (_CCIPReaderTester *CCIPReaderTesterSession) EmitCCIPMessageSent(destChainSelector uint64, message InternalEVM2AnyRampMessage) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.EmitCCIPMessageSent(&_CCIPReaderTester.TransactOpts, destChainSelector, message)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactorSession) EmitCCIPMessageSent(destChainSelector uint64, message InternalEVM2AnyRampMessage) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.EmitCCIPMessageSent(&_CCIPReaderTester.TransactOpts, destChainSelector, message)
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

func (_CCIPReaderTester *CCIPReaderTesterTransactor) EmitExecutionStateChanged(opts *bind.TransactOpts, sourceChainSelector uint64, sequenceNumber uint64, messageId [32]byte, messageHash [32]byte, state uint8, returnData []byte, gasUsed *big.Int) (*types.Transaction, error) {
	return _CCIPReaderTester.contract.Transact(opts, "emitExecutionStateChanged", sourceChainSelector, sequenceNumber, messageId, messageHash, state, returnData, gasUsed)
}

func (_CCIPReaderTester *CCIPReaderTesterSession) EmitExecutionStateChanged(sourceChainSelector uint64, sequenceNumber uint64, messageId [32]byte, messageHash [32]byte, state uint8, returnData []byte, gasUsed *big.Int) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.EmitExecutionStateChanged(&_CCIPReaderTester.TransactOpts, sourceChainSelector, sequenceNumber, messageId, messageHash, state, returnData, gasUsed)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactorSession) EmitExecutionStateChanged(sourceChainSelector uint64, sequenceNumber uint64, messageId [32]byte, messageHash [32]byte, state uint8, returnData []byte, gasUsed *big.Int) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.EmitExecutionStateChanged(&_CCIPReaderTester.TransactOpts, sourceChainSelector, sequenceNumber, messageId, messageHash, state, returnData, gasUsed)
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

func (_CCIPReaderTester *CCIPReaderTesterTransactor) SetInboundNonce(opts *bind.TransactOpts, sourceChainSelector uint64, testNonce uint64, sender []byte) (*types.Transaction, error) {
	return _CCIPReaderTester.contract.Transact(opts, "setInboundNonce", sourceChainSelector, testNonce, sender)
}

func (_CCIPReaderTester *CCIPReaderTesterSession) SetInboundNonce(sourceChainSelector uint64, testNonce uint64, sender []byte) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.SetInboundNonce(&_CCIPReaderTester.TransactOpts, sourceChainSelector, testNonce, sender)
}

func (_CCIPReaderTester *CCIPReaderTesterTransactorSession) SetInboundNonce(sourceChainSelector uint64, testNonce uint64, sender []byte) (*types.Transaction, error) {
	return _CCIPReaderTester.Contract.SetInboundNonce(&_CCIPReaderTester.TransactOpts, sourceChainSelector, testNonce, sender)
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

type CCIPReaderTesterCCIPMessageSentIterator struct {
	Event *CCIPReaderTesterCCIPMessageSent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CCIPReaderTesterCCIPMessageSentIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPReaderTesterCCIPMessageSent)
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
		it.Event = new(CCIPReaderTesterCCIPMessageSent)
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

func (it *CCIPReaderTesterCCIPMessageSentIterator) Error() error {
	return it.fail
}

func (it *CCIPReaderTesterCCIPMessageSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CCIPReaderTesterCCIPMessageSent struct {
	DestChainSelector uint64
	SequenceNumber    uint64
	Message           InternalEVM2AnyRampMessage
	Raw               types.Log
}

func (_CCIPReaderTester *CCIPReaderTesterFilterer) FilterCCIPMessageSent(opts *bind.FilterOpts, destChainSelector []uint64, sequenceNumber []uint64) (*CCIPReaderTesterCCIPMessageSentIterator, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}

	logs, sub, err := _CCIPReaderTester.contract.FilterLogs(opts, "CCIPMessageSent", destChainSelectorRule, sequenceNumberRule)
	if err != nil {
		return nil, err
	}
	return &CCIPReaderTesterCCIPMessageSentIterator{contract: _CCIPReaderTester.contract, event: "CCIPMessageSent", logs: logs, sub: sub}, nil
}

func (_CCIPReaderTester *CCIPReaderTesterFilterer) WatchCCIPMessageSent(opts *bind.WatchOpts, sink chan<- *CCIPReaderTesterCCIPMessageSent, destChainSelector []uint64, sequenceNumber []uint64) (event.Subscription, error) {

	var destChainSelectorRule []interface{}
	for _, destChainSelectorItem := range destChainSelector {
		destChainSelectorRule = append(destChainSelectorRule, destChainSelectorItem)
	}
	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}

	logs, sub, err := _CCIPReaderTester.contract.WatchLogs(opts, "CCIPMessageSent", destChainSelectorRule, sequenceNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CCIPReaderTesterCCIPMessageSent)
				if err := _CCIPReaderTester.contract.UnpackLog(event, "CCIPMessageSent", log); err != nil {
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

func (_CCIPReaderTester *CCIPReaderTesterFilterer) ParseCCIPMessageSent(log types.Log) (*CCIPReaderTesterCCIPMessageSent, error) {
	event := new(CCIPReaderTesterCCIPMessageSent)
	if err := _CCIPReaderTester.contract.UnpackLog(event, "CCIPMessageSent", log); err != nil {
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
	MerkleRoots  []InternalMerkleRoot
	PriceUpdates InternalPriceUpdates
	Raw          types.Log
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
	MessageHash         [32]byte
	State               uint8
	ReturnData          []byte
	GasUsed             *big.Int
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
	case _CCIPReaderTester.abi.Events["CCIPMessageSent"].ID:
		return _CCIPReaderTester.ParseCCIPMessageSent(log)
	case _CCIPReaderTester.abi.Events["CommitReportAccepted"].ID:
		return _CCIPReaderTester.ParseCommitReportAccepted(log)
	case _CCIPReaderTester.abi.Events["ExecutionStateChanged"].ID:
		return _CCIPReaderTester.ParseExecutionStateChanged(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (CCIPReaderTesterCCIPMessageSent) Topic() common.Hash {
	return common.HexToHash("0x192442a2b2adb6a7948f097023cb6b57d29d3a7a5dd33e6666d33c39cc456f32")
}

func (CCIPReaderTesterCommitReportAccepted) Topic() common.Hash {
	return common.HexToHash("0x35c02761bcd3ef995c6a601a1981f4ed3934dcbe5041e24e286c89f5531d17e4")
}

func (CCIPReaderTesterExecutionStateChanged) Topic() common.Hash {
	return common.HexToHash("0x05665fe9ad095383d018353f4cbcba77e84db27dd215081bbf7cdf9ae6fbe48b")
}

func (_CCIPReaderTester *CCIPReaderTester) Address() common.Address {
	return _CCIPReaderTester.address
}

type CCIPReaderTesterInterface interface {
	GetExpectedNextSequenceNumber(opts *bind.CallOpts, destChainSelector uint64) (uint64, error)

	GetInboundNonce(opts *bind.CallOpts, sourceChainSelector uint64, sender []byte) (uint64, error)

	GetSourceChainConfig(opts *bind.CallOpts, sourceChainSelector uint64) (OffRampSourceChainConfig, error)

	EmitCCIPMessageSent(opts *bind.TransactOpts, destChainSelector uint64, message InternalEVM2AnyRampMessage) (*types.Transaction, error)

	EmitCommitReportAccepted(opts *bind.TransactOpts, report OffRampCommitReport) (*types.Transaction, error)

	EmitExecutionStateChanged(opts *bind.TransactOpts, sourceChainSelector uint64, sequenceNumber uint64, messageId [32]byte, messageHash [32]byte, state uint8, returnData []byte, gasUsed *big.Int) (*types.Transaction, error)

	SetDestChainSeqNr(opts *bind.TransactOpts, destChainSelector uint64, sequenceNumber uint64) (*types.Transaction, error)

	SetInboundNonce(opts *bind.TransactOpts, sourceChainSelector uint64, testNonce uint64, sender []byte) (*types.Transaction, error)

	SetSourceChainConfig(opts *bind.TransactOpts, sourceChainSelector uint64, sourceChainConfig OffRampSourceChainConfig) (*types.Transaction, error)

	FilterCCIPMessageSent(opts *bind.FilterOpts, destChainSelector []uint64, sequenceNumber []uint64) (*CCIPReaderTesterCCIPMessageSentIterator, error)

	WatchCCIPMessageSent(opts *bind.WatchOpts, sink chan<- *CCIPReaderTesterCCIPMessageSent, destChainSelector []uint64, sequenceNumber []uint64) (event.Subscription, error)

	ParseCCIPMessageSent(log types.Log) (*CCIPReaderTesterCCIPMessageSent, error)

	FilterCommitReportAccepted(opts *bind.FilterOpts) (*CCIPReaderTesterCommitReportAcceptedIterator, error)

	WatchCommitReportAccepted(opts *bind.WatchOpts, sink chan<- *CCIPReaderTesterCommitReportAccepted) (event.Subscription, error)

	ParseCommitReportAccepted(log types.Log) (*CCIPReaderTesterCommitReportAccepted, error)

	FilterExecutionStateChanged(opts *bind.FilterOpts, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (*CCIPReaderTesterExecutionStateChangedIterator, error)

	WatchExecutionStateChanged(opts *bind.WatchOpts, sink chan<- *CCIPReaderTesterExecutionStateChanged, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (event.Subscription, error)

	ParseExecutionStateChanged(log types.Log) (*CCIPReaderTesterExecutionStateChanged, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
