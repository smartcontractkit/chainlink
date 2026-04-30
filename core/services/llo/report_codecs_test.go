package llo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
)

func Test_NewReportCodecs(t *testing.T) {
	c := NewReportCodecs(logger.Test(t), 1)

	assert.Contains(t, c, llotypes.ReportFormatJSON, "expected JSON to be supported")
	assert.Contains(t, c, llotypes.ReportFormatEVMPremiumLegacy, "expected EVMPremiumLegacy to be supported")
}
