package logpoller

import (
	"math/big"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
)

// PgParserVisitor is a visitor that builds a postgres query and arguments from a commontypes.QueryFilter
type PgParserVisitor struct {
}

var _ commontypes.Visitor = (*PgParserVisitor)(nil)

func NewPgParser(evmChainID *big.Int) *PgParserVisitor {
	return &PgParserVisitor{}
}

func (v *PgParserVisitor) VisitAndFilter(node commontypes.AndFilter) {
	//TODO implement me
	panic("implement me")
}

func (v *PgParserVisitor) VisitAddressFilter(node commontypes.AddressFilter) {
	//TODO implement me
	panic("implement me")
}

func (v *PgParserVisitor) VisitKeysFilter(node commontypes.KeysFilter) {
	//TODO implement me
	panic("implement me")
}

// VisitKeysByValueFilter is unused chain agnostic version of VisitEventTopicsByValueFilter.
func (v *PgParserVisitor) VisitKeysByValueFilter(node commontypes.KeysByValueFilter) {
	return
}

func (v *PgParserVisitor) VisitBlockFilter(node commontypes.BlockFilter) {
	//TODO implement me
	panic("implement me")
}

// VisitConfirmationFilter is unused chain agnostic version of VisitFinalityFilter.
func (v *PgParserVisitor) VisitConfirmationFilter(node commontypes.ConfirmationFilter) {
	return
}

func (v *PgParserVisitor) VisitTimestampFilter(node commontypes.TimestampFilter) {
	//TODO implement me
	panic("implement me")
}

func (v *PgParserVisitor) VisitTxHashFilter(node commontypes.TxHashFilter) {
	//TODO implement me
	panic("implement me")
}

func (v *PgParserVisitor) VisitEventTopicsByValueFilter(filter *EventTopicsByValueFilter) {
	//TODO implement me
	panic("implement me")
}

func (v *PgParserVisitor) VisitFinalityFilter(filter *FinalityFilter) {
	//TODO implement me
	panic("implement me")
}

func (v *PgParserVisitor) VisitChainIdFilter(filter *ChainIdFilter) {
	//TODO implement me
	panic("implement me")
}
