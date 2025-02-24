package ccipaptos

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk/bcs"
)

// MarshalBCS for TokenPriceUpdate
func (tpu *TokenPriceUpdate) MarshalBCS(ser *bcs.Serializer) {
	ser.WriteBytes(tpu.SourceToken[:]) // Serialize fixed-length AccountAddress (32 bytes)
	ser.U256(*tpu.UsdPerToken)         // Serialize big.Int
}

// UnmarshalBCS for TokenPriceUpdate
func (tpu *TokenPriceUpdate) UnmarshalBCS(des *bcs.Deserializer) {
	bytes := des.ReadBytes()
	if len(bytes) != 32 {
		des.SetError(fmt.Errorf("expected 32 bytes for SourceToken, got %d", len(bytes)))
		return
	}
	copy(tpu.SourceToken[:], bytes)
	usdPerToken := des.U256()
	tpu.UsdPerToken = &usdPerToken
}

// MarshalBCS for GasPriceUpdate
func (gpu *GasPriceUpdate) MarshalBCS(ser *bcs.Serializer) {
	ser.U64(gpu.DestChainSelector) // Serialize uint64
	ser.U256(*gpu.UsdPerUnitGas)   // Serialize big.Int
}

// UnmarshalBCS for GasPriceUpdate
func (gpu *GasPriceUpdate) UnmarshalBCS(des *bcs.Deserializer) {
	gpu.DestChainSelector = des.U64()
	usdPerUnitGas := des.U256()
	gpu.UsdPerUnitGas = &usdPerUnitGas
}

// MarshalBCS for PriceUpdates
func (pu *PriceUpdates) MarshalBCS(ser *bcs.Serializer) {
	bcs.SerializeSequence(pu.TokenPriceUpdates, ser) // Serialize slice of TokenPriceUpdate
	bcs.SerializeSequence(pu.GasPriceUpdates, ser)   // Serialize slice of GasPriceUpdate
}

// UnmarshalBCS for PriceUpdates
func (pu *PriceUpdates) UnmarshalBCS(des *bcs.Deserializer) {
	pu.TokenPriceUpdates = bcs.DeserializeSequence[TokenPriceUpdate](des)
	pu.GasPriceUpdates = bcs.DeserializeSequence[GasPriceUpdate](des)
}

// MarshalBCS for MerkleRoot
func (mr *MerkleRoot) MarshalBCS(ser *bcs.Serializer) {
	ser.U64(mr.SourceChainSelector)  // Serialize uint64
	ser.U64(mr.MinSeqNr)             // Serialize uint64
	ser.U64(mr.MaxSeqNr)             // Serialize uint64
	ser.WriteBytes(mr.MerkleRoot[:]) // Serialize fixed-length byte array
}

// UnmarshalBCS for MerkleRoot
func (mr *MerkleRoot) UnmarshalBCS(des *bcs.Deserializer) {
	mr.SourceChainSelector = des.U64()
	mr.MinSeqNr = des.U64()
	mr.MaxSeqNr = des.U64()
	bytes := des.ReadBytes()
	if len(bytes) != 32 {
		des.SetError(fmt.Errorf("expected 32 bytes for MerkleRoot, got %d", len(bytes)))
		return
	}
	copy(mr.MerkleRoot[:], bytes)
}

// MarshalBCS for CommitReport
func (cr *CommitReport) MarshalBCS(ser *bcs.Serializer) {
	cr.PriceUpdates.MarshalBCS(ser)              // Serialize nested PriceUpdates
	bcs.SerializeSequence(cr.MerkleRoots, ser)   // Serialize slice of MerkleRoot
	bcs.SerializeSequence(cr.RmnSignatures, ser) // Serialize slice of variable-length byte slices
	ser.WriteBytes(cr.OfframpAddress[:])         // Serialize fixed-length AccountAddress
}

// UnmarshalBCS for CommitReport
func (cr *CommitReport) UnmarshalBCS(des *bcs.Deserializer) {
	cr.PriceUpdates.UnmarshalBCS(des)
	cr.MerkleRoots = bcs.DeserializeSequence[MerkleRoot](des)
	cr.RmnSignatures = bcs.DeserializeSequence[[]uint8](des)
	bytes := des.ReadBytes()
	if len(bytes) != 32 {
		des.SetError(fmt.Errorf("expected 32 bytes for OfframpAddress, got %d", len(bytes)))
		return
	}
	copy(cr.OfframpAddress[:], bytes)
}

