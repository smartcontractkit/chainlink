package feeds

import (
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type Config interface {
	config.OCR2Config
	OCR2Enabled() bool
}

type JobConfig interface {
	DefaultHTTPTimeout() models.Duration
}

type InsecureConfig interface {
	OCRDevelopmentMode() bool
}

type OCRConfig interface {
	Enabled() bool
}
