package evm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

type EventBySigFilter struct {
	Address  common.Address
	EventSig common.Hash
}

func NewEventBySigFilter(address common.Address, eventSig common.Hash) query.Expression {
	var searchEventFilter *EventBySigFilter
	searchEventFilter.Address = address
	searchEventFilter.EventSig = eventSig
	return query.Expression{Primitive: searchEventFilter}
}

func (f *EventBySigFilter) Accept(visitor primitives.Visitor) {
	switch v := visitor.(type) {
	case *PgDSLParser:
		v.VisitEventBySigFilter(f)
	}
}

type EventByTopicFilter struct {
	EventSig         common.Hash
	Topic            uint64
	ValueComparators []primitives.ValueComparator
}

func NewEventByTopicFilter(eventSig common.Hash, topicIndex uint64, valueComparators []primitives.ValueComparator) query.Expression {
	var eventByIndexFilter *EventByTopicFilter
	eventByIndexFilter.EventSig = eventSig
	eventByIndexFilter.Topic = topicIndex
	eventByIndexFilter.ValueComparators = valueComparators

	return query.Expression{Primitive: eventByIndexFilter}
}

func (f *EventByTopicFilter) Accept(visitor primitives.Visitor) {
	switch v := visitor.(type) {
	case *PgDSLParser:
		v.VisitEventTopicsByValueFilter(f)
	}
}

type EventByWordFilter struct {
	EventSig         common.Hash
	WordIndex        uint8
	ValueComparators []primitives.ValueComparator
}

func NewEventByWordFilter(eventSig common.Hash, wordIndex uint8, valueComparators []primitives.ValueComparator) query.Expression {
	var eventByIndexFilter *EventByWordFilter
	eventByIndexFilter.EventSig = eventSig
	eventByIndexFilter.WordIndex = wordIndex
	eventByIndexFilter.ValueComparators = valueComparators
	return query.Expression{Primitive: eventByIndexFilter}
}

func (f *EventByWordFilter) Accept(visitor primitives.Visitor) {
	switch v := visitor.(type) {
	case *PgDSLParser:
		v.VisitEventByWordFilter(f)
	}
}

type FinalityFilter struct {
	Confs evmtypes.Confirmations
}

func NewFinalityFilter(filter *primitives.Confirmations) (query.Expression, error) {
	// TODO chain agnostic confidence levels that map to evm finality
	switch filter.ConfirmationLevel {
	case primitives.Finalized:
		return query.Expression{Primitive: &FinalityFilter{evmtypes.Finalized}}, nil
	case primitives.Unconfirmed:
		return query.Expression{Primitive: &FinalityFilter{evmtypes.Unconfirmed}}, nil
	default:
		return query.Expression{}, fmt.Errorf("invalid finality confirmations filter value %v", filter.ConfirmationLevel)
	}
}

func (f *FinalityFilter) Accept(visitor primitives.Visitor) {
	switch v := visitor.(type) {
	case *PgDSLParser:
		v.VisitFinalityFilter(f)
	}
}
