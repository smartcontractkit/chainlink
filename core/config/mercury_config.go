package config

import (
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	mercurytransmitter "github.com/smartcontractkit/chainlink-data-streams/llo/transmitter/de"
)

type MercuryCache interface {
	LatestReportTTL() time.Duration
	MaxStaleAge() time.Duration
	LatestReportDeadline() time.Duration
}

type MercuryTLS interface {
	CertFile() string
}

type MercuryTransmitter interface {
	Protocol() mercurytransmitter.MercuryTransmitterProtocol
	TransmitQueueMaxSize() uint32
	TransmitTimeout() time.Duration
	TransmitConcurrency() uint32
	ReaperFrequency() time.Duration
	ReaperMaxAge() time.Duration
}

type Mercury interface {
	Credentials(credName string) *types.MercuryCredentials
	Cache() MercuryCache
	TLS() MercuryTLS
	Transmitter() MercuryTransmitter
	VerboseLogging() bool
}
