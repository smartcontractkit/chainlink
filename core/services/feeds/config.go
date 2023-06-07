package feeds

import (
	"github.com/smartcontractkit/chainlink/v2/core/config"
	ocr2models "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/models"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type Config interface {
	config.OCR2Config
	OCRDevelopmentMode() bool
	FeatureOffchainReporting() bool
	FeatureOffchainReporting2() bool
	MercuryCredentials(credName string) *ocr2models.MercuryCredentials
	// ThresholdKeyShare is unused in feeds, to be refactored to decouple from Secrets interface in core/config/secrets.go
	ThresholdKeyShare() string
}

type JobConfig interface {
	DefaultHTTPTimeout() models.Duration
}
