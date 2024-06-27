package types

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	pkgerrors "github.com/pkg/errors"
	"github.com/ugorji/go/codec"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/hex"

	htrktypes "github.com/smartcontractkit/chainlink/v2/common/headtracker/types"
	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types/internal/blocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

// Head represents a BlockNumber, BlockHash.
type Head struct {
	ID               uint64
	Hash             common.Hash
	Number           int64
	L1BlockNumber    sql.NullInt64
	ParentHash       common.Hash
	Parent           *Head
	EVMChainID       *ubig.Big
	Timestamp        time.Time
	CreatedAt        time.Time
	BaseFeePerGas    *assets.Wei
	ReceiptsRoot     common.Hash
	TransactionsRoot common.Hash
	StateRoot        common.Hash
	Difficulty       *big.Int
	TotalDifficulty  *big.Int
	IsFinalized      bool
}

var _ commontypes.Head[common.Hash] = &Head{}
var _ htrktypes.Head[common.Hash, *big.Int] = &Head{}

// NewHead returns a Head instance.
func NewHead(number *big.Int, blockHash common.Hash, parentHash common.Hash, timestamp uint64, chainID *ubig.Big) Head {
	return Head{
		Number:     number.Int64(),
		Hash:       blockHash,
		ParentHash: parentHash,
		Timestamp:  time.Unix(int64(timestamp), 0),
		EVMChainID: chainID,
	}
}

func (h *Head) BlockNumber() int64 {
	return h.Number
}

func (h *Head) BlockHash() common.Hash {
	return h.Hash
}

func (h *Head) GetParentHash() common.Hash {
	return h.ParentHash
}

func (h *Head) GetParent() commontypes.Head[common.Hash] {
	if h.Parent == nil {
		return nil
	}
	return h.Parent
}

func (h *Head) GetTimestamp() time.Time {
	return h.Timestamp
}

func (h *Head) BlockDifficulty() *big.Int {
	return h.Difficulty
}

// EarliestInChain recurses through parents until it finds the earliest one
func (h *Head) EarliestInChain() *Head {
	for h.Parent != nil {
		h = h.Parent
	}
	return h
}

// EarliestHeadInChain recurses through parents until it finds the earliest one
func (h *Head) EarliestHeadInChain() commontypes.Head[common.Hash] {
	return h.EarliestInChain()
}

// IsInChain returns true if the given hash matches the hash of a head in the chain
func (h *Head) IsInChain(blockHash common.Hash) bool {
	for {
		if h.Hash == blockHash {
			return true
		}
		if h.Parent == nil {
			break
		}
		h = h.Parent
	}
	return false
}

// HashAtHeight returns the hash of the block at the given height, if it is in the chain.
// If not in chain, returns the zero hash
func (h *Head) HashAtHeight(blockNum int64) common.Hash {
	for {
		if h.Number == blockNum {
			return h.Hash
		}
		if h.Parent == nil {
			break
		}
		h = h.Parent
	}
	return common.Hash{}
}

// ChainLength returns the length of the chain followed by recursively looking up parents
func (h *Head) ChainLength() uint32 {
	if h == nil {
		return 0
	}
	l := uint32(1)

	for {
		if h.Parent == nil {
			break
		}
		l++
		if h == h.Parent {
			panic("circular reference detected")
		}
		h = h.Parent
	}
	return l
}

// ChainHashes returns an array of block hashes by recursively looking up parents
func (h *Head) ChainHashes() []common.Hash {
	var hashes []common.Hash

	for {
		hashes = append(hashes, h.Hash)
		if h.Parent == nil {
			break
		}
		if h == h.Parent {
			panic("circular reference detected")
		}
		h = h.Parent
	}
	return hashes
}

func (h *Head) LatestFinalizedHead() commontypes.Head[common.Hash] {
	for h != nil {
		if h.IsFinalized {
			return h
		}

		h = h.Parent
	}

	return nil
}

