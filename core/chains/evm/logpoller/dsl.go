package logpoller

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

type EventTopicsByValueFilter struct {
	EventSigs []common.Hash
	Topics    [][]int
	Values    [][]string
}

func NewEventTopicsByValueFilter(filter *commontypes.KeysByValueFilter, eventIndexBindings evm.EventIndexBindings) (*EventTopicsByValueFilter, error) {
	var searchEventTopicsByValueFilter *EventTopicsByValueFilter
	for i, key := range filter.Keys {
		eventSig, _, index, err := eventIndexBindings.Get(key)
		if err != nil {
			return nil, err
		}
		searchEventTopicsByValueFilter.EventSigs = append(searchEventTopicsByValueFilter.EventSigs, eventSig)
		searchEventTopicsByValueFilter.Topics = append(searchEventTopicsByValueFilter.Topics, []int{index})
		searchEventTopicsByValueFilter.Values = append(searchEventTopicsByValueFilter.Values, filter.Values[i])
	}
	return searchEventTopicsByValueFilter, nil
}

func (f *EventTopicsByValueFilter) Accept(visitor commontypes.Visitor) {
	switch v := visitor.(type) {
	case *PgParserVisitor:
		v.VisitEventTopicsByValueFilter(f)
	}
}

type FinalityFilter struct {
	Confs Confirmations
}

func NewFinalityFilter(filter *commontypes.ConfirmationFilter) (*FinalityFilter, error) {
	switch filter.Confirmations {
	case commontypes.Finalized:
		return &FinalityFilter{Finalized}, nil
	case commontypes.Unconfirmed:
		return &FinalityFilter{Unconfirmed}, nil
	default:
		return nil, fmt.Errorf("invalid finality confirmations filter value %v", filter.Confirmations)
	}
}

func (f *FinalityFilter) Accept(visitor commontypes.Visitor) {
	switch v := visitor.(type) {
	case *PgParserVisitor:
		v.VisitFinalityFilter(f)
	}
}

type ChainIdFilter struct {
	chainId *ubig.Big
}

func (f *ChainIdFilter) accept(visitor commontypes.Visitor) {
	switch v := visitor.(type) {
	case *PgParserVisitor:
		v.VisitChainIdFilter(f)
	}
}

func AppendedNewFilter(root *commontypes.AndFilter, other ...commontypes.QueryFilter) *commontypes.AndFilter {
	filters := make([]commontypes.QueryFilter, 0, len(root.Filters)+len(other))
	filters = append(filters, root.Filters...)
	filters = append(filters, other...)
	return commontypes.NewAndFilter(filters...)
}
