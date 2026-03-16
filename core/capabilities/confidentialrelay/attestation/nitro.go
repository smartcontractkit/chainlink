package attestation

import (
	"bytes"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hf/nitrite"
)

// DomainSeparator is prepended to attestation payloads before hashing.
const DomainSeparator = "CONFIDENTIAL_COMPUTE_PAYLOAD"

// HexBytes is a custom type that can unmarshal hex strings into a byte slice.
// It also marshals byte slices back to hex strings for JSON output. This allows us
// to more easily parse AWS Nitro Measurements, which use hex byte strings in JSON.
type HexBytes []byte

func (h *HexBytes) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("HexBytes: cannot unmarshal JSON into string: %w", err)
	}

	decoded, err := hex.DecodeString(s)
	if err != nil {
		return fmt.Errorf("HexBytes: failed to decode hex string '%s': %w", s, err)
	}
	*h = decoded
	return nil
}

func (h HexBytes) MarshalJSON() ([]byte, error) {
	s := hex.EncodeToString(h)
	return json.Marshal(s)
}

// PCRs uses our custom HexBytes type for PCR values.
type PCRs struct {
	PCR0 HexBytes `json:"pcr0"`
	PCR1 HexBytes `json:"pcr1"`
	PCR2 HexBytes `json:"pcr2"`
}

// DefaultCARoots is the AWS Nitro Enclaves root certificate.
// Downloaded from: https://aws-nitro-enclaves.amazonaws.com/AWS_NitroEnclaves_Root-G1.zip
const DefaultCARoots = "-----BEGIN CERTIFICATE-----\nMIICETCCAZagAwIBAgIRAPkxdWgbkK/hHUbMtOTn+FYwCgYIKoZIzj0EAwMwSTEL\nMAkGA1UEBhMCVVMxDzANBgNVBAoMBkFtYXpvbjEMMAoGA1UECwwDQVdTMRswGQYD\nVQQDDBJhd3Mubml0cm8tZW5jbGF2ZXMwHhcNMTkxMDI4MTMyODA1WhcNNDkxMDI4\nMTQyODA1WjBJMQswCQYDVQQGEwJVUzEPMA0GA1UECgwGQW1hem9uMQwwCgYDVQQL\nDANBV1MxGzAZBgNVBAMMEmF3cy5uaXRyby1lbmNsYXZlczB2MBAGByqGSM49AgEG\nBSuBBAAiA2IABPwCVOumCMHzaHDimtqQvkY4MpJzbolL//Zy2YlES1BR5TSksfbb\n48C8WBoyt7F2Bw7eEtaaP+ohG2bnUs990d0JX28TcPQXCEPZ3BABIeTPYwEoCWZE\nh8l5YoQwTcU/9KNCMEAwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUkCW1DdkF\nR+eWw5b6cp3PmanfS5YwDgYDVR0PAQH/BAQDAgGGMAoGCCqGSM49BAMDA2kAMGYC\nMQCjfy+Rocm9Xue4YnwWmNJVA44fA0P5W2OpYow9OYCVRaEevL8uO1XYru5xtMPW\nrfMCMQCi85sWBbJwKKXdS6BptQFuZbT73o/gBh1qUxl/nNr12UO8Yfwr6wPLb+6N\nIwLz3/Y=\n-----END CERTIFICATE-----\n"

// DomainHash computes SHA-256 over DomainSeparator + "\n" + tag + "\n" + data.
// This is the standard domain-separated hash used for attestation UserData
// throughout the system.
func DomainHash(tag string, data []byte) []byte {
	h := sha256.New()
	h.Write([]byte(DomainSeparator))
	h.Write([]byte("\n" + tag + "\n"))
	h.Write(data)
	return h.Sum(nil)
}

// ValidateNitroAttestation verifies an AWS Nitro attestation document against
// expected user data and trusted PCR measurements.
func ValidateNitroAttestation(attestation, expectedUserData, trustedMeasurements []byte, caRootsPEM string) error {
	if attestation == nil {
		return fmt.Errorf("attestation is nil")
	}

	roots := DefaultCARoots
	if caRootsPEM != "" {
		roots = caRootsPEM
	}
	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM([]byte(roots))
	if !ok {
		return fmt.Errorf("failed to parse CA roots")
	}
	result, err := nitrite.Verify(attestation, nitrite.VerifyOptions{
		CurrentTime: time.Now(),
		Roots:       pool,
	})
	if err != nil {
		return fmt.Errorf("failed to verify nitro attestation: %w", err)
	}
	if !result.SignatureOK {
		return fmt.Errorf("signature verification failed")
	}

	if !bytes.Equal(expectedUserData, result.Document.UserData) {
		return fmt.Errorf("expected user data %x, got %x", expectedUserData, result.Document.UserData)
	}

	var trustedPCRs PCRs
	if err := json.Unmarshal(trustedMeasurements, &trustedPCRs); err != nil {
		return fmt.Errorf("failed to unmarshal trusted PCRs: %w", err)
	}
	if !bytes.Equal(result.Document.PCRs[0], trustedPCRs.PCR0) {
		return fmt.Errorf("expected PCR0 %x, got %x", trustedPCRs.PCR0, result.Document.PCRs[0])
	}
	if !bytes.Equal(result.Document.PCRs[1], trustedPCRs.PCR1) {
		return fmt.Errorf("expected PCR1 %x, got %x", trustedPCRs.PCR1, result.Document.PCRs[1])
	}
	if !bytes.Equal(result.Document.PCRs[2], trustedPCRs.PCR2) {
		return fmt.Errorf("expected PCR2 %x, got %x", trustedPCRs.PCR2, result.Document.PCRs[2])
	}
	return nil
}
