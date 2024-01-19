package config

import (
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
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
	Credentials(credName string) *types.MercuryCredentials
	Cache() MercuryCache
	TLS() MercuryTLS
}
