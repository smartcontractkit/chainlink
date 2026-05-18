package llo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
)

func Test_NewCodecs(t *testing.T) {
	c := NewCodecs()

	assert.Contains(t, c, llotypes.ReportFormatJSON, "expected JSON to be supported")
	assert.Contains(t, c, llotypes.ReportFormatEVMPremiumLegacy, "expected EVMPremiumLegacy to be supported")
}
