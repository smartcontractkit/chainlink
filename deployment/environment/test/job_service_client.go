package test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml/v2"
	"google.golang.org/grpc"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/v2/core/services/feeds"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

type getter[T any] interface {
	Get(string) (T, error)
}

type JobApprover interface {
	AutoApproveJob(ctx context.Context, p *feeds.ProposeJobArgs) error
}

type JobServiceClient struct {
	jobStore
	proposalStore
	jobApproverStore getter[JobApprover]
}

const (
	LabelDoNotAutoApprove = "doNotAutoApprove"
)

func NewJobServiceClient(jg getter[JobApprover]) *JobServiceClient {
	return &JobServiceClient{
		jobStore:         newMapJobStore(),
		proposalStore:    newMapProposalStore(),
		jobApproverStore: jg,
	}
}

func (j *JobServiceClient) BatchProposeJob(ctx context.Context, in *jobv1.BatchProposeJobRequest, opts ...grpc.CallOption) (*jobv1.BatchProposeJobResponse, error) {
	targets := make(map[string]struct{})
	for _, nodeID := range in.NodeIds {
		_, err := j.jobApproverStore.Get(nodeID)
		if err != nil {
			return nil, fmt.Errorf("node not found: %s", nodeID)
		}
		targets[nodeID] = struct{}{}
	}
	if len(targets) == 0 {
		return nil, errors.New("no nodes found")
	}
	out := &jobv1.BatchProposeJobResponse{
		SuccessResponses: make(map[string]*jobv1.ProposeJobResponse),
		FailedResponses:  make(map[string]*jobv1.ProposeJobFailure),
	}
	var totalErr error
	for id := range targets {
		singleReq := &jobv1.ProposeJobRequest{
			NodeId: id,
			Spec:   in.Spec,
			Labels: in.Labels,
		}
		resp, err := j.ProposeJob(ctx, singleReq)
		if err != nil {
			out.FailedResponses[id] = &jobv1.ProposeJobFailure{
				ErrorMessage: err.Error(),
			}
			totalErr = errors.Join(totalErr, fmt.Errorf("failed to propose job for node %s: %w", id, err))
		} else {
			out.SuccessResponses[id] = resp
		}
	}
	return out, totalErr
}

func (j *JobServiceClient) UpdateJob(ctx context.Context, in *jobv1.UpdateJobRequest, opts ...grpc.CallOption) (*jobv1.UpdateJobResponse, error) {
	// TODO CCIP-3108 implement me
	panic("implement me")
}

func (j *JobServiceClient) GetJob(ctx context.Context, in *jobv1.GetJobRequest, opts ...grpc.CallOption) (*jobv1.GetJobResponse, error) {
	// implementation detail that job id and uuid is the same
	jb, err := j.jobStore.get(in.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to get job: %w", err)
	}
	// TODO CCIP-3108 implement me
	return &jobv1.GetJobResponse{
		Job: jb,
	}, nil
}

func (j *JobServiceClient) GetProposal(ctx context.Context, in *jobv1.GetProposalRequest, opts ...grpc.CallOption) (*jobv1.GetProposalResponse, error) {
	p, err := j.proposalStore.get(in.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get proposal: %w", err)
	}
	return &jobv1.GetProposalResponse{
		Proposal: p,
	}, nil
}

func (j *JobServiceClient) ListJobs(ctx context.Context, in *jobv1.ListJobsRequest, opts ...grpc.CallOption) (*jobv1.ListJobsResponse, error) {
	jbs, err := j.jobStore.list(in.Filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list jobs: %w", err)
	}

	return &jobv1.ListJobsResponse{
		Jobs: jbs,
	}, nil
}

func (j *JobServiceClient) ListProposals(ctx context.Context, in *jobv1.ListProposalsRequest, opts ...grpc.CallOption) (*jobv1.ListProposalsResponse, error) {
	proposals, err := j.proposalStore.list(in.Filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list proposals: %w", err)
	}
	return &jobv1.ListProposalsResponse{
		Proposals: proposals,
	}, nil
}

