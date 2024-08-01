package serialization

import (
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting/internal/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting/internal/protocol/observation"
	"github.com/smartcontractkit/libocr/offchainreporting/internal/serialization/protobuf"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

// Serialize encodes a protocol.Message into a binary payload
func Serialize(m protocol.Message) (b []byte, pbm *protobuf.MessageWrapper, err error) {
	pbm, err = toProtoMessage(m)
	if err != nil {
		return nil, nil, err
	}
	b, err = proto.Marshal(pbm)
	if err != nil {
		return nil, nil, err
	}
	return b, pbm, nil
}

// Deserialize decodes a binary payload into a protocol.Message
func Deserialize(b []byte) (protocol.Message, *protobuf.MessageWrapper, error) {
	pbm := &protobuf.MessageWrapper{}
	err := proto.Unmarshal(b, pbm)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not unmarshal protobuf")
	}
	m, err := messageWrapperFromProtoMessage(pbm)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not translate protobuf to protocol.Message")
	}
	return m, pbm, nil
}

func toProtoMessage(m protocol.Message) (*protobuf.MessageWrapper, error) {
	msgWrapper := protobuf.MessageWrapper{}
	switch v := m.(type) {
	case protocol.MessageNewEpoch:
		pm := &protobuf.MessageNewEpoch{
			Epoch: uint64(v.Epoch),
		}
		msgWrapper.Msg = &protobuf.MessageWrapper_MessageNewEpoch{pm}
	case protocol.MessageObserveReq:
		pm := &protobuf.MessageObserveReq{
			Round: uint64(v.Round),
			Epoch: uint64(v.Epoch),
		}
		msgWrapper.Msg = &protobuf.MessageWrapper_MessageObserveReq{pm}
	case protocol.MessageObserve:
		pm := &protobuf.MessageObserve{
			Round:             uint64(v.Round),
			Epoch:             uint64(v.Epoch),
			SignedObservation: signedObservationToProtoMessage(v.SignedObservation),
		}
		msgWrapper.Msg = &protobuf.MessageWrapper_MessageObserve{pm}
	case protocol.MessageReportReq:
		pm := &protobuf.MessageReportReq{
			Round: uint64(v.Round),
			Epoch: uint64(v.Epoch),
		}
		for _, o := range v.AttributedSignedObservations {
			pm.AttributedSignedObservations = append(pm.AttributedSignedObservations,
				attributedSignedObservationToProtoMessage(o))
		}
		msgWrapper.Msg = &protobuf.MessageWrapper_MessageReportReq{pm}
	case protocol.MessageReport:
		pm := &protobuf.MessageReport{
			Epoch:  uint64(v.Epoch),
			Round:  uint64(v.Round),
			Report: attestedReportOneToProtoMessage(v.Report),
		}
		msgWrapper.Msg = &protobuf.MessageWrapper_MessageReport{pm}
	case protocol.MessageFinal:
		msgWrapper.Msg = &protobuf.MessageWrapper_MessageFinal{finalToProtoMessage(v)}
	case protocol.MessageFinalEcho:
		msgWrapper.Msg = &protobuf.MessageWrapper_MessageFinalEcho{
			&protobuf.MessageFinalEcho{Final: finalToProtoMessage(v.MessageFinal)},
		}
	default:
		return nil, errors.Errorf("Unable to serialize message of type %T", m)

	}
	return &msgWrapper, nil
}

func observationToProtoMessage(o observation.Observation) *protobuf.Observation {
	return &protobuf.Observation{Value: o.Marshal()}
}

func signedObservationToProtoMessage(o protocol.SignedObservation) *protobuf.SignedObservation {
	sig := o.Signature
	if sig == nil {
		sig = []byte{}
	}
	return &protobuf.SignedObservation{
		Observation: observationToProtoMessage(o.Observation),
		Signature:   sig,
	}
}

func attributedSignedObservationToProtoMessage(aso protocol.AttributedSignedObservation) *protobuf.AttributedSignedObservation {
	return &protobuf.AttributedSignedObservation{
		SignedObservation: signedObservationToProtoMessage(aso.SignedObservation),
		Observer:          uint32(aso.Observer),
	}
}

func messageWrapperFromProtoMessage(wrapper *protobuf.MessageWrapper) (protocol.Message, error) {
	switch msg := wrapper.Msg.(type) {
	case *protobuf.MessageWrapper_MessageNewEpoch:
		return messageNewEpochFromProtoMessage(wrapper.GetMessageNewEpoch())
	case *protobuf.MessageWrapper_MessageObserveReq:
		return messageObserveReqFromProtoMessage(wrapper.GetMessageObserveReq())
	case *protobuf.MessageWrapper_MessageObserve:
		return messageObserveFromProtoMessage(wrapper.GetMessageObserve())
	case *protobuf.MessageWrapper_MessageReportReq:
		return messageReportReqFromProtoMessage(wrapper.GetMessageReportReq())
	case *protobuf.MessageWrapper_MessageReport:
		return messageReportFromProtoMessage(wrapper.GetMessageReport())
	case *protobuf.MessageWrapper_MessageFinal:
		return messageFinalFromProtoMessage(wrapper.GetMessageFinal())
	case *protobuf.MessageWrapper_MessageFinalEcho:
		return messageFinalEchoFromProtoMessage(wrapper.GetMessageFinalEcho())
	default:
		return nil, errors.Errorf("Unrecognised Msg type %T", msg)
	}
}

