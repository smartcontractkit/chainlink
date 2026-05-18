package presenters

import (
	"time"

	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/jsonserializable"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

// Corresponds with models.d.ts PipelineRun
type PipelineRunResource struct {
	JAID
	Outputs []*string `json:"outputs"`
	// XXX: Here for backwards compatibility, can be removed later
	// Deprecated: Errors
	Errors       []*string                         `json:"errors"`
	AllErrors    []*string                         `json:"allErrors"`
	FatalErrors  []*string                         `json:"fatalErrors"`
	Inputs       jsonserializable.JSONSerializable `json:"inputs"`
	TaskRuns     []PipelineTaskRunResource         `json:"taskRuns"`
	CreatedAt    time.Time                         `json:"createdAt"`
	FinishedAt   null.Time                         `json:"finishedAt"`
	PipelineSpec PipelineSpec                      `json:"pipelineSpec"`
}

// GetName implements the api2go EntityNamer interface
func (r PipelineRunResource) GetName() string {
	return "pipelineRun"
}

func NewPipelineRunResource(pr pipeline.Run, lggr logger.Logger) PipelineRunResource {
	lggr = lggr.Named("PipelineRunResource")
	var trs []PipelineTaskRunResource
	for i := range pr.PipelineTaskRuns {
		trs = append(trs, NewPipelineTaskRunResource(pr.PipelineTaskRuns[i]))
	}

	outputs, err := pr.StringOutputs()
	if err != nil {
		lggr.Errorw(err.Error(), "out", pr.Outputs)
	}

	fatalErrors := pr.StringFatalErrors()

	return PipelineRunResource{
		JAID:         NewJAIDInt64(pr.ID),
		Outputs:      outputs,
		Errors:       fatalErrors,
		AllErrors:    pr.StringAllErrors(),
		FatalErrors:  fatalErrors,
		Inputs:       pr.Inputs,
		TaskRuns:     trs,
		CreatedAt:    pr.CreatedAt,
		FinishedAt:   pr.FinishedAt,
		PipelineSpec: NewPipelineSpec(&pr.PipelineSpec),
	}
}

// Corresponds with models.d.ts PipelineTaskRun
type PipelineTaskRunResource struct {
	Type       pipeline.TaskType `json:"type"`
	CreatedAt  time.Time         `json:"createdAt"`
	FinishedAt null.Time         `json:"finishedAt"`
	Output     *string           `json:"output"`
	Error      *string           `json:"error"`
	DotID      string            `json:"dotId"`
}

// GetName implements the api2go EntityNamer interface
func (r PipelineTaskRunResource) GetName() string {
	return "taskRun"
}

func NewPipelineTaskRunResource(tr pipeline.TaskRun) PipelineTaskRunResource {
	var output *string
	if tr.Output.Valid {
		outputBytes, _ := tr.Output.MarshalJSON()
		outputStr := string(outputBytes)
		output = &outputStr
	}
	var errString *string
	if tr.Error.Valid {
		errString = &tr.Error.String
	}
	return PipelineTaskRunResource{
		Type:       tr.Type,
		CreatedAt:  tr.CreatedAt,
		FinishedAt: tr.FinishedAt,
		Output:     output,
		Error:      errString,
		DotID:      tr.GetDotID(),
	}
}

func NewPipelineRunResources(prs []pipeline.Run, lggr logger.Logger) []PipelineRunResource {
	var out []PipelineRunResource

	for _, pr := range prs {
		out = append(out, NewPipelineRunResource(pr, lggr))
	}

	return out
}