// MarshalBCS for RampMessageHeader
func (rmh *RampMessageHeader) MarshalBCS(ser *bcs.Serializer) {
	ser.WriteBytes(rmh.MessageId[:]) // Serialize fixed-length byte array
	ser.U64(rmh.SourceChainSelector) // Serialize uint64
	ser.U64(rmh.DestChainSelector)   // Serialize uint64
	ser.U64(rmh.SequenceNumber)      // Serialize uint64
	ser.U64(rmh.Nonce)               // Serialize uint64
}

// UnmarshalBCS for RampMessageHeader
func (rmh *RampMessageHeader) UnmarshalBCS(des *bcs.Deserializer) {
	bytes := des.ReadBytes()
	if len(bytes) != 32 {
		des.SetError(fmt.Errorf("expected 32 bytes for MessageId, got %d", len(bytes)))
		return
	}
	copy(rmh.MessageId[:], bytes)
	rmh.SourceChainSelector = des.U64()
	rmh.DestChainSelector = des.U64()
	rmh.SequenceNumber = des.U64()
	rmh.Nonce = des.U64()
}

// MarshalBCS for Any2AptosTokenTransfer
func (att *Any2AptosTokenTransfer) MarshalBCS(ser *bcs.Serializer) {
	ser.WriteBytes(att.SourcePoolAddress)   // Serialize variable-length byte slice
	ser.WriteBytes(att.DestTokenAddress[:]) // Serialize fixed-length AccountAddress
	ser.U32(att.DestGasAmount)              // Serialize uint32
	ser.WriteBytes(att.ExtraData)           // Serialize variable-length byte slice
	ser.U256(*att.Amount)                   // Serialize big.Int
}

// UnmarshalBCS for Any2AptosTokenTransfer
func (att *Any2AptosTokenTransfer) UnmarshalBCS(des *bcs.Deserializer) {
	att.SourcePoolAddress = des.ReadBytes()
	bytes := des.ReadBytes()
	if len(bytes) != 32 {
		des.SetError(fmt.Errorf("expected 32 bytes for DestTokenAddress, got %d", len(bytes)))
		return
	}
	copy(att.DestTokenAddress[:], bytes)
	att.DestGasAmount = des.U32()
	att.ExtraData = des.ReadBytes()
	amt := des.U256()
	att.Amount = &amt
}

// MarshalBCS for Any2AptosRampMessage
func (arm *Any2AptosRampMessage) MarshalBCS(ser *bcs.Serializer) {
	arm.Header.MarshalBCS(ser)                   // Serialize nested RampMessageHeader
	ser.WriteBytes(arm.Sender)                   // Serialize variable-length byte slice
	ser.WriteBytes(arm.Data)                     // Serialize variable-length byte slice
	ser.WriteBytes(arm.Receiver[:])              // Serialize fixed-length AccountAddress
	ser.U256(*arm.GasLimit)                      // Serialize big.Int
	bcs.SerializeSequence(arm.TokenAmounts, ser) // Serialize slice of Any2AptosTokenTransfer
}

// UnmarshalBCS for Any2AptosRampMessage
func (arm *Any2AptosRampMessage) UnmarshalBCS(des *bcs.Deserializer) {
	arm.Header.UnmarshalBCS(des)
	arm.Sender = des.ReadBytes()
	arm.Data = des.ReadBytes()
	bytes := des.ReadBytes()
	if len(bytes) != 32 {
		des.SetError(fmt.Errorf("expected 32 bytes for Receiver, got %d", len(bytes)))
		return
	}
	copy(arm.Receiver[:], bytes)
	gasLimit := des.U256()
	arm.GasLimit = &gasLimit
	arm.TokenAmounts = bcs.DeserializeSequence[Any2AptosTokenTransfer](des)
}

// MarshalBCS for ExecutionReport
func (er *ExecutionReport) MarshalBCS(ser *bcs.Serializer) {
	ser.U64(er.SourceChainSelector)                  // Serialize uint64
	bcs.SerializeSequence(er.Messages, ser)          // Serialize slice of Any2AptosRampMessage
	bcs.SerializeSequence(er.OffchainTokenData, ser) // Serialize slice of slice of byte slices
	bcs.SerializeSequence(er.Proofs, ser)            // Serialize slice of [32]byte
	ser.U256(*er.ProofFlagBits)                      // Serialize big.Int
}

// UnmarshalBCS for ExecutionReport
func (er *ExecutionReport) UnmarshalBCS(des *bcs.Deserializer) {
	er.SourceChainSelector = des.U64()
	er.Messages = bcs.DeserializeSequence[Any2AptosRampMessage](des)
	er.OffchainTokenData = bcs.DeserializeSequence[[][]byte](des)
	er.Proofs = bcs.DeserializeSequence[[32]byte](des)
	proofFlagBits := des.U256()
	er.ProofFlagBits = &proofFlagBits
}
