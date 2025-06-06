package memory_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/freeport"

	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestJobClientProposeJob(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	chains, _ := memory.NewMemoryChains(t, 1, 1)
	ports := freeport.GetN(t, 1)
	c := memory.NewNodeConfig{
		Port:           ports[0],
		Chains:         chains,
		Solchains:      nil,
		Aptoschains:    nil,
		LogLevel:       zapcore.DebugLevel,
		Bootstrap:      false,
		RegistryConfig: deployment.CapabilityRegistryConfig{},
		CustomDBSetup:  nil,
	}
	testNode := memory.NewNode(t, c)

	// Set up the JobClient with a mock node
	nodeID := "node-1"
	nodes := map[string]memory.Node{
		nodeID: *testNode,
	}
	jobClient := memory.NewMemoryJobClient(nodes)

	type testCase struct {
		name      string
		req       *jobv1.ProposeJobRequest
		checkErr  func(t *testing.T, err error)
		checkResp func(t *testing.T, resp *jobv1.ProposeJobResponse)
	}
	cases := []testCase{
		{
			name: "valid request",
			req: &jobv1.ProposeJobRequest{
				NodeId: "node-1",
				Spec:   testJobProposalTOML(t, "f1ac5211-ab79-4c31-ba1c-0997b72db466"),
			},
			checkResp: func(t *testing.T, resp *jobv1.ProposeJobResponse) {
				assert.NotNil(t, resp)
				assert.Equal(t, int64(1), resp.Proposal.Revision)
				assert.Equal(t, jobv1.ProposalStatus_PROPOSAL_STATUS_APPROVED, resp.Proposal.Status)
				assert.Equal(t, jobv1.ProposalDeliveryStatus_PROPOSAL_DELIVERY_STATUS_DELIVERED, resp.Proposal.DeliveryStatus)
				assert.Equal(t, "f1ac5211-ab79-4c31-ba1c-0997b72db466", resp.Proposal.JobId)
				assert.Equal(t, testJobProposalTOML(t, "f1ac5211-ab79-4c31-ba1c-0997b72db466"), resp.Proposal.Spec)
			},
		},
		{
			name: "idempotent request bumps version",
			req: &jobv1.ProposeJobRequest{
				NodeId: "node-1",
				Spec:   testJobProposalTOML(t, "f1ac5211-ab79-4c31-ba1c-0997b72db466"),
			},
			// the feeds service doesn't allow duplicate job names
			checkResp: func(t *testing.T, resp *jobv1.ProposeJobResponse) {
				assert.NotNil(t, resp)
				assert.Equal(t, int64(2), resp.Proposal.Revision)
				assert.Equal(t, jobv1.ProposalStatus_PROPOSAL_STATUS_APPROVED, resp.Proposal.Status)
				assert.Equal(t, jobv1.ProposalDeliveryStatus_PROPOSAL_DELIVERY_STATUS_DELIVERED, resp.Proposal.DeliveryStatus)
				assert.Equal(t, "f1ac5211-ab79-4c31-ba1c-0997b72db466", resp.Proposal.JobId)
				assert.Equal(t, testJobProposalTOML(t, "f1ac5211-ab79-4c31-ba1c-0997b72db466"), resp.Proposal.Spec)
			},
		},
		{
			name: "another request",
			req: &jobv1.ProposeJobRequest{
				NodeId: "node-1",
				Spec:   testJobProposalTOML(t, "11115211-ab79-4c31-ba1c-0997b72aaaaa"),
			},
			checkResp: func(t *testing.T, resp *jobv1.ProposeJobResponse) {
				assert.NotNil(t, resp)
				assert.Equal(t, int64(1), resp.Proposal.Revision)
				assert.Equal(t, jobv1.ProposalStatus_PROPOSAL_STATUS_APPROVED, resp.Proposal.Status)
				assert.Equal(t, jobv1.ProposalDeliveryStatus_PROPOSAL_DELIVERY_STATUS_DELIVERED, resp.Proposal.DeliveryStatus)
				assert.Equal(t, "11115211-ab79-4c31-ba1c-0997b72aaaaa", resp.Proposal.JobId)
				assert.Equal(t, testJobProposalTOML(t, "11115211-ab79-4c31-ba1c-0997b72aaaaa"), resp.Proposal.Spec)
			},
		},
		{
			name: "node does not exist",
			req: &jobv1.ProposeJobRequest{
				NodeId: "node-2",
				Spec:   testJobProposalTOML(t, "f1ac5211-ab79-4c31-ba1c-0997b72db466"),
			},
			checkErr: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "node not found")
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Call the ProposeJob method
			resp, err := jobClient.ProposeJob(ctx, c.req)
			if c.checkErr != nil {
				c.checkErr(t, err)
				return
			}
			require.NoError(t, err)
			c.checkResp(t, resp)
		})
	}
}

