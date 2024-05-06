package primitives

// Visitor should have a per chain per db type implementation that converts primitives to db queries.
type Visitor interface {
	Comparator(primitive Comparator)
	Block(primitive Block)
	Confirmations(primitive Confirmations)
	Timestamp(primitive Timestamp)
	TxHash(primitive TxHash)
}

// Primitive is the basic building block for KeyFilter.
type Primitive interface {
	Accept(visitor Visitor)
}

type ComparisonOperator int

const (
	Eq ComparisonOperator = iota
	Neq
	Gt
	Lt
	Gte
	Lte
)

type ValueComparator struct {
	Value    string
	Operator ComparisonOperator
}

// Comparator is used to filter over values that belong to key data.
type Comparator struct {
	Name             string
	ValueComparators []ValueComparator
}

func (f *Comparator) Accept(visitor Visitor) {
	visitor.Comparator(*f)
}

// Block is a primitive of KeyFilter that filters search in comparison to block number.
type Block struct {
	Block    uint64
	Operator ComparisonOperator
}

func (f *Block) Accept(visitor Visitor) {
	visitor.Block(*f)
}

type ConfirmationLevel int32

// TODO come up with a proper chain agnostic solution BCF-3156
const (
	Finalized   = ConfirmationLevel(-1)
	Unconfirmed = ConfirmationLevel(-2)
)

// Confirmations is a primitive of KeyFilter that filters search to results that have a certain level of confirmation.
// Confirmation map to different concepts on different blockchains.
type Confirmations struct {
	ConfirmationLevel
}

func (f *Confirmations) Accept(visitor Visitor) {
	visitor.Confirmations(*f)
}

// Timestamp is a primitive of KeyFilter that filters search in comparison to timestamp.
type Timestamp struct {
	Timestamp uint64
	Operator  ComparisonOperator
}

func (f *Timestamp) Accept(visitor Visitor) {
	visitor.Timestamp(*f)
}

// TxHash is a primitive of KeyFilter that filters search to results that contain txHash.
type TxHash struct {
	TxHash string
}

func (f *TxHash) Accept(visitor Visitor) {
	visitor.TxHash(*f)
}
