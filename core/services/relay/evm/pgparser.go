package evm

import (
	"math/big"

	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
)

// PgDSLParser is a visitor that builds a postgres query and arguments from a query.KeyFilter
type PgDSLParser struct {
	//TODO implement psql parser
}

var _ primitives.Visitor = (*PgDSLParser)(nil)

func NewPgParser(evmChainID *big.Int) *PgDSLParser {
	return &PgDSLParser{}
}

func (v *PgDSLParser) Comparator(_ primitives.Comparator) {}

func (v *PgDSLParser) Block(_ primitives.Block) {}

func (v *PgDSLParser) Confirmations(_ primitives.Confirmations) {}

func (v *PgDSLParser) Timestamp(_ primitives.Timestamp) {}

func (v *PgDSLParser) TxHash(_ primitives.TxHash) {}

func (v *PgDSLParser) VisitEventTopicsByValueFilter(_ *EventByTopicFilter) {}

func (v *PgDSLParser) VisitEventByWordFilter(_ *EventByWordFilter) {}

func (v *PgDSLParser) VisitEventBySigFilter(_ *EventBySigFilter) {}

func (v *PgDSLParser) VisitFinalityFilter(_ *FinalityFilter) {}
