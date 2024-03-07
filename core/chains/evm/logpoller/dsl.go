package logpoller

import (
	"github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

type EventTopicsByValueFilter struct {
	EventSigs []common.Hash
	Topics    [][]int
	Values    [][]string
}

func (f *EventTopicsByValueFilter) Accept(visitor commontypes.Visitor) {
	switch v := visitor.(type) {
	case *PgParserVisitor:
		v.VisitEventTopicsByValueFilter(f)
	}
}

type FinalityFilter struct {
	confs Confirmations
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
