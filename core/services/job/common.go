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
		Stop() error
	}

	Config interface {
		DatabaseURL() string
		JobPipelineDBPollInterval() time.Duration
		JobPipelineParallelism() uint8
	}
)
