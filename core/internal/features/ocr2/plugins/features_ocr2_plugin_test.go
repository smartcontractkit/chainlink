//go:build integration

// This package exists to separate a long-running test which sets environment variables, and so cannot be run in parallel.
package plugins

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/internal/features/ocr2"
)

func TestIntegration_OCR2_plugins(t *testing.T) {
	t.Setenv(string(env.MedianPlugin.Cmd), "chainlink-feeds")
	ocr2.RunTestIntegrationOCR2(t)
}
