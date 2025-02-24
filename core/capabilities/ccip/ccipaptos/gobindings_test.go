package ccipaptos

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/bcs"
)

// TestTokenPriceUpdate tests serialization and deserialization of TokenPriceUpdate
func TestTokenPriceUpdate(t *testing.T) {
	tpu := TokenPriceUpdate{
		SourceToken: aptos.AccountAddress{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
		UsdPerToken: big.NewInt(123456),
	}

	bytes, err := bcs.Serialize(&tpu)
	if err != nil {
		t.Fatalf("Failed to serialize TokenPriceUpdate: %v", err)
	}

	var deserialized TokenPriceUpdate
	err = bcs.Deserialize(&deserialized, bytes)
	if err != nil {
		t.Fatalf("Failed to deserialize TokenPriceUpdate: %v", err)
	}

	if !reflect.DeepEqual(tpu, deserialized) {
		t.Errorf("TokenPriceUpdate mismatch: expected %+v, got %+v", tpu, deserialized)
	}
}

// TestGasPriceUpdate tests serialization and deserialization of GasPriceUpdate
func TestGasPriceUpdate(t *testing.T) {
	gpu := GasPriceUpdate{
		DestChainSelector: 42,
		UsdPerUnitGas:     big.NewInt(789),
	}

	bytes, err := bcs.Serialize(&gpu)
	if err != nil {
		t.Fatalf("Failed to serialize GasPriceUpdate: %v", err)
	}

	var deserialized GasPriceUpdate
	err = bcs.Deserialize(&deserialized, bytes)
	if err != nil {
		t.Fatalf("Failed to deserialize GasPriceUpdate: %v", err)
	}

	if !reflect.DeepEqual(gpu, deserialized) {
		t.Errorf("GasPriceUpdate mismatch: expected %+v, got %+v", gpu, deserialized)
	}
}

// TestPriceUpdates tests serialization and deserialization of PriceUpdates
func TestPriceUpdates(t *testing.T) {
	pu := PriceUpdates{
		TokenPriceUpdates: []TokenPriceUpdate{
			{SourceToken: aptos.AccountAddress{1}, UsdPerToken: big.NewInt(100)},
		},
		GasPriceUpdates: []GasPriceUpdate{
			{DestChainSelector: 1, UsdPerUnitGas: big.NewInt(50)},
		},
	}

	bytes, err := bcs.Serialize(&pu)
	if err != nil {
		t.Fatalf("Failed to serialize PriceUpdates: %v", err)
	}

	var deserialized PriceUpdates
	err = bcs.Deserialize(&deserialized, bytes)
	if err != nil {
		t.Fatalf("Failed to deserialize PriceUpdates: %v", err)
	}

	if !reflect.DeepEqual(pu, deserialized) {
		t.Errorf("PriceUpdates mismatch: expected %+v, got %+v", pu, deserialized)
	}
}

// TestMerkleRoot tests serialization and deserialization of MerkleRoot
func TestMerkleRoot(t *testing.T) {
	mr := MerkleRoot{
		SourceChainSelector: 42,
		MinSeqNr:            100,
		MaxSeqNr:            200,
		MerkleRoot:          [32]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
	}

	bytes, err := bcs.Serialize(&mr)
	if err != nil {
		t.Fatalf("Failed to serialize MerkleRoot: %v", err)
	}

	var deserialized MerkleRoot
	err = bcs.Deserialize(&deserialized, bytes)
	if err != nil {
		t.Fatalf("Failed to deserialize MerkleRoot: %v", err)
	}

	if !reflect.DeepEqual(mr, deserialized) {
		t.Errorf("MerkleRoot mismatch: expected %+v, got %+v", mr, deserialized)
	}
}

// TestCommitReport tests serialization and deserialization of CommitReport
func TestCommitReport(t *testing.T) {
	cr := CommitReport{
		PriceUpdates: PriceUpdates{
			TokenPriceUpdates: []TokenPriceUpdate{
				{SourceToken: aptos.AccountAddress{1}, UsdPerToken: big.NewInt(100)},
			},
			GasPriceUpdates: []GasPriceUpdate{
				{DestChainSelector: 1, UsdPerUnitGas: big.NewInt(50)},
			},
		},
		MerkleRoots: []MerkleRoot{
			{SourceChainSelector: 42, MinSeqNr: 100, MaxSeqNr: 200, MerkleRoot: [32]uint8{1}},
		},
		RmnSignatures:  [][]uint8{{1, 2, 3}},
		OfframpAddress: aptos.AccountAddress{2},
	}

	bytes, err := bcs.Serialize(&cr)
	if err != nil {
		t.Fatalf("Failed to serialize CommitReport: %v", err)
	}

	var deserialized CommitReport
	err = bcs.Deserialize(&deserialized, bytes)
	if err != nil {
		t.Fatalf("Failed to deserialize CommitReport: %v", err)
	}

	if !reflect.DeepEqual(cr, deserialized) {
		t.Errorf("CommitReport mismatch: expected %+v, got %+v", cr, deserialized)
	}
}

// TestRampMessageHeader tests serialization and deserialization of RampMessageHeader
func TestRampMessageHeader(t *testing.T) {
	rmh := RampMessageHeader{
		MessageId:           [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
		SourceChainSelector: 42,
		DestChainSelector:   43,
		SequenceNumber:      100,
		Nonce:               200,
	}

	bytes, err := bcs.Serialize(&rmh)
	if err != nil {
		t.Fatalf("Failed to serialize RampMessageHeader: %v", err)
	}

	var deserialized RampMessageHeader
	err = bcs.Deserialize(&deserialized, bytes)
	if err != nil {
		t.Fatalf("Failed to deserialize RampMessageHeader: %v", err)
	}

	if !reflect.DeepEqual(rmh, deserialized) {
		t.Errorf("RampMessageHeader mismatch: expected %+v, got %+v", rmh, deserialized)
	}
}

// TestAny2AptosTokenTransfer tests serialization and deserialization of Any2AptosTokenTransfer
func TestAny2AptosTokenTransfer(t *testing.T) {
	att := Any2AptosTokenTransfer{
		SourcePoolAddress: []byte{1, 2, 3},
		DestTokenAddress:  aptos.AccountAddress{4},
		DestGasAmount:     500,
		ExtraData:         []byte{5, 6, 7},
		Amount:            big.NewInt(1000),
	}

	bytes, err := bcs.Serialize(&att)
	if err != nil {
		t.Fatalf("Failed to serialize Any2AptosTokenTransfer: %v", err)
	}

	var deserialized Any2AptosTokenTransfer
	err = bcs.Deserialize(&deserialized, bytes)
	if err != nil {
		t.Fatalf("Failed to deserialize Any2AptosTokenTransfer: %v", err)
	}

	if !reflect.DeepEqual(att, deserialized) {
		t.Errorf("Any2AptosTokenTransfer mismatch: expected %+v, got %+v", att, deserialized)
	}
}

// TestAny2AptosRampMessage tests serialization and deserialization of Any2AptosRampMessage
func TestAny2AptosRampMessage(t *testing.T) {
	arm := Any2AptosRampMessage{
		Header: RampMessageHeader{
			MessageId:           [32]byte{1},
			SourceChainSelector: 42,
			DestChainSelector:   43,
			SequenceNumber:      100,
			Nonce:               200,
		},
		Sender:       []byte{2, 3},
		Data:         []byte{4, 5},
		Receiver:     aptos.AccountAddress{6},
		GasLimit:     big.NewInt(10000),
		TokenAmounts: []Any2AptosTokenTransfer{{SourcePoolAddress: []byte{7}, DestTokenAddress: aptos.AccountAddress{8}, DestGasAmount: 500, ExtraData: []byte{9}, Amount: big.NewInt(2000)}},
	}

	bytes, err := bcs.Serialize(&arm)
	if err != nil {
		t.Fatalf("Failed to serialize Any2AptosRampMessage: %v", err)
	}

	var deserialized Any2AptosRampMessage
	err = bcs.Deserialize(&deserialized, bytes)
	if err != nil {
		t.Fatalf("Failed to deserialize Any2AptosRampMessage: %v", err)
	}

	if !reflect.DeepEqual(arm, deserialized) {
		t.Errorf("Any2AptosRampMessage mismatch: expected %+v, got %+v", arm, deserialized)
	}
}

// TestExecutionReport tests serialization and deserialization of ExecutionReport
func TestExecutionReport(t *testing.T) {
	er := ExecutionReport{
		SourceChainSelector: 42,
		Messages: []Any2AptosRampMessage{
			{
				Header:       RampMessageHeader{MessageId: [32]byte{1}, SourceChainSelector: 42, DestChainSelector: 43, SequenceNumber: 100, Nonce: 200},
				Sender:       []byte{2},
				Data:         []byte{3},
				Receiver:     aptos.AccountAddress{4},
				GasLimit:     big.NewInt(10000),
				TokenAmounts: []Any2AptosTokenTransfer{{SourcePoolAddress: []byte{5}, DestTokenAddress: aptos.AccountAddress{6}, DestGasAmount: 500, ExtraData: []byte{7}, Amount: big.NewInt(2000)}},
			},
		},
		OffchainTokenData: [][][]byte{{{1, 2}, {3, 4}}},
		Proofs:            [][32]byte{{8}},
		ProofFlagBits:     big.NewInt(1),
	}

	bytes, err := bcs.Serialize(&er)
	if err != nil {
		t.Fatalf("Failed to serialize ExecutionReport: %v", err)
	}

	var deserialized ExecutionReport
	err = bcs.Deserialize(&deserialized, bytes)
	if err != nil {
		t.Fatalf("Failed to deserialize ExecutionReport: %v", err)
	}

	if !reflect.DeepEqual(er, deserialized) {
		t.Errorf("ExecutionReport mismatch: expected %+v, got %+v", er, deserialized)
	}
}
