package resolver

import (
	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

type TaskRunResolver struct {
	tr pipeline.TaskRun
}

func NewTaskRun(tr pipeline.TaskRun) *TaskRunResolver {
	return &TaskRunResolver{tr: tr}
}

func NewTaskRuns(runs []pipeline.TaskRun) []*TaskRunResolver {
	var resolvers []*TaskRunResolver

	for _, run := range runs {
		resolvers = append(resolvers, NewTaskRun(run))
	}

	return resolvers
}

func (r *TaskRunResolver) ID() graphql.ID {
	return graphql.ID(r.tr.ID.String())
}

func (r *TaskRunResolver) Type() string {
	return string(r.tr.Type)
}

func (r *TaskRunResolver) Output() string {
	val, err := r.tr.Output.MarshalJSON()
	if err != nil {
		return "error: unable to retrieve output"
	}
	return string(val)
}

func (r *TaskRunResolver) Error() *string {
	if r.tr.Error.Valid {
		return r.tr.Error.Ptr()
	}

	return nil
}

func (r *TaskRunResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.tr.CreatedAt}
}

func (r *TaskRunResolver) FinishedAt() *graphql.Time {
	return &graphql.Time{Time: r.tr.FinishedAt.ValueOrZero()}
}

func (r *TaskRunResolver) DotID() string {
	return r.tr.GetDotID()
}
