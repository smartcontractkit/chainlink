// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package report_codec

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

type IRMNV2Signature struct {
	R [32]byte
	S [32]byte
}

type InternalAny2EVMRampMessage struct {
	Header       InternalRampMessageHeader
	Sender       []byte
	Data         []byte
	Receiver     common.Address
	GasLimit     *big.Int
	TokenAmounts []InternalRampTokenAmount
}

type InternalExecutionReportSingleChain struct {
	SourceChainSelector uint64
	Messages            []InternalAny2EVMRampMessage
	OffchainTokenData   [][][]byte
	Proofs              [][32]byte
	ProofFlagBits       *big.Int
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
	PriceUpdates  InternalPriceUpdates
	MerkleRoots   []InternalMerkleRoot
	RmnSignatures []IRMNV2Signature
}

var ReportCodecMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRampAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"maxSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structInternal.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIRMNV2.Signature[]\",\"name\":\"rmnSignatures\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structOffRamp.CommitReport\",\"name\":\"report\",\"type\":\"tuple\"}],\"name\":\"CommitReportDecoded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destExecData\",\"type\":\"bytes\"}],\"internalType\":\"structInternal.RampTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.Any2EVMRampMessage[]\",\"name\":\"messages\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[][]\",\"name\":\"offchainTokenData\",\"type\":\"bytes[][]\"},{\"internalType\":\"bytes32[]\",\"name\":\"proofs\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"proofFlagBits\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structInternal.ExecutionReportSingleChain[]\",\"name\":\"report\",\"type\":\"tuple[]\"}],\"name\":\"ExecuteReportDecoded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"}],\"name\":\"decodeCommitReport\",\"outputs\":[{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRampAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"maxSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structInternal.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIRMNV2.Signature[]\",\"name\":\"rmnSignatures\",\"type\":\"tuple[]\"}],\"internalType\":\"structOffRamp.CommitReport\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"}],\"name\":\"decodeExecuteReport\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destExecData\",\"type\":\"bytes\"}],\"internalType\":\"structInternal.RampTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.Any2EVMRampMessage[]\",\"name\":\"messages\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[][]\",\"name\":\"offchainTokenData\",\"type\":\"bytes[][]\"},{\"internalType\":\"bytes32[]\",\"name\":\"proofs\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"proofFlagBits\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.ExecutionReportSingleChain[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506113e5806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80636fb349561461003b578063f816ec6014610064575b600080fd5b61004e610049366004610231565b610084565b60405161005b91906104f2565b60405180910390f35b610077610072366004610231565b6100a0565b60405161005b9190610837565b60608180602001905181019061009a9190610e8f565b92915050565b6040805160a08101825260608082018181526080830182905282526020808301829052928201528251909161009a918401810190840161124f565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160a0810167ffffffffffffffff8111828210171561012d5761012d6100db565b60405290565b60405160c0810167ffffffffffffffff8111828210171561012d5761012d6100db565b6040805190810167ffffffffffffffff8111828210171561012d5761012d6100db565b6040516060810167ffffffffffffffff8111828210171561012d5761012d6100db565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156101e3576101e36100db565b604052919050565b600067ffffffffffffffff821115610205576102056100db565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b60006020828403121561024357600080fd5b813567ffffffffffffffff81111561025a57600080fd5b8201601f8101841361026b57600080fd5b803561027e610279826101eb565b61019c565b81815285602083850101111561029357600080fd5b81602084016020830137600091810160200191909152949350505050565b60005b838110156102cc5781810151838201526020016102b4565b50506000910152565b600081518084526102ed8160208601602086016102b1565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b600082825180855260208086019550808260051b84010181860160005b848110156103ef577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0868403018952815160a08151818652610380828701826102d5565b915050858201518582038787015261039882826102d5565b915050604080830151868303828801526103b283826102d5565b925050506060808301518187015250608080830151925085820381870152506103db81836102d5565b9a86019a945050509083019060010161033c565b5090979650505050505050565b6000828251808552602080860195506005818360051b8501018287016000805b868110156104a7577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe088850381018c5283518051808752908801908887019080891b88018a01865b8281101561049057858a830301845261047e8286516102d5565b948c0194938c01939150600101610464565b509e8a019e9750505093870193505060010161041c565b50919998505050505050505050565b60008151808452602080850194506020840160005b838110156104e7578151875295820195908201906001016104cb565b509495945050505050565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b828110156106da577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0888603018452815160a0860167ffffffffffffffff8083511688528883015160a08a8a015282815180855260c08b01915060c08160051b8c010194508b8301925060005b81811015610683577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff408c87030183528351805180518852868f820151168f890152866040820151166040890152866060820151166060890152866080820151166080890152508d81015161014060a08901526106066101408901826102d5565b9050604082015188820360c08a015261061f82826102d5565b915050606082015161064960e08a018273ffffffffffffffffffffffffffffffffffffffff169052565b50608082015161010089015260a0820151915087810361012089015261066f818361031f565b97505050928c0192918c0191600101610586565b50505050506040820151878203604089015261069f82826103fc565b915050606082015187820360608901526106b982826104b6565b60809384015198909301979097525094509285019290850190600101610519565b5092979650505050505050565b60008151808452602080850194506020840160005b838110156104e7578151805167ffffffffffffffff1688528301517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1683880152604090960195908201906001016106fc565b600082825180855260208086019550808260051b84010181860160005b848110156103ef577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0868403018952815160a067ffffffffffffffff8083511686528683015182888801526107bf838801826102d5565b60408581015184169089015260608086015190931692880192909252506080928301519290950191909152509783019790830190600101610768565b60008151808452602080850194506020840160005b838110156104e7578151805188528301518388015260409096019590820190600101610810565b602080825282516060838301528051604060808501819052815160c086018190526000949392840191859160e08801905b808410156108c5578451805173ffffffffffffffffffffffffffffffffffffffff1683528701517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1687830152938601936001939093019290820190610868565b50938501518785037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff800160a0890152936108ff81866106e7565b9450505050508185015191507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe080858303016040860152610940828461074b565b925060408601519150808584030160608601525061095e82826107fb565b95945050505050565b600067ffffffffffffffff821115610981576109816100db565b5060051b60200190565b805167ffffffffffffffff811681146109a357600080fd5b919050565b600060a082840312156109ba57600080fd5b6109c261010a565b9050815181526109d46020830161098b565b60208201526109e56040830161098b565b60408201526109f66060830161098b565b6060820152610a076080830161098b565b608082015292915050565b600082601f830112610a2357600080fd5b8151610a31610279826101eb565b818152846020838601011115610a4657600080fd5b610a578260208301602087016102b1565b949350505050565b805173ffffffffffffffffffffffffffffffffffffffff811681146109a357600080fd5b600082601f830112610a9457600080fd5b81516020610aa461027983610967565b82815260059290921b84018101918181019086841115610ac357600080fd5b8286015b84811015610bdd57805167ffffffffffffffff80821115610ae85760008081fd5b818901915060a0807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848d03011215610b215760008081fd5b610b2961010a565b8784015183811115610b3b5760008081fd5b610b498d8a83880101610a12565b82525060408085015184811115610b605760008081fd5b610b6e8e8b83890101610a12565b8a8401525060608086015185811115610b875760008081fd5b610b958f8c838a0101610a12565b83850152506080915081860151818401525082850151925083831115610bbb5760008081fd5b610bc98d8a85880101610a12565b908201528652505050918301918301610ac7565b509695505050505050565b600082601f830112610bf957600080fd5b81516020610c0961027983610967565b82815260059290921b84018101918181019086841115610c2857600080fd5b8286015b84811015610bdd57805167ffffffffffffffff80821115610c4d5760008081fd5b8189019150610140807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848d03011215610c875760008081fd5b610c8f610133565b610c9b8c8986016109a8565b815260c084015183811115610cb05760008081fd5b610cbe8d8a83880101610a12565b898301525060e084015183811115610cd65760008081fd5b610ce48d8a83880101610a12565b604083015250610cf76101008501610a5f565b60608201526101208401516080820152908301519082821115610d1a5760008081fd5b610d288c8984870101610a83565b60a08201528652505050918301918301610c2c565b600082601f830112610d4e57600080fd5b81516020610d5e61027983610967565b82815260059290921b84018101918181019086841115610d7d57600080fd5b8286015b84811015610bdd57805167ffffffffffffffff80821115610da157600080fd5b818901915089603f830112610db557600080fd5b85820151610dc561027982610967565b81815260059190911b830160400190878101908c831115610de557600080fd5b604085015b83811015610e1e57805185811115610e0157600080fd5b610e108f6040838a0101610a12565b845250918901918901610dea565b50875250505092840192508301610d81565b600082601f830112610e4157600080fd5b81516020610e5161027983610967565b8083825260208201915060208460051b870101935086841115610e7357600080fd5b602086015b84811015610bdd5780518352918301918301610e78565b60006020808385031215610ea257600080fd5b825167ffffffffffffffff80821115610eba57600080fd5b818501915085601f830112610ece57600080fd5b8151610edc61027982610967565b81815260059190911b83018401908481019088831115610efb57600080fd5b8585015b83811015610ff557805185811115610f1657600080fd5b860160a0818c037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0011215610f4b5760008081fd5b610f5361010a565b610f5e89830161098b565b815260408083015188811115610f745760008081fd5b610f828e8c83870101610be8565b8b8401525060608084015189811115610f9b5760008081fd5b610fa98f8d83880101610d3d565b8385015250608091508184015189811115610fc45760008081fd5b610fd28f8d83880101610e30565b918401919091525060a09290920151918101919091528352918601918601610eff565b5098975050505050505050565b80517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811681146109a357600080fd5b600082601f83011261103f57600080fd5b8151602061104f61027983610967565b82815260069290921b8401810191818101908684111561106e57600080fd5b8286015b84811015610bdd576040818903121561108b5760008081fd5b611093610156565b61109c8261098b565b81526110a9858301611002565b81860152835291830191604001611072565b600082601f8301126110cc57600080fd5b815160206110dc61027983610967565b82815260059290921b840181019181810190868411156110fb57600080fd5b8286015b84811015610bdd57805167ffffffffffffffff808211156111205760008081fd5b818901915060a0807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848d030112156111595760008081fd5b61116161010a565b61116c88850161098b565b8152604080850151848111156111825760008081fd5b6111908e8b83890101610a12565b8a84015250606093506111a484860161098b565b9082015260806111b585820161098b565b938201939093529201519082015283529183019183016110ff565b600082601f8301126111e157600080fd5b815160206111f161027983610967565b82815260069290921b8401810191818101908684111561121057600080fd5b8286015b84811015610bdd576040818903121561122d5760008081fd5b611235610156565b815181528482015185820152835291830191604001611214565b6000602080838503121561126257600080fd5b825167ffffffffffffffff8082111561127a57600080fd5b908401906060828703121561128e57600080fd5b611296610179565b8251828111156112a557600080fd5b830160408189038113156112b857600080fd5b6112c0610156565b8251858111156112cf57600080fd5b8301601f81018b136112e057600080fd5b80516112ee61027982610967565b81815260069190911b8201890190898101908d83111561130d57600080fd5b928a01925b8284101561135b5785848f03121561132a5760008081fd5b611332610156565b61133b85610a5f565b81526113488c8601611002565b818d0152825292850192908a0190611312565b84525050508287015191508482111561137357600080fd5b61137f8a83850161102e565b8188015283525050828401518281111561139857600080fd5b6113a4888286016110bb565b858301525060408301519350818411156113bd57600080fd5b6113c9878585016111d0565b6040820152969550505050505056fea164736f6c6343000818000a",
}

