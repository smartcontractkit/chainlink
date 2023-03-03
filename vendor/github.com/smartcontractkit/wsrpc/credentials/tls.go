package credentials

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"math/big"
)

type StaticSizedPublicKey [ed25519.PublicKeySize]byte

// NewClientTLSConfig uses the private key and public keys to construct a mutual
// TLS config for the client.
func NewClientTLSConfig(priv ed25519.PrivateKey, pubs *PublicKeys) (*tls.Config, error) {
	return newMutualTLSConfig(priv, pubs)
}

// NewServerTLSConfig uses the private key and public keys to construct a mutual
// TLS config for the server.
func NewServerTLSConfig(priv ed25519.PrivateKey, pubs *PublicKeys) (*tls.Config, error) {
	c, err := newMutualTLSConfig(priv, pubs)
	if err != nil {
		return nil, err
	}
	c.ClientAuth = tls.RequireAnyClientCert

	return c, nil
}

// newMutualTLSConfig uses the private key and public keys to construct a mutual
// TLS 1.3 config.
//
// We provide our own peer certificate verification function to check the
// certificate's public key matches our list of registered keys.
//
// Certificates are currently used similarly to GPG keys and only functionally
// as certificates to support the crypto/tls go module.
func newMutualTLSConfig(priv ed25519.PrivateKey, pubs *PublicKeys) (*tls.Config, error) {
	cert, err := newMinimalX509Cert(priv)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},

		// Since our clients use self-signed certs, we skip verification here.
		// Instead, we use VerifyPeerCertificate for our own check.
		//
		// If VerifyPeerCertificate changes to rely on standard x509 certificate
		// fields (such as, but not limited too CN, expiration date and time)
		// then it may be necessary to reconsider the use of InsecureSkipVerify.
		InsecureSkipVerify: true, //nolint:gosec

		MaxVersion: tls.VersionTLS13,
		MinVersion: tls.VersionTLS13,

		VerifyPeerCertificate: pubs.VerifyPeerCertificate(),
	}, nil
}

// Generates a minimal certificate (that wouldn't be considered valid outside of
// this networking protocol) from an Ed25519 private key.
func newMinimalX509Cert(priv ed25519.PrivateKey) (tls.Certificate, error) {
	template := x509.Certificate{
		SerialNumber: big.NewInt(0), // serial number must be set, so we set it to 0
	}

	encodedCert, err := x509.CreateCertificate(rand.Reader, &template, &template, priv.Public(), priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	return tls.Certificate{
		Certificate:                  [][]byte{encodedCert},
		PrivateKey:                   priv,
		SupportedSignatureAlgorithms: []tls.SignatureScheme{tls.Ed25519},
	}, nil
}

// PublicKeys wraps a slice of keys so we can update the keys dynamically.
type PublicKeys []ed25519.PublicKey

// Verifies that the certificate's public key matches with one of the keys in
// our list of registered keys.
func (r *PublicKeys) VerifyPeerCertificate() func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	return func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		if len(rawCerts) != 1 {
			return fmt.Errorf("required exactly one client certificate")
		}
		cert, err := x509.ParseCertificate(rawCerts[0])
		if err != nil {
			return err
		}
		pk, err := pubKeyFromCert(cert)
		if err != nil {
			return err
		}

		ok := isValidPublicKey(*r, pk)
		if !ok {
			return fmt.Errorf("unknown public key on cert %x", pk)
		}

		return nil
	}
}

// Replace replaces the existing keys with new keys. Use this to dynamically
// update the allowable keys at runtime.
func (r *PublicKeys) Replace(pubs []ed25519.PublicKey) {
	*r = PublicKeys(pubs)
}

// isValidPublicKey checks the public key against a list of valid keys.
func isValidPublicKey(valid []ed25519.PublicKey, pub ed25519.PublicKey) bool {
	for _, vpub := range valid {
		if pub.Equal(vpub) {
			return true
		}
	}

	return false
}

// PubKeyFromCert extracts the public key from the cert and returns it as a
// statically sized byte array.
func PubKeyFromCert(cert *x509.Certificate) (StaticSizedPublicKey, error) {
	pubKey, err := pubKeyFromCert(cert)
	if err != nil {
		return StaticSizedPublicKey{}, err
	}

	return ToStaticallySizedPublicKey(pubKey)
}

// pubKeyFromCert returns an ed25519 public key extracted from the certificate.
func pubKeyFromCert(cert *x509.Certificate) (ed25519.PublicKey, error) {
	if cert.PublicKeyAlgorithm != x509.Ed25519 {
		return nil, fmt.Errorf("requires an ed25519 public key")
	}

	pub, ok := cert.PublicKey.(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("invalid ed25519 public key")
	}

	return pub, nil
}

// ToStaticallySizedPublicKey converts an ed25519 public key into a statically
// sized byte array.
func ToStaticallySizedPublicKey(pubKey ed25519.PublicKey) (StaticSizedPublicKey, error) {
	var result [ed25519.PublicKeySize]byte

	if ed25519.PublicKeySize != copy(result[:], pubKey) {
		return StaticSizedPublicKey{}, errors.New("copying public key failed")
	}

	return result, nil
}
