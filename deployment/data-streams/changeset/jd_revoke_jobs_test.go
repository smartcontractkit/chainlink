package changeset

import (
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/pointer"
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
				Key:   utils.LabelProduct,
				Value: pointer.To(utils.ProductLabel),
			},
			{
				Key:   utils.LabelEnvironment,
				Value: pointer.To("memory"),
			},
			{
				Key: utils.DonIdentifier(1, "don"),
			},
		},
	}).Environment

	// Create some jobs:
	out := sendTestLLOJobs(t, env, numOracles, numBootstraps)
	require.Len(t, out, 1)
	require.Len(t, out[0].Jobs, numBootstraps+numOracles)

	var oracleJobID, btJobID string
	for _, job := range out[0].Jobs {
		if strings.Contains(job.Spec, "bootstrap") {
			btJobID = job.JobID
		} else {
			oracleJobID = job.JobID
		}
	}

	tests := []struct {
		name      string
		jobID     string
		wantErr   string
		wantJobID string
	}{
		{
			name:      "Revoke an oracle job",
			jobID:     oracleJobID,
			wantJobID: oracleJobID,
		},
		{
			name:      "Revoke a bootstrap job",
			jobID:     btJobID,
			wantJobID: btJobID,
		},
		{
			name:    "Revoke a non-existing job",
			jobID:   "non-existing-job",
			wantErr: "failed to revoke job",
		},
		{
			name:    "Revoke the same job again",
			jobID:   oracleJobID,
			wantErr: "failed to revoke job",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, out, err := changeset.ApplyChangesetsV2(t,
				env,
				[]changeset.ConfiguredChangeSet{
					changeset.Configure(CsRevokeJobSpecs{}, CsRevokeJobSpecsConfig{
						JobID: tc.jobID,
					}),
				})
			if tc.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErr)
				return
			}
			require.NoError(t, err)
			require.Len(t, out, 1)
			require.Len(t, out[0].Jobs, 1)
			require.Equal(t, tc.wantJobID, out[0].Jobs[0].JobID)

		})

	}
}
