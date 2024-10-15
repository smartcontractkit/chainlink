//go:build integration

package ocr2_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func TestIntegration_OCR2_plugins(t *testing.T) {
	t.Setenv(string(env.MedianPlugin.Cmd), "chainlink-feeds")
	testutils.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/BCF-3417")
	testIntegration_OCR2(t)
}
