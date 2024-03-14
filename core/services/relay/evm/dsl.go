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

func NewEventFilter(address common.Address, eventSig common.Hash) *EventFilter {
	var searchEventFilter *EventFilter
	searchEventFilter.Address = address
	searchEventFilter.EventSig = eventSig
	return searchEventFilter
}

func (f *EventFilter) Accept(visitor commontypes.Visitor) {
	switch v := visitor.(type) {
	case *PgDSLParser:
		v.VisitEventFilter(f)
	}
}

type EventTopicByValuesFilter struct {
	Address  common.Address
	EventSig common.Hash
	Topics   []int
	Values   []string
}

func NewEventTopicsByValueFilter(address common.Address, values []string, eventSig common.Hash, topicIndex int) *EventTopicByValuesFilter {
	var searchEventTopicsByValueFilter *EventTopicByValuesFilter
	searchEventTopicsByValueFilter.Address = address
	searchEventTopicsByValueFilter.EventSig = eventSig
	searchEventTopicsByValueFilter.Topics = append(searchEventTopicsByValueFilter.Topics, topicIndex)
	searchEventTopicsByValueFilter.Values = append(searchEventTopicsByValueFilter.Values, values...)

	return searchEventTopicsByValueFilter
}

func (f *EventTopicByValuesFilter) Accept(visitor commontypes.Visitor) {
	switch v := visitor.(type) {
	case *PgDSLParser:
		v.VisitEventTopicsByValueFilter(f)
	}
}

type FinalityFilter struct {
	Confs evmtypes.Confirmations
}

func NewFinalityFilter(filter *commontypes.ConfirmationsFilter) (*FinalityFilter, error) {
	switch filter.Confirmations {
	case commontypes.Finalized:
		return &FinalityFilter{evmtypes.Finalized}, nil
	case commontypes.Unconfirmed:
		return &FinalityFilter{evmtypes.Unconfirmed}, nil
	default:
		return nil, fmt.Errorf("invalid finality confirmations filter value %v", filter.Confirmations)
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