func (h *Head) ChainID() *big.Int {
	return h.EVMChainID.ToInt()
}

func (h *Head) HasChainID() bool {
	return h.EVMChainID != nil
}

func (h *Head) IsValid() bool {
	return h != nil
}

func (h *Head) ChainString() string {
	var sb strings.Builder

	for {
		sb.WriteString(h.String())
		if h.Parent == nil {
			break
		}
		if h == h.Parent {
			panic("circular reference detected")
		}
		sb.WriteString("->")
		h = h.Parent
	}
	sb.WriteString("->nil")
	return sb.String()
}

// String returns a string representation of this head
func (h *Head) String() string {
	return fmt.Sprintf("Head{Number: %d, Hash: %s, ParentHash: %s}", h.ToInt(), h.Hash.Hex(), h.ParentHash.Hex())
}

// ToInt return the height as a *big.Int. Also handles nil by returning nil.
func (h *Head) ToInt() *big.Int {
	if h == nil {
		return nil
	}
	return big.NewInt(h.Number)
}

// GreaterThan compares BlockNumbers and returns true if the receiver BlockNumber is greater than
// the supplied BlockNumber
func (h *Head) GreaterThan(r *Head) bool {
	if h == nil {
		return false
	}
	if h != nil && r == nil {
		return true
	}
	return h.Number > r.Number
}

// NextInt returns the next BlockNumber as big.int, or nil if nil to represent latest.
func (h *Head) NextInt() *big.Int {
	if h == nil {
		return nil
	}
	return new(big.Int).Add(h.ToInt(), big.NewInt(1))
}

// AsSlice returns a slice of heads up to length k
// len(heads) may be less than k if the available chain is not long enough
func (h *Head) AsSlice(k int) (heads []*Head) {
	if k < 1 || h == nil {
		return
	}
	heads = make([]*Head, 1)
	heads[0] = h
	for len(heads) < k && h.Parent != nil {
		h = h.Parent
		heads = append(heads, h)
	}
	return
}

func (h *Head) UnmarshalJSON(bs []byte) error {
	type head struct {
		Hash             common.Hash    `json:"hash"`
		Number           *hexutil.Big   `json:"number"`
		ParentHash       common.Hash    `json:"parentHash"`
		Timestamp        hexutil.Uint64 `json:"timestamp"`
		L1BlockNumber    *hexutil.Big   `json:"l1BlockNumber"`
		BaseFeePerGas    *hexutil.Big   `json:"baseFeePerGas"`
		ReceiptsRoot     common.Hash    `json:"receiptsRoot"`
		TransactionsRoot common.Hash    `json:"transactionsRoot"`
		StateRoot        common.Hash    `json:"stateRoot"`
		Difficulty       *hexutil.Big   `json:"difficulty"`
		TotalDifficulty  *hexutil.Big   `json:"totalDifficulty"`
	}

	var jsonHead head
	err := json.Unmarshal(bs, &jsonHead)
	if err != nil {
		return err
	}

	if jsonHead.Number == nil {
		*h = Head{}
		return nil
	}

	h.Hash = jsonHead.Hash
	h.Number = (*big.Int)(jsonHead.Number).Int64()
	h.ParentHash = jsonHead.ParentHash
	h.Timestamp = time.Unix(int64(jsonHead.Timestamp), 0).UTC()
	h.BaseFeePerGas = assets.NewWei((*big.Int)(jsonHead.BaseFeePerGas))
	if jsonHead.L1BlockNumber != nil {
		h.L1BlockNumber = sql.NullInt64{Int64: (*big.Int)(jsonHead.L1BlockNumber).Int64(), Valid: true}
	}
	h.ReceiptsRoot = jsonHead.ReceiptsRoot
	h.TransactionsRoot = jsonHead.TransactionsRoot
	h.StateRoot = jsonHead.StateRoot
	h.Difficulty = jsonHead.Difficulty.ToInt()
	h.TotalDifficulty = jsonHead.TotalDifficulty.ToInt()
	return nil
}

