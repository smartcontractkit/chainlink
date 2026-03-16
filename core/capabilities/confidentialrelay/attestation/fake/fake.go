// Package fake provides a FakeAttestor that produces structurally valid
// COSE Sign1 attestation documents. These documents pass nitrite.Verify's
// full validation chain (CBOR parsing, cert chain, ECDSA signature, UserData,
// PCRs) without requiring real Nitro hardware.
package fake

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha512"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"

	"github.com/fxamacker/cbor/v2"
)

// FakeAttestor produces structurally valid COSE Sign1 attestation documents
// that pass nitrite.Verify with a custom CA root.
type FakeAttestor struct {
	rootKey     *ecdsa.PrivateKey
	rootCert    *x509.Certificate
	rootCertDER []byte
	leafKey     *ecdsa.PrivateKey
	leafCert    *x509.Certificate
	leafCertDER []byte
	pcrs        map[uint][]byte // PCR0, PCR1, PCR2 (48 bytes each)
}

// NewFakeAttestor generates a self-signed P-384 root CA, a leaf cert signed
// by that root, and deterministic 48-byte fake PCR values.
func NewFakeAttestor() (*FakeAttestor, error) {
	rootKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate root key: %w", err)
	}
	rootTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "Fake Nitro Root CA"},
		NotBefore:             time.Now().Add(-1 * time.Hour),
		NotAfter:              time.Now().Add(10 * 365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		IsCA:                  true,
		BasicConstraintsValid: true,
	}
	rootCertDER, err := x509.CreateCertificate(rand.Reader, rootTemplate, rootTemplate, &rootKey.PublicKey, rootKey)
	if err != nil {
		return nil, fmt.Errorf("create root cert: %w", err)
	}
	rootCert, err := x509.ParseCertificate(rootCertDER)
	if err != nil {
		return nil, fmt.Errorf("parse root cert: %w", err)
	}

	leafKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate leaf key: %w", err)
	}
	leafTemplate := &x509.Certificate{
		SerialNumber:       big.NewInt(2),
		Subject:            pkix.Name{CommonName: "Fake Nitro Enclave"},
		NotBefore:          time.Now().Add(-1 * time.Hour),
		NotAfter:           time.Now().Add(24 * time.Hour),
		KeyUsage:           x509.KeyUsageDigitalSignature,
		ExtKeyUsage:        []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		SignatureAlgorithm: x509.ECDSAWithSHA384,
	}
	leafCertDER, err := x509.CreateCertificate(rand.Reader, leafTemplate, rootCert, &leafKey.PublicKey, rootKey)
	if err != nil {
		return nil, fmt.Errorf("create leaf cert: %w", err)
	}
	leafCert, err := x509.ParseCertificate(leafCertDER)
	if err != nil {
		return nil, fmt.Errorf("parse leaf cert: %w", err)
	}

	pcrs := map[uint][]byte{
		0: sha384Sum([]byte("fake-pcr-0")),
		1: sha384Sum([]byte("fake-pcr-1")),
		2: sha384Sum([]byte("fake-pcr-2")),
	}

	return &FakeAttestor{
		rootKey:     rootKey,
		rootCert:    rootCert,
		rootCertDER: rootCertDER,
		leafKey:     leafKey,
		leafCert:    leafCert,
		leafCertDER: leafCertDER,
		pcrs:        pcrs,
	}, nil
}

// CreateAttestation builds a COSE Sign1 document encoding a Nitro-like
// attestation with the given userData.
func (f *FakeAttestor) CreateAttestation(userData []byte) ([]byte, error) {
	doc := attestationDocument{
		ModuleID:    "fake-enclave-module",
		Timestamp:   uint64(time.Now().UnixMilli()),
		Digest:      "SHA384",
		PCRs:        f.pcrs,
		Certificate: f.leafCertDER,
		CABundle:    [][]byte{f.rootCertDER},
		UserData:    userData,
	}

	payloadBytes, err := cbor.Marshal(doc)
	if err != nil {
		return nil, fmt.Errorf("cbor encode document: %w", err)
	}

	header := coseHeader{Alg: int64(-35)}
	protectedBytes, err := cbor.Marshal(header)
	if err != nil {
		return nil, fmt.Errorf("cbor encode protected header: %w", err)
	}

	sigStruct := coseSignature{
		Context:     "Signature1",
		Protected:   protectedBytes,
		ExternalAAD: []byte{},
		Payload:     payloadBytes,
	}
	sigStructBytes, err := cbor.Marshal(sigStruct)
	if err != nil {
		return nil, fmt.Errorf("cbor encode sig structure: %w", err)
	}

	hash := sha512.Sum384(sigStructBytes)
	r, s, err := ecdsa.Sign(rand.Reader, f.leafKey, hash[:])
	if err != nil {
		return nil, fmt.Errorf("ecdsa sign: %w", err)
	}

	signature := make([]byte, 96)
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	copy(signature[48-len(rBytes):48], rBytes)
	copy(signature[96-len(sBytes):96], sBytes)

	outer := cosePayload{
		Protected: protectedBytes,
		Payload:   payloadBytes,
		Signature: signature,
	}
	result, err := cbor.Marshal(outer)
	if err != nil {
		return nil, fmt.Errorf("cbor encode cose sign1: %w", err)
	}
	return result, nil
}

// CARoots returns an x509.CertPool containing the fake root CA certificate.
func (f *FakeAttestor) CARoots() *x509.CertPool {
	pool := x509.NewCertPool()
	pool.AddCert(f.rootCert)
	return pool
}

// CARootsPEM returns the root CA certificate in PEM format.
func (f *FakeAttestor) CARootsPEM() string {
	return string(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: f.rootCertDER,
	}))
}

// TrustedPCRsJSON returns the PCR values as a JSON object matching the
// format expected by the attestation validator.
func (f *FakeAttestor) TrustedPCRsJSON() []byte {
	m := map[string]string{
		"pcr0": hex.EncodeToString(f.pcrs[0]),
		"pcr1": hex.EncodeToString(f.pcrs[1]),
		"pcr2": hex.EncodeToString(f.pcrs[2]),
	}
	b, _ := json.Marshal(m)
	return b
}

func sha384Sum(data []byte) []byte {
	h := sha512.Sum384(data)
	return h[:]
}

// CBOR structures matching what nitrite.Verify expects.

type attestationDocument struct {
	ModuleID    string          `cbor:"module_id"`
	Timestamp   uint64          `cbor:"timestamp"`
	Digest      string          `cbor:"digest"`
	PCRs        map[uint][]byte `cbor:"pcrs"`
	Certificate []byte          `cbor:"certificate"`
	CABundle    [][]byte        `cbor:"cabundle"`
	PublicKey   []byte          `cbor:"public_key,omitempty"`
	UserData    []byte          `cbor:"user_data,omitempty"`
	Nonce       []byte          `cbor:"nonce,omitempty"`
}

type coseHeader struct {
	Alg int64 `cbor:"1,keyasint"`
}

type cosePayload struct {
	_           struct{} `cbor:",toarray"`
	Protected   []byte
	Unprotected cbor.RawMessage
	Payload     []byte
	Signature   []byte
}

type coseSignature struct {
	_           struct{} `cbor:",toarray"`
	Context     string
	Protected   []byte
	ExternalAAD []byte
	Payload     []byte
}
