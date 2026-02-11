package topologyviz

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestASCIIMatrixHasClosingBorder(t *testing.T) {
	t.Parallel()

	summary := &TopologySummary{
		ConfigRef: "configs/workflow-gateway-don.toml",
		Topology:  ClassSingleDON,
		InfraType: "docker",
		DONs: []DONSummary{
			{
				Name:         "workflow",
				DONTypes:     []string{"workflow"},
				NodeCount:    4,
				Capabilities: []CapabilityPlacement{{RawFlag: "ocr3", BaseFlag: "ocr3"}},
			},
		},
	}

	matrix := RenderASCIICapabilityMatrix(summary)
	lines := strings.Split(strings.TrimSpace(matrix), "\n")
	require.GreaterOrEqual(t, len(lines), 5)
	// first line is section title, and the final line is the closing border
	require.True(t, strings.HasPrefix(lines[1], "+"))
	require.True(t, strings.HasPrefix(lines[len(lines)-1], "+"))
}
