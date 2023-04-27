package feeds

import (
	"github.com/smartcontractkit/chainlink/v2/core/config"
	ocr2models "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/models"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type Config interface {
	pg.QConfig
	config.OCR2Config
	Dev() bool
	OCRDevelopmentMode() bool
	FeatureOffchainReporting() bool
	FeatureOffchainReporting2() bool
	DefaultHTTPTimeout() models.Duration
	JobPipelineResultWriteQueueDepth() uint64
	JobPipelineMaxSuccessfulRuns() uint64
	MercuryCredentials(credName string) (*ocr2models.MercuryCredentials, error)
}
