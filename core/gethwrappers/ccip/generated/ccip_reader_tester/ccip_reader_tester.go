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

type InternalEVM2AnyRampMessage struct {
	Header         InternalRampMessageHeader
	Sender         common.Address
	Data           []byte
	Receiver       []byte
	ExtraArgs      []byte
	FeeToken       common.Address
	FeeTokenAmount *big.Int
	TokenAmounts   []InternalRampTokenAmount
}

type InternalGasPriceUpdate struct {
	DestChainSelector uint64
	UsdPerUnitGas     *big.Int
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

type InternalRampTokenAmount struct {
	SourcePoolAddress []byte
	DestTokenAddress  []byte
	ExtraData         []byte
	Amount            *big.Int
	DestExecData      []byte
}

type InternalTokenPriceUpdate struct {
	SourceToken common.Address
	UsdPerToken *big.Int
}

type OffRampCommitReport struct {
	PriceUpdates InternalPriceUpdates
	MerkleRoots  []OffRampMerkleRoot
}

type OffRampInterval struct {
	Min uint64
	Max uint64
}

type OffRampMerkleRoot struct {
	SourceChainSelector uint64
	Interval            OffRampInterval
	MerkleRoot          [32]byte
}

type OffRampSourceChainConfig struct {
	Router    common.Address
	IsEnabled bool
	MinSeqNr  uint64
	OnRamp    []byte
}

var CCIPReaderTesterMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destExecData\",\"type\":\"bytes\"}],\"internalType\":\"structInternal.RampTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structInternal.EVM2AnyRampMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"CCIPSendRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"min\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"max\",\"type\":\"uint64\"}],\"internalType\":\"structOffRamp.Interval\",\"name\":\"interval\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structOffRamp.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structOffRamp.CommitReport\",\"name\":\"report\",\"type\":\"tuple\"}],\"name\":\"CommitReportAccepted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"state\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"name\":\"ExecutionStateChanged\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destExecData\",\"type\":\"bytes\"}],\"internalType\":\"structInternal.RampTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.EVM2AnyRampMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"emitCCIPSendRequested\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"min\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"max\",\"type\":\"uint64\"}],\"internalType\":\"structOffRamp.Interval\",\"name\":\"interval\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structOffRamp.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"}],\"internalType\":\"structOffRamp.CommitReport\",\"name\":\"report\",\"type\":\"tuple\"}],\"name\":\"emitCommitReportAccepted\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"state\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"name\":\"emitExecutionStateChanged\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"}],\"name\":\"getExpectedNextSequenceNumber\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"}],\"name\":\"getInboundNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"getSourceChainConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"setDestChainSeqNr\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"testNonce\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"}],\"name\":\"setInboundNonce\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfig\",\"name\":\"sourceChainConfig\",\"type\":\"tuple\"}],\"name\":\"setSourceChainConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061148e806100206000396000f3fe608060405234801561001057600080fd5b50600436106100835760003560e01c80634cf66e361461008857806385096da91461009d5780639041be3d146100b057806393df2867146100e0578063c1a5a355146100f3578063c92236251461012f578063e44302b714610142578063e83eabba14610155578063e9d68a8e14610168575b600080fd5b61009b610096366004610658565b610188565b005b61009b6100ab3660046108d0565b6101dd565b6100c36100be3660046109fe565b610222565b6040516001600160401b0390911681526020015b60405180910390f35b61009b6100ee366004610a68565b610251565b61009b610101366004610ac8565b6001600160401b03918216600090815260016020526040902080546001600160401b03191691909216179055565b6100c361013d366004610afb565b6102b2565b61009b610150366004610cbd565b6102fc565b61009b610163366004610e27565b610336565b61017b6101763660046109fe565b6103c2565b6040516100d79190610f28565b82846001600160401b0316866001600160401b03167f8c324ce1367b83031769f6a813e3bb4c117aba2185789d66b98b791405be6df285856040516101ce929190610f80565b60405180910390a45050505050565b816001600160401b03167fcae3d9f68a9ed33d0a770e9bcc6fc25d4041dd6391e6cd8015546547a3c93140826040516102169190611083565b60405180910390a25050565b6001600160401b038082166000908152600160208190526040822054919261024b921690611185565b92915050565b6001600160401b03841660009081526002602052604090819020905184919061027d90859085906111ba565b90815260405190819003602001902080546001600160401b03929092166001600160401b031990921691909117905550505050565b6001600160401b03831660009081526002602052604080822090516102da90859085906111ba565b908152604051908190036020019020546001600160401b031690509392505050565b7f3a3950e13dd607cc37980db0ef14266c40d2bba9c01b2e44bfe549808883095d8160405161032b9190611282565b60405180910390a150565b6001600160401b0380831660009081526020818152604091829020845181549286015193860151909416600160a81b02600160a81b600160e81b0319931515600160a01b026001600160a81b03199093166001600160a01b03909516949094179190911791909116919091178155606082015182919060018201906103bb90826113c2565b5050505050565b60408051608080820183526000808352602080840182905283850182905260608085018190526001600160401b038781168452838352928690208651948501875280546001600160a01b0381168652600160a01b810460ff16151593860193909352600160a81b90920490921694830194909452600184018054939492939184019161044d90611337565b80601f016020809104026020016040519081016040528092919081815260200182805461047990611337565b80156104c65780601f1061049b576101008083540402835291602001916104c6565b820191906000526020600020905b8154815290600101906020018083116104a957829003601f168201915b5050505050815250509050919050565b80356001600160401b03811681146104ed57600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b60405160a081016001600160401b038111828210171561052a5761052a6104f2565b60405290565b60405161010081016001600160401b038111828210171561052a5761052a6104f2565b604080519081016001600160401b038111828210171561052a5761052a6104f2565b604051606081016001600160401b038111828210171561052a5761052a6104f2565b604051608081016001600160401b038111828210171561052a5761052a6104f2565b604051601f8201601f191681016001600160401b03811182821017156105e1576105e16104f2565b604052919050565b600082601f8301126105fa57600080fd5b81356001600160401b03811115610613576106136104f2565b610626601f8201601f19166020016105b9565b81815284602083860101111561063b57600080fd5b816020850160208301376000918101602001919091529392505050565b600080600080600060a0868803121561067057600080fd5b610679866104d6565b9450610687602087016104d6565b9350604086013592506060860135600481106106a257600080fd5b915060808601356001600160401b038111156106bd57600080fd5b6106c9888289016105e9565b9150509295509295909350565b600060a082840312156106e857600080fd5b6106f0610508565b905081358152610702602083016104d6565b6020820152610713604083016104d6565b6040820152610724606083016104d6565b6060820152610735608083016104d6565b608082015292915050565b6001600160a01b038116811461075557600080fd5b50565b80356104ed81610740565b60006001600160401b0382111561077c5761077c6104f2565b5060051b60200190565b600082601f83011261079757600080fd5b813560206107ac6107a783610763565b6105b9565b82815260059290921b840181019181810190868411156107cb57600080fd5b8286015b848110156108c55780356001600160401b03808211156107ef5760008081fd5b9088019060a0828b03601f19018113156108095760008081fd5b610811610508565b87840135838111156108235760008081fd5b6108318d8a838801016105e9565b825250604080850135848111156108485760008081fd5b6108568e8b838901016105e9565b8a840152506060808601358581111561086f5760008081fd5b61087d8f8c838a01016105e9565b838501525060809150818601358184015250828501359250838311156108a35760008081fd5b6108b18d8a858801016105e9565b9082015286525050509183019183016107cf565b509695505050505050565b600080604083850312156108e357600080fd5b6108ec836104d6565b915060208301356001600160401b038082111561090857600080fd5b90840190610180828703121561091d57600080fd5b610925610530565b61092f87846106d6565b815261093d60a08401610758565b602082015260c08301358281111561095457600080fd5b610960888286016105e9565b60408301525060e08301358281111561097857600080fd5b610984888286016105e9565b6060830152506101008301358281111561099d57600080fd5b6109a9888286016105e9565b6080830152506109bc6101208401610758565b60a082015261014083013560c0820152610160830135828111156109df57600080fd5b6109eb88828601610786565b60e0830152508093505050509250929050565b600060208284031215610a1057600080fd5b610a19826104d6565b9392505050565b60008083601f840112610a3257600080fd5b5081356001600160401b03811115610a4957600080fd5b602083019150836020828501011115610a6157600080fd5b9250929050565b60008060008060608587031215610a7e57600080fd5b610a87856104d6565b9350610a95602086016104d6565b925060408501356001600160401b03811115610ab057600080fd5b610abc87828801610a20565b95989497509550505050565b60008060408385031215610adb57600080fd5b610ae4836104d6565b9150610af2602084016104d6565b90509250929050565b600080600060408486031215610b1057600080fd5b610b19846104d6565b925060208401356001600160401b03811115610b3457600080fd5b610b4086828701610a20565b9497909650939450505050565b80356001600160e01b03811681146104ed57600080fd5b600082601f830112610b7557600080fd5b81356020610b856107a783610763565b82815260069290921b84018101918181019086841115610ba457600080fd5b8286015b848110156108c55760408189031215610bc15760008081fd5b610bc9610553565b610bd2826104d6565b8152610bdf858301610b4d565b81860152835291830191604001610ba8565b600082601f830112610c0257600080fd5b81356020610c126107a783610763565b82815260079290921b84018101918181019086841115610c3157600080fd5b8286015b848110156108c5578088036080811215610c4f5760008081fd5b610c57610575565b610c60836104d6565b8152604080601f1984011215610c765760008081fd5b610c7e610553565b9250610c8b8785016104d6565b8352610c988185016104d6565b8388015281870192909252606083013591810191909152835291830191608001610c35565b60006020808385031215610cd057600080fd5b82356001600160401b0380821115610ce757600080fd5b81850191506040808388031215610cfd57600080fd5b610d05610553565b833583811115610d1457600080fd5b84016040818a031215610d2657600080fd5b610d2e610553565b813585811115610d3d57600080fd5b8201601f81018b13610d4e57600080fd5b8035610d5c6107a782610763565b81815260069190911b8201890190898101908d831115610d7b57600080fd5b928a01925b82841015610dcb5787848f031215610d985760008081fd5b610da0610553565b8435610dab81610740565b8152610db8858d01610b4d565b818d0152825292870192908a0190610d80565b845250505081870135935084841115610de357600080fd5b610def8a858401610b64565b8188015282525083850135915082821115610e0957600080fd5b610e1588838601610bf1565b85820152809550505050505092915050565b60008060408385031215610e3a57600080fd5b610e43836104d6565b915060208301356001600160401b0380821115610e5f57600080fd5b9084019060808287031215610e7357600080fd5b610e7b610597565b8235610e8681610740565b815260208301358015158114610e9b57600080fd5b6020820152610eac604084016104d6565b6040820152606083013582811115610ec357600080fd5b610ecf888286016105e9565b6060830152508093505050509250929050565b6000815180845260005b81811015610f0857602081850181015186830182015201610eec565b506000602082860101526020601f19601f83011685010191505092915050565b602080825282516001600160a01b03168282015282015115156040808301919091528201516001600160401b0316606080830191909152820151608080830152600090610f7860a0840182610ee2565b949350505050565b600060048410610fa057634e487b7160e01b600052602160045260246000fd5b83825260406020830152610f786040830184610ee2565b6001600160a01b03169052565b600082825180855260208086019550808260051b84010181860160005b8481101561107657601f19868403018952815160a0815181865261100782870182610ee2565b915050858201518582038787015261101f8282610ee2565b915050604080830151868303828801526110398382610ee2565b925050506060808301518187015250608080830151925085820381870152506110628183610ee2565b9a86019a9450505090830190600101610fe1565b5090979650505050505050565b602081526110d0602082018351805182526020808201516001600160401b039081169184019190915260408083015182169084015260608083015182169084015260809182015116910152565b600060208301516110e460c0840182610fb7565b5060408301516101808060e08501526111016101a0850183610ee2565b91506060850151601f1980868503016101008701526111208483610ee2565b935060808701519150808685030161012087015261113e8483610ee2565b935060a08701519150611155610140870183610fb7565b60c087015161016087015260e087015191508086850301838701525061117b8382610fc4565b9695505050505050565b6001600160401b038181168382160190808211156111b357634e487b7160e01b600052601160045260246000fd5b5092915050565b8183823760009101908152919050565b60008151808452602080850194506020840160005b8381101561121857815180516001600160401b031688528301516001600160e01b031683880152604090960195908201906001016111df565b509495945050505050565b600081518084526020808501945080840160005b8381101561121857815180516001600160401b0390811689528482015180518216868b0152850151166040898101919091520151606088015260809096019590820190600101611237565b6000602080835283516040808386015260a0850182516040606088015281815180845260c0890191508683019350600092505b808310156112f057835180516001600160a01b031683528701516001600160e01b0316878301529286019260019290920191908401906112b5565b5093850151878503605f190160808901529361130c81866111ca565b945050505050818501519150601f1984820301604085015261132e8183611223565b95945050505050565b600181811c9082168061134b57607f821691505b60208210810361136b57634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156113bd576000816000526020600020601f850160051c8101602086101561139a5750805b601f850160051c820191505b818110156113b9578281556001016113a6565b5050505b505050565b81516001600160401b038111156113db576113db6104f2565b6113ef816113e98454611337565b84611371565b602080601f831160018114611424576000841561140c5750858301515b600019600386901b1c1916600185901b1785556113b9565b600085815260208120601f198616915b8281101561145357888601518255948401946001909101908401611434565b50858210156114715787850151600019600388901b60f8161c191681555b5050505050600190811b0190555056fea164736f6c6343000818000a",
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
	return common.HexToHash("0xcae3d9f68a9ed33d0a770e9bcc6fc25d4041dd6391e6cd8015546547a3c93140")
}

