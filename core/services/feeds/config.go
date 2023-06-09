package feeds

import (
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type Config interface {
	config.OCR2Config
	FeatureOffchainReporting() bool
	FeatureOffchainReporting2() bool
}

type JobConfig interface {
	DefaultHTTPTimeout() models.Duration
}

type InsecureConfig interface {
	OCRDevelopmentMode() bool
}