func (h *Head) MarshalJSON() ([]byte, error) {
	type head struct {
		Hash             *common.Hash    `json:"hash,omitempty"`
		Number           *hexutil.Big    `json:"number,omitempty"`
		ParentHash       *common.Hash    `json:"parentHash,omitempty"`
		Timestamp        *hexutil.Uint64 `json:"timestamp,omitempty"`
		ReceiptsRoot     *common.Hash    `json:"receiptsRoot,omitempty"`
		TransactionsRoot *common.Hash    `json:"transactionsRoot,omitempty"`
		StateRoot        *common.Hash    `json:"stateRoot,omitempty"`
		Difficulty       *hexutil.Big    `json:"difficulty,omitempty"`
		TotalDifficulty  *hexutil.Big    `json:"totalDifficulty,omitempty"`
	}

	var jsonHead head
	if h.Hash != (common.Hash{}) {
		jsonHead.Hash = &h.Hash
	}
	if h.ReceiptsRoot != (common.Hash{}) {
		jsonHead.ReceiptsRoot = &h.ReceiptsRoot
	}
	if h.TransactionsRoot != (common.Hash{}) {
		jsonHead.TransactionsRoot = &h.TransactionsRoot
	}
	if h.StateRoot != (common.Hash{}) {
		jsonHead.StateRoot = &h.StateRoot
	}
	jsonHead.Number = (*hexutil.Big)(big.NewInt(h.Number))
	if h.ParentHash != (common.Hash{}) {
		jsonHead.ParentHash = &h.ParentHash
	}
	if h.Timestamp != (time.Time{}) {
		t := hexutil.Uint64(h.Timestamp.UTC().Unix())
		jsonHead.Timestamp = &t
	}
	jsonHead.Difficulty = (*hexutil.Big)(h.Difficulty)
	jsonHead.TotalDifficulty = (*hexutil.Big)(h.TotalDifficulty)
	return json.Marshal(jsonHead)
}

// Block represents an ethereum block
// This type is only used for the block history estimator, and can be expensive to unmarshal. Don't add unnecessary fields here.
type Block struct {
	Number        int64
	Hash          common.Hash
	ParentHash    common.Hash
	BaseFeePerGas *assets.Wei
	Timestamp     time.Time
	Transactions  []Transaction
}

