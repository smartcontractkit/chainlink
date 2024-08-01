package ocr3types

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// ContractTransmitter sends new reports to a smart contract or other system.
//
// All its functions should be thread-safe.
type ContractTransmitter[RI any] interface {

	// Transmit sends the report to the on-chain smart contract's Transmit
	// method.
	//
	// In most cases, implementations of this function should store the
	// transmission in a queue/database/..., but perform the actual
	// transmission (and potentially confirmation) of the transaction
	// asynchronously.
	Transmit(
		context.Context,
		types.ConfigDigest,
		uint64,
		ReportWithInfo[RI],
		[]types.AttributedOnchainSignature,
	) error

	// Account from which the transmitter invokes the contract
	FromAccount() (types.Account, error)
}

// OnchainKeyring provides cryptographic signatures that need to be verifiable
// on the targeted blockchain. The underlying cryptographic primitives may be
// different on each chain; for example, on Ethereum one would use ECDSA over
// secp256k1 and Keccak256, whereas on Solana one would use Ed25519 and SHA256.
//
// All its functions should be thread-safe.
type OnchainKeyring[RI any] interface {
	// PublicKey returns the public key of the keypair used by Sign.
	PublicKey() types.OnchainPublicKey

	// Sign returns a signature over Report.
	Sign(types.ConfigDigest, uint64, ReportWithInfo[RI]) (signature []byte, err error)

	// Verify verifies a signature over ReportContext and Report allegedly
	// created from OnchainPublicKey.
	//
	// Implementations of this function must gracefully handle malformed or
	// adversarially crafted inputs.
	Verify(_ types.OnchainPublicKey, _ types.ConfigDigest, seqNr uint64, _ ReportWithInfo[RI], signature []byte) bool

	// Maximum length of a signature
	MaxSignatureLength() int
}
