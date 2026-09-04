package runner_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/runner"
)

func TestResolveSpot_InvalidStrategyOverride(t *testing.T) {
	t.Parallel()

	_, err := runner.ResolveSpot(runner.SpotInput{
		EventName:        "pull_request",
		StrategyOverride: "bogus",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid spot strategy override")
}

func TestResolveSpot_InvalidDefaultStrategy(t *testing.T) {
	t.Parallel()

	_, err := runner.ResolveSpot(runner.SpotInput{
		EventName:       "pull_request",
		DefaultStrategy: "bogus",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid default spot strategy")
}

func TestResolveSpot(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		input            runner.SpotInput
		expectedSpot     string
		expectedFlag     string
		expectedEnab     bool
		expectedStrategy runner.SpotStrategy
		checkRelease     bool
		checkMQ          bool
	}{
		{
			name: "merge_group event forces on-demand",
			input: runner.SpotInput{
				EventName: "merge_group",
				Ref:       "refs/heads/gh-readonly-queue/develop/pr-123-abcdef",
				RefName:   "gh-readonly-queue/develop/pr-123-abcdef",
			},
			expectedSpot:     "false",
			expectedFlag:     "spot=false",
			expectedEnab:     false,
			expectedStrategy: runner.SpotDisabled,
			checkMQ:          true,
		},
		{
			name: "merge queue branch ref forces on-demand",
			input: runner.SpotInput{
				EventName: "push",
				Ref:       "refs/heads/gh-readonly-queue/main/pr-456",
				RefName:   "gh-readonly-queue/main/pr-456",
			},
			expectedSpot:     "false",
			expectedFlag:     "spot=false",
			expectedEnab:     false,
			expectedStrategy: runner.SpotDisabled,
			checkMQ:          true,
		},
		{
			name: "release event forces on-demand",
			input: runner.SpotInput{
				EventName: "release",
				Ref:       "refs/tags/v2.16.0",
				RefType:   "tag",
				RefName:   "v2.16.0",
			},
			expectedSpot:     "false",
			expectedFlag:     "spot=false",
			expectedEnab:     false,
			expectedStrategy: runner.SpotDisabled,
			checkRelease:     true,
		},
		{
			name: "tag ref forces on-demand",
			input: runner.SpotInput{
				EventName: "push",
				Ref:       "refs/tags/v2.16.0-rc1",
				RefType:   "tag",
				RefName:   "v2.16.0-rc1",
			},
			expectedSpot:     "false",
			expectedFlag:     "spot=false",
			expectedEnab:     false,
			expectedStrategy: runner.SpotDisabled,
			checkRelease:     true,
		},
		{
			name: "release branch push forces on-demand",
			input: runner.SpotInput{
				EventName: "push",
				Ref:       "refs/heads/release/2.57.1",
				RefName:   "release/2.57.1",
			},
			expectedSpot:     "false",
			expectedFlag:     "spot=false",
			expectedEnab:     false,
			expectedStrategy: runner.SpotDisabled,
			checkRelease:     true,
		},
		{
			name: "releases prefix branch forces on-demand",
			input: runner.SpotInput{
				EventName: "push",
				Ref:       "refs/heads/releases/v2.0.0",
				RefName:   "releases/v2.0.0",
			},
			expectedSpot:     "false",
			expectedFlag:     "spot=false",
			expectedEnab:     false,
			expectedStrategy: runner.SpotDisabled,
			checkRelease:     true,
		},
		{
			name: "hotfix branch forces on-demand",
			input: runner.SpotInput{
				EventName: "push",
				Ref:       "refs/heads/hotfix/core-bug",
				RefName:   "hotfix/core-bug",
			},
			expectedSpot:     "false",
			expectedFlag:     "spot=false",
			expectedEnab:     false,
			expectedStrategy: runner.SpotDisabled,
			checkRelease:     true,
		},
		{
			name: "PR targeting release branch forces on-demand",
			input: runner.SpotInput{
				EventName: "pull_request",
				BaseRef:   "release/2.57.1",
				HeadRef:   "fix-something",
			},
			expectedSpot:     "false",
			expectedFlag:     "spot=false",
			expectedEnab:     false,
			expectedStrategy: runner.SpotDisabled,
			checkRelease:     true,
		},
		{
			name: "push to develop uses capacity-optimized spot",
			input: runner.SpotInput{
				EventName: "push",
				Ref:       "refs/heads/develop",
				RefName:   "develop",
			},
			expectedSpot:     "co",
			expectedFlag:     "spot=co",
			expectedEnab:     true,
			expectedStrategy: runner.SpotCapacityOptimized,
		},
		{
			name: "push to main uses capacity-optimized spot",
			input: runner.SpotInput{
				EventName: "push",
				Ref:       "refs/heads/main",
				RefName:   "main",
			},
			expectedSpot:     "co",
			expectedFlag:     "spot=co",
			expectedEnab:     true,
			expectedStrategy: runner.SpotCapacityOptimized,
		},
		{
			name: "standard PR uses price-capacity-optimized spot",
			input: runner.SpotInput{
				EventName: "pull_request",
				BaseRef:   "develop",
				HeadRef:   "feature/DX-5101",
			},
			expectedSpot:     "pco",
			expectedFlag:     "spot=pco",
			expectedEnab:     true,
			expectedStrategy: runner.SpotPriceCapacityOptimized,
		},
		{
			name: "scheduled run on weekend uses price-capacity-optimized spot",
			input: runner.SpotInput{
				EventName: "schedule",
				Ref:       "refs/heads/develop",
			},
			expectedSpot:     "pco",
			expectedFlag:     "spot=pco",
			expectedEnab:     true,
			expectedStrategy: runner.SpotPriceCapacityOptimized,
		},
		{
			name: "force on demand override disables spot",
			input: runner.SpotInput{
				EventName:     "pull_request",
				BaseRef:       "develop",
				HeadRef:       "feature/some-feature",
				ForceOnDemand: true,
			},
			expectedSpot:     "false",
			expectedFlag:     "spot=false",
			expectedEnab:     false,
			expectedStrategy: runner.SpotDisabled,
		},
		{
			name: "explicit strategy override",
			input: runner.SpotInput{
				EventName:        "pull_request",
				BaseRef:          "develop",
				HeadRef:          "feature/some-feature",
				StrategyOverride: runner.SpotLowestPrice,
			},
			expectedSpot:     "lowest-price",
			expectedFlag:     "spot=lowest-price",
			expectedEnab:     true,
			expectedStrategy: runner.SpotLowestPrice,
		},
		{
			name: "custom default strategy",
			input: runner.SpotInput{
				EventName:       "pull_request",
				BaseRef:         "develop",
				HeadRef:         "feature/some-feature",
				DefaultStrategy: runner.SpotCapacityOptimized,
			},
			expectedSpot:     "co",
			expectedFlag:     "spot=co",
			expectedEnab:     true,
			expectedStrategy: runner.SpotCapacityOptimized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			res, err := runner.ResolveSpot(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedSpot, res.Spot)
			assert.Equal(t, tt.expectedFlag, res.SpotFlag)
			assert.Equal(t, tt.expectedEnab, res.Enabled)
			assert.Equal(t, tt.expectedStrategy, res.Strategy)
			assert.NotEmpty(t, res.Reason)

			if tt.checkRelease {
				assert.True(t, res.IsRelease, "expected is_release to be true")
			}
			if tt.checkMQ {
				assert.True(t, res.IsMergeQueue, "expected is_merge_queue to be true")
			}
		})
	}
}
