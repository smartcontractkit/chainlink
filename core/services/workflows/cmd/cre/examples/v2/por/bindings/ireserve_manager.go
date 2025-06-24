package bindings

import (
	_ "embed"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-common/pkg/values/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2"

	evmcappb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm"
)

// TODO figure out how we know the type...?

//go:embed solc/bin/IReserveManager.abi
var iReserveManagerRaw string

var iReserveManagerApi = NewIReserveManagerAbi()

func NewIReserveManagerAbi() abi.ABI {
	a, _ := abi.JSON(strings.NewReader(iReserveManagerRaw))
	return a
}

type IReserveManagerCodec interface {
}

type iReserveManagerCodec struct {
	abi abi.ABI
}

type IReserverManager struct {
	codec          IReserveManagerCodec
	ContractInputs ContractInputs
}

func NewIReserveManagerCodec() (IReserveManagerCodec, error) {
	return iReserveManagerCodec{abi: NewIReserveManagerAbi()}, nil
}

func NewIReserveManager(contracInputs ContractInputs) IReserverManager {
	codec, _ := NewIReserveManagerCodec()
	reserveManager := IReserverManager{ContractInputs: contracInputs, codec: codec}
	return reserveManager
}

type UpdateReserves struct {
	reserveManager *IReserverManager
}

type UpdateReservesStruct struct {
	TotalMinted  *big.Int
	TotalReserve *big.Int
}

func (irm IReserverManager) WriteReportUpdateReserves(runtime sdk.Runtime, updateReserves UpdateReservesStruct, options *WriteOptions) sdk.Promise[*evmcappb.WriteReportReply] {
	// Pack the complete function call (selector + args)
	body, err := iReserveManagerApi.Pack("updateReserves", updateReserves)
	if err != nil {
		return sdk.PromiseFromResult[*evmcappb.WriteReportReply](nil, err)
	}

	writeReportReplyPromise := irm.ContractInputs.EVM.WriteReport(runtime, &evmcappb.WriteReportRequest{
		Receiver: irm.ContractInputs.Address,
		Report: &evmcappb.SignedReport{
			RawReport:     body,
			ReportContext: []byte{},
			Signatures:    [][]byte{},
			Id:            []byte{},
		},
	})

	return writeReportReplyPromise
}

// func (irm IReserverManager) RequestReserveUpdateTrigger(confidence evmcappb.ConfidenceLevel) sdk.GenericTrigger[*evm.Log, *RequestReserveUpdateLog] {
// 	evmTrigger := irm.ContractInputs.EVM.LogTrigger(&evmcappb.FilterLogTriggerRequest{
// 		Addresses:  [][]byte{irm.ContractInputs.Address},
// 		EventSigs:  [][]byte{[]byte(iReserveManagerApi.Events["RequestReserveUpdate"].Sig)},
// 		Confidence: confidence,
// 	})
// 	return &requestReserveUpdateLogTrigger{Trigger: evmTrigger}
// }

// Someone should review the helpers we generate.

type RequestReserveUpdateLog struct {
	// No topics in this event except the hash, should we expose it or verify it?
	TxHash        common.Hash
	BlockHash     common.Hash
	BlockNumber   *pb.BigInt
	TxIndex       uint32
	Index         uint32
	Removed       bool
	ChainSelector uint32
	RequestId     *big.Int
}

// type requestReserveUpdateLogTrigger struct {
// 	sdk.GenericTrigger[*evm.Log, *evm.Log]
// }

// func (r requestReserveUpdateLogTrigger) Adapt(m *evm.Log) (*RequestReserveUpdateLog, error) {
// 	data, err := iReserveManagerApi.Events["RequestReserveUpdate"].Inputs.Unpack(m.Data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	requestId := data[0].(*big.Int)
// 	return &RequestReserveUpdateLog{
// 		TxHash:        common.BytesToHash(m.TxHash),
// 		BlockHash:     common.BytesToHash(m.BlockHash),
// 		BlockNumber:   m.BlockNumber,
// 		TxIndex:       m.TxIndex,
// 		Index:         m.Index,
// 		Removed:       m.Removed,
// 		ChainSelector: m.ChainSelector,
// 		RequestId:     requestId,
// 	}, nil
// }
