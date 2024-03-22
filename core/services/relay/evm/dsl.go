package evm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

type EventFilter struct {
	Address  common.Address
	EventSig common.Hash
}

func NewEventFilter(address common.Address, eventSig common.Hash) query.Expression {
	var searchEventFilter *EventFilter
	searchEventFilter.Address = address
	searchEventFilter.EventSig = eventSig
	return query.Expression{Primitive: searchEventFilter}
}

func (f *EventFilter) Accept(visitor query.Visitor) {
	switch v := visitor.(type) {
	case *PgDSLParser:
		v.VisitEventFilter(f)
	}
}

type EventByIndexFilter struct {
	Address          common.Address
	EventSig         common.Hash
	Topic            uint64
	ValueComparators []query.ValueComparator
}

func NewEventByIndexFilter(address common.Address, eventSig common.Hash, topicIndex uint64, valueComparators []query.ValueComparator) query.Expression {
	var eventByIndexFilter *EventByIndexFilter
	eventByIndexFilter.Address = address
	eventByIndexFilter.EventSig = eventSig
	eventByIndexFilter.Topic = topicIndex
	eventByIndexFilter.ValueComparators = valueComparators

	return query.Expression{Primitive: eventByIndexFilter}
}

func (f *EventByIndexFilter) Accept(visitor query.Visitor) {
	switch v := visitor.(type) {
	case *PgDSLParser:
		v.VisitEventTopicsByValueFilter(f)
	}
}

type EventByWordFilter struct {
	Address          common.Address
	EventSig         common.Hash
	WordIndex        uint8
	ValueComparators []query.ValueComparator
}

func NewEventByWordFilter(address common.Address, eventSig common.Hash, wordIndex uint8, valueComparators []query.ValueComparator) query.Expression {
	var eventByIndexFilter *EventByWordFilter
	eventByIndexFilter.Address = address
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

type ChainIdFilter struct {
	chainId *ubig.Big
}

func (f *ChainIdFilter) accept(visitor query.Visitor) {
	switch v := visitor.(type) {
	case *PgDSLParser:
		v.VisitChainIdFilter(f)
	}
}
