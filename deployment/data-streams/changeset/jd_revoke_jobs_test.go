package changeset

import (
	"regexp"
	"strconv"
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
	const numOracles = 2

	env := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{
		ShouldDeployMCMS:      false,
		ShouldDeployLinkToken: false,
		NumBootstrapNodes:     numBootstraps,
		NumNodes:              numOracles,
		NodeLabels: []*ptypes.Label{
			{
				Key:   devenv.LabelProductKey,
				Value: pointer.To(utils.ProductLabel),
			},
			{
				Key:   devenv.LabelEnvironmentKey,
				Value: pointer.To(testutil.TestDON.Env),
			},
			{
				Key: utils.DonIdentifier(testutil.TestDON.ID, testutil.TestDON.Name),
			},
		},
		CustomDBSetup: []string{
			// Seed the database with the list of bridges we're using.
			`INSERT INTO bridge_types (name, url, confirmations, incoming_token_hash, salt, outgoing_token, created_at, updated_at)
				VALUES ('bridge-api1', 'http://url', 0, '', '', '', now(), now());`,
			`INSERT INTO bridge_types (name, url, confirmations, incoming_token_hash, salt, outgoing_token, created_at, updated_at)
				VALUES ('bridge-api2', 'http://url', 0, '', '', '', now(), now());`,
			`INSERT INTO bridge_types (name, url, confirmations, incoming_token_hash, salt, outgoing_token, created_at, updated_at)
				VALUES ('bridge-api3', 'http://url', 0, '', '', '', now(), now());`,
			`INSERT INTO bridge_types (name, url, confirmations, incoming_token_hash, salt, outgoing_token, created_at, updated_at)
				VALUES ('bridge-api4', 'http://url', 0, '', '', '', now(), now());`,
		},
	}).Environment

	uuidFromJobSpec := func(jobSpec string) string {
		matches := regexp.MustCompile(`externalJobID\s*=\s*'([a-f0-9-]+)'`).FindStringSubmatch(jobSpec)
		require.Len(t, matches, 2, "expected to find a UUID in the job spec")
		return matches[1]
	}

	// Create some jobs:
	sentLLOJobs := sendTestLLOJobs(t, env, numOracles, numBootstraps, false)
	require.Len(t, sentLLOJobs, 1)
	require.Len(t, sentLLOJobs[0].Jobs, numBootstraps+numOracles)

	var oracleJobUUIDs, btJobUUIDs []string
	for _, job := range sentLLOJobs[0].Jobs {
		if strings.Contains(job.Spec, "bootstrap") {
			btJobUUIDs = append(btJobUUIDs, uuidFromJobSpec(job.Spec))
		} else if strings.Contains(job.Spec, "offchainreporting2") {
			oracleJobUUIDs = append(oracleJobUUIDs, uuidFromJobSpec(job.Spec))
		}
	}

	// Create some stream jobs:
	sentStreamJobs := sendTestStreamJobs(t, env, numOracles, false)
	require.Len(t, sentStreamJobs, 1)
	require.Len(t, sentStreamJobs[0].Jobs, numOracles)

	streamIDFromJobSpec := func(jobSpec string) string {
		matches := regexp.MustCompile(`streamID\s*=\s*([0-9]+)`).FindStringSubmatch(jobSpec)
		require.Len(t, matches, 2, "expected to find a stream ID in the job spec")
		return matches[1]
	}

	var streamIDs []uint32
	streamIDsToJobIDs := make(map[uint32][]string)
	for _, job := range sentStreamJobs[0].Jobs {
		s, e := strconv.ParseUint(streamIDFromJobSpec(job.Spec), 10, 32)
		require.NoError(t, e)
		streamIDs = append(streamIDs, uint32(s))
		streamIDsToJobIDs[uint32(s)] = append(streamIDsToJobIDs[uint32(s)], uuidFromJobSpec(job.Spec))
	}

	tests := []struct {
		name        string
		uuids       []string
		streamIDs   []uint32
		wantErr     string
		wantJobIDs  []string
		wantNumJobs int
	}{
		{
			name:        "Revoke an oracle job",
			uuids:       oracleJobUUIDs,
			wantJobIDs:  oracleJobUUIDs,
			wantNumJobs: numOracles,
		},
		{
			name:        "Revoke the same job again",
			uuids:       oracleJobUUIDs,
			wantNumJobs: numOracles,
			wantErr:     "failed to revoke job",
		},
		{
			name:        "Revoke a bootstrap job",
			uuids:       btJobUUIDs,
			wantJobIDs:  btJobUUIDs,
			wantNumJobs: numBootstraps,
		},
		{
			name:    "Revoke a non-existing job",
			uuids:   []string{"non-existing-job"},
			wantErr: "failed to find jobs for all provided UUIDs",
		},
		{
			name:        "Revoke a stream job",
			streamIDs:   []uint32{streamIDs[0]},
			wantNumJobs: numOracles,
			wantJobIDs:  streamIDsToJobIDs[streamIDs[0]],
		},
		{
			name:        "Fail when both stream ids and uuids are provided",
			uuids:       oracleJobUUIDs,
			streamIDs:   streamIDs,
			wantNumJobs: numOracles,
			wantErr:     "either job ids or stream ids are required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, out, err := changeset.ApplyChangesetsV2(t,
				env,
				[]changeset.ConfiguredChangeSet{
					changeset.Configure(CsRevokeJobSpecs{}, CsRevokeJobSpecsConfig{
						UUIDs:     tc.uuids,
						StreamIDs: tc.streamIDs,
					}),
				})
			if tc.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErr, "unexpected error message")
				return
			}
			require.NoError(t, err)
			require.Len(t, out, 1)
			require.Len(t, out[0].Jobs, tc.wantNumJobs)
			for _, wantedJobID := range tc.wantJobIDs {
				found := false
				for _, job := range out[0].Jobs {
					if job.JobID == wantedJobID {
						found = true
						break
					}
				}
				require.True(t, found, "expected to find job %s in the output", wantedJobID)
			}
		})
	}
}
