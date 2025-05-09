package changeset

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/pointer"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
)

func TestRevokeJobSpecs(t *testing.T) {
	t.Parallel()

	const numBootstraps = 1
	const numOracles = 1

	env := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{
		ShouldDeployMCMS:      false,
		ShouldDeployLinkToken: false,
		NumNodes:              numBootstraps,
		NumBootstrapNodes:     numOracles,
		NodeLabels: []*ptypes.Label{
			{
				Key:   devenv.LabelProductKey,
				Value: pointer.To(utils.ProductLabel),
			},
			{
				Key:   devenv.LabelEnvironmentKey,
				Value: pointer.To("memory"),
			},
			{
				Key: utils.DonIdentifier(1, "don"),
			},
		},
	}).Environment

	uuidFromJobSpec := func(jobSpec string) string {
		matches := regexp.MustCompile(`externalJobID\s*=\s*'([a-f0-9-]+)'`).FindStringSubmatch(jobSpec)
		require.Len(t, matches, 2, "expected to find a UUID in the job spec")
		return matches[1]
	}

	// Create some jobs:
	sentJobs := sendTestLLOJobs(t, env, numOracles, numBootstraps, false)
	require.Len(t, sentJobs, 1)
	require.Len(t, sentJobs[0].Jobs, numBootstraps+numOracles)

	var oracleJobUUIDs, btJobUUIDs []string
	for _, job := range sentJobs[0].Jobs {
		if strings.Contains(job.Spec, "bootstrap") {
			btJobUUIDs = append(btJobUUIDs, uuidFromJobSpec(job.Spec))
		} else {
			oracleJobUUIDs = append(oracleJobUUIDs, uuidFromJobSpec(job.Spec))
		}
	}

	tests := []struct {
		name      string
		uuids     []string
		wantErr   string
		wantJobID string
	}{
		{
			name:      "Revoke an oracle job",
			uuids:     oracleJobUUIDs,
			wantJobID: oracleJobUUIDs[0], // we only have one
		},
		{
			name:    "Revoke the same job again",
			uuids:   oracleJobUUIDs,
			wantErr: "failed to revoke job",
		},
		{
			name:      "Revoke a bootstrap job",
			uuids:     btJobUUIDs,
			wantJobID: btJobUUIDs[0], // we only have one
		},
		{
			name:    "Revoke a non-existing job",
			uuids:   []string{"non-existing-job"},
			wantErr: "failed to find jobs for all provided UUIDs",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, out, err := changeset.ApplyChangesetsV2(t,
				env,
				[]changeset.ConfiguredChangeSet{
					changeset.Configure(CsRevokeJobSpecs{}, CsRevokeJobSpecsConfig{
						UUIDs: tc.uuids,
					}),
				})
			if tc.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErr, "unexpected error message")
				return
			}
			require.NoError(t, err)
			require.Len(t, out, 1)
			require.Len(t, out[0].Jobs, 1)
			require.Equal(t, tc.wantJobID, out[0].Jobs[0].JobID)
		})
	}
}
