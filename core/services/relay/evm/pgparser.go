package evm

import (
	"math/big"

	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
)

// PgDSLParser is a visitor that builds a postgres query and arguments from a commontypes.QueryFilter
type PgDSLParser struct {
}

var _ query.Visitor = (*PgDSLParser)(nil)

func NewPgParser(evmChainID *big.Int) *PgDSLParser {
	return &PgDSLParser{}
}

// TODO remove from common
func (v *PgDSLParser) AddressPrimitive(primitive query.AddressPrimitive) {
	//TODO implement me
	panic("implement me")
}

func (v *PgDSLParser) BlockPrimitive(primitive query.BlockPrimitive) {
	//TODO implement me
	panic("implement me")
}

func (v *PgDSLParser) ConfirmationPrimitive(primitive query.ConfirmationsPrimitive) {
	//TODO implement me
	panic("implement me")
}

func (v *PgDSLParser) TimestampPrimitive(primitive query.TimestampPrimitive) {
	//TODO implement me
	panic("implement me")
}

func (v *PgDSLParser) TxHashPrimitives(primitive query.TxHashPrimitive) {
	//TODO implement me
	panic("implement me")
}

func (v *PgDSLParser) VisitEventTopicsByValueFilter(filter *EventByIndexFilter) {
	//TODO implement me
	panic("implement me")
}

func (v *PgDSLParser) VisitEventFilter(filter *EventFilter) {
	//TODO implement me
	panic("implement me")
}

func (v *PgDSLParser) VisitFinalityFilter(filter *FinalityFilter) {
	//TODO implement me
	panic("implement me")
}

func (v *PgDSLParser) VisitChainIdFilter(filter *ChainIdFilter) {
	//TODO implement me
	panic("implement me")
}
