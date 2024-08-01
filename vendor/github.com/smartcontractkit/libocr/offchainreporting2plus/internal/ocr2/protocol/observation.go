package protocol

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/binary"
	"fmt"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type SignedObservation struct {
	Observation types.Observation
	Signature   []byte
}

func MakeSignedObservation(
	repts types.ReportTimestamp,
	query types.Query,
	observation types.Observation,
	signer func(msg []byte) (sig []byte, err error),
) (
	SignedObservation,
	error,
) {
	payload := signedObservationWireMessage(repts, query, observation)
	sig, err := signer(payload)
	if err != nil {
		return SignedObservation{}, err
	}
	return SignedObservation{observation, sig}, nil
}

func (so SignedObservation) Equal(so2 SignedObservation) bool {
	return bytes.Equal(so.Observation, so2.Observation) &&
		bytes.Equal(so.Signature, so2.Signature)
}

func (so SignedObservation) Verify(repts types.ReportTimestamp, query types.Query, publicKey types.OffchainPublicKey) error {
	pk := ed25519.PublicKey(publicKey[:])
	// should never trigger since types.OffchainPublicKey is an array with length ed25519.PublicKeySize
	if len(pk) != ed25519.PublicKeySize {
		return fmt.Errorf("ed25519 public key size mismatch, expected %v but got %v", ed25519.PublicKeySize, len(pk))
	}

	ok := ed25519.Verify(pk, signedObservationWireMessage(repts, query, so.Observation), so.Signature)
	if !ok {
		return fmt.Errorf("SignedObservation has invalid signature")
	}

	return nil
}

func signedObservationWireMessage(repts types.ReportTimestamp, query types.Query, observation types.Observation) []byte {
	h := sha256.New()
	// ConfigDigest
	_, _ = h.Write(repts.ConfigDigest[:])
	_ = binary.Write(h, binary.BigEndian, repts.Epoch)
	_, _ = h.Write([]byte{repts.Round})

	// Query
	_ = binary.Write(h, binary.BigEndian, uint64(len(query)))
	_, _ = h.Write(query)

	// Observation
	_ = binary.Write(h, binary.BigEndian, uint64(len(observation)))
	_, _ = h.Write(observation)

	return h.Sum(nil)
}

type AttributedSignedObservation struct {
	SignedObservation SignedObservation
	Observer          commontypes.OracleID
}

func (aso AttributedSignedObservation) Equal(aso2 AttributedSignedObservation) bool {
	return aso.SignedObservation.Equal(aso2.SignedObservation) &&
		aso.Observer == aso2.Observer
}
