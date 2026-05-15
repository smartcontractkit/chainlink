package llo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
)

func Test_NewReportCodecs(t *testing.T) {
	c := NewReportCodecs(logger.Test(t), 1)

	_, ok := c[llotypes.ReportFormatJSON]
	assert.True(t, ok, "expected JSON to be supported")
	_, ok = c[llotypes.ReportFormatEVMPremiumLegacy]
	assert.True(t, ok, "expected EVMPremiumLegacy to be supported")
	_, ok = c[llotypes.ReportFormatHistoryBackfill]
	assert.True(t, ok, "expected HistoryBackfill meta-format to be supported for definition verification")
}
