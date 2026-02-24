package ocr2

import (
	"time"

	"github.com/smartcontractkit/chainlink/devenv/products/ocr2"
)

var DefaultProductionOCR2Config = &ocr2.OCRv2SetConfigOptions{
	RMax:                                    3,
	DeltaProgress:                           20 * time.Second,
	DeltaResend:                             20 * time.Second,
	DeltaStage:                              15 * time.Second,
	MaxDurationInitialization:               5 * time.Second,
	MaxDurationQuery:                        5 * time.Second,
	MaxDurationObservation:                  5 * time.Second,
	MaxDurationReport:                       5 * time.Second,
	MaxDurationShouldAcceptFinalizedReport:  5 * time.Second,
	MaxDurationShouldTransmitAcceptedReport: 5 * time.Second,
}