func TestJobClientJobAPI(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	chains, _ := memory.NewMemoryChains(t, 1, 1)
	ports := freeport.GetN(t, 1)
	c := memory.NewNodeConfig{
		Port:           ports[0],
		Chains:         chains,
		Solchains:      nil,
		Aptoschains:    nil,
		LogLevel:       zapcore.DebugLevel,
		Bootstrap:      false,
		RegistryConfig: deployment.CapabilityRegistryConfig{},
		CustomDBSetup:  nil,
	}
	testNode := memory.NewNode(t, c)

	// Set up the JobClient with a mock node
	nodeID := "node-1"
	externalJobID := "f1ac5211-ab79-4c31-ba1c-0997b72db466"

	jobSpecToml := testJobProposalTOML(t, externalJobID)
	nodes := map[string]memory.Node{
		nodeID: *testNode,
	}
	jobClient := memory.NewMemoryJobClient(nodes)

	// Create a mock request
	req := &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec:   jobSpecToml,
		Labels: []*ptypes.Label{
			{
				Key:   "label-key",
				Value: ptr("label-value"),
			},
		},
	}

	// Call the ProposeJob method
	resp, err := jobClient.ProposeJob(ctx, req)

	// Validate the response
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, jobv1.ProposalStatus_PROPOSAL_STATUS_APPROVED, resp.Proposal.Status)
	assert.Equal(t, jobv1.ProposalDeliveryStatus_PROPOSAL_DELIVERY_STATUS_DELIVERED, resp.Proposal.DeliveryStatus)
	assert.Equal(t, jobSpecToml, resp.Proposal.Spec)
	assert.Equal(t, externalJobID, resp.Proposal.JobId)

	expectedProposalID := resp.Proposal.Id
	expectedProposal := resp.Proposal

	t.Run("GetJob", func(t *testing.T) {
		t.Run("existing job", func(t *testing.T) {
			// Create a mock request
			getReq := &jobv1.GetJobRequest{
				IdOneof: &jobv1.GetJobRequest_Id{Id: externalJobID},
			}

			getResp, err := jobClient.GetJob(ctx, getReq)
			require.NoError(t, err)
			assert.NotNil(t, getResp)
			assert.Equal(t, externalJobID, getResp.Job.Id)
		})

		t.Run("non-existing job", func(t *testing.T) {
			// Create a mock request
			getReq := &jobv1.GetJobRequest{
				IdOneof: &jobv1.GetJobRequest_Id{Id: "non-existing-job"},
			}

			getResp, err := jobClient.GetJob(ctx, getReq)
			require.Error(t, err)
			assert.Nil(t, getResp)
		})
	})

	t.Run("ListJobs", func(t *testing.T) {
		type listCase struct {
			name      string
			req       *jobv1.ListJobsRequest
			checkErr  func(t *testing.T, err error)
			checkResp func(t *testing.T, resp *jobv1.ListJobsResponse)
		}
		cases := []listCase{
			{
				name: "no filters",
				req:  &jobv1.ListJobsRequest{},
				checkResp: func(t *testing.T, resp *jobv1.ListJobsResponse) {
					assert.NotNil(t, resp)
					assert.Len(t, resp.Jobs, 1)
					assert.Equal(t, externalJobID, resp.Jobs[0].Id)
				},
			},
			{
				name: "with id filter",
				req: &jobv1.ListJobsRequest{
					Filter: &jobv1.ListJobsRequest_Filter{
						Ids: []string{externalJobID},
					},
				},
				checkResp: func(t *testing.T, resp *jobv1.ListJobsResponse) {
					assert.NotNil(t, resp)
					assert.Len(t, resp.Jobs, 1)
					assert.Equal(t, externalJobID, resp.Jobs[0].Id)
				},
			},
			{
				name: "non-existing job id",
				req: &jobv1.ListJobsRequest{
					Filter: &jobv1.ListJobsRequest_Filter{
						Ids: []string{"non-existing-job-id"},
					},
				},
				checkResp: func(t *testing.T, resp *jobv1.ListJobsResponse) {
					require.NotNil(t, resp)
					assert.Empty(t, resp.Jobs)
				},
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				// Call the ListJobs method
				listResp, err := jobClient.ListJobs(ctx, c.req)
				if c.checkErr != nil {
					c.checkErr(t, err)
					return
				}
				require.NoError(t, err)
				c.checkResp(t, listResp)
			})
		}
	})

	t.Run("GetProposal", func(t *testing.T) {
		t.Run("existing proposal", func(t *testing.T) {
			// Create a mock request
			getReq := &jobv1.GetProposalRequest{
				Id: expectedProposalID,
			}

			getResp, err := jobClient.GetProposal(ctx, getReq)
			require.NoError(t, err)
			assert.NotNil(t, getResp)
			assert.Equal(t, expectedProposal, getResp.Proposal)
		})

		t.Run("non-existing proposal", func(t *testing.T) {
			// Create a mock request
			getReq := &jobv1.GetProposalRequest{
				Id: "non-existing-job",
			}

			getResp, err := jobClient.GetProposal(ctx, getReq)
			require.Error(t, err)
			assert.Nil(t, getResp)
		})
	})

	t.Run("ListProposals", func(t *testing.T) {
		type listCase struct {
			name      string
			req       *jobv1.ListProposalsRequest
			checkErr  func(t *testing.T, err error)
			checkResp func(t *testing.T, resp *jobv1.ListProposalsResponse)
		}
		cases := []listCase{

			{
				name: "no filters",
				req:  &jobv1.ListProposalsRequest{},
				checkResp: func(t *testing.T, resp *jobv1.ListProposalsResponse) {
					assert.NotNil(t, resp)
					assert.Len(t, resp.Proposals, 1)
					assert.Equal(t, expectedProposalID, resp.Proposals[0].Id)
					assert.Equal(t, expectedProposal, resp.Proposals[0])
				},
			},
			{
				name: "with id filter",
				req: &jobv1.ListProposalsRequest{
					Filter: &jobv1.ListProposalsRequest_Filter{
						Ids: []string{expectedProposalID},
					},
				},
				checkResp: func(t *testing.T, resp *jobv1.ListProposalsResponse) {
					assert.NotNil(t, resp)
					assert.Len(t, resp.Proposals, 1)
					assert.Equal(t, expectedProposalID, resp.Proposals[0].Id)
					assert.Equal(t, expectedProposal, resp.Proposals[0])
				},
			},

			{
				name: "non-existing job id",
				req: &jobv1.ListProposalsRequest{
					Filter: &jobv1.ListProposalsRequest_Filter{
						Ids: []string{"non-existing-job-id"},
					},
				},
				checkResp: func(t *testing.T, resp *jobv1.ListProposalsResponse) {
					require.NotNil(t, resp)
					assert.Empty(t, resp.Proposals, "expected no proposals %v", resp.Proposals)
				},
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				listResp, err := jobClient.ListProposals(ctx, c.req)
				if c.checkErr != nil {
					c.checkErr(t, err)
					return
				}
				require.NoError(t, err)
				c.checkResp(t, listResp)
			})
		}
	})
}

func ptr(s string) *string {
	return &s
}

// need some non-ocr job type to avoid the ocr validation and the p2pwrapper check
func testJobProposalTOML(t *testing.T, externalJobId string) string {
	tomlString := `
type = "standardcapabilities"
schemaVersion = 1
externalJobID = "%s"
name = "hacking-%s"
forwardingAllowed = false
command = "/home/capabilities/nowhere"
config = ""
`
	return fmt.Sprintf(tomlString, externalJobId, externalJobId)
}