func messageNewEpochFromProtoMessage(m *protobuf.MessageNewEpoch) (protocol.MessageNewEpoch, error) {
	if m == nil {
		return protocol.MessageNewEpoch{}, errors.New("Unable to extract a MessageNewEpoch value")
	}
	return protocol.MessageNewEpoch{
		Epoch: uint32(m.Epoch),
	}, nil
}

func messageObserveReqFromProtoMessage(m *protobuf.MessageObserveReq) (protocol.MessageObserveReq, error) {
	if m == nil {
		return protocol.MessageObserveReq{}, errors.New("Unable to extract a MessageObserveReq value")
	}
	return protocol.MessageObserveReq{
		Epoch: uint32(m.Epoch),
		Round: uint8(m.Round),
	}, nil
}

func messageObserveFromProtoMessage(m *protobuf.MessageObserve) (protocol.MessageObserve, error) {
	if m == nil {
		return protocol.MessageObserve{}, errors.New("Unable to extract a MessageObserve value")
	}
	so, err := signedObservationFromProtoMessage(m.SignedObservation)
	if err != nil {
		return protocol.MessageObserve{}, err
	}
	return protocol.MessageObserve{
		Epoch:             uint32(m.Epoch),
		Round:             uint8(m.Round),
		SignedObservation: so,
	}, nil
}

func messageReportReqFromProtoMessage(m *protobuf.MessageReportReq) (protocol.MessageReportReq, error) {
	if m == nil {
		return protocol.MessageReportReq{}, errors.New("Unable to extract a MessageReportReq value")
	}
	asos, err := attributedSignedObservationsFromProtoMessage(m.AttributedSignedObservations)
	if err != nil {
		return protocol.MessageReportReq{}, err
	}
	return protocol.MessageReportReq{
		Epoch:                        uint32(m.Epoch),
		Round:                        uint8(m.Round),
		AttributedSignedObservations: asos,
	}, nil
}

func observationFromProtoMessage(o *protobuf.Observation) (observation.Observation, error) {
	if o == nil {
		return observation.Observation{}, errors.New("Unable to extract a Observation value")
	}
	obs, err := observation.UnmarshalObservation(o.Value)
	if err != nil {
		return observation.Observation{}, errors.Errorf(`could not deserialize bytes as `+
			`observation.Observation: "%v" from 0x%x`, err, o.Value)
	}
	return obs, nil
}

func attributedObservationsFromProtoMessage(pbaos []*protobuf.AttributedObservation) ([]protocol.AttributedObservation, error) {
	if pbaos == nil {
		// note: we return an empty list instead of an error, because protobuf
		// represents empty list and nil as the same thing
		return []protocol.AttributedObservation{}, nil
	}
	aos := make([]protocol.AttributedObservation, 0, len(pbaos))
	for _, pbao := range pbaos {
		ao, err := attributedObservationFromProtoMessage(pbao)
		if err != nil {
			return nil, err
		}
		aos = append(aos, ao)
	}
	return aos, nil
}

func attributedObservationFromProtoMessage(ao *protobuf.AttributedObservation) (protocol.AttributedObservation, error) {
	if ao == nil {
		return protocol.AttributedObservation{}, errors.New("Unable to extract a AttributedObservation value")
	}

	o, err := observationFromProtoMessage(ao.Observation)
	if err != nil {
		return protocol.AttributedObservation{}, err
	}

	return protocol.AttributedObservation{
		o,
		commontypes.OracleID(ao.Observer),
	}, nil
}

func attestedReportOneFromProtoMessage(m *protobuf.AttestedReportOne) (protocol.AttestedReportOne, error) {
	if m == nil {
		return protocol.AttestedReportOne{}, errors.New("Unable to extract a AttestedReportOne value")
	}

	aos, err := attributedObservationsFromProtoMessage(m.AttributedObservations)
	if err != nil {
		return protocol.AttestedReportOne{}, err
	}
	sig := m.Signature
	if sig == nil {
		sig = []byte{}
	}

	return protocol.AttestedReportOne{aos, sig}, nil
}

func messageReportFromProtoMessage(m *protobuf.MessageReport) (protocol.MessageReport, error) {
	if m == nil {
		return protocol.MessageReport{}, errors.New("Unable to extract a MessageReport value")
	}

	report, err := attestedReportOneFromProtoMessage(m.Report)
	if err != nil {
		return protocol.MessageReport{}, err
	}

	return protocol.MessageReport{uint32(m.Epoch), uint8(m.Round), report}, nil
}

