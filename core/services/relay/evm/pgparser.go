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

func (v *PgDSLParser) ComparerPrimitive(_ query.ComparerPrimitive) {}

func (v *PgDSLParser) BlockPrimitive(_ query.BlockPrimitive) {}

func (v *PgDSLParser) ConfirmationPrimitive(_ query.ConfirmationsPrimitive) {}

func (v *PgDSLParser) TimestampPrimitive(_ query.TimestampPrimitive) {}

func (v *PgDSLParser) TxHashPrimitives(_ query.TxHashPrimitive) {}

func (v *PgDSLParser) VisitEventTopicsByValueFilter(_ *EventByTopicFilter) {}

func (v *PgDSLParser) VisitEventByWordFilter(_ *EventByWordFilter) {}

func (v *PgDSLParser) VisitEventFilter(_ *EventFilter) {}

func (v *PgDSLParser) VisitFinalityFilter(_ *FinalityFilter) {}