// ProposeJob is used to propose a job to the node
// It auto approves the job
func (j *JobServiceClient) ProposeJob(ctx context.Context, in *jobv1.ProposeJobRequest, opts ...grpc.CallOption) (*jobv1.ProposeJobResponse, error) {
	n, err := j.jobApproverStore.Get(in.NodeId)
	if err != nil {
		return nil, fmt.Errorf("node not found: %w", err)
	}
	_, err = job.ValidateSpec(in.Spec)
	if err != nil {
		return nil, fmt.Errorf("failed to validate job spec: %w", err)
	}
	var extractor ExternalJobIDExtractor
	err = toml.Unmarshal([]byte(in.Spec), &extractor)
	if err != nil {
		return nil, fmt.Errorf("failed to load job spec: %w", err)
	}
	if extractor.ExternalJobID == "" {
		return nil, errors.New("externalJobID is required")
	}

	// must auto increment the version to avoid collision on the node side
	proposals, err := j.proposalStore.list(&jobv1.ListProposalsRequest_Filter{
		JobIds: []string{extractor.ExternalJobID},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list proposals: %w", err)
	}
	proposalVersion := getProposalVersion(proposals)
	pargs := &feeds.ProposeJobArgs{
		FeedsManagerID: 1,
		Spec:           in.Spec,
		RemoteUUID:     uuid.MustParse(extractor.ExternalJobID),
		Version:        proposalVersion,
	}
	// Allow for skipping auto-approval by supplying the LabelDoNotAutoApprove label.
	autoApprove := true
	for _, label := range in.Labels {
		if label.Key == LabelDoNotAutoApprove {
			autoApprove = false
			break
		}
	}
	status := jobv1.ProposalStatus_PROPOSAL_STATUS_PENDING
	if autoApprove {
		err = n.AutoApproveJob(ctx, pargs)
		if err != nil {
			return nil, fmt.Errorf("failed to auto approve job: %w", err)
		}
		status = jobv1.ProposalStatus_PROPOSAL_STATUS_APPROVED
	}

	storeProposalID := uuid.Must(uuid.NewRandom()).String()
	p := &jobv1.ProposeJobResponse{Proposal: &jobv1.Proposal{
		// make the proposal id the same as the job id for further reference
		// if you are changing this make sure to change the GetProposal and ListJobs method implementation
		Id:             storeProposalID,
		Revision:       int64(proposalVersion),
		Status:         status,
		DeliveryStatus: jobv1.ProposalDeliveryStatus_PROPOSAL_DELIVERY_STATUS_DELIVERED,
		Spec:           in.Spec,
		JobId:          extractor.ExternalJobID,
	}}

	// save the proposal and job
	{
		var (
			storeErr error // used to cleanup if we fail to save the job
			job      *jobv1.Job
		)

		storeErr = j.proposalStore.put(storeProposalID, p.Proposal)
		if storeErr != nil {
			return nil, fmt.Errorf("failed to save proposal: %w", err)
		}
		defer func() {
			// cleanup if we fail to save the job
			if storeErr != nil {
				_ = j.proposalStore.delete(storeProposalID)
			}
		}()

		job, storeErr = j.jobStore.get(extractor.ExternalJobID)
		if storeErr != nil && !errors.Is(storeErr, errNoExist) {
			return nil, fmt.Errorf("failed to get job: %w", storeErr)
		}
		if errors.Is(storeErr, errNoExist) {
			job = &jobv1.Job{
				Id:          extractor.ExternalJobID,
				Uuid:        extractor.ExternalJobID,
				NodeId:      in.NodeId,
				ProposalIds: []string{storeProposalID},
				Labels:      in.Labels,
			}
		} else {
			job.ProposalIds = append(job.ProposalIds, storeProposalID)
		}
		storeErr = j.jobStore.put(extractor.ExternalJobID, job)
		if storeErr != nil {
			return nil, fmt.Errorf("failed to save job: %w", storeErr)
		}
	}
	return p, nil
}

// getProposalVersion returns the next proposal version, based on all proposals made for this job and all of their
// respective revisions.
func getProposalVersion(proposals []*jobv1.Proposal) int32 {
	totalRevisions := int64(0)
	for _, p := range proposals {
		totalRevisions += p.Revision
	}
	return int32(totalRevisions + 1) //nolint:gosec // G115
}

func (j *JobServiceClient) RevokeJob(ctx context.Context, in *jobv1.RevokeJobRequest, opts ...grpc.CallOption) (*jobv1.RevokeJobResponse, error) {
	// Get all proposals for this job.
	proposals, err := j.proposalStore.list(&jobv1.ListProposalsRequest_Filter{
		JobIds: []string{in.GetId()},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list proposals: %w", err)
	}
	if len(proposals) == 0 {
		return nil, fmt.Errorf("no proposals found for job %s", in.GetId())
	}
	// Get the latest proposal.
	prop := proposals[0]
	for _, p := range proposals {
		if p != nil && p.UpdatedAt != nil &&
			(prop == nil || prop.UpdatedAt == nil || p.UpdatedAt.GetNanos() > prop.UpdatedAt.GetNanos()) {
			prop = p
		}
	}
	// Check if it's revokable.
	if prop.Status != jobv1.ProposalStatus_PROPOSAL_STATUS_PENDING &&
		prop.Status != jobv1.ProposalStatus_PROPOSAL_STATUS_CANCELLED {
		return nil, fmt.Errorf("proposal %s is not revokable: status %s", prop.Id, prop.Status.String())
	}
	// Revoke it.
	prop.Status = jobv1.ProposalStatus_PROPOSAL_STATUS_REVOKED
	// Store the revoked version.
	err = j.proposalStore.put(prop.Id, prop)
	if err != nil {
		return nil, fmt.Errorf("failed to save proposal: %w", err)
	}
	return &jobv1.RevokeJobResponse{
		Proposal: prop,
	}, nil
}

func (j *JobServiceClient) DeleteJob(ctx context.Context, in *jobv1.DeleteJobRequest, opts ...grpc.CallOption) (*jobv1.DeleteJobResponse, error) {
	// TODO CCIP-3108 implement me
	panic("implement me")
}

type ExternalJobIDExtractor struct {
	ExternalJobID string `toml:"externalJobID"`
}

var errNoExist = errors.New("does not exist")

// proposalStore is an interface for storing job proposals.
type proposalStore interface {
	put(proposalID string, proposal *jobv1.Proposal) error
	get(proposalID string) (*jobv1.Proposal, error)
	list(filter *jobv1.ListProposalsRequest_Filter) ([]*jobv1.Proposal, error)
	delete(proposalID string) error
}

// jobStore is an interface for storing jobs.
type jobStore interface {
	put(jobID string, job *jobv1.Job) error
	get(jobID string) (*jobv1.Job, error)
	list(filter *jobv1.ListJobsRequest_Filter) ([]*jobv1.Job, error)
	delete(jobID string) error
}

var _ jobStore = &mapJobStore{}

type mapJobStore struct {
	mu            sync.Mutex
	jobs          map[string]*jobv1.Job
	nodesToJobIDs map[string][]string
	uuidToJobIDs  map[string][]string
}

func newMapJobStore() *mapJobStore {
	return &mapJobStore{
		jobs:          make(map[string]*jobv1.Job),
		nodesToJobIDs: make(map[string][]string),
		uuidToJobIDs:  make(map[string][]string),
	}
}

func (m *mapJobStore) put(jobID string, job *jobv1.Job) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.jobs == nil {
		m.jobs = make(map[string]*jobv1.Job)
		m.nodesToJobIDs = make(map[string][]string)
		m.uuidToJobIDs = make(map[string][]string)
	}
	m.jobs[jobID] = job
	if _, ok := m.nodesToJobIDs[job.NodeId]; !ok {
		m.nodesToJobIDs[job.NodeId] = make([]string, 0)
	}
	m.nodesToJobIDs[job.NodeId] = append(m.nodesToJobIDs[job.NodeId], jobID)
	if _, ok := m.uuidToJobIDs[job.Uuid]; !ok {
		m.uuidToJobIDs[job.Uuid] = make([]string, 0)
	}
	m.uuidToJobIDs[job.Uuid] = append(m.uuidToJobIDs[job.Uuid], jobID)
	return nil
}

func (m *mapJobStore) get(jobID string) (*jobv1.Job, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.jobs == nil {
		return nil, fmt.Errorf("%w: job not found: %s", errNoExist, jobID)
	}
	job, ok := m.jobs[jobID]
	if !ok {
		return nil, fmt.Errorf("%w: job not found: %s", errNoExist, jobID)
	}
	return job, nil
}

func (m *mapJobStore) list(filter *jobv1.ListJobsRequest_Filter) ([]*jobv1.Job, error) {
	if filter != nil {
		counter := 0
		if filter.NodeIds != nil {
			counter++
		}
		if filter.Uuids != nil {
			counter++
		}
		if filter.Ids != nil {
			counter++
		}
		if counter > 1 {
			return nil, errors.New("only one of NodeIds, Uuids or Ids can be set")
		}
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.jobs == nil {
		return []*jobv1.Job{}, nil
	}

	jobs := make([]*jobv1.Job, 0, len(m.jobs))

	if filter == nil || (filter.NodeIds == nil && filter.Uuids == nil && filter.Ids == nil) {
		for _, job := range m.jobs {
			if filter != nil && !matchesSelectors(filter.Selectors, job) {
				continue
			}
			jobs = append(jobs, job)
		}
		return jobs, nil
	}

	wantedJobIDs := make(map[string]struct{})
	// use node ids to construct wanted job ids
	switch {
	case filter.NodeIds != nil:
		for _, nodeID := range filter.NodeIds {
			jobIDs, ok := m.nodesToJobIDs[nodeID]
			if !ok {
				continue
			}
			for _, jobID := range jobIDs {
				wantedJobIDs[jobID] = struct{}{}
			}
		}
	case filter.Uuids != nil:
		for _, uuid := range filter.Uuids {
			jobIDs, ok := m.uuidToJobIDs[uuid]
			if !ok {
				continue
			}
			for _, jobID := range jobIDs {
				wantedJobIDs[jobID] = struct{}{}
			}
		}
	case filter.Ids != nil:
		for _, jobID := range filter.Ids {
			wantedJobIDs[jobID] = struct{}{}
		}
	default:
		panic("this should never happen because of the nil filter check")
	}

	for _, job := range m.jobs {
		if !matchesSelectors(filter.Selectors, job) {
			continue
		}
		if _, ok := wantedJobIDs[job.Id]; ok {
			jobs = append(jobs, job)
		}
	}
	return jobs, nil
}

func matchesSelectors(selectors []*ptypes.Selector, job *jobv1.Job) bool {
	for _, selector := range selectors {
		label := labelForKey(selector.Key, job.Labels)
		switch selector.Op {
		case ptypes.SelectorOp_EXIST:
			return label != nil
		case ptypes.SelectorOp_NOT_EXIST:
			return label == nil
		case ptypes.SelectorOp_EQ:
			if label == nil || *label.Value != *selector.Value {
				return false
			}
		case ptypes.SelectorOp_NOT_EQ:
			if label != nil && *label.Value == *selector.Value {
				return false
			}
		case ptypes.SelectorOp_IN:
			if label == nil || label.Value == nil || selector.Value == nil {
				return false
			}
			list := strings.Split(*selector.Value, ",")
			found := false
			for _, v := range list {
				if *label.Value == v {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		case ptypes.SelectorOp_NOT_IN:
			if label != nil && label.Value != nil {
				list := strings.Split(*selector.Value, ",")
				for _, v := range list {
					if *label.Value == v {
						return false
					}
				}
			}
		default:
			return false
		}
	}

	return true
}

func labelForKey(key string, labels []*ptypes.Label) *ptypes.Label {
	for _, label := range labels {
		if label.Key == key {
			return label
		}
	}
	return nil
}

func (m *mapJobStore) delete(jobID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.jobs == nil {
		return fmt.Errorf("job not found: %s", jobID)
	}
	job, ok := m.jobs[jobID]
	if !ok {
		return nil
	}
	delete(m.jobs, jobID)
	delete(m.nodesToJobIDs, job.NodeId)
	delete(m.uuidToJobIDs, job.Uuid)
	return nil
}

var _ proposalStore = &mapProposalStore{}

type mapProposalStore struct {
	mu                sync.Mutex
	proposals         map[string]*jobv1.Proposal
	jobIDToProposalID map[string]string
}

func newMapProposalStore() *mapProposalStore {
	return &mapProposalStore{
		proposals:         make(map[string]*jobv1.Proposal),
		jobIDToProposalID: make(map[string]string),
	}
}

func (m *mapProposalStore) put(proposalID string, proposal *jobv1.Proposal) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.proposals == nil {
		m.proposals = make(map[string]*jobv1.Proposal)
	}
	if m.jobIDToProposalID == nil {
		m.jobIDToProposalID = make(map[string]string)
	}
	m.proposals[proposalID] = proposal
	m.jobIDToProposalID[proposal.JobId] = proposalID
	return nil
}
func (m *mapProposalStore) get(proposalID string) (*jobv1.Proposal, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.proposals == nil {
		return nil, fmt.Errorf("proposal not found: %s", proposalID)
	}
	proposal, ok := m.proposals[proposalID]
	if !ok {
		return nil, fmt.Errorf("%w: proposal not found: %s", errNoExist, proposalID)
	}
	return proposal, nil
}
func (m *mapProposalStore) list(filter *jobv1.ListProposalsRequest_Filter) ([]*jobv1.Proposal, error) {
	if filter != nil && filter.GetIds() != nil && filter.GetJobIds() != nil {
		return nil, errors.New("only one of Ids or JobIds can be set")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.proposals == nil {
		return nil, nil
	}
	proposals := make([]*jobv1.Proposal, 0)
	// all proposals
	if filter == nil || (filter.GetIds() == nil && filter.GetJobIds() == nil) {
		for _, proposal := range m.proposals {
			proposals = append(proposals, proposal)
		}
		return proposals, nil
	}

	// can't both be nil at this point
	wantedProposalIDs := filter.GetIds()
	if wantedProposalIDs == nil {
		wantedProposalIDs = make([]string, 0)
		for _, jobID := range filter.GetJobIds() {
			proposalID, ok := m.jobIDToProposalID[jobID]
			if !ok {
				continue
			}
			wantedProposalIDs = append(wantedProposalIDs, proposalID)
		}
	}

	for _, want := range wantedProposalIDs {
		p, ok := m.proposals[want]
		if !ok {
			continue
		}
		proposals = append(proposals, p)
	}
	return proposals, nil
}
func (m *mapProposalStore) delete(proposalID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.proposals == nil {
		return fmt.Errorf("proposal not found: %s", proposalID)
	}

	delete(m.proposals, proposalID)
	return nil
}
