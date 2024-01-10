package config

import (
	"time"

	ocr2models "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/models"
)

type MercuryCache interface {
	LatestReportTTL() time.Duration
	MaxStaleAge() time.Duration
	LatestReportDeadline() time.Duration
}

type MercuryTLS interface {
	CertFile() string
}

type Mercury interface {
	Credentials(credName string) *ocr2models.MercuryCredentials
	Cache() MercuryCache
	TLS() MercuryTLS
}
