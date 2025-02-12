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
	PriceUpdates         InternalPriceUpdates
	BlessedMerkleRoots   []InternalMerkleRoot
	UnblessedMerkleRoots []InternalMerkleRoot
	RmnSignatures        []IRMNRemoteSignature
}

var ReportCodecMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"decodeCommitReport\",\"inputs\":[{\"name\":\"report\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.CommitReport\",\"components\":[{\"name\":\"priceUpdates\",\"type\":\"tuple\",\"internalType\":\"structInternal.PriceUpdates\",\"components\":[{\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"components\":[{\"name\":\"sourceToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"usdPerToken\",\"type\":\"uint224\",\"internalType\":\"uint224\"}]},{\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.GasPriceUpdate[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"usdPerUnitGas\",\"type\":\"uint224\",\"internalType\":\"uint224\"}]}]},{\"name\":\"blessedMerkleRoots\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.MerkleRoot[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRampAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"unblessedMerkleRoots\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.MerkleRoot[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRampAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"rmnSignatures\",\"type\":\"tuple[]\",\"internalType\":\"structIRMNRemote.Signature[]\",\"components\":[{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"decodeExecuteReport\",\"inputs\":[{\"name\":\"report\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.ExecutionReport[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"messages\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMRampMessage[]\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structInternal.RampMessageHeader\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destGasAmount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"offchainTokenData\",\"type\":\"bytes[][]\",\"internalType\":\"bytes[][]\"},{\"name\":\"proofs\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"proofFlagBits\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"pure\"},{\"type\":\"event\",\"name\":\"CommitReportDecoded\",\"inputs\":[{\"name\":\"report\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOffRamp.CommitReport\",\"components\":[{\"name\":\"priceUpdates\",\"type\":\"tuple\",\"internalType\":\"structInternal.PriceUpdates\",\"components\":[{\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"components\":[{\"name\":\"sourceToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"usdPerToken\",\"type\":\"uint224\",\"internalType\":\"uint224\"}]},{\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.GasPriceUpdate[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"usdPerUnitGas\",\"type\":\"uint224\",\"internalType\":\"uint224\"}]}]},{\"name\":\"blessedMerkleRoots\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.MerkleRoot[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRampAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"unblessedMerkleRoots\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.MerkleRoot[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRampAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"rmnSignatures\",\"type\":\"tuple[]\",\"internalType\":\"structIRMNRemote.Signature[]\",\"components\":[{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ExecuteReportDecoded\",\"inputs\":[{\"name\":\"report\",\"type\":\"tuple[]\",\"indexed\":false,\"internalType\":\"structInternal.ExecutionReport[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"messages\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMRampMessage[]\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structInternal.RampMessageHeader\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destGasAmount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"offchainTokenData\",\"type\":\"bytes[][]\",\"internalType\":\"bytes[][]\"},{\"name\":\"proofs\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"proofFlagBits\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"anonymous\":false}]",
	Bin: "0x6080806040523460155761168f908161001b8239f35b600080fdfe6101c0604052600436101561001357600080fd5b60003560e01c80636fb34956146104fc5763f816ec601461003357600080fd5b346104795761004136611329565b60608061004c61126b565b61005461128b565b82815282602082015281528160208201528160408201520152805181019060208201906020818403126104795760208101519067ffffffffffffffff8211610479570190608082840312610479576100aa61126b565b91602081015167ffffffffffffffff81116104795781016020810190604090860312610479576100d861128b565b90805167ffffffffffffffff811161047957810184601f8201121561047957805161010a610105826114d2565b6112ab565b9160208084848152019260061b8201019087821161047957602001915b8183106104bd57505050825260208101519067ffffffffffffffff8211610479570183601f82011215610479578051610162610105826114d2565b9160208084848152019260061b8201019086821161047957602001915b81831061047e5750505060208201528352604081015167ffffffffffffffff8111610479578260206101b392840101611584565b9360208401948552606082015167ffffffffffffffff8111610479578360206101de92850101611584565b9160408501928352608081015167ffffffffffffffff81116104795760209101019280601f850112156104795783519161021a610105846114d2565b9460208087868152019460061b82010192831161047957602001925b828410610448575050505060608301918252604051926020845251936080602085015260e0840194805195604060a087015286518091526020610100870197019060005b8181106103f25750505060200151947fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff608582030160c08601526020808751928381520196019060005b8181106103a85750505090610307849561033893517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe087830301604088015261140c565b90517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe085830301606086015261140c565b9051907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08382030160808401526020808351928381520192019060005b818110610383575050500390f35b8251805185526020908101518186015286955060409094019390920191600101610375565b8251805167ffffffffffffffff1689526020908101517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16818a0152604090980197909201916001016102c3565b8251805173ffffffffffffffffffffffffffffffffffffffff168a526020908101517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16818b01526040909901989092019160010161027a565b60406020858403011261047957602060409161046261128b565b865181528287015183820152815201930192610236565b600080fd5b60406020848b03011261047957602060409161049861128b565b6104a1866114ea565b81526104ae83870161155b565b8382015281520192019161017f565b60406020848c0301126104795760206040916104d761128b565b6104e08661153a565b81526104ed83870161155b565b83820152815201920191610127565b346104795761050a36611329565b61010052610100515161010051016101a0526020610100516101a05103126104795760206101005101516101405267ffffffffffffffff61014051116104795760206101a05101603f610140516101005101011215610479576020610140516101005101015161012052610583610105610120516114d2565b60c05260c051610180526101205160c05152602060c0510160c05260c0516101605260206101a051016020806101205160051b610140516101005101010101116104795760406101405161010051010160e0525b6020806101205160051b61014051610100510101010160e051106109a35760405180602081016020825261018051518091526040820160408260051b8401019161016051916000905b82821061062f57505050500390f35b91939092947fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc090820301825284519060a081019167ffffffffffffffff815116825260208101519260a06020840152835180915260c08301602060c08360051b86010195019160005b8181106107f357505050506040810151928281036040840152835180825260208201906020808260051b85010196019260005b8281106107385750505050506060810151928281036060840152602080855192838152019401906000905b80821061072057505050600192602092608080859401519101529601920192018594939192610620565b909194602080600192885181520196019201906106f6565b90919293967fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe083829e9d9e0301855287519081519081815260208101906020808460051b8301019401926000915b8183106107aa5750505050506020806001929e9d9e990195019101929190926106cb565b90919293946020806107e6837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0866001960301895289516113c9565b9701950193019190610786565b909192957fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff40868203018452865167ffffffffffffffff6080825180518552826020820151166020860152826040820151166040860152826060820151166060860152015116608083015260a061088e61087c6020840151610140848701526101408601906113c9565b604084015185820360c08701526113c9565b9173ffffffffffffffffffffffffffffffffffffffff60608201511660e0850152608081015161010085015201519161012081830391015281519081815260208101906020808460051b8301019401926000915b818310610902575050505050602080600192980194019101919091610698565b9091929394602080827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08560019503018852885190608080610990610950855160a0865260a08601906113c9565b73ffffffffffffffffffffffffffffffffffffffff87870151168786015263ffffffff6040870151166040860152606086015185820360608701526113c9565b93015191015297019501930191906108e2565b60e0515160a05267ffffffffffffffff60a051116104795760a060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081835161014051610100510101016101a0510301011261047957610a0261124b565b610a1c60208060a0516101405161010051010101016114ea565b81526040602060a0516101405161010051010101015160805267ffffffffffffffff608051116104795760206101a05101601f60206080518160a051610140516101005101010101010112156104795760206080518160a0516101405161010051010101010151610a8f610105826114d2565b90602082828152019060206101a0510160208260051b816080518160a05161014051610100510101010101010111610479576020806080518160a0516101405161010051010101010101915b60208260051b816080518160a0516101405161010051010101010101018310610dc65750505060208201526060602060a0516101405161010051010101015167ffffffffffffffff81116104795760206101a05101601f6020838160a05161014051610100510101010101011215610479576020818160a051610140516101005101010101015190610b6f610105836114d2565b91602083828152019160206101a0510160208360051b81848160a051610140516101005101010101010101116104795760a051610140516101005101018101606001925b60208360051b81848160a0516101405161010051010101010101018410610ca55750505050604082015260a0805161014051610100510101015167ffffffffffffffff8111610479576020908160a0516101405161010051010101010160206101a05101601f8201121561047957805190610c30610105836114d2565b9160208084838152019160051b8301019160206101a05101831161047957602001905b828210610c9557505050606082015260a06020815161014051610100510101010151608082015260c05152602060c0510160c052602060e0510160e0526105d7565b8151815260209182019101610c53565b835167ffffffffffffffff81116104795760206101a05101603f826020868160a05161014051610100510101010101010112156104795760a051610140516101005101018301810160600151610cfd610105826114d2565b916020838381520160206101a051016020808560051b85828b8160a0516101405161010051010101010101010101116104795760a0516101405161010051010186018201608001905b60a05161014051610100516080910190910188018401600586901b01018210610d7c575050509082525060209384019301610bb3565b815167ffffffffffffffff8111610479576101a05160a051610140516101005160209586959094610dbb949087019360809301018d01890101016114ff565b815201910190610d46565b825167ffffffffffffffff81116104795760207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082826080518160a05161014051610100510101010101016101a051030101906101408212610479576040519160c083019083821067ffffffffffffffff83111761121c5760a0916040521261047957610e5161124b565b602082816080518160a051610140516101005101010101010101518152610e9160408360206080518160a0516101405161010051010101010101016114ea565b6020820152610eb960608360206080518160a0516101405161010051010101010101016114ea565b6040820152610ee0608083602082518160a0516101405161010051010101010101016114ea565b6060820152610f0760a08360206080518184516101405161010051010101010101016114ea565b6080820152825260c08160206080518160a0516101405161010051010101010101015167ffffffffffffffff811161047957610f64906020806101a051019184826080518160a051610140516101005101010101010101016114ff565b602083015260e08160206080518160a0516101405161010051010101010101015167ffffffffffffffff811161047957610fbf906020806101a051019184826080518160a051610140516101005101010101010101016114ff565b6040830152610fe66101008260206080518160a0516101405186510101010101010161153a565b60608301526101208160206080518160a0516101405161010051010101010101015160808301526101408160206080518160a05185516101005101010101010101519067ffffffffffffffff82116104795760206101a05101601f60208484826080518160a051610140516101005101010101010101010112156104795760208282826080518160a051610140516101005101010101010101015161108d610105826114d2565b92602084838152019260206101a0510160208460051b818585826080518160a05161014051610100510101010101010101010111610479576080805160a05161014051610100510101018201830101935b6080805160a051610140516101005101010183018401600586901b0101851061111957505050505060a0820152815260209283019201610adb565b845167ffffffffffffffff81116104795760208484826080518160a051610140516101005101010101010101010160a060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0836101a051030101126104795761118161124b565b91602082015167ffffffffffffffff8111610479576111ab906020806101a05101918501016114ff565b83526111b96040830161153a565b6020840152606082015163ffffffff8116810361047957604084015260808201519267ffffffffffffffff84116104795760a060209493611205869586806101a05101918401016114ff565b6060840152015160808201528152019401936110de565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040519060a0820182811067ffffffffffffffff82111761121c57604052565b604051906080820182811067ffffffffffffffff82111761121c57604052565b604051906040820182811067ffffffffffffffff82111761121c57604052565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f604051930116820182811067ffffffffffffffff82111761121c57604052565b67ffffffffffffffff811161121c57601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b60207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8201126104795760043567ffffffffffffffff8111610479578160238201121561047957806004013590611382610105836112ef565b92828452602483830101116104795781600092602460209301838601378301015290565b60005b8381106113b95750506000910152565b81810151838201526020016113a9565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f602093611405815180928187528780880191016113a6565b0116010190565b9080602083519182815201916020808360051b8301019401926000915b83831061143857505050505090565b9091929394602080827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0856001950301865288519067ffffffffffffffff82511681526080806114958585015160a08786015260a08501906113c9565b9367ffffffffffffffff604082015116604085015267ffffffffffffffff6060820151166060850152015191015297019301930191939290611429565b67ffffffffffffffff811161121c5760051b60200190565b519067ffffffffffffffff8216820361047957565b81601f82011215610479578051611518610105826112ef565b92818452602082840101116104795761153791602080850191016113a6565b90565b519073ffffffffffffffffffffffffffffffffffffffff8216820361047957565b51907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8216820361047957565b81601f820112156104795780519061159e610105836114d2565b9260208085858152019360051b830101918183116104795760208101935b8385106115cb57505050505090565b845167ffffffffffffffff811161047957820160a07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082860301126104795761161261124b565b9161161f602083016114ea565b835260408201519267ffffffffffffffff84116104795760a08361164a8860208098819801016114ff565b8584015261165a606082016114ea565b604084015261166b608082016114ea565b6060840152015160808201528152019401936115bc56fea164736f6c634300081a000a",
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
	return common.HexToHash("0x58a603662478f1f3f18b8fc8006e127dc1ef28117b04eb2d89cc22849068dcd8")
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