var ReportCodecABI = ReportCodecMetaData.ABI

var ReportCodecBin = ReportCodecMetaData.Bin

func DeployReportCodec(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ReportCodec, error) {
	parsed, err := ReportCodecMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ReportCodecBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ReportCodec{address: address, abi: *parsed, ReportCodecCaller: ReportCodecCaller{contract: contract}, ReportCodecTransactor: ReportCodecTransactor{contract: contract}, ReportCodecFilterer: ReportCodecFilterer{contract: contract}}, nil
}

type ReportCodec struct {
	address common.Address
	abi     abi.ABI
	ReportCodecCaller
	ReportCodecTransactor
	ReportCodecFilterer
}

type ReportCodecCaller struct {
	contract *bind.BoundContract
}

type ReportCodecTransactor struct {
	contract *bind.BoundContract
}

type ReportCodecFilterer struct {
	contract *bind.BoundContract
}

type ReportCodecSession struct {
	Contract     *ReportCodec
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ReportCodecCallerSession struct {
	Contract *ReportCodecCaller
	CallOpts bind.CallOpts
}

type ReportCodecTransactorSession struct {
	Contract     *ReportCodecTransactor
	TransactOpts bind.TransactOpts
}

type ReportCodecRaw struct {
	Contract *ReportCodec
}

type ReportCodecCallerRaw struct {
	Contract *ReportCodecCaller
}

type ReportCodecTransactorRaw struct {
	Contract *ReportCodecTransactor
}

func NewReportCodec(address common.Address, backend bind.ContractBackend) (*ReportCodec, error) {
	abi, err := abi.JSON(strings.NewReader(ReportCodecABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindReportCodec(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ReportCodec{address: address, abi: abi, ReportCodecCaller: ReportCodecCaller{contract: contract}, ReportCodecTransactor: ReportCodecTransactor{contract: contract}, ReportCodecFilterer: ReportCodecFilterer{contract: contract}}, nil
}

func NewReportCodecCaller(address common.Address, caller bind.ContractCaller) (*ReportCodecCaller, error) {
	contract, err := bindReportCodec(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ReportCodecCaller{contract: contract}, nil
}

func NewReportCodecTransactor(address common.Address, transactor bind.ContractTransactor) (*ReportCodecTransactor, error) {
	contract, err := bindReportCodec(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ReportCodecTransactor{contract: contract}, nil
}

func NewReportCodecFilterer(address common.Address, filterer bind.ContractFilterer) (*ReportCodecFilterer, error) {
	contract, err := bindReportCodec(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ReportCodecFilterer{contract: contract}, nil
}

func bindReportCodec(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ReportCodecMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ReportCodec *ReportCodecRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ReportCodec.Contract.ReportCodecCaller.contract.Call(opts, result, method, params...)
}

func (_ReportCodec *ReportCodecRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReportCodec.Contract.ReportCodecTransactor.contract.Transfer(opts)
}

func (_ReportCodec *ReportCodecRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReportCodec.Contract.ReportCodecTransactor.contract.Transact(opts, method, params...)
}

func (_ReportCodec *ReportCodecCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ReportCodec.Contract.contract.Call(opts, result, method, params...)
}

func (_ReportCodec *ReportCodecTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReportCodec.Contract.contract.Transfer(opts)
}

func (_ReportCodec *ReportCodecTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReportCodec.Contract.contract.Transact(opts, method, params...)
}

func (_ReportCodec *ReportCodecCaller) DecodeCommitReport(opts *bind.CallOpts, report []byte) (OffRampCommitReport, error) {
	var out []interface{}
	err := _ReportCodec.contract.Call(opts, &out, "decodeCommitReport", report)

	if err != nil {
		return *new(OffRampCommitReport), err
	}

	out0 := *abi.ConvertType(out[0], new(OffRampCommitReport)).(*OffRampCommitReport)

	return out0, err

}

func (_ReportCodec *ReportCodecSession) DecodeCommitReport(report []byte) (OffRampCommitReport, error) {
	return _ReportCodec.Contract.DecodeCommitReport(&_ReportCodec.CallOpts, report)
}

func (_ReportCodec *ReportCodecCallerSession) DecodeCommitReport(report []byte) (OffRampCommitReport, error) {
	return _ReportCodec.Contract.DecodeCommitReport(&_ReportCodec.CallOpts, report)
}

func (_ReportCodec *ReportCodecCaller) DecodeExecuteReport(opts *bind.CallOpts, report []byte) ([]InternalExecutionReportSingleChain, error) {
	var out []interface{}
	err := _ReportCodec.contract.Call(opts, &out, "decodeExecuteReport", report)

	if err != nil {
		return *new([]InternalExecutionReportSingleChain), err
	}

	out0 := *abi.ConvertType(out[0], new([]InternalExecutionReportSingleChain)).(*[]InternalExecutionReportSingleChain)

	return out0, err

}

func (_ReportCodec *ReportCodecSession) DecodeExecuteReport(report []byte) ([]InternalExecutionReportSingleChain, error) {
	return _ReportCodec.Contract.DecodeExecuteReport(&_ReportCodec.CallOpts, report)
}

func (_ReportCodec *ReportCodecCallerSession) DecodeExecuteReport(report []byte) ([]InternalExecutionReportSingleChain, error) {
	return _ReportCodec.Contract.DecodeExecuteReport(&_ReportCodec.CallOpts, report)
}

type ReportCodecCommitReportDecodedIterator struct {
	Event *ReportCodecCommitReportDecoded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ReportCodecCommitReportDecodedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReportCodecCommitReportDecoded)
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
		it.Event = new(ReportCodecCommitReportDecoded)
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

func (it *ReportCodecCommitReportDecodedIterator) Error() error {
	return it.fail
}

func (it *ReportCodecCommitReportDecodedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ReportCodecCommitReportDecoded struct {
	Report OffRampCommitReport
	Raw    types.Log
}

func (_ReportCodec *ReportCodecFilterer) FilterCommitReportDecoded(opts *bind.FilterOpts) (*ReportCodecCommitReportDecodedIterator, error) {

	logs, sub, err := _ReportCodec.contract.FilterLogs(opts, "CommitReportDecoded")
	if err != nil {
		return nil, err
	}
	return &ReportCodecCommitReportDecodedIterator{contract: _ReportCodec.contract, event: "CommitReportDecoded", logs: logs, sub: sub}, nil
}

func (_ReportCodec *ReportCodecFilterer) WatchCommitReportDecoded(opts *bind.WatchOpts, sink chan<- *ReportCodecCommitReportDecoded) (event.Subscription, error) {

	logs, sub, err := _ReportCodec.contract.WatchLogs(opts, "CommitReportDecoded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ReportCodecCommitReportDecoded)
				if err := _ReportCodec.contract.UnpackLog(event, "CommitReportDecoded", log); err != nil {
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

func (_ReportCodec *ReportCodecFilterer) ParseCommitReportDecoded(log types.Log) (*ReportCodecCommitReportDecoded, error) {
	event := new(ReportCodecCommitReportDecoded)
	if err := _ReportCodec.contract.UnpackLog(event, "CommitReportDecoded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ReportCodecExecuteReportDecodedIterator struct {
	Event *ReportCodecExecuteReportDecoded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ReportCodecExecuteReportDecodedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReportCodecExecuteReportDecoded)
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
		it.Event = new(ReportCodecExecuteReportDecoded)
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

func (it *ReportCodecExecuteReportDecodedIterator) Error() error {
	return it.fail
}

func (it *ReportCodecExecuteReportDecodedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ReportCodecExecuteReportDecoded struct {
	Report []InternalExecutionReportSingleChain
	Raw    types.Log
}

func (_ReportCodec *ReportCodecFilterer) FilterExecuteReportDecoded(opts *bind.FilterOpts) (*ReportCodecExecuteReportDecodedIterator, error) {

	logs, sub, err := _ReportCodec.contract.FilterLogs(opts, "ExecuteReportDecoded")
	if err != nil {
		return nil, err
	}
	return &ReportCodecExecuteReportDecodedIterator{contract: _ReportCodec.contract, event: "ExecuteReportDecoded", logs: logs, sub: sub}, nil
}

func (_ReportCodec *ReportCodecFilterer) WatchExecuteReportDecoded(opts *bind.WatchOpts, sink chan<- *ReportCodecExecuteReportDecoded) (event.Subscription, error) {

	logs, sub, err := _ReportCodec.contract.WatchLogs(opts, "ExecuteReportDecoded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ReportCodecExecuteReportDecoded)
				if err := _ReportCodec.contract.UnpackLog(event, "ExecuteReportDecoded", log); err != nil {
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

func (_ReportCodec *ReportCodecFilterer) ParseExecuteReportDecoded(log types.Log) (*ReportCodecExecuteReportDecoded, error) {
	event := new(ReportCodecExecuteReportDecoded)
	if err := _ReportCodec.contract.UnpackLog(event, "ExecuteReportDecoded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_ReportCodec *ReportCodec) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _ReportCodec.abi.Events["CommitReportDecoded"].ID:
		return _ReportCodec.ParseCommitReportDecoded(log)
	case _ReportCodec.abi.Events["ExecuteReportDecoded"].ID:
		return _ReportCodec.ParseExecuteReportDecoded(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (ReportCodecCommitReportDecoded) Topic() common.Hash {
	return common.HexToHash("0x31a4e1cb25733cdb9679561cd59cdc238d70a7d486f8bfc1f13242efd60fc29d")
}

func (ReportCodecExecuteReportDecoded) Topic() common.Hash {
	return common.HexToHash("0x70d042e9d8463eeac7e835e7172c08846ad327fc4eb6fc89f3bb5226e17ad618")
}

func (_ReportCodec *ReportCodec) Address() common.Address {
	return _ReportCodec.address
}

type ReportCodecInterface interface {
	DecodeCommitReport(opts *bind.CallOpts, report []byte) (OffRampCommitReport, error)

	DecodeExecuteReport(opts *bind.CallOpts, report []byte) ([]InternalExecutionReportSingleChain, error)

	FilterCommitReportDecoded(opts *bind.FilterOpts) (*ReportCodecCommitReportDecodedIterator, error)

	WatchCommitReportDecoded(opts *bind.WatchOpts, sink chan<- *ReportCodecCommitReportDecoded) (event.Subscription, error)

	ParseCommitReportDecoded(log types.Log) (*ReportCodecCommitReportDecoded, error)

	FilterExecuteReportDecoded(opts *bind.FilterOpts) (*ReportCodecExecuteReportDecodedIterator, error)

	WatchExecuteReportDecoded(opts *bind.WatchOpts, sink chan<- *ReportCodecExecuteReportDecoded) (event.Subscription, error)

	ParseExecuteReportDecoded(log types.Log) (*ReportCodecExecuteReportDecoded, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
