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

var ReportCodecMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"min\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"max\",\"type\":\"uint64\"}],\"internalType\":\"structOffRamp.Interval\",\"name\":\"interval\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structOffRamp.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structOffRamp.CommitReport\",\"name\":\"report\",\"type\":\"tuple\"}],\"name\":\"CommitReportDecoded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destExecData\",\"type\":\"bytes\"}],\"internalType\":\"structInternal.RampTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.Any2EVMRampMessage[]\",\"name\":\"messages\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[][]\",\"name\":\"offchainTokenData\",\"type\":\"bytes[][]\"},{\"internalType\":\"bytes32[]\",\"name\":\"proofs\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"proofFlagBits\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structInternal.ExecutionReportSingleChain[]\",\"name\":\"report\",\"type\":\"tuple[]\"}],\"name\":\"ExecuteReportDecoded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"}],\"name\":\"decodeCommitReport\",\"outputs\":[{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"min\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"max\",\"type\":\"uint64\"}],\"internalType\":\"structOffRamp.Interval\",\"name\":\"interval\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structOffRamp.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"}],\"internalType\":\"structOffRamp.CommitReport\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"}],\"name\":\"decodeExecuteReport\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destExecData\",\"type\":\"bytes\"}],\"internalType\":\"structInternal.RampTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.Any2EVMRampMessage[]\",\"name\":\"messages\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[][]\",\"name\":\"offchainTokenData\",\"type\":\"bytes[][]\"},{\"internalType\":\"bytes32[]\",\"name\":\"proofs\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"proofFlagBits\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.ExecutionReportSingleChain[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061126d806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80636fb349561461003b578063f816ec6014610064575b600080fd5b61004e61004936600461022c565b610084565b60405161005b91906104ed565b60405180910390f35b61007761007236600461022c565b6100a0565b60405161005b91906107a6565b60608180602001905181019061009a9190610de1565b92915050565b604080516080810182526060918101828152828201839052815260208101919091528180602001905181019061009a91906110f7565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160a0810167ffffffffffffffff81118282101715610128576101286100d6565b60405290565b60405160c0810167ffffffffffffffff81118282101715610128576101286100d6565b6040805190810167ffffffffffffffff81118282101715610128576101286100d6565b6040516060810167ffffffffffffffff81118282101715610128576101286100d6565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156101de576101de6100d6565b604052919050565b600067ffffffffffffffff821115610200576102006100d6565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b60006020828403121561023e57600080fd5b813567ffffffffffffffff81111561025557600080fd5b8201601f8101841361026657600080fd5b8035610279610274826101e6565b610197565b81815285602083850101111561028e57600080fd5b81602084016020830137600091810160200191909152949350505050565b60005b838110156102c75781810151838201526020016102af565b50506000910152565b600081518084526102e88160208601602086016102ac565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b600082825180855260208086019550808260051b84010181860160005b848110156103ea577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0868403018952815160a0815181865261037b828701826102d0565b915050858201518582038787015261039382826102d0565b915050604080830151868303828801526103ad83826102d0565b925050506060808301518187015250608080830151925085820381870152506103d681836102d0565b9a86019a9450505090830190600101610337565b5090979650505050505050565b6000828251808552602080860195506005818360051b8501018287016000805b868110156104a2577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe088850381018c5283518051808752908801908887019080891b88018a01865b8281101561048b57858a83030184526104798286516102d0565b948c0194938c0193915060010161045f565b509e8a019e97505050938701935050600101610417565b50919998505050505050505050565b60008151808452602080850194506020840160005b838110156104e2578151875295820195908201906001016104c6565b509495945050505050565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b828110156106d5577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0888603018452815160a0860167ffffffffffffffff8083511688528883015160a08a8a015282815180855260c08b01915060c08160051b8c010194508b8301925060005b8181101561067e577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff408c87030183528351805180518852868f820151168f890152866040820151166040890152866060820151166060890152866080820151166080890152508d81015161014060a08901526106016101408901826102d0565b9050604082015188820360c08a015261061a82826102d0565b915050606082015161064460e08a018273ffffffffffffffffffffffffffffffffffffffff169052565b50608082015161010089015260a0820151915087810361012089015261066a818361031a565b97505050928c0192918c0191600101610581565b50505050506040820151878203604089015261069a82826103f7565b915050606082015187820360608901526106b482826104b1565b60809384015198909301979097525094509285019290850190600101610514565b5092979650505050505050565b60008151808452602080850194506020840160005b838110156104e2578151805167ffffffffffffffff1688528301517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1683880152604090960195908201906001016106f7565b600081518084526020808501945080840160005b838110156104e2578151805167ffffffffffffffff90811689528482015180518216868b0152850151166040808a0191909152015160608801526080909601959082019060010161075a565b6000602080835283516040808386015260a0850182516040606088015281815180845260c0890191508683019350600092505b80831015610836578351805173ffffffffffffffffffffffffffffffffffffffff1683528701517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16878301529286019260019290920191908401906107d9565b50938501518785037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa00160808901529361087081866106e2565b9450505050508185015191507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08482030160408501526108b08183610746565b95945050505050565b600067ffffffffffffffff8211156108d3576108d36100d6565b5060051b60200190565b805167ffffffffffffffff811681146108f557600080fd5b919050565b600060a0828403121561090c57600080fd5b610914610105565b905081518152610926602083016108dd565b6020820152610937604083016108dd565b6040820152610948606083016108dd565b6060820152610959608083016108dd565b608082015292915050565b600082601f83011261097557600080fd5b8151610983610274826101e6565b81815284602083860101111561099857600080fd5b6109a98260208301602087016102ac565b949350505050565b805173ffffffffffffffffffffffffffffffffffffffff811681146108f557600080fd5b600082601f8301126109e657600080fd5b815160206109f6610274836108b9565b82815260059290921b84018101918181019086841115610a1557600080fd5b8286015b84811015610b2f57805167ffffffffffffffff80821115610a3a5760008081fd5b818901915060a0807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848d03011215610a735760008081fd5b610a7b610105565b8784015183811115610a8d5760008081fd5b610a9b8d8a83880101610964565b82525060408085015184811115610ab25760008081fd5b610ac08e8b83890101610964565b8a8401525060608086015185811115610ad95760008081fd5b610ae78f8c838a0101610964565b83850152506080915081860151818401525082850151925083831115610b0d5760008081fd5b610b1b8d8a85880101610964565b908201528652505050918301918301610a19565b509695505050505050565b600082601f830112610b4b57600080fd5b81516020610b5b610274836108b9565b82815260059290921b84018101918181019086841115610b7a57600080fd5b8286015b84811015610b2f57805167ffffffffffffffff80821115610b9f5760008081fd5b8189019150610140807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848d03011215610bd95760008081fd5b610be161012e565b610bed8c8986016108fa565b815260c084015183811115610c025760008081fd5b610c108d8a83880101610964565b898301525060e084015183811115610c285760008081fd5b610c368d8a83880101610964565b604083015250610c4961010085016109b1565b60608201526101208401516080820152908301519082821115610c6c5760008081fd5b610c7a8c89848701016109d5565b60a08201528652505050918301918301610b7e565b600082601f830112610ca057600080fd5b81516020610cb0610274836108b9565b82815260059290921b84018101918181019086841115610ccf57600080fd5b8286015b84811015610b2f57805167ffffffffffffffff80821115610cf357600080fd5b818901915089603f830112610d0757600080fd5b85820151610d17610274826108b9565b81815260059190911b830160400190878101908c831115610d3757600080fd5b604085015b83811015610d7057805185811115610d5357600080fd5b610d628f6040838a0101610964565b845250918901918901610d3c565b50875250505092840192508301610cd3565b600082601f830112610d9357600080fd5b81516020610da3610274836108b9565b8083825260208201915060208460051b870101935086841115610dc557600080fd5b602086015b84811015610b2f5780518352918301918301610dca565b60006020808385031215610df457600080fd5b825167ffffffffffffffff80821115610e0c57600080fd5b818501915085601f830112610e2057600080fd5b8151610e2e610274826108b9565b81815260059190911b83018401908481019088831115610e4d57600080fd5b8585015b83811015610f4757805185811115610e6857600080fd5b860160a0818c037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0011215610e9d5760008081fd5b610ea5610105565b610eb08983016108dd565b815260408083015188811115610ec65760008081fd5b610ed48e8c83870101610b3a565b8b8401525060608084015189811115610eed5760008081fd5b610efb8f8d83880101610c8f565b8385015250608091508184015189811115610f165760008081fd5b610f248f8d83880101610d82565b918401919091525060a09290920151918101919091528352918601918601610e51565b5098975050505050505050565b80517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff811681146108f557600080fd5b600082601f830112610f9157600080fd5b81516020610fa1610274836108b9565b82815260069290921b84018101918181019086841115610fc057600080fd5b8286015b84811015610b2f5760408189031215610fdd5760008081fd5b610fe5610151565b610fee826108dd565b8152610ffb858301610f54565b81860152835291830191604001610fc4565b600082601f83011261101e57600080fd5b8151602061102e610274836108b9565b82815260079290921b8401810191818101908684111561104d57600080fd5b8286015b84811015610b2f57808803608081121561106b5760008081fd5b611073610174565b61107c836108dd565b81526040807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0840112156110b05760008081fd5b6110b8610151565b92506110c58785016108dd565b83526110d28185016108dd565b8388015281870192909252606083015191810191909152835291830191608001611051565b6000602080838503121561110a57600080fd5b825167ffffffffffffffff8082111561112257600080fd5b8185019150604080838803121561113857600080fd5b611140610151565b83518381111561114f57600080fd5b84016040818a03121561116157600080fd5b611169610151565b81518581111561117857600080fd5b8201601f81018b1361118957600080fd5b8051611197610274826108b9565b81815260069190911b8201890190898101908d8311156111b657600080fd5b928a01925b828410156112045787848f0312156111d35760008081fd5b6111db610151565b6111e4856109b1565b81526111f18c8601610f54565b818d0152825292870192908a01906111bb565b84525050508187015193508484111561121c57600080fd5b6112288a858401610f80565b818801528252508385015191508282111561124257600080fd5b61124e8883860161100d565b8582015280955050505050509291505056fea164736f6c6343000818000a",
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
	return common.HexToHash("0x1b2cb5e9d31bdaabb2ae07532436ae669406f84003ca27179b4dfb72f127f7dc")
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
