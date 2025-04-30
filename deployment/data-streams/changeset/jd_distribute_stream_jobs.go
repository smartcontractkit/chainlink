package changeset

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/jd"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/jobs"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/pointer"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
)

var _ deployment.ChangeSetV2[CsDistributeStreamJobSpecsConfig] = CsDistributeStreamJobSpecs{}

type CsDistributeStreamJobSpecsConfig struct {
	Filter  *jd.ListFilter
	Streams []StreamSpecConfig

	// NodeNames specifies on which nodes to distribute the job specs.
	NodeNames []string
}

type StreamSpecConfig struct {
	StreamID   uint32
	Name       string
	StreamType jobs.StreamType
	// ReportFields should be QuoteReportFields, MedianReportFields, etc., based on the stream type.
	ReportFields    jobs.ReportFields
	EARequestParams EARequestParams
	APIs            []string
}

type EARequestParams struct {
	Endpoint string `json:"endpoint"`
	From     string `json:"from"`
	To       string `json:"to"`
}

type CsDistributeStreamJobSpecs struct{}

func (CsDistributeStreamJobSpecs) Apply(e deployment.Environment, cfg CsDistributeStreamJobSpecsConfig) (deployment.ChangesetOutput, error) {
	ctx, cancel := context.WithTimeout(e.GetContext(), defaultJobSpecsTimeout)
	defer cancel()

	// Add a label to the job spec to identify the related DON
	labels := append([]*ptypes.Label(nil),
		&ptypes.Label{
			Key: utils.DonIdentifier(cfg.Filter.DONID, cfg.Filter.DONName),
		},
	)

	oracleNodes, err := jd.FetchDONOraclesFromJD(ctx, e.Offchain, cfg.Filter, cfg.NodeNames)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get oracle nodes: %w", err)
	}

	var proposals []*jobv1.ProposeJobRequest
	for _, s := range cfg.Streams {
		for _, n := range oracleNodes {
			localLabels := append(labels, //nolint: gocritic // obvious and readable locally modified copy of labels
				&ptypes.Label{
					Key:   devenv.LabelStreamIDKey,
					Value: pointer.To(strconv.FormatUint(uint64(s.StreamID), 10)),
				},
				&ptypes.Label{
					Key:   devenv.LabelJobTypeKey,
					Value: pointer.To(devenv.LabelJobTypeValueStream),
				},
			)

			// Check if there is already a job spec for this stream on this node:
			externalJobID, err := fetchExternalJobID(e, n.Id, []*ptypes.Selector{
				{
					Key:   devenv.LabelStreamIDKey,
					Value: pointer.To(strconv.FormatUint(uint64(s.StreamID), 10)),
					Op:    ptypes.SelectorOp_EQ,
				},
			})
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to get externalJobID: %w", err)
			}

			spec, err := generateJobSpec(s, externalJobID)

			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to create stream job spec: %w", err)
			}
			renderedSpec, err := spec.MarshalTOML()
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to marshal stream job spec: %w", err)
			}

			proposals = append(proposals, &jobv1.ProposeJobRequest{
				NodeId: n.Id,
				Spec:   string(renderedSpec),
				Labels: localLabels,
			})
		}
	}

	proposedJobs, err := proposeAllOrNothing(ctx, e.Offchain, proposals)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to propose all oracle jobs: %w", err)
	}

	return deployment.ChangesetOutput{
		Jobs: proposedJobs,
	}, nil
}

func generateJobSpec(cc StreamSpecConfig, externalJobID uuid.UUID) (spec *jobs.StreamJobSpec, err error) {
	if externalJobID == uuid.Nil {
		externalJobID = uuid.New()
	}
	spec = &jobs.StreamJobSpec{
		Base: jobs.Base{
			Name:          fmt.Sprintf("%s | %d", cc.Name, cc.StreamID),
			Type:          jobs.JobSpecTypeStream,
			SchemaVersion: 1,
			ExternalJobID: externalJobID,
		},
		StreamID: cc.StreamID,
	}

	datasources := generateDatasources(cc)
	base := jobs.BaseObservationSource{
		Datasources:   datasources,
		AllowedFaults: len(datasources) - 1,
	}

	err = spec.SetObservationSource(base, cc.ReportFields)

	return spec, err
}

func generateDatasources(cc StreamSpecConfig) []jobs.Datasource {
	dss := make([]jobs.Datasource, len(cc.APIs))
	params := cc.EARequestParams
	for i, api := range cc.APIs {
		dss[i] = jobs.Datasource{
			BridgeName: api,
			ReqData:    fmt.Sprintf(`"{\"data\":{\"endpoint\":\"%s\",\"from\":\"%s\",\"to\":\"%s\"}}"`, params.Endpoint, params.From, params.To),
		}
	}
	return dss
}

func (f CsDistributeStreamJobSpecs) VerifyPreconditions(_ deployment.Environment, config CsDistributeStreamJobSpecsConfig) error {
	if config.Filter == nil {
		return errors.New("filter is required")
	}
	if config.Filter.DONID == 0 || config.Filter.DONName == "" {
		return errors.New("DONID and DONName are required")
	}
	if len(config.Streams) == 0 {
		return errors.New("streams are required")
	}
	for _, s := range config.Streams {
		if s.StreamID == 0 {
			return errors.New("streamID is required for each stream")
		}
		if s.Name == "" {
			return errors.New("name is required for each stream")
		}
		if !s.StreamType.Valid() {
			return errors.New("stream type is not valid")
		}
		if s.ReportFields == nil {
			return errors.New("report fields are required for each stream")
		}
		if s.EARequestParams.Endpoint == "" {
			return errors.New("endpoint is required for each EARequestParam on each stream")
		}
		if len(s.APIs) == 0 {
			return errors.New("at least one API is required for each stream")
		}
	}
	if len(config.NodeNames) == 0 {
		return errors.New("at least one node name is required")
	}
	// The list of node names tells us which nodes to distribute the job specs to.
	// The size of that list needs to match the filter size, i.e. the number of nodes we expect to get from JD.
	if config.Filter.NumOracleNodes+config.Filter.NumBootstrapNodes != len(config.NodeNames) {
		return fmt.Errorf("number of node names (%d) does not match filter size (%d)", len(config.NodeNames), config.Filter.NumOracleNodes+config.Filter.NumBootstrapNodes)
	}

	return nil
}
