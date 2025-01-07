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

type IRMNRemoteSignature struct {
	R [32]byte
	S [32]byte
}

type InternalAny2EVMRampMessage struct {
	Header       InternalRampMessageHeader
	Sender       []byte
	Data         []byte
	Receiver     common.Address
	GasLimit     *big.Int
	TokenAmounts []InternalAny2EVMTokenTransfer
}

type InternalAny2EVMTokenTransfer struct {
	SourcePoolAddress []byte
	DestTokenAddress  common.Address
	DestGasAmount     uint32
	ExtraData         []byte
	Amount            *big.Int
}

type InternalExecutionReport struct {
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

type InternalTokenPriceUpdate struct {
	SourceToken common.Address
	UsdPerToken *big.Int
}

type OffRampCommitReport struct {
	PriceUpdates  InternalPriceUpdates
	MerkleRoots   []InternalMerkleRoot
	RmnSignatures []IRMNRemoteSignature
}

var ReportCodecMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"decodeCommitReport\",\"inputs\":[{\"name\":\"report\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.CommitReport\",\"components\":[{\"name\":\"priceUpdates\",\"type\":\"tuple\",\"internalType\":\"structInternal.PriceUpdates\",\"components\":[{\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"components\":[{\"name\":\"sourceToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"usdPerToken\",\"type\":\"uint224\",\"internalType\":\"uint224\"}]},{\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.GasPriceUpdate[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"usdPerUnitGas\",\"type\":\"uint224\",\"internalType\":\"uint224\"}]}]},{\"name\":\"merkleRoots\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.MerkleRoot[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRampAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"rmnSignatures\",\"type\":\"tuple[]\",\"internalType\":\"structIRMNRemote.Signature[]\",\"components\":[{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"decodeExecuteReport\",\"inputs\":[{\"name\":\"report\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.ExecutionReport[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"messages\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMRampMessage[]\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structInternal.RampMessageHeader\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destGasAmount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"offchainTokenData\",\"type\":\"bytes[][]\",\"internalType\":\"bytes[][]\"},{\"name\":\"proofs\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"proofFlagBits\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"pure\"},{\"type\":\"event\",\"name\":\"CommitReportDecoded\",\"inputs\":[{\"name\":\"report\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOffRamp.CommitReport\",\"components\":[{\"name\":\"priceUpdates\",\"type\":\"tuple\",\"internalType\":\"structInternal.PriceUpdates\",\"components\":[{\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"components\":[{\"name\":\"sourceToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"usdPerToken\",\"type\":\"uint224\",\"internalType\":\"uint224\"}]},{\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.GasPriceUpdate[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"usdPerUnitGas\",\"type\":\"uint224\",\"internalType\":\"uint224\"}]}]},{\"name\":\"merkleRoots\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.MerkleRoot[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRampAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"rmnSignatures\",\"type\":\"tuple[]\",\"internalType\":\"structIRMNRemote.Signature[]\",\"components\":[{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ExecuteReportDecoded\",\"inputs\":[{\"name\":\"report\",\"type\":\"tuple[]\",\"indexed\":false,\"internalType\":\"structInternal.ExecutionReport[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"messages\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMRampMessage[]\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structInternal.RampMessageHeader\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destGasAmount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"offchainTokenData\",\"type\":\"bytes[][]\",\"internalType\":\"bytes[][]\"},{\"name\":\"proofs\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"proofFlagBits\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"anonymous\":false}]",
	Bin: "0x60808060405234601557611604908161001b8239f35b600080fdfe6101c0604052600436101561001357600080fd5b60003560e01c80636fb34956146106355763f816ec601461003357600080fd5b3461051d5761004136611462565b6060604061004d6113a4565b6100556113c4565b838152836020820152815282602082015201528051810190602082019060208184031261051d5760208101519067ffffffffffffffff821161051d57019060608284031261051d576100a56113a4565b92602083015167ffffffffffffffff811161051d578301602081019060409083031261051d576100d36113c4565b90805167ffffffffffffffff811161051d57810184601f8201121561051d57805161010561010082611545565b6113e4565b9160208084848152019260061b8201019087821161051d57602001915b8183106105f657505050825260208101519067ffffffffffffffff821161051d570183601f8201121561051d57805161015d61010082611545565b9160208084848152019260061b8201019086821161051d57602001915b8183106105b75750505060208201528452604083015167ffffffffffffffff811161051d576020908401019282601f8501121561051d578351936101c061010086611545565b9460208087838152019160051b8301019185831161051d5760208101915b838310610522575050505060208501938452606081015167ffffffffffffffff811161051d5760209101019180601f8401121561051d5782519161022461010084611545565b9360208086868152019460061b82010192831161051d5793959493602001925b8284106104e45750505050604082019283526040519283926020845251916060602085015260c0840192805193604060808701528451809152602060e0870195019060005b81811061048b5750505060200151927fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff808582030160a08601526020808551928381520194019060005b81811061043e5750505051917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848203016040850152825180825260208201916020808360051b8301019501926000915b83831061039e57505050505051907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08382030160608401526020808351928381520192019060005b818110610379575050500390f35b825180518552602090810151818601528695506040909401939092019160010161036b565b91939650919394602080827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0856001950301865289519067ffffffffffffffff82511681526080806103fd8585015160a08786015260a0850190611502565b9367ffffffffffffffff604082015116604085015267ffffffffffffffff606082015116606085015201519101529801930193019092879695949293610323565b8251805167ffffffffffffffff1687526020908101517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1681880152889750604090960195909201916001016102d2565b8251805173ffffffffffffffffffffffffffffffffffffffff1688526020908101517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff168189015289985060409097019690920191600101610289565b6040602085849997989903011261051d5760206040916105026113c4565b86518152828701518382015281520193019295949395610244565b600080fd5b825167ffffffffffffffff811161051d57820160a08188031261051d57610547611384565b916105546020830161155d565b835260408201519267ffffffffffffffff841161051d5760a08361057f8c6020809881980101611572565b8584015261058f6060820161155d565b60408401526105a06080820161155d565b6060840152015160808201528152019201916101de565b60406020848803011261051d5760206040916105d16113c4565b6105da8661155d565b81526105e78387016115ce565b8382015281520192019161017a565b60406020848903011261051d5760206040916106106113c4565b610619866115ad565b81526106268387016115ce565b83820152815201920191610122565b3461051d5761064336611462565b61010052610100515161010051016101a0526020610100516101a051031261051d5760206101005101516101405267ffffffffffffffff610140511161051d5760206101a05101603f61014051610100510101121561051d5760206101405161010051010151610120526106bc61010061012051611545565b60c05260c051610180526101205160c05152602060c0510160c05260c0516101605260206101a051016020806101205160051b6101405161010051010101011161051d5760406101405161010051010160e0525b6020806101205160051b61014051610100510101010160e05110610adc5760405180602081016020825261018051518091526040820160408260051b8401019161016051916000905b82821061076857505050500390f35b91939092947fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc090820301825284519060a081019167ffffffffffffffff815116825260208101519260a06020840152835180915260c08301602060c08360051b86010195019160005b81811061092c57505050506040810151928281036040840152835180825260208201906020808260051b85010196019260005b8281106108715750505050506060810151928281036060840152602080855192838152019401906000905b80821061085957505050600192602092608080859401519101529601920192018594939192610759565b9091946020806001928851815201960192019061082f565b90919293967fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe083829e9d9e0301855287519081519081815260208101906020808460051b8301019401926000915b8183106108e35750505050506020806001929e9d9e99019501910192919092610804565b909192939460208061091f837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe086600196030189528951611502565b97019501930191906108bf565b909192957fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff40868203018452865167ffffffffffffffff6080825180518552826020820151166020860152826040820151166040860152826060820151166060860152015116608083015260a06109c76109b5602084015161014084870152610140860190611502565b604084015185820360c0870152611502565b9173ffffffffffffffffffffffffffffffffffffffff60608201511660e0850152608081015161010085015201519161012081830391015281519081815260208101906020808460051b8301019401926000915b818310610a3b5750505050506020806001929801940191019190916107d1565b9091929394602080827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08560019503018852885190608080610ac9610a89855160a0865260a0860190611502565b73ffffffffffffffffffffffffffffffffffffffff87870151168786015263ffffffff604087015116604086015260608601518582036060870152611502565b9301519101529701950193019190610a1b565b60e0515160a05267ffffffffffffffff60a0511161051d5760a060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081835161014051610100510101016101a0510301011261051d57610b3b611384565b610b5560208060a05161014051610100510101010161155d565b81526040602060a0516101405161010051010101015160805267ffffffffffffffff6080511161051d5760206101a05101601f60206080518160a0516101405161010051010101010101121561051d5760206080518160a0516101405161010051010101010151610bc861010082611545565b90602082828152019060206101a0510160208260051b816080518160a0516101405161010051010101010101011161051d576020806080518160a0516101405161010051010101010101915b60208260051b816080518160a0516101405161010051010101010101018310610eff5750505060208201526060602060a0516101405161010051010101015167ffffffffffffffff811161051d5760206101a05101601f6020838160a0516101405161010051010101010101121561051d576020818160a051610140516101005101010101015190610ca861010083611545565b91602083828152019160206101a0510160208360051b81848160a0516101405161010051010101010101011161051d5760a051610140516101005101018101606001925b60208360051b81848160a0516101405161010051010101010101018410610dde5750505050604082015260a0805161014051610100510101015167ffffffffffffffff811161051d576020908160a0516101405161010051010101010160206101a05101601f8201121561051d57805190610d6961010083611545565b9160208084838152019160051b8301019160206101a05101831161051d57602001905b828210610dce57505050606082015260a06020815161014051610100510101010151608082015260c05152602060c0510160c052602060e0510160e052610710565b8151815260209182019101610d8c565b835167ffffffffffffffff811161051d5760206101a05101603f826020868160a051610140516101005101010101010101121561051d5760a051610140516101005101018301810160600151610e3661010082611545565b916020838381520160206101a051016020808560051b85828b8160a05161014051610100510101010101010101011161051d5760a0516101405161010051010186018201608001905b60a05161014051610100516080910190910188018401600586901b01018210610eb5575050509082525060209384019301610cec565b815167ffffffffffffffff811161051d576101a05160a051610140516101005160209586959094610ef4949087019360809301018d0189010101611572565b815201910190610e7f565b825167ffffffffffffffff811161051d5760207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082826080518160a05161014051610100510101010101016101a05103010190610140821261051d576040519160c083019083821067ffffffffffffffff8311176113555760a0916040521261051d57610f8a611384565b602082816080518160a051610140516101005101010101010101518152610fca60408360206080518160a05161014051610100510101010101010161155d565b6020820152610ff260608360206080518160a05161014051610100510101010101010161155d565b6040820152611019608083602082518160a05161014051610100510101010101010161155d565b606082015261104060a083602060805181845161014051610100510101010101010161155d565b6080820152825260c08160206080518160a0516101405161010051010101010101015167ffffffffffffffff811161051d5761109d906020806101a051019184826080518160a05161014051610100510101010101010101611572565b602083015260e08160206080518160a0516101405161010051010101010101015167ffffffffffffffff811161051d576110f8906020806101a051019184826080518160a05161014051610100510101010101010101611572565b604083015261111f6101008260206080518160a051610140518651010101010101016115ad565b60608301526101208160206080518160a0516101405161010051010101010101015160808301526101408160206080518160a05185516101005101010101010101519067ffffffffffffffff821161051d5760206101a05101601f60208484826080518160a0516101405161010051010101010101010101121561051d5760208282826080518160a05161014051610100510101010101010101516111c661010082611545565b92602084838152019260206101a0510160208460051b818585826080518160a0516101405161010051010101010101010101011161051d576080805160a05161014051610100510101018201830101935b6080805160a051610140516101005101010183018401600586901b0101851061125257505050505060a0820152815260209283019201610c14565b845167ffffffffffffffff811161051d5760208484826080518160a051610140516101005101010101010101010160a060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0836101a0510301011261051d576112ba611384565b91602082015167ffffffffffffffff811161051d576112e4906020806101a0510191850101611572565b83526112f2604083016115ad565b6020840152606082015163ffffffff8116810361051d57604084015260808201519267ffffffffffffffff841161051d5760a06020949361133e869586806101a0510191840101611572565b606084015201516080820152815201940193611217565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040519060a0820182811067ffffffffffffffff82111761135557604052565b604051906060820182811067ffffffffffffffff82111761135557604052565b604051906040820182811067ffffffffffffffff82111761135557604052565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f604051930116820182811067ffffffffffffffff82111761135557604052565b67ffffffffffffffff811161135557601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b60207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc82011261051d5760043567ffffffffffffffff811161051d578160238201121561051d578060040135906114bb61010083611428565b928284526024838301011161051d5781600092602460209301838601378301015290565b60005b8381106114f25750506000910152565b81810151838201526020016114e2565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f60209361153e815180928187528780880191016114df565b0116010190565b67ffffffffffffffff81116113555760051b60200190565b519067ffffffffffffffff8216820361051d57565b81601f8201121561051d57805161158b61010082611428565b928184526020828401011161051d576115aa91602080850191016114df565b90565b519073ffffffffffffffffffffffffffffffffffffffff8216820361051d57565b51907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8216820361051d5756fea164736f6c634300081a000a",
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

func (_ReportCodec *ReportCodecCaller) DecodeExecuteReport(opts *bind.CallOpts, report []byte) ([]InternalExecutionReport, error) {
	var out []interface{}
	err := _ReportCodec.contract.Call(opts, &out, "decodeExecuteReport", report)

	if err != nil {
		return *new([]InternalExecutionReport), err
	}

	out0 := *abi.ConvertType(out[0], new([]InternalExecutionReport)).(*[]InternalExecutionReport)

	return out0, err

}

func (_ReportCodec *ReportCodecSession) DecodeExecuteReport(report []byte) ([]InternalExecutionReport, error) {
	return _ReportCodec.Contract.DecodeExecuteReport(&_ReportCodec.CallOpts, report)
}

func (_ReportCodec *ReportCodecCallerSession) DecodeExecuteReport(report []byte) ([]InternalExecutionReport, error) {
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
	Report []InternalExecutionReport
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
	return common.HexToHash("0x9467c8093a35a72f74398d5b6e351d67dc82eddc378efc6177eafb4fc7a01d39")
}

func (_ReportCodec *ReportCodec) Address() common.Address {
	return _ReportCodec.address
}

type ReportCodecInterface interface {
	DecodeCommitReport(opts *bind.CallOpts, report []byte) (OffRampCommitReport, error)

	DecodeExecuteReport(opts *bind.CallOpts, report []byte) ([]InternalExecutionReport, error)

	FilterCommitReportDecoded(opts *bind.FilterOpts) (*ReportCodecCommitReportDecodedIterator, error)

	WatchCommitReportDecoded(opts *bind.WatchOpts, sink chan<- *ReportCodecCommitReportDecoded) (event.Subscription, error)

	ParseCommitReportDecoded(log types.Log) (*ReportCodecCommitReportDecoded, error)

	FilterExecuteReportDecoded(opts *bind.FilterOpts) (*ReportCodecExecuteReportDecodedIterator, error)

	WatchExecuteReportDecoded(opts *bind.WatchOpts, sink chan<- *ReportCodecExecuteReportDecoded) (event.Subscription, error)

	ParseExecuteReportDecoded(log types.Log) (*ReportCodecExecuteReportDecoded, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
