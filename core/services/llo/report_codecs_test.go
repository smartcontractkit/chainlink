package llo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_NewReportCodecs(t *testing.T) {
	c := NewReportCodecs(logger.TestLogger(t), 1)

	assert.Contains(t, c, llotypes.ReportFormatJSON, "expected JSON to be supported")
	assert.Contains(t, c, llotypes.ReportFormatEVMPremiumLegacy, "expected EVMPremiumLegacy to be supported")
}
