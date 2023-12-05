package config

type Transmission interface {
	TLS() TransmissionTLS
}

type TransmissionTLS interface {
	CertPath() string
}
