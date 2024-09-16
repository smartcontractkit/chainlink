// Code generated evm-bindings; DO NOT EDIT.

package bindings

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"math/big"
)

// CodeDetails methods inputs and outputs structs
type ChainReaderTester struct {
	ContractReader types.ContractReader
	ChainWriter    types.ChainWriter
}

type AddTestStructInput struct {
	Field          int32
	DifferentField string
	OracleId       uint8
	OracleIds      [32]uint8
	Account        common.Address
	Accounts       []common.Address
	BigField       big.Int
	NestedStruct   MidLevelTestStruct
}

type GetAlterablePrimitiveValueOutput struct {
	Value uint64
}

type GetDifferentPrimitiveValueOutput struct {
	Value uint64
}

type GetElementAtIndexInput struct {
	I uint64
}

type GetPrimitiveValueOutput struct {
	Value uint64
}

type GetSliceValueOutput struct {
	Value string
}

type InnerTestStruct struct {
	IntVal int64
	S      string
}

type MidLevelTestStruct struct {
	FixedBytes [2]uint8
	Inner      InnerTestStruct
}

type ReturnSeenInput struct {
	Field          int32
	DifferentField string
	OracleId       uint8
	OracleIds      [32]uint8
	Account        common.Address
	Accounts       []common.Address
	BigField       big.Int
	NestedStruct   MidLevelTestStruct
}

type SetAlterablePrimitiveValueInput struct {
	Value uint64
}

type TestStruct struct {
	Field          int32
	DifferentField string
	OracleId       uint8
	OracleIds      [32]uint8
	Account        common.Address
	Accounts       []common.Address
	BigField       big.Int
	NestedStruct   MidLevelTestStruct
}

type TriggerEventInput struct {
	Field          int32
	DifferentField string
	OracleId       uint8
	OracleIds      [32]uint8
	Account        common.Address
	Accounts       []common.Address
	BigField       big.Int
	NestedStruct   MidLevelTestStruct
}

type TriggerEventWithDynamicTopicInput struct {
	Field string
}

type TriggerWithFourTopicsInput struct {
	Field1 int32
	Field2 int32
	Field3 int32
}

type TriggerWithFourTopicsWithHashedInput struct {
	Field1 string
	Field2 [32]uint8
	Field3 [32]uint8
}

func (b ChainReaderTester) AddTestStruct(ctx context.Context, input AddTestStructInput, txId string, toAddress string, meta *types.TxMeta) error {
	return b.ChainWriter.SubmitTransaction(ctx, "ChainReaderTester", "AddTestStruct", input, txId, toAddress, meta, big.NewInt(0))
}

func (b ChainReaderTester) GetSliceValue(ctx context.Context, confidence primitives.ConfidenceLevel) (string, error) {
	var output string
	err := b.ContractReader.GetLatestValue(ctx, "ChainReaderTester", "GetSliceValue", confidence, nil, &output)
	return output, err
}

func (b ChainReaderTester) TriggerWithFourTopics(ctx context.Context, input TriggerWithFourTopicsInput, txId string, toAddress string, meta *types.TxMeta) error {
	return b.ChainWriter.SubmitTransaction(ctx, "ChainReaderTester", "TriggerWithFourTopics", input, txId, toAddress, meta, big.NewInt(0))
}

func (b ChainReaderTester) ReturnSeen(ctx context.Context, input ReturnSeenInput, confidence primitives.ConfidenceLevel) (TestStruct, error) {
	output := TestStruct{}
	err := b.ContractReader.GetLatestValue(ctx, "ChainReaderTester", "ReturnSeen", confidence, input, &output)
	return output, err
}

func (b ChainReaderTester) SetAlterablePrimitiveValue(ctx context.Context, input SetAlterablePrimitiveValueInput, txId string, toAddress string, meta *types.TxMeta) error {
	return b.ChainWriter.SubmitTransaction(ctx, "ChainReaderTester", "SetAlterablePrimitiveValue", input, txId, toAddress, meta, big.NewInt(0))
}

func (b ChainReaderTester) TriggerEvent(ctx context.Context, input TriggerEventInput, txId string, toAddress string, meta *types.TxMeta) error {
	return b.ChainWriter.SubmitTransaction(ctx, "ChainReaderTester", "TriggerEvent", input, txId, toAddress, meta, big.NewInt(0))
}

func (b ChainReaderTester) TriggerEventWithDynamicTopic(ctx context.Context, input TriggerEventWithDynamicTopicInput, txId string, toAddress string, meta *types.TxMeta) error {
	return b.ChainWriter.SubmitTransaction(ctx, "ChainReaderTester", "TriggerEventWithDynamicTopic", input, txId, toAddress, meta, big.NewInt(0))
}

func (b ChainReaderTester) GetAlterablePrimitiveValue(ctx context.Context, confidence primitives.ConfidenceLevel) (uint64, error) {
	var output uint64
	err := b.ContractReader.GetLatestValue(ctx, "ChainReaderTester", "GetAlterablePrimitiveValue", confidence, nil, &output)
	return output, err
}

func (b ChainReaderTester) GetDifferentPrimitiveValue(ctx context.Context, confidence primitives.ConfidenceLevel) (uint64, error) {
	var output uint64
	err := b.ContractReader.GetLatestValue(ctx, "ChainReaderTester", "GetDifferentPrimitiveValue", confidence, nil, &output)
	return output, err
}

func (b ChainReaderTester) GetElementAtIndex(ctx context.Context, input GetElementAtIndexInput, confidence primitives.ConfidenceLevel) (TestStruct, error) {
	output := TestStruct{}
	err := b.ContractReader.GetLatestValue(ctx, "ChainReaderTester", "GetElementAtIndex", confidence, input, &output)
	return output, err
}

func (b ChainReaderTester) GetPrimitiveValue(ctx context.Context, confidence primitives.ConfidenceLevel) (uint64, error) {
	var output uint64
	err := b.ContractReader.GetLatestValue(ctx, "ChainReaderTester", "GetPrimitiveValue", confidence, nil, &output)
	return output, err
}

func (b ChainReaderTester) TriggerWithFourTopicsWithHashed(ctx context.Context, input TriggerWithFourTopicsWithHashedInput, txId string, toAddress string, meta *types.TxMeta) error {
	return b.ChainWriter.SubmitTransaction(ctx, "ChainReaderTester", "TriggerWithFourTopicsWithHashed", input, txId, toAddress, meta, big.NewInt(0))
}
