package bindings

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	evmmock "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm/capabilitymock"
	"github.com/smartcontractkit/chainlink-common/pkg/chains/evm"

	evmcappb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm"
	chain_common "github.com/smartcontractkit/chainlink-common/pkg/loop/chain-common"
)

// Define a custom error type
type TxFatalError struct {
	Message string
}

// Implement the error interface
func (e *TxFatalError) Error() string {
	return fmt.Sprintf("Error %s", e.Message)
}

// Define a custom error type
type ReceiverContractError struct {
	Message string
	TxHash  *[]byte
}

// Implement the error interface
func (e *ReceiverContractError) Error() string {
	return fmt.Sprintf("Error %s", e.Message)
}

type ContractOptions struct {
	GasConfig *evm.GasConfig
}

type ContractInputs struct {
	EVM     *evmcappb.Client
	Address []byte
	Options *ContractOptions
}

type ReadOptions struct {
	BlockNumber *big.Int
}

type WriteOptions struct {
	GasConfig  *evm.GasConfig
	BlockDepth uint16 //0 means finalized, 1 confirmed, positive numbers block depth - TODO to be defined together with all other operations
}

//Logs support

const FINALIZED = 0
const CONFIRMED = 1

type LogTrackingOptions struct {
	MaxLogsKept   uint64  `protobuf:"varint,1,opt,name=max_logs_kept,json=maxLogsKept,proto3" json:"max_logs_kept,omitempty"`     // maximum number of logs to retain ( 0 = unlimited )
	RetentionTime int64   `protobuf:"varint,2,opt,name=retention_time,json=retentionTime,proto3" json:"retention_time,omitempty"` // maximum amount of time to retain logs in seconds
	LogsPerBlock  uint64  `protobuf:"varint,3,opt,name=logs_per_block,json=logsPerBlock,proto3" json:"logs_per_block,omitempty"`  // rate limit ( maximum # of logs per block, 0 = unlimited )
	Topic2        *[]byte `protobuf:"bytes,7,rep,name=topic2,proto3" json:"topic2,omitempty"`                                     // list of possible values for topic2
	Topic3        *[]byte `protobuf:"bytes,8,rep,name=topic3,proto3" json:"topic3,omitempty"`                                     // list of possible values for topic3
	Topic4        *[]byte `protobuf:"bytes,9,rep,name=topic4,proto3" json:"topic4,omitempty"`                                     // list of possible values for topic4
}

type QueryTrackedLogsOptions struct {
	SortBy []*chain_common.SortBy `protobuf:"bytes,1,rep,name=sort_by,json=sortBy,proto3" json:"sort_by,omitempty"` // A list of sorting criteria.
	Limit  *chain_common.Limit    `protobuf:"bytes,2,opt,name=limit,proto3" json:"limit,omitempty"`                 // Pagination limit and direction.
}

type ParsedLog[T any] struct {
	LogData T
	RawLog  evm.Log
}

type FilterOptions struct {
	BlockHash *[]byte
	FromBlock *big.Int
	ToBlock   *big.Int
}

func AddInterfaceMock(
	address common.Address,
	clientMock *evmmock.ClientCapability,
	callContract map[string]func(payload []byte) ([]byte, error),
	writeReport func(payload []byte, config *evmcappb.GasConfig) (*evmcappb.WriteReportReply, error),
) {

	// copy the mock so that other contract interfaces can be implemented on the same contract
	original := *clientMock

	// We need to do this for all callbacks. Some refactoring might be good...
	clientMock.CallContract = func(ctx context.Context, input *evm.CallContractRequest) (*evm.CallContractReply, error) {
		if !bytes.Equal(address[:], input.Call.To) {
			if original.CallContract == nil {
				return nil, fmt.Errorf("contract %s not found", address.Hex())
			} else {
				return original.CallContract(ctx, input)
			}
		}

		data := input.Call.Data
		if len(data) < 4 {
			return nil, errors.New("data too short")
		}

		methodID := data[:4]
		callback := callContract[string(methodID)]
		if callback == nil {
			if original.CallContract != nil {
				return original.CallContract(ctx, input)
			}
			return nil, fmt.Errorf("method with ID %x not implemented", methodID)
		}

		responsePayload, err := callback(data[4:])
		if err != nil {
			return nil, err
		}

		return &evm.CallContractReply{
			Data: responsePayload,
		}, nil
	}

	clientMock.WriteReport = func(ctx context.Context, input *evmcappb.WriteReportRequest) (*evmcappb.WriteReportReply, error) {
		if !bytes.Equal(address[:], input.Receiver) {
			if original.WriteReport == nil {
				return nil, fmt.Errorf("contract %s not found", address.Hex())
			} else {
				return original.WriteReport(ctx, input)
			}
		}

		if writeReport == nil {
			if original.WriteReport != nil {
				return original.WriteReport(ctx, input)
			}
			return nil, fmt.Errorf("writeReport callback not implemented for contract %s", address.Hex())
		}

		// TODO this needs to be right.
		raw := input.Report.RawReport
		chainSelector := binary.LittleEndian.Uint32(raw)
		if chainSelector != 1 { // TODO: Fix chain selector
			return nil, fmt.Errorf("chain selector 0x%x not implemented", chainSelector)
		}

		return writeReport(raw[4:], input.GasConfig)
	}
}
