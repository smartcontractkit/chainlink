package gating

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEvaluate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		in       Inputs
		expected Decisions
	}{
		{
			name: "workflow dispatch runs everything",
			in:   Inputs{EventName: "workflow_dispatch"},
			expected: Decisions{
				CREShouldRun:      true,
				CREWithRegression: true,
				CRERunMixedEnv:    true,
				CCIPShouldRun:     true,
				BuildCoreImage:    true,
				BuildPluginsImage: true,
			},
		},
		{
			name:     "push to develop without changes runs regression and mixed-env only",
			in:       Inputs{EventName: "push", RefName: "develop"},
			expected: Decisions{CREWithRegression: true, CRERunMixedEnv: true},
		},
		{
			name: "push to feature branch with CRE changes only runs CRE",
			in:   Inputs{EventName: "push", RefName: "feature/x", CREChanges: true},
			expected: Decisions{
				CREShouldRun:      true,
				BuildCoreImage:    true,
				BuildPluginsImage: true,
			},
		},
		{
			name: "push of tag without changes runs CRE and CCIP",
			in:   Inputs{EventName: "push", RefName: "v2.58.0", RefType: "tag"},
			expected: Decisions{
				CREShouldRun:      true,
				CCIPShouldRun:     true,
				BuildCoreImage:    true,
				BuildPluginsImage: true,
			},
		},
		{
			name: "pull request without changes or labels only runs regression and mixed-env",
			in:   Inputs{EventName: "pull_request"},
			expected: Decisions{
				CREWithRegression: true,
				CRERunMixedEnv:    true,
			},
		},
		{
			name: "pull request with CRE changes runs CRE, regression, and mixed-env",
			in:   Inputs{EventName: "pull_request", CREChanges: true},
			expected: Decisions{
				CREShouldRun:      true,
				CREWithRegression: true,
				CRERunMixedEnv:    true,
				BuildCoreImage:    true,
				BuildPluginsImage: true,
			},
		},
		{
			name: "pull request with run-e2e label runs CRE and builds plugins",
			in:   Inputs{EventName: "pull_request", RunE2ELabel: true},
			expected: Decisions{
				CREShouldRun:      true,
				CREWithRegression: true,
				CRERunMixedEnv:    true,
				BuildCoreImage:    true,
				BuildPluginsImage: true,
			},
		},
		{
			name: "pull request with skip-regression label skips regression",
			in:   Inputs{EventName: "pull_request", SkipRegressionLabel: true},
			expected: Decisions{
				CRERunMixedEnv: true,
			},
		},
		{
			name: "pull request with skip-mixed-env label skips mixed-env",
			in:   Inputs{EventName: "pull_request", SkipMixedEnvLabel: true},
			expected: Decisions{
				CREWithRegression: true,
			},
		},
		{
			name: "merge_group with CCIP changes only runs CCIP",
			in:   Inputs{EventName: "merge_group", CCIPChanges: true},
			expected: Decisions{
				CCIPShouldRun:  true,
				BuildCoreImage: true,
			},
		},
		{
			name:     "merge_group without changes runs nothing",
			in:       Inputs{EventName: "merge_group"},
			expected: Decisions{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.expected, Evaluate(tt.in))
		})
	}
}

func TestOutputVars(t *testing.T) {
	t.Parallel()

	vars := Decisions{
		CREShouldRun:      true,
		CREWithRegression: false,
		CRERunMixedEnv:    true,
		CCIPShouldRun:     false,
		BuildCoreImage:    true,
		BuildPluginsImage: true,
	}.OutputVars()

	require.Equal(t, map[string]string{
		"cre-should-run":      "true",
		"cre-with-regression": "false",
		"cre-run-mixed-env":   "true",
		"ccip-should-run":     "false",
		"build-core-image":    "true",
		"build-plugins-image": "true",
	}, vars)
}

func TestSummaryTable(t *testing.T) {
	t.Parallel()

	decisions := Evaluate(Inputs{EventName: "pull_request", CREChanges: true})
	table := decisions.SummaryTable(Inputs{EventName: "pull_request", CREChanges: true})

	require.Contains(t, table, "### Integration Test Gating Decisions")
	require.Contains(t, table, "| **Build Core Image** | `true` |")
	require.Contains(t, table, "| **Core CRE Smoke Tests** | `true` | Event: `pull_request`")
	require.Contains(t, table, "| **CCIP v1.6 E2E Tests** | `false` |")
}