func (CCIPReaderTesterCommitReportAccepted) Topic() common.Hash {
	return common.HexToHash("0x3a3950e13dd607cc37980db0ef14266c40d2bba9c01b2e44bfe549808883095d")
}

func (CCIPReaderTesterExecutionStateChanged) Topic() common.Hash {
	return common.HexToHash("0x8c324ce1367b83031769f6a813e3bb4c117aba2185789d66b98b791405be6df2")
}

func (_CCIPReaderTester *CCIPReaderTester) Address() common.Address {
	return _CCIPReaderTester.address
}

type CCIPReaderTesterInterface interface {
	GetExpectedNextSequenceNumber(opts *bind.CallOpts, destChainSelector uint64) (uint64, error)

	GetInboundNonce(opts *bind.CallOpts, sourceChainSelector uint64, sender []byte) (uint64, error)

	GetSourceChainConfig(opts *bind.CallOpts, sourceChainSelector uint64) (OffRampSourceChainConfig, error)

	EmitCCIPSendRequested(opts *bind.TransactOpts, destChainSelector uint64, message InternalEVM2AnyRampMessage) (*types.Transaction, error)

	EmitCommitReportAccepted(opts *bind.TransactOpts, report OffRampCommitReport) (*types.Transaction, error)

	EmitExecutionStateChanged(opts *bind.TransactOpts, sourceChainSelector uint64, sequenceNumber uint64, messageId [32]byte, state uint8, returnData []byte) (*types.Transaction, error)

	SetDestChainSeqNr(opts *bind.TransactOpts, destChainSelector uint64, sequenceNumber uint64) (*types.Transaction, error)

	SetInboundNonce(opts *bind.TransactOpts, sourceChainSelector uint64, testNonce uint64, sender []byte) (*types.Transaction, error)

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
