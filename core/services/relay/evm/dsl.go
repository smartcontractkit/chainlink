package evm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

type EventFilter struct {
	Address  common.Address
	EventSig common.Hash
}

func NewEventFilter(address common.Address, eventSig common.Hash) commontypes.Expression {
	var searchEventFilter *EventFilter
	searchEventFilter.Address = address
	searchEventFilter.EventSig = eventSig
	return commontypes.Expression{Primitive: searchEventFilter}
}

func (f *EventFilter) Accept(visitor commontypes.Visitor) {
	switch v := visitor.(type) {
	case *PgDSLParser:
		v.VisitEventFilter(f)
	}
}

type EventByIndexFilter struct {
	Address  common.Address
	EventSig common.Hash
	Topics   []int
	Values   []string
}

func NewEventByIndexFilter(address common.Address, values []string, eventSig common.Hash, topicIndex int) commontypes.Expression {
	var eventByIndexFilter *EventByIndexFilter
	eventByIndexFilter.Address = address
	eventByIndexFilter.EventSig = eventSig
	eventByIndexFilter.Topics = append(eventByIndexFilter.Topics, topicIndex)
	eventByIndexFilter.Values = append(eventByIndexFilter.Values, values...)

	return commontypes.Expression{Primitive: eventByIndexFilter}
}

func (f *EventByIndexFilter) Accept(visitor commontypes.Visitor) {
	switch v := visitor.(type) {
	case *PgDSLParser:
		v.VisitEventTopicsByValueFilter(f)
	}
}

type FinalityFilter struct {
	Confs evmtypes.Confirmations
}

func NewFinalityFilter(filter *commontypes.ConfirmationsFilter) (commontypes.Expression, error) {
	switch filter.Confirmations {
	case commontypes.Finalized:
		return commontypes.Expression{Primitive: &FinalityFilter{evmtypes.Finalized}}, nil
	case commontypes.Unconfirmed:
		return commontypes.Expression{Primitive: &FinalityFilter{evmtypes.Unconfirmed}}, nil
	default:
		return commontypes.Expression{}, fmt.Errorf("invalid finality confirmations filter value %v", filter.Confirmations)
	}
}

func (f *FinalityFilter) Accept(visitor commontypes.Visitor) {
	switch v := visitor.(type) {
	case *PgDSLParser:
		v.VisitFinalityFilter(f)
	}
}

type ChainIdFilter struct {
	chainId *ubig.Big
}

func (f *ChainIdFilter) accept(visitor commontypes.Visitor) {
	switch v := visitor.(type) {
	case *PgDSLParser:
		v.VisitChainIdFilter(f)
	}
}