// MarshalJSON implements json marshalling for Block
func (b Block) MarshalJSON() ([]byte, error) {
	bi := &blocks.BlockInternal{
		Number:        hexutil.EncodeBig(big.NewInt(b.Number)),
		Hash:          b.Hash,
		ParentHash:    b.ParentHash,
		BaseFeePerGas: (*hexutil.Big)(b.BaseFeePerGas),
		Timestamp:     (hexutil.Uint64)(uint64(b.Timestamp.Unix())),
		Transactions:  toInternalTxnSlice(b.Transactions),
	}

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	enc := codec.NewEncoder(buf, &codec.JsonHandle{})
	err := enc.Encode(bi)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

var ErrMissingBlock = pkgerrors.New("missing block")

// UnmarshalJSON unmarshals to a Block
func (b *Block) UnmarshalJSON(data []byte) error {
	var h codec.Handle = new(codec.JsonHandle)
	bi := blocks.BlockInternal{}

	dec := codec.NewDecoderBytes(data, h)
	err := dec.Decode(&bi)

	if err != nil {
		return err
	}
	if bi.Empty() {
		return pkgerrors.WithStack(ErrMissingBlock)
	}

	n, err := hexutil.DecodeBig(bi.Number)
	if err != nil {
		return pkgerrors.Wrapf(err, "failed to decode block number while unmarshalling block, got:  '%s' in '%s'", bi.Number, data)
	}
	*b = Block{
		Number:        n.Int64(),
		Hash:          bi.Hash,
		ParentHash:    bi.ParentHash,
		BaseFeePerGas: (*assets.Wei)(bi.BaseFeePerGas),
		Timestamp:     time.Unix((int64((uint64)(bi.Timestamp))), 0),
		Transactions:  fromInternalTxnSlice(bi.Transactions),
	}
	return nil
}

// thin public wrapper for internal type of the same name
// and which has to be internal for JSON un/marshal'ing code gen consistency
type TxType uint8

// Transaction represents an ethereum transaction
// Use our own type because geth's type has validation failures on e.g. zero
// gas used, which can occur on other chains.
// This type is only used for the block history estimator, and can be expensive to unmarshal. Don't add unnecessary fields here.
type Transaction struct {
	GasPrice             *assets.Wei `json:"gasPrice"`
	GasLimit             uint32      `json:"gasLimit"`
	MaxFeePerGas         *assets.Wei `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *assets.Wei `json:"maxPriorityFeePerGas"`
	Type                 TxType      `json:"type"`
	Hash                 common.Hash `json:"hash"`
}

const LegacyTxType = blocks.TxType(0x0)

// UnmarshalJSON unmarshals a Transaction
func (t *Transaction) UnmarshalJSON(data []byte) error {
	var h codec.Handle = new(codec.JsonHandle)
	ti := blocks.TransactionInternal{}

	dec := codec.NewDecoderBytes(data, h)
	err := dec.Decode(&ti)

	if err != nil {
		return err
	}

	if ti.Gas == nil {
		return pkgerrors.Errorf("expected 'gas' to not be null, got: '%s'", data)
	}
	if ti.Type == nil {
		tpe := LegacyTxType
		ti.Type = &tpe
	}
	*t = fromInternalTxn(ti)

	return nil
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	ti := toInternalTxn(*t)

	buf := bytes.NewBuffer(make([]byte, 0, 256))
	enc := codec.NewEncoder(buf, &codec.JsonHandle{})

	err := enc.Encode(ti)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// WeiPerEth is amount of Wei currency units in one Eth.
var WeiPerEth = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

// ChainlinkFulfilledTopic is the signature for the event emitted after calling
// ChainlinkClient.validateChainlinkCallback(requestId). See
// ../../contracts/src/v0.8/ChainlinkClient.sol
var ChainlinkFulfilledTopic = utils.MustHash("ChainlinkFulfilled(bytes32)")

// ReceiptIndicatesRunLogFulfillment returns true if this tx receipt is the result of a
// fulfilled run log.
func ReceiptIndicatesRunLogFulfillment(txr types.Receipt) bool {
	for _, log := range txr.Logs {
		if log.Topics[0] == ChainlinkFulfilledTopic {
			return true
		}
	}
	return false
}

// FunctionSelector is the first four bytes of the call data for a
// function call and specifies the function to be called.
type FunctionSelector [FunctionSelectorLength]byte

// FunctionSelectorLength should always be a length of 4 as a byte.
const FunctionSelectorLength = 4

// BytesToFunctionSelector converts the given bytes to a FunctionSelector.
func BytesToFunctionSelector(b []byte) FunctionSelector {
	var f FunctionSelector
	f.SetBytes(b)
	return f
}

// HexToFunctionSelector converts the given string to a FunctionSelector.
func HexToFunctionSelector(s string) FunctionSelector {
	return BytesToFunctionSelector(common.FromHex(s))
}

// String returns the FunctionSelector as a string type.
func (f FunctionSelector) String() string { return hexutil.Encode(f[:]) }

// Bytes returns the FunctionSelector as a byte slice
func (f FunctionSelector) Bytes() []byte { return f[:] }

// SetBytes sets the FunctionSelector to that of the given bytes (will trim).
func (f *FunctionSelector) SetBytes(b []byte) { copy(f[:], b[:FunctionSelectorLength]) }

var hexRegexp = regexp.MustCompile("^[0-9a-fA-F]*$")

func unmarshalFromString(s string, f *FunctionSelector) error {
	if hex.HasPrefix(s) {
		if !hexRegexp.Match([]byte(s)[2:]) {
			return fmt.Errorf("function selector %s must be 0x-hex encoded", s)
		}
		bytes := common.FromHex(s)
		if len(bytes) != FunctionSelectorLength {
			return pkgerrors.New("function ID must be 4 bytes in length")
		}
		f.SetBytes(bytes)
	} else {
		bytes, err := utils.Keccak256([]byte(s))
		if err != nil {
			return err
		}
		f.SetBytes(bytes[0:4])
	}
	return nil
}

// UnmarshalJSON parses the raw FunctionSelector and sets the FunctionSelector
// type to the given input.
func (f *FunctionSelector) UnmarshalJSON(input []byte) error {
	var s string
	err := json.Unmarshal(input, &s)
	if err != nil {
		return err
	}
	return unmarshalFromString(s, f)
}

// MarshalJSON returns the JSON encoding of f
func (f FunctionSelector) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}

// Value returns this instance serialized for database storage
func (f FunctionSelector) Value() (driver.Value, error) {
	return f.Bytes(), nil
}

// Scan returns the selector from its serialization in the database
func (f *FunctionSelector) Scan(value interface{}) error {
	temp, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("unable to convent %v of type %T to FunctionSelector", value, value)
	}
	if len(temp) != FunctionSelectorLength {
		return fmt.Errorf("function selector %v should have length %d, but has length %d",
			temp, FunctionSelectorLength, len(temp))
	}
	copy(f[:], temp)
	return nil
}

// This data can contain anything and is submitted by user on-chain, so we must
// be extra careful how we interact with it
type UntrustedBytes []byte

// SafeByteSlice returns an error on out of bounds access to a byte array, where a
// normal slice would panic instead
func (ary UntrustedBytes) SafeByteSlice(start int, end int) ([]byte, error) {
	if end > len(ary) || start > end || start < 0 || end < 0 {
		var empty []byte
		return empty, pkgerrors.New("out of bounds slice access")
	}
	return ary[start:end], nil
}

// toInternalTxn converts a Transaction into the internal intermediate representation
func toInternalTxn(txn Transaction) blocks.TransactionInternal {
	gas := (hexutil.Uint64)(uint64(txn.GasLimit))
	itype := blocks.TxType(txn.Type)
	return blocks.TransactionInternal{
		GasPrice:             (*hexutil.Big)(txn.GasPrice),
		Gas:                  &gas,
		MaxFeePerGas:         (*hexutil.Big)(txn.MaxFeePerGas),
		MaxPriorityFeePerGas: (*hexutil.Big)(txn.MaxPriorityFeePerGas),
		Type:                 &itype,
		Hash:                 txn.Hash,
	}
}

// toInternalTxn converts a []Transaction into the internal intermediate representation
func toInternalTxnSlice(txns []Transaction) []blocks.TransactionInternal {
	out := make([]blocks.TransactionInternal, len(txns))
	for i, txn := range txns {
		out[i] = toInternalTxn(txn)
	}
	return out
}

// fromInternalTxn converts an internal intermediate representation into a Transaction
func fromInternalTxn(ti blocks.TransactionInternal) Transaction {
	if ti.Type == nil {
		tpe := LegacyTxType
		ti.Type = &tpe
	}
	return Transaction{
		GasPrice:             (*assets.Wei)(ti.GasPrice),
		GasLimit:             uint32(*ti.Gas),
		MaxFeePerGas:         (*assets.Wei)(ti.MaxFeePerGas),
		MaxPriorityFeePerGas: (*assets.Wei)(ti.MaxPriorityFeePerGas),
		Type:                 TxType(*ti.Type),
		Hash:                 ti.Hash,
	}
}

// fromInternalTxnSlice converts a slice of internal intermediate representation into a []Transaction
func fromInternalTxnSlice(tis []blocks.TransactionInternal) []Transaction {
	out := make([]Transaction, len(tis))
	for i, ti := range tis {
		out[i] = fromInternalTxn(ti)
	}
	return out
}
