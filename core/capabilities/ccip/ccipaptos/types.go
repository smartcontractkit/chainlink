package ccipaptos

import (
	"math/big"

	"github.com/aptos-labs/aptos-go-sdk"
)

// TokenPriceUpdate struct
type TokenPriceUpdate struct {
	SourceToken aptos.AccountAddress
	UsdPerToken *big.Int
}

// GasPriceUpdate struct
type GasPriceUpdate struct {
	DestChainSelector uint64
	UsdPerUnitGas     *big.Int
}

// PriceUpdates struct
type PriceUpdates struct {
	TokenPriceUpdates []TokenPriceUpdate
	GasPriceUpdates   []GasPriceUpdate
}

// MerkleRoot struct
type MerkleRoot struct {
	SourceChainSelector uint64
	// OnRampAddress       []byte <-- Not there onchain. Investigate?
	MinSeqNr   uint64
	MaxSeqNr   uint64
	MerkleRoot [32]uint8
}

// CommitInput struct with optional MerkleRoot
type CommitReport struct {
	PriceUpdates   PriceUpdates
	MerkleRoots    []MerkleRoot
	RmnSignatures  [][]uint8
	OfframpAddress aptos.AccountAddress // This is only for Aptos
}

type RampMessageHeader struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	DestChainSelector   uint64
	SequenceNumber      uint64
	Nonce               uint64
}

type Any2AptosTokenTransfer struct {
	SourcePoolAddress []byte
	DestTokenAddress  aptos.AccountAddress
	DestGasAmount     uint32
	ExtraData         []byte
	Amount            *big.Int
}

type Any2AptosRampMessage struct {
	Header       RampMessageHeader
	Sender       []byte
	Data         []byte
	Receiver     aptos.AccountAddress
	GasLimit     *big.Int
	TokenAmounts []Any2AptosTokenTransfer
}

type ExecutionReport struct {
	SourceChainSelector uint64
	Messages            []Any2AptosRampMessage
	OffchainTokenData   [][][]byte
	Proofs              [][32]byte
	ProofFlagBits       *big.Int
}
