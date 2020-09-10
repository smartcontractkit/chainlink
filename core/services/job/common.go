package job

import (
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/model"
)

//go:generate mockery --name Spec --output ./mocks/ --case=underscore
//go:generate mockery --name Service --output ./mocks/ --case=underscore

type (
	Type string

	Spec interface {
		JobID() *models.ID
		JobType() Type
		TaskDAG() pipeline.TaskDAG
	}

	Service interface {
		Start() error
		Stop() error
	}
)
