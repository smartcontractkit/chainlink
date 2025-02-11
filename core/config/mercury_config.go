package config

import (
	"fmt"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type MercuryTLS interface {
	CertFile() string
}

type MercuryTransmitterProtocol string

const (
	MercuryTransmitterProtocolGRPC MercuryTransmitterProtocol = "grpc"
)

func (m MercuryTransmitterProtocol) String() string {
	return string(m)
}

func (m *MercuryTransmitterProtocol) UnmarshalText(text []byte) error {
	switch string(text) {
	case "grpc":
		*m = MercuryTransmitterProtocolGRPC
	default:
		return fmt.Errorf("unknown mercury transmitter protocol: %s", text)
	}
	return nil
}

type MercuryTransmitter interface {
	Protocol() MercuryTransmitterProtocol
	TransmitQueueMaxSize() uint32
	TransmitTimeout() commonconfig.Duration
	TransmitConcurrency() uint32
	ReaperFrequency() commonconfig.Duration
	ReaperMaxAge() commonconfig.Duration
}

type Mercury interface {
	Credentials(credName string) *types.MercuryCredentials
	TLS() MercuryTLS
	Transmitter() MercuryTransmitter
	VerboseLogging() bool
}
