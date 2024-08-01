package ragep2p

import (
	"crypto/tls"
	"crypto/x509"
)

func newTLSConfig(cert tls.Certificate, verifyPeerCertificate func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error) *tls.Config {
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAnyClientCert,

		// Since our clients use self-signed certs, we skip verification here.
		// Instead, we use VerifyPeerCertificate for our own check
		InsecureSkipVerify: true,

		MaxVersion: tls.VersionTLS13,
		MinVersion: tls.VersionTLS13,

		VerifyPeerCertificate: verifyPeerCertificate,
	}
}
