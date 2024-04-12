package evm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
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

func (f *EventBySigFilter) Accept(visitor query.Visitor) {
	switch v := visitor.(type) {
	case *PgDSLParser:
		v.VisitEventBySigFilter(f)
	}
}

type EventByTopicFilter struct {
	EventSig       common.Hash
	Topic          uint64
	ValueComparers []query.ValueComparer
}

func NewEventByTopicFilter(eventSig common.Hash, topicIndex uint64, valueComparators []query.ValueComparer) query.Expression {
	var eventByIndexFilter *EventByTopicFilter
	eventByIndexFilter.EventSig = eventSig
	eventByIndexFilter.Topic = topicIndex
	eventByIndexFilter.ValueComparers = valueComparators

	return query.Expression{Primitive: eventByIndexFilter}
}

func (f *EventByTopicFilter) Accept(visitor query.Visitor) {
	switch v := visitor.(type) {
	case *PgDSLParser:
		v.VisitEventTopicsByValueFilter(f)
	}
}

type EventByWordFilter struct {
	EventSig         common.Hash
	WordIndex        uint8
	ValueComparators []query.ValueComparer
}

func NewEventByWordFilter(eventSig common.Hash, wordIndex uint8, valueComparators []query.ValueComparer) query.Expression {
	var eventByIndexFilter *EventByWordFilter
	eventByIndexFilter.EventSig = eventSig
	eventByIndexFilter.WordIndex = wordIndex
	eventByIndexFilter.ValueComparators = valueComparators
	return query.Expression{Primitive: eventByIndexFilter}
}

func (f *EventByWordFilter) Accept(visitor query.Visitor) {
	switch v := visitor.(type) {
	case *PgDSLParser:
		v.VisitEventByWordFilter(f)
	}
}

type FinalityFilter struct {
	Confs evmtypes.Confirmations
}

func NewFinalityFilter(filter *query.ConfirmationsPrimitive) (query.Expression, error) {
	switch filter.ConfirmationLevel {
	case query.Finalized:
		return query.Expression{Primitive: &FinalityFilter{evmtypes.Finalized}}, nil
	case query.Unconfirmed:
		return query.Expression{Primitive: &FinalityFilter{evmtypes.Unconfirmed}}, nil
	default:
		return query.Expression{}, fmt.Errorf("invalid finality confirmations filter value %v", filter.ConfirmationLevel)
	}
}

func (f *FinalityFilter) Accept(visitor query.Visitor) {
	switch v := visitor.(type) {
	case *PgDSLParser:
		v.VisitFinalityFilter(f)
	}
}
