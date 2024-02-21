package logpoller

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

type ComparisonOperator int

const (
	Eq ComparisonOperator = iota
	Neq
	Gt
	Lt
	Gte
	Lte
)

type SortDirection int

const (
	Asc SortDirection = iota
	Desc
)

var (
	DefaultSortAndLimit = SortAndLimit{
		sortBy: []SortBy{
			{field: "block_number", dir: Asc},
			{field: "log_index", dir: Asc},
		},
		limit: 0,
	}
)

type SortAndLimit struct {
	sortBy []SortBy
	limit  int
}

type SortBy struct {
	field string
	dir   SortDirection
}

func NewSortAndLimit(limit int, sortBy ...SortBy) SortAndLimit {
	return SortAndLimit{sortBy: sortBy, limit: limit}
}

func NewSortBy(field string, dir SortDirection) SortBy {
	return SortBy{field: field, dir: dir}
}

type QFilter interface {
	accept(visitor Visitor)
}

type AndFilter struct {
	filters []QFilter
}

func NewAndFilter(filters ...QFilter) *AndFilter {
	return &AndFilter{filters: filters}
}

func NewBasicAndFilter(address common.Address, eventSig common.Hash, filters ...QFilter) *AndFilter {
	allFilters := make([]QFilter, 0, len(filters)+2)
	allFilters = append(allFilters, NewAddressFilter(address), NewEventSigFilter(eventSig))
	allFilters = append(allFilters, filters...)
	return NewAndFilter(allFilters...)
}

func AppendedNewFilter(root *AndFilter, other ...QFilter) *AndFilter {
	filters := make([]QFilter, 0, len(root.filters)+len(other))
	filters = append(filters, root.filters...)
	filters = append(filters, other...)
	return NewAndFilter(filters...)
}

func (f *AndFilter) accept(visitor Visitor) {
	visitor.VisitAndFilter(*f)
}

type EvmChainIdFilter struct {
	chainId *ubig.Big
}

func NewEvmChainIdFilter(chainId *big.Int) *EvmChainIdFilter {
	return &EvmChainIdFilter{chainId: ubig.New(chainId)}
}

func (f *EvmChainIdFilter) accept(visitor Visitor) {
	visitor.VisitEvmChainIdFilter(*f)
}

type AddressFilter struct {
	address []common.Address
}

func NewAddressFilter(address ...common.Address) *AddressFilter {
	return &AddressFilter{address: address}
}

func (f *AddressFilter) accept(visitor Visitor) {
	visitor.VisitAddressFilter(*f)
}

type EventSigFilter struct {
	eventSig []common.Hash
}

func NewEventSigFilter(eventSig ...common.Hash) *EventSigFilter {
	return &EventSigFilter{eventSig: eventSig}
}

func (f *EventSigFilter) accept(visitor Visitor) {
	visitor.VisitEventSigFilter(*f)
}

type DataWordFilter struct {
	index    int
	operator ComparisonOperator
	value    common.Hash
}

func NewDataWordFilter(index int, operator ComparisonOperator, value common.Hash) *DataWordFilter {
	return &DataWordFilter{index: index, operator: operator, value: value}
}

func NewDataWordGteFilter(index int, value common.Hash) *DataWordFilter {
	return NewDataWordFilter(index, Gte, value)
}

func NewDataWordLteFilter(index int, value common.Hash) *DataWordFilter {
	return NewDataWordFilter(index, Lte, value)
}

func (f *DataWordFilter) accept(visitor Visitor) {
	visitor.VisitDataWordFilter(*f)
}

type TopicFilter struct {
	index    int
	operator ComparisonOperator
	value    common.Hash
}

func NewTopicFilter(index int, operator ComparisonOperator, value common.Hash) *TopicFilter {
	return &TopicFilter{index: index, operator: operator, value: value}
}

func NewTopicRangeFilter(index int, topicValueMin, topicValueMax common.Hash) *AndFilter {
	return NewAndFilter(
		NewTopicFilter(index, Gte, topicValueMin),
		NewTopicFilter(index, Lte, topicValueMax),
	)
}

func (f *TopicFilter) accept(visitor Visitor) {
	visitor.VisitTopicFilter(*f)
}

type TopicsFilter struct {
	index  int
	values []common.Hash
}

func NewTopicsFilter(index int, values ...common.Hash) *TopicsFilter {
	return &TopicsFilter{index: index, values: values}
}

func (f *TopicsFilter) accept(visitor Visitor) {
	visitor.VisitTopicsFilter(*f)
}

type ConfirmationFilter struct {
	confs Confirmations
}

func NewConfirmationFilter(confs Confirmations) *ConfirmationFilter {
	return &ConfirmationFilter{confs: confs}
}

func (f *ConfirmationFilter) accept(visitor Visitor) {
	visitor.VisitConfirmationFilter(*f)
}

func NewBlockFilter(block int64, operator ComparisonOperator) *BlockFilter {
	return &BlockFilter{operator, block}
}

func NewBlockRangeFilter(start, end int64) *AndFilter {
	return NewAndFilter(
		NewBlockFilter(start, Gte),
		NewBlockFilter(end, Lte),
	)
}

type BlockFilter struct {
	operator ComparisonOperator
	block    int64
}

func (f *BlockFilter) accept(visitor Visitor) {
	visitor.VisitBlockFilter(*f)
}

func NewBlockTimestampAfterFilter(after time.Time) *BlockTimestampFilter {
	return NewBlockTimeStampFilter(after, Gt)
}

func NewBlockTimeStampFilter(timestamp time.Time, operator ComparisonOperator) *BlockTimestampFilter {
	return &BlockTimestampFilter{operator, timestamp}
}

type BlockTimestampFilter struct {
	operator  ComparisonOperator
	timestamp time.Time
}

func (f *BlockTimestampFilter) accept(visitor Visitor) {
	visitor.VisitBlockTimestampFilter(*f)
}

func NewTxHashFilter(txHash common.Hash) *TxHashFilter {
	return &TxHashFilter{txHash}
}

type TxHashFilter struct {
	txHash common.Hash
}

func (f *TxHashFilter) accept(visitor Visitor) {
	visitor.VisitTxHashFilter(*f)
}

type Visitor interface {
	VisitAndFilter(node AndFilter)
	VisitEvmChainIdFilter(node EvmChainIdFilter)
	VisitAddressFilter(node AddressFilter)
	VisitEventSigFilter(node EventSigFilter)
	VisitDataWordFilter(node DataWordFilter)
	VisitTopicFilter(node TopicFilter)
	VisitTopicsFilter(node TopicsFilter)
	VisitBlockFilter(node BlockFilter)
	VisitConfirmationFilter(node ConfirmationFilter)
	VisitBlockTimestampFilter(node BlockTimestampFilter)
	VisitTxHashFilter(node TxHashFilter)
}
