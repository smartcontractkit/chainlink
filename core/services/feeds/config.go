package feeds

import (
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

//go:generate mockery --quiet --name Config --output ./mocks/ --case=underscore

type Config interface {
	pg.QConfig
	config.OCR2Config
	Dev() bool
	FeatureOffchainReporting() bool
	FeatureOffchainReporting2() bool
	DefaultHTTPTimeout() models.Duration
	JobPipelineResultWriteQueueDepth() uint64
	JobPipelineMaxSuccessfulRuns() uint64
}
