package protocol

import (
	"bytes"
	"errors"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting/internal/protocol/observation"
	"github.com/smartcontractkit/libocr/offchainreporting/internal/signature"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
)

type SignedObservation struct {
	Observation observation.Observation
	Signature   []byte
}

func MakeSignedObservation(
	observation observation.Observation,
	repctx ReportContext,
	signer func(msg []byte) (sig []byte, err error),
) (
	SignedObservation,
	error,
) {
	payload := signedObservationWireMessage(repctx, observation)
	sig, err := signer(payload)
	if err != nil {
		return SignedObservation{}, err
	}
	return SignedObservation{observation, sig}, nil
}

func (so SignedObservation) Equal(so2 SignedObservation) bool {
	return so.Observation.Equal(so2.Observation) &&
		bytes.Equal(so.Signature, so2.Signature)
}

func (so SignedObservation) Verify(repctx ReportContext, publicKey types.OffchainPublicKey) error {
	if so.Observation.IsMissingValue() {
		return errors.New("Observation is missing value")
	}

	sigPublicKey := signature.OffchainPublicKey(publicKey)
	if !sigPublicKey.Verify(signedObservationWireMessage(repctx, so.Observation), so.Signature) {
		return errors.New("SignedObservation has invalid signature")
	}

	return nil
}

func signedObservationWireMessage(repctx ReportContext, observation observation.Observation) []byte {
	tag := repctx.DomainSeparationTag()
	return append(tag[:], observation.Marshal()...)
}

type AttributedSignedObservation struct {
	SignedObservation SignedObservation
	Observer          commontypes.OracleID
}

func (aso AttributedSignedObservation) Equal(aso2 AttributedSignedObservation) bool {
	return aso.SignedObservation.Equal(aso2.SignedObservation) &&
		aso.Observer == aso2.Observer
}
