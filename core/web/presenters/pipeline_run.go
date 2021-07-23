package presenters

import (
	"encoding/json"
	"github.com/smartcontractkit/chainlink/core/logger"
	"math/big"
	"time"

	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

// Corresponds with models.d.ts PipelineRun
type PipelineRunResource struct {
	JAID
	Outputs      []*string                 `json:"outputs"`
	Errors       []*string                 `json:"errors"`
	TaskRuns     []PipelineTaskRunResource `json:"taskRuns"`
	CreatedAt    time.Time                 `json:"createdAt"`
	FinishedAt   time.Time                 `json:"finishedAt"`
	PipelineSpec PipelineSpec              `json:"pipelineSpec"`
}

// GetName implements the api2go EntityNamer interface
func (r PipelineRunResource) GetName() string {
	return "pipelineRun"
}

func NewPipelineRunResource(pr pipeline.Run) PipelineRunResource {
	var trs []PipelineTaskRunResource
	for i := range pr.PipelineTaskRuns {
		trs = append(trs, NewPipelineTaskRunResource(pr.PipelineTaskRuns[i]))
	}
	// The UI expects all outputs to be strings.
	var outputs []*string
	if !pr.Outputs.Null {
		outs := pr.Outputs.Val.([]interface{})
		for _, out := range outs {
			switch v := out.(type) {
			case string:
				s := v
				outputs = append(outputs, &s)
			case map[string]interface{}:
				b, _ := json.Marshal(v)
				bs := string(b)
				outputs = append(outputs, &bs)
			case decimal.Decimal:
				s := v.String()
				outputs = append(outputs, &s)
			case *big.Int:
				s := v.String()
				outputs = append(outputs, &s)
			case nil:
				outputs = append(outputs, nil)
			default:
				logger.Default.Error("PipelineRunResource: unable to process output type", "out", out)
			}
		}
	}
	var errors []*string
	for _, err := range pr.Errors {
		if err.Valid {
			s := err.String
			errors = append(errors, &s)
		} else {
			errors = append(errors, nil)
		}
	}
	return PipelineRunResource{
		JAID:         NewJAIDInt64(pr.ID),
		Outputs:      outputs,
		Errors:       errors,
		TaskRuns:     trs,
		CreatedAt:    pr.CreatedAt,
		FinishedAt:   pr.FinishedAt.ValueOrZero(),
		PipelineSpec: NewPipelineSpec(&pr.PipelineSpec),
	}
}

// Corresponds with models.d.ts PipelineTaskRun
type PipelineTaskRunResource struct {
	Type       pipeline.TaskType `json:"type"`
	CreatedAt  time.Time         `json:"createdAt"`
	FinishedAt time.Time         `json:"finishedAt"`
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
	if tr.Output != nil && !tr.Output.Null {
		outputBytes, _ := tr.Output.MarshalJSON()
		outputStr := string(outputBytes)
		output = &outputStr
	}
	var error *string
	if tr.Error.Valid {
		error = &tr.Error.String
	}
	return PipelineTaskRunResource{
		Type:       tr.Type,
		CreatedAt:  tr.CreatedAt,
		FinishedAt: tr.FinishedAt.ValueOrZero(),
		Output:     output,
		Error:      error,
		DotID:      tr.GetDotID(),
	}
}
