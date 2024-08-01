package llo

import (
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

var _ PredecessorRetirementReportCache = &predecessorRetirementReportCache{}

type predecessorRetirementReportCache struct{}

// TODO: This ought to be DB-persisted
// https://smartcontract-it.atlassian.net/browse/MERC-3386
func NewPredecessorRetirementReportCache() PredecessorRetirementReportCache {
	return newPredecessorRetirementReportCache()
}

func newPredecessorRetirementReportCache() *predecessorRetirementReportCache {
	return &predecessorRetirementReportCache{}
}

func (c *predecessorRetirementReportCache) AttestedRetirementReport(predecessorConfigDigest types.ConfigDigest) ([]byte, error) {
	panic("TODO")
}

func (c *predecessorRetirementReportCache) CheckAttestedRetirementReport(predecessorConfigDigest types.ConfigDigest, attestedRetirementReport []byte) (RetirementReport, error) {
	panic("TODO")
}
