package evm

import (
	"math/big"

	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
)

// PgDSLParser is a visitor that builds a postgres query and arguments from a commontypes.QueryFilter
type PgDSLParser struct {
	//TODO implement psql parser
}

var _ query.Visitor = (*PgDSLParser)(nil)

func NewPgParser(evmChainID *big.Int) *PgDSLParser {
	return &PgDSLParser{}
}

func (v *PgDSLParser) ComparerPrimitive(_ query.ComparerPrimitive) {
	return
}

func (v *PgDSLParser) BlockPrimitive(_ query.BlockPrimitive) {
	return
}

func (v *PgDSLParser) ConfirmationPrimitive(_ query.ConfirmationsPrimitive) {
	return
}

func (v *PgDSLParser) TimestampPrimitive(_ query.TimestampPrimitive) {
	return
}

func (v *PgDSLParser) TxHashPrimitives(_ query.TxHashPrimitive) {
	return
}

func (v *PgDSLParser) VisitEventTopicsByValueFilter(_ *EventByIndexFilter) {
	return
}

func (v *PgDSLParser) VisitEventByWordFilter(_ *EventByWordFilter) {
	return
}

func (v *PgDSLParser) VisitEventFilter(_ *EventFilter) {
	return
}

func (v *PgDSLParser) VisitFinalityFilter(_ *FinalityFilter) {
	return
}

func (v *PgDSLParser) VisitChainIdFilter(_ *ChainIdFilter) {
	return
}
