package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/NethermindEth/juno/core/felt"
)

var ErrInvalidBlockID = errors.New("invalid blockid")

// BlockHashAndNumberOutput is a struct that is returned by BlockHashAndNumber.
type BlockHashAndNumberOutput struct {
	BlockNumber uint64     `json:"block_number,omitempty"`
	BlockHash   *felt.Felt `json:"block_hash,omitempty"`
}

// BlockID is a struct that is used to choose between different
// search types.
type BlockID struct {
	Number *uint64    `json:"block_number,omitempty"`
	Hash   *felt.Felt `json:"block_hash,omitempty"`
	Tag    string     `json:"block_tag,omitempty"`
}

// MarshalJSON marshals the BlockID to JSON format.
//
// It returns a byte slice and an error. The byte slice contains the JSON representation of the BlockID,
// while the error indicates any error that occurred during the marshaling process.
//
// Parameters:
//
//	none
//
// Returns:
// - []byte: the JSON representation of the BlockID
// - error: any error that occurred during the marshaling process
func (b BlockID) MarshalJSON() ([]byte, error) {
	if b.Tag == "pending" || b.Tag == "latest" {
		return []byte(strconv.Quote(b.Tag)), nil
	}

	if b.Tag != "" {
		return nil, ErrInvalidBlockID
	}

	if b.Number != nil {
		return []byte(fmt.Sprintf(`{"block_number":%d}`, *b.Number)), nil
	}

	if b.Hash.BigInt(big.NewInt(0)).BitLen() != 0 {
		return []byte(fmt.Sprintf(`{"block_hash":"%s"}`, b.Hash.String())), nil
	}

	return nil, ErrInvalidBlockID

}

type BlockStatus string

const (
	BlockStatus_Pending      BlockStatus = "PENDING"
	BlockStatus_AcceptedOnL2 BlockStatus = "ACCEPTED_ON_L2"
	BlockStatus_AcceptedOnL1 BlockStatus = "ACCEPTED_ON_L1"
	BlockStatus_Rejected     BlockStatus = "REJECTED"
)

// UnmarshalJSON unmarshals the JSON representation of a BlockStatus.
//
// It takes in a byte slice containing the JSON data to be unmarshaled.
// The function returns an error if there is an issue unmarshaling the data.
//
// Parameters:
// - data: It takes a byte slice as a parameter, which represents the JSON data to be unmarshaled
// Returns:
// - error: an error if the unmarshaling fails
func (bs *BlockStatus) UnmarshalJSON(data []byte) error {
	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	switch unquoted {
	case "PENDING":
		*bs = BlockStatus_Pending
	case "ACCEPTED_ON_L2":
		*bs = BlockStatus_AcceptedOnL2
	case "ACCEPTED_ON_L1":
		*bs = BlockStatus_AcceptedOnL1
	case "REJECTED":
		*bs = BlockStatus_Rejected
	default:
		return fmt.Errorf("unsupported status: %s", data)
	}

	return nil
}

// MarshalJSON returns the JSON encoding of BlockStatus.
//
// Parameters:
//
//	none
//
// Returns:
// - []byte: a byte slice
// - error: an error if any
func (bs BlockStatus) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(string(bs))), nil
}

type Block struct {
	BlockHeader
	Status BlockStatus `json:"status"`
	// Transactions The transactions in this block
	Transactions BlockTransactions `json:"transactions"`
}

type PendingBlock struct {
	PendingBlockHeader
	BlockTransactions
}

// encoding/json doesn't support inlining fields
type BlockWithReceipts struct {
	BlockStatus BlockStatus `json:"status"`
	BlockHeader
	BlockBodyWithReceipts
}

type BlockBodyWithReceipts struct {
	Transactions []TransactionWithReceipt `json:"transactions"`
}

type TransactionWithReceipt struct {
	Transaction UnknownTransaction        `json:"transaction"`
	Receipt     UnknownTransactionReceipt `json:"receipt"`
}

// The dynamic block being constructed by the sequencer. Note that this object will be deprecated upon decentralization.
type PendingBlockWithReceipts struct {
	PendingBlockHeader
	BlockBodyWithReceipts
}

type BlockTxHashes struct {
	BlockHeader
	Status BlockStatus `json:"status"`
	// Transactions The hashes of the transactions included in this block
	Transactions []*felt.Felt `json:"transactions"`
}

type PendingBlockTxHashes struct {
	PendingBlockHeader
	Transactions []*felt.Felt `json:"transactions"`
}

type BlockHeader struct {
	// BlockHash The hash of this block
	BlockHash *felt.Felt `json:"block_hash"`
	// ParentHash The hash of this block's parent
	ParentHash *felt.Felt `json:"parent_hash"`
	// BlockNumber the block number (its height)
	BlockNumber uint64 `json:"block_number"`
	// NewRoot The new global state root
	NewRoot *felt.Felt `json:"new_root"`
	// Timestamp the time in which the block was created, encoded in Unix time
	Timestamp uint64 `json:"timestamp"`
	// SequencerAddress the StarkNet identity of the sequencer submitting this block
	SequencerAddress *felt.Felt `json:"sequencer_address"`
	// The price of l1 gas in the block
	L1GasPrice ResourcePrice `json:"l1_gas_price"`
	// The price of l1 data gas in the block
	L1DataGasPrice ResourcePrice `json:"l1_data_gas_price"`
	// Specifies whether the data of this block is published via blob data or calldata
	L1DAMode L1DAMode `json:"l1_da_mode"`
	// Semver of the current Starknet protocol
	StarknetVersion string `json:"starknet_version"`
}

type L1DAMode int

const (
	L1DAModeBlob L1DAMode = iota
	L1DAModeCalldata
)

func (mode L1DAMode) String() string {
	switch mode {
	case L1DAModeBlob:
		return "BLOB"
	case L1DAModeCalldata:
		return "CALLDATA"
	default:
		return "Unknown L1DAMode"
	}
}

func (mode *L1DAMode) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), "\"")
	switch str {
	case "BLOB":
		*mode = L1DAModeBlob
	case "CALLDATA":
		*mode = L1DAModeCalldata
	default:
		return fmt.Errorf("unknown L1DAMode: %s", str)
	}
	return nil
}
func (mode L1DAMode) MarshalJSON() ([]byte, error) {
	return json.Marshal(mode.String())
}

type PendingBlockHeader struct {
	// ParentHash The hash of this block's parent
	ParentHash *felt.Felt `json:"parent_hash"`
	// Timestamp the time in which the block was created, encoded in Unix time
	Timestamp uint64 `json:"timestamp"`
	// SequencerAddress the StarkNet identity of the sequencer submitting this block
	SequencerAddress *felt.Felt `json:"sequencer_address"`
	// The price of l1 gas in the block
	L1GasPrice ResourcePrice `json:"l1_gas_price"`
	// Semver of the current Starknet protocol
	StarknetVersion string `json:"starknet_version"`
	// The price of l1 data gas in the block
	L1DataGasPrice ResourcePrice `json:"l1_data_gas_price"`
	// Specifies whether the data of this block is published via blob data or calldata
	L1DAMode L1DAMode `json:"l1_da_mode"`
}

type ResourcePrice struct {
	// the price of one unit of the given resource, denominated in fri (10^-18 strk)
	PriceInFRI *felt.Felt `json:"price_in_strk,omitempty"`
	// The price of one unit of the given resource, denominated in wei
	PriceInWei *felt.Felt `json:"price_in_wei"`
}
