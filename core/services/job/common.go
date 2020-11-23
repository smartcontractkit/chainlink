package job

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

//go:generate mockery --name Spec --output ./mocks/ --case=underscore
//go:generate mockery --name Service --output ./mocks/ --case=underscore

type (
	Type string

	Spec interface {
		JobID() int32
		JobType() Type
		TaskDAG() pipeline.TaskDAG
		TableName() string
	}

	Service interface {
		Start() error
		Close() error
	}

	Config interface {
		DatabaseMaximumTxDuration() time.Duration
		DatabaseURL() string
		TriggerFallbackDBPollInterval() time.Duration
		JobPipelineParallelism() uint8
	}
)
