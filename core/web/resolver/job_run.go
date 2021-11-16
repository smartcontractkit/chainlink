package resolver

import (
	"strconv"

	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

var outputRetrievalErrorStr = "error: unable to retrieve outputs"

type JobRunResolver struct {
	run pipeline.Run
}

func NewJobRun(run pipeline.Run) *JobRunResolver {
	return &JobRunResolver{run}
}

func NewJobRuns(runs []pipeline.Run) []*JobRunResolver {
	var resolvers []*JobRunResolver

	for _, run := range runs {
		resolvers = append(resolvers, NewJobRun(run))
	}

	return resolvers
}

func (r *JobRunResolver) Outputs() []*string {
	if !r.run.Outputs.Valid {
		return []*string{&outputRetrievalErrorStr}
	}

	outputs, err := r.run.StringOutputs()
	if err != nil {
		errMsg := err.Error()
		return []*string{&errMsg}
	}

	return outputs
}

func (r *JobRunResolver) PipelineSpecID() graphql.ID {
	return graphql.ID(strconv.Itoa(int(r.run.PipelineSpecID)))
}

func (r *JobRunResolver) FatalErrors() []*string {
	return r.run.StringAllErrors()
}

func (r *JobRunResolver) AllErrors() []*string {
	return r.run.StringAllErrors()
}

func (r *JobRunResolver) Inputs() string {
	val, err := r.run.Inputs.MarshalJSON()
	if err != nil {
		return "error: unable to retrieve inputs"
	}

	return string(val)
}

// ObservationSource resolves the job's observation source.
//
// This could potentially be moved to a dataloader in the future as we are
// fetching it from a relationship.
func (r *JobRunResolver) ObservationSource() string {
	return r.run.PipelineSpec.DotDagSource
}
