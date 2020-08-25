package job

import (
	"fmt"
)

type PipelineStage interface {
	GetID() uint64
	SetNotifiee(n Notifiee)
}

type Notifiee interface {
	OnBeginStage(stage PipelineStage, input interface{})
	OnEndStage(stage PipelineStage, output interface{}, err error)
}

type PipelineError struct {
	FailedStage interface{}
	StageIndex  uint
	Error       error
}

func (e PipelineError) Error() string {
	if e.StageType == PipelineStageTypeFetcher {
		return fmt.Sprintf("(%v) %+v", e.StageName, e.Error)
	} else if e.StageType == PipelineStageTypeTransformer {
		return fmt.Sprintf("(transformer stage %v, %v) %+v", e.StageIndex, e.StageName, e.Error)
	}
}
