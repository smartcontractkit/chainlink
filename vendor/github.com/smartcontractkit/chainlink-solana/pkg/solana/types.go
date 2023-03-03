package solana

import (
	"errors"
	"math/big"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
)

const (
	AccountDiscriminatorLen uint64 = 8

	// TransmissionLen = Slot, Timestamp, Padding0, Answer, Padding1, Padding2
	TransmissionLen uint64 = 8 + 4 + 4 + 16 + 8 + 8

	// TransmissionsHeaderLen = Version, State, Owner, ProposedOwner, Writer, Description, Decimals, FlaggingThreshold, LatestRoundID, Granularity, LiveLength, LiveCursor, HistoricalCursor
	TransmissionsHeaderLen     uint64 = 1 + 1 + 32 + 32 + 32 + 32 + 1 + 4 + 4 + 1 + 4 + 4 + 4
	TransmissionsHeaderMaxSize uint64 = 192 // max area allocated to transmissions header

	// ReportLen data (61 bytes)
	MedianLen       uint64 = 16
	JuelsLen        uint64 = 8
	ReportHeaderLen uint64 = 4 + 1 + 32 // timestamp (uint32) + number of observers (uint8) + observer array [32]uint8
	ReportLen       uint64 = ReportHeaderLen + MedianLen + JuelsLen

	// MaxOracles is the maximum number of oracles that can be stored onchain
	MaxOracles = 19
	// MaxOffchainConfigLen is the maximum byte length for the encoded offchainconfig
	MaxOffchainConfigLen = 4096
)

// State is the struct representing the contract state
type State struct {
	AccountDiscriminator [8]byte // first 8 bytes of the SHA256 of the account’s Rust ident, https://docs.rs/anchor-lang/0.18.2/anchor_lang/attr.account.html
	Version              uint8
	Nonce                uint8
	Padding0             uint16
	Padding1             uint32
	Transmissions        solana.PublicKey
	Config               Config
	OffchainConfig       OffchainConfig
	Oracles              Oracles
}

// SigningKey represents the report signing key
type SigningKey struct {
	Key [20]byte
}

type OffchainConfig struct {
	Version uint64
	Raw     [MaxOffchainConfigLen]byte
	Len     uint64
}

func (oc OffchainConfig) Data() ([]byte, error) {
	if oc.Len > MaxOffchainConfigLen {
		return []byte{}, errors.New("OffchainConfig.Len exceeds MaxOffchainConfigLen")
	}
	return oc.Raw[:oc.Len], nil
}

// Config contains the configuration of the contract
type Config struct {
	Owner                     solana.PublicKey
	ProposedOwner             solana.PublicKey
	TokenMint                 solana.PublicKey
	TokenVault                solana.PublicKey
	RequesterAccessController solana.PublicKey
	BillingAccessController   solana.PublicKey
	MinAnswer                 bin.Int128
	MaxAnswer                 bin.Int128
	F                         uint8
	Round                     uint8
	Padding0                  uint16
	Epoch                     uint32
	LatestAggregatorRoundID   uint32
	LatestTransmitter         solana.PublicKey
	ConfigCount               uint32
	LatestConfigDigest        [32]byte
	LatestConfigBlockNumber   uint64
	Billing                   Billing
}

// Oracles contains the list of oracles
type Oracles struct {
	Raw [MaxOracles]Oracle
	Len uint64
}

func (o Oracles) Data() ([]Oracle, error) {
	if o.Len > MaxOracles {
		return []Oracle{}, errors.New("Oracles.Len exceeds MaxOracles")
	}
	return o.Raw[:o.Len], nil
}

// Oracle contains information about the reporting nodes
type Oracle struct {
	Transmitter   solana.PublicKey
	Signer        SigningKey
	Payee         solana.PublicKey
	ProposedPayee solana.PublicKey
	FromRoundID   uint32
	Payment       uint64
}

// Billing contains the payment information
type Billing struct {
	ObservationPayment  uint32
	TransmissionPayment uint32
}

// Answer contains the current price answer
type Answer struct {
	Data      *big.Int
	Timestamp uint32
}

// Access controller state
type AccessController struct {
	Owner         solana.PublicKey
	ProposedOwner solana.PublicKey
	Access        [32]solana.PublicKey
	Len           uint64
}

// TransmissionsHeader struct for decoding transmission state header
type TransmissionsHeader struct {
	Version           uint8
	State             uint8
	Owner             solana.PublicKey
	ProposedOwner     solana.PublicKey
	Writer            solana.PublicKey
	Description       [32]byte
	Decimals          uint8
	FlaggingThreshold uint32
	LatestRoundID     uint32
	Granularity       uint8
	LiveLength        uint32
	LiveCursor        uint32
	HistoricalCursor  uint32
}

// Transmission struct for decoding individual tranmissions
type Transmission struct {
	Slot      uint64
	Timestamp uint32
	Padding0  uint32
	Answer    bin.Int128
	Padding1  uint64
	Padding2  uint64
}

// TransmissionV1 struct for parsing results pre-migration
type TransmissionV1 struct {
	Timestamp uint64
	Answer    bin.Int128
}

// CL Core OCR2 job spec RelayConfig member for Solana
type RelayConfig struct {
	// network data
	ChainID string `json:"chainID"` // required

	// state account passed as the ContractID in main job spec
	// on-chain program + transmissions account + store programID
	OCR2ProgramID   string `json:"ocr2ProgramID"`
	TransmissionsID string `json:"transmissionsID"`
	StoreProgramID  string `json:"storeProgramID"`
}
