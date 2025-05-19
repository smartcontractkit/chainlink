// Code generated evm-bindings; DO NOT EDIT.

package bindings

import (
	"context"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"math/big"
)

// CodeDetails methods inputs and outputs structs
type ChainReaderTester struct {
	BoundContract  types.BoundContract
	ContractReader types.ContractReader
	ChainWriter    types.ContractWriter
}

type AccountStruct struct {
	Account    []byte
	AccountStr []byte
}

type AddTestStructInput struct {
	Field               int32
	DifferentField      string
	OracleId            uint8
	OracleIds           [32]uint8
	AccountStruct       AccountStruct
	Accounts            [][]byte
	BigField            *big.Int
	NestedDynamicStruct MidLevelDynamicTestStruct
	NestedStaticStruct  MidLevelStaticTestStruct
}

type GetAlterablePrimitiveValueOutput struct {
	Value uint64
}

type GetDifferentPrimitiveValueOutput struct {
	Value uint64
}

type GetElementAtIndexInput struct {
	I *big.Int
}

type GetPrimitiveValueOutput struct {
	Value uint64
}

type GetSliceValueOutput struct {
	Value []uint64
}

type InnerDynamicTestStruct struct {
	IntVal int64
	S      string
}

type InnerStaticTestStruct struct {
	IntVal int64
	A      []byte
}

type MidLevelDynamicTestStruct struct {
	FixedBytes [2]uint8
	Inner      InnerDynamicTestStruct
}

type MidLevelStaticTestStruct struct {
	FixedBytes [2]uint8
	Inner      InnerStaticTestStruct
}

type ReturnSeenInput struct {
	Field               int32
	DifferentField      string
	OracleId            uint8
	OracleIds           [32]uint8
	AccountStruct       AccountStruct
	Accounts            [][]byte
	BigField            *big.Int
	NestedDynamicStruct MidLevelDynamicTestStruct
	NestedStaticStruct  MidLevelStaticTestStruct
}

type SetAlterablePrimitiveValueInput struct {
	Value uint64
}

type TestStruct struct {
	Field               int32
	DifferentField      string
	OracleId            uint8
	OracleIds           [32]uint8
	AccountStruct       AccountStruct
	Accounts            [][]byte
	BigField            *big.Int
	NestedDynamicStruct MidLevelDynamicTestStruct
	NestedStaticStruct  MidLevelStaticTestStruct
}

type TriggerEventInput struct {
	Field               int32
	OracleId            uint8
	NestedDynamicStruct MidLevelDynamicTestStruct
	NestedStaticStruct  MidLevelStaticTestStruct
	OracleIds           [32]uint8
	AccountStruct       AccountStruct
	Accounts            [][]byte
	DifferentField      string
	BigField            *big.Int
}

type TriggerEventWithDynamicTopicInput struct {
	Field string
}

type TriggerStaticBytesInput struct {
	Val1 uint32
	Val2 uint32
	Val3 uint32
	Val4 uint64
	Val5 [32]uint8
	Val6 [32]uint8
	Val7 [32]uint8
	Raw  []uint8
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

func (b ChainReaderTester) GetPrimitiveValue(ctx context.Context, confidence primitives.ConfidenceLevel) (uint64, error) {
	var output uint64
	err := b.ContractReader.GetLatestValue(ctx, b.BoundContract.ReadIdentifier("GetPrimitiveValue"), confidence, nil, &output)
	return output, err
}

func (b ChainReaderTester) GetSliceValue(ctx context.Context, confidence primitives.ConfidenceLevel) ([]uint64, error) {
	var output []uint64
	err := b.ContractReader.GetLatestValue(ctx, b.BoundContract.ReadIdentifier("GetSliceValue"), confidence, nil, &output)
	return output, err
}

func (b ChainReaderTester) ReturnSeen(ctx context.Context, input ReturnSeenInput, confidence primitives.ConfidenceLevel) (TestStruct, error) {
	output := TestStruct{}
	err := b.ContractReader.GetLatestValue(ctx, b.BoundContract.ReadIdentifier("ReturnSeen"), confidence, input, &output)
	return output, err
}

func (b ChainReaderTester) TriggerEvent(ctx context.Context, input TriggerEventInput, txId string, toAddress string, meta *types.TxMeta) error {
	return b.ChainWriter.SubmitTransaction(ctx, "ChainReaderTester", "TriggerEvent", input, txId, toAddress, meta, big.NewInt(0))
}

func (b ChainReaderTester) GetDifferentPrimitiveValue(ctx context.Context, confidence primitives.ConfidenceLevel) (uint64, error) {
	var output uint64
	err := b.ContractReader.GetLatestValue(ctx, b.BoundContract.ReadIdentifier("GetDifferentPrimitiveValue"), confidence, nil, &output)
	return output, err
}

func (b ChainReaderTester) GetElementAtIndex(ctx context.Context, input GetElementAtIndexInput, confidence primitives.ConfidenceLevel) (TestStruct, error) {
	output := TestStruct{}
	err := b.ContractReader.GetLatestValue(ctx, b.BoundContract.ReadIdentifier("GetElementAtIndex"), confidence, input, &output)
	return output, err
}

func (b ChainReaderTester) SetAlterablePrimitiveValue(ctx context.Context, input SetAlterablePrimitiveValueInput, txId string, toAddress string, meta *types.TxMeta) error {
	return b.ChainWriter.SubmitTransaction(ctx, "ChainReaderTester", "SetAlterablePrimitiveValue", input, txId, toAddress, meta, big.NewInt(0))
}

func (b ChainReaderTester) TriggerEventWithDynamicTopic(ctx context.Context, input TriggerEventWithDynamicTopicInput, txId string, toAddress string, meta *types.TxMeta) error {
	return b.ChainWriter.SubmitTransaction(ctx, "ChainReaderTester", "TriggerEventWithDynamicTopic", input, txId, toAddress, meta, big.NewInt(0))
}

func (b ChainReaderTester) TriggerWithFourTopics(ctx context.Context, input TriggerWithFourTopicsInput, txId string, toAddress string, meta *types.TxMeta) error {
	return b.ChainWriter.SubmitTransaction(ctx, "ChainReaderTester", "TriggerWithFourTopics", input, txId, toAddress, meta, big.NewInt(0))
}

func (b ChainReaderTester) TriggerWithFourTopicsWithHashed(ctx context.Context, input TriggerWithFourTopicsWithHashedInput, txId string, toAddress string, meta *types.TxMeta) error {
	return b.ChainWriter.SubmitTransaction(ctx, "ChainReaderTester", "TriggerWithFourTopicsWithHashed", input, txId, toAddress, meta, big.NewInt(0))
}

func (b ChainReaderTester) AddTestStruct(ctx context.Context, input AddTestStructInput, txId string, toAddress string, meta *types.TxMeta) error {
	return b.ChainWriter.SubmitTransaction(ctx, "ChainReaderTester", "AddTestStruct", input, txId, toAddress, meta, big.NewInt(0))
}

func (b ChainReaderTester) GetAlterablePrimitiveValue(ctx context.Context, confidence primitives.ConfidenceLevel) (uint64, error) {
	var output uint64
	err := b.ContractReader.GetLatestValue(ctx, b.BoundContract.ReadIdentifier("GetAlterablePrimitiveValue"), confidence, nil, &output)
	return output, err
}