func attestedReportManyFromProtoMessage(m *protobuf.AttestedReportMany) (protocol.AttestedReportMany, error) {
	if m == nil {
		return protocol.AttestedReportMany{}, errors.New("Unable to extract a AttestedReportMany value")
	}

	signatures := make([][]byte, 0, len(m.Signatures))
	for _, sig := range m.Signatures {
		if sig == nil {
			sig = []byte{}
		}
		signatures = append(signatures, sig)
	}
	aos, err := attributedObservationsFromProtoMessage(m.AttributedObservations)
	if err != nil {
		return protocol.AttestedReportMany{}, err
	}

	return protocol.AttestedReportMany{aos, signatures}, nil
}

func messageFinalFromProtoMessage(m *protobuf.MessageFinal) (protocol.MessageFinal, error) {
	if m == nil {
		return protocol.MessageFinal{}, errors.New("Unable to extract a MessageFinal value")
	}
	report, err := attestedReportManyFromProtoMessage(m.Report)
	if err != nil {
		return protocol.MessageFinal{}, err
	}
	return protocol.MessageFinal{uint32(m.Epoch), uint8(m.Round), report}, nil
}

func messageFinalEchoFromProtoMessage(m *protobuf.MessageFinalEcho) (protocol.MessageFinalEcho, error) {
	if m == nil {
		return protocol.MessageFinalEcho{}, errors.New("Unable to extract a MessageFinalEcho value")
	}
	final, err := messageFinalFromProtoMessage(m.Final)
	if err != nil {
		return protocol.MessageFinalEcho{}, err
	}
	return protocol.MessageFinalEcho{MessageFinal: final}, nil
}

func attributedSignedObservationsFromProtoMessage(pbasos []*protobuf.AttributedSignedObservation) ([]protocol.AttributedSignedObservation, error) {
	if pbasos == nil {
		// note: we return an empty list instead of an error, because protobuf
		// represents empty list and nil as the same thing
		return []protocol.AttributedSignedObservation{}, nil
	}
	asos := make([]protocol.AttributedSignedObservation, 0, len(pbasos))
	for _, pbaso := range pbasos {
		aso, err := attributedSignedObservationFromProtoMessage(pbaso)
		if err != nil {
			return nil, err
		}
		asos = append(asos, aso)
	}
	return asos, nil
}

func attributedSignedObservationFromProtoMessage(m *protobuf.AttributedSignedObservation) (protocol.AttributedSignedObservation, error) {
	if m == nil {
		return protocol.AttributedSignedObservation{}, errors.New("Unable to extract an AttributedSignedObservation value")
	}

	signedObservation, err := signedObservationFromProtoMessage(m.SignedObservation)
	if err != nil {
		return protocol.AttributedSignedObservation{}, err
	}
	return protocol.AttributedSignedObservation{
		signedObservation,
		commontypes.OracleID(m.Observer),
	}, nil
}

func signedObservationsFromProtoMessage(pbsos []*protobuf.SignedObservation) ([]protocol.SignedObservation, error) {
	if pbsos == nil {
		// note: we return an empty list instead of an error, because protobuf
		// represents empty list and nil as the same thing
		return []protocol.SignedObservation{}, nil
	}
	sos := make([]protocol.SignedObservation, 0, len(pbsos))
	for _, pbso := range pbsos {
		so, err := signedObservationFromProtoMessage(pbso)
		if err != nil {
			return nil, err
		}
		sos = append(sos, so)
	}
	return sos, nil
}

func signedObservationFromProtoMessage(m *protobuf.SignedObservation) (protocol.SignedObservation, error) {
	if m == nil {
		return protocol.SignedObservation{}, errors.New("Unable to extract an SignedObservation value")
	}
	sig := m.Signature
	if sig == nil {
		sig = []byte{}
	}
	obs, err := observationFromProtoMessage(m.Observation)
	if err != nil {
		return protocol.SignedObservation{}, err
	}
	return protocol.SignedObservation{obs, sig}, nil
}

func attestedReportOneToProtoMessage(aro protocol.AttestedReportOne) *protobuf.AttestedReportOne {
	return &protobuf.AttestedReportOne{
		AttributedObservations: attributedObservationsToProtoMessage(aro.AttributedObservations),
		Signature:              aro.Signature,
	}
}

func attributedObservationsToProtoMessage(aos protocol.AttributedObservations) []*protobuf.AttributedObservation {
	result := []*protobuf.AttributedObservation{}
	for _, ao := range aos {
		result = append(result, &protobuf.AttributedObservation{
			Observation: &protobuf.Observation{Value: ao.Observation.Marshal()},
			Observer:    uint32(ao.Observer),
		})
	}

	return result
}

func finalToProtoMessage(v protocol.MessageFinal) *protobuf.MessageFinal {
	pm := &protobuf.MessageFinal{
		Epoch: uint64(v.Epoch),
		Round: uint64(v.Round),
		Report: &protobuf.AttestedReportMany{
			AttributedObservations: attributedObservationsToProtoMessage(v.Report.AttributedObservations),
			Signatures:             make([][]byte, len(v.Report.Signatures)),
		},
	}
	for i, sig := range v.Report.Signatures { //nolint:gosimple
		pm.Report.Signatures[i] = sig
	}
	return pm
}
