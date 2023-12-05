package chainlink

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

func TestTransmissionConfig_TLS(t *testing.T) {
	certPath := "/path/to/cert.pem"
	transmission := toml.Transmission{
		TLS: toml.TransmissionTLS{
			CertPath: &certPath,
		},
	}
	cfg := transmissionConfig{transmission}

	assert.Equal(t, certPath, cfg.TLS().CertPath())
}
