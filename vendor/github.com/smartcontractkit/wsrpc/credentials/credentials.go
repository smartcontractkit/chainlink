package credentials

import (
	"crypto/tls"
)

// TransportCredentials defines the TLS configuration for establishing a
// connection.
type TransportCredentials struct {
	Config     *tls.Config
	PublicKeys *PublicKeys
}

func NewTLS(config *tls.Config, publicKeys *PublicKeys) TransportCredentials {
	return TransportCredentials{
		Config:     config,
		PublicKeys: publicKeys,
	}
}
