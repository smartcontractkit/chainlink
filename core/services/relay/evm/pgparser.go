package evm

import (
	"math/big"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
)

// PgDSLParser is a visitor that builds a postgres query and arguments from a commontypes.QueryFilter
type PgDSLParser struct {
}

var _ commontypes.Visitor = (*PgDSLParser)(nil)

func NewPgParser(evmChainID *big.Int) *PgDSLParser {
	return &PgDSLParser{}
}

func (v *PgDSLParser) VisitAndFilter(node commontypes.AndFilter) {
	//TODO implement me
	panic("implement me")
}

func (v *PgDSLParser) VisitAddressFilter(node commontypes.AddressFilter) {
	//TODO implement me
	panic("implement me")
}

func (v *PgDSLParser) VisitKeysFilter(node commontypes.KeysFilter) {
	//TODO implement me
	panic("implement me")
}

// VisitKeysByValueFilter is unused chain agnostic version of VisitEventTopicsByValueFilter.
func (v *PgDSLParser) VisitKeysByValueFilter(node commontypes.KeysByValueFilter) {
	return
}

func (v *PgDSLParser) VisitBlockFilter(node commontypes.BlockFilter) {
	//TODO implement me
	panic("implement me")
}

// VisitConfirmationFilter is unused chain agnostic version of VisitFinalityFilter.
func (v *PgDSLParser) VisitConfirmationFilter(node commontypes.ConfirmationFilter) {
	return
}

func (v *PgDSLParser) VisitTimestampFilter(node commontypes.TimestampFilter) {
	//TODO implement me
	panic("implement me")
}

func (v *PgDSLParser) VisitTxHashFilter(node commontypes.TxHashFilter) {
	//TODO implement me
	panic("implement me")
}

func (v *PgDSLParser) VisitEventTopicsByValueFilter(filter *EventTopicsByValueFilter) {
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
