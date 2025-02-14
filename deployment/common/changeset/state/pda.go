package state

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
)

func getPDA(programID solana.PublicKey, seeds [][]byte) solana.PublicKey {
	pda, _, _ := solana.FindProgramAddress(seeds, programID)
	return pda
}

func GetMCMSignerPDA(programID solana.PublicKey, msigID PDASeed) solana.PublicKey {
	seeds := [][]byte{[]byte("multisig_signer"), msigID[:]}
	return getPDA(programID, seeds)
}

func GetMCMConfigPDA(programID solana.PublicKey, msigID PDASeed) solana.PublicKey {
	seeds := [][]byte{[]byte("multisig_config"), msigID[:]}
	return getPDA(programID, seeds)
}

func GetMCMRootMetadataPDA(programID solana.PublicKey, msigID PDASeed) solana.PublicKey {
	seeds := [][]byte{[]byte("root_metadata"), msigID[:]}
	return getPDA(programID, seeds)
}

func GetMCMExpiringRootAndOpCountPDA(programID solana.PublicKey, pdaSeed PDASeed) solana.PublicKey {
	seeds := [][]byte{[]byte("expiring_root_and_op_count"), pdaSeed[:]}
	return getPDA(programID, seeds)
}

func GetMCMRootSignaturesPDA(
	programID solana.PublicKey, msigID PDASeed, root common.Hash, validUntil uint32,
) solana.PublicKey {
	seeds := [][]byte{[]byte("root_signatures"), msigID[:], root[:], validUntilBytes(validUntil)}
	return getPDA(programID, seeds)
}

func GetMCMSeenSignedHashesPDA(
	programID solana.PublicKey, msigID PDASeed, root common.Hash, validUntil uint32,
) solana.PublicKey {
	seeds := [][]byte{[]byte("seen_signed_hashes"), msigID[:], root[:], validUntilBytes(validUntil)}
	return getPDA(programID, seeds)
}

func GetTimelockConfigPDA(programID solana.PublicKey, timelockID PDASeed) solana.PublicKey {
	seeds := [][]byte{[]byte("timelock_config"), timelockID[:]}
	return getPDA(programID, seeds)
}

func GetTimelockOperationPDA(programID solana.PublicKey, timelockID PDASeed, opID [32]byte) solana.PublicKey {
	seeds := [][]byte{[]byte("timelock_operation"), timelockID[:], opID[:]}
	return getPDA(programID, seeds)
}

func GetTimelockBypasserOperationPDA(programID solana.PublicKey, timelockID PDASeed, opID [32]byte) solana.PublicKey {
	seeds := [][]byte{[]byte("timelock_bypasser_operation"), timelockID[:], opID[:]}
	return getPDA(programID, seeds)
}

func GetTimelockSignerPDA(programID solana.PublicKey, timelockID PDASeed) solana.PublicKey {
	seeds := [][]byte{[]byte("timelock_signer"), timelockID[:]}
	return getPDA(programID, seeds)
}

func validUntilBytes(validUntil uint32) []byte {
	const uint32Size = 4
	vuBytes := make([]byte, uint32Size)
	binary.LittleEndian.PutUint32(vuBytes, validUntil)

	return vuBytes
}
