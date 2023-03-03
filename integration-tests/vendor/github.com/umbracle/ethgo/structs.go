package ethgo

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

var (
	// ZeroAddress is an address of all zeros
	ZeroAddress = Address{}

	// ZeroHash is a hash of all zeros
	ZeroHash = Hash{}
)

// Address is an Ethereum address
type Address [20]byte

// HexToAddress converts an hex string value to an address object
func HexToAddress(str string) Address {
	a := Address{}
	a.UnmarshalText(completeHex(str, 20))
	return a
}

// BytesToAddress converts bytes to an address object
func BytesToAddress(b []byte) Address {
	var a Address

	size := len(b)
	min := min(size, 20)

	copy(a[20-min:], b[len(b)-min:])
	return a
}

// Address implements the ethgo.Key interface Address method.
func (a Address) Address() Address {
	return a
}

// Sign implements the ethgo.Key interface Sign method.
func (a Address) Sign(hash []byte) ([]byte, error) {
	panic("an address cannot sign messages")
}

// UnmarshalText implements the unmarshal interface
func (a *Address) UnmarshalText(b []byte) error {
	return unmarshalTextByte(a[:], b, 20)
}

// MarshalText implements the marshal interface
func (a Address) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

// Bytes returns the bytes of the Address
func (a Address) Bytes() []byte {
	return a[:]
}

func (a Address) String() string {
	return a.checksumEncode()
}

func (a Address) checksumEncode() string {
	address := strings.ToLower(hex.EncodeToString(a[:]))
	hash := hex.EncodeToString(Keccak256([]byte(address)))

	ret := "0x"
	for i := 0; i < len(address); i++ {
		character := string(address[i])

		num, _ := strconv.ParseInt(string(hash[i]), 16, 64)
		if num > 7 {
			ret += strings.ToUpper(character)
		} else {
			ret += character
		}
	}

	return ret
}

// Hash is an Ethereum hash
type Hash [32]byte

// HexToHash converts an hex string value to a hash object
func HexToHash(str string) Hash {
	h := Hash{}
	h.UnmarshalText(completeHex(str, 32))
	return h
}

// BytesToHash converts bytes to a hash object
func BytesToHash(b []byte) Hash {
	var h Hash

	size := len(b)
	min := min(size, 32)

	copy(h[32-min:], b[len(b)-min:])
	return h
}

// UnmarshalText implements the unmarshal interface
func (h *Hash) UnmarshalText(b []byte) error {
	return unmarshalTextByte(h[:], b, 32)
}

// MarshalText implements the marshal interface
func (h Hash) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

// Bytes returns the bytes of the Hash
func (h Hash) Bytes() []byte {
	return h[:]
}

func (h Hash) String() string {
	return "0x" + hex.EncodeToString(h[:])
}

func (h Hash) Location() string {
	return h.String()
}

type Block struct {
	Number             uint64
	Hash               Hash
	ParentHash         Hash
	Sha3Uncles         Hash
	TransactionsRoot   Hash
	StateRoot          Hash
	ReceiptsRoot       Hash
	Miner              Address
	Difficulty         *big.Int
	ExtraData          []byte
	GasLimit           uint64
	GasUsed            uint64
	Timestamp          uint64
	Transactions       []*Transaction
	TransactionsHashes []Hash
	Uncles             []Hash
}

func (b *Block) Copy() *Block {
	bb := new(Block)
	*bb = *b
	if b.Difficulty != nil {
		bb.Difficulty = new(big.Int).Set(b.Difficulty)
	}
	bb.ExtraData = append(bb.ExtraData[:0], b.ExtraData...)
	bb.Transactions = make([]*Transaction, len(b.Transactions))
	for indx, txn := range b.Transactions {
		bb.Transactions[indx] = txn.Copy()
	}
	return bb
}

type TransactionType int

const (
	TransactionLegacy TransactionType = 0
	// eip-2930
	TransactionAccessList TransactionType = 1
	// eip-1559
	TransactionDynamicFee TransactionType = 2
)

type Transaction struct {
	Type TransactionType

	// legacy values
	Hash     Hash
	From     Address
	To       *Address
	Input    []byte
	GasPrice uint64
	Gas      uint64
	Value    *big.Int
	Nonce    uint64
	V        []byte
	R        []byte
	S        []byte

	// jsonrpc values
	BlockHash   Hash
	BlockNumber uint64
	TxnIndex    uint64

	// eip-2930 values
	ChainID    *big.Int
	AccessList AccessList

	// eip-1559 values
	MaxPriorityFeePerGas *big.Int
	MaxFeePerGas         *big.Int
}

func (t *Transaction) Copy() *Transaction {
	tt := new(Transaction)
	if t.To != nil {
		to := Address(*t.To)
		tt.To = &to
	}
	tt.Input = append(tt.Input[:0], t.Input...)
	if t.Value != nil {
		tt.Value = new(big.Int).Set(t.Value)
	}
	tt.V = append(tt.V[:0], t.V...)
	tt.R = append(tt.R[:0], t.R...)
	tt.S = append(tt.S[:0], t.S...)
	if t.ChainID != nil {
		tt.ChainID = new(big.Int).Set(t.ChainID)
	}
	if t.MaxPriorityFeePerGas != nil {
		tt.MaxPriorityFeePerGas = new(big.Int).Set(t.MaxPriorityFeePerGas)
	}
	if t.MaxFeePerGas != nil {
		tt.MaxFeePerGas = new(big.Int).Set(t.MaxFeePerGas)
	}
	return tt
}

type AccessEntry struct {
	Address Address `json:"address"`
	Storage []Hash  `json:"storageKeys"`
}

type AccessList []AccessEntry

type CallMsg struct {
	From     Address
	To       *Address
	Data     []byte
	GasPrice uint64
	Gas      *big.Int
	Value    *big.Int
}

type LogFilter struct {
	Address   []Address
	Topics    [][]*Hash
	BlockHash *Hash
	From      *BlockNumber
	To        *BlockNumber
}

func (l *LogFilter) SetFromUint64(num uint64) {
	b := BlockNumber(num)
	l.From = &b
}

func (l *LogFilter) SetToUint64(num uint64) {
	b := BlockNumber(num)
	l.To = &b
}

func (l *LogFilter) SetTo(b BlockNumber) {
	l.To = &b
}

type Receipt struct {
	TransactionHash   Hash
	TransactionIndex  uint64
	ContractAddress   Address
	BlockHash         Hash
	From              Address
	BlockNumber       uint64
	GasUsed           uint64
	CumulativeGasUsed uint64
	LogsBloom         []byte
	Logs              []*Log
	Status            uint64
}

func (r *Receipt) Copy() *Receipt {
	rr := new(Receipt)
	*rr = *r
	rr.LogsBloom = append(rr.LogsBloom[:0], r.LogsBloom...)
	rr.Logs = make([]*Log, len(r.Logs))
	for indx, log := range r.Logs {
		rr.Logs[indx] = log.Copy()
	}
	return rr
}

type Log struct {
	Removed          bool
	LogIndex         uint64
	TransactionIndex uint64
	TransactionHash  Hash
	BlockHash        Hash
	BlockNumber      uint64
	Address          Address
	Topics           []Hash
	Data             []byte
}

func (l *Log) Copy() *Log {
	ll := new(Log)
	*ll = *l
	ll.Data = append(ll.Data[:0], l.Data...)
	return ll
}

type BlockNumber int

const (
	Latest   BlockNumber = -1
	Earliest BlockNumber = -2
	Pending  BlockNumber = -3
)

func (b BlockNumber) Location() string {
	return b.String()
}

func (b BlockNumber) String() string {
	switch b {
	case Latest:
		return "latest"
	case Earliest:
		return "earliest"
	case Pending:
		return "pending"
	}
	if b < 0 {
		panic("internal. blocknumber is negative")
	}
	return fmt.Sprintf("0x%x", uint64(b))
}

func EncodeBlock(block ...BlockNumber) BlockNumber {
	if len(block) != 1 {
		return Latest
	}
	return block[0]
}

type BlockNumberOrHash interface {
	Location() string
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}

type Key interface {
	Address() Address
	Sign(hash []byte) ([]byte, error)
}

func completeHex(str string, num int) []byte {
	num = num * 2
	str = strings.TrimPrefix(str, "0x")

	size := len(str)
	if size < num {
		for i := size; i < num; i++ {
			str = "0" + str
		}
	} else {
		diff := size - num
		str = str[diff:]
	}
	return []byte("0x" + str)
}
